---
task: task_04
spec: 0096-a-failure-the-agent-can-read
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 04: Let the vacuity refusal name its offenders

## Overview

The event that reports a Verification refused as vacuous summarises a count and
carries the offending commands under a key named `commands`, which reads as the
commands the tool ran. A reader takes the list for the wrong thing and looks for
a defect in commands that were fine.

## Requirements

1. MUST name the offending commands in the event summary rather than only their
   number.
2. MUST carry every probed command with its own verdict, under a key whose name
   says what the list holds.
3. MUST point at the probe log that settles what ran.
4. MUST leave the event's classification and phase unchanged, so existing readers
   keep working.

## Subtasks

- [ ] Name the offenders in the summary.
- [ ] Carry every command with its verdict under a self-describing key.
- [ ] Include the probe log path.

## Acceptance Criteria

- [ ] The summary names the offending commands.
- [ ] The payload lists every probed command with its verdict.
- [ ] The payload names the probe log path.
- [ ] The classification and phase values are unchanged.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestVacuousPreWorkEvent' -v > /tmp/0096-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestVacuousPreWorkEvent' /tmp/0096-t04.log || { cat /tmp/0096-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0096-t04.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t04.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0096-t04.log > /tmp/0096-t04-n.txt; test "$(cat /tmp/0096-t04-n.txt)" -ge 3 || { echo "expected the summary, the per-command verdicts, and the probe log cases, got $(cat /tmp/0096-t04-n.txt)"; cat /tmp/0096-t04.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each assertion runs on its own.
- `grep -n '"commands":' internal/daemon/task_engine.go && { echo 'the payload still uses the ambiguous commands key'; exit 1; }; grep -q 'VerificationClassificationVacuous' internal/daemon/task_engine.go || { echo 'the vacuous classification was removed rather than kept'; exit 1; }; exit 0` — expected: exits 0, proving the ambiguous key is gone and the classification survives. It prints the offending line on failure. Fails today.

## Context

- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 4; System Architecture, the vacuity event.
`_prd.md` → Core Feature 5. ADR-0111.

## Result

### Implementation

- The vacuous pre-work event summary now quotes every command that exited zero
  against the unchanged tree instead of reporting only their count.
- Replaced the ambiguous `commands` payload field with `probed_commands`. Each
  entry records the authored command, its `passed`, `failed`, or `unknown`
  verdict, and its numbered `probe_log_path`.
- Kept the event's `vacuous` classification and `failed` phase unchanged.

### Focused-check evidence

- Red signal before the production edit:
  `GOCACHE=/tmp/roundfix-0096-task04-go-cache rtk go test ./internal/daemon -run '^TestVacuousPreWorkEvent$/summary_names_offending_commands$'`
  failed because the summary contained only `2 commands` and named neither
  offender.
- `GOCACHE=/tmp/roundfix-0096-task04-go-cache rtk go test ./internal/daemon -run '^TestVacuousPreWorkEvent$' -v`
  passed all four subtests and the parent test.
- `GOCACHE=/tmp/roundfix-0096-task04-go-cache rtk go test ./internal/daemon -run '^TestPreWorkProbe(RefusesATaskWhoseGateAlreadyPasses|PublishesEveryOffendingCommand)$' -v`
  passed both existing pre-work refusal cases against the new payload contract.
- `GOCACHE=/tmp/roundfix-0096-task04-go-cache rtk go test ./internal/daemon`
  passed all 230 daemon tests.
- `GOCACHE=/tmp/roundfix-0096-task04-go-cache rtk make verify-incremental`
  first reached the full suite but the sandbox denied process-table access to
  two unrelated force-stop integration tests. Re-running with process-table
  access exited 0: all Go packages, the focused skill contract,
  `roundfix skills check`, and the production build passed.

### Acceptance evidence

1. `TestVacuousPreWorkEvent/summary_names_offending_commands` passed after
   asserting that the event summary contains both offending commands.
2. `TestVacuousPreWorkEvent/payload_carries_every_probed_command_verdict`
   passed after asserting the ordered `passed`, `failed`, and `passed` verdicts
   for all three probed commands.
3. `TestVacuousPreWorkEvent/payload_names_each_probe_log` passed after asserting
   every entry's numbered `probe_log_path` under the Run artifact directory.
4. `TestVacuousPreWorkEvent/classification_and_phase_stay_unchanged` passed with
   `VerificationClassificationVacuous` and `VerificationPhaseFailed`.

### Not run

- The commands under this Task's `## Verification` section were not run; the
  Roundfix Daemon owns declared Verification and Task settlement.
