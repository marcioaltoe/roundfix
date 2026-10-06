---
task: task_04
spec: 0242-an-archive-that-leaves-an-archive-record
status: pending
type: docs
complexity: medium
---

# Task 04: Archive Advice and promotion, and the skills, guides, Baseline clauses and glossary describe the record

## Overview

After task_01, the cut removes everything the Spec folder held. This Task
gives the person or Agent who archives two things. `roundfix archive <slug>
--plan` lists what the cut removes and, when a judge key is set, Jev's advice
on each candidate file. `--promote <path>` copies a confirmed file to
`docs/references/` in the archive change. Jev advises and never gates
(ADR-0247). The Task also writes the guidance:

- the archive-spec, qa-gate and Roundfix skills;
- the archive and context-driven user guides;
- two Baseline clauses with their derived files;
- the glossary terms **Archive Record** and **Archive Advice**, and the
  revised **Archive Command**.

It answers the Backlog Entry "History keeps only what the Secondbrain needs"
of 2026-10-06.

## Requirements

1. MUST add the `archive-value` judgment to `internal/judge/questions.json`
   and its loader in `internal/judge/questions.go`. It is a choice question
   with the options `reusable_knowledge`, `repository_record` and
   `transient_evidence`, worded from the 2026-10-06 judgment
   (`_prd.md` → Acceptance evidence), with `content_max_chars` 6000. MUST
   add `AdviseArchive` in `internal/judge/archive_advice.go` per Invariant 14.
   It reuses the Spec judge's key variables, transports, ceiling, retry,
   redaction and log, and fails open.
2. MUST add `--plan` and `--promote <path>` (repeatable) to the Archive
   Command per Invariants 4 and 15, API Contracts 3 and 4, and Surface
   Transcripts 2, 3 and 4. `--plan` sorts files into three groups: core
   artifacts, `qa/evidence/`, and candidates. Candidates are every other
   file. It prints the advice per candidate and changes nothing under the
   repository. `--promote` passes `Promote` to `spec.Archive`.
3. MUST append the two Baseline sentences of `_techspec.md` → Exact texts,
   byte for byte, to `clause.spec.keep-artifacts-in-spec-folder` and
   `clause.context.docs-one-job-per-directory`, keeping every existing
   sentence. Then run `make baseline-digests` and the Managed Refresh so the
   derived files and this repository's `docs/agents/docs-layout.md` and
   `docs/agents/setup-context.json` follow. No clause is removed or renamed,
   so no Source Baseline retention disposition changes. MUST add
   `internal/baseline/archive_record_clause_test.go`, which asserts both
   sentences in the embedded modules, the formatter golden and the generated
   guide.
4. MUST rewrite the `### QA settlement` section identically in the
   archive-spec, qa-gate and Roundfix skills, per `_techspec.md` → Exact
   texts. A passing, qualifying partial or override archive leaves an
   Archive Record, with the override's reason when overridden. It no longer
   keeps the Spec, its QA report and evidence. The three sections stay
   byte-identical, as `TestSettlementGuidanceIsOneTable` requires, and the
   outcome names and the Settles column do not change. This follows the
   operator's decision of 2026-10-06, which lifted the briefing's rule
   against editing that section for this Spec only. MUST also add these skill
   sections, each outside `### QA settlement`:
   - `## Archive Record` in the archive-spec skill, before `## Unarchive`;
   - `## Outcome for the Archive Record` in the qa-gate skill, before
     `## Anti-patterns`;
   - `### Archive Record` at the end of `## Archive Command` in the Roundfix
     skill's archive reference.

   Each names `<slug>.md`, `source_revision`, `--plan` and `--promote`, with
   the content `_techspec.md` → Exact texts gives. Edit the canonical copies
   under `.agents/skills/`. Run `make skills-sync`, then record the raised
   versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
5. MUST update `docs/user-guide/commands/archive.md` to describe the record,
   the cut, `source_revision`, `--plan`, Archive Advice and `--promote`,
   carrying the Vocabulary Contract's patterns. MUST update the `roundfix
   archive` row of `docs/user-guide/context-driven-development.md` to name
   the Archive Record.
6. MUST make these glossary changes in `CONTEXT.md` through the
   `domain-modeling` skill:
   - add **Archive Record**: the small record an archive leaves under the
     archive root, the fields readers use, and the Spec folder kept in Git
     at its `source_revision`;
   - add **Archive Advice**: Jev's advisory classification of candidate
     files before the cut, which never gates;
   - revise **Archive Command** so it writes the record and removes the
     folder, instead of "moves the whole Spec".
7. MUST add the tests named in Verification:
   - in `internal/judge/archive_advice_test.go`, with a fake transport and a
     temporary home: only candidate files are judged; no key means skipped
     advice and no request; a reached ceiling stops; binary files are
     skipped;
   - in `internal/cli/archive_plan_test.go`: Surface Transcripts 2, 3 and 4;
     a plan leaves every file byte-identical; each promotion refusal; the
     usage errors.
8. MUST NOT change the `### QA settlement` section beyond Requirement 4, the
   Spec judge's existing judgments, or any file under `docs/history`.

## Subtasks

- [ ] Add the archive-value judgment and AdviseArchive.
- [ ] Add --plan and --promote to the Archive Command.
- [ ] Append the Baseline sentences and regenerate the derived files.
- [ ] Rewrite the three QA settlement tables identically.
- [ ] Add the skill sections, sync the mirrors and record the versions.
- [ ] Update the user guides and the glossary.

## Acceptance Criteria

- [ ] `--plan` advises without writing, and fails open without a key.
- [ ] `--promote` copies a confirmed file upstream in the archive change.
- [ ] The skills, guides, Baseline clauses and glossary describe the record,
      and every mirror and derived file matches its source.

## Context

- creates: `internal/judge/archive_advice.go`
- creates: `internal/judge/archive_advice_test.go`
- creates: `internal/cli/archive_plan_test.go`
- creates: `internal/baseline/archive_record_clause_test.go`
- interface: `internal/judge/questions.json`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/questions_test.go`
- interface: `internal/cli/archive.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
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
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- instruction: `docs/agents/secondbrain.md`
- instruction: `skills/settlement_guidance_repocontract_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveAdviceJudgesOnlyCandidateFiles|TestArchiveAdviceFailsOpenWithoutAKey|TestArchiveAdviceStopsAtTheCeiling|TestArchivePlanListsWhatTheCutRemoves|TestArchivePlanWithoutAKeyNamesTheSkip|TestArchivePlanWritesNothing|TestArchivePromoteCopiesUpstream|TestArchivePromoteRefusals|TestArchivePlanUsageErrors|TestArchiveRecordClausesAreAppended)$' ./internal/judge ./internal/cli ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveAdviceJudgesOnlyCandidateFiles TestArchiveAdviceFailsOpenWithoutAKey TestArchiveAdviceStopsAtTheCeiling TestArchivePlanListsWhatTheCutRemoves TestArchivePlanWithoutAKeyNamesTheSkip TestArchivePlanWritesNothing TestArchivePromoteCopiesUpstream TestArchivePromoteRefusals TestArchivePlanUsageErrors TestArchiveRecordClausesAreAppended; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the ten tests do not exist, so the command fails.
- `for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '## Archive Record ' || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '## Outcome for the Archive Record ' || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '### Archive Record ' || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/archive-spec/SKILL.md .agents/skills/qa-gate/SKILL.md .agents/skills/roundfix/SKILL.md skills/archive-spec/SKILL.md skills/qa-gate/SKILL.md skills/roundfix/SKILL.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$n" | grep -qF -- 'The Archive Record, which names the QA Report and verdict;' || { printf 'missing settlement wording in %s\n' "$f" >&2; exit 1; }; printf '%s' "$n" | grep -qF -- 'determines what the archive leaves:' || { printf 'missing settlement intro in %s\n' "$f" >&2; exit 1; }; if printf '%s' "$n" | grep -qF -- 'The Spec and its QA report and evidence.'; then printf 'old settlement wording in %s\n' "$f" >&2; exit 1; fi; if printf '%s' "$n" | grep -qF -- 'QA files move byte-identically'; then printf 'old override wording in %s\n' "$f" >&2; exit 1; fi; done; cmp .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md && cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md && go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -tags repocontract ./skills -run '^TestSettlementGuidanceIsOneTable$'` — expected: exit 0; before this Task none of the three skills has its section or the new settlement wording, and all six files still say "The Spec and its QA report and evidence.", so the command fails.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Archive Record**:' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Archive Advice**:' && ! tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'moves the whole Spec' && tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'roundfix archive <slug> --plan' && tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- '--promote <path>' && tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Archive Record'` — expected: exit 0; before this Task `CONTEXT.md` lacks both terms and still says the command "moves the whole Spec", so the command fails.

## References

- `_prd.md` → Goal 3; Goal 5; User Story 2; Core Feature 3; Core Feature 6; Success Metric 3
- `_techspec.md` → Invariants 4, 14 and 15; API Contract 3; API Contract 4; Surface Transcripts 2, 3 and 4; Exact texts; Vocabulary Contract; Build Order 4
- ADR-0247; ADR-0189; ADR-0233
