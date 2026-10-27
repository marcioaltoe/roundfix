---
task: task_01
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: pending
type: docs
complexity: medium
---

# Task 01: The guides, skills and glossary describe the lock rewrite and the gate-less settlement

## Overview

The user guides, the Roundfix skill and the `write-tasks` skill describe what
tasks 02 to 04 implement: the bundle install rewrites the lock hash of each
installed owned skill, every Task settles on the repository Verification
with or without a QA gate, and a Run without a QA gate runs the Integrated
Repository Verification before it ends Clean. `CONTEXT.md` gains one term and
revises two, and the owned-skill versions and the Behavior Surface record
follow through their record commands. This Task runs first, so each code
Task names the guide of its command through this dependency. It closes the
documentation of the two adopted Backlog Entries (Quaestor, 2026-10-26 and
2026-10-27).

## Requirements

1. MUST update `docs/user-guide/commands/implement.md`, section "Settlement
   Checks", keeping its heading: the repository Verification at settlement
   runs for every non-QA Task with or without a QA gate Task; the two
   in-process checks run only in a graph with one; a Run without a QA gate
   runs the Integrated Repository Verification under the conditions of
   `_techspec.md` → API Contract 5 and ends `Unresolved` with the line of
   API Contract 6 when it fails. The section MUST contain the phrases `with
   or without a QA gate` and `Integrated Repository Verification`.
2. MUST update `docs/user-guide/configuration.md` where it describes
   `verification.repository_at_settlement` (the paragraph after the empty
   value table and its row in the key table): the key now applies to every
   graph and also turns off the Integrated Repository Verification, and a
   repository without the built-in `make verify` target configures
   `defaults.verification` or turns the key off. Both
   places MUST contain `Integrated Repository Verification`, and neither may
   keep `of a graph with a QA gate Task`.
3. MUST update `docs/user-guide/commands/baseline.md` (managed refresh
   skills paragraphs), `docs/user-guide/commands/skills.md` and the
   canonical `.agents/skills/roundfix/references/baseline.md` (after
   "`--no-skills` skips the Repository Skill Set refresh."): after the bundle
   install, an existing `skills-lock.json` entry that names an owned skill
   gets the `computedHash` of the installed folder, every other field is
   kept, no entry or lock is created, and the names are reported. The
   baseline guide and reference MUST contain `skills.lockUpdated` and
   `Skills lock updated:`; the skills guide MUST contain `Updated
   skills-lock.json computedHash:`.
4. MUST update the canonical `.agents/skills/roundfix/SKILL.md` section
   "### Settlement Checks" (never the `### QA settlement` section) and the
   canonical `.agents/skills/roundfix/references/implement.md` (after the
   repository precondition paragraph) the same way as requirement 1, each
   containing `with or without a QA gate` and `Integrated Repository
   Verification`; the canonical `.agents/skills/roundfix/references/setup.md`
   (beside `roundfix skills check`) the same way as the skills guide of
   requirement 3, containing `Updated skills-lock.json computedHash:`; and
   the canonical
   `.agents/skills/write-tasks/SKILL.md` bullet "A gated graph leaves the
   repository Verification green." so it reads that every Task leaves the
   repository Verification green, containing `Every Task leaves the
   repository Verification green` and no longer `A gated graph leaves`.
5. MUST copy the canonical skills to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/write-tasks/SKILL.md` and their mirrors with the
   `-record-skill-versions` step (`go test ./skills -run
   '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
   -count=1`), recording `skills/testdata/owned-skill-versions.json`.
6. MUST, through `domain-modeling`, add **Integrated Repository
   Verification** to `CONTEXT.md` with its `_Avoid_` line and a link to
   ADR-0342, revise **Settlement Check** so its definition says the
   repository Verification at settlement runs with or without a QA gate and
   the two in-process checks only with one, and revise **Clean** so a Spec
   Run without a QA gate also needs a passing Integrated Repository
   Verification when one runs. No entry references a Spec.
7. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Update the four user guide pages.
- [ ] Update the canonical Roundfix and `write-tasks` skills, sync the mirrors and record the versions.
- [ ] Add and revise the glossary terms.
- [ ] Re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Every guide and skill that described the gate-scoped settlement now describes the settlement in every graph and the Integrated Repository Verification.
- [ ] The baseline and skills guides and the baseline reference describe the lock rewrite and its output.
- [ ] `CONTEXT.md` defines the new term and the two revised ones; the owned-skill version and Behavior Surface contracts pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0342-every-task-settles-on-the-repository-verification-and-a-run-without-a-qa-gate-checks-its-integrated-tree.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- interface: `docs/user-guide/commands/implement.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned-skill version record and the Behavior Surface record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/implement.md | grep -qF -- 'with or without a QA gate' || { echo 'implement guide lacks: with or without a QA gate' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/implement.md | grep -qF -- 'Integrated Repository Verification' || { echo 'implement guide lacks: Integrated Repository Verification' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- 'with or without a QA gate' || { echo 'roundfix skill lacks: with or without a QA gate' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- 'Integrated Repository Verification' || { echo 'roundfix skill lacks: Integrated Repository Verification' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/implement.md | grep -qF -- 'with or without a QA gate' || { echo 'implement reference lacks: with or without a QA gate' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/implement.md | grep -qF -- 'Integrated Repository Verification' || { echo 'implement reference lacks: Integrated Repository Verification' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'Integrated Repository Verification' || { echo 'configuration guide lacks the integrated check' >&2; exit 1; }; if tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'of a graph with a QA gate Task'; then echo 'configuration guide keeps the gate scope' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'skills.lockUpdated' || { echo 'baseline guide lacks: skills.lockUpdated' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'Skills lock updated:' || { echo 'baseline guide lacks: Skills lock updated:' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'skills.lockUpdated' || { echo 'baseline reference lacks: skills.lockUpdated' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'Skills lock updated:' || { echo 'baseline reference lacks: Skills lock updated:' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Updated skills-lock.json computedHash:' || { echo 'skills guide lacks the lock line' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- 'Updated skills-lock.json computedHash:' || { echo 'setup reference lacks the lock line' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- 'Every Task leaves the repository Verification green' || { echo 'write-tasks lacks the every-Task rule' >&2; exit 1; }; if grep -qF -- 'A gated graph leaves' .agents/skills/write-tasks/SKILL.md; then echo 'write-tasks keeps the gated-only rule' >&2; exit 1; fi; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Integrated Repository Verification**:' || { echo 'CONTEXT.md lacks Integrated Repository Verification' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Settlement Check**:' || { echo 'CONTEXT.md lacks Settlement Check' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Clean**:' || { echo 'CONTEXT.md lacks Clean' >&2; exit 1; }; awk '/^[*][*]Settlement Check[*][*]:/ { getline; if (index($0, "with or without a QA gate") > 0) found=1 } END { exit(found ? 0 : 1) }' CONTEXT.md || { echo 'Settlement Check is not revised' >&2; exit 1; }; awk '/^[*][*]Clean[*][*]:/ { getline; if (index($0, "Integrated Repository Verification") > 0) found=1 } END { exit(found ? 0 : 1) }' CONTEXT.md || { echo 'Clean is not revised' >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && cmp .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md && cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide names `with or without a QA gate`, so the first phrase check fails. Afterwards every phrase and term is present, the gate-only wording is gone, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-3; User Experience; Non-Goals / Out of Scope; Success Metric 7; Glossary
- `_techspec.md` → API Contracts 2-6; Glossary; Build Order 1
- ADR-0342; ADR-0189
