---
spec: 0049-baseline-preservation-idempotency
prd: _prd.md
created: 2026-07-26
---

# Baseline Preservation idempotency — Technical Spec

## Executive Summary

This correction separates the bytes that require Preservation classification
from the complete root preimage that requires immutable backup. The
preservation engine will expose only new unmarked root spans as Source Baseline
Entries; setup-managed blocks and repository-owned rules already represented
by the retention inventory will not re-enter classification. After every
source entry has an approved disposition, portable Plan assembly will rebuild
the active `AGENTS.md` from its owned managed and repository-owned regions
instead of appending those regions to the consumed source. The design derives
root consumption only from the ready Source Baseline and its complete
disposition set so the apply contract can remain strict and unchanged.

## Project Constraints

- Identifier strategy: not applicable — the implementation creates no
  Internal Identifier or application identity. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — only local Baseline inventory,
  Plan assembly, and filesystem postimages change. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADRs 0058, 0060, 0064, 0069, 0071,
  0073, 0074, and 0078 require exhaustive accounting, supervised semantic
  decisions, portable preimage-bound Plans, recoverable apply, and
  single-owner rule retention. ADR-0070 still limits mutation to root carriers
  and leaves arbitrary nested carriers warning-only. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation is proposed or
  authorized. Source: `docs/agents/agent-instructions.md`.

## System Architecture

The correction extends the existing `internal/baseline` preservation and Plan
assembly seams. It adds no package, public command, dependency, or persistence
schema.

`PlanRootPreservation` continues to load complete root carrier bytes and derive
their content identities for immutable backups. Its classification projection,
however, contains only unmarked spans with nonblank repository-authored bytes.
Setup-managed spans remain structural evidence but are omitted from the Source
Baseline presented to semantic analysis. Recognized
`docs/agents/specific-repository.md` and repository-owned blocks in semantic
guides remain under `planSpecificRepository` and
`inventoryRepositoryRuleBlocks`. A Setup Manifest whose catalog and Profile
identities still resolve proves those carriers already have repository
ownership, so they are not independent migration sources. An absent,
incompatible, or stale Setup Manifest continues to expose recognized carriers
for explicit Readoption.

For repositories already written by the defective apply behavior, a compatible
Setup Manifest plus a verified `AGENTS.<digest>.md` payload proves that the
matching unmarked root bytes were part of an earlier approved transfer. Those
bytes are removed from the classification projection and the live stale copy
is consumed without asking ACP to make a second, potentially different,
classification. If new unmarked bytes coexist with that backed-up payload,
only the new remainder enters classification.

`assemblePostimages` consumes a preservation-owned root path only after its
current Source Baseline has complete valid dispositions or its stale payload is
proven by the compatible Manifest and verified prior backup. For `AGENTS.md`,
it starts that postimage from the managed regions that remain active rather
than from the full preimage, then applies the normal managed-block and
repository-rule assembly. This removes classified or previously approved
source bytes without weakening backup or retention validation. A root with no
new unmarked bytes and no stale payload is not consumed again and produces no
backup or classification work.

The existing portable Plan, `validatePlanApplyContract`, recoverable
transaction, and preimage comparison remain the mutation authority.

## Implementation Design

### Interfaces

No request field or persisted signal is added. Plan assembly treats a root path
as consumed only when:

1. it is a writable root source selected by preservation inventory;
2. it contains at least one nonblank unmarked Source Baseline Entry;
3. every current entry passes the existing classification and disposition
   validation; and
4. its complete current content has a valid immutable backup plan.

A repair path may also consume a root whose unmarked payload exactly matches a
verified content-addressed backup and whose current Setup Manifest resolves the
same catalog and Profile ownership. This evidence is kept only in the
in-memory preservation result; it adds no public Plan or Result field.

A helper projects classification entries without changing backup identity:

```go
func partitionUnmarkedRootSource(
	path string,
	content []byte,
) []ReadoptionSourceEntry
```

It uses the existing managed-marker parser and preserves original byte offsets,
digests, carrier digest, and source bytes for every admitted span. Empty
inter-block whitespace does not become a decision. Managed blocks remain
available to `containsOnlySetupManagedGuidance` and other ownership checks.

### Data Models

No persisted schema changes. `ReadoptionSourceBaseline.ByteCount` becomes the
sum of admitted Source Baseline Entry bytes rather than the full size of root
and recognized repository-rule carrier preimages. For every admitted carrier,
the Source Baseline identity still binds its complete current bytes, while its
entries bind the exact admitted ranges and bytes. Managed-byte drift and stale
classification proposals therefore fail closed without presenting managed
bytes as decisions.

The backup continues to bind the complete carrier content identity. Therefore
one addition outside managed markers changes the digest-addressed backup path
even when all managed bytes are unchanged. An existing backup with the same
digest is reused; a new backup is planned only for a new complete carrier
identity.

Already retained repository-owned semantic blocks are authoritative through
the Upgrade Retention Contract. The canonical
`docs/agents/specific-repository.md` is authoritative through
`planSpecificRepository`. Neither is copied into the Source Baseline on an
ordinary compatible Preservation update. When a later root rule also resolves
to Repository-Specific Normative Rules, the planner preserves the existing
carrier and appends the newly approved exact bytes. It replaces a recognized
carrier only when that carrier itself entered the reviewed Source Baseline.

### API Contracts

The public command shape and result schemas do not change.

For a first or newly changed Preservation:

1. audit reads the complete current carrier;
2. the human or automation flow reviews only new unmarked Source Baseline
   Entries;
3. the Change Plan shows the complete digest-addressed backup, every
   disposition and retention destination, and the rebuilt root postimage;
4. apply verifies the approved Plan Digest and exact preimages;
5. the live root contains only active generated hierarchy/pointers and any
   explicitly retained root-owned regions.

For an unchanged compatible Preservation, the Source Baseline has no new
entries. Planning reuses current decisions and retention evidence, emits no
backup, and the final Change Plan has zero file changes.

If an unmarked byte cannot be accounted for, the existing action-required or
blocked outcome remains. Failure never strips the source or creates a backup.

## Coverage Map

- Goal 1 → unmarked-source projection and consumed-root postimage assembly.
- Goal 2 → fresh-plan integration journey and backup de-duplication assertions.
- Goal 3 → complete-preimage backup identity and later-addition journey.
- Story 1 → ready-disposition detection and root reconstruction.
- Story 2 → compatible unchanged Preservation test.
- Story 3 → new-unmarked-span classification, backup, apply, and second rerun
  test.

## Integration Points

Semantic segmentation and classification continue to receive a sealed Analysis
Snapshot. The smaller Source Baseline does not grant new filesystem access or
change the preferred/fallback ACP Runtime contract.

Plan assembly continues to use catalog artifacts, repository-rule inventory,
preimages, postimages, and the complete managed-entry ledger.
`validatePlanApplyContract` remains responsible for proving that every planned
root consumption has the exact immutable backup and that every postimage is
preimage-bound. Transaction recovery continues to restore the complete
preimage, including consumed root bytes.

No user-guide or thin-skill update is required because the public contract
already promises semantic redistribution, immutable backup, and empty
reapply. A documentation contract test should instead prevent those existing
claims from drifting.

## Testing Approach

Unit tests for preservation inventory will cover managed-only roots, mixed
managed and unmarked roots, empty whitespace around managed blocks, malformed
markers, recognized residual carriers, and stable source identities. They will
assert that backup identity uses the complete carrier while classification
uses only admitted unmarked spans.

Plan integration tests will begin with a repository whose `AGENTS.md` contains
real repository rules. They will build and apply an approved semantic
redistribution, assert the exact original bytes in
`AGENTS.<digest>.md`, assert those bytes are absent from the live root and
present only in their selected owners, then build a fresh compatible
Preservation Plan and require zero postimage changes and no backup entry.

A second journey will add a new unmarked rule to the managed-only root. It must
produce only that new classification input, create a backup for the new
complete root identity, consume the new source after apply, and converge to a
zero-change third Plan. A migration-repair journey will start with source bytes
already present in their semantic or residual owner, a valid Manifest, and the
original content-addressed backup. It will assert that the stale bytes do not
re-enter classification, consuming the live-root copy does not append a
duplicate block, and an adjacent new rule remains the sole classification
input. Existing tests continue to cover repository-owned semantic edits,
residual carrier edits, preimage drift, backup collisions, atomic rollback,
and recovery.

## Build Order

1. Add preservation inventory and Plan journey tests that reproduce duplicate
   classification, second-backup creation, and live-root duplication.
2. Project only new unmarked root spans into the Source Baseline and derive
   consumed root paths after complete disposition validation (depends on: 1).
3. Rebuild consumed `AGENTS.md` postimages from owned regions while preserving
   exact backup, retention, preimage, rollback, and recovery contracts
   (depends on: 2).
4. Add the later-unmarked-addition convergence journey and documentation
   contract guard, then run full Verification (depends on: 3).

## Risks & Considerations

The primary risk is deleting a byte before it has a proven destination.
Consumption is therefore downstream of the existing exhaustive classification
and disposition validator and upstream of the unchanged preimage-bound apply
contract. A missing disposition, stale Source Baseline, backup collision, or
invalid semantic owner prevents the path from entering `ConsumedRootPaths`.

Managed-marker corruption must continue to fail closed; an invalid marker
cannot be treated as harmless unmarked content and silently consumed. Literal
blank separators around managed blocks may be normalized by root
reconstruction, but nonblank bytes always require disposition.

Repository-owned semantic blocks must not re-enter classification merely
because their bodies were edited. Their current bodies are already
repository-owned retention evidence. Direct unmarked additions to `AGENTS.md`
must still enter classification. A historical backup is repair evidence only
when its filename digest matches its exact bytes, the bytes are a contiguous
part of an unmarked root span, setup-managed guidance is also present, and the
current Manifest proves compatible ownership. Invalid, unrelated, or stale
backup files cannot suppress classification.
are different: they intentionally reopen Preservation and produce a new
content-addressed backup before migration.

The correction deliberately applies to the generated `AGENTS.md` root. Safe
aliases that target it follow the same bytes. Independent root carriers without
a generated managed owner retain ADR-0070's existing preservation behavior
until a separately approved design defines their canonical postimage.

## Decisions

- Rebuild the live `AGENTS.md` after confirmed distribution instead of keeping
  migrated source bytes beside their new owners. See ADR-0078.
- Classify only new nonblank unmarked root spans.
- Back up the complete current root content once per distinct digest.
- Keep recognized residual and semantic repository-rule carriers in the
  existing retention inventory, not the Source Baseline.
- Preserve public Plan, Result, apply, recovery, confirmation, and exit-code
  contracts.
