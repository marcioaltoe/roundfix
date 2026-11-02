---
task: task_02
spec: 0341-a-glossary-definition-written-as-a-list-item-counts
status: pending
type: docs
complexity: low
---

# Task 02: The authoring guides state the forms and ask for the line-anchored glossary Verification

## Overview

The `write-tasks` skill asks for a glossary Verification that "checks the
bolded term followed by a colon in the glossary on whitespace-normalized
text", which misses Pantheon's em dash form and also passes on a definition
inside a sentence that the check rejects; the `write-prd` glossary guide and
the Roundfix `spec` reference never name a form. This Task makes the three
skills state the accepted Glossary Definition forms and makes `write-tasks`
ask for the command `roundfix spec check` prints, so a Task's Verification
and the check agree. It answers the Pantheon report of 2026-11-02 (ai-memory
message `01a19a04`, adopted as
[references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md](references/2026-11-02-a-glossary-definition-written-as-a-list-item-is-not-read.md)).

## Requirements

1. MUST replace, in the canonical `.agents/skills/write-tasks/SKILL.md`
   section `### Glossary requirement`, the sentence ending "checks the bolded
   term followed by a colon in the glossary on whitespace-normalized text"
   with text that (`_techspec.md` → API Contract 4):
   - names the accepted forms with the exact text of
     `speccheck.GlossaryDefinitionForms`;
   - asks for the Verification `roundfix spec check` prints for an unbound
     term, shown in a fenced block exactly as
     `speccheck.GlossaryVerificationCommand("<Term>", "CONTEXT.md")` renders
     it, with the term's regular-expression metacharacters escaped and the
     glossary path replaced;
   - says the command reads lines, not whitespace-normalized text, because a
     definition is a line, and that it is never more permissive than the
     check.
2. MUST add to the canonical `.agents/skills/write-prd/references/glossary.md`
   and `.agents/skills/roundfix/references/spec.md`, each under its existing
   glossary text, that a Glossary Definition is
   `speccheck.GlossaryDefinitionForms`, and that a heading, an indented or
   ordered list item and bold text inside a sentence define nothing.
3. MUST add `skills/glossary_definition_test.go` (package `skills_test`) with
   `TestTheGlossaryGuidesStateTheDefinitionForms` as `_techspec.md` →
   Testing Approach describes, calling `t.Parallel()` first.
4. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT touch Go code outside the new test, any guide under
   `docs/`, any Baseline module, the `domain-modeling` skill, the Makefile, CI
   workflows, `go.mod` or `go.sum`.
5. MUST copy the canonical skills with `make skills-sync`, then raise both
   version fields of the `roundfix`, `write-prd` and `write-tasks` `SKILL.md`
   pairs with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`,
   recording `skills/testdata/owned-skill-versions.json`, as the Task's last
   step. Never hand-edit a mirror or a version field.

## Subtasks

- [ ] Rewrite the `write-tasks` Glossary requirement.
- [ ] State the forms in the `write-prd` glossary guide and the Roundfix `spec` reference.
- [ ] Add the skills test.
- [ ] Sync the mirrors and record the versions.

## Acceptance Criteria

- [ ] Both copies of `write-tasks` carry the rendered command and the forms; both copies of the `write-prd` glossary guide and the Roundfix `spec` reference carry the forms.
- [ ] The mirrors equal their canonical copies and the owned-skill version record passes.

## Context

- instruction: `internal/speccheck/glossary.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-prd/references/glossary.md`
- interface: `skills/write-prd/references/glossary.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/glossary_definition_test.go`

## Consumers

- The `write-tasks`, `write-prd` and Roundfix skill texts and the owned-skill version record — consumers: `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheGlossaryGuidesStateTheDefinitionForms$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheGlossaryGuidesStateTheDefinitionForms ' || { echo 'missing pass: TestTheGlossaryGuidesStateTheDefinitionForms' >&2; exit 1; }; cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/write-prd/references/glossary.md skills/write-prd/references/glossary.md && cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the test does not exist, so the pass line is missing; afterwards the guides carry the forms and the command, the mirrors equal their canonical copies, and the owned-skill version record and the settlement guidance contracts pass.

## References

- `_prd.md` → Goals; Core Features 3 and 4; Success Metric 5
- `_techspec.md` → Current behavior 4-5; API Contracts 2 and 4; Testing Approach; Build Order 2
