# Failed Batches do not halt the resolve cycle

A Batch fails when its Agent errors or the verification command fails after
it. Roundfix marks every Review Issue assigned to a failed Batch with
`status: failed`, records the failure as Run Events, and continues the
resolve cycle with the next Batch. Review Issues the Agent leaves unsettled
(`pending`, `valid`) are individually marked failed before verification.
`failed` is a settled Batch outcome but not a Terminal Review Issue status:
failed issues stay Unresolved, keep blocking Final Push, and are retried when
a later Round downloads their still-open Review Source threads.

A resolve cycle that completes with Unresolved Review Issues ends the Run in
the terminal state `Unresolved`, distinct from `Failed`, which means the Run
itself broke. Stop Requests and infrastructure errors still halt the cycle.
Watch Rounds run exactly one resolve cycle each and never retry failed issues
inside the same Round; a Round whose cycle settles nothing ends the Run as
`Unresolved` instead of repeating identical Rounds.

Failed Batches create no Batch commit and resolve no Review Source threads.
Worktree changes from a failed Batch are preserved, never reverted, and
excluded from later Batch commits by the before-snapshot rule.
