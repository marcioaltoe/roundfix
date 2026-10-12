---
task: task_02
spec: 0261-a-cancelled-spec-leaves-a-record
status: pending
type: docs
complexity: high
---

# Task 02: The docs layout guide, the skills, the user guides and the glossary describe cancellation

## Overview

The CLI guide requires the skills and user guides to change with the
behavior, and the adopted Backlog Entry asks the archive-spec skill and the
docs guide to describe the cancellation. This Task writes every text of
`_techspec.md` → Skill and document texts before any code Task, so each code
Task depends on the Task that documents its surface, and the owned skills'
versions rise once. It extends one clause of the Baseline's `spec-workflow`
module, regenerates its derived files and refreshes this repository's guides.
It adds **Cancel Command** and **Source Return** and revises **Archive
Record**, **Supersession**, **Prerequisite Spec** and **Delivery Retry**
(ADR-0262). It is verifiable on its own through phrase checks, the Baseline
tests, mirror equality, the owned-skill version record and the Behavior
Surface Record.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer
   [the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)
   of 2026-10-12 on the documentation side, placing every text of
   `_techspec.md` → Skill and document texts in the file it names.
2. MUST append the clause text of `_techspec.md` → Skill and document texts to
   the `guidance` of `clause.spec.keep-artifacts-in-spec-folder` in the
   `spec-workflow` module. Every earlier byte of the clause, its `id`, its
   `mandatory` enforcement and its absent `replaces` list MUST stay. It MUST
   raise `rule.spec.docs-layout` and `guide.spec-docs-layout` by one from their
   value on this Task's starting commit, then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
3. MUST create `internal/baseline/cancel_clause_test.go` with
   `TestTheCancelClauseCarriesItsText`,
   `TestTheCancelClauseRendersInTheDocsLayoutGuide` and
   `TestAnAdopterRetainsTheCancelClause`, following the pattern of the
   retirement clause tests. The text test MUST require the new text in the
   embedded clause. The rendering test MUST check this repository's
   `docs/agents/docs-layout.md` and the Standard TypeScript golden on
   whitespace-normalized text. The retention test MUST require a ready Managed
   Refresh of a Source Baseline adopter that records the clause `retained` and
   no clause `unaccounted`. Every new top-level test MUST call `t.Parallel()`
   first or carry a `// Sequential: <reason>` comment.
4. MUST write the Roundfix skill's `## Cancel Command` section in the archive
   reference, the routing row naming `archive`, `cancel`, `history` and
   `supersede`, the deliver reference sentence "`deliver start` refuses a
   queue whose Spec requires a cancelled Spec, and `deliver retry` refuses a
   cancelled Spec", and the archive-spec skill's `## Cancelling a Spec`
   section. It MUST NOT add or change any text inside the `### QA settlement`
   section of any skill.
5. MUST create `docs/user-guide/commands/cancel.md` as `_techspec.md` → Skill
   and document texts says, add the `cancel` row to the command index in
   `docs/user-guide/commands.md` and to the command table in
   `docs/user-guide/usage.md`, and write the archive, supersede and deliver
   guide sentences. It MUST add the Skill Coverage Map entry
   `guide docs/user-guide/commands/cancel.md`, covered by the Roundfix archive
   reference with the guide as its source. It MUST NOT add the
   `command cancel` entry, which task_03 adds with the command.
6. MUST add **Cancel Command** and **Source Return** and revise **Archive
   Record**, **Supersession**, **Prerequisite Spec** and **Delivery Retry** in
   `CONTEXT.md` through `domain-modeling`, with the texts and `_Avoid_` lines
   of `_techspec.md` → Skill and document texts. Each existing `_Avoid_` line
   MUST stay.
7. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of the Roundfix and archive-spec skills,
   then
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   It MUST NOT hand-edit a version, a mirror or the Behavior Surface Record.
8. MUST NOT describe behavior beyond `_techspec.md` → API Contracts 1-4, and
   MUST NOT change any production Go file, any other clause or module, the
   Source Baseline assets or the retention transition.

## Subtasks

- [ ] Extend the clause, raise the versions, record, regenerate and refresh twice.
- [ ] Add the Baseline clause test file.
- [ ] Write the Roundfix and archive-spec skill texts and sync the mirrors.
- [ ] Write the cancel guide, the index rows and the three guide sentences.
- [ ] Add the two glossary terms and revise the four entries.
- [ ] Record the skill versions, the coverage entry and the Behavior Surface Record.

## Acceptance Criteria

- [ ] `docs/agents/docs-layout.md` tells a reader to cancel a dropped Spec
      with the runtime's cancel command instead of deleting it.
- [ ] A Source Baseline adopter's refresh retains the clause, and a second
      refresh of this repository is a no-op.
- [ ] Every phrase is present in the canonical skill, its mirror and the
      matching user guide, and the `### QA settlement` tables are
      byte-identical.
- [ ] `CONTEXT.md` carries the six entries.

## Context

- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `internal/baseline/reduced_retirement_clauses_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
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
- creates: `internal/baseline/cancel_clause_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/user-guide/commands/cancel.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/supersede.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheCancelClauseCarriesItsText|TestTheCancelClauseRendersInTheDocsLayoutGuide|TestAnAdopterRetainsTheCancelClause|TestTheRetirementClausesCarryTheirText|TestTheRetirementClausesRenderInTheDocsLayoutGuide|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheCancelClauseCarriesItsText TestTheCancelClauseRendersInTheDocsLayoutGuide TestAnAdopterRetainsTheCancelClause TestTheRetirementClausesCarryTheirText TestTheRetirementClausesRenderInTheDocsLayoutGuide TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/docs-layout.md | grep -qF -- "A Spec dropped before delivery is cancelled, never deleted by hand" || { printf 'docs/agents/docs-layout.md lacks the cancel clause\n' >&2; exit 1; }; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the three cancel clause tests do not exist, so the command fails.
- `for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do grep -q '^## Cancel Command$' "$f" || { printf 'missing Cancel Command section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- "roundfix cancel <slug> --reason <text>. writes an Archive Record with the disposition .cancelled." || { printf 'missing record phrase in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- "A cancelled Spec never meets a prerequisite, never supersedes another Spec, and a Delivery Retry refuses it" || { printf 'missing reader phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do grep -qE -- '^[|] \[archive\]\(references/archive.md\) [|] .archive., .cancel., .history., .supersede. [|]' "$f" || { printf 'missing routing row in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- "deliver start. refuses a queue whose Spec requires a cancelled Spec, and .deliver retry. refuses a cancelled Spec" || { printf 'missing deliver phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do grep -q '^## Cancelling a Spec$' "$f" || { printf 'missing cancelling heading in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- "A Spec that will not be delivered is cancelled with .roundfix cancel <slug> --reason <text>., never deleted by hand" || { printf 'missing cancelling phrase in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- "returns each adopted Backlog Entry to .docs/backlog/. with .status: open." || { printf 'missing return phrase in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qE -- "A Spec that will not be delivered is cancelled with .roundfix cancel." || { printf 'missing archive guide phrase\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/supersede.md | grep -qF -- "A cancelled Spec cannot be the superseding Spec" || { printf 'missing supersede guide phrase\n' >&2; exit 1; }; test -f docs/user-guide/commands/cancel.md || { printf 'missing cancel guide\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/cancel.md | grep -qF -- "roundfix cancel <slug> --reason <text>" || { printf 'cancel guide lacks the usage line\n' >&2; exit 1; }; grep -qE -- '^[|] .cancel. [|] \[command guide\]\(commands/cancel.md\) [|]' docs/user-guide/commands.md || { printf 'missing cancel index row\n' >&2; exit 1; }` — expected: exit 0; before this Task the archive reference has no Cancel Command section, so the command fails at its first check.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Cancel Command**: The" || { printf 'CONTEXT.md lacks **Cancel Command**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Source Return**: The Cancel Command's hand-back of a cancelled Spec's adopted sources" || { printf 'CONTEXT.md lacks **Source Return**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Archive Record**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "only its listed completed Tasks count as completed (ADR-0262)" || { printf 'CONTEXT.md lacks the revised **Archive Record**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Supersession**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A cancelled Spec never delivers another Spec's content (ADR-0262)" || { printf 'CONTEXT.md lacks the revised **Supersession**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Prerequisite Spec**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A cancelled Spec's record never meets a prerequisite (ADR-0262)" || { printf 'CONTEXT.md lacks the revised **Prerequisite Spec**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Delivery Retry**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "It refuses a Spec whose Archive Record is" || { printf 'CONTEXT.md lacks the revised **Delivery Retry**\n' >&2; exit 1; }` — expected: exit 0; before this Task neither new term exists, so the command fails.
- `grep -q '^## Cancelling a Spec$' .agents/skills/archive-spec/SKILL.md || exit 1; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestArchiveSpecContract|TestSettlementGuidanceIsOneTable)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestArchiveSpecContract TestSettlementGuidanceIsOneTable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -tags docscontract -v -run '^(TestTheSkillCoverageMapIsCurrent|TestCommandIndexNamesEveryCommandFile|TestCommandReferenceLinksResolve|TestDocsLayoutGuideNamesEveryHistoryFamily)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSkillCoverageMapIsCurrent TestCommandIndexNamesEveryCommandFile TestCommandReferenceLinksResolve TestDocsLayoutGuideNamesEveryHistoryFamily; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the archive-spec skill has no cancelling heading, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 5; Core Feature 6; User Story 1; User Story 2
- [_techspec.md](_techspec.md) — Skill and document texts; Baseline versions and derived files; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 6; Build Order 2
- ADR-0262; ADR-0247; ADR-0058; ADR-0060; ADR-0250; ADR-0256; ADR-0187; ADR-0189; ADR-0233; ADR-0244
