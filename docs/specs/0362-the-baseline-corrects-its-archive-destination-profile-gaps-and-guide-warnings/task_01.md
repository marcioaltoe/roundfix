---
task: task_01
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, no-workarounds, setup-context-driven, roundfix, tech-writer, evidence-gate.

The archive clause now uses API Contract 1's exact ai-memory destination.
Its identity, mandatory enforcement and absence of `replaces` are preserved.
`rule.spec.docs-layout` rose from 8 to 9; the sanctioned recorder raised the
`spec-workflow` module from 23 to 24 and appended its version record.
Every other byte of the module is unchanged, as checked against `HEAD`.
Generated files were written only by their owning commands, in the required
order. Initial status inspection found only this Task's Daemon-owned status
change; that status and the authored Task body remain unchanged by the Agent.

Acceptance evidence:

1. The embedded clause and both rendered docs layout guides carry the new
   sentence and reject `Secondbrain inbox`. The updated
   `TestArchiveRecordClausesAreAppended` first failed with six diagnostics
   covering the retired destination and missing replacement in all three
   locations, then passed after source editing and regeneration.
2. `TestTheArchiveClauseIsRetainedFromTheSourceBaseline` starts with
   `t.Parallel()`, resolves the Standard TypeScript profile and classifies
   Source Baseline `baseline.standard-typescript-monorepo-0.0.1`. It asserts
   the clause's `retained` disposition and exactly one evidence entry with
   mandatory enforcement and only its own identity as target. The focused
   run passed. The second Managed Refresh reported `File changes: 0` and
   `Idempotence: verified`; hashes of every declared interface and this Task
   were identical before and after that invocation.
3. The focused Baseline tests and sanctioned regeneration checks passed.
   The complete Baseline package and Setup Manifest docscontract command
   are reserved for Daemon Verification and were not run in this turn.
   Managed Refresh itself verified its approved postimages and profile
   alignment, and regenerated the Setup Manifest's catalog and guide digests.

Commands and focused outcomes:

- `go test ./internal/baseline -run '^TestArchiveRecordClausesAreAppended$' -count=1`
  — exit 1 before changing the module, as expected for the new contract.
- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
  — exit 0; only the declared module and module version record changed.
- `make baseline-digests` — exit 0; reported `ok: true`, `changed: true`.
  It rewrote exactly the nine declared derived outputs: the Standard
  TypeScript docs layout golden and profile, three catalog files and four
  plan characterization goldens.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
  — executed twice, each at exit 0, with one temporary `HOME` removed in
  cleanup. `GOMODCACHE` and `GOCACHE` retained their original `go env` values.
  The first invocation reported two file changes, limited to
  `docs/agents/docs-layout.md` and `docs/agents/setup-context.json`, with
  approved Plan Digest
  `sha256:282b78793acfc0c707a9893d81015a9d293c5b830ca7276abbf2f712afdc9426`.
  The second reported zero file changes. Both reported `Retention evidence: 0`
  and skipped skills and history as requested. Both retained the existing
  warnings for the formatter-fixture and Source Baseline nested `AGENTS.md`
  carriers, leaving those carriers untouched.
- `go test ./internal/baseline -run '^(TestArchiveRecordClausesAreAppended|TestTheArchiveClauseIsRetainedFromTheSourceBaseline|TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestTheCancelClauseCarriesItsText|TestTheCancelClauseRendersInTheDocsLayoutGuide|TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand)$' -count=1`
  — exit 0, covering the changed text, retention evidence and neighboring
  cross-project, cancellation and QA override contracts.
- Byte and JSON comparisons against `HEAD` confirmed that the source module
  has only the exact sentence replacement, rule version increment and
  recorder-owned module version, and that the module record only appends
  the new `spec-workflow` entry. Scope checks after regeneration and each
  refresh found no undeclared path; all 15 changed paths, including this
  Task, are declared. `git diff --check` exited 0.

The shell runner could not start with its configured symlinked writable
root. Approved shell execution used the physical checkout for all work;
no repository configuration was changed to address that environment limit.
No other Task, Task Graph, Source Baseline, retention corpus, ADR, skill,
production Go source or tooling configuration was changed. No commit,
push or pull request was made. No additional follow-up was identified.

## Carry-forward provenance

- Source Run: `run_20261106T174350Z_09865ccb421258be`
- Source commit: `5c2f91c2ac7134852ed944db4ff8d0afbc72de9f`
