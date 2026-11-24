---
type: feat
status: open
created: 2026-11-24
spec: null
reason: null
---

# The Baseline carries the language rule for code, data and products

## Intent

Carry the maintainer's language rule of 2026-11-20 into the Baseline, so every adopting repository receives it through `roundfix baseline update`. This is the future Spec 0417. The rule, as the Jarvis project sent it:

- English for code and infrastructure: variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs.
- The reader's language for what a person reads (screens, user messages, e-mails, PDFs, documentation for people), declared per repository: the Gesttione products declare pt-BR; Roundfix and skills declare English. In the Gesttione products, MCP tool names, descriptions and parameters stay in Portuguese.
- Exception for Brazilian tax and legal terms inside English code: official acronyms stay as they are (`icms_amount`, `efd_file`); institutes with no faithful translation stay in Portuguese without accents (`escrituracao`, `apuracao`), the rest of the name in English (`apuracao_summary`). Each such term enters the repository glossary (`CONTEXT.md`) with the name used in the code.
- Application: a new name is born under the rule; old code changes when the part is touched, with its own Spec when the change is large; no batch renaming without a request.

## Maintainer decision

The maintainer, on 2026-11-20 (relayed by the Jarvis project, message `01a1f991-12bb-7621-a985-f3456059763c`, 14:13 BRT):

> "Talvez, até mesmo essas regras aí de idioma devam ir para a roundfix adicionar na baseline. Ficaram muito boas desse jeito aí."

Scope as read by the Jarvis project: where the rule enters the Baseline (an agent guide, a language module, or whatever fits best); the glossary as a Baseline item in repositories with a tax domain; a review/QA check if it makes sense. PRD, TechSpec, Tasks, implementation and PR were authorized; protected Baseline paths, merge and release go through a `PEDIDO DE DECISÃO` with the exact list, as in Spec 0415. Priority after Specs 0415 and 0416.

## Idea answers

The Idea stage was answered by the maintainer on 2026-11-20, 15:28 BRT (message `01a1f9b0-2f68-7503-817e-929d363c56e8`, answering request `01a1f9a2-fa98-7992-a272-d5c19e3f9701`):

1. Proceed with the Idea.
2. "Cada repositório declara": code in English for all; the Gesttione products declare pt-BR for the reader; Roundfix and skills declare English; each repository answers once at its next Baseline update.
3. "No CONTEXT.md, com o nome no código": the term stays in the existing glossary (`CONTEXT.md`) and records the name used in code, translated (Filial -> `branch`) or kept (Apuração -> `apuracao`). No `docs/glossary` folder.
4. "Review e QA": review flags a new name outside the rule in the diff and QA flags user-facing text in the wrong language, as any finding, with no new tool.
5. "Não": the rule covers what the product user reads (screens, messages, e-mails, PDFs, manual, MCP tools); Specs, ADRs, agent guides, commits and PRs stay as each repository already does.

An Idea draft exists outside the repository at `/Users/marcio/.claude/jobs/e807deb6/tmp/idea-0417-draft.md`; nothing was copied from it.

## Source

Cross-project messages from the Jarvis project, 2026-11-20: `01a1f990-ed12-7882-81a1-bf0c370fda7f` (the rule), `01a1f991-12bb-7621-a985-f3456059763c` (the decision), `01a1f9b0-2f68-7503-817e-929d363c56e8` (the Idea answers). Nature: reported. Recorded by triage on 2026-11-24.
