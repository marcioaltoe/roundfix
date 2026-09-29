---
spec: 0183-run-storage-that-says-what-it-holds
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# Run storage that says what it holds

## Executive Summary

Make "reclaimable" mean that something is left: a terminal Run past the
Journal Retention cutoff counts only while it holds Run Event Journal rows or
an artifact directory. The store prune counts events on the read connection
first, and takes the write lock only when some candidate still has events. One
predicate feeds `gc`, `gc --dry-run`, the operational sweep and a new
read-only Doctor `storage` check. `gc sanitize` also accepts the default Artifact
Root derived from a Run's recorded checkout. The retained-Run count behind
`runs list` falls back to the Run's repository key when the recorded checkout
is gone, instead of warning. The trade-off accepted: Doctor reads only cheap
facts, the header's free-page count and one `lstat` per candidate. It reports
counts, not the bytes `gc --dry-run` would measure, so the notice is a pointer
to `gc`, not a measurement. No Run Database schema changes.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0171 and ADR-0172 (this Spec) are
  implemented here; ADR-0023, ADR-0032, ADR-0033, ADR-0053, ADR-0080,
  ADR-0090, ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0107,
  ADR-0117, ADR-0155, ADR-0156 and ADR-0166 to ADR-0170 hold as the PRD
  states; ADR-0097, ADR-0141, ADR-0161 and ADR-0173 do not apply, for the reasons the
  PRD records. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

Four existing seams change; one file is new.

- `internal/store/journal.go` — `PruneTerminalRuns` and `PruneResult`;
  `ArtifactRootRun` gains the recorded checkout. The new
  `internal/store/free_bytes.go` holds the `FreeBytes` read.
- `internal/cli/gc.go` — the reclaimable predicate, `runGC`,
  `pruneRunRetention` and `classifyGCSanitationRoot`.
- `internal/cli/doctor.go` and the new `internal/cli/doctor_storage.go` — the
  `storage` check; `internal/cli/health.go` names it.
- `internal/config/config.go` — the default Artifact Root of a path without
  re-resolving it.
- `internal/worktree/worktree.go` — `countRetainedTerminalRuns`.

## Implementation Design

### Interfaces

```go
// internal/store
type PruneResult struct {
	RunIDs         []string // Runs whose Run Events were deleted, in completion order
	Events         int      // rows the DELETE removed
	EligibleRunIDs []string // every Run past the cutoff at prune time
}
func (store *Store) FreeBytes(ctx context.Context) (int64, error)

// internal/cli
func retentionReclaimable(
	candidates []store.PruneCandidate,
	hasArtifactDir func(runID string) (bool, error),
) ([]string, error)

// internal/config
func DefaultArtifactDirectoryForPath(path string, homeDir string) (string, error)
```

### Data Models

No schema change. `store.ArtifactRootRun` gains `GitRoot string`, selected from
the existing `git_root` column beside the coalesced repository key.

### A prune reports only what it reclaimed

`PruneTerminalRuns` keeps `terminalRunPruneCandidates` as its eligibility scan.
It then fills event counts with `countPruneCandidateEvents` on the read
connection. When no candidate has events, it returns
`PruneResult{EligibleRunIDs: <all candidates>}` without opening a write
transaction. Otherwise, inside one `withWriteTx`, it counts the candidates'
events again and deletes the events of the candidates that have them at that
moment. The count and the deletion therefore come from the same transaction;
the read-connection count only decides whether a write transaction is needed.
`RunIDs` names those candidates, and `Events` is the number of rows the DELETE
removed. `TerminalRunPruneCandidates` is
unchanged.

`retentionReclaimable` returns, in candidate order, every candidate whose
`Events` is above zero or whose artifact directory exists. It relies on
`TerminalRunPruneCandidates`, which already fills each candidate's `Events`
through `countPruneCandidateEvents` on the read connection. A Run that holds
journal rows but no artifact directory is therefore reclaimable in `gc` and in
Doctor alike. `runGC` uses it for
the dry run, with `hasArtifactDir` answering from the directories
`gcCandidateArtifactDirs` found, so `Runs eligible` and
`Eligible Runs` list only reclaimable Runs. `pruneRunRetention` keeps its
order: it measures the candidates' directories, calls `PruneTerminalRuns`,
filters the directories by `EligibleRunIDs` instead of `RunIDs`, removes them
and reports the union of `RunIDs` and the removed directories' Run IDs, in
candidate order. A Run whose events were pruned earlier but whose directory
survived is therefore still removed and reported. `sweepRunRetention` keeps its
silence rule, which now holds whenever nothing was reclaimed.

### Doctor reports reclaimable storage

`health.go` adds `HealthCheckStorage = "storage"`. `doctorDependencies` gains
`storage func(context.Context, roundconfig.Loaded) []CheckResult`, defaulting to
`defaultDoctorStorageResults` in the new `doctor_storage.go`.
`runDoctorCommand` appends its results after the residue results and before
`codex`. The check:

1. Returns `ok (no Run Database)` when `store.DatabasePath` does not exist,
   creating nothing.
2. Opens `store.OpenReader`. A failure, including a schema-version mismatch,
   returns `partial (could not read Run Database: <err>)`.
3. Reads `FreeBytes`, which runs `PRAGMA page_size` and
   `PRAGMA freelist_count` and multiplies them. It runs no transaction, no
   `dbstat` and no pragma write.
4. With a non-zero Journal Retention, reads `TerminalRunPruneCandidates` at
   `now - retention` and applies `retentionReclaimable`. `hasArtifactDir`
   checks `lstat` of `gcRunArtifactPath(<root>, <id>)` under the root that
   `ResolveArtifactDirectory` gives the loaded config, and never walks it.
   Outside a Git repository it answers `false` without resolving any root,
   even when `artifact_dir` is absolute, and the detail says
   `artifact directories not inspected outside a Git repository`.
5. Reports `found` when the reclaimable count is above zero or free bytes reach
   `doctorStorageFreeBytesThreshold` (64 MiB). The detail is
   `Runs reclaimable: <n>; Run Database free bytes: <b>; next: <commands>`,
   where the next action is `roundfix gc`, `roundfix gc compact` or both, for
   whichever condition held. Otherwise it reports
   `ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: <b>)`.
   With a zero retention, the Runs clause reads
   `Journal Retention 0 keeps every Run`.

It never returns `failed`, so Doctor's exit code is unchanged. It closes the
reader and reports a close failure as `partial`.

### `gc sanitize` recognizes a pre-key default root

`DiscoverArtifactRoots` selects `git_root` into `ArtifactRootRun.GitRoot`.
`DefaultArtifactDirectoryForPath` returns
`<home>/.roundfix/artifacts/<repoID(clean path)>`. It refuses an empty path or
home, and it never calls `RepositoryRoot`. In `classifyGCSanitationRoot`, the
per-repository default loop becomes a per-Run check that runs after every
existing guard. A Run accepts the root when it equals the default
`ResolveArtifactDirectory("", run.Repository, home)` gives, or the default
`DefaultArtifactDirectoryForPath(run.GitRoot, home)` gives. A Run matching
neither makes the root `overridden`, and the evidence names both defaults. An
error deriving the key default keeps the root `unsafe` unless the checkout
default matches. When the checkout default was the match, the classification
evidence names that checkout.

### `runs list` ignores a vanished checkout

In `countRetainedTerminalRuns`, a group whose `recordedGitRoot` fails with an
error wrapping `fs.ErrNotExist` is not a failure. For each distinct non-empty
`RepositoryRoot` among the group's Runs that differs from the group root, the
count looks up that key, listing each key once:

- an absent key is skipped;
- a key that passes `recordedGitRoot`, or is a bare repository whose
  `git rev-parse --absolute-git-dir` is the key, has its Run Branches listed;
- any other validation failure is appended to the failures, as today.

Each Run's recorded Run Worktree path is still checked. Every other
`recordedGitRoot` failure is still a warning.

## API Contracts

1. `roundfix gc` and `roundfix gc --dry-run` count, list and report only Runs
   that still hold Run Event Journal rows or an artifact directory. A second
   `gc` after a complete one reports `Runs pruned: 0`. The operational sweep
   prints its stderr line only when it reclaimed something.
2. `roundfix doctor` prints one `storage:` line with status `ok`, `found` or
   `partial`, never `failed`, and writes nothing.
3. `roundfix gc sanitize` classifies a root equal to a Run's checkout-derived
   default like a key-derived default root.
4. `roundfix runs list` prints no warning for a terminal Run whose recorded Git
   root does not exist.

## Coverage Map

- Goals 1–2 → A prune reports only what it reclaimed; API Contract 1.
- Goal 3 → Doctor reports reclaimable storage; API Contract 2.
- Goal 4 → `gc sanitize` recognizes a pre-key default root; API Contract 3.
- Goal 5 → `runs list` ignores a vanished checkout; API Contract 4.
- Core Features 1–4 → the four sections above, in order.
- Success Metrics 1–2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- **ADR-0033's retention.** Eligibility, the cutoff and the protected rows are
  unchanged; only the count and the lock use change.
- **Spec 0162's durable repository key.** Sanitation and the retained-Run count
  read the key it records and the checkout it kept beside it.
- **Spec 0175's legacy key backfill.** It made Runs of removed checkouts visible
  through their key; the count now inspects them there.

## Testing Approach

Every test uses a temporary Roundfix Home and disposable repositories. None
opens `~/.roundfix`.

1. **Prune.** Store tests against a real SQLite database:
   - a candidate with no events is neither deleted nor reported, and a store
     holding only such candidates opens no write transaction;
   - `EligibleRunIDs` names every candidate.

   Public `gc` tests:
   - a second `gc` reports `Runs pruned: 0`, and a dry run then reports
     `Runs eligible: 0`;
   - a Run with no events but a surviving directory is removed and reported;
   - the operational sweep prints no line over emptied Runs.

   `TestPruneTerminalRunsDeletesOnlyEligibleJournalRows` and the existing `gc`
   and sweep tests stay green.
2. **Doctor.** Through `runDoctorCommand`, with the other checks faked and a
   real Run Database in a temporary home:
   - `found` names `roundfix gc` for a reclaimable Run, and `ok` shows when
     nothing is left;
   - free bytes at the threshold produce `found` naming `roundfix gc compact`;
   - a missing database prints `ok (no Run Database)` and creates nothing;
   - a schema mismatch prints `partial`;
   - the exit code stays `0`, and the database file bytes are unchanged.
3. **Sanitize.** A bare clone with `git worktree add`, and a Run recorded with
   a root derived from its checkout path:
   - the root is `orphaned`, and `--apply` removes its eligible directory;
   - a root equal to neither default stays `overridden`.

   `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval` stays green.
4. **Retained count.** Unit tests with the package's fake Git runner:
   - an absent recorded root yields no failure, and a branch is counted through
     the key;
   - an absent key counts only an existing Run Worktree path;
   - a symlinked recorded root still fails.

   A public `runs list` test with a deleted linked checkout prints no warning.
   `TestCountRetainedTerminalRunsBatchesGitInspectionByRepository` stays green.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. A prune reports only what it reclaimed (depends on: none).
2. `runs list` ignores a vanished checkout (depends on: none).
3. Doctor reports reclaimable storage (depends on: 1 — it reuses
   `retentionReclaimable`, and both edit `docs/user-guide/commands.md`, the
   Roundfix skill and `CONTEXT.md`).
4. `gc sanitize` recognizes a pre-key default root (depends on: 1 — both edit
   `internal/store/journal.go` and `internal/cli/gc.go`).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Doctor must never write.** It opens only the read-only reader and checks
  existence first, so a missing database is never created. The test compares
  the database bytes before and after.
- **Late events on a terminal Run.** The deletion and its count happen in one
  write transaction, so nothing is deleted without being counted. An event
  appended after that transaction commits is left for the next prune.
- **A checkout default that is not Roundfix's.** Only a root equal to the
  default derived from the Run's own recorded checkout is accepted, and every
  earlier guard still has to pass.
- **Bare keys.** Branch listing works in a bare repository, and validation
  accepts one only when its absolute Git directory is the key itself.

## Decisions

- Count what is left, and lock only to delete. See ADR-0171.
- A read-only Doctor `storage` check from cheap facts, never `failed`. See
  ADR-0172.
- One reclaimable predicate for `gc`, the dry run, the sweep and Doctor.
- Accept the recorded checkout's default without re-resolving it.
- Fall back to the repository key before concluding nothing is retained.
