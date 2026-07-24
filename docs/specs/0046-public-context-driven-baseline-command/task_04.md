---
task: task_04
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Implement narrow Git root and lineage inspection.
- [x] Implement root-anchored no-follow inventory and preimage recording.
- [x] Classify safe, nested, and unsafe instruction carriers.
- [x] Expose deterministic plan preflight results and exit categories.
- [x] Add repository and real-CLI macro tests.

## Acceptance Criteria

- [x] Equivalent clones at different paths produce the same repository identity.
- [x] Unrelated dirty files do not invalidate the bounded snapshot.
- [x] Safe aliases retain source evidence without duplicate source entries.
- [x] External, escaping, cyclic, unreadable, and special-file targets block safely.
- [x] Nested carriers remain unchanged and every detected conflict is visible.
- [x] Preflight performs zero repository mutations and zero network operations.

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

## Result

Roundfix now performs the read-only Baseline repository preflight behind
`baseline plan`. The inspector accepts dirty, detached, and upstream-free Git
worktrees with a committed `HEAD`; derives a clone-stable identity from the Git
object format and sorted reachable root commits; and records only normalized
repository-relative evidence.

The bounded inventory uses a root-confined filesystem handle and never follows
an escaping target. It audits root and nested `AGENTS.md`/`CLAUDE.md` carriers
plus every file under `docs/agents`, excludes the maintained ignored trees,
records missing root mutation candidates, deduplicates trusted source evidence
for safe aliases, warns for every nested carrier, and makes unsafe or opaque
carriers apply-blocking.

`roundfix baseline plan --repo <path> --format text|json` now reports the
repository identity, bounded snapshot, warnings, blocking findings, and next
action without prompts. A safe tracer-bullet inspection returns
`action_required` with exit 3; invalid Git state, invalid arguments, and unsafe
carriers return preflight exit 2. The compiled real-CLI test confirms JSON
output contains no checkout path.

Acceptance evidence:

- `TestRepositoryIdentityEquivalentClones` cloned one committed repository to
  a different path and produced the same repository digest. Detached, dirty,
  and missing-upstream behavior passed independently.
- `TestRepositorySnapshotDigestChangesWithBoundedBytesOnly` preserved the
  snapshot digest after an unrelated dirty file appeared and changed it after
  bounded carrier bytes changed.
- `TestInstructionAliasRetainsOneSourceEvidence` resolved two safe root aliases
  to one target path and content identity while emitting one trusted source
  record.
- `TestInstructionAliasUnsafeTargetsBlock` and
  `TestInstructionAliasUnreadableTargetBlocks` covered absolute external,
  lexical escape, cycle, directory, unreadable, and special-file targets.
- `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths` recorded
  every bounded regular carrier, omitted ignored and unrelated paths, and
  emitted the nested conflict warning without mutation.
- `TestRepositoryInspectionNoMutation` compared modes, link targets, paths, and
  bytes across the complete tracked and untracked repository tree before and
  after inspection. `TestRepositoryInspectionUsesNarrowReadOnlyGitCommands`
  observed only the four local lineage reads. Production inspection contains
  no network client or filesystem write operation.

Verification:

- Pre-change focused tests failed to compile because `InspectRepository`,
  repository identity/snapshot models, carrier kinds, findings, and preimages
  did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestRepositoryIdentity|TestBoundedInventory|TestInstructionAlias|TestBaselinePlanPreflight'`:
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache go test -count=1
  ./internal/baseline -run TestRepositoryInspectionNoMutation`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache go vet
  ./internal/baseline ./internal/cli`: passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache make verify`: passed,
  including 1,788 Go tests in 21 packages, the setup-context-driven Python
  suites, embedded asset and shipped-skill checks, formatting, and the build.

The isolated `GOCACHE` keeps build artifacts inside the Task Worktree sandbox.
The Daemon remains responsible for the task file's verbatim authoritative
Verification commands.
