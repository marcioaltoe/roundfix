---
task: task_04
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 04: Inspect bounded Git repository state

## Overview

Deliver the read-only repository preflight and inventory used by every
Baseline workflow. The first `baseline plan` tracer bullet identifies the Git
lineage, current Baseline state, safe aliases, bounded evidence, and actionable
blocking conditions without changing repository bytes.

## Requirements

1. MUST require a Git worktree with at least one commit while accepting dirty
   state, detached HEAD, and missing upstream.
2. MUST derive the clone-stable repository identity and record normalized,
   root-relative bounded preimages for every consulted or mutable path.
3. MUST inventory every bounded instruction and agent-document carrier without
   following unsafe targets or reading special files as trusted evidence.
4. MUST report safe root aliases with target path and content identity, nested
   conflicts as warnings, and unsafe aliases as apply-blocking findings.
5. MUST expose preflight and action-required text/JSON results through
   `roundfix baseline plan` with no writes, network, repository commands, or
   prompts.

## Subtasks

- [ ] Implement narrow Git root and lineage inspection.
- [ ] Implement root-anchored no-follow inventory and preimage recording.
- [ ] Classify safe, nested, and unsafe instruction carriers.
- [ ] Expose deterministic plan preflight results and exit categories.
- [ ] Add repository and real-CLI macro tests.

## Acceptance Criteria

- [ ] Equivalent clones at different paths produce the same repository identity.
- [ ] Unrelated dirty files do not invalidate the bounded snapshot.
- [ ] Safe aliases retain source evidence without duplicate source entries.
- [ ] External, escaping, cyclic, unreadable, and special-file targets block safely.
- [ ] Nested carriers remain unchanged and every detected conflict is visible.
- [ ] Preflight performs zero repository mutations and zero network operations.

## Context

- instruction: `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- instruction: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- interface: `internal/preflight/preflight.go`
- interface: `.agents/skills/setup-context-driven/scripts/context_baseline.py`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestRepositoryIdentity|TestBoundedInventory|TestInstructionAlias|TestBaselinePlanPreflight'` — expected: Git, inventory, alias, no-write, and exit-contract cases pass.
- `rtk go test -count=1 ./internal/baseline -run TestRepositoryInspectionNoMutation` — expected: tracked and untracked repository bytes are identical before and after every inspection outcome.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 1 and 3; Core Features 2–3, 6–7, 13, 17, 19.
- `_techspec.md` → System Architecture; Data Models: RepositoryIdentity and RepositorySnapshot; Integration Points: Git and Filesystem; Build Order 3.
- ADR-0070 → repository-wide audit and root-only preservation.
- ADR-0071 → clone-stable identity and bounded preimages.
