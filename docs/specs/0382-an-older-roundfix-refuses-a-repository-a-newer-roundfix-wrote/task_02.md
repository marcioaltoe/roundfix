---
task: task_02
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
status: completed
type: backend
complexity: high
---

# Task 02: A Baseline plan records its writer's release and the Baseline commands refuse a newer one

## Overview

The Setup Manifest gains the optional `roundfixVersion`. A Baseline plan the
CLI builds carries the recorded value and records the running release only
when the plan changes a file and that release is higher; the Baseline commands
refuse a repository whose record is newer, before they read anything else, and
the Update Check's Baseline preview shows that refusal. This answers the
message reported by conexus on 2026-11-11
(`references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md`)
and the measured downgrade: Roundfix 0.137.7 planned to rewrite this
repository's managed guides back to its own text.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `parseManagedSetupManifest` rejects a Setup Manifest with an unknown field,
   and no plan builder passes a Roundfix release to `BuildPlan`.
2. MUST change `internal/baseline/plan.go` and
   `internal/baseline/classification.go` as `_techspec.md` → Data Models and
   API Contracts 6 and 7 say: `SetupManifest` gains `RoundfixVersion string`
   with JSON tag `roundfixVersion,omitempty`, declared last; `PlanRequest`
   gains `WriterVersion string`; `BuildPlan` carries the recorded release
   (read from the current Setup Manifest; absent, unreadable or malformed
   carry nothing) and, only after `FileChanges` is derived non-empty and the
   normalized writer version is a release number above the carried one (or
   none is carried), records the writer's release, re-marshals the manifest,
   re-assembles the postimages and ledger with `assemblePostimages`, derives
   `FileChanges` again, then computes the Plan Digest; `parseManagedSetupManifest`
   and `validateSetupManifest` reject a non-empty value that is not a release
   number, the latter with
   `Setup Manifest roundfixVersion must be a release version such as 0.140.0`.
   Nothing ever lowers the value.
3. MUST pass `WriterVersion: app.Version` from the plan builders of `roundfix
   baseline update` (`internal/cli/baseline_update.go`), `roundfix baseline
   plan` (`internal/cli/baseline_profile.go`) and both of `roundfix baseline`
   (`internal/cli/baseline_human.go`), and from no other caller.
4. MUST create `internal/cli/baseline_version.go` with
   `refuseNewerBaselineRepository` as API Contract 9 says, and call it right
   after the arguments parse, before the catalog, the manifest, a plan file or
   a skill source is read, in `runBaselineUpdateCommandWithSkillsStage`,
   `runBaselineHumanCommandWithIO`, `runBaselinePlanCommand`,
   `runBaselineApplyCommand`, `runBaselineSkillsRestoreCommandWith` and
   `runBaselineSkillsReconcileCommand`, each writing the refusal the way API
   Contract 9 lists: `baseline update` with state `failed`, category
   `version` and next action
   `run roundfix upgrade, then rerun roundfix baseline update`; `roundfix
   baseline` through `writeBaselineHumanAction(baselineHumanActionResult("preflight", <refusal>, "run roundfix upgrade, then rerun roundfix baseline"), ...)`;
   `baseline plan` and `baseline apply` through their usual failure;
   `baseline skills restore` and `reconcile` with finding codes
   `restore.roundfix-version-newer` and `reconcile.roundfix-version-newer` and
   actions `Run roundfix upgrade, then rerun roundfix baseline skills restore.`
   and `Run roundfix upgrade, then rerun roundfix baseline skills reconcile.`;
   every one exits `2` and writes no file.
5. MUST create `internal/baseline/recorded_version_test.go` with
   `TestAPlanRecordsItsWriterVersionOnlyWithAFileChange` and
   `TestASetupManifestRoundfixVersionMustBeAReleaseVersion`, each calling
   `t.Parallel()` first, on the plan characterization repository
   (`newBaselinePlanCharacterizationRepository`,
   `baselinePlanCharacterizationRequest`, `baselinePlanCharacterizationBin`),
   with the cases `_techspec.md` → Testing Approach lists; replans use the
   Managed Refresh and the stale guide replaces `Fix root causes.` with `Fix
   root causes!` in `docs/agents/agent-instructions.md`.
6. MUST create `internal/cli/baseline_version_test.go` with
   `TestBaselineUpdateRecordsTheRunningRoundfixVersion`,
   `TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix` (subtests
   named `adopt`, `update`, `plan`, `apply`, `skills restore`, `skills
   reconcile`; `apply` passes a missing plan file, `skills reconcile` passes
   `--source` and a forty-character `--revision`) and
   `TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix` (Surface
   Transcript 2 whole, then the JSON assertions), each calling `t.Parallel()`
   first and reusing task_01's `newRecordedVersionRepository`. The expected
   release reads `app.NormalizeVersion(app.Version)`, never a literal.
7. MUST change `docs/user-guide/commands/baseline.md` as `_techspec.md` → API
   Contract 10 says: right before the heading `#### Capability evidence and
   remediation`, the heading `#### Recorded Roundfix Version` and a paragraph
   that says the Setup Manifest's `roundfixVersion` records the release of the
   Roundfix that last applied a Baseline change, that a plan records the
   running release only when it changes a file and that release is higher,
   that a plan without a file change keeps it, that development and
   pre-release builds record none, and that `roundfix baseline`, `baseline
   update`, `baseline plan`, `baseline apply`, `baseline skills restore` and
   `baseline skills reconcile` refuse a newer value before they read anything
   else, exit `2` and change no file; then this `text` block, word for word:

   ```text
   roundfix: baseline update failed: Setup Manifest "<path>" records Roundfix <recorded>, newer than this binary's <running>; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'
   Category: version
   Next action: run roundfix upgrade, then rerun roundfix baseline update
   ```

8. MUST change `docs/user-guide/commands/update.md`: right before the
   paragraph that begins "Only when a release is newer or the Baseline has a
   plan", the sentence of API Contract 10 word for word: "When a newer
   Roundfix wrote this repository's Baseline, the check prints `roundfix:
   warning: <refusal>; Project Config not read` on stderr, reads no Project
   Config, and its Baseline line is `Baseline: failed:` followed by the
   refusal."
9. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST update the three tests that compare a plan the CLI builds with one
    `BuildPlan` builds directly, which the writer version otherwise makes
    differ (declared break): add `WriterVersion: app.Version` to the direct
    `BuildPlan` request of `TestBaselinePlanProfileFile`
    (`internal/cli/baseline_plan_test.go`, importing `roundfix/internal/app`)
    and of `TestHumanAutomationPlanParity` and
    `TestBaselineHumanProfileAdaptation` (`internal/cli/baseline_human_test.go`),
    and change nothing else in them; `TestBaselineMacroJourneysPublicCLI` runs
    the second as a subtest. Every other plan golden, the parity corpus,
    `docs/agents/setup-context.json` and every other existing test stay
    byte-identical: a plan built without a writer version records nothing. Any
    other test on the default branch when the Task runs that compares a CLI
    plan or the Setup Manifest it writes with bytes built without a writer
    version is changed the same way and named in the `## Result`.

## Subtasks

- [ ] Record the characterization.
- [ ] Add the field, the writer rule, the validation and the two Baseline tests.
- [ ] Pass the writer version from the four CLI plan builders.
- [ ] Add the Baseline guard to the six commands, with the three command tests.
- [ ] Add the Baseline and Update Check guide text, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A plan that changes a file records a higher writer release; a plan that changes none keeps the value; nothing lowers it.
- [ ] The six Baseline commands refuse a newer record, exit `2` and change no file.
- [ ] Every plan golden, parity fixture and this repository's Setup Manifest keep their bytes.

## Context

- instruction: `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`
- instruction: `internal/baseline/plan_characterization_test.go`
- instruction: `internal/baseline/apply.go`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/cli/baseline_plan_test.go`
- creates: `internal/baseline/recorded_version_test.go`
- creates: `internal/cli/baseline_version.go`
- creates: `internal/cli/baseline_version_test.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/classification.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/baseline_human.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/cli/baseline_apply.go`
- interface: `internal/cli/baseline_skills_restore.go`
- interface: `internal/cli/baseline_skills_reconcile.go`
- interface: `internal/cli/baseline_plan_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `SetupManifest.RoundfixVersion`, `PlanRequest.WriterVersion` and the writer rule — consumers: `./internal/baseline`, `./internal/cli`, `./internal/docscontract`
- The Baseline guard and its refusals — consumers: `./internal/cli`
- The two guides and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 40m -v -run '^(TestAPlanRecordsItsWriterVersionOnlyWithAFileChange|TestASetupManifestRoundfixVersionMustBeAReleaseVersion|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPlanRecordsItsWriterVersionOnlyWithAFileChange TestASetupManifestRoundfixVersionMustBeAReleaseVersion; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'` — expected: exit 0. Before this Task the two Baseline tests do not exist; afterwards they pass, every plan golden and parity fixture still matches, and this repository's Setup Manifest is still current.

- `out="$(go test -count=1 -timeout 40m -v -run 'RecordedByANewerRoundfix|RecordsTheRunningRoundfixVersion|BaselineUpdate|BaselinePlan|BaselineApply|BaselineHuman|BaselineSkills|BaselineMacroJourneys|HumanAutomationPlanParity|UpdateCheck' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateRecordsTheRunningRoundfixVersion TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix TestACommandRefusesARepositoryRecordedByANewerRoundfix TestBaselinePlanProfileFile TestHumanAutomationPlanParity TestBaselineHumanProfileAdaptation; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three command tests do not exist; afterwards they pass with task_01's tests, the three parity tests with the writer version, and every existing Baseline and Update Check test of `./internal/cli`.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- '#### Recorded Roundfix Version' || { echo 'baseline guide unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'Next action: run roundfix upgrade, then rerun roundfix baseline update' || { echo 'baseline guide block missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qE -- 'reads no Project Config, and its Baseline line is .Baseline: failed:. followed by the refusal' || { echo 'update guide unchanged' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide describes the record; afterwards both do, and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goals; User Story 2; User Story 3; Core Features 1, 2, 4; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7
- `_techspec.md` → Current behavior 2, 3, 6; Interfaces; Data Models; API Contracts 6-10; Surface Transcript 2; Invariants 2-4; Declared breaks; Testing Approach; Build Order 2
- ADR-0392
- [The adopted message](references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md)

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging.

### Starting-tree characterization

- Before implementation, a temporary Go overlay test built the characterization adoption manifest, confirmed it parsed, added `roundfixVersion: "0.140.0"`, and confirmed `parseManagedSetupManifest` rejected it. `go test -count=1 -overlay <temporary overlay.json> ./internal/baseline -run '^TestTask02StartingManifestRejectsUnknownField$' -v` exited 0; the named test passed. The overlay changed no repository source.
- Inspection of the four production `BuildPlan` calls in `baseline_update.go`, `baseline_profile.go` and `baseline_human.go` (two) found no writer release input; `PlanRequest` had no `WriterVersion` field. The starting working tree contained only this Daemon-owned Task file's status change.

### Implementation and acceptance evidence

- Writer rule: `SetupManifest.RoundfixVersion` is the last optional JSON field, and `PlanRequest.WriterVersion` is supplied by only the four production CLI plan builders. Planning carries a valid existing release, derives changes before deciding whether to raise it, then rebuilds manifest postimages and the ledger before computing the digest. Both manifest validators reject non-release values.
- Refusal: the six commands check the repository immediately after successful argument parsing and before their remaining reads. Existing failure printers preserve each command's contract; update reports category `version`, human adoption reports `preflight`, and skill commands report their specified finding codes and actions. The Update Check reaches the same refusal through its preview.
- Criterion 1 focused check: `go test -count=1 ./internal/baseline -run '^(TestAPlanRecordsItsWriterVersionOnlyWithAFileChange|TestASetupManifestRoundfixVersionMustBeAReleaseVersion)$' -v` exited 0. Adoption without a writer omitted the field; adoption with a release applied `verified`; Managed Refresh replans by newer, empty, development and pre-release writers had no changes and kept the record. A stale root-cause guide kept the record for lower, equal, empty, development, pre-release and invalid writers; a normalized higher writer raised it in the manifest postimage. Strict parsing and validation accepted empty/release values and rejected the three specified invalid forms with the required validation diagnostic.
- Criterion 2 focused check: `go test -count=1 ./internal/cli -run '^(TestBaselineUpdateRecordsTheRunningRoundfixVersion|TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix|TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix|TestBaselinePlanProfileFile|TestHumanAutomationPlanParity|TestBaselineHumanProfileAdaptation)$' -v` exited 0. All six refusal subtests asserted exit 2, the JSON refusal and unchanged repository-tree identity; apply supplied a missing plan plus the required confirmation. Update recorded `app.NormalizeVersion(app.Version)` after clearing the fixture record, then converged with zero changes. The Update Check matched Surface Transcript 2 in full and asserted JSON `baseline.state`, `baseline.message` and `next`.
- The first command-test run exposed two new-fixture mistakes: apply omitted the existing mandatory `--confirm-plan`, and the zero-change `--yes` update expected `current` instead of the existing `verified` result. Reading the parser and existing idempotence test established those contracts; the fixtures were corrected without changing production behavior or existing assertions.
- Criterion 3 focused check: `go test -count=1 ./internal/baseline -run '^(TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$'` exited 0 without regeneration. The only existing tests edited are the three declared parity requests (and the specified `app` import); their original lines and assertions are preserved. No additional parity adjustment was needed in the focused journeys.
- Existing-behavior check: `go test -count=1 ./internal/cli -run '^(TestBaselineUpdateFleetSweep|TestBaselineUpdateIdempotenceReportsZeroFileChanges|TestBaselineMacroJourneysPublicCLI|TestACommandRefusesARepositoryRecordedByANewerRoundfix|TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix)$'` exited 0.
- The Baseline and Update Check guides contain the prescribed heading, explanatory text, failure block and exact Update Check sentence. No Baseline command was run against this checkout. Plan goldens, parity fixtures, this repository's Setup Manifest, Task Graph and other Task files remain unchanged. `git diff --check` exited 0 before regeneration.
- Declared Verification and repository Verification remain for the Daemon; no Task status, checkbox settlement, commit, push or pull request was performed. No follow-up outside this Task was required.
- Final regeneration: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 after the last guide edit. Its diff contains only the Baseline and Update Check guide fingerprints. Final `git diff --check` exited 0; the protected-byte diff for `internal/baseline/testdata`, `docs/agents/setup-context.json`, `_tasks.md` and the other Task files was empty. All changed and new paths are declared by task_02; the two governed parity test files contain only the authorized writer inputs and import.
