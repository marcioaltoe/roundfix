---
task: task_12
spec: 0057-baseline-capability-evidence-and-retention
status: pending
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

