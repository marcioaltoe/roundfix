---
task: task_07
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: docs
complexity: low
---

# Task 07: The archive-spec skill carries the same QA settlement table

## Overview

`TestSettlementGuidanceIsOneTable` in `skills/settlement_guidance_repocontract_test.go` requires the `### QA settlement` section to be byte-identical in the canonical `qa-gate`, `archive-spec` and Roundfix skills. task_04 updated the qualifying declared `partial` row in `qa-gate` and in the Roundfix skill. The `archive-spec` skill was outside this Spec's grant and kept the old row, so `make verify-changed` fails and the second QA gate refused at its precondition. The maintainer widened the grant to the two `archive-spec` paths on 2026-09-29 (PR #277). This corrective Task copies the section across and nothing else.

## Requirements

1. MUST make the `### QA settlement` section of `.agents/skills/archive-spec/SKILL.md` byte-identical to the same section of `.agents/skills/qa-gate/SKILL.md`, from its heading up to the next heading. No other line of the skill may change except its version field, and that only if the repository contract requires a version change with a content change.
2. MUST bring `skills/archive-spec/SKILL.md` to the same bytes with the sanctioned `make skills-sync`, then run the sanctioned `make baseline-digests`. It MUST name in its Result every file either command rewrote.
3. MUST NOT edit any path other than the two `archive-spec` skill files, the derived files the two sanctioned commands rewrite, and this Task file.
4. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`, because this Spec's QA gate still runs on the v0.20.0 auditor.

## Subtasks

- [ ] Copy the QA settlement section from the canonical qa-gate skill.
- [ ] Regenerate the mirror and the digests with the sanctioned commands.

## Acceptance Criteria

- [ ] `TestSettlementGuidanceIsOneTable` passes.
- [ ] `diff -r .agents/skills/archive-spec skills/archive-spec` reports no difference, and `make skills-sync-check` exits `0`.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null && make skills-sync-check` — expected: exit 0; before this Task `TestSettlementGuidanceIsOneTable` fails on the archive-spec section.

## References

- `_prd.md` → Core Feature 2
- `_authorization.md` → the 2026-09-29 archive-spec grant

## Result
