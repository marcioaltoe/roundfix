---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Archive and cancel refuse ignored Spec files and finish a failed removal forward

`roundfix archive` and `roundfix cancel` remove the Spec folder and promise
that its bytes stay in Git at the record's `source_revision`. Their only guard
was `git status --porcelain --untracked-files=all -- <folder>`, which omits
files Git ignores. On 2026-10-12, in a temporary repository whose
`.gitignore` held `*.log`, both commands removed `run.log` from the Spec
folder, counted it among the files "kept in Git", and no commit held it. When
the final `os.RemoveAll` fails, both commands have already written the record,
the promoted or returned sources and the relinked documents. Archive then
reported "Preflight failed" and "No side effects" with exit 2, cancel reported
"cancel write failed", and neither said how to finish. A second run is refused
because the record already exists.

On 2026-10-12 the maintainer decided:

- **Refuse ignored files.** Both commands read the Spec folder with
  `--ignored`, refuse while Git ignores any file under it, list every such
  file one by one, and write nothing. The operator removes each file, or
  commits it with `git add -f`, and repeats. The existing refusal of
  uncommitted changes keeps its text and is reported first.
- **Recover forward.** A failed final folder removal keeps everything already
  written; nothing is rolled back. The command exits 1 and prints, on standard
  error, what it wrote and the exact command that finishes the removal:
  `chmod -R u+w '<folder>' && rm -rf '<folder>'`, with the absolute folder
  quoted for a POSIX shell. Archive still records `--pull-request` in the
  written record. After the earlier refusals, the folder holds only files
  committed at `source_revision`, so the command removes nothing that Git does
  not hold.
- **Symbolic-link destination parents stay accepted.** The external Specs
  Root of a knowledge workspace is one; this decision does not change that
  check.

This completes ADR-0247 ("so the recorded revision holds exactly what was
removed") for files Git ignores, and applies the same guard to the Cancel
Command of ADR-0262.

## Considered Options

- Copying ignored files into the record or into `docs/references/` would put
  scratch output into history, which ADR-0247 keeps small.
- Committing ignored files automatically would override the repository's own
  ignore rules and make the command commit, which neither command does.
- Warning and removing anyway keeps the defect: the bytes are lost.
- Rolling back after a failed removal would also have to restore files
  `os.RemoveAll` already deleted. It "removes everything it can but returns
  the first error it encounters" (`go doc os.RemoveAll`), so a rollback would
  rewrite deleted Spec files from Git, a second failure mode for one failure.
- Printing only `rm -rf '<folder>'` fails again on the usual cause, a
  directory without write permission, so the command restores write
  permission first.

## Consequences

An ignored file such as `.DS_Store` or a log under `qa/evidence/` now blocks
both commands until the operator removes or commits it, including inside a
Delivery Queue item worktree, where the archive step parks the item as any
failed archive does. Archive's exit code for a failed folder removal moves
from 2 to 1, so a script that read 2 as "nothing happened" now sees the truth.

**Narrowed by [ADR-0306](0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md) (2026-10-20).** A symbolic-link destination parent whose resolved target stays in the repository, or in the external Spec Root's repository, stays accepted; one whose target leaves it is refused before anything is written.
