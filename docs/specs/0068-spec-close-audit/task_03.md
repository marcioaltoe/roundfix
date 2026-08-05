---
task: task_03
spec: 0068-spec-close-audit
status: completed
type: backend
complexity: medium
---

# Task 03: Report what the Spec claims but has not delivered

## Overview

The expensive residue is not a stale branch — it is work reported as delivered
that lives only on an unmerged branch. Spec 0058's archive and five queued
Specs were in exactly that state: the default branch still showed 0058 active
and showed none of the new Specs.

This slice checks every artifact a Spec claims against the synced default
branch and names the branch holding anything absent.

## Requirements

1. MUST verify each artifact the Spec claims — its folder under the active or
   archived Spec Root, and files its Tasks recorded — is present on the synced
   default branch.
2. MUST name the branch that holds any absent artifact, so the reader learns
   where the work is rather than only that it is missing.
3. MUST report an artifact it cannot locate on any branch as absent without
   naming a holder, rather than guessing one.
4. MUST read the default branch's tree rather than the working copy, so an
   uncommitted file never reads as delivered.
5. MUST NOT fetch, pull, or otherwise mutate Git state; the caller syncs.

## Subtasks

- [ ] Resolve the Spec's claimed artifacts.
- [ ] Compare against the default branch's tree.
- [ ] Resolve the holding branch for each absent artifact.
- [ ] Add fixtures: delivered, held-by-branch, and nowhere.

## Acceptance Criteria

- [ ] A Spec whose archive exists only on an unmerged branch is reported as not
      delivered, naming that branch.
- [ ] A Spec whose artifacts are all on the default branch reports nothing.
- [ ] An artifact present in the working copy but not committed is reported as
      absent.
- [ ] An artifact on no branch is reported absent with no holder named.
- [ ] The check performs no fetch, pull, or write, proven by a fixture whose
      Git state is byte-identical before and after.

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/specaudit -count=1 -run 'Deliver|Undelivered|Held' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the delivery tests ran and passed.
- `go test ./internal/specaudit -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 2; Goals; Success Metric 3.
- `_techspec.md` → Data Models; Build Order 3.

## Result

Implemented delivery verification in `internal/specaudit`. The audit resolves
the configured active or archived Spec Root, including an external Spec Root
owned by another Git repository, and treats that repository's committed
default tree as the delivery source of truth. It derives implementation claims
from commits carrying matching `Roundfix-Spec` and `Roundfix-Task` trailers,
maps archived Task artifacts to their archived paths, and reports each absent
artifact with the first deterministic local or remote-tracking branch whose
tip contains it. An artifact on no branch keeps an empty `HeldBy` value.

Focused checks:

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/specaudit -count=1 -run
  '^TestAuditReportsUndeliveredArchiveHeldByBranch$' -v` failed after fixture
  setup was corrected because `Undelivered` was empty instead of naming
  `ma/spec-close-archive`.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/specaudit -count=1 -run
  '^TestAudit(ClassifiesPullRequestBranch|ClassifiesPendingBranch|ClassifiesResidueBranch|PreservesUnmatchedWorktree|PreservesActiveRunSurvivors|ReportsUndeliveredArchiveHeldByBranch|ReportsNothingWhenClaimedArtifactsAreDelivered|ReportsUncommittedWorkingCopyArtifactAsUndelivered|ReportsUndeliveredArtifactWithNoHoldingBranch|DeliveryCheckLeavesGitStateByteIdentical|UsesConfiguredExternalSpecRootTree)$'
  -v` — exit 0; all eleven selected survivor and delivery fixtures passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go vet
  ./internal/specaudit` — exit 0.
- `rtk gofmt -w internal/specaudit/audit.go
  internal/specaudit/audit_test.go` — exit 0 before the final focused run.
- `rtk git -c core.fsmonitor=false diff --check` — exit 0 before this Result
  update.

Acceptance evidence:

1. `TestAuditReportsUndeliveredArchiveHeldByBranch` leaves the active folder on
   `main`, archives it only on `ma/spec-close-archive`, and asserts the archived
   path is undelivered with that holding branch.
2. `TestAuditReportsNothingWhenClaimedArtifactsAreDelivered` commits both the
   Spec folder and a Task-recorded implementation file on `main` and asserts an
   empty delivery report.
3. `TestAuditReportsUncommittedWorkingCopyArtifactAsUndelivered` creates the
   Spec only in the working copy and asserts it remains undelivered with no
   holder, proving the working tree is not delivery evidence.
4. `TestAuditReportsUndeliveredArtifactWithNoHoldingBranch` records a file in a
   Task commit, removes it from every branch tip, and asserts an empty holder
   rather than a guessed branch.
5. `TestAuditDeliveryCheckLeavesGitStateByteIdentical` snapshots every byte
   under `.git` before and after the audit and compares the snapshots equal.
   The production Git runner continues to use `--no-optional-locks`; the new
   delivery path contains no fetch, pull, or write operation.

`TestAuditUsesConfiguredExternalSpecRootTree` additionally proves that Spec
folder delivery follows `specs.root` and the external repository's default
tree instead of hard-coding `docs/specs` in the code repository.

The commands under `## Verification` were not run; the Daemon owns those
commands and Task settlement.
