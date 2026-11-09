---
task: task_03
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
status: pending
type: docs
complexity: low
---

# Task 03: The guides and the authoring skills say how to mark an ADR not applicable and what the new finding means

## Overview

Tasks 01 and 02 add `SC-UNREACHABLE-UNREAD` and make an ADR whose TechSpec
ADR Alternatives line begins with `not applicable` stop drawing in the ADRs
that cite it. This Task writes both facts where authors read them: the
`spec` command guide, the context-driven development guide, the Roundfix
skill's `spec` reference, and the `write-techspec` and `write-tasks` skills.
It answers the adopted Backlog Entry of 2026-11-07 and the Fiscus finding of
2026-11-07 on the documentation side, raises the three skills' versions by
the record command and re-records the Behavior Surface Record last.

## Requirements

1. MUST, in `docs/user-guide/commands/spec.md`, replace the table row that
   begins "| `SC-ADR-ALTERNATIVES-UNRELATED` |" with these two rows, exactly
   as the block holds them:

   ```markdown
   | `SC-ADR-ALTERNATIVES-UNRELATED` | A listed ADR has no line relating the design to that ADR's rejected alternatives. Runs from the TechSpec stage. A line whose relation begins with `not applicable` keeps the ADR listed, and `SC-ADR-RELATED` then reports no ADR that cites only ADRs marked that way. |
   | `SC-UNREACHABLE-UNREAD` | The PRD's `## Unreachable Acceptance` section holds text that no fielded declaration reads, such as a numbered item, a prose paragraph or an unindented line inside a declaration, or a declaration lacks its `criterion:`, `reason:` or `satisfied-by:` field. One finding per run of consecutive unread lines, at its first line, and the fix shows the fielded form. Blank lines, HTML comments and a first paragraph that begins with `None.` in a section that holds no declaration are read. An error from the PRD stage; a missing PRD skips it with its path. |
   ```

2. MUST, in `docs/user-guide/context-driven-development.md`, replace the
   paragraph that begins "The Spec Consistency Check uses an ADR horizon:"
   with these two paragraphs, exactly:

   ```markdown
   The Spec Consistency Check uses an ADR horizon: `SC-ADR-RELATED` reports an ADR
   for a committed Spec only when the commit that added the ADR is an ancestor of
   the commit that added the Spec's `_prd.md`. Only an ADR the Spec applies or
   changes draws in the ADRs that cite it: an ADR whose line under the TechSpec's
   `## ADR Alternatives` begins with `not applicable` stays listed, and an ADR that
   cites only such ADRs opens no `SC-ADR-RELATED` gap.

   The check reads the PRD's `## Unreachable Acceptance` section as the QA gate
   does, and `SC-UNREACHABLE-UNREAD` names, at its line, any text in it that no
   fielded declaration reads.
   ```

3. MUST, in `.agents/skills/roundfix/references/spec.md`, append right after
   the two-line bullet that begins "- `SC-ADR-ALTERNATIVES-UNRELATED` — a
   listed ADR has no line relating the" these lines, exactly:

   ```markdown
     A line whose relation begins with `not applicable` keeps the ADR listed,
     and `SC-ADR-RELATED` reports no ADR that cites only ADRs marked that way.
   - `SC-UNREACHABLE-UNREAD` — the PRD's `## Unreachable Acceptance` section
     holds text that no fielded declaration reads, or a declaration lacks one of
     its fields; the fix shows the `- criterion:` form. An error from the PRD
     stage.
   ```

   The first two lines continue the existing bullet, indented by two spaces.
4. MUST, in `.agents/skills/write-techspec/SKILL.md`, append to the
   "Rejected alternatives." bullet, right after its line
   "`SC-ADR-ALTERNATIVES-UNRELATED` reports a listed ADR without its line.",
   these lines, indented by two spaces like the bullet:

   ```markdown
     For an ADR the Spec lists only to record why it does not apply, begin the
     relation with `not applicable`: the ADR stays listed, and `SC-ADR-RELATED`
     does not ask for the ADRs that cite only it.
   ```

   and, in `.agents/skills/write-techspec/references/techspec-template.md`,
   replace the comment under `## ADR Alternatives` with:

   ```markdown
   <!-- One line for each ADR listed under Active ADR obligations. State how the design
        relates to the alternatives that ADR rejected. Begin the relation with
        `not applicable` for an ADR the Spec lists only to record why it does not apply. -->
   ```

5. MUST, in `.agents/skills/write-tasks/SKILL.md`, replace the sentence
   "A `[gap]` is blocking too; account for the named ADR or record why it
   does not apply." with "A `[gap]` is blocking too; account for the named
   ADR or record why it does not apply, and begin a not-applying ADR's line
   under the TechSpec's `## ADR Alternatives` with `not applicable` so it
   draws in no ADR that cites only it.", keeping the line width of the
   paragraph.
6. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT change any other skill or any Roundfix reference other
   than `spec`.
7. MUST copy the canonical skills to their mirrors with `make skills-sync`,
   then raise and record the versions of `roundfix`, `write-techspec` and
   `write-tasks` with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   no version field or record is hand-edited.
8. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's LAST step, after the last edit of either guide. No record is
   hand-edited.

## Subtasks

- [ ] Add the two rows to the `spec` command guide and the paragraphs to the context-driven development guide.
- [ ] Add the bullets to the Roundfix `spec` reference.
- [ ] Add the marker to the `write-techspec` skill, its template and the `write-tasks` skill.
- [ ] Sync the mirrors, record the skill versions, then re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both guides and the Roundfix `spec` reference name `SC-UNREACHABLE-UNREAD` and the `not applicable` marker.
- [ ] The `write-techspec` skill, its template and the `write-tasks` skill tell authors to begin a not-applying ADR's line with `not applicable`.
- [ ] The mirrors, the version record and the Behavior Surface Record are current, and the skill and documentation contracts pass.

## Context

- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The guides, the owned skills and their records — consumers: `./skills`, `./internal/docscontract`

## Verification

- `grep -qE -- '^[|] .SC-UNREACHABLE-UNREAD. [|] The PRD' docs/user-guide/commands/spec.md || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Only an ADR the Spec applies or changes draws in the ADRs that cite it' || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'names, at its line, any text in it that no fielded declaration reads' || exit 1; grep -qE -- '^- .SC-UNREACHABLE-UNREAD. — the PRD' skills/roundfix/references/spec.md || exit 1; tr -s '[:space:]' ' ' < skills/write-techspec/SKILL.md | grep -qF -- 'does not ask for the ADRs that cite only it' || exit 1; tr -s '[:space:]' ' ' < skills/write-techspec/references/techspec-template.md | grep -qF -- 'for an ADR the Spec lists only to record why it does not apply' || exit 1; tr -s '[:space:]' ' ' < skills/write-tasks/SKILL.md | grep -qF -- 'so it draws in no ADR that cites only it' || exit 1; go test -count=1 -timeout 20m ./skills || exit 1; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide nor any skill names the new code or the marker, so the first check fails; afterwards the text is in every canonical copy and mirror, the versions are recorded, and the skill and documentation contracts pass, the skill coverage check included.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3
- `_techspec.md` → Documentation; Build Order 3
- ADR-0377; ADR-0265
