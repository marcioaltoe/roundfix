---
task: task_06
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: completed
type: qa
complexity: high
---

# Task 06: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec rewords four clauses of the Baseline's docs-layout rule and one of its cross-project rule and adds the lesson rule and the export rule; gives `roundfix baseline update` a documentation part that keeps the Standard Export List and fills each Amendment List; gives `roundfix doctor` the checks `baseline-export` and `baseline-docs`; counts `.secondbrain-export` and `docs/lessons/` as needing no release; carries the lesson steps into `archive-spec` and the Roundfix skill's `corrections` reference and the marking of existing ADRs into its `modes` reference; defines five terms and revises Managed Refresh in the glossary; and applies the rules to this repository's export file and decision records, with ADR-0438, ADR-0439 and ADR-0440. No command, flag, configuration key or exit code is added, changed or removed.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon or a Delivery Owner;
- runs `roundfix doctor` through the built binary, because its other checks prove adapters and reach the forge; the Doctor rows run its entry point through the named tests and a throwaway test on a temporary copy;
- calls ai-memory, runs the real `gh` or `rsync`, or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`; their tests make their own temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that needs the host's toolchains. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary copy is made with `git archive` of the commit it names and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 9 by running, with `-v` over `./internal/baseline`, `TestTheBaselineModulesCarryTheDocumentationRules`, `TestTheGuidesRenderTheDocumentationRulesOnce`, `TestTheCaptureTriggerAsksForAConditionalLesson`, `TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide`, `TestTheDocsLayoutRegionStaysLean`, `TestADRLifecycleContract`, `TestBaselineClauseForceIsCharacterized`, `TestShippedGuidanceCitesNoRepositoryRecord` and `TestNoTwoBaselineClausesShareText`, recording every subtest by name; then the whole `./internal/baseline` package and `-run 'Baseline|GuidanceComposition|SetupManifest' ./internal/cli`, recording each pass.
3. MUST reproduce the Managed Refresh of Success Metric 1 through `bin/roundfix`, in a fresh home prepared as above, on a temporary copy of the Delivery Base (the merge base of the audited head with the default branch). Record that `bin/roundfix baseline update --repo . --no-skills --no-history --yes --format text` there writes, each exactly once after whitespace normalization, "and a Lesson Learned once an accepted ADR, a guide or a skill states what it teaches are deleted" and "Write a Lesson Learned only when every part of this test holds" into `docs/agents/docs-layout.md`, "send no message about what it taught" into `docs/agents/agent-instructions.md`, and "An ADR without the mark is not published." into `docs/agents/secondbrain.md`, none of which the copy held before; and that a second run of the same command reports nothing to do. Record also that the same command on a temporary copy of the audited head reports nothing to do on its first run.
4. MUST reproduce Success Metrics 2 and 3 through `bin/roundfix`, in a fresh home prepared as above, on a temporary copy of the audited head:
   - with `.secondbrain-export` replaced by a file that lists `docs/` whole, `AGENTS.md`, one file of its own with a comment above it, and with the `knowledge` line removed from one ADR, `bin/roundfix baseline update --repo . --no-skills --no-history --format text` prints the section `Documentation: pending`, the line `Export list: update .secondbrain-export (...)`, a `- stops exporting <path> (process)` line for `AGENTS.md`, `docs/agents/`, `docs/backlog/`, `docs/findings/`, `docs/handoffs/`, `docs/history/` and `docs/specs/` each, and `ADRs without a Knowledge Mark: 1 of <m>; they are not exported`; record the whole section;
   - the same command with `--yes` writes the block, removes the lines it named and keeps the file's own line and its comment byte for byte; a second run reports nothing to do;
   - with the line `# roundfix:standard-export:begin` then deleted, a run prints `Export list refused: # roundfix:standard-export:begin is missing; nothing is written to .secondbrain-export` and leaves the file byte-identical;
   - with a notice `Narrowed by ADR-0438 (2026-11-20).` added after the decision text of one accepted ADR other than ADR-0438 and the same notice added inside a fenced code block of another, a run lists only the first under `ADR headers:`, the confirmed run adds only its `amended_by` line, every other byte of both files is unchanged, and a second run reports nothing to do.
5. MUST verify Success Metrics 2 to 5 through the named tests, with `-v`, recording every subtest by name: `TestReadADRHeader`, `TestADRAmendmentNotices`, `TestWithADRAmendedBy` and `TestLessonLearnedProblems` over `./internal/spec`; `TestPlanExportList`, `TestExportVariationDepartures` and `TestDocsLayoutDepartures` over `./internal/baseline`; `TestBaselineUpdateDocumentationPart` and `TestDoctorReportsTheDocumentationRules` over `./internal/cli`; and `TestMaintenanceOnlyCountsLessonsAndTheExportList` over `./internal/releaseplan`; then `./internal/spec`, `./internal/releaseplan` and `-run 'Baseline|Doctor|UpdateCheck|ReleasePlan' ./internal/cli`, recording each pass. Record that the Doctor rows run the command's entry point and that `exit_unchanged` holds.
6. MUST verify Success Metric 6 on a temporary copy of the audited head: add there one throwaway test in `internal/cli` that calls the function task_03 added for the two Doctor results with the copy's root, run it, and record that both results are `ok` with the details of API Contracts 11 and 12; record also, on the audited head itself, that every file under `docs/adr/` carries `knowledge: true` or `knowledge: false` on the line after `status`, that `bin/roundfix baseline update --repo . --no-skills --no-history --format json` in a fresh home reports `current` with a current documentation part and `unmarked` 0, that ADR-0248, ADR-0278 and ADR-0412 carry the notices of API Contract 23 and list ADR-0438 or ADR-0440 in `amended_by`, and that `.secondbrain-export` is the Standard Export List followed by the own lines of API Contract 22, each with its reason. The throwaway test never enters the checkout.
7. MUST verify Success Metric 7 by running, with `-v`, `TestTheSkillsCarryTheDocumentationSteps` over `./skills`, recording every subtest by name; then `TestEveryOwnedSkillVersionIsRecorded`, the whole `./skills` package and `make skills-sync-check`, recording each pass. Record that `skills/testdata/owned-skill-versions.json` holds the version each `SKILL.md` of `archive-spec` and `roundfix` declares, in both copies, and that the `## --update` section of `modes.md` keeps five steps and every sentence it held at the Delivery Base.
8. MUST verify Success Metric 8 by running `-tags docscontract ./internal/docscontract` whole and, with `-v`, `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryAvoidLinesArePlainLists`, `TestGlossaryLinksNameAcceptedADRs`, `TestTheSkillCoverageMapIsCurrent`, `TestThisRepositoryExportsNoProcessPath` and `TestTheExportRuleRefusesAProcessPath`, then `-run '^TestCoverageEquivalence$' ./internal/spec` and `make verify-docs`, recording each pass.
9. MUST verify Success Metric 9 against the Delivery Base, never against a named commit:
   - every clause of every module under `internal/baseline/assets/modules/` at the Delivery Base exists at the audited head with the same id and `enforcement`, and with the same `guidance` except the five clauses of API Contracts 1 to 4 and 6; the only clauses the head adds are the five of `rule.context.lessons` and `rule.secondbrain.knowledge-export`;
   - every file under `internal/baseline/assets/source-baselines/` and `internal/baseline/assets/setups/`, every module other than `context-workflow.json`, `core.json` and `secondbrain.json`, and every profile other than `standard-typescript-monorepo.json` is byte-identical to the Delivery Base, and no profile lists either new rule under `requiredRules`;
   - the only lines `internal/baseline/clause_characterization_test.go` gains over the Delivery Base are the five ids of the new clauses, and it loses none; `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_human_test.go` and `skills/baseline_skill_contract_test.go` are byte-identical to it.
10. MUST run each sabotage below on its own temporary copy of the audited head. Each sabotage keeps the package compiling and regenerates or synchronizes nothing, so the named test fails on its assertion; record that it fails on the named subtest.
    1. "Otherwise write nothing." is removed from the `guidance` of `clause.context.lesson-02-test` in the module file: `TestTheBaselineModulesCarryTheDocumentationRules/lesson/otherwise_nothing`.
    2. The `enforcement` of `clause.secondbrain.export-02-variation` becomes `stop-and-ask` in the module file: `TestTheBaselineModulesCarryTheDocumentationRules/the_export_rule`.
    3. In this repository's `docs/agents/docs-layout.md`, the rendered `clause.context.lesson-02-test` is removed: `TestTheGuidesRenderTheDocumentationRulesOnce/repository_guide`.
    4. The notice finder of `internal/spec/adr_header.go` stops skipping fenced code: `TestADRAmendmentNotices/fenced_block_ignored`.
    5. `AGENTS.md` is removed from the process paths the export list planner refuses: `TestPlanExportList/never_a_process_path`.
    6. The refusal `# roundfix:standard-export:begin is missing` loses its last word in the planner: `TestPlanExportList/refuse_begin_missing`.
    7. The Doctor check `baseline-export` gives `DR-EXPORT-REASON-MISSING` the status `failed`: `TestDoctorReportsTheDocumentationRules/export_reason_missing`.
    8. `MaintenanceOnly` no longer counts `.secondbrain-export`: `TestMaintenanceOnlyCountsLessonsAndTheExportList/export_list`.
    9. "Never write a lesson to fill the folder." is removed from the mirror `skills/archive-spec/SKILL.md`: `TestTheSkillsCarryTheDocumentationSteps/archive/never_fill`.
    10. "mark nothing and report the count." is removed from the canonical Roundfix `modes` reference: `TestTheSkillsCarryTheDocumentationSteps/modes/auto_marks_nothing`.
    11. A line `docs/agents/` with a comment above it is added to the own lines of `.secondbrain-export`: `TestThisRepositoryExportsNoProcessPath`.
11. MUST execute the tests task_01 to task_05 name in their Verification.
12. MUST search and record:
    - this repository's three guides, the three formatter goldens of the Standard TypeScript profile and the three modules carry API Contracts 1 to 7, and no new or reworded clause names an ADR number or a Spec number;
    - both copies of `archive-spec` carry API Contract 14 under its own heading before `## Steps`, both copies of the Roundfix `corrections`, `modes`, `baseline`, `setup` and `release` references carry API Contracts 15 to 19, and `.agents/skills` and `skills` hold the same bytes for each changed skill file; no `### QA settlement` section changed;
    - the `release plan` help, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md` carry API Contract 19, and `docs/user-guide/commands/baseline.md`, `docs/user-guide/commands/doctor.md` and `docs/user-guide/context-driven-development.md` carry API Contract 20, the last before `### Spec Event Announcements`;
    - `CONTEXT.md` carries API Contract 21, and ADR-0438, ADR-0439 and ADR-0440 are `accepted`;
    - no command, flag, configuration key or exit code of `roundfix baseline update`, `roundfix doctor`, `roundfix update check` or `roundfix release plan` was added, changed or removed (`_techspec.md` → Surface Transcripts: none).
13. MUST record, as outside-evidence rows, the evidence this Spec did not author, each with its date and none blocked, through what `_prd.md` → Acceptance evidence records, looking nothing up:
    - the count of the project mirrors in `references/2026-11-18-mirrored-documentation-measurement.md`, taken from the Secondbrain at its commit `a2b713574` (recorded 2026-11-18): compare 1,509 process files of 3,534 and no folder of lessons with the standard list, which names no process path and names `docs/lessons/`;
    - Milton, 2017 (read 2026-11-18 and 2026-11-19): compare "Once the lesson has been embedded as a practice change, or a design-principle change, or a change in a checklist, then it can be removed from the database." with the retirement of a Lesson Learned once a rule states it;
    - Lindenbauer, Groh and Schütze, 2025 (read 2026-11-18 and 2026-11-19): compare "noise introduced by distracting CTIM items or exemplar trajectories" with the test that writes nothing by default;
    - Iscan, 2026 (read 2026-11-18 and 2026-11-19): compare "This memory is valuable only when the current failure is genuinely compatible with a previous one." with the lines "Applies when" and "Does not apply when";
    - the public discussion of conflicting ADRs, 2026 (read 2026-11-18 and 2026-11-19): compare "If the supersede note is buried in the body, it never sees it" with the Amendment List in the header;
    - the Ansible documentation of `blockinfile` (read 2026-11-18 and 2026-11-19): compare "may result in the block being repeatedly inserted on subsequent playbook runs" with the two distinct markers and the refusal of a broken pair.
14. MUST record the five entries of `_prd.md` → Unreachable Acceptance, one row each, as `blocked` with `declared: ` followed by the entry's criterion text: this repository's mirror under 60 files after the first sync; Success Metric 10 for the adopters after 30 days; the marking of an adopter's existing ADRs by its update session and a person's review; Success Metric 11 over the first 60 days; and the sync dropping what this repository's own lines exclude with the rsync of its workflow.
15. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD declares `adds:` for **Standard Export List**, **Export Variation**, **Knowledge Mark**, **Amendment List** and **Lesson Learned** and `changes:` for **Managed Refresh**, and the TechSpec declares none; that `CONTEXT.md` defines the five terms with the text of API Contract 21 after the entry of Managed Refresh, which gains only its one sentence, and no other entry changed; and that `TestGlossaryPhrasesAreDefinedTerms` passes.
16. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - each file under `docs/adr/` that existed at the Delivery Base differs from it only by added header lines `knowledge:` and `amended_by:`, and for ADR-0248, ADR-0278 and ADR-0412 by their notice paragraph; no ADR lost a line and no `updated_at` changed;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`, `docs/agents/specific-repository.md`, any file under `docs/agents/` other than `agent-instructions.md`, `docs-layout.md`, `secondbrain.md` and `setup-context.json`, any file under `docs/lessons/`, any module other than the three, any profile other than the Standard TypeScript profile, any Source Baseline or Setup Snapshot, any Vendored Skill under `.agents/skills`, and any skill other than `archive-spec` and `roundfix`, of which `roundfix` only in its `baseline`, `corrections`, `modes`, `release` and `setup` references and version fields;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
17. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
18. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh` or `rsync`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the five Tasks and their packages.
- [ ] Reproduce the Managed Refresh, the documentation part and the Amendment List fill through the binary on temporary copies.
- [ ] Verify this repository's marks, notices and export file, and the two Doctor results on a temporary copy.
- [ ] Verify the clauses and files of the Delivery Base against the audited head.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside evidence, the five declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An adopter's Managed Refresh writes the new clauses into its three guides once each, and a second refresh changes nothing.
- [ ] The update converts an export file, names every path that stops being exported, refuses a broken marker and fills an Amendment List, and a second run changes nothing; the Doctor Command reports each departure as a warning and this repository as `ok`.
- [ ] Every clause of the Delivery Base keeps its id and force, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-13; Core Features 1-20; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; Success Metric 11; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; Invariants 1-10; API Contracts 1-25; Surface Transcripts; Testing Approach; Build Order 6
- [`references/2026-11-18-mirrored-documentation-measurement.md`](references/2026-11-18-mirrored-documentation-measurement.md)
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md)
- ADR-0438, ADR-0439, ADR-0440
