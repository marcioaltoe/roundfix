---
task: task_03
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: completed
type: backend
complexity: medium
---

# Task 03: A related-ADR gap opens only for ADRs that predate the Spec

## Overview

`SC-ADR-RELATED`, in `detectADRConsistency` in `internal/speccheck/citations.go`, reports every accepted ADR that cites an ADR a Spec lists. A new ADR therefore opens a gap in every other active Spec, and wherever gaps are promoted that gap refuses: the strict check, the Delivery Plan, queue revalidation and the QA gate. This Task gives the detector a horizon from commit ancestry (ADR-0168). For a Spec whose `_prd.md` is committed, an ADR is related only when the commit that added it is an ancestor of the commit that added the PRD. Every caller of the Spec Consistency Check inherits the horizon.

## Requirements

1. MUST add `internal/speccheck/adr_horizon.go` with `adrHorizon`, `newADRHorizon` and `predates`, as the TechSpec's "The related-ADR horizon" section states. All Git runs with `mechanicalGitEnvironment()`. `newADRHorizon` returns false, which keeps the full check, in each of these cases:
   - the PRD and the repository resolve to different Git common directories;
   - `git rev-parse --is-shallow-repository` prints `true`, because a shallow boundary commit appears to add every file it holds;
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
- [ ] An uncommitted PRD, a Spec directory without Git, and a shallow clone keep the full check.
- [ ] An uncommitted ADR against a committed PRD opens no finding.
- [ ] A Spec that cites an ADR outside its horizon still reports `SC-ADR-UNLISTED`.

## Context

- instruction: `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- interface: `internal/speccheck/citations.go`
- creates: `internal/speccheck/adr_horizon.go`
- creates: `internal/speccheck/adr_horizon_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnADRCommittedAfterTheSpecOpensNoRelatedGap|TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap|TestAnADRCommittedWithTheSpecOpensTheRelatedGap|TestAnUncommittedSpecKeepsTheFullRelatedCheck|TestAnUncommittedADRIsOutsideACommittedSpecsHorizon|TestASpecWithoutGitKeepsTheFullRelatedCheck|TestAShallowHistoryKeepsTheFullRelatedCheck|TestTheHorizonLeavesUnlistedCitationsChecked|TestCheckADRClosureDepthOne)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnADRCommittedAfterTheSpecOpensNoRelatedGap TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap TestAnADRCommittedWithTheSpecOpensTheRelatedGap TestAnUncommittedSpecKeepsTheFullRelatedCheck TestAnUncommittedADRIsOutsideACommittedSpecsHorizon TestASpecWithoutGitKeepsTheFullRelatedCheck TestAShallowHistoryKeepsTheFullRelatedCheck TestTheHorizonLeavesUnlistedCitationsChecked TestCheckADRClosureDepthOne; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The related-ADR horizon
- `_prd.md` → Goal 3; Goal 4; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 5; Testing Approach 3
- [references/2026-09-29-a-new-adr-forces-edits-to-every-active-spec.md](references/2026-09-29-a-new-adr-forces-edits-to-every-active-spec.md)

## Result

Implementation:

- Added a fail-closed related-ADR horizon that reads the PRD's adding commit,
  reads every ADR adding commit in one log, keeps the newest add per path, and
  answers ancestry with `merge-base --is-ancestor`. Every Git subprocess uses
  the mechanical Git environment.
- Applied the horizon lazily on the first related candidate only. Unlisted and
  unsupported citation checks, skip behavior, and finding text are unchanged.
- Added real-Git and plain-directory coverage for later, earlier, same-commit,
  re-added, uncommitted, different-repository, shallow-history, and unlisted
  citation behavior.

Focused checks:

- Before implementation,
  `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestAnADRCommittedAfterTheSpecOpensNoRelatedGap$' ./internal/speccheck`
  failed with the existing `SC-ADR-RELATED` finding for ADR-0002.
- After implementation, the same focused test passed.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run 'ADRCommitted|UncommittedSpec|UncommittedADR|SpecWithoutGit|SpecInAnotherRepository|ShallowHistory|HorizonLeaves' ./internal/speccheck`
  passed.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^TestTheHorizonUsesTheNewestAddingCommit$' ./internal/speccheck`
  passed.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run '^(TestCheckADRClosureDepthOne|TestCheckReplay0056F001FromReport)$' ./internal/speccheck`
  passed after the repository-root guard restored the non-Git fixture behavior.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck`
  passed.
- `GOCACHE=/tmp/roundfix-task03-gocache make verify-incremental` passed with
  host process-table access. The sandboxed attempt could not enumerate child
  processes in two force-stop integration tests. One permission-enabled run
  then hit an unrelated 200 ms daemon budget at 211 ms; that test passed alone
  unchanged, and the unchanged incremental command passed on the next run.

Acceptance evidence:

- `TestAnADRCommittedAfterTheSpecOpensNoRelatedGap` covers a later committed
  ADR producing no related gap.
- `TestAnADRCommittedBeforeTheSpecOpensTheRelatedGap` and
  `TestAnADRCommittedWithTheSpecOpensTheRelatedGap` cover earlier and
  same-commit ADRs retaining the gap.
- `TestAnUncommittedSpecKeepsTheFullRelatedCheck`,
  `TestASpecWithoutGitKeepsTheFullRelatedCheck`, and
  `TestAShallowHistoryKeepsTheFullRelatedCheck` cover every declared full-check
  fallback; `TestASpecInAnotherRepositoryKeepsTheFullRelatedCheck` separately
  covers different Git common directories.
- `TestAnUncommittedADRIsOutsideACommittedSpecsHorizon` covers an uncommitted
  ADR producing no related gap.
- `TestTheHorizonLeavesUnlistedCitationsChecked` covers an outside-horizon ADR
  still producing `SC-ADR-UNLISTED`.

The Daemon-owned Verification command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260929T170240Z_b8558d4fb9103028`
- Source commit: `0eb554a39e2496cdf482e09fac03607e1025fd9c`
