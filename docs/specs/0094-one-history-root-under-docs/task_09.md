---
task: task_09
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: qa
complexity: high
---

# Task 09: Run the final QA gate

## Overview

The authored terminal gate. It walks every user story through the surfaces the
Spec declares, rests one row on a repository this Spec did not build, and carries
the glossary check that closes the Spec's vocabulary. It runs after every
implementation branch settles and before any Pull Request exists.

## Requirements

1. MUST derive the QA matrix from the PRD's user stories and the task evidence,
   and execute each row through a declared surface rather than by reading
   artifacts.
2. MUST rest one named row on evidence originating outside this Spec's artifacts:
   a repository on the pre-0085 nested layout that this Spec did not build,
   migrated by the built binary, with the row recording which repository and how
   many archived files it held.
3. MUST perform the outside-evidence migration against a disposable copy, never
   against a working checkout another session may own.
4. MUST classify every finding by user impact and record auditable evidence per
   row.
5. MUST check whether this Spec introduced, changed, or retired a term the
   glossary should carry, and update the domain context when it found something.
6. MUST record a row as blocked with its reason rather than passing it on
   equivalent artifacts when its surface cannot be reached.
7. MUST write the dated QA report into this Spec's evidence directory.
8. MUST NOT commit scratch repositories, built binaries, or any file it did not
   author as evidence.

## Subtasks

- [ ] Derive the resumable QA matrix from the PRD and task evidence.
- [ ] Execute every row through its declared surface.
- [ ] Run the outside-evidence migration against a disposable fleet copy.
- [ ] Run the glossary check and update the domain context if it found something.
- [ ] Write the dated QA report with per-row evidence.

## Acceptance Criteria

- [ ] Every PRD user story has at least one executed row with recorded evidence.
- [ ] The outside-evidence row names the fleet repository, the count of archived
      files it held before migration, and the result after it. Source: a
      repository adopted before Spec 0085 and still on the nested layout —
      measured 2026-08-12 across four candidates holding 173 to 735 archived
      files each. Blocked with that reason recorded if no copy can be obtained.
- [ ] The migrated copy's archived files are byte-identical to their originals.
- [ ] Re-running the migration against the already-migrated copy reports no
      relocation.
- [ ] The glossary check ran, and its outcome is recorded whether or not it
      changed the domain context.
- [ ] The QA report carries a verdict, per-row classification, and evidence
      paths.
- [ ] No scratch repository, built binary, or unauthored file is committed as
      evidence.

## Verification

- `ls docs/specs/0094-one-history-root-under-docs/qa/qa-report-*.md > /dev/null 2>&1 || { echo 'no QA report written'; exit 1; }` — expected: exits 0, proving the gate wrote its report. An earlier draft redirected stderr into the log it then tested for content, so `ls` failing wrote its own error message and satisfied the check.
- `grep -q '^verdict:' $(ls -t docs/specs/0094-one-history-root-under-docs/qa/qa-report-*.md | head -1)` — expected: exits 0, proving the newest report carries a verdict.
- `! git ls-files -s docs/specs/0094-one-history-root-under-docs/qa | grep -q '^160000'` — expected: exits 0, proving no scratch repository was committed as a gitlink.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## References

`_prd.md` → User Stories 1, 2, 3, 4, 5 and 6, Success Metrics. User Story 6 is
covered here rather than by a Task of its own: publishing a release is delivery
the graph deliberately excludes, and what this Spec can prove is the condition
behind it — that a built binary migrates a repository this Spec did not build.
The outside-evidence row is that proof. `_techspec.md` → Coverage Map; Risks &
Considerations. ADR-0091, ADR-0104, ADR-0120, ADR-0123.
