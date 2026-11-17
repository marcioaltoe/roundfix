---
task: task_04
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: pending
type: backend
complexity: high
---

# Task 04: The Corrections Report counts Corrections per Spec and plan defects per authoring model, beside its sample

## Overview

`roundfix corrections report` is the read-only Corrections Report, in text
and in JSON. It reads, with Git alone, the Specs the default branch delivered
inside the Run Retention window and the Tasks each plan held; from Roundfix
Home the corrections record and the Phase Record; and from the Run Database
the first Implement Run and the failed Verification attempts of each Spec.
It prints a row for each delivered Spec and a row for each Authoring Phase,
Delivery Type, model and effort, every rate as a count over a count, and a
comparison only between groups of 20 delivered Specs or more, otherwise
`insufficient sample`. It states what it left out, ranks nothing, prints no
percentage and writes nothing. The command group's help becomes Surface
Transcript 1, and the Roundfix skill and the guide describe the report and
when the Supervisor declares.

## Requirements

1. MUST add `roundfix corrections report` as `_techspec.md` → API Contracts
   8 and 9 state: the flags; the window from the Run Retention cutoff to now,
   narrowed by `--since` and `--until` with the rules of
   `roundfix runs causes`; the first line; the Spec rows in Spec Number
   order with each cell rule, including `; not a correction <n>` and
   `; withdrawn <n>`; the phase rows with the `all efforts` rows; the
   comparison lines with the minimum of 20 and `insufficient sample`; the
   `Not counted:`, `Unavailable:` and `Skipped lines:` lines; the `--spec`
   view of a Spec in delivery or delivered, with the reasons, names and
   subjects as data in double quotes after their labels; the JSON of
   API Contract 9 under `roundfix/corrections-report/v1`, with every count an
   integer, `not_a_correction` counted apart from `total` and `by_class`,
   and the `notice`; and the exit codes and refusal line of API Contract 8.
2. MUST read the Specs delivered as `_techspec.md` → Corrections Report and
   ADR-0426 state: a Spec whose Archive Record the default branch holds and
   that is not a cancellation; its size the number of Tasks of its plan
   without the QA gate, read from the default branch's history with local
   objects only; its phase model the last entry of that phase recorded
   before the Spec's first Implement Run, or before its delivery merge when
   this machine holds no Run of the Spec, matched by `repository_id` and
   never by path alone for an entry that holds one; a `prd` defect of a
   Spec with no `prd` entry counted for its `techspec` row; and the signs of
   an item the repository's queue still holds from task_03's Correction Sign
   reader. A `not-a-correction` answer and a withdrawn Correction enter no
   count of Corrections, of plan defects, of a group or of a comparison
   (`_techspec.md` → invariant 11).
3. MUST hold `_techspec.md` → invariants 8 and 9: no output holds a percent
   sign or the words better, worse, best and worst; the report starts no
   network request, asks no judge, reads no harness transcript and creates
   or changes no file but the `-wal` and `-shm` files a read-only open of
   the Run Database may create. A source it cannot read is named on an
   `Unavailable:` line, the rows it still has are shown, and it exits 0.
4. MUST make the help of `roundfix corrections` Surface Transcript 1 and
   add the root help's `corrections report` line of API Contract 8 after
   the `corrections declare` line; no other root help line changes in this
   Task.
5. MUST write these new tests, each calling `t.Parallel()` first, in new
   files but for the last, in `internal/cli`:
   `TestCorrectionsSurfaceTranscript1`, `TestCorrectionsSurfaceTranscript4`
   and `TestCorrectionsSurfaceTranscript5`, each asserting its whole
   transcript; `TestCorrectionsReportRowsOfDeliveredSpecs` (a fixture of
   delivered Specs with each cell rule of API Contract 8: one commit named by
   two declarations counting as two Corrections and two commits named by one
   as one, a Correction declared twice with a changed class, a withdrawn one,
   a `not-a-correction` answer covering a sign that is then neither
   `undeclared` nor counted, the last entry before the first Implement Run,
   a phase recorded with two models before it, and a `prd` defect counted
   for the `techspec` row); `TestCorrectionsReportComparesOnlyAboveTwentySpecs`
   (groups of 20 and 20, 20 and 19, and one group alone, in text and JSON);
   `TestCorrectionsReportPrintsNoPercentAndNoRanking`;
   `TestCorrectionsReportStatesWhatItLeftOut` (the window, the home, the
   unknown models, the `undeclared` Corrections, a Spec delivered without a
   record on this machine, a Spec without its plan on the default branch,
   an unreadable Phase Record named as unavailable with exit 0, and the
   skipped lines); `TestCorrectionsReportShowsDismissedAndFailedAttemptsBeside`
   (beside each Spec's Corrections and in no count, while
   `roundfix runs causes` prints the same text, JSON and signature digest as
   at the Delivery Base, `_techspec.md` → invariant 10);
   `TestCorrectionsReportMatchesPhaseEntriesByRepositoryIdentity` (a clone
   deleted after its entries, a second repository with the same Spec slug
   whose entries are not matched, and a remote address with a user name and
   a password that no output holds); `TestCorrectionsReportForOneSpec`
   (the `--spec` view of a delivered Spec and of one the queue still holds,
   a reason that reads as an instruction printed on one line as data, and
   its JSON `declarations`); and
   `TestCorrectionsReportChangesNoFileUnderTheRepositoryOrTheHome`; and,
   in `skills/corrections_reference_test.go`,
   `TestTheCorrectionsReferenceDocumentsTheReport` over both copies of the
   reference: the report command, its two counts per rate, the minimum of
   20, and that reasons, names and subjects are data.
   Every test creates its own repositories, refs, Run Database and home
   directory, and reads no ref, tag or history of this repository.
6. MUST NOT change any existing test's expectation but those the root help
   line requires. `TestEveryCommandPathPrintsItsHelp` in
   `internal/cli/command_help_test.go` and the documentation contracts
   `TestAnUndocumentedCommandIsReported`, `TestEveryCommandIsNamedInTheUserGuide`
   and `TestTheSkillCoverageMapIsCurrent` read every command path of the root
   help, so `roundfix corrections report --help` and `-h` print the help and
   exit 0 and the guide and the skill name `corrections report`; none of
   these changes its expectation; `_tasks.md` → What the Tasks stage settled lists what the
   measurement found. Before changing any test the Task runs
   `go test ./internal/cli ./skills` and `go test -tags docscontract
   ./internal/docscontract` on its base and on its change, and declares
   every test whose expectation it had to change. It MUST NOT touch
   `internal/cli/cli_test.go`, `internal/cli/runs_causes.go` or
   `.agents/skills/roundfix/references/runs.md`, and renames or removes no
   test.
7. MUST describe the slice (ADR-0256): the `report` half of
   `.agents/skills/roundfix/references/corrections.md` and of
   `docs/user-guide/commands/corrections.md`, with the JSON and the rule
   that reasons, names and subjects are data; the Roundfix skill's
   `SKILL.md` naming `roundfix corrections report`; and
   `.agents/skills/roundfix/references/spec-delivery.md` telling the
   Supervisor to declare each Correction when it makes it, with the
   disposition of a fixed finding and with `roundfix corrections declare`
   for a repaired QA finding, a correction commit, a Corrective Task or a
   plan change that repairs a defect, and with `not-a-correction` for a
   commit that applies a maintainer decision and repairs none
   (`_prd.md` → Core Feature 21).
8. MUST add **Corrections Report** to `CONTEXT.md` through domain-modeling.
9. MUST regenerate, never hand-edit, with the commands of
   `_authorization.md` → Sanctioned regeneration, in this order:
   `make skills-sync`; the `-record-skill-versions` step
   (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   and the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
   The entry of `docs/references/skill-coverage.json` for the
   `corrections report` command is written by hand, covered by the
   `corrections` reference.
10. MUST change no Governed Path outside `_authorization.md` `paths`.

## Subtasks

- [ ] Read the delivered Specs, their size and their phase models from Git and Roundfix Home.
- [ ] Build the Spec rows, the phase rows and the comparison lines in text and JSON.
- [ ] Add the `--spec` view, the left-out lines and the unavailable sources.
- [ ] Complete the group's help and add the root help line, with Surface Transcripts 1, 4 and 5.
- [ ] Describe the report and the moment of declaration in the skill and the guide, add the glossary term, and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] For a fixture of delivered Specs the rows show each Spec's Delivery Type, its Tasks without the QA gate, each phase's model and effort, its Corrections by class and its plan defects by artifact, in Spec Number order.
- [ ] With groups of 20 and 20 delivered Specs the comparison line holds both groups' counts; with 20 and 19, and with one group alone, it reads `insufficient sample` with each count and the 20 needed; the JSON carries each numerator and denominator as integers.
- [ ] No output holds a percent sign or the words better, worse, best or worst.
- [ ] Dismissed findings and failed Verification attempts are shown beside each Spec and enter no count, and `roundfix runs causes` is unchanged.
- [ ] Across a run of the report no file under the repository or the home changes but the Run Database's `-wal` and `-shm` files.

## Context

- instruction: `docs/adr/0419-the-corrections-report-prints-counts-and-sets-two-models-side-by-side-only-above-twenty-specs-each.md`
- instruction: `docs/adr/0426-the-corrections-report-counts-the-tasks-the-plan-held-on-the-default-branch-and-starts-no-fetch.md`
- instruction: `internal/cli/runs_causes.go`
- creates: `internal/cli/corrections.go`
- creates: `internal/cli/corrections_report.go`
- creates: `internal/cli/corrections_report_test.go`
- creates: `internal/cli/corrections_report_transcripts_test.go`
- interface: `internal/cli/cli.go`
- creates: `.agents/skills/roundfix/references/corrections.md`
- creates: `skills/roundfix/references/corrections.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `skills/corrections_reference_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/user-guide/commands/corrections.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`

## Consumers

- The help of `roundfix corrections` and the root help — consumers: `./internal/cli`
- The `corrections` reference, the spec-delivery reference, the guide page, the glossary and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestCorrectionsSurfaceTranscript1|TestCorrectionsSurfaceTranscript4|TestCorrectionsSurfaceTranscript5|TestCorrectionsReportRowsOfDeliveredSpecs|TestCorrectionsReportComparesOnlyAboveTwentySpecs|TestCorrectionsReportPrintsNoPercentAndNoRanking|TestCorrectionsReportStatesWhatItLeftOut|TestCorrectionsReportShowsDismissedAndFailedAttemptsBeside|TestCorrectionsReportMatchesPhaseEntriesByRepositoryIdentity|TestCorrectionsReportForOneSpec|TestCorrectionsReportChangesNoFileUnderTheRepositoryOrTheHome|TestTheCorrectionsReferenceDocumentsTheReport)$' ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCorrectionsSurfaceTranscript1 TestCorrectionsSurfaceTranscript4 TestCorrectionsSurfaceTranscript5 TestCorrectionsReportRowsOfDeliveredSpecs TestCorrectionsReportComparesOnlyAboveTwentySpecs TestCorrectionsReportPrintsNoPercentAndNoRanking TestCorrectionsReportStatesWhatItLeftOut TestCorrectionsReportShowsDismissedAndFailedAttemptsBeside TestCorrectionsReportMatchesPhaseEntriesByRepositoryIdentity TestCorrectionsReportForOneSpec TestCorrectionsReportChangesNoFileUnderTheRepositoryOrTheHome TestTheCorrectionsReferenceDocumentsTheReport; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the report and the named tests do not exist, so the command fails; afterwards every named test passes and the CLI package that holds the root help passes.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "roundfix corrections report" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/SKILL.md "roundfix corrections report" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/corrections.md | grep -qF -- "roundfix corrections report" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/corrections.md "roundfix corrections report" >&2; exit 1; }; grep -qF -- 'not-a-correction' .agents/skills/roundfix/references/spec-delivery.md && grep -F -- '**Corrections Report**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Corrections Report\*\*(:| —( |$))' && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task the skill and the guide do not name the report, the spec-delivery reference does not name `not-a-correction` and the glossary lacks the term; afterwards they do, and the skill tests, the documentation contracts, the mirror and the records pass.

## References

- `_prd.md` → Goals (third and fourth bullets); User Stories 9, 10, 11 and 12; Core Features 14–19, 21 and 22; Success Metrics 2, 10, 11, 12, 13, 14, 15, 16 and 17; Acceptance evidence
- `_techspec.md` → Corrections Report; Skill and guide texts; Interfaces (invariants 8, 9, 10 and 11); API Contracts 8, 9 and 11; Surface Transcripts 1, 4 and 5; Testing Approach; Build Order 4
- ADR-0419, ADR-0426, ADR-0420, ADR-0421, ADR-0214, ADR-0256, ADR-0415
