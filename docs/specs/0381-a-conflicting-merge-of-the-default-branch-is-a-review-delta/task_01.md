---
task: task_01
spec: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
status: pending
type: docs
complexity: low
---

# Task 01: The guides, references, glossary and ADR notices describe the Conflicted Carry

## Overview

task_02, which needs this Task, makes a pre-PR review measure across a
conflicting merge of the default branch from its Conflicted Carry, keeps the
Reviewer Lineage across that merge, retires the restart reason for it, and
adds two prompt lines and one prompt sentence, with the texts `_techspec.md`
→ API Contracts 2–4 fix. This Task writes them into the `review` and
`deliver` guides and the Roundfix skill's references first, adds the term to
the glossary, and writes the "Amended by" notices of ADR-0391 into ADR-0366
and ADR-0197, so an agent reading the skill acts as the command does once
task_02 lands. It answers the Conexus report of 2026-11-11
([`references/2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md`](references/2026-11-11-conexus-review-false-compile-findings-after-a-main-merge.md),
reported by conexus).

## Requirements

1. MUST, in `docs/user-guide/commands/review.md` and the canonical
   `.agents/skills/roundfix/references/review.md`, each (the texts are the
   same in both files):
   - replace the sentence end "it does not ask the reviewer to discover the
     candidate." in the opening paragraph with that sentence followed by:
     "Every review prompt lets the reviewer open repository files and, right
     after that line, adds: `A diff hunk shows only the lines near a change.
     Before a finding states that an import, declaration, symbol or closing
     part of a block is missing or undefined, read the whole file in the
     repository, which holds the head commit, and raise the finding only if
     the file confirms it.`" (the quoted prompt sentence in one inline code
     span);
   - change "moves only by a conflict-free merge of the default branch." to
     "moves only by a merge of the default branch, with or without a
     conflict.";
   - in the paragraph that starts "A merge of the default branch keeps the
     lineage when", delete ", and the reviewed head merges onto the new merge
     base without conflict" so the sentence ends "brings in only commits the
     new merge base contains.", and insert right before "A rebase, a changed
     provider" this text: "When the reviewed head conflicts with the new
     merge base, the delta starts from its **Conflicted Carry**, the tree
     `git merge-tree --write-tree` writes, read with
     `merge.conflictStyle=merge`, whose conflicted files hold Git's conflict
     markers, so the diff shows how the candidate resolved the conflict; the
     prompt then adds `Carrying <reviewed head> onto the default branch
     conflicted; where tree <deltaBase> holds Git's conflict markers, the
     diff shows how the candidate resolved the conflict.` and, when Git lists
     conflicted paths, `Conflicted paths: <path>, <path>.`" (each prompt line
     in one inline code span);
   - delete the bullet `the merge base moved from <prior base> to <base> and
     the reviewed head <prior head> conflicts with it` from the list after "A
     restart names one of these reasons:", leaving seven bullets;
   - replace "A carry that conflicts measures the segment from its base." with
     "A carry that conflicts measures the segment from its Conflicted Carry,
     and the prompt adds, after that header line, `Carrying <segment base>
     onto the default branch conflicted; where tree <measuredFrom> holds
     Git's conflict markers, the diff shows how the candidate resolved the
     conflict.` and, when Git lists conflicted paths, `Conflicted paths:
     <path>, <path>.`".
2. MUST, in `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`, each, change "The reason
   is one of the eight the `review` guide lists; a conflict-free merge of the
   default branch keeps the lineage and prints nothing." to "The reason is
   one of the seven the `review` guide lists; a merge of the default branch
   keeps the lineage, with or without a conflict, and prints nothing.".
3. MUST, in `CONTEXT.md`, through the domain-modeling skill:
   - add, right after the `_Avoid_:` line of **Effective Base**, a blank line
     and this entry:
     "**Conflicted Carry**:" then the line "The tree `git merge-tree
     --write-tree` writes, with Git's conflict markers in the conflicted
     files, when a candidate commit's own change does not carry onto a later
     Effective Base cleanly. A round-2 delta, a Correction Review delta or a
     Review Segment measured from it holds the candidate's own change and its
     conflict resolution and none of the default branch's work
     (" followed by the link `[ADR-0391](docs/adr/0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md)` and ").",
     then the line "_Avoid_: Conflict restart, remerge diff";
   - change, in **Effective Base**, "so neither shows the default branch's
     work as the candidate's (" to "so neither shows the default branch's
     work as the candidate's, starting from the Conflicted Carry when that
     carry conflicts (" and add, after its ADR-0366 link and inside the same
     parentheses, ", " followed by the link
     `[ADR-0391](docs/adr/0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md)`;
   - change, in **Reviewer Lineage**, "A conflict-free merge of the default
     branch keeps it, and the delta then starts from the round-1 head carried
     onto the new merge base." to "A merge of the default branch keeps it, and
     the delta then starts from the round-1 head carried onto the new merge
     base, or from its Conflicted Carry when that merge conflicts." and add
     the same ADR-0391 link after its ADR-0366 link.
4. MUST write the amendment notices of ADR-0391, setting `updated_at` to
   `2026-11-11T18:00:00Z` in each file's front matter:
   - in `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`,
     a new paragraph right after the Decision paragraph that ends "starts
     round 1 as before.", starting with the bold notice
     `**Amended by [ADR-0391](0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md) (2026-11-11).**`
     and continuing: "A carry that conflicts no longer measures from the
     earlier commit or restarts the lineage: the diff starts from its
     Conflicted Carry, the tree `git merge-tree --write-tree` writes with
     Git's conflict markers, and a merge of the default branch keeps the
     lineage with or without a conflict. The other lineage conditions, the
     Effective Base and the `Review:` line still hold."; and, at the end of
     the Considered Option "Show the merge's `--remerge-diff`.", the sentence
     "ADR-0391 takes this option's idea, a diff from Git's re-merge with
     conflict markers, over the whole range from the reviewed head.";
   - in `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`,
     a new paragraph right after the "Amended by ADR-0366" paragraph,
     starting with the same bold notice
     `**Amended by [ADR-0391](0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md) (2026-11-11).**`
     and continuing: "A merge of the default branch that conflicts keeps the
     Reviewer Lineage too; its delta starts from the reviewed head's
     Conflicted Carry. A rebase or a merge that brings in another branch
     still starts round 1, and the two-round ceiling is unchanged."
5. MUST NOT change any other text of those files, any other guide, skill,
   reference or ADR, code or test, or text inside the `### QA settlement`
   section of any skill.
6. MUST run `make skills-sync` to write the mirrors
   `skills/roundfix/references/review.md` and
   `skills/roundfix/references/deliver.md`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise both version fields of the Roundfix `SKILL.md` pair to the next
   version and record them in `skills/testdata/owned-skill-versions.json`;
   never hand-edit a mirror, a version field or the record.
7. MUST, as the last edit of the Task, after every guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Write the absence-claim sentence, the Conflicted Carry text and the seven reasons into both `review` texts, and the count into both `deliver` texts.
- [ ] Add **Conflicted Carry** to the glossary and reword **Effective Base** and **Reviewer Lineage**.
- [ ] Write the amendment notices in ADR-0366 and ADR-0197.
- [ ] Sync the mirrors, record the skill version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both `review` texts describe the Conflicted Carry for the delta and for segments, the absence-claim sentence and seven restart reasons, and both `deliver` texts name seven reasons.
- [ ] The glossary defines **Conflicted Carry** and its two neighbours use it; ADR-0366 and ADR-0197 carry the ADR-0391 notice.
- [ ] The mirrors match the canonical references, the owned skill version record is current, and the Behavior Surface Record is current.

## Context

- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- interface: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `docs/adr/0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- The Roundfix skill's `review` and `deliver` references and their version record, and the glossary — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'the delta starts from its **Conflicted Carry**' 'onto the default branch conflicted; where tree <deltaBase> holds Git' 'onto the default branch conflicted; where tree <measuredFrom> holds Git' 'Conflicted paths: <path>, <path>.' 'read the whole file in the repository, which holds the head commit, and raise the finding only if the file confirms it.' 'moves only by a merge of the default branch, with or without a conflict.'; do printf '%s\n' "$t" | grep -qF -- "$p" || { echo "missing in $f: $p" >&2; exit 1; }; done; if printf '%s\n' "$t" | grep -qF -- 'conflicts with it'; then echo "retired reason still in $f" >&2; exit 1; fi; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'The reason is one of the seven the' || { echo "missing in $f: seven reasons" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md` — expected: exit 0. Before this Task no guide or reference names the Conflicted Carry and all three `review` texts list the retired reason; afterwards all six files carry the text and each mirror matches its canonical copy.
- `grep -F -- '**Conflicted Carry**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Conflicted Carry\*\*(:| —( |$))' && grep -F -- '**Effective Base**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Effective Base\*\*(:| —( |$))' && grep -F -- '**Reviewer Lineage**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reviewer Lineage\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'starting from the Conflicted Carry when that carry conflicts' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'or from its Conflicted Carry when that merge conflicts' && tr -s '[:space:]' ' ' < docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md | grep -qF -- '**Amended by [ADR-0391](0391-' && tr -s '[:space:]' ' ' < docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md | grep -qF -- '**Amended by [ADR-0391](0391-' && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestGlossaryAvoidLinesArePlainLists)$'` — expected: exit 0. Before this Task the glossary has no **Conflicted Carry** entry, so the first check fails; afterwards the three entries are definitions, both ADRs carry the notice, the owned skill version record and skill contracts pass, the Behavior Surface Record is current, and the glossary contracts pass.

## References

- `_prd.md` → Core Feature 5; Glossary
- `_techspec.md` → API Contract 2; API Contract 3; API Contract 4; Build Order 1
- ADR-0391; ADR-0366; ADR-0197
