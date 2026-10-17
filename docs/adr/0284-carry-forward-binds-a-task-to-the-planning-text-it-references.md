---
status: accepted
created_at: 2026-10-17T20:40:00Z
updated_at: 2026-10-17T20:40:00Z
deprecated_at: null
superseded_by: null
---

# Carry-forward binds a Task to the planning text it references

Task Carry-Forward refuses the whole set when a settled Task's declared inputs
moved since its settlement, and the PRD, the TechSpec and the Task Graph
manifest are inputs of every Task. On 2026-10-17 three amendments of planning
files refused every settled Task of their Spec: Spec 0272 after pull request
#551, which added a corrective Task node and Tooling authority rows, and Spec
0274 after pull request #559, which moved two contract tests between files in
the TechSpec. Each time the operator carried the settled Tasks by hand: cherry-
pick the settlement commit, append `## Carry-forward provenance`, amend.

## Decision

**A settled Task binds only the planning text it names.** When the PRD, the
TechSpec or the Task Graph manifest of the Task's own Spec moved, carry-forward
compares the Task's **Binding Text** at the settlement parent and on the
staged target instead of the whole file:

- in the Task Graph manifest, the Task's own graph node and its row of the
  projection table;
- in the PRD and the TechSpec, every section, from its heading to the next
  heading of the same or a higher level, that a unit of the Task's
  `## References` line for that file names, singular or plural, with its
  numbers dropped.

The Task stays carriable when its Binding Text is unchanged in each moved
planning file, and the carried candidate lists those files as amended inputs.
A Task's own Task file, `AGENTS.md`, `CONTEXT.md`, the `implement-task` skill
and every `## Context` path still refuse on any change. A planning file the
Task's References do not name, or a unit that names no heading, binds the
whole file. A Task whose named section changed is refused with today's reason.

## Considered Options

- Item-level binding, comparing only API Contract 2 inside the API Contracts
  section, was rejected: the numbered units are list items whose boundaries
  and renumbering are not mechanical, and a wrong cut would carry a Task whose
  contract changed.
- Ignoring every planning file was rejected: a Task's References are its
  contract, and a changed contract must run again.
- Carrying the unaffected part of a refused set stays rejected (ADR-0170); a
  Task whose Binding Text changed still refuses the whole set.

## Consequences

- An amendment of Project Constraints, Decisions, the Coverage Map, another
  Task's node or the manifest's prose no longer refuses a settled Task that
  does not name them.
- An amendment inside a section a settled Task names, such as Testing Approach,
  still refuses it, so an author who amends that section expects the Task to
  run again.
