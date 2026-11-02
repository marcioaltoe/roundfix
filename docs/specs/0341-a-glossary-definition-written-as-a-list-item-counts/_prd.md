---
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
status: active
created: 2026-11-02
surfaces: [backend, docs]
type: fix
---

# A glossary definition written as a list item counts

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by pantheon on
2026-11-02: `roundfix spec check` reads a glossary term as defined only when
its line starts with `**Term**:`, so Pantheon's glossary, which lists its
vocabulary as `- **Term** — definition`, never defined the term its Spec 0055
added; the Task that wrote the term passed its own Verification, the QA gate
refused its precondition with `SC-GLOSSARY-MISSING`, and nothing said which
form the check wanted. After this Spec a Glossary Definition is a line that
starts with the bolded term followed by a colon or by a spaced em dash, bare
or as a top-level `-`, `*` or `+` list item; the finding names those forms,
the guides state them, and the glossary Verification `write-tasks` asks for
is a line-anchored command that passes only on a form the check accepts.

## Sources

One adopter report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md`](references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md),
   `type: fix`, Pantheon, ai-memory message `01a19a04`.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry
on the glossary definition format.

## Prerequisites

None in `requires`. Authored against `900e13f8` (v0.129.0), where no other
Spec is active.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads lines of an existing Markdown file. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads local files. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0244: "A Spec now carries a
  Glossary Declaration" and the check holds it to the glossary
  "deterministically"; this Spec widens only which glossary lines define a
  term and keeps the check deterministic. ADR-0280: "The glossary is held to
  its own rule": this repository's `CONTEXT.md` keeps its own Repository
  Contract Test, which reads `**Term**:` lines and stays unchanged. No new
  ADR: the accepted forms are a glossary and guide fact. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AV answers of 2026-11-02 quoted in the authorization record, which keep the
  cycle AU rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the canonical Roundfix `spec` reference,
  `write-prd` glossary guide and `write-tasks` skill and the `SKILL.md` pairs
  of the three skills. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0341-a-glossary-definition-written-as-a-list-item-counts/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/glossary.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A term written as `- **Term**:`, `- **Term** —`, `**Term** —` (or with a
  `*` or `+` list marker) counts as defined, for every glossary code: the
  bold candidate coverage of `SC-GLOSSARY-UNDECLARED`, the changed-term check
  of `SC-GLOSSARY-UNPLANNED` and `SC-GLOSSARY-MISSING`.
- Every form accepted today stays accepted, and bold text that is not a
  definition (a heading, an indented item, bold inside a sentence, a bold
  label followed by prose) still defines nothing.
- `SC-GLOSSARY-MISSING` names the accepted forms, and `SC-GLOSSARY-UNPLANNED`
  prints the glossary Verification command that binds the term.
- The `write-prd` glossary guide, the `write-tasks` Glossary requirement and
  the Roundfix `spec` reference state the forms, and `write-tasks` asks for
  the same command; the command never passes on a line the check rejects.

## Core Features

1. The Glossary Definition reader. One pattern in `internal/speccheck`
   accepts the colon and the spaced em dash after a bolded term at the start
   of a line, bare or after one top-level `-`, `*` or `+` list marker and a
   space.
2. The messages. The `SC-GLOSSARY-MISSING` fix names the forms; the
   `SC-GLOSSARY-UNPLANNED` fix for a term without a binding Task prints the
   line-anchored Verification command, rendered by the same package.
3. The mirror proof. A test runs the rendered command against every form
   the reader is tested on and requires both to agree.
4. The rules in writing. The spec command guide, the three skills and
   `CONTEXT.md`, which gains the term **Glossary Definition**.

## Non-Goals / Out of Scope

- A heading (`### Term`) as a definition, decided by the maintainer.
- An indented (nested) list item, an ordered list item, an ASCII hyphen or
  double hyphen, an en dash, an em dash without spaces, and a parenthetical
  between the term and its separator (`- **Term** (_alias_) —`); see
  `_techspec.md` → Decisions for each reason.
- The Glossary Declaration grammar, the bold candidate rule, the glossary
  horizon, the glossary file set and the binding rule (a Verification naming
  the exact bolded term).
- This repository's own `CONTEXT.md` format contract, which keeps the
  canonical `**Term**:` lines, and the upstream `domain-modeling` skill.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`.

## Success Metrics

1. Success Metric: in a temporary repository whose Spec adds a term bound by
   a completed Task, `roundfix spec check` reports no `SC-GLOSSARY-MISSING`
   when the glossary holds the term as `- **Term**:`, `- **Term** —`,
   `**Term** —`, `* **Term**:` or `+ **Term**:`; at `900e13f8` each of those
   reports one.
2. Success Metric: `**Term**:` is still accepted, and `### Term`, an indented
   `  - **Term**:`, `1. **Term**:`, `- **Term** - x`, `- **Term** -- x`,
   `- **Term** – x`, `- **Term**—x`, `-**Term**: x`, `- **Term** is prose`,
   `- **Term.** Prose` and `See **Term**: x` still report
   `SC-GLOSSARY-MISSING`; a bold candidate in the PRD is covered exactly by
   the accepted forms.
3. Success Metric: the `SC-GLOSSARY-MISSING` fix names the four forms and the
   `*` and `+` markers; the `SC-GLOSSARY-UNPLANNED` fix for an unbound term
   prints `grep -F -- '**<term>**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*<term>\*\*(:| —( |$))'`
   with the term and its regular-expression escaping.
4. Success Metric: the rendered command exits 0 on exactly the lines the
   reader accepts, over every row of Success Metrics 1 and 2 plus a term with
   a regular-expression metacharacter and a term with an apostrophe.
5. Success Metric: both copies of the `write-tasks` skill carry the command
   rendered for `<Term>` and both copies of the `write-prd` glossary guide and
   the Roundfix `spec` reference name the accepted forms; the spec command
   guide and `CONTEXT.md` state them.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Pantheon report quoted in the adopted reference, written by the
  Pantheon project before this Spec existed, with its four measured forms.
- The CommonMark specification, section 5.2 List items
  (https://spec.commonmark.org/0.31.2/, read through Exa on 2026-11-02): "A
  bullet list marker is a `-`, `+`, or `*` character." The reader accepts
  exactly those three markers.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-02 against these sources.

- Characterization with the installed `roundfix` 0.129.0 (`900e13f8`) in a
  temporary repository whose Spec adds a two-word term bound by a completed
  Task: a glossary line `**Term**: def` reports nothing; `- **Term**: def`,
  `- **Term** — def`, `**Term** — def`, `* **Term**: def`, an indented
  `  - **Term**: def`, `### Term`, `- **Term** - def`,
  `- **Term** is prose` and `See **Term**: def` each report one
  `SC-GLOSSARY-MISSING`. Pantheon's four measured cases reproduce.
- Corpus characterization, read only, of the glossary forms in use (counts
  of lines by form):
  - this repository's `CONTEXT.md`: 370 `**Term**:` lines, no list-item, em
    dash or ASCII hyphen definition and no indented bold list item, so the
    reader change adds no term here;
  - the local `CONTEXT.md` of conexus, fluxus, lexia, onioncry, oraculum and
    vortex: only `**Term**:` definitions; their 70 `- **Term** <prose>` list
    lines (relationships such as `- **Workspace** belongs to ...`) carry
    no colon or em dash right after the bold and stay non-definitions;
  - the Secondbrain mirror of Pantheon's `CONTEXT.md`
    (`projects/pantheon/mirror/CONTEXT.md`, 2026-11-02): one `**Acervo**:`
    line (the workaround), 16 `- **Term** —` lines that become definitions,
    and 9 `- **Term** (_alias_) —` lines that stay out of scope.
- Reading at `900e13f8`: `glossaryDefinitionPattern` in
  `internal/speccheck/glossary.go` is the only definition reader of the
  check; `glossaryFiles` builds the defined term set that the bold coverage,
  the changed-term check and the missing-term check all read; the
  `SC-GLOSSARY-MISSING` fix reads "Define the declared term in the glossary.";
  `write-tasks` says the Verification "checks the bolded term followed by a
  colon in the glossary on whitespace-normalized text", which also passes on a
  mid-sentence `**Term**:` that the check rejects.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "glossary definition format list item bold term em dash CONTEXT.md" --all --files --min-score 0.3`.
  It returned mirrors of the spec command guide, the upstream `write-prd`
  glossary guide and the fleet's domain guides; none decides a definition
  form.
- Exa: the CommonMark passage quoted above.

## Glossary

- adds: **Glossary Definition** — the glossary line shapes that define a term

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AV brief for Spec 0341 of 2026-11-02 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-02), verbatim:

- Scope: "Aceitar lista + mensagem (Recommended)" — option text: "Reconhecer
  também `- **Termo**:` e `- **Termo** —` (e `**Termo** —`) como definição; a
  mensagem do achado e os guias write-prd/write-tasks dizem os formatos
  aceitos; a Verification gerada pelo write-tasks passa a casar com qualquer
  um deles."
- Release: "v0.129.1 logo após o merge (Recommended)" — option text: "Patch
  único ao fim do ciclo; desbloqueia o Pantheon. Aviso aos 8 adotantes, com
  resposta ao Pantheon."
- Rules: "Mesmas do AU (Recommended)" — option text: "Sem supervisão: push,
  PR do plano, entrega e merge pela fila; Makefile/CI/go.mod/lint liberados
  se o desenho pedir; falha técnica = consertar e seguir; findings do review
  corrigidos ou dispensados com evidência; mensagem do Pantheon consumida
  quando a Spec estiver no main."

Taken from those answers and the characterization, with no further
interview:

- The three CommonMark bullet markers are accepted alike; leading whitespace,
  ordered markers and ASCII dashes are not, because no glossary in the corpus
  uses them and each would turn sub-points or prose into terms.
- The glossary Verification becomes a line-anchored command that is exact
  in case and never more permissive than the reader, rendered in one place
  and printed by the check.
- No new ADR; the forms live in `CONTEXT.md` and the guides.
- Two implementation Tasks in series and a QA gate.

## Open Questions

None.
