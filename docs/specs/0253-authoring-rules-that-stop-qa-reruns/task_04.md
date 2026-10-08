---
task: task_04
spec: 0253-authoring-rules-that-stop-qa-reruns
status: completed
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


## Result

Implemented the task_04 retirement guidance slice for Daemon Verification.
Task status and the authored Verification remain Daemon-owned.

### Implementation and acceptance evidence

- Retirement behavior: reworded only the three declared `context-workflow`
  clauses, preserving their IDs, `mandatory` enforcement and lack of
  `replaces`. The generated docs-layout guide now deletes a finished orphan
  Review Artifact and all explicitly confirmed handoffs, and retires Findings,
  Rollups and Backlog Entries as reduced entries. The Archive Record and
  history sanitize sentences remain byte-identical to the starting commit.
  A focused source comparison confirmed that no other clause changed and
  that the rendered guide still names all six ArchiveDir history families.
  `TestTheRetirementClausesCarryTheirText` and
  `TestTheRetirementClausesRenderInTheDocsLayoutGuide` passed.
- Reduced form: `TestAFindingInTheGuideReducedFormNeedsNoSanitize` reads the
  fenced line from the embedded Findings archive clause, substitutes a
  40-hex revision and an active `docs/findings/` path, and writes a temporary
  entry with front matter, title and first paragraph. The guide form is
  recognized by `spec.IsReducedHistoryEntry` and plans no sanitize file;
  the same entry without provenance plans exactly that file for reduction.
  Both subtests passed.
- Domain vocabulary: revised only the Reduced History Entry definition in
  `CONTEXT.md` through `domain-modeling`, using the exact TechSpec text;
  its `_Avoid_` line is unchanged. Updated only the named retirement
  paragraph of the user guide. A focused comparison with the starting
  commit checked both document boundaries and their normalized text.
- Retention and refresh: `TestAnAdopterRetainsTheRetirementClauses` passed
  against the Source Baseline adopter helper: Managed Refresh is ready,
  all three clauses have retained dispositions and retention evidence,
  and no clause is unaccounted. This repository's second refresh exited 0,
  reported `File changes: 0` and `Idempotence: verified`.

### Commands and focused checks

All Go commands used `GOCACHE=/tmp/roundfix-task04-gocache`.

- Starting-contract check:
  `go test ./internal/baseline -run '^(TestTheRetirementClausesCarryTheirText|TestAFindingInTheGuideReducedFormNeedsNoSanitize)$' -count=1`
  exited 1 against the starting module: obsolete retirement wording and
  missing fenced provenance form were detected. The proposed source was
  then restored before the record step.
- Raised `rule.context.docs-layout` from 16 to 17 and `guide.docs-layout`
  from 15 to 16, each one above the starting commit. Ran
  `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`:
  exit 0; the record command chose module version 25 (previously 24).
- `make baseline-digests`: exit 0, regenerated the nine derived files listed
  below. After restoring the sabotage edits, ran it again: exit 0,
  `changed:false`, with no derived drift.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`:
  the sandboxed attempt exited 1 before apply because the Git-private
  transaction directory was not writable. Reran the same authorized command
  with the required filesystem access: exit 0, approved postimages verified,
  two file changes. The second refresh exited 0 with zero file changes and
  verified idempotence. The existing nested-carrier warnings remained;
  Repository Verification was not run by Baseline.
- Final focused check:
  `go test ./internal/baseline -run '^(TestTheRetirementClausesCarryTheirText|TestTheRetirementClausesRenderInTheDocsLayoutGuide|TestAFindingInTheGuideReducedFormNeedsNoSanitize|TestAnAdopterRetainsTheRetirementClauses)$' -count=1 -v`:
  exit 0, all four named tests and their subtests passed after restoration
  and regeneration.
- `git diff --check` and the focused source/document/changed-path postflight:
  exit 0; every changed path belongs to this Task's declared slice.

### Sabotage evidence

1. Replaced the source clause's sentence prefix
   `A finished orphan Review Artifact retires by deletion and never enters history`
   with `A finished orphan Review Artifact enters history`.
   `go test ./internal/baseline -run '^TestTheRetirementClausesCarryTheirText$' -count=1`
   exited 1 in the docs-one-job-per-directory subtest with a force/text
   mismatch. Restored that text before the next sabotage.
2. Replaced `Full text in Git at` with `Full text was in Git at` in the
   source clause's fenced provenance line.
   `go test ./internal/baseline -run '^TestAFindingInTheGuideReducedFormNeedsNoSanitize$' -count=1`
   exited 1 in `guide_form`: `IsReducedHistoryEntry = false, want true`.
   Restored the complete recorded source byte-for-byte and ran
   `make baseline-digests` again. Final focused checks passed with the
   restored source and rendered guides.

### Files rewritten by sanctioned commands

Module Version Record step:

- `internal/baseline/assets/modules/context-workflow.json`
- `internal/baseline/module-versions.json`

`make baseline-digests`:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

Managed Refresh:

- `docs/agents/docs-layout.md`
- `docs/agents/setup-context.json`

The second refresh rewrote no files. No module pin, golden, snapshot or
managed guide was hand-edited. Source Baseline assets, retention transitions,
History Relocation, history sanitize and production Go files are unchanged.
No follow-up implementation was added. The declared Task Verification and
repository Verification were not executed; settlement remains with the Daemon.


### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261008T195600Z_84187a61e36f6f25/verification/batch-004-attempt-1.log`.
The configured gate reported a stale Behavior Surface fingerprint for
`docs/user-guide/context-driven-development.md`, whose retirement paragraph
this Task changes. The authoring prototype's derived-file list did not include
this record dependency. The repair updates the generated record for that
existing Task slice; no guide text or test contract was weakened.

Focused evidence, with `GOCACHE=/tmp/roundfix-task04-gocache`:

- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -count=1`:
  reproduced the stale guide fingerprint, exit 1, before repair.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  exit 0; the repository's record generator rewrote
  `docs/references/behavior-surfaces.json`.
- Reran the focused check without the recording flag: exit 0.
- Inspected the generated diff and compared parsed records: exactly one
  fingerprint changed, for `guide docs/user-guide/context-driven-development.md`;
  every other surface and record field stayed unchanged.
- `git diff --check`: exit 0 after the repair and Result update.

The record was generated, not hand-edited. It is the only additional artifact
from this feedback repair. Task status, Task authoring, other Tasks and the
Task Graph remain unchanged. Neither the Task's declared Verification nor
`make verify-changed` was rerun by this Agent. The Daemon owns the next full
configured Verification sequence and settlement.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`
