---
task: task_02
spec: 0153-a-reviewer-the-workflow-runs
status: completed
type: backend
complexity: medium
---

# Task 02: The diff the reviewer is handed

## Overview

A reviewer sent to find its own evidence can answer without finding any. The
first attempt at this Spec passed two commit identifiers and asked the session
to inspect the range; across three review rounds it could edit the candidate,
then read nothing, then read files without ever seeing what changed — and in
that last shape it could return `No findings` having reviewed nothing.

This Task makes the command gather the evidence and hand it over.

## Requirements

1. MUST compute the candidate diff from the base commit to the head commit.
2. MUST include that diff in the prompt as content, so the reviewer judges what
   it was given rather than what it must go find.
3. MUST run the session with read-only capabilities that permit opening a file
   the diff references and deny every mutation.
4. MUST prove, by test, that the prompt carries the diff itself and not merely
   the commit identifiers.
5. MUST prove, by test, that the session cannot write.
6. MUST NOT rely on the reviewer having git, shell, or any diff-producing tool.

## Subtasks

- [ ] Compute the candidate diff.
- [ ] Build the prompt that carries it.
- [ ] Set the read-only capability set for the session.
- [ ] Add tests for the prompt content and the capability set.

## Acceptance Criteria

- [ ] The prompt contains the diff text for a fixture candidate.
- [ ] The prompt is not satisfied by commit identifiers alone: a test asserts a
      changed line from the fixture appears in it.
- [ ] The request the runner receives permits reading and denies writing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/agent/agent.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReviewPromptCarriesTheCandidateDiff" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewPromptCarriesTheCandidateDiff"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReviewSessionReadsWithoutWriting" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewSessionReadsWithoutWriting"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — What the reviewer is handed

## Result

Implemented candidate-diff collection and a reviewer prompt that embeds the
computed base-to-head diff. Review Agent Sessions now carry an explicit
read-only access mode: reads are approved, while write permission requests are
denied when no interactive prompt is available. Existing implementation
sessions retain their read-write zero-value behavior.

Focused-check evidence:

- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1 -run 'TestReview(Record|Prompt|Session)' ./internal/cli` — passed. The real two-commit fixture proved the captured prompt contains the `diff --git` header and the added line `+the reviewer receives this changed line`.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1 -run '^TestACPXPromptArgsPlaceGlobalsBeforeAgentAndSubcommand$' ./internal/agent` — passed. The read-only request produced `--approve-reads --non-interactive-permissions deny`; normal work requests retained `--approve-all`.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache make verify-incremental` — the sandboxed attempt reached the Go suite and failed only where two existing force-stop integration tests could not inspect their spawned process trees. The narrowly host-permitted rerun passed all packages, skill checks, and the build.

Acceptance evidence:

- The fixture-candidate prompt contains the computed diff text, including its file header.
- The prompt assertion names the fixture's added line, so commit identifiers alone cannot satisfy it.
- The captured runner request reports readable and non-writable access, and the ACPX argument test proves that request denies write approvals at the runtime boundary.
