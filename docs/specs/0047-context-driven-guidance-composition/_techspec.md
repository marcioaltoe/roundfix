---
spec: 0047-context-driven-guidance-composition
prd: _prd.md
created: 2026-07-24
---

# Context-Driven guidance composition — Technical Spec

## Executive Summary

This feature extends the existing embedded Baseline catalog, preservation
engine, and portable Plan assembly instead of adding a second renderer.
Repository policy follows a hybrid ownership order: an already confirmed typed
decision, a byte-exact repository-owned block in a semantic guide, then a
non-empty residual carrier. The design accepts more classification and marker
validation so update can redistribute rules without rewriting them or changing
the public Plan and Result schemas. It also closes the greenfield alignment gap
by planning a repository-owned Profile adaptation before instruction
classification, while keeping required capabilities non-waivable.

## Project Constraints

- Identifier strategy: not applicable; no project-owned application
  identifiers are created.
- Authentication and HTTP: not applicable; no provider or application route is
  installed or changed.
- Active ADR obligations: ADRs 0058, 0067, 0069, 0070, 0071, 0073, 0074, and
  0075 govern retention, ownership, supervised analysis, portable planning, and
  recoverable apply.
- Tooling authority: no tooling configuration, script, ignore file, plugin
  declaration, or version pin is authorized for modification by this Spec.
- Sources: `docs/agents/agent-instructions.md`,
  `docs/agents/domain.md`, and `docs/agents/spec-routing.md`.

## System Architecture

`internal/baseline` remains the deterministic authority. The feature extends
five existing seams:

1. The embedded catalog declares the Instruction Hierarchy, semantic-owner
   metadata, complete ADR and Findings contracts, and the templates that render
   them.
2. The preservation engine inventories root Readoption entries and the three
   recognized repository-rule carriers. A new segment validator admits
   byte-exact clause proposals without granting ACP access to the checkout.
3. Plan assembly inserts repository-owned rule blocks outside setup-owned
   markers, removes empty legacy carriers and pointers, and records every
   change through the existing `ManagedEntry`, preimage, postimage, and
   retention ledgers.
4. Profile alignment can resolve a built-in mismatch through an in-memory
   repository-owned Profile draft. The draft is validated against the embedded
   catalog and becomes a postimage in the same final Plan.
5. `internal/cli` drives the new alignment loop and consolidated semantic
   review. The setup skill and user guide remain thin descriptions of this
   public behavior.

No arbitrary nested carrier becomes writable. A recognized
`specific-repository.md`, `repository.md`, or `repository-rules.md` file is a
bounded migration source because its only documented purpose is repository
rules; every other nested instruction carrier remains warning-only under
ADR-0070.

## Implementation Design

### Interfaces

Plan input accepts either an existing Profile ID or one strict draft. Existing
callers keep using `ProfileID`; the planner rejects simultaneous inputs.

```go
type ProfileDraftInput struct {
	SourceProfileID string
	Document        []byte
}

type PlanRequest struct {
	Repository   string
	ProfileID    string
	ProfileDraft *ProfileDraftInput
	Decisions    []DecisionValue
	Preservation RootPreservationRequest
}
```

Semantic analysis returns offsets and identities, never rewritten rule text:

```go
type RuleSegmentProposal struct {
	EntryID    string
	Start, End int
	Digest     string
}
```

The planner derives the bytes with `entry.SourceBytes[start:end]`. It accepts a
segment set only when entries appear in source order, ranges do not overlap,
and the union covers `[0,len(SourceBytes))` exactly once. Empty segments,
unknown managed IDs, inactive destinations, digest mismatches, and a semantic
destination for `non-governed` evidence fail closed.

Sealed analysis runs in two bounded stages. The first proposal supplies only
byte ranges; after deterministic coverage validation, the planner materializes
one Source Baseline Entry per range. The existing classification boundary then
receives those entries and proposes classifications and advertised
destinations. Both proposals are discarded when preferred and fallback ACP
attempts fail; the manual path keeps each original entry as one lossless
segment and can retain it in the residual carrier.

### Data Models

The active artifact set produces a `SemanticOwnerRegistry` keyed by managed ID.
Each entry contains the guide path, active module, accepted classifications,
and deterministic display title. The registry is included in the sealed
Analysis Snapshot so an ACP can choose only destinations the current Profile
will render. It contains no filesystem capability.

The strict classification dispositions are:

- `managed-entry` for a known source clause already represented by an active
  setup-owned rule or confirmed typed decision;
- `repository-document` for literal repository-owned bytes placed in the
  advertised semantic guide;
- `repository-rules` for accepted policy with no semantic owner; and
- `rejected` with a non-empty reason for non-governed evidence.

A repository-owned semantic block uses stable markers outside the
setup-owned block:

```markdown
<!-- roundfix:repository-rule:begin id=<stable-rule-id> -->
<literal source bytes>
<!-- roundfix:repository-rule:end id=<stable-rule-id> -->
```

The stable ID binds the Source Baseline Entry, relative byte range, and initial
destination. The body digest is not stored in the marker, so later
repository-owned edits remain valid. Future updates inventory the current body
as repository policy, preserve it by default, and use current preimages for
approval.

`SetupManifest.ManagedArtifacts` continues to list only setup-owned artifacts.
Repository-owned semantic blocks appear in the existing Plan retention ledger
and as `ManagedEntry.Kind == "repository-owned"`; adding them does not change
`roundfix/baseline-plan/v1` or `roundfix/baseline-result/v1`.

A Profile draft uses the existing
`roundfix/custom-baseline-profile/v1` document. The proposed canonical path is
`.roundfix/baseline/profiles/<id>.json`. Its modules preserve catalog dependency
order, its templates are recalculated from the selected modules and decisions,
and its capabilities can remove only profile-specific expectations.
`capability.context7`, `capability.exa`, and every other universal required
capability remain outside the draft and cannot be waived.

### API Contracts

The interactive root command moves alignment before instruction
classification:

1. select the built-in or repository-owned Profile and collect its decisions;
2. audit repository alignment and print every blocking and advisory
   divergence;
3. for profile-specific blockers, offer Profile change, repository-owned
   adaptation, or decline;
4. review every proposed module and capability removal, require a valid custom
   Profile ID, validate the draft, and repeat audit;
5. for universal required gaps, stop with the exact supported skill-restoration
   or installation next action;
6. only after alignment is ready, collect instruction classifications and
   present the complete final Change Plan.

Adaptation is a proposal, not a waiver. The Profile file appears beside all
other file changes and is written only after the final Plan Digest
confirmation.

Automation gains one additive planning input:

```text
roundfix baseline plan --repo . --profile-file <draft.json> \
  --decision-file <decisions.json> --format json
```

`--profile-file` and `--profile` are mutually exclusive. The file must be a
strict custom Profile document whose target path is safe and whose catalog
references validate. A valid draft is resolved in memory and included as an
exact postimage; an invalid, stale, conflicting, or already divergent target
returns exit `2` or `3` under the existing category contract. Apply still
consumes only the approved portable Plan and requires `--confirm-plan`.

The segmentation snapshot and proposal receive new strict schema identities
because byte ranges change their shape. Once admitted, the planner materializes
each range as a normal `ReadoptionSourceEntry`; the classification snapshot and
`setup-context-driven/decisions/0.0.1` then use their existing one-entry,
one-disposition contract. Semantic guide targets use the existing
`repository-document` destination, constrained to the paths advertised by the
active semantic-owner registry.

## Coverage Map

- Goal 1 → catalog hierarchy metadata, ordered root templates, semantic owner
  registry.
- Goal 2 → segment validator, semantic guide blocks, residual planner.
- Goal 3 → recognized-carrier migration, retention ledger, empty carrier
  removal.
- Goal 4 → ADR lifecycle overlay and Findings Operational Contract templates.
- Goal 5 → formatter fixtures, portable Plan assembly, audit and reapply
  suites.
- Goal 6 → alignment loop, Profile draft planning, public `--profile-file`.
- Story 1 → `SemanticOwnerRegistry` and repository-owned block assembly.
- Story 2 → byte-exact segmentation, consolidated review, retention evidence.
- Story 3 → rendered Instruction Hierarchy and precedence clauses.
- Story 4 → self-contained catalog rules and generated guides.
- Story 5 → ADR and Findings templates in `guide.docs-layout`.
- Story 6 → residual carrier planner and root-pointer removal.
- Story 7 → repository-owned Profile adaptation and repeated alignment.

## Integration Points

The sealed ACP classifier receives only the canonical snapshot and returns a
proposal. The preferred Codex `gpt-5.6-sol` `xhigh` attempt and the Codex
`gpt-5.5` `xhigh` fallback from ADR-0069 remain unchanged; deterministic
validation and human review remain authoritative.

The filesystem boundary stays inside the existing portable Plan and
recoverable transaction. Profile files and semantic guide edits use the same
safe-relative-path, preimage, postimage, lock, rollback, and Git-lineage checks
as managed guides.

Catalog changes flow through canonical setup snapshots, formatter golden
fixtures, the embedded catalog digest, and `baseline assets sync`. The
externally owned `domain-modeling/ADR-FORMAT.md` is an immutable comparison
fixture, not an output target. The repo-owned setup skill and Roundfix user
guide must ship in the same change as CLI behavior.

## Testing Approach

Unit tests cover segment range validation, exact byte reconstruction, marker
parsing, semantic owner admission, stable ordering, literal preservation,
empty residual removal, and rejection of arbitrary nested carriers. Mutation
tests delete a clause, destination, template section, or lifecycle field and
must fail catalog coverage.

Profile tests reproduce the Oraculum shape with a backend-only TypeScript
fixture. The built-in Profile must block; the reviewed draft must remove only
the confirmed profile-specific capabilities, keep universal requirements
blocking, enter the final postimage ledger, apply, and resolve identically from
human and automation inputs.

Integration tests run greenfield, update, legacy-carrier migration, Profile
adaptation, apply, formatter composition, audit, and empty reapply for every
affected maintained Profile. They assert no `repository.md`, no empty
`specific-repository.md`, one semantic owner per accepted segment, unchanged
external skill bytes, and zero second-pass managed delta.

Spec QA uses the real binary against separately authorized Fluxus greenfield
and update checkouts plus the Oraculum divergence reproduction. Repository
Verification recommendations run outside Baseline, then a fresh Plan must be
empty. The full Roundfix gate remains `rtk make verify`.

## Build Order

1. Add ADR lifecycle metadata, glossary terms, catalog hierarchy and complete
   ADR/Findings Operational Contracts.
2. Add the semantic-owner registry, segmented snapshot/proposal contracts, and
   byte-exhaustive validator (depends on: 1).
3. Add repository-owned block inventory and Plan assembly, recognized-carrier
   migration, retention projection, and empty residual cleanup (depends on: 2).
4. Add in-memory custom Profile draft resolution and include its canonical file
   in portable Plan assembly (depends on: 1).
5. Add the human alignment/adaptation loop and automation `--profile-file`
   contract, then move classification after ready alignment (depends on: 4).
6. Update every affected Profile, template, source corpus, transition,
   formatter fixture, and canonical setup snapshot (depends on: 1, 3, 5).
7. Update user documentation and the thin setup skill, then add documentation
   contract and skill-sync guards (depends on: 5, 6).
8. Run all-profile macro QA, authorized Fluxus greenfield/update journeys,
   Oraculum divergence/adaptation, formatter and repository Verification, audit,
   and empty reapply (depends on: 6, 7).

## Risks & Considerations

Clause segmentation increases review volume and can split Markdown structures
at unsafe boundaries. The byte-exhaustive validator prevents loss, while
fixtures must cover headings, lists, code fences, blank separators, and
multi-clause carriers. A proposal that cannot preserve valid boundaries falls
back to manual classification or the residual carrier.

Repository-owned blocks inside setup-managed files create two ownership zones.
Distinct markers and ownership-aware scanners are mandatory; setup rendering
can replace only its own block, while repository-rule planning can replace only
confirmed repository-owned blocks. Marker corruption blocks planning.

Profile adaptation can produce a Profile that is truthful today but too narrow
for planned work. The complete removal review and versioned repository file
make that trade-off visible. Adding a capability later requires updating the
repository-owned Profile and approving a fresh Plan.

Rollout is atomic through the existing Plan/apply transaction. Before apply,
rollback means declining the Plan. After apply, transaction recovery restores
preimages; a later policy rollback requires a new Plan because repository-owned
blocks and Profile files are approved repository content.

## Decisions

- Use typed decisions, then semantic repository-owned blocks, then a residual
  carrier. See [ADR-0074](../../adr/0074-repository-rules-use-hybrid-semantic-ownership.md).
- Preserve migrated rule bytes exactly and validate complete segment coverage.
- Keep repository-owned blocks outside setup-owned markers and out of
  `SetupManifest.ManagedArtifacts`.
- Resolve required built-in Profile mismatch through a confirmed
  repository-owned adaptation, not a waiver. See
  [ADR-0075](../../adr/0075-profile-divergence-uses-confirmed-repository-owned-adaptation.md).
- Keep arbitrary nested instruction carriers warning-only under ADR-0070.
- Keep public Plan, Result, apply, recovery, and exit-code schemas unchanged;
  add only the explicit planning input required for Profile drafts.
