---
task: task_01
spec: 0060-spec-owned-reference-lifecycle
status: completed
type: docs
complexity: high
---

# Task 01: Own the adoption lifecycle in the authorial Skills

## Overview

Promotion currently links a source document without transferring it, so a Spec
can archive while its evidence stays behind in another lifecycle with links that
no longer resolve. Make adoption an explicit transition owned by PRD creation,
validated by the archive gate, and referenced consistently by every authorial
Skill in between.

All five authorized Skill pairs move together on purpose: each edit rewrites the
same seven derived digest artifacts, so splitting them would make the serialized
cherry-pick integration conflict on files neither Task authored.

## Requirements

1. MUST add an adoption step to `write-prd`, between recording decisions and
   writing the Spec folder, covering: inventory of relied-upon sources;
   classification of a raw inbox note as source material versus field evidence
   that becomes a finding first; the ownership claim against both the active and
   archived Spec trees; the finding's `status: done` flip and Spec link recorded
   **before** the move; `git mv` into `docs/specs/<slug>/references/`; the index
   row; and the repository-wide link rewrite.
2. MUST specify one move, never a copy and never a stub, with the basename and
   the bytes preserved and the observations inside never rewritten.
3. MUST specify the index at `docs/specs/<slug>/references/_index.md` as a
   Markdown table whose columns are source (the pre-adoption path, never
   updated), type, owner (the four-digit Spec number), adopted date, and path
   (relative to the index).
4. MUST make `write-prd` fail, and not report completion, while an adopted
   source remains at its original path or a link to it is unresolved.
5. MUST add a third `archive-spec` precondition — the Spec is self-contained —
   with named commands, verified with fresh evidence like the two it already
   has, whose failure names the offending source or link and the adoption step
   that fixes it.
6. MUST make that precondition distinguish a link into `docs/_inbox/` or
   `docs/findings/` from prose that merely names those trees.
7. MUST state that `qa_override: true` does not override self-containment,
   because the override speaks to verification and self-containment is a
   property of the artifact.
8. MUST apply the archive precondition only to a Spec carrying a
   `references/_index.md`, so historical archives stay green — new promotions
   only.
9. MUST add the referencing convention to `write-idea`, `write-techspec`, and
   `write-tasks`: a downstream artifact links an adopted source at its
   post-adoption path, never at its pre-adoption path.
10. MUST specify that exactly one Spec — the first to commit to implementation —
    primarily owns a shared source, and that a secondary Spec links the owner's
    copy and adopts nothing.
11. MUST confine protected-tooling edits to exactly the ten Skill files the PRD
    authorizes, run `make skills-sync` to propagate to the embedded pairs, and
    obtain every derived pin from `make baseline-digests`.

## Subtasks

- [ ] Write the adoption step and the index format into `write-prd`.
- [ ] Write the self-containment precondition into `archive-spec`.
- [ ] Write the referencing convention into the three remaining Skills.
- [ ] Run `make skills-sync`, then `make baseline-digests`, committing only
      authorized paths.

## Acceptance Criteria

- [ ] `write-prd` describes all eight adoption steps in order, with the status
      flip preceding the move.
- [ ] The index format specifies a never-updated provenance column and a
      relative current path.
- [ ] `write-prd` refuses to report completion while a source sits at its
      original path.
- [ ] `archive-spec` lists three preconditions, each with a named command.
- [ ] `archive-spec` states that `qa_override` does not reach the new check and
      that the check applies only to Specs carrying an index.
- [ ] The three downstream Skills state the post-adoption referencing rule.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.
- [ ] Every changed derived pin came from `make baseline-digests`, and no file
      outside the authorized list changed.

## Context

- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`

## Verification

- `make skills-sync-check` — expected: no drift.
- `roundfix skills check` — expected: pass.
- `make baseline-digests` — expected: `ok: true`, only authorized paths changed.
- `make verify` — expected: exit 0.
- `git status --porcelain` — expected: only the ten Skill files and the seven
  authorized derived artifacts.

## References

`_prd.md` → Goals 1–3, Stories 1–4, Features 1–7; `_techspec.md` → Build Order
1–3, The adoption transition, The reference index, The archive precondition;
ADR-0083.

## Result

### Implementation

- `write-prd` now owns an eight-step adoption transition between recording
  decisions and writing the Spec: inventory, classification, ownership claim,
  finding status/link update, one history-preserving move, index entry,
  repository-wide link rewrite, and the authoring gate.
- `archive-spec` now verifies three preconditions with fresh named command
  evidence. Its indexed-Spec-only self-containment check validates current
  paths, rejects surviving source paths, and matches Markdown links into the
  inbox or findings trees without matching prose that only names those trees.
- `write-idea`, `write-techspec`, and `write-tasks` now link adopted sources at
  the primary owner's post-adoption path and prohibit secondary adoption.
- `rtk make skills-sync` propagated the five canonical edits to their five
  embedded counterparts. No other canonical or embedded Skill changed.

### Focused checks

- `rtk make skills-sync` — exit 0.
- `rtk cmp -s .agents/skills/<name>/SKILL.md skills/<name>/SKILL.md` for
  `write-prd`, `archive-spec`, `write-idea`, `write-techspec`, and
  `write-tasks` — all five exited 0.
- `rtk git diff --check` — exit 0.
- Focused `rtk rg -n` contract inspections — exit 0 for the ordered PRD
  transition, the three archive preconditions and override boundary, and all
  three downstream reference rules.
- Archive link-pattern probes — prose naming `docs/_inbox/` and
  `docs/findings/` produced the expected no-match exit 1; an inline findings
  link and a reference-style inbox link each matched with exit 0.
- `rtk git diff --name-only` after synchronization — only the ten authorized
  Skill files and this Task file differ. The Task file's `pending` to
  `in_progress` frontmatter change pre-existed this Agent turn and remains
  Daemon-owned.

### Acceptance evidence

1. `write-prd` contains exactly eight numbered adoption steps in order;
   `Flip then link` is step 4 and `Move` is step 5.
2. Its fixed Markdown index columns are `source`, `type`, `owner`,
   `adopted date`, and `path`; the prose makes `source` never-updated
   provenance and `path` relative to `_index.md`.
3. Adoption step 8 says to fail and not report completion while a source
   remains at its original path or a rewritten link is unresolved.
4. `archive-spec` lists exactly three preconditions and gives fresh command
   evidence for Task status, the newest QA verdict, and indexed-reference
   self-containment.
5. The archive contract says `qa_override: true` reaches only QA evidence,
   never self-containment, and applies the new check only when
   `references/_index.md` exists.
6. Each of `write-idea`, `write-techspec`, and `write-tasks` states the
   post-adoption-path rule and the single-primary-owner rule.
7. Canonical/embedded byte equality is proven by the five focused `cmp`
   checks. `make skills-sync-check` and `roundfix skills check` were not run;
   they are reserved for Daemon Verification.
8. No derived pin has been edited or hand-transcribed. `make baseline-digests`
   was not run because it is a declared Daemon Verification command; that
   Verification step must generate the seven authorized deterministic
   artifacts before the final changed-path check.

### Daemon-owned verification

No command from `## Verification` was run in this Agent turn. The Daemon owns
all five commands, the derived digest regeneration, Task status, and terminal
settlement.
