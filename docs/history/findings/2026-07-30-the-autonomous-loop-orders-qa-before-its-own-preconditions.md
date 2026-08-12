---
status: done
created_at: 2026-07-30
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# The autonomous loop orders QA before its own preconditions (2026-07-30)

The loop discipline shipped in `rule.autonomous.loop` and in
`docs/agents/autonomous-work.md` prescribes: implement the Task Graph, request
the QA gate once when the graph closes, open the Pull Request, watch until
Clean, merge. Spec 0053 needed four QA cycles to discover that this order
cannot reach `pass` for any Spec whose acceptance observes its own Pull
Request, and the correction arrived in two stages.

## Stage 1 — QA cannot run a journey whose subject does not exist

The first gate returned `partial` with `rows_blocked_environment: 1`. The only
blocked row was the Pull Request observation journey, blocked because no Pull
Request existed. Nothing was wrong: the qa-gate Skill requires an
environment-blocked row to record its cause **and equivalent observed
evidence**, and "target lookup returned `[]`" proves the Pull Request does not
exist rather than substituting for the journey.

## Stage 2 — QA cannot accept a Pull Request nobody has processed

Opening the Pull Request and rerunning was still not enough. The fourth gate
reported that the Pull Request carried `CHANGES_REQUESTED`, five unresolved
review threads, and no review-artifact descendant, and correctly refused to
declare Merge-Ready acceptance. The Pull Request has to be not merely open but
already Clean.

## Root cause

The loop treats the Pull Request as publication, which follows verification.
For a Spec whose QA matrix reaches a review surface, the Pull Request and its
review are **preconditions** of verification. The loop's own
`clause.autonomous.loop-03` names this hazard class — a requirement that states
its objective without stating where its inputs come from — and the loop
committed it about itself.

## Action / suggestion

The correct order is:

> implement the Task Graph → open the Pull Request → watch until Clean →
> request the QA gate once → merge

`clause.autonomous.loop-01` should read "request the QA gate once, after the
Task Graph closes **and every surface its acceptance observes exists and is
Clean**", and the loop section in `docs/agents/autonomous-work.md` should carry
the reordered sequence. A Spec whose QA matrix touches no review surface is
unaffected by the reorder, so one order serves both.

The cost of the reorder is that a Pull Request is published before QA passes.
That is already true of every human workflow in this repository, and the Pull
Request is opened from a Task Graph that closed green, so it is not a draft.

## What worked — keep

- Requesting QA once per attempt was right. Each cycle cost roughly 25 minutes
  and returned a precise, actionable diagnosis; re-requesting after each
  corrective Task would have spent several cycles reaching the same conclusion.
- The typed blocked-cause counts introduced by Spec 0053 are what made stage 1
  legible. `rows_blocked_environment: 1` with `rows_blocked_finding: 0` says
  "nothing is wrong with the work" in two lines, which the previous single
  verdict scalar could not express.
- The gate refused to claim `pass` on a journey it could not run, under a rule
  that explicitly permits `pass` when equivalent evidence exists. That boundary
  behaved exactly as designed in all four cycles.

## Routing — 2026-08-01

Routed to [Spec 0065](../specs/0065-loop-order-and-verification-honesty/_prd.md) on 2026-08-01.
