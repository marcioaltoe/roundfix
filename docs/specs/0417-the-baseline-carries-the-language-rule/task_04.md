---
task: task_04
spec: 0417-the-baseline-carries-the-language-rule
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds the Language Rule to the `core` Baseline module, the required decision `language.reader` and its line in the agent-instructions guide, the Name in Code inside the reworded canonical-terms clause, the optional decision `language.code_exception` with the `present` clause gate, one instruction to the Pre-PR Review prompt and one section to the `qa-gate` skill, with ADR 0454, five Glossary Definitions, the skills and the user-guide pages that describe them, and this repository's Managed Refresh with the Reader Language `English`.

Every row runs on this repository's own files, temporary homes, temporary Git repositories and temporary copies of the audited head or of the Delivery Base, the merge base of the audited head with the default branch. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon, a Delivery Owner or a Pre-PR Review;
- runs `roundfix baseline update` without `--no-skills` and `--no-history`, or a skills restore that reaches an upstream repository;
- calls ai-memory, sends a Cross-project Message or a Spec Event Announcement, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary copy of a commit is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard and the update refuse a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST execute the tests task_01, task_02 and task_03 name in their Verification, recording each named test and subtest that passes, then the whole `./internal/baseline`, `./internal/cli` and `./skills` packages, `go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` and `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`, recording each pass.
3. MUST verify Success Metric 1 by running, with `-v` over `./internal/baseline`, `TestEveryBuiltInProfileRendersTheLanguageRule`, recording its five profile subtests and `a_catalog_without_a_clause_fails`; and by reading, in the formatter golden of the Standard TypeScript profile's agent-instructions guide at the audited head, the line `The Reader Language is ` followed by the answer, and the four clauses of `_techspec.md` → API Contract 3 word for word with mandatory force.
4. MUST verify Success Metric 2 and Surface Transcript 1 through `bin/roundfix`, in a fresh home prepared as above, on a temporary copy of the Delivery Base, whose Setup Manifest lacks the answer:
   - `bin/roundfix baseline update --repo <copy> --no-skills --no-history --format text` prints Surface Transcript 1 line by line, with only `<repository>` and `<digest>` varying, exits `3`, and leaves `git status --porcelain` of the copy empty;
   - `bin/roundfix baseline update --repo <copy> --no-skills --no-history --yes --format json --decision language.reader=English` succeeds; record the files `git status --porcelain` lists as the command left them, which must be `docs/agents/agent-instructions.md`, `docs/agents/domain.md` and `docs/agents/setup-context.json` and no other; record from the JSON result that its file changes are those three files and no other. The Delivery Base's Setup Manifest (profile `go-cli-tui`) records no Source Baseline, so its JSON result carries no clause delta and no retention list: the accounting of every clause a repository held, with none unaccounted and `clause.domain.canonical-language` `retained` with mandatory force, is proven by `TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles` and `TestStandardTypeScriptStructuralClauseRetention` on a repository whose manifest records a Source Baseline, recorded in item 2; on the copy, record that `docs/agents/domain.md` after the update holds, as a `**mandatory**` bullet, the canonical-terms clause that starts "Use the repository's canonical domain terms in code names" and carries the Name in Code sentence (amended 2026-11-27 after QA finding F-01: the Delivery Base cannot expose a clause delta);
   - a second `bin/roundfix baseline update --repo <copy> --no-skills --no-history --format json` reports the state `current` and plans no change.
   Then run, with `-v` over `./internal/cli`, `TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles`, `TestTheUpdateNamesTheReaderLanguageDecision` and `TestBaselineUpdateFleetSweep`, recording each pass, and over `./internal/baseline` `TestStandardTypeScriptStructuralClauseRetention`.
5. MUST verify Success Metric 3 and Surface Transcript 2 through `bin/roundfix`, prepared the same way, each case on its own temporary copy of the Delivery Base:
   - `bin/roundfix baseline update --repo <copy> --no-skills --no-history --format text --decision language.reader=pt-BR --decision 'language.code_exception={"language":"Portuguese"}'` prints Surface Transcript 2 line by line, with only `<repository>` and `<digest>` varying, exits `2`, and leaves `git status --porcelain` empty;
   - `bin/roundfix baseline update --repo <copy> --no-skills --no-history --yes --format text --decision language.reader=pt-BR --decision 'language.code_exception={"language":"Portuguese","source":"docs/adr/0141-example.md"}'` succeeds, and the copy's `docs/agents/agent-instructions.md` states the sentence of the second block of `_techspec.md` → API Contract 11 with that language and that source in a code span, holds the first block of API Contract 11, and holds none of `clause.core.code-names-in-english`, `clause.core.code-names-keep-untranslated-terms` or `clause.core.code-names-new-only` by their texts;
   - `bin/roundfix baseline update --repo <copy> --no-skills --no-history --yes --format text --decision language.reader=pt-BR` renders the English clause of `clause.core.code-names-in-english` and not the first block of API Contract 11.
   Then run, with `-v` over `./internal/baseline`, `TestACodeLanguageExceptionSwapsTheCodeClauses`, `TestACodeLanguageExceptionValueIsRefused`, `TestAClauseGateOnPresence` and `TestTheCatalogRefusesADefaultForACodeLanguageException`, and over `./internal/cli` `TestAnExceptionWithoutASourceIsRefused` and `TestTheInteractiveFlowAsksNothingForTheCodeLanguageException`, recording every subtest by name.
6. MUST verify Success Metric 4 by running, with `-v` over `./internal/cli`, `TestARepositoryProfileWithoutTheReaderLanguageIsRefused`, recording its pass, and through `bin/roundfix`, in a fresh home, in a temporary Git repository made outside the checkout with `git init` and one commit: `bin/roundfix baseline profile init --id own-profile --from go-cli-tui`, the entry `language.reader` removed from the new profile's decision list, and `bin/roundfix baseline profile validate own-profile`, recording its output and exit code; the output names `language.reader`, the template that renders it and the entry to add.
7. MUST verify Success Metric 5 by running, with `-v`, `TestTheReviewPromptCarriesTheLanguageRuleInstruction` over `./internal/cli` and, with `-tags docscontract`, `TestTheLanguageRuleInstructionIsDocumentedAsTheReviewerPromptsIt`, `TestTheLanguageRuleInstructionContractReportsAStaleQuotation` and `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem` over `./internal/docscontract`, recording each pass; and by reading that `reviewLanguageRuleInstruction` in `internal/cli/review.go` holds the text of `_techspec.md` → API Contract 13 word for word, and that `deliveryConventionsVersion` and the finding grammar lines of `buildReviewPrompt` are unchanged from the Delivery Base.
8. MUST verify Success Metric 6 by running, with `-v` over `./skills`, `TestTheQAGateReportsTextOutsideTheReaderLanguage`, `TestAuthorialSkillSync`, `TestEveryOwnedSkillVersionIsRecorded` and `TestSettlementGuidanceIsOneTable`, recording each pass; by recording that `.agents/skills/qa-gate/SKILL.md` and `skills/qa-gate/SKILL.md` are byte-identical and carry the paragraph of `_techspec.md` → API Contract 14 under `### Reader Language` after `### Frontend sweep`; and that the `### QA settlement` section of the skill is byte-identical with the Delivery Base.
9. MUST verify Success Metric 7: ADR 0454 exists, is `accepted` and carries `knowledge: true`; its path is in `.secondbrain-export`; and on a temporary copy of the audited head, `bin/roundfix baseline update --repo <copy> --no-skills --no-history --format json`, in a fresh home, reports the state `current`, which proves the Standard Export List, the guides and the Setup Manifest are the ones the update writes. Record that this repository's `docs/agents/agent-instructions.md` states `The Reader Language is English.` and no Code Language Exception.
10. MUST verify Success Metric 8 by measuring, with `wc -c`, the total bytes of the Markdown guides under `docs/agents/` at the audited head and at the Delivery Base, recording both totals and their difference as the command prints them, and that the difference is at most 1,500 bytes.
11. MUST verify Success Metric 10 from `git diff --name-status -M` between the Delivery Base and the audited head: no file is renamed, and none of the Makefile, `.github/workflows/`, `go.mod`, `go.sum`, `.roundfixrc.yml`, lint configuration, `docs/references/gate-time-targets.json` or `cmd/gate-time/` changed.
12. MUST prove that every derived file is the output of its command: on a temporary copy of the audited head, run, in order, the seven commands of `_authorization.md` → Sanctioned regeneration, the Managed Refresh among them with `--decision language.reader=English`, and record that `git status --porcelain` of the copy is empty afterwards, so the module versions, the digests, the goldens, this repository's guides and Setup Manifest, the skill mirrors and versions, the surface record and the coverage record are the commands' output.
13. MUST run each sabotage below on its own temporary copy of the audited head. Each keeps the packages compiling and regenerates nothing, so the named test fails on its assertion; record that it fails on the named test or subtest.
    1. `clause.core.code-names-in-english` is removed from `internal/baseline/assets/modules/core.json`: `TestEveryBuiltInProfileRendersTheLanguageRule` fails on a profile subtest.
    2. `normalizeCodeLanguageException` accepts a value without a source: `TestACodeLanguageExceptionValueIsRefused/no_source` fails.
    3. `clauseApplies` treats a decision's default as presence: `TestAClauseGateOnPresence/a_default_is_not_presence` fails.
    4. `buildReviewPrompt` no longer writes `reviewLanguageRuleInstruction`: `TestTheReviewPromptCarriesTheLanguageRuleInstruction` fails.
    5. The paragraph under `### Reader Language` is removed from both copies of the `qa-gate` skill: `TestTheQAGateReportsTextOutsideTheReaderLanguage` fails.
14. MUST record the outside-evidence row of Success Metric 9, with its date and not blocked, through the passages `_prd.md` → Acceptance evidence quotes and `references/2026-11-27-characterization.md`, section 5, holds, looking nothing up: Bartlett, "business terms stay in the native language, the rest of the code stays in English", against `clause.core.code-names-in-english` and `clause.core.code-names-keep-untranslated-terms`; and Doomen, "create a glossary that maps the verbs and nouns from the native domain model into the English translations", against the Name in Code sentences of `clause.domain.canonical-language`. The row is satisfied when the rendered clauses state the split and the mapping these passages state.
15. MUST record the three criteria `_prd.md` → Unreachable Acceptance declares, Success Metrics 11, 12 and 13, as declared unreachable with their `satisfied-by`, and never as blocked by the environment.
16. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD declares `adds:` for **Language Rule**, **Reader Language**, **Untranslated Term**, **Name in Code** and **Code Language Exception**; that `CONTEXT.md` gives each a Glossary Definition line citing ADR 0454 and no `_In code_` line; that no other entry changed from the Delivery Base; and that `TestGlossaryPhrasesAreDefinedTerms` passes with `-tags docscontract` over `./internal/docscontract`.
17. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`; record the list of changed Governed Paths as the comparison prints it;
    - no file under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/`, no file under `internal/baseline/testdata/parity-corpus/`, `docs/references/skill-coverage.json`, `skills-lock.json`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json` or `AGENTS.md` changed in the Task commit range;
    - `docs/references/coverage-record.json` removes no test name that the Delivery Base holds, and every name it adds is a test a Task of this graph names or `TestRunsInWindowKeepsTheOldestRunsSinceTheStart`, a test the Delivery Base already holds that its stale record lacked and the authorized command (g) records (coordinator decision of 2026-11-27, ai-memory `01a21e0b-4631`); record the added names as the comparison prints them and note that line as inherited from the Delivery Base;
    - no accepted ADR other than ADR 0454 changed, and no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
18. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority, `_authorization.md`;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
19. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0417-the-baseline-carries-the-language-rule --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon, a Delivery Owner or a Pre-PR Review;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary, run the repository Verification and the tests the three Tasks name with their packages.
- [ ] Reproduce both Surface Transcripts, the three-file update and the exception through the binary on copies of the Delivery Base.
- [ ] Prove the derived files are the commands' output on a copy of the audited head, and run each sabotage.
- [ ] Measure the guides' growth and read the review prompt, the QA gate skill, the ADR and the export list.
- [ ] Record the outside evidence, the Unreachable Acceptance rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every built-in profile renders the Language Rule and the Reader Language line, and an update without the answer stops with Surface Transcript 1 and writes nothing; with the answer it changes the two guides and the Setup Manifest only, and a second update plans nothing.
- [ ] A declared exception with a source swaps the code clauses and names its source; one without a source prints Surface Transcript 2 and writes nothing; a repository that declares nothing keeps the English clause.
- [ ] The review prompt and the QA gate skill carry their instructions, every sabotage makes its named test fail, the guides grow by at most 1,500 bytes, and the scope audit finds no change outside the declared files and the granted Governed Paths.

## Verification

- `newest="$(find 'docs/specs/0417-the-baseline-carries-the-language-rule/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-8; Core Features 1-13; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; Success Metric 11; Success Metric 12; Success Metric 13; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces, invariants 1-7; API Contract 3; API Contract 4; API Contract 7; API Contract 11; API Contract 13; API Contract 14; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 4
- `_authorization.md` → Sanctioned regeneration
- [`references/2026-11-27-characterization.md`](references/2026-11-27-characterization.md), sections 5 and 8
- The ADR task_01 writes, numbered 0454
