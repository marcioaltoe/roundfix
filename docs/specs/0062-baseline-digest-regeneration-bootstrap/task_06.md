---
task: task_06
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
type: docs
complexity: low
---

# Task 06: Document the regeneration contract

## Overview

A mode that relaxes validation is safe only while the next reader understands
its boundary. This Task records the contract where this repository keeps the
working instructions for its own delivery loop: when regeneration mode applies,
what it defers and why that is safe, and the one step a new Normative Clause
still requires by hand.

## Requirements

1. MUST document that regeneration mode is reachable only from the regeneration
   path and that every other load — production, CLI, CI, and the Verification
   gate — stays strict.
2. MUST document what the mode defers, stated as an enumerable allowlist rather
   than a category, and why deferral is safe: the regeneration target
   re-validates strictly afterwards, so a deferred pin is checked later.
3. MUST document the manual step a new Normative Clause requires — its Source
   Baseline manifest row must be added by hand first, because the regenerator
   maintains rows but never creates them.
4. MUST state the ordinary recovery for a maintainer who hits the cycle: what
   the failure looks like and what to run.
5. MUST NOT restate the ADR's reasoning at length; reference the decision and
   document the operating contract.
6. MUST change no Go source, no workflow, and no `Makefile`.

## Subtasks

- [ ] Document the mode's boundary and what stays strict.
- [ ] Document the deferred allowlist and the strict re-validation that makes
      it safe.
- [ ] Document the manual manifest-row step for a new clause.
- [ ] Document the maintainer's recovery path.

## Acceptance Criteria

- [ ] The document states that regeneration mode is unreachable outside the
      regeneration path.
- [ ] The document names what is deferred and states that a strict
      re-validation closes every regeneration run.
- [ ] The document states that a new clause needs its manifest row added by
      hand first.
- [ ] The document references ADR-0085 rather than restating it.
- [ ] `git status --porcelain` shows no path outside `docs/workflow/` and this
      task file.

## Context

- instruction: `docs/agents/docs-layout.md`
- interface: `docs/workflow/spec-implementation-loop.md`

## Verification

- `test -f docs/workflow/baseline-digest-regeneration.md` — expected: exit 0.
- `grep -qi "strict" docs/workflow/baseline-digest-regeneration.md` — expected:
  exit 0; the boundary between deferred and strict is stated.
- `grep -q "ADR-0085" docs/workflow/baseline-digest-regeneration.md` —
  expected: exit 0; the decision is referenced.
- `grep -qi "manifest row" docs/workflow/baseline-digest-regeneration.md` —
  expected: exit 0; the manual step is documented.
- `grep -qi "baseline-digests" docs/workflow/baseline-digest-regeneration.md` —
  expected: exit 0; the recovery names the command.
- `git diff --name-only HEAD -- Makefile .github/ internal/ | grep -q . && exit 1 || exit 0`
  — expected: exit 0; this task changed no code, workflow, or tooling.

## References

- `_prd.md` → Goals; Core Features 1 and 3.
- `_techspec.md` → Build Order 6; Decisions.
- ADR-0085.
