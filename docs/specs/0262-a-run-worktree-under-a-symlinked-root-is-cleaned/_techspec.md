---
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
prd: _prd.md
created: 2026-10-12
---

# A Run Worktree under a symlinked root is cleaned — Technical Spec

## Executive Summary

One rule and one comparison in `internal/worktree` decide every case the sweep
found. `recordedWorktree` stops refusing any symbolic link and instead resolves
a recorded path within its Worktree Root: links are followed above the
repository directory, and the repository and worktree directories must be
real. `canonicalPath`, behind every `samePath`, resolves the longest existing
prefix of a path instead of falling back to its lexical form, so a worktree
whose directory is gone still matches the resolved path Git keeps. Every
consumer reaches both through the existing functions, so no caller changes.
The timing-dependent disposal lock tests are out of scope: Spec 0259's
corrective task_06 fixes them.

The trade-off accepted is that the Worktree Root is read from the recorded
path's layout, `<worktree.location>/<repo-slug>/<name>`, and not from the
current configuration. A link inside the location is trusted as the location
itself; a link at the repository or worktree directory is not. Reading the
configured location would bind the rule to today's configuration, strand Runs
created under an earlier location, and touch every caller and its tests
(ADR-0268).

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes; recorded paths are read and never rewritten. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; tests use temporary homes, repositories and the local `git`. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0268 governs every decision below.
  ADR-0053: "is read-only unless `--apply` is supplied" and "Ambiguous, dirty,
  and unintegrated work stays intact". ADR-0023: "living under Roundfix Home,
  recorded on the Run row". ADR-0259: "Every
  top-level test in a Parallel Test Package calls `t.Parallel()` as its first
  statement". ADR-0271: "A test bounds by real time only what it asserts"; the new tests give no child, goroutine or Git a real-time allowance they do not assert. ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0244: "A Spec now carries a Glossary Declaration".
  ADR-0257: "Inside a Run the Daemon is the only full gate". ADR-0258: "A Spec
  is authored against the QA rerun classes". ADR-0161, ADR-0212 and ADR-0232:
  merge-evidence release proofs, unchanged and now reached under a linked
  root. ADR-0170: carry-forward, untouched. ADR-0141: does not apply, review
  Runs record no Run Worktree. ADR-0237: the Delivery Retry merge record is
  read as before. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the cycle AF grant of 2026-10-12 ("Sim, como no AE") covers the
  Governed Paths each Spec declares, and none is declared. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned/_authorization.md`;
  bounded files: none (`paths: []`).

## Current behavior

Measured on 2026-10-12 at `94fbf85b` on the maintainer's arm64 Mac, Git 2.56.0.

**Git registers the resolved path.** A worktree added at
`<tmp>/link/slug/run-x`, where `<tmp>/link` points to `<tmp>/real`, appears in
`git worktree list --porcelain` as `worktree <tmp>/real/slug/run-x`.

**Reconciliation refuses the recorded spelling.** A characterization test,
removed after use, created a Run Worktree with `worktree.Create` under a
location that is a symbolic link and recorded `ref.Path` as the Run's
`WorkDir`. `InspectTerminalRun` returned state `unknown` with reason
`recorded Run Worktree path is unsafe`; `PruneTerminalReport` pruned nothing
and returned no error; `CountRetainedTerminalRuns` counted it as retained.
The cause is `recordedWorktree`, which calls `pathContainsSymlink` and returns
`recordedWorktreeUnsafe` at the first link in any component, before it reaches
the comparison with the registered paths, and before it can notice that the
worktree directory is gone.

**The sweep.** Every production comparison of a recorded or derived worktree
path with a registered path or the worktree location was read:

| Site | Comparison | Under a symbolic-linked root |
| --- | --- | --- |
| `recordedWorktree` | `pathContainsSymlink` refusal, then `canonicalPath` equality | refuses; the root cause |
| `inspectTerminalRunMerged`, `revalidateTerminalRunApply`, `ApplyTerminalRun` | through `recordedWorktree` | `unknown`; apply never reaches removal |
| `ApplyRunBranchCandidate`, `PruneTerminalReport` | through `InspectTerminalRun` | refuses; prune skips silently, Task Worktrees leak |
| daemon `ClassifySupersededBranch`, `discardSupersededBranch` | through `InspectTerminalRun` | refuses the discard |
| CLI reconcile, `releaseMergedSpecRuns`, `pruneTerminalRunWorktreeDebris`, `forceStopRun` | through the functions above | `unknown`, `warning: cleanup failed`, silent preserve |
| `remainingTerminalRunSurface` | `samePath`, then `os.Stat` fallback | right answer through the fallback |
| `UseItem` | `samePath` with the registered paths | directory gone: no match, stale registration kept, `worktree add` fails |
| `CleanupItem` | `samePath` on branch and path registrations | directory gone: refuses "checked out at X, expected Y" |
| `inspectItemBranch` | `samePath` with the derived path | directory gone: refuses "not at derived path" |
| `recordedGitDirectory` | `pathContainsSymlink` on the recorded Git root | refuses; out of scope |
| staging, disposal, `CleanupClean`, `CleanupTask`, `terminalCandidates`, `countRetainedTerminalRuns` | none, or Git-to-Git | correct |

`samePath` compares `canonicalPath` of both sides, and `canonicalPath` falls
back to `filepath.Clean` when `EvalSymlinks` fails, so a missing directory in
the link spelling never equals the resolved path Git keeps. Neither `doctor`
nor `gc` compares worktree paths. `internal/config` already has the
prefix-resolving form, `canonicalizeResolvedPath`, for review roots.

## System Architecture

- **Worktree Root resolution** (`internal/worktree/worktree.go`): a new
  unexported helper replaces the `pathContainsSymlink` call inside
  `recordedWorktree`. `pathContainsSymlink` stays for the recorded Git root.
- **Prefix-resolving comparison**: `canonicalPath` delegates to the review
  root's algorithm, exported from `internal/config` as
  `CanonicalizeResolvedPath`; its only config caller is renamed.
- **Consumers**: no change to `internal/cli`, `internal/daemon` or
  `internal/delivery` production code; their tests prove the rule end to end.

## Implementation Design

### Interfaces

```go
// internal/worktree/worktree.go
// resolveWithinWorktreeRoot returns the resolved path of a recorded worktree,
// or the state that ends recordedWorktree early.
func resolveWithinWorktreeRoot(path string) (string, recordedWorktreeState)
func canonicalPath(path string) string // longest existing prefix resolved

// internal/config/config.go (renamed from canonicalizeResolvedPath)
func CanonicalizeResolvedPath(path string) string
```

### Data Models

None. No record, table or report field changes.

### API Contracts

1. API Contract: the Worktree Root rule. `recordedWorktree` keeps its checks
   for an empty, multi-line, relative or unclean path. Then, with `parent` the
   path's directory and `root` the parent's directory: `os.Lstat` of `parent`
   and of the path. A missing `parent` or path is `recordedWorktreeAbsent`. A
   `parent` or path that is a symbolic link, dangling or not, or is not a
   directory, is `recordedWorktreeUnsafe`. `filepath.EvalSymlinks(root)`
   failing with not-exist is absent; any other failure is unsafe. The path's
   `EvalSymlinks` result must equal the resolved root joined with the base
   names of `parent` and the path, or it is unsafe. Otherwise the resolved
   path is compared with `canonicalPath` of each registered worktree; a match
   is `recordedWorktreePresent`, none is `recordedWorktreeUnregistered`. The
   reason strings, the `Path` reported and every later proof are unchanged.
2. API Contract: the prefix-resolving comparison. `canonicalPath` cleans the
   path, resolves the longest prefix that exists with `EvalSymlinks`, and joins
   the missing components verbatim; a path with no resolvable prefix returns
   its clean form. For an existing path the result is unchanged.
   `config.CanonicalizeResolvedPath` is the same function that
   `reviewRootInsideHistory` already calls.
3. API Contract: consumers. Under a symbolic-linked worktree location, each of
   these behaves as under a real directory: `roundfix reconcile` and
   `roundfix reconcile --apply`; `releaseMergedSpecRuns` after a Delivery
   Queue merge; `PruneTerminalReport`, reached by `implement` and Force Stop,
   including a released Run's Task Worktrees; `ClassifySupersededBranch` and
   its discard; `ApplyRunBranchCandidate`; `UseItem` and `CleanupItem` for an
   item worktree whose directory is gone. No flag, report field, reason or
   exit code changes.
4. API Contract: documentation. The reconcile user guide gains a paragraph:
   when `worktree.location` is a symbolic link, a recorded Run Worktree is
   resolved within its Worktree Root, a recorded path whose repository or
   worktree directory is a link stays `unknown`, and `roundfix reconcile`
   followed by `roundfix reconcile --apply` releases the Runs an earlier
   version kept with `recorded Run Worktree path is unsafe`. The configuration
   guide's worktree section gains one sentence naming the Worktree Root.

### Surface Transcripts

None. No command's output lines, flags or exit codes change; a Run under a
symbolic-linked root is reported with the classification and reason it already
gets under a real directory, and the CLI tests assert them.

### Invariants

1. A recorded path whose repository or worktree directory is a symbolic link
   is never `present`.
2. For a recorded path without any symbolic link, `recordedWorktree` returns
   what it returns today.
3. `canonicalPath` of an existing path is unchanged.
4. Reconciliation never removes anything it would not remove if the worktree
   location were the real directory.

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → Worktree Root
  resolution (task_01), consumers' tests (task_03).
- Goal 2, User Story 3, Core Feature 2, Success Metric 2 → Worktree Root
  resolution refusal cases (task_01).
- Goal 3, User Story 4, Core Feature 3, Success Metric 3 → prefix-resolving
  `canonicalPath` (task_02), item and release tests (task_03).
- Goal 4, User Story 2, Core Feature 4, Core Feature 5, Success Metric 1 →
  CLI reconcile tests and the reconcile and configuration guides (task_03).
- Core Feature 6 → `CONTEXT.md` **Worktree Root** (task_01).

## Absence Evidence

1. Worktree Root resolution — command: `grep -rn 'Worktree Root' internal CONTEXT.md docs/user-guide`; result: no match at `94fbf85b`; the only symbolic-link rule for worktree paths is `pathContainsSymlink`, which refuses every link.
2. Prefix-resolving comparison in the worktree package — command: `grep -n 'EvalSymlinks' internal/worktree/worktree.go`; result: one match, `canonicalPath`, which falls back to the lexical path; the prefix-resolving form exists only as the unexported `canonicalizeResolvedPath` in `internal/config/config.go`, which this design reuses.
3. Reconcile guide paragraph — command: `grep -n 'symbolic link\|symlink' docs/user-guide/commands/reconcile.md`; result: no match; the configuration guide mentions symlinks only for Spec artifacts.

## ADR Alternatives

- ADR-0268: this Spec's own decision; it rejects threading the configured location, accepting any link and recording resolved paths.
- ADR-0053: no rejected alternative is taken; reconciliation stays proof-based and read-only without `--apply`, and only the set of paths that count as the recorded worktree widens.
- ADR-0023: no rejected alternative is taken; Run Worktrees still live one per Run under Roundfix Home, and the recorded path is read as written.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first and no Sequential Test is added.
- ADR-0271: no rejected alternative is taken; the new tests bound by real time only what they assert, and no test file gains a Real-Time Budget.
- ADR-0256: its rejected model-judging alternative stays rejected; the Coverage Review records why the reconcile skill text needs no change.
- ADR-0244: no rejected alternative is taken; **Worktree Root** follows the declaration contract.
- ADR-0257: no rejected alternative is taken; Task agents run focused tests and the Daemon runs the declared Verification.
- ADR-0258: no rejected alternative is taken; Verification names the consuming suites and the QA rows reach no network.
- ADR-0161: no rejected alternative is taken; a merged Spec's Runs are still released only on the merged-head proof, which now also runs for a Run under a linked root.
- ADR-0212: no rejected alternative is taken; supersession of Spec-directory work and declared leftovers is proved as before.
- ADR-0232: no rejected alternative is taken; merge evidence still decides the release of Runs and item branches, and item worktree comparison only stops depending on the link spelling.
- ADR-0170: no rejected alternative is taken; carry-forward is untouched.
- ADR-0141: no rejected alternative is taken; review Runs keep executing in the user checkout, outside this rule.
- ADR-0237: no rejected alternative is taken; the merge record a Delivery Retry writes stays the release evidence it was.

## Integration Points

None outside the repository. Git is the local `git` binary the package
already runs; its behavior under a symbolic link was measured above.

## Testing Approach

New tests, each calling `t.Parallel()` first, all on temporary homes and
repositories:

- `internal/worktree/worktree_root_test.go` (task_01):
  `TestInspectTerminalRunResolvesASymlinkedWorktreeRoot`,
  `TestApplyTerminalRunReleasesAWorktreeUnderASymlinkedRoot`,
  `TestPruneTerminalReportPrunesUnderASymlinkedRoot` (a released Run with a
  Task Worktree),
  `TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot` with
  subtests `repository-directory-link`, `worktree-directory-link` and
  `dangling-worktree-link`, and
  `TestInspectTerminalRunTreatsAMissingWorktreeUnderASymlinkedRootAsAbsent`.
  Each builds the location as a link to a real directory with `worktree.Create`.
  `TestInspectTerminalRunUnsafePath` stays unchanged and must still pass.
- `internal/worktree/worktree_path_test.go` (task_02):
  `TestCanonicalPathResolvesTheExistingPrefix` with subtests `existing`,
  `missing-leaf`, `missing-parent` and `no-existing-prefix`;
  `TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot`;
  `TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot`.
- `internal/cli/reconcile_symlinked_root_test.go` (task_03):
  `TestReconcileClassifiesARunUnderASymlinkedWorktreeRoot`,
  `TestReconcileApplyReleasesARunUnderASymlinkedWorktreeRoot`,
  `TestReleaseMergedSpecRunsCleansUnderASymlinkedWorktreeRoot`; each records
  the Run Worktree in the link spelling of a linked location.
- `internal/daemon/reconcile_symlinked_root_test.go` (task_03):
  `TestClassifySupersededBranchUnderASymlinkedWorktreeRoot`.
- Each Task runs its changed packages once with `-race` (ADR-0259).

## Build Order

1. Worktree Root resolution in `recordedWorktree`, its tests, ADR-0268's term
   **Worktree Root** in `CONTEXT.md` (task_01).
2. Prefix-resolving `canonicalPath` reusing the config helper, and the item
   tests (task_02) (depends on: 1).
3. Consumers' tests in `internal/cli` and `internal/daemon`, the reconcile and
   configuration guides, the Coverage Review and the Behavior Surface Record
   (task_03) (depends on: 1, 2).
4. QA gate (depends on: 3).

task_01 and task_02 both edit `internal/worktree/worktree.go`, so they run in
series; task_03 is the Integration Checkpoint for both, whose Verification runs
the consuming suites.

## Risks & Considerations

- **A test that relied on the lexical fallback.** A test comparing
  `canonicalPath` of a missing macOS temporary path with its `/var` spelling
  will now see `/private/var`. task_02 runs the whole `internal/worktree`
  package and task_03 the consumers; a broken comparison is fixed in the test
  only when it compared a missing path by accident, and the Task names it.
- **Trusting the location's link.** A link placed at or above the location
  redirects all worktrees together; that is the configuration the user chose.
  Links below it are refused, so one Run's directory cannot be swapped.
- **Old reports.** Runs an earlier version marked with the unsafe reason carry
  no stored state; the next inspection classifies them afresh.

## Glossary

- adds: **Worktree Root**
- not a term: **Worktree Root resolution** — the helper that applies the rule
- not a term: **Prefix-resolving comparison** — the behavior of one helper
- not a term: **Git registers the resolved path** — a heading of the measurements
- not a term: **Reconciliation refuses the recorded spelling** — a heading of the measurements
- not a term: **A test that relied on the lexical fallback** — a risk heading
- not a term: **Trusting the location's link** — a risk heading
- not a term: **Old reports** — a risk heading
- not a term: **Current Behavior** — a section heading

## Decisions

- A recorded worktree path is resolved within its Worktree Root read from its
  layout; see ADR-0268.
- `canonicalPath` resolves the longest existing prefix by reusing the config
  helper; see ADR-0268.
- No consumer's production code changes; task_03 proves them end to end.
