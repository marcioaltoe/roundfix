---
task: task_07
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: completed
type: docs
complexity: low
---

# Task 07: The archive-spec and Roundfix skills carry the same QA settlement table

## Overview

`TestSettlementGuidanceIsOneTable` in `skills/settlement_guidance_repocontract_test.go` requires the `### QA settlement` section to be byte-identical in the canonical `qa-gate`, `archive-spec` and Roundfix skills. task_04 updated the qualifying declared `partial` row in `qa-gate` and in the Roundfix skill. The `archive-spec` skill was outside this Spec's grant and kept the old row, so `make verify-changed` fails and the second QA gate refused at its precondition. The maintainer widened the grant to the two `archive-spec` paths on 2026-09-29 (PR #277). The first attempt aligned `archive-spec` and then showed that task_04 had left the Roundfix skill's qualifying-`partial` row unchanged, adding a paragraph inside the section instead. This corrective Task makes all three sections byte-identical to the canonical `qa-gate` section and nothing else.

## Requirements

1. MUST make the `### QA settlement` section of `.agents/skills/archive-spec/SKILL.md` and of `.agents/skills/roundfix/SKILL.md` byte-identical to the same section of `.agents/skills/qa-gate/SKILL.md`, from its heading up to the next heading. In the Roundfix skill this replaces the old qualifying-`partial` row with the canonical one and removes the paragraph beginning "Settlement also treats the pre-PR Pull Request row as an exception", because the canonical row already states that exception. No other line of either skill may change except a version field, and that only if the repository contract requires a version change with a content change.
2. MUST bring `skills/archive-spec/SKILL.md` and `skills/roundfix/SKILL.md` to the same bytes with the sanctioned `make skills-sync`, then run the sanctioned `make baseline-digests`. It MUST name in its Result every file either command rewrote.
3. MUST NOT edit any path other than the two `archive-spec` and two Roundfix skill files, the derived files the two sanctioned commands rewrite, and this Task file.
4. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`, because this Spec's QA gate still runs on the v0.20.0 auditor.

## Subtasks

- [ ] Copy the QA settlement section from the canonical qa-gate skill into archive-spec and the Roundfix skill.
- [ ] Regenerate the mirror and the digests with the sanctioned commands.

## Acceptance Criteria

- [ ] `TestSettlementGuidanceIsOneTable` passes.
- [ ] `diff -r` reports no difference for the archive-spec and Roundfix skill mirrors, and `make skills-sync-check` exits `0`.
- [ ] The Roundfix skill still contains the phrase `pre-PR Pull Request row`.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- 'pre-PR Pull Request row' && make skills-sync-check` — expected: exit 0; before this Task `TestSettlementGuidanceIsOneTable` fails on the archive-spec section.

## References

- `_prd.md` → Core Feature 2
- `_authorization.md` → the 2026-09-29 archive-spec grant

## Result

Implemented the canonical QA settlement table in the `archive-spec` and
Roundfix `.agents` skills. The stale qualifying-`partial` wording was replaced
with the canonical pre-PR Pull Request exception, and the duplicate Roundfix
paragraph was removed. No other source-skill lines were changed.

Focused implementation evidence:

- Extracted `### QA settlement` sections from `qa-gate`, `archive-spec`, and
  Roundfix and compared their bytes: all three comparisons were identical.
- Confirmed the Roundfix skill still contains `pre-PR Pull Request row`.
- `git diff --check`: passed.
- `make skills-sync`: exited `0`; rewrote `skills/archive-spec/SKILL.md` and
  `skills/roundfix/SKILL.md`.
- `make baseline-digests`: exited `0`; rewrote no files and reported
  `changed: false`.
- Post-regeneration changed paths were limited to the two `.agents` source
  skills, the two `skills/` mirrors, and this Task file.

Acceptance evidence:

- `TestSettlementGuidanceIsOneTable`: not run because it is part of the
  Daemon-owned Task Verification; the byte comparison above is the focused
  implementation check.
- Mirror regeneration: `make skills-sync` passed and named both rewritten
  mirror files; `make baseline-digests` passed with no derived rewrites.
  `make skills-sync-check` remains for the Daemon's declared Verification.
- The Roundfix phrase `pre-PR Pull Request row` remains present.

The Daemon must run the declared Verification and settle the Task; Task status
was left unchanged.

## Carry-forward provenance

- Source Run: `run_20260929T200238Z_ca4a5877c9dc9d5c`
- Source commit: `56c29f7ecec3bf785489f0bce8fbf1593ce7330a`
