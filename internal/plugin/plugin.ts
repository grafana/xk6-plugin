/**
 * The version of the JSON-RPC protocol.
 */
export const RPCVersion = "2.0";

/**
 * A rpc call is represented by sending a Request object to a Server.
 */
export interface RPCRequest {
  /**
   * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
   */
  jsonrpc: "2.0";
  /**
   * An identifier established by the Client that MUST contain a String, Number, or NULL value if included.
   * If it is not included it is assumed to be a notification.
   * The value SHOULD normally not be Null and Numbers SHOULD NOT contain fractional parts.
   */
  id: number | string;
  /**
   * A String containing the name of the method to be invoked.
   * Method names that begin with the word rpc followed by a period character (U+002E or ASCII 46)
   * are reserved for rpc-internal methods and extensions and MUST NOT be used for anything else.
   */
  method: string;
  /**
   * A Structured value that holds the parameter values to be used during the invocation of the method.
   * This member MAY be omitted.
   */
  params?: unknown[];
}

/**
 * Type guard for RPCRequest.
 * @param arg value to check
 * @returns true if the argument is a RPCRequest
 */
export function isRPCRequest(arg: unknown): arg is RPCRequest {
  if (typeof arg != "object" || arg == null) {
    return false;
  }

  if (!("jsonrpc" in arg) || arg["jsonrpc"] != RPCVersion) {
    return false;
  }

  if (!("id" in arg)) {
    return false;
  }

  if (!("method" in arg)) {
    return false;
  }

  if (
    "params" in arg &&
    typeof arg["params"] == "object" &&
    !Array.isArray(arg["params"])
  ) {
    return false;
  }

  return true;
}

/**
 * When a rpc call is made, the Server MUST reply with a RPCResponse, except for in the case of Notifications.
 * The RPCResponse is expressed as a single JSON Object.
 */
export interface RPCResponse {
  /**
   * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
   */
  jsonrpc: "2.0";
  /**
   * It MUST be the same as the value of the id member in the Request Object.
   * This member is REQUIRED.
   * If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request),
   * it MUST be Null.
   */
  id: number | string;
  /**
   * The value of this member is determined by the method invoked on the Server.
   * This member is REQUIRED on success.
   * This member MUST NOT exist if there was an error invoking the method.
   */
  result?: unknown;
  /**
   * This member is REQUIRED on error.
   * This member MUST NOT exist if there was no error triggered during invocation.
   * The value for this member MUST be an RPCError object.
   */
  error?: RPCError;
}

/**
 * When a rpc call encounters an error, the RPCResponse Object MUST contain the error member with the
 * value that is RPCError object.
 */
export interface RPCError {
  /**
   * A Number that indicates the error type that occurred.
   * This MUST be an integer.
   */
  code: number;
  /**
   * A String providing a short description of the error.
   * The message SHOULD be limited to a concise single sentence.
   */
  message: string;
  /**
   * A Primitive or Structured value that contains additional information about the error.
   * This may be omitted.
   * The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
   */
  data: unknown;
}

/**
 * Parse error
 * Invalid JSON was received by the server.
 * An error occurred on the server while parsing the JSON text.
 */
export const ErrParse = -32700;
/**
 * Invalid Request
 * The JSON sent is not a valid Request object.
 */
export const ErrInvalid = -32600;
/**
 * Method not found
 * The method does not exist / is not available.
 */
export const ErrNotFound = -32601;
/**
 * Invalid params
 * Invalid method parameter(s).
 */
export const ErrParams = -32602;
/**
 * Internal error
 * Internal JSON-RPC error.
 */
export const ErrInternal = -32603;

// -32000 to -32099 	Server error 	Reserved for implementation-defined server-errors.

interface RPCConfig {
  url: string;
  methods: string[];
  runtime: string;
}

export async function newRPCServer(filename: string): Promise<RPCServer> {
  return new RPCServer(await import(filename));
}

type ModuleType = Record<string, unknown>;

export class RPCServer {
  module: ModuleType;
  constructor(instance: ModuleType) {
    this.module = instance;
  }

  // deno-lint-ignore ban-types
  method(name: string): Function | null {
    if (!(name in this.module)) {
      return null;
    }

    const fn = this.module[name];

    if (typeof fn != "function") {
      return null;
    }

    return fn;
  }

  async handle(request: RPCRequest): Promise<RPCResponse> {
    const fn = this.method(request.method);

    if (fn == null) {
      return {
        jsonrpc: RPCVersion,
        id: request.id,
        error: newRPCError(ErrNotFound, "No such method"),
      } as RPCResponse;
    }

    const result = await fn.apply(this.module, request.params);
    return {
      jsonrpc: request.jsonrpc,
      id: request.id,
      result,
    } as RPCResponse;
  }

  async serve(input: string): Promise<string> {
    let req: unknown;

    try {
      req = JSON.parse(input);
    } catch (e) {
      return JSON.stringify({
        jsonrpc: RPCVersion,
        error: newRPCError(ErrParse, e),
      } as RPCResponse);
    }

    if (!isRPCRequest(req)) {
      return JSON.stringify({
        jsonrpc: RPCVersion,
        error: newRPCError(ErrInvalid, "Not a request object"),
      } as RPCResponse);
    }

    return JSON.stringify(await this.handle(req));
  }

  methods(): string[] {
    return Object.getOwnPropertyNames(this.module).filter(
      (name) => this.method(name) != null
    );
  }
}

function newRPCError(code: number, err: unknown): RPCError {
  return { code, message: err } as RPCError;
}

async function mainDeno() {
  const server = await newRPCServer(Deno.args[0]);

  const ac = new AbortController();

  const handler = (request: Request): Response => {
    if (request.headers.get("upgrade") != "websocket") {
      return new Response(null, { status: 501 });
    }

    const { socket, response } = Deno.upgradeWebSocket(request);

    socket.onmessage = async (event) => {
      socket.send(await server.serve(event.data));
    };

    socket.onclose = () => ac.abort();
    socket.onerror = (e) => {
      const event = e as ErrorEvent;
      if (event.error && event.error instanceof Deno.errors.UnexpectedEof) {
        return;
      }

      if (event.error == "Error: Unexpected EOF") {
        return;
      }

      console.error("ERROR:", event.error);
    };

    return response;
  };

  Deno.serve({
    onListen: ({ port, hostname }) => {
      const config: RPCConfig = {
        url: `ws://${hostname}:${port}`,
        methods: server.methods(),
        runtime: `deno ${Deno.version.deno}`,
      };

      console.log(JSON.stringify(config));
    },
    port: 0,
    hostname: "127.0.0.1",
    handler,
    signal: ac.signal,
  });
}

async function mainBun() {
  // @ts-ignore: Deno Visual Studio Code Extension
  const server = await newRPCServer(Bun.argv[2]);

  // @ts-ignore: Deno Visual Studio Code Extension
  const srv = Bun.serve({
    port: 0,
    hostname: "127.0.0.1",
    fetch(req: Request, server: { upgrade: (arg: Request) => unknown }) {
      if (server.upgrade(req)) {
        return;
      }

      return new Response(null, { status: 501 });
    },
    websocket: {
      async message(ws: { send: (arg: string) => void }, message: string) {
        ws.send(await server.serve(message));
      },
      error(_ws: unknown, error: Error) {
        console.error("ERROR:", error);
      },
      close(_ws: unknown, _code: number, _reason: string) {
        srv.stop();
      },
    },
  });

  const config: RPCConfig = {
    url: `ws://${srv.hostname}:${srv.port}`,
    methods: server.methods(),
    // @ts-ignore: Deno Visual Studio Code Extension
    runtime: `bun ${Bun.version}`,
  };

  console.log(JSON.stringify(config));
}

if (import.meta.main) {
  // @ts-ignore: Detect Bun
  // deno-lint-ignore no-process-global
  if (process.versions.bun) {
    mainBun();
  } else {
    mainDeno();
  }
}
