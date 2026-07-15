---
task: task_08
spec: 0028-settlement-and-reporting
status: completed
type: docs
complexity: medium
---

# Task 08: Sync authoring guidance and the Roundfix Skill

## Overview

Ship the teaching material with the behavior (skill-sync hard rule): the repo-owned write-tasks authoring guidance gains the field-proven Verification rules, and the Roundfix Skill's settle output contract, implement report contract, stop/lock guidance, and task-status vocabulary are updated to what the preceding tasks shipped, keeping the skills check green.

## Requirements

1. MUST add to the repo-owned write-tasks authoring guidance (canonical source under the repo's authorial skills tree, generated mirror refreshed): Verification commands must use portable shell forms (no `wc`-pipeline patterns, `grep` rather than `rg`), the repository's build flags (`go build -buildvcs=false`), and must prove the Task's effect with executable checks.
2. MUST update the Roundfix Skill's settle section for the new `commit <path>` report lines and sibling-failed warning, the implement section for the indented reason lines, the stop/lock guidance for automatic orphaned-lock reclamation on proven owner death, and the Assigned Task Batches status vocabulary for the documented synonym normalization.
3. MUST update the skill-check required-phrase anchors together with any anchored text it changes, keeping the check green.
4. MUST mention the adapter preflight diagnosis wherever the skill describes runtime readiness (setup/doctor sections).
5. MUST NOT alter upstream-managed skills (skill-governance ownership split).

## Subtasks

- [x] Write-tasks authoring guidance: portability, build flags, effect-proving Verification
- [x] Roundfix Skill: settle, implement report, stop/lock, status vocabulary, doctor sections
- [x] Skill-check anchors updated with the text they anchor
- [x] Cross-read the PRD's Core Features and confirm each shipped behavior is described accurately

## Acceptance Criteria

- [x] The skills check passes with the updated anchors
- [x] The skill's settle example includes the `commit <path>` line shape and the implement section documents the reason line
- [x] The skill no longer implies a dead owner always requires a manual force stop
- [x] The write-tasks guidance names the portable-forms and effect-proving rules
- [x] The full test suite passes

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/skills.go`

## Verification

- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: skill check passes
- `grep -q "buildvcs" .agents/skills/write-tasks/SKILL.md` — expected: exit 0 (build-flag rule documented)
- `grep -q "commit <path>" skills/roundfix/SKILL.md` — expected: exit 0 (settle contract updated)
- `go test ./...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Core Feature 7, Goals 3–4; `_techspec.md` → Build Order 8, Risks (report consumers), Coverage Map (Core Feature 7); `docs/agents/skill-governance.md`.

## Result

- Updated the repo-owned write-tasks guidance and task template to require portable Verification forms, `go build -buildvcs=false` for Go build checks, and executable checks that prove the Task's effect.
- Updated the canonical Roundfix Skill and regenerated the embedded `skills/` mirror with settle `commit <path>` lines, implement `reason:` lines, orphan-lock reclamation guidance, status synonym normalization, no-op Task commit warnings, and adapter readiness diagnostics.
- Added required phrase anchors in `skills/skills.go` for the new Roundfix Skill contracts.
- Acceptance evidence:
  - `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed with the updated anchors.
  - `rtk grep -q "buildvcs" .agents/skills/write-tasks/SKILL.md` passed.
  - `rtk grep -q "commit <path>" skills/roundfix/SKILL.md` passed.
  - `rtk go test ./...` passed, 1235 tests.
  - `rtk go build -buildvcs=false ./...` passed.
  - `rtk make verify` passed: `go test ./...` 1235 tests, skill check passed, build passed.
