---
spec: 0010-run-robustness
prd: _prd.md
created: 2026-07-06
---

# Run Robustness — Technical Spec

## Executive Summary

Three contained fixes to process and migration lifecycle. The consequential
one is detach: roundfix gains a second process model (self re-exec as a
session leader with a startup handshake), accepted because the alternative —
callers owning survival via `nohup` — is not a contract, and the full daemon
(work-plan item 3) is a much larger investment this stepping stone
de-risks. Config deprecation is a registry lookup ahead of strict parsing;
orphan reaping rides acpx's own session-close lifecycle rather than raw
process killing — acpx owns its tree, roundfix asks it to shut down.

## System Architecture

- `internal/config` — a `deprecatedKeys` table consulted before strict
  decoding: listed paths are stripped from the raw document, each emitting
  one stderr warning with the replacement; `KnownFields` stays strict for
  everything else. First entry: `resolve.concurrent` →
  `worktree.concurrency`.
- `internal/agent` — `CloseSession` already exists; gains
  `ListRoundfixSessions` (parse `acpx <agent> sessions list` output for
  `roundfix-` prefixed names) so callers can close by discovery, not just
  by derivation.
- `internal/cli` — `stop --force` closes the target Run's sessions
  (Run-level `roundfix-<id>` plus per-Task `roundfix-<id>-<task>` matches)
  after the cancel; the implement preflight sweep closes discovered
  roundfix-sessions whose run ids resolve to terminal Runs; `--detach` flag
  on the three operational commands with the re-exec/handshake flow.
- No store, journal, engine, or TUI changes.

## Implementation Design

### Config deprecation (ADR-0027)

Parse the YAML to a generic node tree first; for each registered deprecated
path present, remove the node and record a warning
(`config: resolve.concurrent is deprecated and ignored; use
worktree.concurrency`); then strict-decode the cleaned document exactly as
today. Warnings print once per load on stderr, never fail, and are covered
for both User and Project Config files.

### Session-close reaping

- `stop --force` ordering becomes: cooperative cancel → `sessions close`
  for the Run-level session and any per-Task sessions discovered by prefix
  match → database completion and lock release. Close failures are stderr
  notes (`could not close session <name>: <reason>`), never fatal
  (acpx's TTL remains the last resort).
- The implement preflight sweep, after the existing worktree reaping:
  discover `roundfix-` sessions via `ListRoundfixSessions`, map each
  embedded run id against the Run Database, and close those whose Runs are
  terminal — one report line per closed session. Unknown or non-roundfix
  sessions are never touched.
- Product code never issues raw process kills; the gated real-acpx
  integration test asserts the process tree is gone after close.

### Detach (ADR-0028)

```
roundfix <implement|resolve|watch> ... --detach
```

- The caller validates flags cheaply (flag parse + `--detach` conflicts:
  `--interactive` rejected; `--no-input` implied), creates the handshake
  pipe, and re-execs its own binary with the same arguments minus
  `--detach` plus an internal handshake fd/env marker, using
  `SysProcAttr{Setsid: true}` and stdio bound to the per-Run console log
  path (`<artifact_dir>/runs/<run-id>/console.log` — opened by the child
  after the run id exists; before that, a temp log renamed into place).
- The child runs the normal command path; immediately after `CreateRun` it
  writes `run_id\tconsole_log_path\n` to the handshake and continues.
- The caller reads the handshake (bounded wait): on success prints exactly
  `Run detached: <run-id>`, `Console log: <path>`,
  `Follow: roundfix attach <run-id>`, `Stop: roundfix stop <run-id>`, exit
  0. If the child exits before the handshake (Preflight failure), the
  caller relays the child's captured stderr and exit code verbatim — a
  detached launch never hides a Preflight message.
- The detached child behaves as a normal non-TTY Run in every other way
  (journal, worktrees, integration, outcomes); Ctrl-C in the caller after
  detach is irrelevant by construction.

## Coverage Map

- Story 1 → deprecation registry (ADR-0027, finding R3-5)
- Story 2 → force-stop session close (finding R3-6)
- Story 3 → preflight sweep session close (finding R3-6)
- Story 4 → detach re-exec + handshake (ADR-0028, finding R3-7)
- Story 5 → console log + deterministic detach stdout

## Integration Points

acpx only (`sessions list`, `sessions close` — both already in the
invocation grammar); the detach re-exec uses the binary's own path
(`os.Executable`). No new external systems.

## Testing Approach

Config: table tests over documents carrying deprecated keys (User and
Project files, mixed with valid keys), asserting the cleaned decode, the
exact warning line, and unknown keys still failing. Sessions: fake-acpx rig
scripts for `sessions list` parsing and close invocations per selector;
CLI tests assert the stop --force ordering and sweep report lines; the
gated real-acpx test (ROUNDFIX_REAL_ACPX=1) gains the orphan assertion
(spawn, kill parent, force-stop, assert no adapter processes). Detach:
buffer-captured caller tests over the handshake protocol (success shape,
pre-handshake Preflight relay with exit codes), plus one integration-style
test that detaches, kills the caller process, and asserts the child reaches
a journaled terminal outcome and remains attachable. This Spec's own
implement Run is the 0009 concurrency live test (three independent wave-1
tasks at concurrency 2).

## Build Order

1. Config deprecation registry with the `resolve.concurrent` entry (no
   deps)
2. Session discovery and close-based reaping in stop --force and the
   preflight sweep (no deps)
3. Detach mode: flag, re-exec, handshake, console log, relay semantics (no
   deps)
4. Docs and skill sync (depends on: 1, 2, 3)

## Risks & Considerations

- The handshake must be crash-safe: a child dying pre-handshake leaves the
  caller with the child's stderr and exit code, never a hang (bounded
  read).
- Session discovery must be conservative: only `roundfix-` prefixed names
  with run ids that resolve in the Run Database are ever closed.
- Deprecation stripping runs before strict decode — the registry is the
  only bypass of `KnownFields`, keeping typo protection intact.
- Detached console logs are the Run's only record; the log open must
  precede any output-worthy work (temp-then-rename covers the pre-run-id
  window).

## Decisions

- Registry-based warn-and-ignore; unknown keys still strict. See ADR-0027.
- Reaping via acpx session close, discovery-scoped to roundfix-named
  sessions of terminal Runs; no raw kills in product code.
- Detach = self re-exec + Setsid + handshake; `--interactive` conflicts;
  Preflight failures relay verbatim. See ADR-0028.
- Glossary gains Detached Run (committed with the PRD).
