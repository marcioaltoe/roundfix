---
task: task_03
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
status: completed
type: docs
complexity: low
---

# Task 03: The Roundfix archive reference and ADR-0269 state the refusal

## Overview

The Roundfix skill's `archive` reference lists the refusals of the Archive
Command and the Cancel Command, and ADR-0269 still says "Symbolic-link
destination parents stay accepted". After task_01 and task_02 both commands
refuse a path whose resolved target leaves the repository. This Task states
that refusal in the reference, raises the Roundfix skill version, and adds to
ADR-0269 the dated line naming ADR-0306 (`_prd.md` → Core Feature 4).

## Requirements

1. MUST add, in `.agents/skills/roundfix/references/archive.md`, under
   `## Archive Command` and right after the paragraph that begins "Archive
   also refuses while Git ignores any file under the Spec folder", this
   paragraph: "Archive and Cancel also refuse, writing nothing, a path whose
   resolved target lies outside the repository, such as a `docs/history`,
   `docs/references` or `docs/backlog` symbolic link to another directory; the
   refusal names the path. Move the link's target inside the repository or
   remove the link, then retry. A symbolic link whose target stays inside the
   repository, a checkout reached through a symbolic link and an external
   Spec Root keep working (ADR-0306)." No text inside any `### QA settlement`
   section may change.
2. MUST append to
   `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`,
   after its last line and separated by one blank line, one paragraph
   written on a single line, changing nothing else in the file. It starts with the bold lead "Narrowed
   by ADR-0306 (2026-10-20).", in which ADR-0306 is a Markdown link to its
   file name relative to `docs/adr/`, the form ADR-0260 uses for its
   "Extended by" line, and continues: "A symbolic-link destination parent
   whose resolved target stays in the repository, or in the external Spec
   Root's repository, stays accepted; one whose target leaves it is refused
   before anything is written."
3. MUST run these commands, because this Task edits the Roundfix skill.
   Before them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   them itself before handing back; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version.
4. MUST NOT change any Go source, guide page or other ADR: the code and the
   guides belong to task_01 and task_02.

## Subtasks

- [ ] Add the refusal paragraph to the Roundfix `archive` reference.
- [ ] Append the narrowing line to ADR-0269.
- [ ] Raise the skill version, sync the mirror and record the version.

## Acceptance Criteria

- [ ] The reference states the refusal and what keeps working, in the canonical copy and its mirror.
- [ ] ADR-0269 names ADR-0306 with its date and the narrowed rule.
- [ ] The Roundfix skill version is raised and recorded, and the mirror is in sync.

## Consumers

None. This Task changes documentation, the skill mirror and the skill version record only; no code consumes them.

## Context

- instruction: `docs/user-guide/commands/archive.md`
- instruction: `docs/user-guide/commands/cancel.md`
- instruction: `docs/adr/0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`

## Verification

- `for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Archive and Cancel also refuse, writing nothing, a path whose resolved target lies outside the repository' || { printf 'missing refusal paragraph in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md | grep -qF -- 'A symbolic-link destination parent whose resolved target stays in the repository, or in the external Spec Root' || { echo 'missing ADR-0269 narrowing line' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md | grep -q -- 'Narrowed by .ADR-0306.(0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md) (2026-10-20)' || { echo 'missing ADR-0269 narrowing line' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither the paragraph nor the narrowing line exists, so the command fails at its first phrase check. Afterwards both copies of the reference state the refusal, ADR-0269 names ADR-0306, the mirror is in sync, the skill version is recorded, and the documentation contracts pass.

## References

- `_prd.md` → Core Feature 4; Declared break; Goals
- `_techspec.md` → System Architecture; API Contract 4; Build Order 3
- ADR-0306; ADR-0269; ADR-0256

## Result

- Added the resolved-target refusal paragraph immediately after the ignored-file refusal in the canonical Roundfix archive reference; `make skills-sync` propagated it to the mirror. Focused checks confirmed both reference files are byte-identical and each contains the refusal text and the inside-repository, linked-checkout and external Spec Root behavior.
- Appended the single-line dated ADR-0306 narrowing to ADR-0269 with the relative Markdown link and the accepted/refused target rule.
- Raised both Roundfix skill version fields from `0.1.121` to `0.1.122`. The authorized `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` command passed (`1 passed in 1 packages`) and recorded the `0.1.122` digest in `skills/testdata/owned-skill-versions.json`.
- `git diff --check` passed. The declared `## Verification` commands remain for the Daemon.
