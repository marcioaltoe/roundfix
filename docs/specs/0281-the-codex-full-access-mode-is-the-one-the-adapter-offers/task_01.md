---
task: task_01
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
status: completed
type: backend
complexity: high
---

# Task 01: A full-access request selects the mode the Agent Session advertises

## Overview

On 2026-10-18 codex-acp 2.2.2 refused the fixed Codex mode `full-access`
with ACP -32602 and accepted `agent-full-access` (`_prd.md` → Sources 1).
This Task makes `applyFullAccess` send the first full-access candidate the
Agent Session advertises in its `mode` config option. It sends the Codex
sandbox preset when the session advertises `sandbox_mode`, and reports a
degraded effective access policy when no candidate is advertised. It is
verifiable on its own through the fake acpx runner, in the disposable proof
and in a Run. Because the change alters what `profiles validate`, `doctor`
and the Run events report, the Roundfix `profiles` reference describes it in
the same Task.

## Requirements

1. MUST answer the Secondbrain Inbox Entry of 2026-10-18,
   `references/2026-10-18-agent-full-access-pede-um-modo-que-o-codex-acp-recusa.md`
   (`_prd.md` → Sources 1), by implementing API Contracts 1 to 6 of
   `_techspec.md` (ADR-0293).
2. MUST add `RuntimeSpec.FullAccessModes`. `RuntimeFor` sets it to
   `agent-full-access`, `full-access` for `codex` and to `bypassPermissions`
   for `claude`, and to nothing for any other runtime or a command override
   (API Contract 1). `SupportedFullAccessMode`, and with `EnableFullAccess`
   `FullAccessMode`, hold the first candidate. A `RuntimeSpec` without
   `FullAccessModes` uses the name `AccessModeFor` returns as its single
   candidate.
3. MUST add `SelectionCapabilities.Modes`, set by `ParseSessionConfigOptions`
   to the projected select option whose identifier is `mode`, else the first
   whose category is `mode`, else nil (API Contract 2). Every existing issue
   rule, `Options` and the snapshot parser stay as they are.
4. MUST make `applyFullAccess` take the `SelectionCapabilities` its caller
   observed: `ProveExactSelection` passes the ones from its
   `startSessionSelection`, and `ensureSession` passes the ones it observed.
   No new `sessions show` call is added. It MUST:
   - choose the mode with `chooseFullAccessMode`: the first candidate, in
     candidate order, that the advertised values contain, compared exactly
     after trimming space (API Contract 3);
   - send `set-mode` with that mode only, and keep today's
     `access_policy_rejected` error when the adapter refuses it;
   - on a `codex` runtime, after a successful `set-mode`, send `set
     sandbox_mode danger-full-access` only when `capabilities.Options` holds
     an option with the identifier `sandbox_mode`, keeping today's handling of
     its failure; without that option send no preset and return `full-access`
     (API Contract 4).
5. MUST report a session that advertises no candidate as API Contract 5
   states.
   - Send no `set-mode` and no preset.
   - Return the effective access policy `full-access (degraded: no advertised
     mode grants full access; advertised modes: <values joined by ", ">)`, or
     `... advertised modes: none)` when the session advertises no mode
     values.
   - Publish the status `full_access_mode_not_advertised`
     (`AccessModeNotAdvertisedStatus`) and warn `full-access mode not
     advertised for Agent Session <name>; advertised modes: <list>`.
   - Let the Run continue.

   The advertised values come from the adapter's own `acpx sessions show`
   output. The policy, status and warning are read by the operator through
   `profiles validate`, `doctor` and the Run event stream, and carry mode
   identifiers only.
6. MUST add these tests in the new `internal/agent/full_access_mode_test.go`.
   Each calls `t.Parallel()` first, uses `newFakeACPXHarness`, overrides
   `sessions show` through the fake's stdout map, and compiles for Windows.
   No test starts acpx or an adapter.
   - `TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate`.
   - `TestProofSelectsTheAdvertisedCodexFullAccessMode`: `read-only`,
     `workspace-write`, `agent`, `agent-full-access`, no `sandbox_mode`. The
     proof sends one `set-mode agent-full-access` and no `set sandbox_mode`,
     and records `full-access`.
   - `TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt`: `full-access`
     advertised with a `sandbox_mode` option. A Run sends `set-mode
     full-access`, then `set sandbox_mode danger-full-access`.
   - `TestRunSelectsBypassPermissionsForClaude`: a Claude Run with
     `bypassPermissions` advertised sends `set-mode bypassPermissions` and no
     preset.
   - `TestFullAccessWithoutAnAdvertisedModeIsDegraded`: `read-only`,
     `workspace-write`, `agent` advertised. The proof sends no `set-mode` and
     records the exact degraded text. A Run publishes
     `full_access_mode_not_advertised`, warns with the session name and the
     list, and reaches its prompt. A snapshot without a mode option yields
     `advertised modes: none)`.
7. MUST keep the existing full-access tests' contracts by giving the fake's
   default `sessions show` (`fakeSelectionStateOutput`) and
   `sessionCapabilitySnapshotFixture` a `mode` option. For Codex the values
   are `read-only`, `auto` and `full-access`, with a `sandbox_mode` option. For
   Claude they are `default` and `bypassPermissions`. With that,
   `TestDisposableProofAppliesRequestedAccessMode`,
   `TestProofRecordsADegradedFullAccess`,
   `TestACPXRunAppliesFullAccessSessionSetup`,
   `TestACPXRunFailsSetupWhenFullAccessModeFails` and
   `TestACPXRunWarnsWhenCodexSandboxPresetUnavailable` keep their names and
   their invocation lists. MUST move the Codex name pinned by
   `TestRuntimeForCodexFullAccessOptIn` and
   `TestRuntimeReportsAccessPolicySupport` in
   `internal/agent/agent_test.go`, and by
   `TestRunResolveAgentFullAccessIsExplicitOptIn` in
   `internal/cli/cli_test.go`, from `full-access` to `agent-full-access`, the
   first Codex candidate. MUST sweep the other fake `sessions show` producers
   (the shell and Python fakes in `internal/cli/implement_test.go`) and
   change one only when a test there requests full access through the acpx
   runner.
8. MUST add a `## Full access` section to
   `.agents/skills/roundfix/references/profiles.md`. It names the candidates
   `agent-full-access`, `full-access` and `bypassPermissions`, the advertised
   `mode` option, the sandbox preset rule, the degraded text `no advertised
   mode grants full access`, the status `full_access_mode_not_advertised`, and
   that `profiles validate` and `doctor` print the degraded policy. No text
   inside any `### QA settlement` section may change.
9. MUST run these commands, because this Task edits the Roundfix reference.
   No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix skill version.

## Subtasks

- [ ] Add the candidate lists and the session's mode option.
- [ ] Select the advertised mode, send the sandbox preset by advertisement, and report the degraded case.
- [ ] Add the five tests and keep the pinned tests' invocation lists.
- [ ] Move the pinned Codex name to `agent-full-access`.
- [ ] Describe the selection in the Roundfix `profiles` reference, sync the mirror and record the version.

## Acceptance Criteria

- [ ] A Codex session advertising `agent-full-access` is set to it, and the proof records `full-access`.
- [ ] A Codex session advertising `full-access` and `sandbox_mode` keeps today's mode and preset, and Claude keeps `bypassPermissions`.
- [ ] A session advertising no candidate receives no `set-mode` and is reported degraded with its advertised modes.
- [ ] The Roundfix `profiles` reference describes the selection and the degraded report.

## Consumers

- `agent.RuntimeSpec.FullAccessModes`, `agent.SelectionCapabilities.Modes` and the effective access policy text — consumers: `./internal/agent`, `./internal/cli`

## Context

- instruction: `docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md`
- instruction: `docs/adr/0011-acp-runtimes-run-batches-with-full-access.md`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/selection_capabilities.go`
- creates: `internal/agent/full_access_mode_test.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/agent/selection_capabilities_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/implement_test.go`
- instruction: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/cli/profiles_validate.go`
- instruction: `internal/cli/doctor.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate|TestProofSelectsTheAdvertisedCodexFullAccessMode|TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt|TestRunSelectsBypassPermissionsForClaude|TestFullAccessWithoutAnAdvertisedModeIsDegraded)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate TestProofSelectsTheAdvertisedCodexFullAccessMode TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt TestRunSelectsBypassPermissionsForClaude TestFullAccessWithoutAnAdvertisedModeIsDegraded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/agent || exit 1; go test -count=1 -run 'FullAccess|AccessPolicy|DegradedPolicy|ProfilesValidate' ./internal/cli || exit 1; GOOS=windows go vet ./internal/agent || exit 1; for phrase in 'agent-full-access' 'full_access_mode_not_advertised' 'no advertised mode grants full access'; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/profiles.md | grep -qF -- "$phrase" || { printf 'missing phrase in the profiles reference: %s\n' "$phrase" >&2; exit 1; }; done; grep -qxF -- '## Full access' .agents/skills/roundfix/references/profiles.md || { printf 'missing profiles reference heading\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills` — expected: exit 0. Before this Task the five tests do not exist, so the command fails at its first pass check. Afterwards the new tests, the whole `internal/agent` suite with the pinned full-access tests, the CLI access-policy tests, the Windows compilation of the package's tests, the `profiles` reference phrases, the skill mirror and the recorded Roundfix version pass.

## References

- `_prd.md` → Sources; Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants; Current behavior; Testing Approach; Build Order 1
- ADR-0293; ADR-0011; ADR-0017; ADR-0285; ADR-0125; ADR-0189

## Result

Implemented the Task 01 selection contract. `RuntimeFor` now exposes ordered
full-access candidates (`agent-full-access`, `full-access` for Codex;
`bypassPermissions` for Claude), while legacy `RuntimeSpec` values still use
the mode returned by `AccessModeFor`. Capability parsing projects `Modes`
with identifier precedence over the first category match, preserving the
existing Options and issue rules. Proof and Run setup pass their already
observed capabilities into full-access setup; no observation command was added.

Acceptance evidence from focused checks:

- Codex proof: `TestProofSelectsTheAdvertisedCodexFullAccessMode` passes with
  `agent-full-access` advertised and no sandbox option. It checks one exact
  mode call, no preset, effective policy `full-access`, and disposable cleanup.
  `TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate` covers candidate
  order, trimming, exact comparison, absent values/candidates, and a proof
  advertising both Codex candidates.
- Preset and Claude: `TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt`
  passes for both Codex candidate names and checks mode-before-preset order.
  `TestRunSelectsBypassPermissionsForClaude` checks `bypassPermissions` and no
  preset. The five pinned proof/Run tests pass without changes to their
  invocation expectations, including rejection and unavailable-preset handling.
- Degradation: `TestFullAccessWithoutAnAdvertisedModeIsDegraded` passes for
  unmatched modes and no mode option, checking the exact policy, no mode or
  preset calls even with a sandbox option present, status, warning with session
  and list, cleanup, and continuation to the Run prompt.
- Reference: added `## Full access` to the canonical profiles reference,
  describing candidates, mode discovery, preset conditions, degraded policy,
  status, and profiles/Doctor reporting. `make skills-sync` and the owned-skill
  record command succeeded; the generator raised Roundfix from `0.1.102` to
  `0.1.103` in both copies and recorded its digest. No mirror or record was
  edited by hand. `make baseline-digests` succeeded with `changed:false`.

Focused commands and outcomes (all exit 0 on the final source):

```sh
go test ./internal/agent -run 'TestChooseFullAccessMode|TestProofSelectsTheAdvertised|TestRunAppliesTheSandbox|TestRunSelectsBypass|TestFullAccessWithout|TestDisposableProofAppliesRequestedAccessMode|TestProofRecordsADegradedFullAccess|TestACPXRunAppliesFullAccessSessionSetup|TestACPXRunFailsSetupWhenFullAccessModeFails|TestACPXRunWarnsWhenCodexSandboxPresetUnavailable|TestRuntimeFor|TestRuntimeReportsAccessPolicySupport|TestSelectionCapabilities|TestParseSessionCapabilitySnapshot' -count=1 -v
go test ./internal/cli -run '^(TestRunResolveAgentFullAccessIsExplicitOptIn|TestProfilesValidateTextNamesADegradedPolicy|TestDoctorNamesADegradedPolicy|TestProfilePreflightRefusesUnhonouredAccessPolicy|TestReviewRetriesASelectionFailureBeforeThePromptOnce|TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt|TestReviewRoundTwoContinuesTheRecordedSession|TestReviewRoundTwoContinuesTheRecordedFallbackSelection|TestReviewValidatorDismissesADaemonSettlementAsC2)$' -count=1
GOOS=windows go test -c -o /tmp/roundfix-0281-agent.test.exe ./internal/agent
make skills-sync
go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
make baseline-digests
git diff --check
```

Additional scope necessary for this API change: adding a slice makes
`RuntimeSpec` non-comparable. The focused CLI build exposed five direct
comparisons in `internal/cli/review_convention_validator_test.go`,
`internal/cli/review_selection_retry_test.go`, and
`internal/cli/review_session_test.go`. Replaced only those comparisons with
`reflect.DeepEqual`, preserving their full-runtime equality assertions. These
ordinary test paths are not governed by the tooling allowlist; the focused
review tests listed above exercise every adjusted assertion.

Swept the shell and Python `sessions show` producers in
`internal/cli/implement_test.go`. Its full-access callers use fake proof
runners rather than either acpx producer, so neither producer needed editing.
The new tests use the compiled fake harness only; no adapter, acpx installation,
provider call, or network request was used. Windows test compilation succeeded.

The initial focused compile demonstrated the missing chooser and degraded
status. Early focused checks also exposed test-helper mistakes (the preset
call key omits its value, and proof plus Run have distinct observation counts);
those helpers were corrected and the final focused checks above passed.

Execution used the resolved `/Volumes/Brain-SSD/...` checkout with authorized
escalation because the executor rejected the configured symlinked writable
root before command execution. The only pre-existing modification was this
Task's Daemon-owned status. Task status, checkboxes, Task Graph and other Task
files remain untouched. Declared Verification and repository Verification were
not run; they remain Daemon-owned. No terminal Task verdict is claimed.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/review_convention_validator_test.go`
- `internal/cli/review_selection_retry_test.go`
- `internal/cli/review_session_test.go`

## Carry-forward provenance

- Source Run: `run_20261019T203112Z_ce9489b10bc966cb`
- Source commit: `7be840588754a314038e781dfd2f28991e2d7496`
