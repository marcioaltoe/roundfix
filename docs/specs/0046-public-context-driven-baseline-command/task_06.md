---
task: task_06
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 06: Resolve profile alignment decisions

## Overview

Turn repository evidence and one selected Baseline Profile into explicit
blocking or advisory decisions. The planner provides useful HTTP, PostgreSQL,
and Verification evidence without inventing repository policy.

## Requirements

1. MUST resolve exactly one valid built-in or repository-owned Baseline
   Profile and compare repository evidence with its requirements.
2. MUST block on unresolved required divergence and keep advisory divergence
   visible without treating it as policy.
3. MUST emit bounded HTTP route candidates, observed methods and scopes, and
   source digests without assigning contract mode, owner, or rationale.
4. MUST distinguish PostgreSQL implementation evidence from a missing accepted
   repository contract and name accepted contract paths when required.
5. MUST label commands repository-executable only after local declaration
   validation; portable roles and profile expectations remain distinct.
6. MUST remain local, read-only, network-free, and command-execution-free.

## Subtasks

- [ ] Port Repository Capability evaluation and evidence ranking.
- [ ] Resolve profile alignment and divergence decision states.
- [ ] Add HTTP route-candidate and source-digest projection.
- [ ] Separate implementation, contract, portable-role, and executable-command evidence.
- [ ] Add finding regression and profile parity tests.

## Acceptance Criteria

- [ ] Required divergence prevents a ready Plan until explicitly resolved.
- [ ] Advisory divergence never blocks and never becomes inferred policy.
- [ ] HTTP candidates contain facts but no inferred Normative Clause.
- [ ] PostgreSQL diagnostics report found implementation evidence separately from contract absence.
- [ ] A nonexistent formatter or Verification script is never labeled executable.
- [ ] Equivalent evidence and answers produce equivalent normalized decisions across interaction modes.

## Context

- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_capabilities.py`
- interface: `.agents/skills/setup-context-driven/assets/profiles/standard-typescript-monorepo.json`
- interface: `docs/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestProfileAlignment|TestRequiredDivergence|TestHTTPRouteCandidates|TestPostgreSQLEvidence|TestExecutableVerificationCommand'` — expected: alignment, evidence classification, and finding regressions pass.
- `rtk go test -count=1 ./internal/baseline -run TestCapabilityAuditNoExecution` — expected: audit uses declared local evidence and invokes no repository or network command.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 4 and 8; Core Features 3–4, 10, 14, 17–18.
- `_techspec.md` → Data Models: Catalog and RepositorySnapshot; Testing Approach: Fluxus assertions; Build Order 3–4.
- ADR-0063 → repository-owned HTTP contract policy.
