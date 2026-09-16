---
task: task_03
spec: 0139-a-suite-that-passes-where-it-runs
status: completed
type: test
complexity: medium
---

# Task 03: Prove the historical audit with a controlled fixture

## Overview

The historical subtest of `TestAuditJudgesTheGrant` reads three commits that
are no longer reachable from the repository's history. Its expectation has been
stale since the audit started reading a grant at its authorizing ancestor, which
Spec 0119 requires. Its skip guard checks only whether a commit's object exists,
so the subtest runs where old objects survive locally and skips everywhere
else.

This slice replaces the subtest with a repository the test builds itself, and
makes the guard require reachability.

## Requirements

1. MUST remove the historical subtest and the constants naming the unreachable
   commits: `archiveHelpCommit`, `archiveCarrierCommit` and
   `verificationTaskCommit`.
2. MUST add `TestAuditRefusesAGrantWidenedAfterItsConsumingCommit`. The test
   builds its own repository in which a Task commit changes a governed path, then
   records two cases:
   - when the authorization record bounds that path only in a later commit, the
     audit returns exactly one `QA-AUTH-PATHS` finding for that path;
   - when the record already bounded the path at the authorizing ancestor, the
     audit returns no finding.
3. MUST make `firstMissingMechanicalCommit` report a commit as missing unless it
   is an ancestor of the audited repository's `HEAD`, and MUST keep its
   signature.
4. MUST keep `TestMechanicalAuditJudgesTheProjectRoot` and every other remaining
   `TestAuditJudgesTheGrant` subtest passing without editing them.
5. MUST NOT change production code or weaken any other assertion.

## Subtasks

- [ ] Remove the historical subtest and its commit constants.
- [ ] Add the controlled late-widening fixture test.
- [ ] Change the guard to reachability.

## Acceptance Criteria

- [ ] The test file no longer names the three unreachable commits; today it does.
- [ ] `TestAuditRefusesAGrantWidenedAfterItsConsumingCommit` passes and records
      both the refused and the authorized case.
- [ ] `firstMissingMechanicalCommit` uses ancestry, not object existence.
- [ ] `TestAuditJudgesTheGrant` and `TestMechanicalAuditJudgesTheProjectRoot`
      pass.

## Context

- interface: `internal/speccheck/mechanical_test.go`
- interface: `internal/speccheck/constraints_characterization_test.go`

## Verification

- `test -f internal/speccheck/mechanical_test.go && ! grep -q -E 'archiveCarrierCommit|archiveHelpCommit|verificationTaskCommit' internal/speccheck/mechanical_test.go` — the unreachable historical commits are gone; this fails today.
- `awk 'index($0, "func firstMissingMechanicalCommit(") == 1 {inside=1} inside {print} inside && $0 == "}" {inside=0}' internal/speccheck/mechanical_test.go | grep -q 'is-ancestor'` — the guard requires reachability; this fails today.
- `grep -q 'func TestAuditRefusesAGrantWidenedAfterItsConsumingCommit' internal/speccheck/mechanical_test.go && go test -count=1 ./internal/speccheck -run '^(TestAuditRefusesAGrantWidenedAfterItsConsumingCommit|TestAuditJudgesTheGrant|TestMechanicalAuditJudgesTheProjectRoot)$'` — the fixture test and the remaining audit tests pass.

## References

- `_prd.md` → Goals 1 and 4; User Story 4; Core Feature 3; Declared intentional
  breaks.
- `_techspec.md` → Implementation Design: Historical audit fixture and guard;
  Testing Approach 3; Build Order 3.

## Result

Implemented the controlled authorization history in
`internal/speccheck/mechanical_test.go`. The new test records a Task commit that
changes `.golangci.yml` and proves both outcomes: a later grant widening yields
one `QA-AUTH-PATHS` finding for that path, while a grant already bounded at the
authorizing ancestor yields none. The stale repository-history subtest and its
three unreachable commit constants are gone. `firstMissingMechanicalCommit`
now treats a commit as present only when `git merge-base --is-ancestor` proves
it is reachable from the audited repository's `HEAD`.

Pre-change signal: the focused historical subtest exited 1 with three
`QA-AUTH-PATHS` findings from the stale expectation, while ancestry checks
returned exit 1 for each of its three pinned commits.

Focused evidence for the acceptance criteria:

- Removed historical names: `rtk proxy rg -n
  'archiveCarrierCommit|archiveHelpCommit|verificationTaskCommit'
  internal/speccheck/mechanical_test.go` found no matches (expected exit 1).
- Controlled refused and authorized cases: `rtk env
  GOCACHE=/private/tmp/roundfix-task03-go-cache go test
  ./internal/speccheck -run
  '^TestAuditRefusesAGrantWidenedAfterItsConsumingCommit$' -v` exited 0; both
  named subtests passed.
- Reachability guard: `rtk proxy rg -n
  'merge-base.*--is-ancestor' internal/speccheck/mechanical_test.go` matched the
  `firstMissingMechanicalCommit` command.
- Remaining grant audit cases: `rtk env
  GOCACHE=/private/tmp/roundfix-task03-go-cache go test
  ./internal/speccheck -run '^TestAuditJudgesTheGrant$' -v` exited 0; all five
  remaining subtests passed.
- Project-root audit: `rtk env
  GOCACHE=/private/tmp/roundfix-task03-go-cache go test
  ./internal/speccheck -run '^TestMechanicalAuditJudgesTheProjectRoot$' -v`
  exited 0.
- Diff hygiene: `rtk git diff --check` exited 0.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
