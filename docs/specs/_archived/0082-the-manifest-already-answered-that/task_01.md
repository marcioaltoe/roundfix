---
task: task_01
spec: 0082-the-manifest-already-answered-that
status: completed
type: test
complexity: medium
---

# Task 01: Pin the adoption path with a characterization corpus

## Overview

Every later task in this Spec changes code that first adoption also runs. This
task captures the current, correct behavior of the adoption and non-interactive
planning paths as executable characterization cases, so that a regression in
first adoption becomes a failing test rather than a surprise in someone else's
repository. It delivers no product behavior and is verifiable on its own: the
corpus passes against today's code, unchanged.

## Requirements

1. MUST capture the current behavior of first adoption on a repository with no
   Setup Manifest, covering both instruction preservation modes, as cases that
   assert observable outputs rather than internal call order.
2. MUST capture the current behavior of non-interactive planning when decisions
   are absent, when decisions are supplied but the preservation mode is not, and
   when a complete strict Decision Document is supplied.
3. MUST capture the current retention-accounting outcome for a manifest whose
   managed artifacts diverge from the target, including the blocking outcome
   when a clause is unaccounted.
4. MUST assert against stable identities — exit category, result schema,
   decision ids named, clause dispositions — never against a whole rendered
   text blob whose formatting is free to change.
5. MUST pass against the current code with no production change in this task.
6. SHOULD extend the existing corpus seams rather than introducing a new test
   harness.

## Subtasks

- [x] Inventory which existing corpus file owns each captured behavior.
- [x] Capture the first-adoption cases for both preservation modes.
- [x] Capture the non-interactive planning refusal and success cases.
- [x] Capture the retention-accounting cases including the unaccounted-clause block.
- [x] Confirm the whole corpus passes with no production change.

## Acceptance Criteria

- [x] The corpus names every captured case explicitly; a reader can list what is
      pinned without reading assertions.
- [x] Each captured case asserts a stable identity, not rendered prose.
- [x] The corpus fails if the exit category of any captured path changes.
- [x] The corpus fails if a decision id disappears from the missing-decision
      result for the unanswered planning case.
- [x] No file outside test sources and test data changed in this task.

## Context

- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/plan_characterization_test.go`
- interface: `internal/cli/baseline_plan_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'Characterization|Corpus' -v > /tmp/task_01-1.log 2>&1 && grep -q '^--- PASS: ' /tmp/task_01-1.log` — expected: exits 0, proving at least one characterization case ran and passed rather than being silently selected out.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.
- `git diff --name-only HEAD -- internal/ | grep -v -E '(_test\.go|/testdata/)' | grep . ; test $? -eq 1` — expected: exits 0, proving this task changed no production source.

## References

- `_techspec.md` → Build Order 1; Testing Approach.
- `_prd.md` → Success Metrics.
- ADR-0058, ADR-0068, ADR-0071.

## Result

### Implementation

- `internal/cli/baseline_plan_test.go` owns first-adoption and non-interactive
  planning characterization. It names greenfield success, preservation success
  from a strict Decision Document, absent decisions, and absent preservation
  mode as separate cases.
- `internal/baseline/plan_characterization_test.go` owns managed-artifact drift
  retention characterization. It names retained-clause success and the
  unaccounted-clause block separately.
- `internal/baseline/compatibility_corpus_test.go` remains the frozen historical
  parity corpus; this Task did not change its recorded artifacts or regeneration
  contract.
- No production source or test data changed. The only non-test edit is this
  Daemon-required Task evidence section; the Daemon-owned `status` remains
  `in_progress`.

### Focused checks

- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-01-baseline-gocache go test ./internal/baseline -run '^(TestBaselineRetentionAccountingCharacterizationCorpus|TestReadyPlanNeverCarriesEmptyLedger|TestSameIdentityDriftRequiresRetention)$' -count=1`
  exited 0: `ok roundfix/internal/baseline`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-0082-task-01-cli-gocache go test ./internal/cli -run '^(TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus|TestBaselinePlanCommandEmitsPortableJSONAndNormalizesDecisionFiles)$' -count=1`
  exited 0: `ok roundfix/internal/cli`.
- `git -c core.fsmonitor=false diff --check` exited 0.
- The first focused attempts did not reach compilation because the sandbox
  denied writes to the host Go build cache. The commands above reran with
  isolated writable caches and provide the applicable evidence.
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance evidence

- The six table-case names enumerate every captured path without requiring a
  reader to inspect assertions.
- Success paths assert Plan schema and Profile identity. Refusal paths assert
  exit code 3 plus Result schema, operation, state, and category. Preservation
  and drift paths assert exact clause dispositions (`repository-rules`,
  `retained`, and `unaccounted`) rather than rendered prose.
- The unanswered case resolves the current `go-cli-tui` Profile and requires
  every selected decision ID to remain named in the missing-decision Result.
- `git diff --name-only HEAD` lists only this Task file and the two `_test.go`
  corpus files. No production source is changed.
