---
task: task_10
spec: 0001-implement-command
status: completed
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

- [x] Implement Command section in the canonical skill
- [x] Task Batch contract semantics alongside the review Batch contract
- [x] Embedded skill regeneration via the sync target
- [x] Glossary coverage pass over the new user-facing text

## Acceptance Criteria

- [x] The canonical skill describes the Implement Command's flags, output lines, and exit codes exactly as implemented — each documented flag exists in the command's help output.
- [x] The skill's task Batch contract names the four task statuses and the never-commit/never-push rules.
- [x] The embedded copy matches the canonical skill (drift check passes inside the verification gate).
- [x] No new un-glossaried term appears in skill or help text.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts validate.
- `make verify` — expected: full gate passes, including the skills drift check.

## References

`_prd.md` → Core Features 1, 9, 10. `_techspec.md` → Build Order 11. Repo hard rule (canonical skill ships with CLI behavior changes). ADR-0013.

## Result

The canonical Roundfix skill now teaches both halves of the Implement
Command surface. A new "User-Facing Spec Runs" section documents the
invocation, all eight flags (`--spec`, `--qa`, `--agent`, `--model`,
`--agent-command`, `--agent-full-access`, `--interactive`, `--no-input`),
the deterministic stdout contract (per-Task lines in Task Graph order, the
`--qa` verdict line including the missing-report form, and the three
outcome lines), the exit codes (0 Clean/no-op, 1 Unresolved/Failed, 2
Preflight Validation, 130 Stop Request), the Preflight Validation checks,
Interactive Input behavior (Active Specs picker, remembered Agent, never-
remembered Spec slug), the never-push rule with the ADR-0013 pull request
decision, and Attach through the Live Run View. A new "Assigned Task
Batches" section documents the task Batch contract alongside the renamed
"Assigned Review Issue Batches" section: each Task is one Batch of one,
the four task statuses (`pending`, `in_progress`, `completed`, `failed`),
the Agent's ownership (in_progress on start, code edits, Verification
while working, `## Result`, settling completed/failed) versus the Daemon's
(verbatim Verification re-run, final settling, one commit per verified
Task with `Roundfix-Spec`/`Roundfix-Task` trailers, the QA Report commit),
and the never-commit/never-push/never-edit-manifest rules, which also
extend Forbidden Actions and the Completion Report.
`.agents/skills/roundfix/agents/openai.yaml` gained the spec command hint
and Task Batch ownership entries since it enumerates commands. The
embedded copy was regenerated with `make skills-sync`; `metadata.version`
stays 0.1.0 because no `v*` release tag exists yet.

Verification evidence (fresh, this session):

- `rtk go run ./cmd/roundfix implement --help` — all eight documented
  flags present; skill wording cross-checked against the help text and
  `internal/cli/implement.go` output printers.
- `rtk go run ./cmd/roundfix skills check` — passed
  ("Roundfix skill check passed: roundfix").
- `make verify` — passed: fmt-check, 406 tests in 16 packages, the
  skills-sync drift check, skills check, and build.

Glossary gap (Requirement 4): CONTEXT.md has no entry for **Clean**, the
terminal Run outcome the CLI prints (`Clean: all N Task(s) completed.`)
and the help text names ("lets the Run end Clean"). Unresolved Outcome is
glossaried; its counterpart is not. The skill uses "Clean" exactly as the
shipped CLI does rather than inventing a definition; adding a Clean entry
to CONTEXT.md is a follow-up for the glossary owner. Every other term in
the new text maps to an existing CONTEXT.md entry.
