---
task: task_01
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: completed
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

## Result

Implementation:

- `internal/speccheck/verification_selector.go` (new): `CodeVerifySubsumedSelector` and `detectSubsumedSelectors`. It judges pending Tasks that are neither the manifest's QA Task nor `type: qa`, reads each Verification command and every `$(...)` body (`commandSubstitutions`), splits them with `shellWords` at `;`, `&&`, `||`, `|`, `&`, `(`, `)`, `{`, `}`, and recognizes `pytest`/`py.test` by base name, `python[0-9.]* -m pytest`, and `uv|poetry|pdm|pipenv|hatch … run … <runner>`, after `!`, `rtk`, `env`, `exec` and `NAME=value`. Collection arguments exclude options, the values of the API Contract 2 option list, redirection targets, and globs/`$`/backticks without `::`. `--keep-duplicates` skips the invocation. Subsumption follows API Contract 3: `.`, an equal path, or a directory above it at a `/` boundary; or an equal path whose parts prefix the selector's, with the last part matching exactly or as the name before `[`. Each finding uses the exact summary and fix from API Contract 4, located at the Verification line via `sectionLineContaining`. A nil graph records a skip naming `_tasks.md`. Nothing is executed, and the Task file is read only when a finding needs its line.
- `internal/speccheck/constraints.go`: one call right after `detectVerificationPaths`.
- `internal/speccheck/coherence.go`: `{code: CodeVerifySubsumedSelector, stage: StageTasks}` right after `CodeVerifyPathUndeclared`.
- `internal/speccheck/verification_selector_test.go` (new): the three tests with the commands listed in the TechSpec's Testing Approach, run through `CheckStage` using `writeConsumerFixture`, `checkProofFixture`, `findingsWithCode`, `hasExactLocation` (and `hasSkip` for the PRD-stage skip).

Acceptance criteria evidence:

- Overlapping commands raise one error each, and non-overlapping commands, non-pending Tasks and the PRD stage raise none: `go test -count=1 -run 'TestSubsumedSelector|TestVerificationPath' ./internal/speccheck` → `ok`. Sabotage runs (reverted afterwards) prove the tests can fail:
  - Dropping the parametrized match fails the `test_it[new]` case.
  - Dropping the `.` rule fails the `python3.12 -m pytest .` case.
  - Removing the `/` boundary fails `tests/test_ab.py::test_new tests/test_a`.
  - Ignoring `--keep-duplicates` fails that accept case.
  - Judging every status fails the `completed` and `in_progress` subtests.
- The prober, the CLI and every other detector are unchanged. `git diff --stat` lists only `coherence.go` (+1) and `constraints.go` (+3) plus the two new files. `go test -count=1 ./internal/speccheck` → `ok` (40.8s), `go test -count=1 ./internal/cli/...` → `ok`, `GOOS=windows go vet ./internal/speccheck` → clean, and `gofmt -l` → clean.

## Carry-forward provenance

- Source Run: `run_20261023T070317Z_62af7ae2274ed43c`
- Source commit: `e1f61677b58086d92e2facb27d848d2bd91deb16`
