---
task: task_04
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: pending
type: docs
complexity: low
---

# Task 04: The history guide, the Roundfix skill and the glossary describe the readings

## Overview

task_01 to task_03 change what a conversion writes; task_02 documented the
field fallbacks, and the `history` user guide still says promoted files are
copied byte-identically. This Task brings
the guide, the two Roundfix skill references that cover the changed sources
(`command history sanitize` and `command baseline update` in the Skill
Coverage Map) and the glossary up to date. It is verifiable on its own by its
phrases, the mirror check, the skill version record and the Behavior Surface
Record.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [history sanitize records lose their summary, links and dates](references/2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md),
   for the documentation its six fixes need.
2. MUST describe in `docs/user-guide/commands/history.md`, beside the field
   fallbacks task_02 documented: the summary reading and the Legacy Link
   Rebase of the outcome and of each promoted copy, replacing the sentence
   that says `--promote` copies each file byte-identically; and that records
   written before this version keep their bytes. The guide MUST contain the
   phrase `Legacy Link Rebase` and MUST NOT contain
   `copies each file byte-identically`.
3. MUST add the same readings, with both phrases, to the History Sanitize
   Command section of `.agents/skills/roundfix/references/archive.md`, and to
   the `--promote` paragraph of `.agents/skills/roundfix/references/baseline.md`
   a sentence that the update's conversions read the same way and rebase each
   promoted copy, with the phrase `Legacy Link Rebase`. No text inside any
   `### QA settlement` section may change.
4. MUST run `make skills-sync`, raise and record the Roundfix skill's version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Measured on a scratch copy on 2026-10-13, these rewrite exactly the eight
   files this Task declares besides the guide and `CONTEXT.md`; no record or
   mirror is hand-edited.
5. MUST, through `domain-modeling`, add **Legacy Link Rebase** to `CONTEXT.md`
   near **Archive Relink**, linked to ADR-0277: the rewrite a Legacy Archive
   Folder conversion makes to the relative destinations of the record's
   outcome and of each promoted copy, pointing a folder file at its promoted
   copy or the Archive Record and keeping every other destination on the same
   repository path. It MUST revise **Lenient Legacy Reading** to add the
   summary reading, `legacy: not recorded`, the sources and creation-date
   fallbacks and Legacy Link Rebase, citing ADR-0277.

## Subtasks

- [ ] Update the `history` user guide.
- [ ] Update the archive and baseline references of the Roundfix skill.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.
- [ ] Add and revise the two glossary entries.

## Acceptance Criteria

- [ ] The guide and the archive reference name `legacy: not recorded` and Legacy Link Rebase, and the guide no longer promises byte-identical copies.
- [ ] The baseline reference says the update's conversions rebase promoted copies.
- [ ] Mirrors, the skill version and the Behavior Surface Record are current.
- [ ] `CONTEXT.md` defines **Legacy Link Rebase** and the revised **Lenient Legacy Reading**.

## Consumers

None. This Task changes documentation, the skill text and its derived records, which no code consumes; the skills and docscontract checks it runs read them.

## Context

- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/history.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'legacy: not recorded' || { printf 'missing override value in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Legacy Link Rebase' || { printf 'missing Legacy Link Rebase in %s\n' "$f" >&2; exit 1; }; done; if tr -s '[:space:]' ' ' < docs/user-guide/commands/history.md | grep -qF -- 'copies each file byte-identically'; then printf 'guide still promises byte-identical copies\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Legacy Link Rebase**:' || { printf 'CONTEXT.md does not define **Legacy Link Rebase**\n' >&2; exit 1; }; grep -A1 '^[*][*]Lenient Legacy Reading[*][*]:' CONTEXT.md | grep -qF -- 'legacy: not recorded' || { printf 'CONTEXT.md does not revise **Lenient Legacy Reading**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task neither the guide nor the skill names the new readings and the glossary lacks the term, so the command fails.

## References

- `_prd.md` → Core Feature 4; Glossary
- `_techspec.md` → API Contracts 1-5; Build Order 4; Glossary
- ADR-0277; ADR-0256; ADR-0244
