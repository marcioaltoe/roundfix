---
task: task_02
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
status: pending
type: backend
complexity: medium
---

# Task 02: The skills project root uses the same walk, and the commands under a rejected .git say so

## Overview

The project scope of `skills install` holds `findSkillsGitRoot`, a copy of
the walk task_01 replaced, so under a `~/.git` that Git rejects it would still
install into the home directory. This Task makes
`defaultResolveSkillsProjectRoot` call `config.FindGitRoot`, deletes the
copy, and pins the commands the maintainer's MacBook reached in Surface
Transcripts 1 to 3. It answers the maintainer report of 2026-11-05, adopted
as
[references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md](references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md).

## Requirements

1. MUST make `defaultResolveSkillsProjectRoot` in `internal/cli/cli.go` call
   `roundconfig.FindGitRoot(cwd)` where it called `findSkillsGitRoot(cwd)`,
   delete `findSkillsGitRoot`, and keep its `git rev-parse --show-toplevel`
   fallback and its refusal `project skill install requires running inside a
   git repository or passing --dir` byte for byte (`_techspec.md` → API
   Contracts 2 and 3).
2. MUST NOT change any other production code, `internal/cli/cli_test.go` or
   the tests task_01 changed.
3. MUST add `internal/cli/rejected_git_directory_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, with:
   - `TestRejectedGitDirectoryTranscripts`: Surface Transcripts 1, 2 and 3 of
     `_techspec.md`, one subtest each named `Surface Transcript 1`, `Surface
     Transcript 2` and `Surface Transcript 3`, each with a fresh `<home>`
     holding `<home>/.git/gk/` and `<home>/.git/info/` and the working
     directory `<home>/drive`, set through `setCommandEnvironmentForTest`,
     run through `runCLI`, comparing the whole stdout, the whole stderr and
     the exit code with the transcript, then asserting that
     `<home>/.roundfixrc.yml` does not exist;
   - `TestSkillsProjectRootSkipsARejectedGitDirectory`:
     `defaultResolveSkillsProjectRoot` from `<home>/drive` returns the
     refusal of Requirement 1, and from `<home>/drive/repo/sub`, where
     `<home>/drive/repo` is made with `gittest.InitRepo`, returns
     `<home>/drive/repo`.
   Each top-level test calls `t.Parallel()` first and uses only `t.TempDir()`
   directories.

## Subtasks

- [ ] Point the skills project root at `FindGitRoot` and delete the copy.
- [ ] Add the transcript and skills root tests.

## Acceptance Criteria

- [ ] Under a rejected `~/.git`, `runs list` prints `No Runs found.` with no warning, `reconcile` refuses the missing repository, and `init --scope project` refuses and writes no Project Config, each exactly as its transcript.
- [ ] The skills project root ignores the rejected `.git` and still finds a nested repository.

## Context

- instruction: `docs/adr/0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/rejected_git_directory_test.go`

## Consumers

- `defaultResolveSkillsProjectRoot`, the project root of `skills install`, `skills check` and `baseline update` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRejectedGitDirectoryTranscripts|TestSkillsProjectRootSkipsARejectedGitDirectory)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRejectedGitDirectoryTranscripts TestSkillsProjectRootSkipsARejectedGitDirectory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; ! grep -q 'func findSkillsGitRoot' internal/cli/cli.go || { echo 'findSkillsGitRoot still defined' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards both pass, the copy of the walk is gone, and the whole CLI suite passes.

## References

- `_prd.md` → Goals; Core Feature 2; Success Metric 1; Success Metric 4
- `_techspec.md` → Current behavior 2; API Contracts 2 and 3; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariant 3; Testing Approach; Build Order 2
