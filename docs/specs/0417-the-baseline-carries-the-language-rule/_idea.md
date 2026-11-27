---
spec: 0417-the-baseline-carries-the-language-rule
created: 2026-11-27
status: promoted
leverage: compounding
---

# The Baseline carries the language rule for code, data and products

## Overview

**Problem.** On 2026-11-20 the maintainer decided one language rule for every project: English for code and infrastructure, the reader's language for what a person reads, Brazilian fiscal and legal terms kept in Portuguese inside English code and entered in the repository glossary, and no bulk rename. The rule lives in a user-level instruction file and in ai-memory messages. The Context-Driven Baseline says nothing about the language of code, and one of its clauses pushes the other way. In the week since the decision, two repositories wrote the rule into ADRs of their own, one to adopt it and one to record that it does not apply there, and a third cites the user-level file by its path.

**Who it is for.** Every Agent that names something in a repository that takes the Baseline, the maintainer who reviews the result, and the people who read the product: Brazilian business users, who read Portuguese, and the users of Roundfix and of the skills catalog, who read English.

**Why it is valuable.** Fourteen repositories were read, this one and thirteen others. Two of them write their code in Portuguese and a third mixes the two languages. Six state a language rule of their own, in six different texts, and one of the six states the opposite of the decision. A rule the Baseline carries reaches each repository on its next Baseline update and replaces those texts with one.

**Ambition of V1.** A complete feature of clause size: the rule as Normative Clauses every adopter reads, the language of the reader as each repository's own answer, the glossary the Baseline already has as the record of each term's name in code, and a check where a diff is already read, the Pre-PR Review and the QA gate. No new tool and no rename of existing code.

This Idea calls the language a repository declares for its readers the Reader Language, and a business term that keeps its original form inside an English name an Untranslated Term. Both names are candidates for the glossary, which the PRD settles.

## Problem

An Agent in a repository that takes the Baseline reads the Baseline guides before it names anything. Those guides hold one sentence about language, `clause.core.write-generated-guidance-in-english`, which covers only what the Baseline generates. They hold one sentence about names, `clause.domain.canonical-language`: "Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes."

In a repository whose glossary is written in Portuguese, an Agent that obeys that clause writes Portuguese identifiers. Oraculum and the legacy Quaestor hold exactly that: a table named `entregas_atualizacoes`, a module named `apuracao_icms.py`. The new Quaestor held it until 2026-11-20, when a Spec of its own moved the code to English: the port that this Idea's first sweep read as `executor-de-codigo.ts` is `code-executor.ts` today.

The decision lives outside the repositories. A user-level instruction file is not a rule of the repository (ADR-0436): a collaborator, another machine and a harness that does not load that file never see it. The maintainer already refused that home for a rule the flow depends on: "Tudo tem que funcionar por configurações no repositório" (ADR-0431). Praetor's own rules now give the language rule in two sentences and name the user-level file, by its path, as their source.

Each repository fills the gap alone. Six state a rule of their own, Praetor's among them, and no two of the texts are the same:

- The new Quaestor wrote the rule into its ADR-0010 on 2026-11-20 and moved its code to English the same day. The same ADR records that it created no `docs/glossary` folder, because the Baseline fixes `CONTEXT.md` as the glossary, and that the question "fica com o dono do Baseline".
- The legacy Quaestor's own rules say "Idioma do código, dos comentários e da documentação: português do Brasil, sem acentos em identificadores", the opposite of the rule. Its ADR-0141 records why the line stays: on 2026-11-24 the maintainer decided that the rule "Não vale no legado", a codebase kept in maintenance until it is shut down.
- Vortex says code, comments, identifiers, migrations, API contracts and structured data are English and the domain vocabulary is pt-BR, which is the rule in other words.
- The skills catalog says every tracked file is English.
- Oraculum decided in its ADR-0005 that MCP tool names, parameters and enums are Portuguese, and wrote nothing about its code, which is Portuguese.
- Fluxus states no rule; its glossary says canonical domain terms are English even when the product is discussed in Portuguese.

Nothing checks a name. The Pre-PR Review prompt asks for "correctness, regressions, and security defects" and names no repository rule. The QA gate's one sentence about language asks that product copy use the Spec's vocabulary, and says nothing about which language the copy is written in. No command reads an identifier. The maintainer talks to his Agents in Portuguese, and a published study of GPT, Claude and DeepSeek found that code generated from a non-English prompt tends "to contain mixes of comments and literals written in English and the language used to formulate the prompt".

The rule as first written named two things the Baseline does not have. It named `docs/glossary`, a path none of the repositories read has: nine of the thirteen others keep their glossary in a root `CONTEXT.md`, Fiscus keeps one per context under a root `CONTEXT-MAP.md`, and three have none. Those are the glossary files `roundfix spec check`, the QA gate, the Archive Command and the Standard Export List read.

It also fixed Portuguese as the language a person reads, while the Baseline is public and two of the maintainer's own repositories, Roundfix and the skills catalog, are English products. The maintainer's answers of 2026-11-20 settled both, and `## Decisions` quotes them.

### Evidence

Codebase, this repository at `8ac0224b`, read on 2026-11-27:

- `internal/baseline/assets/modules/core.json`: `rule.core.english-guidance` and its only clause; `internal/baseline/assets/decisions.json`: `language.generated`, whose only value is English.
- `internal/baseline/assets/modules/context-workflow.json`: `clause.domain.canonical-language` and `clause.domain.glossary-currency`.
- `internal/speccheck/glossary.go`, `GlossaryFiles`: the glossary is the root `CONTEXT.md`, `GLOSSARY.md`, `CONTEXT-MAP.md`, `GLOSSARY-MAP.md` and each `CONTEXT.md` or `GLOSSARY.md` a map links (ADR-0244 holds a Spec's terms to it).
- `internal/baseline/docs_layout.go`: a folder under `docs/` outside the listed ones, and a file directly under `docs/`, is reported as `DR-DOCS-OUTSIDE-LAYOUT` unless the repository's own rules name it.
- `internal/cli/review.go`, `buildReviewPrompt` and `reviewSpecContextInstruction`: the review judges correctness, regressions and security, and reports a choice that contradicts a decision the Spec records.
- `.agents/skills/qa-gate/SKILL.md`, the frontend sweep: "product copy uses the spec's vocabulary" is the gate's only sentence about language.
- `docs/references/0001-ubiquitous-language-in-colocated-context-md.md`, an early Baseline ADR the catalog no longer ships: "Canonical domain terms use the language spoken by the domain experts. When the code uses a different English identifier, `CONTEXT.md` records that identifier as an alias rather than replacing the domain term. Standard legal, regulatory, or protocol vocabulary retains its authoritative form." It rejected a central glossary folder because glossaries kept apart from the code drift.
- ADR-0118 (a repository constant that varies becomes a decision), ADR-0222 (a Baseline guide says only what holds for the repository that reads it), ADR-0431 and ADR-0436.
- Archive Records 0409 (three implementation Tasks and the QA gate for one group of clauses and five skills) and 0410 (five and the QA gate for three modules, two skills and three ADRs), which size this work, and 0415 and 0416, the latest deliveries that regenerated the derived Baseline files.
- The guides under `docs/agents/` hold 96,009 bytes and `CONTEXT.md` 208,744, against 78,826 and 182,265 measured on 2026-11-15 (`docs/references/2026-11-15-usage-by-context-and-session-identity-measurement.md`).

Other repositories, read on 2026-11-27 in their checkouts, and Conexus in its Secondbrain mirror: the new Quaestor, the legacy Quaestor, Oraculum, Vortex, Argus, the skills catalog, the Secondbrain, Fluxus, Fiscus, Pantheon, Indexus, Praetor and Conexus. The first sweep was of 2026-11-20. This one read again each repository's own rules, its glossary files and a sample of file names; it is a sample, not an audit.

- Eleven of the thirteen hold a Setup Manifest, so the rule reaches them through `roundfix baseline update`. The skills catalog and Praetor hold none.
- Code in Portuguese: Oraculum (`leitor-calibracao.ts`, the table `entregas_atualizacoes`) and the legacy Quaestor (`apuracao_icms.py`). Indexus mixes: `gerador-de-vetores.ts` beside `search-fiscal-corpus.use-case.ts`.
- A language rule of their own: the new Quaestor (`docs/adr/0010-codigo-em-ingles-e-superficie-do-mcp-em-portugues.md`), the legacy Quaestor (`docs/agents/specific-repository.md`, line 3, and `docs/adr/0141-codigo-do-legado-segue-em-portugues.md`), Oraculum (`docs/adr/0005-nomes-de-ferramentas-em-portugues.md`, the MCP surface only), Vortex (`docs/agents/specific-repository.md`, lines 386 to 389), the skills catalog (`AGENTS.md`, lines 11 and 32) and Praetor (`AGENTS.md`, lines 94 and 95).
- The name in code is already written into two glossaries. Vortex records the English identifier as an alias of the Portuguese term (`income-statement`, `buyer`), by the early Baseline ADR it still holds. The new Quaestor's `CONTEXT.md` writes "em código, `buildRelationGap`" beside a Portuguese term and lists its kept fiscal terms (`escrituracao`, `competencia`, `finalidade`) in a section of their own.
- No repository has a `docs/glossary` path.

Secondbrain, searched through the Indexus tool `pesquisar_acervos` on 2026-11-20, and again on 2026-11-27 after a search of `wiki/index.md` for the subject:

- `wiki/sources/branas-formacao-arquitetura-ddd.md`, line 55, and `wiki/concepts/bounded-contexts-e-ddd.md`, line 48: proper names are not translated (boleto, PIX, PIS, COFINS); a Portuguese domain in code needs discipline and English code risks a wrong translation. It supports the exception and names the risk the glossary answers.
- `wiki/sources/conexus-alinhamento-com-o-diego-2026-10-20.md`, line 20: Specs in Portuguese and code in English, because translated fiscal and banking terms confuse. It shows the rule was already practice in one project.
- `projects/oraculum/mirror/docs/adr/0005-nomes-de-ferramentas-em-portugues.md`: MCP tools, parameters and enums in Portuguese, because the questions arrive in Portuguese and "A precisão de seleção de ferramentas depende mais da qualidade das descrições do que do idioma do nome." It supports the MCP sentence of the rule.
- `wiki/sources/digest-roundfix-idioma-do-codigo-em-dominio-nao-ingles-2026-11-20.md` is the digest of this Idea's own research, which the Secondbrain took in on 2026-11-25. It is the same research, not a second source.
- No page states a global language rule, a glossary of kept terms, `docs/glossary` or the decision about the legacy codebase.

Published sources, read through Exa on 2026-11-20 and not read again for this version:

- William Bartlett, "When the ubiquitous language isn't English", DDD Europe 2026, <https://2027.dddeurope.com/blog/when-the-ubiquitous-language-isnt-english/>: "business terms stay in the native language, the rest of the code stays in English"; document the convention; keep the original beside a translated term; refactor gradually. It supports the split and the rename rule.
- "Working language", <https://repositorystandards.com/docs/working-language.md>: code, identifiers and file names in English, user-facing copy in the audience's language, "Mixing is fine. Undeclared mixing is not", and tax terms have no clean English. It supports writing the rule where the Agent reads first, and it draws the boundary by artifact type, which is the adjacent alternative of the Opportunity Scan.
- Dennis Doomen, "ALM Practices Part 9: Ubiquitous Language", <https://www.dennisdoomen.com/2010/08/alm-practices-part-9-ubiquitous.html>: with English code and a native domain language, "create a glossary that maps the verbs and nouns from the native domain model into the English translations", because developers "have a tendency for coming up with alternative translations". It changed the glossary proposal from kept terms only to the name in code of every term.
- Claude Code documentation, <https://code.claude.com/docs/en/claude-md>: an instruction file is "context, not enforced configuration". It challenges a rule nobody checks.
- "Large Language Models for Multilingual Code Generation", <https://github.com/saimaafrin/multilingual-prompt-llm-code-generation>, abstract: the drift quoted above.
- Read only as search excerpts, not in full: the Model Context Protocol tool-name character set (letters, digits, `_`, `-`, `.`), Martin Fowler on opportunistic refactoring and on mass renames, cspell and SQLFluff as spelling and naming checks, and a study of LLM judges on non-English text.
- Not found: a measurement of Agents obeying a language rule over time, a measurement of tool selection with Portuguese tool names, and a linter that checks identifiers against a bilingual glossary.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | English for code and infrastructure | Critical | Every adopter's guides state that variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs are written in English. |
| F2  | The Reader Language for what a person reads | Critical | Screens, user messages, e-mails, PDFs and the product's documentation for its users are written in the language the repository declares for its readers: pt-BR in the Gesttione products, English in Roundfix and in the skills catalog. MCP tool names, descriptions and parameters follow the Reader Language. Each repository answers once, at its next Baseline update. |
| F3  | Untranslated Terms | Critical | A term of the business with no faithful English translation keeps its original form inside an English name. An official acronym stays as it is (`icms_amount`, `efd_file`); a legal institute is written without accents (`apuracao_summary`); the rest of the name is English. Brazilian fiscal and legal terms are the named case. |
| F4  | The glossary records the name in code | High | Each Untranslated Term enters the repository's glossary, the one the Baseline already defines. A term whose name in code differs from the term itself records that name: the English translation chosen (Filial, `branch`) or the unaccented form kept (Apuração, `apuracao`). |
| F5  | The domain clause agrees with the rule | High | "Canonical domain terms in code names" comes to mean the name the glossary records for the term, so a Portuguese glossary no longer produces Portuguese identifiers. |
| F6  | New names only | High | A new name follows the rule. An existing name changes when its part of the code is touched, a large change takes its own Spec, and nothing is renamed in bulk without a request. |
| F7  | Review and QA report what breaks the rule | Medium | The Pre-PR Review reports a new identifier, file or database name in the candidate that breaks the rule. The QA gate reports text a person reads in the wrong language. Both are findings like any other, and no new tool reads a name. |
| F8  | Adopters are told | Medium | The rule arrives with each adopter's next Baseline update. Each one receives one message that names the Reader Language the maintainer already gave for it, and, where its own rules repeat or contradict the rule, the line. |

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Built-in Baseline Profiles whose rendered guides carry the rule | 5 of 5 | A catalog test renders every built-in profile and asserts the rule's clause identities. |
| Repositories that take the Baseline and whose Repository-Specific Normative Rules repeat the rule, or contradict it without a recorded decision of their owner | 0 within 30 days of the release that carries it, from 1 today (Vortex repeats it; the legacy Quaestor contradicts it by a recorded decision) | Repeat this Idea's sweep of each repository's `docs/agents/specific-repository.md`, in the checkouts and the Secondbrain mirrors. |
| New accented identifiers, file names and database names in the first three Specs each adopter under the rule delivers after its update | 0 | Filter the added paths, and the names on the added lines, of each delivery's diff with a non-ASCII pattern. |
| Untranslated Terms in those Specs that the glossary lists | 100% | Compare the Portuguese words in added names with the glossary files `roundfix spec check` resolves. |
| Growth of the guides every Authoring Phase reads | At most 1,500 bytes, under 2% of the 96,009 bytes measured on 2026-11-27 | `wc -c` over the rendered guides under `docs/agents/` before and after. |
| Language findings per Spec in the adopters, after each one's first five Specs under the rule | Fewer than 1 | Count the findings that cite the rule in each Spec's review record and QA Report. |

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Strong: one text replaces six local ones and a pointer to a user-level file, and rewords a clause that works against the decision. It adds no capability by itself. |
| Reach           | What % of users would this affect?                  | Must do: every adopter and every Agent session reads the guides the rule joins. |
| Frequency       | How often would users encounter this value?         | Must do: every new identifier, file, table and message of every Task. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe: a published standard draws the same split by artifact type. The tie between Untranslated Terms and the glossary the Spec check already holds is this product's own. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong when the glossary records the name in code: each entry settles one translation for every later session. A rule text alone is easy to copy. |
| Feasibility     | Can we actually build this?                         | Strong: clause work with two precedents of its size, 0409 and 0410. The known costs are the Source Baseline rows, the derived Baseline files, and one answer per repository for the Reader Language (ADR-0118). |

Leverage: compounding, through the glossary. Without F4's name in code it is a quick win. Verdict: proceed, reshaped from the literal text into a portable obligation with a repository's own answer. Two results would change it: language findings that stay at one or more per Spec after an adopter's first five Specs, which means the text is not read and V1 needs the check it left out; or glossaries that gain no name in code in the first Specs with Untranslated Terms, which means F4 is a line nobody writes.

## Trade-off Insights

- **Recommended approach:** the rule as Normative Clauses in the guide every Agent already reads first, the Reader Language as each repository's answer, the existing glossary as the record of each term's name in code, and the Pre-PR Review and the QA gate as the check. Where the clauses sit, a new guide or a new module, is the TechSpec's decision; the maintainer's message left the place to this repository.
- **Key trade-offs:**
  - *Literal text against a repository's answer.* The literal text costs nothing in the fleet and is word for word what the maintainer first approved. But a universal Normative Clause cannot be weakened by a repository's own rules, so Portuguese would bind Roundfix's CLI, README and user guide, and every outside adopter. ADR-0118 names this case: the obligation is portable, the value is one maintainer's answer. The maintainer chose the repository's answer.
  - *`docs/glossary` against the glossary that exists.* A short list of kept terms is easier to scan than a long `CONTEXT.md`. But it would be a second glossary beside the first, unknown to the Spec check, the QA gate, the archive, the Standard Export List and the documentation layout, and the early Baseline ADR already rejected a central glossary folder because it drifts from the code. The maintainer chose `CONTEXT.md`.
  - *Kept terms only against the name in code of every term.* Recording only the kept terms is what the rule first asked. But the larger risk of English code over a Portuguese business is the translated term: without a recorded choice, each session translates "filial" its own way. The cost is one line for each term a Spec introduces or touches, never a sweep of the existing glossary.
  - *A written rule against a checked rule.* Text is cheap and an instruction file gives no guarantee. A deterministic check is the only strong one, and none exists for a bilingual glossary. The review and the QA gate cost one instruction each and are weak on non-English text. The maintainer chose the review and the QA gate, with no new tool.
- **Risks identified:**
  - The rule reaches the legacy Quaestor, where the maintainer decided on 2026-11-24 that it does not apply. That repository takes the Baseline and updated it three times between 2026-11-15 and 2026-11-20. As approved, F1 would put "English for code" in its guides beside an ADR that says the opposite. This is the first open question, and it is the maintainer's.
  - A new required decision stops each adopter's next Baseline update until it is answered (ADR-0118). Every adopter pays that step once, and F8's message carries the answer the maintainer already gave.
  - ADR-0431 records that an Agent reads a guide only when its pointer leads there. A rule in a guide of its own is read less surely than one in the guide the root file already makes mandatory.
  - An optional decision with a default would give an unanswered repository the wrong Reader Language in silence.
  - The rule reaches Oraculum, whose tree is Portuguese and which states no rule for its code. "Changes only when touched" is the only thing that keeps an Agent from renaming it; the clause must say it as plainly as the rest.
  - An Agent may read "documentation for people" as Specs and ADRs. The maintainer's fifth answer keeps them out, and the clause must say "the product's users".
  - Neither list of the rule holds every name: the keys and values of an MCP response, an error message to an operator, a command-line flag. The new Quaestor decided those for itself in its ADR-0010; another repository will meet them without a line to read.
  - Every byte added to the guides is read by every Authoring Phase. The guides grew from 78,826 to 96,009 bytes in twelve days, and the open Backlog Entry of 2026-11-16 is trying to shrink that floor.
  - Two open Backlog Entries of 2026-11-24 give adopters a reason to hold their Baseline update. Vortex is not running `roundfix baseline update` until this repository answers its report about the 2 MB limit of the release plan's Baseline check, and Pantheon reported that the update empties `amended_by` in ADRs named by date. The rule reaches a repository only when it updates.
  - No Spec is active in this repository today. If another Baseline Spec is planned before this one is delivered, both regenerate the same derived Baseline files and their order must be stated.
- **Dissenting view:** the Reader Language as a decision is generality for adopters nobody has met, and the maintainer praised the text as it stood ("Ficaram muito boas desse jeito aí"). The answer to it is that the second repository already exists: Roundfix reads this Baseline and is an English product. A second dissent: the name in code of every term is more than was first asked, and a glossary nobody keeps current is worse than none; the Spec check already holds the glossary to each Spec, which is why the cost falls on new terms only.
- **Stretch goal (V2+):** a deterministic check that fails a delivery whose diff adds an accented name or a non-English word the glossary does not list; and a declared language for each kind of engineering document.

## Opportunity Scan

**Original.** The rule's text as first decided: Portuguese and Brazilian fiscal terms fixed, a guide of its own, `docs/glossary` as a Baseline item in repositories with a fiscal domain, and a check "if it makes sense". Its ceiling is consistency in the maintainer's Portuguese products. It is a one-time text, and it is wrong for English products. Two rules stand beside it: ADR-0118 beside Portuguese for every adopter, and ADR-0244 with the documentation layout clause beside `docs/glossary`.

| Alternative | What | Why powerful | Scale | Impact | Reach | Frequency | Differentiation | Defensibility | Feasibility |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1, more ambitious | The glossary is the naming authority and a command holds new names to it: every term records its name in code, and a delivery fails on an accented name or a non-English word the glossary lacks. | The rule stops depending on anyone reading it. | Massive: a tokenizer for each language, dictionaries, false positives; no ready tool was found. | Strong | Must do | Must do | Strong | Strong | Maybe |
| 2, simpler | Two clauses and nothing else: the English list and the rename rule join the existing language rule, "what a person reads follows the language of the product's users", kept terms go to the existing glossary. No decision, no check. | Ships in two Tasks and reaches every adopter with no answer to give. | Small | Maybe | Must do | Must do | Pass | Maybe | Must do |
| 3, adjacent | Each repository declares the language of each kind of artifact: code (English, fixed), text for readers, Specs and ADRs, commits and Pull Requests. | Settles a wider question the adopters also answer differently today. | Medium to large: it reaches the authoring skills. | Strong | Must do | Must do | Maybe | Maybe | Maybe |

**Chosen direction:** hybrid — the maintainer's answers 2 to 4 of 2026-11-20 are each the recommended option, and together they are this direction: the original's content, the Reader Language as a repository's answer, the existing glossary with the name in code from alternative 1 as text only, and the review and QA instruction. Alternative 1's command and alternative 3's table are the stretch goal. His fifth answer keeps alternative 3's wider question out of V1.

## Integration with Existing Features

| Integration point | How |
| --- | --- |
| `rule.core.english-guidance` and `language.generated` | The Baseline's existing language rule and decision; the new clauses sit beside them and leave the language of generated guidance as it is. |
| `clause.domain.canonical-language` | Reworded so the name in code is the one the glossary records (F5). |
| Glossary Declaration and the Spec check | A Spec that introduces an Untranslated Term declares it like any other domain term; no new gate. |
| Pre-PR Review and QA gate | One instruction each (F7). |
| Baseline update and Update Check | Carry the rule to each adopter; nothing new. |
| Repository-Specific Normative Rules | An adopter's own line that repeats or contradicts the rule is the adopter's to remove; Roundfix sends the message (F8). |

## Out of Scope (V1)

- **Renaming existing code, tables or files in any repository** — the rule itself forbids a bulk rename, the coordinator repeated it on 2026-11-27, and each adopter owns its code.
- **A deterministic check of names** — the maintainer chose the review and the QA gate with no new tool. No ready tool exists, the false positives are unmeasured, and the first Specs under the rule give the numbers that would justify one.
- **The language of Specs, ADRs, agent guides, commits and Pull Requests** — the maintainer's fifth answer: the rule covers what the product's user reads, and each repository keeps what it does today.
- **A glossary starter file or a `docs/glossary` folder** — the maintainer's third answer: the term stays in `CONTEXT.md`.
- **Renaming the legacy Quaestor or changing the language of its new names** — the maintainer decided on 2026-11-24 that the rule does not apply there, as its ADR-0141 records. How the Baseline leaves it out is the first open question.
- **Translating Roundfix's own output or documentation** — Roundfix declares English, and its readers read English.
- **Editing an adopter's own rules from this repository** — a message to its owner is the route.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.94)
- unverified_assumption: no (0.28)
- market_question: no (0.07)
- verdict: indicated
- stage: ran
- shape: changed — from the literal text, with Portuguese for every repository and a `docs/glossary` folder, to a portable rule with a Reader Language each repository declares and the existing `CONTEXT.md` glossary recording each term's name in code.

The judge read the request on 2026-11-20, while it existed only in the session. Triage recorded it on 2026-11-24 as [the Backlog Entry this Spec adopts](references/2026-11-24-the-baseline-carries-the-language-rule-for-code-data-and-products.md), and this stage did not run the judge again. The verdict was indicated, and the maintainer's first answer of 2026-11-20 was to proceed with the Idea. The research agrees with the judge that the shape was open. It disagrees mildly on the second criterion: the request rested on one assumption that did not hold, that every reader of every adopter reads Portuguese, and named one path no repository has.

## Rule challenges

None.

No rule blocked the goal. `clause.domain.canonical-language` changes with this Spec (F5), and no Rule Challenge was asked about it: the change follows from the maintainer's decision of 2026-11-20, message `01a1f990-ed12-7882-81a1-bf0c370fda7f`, which sets English for code names. The chosen direction stays inside ADR-0118, ADR-0244 and the documentation layout clause. The two shapes those rules stand beside, Portuguese for every adopter and a `docs/glossary` folder, are the ones the maintainer's answers 2 and 3 did not take.

## Decisions

The maintainer's decisions come first. The first three reached this project on 2026-11-20 as ai-memory messages that the Jarvis project relayed, and the Backlog Entry records them. The text in quotes is verbatim. The author did not read the messages: the Backlog Entry and the brief that started this stage hold them, and the inbox of this project held no pending message on 2026-11-27.

- 2026-11-20, message `01a1f990-ed12-7882-81a1-bf0c370fda7f`: the language rule for every project, approved with "Serve, aplique".
- 2026-11-20, message `01a1f991-12bb-7621-a985-f3456059763c`: "Talvez, até mesmo essas regras aí de idioma devam ir para a roundfix adicionar na baseline. Ficaram muito boas desse jeito aí." As the Jarvis project read it, the message leaves the place in the Baseline to this repository, an agent guide, a language module or whatever fits best, and gives this Spec its priority, after 0415 and 0416.
- 2026-11-20, 15:28 BRT, message `01a1f9b0-2f68-7503-817e-929d363c56e8`, which answers the request `01a1f9a2-fa98-7992-a272-d5c19e3f9701`: the five answers to this stage's questions. The text after each dash is the option's description as the Backlog Entry records it.
  1. Proceed with the Idea.
  2. "Cada repositório declara" — code in English for all; the Gesttione products declare pt-BR for the reader; Roundfix and skills declare English; each repository answers once at its next Baseline update.
  3. "No CONTEXT.md, com o nome no código" — the term stays in `CONTEXT.md` and records the name used in code, translated (Filial, `branch`) or kept (Apuração, `apuracao`); no `docs/glossary` folder.
  4. "Review e QA" — the review flags a new name outside the rule in the diff and the QA gate flags user-facing text in the wrong language, as any finding, with no new tool.
  5. "Não" — the rule covers what the product's user reads (screens, messages, e-mails, PDFs, manual, MCP tools); Specs, ADRs, agent guides, commits and Pull Requests stay as each repository already does.
- 2026-11-27, 12:30 BRT, the Praetor, the session that coordinates the software projects, message `01a21d88-ca8e`: this Spec enters the queue now, with no batch rename of old code. New names follow the rule, and an old name changes only where the Spec touches it. The author did not read this message; the stage's brief relays it and gives the id in this short form.
- 2026-11-27, the Praetor, message `01a21d98-beb5`, which answers the request `01a21d98-452e`: the first open question below is answered with its option 1, the declared exception. The decision is the coordinator's and rests on the maintainer's decision of 2026-11-24, which the legacy Quaestor's ADR-0141 records. Two conditions come with it: the exception holds only when the declaration cites the source of the owner's decision, and it covers existing code only. The same message settles the second open question: each repository, Pantheon and the Secondbrain included, answers its Reader Language when its own Baseline update asks. This line was added in the PRD stage, whose brief relays the message; the author did not read it. The PRD's `## Decisions` holds the record.

The author's own, made in this stage:

- A draft of this Idea was written on 2026-11-20 with the recommended option of questions 2 to 5 assumed. The maintainer's answers are those options, and the brief that started this stage treats that draft as approved in that shape. The questions were not asked again. They went to the maintainer as the request `01a1f9a2-fa98-7992-a272-d5c19e3f9701`, because an Author Subagent cannot ask.
- The chosen direction is the hybrid, because answers 2 to 4 are its three parts.
- `clause.domain.canonical-language` changes with this Spec, as `## Rule challenges` records.
- The facts were read again on 2026-11-27 and four of them moved. The new Quaestor adopted the rule in an ADR of its own and its code is English. The legacy Quaestor is outside the rule by the maintainer's decision of 2026-11-24. The QA gate has one sentence about language, where the draft said it had none. The guides hold 96,009 bytes, where the draft counted against 78,826. The Problem, the Evidence, the second and fifth KPIs, the risks and the open questions follow those facts.
- F8 gained one thing from the maintainer's second answer: the message names the Reader Language he already gave for the repository, so its update does not wait for a question he has answered.
- The Backlog Entry is not adopted at this stage. The PRD's adoption step sets its status, moves it into `references/` and rewrites the link in `## Idea-stage judgment`.
- The other open Backlog Entries stay where they are. The nearest is "A rule or ADR must not block the project", also a Baseline rule the maintainer asked for on 2026-11-20: it changes the Rule Challenge contract and the authoring skills, where this Spec changes the language and domain clauses, the review and the QA gate.
- No external research was done for this version, so no second digest was sent. The digest of 2026-11-20 is ai-memory message `01a1f99d-8340-7212-8fdc-6aea3bbcd996`.
- No ADR is written at this stage.

## Open Questions

The first question is the maintainer's. The author returned it to the session that started this stage and did not answer it.

Answered on 2026-11-27, after this stage: the first question with its option 1 and the second with its default, as the last dated line of `## Decisions` records. Both stay below as they were asked.

- How does the legacy Quaestor stay outside the rule once the Baseline carries it? Two answers of the maintainer meet there. On 2026-11-20 he chose "Cada repositório declara", with code in English for all. On 2026-11-24 he decided that the rule "Não vale no legado": the legacy Quaestor's new names follow the Portuguese of its existing code until it is shut down. The author read that decision in the repository's ADR-0141, which cites it as `01a20f35-624a`, and not in the message. A repository's own rules cannot weaken a universal Normative Clause, and that repository takes every Baseline update. The author recommends the first of two options.
  1. The rule admits a declared exception. A repository whose owner decided that an existing codebase keeps its language records that decision; its guides then state that language for code, and F6 does not apply to it. English stays the rule of every repository that records nothing. How the decision is recorded is the TechSpec's.
  2. The Baseline has no exception. The rule stays English for all, and the legacy Quaestor either holds its Baseline at the version before the rule or carries the clause beside the ADR that contradicts it.

  Default until he answers: the approved shape, which has no exception. F8's message tells the legacy Quaestor that the update that carries the rule contradicts its ADR-0141 and that it holds that update, and the PRD does not fix the reach of F1 before the answer.
- Which Reader Language do Pantheon and the Secondbrain declare? The maintainer's second answer names the Gesttione products, Roundfix and the skills catalog; neither of the two is a product with users of its own. Default: each asks the maintainer when its update asks for the answer, and F8's message to them names no language.
- Does the Reader Language carry a suggested value? ADR-0118 gives a decision a suggested default that one flag adopts. Default: none, so no repository gets a language it did not choose. The TechSpec decides.
- What do the skills catalog and Praetor do? Neither holds a Setup Manifest, so no Baseline update asks them. Default: nothing in V1. The skills catalog's own rules already say English for every tracked file, and Praetor's pointer to the user-level file is Praetor's to change.
- Does `docs/references/0001-ubiquitous-language-in-colocated-context-md.md` still bind? The catalog does not ship it and no guide cites it; its text agrees with this Idea. Default: prior art, not a rule.
- Which names does the glossary take? Default: the PRD's Glossary Declaration settles them, starting from Reader Language, Untranslated Term and the name-in-code line of a Glossary Definition.
- Which repositories receive the message of F8 because of their own rules? Default: Vortex, which repeats the rule; the new Quaestor, whose ADR-0010 left the glossary folder to the Baseline's owner; the legacy Quaestor, by the first question; and Oraculum, whose code is Portuguese. The PRD repeats the sweep.
- On what outside evidence does the acceptance rest? Default: the published passages quoted under `### Evidence`, which the QA gate can check without the network.
