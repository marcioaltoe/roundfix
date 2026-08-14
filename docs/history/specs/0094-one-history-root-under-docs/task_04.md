---
task: task_04
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Context

- interface: `internal/gittest/gittest.go`

## References

`_techspec.md` → Build Order 4; Interfaces: `DiscoverHistoryLayout`,
`HistoryRelocation`; Testing Approach: layout discovery; Risks: fifty orphan
folders, one decision each. `_prd.md` → Core Features 2, 4, 6 and 7; Goals 2, 3
and 4; User Stories 2, 3 and 4. ADR-0122, ADR-0123.

## Result

### Implementation

- `DiscoverHistoryLayout` performs a read-only repository walk and returns one
  source-sorted `HistoryRelocation` per regular file. Each relocation carries
  repository-relative source and destination paths plus a `sha256:` content
  identity read from the source bytes.
- Discovery recognizes both prior archive layouts:
  `docs/specs/_archived` and `docs/findings/_archived`, plus the interim
  repository-root `_archived/{specs,findings,adr,backlog}` tree. Every file keeps
  its path relative to the retired-family root at the corresponding
  `docs/history/` destination.
- Active `docs/adr` and `docs/backlog` files route through `ClassifyADR` and
  `ClassifyBacklogEntry`. Retired or declined documents relocate; proposed,
  accepted, open, and otherwise active documents do not.
- Orphan Review Artifacts under both `docs/specs/_reviews` and
  `docs/specs/reviews` route through `ClassifyReview`. Only a `finished` review
  contributes its files. The walk never enters `<spec>/reviews`, so Spec-owned
  Review Artifacts remain with their Spec.
- An occupied destination produces a `HistoryCollision` with both paths, the
  source identity, and the refusal reason while unaffected siblings remain
  relocations. Two legacy sources claiming the same destination are also
  collisions rather than an ambiguous plan.
- Discovery uses only filesystem reads and the existing read-only local-Git
  liveness check. It creates, moves, rewrites, and deletes nothing.

### Focused-check evidence

- The pre-implementation command, `rtk go test -count=1 ./internal/baseline
  -run '^TestDiscoverHistoryLayout'`, failed to compile because
  `DiscoverHistoryLayout`, `HistoryRelocation`, and `HistoryCollision` did not
  exist.
- The same focused command passed seven tests after implementation. The nested
  archive fixture reports three per-file relocations; the repository-root
  archive fixture reports four, covering Specs, findings, ADRs, and backlog
  entries. Exact expected slices prove source ordering and destination routing,
  and a second discovery over each unchanged fixture returns identical results.
- `TestDiscoverHistoryLayoutCurrentLayoutReportsNothing` passed with all five
  retired families under `docs/history/` and active documents in their live
  directories, returning no relocation or collision.
- `TestDiscoverHistoryLayoutClassifiesActiveDocuments` passed: a superseded ADR
  and declined backlog entry relocate, while a proposed ADR and open backlog
  entry do not.
- `TestDiscoverHistoryLayoutClassifiesOrphanReviews` passed against an isolated
  real Git repository. It reports all three files belonging to finished orphan
  reviews across the current and underscored roots, and reports none from the
  live, undecidable, or Spec-owned Review Artifacts.
- `TestDiscoverHistoryLayoutReportsCollisionWithoutHidingSiblings` passed: the
  occupied `_prd.md` destination reports one collision naming source,
  destination, and `destination already exists`, while `task_01.md` remains a
  relocation.
- Every exact relocation expectation includes the SHA-256 identity independently
  computed from its fixture bytes. The collision fixture also checks its source
  identity.
- Every fixture snapshots the path-to-content-identity map before discovery and
  compares it after discovery. The Review fixture includes the repository's
  `.git` files, so the local-Git liveness reads are covered by the same
  no-mutation assertion.
- `rtk go test -count=1 ./internal/baseline` exited 0 for the complete package,
  and `rtk go vet ./internal/baseline` exited 0.
- `rtk make verify-incremental` exited 0 after the implementation and Result
  edits, covering formatting, the Go suite, skill checks, and the build with
  reusable caches.

The Daemon-owned commands in `## Verification` were not run in this Agent turn.
