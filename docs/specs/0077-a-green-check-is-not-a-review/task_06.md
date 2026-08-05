---
task: task_06
spec: 0077-a-green-check-is-not-a-review
status: pending
type: backend
complexity: low
---

# Task 06: Make the watch help text teach the closed default

## Overview

Corrective Task from the QA gate's F-001 (`Trust-Damage`). The runtime gate
closed correctly and the Roundfix Skill teaches the new contract, but the built
`watch --help` still says Clean accepts "a successful CodeRabbit check or commit
status", which is exactly the permissive rule this Spec removed. An operator
reads the help before deciding whether `--until-clean` provides a review
guarantee, so the shipped text contradicting the shipped behaviour is the
documented failure mode of this Spec, repeated in the CLI surface.

This slice changes help text and its assertion only. It changes no behaviour.

## Requirements

1. MUST make the `watch` usage text state that accepted Evidence requires a
   recognised review-completed current-head CodeRabbit check or commit status,
   or a current-head CodeRabbit `APPROVED` review, with zero unresolved threads.
2. MUST state that an unrecognised signal resolves `pending` even when its
   check conclusion is success — a green check is not evidence a review ran.
3. MUST state that an explicit Review Source refusal resolves `skipped` and is
   never merged or cleared for merge.
4. MUST agree with `.agents/skills/roundfix/SKILL.md`, which already carries
   this contract; the help text is aligned to the Skill, not the reverse.
5. MUST update the `watch` case in `TestRunCommandHelp` so the assertion pins
   the new contract instead of the removed one.
6. MUST NOT change classification, watch, event, or any other behaviour, and
   MUST NOT change the Skill or its mirror — task_04 already shipped those.

## Subtasks

- [ ] Rewrite the `watch` Behavior paragraph to the closed-default contract.
- [ ] Pin the new wording in `TestRunCommandHelp` and drop the removed phrase.

## Acceptance Criteria

- [ ] `roundfix watch --help` describes recognised review-completed evidence as
      the only check-or-status route to a verified head.
- [ ] `roundfix watch --help` no longer presents a successful check conclusion
      alone as accepted review evidence.
- [ ] `roundfix watch --help` names the `pending` default for an unrecognised
      signal and the `skipped` outcome for a refusal.
- [ ] `TestRunCommandHelp` asserts the new wording.
- [ ] No Go file outside `internal/cli/` changed, and no Skill file changed.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1 -run TestRunCommandHelp -v | grep -q -- "--- PASS"`
  — expected: exit 0; the help assertion ran and passed.
- `go run -buildvcs=false ./cmd/roundfix watch --help | grep -q "review-completed"`
  — expected: exit 0; the built help names the recognised review-completed rule.
- `go run -buildvcs=false ./cmd/roundfix watch --help | grep -q "a successful CodeRabbit check or commit status" && exit 1 || exit 0`
  — expected: exit 0; the removed permissive phrase is gone from the built help.
- `git diff --name-only HEAD | grep -vE "^(internal/cli/cli\.go|internal/cli/cli_test\.go|docs/specs/0077-a-green-check-is-not-a-review/task_06\.md)$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `qa/qa-report-2026-08-05.md` → F-001; row 18.
- `_prd.md` → Goals 1 and 3; Core Features 2 and 4.
- `_techspec.md` → API Contracts.
