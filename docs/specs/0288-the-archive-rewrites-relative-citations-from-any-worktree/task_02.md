---
task: task_02
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
status: completed
type: backend
complexity: medium
---

# Task 02: Archive and cancel work from every spelling of the checkout

## Overview

The Archive Command and the Cancel Command take the repository root from `git
rev-parse --show-toplevel` and the Spec Root from the work directory, and
`filepathRelSlash` hands an absolute path on whenever the two disagree. From a
work directory that crosses a symbolic link, such as a Delivery Queue item
worktree under a linked `~/.roundfix/worktrees`, the cancel's citation advice
reads a path that does not exist and the refusals name absolute paths. This
Task makes every repository-relative path the two commands build come from
`RepositoryPath`, which task_01 adds, and proves the whole command from each
spelling of a linked worktree.

## Requirements

1. MUST answer the Inbox Entry of 2026-10-20,
   [`references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md`](references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md)
   (`_prd.md` → Sources 1), whose `roundfix archive --qa-override` failed at
   the root of an item worktree, by implementing API Contracts 4 and 5 of
   `_techspec.md`.
2. MUST, in `internal/cli/archive.go`, build with `RepositoryPath` against
   the root Git reported for that step: the Spec folder the uncommitted,
   ignored-file and removal-recovery messages name; the record path the
   record-only `--pull-request` form hands to `git ls-tree` and names in its
   refusal; and the four relink paths of `archive --plan`. The resolved-path
   special case for ignored files is removed. A path outside the repository
   is a Preflight failure naming it, never passed on as an absolute path.
3. MUST, in `internal/cli/cancel.go`, build with `RepositoryPath` against the
   root Git reported: the Spec folder its refusals and recovery name, and
   every path `cancelCitations` builds.
4. MUST state the hazard in the code's behavior: the root is Git's top level
   of the checkout the command runs in, read from the command's work
   directory, never the main checkout of a linked worktree and never a path
   the process inherited; the paths go to Git, to standard error and to the
   files the commands write.
5. MUST keep everything else (API Contract 6): the display paths built from
   the loaded Git root, every output line and its order, the Delivery Queue's
   archive step and `archivePaths`, `roundfix history sanitize`, the history
   conversion of `roundfix baseline update`, and every existing test name and
   assertion. No existing test, golden or fixture changes, and no skill or
   skill reference changes.
6. MUST document the change (`_prd.md` → Core Feature 4). In
   `docs/user-guide/commands/archive.md`, right after the sentence "Outward
   relative destinations are rebased to keep naming the same file.", add:
   "The archive reads every path relative to the checkout it runs in, so it
   writes the same record and relinks the same files from the main checkout,
   a linked worktree such as a Delivery Queue item worktree, a subdirectory
   of either, or a path that crosses a symbolic link." In
   `docs/user-guide/commands/cancel.md`, at the end of its first paragraph,
   add: "The cancel reads every path relative to the checkout it runs in, so
   it writes the same record from the main checkout, a linked worktree, a
   subdirectory of either, or a path that crosses a symbolic link." MUST
   search the guides, the agent guides and the skill copies for wording that
   says either command must run from the main checkout or the repository
   root, and report what it found; no skill text changes.
7. MUST record the Coverage Review `_prd.md` → Skills lists: in
   `docs/references/skill-coverage.json`, set the `review` of the surfaces
   `command archive`, `command cancel`, `guide
   docs/user-guide/commands/archive.md` and `guide
   docs/user-guide/commands/cancel.md` to `2026-10-20 — Spec 0288 task_02
   builds archive and cancel repository paths from one root spelling; no
   flag, help text, output line or exit code changes`, in the `YYYY-MM-DD —
   <reason>` form the map requires, changing nothing else in the file. Then
   re-record the Behavior Surface Record, whose two guide fingerprints change
   with the new sentences, by running
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back, because the Daemon never
   regenerates a record; the record is never hand-edited.
8. MUST add, in the new file `internal/cli/archive_symlinked_checkout_test.go`,
   these tests, each top-level test calling `t.Parallel()` first, compiling
   for Windows, using temporary repositories, homes and Git configuration,
   resolving its temporary directory first, creating its own symbolic link
   with `os.Symlink`, adding the linked worktree with `git worktree add`
   outside the main checkout, and running the command through the command
   entry point with the command environment's work directory set to the
   spelling under test:
   - `TestArchiveFromEverySpellingOfTheCheckoutTranscript`, with the subtests
     `_techspec.md` → Testing Approach lists: the worktree root resolved
     (control), the worktree root through the link, a worktree subdirectory
     through the link, and the main checkout through a link. Each runs
     `archive 0001-demo --qa-override --approval probe --reason probe` on a
     fresh fixture whose `_prd.md` first paragraph links
     `../../backlog/<file>.md` and whose committed Backlog Entry links the
     Spec's `_prd.md`, and each exits 0 with the control's standard output
     and standard error and writes the control's Archive Record and relinked
     Backlog Entry bytes. The worktree-root subtest asserts Surface
     Transcript 1 of `_techspec.md` whole: its one standard output line with
     only the declared variations, an empty standard error and exit code 0.
   - `TestCancelFromALinkedSpellingTranscript`: the control and the worktree
     root through the link, running `cancel 0001-demo --reason probe`; the
     linked run asserts Surface Transcript 2 whole and writes the control's
     record and relinked bytes.
   - `TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling`: from the
     worktree root through the link, an uncommitted change under the Spec
     folder gives the refusal naming `docs/specs/0001-demo`, and the
     record-only `archive 0001-demo --pull-request 7` for a record already on
     `origin/main` gives the refusal naming
     `docs/history/specs/0001-demo.md`; each transcript is byte-identical to
     the control's.
   - `TestArchivePlanWithALinkedExternalSpecRootListsTheRelink`: `archive
     0001-demo --plan --offline` with an external Spec Root configured
     through a link prints `relinks docs/backlog/<file>.md: 1 citation(s)`
     for the external repository's committed Backlog Entry that cites the
     Spec.

## Subtasks

- [ ] Build the Archive Command's refusal, record-only and plan paths with `RepositoryPath`.
- [ ] Build the Cancel Command's refusal, recovery and citation-advice paths with it.
- [ ] Remove the ignored-file special case the resolver covers.
- [ ] Add the transcript and refusal tests for every spelling.
- [ ] Add the sentence to the `archive` and `cancel` guides.
- [ ] Record the Coverage Review of the four surfaces and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] From a linked worktree's root or subdirectory spelled through a symbolic link, and from a linked main checkout, the archive prints and writes what it does from the resolved spelling.
- [ ] The cancel does the same from the linked worktree root.
- [ ] Refusals from a linked spelling name repository-relative paths.
- [ ] Every existing archive, cancel and folder-status test passes unchanged.
- [ ] The `archive` and `cancel` guides say each command reads every path relative to the checkout it runs in.

## Consumers

- The Archive Command and the Cancel Command (`roundfix archive`, `roundfix cancel`) — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`
- instruction: `.agents/skills/roundfix/references/archive.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `internal/cli/spec_folder_status.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `internal/cli/archive_relink_test.go`
- instruction: `internal/cli/spec_folder_status_test.go`
- instruction: `internal/cli/cancel_test.go`
- instruction: `internal/cli/archive_pull_request_test.go`
- instruction: `internal/cli/archive_plan_offline_test.go`
- creates: `internal/cli/archive_symlinked_checkout_test.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/cancel.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/cancel.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveFromEverySpellingOfTheCheckoutTranscript|TestCancelFromALinkedSpellingTranscript|TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling|TestArchivePlanWithALinkedExternalSpecRootListsTheRelink|TestArchiveRelinksLiveCitationsTranscript|TestSpecFolderStatusRefusesUncommittedChangesFirst|TestSpecFolderStatusReadsTheExternalSpecsRootRepository|TestArchiveRefusesIgnoredSpecFiles|TestCancelRefusesIgnoredSpecFiles|TestCancelTranscriptCancelsASpecWithSources|TestCancelWarnsOnlyOncePerRemainingTrackedCitationLine|TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch|TestArchivePlanOfflineTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveFromEverySpellingOfTheCheckoutTranscript TestCancelFromALinkedSpellingTranscript TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling TestArchivePlanWithALinkedExternalSpecRootListsTheRelink TestArchiveRelinksLiveCitationsTranscript TestSpecFolderStatusReadsTheExternalSpecsRootRepository; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/references/skill-coverage.json | grep -qF -- 'Spec 0288 task_02 builds archive and cancel repository paths from one root spelling' || { printf 'missing Coverage Review\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -qF -- 'The archive reads every path relative to the checkout it runs in' || { printf 'missing archive guide sentence\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/cancel.md | grep -qF -- 'The cancel reads every path relative to the checkout it runs in' || { printf 'missing cancel guide sentence\n' >&2; exit 1; }; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' || exit 1; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the four new tests do not exist, so the command fails at its first pass check. Afterwards the commands work from every spelling with the existing archive, cancel, folder-status, pull-request and plan tests, both guides carry their sentence, the Coverage Review and the Behavior Surface Record are current, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 2-4; Success Metric 1; Success Metric 2; Success Metric 3; Skills
- `_techspec.md` → Current behavior; API Contracts 4-6; Surface Transcript 1; Surface Transcript 2; Invariants 2-4; Testing Approach; Build Order 2
- ADR-0247; ADR-0260; ADR-0262; ADR-0154; ADR-0269; ADR-0184; ADR-0280; ADR-0256; ADR-0259; ADR-0271; ADR-0285

## Result

Implemented the task_02 command-path slice for Daemon Verification. Task status
remains Daemon-owned; Subtasks and Acceptance Criteria checkboxes are unchanged.

- Archive now uses `spec.RepositoryPath` for its Spec-folder refusal/recovery
  path, record-only pull-request lookup/refusal path and four plan relink paths.
  Plan resolves Git's checkout top level for the built-in root as well as an
  external root. Cancel uses the resolver for its folder refusal/recovery path
  and every path built by `cancelCitations`. Both ignored-file spelling special
  cases are removed. Resolver errors propagate as Preflight failures.
- Code comments identify the root's provenance: the checkout reached from the
  command's work directory, or the external Spec repository, rather than the
  main worktree or inherited process directory. They identify Git, diagnostics
  and written archive files as path consumers. Loaded-root display formatting,
  output lines/order, Delivery Queue archive code and history conversion are
  unchanged.
- Added only `internal/cli/archive_symlinked_checkout_test.go`. Its four
  top-level tests begin with `t.Parallel()`, resolve temporary directories before
  creating links, use disposable repositories/homes/Git configuration and add
  linked worktrees outside the main checkout. Controls use the same committed
  revision, allowing literal transcript, record and Backlog Entry comparisons.
  Comparison controls are created outside subtests so filtered subtest runs do
  not depend on another subtest executing.

Acceptance evidence:

| Criterion | Focused evidence |
| --- | --- |
| Archive from resolved/linked worktree roots, linked subdirectory and linked main checkout | `TestArchiveFromEverySpellingOfTheCheckoutTranscript` passes; every spelling asserts the whole successful transcript and exact control record/relinked bytes, including the outward PRD link. |
| Cancel from linked worktree root | `TestCancelFromALinkedSpellingTranscript` passes; control and linked spelling assert the whole transcript, exact record/relinked bytes and removal. |
| Repository-relative refusals | `TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling` passes for dirty Spec and default-branch record-only refusals, with byte-identical control transcripts and unchanged checkout contents. |
| Existing archive/cancel/folder-status behavior | Focused `Test(Archive\|Cancel\|SpecFolderStatus)` selection passes without changing any existing test, fixture, golden or assertion. The new external-root offline-plan test also passes and observes the committed Backlog Entry's relink line without modifying files. |
| Guide coverage | Added both prescribed sentences in their required positions. The Coverage Review changes exactly the four requested surfaces; sanctioned regeneration changes only the two guide fingerprints. |

Focused checks (all executed through `rtk proxy`):

- Before the command fix, the linked cancel and external-root plan tests
  reproduced the invalid concatenated absolute citation-target path; the dirty
  linked archive refusal reproduced its absolute Spec-folder path. Initial new
  fixture errors (missing QA dependency, invalid dirty PRD and absent empty
  Spec Root after archival) were corrected within the new test file.
- `go test ./internal/cli -run 'Test(ArchiveFromEverySpelling|CancelFromALinkedSpelling|ArchiveRefusalsNameRepositoryPaths|ArchivePlanWithALinkedExternal)' -count=1`
  — exit 0 after the final test edit; all four new tests pass.
- `go test ./internal/cli -run 'Test(Archive|Cancel|SpecFolderStatus)' -count=1`
  — exit 0. An earlier run reported PASS but exited 1 because concurrent
  sanctioned record generation changed `behavior-surfaces.json` during the
  suite guard. The successful run occurred after regeneration finished.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  — exit 0; ran the required generator, with only the archive/cancel guide
  fingerprints changing in the generated record.
- `env GOOS=windows go test -c -o /private/tmp/roundfix-0288-task02-cli.test.exe ./internal/cli`
  — exit 0 after the final test edit; cross-compiles the package and tests.
- `git diff --check` — exit 0. Changed-file audit finds only this Task's
  declared files; the Task file already had a Daemon-owned status change on
  entry. No other Task, Task Graph, skill or existing test changed.

Wording review: searched `docs/user-guide`, `docs/agents`, `.agents/skills` and
`skills` for `main checkout`, `repository root` and `repo root`, and for either
command near those phrases. No wording requires the archive or cancel command
itself to run from the main checkout or repository root. Both copies of
`archive-spec/SKILL.md:155` say to run their reference-index audit shell snippet
from the repository root. Both copies of
`roundfix/references/spec-delivery.md:51` and the usage guide's
`IntegrationPending` row require the printed integration-recovery command from
that root. Other matches describe configuration resolution, Bun scripts or
repository conditions, including both archive references' "When the repository
root is inside a Git work tree". These statements govern different operations
or describe roots; no skill text changed.

Environment: the executor rejected the symlinked writable-root spelling before
starting a process. Authorized commands ran in the resolved Run checkout under
`/Volumes/Brain-SSD/cache/.roundfix/worktrees/`; no configuration workaround was
written.

The declared Verification, Windows vet and repository Verification remain for
the Daemon. No commit, push or pull request was made; no follow-up scope was
implemented.
