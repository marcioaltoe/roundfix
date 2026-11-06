---
task: task_01
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: pending
type: docs
complexity: medium
---

# Task 01: ADR-0369, its notices, CONTEXT.md, the guides and the Roundfix deliver reference state the rule

## Overview

task_02 to task_04 make the checking stage wait for the Head Check Set and
name what it waits for. This Task runs first and records that rule in
ADR-0369, writes the amendment notices into ADR-0311 and ADR-0365, adds
**Head Check Set** and **Expected Workflow** to `CONTEXT.md` and revises
**Check Wait** and **Queue Turn**, and writes the rule into the `deliver`
guide, the configuration guide and the Roundfix skill's `deliver`
reference. It answers the Fluxus report of 2026-11-05, reported by fluxus
([references](references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md)):
pull requests #466, #469 and #477 merged before `CI | Backend` ran. The code
Tasks that follow implement exactly the behavior these documents state.

## Requirements

1. MUST create `docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md`
   with exactly the content of `_techspec.md` → ADR-0369, apart from line
   wrapping.
2. MUST write the two paragraphs of `_techspec.md` → Amendment notices, each
   verbatim apart from line wrapping and as its own paragraph right after the
   paragraph it names, into
   `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
   and `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`,
   and set `updated_at: 2026-11-06T12:00:00Z` in both, changing nothing else
   in them.
3. MUST, through the domain-modeling skill, change `CONTEXT.md` so that,
   where `<ADR-0369 link>` stands for
   `[ADR-0369](docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md)`
   and `<ADR-0311 link>` for the ADR-0311 link the current `**Check Wait**`
   entry already carries, unchanged:
   - the `**Check Wait**` entry is exactly, line for line:

     ````markdown
     **Check Wait**:
     The bounded wait of a Delivery Queue item in `checking` for its Head Check Set and a mergeable state, and of its publication for an existing pull request to report the candidate head. Each wait lasts `delivery.check_timeout` (User Config, 60 minutes by default) from its own start, restarts once after a check re-run, and ends in a `checks-timeout` park that names the key, the duration and the members still pending (<ADR-0311 link>, <ADR-0369 link>).
     _Avoid_: CI timeout, check deadline
     ````

   - right after the `_Avoid_:` line of that entry, separated by one blank
     line each, come these two new entries, line for line:

     ````markdown
     **Head Check Set**:
     What a Delivery Queue item in `checking` waits for before it merges: every check run and commit status listed for its pull request head, required or not, the latest Actions workflow run of each workflow file and event of the head commit, and each Expected Workflow. It is complete when it has a member, every member is completed and every Expected Workflow has a run; a member that ends other than `success`, `neutral` or `skipped` parks the item as `checks-failed`, naming the member and its link (<ADR-0369 link>).
     _Avoid_: Required checks, check suite, CI status

     **Expected Workflow**:
     A workflow file of the pull request head that GitHub runs for every pull request into the base branch: a `pull_request` or `pull_request_target` trigger with no path filter and no `branches-ignore`, whose `branches` and `types`, when present, name the base branch and both `opened` and `synchronize`, or a `push` trigger with no filter. The Head Check Set is not complete until each one has a workflow run for the head (<ADR-0369 link>).
     _Avoid_: Required workflow, triggered workflow
     ````

   - in the `**Queue Turn**` entry, `advancing in the owner's turn, reviewing
     in one of` becomes `advancing in the owner's turn, with the members of
     its Head Check Set it still waits for while in `checking`, reviewing in
     one of`, and its closing link list
     `([ADR-0365](docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md))`
     becomes
     `([ADR-0365](docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md), <ADR-0369 link>)`.

   This text passed `TestGlossaryPhrasesAreDefinedTerms`,
   `TestGlossaryLinksNameAcceptedADRs`, `TestGlossaryAvoidLinesArePlainLists`
   and `TestGlossaryStatesEachRuleOnce` on a scratch tree with ADR-0369
   present.
4. MUST, in `docs/user-guide/commands/deliver.md` and in
   `.agents/skills/roundfix/references/deliver.md`, replace the paragraph
   that starts "During `checking`, Roundfix waits while GitHub reports the
   Pull Request merge" (through "A second refusal for the same head parks
   `delivery-error`.") with this text, verbatim apart from line wrapping:

     ````markdown
     During `checking`, Roundfix waits for the head's Head Check Set: every check
     run and commit status the Pull Request lists for its head, required or not;
     the latest GitHub Actions workflow run of each workflow file and event whose
     head commit is the head; and each Expected Workflow. An Expected Workflow is a
     file under `.github/workflows/` at the head whose `pull_request` or
     `pull_request_target` trigger has no path filter and no `branches-ignore`,
     whose `branches`, when present, names the base branch exactly and no `!`
     pattern, and whose `types`, when present, name both `opened` and
     `synchronize`, or whose `push` trigger has no filter. An Expected Workflow
     without a run for the head is pending as `<path> (no run yet)`, so a head whose
     workflows have not registered yet is never taken for a head with nothing to
     wait for. A filtered workflow is never expected, but its run must complete
     once GitHub creates it. Check suites are not read. The item moves to `merging`
     only when the set has a member, every member is completed, and GitHub reports
     the merge state as neither `BLOCKED` nor `UNKNOWN`; a `CLEAN` merge state alone
     never merges.

     A member that ends `success`, `neutral` or `skipped` passes. A failed check
     goes through the re-run rules below; any other member that ends `failure`,
     `timed_out`, `cancelled`, `action_required`, `startup_failure`, `stale` or
     another conclusion parks the item as `checks-failed: <check> (<conclusion>)
     <link>`, each failed member joined by `; `. `<check>` is `<workflow> / <name>`
     for a check of a workflow, the check name for any other check, and the
     workflow name for a run. While the item waits, the owner log names the members
     it waits for:

     ```text
     roundfix: checks: Delivery Queue item <slug>: waiting for <n> check(s): <check>; <check>
     roundfix: checks: Delivery Queue item <slug>: <m> check(s) completed for <head>
     roundfix: checks: Delivery Queue item <slug>: read failed: <error>
     ```

     At the end of the Check Wait the item parks as `checks-timeout: current-head
     checks not settled within delivery.check_timeout <duration>; pending: <check>;
     <check>`, or with `; merge state <state>` in place of the pending list when
     only the merge state was awaited. If GitHub refuses the merge because
     `base branch policy prohibits the merge`, Roundfix returns to `checking` once
     for that head. A second refusal for the same head parks `delivery-error`.
     ````

5. MUST, in `docs/user-guide/commands/deliver.md`, also:
   - insert the line
     `Turn: <slug> advancing at checking; waiting for <n> check(s): <check>; <check>`
     right after `Turn: <slug> advancing at <stage>` in the block of `Turn:`
     forms;
   - insert the sentence "The checking form names the members of the Head
     Check Set the owner still waits for." right before "The last form names
     an item whose stage changed";
   - replace "When a required GitHub Actions check fails on its first
     attempt" with "When a GitHub Actions check fails on its first attempt".
6. MUST, in `.agents/skills/roundfix/references/deliver.md`, also replace
   "When a required Actions check fails on attempt one" with "When an
   Actions check fails on attempt one", and in the list of `Turn:` forms
   replace "Record: `advancing at <stage>`, `reviewing;" with "Record:
   `advancing at <stage>`, `advancing at checking; waiting for <n> check(s):
   <check>; <check>`, `reviewing;". Never touch the `### QA settlement`
   section. Then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
7. MUST, in `docs/user-guide/configuration.md`, replace "wait of an item in
   `checking` for its current-head checks and a mergeable state," with "wait
   of an item in `checking` for its Head Check Set and a mergeable state,",
   and replace "`checks-timeout: current-head checks not settled within
   delivery.check_timeout <duration>`, with the duration in Go's form such as
   `1h0m0s`." with "`checks-timeout: current-head checks not settled within
   delivery.check_timeout <duration>; pending: <check>; <check>`, with the
   duration in Go's form such as `1h0m0s` and the members of the Head Check
   Set still pending."
8. MUST, as the LAST step after the last guide edit, re-record the Behavior
   Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change any Go file, any other skill or reference, any other
   `CONTEXT.md` entry, any other ADR, any Archive Record, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Write ADR-0369 and the notices in ADR-0311 and ADR-0365.
- [ ] Add the two glossary entries and revise Check Wait and Queue Turn.
- [ ] Update the `deliver` and configuration guides.
- [ ] Update the Roundfix `deliver` reference, sync the mirror and record the skill version.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] ADR-0369 is accepted, ADR-0311 and ADR-0365 carry their notices, and `CONTEXT.md` defines Head Check Set and Expected Workflow and the revised Check Wait and Queue Turn with links to ADR-0369.
- [ ] The `deliver` guide, the configuration guide and both copies of the Roundfix `deliver` reference state the Head Check Set, the conclusions, the log lines, the blockers and the checking `Turn:` form.
- [ ] The skill version record and the Behavior Surface Record are current.

## Context

- creates: `docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md`
- interface: `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
- interface: `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`
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

- The glossary, the ADRs, the guides, the Roundfix skill and its records — consumers: `./internal/docscontract`, `./skills`

## Verification

- `f=docs/adr/0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md; test -f "$f" || { echo 'missing ADR-0369' >&2; exit 1; }; grep -q '^status: accepted$' "$f" || { echo 'ADR-0369 not accepted' >&2; exit 1; }; for adr in docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md; do grep -qF -- '**Amended by [ADR-0369](0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md) (2026-11-06).**' "$adr" || { echo "missing notice in $adr" >&2; exit 1; }; grep -q '^updated_at: 2026-11-06T12:00:00Z$' "$adr" || { echo "updated_at not bumped in $adr" >&2; exit 1; }; done; grep -F -- '**Head Check Set**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Head Check Set\*\*(:| —( |$))' || { echo 'missing term Head Check Set' >&2; exit 1; }; grep -F -- '**Expected Workflow**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Expected Workflow\*\*(:| —( |$))' || { echo 'missing term Expected Workflow' >&2; exit 1; }; grep -F -- '**Check Wait**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Check Wait\*\*(:| —( |$))' || { echo 'missing term Check Wait' >&2; exit 1; }; grep -F -- '**Queue Turn**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Queue Turn\*\*(:| —( |$))' || { echo 'missing term Queue Turn' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'in .checking. for its Head Check Set and a mergeable state' || { echo 'Check Wait not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'with the members of its Head Check Set it still waits for while in' || { echo 'Queue Turn not revised' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'a head whose workflows have not registered yet is never taken for a head with nothing to wait for' || { echo "missing Head Check Set paragraph in $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'waiting for <n> check(s): <check>; <check>' || { echo "missing waiting form in $f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'delivery.check_timeout <duration>; pending: <check>; <check>' || { echo 'configuration guide not updated' >&2; exit 1; }; cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || { echo 'deliver reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce|TestGlossaryClosedListsEqualTheirConstants|TestTheSkillCoverageMapIsCurrent)$' && go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — expected: exit 0. Before this Task ADR-0369 does not exist, so the first check fails; afterwards the ADR, the notices, the four glossary entries, the guides and both copies of the reference carry the prescribed text, the glossary contracts pass, and the skill version record and the Behavior Surface Record are current.

## References

- `_prd.md` → Goals; Core Feature 6; Success Metric 5; Glossary
- `_techspec.md` → ADR-0369; Amendment notices; API Contract 5; API Contract 6; Coverage Map; Build Order 1
