---
type: fix
status: open
created: 2026-10-12
spec: null
reason: null
---

# Archive and cancel remove ignored files in the Spec folder

## Symptom

`roundfix archive` and `roundfix cancel` refuse uncommitted changes in the Spec
folder so that the record's `source_revision` holds every removed byte, but the
check omits files Git ignores. Such a file is removed with the folder and exists
in no commit. When the final folder removal fails, both commands print only
`leftover folder "<path>"`, without the command that finishes the removal.

## Where

- `internal/cli/archive.go:156` and `internal/cli/cancel.go:100` run
  `git status --porcelain --untracked-files=all -- <folder>` without `--ignored`.
- `internal/spec/archive.go:419` and `internal/spec/cancel.go` return the
  leftover-folder error after `os.RemoveAll` fails.

## Expected

Both commands refuse, list each ignored file under the Spec folder, and write
nothing; the operator removes or commits the files and repeats. A failed folder
removal keeps what was written (forward recovery, maintainer decision of
2026-10-12) and prints the exact command that finishes it. Symbolic-link
destination parents stay accepted (maintainer decision of 2026-10-12: the
external Specs Root of a knowledge workspace is one).

## Evidence

Pre-PR review of Spec 0261 on 2026-10-12, findings F1 and F3 at head
`7f22d85e`, dismissed there because the Archive Command shares the behavior.
