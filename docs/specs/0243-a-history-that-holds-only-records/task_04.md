---
task: task_04
spec: 0243-a-history-that-holds-only-records
status: pending
type: docs
complexity: medium
---

# Task 04: The skill, a Baseline clause, the glossary and the export contract describe the sanitize

## Overview

The History Sanitize Command needs the same guidance the archive has. This
Task writes it:

- the Roundfix skill's archive reference and index;
- one Baseline sentence on where retired documentation goes, with its
  derived files;
- the glossary terms **History Sanitize Command**, **Legacy Archive
  Folder**, **Sanitize Batch**, **Reduced History Entry** and **History Full
  Tag**, and the revised **Archive Record** and **History Root**;
- the two `CHANGELOG.md` lines that cite folders whose measurements moved to
  `docs/references/`;
- a documentation contract that ties `.secondbrain-export` to the history's
  form.

The export contract holds both directions. While any Legacy Archive Folder
remains, the export keeps its history exclusions, so the mirror does not
grow back mid-migration. Once none remains, the export must mirror
`docs/` whole, which the operator's last batch does. This Task answers the
Backlog Entry "History keeps only what the Secondbrain needs" of 2026-10-06,
Expected 5.

## Requirements

1. MUST add `## History Sanitize Command` to
   `.agents/skills/roundfix/references/archive.md`, after the
   `## Archive Command` section, with the content of `_techspec.md` → Exact
   texts. MUST add `history` to the commands of the `archive` row of the
   index in `.agents/skills/roundfix/SKILL.md`. Run `make skills-sync`, then
   record the raised version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   MUST NOT change the `### QA settlement` section.
2. MUST append the Baseline sentence of `_techspec.md` → Exact texts, byte
   for byte, to `clause.context.docs-one-job-per-directory` in
   `internal/baseline/assets/modules/context-workflow.json`, keeping every
   existing sentence. Then run `make baseline-digests` and the Managed
   Refresh (`go run -buildvcs=false ./cmd/roundfix baseline update --repo .
   --no-skills --yes --format text`) so the derived files and this
   repository's `docs/agents/docs-layout.md` and
   `docs/agents/setup-context.json` follow. No clause is removed or renamed,
   so no Source Baseline retention disposition changes. MUST add
   `internal/baseline/history_sanitize_clause_test.go`, which asserts the
   sentence in the embedded module, the formatter golden and the generated
   guide.
3. MUST make these glossary changes in `CONTEXT.md` through the
   `domain-modeling` skill, citing ADR-0248 and never a Spec:
   - add **History Sanitize Command**: `roundfix history sanitize`, a dry
     run unless `--apply --batch <n>`, which converts the existing history
     after the History Full Tag;
   - add **Legacy Archive Folder**: a Spec folder an archive left under the
     archive root before archives wrote Archive Records;
   - add **Sanitize Batch**: the next units one `--apply` converts, delivered
     as one Pull Request with the repository gates green, revertible;
   - add **Reduced History Entry**: a retired Finding or Backlog Entry cut to
     its front matter, title, first paragraph and the revision holding its
     full text;
   - add **History Full Tag**: the annotated `history-full` tag on the last
     commit before the first batch, which keeps every removed byte
     reachable;
   - revise **Archive Record** to say the History Sanitize Command also
     writes one for a Legacy Archive Folder, and that a folder without a QA
     Report gets the disposition `no-qa`;
   - revise **History Root** to say it holds Archive Records, Reduced
     History Entries and retired ADRs, with the full bytes in Git.
4. MUST repoint `CHANGELOG.md` lines 188 and 198: the 0217 measurement path
   becomes `docs/references/grok-through-cursor-measurement.md`, and the
   0218 one becomes `docs/references/jev-router-measurement-2026-10-02.md`.
   No other `CHANGELOG.md` line changes.
5. MUST add `internal/docscontract/secondbrain_export_test.go` with the
   tests named in Verification. The rule: when any directory directly under
   `docs/history/specs` holds `_prd.md`, `.secondbrain-export` holds the
   line `!docs/history/specs/*/qa/`; when none does, it holds no line that
   starts with `!docs/history`. A repository without `.secondbrain-export`
   is skipped. One test applies the rule to this repository, where it must
   pass unchanged. The other applies it to temporary fixtures and proves it
   fails for an export that keeps the exclusions without folders, and for
   one that drops them while a folder remains.
6. MUST NOT change `.secondbrain-export`, `.github/workflows/`, the
   Makefile, or any file under `docs/history`.

## Subtasks

- [ ] Add the skill section and index entry, sync the mirrors and record the version.
- [ ] Append the Baseline sentence and regenerate the derived files.
- [ ] Write the glossary terms and revisions.
- [ ] Repoint the two CHANGELOG paths.
- [ ] Add the export contract tests.

## Acceptance Criteria

- [ ] The skill, the Baseline clause and the glossary describe the command,
      and every mirror and derived file matches its source.
- [ ] The export contract passes on this repository and fails each
      mismatched fixture.

## Context

- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- creates: `internal/baseline/history_sanitize_clause_test.go`
- interface: `CONTEXT.md`
- interface: `CHANGELOG.md`
- creates: `internal/docscontract/secondbrain_export_test.go`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- instruction: `internal/baseline/archive_record_clause_test.go`
- instruction: `.secondbrain-export`

## Verification

- `out="$(go test -count=1 -v -tags docscontract -run '^(TestSecondbrainExportFollowsTheHistoryForm|TestSecondbrainExportRuleRefusesBothMismatches)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; base="$(go test -count=1 -v -run '^TestHistorySanitizeClauseIsAppended$' ./internal/baseline 2>&1)" || { printf '%s\n' "$base"; exit 1; }; for name in TestSecondbrainExportFollowsTheHistoryForm TestSecondbrainExportRuleRefusesBothMismatches; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; printf '%s\n' "$base" | grep -q -- '--- PASS: TestHistorySanitizeClauseIsAppended (' || { printf 'missing pass: TestHistorySanitizeClauseIsAppended\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/docs-layout.md | grep -qF -- 'When the runtime offers a history sanitize, such as' || { printf 'generated guide lacks the sentence\n' >&2; exit 1; }` — expected: exit 0; before this Task none of the three tests exists and the guide lacks the sentence, so the command fails.
- `for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '## History Sanitize Command ' || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'history-full' || { printf 'missing tag in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do grep -E -q -- '^[|] .archive.\(references/archive[.]md\) [|].*history' "$f" || { printf 'index lacks history in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md && go test -count=1 ./skills -run '^(TestEveryOwnedSkillVersionIsRecorded|TestRoundfixSkillIndexNamesEveryReference|TestRoundfixSkillEntryFileStaysWithinItsBudget)$'` — expected: exit 0; before this Task neither copy has the section and the index row does not name history, so the command fails.
- `for term in '**History Sanitize Command**:' '**Legacy Archive Folder**:' '**Sanitize Batch**:' '**Reduced History Entry**:' '**History Full Tag**:'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { printf 'CONTEXT.md lacks %s\n' "$term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | awk -v t='**Archive Record**:' '{ i = index($0, t); if (!i) exit 1; s = substr($0, i); j = index(s, "_Avoid_"); if (!j || !index(substr(s, 1, j), "History Sanitize Command") || !index(substr(s, 1, j), "no-qa")) exit 1 }' || { printf 'the Archive Record entry is not revised\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | awk -v t='**History Root**:' '{ i = index($0, t); if (!i) exit 1; s = substr($0, i); j = index(s, "_Avoid_"); if (!j || !index(substr(s, 1, j), "Reduced History Entr")) exit 1 }' || { printf 'the History Root entry is not revised\n' >&2; exit 1; }` — expected: exit 0; before this Task none of the five terms exists and neither entry names the sanitize, so the command fails.
- `! grep -qF -- '0217-a-cursor-runtime-to-measure-grok-on/measurement' CHANGELOG.md && ! grep -qF -- '0218-the-jev-router-on-docs-and-chore-tasks/measurement' CHANGELOG.md && grep -qF -- 'docs/references/grok-through-cursor-measurement.md' CHANGELOG.md && grep -qF -- 'docs/references/jev-router-measurement-2026-10-02.md' CHANGELOG.md` — expected: exit 0; before this Task `CHANGELOG.md` still names both folder paths, so the command fails.

## References

- `_prd.md` → Goal 5; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → Measured inventory; Exact texts; Glossary; Operator batch procedure; Build Order 4
- ADR-0248; ADR-0189; ADR-0233
