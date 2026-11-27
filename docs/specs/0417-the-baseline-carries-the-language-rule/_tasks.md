---
schema: spec-tasks/v1
spec: 0417-the-baseline-carries-the-language-rule
qa: task_04
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
---

# Tasks — The Baseline carries the language rule for code, data and products

| id      | title                                                                                                    | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Baseline carries the Language Rule and asks each repository its Reader Language                      | backend | high       | —       |
| task_02 | A codebase that keeps its language declares a Code Language Exception with its source                    | backend | high       | task_01 |
| task_03 | The review and the QA gate report what breaks the rule, and the skills, guides and glossary describe it  | backend | high       | task_02 |
| task_04 | Run the final QA gate                                                                                    | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The graph is `_techspec.md` → Build Order: steps 1, 2 and 3 are task_01, task_02 and task_03, and step 4, the gate, is task_04. The three implementation Tasks run in series because each re-records `docs/references/coverage-record.json`, task_01 and task_02 edit the same module, catalog, template, profiles and derived files, and task_03 describes the refusals task_02 builds. task_01 stays one Task because each of its parts breaks the tree without the next: the required decision fails 198 Baseline tests and every interactive-flow test until the answer fixtures carry it, the module edits fail the version and digest tests until the record commands run, and the regenerated catalog fails this repository's update transcript until the Managed Refresh runs; the ADR joins it so the record and the clauses it describes never disagree on the default branch.

## What the Tasks stage settled beyond the TechSpec

- **ADR number.** The next free number on the authoring base is 0454; task_01 creates `docs/adr/0454-code-is-written-in-english-and-each-repository-declares-its-reader-language.md`.
- **Test and subtest names.** The TechSpec names the tests of the Success Metrics; the Tasks add the names of the table tests (`TestACodeLanguageExceptionValueIsRefused`, `TestAClauseGateOnPresence`, `TestTheCatalogRefusesADefaultForACodeLanguageException`, `TestTheInteractiveFlowAsksNothingForTheCodeLanguageException`), the two documentation-contract tests of the review quotation, and every subtest, so that each Verification and the gate check them by name. No other top-level test is added, and no existing test is renamed, merged or removed, so the coverage record gains only the fourteen names the Tasks give and loses none.
- **The interactive-flow tests.** The 17 tests the TechSpec's prototype left without their scripted answer live in `internal/cli/baseline_human_test.go` and `internal/cli/baseline_release_gate_test.go`, both in the grant; task_01 owns them, because the required decision is what the interactive flow asks. task_02 changes no scripted answer, since the flow never asks for the exception.
- **Narrow Verification.** No Task runs a whole Go package; each runs its named tests and every test of the test files it changes, listed from those files, so a Task's settlement stays close to the Gate Time Rule's numbers. The gate runs the whole packages.
- **The guides' growth.** The gate measures the growth of the guides under `docs/agents/` against the Delivery Base, not against the 96,009 bytes of the authoring base, because another Spec may change them first.

## Governed Paths

Every governed file a Task declares was measured with `speccheck.GovernedPath` on 2026-11-27 and is among the 28 `_authorization.md` grants: task_01 declares 20 of them, task_02 15 and task_03 9, and together they declare all 28. No Task needs a governed file outside the grant, none of the eleven reserve files of `references/2026-11-27-characterization.md`, section 4, is declared, and no file is created under `internal/docscontract/`.

## Shared files

No other active Spec is on the default branch. A Spec that edits a Baseline module or an owned skill and merges first shares the module version record, the catalog files, the plan goldens, the Standard TypeScript profile's pin, this repository's Setup Manifest and the owned-skill version record; each Task regenerates them on the tree it runs on and writes "the next version", never a number. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads the embedded catalog, a mutated in-memory copy of it, a prompt string, an embedded skill or files in a temporary repository.
