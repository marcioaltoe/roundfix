---
schema: spec-tasks/v1
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
requires: [0409-the-idea-and-prd-question-rules-that-block-the-goal]
qa: task_05
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
---

# Tasks — A project declares whether an open authoring question waits or follows the author

| id      | title                                                                                                   | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A project declares its Authoring Question Preference, and the status view reports it                    | backend | high       | —       |
| task_02 | A Delegated Decision is read from the PRD, checked and listed in the status view                        | backend | high       | task_01 |
| task_03 | The Spec workflow says what a session does with an Open Authoring Question                              | docs    | medium     | task_02 |
| task_04 | The authoring skills and the modes treat an Open Authoring Question alike, and the glossary names it    | docs    | high       | task_03 |
| task_05 | Run the final QA gate                                                                                   | qa      | high       | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order. task_01 is step 1: the Project Config key with its validation and warning, the reader at the default branch head, the `authoring_questions` field and line of the status view, and the guide and reference texts that describe the line. task_02 is step 2: the parser of the entry, the finding of the Clarification Record's check, the `delegated_decisions` field and block, the Surface Transcripts, the Clarification Record guide and the PRD template. task_03 and task_04 are step 3, cut in two: task_03 is the Baseline rule with its two clauses, the closed list of clause forces, the clause tests and the regeneration through the Managed Refresh; task_04 is the section of the three authoring skills, the `modes` reference, the glossary, the workflow guide and the amendment notice of ADR-0280. task_05 is the gate.

The Tasks run in series. task_01 and task_02 share `internal/cli/follow_up_status.go`, the three test files that pin the status text, `docs/user-guide/commands/follow-up.md`, the Roundfix skill's `follow-up` reference with its version fields, the owned-skill version record and the Behavior Surface Record; task_02 and task_04 share both copies of `write-prd/SKILL.md`; task_01, task_02 and task_04 share the Roundfix skill's version fields; the Wave collision rule puts them in series and the Build Order gives the same order. task_03 changes no skill, so `make baseline-digests` runs on a tree whose mirrors equal their canonical files, as `_techspec.md` → Invariant 14 requires. Each implementation Task runs the suites that consume its changed contract in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec is one context, what a session does with a question its author could not settle, in four implementation Tasks that each pay for their own cycle: a configuration key with eleven cases, a reader at the default branch head with seven, three moved status tests and two guides; a parser with sixteen cases, a check with four, a list with seven, three Surface Transcripts and a skill guide and template with nineteen pinned lines; a Baseline rule with 34 pinned guarantees, its rendering and the derived files; and three skills and a mode reference with 36 pinned guarantees in both copies, six glossary entries, a guide subsection and an ADR notice.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on 2026-11-19 on a scratch worktree of `abe2ebe1`, outside the planning clone and removed afterwards. The texts of `_techspec.md` → API Contracts 4 to 14 were applied, except the field-list sentences that follow the Go structs; the regeneration commands of `_authorization.md` ran in the order task_03 and task_04 give; and throwaway tests written to these Task files' test and subtest names and phrases passed: 39 pass lines over `./internal/baseline` and 57 over `./skills`. The whole `./internal/baseline`, `./skills` and `-tags docscontract ./internal/docscontract` suites then exited 0, the last with a one-line stand-in for the glossary entry of Rule Challenge, which Spec 0409 delivers. The change touched the 40 files `references/2026-11-19-prototype-and-governed-paths.md`, section 3, lists.

- **Build Order step 3 is cut in two.** The TechSpec puts the Baseline rule, the three skills, the `modes` reference, the glossary, the workflow guide and the ADR notice in one step. As one Task it would declare about 35 paths, run two regenerations and pin more than 70 guarantees, which no fresh session finishes in one sitting. task_03 takes the Baseline rule and task_04 the rest. The cut is a Delegated Decision of the `tasks` phase, recorded in `_prd.md` → Decisions, and `_techspec.md` → Build Order names it.
- **A later Task does not edit a file an earlier Task creates.** The Spec Consistency Check resolves a Task's declared paths at the Delivery Base, so task_02 cannot declare the files task_01 creates. task_02 puts the Delegated Decision list in `internal/cli/follow_up_status_delegated.go` with its test, beside the reader task_01 puts in `internal/cli/follow_up_status_authoring.go`, and the Clarification Record guide test in `skills/delegated_decision_record_test.go`; task_04 creates `skills/open_authoring_question_test.go` for the other two skill tests. The functions, the tests and their names are those of `_techspec.md` → Testing Approach. Each new path was measured not governed with the probe of `_authorization.md`.
- **The status text moves in two steps.** task_01 adds the `Open authoring questions` line and task_02 the `Delegated decisions` block, each to every test that pins the whole status text on the default branch when the Task runs. The Surface Transcripts hold both, so their test is task_02's.
- **The guide and reference texts are split by the Go struct they describe.** `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` requires each field sentence to name the fields of its struct in order, so task_01 writes the field sentences of `authoring_questions` and task_02 those of `delegated_decisions`.
- **ADR header fields.** Spec 0410 is delivered before this Spec and gives an amended ADR its `amended_by` list and every ADR its Knowledge Mark. task_04 writes the notice of API Contract 14 into ADR-0280 and lets the Managed Refresh fill the list. It gives ADR-0442 the mark `knowledge: true` when the docs-layout guide of the tree holds that rule and the ADR has no mark yet, because the decision is about what Roundfix does for the repositories that use it. The refresh can then also add ADR-0442 to the Standard Export List of `.secondbrain-export`, which is not a Governed Path.
- **`requires`.** The texts of task_03 and task_04 name the Rule Challenge, which Spec 0409 adds to the glossary, so the glossary contract task_04 runs passes only on a tree that holds Spec 0409.

## Shared files

Specs 0414, 0409 and 0410 are delivered before this Spec, and Spec 0413 after it (`_prd.md` → Prerequisites). Each Task inserts its text by the anchors the TechSpec names as they read on the default branch when the Task runs, edits only its own sentences, and regenerates every shared record on the tree it runs on, with no version or digest written by hand.

- With Spec 0414: `internal/cli/follow_up_status.go`, the status view's tests in `internal/cli/follow_up_status_queue_test.go`, `docs/user-guide/commands/follow-up.md`, `CONTEXT.md`, the Roundfix skill's version fields, `skills/testdata/owned-skill-versions.json` and `docs/references/behavior-surfaces.json`.
- With Spec 0409: `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/clause_characterization_test.go`, the derived Baseline files, `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`, both copies of `write-prd/SKILL.md`, `write-techspec/SKILL.md` and `prd-template.md`, the Roundfix skill's `modes` reference and version fields, the owned-skill version record, `CONTEXT.md`, `docs/user-guide/context-driven-development.md` and the Behavior Surface Record.
- With Spec 0410: `internal/baseline/module-versions.json`, the derived Baseline files and the Standard TypeScript profile pin, `docs/agents/setup-context.json`, the Roundfix skill's version fields and `modes` reference, the owned-skill version record, `CONTEXT.md`, `docs/user-guide/context-driven-development.md`, the Behavior Surface Record, `.secondbrain-export` and the header of ADR-0280.
- With Spec 0413: the Roundfix skill, the user guide and the Behavior Surface Record, which it changes after this Spec.

## Time in tests

No test of this Spec reads the clock or waits on a timer. The status tests set the time through the fixture's `followUpNow`, and every other new test reads files, parses bytes or builds a plan on a temporary repository.
