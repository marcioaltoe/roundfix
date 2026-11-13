---
task: task_01
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
status: pending
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
