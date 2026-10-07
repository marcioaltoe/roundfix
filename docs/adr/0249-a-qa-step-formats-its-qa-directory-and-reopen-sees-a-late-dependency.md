---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A QA step formats its QA directory, and reopen sees a Late Dependency

On 2026-10-06 an adopter whose repository Verification runs its formatter in
check mode lost a QA Run to a file Roundfix wrote. A failed QA pass had
committed its report without the repository's formatter. ADR-0194 imports that
unintegrated pass into the next Run's worktree byte for byte, and the next
gate's repository Verification precondition refused on the imported report's
formatting before any row ran. Every passing QA in that repository also needed
a separate formatting commit, because the QA Report commit is the one commit
the Daemon makes after the last repository Verification. The QA Agent could
not repair it either: its sandbox denied the formatter's package-runner cache.

Three repairs were weighed. Excluding the imported files from the precondition
is not possible in general, because the precondition is the repository's
opaque command, and the files ride in the QA Report commit anyway, so the next
gate or CI fails on them instead. Not importing a failed pass reverses
ADR-0194 and still leaves every committed report unformatted. Formatting what
the QA step commits repairs both, for every stack, provided Roundfix does not
guess the formatter.

Roundfix therefore reads one optional Project Config command,
`verification.format`, the Format Command. The Daemon runs it outside the
Agent sandbox, from the Run Worktree root, with each file to format appended as
its own argument, as pre-commit tools such as lint-staged do. Which file types
it accepts is the command's concern. It runs at two points, over regular files
under the Spec's `qa/` directory only. The first point is after the mechanical
stage and before the repository Verification precondition, over the imported
pass. The second is immediately before the QA Report commit, over the files that
commit stages. An empty value runs nothing, and that is today's behavior. A
failing, timed-out or verdict-changing run restores every file's original
bytes, and the commit proceeds unformatted. A Run Event records each run. A
formatter never costs a QA verdict.

Roundfix does not infer the command from the Baseline profile, from
`make fmt`, or from the repository Verification. A repository need not adopt a
profile, and a whole-tree target such as `make fmt` takes no file list and may
rewrite files the QA step does not own. Task files and other Spec artifacts are
never formatted, because the Daemon must not change authored Spec text, and
Task commits keep their Task Verification as the only gate before them.

The same adopter added a corrective Task, already `completed` because the
supervising session implemented it, as a dependency of a QA gate that had
settled `completed`. `roundfix reopen` refused, because it reopens only a gate
above a dependency that is no longer completed, and the operator edited the QA
Task by hand. Reopen now also proves a Late Dependency from Git. It reads the
Task Graph manifest at the commit that added the newest QA Report. Any Task in
the current dependency closure of the QA Task that the recorded closure lacks
was added after the gate, whatever its status, and the gate returns to
`pending`. When no such commit exists, reopen refuses exactly as before:
the report may be uncommitted, outside the repository, or older than its
manifest. The Task Graph loader keeps its file-only staleness rule.

## Consequences

ADR-0194's byte-for-byte import is refined, not superseded: the import still
copies the committed bytes and proves carries on them, and only then are they
formatted. A pass that recorded its report through a configured Format Command
leaves the formatter nothing to change on import, so carries keep their proof.
A report recorded before the command was configured may lose its carries in
the following pass, which re-runs those rows.

An adopter whose formatter rejects files it does not support configures that
formatter's own option to ignore them, as with lint-staged. A refusal reverts
and is recorded; it does not fail the gate.

A Late Dependency needs the newest QA Report committed in the checkout's
history; without that commit reopen refuses as it did before. `roundfix
implement` and `roundfix deliver` do not detect a Late Dependency, so the
guides keep the rule that a Task added after the gate settled is followed by
`roundfix reopen`.
