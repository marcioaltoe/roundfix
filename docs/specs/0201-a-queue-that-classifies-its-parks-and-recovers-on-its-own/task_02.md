---
task: task_02
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: pending
type: backend
complexity: high
---

# Task 02: A Spec names its prerequisites and the owner waits for them

## Overview

Queue order is the only order the Delivery Queue knows, and a parked item never holds back the items behind it. On 2026-09-30 that started Spec 0195 before the two Specs its Verification needed. This Task lets a Task Graph manifest name prerequisite Specs in `requires`, has `deliver start` refuse a queue that could never finish, and has the owner wait for, park on and release unmet prerequisites (ADR-0193).

## Requirements

1. MUST read an optional `requires` list from the Task Graph manifest frontmatter into `spec.Graph.Requires`. A value that is not a list of strings, an empty or duplicate entry, or the Spec's own slug MUST make `spec.Load` fail with an error naming `_tasks.md`. A manifest without `requires` MUST load exactly as today.
2. MUST make `roundfix deliver start` refuse, with exit `2`, after every Spec loads and before the delivery authorization check and the queue record, a `requires` entry that names no Spec in the Specs Root or its archive root, and a cycle among the queued Specs, with the reason `Delivery Queue Specs require each other in a cycle: <a> -> <b> -> <a>` of Surface Transcript 2.
3. MUST add the optional `PrerequisiteReader` dependency and implement `UnmetPrerequisites` in `internal/cli/deliver_workflow.go`: fetch the default branch from the delivery remote resolved as `CreateItemBranch` resolves it, read the item's `_tasks.md` at the refreshed remote-tracking ref, and return each prerequisite whose archived `_prd.md` is absent at that ref.
4. MUST make the engine consult it for a `queued` item before `CreateItemBranch`: start as today with none unmet; stay `queued`, log the wait and move on when every unmet prerequisite is a queue item neither parked nor merged; otherwise park `prerequisite-unmerged: <slug>, …` without creating a worktree.
5. MUST return each `prerequisite-unmerged` item whose prerequisites are all met to `queued`, with its blocker cleared and its retry count unchanged, at the start of every engine pass and after any item merges, and log the release.
6. MUST make `Retry` of a `prerequisite-unmerged` item skip the worktree, the revalidation and the carry-forward and target `queued`, and MUST let the store's retry transition accept `queued` as a target.
7. MUST classify `prerequisite-unmerged` as `dependency`, with the next command the TechSpec's Park Classes table states for a prerequisite parked in the queue and for one that is not.
8. MUST describe `requires`, the wait, the `prerequisite-unmerged` park and its release in `docs/user-guide/commands/deliver.md`.
9. MUST keep a nil `PrerequisiteReader` and a Spec without `requires` on today's behavior, and MUST NOT reach a network remote in any test.

## Subtasks

- [ ] Read and validate `requires` in the Task Graph manifest.
- [ ] Refuse an unknown prerequisite and a cycle at `deliver start`.
- [ ] Read unmet prerequisites at the refreshed default branch.
- [ ] Wait, park and release in the engine, and accept the retry.
- [ ] Describe the behavior in the delivery command guide, with a test for each criterion.

## Acceptance Criteria

- [ ] A manifest with a valid `requires` list loads it; each malformed form is refused with `_tasks.md` in the error.
- [ ] `deliver start` of two Specs that require each other exits `2`, prints Surface Transcript 2's reason and records no queue; an unknown prerequisite is refused the same way.
- [ ] A prerequisite archived only on the remote default branch, not yet fetched locally, counts as met; one archived nowhere does not.
- [ ] An item whose prerequisite is ahead in the queue waits `queued`; one whose prerequisite is parked parks `prerequisite-unmerged` with no worktree; after the prerequisite merges it returns to `queued` with the same retry count.
- [ ] A Spec without `requires` starts exactly as before.

## Context

- interface: `internal/spec/spec.go`
- creates: `internal/spec/requires_test.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/park_class.go`
- creates: `internal/delivery/prerequisite_test.go`
- interface: `internal/store/delivery.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_prerequisite_test.go`
- creates: `docs/user-guide/commands/deliver.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskGraphReadsRequiredSpecs|TestTaskGraphRefusesAMalformedRequiresList|TestAnItemWaitsForAPrerequisiteAheadInTheQueue|TestAnItemParksWhenItsPrerequisiteIsParked|TestTheOwnerReleasesADependencyParkWhenThePrerequisiteMerges|TestRetryReturnsADependencyParkToTheQueue|TestAnItemWithoutPrerequisitesStartsAsBefore)$" ./internal/spec ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTaskGraphReadsRequiredSpecs TestTaskGraphRefusesAMalformedRequiresList TestAnItemWaitsForAPrerequisiteAheadInTheQueue TestAnItemParksWhenItsPrerequisiteIsParked TestTheOwnerReleasesADependencyParkWhenThePrerequisiteMerges TestRetryReturnsADependencyParkToTheQueue TestAnItemWithoutPrerequisitesStartsAsBefore; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDeliverStartRefusesAPrerequisiteCycle|TestDeliverStartRefusesAnUnknownPrerequisite|TestAPrerequisiteIsMetByItsArchiveOnTheRefreshedDefaultBranch)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverStartRefusesAPrerequisiteCycle TestDeliverStartRefusesAnUnknownPrerequisite TestAPrerequisiteIsMetByItsArchiveOnTheRefreshedDefaultBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands/deliver.md|requires" "docs/user-guide/commands/deliver.md|prerequisite-unmerged"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task `prerequisite-unmerged` is not in the delivery command guide, so the command fails.

## References

- `_prd.md` → Goal 1; User Story 1; Core Feature 1; Success Metric 1
- `_techspec.md` → Prerequisites; Park Classes; API Contracts 2-3; Surface Transcript 2; Testing Approach 2; Build Order 2
- ADR-0193; ADR-0090
