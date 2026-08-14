---
task: task_04
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
type: backend
complexity: medium
---

# Task 04: Protect tooling authority

## Overview

Ship tooling authority as one universal Normative Clause in the core Baseline.
No Profile, setup preference, narrower guide, or generic implementation request
can authorize or weaken the protected mutation surface.

## Requirements

1. MUST cover creation, editing, renaming, moving, and deletion across every
   confirmed tooling category.
2. MUST cover configurations, scripts, ignore files, plugin declarations, and
   version pins.
3. MUST include the clause unconditionally in every maintained Profile.
4. MUST expose no decision, default, prompt, or exclusion effect for the
   clause.
5. MUST reject any catalog declaration that allows a narrower artifact to
   weaken the clause.
6. MUST include the clause in source accounting and retention contracts.

## Subtasks

- [x] Add the universal tooling-authority rule.
- [x] Render it in the core agent guide.
- [x] Add catalog non-excludability validation.
- [x] Extend source and retention accounting.
- [x] Add mutation and no-prompt tests.

## Acceptance Criteria

- [x] Every maintained Profile renders the complete tooling-authority clause.
- [x] Baseline human setup adds no tooling preference prompt.
- [x] A Profile effect cannot exclude or weaken the clause.
- [x] Every protected tooling artifact category appears in the operative text.
- [x] Removing the clause or its accounting makes catalog validation fail.

## Context

- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`
- interface: `internal/baseline/assets/coverage.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestToolingAuthorityClause|TestToolingAuthorityCannotBeDisabled|TestToolingAuthorityNoPrompt|TestToolingAuthorityAccounting'` — expected: completeness, universality, no-prompt, mutation, and retention cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 3–4; User Story 5; Core Features 10–11 and 13.
- `_techspec.md` → Implementation Design: Data Models; Build Order 4.
- ADR-0077 → readable authorization contract.

## Result

### Behavior delivered

- Added `rule.core.tooling-authority` and its prohibited Normative Clause to the
  core Baseline, the core agent guide, and every maintained Profile without
  introducing a Decision Value, default, setup prompt, or opt-out.
- Protected creation, editing, renaming, moving, and deletion for linter,
  formatter, typechecker, test-runner, architecture-checker, build-tool,
  package-manager, code-generator, and other repository-tooling
  configurations, scripts, ignore files, plugin declarations, and version
  pins.
- Added current-catalog validation that rejects clause removal or mutation,
  missing Profile and guide carriers, conditional or excluding effects,
  alternate narrower templates, decision bindings, and missing retention or
  source accounting. Legacy catalog generations remain compatible.
- Added the clause to both maintained upgrade-retention transitions and the
  maintained source Baseline manifest, corpus, index, preservation checks, and
  normalized catalog fixtures.

### Acceptance evidence

1. `TestToolingAuthorityClause` rendered every maintained Profile and matched
   the complete operative clause.
2. `TestToolingAuthorityNoPrompt` completed human setup for every Profile and
   found no tooling-authority answer or prompt.
3. `TestToolingAuthorityCannotBeDisabled` rejected missing Profile ownership,
   guide exclusion, and conditional core activation.
4. `TestToolingAuthorityClause` asserted every protected mutation verb,
   tooling category, and artifact category in rendered output.
5. `TestToolingAuthorityAccounting` removed source and retention accounting in
   mutation catalogs and observed the required validation diagnostics;
   `TestToolingAuthorityCannotBeDisabled` covered operative-clause protection.

### Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestToolingAuthorityClause|TestToolingAuthorityCannotBeDisabled|TestToolingAuthorityNoPrompt|TestToolingAuthorityAccounting'`
  passed 15 tests in two packages.
- `rtk make verify` passed 2,259 Go tests across 22 packages, four skill
  contract tests, the Roundfix skill check, and the binary build.
- `rtk git diff --check` passed.

No follow-up work was discovered inside this Task's slice.
