---
task: task_04
spec: 0169-override-and-ordinal-follow-ups
status: completed
type: backend
complexity: low
---

# Task 04: One instruction for overrides, one finding per conflict

## Overview

Corrective Task from the pre-PR review of 2026-09-25. The archive-spec skill's Precondition 2 still says to record the override in the stamped frontmatter (`qa_override: true`), contradicting the new Steps; the override field lists in the archive-spec, qa-gate and Roundfix skills omit `qa_override_qa_task_status`. And a same-Spec duplicate beside a fulfilled claim is reported twice, with a message and fix written for two different Specs.

## Requirements

1. MUST reword every archive-spec passage that describes hand-stamping, including Precondition 2 and the paragraph opening with `qa_override: true`, to say the override is performed only through `roundfix archive <slug> --qa-override --approval <source> --reason <text>`.
2. MUST add `qa_override_qa_task_status` to the override field lists in the archive-spec, qa-gate and Roundfix skills, and regenerate the mirrors with `make skills-sync`.
3. MUST skip the same-Spec pair when either path is already held on the tree, and give a same-Spec duplicate its own message and a fix that says to renumber one Task's `creates:` path.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] No skill instructs hand-stamping; every field list names the new field.
- [ ] A same-Spec duplicate beside a fulfilled claim yields exactly one finding, whose fix names renumbering.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/ordinal.go`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce|TestSameSpecDuplicateFixNamesRenumbering)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce TestSameSpecDuplicateFixNamesRenumbering; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && ! grep -q "record it in the stamped frontmatter" .agents/skills/archive-spec/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/archive-spec/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/qa-gate/SKILL.md && grep -q "qa_override_qa_task_status" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/archive-spec skills/archive-spec >/dev/null && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task the tests do not exist and Precondition 2 still hand-stamps, so the command fails.

## References

- [_techspec.md](_techspec.md) — The stamp and the guidance

## Result

Implementation:

- The archive guidance now routes every QA Archive Override through `roundfix
  archive <slug> --qa-override --approval <source> --reason <text>` and no
  longer describes editing override frontmatter as the action. The archive,
  QA-gate, and Roundfix override field lists now include
  `qa_override_qa_task_status`; `make skills-sync` regenerated all three
  mirrors.
- The ordinal checker now loads held ADR paths before comparing claims from one
  Spec. It omits the redundant same-Spec pair when either exact claim is already
  fulfilled on the tree, while the remaining unfulfilled claim still produces
  the tree collision. A same-Spec-only conflict now has its own summary and
  tells the author to renumber one Task's `creates:` path.
- Added `TestSameSpecDuplicateBesideAFulfilledClaimIsReportedOnce` and
  `TestSameSpecDuplicateFixNamesRenumbering` at the public stage-check boundary.

Focused checks:

- Before the implementation, each new test failed independently with the
  expected regression: the fulfilled-claim fixture produced two
  `SC-ORDINAL-CLAIMED` findings, and the duplicate fix said to give one active
  Spec an ordinal instead of naming renumbering.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run
  'TestSameSpecDuplicate(BesideAFulfilledClaimIsReportedOnce|FixNamesRenumbering)$'
  ./internal/speccheck` — exit 0.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1
  ./internal/speccheck` — exit 0.
- `rtk make skills-sync` — exit 0; mirrors regenerated.
- `rtk make skills-sync-check` — exit 0.
- `rtk make baseline-digests` — exit 0 and reported `changed:false`; derived
  artifacts already matched the canonical skill sources.
- `rtk rg -n "record it in the stamped frontmatter"
  .agents/skills/archive-spec/SKILL.md` — exit 1 with no matches, the expected
  absence result.
- `rtk rg -l "qa_override_qa_task_status"
  .agents/skills/archive-spec/SKILL.md .agents/skills/qa-gate/SKILL.md
  .agents/skills/roundfix/SKILL.md` — exit 0 and named all three canonical
  skills.
- `rtk git diff --check` — exit 0.

Acceptance evidence:

- No skill instructs hand-stamping; every field list names the new field: the
  no-match search, three-file field search, and successful skill synchronization
  check above cover the canonical skills and generated mirrors.
- A same-Spec duplicate beside a fulfilled claim yields exactly one finding,
  whose same-Spec message and fix name renumbering: both named regression tests
  and the complete `internal/speccheck` package check passed after the last code
  edit.

The Daemon-owned command under `## Verification` was not run.
