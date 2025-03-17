# xk6-plugin

**Create k6 plugins using JavaScript and Node.js**

The **xk6-plugin** k6 extension makes the Node.js runtime available to k6 tests using k6 plugins. Plugins are executed by a modern Node.js compatible runtime ([Deno](https://deno.com) or [Bun](https://bun.sh)), so they can use the entire Node.js ecosystem.

In simple terms, a plugin is nothing more than a JavaScript (or TypeScript) module executed on the Node.js runtime, whose exported functions are available from k6 tests. Both [Deno](https://deno.com) and [Bun](https://bun.sh) runtimes are supported.


![screencast](docs/readme.svg)

## Features

- the entire Node.js ecosystem is available, millions of packages
- plugins can be written in TypeScript/JavaScript
- modern Node.js compatible JavaScript runtime (Deno or Bun)
- no go language skills required
- plugin development does not require a custom k6 build
- plugins are ESM modules, no need to learn a new API
- plugins can be tested without k6
- k6 tests and plugins can be written in the same language

## Examples

Examples with explanations can be found in the [examples](examples) directory.

## Download

You can download pre-built k6 binaries from [Releases](https://github.com/grafana/xk6-plugin/releases) page.

> [!IMPORTANT]
> [Deno](https://deno.com) or [Bun](https://bun.sh) must also be installed.

## How It Works

xk6-plugin automatically finds plugins based on file names. Files with names ending in `.plugin.ts`, `.plugin.js`, `.plugin.mjs` or `.plugin.cjs` are considered plugins. A k6 extension registration is created for each plugin by prefixing the plugin file name with `k6/x/plugin/./`. For example, the `tools.plugin.js` plugin in the current directory will be available in k6 with the following import path:

```
k6/x/plugin/./tools.plugin.js
```

A JavaScript runtime (Deno or Bun) is started for each plugin and the exported functions of the plugin will be available as asynchronous functions for the k6 test. The communication between the k6 process and the plugin JavaScript runtime is done via WebSocket. A JSON-RPC is implemented on top of the WebSocket communication.

![architecture](docs/plugin.dark.svg#gh-dark-mode-only)
![architecture](docs/plugin.light.svg#gh-light-mode-only)

The `deno` runtime is used to run the plugin if available. Otherwise, the `bun` runtime. This behavior can be overridden using the `K6_PLUGIN_RUNTIME` environment variable.

## Quick start

1. Download the archive from the [Releases page](https://github.com/grafana/xk6-plugin/releases) for your operating system and extract the `k6` executable from it.

2. Download and install [Deno](https://deno.com) and/or [Bun](https://bun.sh).

3. Create a plugin, let's say `hello.plugin.js`:

   ```js file=hello.plugin.js
   export function hello() {
    return "Hello, World!"
   }
   ```

4. Create a k6 test script that uses this plugin, let's say `script.js`:

   ```js file=script.js
   import { hello } from "k6/x/plugin/./hello.plugin.js"

   export default async function() {
    console.log(await hello())
   }
   ```

5. Run the k6 test script:

   ```bash
   k6 run -q --no-summary script.js
   ```

   You will get something like this output:

   ```
   INFO[0000] deno 2.2.2 started with ./hello.plugin.js    
   INFO[0000] Hello, World!                                 source=console
   ```

🎉 Congratulations, you have created and used your first k6 plugin!

## Limitations

- `./` in the path refers to the current directory and not the directory containing the script. This is due to the fact that there is currently no way for the k6 extension to find out the name of the referring test script.

## Contribute

If you would like to contribute, start by reading [Contributing Guidelines](docs/CONTRIBUTING.md).
