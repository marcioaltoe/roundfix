---
task: task_05
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds the Project Config key `authoring.open_questions`, which the status view of `roundfix follow-up --status` reads at the default branch head and reports as the Authoring Question Preference in force; gives the Delegated Decision a fixed form in the PRD's `## Decisions`, which the Clarification Record's check reads under `SC-DECISIONS-UNRECORDED` and the status view lists for every active Spec; adds to the Baseline's Spec workflow the rule `rule.spec.open-authoring-question` with two clauses; carries the steps into `write-prd`, `write-techspec`, `write-tasks` and the Roundfix skill's `modes`, `follow-up` and `spec` references; and defines the four terms in the glossary and the user guide, with ADR-0442 and the amendment notice of ADR-0280. No command, flag, exit code or check code is added or removed.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon or a Delivery Owner;
- calls ai-memory, runs the real `gh` or reaches the network;
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

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 2 by running, with `-v`, `TestTheAuthoringQuestionPreferenceIsAProjectConfigValue` and `TestRequireDefectClassIsAProjectConfigBoolean` over `./internal/config`, recording every subtest by name, then the whole `./internal/config` package. Through `bin/roundfix`, in a fresh home prepared as above, on a temporary Git repository on the branch `main` with no remote, record that:
   - `bin/roundfix follow-up --status` with a User Config in the fresh home that holds `authoring.open_questions: follow` prints, once on standard error, `config: authoring.open_questions in User Config is ignored; set authoring.open_questions in Project Config`, and the line `Open authoring questions: wait (default)`;
   - with a working-tree Project Config whose `authoring.open_questions` is `maybe`, it prints on standard error `roundfix: follow-up failed: parse config "<path>": authoring.open_questions must be wait or follow` and `Run 'roundfix follow-up --help' for usage.`, and exits 2 (`_techspec.md` → Surface Transcript 3).
3. MUST verify Success Metric 1 by running, with `-v`, `TestTheStatusViewReportsTheAuthoringQuestionPreference` over `./internal/cli`, recording every subtest by name. Through `bin/roundfix` as in requirement 2, record that the status view of the temporary repository prints `Open authoring questions: wait (default)` and, in `--format json`, `"authoring_questions":{"preference":"wait","source":"default"` when `main` holds no key; `Open authoring questions: follow (Project Config)` once `main` commits `authoring.open_questions: follow`; and, after `main` commits the key's removal and the working tree alone holds `follow`, the line `Open authoring questions: wait (default); this checkout's Project Config gives follow, which counts only once the default branch holds it` and the same values in JSON. The first case reproduces the lines of `_techspec.md` → Surface Transcript 1 from `Waiting on the maintainer` to `Sources:`.
4. MUST verify Success Metric 3 by running, with `-v`, `TestDelegatedDecisionsParsesTheEntryForm`, `TestADelegatedDecisionThatCannotBeReadNamesItsFirstProblem` and `TestTheClarificationRecordCheckReportsAnUnreadableDelegatedDecision` over `./internal/speccheck`, recording every subtest by name, then the whole `./internal/speccheck` package. Through `bin/roundfix` as in requirement 2, on a temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, record that `bin/roundfix spec check 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author --strict` reports no `SC-DECISIONS-UNRECORDED` for this Spec's own PRD, whose `## Decisions` holds dated answers, the autonomous-route line and Delegated Decisions; and that after the `adopted` line of one Delegated Decision in that copy's PRD loses its em dash and is committed, the same command reports `SC-DECISIONS-UNRECORDED` with the summary `has a Delegated Decision that cannot be read: adopted needs the option, an em dash and the author's reason` and the entry's line.
5. MUST verify Success Metric 4 by running, with `-v`, `TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs` and `TestTheStatusViewAuthoringSurfaceTranscripts` over `./internal/cli`, recording every subtest by name. Through `bin/roundfix` as in requirement 2, on the temporary repository with `authoring.open_questions: follow` committed on `main` and an active Spec folder `docs/specs/0001-expired-imports` committed with the `_prd.md` and the other files the Spec fixture of `TestTheStatusViewAuthoringSurfaceTranscripts` writes, its `## Decisions` holding the entry of `_prd.md` → User Experience, record that the status view prints, after `Waiting on the maintainer` and before `Sources:`, exactly the three lines of `_techspec.md` → Surface Transcript 2 from `Open authoring questions:` on, and that `--format json` carries one `delegated_decisions` row with `"spec":"0001-expired-imports"`, `"phase":"prd"` and `"answer":"none yet"`.
6. MUST reproduce the Managed Refresh of Success Metric 5 through `bin/roundfix`, in a fresh home prepared as above, on a temporary copy of the Delivery Base (the merge base of the audited head with the default branch) made with `git archive` and committed in a fresh Git repository outside the checkout. Record that `bin/roundfix baseline update --repo . --no-skills --no-history --yes --format text` there writes into `docs/agents/spec-routing.md`, each exactly once after whitespace normalization, "With the user absent the Spec waits, unless following is in force", "Never write that preference into Project Config yourself.", "Never write it as a dated answer of a person or under the autonomous-route line." and "The waits of the Idea stage and of a Rule Challenge keep their own rules.", before the first Project Constraints clause; that the Spec routing guide of the copy did not hold them before; and that a second run of the same command prints `File changes: 0`. Record also that the same command on a temporary copy of the audited head prints `File changes: 0` on its first run.
7. MUST verify Success Metric 6 by running, with `-v`, `TestTheAuthoringSkillsTreatAnOpenAuthoringQuestionAlike`, `TestTheSkillModesFollowTheAuthoringQuestionPreference` and `TestTheClarificationRecordGuideCarriesTheDelegatedDecisionForm` over `./skills`, recording every subtest by name; then `TestTheFollowUpModeRunsItsStepsInOrder`, `TestTheUnattendedModeAsksWithinTheBaselineLimit`, `TestEveryOwnedSkillVersionIsRecorded`, the whole `./skills` package and `make skills-sync-check`, recording each pass. Record that `skills/testdata/owned-skill-versions.json` holds the version each `SKILL.md` of `write-prd`, `write-techspec`, `write-tasks` and `roundfix` declares, in both copies, and that neither copy of `write-tasks` holds the word "recommendation".
8. MUST verify Success Metric 7 by running `-tags docscontract ./internal/docscontract` whole and, with `-v`, `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryAvoidLinesArePlainLists`, `TestGlossaryLinksNameAcceptedADRs` and `TestTheSkillCoverageMapIsCurrent`, then `-tags docscontract -run '^TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes$' ./internal/cli` and `make verify-docs`, recording each pass.
9. MUST verify Success Metric 8 and Invariants 6, 7, 10 and 12 against the Delivery Base, never against a named commit:
   - every clause of every module under `internal/baseline/assets/modules/` at the Delivery Base exists at the audited head with the same id, `enforcement` and `guidance`, the stop-and-ask clause about questions, the Rule Challenge clauses and the Idea-stage route clauses included, and the only clauses the head adds are the two of `rule.spec.open-authoring-question`;
   - every file under `internal/baseline/assets/source-baselines/` and `internal/baseline/assets/setups/` and every module but `spec-workflow.json` is byte-identical to the Delivery Base, and no profile lists `rule.spec.open-authoring-question` under `requiredRules`;
   - `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall`, `TestShippedGuidanceCitesNoRepositoryRecord`, `TestNoTwoBaselineClausesShareText`, `TestCheckCorpusGolden`, `TestStageScopeDefaultSweepIsUnchanged`, `TestFollowUpHelpTranscript` and `TestFollowUpStatusWritesNothing` pass, each run with `-v` in its own package, the tags its file needs included;
   - the only lines `internal/baseline/clause_characterization_test.go` gains over the Delivery Base are the two ids of the new clauses, and it loses none; `internal/speccheck/coherence.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/spec/archive_layout_characterization_test.go` and `docs/references/coverage-record.json` are byte-identical to it;
   - `roundfix follow-up --help` prints the same text as the binary built at the Delivery Base, compared through both binaries in fresh homes prepared as above.
10. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Each sabotage keeps the package compiling and regenerates or synchronizes nothing, so the named test fails on its assertion; record that it fails on the named subtest.
    1. In `internal/config/authoring.go`, the check that a Project Config value is `wait` or `follow` is removed: `TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/refuses_maybe`.
    2. In `internal/config/authoring.go`, the key is no longer removed from User Config: `TestTheAuthoringQuestionPreferenceIsAProjectConfigValue/user_config_ignored_once`.
    3. In `internal/cli/follow_up_status_authoring.go`, the preference reported is the value of the configuration the command loaded from the checkout instead of the value at the default branch head: `TestTheStatusViewReportsTheAuthoringQuestionPreference/follow_only_in_the_working_tree`.
    4. In `internal/cli/follow_up_status_authoring.go`, the `detail` of a Project Config that could not be read is left empty: `TestTheStatusViewReportsTheAuthoringQuestionPreference/unreadable_project_config`.
    5. In `internal/speccheck/delegated_decision.go`, an `adopted` line with a hyphen in place of the em dash is accepted: `TestADelegatedDecisionThatCannotBeReadNamesItsFirstProblem/adopted_without_em_dash`.
    6. In `internal/speccheck/delegated_decision.go`, the head's date is checked by its shape only, not as a calendar date: `TestADelegatedDecisionThatCannotBeReadNamesItsFirstProblem/invalid_date`.
    7. In `internal/speccheck/clarification_record.go`, a line that begins `- delegated:` counts as a dated answer: `TestTheClarificationRecordCheckReportsAnUnreadableDelegatedDecision/only_a_delegated_decision_still_lacks_the_record`.
    8. In `internal/cli/follow_up_status_delegated.go`, rows that carry a problem are dropped: `TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs/an_unreadable_entry_with_its_line`.
    9. The sentence "Never write that preference into Project Config yourself." is removed from the `guidance` of `clause.spec.open-authoring-question-01-wait-or-follow` in the module file: `TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/never_writes_the_preference`.
    10. The `enforcement` of `clause.spec.open-authoring-question-01-wait-or-follow` becomes `mandatory` in the module file: `TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/the_rule`.
    11. In this repository's `docs/agents/spec-routing.md`, the rendered second clause is removed: `TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce/repository_guide`.
    12. In the canonical `write-tasks` skill, inside the section `## Open authoring questions`, "what was found" becomes "what was seen": `TestTheAuthoringSkillsTreatAnOpenAuthoringQuestionAlike/the_core_text_is_shared`.
    13. In the mirror `skills/write-prd/SKILL.md`, the sentence "A session that authors the phase itself makes both reads." is removed: `TestTheAuthoringSkillsTreatAnOpenAuthoringQuestionAlike/both_reads`.
    14. In the canonical `modes` reference, "never write that key." is removed: `TestTheSkillModesFollowTheAuthoringQuestionPreference/unattended/never_writes_the_key`.
    15. In the canonical Clarification Record guide, the line `  - answer: none yet` of the entry form is removed: `TestTheClarificationRecordGuideCarriesTheDelegatedDecisionForm/guide/answer_line`.
11. MUST execute the tests task_01, task_02, task_03 and task_04 name in their Verification.
12. MUST search and record:
    - this repository's `docs/agents/spec-routing.md`, the formatter golden of the Standard TypeScript profile and the `spec-workflow` module carry API Contracts 4 and 5, and neither clause names an ADR number or a Spec number;
    - both copies of the three authoring skills carry API Contracts 6 and 7, both copies of the Clarification Record guide and the PRD template carry API Contracts 8 and 9, both copies of `modes.md` carry API Contract 10, and both copies of `follow-up.md` and `spec.md` carry API Contract 11; `.agents/skills` and `skills` hold the same bytes for each changed skill file;
    - both copies of the TechSpec template, `write-idea` and every Vendored Skill are byte-identical to the Delivery Base, and no `### QA settlement` section changed;
    - `CONTEXT.md` carries API Contract 12, the four user-guide pages carry API Contract 13, ADR-0280 carries API Contract 14 and lists ADR-0442 in its `amended_by` header when `docs/agents/docs-layout.md` defines that field, and ADR-0442 is `accepted`;
    - `.roundfixrc.yml` and the default configuration that `roundfix init` writes hold no `open_questions` (`_techspec.md` → Invariant 3).
13. MUST record, as outside-evidence rows, the published sources this Spec did not write, each with its date and none blocked, through the passages `_prd.md` → Acceptance evidence quotes, looking nothing up:
    - Appelo, "The Seven Levels of Authority" (read 2026-11-19): compare level 6, "You let the team decide. And afterwards you inquire about their motivations", with the Delegated Decision that is recorded and listed;
    - Feng, McDonald and Zhang, 2025 (read 2026-11-19): compare "five levels of escalating agent autonomy, characterized by the roles a user can take" with the preference a project chooses;
    - AG2, "Approval Required" (read 2026-11-19): compare "an unanswered approval is not an approval" with the default `wait`;
    - Cuadron and Gupta, 2026 (read 2026-11-19): compare "complete 95% of requests taken one at a time, but when requests are chained together, they complete 51%" with the list in the status view and the later answer that wins;
    - Bezos, letter to shareholders for 2015 (read 2026-11-19): compare "Type 2 decisions can and should be made quickly" with the line between a delegated option and a Reserved Question;
    - The Apache Software Foundation, "Decision-Making" (read 2026-11-19): compare "Sometimes lazy consensus is not appropriate." with following only where a project chose it.
14. MUST record the two entries of `_prd.md` → Unreachable Acceptance, one row each, as `blocked` with `declared: ` followed by the entry's criterion text: the real flow after the release, and Success Metric 9 and Success Metric 10 at the cycle closes after the release.
15. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD declares `adds:` for **Open Authoring Question**, **Authoring Question Preference**, **Reserved Question** and **Delegated Decision** and `changes:` for **Clarification Record** and **Follow-up Command**, that the TechSpec declares none, that `CONTEXT.md` defines the four with the text of API Contract 12 and carries the two revisions, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
16. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`, `internal/speccheck/coherence.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/spec/archive_layout_characterization_test.go`, any file under `docs/agents/` other than `spec-routing.md` and `setup-context.json`, any file under `docs/user-guide/` other than `configuration.md`, `commands/follow-up.md`, `commands/spec.md` and `context-driven-development.md`, any ADR other than ADR-0280 and ADR-0442, any Vendored Skill under `.agents/skills`, and any skill other than `write-prd`, `write-techspec`, `write-tasks` and `roundfix`, of which `roundfix` only in its `modes`, `follow-up` and `spec` references and its version fields;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
17. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
18. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the four Tasks and their packages.
- [ ] Reproduce the key, the status view, the check and the Managed Refresh through the built binary.
- [ ] Verify the clauses and the frozen files against the Delivery Base.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the two declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The built binary reads the key only from Project Config at the default branch head, refuses a wrong value, lists the Delegated Decisions of an active Spec and reports an unreadable entry through the check.
- [ ] An adopter's Managed Refresh writes both clauses into its Spec routing guide once each, and a second refresh changes nothing.
- [ ] Both copies of the changed skills and references carry their text, the glossary and the user guide define and describe the terms, every clause of the Delivery Base keeps its id, force and text, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-9; Core Features 1-17; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; Invariants 1-14; API Contracts 1-14; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Testing Approach; Build Order 5
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md)
- ADR-0442, ADR-0280
