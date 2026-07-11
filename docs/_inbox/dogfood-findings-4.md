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

2. **Roundfix cannot drive the gpt-5.6 model family.** codex-acp 0.16.0 only
   advertises the `reasoning_effort` session config option when a model preset
   supports more than one effort; the gpt-5.6 family manages effort itself, so
   every `session/set_config_option "reasoning_effort"` value (`minimal`,
   `low`, `medium`, `high`, `xhigh`, `none`) is rejected with ACP -32602 for
   `gpt-5.6-sol`. Roundfix requires a non-empty `reasoning_effort`
   (`validateRuntimeSelection`) and unconditionally issues the set call
   (`applyDisposableSelection` / `applySelection` in
   `internal/agent/acpx_runner.go`), so Agent selection fails for any 5.6
   model. Fix direction: treat `reasoning_effort` as optional and skip the set
   call when the model does not expose the option (or tolerate the rejection),
   instead of hard-failing selection.
3. **Upstream adapter notes.** zed-industries/codex-acp is deprecated —
   development moved to `@agentclientprotocol/codex-acp` on the new Codex App
   Server (upstream-only note). The adapter advertises models `gpt-5.6-sol`,
   `gpt-5.5`, `gpt-5.4`, `gpt-5.4-mini`, `gpt-5.3-codex-spark`; the
   `gpt-5.6-terra` and `gpt-5.6-luna` entries in Roundfix's model catalog are
   not advertised by the installed adapter. acpx `set model 'model[effort]'`
   bracket syntax is accepted transiently but breaks session replay ("did not
   advertise that model") — not a viable workaround. Local codex-acp upgraded
   0.15.0 → 0.16.0 during diagnosis (no behavior change for this finding).

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
