---
task: task_02
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
status: pending
type: docs
complexity: medium
---

# Task 02: The deliver reference, the user guide and the glossary describe the Base Update

## Overview

The Roundfix skill's `deliver` reference and the user guide page for `deliver` describe the Base Update of task_01: the `BEHIND` merge state that never merges, what runs again and what does not, the bound and the `pull-request-behind` park with its next action. The park-action contract test holds the new sentence in all three documents, the glossary gains **Base Update** and extends **Check Wait**, the Roundfix skill's version rises and the Behavior Surface Record is re-recorded.

## Requirements

1. MUST make in `.agents/skills/roundfix/references/deliver.md` and in `docs/user-guide/commands/deliver.md` each change of `_techspec.md` → API Contract 3, with its words exactly:
   - the clause "GitHub reports the merge state as neither `BLOCKED` nor `UNKNOWN`" becomes "GitHub reports the merge state as none of `BLOCKED`, `UNKNOWN` and `BEHIND`";
   - the new paragraph, right before the paragraph that begins "During `checking`", followed by a `text` block holding the line `roundfix: base update: Delivery Queue item <slug>: pull request #<n> was behind its base branch; merged the default branch at <head>`;
   - the sentence "Roundfix returns to `checking` once for that head." becomes "Roundfix returns to `checking` once for that head, where a `BEHIND` merge state makes a Base Update.", and the sentence "A second refusal for the same head parks `delivery-error`." stays;
   - `pull-request-behind` right after `checks-timeout` in the `environment` row of the Park Class table;
   - right after the sentence of the `numbered-file-collision` park, the sentence "A `pull-request-behind` park's next action is `run roundfix deliver retry <slug> once the base branch stops moving; the retry resumes at checking, which merges the default branch into the item branch again`."
2. MUST carry every guarantee of the paragraphs it edits forward: the `CONFLICTING` and `UNKNOWN` rules, the derived-merge rules, the Head Check Set, the `checks-timeout` forms and the second policy refusal parking `delivery-error` stay stated as before.
3. MUST add `delivery.BlockerPullRequestBehind` to the blockers whose action sentence `checkParkActionDocumentation` in `internal/docscontract/command_documentation_test.go` requires exactly once in each document, beside `delivery.BlockerDefectClassUndeclared`; `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` then holds the sentence of requirement 1 in the user guide, the canonical reference and its mirror.
4. MUST write, through `domain-modeling`, in `CONTEXT.md`:
   - a new entry **Base Update** with the definition of `_techspec.md` → API Contract 4, citing ADR-0435 by its link, and the line `_Avoid_: Branch update, rebase`, placed after **Check Wait**;
   - in the entry **Check Wait**, one sentence: a `BEHIND` merge state ends it in a Base Update, citing ADR-0435 by its link.
   It MUST change no other glossary entry.
5. MUST regenerate, in this order and with no hand edit of a mirror, version field or record:
   - the mirror `skills/roundfix/references/deliver.md`, with `make skills-sync`;
   - the Roundfix skill's version, with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` to the next version and records it in `skills/testdata/owned-skill-versions.json`;
   - `make skills-sync` again, then `make skills-sync-check`;
   - as the last step, after the last guide edit, the Behavior Surface Record, with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, which rewrites `docs/references/behavior-surfaces.json`.
   Spec 0407 delivers before this Spec and may have raised the same version fields and changed the same records; each command runs on the tree this Task starts from.
6. MUST NOT change any other skill, any other section of the two documents, any Go source outside the contract test, the `### QA settlement` section of any skill, or any Governed Path outside `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Edit the canonical `deliver` reference and the user guide page.
- [ ] Extend the park-action contract test.
- [ ] Write the glossary entries through `domain-modeling`.
- [ ] Synchronize the mirror, record the skill version, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The three documents state the Base Update, its bound, what runs again and the `pull-request-behind` park with the exact next action `ClassifyPark` returns.
- [ ] The park-action contract fails when the sentence is missing or reworded in any of the three documents.
- [ ] **Base Update** is a Glossary Definition, and **Check Wait** names it.
- [ ] The mirror equals the canonical reference, the Roundfix skill's version is recorded, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0435-a-delivery-pull-request-behind-its-base-branch-is-brought-up-to-date-by-a-merge-of-the-default-branch.md`
- instruction: `internal/delivery/park_class.go`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The documented park actions read through `delivery.ClassifyPark` — consumers: `./internal/docscontract`, `./skills`

## Verification

- `for file in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do text="$(tr -s '[:space:]' ' ' < "$file")"; for phrase in 'A .pull-request-behind. park.s next action is .run roundfix deliver retry <slug> once the base branch stops moving; the retry resumes at checking, which merges the default branch into the item branch again.' 'reports the merge state as none of .BLOCKED., .UNKNOWN. and .BEHIND.' 'once for that head, where a .BEHIND. merge state makes a Base Update' 'pull request #<n> is still behind its base branch after 3 Base Updates' 'roundfix: base update: Delivery Queue item <slug>: pull request #<n> was behind its base branch; merged the default branch at <head>'; do printf '%s\n' "$text" | grep -qE -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; done; grep -q 'BlockerPullRequestBehind' internal/docscontract/command_documentation_test.go || { echo 'contract test lacks the blocker' >&2; exit 1; }; grep -F -- '**Base Update**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Base Update\*\*(:| —( |$))' || { echo 'glossary lacks Base Update' >&2; exit 1; }; grep -F -- '**Check Wait**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Check Wait\*\*(:| —( |$))' || { echo 'glossary lacks Check Wait' >&2; exit 1; }; make skills-sync-check && go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the three documents holds the sentences and the contract test does not name the blocker, so the command fails; afterwards the documents, the contract test, the glossary, the mirror, the skill version record and the Behavior Surface Record agree, and the skill and documentation contract suites pass.

## References

- `_prd.md` → Goals (third and fifth bullets); Core Feature 3; Glossary
- `_techspec.md` → API Contract 2; API Contract 3; API Contract 4; Testing Approach (Task 2); Build Order 2; Glossary
- ADR-0435, ADR-0192, ADR-0311
