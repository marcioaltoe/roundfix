---
task: task_08
spec: 0047-context-driven-guidance-composition
status: pending
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

- [ ] Update the Context-Driven user guide.
- [ ] Update the thin setup skill recipes and interpretation.
- [ ] Add Profile adaptation and redistribution examples.
- [ ] Extend documentation and skill-sync contract tests.

## Acceptance Criteria

- [ ] Documentation describes every new human and automation state.
- [ ] Every command example matches current help and implemented flags.
- [ ] The skill contains no independent classification, rendering, mutation,
  or Profile engine.
- [ ] The ADR and Findings templates are copyable from generated guidance.
- [ ] CLI behavior and skill guidance pass the hard sync guard.

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
