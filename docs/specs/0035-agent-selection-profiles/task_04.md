---
task: task_04
spec: 0035-agent-selection-profiles
status: completed
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

- [x] Parse strict file-driven profile fragments and scopes.
- [x] Build Interactive Input with recommendation context and fallback collection.
- [x] Normalize, preview, confirm, and atomically persist complete profiles.
- [x] Implement dry-run and JSON change reports.
- [x] Deduplicate and prove requested profile tuples through disposable sessions.
- [x] Add legacy migration and recovery diagnostics.
- [x] Cover atomicity, session closure, and no-mutation validation cases.

## Acceptance Criteria

- [x] Valid User and Project profiles are written atomically and resolve from their reported source.
- [x] Dry-run and failed configuration leave config bytes unchanged and report `changed: false` in JSON.
- [x] Interactive configuration cannot confirm until at least one fallback is complete.
- [x] Validation proves identical tuples once in stable order but reports the result under every category that uses them.
- [x] A failed proof names runtime, model, reasoning effort, affected categories, adapter error, and the next configuration or validation action.
- [x] Every disposable session closes, and no Run row, worktree, Agent prompt, or runtime-owned setting is created.

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

## Result

- Added `profiles configure --scope user|project` with strict file fragments, Interactive Input fallback collection, normalized preview, confirmation, atomic config replacement, dry-run, and `roundfix/profiles-configure/v1` JSON.
- Added config-owned profile fragment parsing and atomic profile writes that preserve unrelated config, write only the `profiles` schema, and reject invalid scope, incomplete profiles, malformed YAML, legacy-only targets, and same-scope legacy/new conflicts before touching files.
- Added read-only `profiles validate` with `roundfix/profiles-validate/v1` JSON. It resolves all required categories or one requested category, deduplicates exact Agent Selection tuples in stable order, probes each distinct tuple through `agent.Runner.Probe`, and reports every category reference.
- Added ACP proof tests showing exact non-empty reasoning is applied through the pinned acpx path, model-managed empty reasoning is respected by existing coverage, no prompt is sent, and disposable sessions close on success and selection failure.
- Acceptance evidence:
  - `TestProfileConfigAtomicWritesUserAndProjectProfiles` verifies User and Project writes resolve from `user` and `project` sources and preserve unrelated config.
  - `TestProfilesConfigureFileWritesProjectProfileJSON` verifies file-driven Project Config writes, JSON reporting, normalized preview, and no Run Database creation.
  - `TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged` verifies dry-run and invalid input preserve bytes and JSON reports `changed: false`.
  - `TestProfilesConfigureInteractiveRequiresCompleteFallbackBeforeConfirm` verifies confirmation is not reached without a complete fallback.
  - `TestProfilesValidateDeduplicatesProofsAndReportsEveryReference` verifies stable deduplication, all referencing categories, no Agent prompt, and no Run Database creation.
  - `TestProfilesValidateFailedProofNamesTupleAffectedCategoriesAndRecovery` verifies runtime, model, reasoning effort, affected category, adapter error, recovery action, read-only config behavior, and no Agent prompt.
- Verification passed:
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/config ./internal/agent ./internal/cli -run 'Test(ProfilesConfigure|ProfilesValidate|ProfileProof|ProfileConfigAtomic)' -count=1` → `10 passed in 3 packages`.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go test -race ./internal/agent ./internal/cli -run 'Test(ProfileProof|ProfilesValidate)' -count=1` → `4 passed in 2 packages`.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/config ./internal/agent ./internal/cli -count=1` → `786 passed in 3 packages`.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go run -buildvcs=false ./cmd/roundfix profiles configure --help && GOCACHE=/private/tmp/roundfix-gocache rtk go run -buildvcs=false ./cmd/roundfix profiles validate --help` → passed.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk make verify` → passed; included `rtk go test ./...` with `1502 passed in 20 packages`, skill sync check, and build.
