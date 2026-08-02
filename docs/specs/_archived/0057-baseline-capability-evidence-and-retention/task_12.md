---
task: task_12
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: docs
complexity: medium
---

# Task 12: Document the evidence and retention contract

## Overview

This Spec changes what a maintainer sees at nearly every step of a Baseline
run, and adds a command and a prompt outcome that nothing else announces. This
Task records the contract where the human documentation for Baseline already
lives.

## Requirements

1. MUST document the read-only capability re-check: what it does, that it needs
   no decisions, that it writes nothing, and when to use it.
2. MUST document the four prompt outcomes, including the remediate-and-re-run
   outcome and the adaptation option's removal-only constraint.
3. MUST document that a Profile or catalog digest change under an unchanged
   Baseline identifier requires retention accounting, and that an unaccounted
   clause stops planning.
4. MUST document the five status-matrix axes and state that completion language
   requires verified retention and a passing idempotence check.
5. MUST document that executable discovery resolves symlink chains without
   executing, and what each rejection reason means.
6. MUST document the requirement-strength grouping and that an advisory never
   blocks readiness or apply.
7. MUST NOT change any Go source or tooling file.

## Subtasks

- [ ] Document the re-check command.
- [ ] Document the four prompt outcomes.
- [ ] Document retention accounting under same-identity drift.
- [ ] Document the status matrix and its completion rule.
- [ ] Document executable discovery and its rejection reasons.

## Acceptance Criteria

- [ ] The documentation describes the re-check, its zero-decision requirement,
      and that it writes nothing.
- [ ] It lists all four prompt outcomes and the adaptation removal-only
      constraint.
- [ ] It states that same-identity digest drift requires retention accounting
      and that an unaccounted clause stops planning.
- [ ] It lists the five status axes and the condition for completion language.
- [ ] It states that discovery resolves links without executing, and explains
      the cycle, broken-link, and non-executable reasons.
- [ ] It states that an advisory divergence never blocks readiness or apply.
- [ ] `git status --porcelain` shows no path outside `docs/user-guide/` and
      this task file.

## Context

- instruction: `docs/agents/docs-layout.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -qi "re-check" docs/user-guide/commands.md` — expected: exit 0.
- `grep -qi "retention" docs/user-guide/commands.md` — expected: exit 0.
- `grep -qi "idempotence" docs/user-guide/commands.md` — expected: exit 0.
- `grep -qi "symlink" docs/user-guide/commands.md` — expected: exit 0.
- `grep -qi "advisory" docs/user-guide/commands.md` — expected: exit 0.
- `git diff --name-only HEAD -- internal/ Makefile .github/ | grep -q . && exit 1 || exit 0`
  — expected: exit 0; this task changed no code or tooling.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Experience; Decisions.
- `_techspec.md` → API Contracts; Build Order 12.
- ADR-0087.

## Result

### Implementation

- The current Baseline command reference records the public
  `roundfix baseline capabilities check` synopsis and directs maintainers to
  use the exact prompt-provided re-check after repository remediation and
  before returning to planning. It states that the re-check shares the full
  plan's capability outcomes, requires and resolves no decisions, and writes
  no repository, journal, or configuration bytes.
- The reference lists all four blocked-alignment prompt outcomes. It keeps a
  repository-owned Profile adaptation removal-only and distinguishes
  remediate-and-re-run from declining without writing.
- The reference requires clause-level retention accounting when an unchanged
  Baseline identifier has a changed Profile or catalog digest. An unaccounted
  clause stops planning with an action-required result and no apply offer.
- The reference lists approved postimages, semantic retention, Profile
  alignment, repository Verification, and idempotence as five independent
  status axes. Completion language requires verified semantic retention and a
  passing idempotence check.
- The reference documents blocking, advisory, and informational requirement
  strengths; advisory divergence never blocks readiness or apply. It also
  explains bounded symlink resolution without execution and the `link-cycle`,
  `broken-link`, and `not-executable` rejection reasons.
- This fresh Task run found the required user-guide contract already present
  in the current `HEAD`, while the Task file had no Result. No additional
  user-guide wording or Go/tooling change was necessary; this Result records
  fresh criterion-level evidence without changing the Daemon-owned status.

### Focused checks

- Static comparison of `docs/user-guide/commands.md` against the PRD's User
  Experience and Decisions, the TechSpec's API Contracts and Build Order 12,
  ADR-0087, and the current CLI implementation found all required contract
  statements present and consistent with the shipped command surface.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task12-gocache rtk go test
  ./internal/cli -run '^TestBaselineDocumentationContract$' -count=1` exited 0
  with the focused Baseline public-documentation contract test passing.
- `rtk git diff --check` exited 0 with no whitespace errors.
- The commands declared under `## Verification` were not run; the Daemon owns
  that Verification and terminal settlement.

### Acceptance evidence

1. Re-check: the command reference names the command, says to use it after
   remediation and before re-planning, and states both its zero-decision and
   zero-write contracts.
2. Prompt outcomes: the numbered list contains Profile change, removal-only
   Profile adaptation, remediate-and-re-run, and decline.
3. Retention: same-identity Profile or catalog digest drift requires an
   explicit disposition for every previously managed Normative Clause; an
   unaccounted clause stops planning before apply.
4. Status matrix: the reference lists all five axes and gates completion
   language on verified semantic retention plus passed idempotence.
5. Discovery: the reference states that bounded symlink chains resolve without
   execution and defines cycle, broken-link, and non-executable rejection
   reasons.
6. Requirement strength: the reference groups blocking, advisory, and
   informational divergences and states that advisory never blocks readiness
   or apply.
7. Scope: this Agent edited only this Task's `## Result`; the pre-existing
   `status: in_progress` change remains Daemon-owned. No Go source, tooling,
   sibling Task, Task Graph manifest, or other path was edited.
