---
task: task_02
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
status: completed
type: backend
complexity: medium
---

# Task 02: An unregistered Run Worktree holding only ignored files is released, and empty directories go with it

## Overview

Fiscus also found a Run whose recorded Run Worktree directory Git no longer
registered, holding only a 46 MB test database volume under an ignored
`.storage/`: `roundfix reconcile` kept it `unknown` for ever, and an empty
repository directory under the Worktree Root stayed too
([`references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md`](references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md),
reported by fiscus on 2026-11-09, gap 2). This Task makes that directory
`safe` when Git's own ignore rules cover every file in it, lets `--apply`
remove it after a fresh re-listing, and makes a full-scan `--apply` remove
empty repository directories, as ADR-0379 decides.

## Requirements

1. MUST add to `internal/worktree/worktree.go` the reason constant
   `recorded Run Worktree is not registered in the Git root and holds only files the repository ignores`
   and, where `inspectTerminalRunMerged` returns `unknown` with
   `recorded Run Worktree is not registered in the Git root` because
   `recordedWorktree` reports the path unregistered, classify the Run `safe`
   with that reason instead only when all of these hold, and keep today's
   `unknown` and reason otherwise:
   - `filepath.Base` of the recorded path equals the Run ID;
   - the Run Branch is absent (`localBranchExists` returns false without an
     error) and no registered worktree holds it;
   - the directory has no `.git` entry of any type at its top;
   - walking it without following symbolic links succeeds, and every entry
     that is not a directory, listed as a slash-separated path relative to
     the directory, is printed by one
     `git check-ignore -z --stdin` run in the Git root with those paths
     NUL-separated on stdin (exit 0 with every path printed); exit 1 means
     not every path is ignored, and any other failure keeps `unknown`. A
     directory with no such entry qualifies without running Git. The
     command MUST NOT use `--no-index`.

   The `safe` result carries inspection evidence with `worktreePresent` and
   `runBranchPresent` false and the sorted relative path list.
2. MUST make `ApplyTerminalRun` release that result: the re-proof
   (`revalidateTerminalRunApply`) re-inspects and refuses, removing nothing,
   unless the fresh result is the same ignored-only `safe` case with an equal
   path list; then, skipping `ReconcileIntegration` because no Run Branch or
   head exists to record, check again that no `.git` entry exists at its top,
   make the tree owner-writable (`makeTreeOwnerWritable`) and remove the
   directory with `os.RemoveAll`. Every other result keeps its apply path
   byte for byte, including the Delivery Queue's release.
3. MUST add
   `func EmptyRepositoryDirectories(location string) ([]string, error)`,
   returning in byte order the absolute paths of the entries directly under
   `location` that are real directories (not symbolic links), whose name
   matches `^[a-z0-9]+(-[a-z0-9]+)*-[0-9a-f]{8}$`, and that hold no entry; an
   absent `location` returns none. Add
   `func RemoveEmptyRepositoryDirectory(location, path string) error`, which
   refuses a path that is not such a directory directly under `location` and
   removes it with `os.Remove`, so a directory that gained an entry is never
   removed.
4. MUST, in `internal/cli/reconcile.go`, on a full scan (no Run ID) only:
   list `EmptyRepositoryDirectories(runs.worktreeLocation)` in a new report
   field `EmptyDirectoryCandidates []reconcileDebrisResult` (JSON
   `emptyDirectoryCandidates`, `omitempty`), each with kind
   `emptyRepositoryDirectory`, the path in `worktree`, proof
   `repository directory under the Worktree Root holds nothing` and action
   `would remove with --apply`, or `remove after the Run releases` under
   `--apply`; and under `--apply`, after `applyReconcileWorktrees` and the
   other releases, list the directories again and remove each one with
   `RemoveEmptyRepositoryDirectory`, recording each removed directory with
   action `removed` and counting it in a new
   `reconcileDebrisSummary.EmptyDirectoriesRemoved` (JSON
   `emptyDirectoriesRemoved`, `omitempty`); a failure keeps the directory
   with action `preserve`, the error as its refusal reason, and counts an
   operational failure. The text report prints, before the `Summary:` line,
   for each entry `Empty repository directory: <path>` followed by
   `  action: <action>` and, when set, `  refusal-reason: <reason>`. A report
   with no such directory keeps every byte it prints after task_01.
5. MUST create `internal/worktree/unregistered_ignored_test.go`, with
   temporary repositories whose committed `.gitignore` holds `.storage/`,
   and a Run whose recorded path is a directory named after the Run under a
   repository directory below a temporary Worktree Root, with:
   - `TestAnUnregisteredRunWorktreeHoldingOnlyIgnoredFilesIsSafe`: the
     directory holding only `.storage/postgres/18/docker/base/1`, and an empty
     directory, are `safe` with the new reason; `ApplyTerminalRun` with
     `recordingTerminalRunStore` removes the directory, records nothing, and
     the next inspection reads `released`;
   - `TestAnUnregisteredRunWorktreeKeepsTodaysUnknown`, one subtest each for
     an added `notes.txt` the repository does not ignore, a `.git` file at
     the top, an existing Run Branch, a tracked path that a later
     `.gitignore` pattern also matches, and a directory not named after the
     Run: `unknown` with `recorded Run Worktree is not registered in the Git root`;
   - `TestApplyRefusesAnIgnoredOnlyWorktreeWhoseFilesChanged`: an ignored
     file added after inspection makes apply fail and the directory stays;
   - `TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames`:
     under one root, `fiscus-2f7556b8` and
     `fiscus-0035-implementation-20260908-392a816f` empty are listed and
     removed; `notes` empty, `repo-0123abcd` holding a file, and a symbolic
     link named `link-0123abcd` to an empty directory are neither listed nor
     removed.
6. MUST create `internal/cli/reconcile_unregistered_test.go`, with
   temporary homes and repositories only, where the Run's Run Worktree is
   created by `runworktree.Create`, unregistered with
   `git worktree remove --force`, its Run Branch deleted, and its directory
   recreated holding only `.storage/postgres/18/docker/base/1`, with:
   - `TestReconcileApplyReleasesAnIgnoredOnlyRunWorktreeAndEmptyRepositoryDirectories`:
     the dry run reads `classification: safe` with the new reason and lists
     an extra empty `other-0123abcd` under the Worktree Root; a full-scan
     `--apply` releases the Run, removes its emptied repository directory and
     `other-0123abcd`, keeps an empty `notes` directory, and reports
     `emptyDirectoriesRemoved` 2;
   - `TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory`: the same
     setup with a Run ID releases the Run, lists no empty directory and
     removes none.

   Every new top-level test calls `t.Parallel()` as its first statement, and
   every test file compiles for Windows.
7. MUST, in `docs/user-guide/commands/reconcile.md`, change the `unknown`
   row of the state table to say that an unregistered worktree stays
   `unknown` unless it holds only files the repository ignores, and add a
   paragraph after the Dirty Discard paragraph stating the conditions of
   Requirement 1, that `--apply` re-lists the files and removes the
   directory without a Git command and without a Run Database record, and
   the empty repository directory sweep of Requirements 3–4 with
   `emptyDirectoryCandidates`, `emptyDirectoriesRemoved` and the text lines.
8. MUST, through the `domain-modeling` skill, append to `CONTEXT.md` the two
   sentences the block after this list gives, byte for byte, each after a
   space: the first at the end of the **Run Worktree Reconciliation**
   definition line, after "only freshly revalidated `safe` or `superseded`
   work can be cleaned up."; the second at the end of the **Reconcile
   Command** definition line, after its ADR-0313 link sentence. No other line
   of `CONTEXT.md` changes; the **Run Worktree Reconciliation** entry keeps
   exactly its six state code spans.
9. MUST, as the last edit of the Task, after every guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST NOT change the Dirty Discard of task_01, any other classification or
    reason, any ADR, skill or other guide.


The `CONTEXT.md` sentences of Requirement 8, in order (Run Worktree
Reconciliation; Reconcile Command):

```markdown
A recorded Run Worktree that Git does not register is `safe` when its directory is named after the Run, its Run Branch is absent, it holds no Git metadata at its top and the repository's ignore rules cover every file in it ([ADR-0379](docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md)).

On a full scan, `--apply` also removes each empty repository directory directly under the Worktree Root ([ADR-0379](docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md)).
```

## Subtasks

- [ ] Classify the ignored-only unregistered directory and teach the apply path to remove it.
- [ ] Add the empty repository directory helpers and the full-scan sweep.
- [ ] Update the user guide and `CONTEXT.md`.
- [ ] Add the two test files and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A directory with only ignored files is `safe` and `--apply` removes it; any other unregistered directory stays `unknown`.
- [ ] A full-scan `--apply` removes empty repository-shaped directories under the Worktree Root and nothing else; a Run ID removes none.
- [ ] Reports without either case print the same bytes as after task_01.

## Context

- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/unregistered_ignored_test.go`
- interface: `internal/cli/reconcile.go`
- creates: `internal/cli/reconcile_unregistered_test.go`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/cli/reconcile_test.go`
- instruction: `docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Consumers

- `runworktree.InspectTerminalRunMerged` and `ApplyTerminalRun` — consumers: `./internal/worktree`, `./internal/daemon`, `./internal/cli`
- The `roundfix-reconcile/v1` report — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnUnregisteredRunWorktreeHoldingOnlyIgnoredFilesIsSafe|TestAnUnregisteredRunWorktreeKeepsTodaysUnknown|TestApplyRefusesAnIgnoredOnlyWorktreeWhoseFilesChanged|TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnUnregisteredRunWorktreeHoldingOnlyIgnoredFilesIsSafe TestAnUnregisteredRunWorktreeKeepsTodaysUnknown TestApplyRefusesAnIgnoredOnlyWorktreeWhoseFilesChanged TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/worktree && go test -count=1 -timeout 20m -run 'Branch|Disposition|Superseded' ./internal/daemon` — expected: exit 0. Before this Task none of the four tests exists; afterwards they pass with the whole package and the disposition tests that read the inspection.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReconcileApplyReleasesAnIgnoredOnlyRunWorktreeAndEmptyRepositoryDirectories|TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReconcileApplyReleasesAnIgnoredOnlyRunWorktreeAndEmptyRepositoryDirectories TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'Reconcile|DirtyDiscard|ReleaseMergedRuns|AutomaticRelease|TestRunCommandHelp' ./internal/cli` — expected: exit 0. Before this Task neither test exists; afterwards both pass with every reconcile, Dirty Discard, merged-release and help test.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'the repository.s ignore rules cover every file in it' || { echo 'missing Run Worktree Reconciliation change' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'also removes each empty repository directory directly under the Worktree Root' || { echo 'missing Reconcile Command change' >&2; exit 1; }; grep -F -- '**Run Worktree Reconciliation**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Run Worktree Reconciliation\*\*(:| —( |$))' || exit 1; grep -F -- '**Reconcile Command**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reconcile Command\*\*(:| —( |$))' || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- 'emptyDirectoriesRemoved' || { echo 'missing empty directory sweep in reconcile guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestGlossaryClosedListsEqualTheirConstants|TestGlossaryAvoidLinesArePlainLists|TestTheSkillCoverageMapIsCurrent)$'` — expected: exit 0. Before this Task the glossary and guide lack the new text; afterwards they hold it, the six-state closed list still matches and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Core Feature 4; Core Feature 5; User Experience; Success Metric 3; Success Metric 4; Glossary
- `_techspec.md` → Interfaces; Data Models; API Contract 3; API Contract 4; Testing Approach; Build Order 2
- ADR-0379; ADR-0053


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, no-workarounds, systematic-debugging.

Implemented the task_02 slice for Daemon verification:

- An unregistered recorded directory becomes safe only after the Run-name,
  absent-branch, registration and top-level Git metadata guards hold. A
  non-following walk collects sorted slash-separated paths, and one
  NUL-framed check-ignore invocation proves every non-directory entry against
  the repository index and checkout ignore rules, without --no-index. Empty
  trees need no check-ignore invocation. Every failed proof retains the old
  unknown reason.
- Apply re-inspects the same ignored-only case and exact path set, checks
  top-level metadata again, makes the tree owner-writable and removes it
  without ReconcileIntegration or a Git removal. Existing registered cleanup
  and Dirty Discard paths retain their behavior.
- Full scans report empty real repository-shaped directories; apply re-lists
  after the other releases and removes them with os.Remove. The sweep also
  runs when no Run is eligible. Run-ID invocations omit it. Failures report
  preserve, their reason and an operational failure. New JSON fields omit
  empty collections and zero counts, and text entries precede Summary.
- The reconcile guide documents both behaviors and the Run Database exception.
  CONTEXT.md changes only the two requested definition lines, appending the
  supplied sentences verbatim through domain-modeling.

Focused implementation evidence (all commands below exited 0 on this checkout):

```text
rtk proxy go test ./internal/worktree ./internal/cli -run 'TestAnUnregistered|TestApplyRefusesAnIgnoredOnly|TestEmptyRepositoryDirectories|TestIgnoredOnly|TestReconcileApplyReleasesAnIgnoredOnly|TestReconcileWithARunIDRemoves|TestReconcileWithoutADirtyDiscardKeepsItsWireBytes|TestReconcileDiscardDirtyNeedsItsDigestTranscript|TestRunReconcileIdempotentApply' -count=1 -v
rtk proxy go test ./internal/cli -run '^TestReconcileEmptyDirectory' -count=1 -v
rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestGlossary(ClosedListsEqualTheirConstants|LinksNameAcceptedADRs|PhrasesAreDefinedTerms|StatesEachRuleOnce)$' -count=1
```

| Acceptance criterion | Focused evidence |
| --- | --- |
| Ignored-only directories release; other unregistered directories remain unknown | The four required worktree tests passed. They cover ignored and empty trees, no Run Database record, a subsequent released inspection, unignored files, tracked paths matching a later ignore rule, Git metadata, an existing branch, a wrong directory name and a changed ignored path set. Additional tests passed for metadata directories and dangling links, preservation of a symlink target, whitespace in paths, and metadata/branch/ignore-rule changes between inspection and apply. |
| Full scan removes only empty repository-shaped directories; Run ID removes none | Both required CLI tests passed with temporary homes, real Git worktree creation/removal and recreated ignored storage. The full scan removed two repository directories and kept notes; Run-ID apply left both repository directories. The helper test passed for ordering, shape, nonempty directories, symlinks, absent roots and an entry added after listing. Additional CLI tests passed for a sweep without Run candidates and reported inspection failures. |
| Reports without either case retain their bytes | Existing TestReconcileWithoutADirtyDiscardKeepsItsWireBytes passed for text and JSON. New CLI assertions also passed for omission of emptyDirectoryCandidates and emptyDirectoriesRemoved when unused. The existing digest transcript and idempotent-apply regression tests passed. |

Windows compilation also exited 0, with both binaries written only into one
fresh temporary build directory:

```sh
task_build_dir=$(mktemp -d)
GOOS=windows rtk proxy go test -c -o "$task_build_dir/worktree.test.exe" ./internal/worktree && GOOS=windows rtk proxy go test -c -o "$task_build_dir/cli.test.exe" ./internal/cli
```

The initial red compilation identified the missing reason, ignored-path
inspection evidence and empty-directory APIs. Fixture setup was corrected to
create nested storage directories and resolve the temporary Git root's macOS
symlink spelling. The sandbox denied the initial Go cache access; the focused
checks above ran with authorized cache access. git diff --check exited 0.

The required Behavior Surface Record regeneration is the final edit after
this Result and all guide edits; its execution outcome is included in the
Agent handback:

```sh
go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
```

Task status and checkboxes remain Daemon-owned. No declared Verification
command, repository Verification or incremental Verification was run; the
focused checks above are not a Daemon verdict. No commit, push, pull request,
other task, Task Graph manifest, ADR or skill was changed. No follow-up work
was added to this slice.
