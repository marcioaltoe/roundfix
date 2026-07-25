---
task: task_08
spec: 0047-context-driven-guidance-composition
status: completed
type: docs
complexity: medium
---

# Task 08: Document the composed Baseline workflow

## Overview

Align the user guide and thin setup skill with the shipped hierarchy,
redistribution, residual, document lifecycle, and Profile adaptation behavior.
Humans and automation receive runnable public CLI recipes without a second
execution engine.

## Requirements

1. MUST document greenfield composition, update redistribution, exact rule
   preservation, and residual-carrier behavior.
2. MUST document the Instruction Hierarchy and narrower-guide precedence.
3. MUST document ADR lifecycle and the complete Findings template.
4. MUST document guided Profile adaptation, `--profile-file`, universal
   remediation, Plan review, apply, and re-audit.
5. MUST keep `setup-context-driven` interpretation-only with no executable
   setup engine or behavioral fallback.
6. MUST update Roundfix documentation contract tests with copy-pasteable
   commands and stable exit behavior.

## Subtasks

- [x] Update the Context-Driven user guide.
- [x] Update the thin setup skill recipes and interpretation.
- [x] Add Profile adaptation and redistribution examples.
- [x] Extend documentation and skill-sync contract tests.

## Acceptance Criteria

- [x] Documentation describes every new human and automation state.
- [x] Every command example matches current help and implemented flags.
- [x] The skill contains no independent classification, rendering, mutation,
  or Profile engine.
- [x] The ADR and Findings templates are copyable from generated guidance.
- [x] CLI behavior and skill guidance pass the hard sync guard.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `internal/cli/baseline_documentation_contract_test.go`

## Verification

- `rtk go test -count=1 ./internal/cli ./skills -run 'TestBaselineDocumentationContract|TestThinSetupSkill|TestGuidanceCompositionDocumentation'` — expected: help, examples, thin-skill, and behavior-sync contracts pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline --help` — expected: public help matches the documented command family.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 3–7; Core Feature 17; User Experience; Success Metrics.
- `_techspec.md` → Integration Points; Build Order 7.
- ADR-0066 → CLI authority and thin skill.
- ADR-0075 → documented adaptation behavior.

## Result

The user guide now documents the active Instruction Hierarchy, strengthen-only
narrower-guide precedence, greenfield composition, exact-byte update
redistribution, semantic owners, residual creation and removal, and
warning-only nested carriers. It includes the complete generated ADR lifecycle
overlay and Findings Operational Contract, plus the human and automation
Profile adaptation flow from alignment through universal remediation, Plan
review, apply, and fresh re-audit.

The canonical setup skill remains one interpretation-only `SKILL.md`. It
explains the same states and public CLI recipes without adding classification,
rendering, mutation, or Profile execution. `make skills-sync` refreshed the
distributed copy.

Verification:

- `rtk proxy env GOCACHE=/tmp/roundfix-task08-go-cache go test -count=1
  ./internal/cli ./skills -run
  'TestBaselineDocumentationContract|TestThinSetupSkill|TestGuidanceCompositionDocumentation|TestBaselineExamplesParse|TestBaselineDecisionExamples'`
  — passed both packages.
- `rtk proxy env GOCACHE=/tmp/roundfix-task08-go-cache make
  skills-sync-check` — passed all 4 skill contract checks.
- `rtk proxy env GOCACHE=/tmp/roundfix-task08-go-cache go run
  -buildvcs=false ./cmd/roundfix baseline --help` — exited `0`; public help
  matches the documented `--profile-file`, apply, Profile, restoration, and
  asset command family.
- `rtk proxy env GOCACHE=/tmp/roundfix-task08-go-cache
  GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign
  GIT_CONFIG_VALUE_0=false make verify` — passed 2,189 Go tests across 22
  packages, 4 skill checks, the Roundfix skill check, and the CLI build.
- `rtk git diff --check` — passed.

Acceptance evidence:

- `TestGuidanceCompositionDocumentation` proves the guide describes every
  hierarchy, redistribution, residual, lifecycle, adaptation, remediation,
  apply, and re-audit state required by this Task.
- `TestBaselineExamplesParse`, `TestBaselineDecisionExamples`, and the
  extended help assertions prove all published Baseline commands parse and
  retain the shipped exit-code contract.
- `TestThinSetupSkill` and `TestNoPythonBaselineRuntime` prove the setup skill
  ships only guidance, names the CLI as runtime authority, and contains no
  independent executable engine.
- `TestGuidanceCompositionDocumentation` compares the published ADR and
  Findings templates exactly with the generated formatter golden.
- `TestBaselineSkillContract`, `skills-sync-check`, and the full gate prove
  canonical, distributed, and embedded guidance remain synchronized.

The first sandboxed full-gate attempt reached 2,188 passing tests, then
`TestFormatterComposition` inherited the host's `commit.gpgsign=true` and
could not access `/Users/marcio/.gnupg`. The exact subtest reproduced before
Baseline assertions and passed with only the process-scoped signing override;
the isolated full gate then passed. Hardening that unrelated fixture against
host signing configuration belongs to a follow-up Task and is not part of this
documentation diff. The Daemon owns the authoritative Verification rerun.
