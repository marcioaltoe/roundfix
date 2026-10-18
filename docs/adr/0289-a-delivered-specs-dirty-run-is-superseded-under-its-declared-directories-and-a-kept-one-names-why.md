---
status: accepted
created_at: 2026-10-18T17:30:00Z
updated_at: 2026-10-18T17:30:00Z
deprecated_at: null
superseded_by: null
---

# A delivered Spec's dirty Run is superseded under its declared directories, and a kept one names why

After Spec 0272 merged as pull request #555, `roundfix reconcile` released three
of its seven Runs as `superseded` ("N uncommitted path(s) superseded by the
archived Spec") and kept three as `dirty`, each with only "Run Worktree has
tracked or untracked changes": `run_20261017T022748Z`, `run_20261017T032927Z`
and `run_20261017T040614Z`. The maintainer had them removed by hand on
2026-10-18.

The rule of ADR-0212 releases a dirty Run of a merged Spec only when two proofs
hold (`internal/worktree/worktree.go:1053`): the commit proof at the
merged head or by the delivery is `safe` or `superseded`, and every uncommitted
path lies under the Spec's directories or equals a path one of the archived
Spec's Tasks declares under `interface:`, `creates:` or `deletes:`, or records
under `## Recorded paths` (`internal/worktree/merged_head.go:785-859`). The
comparison is exact equality (`merged_head.go:853-857`), so a declaration that
names a directory covers no file beneath it. 0272's task_04 declares
`deletes: docs/history/findings` and `deletes: docs/history/backlog`. When
either proof fails, the reason is the same constant
(`worktree.go:1062-1063`), so the operator cannot tell which proof failed or
which path kept the Run.

The three kept Runs settled no Task the delivery did not also settle: the first
settled task_01 and failed task_02, the second failed task_02, and the third settled task_02 and task_03 and failed
task_04. They ran under plan revisions `123f6524`, `adc564d0` and `6ee8d456`,
before pull request #551 added task_06. The three released Runs ran after it.
Their worktrees were removed by hand, and their journals record Task phases,
not paths, so the exact path that kept each one is not recoverable.

## Decision

**A declared path covers the paths beneath it.** A dirty path of a Run of a
Spec whose Archive Record is at the proof head is in scope when it lies under
the Spec's directories, equals a path a Task of the archived Spec declares
under `interface:`, `creates:` or `deletes:` or records under
`## Recorded paths`, or lies beneath such a path (it begins with the path and a
`/`). The commit proof is unchanged: every commit the Run holds and the
delivered head lacks is represented there, or carried into it as a completed
Task, or superseded by the delivery evidence of ADR-0212.

**A kept dirty Run names why.** When the Spec is archived at the proof head and
the Run stays `dirty`, the reason is the constant followed by `; ` and the
first failing proof: the commit proof's own reason, or `<path> is outside the
archived Spec's declared and recorded paths` for the first such path in byte
order, bounded to the existing 160 bytes. A Run of a Spec that is not archived
keeps the bare constant.

## Considered Options

- Releasing every dirty Run of a merged Spec was rejected. A stray edit outside
  the Spec's scope is work only its author can judge, and ADR-0053 keeps dirty
  work unless a proof covers it.
- Reading the declarations of the plan revision each Run started from was
  rejected. Every path any of 0272's seven plan revisions declared is still
  declared at its source revision, so it would have released none of the three
  Runs, and it trusts a Task file the Run itself could have edited.
- Comparing each dirty path's bytes with the delivered head was rejected. It
  hashes worktree files on every reconcile and covers almost nothing the
  declarations miss.
- Requiring a declared directory to exist at the source revision was rejected.
  A `deletes:` declaration names a directory the delivery removed.

## Consequences

This refines ADR-0212 and keeps ADR-0053's safety bar: no change outside what
the delivered Spec declared, recorded or holds in its own directories is
released. Releasing stays explicit, through `roundfix reconcile --apply` or the
Delivery Queue's release after a merge, and apply re-proves the exact dirty
path set before it removes anything. The Delivery Queue's cleanup warning and
`roundfix reconcile` print the named reason, so the next kept Run says which
path or commit keeps it.
