---
spec: 0016-worktree-bootstrap
prd: _prd.md
created: 2026-07-06
---

# Worktree Bootstrap — Technical Spec

## Executive Summary

Add one preparation step to worktree creation: a configured command run in the
new worktree after `worktree.copy`, before Agent work and Verification. The
primary trade-off is where the command runs and how failure surfaces: bootstrap
executes inside the isolated worktree (so it prepares that worktree, not the
user checkout) and a non-zero exit is a distinct terminal signal — the Run fails
for a Run Worktree, the Task settles failed for a Task Worktree — never a
Verification-shaped failure. Roundfix owns invocation and bounding (a timeout),
not semantics: dependency install and database strategy live in the command, and
shared-state projects run at `worktree.concurrency: 1` so bootstrap runs once on
the reused Run Worktree. The change is contained to the worktree lifecycle and
config; no store, journal, or scheduler change.

## System Architecture

- `internal/config` — a `worktree.bootstrap` string (empty = no bootstrap) and a
  `worktree.bootstrap_timeout` duration (default `10m`), strict-decoded with the
  usual Project > User > builtin precedence.
- `internal/worktree` — worktree creation already applies the `CopyList`
  (`worktree.copy`); it gains a bootstrap step that runs the configured command
  in the new worktree root immediately after copying and before the worktree is
  handed back for Agent work. Returns a distinct bootstrap error on non-zero
  exit or timeout.
- `internal/cli` / `internal/daemon` — the Run Worktree creation path maps a
  bootstrap failure to a Run-ending outcome with a clear message; the Task
  Worktree creation path (concurrent Tasks) maps a bootstrap failure to that
  Task settling failed with a bootstrap-failed reason. Bootstrap output is
  journaled and written to stderr.
- No scheduler, integration, store, or journal-format change.

## Implementation Design

### Interfaces

```go
// internal/worktree — bootstrap runs after CopyList placement.
type BootstrapSpec struct {
    Command string        // worktree.bootstrap; empty = skip
    Timeout time.Duration // worktree.bootstrap_timeout
}

// Run in the worktree root after copy, before returning the worktree.
// Returns a *BootstrapError (command, exit info, captured tail) on failure.
func runBootstrap(ctx context.Context, worktreeDir string, spec BootstrapSpec, out io.Writer) error

type BootstrapError struct {
    Command string
    Err     error // exit status / timeout / start failure
    Tail    string
}
// Message shape: `worktree bootstrap failed: <command>: <reason>`.
```

### Where bootstrap runs

Worktree creation order becomes: create the worktree from committed state → copy
`worktree.copy` files → **run `worktree.bootstrap`** → return the ready worktree.
This holds for every worktree that hosts Agent work: the Run Worktree (sequential
and review Runs) and each Task Worktree (concurrent Runs). `fetch` creates no
worktree and never bootstraps. The command is shell-invoked in the worktree root
with the Run's environment, so it can use the copied `.env`.

### Failure handling

- **Run Worktree** bootstrap failure: the Run ends before Agent work with a
  distinct outcome — Preflight-adjacent for the CLI, exit code consistent with a
  setup failure — and stderr carries `worktree bootstrap failed: <command>:
  <reason>`. No Batch is assigned, nothing is committed.
- **Task Worktree** bootstrap failure (concurrent Tasks): that Task settles
  `failed` with reason `worktree bootstrap failed: <command>: <reason>`; other
  independent Tasks are unaffected (failure isolation, per 0009). Bootstrap runs
  before the Task's Agent work, so no Task work is lost.
- Bootstrap output is captured to the Run Event Journal and streamed to stderr
  like other Run diagnostics.

### Bounded execution

`runBootstrap` runs under a context with `worktree.bootstrap_timeout` (default
`10m`); a timeout is a bootstrap failure with a timeout reason. This prevents a
hung install or migration from stalling a Run.

### Env-file recipe (docs, no new code)

`worktree.copy` already copies untracked repository-relative files into each
worktree. The docs tie the pieces together for stateful monorepos:

```yaml
worktree:
  concurrency: 1                    # shared DB: bootstrap once on the reused Run Worktree
  copy: [".env", "packages/backend/.env"]
  bootstrap: "bun install && bun run db:migrate && bun run db:seed"
  bootstrap_timeout: "10m"
```

with the note that copied files must be gitignored so Batch commits never leak
them.

## Coverage Map

- Stories 1-2 → `worktree.bootstrap` config + `runBootstrap` in worktree creation
- Story 3 → `BootstrapError` message + Run/Task failure mapping
- Story 4 → once-per-Run-Worktree at concurrency 1 (reused Run Worktree)
- Story 5 → env-file recipe docs (existing `worktree.copy`)

## Integration Points

Local process execution only: the bootstrap command runs via the OS shell in the
worktree directory. No new external systems. Reuses the existing worktree
creation seam that already applies `worktree.copy`.

## Testing Approach

- Config: table tests for `worktree.bootstrap`/`bootstrap_timeout` precedence and
  strict decoding.
- `runBootstrap`: unit tests with a fake command — success runs in the worktree
  dir and returns nil; a non-zero exit returns a `BootstrapError` with the
  command and tail; a sleeping command past the timeout returns a timeout
  BootstrapError. No real `bun`/DB in unit tests.
- Run wiring: a Run Worktree bootstrap failure ends the Run with the message and
  assigns no Batch; a Task Worktree bootstrap failure settles only that Task
  failed and leaves independent Tasks unaffected; an empty `worktree.bootstrap`
  leaves Run behavior byte-stable.

## Build Order

1. Config `worktree.bootstrap` + `worktree.bootstrap_timeout` and `runBootstrap`
   in `internal/worktree`, wired into Run Worktree creation after copy, with
   Run-level failure mapping (no deps)
2. Task Worktree wiring for concurrent Runs: bootstrap each Task Worktree,
   bootstrap failure settles that Task failed, reusing `runBootstrap` (depends
   on: 1)
3. Docs and skill sync (env-file recipe + bootstrap) (depends on: 1, 2)

## Risks & Considerations

- Bootstrap runs before Agent work, so a slow install/migrate adds latency to
  every worktree; the timeout bounds it and the concurrency-1 recipe keeps it to
  once per Run for shared-state projects.
- A bootstrap command that mutates shared state (one database) under
  concurrency > 1 will race — documented; the mitigation (concurrency 1 or a
  per-worktree database in the command) is the project's responsibility, not
  Roundfix's.
- Copied `.env` files must be gitignored or a Batch commit could leak them; the
  docs state this explicitly (Roundfix already excludes its own config from Batch
  commits but does not police arbitrary copied files).
- Bootstrap output can be large; it is journaled and streamed but not embedded in
  the deterministic stdout report.

## Decisions

- `worktree.bootstrap` runs once per new worktree after `worktree.copy`, before
  Agent work and Verification; `worktree.bootstrap_timeout` bounds it (default
  `10m`). See ADR-0034.
- Bootstrap failure is a distinct Run-failed / Task-failed signal with a
  `worktree bootstrap failed: <command>` message.
- Env files ride the existing `worktree.copy`; the glossary gains Worktree
  Bootstrap.
