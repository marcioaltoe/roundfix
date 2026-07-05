---
spec: 0007-command-lifecycle
prd: _prd.md
created: 2026-07-05
---

# Command Lifecycle — Technical Spec

## Executive Summary

Four lifecycle features over existing seams, with one consequential change:
the Run Database becomes a control channel (schema v5, Stop Requests) —
accepted because ADR-0004 already made it the single source of Run state and
the future daemon (work-plan item 3) needs exactly this primitive; the
alternative (process signals) dies with the terminal that owns the Run.
Setup and upgrade are new commands over existing probes and the release
channel; push-at-Clean reuses the review path's Pusher behind a Project
Config key. Nothing touches prompts, verification, or commit contracts.

## System Architecture

- `internal/cli` — new `setup` and `upgrade` command files (stop/implement
  conventions); `stop` gains the graceful default and `--force`; implement
  wires the push step at Clean.
- `internal/store` — schema v5: `stop_requested_at` on runs;
  `RequestStop(runID)` + the engine-side check; migration preserving v4
  rows and locks.
- `internal/daemon` — both engines poll the Stop Request flag at settlement
  boundaries (between Batches/Tasks, before the QA step) and end Stopped
  after settling in-flight work; TaskCycle's caller performs the optional
  Clean push via the existing `Pusher`.
- `internal/agent` — small additions: acpx presence/install helpers reused
  by setup; the session-cancel invocation reused by `stop --force`
  (`acpx <agent> cancel -s roundfix-<run-id>` — the session name is
  derivable from the run id by construction).
- `internal/config` — `implement.auto_push` (bool, default false, Project
  Config); validation mirrors `watch.auto_push` (requires auto-commit
  semantics, which spec Runs always have).
- `internal/app` — release repository constant + version comparison for the
  freshness check; check cache file in Roundfix Home.

## Implementation Design

### Interfaces

```go
// internal/store — schema v5
func (s *Store) RequestStop(ctx context.Context, runID string) error
func (s *Store) StopRequested(ctx context.Context, runID string) (bool, error)

// internal/daemon — checked at settlement boundaries by both engines
// (no interface change; engines already hold the RunStateStore seam,
//  which gains the two methods above via the store)

// internal/app — release channel
const ReleaseRepository = "<owner>/<repo>" // derived from the module path
func LatestRelease(ctx context.Context) (tag string, assets []ReleaseAsset, err error) // gh api via exec, no new deps
```

### Setup Command

`roundfix setup [--yes] [--no-input]` — ordered idempotent checks, one
deterministic report line each (`ok | installed | skipped | offered:
declined | failed`):

1. Node present with the documented minimum.
2. acpx present at the pin; missing/mismatched → offer
   `npm install -g acpx@<pin>` (run on confirm/`--yes`).
3. Runtime adapters: probe configured agent; when local adapter binaries
   (`codex-acp`, `claude-agent-acp`, `opencode`) are on PATH and the acpx
   config lacks the agents override, offer writing it to
   `~/.acpx/config.json` (the dogfood finding-27 remediation, applied
   surgically: only the `agents` entries, preserving the rest).
4. User Config / Project Config missing → offer creation through the
   existing Init Command flows.
Exit 0 when everything ends ok/installed/declined-by-choice; 1 when any
action failed; 2 for usage errors. Second run with a healthy environment
prints all-ok and changes nothing.

### Upgrade Command and freshness check

- `roundfix upgrade [--check]`: resolve `LatestRelease` (gh api on the
  release repository — same gh dependency preflight already uses); compare
  against the built version; `--check` only reports. Download the
  platform asset to a temp path, verify (size + checksum asset when
  published), rename over the current executable atomically (write-sibling
  + rename; on Windows-style rename failure, report the manual path).
  Outcomes: `upgraded <old> → <new>`, `already current <version>`,
  `no releases published`, each deterministic; failures exit 1 with the
  bounded stderr tail convention.
- Freshness check: on operational commands (fetch/resolve/watch/implement),
  read the cache in Roundfix Home (`version-check.json`: last check time +
  latest seen); if older than 24h, refresh best-effort with a short
  timeout; when behind, one stderr line naming both versions and
  `roundfix upgrade`. Silent on any failure; never delays the Run more than
  the short timeout; disabled in tests via the existing seam pattern
  (package-var clock/fetcher).

### Stop semantics (ADR-0022)

- Default: resolve the target Run (existing selectors); if Active, call
  `RequestStop` and report `Stop Request recorded; the Run stops after the
  current Work Item settles.` Engines check `StopRequested` after each
  settlement (and before starting the QA step); on true: journal the stop,
  end Stopped through the existing paths (locks released by completion).
- `--force`: best-effort `acpx <agent> cancel -s roundfix-<run-id>` for the
  Run's session (agent from the Run's stored runtime; failure logged, not
  fatal), then today's force-completion and lock release. Works for dead
  processes by construction.
- In-terminal Ctrl-C behavior is unchanged.

### Push at Clean (ADR-0021)

`implement.auto_push: true` in Project Config: after a spec Run completes
Clean (including the QA-pass case), the CLI runs the existing `Pusher`
against the branch's detected upstream (preflight already inspects it);
missing upstream is a stderr note, not a failure. Never wired for any other
outcome. stdout gains one final `pushed <remote>/<branch>` line only when a
push happened.

### Data Models

Schema v5 migration: `runs` gains nullable `stop_requested_at`; existing
rows migrate untouched; `user_version = 5`. Config gains the one key. The
freshness cache is a plain JSON file in Roundfix Home, not the database.

### API Contracts

New: `setup` (flags `--yes`, `--no-input`), `upgrade` (flag `--check`),
`stop --force`, config `implement.auto_push`, the daily stderr freshness
line, and the `pushed` stdout line on opted-in Clean spec Runs. Existing
stop selectors, exit codes, and every other contract unchanged; graceful
stop reuses exit 0.

## Coverage Map

- Story 1 → Setup Command (round-1 findings 22, 27 remediation)
- Story 2 → Upgrade Command (finding 23)
- Story 3 → freshness check (finding 23)
- Story 4 → Stop Requests + engine boundary checks (finding 24, ADR-0022)
- Story 5 → `--force` with session cancel (finding 24, ADR-0022)
- Story 6 → `implement.auto_push` (finding 25, ADR-0021)

## Integration Points

- **npm** (setup's acpx install — the documented command, run only on
  confirmation) and **gh** (releases API — existing dependency).
- **acpx** — the cancel invocation via existing builders; the setup config
  write edits `~/.acpx/config.json` surgically.

## Testing Approach

Existing patterns: buffer-captured CLI tests with package-var seams for the
npm/gh/exec boundaries (hand-rolled fakes recording invocations); store
migration v4→v5 fixture test plus Stop Request round-trip; engine tests
proving the in-flight Task settles (commit exists) before Stopped, and the
QA step is skipped when a Stop Request precedes it; force-stop test
asserting the cancel invocation and immediate completion; push-at-Clean
matrix (Clean pushes, all other outcomes do not, missing upstream notes);
upgrade against fixture release JSON + a temp fake binary path; freshness
cache aging tests with a fake clock. Full suite as the net; `-race` on the
daemon package (engines gain a polled read).

## Build Order

1. Stop Requests: schema v5, store methods, engine boundary honor, stop
   default (no deps)
2. Setup Command (no deps)
3. Upgrade Command and freshness check (no deps)
4. Force stop with Agent Session cancel (depends on: 1)
5. Push at Clean via Project Config (no deps)
6. Docs and skill sync (depends on: 1, 2, 3, 4, 5)

## Risks & Considerations

- Self-replacing binaries are platform-sensitive: atomic rename on the same
  filesystem, clear manual fallback message elsewhere; never leave a
  half-written executable (write sibling, rename last).
- Setup edits a user-owned file (`~/.acpx/config.json`): surgical merge,
  confirmation required, and a printed diff of what changed.
- The engine's Stop Request poll is one store read per settlement boundary
  — negligible; `-race` guards the new read path.
- The freshness check must be unobservable in failure: short timeout,
  swallow everything, cache regardless.
- Graceful stop on a dead process never fires; the stop report must say
  when to use `--force` (the message names it).

## Decisions

- Stop Requests through the Run Database; graceful default; `--force`
  covers dead Runs and adds the cooperative cancel. See ADR-0022.
- Spec-Run push only at Clean, only via Project Config, default off; PR
  creation permanently out. See ADR-0021.
- Setup wraps Init flows and applies the acpx agents override only with
  confirmation; upgrade ships before the first release exists by design.
- Freshness state lives in a Roundfix Home JSON file, not the database.
