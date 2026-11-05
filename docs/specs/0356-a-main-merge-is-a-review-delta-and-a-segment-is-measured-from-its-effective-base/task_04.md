---
task: task_04
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: completed
type: docs
complexity: medium
---

# Task 04: The guides, the Roundfix skill, the glossary and four ADRs state the Effective Base and lineage rules

## Overview

The `review` guide and the Roundfix `review` reference still say that a
changed merge base starts a new lineage and describe segments as plain commit
ranges, the `deliver` guide and reference do not know the `Review:` line,
`CONTEXT.md` has no **Effective Base**, and ADR-0197, ADR-0282, ADR-0286 and
ADR-0352 still prescribe the old behavior. This Task writes what tasks 01 to
03 built, as ADR-0366 records. It answers points 2 and 4 of the Fiscus report
adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md).

## Requirements

1. MUST make, in the canonical `.agents/skills/roundfix/references/review.md`
   and in `docs/user-guide/commands/review.md`, every change of
   `_techspec.md` → API Contract 6 for those files, word for word: the
   Effective Base paragraph after the paragraph ending "all segments are one
   round." (task_01 already wrote it into the guide; keep it there), the two
   lineage sentences replaced, and the six restart reasons introduced by "A
   restart names one of these reasons:" right after the paragraph that ends
   "its recorded previous data.".
2. MUST add, in `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`, the two paragraphs of API
   Contract 6 at the places it names, word for word, and change nothing else
   in either file (Spec 0355 owns the rest of the `deliver` reference).
3. MUST write into `CONTEXT.md`, through the `domain-modeling` skill, the
   **Effective Base** entry of API Contract 6 after the **Review Segment**
   entry, and the two changes it names to the **Review Segment** and
   **Reviewer Lineage** entries, each citing ADR-0366.
4. MUST append the four amendment notices of API Contract 6 as the last
   paragraph of ADR-0197, ADR-0282, ADR-0286 and ADR-0352, each opening with
   the bold `Amended by ADR-0366 (2026-11-05).` where `ADR-0366` links
   `0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`,
   and set each one's `updated_at: 2026-11-05T12:00:00Z`; every other byte of
   the four ADRs stays.
5. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT touch Go code, any other guide, any other ADR, any
   other skill or reference, any Baseline module, the Makefile, CI workflows,
   `go.mod` or `go.sum`.
6. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of the `roundfix` `SKILL.md` pair with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   recording `skills/testdata/owned-skill-versions.json`; never hand-edit a
   mirror or a version field.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of any guide, and never
   hand-edit it.

## Subtasks

- [ ] Write the Effective Base, lineage and reason text into the `review` guide and reference.
- [ ] Write the notice into the `deliver` guide and reference.
- [ ] Write the glossary entries and the four amendment notices.
- [ ] Sync the mirrors, record the versions, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the `review` reference and the `review` guide state the Effective Base paragraph, the new lineage sentences and the six reasons, and none still says that a changed merge base starts a new lineage.
- [ ] Both copies of the `deliver` reference and the `deliver` guide state the `Review:` line for status and retry.
- [ ] `CONTEXT.md` defines **Effective Base**, and ADR-0197, ADR-0282, ADR-0286 and ADR-0352 carry the amendment naming ADR-0366.
- [ ] The mirrors equal their canonical copies, and the owned-skill version record and the documentation contracts, the glossary among them, pass.

## Context

- instruction: `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- interface: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- interface: `docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md`
- interface: `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`

## Consumers

- The Roundfix skill texts and the owned-skill version record — consumers: `./skills`
- The `review` and `deliver` guide fingerprints, the glossary and the ADR contracts — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'starts from the Effective Base of the commit it ends at' 'moves only by a conflict-free merge of the default branch' 'A restart names one of these reasons:' 'which the default branch does not contain'; do printf '%s' "$t" | grep -qF -- "$p" || { echo "missing $p in $f" >&2; exit 1; }; done; ! printf '%s' "$t" | grep -qF -- 'changed provider or merge base, or a head' || { echo "old lineage rule still in $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'opens a new Reviewer Lineage at round 1, a full review of the candidate: <reason>' || { echo "missing notice in $f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no copy states the new text, so the first check fails; afterwards all do, none keeps the old lineage sentence, the mirrors equal their canonical copies, and the owned-skill version record passes.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Effective Base**: The merge base of one candidate commit with the Delivery Base' || { echo 'missing Effective Base' >&2; exit 1; }; grep -F -- '**Reviewer Lineage**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reviewer Lineage\*\*(:| —( |$))' && grep -F -- '**Review Segment**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Review Segment\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'each measured from its Effective Base and within the applicable bound' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A conflict-free merge of the default branch keeps it' || { echo 'missing glossary changes' >&2; exit 1; }; for adr in docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md; do tr -s '[:space:]' ' ' < "$adr" | grep -qF -- '**Amended by [ADR-0366](0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md) (2026-11-05).**' && grep -q '^updated_at: 2026-11-05T12:00:00Z$' "$adr" || { echo "missing amendment in $adr" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task `CONTEXT.md` has no **Effective Base** and no ADR carries the amendment, so the first check fails; afterwards they do, and the documentation contracts, the glossary and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → API Contract 6; Declared breaks; ADR Alternatives; Glossary; Build Order 4

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implementation evidence:

- Added the API Contract 6 Effective Base paragraph, updated Reviewer Lineage
  rules, and six restart reasons to both review references and the user guide.
- Added the `Review:` status and retry notices to both deliver references and
  the user guide.
- Added the Effective Base glossary entry and ADR-0366 citations to the
  Review Segment and Reviewer Lineage entries; appended the four ADR-0366
  amendment notices and updated their timestamps.
- Ran `make skills-sync`; the canonical and mirror hashes for both references
  and both `roundfix/SKILL.md` files match.
- The owned-skill recorder passed (`1 passed`) and recorded version `0.1.190`
  in both Roundfix skill files and `skills/testdata/owned-skill-versions.json`.
- The final Behavior Surface Record recorder passed (`1 passed`) and updated
  `docs/references/behavior-surfaces.json`.

Focused checks:

- `git diff --check`: passed.
- Focused Contract 6 assertions for all review and deliver copies, obsolete
  lineage text absence, glossary terms, and ADR anchors: passed.
- The four ADRs each contain the required `updated_at` value and final
  amendment paragraph; the skill `### QA settlement` section is unchanged.

Acceptance evidence:

- Review copies contain the Effective Base paragraph, conflict-free-merge
  lineage rule, and all six restart reasons; the old changed-merge-base rule
  is absent from those copies.
- Deliver copies contain the status and retry `Review:` contract text.
- `CONTEXT.md` defines Effective Base and cites ADR-0366 from the changed
  glossary entries; all four named ADRs carry the amendment.
- Mirrors and generated records were produced by their sanctioned commands;
  the Daemon must run the declared Verification commands, including the full
  documentation contract suite, before settlement.

Not run: the commands in this Task's `## Verification` section, per the
Daemon-assigned execution contract.
