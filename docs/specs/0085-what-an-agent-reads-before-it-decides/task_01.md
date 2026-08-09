---
task: task_01
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: test
complexity: medium
---

# Task 01: Record the archive paths and the conditional clause today

## Overview

The corpus every later Task is measured against. It captures where each retired
artifact family lives today, which packages compose that path themselves, and the
current conditional wording of the Secondbrain consultation clause. It also
records the corpus-golden counts that Task 04's relocation will move, so that
move is a declared break rather than a surprise.

## Requirements

1. MUST record the current archive location for each retired artifact family:
   Specs, findings, ADRs, and backlog entries.
2. MUST record every package that composes an archive path itself rather than
   asking one owner, and the two hardcoded literals in the Spec checker.
3. MUST record the current conditional wording of the Secondbrain consultation
   clause in the Baseline catalog.
4. MUST record the corpus-golden counts that a relocation changes, and declare
   the break against Task 04.
5. MUST NOT change any production behaviour or move any artifact.

## Subtasks

- [ ] Capture the four archive locations.
- [ ] Enumerate the path-composing packages and the two literals.
- [ ] Capture the conditional clause and the golden counts.

## Acceptance Criteria

- [ ] A test asserts today's directory for each of the four families.
- [ ] The path-composing packages are enumerated with the literals they hold.
- [ ] The conditional clause wording is captured.
- [ ] The golden counts are captured with the break declared against Task 04.

## Bounded scope

This Task may create or modify only:

- `internal/spec/archive_layout_characterization_test.go`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_01.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveLayoutCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily'` — expected: exits 0. A `-run` pattern selecting no cases exits 0, so this asserts the named case ran.
- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveLayoutCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestArchiveLayoutCharacterizationEnumeratesEveryPathComposer'` — expected: exits 0.
- `grep -q 'Declared break: task_04' internal/spec/archive_layout_characterization_test.go` — expected: exits 0.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 1; System Architecture.
