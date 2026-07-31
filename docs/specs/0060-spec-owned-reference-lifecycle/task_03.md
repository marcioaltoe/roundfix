---
task: task_03
spec: 0060-spec-owned-reference-lifecycle
status: completed
type: docs
complexity: medium
---

# Task 03: Rehearse the lifecycle end to end and prove the gates

## Overview

This Spec's gates are instructions, not compiled checks, so the only way to know
they catch their cases is to run them against a case built to fail. Rehearse the
whole transition on throwaway documents, prove each gate fires and each
non-case does not, and record the evidence.

## Requirements

1. MUST rehearse on a scratch branch with throwaway documents — one finding and
   one inbox note — and MUST NOT commit the rehearsal fixtures.
2. MUST prove the move preserves history: `git log --follow` on the moved file
   reaches its pre-adoption commits.
3. MUST prove the index populates with the provenance column holding the
   pre-adoption path and the path column resolving from the index.
4. MUST prove no repository link points at the pre-adoption path afterward,
   counting occurrences before and after rather than sampling one.
5. MUST prove the archive gate passes for the self-contained Spec.
6. MUST prove the archive gate fails, naming the offending link, when a stale
   link is injected.
7. MUST prove the archive gate does **not** fail on prose that names
   `docs/findings/` without linking into it — the false-positive case.
8. MUST prove the archive gate passes trivially for a Spec carrying no index,
   which is the migration boundary.
9. MUST record every command and its outcome in this Task's `## Result`, and
   MUST leave the repository with the rehearsal fully discarded.

## Subtasks

- [ ] Build the throwaway finding and inbox note and run the adoption steps.
- [ ] Prove history, index, and link rewrite.
- [ ] Prove the archive gate on all four cases: pass, stale link, prose
      mention, no index.
- [ ] Discard the rehearsal and confirm a clean tree.

## Acceptance Criteria

- [ ] `git log --follow` on the moved document reaches its pre-adoption history.
- [ ] The occurrence count for the pre-adoption path is zero outside the index's
      provenance column and Git history.
- [ ] The archive gate passes for the self-contained Spec.
- [ ] The archive gate fails on the injected stale link and names it.
- [ ] The archive gate passes on the prose-mention case.
- [ ] The archive gate passes on a Spec with no index.
- [ ] `git status --porcelain` is empty of rehearsal artifacts at the end.

## Context

- interface: `docs/specs/0060-spec-owned-reference-lifecycle/task_03.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`

## Verification

- `make verify` — expected: exit 0.
- `git status --porcelain` — expected: no rehearsal fixture remains.

## References

`_prd.md` → Success Metrics; `_techspec.md` → Testing Approach, Risks (a gate
written as prose can be skipped silently; the archive gate could fire on prose).

## Result

### Contract blocker

- The required history proof cannot coexist with the no-commit constraints.
  `git log --follow` reads committed history; an uncommitted `git mv` makes the
  destination absent from every commit, so a log starting at the moved path has
  nothing to traverse. Satisfying that criterion requires committing the move
  on the scratch branch, while Requirement 1 and the Daemon execution invariants
  prohibit committing rehearsal fixtures or making any commit in this turn.
- No lifecycle rehearsal or archive-gate case was run after this contradiction
  was confirmed. The Task remains Daemon-owned and no acceptance checkbox is
  claimed.

### Focused preflight evidence

- `rtk git clone --quiet <current-worktree> repo` in a disposable directory —
  exit 0.
- `rtk git switch -c ma/spec-0060-task03-rehearsal` in the disposable clone —
  exit 0; created the required prefixed scratch branch without touching the Run
  branch.
- `rtk mkdir -p docs/specs/9999-task03-rehearsal/references` — exit 0.
- `rtk git mv docs/findings/2026-07-25-spec-owned-reference-lifecycle.md
  docs/specs/9999-task03-rehearsal/references/2026-07-25-spec-owned-reference-lifecycle.md`
  — exit 0.
- `rtk git diff --cached --summary` — exit 0; reported a 100% rename from
  `docs/findings/` into the scratch Spec's `references/` directory.
- `rtk git log --follow --oneline --
  docs/specs/9999-task03-rehearsal/references/2026-07-25-spec-owned-reference-lifecycle.md`
  — exit 0 with no commits, demonstrating that the uncommitted destination has
  no history for `--follow` to traverse.
- `rtk git log --follow --oneline --
  docs/findings/2026-07-25-spec-owned-reference-lifecycle.md` — exit 0 and
  reached pre-adoption commits `397227f` and `be9c42c` at the committed source
  path.
- `rtk rm -rf /tmp/roundfix-task03-conflict.IMxBG3` — exit 0; discarded the
  disposable clone, its scratch branch, and the staged move.
- `rtk test ! -e /tmp/roundfix-task03-conflict.IMxBG3` — the wrapper returned
  exit 0 but printed `sh: -e: command not found`; rejected as cleanup evidence.
- `rtk proxy /usr/bin/test ! -e
  /tmp/roundfix-task03-conflict.IMxBG3` — exit 1 because this host has no
  `/usr/bin/test`; rejected as cleanup evidence.
- `rtk proxy /bin/test ! -e
  /tmp/roundfix-task03-conflict.IMxBG3` — exit 0; confirmed the rehearsal probe
  was fully discarded.
- `rtk git status --short --branch` in the shared worktree — exit 0; only this
  Task file differs, including the Daemon-owned `pending` to `in_progress`
  status change that pre-existed this Agent's edits. No rehearsal artifact or
  scratch branch exists in the shared repository.
- `rtk git branch --list 'ma/spec-0060-task03-rehearsal'` in the shared
  repository — exit 0 with no output; the scratch branch existed only in the
  discarded clone.
- `rtk git diff --check` after this Result update — exit 0.

### Daemon-owned verification

Neither command under `## Verification` was run. The Daemon retains ownership
of declared Verification, Task status, and terminal settlement.
