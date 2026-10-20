---
task: task_02
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
status: pending
type: backend
complexity: medium
---

# Task 02: Archive and cancel refuse it through the command and the guides say so

## Overview

Every path the Archive Command and the Cancel Command build already goes
through `RepositoryPath`, which task_01 hardens, so the commands need no code
change; what is missing is proof through the command entry point that each
refusal reaches the operator whole, writes nothing, and that every link that
stays inside, a checkout reached through a link and an external Spec Root
keep working. This Task adds those transcript tests and states the refusal in
the `archive` and `cancel` user guides.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md`](references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md)
   (`_prd.md` → Sources 1), by proving API Contract 4 of `_techspec.md`
   through `roundfix archive` and `roundfix cancel`.
2. MUST add, in the new file `internal/cli/archive_link_escape_test.go`, the
   three tests `_techspec.md` → Testing Approach names, each top-level test
   calling `t.Parallel()` first, compiling for Windows, using temporary
   repositories, homes and Git configuration, resolving its temporary
   directory first, creating its own links with `os.Symlink`, and running the
   command through the command entry point with the command environment's
   work directory set to the spelling under test:
   - `TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript`: the
     `docs/history` case asserts Surface Transcript 1 whole (every standard
     error line in order, an empty standard output and exit `2`, with `<repo>`
     the only variation); its subtests assert the whole refusal transcript for
     the Spec folder link, the promotion into a linked `docs/references`, and
     `archive 0001-demo --plan --offline` with a linked `docs/history`.
   - `TestCancelRefusesALinkThatLeavesTheRepositoryTranscript`: the
     `docs/backlog` case asserts Surface Transcript 2 whole; a subtest asserts
     the `docs/history` case refused by the citation advice.
   - `TestArchiveKeepsLinksThatStayInsideTranscript`: a `docs/history` link
     inside the repository, the checkout reached through a link, and an
     external Spec Root configured through a link spelling of its repository
     each exit 0 with the control's standard output and standard error and
     write the control's record bytes.
   Every refusal case also asserts that the outside directory holds no new
   file and that the Spec folder bytes and `git status --porcelain
   --untracked-files=all` are unchanged.
3. MUST NOT change `internal/cli/archive.go` or `internal/cli/cancel.go`
   unless a transcript cannot pass otherwise; if it must, the change keeps
   every existing output line and refusal order (API Contract 5) and the
   Task's Result names why.
4. MUST document the change (`_prd.md` → Core Feature 4). In
   `docs/user-guide/commands/archive.md`, right after the sentence ending
   "or a path that crosses a symbolic link.", add: "A path whose resolved
   target lies outside the repository, such as a `docs/history` or
   `docs/references` symbolic link to another directory, is refused with exit
   `2` before anything is written, and the refusal names the path. A symbolic
   link whose target stays inside the repository, a checkout reached through
   a symbolic link and an external Spec Root keep working." In
   `docs/user-guide/commands/cancel.md`, right after its sentence ending "or a
   path that crosses a symbolic link.", add: "A path whose resolved target
   lies outside the repository, such as a `docs/history` or `docs/backlog`
   symbolic link to another directory, is refused with exit `2` before
   anything is written, and the refusal names the path. A symbolic link whose
   target stays inside the repository, a checkout reached through a symbolic
   link and an external Spec Root keep working." MUST search the guides, the
   agent guides and the skill copies for wording that says either command
   accepts or follows a symbolic link whose target lies outside the
   repository, and report what it found in the Result.
5. MUST re-record the Behavior Surface Record, whose two guide fingerprints
   change with the new paragraphs, by running
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back, because the Daemon never
   regenerates a record; the record is never hand-edited, and
   `docs/references/skill-coverage.json` does not change.

## Subtasks

- [ ] Add the archive refusal, cancel refusal and inside transcript tests.
- [ ] Add the paragraph to the `archive` and `cancel` guides.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Each outside case exits 2 with its whole refusal transcript and writes nothing.
- [ ] Each inside, linked-checkout and external Spec Root case archives with the control's output and record bytes.
- [ ] Both guides state the refusal and what keeps working, and the Behavior Surface Record is current.

## Consumers

- The Archive Command and the Cancel Command (`roundfix archive`, `roundfix cancel`) — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md`
- instruction: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- instruction: `internal/cli/archive.go`
- instruction: `internal/cli/cancel.go`
- instruction: `internal/cli/spec_folder_status.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `internal/cli/archive_symlinked_checkout_test.go`
- instruction: `internal/cli/spec_folder_status_test.go`
- instruction: `internal/cli/cancel_test.go`
- instruction: `internal/cli/archive_plan_offline_test.go`
- instruction: `internal/cli/archive_relink_test.go`
- creates: `internal/cli/archive_link_escape_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/cancel.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript|TestCancelRefusesALinkThatLeavesTheRepositoryTranscript|TestArchiveKeepsLinksThatStayInsideTranscript|TestArchiveFromEverySpellingOfTheCheckoutTranscript|TestCancelFromALinkedSpellingTranscript|TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling|TestArchivePlanWithALinkedExternalSpecRootListsTheRelink|TestSpecFolderStatusReadsTheExternalSpecsRootRepository|TestCancelWriteFailureUsesExitOne|TestCancelRefusesOccupiedSourceDestination|TestCancelTranscriptCancelsASpecWithSources|TestArchivePlanOfflineTranscript|TestArchiveRelinksLiveCitationsTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript TestCancelRefusesALinkThatLeavesTheRepositoryTranscript TestArchiveKeepsLinksThatStayInsideTranscript TestArchiveFromEverySpellingOfTheCheckoutTranscript TestCancelFromALinkedSpellingTranscript TestArchivePlanWithALinkedExternalSpecRootListsTheRelink TestCancelWriteFailureUsesExitOne; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -q -- 'such as a .docs/history. or .docs/references. symbolic link to another directory, is refused with exit' || { printf 'missing archive guide paragraph\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/cancel.md | grep -q -- 'such as a .docs/history. or .docs/backlog. symbolic link to another directory, is refused with exit' || { printf 'missing cancel guide paragraph\n' >&2; exit 1; }; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' || exit 1; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards both commands refuse every outside case and keep every inside case through the command entry point, the existing symlinked-checkout, external-root, cancel, plan and relink tests pass, both guides carry their paragraph, the Behavior Surface Record is current, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 3-4; Declared break; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior; API Contracts 4-5; Surface Transcript 1; Surface Transcript 2; Invariants 2-4; Testing Approach; Build Order 2
- ADR-0306; ADR-0269; ADR-0262; ADR-0184; ADR-0280; ADR-0259; ADR-0271; ADR-0285
