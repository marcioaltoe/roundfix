---
task: task_01
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Add the supported official Codex adapter identity contract.
- [ ] Resolve and inspect the effective adapter command deterministically.
- [ ] Add typed adapter-lineage and adapter-version failures.
- [ ] Surface bounded adapter evidence through Setup and Doctor diagnostics.
- [ ] Cover official, legacy, unknown, missing, and unsupported adapters.
- [ ] Protect non-Codex adapter behavior with regression tests.

## Acceptance Criteria

- [ ] The supported official adapter reports ready with its command, package,
      and version.
- [ ] The legacy Zed adapter reports a lineage failure and the official package
      update action even when `acpx` and Codex CLI are current.
- [ ] An unknown executable named `codex-acp` fails closed instead of being
      accepted by name.
- [ ] A supported package at an unsupported version reports the observed and
      required versions without changing any config.
- [ ] Missing-adapter and non-Codex regression cases retain their established
      exit and diagnostic contracts.
- [ ] Adapter inspection never writes User Config, Project Config, Run,
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

