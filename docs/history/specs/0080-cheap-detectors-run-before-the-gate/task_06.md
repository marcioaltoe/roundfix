---
task: task_06
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: chore
complexity: high
---

# Task 06: Declare the two-tier verification contract

## Overview

Tooling Task two of two, under the same authorization. The gate's mechanical
stage runs the repository gate cold and complete at about ninety seconds when
the same gate on an unchanged tree costs under five. The fix is a declared
contract every adopting repository inherits: local verification is incremental
and fast, CI judges the complete tree from a fresh run, and the clause says
which tier answers which question.

The contract is expressed per profile in terms of the commands that profile
declares. An adopting repository may be a Bun or Rust repository with no
`make` at all, so a clause naming a tool would be unmeetable by construction.

## Requirements

1. MUST author the two-tier clause in the Baseline modules the authorization
   bounds, stating that the local tier is incremental, the CI tier is complete
   and fresh, and which question each answers.
2. MUST express the tiers per profile as declared commands, never as a named
   tool, so a repository without `make` can satisfy the contract.
3. MUST state that a profile declaring no incremental command inherits the
   clause as unmet rather than silently satisfied.
4. MUST adopt the regenerated managed-block postimages into the corresponding
   `docs/agents/` guides.
5. MUST add this repository's own incremental verification target to the
   `Makefile` under the name `verify-incremental`, and MUST NOT change what
   `make verify` means for CI. The name is fixed here because the declared
   Verification has to be able to fail: `verify:` already exists, so asserting
   that some verify-shaped target is present approves this Task before any work
   happens.
6. MUST run the module chain per the measured choreography: bootstrap the
   Source Baseline manifest rows for the new clauses, then run
   `make baseline-digests` twice, since the maintained fixture is the chain's
   first step and converges only on the second pass.
7. MUST NOT correct the maintained Source Baseline expectation this change
   invalidates. That correction is task_07's, landing as its own commit after
   this one, because a consequent fix folded into an authorized tooling commit
   fails the tooling-authority gate.
8. MUST change only the authorization's bounded files, their sanctioned
   deterministic digest fallout, and this task file.

## Subtasks

- [ ] Author the per-profile two-tier clause in the modules.
- [ ] Add this repository's incremental target.
- [ ] Run the chain and adopt both postimages.

## Acceptance Criteria

- [ ] Both guides carry the two-tier clause expressed as declared commands.
- [ ] A profile with no incremental command reads as unmet, not satisfied.
- [ ] `make verify` still means the complete gate.
- [ ] The digest chain is converged.
- [ ] The maintained expectation is left for task_07, and the diff stays
      inside the bounded files, fallout, and this task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-proof-cost.md
- interface: internal/baseline/assets/modules/core.json
- interface: internal/baseline/assets/modules/spec-workflow.json
- interface: Makefile

## Verification

- `grep -rqi 'incremental' docs/agents/ && grep -rqi 'complete' docs/agents/`
  — expected: exit 0; the two tiers are adopted into the guides.
- `grep -q '^verify-incremental:' Makefile`
  — expected: exit 0; the incremental target exists beside an unchanged
  `verify`.
  — expected: exit 0; the digest chain is converged.
  — expected: exit 0; the consequent fixture correction was not folded into
  this commit.
  — expected: exit 0; nothing outside the bounded files and sanctioned fallout
  changed.

## References

- `_prd.md` → Core Feature 6; User Story 5; Goals 3 and 4.
- `_techspec.md` → Build Order 6; Decisions (per-profile commands).
- ADR-0081.
- `references/2026-08-03-verification-performance-contract.md` → the measured
  tier costs this clause is built on.

## Result

### Implementation

- Versioned the core module, its agent-instructions guide, and the existing
  verification rule. Added `clause.core.verification-two-tiers`: the active
  Baseline Profile's declared incremental command answers whether the current
  change remains valid through fast local checks, while its declared complete
  command answers in CI whether the complete tree satisfies the repository
  contract from a fresh run. A missing incremental command leaves the clause
  unmet rather than satisfied by omission.
- Versioned the spec-workflow module, its spec-routing guide, and the existing
  routing rule. Added `clause.spec.verification-two-tiers` with the same
  profile-relative contract at the Task/assembled-tree boundary. Neither
  clause names `make` or another tool.
- Added `verify-incremental` as the fixed cached local target over
  `fmt-check`, `test`, `skills-sync-check`, `skills-check`, and `build`.
  The existing `verify: fmt-check $(VERIFY_TEST_TARGET) skills-sync-check
  skills-check build` declaration is byte-identical to `HEAD`, preserving
  CI's `VERIFY_TEST_TARGET=test-budget` complete fresh path.
- Bootstrapped the two marker-delimited Source Baseline entries and manifest
  rows. The first sanctioned regeneration replaced their temporary spans and
  digests and reported `changed:true`. Baseline managed refresh applied the
  reviewed three-file postimage plan for both guides plus
  `docs/agents/setup-context.json`; a second regeneration reported
  `changed:false`.
- Left `maintainedSourceBaselineEntries = 132` untouched while the regenerated
  Source Baseline identity now records 134 entries. That deliberate consequent
  correction remains task_07's slice.

### Focused checks

- Pre-change inspection found no `verification-two-tiers` clause or Source
  Baseline row, and the maintained identity carried 132 entries.
- Focused `jq -e` assertions over both module clause objects exited 0. They
  proved each clause names the active Profile's declared incremental and
  complete commands, assigns a distinct question to each tier, requires the
  complete tier from a fresh run, and treats a missing incremental command as
  unmet.
- The first
  `GOCACHE=/private/tmp/roundfix-task06-gocache rtk make baseline-digests`
  pass exited 0 with `changed:true`. The second exited 0 with
  `changed:false` after managed postimage adoption.
- `roundfix baseline update --repo . --no-skills --format text` then exited 0
  as `Baseline update: current` with `File changes: 0`, proving both adopted
  managed blocks and their Setup Manifest record match the current catalog.
  The first apply attempt could not create its Git-private transaction journal
  inside the sandbox; the identical approved-plan command was rerun with the
  required filesystem authorization and exited 0 with all three postimages
  verified.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk go test
  ./internal/baseline -count=1 -run
  "TestCatalogCompatibility|TestFormatterComposition"` exited 0 with two
  passing focused tests.
- A focused manifest assertion found both new rows with nonzero ordered spans
  and non-placeholder SHA-256 digests. `rtk make -n verify-incremental`
  expanded the cached local target without executing it.
- Exact `HEAD` and worktree inspections printed the same `verify:` declaration.
  `rtk git -c core.fsmonitor=false diff --check` exited 0, and
  `internal/baseline/preservation_test.go` has no diff.

### Acceptance-criterion evidence

1. Both public guide managed blocks contain their module-authored two-tier
   clause, and a fresh managed refresh reports zero file changes.
2. Both module and guide clauses state that a missing incremental command
   leaves the Profile contract unmet; omission is never treated as success.
3. The `verify:` declaration matches `HEAD` exactly. Only the separate
   `verify-incremental` target fixes the local tier to the cache-reusing
   `test` target.
4. The required regeneration sequence ended with `changed:false`; both focused
   catalog/formatter tests and the current-catalog managed refresh passed.
5. The maintained 132-entry expectation is unchanged against the regenerated
   134-entry identity. The changed-path audit contains only the authorized
   modules, guides, and Makefile; this Task file; the Setup Manifest postimage;
   and `internal/baseline/assets/**` plus `internal/baseline/testdata/**`
   deterministic regeneration fallout. The Task-file frontmatter transition
   was the pre-existing Daemon-owned change.

### Handoff boundary

- The Daemon-owned commands under `## Verification` were not run in this Agent
  turn. Task status remains Daemon-owned; no commit, push, or pull request was
  created.
