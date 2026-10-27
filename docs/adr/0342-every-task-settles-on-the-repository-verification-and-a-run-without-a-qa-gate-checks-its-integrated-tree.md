---
status: accepted
created_at: 2026-10-27T15:00:00Z
updated_at: 2026-10-27T15:00:00Z
deprecated_at: null
superseded_by: null
---

# Every Task settles on the repository Verification, and a Run without a QA gate checks its integrated tree

ADR-0182 runs Settlement Checks, the configured repository Verification among
them, only for the Tasks of a Task Graph that has a QA gate Task, "because no
gate of that Run would check these facts". ADR-0257 then told every Agent in
a Daemon-assigned turn to run neither selected Verification, because "the
Daemon runs the Task's declared Verification and the repository Verification
at settlement". In a graph without a QA gate the second sentence is false: no
step of the Run runs the repository Verification, and the Agent was told not
to. On 2026-10-26 an adopter's Run of a `qa: declined` Spec reached Clean
while `make verify` failed on a file the Task created; on 2026-10-27 another
Run of the same adopter reached Clean while a Task broke a test of an earlier
Spec. Both were caught only by the branch's own gate before the Pull Request.

**The repository Verification at settlement applies to every Task Graph.**
For every non-QA Task, with or without a QA gate Task, the Daemon appends the
configured repository Verification as the attempt's last command when
`verification.repository_at_settlement` is enabled and the command is not
empty. A failure is Verification Feedback for the one repair turn, and a
failure on the final attempt settles the Task `failed`, exactly as in a
gated graph. The two in-process Settlement Checks, Spec Consistency and the
authorization audit, stay scoped to a graph with a QA gate Task: they are the
gate's own precondition facts, and ADR-0257's promise names only the
repository Verification.

**A Run without a QA gate checks its integrated tree before it ends Clean.**
Settlement Checks see one Task's tree, so two Tasks of one Wave can each pass
and integrate into a red tree (ADR-0182). A QA gate's precondition finds that
in a gated graph. A graph without one gets an Integrated Repository
Verification: when every Task the Run attempted completed, at least two
completed in the Run, the Task Capacity is above one, and the repository
Verification at settlement is enabled and configured, the Daemon runs the
repository Verification once in the Run Worktree after the last Task
integrates. A failure makes the Run end `Unresolved`, with a reason that
names the command and its diagnostics file, and no Task status changes. At
Task Capacity one each Task starts from the integrated result of the one
before, so the last settlement already verified the integrated tree and the
check does not run.

## Considered options

- Refuse Clean without running anything when a graph has no QA gate. Rejected:
  it would end every such Run `Unresolved`, including green ones, and a Task
  that broke the repository would get no repair turn.
- Only the Integrated Repository Verification, no per-Task command. Rejected:
  a failure found after the last Task has no Agent Session to repair it, and
  the Task that broke the repository is not named.
- Extend the in-process checks to every graph. Rejected: no adopter reported
  them missing, they are the QA gate's precondition facts rather than a fact
  the Baseline promises every Daemon-assigned turn, and widening them is a
  separate decision.
- Tell the Agent to run the repository Verification when the graph has no QA
  gate. Rejected: ADR-0257 measured that agent-side runs cost 284
  agent-minutes in 23 Specs and caught almost nothing; the Daemon is the
  verification authority.

## Consequences

Every Task of a graph without a QA gate now leaves the repository
Verification green, and pays one more run of it at settlement, as gated Tasks
have since ADR-0182. `verification.repository_at_settlement: false` turns off
both the settlement command and the Integrated Repository Verification, and
nothing else. A repository that was already red on entry fails every such
Task, as it fails gated Tasks, unless the frozen authorization names it in
`precondition_repairs`. A repository that never configured
`defaults.verification` gets the built-in `make verify`, as a gated graph
does. ADR-0182 measured that running the command for every Task broke 20
gate-less test fixtures in 2026-09; on 2026-10-27 the same change broke the
one daemon test that pinned the old scope and 21 implement tests whose
fixtures inherit `make verify` without a Makefile, one of them by hanging in
a repair turn. Those fixtures now state their repository Verification, so
the cost is paid once, in test configuration, rather than by every adopter
whose Run reports Clean on a red tree.
