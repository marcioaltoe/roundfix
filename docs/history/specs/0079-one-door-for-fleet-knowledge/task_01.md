---
task: task_01
spec: 0079-one-door-for-fleet-knowledge
status: completed
type: docs
complexity: low
---

# Task 01: Name the door's vocabulary in the glossary

## Overview

The vocabulary exists before any rule uses it. Three terms enter the project
glossary beside Finding and Backlog Entry, so every later clause, check, and
skill instruction speaks in defined words instead of coining synonyms.

## Requirements

1. MUST add an `Inbox Entry` glossary entry: a fleeting capture born
   committed in the Secondbrain's inbox under its destination's namespace,
   carrying origin, destination, an advisory type hint, and its capture
   mode; pending until triage resolves it into exactly one Finding, one
   Backlog Entry, or one recorded discard.
2. MUST add a `Rollup` entry: a Finding of rollup kind that consolidates
   related Findings as declared members, supersedes them, and is the license
   for their archival.
3. MUST add a `Triage` entry: the destination-project act that converts one
   pending Inbox Entry into its contract-true artifact, committed by the
   destination alone.
4. MUST give each entry the names to avoid, following the glossary's
   existing `_Avoid_` convention, and MUST keep every existing entry
   byte-identical.
5. MUST NOT add clauses, checks, or skill text — this slice is vocabulary
   only.

## Subtasks

- [ ] Write the three glossary entries with their avoid-lists.
- [ ] Cross-check the wording against the PRD's Core Features 1–4.

## Acceptance Criteria

- [ ] The glossary defines Inbox Entry, Rollup, and Triage in the shapes
      above.
- [ ] Each entry carries an `_Avoid_` line.
- [ ] No other glossary entry changed.

## Verification

- `grep -q "^\*\*Inbox Entry\*\*" CONTEXT.md && grep -q "^\*\*Rollup\*\*" CONTEXT.md && grep -q "^\*\*Triage\*\*" CONTEXT.md`
  — expected: exit 0; the three terms exist.
- `output="$(grep -A3 "^\*\*Rollup\*\*" CONTEXT.md)"; printf '%s' "$output" | grep -q "_Avoid_"`
  — expected: exit 0; the new entries carry avoid-lists.
- `git diff HEAD -- CONTEXT.md | grep "^-" | grep -v "^---" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the change is purely additive — no existing line
  removed.

## References

- `_prd.md` → Core Feature 11; Goals.
- `_techspec.md` → System Architecture; Build Order 1.
- ADR-0095.

## Result

Implemented the three glossary terms beside `Backlog Entry` in `CONTEXT.md`.
The definitions retain the PRD Core Features 1–4 boundaries: capture is born
committed in the Secondbrain, triage belongs to the destination project, and a
rollup supersedes declared Finding members and licenses their archival. No
clauses, checks, or skill text changed.

Focused checks (the Daemon-owned `## Verification` commands were not run):

- `rtk awk '/^\*\*(Inbox Entry|Rollup|Triage)\*\*:/ { term=$0; found[term]=1; next } /^_Avoid_:/ && term != "" { avoid[term]=1; term="" } END { terms[1]="**Inbox Entry**:"; terms[2]="**Rollup**:"; terms[3]="**Triage**:"; for (i=1; i<=3; i++) { printf "%s definition=%s avoid=%s\n", terms[i], found[terms[i]] ? "yes" : "no", avoid[terms[i]] ? "yes" : "no"; if (!found[terms[i]] || !avoid[terms[i]]) bad=1 } exit bad }' CONTEXT.md` — exit 0; each requested term reported `definition=yes avoid=yes`.
- `rtk git diff --numstat -- CONTEXT.md` — exit 0; reported `12 0 CONTEXT.md`, proving the glossary edit is purely additive with no existing line removed.

Acceptance evidence:

- The `Inbox Entry`, `Rollup`, and `Triage` blocks state every shape required
  by Acceptance Criterion 1 and agree with the PRD, TechSpec, and ADR-0095.
- The focused `awk` check associates one `_Avoid_` line with each new block,
  satisfying Acceptance Criterion 2.
- The zero-deletion Git numstat proves every pre-existing glossary entry stayed
  byte-identical, satisfying Acceptance Criterion 3.
