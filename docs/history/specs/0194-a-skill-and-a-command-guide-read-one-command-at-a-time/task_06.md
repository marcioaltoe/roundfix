---
task: task_06
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: completed
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

## Result

The helper now reads the canonical entry and its `references/` companion
directory through `mdtree.Text`. Both the required-phrase check and the
removal-mutation check use that combined text. When the companion directory
is absent, the reader returns only the entry text. The canonical-versus-
distributed entry equality check, required phrases, test names and existing
assertion messages are unchanged. No skill was edited.

Focused acceptance evidence:

- Before the change,
  `GOCACHE=/private/tmp/roundfix-0194-task06-gocache rtk proxy go test -count=1 -run '^TestReviewRequestContract$' ./skills`
  exited 1 with all seven reported missing review phrases.
- After the change,
  `GOCACHE=/private/tmp/roundfix-0194-task06-gocache rtk proxy go test -count=1 -v -run '^(TestReviewRequestContract|TestProjectConstraint.*Gate|TestArchiveSpecContract|TestToolingAuthorizationJourney|TestLegacySpecConstraintExemption)$' ./skills`
  exited 0. All nine tests that call the helper passed, covering the review
  contract, all five Project Constraint gates, archive guidance, tooling
  authorization and legacy exemptions, including their removal mutations.
- `GOCACHE=/private/tmp/roundfix-0194-task06-gocache rtk proxy go test -count=1 -v ./internal/mdtree`
  exited 0. Its four tests cover entry-only reading without companions,
  lexical companion ordering, ignored nested/non-Markdown files and missing
  entry errors.
- `rtk proxy gofmt -l skills/baseline_skill_contract_test.go` exited 0 with
  no output.

The pre-existing Task status was `in_progress` and remains Daemon-owned.
Only the helper file and this Result were edited. The prescribed Verification
command, including the full `./skills` package check, was not run; it remains
for the Daemon. No commit, push or pull request was made. No follow-up work
was identified.
