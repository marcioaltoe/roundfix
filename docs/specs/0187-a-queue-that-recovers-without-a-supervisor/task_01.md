---
task: task_01
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: pending
type: backend
complexity: medium
---

# Task 01: Post-merge cleanup refreshes the default branch before it reads it

## Overview

After the Delivery Queue squash-merges an item, `ReleaseMergedRuns` proves the recorded merge commit against the default branch before it releases the Spec's Runs. `resolveMergedReleaseEvidence` in `internal/cli/deliver_workflow.go` reads the local `refs/heads/<default>` and the merge commit without ever fetching. On 2026-09-29, for Spec 0183 (#282), it reported `merge commit … does not resolve to a commit`, and later `is not on default branch`. The item worktree and Runs were left behind. This Task refreshes the default branch from the delivery remote, the same remote `CreateItemBranch` fetches, before resolving, and reads the refreshed remote-tracking branch.

## Requirements

1. MUST resolve the delivery remote in `resolveMergedReleaseEvidence` the way `CreateItemBranch` does: `Watch.PushRemote`, or `origin` when it is empty or when the workflow has no loaded config.
2. MUST, when `git remote get-url <remote>` succeeds, run `git fetch <remote> <default>` before resolving the merge commit, the candidate head or the default head. A failed fetch MUST return an error that starts with `refresh default branch "<default>"`. The default head MUST then be `refs/remotes/<remote>/<default>`, and the merge commit MUST be its ancestor.
3. MUST keep today's local resolution unchanged when the repository has no such remote, and keep every other merge-commit and candidate-head check unchanged.
4. MUST keep every test in `internal/cli/deliver_release_evidence_test.go` green, rename or remove no top-level test, and change no exported function signature.
5. MUST put the new tests in `internal/cli/deliver_release_refresh_test.go`, over a bare `origin` and two real clones in temporary directories, never a live remote.
6. MUST describe the cleanup of a merged item in `docs/user-guide/commands.md` and the Delivery queue section of `.agents/skills/roundfix/SKILL.md` with the phrase `refreshes the default branch from the delivery remote`, then run `make skills-sync`.

## Subtasks

- [ ] Resolve the delivery remote and refresh the default branch before reading it.
- [ ] Read the merge commit against the refreshed remote-tracking branch.
- [ ] Describe the behavior in the guide and the skill, then sync the mirror.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An item merged and pushed from one clone releases its Runs from the other clone without a prior fetch.
- [ ] An unreachable `origin` returns an error starting with `refresh default branch` and releases nothing.
- [ ] A repository with no `origin` resolves exactly as before, and the existing evidence tests pass unchanged.
- [ ] The guide, the skill and its mirror carry `refreshes the default branch from the delivery remote`, and `make skills-sync-check` passes.

## Context

- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_release_evidence_test.go`
- creates: `internal/cli/deliver_release_refresh_test.go`
- instruction: `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReleaseMergedRunsFetchesTheMergeCommitBeforeReleasing|TestReleaseMergedRunsReportsAFailedRefresh|TestReleaseMergedRunsWithoutARemoteResolvesLocally|TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit|TestReleaseMergedRunsStillReleasesWithValidEvidence)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReleaseMergedRunsFetchesTheMergeCommitBeforeReleasing TestReleaseMergedRunsReportsAFailedRefresh TestReleaseMergedRunsWithoutARemoteResolvesLocally TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit TestReleaseMergedRunsStillReleasesWithValidEvidence; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three new named tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands.md|refreshes the default branch from the delivery remote" ".agents/skills/roundfix/SKILL.md|refreshes the default branch from the delivery remote" "skills/roundfix/SKILL.md|refreshes the default branch from the delivery remote"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && make skills-sync-check` — expected: exit 0; before this Task the phrase `refreshes the default branch from the delivery remote` is in neither the guide nor the skill, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Core Feature 1; Success Metric 1
- [_techspec.md](_techspec.md) — Cleanup refreshes the default branch first; API Contract 1; Testing Approach 1; Build Order 1
- [references/2026-09-30-post-merge-cleanup-needs-the-merge-commit-locally.md](references/2026-09-30-post-merge-cleanup-needs-the-merge-commit-locally.md)
- ADR-0161; ADR-0090

## Result
