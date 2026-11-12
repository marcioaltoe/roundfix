---
task: task_03
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: pending
type: docs
complexity: low
---

# Task 03: The guides, Roundfix skill references and glossary describe skills lock versions 1 and 2

## Overview

task_01 and task_02 make Roundfix read a `skills-lock.json` of version 1 or
2, hash in the version the lock declares, keep that version, and name
`roundfix upgrade` for any other version (vortex and pantheon, 2026-11-12,
[`references/2026-11-12-vortex-skills-lock-version-2-refused.md`](references/2026-11-12-vortex-skills-lock-version-2-refused.md),
reported by vortex;
[`references/2026-11-12-pantheon-skills-lock-version-2-refused.md`](references/2026-11-12-pantheon-skills-lock-version-2-refused.md),
reported by pantheon). This Task writes that into the `doctor`, `skills` and
`baseline` guides, the Roundfix skill's `setup`, `baseline` and `runtime`
references and the glossary entry **Stale Owned Lock Entry**, as
`_techspec.md` → API Contract 4 requires.

## Requirements

1. MUST add, in `docs/user-guide/commands/doctor.md`, a new paragraph right
   after the paragraph that ends "`roundfix baseline skills restore
   --profile <id> --skill <name>` replaces one deliberately.", with the text
   (inline code as shown):
   "Roundfix reads a `skills-lock.json` of version 1 or 2, the version the
   skills CLI 1.8 writes; every `computedHash` it writes into the lock or
   compares with an installed folder follows the hash algorithm of the
   version the file declares, and it never changes that version. Any other
   lock version fails the `skills:` line with the version named and
   `roundfix upgrade` as its only next action, because the skills command
   would write that version again."
2. MUST add, in `docs/user-guide/commands/skills.md`, a new paragraph right
   after the paragraph that ends "`Removed skills-lock.json ref: <name>` for
   each entry whose `ref` was removed.", and in the canonical
   `.agents/skills/roundfix/references/setup.md`, a new paragraph right after
   the paragraph that ends "`Removed skills-lock.json ref: <name>` for each
   entry whose `ref` was removed.", each with the first sentence of
   Requirement 1: "Roundfix reads a `skills-lock.json` of version 1 or 2, the
   version the skills CLI 1.8 writes; every `computedHash` it writes into the
   lock or compares with an installed folder follows the hash algorithm of
   the version the file declares, and it never changes that version."
3. MUST add, in `docs/user-guide/commands/baseline.md` right after the
   paragraph that ends "with one `- lock ref removed <name>` line after the
   lock-updated lines.", as a new paragraph, and in the canonical
   `.agents/skills/roundfix/references/baseline.md` right after the line that
   ends "line for each name whose `ref` was removed.", as the next lines of
   that paragraph, the first sentence of Requirement 1 followed by: "Any other
   lock version stops the update with the next action `run roundfix upgrade,
   then rerun roundfix baseline update`."
4. MUST add, in the canonical `.agents/skills/roundfix/references/runtime.md`,
   right after the sentence that ends "a mixed named failure prints the owned
   action followed by the skill-scoped external actions.", the text: "A lock
   of version 1 or 2 is read. Any other lock version fails the `skills:` line
   with the version named and `roundfix upgrade` as its only next action,
   because the skills command would write that version again."
5. MUST, in `CONTEXT.md`, through the domain-modeling skill, change in
   **Stale Owned Lock Entry** "a `computedHash` other than the installed
   copy's, so" to "a `computedHash` other than the hash of the installed copy
   in the algorithm of the lock version, so", and add, inside the same
   parentheses after its ADR-0359 link, ", " followed by the link
   `[ADR-0393](docs/adr/0393-roundfix-reads-skills-lock-versions-1-and-2-and-keeps-the-version-it-found.md)`.
6. MUST NOT change any other text of those files, any other guide, skill,
   reference, ADR, code or test, or text inside the `### QA settlement`
   section of any skill.
7. MUST run `make skills-sync` to write the mirrors
   `skills/roundfix/references/setup.md`, `skills/roundfix/references/baseline.md`
   and `skills/roundfix/references/runtime.md`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise both version fields of the Roundfix `SKILL.md` pair to the next
   version and record them in `skills/testdata/owned-skill-versions.json`;
   never hand-edit a mirror, a version field or the record.
8. MUST, as the last edit of the Task, after every guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Write the lock version text into the `doctor`, `skills` and `baseline` guides.
- [ ] Write it into the canonical `setup`, `baseline` and `runtime` references.
- [ ] Reword **Stale Owned Lock Entry** in the glossary.
- [ ] Sync the mirrors, record the skill version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The three guides and three references state that versions 1 and 2 are read, hashed in their own algorithm and kept, and name `roundfix upgrade` for any other version.
- [ ] The glossary's **Stale Owned Lock Entry** compares in the lock version's algorithm and links ADR-0393.
- [ ] The mirrors match the canonical references, the owned skill version record is current, and the Behavior Surface Record is current.

## Context

- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0393-roundfix-reads-skills-lock-versions-1-and-2-and-keeps-the-version-it-found.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- The Roundfix skill's `setup`, `baseline` and `runtime` references and their version record, and the glossary — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in docs/user-guide/commands/doctor.md docs/user-guide/commands/skills.md docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/setup.md .agents/skills/roundfix/references/baseline.md skills/roundfix/references/setup.md skills/roundfix/references/baseline.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; for p in 'of version 1 or 2, the version the skills CLI 1.8 writes;' 'follows the hash algorithm of the version the file declares, and it never changes that version.'; do printf '%s\n' "$t" | grep -qF -- "$p" || { echo "missing in $f: $p" >&2; exit 1; }; done; done; for f in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'as its only next action, because the skills command would write that version again.' || { echo "missing in $f: upgrade next action" >&2; exit 1; }; done; for f in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Any other lock version stops the update with the next action' || { echo "missing in $f: update next action" >&2; exit 1; }; done; for r in setup baseline runtime; do cmp -s ".agents/skills/roundfix/references/$r.md" "skills/roundfix/references/$r.md" || { echo "mirror differs: $r" >&2; exit 1; }; done` — expected: exit 0. Before this Task no guide or reference states the lock versions, so the first phrase check fails; afterwards all nine files carry the text and each mirror matches its canonical copy.
- `grep -F -- '**Stale Owned Lock Entry**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Stale Owned Lock Entry\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'other than the hash of the installed copy in the algorithm of the lock version' && grep -qF -- '[ADR-0393](docs/adr/0393-roundfix-reads-skills-lock-versions-1-and-2-and-keeps-the-version-it-found.md)' CONTEXT.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestGlossaryAvoidLinesArePlainLists)$'` — expected: exit 0. Before this Task the glossary entry compares with "the installed copy's" hash and does not link ADR-0393, so the second check fails; afterwards **Stale Owned Lock Entry** stays a definition with the new text and link, the owned skill version record and skill contracts pass, the Behavior Surface Record is current, and the glossary contracts pass.

## References

- `_prd.md` → Core Feature 4; Glossary
- `_techspec.md` → API Contract 4; Build Order 3
- ADR-0393; ADR-0359
