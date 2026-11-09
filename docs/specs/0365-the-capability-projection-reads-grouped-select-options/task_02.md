---
task: task_02
spec: 0365-the-capability-projection-reads-grouped-select-options
status: pending
type: test
complexity: medium
---

# Task 02: The Exact Agent Selection Proof passes over the recorded codex-acp 2.6.0 payload

## Overview

`roundfix doctor` and `roundfix profiles validate` failed on codex-acp
2.6.0 inside `ProveExactSelection`, which reads the session's
`config_options` through `acpx sessions show`. Task 01 makes the projection
read grouped options; this Task proves the whole disposable proof over the
two payloads recorded by the operator, through the existing fake `acpx`
harness, including the full-access mode that codex-acp 2.6.0 lists inside
its `sandbox-modes` group (ADR-0293). It answers the operator record of
2026-11-09
([references](references/2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md)).

## Requirements

1. MUST add `internal/agent/selection_proof_grouped_mode_test.go` with
   `TestProofOverTheRecordedCodexACPPayloads`, which calls `t.Parallel()`
   first and runs four subtests, `2.6.0/full-access`, `2.6.0/runtime-default`,
   `2.5.0/full-access` and `2.5.0/runtime-default`. Each subtest:
   - reads `internal/agent/testdata/codex-acp-<version>-config-options.json`
     and wraps it as
     `{"schema":"acpx.session.v1","acpx":{"current_model_id":"gpt-6.1-sol","config_options":<fixture>}}`;
   - builds a fresh `newFakeACPXHarness(t)` whose `fakeACPXStdoutCall`
     answers `sessions show` with that snapshot;
   - uses `fullAccessRuntimeForTest(t, "codex")` for `full-access` and
     `selectedRuntime(RuntimeSpec{ID: "codex", Protocol: ProtocolACP})` for
     `runtime-default`, each with `Model` set to `gpt-6.1-sol` and
     `ReasoningEffort` set to `high`;
   - calls `harness.runner.ProveExactSelection` with that runtime and
     `harness.gitRoot`, and asserts no error, an assignment of model
     `gpt-6.1-sol`, reasoning `high` and encoding `independent`, and
     `assertLastInvocationClosesDisposable`;
   - for `full-access`, asserts `EffectiveAccessPolicy` is
     `AccessPolicyFullAccess` and `assertFullAccessCalls(t, harness, []string{"set-mode agent-full-access"})`;
   - for `runtime-default`, asserts `EffectiveAccessPolicy` is
     `AccessPolicyRuntimeDefault` and `assertFullAccessCalls(t, harness, nil)`.
2. MUST keep passing, unchanged, every test in
   `internal/agent/full_access_mode_test.go` and the tests task_01 names.
3. MUST NOT change any production file, any helper in
   `full_access_mode_test.go` or `acpx_runner_test.go`, the fixtures, or any
   file outside this Task's declarations. If the proof fails over the 2.6.0
   payload, report the projection defect instead of changing production
   code here.

## Subtasks

- [ ] Add the four proof subtests over the recorded payloads.

## Acceptance Criteria

- [ ] The proof of `codex gpt-6.1-sol high` passes over the 2.6.0 payload with and without full access, as it does over 2.5.0.
- [ ] A full-access request on the 2.6.0 payload sends exactly `set-mode agent-full-access`.

## Context

- instruction: `docs/adr/0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md`
- instruction: `docs/adr/0375-a-grouped-select-option-projects-its-groups-values-and-capability-evidence-stays-fail-closed.md`
- instruction: `internal/agent/selection_capabilities.go`
- instruction: `internal/agent/full_access_mode_test.go`
- instruction: `internal/agent/acpx_runner_test.go`
- creates: `internal/agent/selection_proof_grouped_mode_test.go`

## Consumers

None. This Task adds a test and changes no contract.

## Verification

- `grep -q 'func TestProofOverTheRecordedCodexACPPayloads(' internal/agent/selection_proof_grouped_mode_test.go || exit 1; go test -timeout 20m -count=1 -run '^(TestProofOverTheRecordedCodexACPPayloads|TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate|TestProofSelectsTheAdvertisedCodexFullAccessMode|TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt|TestRunSelectsBypassPermissionsForClaude|TestFullAccessWithoutAnAdvertisedModeIsDegraded)$' ./internal/agent` — expected: exit 0. Before this Task the test file does not exist, so the first check fails; afterwards the proof passes over both recorded payloads and the full-access tests pass unchanged.

## References

- `_prd.md` → Goals; User Stories 1 and 2; Core Feature 2; Success Metric 3
- `_techspec.md` → API Contract 4; Testing Approach; Build Order 2
- ADR-0293, ADR-0375
