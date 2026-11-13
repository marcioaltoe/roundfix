---
task: task_03
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: completed
type: backend
complexity: medium
---

# Task 03: The review validator's sealed session carries the reviewed repository's marker

## Overview

The review validator still opens its sealed session in a new private directory
outside every checkout, so the session reads no repository file and loads no
repository instruction. When the reviewed checkout's root holds a regular
`.ai-memory.toml`, the validator copies it into that directory first, so
ai-memory, which reads a marker only in the exact working directory of a
non-Git directory outside the home directory, records the session under the
reviewed repository. The agent's sealed rule admits that one regular file and
refuses every other entry. This answers item 2 of the message reported by
pantheon on 2026-11-12
(`references/2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md`).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   the validator's working directory is created with
   `os.MkdirTemp("", "roundfix-review-validator-")`, and
   `validateSealedPromptRequest` refuses a directory holding any entry.
2. MUST change `internal/agent/sealed.go` as `_techspec.md` → API Contract 4
   says: add `const SealedIdentityMarker = ".ai-memory.toml"` with a comment
   saying it is the one file a sealed working directory may hold, a copy of
   the repository's marker so session tooling keyed by the working directory
   records the session under that repository; let
   `validateSealedPromptRequest` accept no entry or exactly one regular file of
   that name and otherwise refuse with
   `sealed Agent working directory must be empty apart from its .ai-memory.toml`;
   and add `CopySealedIdentityMarker(repositoryRoot, workDir string) error`,
   which uses `os.Lstat`, copies nothing for an absent marker, a non-regular
   file or one larger than 65536 bytes, returns an error on a read or write
   failure, and writes the copy with mode 0600. MUST NOT change any other
   sealed rule: privacy, the absolute path, the input bounds and the denied
   tools stay as they are.
3. MUST change `validateReviewConventions` in
   `internal/cli/review_convention_validator.go` as API Contract 5 says: right
   after creating the directory, call
   `agent.CopySealedIdentityMarker(repo.Root, workDir)` and report an error as
   the validator being unavailable. MUST NOT change `internal/baselineacp` or
   any other caller of a sealed prompt.
4. MUST create `internal/agent/sealed_identity_test.go` with
   `TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker` and
   `TestCopySealedIdentityMarkerCopiesOnlyARegularMarker`, each calling
   `t.Parallel()` first, with the cases Testing Approach lists, the refusal
   message word for word, and only temporary directories; a symbolic link is
   made with `os.Symlink`, never a hard link.
5. MUST change `conventionValidatorRunner` in
   `internal/cli/review_convention_validator_test.go` so it accepts a working
   directory that holds only a regular `agent.SealedIdentityMarker` and
   records its bytes (nil when absent), still refusing any other entry and a
   directory that is not private; and add
   `TestTheReviewValidatorSessionCarriesTheRepositoryMarker` with subtests
   `marker` and `no marker`, calling `t.Parallel()` first, as Testing Approach
   says.
6. MUST change `docs/user-guide/commands/review.md` with the second text of API
   Contract 7, word for word, right after "It uses the review's selected
   runtime in a sealed session and may use no tools.".
7. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST NOT start an Agent session or an ACP adapter, call ai-memory, or read
   the host's home directory in any test.

## Subtasks

- [ ] Record the characterization.
- [ ] Add the marker constant, the sealed rule and the copy, with their tests.
- [ ] Copy the marker in the validator and extend its fake runner and test.
- [ ] Add the guide sentence, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A validator session for a checkout that commits `.ai-memory.toml` sees exactly that file, with the same bytes, in its private working directory; without it, an empty directory.
- [ ] A sealed working directory holding any other entry, a marker directory or a marker symbolic link is refused.
- [ ] Baseline analysis is unchanged.

## Context

- instruction: `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- instruction: `internal/agent/sealed_test.go`
- instruction: `internal/baselineacp/analyzer.go`
- creates: `internal/agent/sealed_identity_test.go`
- interface: `internal/agent/sealed.go`
- interface: `internal/cli/review_convention_validator.go`
- interface: `internal/cli/review_convention_validator_test.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `agent.SealedIdentityMarker`, `agent.CopySealedIdentityMarker` and the sealed working directory rule — consumers: `./internal/agent`, `./internal/baselineacp`, `./internal/cli`
- The `review` guide sentence and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v ./internal/agent ./internal/baselineacp 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker TestCopySealedIdentityMarkerCopiesOnlyARegularMarker; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the two tests do not exist; afterwards they pass with every other test of `./internal/agent` and `./internal/baselineacp`.

- `out="$(go test -count=1 -timeout 30m -v -run 'ReviewValidator|ReviewConvention|ReviewResultStatus' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in 'TestTheReviewValidatorSessionCarriesTheRepositoryMarker/marker' 'TestTheReviewValidatorSessionCarriesTheRepositoryMarker/no_marker'; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the test does not exist; afterwards both subtests pass with every existing validator test of `./internal/cli`.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- 'Roundfix copies it into that directory, which holds nothing else' || { echo 'review guide: validator directory sentence missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the sentence is absent; afterwards the guide carries it and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 4; Core Features 4 and 5; Success Metric 4; Success Metric 5; Non-Goals
- `_techspec.md` → Current behavior 6; Interfaces; API Contracts 4, 5 and 7; Invariant 5; Declared breaks; Testing Approach; Build Order 3
- ADR-0395
- [The adopted pantheon message](references/2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md)

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-error-handling, golang-testing, testing-boss, tech-writer, no-workarounds, systematic-debugging, evidence-gate.

Starting-tree characterization (recorded before implementation):
`validateReviewConventions` creates its working directory with
`os.MkdirTemp("", "roundfix-review-validator-")` and immediately passes the
empty directory to the sealed runner after registering cleanup.
`validateSealedPromptRequest` reads that directory with `os.ReadDir` and
refuses every nonempty directory with `sealed Agent working directory must be
empty`. `internal/baselineacp/analyzer.go` independently creates empty private
directories with `withPrivateDirectory`; it copies no repository marker.
The starting worktree held only this Task's Daemon-owned status change.

Implementation:

- The sealed directory admits only an empty directory or one regular
  `SealedIdentityMarker`. All other entries retain a refusal, with the exact
  new message. Privacy, absolute-path validation, input bounds and denied
  tools are unchanged.
- `CopySealedIdentityMarker` uses `os.Lstat`, copies only a regular marker of
  at most 65536 bytes, preserves its bytes in a 0600 file, and returns wrapped
  inspection, read and write errors. Absent, non-regular and oversized
  markers copy nothing.
- The validator copies the reviewed repository's marker immediately after
  creating its private directory and registering cleanup. A copy error uses
  the existing unavailable path. The fake runner checks privacy and the
  single regular marker rule and captures marker bytes, or nil when absent.
- The review guide carries API Contract 7's second text verbatim. The
  Behavior Surface Record was regenerated through its sanctioned recorder
  after the last guide edit. No Baseline caller or other sealed caller was
  changed.

Focused checks (commands executed through `rtk proxy`; these are not the
Task's declared Verification):

1. `go test -count=1 ./internal/agent -run '^(TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker|TestCopySealedIdentityMarkerCopiesOnlyARegularMarker)$'`
   initially exited 1 because the new marker constant and helper did not
   exist. With the agent implementation present but before wiring the
   validator, `go test -count=1 -v ./internal/cli -run '^TestTheReviewValidatorSessionCarriesTheRepositoryMarker$'`
   exited 1: `/marker` received nil instead of the committed bytes, while
   `/no_marker` passed.
2. `go test -count=1 -v ./internal/agent -run '^(TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker|TestCopySealedIdentityMarkerCopiesOnlyARegularMarker)$'`
   exited 0. Acceptance criterion 2: empty and marker directories accepted;
   another file, marker plus another file, marker directory and marker
   symbolic link refused with the exact message. Copy cases for absent,
   directory, symbolic link and 65537 bytes wrote nothing; regular and
   exactly 65536-byte markers retained their bytes, mode 0600 and sealed
   validity. Inspection and destination-write failures returned errors.
3. `go test -count=1 -v ./internal/cli -run '^Test(TheReviewValidatorSessionCarriesTheRepositoryMarker|ReviewValidatorDismissesADaemonSettlementAsC2|ReviewValidatorFailureLeavesEveryAskedFindingStanding|ReviewValidatorUsesTheSuccessfulFallbackSelection|ReviewValidatorCanDismissNoFailureOnlyWhenEligible)$'`
   exited 0 after wiring the copy. Acceptance criterion 1: `/marker` and
   `/no_marker` passed with exactly the committed bytes or nil, a private
   directory outside the reviewed checkout, the selected runtime, and
   directory cleanup. The selected existing dismissal, fail-closed and
   fallback-selection checks also passed.
4. `go test -count=1 -v ./internal/baselineacp -run '^Test(RuleSegmentationUsesPrivateCheckoutFreeDirectories|AnalyzerNoPersistenceOrRepositoryMutation|ACPXReadOnlyArguments)$'`
   exited 0. Acceptance criterion 3: existing Baseline checks still observe
   fresh empty private directories, no checkout mutation or exposure, and
   the unchanged denied-tools arguments; `internal/baselineacp` has no diff.
5. `git diff --check` exited 0.
6. As the final implementation step, after the last guide edit,
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   exited 0 and refreshed `docs/references/behavior-surfaces.json`.

All new tests use temporary directories and fake runners; no Agent session,
ACP adapter, ai-memory call, network lookup or host-home read was introduced.
Live session attribution remains the operator's post-release check under the
PRD's Unreachable Acceptance. No follow-up outside this Task was implemented.

Task status remains `in_progress` and Daemon-owned. Declared Verification and
repository Verification were not run; there is no terminal Task verdict.
