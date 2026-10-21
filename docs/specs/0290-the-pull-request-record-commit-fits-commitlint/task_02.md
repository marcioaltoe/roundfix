---
task: task_02
spec: 0290-the-pull-request-record-commit-fits-commitlint
status: completed
type: backend
complexity: medium
---

# Task 02: The record subject names the Spec number, and a refused queue commit names its hook

## Overview

The Delivery Queue's Pull Request Record Commit subject carries the whole
slug, so a long slug passes 100 characters and a commitlint `commit-msg` hook
refuses it; the queue then parks a generic `delivery-error` and a retry
repeats the refusal. This Task builds the subject from the Spec number, pins
that a record commit made by hand under any subject is still adopted, and
makes the three commits the queue makes report a hook refusal with the hook,
the subject and the recovery.

## Requirements

1. MUST answer the Inbox Entry of 2026-10-20,
   [`references/2026-10-20-o-commit-do-registro-do-pull-request-estoura-o-commitlint-com-slug-longo.md`](references/2026-10-20-o-commit-do-registro-do-pull-request-estoura-o-commitlint-com-slug-longo.md)
   (`_prd.md` → Sources 1), whose record commit the oraculum `commit-msg`
   hook refused for a 107-character subject, by implementing API Contracts 1-6
   of `_techspec.md`.
2. MUST add `pullRequestRecordSubject(specSlug, number string) string` in
   `internal/cli/deliver_workflow.go` (API Contract 1): `docs: record pull
   request #<number> in the <digits> Archive Record` when the slug starts with
   one to ten ASCII digits followed by `-`, else `docs: record pull request
   #<number> in the Archive Record`; and make `RecordPullRequest` commit with
   it (API Contract 2). The record proof functions
   `pullRequestRecordCommitIsExact`, `pullRequestRecordTreeIsExact` and
   `pullRequestRecordChangeIsExact` keep their bytes.
3. MUST route the three queue commits, in `Archive`, `RecordPullRequest` and
   `ResolveConflict`, through one helper that implements API Contracts 3 and 4.
   The helper classifies a failure only when the error wraps an
   `*exec.ExitError` and `daemon.ClassifyCommitHookRefusal` reports a refusal
   for the hook output read from the runner's error with its `git <args>: `
   prefix and `: <exit error>` suffix removed; it never classifies on the
   arguments, the subject or the slug. The refusal error's text, its hook
   label, the 512-byte cut of the output and both recovery sentences are
   exactly those of API Contract 4, and the error unwraps to the runner's
   error. Every caller keeps its existing wrapping. Any other failure returns
   today's error unchanged.
4. MUST state the hazard in the code's behavior: the hook's output is
   third-party text from the repository's own hook; it goes, bounded and on
   one line, into the step's error and so into the Delivery Queue item's
   blocker, which `deliver status` prints and the Run Database stores. The
   helper never passes `--no-verify`, never retries with another subject and
   never rewrites the commit.
5. MUST export `InstalledCommitHooks(ctx context.Context, workDir string)
   []string` from `internal/daemon/commit_hook.go` (API Contract 5), built
   from the logic of `commitHookInstalled`, which becomes `len(...) > 0`, so
   `ClassifyCommitHookRefusal` and the Daemon's Task commit refusal keep
   their behavior. No file under `internal/daemon` other than
   `commit_hook.go` changes.
6. MUST keep API Contract 6: no change to `internal/delivery`, to the blocker
   constants, to `ClassifyPark`, to the stage an item parks at, to `deliver
   status` or to the Pull Request body.
7. MUST declare the break: in `internal/cli/deliver_record_pull_request_test.go`,
   `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord` expects the
   subject `docs: record pull request #7 in the 0001 Archive Record`; nothing
   else in that file changes.
8. MUST add, in the new file `internal/cli/deliver_commit_hook_test.go`, the
   tests `_techspec.md` → Testing Approach lists, each top-level test calling
   `t.Parallel()` first, compiling for Windows, using temporary repositories,
   homes and Git configuration, writing its hooks as executable shell scripts
   with `writeScriptFixture`, and waiting on no real time:
   - `TestPullRequestRecordSubjectNamesTheSpecNumber`, with the six subtests
     listed there, asserting the 70-character maximum.
   - `TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook`: a Spec slug of
     54 characters starting with `0088-` (a synthetic slug, not oraculum's),
     archived on an item branch, with a `commit-msg` hook that exits 1 with
     `header must not be longer than 100 characters` on standard error when
     the message's first line exceeds 100 characters;
     `RecordPullRequest(…, "280")` returns an exact head whose subject is
     `docs: record pull request #280 in the 0088 Archive Record`.
   - `TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject`, with the
     subtests old long subject, short subject and `wip`.
   - `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery`, with the
     five subtests listed there, asserting the error text of API Contract 4,
     `errors.As` to `*exec.ExitError`, the `environment` Park Class through
     `delivery.ClassifyPark`, and, for the record subtests, adoption of the
     staged change committed by hand with a conforming subject. The archive
     commit subtest uses a slug that contains `pre-commit`, so the refused
     subject `docs: archive <slug>` names it, while only a silent
     `commit-msg` hook is installed, and asserts the hook label
     `commit-msg`.
   - `TestDeliveryQueueCommitWithoutAHookKeepsTheGitError`: no hook, a stale
     `.git/index.lock`, and an error that starts with `commit pull request
     record: git commit -m ` and contains no `hook refused`.

## Subtasks

- [ ] Build the record subject from the Spec number.
- [ ] Route the three queue commits through one helper that names a hook refusal.
- [ ] Export the installed commit hooks from the Daemon package without changing its behavior.
- [ ] Add the subject, adoption, refusal and plain-failure tests, and update the one subject assertion.

## Acceptance Criteria

- [ ] Behind a 100-character header hook, the record commit of a 54-character slug succeeds with a 57-character subject.
- [ ] A record commit made by hand under any subject is adopted.
- [ ] A hook refusal of any queue commit names the hook, the subject, the exit status, the output and the recovery, and parks with class `environment`.
- [ ] A failure without a hook keeps today's error, and every daemon hook test passes unchanged.

## Consumers

- The record subject, the queue commit helper and its refusal error (`RecordPullRequest`, `Archive`, `ResolveConflict`) — consumers: `./internal/cli`
- `daemon.InstalledCommitHooks` — consumers: `./internal/cli`, `./internal/daemon`

## Context

- instruction: `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/delivery/engine.go`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/preflight/preflight.go`
- instruction: `internal/daemon/daemon.go`
- instruction: `internal/daemon/commit_hook_test.go`
- instruction: `internal/cli/script_fixture_test.go`
- instruction: `internal/cli/deliver_archive_record_test.go`
- instruction: `internal/cli/deliver_conflict_test.go`
- instruction: `internal/cli/deliver_operator_archive_test.go`
- instruction: `internal/cli/settle_test.go`
- creates: `internal/cli/deliver_commit_hook_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_record_pull_request_test.go`
- interface: `internal/daemon/commit_hook.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPullRequestRecordSubjectNamesTheSpecNumber|TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook|TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject|TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery|TestDeliveryQueueCommitWithoutAHookKeepsTheGitError|TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord|TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath|TestDeliveryWorkflowAdoptsAnUncommittedRecord|TestDeliveryWorkflowRefusesAnUncommittedInexactRecord|TestADerivedConflictIsMergedAndRegenerated|TestConflictRecoveryAbortsAFailedRegeneration)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPullRequestRecordSubjectNamesTheSpecNumber TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery TestDeliveryQueueCommitWithoutAHookKeepsTheGitError TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^(TestClassifyCommitHookRefusal|TestClassifyCommitHookRefusalReadsTheRepositoryWhenOutputNamesNothing|TestClassifyCommitHookRefusalResolvesTheHooksDirectoryGitWouldRun|TestGitCommitterClassifiesAHookThatNamesNothing|TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged|TestGitCommitterKeepsPlainGitFailureUnclassified|TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery)$' ./internal/daemon || exit 1; GOOS=windows go vet ./internal/cli ./internal/daemon` — expected: exit 0. Before this Task the five new tests do not exist and the existing subject test expects the slug form, so the command fails at its first pass check. Afterwards the subject, adoption, refusal and plain-failure tests pass with the existing record and conflict tests, the Daemon's hook tests pass unchanged, and both packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 2
- ADR-0302

## Result

Implemented this Task's slice; Task status and declared Verification remain
Daemon-owned.

- Record subjects now use a one-to-ten-digit ASCII Spec prefix, or omit the
  identifier when no valid prefix exists. The existing subject assertion now
  expects `docs: record pull request #7 in the 0001 Archive Record` and is the
  only change in `deliver_record_pull_request_test.go`.
- Archive, record and derived-merge commits share one helper. It requires an
  `*exec.ExitError`, extracts only the runner's output, reuses the Daemon's
  classification, and returns a typed refusal that unwraps to the original
  error. It names the hook, first-line subject, exit status, bounded one-line
  output and the specified staged-change or aborted-merge recovery. The code
  documents the third-party output's path into the blocker, status and Run
  Database. No bypass, alternative-subject retry or rewrite was added.
- `daemon.InstalledCommitHooks` exports the existing Git-resolved hook lookup;
  `commitHookInstalled` now checks its list length. No other Daemon file or
  Delivery Engine file changed. All three record proof functions were
  compared with HEAD and remain byte-identical.

Acceptance evidence from focused checks:

1. The synthetic 54-character `0088-` slug records pull request 280 behind a
   real 100-character header hook with the exact 57-character subject, an
   exact new head and a clean worktree. Six subject subtests cover the prefix
   cases and the 70-character ceiling.
2. Hand commits under the old long subject, the short subject and `wip` are
   adopted without another commit. Existing uncommitted-record and inexact
   record tests also pass.
3. Five refusal cases assert the whole error, `errors.As` to
   `*exec.ExitError`, and `environment` from `delivery.ClassifyPark`. These
   cover output naming, a silent hook, two installed hooks, an archive slug
   containing `pre-commit` with only `commit-msg` installed, and a derived
   merge that is aborted back to a clean candidate. Every record-refusal
   case also checks staged work and subsequent hand-commit adoption.
   Additional real-hook cases assert whitespace collapse, an exact 512-byte
   output and a cut inside a UTF-8 rune. Real-Git tests run in isolated child
   processes with temporary homes and Git configuration, without timed waits.
4. A real stale index lock injected at the commit boundary preserves the raw
   Git error without `hook refused`. Existing Daemon hook tests pass unchanged.
   Both touched packages' test binaries cross-compile for Windows.

Focused commands (all prefixed with `rtk proxy` where applicable):

- Before implementation: `go test ./internal/cli -run '^TestPullRequestRecordSubjectNamesTheSpecNumber$' -count=1`
  failed to compile because `pullRequestRecordSubject` did not exist.
- `go test ./internal/cli -run 'Test(PullRequestRecordSubjectNames|DeliveryWorkflowRecordCommitPasses|DeliveryWorkflowAdoptsARecordCommit|DeliveryQueueCommit)' -count=1 -v`
  passed all six new top-level tests, including all fixture cases.
- `go test ./internal/cli -run 'TestDeliveryWorkflow(RecordsThePullRequest|RefusesARecordCommit|AdoptsAnUncommitted|RefusesAnUncommitted)|TestADerivedConflict|TestConflictRecoveryAborts' -count=1`
  passed the existing record and conflict regressions.
- `go test ./internal/daemon -run 'CommitHook|HookRefusal|GitCommitterKeepsPlain' -count=1`
  passed the existing classification, committer and Task recovery checks.
- `go test ./internal/daemon -run '^TestHookOutputExcerptBoundsOutputAtItsHead$' -count=1`
  passed the remaining existing hook-output check.
- `GOOS=windows go test -c -o "$queue_compile_dir/" ./internal/cli ./internal/daemon`
  exited 0; `queue_compile_dir` was a fresh directory under `/private/tmp`.
- `git diff --check` passed. Scope inspection found only the four assigned
  code/test paths plus this Task file; the three proof functions are unchanged.

The declared Verification command, including Windows vet, and repository
Verification were not run; the Daemon owns them. No commit, push, pull request,
Task Graph edit or other Task-file edit was made. The symlinked writable root
prevented sandbox startup, so commands used the resolved Run checkout with
approved escalation. No product blocker or out-of-slice follow-up was found.
