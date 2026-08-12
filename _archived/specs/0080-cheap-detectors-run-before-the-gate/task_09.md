---
task: task_09
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: test
complexity: medium
---

# Task 09: Let the QA harness complete the report the Daemon seeded

## Overview

Task 03 made the Daemon seed a QA Report before the Agent Session opens, named
by today's date under the `qa-report-YYYY-MM-DD[-NN].md` contract. The CLI test
harness still writes its own report at a fixed `qa-report-2026-01-01.md`, so two
reports exist: the harness's, which carries the authored verdict, and the
Daemon's seed, which carries none and wins recency. Every public journey then
reads the seed and settles `unreadable`.

The QA gate recorded this as F-001 with `Blocks-Completion` impact, reproduced by
`qa/evidence/2026-08-11-spec-0080/R04-minimal-reproduction.sh`. It reaches the
auto-push, attach, QA-only, external-root, branch, and interactive-input
journeys, and both verification tiers exit 2.

The seeding is correct and is what this Spec set out to build. The harness is
what has to change: a fake QA Agent that invents its own report filename was
only ever right while nothing else created one.

## Requirements

1. MUST make the fake QA Agent complete the report the Daemon seeded, resolving
   it by the `qa-report-YYYY-MM-DD[-NN].md` contract rather than by a fixed
   name.
2. MUST keep every authored verdict the matrix exercises — `pass`, `partial`,
   `fail`, missing, and unreadable — reaching its own assertion. The missing and
   unreadable cases must still be produced deliberately, not as a side effect of
   the seed being ignored.
3. MUST NOT change the Daemon's seeding, the report naming contract, or verdict
   semantics. The production behaviour is the one this Spec authored; only the
   harness's assumption about who names the file is wrong.
4. MUST leave the fixed-name constant removed rather than reassigned, so no
   later test can reintroduce the same assumption by reusing it.

## Subtasks

- [ ] Resolve the seeded report instead of naming one.
- [ ] Keep all five verdict cases reaching their assertions.
- [ ] Remove the fixed-name constant.

## Acceptance Criteria

- [ ] `TestRunImplementQAVerdictMatrix` passes every case.
- [ ] No test names a QA report by a hard-coded date.
- [ ] `git diff --name-only` lists only test files and this Task file.

## Bounded scope

This Task may create or modify only:

- `internal/cli/implement_test.go`
- `internal/daemon/task_engine_test.go`
- `docs/specs/0080-cheap-detectors-run-before-the-gate/task_09.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementQAVerdictMatrix$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunImplementQAVerdictMatrix'` — expected: exits 0. The matrix fails against the unchanged tree, where the `pass` case settles `unreadable`.
- `! grep -q 'qa-report-2026-01-01' internal/cli/implement_test.go` — expected: exits 0, proving the fixed name was removed rather than moved.

## References

- `_prd.md` → Goal 1.
- `task_03.md` → the seeding this harness must follow.
- `qa/qa-report-2026-08-11.md` → F-001.

## Result

Updated the in-process CLI and Daemon QA fakes to resolve the report selected
by the production `qa-report-YYYY-MM-DD[-NN].md` recency contract, require it
to match the prompt's exact `Seeded QA Report:` path, and complete that file in
place. The CLI matrix now deletes the resolved seed deliberately for its
missing case and writes malformed frontmatter deliberately for its unreadable
case. The macro ACPX fake also writes the prompt-named seed, and verdict-event
assertions select the verdict phase rather than the preceding mechanical phase.

Focused checks run after the final implementation edit:

- Pre-change signal: `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/cli -run '^TestRunImplementQAVerdictMatrix$/pass$' -count=1`
  exited 1 because the seeded `qa-report-2026-08-11.md` had no verdict while
  the fake wrote a separate fixed-name report.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/cli -run '^TestRunImplementQAVerdictMatrix/(pass|partial|fail)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/cli -run '^TestRunImplementQAVerdictMatrix/(missing_report|unreadable_verdict)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/daemon -run '^(TestMechanicalStageSeedsReportBeforeAgentSession|TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport|TestWriteMechanicalQAReportPreservesSameDayNamingAndPriorReport)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/cli -run '^(TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementInteractiveInputDoesNotChooseQAGate|TestRunImplementAutoPushOutcomeMatrix|TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict|TestRunImplementQAPromptStatesSpecTargetBranchFromRunRecord|TestAttachReplaysCompletedSpecRunReadOnly)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-gocache go test ./internal/cli -run '^TestAgentSelectionProfilesMacro/mixed_profiles_configure_validate_fallback_persist_and_stream$' -count=1`
  passed.
- `rtk git diff --check` passed.

Acceptance evidence:

1. The two focused matrix selections cover and pass all five named subtests:
   `pass`, `partial`, `fail`, `missing_report`, and `unreadable_verdict`.
2. `rtk proxy rg -n 'qa-report-[0-9]{4}-[0-9]{2}-[0-9]{2}' internal/cli/implement_test.go internal/daemon/task_engine_test.go`
   found no hard-coded dated QA Report name. The old fixed-name constant is
   removed; remaining expected names derive from the test clock and the naming
   contract.
3. `rtk git diff --name-only` listed only
   `internal/cli/implement_test.go`, `internal/daemon/task_engine_test.go`, and
   this Task file.

The commands under this Task's `## Verification` were not run; the Daemon owns
that complete selection and settlement evidence.
