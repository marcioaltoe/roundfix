---
schema: spec-tasks/v1
spec: 0406-claude-md-is-always-a-link-to-agents-md
requires: [0399-adopters-keep-their-own-rules-in-the-specific-repository-guide, 0405-authoring-model-fallback]
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — CLAUDE.md is always a symbolic link to AGENTS.md

| id      | title                                                                                                      | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The Baseline reads the shape of the root CLAUDE.md and the passages it holds                               | backend | high       | —                |
| task_02 | Every Baseline plan makes CLAUDE.md the link to AGENTS.md when nothing is lost, and the guide says so      | backend | high       | task_01          |
| task_03 | Doctor and the update report the Claude Instruction Link, and the guides and skills describe it            | backend | high       | task_01, task_02 |
| task_04 | Run the final QA gate                                                                                      | qa      | high       | task_03          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

The graph follows `_techspec.md` → Build Order. Task 01 adds the reader that both the plan and Doctor use, and narrows the detector. Task 02 needs it for the plan; it also carries the guide sentence and its regeneration, because the sentence alone would not earn a verification cycle and the regeneration rewrites the same plan goldens the link changes. Task 03 needs the reader for Doctor and the plan behaviour of Task 02 for the update tests and the Surface Transcript. Tasks 01 and 02 both write `CONTEXT.md`, and the edge between them already runs them in series. Each implementation Task runs the suites that consume its changed contracts in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, the shape of the root `CLAUDE.md`, in three implementation Tasks and the QA gate. Each Task earns its verification cycle: a pure reader with a table of every shape, the plan and transaction change with its regenerated records, and two command surfaces with their guides and skills.

## What the Tasks stage settled beyond the TechSpec

Each item was measured or derived on the planning branch, based on `eb888135`; the measurements are in [`references/2026-11-18-claude-link-inventory-and-measurements.md`](references/2026-11-18-claude-link-inventory-and-measurements.md).

- **The derived set of task 02 was measured** on a scratch worktree, removed afterwards. The guide sentence, its two version raises, the module record command, `make baseline-digests` and the sanctioned managed refresh changed exactly the files task 02 declares, and a second refresh reported `current`. The link in the plan rewrote the same four plan goldens again and nothing else. No Setup Snapshot, Source Baseline, other profile, formatter golden `AGENTS.md`, `AGENTS.md`, `CLAUDE.md` or `docs/agents/specific-repository.md` changed, so `_authorization.md` needs no amendment on that base. Spec 0405 merges first and regenerates some of the same files; task 02 regenerates from the tree it runs on.
- **Existing tests the change breaks.** With the prototype, the whole `./internal/baseline` package failed only in `TestApplyExactDigest`, through its helper `assertAppliedPostimage`, which task 02 changes. The whole `./internal/cli` package passed (223 s). The Doctor line, which the prototype did not add, moves the two table tests of `internal/cli/doctor_test.go` that count sixteen lines, its check-name list and the five expected outputs that hold the `baseline-rules:` line; task 03 changes exactly those.
- **This repository plans nothing for `CLAUDE.md`.** Its `CLAUDE.md` is already the link, so the Surface Transcript of task 03 holds after task 02's managed refresh.
- **Helper names.** The `internal/baseline` tests already define `assertClaudeLink` and `keepClaudeLinks` for the Claude Skill Links; task 01 forbids reusing them.
- **The place of the Doctor line.** The `doctor` guide orders `baseline-rules:` after `skills:`, `user-skills:` and `authoring-models:` and before the residue lines. `baseline-claude-link:` comes right after `baseline-rules:`, which keeps that order.
- **Exact text.** Tasks 02 and 03 fix the words of the guide sentence, the warnings, the Doctor lines and the skill sentence from `_techspec.md` → API Contracts 2 to 5, so the QA gate searches for fixed strings.

## Time in tests

No test of this Spec reads the clock, compares a time it seeded with the real clock, sleeps or waits on a timer. The reader, the plan and Doctor read files and, for one path, Git's index.
