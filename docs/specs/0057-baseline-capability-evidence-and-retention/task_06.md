---
task: task_06
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: high
---

# Task 06: Offer a read-only capability re-check

## Overview

After remediating a blocking divergence there is no way to ask whether it is
fixed without resolving decisions and driving a full plan, so the
remediate-and-re-check loop does not exist. This Task adds a read-only re-check
that needs no decisions, writes nothing, and produces the same capability
outcomes a full plan would.

## Requirements

1. MUST provide a read-only capability re-check that requires no decisions to
   be supplied and resolves none.
2. MUST write nothing: no file, no journal mutation, no configuration change.
3. MUST produce capability outcomes identical to those a full plan produces for
   the same repository, by sharing the evaluation path rather than
   reimplementing it.
4. MUST render the same probe evidence and requirement grouping the full plan
   renders.
5. MUST report clearly when the repository has no resolvable Profile, rather
   than failing obscurely.
6. MUST leave the full plan's behavior and output unchanged.

## Subtasks

- [ ] Add the read-only re-check entry point over the shared evaluation.
- [ ] Require no decisions and resolve none.
- [ ] Render probes and grouping identically to the full plan.
- [ ] Confirm nothing is written.

## Acceptance Criteria

- [ ] The re-check completes with zero decisions supplied.
- [ ] Its capability outcomes match a full plan's for the same repository,
      asserted field by field.
- [ ] It renders the same probe evidence and grouping as the full plan.
- [ ] Running it leaves the repository byte-identical and writes no journal
      entry.
- [ ] A repository with no resolvable Profile produces a named error rather
      than a panic or an empty result.
- [ ] The full plan's output is unchanged, proven by the characterization
      corpus.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline ./internal/cli -run '^TestCapabilityRecheck$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityRecheck"`
  — expected: exit 0; the re-check needs no decisions and writes nothing.
- `go test ./internal/baseline -run '^TestCapabilityRecheckMatchesFullPlan$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityRecheckMatchesFullPlan"`
  — expected: exit 0; outcomes match field by field.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 4; Core Features 6; Success Metrics (obtainable with
  zero decisions and matching full-plan outcomes).
- `_techspec.md` → API Contracts; Build Order 5.

## Result

### Implementation

- Added `roundfix baseline capabilities check` with text and
  `roundfix/baseline-capability-recheck/v1` JSON output. The command accepts no
  decision input, never constructs a Change Plan, and exits `3` only when the
  evaluated Repository Capabilities contain a blocking divergence.
- Added a baseline-owned `RecheckCapabilities` entry point. It accepts only a
  repository and optional Baseline Profile ID; without an explicit ID it reads
  the current Setup Manifest. The result projects only Profile identity,
  capability outcomes, capability divergences, and PostgreSQL evidence — no
  resolved decisions or write-bearing plan fields.
- Extracted capability evaluation and capability-divergence construction from
  `ResolveProfileAlignment` into one shared path. Full planning and re-checking
  now receive the same `CapabilityOutcome` values and evidence-carrying,
  requirement-grouped divergences by construction.
- Connected the interactive full-plan alignment review and the re-check text
  command to Task 04's shared `RenderProfileDivergences` projection. This
  completes the prior presentation seam without changing readiness, blocking
  flags, plan documents, or the characterization corpus.
- Added the named `ErrNoResolvableProfile` error. Missing, invalid, or
  unavailable Profile selection produces a structured `profile` failure and a
  next action instead of a panic or empty result.

### Focused checks

- The initial focused `CapabilityRecheck` test run failed before implementation
  with undefined re-check API and result symbols, establishing the missing
  behavior. The same focused run passed after implementation in both
  `internal/baseline` and `internal/cli`.
- `go test ./internal/baseline -run
  'Test(ProfileAlignment|RequiredDivergence|DivergenceCarries|PostgreSQL|UniversalCapability|CapabilityRecheck|DivergenceRenders|DivergenceGroups)'
  -count=1` with the task-scoped Go cache passed.
- `go test ./internal/cli -run
  'Test(CapabilityRecheck|BaselineHumanProfileAdaptation|BaselineDocumentationContract|BaselineProfileHelpContract)$'
  -count=1` with the task-scoped Go cache passed.
- The two focused characterization cases,
  `unsatisfied-blocking-capabilities` and `advisory-only-divergences`, passed
  against their existing goldens without regeneration.
- `git diff --check` passed.

### Acceptance evidence

1. `TestCapabilityRecheck` invokes the public command without any decision
   arguments and asserts the JSON document has no top-level decision field.
2. `TestCapabilityRecheckMatchesFullPlan` compares every
   `CapabilityOutcome` field with `reflect.DeepEqual`, then compares the
   capability divergence projection field by field against full-plan Profile
   alignment for the same repository.
3. The equivalence test compares rendered divergence output, while the CLI test
   requires the command's complete divergence body to equal the shared
   renderer. The interactive Profile-adaptation flow now exercises the same
   blocking, advisory, informational, probe, and advisory-language rendering.
4. The CLI test hashes every repository path and byte, including a pre-existing
   journal sentinel, before and after both JSON and text re-checks and requires
   byte-identical results.
5. The baseline test requires `errors.Is(err, ErrNoResolvableProfile)` without
   a Profile source. The CLI test requires the named error text and structured
   `profile` category without a panic or partial result.
6. The focused characterization cases passed byte-for-byte against the Task 04
   goldens. The shared evaluator refactor changed no plan state, readiness,
   blocking flag, or machine plan output.
7. The final changed-path postflight is limited to `internal/baseline/`,
   `internal/cli/`, and this Task file.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
