# Corpus, scope and outside evidence

The outside source is the repository's own Task Graph corpus, authored by work
that predates and lies outside Spec 0131:

- active graphs: `docs/specs/*/_tasks.md`;
- archived graphs: `docs/history/specs/*/_tasks.md`.

The assembled public binary ran `roundfix spec check` across all active Specs
and exited 0. A second per-graph sweep through `roundfix spec check <slug>
--stage tasks` loaded both active Task Graphs and reported:

```text
active_task_graphs_loaded=2
```

For archived compatibility, each archived manifest and its real Task files
were copied to a disposable Spec root with only the PRD lifecycle changed to
`active`; no Task or graph byte was rewritten. The public checker loaded every
one before running its task-stage consistency checks:

```text
archived_task_graphs_loaded=102
```

An independent raw filesystem count also found 102 archived `_tasks.md` files.
This evidence comes from 104 repository Task Graphs that Spec 0131 did not
author and exercises the public checker rather than Task 02's test helper.

The committed change range `62bb7874..2d9ab9b2` contains only:

```text
docs/specs/0131-a-failed-gate-accepts-its-repair/task_01.md
docs/specs/0131-a-failed-gate-accepts-its-repair/task_02.md
internal/spec/spec.go
internal/spec/spec_test.go
```

The production diff is one condition:

```diff
- if gate.Status != StatusCompleted && gate.Status != StatusFailed {
+ if gate.Status != StatusCompleted {
```

Fresh Git comparisons found no changes to `internal/daemon`, `docs/adr`,
`CONTEXT.md`, `go.mod`, `go.sum`, `.agents`, or `Makefile`, and `git diff
--check 62bb7874..HEAD` passed. The CLI leaf-coverage and completed-gate probes
also prove gate terminality and dependency coverage remain enforced. No domain
term was added, changed or retired.
