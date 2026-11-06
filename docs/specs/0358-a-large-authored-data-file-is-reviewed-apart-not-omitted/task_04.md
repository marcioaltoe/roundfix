---
task: task_04
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: completed
type: docs
complexity: medium
---

# Task 04: The review guide, the Roundfix skill, the glossary and ADR-0282 describe Data Review Segments

## Overview

After tasks 01 to 03 a declared data file is reviewed apart, but the
`review` guide and the Roundfix `review` reference describe only omissions
and Review Segments, `CONTEXT.md` defines neither **Review Data Path** nor
**Data Review Segment**, and ADR-0282 still rejects splitting a diff by file
without exception. This Task writes the rule, as ADR-0368 records. It answers
point 5 of the Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md).

## Requirements

1. MUST add, in the canonical `.agents/skills/roundfix/references/review.md`
   and in `docs/user-guide/commands/review.md`, after the paragraph ending
   "round two's delta and the Correction Review delta.", the two paragraphs of
   `_techspec.md` → API Contract 6, word for word (line wrapping free), and
   change nothing else in either file.
2. MUST write into `CONTEXT.md`, through the `domain-modeling` skill, the
   **Review Data Path** and **Data Review Segment** entries of API Contract 6
   after the **Review Omit Path** entry, each citing ADR-0368 as a Markdown
   link to `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`.
3. MUST write into `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
   the amendment of API Contract 6: the last paragraph opening with the bold
   `Amended by ADR-0368 (2026-11-05).` where `ADR-0368` links
   `0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`, the sentence "ADR-0368 takes it for declared data files only." at
   the end of the Rejected bullet "**Splitting the diff by file or by
   bytes.**", and `updated_at: 2026-11-05T12:00:00Z`; every other byte stays.
4. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT touch Go code, the configuration guide, any other
   guide, ADR or skill, any Baseline module, the Makefile, CI workflows,
   `go.mod` or `go.sum`.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of the `roundfix` `SKILL.md` pair with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   recording `skills/testdata/owned-skill-versions.json`; never hand-edit a
   mirror or a version field.
6. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide edit, and never hand-edit it.

## Subtasks

- [ ] Write the two paragraphs into the `review` guide and reference.
- [ ] Write the glossary entries and ADR-0282's amendment.
- [ ] Sync the mirrors, record the versions, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the `review` reference and the `review` guide describe Review Data Paths and Data Review Segments.
- [ ] `CONTEXT.md` defines **Review Data Path** and **Data Review Segment**, and ADR-0282 carries the amendment naming ADR-0368.
- [ ] The mirrors equal their canonical copies, and the owned-skill version record and the documentation contracts, the glossary among them, pass.

## Context

- instruction: `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`

## Consumers

- The Roundfix skill texts and the owned-skill version record — consumers: `./skills`
- The `review` guide fingerprint, the glossary and the ADR contracts — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'is a Review Data Path with reason' 'as at most four Data Review Segments' 'needs <d> data review segments; a review allows at most 4'; do printf '%s' "$t" | grep -qF -- "$p" || { echo "missing $p in $f" >&2; exit 1; }; done; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no copy describes Review Data Paths, so the first check fails; afterwards all three do, the mirrors equal their canonical copies, and the owned-skill version record passes.
- `grep -F -- '**Review Data Path**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Review Data Path\*\*(:| —( |$))' && grep -F -- '**Data Review Segment**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Data Review Segment\*\*(:| —( |$))' || { echo 'missing glossary entries' >&2; exit 1; }; adr=docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md; tr -s '[:space:]' ' ' < "$adr" | grep -qF -- '**Amended by [ADR-0368](0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md) (2026-11-05).**' && tr -s '[:space:]' ' ' < "$adr" | grep -qF -- 'ADR-0368 takes it for declared data files only.' && grep -q '^updated_at: 2026-11-05T12:00:00Z$' "$adr" || { echo 'missing amendment in ADR-0282' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `CONTEXT.md` defines neither term and ADR-0282 has no amendment, so the first check fails; afterwards both hold, and the documentation contracts, the glossary and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → API Contract 6; ADR Alternatives; Glossary; Build Order 4

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

- Added the API Contract 6 Review Data Path and Data Review Segment paragraphs
  to the canonical Roundfix review reference and the review guide. `make
  skills-sync` exited 0 and the canonical and mirror review references have
  identical SHA-256 hashes.
- Added both API Contract 6 glossary entries to `CONTEXT.md`, with the ADR-0368
  Markdown link, and amended ADR-0282's rejected alternative and final ADR
  amendment paragraph. The focused inspection found both terms and both ADR
  sentences at their required locations.
- Ran `git diff --check`: passed.
- Ran `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: 1 test passed; both `roundfix` SKILL.md copies are
  version `0.1.191` and `skills/testdata/owned-skill-versions.json` records it.
- The Task's declared Verification commands and final Behavior Surface record
  command were not run; they remain for Daemon Verification. Task status was
  left Daemon-owned and unchanged.
