---
spec: 0028-settlement-and-reporting
prd: _prd.md
created: 2026-07-14
---

# Settlement and Reporting Robustness — Technical Spec

## Executive Summary

Six small, independent hardening changes on existing seams — no new packages, no new architectural elements. The primary trade-off this design accepts is conservatism over convenience in orphan detection: the only accepted proof of owner death is a liveness signal that positively fails (no such process), so a reused PID or a permission-denied probe keeps the block and the manual force-stop path. That trades a rare residual manual step for zero risk of two live Runs racing one target (ADR-0044). Everything else is additive: a schema-versioned column for the owner PID, a normalization step ahead of an existing validation gate, warning events at an existing publisher, extra deterministic report lines, and a preflight probe extension. This spec assumes 0027 has landed (it consumes the `terminal_reason` artifact field and leaves the review report's reason lines to 0027).

## System Architecture

Existing modules extended, in glossary vocabulary:

- **`internal/store`** — Run rows gain the owning process id via the established `PRAGMA user_version` migration pattern; a reclamation operation completes a provably-orphaned Run as Failed and releases its Active-Run lock.
- **`internal/cli`** — every Active-Run lock consumption point (Implement Command preflight, Settle Command preflight, review preflights including 0027's Branch Integrity Preflight, Stop Command resolution, and the create-race path) gains the orphan check before surfacing the block; the implement report renderer gains reason lines; the Settle Command prints its commit path set and a shared-worktree warning.
- **`internal/spec`** — task frontmatter parsing gains status-synonym normalization ahead of the existing allowed-status gate, reusing the byte-preserving status rewrite helper.
- **`internal/daemon`** — the task engine warns (Run Event + stderr) on no-op Task commits and returns per-Task outcomes (status + reason) to the CLI in its cycle result instead of journaling them only.
- **`internal/agent`** — the acpx probe gains an adapter-binary resolution step with per-adapter install hints; the Doctor Command reports the same check.
- **`.agents/skills/write-tasks`** (repo-owned per skill governance) and **`skills/roundfix`** — authoring guidance and the shipped skill's settle/report output contracts updated with the behavior they describe.

## Implementation Design

### Interfaces

Owner identity and liveness (`internal/store` + a small platform helper):

```go
// CreateRunRequest gains the owner's process id, recorded at Run creation.
type CreateRunRequest struct {
    // ...existing fields
    OwnerPID int
}

// ProcessAlive reports whether pid provably exists. Only a definitive
// "no such process" result counts as proof of death; permission errors
// and non-unix platforms report alive (no proof).
func ProcessAlive(pid int) bool

// ReclaimOrphanedRun completes an Active Run whose owner is provably dead:
// state Failed, reason recorded, lock released, journal entry written.
func (s *Store) ReclaimOrphanedRun(ctx context.Context, run Run, reason string) error
```

Status normalization (`internal/spec`):

```go
// NormalizeStatus maps documented synonyms to canonical statuses:
// "done" → completed; hyphen/space variants of canonical values
// ("in-progress", "in progress") → their canonical form.
// Unknown values pass through unchanged and fail AllowedStatus as today.
func NormalizeStatus(raw string) string
```

Per-Task outcomes returned to the CLI (`internal/daemon`):

```go
type TaskOutcome struct {
    Task   string // task id
    Status string // completed | failed | skipped
    Reason string // one line; empty for completed
}

type TaskCycleResult struct {
    // ...existing QA fields
    Outcomes []TaskOutcome
}
```

Adapter probe (`internal/agent`):

```go
// resolveAdapterCommand returns the adapter binary acpx will spawn for a
// runtime: the agents map in acpx's config when present, else the built-in
// default (codex-acp, claude-code-acp, opencode). Stdio overrides return
// the override command.
func resolveAdapterCommand(runtime RuntimeSelection) (string, error)
```

### Data Models

- **`runs` table**: new nullable `owner_pid` integer column via the next schema version's migration (same pattern as the existing `work_dir`/`model` column adds). Pre-migration rows have no PID and therefore never qualify as provably orphaned.
- **Active-Run locks**: unchanged — the owner is resolved through the lock's `run_id` → Run row.
- **Reclaimed Runs**: terminal state `Failed`, with a reason shaped like `owner process <pid> not running; lock reclaimed`. Reclamation is journaled as a Run Event on the reclaimed Run.
- **Task status synonyms**: a package-level documented map in `internal/spec`; the canonical vocabulary is unchanged.

### API Contracts

Public CLI surface changes (additive; breaking-change review applies):

- **Implement report**: failed and skipped Task lines are followed by one indented line `  reason: <one line>`. The existing `task_NN <status> — <title>` line shape is unchanged, so existing line parsers keep working; completed Tasks gain no extra line. The reason for a Verification failure names the failed command and exit status with a pointer to the diagnostics.
- **Settle report**: after the verification lines, settle prints one `commit <path>` line per committed path in sorted order before the existing `settled <task> completed — <sha>` line. When other Tasks in the Spec are `failed` at settle time, one stderr warning names them and states their work may be included in this commit.
- **No-op Task commit**: settlement proceeds (PRD decision: warn, don't block) with one stderr warning and one Run Event when the Task commit contains no path outside the Spec Root — including the currently-silent empty-stageable case.
- **Orphaned lock**: any command blocked by an Active-Run lock whose owner is provably dead reclaims it automatically: one stderr warning naming the dead run id and PID, a journal entry, then the command proceeds. A live owner keeps today's `ActiveRunError` text unchanged.
- **Adapter diagnostics**: when the runtime's adapter binary cannot be found, preflight fails with `<adapter> is required but was not found on PATH; install it with: <command>` (mirroring the existing acpx-missing message) instead of a raw acpx stderr tail. `roundfix doctor` gains one `adapter: ok|failed` line for the configured agent.

## Coverage Map

- Goal 1, Story 1 → owner PID column + `ProcessAlive` + `ReclaimOrphanedRun` (`internal/store`), orphan hook at every lock consumption point (`internal/cli`); ADR-0044.
- Goal 2, Story 2 → `NormalizeStatus` in task parse/reload + daemon canonical rewrite via the existing status-rewrite helper (`internal/spec`, `internal/daemon`).
- Goal 4 (no-op visibility), Story 3 → no-op commit warning at the task engine's commit step (`internal/daemon`).
- Goal 4 (settle visibility), Story 4 → settle commit path listing + sibling-failed warning (`internal/cli`).
- Goal 3, Story 5 → `TaskOutcome` in the daemon cycle result + implement report reason lines (`internal/daemon`, `internal/cli`).
- Goal 5, Story 6 → adapter-binary probe with install hints + Doctor line (`internal/agent`).
- Core Feature 7 → write-tasks authoring guidance (portable Verification forms, effect-proving requirement) and Roundfix Skill output-contract sync.

## Integration Points

- **Operating system process table** via the liveness helper: unix signal-0 semantics behind a small platform-guarded function; non-unix builds report "no proof" and never reclaim.
- **acpx configuration** (`~/.acpx/config.json` agents map): read-only resolution of the adapter command for the probe; absent or unparsable config falls back to built-in adapter defaults. Roundfix already owns acpx setup, so this stays inside the existing acpx boundary in `internal/agent`.
- No new external systems; GitHub is untouched by this spec.

## Testing Approach

All existing seams — stdlib table tests, fakes, buffer-captured CLI runs.

- **Unit**: `ProcessAlive` against the current process (alive) and a known-dead PID (spawn a child, wait it out); `NormalizeStatus` table across the synonym set, canonical values, and garbage; reason-line rendering fixtures for the implement report; settle path-list ordering; adapter-command resolution against fixture acpx config files (present, missing, malformed).
- **Integration (buffer-captured CLI runs)**: a store seeded with an Active Run owned by a dead PID lets implement/settle/stop proceed with the reclamation warning, while a lock owned by the live test process still blocks with today's error; settle in a temp repo prints `commit <path>` lines and the sibling-failed stderr warning; implement report shows reason lines for failed and skipped tasks fed by a stubbed engine result.
- **Engine tests**: no-op commit publishes the warning event in both the spec-root-only and empty-stageable cases; `TaskCycleResult.Outcomes` carries the settle reasons already published to the journal; synonym-written task files (`done`) reload as completed with the file rewritten to canonical form.
- **Migration test**: opening a pre-existing schema database applies the new version and leaves old rows PID-less (never reclaimable), following the existing migration test pattern.

## Build Order

1. **Owner PID groundwork**: schema migration adding `owner_pid`, `CreateRunRequest.OwnerPID` recorded at every Run creation call site, `ProcessAlive` platform helper with tests.
2. **Orphan reclamation** (depends on: 1): `ReclaimOrphanedRun` (Failed state, reason, journal, lock release) and the orphan check wired into every lock consumption point — implement, settle, review preflights (including the Branch Integrity Preflight's active-run guardrail from 0027), stop resolution, and the create-race path — each with the stderr warning.
3. **Status synonym normalization**: `NormalizeStatus` ahead of the allowed-status gate in parse/reload, daemon rewrite-to-canonical on reload via the existing status-rewrite helper, synonym set documented.
4. **No-op Task commit warning**: spec-root/workdir path classification at the commit step, warning event + stderr in both no-op shapes.
5. **Per-Task outcomes in the report** (depends on: 4 only for shared task-engine test fixtures; logically independent): engine accumulates `TaskOutcome` per settle/skip with the same reason strings it journals; implement report renders the indented reason lines.
6. **Settle transparency**: sorted `commit <path>` stdout lines, sibling-failed stderr warning from the already-loaded Task Graph.
7. **Adapter-binary diagnostics**: `resolveAdapterCommand`, LookPath probe with per-adapter install hints, preflight error text, Doctor line.
8. **Docs and skill sync** (depends on: 2, 4, 5, 6, 7): write-tasks authoring guidance (portable shell forms in Verification, effect-proving requirement) in the repo-owned `.agents/skills` source with the generated mirror refreshed, and the Roundfix Skill's settle/implement output contracts and stop/lock guidance updated to the shipped behavior, keeping the skills check green.

## Risks & Considerations

- **PID reuse**: a recycled PID makes a dead owner look alive; the design accepts the resulting manual `stop --force` rather than adding process-start-time matching now (ADR-0044's "anything short of proof keeps the block" — a live-looking PID is not proof of death).
- **Reclamation races**: two commands may probe the same orphaned lock concurrently; reclamation must be idempotent at the store layer (second caller finds the Run already terminal and proceeds).
- **Report consumers**: the implement report gains lines but changes none; the settle report inserts lines before the settled line — the Roundfix Skill documents the settle byte shape, so the skill update in step 8 is mandatory, not optional (skill-sync hard rule).
- **acpx config drift**: the agents map is user-owned; resolution must degrade to built-in defaults silently rather than failing preflight on a malformed config, since acpx itself may still work.
- **Normalization scope creep**: the synonym map is deliberately tiny and documented; anything fuzzy (edit distance, casing beyond exact synonyms) stays rejected so the canonical vocabulary keeps its meaning.

## Decisions

- Proof of owner death is a definitive no-such-process liveness result; permission errors, unknown platforms, and PID-less legacy rows never reclaim. Refines ADR-0044.
- Reclaimed Runs settle as Failed (the Run itself broke), not Stopped (no user asked) — consistent with the glossary's outcome vocabulary.
- Per-Task reasons travel in the daemon's in-memory cycle result, not via a CLI journal query at report time — one producer, no read-back coupling.
- Implement-report reasons are an additive indented line, keeping the existing per-Task line byte-stable for parsers.
- The review report's reason lines belong to 0027 (its report-split task) and are out of scope here; this spec covers the implement and settle surfaces.
- Adapter resolution reads acpx's own config for truth and falls back to built-in defaults; Roundfix does not maintain a parallel adapter registry in Project Config.
- Synonym set: `done` → `completed`, plus hyphen/space variants of canonical statuses — exactly the PRD's open-question default, now fixed.
