# Dogfood findings — round 4

Round driving spec `0024-context-efficient-runs` on branch
`ma/context-efficient-runs` (codex / gpt-5.5 / xhigh, concurrency 1, `--qa`).
This Run doubles as the live test of the previous round's shipped work
(PR #20: Agent model and reasoning selection, run browser, cockpit fidelity).

## Findings

1. **`roundfix events` hard-fails on pre-0024 Run Event Journals.** Replaying
   the Run that shipped the feature itself
   (`run_20260710T200704Z_1362ac67209e572a`, journal written by the pre-0024
   binary) emits one record, then aborts with exit 1:
   `roundfix events failed: project daemon.verification event for Run "…": missing payload field "attempt"`.
   Legacy `daemon.verification` events lack the new `attempt` payload field and
   the projection treats that as fatal instead of tolerating or skipping legacy
   events. New-binary Runs carry the field, so the contract works forward; the
   gap is backward compatibility with retained journals inside the 14-day
   retention window.

## Run summary

- Run `run_20260710T200704Z_1362ac67209e572a`: all 6 Tasks completed, one
  Daemon verification pass each, `qa pass`
  (`docs/specs/0024-context-efficient-runs/qa/qa-report-2026-07-10.md`),
  outcome `Clean: all 6 Task(s) completed.`, auto-integrated onto
  `ma/context-efficient-runs`.
- Live confirmation of the PRD's motivating incident: supervising this Run by
  grepping the Console Log produced repeated false failure signals from echoed
  source lines and raw ACP diff payloads (`failed:`, `Reason:`, `Stopped`
  inside code and docs the Agent read or edited).
- `make verify` on the integrated HEAD: 1097 tests in 19 packages, skills
  check, build — all green on the rebuilt binary.
