---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Baseline Plan reports the citations its History Relocations break

A Baseline Plan moves retired repository documents under the History Root in
the same approval as its managed refresh, and nothing tells the adopter which
citations the moves break. Fiscus postponed a whole plan rather than move six
ADRs whose paths other documents cite. Fluxus measured 16 relative links left
broken by one ADR consolidation into `docs/history/adr/`. Neither repository
could see the damage before approving the plan.

Planning now reports it. When a plan carries History Relocations, it reads the
tracked text files of the repository and finds every citation that resolves
today and would not resolve after the plan: a repository path or a Markdown
link destination. It reports one warning per citing file, naming the line, the
cited text, and where the text resolves before and after. The warnings are part
of the plan, so the Plan Digest binds them. An approval therefore covers the
impact the adopter saw, and a change to what a relocation breaks yields a
different digest. Citing files are not preimages, and the warnings never change
what apply writes, moves or verifies (ADR-0071, ADR-0073). A plan without History
Relocations reads no additional file and reports nothing new, so an unchanged
repository still converges to `current` (ADR-0103).

Two alternatives were rejected for now:

- A deferral flag recorded as a Setup Manifest decision, which omits the moves
  from the plan. It would let an adopter take a managed refresh without the
  moves, but it adds a flag, a manifest decision, result fields and skill text.
- A separate digest and confirmation for history moves. It needs a new plan
  schema version, split validation and apply stages, and two transactions
  where ADR-0073 keeps one.

Either can be reopened by a later Spec if showing the impact proves
insufficient. Planning rewrites no citation, because a repository-owned
document is the repository's to edit.
