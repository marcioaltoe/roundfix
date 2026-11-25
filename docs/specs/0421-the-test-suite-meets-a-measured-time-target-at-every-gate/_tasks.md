---
schema: spec-tasks/v1
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Tasks — The test suite meets a measured time target at every gate

| id      | title                                                                                                                                          | type  | complexity | needs                              |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------- | ----- | ---------- | ---------------------------------- |
| task_01 | Code tests stop running the Go tool over other packages and building per test                                                                  | test  | high       | —                                  |
| task_02 | Tests share what they only read, and wait on events                                                                                            | test  | high       | task_01                            |
| task_03 | A Task settles on the light checks, the Complete Suite Run gets its home, and each gate reports its time                                       | infra | high       | task_01, task_02                   |
| task_04 | Continuous integration runs only for a pull request and in parallel jobs, a release proves its provenance, and two time rules are retired      | infra | high       | task_03                            |
| task_05 | Run the final QA gate                                                                                                                          | qa    | high       | task_01, task_02, task_03, task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The graph follows `_techspec.md` → Build Order one step per Task: steps 1 to 4 are task_01 to task_04 and step 5, the gate, is task_05. The edges are the Build Order's own. task_02 needs task_01 because both re-record the coverage record and task_02's timing rests on task_01's shared builds. task_03 needs both, because it turns the rules task_01 adds into contracts and its lab numbers need both repairs. task_04 needs task_03, because the workflow calls the targets and the selector flags task_03 adds and its contract count follows task_03's. The chain is serial: every implementation Task re-records `docs/references/coverage-record.json` and three of them change `docs/agents/specific-repository.md`, so no two of them can share a Wave.

task_01 and task_02 are `test`: their outcome is the test harness and the tests, with no product or gate behavior changed. task_03 and task_04 are `infra`: their outcome is what each gate of the repository runs, through the `Makefile`, the selector and the workflows; task_04 also carries the decision records, the glossary terms and the runbook that describe that shape, in the change that first depends on them.

## What the Tasks stage settled beyond the TechSpec

- **Step 4 stays one Task.** The TechSpec author asked whether Build Order 4 is still one Task now that it holds both workflows, the Release Provenance, ADR-0453 with six notices, the glossary and the runbook. Within four implementation Tasks and the gate it cannot be split without joining steps 1 and 2, which would make one larger Task out of two unrelated repairs. Every part of step 4 depends on the same fact, that continuous integration and the release stop running after a merge, and the decision records must change in the change that first depends on them (`docs/agents/spec-routing.md`, the Rule Challenge clause). The Tasks stage moved one part out instead: the glossary terms **Gate Time Target**, **Gate Time Report** and **Complete Suite Run** are written by task_03, which introduces those things, and task_04 writes **Gate Time Rule**, **Measured Dossier** and **Release Provenance**, revises **Full Contract Run** and adds the ADR-0453 link to all seven. Five implementation Tasks would be the alternative, and the request bounded the graph at four.
- **The renamed tests.** The TechSpec says two tests of the retired sets are renamed. Eleven tests of `internal/verifyselect` read the retired sets; task_03 renames each to a name that says what it proves under the plan and keeps four others under their names, so no name is lost without a destination. task_04 renames three tests whose names describe a push to main or the retired switch. The coverage record loses seventeen names in all, each listed with its destination in task_05, requirement 15.
- **The contract inputs.** The TechSpec leaves to the Tasks the Contract Relevance of the two nested runs that become contracts. task_01 declares the release matrix and the module files for the coverage collection, and the consumers and fixtures the journal replay reads; each also runs when a file of its own package changes, and every pull request runs both in the Full Contract Run.
- **Where the rule tests live.** Each rule file of `internal/testfixture` gains the `docscontract` tag and the `//verify:always` marker, except `external_tool_test.go`, whose two rule tests move to a new tagged file so that the two tests of `testfixture.RequireTool` stay code tests. `TestRepositoryContractTestsDeclareTheirRelevance` and `TestCLIShardsCoverEveryTestExactlyOnce` share one tagged file in `internal/verifyselect`, and the two Baseline-record tests share one in `internal/cli`.
- **`internal/testfixture` becomes a spawning test package.** The tests of `testfixture.ProductBinary` and `testfixture.FixtureBinary` start the Go tool, so task_01 installs the suite guard there, as the repository's rule on spawning test packages requires.
- **The One Build rule's inventory.** `internal/cli/baseline_release_gate_test.go` builds the product binary inside a test and is a Governed Path outside this Spec's grant, so it stays byte-identical and is the rule's frozen inventory entry with that reason.
- **The detach source and the runbook.** The repair of the detach wait may change `internal/cli/detach.go`, a source of the surfaces `command resolve` and `command watch`; task_02 declares it and records their Coverage Review and the Behavior Surface Record only if it changes. The release runbook task_04 changes is a source of `command release plan` and `guide docs/user-guide/release-runbook.md`; task_04 records both. `_prd.md` → Skills holds the four `unchanged` entries, and `_authorization.md` → Sanctioned regeneration names both Tasks for the conditional record command.
- **The coverage record's control case.** `TestCoverageEquivalence` fails on a lost name and only logs a gained one, so task_05 proves the control of Success Metric 2 with a test renamed without re-recording, the case the record refuses.
- **The timing proof of task_02.** The four waiting tests are proved by their reported time, under five seconds instead of about ten; this is the effect the Task exists for, and the threshold leaves room for a loaded machine.
- **The QA gate measures, and never moves a number.** The lab protocol's results are recorded against the approved targets; a miss is a finding that goes to the maintainer with its measurement (`_prd.md` → Risks, 2).

## Governed Paths

Measured on 2026-11-25, at `3f621ac6`, with `speccheck.GovernedPath` in a `go test -overlay` probe whose test file lives outside the repository, over the 107 paths the five Tasks declare under `interface:`, `creates:`, `deletes:` or `instruction:`: nine are governed, and they are the nine paths of `_authorization.md`. task_01 declares `internal/spec/coverage_test.go`, `internal/cli/baseline_plan_test.go`, `docs/agents/specific-repository.md` and `docs/references/coverage-record.json`; task_02 declares `docs/references/coverage-record.json`; task_03 declares `Makefile`, `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md` and `docs/references/coverage-record.json`; task_04 declares `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md` and `docs/references/coverage-record.json`. No Task declares another Governed Path.

## Shared files

The coverage record, the Behavior Surface Record, the Skill Coverage Map and the numbers of `docscontract` files that `TestRepositoryContractTestsDeclareTheirRelevance` pins move with every Spec that touches them. Each Task regenerates the records on the tree it runs on and raises a count by one from the value it finds, never a number copied from this plan. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Delivery

The authorization record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. No Task performs the release, edits `CHANGELOG.md`, a version file or the npm manifest, or sends a message. The delivering pull request is the first run of the new jobs; the operator records its job times in the reference document before the release and runs the release workflow's manual dispatch before pushing the tag. The Gate Time Rule is not a Task: it is the closing pull request of documents the maintainer decided, after the two windows close.

## Time in tests

No new test waits on elapsed time. The plan and shard tests read files and a temporary Git repository; the Makefile tests run the real recipes with a stub `go`; the gate-time tests run the tool in process on fixtures; the release-provenance tests run the step with a stub `gh` and end on its exit; the four repaired waiting tests end on their event.
