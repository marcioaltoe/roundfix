---
task: task_01
spec: 0069-review-run-targets-its-pull-request
status: completed
type: backend
complexity: medium
---

# Task 01: Compare the checkout against its Pull Request

## Overview

`preflight.Run` already resolves both sides. `InspectGit` returns the
checkout's branch and HEAD; `ResolvePullRequest` returns the Pull Request's
head branch, read straight off `raw.HeadRefName`. Nothing compares them, so a
Review Run acts on whichever branch the checkout happens to be on and the Pull
Request number decides only which Review Issues it fetches.

This slice adds the comparison, in the one place whose contract already
promises no side effects on refusal.

## Requirements

1. MUST compare the checkout's branch against the Pull Request's head branch in
   `preflight.Run`, after the Pull Request resolves and before push planning.
2. MUST refuse a mismatch with exit `2`, naming the Pull Request, both branches,
   and both revisions.
3. MUST include the git command that resolves the mismatch, so recovery is one
   step.
4. MUST create no Run, query no Review Source, start no Agent Session, and make
   no commit or push on refusal, asserted rather than assumed.
5. MUST compare an explicitly supplied `--head-branch` or `--head-repository`
   like any other resolution rather than treating the flag as intent.
6. MUST refuse a Pull Request whose head is on a fork with a message naming
   that cause, since forks are out of scope and a bare branch-name comparison
   would be wrong there.
7. MUST leave a Run whose checkout already matches behaving exactly as it does
   today, asserted over the existing tests unchanged.
8. MUST NOT check out, fetch, or move the working tree on the user's behalf.

## Subtasks

- [ ] Add the mismatch type carrying both sides and its next action.
- [ ] Place the comparison in `preflight.Run` and map it to exit `2`.
- [ ] Table-test the four resolution shapes plus the fork case.
- [ ] Assert no side effects on refusal.

## Acceptance Criteria

- [ ] A checkout on the Pull Request's head branch runs, unchanged.
- [ ] A checkout on another branch refuses with exit `2`, naming the Pull
      Request, both branches, and both revisions.
- [ ] The refusal names a git command that resolves it.
- [ ] A refusal creates no Run row, makes no Review Source call, starts no
      Agent Session, and writes nothing.
- [ ] An explicit `--head-branch` that disagrees with the checkout still
      refuses.
- [ ] A fork head refuses naming the fork, not a branch mismatch.
- [ ] `fetch`, `resolve`, and `watch` all refuse identically.

## Context

- interface: `internal/preflight/preflight.go`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/preflight -count=1 -run 'Target|Mismatch|Branch' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the comparison tests ran and passed.
- `go test ./internal/preflight ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the Skill is task_03's bounded scope.

## References

- `_prd.md` → Core Features 1, 2 and 5; Success Metric 1.
- `_techspec.md` → Interfaces; Build Order 1.
- ADR-0052.

## Result

Implemented the Preflight target guard after Open Pull Request resolution and
before push planning. The resolved Pull Request now carries its `headRefOid`;
`TargetMismatch` carries the Pull Request number, both branches, both
revisions, and a shell-safe `git switch` next action. Preflight refuses fork
heads before branch comparison, and the CLI renders typed next actions through
the existing exit-2 refusal path.

Focused checks:

- Red signal: `GOCACHE=<repo>/.gocache rtk go test ./internal/preflight -count=1 -run TestRunValidatesPullRequestTargetResolutionShapes` failed to build before implementation because `TargetMismatch` and `PullRequest.HeadSHA` did not exist.
- `GOCACHE=<repo>/.gocache rtk go test ./internal/preflight -count=1` — exit 0; 50 tests passed.
- `GOCACHE=<repo>/.gocache rtk go test ./internal/cli -count=1 -run 'TestReviewCommandsRefuseTargetMismatchWithoutSideEffects|TestRunOperationalCommandRejectsInvalidInput|TestRunOperationalCommandAcceptsMVPFlags'` — exit 0; 18 tests passed.

Acceptance evidence:

- A matching checkout remains accepted: the resolved and explicit matching
  rows in `TestRunValidatesPullRequestTargetResolutionShapes` passed, and the
  unchanged operational command acceptance cases remained green.
- A differing checkout returns `TargetMismatch` with Open Pull Request `#123`,
  expected and found branches, and expected and found revisions; the real CLI
  cases asserted exit `2` and the same diagnostic for all three commands.
- The CLI refusal includes `git switch -- '<PR Head Branch>'` under `Next
  action:`; the branch is shell-quoted and no checkout is performed.
- Refusal side effects are asserted in
  `TestReviewCommandsRefuseTargetMismatchWithoutSideEffects`: no Run Database,
  zero Review Source fetch/resolve calls, zero Agent run/probe calls, zero
  verification/commit/push calls, and unchanged Git HEAD and porcelain status.
- The explicit-resolution mismatch row and the real CLI cases supply
  `--head-branch` and `--head-repo`; all refused rather than treating the flags
  as user intent to bypass comparison.
- `TestRunRefusesForkHeadBeforeBranchMismatch` passed with distinct Base and
  Head Repositories, asserted the fork and both repositories in the message,
  and asserted that the error is not `TargetMismatch`.
- `TestRunRefusesTargetMismatchForEveryReviewCommand` and the real CLI table
  both passed for `fetch`, `resolve`, and `watch` with identical refusal
  behavior.

Not run: the commands under this Task's `## Verification`; the Daemon owns
those checks and terminal settlement.

Follow-up discovered: operational command startup already creates and
write-checks the configured Artifact Directory before `preflight.Run`. This
slice leaves that configuration-validation contract unchanged; if “writes
nothing” is intended to prohibit even that existing startup check, it requires
a separate config/CLI Task outside Task 01's enumerated Run, Review Source,
Agent Session, commit, and push refusal boundary.
