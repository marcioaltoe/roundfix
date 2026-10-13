### cancel

```bash
roundfix cancel <slug> --reason <text>
```

Cancels an active Spec that will not be delivered. The command writes one
Archive Record with disposition `cancelled`, the reason, the date and the
Tasks completed and pending. It returns each adopted Backlog Entry to
`docs/backlog/` with `status: open` and each adopted Finding to
`docs/findings/` with `status: pending`, then removes the committed Spec
folder. The removed bytes stay in Git at the record's `source_revision`.

`cancel` creates no Run, never changes a Task status, and never commits or
pushes. A cancelled Spec never meets a prerequisite or supersedes another
Spec. `deliver retry` refuses a cancelled Spec.

The command refuses with exit `2` and leaves the files unchanged when the slug
or reason is missing, an unexpected argument or unknown active Spec is given,
the working tree is not a Git repository, a non-Markdown file names the active
directory, the Spec folder has uncommitted changes, Git ignores a file under
the Spec folder, an Active Run or an unfinished Delivery Queue item holds the
Spec, an active Spec requires it,
or the Spec carries a supersession. Every refusal of the cancellation itself
also leaves the files unchanged. A write or removal failure after validation
exits `1`.

When Git ignores a file under the Spec folder, the refusal lists every ignored
file by its sorted, repository-relative path, including files inside ignored
directories. To continue, remove each file or commit it with git add -f, then
retry. Uncommitted changes keep their existing refusal and are reported first.
For an external Spec Root, the command reads that root's repository and its
ignore rules.

If the final Spec folder removal fails, the command keeps the Archive Record
and every other write, prints nothing on standard output, and exits `1`.
Standard error names the record, the remaining Spec folder and the cause, then
prints the exact command that finishes the removal:

```text
roundfix: cancel wrote <record> but could not remove the Spec folder <folder>: <cause>
roundfix: finish the removal with: chmod -R u+w '<absolute-folder>' && rm -rf '<absolute-folder>'
```

Run the printed command to restore write permission and finish removing the
committed folder. Its absolute path is quoted for a POSIX shell, including
spaces and embedded single quotes. Returned sources, promotions and relinked
documents stay written. Retrying `roundfix cancel` refuses because the record
already exists.

After success, a parked Delivery Queue item produces one warning on standard
error because Delivery Retry refuses a cancelled Spec. The command also warns
once per tracked Markdown line outside `docs/history/` that cites a removed
path. These warnings do not change the exit code.

Options:

- `--reason <text>` — Why the Spec is cancelled; the reason is recorded in the
  Archive Record.

Exit codes:

- `0` — cancellation recorded.
- `1` — cancellation write or removal failed.
- `2` — Preflight Validation failed.
