---
task: task_04
spec: 0253-authoring-rules-that-stop-qa-reruns
status: pending
type: backend
complexity: medium
---

# Task 04: A retirement writes reduced history and deletes reviews and handoffs

## Overview

The docs-layout clauses still tell a retirement to move whole Review Artifacts
and handoffs into `docs/history/reviews/` and `docs/history/handoffs/` and to
move a whole Finding or Backlog Entry with a dated addendum, while ADR-0248's
sanitize reduces or deletes those families and a full entry written after the
History Full Tag is refused (finding B16 of the Baseline audit of 2026-10-08,
`docs/references/2026-10-08-baseline-audit.md`). The operator decided on
2026-10-08 that retirements write history directly in reduced form and that
retired reviews and handoffs are not moved into `docs/history` at all. This
Task rewords three `context-workflow` clauses in place with `_techspec.md` →
Clause changes, follows in the user guide and revises **Reduced History
Entry** (ADR-0258). It is verifiable on its own: an entry written in the
guide's form is one the History Sanitize Command plans no change for.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer finding B16 of the Baseline audit of 2026-10-08 by editing the
   `guidance` of `clause.context.docs-one-job-per-directory`,
   `clause.context.backlog-01-operational-contract` and
   `clause.context.findings-09-archive` in
   `internal/baseline/assets/modules/context-workflow.json` exactly as
   `_techspec.md` → Clause changes (task_04) says. Each keeps its `id`, its
   `mandatory` enforcement and no `replaces` list. The Archive Record sentence
   and the history sanitize sentence of the first clause MUST stay
   byte-identical, and the rendered guide MUST still name every directory
   `spec.ArchiveDir` returns.
2. MUST raise by one, from the value on this Task's starting commit, the
   task_04 versions of `_techspec.md` → Version changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
3. MUST create `internal/baseline/reduced_retirement_clauses_test.go` with
   `TestTheRetirementClausesCarryTheirText`,
   `TestTheRetirementClausesRenderInTheDocsLayoutGuide`,
   `TestAFindingInTheGuideReducedFormNeedsNoSanitize` and
   `TestAnAdopterRetainsTheRetirementClauses`. The text test MUST require each
   new text and the absence of "a byte-identical move for Review Artifacts and
   handoffs", "A finished orphan Review Artifact retires to" followed by the
   reviews path, and "record the true disposition and its reason in a dated
   addendum". The rendering test MUST check the Standard TypeScript golden and
   this repository's `docs/agents/docs-layout.md` on whitespace-normalized
   text. The form test MUST read the fenced provenance line from the embedded
   `clause.context.findings-09-archive`, fill its placeholders with a 40-hex
   commit and a `docs/findings/` path, and write front matter, a title, a first
   paragraph and that line under a temporary `docs/history/findings/`. It MUST
   require `spec.IsReducedHistoryEntry` to accept the entry and
   `spec.PlanHistoryKind` for `spec.ArchiveKindFinding` to plan no file, and
   the same entry without the line to be planned. The retention test MUST
   require a ready Managed Refresh of a Source Baseline adopter that records
   the three clauses `retained` and no clause `unaccounted`.
4. MUST change the paragraph of `docs/user-guide/context-driven-development.md`
   named in `_techspec.md` → Skill and document texts, and nothing else in
   that file.
5. MUST revise the **Reduced History Entry** entry of `CONTEXT.md` through
   `domain-modeling` to the text of `_techspec.md` → Skill and document texts,
   keeping its `_Avoid_` line.
6. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   clause text and one of the provenance line in the clause, with the test
   that failed for each, and that the source was restored and regenerated.
7. MUST NOT change any other clause, the Source Baseline assets, the retention
   transition, History Relocation, `roundfix history sanitize` or any
   production Go file.

## Subtasks

- [ ] Reword the three clauses and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Add the clause test file with the reduced-form test.
- [ ] Change the user guide paragraph and **Reduced History Entry**.
- [ ] Record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/docs-layout.md` deletes a finished Review Artifact or a
      confirmed handoff instead of moving it, and writes a retired Finding or
      Backlog Entry as a reduced entry.
- [ ] An entry in the guide's form needs no sanitize.
- [ ] `CONTEXT.md` carries the revised **Reduced History Entry**.
- [ ] A Source Baseline adopter's refresh retains the three clauses, and a
      second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/spec/history_entries.go`
- instruction: `internal/baseline/glossary_clauses_test.go`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`
- creates: `internal/baseline/reduced_retirement_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheRetirementClausesCarryTheirText|TestTheRetirementClausesRenderInTheDocsLayoutGuide|TestAFindingInTheGuideReducedFormNeedsNoSanitize|TestAnAdopterRetainsTheRetirementClauses|TestHistorySanitizeClauseIsAppended|TestArchiveRecordClausesAreAppended|TestLifecycleClausesCarryTheScopedWording|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRetirementClausesCarryTheirText TestTheRetirementClausesRenderInTheDocsLayoutGuide TestAFindingInTheGuideReducedFormNeedsNoSanitize TestAnAdopterRetainsTheRetirementClauses TestHistorySanitizeClauseIsAppended TestArchiveRecordClausesAreAppended TestLifecycleClausesCarryTheScopedWording TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^TestDocsLayoutGuideNamesEveryHistoryFamily$' ./internal/docscontract || exit 1` — expected: exit 0; before this Task the four retirement tests do not exist, so the command fails.
- `for phrase in "A finished orphan Review Artifact retires by deletion and never enters history" "every handoff is deleted together and never enters history" "naming a commit that already holds the full text at its active path"; do tr -s '[:space:]' ' ' < docs/agents/docs-layout.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/docs-layout.md: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "Delete a finished Review Artifact or a confirmed handoff instead of moving it into history" || { printf 'the user guide lacks the retirement form\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Reduced History Entry**: A retired Finding or Backlog Entry cut to its front matter, title, first paragraph and the revision holding its full text. A retirement writes it directly" || { printf 'CONTEXT.md lacks the revised **Reduced History Entry**\n' >&2; exit 1; }; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the docs-layout guide moves a finished Review Artifact into history, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goal 6; User Story 6; Core Features 6 and 7; Success Metric 4
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; API Contract 4; Clause changes; Skill and document texts; Version changes; Retention; Derived files; Invariants 1-8; Testing Approach; Build Order 4
- ADR-0258; ADR-0248; ADR-0254; ADR-0244; ADR-0250
