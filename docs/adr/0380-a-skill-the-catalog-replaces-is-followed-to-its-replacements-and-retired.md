---
status: accepted
created_at: 2026-11-09T12:00:00Z
updated_at: 2026-11-09T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A skill the catalog replaces is followed to its replacements and retired

Amends ADR-0206 (2026-11-09) for a skill the catalog itself replaces.

ADR-0206 drops a skill the catalog removes from every module, trigger and
snapshot, and never follows it to a skill with a similar purpose under
another name. On 2026-11-09 the catalog `marcioaltoe/skills` removed
`inngest` and, in the same change (pull request 156, commit `e45466d2`),
named seven official skills from `inngest/inngest-skills` as its
replacement: `inngest-setup`, `inngest-durable-functions`, `inngest-steps`,
`inngest-events`, `inngest-flow-control`, `inngest-cli` and
`inngest-middleware`. The `typescript` setup list names the seven where it
named `inngest`. The Baseline's `backend` module still required and
dispatched `inngest`, so the Catalog Skill Bump stopped at it (ADR-0360),
and an adopter's installed `inngest` named a catalog path that no longer
exists. The maintainer decided the same day: "Spec fix na Baseline: a
capacidade opcional Inngest instala as 7 skills e despacha pelos gatilhos
sugeridos (funções/steps, eventos, flow control, setup, CLI, middleware);
fixtures atualizadas."

## Decision

**A skill the catalog replaces is followed to its replacements, and the
replaced skill becomes a Retired Skill.**

- A replacement is one the catalog names itself, in the change that removes
  the skill, and that the maintainer approves. Roundfix never infers one from
  a similar name or purpose; a skill removed without a named replacement is
  still dropped, as ADR-0206 decides.
- The module that required the replaced skill requires every replacement and
  dispatches each one through its own trigger. The Catalog Skill Bump then
  pins the setups that list the replacements (ADR-0360).
- The replaced skill joins the Retired Skills, so no module, bundle or
  snapshot names it again, and an approved `roundfix baseline update` deletes
  an adopter's copy it proves unedited and keeps an edited or unrecorded one
  with the paths to delete by hand (ADR-0350). The skills stage installs the
  replacements the module now requires.
- `inngest` is the first such skill. The `backend` module requires the seven
  `inngest-*` skills and dispatches them by the catalog's suggested mapping:
  functions and steps to `inngest-durable-functions` and `inngest-steps`,
  events and fan-out to `inngest-events`, concurrency, throttling and
  debounce to `inngest-flow-control`, project setup to `inngest-setup`, the
  local Dev Server and CLI to `inngest-cli`, and logging, tracing and
  dependency injection to `inngest-middleware`.

## Considered Options

- Dropping `inngest` without a replacement, as ADR-0206 prescribes for a
  removed skill: rejected, because the catalog named the replacement itself
  and adopters with Inngest would lose their only Inngest guidance.
- Requiring only some of the seven: rejected, because the catalog's setup
  list names all seven and a partial set would make the snapshot list skills
  no module dispatches.
- Listing `inngest` with the removed skills that adopters keep until they
  delete them by hand: rejected, because its lock entry names a catalog path
  that no longer exists, and ADR-0350 already deletes a copy proven unedited
  without losing an edited one.
- Requiring the seven only when the repository declares Inngest: rejected,
  because the `backend` module required `inngest` for every repository that
  selects it; a conditional requirement is a change of its own.

## Consequences

After the release that carries this decision, an adopter whose profile
selects the `backend` module sees the seven `inngest-*` skills as missing
required skills and the update installs them, and sees `inngest` under
`Skills retired`. The approved update deletes the copy the lock proves
unedited.
