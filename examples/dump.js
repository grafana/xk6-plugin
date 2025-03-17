import { get } from "k6/http";
import { sleep } from "k6";
import { dump } from "k6/x/plugin/./dump.plugin.js";

export default async function () {
  const res = get("https://httpbin.test.k6.io/get");
  await dump(res);
  sleep(1);
}
