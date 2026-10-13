---
task: task_04
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: completed
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

## Result

Added the Real-Time Budget detector beside the unchanged Parallel Test Package
rule. It parses test files across the repository, skips the seven declared
fixture/tool directories, resolves exact numeric duration expressions and
same-package test constants, and counts the three declared forms. Production
function declarations supply parameter names; production constants do not
supply budget values. Diagnostics group each offending file's count and ceiling
before its file/line/duration/form entries. Counts below an existing ceiling
are logged for later reduction rather than rejected.

The baseline contains **65 budgets in 21 test files**. Every ceiling equals
its measured count and records the remaining budget's reason. The three files
whose only budgets task_02/task_03 removed (`agent_test.go`,
`source_revision_test.go`, and `prune_reclaimable_test.go`) have no budgets and
no ceilings. Both new top-level tests call `t.Parallel()` first.

Added **Real-Time Budget** immediately after **Parallel Test Package** in
`CONTEXT.md`, with ADR-0271 as its source and an `_Avoid_:` line. Added the
repository rule immediately after the KISS bullet, outside setup markers.
Existing glossary and repository-guide bytes are preserved around those
additions.

Focused implementation evidence:

- Starting signal: the new rule file was absent. With the detector implemented
  and an empty ceiling table, the repository-only focused test exited 1 and
  named the remaining sites; this was the baseline measurement, not Daemon
  Verification.
- `GOCACHE=/private/tmp/roundfix-0265-task04-gocache rtk proxy go test ./internal/testfixture -run '^TestRealTimeBudgetRuleRejectsSeededViolations$' -count=1 -v`
  exited 0 for the seven required trees. Each violating tree asserts its exact
  count header and file/line/duration/form diagnostics; each clean tree asserts
  no violations.
- Final focused check:
  `GOCACHE=/private/tmp/roundfix-0265-task04-gocache rtk proxy go test -overlay /private/tmp/roundfix-0265-task04-overlay.json ./internal/testfixture -run 'RealTimeBudget|TestTask04MeasuredCeilings' -race -count=1 -v`
  exited 0. The repository rule and all 24 seeded subtests passed, including
  fields/assignments, arithmetic, grouped parameters, constant/package scope,
  dynamic durations, production constants/parameters, lower counts and all
  seven excluded directories. The temporary overlay probe asserted exact
  equality for every ceiling and logged `measured 65 budgets in 21 files;
  ceilings=21`; it also asserted the three forbidden ceiling entries absent.
  The probe and overlay exist only under `/private/tmp`, not in the diff.
- A standalone Python byte-comparison check exited 0 for glossary placement,
  ADR-0271 link, `_Avoid_:`, absence of a Spec link, additive-only documentation
  edits, rule placement after KISS, and preservation of
  `parallel_tests_test.go` against `HEAD`. Its unchanged SHA-256 is
  `530960f78bb5cb474ca8975bb4a3a2f006f2585201e7849a8219cf1d20d288e9`.
- `rtk proxy git diff --check` exited 0.

Acceptance evidence: the repository-rule check covers the repository criterion;
the exact seeded diagnostics cover the four violating and three clean trees;
the overlay measurement covers ceiling equality and the permanent rule checks
non-empty reasons; the documentation comparison covers the glossary and guide
criterion. These are focused checks only. Declared Task Verification and the
repository gates remain Daemon-owned and were not run.

Scope note: the current tree still contains the 500 ms budget in
`TestBudgetExceededRunRecordsItsReason` and the two 250 ms disposal-lock
budgets. The PRD and TechSpec explicitly leave their corrections to Spec 0259's
item branch. Their ceiling reasons name the asserted blocking work and that
pending delivery overlap; this Task changes none of those tests. When that
branch removes a counted budget, the rule logs the lower count until the next
change lowers its ceiling. No budget removed by task_01 through task_03 is
restored or given a replacement allowance.

The only pre-existing diff was the Daemon's `status: in_progress` change.
Task status, Subtasks/Acceptance checkboxes, the Task Graph, other Task files,
and the existing parallel rule were not edited. No commit, push or Pull Request
was made. Execution required approved access outside the filesystem sandbox
because its configured writable root contains a symlink.

## Carry-forward provenance

- Source Run: `run_20261013T045342Z_a78df1f6c39a67e2`
- Source commit: `87d407a94446eaf97992ca976b7a09c2052acf23`
