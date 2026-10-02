---
task: task_01
spec: 0216-baseline-wording-left-after-the-stack-wave
status: completed
type: backend
complexity: medium
---

# Task 01: The backend bucket clause names what it forbids

## Overview

The backend guide's clause "Do not introduce generic `modules` or `services` buckets as the normative backend architecture." has no scope and reads as a ban on domain services. This Task replaces it with a clause under a new identity that forbids generic buckets in place of the domain, application and infrastructure layers and says that a domain service in the domain layer is not one, declares the replacement so a Source Baseline adopter's update records the old clause `replaced`, and gives the new clause its Source Baseline row. It answers the Backlog Entry "Baseline wording left after the stack wave" of 2026-09-30 (item: the generic-layers prohibition).

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST replace, in `rule.backend.boundary-contracts` of `internal/baseline/assets/modules/backend.json`, the clause `clause.backend.prohibit-generic-layers` with `clause.backend.prohibit-generic-buckets`, enforcement `prohibited`, `"replaces": ["clause.backend.prohibit-generic-layers"]`, and the guidance the TechSpec's "Exact texts" gives for task_01, editing the object in place.
2. MUST add the Source Baseline row the TechSpec's "Source Baseline rows" table gives for task 01 (corpus entry, `manifest.json` row with placeholder offsets and digest, and `index.json` `entryIds` entry, each after the old clause), and MUST keep the old clause's row.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_01.
4. MUST update the existing tests the TechSpec's "Existing tests that change" names for task_01: the force record swaps the clause key; the maintained Source Baseline entry count rises by one; `TestStandardTypeScriptStructuralClauseRetention` expects the old clause `replaced` by the new one and compares the new clause's guidance with its own Source Baseline row, through one map of replacements that also holds `rule.backend.boundary-contracts`; and the fleet structural fixture in `internal/cli/baseline_update_test.go` names the new identity and rendered line.
5. MUST create `internal/baseline/backend_bucket_clause_test.go` with `TestTheBackendGuideScopesTheBucketProhibition` and `TestTheBucketClauseReplacesTheClauseAdoptersHold` as the TechSpec's Testing Approach 1 describes; the second MUST report `unaccounted` when `replaces` is removed.
6. MUST run `make baseline-digests` twice (the second reports `"changed":false`), then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second refresh MUST report `File changes: 0`. MUST NOT hand-edit a snapshot, golden, digest or Source Baseline offset.
7. MUST NOT change any other clause's identifier, enforcement or guidance, any Skill Activation, or any skill, and MUST NOT rename or remove a top-level test or change an exported function signature.

## Subtasks

- [ ] Replace the clause and declare the replacement.
- [ ] Add the Source Baseline row and raise versions.
- [ ] Update the four existing tests and create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] A Standard TypeScript Monorepo plan's backend guide states the scoped clause with `prohibited` and not the old sentence.
- [ ] The Source Baseline transition records `clause.backend.prohibit-generic-layers` `replaced` by `clause.backend.prohibit-generic-buckets`, and the new clause `retained`.
- [ ] A second regeneration and a second Managed Refresh change nothing.

## Context

- instruction: `docs/adr/0222-a-baseline-guide-says-only-what-holds-for-the-repository-that-reads-it.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/backend.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- creates: `internal/baseline/backend_bucket_clause_test.go`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/cli/baseline_update_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheBackendGuideScopesTheBucketProhibition|TestTheBucketClauseReplacesTheClauseAdoptersHold|TestStandardTypeScriptStructuralClauseRetention|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestReadoptionCompatibilityMaintainedFixture|TestCatalogCompatibility|TestBaselinePlanCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBackendGuideScopesTheBucketProhibition TestTheBucketClauseReplacesTheClauseAdoptersHold TestStandardTypeScriptStructuralClauseRetention TestReadoptionCompatibilityMaintainedFixture; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^TestBaselineUpdateFleetSweep$' ./internal/cli && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.

## References

- `_prd.md` → Goal 2; Core Feature 1; Success Metric 2; Success Metric 4; Declared breaks
- `_techspec.md` → Exact texts; Source Baseline rows; Version changes; Existing tests that change; API Contract 1; API Contract 4; Testing Approach 1; Build Order 1
- ADR-0058, ADR-0060, ADR-0099, ADR-0190, ADR-0202, ADR-0222

## Result

Implemented the task_01 slice for Daemon Verification. Task status remains
Daemon-owned; no declared Verification command was run, and no commit, push,
or Pull Request was created.

- Replaced the bucket clause in place with
  `clause.backend.prohibit-generic-buckets`, force `prohibited`, and the
  declared `replaces` identity. The text forbids generic buckets in place of
  domain, application, and infrastructure layers and explicitly permits a
  domain service in the domain layer.
- Kept the old Source Baseline row and added the successor immediately after
  it in the corpus, manifest, and index. Sanctioned regeneration supplied all
  offsets and digests. The maintained entry count is now 162.
- Raised `backend` and `guide.backend` from 5 to 6 and
  `rule.backend.boundary-contracts` from 4 to 5. Updated the force record,
  structural-retention test through one replacement map, maintained entry
  count, and fleet structural fixture. Added both authored regression tests.

### Acceptance evidence

1. **Scoped rendered prohibition:**
   `GOCACHE=/private/tmp/roundfix-0216-task01-cache rtk proxy go test ./internal/baseline -run 'TestThe(BackendGuideScopes|BucketClauseReplaces)|TestStandardTypeScriptStructuralClauseRetention' -count=1 -v`
   exited 0. `TestTheBackendGuideScopesTheBucketProhibition` checks the complete
   rendered prohibited line and absence of the old sentence. Before the source
   edits, both new tests failed against the original catalog.
2. **Source Baseline transition:** the same focused check passed
   `TestTheBucketClauseReplacesTheClauseAdoptersHold` and
   `TestStandardTypeScriptStructuralClauseRetention`. The old identity is
   `replaced` with the successor as its sole target; the successor is
   `retained`. The negative subtest removes `replaces` from a fresh in-memory
   catalog and requires the old identity to become `unaccounted`.
3. **Convergence:**
   `GOCACHE=/private/tmp/roundfix-0216-task01-cache rtk proxy make baseline-digests`
   ran twice and exited 0 each time; the second reported `"changed":false`.
   `GOCACHE=/private/tmp/roundfix-0216-task01-cache rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
   applied the Setup Manifest catalog digest on the first successful refresh
   and reported `File changes: 0` and verified idempotence on the second.
   The initial sandboxed apply could not create the Git-private transaction
   journal and changed no planned file; the authorized retry with expanded
   filesystem access succeeded. Both successful refreshes skipped skills and
   preserved nested instruction carriers, reported as inventory warnings.

### Additional focused checks

- Compared the backend module with `HEAD`: only the bucket clause changed;
  other module content and Skill Activations are identical, and the three
  versions each rose by one.
- `git -c core.fsmonitor=false diff --check` exited 0. Changed-path inspection
  found all 21 changed paths inside task_01's Context or this Task file; no
  Task Graph or other Task file changed.
- The initial incremental check exited 2: process-owner tests lacked sandbox
  process-table access, and suite guards detected concurrent edits made during
  the check. It also exposed shared-catalog mutation by the new negative test.
  That test now clones the cached catalog before removing `replaces`.
- After isolating the negative test, an expanded focused check exited 0:
  `GOCACHE=/private/tmp/roundfix-0216-task01-cache rtk proxy go test ./internal/baseline -run 'TestThe(BackendGuideScopes|BucketClauseReplaces)|TestStandardTypeScriptStructuralClauseRetention|TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor|TestTheWarningsClauseIsReplacedByTheLintClause|TestTheExternalTriageRuleIsReplacedForSourceBaselineAdopters|TestAnUndeclaredExternalTriageReplacementIsRefused' -count=1 -v`.
  All seven tests passed together, including the four affected by the shared
  mutation in the initial incremental run.
- `GOCACHE=/private/tmp/roundfix-0216-task01-cache rtk make verify-incremental`
  reran with required process-table access and no concurrent repository edits,
  and exited 0: formatting, vet, Go tests, skill synchronization and checks,
  and build passed. The local output is retained at
  `/private/tmp/roundfix-0216-task01-incremental-rerun.log`.

No follow-up implementation was added to this slice. Declared Verification
and Task settlement remain with the Daemon.
