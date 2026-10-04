---
task: task_03
spec: 0226-an-archived-spec-keeps-its-links
status: completed
type: backend
complexity: medium
---

# Task 03: A resumed archive commit accepts the link rewrites

## Overview

When the Delivery Queue resumes an item whose archive commit already exists,
`archiveCommitIsExact` requires every file but the PRD to be blob-identical
and the PRD body unchanged. A commit whose archive rewrote links would read as
inexact and park. This Task accepts exactly the rewrites the archive makes, as
the TechSpec's "The resumed archive commit" states, through
`spec.ArchiveLinksMatch`. It completes the Backlog Entry of 2026-10-03,
"Archiving a Spec breaks its relative links that leave the Spec", for
delivered Specs.

## Requirements

1. MUST compare the archive commit's trees entry by entry in
   `archiveCommitIsExact` in `internal/cli/deliver_workflow.go`: an equal entry
   passes; a differing entry passes only when it is a Markdown file of the same
   mode and type whose blobs satisfy `spec.ArchiveLinksMatch`.
2. MUST accept a PRD body that is byte-identical or matches the same way,
   keeping the existing frontmatter comparison in `archivePRDChangeIsExact`.
3. MUST keep refusing any other content change, a link rewritten to another
   target, an added or removed file, and a link rewritten inside PRD
   frontmatter; and MUST leave `archiveDiffIsExact` unchanged.
4. MUST add the tests of the TechSpec's Testing Approach 3 in the new file
   `internal/cli/deliver_archive_links_test.go`, and MUST NOT edit any existing
   test.

## Subtasks

- [ ] Compare the trees entry by entry with the link match.
- [ ] Accept a matching PRD body.
- [ ] Add the acceptance and refusal tests.

## Acceptance Criteria

- [ ] A resumed archive commit with rewritten links is exact.
- [ ] Any other change in such a commit is still refused.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_archive_links_test.go`
- instruction: `internal/spec/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestResumeAcceptsAnArchiveCommitWithRewrittenLinks|TestResumeRefusesALinkRewritingArchiveCommitWithExtraChanges)$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestResumeAcceptsAnArchiveCommitWithRewrittenLinks TestResumeRefusesALinkRewritingArchiveCommitWithExtraChanges; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails; after it they pass.

## References

- `_prd.md` → Core Feature 6; User Story 3; Success Metric 3
- `_techspec.md` → The resumed archive commit; API Contract 4; Testing Approach 3; Build Order 3
- ADR-0223; ADR-0230

## Result

Implemented the resumed archive comparison in `internal/cli/deliver_workflow.go`.
Equal tree entries retain the fast path; differing entries must be regular
Markdown blobs with identical mode and type and must satisfy
`spec.ArchiveLinksMatch`, using each file's active and archived directories.
The PRD body accepts identical bytes or the same link match. Its existing
frontmatter comparison and `archiveDiffIsExact` behavior are preserved.

Added `internal/cli/deliver_archive_links_test.go`; no existing test was edited.
The tests create real archive commits and resume them through the Delivery
Engine using the existing delivery test boundary.

Acceptance evidence from focused implementation checks:

- **A resumed archive commit with rewritten links is exact:** the new
  acceptance test covers PRD-only, Task-only, nested Markdown-only, and combined
  rewrites, asserting that the archive head joins the reviewed candidate head.
  Before the implementation, the PRD-only subtest parked as `review-stale`.
  After the implementation, all four cases passed.
- **Any other change in such a commit is still refused:** all twelve new
  refusal cases passed, asserting a `review-stale` park and preservation of
  only the reviewed candidate head. Cases cover extra Task/PRD bytes, different
  Task/PRD link targets, a frontmatter destination change, an added or removed
  file, mode/type changes, a non-Markdown rewrite, and query/fragment changes.
  The two existing archive-resume tests also passed, including refusal of an
  unrelated path and changed PRD body.

Commands and outcomes:

- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 ./internal/cli -run '^TestResumeAcceptsAnArchiveCommitWithRewrittenLinks$/_prd.md$'`
  — before implementation, exit 1 with the expected `review-stale` park.
- `rtk proxy gofmt -w internal/cli/deliver_workflow.go internal/cli/deliver_archive_links_test.go`
  — exit 0.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 ./internal/cli -run '^TestResume(AcceptsAnArchiveCommitWithRewrittenLinks|RefusesALinkRewritingArchiveCommitWithExtraChanges|AcceptsARealArchiveCommit|RefusesAnArchiveCommitWithExtraChanges)$'`
  — exit 0, `ok roundfix/internal/cli`.

Only the implementation file, new test file, and this Result were changed by
the Agent. The pre-existing `status: in_progress` remains Daemon-owned.
Declared Verification was not run; Task settlement remains with the Daemon.
No follow-up work identified.
