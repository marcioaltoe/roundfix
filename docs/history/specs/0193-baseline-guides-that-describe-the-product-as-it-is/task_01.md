---
task: task_01
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: completed
type: backend
complexity: medium
---

# Task 01: No clause is rendered twice, and every clause keeps its force

## Overview

The backend module lists one paragraph under two entries of the same rule, so every adopter's backend guide renders it twice. This Task removes the duplicate entry and adds two checks over the embedded catalog: no two clauses share their text, and every clause keeps its enforcement level. The force check is the characterization the later Tasks of this Spec rely on: they rewrite clause guidance and must not change a clause's level.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST remove, from `rule.backend.boundary-contracts` in `internal/baseline/assets/modules/backend.json`, the entry whose `id` is `rule.backend.boundary-contracts`. The entry `clause.backend.boundary-contracts` keeps its guidance and its `mandatory` level.
2. MUST raise by one the versions of `rule.backend.boundary-contracts`, `guide.backend` and the `backend` module.
3. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
4. MUST add `internal/baseline/clause_characterization_test.go` with the four tests the TechSpec's Testing Approach 1 names:
   - the uniqueness check compares every clause-level and rule-level guidance string after lower-casing and collapsing whitespace;
   - its negative test runs the same function on a copy with one guidance repeated;
   - the force table lists every clause identifier with its enforcement level, written as a literal in the test, and fails for a missing clause, an unlisted clause and a changed level;
   - the rendering test counts the boundary sentence in the backend formatter golden and expects one.
5. MUST keep every other clause, in every module, byte-identical.
6. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Remove the duplicate entry and raise the three versions.
- [ ] Regenerate the pins, the golden and the Setup Manifest.
- [ ] Add the uniqueness, force and rendering tests, each negative case separate.

## Acceptance Criteria

- [ ] The backend formatter golden holds the boundary paragraph once.
- [ ] A catalog copy with a repeated guidance string is reported by the uniqueness check.
- [ ] The force table matches every clause of the embedded catalog, and the duplicate identifier is absent from both.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/clause_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestNoTwoBaselineClausesShareText|TestADuplicatedClauseTextIsReported|TestBaselineClauseForceIsCharacterized|TestTheBackendBoundaryParagraphRendersOnce|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestNoTwoBaselineClausesShareText TestADuplicatedClauseTextIsReported TestBaselineClauseForceIsCharacterized TestTheBackendBoundaryParagraphRendersOnce TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the four new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Clause texts; Testing Approach 1; Version changes
- `_prd.md` → Goals 3 and 4; Core Feature 1; Success Metrics 1 and 5
- [references/2026-09-30-the-backend-guide-renders-one-clause-twice.md](references/2026-09-30-the-backend-guide-renders-one-clause-twice.md)
- ADR-0186

## Result

Removed only the duplicate `rule.backend.boundary-contracts` clause entry.
The surviving `clause.backend.boundary-contracts` keeps its exact guidance and
`mandatory` force. Raised the backend module and guide versions from 4 to 5,
and the boundary-contracts rule version from 3 to 4, using in-place source
replacements. All other module sources and surviving clause bytes are unchanged.

Added the four specified characterization tests and three separate negative
force tests. The uniqueness helper compares clause-level and rule-level text
after lower-casing and collapsing whitespace. Its negative case changes a
copy's case and whitespace and requires the exact repeated identifier pair.
The force table is a literal inventory of all 128 clauses at `9e439dbb`, with
only the duplicate backend entry excluded; an independent Git-source comparison
confirmed that inventory. Legacy rule-level guidance has no enforcement field
and participates in uniqueness, while the force check covers clause entries.

### Focused evidence

All commands below used
`GOCACHE=/private/tmp/roundfix-0193-task01-gocache`. The default Go cache was
sandbox-inaccessible; no repository tool configuration was changed.

- Before the source edit, focused tests for uniqueness, force characterization
  and rendering failed with the exact duplicate pair, the unlisted duplicate
  identifier, and two boundary-sentence occurrences.
- `rtk proxy make baseline-digests` exited 0 with `ok: true` and
  `changed: true`. It generated the backend formatter golden, profile pin,
  catalog snapshots and four expected plan-characterization goldens. None was
  hand-edited.
- `rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  exited 0 and reported verified postimages with one Setup Manifest change.
  The first sandboxed attempt was refused when creating the Git-private
  transaction directory; the authorized elevated retry applied the same plan.
  A second sandboxed invocation exited 0, reported `File changes: 0`, and
  verified idempotence.
- `rtk proxy go test -count=1 -v -run '^(TestNoTwoBaselineClausesShareText|TestADuplicatedClauseTextIsReported|TestBaselineClauseForceIsCharacterized|TestTheBackendBoundaryParagraphRendersOnce|TestAMissingBaselineClauseIsReported|TestAnUnlistedBaselineClauseIsReported|TestAChangedBaselineClauseForceIsReported)$' ./internal/baseline`
  exited 0; all seven named tests passed.
- A focused Python comparison against Git checked all 16 module sources:
  only the specified backend removal and three version increments differ,
  and every surviving clause remains byte-identical.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Backend paragraph appears once | Regenerated formatter golden; `TestTheBackendBoundaryParagraphRendersOnce` passed, compared with two occurrences before the edit. |
| Repeated guidance is reported | `TestADuplicatedClauseTextIsReported` passed and checked the exact identifier pair after case and whitespace normalization. |
| Every clause retains its force; duplicate identifier absent | Literal 128-clause table matches the embedded catalog and `9e439dbb` minus the duplicate; characterization and all three negative force tests passed. |
| Second managed refresh is a no-op | Second text refresh exited 0 with `File changes: 0` and `Idempotence: verified`. |

Daemon-owned status is unchanged by this Agent. The declared Verification
command was not run; its execution and Task settlement remain with the Daemon.
No commit, push or Pull Request was created.

### Incremental check and follow-ups outside this slice

`GOCACHE=/private/tmp/roundfix-0193-task01-gocache rtk make verify-incremental`
exited 2 at its test target. Two existing tests still require the removed
duplicate and lie outside this Task's closed file list:

- `internal/baseline/plan_test.go`: `TestStandardTypeScriptStructuralClauseRetention`
  requires `rule.backend.boundary-contracts` to remain a retained current
  clause; it now reports that identifier unaccounted and absent.
- `internal/cli/baseline_update_test.go`: the `structural-clauses-missing`
  case of `TestBaselineUpdateFleetSweep` expects the boundary guidance twice;
  the refreshed guide now contains it once.

Focused reruns reproduced both conflicts without suiteguard violations:
`rtk proxy go test -count=1 -v -run '^TestStandardTypeScriptStructuralClauseRetention$' ./internal/baseline`
and
`rtk proxy go test -count=1 -v -run '^TestBaselineUpdateFleetSweep$/^structural-clauses-missing$' ./internal/cli`
each exited 1. Reconciling these existing structural-retention and fleet
expectations with the intentional duplicate removal requires a follow-up
outside this Task's authorization; neither test was edited.

The incremental run also encountered sandbox process-table restrictions in
two force-stop tests. An authorized elevated focused rerun,
`rtk proxy go test -count=1 -v -run '^(TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion)$' ./internal/cli`,
exited 0 with both tests passing. Writing this Result while the original
incremental run was active also triggered its CLI suiteguard; subsequent
focused reruns ran with no concurrent file edits and no guard violations.

Changed-file postflight found all 13 changed/untracked paths inside this
Task's Context or sanctioned outputs. The Task's authored content is
unchanged apart from the pre-existing Daemon status and this Result.
`rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

## Carry-forward provenance

- Source Run: `run_20260930T153457Z_b5736250c995b99a`
- Source commit: `3329a8af043f256a8fc7459dfa7caba0f7746f82`
