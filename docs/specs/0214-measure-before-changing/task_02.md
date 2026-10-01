---
task: task_02
spec: 0214-measure-before-changing
status: pending
type: test
complexity: high
---

# Task 02: A flag-gated harness re-scores the Task acceptance question against the attempt count

## Overview

The 2026-09-30 measurement found a weak Task-lint signal on a noisy label.
This Task adds a test-only harness in `internal/judge` that reads the cause
record's per-Task attempt counts as labels, builds each implementation Task's
state from its title and acceptance criteria only, asks one fixed question
through the judge Spec 0205 delivers, and writes a record whose statistics and
verdict a second test recomputes. Nothing is sent unless the harness's flag is
passed, and nothing of it is compiled into the binary. It is verifiable on its
own: with a fake transport and fixture labels, the harness writes a record,
and the consistency test accepts it and refuses a sabotaged one.

## Requirements

1. MUST add `internal/judge/testdata/task-acceptance-question.json` byte-identical to the block under `_techspec.md` → The Task acceptance question.
2. MUST build a Task's state exactly as that section states: `task_title` from the first `# ` heading and `acceptance_criteria` from the `## Acceptance Criteria` section, each cut at its limit, reading nothing after that section ends, from a Task file that is a regular file named `task_NN.md` directly inside a Spec directory of the active Spec Root or its archive, never through a symbolic link.
3. MUST implement the harness of `_techspec.md` → The harness, steps 1 to 3, as a function the tests call with a repository root, the label bytes, a key map, a Roundfix Home, an `http.RoundTripper` and a clock. It MUST re-use Spec 0205's transport selection, request, model pin, Judge Log append and monthly ceiling from `internal/judge`, logging judgment `task-acceptance`; where the delivered names differ from Spec 0205's TechSpec, it uses the delivered ones. It MUST NOT add a question, threshold or judgment to the judge's embedded question file or to `roundfix spec judge`.
4. MUST read keys only for the two transport variables `ROUNDFIX_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY`, never `OPENROUTER_API_KEY` or `TYPESAFE_API_KEY`, and MUST write a `blocked` record with no request when neither is set.
5. MUST add `TestMeasureTaskAcceptance`, which sends requests only when `-measure-task-acceptance` is passed, then reads `-measure-labels`, `-measure-repo` (default the repository root) and `-measure-out`, takes the keys from the process environment and the Roundfix Home from `HOME`, uses `http.DefaultTransport`, and writes the record. Without that flag it MUST call `t.Skip` before reading any key, file or home.
6. MUST compute the statistics of `_techspec.md` → The harness, step 3, and the verdict of Decision rules, rule 2, in `internal/judge/task_acceptance_stats_test.go`, with a PCG generator seeded from the question file.
7. MUST add `TestTaskAcceptanceRecordIsConsistent`. With `-task-acceptance-record`, `-task-acceptance-document` and `-task-acceptance-labels` it MUST fail when a file is missing, and otherwise require `labels_sha256` to equal the labels file's SHA-256, `question_sha256` to equal the question file's, every answered Task's model to pass the judge's pin, every `state_sha256` to equal the state rebuilt from that Task file now, every statistic to equal its recomputation to four decimals, and the document's single `Verdict: ` line to equal the record's verdict. Without the flags it MUST run on fixtures under `internal/judge/testdata/task-acceptance/`, and `TestTaskAcceptanceRecordRejectsASabotagedFixture` MUST show that a fixture with one altered figure fails the same check.
8. SHOULD add no non-test file to `internal/judge`. If Spec 0205's delivered code offers no in-package way to ask one question for a given state, it MAY add one unexported function in `internal/judge/measure.go` that does only that and is called only from tests.
9. MUST prove each new gate can fail. The Result MUST record one sabotage of the state builder (for example reading past `## Acceptance Criteria` into `## Result`), one of the key source (for example reading `OPENROUTER_API_KEY`), and one of the interval (for example resampling Tasks instead of Specs), each with the test that failed, and that the code was restored.

Spec 0205 delivers `internal/judge` before this Spec (`_prd.md` → Prerequisites). When that package is absent, this Task stops and reports that Spec 0205 has not landed.

## Subtasks

- [ ] Add the question file and the state builder.
- [ ] Build labels and exclusions from a cause record.
- [ ] Ask through the judge's transports and write the record.
- [ ] Compute AUROC, the cluster interval, both baselines and the verdict.
- [ ] Add the record consistency test with its fixtures.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] With a fake transport, no request body carries a sentinel written in a fixture Task's `## Result` or in a Go source file of the fixture repository, and each body carries only `task_title` and `acceptance_criteria` in its state.
- [ ] With both keys set every request goes to OpenRouter with only the OpenRouter key; with only `ROUNDFIX_TYPESAFE_API_KEY`, to TypeSafe; with only `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY`, the record is `blocked` and the transport receives nothing.
- [ ] An answer whose model does not normalize to Jev 1.13 is `skipped` and enters no statistic, and the Judge Log gains one line per request.
- [ ] AUROC counts ties as one half, two runs with the same seed give the same interval, and each verdict branch of rule 2 is reached by a fixture.
- [ ] `TestMeasureTaskAcceptance` is skipped without its flag.
- [ ] `TestTaskAcceptanceRecordIsConsistent` passes on its fixture and fails on the sabotaged fixture.

## Context

- creates: `internal/judge/testdata/task-acceptance-question.json`
- creates: `internal/judge/task_acceptance_measure_test.go`
- creates: `internal/judge/task_acceptance_stats_test.go`
- creates: `internal/judge/testdata/task-acceptance/labels.json`
- creates: `internal/judge/testdata/task-acceptance/record.json`
- creates: `internal/judge/testdata/task-acceptance/record-sabotaged.json`
- creates: `internal/judge/testdata/task-acceptance/measurement.md`
- creates: `internal/judge/measure.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskAcceptanceQuestionFileLoads|TestTaskAcceptanceStateCarriesOnlyTitleAndCriteria|TestTaskAcceptanceRequestsCarryNoResultOrSourceText|TestTaskAcceptanceSendsEachKeyOnlyToItsEndpoint|TestTaskAcceptanceWithoutARoundfixKeyIsBlocked|TestTaskAcceptanceSkipsAnUnpinnedModel|TestTaskAcceptanceLogsEveryRequest|TestTaskAcceptanceStopsAtTheCeiling|TestAUROCCountsTiesAsHalf|TestClusterBootstrapIsDeterministic|TestTaskAcceptanceVerdictRules|TestTaskAcceptanceRecordIsConsistent|TestTaskAcceptanceRecordRejectsASabotagedFixture|TestMeasureTaskAcceptance)$" ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskAcceptanceQuestionFileLoads TestTaskAcceptanceStateCarriesOnlyTitleAndCriteria TestTaskAcceptanceRequestsCarryNoResultOrSourceText TestTaskAcceptanceSendsEachKeyOnlyToItsEndpoint TestTaskAcceptanceWithoutARoundfixKeyIsBlocked TestTaskAcceptanceSkipsAnUnpinnedModel TestTaskAcceptanceLogsEveryRequest TestTaskAcceptanceStopsAtTheCeiling TestAUROCCountsTiesAsHalf TestClusterBootstrapIsDeterministic TestTaskAcceptanceVerdictRules TestTaskAcceptanceRecordIsConsistent TestTaskAcceptanceRecordRejectsASabotagedFixture; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; printf '%s\n' "$out" | grep -q -- "--- SKIP: TestMeasureTaskAcceptance" || { printf 'the live harness ran without its flag\n' >&2; exit 1; }` — expected: exit 0; before this Task none of the named tests exists, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; awk 'BEGIN { fence = sprintf("%c%c%c", 96, 96, 96) } $0 == "### The Task acceptance question" { h = 1 } h && $0 == fence "json" { f = 1; next } f && $0 == fence { exit } f { print }' docs/specs/0214-measure-before-changing/_techspec.md > "$tmp/want.json" || exit 1; test -s "$tmp/want.json" || { printf 'no question block in the TechSpec\n' >&2; exit 1; }; cmp "$tmp/want.json" internal/judge/testdata/task-acceptance-question.json` — expected: exit 0; before this Task the question file does not exist, so `cmp` fails.

## References

- [_prd.md](_prd.md) — Goals 3 and 5; User Story 3; Core Features 6, 7 and 8; Success Metric 3; Prerequisites
- [_techspec.md](_techspec.md) — The Task acceptance question; The harness; Decision rules; Data Models; API Contract 4; Testing Approach 4; Testing Approach 5; Build Order 2
- [references/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md](references/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md)
- ADR-0215; ADR-0200; ADR-0201; ADR-0089
