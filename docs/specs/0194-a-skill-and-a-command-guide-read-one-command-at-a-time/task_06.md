---
task: task_06
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: pending
type: test
complexity: low
---

# Task 06: The workflow contract reads a skill with its references

## Overview

task_01 taught the documentation contracts to read an entry file together with its companion files. One helper was missed: `testWorkflowProjectConstraintContract` in `skills/baseline_skill_contract_test.go` still reads only `SKILL.md`. task_02 then moved the review-request guidance of the Roundfix skill into `.agents/skills/roundfix/references/review-runs.md`. As a result, the first delivery Run's QA gate failed `TestReviewRequestContract` with "roundfix skill missing …" for all seven required phrases (Run `run_20261001T010622Z_8bfe4bebad66f971`, `batch-005-attempt-1.log`). task_02's Result had already pointed at this helper.

## Requirements

1. MUST make `testWorkflowProjectConstraintContract` read a skill that has a `references/` directory through `mdtree.Text`, with that directory as the companion, so that both the required-phrase check and the removal-mutation check run over the entry file plus its references. A skill without `references/` MUST keep today's single-file read.
2. MUST keep the canonical-versus-distributed `SKILL.md` equality check unchanged.
3. MUST NOT change any required phrase, test name or assertion message, and MUST NOT edit any skill or other file.

## Subtasks

- [ ] Read the skill with its references in the helper.
- [ ] Prove the whole `./skills` package passes.

## Acceptance Criteria

- [ ] `TestReviewRequestContract` and every other test that uses the helper pass.

## Context

- interface: `skills/baseline_skill_contract_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReviewRequestContract|TestProjectConstraintPRDGate|TestProjectConstraintTechSpecGate)$' ./skills 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestReviewRequestContract TestProjectConstraintPRDGate TestProjectConstraintTechSpecGate; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 ./skills` — expected: exit 0; before this Task `TestReviewRequestContract` fails with "roundfix skill missing".

## References

- task_01, task_02
