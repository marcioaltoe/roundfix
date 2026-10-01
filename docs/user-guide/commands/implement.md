### implement

```bash
roundfix implement --spec <slug> [--detach]
```

Executes a Spec's Task Graph in a Run Worktree as one Run — spec Runs keep
worktree isolation. The scheduler executes the current Wave up to
`worktree.concurrency` at a time (default `2`; `1` keeps sequential behavior),
while `verification.concurrency` independently limits concurrent Task
Verification attempts (default `1`). Concurrently running Tasks use Task
Worktrees, and each completed Task creates one commit on the Run Branch. Both
capacities apply only within this Implement Run; they do not coordinate other
Runs, CI, or external processes. It resolves `specs.root` once from the user's
checkout.
The gate is the Spec's authored terminal `qa` Task; no flag requests it.
`implement.auto_push: true` makes a Clean Run push its branch
upstream. Integration Pending, Unresolved, Failed, Stopped, and failing-QA
Runs never push.

Before creating a Run, `implement` inspects prior terminal Runs for the same
Spec in the current repository. It considers Runs with outcome Stopped or
Unresolved when their recorded Run Worktree is present. Preflight Validation
refuses only when a prior Run's complete candidate set would carry: every
candidate in that set passes Task Carry-Forward's proofs. A set with any
refusing candidate is reported and implement proceeds. A Task already
completed on the checkout is reported as `already completed; nothing to carry`,
never refuses the set, and is omitted from the Tasks the Preflight refusal
names. A Run whose candidates are all already completed produces no
not-available note. The refusal exits `2`,
leaves stdout empty, creates no Run or Agent Session, and writes no Git or Run
Database state. It names the selected Run and the Tasks it would recover, then
gives the exact recovery command:

```text
Preflight failed

Reason:
  Run <run-id> holds proved Tasks available for Task Carry-Forward: task_01, task_02

No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

Next action:
  recover them with `roundfix reconcile <run-id> --carry-forward`
```

When more than one prior Run qualifies, `implement` names the Run with the
largest carriable Task set; a tie goes to the newest Run. When an inspected
Run's stranded work is not carriable as a complete set, `implement` reports
each refusal reason on stderr and proceeds to create the Run. If the inspection
itself fails, it reports the failure and also proceeds. A prior Run whose
recorded Run Worktree is gone is skipped by this check.

The profile-led default is:

```bash
roundfix implement --spec <slug> --detach
```

Each Task owns a Task Type-selected Agent Session. During an Implement Run the
Daemon is the only Task-status writer: it writes `in_progress`, receives
implementation-ready work from the Agent, runs the Task's complete
`## Verification` sequence verbatim, and writes the terminal status. The Agent
may run focused checks and record their evidence, but it does not run the
declared Task Verification, edit status, or settle the verdict.

The Implement Run Budget starts with one `budget.max_run_duration` allowance
from Run creation, then renews at each Task settlement. The renewed allowance
bounds the next Task or QA gate and the integration, push, and cleanup work
after the Task cycle. When it expires, Roundfix cancels the Run's Agent
Sessions, settles `BudgetExceeded`, and keeps the Run Worktree and Run Branch
for recovery.

A deterministic Verification failure releases Verification Capacity before
the Daemon sends diagnostics to the same Agent Session. The Agent receives one
Verification Feedback repair turn, then the repaired Task queues and acquires
Verification Capacity again for its final Daemon attempt. Any formatter, test,
Skill synchronization, or build failure in the declared gate blocks
settlement.

#### Settlement Checks

For every non-QA Task in a Task Graph with an authored QA gate, the Daemon runs
three checks as part of the Task's settlement attempt, in this order:

1. The repository Verification runs at settlement as the attempt's last
   command.
2. `settlement check: spec consistency` checks for Spec Consistency findings
   introduced by the Task.
3. `settlement check: authorization` checks the prospective Task commit against
   the frozen authorization record.

Both in-process checks run after the attempt's commands, whatever their result,
and before the attempt's verdict.

A failed check returns as Verification Feedback for the single repair turn. A
final failure settles the Task `failed`. An in-process failure uses the reason
`Settlement check failed: <label>: <first diagnostic line>; diagnostics: <path>`;
a failure from the repository Verification keeps the existing Verification
failure reason.

The `verification.repository_at_settlement` User Config or Project Config
switch turns off only the repository Verification at settlement. It does not
turn off the Spec Consistency or authorization checks. Settlement Checks inspect
the Task's own tree, so a parallel Wave is not checked as one integrated tree;
the existing precondition-repair limit also applies when the repository is red
on entry.

During final QA, the mechanical authorization audit reads each governed Task
commit's grant at its fork point first. When that grant does not cover the
commit, the audit can use the grant the Task ran under: the record in the Task
commit's parent, but only while the delivery target carries byte-identical
content at the same path. The audit reports the latest delivery-target commit
that established that content as the authorizing revision. A parent-only
record, an older record that the delivery target later narrowed or revoked,
and a Task commit that edits its own record remain refusals.

Exit `75` from a project-authored Verification wrapper is the sole Temporary
Verification Failure signal. Roundfix retains its diagnostics and grants that
Task one exclusive retry, which waits for all other Verification attempts in
the Run to drain and then consumes the entire Verification Capacity. The retry
does not consume the Agent repair. A second exit `75` exhausts the retry and
fails the Task; Roundfix never classifies a failure from log text, timing,
ports, package names, or framework messages.

stdout is one line per Task in Task Graph order — failed and skipped Tasks are
followed by one indented `  reason: <one line>` naming the failed step (for
Verification failures: the command, exit status, and diagnostics path) — then
the QA line when requested and one outcome line:

```text
task_01 completed — first task
task_02 failed — second task
  reason: verification failed: make verify (exit status 2); see <path>
Unresolved: 1 completed, 1 failed, 0 skipped, 0 pending.
```

Task status vocabulary is normalized on reload: `done` and hyphen/space
variants of the canonical statuses map to canonical form and the task file is
rewritten; anything else still fails validation. A Task whose commit contains
no change outside the Spec Root still settles `completed`, with one stderr
warning and one Run Event marking the no-op.

For every completed non-QA Task commit, the Daemon writes ordinary committed
paths that the Task did not declare as `interface:` or `creates:` under a
Daemon-owned `## Recorded paths` section at the end of the Task file. The Task
commit event carries the same paths as `recorded_paths`. A recorded path is
disclosed after the change and reserves nothing against another Task; a
Governed Path remains subject to its authorization and is never recorded.

Daemon Task and QA commits stage only repository paths that do not cross a
symbolic link; dropped paths are journaled and warned
(`roundfix: task file <path> kept outside the repository; committed without it`).

