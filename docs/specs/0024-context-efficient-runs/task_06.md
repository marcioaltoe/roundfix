---
task: task_06
spec: 0024-context-efficient-runs
status: pending
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

- [ ] Add user documentation and valid JSONL command recipes.
- [ ] Document Verification Feedback and evidence retention boundaries.
- [ ] Update the canonical Roundfix Skill command contract.
- [ ] Update owned task-authoring and execution Skills.
- [ ] Update Supervisor/autonomous Run monitoring guidance.
- [ ] Regenerate embedded Skills and audit glossary vocabulary.
- [ ] Verify examples and Skill sync checks.

## Acceptance Criteria

- [ ] Documented replay, follow, and filtered-subset commands match CLI help and every sample output line parses as JSON.
- [ ] The four Supervisor categories and their stable fields are documented without exposing internal Run Event kinds as filters.
- [ ] Guidance states that successful Verification output never enters Agent context and that only one repair is permitted.
- [ ] Console Log and Run Event Journal documentation clearly distinguishes compact presentation from lossless payload evidence.
- [ ] Task-authoring guidance defines labeled Context entries and the 50-entry/200-path bounds.
- [ ] Supervisor guidance uses the Run Event Stream and no longer recommends polling/grepping Console Log.
- [ ] Canonical and embedded Skills have zero drift and describe the shipped behavior exactly.

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
