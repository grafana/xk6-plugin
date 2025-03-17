import { expect } from "jsr:@std/expect";
import { newRPCServer, RPCVersion, RPCRequest } from "./plugin.ts";

Deno.test("module loaded correctly", async () => {
  const server = await newRPCServer("./testdata/math.ts");

  const req = { jsonrpc: RPCVersion, id: 2, method: "add", params: [2, 3] } as RPCRequest;
  const res = await server.handle(req);

  expect(res.result).toBe(5);

  const obj = JSON.parse(await server.serve(JSON.stringify(req)));
  expect(obj.result).toBe(5);
});
