---
task: task_01
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
status: completed
type: backend
complexity: medium
---

# Task 01: The Baseline inventory keeps only the nested carriers Git lists

## Overview

The Baseline inventory walk asks Git, once, on the first nested `AGENTS.md`
or `CLAUDE.md` it meets, which nested carriers belong to the repository, and
inventories only those. An ignored copy and a carrier inside a linked
worktree, a nested repository or a submodule then raise no
`baseline.inventory.nested-carrier-conflict`, while a tracked nested carrier
is inventoried exactly as before. This answers the message reported by
quaestor on 2026-11-13
([`references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md`](references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md)).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `InspectRepository` of a temporary repository with an ignored linked
   worktree under `.storage/` reports `baseline.inventory.nested-carrier-conflict`
   for the worktree's `AGENTS.md`, and `InspectRepository` runs no
   `ls-files` command.
2. MUST add `ExecGitRunner.RunGitRaw` and make `RunGit` return its trimmed
   output, as `_techspec.md` → API Contract 1 says; `RunGit`'s behavior MUST
   NOT change.
3. MUST add `listRepositoryCarriers` in `internal/baseline/repository.go` with
   the exact arguments, the untrimmed read and the error text of API
   Contract 2.
4. MUST change `inventoryBuilder.walk` as API Contracts 3 and 4 say: the
   listing runs lazily on the first nested candidate, a candidate outside the
   listing is left out, and the body that turns a path into a carrier moves
   unchanged into `inventoryCarrier`. Root carriers and every path under
   `docs/agents/` MUST NOT depend on the listing.
5. MUST give `inspectRepositorySnapshot` and `verifyAppliedPlanState` a
   leading `context.Context` and pass it from every caller named in
   `_techspec.md` → System Architecture; `InspectRepository` passes its
   runner, every other caller `nil`; the two transaction constructors in
   `internal/baseline/transaction.go` wrap `verifyAppliedPlanState` in a
   closure over their context. The only existing test edit is the
   `inspectRepositorySnapshot` call of the over-budget subtest in
   `internal/baseline/repository_test.go`.
6. MUST add `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists` and
   `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments` to
   `internal/baseline/repository_test.go`, each calling `t.Parallel()` first,
   with the fixtures and assertions of `_techspec.md` → Testing Approach. The
   first sets `core.excludesFile` in the fixture's local Git config to an
   empty file under the test's temporary directory; the submodule is added
   from a local repository with `git -c protocol.file.allow=always submodule
   add`. No test reaches the network, reads the host's home or waits on
   elapsed time.
7. MUST NOT change any other existing test, any output text, any exit code,
   `inventoryPathIgnored`, `isBoundedCarrierPath`, the inventory limits or
   the snapshot schema.

## Subtasks

- [ ] Record the characterization.
- [ ] Add `RunGitRaw` and `listRepositoryCarriers`.
- [ ] Ask Git lazily in the walk and thread the context.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] An ignored copy, a linked worktree under an ignored or a visible path, a nested repository and a submodule add no carrier, no warning and no change to the snapshot digest.
- [ ] A tracked nested carrier, a tracked one under an ignored directory and an untracked one Git does not ignore keep their carriers and warnings.
- [ ] The listing runs once, with the exact arguments, only when a nested carrier exists, and a failed listing returns the error of API Contract 3.

## Context

- instruction: `docs/adr/0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md`
- instruction: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- instruction: `internal/spec/qa_evidence_links.go`
- interface: `internal/baseline/repository.go`
- interface: `internal/baseline/repository_test.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`
- interface: `internal/baseline/transaction.go`

## Consumers

- The Baseline inventory and `inspectRepositorySnapshot` — consumers: `./internal/baseline`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 40m -v ./internal/baseline 2>&1)"; status=$?; for name in TestNestedCarrierInventoryKeepsOnlyCarriersGitLists TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths TestRepositoryInspectionUsesNarrowReadOnlyGitCommands; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf '%s\n' "$out" | tail -40; echo "missing pass: $name" >&2; exit 1; }; done; test "$status" -eq 0 || { printf '%s\n' "$out" | grep -E -- '--- FAIL|FAIL' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/baseline` — expected: exit 0. Before this Task the two new tests do not exist; afterwards they pass, every test of `./internal/baseline` still passes, and the package builds for Windows.

- `grep -qF -- 'list repository instruction carriers' internal/baseline/repository.go || { echo 'listing missing' >&2; exit 1; }; go test -count=1 -timeout 40m -run 'Baseline' ./internal/cli` — expected: exit 0. Before this Task the listing does not exist; afterwards every `Baseline` test of `./internal/cli`, which drives `baseline plan`, `baseline update` and `baseline apply` over temporary repositories, still passes.

## References

- `_prd.md` → Goals; Core Features 1 and 2; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1
- ADR-0399
- [The adopted quaestor message](references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md)

## Result

Skills loaded: implement-task, coding-guidelines, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

### Starting-tree characterization

Before repository edits, a temporary Go overlay appended
`TestTask0390StartingTreeCharacterization` to the existing inspection tests.
It created a committed temporary repository with `.storage/` ignored and a
local detached linked worktree at `.storage/simulation/wt-ignored`.
`InspectRepository` reported `baseline.inventory.nested-carrier-conflict`
for `.storage/simulation/wt-ignored/AGENTS.md`. A recording runner delegated
to `ExecGitRunner` and observed exactly the two identity calls:
`rev-parse --show-toplevel --show-object-format --verify HEAD^{commit}` and
`rev-list --max-parents=0 HEAD`; `ls-files` calls: **0**.

Focused command: `go test -overlay /private/tmp/roundfix-0390-characterization/overlay.json -count=1 -run '^TestTask0390StartingTreeCharacterization$' -v ./internal/baseline`
— exit 0, characterization assertions passed. The overlay and fixture were
outside the checkout; no production code had changed.

### Implementation and acceptance evidence

- `ExecGitRunner.RunGitRaw` keeps the existing invocation, environment and
  error handling and returns stdout bytes unchanged. `RunGit` now trims that
  output, retaining its existing behavior. `listRepositoryCarriers` uses the
  exact narrow NUL-delimited command of API Contract 2, prefers the optional
  raw interface, and wraps listing errors with the required text.
- The walk lists lazily on its first nested candidate, caches the result,
  and omits unlisted candidates before carrier metadata or bytes are read.
  Root carriers and `docs/agents/` bypass the listing. The original carrier
  body moved into `inventoryCarrier`; only its early return changes to the
  required void return. Context reaches inventory from inspection, planning,
  validation and apply, including both transaction verification closures.
- Acceptance criterion 1: the real-Git regression adds an ignored copy,
  ignored and visible linked worktrees, a nested repository and a submodule
  from a second disposable local repository. Before and after, nested
  carriers and conflict paths are exactly ` lead/AGENTS.md`,
  `forced/AGENTS.md`, `packages/api/AGENTS.md`, and
  `packages/web/CLAUDE.md`, in order. The complete warning list and snapshot
  digest are equal, and neither inspection has blocking findings.
- Acceptance criterion 2: the same regression retains the committed nested
  carrier, the force-added carrier under an ignored directory and both
  untracked non-ignored carriers. Its local `core.excludesFile` names an
  empty temporary file. The leading-space untracked path exercises the
  untrimmed real-Git output. Submodule setup uses
  `git -c protocol.file.allow=always submodule add` with a local source.
- Acceptance criterion 3: the runner regression observes exactly three
  calls for two nested candidates (two identity commands and one listing),
  asserts every listing argument and the root directory, and inventories
  only `packages/api/AGENTS.md`. A separate case has both root carriers and
  an agent guide but only the two identity calls. The error case asserts
  `inspect bounded repository state: walk bounded carriers: list repository instruction carriers: listing failed`
  verbatim and verifies that `errors.Is` still matches the original error.

### Focused checks

1. `go test -count=1 -run '^TestNestedCarrierInventory' -v ./internal/baseline`
   — before the production edit, exit 1: foreign carriers changed warnings
   and digest, the listing call was absent, an omitted carrier was retained,
   and a listing failure was not surfaced. After the production edit, exit 0:
   both new tests and all their cases passed.
2. `go test -count=1 -run '^(TestNestedCarrierInventory.*|TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths|TestRepositoryInspectionUsesNarrowReadOnlyGitCommands|TestRepositoryInspectionNoMutation|TestRepositoryInspectionHonorsCancellation|TestInventoryBudget|TestPlanDeterminismAndNoMutation|TestApplyExactDigest|TestApplyStalePreimage|TestApplyPostimageFailureRollsBack|TestTransactionStagesBeforeMutation|TestTransactionRollback)$' -v ./internal/baseline`
   — exit 0; 13 selected top-level tests passed, including both inventory
   budget cases and the apply/transaction consumers. The suite guard
   reported no repository mutation or helper residue.
3. `go test -count=1 -run '^(TestBaselineApplyCommand|TestBaselineUpdateAppliesManifestPlanAndReportsJSON)$' -v ./internal/cli`
   — exit 0; both focused CLI consumer tests passed. The suite guard
   reported no repository mutation or helper residue.
4. A read-only comparison with `HEAD` confirmed exactly the two required
   new tests, only the over-budget snapshot-call change in existing tests,
   byte-identical `inventoryPathIgnored` and `isBoundedCarrierPath`, and the
   extracted carrier body unchanged apart from its void return.
   `git diff --check` exited 0. Changed-file scope is the five declared Go
   interfaces/test files plus this Task file; the starting Daemon status
   change is preserved.

No declared Verification command, selected repository Verification, full
suite or Windows build was run; those remain with the Daemon. No Task
status, Task Graph or other Task file was edited, and no commit, push,
Pull Request, ai-memory message or Spec event was created. No follow-up
outside this slice was identified.
