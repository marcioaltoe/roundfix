---
task: task_01
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
status: completed
type: backend
complexity: medium
---

# Task 01: The status view shows a renamed Spec once

## Overview

The status view of the Follow-up Command lists a Spec whose slug was renamed twice: once under its new slug and once as the old slug still being authored, because the Phase Record keeps the old slug's entries and the row reader compares slugs. This Task makes the row reader drop a phase entry's slug when another slug of the same Spec Number holds the Spec, and states the rule in the `follow-up` guide, the glossary entry **Spec Queue Table**, ADR-0410 and the Skill Coverage Map. It answers the maintainer's Backlog Entry of 2026-11-18, [`references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md`](references/2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md), symptom 2.

## Requirements

1. MUST implement `_techspec.md` → API Contract 1 in `internal/cli/follow_up_queue.go` with the two functions of `_techspec.md` → Interfaces, under those names and signatures:
   - `followUpHeldSpecNumbers` holds, by Spec Number (read with the existing `specNumber` helper), every slug of `in.Active`, every folder of `in.Skipped`, every item of `in.Progress.queue`, merged ones included, and every name under `in.ArchiveRoot` that `spec.IsArchiveRecordName` accepts. A missing or unreadable Archive Root holds nothing and returns no error.
   - `followUpRenamedSlug` returns true when another slug of the candidate's Spec Number is held, or when another slug of that number in `latest` has a later entry, an equal time going to the greater slug. A slug without a Spec Number returns false.
   - In `followUpQueueRows`, a candidate whose folder is not a directory under the Spec Root keeps its row only when `latest` holds it and `followUpRenamedSlug` returns false. A candidate with a folder keeps today's behaviour.
2. MUST NOT change the status view's text or JSON shape, its sources, its forecasts, `followUpAuthoringPhases`, or any other row kind. The reader writes nothing.
3. MUST add `TestARenamedSpecLeavesNoAuthoringRow` to `internal/cli/follow_up_queue_test.go`, calling `t.Parallel()` first, built on `newFollowUpReadInput`, `followUpQueuePhase`, `followUpArchiveFixture`, `followUpWrite`, `mustMkdir` and `followUpQueueDiscover`. Every subtest records a phase entry of the phase `tasks` for `0401-old-slug` in the checkout two hours before the fixed time and asserts the exact ordered list of row slugs:
   - `renamed_to_a_planned_spec`: `0401-new-slug/_prd.md` exists; rows `0401-new-slug`;
   - `renamed_to_an_idea_folder`: only `0401-new-slug/_idea.md` exists; rows `0401-new-slug`;
   - `renamed_to_a_queue_item`: a queued Delivery Queue item `0401-new-slug`; rows `0401-new-slug`;
   - `renamed_to_a_merged_queue_item`: a merged item `0401-new-slug`; no row;
   - `renamed_to_an_archive_record`: an Archive Record for `0401-new-slug`; no row;
   - `renamed_while_authoring`: a later phase entry for `0401-new-slug`; rows `0401-new-slug`;
   - `another_number_stays`: a later phase entry for `0402-other`; rows `0401-old-slug`, `0402-other`;
   - `the_old_folder_stays`: the folder `0401-old-slug` exists and a later phase entry names `0401-new-slug`; rows `0401-old-slug`.

   `TestASpecIsBeingAuthored` and `TestTheQueueRowsComeInTableOrder` stay unchanged and keep passing; so does every other test that pins the status view on the default branch when this Task runs.
4. MUST append to the paragraph of `docs/user-guide/commands/follow-up.md` that begins "A Spec is being authored when" the sentence of API Contract 1, word for word.
5. MUST, through the `domain-modeling` skill, revise the glossary entry **Spec Queue Table** in `CONTEXT.md` with the sentence of API Contract 1, word for word, before its ADR links, and add the link `[ADR-0434](docs/adr/0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md)` after the ADR-0410 link. No other glossary entry changes.
6. MUST write the amendment notice of API Contract 1 into `docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md`, as a paragraph right after the paragraph that starts "**Specs being authored.**", word for word, and set its `updated_at` to the time of the change. No Considered Option of ADR-0410 is taken, so no option note changes.
7. MUST add the two `review` values of API Contract 1 to the entries `command follow-up` and `guide docs/user-guide/commands/follow-up.md` of `docs/references/skill-coverage.json`, word for word, after their `sources`, as the Coverage Review that `_prd.md` → Skills declares. No other entry changes.
8. MUST, as its last edit, after the last change to the guide, re-record the Behavior Surface Record with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, never by hand.
9. MUST NOT change any skill, the Roundfix follow-up reference, `docs/references/coverage-record.json`, any Baseline file, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Add the held Spec Numbers and the renamed-slug rule to the row reader.
- [ ] Add `TestARenamedSpecLeavesNoAuthoringRow` with its eight subtests.
- [ ] State the rule in the guide, the glossary entry and ADR-0410.
- [ ] Record the Coverage Review and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A phase entry of a renamed Spec's earlier slug makes no row when the new slug is planned, has a folder, is queued or merged, has an Archive Record or has a later entry.
- [ ] Another Spec Number keeps its row, and a folder under the Spec Root keeps its row.
- [ ] The guide, the glossary entry, ADR-0410 and the Skill Coverage Map hold the prescribed text, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md`
- instruction: `internal/cli/follow_up_history_test.go`
- instruction: `internal/cli/follow_up_status.go`
- interface: `internal/cli/follow_up_queue.go`
- interface: `internal/cli/follow_up_queue_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `CONTEXT.md`
- interface: `docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The rows of the Spec queue — consumers: `./internal/cli`
- The `follow-up` guide, the glossary, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestARenamedSpecLeavesNoAuthoringRow|TestASpecIsBeingAuthored|TestTheQueueRowsComeInTableOrder)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARenamedSpecLeavesNoAuthoringRow/renamed_to_a_planned_spec TestARenamedSpecLeavesNoAuthoringRow/renamed_to_an_idea_folder TestARenamedSpecLeavesNoAuthoringRow/renamed_to_a_queue_item TestARenamedSpecLeavesNoAuthoringRow/renamed_to_a_merged_queue_item TestARenamedSpecLeavesNoAuthoringRow/renamed_to_an_archive_record TestARenamedSpecLeavesNoAuthoringRow/renamed_while_authoring TestARenamedSpecLeavesNoAuthoringRow/another_number_stays TestARenamedSpecLeavesNoAuthoringRow/the_old_folder_stays TestASpecIsBeingAuthored TestTheQueueRowsComeInTableOrder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the test does not exist, so its pass lines are missing; afterwards every subtest and the two existing tests pass.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- "A Spec keeps its Spec Number when its slug is renamed, so an entry for a slug whose folder is not under the Spec Root does not count" || { echo 'follow-up guide lacks the sentence' >&2; exit 1; }; grep -F -- '**Spec Queue Table**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Spec Queue Table\*\*(:| —( |$))' || { echo 'glossary lacks the entry' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A renamed Spec keeps its Spec Number, so an entry for a slug without a folder there does not count" || { echo 'glossary entry lacks the sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0410-the-status-view-lists-specs-being-authored-and-open-backlog-entries.md | grep -qF -- '**Amended by [ADR-0434](0434-a-renamed-spec-keeps-its-spec-number-and-is-one-row-of-the-spec-queue.md) (2026-11-18).**' || { echo 'ADR-0410 lacks the notice' >&2; exit 1; }; grep -qF -- 'Spec 0407 task_01: a phase entry for a renamed Spec' docs/references/skill-coverage.json || { echo 'coverage map lacks the review' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task none of the texts exists; afterwards the guide, the glossary, ADR-0410 and the Skill Coverage Map hold them, the documentation contracts pass with the re-recorded Behavior Surface Record, and the whole CLI suite passes.

## References

- `_prd.md` → Goals (third bullet); Core Features 3 and 4; Success Metric 2; Skills; Glossary
- `_techspec.md` → Interfaces; API Contract 1; Coverage Map; Testing Approach; Build Order 1
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md), section 2
- ADR-0434, ADR-0410


## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, systematic-debugging, no-workarounds, domain-modeling, tech-writer, evidence-gate.

Implemented the Task 01 slice for Daemon Verification. Task status remains Daemon-owned.

- Acceptance criterion 1: `followUpHeldSpecNumbers` groups the active Specs, skipped folders, all Delivery Queue items (including merged items), and accepted Archive Record names by the existing `specNumber` helper. An unreadable or missing Archive Root contributes nothing. `followUpRenamedSlug` suppresses a folder-less phase candidate when another slug holds its number or has a later counted entry; equal times prefer the greater slug. The new `TestARenamedSpecLeavesNoAuthoringRow` exercises all eight prescribed scenarios with exact ordered row-slug assertions.
- Acceptance criterion 2: folder candidates retain their existing behavior. The `another_number_stays` and `the_old_folder_stays` subtests exercise the two controls. Each new scenario snapshots the temporary checkout and home before and after reading rows and asserts unchanged inputs. `followUpAuthoringPhases`, the status text and JSON types, sources, forecast code and other row kinds are unchanged.
- Acceptance criterion 3: the guide, only the **Spec Queue Table** glossary entry, ADR-0410 and the two Coverage Review entries contain the exact API Contract 1 text. ADR-0410 has the prescribed notice immediately after its authoring paragraph, with an updated timestamp and unchanged Considered Options. The sanctioned recorder generated `docs/references/behavior-surfaces.json`; its only changed surface is the follow-up guide digest. The same recorder is run again after this Result as the final edit, after all guide changes.

Focused implementation evidence:

- Before changing the reader: `go test -count=1 -timeout 5m ./internal/cli -run '^TestARenamedSpecLeavesNoAuthoringRow$' -v` exited 1: seven rename scenarios exposed the extra old or new authoring row, while `another_number_stays` passed.
- After changing the reader: `go test -count=1 -timeout 5m ./internal/cli -run '^(TestARenamedSpecLeavesNoAuthoringRow|TestASpecIsBeingAuthored|TestTheQueueRowsComeInTableOrder|TestQueueReadersSortEverySource|TestAnUnreadableSourceLeavesTheOtherRows)$' -v` exited 0. All eight new subtests and the existing focused controls passed.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestUserGuideLinksResolve|TestSkillCoveragePathsRequireOwnedFilesAndMatchingSources)$'` exited 0.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 and regenerated the Behavior Surface Record through its sanctioned writer.
- A read-only Python inspection compared the exact guide, glossary, amendment and review strings to API Contract 1; confirmed the amendment placement and ADR links; compared the remaining coverage map and glossary entries to HEAD; and confirmed all pre-existing tests, ADR options and the Daemon-owned status were preserved. All assertions passed. `git diff --check` exited 0.

The declared Verification commands, repository Verification and full CLI suite were not run; the Daemon owns them. No skill, Baseline, Task Graph, other Task, dependency or tooling configuration changed. No follow-up work identified.
