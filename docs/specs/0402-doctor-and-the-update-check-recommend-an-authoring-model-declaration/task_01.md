---
task: task_01
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
status: pending
type: backend
complexity: medium
---

# Task 01: Doctor and the Update Check recommend an Authoring Model Declaration when none is declared

## Overview

A repository with an installed authoring skill it does not publish and no
`skills.authoring.claude` starts getting the `undeclared` report. The
binary gains the Recommended Authoring Declaration as one table in
`internal/config` with its YAML rendering; Doctor prints
`authoring-models: warn` and the block as indented continuation lines; the
Update Check prints the line and the block in text and the `undeclared`
state in JSON. Every case reported today keeps its bytes. The Doctor and
Update Check guides and the Behavior Surface Record describe the new
output.

## Requirements

1. MUST add `internal/config/authoring_recommendation.go` with
   `RecommendedAuthoringDeclaration() AuthoringHarness`, returning a fresh
   copy on every call, and `RecommendedAuthoringYAML() string`, with the
   values and the exact bytes of `_techspec.md` → Recommendation and
   Surface Transcript 1 (the block lines without their two leading spaces,
   ending with a newline), and no `tools` key (invariant 5).
2. MUST add `internal/config/authoring_recommendation_test.go` with
   `TestRecommendedAuthoringDeclarationRendersValidProjectConfig`, calling
   `t.Parallel()` first, with the subtests `exact_bytes` (the whole
   rendering, compared as one string), `loads_as_project_config` (the
   rendering written as `.roundfixrc.yml` in a temporary repository loads
   with the existing loader, and its `Skills.Authoring.Claude` deep-equals
   `RecommendedAuthoringDeclaration()` with `Research.Tools` empty) and
   `returns_a_copy` (changing a returned entry leaves the next call
   unchanged).
3. MUST add `Lines []string` to `CheckResult` and make `printDoctorResult`
   print each entry after the result line with two leading spaces; a result
   with no `Lines` prints the bytes it printed before (invariant 7).
4. MUST add the `undeclared` report in `internal/cli/authoring_reports.go`
   exactly under invariants 1 to 4 of `_techspec.md` → Interfaces: Doctor's
   line and continuation lines of API Contract 1; the Update Check text of
   API Contract 2 and the JSON member of API Contract 3, with
   `Installed []string` and `Recommendation string` added after `Next` in
   `updateCheckAuthoring`, both `omitempty`. Every report that exists
   today, a problem from the inspection included, MUST keep its bytes.
5. MUST add `internal/cli/authoring_recommendation_test.go` with three
   tests, each calling `t.Parallel()` first:
   - `TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`, with
     the subtests `undeclared` (the whole warn line after `skills:`, then
     every block line in order before the next check's line, the exit code
     unchanged, and the repository unchanged by `snapshotDoctorPath`),
     `partly_published` (with `skills/write-prd/SKILL.md` and
     `skills/write-tasks/SKILL.md` in the tree, the line names only
     `write-idea, write-techspec`), `published` (all four published: no
     `authoring-models:` line), `no_authoring_skill` (no line),
     `leftover_generated` (no declaration and a marked
     `.claude/agents/roundfix-prd-author.md`: today's `unapplied` line and
     no block) and `declared` (today's `ok` line and no block);
   - `TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`,
     with the same six subtests asserting the whole `authoringModels` JSON
     member as bytes, the report's `next` unchanged, and the text line and
     block in the position of today's `Authoring models:` line;
   - `TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration`, asserting
     Surface Transcript 1 of `_techspec.md` through `runCLI` on a
     development-version fixture: every line of standard output in order,
     with the `...` regions matched as the transcript allows, empty
     standard error and exit 0.
6. MUST change `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted`
   so its `doctor` and `update_check` cases publish the four authoring
   skills under `skills/<name>/SKILL.md` before running, keeping their
   assertions; and change the expected output of
   `TestRunDoctorRealRepositoryCheckDoesNotMutateState` to include the
   warn line and the block lines after `skills:`. These are the two
   declared expectation changes of `_techspec.md` → Testing Approach.
7. MUST describe the `undeclared` report in
   `docs/user-guide/commands/doctor.md` and
   `docs/user-guide/commands/update.md`: the Doctor line with its block,
   the Update Check line with its block, the JSON state `undeclared` and
   the members `installed` and `recommendation` in the field table, that
   neither changes an exit code or `next`, that nothing is written, and
   that a published skill does not count. Then run the
   `-record-skill-coverage` step to re-record
   `docs/references/behavior-surfaces.json`.
8. MUST change no Governed Path and no file outside the ones this Task
   declares, except the one Governed Path the authorization record's
   amendment of 2026-11-17 grants: in `internal/cli/cli_test.go`, only
   `assertCheckResult` changes, comparing `CheckResult` with
   `reflect.DeepEqual` instead of `!=`, because `CheckResult` gains the
   `Lines []string` field of requirement 3.

## Subtasks

- [ ] Add the Recommended Authoring Declaration and its rendering, with its test.
- [ ] Add continuation lines to Doctor's result printing.
- [ ] Add the `undeclared` report to Doctor and the Update Check.
- [ ] Add the report tests and the transcript test, and update the two changed tests.
- [ ] Describe the report in the Doctor and Update Check guides and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] With an installed authoring skill the repository does not publish and no declaration, Doctor warns with the block and the Update Check reports `undeclared` with the block, in text and JSON.
- [ ] Neither check changes its exit code, the Update Check's `next`, or any file.
- [ ] A declaration, a leftover generated file, published skills, no installed authoring skill, Project Config not read and outside a Git repository each keep today's bytes.
- [ ] The rendered block loads as Project Config and yields the recommended entries with no tools.

## Context

- instruction: `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`
- instruction: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- instruction: `docs/adr/0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md`
- instruction: `skills/authoring_inspect.go`
- instruction: `internal/config/recommendations.go`
- creates: `internal/config/authoring_recommendation.go`
- creates: `internal/config/authoring_recommendation_test.go`
- interface: `internal/cli/authoring_reports.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor.go`
- creates: `internal/cli/authoring_recommendation_test.go`
- interface: `internal/cli/authoring_models_test.go`
- interface: `internal/cli/authoring_reports_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `config.RecommendedAuthoringYAML` and `config.RecommendedAuthoringDeclaration` — consumers: `./internal/cli`
- `CheckResult` and the `authoringModels` member of the Update Check report — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^TestRecommendedAuthoringDeclarationRendersValidProjectConfig$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; cli="$(go test -count=1 -timeout 30m -v -run '^(TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared|TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared|TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration)$' ./internal/cli 2>&1)" || { printf '%s\n' "$cli"; exit 1; }; out="$out$cli"; for name in TestRecommendedAuthoringDeclarationRendersValidProjectConfig/exact_bytes TestRecommendedAuthoringDeclarationRendersValidProjectConfig/loads_as_project_config TestRecommendedAuthoringDeclarationRendersValidProjectConfig/returns_a_copy TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/undeclared TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/partly_published TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/published TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/no_authoring_skill TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/leftover_generated TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/declared TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/undeclared TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/partly_published TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/published TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/no_authoring_skill TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/leftover_generated TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared/declared TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/config ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails on the first missing pass; afterwards each case passes by name, and both packages pass whole, the two changed tests included.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- "undeclared" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/update.md "undeclared" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- "installed without an Authoring Model Declaration" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/doctor.md "installed without an Authoring Model Declaration" >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide names the report, so the command fails; afterwards both do and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals 1, 2, 3 and 6; Core Features 1 and 2; Success Metric 1; Success Metric 2; Success Metric 3; Non-Goals
- `_techspec.md` → Recommendation; Reports; Interfaces (invariants 1 to 5 and 7); Data Models; API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Testing Approach (Task 1, Declared expectation changes); Build Order 1
- ADR-0415, ADR-0401, ADR-0408, ADR-0402
