---
task: task_01
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
status: completed
type: backend
complexity: medium
---

# Task 01: The Git root walk accepts only a .git entry that Git accepts

## Overview

`findGitRoot` takes the first ancestor holding any `.git` entry, so a
`~/.git` that holds only `gk/` and `info/`, which Git rejects, made the home
directory a Git root on the maintainer's MacBook. This Task replaces it with
the exported `config.FindGitRoot`, which accepts a `.git` entry only when Git
would (ADR-0362), gives the tests one helper that writes the smallest `.git`
directory Git accepts, and moves every fixture that faked a repository with an
empty `.git` directory onto it, because the repository Verification at this
Task's settlement runs them under the new walk. The configuration guide, which
describes how configuration finds the Git root, states the new rule. It
answers the maintainer
report of 2026-11-05, adopted as
[references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md](references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md).

## Requirements

1. MUST replace `findGitRoot` in `internal/config/config.go` with
   `FindGitRoot(start string) string` and an unexported
   `gitEntryAccepted(path string) bool`, with the doc comments of
   `_techspec.md` → Interfaces, following API Contracts 1 and 2 exactly:
   - a regular `.git` file counts when its content starts with `gitdir:`,
     whatever its target;
   - a `.git` directory counts when its `HEAD` (read with `os.Lstat`) is a
     symbolic link to a target starting with `refs/`, or a file starting with
     `ref:`, optional whitespace and `refs/`, or with a run of exactly 40 or 64
     hexadecimal digits, and its common directory (the trimmed `commondir`
     path resolved against the `.git` directory, or the `.git` directory
     itself) holds an `objects` directory and a `refs` directory;
   - everything else, and every read error, does not count, and the walk goes
     on to the parent.
   The walk starts no process and reads no environment variable.
2. MUST make `Load`, `Init` and `profileConfigTarget`
   (`internal/config/profile_config.go`) call `FindGitRoot`, and leave no
   other copy of the walk in `internal/config`.
3. MUST keep every message of API Contract 3 byte for byte, and MUST NOT
   change `ResolveWorktreeLocation`, `repositoryRootMetadata`,
   `internal/cli/cli.go` or any production code outside `internal/config`.
4. MUST add `MarkRepository(t testing.TB, dir string)` to
   `internal/gittest/gittest.go` as API Contract 4 states: `os.MkdirAll` of
   `<dir>/.git/objects` and `<dir>/.git/refs`, `<dir>/.git/HEAD` holding
   `ref: refs/heads/main\n`, and `t.Fatalf` on any error.
5. MUST replace each empty-`.git` fixture listed in `_techspec.md` → Declared
   breaks (`mustMkdir(t, filepath.Join(<dir>, ".git"))` or the
   `os.Mkdir(filepath.Join(<dir>, ".git"), ...)` form) with
   `gittest.MarkRepository(t, <dir>)`, including `withCLIWorkspace` in
   `internal/cli/cli_test.go` and the `.git` directory rows of the Spec 0347
   tests in `internal/config/worktree_location_test.go` and
   `internal/cli/read_only_run_commands_test.go`, and import
   `roundfix/internal/gittest` in the project-import group where a file lacks
   it. No assertion, expected text or test name changes.
6. MUST make `fakeWorktree.Snapshot` in `internal/cli/cli_test.go` write
   `src/agent-change.go` under the work directory it receives (creating
   `src/`) before reporting it, as `_techspec.md` → Declared breaks
   describes, with a comment saying the reported change exists on disk.
7. MUST add `internal/config/git_root_test.go`, opening with the `Suite`,
   `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, with
   `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts`,
   `TestFindGitRootAgreesWithGit` and
   `TestLoadUnderARejectedGitDirectoryFindsNoGitRoot` over the shape table of
   `_techspec.md` → Testing Approach, every row named there by its name, each
   top-level test calling `t.Parallel()` first and using only `t.TempDir()`
   directories, `gittest.InitRepo`, `gittest.Run` and `gittest.IsolatedEnv`.
   `TestFindGitRootAgreesWithGit` runs `git rev-parse --show-toplevel` with
   `GIT_CEILING_DIRECTORIES` set to the parent of the physical base directory
   and compares after `filepath.EvalSymlinks`.
8. MUST replace, in `docs/user-guide/configuration.md` next to the
   `worktree` block, the sentence "Repository detection finds any `.git`
   entry in an ancestor, so a stray `.git` in a parent such as the home
   directory is found this way." with the rule sentence of `_techspec.md` →
   API Contract 5, word for word, and change no other text of the guide.
9. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of the configuration guide,
   and never hand-edit it.

## Subtasks

- [ ] Write `FindGitRoot` and its `.git` test, and move the three callers.
- [ ] Add `gittest.MarkRepository` and move every empty-`.git` fixture onto it.
- [ ] Make `fakeWorktree` write the file it reports.
- [ ] Add the three tests and run the affected suites.
- [ ] State the rule in the configuration guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A `.git` directory holding only `gk/` and `info/`, an empty one, one without `refs/`, one whose `HEAD` names nothing and a file without `gitdir:` mark no Git root; every accepted shape of the table marks its directory.
- [ ] For every row Git answers, `FindGitRoot` agrees with `git rev-parse --show-toplevel`.
- [ ] `Load` under a rejected `~/.git` finds no Git root and keeps the default worktree location; project `init` and `profiles configure` refuse with their unchanged messages.
- [ ] The config, CLI, judge and gittest suites pass with the fixtures on `gittest.MarkRepository`.
- [ ] The configuration guide states the rule sentence and no longer says that detection finds any `.git` entry.

## Context

- instruction: `docs/adr/0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`
- interface: `internal/config/config.go`
- interface: `internal/config/profile_config.go`
- interface: `internal/gittest/gittest.go`
- creates: `internal/config/git_root_test.go`
- interface: `internal/config/config_test.go`
- interface: `internal/config/delivery_check_timeout_test.go`
- interface: `internal/config/jev_ceiling_test.go`
- interface: `internal/config/pre_pr_review_omit_paths_test.go`
- interface: `internal/config/profile_deviation_test.go`
- interface: `internal/config/profile_proof_setup_timeout_test.go`
- interface: `internal/config/profiles_characterization_test.go`
- interface: `internal/config/review_source_scope_test.go`
- interface: `internal/config/run_retention_test.go`
- interface: `internal/config/run_slot_wait_test.go`
- interface: `internal/config/secondbrain_test.go`
- interface: `internal/config/settlement_config_test.go`
- interface: `internal/config/worktree_caches_test.go`
- interface: `internal/config/worktree_shared_caches_test.go`
- interface: `internal/config/worktree_location_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/runs_list_outside_repository_test.go`
- interface: `internal/cli/read_only_run_commands_test.go`
- interface: `internal/judge/task_acceptance_measure_test.go`
- interface: `internal/judge/task_acceptance_stats_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `config.FindGitRoot` and the Git root every `config.Load` caller reads — consumers: `./internal/config`, `./internal/cli`, `./internal/judge`
- `gittest.MarkRepository` — consumers: `./internal/gittest`, `./internal/config`, `./internal/cli`, `./internal/judge`
- The configuration guide fingerprint in the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFindGitRootAcceptsOnlyAGitEntryGitAccepts|TestFindGitRootAgreesWithGit|TestLoadUnderARejectedGitDirectoryFindsNoGitRoot)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFindGitRootAcceptsOnlyAGitEntryGitAccepts TestFindGitRootAgreesWithGit TestLoadUnderARejectedGitDirectoryFindsNoGitRoot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; printf '%s\n' "$out" | grep -q -- '--- PASS: TestFindGitRootAcceptsOnlyAGitEntryGitAccepts/only_gk_and_info ' || { echo 'missing row: only gk and info' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/config ./internal/gittest ./internal/judge ./internal/cli` — expected: exit 0. Before this Task none of the three tests exists, so the first pass line is missing; afterwards they pass, the MacBook row among them, with the whole config, gittest, judge and CLI suites, whose fixtures now use `gittest.MarkRepository`.

- `t="$(tr -s '[:space:]' ' ' < docs/user-guide/configuration.md)"; printf '%s' "$t" | grep -qF -- 'is skipped as Git skips it' || { echo 'missing rule in the configuration guide' >&2; exit 1; }; ! printf '%s' "$t" | grep -q -- 'finds any .\.git. entry in an ancestor' || { echo 'old rule still in the configuration guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide says detection finds any `.git` entry, so the first check fails; afterwards it states the rule sentence and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 2; Success Metric 3; Success Metric 5
- `_techspec.md` → Current behavior 1 and 5; API Contracts 1-5; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1


## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, golang-cli, golang-error-handling, no-workarounds, systematic-debugging, tech-writer, evidence-gate.

Implemented the assigned slice; Task status and settlement remain Daemon-owned.

- Acceptance 1: `FindGitRoot` and `gitEntryAccepted` implement API Contracts 1 and 2 using filesystem reads only. `Load`, `Init` and `profileConfigTarget` share the exported walk. The acceptance test covers every authored shape, plus symlink HEAD, whitespace after `ref:`, short detached HEAD, invalid ref paths, non-file HEAD, absent/non-directory objects and refs, absolute commondir and symlinked `.git`. Rejected inner entries are also checked beneath an accepted ancestor.
- Acceptance 2: the Git oracle passes for all 12 authored rows where Git answers the discovery question, plus 10 additional compatible edge shapes. It uses `gittest.IsolatedEnv`, a ceiling at the parent of the physical base and physical-path comparison. The two authored gitfile error cases are tested by the acceptance test and excluded from the oracle as specified.
- Acceptance 3: the initial focused regression failed because `Load` treated temporary `.git/{gk,info}` as a Git root and refused the default worktree location. After the change, `TestLoadUnderARejectedGitDirectoryFindsNoGitRoot` passes: no Git root or Project Config path, the default home worktree location, both exact project-scope refusal strings, and no Project Config written at the rejected root. Inspection confirms `ResolveWorktreeLocation`, `repositoryRootMetadata` and `internal/cli/cli.go` remain unchanged.
- Acceptance 4: `gittest.MarkRepository` creates the required HEAD, objects and refs with fatal error handling. All 82 declared empty-marker fixtures in 21 existing test files now use it, including the Spec 0347 directory rows. Existing test names and assertions are unchanged. `fakeWorktree.Snapshot` writes its reported source file under the received work directory and returns filesystem errors with context. Focused consumer checks pass; the complete consumer suites are reserved for Daemon Verification.
- Acceptance 5: the configuration guide differs from HEAD only by the exact API Contract 5 sentence replacement. The sanctioned recorder passes and changes only that guide's fingerprint in the Behavior Surface Record. The recorder is rerun after this Result as the final implementation step; no generated record was hand-edited.

Focused checks (all exited 0 after the final code edits):

```sh
go test -count=1 ./internal/config -run '^(TestFindGitRoot|TestLoadUnderARejected|TestLoadApplies|TestProjectInit|TestInit|TestProfileConfig|TestProfilesConfig|TestProfilesConfigure|TestWorktreeLocation|TestLoadWithOptional|TestRunRetention|TestWorktreeShared|TestWorktreeCaches|TestSettlement|TestSecondbrain|TestReviewSource|TestRunSlotWait|TestDeliveryCheck|TestJev|TestPrePRReview)'
go test -count=1 ./internal/cli -run '^(TestRunResolveNoAgentConsoleSuppressesAgentDisplayOnly|TestRunResolveFinalPushRunsOnceAfterAllUnresolvedTerminal|TestRunResolveProcessesAllBatchesBeforeFinalPush|TestWatchFinalPushRunsOncePerCleanRoundThroughEngine|TestAttachRendersWatchDaemonEventsInTimeline|TestRunOperationalCommandAcceptsMVPFlags|TestReadOnlyRunCommands.*|TestStrayGitDirectoryTranscripts|TestRunRunsList.*|TestRunsListOutsideARepositoryTranscript)$'
go test -count=1 ./internal/judge -run '^(TestTaskAcceptanceResolvesConfiguredAndArchivedRoot|TestTaskAcceptanceRecordIsConsistent)$'
go test -count=1 ./internal/gittest -run '^TestInitRepo'
go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
git diff --check
```

A read-only scope/fixture inspection also confirmed that every changed path belongs to this Task, fixture edits are limited to marker/import replacements plus the declared fake-worktree repair, and the guide contains no other changes.

Follow-up observation outside the authored shape table: an exploratory 65-hex-digit HEAD was accepted by the installed Git, whereas API Contract 1 explicitly requires a run of exactly 40 or 64 digits. The implementation retains that explicit contract; the oracle evidence above is bounded to the authored rows and the compatible extra shapes, not every possible malformed HEAD. No behavior outside this Task was changed for that observation.

Not run: the Task's declared Verification commands, `make verify`, `make verify-changed`, and full consumer/documentation suites; the Daemon owns them. No commit, push, Pull Request, Task status edit, Task Graph edit or other Task edit was made.

## Carry-forward provenance

- Source Run: `run_20261105T105803Z_8c2fa84cfd81bb6f`
- Source commit: `c9c7126c632cfdcfd5d047c6c40660b8dd2e8d1c`
