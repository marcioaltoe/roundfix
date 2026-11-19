# Prototype of the documentation rules, the tests that pin them, the sync's exclusions and the Governed Paths

Recorded by the author of the PRD and the TechSpec on 2026-11-19, on the Mac mini. The prototype ran on a scratch worktree of this repository at `0ea36244` (the plan of Spec 0409), outside the planning clone, and the worktree was removed afterwards. It answers five questions: where the new Baseline clauses can live, which files the text and its regeneration rewrite, which existing tests pin text the Spec changes, how the sync treats an exclusion, and which files are Governed Paths. No Go code of the Spec was written: the update's documentation part, the Doctor checks and the Release Plan rule are designed in the TechSpec and not prototyped.

## 1. What was applied

The exact texts of the TechSpec, applied by two scripts that read the clause texts and the TechSpec's own blocks:

- API Contracts 1 to 7: four reworded clauses and the new rule `rule.context.lessons` in `internal/baseline/assets/modules/context-workflow.json`, one reworded clause in `core.json`, the new rule `rule.secondbrain.knowledge-export` in `secondbrain.json`, the two guides' rule lists and guide versions;
- API Contract 9: five lines in the closed list of `internal/baseline/clause_characterization_test.go`;
- the three existing tests of section 3, updated to the new text;
- API Contracts 14 to 19: the section of `archive-spec`, the section of the `corrections` reference, the three additions to the `modes` reference, the paragraphs of the `baseline` and `setup` references, and the release sentence in the `release` reference and in the two user-guide pages;
- API Contract 21: the five glossary entries and the sentence added to Managed Refresh;
- ADR-0438, ADR-0439 and ADR-0440, copied from the planning clone.

Not applied: the user-guide subsections of API Contract 20, the `release plan` help, this repository's export file, the marks of its ADRs and the three amendment notices.

## 2. New clauses need new rules, and reworded clauses keep their ids

The Standard TypeScript profile lists `rule.context.docs-layout` under `requiredRules`, with `rule.context.domain-docs`, five rules of the `secondbrain` module (`index-first`, `read-only`, `secret-safety`, `citations`, `hermes-escalation`) and thirteen of `core`. `rule.core.cross-project-continuity` and `rule.secondbrain.research-digests` are not required. A new clause in a required rule is refused by the catalog until a Source Baseline row is written by hand, as the prototype of Spec 0409 measured on `rule.spec.routing`.

The design therefore changes text in two ways, and both passed `make baseline-digests`:

- four clauses of `rule.context.docs-layout` and one of `rule.core.cross-project-continuity` are reworded under their ids and their force, which the Source Baseline transition records as retained;
- `rule.context.lessons`, listed by `guide.docs-layout` after `rule.context.docs-layout`, and `rule.secondbrain.knowledge-export`, listed last by `guide.secondbrain`, are new rules that no profile requires.

No file under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/` changed.

## 3. Existing tests that pin text the Spec changes

`go test ./internal/baseline -count=1` on the first attempt failed in four tests, all of them about text this Spec changes on purpose:

| test | what it pins | what the Spec does |
| --- | --- | --- |
| `TestTheCaptureTriggerAsksForAConditionalLesson` (`internal/baseline/retirement_rule_clauses_test.go`) | the phrases "ADR, Archive Record or reference", "Only then send one message" and `<repository>@<commit>:<path>` in `clause.core.message-the-owning-project` | the three phrases become "send no message about what it taught" and "Lesson Learned"; "When a Finding closes" stays |
| `TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide` (`internal/baseline/repository_rule_placement_test.go`) | the whole text of `clause.context.docs-one-job-per-directory`, in the constant `docsOneJobPerDirectoryGuidance`, in the module and in the golden | the constant takes the text of API Contract 3 |
| `TestTheDocsLayoutRegionStaysLean` (`internal/baseline/retirement_rule_clauses_test.go`) | the managed region `guide.docs-layout` of the Standard TypeScript golden is at most 11,000 bytes | the ceiling becomes 13,000 |
| `TestADRLifecycleContract` (`internal/baseline/plan_test.go`, a Governed Path) | the sentence "Do not rewrite existing ADRs solely to adopt lifecycle metadata." with its period | nothing: API Contract 2 keeps the sentence whole and adds the exception as the next sentence |

The region measured 10,441 bytes at `0ea36244`. The first wording of the new text took it to 14,080, a second to 13,188, and the wording of the TechSpec to 12,948. Spec 0272 had set the ceiling of 11,000 on 2026-10-16, down from 14,494.

No other test of the Baseline package, of the skills package, of the documentation contracts or of the CLI's Baseline, Doctor and Update Check tests pins a sentence the Spec changes. The two tests of `internal/docscontract/secondbrain_export_test.go`, `TestSecondbrainExportFollowsTheHistoryForm` and `TestSecondbrainExportRuleRefusesBothMismatches`, pin the history exclusion of this repository's export file; both names are in `docs/references/coverage-record.json`, so replacing them means recording that record again.

## 4. The commands and the files each one rewrote

Run in this order, each with exit 0.

| command | files it rewrote |
| --- | --- |
| `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1` | `internal/baseline/module-versions.json`, and the version lines of `context-workflow.json` (31 to the next), `core.json` (32 to the next) and `secondbrain.json` (14 to the next) |
| `make skills-sync` | the mirrors under `skills/` of the six canonical skill files edited |
| `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` | `skills/testdata/owned-skill-versions.json` and the version fields of both copies of the `SKILL.md` of `archive-spec` and `roundfix` |
| `make baseline-digests` | `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, the goldens `agent-instructions.md`, `docs-layout.md` and `secondbrain.md` under `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/`, `internal/baseline/testdata/catalog.digest`, `catalog.normalized.json` and `catalog.diagnostics.golden.json`, and the four goldens under `internal/baseline/testdata/plan-characterization/` |
| `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` | `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md` and `docs/agents/setup-context.json`; a second run printed `File changes: 0` |
| `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` | `docs/references/behavior-surfaces.json`, two fingerprints: `guide docs/user-guide/release-runbook.md` and `guide docs/user-guide/usage.md` |

One order matters. `make baseline-digests` runs `TestAuthorialSkillSync`, which fails with "canonical and distributed bytes differ" while a canonical skill file is edited and its mirror is not. On an attempt that edited the skills before `make skills-sync`, the command stopped with exit 2. A Task that edits both a module and a skill runs `make skills-sync` and the skill version record command before `make baseline-digests`; a Task that edits only modules is not affected.

The prototype changed 41 tracked files, with 351 lines added and 156 removed, before the three ADRs and before any Go code of the Spec.

## 5. Tests on the prototype

Each ran with `GOMAXPROCS=2`, after the commands of section 4, and exited 0:

- `go test ./internal/baseline -count=1`, with the five ids added to the closed list and the three tests of section 3 updated;
- `go test ./skills -count=1`: no sentence an existing test pins changed, and the `--update` section kept five steps;
- `go test -count=1 -tags docscontract ./internal/docscontract`, the glossary contracts included: `TestGlossaryPhrasesAreDefinedTerms` passes with the five entries, and `TestGlossaryLinksNameAcceptedADRs` passes once ADR-0438, ADR-0439 and ADR-0440 are in the tree, which is why they arrive with the planning change;
- `go test ./internal/cli -run 'Baseline|GuidanceComposition|Doctor|UpdateCheck|ReleasePlan' -count=1`;
- `go test ./internal/speccheck ./internal/spec -count=1`, before the skill and glossary texts;
- `make verify-docs`.

## 6. How the sync treats a missing path and an exclusion

Read in the Secondbrain at its commit `4a496f4a4`, `scripts/sync-project.sh`: the file is read one path per line, a line that starts with `!` becomes an rsync `--exclude`, a `#` line is a comment, and each include is copied with `rsync -aR`. A listed path that does not exist is skipped with the line `::warning::export path does not exist, skipping: <path>`. The mirror is then replaced with `rsync -a --delete`.

Measured on the author's machine, in a scratch folder, with the same `rsync -aR <exclusions> <include>` call over a tree with `CONTEXT.md`, `docs/references/` with a subfolder, `docs/lessons/` and two ADRs, each ADR listed by its own path. The machine's rsync is openrsync, protocol version 29; the workflow of the sync runs another rsync, which was not measured.

| exclusion line | what it dropped |
| --- | --- |
| `!docs/adr/` | nothing: an ADR listed by its own path was still copied |
| `!docs/adr/*` | both ADRs |
| `!docs/adr/0002-b.md` | that ADR |
| `!docs/references/` | the whole folder |
| `!docs/references/vendor/` | that subfolder |
| `!CONTEXT.md` | that file |

So a repository drops a standard folder or file with `!` and the item as the list writes it, and drops every ADR with `!docs/adr/*`.

## 7. Governed Paths

Measured on the planning clone with `speccheck.GovernedPath`, in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository, over 93 paths.

Governed, and changed by this Spec:

- `internal/baseline/assets/modules/context-workflow.json`, `core.json` and `secondbrain.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- the goldens `agent-instructions.md`, `docs-layout.md` and `secondbrain.md` under `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/`
- `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`
- `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`
- `.agents/skills/roundfix/references/baseline.md`, `corrections.md`, `modes.md`, `release.md` and `setup.md`
- `docs/references/coverage-record.json`

Not governed, and changed or created by this Spec: `internal/baseline/module-versions.json`, the three updated tests and the new clause test under `internal/baseline/`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`, the new Go files and tests under `internal/baseline`, `internal/spec`, `internal/cli` and `internal/releaseplan`, `internal/speccheck/glossary.go`, `internal/cli/cli.go`, `internal/docscontract/secondbrain_export_test.go`, the five mirrors `skills/roundfix/references/*.md`, `skills/testdata/owned-skill-versions.json`, the new skill contract test, `CONTEXT.md`, `.secondbrain-export`, every file under `docs/adr/` and `docs/lessons/`, the five user-guide pages and `docs/references/behavior-surfaces.json`.

Controls, governed and not changed: the `Makefile`, `go.mod`, `.roundfixrc.yml`, `.github/workflows/ci-verify.yml`, `docs/agents/specific-repository.md`, `internal/baseline/plan_test.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_human_test.go`, `skills/baseline_skill_contract_test.go` and `internal/baseline/assets/source-baselines/index.json`. Controls measured not governed and not changed: `skills-lock.json`, `docs/references/skill-coverage.json`, `AGENTS.md`, `internal/cli/baseline_update_test.go` and `internal/cli/doctor_test.go`.
