---
status: accepted
created_at: 2026-09-24T00:00:00Z
updated_at: 2026-09-24T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Verification runs past a failure only when declared independent

A Task's Verification commands normally run in order and stop at the first
failure. The one repair turn allowed by ADR-0038 can therefore receive only
the first failure even when later commands were already failing for their own
reasons.

Tasks may declare `verification: independent`. For such a Task, the Daemon
runs every Verification command after a failure and sends every deterministic
failure, with its diagnostics, to the repair turn. The retry remains bounded
as before and reruns the complete command list. A Task without the declaration
keeps the existing first-failure behavior.

Independence is declared, never inferred from command text. Shell text cannot
reliably establish whether a later command depends on setup performed by an
earlier one; inferring independence would guess at command ordering and make
the executor responsible for the Task author's intent. The declaration keeps
that responsibility with the author and makes the all-or-none scope explicit.

The trade-off is that an incorrect declaration can run commands after a setup
failure. That risk is accepted because the author is the one who knows whether
the Task's commands are independent.

## References

- [ADR-0038: Daemon allows one Verification repair](0038-daemon-allows-one-verification-repair.md)
- [Spec 0158 technical specification](../specs/0158-daemon-verification-and-access-readiness/_techspec.md)
