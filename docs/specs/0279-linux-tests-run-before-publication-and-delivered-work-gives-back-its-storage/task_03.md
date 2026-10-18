---
task: task_03
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
status: pending
type: backend
complexity: high
---

# Task 03: Merged items give back their Item Binaries under a written retention policy

## Overview

Measured on 2026-10-18, ten Item Binary copies of the merged Specs 0277 and
0278 held 327 MB of the 412 MB of `~/.roundfix/artifacts`, and nothing reclaims
them (`_prd.md` → Sources 3; ADR-0282 says "this decision adds no reclamation
for them"). This Task makes the Delivery Queue remove a merged item's Item
Binary directory when it releases the item's Runs, adds an Item Binary sweep to
the GC Command for items that are merged or no longer queued, and writes the
retention policy of ADR-0290 into the GC Command guide. It is verifiable on its
own through release and GC tests on temporary artifact directories.

## Requirements

1. MUST answer the operator measurement of 2026-10-18 cited in `_prd.md` →
   Sources 3 by implementing API Contracts 8 and 9 of `_techspec.md`
   (ADR-0290).
2. MUST make `ReleaseMergedRuns`, after releasing the merged Spec's Runs,
   remove `<artifact dir>/delivery/<slug>/item-binary/` when it exists, and
   join a removal failure to its returned error so it becomes the item's
   cleanup warning. Every other file under `delivery/<slug>/`, including
   `item-binary-build.log` and `repository-gate.log`, stays. The engine,
   `internal/delivery/engine.go`, does not change.
3. MUST add the GC Command's Item Binary sweep of API Contract 9:
   - The data source is `<artifact dir>/delivery/*/item-binary` of the
     repository GC runs in and that repository's Delivery Queue, read through
     the readers `--dry-run` already uses.
   - A directory is reclaimable only when the queue has no item for its slug
     or holds that item at stage `merged`. Any other stage keeps it, and GC
     rereads the queue immediately before each removal.
   - A symlinked `item-binary` or a path resolving outside
     `<artifact dir>/delivery` is kept and never followed.
   - `--dry-run` writes nothing and prints the four lines and the path list of
     API Contract 9; the live run prints the removed forms and counts only what
     it removed.
   - With no Item Binary directory, GC output stays byte-identical to today's.
4. MUST add these tests, each calling `t.Parallel()` first, compiling for
   Windows and using temporary homes, artifact directories and Run Databases.
   A test MUST NOT read or write the real `~/.roundfix` or write a literal
   executable.
   - `internal/cli/deliver_item_binary_reclaim_test.go`:
     `TestMergedReleaseReclaimsTheItemBinaries`, which asserts the directory
     is gone and both logs remain.
   - `internal/cli/gc_item_binary_test.go`:
     `TestGCDryRunListsReclaimableItemBinaries`, which asserts the exact lines
     for a merged item, an unqueued slug and an item at stage `parked`, and
     that the dry-run removes nothing; and
     `TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems`, which asserts the
     live removal and that the parked item's directory remains.

   Every existing GC and release test keeps its name and passes unchanged.
5. MUST document the change:
   - `docs/user-guide/commands/gc.md` gains a section `## Item Binaries` with
     the lines of API Contract 9 and a section `## Retention policy` with the
     table of ADR-0290: what each area of Roundfix Home keeps, who reclaims
     it, and the read-only previews `roundfix gc --dry-run`,
     `roundfix gc compact` and `roundfix reconcile`;
   - `docs/user-guide/commands/deliver.md` and
     `.agents/skills/roundfix/references/deliver.md` replace "this operation
     does not reclaim them" with the merged release and the GC sweep;
   - `.agents/skills/roundfix/references/storage.md` names the Item Binary
     sweep, quoting its `Item Binaries reclaimable` line, and the retention
     policy.

   MUST search the guides and skill copies for text saying Item Binary copies
   are never reclaimed, and correct each. No text inside any
   `### QA settlement` section may change.
6. MUST add **Item Binary** and revise **GC Command** in `CONTEXT.md` through
   domain-modeling with ADR-0290's link; the revised **GC Command** names the
   Item Binary sweep.
7. MUST run these commands, because this Task edits user guides and Roundfix
   references. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Remove the merged item's Item Binary directory at its release.
- [ ] Add the GC Item Binary sweep with its dry-run report.
- [ ] Add the release and GC tests.
- [ ] Write the retention policy and correct the guides, the references and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A merged item's release removes its Item Binary directory and keeps its logs.
- [ ] GC reclaims Item Binaries only for merged or unqueued items, the dry-run lists them and removes nothing, and a parked item keeps its copies.
- [ ] No guide or skill copy still says Item Binary copies are never reclaimed.

## Consumers

- The GC Command report and the merged item release — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md`
- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `internal/cli/worktree_build_cache.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_release.go`
- interface: `internal/cli/gc.go`
- creates: `internal/cli/gc_item_binary.go`
- creates: `internal/cli/gc_item_binary_test.go`
- creates: `internal/cli/deliver_item_binary_reclaim_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestMergedReleaseReclaimsTheItemBinaries|TestGCDryRunListsReclaimableItemBinaries|TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMergedReleaseReclaimsTheItemBinaries TestGCDryRunListsReclaimableItemBinaries TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/gc.md .agents/skills/roundfix/references/storage.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Item Binaries reclaimable' || { printf 'missing Item Binary sweep in %s\n' "$f" >&2; exit 1; }; done; grep -q '^## Retention policy$' docs/user-guide/commands/gc.md || { printf 'missing retention policy section\n' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do ! tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'does not reclaim them' || { printf 'stale reclamation text in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Item Binary**:' || { printf 'missing glossary term Item Binary\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**GC Command**:' || { printf 'missing glossary term GC Command\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'GC|ItemBinary|Release|Merged|Storage' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the release and GC tests, the retention policy, the corrected deliver text, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 5; Core Feature 5; Core Feature 6; Core Feature 7; Success Metric 3; Glossary
- `_techspec.md` → API Contract 8; API Contract 9; API Contract 10; Invariant 5; Current behavior; Testing Approach; Build Order 3
- ADR-0290; ADR-0282; ADR-0266; ADR-0255; ADR-0033; ADR-0280; ADR-0189; ADR-0256; ADR-0259
