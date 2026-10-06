---
type: fix
status: promoted
created: 2026-10-06
spec: 0239-a-glossary-every-spec-keeps-current
---

# The CONTEXT-driven loop does not keep `CONTEXT.md` current

## Problem

The maintainer created Roundfix with `CONTEXT.md` (the glossary) and the ADRs
as the base of the CONTEXT-driven process. On 2026-10-06 the maintainer
observed that `grilling` and `domain-modeling` are not being executed and that
`CONTEXT.md` is not kept current during implementations.

Measured on 2026-10-06:

- **No glossary updates.** `CONTEXT.md` (926 lines, 225 terms) last changed on
  2026-10-02 (`0ac53081`, Spec 0214). Eighteen Specs were delivered after it
  without touching it.
- **Terms reported and dropped.** Authoring and QA reports named new terms
  and left them out, for example: Light Tier, Light Spend Log, Evidence
  Manifest, tested base, refused permission, merge evidence, delivery commit,
  item branch, Jev ceiling, network-denied row, Merge Observer, stage key.
  Several archived QA reports flag a glossary candidate.
- **Operator-side cause.** The operator's authoring briefing for unattended
  Specs said "Never edit `CONTEXT.md` … list any new glossary term in your
  report", and nobody applied the reported terms. The briefing was corrected
  on 2026-10-06: a Spec now adds its own terms through `domain-modeling` as a
  requirement of its docs Task.
- **Route-side cause.** The autonomous authoring route (write-prd decided
  autonomously → write-techspec → write-tasks) never reaches `grilling`,
  `grill-with-docs` or `domain-modeling`. `docs/agents/skill-dispatch.md`
  triggers `domain-modeling` only on "defining or changing domain vocabulary",
  and no gate checks that a Spec that introduces vocabulary updated the
  glossary. In the Run Database, `domain-modeling` was read in 25 of 220 Runs;
  `grilling` and `grill-with-docs` in none.

## Expected

1. **A gate.** `roundfix spec check` (and archive) detects domain terms a
   Spec introduces (capitalized concepts in the PRD/TechSpec Ubiquitous
   Language or Glossary section, new ADR terms) that are missing from
   `CONTEXT.md`. It reports them as a finding the Spec resolves by adding the
   term in a Task or by declaring it not a domain term. The Jev judge can
   advise on "is this a domain term", but the gate is deterministic.
2. **Skills in the route.**
   - `write-prd` and `write-techspec` call `domain-modeling` whenever they
     name a new concept, even on the autonomous route.
   - `write-tasks` adds a glossary requirement to the docs Task.
   - The QA gate checks the terms against `CONTEXT.md`.
   - `grilling` and `grill-with-docs` stay as the interactive entry and are
     used when a person drives the authoring.
3. **The Baseline.** The `context-workflow` module's clauses state that
   `CONTEXT.md` and the ADRs are the base, and that a Spec updates the
   glossary for the terms it introduces. Adopters get the same rule. Measure
   whether their `CONTEXT.md` files went stale the same way.
4. **Catch-up.** Add the terms dropped since 2026-10-02 (the list above,
   checked against the delivered code and ADRs) to `CONTEXT.md` in one docs
   change.

## Notes

Related: `2026-10-06-retire-the-skills-nobody-uses-or-jev-replaced.md` (the
maintainer keeps `grilling` and `grill-with-docs` and asked for this check).
Resolve after the current cycle (0235–0237). Do it before the history
clean-up, so the knowledge that history loses is first in the glossary and
the ADRs.
