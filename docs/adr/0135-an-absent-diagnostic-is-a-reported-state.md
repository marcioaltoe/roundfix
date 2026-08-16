---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-16T09:00:00Z
updated_at: 2026-08-16T09:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# An absent diagnostic is a reported state, not an empty message

A Verification that fails after redirecting its output produces no captured text,
and the repair prompt then carries a command, an exit status, and nothing else. An
empty prompt and a prompt about nothing are indistinguishable to the reader who
has to act on them, and in the measured case the Agent spent its one repair turn
rewriting its own Task file with a cause it had deduced — reasonable behaviour
given no information. The feedback therefore states that the command produced no
output, and where that output was redirected if it was. ADR-0111 already makes an
unobserved Verification unknown rather than a verdict; this is the same principle
one layer out, applied to the text an Agent reads instead of to the classification
the Daemon records.
