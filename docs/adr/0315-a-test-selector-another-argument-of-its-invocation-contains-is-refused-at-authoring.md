---
status: accepted
created_at: 2026-10-22T22:00:00Z
updated_at: 2026-10-22T22:00:00Z
deprecated_at: null
superseded_by: null
---

# A test selector that another argument of its invocation contains is refused at authoring

pytest 9.0 collapses collection arguments that overlap: "`pytest a/b a/` or
`pytest a/ a/b` are equivalent to `pytest a`; if an argument overlaps another
arguments, only the prefix remains" (pytest 9.0.0 changelog, #12083). A node
id that names a test the Task has not written yet, passed together with its
own file, a directory above it or `.`, is dropped without an error, and the
command passes on the tests the file already holds. Quaestor met this on
2026-10-21: `roundfix spec check --run-verification` called such a command
`vacuous` and could not say why. The prober is right that the command proves
nothing; what the author lacks is the cause, and the cause is visible in the
written command before anything runs.

## Decision

**The Spec Consistency Check reports the overlap from the written command.**
`SC-VERIFY-SUBSUMED-SELECTOR` is an error from the Tasks stage, with no
horizon, over the Verification commands of each pending non-QA Task. It reads
every pytest invocation in a command, at command position or inside a
`$(...)` substitution: `pytest`, `py.test` or `python -m pytest`, also behind
environment assignments, `rtk`, `env`, `exec`, or a `run` wrapper (`uv run`,
`poetry run`, `pdm run`, `pipenv run`, `hatch run`). Within one invocation it
reports each node id, a collection argument that carries `::`, that another
collection argument of the same invocation contains under pytest's own rule:
an argument without `::` that names the node id's file, a directory above it
or `.`, or a node id whose parts prefix it (`x.py::test_it` contains
`x.py::test_it[1]`). Paths are compared lexically after cleaning; option
values, redirection targets and globs are not collection arguments, and
`--keep-duplicates` exempts the invocation. The finding names the selector
and the argument that contains it, and its fix points to the `write-tasks`
rule "Test selectors and their containing path". The check executes no
command and reads no test file.

**The prober's verdict is unchanged.** The vacuous, honest, unknown and
malformed classification of ADR-0148 stays the one shared classification of
the Daemon's pre-work probe and `--run-verification`; it gains no cause text.
A command the structural check reports still runs, and still comes out
`vacuous` when the containing argument passes, so the author reads both: the
verdict says the command proves nothing, the finding says why.

## Considered Options

- **Explain the vacuous verdict.** The prober would parse each command's
  output for a runner that dropped a selector. That is semantic: it depends on
  the runner's output format and version, and it puts runner knowledge inside
  the one prober ADR-0148 keeps free of it. It would also only speak after the
  command ran, while the overlap is decidable from the text.
- **A rule for every test runner.** Only pytest was measured to drop a
  selector its own file contains; `go test -run` and the others select by
  name, not by overlapping paths. The code and the glossary name the class
  (**Subsumed Test Selector**) so another runner can join it when a measurement
  shows the same behavior.
- **A gap instead of an error.** The command cannot fail before the work, so
  the Daemon refuses the Task at its first Run; an error at authoring costs
  nothing that the Run would not cost later.

## Consequences

- An author who passes `tests/a.py::test_new tests/a.py` learns at `spec
  check` time to name only the selector, or to run the file in a separate
  command.
- A pytest older than 9.0 does not collapse the arguments and exits 4 on the
  missing node id; the rule still reports the command, because the plan does
  not know which pytest will run it and the selector adds nothing to an
  invocation that runs its file whole.
- `--keep-duplicates` keeps pytest's earlier behavior, so an invocation that
  passes it is not reported.
