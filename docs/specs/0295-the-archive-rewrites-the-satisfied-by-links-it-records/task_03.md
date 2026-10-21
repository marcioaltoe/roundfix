---
task: task_03
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
status: pending
type: docs
complexity: low
---

# Task 03: The Roundfix skill references state both rewrites

## Overview

The Roundfix skill's `archive` reference covers the Archive Command and the
History Sanitize Command, its `baseline` reference covers the update's
history section, and its `spec` reference covers `SC-HISTORY-RETAINED`. After
task_01 and task_02 the archive rebases each `unproven` entry and the Record
Link Repair and the check name a promoted copy. This Task repeats the guides'
wording in the three references, raises the Roundfix skill version once, and
records it (`_prd.md` → Core Feature 4).

## Requirements

1. MUST edit `.agents/skills/roundfix/references/archive.md`, changing no
   text inside any `### QA settlement` section:
   - right after the sentence ending "so a reader learns what was never
     verified.", add: "Each entry is read as text of the Spec's `_prd.md`, so
     its relative links are rewritten like the outcome's, and an entry
     without one keeps its bytes (ADR-0308).";
   - change "In the copied outcome and each promoted copy, a destination
     inside the removed folder" to "In the copied outcome, each `unproven`
     entry and each promoted copy, a destination inside the removed folder";
   - change "Legacy Link Rebase rewrites the relative destinations in the
     outcome and each promoted copy," to "Legacy Link Rebase rewrites the
     relative destinations in the outcome, the `unproven` entries and each
     promoted copy,";
   - change "`record-links` unit, which rewrites those links to Git
     Locators." to "`record-links` unit, which rewrites a link to a file the
     record promoted to its `docs/references/` copy and every other such link
     to a Git Locator (ADR-0308)."
2. MUST edit `.agents/skills/roundfix/references/baseline.md`: in the
   paragraph that begins "After removal units, `record-links` repairs", right
   after the sentence ending "holding the target wins.", add the sentence
   task_02 added to the `baseline` guide: "Before that, a link whose target,
   resolved against the `source` folder, is a file whose copy
   `docs/references/<file>` the record lists under `promoted` and the tree
   holds is rewritten to a relative link to that copy, and the plan names
   `docs/references/<file>` as its rewrite (ADR-0308)."
3. MUST edit `.agents/skills/roundfix/references/spec.md`: right after the
   sentence ending "the deletion or replacement Git Locator.", add: "For a
   broken Archive Record link to a file the record promoted, the fix names
   the relative link to its `docs/references/` copy instead (ADR-0308)."
4. MUST run these commands, because this Task edits the Roundfix skill.
   Before them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   them itself before handing back; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version.
5. MUST NOT change any Go source, guide page, ADR or `CONTEXT.md`: the code,
   the guides and the glossary belong to task_01 and task_02.

## Subtasks

- [ ] Edit the `archive`, `baseline` and `spec` references.
- [ ] Raise the skill version, sync the mirror and record the version.

## Acceptance Criteria

- [ ] The three references state the rewrites in the canonical copy and its mirror.
- [ ] The Roundfix skill version is raised and recorded, and the mirror is in sync.

## Consumers

None. This Task changes skill text, its mirror and the skill version record only; no code consumes them.

## Context

- instruction: `docs/user-guide/commands/archive.md`
- instruction: `docs/user-guide/commands/history.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `docs/user-guide/commands/spec.md`
- instruction: `docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'In the copied outcome, each .unproven. entry and each promoted copy, a destination inside the removed folder' || { printf 'missing outcome sentence in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -q -- 'which rewrites a link to a file the record promoted to its .docs/references/. copy and every other such link to a Git Locator' || { printf 'missing record-links sentence in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'the record lists under .promoted. and the tree holds is rewritten to a relative link to that copy' || { printf 'missing repair sentence in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'For a broken Archive Record link to a file the record promoted, the fix names the relative link' || { printf 'missing fix sentence in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the sentences exists in the references, so the command fails at its first phrase check. Afterwards both copies of the three references state the rewrites, the mirror is in sync, the skill version is recorded, and the documentation contracts pass.

## References

- `_prd.md` → Core Feature 4; Behavior changes; Goals
- `_techspec.md` → System Architecture; API Contracts 2, 4 and 5; Build Order 3
- ADR-0308; ADR-0256
