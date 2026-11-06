---
task: task_02
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: completed
type: backend
complexity: high
---

# Task 02: A repository profile that lacks a template decision is refused with the decision, the line and the flag

## Overview

A repository-owned Baseline Profile written before the catalog's core
template required `verification.incremental` passes `baseline profile
validate`, and `baseline plan` and `baseline update` then fail with
`render managed entry "guide.agent-instructions": missing render value for
token "verification.incremental"`. This Task adds one helper that finds
every decision a template the profile renders requires and the profile does
not reach, and makes the three commands refuse with the decision, the
template, the quoted entry to add and the `--decision` flag, writing nothing.
It answers the message reported by sql-library
([references](references/2026-11-06-a-repository-profile-without-a-template-decision-fails-opaquely.md)).

## Requirements

1. MUST add `internal/baseline/profile_decision_gap.go` with
   `ProfileDecisionGap`, `ProfileDecisionGaps`, `ProfileDecisionGapMessage`
   and `ProfileDecisionGapRepair` exactly as `_techspec.md` → Interfaces and
   API Contracts 2 and 3 define them, reusing `decisionOptional`,
   `renderDecisionValue`, `objectsOrEmpty`, `stringsOrEmpty` and
   `stringValue`. The helper reads only the catalog and the profile and
   writes nothing.
2. MUST call it in `buildPlanWithCatalog` (`internal/baseline/plan.go`)
   for both the profile ID and the profile draft branch, after
   `carryRecordedDecisions`, the `--decision` values and
   `normalizePlanDecisions` and before any template renders, keep only the
   gaps whose decision has no value in the decision set rendering receives,
   and return
   `actionOutcome("decision", ProfileDecisionGapMessage(profile, remaining), ProfileDecisionGapRepair(profile, remaining, "rerun with"), initial.Snapshot.Warnings)`
   when any remain (API Contract 4). A profile draft or repository profile
   whose unlisted decisions all get values from the Setup Manifest, a flag or
   a default must keep planning: `TestProfileAdaptationCannotRemoveUniversalRequiredCapabilities`,
   `TestProfileDraftPlanRejectsSimultaneousAndConflictingInputs`,
   `TestProfileDraftPlanIncludesCanonicalRepositoryProfile`,
   `TestProfileDraftPlanAcceptsMatchingSourceBaselineWithoutTransition`,
   `TestProfileDraftStaleTargetProducesNoMutation`,
   `TestProfileDraftRollbackRestoresMissingProfile` (`./internal/baseline`)
   and `TestBaselineHumanProfileAdaptation`, `TestProfileAdaptationJourney`,
   `TestBaselinePlanProfileFile` (`./internal/cli`) pass unchanged; the
   first Run refused them all by checking before the decisions were known.
   The profile path used is the portable path the branch already computes.
3. MUST make `runBaselineProfileValidateCommand`
   (`internal/cli/baseline_profile.go`) check every resolved profile after
   its `Path` is converted by `baselineProfileOutputPath`, and fail through
   `printBaselineProfileFailure("profile.validate", …)` with exit 2 and the
   message of API Contract 5 when a profile has gaps; the repair's command
   is `run roundfix baseline update`. Every other output keeps its bytes.
4. MUST add `internal/baseline/profile_decision_gap_test.go` (suite header
   comment like its neighbours) with tests that each call `t.Parallel()`
   first:
   - `TestNoBuiltInProfileHasADecisionGap`: every catalog profile resolved
     with `resolveBuiltInProfile` has no gap;
   - `TestAProfileWithoutTheIncrementalDecisionHasOneGap`: each built-in
     profile with `verification.incremental` removed from `Decisions` has
     exactly the gap `verification.incremental`,
     `template.guide.agent-instructions`, `guide.agent-instructions`,
     suggested `make verify-incremental`;
   - `TestProfileDecisionGapTextNamesTheLineAndTheFlag`: the message and the
     repair for one gap, for two gaps in sorted order, and `<value>` for an
     empty suggestion, compared with whole expected strings;
   - `TestBuildPlanRefusesAProfileDecisionGapBeforeRendering`: a temporary
     repository with a repository profile copied from `go-cli-tui` without
     the decision gets a `decision` outcome with that message and no plan.
5. MUST add `internal/cli/baseline_profile_decision_gap_test.go` (suite
   header comment) with tests that each call `t.Parallel()` first and use
   temporary repositories only:
   - `TestBaselineProfileValidateNamesADecisionGap` asserts Surface
     Transcript 1 line by line (stdout, stderr, exit code) and, with
     `--format json`, the message in the failure result;
   - `TestBaselinePlanRefusesAProfileDecisionGap` asserts Surface
     Transcript 2 line by line and that the repository bytes are unchanged;
   - `TestBaselineUpdateRefusesAProfileDecisionGap`: an adopted repository
     whose Setup Manifest names a repository profile without the decision
     and records no value for it; `baseline update --format json` exits 3
     with state `action_required`, category `decision`, the message and the
     `nextAction` of API Contracts 3 and 4, and writes nothing;
   - `TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt`: the same
     repository after the profile lists the decision; the update exits 3 and
     `newDecisions` names `verification.incremental`.
6. MUST keep passing, unchanged, `TestBaselineProfileCommandInitShowAndValidate`,
   `TestBaselineProfileCommandRejectsInvalidAndNonRepositoryProfiles`,
   `TestBaselineUpdateNamesTheMissingIncrementalDecision` and every other
   test of `./internal/baseline` and `./internal/cli` that touches the Baseline
   Command.
7. MUST NOT edit a repository profile from any command, change the Setup
   Manifest schema, the plan or update JSON schema, the help text, any
   module, profile, template or decision asset, or any guide or skill.

## Subtasks

- [ ] Add the gap helper and its unit tests.
- [ ] Refuse in `buildPlanWithCatalog` before decisions are normalized.
- [ ] Refuse in `baseline profile validate`.
- [ ] Add the command tests for validate, plan and update.

## Acceptance Criteria

- [ ] No built-in profile has a gap, and each built-in profile without `verification.incremental` has exactly that gap.
- [ ] `baseline profile validate` exits 2, and `baseline plan` and `baseline update` exit 3, with the decision, the template, the quoted entry and the `--decision` flag, writing nothing.
- [ ] Once the profile lists the decision, `baseline update` exits 3 and names it in `newDecisions`.

## Context

- instruction: `docs/adr/0372-the-baseline-names-the-archive-destination-a-profiles-missing-decision-and-a-guide-it-does-not-own.md`
- instruction: `internal/baseline/custom_profile.go`
- instruction: `internal/baseline/update.go`
- instruction: `internal/cli/baseline_update.go`
- instruction: `internal/cli/baseline_profile_test.go`
- instruction: `internal/cli/baseline_incremental_verification_test.go`
- instruction: `internal/cli/baseline_update_test.go`
- creates: `internal/baseline/profile_decision_gap.go`
- creates: `internal/baseline/profile_decision_gap_test.go`
- creates: `internal/cli/baseline_profile_decision_gap_test.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/cli/baseline_profile.go`

## Consumers

- The `decision` outcome of `BuildPlan` for a profile decision gap — consumers: `./internal/baseline`, `./internal/cli`

## Verification

- `for t in TestNoBuiltInProfileHasADecisionGap TestAProfileWithoutTheIncrementalDecisionHasOneGap TestProfileDecisionGapTextNamesTheLineAndTheFlag TestBuildPlanRefusesAProfileDecisionGapBeforeRendering; do grep -q "func $t(" internal/baseline/profile_decision_gap_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestNoBuiltInProfileHasADecisionGap|TestAProfileWithoutTheIncrementalDecisionHasOneGap|TestProfileDecisionGapTextNamesTheLineAndTheFlag|TestBuildPlanRefusesAProfileDecisionGapBeforeRendering)$' ./internal/baseline && go test -timeout 20m -count=1 ./internal/baseline` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the gap rule, its text and the plan refusal pass with the whole Baseline suite.
- `for t in TestBaselineProfileValidateNamesADecisionGap TestBaselinePlanRefusesAProfileDecisionGap TestBaselineUpdateRefusesAProfileDecisionGap TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt; do grep -q "func $t(" internal/cli/baseline_profile_decision_gap_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; go test -timeout 20m -count=1 -run '^(TestBaselineProfileValidateNamesADecisionGap|TestBaselinePlanRefusesAProfileDecisionGap|TestBaselineUpdateRefusesAProfileDecisionGap|TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt|TestBaselineProfileCommandInitShowAndValidate|TestBaselineProfileCommandRejectsInvalidAndNonRepositoryProfiles|TestBaselineUpdateNamesTheMissingIncrementalDecision)$' ./internal/cli` — expected: exit 0. Before this Task the command test file does not exist, so the first check fails; afterwards the three commands refuse a gap profile with the named text and the existing profile and update tests still pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contracts 2-5; Surface Transcript 1; Surface Transcript 2; Interfaces; Testing Approach; Build Order 2
- ADR-0372


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, systematic-debugging, no-workarounds, testing-boss, tech-writer, evidence-gate.

Implemented the profile decision gap helper and the plan and validate consumers.
The helper computes transitive decision reachability without evaluating `when`,
excludes optional decisions and replaceable templates, deduplicates bindings,
and orders gaps by decision, template and artifact. The message and repair name
the template, the portable profile path, the quoted entries and the flags.

The shared plan branch checks after recorded decisions, flags and normalization,
before rendering or the general missing-decision outcome. It filters normalized
values and the existing artifact render defaults through `artifactRenderValues`;
this preserves backend Profile Adaptations whose omitted authentication and HTTP
tokens already render as empty strings. Validation checks every resolved profile
after converting its path, including discovery and explicit-path input. Update
uses its existing BuildPlan outcome mapping; no update or persisted schema changed.

Acceptance evidence from focused checks:

- Built-in profiles: `TestNoBuiltInProfileHasADecisionGap` and
  `TestAProfileWithoutTheIncrementalDecisionHasOneGap` passed for every embedded
  profile. Additional unit checks cover conditional transitive reachability,
  cycles, optional decisions, duplicate bindings, template alternatives, root
  blocks, ordering, suggestion/default precedence and absent suggestions.
- Refusal and no writes: `TestBuildPlanRefusesAProfileDecisionGapBeforeRendering`,
  `TestBaselineProfileValidateNamesADecisionGap`,
  `TestBaselinePlanRefusesAProfileDecisionGap` and
  `TestBaselineUpdateRefusesAProfileDecisionGap` passed. The CLI checks compare the
  complete transcript streams and exit codes (2 for validate, 3 for plan/update),
  JSON messages/category/nextAction, and repository tree snapshots before and after.
  Transcript 2 passes the absolute temporary repository to `--repo` instead of
  changing the process directory, preserving parallel test isolation; its output
  is byte-identical to the authored transcript.
- Listed decision: `TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt`
  passed with exit 3 and exactly `verification.incremental` in `newDecisions`,
  with unchanged repository bytes.
- Compatibility: all six named baseline Profile Adaptation/Profile Draft checks,
  all three named CLI adaptation checks, both existing profile-command checks,
  and `TestBaselineUpdateNamesTheMissingIncrementalDecision` passed unchanged.

Focused commands (all final runs exited 0):

```sh
rtk proxy go test ./internal/baseline -run '^(TestNoBuiltInProfileHasADecisionGap|TestAProfileWithoutTheIncrementalDecisionHasOneGap|TestProfileDecisionGap.*|TestBuildPlanRefusesAProfileDecisionGapBeforeRendering|TestProfileAdaptationCannotRemoveUniversalRequiredCapabilities|TestProfileDraftPlanRejectsSimultaneousAndConflictingInputs|TestProfileDraftPlanIncludesCanonicalRepositoryProfile|TestProfileDraftPlanAcceptsMatchingSourceBaselineWithoutTransition|TestProfileDraftStaleTargetProducesNoMutation|TestProfileDraftRollbackRestoresMissingProfile)$' -count=1 -timeout=3m
rtk proxy go test ./internal/cli -run '^(TestBaselineProfileValidateNamesADecisionGap|TestBaselinePlanRefusesAProfileDecisionGap|TestBaselineUpdateRefusesAProfileDecisionGap|TestBaselineUpdateAsksForTheDecisionOnceTheProfileListsIt|TestBaselineProfileCommandInitShowAndValidate|TestBaselineProfileCommandRejectsInvalidAndNonRepositoryProfiles|TestBaselineUpdateNamesTheMissingIncrementalDecision|TestBaselineHumanProfileAdaptation|TestProfileAdaptationJourney|TestBaselinePlanProfileFile)$' -count=1 -timeout=3m
rtk proxy go test ./internal/baseline -run '^TestProfileDecisionGap' -count=1 -timeout=3m
rtk proxy git diff --check
```

The last baseline command also exercised the ordering/suggestion and existing
render-value assertions added after the combined compatibility run. Before the
consumer changes, the focused CLI probe reproduced validation reporting the gap
profile valid and update failing with the missing render token error.

`gofmt` was applied to the five changed Go files. Initial Git status contained
only this Task file's Daemon-owned `in_progress` transition; that status is
preserved. No other Task, Task Graph, catalog asset, guide or skill was edited.
Declared Verification, the full suites and repository Verification were not run;
they remain Daemon-owned. No commit, push or pull request was made.
