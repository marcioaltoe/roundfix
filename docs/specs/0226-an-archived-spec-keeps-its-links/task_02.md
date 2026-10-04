---
task: task_02
spec: 0226-an-archived-spec-keeps-its-links
status: completed
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


## Result

Implemented the archive link pass after eligibility and destination checks,
before the archive stamp. It scans regular Markdown files recursively, keeps
internal and non-relative destinations, and skips symlinks, non-Markdown files,
fenced code and inline code. Outward destinations retain their lexical target,
query and fragment; angle brackets and titles retain their surrounding bytes.
Every unresolved outward destination is reported with its Spec-relative file
and line before any write. Rewritten files retain their original bytes for
rollback if writing or renaming fails.

`ArchiveResult.RewrittenLinks` counts changed destinations. The exported
`ArchiveLinksMatch` uses the same scanner and rejects changed targets, queries,
fragments and bytes outside destinations. Both command confirmation paths
append the count only when positive, and usage includes API Contract 3 while
retaining “moves unchanged”. The granted existing-test edit only supplies the
three PRD link targets in `prepareSpec0058Replay`; no assertion changed.

Focused evidence from this turn:

- Initial focused compile failed because `ArchiveResult.RewrittenLinks` did
  not exist. A later Markdown-form regression exposed an unfinished link being
  scanned; the scanner now requires complete link syntax and handles an image
  inside a linked label.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec ./internal/cli -run '(Archive|Spec0058Replay)' -count=1`
  exited 0 after the final code edits: both packages passed, including the
  existing Spec 0058 replay assertions.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk make verify-incremental`
  exited 0 when rerun with local listener and process-table permissions and
  the worktree held unchanged during execution. Includes formatting, vet,
  package tests, skill checks and build. The first sandboxed attempt exited 2:
  existing process-owner and HTTP-server tests lacked permissions, and the
  suite guard detected this Agent's concurrent Result edit. The rerun changed
  neither code nor assertions.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

| Acceptance criterion | Focused evidence |
| --- | --- |
| Outward links reach the same files after archive | `TestArchiveRewritesRelativeLinksThatLeaveTheSpec`, `TestArchiveRewritesLinksUnderAConfiguredSpecRoot`, `TestArchiveRewritesLinksInASupersededSpec`, and `TestArchiveLinksMarkdownForms` passed; covers nested Markdown, images, reference definitions, titles, encoded paths, queries and fragments. |
| Broken outward link refuses with no byte changed | `TestArchiveRefusesALinkThatWouldStayBroken` passed with a full regular-file byte comparison and two reported missing destinations; `TestArchiveCommandRefusesABrokenOutwardLink` passed with exit 2, stderr diagnostics and no move. `TestArchiveRestoresRewrittenBytesWhenRenameFails` passed for both PRD and another rewritten file. |
| Already archived target keeps its bytes | `TestArchiveKeepsALinkWhoseTargetWasAlreadyArchived` passed with identical Markdown bytes and a zero rewrite count. |
| Confirmation count appears only when links were rewritten | `TestArchiveCommandReportsRewrittenLinks` passed all four normal/QA-override and positive/zero-count cases. |

The declared Verification command remains for the Daemon. Task status,
Task Graph and other Task files were not edited. Delivery Queue resume
integration remains task_03's slice; no additional follow-up was identified.
