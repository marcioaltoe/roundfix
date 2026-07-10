---
task: task_01
spec: 0024-context-efficient-runs
status: pending
type: backend
complexity: high
---

# Task 01: Capture failed Verification diagnostics

## Overview

Move authoritative Verification output from the caller stream into bounded Run
artifacts and expose a typed command-failure contract for the repair cycle. The
slice is verifiable for successful, failed, cancelled, and infrastructure-error
commands before any Agent repair behavior is added.

## Requirements

1. MUST capture combined Verification stdout/stderr in an attempt-specific temporary artifact instead of streaming it to the Agent or caller.
2. MUST remove successful output and atomically retain failed command output at the TechSpec path.
3. MUST distinguish an exited Verification command from cancellation, process-start failure, and artifact filesystem failure.
4. MUST emit attempt-aware `daemon.verification` events with command phases and exactly one aggregate verdict per attempt.
5. MUST include only the diagnostic path and wrapped failure in events, never the command output body.
6. MUST apply the capture contract to Task command sequences and review Batch Verification without changing their current final settlement policy.

## Subtasks

- [ ] Add the artifact-backed Verifier result and typed command error.
- [ ] Create deterministic attempt paths with atomic failure retention.
- [ ] Remove successful artifacts and preserve infrastructure error identity.
- [ ] Add command-phase and aggregate-verdict Run Events.
- [ ] Route Task and review Verification through the capture contract.
- [ ] Cover real shell behavior and engine integration paths.

## Acceptance Criteria

- [ ] A successful command leaves no Verification output artifact and emits a passed aggregate verdict.
- [ ] A non-zero command retains combined output at the expected attempt path and returns the typed command failure.
- [ ] Cancellation, process-start failure, and artifact write failure remain infrastructure errors and are not typed as repairable command failures.
- [ ] Run Events contain attempt, phase, verdict, and diagnostic path metadata without output bytes.
- [ ] Multiple Task Verification commands stop on the first failure and produce one aggregate verdict for that attempt.
- [ ] Caller progress contains verdict summaries but no dots, passing-package banners, timings, or raw command output.

## Verification

- `rtk go test ./internal/daemon ./internal/runevent` - expected: real-shell capture, artifact lifecycle, error classification, and event-shape tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/runevent/event.go`

## References

`_prd.md` -> User Story 1; Core Features 1, 3; Success Metrics. `_techspec.md` -> Interfaces: Verification capture; Data Models: diagnostic artifacts; Build Order 1. ADR-0014; ADR-0038.
