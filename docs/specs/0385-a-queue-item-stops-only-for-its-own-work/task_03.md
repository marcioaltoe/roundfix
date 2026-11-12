---
task: task_03
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: pending
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
