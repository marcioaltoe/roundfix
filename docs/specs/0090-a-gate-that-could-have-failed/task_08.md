---
task: task_08
spec: 0090-a-gate-that-could-have-failed
status: pending
type: test
complexity: medium
---

# Task 08: Give the CLI fixtures gates that can fail

## Overview

Corrective Task for F-002 of `qa/qa-report-2026-08-09.md`. The pre-work probe
works exactly as designed and refuses every Task whose Verification passes
against the unchanged tree — including 66 CLI test fixtures whose declared
command is literally `true`. The probe is right and the fixtures are wrong: a
fixture asserting an Implement journey with a gate that cannot fail was never
exercising the journey it claimed.

The fixtures are built centrally, so this is surgical rather than spread: one
helper writes the `## Verification` section for the Task files these tests use.

## Requirements

1. MUST give the Task fixtures a Verification command that fails against the
   unchanged tree and passes once the fixture's Agent has done its work, so the
   probe clears them for the reason the probe exists.
2. MUST NOT weaken, disable, or bypass the probe for tests. A fixture that needs
   a vacuous gate is a fixture asserting a journey it does not exercise.
3. MUST keep every assertion these tests make about Implement journeys, exit
   codes, and Run outcomes; this Task changes what the fixtures declare, never
   what the tests conclude.
4. MUST keep the fixtures that deliberately exercise a refused Task — the probe's
   own regression — refusing for their own reason and not by accident.
5. MUST leave `internal/daemon` fixtures alone unless they carry the same defect,
   and say which were changed and why.

## Subtasks

- [ ] Change the central fixture builder's Verification to a failing-then-passing
      command.
- [ ] Fix any fixture that declares its command inline rather than through the
      builder.
- [ ] Confirm the probe's own regression cases still refuse.

## Acceptance Criteria

- [ ] `internal/cli` passes with the probe active.
- [ ] No fixture declares a Verification command that exits zero against the
      unchanged tree.
- [ ] The probe is neither disabled nor special-cased for tests.
- [ ] The refused-Task regression still refuses.

## Rehearsal Cases

- Case: the ordinary Implement journey fixture; Observation: the probe clears
  the Task because its command fails against the unchanged tree, and the test's
  original assertions about exit code and Run outcome are unchanged.
- Case: the probe's own regression fixture carrying a deliberately vacuous gate;
  Observation: the Task is still refused, and the terminal reason still names the
  offending command.
- Case: a fixture whose Agent does nothing; Observation: the Verification fails
  after the Agent turn, exactly as before this Spec.

## Bounded scope

This Task may create or modify only:

- `internal/cli/implement_test.go`
- `internal/cli/settle_test.go`
- `internal/cli/cli_test.go`
- `internal/daemon/task_engine_test.go`
- `internal/daemon/daemon_test.go`
- `docs/specs/0090-a-gate-that-could-have-failed/task_08.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -count=1 2>&1 | tee /dev/stderr | grep -qE '^ok[[:space:]]+roundfix/internal/cli'` — expected: exits 0, proving the package passes with the probe active.
- `test -z "$(grep -rnE '"## Verification\\n\\n- \`true\`|Verification: *"true"' internal/cli/*_test.go)"` — expected: exits 0, proving no fixture still declares a gate that cannot fail.

The probe's own regression is asserted under Acceptance Criteria and Rehearsal
Cases rather than here. It is an invariant that holds before and after this
Task, so as a Verification command it would pass against the unchanged tree —
which is exactly what the probe refuses, and it refused this Task for that
reason on its first dispatch, at a cost of zero Agent turns.

## References

- `_prd.md` → Goal 1.
- `_techspec.md` → Build Order 3; Risks.
- `qa/qa-report-2026-08-09.md` → F-002.
- ADR-0109.
