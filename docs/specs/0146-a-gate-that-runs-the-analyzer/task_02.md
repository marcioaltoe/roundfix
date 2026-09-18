---
task: task_02
spec: 0146-a-gate-that-runs-the-analyzer
status: completed
type: chore
complexity: low
---

# Task 02: Compose the analyzer into the gate

## Overview

With the control in place, the gate gains the step it asserts. This Task is an
authorized tooling mutation and may change only the bounded files plus its own
Task file.

## Requirements

1. MUST add a target named `vet` that runs the Go toolchain's analyzer over the
   module, and MUST include it in the full repository Verification and in its
   incremental sibling.
2. MUST let the analyzer's own exit status decide the step, with no filter or
   pipeline able to decide it instead.
3. MUST leave every existing gate step with its name, command and order.
4. MUST leave the CI workflow unchanged, because it invokes the gate; the path
   is bounded so the decision can be recorded, not so a second composition is
   created.
5. MUST NOT add a dependency, an analyzer configuration file, or any suppression
   mechanism.

## Subtasks

- [ ] Add the analyzer target.
- [ ] Include it in both gate tiers.
- [ ] Confirm the control now passes and the module-wide run stays silent.

## Acceptance Criteria

- [ ] The control passes, having failed before this Task.
- [ ] The repository Verification runs the analyzer and passes on this module.
- [ ] Every other gate step is unchanged, and the workflow file is untouched.
- [ ] No dependency or configuration file is added.

## Context

- interface: `Makefile`

## Verification

- `grep -q "^vet:" Makefile || exit 1; go test -count=1 -tags repocontract -run "^TestRepositoryGateRunsTheAnalyzer$" ./... > /tmp/0146-control-pass.txt 2>&1; status=$?; grep -q "FAIL: TestRepositoryGateRunsTheAnalyzer" /tmp/0146-control-pass.txt && { cat /tmp/0146-control-pass.txt; exit 1; }; exit $status` — expected: exit 0; the control passes now that the gate composes the analyzer. Before this Task it fails, so the command fails.
- `grep -q "^vet:" Makefile && changed="$(git diff --name-only HEAD -- .github/workflows)" && test -z "$changed"` — expected: exit 0; the gate defines the analyzer target and the bounded workflow stays unchanged. Before this Task the target does not exist, so the command fails.
- `grep -qE "^verify:.*[ \t]vet( |$)" Makefile && grep -qE "^verify-incremental:.*[ \t]vet( |$)" Makefile && make verify` — expected: exit 0; both tiers compose the analyzer and the gate passes on this module. Before this Task neither tier names it, so the command fails.

## References

`_prd.md` → Core Features 1 and 4; User Stories 1 and 3; Goal 1;
Success Metric 2; Project Constraints: Tooling authority;
`_techspec.md` → Implementation Design: The step, What stays; API Contracts 1-2;
Build Order 2; `_authorization.md`; ADR-0014.

## Result

Implemented the analyzer gate composition in `Makefile`:

- Added the phony `vet` target, which runs `go vet ./...` directly so the Go
  toolchain exit status controls the step.
- Added `vet` to `verify` and `verify-incremental` while retaining all existing
  dependencies and their order.
- Left `.github/workflows/` unchanged and added no dependency, analyzer
  configuration, or suppression mechanism.

Focused checks after the edit:

- `make vet` — passed with no analyzer diagnostics.
- `make -n verify` and `make -n verify-incremental` — both expanded to invoke
  `go vet ./...` in the expected gate position.
- `git diff --check` — passed; no workflow files changed.

The Daemon must run the declared Verification commands and settle the Task.
