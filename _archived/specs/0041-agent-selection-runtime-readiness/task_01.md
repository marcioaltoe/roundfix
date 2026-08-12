---
task: task_01
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: high
---

# Task 01: Prove official Codex adapter readiness

## Overview

Make adapter readiness depend on the effective Codex ACP adapter's official
package lineage and supported version, not on a same-named executable existing
on PATH. The existing Setup and Doctor surfaces must distinguish a supported
adapter from legacy, unknown, and unsupported installations before profile
proof begins.

## Requirements

1. MUST define one Agent-owned supported Codex adapter package/version contract
   and reuse it across adapter resolution, health checks, and diagnostics.
2. MUST resolve the effective adapter command using the documented override and
   ACPX built-in precedence before inspecting identity.
3. MUST accept Codex adapter readiness only when the executable proves the
   official package lineage and supported version.
4. MUST classify legacy lineage, unknown lineage, missing executables, and
   unsupported versions without treating PATH presence as success.
5. MUST report the effective command, observed package/version when available,
   and one deterministic official install or update action.
6. MUST keep diagnosis read-only and exclude command output that can expose
   environment values, credentials, or unbounded runtime metadata.
7. MUST preserve existing Claude and OpenCode adapter resolution behavior.

## Subtasks

- [x] Add the supported official Codex adapter identity contract.
- [x] Resolve and inspect the effective adapter command deterministically.
- [x] Add typed adapter-lineage and adapter-version failures.
- [x] Surface bounded adapter evidence through Setup and Doctor diagnostics.
- [x] Cover official, legacy, unknown, missing, and unsupported adapters.
- [x] Protect non-Codex adapter behavior with regression tests.

## Acceptance Criteria

- [x] The supported official adapter reports ready with its command, package,
      and version.
- [x] The legacy Zed adapter reports a lineage failure and the official package
      update action even when `acpx` and Codex CLI are current.
- [x] An unknown executable named `codex-acp` fails closed instead of being
      accepted by name.
- [x] A supported package at an unsupported version reports the observed and
      required versions without changing any config.
- [x] Missing-adapter and non-Codex regression cases retain their established
      exit and diagnostic contracts.
- [x] Adapter inspection never writes User Config, Project Config, Run,
      worktree, Session, or artifact state.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/setup.go`
- interface: `internal/cli/doctor.go`

## Verification

- `rtk go test ./internal/agent -run 'Test(CheckAdapter|ResolveAdapterCommand|AdapterProbe)' -count=1` — expected: official, legacy, unknown, missing, unsupported-version, and non-Codex adapter cases pass.
- `rtk go test ./internal/cli -run 'Test(RunSetup|RunDoctor).*Adapter' -count=1` — expected: Setup and Doctor render bounded identity evidence and make no config changes on failure.
- `rtk go test -race ./internal/agent ./internal/cli -run 'Test(CheckAdapter|ResolveAdapterCommand|AdapterProbe|RunSetup.*Adapter|RunDoctor.*Adapter)' -count=1` — expected: adapter readiness and injected diagnostics are race-free.

## References

- `_prd.md` → User Stories 4 and 9; Core Features 1 and 10; Success Metrics.
- `_techspec.md` → Adapter Provisioning and Identity; Error Taxonomy and
  Diagnostics; Build Order 1.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → effective adapter identity is part of proof.

## Result

Implemented one Agent-owned Codex adapter contract for the official
`@agentclientprotocol/codex-acp` package with compatibility floor and tested pin
`1.1.4`. Adapter resolution now preserves stdio override, ACPX User Config,
then ACPX built-in precedence; the built-in Codex command resolves to
`npx -y @agentclientprotocol/codex-acp`. Codex readiness executes only the
effective command's bounded `--version` probe, while Claude and OpenCode keep
their existing executable-presence behavior.

Setup and Doctor now render the effective command plus observed package and
version. Legacy and unknown lineage failures use
`adapter_lineage_unknown`; official packages below the compatibility floor use
`adapter_version_unsupported`. Every Codex failure names the deterministic
`npm install -g @agentclientprotocol/codex-acp@1.1.4` action. Setup no longer
creates a bare Codex ACPX override and stops before Agent proof or configuration
surfaces when adapter inspection fails.

Acceptance evidence:

- `TestCheckAdapterProvesOfficialCodexPackageAndVersion` proves the supported
  command, package, and version evidence.
- `TestCheckAdapterClassifiesUnreadyCodexAdapters` covers legacy Zed lineage,
  an unknown same-named executable with secret-like extra output, and an
  unsupported official version without exposing raw output.
- `TestRunSetupReportsAdapterFailuresWithoutWrites` proves legacy and
  unsupported-version failures leave ACPX, User, and Project Config bytes
  identical, start no Agent proof, and do not reach later mutation surfaces.
- `TestRunDoctorReportsAdapterFailureWithNextAction` proves Doctor reports the
  legacy identity and official update action while continuing unrelated
  read-only checks.
- `TestACPXProbeMissingAdapterNamesInstallCommandBeforeSession`,
  `TestCheckAdapterPreservesNonCodexResolutionWithoutVersionExecution`, and the
  resolution table protect missing-adapter, Claude, OpenCode, and override
  behavior.
- Code-path inspection confirms adapter readiness reads ACPX config and invokes
  only the bounded `--version` process; it has no User Config, Project Config,
  Run, worktree, Session, or artifact writer dependency.

Verification:

- `rtk go test ./internal/agent -run 'Test(CheckAdapter|ResolveAdapterCommand|AdapterProbe)' -count=1`: passed.
- `rtk go test ./internal/cli -run 'Test(RunSetup|RunDoctor).*Adapter' -count=1`: passed.
- `rtk go test -race ./internal/agent ./internal/cli -run 'Test(CheckAdapter|ResolveAdapterCommand|AdapterProbe|RunSetup.*Adapter|RunDoctor.*Adapter)' -count=1`: passed.
- `make verify`: passed with 1,583 Go tests, 79 setup-context-driven tests,
  Roundfix skill synchronization, and the CLI build.

Follow-ups: none for this Task slice.
