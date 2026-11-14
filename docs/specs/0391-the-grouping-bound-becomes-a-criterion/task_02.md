---
task: task_02
spec: 0391-the-grouping-bound-becomes-a-criterion
status: pending
type: docs
complexity: medium
---

# Task 02: The authoring skills, the guide and the glossary state the grouping bound as a criterion

## Overview

The `write-idea`, `write-prd` and `write-techspec` skills stop bounding
grouping at four implementation Tasks and state the criterion; `write-tasks`
gains the rule that each Task earns its verification cycle, beside "Sized for
one fresh session"; the Roundfix skill's `spec` reference and the `spec`
command guide stop calling the bound a count. A new skills test pins each
phrase in both skill copies. The glossary redefines **Grouping Bound**, and
ADR-0208 carries ADR-0400's amendment notice.

## Requirements

1. MUST replace, in each canonical `.agents/skills/write-idea/SKILL.md`,
   `.agents/skills/write-prd/SKILL.md` and `.agents/skills/write-techspec/SKILL.md`,
   the grouping sentence under `## Sources that share a context` as
   `_techspec.md` → API Contract 2 says, word for word, and change nothing else
   in those files apart from the version fields of requirement 7.
2. MUST add the "**Earns its verification cycle.**" line of `_techspec.md` →
   API Contract 3, word for word, to the canonical
   `.agents/skills/write-tasks/SKILL.md` right after the "**Sized for one fresh
   session.**" line, keeping the vertical-slice rule, the Wave collision text
   and the Corrective-Task ceiling unchanged.
3. MUST change, as `_techspec.md` → API Contract 4 says, word for word, the
   grouping bound sentence of the canonical
   `.agents/skills/roundfix/references/spec.md` and the grouping sentence of
   `docs/user-guide/commands/spec.md`, and nothing else in either file. MUST
   NOT add or change text inside any `### QA settlement` section.
4. MUST add `TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle` to
   `skills/acceptance_rules_test.go`, calling `t.Parallel()` first, over both
   copies through `roundfixLayoutSources`, with the files, the phrases, the two
   absent phrases and the subtest names `_techspec.md` → Testing Approach lists
   for `skills`, comparing through `normalizeAcceptanceText`. No other test
   changes.
5. MUST change `CONTEXT.md` through `domain-modeling` as API Contract 5 says,
   word for word: the **Grouping Bound** definition and its `_Avoid_` line;
   MUST add no other capitalized phrase and no other entry.
6. MUST add the amendment notice of API Contract 6, word for word, where it
   places it, to `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`,
   set its `updated_at` to `2026-11-14T10:00:00Z`, and change nothing else in
   the ADR.
7. MUST then run, in this order and never by hand: `make skills-sync`, which
   copies every canonical skill file to `skills/`; the version record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the versions of `write-idea`, `write-prd`, `write-techspec`,
   `write-tasks` and `roundfix` from the tree the Task runs on (Specs 0392 and
   0393 raise some of them later; never write a number); and, as the Task's
   last step, after the last guide edit, the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Of the Governed Paths the Task may change only the canonical
   `SKILL.md` of `write-idea`, `write-prd`, `write-techspec`, `write-tasks` and
   `roundfix`, `.agents/skills/roundfix/references/spec.md`, and their
   `skills/` mirrors `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`,
   `skills/write-techspec/SKILL.md`, `skills/write-tasks/SKILL.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Reword the grouping sentence of the three authoring skills.
- [ ] Add the Task sizing rule to write-tasks.
- [ ] Reword the Roundfix spec reference and the spec command guide.
- [ ] Add the skills contract test.
- [ ] Redefine the glossary term and write the ADR-0208 notice.
- [ ] Sync the skills, record their versions, and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both copies of the three authoring skills, write-tasks and the Roundfix spec reference size a Spec by one context and by Tasks that earn their verification cycle, keep four Tasks as a reference, and none bounds grouping by a count.
- [ ] The spec command guide says the same, and the Behavior Surface Record is current.
- [ ] `CONTEXT.md` defines the changed **Grouping Bound**, ADR-0208 carries an "Amended by ADR-0400" notice dated 2026-11-14, and the glossary and documentation contracts pass.

## Context

- instruction: `docs/adr/0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/acceptance_rules_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`

## Consumers

- The authoring, write-tasks and Roundfix skill texts — consumers: `./skills`
- The glossary entry, the amended ADR, the spec command guide and its Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle ' || { echo 'missing pass: TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the skills test does not exist, so its pass line is missing; afterwards both copies of every changed skill carry each phrase once, the recorded versions match, and the skills suite passes.

- `grep -F -- '**Grouping Bound**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Grouping Bound\*\*(:| —( |$))' || { echo 'missing term: Grouping Bound' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'The limit on grouping sources into one Spec, set by two criteria and not by a count of Tasks' || { echo 'Grouping Bound unchanged' >&2; exit 1; }; f=docs/adr/0208-work-items-that-share-a-context-share-a-spec.md; test "$(grep -cF -- '**Amended by [ADR-0400](0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md) (2026-11-14).**' "$f")" -eq 1 || { echo "missing amendment: $f" >&2; exit 1; }; grep -q '^updated_at: 2026-11-14T10:00:00Z$' "$f" || { echo "updated_at unchanged: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'within the grouping bound, one context and Tasks that each earn their verification cycle' || { echo 'spec command guide unchanged' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary, the ADR and the guide still state the count; afterwards the term carries the criterion, ADR-0208 carries one dated notice, the guide carries the criterion, and the documentation, glossary and Behavior Surface Record contracts pass.

## References

- `_prd.md` → Goals 1 to 4; Core Feature 2; Core Feature 3; Success Metric 2; Success Metric 3; Success Metric 4; Glossary
- `_techspec.md` → Current behavior 4 to 6; API Contracts 2-6; Invariant 3; Testing Approach; Declared breaks; Build Order 2
- ADR-0400
