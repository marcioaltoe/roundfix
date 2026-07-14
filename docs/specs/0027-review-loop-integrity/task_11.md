---
task: task_11
spec: 0027-review-loop-integrity
status: pending
type: docs
complexity: medium
---

# Task 11: Sync the Roundfix Skill, glossary, and command docs to the new contract

## Overview

The skill-sync hard rule: behavior and its teaching material ship together. This task aligns the shipped Roundfix Skill, its OpenAI manifest, the skill-check anchors, the glossary, and command help text with everything this spec changed — the no-worktree review contract, the Branch Integrity Preflight and its audited bypass, Clean Unverified, Outcome Comments, and the retry-versus-Round clarification.

## Requirements

1. MUST update the Roundfix Skill's review-run sections: review Runs execute in the user's checkout (no Run Worktree), the Branch Integrity Preflight contract and its bypass-with-audit, Clean Unverified semantics and exit code, per-issue Outcome Comment behavior, and the statement that a Verification retry never consumes a Round nor counts as a new Review Source review.
2. MUST update the OpenAI agent manifest wherever its command examples or contract statements diverge from the shipped behavior.
3. MUST update the skill-check required-phrase anchors so the check validates the new contract text, keeping the check green.
4. MUST scope the glossary's Run Worktree, Run Branch, and Integration Pending entries to spec Runs, and align Fetch Run and Clean Unverified entries with shipped behavior.
5. MUST make watch, resolve, and fetch help text truthful for the new flag, outcome, and preflight behavior.
6. MUST NOT alter the authorial workflow skills or any upstream-managed skill (skill-governance ownership split).

## Subtasks

- [ ] Rewrite the affected Roundfix Skill sections against the shipped behavior of the preceding tasks
- [ ] Update the OpenAI manifest and the skill-check anchors together
- [ ] Update glossary entries scoped by this spec
- [ ] Review command help text against implemented flags and outcomes
- [ ] Cross-read the PRD's Core Feature 10 checklist and confirm each item landed

## Acceptance Criteria

- [ ] The skills check passes with the updated anchors
- [ ] No skill text still claims review Runs execute in a Run Worktree or that Roundfix never comments on threads
- [ ] Glossary entries for Run Worktree, Run Branch, and Integration Pending name spec Runs as their scope
- [ ] Help output for watch, resolve, and fetch documents the bypass flag and, for watch, the Clean Unverified outcome
- [ ] The full test suite passes

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/agents/openai.yaml`
- interface: `skills/skills.go`
- interface: `CONTEXT.md`

## Verification

- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: skill check passes
- `rg -q "Branch Integrity Preflight" skills/roundfix/SKILL.md` — expected: exit 0
- `rg -c "Run Worktree" skills/roundfix/SKILL.md` — expected: exit 0 (remaining mentions are spec-Run-scoped; reviewed in acceptance criteria)
- `go test ./...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Core Feature 10, Decisions; `_techspec.md` → Build Order 10, Coverage Map (Core Feature 10); ADR-0042, ADR-0043, ADR-0045; `docs/agents/skill-governance.md`.
