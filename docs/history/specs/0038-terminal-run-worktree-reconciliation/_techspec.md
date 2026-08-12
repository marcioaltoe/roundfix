---
spec: 0038-terminal-run-worktree-reconciliation
prd: _prd.md
created: 2026-07-17
---

# Terminal Run Worktree reconciliation — Technical Spec

## Executive Summary

The existing worktree module gains one classifier shared by the new Reconcile Command, the automatic terminal reaper, and Run listing projection. Classification replaces the unsafe creation-base shortcut with two independent proofs: the retained worktree is clean and the Run Branch tip is an ancestor of the recorded target branch tip. The CLI is read-only unless `--apply` is supplied and emits stable text or JSON. The main trade-off is conservative residue: missing refs, missing target metadata, or Git errors remain `unknown` even when a human might infer safety, because preserving redundant storage is cheaper than deleting unique work.

## Project Constraints

- Identifier strategy: not applicable — the classifier consumes existing Run,
  branch, path, and Git object identities without minting a new project-owned
  identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all inspection and apply behavior
  stays within local Git and the Run Database. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0023, ADR-0024, ADR-0052, and
  ADR-0053 bind Run Worktree ownership, porcelain integration, the sole guarded
  terminal transition, and positive cleanup proof. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`. On 2026-07-27, the maintainer additionally
  expressly authorizes the deterministic Skill-digest fallout of that edit in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- **`internal/worktree`** owns `RunWorktreeReconciliation`: discovery, cleanliness inspection, ancestry proof, classification, and safe cleanup. `PruneTerminalReport` delegates to this classifier instead of comparing the branch tip only with its reflog creation base.
- **`internal/cli/reconcile.go`** owns argument parsing, repository scoping, text/JSON projection, dry-run/apply behavior, and exit semantics.
- **`internal/store`** supplies terminal spec Runs and receives reconciliation events. A safe Integration Pending entry calls spec 0037's guarded `ReconcileIntegration` before Git cleanup.
- **`internal/cli/runs.go`** preserves the existing stdout row shape and adds exact retained-worktree guidance on stderr.
- **`internal/cli` Preflight reaping** uses the same classifier, so automatic and explicit cleanup cannot disagree about safety.

The design adds no new database table. Current state is derived from Run metadata plus Git; successful reconciliation is durable through the existing Run state and Run Event Journal.

## Implementation Design

### Interfaces

```go
type ReconciliationState string

const (
    ReconciliationSafe         ReconciliationState = "safe"
    ReconciliationUnintegrated ReconciliationState = "unintegrated"
    ReconciliationDirty        ReconciliationState = "dirty"
    ReconciliationUnknown      ReconciliationState = "unknown"
    ReconciliationReleased     ReconciliationState = "released"
)
```

```go
type RunWorktreeReconciliation struct {
    RunID, Outcome, Path, Branch, TargetBranch string
    RunHead, TargetHead, Reason                 string
    State                                      ReconciliationState
}

func InspectTerminalRun(ctx context.Context, run store.Run) (RunWorktreeReconciliation, error)
func ApplyTerminalRun(ctx context.Context, result RunWorktreeReconciliation) error
```

`InspectTerminalRun` returns `unknown` with a bounded reason for expected uncertainty. It returns an error only when the inspection operation itself cannot continue, such as an invalid repository root. This distinction lets bulk scans preserve one uncertain Run without hiding the rest.

### Data Models

No Run schema change is required:

- `Run.WorkDir` supplies the recorded Run Worktree path.
- `roundfix/run-<id>` remains the deterministic Run Branch.
- `Run.LocalBranch` is the recorded target branch for spec Runs.
- `Run.State` supplies the historical terminal outcome.
- `Run.GitRoot` scopes Git operations.

The reconciliation event payload contains `previous_outcome`, `current_outcome`, `classification`, `run_branch`, `run_head`, `target_branch`, `target_head`, `worktree`, and `action`. For Integration Pending, the guarded store operation changes the outcome to Clean and appends this evidence transactionally. For all other outcomes, a Daemon-source status event records the proof without changing the Run row.

### API Contracts

```text
roundfix reconcile [run-id] [--apply] [--format text|json]
```

- A Run ID may identify only a terminal spec Run. Active Runs, review Runs, and missing Runs fail Preflight Validation with no Git mutation.
- Without a Run ID, Roundfix requires a Git repository and scans all terminal spec Runs whose `git_root` matches the current repository.
- Default format is text. JSON uses a versioned envelope with `schemaVersion`, `mode`, `repository`, `results`, and summary counts. Result order is newest Run first, matching `runs list`.
- Dry-run performs no Run Database or Git writes and prints the exact `--apply` command.
- Apply acts only on results already classified `safe` during that invocation. It never accepts a force flag or user assertion that bypasses proof.
- Dirty detection uses porcelain status including untracked files. A missing worktree with a present branch can still be classified from the branch and target refs; a present dirty worktree always wins as `dirty`.
- Apply uses `git worktree remove` without force. After ancestry proof and successful worktree removal, it deletes the Run Branch through the existing porcelain helper. A failure preserves the remaining ref/path and returns the Run failure exit code.
- `unintegrated` and `dirty` are expected report states and do not make a complete scan fail. `unknown` produces a non-zero result only when caused by an operational inspection error; missing proof alone is a successful preserved result.

Run listing preserves its existing stdout rows. When `--state active` hides retained terminal Run Worktrees, stderr becomes:

```text
(2 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
```

The note counts Runs with an existing recorded worktree path or Run Branch, not all terminal history. Terminal and all-state views use the same stderr guidance when their visible rows include retained Run Worktrees, while `roundfix reconcile` remains the only classification surface.

## Coverage Map

- Goal 1, Stories 1 and 4 → shared classifier and positive-proof rules.
- Goal 2, Stories 2 and 4 → dry-run/apply CLI and JSON envelope.
- Goal 3, Story 3 → exact retained-worktree notes and Reconcile Command pointer without changing Run row stdout.
- Goal 4, Story 5 → guarded Integration Pending reconciliation plus unchanged outcomes for other Runs.
- Core Feature 9 → `PruneTerminalReport` delegates to the shared classifier.

## Integration Points

- **Git worktrees and refs** through the existing `gitRunner`: worktree list/status/remove, revision resolution, `merge-base --is-ancestor`, and porcelain branch deletion.
- **Run Database** for repository-scoped terminal Runs, guarded Integration Pending reconciliation, and audit events.
- No network or GitHub access is required. The target is the local branch recorded when the spec Run began.

## Testing Approach

- Table-driven worktree tests create real temporary repositories for all five classifications, including tracked dirt, untracked dirt, missing target branch, missing Run Branch, target divergence, and Run Branch already reachable from a later merge.
- Safety tests assert that only `safe` reaches cleanup and that `ApplyTerminalRun` refuses a stale result when either head changed after inspection. Apply must re-resolve both heads and cleanliness immediately before mutation.
- CLI tests cover one Run, repository scan, JSON schema, dry-run byte stability, apply, idempotent repeat, invalid Active/review Run selectors, and mixed safe/unsafe results.
- Store integration asserts Integration Pending → Clean with evidence and confirms Unresolved remains Unresolved after safe cleanup.
- Run listing tests cover exact retained counts in active, terminal, and all-state views and prove existing stdout rows remain byte-stable.
- Automatic reaper regression reproduces the Vortex case: a changed Run Branch that is already ancestor of the target is removed, while a unique changed branch remains.

## Build Order

1. Shared reconciliation states and inspection classifier with real-Git fixtures.
2. Stale-proof-safe apply operation and automatic reaper migration (depends on: 1).
3. Store reconciliation event wiring and Integration Pending promotion through spec 0037's guarded transition (depends on: 1, 2, spec 0037).
4. Reconcile Command parsing, text output, JSON envelope, dry-run/apply, and exit contracts (depends on: 1, 2, 3).
5. Run listing retained-worktree notes with byte-stable stdout rows (depends on: 1).
6. User guide, command help, CONTEXT vocabulary, and finding traceability (depends on: 4, 5).
7. Dedicated tooling-only update of
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, with
   direct byte-identical edits and read-only sync verification (depends on: 6).

## Risks & Considerations

- Inspection evidence can become stale between display and apply. Apply rechecks cleanliness and both resolved heads immediately before any mutation and refuses if they changed.
- A deleted target branch makes the result `unknown`; Roundfix does not guess another branch or remote ref.
- The Run Event Journal can later be pruned by retention. The current outcome and released Git state remain authoritative; audit detail is retention-bound by design.
- Repository scans run local Git commands per retained terminal Run. The default scope is one repository and terminal residue is expected to be small; cross-repository bulk operation is excluded.
- Spec 0037 must land first because this command depends on compare-and-set completion and the guarded Integration Pending transition.

Rollout is additive: the new command and stderr guidance can ship without migrating the Run Database. Apply produces ordinary Git and existing Run state changes, so rollback removes only the first-party command; already released work stays safely integrated and missing worktrees classify as released if a newer binary is restored. There is no automatic mutation to roll back.

## Decisions

- Dedicated Reconcile Command, dry-run by default, with one explicit `--apply` mutation flag.
- Positive proof requires both clean worktree state and target-branch ancestry.
- The automatic reaper and CLI share one classifier.
- Protected Roundfix Skill publication is isolated in one Task whose changed
  files are the two authorized `SKILL.md` paths and its own Task file; it does
  not run the broad `make skills-sync` mutation target.
- See [ADR-0053](../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md).
