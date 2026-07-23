---
task: task_03
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: medium
---

# Task 03: Manage repository-owned Baseline Profiles

## Overview

Expose repository-owned Baseline Profile authoring, inspection, and validation
through the public CLI. A maintainer can create one valid declarative profile
and resolve it against the embedded catalog without confusing it with Agent
Selection Profiles.

## Requirements

1. MUST implement `roundfix baseline profile init`, `show`, and `validate`
   with the TechSpec flag and output contracts.
2. MUST discover custom profiles only under
   `.roundfix/baseline/profiles/<id>.json`.
3. MUST allow custom profiles to compose only embedded modules, decisions,
   Repository Capabilities, and templates.
4. MUST reject remote references, custom assets, executable content, unknown
   IDs, profile composition, unsafe paths, and invalid schemas.
5. MUST keep stdout, stderr, JSON schemas, help, and exit categories aligned
   with the public CLI contract.

## Subtasks

- [ ] Add Baseline command dispatch and profile subcommand parsing.
- [ ] Implement repository profile discovery and strict loading.
- [ ] Implement init, text/JSON show, and validation results.
- [ ] Add CLI and filesystem tests for valid and rejected profiles.
- [ ] Add help and schema contract tests.

## Acceptance Criteria

- [ ] `profile init` creates a valid repository-owned profile from an allowed built-in source.
- [ ] `profile show` resolves the same normalized profile in text and JSON.
- [ ] `profile validate` accepts an ID or explicit profile path and reports stable failures.
- [ ] User-scoped profiles and repository-provided assets never participate.
- [ ] The command family does not alter Agent Selection Profile behavior.

## Context

- instruction: `docs/adr/0067-custom-baseline-profiles-are-repository-owned.md`
- interface: `internal/cli/profiles.go`
- interface: `internal/config/profiles.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestCustomProfile|TestBaselineProfileCommand|TestBaselineProfileHelp'` — expected: authoring, discovery, validation, output, and rejection cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline profile --help` — expected: help lists only init, show, and validate with the approved public grammar.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Story 4; Core Features 1, 3–4, 16, 18; Non-Goals / Out of Scope.
- `_techspec.md` → Data Models: CustomProfile; API Contracts: profile operations; Build Order 2 and 7.
- ADR-0067 → repository ownership and embedded-only composition.
