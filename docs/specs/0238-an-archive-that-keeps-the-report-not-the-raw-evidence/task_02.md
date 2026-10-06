---
task: task_02
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: pending
type: backend
complexity: high
---

# Task 02: A normal archive drops the raw QA evidence, keeps a manifest, and the queue accepts it as exact

## Overview

This Task makes `roundfix archive <slug>` cut a Spec's `qa/evidence/`
directory as part of a normal archive. The archive writes
`qa/evidence-manifest.md` with each dropped file's path, size and SHA-256,
the revision and the source directory, and points the Spec's evidence links
at it. The Delivery Queue's exact-move check accepts that commit, so
delivery proceeds as before. A Spec without evidence, an override archive
and a superseded Spec archive exactly as today.

## Requirements

1. MUST add `internal/spec/archive_evidence.go` with the TechSpec's
   Interfaces for the manifest (`EvidenceManifestName`,
   `EvidenceManifestSchema`, `EvidenceManifestEntry`, `EvidenceManifest`,
   `ParseEvidenceManifest`, `RenderEvidenceManifest`) and an unexported
   planner that task_03 reuses. The manifest form follows TechSpec
   Invariants 2 and 3 and the Data Models example.
2. MUST make `Archive` in `internal/spec/archive.go` cut under Invariants 1,
   5 and 6: a normal archive of a Spec with a Task Graph and no QA override
   drops every file under `qa/evidence/`, writes the manifest when at least
   one file was dropped, and fills `ArchiveResult.DroppedEvidenceFiles` and
   `DroppedEvidenceBytes`. `ArchiveRequest` gains `EvidenceRevision` and
   `EvidenceSource`; a non-empty cut set with an empty `EvidenceRevision`
   refuses before the stamp.
3. MUST rewrite evidence links under Invariant 4, in the same pass and with
   the same scanner as ADR-0230's rewrites, and count them in
   `RewrittenLinks`.
4. MUST extend `ArchiveLinksMatch` under Invariant 7 without changing its
   signature.
5. MUST make the Archive Command in `internal/cli/archive.go` pass the
   revision from `git rev-parse HEAD` (of the repository that holds the Spec
   Root, as the override already does) and the repository-relative active
   Spec directory, and append the API Contract 1 suffix after the existing
   rewrite suffix, as in Surface Transcript 1. Its usage text states that a
   normal archive drops `qa/evidence/` and keeps every QA Report with
   `qa/evidence-manifest.md`, and keeps every substring
   `TestRunArchiveHelp` asserts.
6. MUST make `archiveCommitIsExact` in `internal/cli/deliver_workflow.go`
   accept the cut only under Invariant 8, reading the parent's evidence
   blobs in one batched `git cat-file` call. Every other exact-move rule
   stays as it is.
7. MUST describe the cut at archive in `docs/user-guide/commands/archive.md`:
   what is dropped, the manifest and its fields, the link rewrite, the
   confirmation suffix, and that override and superseded archives keep their
   files.
8. MUST add the tests of the TechSpec's Testing Approach for `internal/spec`
   and for the Archive Command and the exact move, in new files only:
   - `internal/spec/archive_evidence_test.go`:
     `TestArchiveDropsQAEvidenceAndWritesItsManifest`,
     `TestArchiveLinksIntoDroppedEvidenceReachTheManifest`,
     `TestArchiveWithoutQAEvidenceIsUnchanged`,
     `TestArchiveKeepsQAEvidenceOfAnOverrideAndASupersededSpec`,
     `TestArchiveRefusesQAEvidenceItCannotList` and
     `TestArchiveLinksMatchAcceptsAnEvidenceLinkThatReachesTheManifest`.
   - `internal/cli/archive_evidence_test.go`:
     `TestArchiveCommandReportsTheDroppedQAEvidence`, asserting Surface
     Transcript 1 on a temporary repository.
   - `internal/cli/deliver_archive_evidence_test.go`:
     `TestArchiveCommitWithAListedEvidenceCutIsAnExactMove` and
     `TestArchiveCommitWithAnUnlistedEvidenceCutIsNotAnExactMove` (a
     manifest that omits a file, a wrong digest, a wrong source, and an
     extra dropped path outside `qa/evidence/`).
9. MUST NOT change any existing test, `internal/cli/archive_test.go`, the QA
   eligibility code, the review scope, any skill or Baseline asset, or any
   file under `docs/history/`.

## Subtasks

- [ ] Add the manifest type, renderer, parser and cut planner.
- [ ] Cut inside a normal archive and rewrite evidence links.
- [ ] Report the cut from the Archive Command.
- [ ] Accept a listed cut in the queue's exact-move check.
- [ ] Document the cut in the archive user guide and add the tests.

## Acceptance Criteria

- [ ] An archive of a Spec with evidence leaves no `qa/evidence/`, a
      manifest whose rows match the dropped files' sizes and digests, and
      every former evidence link resolving to the manifest.
- [ ] A Spec without evidence, an override archive and a superseded Spec
      archive byte-for-byte as before.
- [ ] A refused cut changes no file.
- [ ] The exact-move check accepts a listed cut and refuses every unlisted
      or mismatched one.

## Context

- instruction: `docs/adr/0243-an-archive-keeps-the-qa-report-and-drops-the-raw-qa-evidence.md`
- instruction: `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`
- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`
- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `docs/user-guide/commands/archive.md`
- creates: `internal/spec/archive_evidence.go`
- creates: `internal/spec/archive_evidence_test.go`
- creates: `internal/cli/archive_evidence_test.go`
- creates: `internal/cli/deliver_archive_evidence_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveDropsQAEvidenceAndWritesItsManifest|TestArchiveLinksIntoDroppedEvidenceReachTheManifest|TestArchiveWithoutQAEvidenceIsUnchanged|TestArchiveKeepsQAEvidenceOfAnOverrideAndASupersededSpec|TestArchiveRefusesQAEvidenceItCannotList|TestArchiveLinksMatchAcceptsAnEvidenceLinkThatReachesTheManifest|TestArchiveCommandReportsTheDroppedQAEvidence|TestArchiveCommitWithAListedEvidenceCutIsAnExactMove|TestArchiveCommitWithAnUnlistedEvidenceCutIsNotAnExactMove|TestRunArchiveHelp)$' ./internal/spec ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveDropsQAEvidenceAndWritesItsManifest TestArchiveLinksIntoDroppedEvidenceReachTheManifest TestArchiveWithoutQAEvidenceIsUnchanged TestArchiveKeepsQAEvidenceOfAnOverrideAndASupersededSpec TestArchiveRefusesQAEvidenceItCannotList TestArchiveLinksMatchAcceptsAnEvidenceLinkThatReachesTheManifest TestArchiveCommandReportsTheDroppedQAEvidence TestArchiveCommitWithAListedEvidenceCutIsAnExactMove TestArchiveCommitWithAnUnlistedEvidenceCutIsNotAnExactMove TestRunArchiveHelp; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'listed in qa/evidence-manifest.md' || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md 'listed in qa/evidence-manifest.md' >&2; exit 1; }` — expected: exit 0; before this Task the nine new tests do not exist and the archive guide does not name the manifest, so the command fails; after it the new tests and the existing help test pass and the guide quotes the confirmation suffix.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; Invariants 1-8; Data Models; API Contract 1; API Contract 4; Surface Transcript 1; Testing Approach; Build Order 2
- ADR-0243; ADR-0230; ADR-0154
