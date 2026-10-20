---
task: task_01
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
status: pending
type: backend
complexity: medium
---

# Task 01: The archive and cancel build repository paths from one root spelling

## Overview

`spec.Archive` and `spec.Cancel` relativize paths built from the configured
Spec Root against the repository root Git reported. When one spelling crosses
a symbolic link and the other does not, the result climbs out of the
repository: the archive's citation rewrite fails with the error the adopted
Inbox Entry quotes, the Archive Relink plans nothing, and the cancel refuses
its own Spec folder. This Task adds `RepositoryPath`, one resolver for
repository-relative paths, and makes both functions build every such path
with it. It is verifiable on its own through the package entry points.

## Requirements

1. MUST answer the Inbox Entry of 2026-10-20,
   [`references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md`](references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`
   in the new file `internal/spec/repository_path.go`, with the signature
   `_techspec.md` → Interfaces sketches, reusing
   `config.CanonicalizeResolvedPath` from `internal/config` rather than a new
   copy of the resolution.
2. MUST implement API Contract 2 in `internal/spec/archive.go`: the four
   paths handed to `PlanArchiveRelink` and the Legacy Link Rebase's
   `RecordPath` come from `RepositoryPath(req.RepositoryRoot, …)`; `source`
   keeps its computation and bytes; a resolver error returns wrapped with the
   operation that failed, before any file is written.
3. MUST implement API Contract 3 in `internal/spec/cancel.go`: the `relative`
   helper and `planSourceReturn` use `RepositoryPath(req.RepositoryRoot, …)`,
   and a path the resolver refuses, or that `cancelRelativePath` rejects,
   keeps the message `cancellation path "<path>" is outside the repository`.
4. MUST state the hazard in the code's behavior: the paths come from the
   configured Spec Root and from Git's top level of the checkout the caller
   runs in, never from the main checkout or the process work directory, and
   they go to Git and to the files the archive writes. A path that still lies
   outside the repository after both comparisons MUST be refused, never
   rewritten as `../…`, skipped, or passed on as an absolute path.
5. MUST keep everything else (API Contract 6): QA eligibility, the QA Archive
   Override, every refusal and its order, the Archive Record fields, the
   Archive Relink and Legacy Link Rebase rules, `CitationRewrite`, and every
   existing test name and assertion. No existing test, golden or fixture
   changes.
6. MUST add, in the new file `internal/spec/repository_path_test.go`, these
   tests, each top-level test calling `t.Parallel()` first, compiling for
   Windows, using temporary repositories and Git configuration, resolving its
   temporary directory first and creating its own symbolic link with
   `os.Symlink`:
   - `TestRepositoryPathComparesResolvedSpellings`, with the subtests
     `_techspec.md` → Testing Approach lists: one spelling; root resolved and
     path through the link; root through the link and path resolved; a path
     that does not exist yet under a linked parent; the root itself (`.`); a
     sibling sharing a name prefix, refused; a path outside, refused with a
     message naming both arguments. No result starts with `../` or is
     absolute (Invariant 1).
   - `TestArchiveFromALinkedSpellingRewritesAndRelinks`, through
     `spec.Archive` with a QA override, a `_prd.md` whose first paragraph
     links `../../backlog/<file>.md` and a committed Backlog Entry that links
     the Spec's `_prd.md`: with the Spec Root spelled through the link, and
     with the repository root spelled through the link, the archive succeeds
     and writes the same Archive Record bytes, the same relinked Backlog
     Entry bytes and the same `Relinked` list as a control run with one
     spelling and the same archive date (Invariant 3).
   - `TestCancelFromALinkedSpellingRelinksAndReturnsSources`, through
     `spec.Cancel` with an adopted Backlog Entry under `references/` and a
     committed Live Document citing the Spec: the same comparison with the
     control, and every `Returned` path repository-relative.

## Subtasks

- [ ] Add `RepositoryPath` over `config.CanonicalizeResolvedPath`.
- [ ] Build the archive's relink and record paths with it.
- [ ] Build the cancel's paths and returned sources with it, keeping its message.
- [ ] Add the resolver, archive and cancel tests.

## Acceptance Criteria

- [ ] `RepositoryPath` returns a lexical path when it can, a resolved one when the spellings differ, and a named refusal otherwise.
- [ ] `spec.Archive` and `spec.Cancel` given two spellings of one checkout write the bytes they write given one.
- [ ] Every existing archive, cancel and citation test passes unchanged.

## Consumers

- checkpoint: task_02

## Context

- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `internal/config/config.go`
- instruction: `internal/spec/legacy_link_rebase.go`
- instruction: `internal/spec/git_locator.go`
- instruction: `internal/spec/live_citations.go`
- instruction: `internal/spec/archive_test.go`
- instruction: `internal/spec/archive_links_test.go`
- instruction: `internal/spec/cancel_test.go`
- instruction: `internal/spec/git_locator_test.go`
- creates: `internal/spec/repository_path.go`
- creates: `internal/spec/repository_path_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/spec/cancel.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRepositoryPathComparesResolvedSpellings|TestArchiveFromALinkedSpellingRewritesAndRelinks|TestCancelFromALinkedSpellingRelinksAndReturnsSources|TestArchiveQAOverrideRecordsARelativeOutcome|TestArchiveRewritesRelativeLinksThatLeaveTheSpec|TestCancelRelinksLiveCitations|TestCancelRefusalsWriteNothing|TestCancelReturnsAnAdoptedBacklogEntryAsOpen|TestLegacyLinkRebaseNeverLinksTheRecordToItself|TestCitationRewriteKeepsRelinkContentBytes)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRepositoryPathComparesResolvedSpellings TestArchiveFromALinkedSpellingRewritesAndRelinks TestCancelFromALinkedSpellingRelinksAndReturnsSources TestArchiveQAOverrideRecordsARelativeOutcome TestCancelRelinksLiveCitations; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards the resolver, archive and cancel tests pass with the existing archive, cancel and citation tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-2; Success Metric 4
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-3; API Contract 6; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0247; ADR-0260; ADR-0278; ADR-0262; ADR-0154; ADR-0268; ADR-0272; ADR-0299; ADR-0259; ADR-0271; ADR-0285
