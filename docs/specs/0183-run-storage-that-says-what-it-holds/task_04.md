---
task: task_04
spec: 0183-run-storage-that-says-what-it-holds
status: pending
type: backend
complexity: medium
---

# Task 04: `gc sanitize` recognizes a pre-key default Artifact Root

## Overview

In a bare-repository layout, or a linked worktree of a `--separate-git-dir` repository, Spec 0162 moved the default Artifact Root from the checkout path to the common Git directory. `gc sanitize` compares a recorded root only with the default re-derived from the repository key, so a root created before the upgrade is classified `overridden` and never reclaimed. This Task makes sanitation also accept the default derived from the Run's recorded checkout, without re-resolving that checkout. The recorded roots and Runs come from the machine-wide Run Database, and the directories live under Roundfix Home. `--apply` still removes only proven retention-eligible or absent Run artifact directories.

## Requirements

1. MUST add `GitRoot string` to `store.ArtifactRootRun` in `internal/store/journal.go` and make `DiscoverArtifactRoots` select the recorded `git_root` into it beside the coalesced repository key.
2. MUST add `DefaultArtifactDirectoryForPath(path string, homeDir string) (string, error)` in `internal/config/config.go`, returning `<home>/.roundfix/artifacts/<repoID(clean path)>`, refusing an empty path or home, and never calling `RepositoryRoot`.
3. MUST replace the per-repository default comparison in `classifyGCSanitationRoot` (`internal/cli/gc.go`) with a per-Run check that runs after every existing guard. A Run accepts the root when it equals `ResolveArtifactDirectory("", run.Repository, home)` or `DefaultArtifactDirectoryForPath(run.GitRoot, home)`. A Run matching neither makes the root `overridden`, with evidence naming both defaults. An error deriving the key default keeps the root `unsafe` unless the checkout default matches. When the checkout default matched, the evidence names that checkout.
4. MUST NOT weaken any existing safety guard: clean absolute path, inside Roundfix Home and not equal to it, physical directory, physical `runs` directory, and safe Run IDs.
5. MUST put the tests in `internal/cli/gc_sanitize_pre_key_root_test.go`, against a real bare clone with `git worktree add` and a temporary Roundfix Home, and a unit test of the helper in `internal/config/artifact_directory_for_path_test.go`. MUST keep `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval` and every existing sanitation test green without renaming any.

## Subtasks

- [ ] Carry the recorded checkout into sanitation's Run evidence.
- [ ] Derive a path's default Artifact Root without re-resolving it.
- [ ] Accept either default per Run and keep every guard.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] In a bare clone with a linked worktree, a root derived from the checkout path of a terminal Run is classified `orphaned`, and `gc sanitize --apply` removes its retention-eligible directory.
- [ ] A root equal to neither default is still `overridden`, and its evidence names both defaults.
- [ ] A root derived from the repository key keeps its current classification.
- [ ] `DefaultArtifactDirectoryForPath` keeps a linked checkout's own path, where `ResolveArtifactDirectory` maps it to the common Git directory, and refuses an empty path or home.

## Context

- interface: `internal/store/journal.go`
- interface: `internal/config/config.go`
- creates: `internal/config/artifact_directory_for_path_test.go`
- interface: `internal/cli/gc.go`
- creates: `internal/cli/gc_sanitize_pre_key_root_test.go`
- instruction: `docs/history/specs/0162-a-durable-repository-key-per-run/_prd.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGCSanitizeReclaimsAPreKeyDefaultRootInABareLayout|TestGCSanitizeStillPreservesARootEqualToNeitherDefault|TestGCSanitizeKeepsAKeyDerivedRootClassification)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestGCSanitizeReclaimsAPreKeyDefaultRootInABareLayout TestGCSanitizeStillPreservesARootEqualToNeitherDefault TestGCSanitizeKeepsAKeyDerivedRootClassification; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^TestDefaultArtifactDirectoryForPathKeepsTheCheckoutPath$" ./internal/config 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestDefaultArtifactDirectoryForPathKeepsTheCheckoutPath"` — expected: exit 0; before this Task the test does not exist.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 3; Success Metric 4
- [_techspec.md](_techspec.md) — `gc sanitize` recognizes a pre-key default root; Interfaces; Data Models; API Contract 3; Testing Approach 3; Build Order 4
- [references/2026-09-25-bare-layout-artifact-roots-over-preserved.md](references/2026-09-25-bare-layout-artifact-roots-over-preserved.md)
- ADR-0033
