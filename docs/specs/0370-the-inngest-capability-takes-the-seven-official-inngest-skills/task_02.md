---
task: task_02
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
status: completed
type: docs
complexity: low
---

# Task 02: ADR-0206, the glossary and the guides say that inngest was replaced and retired

## Overview

ADR-0206 still says a skill the catalog removes is never followed to another
name, and the glossary, the `baseline` command guide and the Roundfix
skill's `baseline` reference list the Retired Skills without `inngest`. This
Task writes the amendment notice ADR-0380 requires, the **Retired Skill**
sentence, and the guide and skill text, so the documents describe what
task_01 delivered (`_prd.md` → Core Feature 4).

## Requirements

1. MUST answer the marcioaltoe-skills report of 2026-11-09 (ai-memory
   message 01a1c01f-b6ad, adopted in `references/`) for the documents an
   adopter reads.
2. MUST, in
   `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`,
   change `updated_at` to `2026-11-09T12:00:00Z` and insert, after the
   bullet that ends "`crafting-effective-readmes` leaves the `typescript`
   module." and before the paragraph that begins "Two alternatives were
   rejected.", this paragraph, exactly, followed by one blank line:

   ```markdown
   **Amended by [ADR-0380](0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md) (2026-11-09).** A skill the catalog removes and, in the same change, replaces with skills it names is followed to those replacements once the maintainer approves them, and becomes a Retired Skill, so an approved update deletes an adopter's unedited copy; `inngest` was followed to the seven `inngest-*` skills. A skill removed without a named replacement is still dropped and never followed by a similar name or purpose, and the two alternatives below stay rejected.
   ```

   No other line of the ADR changes.
3. MUST, through the `domain-modeling` skill, append to the definition line
   of **Retired Skill** in `CONTEXT.md`, after its link to ADR-0350, a space
   and this sentence, exactly:

   ```markdown
   A skill the catalog replaces with skills it names, as it replaced `inngest` with the seven `inngest-*` skills, becomes a Retired Skill, and the module that required it requires the replacements ([ADR-0380](docs/adr/0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md)).
   ```

   No other line of `CONTEXT.md` changes.
4. MUST, in `docs/user-guide/commands/baseline.md`, section `### Retired
   skills`, append to the sentence that ends "and the Roundfix-owned
   `implement-spec` since 0.127.0." a space and this sentence, exactly:

   ```markdown
   `inngest` joined them when the catalog replaced it with seven `inngest-*` skills, which the `backend` module requires instead.
   ```

   and add the line
   `  .agents/skills/inngest .claude/skills/inngest \` right after each of
   the two `.agents/skills/implement-spec .claude/skills/implement-spec`
   lines of the removal commands, one in the multi-line
   `git rm -r -q --ignore-unmatch` command and one in the multi-line
   `rm -rf` command.
5. MUST append the same sentence, after "`implement-spec` since 0.127.0.",
   to `.agents/skills/roundfix/references/baseline.md`, run
   `make skills-sync` to mirror it to `skills/roundfix/references/baseline.md`,
   then raise the Roundfix skill's version to the next one with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   MUST NOT change text inside the `### QA settlement` section of any skill.
6. MUST, as its last step, after every guide and skill edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change any Go file, Baseline asset, snapshot, other ADR, other
   guide or other skill; task_01 owns the code.

## Subtasks

- [ ] Write the ADR-0206 notice and the glossary sentence.
- [ ] Add the guide sentence and paths, the skill reference sentence, the mirror and the version.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] ADR-0206 carries the `Amended by [ADR-0380]` notice and `updated_at: 2026-11-09T12:00:00Z`.
- [ ] The **Retired Skill** definition names the replaced-skill rule and links ADR-0380.
- [ ] The guide and both copies of the skill reference name `inngest` as replaced, and the guide's removal commands include its paths.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md`
- interface: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill's `baseline` reference and its version record — consumers: `./skills`, `./internal/docscontract`
- The **Retired Skill** glossary definition and the ADR-0206 notice — consumers: `./internal/docscontract`

## Verification

- `a=docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md; grep -qF -- '**Amended by [ADR-0380](0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md) (2026-11-09).**' "$a" || { echo 'ADR-0206 lacks the ADR-0380 notice' >&2; exit 1; }; grep -qx 'updated_at: 2026-11-09T12:00:00Z' "$a" || { echo 'ADR-0206 updated_at not moved' >&2; exit 1; }; grep -F -- '**Retired Skill**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Retired Skill\*\*(:| —( |$))' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A skill the catalog replaces with skills it names, as it replaced' || { echo 'CONTEXT.md lacks the Retired Skill sentence' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce)$'` — expected: exit 0. Before this Task ADR-0206 has no notice, so the first check fails; afterwards the notice, the date and the glossary sentence are in place and the glossary contracts pass.
- `for p in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$p" | grep -qF -- 'joined them when the catalog replaced it with seven' || { printf 'missing sentence in %s\n' "$p" >&2; exit 1; }; done; test "$(grep -cF -- '.agents/skills/inngest .claude/skills/inngest' docs/user-guide/commands/baseline.md)" = 2 || { echo 'the removal commands lack the inngest paths' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task no guide names `inngest`; afterwards the guide, the skill reference and its mirror carry the sentence, the version record is current and the Behavior Surface Record matches.

## References

- `_prd.md` → Goal 4; Core Feature 4; Glossary
- `_techspec.md` → Documentation; Glossary; Build Order 2
- ADR-0380; ADR-0206

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, systematic-debugging

- Added the exact ADR-0380 amendment notice and set ADR-0206 `updated_at` to
  `2026-11-09T12:00:00Z`.
- Extended the **Retired Skill** glossary definition with the catalog-replacement
  rule and ADR-0380 link.
- Added the `inngest` replacement sentence and both removal paths to the
  `baseline` guide, then added the same sentence to the canonical Roundfix
  baseline reference and mirrored it with `make skills-sync`.
- Raised both Roundfix skill version fields to `0.1.208` and recorded the
  digest in `skills/testdata/owned-skill-versions.json` with the required
  version-record command.
- Re-recorded `docs/references/behavior-surfaces.json` with the required
  docscontract command after the guide and skill edits.

Focused checks:

- `git diff --check` — passed.
- `make skills-sync` — passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1` — passed (`ok roundfix/skills`); the first sandboxed attempt was blocked while writing the declared canonical skill path, then the same command passed with the required filesystem permission.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — passed (`ok roundfix/internal/docscontract`).

Acceptance evidence:

- ADR-0206 now carries the ADR-0380 notice and the required timestamp.
- `CONTEXT.md` names the replaced-skill rule and links ADR-0380 in the
  **Retired Skill** definition.
- The guide and both Roundfix skill reference copies name `inngest` as
  replaced, and the guide contains both `.agents/skills/inngest
  .claude/skills/inngest` removal entries.

Daemon Verification and Task settlement remain pending.
