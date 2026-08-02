---
task: task_03
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: medium
---

# Task 03: Carry probe evidence into the divergence

## Overview

A capability carries the probe that defines it and an outcome carries the
evidence that produced its verdict, but the divergence projection keeps only a
code, an identifier, and a message. Every renderer downstream therefore has to
re-derive what the evaluation already knew, and none does — which is why a
blocking divergence reads as prose. This Task carries the evaluated probe and
its evidence through the projection. Nothing renders them differently yet.

## Requirements

1. MUST carry the evaluated probe and the evidence that produced the verdict on
   every divergence, so one value feeds text, machine output, and the prompt.
2. MUST carry the probe as evaluated, sourced from the same definition the
   evaluation read, so the diagnostic and the check cannot disagree.
3. MUST keep both additions optional in machine output, so an existing consumer
   reading only today's fields is unaffected.
4. MUST leave every existing field's name, type, and meaning unchanged.
5. MUST NOT change any rendering, grouping, message, or verdict in this Task.

## Subtasks

- [ ] Extend the divergence with the evaluated probe and its evidence.
- [ ] Populate both where divergences are produced.
- [ ] Confirm machine output stays additive.
- [ ] Confirm no rendered output changed.

## Acceptance Criteria

- [ ] Every unsatisfied divergence carries the probe the evaluation read.
- [ ] Every unsatisfied divergence carries the evidence behind its verdict.
- [ ] The probe on the divergence is byte-equal to the definition the
      evaluation used.
- [ ] Machine output gains only optional fields; every prior field keeps its
      name, type, and value.
- [ ] No rendered text changed, proven by the characterization corpus.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run '^TestDivergenceCarriesProbeEvidence$' -count=1 -v | grep -q -- "--- PASS: TestDivergenceCarriesProbeEvidence"`
  — expected: exit 0; probe and evidence are present and match the evaluation.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0; no rendered outcome moved.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go vet ./internal/baseline` — expected: exit 0.

## References

- `_prd.md` → User Story 2; Core Features 5.
- `_techspec.md` → System Architecture; Build Order 3.

## Result

Implemented the capability-to-divergence projection without changing its
rendering. `ProfileDivergence` now carries optional `probe` and `evidence`
fields, populated directly from the capability definition and outcome used by
the evaluation. Decision and Verification divergences omit both fields and
retain their existing JSON shape.

Focused checks:

- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -run 'Test(RequiredDivergencePreventsReadyPlan|DivergenceCarriesProbeEvidence|ProfileAlignmentAdvisoryDivergenceNeverBlocksOrInfersPolicy|PostgreSQLEvidenceSeparatesImplementationAndContract)$' -count=1` — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -run 'TestBaselinePlanCharacterization/(unsatisfied-blocking-capabilities|advisory-only-divergences)$' -count=1` — passed after deliberate corpus regeneration.
- The first focused test attempt with the default Go cache was blocked by the
  sandbox at `/Users/marcio/Library/Caches/go-build`; rerunning with the
  writable cache above passed.

Acceptance evidence:

1. `TestDivergenceCarriesProbeEvidence` iterates every unsatisfied capability
   outcome and requires a matching divergence with a probe.
2. The same test compares each divergence's evidence with the exact
   `CapabilityOutcome.Evidence` that produced its verdict.
3. The same test marshals the evaluated capability definition and divergence
   probe and requires byte-equal JSON.
4. Both new fields use `omitempty`; the legacy-shape assertion proves every
   prior field name, type, value, and ordering remains unchanged when the new
   fields are absent. Characterization diffs add only `probe` and `evidence`.
5. The focused characterization cases pass, and the regenerated golden diff
   changes no code, ID, requirement, blocking flag, message, next action,
   grouping, or verdict.
6. The post-edit worktree status contains only `internal/baseline/` paths and
   this Task file.

The Daemon-owned Verification commands were not run.
