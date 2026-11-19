---
task: task_02
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: completed
type: backend
complexity: high
---

# Task 02: The update keeps the Standard Export List and each Amendment List, and its commit asks for no release

## Overview

`roundfix baseline update` gains its documentation part: where the Secondbrain is turned on it writes the Standard Export List between the two markers of `.secondbrain-export`, converting the list a repository had the first time and naming every path that stops being exported; in every repository it fills an ADR's Amendment List from the notices in its body; the Plan Digest covers the part and the one confirmation applies it. The Release Plan counts `.secondbrain-export` and `docs/lessons/` as needing no release. The Roundfix skill, the user guide and the `release plan` help describe both, and this repository runs the new part once so that its own update stays current.

## Requirements

1. MUST create `internal/spec/adr_header.go` with `ADRHeader`, `ReadADRHeader`, `ADRAmendmentNotices` and `WithADRAmendedBy` as `_techspec.md` → Interfaces and Data Models declare them, pure functions over bytes that edit lines and never round-trip YAML (`_techspec.md` → Decisions), holding Invariants 6, 7 and 8:
   - a notice is "Amended by", "Narrowed by", "Extended by" or "Superseded in part by" followed by one or more `ADR-<digits>` identifiers, plain or as the text of a Markdown link, joined by a comma or `and`, read outside fenced code and inline code only;
   - the writer puts `amended_by` after `knowledge` when the header has it and after `status` otherwise, replaces an existing `amended_by` line in place, writes the union of the listed and the noticed identifiers ascending by number, never removes an identifier a person wrote, keeps `updated_at` and every other byte, keeps CRLF line endings, and prepends to a file without frontmatter a header that holds only that line;
   - the caller never lists an ADR in its own Amendment List.
2. MUST create `internal/baseline/export_list.go` with the constants, `ExportListFacts`, `ExportLoss`, `ExportListPlan` and `PlanExportList` of `_techspec.md` → Interfaces, holding Invariants 1 to 5 and 9 and the Data Models rules: the standard lines in the order the Data Models give, every standard folder listed whether or not it exists, the fixed second line, a new file as the block and a final newline, a converted file as the block, one blank line and the kept lines in their order, the conversion removals and the `ExportLoss` kinds, and the five refusals of API Contract 10 word for word. It takes its facts as values and reads no configuration and no file. `ExportVariationDepartures` and `DocumentationDeparture` are task_03's, in files of its own.
3. MUST create `internal/cli/baseline_update_documentation.go` with the documentation part of `_techspec.md` → API Contract 10, called from `runBaselineUpdateCommandWithSkillsStage` in `internal/cli/baseline_update.go` where the history plan is planned, digested and applied:
   - it gathers the facts once, in a function of the repository root that task_03's Doctor checks call without changing it: the Setup Manifest's modules, the Spec Root and the History Root from the configuration, the glossary files through the consistency check's resolver, which `internal/speccheck/glossary.go` exports for it without changing what the check reads, each ADR under `docs/adr/` with its path, its `ClassifyADR` result, its `ReadADRHeader` header and its `ADRAmendmentNotices`, and the immediate children of `docs/`;
   - where the Setup Manifest's modules lack `secondbrain` the export list's action is `none` and the file is neither read nor written; a retired ADR is neither listed nor given a header;
   - the `documentation` member of the JSON result follows `history`, with the shape of `_techspec.md` → Data Models: `status` first, `exportList` with `path`, `action`, `standard`, `added`, `removed`, `converted`, `stopsExporting` and `refusal`, `adrs` with `total` then `unmarked`, and `adrHeaders`; every array present;
   - the Plan Digest moves only when the status is `pending`, by the rule of API Contract 10; with no guidance change, no history unit and a pending part the state is `plan_ready` with the message `guidance matches the current Baseline catalog; <n> documentation change(s) pending`;
   - each file it writes is one more `fileChanges` entry, with the managed entry `export.standard-list` or `adr.amendment-list`, so the Update Check counts it with no change of its own;
   - the apply runs after the guidance and the history, fills each header over the file's bytes at that moment, skips a file the history removed, plans the export list again from the tree and writes it; a failure reports the category `documentation` and the next action `restore .secondbrain-export and docs/adr with git restore, then rerun roundfix baseline update`, with exit 1;
   - the text section is printed after the history section with the lines of API Contract 10, word for word and in that order, each only when it has something to say;
   - `--no-history` and `--no-skills` do not skip it, and no flag, configuration key or exit code is added or changed.
4. MUST make `MaintenanceOnly` in `internal/releaseplan/classifier.go` true for a path under `docs/lessons/` and for `.secondbrain-export`, in every repository and reading no configuration, and change no other classification (API Contract 13).
5. MUST add these tests, each calling `t.Parallel()` first where its package's tests do, all on temporary directories and repositories, with no network, no real home and no dependence on file modes or the platform:
   - `internal/spec/adr_header_test.go`: `TestReadADRHeader` with subtests `mark_true`, `mark_false`, `mark_absent`, `mark_not_boolean`, `amended_by_listed`, `no_frontmatter` and `fields_anywhere_in_header`; `TestADRAmendmentNotices` with `amended_by_one`, `amended_by_two`, `narrowed_by`, `extended_by`, `superseded_in_part_by`, `linked_identifier`, `joined_by_and`, `fenced_block_ignored`, `inline_code_ignored` and `ascending_once`; `TestWithADRAmendedBy` with `after_knowledge`, `after_status`, `replaces_the_list`, `keeps_written_identifiers`, `ascending_by_number`, `keeps_updated_at`, `crlf`, `no_frontmatter_gains_header` and `same_bytes_twice`.
   - `internal/baseline/export_list_test.go`: `TestPlanExportList` with subtests `create`, `convert_the_template`, `convert_reports_losses`, `replace_between_markers`, `unchanged`, `refuse_begin_missing`, `refuse_end_missing`, `refuse_begin_repeated`, `refuse_end_repeated`, `refuse_end_before_begin`, `same_bytes_same_plan`, `never_a_process_path`, `own_lines_kept` and `apply_then_plan_unchanged`.
   - `internal/cli/baseline_update_documentation_test.go`, driving `runBaselineUpdateCommand` on adopted temporary repositories with the Secondbrain turned on unless the subtest says otherwise: `TestBaselineUpdateDocumentationPart` with subtests
     - `plan_names_the_list_and_the_losses`: an export file listing `docs/` whole, `AGENTS.md`, one file of its own and one comment; the plan names the standard list, every path that stops being exported with its kind, and the count of unmarked ADRs (`_prd.md` → Success Metric 2);
     - `text_section`: every line of the text section of API Contract 10, in order;
     - `json_member`: the `documentation` member and its shape;
     - `digest_moves_only_when_pending`, `plan_ready_message` and `file_changes`;
     - `confirmed_apply`: the confirmed update writes the block, removes the lines it named and leaves every other line byte for byte;
     - `second_run_current`: a second run reports nothing to do;
     - `secondbrain_off`: the same repository with the Secondbrain turned off gets no file;
     - `marker_refused_rest_applies`: a file with one marker missing is left untouched while the rest of the plan applies;
     - `fills_amendment_lists`: each of the four notice forms with one and with two identifiers fills its ADR, a notice inside a code block fills nothing, no other byte changes, an ADR without frontmatter gains a header that holds only the list, and a second run reports nothing to do (`_prd.md` → Success Metric 3);
     - `fill_beside_history_unit`: the header fill beside a history unit that appends a line to the same successor keeps both changes;
     - `no_history_does_not_skip`;
     - `apply_failure`: a failure while the part is written reports the category and the next action of requirement 3 and exit 1, through a seam that does not depend on file modes.
   - `internal/releaseplan/classifier_test.go` gains `TestMaintenanceOnlyCountsLessonsAndTheExportList` with subtests `lessons`, `export_list` and `with_another_path`, the last proving that a commit that also changes another path still crosses the boundary (`_prd.md` → Success Metric 5).
   Any existing test that pins a line of the update's output or the release sentence on the tree the Task runs on is updated in this Task to the new text, and no other existing test changes.
6. MUST describe the change, inserting each text by the anchor `_techspec.md` names, as it reads on the default branch when the Task runs:
   - the canonical Roundfix `baseline` reference gains the paragraph of API Contract 17 in the part that describes `roundfix baseline update`;
   - the canonical Roundfix `release` reference, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md` gain the release sentence of API Contract 19, and the `release plan` help in `internal/cli/cli.go` gains its help sentence right after its sentence about planning evidence;
   - `docs/user-guide/commands/baseline.md` gains the subsection `### Documentation part` of API Contract 20, which states every rule of API Contract 10 and holds its four sentences.
7. MUST regenerate, never by hand: `make skills-sync`; the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of both copies of the Roundfix skill's `SKILL.md` to the next version (no number is written here); and, after the last guide edit, the Behavior Surface Record with the `-record-skill-coverage` step, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST, as the Task's last step, run `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` until a run reports nothing to do, so that this repository's update stays current: the run converts this repository's `.secondbrain-export` into the Standard Export List with its kept lines and fills the Amendment List of each ADR whose body carries a notice. The Task changes no ADR by hand; the ADR files the run writes are recorded by the Daemon. The repository's own lines of `_techspec.md` → API Contract 22, the Knowledge Marks and the three notices are task_05's.
9. MUST NOT change any Baseline module, template, profile, golden or Source Baseline, any guide under `docs/agents/`, any skill other than the Roundfix skill's `baseline` and `release` references and version fields, `CONTEXT.md`, `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_human_test.go`, `skills/baseline_skill_contract_test.go`, `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/release.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, all listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Add the ADR header reader, notice finder and line writer with their tests.
- [ ] Add the export list planner with its tests.
- [ ] Add the documentation part to the update: facts, plan, digest, JSON, text, file changes and apply, with its tests.
- [ ] Count the lessons folder and the export list as needing no release, with its test.
- [ ] Describe the part and the release rule in the skill, the user guide and the help; sync the skills, record the versions and the Behavior Surface Record.
- [ ] Run the Managed Refresh on this repository until it reports nothing to do.

## Acceptance Criteria

- [ ] On a temporary repository with the Secondbrain turned on, the update's plan names the Standard Export List, every path that stops being exported with its kind and the count of unmarked ADRs; the confirmed update writes the block, removes only the lines it named, and a second run reports nothing to do; with the Secondbrain turned off no file is written; a broken pair of markers leaves the file untouched while the rest applies.
- [ ] The update fills an Amendment List from each notice form, ignores a notice in code, changes no other byte and reports nothing on a second run.
- [ ] A commit that changes only `.secondbrain-export` and a file under `docs/lessons/` is maintenance-only, and one that also changes another path is not.
- [ ] This repository's update, run under a fresh home, reports `current` with a current documentation part.

## Context

- instruction: `docs/adr/0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md`
- instruction: `docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md`
- instruction: `internal/cli/baseline_update_history.go`
- instruction: `internal/spec/retirement.go`
- creates: `internal/spec/adr_header.go`
- creates: `internal/spec/adr_header_test.go`
- creates: `internal/baseline/export_list.go`
- creates: `internal/baseline/export_list_test.go`
- creates: `internal/cli/baseline_update_documentation.go`
- creates: `internal/cli/baseline_update_documentation_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/speccheck/glossary.go`
- interface: `internal/releaseplan/classifier.go`
- interface: `internal/releaseplan/classifier_test.go`
- interface: `internal/cli/cli.go`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `.secondbrain-export`

## Consumers

- `ReadADRHeader`, `ADRAmendmentNotices` and `WithADRAmendedBy` — consumers: `./internal/spec`, `./internal/cli`
- `PlanExportList` and the export list format — consumers: `./internal/cli`
- The `documentation` member, the text section and the `fileChanges` entries of `roundfix baseline update`, read by the Update Check and the `update` Skill Mode — consumers: `./internal/cli`
- `MaintenanceOnly` — consumers: `./internal/releaseplan`, `./internal/cli`
- The Roundfix skill's `baseline` and `release` references, the user guide pages and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReadADRHeader|TestADRAmendmentNotices|TestWithADRAmendedBy)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReadADRHeader/mark_true TestReadADRHeader/mark_false TestReadADRHeader/mark_absent TestReadADRHeader/mark_not_boolean TestReadADRHeader/amended_by_listed TestReadADRHeader/no_frontmatter TestReadADRHeader/fields_anywhere_in_header TestADRAmendmentNotices/amended_by_one TestADRAmendmentNotices/amended_by_two TestADRAmendmentNotices/narrowed_by TestADRAmendmentNotices/extended_by TestADRAmendmentNotices/superseded_in_part_by TestADRAmendmentNotices/linked_identifier TestADRAmendmentNotices/joined_by_and TestADRAmendmentNotices/fenced_block_ignored TestADRAmendmentNotices/inline_code_ignored TestADRAmendmentNotices/ascending_once TestWithADRAmendedBy/after_knowledge TestWithADRAmendedBy/after_status TestWithADRAmendedBy/replaces_the_list TestWithADRAmendedBy/keeps_written_identifiers TestWithADRAmendedBy/ascending_by_number TestWithADRAmendedBy/keeps_updated_at TestWithADRAmendedBy/crlf TestWithADRAmendedBy/no_frontmatter_gains_header TestWithADRAmendedBy/same_bytes_twice; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec` — expected: exit 0. Before this Task the three tests do not exist; afterwards every reader, notice and writer case passes, and the package passes whole.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestPlanExportList$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in create convert_the_template convert_reports_losses replace_between_markers unchanged refuse_begin_missing refuse_end_missing refuse_begin_repeated refuse_end_repeated refuse_end_before_begin same_bytes_same_plan never_a_process_path own_lines_kept apply_then_plan_unchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestPlanExportList/$name " || { echo "missing pass: TestPlanExportList/$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the planner and its test do not exist; afterwards every creation, conversion, replacement and refusal case and each invariant passes.
- `out="$(go test -count=1 -timeout 40m -v -run '^TestBaselineUpdateDocumentationPart$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in plan_names_the_list_and_the_losses text_section json_member digest_moves_only_when_pending plan_ready_message file_changes confirmed_apply second_run_current secondbrain_off marker_refused_rest_applies fills_amendment_lists fill_beside_history_unit no_history_does_not_skip apply_failure; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestBaselineUpdateDocumentationPart/$name " || { echo "missing pass: TestBaselineUpdateDocumentationPart/$name" >&2; exit 1; }; done; rp="$(go test -count=1 -timeout 20m -v -run '^TestMaintenanceOnlyCountsLessonsAndTheExportList$' ./internal/releaseplan 2>&1)" || { printf '%s\n' "$rp"; exit 1; }; for name in lessons export_list with_another_path; do printf '%s\n' "$rp" | grep -q -- "--- PASS: TestMaintenanceOnlyCountsLessonsAndTheExportList/$name " || { echo "missing pass: TestMaintenanceOnlyCountsLessonsAndTheExportList/$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/releaseplan && go test -count=1 -timeout 40m -run 'Baseline|UpdateCheck|ReleasePlan' ./internal/cli` — expected: exit 0. Before this Task neither test exists; afterwards the documentation part passes every case of Success Metrics 2 and 3 through the command's entry point, the Release Plan rule passes, and the release plan package and the CLI's Baseline, Update Check and Release Plan tests pass, this repository's own update transcripts included.
- `for pair in '.agents/skills/roundfix/references/baseline.md|It reports the count of ADRs without a Knowledge Mark and marks none.' 'skills/roundfix/references/baseline.md|A missing or repeated marker is refused: nothing is written to that file, and the rest of the plan applies.' '.agents/skills/roundfix/references/release.md|are planning evidence too, and the export list' 'skills/roundfix/references/release.md|are planning evidence too, and the export list' 'docs/user-guide/release-runbook.md|are planning evidence too, and the export list' 'docs/user-guide/usage.md|are planning evidence too, and the export list' 'docs/user-guide/commands/baseline.md|The Plan Digest covers the documentation part, and the one confirmation applies it.' 'docs/user-guide/commands/baseline.md|The first update converts the list the file had and names every path that stops being exported.'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix release plan --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | tr -s '[:space:]' ' ' | grep -qF -- 'Neither does docs/lessons or the export list .secondbrain-export.' || { printf '%s\n' "$help"; exit 1; }; go test -count=1 -timeout 20m -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills && go test -count=1 -timeout 20m -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task neither reference, guide nor help carries the new text; afterwards both copies of the two references, the three guide pages and the help carry it, the owned-skill versions and the Behavior Surface Record are current, and the skill mirrors match.
- `test "$(head -n 1 .secondbrain-export)" = '# roundfix:standard-export:begin' || { echo 'the export list has no Standard Export List' >&2; exit 1; }; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"documentation":{"status":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task this repository's export file has no markers and the update has no documentation part; afterwards the file opens with the Standard Export List, this repository's update under a fresh home reports `current` with a current documentation part, and the documentation contracts pass.

## References

- `_prd.md` → Goals; User Stories 1, 2, 4, 7 and 12; Core Features 1, 2, 4, 5, 8, 16 and 18; Success Metric 2; Success Metric 3; Success Metric 5
- `_techspec.md` → ADR header and the lesson form; Export list; Documentation part of the update; Release Plan; Interfaces; Invariants 1-9; Data Models; API Contract 10; API Contract 13; API Contract 17; API Contract 19; API Contract 20; Testing Approach; Build Order 2
- ADR-0438, ADR-0439, ADR-0254, ADR-0282, ADR-0412

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, roundfix, setup-context-driven, evidence-gate, systematic-debugging, no-workarounds, golang-safety.

Implemented this Task's documentation part beside history, outside the guidance plan and Setup Manifest. Its pending files participate in the Plan Digest and `fileChanges`; one confirmation applies guidance, history, then documentation. The result includes all documentation arrays and the prescribed text, state, refusal and failure forms. `--no-history` and `--no-skills` leave this part active.

The ADR functions read metadata, recognize the four notice forms outside code, and insert or replace the Amendment List line without YAML encoding. They preserve manually listed identifiers, numeric order, CRLF and every other byte. The caller filters self-notices and retired ADRs. The export planner is pure and handles creation, conversion, marker replacement, every specified refusal, losses and repository-owned bytes. The release classifier adds only `.secondbrain-export` and `docs/lessons/` to the maintenance boundary.

Shared facts seam for task_03: `gatherDocumentationFacts(root)` in `internal/cli/baseline_update_documentation.go` reads repository configuration without a user home and returns the Setup Manifest modules, local Spec Root and History Root, export facts (including immediate docs entries), and every ADR's path, retirement classification, header, notices and original bytes. ADR discovery includes nested directories. `speccheck.GlossaryFiles` supplies deterministic root/map/link resolution; the consistency check retains its original definition files and excludes the map documents from definitions.

### Acceptance evidence

1. Export planning and confirmed apply: the named command subtests `plan_names_the_list_and_the_losses`, `text_section`, `json_member`, `digest_moves_only_when_pending`, `plan_ready_message`, `file_changes`, `confirmed_apply`, `second_run_current`, `secondbrain_off`, `marker_refused_rest_applies`, `no_history_does_not_skip` and `apply_failure` passed in focused runs. They exercise real adopted temporary repositories, a temporary home, exact kept bytes, process/other losses, refusal with another header applied, populated and empty JSON arrays, digest sensitivity, and both skip flags. The failure seam reaches the real write boundary and returns category `documentation`, the prescribed restore action and exit 1 without depending on mode bits.
2. Amendment Lists: `fills_amendment_lists` and `fill_beside_history_unit` passed. The fixtures cover each notice form with one/two identifiers, linked identifiers, fenced/inline code, self-notices, retirement, nested ADRs, legacy headers, CRLF, exact other-byte preservation, a second current run, and a history unit appending lineage to the same successor. Pure reader/notice/writer cases also passed in focused subsets.
3. Release boundary: all three subtests of `TestMaintenanceOnlyCountsLessonsAndTheExportList` passed across focused subsets. Lessons and the export list stay inside the boundary; an additional product-code path crosses it.
4. This repository: the required Managed Refresh ran twice under a newly created temporary HOME with explicit existing Go caches and `GOPROXY=off`, `GOSUMDB=off`. Run 1 exited 0, reported 73 file changes, converted `.secondbrain-export` and filled 72 ADR headers. Run 2 exited 0 with `File changes: 0` and `Documentation: current`. A read-only text preview under that same fresh home exited 0 with `Baseline update: current`, `File changes: 0` and `Documentation: current`.

No Knowledge Mark was assigned. The preview reports 373 of 376 ADRs unmarked; marking them, the final repository-owned export lines and the three new notices remain task_05's slice.

### Focused checks

The following focused commands exited 0; they are implementation checks, not the Daemon's declared Verification:

- `go test ./internal/cli -run '^TestBaselineUpdateDocumentationPart/(plan_names_the_list_and_the_losses|confirmed_apply|fills_amendment_lists|fill_beside_history_unit|apply_failure)$' -count=1 -v`
- `go test ./internal/cli -run '^TestBaselineUpdateDocumentationPart/(text_section|json_member|digest_moves_only_when_pending|plan_ready_message|file_changes|second_run_current|secondbrain_off|marker_refused_rest_applies|no_history_does_not_skip|apply_failure)$' -count=1 -v`
- After the shared facts reader's final edit: `go test ./internal/cli -run '^TestBaselineUpdateDocumentationPart$/(plan_names_the_list_and_the_losses|fills_amendment_lists|fill_beside_history_unit|apply_failure)$' -count=1 -v`.
- After strengthening the nested-ADR and pending-digest fixtures: `go test ./internal/cli -run '^TestBaselineUpdateDocumentationPart$/(fills_amendment_lists|digest_moves_only_when_pending)$' -count=1 -v`.
- `go test ./internal/spec ./internal/baseline ./internal/releaseplan -run '^(TestReadADRHeader|TestADRAmendmentNotices|TestWithADRAmendedBy|TestPlanExportList|TestMaintenanceOnlyCountsLessonsAndTheExportList)/(mark_true|crlf|ascending_once|convert_reports_losses|apply_then_plan_unchanged|with_another_path)$' -count=1`
- The complementary pure-function subset also exited 0, covering every other required reader, notice, writer, planner and maintenance subtest. Two early selection probes reported no tests and were not counted as behavior evidence.
- `go test ./internal/speccheck -run '^TestGlossary(ReadsMappedAndRenamedGlossaries|DefinitionFormsAreRead)$' -count=1 -v`
- `go test ./internal/cli -run '^(TestUpdateCheckSummarizesTheBaselinePreview|TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary|TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths)$' -count=1 -v`
- Existing focused update/history cases `TestBaselineUpdatePlansPendingHistoryWithoutWrites`, `TestBaselineUpdateAppliesManifestPlanAndReportsJSON` and `TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest` exited 0.

Initial command-fixture failures were traced to passing preservation mode as a decision instead of the separate request field, and expecting a loss for a docs child the fixture had not created. Both fixture causes were corrected before the successful focused runs; no assertion was weakened.

### Generated artifacts and final refresh

- Inserted the prescribed baseline/release reference paragraphs, release guide/help sentence, and `### Documentation part` reference subsection.
- `make skills-sync` exited 0; it was run again after version recording. The three canonical/mirrored Roundfix files compare byte-identically.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` exited 0 and raised both version fields in both Roundfix copies from `0.1.278` to `0.1.279`; the record was generated, not hand-edited.
- After the last guide edit, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0. Its generated diff holds only the release-help fingerprint and the three edited user-guide fingerprints.
- Final implementation step: `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, repeated until the second run reported zero changes. The following read-only preview used the same command without `--yes` and reported current.
- A byte audit against the saved pre-refresh snapshot verified all 72 ADR edits change only the Amendment List line, preserving every other byte and `updated_at`. The export conversion removed only the reported `CONTEXT.md`, `docs/`, `README.md` and `AGENTS.md` includes; every kept byte matches the snapshot.
- The final path audit accounted for 96 changed files: the Task's declared files and the 72 ADRs written by the Managed Refresh. No other Task, `_tasks.md`, Baseline asset, managed guide, configuration, dependency or unrelated skill changed. The Daemon records the generated ADR paths.
- `git diff --check` exited 0. Refresh stdout/stderr, the ADR preimages and byte/path audit are under `/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task02-evidence-566e0orb/`.

Task status and checkboxes remain unchanged. Declared Verification, repository Verification and terminal settlement remain with the Daemon. No commit, push, Pull Request or cross-project message was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`
- `docs/adr/0034-run-worktrees-can-bootstrap-their-environment.md`
- `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`
- `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- `docs/adr/0092-triage-intent-lives-in-a-typed-backlog.md`
- `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`
- `docs/adr/0100-a-managed-refresh-proves-preservation-instead-of-backing-it-up.md`
- `docs/adr/0103-an-applied-refresh-republishes-the-setup-manifest.md`
- `docs/adr/0104-a-spec-accepts-on-evidence-it-did-not-author.md`
- `docs/adr/0120-retired-documentation-lives-under-one-history-root.md`
- `docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md`
- `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
- `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- `docs/adr/0168-a-related-adr-gap-opens-only-for-adrs-that-predate-the-spec.md`
- `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`
- `docs/adr/0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md`
- `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`
- `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`
- `docs/adr/0221-a-required-external-skill-is-held-to-its-setup-snapshot-not-only-its-lock.md`
- `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`
- `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`
- `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`
- `docs/adr/0231-the-jev-ceiling-is-a-user-config-value.md`
- `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- `docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md`
- `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- `docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md`
- `docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md`
- `docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md`
- `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
- `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`
- `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- `docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`
- `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
- `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`
- `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- `docs/adr/0377-only-an-adr-a-spec-applies-or-changes-draws-in-the-adrs-that-cite-it.md`
- `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`
- `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- `docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md`
- `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md`
- `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`

## Carry-forward provenance

- Source Run: `run_20261119T213428Z_8e0d266520442009`
- Source commit: `e2de5d6e92845b3f198e426a2c6ef1987b557716`
