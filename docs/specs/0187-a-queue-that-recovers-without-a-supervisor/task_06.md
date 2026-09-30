---
task: task_06
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: pending
type: docs
complexity: low
---

# Task 06: The grant rule sits outside the shared QA settlement section

## Overview

task_04 described the grant a Task ran under in the Roundfix skill, and placed the paragraph inside the `### QA settlement` section. `TestSettlementGuidanceIsOneTable` requires that section to be byte-identical in the canonical `qa-gate`, `archive-spec` and Roundfix skills, so `make verify-changed` failed and the QA gate refused at its precondition. This corrective Task moves the paragraph out of the shared section and changes nothing else.

## Requirements

1. MUST make the `### QA settlement` section of `.agents/skills/roundfix/SKILL.md` byte-identical to the same section of `.agents/skills/qa-gate/SKILL.md`, from its heading up to the next heading.
2. MUST keep the moved paragraph, with its wording unchanged, under its own heading `### Authorization audit`, placed directly after the QA settlement section. The phrase `the grant the Task ran under` MUST still be present in the skill.
3. MUST bring `skills/roundfix/SKILL.md` to the same bytes with the sanctioned `make skills-sync`.
4. MUST NOT edit any path other than the two Roundfix skill files and this Task file.

## Subtasks

- [ ] Move the grant paragraph under its own heading after the shared section.
- [ ] Regenerate the mirror with the sanctioned command.

## Acceptance Criteria

- [ ] `TestSettlementGuidanceIsOneTable` passes.
- [ ] The Roundfix skill still states the rule for the grant the Task ran under, and its mirror is byte-identical.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^TestSettlementGuidanceIsOneTable$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestSettlementGuidanceIsOneTable" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- 'the grant the Task ran under' && grep -q '^### Authorization audit$' .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task `TestSettlementGuidanceIsOneTable` fails on the Roundfix skill's section.

## References

- task_04
- `_authorization.md`

## Result
