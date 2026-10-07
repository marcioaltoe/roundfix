## Reopen Command

Use `roundfix reopen --spec <slug>` to clear a settled QA gate when one or
more of its dependencies are no longer completed, or when the gate has a Late
Dependency. Reopen proves a Late Dependency from the Task Graph manifest at
the commit that added the newest QA Report: a Task in the current dependency
closure that the recorded closure lacked reopens the completed gate to
`pending`, whatever its status. The invalidation record says `Dependencies added after the QA Report`.
This is the supported way to reopen the gate;
never edit the QA Task file by hand. The command retains the prior QA Report
and the Task's prior Result. Without the commit that added the newest report,
reopen refuses as before.

Reopen refuses before mutation when the terminal QA gate is not settled, or is
neither stale nor above a proven Late Dependency: it is not completed, or
every dependency is still completed and no Late Dependency is proven, either
because the closure is unchanged or because the newest QA Report has no
commit from which the prior Task Graph manifest can be proven. It creates no
Run, writes no Run Event Journal entry, and never commits or pushes.

Immediately before it writes, reopen rechecks the gate's staleness and refuses
with exit 2 when the gate changed since preflight.

Flags:

- `--spec` — Spec slug under the configured Spec Root.

Exit codes: `0` means a settled QA gate that was stale or above a Late
Dependency was reopened, `1` means the reopen write failed, and `2` means
Preflight Validation failed, including an unsettled QA gate or one that is
neither stale nor above a proven Late Dependency.

## Settle Command

Use `roundfix settle --spec <slug> --task <task_id>` only as a local recovery
command for one Task whose completed work is already in a kept Task Worktree, a
kept Run Worktree, or the current repository. Settle resolves that surface by
loading the target Task status in order from the deterministic Task Worktree
path, the Run Worktree recorded on the latest kept Run, and the current
repository. It selects the first candidate whose task file is `failed`, or
`completed` while that surface still holds the work uncommitted — the state a
refusing commit hook leaves behind, where the Daemon settled the Task after its
Verification passed and the commit never landed.

Flags:

- `--spec` — Spec slug under the resolved Spec Root (`docs/specs/` by
  default).
- `--task` — Task id from the Spec Task Graph.

Preflight Validation exits `2` with one actionable message when either flag is
missing, the repository does not resolve, the Spec or Task Graph does not load,
the Task id is absent from the Task Graph, no candidate surface holds the target
Task in a settleable state, a settle surface path exists but is unusable, or
another Active Run owns the Spec target or working tree. When no surface
qualifies, the refusal names every candidate path and the status found there,
`no uncommitted work` for a `completed` surface that holds none, or that the
path does not exist. `pending` and `in_progress` Tasks belong to the Implement
Command; a `completed` Task whose surfaces are all clean was committed already
and settles nothing.

On every settle that proceeds, stderr prints the selected surface before
Verification starts:

```text
Settle surface: <path>
```

stdout carries only deterministic report lines:

```text
verify test -f done.txt — ok
commit <path>
settled task_01 completed — <short sha>
```

On pass, settle prints one sorted `commit <path>` line for each path included
in the commit, between the verification lines and the settled line. When
nothing is stageable and settle creates no commit, it prints no `commit <path>`
lines.

If verification fails, the command stops at the first failed Verification
command, leaves the Task and tree unchanged, and prints:

```text
verify test -f done.txt — ok
verify test -f missing.txt — failed (diagnostics: <path>)
task_01 stays failed — verification failed
```

If a Task's Verification is unsatisfiable because the task file names a
non-hermetic or impossible command, fix that task file's `## Verification`
section and re-run Settle. Settle re-reads the task file on each invocation.
There is no skip-verification flag: Verification is the only gate before
settling, committing, or integrating Task work.

Exit codes: `0` means settled completed and committed, `1` means verification
failed or post-verification integration failed, and `2` means Preflight
Validation failed.

On pass, settle verifies in the selected surface, stages that surface's changes
plus the task file, creates the standard Task commit, creates no Run, writes no
Run Event Journal entries, and never pushes. If other Tasks in the same Spec
are failed at settle time and a commit is created, stderr prints one warning:

```text
roundfix: warning: other failed Tasks in Spec "<slug>" may have work included in this settle commit: task_02, task_03
```

For a completed non-QA Task, that Task commit writes undeclared ordinary paths
under the Task file's `## Recorded paths` section. Recording discloses a change
and reserves nothing; a Governed Path remains subject to authorization.

When the selected surface is a Task Worktree, settle integrates that commit
onto the Run Branch through the same queue mechanics as `implement`; success
removes the Task Worktree and Task Branch. A Task Worktree integration conflict
exits `1`, keeps both the Run and Task worktrees and branches, leaves stdout
with only verification lines, and prints stderr shaped like:

```text
roundfix: settle failed after verification: task worktree integration conflict on <path>
```

After a successful Task Worktree integration, or when settling from the Run
Worktree, settle then runs the Run-level integration protocol. A Run-level
integration refusal exits `1` and prints:

```text
integration pending — git merge --ff-only roundfix/run-<id>
```

Review the Run Worktree before running it.

## QA Report acceptance

```bash
roundfix qa-report accept <path>
```

Reads the selected QA Report and its Spec directory's declared-acceptance
declarations, applying the shared archive and settlement eligibility policy.
It writes no files and produces no stdout; exit status is the machine-readable
result. Exit `0` means acceptable, `1` means missing, unreadable, or
unacceptable, and `2` means a usage error. Read its usage through
`roundfix qa-report --help`; `accept` has no `--help`.

A `pending` verdict, a report with no QA row, and empty or duplicated front
matter are refused. The pre-PR Pull Request row recorded as
`blocked (environment: no open Pull Request)` with that row named in its
provenance, and an outside-evidence row recorded as
`blocked (environment: network denied: <host>)` with the outside-evidence row
named in its provenance, never decide a qualifying partial and need no
Unreachable Acceptance declaration. Any other blocked outside-evidence row
still blocks Pull Request preparation. When this eligibility check refuses,
`qa-report accept` prints the reason alone. Only the Settle Command appends
`(report <path>)` to its refusal, naming the report it judged.
