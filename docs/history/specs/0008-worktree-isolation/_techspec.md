---
spec: 0008-worktree-isolation
prd: _prd.md
created: 2026-07-05
---

# Worktree Isolation — Technical Spec

## Executive Summary

Execution moves from the user's checkout into a per-Run git worktree; the
consequential trade-off is a **new integration boundary at Run end** — today
commits land on the user's branch by construction, after this Spec they land
there through an explicit porcelain protocol that can refuse (Integration
Pending). That cost was accepted because the alternative — plumbing ref
updates under a live checkout — was empirically shown to corrupt the user's
view (phantom staged deletions from `git update-ref`), and because the
boundary buys structural fixes for the three chronic findings (multi-writer
sweeps, dirty-tree rejection, snapshot fragility). Everything inside the Run
is deliberately unchanged: the engines, before-snapshot rule, commit
derivations, and ownership ADRs operate as-is, just pointed at a different
directory. All git mechanics below were verified by experiment during
exploration.

## System Architecture

- **New: `internal/worktree`** — the one new package: create (named Run
  Branch, Roundfix Home path, provisioning copy-list), remove/keep, prune
  sweep, and the integration protocol. Cohesive and small; the daemon's git
  wrappers stay where they are.
- `internal/store` — schema v6: `runs.work_dir` (nullable); recorded at
  creation, read by attach/settle/reports.
- `internal/cli` — implement and resolve/watch wire `WorkDir` from the
  created worktree instead of the git root; preflight demotion; outcome
  handling calls integration; settle retargets; reports read the worktree.
- `internal/daemon` — no engine logic changes: `TaskPlan.WorkDir` /
  `CyclePlan.GitRoot` simply receive the worktree path (the exploration
  mapped every consumer riding on those two fields).
- `internal/tui` + attach — `LiveRunView` gains `WorkDir` (falls back to
  `GitRoot` when empty/gone); cockpit task refresh and detail reads use it.
- `internal/config` — `worktree.copy` (string list, default empty); the
  Artifact Directory *builtin* default moves to Roundfix Home
  (`~/.roundfix/artifacts/<repo-id>`), explicit values untouched.
- `internal/agent` — nothing: cwd and `--cwd` already flow from the plan's
  work dir.

## Implementation Design

### Interfaces

```go
// internal/worktree
type Ref struct {
    RunID   string
    Path    string // ~/.roundfix/worktrees/<repo-id>/<run-id>
    Branch  string // roundfix/run-<id>
    UserRoot string // the checkout that owns the work-target lock
}

func Create(ctx context.Context, userRoot, runID, headSHA string, copyList []string) (Ref, error)
func Integrate(ctx context.Context, ref Ref, targetBranch, runSHA string) (IntegrationResult, error)
func CleanupClean(ctx context.Context, ref Ref) error   // remove worktree + branch
func PruneTerminal(ctx context.Context, userRoot string, isTerminalClean func(runID string) bool) error

type IntegrationResult struct {
    Mode   string // "ff-merge" | "branch-move" | "pending"
    Reason string // pending only: overlap | diverged | dirty-target detail
}
```

### Worktree creation (verified mechanics)

`git worktree add -b roundfix/run-<id> <path> <headSHA>` — the named branch
keeps `branch --show-current`, InspectGit, and human inspection working
(detached HEAD errors preflight's branch detection). Path under Roundfix
Home so the repo tree never sees it. Copy-list entries are copied relative
paths, missing sources are per-file stderr notes, never failures. Cost is
~0.1–0.5s plus whatever untracked state the verification needs (the reason
the copy-list exists).

### Integration protocol (ADR-0024, each case verified)

1. Resolve where `targetBranch` is checked out (`git worktree list`).
2. Checked out in the user root → `git -C <userRoot> merge --ff-only
   <runSHA>`: succeeds on clean and non-overlapping-dirty checkouts;
   refuses with the branch unmoved on overlap ("local changes would be
   overwritten") and on divergence ("Not possible to fast-forward") — both
   map to `pending` with the reason.
3. Not checked out anywhere → `git merge-base --is-ancestor <targetBranch>
   <runSHA>` then `git branch -f <targetBranch> <runSHA>`; non-ancestor →
   `pending/diverged`.
4. NEVER `git update-ref` on a checked-out branch — the exploration showed
   it leaves phantom staged deletions that a routine `git commit -a` turns
   into a reversion of the Run's work.
5. `pending` → Run ends Integration Pending: worktree and Run Branch kept,
   report prints `git merge --ff-only roundfix/run-<id>`; Final Push is
   skipped (pushing an unintegrated sha would strand the local branch).
   Successful integration precedes any Final Push / push-at-Clean.

### Run flows

- **Preflight (user root, mostly unchanged):** PR/branch resolution,
  default-branch veto, Active Run locks — all still against the user
  checkout (`CreateRun.GitRoot` keeps its lock meaning). The implement
  dirty-tree rejection becomes one stderr note stating that overlapping
  local changes will end the Run Integration Pending. A worktree-debris
  sweep (`git worktree prune` + removal for terminal-Clean Runs) runs here.
- **Execution:** create worktree → `TaskPlan.WorkDir`/`CyclePlan.GitRoot` =
  worktree path (Agent cwd, prompts, verification, snapshots, commits,
  task-file reads/writes, QA — all follow automatically per the exploration
  map). Before-snapshot machinery unchanged: inside the worktree it still
  isolates failed-Task dirt between Tasks.
- **Outcomes:** Clean → integrate → on success remove worktree + branch
  (then Final Push / auto-push if configured); integration refusal →
  Integration Pending (new terminal state, exit 1 mapping). Unresolved /
  Failed / Stopped → keep worktree + branch, print the path. Settle:
  resolve the Run's `work_dir`, verify and commit there, then the same
  integration protocol; its stage-everything contract now scopes to the Run
  Worktree (strictly safer).
- **Readers:** cockpit/attach/report task reads use `work_dir` with
  `git_root` fallback; attach on a pruned-worktree terminal Run reads the
  user root (completed task files live in the integrated commits).

### Data Models

Schema v6: `runs.work_dir` nullable text, empty for legacy rows; migration
preserves v5 rows and locks. New terminal state string `IntegrationPending`
in the states enum (additive; exit-code family 1). Config: `worktree.copy`
list; Artifact Directory builtin default change (explicit configs and
`ResolveArtifactDirectory` semantics untouched).

### API Contracts

Commands and flags unchanged. Contract deltas, all additive or documented:
implement's dirty-tree hard error becomes a note; Run headers/report lines
name the Run Worktree; Integration Pending joins the outcome vocabulary
(stdout outcome line + exit 1); settle reads the kept worktree. Everything
else — stdout shapes, exit codes, journal kinds — byte-stable.

## Coverage Map

- Story 1 → worktree execution retarget + engines untouched (findings 15)
- Story 2 → preflight demotion (finding 10's friction)
- Story 3 → Clean integrate + cleanup (ADR-0024, ADR-0023)
- Story 4 → Integration Pending path (ADR-0024)
- Story 5 → kept worktrees + settle retarget (finding 10)
- Story 6 → LiveRunView/attach WorkDir reads
- Story 7 → `worktree.copy` provisioning

## Integration Points

git only — worktree add/remove/prune, ff-merge, merge-base, branch -f, all
through the existing runner wrappers; no new external systems.

## Testing Approach

Existing rigs plus real temp repos (the hermetic helpers from 0003):
`internal/worktree` unit tests reproduce the exploration matrix — creation
on checked-out branches, ff-merge on clean/non-overlapping/overlapping/
diverged checkouts (asserting the branch is unmoved on refusal), branch-move
with and without ancestry, dirty-worktree remove refusal, crash-prune. A
regression test pins the update-ref hazard indirectly: after any pending
outcome, the user checkout's `git status` is asserted clean of phantom
entries. Engine tests re-run against a worktree WorkDir unchanged; CLI
end-to-end tests cover both engines with concurrent user commits (the
multi-writer proof), all four outcome paths, settle-from-worktree, and the
copy-list. Migration v5→v6 fixture. `-race` on daemon. Full suite is the
net; deliberate updates limited to the demoted preflight and new report
lines.

## Build Order

1. `internal/worktree` package: create/remove/prune + integration protocol
   with the verified git matrix (no deps)
2. Store schema v6 (`work_dir`) + config `worktree.copy` + artifact-dir
   builtin default move (no deps)
3. Implement path on worktrees: wiring, preflight demotion, outcomes,
   Integration Pending state, settle retarget (depends on: 1, 2)
4. Resolve/watch path on worktrees: wiring, integration before Final Push
   (depends on: 1, 2)
5. Live Run View and Attach read the execution workspace (depends on: 3)
6. Docs and skill sync: outcome vocabulary, worktree lifecycle, config keys
   (depends on: 3, 4, 5)

## Risks & Considerations

- The integration boundary is the behavior change users will feel; the
  Integration Pending message must carry the exact command (tested
  verbatim) and the report must always name the kept worktree.
- Worktrees accumulate on crashes: the preflight sweep plus `git worktree
  remove` refusing dirty deletion (free guard) contain it; the sweep must
  never touch worktrees of non-terminal or non-Clean Runs.
- Cold worktrees lack untracked state: `worktree.copy` is the escape hatch
  and the docs task must call this out as the one gate-behavior difference.
- Watch reuses one worktree across Rounds within its single Run — fetch
  artifacts stay in the shared Artifact Directory, unaffected.
- Windows rename/path semantics are out of scope (consistent with the
  repo's platform posture to date).

## Decisions

- Worktree-per-Run, named Run Branch, Roundfix Home placement. See
  ADR-0023.
- Porcelain-only two-case integration; Integration Pending terminal state;
  push only after integration. See ADR-0024.
- Engines and ownership ADRs untouched by construction — isolation is a
  pointer change, not an engine change.
- `worktree.copy` explicit list; artifact-dir builtin default leaves the
  repo tree; glossary gains the three terms (already committed with the
  PRD).
