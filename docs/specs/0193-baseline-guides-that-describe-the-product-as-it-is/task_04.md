---
task: task_04
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
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
