---
task: task_01
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: pending
type: infra
complexity: medium
---

# Task 01: The archive clause sends knowledge to the project that owns the Secondbrain

## Overview

The Baseline's archive clause still sends a file a later reader needs to the
Secondbrain inbox, which ADR-0330 retired. This Task gives
`clause.spec.keep-artifacts-in-spec-folder` its current destination, records
the `spec-workflow` module's next version, regenerates the derived Baseline
files and this repository's managed guide with their commands. It answers the message reported by onioncry
([references](references/2026-11-06-the-archive-clause-names-the-retired-secondbrain-inbox.md)).

## Requirements

1. MUST change the `guidance` of `clause.spec.keep-artifacts-in-spec-folder`
   in `internal/baseline/assets/modules/spec-workflow.json` exactly as
   `_techspec.md` → API Contract 1 gives, keeping every other byte of the
   clause, and raise the `version` of `rule.spec.docs-layout` by one over the
   value the tree holds. The clause keeps its `id`, `enforcement` and the
   absence of `replaces`; nothing else in the module is edited by hand.
2. MUST update `internal/baseline/archive_record_clause_test.go`:
   - the `specSentence` constant takes the new sentence of API Contract 1;
   - `TestArchiveRecordClausesAreAppended` also fails when the embedded
     clause, the Standard TypeScript golden or `docs/agents/docs-layout.md`
     contains `Secondbrain inbox`;
   - a new `TestTheArchiveClauseIsRetainedFromTheSourceBaseline`, calling
     `t.Parallel()` first, runs `classifySourceClauseTransition` from Source
     Baseline `baseline.standard-typescript-monorepo-0.0.1` over the Standard
     TypeScript profile, as `cross_project_continuity_test.go` does, and
     asserts that `clause.spec.keep-artifacts-in-spec-folder` is recorded
     once with disposition `retained`, enforcement `mandatory` and its own
     identity as the only target.
3. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`;
   - `make baseline-digests`, which rewrites the Standard TypeScript docs
     layout golden, the Standard TypeScript profile, the three catalog files
     and the four plan characterization goldens this Task declares;
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     which rewrites `docs/agents/docs-layout.md` and
     `docs/agents/setup-context.json`, then the same command again, which
     must report `File changes: 0`. Run both with `HOME` set to a temporary
     directory removed afterwards, and `GOMODCACHE` and `GOCACHE` kept at
     their `go env` values.
   A file the commands change that this Task does not declare is a
   stop-and-report condition.
4. MUST NOT change the Source Baseline corpus, manifest, accounting or
   index, `internal/baseline/assets/retention/**`, any other module, clause,
   profile or guide, any ADR, any skill, `CONTEXT.md`, `.roundfixrc.yml`, the
   Makefile, the CI workflows or any Go production file.

## Subtasks

- [ ] Edit the clause and raise the rule version.
- [ ] Update the clause test and add the retention test.
- [ ] Record the module version, run `make baseline-digests` and the Managed Refresh twice.

## Acceptance Criteria

- [ ] The embedded clause, the Standard TypeScript golden and this repository's docs layout guide end with the ai-memory destination and never name the Secondbrain inbox.
- [ ] The clause is retained from the Source Baseline under its own identity, and a second Managed Refresh changes nothing.
- [ ] The Baseline package and the Setup Manifest contract pass.

## Context

- instruction: `docs/adr/0372-the-baseline-names-the-archive-destination-a-profiles-missing-decision-and-a-guide-it-does-not-own.md`
- instruction: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- instruction: `internal/baseline/plan.go`
- instruction: `internal/baseline/cross_project_continuity_test.go`
- instruction: `internal/baseline/module_versions_test.go`
- instruction: `internal/baseline/cancel_clause_test.go`
- instruction: `internal/baseline/qa_override_clause_test.go`
- instruction: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/archive_record_clause_test.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The archive clause text, the module version and the catalog digests — consumers: `./internal/baseline`, `./internal/docscontract`

## Verification

- `for f in docs/agents/docs-layout.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'or send its knowledge in an ai-memory message to the project that owns the Secondbrain for its ingestion.' || { printf 'missing destination in %s\n' "$f" >&2; exit 1; }; if grep -q 'Secondbrain inbox' "$f"; then printf 'retired inbox still named in %s\n' "$f" >&2; exit 1; fi; done; grep -q 'func TestTheArchiveClauseIsRetainedFromTheSourceBaseline(' internal/baseline/archive_record_clause_test.go || exit 1; go test -timeout 20m -count=1 -run '^(TestArchiveRecordClausesAreAppended|TestTheArchiveClauseIsRetainedFromTheSourceBaseline|TestEveryBaselineModuleVersionIsRecorded|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized)$' ./internal/baseline && go test -timeout 20m -count=1 ./internal/baseline && go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'` — expected: exit 0. Before this Task both rendered guides name the inbox, so the first check fails; afterwards they carry the destination, the retention test exists, the clause, record and whole Baseline suites pass and the Setup Manifest matches the Baseline.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 1; Testing Approach; Build Order 1
- ADR-0372, ADR-0330
