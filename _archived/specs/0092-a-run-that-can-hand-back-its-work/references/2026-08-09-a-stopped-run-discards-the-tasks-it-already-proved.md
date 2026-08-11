---
type: fix # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-09
spec: 0092-a-run-that-can-hand-back-its-work # Spec slug when status: promoted
reason: null # required when status: declined
---

# A stopped Run discards the Tasks it already proved, then blocks the next one

## Opportunity

An Implement Run commits each settled Task inside its Run Worktree and
integrates into the user checkout only when the Run reaches Clean. A Run that
stops before then — by Stop Request, or by a Task failing after earlier ones
succeeded — leaves every completed Task behind. The checkout still reads
`status: pending` for work that ran, passed its Verification, and was committed
in the worktree.

The next Run therefore re-executes it. Not re-checks it: re-executes it, with a
fresh Agent turn and a fresh Verification cycle.

There is a second face to this, and it is the sharper one. `roundfix reconcile`
classifies a stopped Run's branch `unintegrated` and its action `preserve`; it
never integrates, by design. Branch Integrity Preflight then refuses to create
any new Run while an unintegrated Run Branch exists for the same head branch. So
a stopped Run does not merely waste the work it holds — it holds the workflow
shut until a human merges or deletes the branch by hand, and neither `reconcile`
nor any other Roundfix command offers that step.

## Value

Measured on 2026-08-09 implementing Spec 0089. Four Runs were started. Task 01
completed in every one of them and Task 02 in three, all against the same
unchanged inputs, because each Run stopped at Task 03 on an authoring defect
that had nothing to do with either. Roughly twenty minutes of Agent work per
restart, repeated, plus the tokens for turns whose output already existed on
disk in the previous Run Worktree.

The deadlock was measured the same day. Three stopped Run Branches from those
attempts refused two consecutive `roundfix resolve --pr 143` invocations:

```
Branch Integrity Preflight refused pending Run Branch work for PR Head Branch "ma/specs-0082-0083".
- branch=roundfix/run-run_20260809T100113Z_5f31bde03200c464 ahead_commits=1
- branch=roundfix/run-run_20260809T101832Z_7e534c461190cb3f ahead_commits=2
- branch=roundfix/run-run_20260809T104632Z_583f49f98968df54 ahead_commits=2
```

Every one of them held only re-executions of Task 01 and Task 02 whose work was
already integrated. The suggested `git merge --ff-only` could not apply, because
the branches had diverged behind the integrated work. Clearing them meant
`git worktree remove --force` and `git branch -D` by hand, which is exactly the
operation the Daemon exists to keep a maintainer out of.

The cost compounds with the reason for stopping. A Run is stopped precisely when
something upstream needs fixing — an authored gate, a bounded scope, a
maintainer decision. That is exactly the moment when the completed Tasks are
most clearly still valid: nothing about them changed, only the Task that failed.

## Shape

Non-binding. The question worth settling first is which facts a stopped Run may
carry forward and on what evidence — a Task whose Verification passed against an
input set that has not changed is a different case from one whose inputs moved,
and only the first is safely reusable.

Worth settling in the same work: whether carrying a settled Task forward should
be automatic or an explicit maintainer action, since the failure mode of getting
it wrong is a Task reported complete against inputs it never saw. And whether
`reconcile` is the right home for it, given that its current contract is
deliberately read-only — if it stays read-only, something else has to own
discarding a superseded Run Branch, because today nothing does and the workflow
stays shut until a human reaches for `git branch -D`.
