---
task: task_02
spec: 0138-a-matrix-the-spec-declares
status: pending
type: backend
complexity: low
---

# Task 02: State the declared coverage in the QA contract

## Overview

The Daemon puts a short QA contract in the gate Agent's prompt. Its first bullet
tells the gate to validate every user story and acceptance criterion from the
PRD and the task files. That contradicts the rules this Spec writes into the
qa-gate skill.

This slice rewrites the contract to state five things: the declaration rule, the
bounded default, the non-waivable sources, the coverage rule and scoped
blocking. It also adds a prompt test that pins the result.

## Requirements

1. MUST replace the contract bullet that tells the gate to validate every user
   story and acceptance criterion. The new contract text MUST contain these
   phrases verbatim:
   - `Each numbered qa Task Requirement that starts with MUST verify or MUST run is a declared verification Requirement, and the declared Requirements with the non-waivable sources are the coverage sources`
   - `otherwise the coverage sources are the PRD user stories and Goals, each non-QA Task's Acceptance Criteria and the declared intentional breaks`
   - `The non-waivable sources are the outside-evidence row, the Pull Request row, the repository Verification, each PRD Unreachable Acceptance declaration and, when the PRD declares frontend, the frontend sweep`
   - `Every row names in its provenance the sources it covers; every source appears in at least one row, and no row covers anything else`
2. MUST add a contract bullet containing `Once the matrix exists, a finding
   blocks only the rows that depend on it`.
3. MUST keep the report-naming, verdict and never-commit bullets unchanged, so
   the existing QA prompt tests pass without edits.
4. MUST add a test named `TestBuildQAPromptStatesTheDeclaredMatrix`. The test
   builds the QA prompt and asserts two things: every phrase from Requirements 1
   and 2 is present, and the phrase `validate every user story and acceptance
   criterion` is absent.
5. MUST NOT change the prompt's checkout facts, previous-report identity, Spec
   Context Bundle, or any other prompt builder.
6. MUST NOT use a backtick inside the contract text, because the contract is a
   Go raw string.

## Subtasks

- [ ] Replace the matrix bullet with the declaration, default and non-waivable
      sources.
- [ ] Add the coverage and scoped-blocking bullets.
- [ ] Add the prompt test for the new phrases and the removed one.

## Acceptance Criteria

- [ ] `TestBuildQAPromptStatesTheDeclaredMatrix` passes, and it fails if any
      required phrase is removed or the old bullet returns.
- [ ] The QA contract source no longer contains `validate every user story and
      acceptance criterion`.
- [ ] Every existing QA prompt test passes unchanged.

## Context

- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`

## Verification

- `grep -q 'func TestBuildQAPromptStatesTheDeclaredMatrix' internal/agent/spec_prompt_test.go && go test -count=1 ./internal/agent -run '^TestBuildQAPromptStatesTheDeclaredMatrix$'` — the contract states the declaration rule, the bounded default, the non-waivable sources, the coverage rule and scoped blocking. This fails on today's tree.
- `test -f internal/agent/spec_prompt.go && grep -qF 'starts with MUST verify or MUST run' internal/agent/spec_prompt.go && ! grep -qF 'validate every user story and acceptance criterion' internal/agent/spec_prompt.go` — the declaration rule replaced the old bullet. This fails on today's tree.
- `grep -q 'func TestBuildQAPromptStatesTheDeclaredMatrix' internal/agent/spec_prompt_test.go || exit 1; go test -count=1 ./internal/agent -run '^TestBuildQAPrompt'` — every QA prompt test still passes.

## References

- `_prd.md` → User Stories 1-3; Core Feature 9.
- `_techspec.md` → Implementation Design: QA contract shape; Testing Approach 1;
  Build Order 2.
- ADR-0155.
