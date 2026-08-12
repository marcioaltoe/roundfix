---
task: task_04
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 04: Discover retired material sitting on an older layout

## Overview

A read-only walk that answers what a repository would have to move to reach the
current layout, across both shapes an older Roundfix produced: an archive
directory nested inside a documentation tree, and an archive root outside the
documentation tree. It reports relocations and collisions and performs neither,
which is what lets the planning and applying Tasks stay small and lets this one
be tested against fixture repositories alone.

## Requirements

1. MUST report one relocation per file, never per directory, so a partially
   migrated tree is describable and one refusal does not hide its siblings.
2. MUST discover retired material nested inside a documentation tree, retired
   material at an archive root outside the documentation tree, and orphan Review
   Artifacts under the live review root.
3. MUST report a relocation for a decision record or intent entry that the
   retirement classification calls retired and that sits in an active directory.
4. MUST NOT report a relocation for a record the retirement classification calls
   active, including a pending proposal.
5. MUST report a relocation for an orphan Review Artifact the liveness check calls
   finished, and MUST NOT report one for a review it calls live or undecidable.
6. MUST NOT report a relocation for a Spec-owned Review Artifact, which travels
   with its Spec.
7. MUST read and report each file's content identity before any mutation is
   planned, so a later stage can prove the bytes survived.
8. MUST report a destination that already exists as a collision rather than as a
   relocation, naming both paths and the reason.
9. MUST report nothing for a repository already on the current layout.
10. MUST return relocations in a stable order, so two runs over one tree agree.
11. MUST NOT move, create, or delete any file.

## Subtasks

- [ ] Walk every legacy shape and emit per-file relocations.
- [ ] Route decision records and intent entries through the retirement rule.
- [ ] Route orphan Review Artifacts through the liveness rule.
- [ ] Read content identity for each discovered file.
- [ ] Detect and report destination collisions.
- [ ] Cover each shape, the already-migrated case, and a collision with fixtures.

## Acceptance Criteria

- [ ] A fixture repository with an archive nested in a documentation tree reports
      one relocation per file it holds.
- [ ] A fixture repository with an archive root outside the documentation tree
      reports one relocation per file it holds.
- [ ] A fixture repository already on the current layout reports no relocation
      and no collision.
- [ ] A retired decision record in the active directory reports a relocation; a
      pending proposal in the same directory reports none.
- [ ] An orphan review whose liveness is finished reports relocations for its
      files; one whose liveness is live or undecidable reports none.
- [ ] A Spec-owned review directory reports no relocation.
- [ ] An occupied destination reports a collision naming both paths, and the
      other files in the same tree still report their relocations.
- [ ] Every reported relocation carries the source's content identity.
- [ ] The discovery leaves each fixture repository unchanged on disk.

## Verification

- `go test -count=1 ./internal/baseline -run 'DiscoverHistoryLayout' -v > /tmp/0094-task-04.log 2>&1; s=$?; grep -q '^--- PASS: .*DiscoverHistoryLayout' /tmp/0094-task-04.log || { cat /tmp/0094-task-04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing discovery tests; fails when the named tests do not exist.
- `grep -q 'PASS' /tmp/0094-task-04.log && ! grep -qi 'no tests to run' /tmp/0094-task-04.log` — expected: exits 0, refusing a vacuous run that selected no cases.
- `go build -buildvcs=false ./...` — expected: exits 0.

## Context

- interface: `internal/gittest/gittest.go`

## References

`_techspec.md` → Build Order 4; Interfaces: `DiscoverHistoryLayout`,
`HistoryRelocation`; Testing Approach: layout discovery; Risks: fifty orphan
folders, one decision each. `_prd.md` → Core Features 2, 4, 6 and 7; Goals 2, 3
and 4; User Stories 2, 3 and 4. ADR-0122, ADR-0123.
