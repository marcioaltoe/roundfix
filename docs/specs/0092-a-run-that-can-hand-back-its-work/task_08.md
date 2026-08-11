---
task: task_08
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: test
complexity: medium
---

# Task 08: Name the seventh break the enumeration missed

## Overview

Task 01 enumerated six tests that Task 02's new work-started boundary breaks.
The boundary breaks seven. `TestAgentSelectionProfilesMacro/post_start_failure_never_activates_fallback`
in `internal/cli/implement_test.go` passes on `main` and fails at `94ff10b2`,
proved by running that single case at `3fc542bf` (passes) and at `94ff10b2`
(fails).

The failure is the contract working, not the code breaking. The case drives the
fake with `ROUNDFIX_FAKE_ACPX_FAIL_PROMPT_MODEL`, which refuses at model
application before the Agent produces any output. Under the old boundary that
counted as started, so no Fallback Selection was eligible. Under Task 02's
Requirement 2 that is a selection failure, and its Requirement 3 requires the
chain to stay eligible. The case name still says `post_start`, but the scenario
it builds is now pre-start.

This Task splits the case in two so both requirements keep a guard, and records
the seventh break in the corpus that claimed six.

## Requirements

1. MUST keep a case that proves a selection failure with no Agent output leaves
   the Fallback Chain eligible, which is Task 02's Requirement 3.
2. MUST keep a case that proves a failure arriving *after* Agent output leaves
   the chain ineligible, which is Task 02's Requirement 4 and the guard the
   current case was written to hold. Produce the Agent output before the
   failure rather than deleting the guard; the fake already exposes
   `ROUNDFIX_FAKE_ACPX_EXIT_BY_CALL` and `ROUNDFIX_FAKE_ACPX_EXIT_BY_COMMAND`
   for failures that arrive later in a turn.
3. MUST name both cases for what they assert, so neither can be mistaken for
   the pre-Task-02 contract.
4. MUST add the seventh entry to the declared-break corpus with the same
   `Outcome contract test:` marker the other six use, stating the assertion it
   held and the contract that replaces it.
5. MUST NOT change production code. If either case cannot be expressed against
   the current fake, stop and report which knob is missing rather than
   weakening the assertion.

## Subtasks

- [ ] Split the case into a selection-failure case and a post-output case.
- [ ] Drive the post-output case so Agent output precedes the failure.
- [ ] Add the seventh corpus entry.

## Acceptance Criteria

- [ ] A selection failure with no Agent output activates the Fallback Chain.
- [ ] A failure after Agent output does not activate it.
- [ ] The corpus enumerates seven outcome-contract tests, not six.
- [ ] `git diff --name-only` lists only this Task's bounded paths.

## Bounded scope

This Task may create or modify only:

- `internal/cli/implement_test.go`
- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_08.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestAgentSelectionProfilesMacro$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^[[:space:]]*--- PASS: TestAgentSelectionProfilesMacro/a_selection_failure_activates_the_fallback_chain'` — expected: exits 0. The case does not exist before this Task, so the command cannot pass against the unchanged tree.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestAgentSelectionProfilesMacro$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^[[:space:]]*--- PASS: TestAgentSelectionProfilesMacro/agent_output_before_failure_keeps_the_chain_ineligible'` — expected: exits 0, proving Requirement 4 kept a guard rather than losing one.
- `test "$(grep -c 'Outcome contract test:' internal/daemon/run_disposition_characterization_test.go)" -eq 7` — expected: exits 0. The corpus enumerates six today, so this fails until the seventh entry lands.

## References

- `_prd.md` → Goal 1.
- `_techspec.md` → Build Order 2.
- `task_01.md` → the six-test enumeration this Task corrects to seven.
- `task_02.md` → Requirements 2, 3 and 4.
