package plugin

// The version of the JSON-RPC protocol.
const rpcVersion = "2.0"

// RPCRequest describes a JSON-RPC method call request.
type RPCRequest struct {
	// JSONRPC contains the JSON-RPC version number, its value must be "2.0"
	JSONRPC string `json:"jsonrpc,omitempty"`
	// ID is an identifier established by the Client that MUST contain a String, Number, or NULL value if included.
	// If it is not included it is assumed to be a notification.
	// The value SHOULD normally not be Null and Numbers SHOULD NOT contain fractional parts.
	ID string `json:"id,omitempty"`
	// Method contains the name of the method to be called.
	// Method names that begin with the word rpc followed by a period character (U+002E or ASCII 46)
	// are reserved for rpc-internal methods and extensions and MUST NOT be used for anything else.
	Method string `json:"method,omitempty"`
	// Params contains the parameters of the method call.
	// This member MAY be omitted.
	Params []any `json:"params,omitempty"`
}

type rpcConfig struct {
	URL     string   `json:"url"`
	Methods []string `json:"methods"`
	Runtime string   `json:"runtime"`
}

// RPCResponse describes a JSON-RPC method call response.
type RPCResponse struct {
	// JSONRPC contains the JSON-RPC version number, its value must be "2.0"
	JSONRPC string `json:"jsonrpc,omitempty"`
	// ID is an identifier and it MUST be the same as the value of the id member in the Request Object.
	// This member is REQUIRED.
	// If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request),
	// it MUST be Null.
	ID string `json:"id,omitempty"`
	// Result contains the return value of the method call.
	// This member is REQUIRED on success.
	// This member MUST NOT exist if there was an error invoking the method.
	Result any `json:"result,omitempty"`
	// Error contains an error that occurs during the method call.
	// This member MUST NOT exist if there was no error triggered during invocation.
	// The value for this member MUST be an RPCError object.
	Error *RPCError `json:"error,omitempty"`
}

// RPCError describes an error that occurs during a JSON-RPC method call.
type RPCError struct {
	// Code contains the error code.
	Code int `json:"code,omitempty"`
	// Message contains the description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message,omitempty"`
	// Data contains additional information about the error.
	// This may be omitted.
	// The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
	Data any `json:"data,omitempty"`
}

const (
	// RPCErrParse error code indicates a Parse error.
	// Invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	RPCErrParse = -32700
	// RPCErrInvalid error code indicates an Invalid Request.
	// The JSON sent is not a valid Request object.
	RPCErrInvalid = -32600
	// RPCErrNotFound error code indicates a Method not found.
	// The method does not exist / is not available.
	RPCErrNotFound = -32601
	// RPCErrParams error code indicates an Invalid params.
	// Invalid method parameter(s).
	RPCErrParams = -32602
	// RPCErrInternal error code indicates an Internal error.
	// Internal JSON-RPC error.
	RPCErrInternal = -32603
	// RPCErrTransport error code indicates a Transport error.
	RPCErrTransport = -32601
)
