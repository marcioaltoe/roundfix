---
task: task_02
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: pending
type: docs
complexity: low
---

# Task 02: The spec guide, the Roundfix spec reference, the write-tasks skill and the glossary state the subsumed selector rule

## Overview

task_01 adds `SC-VERIFY-SUBSUMED-SELECTOR`, and no document names it. This
Task adds its row to the `spec` command guide, its bullet to the Roundfix
`spec` reference, the authoring rule to the `write-tasks` skill as a section
of its own, and **Subsumed Test Selector** to the glossary, then raises both
skills' versions and re-records the Behavior Surface Record (`_prd.md` →
Core Feature 2; the Quaestor report of 2026-10-21 in the adopted Backlog
Entry).

## Requirements

1. MUST add to the code table of `docs/user-guide/commands/spec.md`, right
   after the `SC-VERIFY-PATH-UNDECLARED` row, this row:
   "| `SC-VERIFY-SUBSUMED-SELECTOR` | A pending non-QA Task's Verification
   passes a pytest node id (`<file>::<name>`) together with a collection
   argument of the same pytest invocation that contains it: its file, a
   directory above it, `.`, or a node id that prefixes it. pytest 9 runs the
   containing argument whole and drops the node id without an error, so a
   test that does not exist yet goes unnoticed and the command passes before
   the work. Reads `pytest`, `py.test` and `python -m pytest`, also behind
   `uv run`, `poetry run`, `pdm run`, `pipenv run` or `hatch run` and inside
   `$(...)`; option values, redirection targets and globs are not collection
   arguments, and `--keep-duplicates` exempts the invocation. An error with no
   horizon from the Tasks stage; Specs without a Task Graph skip with the
   manifest path. Reads files only and executes no command (ADR-0315). |",
   written on one line as the table's other rows are.
2. MUST add to the canonical Roundfix reference
   `.agents/skills/roundfix/references/spec.md`, right after the
   `SC-VERIFY-PATH-UNDECLARED` bullet that ends "and executes no command.",
   this bullet:
   "- `SC-VERIFY-SUBSUMED-SELECTOR` — a pending non-QA Task's Verification
   passes a pytest node id together with a collection argument of the same
   invocation that contains it (its file, a directory above it, `.`, or a
   node id that prefixes it), so pytest 9 runs that argument whole and the
   command passes before the test exists. `--keep-duplicates` exempts the
   invocation. An error with no horizon from the Tasks stage; reads files only
   and executes no command. Name the selector alone and run its containing
   path in a separate command (ADR-0315)."
3. MUST add to the canonical `.agents/skills/write-tasks/SKILL.md`, right
   after the `## Verification paths` section and before `## Owned skill
   versions`, this section:
   "## Test selectors and their containing path

   A Verification command that runs a new test by its selector names the
   selector alone. pytest 9 and later collapse collection arguments that
   overlap: when one invocation passes a node id such as
   `tests/test_a.py::test_new` and also its file, a directory above it or
   `.`, pytest runs the containing path whole and drops the node id without
   an error. Before the test exists the command passes on the tests the file
   already holds, and the pre-work probe refuses it as vacuous without naming
   why; with the node id alone pytest exits 4 and the command fails as it
   should. Run the containing path in a separate command when the Task also
   needs it.

   `SC-VERIFY-SUBSUMED-SELECTOR` reports an error for a pending non-QA Task at
   the Tasks stage, with no horizon. It reads each pytest invocation —
   `pytest`, `py.test` or `python -m pytest`, behind `uv run`, `poetry run`,
   `pdm run`, `pipenv run` or `hatch run`, and inside a `$(...)` substitution
   — and reports a node id that another collection argument of the same
   invocation contains: its file, a directory above it, `.`, or a node id
   that prefixes it. `--keep-duplicates` exempts the invocation. The check
   executes no command."
   No text inside any `### QA settlement` section of any skill may change.
4. MUST add **Subsumed Test Selector** to `CONTEXT.md` through
   `domain-modeling`, right after the **Non-Hermetic Verification** entry, as
   "**Subsumed Test Selector**:" followed on the next line by "The Spec
   Consistency Check error (`SC-VERIFY-SUBSUMED-SELECTOR`) raised when a
   pending non-QA Task's Verification passes a test selector, a pytest node id
   such as `tests/test_a.py::test_new`, together with a collection argument of
   the same invocation that contains it: its file, a directory above it, `.`,
   or a node id that prefixes it. pytest 9 runs the containing argument whole
   and drops the selector without an error, so a test that does not exist yet
   goes unnoticed and the command passes before the work (ADR-0315)." and the
   line "_Avoid_: Duplicate test path, redundant node id, shadowed test". The
   ADR reference takes the glossary's link form, to
   `docs/adr/0315-a-test-selector-another-argument-of-its-invocation-contains-is-refused-at-authoring.md`.
5. MUST, because this Task edits two owned skills and a guide, raise both
   version fields (`metadata.version` and the top-level `version`) of the
   canonical `SKILL.md` of the Roundfix skill and of the `write-tasks` skill
   by one patch over the value the tree holds when this Task starts, because
   the record command refuses changed content under an unchanged version, and
   then run these commands itself, in this order, before handing back; no
   mirror or record is hand-edited, and the Behavior Surface Record is
   re-recorded after the last guide edit:
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `spec` guide's fingerprint changes.
6. MUST NOT change any Go source, ADR, `docs/references/skill-coverage.json`,
   any other guide page, or the `qa-gate` skill.

## Subtasks

- [ ] Add the guide row, the reference bullet and the `write-tasks` section.
- [ ] Add **Subsumed Test Selector** to the glossary.
- [ ] Raise both skills' versions, sync the mirrors, record the versions and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide, both copies of the Roundfix reference and both copies of the `write-tasks` skill describe `SC-VERIFY-SUBSUMED-SELECTOR`.
- [ ] `CONTEXT.md` defines **Subsumed Test Selector** and links ADR-0315.
- [ ] Both skill versions are raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirrors, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0315-a-test-selector-another-argument-of-its-invocation-contains-is-refused-at-authoring.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'SC-VERIFY-SUBSUMED-SELECTOR' || { printf 'missing check code in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'exempts the invocation' || { printf 'missing exemption in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md; do grep -qx -- '## Test selectors and their containing path' "$f" || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Subsumed Test Selector**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'pytest 9 runs the containing argument whole and drops the selector without an error' || { echo 'missing glossary definition' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task no document names `SC-VERIFY-SUBSUMED-SELECTOR`, so the command fails at its first phrase check. Afterwards the guide, both reference copies and both `write-tasks` copies describe the check and its exemption, both skill copies hold the new section, the glossary defines **Subsumed Test Selector**, the mirrors are in sync, the skill versions are recorded, and the Skill Coverage Map is current.

## References

- `_prd.md` → Core Features 2 and 4; Goals; Glossary
- `_techspec.md` → API Contract 8; Build Order 2
- ADR-0315
