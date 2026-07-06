---
task: task_02
spec: 0012-npm-distribution
status: pending
type: infra
complexity: medium
---

# Task 02: Launcher package and pass-through bin shim

## Overview

Build the `roundfix` launcher package whose `bin` shim resolves the correct
per-platform binary package and execs it, forwarding arguments, exit code, and
terminating signal verbatim. The shim is the only piece users invoke through
`npx`/`bunx`/global installs, so its transparency is what keeps Roundfix's
exit-code contract intact through Node.

## Requirements

1. MUST add the `roundfix` launcher `package.json` (`type: commonjs`) declaring
   the per-platform packages as `optionalDependencies`, with a `bin` entry
   pointing at the shim.
2. MUST implement the shim to map `process.platform`/`process.arch` to the
   `@roundfix/cli-<os>-<arch>` package via the task_01 mapping, resolve its
   binary, and exec it with all passthrough arguments.
3. MUST propagate the child's exit code and terminating signal unchanged and MUST
   add nothing to stdout/stderr on success.
4. MUST fail with a clear message only when no matching platform package is
   installed (unsupported platform / reinstall guidance).

## Subtasks

- [ ] Launcher `package.json` with `optionalDependencies` and `bin`
- [ ] Shim: platform→package resolution using the task_01 mapping
- [ ] Exec with `stdio: 'inherit'`; propagate exit code and signal
- [ ] Missing-platform-package error path
- [ ] Shim test: exit code passthrough for success and a known failure

## Acceptance Criteria

- [ ] `roundfix <args>` through the shim runs the platform binary and returns its exit code unchanged for a success and a known non-zero failure.
- [ ] The shim writes nothing of its own to stdout on success (output is the binary's verbatim).
- [ ] A terminating signal to the child is reflected in the shim's termination.
- [ ] With no platform package present, the shim errors with a clear unsupported-platform message.

## Verification

- `node dist/npm/roundfix/bin/roundfix --version` (with the local platform package linked) — expected: prints the binary's version, exit 0.
- The shim test harness — expected: asserts exit-code passthrough for success and failure.

## References

`_prd.md` → User Stories 1-2, 6; Core Features 1, 4. `_techspec.md` → Launcher
shim, Build Order 2. ADR-0031.
