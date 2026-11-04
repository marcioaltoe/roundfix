---
spec: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
status: active
created: 2026-11-04
surfaces: [docs]
type: feat
---

# The Backlog and Findings contracts show the Waiting Declaration

Since Spec 0343 an open Backlog Entry or a pending or partial Finding may carry a Waiting Declaration, and `roundfix follow-up` holds such an item back as waiting. The Baseline's copyable Backlog Operational Contract and Findings Operational Contract, which every adopter's agent copies when it writes an entry or a Finding, do not name the two keys, so an agent that writes outside the Roundfix skill's follow-up mode never learns that a wait can be declared. After this Spec both copyable frontmatter blocks carry `waiting_on` and `since` as commented lines, each contract adds one sentence saying when to uncomment them, and every adopter receives the change through `roundfix baseline update`.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md`](references/2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md),
   `type: feat`, written on 2026-11-03 when ADR-0353 deferred the Baseline
   change.

At the authoring base `docs/findings/` holds no unresolved Finding and
`docs/backlog/` holds no other open entry that shares this context.

## Prerequisites

None in `requires`. Authored against `6a0ed904` (v0.130.1). Cycle AX
delivers 0345 → 0348 → 0347 → 0346, so this Spec rebases last. It owns the
`context-workflow` Baseline module and the docs layout guide. It may share
with the other three only the derived Baseline records (module version
record, catalog digests, plan goldens, the Standard TypeScript profile and
the Setup Manifest), the Behavior Surface Record, the Skill Coverage Map and
`CONTEXT.md`, each of which it changes in one narrow place.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change adds two commented lines and one sentence to two
  existing Baseline clauses, which keep their clause identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the Baseline renders local guides. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0353: "Change the Baseline
  contracts in the same change: deferred to the open Backlog Entry", which
  this Spec adopts; its Decision stays as written ("The command that reads a
  declaration validates it."), so no validation moves into the Baseline or
  the Spec Check. No new ADR: the change shows an accepted contract in the
  guide that adopters copy. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record, which
  allow a Baseline change that reaches adopters and keep the cycle AJ grant
  for the Governed Paths each Spec declares; no Makefile, workflow, Go
  module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0346-the-backlog-and-findings-contracts-show-the-waiting-declaration/_authorization.md`;
  bounded files: `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`.

## Goals

- An agent that copies either Operational Contract sees the two optional
  keys of a Waiting Declaration in the frontmatter it copies, and a copy
  made unchanged declares no wait.
- Each contract says in one sentence which items may uncomment the keys and
  what each key holds.
- Every other Baseline clause, and every other sentence of the two changed
  clauses, keeps its bytes; the two clauses keep their identity and force.
- An adopter's `roundfix baseline update` brings the change as a Managed
  Refresh of its docs layout guide and Setup Manifest, and a Source Baseline
  transition records both clauses as retained.

## User Stories

1. As an agent writing a Backlog Entry in an adopter repository, I want the
   copyable frontmatter to show `waiting_on` and `since`, so that I can
   declare that the entry waits on someone outside the repository without
   reading the Roundfix skill.
2. As an agent writing a Finding, I want the same two keys in the Findings
   contract, so that a Finding that waits is held back by the follow-up list
   like a waiting Backlog Entry.
3. As an adopter maintainer, I want `roundfix baseline update` to show the
   change as an update of my docs layout guide, so that I can review it like
   any Baseline refresh and keep my repository's other guidance.

## Core Features

1. The Backlog contract. Its copyable frontmatter ends with two commented
   lines, `waiting_on` and `since`, after `reason`, and one sentence after
   the frontmatter block says that an `open` entry waiting on someone or
   something outside the repository may uncomment them, naming who or what
   in one line and the date the wait began.
2. The Findings contract. Its copyable frontmatter ends with the same two
   commented lines after `updated_at`, and one sentence after the template
   says the same for a `pending` or `partial` Finding.
3. Non-binding guidance. The keys stay commented, so a copied frontmatter
   declares nothing; no validation rule, Spec Check finding or command
   behavior changes.
4. Delivery to adopters. The `context-workflow` module records its next
   version; the derived Baseline files, the Standard TypeScript golden and
   this repository's own guide are regenerated by their commands; the user
   guide's published Findings template matches the generated contract.
5. The glossary. The **Waiting Declaration** entry names the two contracts
   that show it.

## User Experience

An adopter who runs `roundfix baseline update` after the release sees a
plan whose file changes are `update docs/agents/docs-layout.md` and
`update docs/agents/setup-context.json`, with no retention evidence on a
Managed Refresh; the refreshed guide shows the two commented lines in each
contract and one sentence after each. Bytes outside the Baseline markers
stay unchanged.

## Non-Goals / Out of Scope

- Any change to how `roundfix follow-up` reads, validates or reports a
  Waiting Declaration (ADR-0353 stays the contract).
- A Spec Check finding or `baseline update` diagnostic for a malformed
  declaration.
- The Waiting Declaration on messages, notes, Specs, ADRs or handoffs.
- Uncommented placeholder keys or `null` values in the templates.
- The Source Baseline corpus, manifest, accounting and index, the retention
  transitions, any other module, profile or guide, any skill, the Makefile,
  the CI workflows, `go.mod`, `go.sum` and the lint or formatter
  configuration.

## Success Metrics

1. Success Metric: the embedded Backlog and Findings contracts each end
   their copyable frontmatter with the two commented lines and carry their
   one sentence exactly once; at `6a0ed904` neither does.
2. Success Metric: the copyable frontmatter of each contract parses as YAML
   with neither `waiting_on` nor `since`, and the same block with the two
   lines uncommented parses with both.
3. Success Metric: the Standard TypeScript golden and this repository's
   docs layout guide each carry each commented line twice and each sentence
   once, and the user guide's Findings template equals the generated one.
4. Success Metric: both clauses keep their identity and `mandatory` force,
   have no `replaces`, and a transition from Source Baseline
   `baseline.standard-typescript-monorepo-0.0.1` records each as `retained`;
   every other clause of the module keeps its bytes against the Delivery
   Base.
5. Success Metric: `roundfix baseline update` previewed on a copy of this
   repository at its pre-change state lists exactly the two file changes of
   the User Experience, and a second Managed Refresh after applying changes
   nothing.

## Acceptance evidence

The outside-evidence row rests on a source this Spec did not write: the
YAML 1.2.2 specification, section 3.2.3.3 Comments
(https://yaml.org/spec/1.2.2/, read through Exa on 2026-11-04): "Comments
are a presentation detail and must not have any effect on the serialization
tree or representation graph." A commented key is therefore absent for
every YAML reader, which is why a copied template declares no wait.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files, its
embedded catalog and temporary copies.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Reading at `6a0ed904`: the two clauses are
  `clause.context.backlog-01-operational-contract` and
  `clause.context.findings-01-frontmatter` of the docs layout rule; the
  follow-up inventory reads frontmatter with a YAML decoder and treats an
  absent key or YAML null as absent; the user guide's Findings template is
  held equal to the generated contract by a repository contract test.
- Prototype in a scratch clone: the edited module, `-record-module-versions`,
  `make baseline-digests` and two Managed Refreshes changed exactly the
  files the authorization lists plus the module version record, the three
  catalog files and the four plan goldens; the second refresh reported
  `File changes: 0`; the Baseline package (1,510 tests), the docscontract
  package and the CLI's baseline and follow-up tests passed. The preview of
  `roundfix baseline update` on a copy at `6a0ed904` reported
  `File changes: 2`, `update docs/agents/docs-layout.md (2 managed entries)`,
  `update docs/agents/setup-context.json (1 managed entries)` and
  `Retention evidence: 0`.
- Pending ai-memory messages for this project (2026-11-04): two, from
  Pantheon and Fiscus, which Specs 0345 and 0348 adopt; none concerns these
  contracts, and none was consumed.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "waiting_on since backlog entry finding frontmatter wait declaration" --all --files --min-score 0.3`;
  it returned only mirrors of this repository's own documents (the adopted
  Backlog Entry, ADR-0353, the follow-up command guide and skill reference,
  the glossary), so no outside page decides a frontmatter format for
  waiting items.
- Exa: the YAML passage quoted above.

## Glossary

- changes: **Waiting Declaration** — the glossary entry names the Baseline contracts that show it commented

## Skills

- unchanged: guide docs/user-guide/context-driven-development.md — the guide's Findings template gains the two commented keys and one sentence; its covering skills describe how a Baseline is adopted and how Specs are written, implemented, checked and archived, not the frontmatter keys of a Finding or a Backlog Entry, so no skill text changes.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AX brief for Spec 0346 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04), verbatim:

- 0346: "Sim, Spec própria (Recommended)" — option text: "Spec pequena: as
  duas chaves comentadas nos dois blocos de frontmatter + uma frase; versão
  do módulo e arquivos derivados seguem as regras da Baseline; chega aos
  adotantes no baseline update."
- Rules: "Mesmas do AW (Recommended)" — option text: "Sem supervisão: push e
  PR dos planos, entrega e merge pela fila, consertar e seguir em falha
  técnica; Makefile, .github/workflows, go.mod/go.sum e lint liberados
  quando o desenho pedir; mudança de Baseline que chega aos adotantes
  liberada."
- Release: "Uma v0.131.0 no fim (Recommended)" — option text: "Uma minor (a
  mudança da Baseline e o runs list fora de repo são feat); notas aos 8
  adotantes, resposta ao Pantheon com o defeito corrigido e o contorno que
  podem largar."

Taken from those answers and the characterization, with no further
interview:

- The keys are commented lines, not placeholders or `null` values: an
  uncommented placeholder such as `since: YYYY-MM-DD` would be a malformed
  declaration in every copy.
- The sentence says "may": the declaration stays optional guidance inside
  the two mandatory clauses, whose identity and force do not change.
- The published Findings template in the user guide changes with the
  generated contract, because a repository contract test holds them equal.
- One implementation Task and the QA gate: the module, its derived files,
  the guides and the glossary must change together for the repository gate
  to stay green.
- No new ADR.
- The release is the operator's single v0.131.0 after the four cycle AX
  Specs; no Task releases.

## Open Questions

None.
