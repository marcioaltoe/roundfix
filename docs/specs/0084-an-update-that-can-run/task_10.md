---
task: task_10
spec: 0084-an-update-that-can-run
status: pending
type: qa
complexity: high
---

# Task 10: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It runs after every implementation leaf
settles and walks the PRD's user stories against the built binary. One row is
reserved for the evidence this Spec exists to institutionalize: a reading of the
real adopted repositories on the maintainer's machine, which no fixture in this
Spec can substitute for.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD user story against the built binary rather than against a
   test fixture.
3. MUST include one row whose evidence is a reading of the adopted repositories on
   the maintainer's machine: for each repository carrying a Setup Manifest, the
   read-only update outcome, recorded as state, exit code, and blocker.
4. MUST record that row as blocked with its reason, rather than dropping it, when
   a repository is absent from the machine or cannot be read.
5. MUST compare that reading against the outcomes the adopted finding recorded
   before this Spec, and state for each repository whether its blocker was
   removed, remains, or changed.
6. MUST record the origin of the external evidence used, per the clause this Spec
   seats.
7. MUST NOT mutate any repository outside this one: every fleet reading is
   read-only and passes no approval flag.
8. MUST report `rows_blocked_environment` and `rows_blocked_finding` counts in the
   report frontmatter.

## Subtasks

- [ ] Run the gate and write the report to the Spec's QA directory.
- [ ] Walk every user story against the built binary.
- [ ] Read every adopted repository on the machine, read-only.
- [ ] Compare each reading against the recorded pre-Spec outcome.
- [ ] Record the external evidence's origin and every blocked row's reason.

## Acceptance Criteria

- [ ] A QA report exists in the Spec's QA directory with a machine-readable
      verdict and both blocked-row counts in frontmatter.
- [ ] Every PRD user story has a row.
- [ ] One row carries the fleet reading, with a per-repository state, exit code,
      and blocker, or a recorded reason for each repository it could not read.
- [ ] That row states, per repository, whether the recorded blocker was removed,
      remains, or changed.
- [ ] The report names where its external evidence came from.
- [ ] No repository outside this one was mutated: every fleet command in the
      report is read-only.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/spec-routing.md`

## Verification

- `test -d docs/specs/0084-an-update-that-can-run/qa` — expected: exits 0.
- `ls docs/specs/0084-an-update-that-can-run/qa/*.md > /tmp/0084-task-10-a.log 2>&1 && grep -q 'qa/' /tmp/0084-task-10-a.log` — expected: exits 0, proving a report was written.
- `grep -q 'rows_blocked_environment' docs/specs/0084-an-update-that-can-run/qa/*.md` — expected: exits 0.
- `grep -q 'rows_blocked_finding' docs/specs/0084-an-update-that-can-run/qa/*.md` — expected: exits 0.
- `make verify > /tmp/0084-task-10-b.log 2>&1` — expected: exits 0, proving the authoritative gate passes on the finished Spec.

## References

- `_prd.md` → every User Story; Success Metrics.
- `_techspec.md` → Testing Approach.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the pre-Spec fleet outcomes this gate compares against.
- ADR-0091, ADR-0097, ADR-0104.
