---
task: task_01
spec: 0209-sources-that-share-a-context-share-a-spec
status: completed
type: backend
complexity: medium
---

# Task 01: Sources that share a context share a Spec, as Baseline clauses with force

## Overview

The Baseline does not tell an author to group sources that share a context,
to stop grouping at a size, or to extend an open record instead of minting a
new one. This Task adds the three mandatory clauses of `_techspec.md` → Exact
clause texts to the Spec workflow and CONTEXT workflow modules, gives each its
Source Baseline row, regenerates the derived files and refreshes this
repository's guides (ADR-0208). It answers the maintainer's requests of
2026-10-01 quoted in `_prd.md`. It is verifiable on its own: the rendered
guides state the clauses with force, and a Source Baseline adopter's plan
records each one `retained`.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST append to `rule.spec.routing` in `internal/baseline/assets/modules/spec-workflow.json`, after `clause.spec.verification-two-tiers`, the clauses `clause.spec.sources-01-group-by-shared-context` and `clause.spec.sources-02-bound-the-group`, and to `rule.context.docs-layout` in `internal/baseline/assets/modules/context-workflow.json`, after `clause.context.inbox-02-fleet-flow`, the clause `clause.context.inbox-03-extend-before-minting`, each `mandatory` with the exact guidance of `_techspec.md` → Exact clause texts. Every existing clause MUST stay byte-identical, and no clause gets a `replaces` list.
2. MUST add the three Source Baseline rows of `_techspec.md` → Source Baseline rows: corpus entry, `manifest.json` row and `index.json` `entryIds` identifier, each after its named anchor. No existing row changes.
3. MUST raise by one, from the value on this Task's starting commit, the six versions `_techspec.md` → Version changes lists, keeping each module's formatting.
4. MUST update the existing tests `_techspec.md` → Existing tests that change names for this Task: the force record gains the three clauses as `mandatory`, and the maintained Source Baseline entry count rises by three. No other existing test line changes.
5. MUST create `internal/baseline/grouped_sources_clauses_test.go` with `TestTheGroupingClausesCarryTheirForceAndText`, `TestTheGroupingClausesRenderInTheGuides`, `TestTheGroupingClausesHaveSourceBaselineRows` and `TestAnAdopterRetainsTheGroupingClauses`, following `_techspec.md` → Testing Approach 1. The texts are literals in the test, so a reworded clause fails it; the retention test MUST assert a ready plan and the `retained` disposition for each clause, and no `unaccounted` clause.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`; a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a snapshot or a generated guide.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of a clause text (for example dropping its last sentence) and one of a Source Baseline row (for example removing its manifest row), each with the test or check that failed, and that the source was restored.

## Subtasks

- [ ] Add the three clauses and raise the versions.
- [ ] Add the three Source Baseline rows.
- [ ] Update the two existing tests and create the new test file.
- [ ] Regenerate, refresh twice, and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/spec-routing.md` states both Spec-workflow clauses and `docs/agents/docs-layout.md` the extend clause, each labelled `mandatory`.
- [ ] The Standard TypeScript Monorepo goldens state the same clauses.
- [ ] A Source Baseline adopter's Managed Refresh plan is ready and records the three clauses `retained`.
- [ ] The force record, the duplicate-text check, the record-citation check and the catalog validation pass, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/grouped_sources_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestTheGroupingClausesHaveSourceBaselineRows|TestAnAdopterRetainsTheGroupingClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGroupingClausesCarryTheirForceAndText TestTheGroupingClausesRenderInTheGuides TestTheGroupingClausesHaveSourceBaselineRows TestAnAdopterRetainsTheGroupingClauses TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestShippedGuidanceCitesNoRepositoryRecord TestReadoptionCompatibilityMaintainedFixture TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in "docs/agents/spec-routing.md|- **mandatory**: One Spec may adopt several Inbox Entries, Backlog Entries, and Findings whose context is similar or complementary" "docs/agents/spec-routing.md|is advisory: adopt the suggested source or state why it stays apart." "docs/agents/spec-routing.md|- **mandatory**: Group sources into one Spec only while the Spec fits four implementation Tasks plus its QA gate." "docs/agents/docs-layout.md|- **mandatory**: Before minting a Finding or Backlog Entry, including in Triage" "docs/agents/docs-layout.md|extend a fitting Finding with a dated addendum, because a Finding is immutable history."; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || exit 1; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the four new tests do not exist and the guides lack the clauses, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1, 2, 3 and 4; User Stories 1, 2, 3 and 6; Core Features 1, 2, 3 and 4; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — Exact clause texts; Source Baseline rows; Version changes; Existing tests that change; API Contract 6; API Contract 7; Testing Approach 1; Build Order 1
- ADR-0208; ADR-0186; ADR-0058; ADR-0060; ADR-0081; ADR-0149

## Result

Implemented this Task's Baseline slice for Daemon Verification. Task status,
Task Graph, other Tasks, skills and judge implementation were left untouched.
The only pre-existing change was this Task's Daemon-owned `in_progress` status.

The Spec workflow now carries the two exact mandatory grouping clauses, and
CONTEXT workflow carries the exact mandatory extend-before-minting clause.
They have no `replaces` field. All existing clause objects were compared to
`HEAD` and remain byte-identical. The six versions rose from 13/5/9 to
14/6/10 (Spec module/rule/guide) and from 20/15/14 to 21/16/15
(CONTEXT module/rule/guide). Source Baseline corpus entries, manifest rows and
index identifiers follow the authored anchors; existing row content is
preserved, with shifted offsets derived by regeneration. The force record
adds only the three mandatory entries; the maintained entry count rises from
146 to 149. No other existing test line changes.

Four new tests lock literal clause text and force, rendered mandatory bullets,
Source Baseline identity/force/carrier, and an actual Source Baseline adopter's
Managed Refresh plan. The adopter uses an isolated temporary Git repository,
requires `ready`, asserts all three dispositions and retention evidence are
`retained`, and rejects every `unaccounted` disposition. No network or model
call is involved.

Focused evidence by acceptance criterion:

| Criterion | Evidence |
| --- | --- |
| Repository guides state the three mandatory clauses | Sanctioned Managed Refresh rewrote the two guides and Setup Manifest; exact-text inspection found each mandatory bullet in its guide. |
| Standard TypeScript Monorepo goldens state the clauses | `TestTheGroupingClausesRenderInTheGuides` passed for all three exact mandatory bullets after sanctioned regeneration. |
| Adopter plan is ready and retains the clauses | `TestAnAdopterRetainsTheGroupingClauses` passed for all three clauses, requiring a ready Managed Refresh and no unaccounted clause. |
| Force, duplicate text, record citations, catalog and no-op refresh | Focused checks below passed; second Managed Refresh exited 0 with `File changes: 0` and `Idempotence: verified`. |

Commands and outcomes (Go commands used `GOCACHE=/tmp/roundfix-0209-gocache`):

- Initial `go test -count=1 -run '^TestTheGroupingClausesCarryTheirForceAndText$' ./internal/baseline`: exit 1, named all three missing clauses before source edits.
- `make baseline-digests`: exit 0; sanctioned generation refreshed the goldens, Source Baseline identities/offsets/digests, profile pin, catalog snapshots and four plan-characterization goldens. No generated file or pin was hand-edited.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`: first sandboxed attempt exited 1 because its Git-private transaction directory was outside the writable root. The authorized rerun with filesystem escalation exited 0, `Baseline update: verified`, three file changes. Second refresh exited 0, `File changes: 0`. Existing nested fixture-carrier warnings leave those carriers unchanged.
- `go test -count=1 -v -run '^(TestTheGroupingClauses.*|TestAnAdopterRetainsTheGroupingClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestCatalogCompatibility)$' ./internal/baseline`: exit 0, all nine selected tests passed, including every new clause subtest.
- Source/scope inspection compared existing clause bytes and all six versions with `HEAD`, checked the repository guides and goldens for exact mandatory text, and checked tracked plus untracked changed paths against this Task's Context: 24 paths, all authorized. `git diff --check`: exit 0.

Negative evidence, with sources restored:

1. Removed the last sentence of `clause.spec.sources-01-group-by-shared-context` from the module. `go test -count=1 -run '^TestTheGroupingClausesCarryTheirForceAndText$' ./internal/baseline` exited 1 with `force/text differs` for that clause. Restored the module byte-for-byte.
2. Removed that clause's manifest row. `go test -count=1 -run '^(TestTheGroupingClausesHaveSourceBaselineRows|TestAnAdopterRetainsTheGroupingClauses)$' ./internal/baseline` exited 1 for both tests with `catalog.sourceBaseline.integrity.invalid` (counts or digests disagree). Restored the manifest byte-for-byte.
3. Removed the same sentence from the module again and ran `make baseline-digests` to derive the sabotaged guide, without editing a golden. `go test -count=1 -run '^TestTheGroupingClausesRenderInTheGuides$' ./internal/baseline` exited 1 with `guide lacks exactly one forced clause`. Restored the module byte-for-byte and ran `make baseline-digests` again, exit 0, restoring its derived outputs.

The focused checks and exact-text/scope inspection were repeated after the
sabotages were restored. The Task's declared Verification command and broader
repository Verification were not run in this child turn; settlement remains
Daemon-owned. No commit, push or Pull Request was made. No follow-up was
required within this slice.

## Carry-forward provenance

- Source Run: `run_20261001T233647Z_36941fe094803f4f`
- Source commit: `8b3138f5b29c6f42c270f4de85a8ea178f0ef3a3`
