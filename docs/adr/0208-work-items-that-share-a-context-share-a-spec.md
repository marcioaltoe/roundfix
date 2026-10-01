---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Work items that share a context share a Spec, and an open item is extended

Nothing in the Baseline told an author to look for related work before
minting a Spec or a new record. One Spec per Finding, and a new file for an
observation an open Backlog Entry already holds, both looked like the safe
default. On 2026-10-01 the maintainer asked for the opposite: "podemos ter em
uma mesma spec um ou mais inbox/backlogs ou findings não é obrigatório e nem
desejavel uma spec para um finding com contexto similar ou complementar a
outro. O cuidado é não ter specs desnecessáriamente complexas e grandes, o uso
do jev pode ajudar a decidir pela reunião de inbox, backlogs e findings em uma
só spec. Isso deve fazer parte do baseline de roundfix e parte das regras do
roundfix." The same day the maintainer added: "podemos alterar e adicionar e
atualizar um finding ou backlog já existente e ainda não implementado para que
não seja necessário criar outro arquivo".

The Baseline therefore carries three mandatory clauses:

- **Grouping.** One Spec may adopt several Inbox Entries, Backlog Entries and
  Findings whose context is similar or complementary: they change the same
  component or contract, share a root cause, or one completes the other. One
  Spec per source is neither required nor preferred, and the author looks
  for such sources before minting a Spec for one.
- **The bound.** Grouping stops where the Spec would outgrow four
  implementation Tasks plus its QA gate. A larger grouped scope splits into
  Specs that each fit, and each source keeps exactly one owning Spec.
- **Extend before minting.** Before a new Finding or Backlog Entry is minted,
  including in Triage, an existing one that is not yet implemented and that
  the observation fits is extended instead: an `open` Backlog Entry is revised
  in place, and an unresolved Finding gains a dated addendum, because a
  Finding stays immutable history.

The bound is a limit on grouping, not a size rule for every Spec. Four
implementation Tasks plus the gate is the limit this repository's Spec
authoring already works to, so the clause writes down a limit in use rather
than inventing one.

## Consequences

- Three new clauses reach every adopter of the Spec workflow and the CONTEXT
  workflow. Each gains a row in the Standard TypeScript Monorepo Source
  Baseline, because the catalog requires one for every clause that profile
  selects, so an adopter's update records each one `retained`. No existing
  clause changes, so no clause is `replaced` or `unaccounted`.
- A grouping that the author judged wrong costs a split later. The Spec judge
  can suggest a grouping (ADR-0209), and the author still decides.
- Triage still resolves one Inbox Entry into one Finding, one Backlog Entry
  or one discard. The entry it resolves into may now be an existing one.
