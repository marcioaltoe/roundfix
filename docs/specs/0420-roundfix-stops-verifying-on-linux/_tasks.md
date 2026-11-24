---
schema: spec-tasks/v1
spec: 0420-roundfix-stops-verifying-on-linux
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Roundfix stops verifying on Linux

| id      | title                                                                                              | type    | complexity | needs                     |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Delivery Queue, its commands and doctor lose the Linux Lane, and its two keys become deprecated | backend | high       | —                         |
| task_02 | Pull requests, main and releases are verified on macOS, and Linux is compiled                       | infra   | medium     | —                         |
| task_03 | The user guide and the Roundfix skill describe a product without the Linux Lane                     | docs    | medium     | task_01, task_02          |
| task_04 | ADR-0452 supersedes the lane's records, and the glossary and export list follow                     | docs    | medium     | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                               | qa      | high       | task_04                   |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

The graph follows `_techspec.md` → Build Order one step per Task: steps 1 to 4 are task_01 to task_04 and step 5, the gate, is task_05. task_01 and task_02 share no declared file and run in the same Wave. task_03 needs both, because it describes what they made the product and the release do, and it shares `docs/user-guide/commands/deliver.md`, `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json` with task_01. task_04 needs the three, because its sweep reads what they leave and it shares `internal/verifyselect/contracts_repository_test.go` with task_03. No Task is joined to another: each changes a different contract (the product, the workflows and the compile contract, the guide and the skill, the decision records and the glossary) and each pays its verification cycle with tests of its own.

task_01 is `backend`: its outcome is the behavior of the Delivery Queue, its commands, the configuration loader and doctor. task_02 is `infra`: its outcome is where continuous integration and the release verify. task_03 and task_04 are `docs`: their outcomes are the guide, the skill, the decision records and the glossary, with their regenerated records.

## What the Tasks stage settled beyond the TechSpec

- **The guide in task_01.** task_01 changes `internal/cli` sources, and the Tasks stage requires such a Task to name a guide itself or through a Task it needs (write-tasks: "A Task naming a CLI surface names its skill or guide itself or through a Task it depends on"; `roundfix spec check` reports `SC-CLI-UNDOCUMENTED` otherwise). task_01 has no Task before it, so it removes from `docs/user-guide/commands/deliver.md` the section `Run the Linux Lane on demand`, which documents only the command task_01 deletes, and it re-records the page's fingerprint with the `-record-skill-coverage` step, as `SC-SURFACE-RECORD-UNTASKED` requires of a Task that edits a guide page. Every other edit of that page, of the other guide pages and of the skill stays in task_03, which runs the same step again for its own edits. This departs from `_techspec.md` → Build Order, which keeps every guide edit in step 3; the behavior, the contracts and the order are unchanged.
- **Test and subtest names.** The TechSpec names the tests; the Tasks name the subtests of the tests that have several cases, so the Verification and the gate check them by name. The 36 tests the coverage record may lose are listed by name in task_01 and in task_05.
- **The README's sentence.** task_02 writes it as "The test suite runs on macOS only; the Linux and Windows packages are compiled and not tested.", the shape `_techspec.md` → API Contract 8 gives it, so the Verification and the gate read one fixed sentence.
- **The coverage record in task_02.** task_02 adds untagged tests to `internal/verifyselect` and removes none from the record (`TestEveryTestCompilesForWindows` carries the `repocontract` tag and is not in it). `TestCoverageEquivalence` only logs a test the tree gained, so task_02 does not run the record command and does not declare the record; only task_01 writes it.
- **The ADR number.** 0452. The default branch ends at decision record 0451 (checked on 2026-11-24). When 0452 is taken before task_04 runs, the operator renames the file, its links in the six notices, in the two glossary entries and in the export line, and amends task_04 and task_05; invariant 11 keeps every test independent of the number.
- **The Governed Paths.** Measured on 2026-11-24 with `speccheck.GovernedPath`, in a `go test -overlay` probe whose test file lives outside the repository, over the 94 paths the five Tasks declare: the twelve paths of `_authorization.md` are the only governed ones, two declared by task_01 (`.roundfixrc.yml`, `docs/references/coverage-record.json`), four by task_02 (`.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `docs/agents/specific-repository.md`) and six by task_03 (the four references and the two `SKILL.md` files). task_04 declares none, and the record is unchanged.

## Shared files

The Roundfix skill's version fields, the owned-skill version record, the coverage record, the Behavior Surface Record and the number of `docscontract` files that `TestRepositoryContractTestsDeclareTheirRelevance` pins move with every Spec that touches them. Each Task regenerates the shared records on the tree it runs on and writes the next version and the next count, never a number copied from this plan. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Delivery

The authorization record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. No Task performs the release, edits `CHANGELOG.md`, a version file or the npm manifest, or sends a message; the operator classifies the `refactor` as a patch, cuts v0.153.3 after the merge, and tells the Praetor what the removal changes for the other repositories.

## Time in tests

No new test waits on elapsed time. The gate and doctor tests end on the command's return; the status and retry tests read a temporary Run Database and the owner fake's return; the workflow tests read files; the compile contract waits on `go vet`'s exit.
