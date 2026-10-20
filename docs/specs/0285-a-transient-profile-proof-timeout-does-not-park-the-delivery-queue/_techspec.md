---
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
prd: _prd.md
created: 2026-10-20
---

# A transient profile-proof timeout does not park the Delivery Queue — Technical Spec

## Executive Summary

Three slices extend existing seams; no package, command or flag is new
(ADR-0298).

1. The setup deadline of each proof attempt becomes
   `defaults.profile_proof_setup_timeout` (default `30s`, greater than zero
   and at most `10m`). `internal/config` reads it, `agent.ProbeRequest`
   carries it, and the agent's setup context uses it.
2. The CLI's proof retry, shared by every profile-readiness consumer, makes
   three attempts in all, waiting 5 and then 10 seconds through an injected
   sleeper, and its final error names the attempts, the deadline and the key.
3. The cleanup of a disposable Agent Session that acpx reports absent returns
   without a warning.

The trade-off accepted: a tuple that can never open a Session now fails after
about two minutes with the default deadline, instead of about seventy
seconds, in exchange for not parking on a transient slow setup. The backoff
is deterministic, so a test proves it with an injected sleeper and no real
wait.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the key and the messages are fixed strings. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; tests use injected Agent runners, an
  injected sleeper and the compiled fake acpx. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Build Order follows recorded
  decisions. ADR-0298 governs Build Order 1-3. ADR-0050: "Configured
  fallbacks activate after notification". ADR-0114: "Opening an Agent Session
  is not Agent work". ADR-0140: "An Agent Selection is proven as the exact
  advertised tuple". ADR-0107: "Profile readiness covers every configured
  Agent Work Category". ADR-0184: "A TechSpec states a command surface as a
  transcript". ADR-0259, ADR-0271, ADR-0285, ADR-0250, ADR-0189, ADR-0256,
  ADR-0280, ADR-0104, ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named
  with its title in `_prd.md` → Project Constraints. Set no obligation here, because this Spec changes none of their rules: ADR-0069, ADR-0092, ADR-0120, ADR-0165, ADR-0167, ADR-0169, ADR-0180, ADR-0181, ADR-0196, ADR-0197, ADR-0217, ADR-0229, ADR-0233, ADR-0237, ADR-0238, ADR-0242, ADR-0245, ADR-0247, ADR-0248, ADR-0251, ADR-0254, ADR-0257, ADR-0258, ADR-0260, ADR-0299, ADR-0261, ADR-0262, ADR-0269, ADR-0272, ADR-0274, ADR-0277, ADR-0278, ADR-0279, ADR-0281, ADR-0282, ADR-0286, ADR-0290 and ADR-0295, each named with its title in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AK approval of 2026-10-20 ("Sim, mesmas regras do AJ (Recomendado)", which
  keeps "Mexer em Governed Paths declarados") and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `docs/references/coverage-record.json`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `9448aa64`.

- The retry (`internal/cli/profiles_validate.go`).
  - `proveProfileSelectionWithRetry` (`:243-253`) calls the proof, and when
    `profileProofSetupTimedOut` (`:255-257`: `errors.Is(err,
    context.DeadlineExceeded) && ctx.Err() == nil`) holds, calls it once more
    at once, with no wait. A second timeout returns `profileProofTimeoutError`,
    whose text is `profile proof setup timed out twice: <err>` (`:74-90`).
  - `profileProofClassification` (`:295-307`) classifies it `temporary`;
    `profileProofNextAction` (`:311-314`) returns `rerun the command when load
    drops; the configured profile was not shown to be wrong`.
  - `proveProfileSelectionsWithOptions` (`:177-231`) proves each distinct
    tuple in order and builds `agent.ProbeRequest{Runtime, WorkDir}`. Its
    consumers are `profiles validate` (`:137`), the operational Preflight of
    `implement` (`internal/cli/profile_preflight.go:62`), the pre-PR review
    (`internal/cli/review.go:1405`), Doctor (`internal/cli/doctor.go:46`),
    Setup (`internal/cli/setup.go:313`) and `profiles configure`
    (`internal/cli/profiles_configure.go:90`).
  - With the built-in profiles and an injected runner whose preferred backend
    proof always returns `open disposable session: context deadline
    exceeded`, `profiles validate --category backend` exits `2` with standard
    output empty and standard error `roundfix: profiles failed: profile proof
    failed for runtime "codex", model "gpt-6.1-sol", reasoning_effort "high";
    affected categories: backend preferred source=built-in; classification:
    temporary; adapter error: profile proof setup timed out twice: open
    disposable session: context deadline exceeded; next: rerun the command
    when load drops; the configured profile was not shown to be wrong;
    fallback: Fallback Chains activate only after Run creation (ADR-0050);
    Preflight proves every configured tuple and substitutes none`, then `Run
    'roundfix profiles --help' for usage.` (probe run in the package's test
    harness, then removed).
  - `internal/cli/profile_proof_retry_test.go` pins two attempts:
    `TestProfileProofRetriesATimedOutSetupOnce` and
    `TestProfileProofReportsASecondTimeoutAsTemporary` (`attempts = 2`), and
    `TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce`, which
    passes `context.Background()` and so the default dependencies.
- The deadline (`internal/agent`).
  - `acpxPreflightSetupTimeout = 30 * time.Second` (`acpx_runner.go:56`)
    bounds `ProveExactSelection`'s setup (`selection_assignment.go:117`) and
    the Fallback Chain probe (`fallback.go:63`). `ProbeRequest`
    (`agent.go:244-247`) has no deadline field, and `ACPXRunner.Probe`
    (`acpx_runner.go:574-582`) rebuilds the request from `Runtime` and
    `WorkDir`.
  - `TestACPXProbeSelectionSetupUsesBoundedContext`
    (`acpx_runner_test.go:1552`) proves the setup context has a deadline
    through `deadlineRecordingCodexProbe`, a codex hygiene probe the runner
    already takes (`codexSpawn`).
- The config (`internal/config/config.go`). `defaultsOverlay` (`:315-321`)
  has no such key; `roundfix profiles show` with
  `defaults: {profile_proof_setup_timeout: 90s}` in Project Config exits `2`
  with `field profile_proof_setup_timeout not found in type
  config.defaultsOverlay`. Durations decode through `durationValue`
  (`:268-288`), and `Validate` refuses non-positive durations with `<key>
  must be greater than 0` (for example `:1141-1142`, `worktree.bootstrap_timeout`).
- The cleanup. `closeDisposableSession` (`acpx_runner.go:1136-1148`) logs
  `disposable Agent Session cleanup found no Session to close: close
  disposable Agent Session "<name>": <acpx error>; recovery: rerun Agent
  Selection readiness after the Session can be closed` and returns `nil`
  when `IsAgentSessionAbsent` holds. `TestMissingSessionIsRecognisedFromBothExitShapes`
  (`acpx_runner_test.go:1353`) pins exactly one such warning for a proof
  whose first `sessions ensure` failed, so the Session never opened. The same
  shape follows a setup deadline that stops `sessions ensure` before acpx
  records the Session, which is the 22:57 line of the adopted entry.
- The guides. `docs/user-guide/commands/profiles.md:100-102` and the Roundfix
  `runtime` reference (`.agents/skills/roundfix/references/runtime.md:33-35`)
  say a timed-out proof "is retried once"; neither the configuration guide nor
  the `profiles` reference names a setup deadline.

## System Architecture

`internal/config` gains the key. `internal/agent` gains the request field and
uses it in `ProveExactSelection`; `ACPXRunner.Probe` forwards it.
`internal/cli` passes the configured value from
`proveProfileSelectionsWithOptions`, so every consumer above gets it, and
`proveProfileSelectionWithRetry` gains the attempts and the backoff through a
new `commandDependencies` sleeper. The cleanup change stays inside
`closeDisposableSession`, which the profile proof, the Fallback Chain probe
and sealed prompts share.

## Implementation Design

### Interfaces

```go
// internal/config
const DefaultProfileProofSetupTimeout = 30 * time.Second
const MaxProfileProofSetupTimeout = 10 * time.Minute
type Defaults struct { /* existing fields */ ProfileProofSetupTimeout time.Duration }

// internal/agent: acpxPreflightSetupTimeout = config.DefaultProfileProofSetupTimeout
type ProbeRequest struct {
	Runtime      RuntimeSpec
	WorkDir      string
	SetupTimeout time.Duration // non-positive uses acpxPreflightSetupTimeout
}

// internal/cli
var profileProofRetryWaits = []time.Duration{5 * time.Second, 10 * time.Second}
// commandDependencies.profileProofSleep waits or returns ctx.Err() first.
type profileProofTimeoutError struct { Attempts int; Timeout time.Duration; Err error }
```

### Data Models

One config key, `defaults.profile_proof_setup_timeout`: a Go duration string,
default `30s`, accepted in User Config and Project Config, Project Config
winning. No record, schema or file changes shape.

### API Contracts

1. API Contract: the setup deadline.
   - `Validate` refuses a value not greater than zero or greater than `10m`
     with exactly `defaults.profile_proof_setup_timeout must be greater than
     0 and at most 10m`. A value that is not a duration fails as every
     duration key does today (`invalid duration "<raw>"`).
   - `proveProfileSelectionsWithOptions` sets `SetupTimeout` to
     `config.Defaults.ProfileProofSetupTimeout` in every `ProbeRequest` it
     builds. `ACPXRunner.Probe` forwards `SetupTimeout` to
     `ProveExactSelection`, whose setup context is `context.WithTimeout(ctx,
     <SetupTimeout, or 30s when not positive>)`.
   - Data source: User Config and Project Config. Destination: the setup
     context of each proof attempt only; the cleanup deadline, the Fallback
     Chain probe and Run sessions keep their values. A cap cancels: an
     attempt whose deadline expires has its acpx command stopped and its
     disposable Session closed, as today.
2. API Contract: the attempts.
   - `proveProfileSelectionWithRetry` makes at most three attempts. After an
     attempt fails with `profileProofSetupTimedOut`, it waits
     `profileProofRetryWaits[n]` (5 s before the second attempt, 10 s before
     the third) through `profileProofSleep`, then tries again. Any other
     error returns at once, unchanged.
   - A sleep that returns an error, which happens when the command's context
     is cancelled during the wait, returns that error at once: no further
     attempt, and the existing unclassified cancellation report.
   - The third timeout returns `profileProofTimeoutError{Attempts: 3,
     Timeout: <deadline>, Err: <last error>}`, whose text is `profile proof
     setup timed out on 3 attempts of <deadline> each: <last error>`, with
     `<deadline>` printed by `time.Duration.String`. It stays `temporary`;
     the next action becomes exactly `rerun the command when load drops, or
     raise defaults.profile_proof_setup_timeout; the configured profile was
     not shown to be wrong`.
   - The default `profileProofSleep` waits on a timer and returns
     `ctx.Err()` when the context ends first; it starts no goroutine.
   - With the setting of Current behavior, `profiles validate --category
     backend` then exits `2`, standard output empty, standard error exactly:
     `roundfix: profiles failed: profile proof failed for runtime "codex",
     model "gpt-6.1-sol", reasoning_effort "high"; affected categories:
     backend preferred source=built-in; classification: temporary; adapter
     error: profile proof setup timed out on 3 attempts of 30s each: open
     disposable session: context deadline exceeded; next: rerun the command
     when load drops, or raise defaults.profile_proof_setup_timeout; the
     configured profile was not shown to be wrong; fallback: Fallback Chains
     activate only after Run creation (ADR-0050); Preflight proves every
     configured tuple and substitutes none`, then `Run 'roundfix profiles
     --help' for usage.` A retry that passes writes nothing.
3. API Contract: cleanup of an absent Session. When `CloseSession` fails and
   `IsAgentSessionAbsent` holds, `closeDisposableSession` returns `nil` and
   writes no warning. Every other close failure keeps its
   `AgentSessionCleanupError`, joined with any setup error as today.
4. API Contract: everything else stays. Every configured tuple is proven
   (ADR-0050); the classification of other failures, the JSON schema
   `roundfix/profiles-validate/v1`, exit codes, the fallback boundary text,
   the Park Classes and `roundfix init`'s template keep their behavior.

### Surface Transcripts

The transcript runs in a temporary Git repository whose `.roundfixrc.yml`
holds `defaults:` with `profile_proof_setup_timeout: 0s`, with `HOME` set to
a temporary directory.

1. Surface Transcript: a setup deadline outside the range is refused.

   ```transcript
   $ roundfix profiles show
   stdout:
   stderr:
   roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m
   Run 'roundfix profiles --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. Every profile-readiness consumer proves each distinct configured tuple, preferred and fallback.
2. A proof attempt's setup deadline is defaults.profile_proof_setup_timeout, 30s unless configured.
3. Only a setup timeout under a live command context is retried, at most twice, after 5 s and then 10 s.
4. A cancellation during a backoff wait starts no further attempt.
5. A disposable Session acpx reports absent at cleanup produces no warning and no cleanup error.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 2, Success Metric 1, Invariants 3-4 → `proveProfileSelectionWithRetry` and `profileProofSleep` (API Contract 2)
- Goal 2, Core Feature 1, Success Metric 3, Invariant 2 → the config key, `ProbeRequest.SetupTimeout` and the setup context (API Contract 1), Surface Transcript 1
- Goal 3, Success Metric 2 → `profileProofTimeoutError`, its classification and next action (API Contract 2)
- Goal 4, Core Feature 2, Invariant 1 → `proveProfileSelectionsWithOptions`, unchanged in what it proves (API Contract 4)
- Goal 5, Core Feature 3, Success Metric 4, Invariant 5 → `closeDisposableSession` (API Contract 3)
- Core Feature 4 → the key in the configuration guide and the Roundfix `profiles` reference (Build Order 1); the attempts in those, the `profiles` guide and the Roundfix `runtime` reference, and **Proof Setup Deadline** in `CONTEXT.md` (Build Order 4)

## Absence Evidence

1. The config key — command: `grep -rn 'profile_proof_setup_timeout' internal docs/user-guide .agents/skills CONTEXT.md`; result: no match at `9448aa64`.
2. `ProbeRequest.SetupTimeout` — command: `grep -rn 'SetupTimeout' internal/agent/agent.go internal/cli --include='*.go'`; result: no match at `9448aa64`; the only setup deadline is the constant `acpxPreflightSetupTimeout` in `internal/agent`.
3. The backoff — command: `grep -n 'time.Sleep\|NewTimer\|time.After' internal/cli/profiles_validate.go internal/cli/profile_preflight.go`; result: no match at `9448aa64`; the second attempt starts at once.
4. **Proof Setup Deadline** — command: `grep -n 'Proof Setup Deadline' CONTEXT.md`; result: no match at `9448aa64`.

## ADR Alternatives

- ADR-0298: this Spec's decision; it rejects proving only the Preferred Selection, raising the built-in deadline instead of retrying, retrying until the command ends, randomized backoff and a lower-severity cleanup log.
- ADR-0050: no rejected alternative is taken; Preflight still proves every configured tuple and substitutes none, and Fallback Chains still activate only after Run creation.
- ADR-0114: no rejected alternative is taken; opening a disposable Session stays selection, and a retried attempt opens a new disposable Session.
- ADR-0140: no rejected alternative is taken; each attempt is the same exact-tuple proof, and only a setup timeout repeats it.
- ADR-0107: no rejected alternative is taken; every configured Agent Work Category's tuples are still proven.
- ADR-0184: no rejected alternative is taken; the changed config surface is stated as Surface Transcript 1, and the timeout text is fixed byte for byte in API Contract 2.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; the backoff is proven with an injected sleeper and the deadline by reading the setup context's deadline, so no test waits on real time.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows.
- ADR-0250: no rejected alternative is taken; the coverage record is re-recorded with its own command after one test is renamed.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edit raises its recorded version.
- ADR-0256: no rejected alternative is taken; the covering skill changes with the `profiles` surface.
- ADR-0280: no rejected alternative is taken; the guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on the adopted measurement and published sources.
- ADR-0240: no rejected alternative is taken; no row needs the network or a Pull Request beyond the standard pre-PR row.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0069: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0180: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0181: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0217: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0238: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0242: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0245: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0257: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0299: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0274: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

ACP adapters through acpx, only as today; no test starts one. The tests use
the injected runner of `internal/cli` and the compiled fake acpx of
`internal/agent`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and compiles for Windows;
no new test waits on real time.

- `internal/config/profile_proof_setup_timeout_test.go` (new):
  `TestProfileProofSetupTimeoutConfig` with subtests: the built-in value is
  `DefaultProfileProofSetupTimeout`; User Config `45s` is read; Project Config
  `90s` overrides User Config `45s`; `0s`, `-1s` and `11m` refuse with API
  Contract 1's message; `soon` refuses with `invalid duration`.
- `internal/agent/proof_setup_deadline_test.go` (new):
  `TestExactSelectionProofSetupUsesTheRequestedDeadline`, reusing the fake
  acpx harness and a codex hygiene probe that records the time left until the
  setup context's deadline, as `deadlineRecordingCodexProbe` records its
  presence: a request with `SetupTimeout` of 2 h leaves more than 1 h on every
  recorded context, through `ProveExactSelection` and through `Probe`; a
  request without one leaves more than zero and at most
  `acpxPreflightSetupTimeout`.
- `internal/cli/profile_proof_setup_deadline_test.go` (new):
  `TestProfileProofPassesTheConfiguredSetupDeadline` (Project Config `90s`
  reaches every tuple's `ProbeRequest`; without it, the built-in value) and
  `TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript`, asserting
  Surface Transcript 1 whole.
- `internal/cli/profile_proof_retry_test.go`: every test injects a recording
  `profileProofSleep`.
  - declared changes: `TestProfileProofRetriesATimedOutSetupOnce` also
    asserts one wait of `profileProofRetryWaits[0]`;
    `TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce` runs with
    a context carrying the test dependencies and asserts the same wait;
    `TestProfileProofReportsASecondTimeoutAsTemporary` is renamed
    `TestProfileProofReportsAThirdTimeoutAsTemporary`, asserts three attempts,
    both waits in order, the next action of API Contract 2 and, in text mode,
    the exact standard error of API Contract 2;
  - new: `TestProfileProofPassesOnTheThirdAttempt` and
    `TestProfileProofStopsWhenCancelledDuringTheBackoff`.
  - The rename re-records `docs/references/coverage-record.json` with `go
    test ./internal/spec -run '^TestCoverageEquivalence$'
    -update-coverage-record -count=1`.
- `internal/agent/acpx_runner_test.go`: declared change,
  `TestMissingSessionIsRecognisedFromBothExitShapes` asserts no warning for
  both exit shapes, keeping its name; new
  `TestTimedOutProofSetupClosesAnAbsentSessionSilently`: `sessions ensure`
  blocks, the test-deadline context expires, `sessions close` exits `4`, and
  the proof returns `context.DeadlineExceeded` with no
  `AgentSessionCleanupError` and no warning.

## Build Order

1. The setup deadline: the config key, `ProbeRequest.SetupTimeout`, its use
   in `ProveExactSelection` and `Probe`, the CLI pass-through, Surface
   Transcript 1, their tests, the key in the configuration guide and the
   Roundfix `profiles` reference, the skill version and the Behavior Surface
   Record, which the new config key moves (no dependency).
2. The attempts: three attempts, the backoff, the sleeper dependency, the
   final message and next action, their tests and the coverage record
   (depends on: 1, which it shares `internal/cli/profiles_validate.go` with
   and whose deadline the message prints).
3. Silent cleanup of an absent Session and its tests (depends on: 1, which it
   shares `internal/agent/acpx_runner.go` with).
4. Documentation of the attempts and the backoff: the configuration guide,
   the `profiles` guide, the Roundfix `profiles` and `runtime` references,
   **Proof Setup Deadline** in `CONTEXT.md`, the skill version and the
   Behavior Surface Record (depends on: 1, which wrote the key into the same
   guide and reference, and 2, whose behavior it describes).
5. QA gate (depends on: 1, 2, 3, 4).

No other Spec is a prerequisite.

## Risks & Considerations

- A tuple that is really unavailable costs about 50 more seconds per
  readiness run with the default deadline; a misconfigured `10m` costs up to
  about 30 minutes, which is why the range stops there.
- A timeout of the cleanup's own 5-second deadline also wraps
  `context.DeadlineExceeded` and is retried, as today.
- Doctor runs under its own deadline; when that deadline ends the command
  context, no retry starts, as today.

## Glossary

- adds: **Proof Setup Deadline**

## Decisions

- Three attempts with waits of 5 s and 10 s, deterministic. See ADR-0298.
- `defaults.profile_proof_setup_timeout`, default `30s`, range above zero to
  `10m`, User and Project Config. See ADR-0298.
- Silent cleanup of an absent disposable Session. See ADR-0298.
- The Fallback Chain probe after Run creation keeps the built-in deadline.
- The retry stays in `internal/cli`, so every consumer gets it from one place.
