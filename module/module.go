// Package module provides a javascript module for xk6-plugin extension.
package module

import (
	"log/slog"
	"sync"

	"github.com/grafana/sobek"
	"github.com/grafana/xk6-plugin/internal/plugin"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/promises"
)

// ImportPath contains module's JavaScript import path.
const ImportPath = "k6/x/plugin"

// EnvRuntime contains environment variable for plugin runtime.
const EnvRuntime = "K6_PLUGIN_RUNTIME"

// New creates a new instance of the extension's JavaScript module.
func New(plugin string) modules.Module {
	return &rootModule{plugin: plugin}
}

func (root *rootModule) init(vu modules.VU) {
	var err error

	runtime, _ := vu.InitEnv().LookupEnv(EnvRuntime)
	logger := vu.InitEnv().Logger

	root.methods, root.in, root.out, err = plugin.Start(vu.Context(), root.plugin, runtime, logger)
	if err != nil {
		vu.InitEnv().Logger.Fatal(err)
	}

	go root.receive()
}

// rootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/deno` module instances for each VU.
type rootModule struct {
	plugin string

	methods  []string
	in       chan<- *plugin.RPCRequest
	out      <-chan *plugin.RPCResponse
	promises sync.Map

	once sync.Once
}

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (root *rootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	root.once.Do(func() {
		root.init(vu)
	})

	instance := &module{}

	instance.exports.Named = map[string]any{}

	for _, method := range root.methods {
		instance.exports.Named[method] = func(call sobek.FunctionCall) sobek.Value {
			return root.invoke(vu, method, call)
		}
	}

	return instance
}

func (root *rootModule) invoke(vu modules.VU, method string, call sobek.FunctionCall) sobek.Value {
	params := make([]any, 0, len(call.Arguments))

	for _, arg := range call.Arguments {
		params = append(params, arg.Export())
	}

	req := plugin.NewRequest(method, params...)

	prom, ret := newPromise(vu)

	root.promises.Store(req.ID, prom)

	root.in <- req

	return ret
}

func (root *rootModule) receive() {
	for res := range root.out {
		prom, found := root.promises.Load(res.ID)
		if !found {
			slog.Warn("Unknown call", "id", res.ID)

			continue
		}

		root.promises.Delete(res.ID)

		if p, ok := prom.(*promise); ok {
			p.process(res)
		} else {
			slog.Warn("Unknown call", "info", prom)
		}
	}
}

// module represents an instance of the JavaScript module for every VU.
type module struct {
	exports modules.Exports
}

// Exports is representation of ESM exports of a module.
func (mod *module) Exports() modules.Exports {
	return mod.exports
}

type promise struct {
	resolve func(any)
	reject  func(any)
	vu      modules.VU
}

func newPromise(vu modules.VU) (*promise, sobek.Value) {
	p, resolve, reject := promises.New(vu)

	return &promise{resolve: resolve, reject: reject, vu: vu}, vu.Runtime().ToValue(p)
}

func (p *promise) process(res *plugin.RPCResponse) {
	if res.Error != nil {
		p.reject(p.vu.Runtime().ToValue(res.Error))

		return
	}

	if res.Result == nil {
		p.resolve(sobek.Undefined())

		return
	}

	p.resolve(p.vu.Runtime().ToValue(res.Result))
}
