---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec proves what it claims, and every active Spec is held to it

Four authoring gaps were measured in another repository that uses these
skills. Two Tasks with green Verification left the repository red, because
each changed a route that a suite in another directory consumed, and only the
terminal QA gate's full run found it. A Spec ordered a component that already
existed, with its tests, in a file its own TechSpec cited; the implementing
Agent reconciled the contradiction by weakening the existing behavior and
deleting tests, and every gate stayed green. A Spec raised a timeout to exactly
the value an accepted ADR had rejected, which cost four gate cycles. A TechSpec
mapped every PRD promise in its Coverage Map while deliberately delivering one
of them in a different shape. On 2026-10-11 the maintainer chose "Skills + gap
no spec check".

**Each claim carries its proof as a section.** The TechSpec records Absence
Evidence: for each component it declares new, the search command that proved
the component absent and the result, or `None.` with the reason no component
is new. It records ADR Alternatives: one line for each ADR listed under Active
ADR obligations, on how the Spec relates to the alternatives that ADR
rejected. Its Coverage Map records the Coverage Cross-Read: a
`Deviations from the PRD:` line that is `None.` or lists each PRD unit the
TechSpec delivers in another shape. Each non-QA Task that is not completed
carries a Consumer Declaration: the suites that consume each Shared Contract
the Task changes, named in its own Verification or in the Verification of an
Integration Checkpoint between the Task and the QA gate, or `None.` with the
reason it changes no Shared Contract.

**The check reads structure, never meaning.** `roundfix spec check` reports a
missing section, a missing line, an item without its command or result, an ADR
without its line, a deviation that names no PRD unit, and a named consumer
suite that no Verification runs. Each is a gap, so `--strict`, the Delivery
Plan, queue revalidation and the QA gate refuse it. It never judges whether a
search was the right one, whether a relation is true, or whether a deviation is
justified. ADR-0093: "it never judges whether a decision is correct". A
finding's fix names the missing section or line and the authoring-skill rule
that requires it, so an author can act on it without reading the skill first.

**Every active Spec is held, with no horizon.** The maintainer decided on
2026-10-12 that these gaps apply to all active Specs, with no exemption by the
date a PRD was committed. Archived Specs stay out of scope because nothing under
`docs/history/` is validated as live work. A detector still skips when its
artifact is absent. ADR-0094: "is skipped, not failed, when one of them is
absent". A completed Task is not read for its Consumer Declaration, because it
never runs again.

## Considered options

- A horizon that holds only Specs whose PRD was committed after the skills
  changed, as the receipt and glossary gaps do. Rejected by the maintainer: a
  Spec that has not run yet is the Spec most likely to repeat the measured
  defects, and the adopter notes name the sections it must add before its next
  Run.
- Judging the content, through a deterministic heuristic or the advisory
  judge. Rejected: a semantic gate is the class ADR-0093 excludes, and the judge
  never gates.
- Re-running each Absence Evidence command at check time. Rejected: the command
  describes the tree the author searched, the tree moves, and a check that
  executes authored search commands would need the committed-provenance rules
  that Verification already carries.
- Running every consuming suite at each Task's settlement through the
  repository Verification alone. Rejected: the selective gate picks tests by
  changed path, and a consumer in another directory is exactly the test it does
  not pick.

## Consequences

An active Spec whose TechSpec lacks the three sections, or whose pending Tasks
lack a Consumer Declaration, fails `spec check --strict` until its author adds
them. In this repository the Spec that ships the check is queued after the
other active Specs, or the operator completes their sections first, because
the active-corpus contract tests hold every finding count at zero. A Verification
that names a consuming suite still has to fail before the work, so the suite
shares a command with a condition only the Task's work satisfies.
