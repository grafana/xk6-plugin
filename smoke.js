import { add, sub } from "k6/x/plugin/./util.plugin.ts";

export default async function () {
  console.log("2 + 2 =", await add(2, 2));
  console.log("3 - 2 =", await sub(3, 2));
}
