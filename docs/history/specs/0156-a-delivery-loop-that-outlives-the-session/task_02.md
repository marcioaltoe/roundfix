---
task: task_02
spec: 0156-a-delivery-loop-that-outlives-the-session
status: completed
type: backend
complexity: medium
---

# Task 02: The pull request boundary

## Overview

No code in Roundfix opens, follows or merges a pull request today. Add that boundary behind an interface, backed by the GitHub CLI already authenticated on the host, with a fake for tests.

## Requirements

1. MUST push a branch and report the remote head.
2. MUST find an open pull request for a head branch before creating one, and create one only when none exists.
3. MUST report the current-head checks of a pull request.
4. MUST report an existing merge instead of merging again.
5. MUST provide a fake implementation for tests, and MUST NOT call GitHub from any test.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An existing pull request for the head is returned instead of a new one.
- [ ] An already-merged pull request is reported, not merged again.
- [ ] No test reaches the network.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement.go`
- creates: `internal/delivery/github.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestPullRequestBoundaryReusesAnOpenPullRequest$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPullRequestBoundaryReusesAnOpenPullRequest"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestPullRequestBoundaryDoesNotMergeTwice$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPullRequestBoundaryDoesNotMergeTwice"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components

## Result

Implemented a `PullRequestBoundary` with a GitHub CLI adapter and a test fake.
The adapter pushes and reads back the remote head, finds an Open Pull Request
before creating one, reports checks bound to a stable current head, and reads
merge state before issuing a squash merge for the expected head. Pending
checks remain reportable through the GitHub CLI's exit code 8.

Focused checks:

- Pre-change inspection found no `internal/delivery` package or pull request boundary.
- The first `go test ./internal/delivery` attempt was blocked before compilation because the sandbox denied the default macOS Go build cache.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/delivery` passed.
- `GOCACHE=/private/tmp/roundfix-task02-gocache go vet ./internal/delivery` passed.
- `gofmt` left `internal/delivery/github.go` and `internal/delivery/github_test.go` formatted.

Acceptance evidence:

- An existing pull request for the head is returned instead of a new one: `TestPullRequestBoundaryReusesAnOpenPullRequest` scripts only `gh pr list`; any attempted `gh pr create` is an unexpected command and fails the test. `TestPullRequestBoundaryCreatesOnlyWhenNoOpenPullRequestExists` covers the converse.
- An already-merged pull request is reported, not merged again: `TestPullRequestBoundaryDoesNotMergeTwice` scripts only the initial `gh pr view`, returns the recorded merge commit, and fails on any attempted `gh pr merge`. `TestPullRequestBoundaryMergesTheExpectedHead` and `TestPullRequestBoundaryRefusesToMergeAnUnexpectedHead` cover the new-merge and stale-head paths.
- No test reaches the network: every `GitHubCLI` test injects `scriptedCommandRunner`, which returns in-memory command results without starting `git` or `gh`; `fakePullRequestBoundary` also implements the complete interface for downstream delivery-engine tests.

Additional requirement evidence:

- `TestPullRequestBoundaryPushesBranchAndReportsRemoteHead` proves a push is followed by `git ls-remote` and returns the observed SHA; its negative companion refuses an absent remote head.
- `TestPullRequestBoundaryReportsCurrentHeadChecks` and `TestPullRequestBoundaryReportsFailedCurrentHeadChecks` prove pending and failing current-head checks are returned, while `TestPullRequestBoundaryRejectsChecksFromAMovedHead` refuses a report if the PR Head Branch moves during collection.

The Daemon-owned commands under `## Verification` were not run.

## Carry-forward provenance

- Source Run: `run_20260924T121521Z_744099e82ec81b88`
- Source commit: `d10c9326977ed526fba8f72a40ff51db7f036fd8`
