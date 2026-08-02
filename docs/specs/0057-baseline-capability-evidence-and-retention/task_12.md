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

- Published `roundfix baseline capabilities check [--profile <id>] [--repo
  <path>] [--format <text|json>]` in the Baseline command synopsis. The
  remediation workflow now tells maintainers to run the exact command printed
  by the prompt before returning to planning and states that the re-check
  requires and resolves no decisions and writes nothing.
- Documented all four blocked-alignment prompt outcomes. Repository-owned
  Profile adaptation is explicitly removal-only, while repository remediation
  exits without writing, prints remediation and the re-check command, and stays
  distinct from decline.
- Documented same-identity Profile or catalog digest drift as a retention
  transition. Every previously managed Normative Clause needs an explicit
  disposition; an unaccounted clause produces an action-required stop and no
  apply offer.
- Documented the five independent result axes: approved postimages, semantic
  retention, Profile alignment, repository Verification, and idempotence.
  Completion language requires verified semantic retention and a passing
  idempotence check.
- Documented bounded symlink-chain resolution without candidate or target
  execution, including distinct `link-cycle`, `broken-link`, and
  `not-executable` rejection meanings.
- Documented blocking, advisory, and informational requirement-strength
  groups, including the rule that an advisory divergence never blocks
  readiness or apply.

### Focused checks

- `rtk proxy env GOCACHE=/private/tmp/roundfix-task12-gocache rtk go test
  ./internal/cli -run '^TestBaselineDocumentationContract$' -count=1` exited 0
  with 13 passing tests in one package.
- `rtk git diff --check` exited 0 with no diagnostics after the user-guide
  change.
- A focused inspection of the published Baseline section confirmed the command
  synopsis, the three requirement-strength groups, all four prompt outcomes,
  the no-decision/no-write re-check contract, all three executable rejection
  reasons, retention accounting, and all five status axes.
- `rtk git -c core.fsmonitor=false status --short` listed only
  `docs/user-guide/commands.md` and this Task file. The Task status change was
  present before this Agent's edits and remains Daemon-owned.

### Acceptance evidence

1. The capability-evidence section names the public re-check, directs its use
   after remediation and before re-planning, and states that it requires and
   resolves no decisions and writes nothing.
2. The numbered prompt list contains Profile change, removal-only Profile
   adaptation, remediate-and-re-run, and decline outcomes.
3. The retention section requires clause-level accounting for an unchanged
   Baseline identifier with a changed Profile or catalog digest and states that
   any unaccounted clause stops planning action-required without apply.
4. The status table contains all five required axes, and its following
   completion rule requires verified semantic retention and passed
   idempotence.
5. The discovery section states that bounded symlink chains resolve without
   execution and defines cycle, broken-link, and non-executable diagnostics.
6. The requirement-strength paragraph states that advisory divergences never
   block readiness or apply.
7. The focused changed-path inspection contains no path outside
   `docs/user-guide/` and this Task file.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
