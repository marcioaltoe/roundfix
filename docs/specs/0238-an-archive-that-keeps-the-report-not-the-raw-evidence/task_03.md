---
task: task_03
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: pending
type: backend
complexity: medium
---

# Task 03: The Archive Command cuts one archived Spec on request, as a dry run unless applied

## Overview

With task_02's planner in place, this Task lets the maintainer apply the
same cut to one already-archived Spec:
`roundfix archive <slug> --drop-evidence` reports what would be dropped,
rewritten and written, and `--apply` does it. It never runs implicitly,
never commits, and refuses a Spec whose evidence directory a file other
than Markdown still names. Nothing in this Task runs it on this
repository's History Root.

## Requirements

1. MUST add `CutArchivedEvidence` with the TechSpec's
   `ArchivedEvidenceCutRequest`, `ArchivedEvidenceDisposition` and
   `ArchivedEvidenceCutResult` in a new
   `internal/spec/archive_evidence_cut.go`, reusing task_02's planner, under Invariant 9. A dry run and an apply
   plan the same set, and a dry run writes nothing.
2. MUST add `ArchivedEvidencePathPins` in a new
   `internal/speccheck/archived_evidence_paths.go`, with the same search
   rules as `ActiveSpecPathPins` (tracked and untracked non-ignored files
   other than Markdown, outside the Spec Root, its archive root and
   `docs/history`) and the needle `<archive root relative>/<slug>/qa/evidence`.
3. MUST add `--drop-evidence` and `--apply` to the Archive Command in
   `internal/cli/archive.go` under Invariant 10 and API Contracts 2 and 3:
   pins and `git status --porcelain --untracked-files=all` on the evidence
   directory refuse before planning, and `--apply` records the HEAD
   revision and the repository-relative archived Spec directory. Output
   follows Surface Transcripts 2 to 6 and the superseded line. The usage
   text lists both flags and states that the dry run changes nothing.
4. MUST describe `--drop-evidence` and `--apply` in
   `docs/user-guide/commands/archive.md`, with the exact command
   `roundfix archive <slug> --drop-evidence --apply`, the dispositions, the
   refusals, and that the command never commits and never runs on its own.
5. MUST add `internal/cli/archive_drop_evidence_test.go` with, on temporary
   repositories and homes:
   - `TestArchiveDropEvidenceDryRunChangesNothing`: Surface Transcript 2,
     and every file under the repository byte-identical afterwards;
   - `TestArchiveDropEvidenceApplyCutsWhatTheDryRunReported`: Surface
     Transcripts 3 and 4, a manifest listing exactly the files the dry run
     counted, with source equal to the archived directory and revision equal
     to HEAD, and the report's evidence link resolving to the manifest;
   - `TestArchiveDropEvidenceKeepsOverrideAndSupersededSpecs`: Surface
     Transcript 5 and the superseded line, with no file changed;
   - `TestArchiveDropEvidenceRefusals`: Surface Transcript 6, a missing
     archived Spec, a pin in a Go file naming the evidence directory (and no
     refusal for the same text in a Markdown file), an uncommitted change
     under the evidence directory, a manifest beside remaining evidence, and
     the usage errors of API Contract 3, each exiting 2 with no file
     changed.
6. MUST NOT change any existing test, any skill, any Baseline asset or any
   file under `docs/history/`.

## Subtasks

- [ ] Add the archived cut on top of the planner.
- [ ] Add the archived-evidence pin search.
- [ ] Add the two flags, their output and their refusals.
- [ ] Document the flags in the archive user guide.
- [ ] Add the command tests.

## Acceptance Criteria

- [ ] A dry run changes no file and predicts exactly what an apply cuts.
- [ ] A second apply reports nothing to drop.
- [ ] Override and superseded Specs keep their files.
- [ ] Every refusal exits 2 and changes no file.

## Context

- instruction: `docs/adr/0243-an-archive-keeps-the-qa-report-and-drops-the-raw-qa-evidence.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `internal/speccheck/active_spec_paths.go`
- instruction: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`
- interface: `docs/user-guide/commands/archive.md`
- creates: `internal/spec/archive_evidence_cut.go`
- creates: `internal/speccheck/archived_evidence_paths.go`
- creates: `internal/cli/archive_drop_evidence_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveDropEvidenceDryRunChangesNothing|TestArchiveDropEvidenceApplyCutsWhatTheDryRunReported|TestArchiveDropEvidenceKeepsOverrideAndSupersededSpecs|TestArchiveDropEvidenceRefusals|TestRunArchiveHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveDropEvidenceDryRunChangesNothing TestArchiveDropEvidenceApplyCutsWhatTheDryRunReported TestArchiveDropEvidenceKeepsOverrideAndSupersededSpecs TestArchiveDropEvidenceRefusals TestRunArchiveHelp; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'roundfix archive <slug> --drop-evidence --apply' || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/archive.md 'roundfix archive <slug> --drop-evidence --apply' >&2; exit 1; }` — expected: exit 0; before this Task the four new tests do not exist and the archive guide does not name the flags, so the command fails; after it the new tests and the existing help test pass and the guide quotes the command.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 6; Success Metric 3
- `_techspec.md` → Interfaces; Invariants 9-10; API Contract 2; API Contract 3; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Testing Approach; Build Order 3
- ADR-0243; ADR-0223
