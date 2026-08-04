---
task: task_03
spec: 0068-spec-close-audit
status: pending
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
