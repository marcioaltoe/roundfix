---
schema: spec-tasks/v1
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
qa: task_06
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — The Baseline owns the documentation rules and the Secondbrain export

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Baseline states the documentation rules and the export rule                                    | docs    | high       | —       |
| task_02 | The update keeps the Standard Export List and each Amendment List, and its commit asks for no release | backend | high       | task_01 |
| task_03 | The Doctor Command reports what departs from the documentation rules                               | backend | high       | task_02 |
| task_04 | The skills ask the Lesson Learned question and the update marks the existing ADRs                  | docs    | medium     | task_03 |
| task_05 | This repository marks its decisions and publishes only its README, its user guide and its lessons   | docs    | high       | task_04 |
| task_06 | Run the final QA gate                                                                              | qa      | high       | task_05 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05 · 6 → task_06

The graph follows `_techspec.md` → Build Order, one Task per step. task_01 is step 1: the seven clause texts in three modules, the closed list of clause forces, the three tests that pin the changed text, the clause tests, the regeneration and the Managed Refresh, with the glossary entries the later Tasks use. task_02 is step 2, the documentation part of the update with the Release Plan rule that keeps the commit of an update from asking for a release. task_03 is step 3, the two Doctor checks. task_04 is step 4, the steps of the skills, the user-guide subsection and the contract test that pins every sentence the skills gained in this Spec. task_05 is step 5, this repository's marks, notices and export file. task_06 is the gate.

The Tasks run in series. task_01 changes no skill, so `make baseline-digests` runs before any skill changes, and no Task edits a module and a skill together. task_02, task_03 and task_04 each change both copies of the Roundfix skill's `SKILL.md` through the version record, the owned-skill version record and the Behavior Surface Record; task_02 and task_05 both write `.secondbrain-export` and ADR files; the Wave collision rule puts them in series, and the Build Order's dependencies give the same order. Each implementation Task runs the suites that consume its changed contract in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec is one context, what a repository's documentation holds and publishes, in five implementation Tasks the maintainer approved with the Idea ("one Spec, four or five Tasks and the QA gate, the fifth being this repository applying the rules to itself"). Each pays for its own cycle: seven clauses with 40 named guarantees, three renderings and the derived files; three pure functions, the update's documentation part and the Release Plan rule with 57 named cases; three departure readers and two Doctor checks with 43 named cases and three reworked ordering tests; three skill texts with 25 pinned guarantees; and about 375 decision records with a rewritten contract. No Task changes a command, a flag or an exit code, and no Surface Transcript is fixed (`_techspec.md` → Surface Transcripts).

## What the Tasks stage settled beyond the TechSpec

Each item was measured on 2026-11-19 on the planning clone at `894710c3`, after the plans of Specs 0409 and 0414 and the delivery of Spec 0411 had merged.

- **This repository's first refresh moves to task_02.** This repository has the Secondbrain turned on, its export file lists `docs/` whole, and 74 of its 372 ADRs carry an amendment notice in their body. Once task_02's code exists, its own update plans a conversion and the header fills, and `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` and `TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink`, which run the update on this repository and expect `current`, would fail at task_02's settlement. task_02 therefore runs the Managed Refresh as its last step, which converts the export file and fills those Amendment Lists, and task_05 writes the own lines, the marks and the notices and refreshes again. The current export contract test still passes after the conversion, because it reads only the history exclusion. The PRD records the choice under `## Decisions`, and the TechSpec's Build Order says it in steps 2 and 5.
- **Where the variation check lives.** `_techspec.md` → Interfaces puts `ExportVariationDepartures` in `export_list.go`, which task_02 creates. A Task cannot declare a file an earlier Task creates, so task_03 puts the function in `internal/baseline/export_variation.go` with its test, beside `docs_layout.go`, which holds the `DocumentationDeparture` it returns, and calls task_02's gathering of the facts without changing it. The functions and their contracts are those of the TechSpec.
- **ADR files are recorded, not declared.** A Task declares at most 50 paths, and the refreshes of task_02 and task_05 and the marks of task_05 change hundreds of files under `docs/adr/`. Each Task declares the files it edits by name, the three ADRs that gain a notice included, and the Daemon records the rest under `## Recorded paths`, which the QA scope audit counts; task_05's Verification proves every mark and that no ADR lost a line.
- **The Doctor Command is reached through its entry point.** A run of `roundfix doctor` through the binary proves adapters and reaches the forge, which no Verification or QA row may do. task_03's tests drive the command's entry point with the fakes the existing Doctor tests use, and the gate runs a throwaway test on a temporary copy of the audited head for this repository's two results. task_05's Verification proves Success Metric 6 through what the update reports for this repository: `current`, a current documentation part and no unmarked ADR.
- **Doctor tests that pin the report.** `internal/cli/doctor_test.go` pins three whole reports, the count of 17 lines and the ordered check names; `internal/cli/doctor_baseline_rules_test.go` and `internal/cli/doctor_claude_link_test.go` pin `residue:` right after `baseline-claude-link:`. task_03 declares and updates all three. None is a Governed Path.
- **Subtest names and phrases.** The TechSpec names the test files and the guarantees; each Task names one subtest per guarantee and the phrase it pins, from its API Contract, so each sabotage of the gate names the subtest that must fail. A phrase in a Verification command carries no backtick, so where an API Contract's sentence holds one, the command checks a part of the sentence without it.

## Shared files

Specs 0411 (delivered), 0414 and 0409 are delivered before this Spec, and Specs 0412 and 0413 after it (`_prd.md` → Prerequisites). Each Task edits only its own sentences, inserts them by the anchors the TechSpec names as they read on the default branch when the Task runs, and regenerates every shared record on the tree it runs on, with no version or digest written by hand. No Verification of this Spec depends on another Spec, and the graph needs no `requires`.

- With Spec 0409: `internal/baseline/module-versions.json`, the derived Baseline files under `internal/baseline/testdata/`, the Standard TypeScript profile pin, `docs/agents/setup-context.json` and the closed list in `internal/baseline/clause_characterization_test.go` (task_01); both copies of the Roundfix `SKILL.md` and `modes.md`, `skills/testdata/owned-skill-versions.json`, `CONTEXT.md`, `docs/user-guide/context-driven-development.md` and `docs/references/behavior-surfaces.json` (task_01, task_02, task_03 and task_04). Spec 0409 changes the module `spec-workflow`, which this Spec does not touch.
- With Spec 0414: `internal/releaseplan/classifier.go`, `internal/cli/cli.go` (the `release plan` help), both copies of the Roundfix `release` reference, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md` (task_02); the Roundfix `SKILL.md`, the owned-skill version record, `CONTEXT.md` and the Behavior Surface Record. The release sentence is inserted right after the sentence that ends "commits changing only those paths are maintenance-only.", as Spec 0414 leaves it.
- With Spec 0412, at its Idea stage: the Roundfix `modes` reference, `CONTEXT.md`, the owned-skill version record, the Behavior Surface Record, and, if it changes the module `spec-workflow`, the module version record, the derived Baseline files and the Setup Manifest.
- With Spec 0413, at its Idea stage: `docs/user-guide/commands/doctor.md` and `internal/cli/health.go` (task_03), `docs/references/coverage-record.json` (task_05), `CONTEXT.md`, the owned-skill version record and the Behavior Surface Record.
- For every Spec delivered after this one: an ADR it adds without the Knowledge Mark makes this repository's Doctor Command warn, which is the designed signal; an ADR it adds with `knowledge: true` must add its path to the Standard Export List, and an amendment notice it writes into an existing ADR must add the `amended_by` line, or this repository's update plans a change and its two update transcripts fail, until the Spec runs the Managed Refresh.

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads files or runs the update or the Doctor Command's entry point on a temporary repository.
