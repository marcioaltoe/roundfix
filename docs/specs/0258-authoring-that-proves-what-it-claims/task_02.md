---
task: task_02
spec: 0258-authoring-that-proves-what-it-claims
status: completed
type: backend
complexity: medium
---

# Task 02: Spec check reports a pending Task whose consumers no Verification runs

## Overview

The Spec Consistency Check gains the Tasks-stage detector of
`_techspec.md` → Interfaces and Data Models: `SC-CONSUMERS-UNPROVEN`, a gap
for every non-QA Task that has not completed and whose Consumer Declaration is
missing, malformed, names a suite no Verification runs or names a checkpoint
outside the path to the QA gate. The two fixture Tasks and the two command
test writers the prototype measured gain a `None.` declaration, and the corpus
pins gain the code at 0. It is verifiable on its own through the new detector
tests and the consuming suites.

## Requirements

1. MUST answer gap 1 of the adopted Backlog Entry of 2026-10-12,
   [Spec authoring does not prove what it claims](references/2026-10-12-spec-authoring-does-not-prove-what-it-claims.md),
   by implementing API Contract 4 of `_techspec.md` with the Consumer
   Declaration grammar of `_techspec.md` → Data Models, in a new file beside
   the existing detectors.
2. MUST match each named suite as a literal string against the Task's
   Verification commands, or against the Verification commands of the
   checkpoint the declaration names. A checkpoint MUST be a non-QA Task that
   transitively needs the declaring Task and that the QA Task transitively
   needs; when the gate is declined, the first condition alone applies.
3. MUST skip the QA Task and every Task whose status is completed, MUST run
   from the full sweep where the Task Graph is loaded, MUST register the code
   at the Tasks stage, and MUST list it as skipped when the Task Graph is
   absent (ADR-0094).
4. MUST hold every active Spec with no commit horizon, and MUST make each
   finding a gap whose fix names `## Consumers` and the write-tasks rule
   "Consumers of a shared contract" (ADR-0265).
5. MUST add a `None.` Consumer Declaration with a reason to the `task_01.md`
   of the `coverage-untasked` and `reference-unresolved` fixtures, to the Task
   the spec check Verification workspace writes, and to the Task the Delivery
   Plan workspace writes, without changing what those tests assert.
6. MUST add the code at 0 to `corpusFindingCodes`, to the corpus golden and to
   the characterization that pins it, and append one identical sentence naming
   it to the update text of both.
7. Every new top-level test MUST call `t.Parallel()` first, and each detector
   test MUST prove the gap with the declaration missing or wrong and no finding
   with it correct.

## Subtasks

- [ ] Write the detector and its five tests.
- [ ] Wire it into the full sweep and register its stage.
- [ ] Complete the fixture Tasks and the two command test writers.
- [ ] Add the code to the corpus pins.

## Acceptance Criteria

- [ ] A pending non-QA Task without `## Consumers` yields the gap; with
      `None.` and a reason it yields none.
- [ ] A suite that no Verification of the Task or of its checkpoint contains
      yields the gap; the same suite inside a Verification command yields
      none.
- [ ] A checkpoint that does not sit between the Task and the QA gate yields
      the gap.
- [ ] Completed Tasks and the QA Task are never reported.
- [ ] The existing speccheck, spec check, Delivery Plan and settlement tests
      pass unchanged apart from the declared writers.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- instruction: `.agents/skills/roundfix/references/spec.md`
- creates: `internal/speccheck/consumers.go`
- creates: `internal/speccheck/consumers_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/testdata/repo/docs/specs/coverage-untasked/task_01.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/reference-unresolved/task_01.md`
- interface: `internal/cli/spec_check_test.go`
- interface: `internal/cli/deliver_plan_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestConsumersGapFollowsTheSection|TestConsumersGapNamesASuiteNoVerificationRuns|TestConsumersCheckpointSitsBetweenTheTaskAndTheGate|TestConsumersSkipCompletedAndQATasks|TestConsumersRunAtTheTasksStage)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestConsumersGapFollowsTheSection TestConsumersGapNamesASuiteNoVerificationRuns TestConsumersCheckpointSitsBetweenTheTaskAndTheGate TestConsumersSkipCompletedAndQATasks TestConsumersRunAtTheTasksStage; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"SC-CONSUMERS-UNPROVEN": 0' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks SC-CONSUMERS-UNPROVEN\n' >&2; exit 1; }; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestSpecCheckClassifiesTheGateBoundary|TestRunSpecAuditPreservesSpecCheckBehavior|TestDeliverPlan.*)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$' ./internal/docscontract` — expected: exit 0; before this Task the five detector tests do not exist, so the command fails at its first pass check.

## References

- [_prd.md](_prd.md) — Goals 4 and 5; User Stories 4 and 5; Core Features 4 and 5; Success Metric 1; Success Metric 4
- [_techspec.md](_techspec.md) — Interfaces; Data Models; API Contract 4; API Contract 5; Invariants 1-7; Testing Approach; Build Order 2
- ADR-0265; ADR-0148; ADR-0176; ADR-0252; ADR-0257; ADR-0259

## Result

Implemented the Task 02 slice for Daemon Verification. `SC-CONSUMERS-UNPROVEN`
reads each non-completed, non-QA Task's `## Consumers` outside fenced blocks,
accepts `None.` only with a reason in the same paragraph, and validates consumer
bullets and checkpoint declarations. Named suites are literal substrings of
Verification commands in the declaring Task or its valid named checkpoints.
Checkpoint validation uses transitive dependency closures on both sides of the
checkpoint; a declined QA gate requires only the dependency on the declaring
Task. Invalid checkpoints never contribute Verification commands.

The full sweep invokes the detector after loading the Task Graph. Its registered
stage is Tasks; earlier stages and absent Task Graphs list the code as skipped.
An absent TechSpec does not prevent it from running. There is no commit-horizon
lookup. Every finding remains a gap and names `## Consumers`, `write-tasks`,
and "Consumers of a shared contract" in its fix.

The two declared fixture Tasks and both declared command workspace writers now
supply `None.` with a reason. Their existing test assertions are unchanged.
The corpus code set, golden, and characterization include the new code at 0;
both update texts append the identical sentence, and all earlier counts remain
unchanged. Every new top-level test calls `t.Parallel()` first.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Missing Consumers reports a gap; None. with a reason does not | `TestConsumersGapFollowsTheSection` covers missing, empty, malformed, and fenced declarations, blank contract/suite names, valid reasons, and valid literal suites, including wrapped entries. It also checks finding locations and fix text. |
| A suite no Verification contains reports a gap; a literal match does not | `TestConsumersGapNamesASuiteNoVerificationRuns` proves one finding per unmatched suite and no finding after commands name both suites. It checks the declaration's exact line, the deliberately coarse substring match, and that an authored command marker is never created. The checkpoint test proves an unnamed or invalid checkpoint does not satisfy a suite. |
| Checkpoints must sit between the Task and gate | `TestConsumersCheckpointSitsBetweenTheTaskAndTheGate` pairs missing, self, QA, parallel, and undeclared-gate checkpoints with repaired declarations. Direct and transitive checkpoints pass. Declined gates still require the checkpoint to need the declaring Task. A valid checkpoint must actually name the suite in its Verification. The existing graph loader independently refuses a declared gate that omits a leaf. |
| Completed and QA Tasks are never reported | `TestConsumersSkipCompletedAndQATasks` reports only the non-QA pending/in-progress/failed Task, clears its gap with a valid declaration, and leaves completed Tasks and the QA Task unreported without Consumers. |
| Existing consuming tests retain their contracts | Focused sweep/stage, spec check, Delivery Plan, and settlement tests passed under the selections below. Existing assertions were not changed. Complete consuming suites and the corpus gate remain pending Daemon Verification. |

`TestConsumersRunAtTheTasksStage` additionally proves earlier-stage skips,
full-sweep and Tasks-stage findings, strict gap promotion, a corrected
declaration, execution without a TechSpec, and absent-Graph/absent-PRD skips.

### Focused checks

The Go checks and Git whitespace check used `rtk proxy` to preserve their output:

- `go test ./internal/speccheck -run '^(TestConsumers|TestCheckCleanFixture|TestCheckCitationCoverageErrorLocations|TestStageScope)' -count=1` — exit 0 after the final detector/test edits.
- `go test ./internal/cli -run '^(TestDeliverPlanReportsAnApprovedAndABlockedSpec|TestDeliverPlanExitsZeroWhenEverySpecIsApproved|TestSpecCheckRunVerification)$' -count=1` — exit 0.
- `go test ./internal/cli -run '^(TestDeliverPlanDefaultsToEveryActiveSpec|TestDeliverPlanJSONCarriesTheSameFacts|TestSpecCheckWithoutStageIsUnchanged)$' -count=1` — exit 0.
- `go test ./internal/daemon -run '^(TestSettlementRefusesASpecConsistencyFindingTheTaskIntroduced|TestSettlementIgnoresASpecConsistencyFindingPresentAtStart)$' -count=1` — exit 0.
- Python inspection of the corpus golden, characterization, and `corpusFindingCodes` — confirmed the new code at 0, equal pinned maps and update text, exactly the requested appended sentence, and all earlier counts unchanged against `HEAD`.
- `git diff --check` — exit 0 before the Result append; repeated at handoff.

Starting evidence: the new test selection initially failed to compile because
`CodeConsumersUnproven` did not exist. The new whitespace-only contract case
then exposed a parser acceptance error; rejecting a blank contract name made
that case pass. Early focused runs also found missing active-PRD front matter
and missing QA Verification bullets in the new temporary fixture; the fixture
now supplies both, using Roundfix's derived QA command.

No declared Verification command, selected repository gate, incremental gate,
or full suite was run. Task status, Task Graph, and other Spec Task files stay
Daemon-owned. No commit, push, or pull request was made. No follow-up outside
Task 02's slice was identified.

## Carry-forward provenance

- Source Run: `run_20261013T122525Z_1c1258fdf2b03bd6`
- Source commit: `44b415ea9b8ab940dade207d747e3b3948639590`
