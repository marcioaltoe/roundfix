---
task: task_04
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: completed
type: backend
complexity: medium
---

# Task 04: Lifecycle wording is consistent and cites no repository record

## Overview

Four clauses of the `context-workflow` module need small changes. "Only `accepted` is active" contradicts the legacy rule beside it unless it is scoped to ADRs with lifecycle frontmatter. The Backlog contract lacks `deferred`, which twenty history entries use. Two clauses cite this repository's ADR numbers to every adopter. This Task applies the replacements the TechSpec fixes, teaches the checker that `deferred` is a terminal Backlog status, adds the check that keeps repository-only citations out of shipped guidance, and corrects one sentence of the repository's own guide.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST apply, in `internal/baseline/assets/modules/context-workflow.json`, the replacements the TechSpec's "Clause texts" gives for `clause.context.adr-02-active-status`, `clause.context.backlog-01-operational-contract`, `clause.context.docs-one-job-per-directory` and `clause.context.inbox-01-triage`. Every clause keeps its enforcement level.
2. MUST raise by one the versions of `rule.context.docs-layout`, `guide.docs-layout` and the `context-workflow` module.
3. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
4. MUST make `deferred` a terminal Backlog status in `internal/speccheck/backlog.go`, and make `internal/spec/retirement.go` treat a `deferred` entry as retired only when it carries a non-empty reason, exactly as a `declined` entry. A consuming Spec without a reason does not retire it. An active entry with `status: deferred` MUST be reported the way an active `declined` entry is, not as an unknown status.
5. MUST replace, in `docs/agents/specific-repository.md`, "Nothing under a `_archived` tree is ever validated." with "Nothing under `docs/history/` is ever validated as live work." and change nothing else in that file.
6. MUST add `internal/baseline/adopter_neutral_clauses_test.go`, `internal/speccheck/backlog_deferred_test.go` and `internal/spec/retirement_deferred_test.go` with the tests the TechSpec's Testing Approach 4 names. The citation check MUST read every clause-level and rule-level guidance string of every module, and MUST match `ADR-` followed by a digit and `Spec ` followed by four digits.
7. MUST leave every entry under `docs/history/backlog/` byte-identical, change no exported function signature, and rename or remove no top-level test.

## Subtasks

- [ ] Apply the four clause changes and raise the versions.
- [ ] Regenerate the pins, the golden, this repository's guide and the Setup Manifest.
- [ ] Make `deferred` terminal in the checker and the retirement reader.
- [ ] Correct the repository guide's archive sentence.
- [ ] Add the citation, wording and status tests, each negative case separate.

## Acceptance Criteria

- [ ] `docs/agents/docs-layout.md` scopes "only `accepted` is active" to ADRs with lifecycle frontmatter and lists `deferred` in the Backlog contract.
- [ ] No module guidance cites a Spec number or an ADR number, and a guidance string that cites one is reported.
- [ ] A `deferred` Backlog Entry is terminal: with a reason it is retired, without one it is not, even when it names a consuming Spec, and left in `docs/backlog/` it is reported.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/specific-repository.md`
- interface: `internal/speccheck/backlog.go`
- interface: `internal/spec/retirement.go`
- creates: `internal/baseline/adopter_neutral_clauses_test.go`
- creates: `internal/speccheck/backlog_deferred_test.go`
- creates: `internal/spec/retirement_deferred_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestShippedGuidanceCitesNoRepositoryRecord|TestARepositoryRecordCitationIsReported|TestLifecycleClausesCarryTheScopedWording|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestADeferredBacklogEntryIsTerminal|TestADeferredBacklogEntryLeftActiveIsReported|TestADeferredBacklogEntryWithAReasonIsRetired)$" ./internal/baseline ./internal/speccheck ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestShippedGuidanceCitesNoRepositoryRecord TestARepositoryRecordCitationIsReported TestLifecycleClausesCarryTheScopedWording TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestADeferredBacklogEntryIsTerminal TestADeferredBacklogEntryLeftActiveIsReported TestADeferredBacklogEntryWithAReasonIsRetired; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/docs-layout.md|For an ADR that carries lifecycle frontmatter, only" "docs/agents/docs-layout.md|status: deferred" "docs/agents/docs-layout.md|Preserve the boundary between evidence and intent" "docs/agents/specific-repository.md|is ever validated as live work"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && for pair in "docs/agents/docs-layout.md|ADR-0092" "docs/agents/docs-layout.md|ADR-0163"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the six new named tests exists and the guide still cites two ADR numbers, so the command fails.

## References

- [_techspec.md](_techspec.md) — Clause texts; Interfaces; Testing Approach 4
- `_prd.md` → Goals 3 and 4; Core Feature 4; Success Metrics 4 and 6
- `_techspec.md` → API Contracts 3 and 4; Testing Approach 5
- [references/2026-09-30-baseline-guides-contradict-the-shipped-product.md](references/2026-09-30-baseline-guides-contradict-the-shipped-product.md)
- ADR-0092, ADR-0163, ADR-0186

## Result

Implemented the Task 04 slice; Task status and declared Verification remain
Daemon-owned. No commit, push, Pull Request, Task Graph edit or other Task-file
edit was performed.

### Implementation and acceptance evidence

- ADR lifecycle guidance now scopes accepted-only activity to records carrying
  lifecycle frontmatter; the Backlog template and terminal list include
  `deferred`. The four specified clauses retain identifiers and enforcement.
  Module version rose 19 → 20, docs-layout rule 14 → 15 and docs-layout guide
  13 → 14, using in-place replacements.
- The citation regression reads every module through the existing embedded
  catalog helper, including clause-level and rule-level guidance. Its detector
  matches `ADR-[0-9]+|Spec [0-9]{4}` and reports the owning identifier and
  citation. Separate negative tests exercise an ADR clause citation and a Spec
  rule citation. The source, formatter golden and repository guide pass the
  scoped-wording regression after sanctioned regeneration.
- `terminalBacklogStatus` recognizes deferred. A deferred entry left active
  produces `SC-BACKLOG-UNMOVED`, error severity, its status-line location and
  the history destination. Retirement requires a non-empty reason; missing,
  null, blank and literal-null reasons remain live, including when a consuming
  Spec is named. Existing declined behavior remains unchanged, as required by
  the existing retirement characterization tests.
- Sanctioned regeneration rewrote only the matching formatter golden, profile
  digest, catalog snapshots and four plan goldens. The public managed refresh
  updated only `docs/agents/docs-layout.md` and `docs/agents/setup-context.json`.
  Its second invocation exited 0, reported `File changes: 0` and verified
  idempotence.
- A SHA-256 comparison confirmed all 29 files under `docs/history/backlog/`
  byte-identical. A comparison against HEAD confirmed the repository-specific
  guide differs by exactly the requested archive sentence. The changed-path
  postflight found 19 paths, all in this Task's Context or this Task file.
  No exported signature or existing top-level test was changed.

### Focused checks

Go commands used `GOCACHE=/tmp/roundfix-task04-gocache` after the default cache
was denied by the sandbox.

- Before implementation, `go test -count=1 -run
  'Test.*(Deferred|RepositoryRecord|SpecCitation|LifecycleClauses)'
  ./internal/baseline ./internal/speccheck ./internal/spec` reproduced both
  repository ADR citations, missing scoped wording, absent terminal deferred
  classification, absent active-entry finding and missing deferred retirement.
- After regeneration and refresh, `go test -count=1 -run
  'Test.*(Deferred|RepositoryRecord|SpecCitation|LifecycleClauses|ClassifyBacklogEntry)'
  ./internal/baseline ./internal/speccheck ./internal/spec` exited 0 for all
  three packages, including the pre-existing retirement characterizations.
- `make baseline-digests` exited 0 and reported `ok: true`, `changed: true`;
  its sanctioned regeneration checks passed.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills
  --yes --format text` first hit a sandbox denial opening the Git-private
  transaction lock. The approved elevated retry exited 0 with verified
  postimages and two file changes. The second normal invocation exited 0
  with zero changes and verified idempotence. Both reported existing nested
  instruction-carrier warnings; skills were skipped as requested.
- `git diff --check` exited 0. History-byte, repository-guide and path-scope
  assertions passed.
- `make verify-incremental` exited 2. Formatting and vet passed; the test tier
  exposed the out-of-scope stale assertions listed below and two sandbox
  process-table denials. The two denied tests passed on the approved elevated
  rerun: `go test -count=1 -run
  '^TestRunForceStop(OwnerProcessIntegrationProvesExitBeforeStoreCompletion|LegacyRunWithoutOwnerIdentityStillStopsOwner)$'
  ./internal/cli` exited 0. Subsequent incremental targets were not reached.
- The Task's declared `## Verification` command was not run.

### Follow-up outside this Task's authorized paths

The incremental check requires these expectation updates in a separately
bounded follow-up; those files were left unchanged:

- `internal/baseline/plan_test.go`, `TestStandardTypeScriptStructuralClauseRetention`:
  still requires `rule.backend.boundary-contracts`, the duplicate removed by
  Task 01.
- `internal/cli/baseline_update_test.go`,
  `TestBaselineUpdateFleetSweep/structural-clauses-missing`: still expects two
  backend boundary paragraphs rather than one after Task 01.
- `internal/baseline/plan_test.go`, `TestADRLifecycleContract`: still requires
  the old capitalized accepted-only sentence; Task 04's specified scoped
  replacement changes that sentence.

No terminal Task verdict or passing incremental/Daemon Verification claim is
made by this handoff.

## Carry-forward provenance

- Source Run: `run_20260930T153457Z_b5736250c995b99a`
- Source commit: `24da45da4b4a45ba54573db46221a33e184f214e`
