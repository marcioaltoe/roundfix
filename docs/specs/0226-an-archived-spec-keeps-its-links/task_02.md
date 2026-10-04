---
task: task_02
spec: 0226-an-archived-spec-keeps-its-links
status: pending
type: backend
complexity: high
---

# Task 02: The archive rewrites outward links and refuses a broken one

## Overview

Add the link pass of the TechSpec to `spec.Archive`: scan the Spec's Markdown
files, rewrite each relative destination that leaves the Spec to reach the
same path from the archived location, keep one whose target was already
archived, and refuse before any write when one reaches nothing. Export the
match function the Delivery Queue's resume path needs, and append the rewrite
count to the command's confirmation line. This answers the Backlog Entry of
2026-10-03, "Archiving a Spec breaks its relative links that leave the Spec",
measured on fluxus on 2026-10-02.

## Requirements

1. MUST implement "The link pass" steps 1-6 in `internal/spec/archive.go`,
   classifying by the target a destination reaches (lexical resolution, no
   symbolic link evaluation), for the built-in and any configured Spec Root and
   for a superseded Spec.
2. MUST refuse with API Contract 1 before the archive stamp, writing no file,
   when any destination that leaves the Spec reaches nothing from either
   location; and MUST restore every rewritten file's bytes when the rename
   fails.
3. MUST add `ArchiveResult.RewrittenLinks` and the exported
   `ArchiveLinksMatch` of the TechSpec's Interfaces in
   `internal/spec/archive.go`, built on the same scanner the pass uses.
4. MUST append `; rewrote <n> relative link(s)` to both confirmation lines in
   `internal/cli/archive.go` only when n is greater than 0, and add the usage
   sentence of API Contract 3 while keeping "moves unchanged".
5. MUST change `prepareSpec0058Replay` in the governed
   `internal/spec/archive_test.go` only to create, under the replay's temporary
   repository, the three files archived Spec 0058's PRD links to, under the
   maintainer's named grant of 2026-10-04 ("Concedo") recorded in
   `_authorization.md`, leaving every assertion unchanged.
6. MUST add the tests of the TechSpec's Testing Approach 1 and 2 in the two new
   files, and MUST NOT edit any other existing test.
7. MUST NOT rewrite non-Markdown files, change archive eligibility, the QA
   override, the archive stamp or the archive destinations.

## Subtasks

- [ ] Scan, classify and refuse before any write.
- [ ] Write the rewrites before the rename and restore them on failure.
- [ ] Export the match function and append the count to the confirmation line.
- [ ] Create the replay's link targets and add the eight tests.

## Acceptance Criteria

- [ ] Outward links reach the same files from the archived location.
- [ ] A broken outward link refuses with no byte changed.
- [ ] A link whose target was already archived keeps its bytes.
- [ ] The confirmation line carries the count only when links were rewritten.

## Context

- interface: `internal/spec/archive.go`
- interface: `internal/spec/archive_test.go`
- interface: `internal/cli/archive.go`
- creates: `internal/spec/archive_links_test.go`
- creates: `internal/cli/archive_links_test.go`
- instruction: `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveRewritesRelativeLinksThatLeaveTheSpec|TestArchiveKeepsLinksInsideTheSpecAndNonRelativeLinks|TestArchiveRefusesALinkThatWouldStayBroken|TestArchiveRewritesLinksUnderAConfiguredSpecRoot|TestArchiveRewritesLinksInASupersededSpec|TestArchiveKeepsALinkWhoseTargetWasAlreadyArchived|TestSpec0058ReplayArchivesDeclaredUnreachableRelease|TestArchiveCommandRefusesABrokenOutwardLink|TestArchiveCommandReportsRewrittenLinks)$" ./internal/spec ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestArchiveRewritesRelativeLinksThatLeaveTheSpec TestArchiveKeepsLinksInsideTheSpecAndNonRelativeLinks TestArchiveRefusesALinkThatWouldStayBroken TestArchiveRewritesLinksUnderAConfiguredSpecRoot TestArchiveRewritesLinksInASupersededSpec TestArchiveKeepsALinkWhoseTargetWasAlreadyArchived TestSpec0058ReplayArchivesDeclaredUnreachableRelease TestArchiveCommandRefusesABrokenOutwardLink TestArchiveCommandReportsRewrittenLinks; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eight new tests do not exist, so the command fails; after it they pass and the Spec 0058 replay archives with its link targets present.

## References

- `_prd.md` → Core Features 1-5; User Stories 1-2; Success Metrics 1-2
- `_techspec.md` → The link pass; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Testing Approach 1; Testing Approach 2; Build Order 2
- ADR-0230
