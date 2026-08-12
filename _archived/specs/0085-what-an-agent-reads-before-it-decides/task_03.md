---
task: task_03
spec: 0085-what-an-agent-reads-before-it-decides
status: completed
type: backend
complexity: high
---

# Task 03: Move every consumer onto the resolver

## Overview

Every package that composes an archive path stops doing so and asks the resolver
instead, and the two hardcoded literals in the Spec checker are deleted. Nothing
moves on disk in this Task — the resolver still answers today's paths — so this
is the prefactoring that makes Task 04's relocation a one-line change rather than
five.

## Requirements

1. MUST replace every self-composed archive path in the Spec checker, the Spec
   audit, the worktree QA evidence path, and the CLI with a resolver call.
2. MUST delete the two hardcoded archive literals in the Spec checker.
3. MUST leave every observable path identical, because the resolver still
   answers the locations Task 01 recorded.
4. MUST leave no package other than the resolver's own expressing the archive
   layout.

## Subtasks

- [ ] Move the Spec checker onto the resolver and delete its two literals.
- [ ] Move the Spec audit, worktree, and CLI onto the resolver.
- [ ] Prove no other package expresses the layout.

## Acceptance Criteria

- [ ] Every consumer resolves its archive path through the one owner.
- [ ] The two literals are gone.
- [ ] Every observable path is unchanged.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/citations.go`
- `internal/speccheck/citations_test.go`
- `internal/speccheck/backlog.go`
- `internal/speccheck/backlog_test.go`
- `internal/specaudit/audit.go`
- `internal/specaudit/audit_test.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `internal/cli/archive.go`
- `internal/cli/archive_test.go`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_03.md`

This list was corrected on 2026-08-11. It first named `internal/speccheck/speccheck.go`
and `internal/specaudit/specaudit.go`, which do not exist: the archive-layout
expressions live in `citations.go`, `backlog.go` and `audit.go`. The Agent
refused the Task rather than migrate part of it, and named the missing files —
the right call, because a partial migration leaves the one-owner criterion false
while looking like progress. Enumerating a boundary from guessed filenames is
what cost this Run.

## Verification

- `test -z "$(grep -rnE '"_archived"|_archived/' internal/speccheck/*.go internal/specaudit/*.go internal/worktree/*.go internal/cli/archive.go | grep -v '_test.go' | grep -vE ':[0-9]+:[[:space:]]*//')"` — expected: exits 0. Comment lines are excluded because a doc comment cannot call a resolver; every other occurrence, including the text of a `Fix:` message that tells a reader where an artifact belongs, must come from the resolver, since a message naming a layout goes stale the moment Task 04 moves it.

Asserting the layout characterization still passes is deliberately absent: this
Task changes no observable path, so the case passes before and after, and a
command that cannot fail cannot prove anything. Keeping it passing is the
Run-level gate's job.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 3; System Architecture.

## Result

Moved the Spec checker, Spec audit, worktree QA evidence lookup, and Archive
Command help onto `spec.ArchiveDir`. The checker now builds its archived-spec
and archived-finding reads and repair text from the resolver. The audit keeps
today's configured external Spec Root behavior while accepting a future
resolver answer outside the active Spec tree without another layout literal.

- Criterion 1: a focused `ArchiveDir(` call-site sweep found resolver calls in
  `internal/speccheck`, `internal/specaudit`, `internal/worktree`, and
  `internal/cli/archive.go`; the audit's default and configured-external-root
  cases both passed.
- Criterion 2: a focused source sweep for `docs/specs/_archived`,
  `docs/findings/_archived`, and the `_archived` string literal reported no
  matches in the five changed production files.
- Criterion 3: the focused consumer tests preserved the current archive move,
  CLI output and help, checker findings, audit artifact paths, and worktree QA
  evidence paths.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-spec0085-task03-gocache go test
  ./internal/speccheck ./internal/specaudit ./internal/worktree ./internal/cli
  -run 'Test(Check(RollupMember|ArchiveLicense|BacklogUnmoved)|Audit(ReportsUndeliveredArchiveHeldByBranch|UsesConfiguredExternalSpecRootTree)|QAReportOnlyBranch|InspectTerminalRunClassifiesSupersededQAReport|RunArchive(MovesCompletedSpecAndStampsMetadata|UsesConfiguredExternalSpecRoot|Help))$'
  -count=1` passed in all four packages.
- `rtk git diff --check` reported no whitespace errors.

The Daemon-owned `## Verification` command was not run.
