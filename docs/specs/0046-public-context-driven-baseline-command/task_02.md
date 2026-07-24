---
task: task_02
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Create the Go-owned embedded catalog boundary.
- [x] Port strict catalog loading and cross-reference validation.
- [x] Port canonical serializers and catalog digest domains.
- [x] Compare Go catalog outputs with exact characterization fixtures.
- [x] Add mutation tests for invalid catalog relationships.

## Acceptance Criteria

- [x] The binary loads all maintained built-in Baseline Profiles from embedded Go assets.
- [x] Catalog loading succeeds when no setup skill is installed.
- [x] Exact-parity catalog fixtures produce identical normalized bytes and digests.
- [x] Unknown references, duplicate keys, cycles, and invalid template tokens fail closed.
- [x] Mutation tests prove each catalog invariant independently.

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

## Result

Roundfix now owns the Baseline catalog in `internal/baseline`. The Go package
embeds the canonical profiles, modules, decisions, templates, Repository
Capabilities, Setup Snapshots, Upgrade Retention transitions, compatibility
inputs, and Source Baseline assets without importing CLI, Agent, config, store,
network, or installed-skill packages.

The loader rejects malformed or duplicate-key JSON, unknown schema fields,
invalid versions and paths, duplicate catalog identifiers, unknown
cross-references, dependency cycles, invalid decision effects, undeclared
template tokens, Setup Snapshot digest drift, invalid retention mappings, and
invalid Source Baseline references. Its immutable API exposes catalog entries,
template bytes, dependency-ordered modules, normalized identity bytes, and the
domain-separated Catalog Digest.

Acceptance evidence:

- `TestEmbeddedCatalog` loaded `go-cli-tui`, `rust-cli`, and
  `standard-typescript-monorepo` from `go:embed`, resolved the expected module
  order, and proved returned entry bytes cannot mutate catalog state.
- `TestEmbeddedCatalog` changed the process directory to an empty temporary
  directory before loading. `rtk env
  GOCACHE=/private/tmp/roundfix-task02-go-cache go list -deps
  ./internal/baseline` listed only Go standard-library dependencies, proving
  loading does not read an installed setup skill.
- `TestCatalogCompatibility` matched the complete path-ordered normalized
  fixture byte for byte and reproduced Catalog Digest
  `sha256:139704a2191e3c571bf7bbd6085517d7fdab260f1450ff1400c33955fa0b8b71`.
  It also loaded the maintained v2 compatibility catalog.
- `TestCatalogMutation` passed 13 isolated mutations covering duplicate JSON
  keys, unknown schema fields, duplicate identifiers, unknown profile/setup/
  template/decision/retention/Source Baseline references, module and decision
  cycles, undeclared template tokens, and Setup Snapshot digest drift.
- Byte comparisons between the Go-owned catalog directories and the canonical
  setup-skill catalog returned no differences.

Verification:

- Pre-change: `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test
  -count=1 ./internal/baseline -run
  'TestEmbeddedCatalog|TestCatalogDigest|TestCatalogCompatibility|TestCatalogMutation'`
  failed to compile because the Go catalog API did not exist.
- Post-change: the same focused command passed all catalog and mutation tests.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1
  ./internal/baseline`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go vet
  ./internal/baseline`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1
  ./skills`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache make verify`: passed
  1,744 Go tests, 256 canonical setup tests, 256 distributed setup tests,
  catalog loading, owned-skill checks, and the Roundfix build.

The isolated `GOCACHE` keeps verification inside the Task Worktree sandbox;
the Daemon remains responsible for the task's verbatim authoritative
Verification commands.
