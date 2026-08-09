---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A defect is checked by the stage that can produce it

ADR-0096 already establishes that the QA gate proves machine facts before it
spends an Agent turn. Measured on Spec 0090, the gate did the opposite in
practice: eight of its sixteen matrix rows audited artifacts rather than
product, and the finding that failed the Run came from one of those eight — a
PRD citing the wrong decision record. Reaching it took 461 tool calls and two
context compactions, against 161 tool calls for the average implementation Task
in the same Run.

The cost is not the gate being thorough. It is the gate answering a question
that a file read answers, at the end of a Spec, in the most expensive context
the loop has.

Roundfix therefore places each check with the stage that can produce the defect
it catches. A PRD's Project Constraints are checked when the PRD is written. A
Task graph's coverage is checked when the graph is written. A citation is
checked when the artifact that cites it is written. The checker already reads
one Spec in 0.04 seconds, so the check is affordable at every stage and the
author fixes the artifact while it is still open in front of them.

What stays in the gate is what a file read cannot settle: whether the Spec's
goals work through the surfaces a user reaches, with evidence captured from
those surfaces. That is judgement, and it is what an Agent turn is for.

Two classes deliberately stay behind. A check needing commits — which paths a
Task actually touched against its bounded list — cannot run before the commits
exist, and remains a gate row, executed as a command rather than a judgement.
And a check whose rule is not yet mechanical stays a gate row until it is, so
moving a check earlier is a promise about where it runs, never an excuse to stop
running it.

Leaving the checks in the gate and making the gate cheaper was rejected: it
optimises the moment of discovery instead of moving it, and the author who can
fix the defect in seconds has already moved on by then.
