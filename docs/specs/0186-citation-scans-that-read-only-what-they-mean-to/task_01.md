---
task: task_01
spec: 0186-citation-scans-that-read-only-what-they-mean-to
status: pending
type: backend
complexity: low
---

# Task 01: The citation projection applies to the Spec's own Task files only

## Overview

`readSpecCitations` in `internal/speccheck/citations.go` walks every Markdown file under a Spec folder for `SC-ADR-UNLISTED`. It skips the Agent-owned `## Result` and the Daemon-owned `## Recorded paths` and `## Carry-forward provenance` sections in any file whose basename starts with `task_`, so an adopted source such as `references/task_example.md` loses those sections too. `references/` is authored and must be read in full. This Task limits the projection to the Spec folder's own Task files: files directly inside the Spec folder whose name matches `task_*.md`.

## Requirements

1. MUST project only files that are directly inside the Spec folder and whose name matches `task_*.md`, through a helper `projectedTaskFile(specDir, path string) bool`. Every other Markdown file under the Spec folder, including any file under `references/` or another subdirectory, MUST be read in full.
2. MUST keep the `qa/` directory skipped and keep every finding's summary, location and fix text unchanged.
3. MUST keep `TestASpecCitationInATaskResultIsNotAnObligation`, `TestAnAuthoredSpecCitationStillMustBeListed` and `TestCheckADRClosureDepthOne` green without editing them. It MUST rename or remove no top-level test and change no exported function signature.
4. MUST put the new tests in `internal/speccheck/citation_projection_scope_test.go`, over plain temporary directories.

## Subtasks

- [ ] Add `projectedTaskFile` and use it in `readSpecCitations`.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An unlisted ADR cited only in the `## Result` of `references/task_example.md` reports `SC-ADR-UNLISTED`.
- [ ] An unlisted ADR cited only in a nested `notes/task_02.md` `## Recorded paths` section reports `SC-ADR-UNLISTED`.
- [ ] The same citation in the top-level `task_01.md` `## Result` reports nothing.

## Context

- instruction: `docs/adr/0176-citation-checks-read-only-what-a-specs-authors-wrote.md`
- interface: `internal/speccheck/citations.go`
- creates: `internal/speccheck/citation_projection_scope_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAReferenceNamedLikeATaskIsReadInFull|TestANestedFileNamedLikeATaskIsReadInFull|TestTheSpecsOwnTaskFileKeepsItsProjection|TestASpecCitationInATaskResultIsNotAnObligation|TestAnAuthoredSpecCitationStillMustBeListed|TestCheckADRClosureDepthOne)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAReferenceNamedLikeATaskIsReadInFull TestANestedFileNamedLikeATaskIsReadInFull TestTheSpecsOwnTaskFileKeepsItsProjection TestASpecCitationInATaskResultIsNotAnObligation TestAnAuthoredSpecCitationStillMustBeListed TestCheckADRClosureDepthOne; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the three new named tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 3; Core Feature 1; Success Metric 1; Success Metric 3
- [_techspec.md](_techspec.md) — System Architecture; Testing Approach 1; Build Order 1
- [references/2026-09-29-citation-projection-matches-task-named-references.md](references/2026-09-29-citation-projection-matches-task-named-references.md)
- ADR-0176; ADR-0093

## Result
