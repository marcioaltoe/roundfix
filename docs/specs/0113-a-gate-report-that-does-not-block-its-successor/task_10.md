---
task: task_10
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 10: Give the gate's classification a boundary the code can see

## Overview

task_08 put the classifier in the Daemon's request builder, which contradicted
the contract that Implement runs no Spec check. task_09 took it out and put it
nowhere, because the requirement said "where the gate runs its strict check"
without naming a boundary any code could test. `GatePrecondition` still has no
production caller.

There is a boundary already in the surface. An authoring stage runs
`spec check <slug> --stage <stage>`; the gate runs the full sweep,
`spec check <slug> --strict`, with no stage. A Spec-owned declared term is an
error in the first and a repair input in the second — no skill file changes, no
precondition in the path that starts a Run.

## Requirements

1. MUST keep a Spec's declared-but-undocumented term an error when the check runs
   a single authoring stage, so authoring still reports it.
2. MUST classify that same finding as a repair input rather than an error on the
   full-sweep strict run, so the gate assigned the repair can reach it.
3. MUST call `GatePrecondition` from production to make that classification,
   giving it the caller it has never had.
4. MUST keep a term emitted by code that no Spec declared an error in both.
5. MUST NOT add a Spec-consistency precondition to the Implement path, and MUST
   NOT edit any skill file.

## Subtasks

- [ ] Distinguish the authoring-stage run from the full-sweep strict run.
- [ ] Route the full sweep through `GatePrecondition`.
- [ ] Cover both runs, the undeclared term, and the Implement contract.

## Acceptance Criteria

- [ ] `spec check <slug> --stage techspec --strict` still errors on a Spec's own
      declared-but-undocumented term.
- [ ] `spec check <slug> --strict` reports it as a repair input and does not exit
      non-zero for it alone.
- [ ] A term no Spec declared errors in both runs.
- [ ] `GatePrecondition` resolves from non-test production code.
- [ ] `TestRunImplementHasNoSpecCheckPrecondition` still passes unchanged.

## Verification

- `go test -count=1 ./internal/cli ./internal/speccheck -run 'TestSpecCheckClassifiesTheGateBoundary' -v > /tmp/0113-t10.log 2>&1; s=$?; grep -q '^--- PASS: TestSpecCheckClassifiesTheGateBoundary' /tmp/0113-t10.log || { cat /tmp/0113-t10.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t10.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestSpecCheckClassifiesTheGateBoundary' /tmp/0113-t10.log || { echo 'the boundary test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t10.log > /tmp/0113-t10-n.txt; test "$(cat /tmp/0113-t10-n.txt)" -ge 3 || { echo "expected the authoring-stage, full-sweep, and undeclared-term cases, got $(cat /tmp/0113-t10-n.txt)"; cat /tmp/0113-t10.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving the boundary is tested in both directions.
- `grep -rln 'GatePrecondition(' internal/ --include='*.go' | grep -v '_test.go' > /tmp/0113-t10-callers.txt; n=$(grep -c . /tmp/0113-t10-callers.txt); test "$n" -ge 2 || { echo "GatePrecondition still has no production caller beyond its declaration; files: $(cat /tmp/0113-t10-callers.txt)"; exit 1; }; grep -q 'internal/cli' /tmp/0113-t10-callers.txt || { echo 'the classifier is not called from the command surface that runs the sweep'; cat /tmp/0113-t10-callers.txt; exit 1; }` — expected: exits 0, proving the classifier is reached from the CLI and not only declared. Fails today, where only its own file resolves.
- `go test -count=1 ./internal/cli -run '^TestRunImplementHasNoSpecCheckPrecondition$' -v > /tmp/0113-t10-implement.log 2>&1; s=$?; grep -q '^--- PASS: TestRunImplementHasNoSpecCheckPrecondition' /tmp/0113-t10-implement.log || { cat /tmp/0113-t10-implement.log; exit 1; }; grep -rq 'GatePrecondition(' internal/cli || { echo 'the Implement contract holds, but the classifier still has no CLI caller'; exit 1; }; exit $s` — expected: exits 0, proving the Implement contract survives the change that gives the classifier its caller. The contract alone passes today, so it is anchored to the caller.

## Context

- interface: `internal/cli/spec_check.go`
- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → System Architecture, the gate's contract execution. `_prd.md` →
Core Feature 6; Goal 4; User Story 4. ADR-0134.
Evidence: this Spec's QA reports `qa/qa-report-2026-08-15.md` finding F-001 and
`qa/qa-report-2026-08-15-02.md` finding F-001.

## Result

The CLI now treats `--strict` with no `--stage` as the gate's full-sweep
boundary. After the ordinary strict promotion, that boundary calls
`GatePrecondition`, keeps its blocking findings in the Spec Consistency Check
result, and renders its declared-term inputs separately. Text output labels each
one `[repair input]`; JSON exposes `repairInputs` without an error severity. A
stage-scoped strict check and a non-strict full check keep their authoring
behaviour.

Acceptance evidence:

- `TestSpecCheckClassifiesTheGateBoundary/authoring_stage_keeps_a_declared_term_as_an_error`
  observed exit `1`, the `SC-VOCABULARY-UNDOCUMENTED` error, and no repair-input
  label for `--stage techspec --strict`.
- `TestSpecCheckClassifiesTheGateBoundary/strict_full_sweep_reports_a_declared_term_as_repair_input`
  observed exit `0`, named `publish:` as a repair input in text and JSON, left
  `findings` empty, and emitted no severity on the JSON repair input.
- `TestSpecCheckClassifiesTheGateBoundary/term_with_no_Spec_declaration_remains_an_error_in_both_modes`
  fed the established synthetic `orphan:` negative control through both boundary
  modes; each retained the error and produced no repair input. The lower
  `TestGateAcceptsItsOwnDeclaredTerm` suite also passed.
- `rtk rg -n 'GatePrecondition\\(' internal --glob '*.go'` resolved the
  classifier declaration in `internal/speccheck/mechanical.go` and its production
  caller in `internal/cli/spec_check.go`, in addition to tests.
- The unchanged `TestRunImplementHasNoSpecCheckPrecondition` passed in the same
  focused run as the new boundary test. No Implement or Daemon source changed.

Focused checks:

- Before the production edit,
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task10-gocache go test ./internal/cli -run '^TestSpecCheckClassifiesTheGateBoundary$'`
  exited `1` because `classifySpecCheckBoundary` did not exist.
- After implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-0113-task10-gocache go test ./internal/cli ./internal/speccheck -run '^(TestSpecCheckClassifiesTheGateBoundary|TestGateAcceptsItsOwnDeclaredTerm|TestRunImplementHasNoSpecCheckPrecondition|TestRunSpecCheckGapStrictPromotion|TestSpecCheckStageExitsNonZeroOnAFinding|TestRunSpecCheckJSONWritesOneObjectPerSpec)$'`
  exited `0` for both packages.
- `rtk env GOCACHE=/private/tmp/roundfix-0113-task10-gocache go test ./internal/cli ./internal/speccheck`
  exited `0` for both affected packages.
- After the final code and test edit, the boundary, lower classifier, unchanged
  Implement contract, unchanged non-strict sweep, and help tests passed together
  in both affected packages; `rtk git diff --check` exited `0`.

The commands under `## Verification` remain for Daemon Verification and were not
run in this Agent turn.
