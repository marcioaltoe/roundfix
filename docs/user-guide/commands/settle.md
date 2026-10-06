### settle

```bash
roundfix settle --spec <slug> --task <task_id>
```

Local recovery for one Task whose completed work already exists in a kept Task
Worktree, kept Run Worktree, or the current repository. Surface resolution
picks the first candidate, in that order, **where the task file is actually
`failed`, or `completed` while that surface still holds the work uncommitted**
— the state a refusing commit hook leaves behind, where the Daemon settled the
Task after its Verification passed and the commit never landed. A `completed`
Task whose surfaces are all clean was committed already and settles nothing.
Settle always reports the choice on stderr as `Settle surface: <path>`; when no
candidate qualifies, the refusal names each candidate path, the status found
there, and `no uncommitted work` for a surface that holds none. It re-runs the
Task's Verification commands in the selected tree — verbatim, in task file
order, with no Agent session — changes nothing on failure, and on pass prints
one sorted `commit <path>` line per committed path before the settled line. A
path the commit removes — a file the Task deleted, or the old name of one it
renamed — carries a `— deleted` qualifier, because the path alone no longer
exists in the surface to be read:

```text
verify go test ./... — ok
commit docs/specs/<slug>/task_01.md
commit internal/example/example.go
commit internal/example/obsolete.go — deleted
settled task_01 completed — <short sha>
```

The first command that does not pass ends the re-run and exits `1` with the
Task status, the surface, and the Run left as they were:

```text
verify go test ./... — failed (diagnostics: <artifact-dir>/runs/settle-<slug>-<task_id>/verification/batch-001-attempt-1.log)
task_01 stays completed — verification failed
```

A command whose verdict the runner never observed — it could not start, or its
diagnostics could not be retained — reports `verify <command> — verdict
unknown` and `<task_id> stays <status> — verification verdict unknown` instead,
with the cause on stderr. Settle has no Verification repair: fix the cause and
run it again.

The Task commit that `settle` creates applies the same disclosure rule as an
Implement Run: for a completed non-QA Task, it writes undeclared ordinary
commit paths under `## Recorded paths` and removes a stale section when the
commit has none. Recording discloses a path and reserves nothing; it never
records a Governed Path.

For a QA Task, settle applies the same report eligibility decision as archive
after Verification passes. A `pending` verdict is never accepted, and a report
that records no QA row is refused even when it says `pass` or would otherwise
qualify as declared-only `partial`.

If the newest report in the selected surface is ineligible, settle exits `1`
and leaves the Task file unchanged. Its stderr refusal names that report with
`(report <path>)`, relative to the surface's working tree:

```text
Settle surface: <path>
roundfix: settle QA Report is ineligible: <reason> (report <path>)
```

When other Tasks of the Spec are `failed` at settle time, one stderr warning
names them: their work may be swept into this commit. Settle creates no Run,
writes no Run Event Journal entries, and never pushes; Task Worktree
settlements integrate onto the Run Branch before the Run-level integration.

