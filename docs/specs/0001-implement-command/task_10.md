---
task: task_10
spec: 0001-implement-command
status: pending
type: docs
complexity: low
---

# Task 10: Update the Roundfix skill and docs for the Implement Command

## Overview

Bring the shipped Roundfix skill and agent-facing docs in line with the new CLI surface so the repo's hard rule holds: the canonical skill matches shipped behavior in the same PR that changes it. Verifiable through the skills drift check inside the full verification gate.

## Requirements

1. MUST document the Implement Command in the canonical Roundfix skill: invocation, flags (including `--qa`), stdout contract, exit codes, and the interactive/non-interactive behavior.
2. MUST document the task Batch contract for assigned spec work: the Agent owns code edits, task status (`in_progress`, then `completed` or `failed`), and the `## Result` section; the Daemon owns verification, settling, and commits; the Agent never commits, pushes, or edits the Task Graph manifest — consistent with the existing review Batch contract's ownership split.
3. MUST regenerate the embedded skill copy through the sync target so the drift check passes.
4. MUST verify the glossary already covers every term the new surface uses and call out any gap instead of inventing language.
5. SHOULD note in the skill that spec Runs never push and that handing the branch to a pull request is the developer's explicit decision (ADR-0013).

## Subtasks

- [ ] Implement Command section in the canonical skill
- [ ] Task Batch contract semantics alongside the review Batch contract
- [ ] Embedded skill regeneration via the sync target
- [ ] Glossary coverage pass over the new user-facing text

## Acceptance Criteria

- [ ] The canonical skill describes the Implement Command's flags, output lines, and exit codes exactly as implemented — each documented flag exists in the command's help output.
- [ ] The skill's task Batch contract names the four task statuses and the never-commit/never-push rules.
- [ ] The embedded copy matches the canonical skill (drift check passes inside the verification gate).
- [ ] No new un-glossaried term appears in skill or help text.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts validate.
- `make verify` — expected: full gate passes, including the skills drift check.

## References

`_prd.md` → Core Features 1, 9, 10. `_techspec.md` → Build Order 11. Repo hard rule (canonical skill ships with CLI behavior changes). ADR-0013.
