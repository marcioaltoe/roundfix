---
task: task_01
spec: 0347-read-only-run-commands-work-outside-a-repository
status: pending
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
