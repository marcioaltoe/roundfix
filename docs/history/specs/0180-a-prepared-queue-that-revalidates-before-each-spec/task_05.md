---
task: task_05
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: completed
type: docs
complexity: low
---

# Task 05: The implement-spec entry point hands implementation to Roundfix

## Overview

The owned `implement-spec` skill in `.agents/skills/implement-spec/SKILL.md` still tells the Supervisor to plan waves, run the implement-task cycle for every Task and run the qa-gate itself. That is a second implementation loop, and it contradicts `docs/agents/autonomous-work.md`: the Supervisor must not write feature code or tests, and implementation is delegated through a Roundfix Run. This Task rewrites the skill so it prepares with the Delivery Plan from task_03, hands the work to `roundfix implement` or `roundfix deliver`, and asks the maintainer only the Pending Question from task_04. The skill is read by a Supervisor session and shipped in the binary through its `skills/` mirror.

## Requirements

1. MUST rewrite `.agents/skills/implement-spec/SKILL.md` so the Supervisor, in order:
   - runs `roundfix deliver plan <slug>...` and stops on a blocked Spec, reporting its reasons;
   - hands a single Spec on the current branch to `roundfix implement --spec <slug>`, with `--detach` when the session may end, or a merge-through sequence to `roundfix deliver start [--max-duration <duration>] [--max-retries <n>] <slug>...`;
   - monitors through `roundfix deliver status` or the Run's events, and asks the maintainer only the Pending Question.
2. MUST state in the skill that the Supervisor never writes code or tests, never runs a Task or the implement-task cycle itself, and never runs the QA gate itself. The QA gate is the Spec's terminal Task, which the Daemon runs. The skill MUST remove the wave-planning and per-Task loop sections.
3. MUST keep the skill's `name`, `disable-model-invocation: true` and argument hint, update its `description` to the delegation, move both of its version declarations to `0.1.0`, and defer to the Roundfix skill for command details instead of copying them.
4. MUST regenerate `skills/implement-spec/SKILL.md` with `make skills-sync` and run `make baseline-digests` after the edit.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Check each acceptance criterion against the regenerated mirror.

## Acceptance Criteria

- [ ] The skill names `roundfix deliver plan`, `roundfix implement --spec` and `roundfix deliver start`, and tells the Supervisor to ask only the Pending Question.
- [ ] The skill no longer tells the Supervisor to run the implement-task cycle, plan waves or run the qa-gate.
- [ ] The mirror is byte-identical to the canonical skill, and `make skills-sync-check` exits `0`.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `.agents/skills/implement-spec/SKILL.md`
- interface: `skills/implement-spec/SKILL.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "roundfix deliver plan" && tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "roundfix implement --spec" && tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "roundfix deliver start" && tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "Pending Question" && tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "never writes code or tests" && ! { tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "Run the full **implement-task** cycle"; } && ! { tr -s '[:space:]' ' ' < .agents/skills/implement-spec/SKILL.md | grep -qF -- "## 2. Plan the waves"; } && grep -q "^version: 0.1.0$" .agents/skills/implement-spec/SKILL.md && diff -r .agents/skills/implement-spec skills/implement-spec >/dev/null && make skills-sync-check` — expected: exit 0; before this Task the skill carries the implement-task loop and names none of the Roundfix hand-off commands, so the command fails.

## References

- [_techspec.md](_techspec.md) — The implement-spec entry point
- `_prd.md` → Goal 6; Core Feature 6; Success Metric 5
- `_techspec.md` → Testing Approach 5

## Result

Implementation:

- Rewrote the canonical `implement-spec` skill as a Supervisor-only handoff:
  it runs the Delivery Plan, delegates to `roundfix implement --spec` or
  `roundfix deliver start`, monitors the Run or queue, and asks only the
  Pending Question.
- Removed the wave-planning, per-Task `implement-task` loop, and Supervisor-run
  QA instructions. The skill states that the Daemon runs the terminal QA Task
  and defers command details to the Roundfix skill.
- Updated both version declarations to `0.1.0` while preserving the skill name,
  invocation setting, and argument hint.
- Regenerated `skills/implement-spec/SKILL.md` with `make skills-sync` and ran
  `make baseline-digests`; the latter reported `ok: true` and `changed: false`.

Focused checks:

- `rtk git diff --check`: exited 0.
- `rtk cmp -s .agents/skills/implement-spec/SKILL.md skills/implement-spec/SKILL.md`:
  exited 0; the regenerated mirror matches the canonical skill.
- Targeted content scan: exited 0 and found the required delegation commands,
  Pending Question instruction, Supervisor guardrails, and both `0.1.0`
  declarations.
- Negative scan for the removed wave heading, per-Task loop instruction, and
  Supervisor-run QA instruction: exited 1 because none matched, as expected.

Acceptance evidence:

- The canonical skill names `roundfix deliver plan <slug>...`,
  `roundfix implement --spec <slug>`, and `roundfix deliver start ...`, and
  instructs the Supervisor to ask only the Pending Question.
- The canonical skill explicitly prohibits Supervisor code/test authoring,
  Task execution, the `implement-task` cycle, wave planning, and running the
  `qa-gate`.
- The mirror is byte-identical by `cmp`; the Daemon must run the declared
  `make skills-sync-check` Verification command before settlement.
