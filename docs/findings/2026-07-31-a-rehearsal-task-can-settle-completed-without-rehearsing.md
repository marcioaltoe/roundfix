---
date: 2026-07-31
status: open
spec: —
---

# A rehearsal Task can settle `completed` without rehearsing

## What happened

Spec 0060's `task_03` existed to prove an instruction-level gate actually fires:
archive passes when a Spec is self-contained, fails naming an injected stale
link, does not false-positive on prose, and passes trivially with no index.

It settled `completed` having run none of those cases.

Two defects combined.

**The requirements contradicted each other.** `task_03` R1 said "MUST NOT commit
the rehearsal fixtures"; R2 said "MUST prove the move preserves history:
`git log --follow` on the moved file reaches its pre-adoption commits".
`git log --follow` reads committed history, so with the move uncommitted the
destination path exists in no commit and there is nothing to traverse. The
Implement Agent demonstrated this in a disposable clone rather than asserting
it, then declined to claim any acceptance criterion — the correct behavior.

**The Verification could not detect the omission.** `task_03` declared
`make verify` and `git status --porcelain`. Both pass most easily when nothing
happened: an empty tree is exactly what "no rehearsal ran" produces. The
Daemon ran both verbatim, they passed, and the Task settled `completed`.

The Spec's own TechSpec had named this risk — "a gate written as prose can be
skipped silently… this is the Spec's main weakness" — and it materialized on
the Task written to prevent it.

## Why it did not ship broken

The `qa` Agent Session independently rebuilt the rehearsal in a disposable
clone and exercised every case, going beyond `task_03`'s requirements: it added
reference-style links alongside inline ones and proved `qa_override: true`
cannot bypass an indexed self-containment failure. Evidence with real exit
codes is in the Spec's QA evidence directory. QA recorded the same
contradiction as row QA-11, blocked it as environmental, and named the fix.

So the safety net held — but it held one layer later than intended, and only
because the QA gate happened to be thorough.

## The general defect

**A Verification whose commands pass on an empty tree cannot prove work
happened.** This is the negative-space case of "verify the class, not the
case": absence assertions must distinguish "nothing was done" from "nothing
was wrong". A Task whose product is *evidence* rather than *code* needs a
Verification that reads the evidence — a non-empty `## Result`, an expected
exit code recorded, a named artifact present — not one that inspects a tree
the Task is required to leave clean.

## Suggested repair

1. Scope the prohibition to the repository: a rehearsal may commit freely
   inside a disposable clone that is deleted afterward; what must never be
   committed is a fixture *into this repository*. That resolves R1 against R2
   with no loss of safety.
2. Give evidence-producing Tasks a Verification that greps their own recorded
   evidence for the required outcomes, so an un-run rehearsal fails.

Both are Task-authoring rules, so the natural owner is `write-tasks` and the
Verification guidance in `docs/agents/autonomous-work.md`.
