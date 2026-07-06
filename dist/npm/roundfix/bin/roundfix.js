#!/usr/bin/env node
"use strict";

// Transparent launcher for the Roundfix CLI. It resolves the prebuilt binary
// for the current platform (installed as an optional dependency) and execs it,
// forwarding arguments, stdio, exit code, and terminating signal unchanged so
// Roundfix's stable exit-code contract survives the Node wrapper.

const { spawnSync } = require("node:child_process");

// process.platform + process.arch -> { pkg, bin }
const TARGETS = {
  "darwin arm64": { pkg: "@roundfix/cli-darwin-arm64", bin: "bin/roundfix" },
  "darwin x64": { pkg: "@roundfix/cli-darwin-x64", bin: "bin/roundfix" },
  "linux arm64": { pkg: "@roundfix/cli-linux-arm64", bin: "bin/roundfix" },
  "linux x64": { pkg: "@roundfix/cli-linux-x64", bin: "bin/roundfix" },
  "win32 x64": { pkg: "@roundfix/cli-win32-x64", bin: "bin/roundfix.exe" },
};

function fail(message) {
  process.stderr.write(`roundfix: ${message}\n`);
  process.exit(1);
}

function resolveBinary() {
  const key = `${process.platform} ${process.arch}`;
  const target = TARGETS[key];
  if (!target) {
    fail(
      `no prebuilt binary for ${key}. Supported platforms: ${Object.keys(TARGETS).join(", ")}.`,
    );
  }
  try {
    return require.resolve(`${target.pkg}/${target.bin}`);
  } catch (_err) {
    fail(
      `platform package "${target.pkg}" is not installed. Reinstall roundfix so ` +
        `the optional dependency for ${key} is fetched.`,
    );
  }
  return undefined;
}

const binary = resolveBinary();
const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  fail(`failed to launch ${binary}: ${result.error.message}`);
}
if (result.signal) {
  process.kill(process.pid, result.signal);
}
process.exit(result.status === null ? 1 : result.status);
