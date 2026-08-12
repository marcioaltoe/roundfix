---
task: task_08
spec: 0090-a-gate-that-could-have-failed
status: completed
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

One command, deliberately. It fails today — the package does not pass with the
probe active — and passes only once every fixture declares a gate that can fail,
which is this Task's whole effect. A second command checking for the literal
`true` was tried twice and dropped: it needs backticks, and the Verification
contract extracts commands from a backticked span, so the pattern cannot survive
the round trip. A fragile check beside a sufficient one buys nothing.

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

## Result

### Implementation

- `internal/cli/implement_test.go` now gives the shared Task fixture a
  per-Task Verification marker under Git metadata. The command exits non-zero
  on the unchanged Task Worktree and the scripted fake Agent records the marker
  only after its Task work returns successfully. The detached fake ACP runtime
  records the same marker from the `Task:` identity in its prompt.
- Inline CLI fixtures that previously used `echo`, an already-created bootstrap
  file, or an unconditional `printf` now use the same Agent-work gate. Blocking
  and temporary-failure commands first require that gate, so the pre-work probe
  cannot block or mutate their later post-Agent scenario.
- `internal/cli/cli_test.go` delegates the shared marker command to the real
  `daemon.ExecVerifier`. Scripted verifier failures still apply after the marker
  exists, preserving the existing journey, exit-code, and Run-outcome
  assertions.
- Added the negative companion
  `TestRunImplementAgentThatDoesNoWorkStillFailsVerification`: an Agent with no
  scripted work receives its Verification Feedback turn, fails the marker gate
  again, and leaves the Run Unresolved.
- `internal/daemon/task_engine_test.go`, `internal/daemon/daemon_test.go`, and
  `internal/cli/settle_test.go` were not changed. The daemon fixtures already
  model the deliberate Vacuous Verification refusal correctly; Settle fixtures
  verify preserved work from an earlier Agent turn rather than a pre-Agent
  Implement journey.

### Focused checks

- Red baseline:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementExecutesSpecEndToEnd$' -count=1 -v`
  exited 1 before the edit because `echo docs-check` was reported as a vacuous
  pre-work Verification and `task_02` was skipped.
- Ordinary and inline journeys:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run 'Implement' -count=1`
  exited 0 after the final edit (`ok roundfix/internal/cli 7.239s`). This focused
  sweep covers the central Task builder, detached child, bootstrap, temporary
  retry, failure-report, worktree, QA, and Run-outcome Implement cases.
- No-work negative companion:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementAgentThatDoesNoWorkStillFailsVerification$' -count=1 -v`
  exited 0; the test observed two no-op Agent turns and an Unresolved Run after
  the marker Verification failed twice.
- Refused-Task regression:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^(TestPreWorkProbeRefusesATaskWhoseGateAlreadyPasses|TestPreWorkProbePublishesEveryOffendingCommand|TestPreWorkProbeSpendsNoAgentTurnOnARefusedTask)$' -count=1`
  exited 0 (`ok roundfix/internal/daemon 0.252s`). The deliberately vacuous gate
  still names every offending command and spends no Agent turn.
- Vacuous inline-fixture scan:
  `rtk proxy rg -n 'verification = \\[\\]string\\{\"true\"\\}|echo docs-check|echo backend-check|verification: \\[\\]string\\{\"test -f bootstrap\\.ready\"\\}|verification: \\[\\]string\\{\"printf .independent pass' internal/cli/implement_test.go`
  exited 1 with no matches, the expected absence result.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.

### Acceptance evidence

- `internal/cli` with the probe active: the focused Implement-named sweep exits
  0 after the final edit. The complete package command in `## Verification` was
  not run; the Daemon owns that terminal check.
- No work-independent fixture gate: the shared command requires the Agent marker;
  inline unconditional/bootstrap gates are absent, the ordinary journey passes,
  and the no-work negative companion remains Unresolved.
- Probe integrity: no production probe code or daemon fixture changed. The test
  verifier runs the marker command through `daemon.ExecVerifier`; it does not
  disable or force a pre-work verdict.
- Refused Task: the three focused daemon regressions pass without any daemon
  changes, including exact offending-command publication and zero Agent turns.

### Verification feedback attempt 1

- The Daemon diagnostic isolated
  `TestAgentSelectionProfilesMacro/mixed_profiles_configure_validate_fallback_persist_and_stream`:
  its separate Python ACP fixture still marked Task frontmatter completed but
  did not record the new Agent-work marker, so both backend and frontend gates
  correctly remained non-zero after their Agent turns.
- `macroFakeACPXScript` now derives the Task identity from the prompted Task file
  and writes the same Git-metadata marker only after its successful scripted
  Task work. A rejected Agent prompt exits before this write, preserving the
  macro fixture's post-start failure assertions.
- Exact reported regression:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestAgentSelectionProfilesMacro$/^mixed_profiles_configure_validate_fallback_persist_and_stream$' -count=1 -v`
  exited 0 (`ok roundfix/internal/cli 6.014s`).
- Macro success and refusal/failure companions:
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestAgentSelectionProfilesMacro$' -count=1 -v`
  exited 0 (`ok roundfix/internal/cli 8.156s`), including the invalid-Task
  preflight refusal and the post-start Agent failure with no fallback.
- The declared `## Verification` command was not rerun; the Daemon owns that
  terminal package check.
