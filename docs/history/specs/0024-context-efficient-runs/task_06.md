---
task: task_06
spec: 0024-context-efficient-runs
status: completed
type: docs
complexity: medium
---

# Task 06: Ship context-efficient Run guidance

## Overview

Ship operator guidance and workflow Skills that explain which compact surface
to use and where lossless evidence remains. The slice is verifiable through
JSONL examples, canonical-to-embedded Skill sync, and removal of instructions
that tell Agents to run the authoritative Verification themselves.

## Requirements

1. MUST document events replay, follow, filters, categories, JSONL fields, terminal behavior, and stdout/stderr boundaries with copy-paste examples.
2. MUST document the one-repair Verification Feedback lifecycle, failure artifact location, final settlement, and missing-prerequisite behavior.
3. MUST document Console Log summaries, Spec Context Bundle limits, and the lossless Run Event Journal boundary.
4. MUST update the canonical Roundfix Skill to match every new command, flag, exit code, event field, and operational recovery path.
5. MUST update owned write-tasks and implement-task guidance for labeled `## Context` entries and Daemon-owned Verification.
6. MUST regenerate embedded Skills from canonical sources and preserve Skill governance.
7. MUST update autonomous Supervisor guidance to monitor `roundfix events` instead of grepping Console Log text.

## Subtasks

- [x] Add user documentation and valid JSONL command recipes.
- [x] Document Verification Feedback and evidence retention boundaries.
- [x] Update the canonical Roundfix Skill command contract.
- [x] Update owned task-authoring and execution Skills.
- [x] Update Supervisor/autonomous Run monitoring guidance.
- [x] Regenerate embedded Skills and audit glossary vocabulary.
- [x] Verify examples and Skill sync checks.

## Acceptance Criteria

- [x] Documented replay, follow, and filtered-subset commands match CLI help and every sample output line parses as JSON.
- [x] The four Supervisor categories and their stable fields are documented without exposing internal Run Event kinds as filters.
- [x] Guidance states that successful Verification output never enters Agent context and that only one repair is permitted.
- [x] Console Log and Run Event Journal documentation clearly distinguishes compact presentation from lossless payload evidence.
- [x] Task-authoring guidance defines labeled Context entries and the 50-entry/200-path bounds.
- [x] Supervisor guidance uses the Run Event Stream and no longer recommends polling/grepping Console Log.
- [x] Canonical and embedded Skills have zero drift and describe the shipped behavior exactly.

## Verification

- `rtk make skills-sync-check` - expected: canonical and embedded Skill bundles have zero drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` - expected: every shipped Skill passes validation.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/tech-writer/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/skill-governance.md`
- interface: `README.md`

## References

`_prd.md` -> User Story 7; Core Feature 12; User Experience; Success Metrics. `_techspec.md` -> API Contracts; Integration Points; Build Order 6. ADR-0008; ADR-0038.

## Result

- Replay, follow, and filtered-subset commands are documented in `README.md` and the canonical/embedded Roundfix Skill. Fresh evidence: `rtk go run -buildvcs=false ./cmd/roundfix events --help` showed `roundfix events <run-id> [--follow] [--filter <categories>]`; a JSON parsing check loaded all 9 documented `roundfix-events/v1` example lines.
- The four public Supervisor categories and stable fields are documented as `task-status`, `batch`, `verification`, and `outcome`, with no internal Run Event kinds exposed as filters.
- Verification guidance now states that successful output never enters Agent context, attempt-1 command failure sends one same-session Verification Feedback prompt with a diagnostic path, and final settlement follows the attempt-2 verdict without a third attempt.
- Console Log guidance now documents compact read/edit summaries while the Run Event Journal remains the lossless raw ACP payload boundary.
- Task-authoring guidance now defines optional labeled `## Context` entries, clean repository-relative paths, the 50-entry explicit bound, and the 200-path Spec Context Bundle bound.
- Supervisor/autonomous guidance now uses `roundfix events <run-id> --follow` for unattended monitoring and identifies the Console Log as a compact text record, not a state API.
- Embedded Skills were regenerated with `rtk make skills-sync`, and `rtk make skills-sync-check` confirmed zero drift.

Verification:

- `rtk make skills-sync-check` passed.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed after rerunning with build-cache access; every shipped Skill validated.
- `rtk make verify` passed: `go test ./...` reported 1097 tests in 19 packages, the Roundfix skill check passed, and the build completed.
