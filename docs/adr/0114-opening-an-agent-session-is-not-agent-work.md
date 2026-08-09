---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Opening an Agent Session is not Agent work

A Fallback Selection may switch ACP Runtime automatically only while Agent work
has not begun. The guard is right: never swap models over partially finished
work, because the second model inherits a state the first one built and neither
account of it is complete.

The boundary was drawn in the wrong place. Roundfix emitted `AGENT_WORK_STARTED`
when it began opening the session, so on 2026-08-08 an exhausted Codex quota —
the adapter printed its usage limit and exited, and the whole Run lasted nineteen
seconds — was classified as a Batch failure after work started. The configured
`codex → claude` chain never became eligible. A Spec stopped for four days on a
quota that a proven, configured alternative could have absorbed immediately.

Roundfix therefore treats opening an Agent Session as selection, not work. The
signal that makes a Fallback Selection ineligible moves to the first Agent turn
that could have changed something. A runtime that refuses to serve — exhausted
quota, failed authentication, an adapter that will not start — is a selection
failure, and the Fallback Chain exists for precisely that case.

This narrows when the chain is ineligible; it does not widen what the chain may
do. ADR-0050 still keeps Fallback Chains inactive until Run creation and forbids
preflight from substituting anything, and every tuple in a chain is still proven
before the Run starts. What changes is that a Run which produced no Agent turn
can still reach its second choice.

Emitting the signal at the same point and special-casing quota errors was
rejected: it would fix the measured instance and leave the boundary wrong for
authentication, adapter startup, and every future refusal an adapter learns to
report.
