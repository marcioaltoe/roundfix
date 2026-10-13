---
task: task_03
spec: 0263-archive-and-cancel-keep-what-git-ignores
status: completed
type: docs
complexity: low
---

# Task 03: The skills and the glossary describe the refusal and the recovery

## Overview

task_01 and task_02 change what `roundfix archive` and `roundfix cancel`
refuse and print. This Task brings the Roundfix skill's archive reference, the
archive-spec skill and `CONTEXT.md` in line, so an Agent that archives or
cancels a Spec knows why the command refuses and how to finish a failed
removal. It is verifiable by phrase checks, the mirror check and the skill
version record.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [archive and cancel remove ignored files in the Spec folder](references/2026-10-12-archive-and-cancel-remove-ignored-files-in-the-spec-folder.md),
   for the skill text both halves need.
2. MUST describe, in `.agents/skills/roundfix/references/archive.md` under the
   Archive Command and the Cancel Command and in
   `.agents/skills/archive-spec/SKILL.md` under a new `## Ignored files and a
   failed removal` heading, that both commands refuse while Git ignores a file
   under the Spec folder, list each file and write nothing, with the phrase
   `remove each file or commit it with git add -f`; and that a failed folder
   removal keeps every write, exits 1 and prints the finishing command, with
   the phrase `chmod -R u+w`. No text inside any `### QA settlement` section
   may change.
3. MUST run `make skills-sync` and raise and record both skills' versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which rewrites the seven files this Task declares (measured on 2026-10-12);
   no record or mirror is hand-edited.
4. MUST add **Ignored Spec File** and **Folder Removal Recovery** to
   `CONTEXT.md` through `domain-modeling`, near **Archive Record** and
   **Cancel Command**, each linked to ADR-0269: an Ignored Spec File is a file
   under an active Spec folder that Git ignores, so no commit holds it, and the
   Archive Command and the Cancel Command refuse while one exists; Folder
   Removal Recovery is how both commands end when the final removal of the Spec
   folder fails, keeping every write and printing the command that finishes the
   removal.

## Subtasks

- [ ] Describe the refusal and the recovery in the Roundfix archive reference.
- [ ] Describe them in the archive-spec skill under its own heading.
- [ ] Sync the mirrors and record both skill versions.
- [ ] Add the two glossary terms.

## Acceptance Criteria

- [ ] Both skills name the refusal and the finishing command, and their mirrors match.
- [ ] Both skill versions are raised and recorded.
- [ ] `CONTEXT.md` defines both terms and links ADR-0269.

## Consumers

- Owned skill versions of `roundfix` and `archive-spec` — consumers: `./skills`

## Context

- instruction: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/roundfix/references/archive.md .agents/skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'remove each file or commit it with git add -f' || { printf 'missing refusal phrase in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'chmod -R u+w' || { printf 'missing finishing command in %s\n' "$f" >&2; exit 1; }; done; grep -q '^## Ignored files and a failed removal$' .agents/skills/archive-spec/SKILL.md || { printf 'archive-spec lacks its heading\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Ignored Spec File**:' || { printf 'CONTEXT.md does not define **Ignored Spec File**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Folder Removal Recovery**:' || { printf 'CONTEXT.md does not define **Folder Removal Recovery**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1` — expected: exit 0. Before this Task neither skill names the refusal and the glossary lacks both terms, so the command fails.

## References

- `_prd.md` → Core Feature 3; Glossary
- `_techspec.md` → API Contract 1; API Contract 3; Build Order 3
- ADR-0269; ADR-0256; ADR-0244

## Result

Implementation evidence:

- The Roundfix archive reference now documents the ignored-file refusal and
  forward folder-removal recovery under both the Archive Command and Cancel
  Command. The archive-spec skill has the new `## Ignored files and a failed
  removal` section. The required phrases are present in both source skills and
  their mirrors.
- `CONTEXT.md` now defines **Ignored Spec File** near **Cancel Command** and
  **Folder Removal Recovery** near **Archive Record**, with ADR-0269 links.
- `make skills-sync` exited 0 and regenerated the embedded mirrors.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` exited 0 (`1 passed in 1 packages`). It raised and
  recorded `archive-spec` as `0.0.11` and `roundfix` as `0.1.70`; the recorder
  wrote their digests to `skills/testdata/owned-skill-versions.json`.

Focused checks:

- Phrase, heading, glossary and ADR checks passed for the source files.
- The two source/mirror pairs compare byte-identically.
- The `### QA settlement` section in the archive-spec source remains
  unchanged from `HEAD`.

Acceptance-criterion evidence:

1. Both skills contain the refusal phrase and `chmod -R u+w`; each mirror
   matches its source.
2. Both skill versions are raised and recorded by the repository recorder.
3. `CONTEXT.md` defines both requested terms and links each to ADR-0269.

The implementation is ready for the Daemon's declared Verification and
settlement. Task status remains Daemon-owned.
