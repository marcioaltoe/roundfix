---
schema: spec-tasks/v1
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
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
      needs: [task_06]
    - id: task_06
      file: task_06.md
      needs: [task_04]
---

# Tasks — A session runs the flow's own commands without a permission prompt

| id      | title                                                                                                                      | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The Flow Command List is closed, rendered for each harness, and equal to what the owned skills name                        | backend | high       | —       |
| task_02 | A harness's file is written or cleaned only as a preview with the same digest showed it                                    | backend | high       | task_01 |
| task_03 | `roundfix permissions` shows the list for a harness and writes it only with the digest a preview printed                   | backend | high       | task_02 |
| task_04 | The Doctor Command and the Update Check say where a project's list stands, and a harness's permission files are Governed Paths | backend | high       | task_03 |
| task_05 | Run the final QA gate                                                                                                      | qa      | high       | task_06 |
| task_06 | Every text preview names every Named Exclusion, `git -C` included                                                          | backend | low        | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05

task_06 is corrective Task 1 of 2, added on 2026-11-20 for finding F-01 of the QA gate of that day (`qa/qa-report-2026-11-20.md`, audited head `df1529f9`, row Q15); the QA gate now needs it.

The graph follows `_techspec.md` → Build Order, with its step 1 cut in two. task_01 is the list half of step 1: the data, the three renderings, the Codex golden, the release record and the three contracts of API Contract 6. task_02 is the file half: the settings editor, the OpenCode reader, the plan, its digest and the write. task_03 is step 2, the command with its skill rule, its guide and the glossary. task_04 is step 3, the two reports, the governed set and the test that nothing else uses the list. task_05 is the gate.

Sized by the 0391 criterion, the Spec is one context, the list a session runs without a prompt and the files that hold it, in four implementation Tasks. Each pays for its own cycle: task_01 derives the list again on its tree and proves it against every owned skill and every rendering with 38 named cases; task_02 is the byte-preserving editor and the digest-bound write, the part the TechSpec names as most likely to be wrong, with 51 named cases; task_03 is a new command with four Surface Transcripts and 25 named cases and the skill, guide and glossary that the documentation contracts force into the same Task; task_04 changes two existing reports and their pinned tests, the governed set and its two contracts. The TechSpec's step 1 is cut because one session would otherwise carry both the list's derivation and the editor, two verifiable outcomes with no shared file; cutting it changes no edge the Build Order states.

The Tasks run in series. task_01 and task_02 share the package `internal/flowpermissions`, and task_02 plans from the list value task_01 defines; task_03 and task_04 both change `internal/cli/cli.go`, both copies of the Roundfix `SKILL.md` and `setup` reference, the owned-skill version record, the Skill Coverage Map, the Behavior Surface Record and `CONTEXT.md`, so the Wave collision rule puts them in series, and the Build Order's dependencies give the same order. Each implementation Task runs the suites that consume its changed contract in its own Verification, so no separate Integration Checkpoint is needed.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on 2026-11-19 on the planning clone at `abe2ebe1`, after the plans of Specs 0409, 0410 and 0414 merged.

- **The list is derived again by its Task.** The TechSpec's table is the authoring base. Specs 0414, 0409, 0410 and 0412 deliver first and may make a skill name another command; task_01 classifies every command the documentation contract reports on its own tree by the two rules of ADR-0445 and records each difference in its Result. A scan of the twelve owned skills at `abe2ebe1` found exactly the TechSpec's 35 registered `roundfix` subcommands, 21 Git subcommands and the GitHub commands `gh api`, `gh pr view` and `gh auth` with three second words, plus `gh failed` and `gh has` from quoted diagnostics. Exclusion `gh auth` covers `gh auth login`, `gh auth refresh` and `gh auth switch` by prefix, and any other GitHub command that writes takes the reason `it changes GitHub`.
- **The retired set has a record from the first release.** The TechSpec records "what earlier releases held" as empty, which would leave the test unable to fail until a later release. task_01 records instead every entry and rule this release ships, in `internal/flowpermissions/testdata/released.txt`: a later Spec that drops an entry without moving it to the retired set fails `TestRetiredSetKeepsEveryReleasedEntry` at once, and one that adds an entry appends it. Invariant 9 holds either way.
- **The plan reads one list value.** A Task cannot declare a file an earlier Task creates, so task_01 keeps the list in one unexported value that every exported function reads, and task_02 plans from that value without changing `list.go`; its tests plan with a list that holds a retired rule.
- **Harness Permission File is defined with the command.** The TechSpec's Build Order defines it with the governed set, but the command's guide page in task_03 names it, and the glossary contract reads every term a guide uses. task_03 defines all four terms; task_04 adds to the definition the sentence that each one is a Governed Path, when the pattern exists, and revises Governed Path.
- **Strings the gate requires and the TechSpec did not fix.** The link refusal `<path> leaves the repository through a link. Nothing was written.` with exit `1`, the Doctor line outside a Git repository, the help sentences of the Doctor and the Update Check, the guide sentences and the configuration section are prescribed verbatim by the Task that writes them, so the gate checks text a Task was told to write.
- **The two reports run through their entry points.** A run of `roundfix doctor` through the binary proves adapters and reaches the forge, and `roundfix update check` reaches the release listing; no Verification or QA row may do either. task_04's tests drive both with the existing fakes in temporary repositories, and the gate runs those tests.
- **What the gate measures offline.** The short-option Guard Entry form and the Codex file name were left unmeasured. The gate measures what needs no harness: the rendered short-option rules under the documented reading of `*`, the Codex file's place and form against the folder and extension the documentation names and the form the author measured with Codex's offline checker. It runs no real `claude`, `codex` or `opencode`, which the authorization keeps out of every test and Verification, and the two walks of `_prd.md` → Unreachable Acceptance stay the proof in a real session. The PRD records this as a delegated decision of this stage.
- **Doctor tests that pin the report.** At `abe2ebe1`, `internal/cli/doctor_test.go` pins the ordered check names and whole reports, and `internal/cli/doctor_baseline_rules_test.go` and `internal/cli/doctor_claude_link_test.go` pin `residue:` right after `baseline-claude-link:`; Spec 0410 adds `baseline-export` and `baseline-docs` before `residue` and its own ordering test. task_04 declares the three that exist here and updates every such test on the default branch when it runs; the Daemon records the one Spec 0410 adds.
- **Paths a Task may not touch.** `internal/cli/cli_test.go` and `internal/gittest/gittest.go` are historically bounded Governed Paths this Spec is not granted. `TestRunHelp` checks substrings only, so the new command needs no change there; `TestRootHelpNamesThePermissionsCommand` lives in `internal/cli/permissions_test.go`. Test fixtures never carry a Harness Permission File's path outside a temporary directory, because from task_04 on such a path is governed.

## Shared files

Specs 0414, 0409, 0410 and 0412 are delivered before this Spec (`_prd.md` → Prerequisites). Each Task edits only its own sentences, inserts them by the anchors it names as they read on the default branch when it runs, and regenerates every shared record on the tree it runs on, with no version or digest written by hand. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

- With every one of the four: both copies of the Roundfix `SKILL.md`, `skills/testdata/owned-skill-versions.json` and `docs/references/behavior-surfaces.json` (task_03 and task_04).
- With Spec 0410: `internal/cli/doctor.go`, `internal/cli/health.go`, `internal/cli/doctor_test.go`, `internal/cli/doctor_baseline_rules_test.go`, `internal/cli/doctor_claude_link_test.go`, both copies of the `setup` reference and `docs/user-guide/commands/doctor.md` (task_04, and the `setup` reference in task_03); `internal/cli/cli.go`, `docs/user-guide/usage.md`, both copies of the `modes` reference and `CONTEXT.md` (task_03 and task_04). This Spec's check comes after Spec 0410's two checks and before `residue`. ADR-0444 and ADR-0445 carry `knowledge: true` and ADR-0446 `knowledge: false`; this plan reaches the default branch before Spec 0410 delivers, so that Spec's own marking pass names them in the Standard Export List. If this plan merged after Spec 0410's delivery, the update of this repository would plan a change for them until the operator ran the Managed Refresh.
- With Spec 0414: `internal/cli/cli.go`, `docs/user-guide/usage.md`, `docs/references/skill-coverage.json` and `CONTEXT.md` (task_03 and task_04).
- With Spec 0409: both copies of the `modes` reference and `CONTEXT.md` (task_03 and task_04), and the Skill Coverage Map if it changes it.
- With Spec 0412: `docs/user-guide/configuration.md`, both copies of the `modes` reference and `CONTEXT.md` (task_04, and `CONTEXT.md` in task_03). Spec 0412 adds a Project Config key; this Spec adds none and changes no harness permission file.
- `internal/speccheck/governed.go`, `internal/speccheck/governed_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/cli/update_check.go`, `internal/cli/update_check_test.go`, `docs/user-guide/commands/update.md` and `internal/docscontract/command_documentation_test.go` are changed by none of the four.

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads or writes files in a temporary directory, renders the list or runs a command's entry point in a temporary repository.
