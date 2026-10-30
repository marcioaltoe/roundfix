---
status: accepted
created_at: 2026-10-30T15:00:00Z
updated_at: 2026-10-30T15:00:00Z
deprecated_at: null
superseded_by: null
---

# A row citing an empty grant declares no mutation, and a completed Task never collides

On 2026-10-27 Quaestor reported two refusals that the authoring rules
themselves caused. The `write-prd` and `write-techspec` templates prescribe
the Tooling authority row `applicable — no protected tooling mutation
proposed or authorized`, tell the author to cite the Spec-contained
authorization record, and tell a Spec with no Governed Path to record
`paths: []`. `roundfix spec check` refuses exactly that row with
`SC-TOOLING-UNBOUNDED`, because any cited record makes it read the row as a
mutation, and the record's empty list leaves nothing to bound. Separately, a
corrective Task that declares a file a completed Task also declares makes
`spec.Collisions` report a Wave collision, and `roundfix implement` creates
the Run before the Daemon refuses the Wave, so a Run exists for a graph that
was never runnable.

## Decision

**The row's own words decide, and the cited record must agree.** A Tooling
authority row whose reason says `no protected tooling mutation` declares no
mutation when it cites no record, as ADR-0131 decided, and also when the one
record it cites is an operative grant whose `paths` is the explicitly empty
list ADR-0179 introduced. A row citing a record that lists any Governed Path
still declares a mutation and still needs its bounded files, and a row
without those words keeps today's reading. The checker moves, not the rule:
the templates already state the right reading, and they now also say to end
such a row with `bounded files: none`. A contract test composes the row from
both templates and runs the checker on it, so template and checker cannot
drift apart again.

**A completed Task is outside every Wave.** A Wave holds only Tasks whose
dependencies are completed and that still have to run, so a Task whose status
is `completed` never forms a Wave Collision. Its work is already on the branch
a later Task starts from, because its settlement commit carries the status. A
Task reopened to `pending` counts again. Every reader of `spec.Collisions`
inherits the rule: `roundfix spec check`, the QA gate's precondition, the
settlement checks and the Daemon's pre-dispatch refusal.

**A real collision is refused before a Run exists.** `roundfix implement`
runs the same collision rule in Preflight Validation, on the committed graph,
whenever the Task Capacity is above 1, and exits `2` with nothing recorded.
The Daemon keeps its pre-dispatch refusal as the second line for a graph that
reached the Task cycle another way.

## Consequences

- The Spec 0297 test that expected a corrective Task to collide with a
  completed one now expects no finding; a collision between two Tasks that
  are not completed still makes all four readers agree on the refusal.
- At a Task Capacity of 1 Preflight Validation refuses nothing, as the
  Daemon does, because a serial plan has no Wave.
- A row whose bounded files read `none (...)` keeps counting as bounded,
  as today; this decision does not tighten that reading.

## Alternatives rejected

- Change the templates to `not applicable`: the constraint binds every Spec
  (ADR-0131), and adopters following the template today would have to
  rewrite rows the checker should accept.
- Stop treating a cited record as a mutation: a row that says "no mutation"
  while its record grants Governed Paths would pass unbounded.
- Refuse only in `spec check`: nothing makes the operator run it before
  `implement`, and the Delivery Queue starts `implement` directly.
- Create the Run and refuse in the Daemon, as now: the Run, its worktree and
  its record exist for work that could never start.
