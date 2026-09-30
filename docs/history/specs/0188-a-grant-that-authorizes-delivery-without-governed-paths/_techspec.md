---
spec: 0188-a-grant-that-authorizes-delivery-without-governed-paths
prd: _prd.md
created: 2026-09-30
---

# A grant that authorizes delivery without Governed Paths — Technical Spec

## Executive Summary

The conflict is two rules on one list:

- The authorization reader refuses an approved record whose `paths` holds no
  entry (`internal/authorization/authorization.go`,
  `classifyAuthorizationRecord`). Spec 0119 introduced the rule when a record
  only bounded Governed Paths.
- ADR-0130's contract refuses any listed path that is not governed.

The fix changes only the reader: a Spec-contained record whose `paths` is the
explicit empty YAML sequence is operative and bounds nothing. Every consumer
of the record already treats its `Paths` as the bounded set, so an empty set
authorizes no governed mutation without further change. The trade-off this
design accepts is a second meaning for an empty field, so omission and
declaration must be told apart. The alternative, a new field such as
`governed: none`, would add vocabulary that every reader and guide must learn.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0179 (this Spec), ADR-0130,
  ADR-0160, ADR-0178, ADR-0166, ADR-0167 and ADR-0176 hold as the PRD states.
  The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104,
  ADR-0117, ADR-0155 and ADR-0156, and by ADR-0093 and ADR-0094. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix skill files ride the standing
  grant of 2026-09-18, recorded in [_authorization.md](_authorization.md)
  under the maintainer's 2026-09-30 request to release every remaining fix;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package or seam. The reader in `internal/authorization` is the single
place that decides whether a record is operative, and `internal/spec`
re-exports it.

The consumers of `Record.Paths` were checked:

- `internal/speccheck/undeclared.go` reports a declared Governed Path the
  record does not bound. With an empty set it reports every declared Governed
  Path.
- `internal/speccheck/mechanical.go` builds its bounded map from `Paths`. With
  an empty set, every governed change is outside the grant.
- `internal/daemon/task_engine.go` (`authorizationBoundsProjectConfig`) never
  treats an empty set as bounding `.roundfixrc.yml`.
- `internal/speccheck/constraints.go` maps a `paths` refusal with an empty
  list to an untyped finding. An explicit empty list no longer produces that
  refusal, and an absent or null one still does, so the branch stays correct
  unchanged.

## Implementation Design

### Interfaces

```go
// internal/authorization — unexported, inside the frontmatter parser
type authorizationFrontmatterPresence struct {
	// existing fields …
	pathsEmptySequence bool // paths is a YAML sequence with no items
}
```

No exported signature changes.

### Data Models

No schema change. `AuthorizationRecord.Paths` stays `[]string`. An explicit
empty sequence yields an empty, non-nil slice, and a null or absent field
yields nil. Only the presence flag decides.

### The explicit empty list

In `parseAuthorizationFrontmatter`, when the `paths` node is a sequence
(`yaml.SequenceNode`) with no content, set `presence.pathsEmptySequence`.

In `classifyAuthorizationRecord`, the `len(record.Paths) == 0` refusal is
skipped when `record.Role == AuthorizationRoleSpec` and
`presence.pathsEmptySequence`. Everything before it is unchanged: status,
granted, action, consuming and the operations vocabulary. A legacy record, a
null `paths:` and an absent `paths` all keep the refusal with code `paths` and
the same message.

### Guides

`docs/user-guide/commands.md` and the Delivery queue section of
`.agents/skills/roundfix/SKILL.md` gain one sentence after the start
requirement:

> A Spec that changes no Governed Path records `paths: []`; that explicit
> empty list grants the listed operations and bounds no Governed Path.

Then `make skills-sync` updates the mirror.

### API Contracts

1. API Contract: `roundfix deliver plan <slug>` — a Spec whose committed
   record declares `paths: []` and every delivery operation prints
   `spec	<slug>	approved	…` and exits `0`. The same record with `paths:`
   null prints `authorization refused: paths` and exits `1`.
2. API Contract: `roundfix deliver start <slug>` — accepts the `paths: []`
   Spec and refuses the null one, as `deliver plan` reports.
3. API Contract: QA report "Mechanical findings" — a Task commit that changes a
   Governed Path under a `paths: []` grant reports `QA-AUTH-PATHS`, naming the
   path as outside the grant's bounded files.

## Coverage Map

- Goal 1 → The explicit empty list; API Contracts 1-2.
- Goal 2 → System Architecture (consumers); API Contract 3.
- Core Feature 1 → The explicit empty list.
- Core Feature 2 → System Architecture; API Contract 3.
- Core Feature 3 → Guides.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 1, Testing Approach 2.
- Success Metric 3 → Testing Approach 3.

## Integration Points

- **Delivery Plan and start.** Both read the record through
  `spec.ReadSpecAuthorization` and need no change.
- **QA mechanical audit.** Reads the same record and needs no change.

## Testing Approach

1. **Reader.** A new `internal/authorization/empty_paths_test.go` covers:
   - an approved Spec record with `paths: []` and operations is granted, with
     an empty `Paths`;
   - `paths:` null and an absent `paths` are refused on `paths`;
   - a legacy-role record with an empty sequence is still refused.

   The existing `internal/spec` test `TestAuthorizationReaderRefusesEmptyPaths`,
   whose fixture writes a null `paths:`, stays green unchanged.
2. **Delivery Plan.** A new `internal/cli/deliver_plan_empty_paths_test.go`
   builds a Delivery Plan workspace whose record declares `paths: []` and
   every delivery operation. `deliver plan` approves it and `deliver start`
   accepts it. The same workspace with `paths:` null is refused on `paths`.
3. **Audit.** A new `internal/speccheck/mechanical_empty_grant_test.go` builds a
   temporary repository whose record declares `paths: []`, with a Task commit
   changing a Governed Path such as `Makefile`. The mechanical audit reports
   `QA-AUTH-PATHS`. A Task commit changing only an ordinary file reports
   nothing.
4. **Docs.** A phrase check on the guide, the skill and its mirror, plus
   `make skills-sync-check`.

## Build Order

1. The explicit empty list in the reader, its tests, the audit test and the
   guide and skill sentence, task_01 (depends on: none).
2. Terminal QA, task_02 (depends on: 1).

## Risks & Considerations

- **An authoring slip reads as a grant.** Only the literal empty sequence
  grants. The template-shaped `paths:` with nothing under it is null and still
  refuses.
- **A silent governed mutation.** Every governed audit reads the bounded set,
  which is empty, so a governed change is refused, and Testing Approach 3
  proves it.

## Decisions

- The literal empty sequence is the only form that grants operations without
  Governed Paths. See ADR-0179.
