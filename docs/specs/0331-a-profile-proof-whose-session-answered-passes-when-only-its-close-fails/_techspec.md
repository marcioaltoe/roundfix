---
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
prd: _prd.md
created: 2026-10-27
---

# A profile proof whose session answered passes when only its close fails — Technical Spec

## Executive Summary

The defect has three parts in two packages. In `internal/agent`, the two
disposable-session proofs (`ProveExactSelection` and the Fallback Chain
candidate probe) let a close error decide the verdict and give the close a
fixed 5-second deadline. In `internal/cli`, `profileProofSetupTimedOut` treats
any wrapped `context.DeadlineExceeded` as a setup timeout, so a slow close
repeats the whole proof and is reported as a setup that timed out three
times. The fix keeps every seam: the proof decides on its setup, closes
under the attempt's Proof Setup Deadline, records a session it could not
close in a small per-session file under the Roundfix home, and the first
proof of each Agent runner closes the recorded sessions first. A setup
deadline is wrapped with the step that reached it. The CLI counts only a
setup deadline as a timeout and shows the warning the proof carries. The
Baseline semantic analyzer keeps its stop by treating a proof that left its
session open as cleanup unproven.

The trade-off accepted: a new on-disk record and a reaping pass in the proof
path, instead of a background close or a sweep of `acpx sessions list`,
because only the record names exactly what Roundfix failed to close and
survives the process (ADR-0344).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme is minted; the
  record file is named by the existing disposable session name. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; the record stores no environment value. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0140, ADR-0298, ADR-0069 and
  ADR-0344 govern the design; ADR-0180, ADR-0181 and ADR-0217 are recorded
  as not applying, as `_prd.md` → Project Constraints states and ADR
  Alternatives below relates. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Governed Paths the authorization
  record lists, granted by the cycle AR decisions of 2026-10-27; no
  Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails/_authorization.md`;
  bounded files:
  `.agents/skills/roundfix/references/profiles.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-27 at `5e9475bf`; `_prd.md` → Validation records the runs.

1. `ProveExactSelection` (`internal/agent/selection_assignment.go`) runs
   `checkAdapter`, `codexEnvForSession`, `readRuntimeCatalogueWithEvidence`,
   `startSessionSelection`, `applySessionSelection` and `applyFullAccess`
   under one setup context bounded by `ProbeRequest.SetupTimeout` (built-in
   `acpxPreflightSetupTimeout`, 30 s), then `closeDisposableSession` under
   `disposableCleanupTimeout()` (5 s unless the unexported test seam
   `cleanupTimeout` is set). It returns `errors.Join(setupErr, cleanupErr)`,
   `setupErr` or `cleanupErr`.
2. `probeFallbackCandidate` (`internal/agent/fallback.go`) closes the same
   way and returns the close error (joined with a probe error) before it
   looks at the candidate's result.
3. `closeDisposableSession` (`internal/agent/acpx_runner.go`) wraps a close
   failure in `*AgentSessionCleanupError` (classification
   `session_cleanup_failed`) and returns nil when acpx reports the session
   absent (ADR-0298).
4. Setup steps return `context.DeadlineExceeded` unwrapped, so the step is
   lost.
5. `proveProfileSelectionWithRetry` (`internal/cli/profiles_validate.go`)
   retries while `profileProofSetupTimedOut` holds, which is
   `errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil`.
6. `Analyzer.proveSelection` (`internal/baselineacp/analyzer.go`) returns the
   proof error, and each stage stops at `cleanupUnproven`.
7. Nothing closes a disposable session after its proof returns.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/agent/selection_assignment.go`, `internal/agent/fallback.go`:
  verdict, close deadline, step wrapping, warning, record and reaping calls.
- `internal/agent/acpx_runner.go`: the close deadline helper and the runner's
  once-only reaping state.
- `internal/agent/unclosed_proof_sessions.go` (new, task_02): the Unclosed
  Proof Session record: its path, its write and the warnings.
- `internal/agent/unclosed_proof_session_reaper.go` (new, task_03): reading
  the records and the reaping pass.
- `internal/baselineacp/analyzer.go`: an unclosed proof is cleanup unproven.
- `internal/cli/profiles_validate.go`: the setup-timeout rule and the
  `warning` report field and line.

## Implementation Design

### Interfaces

```go
// internal/agent/selection_assignment.go
type SelectionProof struct {
	// ...existing fields...
	Unclosed *UnclosedProofSession // set when the proof passed but its close failed
}

// internal/agent/unclosed_proof_sessions.go
type UnclosedProofSession struct {
	Session  string
	WorkDir  string
	Deadline time.Duration
	Warning  string // the line the runner printed
}
func (runner *ACPXRunner) recordUnclosedProofSession(runtime RuntimeSpec, session, workDir string, deadline time.Duration, cause error) *UnclosedProofSession

// internal/agent/unclosed_proof_session_reaper.go
func (runner *ACPXRunner) reapUnclosedProofSessions(ctx context.Context, deadline time.Duration)

// internal/cli/profiles_validate.go
type profileProofReport struct {
	// ...existing fields...
	Warning string `json:"warning,omitempty"`
}
```

### Data Models

The record `<HOME>/.roundfix/unclosed-proof-sessions/<session>.json`, where
`HOME` is read from the runner's environment (`ACPXRunner.Environment`, or
the process environment when it is nil). Directory mode `0700`, file mode
`0600`, written to a temporary file in the same directory and renamed:

```json
{"schema":"roundfix/unclosed-proof-session/v1","session":"roundfix-preflight-…",
 "work_dir":"/abs/dir","runtime":{"id":"codex","protocol":"acp","command":"",
 "openrouter_key_variable":""},"recorded_at":"2026-10-27T15:30:00Z","cause":"…"}
```

`openrouter_key_variable` is the variable's name the runtime already carries,
never its value; no other environment data is stored.

### API Contracts

1. API Contract: the verdict is the setup's. In `ProveExactSelection`, when
   every setup step succeeded and the close returns any error other than a
   proven absence, the proof is returned with a nil error and `Unclosed` set.
   When a setup step failed, its error is returned as it is today minus the
   join: no `*AgentSessionCleanupError` is in its chain. In both cases a
   failed close prints exactly one warning through the runner's warning
   output and writes one record (API Contract 4). In
   `probeFallbackCandidate` a failed close never changes the candidate's
   result (functional, non-functional or error); it prints the warning and
   writes the record.
2. API Contract: the close deadline. `ProveExactSelection` closes under the
   attempt's resolved setup deadline (`ProbeRequest.SetupTimeout` when
   positive, else `acpxPreflightSetupTimeout`); `probeFallbackCandidate`
   closes under `acpxPreflightSetupTimeout`. The close context still derives
   from `context.WithoutCancel(ctx)`. The unexported `cleanupTimeout` seam
   still overrides when positive. `RunSealedPrompt`, `EndSession` and Run
   session cleanup keep their deadlines.
3. API Contract: named steps. A setup step error that satisfies
   `errors.Is(err, context.DeadlineExceeded)` is wrapped as
   `profile proof step "<step>": %w`, with `<step>` one of `check adapter`,
   `start runtime`, `read runtime catalogue`, `start session selection`,
   `apply session selection`, `apply full access` in `ProveExactSelection`,
   and `start runtime`, `apply session selection`, `apply reasoning effort`
   in the candidate probe. Any other error keeps its bytes.
4. API Contract: the record and its warnings. A failed close writes the
   Data Models record and prints `disposable Agent Session "<session>" was
   not closed within <deadline> after its profile proof: <cause>; recorded in
   <record path> for the next profile proof to close`. When the record cannot
   be written (no `HOME`, a write error), it prints `disposable Agent Session
   "<session>" was not closed within <deadline> after its profile proof:
   <cause>; it could not be recorded (<error>); close it with: acpx <close
   arguments>`, with the arguments `acpxCloseArgs` builds. `<deadline>` is
   the close deadline in Go duration form; `<cause>` is the close error's
   text.
5. API Contract: reaping. Before it opens its disposable session, the first
   `ProveExactSelection` or candidate probe of an `ACPXRunner` reads every
   record, oldest `recorded_at` first (session name breaks ties), and closes
   each with `sessions close` under that proof's close deadline, building
   the command from the stored runtime fields and working directory. A
   record whose close succeeds or whose session acpx reports absent
   (`IsAgentSessionAbsent`) is deleted. The first close that fails prints
   `recorded disposable Agent Session "<session>" (<record path>) was not
   closed within <deadline>: <cause>; the next profile proof retries it` and
   ends the pass, keeping that record and every later one. A record that
   cannot be parsed prints `unreadable Unclosed Proof Session record <path>:
   <error>; remove it after closing its session by hand` and is kept. A
   cancelled context ends the pass. Reaping never returns an error to the
   proof, never changes its verdict, and runs at most once per runner.
6. API Contract: the CLI's timeout rule. `profileProofSetupTimedOut` holds
   only when a `context.DeadlineExceeded` is in the error chain outside
   every `*agent.AgentSessionCleanupError` and the command's context is
   live. A proof error whose only deadline sits inside a close error is
   returned after one attempt with its own classification.
7. API Contract: the report. `profileProofReport.Warning` carries
   `SelectionProof.Unclosed.Warning`; `roundfix/profiles-validate/v1` JSON
   gains the optional `warning` field; the text success output prints
   `   warning: <text>` after that proof's reference lines. Exit codes and
   every other field keep their bytes.
8. API Contract: the analyzer. `Analyzer.proveSelection` returns an
   `*agent.AgentSessionCleanupError` naming `Unclosed.Session` (with the
   warning as its cause) when the proof carries `Unclosed`, so every stage's
   `cleanupUnproven` stop behaves as it did when the close error failed the
   proof.

### Surface Transcripts

None. No command gains a flag, an argument or an exit code; the new lines
are fixed in API Contracts 4, 5 and 7 and asserted by the named tests of
Testing Approach.

### Invariants

```text
1. A setup step that fails or times out fails the proof with the same error type and classification as before; only the close error leaves its chain.
2. A close outcome never changes a proof's verdict or a fallback candidate's result.
3. Every failed close writes one record or prints the acpx close command; no failed close is silent.
4. A record leaves the directory only after acpx closed its session or reported it absent.
5. The record holds no environment value, no prompt and no token.
6. Every configured tuple is still proven before a Run, and the three attempts and their backoff are unchanged.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `ProveExactSelection`,
  `probeFallbackCandidate` (API Contract 1).
- Goal 2, Core Feature 1, Success Metric 2 → the same, Invariant 1.
- Core Feature 2 → the close deadline helper (API Contract 2).
- Goal 3, Core Feature 3, Success Metric 3 →
  `internal/agent/unclosed_proof_sessions.go` and
  `internal/agent/unclosed_proof_session_reaper.go` (API Contracts 4–5).
- Goal 4, Core Feature 4, Success Metric 4 → step wrapping and
  `profileProofSetupTimedOut` (API Contracts 3, 6).
- Core Feature 5, Success Metric 5 → `printProfilesValidateSuccess`,
  `profileProofReport`, `Analyzer.proveSelection` (API Contracts 7–8).
- Success Metric 6 → declared unreachable; the guides task_01 writes tell the
  operator what the warning means.

Deviations from the PRD: None.

## Absence Evidence

1. Unclosed Proof Session record and reaper — command: `rg -n 'unclosed|roundfix-preflight|acpxPreflightSessionPrefix' internal -g '!*_test.go'`; result: only the session-name generator in `internal/agent/acpx_runner.go`; nothing reads or closes a disposable session after its proof.
2. Proof step naming — command: `rg -n 'profile proof step' internal`; result: no match.
3. Close deadline tied to the setup deadline — command: `rg -n 'disposableCleanupTimeout\(\)' internal/agent -g '!*_test.go'`; result: `selection_assignment.go`, `fallback.go` and `sealed.go` call it, and it returns the 5-second constant unless the test seam is set.
4. `warning` in the profiles validate report — command: `rg -n 'json:"warning' internal/cli/profiles_validate.go`; result: no match.

## ADR Alternatives

- ADR-0069: followed; the analyzer keeps falling back to manual classification when cleanup is unproven, which ADR-0344 carries over to an unclosed proof session.
- ADR-0180: untouched; built-in selections and the Recommended Profile are not read or changed.
- ADR-0181: untouched; no configuration is compared with the recommendation.
- ADR-0217: untouched; a Cursor selection's model value and login are not changed, and its proof keeps refusing before any session opens.
- ADR-0140: followed and narrowed; the tuple is still proven from advertised state, the static catalog and config inspection it rejected stay rejected, and only the close leaves the verdict.
- ADR-0298: refined; the attempts, backoff and setup key stay, the close now shares the setup deadline, and its "about two minutes" consequence grows as ADR-0344 states.
- ADR-0344: implemented as recorded, including its rejection of a background close, a `sessions list` sweep and a load-average refusal.

## Integration Points

acpx `sessions close`, as the existing close already runs it. No new
external call; tests use the compiled fake acpx of `internal/agent`.

## Testing Approach

All tests are hermetic: the fake acpx harness (`newFakeACPXHarness`, whose
environment carries a temporary `HOME`), fake provers in `internal/cli`, and
the fake runtime of `internal/baselineacp`. No test starts acpx, an adapter or
the network.

- task_02, `internal/agent/proof_close_test.go`:
  `TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs` (passing setup,
  `sessions close` blocked, a short `cleanupTimeout` that the test asserts:
  nil error, `Unclosed` set, one warning, one record),
  `TestTheProofCloseDeadlineFollowsTheSetupDeadline` (`cleanupTimeout` zero,
  `SetupTimeout` two hours, close exit `2`: the warning names `2h0m0s`),
  `TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails` (rejected
  `set reasoning_effort`, close exit `2`: `*SelectionRejectedError`, no
  `*AgentSessionCleanupError`, a record),
  `TestATimedOutProofStepIsNamed` (first `sessions show` blocked, the parent
  context expired by the test: `profile proof step "read runtime
  catalogue"`, `errors.Is` deadline),
  `TestAFallbackCandidateKeepsItsResultWhenItsCloseFails` (close exit `2`:
  `ok` true with the candidate, a record),
  `TestAnUnclosedProofSessionRecordHoldsNoEnvironmentValue` (a sentinel value
  built at run time in the runner environment never appears in the record;
  modes `0700`/`0600`).
  `internal/baselineacp/analyzer_test.go`:
  `TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis`.
  `internal/testfixture/real_time_budget_test.go`: a ceiling of `1` for
  `internal/agent/proof_close_test.go`.
  Declared breaks: `TestProveExactSelectionCleanupJoinedFailure`,
  `TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose`,
  `TestACPXProbeCleanupFailureJoinsSelectionError` and
  `TestACPXProbeFallbackReportsCleanupInfrastructureError` in
  `internal/agent/acpx_runner_test.go` assert the old join; they are removed,
  and the new tests above take their cases under the new contract.
- task_03, `internal/agent/unclosed_proof_sessions_test.go`:
  `TestTheNextProofClosesARecordedSessionFirst`,
  `TestARecordedSessionAcpxReportsAbsentIsForgotten`,
  `TestAFailedReapKeepsTheRecordAndStops`,
  `TestTheReapRunsOncePerRunner`,
  `TestAnUnreadableRecordIsReportedAndKept`;
  `internal/cli/profile_proof_retry_test.go`:
  `TestACloseDeadlineIsNotRetriedAsASetupTimeout`;
  `internal/cli/profile_proof_close_test.go`:
  `TestProfilesValidateReportsAnUnclosedProofSession` (text and JSON).
- Existing suites that consume the changed contracts run in the same
  Verification commands: the `TestProveExactSelection*`, `TestACPXProbe*`,
  `TestExactSelectionProofSetup*` and analyzer tests, and the CLI profile
  proof retry and validate tests.

## Build Order

1. Verdict, close deadline, named steps, warning, record and the analyzer
   stop (task_02).
2. Reaping, the CLI timeout rule and the report (task_03) (depends on: 1 —
   both edit `selection_assignment.go` and `fallback.go`, and the report
   reads `SelectionProof.Unclosed`).
3. Guides, skill references, glossary and Behavior Surface record (task_01)
   (depends on: 1, 2 for the behavior it describes; the graph runs it first
   so each code Task names its guide through that dependency).
4. QA gate (task_04) (depends on: 2, 3).

## Risks & Considerations

- A wrong reap could close a live session: the reaper closes only names
  Roundfix recorded after its own close failed, and the disposable name
  carries random entropy, so no other owner uses it.
- A record directory that grows under persistent acpx failure: each pass
  stops at the first failure, so the cost per runner is one close deadline;
  the warnings name every record path.
- Two processes reaping one record: the second close sees an absent session
  and deletes the record; a missing file on delete is not an error.
- `internal/agent/selection_assignment.go` is shared with live sessions;
  the verdict change applies only to the disposable proof path.

## Glossary

- adds: **Unclosed Proof Session**
- changes: **Exact Agent Selection Proof**
- changes: **Proof Setup Deadline**

## Decisions

- Setup decides the verdict; a failed close is a warning and a record. See
  ADR-0344.
- The close shares the setup deadline. See ADR-0344.
- A per-session record under the Roundfix home and a once-per-runner reaping
  pass, over a background close or a `sessions list` sweep. See ADR-0344.
- The analyzer's stop is kept by converting an unclosed proof into the
  cleanup error it already handles, so no stage changes.
