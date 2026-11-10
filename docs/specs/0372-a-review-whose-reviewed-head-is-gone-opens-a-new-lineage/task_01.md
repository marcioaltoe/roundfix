---
task: task_01
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
status: pending
type: docs
complexity: low
---

# Task 01: The review and deliver guides and references name a commit that is gone

## Overview

task_02, which needs this Task, makes a recorded commit that a complete
repository proves absent restart the Reviewer Lineage, and makes `review
dispose` and the Delivery Retry's correction proof refuse it by name, with
the texts `_techspec.md` → API Contracts 1–3 fix. This Task writes those
reasons, refusals and the proof's bounds into the `review` and `deliver`
guides and the Roundfix skill's references first, so an agent reading the
skill acts as the command does once task_02 lands. It answers the Conexus report of 2026-11-10
([`references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md`](references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md),
reported by conexus), whose reader found no documented way out.

## Requirements

1. MUST, in `docs/user-guide/commands/review.md` and the canonical
   `.agents/skills/roundfix/references/review.md`, each:
   - add, in the list after "A restart names one of these reasons:", right
     after the bullet `the head does not descend from the reviewed head <prior head>`,
     the two bullets `the reviewed head <prior head> is no longer in the repository`
     and `the recorded merge base <prior base> is no longer in the repository`,
     each in an inline code span like the others;
   - add, right after that list's last bullet (the one ending
     `conflicts with it`), this paragraph, with `git cat-file -e`, the exit
     status `1` and the error text in inline code spans: "A recorded commit is no longer in the
     repository only when the repository is neither shallow nor partial and
     git cat-file -e answers for its full object ID with exit 1 and no output.
     A complete repository holds every ancestor of every commit it holds, so
     a reviewed head it no longer holds is a head that does not descend. A
     shallow or partial clone cannot prove that, and Roundfix never asks a
     partial clone, because the question would fetch from its promisor
     remote; the review blocks with check review lineage ancestry: a shallow
     or partial clone cannot prove that the recorded commit <commit> is gone;
     fetch the full history and try again. Every other Git failure still
     blocks the review, so no failure resets the two-round ceiling.";
   - add, in the `review dispose` paragraph, right after the sentence ending
     "is reachable from the current `HEAD`.", the sentence "When the reviewed
     head is no longer in the repository, both forms refuse with
     reviewed head "<commit>" is no longer in the repository; run roundfix review to open a new Reviewer Lineage."
     with the refusal in an inline code span.
2. MUST, in `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`, each:
   - change "The reason is one of the six the `review` guide lists" to
     "The reason is one of the eight the `review` guide lists";
   - add, right after the sentence "Outside that exception, any other path
     still requires a corrective Spec, and the refusal names that path.", the
     sentence "A parked candidate that is no longer in the repository refuses
     with parked candidate <commit> is no longer in the repository, and a
     fixing commit that is no longer in the repository is not a valid
     disposition." with the refusal in an inline code span.
3. MUST NOT change any other text of those four files, any other guide,
   skill or reference, `CONTEXT.md`, code or test, or text inside the
   `### QA settlement` section of any skill.
4. MUST run `make skills-sync` to write the mirrors
   `skills/roundfix/references/review.md` and
   `skills/roundfix/references/deliver.md`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise both version fields of the Roundfix `SKILL.md` pair and record
   them in `skills/testdata/owned-skill-versions.json`; never hand-edit a
   mirror, a version field or the record.
5. MUST, as the last edit of the Task, after every guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the two restart reasons, the proof's bounds and the dispose refusal to both `review` texts.
- [ ] Add the count and the correction refusal to both `deliver` texts.
- [ ] Sync the mirrors, record the skill version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both `review` texts list eight restart reasons, the proof's bounds and the dispose refusal, and both `deliver` texts name eight reasons and the correction refusal.
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

## Consumers

- The Roundfix skill's `review` and `deliver` references and their version record — consumers: `./skills`

## Verification

- `for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'the reviewed head <prior head> is no longer in the repository' 'the recorded merge base <prior base> is no longer in the repository' 'a shallow or partial clone cannot prove that the recorded commit <commit> is gone; fetch the full history and try again' 'no longer in the repository; run roundfix review to open a new Reviewer Lineage' 'Roundfix never asks a partial clone'; do printf '%s\n' "$t" | grep -qF -- "$p" || { echo "missing in $f: $p" >&2; exit 1; }; done; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'The reason is one of the eight the' 'parked candidate <commit> is no longer in the repository'; do printf '%s\n' "$t" | grep -qF -- "$p" || { echo "missing in $f: $p" >&2; exit 1; }; done; done; cmp -s .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md` — expected: exit 0. Before this Task no guide or reference names a commit that is no longer in the repository; afterwards all six files carry the text and each mirror matches its canonical copy.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/review.md | grep -qF -- 'is no longer in the repository' || exit 1; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms)$'` — expected: exit 0. Before this Task the mirror lacks the text, so the first check fails; afterwards the owned skill version record and the skill contracts pass, the Behavior Surface Record is current, and every capitalized phrase the new text uses is a defined term.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 4
- `_techspec.md` → System Architecture; Build Order 1
- ADR-0197; ADR-0272; ADR-0279
