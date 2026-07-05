---
spec: 0002-acpx-migration
prd: _prd.md
created: 2026-07-05
---

# ACPX Migration — Technical Spec

## Executive Summary

The agent layer swaps its implementation behind an existing seam: everything above `agent.Runner` (CLI, daemon engines, watch, TUI, store) is untouched, while the ~1,150-line hand-rolled ACP client (`coder/acp-go-sdk`: session plumbing, permission answering, fs jailing, teardown) is replaced by process orchestration of the acpx CLI. The primary trade-off this design accepts is **taking a declared-alpha external CLI as a hard dependency in exchange for deleting the largest subsystem Roundfix owns** (ADR-0017); containment is a pinned version verified at Preflight Validation, acpx's documented-stable exit codes and `--json-strict` wire-format guarantees, and the fact that a rollback is a git revert behind one interface. Sessions become explicit: one named Agent Session per Run (ADR-0018), created on first Agent work, resumed by acpx across runtime crashes, and closed on every terminal path.

## System Architecture

**Materially changed: `internal/agent` only.**

- A new acpx-backed `Runner` implementation replaces the ACP SDK runner. `RuntimeSpec` keeps its role (runtime names `codex`/`claude`/`opencode`, command override, model, full-access mode); prompt builders, `SettleAssignedIssues`, and the stream→Run Event conversion survive with the conversion's input switching from SDK types to raw JSON-RPC lines.
- The `Runner` interface gains explicit session lifecycle (`EndSession`) and `ExecuteRequest` gains the Agent Session reference. All fakes in existing tests implement the extended interface trivially.

**Touched at the edges:**

- `internal/cli` — derives the Agent Session name from the run id, passes it through both engine paths, and closes the session on every terminal path (review and spec); Preflight gains the acpx presence/version check with an actionable install message.
- `internal/daemon` — `CyclePlan`/`TaskPlan` carry the session reference through to `ExecuteRequest`; no engine logic changes.
- `internal/preflight` — nothing (the acpx check lives beside the existing runtime probe call sites in the CLI).

**Deleted:** the SDK client, its fs/terminal/permission handlers, and the `coder/acp-go-sdk` dependency from the Go module. acpx owns permission answering (blanket approval for parity), the cwd-scoped fs/terminal jail, crash respawn (`session/resume` → `session/load` → `session/new` with desired-mode replay), and queue ownership.

## Implementation Design

### Interfaces

```go
// internal/agent — the seam stays; sessions become explicit
type SessionRef struct{ Name string } // "roundfix-<run-id>"

type Runner interface {
    Probe(ctx context.Context, runtime RuntimeSpec) error
    Run(ctx context.Context, req ExecuteRequest, sink runevent.Sink) (ExecuteResult, error)
    EndSession(ctx context.Context, runtime RuntimeSpec, session SessionRef) error
}
// ExecuteRequest gains `Session SessionRef`; all other fields unchanged.
```

### acpx invocation mapping

The whole integration surface, explicit so tasks need no other reference (acpx pinned at the version current when implementation starts; constant in `internal/agent`):

- **Ensure (first Agent work of a Run, idempotent):** `acpx <agent> sessions ensure --name <session> --cwd <workdir>` — prompt commands never auto-create; exit 4 from a prompt means the ensure was skipped (infrastructure bug, not a user error).
- **Full access (after ensure, when opted in, ADR-0011):** `acpx <agent> set-mode <FullAccessMode> -s <session>` with the existing verbatim ids (`full-access` codex, `bypassPermissions` claude; OpenCode keeps defaults). Codex's `danger-full-access` sandbox preset is applied through acpx's session config options (`acpx codex set <key> <value>`; the exact key is verified against the pinned adapter during implementation — acpx persists desired mode and config options and replays them after crash respawn).
- **Prompt (per Work Item):** `acpx <agent> prompt -s <session> --cwd <workdir> --format json --json-strict --approve-all [--model <id>] -f -` with the built prompt on stdin. stdout is then one raw ACP JSON-RPC message per line, no envelope, no renamed keys.
- **Cancel (context canceled mid-prompt):** `acpx <agent> cancel -s <session>` (cooperative `session/cancel`), then process kill as fallback — mirroring acpx's own SIGINT behavior.
- **Close (every terminal Run outcome):** `acpx <agent> sessions close -s <session>` (best effort; acpx's 300 s idle TTL is the backstop for crashed Roundfix processes).
- **Command override:** the stdio escape hatch maps to the global `--agent "<command>"` in place of the adapter name.
- **Probe:** acpx binary on PATH and `acpx --version` equal to the pin; mismatch or absence becomes one Preflight Validation message carrying `npm install -g acpx@<pin>`.
- The runtime name is always passed explicitly — acpx's implicit codex default is never relied on.

### Stream and journaling

Each stdout line is parsed as JSON-RPC: `session/update` notifications feed the existing conversion into Run Events with the raw line as payload (ADR-0008 unchanged); the `session/prompt` response line carries the stop reason that becomes `ExecuteResult`; every line is also appended to the existing agent log path. Non-JSON output cannot appear on stdout under `--json-strict` (documented hard rule); stderr remains diagnostics and is captured only into error context.

### Exit-code mapping

acpx documents these as stable: `0` success → result from the stream; `1` agent/protocol/runtime error → Batch failure (ADR-0010: settle, continue); `3` `--timeout` expiry → Batch failure with timeout reason; `5` all permissions denied → Batch failure (should be unreachable under `--approve-all`; journaled loudly); `2` usage error and `4` missing session → infrastructure errors that fail the Run (they indicate a Roundfix bug, not user conditions); `130` → Stop Request semantics.

### Data Models

None. No store schema change, no config keys added or changed, no Run Event kinds added. The pinned acpx version is a code constant, not configuration — upgrades are deliberate commits (PRD open question resolved this way).

### API Contracts

No command, flag, stdout, or exit-code changes. New Preflight Validation failure messages: acpx missing / version mismatch (with the install command), which follow the existing one-actionable-message convention.

## Coverage Map

- Story 1 (warm Work Items) → acpx runner prompt path + Agent Session ensure-once wiring
- Story 2 (crash resume) → acpx respawn/resume (verified by an induced-kill integration test at the runner seam)
- Story 3 (cooperative stop) → cancel mapping + existing interrupt handling
- Story 4 (full access) → set-mode/config mapping (ADR-0011 parity)
- Story 5 (command override) → `--agent` mapping in `RuntimeSpec`
- Story 6 (identical Attach) → stream parsing into the unchanged conversion + journal
- Story 7 (actionable preflight) → Probe + version pin + install message
- Goal "layer deleted" → cutover step (SDK removal, `go.mod` shrink)

## Integration Points

- **acpx CLI** (process boundary; the only integration surface — no IPC, no Node embedding). Pinned version.
- **Node ≥ 22.13 / npx** — environment prerequisite surfaced at Preflight; default adapters launch via `npx -y`, and the docs update recommends configuring direct binaries in acpx config for latency-sensitive setups.
- **ACP Runtimes** — codex-acp, claude-agent-acp, opencode as acpx built-in adapters; no direct protocol contact remains in Go.

## Testing Approach

- **Runner unit tests via the stdlib helper-process pattern**: the test binary re-execs itself as a scripted fake `acpx` (env-selected script: canned NDJSON on stdout, chosen exit code, arg capture to a file). Zero new dependencies, no real Node in unit tests. Covers: command construction for every mapping above, NDJSON→Run Event conversion with raw-payload equality, exit-code mapping table, cancel-on-context, ensure-once-then-prompt sequencing, EndSession on terminal paths.
- **Existing suites unchanged**: CLI and daemon tests already run against `Runner` fakes through the collaborator seams; they extend only for the added `EndSession`/`Session` fields. The full suite passing unchanged is the parity metric.
- **One gated integration test** (env-guarded, skipped by default) drives a real pinned acpx + a trivial `--agent` echo server for the crash-resume and cooperative-cancel stories on machines that have Node.
- No new test seams; the existing `newEngineCollaborators` and runner fakes remain the attachment points.

## Build Order

1. **acpx invocation core** — command construction, NDJSON stream parsing into the existing conversion, exit-code mapping, agent log writing; helper-process test rig.
2. **Session lifecycle** (depends on: 1) — `SessionRef`, ensure-once, set-mode/full-access mapping, cancel, `EndSession`; `ExecuteRequest.Session`.
3. **Probe and Preflight Validation** (depends on: 1) — pinned-version check, actionable install/upgrade messages, wired at the existing probe call sites.
4. **Run wiring** (depends on: 2, 3) — session name derived from the run id in both engine paths; close on every terminal path (review and spec); journal the session lifecycle through existing event kinds.
5. **Cutover** (depends on: 2, 3, 4) — delete the SDK runner and handlers, drop `coder/acp-go-sdk` from the module, full parity gate plus the gated real-acpx integration test.
6. **Docs and skill update** (depends on: 5) — the Roundfix skill and README name the acpx pin and Node prerequisite; handoff work-plan item 2 recorded as done.

## Risks & Considerations

- **Alpha churn**: acpx says interfaces may change. Containment: exact version pin verified at preflight, reliance only on documented-stable surfaces (exit codes, `--json-strict` wire rules), upgrades as deliberate commits with the gated integration test as the canary.
- **Codex sandbox config key** is the one mapping not pinned by documentation; step 2 verifies it against the pinned adapter and records the result in the task's evidence. Fallback: full-access still applies the mode id; the sandbox preset degrades gracefully with a journaled warning.
- **Orphan owner processes** if Roundfix dies before `EndSession`: acpx's idle TTL (300 s default) is the backstop; the close is best-effort by design.
- **Journal volume**: journaling reads from stdout, never from acpx's rotating on-disk stream (5×64 MB cap), so long Runs lose nothing.
- **One in-flight prompt per session** matches today's sequential execution; future parallel work (work-plan item 4) will need one named session per worktree, which ADR-0018's naming scheme already permits (name derives from the Run, extensible per worktree).
- **First-run latency** of `npx -y` adapter installs: documented recommendation to configure direct binaries; not a correctness risk.

## Decisions

- Hard cutover behind `agent.Runner`; pinned acpx; Node accepted. See ADR-0017.
- One Agent Session per Run, ensure-once, closed on every terminal path. See ADR-0018.
- Journal source is the `--json-strict` stdout stream; acpx's on-disk session files are never read.
- Permission parity via `--approve-all` + verbatim `set-mode` ids (ADR-0011); the policy engine stays work-plan item 8.
- Exit codes 2 and 4 are treated as Roundfix bugs (infrastructure failure), not user conditions.
- The pin is a code constant; version upgrades are commits, not configuration.
- acpx flows are never used; the Daemon remains the only orchestrator.
