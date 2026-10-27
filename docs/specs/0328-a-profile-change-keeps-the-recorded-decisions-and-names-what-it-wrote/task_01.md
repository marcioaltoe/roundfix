---
task: task_01
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: pending
type: docs
complexity: medium
---

# Task 01: The guides, skills and glossary describe the profile change, the Root Backup and the Claude Skill Links

## Overview

The user guides, the Roundfix skill's `baseline` reference and the
`setup-context-driven` skill describe what tasks 02 to 04 changed: a profile
change through `baseline plan` carries the recorded decisions and refuses
once, `baseline update` names that command, an apply names each Root Backup,
reconcile names the Setup Snapshot revision, and a restore keeps the Claude
Skill Links. `CONTEXT.md` gains the three terms, and the owned-skill versions
and the Behavior Surface record follow through their record commands. This
closes the documentation of frictions (a) to (e) of the adopted Backlog Entry
(Quaestor, 2026-10-26).

## Requirements

1. MUST update `docs/user-guide/context-driven-development.md` in its
   "Update, profile change, and rejected plans", "Automation", "First
   adoption" and "Repository Skill Set restoration" sections, keeping every
   heading: a profile change from automation runs `baseline plan --profile
   <id>` with `preservation.mode`, which carries each Carried Decision and
   names every missing answer in one refusal; `baseline update` refuses
   `--decision profile=<id>` naming that command; each Root Backup an apply
   writes is named by a `baseline.preservation.backup.written` warning and is
   the repository's to commit; a restore writes a Claude Skill Link where
   `.claude/skills` keeps them, and a restore of an installed skill repairs a
   missing one. The page MUST contain the phrases `Carried Decision`, `Root
   Backup` and `Claude Skill Link`.
2. MUST update `docs/user-guide/commands/baseline.md` and the canonical
   `.agents/skills/roundfix/references/baseline.md` with the outputs of API
   Contracts 1 to 6: each MUST contain `Carried decisions:`,
   `baseline.preservation.backup.written`, `create-link` and `never moves an
   installed skill`.
3. MUST update the canonical `.agents/skills/setup-context-driven/SKILL.md`
   section "Non-interactive planning for adoption or profile change" and its
   restore guidance so it contains `Carried Decision` and `Claude Skill Link`,
   under existing headings or a new heading of its own, never inside a
   `### QA settlement` section.
4. MUST copy the canonical skills to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/setup-context-driven/SKILL.md` and their mirrors with the
   `-record-skill-versions` step (`go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`), recording `skills/testdata/owned-skill-versions.json`.
5. MUST add **Carried Decision**, **Root Backup** and **Claude Skill Link**
   to `CONTEXT.md` through `domain-modeling`, each with its `_Avoid_` line
   and a link to ADR-0341, never to a Spec.
6. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Update the two user guide pages.
- [ ] Update the canonical Roundfix reference and the `setup-context-driven` skill, sync the mirrors and record the versions.
- [ ] Add the three glossary terms.
- [ ] Re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Each guide and skill names the profile-change path, the Root Backup warning, the reconcile explanation and the Claude Skill Links.
- [ ] `CONTEXT.md` defines the three terms.
- [ ] The owned-skill version and Behavior Surface contracts pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0341-a-profile-change-carries-the-recorded-decisions-and-a-restore-keeps-the-claude-skill-links.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned-skill version record and the Behavior Surface record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for phrase in 'Carried decisions:' 'baseline.preservation.backup.written' 'create-link' 'never moves an installed skill'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "$phrase" || { printf 'missing phrase in the baseline guide: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- "$phrase" || { printf 'missing phrase in the baseline reference: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in 'Carried Decision' 'Claude Skill Link'; do tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "$phrase" || { printf 'missing phrase in the context-driven guide: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/setup-context-driven/SKILL.md | grep -qF -- "$phrase" || { printf 'missing phrase in the setup skill: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'Root Backup' || { echo 'guide lacks Root Backup' >&2; exit 1; }; for term in '**Carried Decision**:' '**Root Backup**:' '**Claude Skill Link**:'; do grep -qF -- "$term" CONTEXT.md || { printf 'CONTEXT.md lacks %s\n' "$term" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/setup-context-driven/SKILL.md skills/setup-context-driven/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the guides or skills names `Carried decisions:` and `CONTEXT.md` lacks the terms, so the first phrase check fails. Afterwards every phrase and term is present, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-6; User Experience; Non-Goals / Out of Scope; Success Metric 6; Glossary
- `_techspec.md` → API Contracts 1-6; Glossary; Build Order 4
- ADR-0341
