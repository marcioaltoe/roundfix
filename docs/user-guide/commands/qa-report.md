### qa-report accept

```bash
roundfix qa-report accept <path>
```

Reads the selected QA Report and exits zero only when the shared archive and
settlement eligibility decision accepts it. The pre-PR Pull Request row,
recorded as `blocked (environment: no open Pull Request)` with the Pull Request
row named in its provenance, never decides a qualifying partial and needs no
Unreachable Acceptance declaration. An outside-evidence row blocked only because
the Run sandbox denied network access, recorded as `blocked (environment: network denied: <host>)` with the outside-evidence row named in its provenance,
also never decides a qualifying partial and needs no Unreachable Acceptance
declaration. So does an outside-evidence row the gate did not look up because
the Spec's own authorization forbids that lookup, recorded as `blocked
(environment: <subject> prohibited by Spec authorization)` with the
outside-evidence row named in its provenance, where the subject names what the
authorization forbids in 1 to 80 characters without parentheses (ADR-0276,
ADR-0335). Any other blocked outside-evidence row still blocks Pull Request
preparation. A `pending` verdict is never
accepted, and a `pass` or otherwise-eligible `partial` that records no QA row
is refused. A report whose front matter is empty or duplicated is unreadable
and refused. The command writes no files.

One Unreachable Acceptance declaration may name several criteria, one criterion line each, and each criterion becomes its own QA row; acceptance matches every declared-blocked row to a declaration that names its criterion, by its text or its leading label.
A report whose declared-blocked rows do not exceed its declarations is accepted as before.

Every local link of the report must resolve in the committed tree. A local
link is an inline Markdown link or reference definition outside code whose
destination is a relative path, resolved against the report's directory;
URLs, anchors and absolute paths are excluded. An existing target resolves
when Git tracks it or would add it, including a directory containing such a
file. A tracked target resolves even when it matches an ignore pattern.
Outside Git, an existing target resolves.

An Unresolved Evidence Link is `ignored` when the target exists but Git ignores
it and does not track it, or `missing` when it is absent or outside the
repository. Keep logs and command output under an unignored name such as
`.txt`, then update the link. The command exits `1` with no stdout and this
stderr line, joining multiple links with `; `:

```text
roundfix: QA Report eligibility refused: QA Report evidence links do not resolve in the committed tree: <destination> (line <n>, <reason>)
```

Every evidence file the report relies on must belong to this Run. Borrowed
Evidence is `reused` when a linked file under the Spec's `qa/evidence/` has
the same bytes at that path in the audited head, or `imported` when it
matches the QA Report Import's copy; it `names another Run Worktree <path>`
when the report or linked evidence names a registered worktree of another
Run, or is `outside the worktree` when a `user_flow_executables` entry lies
outside this worktree, compared without resolving the executable's symbolic
links. The optional key must be a list of strings; otherwise its reason is
`not a list of paths`. The audited head is `HEAD` for an uncommitted report,
otherwise the first parent of the commit that added it; an unborn head or
root addition has no audited head for the `reused` check. A Carried Row's
Results line and `### <row-id>` block are exempt. Outside Git nothing is
checked. A refusal exits `1` with no stdout and this stderr line, joining
multiple items with `; `:

```text
roundfix: QA Report eligibility refused: QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>)
```
