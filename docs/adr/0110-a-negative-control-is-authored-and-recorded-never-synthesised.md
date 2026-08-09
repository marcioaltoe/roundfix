---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A negative control is authored and recorded, never synthesised

A gate earns trust by discriminating, and discrimination takes more than one
green run. The published pattern is three observations: a positive control that
passes on correct work, a negative control that must fail against a known defect,
and an observability control that must not report success when the command never
reached the surface it names. Roundfix had only the first.

ADR-0109 adds the observability control mechanically, because the Daemon can run
it without author cooperation. The negative control cannot be obtained the same
way. Manufacturing a defect means mutating the repository under test — reverting
a hunk, corrupting a fixture, disabling a branch — and doing that automatically
inside a Run that also commits is an unacceptable class of action for a tool
that holds a maintainer's working tree.

Roundfix therefore requires the negative control to be **authored** and records
whether it was exercised. A Task may declare a `## Negative Control` section
naming the defect its Verification must catch; the Daemon carries that
declaration and reports it beside the verdict. Gate health becomes a recorded
fact — which surface was observed, how many negative controls were exercised,
when the gate's own test last moved — rather than an assumption renewed by every
green run.

Synthesising the control through mutation testing was rejected for this Spec, not
forever. It is the natural mechanism and the literature treats it as such, but it
needs an isolated tree, a mutation budget, and a story for mutations that are
semantically equivalent. None of those exist yet, and adding them inside a Spec
about gate honesty would make the gate for that Spec depend on the machinery it
is introducing.

Requiring the declaration on every Task was also rejected. A Task that declares
no negative control is recorded as having none, which is a weaker gate stated
honestly. Forcing the section would produce declarations written to satisfy a
parser, which is the failure this decision exists to prevent one level up.
