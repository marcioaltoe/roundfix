---
status: accepted
created_at: 2026-10-05T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A failed check is judged on a merge with the current default branch

On 2026-10-05 a Pull Request's CI failed because of a defect on the default
branch. After the fix merged, a re-run of the failed job still tested the old
merge, because a re-run keeps the original event's commit, and a reopen made
four and a half minutes later still received GitHub's old test merge. The
Delivery Queue meanwhile parked the item as `flaky-check` and then as
`checks-failed`, although neither verdict was about the current default
branch.

CI's pull request job therefore checks out the pull request's head commit with
every branch fetched, merges the default branch tip it fetched at job start,
and records that tip as a notice annotation titled `tested-base`. Every
attempt, including a re-run, tests the default branch as it stands when the
job starts, and verifies the changed sets against that tip. Pushes to the
default branch keep their full Verification.

Before the Delivery Queue classifies a failed check, it reads the failing
job's `tested-base` annotation and fetches the delivery remote's default
branch. A failure whose tested base does not contain that tip is stale: the
queue re-runs the failed jobs once for each default-branch tip it sees, logs
the re-run and keeps polling, and never parks a stale failure as
`flaky-check` or `checks-failed`. A run it re-ran as stale gets the one
outside-change re-run on its first failure against the current tip. A
failure without the annotation keeps the rules it had before, and the check
timeout still bounds the wait.

## Considered Options

- Close and reopen the Pull Request: rejected, because GitHub computes its
  test merge lazily, and the reopen of 2026-10-05 still tested the old merge.
- Check out the dynamic `refs/pull/<n>/merge` ref on every attempt: rejected
  for the same lag, and because the queue still could not tell which default
  branch a failure tested.
- Compare a run's start time with when the default branch moved: rejected,
  because a commit's date is not its push time and the reopen showed that a
  later run can still test an older merge.

## Consequences

CI tests a merge it computes instead of GitHub's test merge. A head that
conflicts with the default branch fails the merge step, as GitHub reports the
same conflict to the queue first. The queue's memory of which runs it re-ran
as stale lasts one owner process, so after a restart a failure past its first
attempt on the current tip parks as before.
