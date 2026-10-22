---
task: task_01
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: pending
type: backend
complexity: medium
---

# Task 01: The Spec Consistency Check reports a test selector its own invocation runs whole

## Overview

No structural Verification check reads test runner arguments, so a pytest node
id passed beside its own file reaches the prober, which can only call the
command vacuous (Quaestor, 2026-10-21; the adopted Backlog Entry
`references/2026-10-22-a-pytest-node-id-beside-its-own-file-passes-before-the-test-exists.md`).
This Task adds `SC-VERIFY-SUBSUMED-SELECTOR`, an error from the Tasks stage
that reports each node id another collection argument of the same pytest
invocation contains, under pytest 9's own rule (ADR-0315), and tests it
through `CheckStage` (`_prd.md` → Core Feature 1).

## Requirements

1. MUST create `internal/speccheck/verification_selector.go` with the exported
   constant `CodeVerifySubsumedSelector = "SC-VERIFY-SUBSUMED-SELECTOR"` and
   an unexported `detectSubsumedSelectors(result *Result, repoRoot, specDir
   string, graph *spec.Graph) error` that implements `_techspec.md` → API
   Contracts 1-4 exactly: which Tasks it judges, how it finds a pytest
   invocation (command position after `!`, `rtk`, `env`, `exec` and
   `NAME=value`; `pytest`, `py.test`, `python[0-9.]* -m pytest`; `uv`,
   `poetry`, `pdm`, `pipenv` or `hatch` `run` followed later by a runner; the
   command and each `$(...)` body, split at `;`, `&&`, `||`, `|`, `&`, `(`,
   `)`, `{`, `}` with the existing `shellWords` and `commandSubstitutions`),
   which words are collection arguments (the option list of API Contract 2,
   redirections, globs without `::`), the `--keep-duplicates` exemption, the
   subsumption rule with its `/` boundary and parametrized last part, and the
   exact summary, fix and location of the finding. A nil graph records a skip
   with the manifest path, as `detectVerificationPaths` does. It MUST execute
   no command and read no file other than the Task file.
2. MUST call `detectSubsumedSelectors` from `internal/speccheck/constraints.go`
   right after the `detectVerificationPaths` call, and add `{code:
   CodeVerifySubsumedSelector, stage: StageTasks}` to `stagedDetectors` in
   `internal/speccheck/coherence.go` right after the
   `CodeVerifyPathUndeclared` entry. No other line of either file changes.
3. MUST create `internal/speccheck/verification_selector_test.go` (package
   `speccheck_test`, every top-level test calling `t.Parallel()` first, no
   real-time duration) with the three tests `_techspec.md` → Testing Approach
   names and exactly the commands it lists:
   - `TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath`, ten
     commands, each with exactly one error at `StageAll` and `StageTasks`,
     the exact summary and fix of API Contract 4 for its selector and
     containing argument, and the Verification line of the command;
   - `TestSubsumedSelectorAcceptsSelectorsNoArgumentContains`, ten commands
     with no error, among them `pytest tests/test_ab.py::test_new
     tests/test_a`, `pytest --keep-duplicates tests/test_a.py::test_new
     tests/test_a.py`, `pytest --deselect tests/test_a.py::test_old
     tests/test_a.py::test_old[new]`, `pytest -k tests/test_a.py
     tests/test_a.py::test_new` and `pytest tests/test_a.py::test_new; pytest
     tests/test_a.py`;
   - `TestSubsumedSelectorJudgesOnlyPendingNonQATasks`, completed and
     in-progress Tasks with no error, and the PRD stage with no error and the
     code skipped.
   It uses the existing helpers `writeConsumerFixture`, `checkProofFixture`,
   `findingsWithCode` and `hasExactLocation`.
4. MUST NOT change the prober (`internal/daemon/verification_probe.go`),
   `internal/cli/spec_check.go`, any other detector, any existing test or
   golden, or `internal/docscontract/testdata/corpus-golden.json`.

## Subtasks

- [ ] Write the detector and wire it into `Check` and the stage table.
- [ ] Write the three tests through `CheckStage`.

## Acceptance Criteria

- [ ] Each overlapping command yields one error with the exact text and line; each non-overlapping command, a non-pending Task and the PRD stage yield none.
- [ ] The prober, the CLI and every other detector are unchanged.

## Consumers

- `CodeVerifySubsumedSelector` and the finding it raises — consumers: `internal/speccheck`

## Context

- creates: `internal/speccheck/verification_selector.go`
- creates: `internal/speccheck/verification_selector_test.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- instruction: `internal/speccheck/verification_path.go`
- instruction: `internal/speccheck/verification.go`
- instruction: `internal/speccheck/verification_path_test.go`
- instruction: `internal/speccheck/consumers_test.go`
- instruction: `internal/speccheck/authoring_proof_test.go`
- instruction: `docs/adr/0315-a-test-selector-another-argument-of-its-invocation-contains-is-refused-at-authoring.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath|TestSubsumedSelectorAcceptsSelectorsNoArgumentContains|TestSubsumedSelectorJudgesOnlyPendingNonQATasks|TestVerificationPathUndeclaredRaisesTheMissingFile|TestVerificationPathUndeclaredIgnoresWhatThePlanProvides)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath TestSubsumedSelectorAcceptsSelectorsNoArgumentContains TestSubsumedSelectorJudgesOnlyPendingNonQATasks TestVerificationPathUndeclaredRaisesTheMissingFile TestVerificationPathUndeclaredIgnoresWhatThePlanProvides; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/speccheck` — expected: exit 0. Before this Task the three selector tests and `CodeVerifySubsumedSelector` do not exist, so the command fails at its first pass check. Afterwards the selector tests pass beside the unchanged Verification path tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Core Feature 1; Goals; Success Metrics 1-2
- `_techspec.md` → API Contracts 1-4; Testing Approach; Build Order 1
- ADR-0315
