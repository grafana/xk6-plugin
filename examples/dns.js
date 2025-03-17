import { lookup, resolve } from "k6/x/plugin/./dns.plugin.js";

export default async function () {
  console.log("k6.io:", await lookup("k6.io"));
  console.log("grafana.com:", await resolve("grafana.com"));
}
