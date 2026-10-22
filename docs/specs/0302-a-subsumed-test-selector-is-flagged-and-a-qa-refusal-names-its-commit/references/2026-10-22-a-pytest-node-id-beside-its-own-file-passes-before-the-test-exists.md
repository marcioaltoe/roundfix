---
type: fix
status: promoted
created: 2026-10-22
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
reason: null
---

# A pytest node id beside its own file passes before the test exists, and the check calls it vacuous without saying why

## Symptom

Quaestor reported on 2026-10-21 at 19:02 BRT (an ai-memory message to
Roundfix, transcribed by the operator's cycle AP briefing of 2026-10-22) that
pytest ignores a node id that does not exist yet, such as
`tests/a.py::test_novo`, when the same command also passes the whole file,
`tests/a.py`. The command exits 0, and `roundfix spec check
--run-verification` classifies it `vacuous` without naming the cause. With
only the node ids pytest exits 4 and the command is honest. The author sees a
vacuous verdict on a command that names a test which does not exist, and has
to find out by hand that the file argument swallowed the node id.

## Where

The authored Verification of a Task in an adopting Python repository, checked
by `roundfix spec check <slug> --strict --run-verification`
(`internal/cli/spec_check.go` `probeSpecVerifications`, which runs the shared
prober `daemon.ProbeCommands`), and the structural Verification checks in
`internal/speccheck` (`verification.go`, `verification_path.go`), none of
which reads test runner arguments. The `write-tasks` skill says nothing about
test selectors.

## Expected

`roundfix spec check` reports, from the written command alone and before any
Run, a Verification that passes a test selector together with a collection
argument of the same invocation that contains it, names the selector and the
argument, and says how to fix it; the `write-tasks` skill states the rule.
The prober's `vacuous` verdict is unchanged.

## Evidence

Measured on 2026-10-22 in a throwaway virtual environment, with a
`tests/test_a.py` that holds only `test_old`:

- pytest 9.1.1: `pytest tests/test_a.py::test_novo tests/test_a.py`,
  `pytest tests/test_a.py tests/test_a.py::test_novo`,
  `pytest tests tests/test_a.py::test_novo`,
  `pytest . tests/test_a.py::test_novo`,
  `pytest tests/test_a.py::test_novo ./tests/test_a.py` and
  `python -m pytest tests/test_a.py::test_novo tests/test_a.py` each run
  `test_old` and exit 0. `pytest tests/test_a.py::test_novo` alone,
  `pytest tests/test_a.py::test_novo tests/test_b.py`,
  `pytest tests/test_a.py::test_old tests/test_a.py::test_novo` and
  `pytest --keep-duplicates tests/test_a.py::test_novo tests/test_a.py`
  exit 4 with `ERROR: not found: … tests/test_a.py::test_novo`.
- pytest 8.3.5 and 7.4.4: every one of those commands exits 4.

pytest 9.0.0's changelog (#12083, https://docs.pytest.org/en/9.0.x/changelog.html)
states the cause: "`pytest a/b a/` or `pytest a/ a/b` are equivalent to
`pytest a`; if an argument overlaps another arguments, only the prefix
remains", and "If you rely on these behaviors, consider using
`--keep-duplicates`, which retains its existing behavior". pytest's
`is_collection_argument_subsumed_by` (`_pytest/main.py`) applies it to node
ids: a path argument with no parts subsumes every argument at or below it,
and a node id subsumes the node ids it prefixes.
