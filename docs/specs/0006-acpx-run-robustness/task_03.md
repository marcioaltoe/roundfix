---
task: task_03
spec: 0006-acpx-run-robustness
status: completed
type: infra
complexity: low
---

# Task 03: Mitigate the acpx message-buffer limit, or evidence it upstream

## Overview

acpx 0.12.0's 10 MiB per-message buffer killed two finished turns in one day
(`-32603 Message buffer exceeded 10485760 bytes`, `acpxCode: RUNTIME`).
Investigate the pinned version for any configuration surface that raises or
disables the limit; apply it if real, and produce the upstream-ready report
either way. Evidence-driven: no invented flags. Verifiable by the recorded
investigation and, when a mitigation exists, a proven configuration.

## Requirements

1. MUST inspect the pinned acpx 0.12.0 for a message-buffer configuration
   surface — CLI help, `acpx config show` schema, and the installed
   package's source (`~/.asdf`/npm global tree) for the `10485760` constant
   and any override reading — and record exactly what was found with file or
   help references.
2. MUST, if a real surface exists: document the recommended setting in the
   README/skill guidance for orchestrator use and prove it with a local
   reproduction (a command that previously exceeded the limit, or the
   constant's behavior test) — no roundfix code changes unless an invocation
   flag is the surface, in which case the runner gains it behind the
   existing invocation builders with rig tests.
3. MUST, if no surface exists: write the upstream issue draft (title, body
   with the two incident reproductions, the buffer constant location, and
   the orchestrator use case) into this task's Result, and add the limit to
   the shipped docs as a known constraint with the "large docs-task
   payloads" trigger description.
4. MUST update `docs/_inbox/dogfood-findings-2.md` item 1 with the
   investigation outcome (mitigated with X / upstream-only).

## Subtasks

- [x] Source and help inspection of the pinned acpx for the buffer constant
- [x] Mitigation applied and proven, or upstream draft written
- [x] Docs guidance and findings update

## Acceptance Criteria

- [x] The Result names the exact finding: the buffer constant's location and
      whether any override exists, with references.
- [x] Mitigation path: the setting is documented and demonstrated; or
      upstream path: the issue draft is complete enough to file verbatim.
- [x] Findings-2 item 1 carries the dated outcome.

## Verification

- `rtk go test ./...` — expected: full suite passes (unchanged unless an
  invocation flag was added, in which case the agent-package rig covers it).
- `make verify` — expected: full gate passes.

## References

`_prd.md` → User Story 3; Core Feature 3; Decisions (evidence-driven).
`_techspec.md` → Build Order 3, Risks. Round-2 dogfood finding 1 (both
incident run ids and the exact error line).

## Result

Finding: upstream-only. The pinned `acpx@0.12.0` has a hard 10 MiB
queue-owner per-message buffer and no discovered CLI, config, or environment
override.

Evidence:

- Installed version: `rtk which acpx` resolved
  `/Users/marcio/.asdf/installs/nodejs/25.6.1/bin/acpx`; `rtk npm list -g
  acpx --depth=0` reported `acpx@0.12.0`.
- Help/config inspection: `rtk acpx --help`, `rtk acpx config --help`,
  `rtk acpx config show`, `rtk acpx config show --help`, and `rtk acpx config
  init --help` expose no message-buffer setting. `config show` contains
  defaults such as `queueMaxDepth`, `ttl`, `timeout`, `format`, and `agents`,
  but no buffer-size key.
- Installed package source:
  `/Users/marcio/.asdf/installs/nodejs/25.6.1/lib/node_modules/acpx/dist/output-CjdF5rHk.js`
  contains `const MAX_MESSAGE_BUFFER_SIZE = 10 * 1024 * 1024;` at line 1112,
  checks `buffer.length > 10485760` at line 1226, and emits `Message buffer
  exceeded ${MAX_MESSAGE_BUFFER_SIZE} bytes` at line 1228.
- Source-map location: the bundled source map points those lines to
  `../src/cli/queue/ipc.ts:43`, `:273`, and `:275`.
- Override search: searches for `messageBuffer`, `message_buffer`,
  `bufferSize`, `ACPX_.*BUFFER`, `process.env.*BUFFER`, `MAX_MESSAGE`, and
  `10485760` in the installed `dist` files found the hard constant and error
  path only. The queue-owner environment variables present in the package are
  payload/args variables, not buffer-size controls.
- Docs guidance: README, `.agents/skills/roundfix/SKILL.md`, and the embedded
  `skills/roundfix/SKILL.md` now record the known acpx `0.12.0` limit and the
  "large docs-task payloads" trigger.
- Findings update: `docs/_inbox/dogfood-findings-2.md` item 1 now carries the
  dated `2026-07-05` upstream-only outcome and names the second incident Run
  `run_20260705T152742Z_854a4b853b910140`.

Upstream issue draft:

Title: acpx queue owner should expose or avoid the 10 MiB message buffer limit

Body:

Expected: acpx should let orchestrators run long agent turns whose final ACP
message exceeds 10 MiB, either by exposing a documented buffer-size
configuration surface or by chunking/streaming queue-owner messages so a
single large message does not fail the completed turn.

Actual: acpx `0.12.0` fails the turn with `-32603 Message buffer exceeded
10485760 bytes` and exits 1 after the underlying agent has already completed
work.

### Reproduction steps

1. Install acpx `0.12.0` through npm on Node.js 25.6.1:
   `npm install -g acpx@0.12.0`.
2. Run a Roundfix docs/spec task through acpx-backed Codex where the agent
   prints or returns large skill/docs payloads.
3. Let the agent complete and verify the task.
4. Observe the acpx queue-owner error at the end of the turn.

### Environment

- acpx `0.12.0`
- Node.js `25.6.1`
- Installed package path:
  `/Users/marcio/.asdf/installs/nodejs/25.6.1/lib/node_modules/acpx`
- Orchestrator: Roundfix implement Runs using one acpx-backed Agent Session
  across task work.

### Evidence

- Incident 1: Roundfix Run `run_20260705T131519Z_d158ba8f114de398`, Spec
  `0003-dogfood-polish`, task `task_09`. The agent completed and verified the
  task, then acpx emitted `{"error":{"code":-32603,"message":"Message buffer
  exceeded 10485760 bytes","data":{"acpxCode":"RUNTIME"}}}` and exited 1.
- Incident 2: Roundfix Run `run_20260705T152742Z_854a4b853b910140`, Spec
  `0005-tui-cockpit`, task `task_07`. The same `-32603` buffer error recurred
  on another docs task touching large skill files.
- Hard limit location in the installed bundle:
  `/Users/marcio/.asdf/installs/nodejs/25.6.1/lib/node_modules/acpx/dist/output-CjdF5rHk.js:1112`
  defines `MAX_MESSAGE_BUFFER_SIZE = 10 * 1024 * 1024`;
  `:1226` checks `buffer.length > 10485760`; `:1228` emits the error.
- Source-map location: `../src/cli/queue/ipc.ts:43`, `:273`, and `:275`.
- `acpx --help`, `acpx config show`, `acpx config show --help`, and source
  searches found no CLI flag, config key, or environment variable to raise or
  disable this limit.

Hypothesis: the queue owner accumulates each socket response in memory and
rejects once `buffer.length` exceeds `MAX_MESSAGE_BUFFER_SIZE`, even if the
agent-side task has already completed successfully.

Request: please add a documented configuration surface to raise or disable the
per-message buffer, or change the queue-owner protocol to chunk/stream large
messages so orchestrators can safely run documentation-heavy tasks.
