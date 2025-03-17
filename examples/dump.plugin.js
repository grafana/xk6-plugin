import { mkdir, writeFile } from "node:fs/promises";
import path from "node:path";

export async function dump(response, dir = "dump") {
  await mkdir(dir, { recursive: true });

  const filename = path.join(dir, new Date().toISOString() + ".json");

  await writeFile(filename, JSON.stringify(response, null, 2));
}
