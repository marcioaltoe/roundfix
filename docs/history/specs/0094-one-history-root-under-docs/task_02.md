---
task: task_02
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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
- `go test -count=1 ./internal/speccheck > /tmp/0094-task-02b.log 2>&1 && grep -q '"docs/history/findings"' internal/spec/archive.go || { cat /tmp/0094-task-02b.log; exit 1; }` — expected: exits 0, proving the consistency detectors still resolve retired references once the resolver answers under the new root. The suite alone passed before any work, so it is anchored to the relocation it is guarding.

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

## Result

### Implementation

- `ArchiveDir` now resolves Specs, findings, ADRs, backlog entries, and Review
  Artifacts under `docs/history/`; external and configured non-default Spec Roots
  keep their existing beside-root archive behavior.
- The live orphan Review Artifact root is `reviews` without an underscore. A
  Spec Root under repository history falls back to `docs/specs/reviews`, and an
  explicit Artifact Directory under repository history is refused. The
  Spec-owned `<spec>/reviews` path is unchanged.
- The repository's 2,225 tracked retired files moved from `_archived/{specs,
  findings,adr}` to `docs/history/{specs,findings,adr}`. `_archived/` no longer
  exists. No orphan Review Artifact moved.
- Both archive-layout characterization tests now pin the five new resolver
  answers, the three retired families present in this repository, and the Spec
  0094 corpus re-record reason. The focused Spec tests also stopped deriving an
  in-tree `_archived` fixture name from the repository history resolver.
- Verification-feedback repair separates current replay source lookup from
  frozen replay provenance: consistency checks still resolve source reports via
  `ArchiveDir`, while README assertions deliberately pin the historical
  `_archived/...` text. Replay fixtures and retired artifacts remain unchanged.
- Fourteen relocated `.log` files match ignore rules at their new paths. They
  are force-staged as additions so the Daemon's normal staging cannot retain
  their old deletions without their replacements; no path outside
  `docs/history/` is staged.

### Focused-check evidence

- All five resolver criteria: `GOCACHE=/tmp/roundfix-0094-task-02-gocache rtk
  go test ./internal/spec` passed all 275 tests, including the updated resolver,
  built-in/non-default root, characterization, and historical replay cases.
- Orphan-root and history-refusal criteria: `GOCACHE=/tmp/roundfix-0094-task-02-gocache
  rtk go test ./internal/config -run '^(TestResolveReviewRoot|TestReviewArtifactRootNeverResolvesIntoHistory)$'`
  passed all 9 tests. The negative test covers a slug found only under
  `docs/history/specs` and an explicit Artifact Directory under history.
- CLI consumer evidence: `GOCACHE=/tmp/roundfix-0094-task-02-gocache rtk go test
  ./internal/cli -run '^(TestRunFetchWritesReviewArtifactsUnderSpeclessRoot|TestReviewArtifactEvidenceMixedParentEmptyUserRootRefused|TestStageableReviewRootClassifiesInsideOutsideAndSymlink)$'`
  passed all 8 tests.
- Consistency-detector evidence: `GOCACHE=/tmp/roundfix-0094-task-02-gocache rtk
  go test ./internal/speccheck -run '^(TestCheckFindingLifecycle|TestCheckRollupMember|TestCheckArchiveLicense|TestCheckReplay0058QA001FromReport|TestCheckReplay0058QA004FromReport|TestCheckReplay0056F001FromReport|TestCheckReplay0056F002FromReport)$'`
  passed all 18 tests against the relocated Spec and finding roots.
- Verification-feedback evidence: `GOCACHE=/tmp/roundfix-0094-task-02-gocache
  rtk go test -count=1 ./internal/speccheck -run
  '^(TestCheckReplay0060Task03RefusesWorkIndependentVerification|TestCheckReplayReadmeProvenance)$'
  -v` passed all 6 tests, including the four provenance cases reported by the
  Daemon. A broader focused `-run '^TestCheckReplay'` check passed all 11 replay
  tests.
- Corpus-golden evidence: `GOCACHE=/tmp/roundfix-0094-task-02-gocache rtk go test
  -tags docscontract ./internal/docscontract -run '^TestCheckCorpusGolden$'`
  exited 0. The `update` field names Spec 0094 Task 02 and the active counts stay
  unchanged.
- Byte-identity evidence: sorted per-file SHA-1 manifests before and after the
  move both produced `ea29d6a45ba959753730b7a213864d7fc822446e`; both trees
  contained 2,225 files. The post-move staging audit reports 14 staged ignored
  replacements, 0 ignored unstaged replacements, and 2,211 normal unstaged
  replacements.
- Diff hygiene: `rtk git diff --check` and `rtk git diff --cached --check` both
  exited 0.

### Checks not used as evidence

- Initial focused Go checks did not reach test execution because the sandbox
  denied writes to the default Go build cache. The recorded passing reruns use
  the Task-specific cache under `/tmp`.
- A broader `rtk go test ./internal/cli` probe produced no output for more than
  five minutes and was interrupted with exit 130. The three affected CLI seams
  were then run directly and passed as recorded above.
- The Daemon-owned commands in `## Verification` were not run in this Agent
  turn.
- Daemon attempt 1 exposed replay provenance assertions that followed the live
  resolver after relocation. Its diagnostic artifact was inspected; the
  declared full-package command was not rerun during this repair.
