---
task: task_03
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: completed
type: docs
complexity: low
---

# Task 03: The Roundfix Skill and the deliver guide describe both rules

## Overview

The repository's skill-sync rule requires a change to CLI behavior to ship the
Roundfix Skill update. The `deliver` guide and the skill's `deliver` reference
say today that only the QA environment park may use the Run start head and
that a partial blocked only by the pre-PR row keeps `run-unresolved`; both
sentences became wrong with ADR-0229. This Task rewrites them from the
TechSpec, raises the skill's version and syncs the mirrors.

## Requirements

1. MUST state in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md` the rule of API Contract 2:
   an unresolved Run parks `qa-environment-partial` when its newest Run-Branch
   QA Report is `partial`, has no finding-blocked row and has at least one
   environment-blocked row, a row blocked only because no Pull Request is open
   included; otherwise it keeps `run-unresolved`. Each file MUST carry the
   phrase "has at least one environment-blocked row".
2. MUST state in both files the rule of API Contract 1: when no candidate is
   recorded, only an item the operator archived with the QA Archive Override
   may use the Run start head, whatever its park, and the retry records the
   descended head and resumes at `reviewing`. Each file MUST carry the phrase
   "whatever its park", and the re-entry table row for the operator archive
   MUST no longer limit it to `qa-environment-partial`.
3. MUST remove from both files the sentences "Only the QA environment park may
   use the Run start head", "blocked only by the pre-PR row keeps" and "rows
   beyond the pre-PR Pull Request row", and MUST keep the
   `corrective-spec-required` and `pull-request-conflict` restrictions.
4. MUST raise the Roundfix Skill's version by one patch level from the value
   on the tree the Task starts from, in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`, run `make skills-sync` so the mirrors
   equal their canonical files, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
5. MUST NOT edit the `### QA settlement` section of any skill, any other
   command reference or guide, `commands.md`, `.roundfixrc.yml` or
   `CONTEXT.md`.

## Subtasks

- [ ] Rewrite the environment-partial park rule in the guide and the reference.
- [ ] Rewrite the Run start anchor rule and the re-entry table row.
- [ ] Raise the version, sync the mirrors and record the version.

## Acceptance Criteria

- [ ] The guide and the reference carry "has at least one
      environment-blocked row" and "whatever its park", and neither carries
      the three removed sentences.
- [ ] Each mirror equals its canonical file, and the raised version is
      recorded.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do t="$(tr -s '[:space:]' ' ' < "$f")" || exit 1; for p in "has at least one environment-blocked row" "whatever its park"; do printf '%s\n' "$t" | grep -qF -- "$p" || { printf 'missing phrase in %s: %s\n' "$f" "$p" >&2; exit 1; }; done; for p in "Only the QA environment park may use the Run start head" "blocked only by the pre-PR row keeps" "rows beyond the pre-PR Pull Request row"; do if printf '%s\n' "$t" | grep -qF -- "$p"; then printf 'stale phrase in %s: %s\n' "$f" "$p" >&2; exit 1; fi; done; done` — expected: exit 0; before this Task neither file carries the new phrases and both carry stale ones, so the command fails.
- `tr -s '[:space:]' ' ' < skills/roundfix/references/deliver.md | grep -qF -- "whatever its park" && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the mirror does not carry the new phrase, so the command fails; after it the mirrors equal their canonical files and the raised version and its content digest are recorded.

## References

- `_prd.md` → Core Feature 3
- `_techspec.md` → Vocabulary Contract; API Contract 1; API Contract 2; Build Order 3
- ADR-0187; ADR-0189; ADR-0229

## Result

The deliver guide and Roundfix Skill reference now describe the environment
partial predicate as a zero-finding partial with at least one
environment-blocked row, including rows waiting for an open Pull Request. They
also describe the QA Archive Override as the authorization for using the Run
start head when no candidate is recorded, whatever its park, and generalize
the operator archive re-entry row accordingly. The `corrective-spec-required`
and `pull-request-conflict` restrictions remain documented.

The canonical Roundfix Skill version was raised from `0.1.24` to `0.1.25` in
both front-matter fields. `make skills-sync` regenerated the mirrors, and the
owned-version record now includes the `0.1.25` content digest.

Focused checks:

- `make skills-sync` — passed.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed (`ok roundfix/skills`); this recorded the new version.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` and the matching `deliver.md` comparison — passed.
- Required phrase searches found `has at least one environment-blocked row`
  and `whatever its park` in both documents. Searches for each of the three
  removed phrases returned no matches.
- `git diff --check` — passed.

Evidence for acceptance criteria:

- Both guide surfaces contain the new environment and operator-archive rules,
  and the stale sentences are absent.
- Both Skill mirrors compare equal, and the raised version is recorded in
  `skills/testdata/owned-skill-versions.json`.

## Carry-forward provenance

- Source Run: `run_20261004T162613Z_0bb1b015eb38156b`
- Source commit: `02bc1cff3195d25d1078db696620a31a354174ad`
