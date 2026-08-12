---
task: task_02
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 02: Resolve one history root and move this repository into it

## Overview

The resolver's answers move under the documentation tree and gain a fifth family
for Review Artifacts; the Review Artifact root drops its underscore and learns to
refuse the history tree; and this repository's own retired bytes move in the same
slice. They travel together on purpose: a resolver pointing at a location the
bytes have not reached leaves consistency checks resolving retired references
against an empty directory, so splitting them would settle a Task on a red tree.

## Requirements

1. MUST resolve every retired family under one history root inside the
   documentation tree, covering retired Specs, findings, decision records, intent
   entries, and Review Artifacts.
2. MUST resolve the orphan Review Artifact root without a leading underscore, so
   the live root and the Spec-owned root read the same.
3. MUST NOT resolve any Review Artifact root to a path under the history root,
   including for a Spec slug that resolves only inside history.
4. MUST leave the Spec-owned Review Artifact path unchanged; it already writes
   inside the Spec and already travels with it.
5. MUST relocate this repository's retired Specs and findings to the new root,
   preserving history for every moved path.
6. MUST leave every retired artifact byte-identical; this Task moves files and
   changes nothing inside them.
7. MUST update the two archive-layout characterization tests to the new layout in
   the same change, rather than deleting or weakening them.
8. MUST re-record the characterization corpus golden in the same change with the
   reason stated, and MUST NOT regenerate it silently.
9. MUST NOT change any carrier that merely names the location in prose, and MUST
   NOT relocate any orphan Review Artifact; both travel in their own Tasks.

## Subtasks

- [ ] Change the resolver's answers and add the Review Artifact family.
- [ ] Drop the underscore from the orphan Review Artifact root and refuse history.
- [ ] Relocate this repository's retired Specs and findings to the new root.
- [ ] Update the archive-layout characterization tests.
- [ ] Re-record the corpus golden with its reason.

## Acceptance Criteria

- [ ] The resolver answers under the documentation tree for all five families.
- [ ] The orphan Review Artifact root carries no underscore.
- [ ] No Review Artifact root resolves under the history root, proven for a slug
      that exists only inside history.
- [ ] Every retired Spec and finding resolves under the new root, and none
      remains at the old location.
- [ ] The characterization tests assert the new layout and fail against the old.
- [ ] The corpus golden records why it was re-recorded, in its own update field.
- [ ] No retired artifact's content changed.

## Verification

- `grep -q '"docs/history/specs"' internal/spec/archive.go && grep -q '"docs/history/findings"' internal/spec/archive.go && grep -q '"docs/history/adr"' internal/spec/archive.go && grep -q '"docs/history/backlog"' internal/spec/archive.go && grep -q '"docs/history/reviews"' internal/spec/archive.go` — expected: exits 0, proving all five answers moved rather than only the ones this Task's tests read. Fails today, where the resolver answers `_archived/…` and names no review family.
- `test -d docs/history/specs && test -d docs/history/findings && test ! -e _archived` — expected: exits 0, proving the bytes moved rather than only the resolver's answer. Fails today, where `_archived/` still exists.
- `! grep -q '"_reviews"' internal/config/config.go` — expected: exits 0, proving the orphan review root lost its underscore. Fails today, where that literal is present.
- `go test -count=1 ./internal/config -run 'TestReviewArtifactRootNeverResolvesIntoHistory' -v > /tmp/0094-task-02c.log 2>&1; s=$?; grep -q '^--- PASS: TestReviewArtifactRootNeverResolvesIntoHistory' /tmp/0094-task-02c.log || { cat /tmp/0094-task-02c.log; exit 1; }; exit $s` — expected: exits 0 and the log names that exact test. The name is exact rather than a `Review` substring, which would match existing resolution tests and pass before any work.
- `grep -q 'Re-recorded because Spec 0094' internal/docscontract/testdata/corpus-golden.json` — expected: exits 0. The golden already carries a `Re-recorded because` sentence from Spec 0085, so the reason must name this Spec for the check to mean anything.
- `go test -count=1 ./internal/spec -run ArchiveLayoutCharacterization -v > /tmp/0094-task-02.log 2>&1; s=$?; grep -q '^--- PASS: TestArchiveLayoutCharacterization' /tmp/0094-task-02.log && grep -q 'docs/history' internal/spec/archive_layout_characterization_test.go || { cat /tmp/0094-task-02.log; exit 1; }; exit $s` — expected: exits 0 with the characterization passing against the relocated layout. The second clause is what makes this non-vacuous: the test passes today against the old layout, so the gate must also prove it now pins the new one.
- `go test -count=1 ./internal/speccheck > /tmp/0094-task-02b.log 2>&1 || { cat /tmp/0094-task-02b.log; exit 1; }` — expected: exits 0, proving the consistency detectors still resolve retired references after the move. This one passes before the work as well; it is a regression guard for the move, and the commands above are what prove the Task's own effect.

## Context

- interface: `internal/spec/archive.go`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `internal/config/config.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`

## References

`_techspec.md` → Build Order 2; System Architecture: the history resolver and the
Review Artifact root; Decisions: `os.Rename` rather than `git mv`, and the
Spec-owned path left unchanged. `_prd.md` → Core Features 1, 3 and 5; Goal 1;
User Story 1. ADR-0120, ADR-0123.
