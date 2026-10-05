---
task: task_04
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: pending
type: docs
complexity: low
---

# Task 04: The Roundfix Skill and the reconcile and deliver guides describe merge evidence and item branches

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship the
Roundfix Skill update. The `reconcile` guide and the skill's `reconcile`
reference say today that "Other committed paths still require content
comparison", and the reference lists three debris candidate kinds; both
become wrong with ADR-0232. The `deliver` guide does not say the cleanup
proves a squash merge onto a moved default branch. This Task rewrites them
from the TechSpec, raises the skill's version and syncs the mirrors.

## Requirements

1. MUST state in `docs/user-guide/commands/reconcile.md` and
   `.agents/skills/roundfix/references/reconcile.md` the rule of
   `_techspec.md` → Merge evidence and API Contract 1: a Spec carries merge
   evidence when the default branch holds its archived `_prd.md` and the
   delivery commit that added it is not in the branch being proven; with it,
   Task commits still need their Task completed and every other commit is
   superseded, with the reason "superseded by the delivery of Spec"; without
   it the content comparison is unchanged. Each file MUST carry the phrases
   "merge evidence", "delivery commit" and "superseded by the delivery of
   Spec".
2. MUST state in both files the item branch candidates of API Contract 3 and
   Surface Transcript 1: only a full scan inspects
   `roundfix/deliver-<slug>-<16 hex>` branches, a live item's branch is
   preserved, `--apply` deletes a proven branch with its clean worktree, and
   the report's `itemBranchCandidates` list and its two `debrisSummary`
   counts. Each file MUST carry "itemBranchCandidates".
3. MUST keep in both files the dirty rule of ADR-0212 unchanged, stating that
   a dirty path outside the archived scope keeps the Run `dirty` and a dirty
   item worktree keeps its branch.
4. MUST remove from both files the sentence "Other committed paths still
   require content comparison", and from the reference the phrase "three
   debris candidate kinds".
5. MUST state in `docs/user-guide/commands/deliver.md`, beside the post-merge
   release, that the cleanup proves a squash merge onto a default branch that
   moved during the item by comparing the merge commit with Git's merge of the
   candidate into its first parent (API Contract 2), carrying the phrase "a
   default branch that moved during the item".
6. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
7. MUST NOT edit the `### QA settlement` section of any skill,
   `.agents/skills/roundfix/references/deliver.md`, any other command
   reference or guide, `commands.md`, `.roundfixrc.yml` or `CONTEXT.md`.

## Subtasks

- [ ] Rewrite the merged-Spec proof paragraph in the guide and the reference.
- [ ] Add the item branch candidates to both.
- [ ] Add the moved-default-branch sentence to the deliver guide.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The guide and the reference carry the four required phrases and
      neither carries the removed sentence; the reference no longer carries
      "three debris candidate kinds".
- [ ] The deliver guide carries "a default branch that moved during the item".
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `docs/user-guide/commands/deliver.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`

## Verification

- `for f in docs/user-guide/commands/reconcile.md .agents/skills/roundfix/references/reconcile.md; do t="$(tr -s '[:space:]' ' ' < "$f")" || exit 1; for p in "merge evidence" "delivery commit" "superseded by the delivery of Spec" "itemBranchCandidates"; do printf '%s\n' "$t" | grep -qF -- "$p" || { printf 'missing phrase in %s: %s\n' "$f" "$p" >&2; exit 1; }; done; for p in "Other committed paths still require content comparison" "three debris candidate kinds"; do if printf '%s\n' "$t" | grep -qF -- "$p"; then printf 'stale phrase in %s: %s\n' "$f" "$p" >&2; exit 1; fi; done; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "a default branch that moved during the item" || { printf 'missing phrase in deliver guide\n' >&2; exit 1; }` — expected: exit 0; before this Task neither reconcile file carries the new phrases and both carry the stale sentence, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/reconcile.md | grep -qF -- "superseded by the delivery of Spec" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirror does not carry the new phrase, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 4; Goals 1 to 4
- `_techspec.md` → Merge evidence; API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Vocabulary Contract; Build Order 4
- ADR-0187; ADR-0189; ADR-0232
