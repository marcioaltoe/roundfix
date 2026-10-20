---
task: task_01
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
status: pending
type: backend
complexity: medium
---

# Task 01: The resolver refuses a link that leaves the repository

## Overview

`RepositoryPath` accepts a path that is lexically inside the repository
before it resolves anything, so a symbolic link inside the repository whose
target lies outside is accepted, and `spec.Archive` and `spec.Cancel` write
through it. Their promotion and returned-source destinations never reach the
resolver at all. This Task adds the containment test to the resolver, routes
the three destination joins through it, and proves both through the package
entry points.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md`](references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`
   in `internal/spec/repository_path.go`, with the signature `_techspec.md` →
   Interfaces sketches: accept a path only when
   `config.CanonicalizeResolvedPath(path)` lies inside or equals
   `config.CanonicalizeResolvedPath(repositoryRoot)`; return the lexical
   relative path when it stays inside, the resolved one otherwise; refuse
   with today's message `path "<path>" is outside the repository
   "<repositoryRoot>"`. Both sides MUST be resolved: resolving only the path
   would refuse a checkout reached through a symbolic link.
2. MUST implement API Contract 2 in `internal/spec/archive.go`: before the
   vacancy check of each promotion destination, call `RepositoryPath` and
   return `resolve promotion destination: <resolver error>` before any write.
3. MUST implement API Contract 3 in `internal/spec/cancel.go`: before the
   vacancy check of each promotion destination and each returned-source
   destination in `planSourceReturn`, call `RepositoryPath` and return
   `cancellation path "<destination>" is outside the repository` before any
   write.
4. MUST state the hazard in the code's behavior: the paths come from the
   configured Spec Root and from Git's top level of the checkout the caller
   runs in; the destinations are where the commands write the Archive Record,
   promoted copies and returned sources, which a commit of that repository
   must hold. A path whose resolved target lies outside MUST be refused
   before anything is written, never followed.
5. MUST keep everything else (API Contract 5): the relative path of every
   path accepted today whose target stays inside, every existing refusal and
   its order, the record fields, and every existing test name and assertion.
   No existing test, golden or fixture changes.
6. MUST add, in the new file
   `internal/spec/repository_path_containment_test.go`, the five tests
   `_techspec.md` → Testing Approach names, with their subtests, each
   top-level test calling `t.Parallel()` first, compiling for Windows, using
   temporary repositories and Git configuration, resolving its temporary
   directory first and creating its own links with `os.Symlink`:
   - `TestRepositoryPathRefusesALinkThatLeavesTheRepository`;
   - `TestRepositoryPathKeepsLinksThatStayInside`, including the subtest
     where the root and the path are both spelled through a link to the
     checkout and the lexical relative path is returned;
   - `TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting`;
   - `TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting`;
   - `TestArchiveAndCancelKeepLinksThatStayInside`.
   Each refusal subtest asserts the exact refusal text, that the outside
   directory holds no new file, and that the Spec folder bytes and `git
   status --porcelain --untracked-files=all` are unchanged.

## Subtasks

- [ ] Add the containment test to `RepositoryPath`, resolving both sides.
- [ ] Route the archive's promotion destination through it.
- [ ] Route the cancel's promotion and returned-source destinations through it.
- [ ] Add the five containment tests.

## Acceptance Criteria

- [ ] `RepositoryPath` refuses a path whose resolved target lies outside, and returns today's relative path for every path whose target stays inside.
- [ ] `spec.Archive` and `spec.Cancel` refuse every outside destination before writing anything.
- [ ] Every existing archive, cancel and resolver test passes unchanged.

## Consumers

- checkpoint: task_02

## Context

- instruction: `docs/adr/0306-archive-and-cancel-refuse-a-path-whose-resolved-target-leaves-the-repository.md`
- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `internal/config/config.go`
- instruction: `internal/spec/repository_path_test.go`
- instruction: `internal/spec/archive_test.go`
- instruction: `internal/spec/archive_record_test.go`
- instruction: `internal/spec/archive_links_test.go`
- instruction: `internal/spec/cancel_test.go`
- creates: `internal/spec/repository_path_containment_test.go`
- interface: `internal/spec/repository_path.go`
- interface: `internal/spec/archive.go`
- interface: `internal/spec/cancel.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRepositoryPathRefusesALinkThatLeavesTheRepository|TestRepositoryPathKeepsLinksThatStayInside|TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting|TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting|TestArchiveAndCancelKeepLinksThatStayInside|TestRepositoryPathComparesResolvedSpellings|TestArchiveFromALinkedSpellingRewritesAndRelinks|TestCancelFromALinkedSpellingRelinksAndReturnsSources|TestArchiveRollsBackTheRecordWhenPromotionCannotBeWritten|TestArchiveRewritesRelativeLinksThatLeaveTheSpec|TestCancelReturnsAnAdoptedBacklogEntryAsOpen|TestCancelRefusalsWriteNothing|TestCancelRelinksLiveCitations)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRepositoryPathRefusesALinkThatLeavesTheRepository TestRepositoryPathKeepsLinksThatStayInside TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting TestArchiveAndCancelKeepLinksThatStayInside TestRepositoryPathComparesResolvedSpellings TestArchiveFromALinkedSpellingRewritesAndRelinks TestCancelFromALinkedSpellingRelinksAndReturnsSources TestArchiveRollsBackTheRecordWhenPromotionCannotBeWritten; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec` — expected: exit 0. Before this Task the five new tests do not exist, so the command fails at its first pass check. Afterwards the containment tests pass with the existing resolver, archive, cancel and promotion tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Declared break; Success Metric 4
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-3; API Contract 5; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0306; ADR-0268; ADR-0272; ADR-0247; ADR-0260; ADR-0278; ADR-0262; ADR-0154; ADR-0299; ADR-0259; ADR-0271; ADR-0285
