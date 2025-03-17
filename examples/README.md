# Examples

## math.js

This is a very simple example. The plugin has an adder and a subtractor function. This plugin is just for the sake of the example, since of course k6 can also add and subtract.

**math.plugin.ts** *(run by Deno/Bun)*

```ts file=math.plugin.ts
export function add(a: number, b: number): number {
  return a + b;
}

export function sub(a: number, b: number): number {
  return a - b;
}
```

The plugin only contains the two exported functions. These functions will be callable from the k6 test as asynchronous functions.


**math.js** *(run by k6)*
```js file=math.js
import { add, sub } from "k6/x/plugin/./math.plugin.ts";

export default async function () {
  console.log("2 + 2 =", await add(2, 2));
  console.log("3 - 2 =", await sub(3, 2));
}
```

As you can see, the plugin import path is the relative path prefixed with the xk6-plugin import prefix (`k6/x/plugin/`)

**screencast**

![math.js](../docs/math.svg)

The `deno` runtime is used to run the plugin if available. Otherwise, the `bun` runtime. This behavior can be overridden using the `K6_PLUGIN_RUNTIME` environment variable.

## dns.js

k6 currently does not have a DNS API. This example shows how to make the Node.js DNS API available to the k6 test using a plugin.

**dns.plugin.js** *(run by Deno/Bun)*

```ts file=dns.plugin.js
export { lookup, resolve } from "node:dns/promises";
```

The plugin is quite trivial, it exports two corresponding functions of the Node.js DNS API.

The exported functions in this case are asynchronous functions, but this does not cause any problems.

**dns.js** *(run by k6)*
```js file=dns.js
import { lookup, resolve } from "k6/x/plugin/./dns.plugin.js";

export default async function () {
  console.log("k6.io:", await lookup("k6.io"));
  console.log("grafana.com:", await resolve("grafana.com"));
}
```

The test script performs a lookup for the name `k6.io` and resolves the name `grafana.com`.

**screencast**

![dns.js](../docs/dns.svg)

## dump.js

This is a more realistic example. For debugging purposes, we want to dump the HTTP responses into separate files. k6 does not have a file writing feature, so the plugin uses the Node.js filesystem API to write files (and create the directory).

**dump.plugin.js** *(run by Deno/Bun)*

```ts file=dump.plugin.js
import { mkdir, writeFile } from "node:fs/promises";
import path from "node:path";

export async function dump(response, dir = "dump") {
  await mkdir(dir, { recursive: true });

  const filename = path.join(dir, new Date().toISOString() + ".json");

  await writeFile(filename, JSON.stringify(response, null, 2));
}
```

The plugin writes the response object to a JSON file named currenttime in the specified directory.

**dump.js** *(run by k6)*
```js file=dump.js
import { get } from "k6/http";
import { sleep } from "k6";
import { dump } from "k6/x/plugin/./dump.plugin.js";

export default async function () {
  const res = get("https://httpbin.test.k6.io/get");
  await dump(res);
  sleep(1);
}
```
Using the dump function is simple, you just need to pass the response object and call it asynchronously.

**screencast**

![dump.js](../docs/dump.svg)
