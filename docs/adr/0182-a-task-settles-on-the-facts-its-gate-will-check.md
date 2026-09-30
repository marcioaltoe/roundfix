---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task settles on the facts its gate will check

ADR-0014 lets a Task settle `completed` when its own declared Verification
passes. ADR-0096 then makes the QA gate prove three machine facts before it
spends an Agent turn: the strict Spec Consistency Check, the configured
repository Verification and the authorization audit of each Task commit. A
completed Task could therefore carry a fact the gate would refuse, and the Run
learned it only at the end. On 2026-09-30 Spec 0187's task_04 settled
`completed` after one second of declared Verification. Two hours into the Run
the gate refused at its precondition, because `make verify-changed` failed on
a contract test that task_04's own edit had broken. ADR-0117 left
commit-dependent checks in the gate because they "cannot run before the commits
exist". They can run on the commit the Daemon is about to create.

In a Task Graph that has a QA gate Task, the Daemon now runs Settlement Checks
for every other Task, after the Task's declared Verification and before it
settles:

- the configured repository Verification, as the last Verification command;
- the refusing Spec Consistency findings that were absent when the Task
  started;
- the gate's authorization audit, applied to the commit the Daemon is about to
  create.

A failed Settlement Check is a Verification failure. Its diagnostics return to
the same Agent Session as Verification Feedback, under ADR-0038's single
repair, and a failure on the final attempt settles the Task `failed`.

The Daemon gains this stage itself, so ADR-0014's rule about hooks is
unchanged: a commit hook still must not be stricter than what the Daemon
verified, and the Daemon now verifies more. The QA gate keeps its whole
precondition as a second check of the integrated tree.

## Consequences

A Task Graph without a QA gate Task keeps settling on declared Verification
alone, because no gate of that Run would check these facts. A prototype that
ran the repository Verification for every Task broke 20 existing tests whose
gate-less fixtures have no such command; conditioning on the gate broke none.

Every Task of a gated graph must leave the repository Verification green. An
author can no longer split a change from the files its contract tests compare.
Each settlement costs one more run of that command, measured at 206 to 223
seconds on this repository, and `verification.repository_at_settlement: false`
turns that one check off. The two in-process checks have no switch: they cost
under a second and the gate refuses the same facts.

Settlement Checks see one Task's tree. Two Tasks of one Wave can each pass and
still integrate into a red tree, which the next settlement or the gate finds.
A failure that predates the Task is still reported to it, because the Daemon
takes no entry observation of the repository Verification; only the Spec
Consistency findings are compared with the Task's start.
