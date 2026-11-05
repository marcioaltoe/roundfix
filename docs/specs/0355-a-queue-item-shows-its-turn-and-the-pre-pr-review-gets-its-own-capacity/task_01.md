---
task: task_01
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: completed
type: docs
complexity: medium
---

# Task 01: ADR-0365, CONTEXT.md, the Roundfix deliver reference and the guides state the behavior

## Overview

task_02 to task_04 give the Queue Owner review slots up to the Review
Capacity, keep the start order and print each item's Queue Turn. This Task
runs first and records that behavior in ADR-0365, adds **Review Capacity**, **Queue Turn** and
**Turn Record** to `CONTEXT.md` and revises **Queue Owner**, and writes the
stage meaning, the `Turn:` lines, the `Limits:` line and the key into the
Roundfix skill's `deliver` reference, the `deliver` guide and the
configuration guide, answering point 1 of the report reported by fiscus on
2026-11-05 ([references](references/2026-11-05-fiscus-delivery-queue-takes-days.md)).
The code Tasks that follow implement exactly the behavior these documents
state.

## Requirements

1. MUST create `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md` with exactly the content of `_techspec.md` → ADR-0365,
   apart from line wrapping.
2. MUST replace the `**Queue Owner**` entry of `CONTEXT.md` and add the three
   entries after it, so that the text from `**Queue Owner**:` up to the blank
   line before `**Delivery Plan**:` is exactly, line for line, the block
   below, where `<ADR-0365 link>` stands for
   `[ADR-0365](docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md)`
   and `<ADR-0279 link>` for the ADR-0279 link the current `**Queue Owner**`
   entry already carries, unchanged:

     ````markdown
     **Queue Owner**:
     The `roundfix deliver` process that holds one Delivery Queue and advances its items in order, one item at a time through every stage except the pre-PR review, which it runs for up to its Review Capacity of items at once. It waits on each Prerequisite Spec and starts a queued item only when no earlier item is between its Implement Run and its merge; a Delivery Retry hands a parked item to the live Queue Owner or starts a new one, and a live owner finds a retried item within 30 seconds (<ADR-0365 link>). When an item can still advance, `roundfix deliver status` reports a recorded Queue Owner proven not running by its PID and start-time identity; status never reclaims it, so the item can still advance after `roundfix deliver resume` (<ADR-0279 link>).
     _Avoid_: Scheduler, daemon, deliver loop

     **Review Capacity**:
     The maximum number of pre-PR reviews one Queue Owner runs at the same time, beside the one item it advances through every other stage. It is set by `delivery.review_concurrency` in User Config and is `2` by default; each review reads only its own item worktree in a read-only session, so it never waits for the Implement Run of another item (<ADR-0365 link>).
     _Avoid_: Verification Capacity, Task Capacity, review pool

     **Queue Turn**:
     Where a Delivery Queue item that is neither merged nor parked stands with its Queue Owner: advancing in the owner's turn, reviewing in one of the Review Capacity's slots, or waiting for the owner's turn, a review slot or its prerequisites, with its place in line. `roundfix deliver status` prints it on one `Turn:` line per item from the Turn Record (<ADR-0365 link>).
     _Avoid_: Queue position, waiting stage, capacity wait

     **Turn Record**:
     The file `turn.json` beside a Queue Owner's console log, in which a live owner records the Queue Turn of every item it has not finished with its PID, start-time identity and Review Capacity. `roundfix deliver status` reads it only for a Queue Owner it proves live, and the Run Database gains no column for it (<ADR-0365 link>).
     _Avoid_: Turn column, owner log
     ````

   This text passed `TestGlossaryPhrasesAreDefinedTerms`,
   `TestGlossaryLinksNameAcceptedADRs`, `TestGlossaryAvoidLinesArePlainLists`
   and `TestGlossaryStatesEachRuleOnce` on a scratch tree with ADR-0365
   present.
3. MUST, in `docs/user-guide/commands/deliver.md`:
   - replace `Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, tokens <n|none>`
     with `Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, review concurrency <n>, tokens <n|none>`;
   - replace `Limits: deadline none, retries per item none, concurrency 1, tokens none`
     with `Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none`;
   - insert, as its own paragraph before the paragraph that starts "At or
     after the deadline, the owner parks each item that has not started",
     verbatim apart from line wrapping:

     ````markdown
     `concurrency 1` is the one item the Queue Owner advances at a time through
     every stage except the pre-PR review. `review concurrency <n>` is the Review
     Capacity, `delivery.review_concurrency` in User Config (`2` by default): the
     owner reviews up to that many items in `reviewing` at once, beside that item.
     The owner reads the value when it starts; stop and resume the owner to apply a
     change.
     ````

   - insert, after the paragraph that ends "the recorded owner PID and
     identity for `deliver resume` to recover.", verbatim apart from line
     wrapping:

     ````markdown
     An item's stage is the step it is at, not proof that the step is running: a
     Delivery Retry records `reviewing` before any review starts. A queued item
     starts its Implement Run only when no earlier item is between its Implement
     Run and its merge, and never ahead of an earlier queued item; an earlier item
     that is parked, merged or waiting only for its prerequisites does not hold
     it. A live owner finds an item a retry returned to `reviewing` within 30
     seconds and reviews it in a free review slot, even while it is inside another
     item's Implement Run; a reviewed item then waits for the owner's turn for its
     archive, gate, publication and merge.

     For a Queue Owner it proves live, status reads the owner's Turn Record,
     `turn.json` beside its console log, and prints one line per item that is
     neither `merged` nor `parked`, in queue order, after any `Warning:` lines and
     before any `Park:` lines:

     ```text
     Turn: <slug> advancing at <stage>
     Turn: <slug> reviewing; <held> of <capacity> review slot(s) held
     Turn: <slug> waiting for a review slot (<n> of <m> waiting); <held> of <capacity> review slot(s) held
     Turn: <slug> waiting for prerequisites: <slug>, <slug>
     Turn: <slug> waiting for the owner's turn (<n> of <m> waiting)
     Turn: <slug> waiting for the owner to read it
     ```

     `<n> of <m> waiting` is the item's place in line: for a review slot by queue
     position, for the owner's turn in the order the owner will visit the items,
     starting after the item it advances. The last form names an item whose stage
     changed after the owner last wrote its record, such as one just retried. An
     owner that is not proven live, a record of another PID or identity, or no
     record prints no `Turn:` line; a record that cannot be read writes
     `roundfix: Delivery Queue turn record <path> could not be read: <error>` to
     stderr, and status still exits `0`.
     ````

4. MUST, in `docs/user-guide/configuration.md`, add after the
   `delivery.check_timeout` row of the settings table the row

   `| `delivery.review_concurrency` | `2` | Sets the Review Capacity: how many Delivery Queue items the Queue Owner reviews at once, beside the one item it advances through every other stage. An integer of at least `1`. User Config only; Project Config is ignored with a warning. |`

   and insert before the `## Linux Lane` heading a new second-level heading
   `Delivery Queue review concurrency` followed by this text, verbatim apart
   from line wrapping:

     ````markdown
     `delivery.review_concurrency` sets the Review Capacity of a Queue Owner: the
     number of items in `reviewing` whose pre-PR review it runs at once, beside the
     one item it advances through every other stage. It is an integer, default
     `2`, at least `1`; `0` or a negative value fails configuration validation with
     `delivery.review_concurrency must be greater than 0`. For example:

     ```yaml
     delivery:
       review_concurrency: 3
     ```

     Each review reads only its own item worktree in a read-only session, so the
     reviews and the owner's Implement Run never write the same files. Reviews
     count neither against `runs.max_active` nor a Run's `verification.concurrency`,
     and their tokens stay outside the queue's token ceiling. Only User Config can
     set it; a Project Config value is ignored with this warning:
     `config: delivery.review_concurrency in Project Config is ignored; set
     delivery.review_concurrency in User Config`. A Queue Owner reads the value
     when it starts, and `roundfix deliver start` and `roundfix deliver status`
     print it on their `Limits:` line.
     ````

5. MUST, in `.agents/skills/roundfix/references/deliver.md`, replace
   `Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, spend not measured`
   with `Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, review concurrency <n>, tokens <n|none>`,
   and insert after the paragraph that ends "the recorded owner PID and
   identity for `deliver resume` to recover.", verbatim apart from line
   wrapping:

     ````markdown
     A stage is the step an item is at, not proof that the step is running: a
     Delivery Retry records `reviewing` before any review starts. The Queue Owner
     advances one item at a time through every stage except the pre-PR review,
     which it runs for up to its Review Capacity of items at once
     (`delivery.review_concurrency`, User Config, `2` by default, printed on the
     `Limits:` line as `review concurrency <n>`); a queued item starts only when no
     earlier item is between its Implement Run and its merge. For a live owner,
     status prints one `Turn:` line per item that is neither `merged` nor `parked`,
     after any `Warning:` lines and before any `Park:` lines, from the owner's Turn
     Record: `advancing at <stage>`, `reviewing; <held> of <capacity> review
     slot(s) held`, `waiting for a review slot (<n> of <m> waiting); <held> of
     <capacity> review slot(s) held`, `waiting for prerequisites: <slugs>`,
     `waiting for the owner's turn (<n> of <m> waiting)` or `waiting for the owner
     to read it`. Read the `Turn:` lines before acting on a stage that looks
     stuck: an item that waits for a review slot or the owner's turn needs no
     retry and no stop.
     ````

   Never touch the `### QA settlement` section. Then run `make skills-sync`
   and `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
6. MUST, as the LAST step after the last guide edit, re-record the Behavior
   Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST search `docs/user-guide/`, `.agents/skills/`, `skills/` and
   `CONTEXT.md` for `concurrency 1, tokens` and `spend not measured`, and
   leave neither outside `docs/history/`.
8. MUST NOT change any Go file, any other skill or reference, any other
   `CONTEXT.md` entry, any Archive Record, the Makefile, CI workflows,
   `go.mod` or `go.sum`.

## Subtasks

- [ ] Write ADR-0365.
- [ ] Add the three glossary entries and revise the Queue Owner entry.
- [ ] Update the `deliver` guide and the configuration guide.
- [ ] Update the Roundfix `deliver` reference, sync the mirror and record the skill version.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] ADR-0365 is accepted and `CONTEXT.md` defines the three new terms and the revised Queue Owner with links to it.
- [ ] The `deliver` guide, the configuration guide and both copies of the Roundfix `deliver` reference state the stage meaning, the `Turn:` lines, the `Limits:` line and the key.
- [ ] The skill version record and the Behavior Surface Record are current.

## Context

- creates: `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- The glossary, the guides, the Roundfix skill and its records — consumers: `./internal/docscontract`, `./skills`

## Verification

- `test -f docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md || { echo 'missing ADR-0365' >&2; exit 1; }; grep -q '^status: accepted$' docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md || { echo 'ADR-0365 not accepted' >&2; exit 1; }; for term in '**Review Capacity**:' '**Queue Turn**:' '**Turn Record**:' '**Queue Owner**:'; do grep -qF -- "$term" CONTEXT.md || { echo "missing term $term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'one item at a time through every stage except the pre-PR review' || { echo 'Queue Owner entry not revised' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'concurrency 1, review concurrency <n>, tokens <n|none>' || { echo "missing Limits line in $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'waiting for the owner to read it' || { echo "missing Turn line in $f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'delivery.review_concurrency must be greater than 0' || { echo 'missing key section in configuration.md' >&2; exit 1; }; grep -q '^## Delivery Queue review concurrency$' docs/user-guide/configuration.md || { echo 'missing heading in configuration.md' >&2; exit 1; }; cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || { echo 'deliver reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce|TestGlossaryClosedListsEqualTheirConstants|TestTheSkillCoverageMapIsCurrent)$' && go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — expected: exit 0. Before this Task ADR-0365 does not exist, so the first check fails; afterwards the ADR, the four glossary entries, the guides and both copies of the reference carry the prescribed text, the glossary contracts pass, and the skill version record and the Behavior Surface Record are current.

## References

- `_prd.md` → Goals; Core Feature 1; Core Feature 8; Success Metric 6; Glossary
- `_techspec.md` → ADR-0365; API Contract 2; API Contract 6; API Contract 8; Coverage Map; Build Order 1

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implementation evidence:

- Created ADR-0365 with accepted lifecycle frontmatter and the prescribed decision, alternatives and consequences.
- Replaced the Queue Owner glossary entry and added Review Capacity, Queue Turn and Turn Record with accepted ADR links.
- Updated the `deliver` and configuration guides with the Review Capacity key, stage meaning, limits and `Turn:` contract.
- Updated the canonical Roundfix `deliver` reference; `make skills-sync` synchronized its mirror; the sanctioned version recorder raised the Roundfix skill pair to `0.1.191` and updated `skills/testdata/owned-skill-versions.json`.
- Re-recorded `docs/references/behavior-surfaces.json` as the last guide-related regeneration step.

Focused checks:

- `make skills-sync`: passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed (`1 passed in 1 packages`).
- `rg -n "concurrency 1, tokens|spend not measured" docs/user-guide .agents/skills skills CONTEXT.md`: no matches (exit 1 as expected for the required absence check).
- `git diff --check`: passed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed (`1 passed in 1 packages`).

Acceptance evidence for Daemon Verification:

- ADR/glossary criterion: the new accepted ADR and four glossary entries are present with the required links and terms.
- Guide/reference criterion: the `Limits:` lines, stage explanation, `Turn:` lines, and `delivery.review_concurrency` documentation are present in all declared surfaces; the canonical and mirrored references are synchronized by `make skills-sync`.
- Record criterion: the owned-skill version test and Behavior Surface Record regeneration both passed.

Not run: the commands in this Task's `## Verification` section; the Daemon owns those checks and Task status.

## Carry-forward provenance

- Source Run: `run_20261105T203527Z_b348492260e5bafb`
- Source commit: `6977e95ac77d48957b0826e873e9439b70e4ba10`
