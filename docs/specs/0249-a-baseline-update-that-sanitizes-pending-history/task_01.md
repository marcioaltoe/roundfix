---
task: task_01
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: completed
type: docs
complexity: medium
---

# Task 01: The glossary, the command references and the Roundfix Skill describe the Pending History

## Overview

Describes what task_02, task_03 and task_04 build, so the glossary, the
command references and the Roundfix Skill ship with the CLI behavior. That
behavior is: the update's history section, `--no-history`, the tag the update
creates, the upgrade notice line, the list-of-maps `unproven` and one-line
refusals. The Task is verifiable on its own through phrase checks, the skill
mirrors and the recorded skill version.

## Requirements

1. MUST add to `CONTEXT.md`, following the `domain-modeling` skill and in the
   glossary's `**Term**:` / definition / `_Avoid_:` shape, beside
   **History Sanitize Command**, the term **Pending History**. It is the units
   of a repository's existing history that the History Sanitize Command would
   plan, which the Managed Refresh includes in its plan and the upgrade notice
   names (ADR-0254). Its definition MUST contain the phrase
   `roundfix baseline update`.
2. MUST revise these existing `CONTEXT.md` entries in place, each keeping its
   other content:
   - **Managed Refresh**: it also plans the Pending History, binds it into the
     Plan Digest and converts it on approval. Its preservation proof still
     covers the instruction carriers (ADR-0254). The entry MUST contain
     `Pending History`.
   - **History Sanitize Command**: `roundfix baseline update` runs the same
     planning for every convertible unit at once. The command stays the path
     for reviewed batches (ADR-0254). The entry MUST contain
     `roundfix baseline update`.
   - **Refused Unit**: its reason prints on one line, and the update lists one
     without blocking the Baseline Plan or changing its state. A unit with
     uncommitted changes is refused (ADR-0254). The entry MUST contain
     `on one line`.
   - **Lenient Legacy Reading**: it also accepts a legacy list-of-maps
     `unproven`, one text line per map (ADR-0254). The entry MUST contain
     `list of maps`.
   - **Sanitize Batch**: the update converts every convertible unit as one
     batch in its own change (ADR-0254). The entry MUST contain
     `every convertible unit`.
   - **History Full Tag**: the update creates it at `HEAD` when absent, never
     moves it, and the operator pushes it (ADR-0254). The entry MUST contain
     `git push origin history-full`.
3. MUST update `docs/user-guide/commands/baseline.md`, in the managed refresh
   paragraphs, to describe the history section: the `History` text lines, the
   JSON `history` object and its statuses, the combined Plan Digest, the apply
   order, the created tag and `git push origin history-full`, Refused Units
   (uncommitted changes and tag coverage included), `blocked`, and
   `--no-history`. It MUST contain `--no-history`, `History refused:` and
   `git push origin history-full`.
4. MUST update `docs/user-guide/commands/history.md` to describe the
   list-of-maps `unproven` under Lenient Legacy Reading and the one-line
   refusal reason. It MUST also say that `roundfix baseline update` plans the
   same units at once and that the History Sanitize Command remains the
   reviewed-batch path. It MUST contain `list of maps` and
   `roundfix baseline update --no-history`. The Operator batch procedure keeps
   its steps.
5. MUST update `docs/user-guide/commands/upgrade.md` with the history notice
   line inside and outside a repository. It MUST contain `pending sanitize`.
6. MUST describe the same behavior for Agents in
   `.agents/skills/roundfix/references/baseline.md`, which MUST contain
   `--no-history`, `Pending History` and `git push origin history-full`, and
   in `.agents/skills/roundfix/references/archive.md` under
   `## History Sanitize Command`, which MUST contain `list of maps` and
   `roundfix baseline update`. Only text about these behaviors changes.
7. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, from the value it holds when this Task
   starts to the next free version at record time. Then run `make skills-sync`
   and record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Those commands are implementation steps, never part of Verification, and
   no version is written by hand.
8. MUST NOT change:
   - any other command reference or glossary entry;
   - the `### QA settlement` section of any skill;
   - any Baseline module or Baseline guide under `docs/agents/`.

## Subtasks

- [ ] Add **Pending History** and revise the six existing terms.
- [ ] Document the history section in the baseline reference, the one-line refusals in the history reference and the notice in the upgrade reference.
- [ ] Describe the behavior in the skill's baseline and archive references.
- [ ] Raise and record the skill version and sync the mirrors.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Pending History** and carries the revised **Managed Refresh**, **History Sanitize Command**, **Refused Unit**, **Lenient Legacy Reading**, **Sanitize Batch** and **History Full Tag**.
- [ ] The three command references and the skill's two references name `--no-history`, the tag push and the one-line refusals.
- [ ] Every skill mirror equals its canonical file, and the raised version is recorded.

## Context

- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/upgrade.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Pending History**:'; f CONTEXT.md '**Managed Refresh**:'; f CONTEXT.md '**History Sanitize Command**:'; f CONTEXT.md '**Refused Unit**:'; f CONTEXT.md '**Lenient Legacy Reading**:'; f CONTEXT.md '**Sanitize Batch**:'; f CONTEXT.md '**History Full Tag**:'; f CONTEXT.md 'on one line'; f CONTEXT.md 'list of maps'; f CONTEXT.md 'every convertible unit'; f CONTEXT.md 'git push origin history-full'; f docs/user-guide/commands/baseline.md '--no-history'; f docs/user-guide/commands/baseline.md 'History refused:'; f docs/user-guide/commands/baseline.md 'git push origin history-full'; f docs/user-guide/commands/history.md 'list of maps'; f docs/user-guide/commands/history.md 'roundfix baseline update --no-history'; f docs/user-guide/commands/upgrade.md 'pending sanitize'; f .agents/skills/roundfix/references/baseline.md '--no-history'; f .agents/skills/roundfix/references/baseline.md 'Pending History'; f .agents/skills/roundfix/references/baseline.md 'git push origin history-full'; f .agents/skills/roundfix/references/archive.md 'list of maps'; f .agents/skills/roundfix/references/archive.md 'roundfix baseline update'; cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md || exit 1; cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0. Before this Task `CONTEXT.md` has no **Pending History**, so the command fails at its first phrase check. After it, every phrase is present, the mirrors match and the raised version is recorded.

## References

- `_prd.md` → Core Feature 7; Goals; User Story 1; User Story 3; User Story 4
- `_techspec.md` → API Contract 5; API Contract 6; API Contract 7; API Contract 8; Build Order 1
- ADR-0254; ADR-0248; ADR-0251

## Result

Implemented the documentation slice for Pending History. `CONTEXT.md` now
defines Pending History and revises the six requested glossary terms. The
baseline, history and upgrade command references describe the history section,
JSON statuses, Plan Digest and apply order, refusals, tag push, opt-out and
upgrade notice. The Roundfix baseline and archive references carry the same
behavior.

Focused implementation evidence:

- `make skills-sync` passed after the reference edits and after the version
  recording step.
- `GOCACHE=/private/tmp/roundfix-task-0249-gocache go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` passed and
  raised both Roundfix Skill front-matter fields from `0.1.53` to `0.1.54`;
  `skills/testdata/owned-skill-versions.json` records the new digest.
- The first recorder attempt hit the shared Go build-cache permission error;
  the task-scoped cache rerun reached the recorder and passed. The declared
  Verification commands were not run in this daemon-assigned turn.
- Canonical and distributed Roundfix Skill files match for `SKILL.md`,
  `references/baseline.md` and `references/archive.md`.
- Required phrase checks passed for the glossary and all three command
  references; `git diff --check` passed.

Acceptance evidence:

1. Glossary: Pending History and the revised Managed Refresh, History Sanitize
   Command, Refused Unit, Lenient Legacy Reading, Sanitize Batch and History
   Full Tag entries are present with the requested behavior phrases.
2. References: baseline and history docs name the History text lines, JSON
   statuses, `--no-history`, one-line refusals, legacy `list of maps`, and
   `git push origin history-full`; the upgrade reference documents notices
   inside and outside repositories with `pending sanitize`.
3. Skill mirrors and version: canonical and distributed files compare equal,
   and the raised `0.1.54` version is recorded by the repository command.

Verification feedback repair:

- The Daemon diagnostic showed `TestEveryOwnedSkillVersionIsRecorded` refused
  the final canonical content because the wording edit after the first record
  left version `0.1.53` with an old digest.
- Re-ran the required recorder with the task-scoped `GOCACHE`; it passed and
  raised the canonical and recorded version to `0.1.54`.
- Re-ran `make skills-sync` after recording. The declared Verification command
  was not rerun in this repair turn.
