---
status: accepted
created_at: 2026-08-02T00:00:00Z
updated_at: 2026-08-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The QA gate is authored into the graph, not requested per run

The Implement Command's `--qa` parameter makes the gate a property of the
invocation, so a Task Graph that grows after a gate has reported leaves no
structural trace: on Spec 0057 three gates ran against three progressively
larger graphs and read from outside as three ordinary cycles rather than as one
decomposition that was wrong twice. The gate therefore becomes a Task the Spec
authors — a terminal node depending on every leaf — and the parameter is
removed, so what runs is what the graph says and appending work after the
terminal is visibly an insertion that invalidates the reported result. Keeping
the flag and enforcing the existing advisory cap was rejected because an
advisory cap on a flag is precisely what allowed three closings to pass
unnoticed, and making the gate implicit for every Spec was rejected because a
Spec with no behavioral surface to observe would then pay for a gate it does
not need. The Daemon's existing withholding — no gate begins while any Task is
unsettled — is unchanged and remains correct; this moves where the gate lives,
not when it may begin.
