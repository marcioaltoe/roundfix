---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: invalid
terminal_reason: "The test already uses fakeWatchClock and fakeWatchSleeper through withSuccessfulPreflight, so its nanosecond durations do not consult wall-clock time."
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/cli/cli_test.go
line: 6763
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wuazg,comment:PRRC_kwDOS0qyts7d35sw
review_hash: 55716199fe9a17fc490cf4f67f33f31aa9a05a9cc1925898111608aeed64f473
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 001: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

<details>
<summary>🧩 Analysis chain</summary>

🏁 Script executed:

```shell
#!/bin/bash
# Description: Look for clock and sleeper injection helpers used by CLI tests.
set -euo pipefail

rg -n -C 4 'watchClock|watchSleeper|withWatchClock|withWatchSleeper' --glob '*.go'
```

Repository: marcioaltoe/roundfix

Length of output: 158

---

🏁 Script executed:

```shell
#!/bin/bash
# First, search for how fakeClock and fakeSleeper are used in watch_test.go
rg -n 'fakeClock|fakeSleeper' internal/watch/watch_test.go -C 3

echo "---"

# Check if there are other clock-related patterns in CLI tests
rg -n 'Clock|Sleeper|clock|sleeper' internal/cli/cli_test.go | head -20

echo "---"

# Look for common test setup patterns that might handle time
rg -n 'withWatch|withCLI' internal/cli/cli_test.go | head -10
```

Repository: marcioaltoe/roundfix

Length of output: 17062

---

🏁 Script executed:

```shell
#!/bin/bash
# Examine the test file at the specific lines to understand context
sed -n '6730,6775p' internal/cli/cli_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 1870

---

🏁 Script executed:

```shell
#!/bin/bash
# Search for how time/clock is handled in CLI test setup
rg -n 'poll_interval|review_timeout|quiet_period' internal/cli/cli_test.go | head -15
```

Repository: marcioaltoe/roundfix

Length of output: 492

---

🏁 Script executed:

```shell
#!/bin/bash
# Look at the withWatchTiming helper and how it injects the clock
sed -n '12414,12425p' internal/cli/cli_test.go
```

Repository: marcioaltoe/roundfix

Length of output: 665

---

🏁 Script executed:

```shell
#!/bin/bash
# Check if other CLI tests successfully use withWatchTiming
rg -n 'withWatchTiming' internal/cli/cli_test.go -A 2 | head -20
```

Repository: marcioaltoe/roundfix

Length of output: 719

---



</details>

**Use the `withWatchTiming` clock injection helper instead of nanosecond real-time timeouts.**

This test writes `poll_interval: 1ns`, `review_timeout: 2ns`, and `quiet_period: 1ns` to `.roundfixrc.yml` and relies on the real clock. The test is stable in practice because these intervals always fire, but it couples the test to wall-clock behavior instead of deterministic time control.

The CLI test harness exposes `withWatchTiming(t, clock, sleeper)` to inject a `fakeWatchClock` and `fakeWatchSleeper`. Two other tests at lines 6352 and 11980 use this same pattern. Apply it here: construct a `fakeWatchClock`, create a `fakeWatchSleeper`, and call `withWatchTiming` early in the test. This makes timing assertions deterministic and aligns with the coding guideline to dependency-inject clocks in time-dependent tests.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 6742 - 6763, Update
TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush to inject
deterministic timing via withWatchTiming early in the test: construct a
fakeWatchClock and fakeWatchSleeper, following the existing patterns in nearby
tests. Remove the nanosecond poll_interval, review_timeout, and quiet_period
values from the test configuration, while preserving the existing watch evidence
and assertion behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:8dce8260379f39f6a076ca5f -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: `TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush` calls
  `withSuccessfulPreflight` before running the command. That helper installs a
  `fakeWatchClock` and a `fakeWatchSleeper` through `withWatchTiming`, so the
  configured nanosecond durations advance only the injected clock. Adding the
  same injection again would duplicate existing setup rather than remove a
  wall-clock dependency.
- Focused evidence: inspected the test, `withSuccessfulPreflight`, and
  `withWatchTiming`; `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./internal/cli -count=1 -run
  '^TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush$'` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
