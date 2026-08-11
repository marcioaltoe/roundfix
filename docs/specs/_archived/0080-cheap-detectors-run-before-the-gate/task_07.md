---
task: task_07
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: test
complexity: low
---

# Task 07: Re-record the expectation task_06 invalidates

## Overview

The consequent fix, landing as its own commit after its cause. task_06's new
clauses grow the Source Baseline, and the maintained compatibility fixture
declares its expected shape by hand so a legitimate corpus change moves one
declared value instead of hunting literals.

This node exists because folding the correction into task_06 would fail the
tooling-authority gate: `docs/agents/specific-repository.md` requires a
consequent fix — one that only became necessary because the authorized change
made something stale — to be its own commit landing after the Task commit.
Spec 0079 paid a full gate round to learn that, and this graph spends a node
instead.

## Requirements

1. MUST confirm the corpus change is legitimate before moving anything: the
   added entries trace to the clauses task_06 authored under the recorded
   authorization, and the digest chain is converged.
2. MUST run the sanctioned re-recording workflow rather than hand-editing any
   generated artifact; only the declared expectation is authored by hand.
3. MUST move the maintained Source Baseline expectation to the value the
   regenerated fixture carries, keeping the identity-agrees-with-entries
   invariant and the accounting expectation intact as separate assertions.
4. MUST NOT weaken the assertion: no tolerance, no computing the expected
   value from the fixture under test, no skip. The expectation stays a
   declared constant a reviewer can read.
5. MUST leave every other package untouched; this is a one-expectation
   correction, not a refactor.

## Subtasks

- [ ] Confirm provenance and chain convergence.
- [ ] Re-record through the sanctioned workflow; move the expectation.
- [ ] Prove the package and the repository gate green.

## Acceptance Criteria

- [ ] The maintained-fixture test passes with the moved expectation.
- [ ] The invariant and accounting assertions remain present and unweakened.
- [ ] The whole `internal/baseline` package passes.
- [ ] No file outside `internal/baseline` and this task file changed.

## Context

- interface: internal/baseline/preservation_test.go

## Verification

- `output="$(go test -count=1 ./internal/baseline -run '^TestReadoptionCompatibilityMaintainedFixture$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the named fixture test is selected and passes.
- `go test -count=1 ./internal/baseline/...`
  — expected: exit 0; the package the correction lives in is green.
- `grep -q 'sourceBaseline.Identity.EntryCount != len(sourceBaseline.Entries)' internal/baseline/preservation_test.go && grep -q 'maintainedSourceBaselineAccounting' internal/baseline/preservation_test.go && ! grep -q 'maintainedSourceBaselineEntries    = 132' internal/baseline/preservation_test.go` — expected: exits 0. The first two halves prove the invariant and the named
constant survive as separate assertions; the third proves the constant actually
moved off `132`, the entry count Task 06's clauses invalidated. The Agent proved
on 2026-08-11 that `maintainedSourceBaselineAccounting` legitimately stays at
`51` — the sanctioned regeneration reports 134 entries and 51 accounting rows —
and refused to make an assertion true that the regenerated artifacts say is
false. This command named the wrong constant; the entry count is the one that
moves. Without it the command
asserts only that two strings this file already contains are still there, and
the pre-work probe refused exactly that on 2026-08-11.
  — expected: exit 0; the invariant and accounting assertions survived rather
  than being relaxed.
  — expected: exit 0; the correction stayed inside its package.

## References

- `_techspec.md` → Build Order 7.
- ADR-0081.
- `docs/findings/_archived/` → the 2026-08-06 rollup covering gate and
  verification evidence, whose members record what folding this fix costs.

## Result

### Implementation

- Confirmed the legitimate corpus change at task_06 commit `9d36349a`: its
  expressly authorized module edits added
  `clause.core.verification-two-tiers` and
  `clause.spec.verification-two-tiers`, and the sanctioned fallout added the
  matching marker-delimited Source Baseline entries and manifest rows.
- Ran the sanctioned `baseline-digests` workflow before authoring the
  expectation. It exited 0 with `changed:false`, proving the regenerated
  digest chain already matched its canonical sources. The Source Baseline
  identity and manifest each declare 134 entries, while `accounting.json`
  retains 51 rows.
- Moved only `maintainedSourceBaselineEntries` from `132` to `134` in the
  maintained-fixture test. The separately named
  `maintainedSourceBaselineAccounting = 51`, the identity-versus-entry-count
  comparison, and the three independent count comparisons remain unchanged.

### Focused checks

- `GOCACHE=/private/tmp/roundfix-task07-gocache rtk make baseline-digests`
  exited 0 and reported `changed:false`; no generated artifact changed.
- `GOCACHE=/private/tmp/roundfix-task07-gocache rtk go test -count=1
  ./internal/baseline -run
  '^(TestReadoptionCompatibilityMaintainedFixture|TestSourceBaselineGuidanceComposition)$'`
  exited 0 with 2 passing tests.
- `GOCACHE=/private/tmp/roundfix-task07-gocache rtk go test -shuffle=on
  ./internal/baseline` exited 0 with 593 passing tests.
- `GOCACHE=/private/tmp/roundfix-task07-gocache rtk make
  verify-incremental` was also attempted and exited 2. Its
  `internal/baseline` package passed, while existing QA-report expectation
  failures remained in `internal/cli` and a mechanical-stage Git-head setup
  failure remained in `internal/daemon`. Those packages are outside this
  one-expectation Task and were not edited.
- `rtk git -c core.fsmonitor=false diff --check` exited 0. The changed-path
  audit contains only `internal/baseline/preservation_test.go` and this Task
  file; the Task-file frontmatter transition was the pre-existing
  Daemon-owned change.

### Acceptance-criterion evidence

1. The focused maintained-fixture selection passed after the declared entry
   expectation moved to 134.
2. The diff changes only the entry-count constant; the identity-agreement
   comparison, the accounting constant at 51, and its separate comparison
   remain present and unweakened.
3. A fresh shuffled run of the whole `internal/baseline` package passed all
   593 tests.
4. The changed-path audit contains no path outside `internal/baseline` and
   this Task file.

### Handoff boundary

- The Daemon-owned commands under `## Verification` were not run. Task status
  remains Daemon-owned; no commit, push, or pull request was created.
