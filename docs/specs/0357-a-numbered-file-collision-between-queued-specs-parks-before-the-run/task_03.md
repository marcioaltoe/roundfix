---
task: task_03
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: completed
type: docs
complexity: medium
---

# Task 03: CONTEXT.md, ADR-0367, the guides and the deliver reference state the rule

## Overview

After task_02 the queue can park a Numbered File Collision, and task_04 makes
Delivery Revalidation find one, but no glossary entry, decision record, guide
or skill says so. This Task writes the two new
terms and the changed Delivery Revalidation into `CONTEXT.md`, creates
ADR-0367, documents `delivery.numbered_files` in the configuration guide and
the collision park in the `deliver` guide and both copies of the Roundfix
`deliver` reference, and extends the park action contract to the new park, so
adopters and Agents learn the rule the fiscus report asked for, reported by
fiscus on 2026-11-05
([references](references/2026-11-05-fiscus-delivery-queue-report.md)).

## Requirements

1. MUST, before editing, run
   `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem|TestGlossaryPhrasesAreDefinedTerms|TestTheSkillCoverageMapIsCurrent)$'`
   and record in the Result that it passes.
2. MUST, through `domain-modeling`, add to `CONTEXT.md`, where every
   `<ADR-0367 link>` below is a Markdown link labelled `ADR-0367` to the ADR
   file requirement 3 creates, written with its repository path as the
   neighbouring entries write theirs, right after the
   **Derived Path Declaration** entry, verbatim apart from line wrapping:

   ```markdown
   **Numbered File**:
   A file that a `delivery.numbered_files` pattern in Project Config matches and whose name starts with digits; its directory and the value of those digits are its place in a sequence, such as a database migration (<ADR-0367 link>).
   _Avoid_: Migration number, sequence file

   **Numbered File Collision**:
   A Numbered File that a Delivery Queue item adds, through its branch or the `creates:` declaration of a Task not yet completed, in the same directory and with the same number as a file on the default branch or one that an earlier item of the same queue adds while it is not merged. Delivery Revalidation parks the item as `numbered-file-collision` in the conflict Park Class, naming each pair and the next free number (<ADR-0367 link>).
   _Avoid_: Migration conflict, numbering clash
   ```

   and change the definition of **Delivery Revalidation** to, verbatim apart
   from line wrapping: "The check of a queued Spec against the refreshed
   default branch immediately before its first Run. A strict finding parks it
   as `revalidation-failed`; a Numbered File Collision parks it as
   `numbered-file-collision`
   (<ADR-0367 link>);
   a production-premise overlap records a `premise-changed` warning without
   parking it. It does not run for an Archived Item Resume." Its `_Avoid_`
   line is unchanged.
3. MUST create
   `docs/adr/0367-a-numbered-file-collision-with-an-earlier-queue-item-parks-before-the-run.md`
   with the content of `_techspec.md` → ADR-0367, apart from line wrapping.
4. MUST add to `docs/user-guide/configuration.md`, in the Delivery conflict
   recovery section right before the paragraph that starts
   "`delivery.item_binary` declares", verbatim apart from line wrapping:

   ````markdown
   `delivery.numbered_files` declares the files a repository numbers in
   sequence, such as database migrations. It defaults to `[]`, which turns the
   check off; Project Config replaces the User Config list, and an explicit
   empty list disables it.

   ```yaml
   delivery:
     numbered_files: [drizzle/*.sql, drizzle/meta/*_snapshot.json]
   ```

   Entries use the path grammar of `delivery.derived_paths[].paths`; an unsafe
   path or an invalid pattern is refused with a `delivery.numbered_files`
   config error. A matching file whose name starts with digits is a Numbered
   File; a matching `drizzle/meta/_journal.json` is not. The queue reads the
   list from the fetched default-branch commit, with User Config beneath it,
   and Delivery Revalidation parks an item whose Numbered File collides; see
   [deliver](commands/deliver.md). Declare only directories numbered in
   sequence, because names that start with a date also start with digits.
   ````

5. MUST add the paragraph below to `docs/user-guide/commands/deliver.md`
   right after the paragraph that ends "An item resumed from an item branch
   that already archives its Spec is not revalidated.", and to
   `.agents/skills/roundfix/references/deliver.md` at the same place,
   verbatim apart from line wrapping:

   "Revalidation also looks for a Numbered File Collision when the default
   branch's Project Config declares `delivery.numbered_files`. The Numbered
   Files an item adds are those its branch adds since its merge base with the
   default branch, and those that a Task of its Spec not yet completed declares
   under `creates:` while its merge base does not hold them. One that shares a
   directory and a number, read by value, with a file on the default branch,
   or with a file that an earlier item of the same queue adds while it is not
   merged, parks the item before any Run as
   `numbered-file-collision: <path> collides with <other> (<slug or default branch>), ...; next free number: <directory>/<digits>, ...`.
   The next free number is one above the highest number that the default
   branch and the earlier unmerged items hold in that directory, with the
   item's own number of digits. Later items, merged items and an Archived Item
   Resume are not compared. A Delivery Retry of an item that has not run
   repeats the check and refuses with `numbered file collision: <detail>`
   while a collision remains; a retry of an item whose Run started is not
   checked. Declare migration files under `creates:` in the plan: a number a
   Run generates without a declaration is seen only once a branch holds it."
6. MUST, in the Park Class table of `docs/user-guide/commands/deliver.md`,
   make the conflict row read
   ``| `conflict` | `pull-request-conflict`, `numbered-file-collision` |``, and
   in the table of `.agents/skills/roundfix/references/deliver.md` make it
   read ``| `conflict` | `pull-request-conflict`; merge, resolve, commit, then retry; `numbered-file-collision`; renumber the plan, or deliver the earlier item and merge the default branch, then retry |``;
   and add to both files, right after the sentence that starts
   "A `secret-found` park's next action is", the sentence, verbatim:
   ``A `numbered-file-collision` park's next action is `in <worktree>, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry <slug>`.``
7. MUST, in `internal/docscontract/command_documentation_test.go`, make
   `checkParkActionDocumentation` also require, in each of the three park
   documents, the action sentence of `delivery.BlockerNumberedFileCollision`
   exactly once (beside `BlockerArchivedItemUnproven` and
   `BlockerSecretFound`) and a conflict row (the line starting
   ``| `conflict` |``) that contains `` `numbered-file-collision` ``, failing
   with the document path otherwise; and add
   `TestTheParkActionContractRequiresTheNumberedFileCollision`, calling
   `t.Parallel()` first, which removes that sentence from the first document,
   and separately that blocker from its conflict row, and asserts
   `checkParkActionDocumentation` reports the document path each time.
8. MUST run sabotage 13 of `_techspec.md` → Testing Approach on a temporary
   copy and record that the named test fails.
9. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and, as the LAST step after the last guide edit, re-record the Behavior
   Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST run `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
   after the edits and record that it passes.
11. MUST NOT change any other glossary entry, accepted ADR, guide, skill or
    reference, `docs/references/skill-coverage.json`, `.roundfixrc.yml`, any
    `### QA settlement` section, any Go file outside the docs contract test,
    the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the park action, glossary and coverage contracts.
- [ ] Write the glossary entries through domain-modeling and ADR-0367.
- [ ] Document the key and the park in the guides and the reference.
- [ ] Extend the park action contract and its rejection test.
- [ ] Sync and record the skill; run the glossary test; re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Numbered File** and **Numbered File Collision** and the changed **Delivery Revalidation**, and ADR-0367 exists with status `accepted`.
- [ ] The configuration guide documents `delivery.numbered_files`; the `deliver` guide and both reference copies document the rule, the conflict row and the next action, and the contract test requires them.

## Context

- interface: `CONTEXT.md`
- creates: `docs/adr/0367-a-numbered-file-collision-with-an-earlier-queue-item-parks-before-the-run.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- The park action documentation contract — consumers: `./internal/docscontract`

## Verification

- `adr=docs/adr/0367-a-numbered-file-collision-with-an-earlier-queue-item-parks-before-the-run.md; test -f "$adr" && grep -q '^status: accepted$' "$adr" || { echo 'missing accepted ADR-0367' >&2; exit 1; }; for term in '**Numbered File**:' '**Numbered File Collision**:' '**Delivery Revalidation**:'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { echo "missing glossary term $term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -q -- 'a Numbered File Collision parks it as .numbered-file-collision.' || { echo 'Delivery Revalidation not changed' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'declares the files a repository numbers in sequence' || { echo 'configuration guide lacks the key' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Revalidation also looks for a Numbered File Collision' || { echo "missing collision paragraph in $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || { echo 'deliver reference mirror differs' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -tags docscontract -v ./internal/docscontract -run '^(TestTheParkActionContractRequiresTheNumberedFileCollision|TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem|TestTheParkActionContractReportsAnActionInTheClassTable|TestGlossaryPhrasesAreDefinedTerms|TestTheSkillCoverageMapIsCurrent)$' 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheParkActionContractRequiresTheNumberedFileCollision ' || { echo 'missing pass: TestTheParkActionContractRequiresTheNumberedFileCollision' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — expected: exit 0. Before this Task ADR-0367 does not exist, so the first check fails; afterwards the glossary carries the terms, ADR-0367 is accepted, the guides and both reference copies state the rule, the park action contract requires it and rejects its absence, the glossary phrases are defined, the version record and the Behavior Surface Record are current.

## References

- `_prd.md` → Goal 4; Core Feature 5; Core Feature 7; Success Metric 5; Glossary
- `_techspec.md` → ADR-0367; Coverage Map; Build Order 3

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, coding-guidelines, testing-boss, evidence-gate, systematic-debugging

- Acceptance criterion 1: added the Numbered File and Numbered File Collision glossary entries, changed Delivery Revalidation, and created accepted ADR-0367.
- Acceptance criterion 2: documented `delivery.numbered_files`, the collision revalidation rule, conflict rows and next action in the user guide and both mirrored deliver references; extended the docs-contract test and synchronized the owned skill mirror.
- Pre-edit contract check passed: `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem|TestGlossaryPhrasesAreDefinedTerms|TestTheSkillCoverageMapIsCurrent)$'`.
- Focused post-edit checks passed: the glossary test; the three park-action contract tests; `make skills-sync`; and `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
- Sabotage 13 on a temporary copy behaved as required: removing the first document's collision action made `TestTheParkActionContractRequiresTheNumberedFileCollision` fail with exit 1 and the document path.
- Final Behavior Surface Record regeneration passed with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
- `git diff --check` passed and the two deliver reference copies compare identical. Declared Task Verification was not run; it remains Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20261105T211401Z_cae788fdc7c4d6e1`
- Source commit: `b66647650bf98c732a781abae09cdc33d77c5cf6`
