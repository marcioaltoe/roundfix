---
task: task_01
spec: 0266-the-review-record-survives-a-fix
status: pending
type: docs
complexity: low
---

# Task 01: The guides, the skill and the glossary name the resolved key and the Reviewed Head Record

## Overview

task_02 and task_03 change where a checkout's review record lives and what
the Delivery Retry's correction proof reads; this Task states that contract
first, in the same Pull Request. The review and deliver user
guides and the Roundfix skill's review and deliver references describe those
records (`docs/agents/cli.md`; ADR-0256 makes the skill files the covering
files of the `review` and `deliver` surfaces), and `CONTEXT.md` names the new
artifact. This Task writes `_techspec.md` → API Contract 5, raises the owned
skill's version with its record command, adds and revises the glossary terms
and re-records the Behavior Surface Record. It is verifiable on its own
through the phrases, the mirrors, the version record and the coverage
contract.

## Requirements

1. MUST add to `docs/user-guide/commands/review.md` and to
   `.agents/skills/roundfix/references/review.md`, directly after the
   sentence "Roundfix never reads another checkout's record.", the text of
   API Contract 5 for the review command, containing the phrases
   `named from the resolved checkout path` and `Reviewed Head Record`, the
   path `heads/<head commit>.json`, and the statement that a later round never
   replaces it and that a reused, blocked or omitted record writes none.
2. MUST rewrite, in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`, the moved-head
   `corrective-spec-required` paragraph and the matching row of the retry
   table to state API Contract 5 and API Contract 6 for the retry: the phrase
   `Reviewed Head Record of the parked candidate`; the widened path rule with
   the phrase `or among the paths its Tasks declared`, naming `creates:`,
   `interface:` and recorded paths read at the recorded revision; that review
   round 2 and the repository checks still follow; and that any other path
   still requires a corrective Spec, which the refusal names. MUST rewrite
   the last sentence of the archived-Specs paragraph of
   `docs/user-guide/commands/review.md` and
   `.agents/skills/roundfix/references/review.md` the same way. After this
   Task no copy of the four files says the correction "answers only the
   review in the archived Spec's own records".
3. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise and record the `roundfix` skill version. MUST NOT edit text
   inside any `### QA settlement` section, and MUST NOT change any other
   skill.
4. MUST add **Reviewed Head Record** to `CONTEXT.md` through
   `domain-modeling`, next to **Pre-PR Review Record**, linked to ADR-0272:
   the copy of a Pre-PR Review Record kept under the head it reviewed, written
   for every verdict that is not a reuse and never replaced by a later round
   at another head, which a Delivery Retry's correction proof reads for the
   parked candidate; `Avoid` synonyms such as round record or review history.
   MUST revise **Pre-PR Review Record** to say it is kept per checkout under
   its resolved path, with the phrase
   `kept per checkout under its resolved path`, linked to ADR-0272. Neither
   entry may reference a Spec.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`;
   the Skill Coverage Map itself does not change.
6. MUST search the user guides, `CONTEXT.md`, the agent guides and the
   Roundfix skill references for wording that says the review directory is
   named from the checkout path as spelled, that a review record is replaced
   by every review in the checkout, or that the retry admits only paths under
   the archived Spec, and change any it finds in this Task; the search on
   `main` at v0.76.0 found only the sentences Requirements 1 and 2 change.

## Subtasks

- [ ] Add the review and deliver text to both guides and both skill references.
- [ ] Sync the mirror and record the version raise.
- [ ] Add and revise the glossary terms and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guides and both copies of the skill references carry the new phrases.
- [ ] The `SKILL.md` pair matches and the owned-skill version record passes.
- [ ] `CONTEXT.md` defines **Reviewed Head Record** and the revised **Pre-PR Review Record**, and the Skill Coverage Map is current.

## Consumers

- The `roundfix` skill's review and deliver references and the owned-skill version record — consumers: `./skills`
- The Behavior Surface Record of the review and deliver guides — consumers: `./internal/docscontract`

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for file in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'named from the resolved checkout path' || { printf 'missing resolved key phrase in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'heads/<head commit>.json' || { printf 'missing head record path in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'Reviewed Head Record of the parked candidate' || { printf 'missing retry phrase in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'or among the paths its Tasks declared' || { printf 'missing declared-path rule in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do if tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'answers only the review in the archived Spec'; then printf 'old retry rule remains in %s\n' "$file" >&2; exit 1; fi; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1` — expected: exit 0. Before this Task the phrases are absent from every copy, so the command fails; afterwards the mirrors match and the owned-skill version record passes.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Reviewed Head Record**:' || { printf 'CONTEXT.md does not define **Reviewed Head Record**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'kept per checkout under its resolved path' || { printf 'CONTEXT.md does not revise **Pre-PR Review Record**\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the term and the revision are absent, so the command fails; afterwards the Skill Coverage Map and the Behavior Surface Record are current.

## References

- `_prd.md` → Goals; User Story 3; User Story 4; Core Feature 5; Core Feature 6
- `_techspec.md` → API Contract 5; API Contract 6; Coverage Map; Build Order 1
- ADR-0272; ADR-0256; ADR-0189; ADR-0233; ADR-0244
