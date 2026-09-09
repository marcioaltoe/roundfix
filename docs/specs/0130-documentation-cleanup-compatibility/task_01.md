---
task: task_01
spec: 0130-documentation-cleanup-compatibility
status: completed
type: test
complexity: high
---

# Task 01: Preserve authorization and regeneration after cleanup

## Overview

Repair the approved five-file boundary described by the TechSpec. Demonstrate
that current grants and historical evidence work without a workflow directory,
while invalid authority and unowned writes remain rejected.

## Requirements

1. MUST implement all four Core Features through the TechSpec's existing seams.
2. MUST preserve deep equality and negative coverage; no suppressions or private
   exemption list. Use filesystem-only tests in suiteguardcontract.
3. MUST change only the following bounded paths and this Task's own evidence:
   - `internal/baseline/derived_ownership_test.go`
   - `internal/speccheck/mechanical_test.go`
   - `internal/speccheck/governed_repocontract_test.go`
   - `internal/suiteguardcontract/regeneration.go`
   - `internal/suiteguardcontract/regeneration_test.go`

## Subtasks

- [ ] Add positive and negative Spec-grant discovery and ownership tests.
- [ ] Implement bounded discovery and command-only ownership resolution.
- [ ] Adapt regeneration fixtures and assert parity with Baseline output ownership.
- [ ] Recover historical grant evidence with explicit unavailable-object outcomes.
- [ ] Record criterion-level evidence and limits in this Task's Result.

## Acceptance Criteria

- [ ] CF-1: approved current and legacy inputs work; proposed, null, malformed,
  unrelated, mismatched-consumer and symlinked inputs grant no authority.
- [ ] CF-2: command-only declarations match Baseline ownership and concrete
  positive/negative fixture expectations, including frozen/dedicated paths,
  sidecars, nested exceptions and invalid/conflicting declarations.
- [ ] CF-3: real historical grant bytes and non-empty bounded-path coverage remain
  tested; unavailable historical Task objects are explicit skips, with controlled
  accepted/refused Git-change replay providing non-vacuous coverage.
- [ ] CF-4/OE-1: regeneration fixtures exercise the missing-directory failure from
  PR #180 using Spec grants, preserving frozen boundaries and strict assertions.
- [ ] All implementation/test changes stay within the five approved paths.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `docs/agents/go.md`
- instruction: `docs/agents/specific-repository.md`
- interface: `internal/baseline/derived_ownership.go`
- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguardcontract/contract.go`

## Verification

- `test -s internal/suiteguardcontract/regeneration_test.go && grep -q 'func TestCleanupRegenerationDiscovery' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestCleanupRegenerationDiscovery$'` — the new positive/negative reader contract executes; absence of its test cannot pass.
- `grep -q 'func TestCleanupRegenerationOwnershipParity' internal/baseline/derived_ownership_test.go && go test -count=1 -tags repocontract ./internal/baseline -run '^(TestCleanupRegenerationOwnershipParity|TestOutputsForCommand|TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries)$'` — actual ownership and cleanup-compatible regeneration fixtures agree.
- `grep -q 'func TestCleanupHistoricalGrantEvidence' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^(TestCleanupHistoricalGrantEvidence|TestAuditJudgesTheGrant|TestEveryBoundedPathIsGoverned)$'` — historical and controlled grant evidence is exercised, with unavailable original Task objects explicitly reported.

## References

- `_prd.md` → Goals 1–3, Core Features 1–4, Outside Evidence OE-1.
- `_techspec.md` → all Design sections, Coverage Map, Research and decision evidence.
- ADR-0149 — command authority and tree-owned outputs.

## Result

Implemented the five-file cleanup-compatibility slice. The suiteguard reader now
discovers legacy records and approved Spec-contained grants, validates operative
Spec metadata, refuses symlinked or malformed authority, and resolves command-only
grants from `DERIVED_DIGEST_PATHS` plus ownership records. The Baseline fixtures
now copy current Spec grants while proving `docs/workflow` is absent, and the
dedicated fixture declares its command through a valid Spec grant without listing
the command's outputs.

Acceptance evidence:

- CF-1: `TestCleanupRegenerationDiscovery` covers an approved grant under
  `references/`, compatible unfrontmattered legacy input, missing optional roots,
  deterministic reads, and rejection of proposed, null-date, malformed,
  unrelated, mismatched-consumer, empty-action, unsafe-path and symlinked inputs.
  It also proves ownership symlinks and invalid roots return no partial authority.
  Focused check `rtk go test ./internal/suiteguard ./internal/suiteguardcontract`
  passed with 12 tests.
- CF-2: `TestCleanupRegenerationOwnershipParity` compares the suiteguard reader
  with Baseline output ownership and exact fixture slices. Its positive cases
  cover sanctioned, dedicated, frozen, nearest-directory, sidecar and nested
  exception behavior; its negative cases cover unsafe, invalid, incomplete and
  conflicting declarations. Focused check
  `rtk go test ./internal/baseline -run 'CleanupRegenerationOwnershipParity|OutputsForCommand'`
  passed with 18 tests, preserving deep equality including the non-nil empty
  result for unknown commands.
- CF-3: `TestCleanupHistoricalGrantEvidence` reads all 42 records from ancestor
  `81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4`, requires real non-empty bytes and
  non-empty bounded-path coverage, checks the proof-cost regeneration
  enumeration, and replays one accepted regeneration output plus one refused
  governed path in a temporary Git repository. Focused check
  `rtk go test -tags repocontract ./internal/speccheck -run 'CleanupHistorical'`
  passed. A verbose focused run of `TestAuditJudgesTheGrant` reported the
  unavailable historical Task commits `419a4661ac769ff7ee6ce5423bd795185c859d01`
  and `c80e1266658929f68e8046af82f88e13392dc56d` as explicit skipped subtests;
  its four controlled rejection subtests passed.
- CF-4/OE-1: the regeneration fixture no longer copies or recreates the deleted
  workflow tree, and the dedicated command fixture uses command-only ownership
  through a Spec grant. The focused Baseline check above exercises this grant
  shape and asserts that `docs/workflow` remains absent. The existing tagged
  regeneration gates retain their byte-for-byte tree, idempotence, frozen-path
  and wrong-command assertions for Daemon Verification.
- Scope: postflight `rtk git status --short` listed only the five authorized
  implementation paths plus this Task file; `rtk git diff --check` exited 0.

Verification feedback repair, attempt 1:

- The Daemon isolated two failures in the tagged declared-step gate. The copied
  fixture still declared the fixed `make baseline-digests` token inside update
  tests even when the fixture's grant and ownership record named a dedicated
  command. The fixture adapter now makes that copied in-process declaration read
  the exact command under exercise; repository source behavior remains unchanged.
- Suiteguard now refuses a frozen output before the outer byte comparison runs.
  The declared-step helper preserves the existing assertion by mapping an exact
  suiteguard violation for an untouched probe to the same
  `rewrote frozen artifact` error, while unrelated command failures keep their
  original diagnostics.
- Focused check
  `rtk go test -count=1 -tags repocontract ./internal/baseline -run '^TestDeclaredStepRegenerationAndFrozenBoundaries/dedicated/synthetic_plan_characterization$'`
  passed with the parent and dedicated subtest.
- Focused check
  `rtk go test -count=1 -tags repocontract ./internal/baseline -run '^TestDeclaredStepRegenerationAndFrozenBoundaries/frozen_resolved_path_rejects_rewrite$'`
  passed with the parent and frozen-boundary subtest.

Limits: the commands under `## Verification`, including the tagged commands that
execute the full regeneration steps, were not run because the Daemon owns them.
No Task status, Task Graph, sibling Task, commit, push or pull request was changed
by this implementation turn.
