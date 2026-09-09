---
task: task_07
spec: 0119-spec-contained-authorization
status: pending
type: docs
complexity: medium
---

# Task 07: Make the Spec-contained record canonical in the Baseline

## Overview

The audit already reads authorization records, but no canonical guide says
where they live, so an adopting repository learns the location only by reading
Go source. Name the Spec-contained `_authorization.md` as the home in the
Baseline modules and regenerate the managed guides from them. The slice is
verifiable on its own: the rendered guides carry the clause and the managed
refresh converges.

This is an authorized tooling Task. It may change only
`internal/baseline/assets/modules/core.json`,
`internal/baseline/assets/modules/spec-workflow.json`,
`internal/baseline/assets/modules/context-workflow.json`,
`docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
`docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, the derived pins
that `make baseline-digests` rewrites as sanctioned regeneration, and this Task
file. Stop before any other mutation. The bounded set and the sanctioned
regeneration come from the approved grant in [_authorization.md](_authorization.md).

## Requirements

1. MUST name `<spec-root>/<slug>/_authorization.md` as the home for tooling
   authorization records, stating the record's required shape: approval state,
   maintainer decision date, permitted actions, exact bounded repository paths,
   consuming Spec, and any sanctioned regeneration.
2. MUST state that a proposed, absent, contradictory, or withdrawn record grants
   nothing, and that an executor cannot widen the grant it depends on.
3. MUST state the committed-provenance contract for executing authored
   Verification, including that read-only checking stays available for any
   source and that an untrusted source needs an execution approval naming the
   approved revision.
4. MUST state that a new or widened grant lands in target ancestry before the
   consuming squash delivery.
5. MUST preserve the existing legacy record location as still readable, so
   historical grants keep resolving, and MUST NOT declare their old paths
   present as worktree files.
6. MUST extend existing clause identities where an obligation already exists
   rather than replacing them, and MUST NOT weaken any current clause.
7. MUST regenerate the managed guides and manifest from the modules through the
   public Baseline update, then regenerate the sanctioned derived pins, and
   MUST NOT hand-edit a derived pin value.
8. MUST leave the rendered guides converged, so a second managed refresh
   produces no further change.

## Subtasks

- [ ] Add the record-home and record-shape clauses to the owning modules.
- [ ] Add the non-granting-state and no-self-widening clauses.
- [ ] Add the committed-provenance execution clause.
- [ ] Add the ancestry-before-squash clause and preserve the legacy location.
- [ ] Regenerate the managed guides and the sanctioned derived pins.

## Acceptance Criteria

- [ ] Searching the rendered guides for `_authorization.md` returns the new
      clauses, where the same search returns nothing before this Task.
- [ ] A rendered clause states that a proposed or absent record grants nothing
      and that an executor cannot widen its own grant.
- [ ] A rendered clause states the committed-provenance execution contract and
      that read-only checking remains available for any source.
- [ ] A rendered clause states that a new or widened grant lands in ancestry
      before the consuming squash delivery.
- [ ] The legacy record location is still described as readable for historical
      grants.
- [ ] Running the managed refresh a second time reports no further change, and
      the derived pins match a fresh regeneration byte for byte.

## Context

- instruction: `docs/agents/specific-repository.md`
- interface: `internal/baseline/assets/modules/core.json`

## Verification

- `grep -q '_authorization.md' docs/agents/docs-layout.md && grep -q '_authorization.md' docs/agents/spec-routing.md` — the rendered guides name the record home, which they do not today.
- `grep -q 'committed provenance' docs/agents/agent-instructions.md && grep -q '_authorization.md' internal/baseline/assets/modules/core.json` — the execution contract is canonical and sourced from a module rather than hand-written into a rendered guide.
- `grep -rq '_authorization.md' internal/baseline/assets/modules || exit 1; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text > /dev/null || exit 1; git diff --quiet -- docs/agents internal/baseline || { git --no-pager diff -- docs/agents internal/baseline; exit 1; }` — the clause is sourced from a module and the managed refresh converges: a second run changes nothing.
- `grep -rq '_authorization.md' internal/baseline/assets/modules || exit 1; make baseline-digests || exit 1; git diff --quiet -- internal/baseline || { git --no-pager diff -- internal/baseline; exit 1; }` — the module edit landed and the sanctioned derived pins match a fresh regeneration.

## References

- `_prd.md` → User Stories 1; Core Features 1, 3, 5; Goals 1.
- `_techspec.md` → System Architecture: Authoring guidance; Build Order 5.
- ADR-0130, ADR-0149.
