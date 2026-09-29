---
task: task_03
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: backend
complexity: medium
---

# Task 03: A related-ADR gap opens only for ADRs that predate the Spec

## Overview

`SC-ADR-RELATED`, in `detectADRConsistency` in `internal/speccheck/citations.go`, reports every accepted ADR that cites an ADR a Spec lists. A new ADR therefore opens a gap in every other active Spec, and wherever gaps are promoted that gap refuses: the strict check, the Delivery Plan, queue revalidation and the QA gate. This Task gives the detector a horizon from commit ancestry (ADR-0168). For a Spec whose `_prd.md` is committed, an ADR is related only when the commit that added it is an ancestor of the commit that added the PRD. Every caller of the Spec Consistency Check inherits the horizon.

## Requirements

1. MUST add `internal/speccheck/adr_horizon.go` with `adrHorizon`, `newADRHorizon` and `predates`, as the TechSpec's "The related-ADR horizon" section states. All Git runs with `mechanicalGitEnvironment()`. `newADRHorizon` returns false, which keeps the full check, in each of these cases:
   - the PRD and the repository resolve to different Git common directories;
   - the PRD has no adding commit;
   - any Git command fails.
2. MUST read every ADR's adding commit with one `git log --diff-filter=A` over `docs/adr`, keeping the newest adding commit per path. `predates` MUST return false for an ADR with no adding commit, the ancestry answer when `git merge-base --is-ancestor` exits `0` or `1`, and true on any other failure.
3. MUST make `detectADRConsistency` build the horizon at most once, only when a related candidate exists, and skip the candidate when `predates` reports false. `SC-ADR-UNLISTED`, `SC-CITATION-UNSUPPORTED`, the skip rules and every finding's text MUST stay unchanged.
4. MUST leave existing characterizations and the corpus golden unmoved. Fixtures without Git keep the full check, and `TestCheckADRClosureDepthOne` stays green without editing it. The Task renames or removes no top-level test and changes no exported function signature.
5. MUST put the new tests in `internal/speccheck/adr_horizon_test.go`, each over a `gittest` repository or a plain temporary directory, never over this repository's own history.

## Subtasks

- [ ] Add the horizon from commit ancestry.
- [ ] Apply it to related candidates only.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An ADR committed after a committed PRD that lists an ADR it cites opens no `SC-ADR-RELATED` finding.
- [ ] An ADR committed before the PRD, or in the same commit, still opens the finding.
- [ ] An uncommitted PRD, and a Spec directory without Git, keep the full check.
- [ ] An uncommitted ADR against a committed PRD opens no finding.
- [ ] A Spec that cites an ADR outside its horizon still reports `SC-ADR-UNLISTED`.

## Context

- instruction: `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- interface: `internal/speccheck/citations.go`
- creates: `internal/speccheck/adr_horizon.go`
- creates: `internal/speccheck/adr_horizon_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnADRCommittedAfterTheSpecOpensNoRelatedGap|TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap|TestAnADRCommittedWithTheSpecOpensTheRelatedGap|TestAnUncommittedSpecKeepsTheFullRelatedCheck|TestAnUncommittedADRIsOutsideACommittedSpecsHorizon|TestASpecWithoutGitKeepsTheFullRelatedCheck|TestTheHorizonLeavesUnlistedCitationsChecked|TestCheckADRClosureDepthOne)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnADRCommittedAfterTheSpecOpensNoRelatedGap TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap TestAnADRCommittedWithTheSpecOpensTheRelatedGap TestAnUncommittedSpecKeepsTheFullRelatedCheck TestAnUncommittedADRIsOutsideACommittedSpecsHorizon TestASpecWithoutGitKeepsTheFullRelatedCheck TestTheHorizonLeavesUnlistedCitationsChecked TestCheckADRClosureDepthOne; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The related-ADR horizon
- `_prd.md` → Goal 3; Goal 4; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 5; Testing Approach 3
- [references/2026-09-29-a-new-adr-forces-edits-to-every-active-spec.md](references/2026-09-29-a-new-adr-forces-edits-to-every-active-spec.md)

## Result
