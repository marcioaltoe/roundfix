---
task: task_11
spec: 0027-review-loop-integrity
status: completed
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

- [x] Rewrite the affected Roundfix Skill sections against the shipped behavior of the preceding tasks
- [x] Update the OpenAI manifest and the skill-check anchors together
- [x] Update glossary entries scoped by this spec
- [x] Review command help text against implemented flags and outcomes
- [x] Cross-read the PRD's Core Feature 10 checklist and confirm each item landed

## Acceptance Criteria

- [x] The skills check passes with the updated anchors
- [x] No skill text still claims review Runs execute in a Run Worktree or that Roundfix never comments on threads
- [x] Glossary entries for Run Worktree, Run Branch, and Integration Pending name spec Runs as their scope
- [x] Help output for watch, resolve, and fetch documents the bypass flag and, for watch, the Clean Unverified outcome
- [x] The full test suite passes

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/agents/openai.yaml`
- interface: `skills/skills.go`
- interface: `CONTEXT.md`

## Verification

- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: skill check passes
- `grep -q "Branch Integrity Preflight" skills/roundfix/SKILL.md` — expected: exit 0
- `grep -c "Run Worktree" skills/roundfix/SKILL.md` — expected: exit 0 (remaining mentions are spec-Run-scoped; reviewed in acceptance criteria)
- `go test ./...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Core Feature 10, Decisions; `_techspec.md` → Build Order 10, Coverage Map (Core Feature 10); ADR-0042, ADR-0043, ADR-0045; `docs/agents/skill-governance.md`.

## Result

Updated the Roundfix Skill, OpenAI manifest, skill-check anchors, glossary, and review command help text to match the shipped review-loop integrity contract:

- Review Runs now document checkout execution, Branch Integrity Preflight, audited `--skip-branch-integrity`, clean tracked checkout validation, no review Integration Pending outcome, CleanUnverified exit `3`, Outcome Comments, and the rule that a Verification Feedback retry does not consume a Round or request a new Review Source review.
- The OpenAI manifest now scopes Review Run, watch outcome, Review Source, and Worktree Bootstrap hints to the current behavior.
- `CONTEXT.md` now scopes Run Worktree, Run Branch, and Integration Pending to spec Runs; Fetch Run and Clean Unverified now describe checkout execution, no worktree creation, and exit code `3`.
- `fetch`, `resolve`, and `watch` help now describe Branch Integrity Preflight and the bypass flag; `watch --help` also documents CleanUnverified and exit `3`.

Verification evidence:

- `rtk go test ./internal/cli -run 'TestRunCommandHelp|TestReviewCommandHelpDocumentsSkipBranchIntegrity'`: passed, 12 tests in `internal/cli`.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check`: passed with updated anchors.
- `rtk go run -buildvcs=false ./cmd/roundfix watch --help`: passed; output includes Branch Integrity Preflight, `--skip-branch-integrity`, CleanUnverified, and exit `3`.
- `rtk go run -buildvcs=false ./cmd/roundfix resolve --help`: passed; output includes Branch Integrity Preflight and `--skip-branch-integrity`.
- `rtk go run -buildvcs=false ./cmd/roundfix fetch --help`: passed; output includes Branch Integrity Preflight and `--skip-branch-integrity`.
- `rtk proxy grep -q "Branch Integrity Preflight" skills/roundfix/SKILL.md`: passed.
- `rtk proxy grep -c "Run Worktree" skills/roundfix/SKILL.md`: passed with count `17`; remaining mentions are explicit no-review-worktree text or spec Run / Settle recovery text.
- `rtk proxy grep -n "execute in a Run Worktree\\|not in the user's checkout\\|treating Run as Clean\\|never comments\\|does not comment" skills/roundfix/SKILL.md`: exited `1`, confirming no stale review-worktree or no-comments claims matched.
- `rtk go test ./...`: passed, 1193 tests in 19 packages.
- `rtk go build -buildvcs=false ./...`: passed.
- `rtk make verify`: passed; it ran the full local gate (`go test ./...`, skills check, and `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`).
