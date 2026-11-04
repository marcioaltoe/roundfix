---
task: task_01
spec: 0347-read-only-run-commands-work-outside-a-repository
status: completed
type: backend
complexity: medium
---

# Task 01: The worktree-location refusal names its Git root and Load can return it as data

## Overview

`config.Load` fails every command when `worktree.location` lies inside the
Git root of the working directory, and the refusal names neither the root nor
the `.git` that made it one, so a stray `.git` in the home directory cannot be
found from the message. This Task appends the Git root and its `.git` path to
the refusal, keeping the old text as a prefix, and adds the
`WorktreeLocationOptional` load option that returns the refusal in
`Loaded.WorktreeLocationErr` instead of failing. It answers the maintainer
report of 2026-11-03, adopted as
[references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md](references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md).

## Requirements

1. MUST change the inside-the-tree error of `ResolveWorktreeLocation` in
   `internal/config/config.go` to exactly API Contract 1 of `_techspec.md`:
   `worktree.location must not be inside the repository tree: "<location>" is
   inside the Git root "<root>" (found "<root>/.git")`, built with `%q` and
   `filepath.Join(gitRoot, ".git")`; the text up to the closing quote of the
   location keeps its bytes, and every other error of the function keeps its
   bytes.
2. MUST add `WorktreeLocationOptional bool` to `LoadOptions` and
   `WorktreeLocationErr error` to `Loaded`, with the doc comments of
   `_techspec.md` → Interfaces, and make `Load` follow API Contract 2: with the
   option, a `ResolveWorktreeLocation` error is stored in
   `WorktreeLocationErr`, `Config.Worktree.Location` is `""` and `Load`
   returns no error; without it, `Load` is unchanged.
3. MUST NOT change `findGitRoot`, `Validate`, any other `Load` step or any
   caller in `internal/cli`.
4. MUST add `internal/config/worktree_location_test.go` with
   `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` and
   `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal` as `_techspec.md`
   → Testing Approach describes: the rows `git directory`, `git file` (a
   `.git` file holding `gitdir: ...`) and `location is the root`, the exact
   error and its old-text prefix, and the option under a stray `.git` and in
   a repository whose tree does not hold the location. Each top-level test
   calls `t.Parallel()` first and uses only `t.TempDir()` directories.
5. MUST keep `TestLoadRejectsWorktreeLocationInsideRepository` passing
   unchanged.

## Subtasks

- [ ] Extend the refusal text.
- [ ] Add the option and the field to `Load`.
- [ ] Add the two tests and run the config and CLI suites.

## Acceptance Criteria

- [ ] The refusal names the Git root and its `.git` path for a `.git` directory, a `.git` file and a location equal to the root, and starts with its old text.
- [ ] With the option, `Load` succeeds under a stray `.git` and returns the same refusal as data; in a repository it returns the same location as without the option.

## Context

- interface: `internal/config/config.go`
- creates: `internal/config/worktree_location_test.go`

## Consumers

- `config.LoadOptions.WorktreeLocationOptional`, `config.Loaded.WorktreeLocationErr` and the `worktree.location` refusal text — consumers: `./internal/config`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit|TestLoadWithOptionalWorktreeLocationReturnsTheRefusal)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit TestLoadWithOptionalWorktreeLocationReturnsTheRefusal; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config ./internal/cli` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards both pass with the whole config and CLI suites, which read the refusal.

## References

- `_prd.md` → Goals; Core Features 1 and 2; Success Metric 3
- `_techspec.md` → Current behavior 1-2; API Contracts 1 and 2; Invariants 1 and 5; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

The inside-tree refusal now appends the quoted Git root and its `.git` path
using the exact API Contract 1 format. `LoadOptions.WorktreeLocationOptional`
and `Loaded.WorktreeLocationErr` carry the specified comments. Optional loading
retains the resolution refusal as data and leaves the location empty; ordinary
loading still returns the refusal. Repository discovery, validation, other
load steps and CLI callers retain their implementation.

Acceptance evidence:

- Refusal text: `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit`
  exercises `Load` from below the root for `git directory`, `git file` with
  `gitdir: ...` content, and `location is the root`. All three cases assert
  the exact message and the old-text prefix and passed in the focused check.
- Optional loading: `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal`
  passed for stray Git directories and files, comparing the retained refusal
  with strict loading, checking the empty location, root and config paths,
  and preserving a User Config value and the remaining configuration. Its
  repository case checks the expected resolved location, nil refusal and
  equality with the entire strict-load result. Both new top-level tests call
  `t.Parallel()` first and allocate directories only through `t.TempDir()`.

Focused checks:

- Before the implementation,
  `rtk proxy go test -count=1 -timeout 2m -run '^TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit/git_directory$' ./internal/config`
  exited 1: the observed refusal contained only the old location text and
  lacked the Git root and `.git` suffix.
- After the implementation,
  `rtk proxy go test -count=1 -timeout 2m -v -run '^(TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit|TestLoadWithOptionalWorktreeLocationReturnsTheRefusal|TestLoadRejectsWorktreeLocationInsideRepository|TestLoadExpandsWorktreeLocationHome)$' ./internal/config`
  exited 0; all four tests and all six new subtests passed. The existing
  rejection test was not edited.
- `rtk proxy go test -count=1 -timeout 2m -v -run '^(TestRunRunsListOutsideRepositoryRequiresAll|TestRunOperationalCommandRejectsInvalidConfigAndArtifactDirectory)$' ./internal/cli`
  exited 0; both selected CLI consumer tests passed.
- `rtk proxy git diff --check` exited 0. Changed paths are confined to
  `internal/config/config.go`, the new
  `internal/config/worktree_location_test.go`, and this Task's Result;
  the pre-existing Daemon status change remains intact.

The declared Verification, full config/CLI suites and repository gates were
not run; they remain Daemon-owned. No follow-up work was identified in this
slice. Commands used the assigned physical checkout through the approved
execution fallback because the sandbox could not start with its symlinked
writable root. No commit, push or pull request was made.
