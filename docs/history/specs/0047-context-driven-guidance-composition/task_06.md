---
task: task_06
spec: 0047-context-driven-guidance-composition
status: completed
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

- [x] Move human alignment before classification.
- [x] Add the reviewed adaptation interaction loop.
- [x] Add strict `--profile-file` planning input.
- [x] Emit actionable universal-capability remediation.
- [x] Add human, automation, output, and no-write tests.

## Acceptance Criteria

- [x] The Oraculum-shaped fixture receives a guided adaptation instead of an
  aggregate dead-end result.
- [x] Required built-in capabilities never become waivers.
- [x] Human and automation draft inputs produce identical Profile, postimages,
  and Plan Digest.
- [x] Missing universal skills stop before classification with exact next
  actions.
- [x] Decline, invalid draft, stale draft, and output failure write no
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

## Result

Implemented Profile divergence as a complete pre-classification alignment
state. The human root workflow now renders every blocking and advisory
divergence, offers Profile change, a reviewed repository-owned adaptation, or
decline for profile-specific blockers, re-audits the normalized draft, and
starts instruction classification only after alignment is ready. Every
proposed module and capability removal is listed for explicit confirmation,
and the Profile ID is validated by the strict draft engine.

Automation now accepts `baseline plan --profile-file <draft.json>` as the
mutually exclusive counterpart to `--profile`. The strict custom Profile
document is source-resolved, normalized through the same in-memory draft path
as the human workflow, and remains read-only until its portable Plan Digest is
approved. Missing universal Context7 or Exa capabilities remain required and
report the exact `baseline skills restore` preview and `--confirm-plan`
operation before classification.

Verification:

- `rtk go test -count=1 ./internal/cli ./internal/baseline -run 'TestBaselineHumanProfileAdaptation|TestBaselinePlanProfileFile|TestProfileDivergenceResolution|TestUniversalCapabilityRemediation'` — passed, 10 tests.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline plan --help` — passed; help documents mutually exclusive `--profile` and `--profile-file` inputs.
- `rtk go test -count=1 ./internal/cli ./internal/baseline` — passed, 1,015 tests.
- `rtk make verify` — passed: 2,180 Go tests across 22 packages, 4 skill-contract tests, Roundfix skill checks, and binary build.
- `rtk git diff --check` — passed.

Acceptance evidence:

- `TestBaselineHumanProfileAdaptation` proves the backend-only TypeScript
  fixture receives individual blocking/advisory output, all three resolution
  choices, an explicit removal review, a strict repository-owned Profile ID,
  a ready re-audit, and a consolidated Change Plan.
- `TestProfileDivergenceResolution` proves the adapted Profile retains
  universal Context7 and Exa requirements as satisfied required capabilities
  rather than waivers.
- `TestBaselineHumanProfileAdaptation` and `TestBaselinePlanProfileFile` prove
  human and automation drafts normalize to identical Profile, postimages, and
  Plan Digest.
- `TestUniversalCapabilityRemediation` proves a missing universal skill is
  blocking and names the exact restoration preview and confirmation command.
- `TestBaselineHumanProfileAdaptation` and `TestBaselinePlanProfileFile` prove
  decline, invalid input, stale catalog binding, and human or automation output
  failure leave repository bytes unchanged.

One broad package run encountered the known transient concurrent Git
`maintenance.lock` race in `TestConsolidatedReview`. The exact test passed in
isolation, and a fresh full package run plus `make verify` both passed without
an out-of-scope workaround.

Follow-ups: none.
