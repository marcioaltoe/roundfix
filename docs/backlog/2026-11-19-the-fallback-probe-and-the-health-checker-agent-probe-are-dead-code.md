---
type: refactor
status: open
created: 2026-11-19
spec: null
reason: null
---

# The fallback candidate probe and the health checker's Agent probe are dead code

## Tangled

No command reaches either path, yet both stay compiled, tested and documented.

- The fallback candidate probe: `probeRuntimeSelection` (`internal/cli/selection.go:94`) is the only caller of `ProbeFallback` (`internal/agent/agent.go:675`, `internal/agent/fallback.go:26`), which calls `probeFallbackCandidate` (`internal/agent/fallback.go:51`). No code calls `probeRuntimeSelection`: commit `d7fdd8cf` (#32, 2026-07-17, "feat: add release planning and agent selection profiles") removed its last five calls, when the Fallback Chain started to be proven by profile readiness before a Run. The rest of `internal/cli/selection.go` around it, the confirmation prompt, `fallbackConfirmationInput` and `fallbackConfirmationAvailable` in `internal/cli/cli.go:303-304` and the fake `ProbeFallback` methods in `internal/cli/implement_test.go:977` and `internal/cli/cli_test.go:14025`, serve only that path. `probeFallbackCandidate` is also the only proof that still uses the constant `acpxPreflightSetupTimeout` instead of a configured deadline (`internal/agent/fallback.go:58,64,96`).
- The health checker's Agent probe: `runtimeHealthChecker.Agent` (`internal/cli/health.go:186`) and the `probeAgent` dependency it calls (`internal/cli/health.go:73`, wired at `internal/cli/setup.go:586,601`) have no non-test caller; only `internal/cli/cli_test.go:4768` and `:4840` call `Agent`.

Measured by the author of Spec 0411 on 2026-11-18 at `a6a4a15b` with `rg -n 'probeRuntimeSelection|ProbeFallback|\.Agent\(ctx' internal` and `git show d7fdd8cf`.

## Cost

A reader takes the probe for live behavior: the configuration guide and the Roundfix skill's `profiles` reference say "The Fallback Chain probe keeps the built-in setup deadline", ADR-0298 and ADR-0344 decide its deadline, and Spec 0411's brief listed it as a defect before the measurement showed no command reaches it. Each change to the proof or the deadline must keep the dead path compiling and its tests (`TestACPXProbeFallback…`, `TestAFallbackCandidateKeepsItsResultWhenItsCloseFails`) green.

## Shape

Non-binding: remove `probeRuntimeSelection` and the fallback confirmation helpers, `ProbeFallback`, `probeFallbackCandidate` and their tests, `runtimeHealthChecker.Agent` with the `probeAgent` dependency, re-record `docs/references/coverage-record.json` for the removed tests, drop the sentence on the Fallback Chain probe from the guide and the skill, and note the removal in ADR-0298 and ADR-0344.

## Decisions

- 2026-11-19, about 00:05 BRT (AskUserQuestion during the planning of Spec 0411): the maintainer chose to record it as future work — "(a) remover código morto".
