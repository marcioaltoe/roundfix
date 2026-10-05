---
task: task_04
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: docs
complexity: low
---

# Task 04: The skills, guides and version rule describe the three rules

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship the
Roundfix Skill update, and both Backlog Entries adopted on 2026-10-04
([queued Specs raise the same skill version](references/2026-10-04-queued-specs-raise-the-same-skill-version.md)
and [a docs fix after archive forces a manual Pull Request](references/2026-10-04-a-docs-fix-after-archive-forces-a-manual-pull-request.md))
name text that became wrong. This Task rewrites the `deliver` and `review`
guides and references, tells the `implement-task` agent that a command a
requirement names is part of its work, restates the repository's version
rule, and records both skill versions with the record command task_01 changes.

## Requirements

1. MUST state in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md` API Contract 3, carrying
   the phrase "every path changed from the candidate to the head lies under an
   archived Spec the blocker names", and API Contract 2, carrying the phrase
   "line-scoped", and MUST remove from both the sentence ending "still refuses
   a moved head" and the re-entry table row "after the item head moved",
   replacing that row with the review-only correction and the corrective-Spec
   refusal.
2. MUST state in `docs/user-guide/commands/review.md` and
   `.agents/skills/roundfix/references/review.md`, beside the corrective-Spec
   line, that the queue resumes without a corrective Spec when the correction
   answers only the review in the archived Spec's own records, carrying the
   phrase "answers only the review in the archived Spec's own records"; the
   review's stderr line itself does not change.
3. MUST add to `.agents/skills/implement-task/SKILL.md`, outside the
   `### QA settlement` section, that a command a requirement names, such as a
   generator or the owned-skill record command, is part of the work even when
   Verification runs the same test, and that its output is never written by
   hand, carrying the phrase "is part of the work even when Verification runs
   the same test".
4. MUST replace in `docs/agents/specific-repository.md` the version rule's
   "Raise both version fields, then record" with the record command that
   raises both version fields one patch above the highest recorded version,
   carrying the phrase "one patch above the highest recorded version".
5. MUST run `make skills-sync` so every mirror equals its canonical file, then
   run `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   yourself, as part of this Task's work, so it raises and records the
   Roundfix and `implement-task` skill versions; never edit a version field or
   the record by hand.
6. MUST change no Governed Path outside `.agents/skills/implement-task/SKILL.md`,
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md`,
   `.agents/skills/roundfix/references/review.md`,
   `docs/agents/specific-repository.md`, `skills/implement-task/SKILL.md` and
   `skills/roundfix/SKILL.md`, and MUST NOT edit any other skill, command
   guide, `commands.md`, `.roundfixrc.yml` or `CONTEXT.md`.

## Subtasks

- [ ] Rewrite the retry and derived-merge rules in the deliver guide and reference.
- [ ] Add the review-only correction to the review guide and reference.
- [ ] Add the record-command rule to implement-task and the repository version rule.
- [ ] Sync the mirrors and record both skill versions with the record command.

## Acceptance Criteria

- [ ] Each guide and reference carries its phrase, and neither deliver file
      carries the two removed passages.
- [ ] Each mirror equals its canonical file, and both raised versions are
      recorded by the record command.
- [ ] No Governed Path outside the seven named paths changes.

## Context

- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `docs/agents/specific-repository.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do t="$(tr -s '[:space:]' ' ' < "$f")" || exit 1; for p in "every path changed from the candidate to the head lies under an archived Spec the blocker names" "line-scoped"; do printf '%s\n' "$t" | grep -qF -- "$p" || { printf 'missing phrase in %s: %s\n' "$f" "$p" >&2; exit 1; }; done; for p in "still refuses a moved head" "after the item head moved"; do if printf '%s\n' "$t" | grep -qF -- "$p"; then printf 'stale phrase in %s: %s\n' "$f" "$p" >&2; exit 1; fi; done; done; for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "answers only the review in the archived Spec's own records" || { printf 'missing phrase in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- "is part of the work even when Verification runs the same test" || { printf 'missing phrase in implement-task\n' >&2; exit 1; }; t="$(tr -s '[:space:]' ' ' < docs/agents/specific-repository.md)" || exit 1; printf '%s\n' "$t" | grep -qF -- "one patch above the highest recorded version" || { printf 'missing phrase in specific-repository\n' >&2; exit 1; }; if printf '%s\n' "$t" | grep -qF -- "Raise both version fields, then record"; then printf 'stale version rule\n' >&2; exit 1; fi` — expected: exit 0; before this Task no file carries its new phrase and the deliver files and the version rule carry stale ones, so the command fails.
- `tr -s '[:space:]' ' ' < skills/implement-task/SKILL.md | grep -qF -- "is part of the work even when Verification runs the same test" && tr -s '[:space:]' ' ' < skills/roundfix/references/review.md | grep -qF -- "answers only the review in the archived Spec's own records" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirrors do not carry the new phrases, so the command fails; after it the mirrors equal their canonical files and both raised versions are recorded.

## References

- `_prd.md` → Goal 2; Core Feature 4
- `_techspec.md` → Vocabulary Contract; API Contract 2; API Contract 3; Build Order 4
- ADR-0233; ADR-0187; ADR-0189
