# Characterization for Spec 0417: the Baseline, a prototype, the Governed Paths and the sources (2026-11-27)

Recorded by the author of the PRD on 2026-11-27, against this repository at `8ac0224b` (v0.153.4). The Idea of the same Spec holds the sweep of the thirteen other repositories and the Secondbrain search of 2026-11-20 and 2026-11-27; this file adds what the PRD stage read, ran and measured. Nothing here was copied from another repository's code.

## 1. Decisions and where each one was read

| Date | Who | What | Where the author read it |
| --- | --- | --- | --- |
| 2026-11-20 | the maintainer | The language rule, approved with "Serve, aplique" (ai-memory message `01a1f990-ed12-7882-81a1-bf0c370fda7f`) | The Backlog Entry this Spec adopts |
| 2026-11-20 | the maintainer | The rule goes to the Baseline (message `01a1f991-12bb-7621-a985-f3456059763c`) | The Backlog Entry, which quotes it |
| 2026-11-20 | the maintainer | The five answers of the Idea stage (message `01a1f9b0-2f68-7503-817e-929d363c56e8`) | The Backlog Entry and the Idea's `## Decisions` |
| 2026-11-24 | the maintainer | The rule "Não vale no legado": the legacy Quaestor's new names follow the Portuguese of its existing code until it is shut down (decision `01a20f35-624a`) | That repository's ADR-0141, read in its checkout on 2026-11-27; the Secondbrain mirror did not hold the file yet |
| 2026-11-27 | the Praetor, the session that coordinates the software projects | This Spec enters the queue, with no batch rename, counting toward the gate-time windows without changing their measurement, and the plan returns to the coordinator before the Run (message `01a21d88-ca8e`) | The brief of the PRD stage |
| 2026-11-27 | the Praetor | The declared exception enters the scope, with two conditions (message `01a21d98-beb5`, answering the request `01a21d98-452e`) | A message of the coordinator to the author during the PRD stage |

The author listed this project's ai-memory inbox on 2026-11-27 and it held no pending message. The author read none of the messages above in the inbox.

The legacy Quaestor's ADR-0141 is titled "O código do Quaestor legado segue em português até o desligamento". It records the decision as the rule "Não vale no legado" and gives the reason for the option taken: "Um código só, legível até o desligamento, sem refactor que não ajude o operador." Its consequences keep the MCP tool parameters and descriptions in Portuguese, "como em todos os produtos da Gesttione".

## 2. What the Baseline says about language today

- `clause.core.write-generated-guidance-in-english`, the only clause of `rule.core.english-guidance`, rendered in the agent-instructions guide: "Write generated repository guidance, identifiers, headings, and examples in English. Preserve repository-authored language outside setup-owned markers."
- `language.generated`, a required decision whose only value is `English`. It renders no text.
- `clause.domain.canonical-language`, rendered in the domain guide: "Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym."
- No clause names the language of code, of infrastructure names or of what a product's users read.
- The prompt of the Pre-PR Review asks for "correctness, regressions, and security defects" and names no guide of the repository. The reviewer may open repository files, and a finding must start with a `path:line` anchor and state what breaks.
- The owned skill `qa-gate` has one sentence about language, in its frontend sweep: "product copy uses the spec's vocabulary". It says nothing about which language the copy is written in.
- The Vendored Skill `domain-modeling` owns the glossary's entry format, which has a definition and an `_Avoid_` line and no line for a name in code. It also says the glossary is "totally devoid of implementation details". Its text changes only through the skills catalog.

## 3. Prototype on a scratch clone

A clone of this repository at `8ac0224b`, outside the repository, with a throwaway authorization record that declared the regeneration commands. Nothing of it was copied into the planning change. The wording below is trial wording, not the Spec's.

### 3.1 New clauses inside the two existing rules

Four clauses added to `rule.core.english-guidance`, one added to `rule.context.domain-docs`, and `clause.domain.canonical-language` reworded. The catalog refused to load:

```text
catalog.sourceBaseline.required-clause.missing: baseline.standard-typescript-monorepo-0.0.1: clause.core.code-names-in-english; the regenerator maintains manifest rows but never creates them, so add this row first
```

One refusal for each of the five new clauses. Both rules are in the Standard TypeScript profile's required rules, so each new clause there needs a row written by hand in the Source Baseline (its manifest, its corpus and the two counts and digests of its index).

### 3.2 A new rule outside the required rules

The four clauses moved to a new rule of the `core` module, rendered in the agent-instructions guide and listed in no profile's required rules; the name in code folded into the reworded `clause.domain.canonical-language`. The catalog loaded with no Source Baseline row. Three tests then failed, as expected:

- the module version record, until the record command ran;
- the clause force characterization, which lists every clause (`internal/baseline/clause_characterization_test.go`);
- `TestStandardTypeScriptStructuralClauseRetention` in `internal/baseline/plan_test.go`, which holds `clause.domain.canonical-language` to the text the Source Baseline recorded.

The text of that clause is pinned in two more places: the fleet fixture of `internal/cli/baseline_update_test.go` and a module-edit fixture of `internal/config/regeneration_declared_test.go`.

The module record command, `make baseline-digests` and the Managed Refresh then rewrote these files:

```text
docs/agents/agent-instructions.md
docs/agents/domain.md
docs/agents/setup-context.json
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md
internal/baseline/assets/modules/context-workflow.json
internal/baseline/assets/modules/core.json
internal/baseline/assets/profiles/standard-typescript-monorepo.json
internal/baseline/module-versions.json
internal/baseline/testdata/catalog.diagnostics.golden.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json
internal/baseline/testdata/plan-characterization/clean-adoption.golden.json
internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json
internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json
```

The guides under `docs/agents/` grew from 96,009 to 97,335 bytes, 1,326 bytes, with the four clauses and the reworded one and without the sentence that names the Reader Language.

### 3.3 A required decision for the Reader Language

A required string decision, rendered as one sentence at the top of the agent-instructions guide, listed in the `core` module's required decisions and in the five built-in profiles. It changed `decisions.json`, the template of the guide, the template index, `core.json` and the five profiles.

- 188 tests of `internal/baseline` failed, almost all with `required Baseline decisions are missing: language.reader`, because their fixtures list the answers of every required decision.
- `make baseline-digests` stopped at `TestFormatterComposition`, whose answers live in `internal/baseline/plan_test.go`, and after that file gained the answer, at the plan characterization fixtures. The prototype stopped there.
- The files that list the answers of the required decisions today, found by the name of the existing decision `language.generated`: in `internal/baseline`, `plan_test.go`, `profile_alignment_test.go`, `custom_profile_test.go`, `update_test.go`, eight parity-corpus fixtures and five plan goldens; in `internal/cli`, `baseline_apply_test.go`, `baseline_human_test.go`, `baseline_incremental_verification_test.go`, `baseline_plan_test.go`, `baseline_profile_change_test.go`, `baseline_release_gate_test.go`, `baseline_update_documentation_test.go` and `baseline_update_repository_profile_test.go`.
- The user guide publishes a Decision Document with sixteen decisions, and `internal/docscontract/publicdocs_test.go` asserts that count twice.
- An update that lacks a required answer exits 3 with the state `action_required`, writes nothing and names the decision; `--decision <id=value>` answers it and `--adopt-suggested` adopts the suggested value. The code refuses a new decision that has neither a suggestion nor a default.
- A repository-owned profile that does not list a decision a template renders is refused with a message that names the decision, the template and the entry to add.
- No decision today holds a value together with the source that supports it, which the declared exception needs.

## 4. Governed Paths

Measured with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` whose test file lived outside the repository. The candidates were the files sections 3.2 and 3.3 changed or named and the files the review, the QA gate and the decision record would change.

Governed, and expected to change:

```text
.agents/skills/qa-gate/SKILL.md
.agents/skills/roundfix/SKILL.md
.agents/skills/roundfix/references/baseline.md
.agents/skills/roundfix/references/review.md
.agents/skills/setup-context-driven/SKILL.md
docs/agents/agent-instructions.md
docs/agents/setup-context.json
docs/references/coverage-record.json
internal/baseline/assets/decisions.json
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md
internal/baseline/assets/modules/context-workflow.json
internal/baseline/assets/modules/core.json
internal/baseline/assets/profiles/go-cli-tui.json
internal/baseline/assets/profiles/go-cli-typescript-monorepo.json
internal/baseline/assets/profiles/python.json
internal/baseline/assets/profiles/rust-cli.json
internal/baseline/assets/profiles/standard-typescript-monorepo.json
internal/baseline/assets/templates/guides/agent-instructions.md
internal/baseline/assets/templates/index.json
internal/baseline/plan_test.go
internal/cli/baseline_human_test.go
internal/cli/baseline_plan_test.go
internal/cli/baseline_release_gate_test.go
internal/docscontract/publicdocs_test.go
skills/qa-gate/SKILL.md
skills/roundfix/SKILL.md
skills/setup-context-driven/SKILL.md
```

Governed, and changed only if the TechSpec's design needs them:

```text
.agents/skills/roundfix/references/setup.md
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/domain.md
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
internal/baseline/assets/source-baselines/index.json
internal/baseline/repository_test.go
internal/cli/baseline_documentation_contract_test.go
internal/cli/cli_test.go
skills/baseline_skill_contract_integration_test.go
skills/baseline_skill_contract_test.go
```

The five Source Baseline files are needed only by the design of section 3.1. The other six are governed files that hold Baseline or skill contracts and that the prototype did not reach.

Not governed, by the same probe: `docs/agents/domain.md`, `CONTEXT.md`, `.secondbrain-export`, every file under `docs/adr/` and `docs/user-guide/`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `internal/baseline/module-versions.json`, the catalog files, plan goldens and parity-corpus fixtures under `internal/baseline/testdata/`, `skills/testdata/owned-skill-versions.json`, the mirrors under `skills/roundfix/references/`, `internal/cli/review.go` and its tests, `internal/cli/baseline_update_test.go`, `internal/config/regeneration_declared_test.go`, every production Go file under `internal/baseline/` and `internal/cli/`, and a new test file under `internal/baseline/`, `internal/cli/` or `internal/docscontract/`.

Two facts bear on the list. `docs/references/coverage-record.json` lists every test of every package by name, so any Task that adds, renames or removes a test changes it. And this list is the one the author expects from a prototype that stopped before the fixtures were complete: the TechSpec author finishes the prototype and the Tasks author measures again every file the Tasks declare.

## 5. Published sources

Read in full through Exa on 2026-11-27, except where noted. Each passage below was found word for word on its page.

1. William Bartlett, "When the ubiquitous language isn't English", a recap of his talk at DDD Europe 2026, posted 2026-11-02, <https://2027.dddeurope.com/blog/when-the-ubiquitous-language-isnt-english/>.
   - "His rule of thumb: business terms stay in the native language, the rest of the code stays in English."
   - "Document your naming conventions, for example with architecture decision records."
   - "If you do translate a term, keep the original alongside it in code comments or documentation, so people can search for the business term and find its English counterpart."
   - "Refactor gradually rather than renaming everything at once. Deprecate the old term, give people time to migrate, then remove it."
   - The page calls the split "A heuristic, not a rule" and conditions it: "If your business already operates in English, keep coding in English."
2. Dennis Doomen, "ALM Practices Part 9: Ubiquitous Language", 2010-08-03, <https://www.dennisdoomen.com/2010/08/alm-practices-part-9-ubiquitous.html>.
   - "If your policy is to write code in English and the Ubiquitous Language is in your native language, create a glossary that maps the verbs and nouns from the native domain model into the English translations. Be very pro-active in enforcing these specific translations, because from experience I know developers have a tendency for coming up with alternative translations."
3. "Working language - one repo, and the language each part of it speaks", no author or date on the page, <https://repositorystandards.com/docs/working-language.md>.
   - The heading "Mixing is fine. Undeclared mixing is not."
   - "An undeclared convention is not a convention, it is a habit that the newest contributor will break by accident within a week."
   - "The domain has no clean English. Regulatory terms, tax categories, legal instruments, local processes."
   - "Translating them is not neutral: it invents a vocabulary nobody uses and loses the one the business actually says out loud."
   - Its table gives "Code, identifiers, file and folder names" to English and "User-facing copy" to "the audience's language".
4. Claude Code documentation, "How Claude remembers your project", <https://code.claude.com/docs/en/claude-md>: "Claude treats CLAUDE.md files as context, not enforced configuration, so how you write instructions affects how reliably Claude follows them."
5. "Large Language Models for Multilingual Code Generation: A Benchmark and a Study on Code Quality", replication package, <https://github.com/saimaafrin/multilingual-prompt-llm-code-generation>. The abstract says "the generated code tend to contain mixes of comments and literals written in English and the language used to formulate the prompt". The page shows no venue, year or author, so it is cited as a replication package and not as a paper.
6. Golo Roden, "German technical language, English code: Where is the line drawn?", heise online, 2026-08-26, <https://www.heise.de/en/blog/German-technical-language-English-code-Where-is-the-line-drawn-11417957.html>: "A name belongs in the domain language if subject matter experts can correct it." Its example of a lost distinction is three German business words that all become `order`.
7. On an existing codebase that breaks the convention, an answer by a Norwegian developer to "How to name parts of your program as a non native English speaker", Software Engineering Stack Exchange, <https://softwareengineering.stackexchange.com/questions/116242/how-to-name-parts-of-your-program-as-a-non-native-english-speaker>: "For mixed-language projects that break the above principles, I just try to extend whatever principles that project uses, as consistency is more important than the above principles." It is a community answer, read as the highlighted answer and not as the whole thread.
8. Read in part: Florian Grabs, a bachelor thesis of TU Chemnitz of 2022 on the language of identifiers and code comprehension, 71 participants, which found no significant advantage for identifiers in a language other than English; and Pawelka and Jürgens, "Is This Code Written in English?", 2015, which found comments and identifiers in more than one language in a significant part of the industry projects it analysed and in none of the open-source ones. Neither measured defects.
9. Tools, read from their README pages only: `glosslint` checks names against a glossary of terms, and `eslint-plugin-only-english-identifiers` flags identifiers outside an English dictionary and takes a list of allowed words. No tool was found that holds identifiers to a glossary that maps a native term to its name in code.

What the sources changed here: sources 1, 3 and 6 support the split and the Untranslated Term; source 2 supports recording the name in code of a translated term; sources 3 and 4 support a rule the repository states and a check where a diff is read; source 7 supports the declared exception for an existing codebase, with the weight of a community answer; sources 8 and 9 confirm that no measurement and no ready tool justify a deterministic check in this version. No source argues that an existing codebase in another language should take English for its new names while the old ones stay.

Not found: a measurement of Agents obeying a language rule over time, and a measurement of defects in mixed-language code.

## 6. Secondbrain

Searched on 2026-11-27 through the Indexus tool `pesquisar_acervos`, after `wiki/index.md`, with two questions: the language of code in a repository whose domain is Portuguese, and the legacy codebase that keeps Portuguese. The answers held the acervos `geral`, `gesttione` and `tecnico`.

- `sources/digest-roundfix-idioma-do-codigo-em-dominio-nao-ingles-2026-11-20.md`: the digest of this Spec's own Idea research. The same research, not a second source.
- `raw/web/akita/2008-07-31-akita-tradu-o-por-que-voc-n-o-deve-codificar-em-portugu-s.md`, opened with `abrir_conteudo`: a 2008 translation, by a Brazilian engineer, of an article that argues for English in variable, function, class and file names and in database objects, and for a professional translator for "as mensagens mostradas ao usuários". It supports the first two parts of the rule and gives one reason for the exception to weigh: mixing a human language with the language's own English functions leads to "confusões e duplicação de código".
- `concepts/bounded-contexts-e-ddd.md` and `projects/oraculum/mirror/docs/adr/0005-nomes-de-ferramentas-em-portugues.md`: returned again, and already read by the Idea.
- No page holds the decision about the legacy codebase. The mirror of that repository ends one decision record before the one that holds it.

## 7. The adopters and what each message says

From the Idea's sweep of 2026-11-27, with the coordinator's decision of the same day applied. Eleven repositories hold a Setup Manifest.

| Repository | Reader Language the maintainer already gave | What its own rules hold |
| --- | --- | --- |
| Vortex | pt-BR | Its own rules repeat the rule, at lines 386 to 389 of its repository rules |
| The new Quaestor | pt-BR | Its ADR-0010 adopted the rule and left the glossary folder to the Baseline's owner; it uses a repository-owned profile |
| The legacy Quaestor | pt-BR | Line 3 of its repository rules states Portuguese for code, and its ADR-0141 records the owner's decision: it declares the exception |
| Oraculum | pt-BR | Its code is Portuguese and it states no rule for code; nothing is renamed |
| Argus, Fluxus, Fiscus, Indexus, Conexus | pt-BR | Nothing that repeats or contradicts the rule was found |
| Pantheon, the Secondbrain | none named | Each answers when its update asks |

The skills catalog and Praetor hold no Setup Manifest and receive no message.

## 8. The TechSpec stage's prototype

Recorded by the author of the TechSpec on 2026-11-27, on a second scratch clone at `8ac0224b`, outside the repository. It carried the design of `_techspec.md`: the new rule with its gated clauses, the required decision `language.reader`, the optional decision `language.code_exception` with its type and the `present` gate, the reworded canonical-terms clause, the template line and the review sentence. Nothing of it was copied into the planning change except the clause texts, the refusal reasons and the two transcripts.

- `internal/baseline`: 198 tests failed with the required decision alone. They passed after the answer was added to the helpers of `plan_test.go` and `profile_alignment_test.go`, the new clauses to `clause_characterization_test.go`, a skip of the optional decision to the helper of `update_test.go`, and the insertion to the clause list of `TestStandardTypeScriptStructuralClauseRetention`. No parity-corpus fixture changed.
- The template line placed after `{{branch.prefix}}` failed the two sentence tests of `branch_prefix_optional_test.go`; placed before it, both passed unchanged.
- `make baseline-digests` rewrote the Standard TypeScript profile's `goldenDigest`, the two formatter goldens, the three catalog files and five plan goldens, the fifth being `unsatisfied-blocking-capabilities.golden.json`.
- `internal/cli`: 19 tests failed. Seventeen are in `baseline_human_test.go` and `baseline_release_gate_test.go` and wait for one more scripted answer of the interactive flow; one is the message of `baseline_profile_change_test.go`; one is the clause line of the fleet fixture in `baseline_update_test.go`; and before the interactive flow skipped the exception, every interactive test failed with `unsupported type "code-language-exception"`. The prototype did not extend the scripted answers.
- The Managed Refresh with `--decision language.reader=English` changed `docs/agents/agent-instructions.md`, `docs/agents/domain.md` and `docs/agents/setup-context.json` and no other file, a second update planned nothing, and repeating the recorded answer was no change. The guides grew from 96,009 to 97,442 bytes, 1,433 bytes.
- With the Decision Document, its count, the two clause-text fixtures and the surface record updated, `make verify-contracts` passed: 144 Repository Contract Tests.
- A declaration with a language and a source rendered the exception's sentence and its clause in place of the three gated clauses; one without a source exited 2 and wrote nothing.
- The probe of section 4, run again over 89 candidate files, classified as governed the 28 files of the PRD's list among the files the design changes, and none of the eleven reserve files is changed. Not governed and changed: `internal/baseline/catalog_validate.go`, `plan.go`, `custom_profile.go`, `project_decision_render.go`, `clause_characterization_test.go`, `profile_alignment_test.go`, `update_test.go`, `internal/cli/baseline_human.go`, `review.go`, the other `internal/cli` test files of API Contract 5, `internal/config/regeneration_declared_test.go`, `internal/docscontract/command_documentation_test.go`, and a new test file under `internal/baseline`, `internal/cli` or `skills`.

Not done by the prototype: the new tests, the skills, the user-guide pages beyond the Decision Document, the glossary and the ADR.
