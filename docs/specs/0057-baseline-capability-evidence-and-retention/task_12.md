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

- The Baseline command reference now explains requirement-strength grouping,
  the advisory non-blocking contract, and executable discovery's bounded,
  non-executing symlink resolution with each rejection reason.
- It records all four Profile-divergence outcomes. The adaptation option is
  explicitly removal-only, while remediate-and-re-run exits without writing,
  prints per-divergence remediation and the exact re-check command, and remains
  distinct from decline.
- It explains when to use the capability re-check, its shared capability
  outcomes, its zero-decision contract, and its prohibition on repository,
  journal, and configuration writes.
- It documents same-identity Profile or catalog digest drift, clause-level
  retention accounting, the action-required stop for an unaccounted clause,
  all five status axes, and the retention-plus-idempotence completion rule.

### Focused checks

- The Daemon diagnostic artifact for Verification attempt 1 was inspected; it
  exists but contains no log body. The reported failed command identified the
  missing `re-check` documentation as the actionable diagnostic.
- Static inspection of ADR-0087, `resolveExecutableCandidate`,
  `TestSameIdentityDriftRequiresRetention`, `ResultStatusMatrix`, and
  `resultStatusMatrix` grounded the symlink reasons, retention boundary, axis
  names, and `verified`/`not run` vocabulary in repository sources.
- `rtk git -c core.fsmonitor=false diff --check` passed after the documentation
  and Result edits.
- The commands declared under `## Verification` were not rerun; the Daemon
  owns the next configured Verification attempt.

### Acceptance evidence

- Re-check: the command reference states when to use it, that it shares full
  plan capability outcomes, requires and resolves no decisions, and writes no
  repository file, journal entry, or configuration.
- Prompt outcomes: the reference lists four outcomes and states the adaptation
  option's removal-only boundary.
- Retention: the reference binds retention accounting to Profile or catalog
  digest drift under an unchanged Baseline identifier and states that an
  unaccounted clause stops planning before apply.
- Status matrix: the reference lists approved postimages, semantic retention,
  Profile alignment, repository Verification, and idempotence, each as
  `verified` or `not run`; completion language requires verified retention and
  a passing idempotence check.
- Discovery: the reference states that bounded symlink chains resolve without
  execution and explains `link-cycle`, `broken-link`, and `not-executable`.
- Requirement strength: the reference groups blocking, advisory, and
  informational divergences and states that an advisory never blocks readiness
  or apply.
- Scope: this repair changed only `docs/user-guide/commands.md` and this Task
  file. No Go source, tooling file, sibling Task, or Task Graph manifest was
  edited.
