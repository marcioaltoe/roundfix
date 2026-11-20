---
task: task_02
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
status: pending
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec rewords two clauses of the Baseline's Secondbrain module and the research step of the owned skill `write-idea` so that the Secondbrain search and the advisory duplicate check name the Indexus tool `pesquisar_acervos` and no other search tool, regenerates the module version, the digests, the goldens, the skill mirror and version, and this repository's Secondbrain guide and Setup Manifest, and records the decision in a new decision record, `docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md`, with an amendment notice in ADR-0330. No command, flag, output line, field, exit code or Project Config key is added or changed.

Every row runs on this repository's own files, temporary homes, temporary Git repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a Run, a Daemon or a Delivery Owner;
- calls the Indexus, ai-memory, Exa or any MCP server, runs `qmd`, runs the real `gh` or reaches the network;
- reads the Secondbrain or any file of another repository;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a message or a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, or the Delivery Base binary named in requirement 5, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary Git repository for a product row is made outside the checkout with `git init` and one commit; a temporary copy of the audited head or of the Delivery Base is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Core Features 1 to 4 by running, with `-v` over `./internal/baseline`, `TestTheSecondbrainClausesNameTheIndexus` and `TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten`, recording every subtest by name, each phrase subtest of `_techspec.md` → API Contract 4 among them, in the module, the planned guide, the formatter golden and this repository's guide; then `TestBaselineClauseForceIsCharacterized`, `TestEveryInboxClauseAnAdopterHoldsIsReplaced`, `TestFormatterComposition` and `TestEveryBaselineModuleVersionIsRecorded`, recording each pass, and the whole `./internal/baseline` package, recording its pass.
3. MUST verify Success Metric 3 and Core Feature 5 by running, with `-v` over `./skills`, `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus`, `TestNoOwnedSkillNamesTheRetiredSearchCommand`, `TestAuthorialSkillSync` and `TestEveryOwnedSkillVersionIsRecorded`, recording each test and subtest by name; record that both copies of `write-idea` are byte-identical with `cmp`, and that their version is the one `skills/testdata/owned-skill-versions.json` records for those bytes.
4. MUST verify Success Metric 5 and Core Feature 7 by running, with `-tags docscontract` and `-v` over `./internal/docscontract`, `TestTheIndexusDecisionAmendsTheConsultationRule` and `TestTheSetupManifestIsCurrent`, recording each pass, and by recording from the files: the new decision record's header with `status: accepted` and `knowledge: true`; ADR-0330's `amended_by` ending with the new record's identifier and its `updated_at` equal to the new record's `created_at`; the notice paragraph directly after the first paragraph of ADR-0330's `## Decision`, with its link and date; the sentence that names the earlier command byte-identical to the Delivery Base; and the line of the new record's path in `.secondbrain-export` between the Standard Export List markers, in ascending order. Record that the new record names no client and no acervo and cites no Spec folder path.
5. MUST verify Success Metric 2 and `_techspec.md` → API Contract 7 through the built binaries, in a fresh home prepared as above:
   - build a Delivery Base binary from a `git archive` of the Delivery Base, the merge base of the audited head with the default branch, in a temporary copy outside the checkout;
   - in a temporary Git repository made outside the checkout, adopt the Standard TypeScript profile with that Delivery Base binary through `baseline plan --profile standard-typescript-monorepo` with the decision `secondbrain.enabled=true` and the other decisions the plan asks for, then `baseline apply` with the plan's digest, and commit;
   - run `bin/roundfix baseline update --repo <that repository> --no-skills --no-history --format json` without `--yes`, and record the whole output and exit code: the plan changes `docs/agents/secondbrain.md` and `docs/agents/setup-context.json` and no other file; `clause.secondbrain.02-query-order` is `retained` with `mandatory` force; `clause.secondbrain.research-capture` is `replaced` with the target `clause.secondbrain.research-to-the-knowledge-owner` and `mandatory` force; no clause is `unaccounted`; record the clause counts by disposition as the command prints them, never as numbers fixed in this plan;
   - apply it with `--yes`, record that the written `docs/agents/secondbrain.md` holds the phrases of API Contract 4 and no whole-word `qmd` in any letter case, then run the update again and record that it plans no change.
   When the adoption or the update needs the network for a reason this Spec does not own, record the row as blocked by the environment with that reason and keep `TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest`, run with `-v` over `./internal/baseline` and recorded by name, as its evidence.
6. MUST verify Core Feature 6 by running, on a temporary copy of the audited head, `make skills-sync`, `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, `make baseline-digests` and `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, and recording that they change no file, which proves that the skill mirror and versions, the module version, the digests, the goldens, this repository's Secondbrain guide and Setup Manifest, the Standard Export List and ADR-0330's `amended_by` are the output of the repository's commands; and record that `bin/roundfix baseline update --repo . --no-skills --no-history --format text`, run on a temporary copy of the audited head without `--yes`, plans no change.
7. MUST verify Success Metric 6 and Core Feature 8 with a whole-word, case-insensitive search of the audited tree for `qmd` outside `.git`, and record every file it finds against this list: the Source Baseline corpus copy of the Secondbrain guide; the two receipt fixtures `internal/speccheck/testdata/receipt-characterization/0181-prd.md` and `0182-prd.md`; `docs/references/archived-evidence-measurement.md`; ADR-0330, in its decision sentence and its notice; the new decision record; the Archive Records under `docs/history/`; the Vendored Skill `.agents/skills/knowledge-workspace/SKILL.md`; the three new test files; this Spec's folder; and the folder of any other active Spec under `docs/specs/` that this Spec's change does not touch, whose planning pages may record earlier consultations. Record a finding for any other file. Then run `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand` with `-v` over `./internal/baseline` and `TestNoOwnedSkillNamesTheRetiredSearchCommand` with `-v` over `./skills`, recording each pass.
8. MUST run each sabotage below on its own temporary copy of the audited head. Each sabotage keeps the packages compiling and regenerates nothing, so the named test fails on its assertion; record that it fails on the named test or subtest.
   1. The word `qmd` is put back into the search clause of `internal/baseline/assets/modules/secondbrain.json`: `TestTheSecondbrainClausesNameTheIndexus/no_retired_word` and `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand` fail.
   2. The phrase "substitute no other search tool" is removed from the search clause of the module: `TestTheSecondbrainClausesNameTheIndexus/search_clause/no_other_search_tool` fails.
   3. The force of `clause.secondbrain.02-query-order` becomes `stop-and-ask`, a force the catalog accepts, so the catalog still loads: `TestTheSecondbrainClausesNameTheIndexus/keeps_identity_and_force` fails.
   4. `pesquisar_acervos` is replaced by another tool name in `skills/write-idea/SKILL.md` only: `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus` fails.
   5. The notice paragraph is removed from ADR-0330: `TestTheIndexusDecisionAmendsTheConsultationRule` fails.
   6. The Source Baseline corpus copy of the Secondbrain guide loses the word `qmd`: `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand` fails on its stale exemption.
9. MUST execute the tests task_01 names in its Verification.
10. MUST record the outside-evidence row of Success Metric 4, with its date and not blocked, through what `_prd.md` → Acceptance evidence records, looking nothing up: compare the tool names `pesquisar_acervos` and `abrir_conteudo`, the server name `indexus`, "one group per acervo the session's key reaches" and the `caminho` of a `wiki` result relative to `wiki/` in the reworded clauses with the contract `indexus-acervos/1` and the live calls that `references/2026-11-20-characterization.md`, sections 5 and 7, and `references/2026-11-20-messages-and-decisions.md`, section 9, record, quoting the passage of each source the comparison rests on.
11. MUST record the two rows of `_prd.md` → Unreachable Acceptance as blocked by declaration, each with its reason and its `satisfied-by`, and record that the authorization record grants neither `merge` nor `release`.
12. MUST perform the glossary check (`docs/agents/domain.md`): record that the PRD and the TechSpec declare no glossary change, that `CONTEXT.md` is byte-identical to the Delivery Base, and that no term the Spec uses, the Indexus included, needs a definition the glossary lacks.
13. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`, and the Governed Paths changed are only the seven that list names;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CONTEXT.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/behavior-surfaces.json`, `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`, any file under `internal/baseline/assets/source-baselines/`, any file under `docs/agents/` other than `secondbrain.md` and `setup-context.json`, any user-guide page, any Vendored Skill and any owned skill other than `write-idea`;
    - in `internal/baseline/assets/modules/secondbrain.json`, only the guidance of the two clauses and the module's version line differ from the Delivery Base;
    - no accepted ADR other than ADR-0330 changed, ADR-0330 stays `accepted`, and the new decision record is `accepted`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
14. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
15. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory, the Indexus or any MCP server, run `qmd` or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of task_01 and their packages.
- [ ] Reproduce an adopter's update from the Delivery Base through the binaries in a temporary repository.
- [ ] Prove the regenerated records are the commands' output on a temporary copy, and search the tree for the retired word.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside evidence, the unreachable rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The two clauses, the planned and written guides, the formatter golden and both copies of `write-idea` name the Indexus tool and not the retired command, and both clauses keep id and force.
- [ ] An adopter of the Delivery Base with the Secondbrain turned on receives a plan that changes only its Secondbrain guide and its Setup Manifest, with the clause accounting of API Contract 7, and plans no change after the apply.
- [ ] The new decision record records the decision, ADR-0330 carries its notice and amendment entry, every sabotage makes its named test fail, the retired word remains only where Core Feature 8 and the TechSpec's deviation declare it, and the scope audit finds no change outside the declared files and the seven granted Governed Paths.

## Verification

- `newest="$(find 'docs/specs/0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-6; Core Features 1-8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Delivery limits; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces; API Contract 4; API Contract 6; API Contract 7; Coverage Map; Testing Approach; Risks & Considerations; Build Order 3
- `_authorization.md` → What is not granted; Sanctioned regeneration; Limits
- [`references/2026-11-20-characterization.md`](references/2026-11-20-characterization.md), sections 1, 2, 5 and 7
- [`references/2026-11-20-messages-and-decisions.md`](references/2026-11-20-messages-and-decisions.md), section 9
- ADR-0330; the new decision record, `docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md`
