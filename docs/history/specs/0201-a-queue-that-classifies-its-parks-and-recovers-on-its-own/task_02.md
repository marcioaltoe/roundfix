---
task: task_02
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: completed
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

## Result

Implemented the Task 02 prerequisite slice for Daemon Verification. Task
status remains Daemon-owned; no authored Verification command was run and no
commit, push or Pull Request was made. The starting worktree contained only
the Daemon's pre-existing edit to this Task file. Baseline `HEAD` inspection
confirmed that `Graph.Requires`, `PrerequisiteReader` and
`UnmetPrerequisites` were absent before this slice.

The manifest parser reads and validates `requires`; delivery start loads all
queued Specs before resolving prerequisite names and detecting deterministic
cycles, then checks authorization. The optional reader fetches the configured
delivery remote (or `origin`) and reads the manifest and archive evidence at
its refreshed remote-tracking default branch. Archive presence uses one
immutable-tree listing for every prerequisite, as ADR-0090 requires; Git
errors remain errors rather than evidence of absence.

The engine waits or parks before creating a worktree, releases dependency
parks at pass start and after merges without counting a retry, and accepts an
operator retry to `queued` without workspace recovery, revalidation or
carry-forward. Ordinary retry counting and limits remain in the store's
retry transition. Dependency Park Classes provide the TechSpec's operator
commands. The delivery guide documents declaration, refusal, waiting,
parking, release and retry.

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Valid and malformed manifests | `TestTaskGraphReadsRequiredSpecs` covers a valid list, an empty list and an omitted declaration; `TestTaskGraphRefusesAMalformedRequiresList` covers scalar, map, non-string entries, empty/whitespace entries, duplicate/trimmed duplicate entries and self-reference, with `_tasks.md` in every refusal. |
| Cycle and unknown prerequisite refusal | `TestDeliverStartRefusesAPrerequisiteCycle` and `TestDeliverStartRefusesAnUnknownPrerequisite` assert exit 2, empty stdout, the complete preflight stderr, no owner launch and no persisted queue. Fixtures lack delivery grants, proving prerequisite refusal precedes authorization. Additional tests cover longer cycles and active/archive resolution for built-in and configured roots. |
| Archive only on refreshed remote default | `TestAPrerequisiteIsMetByItsArchiveOnTheRefreshedDefaultBranch` advances a disposable local remote after cloning, proves the tracking ref was stale, and checks both `origin` and a configured delivery remote. The refreshed archive counts as met; absent and checkout-only archives remain unmet. A checkout-only manifest amendment is ignored. Git fetch, manifest and archive-read errors are separately covered. |
| Wait, park without worktree, release with unchanged retry count | `TestAnItemWaitsForAPrerequisiteAheadInTheQueue` checks the wait log, queued state and absence of dependent actions while the prerequisite proceeds. `TestAnItemParksWhenItsPrerequisiteIsParked` covers parked, missing and merged-but-unarchived prerequisites. `TestTheOwnerReleasesADependencyParkWhenThePrerequisiteMerges` covers pass-start and after-merge release, cleared blockers, release logs and three existing retries remaining unchanged; after-merge release starts on the next pass. `TestDependencyParkRemainsUntilEveryPrerequisiteIsMet` covers partial satisfaction and a nil reader. |
| No prerequisites preserves normal start | `TestAnItemWithoutPrerequisitesStartsAsBefore` exercises nil and empty readers through merge. `TestPrerequisiteReaderWithoutRequiresSkipsArchiveRead` checks that an omitted declaration needs no archive lookup. Existing delivery regression tests also pass. |

Additional checks: `TestRetryReturnsADependencyParkToTheQueue` uses nil
workspace, recovery and revalidator dependencies to prove the retry bypass;
`TestPrerequisiteParkGuidesTheOperator` asserts both exact next commands.
All new tests use fakes or disposable filesystem remotes; none reaches a
network remote.

Focused checks run in this turn:

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec ./internal/delivery ./internal/store ./internal/cli -run 'TaskGraph|Prerequisite|DependencyPark|WithoutPrerequisites|Deliver|Delivery|WithoutARevalidator' -count=1` — exit 0; all four packages passed after the final code/test edits.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test ./internal/delivery -count=1` — exit 0 after restoring the existing missing-revalidator retry diagnostic.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test ./internal/delivery ./internal/store -count=1` — the store passed; the initial delivery run exposed that diagnostic regression, subsequently repaired and covered by the final focused run above.
- Initial focused delivery check using the default Go cache could not access a sandbox-restricted cache file. Subsequent checks used the writable task-scoped cache above.
- The initial release fixture tried to seed `RetryCount` through an ordinary update, which deliberately does not write it. The fixture now seeds three real retry transitions and verifies release preserves that persisted count.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.
- Python inspection of the delivery guide confirmed the documented `requires`, wait, `prerequisite-unmerged`, release and retry behavior.

Not run: the Task's declared Verification commands, the repository gate and
final Spec QA; the Daemon owns authored Verification and settlement.
The Task Graph and other Task files are unchanged. No follow-up feature work
was added to this diff; remaining Spec slices retain their existing owners.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/spec/requires.go`
- `internal/spec/spec_test.go`
