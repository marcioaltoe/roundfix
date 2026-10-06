---
task: task_04
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: pending
type: docs
complexity: medium
---

# Task 04: The skills and the Spec workflow guides state the evidence cut

## Overview

Agents and adopters learn what an archive keeps from the skills and the
Baseline guides. Two Baseline clauses still say that `qa/` evidence stays
under the Spec folder and that archived Specs stay byte-identical, and no
skill tells the QA gate that its evidence does not outlive the archive.
This Task adds the TechSpec's "Exact texts" to two clauses, adds one section
to each of the archive-spec, qa-gate and Roundfix skills, and regenerates the
derived files with the repository's own commands.

## Requirements

1. MUST append the TechSpec's "Exact texts" sentence to the guidance of
   `clause.spec.project-constraints-05-legacy-and-ownership` and of
   `clause.spec.keep-artifacts-in-spec-folder` in
   `internal/baseline/assets/modules/spec-workflow.json`, after each
   clause's existing last sentence. Every existing sentence stays
   byte-identical, so `TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand`
   keeps passing. No clause is removed or renamed and no module version
   changes, so no Source Baseline retention disposition applies.
2. MUST add, in `.agents/skills/archive-spec/SKILL.md`, a
   `## Raw QA evidence at archive` section before `## Unarchive`. It says
   that a normal archive drops `qa/evidence/`, keeps every QA Report and
   writes `qa/evidence-manifest.md`; that an override or superseded archive
   keeps its files; that an unarchived Spec keeps its manifest and a dropped
   file is read back with `git show` at the recorded revision; and that
   `roundfix archive <slug> --drop-evidence --apply` cuts an archived Spec
   only on request, after a dry run.
3. MUST add, in `.agents/skills/qa-gate/SKILL.md`, a
   `## Raw evidence does not outlive the archive` section before
   `## Anti-patterns`. It says that the archive drops `qa/evidence/` and
   keeps the QA Report with `qa/evidence-manifest.md`, so each row records
   its observed result in the report, and material a later reader needs
   goes to the Spec's `references/` or to `docs/references/`.
4. MUST add, in `.agents/skills/roundfix/references/archive.md`, a
   `### Raw QA evidence` subsection at the end of `## Archive Command`
   stating the cut, the manifest's fields, the confirmation suffix of API
   Contract 1, and `roundfix archive <slug> --drop-evidence` with `--apply`,
   its dispositions and refusals.
5. MUST NOT change any text inside the `### QA settlement` section of any
   skill, nor any other section of the three skills.
6. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of the three changed skills and records
   them in `skills/testdata/owned-skill-versions.json`.
7. MUST run `make baseline-digests` twice (the second reports
   `"changed":false`), then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
   twice; the second refresh MUST report `File changes: 0`. A probe in a
   disposable clone at `08660fcc` showed the regenerated set: the two
   formatter goldens, the profile, `catalog.diagnostics.golden.json`,
   `catalog.digest`, `catalog.normalized.json`, the four plan
   characterization goldens, `docs/agents/docs-layout.md`,
   `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`. MUST
   NOT hand-edit a snapshot, golden, digest, generated guide or record.
8. MUST add `internal/baseline/evidence_cut_clause_test.go` with
   `TestTheGuidesStateTheQAEvidenceCut`. It asserts both added sentences in
   the embedded clauses, in the Standard TypeScript Monorepo formatter
   goldens of `docs-layout.md` and `spec-routing.md`, and in this
   repository's `docs/agents/docs-layout.md` and `docs/agents/spec-routing.md`.
9. MUST NOT change `CONTEXT.md`, `CHANGELOG.md`, any other clause or
   module, any other skill or any file under `docs/history/`.

## Subtasks

- [ ] Append the two clause sentences.
- [ ] Add the three skill sections and sync the mirrors.
- [ ] Record the raised skill versions.
- [ ] Regenerate the derived Baseline files and run the Managed Refresh.
- [ ] Add the clause test.

## Acceptance Criteria

- [ ] Both generated guides and both formatter goldens state the evidence
      cut and keep their existing sentences.
- [ ] The three skills name `qa/evidence-manifest.md`, the archive-spec and
      Roundfix ones also name the `--drop-evidence --apply` command, and
      `### QA settlement` is unchanged and identical in all three.
- [ ] Each mirror equals its canonical file and every raised version is
      recorded.
- [ ] A second `make baseline-digests` and a second Managed Refresh change
      nothing.

## Context

- instruction: `docs/adr/0243-an-archive-keeps-the-qa-report-and-drops-the-raw-qa-evidence.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `internal/baseline/evidence_cut_clause_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheGuidesStateTheQAEvidenceCut|TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestNoTwoBaselineClausesShareText)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGuidesStateTheQAEvidenceCut TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand TestCatalogCompatibility TestBaselinePlanCharacterization TestNoTwoBaselineClausesShareText; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for file in .agents/skills/archive-spec/SKILL.md .agents/skills/qa-gate/SKILL.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'qa/evidence-manifest.md' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; cmp -s "$file" "skills/${file#.agents/skills/}" || { printf 'mirror differs: %s\n' "$file" >&2; exit 1; }; done; for file in .agents/skills/archive-spec/SKILL.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'roundfix archive <slug> --drop-evidence --apply' || { printf 'missing command in %s\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- 'QA evidence cut: an archive, or an explicit request to cut an archived Spec' || { printf 'missing phrase in docs/agents/spec-routing.md\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/docs-layout.md | grep -qF -- 'keeps every QA Report beside' || { printf 'missing phrase in docs/agents/docs-layout.md\n' >&2; exit 1; }` — expected: exit 0; before this Task the clause test does not exist, no skill or guide names the manifest, and the guides lack both sentences, so the command fails; after it the clause, catalog, plan and skill-contract tests pass, each skill mirror equals its canonical file, and both guides carry their sentence.

## References

- `_prd.md` → Goals; Core Feature 8; User Story 2
- `_techspec.md` → Exact texts; Risks & Considerations; Build Order 4
- ADR-0243; ADR-0233; ADR-0189
