---
task: task_02
spec: 0153-a-reviewer-the-workflow-runs
status: pending
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
