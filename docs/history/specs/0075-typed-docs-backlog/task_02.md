---
task: task_02
spec: 0075-typed-docs-backlog
status: completed
type: test
complexity: medium
---

# Task 02: Re-record the corpus and declare what moved

## Overview

The corpus is a characterization gate, and re-recording it is sanctioned by
ADR-0081 as fallout of the authorized module edit. Sanctioned is not the same as
unexamined: **only layout content and its digests may move.** Anything else
moving means the edit leaked into a surface this Spec does not own.

This slice runs the regeneration and reads the diff, which is the whole point.

## Requirements

1. MUST run `make baseline-digests` and re-record the two characterization
   corpora it does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

2. MUST run `make baseline-digests` a second time after both re-records,
   because `BASELINE_DIGEST_STEPS` regenerates the maintained Source Baseline
   fixture before the catalog digests it derives from, so a single pass leaves
   that fixture stale.
3. MUST inspect the resulting diff and record, in the Task Result, every path
   that moved and why layout content explains it.
3. MUST fail the Task if any path moved that layout content does not explain.
   A leaked edit is the failure this Task exists to catch.
4. MUST leave `make verify` green after one regeneration pass.
5. MUST NOT hand-edit any derived artifact **except** to create the Source
   Baseline manifest rows for clauses task_01 introduced. The regenerator
   maintains those rows and never creates them, and says so itself:

   ```text
   catalog.sourceBaseline.required-clause.missing:
     baseline.standard-typescript-monorepo-0.0.1:
     clause.context.backlog-01-operational-contract;
     the regenerator maintains manifest rows but never creates them,
     so add this row first
   ```

   Creating the row is the only action that unblocks the regeneration, so a
   blanket prohibition forbids the one step this Task exists to perform.
   Everything else still comes from a regeneration command, and every created
   row MUST be listed in the Result with the clause that required it.

## Subtasks

- [ ] Run the sanctioned command and both corpus re-records.
- [ ] Read the full diff and classify every moved path.
- [ ] Record the declared break list in the Result.

## Acceptance Criteria

- [ ] `make verify` exits 0 after the regeneration pass.
- [ ] Every moved path is listed in the Result with the layout content that
      explains it.
- [ ] No path moved that layout content does not explain.
- [ ] No derived artifact was hand-edited, asserted by every change tracing to a
      regeneration command named in the Result.

## Verification

This Task's whole purpose is to return a repository that task_01 legitimately
left red. Declaring the configured repository Verification command verbatim
makes the Daemon run it as a **precondition**, before the Agent starts — and it
fails on exactly the state this Task exists to repair, settling the Task
without ever starting it. Measured on 2026-08-05:
`repository not green on entry: make verify exited 2`.

The gate below is the same gate, named by its parts rather than by the
configured string, so the repository is still proven green **after** the
regeneration and never demanded green before it. Nothing is weakened: every
target `make verify` runs is listed.

- `make baseline-digests` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'TestBaselinePlanCharacterization|TestCatalogDiagnosticCharacterization' -v | grep -q -- "--- PASS"`
  — expected: exit 0; both corpora match after re-recording.
- `make baseline-digests` — expected: exit 0. The second pass is load-bearing,
  not a retry: `BASELINE_DIGEST_STEPS` regenerates
  `TestReadoptionCompatibilityMaintainedFixture` as its **first** step, before
  the catalog digests that fixture derives from, so one pass leaves it stale
  against what the later steps just wrote. Measured on 2026-08-05:
  `TestReadoptionCompatibilityMaintainedFixture` failed the full gate after a
  single pass, with `maintained Source Baseline counts = identity 106 entries
  106 accounting 51`.
- `make fmt-check test spec-budget skills-sync-check skills-check build spec-check`
  — expected: exit 0; every target `make verify` runs, after the regeneration.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.

## References

- `_prd.md` → Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 2; Risks & Considerations.
- ADR-0081.

## Result

The maintained Source Baseline now carries task_01's two backlog clauses, and
the complete derived corpus has been re-recorded in two focused regeneration
waves. The repair adds no Go source and changes no surface outside layout
content, Source Baseline accounting, and their derived digests.

### Bootstrap

- Added the canonical marker-delimited Source Baseline content for
  `clause.context.backlog-01-operational-contract` and its manifest row.
- Added the canonical marker-delimited Source Baseline content for
  `clause.context.backlog-02-finding-boundary` and its manifest row.
- Updated the existing `clause.context.docs-one-job-per-directory` Source
  Baseline content with `docs/backlog/`, matching task_01's module guidance.
- The three Source Baseline entries are byte-identical to their corresponding
  `internal/baseline/assets/modules/context-workflow.json` guidance; focused
  `diff -u` comparisons exited 0 for all three.
- The manifest-row bootstrap used temporary span and digest values only until
  `TestReadoptionCompatibilityMaintainedFixture -update` calculated them from
  the source bytes. A focused search confirmed no temporary all-zero digest
  remains.

### Focused evidence

- The first Daemon Verification attempt stopped at
  `TestReadoptionCompatibilityMaintainedFixture` because the two required
  Source Baseline entries did not yet exist. The diagnostic named both clause
  IDs and the manifest-row bootstrap action.
- Two focused regeneration waves ran the individual `BASELINE_DIGEST_STEPS` in
  Makefile order without rerunning the declared `make baseline-digests`
  Verification command. Both waves exited 0 for
  `TestReadoptionCompatibilityMaintainedFixture -update`,
  `TestAuthorialSkillSync -update`, `TestFormatterComposition -update`,
  `TestBaselineCompatibilityCorpus -update`, `TestCatalogCompatibility
  -update`, `TestCatalogDiagnosticCharacterization
  -update-catalog-diagnostics`, and `TestBaselinePlanCharacterization
  -update-baseline-plan-characterization`.
- The two characterization re-record commands required by this Task each
  exited 0 in both waves.
- Focused strict check
  `go test ./internal/baseline -count=1 -run
  'TestReadoptionCompatibilityMaintainedFixture|TestCatalogCompatibility|TestFormatterComposition'`
  exited 0 after the last mutation.
- `rtk git diff --check` exited 0 after the last mutation.

### Declared break list

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
  — regenerated guide content adds the backlog operational contract, the
  finding boundary, and `docs/backlog/` to the directory-purpose clause.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — the
  regenerated docs-layout golden changes the selected formatter golden digest.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
  — the two clauses move the entry count from 104 to 106 and regenerate the
  corpus and manifest digests.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
  — canonical Source Baseline layout content adds the two marker-delimited
  backlog clauses and updates the directory-purpose clause.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
  — adds the two required clause rows; regeneration calculates their spans and
  digests and shifts later docs-layout spans without changing their digests.
- `internal/baseline/assets/source-baselines/index.json` — adds both clause IDs,
  moves the entry count to 106, and regenerates the Source Baseline corpus and
  manifest digests.
- `internal/baseline/testdata/catalog.diagnostics.golden.json` — the formatter
  golden-digest mismatch characterization follows the regenerated expected and
  actual layout digests.
- `internal/baseline/testdata/catalog.digest` — regenerated catalog digest for
  the changed module, formatter golden, profile pin, and Source Baseline.
- `internal/baseline/testdata/catalog.normalized.json` — normalized catalog
  records the changed layout bytes and digests, Source Baseline accounting, and
  profile pin.
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
  — re-records only catalog, layout-content, managed-file identity, and plan
  digest fallout.
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
  — re-records only catalog, layout-content, managed-file identity, and plan
  digest fallout.
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
  — re-records only catalog, layout-content, byte-count, managed-file identity,
  and plan digest fallout.
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
  — re-records only catalog, layout-content, byte-count, managed-file identity,
  and plan digest fallout.
- `docs/specs/0075-typed-docs-backlog/task_02.md` — the Daemon-owned status
  transition was pre-existing; this Task adds the required Result and declared
  break evidence.

### Acceptance-criterion evidence

- `make verify` was not run in this Agent repair turn; the Daemon owns the
  declared Verification sequence.
- Every moved path is listed above. Full diff inspection found only backlog
  layout content, Source Baseline rows/accounting, and their digest or
  characterization fallout.
- No path moved outside those declared categories; no Go source changed.
- The only bootstrap-authored Source Baseline material is the canonical layout
  content and the two manifest rows named above. Every span, digest, index,
  formatter fixture, catalog artifact, diagnostic fixture, and plan
  characterization change traces to the focused update tests named above.
