---
task: task_01
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: completed
type: docs
complexity: low
---

# Task 01: The Roundfix Skill and the command guides describe the item binary and the migration check

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship
the Roundfix Skill update. This Task describes, in the `deliver` and `setup`
references of the Roundfix Skill and in the `deliver`, `migrate` and
configuration guides, the behavior task_02, task_03 and task_04 implement, as
the TechSpec states it: the `delivery.item_binary` declaration, the item
binary the queue owner builds and runs for an item's steps, the schema rule
and its notice, the parks of a broken declaration, and
`roundfix migrate --check`.

## Requirements

1. MUST describe in `.agents/skills/roundfix/references/deliver.md` and
   `docs/user-guide/commands/deliver.md`: that a repository may declare
   `delivery.item_binary` with `build` and `path`; that the owner reads the
   declaration it loaded at start; that before each `implement`, `archive`
   and `review` step it builds the binary in the item worktree, asks it
   `migrate --check`, and runs the step with it on exit `0`; that any other
   exit runs the step with the owner's binary; that a failed build, a path Git
   does not ignore, or a binary that cannot start parks the item
   `delivery-error`; where the build log lives; and the two console-log lines
   of API Contracts 2 and 3 with their words. They MUST say that a repository
   without the declaration, which includes every repository that installs
   Roundfix from npm, is delivered as before.
2. MUST describe in `docs/user-guide/configuration.md`, beside
   `delivery.derived_paths`, the `delivery.item_binary` key, its two fields,
   that Project Config replaces User Config, an example declaring
   `make build` and `bin/roundfix`, and the three configuration errors of API
   Contract 4 verbatim.
3. MUST describe in `docs/user-guide/commands/migrate.md` and
   `.agents/skills/roundfix/references/setup.md` the synopsis
   `roundfix migrate [--check]` and each outcome of "The migration check",
   including the wording of Surface Transcripts 1 to 3, the exit codes, and
   that the check never migrates or writes the database.
4. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
5. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, `commands.md`, `.roundfixrc.yml` or
   `CONTEXT.md`.

## Subtasks

- [ ] Describe the declaration, the item binary, the schema rule and the parks in the `deliver` reference and guide.
- [ ] Describe `delivery.item_binary` in the configuration guide.
- [ ] Describe `migrate --check` in the `migrate` guide and the `setup` reference.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The `deliver` reference and guide carry `delivery.item_binary`,
      `runs the item binary` and `runs the owner's binary`.
- [ ] The configuration guide carries `delivery.item_binary requires build
      and path`; the `migrate` guide and the `setup` reference carry
      `roundfix migrate --check` and `the version this binary supports`.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/commands/migrate.md`
- interface: `docs/user-guide/configuration.md`
- instruction: `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "delivery.item_binary" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "delivery.item_binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "runs the item binary" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "runs the item binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "runs the owner's binary" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "runs the owner's binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "migrate --check" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "migrate --check" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "delivery.item_binary" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "delivery.item_binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "runs the item binary" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "runs the item binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "runs the owner's binary" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "runs the owner's binary" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "migrate --check" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "migrate --check" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "delivery.item_binary requires build and path" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "delivery.item_binary requires build and path" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "delivery.item_binary has unsafe path" || { printf 'missing phrase in %s: %s\n' docs/user-guide/configuration.md "delivery.item_binary has unsafe path" >&2; exit 1; }` — expected: exit 0; before this Task none of these phrases is in these files, so the command fails.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- "roundfix migrate --check" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/setup.md "roundfix migrate --check" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- "the version this binary supports" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/setup.md "the version this binary supports" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/migrate.md | grep -qF -- "roundfix migrate --check" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/migrate.md "roundfix migrate --check" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/migrate.md | grep -qF -- "the version this binary supports" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/migrate.md "the version this binary supports" >&2; exit 1; }` — expected: exit 0; before this Task neither guide nor reference names the check, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/deliver.md | grep -qF -- "delivery.item_binary" && tr -s '[:space:]' ' ' < skills/roundfix/references/setup.md | grep -qF -- "roundfix migrate --check" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && cmp .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirrors do not carry the new phrases, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 8; User Stories 1-5
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Build Order 1
- ADR-0187; ADR-0189; ADR-0225

## Result

Implemented the documentation slice for the item binary and migration check.
The delivery reference and guide now describe `delivery.item_binary`, the
owner-loaded declaration, per-step build and probe, build-log location,
item-binary and owner-binary console lines, `delivery-error` parks, and the
unchanged npm-adopter path. The configuration guide documents `build`, `path`,
Project Config precedence, the `make build` / `bin/roundfix` example, and all
three declaration errors. The migrate guide and setup reference document
`roundfix migrate [--check]`, all migration-check outcomes, exit codes,
transcripts, and the no-write guarantee. The Roundfix Skill was raised from
0.1.22 to 0.1.23 and its embedded mirrors were regenerated.

Focused checks and outcomes:

- `make skills-sync`: passed; embedded Roundfix Skill files regenerated.
- `make baseline-digests`: passed; no unrelated derived changes remained.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  the first attempt hit the host Go cache permission error; the same command
  passed with task-scoped `GOCACHE=/tmp/roundfix-0221-task01-gocache` and
  recorded Roundfix 0.1.23 with its new digest.
- `cmp` checks for the three canonical Roundfix Skill files and their embedded
  mirrors: passed.
- `git diff --check`: passed.
- Phrase and transcript searches confirmed the required delivery, configuration,
  migration, version-support, build-log, and park text is present in the named
  files.

Acceptance evidence:

- Delivery reference and guide contain `delivery.item_binary`, `runs the item
  binary`, `runs the owner's binary`, `migrate --check`, the exact API Contract
  lines, build-log path, fallback behavior, parks, and npm-adopter behavior.
- Configuration guide contains `delivery.item_binary requires build and path`,
  the unsafe-path and unknown-key errors, both fields, precedence, and the
  requested example.
- Migrate guide and setup reference contain `roundfix migrate --check`,
  `the version this binary supports`, Surface Transcripts 1 to 3, exit codes,
  and the no-migration/no-write rule.
- Canonical and embedded Skill files compare equal, and the owned-skill version
  record contains Roundfix 0.1.23 and its content digest.

## Carry-forward provenance

- Source Run: `run_20261003T234640Z_3e854d86e9475f64`
- Source commit: `39e164af1ec5cec8659c3fc869c994734a7ca626`
