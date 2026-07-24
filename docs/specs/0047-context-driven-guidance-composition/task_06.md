---
task: task_06
spec: 0047-context-driven-guidance-composition
status: pending
type: backend
complexity: high
---

# Task 06: Guide public Profile divergence resolution

## Overview

Expose Profile adaptation as a complete human alignment state and as explicit
automation input. The root workflow re-audits every reviewed draft before
instruction classification and never presents a capability waiver.

## Requirements

1. MUST show every blocking and advisory divergence before classification.
2. MUST offer Profile change, repository-owned adaptation, or decline for
   profile-specific blockers.
3. MUST require explicit review of every removed module and capability plus a
   valid repository-owned Profile ID.
4. MUST repeat alignment after adaptation and proceed only when all
   non-removable requirements are satisfied.
5. MUST add mutually exclusive `baseline plan --profile` and
   `--profile-file` inputs with equivalent normalized results.
6. MUST report exact supported remediation operations for missing universal
   capabilities.
7. MUST preserve stdout, stderr, exit-category, final confirmation, and
   no-write-before-approval contracts.

## Subtasks

- [ ] Move human alignment before classification.
- [ ] Add the reviewed adaptation interaction loop.
- [ ] Add strict `--profile-file` planning input.
- [ ] Emit actionable universal-capability remediation.
- [ ] Add human, automation, output, and no-write tests.

## Acceptance Criteria

- [ ] The Oraculum-shaped fixture receives a guided adaptation instead of an
  aggregate dead-end result.
- [ ] Required built-in capabilities never become waivers.
- [ ] Human and automation draft inputs produce identical Profile, postimages,
  and Plan Digest.
- [ ] Missing universal skills stop before classification with exact next
  actions.
- [ ] Decline, invalid draft, stale draft, and output failure write no
  repository bytes.

## Context

- instruction: `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- instruction: `docs/adr/0075-profile-divergence-uses-confirmed-repository-owned-adaptation.md`
- interface: `internal/cli/baseline_human.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/baseline/profile_alignment.go`

## Verification

- `rtk go test -count=1 ./internal/cli ./internal/baseline -run 'TestBaselineHumanProfileAdaptation|TestBaselinePlanProfileFile|TestProfileDivergenceResolution|TestUniversalCapabilityRemediation'` — expected: human loop, automation parity, output, exits, and refusal cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline plan --help` — expected: help documents mutually exclusive Profile ID and Profile draft inputs.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 6; User Story 7; Core Features 18–20; User Experience.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 5.
- ADR-0068 → one confirmation-gated workflow.
- ADR-0075 → no required-capability waiver.
