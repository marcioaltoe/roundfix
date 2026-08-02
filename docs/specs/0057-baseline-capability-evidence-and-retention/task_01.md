---
task: task_01
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: test
complexity: medium
---

# Task 01: Record how planning behaves today

## Overview

This Spec's largest slice turns a plan that completes today into one that can
stop. The only way to bound that is to record what completes today before
anything moves. This Task captures plan outcomes and diagnostics for real
repository shapes and changes no behavior; it is the gate every later slice is
measured against.

## Requirements

1. MUST record, for a corpus of repository shapes, the complete plan outcome:
   state, every diagnostic, every divergence, and every warning.
2. MUST cover at minimum a clean adoption, an idempotent re-plan after a
   verified apply, a repository with unsatisfied blocking capabilities, one
   with advisory-only divergences, and one whose Profile and catalog digests
   changed under an unchanged Baseline identifier.
3. MUST fail with a readable diff naming the affected shape and the changed
   field when a later change alters any recorded outcome.
4. MUST be regenerable through an explicit flag so an intended change is
   re-recorded deliberately, never silently.
5. MUST NOT change any production behavior, diagnostic, or exported API.

## Subtasks

- [ ] Assemble the repository-shape corpus.
- [ ] Record each shape's plan outcome deterministically.
- [ ] Add the comparison with a readable failure diff.
- [ ] Add the explicit regeneration flag.

## Acceptance Criteria

- [ ] The corpus contains a case for each shape named in Requirement 2.
- [ ] A test plans each shape and compares against its golden, passing on the
      unmodified tree.
- [ ] Deliberately altering one recorded outcome makes the test fail and name
      the affected shape.
- [ ] Two consecutive runs produce the same result and rewrite no golden.
- [ ] Goldens are re-recordable only through the explicit flag.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0; the corpus matches the current tree.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0 on a second consecutive run, proving the comparison is stable
  and self-recording is gated.
- `grep -rq "TestBaselinePlanCharacterization" internal/baseline` — expected:
  exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 11; Success Metrics.
- `_techspec.md` → Testing Approach: characterization corpus; Build Order 1.

## Result

### Implementation

- Added a deterministic `TestBaselinePlanCharacterization` corpus at the
  existing plan/alignment seams. Each record contains the complete
  `ProfileAlignment` and `PlanOutcome`, including capability diagnostics,
  divergences, retention, warnings, plan state, and all plan projections.
- Recorded clean adoption, verified-apply idempotent re-plan, blocking
  capability absence, advisory-only divergence, and unchanged Baseline
  identity with changed Profile and catalog digests as five separate goldens.
- Added field-aware comparison. A changed golden reports the repository shape,
  the first changed JSON field, and its old and new values.
- Added the explicit `-update-baseline-plan-characterization` regeneration flag.
  Missing or changed goldens fail without that flag; normal runs only read them.
- Added no production code, production-only test hook, diagnostic change, or
  exported API.

### Focused checks

- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -update-baseline-plan-characterization -count=1`
  with a writable repository-local `GOCACHE`: passed (6 tests) and wrote the
  five named goldens through the explicit flag.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=2 -shuffle=off`
  with the same `GOCACHE`: passed (12 tests across two consecutive corpus
  runs).
- SHA-256 of every golden before and after the two normal runs was identical;
  no golden was rewritten without the update flag.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterizationDiffNamesShapeAndField$' -count=1`
  with the same `GOCACHE`: passed. The negative comparison changes
  `$.outcome.Result.state` and requires the error to name both that field and
  `clean-adoption`.
- A first comparator-only compile attempt was environment-blocked by the
  sandbox-denied host Go cache. The repository-local `GOCACHE` rerun above
  passed; no implementation change was made to bypass the failure.
- `git -c core.fsmonitor=false status --short --untracked-files=all`: only this
  task file, `internal/baseline/plan_characterization_test.go`, and the five
  `internal/baseline/testdata/plan-characterization/*.golden.json` files are
  changed.

### Acceptance criteria evidence

- Corpus coverage: the five golden filenames map one-to-one to every shape in
  Requirement 2. The recorded advisory case contains only the non-blocking
  `capability.firecrawl` divergence; the blocking case contains the blocking
  `capability.context7` divergence.
- Golden comparison: the two-run focused check passed against all five recorded
  full outcomes on the unmodified implementation tree.
- Readable failure: the comparator negative test passed with the shape and
  changed field in its required error text.
- Stable consecutive runs: the two-run check passed and all five pre/post
  golden checksums were byte-identical.
- Explicit regeneration: the update-flag run generated the corpus; two normal
  runs left every checksum unchanged, and the comparison path fails rather than
  writing when the flag is absent.
- Bounded paths: the focused status inspection lists no path outside
  `internal/baseline/` and this task file.

### Not run

- The commands under `## Verification` were not run; the Daemon owns that gate
  and Task settlement.
