// Package plugin implements a Deno module plugin that can be called using the JSON-RPC protocol.
package plugin

import (
	"bytes"
	"context"
	_ "embed" // for plugin.ts
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

// NewRequest creates a new plugin method call request.
func NewRequest(method string, params ...any) *RPCRequest {
	return &RPCRequest{JSONRPC: rpcVersion, ID: xid.New().String(), Method: method, Params: params}
}

// Start starts the go routines required to use the plugin.
func Start(ctx context.Context, plugin string, runtime string, logger logrus.FieldLogger) (
	[]string, chan<- *RPCRequest, <-chan *RPCResponse, error,
) {
	cmd, config, err := startRPCServer(ctx, plugin, runtime)
	if err != nil {
		return nil, nil, nil, err
	}

	logger.Infof("%s started with %s", config.Runtime, plugin)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, config.URL, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	in, out := startPlugin(conn)

	go func() {
		<-ctx.Done()

		_ = conn.Close()
		_ = cmd.Wait()

		close(in)
		close(out)
	}()

	return config.Methods, in, out, nil
}

func startPlugin(conn *websocket.Conn) (chan *RPCRequest, chan *RPCResponse) {
	out := make(chan *RPCResponse)
	in := make(chan *RPCRequest)

	go func() {
		for {
			var resp *RPCResponse

			err := conn.ReadJSON(&resp)
			if err == nil {
				out <- resp

				continue
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				out <- newTransportError(err)
			}

			return
		}
	}()

	go func() {
		for {
			req := <-in

			err := conn.WriteJSON(req)
			if err == nil {
				continue
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				out <- newTransportError(err)
			}

			return
		}
	}()

	return in, out
}

//go:embed plugin.ts
var scripts []byte

// ErrNoRuntime will be returned if no runtime is available.
var ErrNoRuntime = errors.New("no plugin runtime available, 'deno' or 'bun' required")

const (
	runtimeBun  = "bun"
	runtimeDeno = "deno"
)

func findRuntime(module string, runtime string) (string, []string, error) {
	exe := runtime

	if len(runtime) == 0 {
		for _, runtime := range []string{runtimeDeno, runtimeBun} {
			abs, err := exec.LookPath(runtime)
			if err == nil {
				exe = abs

				break
			}
		}

		if len(exe) == 0 {
			return "", nil, ErrNoRuntime
		}
	}

	if strings.HasSuffix(exe, runtimeDeno) {
		return exe, []string{"run", "--allow-all", "-", module}, nil
	}

	return exe, []string{"run", "-", module}, nil
}

func startRPCServer(ctx context.Context, module string, runtime string) (*exec.Cmd, *rpcConfig, error) {
	exe, args, err := findRuntime(module, runtime)
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.CommandContext(ctx, exe, args...) //nolint:gosec

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(scripts)

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	var conf *rpcConfig

	if err := json.NewDecoder(stdout).Decode(&conf); err != nil {
		return nil, nil, err
	}

	return cmd, conf, nil
}

func newTransportError(err error) *RPCResponse {
	return &RPCResponse{
		JSONRPC: rpcVersion,
		Error: &RPCError{
			Code:    RPCErrTransport,
			Message: "Transport I/O error",
			Data:    err.Error(),
		},
	}
}
