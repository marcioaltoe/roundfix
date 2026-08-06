---
task: task_04
spec: 0065-loop-order-and-verification-honesty
status: completed
type: backend
complexity: low
---

# Task 04: Check that the order restatements agree

## Overview

task_01 states the loop order once. This slice makes a second divergence
mechanically impossible, so the next edit to any one statement cannot silently
recreate the contradiction.

It depends on task_01 for a reason that is not cosmetic: this rule fails while
the statements disagree, and `spec check` runs inside `make verify`. Landing it
first would leave the repository red between two Tasks, and the Daemon runs the
configured Verification command as a precondition — so the Task meant to repair
that state would be settled without ever starting. Spec 0075 lost a Run to that
exact shape on 2026-08-05.

## Requirements

1. MUST add `SC-LOOP-ORDER-DIVERGENT`, failing when the loop order differs
   between the places that state it.
2. MUST cover all three sources: the shipped clause, `docs/agents/autonomous-work.md`,
   and the Baseline module asset. A check reading only two recreates the defect
   in a smaller form.
3. MUST report which sources disagree and how, so the fix does not require
   diffing three files by hand.
4. MUST pass on the corrected statements task_01 produced, proving the rule and
   the sources agree at the moment it lands.
5. MUST keep `TestCheckCorpusBudget` passing.
6. MUST NOT restate or edit the order itself; task_01 owns the content.

## Subtasks

- [x] Add the rule reading all three sources.
- [x] Fixture a divergence in each source and assert detection.
- [x] Assert the corrected repository passes.

## Acceptance Criteria

- [x] A divergence in the shipped clause is detected and named.
- [x] A divergence in the repository guide is detected and named.
- [x] A divergence in the Baseline module asset is detected and named.
- [x] The repository as task_01 left it passes the rule.
- [x] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/agents/autonomous-work.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'LoopOrder|Divergent' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the divergence tests ran and passed.
- `go run -buildvcs=false ./cmd/roundfix spec check > /dev/null` — expected:
  exit 0; the corrected repository passes its own new rule.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 1; Success Metric 1.
- `_techspec.md` → Build Order 1; Risks & Considerations.
- ADR-0093.

## Result

Implementation evidence:

- Added `SC-LOOP-ORDER-DIVERGENT` to the Spec Consistency Check. It reads the
  shipped generated clause, the repository guide, and the Baseline module's
  `clause.autonomous.loop-01-qa-once` guidance. The module is parsed as JSON by
  clause identifier; each order is extracted from its explicit declaration,
  keeping the check inside ADR-0093's declared-source boundary.
- A mismatch produces an `error` whose summary names all three sources and
  prints each observed order. Its locations point to all three declarations,
  so the diagnostic shows both the disagreeing source and the difference.
- Added independent fixtures that reorder `archive` and `merge` in each source,
  plus a repository-backed case that asserts task_01's corrected declarations
  produce no `SC-LOOP-ORDER-DIVERGENT` finding.
- Registered the new code in the active and archived characterization corpus
  with zero expected findings.

Focused-check evidence:

- Red signal: `rtk env GOCACHE=<repository>/.gocache rtk go test
  ./internal/speccheck -count=1 -run '^TestCheckLoopOrder'` failed to compile
  before implementation because `speccheck.CodeLoopOrderDivergent` did not
  exist.
- `rtk env GOCACHE=<repository>/.gocache rtk go test ./internal/speccheck
  -count=1 -run '^TestCheckLoopOrder'` — exit 0; five loop-order tests passed,
  including one divergent case for each source and the corrected repository
  case.
- `rtk env GOCACHE=<repository>/.gocache rtk go test ./internal/speccheck
  -count=1 -run '^(TestCheckLoopOrder|TestCheckCorpusGolden)$'` — exit 0; the
  loop-order and corpus characterization selectors passed.
- `rtk env GOCACHE=<repository>/.gocache rtk go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` — exit 0; the dedicated
  corpus sweep remained below its wall-clock budget.
- `rtk env GOCACHE=<repository>/.gocache rtk go test ./internal/speccheck
  -count=1 -skip '^TestCheckCorpusBudget$'` — exit 0 with 58 tests passed.
- `rtk git diff --check` — exit 0.

Acceptance-criterion evidence:

- Shipped clause: `TestCheckLoopOrderDivergent/shipped_clause` passed and
  required the finding summary to name `shipped clause`, print its reordered
  actions, and locate the shipped clause path.
- Repository guide: `TestCheckLoopOrderDivergent/repository_guide` passed with
  the equivalent name, observed-order, and location assertions.
- Baseline module asset:
  `TestCheckLoopOrderDivergent/Baseline_module_asset` passed with the equivalent
  assertions against the clause read from the JSON module.
- Corrected repository: `TestCheckLoopOrderRepositoryAgrees` passed with no
  `SC-LOOP-ORDER-DIVERGENT` finding.
- Corpus budget: the dedicated `TestCheckCorpusBudget` selector exited 0.

Daemon-owned Verification:

- The commands under `## Verification` were not run in this Agent turn. The
  Daemon owns those commands and the terminal Task verdict.
