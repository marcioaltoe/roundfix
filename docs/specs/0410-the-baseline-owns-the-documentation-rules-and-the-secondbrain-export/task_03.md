---
task: task_03
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: completed
type: backend
complexity: high
---

# Task 03: The Doctor Command reports what departs from the documentation rules

## Overview

`roundfix doctor` gains two warning checks after `baseline-claude-link`: `baseline-export`, where the Secondbrain is turned on, compares `.secondbrain-export` with the Standard Export List the update would write and reads the repository's own lines; `baseline-docs` reads the documentation layout, the ADR headers and `docs/lessons/`. Each departure is one line with a stable code, neither check ever fails or writes, and the Roundfix skill's `setup` reference and the user guide describe both.

## Requirements

1. MUST create `internal/spec/lesson.go` with `LessonLearnedProblems` of `_techspec.md` → Interfaces: nil when the file follows the form of `_techspec.md` → Data Models (a `YYYY-MM-DD-<kebab-slug>.md` name with a valid date, a frontmatter `created` with a valid date, one `# ` title line, and the eight labelled lines once each, in order, each with text after the colon), and otherwise one problem per departure, worded as API Contract 12 words the problems of `DR-DOCS-LESSON-FORM`.
2. MUST create `internal/baseline/docs_layout.go` with `DocumentationDeparture`, `DocsLayoutFacts` and `DocsLayoutDepartures` of `_techspec.md` → Interfaces, reporting `DR-DOCS-REQUIRED-MISSING` and `DR-DOCS-OUTSIDE-LAYOUT` by the rules of API Contract 12: the folders the layout names under `docs/` and the first folder of the Spec Root and of the History Root; `docs/README.md` never reported; a root Markdown file read as a reference document only when its name before `.md` is in capitals, digits, hyphens and underscores, with the fifteen names the layout names; and a path that `docs/agents/specific-repository.md` contains never reported.
3. MUST create `internal/baseline/export_variation.go` with `ExportVariationDepartures` of `_techspec.md` → Interfaces, by the Data Models rule for the repository's own lines: an entry is a line that is neither blank nor a comment, its reason is the nearest comment line above it with no blank line between, and a marker line is not a reason. It reports `DR-EXPORT-PROCESS-PATH` and `DR-EXPORT-REASON-MISSING` with the line number. It changes no file task_02 created, so `PlanExportList` does not change; the TechSpec places the function in `export_list.go`, and a file of its own keeps that file task_02's.
4. MUST create `internal/cli/doctor_documentation.go` with a function of the repository root that returns the two results `baseline-export` and `baseline-docs`, in that order, from the facts task_02's gathering function returns, called without changing it, plus the root Markdown files and `docs/agents/specific-repository.md`, which this file reads itself, and that `runDoctorCommand` appends right after `baseline-claude-link`; `internal/cli/health.go` gains the two names. The results follow API Contracts 11 and 12 word for word: the `skipped` details, the `ok` details, the `warn` details and next actions, one line per departure in the forms of the eleven codes, `DR-EXPORT-UNREAD` for a read error, and at most 20 lines per code followed by `<code>: <n> more`. `DR-EXPORT-LIST-STALE` compares with the plan `PlanExportList` makes from the same facts; `DR-EXPORT-ADR-UNMARKED` and `DR-DOCS-ADR-RETIRED` read the headers through `ReadADRHeader` and `ClassifyADR`, so a proposed ADR is never reported (ADR-0122); `DR-DOCS-ADR-AMENDMENT-UNLISTED` compares `ADRAmendmentNotices` with the Amendment List. No result is ever `failed`, neither check writes a file, and the Doctor Command's exit code stays the one the other checks give (`_techspec.md` → Invariant 10; ADR-0220).
5. MUST add these tests, on temporary directories and repositories, with no network, no real home and no dependence on file modes or the platform:
   - `internal/spec/lesson_test.go`: `TestLessonLearnedProblems` with subtests `valid`, `name_not_dated`, `name_bad_date`, `no_created`, `no_title`, `label_missing`, `label_empty` and `label_out_of_order`.
   - `internal/baseline/docs_layout_test.go`: `TestDocsLayoutDepartures` with subtests `named_folders`, `spec_and_history_roots`, `unnamed_folder`, `loose_file`, `docs_readme_not_reported`, `root_reference_document`, `root_names_the_layout_names`, `lower_case_root_file`, `named_by_the_repository_guide`, `readme_missing` and `glossary_missing`.
   - `internal/baseline/export_variation_test.go`: `TestExportVariationDepartures` with subtests `process_path`, `reason_missing`, `comment_stands_for_following_lines`, `marker_not_a_reason` and `blank_line_breaks_the_reason`.
   - `internal/cli/doctor_documentation_test.go`: `TestDoctorReportsTheDocumentationRules`, driving the Doctor Command's entry point with the fakes the existing Doctor tests use, with one subtest per seeded departure, `export_list_missing`, `export_list_stale`, `export_marker`, `export_process_path`, `export_reason_missing`, `export_adr_unmarked`, `export_unread`, `docs_required_missing`, `docs_outside_layout`, `docs_adr_retired`, `docs_adr_amendment_unlisted` and `docs_lesson_form`, each asserting its whole line, the `warn` status and the next action; and `both_ok`, `skipped_without_manifest`, `export_skipped_secondbrain_off`, `twenty_lines_then_more`, `exit_unchanged`, `named_by_the_repository_guide` and `order_after_claude_link` (`_prd.md` → Success Metric 4).
6. MUST update the existing Doctor tests that pin the list, the count or the order of the report lines to the two new lines right after `baseline-claude-link`, and change nothing else in them: `internal/cli/doctor_test.go` (the whole-report fixtures without a Setup Manifest, the line count and the ordered check names), `internal/cli/doctor_baseline_rules_test.go` (the line after `baseline-claude-link`) and `internal/cli/doctor_claude_link_test.go` (the line after `baseline-claude-link`). Any other existing test that pins the Doctor Command's lines on the tree the Task runs on is updated the same way.
7. MUST describe the change: the canonical Roundfix `setup` reference gains the paragraph of `_techspec.md` → API Contract 18 in the part that describes `roundfix doctor`; `docs/user-guide/commands/doctor.md` gains the bullets `baseline-export` and `baseline-docs` after `baseline-claude-link`, the eleven codes in its code table, each with its meaning and what repairs it, and the subsection `#### Documentation rules` with the line forms of API Contracts 11 and 12 and the sentence "Both checks warn and neither changes the exit code." (API Contract 20).
8. MUST regenerate, never by hand: `make skills-sync`; the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of both copies of the Roundfix skill's `SKILL.md` again on this Task's tree; and, after the last guide edit, the Behavior Surface Record with the `-record-skill-coverage` step, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change any file task_02 created, the update's documentation part, any Baseline module, template, profile, golden or Source Baseline, any guide under `docs/agents/`, any ADR, `.secondbrain-export`, any skill other than the Roundfix skill's `setup` reference and version fields, `CONTEXT.md`, `internal/cli/cli_test.go`, `internal/cli/baseline_human_test.go`, `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, all listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Add the lesson form check with its test.
- [ ] Add the layout departures and the variation departures with their tests.
- [ ] Add the two Doctor checks after `baseline-claude-link`, with their tests, and update the Doctor tests that pin the report's lines.
- [ ] Describe the checks in the `setup` reference and the Doctor guide; sync the skills, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] On temporary repositories the Doctor Command prints each of the eleven codes for its seeded departure as a warning with its next action, prints `ok` for both checks on a repository that follows the rules, skips the export check where the Secondbrain is turned off, does not report a path the repository's own rules guide names, and exits as it did before.
- [ ] Both new lines come right after `baseline-claude-link`, and no other line of the report changes.
- [ ] Both copies of the `setup` reference and the Doctor guide describe the two checks and their codes.

## Context

- instruction: `internal/spec/retirement.go`
- creates: `internal/spec/lesson.go`
- creates: `internal/spec/lesson_test.go`
- creates: `internal/baseline/docs_layout.go`
- creates: `internal/baseline/docs_layout_test.go`
- creates: `internal/baseline/export_variation.go`
- creates: `internal/baseline/export_variation_test.go`
- creates: `internal/cli/doctor_documentation.go`
- creates: `internal/cli/doctor_documentation_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/doctor_baseline_rules_test.go`
- interface: `internal/cli/doctor_claude_link_test.go`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The report lines of `roundfix doctor`, their order and the two check names — consumers: `./internal/cli`
- `LessonLearnedProblems` — consumers: `./internal/spec`, `./internal/cli`
- `DocsLayoutDepartures` and `ExportVariationDepartures` — consumers: `./internal/baseline`, `./internal/cli`
- The Roundfix skill's `setup` reference, the Doctor guide and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestLessonLearnedProblems$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in valid name_not_dated name_bad_date no_created no_title label_missing label_empty label_out_of_order; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestLessonLearnedProblems/$name " || { echo "missing pass: TestLessonLearnedProblems/$name" >&2; exit 1; }; done; bl="$(go test -count=1 -timeout 20m -v -run '^(TestDocsLayoutDepartures|TestExportVariationDepartures|TestPlanExportList)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$bl"; exit 1; }; for name in TestDocsLayoutDepartures/named_folders TestDocsLayoutDepartures/spec_and_history_roots TestDocsLayoutDepartures/unnamed_folder TestDocsLayoutDepartures/loose_file TestDocsLayoutDepartures/docs_readme_not_reported TestDocsLayoutDepartures/root_reference_document TestDocsLayoutDepartures/root_names_the_layout_names TestDocsLayoutDepartures/lower_case_root_file TestDocsLayoutDepartures/named_by_the_repository_guide TestDocsLayoutDepartures/readme_missing TestDocsLayoutDepartures/glossary_missing TestExportVariationDepartures/process_path TestExportVariationDepartures/reason_missing TestExportVariationDepartures/comment_stands_for_following_lines TestExportVariationDepartures/marker_not_a_reason TestExportVariationDepartures/blank_line_breaks_the_reason TestPlanExportList; do printf '%s\n' "$bl" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec` — expected: exit 0. Before this Task the three tests do not exist; afterwards every lesson form, layout and variation case passes, the planner's cases still pass, and the `internal/spec` package passes whole.
- `out="$(go test -count=1 -timeout 40m -v -run '^TestDoctorReportsTheDocumentationRules$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in export_list_missing export_list_stale export_marker export_process_path export_reason_missing export_adr_unmarked export_unread docs_required_missing docs_outside_layout docs_adr_retired docs_adr_amendment_unlisted docs_lesson_form both_ok skipped_without_manifest export_skipped_secondbrain_off twenty_lines_then_more exit_unchanged named_by_the_repository_guide order_after_claude_link; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestDoctorReportsTheDocumentationRules/$name " || { echo "missing pass: TestDoctorReportsTheDocumentationRules/$name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 40m -run 'Doctor|Baseline' ./internal/cli` — expected: exit 0. Before this Task the test does not exist; afterwards each of the eleven codes, the read error, both `ok` and both `skipped` details, the cap, the unchanged exit code and the order pass through the command's entry point, and the Baseline package and the CLI's Doctor and Baseline tests, the updated ordering tests included, pass.
- `for pair in '.agents/skills/roundfix/references/setup.md|reports two documentation checks as warnings that never change its exit.' 'skills/roundfix/references/setup.md|and the rest is an edit by hand.' 'docs/user-guide/commands/doctor.md|Both checks warn and neither changes the exit code.' 'docs/user-guide/commands/doctor.md|DR-EXPORT-REASON-MISSING' 'docs/user-guide/commands/doctor.md|DR-DOCS-ADR-AMENDMENT-UNLISTED'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills && go test -count=1 -timeout 20m -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task neither the reference nor the guide describes the checks; afterwards both copies of the reference and the guide do, the owned-skill versions and the Behavior Surface Record are current, and the skill mirrors match.

## References

- `_prd.md` → Goals; User Stories 3 and 11; Core Features 3, 10, 11, 15 and 17; Success Metric 4
- `_techspec.md` → ADR header and the lesson form; Export list; Doctor checks; Interfaces; Invariant 10; Data Models; API Contract 11; API Contract 12; API Contract 18; API Contract 20; Testing Approach; Build Order 3
- ADR-0220, ADR-0411, ADR-0122, ADR-0438, ADR-0439, ADR-0440


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging.

Implemented this Task's slice; settlement and declared Verification remain with the Daemon.

- Acceptance criterion 1: added the pure lesson-form, layout and export-variation readers and the two advisory Doctor checks. Doctor calls the existing documentation-facts gatherer and export planner without changing either. The command-entry tests assert whole diagnostic lines, warning headers and next actions for all eleven departure codes and the read error, both successful and skipped details, repository-guide exceptions, proposed ADRs (including legacy status), and unchanged exits. Repository and temporary-home snapshots assert no writes. The cap case seeds two interleaved codes with 23 departures each and asserts 20 full lines plus a separate `3 more` line per code, while the header still counts all 46 departures.
- Acceptance criterion 2: appended `baseline-export` and `baseline-docs` immediately after `baseline-claude-link`, and registered their names in health results. Updated existing pinned report fixtures, counts and header ordering only for the added checks and their indented details. Focused existing Doctor tests continue to assert the other report text and independent failures.
- Acceptance criterion 3: inserted API Contract 18's paragraph into the canonical setup reference and synchronized its mirror. The Doctor guide describes the checks, every departure code and read error, exact line forms, repairs, layout exceptions and cap. Both Roundfix SKILL.md copies changed only in generated version fields, from `0.1.279` to `0.1.280`; the version recorder added its digest, and the coverage recorder updated only the Doctor guide's Behavior Surface Record entry.

Focused implementation evidence:

- Before implementation, `rtk proxy go test ./internal/spec ./internal/baseline -run 'Test(LessonLearnedProblems|DocsLayoutDepartures|ExportVariationDepartures)/valid' -count=1` exited 1 with the new APIs undefined, establishing the starting gap.
- `rtk proxy go test ./internal/spec ./internal/baseline -run 'Test(LessonLearnedProblems|DocsLayoutDepartures|ExportVariationDepartures)/' -count=1` exited 0 for both packages after the final helper edits. It exercised all requested subtests and additional invalid-created-date and duplicate-label cases.
- `rtk proxy go test ./internal/cli -run '^TestDoctorReportsTheDocumentationRules/' -count=1 -timeout 5m` exited 0 after the final code edits, including a comment that mentions a marker without creating a Standard Export List.
- `rtk proxy go test ./internal/cli -run '^(TestDoctorReportsTheDocumentationRules|TestDoctorReportsMisplacedRepositoryRules|TestDoctorReportsTheClaudeInstructionLink|TestRunDoctor.*|TestDoctor.*Readiness.*)$' -count=1 -timeout 5m` exited 0 after the final code edits. Initial runs exposed a malformed positive lesson fixture, absent-root handling and old ordering assertions that assumed no warning-detail lines; the fixture, missing-manifest skip and header-order expectations were corrected, retaining the warning diagnostics.
- Required generation: `rtk proxy make skills-sync` and `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` exited 0. `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 and was run again after the final guide formatting edit.
- Read-only scope and content assertions found all 21 changed paths declared by this Task, equal setup-reference mirrors, SKILL.md differences confined to version fields, the prescribed reference and guide sentences, and every documented Doctor code. `rtk proxy git diff --check` exited 0.

The pre-existing Task status edit was preserved. No Task status, Task Graph, other Task, Task 02-created file, Baseline source or generated Baseline artifact was changed by this turn. The declared Verification commands, repository Verification and incremental Verification were not run. No commit, push, Pull Request or cross-project message was made.

## Carry-forward provenance

- Source Run: `run_20261119T213428Z_8e0d266520442009`
- Source commit: `4c649e2b0c7a9bd3e825e878dfd68f3363414d5a`
