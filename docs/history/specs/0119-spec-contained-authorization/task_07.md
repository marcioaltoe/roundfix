---
task: task_07
spec: 0119-spec-contained-authorization
status: completed
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
   maintainer decision date, permitted actions, the closed list of permitted
   operations, exact bounded repository paths, consuming Spec, and any
   sanctioned regeneration. An absent operations list grants no operation, so
   the field is part of the shape rather than an optional extra.
2. MUST define the fields an execution approval carries — the repository, the
   approved source revision, the carrying artifact, and a digest of the approved
   command text — and state that an approval expires when any of them changes.
   Without a named representation an approval cannot be scoped or expired.
3. MUST state that a proposed, absent, contradictory, or withdrawn record grants
   nothing, and that an executor cannot widen the grant it depends on.
4. MUST state the committed-provenance contract for executing authored
   Verification, including that the comparison is over the authored projection
   and excludes Daemon-owned status and Result fields, including that read-only checking stays available for any
   source and that an untrusted source needs an execution approval naming the
   approved revision.
5. MUST state that a new or widened grant lands in target ancestry before the
   consuming squash delivery.
6. MUST preserve the existing legacy record location as still readable, so
   historical grants keep resolving, and MUST NOT declare their old paths
   present as worktree files.
7. MUST extend existing clause identities where an obligation already exists
   rather than replacing them, and MUST NOT weaken any current clause.
8. MUST regenerate the managed guides and manifest from the modules through the
   public Baseline update, then regenerate the sanctioned derived pins, and
   MUST NOT hand-edit a derived pin value.
9. MUST leave the rendered guides converged, so a second managed refresh
   produces no further change.

## Subtasks

- [ ] Add the record-home and record-shape clauses, including operations and the execution-approval fields, to the owning modules.
- [ ] Add the non-granting-state and no-self-widening clauses.
- [ ] Add the committed-provenance execution clause.
- [ ] Add the ancestry-before-squash clause and preserve the legacy location.
- [ ] Regenerate the managed guides and the sanctioned derived pins.

## Acceptance Criteria

- [ ] Searching the rendered guides for `_authorization.md` returns the new
      clauses, where the same search returns nothing before this Task.
- [ ] A rendered clause states that a proposed or absent record grants nothing
      and that an executor cannot widen its own grant.
- [ ] A rendered clause states the committed-provenance execution contract over
      the authored projection, names the execution-approval fields, and states
      that read-only checking remains available for any source.
- [ ] A rendered clause names the closed operations vocabulary and states that
      an absent list grants no operation.
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
- `grep -q 'committed provenance' docs/agents/agent-instructions.md && grep -q 'operations' docs/agents/spec-routing.md && grep -q '_authorization.md' internal/baseline/assets/modules/core.json` — the execution contract and the operations vocabulary are canonical and sourced from a module rather than hand-written into a rendered guide.
- `grep -rq '_authorization.md' internal/baseline/assets/modules || exit 1; raw="$(mktemp)"; before="$(mktemp)"; after="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text > /dev/null || exit 1; find docs/agents internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$before" || exit 1; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text > /dev/null || exit 1; find docs/agents internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$after" || exit 1; diff "$before" "$after"` — the clause is sourced from a module and the managed refresh converges: a second run reproduces the first byte for byte. The comparison is against a snapshot taken after the first refresh, because the Daemon runs Verification while the authorized edits are still uncommitted.
- `grep -rq '_authorization.md' internal/baseline/assets/modules || exit 1; raw="$(mktemp)"; before="$(mktemp)"; after="$(mktemp)"; make baseline-digests || exit 1; find internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$before" || exit 1; make baseline-digests || exit 1; find internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$after" || exit 1; diff "$before" "$after"` — the module edit landed and the sanctioned derived pins reproduce byte for byte on a second regeneration.

## References

- `_prd.md` → User Stories 1; Core Features 1, 3, 5; Goals 1.
- `_techspec.md` → System Architecture: Authoring guidance; Build Order 5.
- ADR-0130, ADR-0149.

## Result

Implemented the canonical Spec-contained authorization guidance in the owning
Baseline modules and carried the corresponding clauses into the three managed
guide regions. The existing clause identities remain in place: tooling
authorization now names the Spec-relative record, its required fields, the
closed operation vocabulary, non-granting states, and the no-self-widening
rule; the existing Verification guidance now carries the committed-provenance
contract; the existing documentation and bounded-execution clauses carry the
legacy-location and ancestry prerequisites.

Focused checks:

- Pre-change inspection found no `_authorization.md` clause in the rendered
  `docs/agents` guides and no committed-provenance clause in the Baseline
  modules.
- `rtk jq -e . internal/baseline/assets/modules/core.json`, the equivalent
  checks for `spec-workflow.json` and `context-workflow.json`, and
  `rtk git diff --check` passed.
- The focused source/guide search found the record home, committed-provenance
  fields, non-granting states, closed operations list, ancestry rule, and
  legacy readability statement in both their module sources and managed guide
  regions.
- `rtk proxy test ! -e docs/workflow/authorizations && rtk proxy test ! -e docs/workflow/authorizations/2026-09-09-spec-contained-authorization.md`
  passed; the legacy location is described without being present as a worktree
  file.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-build rtk go test ./internal/spec -run '^TestAuthorizationReaderTypesPermittedOperations$' -count=1`
  passed with 3 tests, confirming the documented six-operation vocabulary
  matches the reader contract.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-go-build make verify-incremental`
  stopped at the pre-existing formatter gate for
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go`; neither file is in this Task's
  allowlist.
- The first Daemon refresh diagnostic identified that the initial shape added a
  new clause without a source-baseline row and changed the protected tooling
  clause. The repair removed that new identity, restored the protected clause
  byte-for-byte, and folded the obligations into existing clause identities.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task07-go-build go test ./internal/baseline -run '^TestEmbeddedCatalog$' -count=1`
  passed after the repair; the embedded catalog now validates without the
  missing-row or protected-clause diagnostics.

Acceptance-criterion evidence:

1. The rendered `docs/agents/docs-layout.md` and `docs/agents/spec-routing.md`
   regions now name `<spec-root>/<slug>/_authorization.md`; the pre-change
   search had no such rendered clause.
2. `docs/agents/agent-instructions.md` states that a proposed, absent,
   contradictory, or withdrawn record grants nothing and that an executor
   cannot widen the grant it depends on.
3. The committed-provenance clause names the authored projection, excludes
   Daemon-owned `status` and `## Result`, names `repository`, `revision`,
   `artifact`, and `command_digest`, expires approval when any changes, and
   preserves read-only checking for every source.
4. The Spec routing clause names the closed `implement`, `commit`, `push`,
   `pull_request`, `merge`, and `release` vocabulary and states that an absent
   operations list grants no operation.
5. The bounded-execution clause states that a new or widened grant must land in
   the delivery target's ancestry before the consuming squash delivery.
6. The docs-layout clause preserves `docs/workflow/authorizations/` as a
   readable historical location and explicitly says its old paths are not
   required as worktree files; the focused absence check passed.
7. The final diff extends existing clause identities and retains their prior
   enforcement and guidance while appending the new obligations; the
   protected tooling clause is unchanged, no new clause identity was added,
   and no existing clause was replaced or weakened.
8. Public Baseline update and sanctioned digest regeneration were not run in
   this Agent turn because they are declared Verification commands owned by
   the Daemon. Their pending derived-state repair is recorded by the focused
   catalog diagnostic above.
9. The second-refresh convergence check remains Daemon-owned and pending;
   the generated guides and derived pins must be rechecked after regeneration.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains pending. No commit, push, pull request, Task status, or
other Task file was changed.
