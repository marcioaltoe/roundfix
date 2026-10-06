---
spec: 0239-a-glossary-every-spec-keeps-current
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# A glossary every Spec keeps current

Roundfix was built on the CONTEXT-driven method: the glossary (`CONTEXT.md`)
and the ADRs are the base, and every Spec, Task and line of copy draws its
vocabulary from them. On 2026-10-06 the maintainer observed that the loop no
longer keeps the glossary current: it last changed on 2026-10-02, twenty-one
Specs archived after it without touching it, and authoring and QA reports
named new terms that nobody wrote down
([the adopted Backlog Entry](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)).
The domain guide already asks for a glossary check at the close of a Spec, but
no gate reads the answer, the autonomous authoring route never reaches
`domain-modeling`, and the operator's briefing had told authors not to edit the
glossary.

This Spec makes the glossary a checked part of every Spec. A Spec declares the
domain terms it introduces in a Glossary Declaration, one of its own Tasks
writes them through `domain-modeling`, and the Spec Consistency Check, the QA
gate and the Archive Command refuse a Spec that leaves a Glossary Gap. The
Baseline tells adopters the same, the authoring skills follow it on the
autonomous route, and one Task adds the terms dropped since 2026-10-02.

## Prerequisites

None that the Verification depends on. Specs 0240, 0241 and 0242 are authored
in the same cycle and this Spec is delivered first. Spec 0241 may edit the
same Baseline module and Spec 0242 the Archive Command, and several of them
raise owned-skill versions; the operator orders the queue, every version in
this Spec rises from the value on its Task's starting commit, and ADR-0233
regenerates a raised skill version at merge.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the three
  new finding codes follow the existing `SC-<FAMILY>-<STATE>` form and every
  glossary term keeps the glossary's bold-heading form. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the check, the refusal and the
  catch-up read local files and Git history only; no request, credential or
  provider call is added, the Jev judge is not called, and every test uses a
  temporary repository. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0244 (this Spec) decides the
  Glossary Declaration, the bold rule, the binding Task, the three findings and
  the horizon. ADR-0094: "a Spec that passes today's gates must not be blocked
  by a new check that merely disagrees with its style", which the horizon
  keeps. ADR-0117 checks each defect at the stage that can produce it, so the
  declaration is read at the PRD stage and the binding at the Tasks stage.
  ADR-0168 opens a related-ADR gap only for ADRs that predate the Spec, the
  same commit-ancestry horizon the glossary rule uses. ADR-0183: an attributed
  claim carries a receipt the check proves, unchanged. ADR-0200: "Spec
  authoring gets an advisory judge that never gates", so the judge is not the
  gate. ADR-0186: "It never cites a Spec number or an ADR number", which the
  reworded Baseline clauses obey while they name Roundfix commands. ADR-0187
  splits the Roundfix Skill by command, ADR-0189 ties an owned skill's version
  to its content, and ADR-0233 regenerates a raised skill version at merge.
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts", answered in
  the TechSpec with the reason none applies. ADR-0240: "One QA partial policy,
  and rows a Run sandbox cannot reach" decides when its QA partial qualifies.
  The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155,
  ADR-0156 and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and
  ADR-0183 check this Spec's consistency by citation and receipt. ADR-0178
  authorizes each Task commit by its grant, ADR-0182 runs Settlement Checks
  before it, and ADR-0166 records undeclared paths. ADR-0130: "every path any
  authorization has ever bounded must be matched by the declared set", so the
  paths this Spec bounds stay governed. The catch-up terms are taken from ADR-0174, ADR-0196,
  ADR-0231, ADR-0232, ADR-0236, ADR-0238, ADR-0239 and ADR-0240, which this
  Spec cites as sources and does not change. ADR-0201 fixed the judge's ceiling that ADR-0231 made configurable, ADR-0237
  reads the Delivery Commit when a Delivery Retry records a merge made outside
  the queue, and ADR-0235 retired the Jev Router, which is why it is not
  added; these are sources too and do not change. ADR-0096 and ADR-0097 cite
  ADR-0080 but decide the gate's machine stage and row carry; ADR-0179 cites
  ADR-0130 but decides an empty paths list; ADR-0192 cites ADR-0178 but
  decides how a conflict confined to declared derived paths is resolved;
  ADR-0209 cites ADR-0200 but decides the grouping suggestion; ADR-0229 cites
  ADR-0167 but decides how an operator archive resumes a park; ADR-0194,
  ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when
  it is observed again and its evidence snapshot; ADR-0208, which ADR-0209
  cites, decides that sources sharing a context share a Spec; this Spec
  changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Spec Consistency Check sources, the
  corpus golden and its characterization, the CONTEXT workflow module and its
  derived files, and the owned skills are Governed Paths. The maintainer
  authorized skill edits ("considere autorizado a ajustar todas as skills se
  necessário"), the Baseline source and guides ("Autorizar os dois") and, on
  2026-10-06, the Governed Paths this Spec declares ("Concedo"). No other
  Governed Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0239-a-glossary-every-spec-keeps-current/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/glossary.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`.

## Goals

- A Spec that introduces a domain term cannot pass its QA gate or archive
  until the glossary defines that term, and the answer is the same on every
  run.
- The autonomous authoring route reaches `domain-modeling` whenever it names a
  concept, and its docs Task writes the term.
- Adopters receive the same rule through the Baseline: the glossary and the
  ADRs are the base, and a Spec writes the terms it introduces.
- The glossary carries the terms the cycle since 2026-10-02 introduced and
  dropped.

## User Stories

1. As the maintainer, I want every Spec to leave the glossary defining the
   terms it introduced, so that the glossary and the ADRs stay the base of the
   process instead of a snapshot from four days ago.
2. As a Spec author on the autonomous route, I want the check to name each
   bolded term the glossary lacks before a Run starts, so that I declare it,
   plan its Task, or declare it not a domain term while authoring is cheap.
3. As the operator, I want the QA gate and `roundfix archive` to refuse a Spec
   whose declared term never reached the glossary, so that a dropped term is
   caught before the Pull Request rather than in a later audit.
4. As an adopter, I want my guides to tell my agents the same rule, so that my
   glossary does not go stale the way this repository's did.

## Core Features

1. **The Glossary Declaration.** A PRD or TechSpec carries a `## Glossary`
   section listing each domain term the Spec adds or changes, and each bolded
   phrase it declares not a domain term with a reason, or `None.`. A PRD
   committed at or after the commit that added the glossary guide to the
   `write-prd` skill must carry one (ADR-0244).
2. **Undeclared terms are reported.** The Spec Consistency Check reports a
   bolded phrase of two to five capitalized words in the PRD or TechSpec that
   the glossary does not define and the declaration does not cover, and a
   missing or malformed declaration where one is required.
3. **Declared terms are planned and then present.** A term the Spec adds or
   changes needs a non-QA Task that declares the glossary file and names the
   bolded term in its Verification. Once every such Task has completed, the
   glossary must define the term. The QA gate's precondition and the Archive
   Command refuse a Spec that leaves any of these gaps, the latter also under
   a QA Archive Override.
4. **The glossary is found where the repository keeps it.** The check reads the
   root `CONTEXT.md`, or `GLOSSARY.md` under the upstream skills' new name,
   and each context file their map file links.
5. **The Baseline states the base.** The CONTEXT workflow's domain clauses say
   that the domain context and the accepted ADRs are the base, and that a Spec
   adds the terms it introduces through `domain-modeling` in one of its own
   Tasks, into the selected domain context whatever file name the skill uses.
6. **The skills follow the route.** `write-prd` and `write-techspec` activate
   `domain-modeling` whenever they name a concept, also on the autonomous
   route, and write the declaration; `write-tasks` gives each declared term a
   glossary requirement in the docs Task; the QA gate records each declared
   term; `grilling` and `grill-with-docs` stay the interactive entry; the
   Roundfix Skill describes the findings and the refusal.
7. **The catch-up.** The glossary gains the terms the cycle introduced and
   dropped, each checked against the ADR that decided it, and this Spec's own
   terms.

## Non-Goals / Out of Scope

- A model deciding whether a phrase is a domain term, or any new Jev judge
  question; the judge stays advisory and is not called (ADR-0200).
- Detecting lowercase or one-word terms in prose; those reach the glossary
  through the declaration and the authoring skills.
- Retiring glossary terms, or checking that a term a Spec retires is removed.
- Renaming `CONTEXT.md` to `GLOSSARY.md` in this repository, or editing the
  upstream `domain-modeling`, `grilling` or `grill-with-docs` skills.
- Reading ADR text for terms; an ADR's new term reaches the check through the
  Spec that writes the ADR, which declares it.
- Changing any archived Spec, or requiring the section of a Spec authored
  before the horizon.
- Writing to, or sending the content of, any adopter repository.

## Success Metrics

1. Success Metric: in a temporary repository, a PRD that bolds a two-to-five
   word capitalized phrase the glossary lacks and does not declare reports
   `SC-GLOSSARY-UNDECLARED`; declaring it not a term or defining it in the
   glossary clears it; a one-word bold label, a phrase with a digit and a
   phrase inside a fenced block are never reported.
2. Success Metric: a term declared as added with no non-QA Task that declares
   the glossary file and names the bolded term in its Verification reports
   `SC-GLOSSARY-UNPLANNED`; with that Task completed and the term absent,
   `SC-GLOSSARY-MISSING`; the QA gate precondition is blocking, and
   `roundfix archive` exits `2` naming the gap, with and without
   `--qa-override`.
3. Success Metric: a PRD committed before the glossary guide's adding commit,
   without the section, reports no glossary finding and records the skip; one
   committed after reports the missing section; a PRD that carries the section
   is checked at any age; the staged and full checks agree, and the active
   corpus reports the three codes at zero.
4. Success Metric: `docs/agents/domain.md` states both reworded clauses as
   `mandatory`, a Source Baseline adopter's Managed Refresh plan records both
   `retained`, and a second refresh changes no file.
5. Success Metric: the glossary defines each of the ten catch-up terms and this
   Spec's three declared terms, and each catch-up definition agrees with the
   ADR the TechSpec names for it.
6. Success Metric: the authoring skills name `domain-modeling` and the
   declaration on the autonomous route, `write-tasks` gives each declared term
   a glossary requirement, the QA gate records the declared terms, and every
   raised owned-skill version is recorded.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- **Adopter measurement (2026-10-06, read-only, local checkouts, no fetch).**
  Last glossary change, defined terms, and Specs archived after that change:
  conexus `CONTEXT.md` 2026-08-18, 108 terms, 8 of 33 archived Specs after it;
  fluxus 2026-08-25, 82 terms, 0 of 56; vortex root `CONTEXT.md` 2026-08-25
  with 6 terms and `packages/backend/CONTEXT.md` 2026-08-14 with 3, under a
  `CONTEXT-MAP.md`, 1 of 46; tax-poc 2026-07-14, 69 terms, 1 of 10;
  oraculum 2026-08-26, 38 terms, 0 of 48. Every local checkout stops at or
  before 2026-08-26, so it cannot show the drift after 2026-10-02; conexus
  shows the same pattern, eight Specs without a glossary change. This
  repository: 225 terms, last change 2026-10-02, 21 of 220 archived Specs after
  it. Glossary commits per archived Spec: conexus 7/33, fluxus 13/56, vortex
  4/46, tax-poc 4/10, oraculum 12/48, Roundfix 71/220.
- **Run Database measurement (operator, 2026-10-06, recorded in the adopted
  Backlog Entry).** `domain-modeling` was read in 25 of 220 Runs; `grilling`
  and `grill-with-docs` in none.
- **Replay of the bold rule (2026-10-06).** Over the PRDs and TechSpecs of
  archived Specs 0200 to 0237, phrases of two to five capitalized words without
  a digit that the glossary does not define: eight in five Specs, six distinct
  concept names, none defined. Counting one-word bold phrases added twenty-four
  more, nearly all labels.
- **Published literature.** Eric Evans, *Domain-Driven Design* (2003), on the
  Ubiquitous Language: "Recognize that a change in the language is a
  change to the model" (<http://ddd.fed.wiki.org/ubiquitous-language.html>).
  James Shore, *The Art of Agile Development*: "The model and the ubiquitous
  language must always stay in sync"
  (<https://www.jamesshore.com/v2/books/aoad1/ubiquitous_language>).
- **Upstream skills changelog.** Matt Pocock, "v1.3: /implement-spec, /pr,
  /retro, and GLOSSARY.md" (aihero.dev, 2026-10-05): the domain-doc file is
  renamed `CONTEXT.md` to `GLOSSARY.md` and `CONTEXT-MAP.md` to
  `GLOSSARY-MAP.md`, and "The skills look for the new names only". The
  vendored `domain-modeling` skill already uses the new name, so a repository
  rule must point it at the selected domain context.

## Glossary

- adds: **Glossary Declaration**
- adds: **Glossary Gap**
- adds: **Light Tier**
- adds: **Light Spend Log**
- adds: **Jev Ceiling**
- adds: **Stage Key**
- adds: **Tested Base**
- adds: **Merge Evidence**
- adds: **Delivery Commit**
- adds: **Item Branch**
- adds: **Network-Denied Row**
- adds: **Pre-PR Review Command**
- changes: **Spec Consistency Check**

## Decisions

- A Spec declares its terms in a `## Glossary` section; bold phrases of two to
  five capitalized words are the only terms detected without a declaration;
  see ADR-0244.
- A declared term binds to a non-QA Task by its glossary declaration and its
  Verification, not by a Task id written in the PRD; see ADR-0244.
- The section is required from the commit that adds the `write-prd` glossary
  guide; see ADR-0244.
- The Jev judge is not part of the gate and gains no question in this Spec.
- The catch-up edits `CONTEXT.md` directly through `domain-modeling` (the
  maintainer's "Sim").
- The open Backlog Entry on retiring unused skills, which the advisory judge
  suggested grouping with this Spec (0.44), stays with Spec 0241: the
  maintainer ordered the two separately, and together they exceed four
  implementation Tasks.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
