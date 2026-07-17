---
task: task_04
spec: 0035-agent-selection-profiles
status: pending
type: backend
complexity: high
---

# Task 04: Configure and validate complete profiles

## Overview

Let humans and automation write a complete profile atomically and prove effective profiles through the installed ACP Runtimes without starting a Run. Interactive and file-driven configuration share one strict normalization, confirmation, persistence, and JSON contract.

## Requirements

1. MUST implement `profiles configure --scope user|project` with Interactive Input or strict `--file` input, plus `--dry-run` and deterministic `--json` reporting.
2. MUST require one complete Preferred Selection and at least one complete ordered fallback for every configured profile before any write.
3. MUST show the complete normalized profile and target scope before the existing confirmation boundary and write the target config atomically only after validation and approval.
4. MUST write only the new `profiles` schema, preserve unrelated config, and reject invalid scope, incomplete fragments, same-scope legacy/new conflicts, and malformed input without partial files.
5. MUST implement read-only `profiles validate` for all required profiles or one category, deduplicating exact tuples for disposable proof while reporting every referencing category.
6. MUST apply runtime, official or custom model, and exact reasoning effort through the same pinned acpx boundary used by live sessions and close every disposable session on success or error.
7. MUST provide actionable text and JSON diagnostics without editing runtime-owned configuration, authentication, credentials, or recommendation data.

## Subtasks

- [ ] Parse strict file-driven profile fragments and scopes.
- [ ] Build Interactive Input with recommendation context and fallback collection.
- [ ] Normalize, preview, confirm, and atomically persist complete profiles.
- [ ] Implement dry-run and JSON change reports.
- [ ] Deduplicate and prove requested profile tuples through disposable sessions.
- [ ] Add legacy migration and recovery diagnostics.
- [ ] Cover atomicity, session closure, and no-mutation validation cases.

## Acceptance Criteria

- [ ] Valid User and Project profiles are written atomically and resolve from their reported source.
- [ ] Dry-run and failed configuration leave config bytes unchanged and report `changed: false` in JSON.
- [ ] Interactive configuration cannot confirm until at least one fallback is complete.
- [ ] Validation proves identical tuples once in stable order but reports the result under every category that uses them.
- [ ] A failed proof names runtime, model, reasoning effort, affected categories, adapter error, and the next configuration or validation action.
- [ ] Every disposable session closes, and no Run row, worktree, Agent prompt, or runtime-owned setting is created.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `docs/adr/0040-reasoning-effort-is-assigned-only-when-configured.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/sessions.go`
- interface: `internal/cli/selection.go`
- interface: `internal/config/config.go`

## Verification

- `rtk go test ./internal/config ./internal/agent ./internal/cli -run 'Test(ProfilesConfigure|ProfilesValidate|ProfileProof|ProfileConfigAtomic)' -count=1` — expected: interactive/file, dry-run, JSON, proof deduplication, exact reasoning, closure, and atomic failure cases pass.
- `rtk go test -race ./internal/agent ./internal/cli -run 'Test(ProfileProof|ProfilesValidate)' -count=1` — expected: disposable proof and cleanup paths are race-free.

## References

- `_prd.md` → Goals 1-6 and 8; User Stories 2, 4, and 7; Core Features 1, 4-6, and 8; Success Metrics.
- `_techspec.md` → Profile CLI: configure and validate; Official Model Catalog and adapter proof; Error and exit behavior; Build Order 4.
- ADR-0040 → exact reasoning effort must be advertised, applied, displayed, and persisted consistently.
