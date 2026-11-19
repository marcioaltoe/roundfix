---
task: task_04
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds to the Baseline's Spec workflow the rule `rule.spec.rule-challenge`, whose two clauses tell an author of brainstorming, an Idea or a PRD to put a Rule Challenge to the user and say what follows the answer; carries the step into `brainstorming`, `write-idea`, `write-prd` and `write-techspec`, the record into the Idea and PRD templates and the wait into the Roundfix skill's `modes` reference; and defines Rule Challenge and Kept Line in the glossary and the user guide, with ADR-0436 and ADR-0437. No command, flag, configuration key, exit code or check changes.

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
2. MUST verify Success Metric 1 and Invariants 1 to 4 by running, with `-v`, `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses` and `TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce` over `./internal/baseline`, recording every subtest by name; then the whole `./internal/baseline` package and `-run 'Baseline|GuidanceComposition|SetupManifest' ./internal/cli`, recording each pass.
3. MUST reproduce the Managed Refresh of Success Metric 1 through `bin/roundfix`, in a fresh home prepared as above, on a temporary copy of the Delivery Base (the merge base of the audited head with the default branch) made with `git archive` and committed in a fresh Git repository outside the checkout. Record that `bin/roundfix baseline update --repo . --no-skills --no-history --yes --format text` there writes into `docs/agents/spec-routing.md`, each exactly once after whitespace normalization, "never drop or narrow a solution in silence because a rule forbids it, and never contradict a rule in silence: put a Rule Challenge to the user.", "In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge.", "a rule keeps every Kept Line, the newest last" and "An Author Subagent returns the question to the session that started it.", after the last route clause and before the first Project Constraints clause; that the Spec routing guide of the copy did not hold them before; and that a second run of the same command prints `File changes: 0`. Record also that the same command on a temporary copy of the audited head prints `File changes: 0` on its first run.
4. MUST verify Success Metric 2 and Invariants 6 to 8 by running, with `-v`, `TestTheAuthoringSkillsPutARuleChallengeToTheUser`, `TestTheTemplatesCarryTheRuleChallengeRecord` and `TestTheSkillModesNeverAnswerARuleChallenge` over `./skills`, recording every subtest by name; then `TestTheFollowUpModeRunsItsStepsInOrder`, `TestTheUnattendedModeAsksWithinTheBaselineLimit`, `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage`, `TestTheFollowUpModeRunsTheIdeaStageJudgment` and `TestEveryOwnedSkillVersionIsRecorded`, the whole `./skills` package and `make skills-sync-check`, recording each pass. Record that `skills/testdata/owned-skill-versions.json` holds the version each `SKILL.md` of `brainstorming`, `write-idea`, `write-prd`, `write-techspec` and `roundfix` declares, in both copies.
5. MUST verify Success Metric 3 by running `-tags docscontract ./internal/docscontract` whole and, with `-v`, `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryAvoidLinesArePlainLists`, `TestGlossaryLinksNameAcceptedADRs` and `TestTheSkillCoverageMapIsCurrent`, and `make verify-docs`, recording each pass.
6. MUST verify Success Metric 4 and Invariant 5 against the Delivery Base, never against a named commit:
   - every clause of every module under `internal/baseline/assets/modules/` at the Delivery Base exists at the audited head with the same id, `enforcement` and `guidance`, and the only clauses the head adds are the two of `rule.spec.rule-challenge`;
   - every file under `internal/baseline/assets/source-baselines/` and `internal/baseline/assets/setups/` and every module but `spec-workflow.json` is byte-identical to the Delivery Base, and no profile lists `rule.spec.rule-challenge` under `requiredRules`;
   - `TestInstructionHierarchyRendersActivePointersOnce`, `TestTheGlossaryClausesCarryTheirForceAndText`, `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall`, `TestShippedGuidanceCitesNoRepositoryRecord`, `TestNoTwoBaselineClausesShareText` and `TestTheRouteClausesJudgeAFeatureForTheIdeaStage` pass, run with `-v` over `./internal/baseline`;
   - the only lines `internal/baseline/clause_characterization_test.go` gains over the Delivery Base are the two ids of the new clauses, and it loses none, and `internal/baseline/plan_test.go` and `skills/baseline_skill_contract_test.go` are byte-identical to it.
7. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Each sabotage keeps the package compiling and regenerates or synchronizes nothing, so the named test fails on its assertion; record that it fails on the named subtest.
   1. The sentence "In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge." is removed from the `guidance` of `clause.spec.rule-challenge-01-ask` in the module file: `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/implementation_guard_rail`.
   2. The `enforcement` of `clause.spec.rule-challenge-01-ask` becomes `mandatory` in the module file: `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/the_rule`.
   3. ", the newest last" is removed from the `guidance` of `clause.spec.rule-challenge-02-answer`: `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/newest_last`.
   4. In this repository's `docs/agents/spec-routing.md`, the rendered second clause is removed: `TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce/repository_guide`.
   5. The bullet that begins "Ask at most three" is removed from the canonical `write-techspec` skill: `TestTheAuthoringSkillsPutARuleChallengeToTheUser/at_most_three`.
   6. In the mirror `skills/brainstorming/SKILL.md`, "Never answer for the user." becomes "Answer for the user.": `TestTheAuthoringSkillsPutARuleChallengeToTheUser/never_answers_for_the_user`.
   7. In the canonical `write-idea` skill, one word of the core text changes inside the sentence "Extra work, taste and friction are not a block.": `TestTheAuthoringSkillsPutARuleChallengeToTheUser/the_core_text_is_shared`.
   8. The `new fact` line of the comment is removed from the canonical PRD template: `TestTheTemplatesCarryTheRuleChallengeRecord/prd/fields`.
   9. "this mode asks it and never answers it" is removed from the canonical Roundfix `modes` reference: `TestTheSkillModesNeverAnswerARuleChallenge/unattended`.
8. MUST execute the tests task_01, task_02 and task_03 name in their Verification.
9. MUST search and record:
   - this repository's `docs/agents/spec-routing.md`, the formatter golden of the Standard TypeScript profile and the `spec-workflow` module carry API Contracts 1 and 2, and neither clause names an ADR number or a Spec number;
   - both copies of the four skills carry API Contracts 3, 4 and 5, both copies of the two templates carry API Contract 6, and both copies of `modes.md` carry API Contract 7; `.agents/skills` and `skills` hold the same bytes for each changed skill file;
   - both copies of the TechSpec template are byte-identical to the Delivery Base, and no `### QA settlement` section changed;
   - `CONTEXT.md` carries API Contract 8, `docs/user-guide/context-driven-development.md` carries API Contract 9 before `### Spec Event Announcements`, and ADR-0436 and ADR-0437 are `accepted`;
   - no Go file outside `_test.go` files changed, and no command, flag, help text or exit code changed (`_techspec.md` → Surface Transcripts: none).
10. MUST record, as outside-evidence rows, the published sources this Spec did not write, each with its date and none blocked, through the passages `_prd.md` → Acceptance evidence quotes, looking nothing up:
    - Xu and Wu, 2026 (read 2026-11-18 and 2026-11-19): compare "Claude Code pointed out the conflict in 79 of 80 runs and still followed the file." and the rise "from 4.2% to 91.4%" with the first clause, which forbids both silences in one sentence of the instructions;
    - Zhang and others, 2026 (read 2026-11-18 and 2026-11-19): compare "limiting agents to one to three questions reduces user-attention cost by 36–73% while keeping 86–94% of the task reward" with the cap of three, and "reviewing each question against the memory of users' earlier answers cuts attention cost by another 13%" with the Kept Line read before asking again;
    - Fowler, "Architecture Decision Record" (read 2026-11-19): compare "Once an ADR is accepted, it should never be reopened or changed - instead it should be superseded." with the change answer, which goes through the ADR lifecycle;
    - Shadewright, "Amending accepted ADRs" (read 2026-11-19): compare "Every amendment gets a dated entry" with the dated Kept Line.
11. MUST record the two entries of `_prd.md` → Unreachable Acceptance, one row each, as `blocked` with `declared: ` followed by the entry's criterion text: the live Rule Challenge in a real authoring session after the release, and Success Metrics 5, 6 and 7 over the first 20 plans and 30 days.
12. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD declares `adds:` for **Rule Challenge** and **Kept Line** and the TechSpec declares none, that `CONTEXT.md` defines both with the text of API Contract 8 and no other entry changed, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
13. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`, any file under `internal/` other than the files task_01 declares, any file under `docs/agents/` other than `spec-routing.md` and `setup-context.json`, any file under `docs/user-guide/` other than `context-driven-development.md`, any ADR, any Vendored Skill under `.agents/skills`, and any skill other than `brainstorming`, `write-idea`, `write-prd`, `write-techspec` and `roundfix`, of which `roundfix` only in its `modes` reference and version fields;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
14. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
15. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0409-the-idea-and-prd-question-rules-that-block-the-goal --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
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
- [ ] Run the named tests of the three Tasks and their packages.
- [ ] Reproduce the Managed Refresh on a copy of the Delivery Base and on a copy of the audited head.
- [ ] Verify the clauses of the Delivery Base against the audited head.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the two declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An adopter's Managed Refresh writes both clauses into its Spec routing guide once each, in their place, and a second refresh changes nothing.
- [ ] Both copies of the four skills, the two templates and the `modes` reference carry their text, and the glossary and the user guide define and describe the Rule Challenge.
- [ ] Every clause of the Delivery Base keeps its id, force and text, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0409-the-idea-and-prd-question-rules-that-block-the-goal/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-10; Core Features 1-15; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; Invariants 1-10; API Contracts 1-9; Surface Transcripts; Testing Approach; Build Order 4
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md)
- ADR-0436, ADR-0437
