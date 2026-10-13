---
task: task_04
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: pending
type: test
complexity: medium
---

# Task 04: The Real-Time Budget rule holds each test file to a ceiling

## Overview

Specs 0171 and 0213 fixed this class one test at a time and it came back
(`docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`). This Task
adds the Real-Time Budget rule beside the Parallel Test Package rule, measures
each test file's ceiling on the tree task_01 to task_03 leave, and writes the
**Real-Time Budget** term and the repository rule. It is verifiable on its own:
the rule passes on the repository and fails, naming file and line, on each
seeded violation.

## Requirements

1. MUST answer Spec 0259's QA failures of 2026-10-12 with a rule that keeps the
   class from growing back, as ADR-0271 decides.
2. MUST add `TestEveryRealTimeBudgetIsWithinItsFileCeiling` and
   `TestRealTimeBudgetRuleRejectsSeededViolations` in
   `internal/testfixture/real_time_budget_test.go`, following `_techspec.md` →
   Interfaces and API Contract 7: the walk, the excluded directories, the
   constant duration expressions including constants declared in test files of
   the same package, the three
   forms, the threshold `realTimeBudgetThreshold` of ten seconds, the
   `<file>:<line>: Real-Time Budget <duration> in <form>` lines after a line
   naming the file, its count and its ceiling, and failure only when a count
   exceeds its ceiling or a file without a ceiling has any budget.
3. MUST seed the trees `unlisted-timeout`, `over-ceiling`, `parameter-budget`
   and `named-constant`, each failing with the seeded file and line, and
   `long-budget`, `sleep-only` and `listed`, each passing, as subtests named
   after the trees.
4. MUST fill `realTimeBudgetCeilings` with one entry per test file that still
   holds a budget after task_01 to task_03, its ceiling equal to the count
   measured on that tree and a non-empty reason saying why each remaining
   budget is safe (asserted against blocking work, data for a fake clock, or a
   budget that bounds no child, goroutine or Git work). No ceiling MAY cover a
   row of `_techspec.md` → Current behavior that task_01 to task_03 fixed.
5. MUST make both tests call `t.Parallel()` first and keep
   `internal/testfixture/parallel_tests_test.go` byte-identical.
6. MUST add **Real-Time Budget** to `CONTEXT.md` through `domain-modeling`,
   after **Parallel Test Package**, linked to ADR-0271 and to no Spec, with an
   `_Avoid_:` line.
7. MUST add to `docs/agents/specific-repository.md`, outside every setup marker
   and after the "Keep the project KISS" bullet, one rule: a test bounds by real
   time only the budget it asserts, every other wait ends on the event its
   assertion reads or the test deadline, and the Real-Time Budget rule in
   `internal/testfixture` holds each test file to a ceiling that only falls.
   The two bullets 0259 adds before the Parallel Test Package bullet MUST NOT
   be touched.

## Subtasks

- [ ] Write the walk, the constant evaluator and the three forms.
- [ ] Write the seeded trees and their expectations.
- [ ] Measure and record each file's ceiling with its reason.
- [ ] Add the glossary term and the repository rule.

## Acceptance Criteria

- [ ] The rule passes on the repository.
- [ ] Each of the four violating seeded trees fails naming its seeded line; the three clean trees pass.
- [ ] Every ceiling equals its file's count on this tree and carries a reason.
- [ ] `CONTEXT.md` defines **Real-Time Budget** and the repository guide states the rule.

## Context

- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `internal/testfixture/parallel_tests_test.go`
- instruction: `internal/testfixture/written_executable_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/testfixture/real_time_budget_test.go`
- interface: `CONTEXT.md`
- interface: `docs/agents/specific-repository.md`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Real-Time Budget**:' || { printf 'CONTEXT.md does not define %s\n' '**Real-Time Budget**' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'Real-Time Budget rule' || { printf 'the repository guide does not state the Real-Time Budget rule\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestEveryRealTimeBudgetIsWithinItsFileCeiling|TestRealTimeBudgetRuleRejectsSeededViolations|TestEveryTestInAParallelTestPackageRunsInParallel)$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryRealTimeBudgetIsWithinItsFileCeiling TestRealTimeBudgetRuleRejectsSeededViolations TestRealTimeBudgetRuleRejectsSeededViolations/unlisted-timeout TestRealTimeBudgetRuleRejectsSeededViolations/over-ceiling TestRealTimeBudgetRuleRejectsSeededViolations/parameter-budget TestRealTimeBudgetRuleRejectsSeededViolations/named-constant TestRealTimeBudgetRuleRejectsSeededViolations/long-budget TestRealTimeBudgetRuleRejectsSeededViolations/sleep-only TestRealTimeBudgetRuleRejectsSeededViolations/listed; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the term, the rule line and the rule tests are absent, so the command fails.
- `test -f internal/testfixture/real_time_budget_test.go || { printf 'the Real-Time Budget rule file is absent\n' >&2; exit 1; }; for path in internal/spec/source_revision_test.go internal/store/prune_reclaimable_test.go internal/agent/agent_test.go; do ! grep -Eq "\"$path\": +[{]" internal/testfixture/real_time_budget_test.go || { printf 'a ceiling covers %s, whose only budget was fixed\n' "$path" >&2; exit 1; }; done; grep -q 'realTimeBudgetThreshold = 10 [*] time.Second' internal/testfixture/real_time_budget_test.go || { printf 'the threshold is not ten seconds\n' >&2; exit 1; }` — expected: exit 0. Before this Task the rule file is absent, so the command fails. After task_02 and task_03 the only budget of each of the three named files is gone, so none of them may carry a ceiling.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 3; Glossary
- `_techspec.md` → Interfaces; API Contract 7; Invariant 5; Current behavior; Rule baseline; Build Order 4; Risks & Considerations
- ADR-0271; ADR-0259; ADR-0244
