---
task: task_03
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: completed
type: docs
complexity: low
---

# Task 03: The skill, the guide and the glossary describe Relocation Citations

## Overview

The Roundfix skill and the context-driven development guide describe what
`baseline update` and `baseline plan` report, and the glossary defines the
History Relocation. After Task 02, a plan with History Relocations also reports
its Relocation Citations. This Task documents the warnings, their three codes,
the digest binding and the unchanged apply, in the skill, the guide and a new
glossary term. Adopters and Agents read these documents to decide how to review
a plan, so they must say what the plan reports and what it never does: rewrite
a citation, or change what apply writes.

## Requirements

1. MUST add a paragraph to the Baseline update section of
   `.agents/skills/roundfix/SKILL.md`, then run `make skills-sync` so that
   `skills/roundfix/SKILL.md` matches. The paragraph MUST:
   - contain the phrase
     `report each tracked file whose citations its History Relocations would break`;
   - name the codes `baseline.history.citation`,
     `baseline.history.citation.omitted` and
     `baseline.history.citation.unscanned`;
   - contain the phrase `The Plan Digest covers these warnings`;
   - state that planning never rewrites a citation and that apply writes the
     same files.
2. MUST add the same content, with the same two phrases and the three codes, to
   the managed-refresh update section of
   `docs/user-guide/context-driven-development.md`.
3. MUST add a **Relocation Citation** entry to `CONTEXT.md`, right after the
   **History Relocation** entry. The entry MUST:
   - define it as a citation in a tracked file that resolves before a Baseline
     Plan's History Relocations and would not resolve after them;
   - contain the phrase `reported as a warning the Plan Digest binds`;
   - list `_Avoid_: Broken link, dangling reference, link check`.
4. MUST NOT edit any other skill, `skills/_ownership.yml` or any governed file
   outside the two bounded Roundfix skill paths.

## Subtasks

- [ ] Document the warnings in the canonical Roundfix skill and regenerate its
      mirror.
- [ ] Document them in the context-driven development guide.
- [ ] Add the glossary term.

## Acceptance Criteria

- [ ] Both Roundfix skill copies and the guide contain both phrases and all
      three codes, and `make skills-sync-check` exits `0`.
- [ ] `CONTEXT.md` defines **Relocation Citation** with the required phrase.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`

## Verification

- `for file in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md docs/user-guide/context-driven-development.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "report each tracked file whose citations its History Relocations would break" || { printf 'missing phrase in %s: %s\n' "$file" "report each tracked file whose citations its History Relocations would break" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "The Plan Digest covers these warnings" || { printf 'missing phrase in %s: %s\n' "$file" "The Plan Digest covers these warnings" >&2; exit 1; }; grep -qF -- "baseline.history.citation.omitted" "$file" || { printf 'missing code in %s\n' "$file" >&2; exit 1; }; grep -qF -- "baseline.history.citation.unscanned" "$file" || { printf 'missing code in %s\n' "$file" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0. Before this Task the phrases are absent, so the command fails.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Relocation Citation**:" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "**Relocation Citation**:" >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "reported as a warning the Plan Digest binds" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "reported as a warning the Plan Digest binds" >&2; exit 1; }` — expected: exit 0. Before this Task the term is absent.

## References

- [_prd.md](_prd.md) — User Stories 1–2; Core Features 3–4; Decisions
- [_techspec.md](_techspec.md) — Data Models; API Contract 1; Integration
  Points; Build Order 3
- ADR-0173

## Result

Implemented the documentation slice for Relocation Citations:

- Added the warning behavior, all three warning codes, Plan Digest binding, and
  unchanged planning/apply behavior to the canonical Roundfix skill and the
  managed-refresh guide section.
- Added the `Relocation Citation` glossary entry immediately after `History
  Relocation`, including its resolution definition, digest-warning phrase, and
  avoidance terms.

Focused checks:

- `make skills-sync` completed successfully; `cmp
  .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited `0`.
- `make baseline-digests` completed successfully with `changed=false`.
- `git diff --check` exited `0`.
- Focused `rg` checks found both required phrases and all three warning codes in
  the canonical skill, shipped skill, and guide. Focused `rg` checks found the
  glossary heading, digest-warning phrase, and required `_Avoid_` list in
  `CONTEXT.md`.
- The task's declared Verification commands were not run; the Daemon retains
  that verification responsibility.

## Carry-forward provenance

- Source Run: `run_20260929T222541Z_a2fa4eeca2973de1`
- Source commit: `de92c3203244003c2d0301ba3a71d0db6a741711`
