---
task: task_02
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 02: Establish the embedded Baseline catalog authority

## Overview

Move the canonical Baseline catalog into the Go CLI and prove that catalog
loading no longer depends on an installed Agent Skill. This slice establishes
the deterministic authority that every later planning and apply path consumes.

## Requirements

1. MUST embed the canonical profiles, modules, decisions, templates,
   Repository Capabilities, retention transitions, and setup snapshots in the
   Go binary.
2. MUST strictly validate schemas, duplicate keys, references, decision
   effects, module cycles, template tokens, and all maintained legacy inputs.
3. MUST reproduce every exact catalog identity and normalized catalog output
   recorded by the compatibility corpus.
4. MUST expose a cohesive deterministic catalog API without importing CLI,
   Agent, config, store, network, or installed-skill packages.
5. MUST reject catalog drift or invalid embedded assets during tests and build
   verification.

## Subtasks

- [ ] Create the Go-owned embedded catalog boundary.
- [ ] Port strict catalog loading and cross-reference validation.
- [ ] Port canonical serializers and catalog digest domains.
- [ ] Compare Go catalog outputs with exact characterization fixtures.
- [ ] Add mutation tests for invalid catalog relationships.

## Acceptance Criteria

- [ ] The binary loads all maintained built-in Baseline Profiles from embedded Go assets.
- [ ] Catalog loading succeeds when no setup skill is installed.
- [ ] Exact-parity catalog fixtures produce identical normalized bytes and digests.
- [ ] Unknown references, duplicate keys, cycles, and invalid template tokens fail closed.
- [ ] Mutation tests prove each catalog invariant independently.

## Context

- instruction: `docs/adr/0066-context-driven-baseline-execution-belongs-to-the-cli.md`
- interface: `.agents/skills/setup-context-driven/assets`
- interface: `skills/skills.go`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestEmbeddedCatalog|TestCatalogDigest|TestCatalogCompatibility|TestCatalogMutation'` — expected: embedded loading, exact parity, strict validation, and mutation cases pass.
- `rtk go test -count=1 ./skills` — expected: existing skill embedding remains valid while the Baseline catalog gains its independent authority.
- `rtk make verify` — expected: the full gate passes with no runtime dependency on setup-skill assets.

## References

- `_prd.md` → Goals 1, 4–5; User Story 4; Core Features 1, 4, 19–21.
- `_techspec.md` → System Architecture; Data Models: Catalog; Build Order 2.
- ADR-0066 → CLI runtime authority.
- ADR-0072 → catalog parity before cutover.
