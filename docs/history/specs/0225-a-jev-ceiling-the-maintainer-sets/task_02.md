---
task: task_02
spec: 0225-a-jev-ceiling-the-maintainer-sets
status: completed
type: backend
complexity: medium
---

# Task 02: The judge and the Jev Router gate read the User Config Jev ceiling

## Overview

The monthly Jev ceiling is compiled into the binary at US$5. This Task adds
the User Config value `jev.monthly_ceiling_usd` and the one rule that applies
it, and passes the loaded value to `roundfix spec judge` and to the Jev
Router gate of the Implement Run and resolve engines, as the TechSpec's "The
configured ceiling", "The ceiling rule" and "The consumers" state. An unset
value keeps US$5 and every message as before; a Project Config value is
ignored with a warning (ADR-0231).

## Requirements

1. MUST add `Jev`, its `MonthlyCeilingUSD` field and `jevOverlay` with the
   shapes of the TechSpec's Interfaces, decode `jev.monthly_ceiling_usd`
   from User Config, and leave it zero in `Builtin()` and when unset.
2. MUST remove `jev.monthly_ceiling_usd` from a Project Config document
   before decoding and print API Contract 2 through the existing
   ignored-setting warning, as `runs.max_active` is handled; and MUST refuse
   a User Config value that is not a finite number greater than zero with
   API Contract 3 inside the existing `parse config "<path>": ` wrapping.
3. MUST add `Questions.WithMonthlyCeiling` with the TechSpec's contract and
   use it as the only place a configured value replaces the embedded one;
   MUST NOT change `questions.json`, `judge.Run`, `jevrouter.MonthSpend` or
   `Spend.CheckKeyLimit`.
4. MUST make `runSpecJudgeCommand` apply the loaded value, and replace the
   usage sentence "The monthly ceiling is US$5.00" with one naming
   `jev.monthly_ceiling_usd` in User Config and "US$5 by default", keeping the
   phrase "monthly ceiling".
5. MUST add `JevMonthlyCeilingUSD` to `daemon.Dependencies`, make `NewEngine`
   build the default gate's ceiling through `WithMonthlyCeiling`, and set the
   field from the loaded configuration in the Implement Run engine in
   `internal/cli/implement.go` and the resolve engine in `internal/cli/cli.go`.
6. MUST add the tests named in Verification, as Testing Approach 1 to 4
   describe: `internal/config/jev_ceiling_test.go` (new) for the key, its
   default, the Project Config warning and the refusals of 0, -1 and `.inf`;
   `internal/judge/judge_test.go` for the rule and a run at a ceiling of 50;
   `internal/cli/spec_judge_test.go` for Surface Transcript 1's summary with
   JSON `month_ceiling_usd` 50, for Surface Transcript 3's warning and for
   Surface Transcript 4's refusal with exit 2; and
   `internal/daemon/jev_router_gate_test.go` for the default gate's ceiling
   and a key limit of 50 accepted at a ceiling of 50 and refused with
   `jev_router_key_unbounded` at the built-in one.
7. MUST rebuild from the loaded ceiling every existing test literal that
   spells the US$5 ceiling (`US$5.0000 of US$5.00`, `of US$5.00`,
   `US$5.0000`) in `internal/judge/judge_test.go`,
   `internal/cli/spec_judge_test.go` and
   `internal/daemon/jev_router_gate_test.go`, so a change of the default
   moves them, and extend `TestSpecJudgeHelp` to require
   `jev.monthly_ceiling_usd`.
8. MUST NOT change `.roundfixrc.yml`, the `roundfix config init` templates,
   any other configuration key, or write any file outside the repository;
   every test uses a disposable home and a local key endpoint.

## Subtasks

- [ ] Add the configuration value, its Project Config removal and its validation.
- [ ] Add the ceiling rule and apply it in the judge command and the engines.
- [ ] Update the judge usage text.
- [ ] Add the tests and rebuild the existing ceiling literals.

## Acceptance Criteria

- [ ] `jev.monthly_ceiling_usd: 50` in User Config loads as 50; unset loads
      as zero; a Project Config value is ignored with API Contract 2; 0, -1
      and `.inf` fail with API Contract 3.
- [ ] `roundfix spec judge` with a User Config ceiling of 50 and a month at
      US$50.25 prints Surface Transcript 1's summary; with nothing set it
      prints Surface Transcript 2's.
- [ ] A Jev Router gate built with `JevMonthlyCeilingUSD` 50 accepts a key
      whose monthly limit is 50; without it the key is refused naming the
      built-in ceiling.
- [ ] No test literal in the three test files spells the US$5 ceiling.

## Context

- interface: `internal/config/config.go`
- creates: `internal/config/jev_ceiling_test.go`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/judge_test.go`
- interface: `internal/cli/spec_judge.go`
- interface: `internal/cli/spec_judge_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/jev_router_gate_test.go`
- instruction: `docs/user-guide/configuration.md`
- instruction: `docs/user-guide/commands/spec.md`

## Verification

- `grep -q JevMonthlyCeilingUSD internal/cli/implement.go || { printf 'implement engine does not pass the ceiling\n' >&2; exit 1; }; grep -q JevMonthlyCeilingUSD internal/cli/cli.go || { printf 'resolve engine does not pass the ceiling\n' >&2; exit 1; }; if grep -qF -- 'US$5.0' internal/judge/judge_test.go internal/cli/spec_judge_test.go internal/cli/spec_judge.go internal/daemon/jev_router_gate_test.go; then printf 'a literal US$5 ceiling is left\n' >&2; exit 1; fi; out="$(go test -count=1 -v -run "^(TestJevMonthlyCeilingIsReadFromUserConfig|TestJevMonthlyCeilingIsUnsetByDefault|TestProjectConfigCannotSetTheJevCeiling|TestJevMonthlyCeilingRefusesANonPositiveOrInfiniteValue|TestWithMonthlyCeilingReplacesOnlyAPositiveCeiling|TestRunStopsAtAConfiguredCeiling|TestRunStopsAtTheMonthlyCeiling|TestRunRechecksSpendBeforeRetryAndHonorsCancellation|TestSpecJudgeSkipsAtTheUserConfigCeiling|TestSpecJudgeIgnoresAProjectConfigCeiling|TestSpecJudgeRefusesAnInvalidUserConfigCeiling|TestSpecJudgeSkipsAtTheMonthlyCeiling|TestSpecJudgeHelp|TestJevRouterDefaultGateUsesTheConfiguredCeiling|TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling|TestJevRouterDefaultGateUsesProcessHomeAndJudgeCeiling|TestJevRouterGateChecksTheReportedKeyLimit|TestJevRouterCeilingFallsBackBeforeWork|TestJevRouterUnboundedKeyFallsBackBeforeWork)$" ./internal/config ./internal/judge ./internal/cli ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestJevMonthlyCeilingIsReadFromUserConfig TestJevMonthlyCeilingIsUnsetByDefault TestProjectConfigCannotSetTheJevCeiling TestJevMonthlyCeilingRefusesANonPositiveOrInfiniteValue TestWithMonthlyCeilingReplacesOnlyAPositiveCeiling TestRunStopsAtAConfiguredCeiling TestRunStopsAtTheMonthlyCeiling TestRunRechecksSpendBeforeRetryAndHonorsCancellation TestSpecJudgeSkipsAtTheUserConfigCeiling TestSpecJudgeIgnoresAProjectConfigCeiling TestSpecJudgeRefusesAnInvalidUserConfigCeiling TestSpecJudgeSkipsAtTheMonthlyCeiling TestSpecJudgeHelp TestJevRouterDefaultGateUsesTheConfiguredCeiling TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling TestJevRouterDefaultGateUsesProcessHomeAndJudgeCeiling TestJevRouterGateChecksTheReportedKeyLimit TestJevRouterCeilingFallsBackBeforeWork TestJevRouterUnboundedKeyFallsBackBeforeWork; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither engine passes the ceiling, the test files spell US$5, and the new tests do not exist, so the command fails.

## References

- `_prd.md` → User Stories 1-5; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → The configured ceiling; The ceiling rule; The consumers; Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Testing Approach 1-4; Build Order 2
- ADR-0231; ADR-0201; ADR-0218; ADR-0027

## Result

Implemented the User Config ceiling and its shared application rule. The
judge command, Implement Run engine and resolve engine receive the loaded
value; an unset value preserves the embedded ceiling. Project Config removes
the setting before strict decoding and warns once. Invalid User Config
values return the specified error with the configuration path.

Acceptance evidence from focused implementation checks:

- Config loading: `TestJevMonthlyCeilingIsReadFromUserConfig` loads 50;
  `TestJevMonthlyCeilingIsUnsetByDefault` proves loaded and built-in values
  remain zero. `TestProjectConfigCannotSetTheJevCeiling` preserves either
  the User Config value or zero and checks the exact warning, including
  invalid ignored Project Config values.
  `TestJevMonthlyCeilingRefusesANonPositiveOrInfiniteValue` checks 0, -1,
  infinities, NaN, null and non-number values against the exact wrapped error.
- Judge output: `TestSpecJudgeSkipsAtTheUserConfigCeiling` checks the
  US$50.25/US$50 text transcript, JSON `month_ceiling_usd: 50`, and no HTTP
  calls. `TestSpecJudgeSkipsAtTheMonthlyCeiling` checks the unchanged
  default behavior using the loaded ceiling plus US$0.25.
  `TestSpecJudgeIgnoresAProjectConfigCeiling` checks the exact warning;
  `TestSpecJudgeRefusesAnInvalidUserConfigCeiling` checks exit 2.
  `TestSpecJudgeHelp` now requires the key name.
- Router gate: `TestJevRouterDefaultGateUsesTheConfiguredCeiling` checks
  the default gate with 50 and with zero.
  `TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling` serves a local
  key endpoint with monthly limit 50, observes acceptance at ceiling 50,
  and checks the exact `jev_router_key_unbounded` refusal naming the
  embedded ceiling when unset.
- Default literals: existing judge, command and router expectations and
  ceiling-reaching spend fixtures now derive from `judge.Load()`.
  A Python inspection of all three test files found no `US$5.0` literal.
  `TestWithMonthlyCeilingReplacesOnlyAPositiveCeiling` checks other settings
  and the original value are preserved; `TestRunStopsAtAConfiguredCeiling`
  checks requests below 50 and no requests at 50.

Focused check command:

```sh
GOCACHE="$PWD/.gocache" rtk proxy go test -count=1 ./internal/config ./internal/judge ./internal/cli ./internal/daemon -run 'JevMonthlyCeiling|ProjectConfigCannotSetTheJevCeiling|WithMonthlyCeiling|ConfiguredCeiling|SpecJudge|JevRouter|RunStopsAtTheMonthlyCeiling|RunRechecksSpend'
```

Outcome: exit 0 for all four packages. The initial config test failed to
compile because `Config.Jev` was absent, establishing the pre-change gap.
The sandbox initially refused local HTTP binding; the focused run succeeded
with the required sandbox escalation. All test endpoints are local or fake;
no paid API call was made.

`git -c core.fsmonitor=false diff --check` exited 0. The initial working
tree contained only the Daemon's task_02 status edit; it was preserved.
No other Task, Task Graph, config template or repository Project Config was
edited. Declared Verification remains for the Daemon; Task status was not
changed by this Agent.

The first `rtk make verify-incremental` run exited 2: test assertions
passed, but the repository suiteguard detected this Agent appending the
Result while the baseline, CLI and daemon tests were running. This was a
check-execution error, not a test assertion failure. The incremental check
was rerun without concurrent file edits. RTK's recovery hint
`rtk recall dee832fa96ab` returned `store unavailable: init recall schema`.

The rerun, `rtk proxy make verify-incremental`, exited 0: formatting,
`go vet ./...`, repository tests (including CLI and daemon suiteguards),
shipped skill synchronization and validation, and the CLI build all passed.
The final Result update was made only after that command exited.

## Carry-forward provenance

- Source Run: `run_20261004T194332Z_e41f982109e26c75`
- Source commit: `d2bf17e795e9b8fb19957e7eac43bec06cbee592`
