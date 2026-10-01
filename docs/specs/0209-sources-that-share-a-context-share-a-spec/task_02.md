---
task: task_02
spec: 0209-sources-that-share-a-context-share-a-spec
status: pending
type: backend
complexity: high
---

# Task 02: The judge asks whether an adopted and an open source share a Spec

## Overview

Spec 0205's judge asks two questions about a Spec's own text and cannot say
whether an open record belongs in the Spec. This Task adds the third judgment,
`source-grouping`: the measured question in the question file, a reader that
accepts only the Findings and Backlog Entries the grouping question may send,
the measured preparation of their text, one pending judgment per pair of an
adopted source and an open one at every stage, and the `suggested` outcome at
the measured threshold (ADR-0209). It is verifiable on its own through Spec
0205's injected transport: a fixture repository yields the suggested, clear
and skipped pairs the TechSpec names, and nothing outside the accepted
sources reaches a request.

## Requirements

1. MUST add to `internal/judge/questions.json` the `source-grouping` member whose bytes are the block under `_techspec.md` → The grouping question, after `goal-mechanism`, and MUST extend `Load` to parse it into `Questions.Grouping` and compile `source_scrub_pattern`. No other file of the package may repeat one of its values as a literal. The other members MUST stay byte-identical.
2. MUST create `internal/judge/sources.go` with `readGroupingSource`, `adoptedSources`, `openSources` and `prepareSource` as `_techspec.md` → Interfaces sketches, implementing every rule of The grouping reader and Preparing a source, including the skipped-artifact reasons `only Findings and Backlog Entries are sent`, `not a regular file in its directory` and `not English`. `Source` MUST keep its fields unexported, and only `readSpecArtifact`, `readADR` and `readGroupingSource` may build one.
3. MUST extend `PlanSpec` in `internal/judge/pairs.go` with `_techspec.md` → Grouping pairs: at every stage, after the citation and goal judgments, one pending `source-grouping` judgment per anchor and candidate, in that order, with the state `{"first", "second"}` built only from prepared source text, encoded without HTML escaping, and no grouping judgment when the Spec has no `references/_index.md`.
4. MUST extend `Run` in `internal/judge/judge.go` and the Judge Log in `internal/judge/log.go` with `_techspec.md` → Asking and outcomes: a pinned `noul` at or above `suggest_when_noul_at_least` is `suggested`, below it `clear`; an unpinned model is skipped with Spec 0205's reason and never suggested or cleared; and the grouping log line carries `"line":null`, the candidate as `target` and the `suggested` outcome. Every transport, retry, stop, skip, ceiling and log rule of Spec 0205 MUST hold for grouping judgments unchanged.
5. MUST put the new tests in `internal/judge/grouping_test.go`: `TestGroupingQuestionIsTheMeasuredOne`, `TestGroupingReaderAcceptsOnlyFindingsAndBacklogEntries`, `TestGroupingSourcesArePreparedAsMeasured`, `TestGroupingPairsEveryAdoptedSourceWithEveryOpenSource`, `TestGroupingIsPlannedAtEveryStage`, `TestGroupingSuggestsAtTheMeasuredThreshold`, `TestGroupingNeverComparesAnotherModelsAnswer`, `TestGroupingRequestsCarryOnlySourceText` and `TestGroupingJudgeLogLine`, each with a fake `http.RoundTripper`, a fixed clock and every fixture in `t.TempDir()`, covering `_techspec.md` → Testing Approach 2, 3 and 4. `TestGroupingRequestsCarryOnlySourceText` MUST run once per transport and MUST hold an Inbox Entry, a `declined` Backlog Entry, a `done` Finding, a symbolic-link Backlog Entry and a Go file with sentinel strings, and require that no sentinel and no key appears in any request body. `TestGroupingSuggestsAtTheMeasuredThreshold` MUST require `suggested` at 0.30 and `clear` at 0.29.
6. MUST update `internal/judge/questions_test.go` only to accept the third judgment and its fifth pattern, and MUST keep `TestRequestsCarryOnlySpecArtifactText` and Spec 0205's threshold tests passing with their assertions unchanged. If a test of Spec 0205 named here does not exist, or an interface `_techspec.md` → Interfaces names is missing, the Task stops and reports that Spec 0205 has not landed as its TechSpec states.

Spec 0205 creates `internal/judge/questions.json`, `questions.go`, `pairs.go`, `judge.go`, `log.go` and `questions_test.go` (`_prd.md` → Prerequisites). They are declared under `creates:` because they do not exist when this Spec is authored; this Task changes them and MUST NOT create one that Spec 0205 has not created. It reads, and does not change, Spec 0205's `source.go` and `judge_test.go`.
7. MUST NOT open a real network connection, read the process's `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write outside `t.TempDir()`.
8. MUST prove each new gate can fail. The Result MUST record one sabotage of the reader (for example accepting a `declined` Backlog Entry), one of the preparation (for example skipping the scrub) and one of the threshold (for example `>` instead of `>=`), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Add the question member and load it.
- [ ] Add the grouping reader and the preparation.
- [ ] Plan the pairs at every stage.
- [ ] Decide `suggested` or `clear` and log the line.
- [ ] Prove the data boundary on captured requests and record each sabotage.

## Acceptance Criteria

- [ ] The question file holds the TechSpec's grouping block byte for byte and loads it.
- [ ] A Spec that adopted one Finding, with two open Backlog Entries, plans two grouping judgments at the `prd`, `techspec` and default stages, and none without `references/_index.md`.
- [ ] A `noul` of 0.30 is `suggested`, 0.29 `clear`; an answer reported as `jev-1.14.0` is skipped and logged.
- [ ] A source is sent with its front matter removed, its Spec and ADR numbers scrubbed and its text cut at 1,500 code points.
- [ ] An Inbox Entry, a `declined` Backlog Entry, a `done` Finding, a symbolic link, a non-English entry and a Go file reach no request.

## Context

- instruction: `docs/adr/0209-the-judge-suggests-which-sources-share-a-spec.md`
- instruction: `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`
- creates: `internal/judge/questions.json`
- creates: `internal/judge/questions.go`
- creates: `internal/judge/pairs.go`
- creates: `internal/judge/judge.go`
- creates: `internal/judge/log.go`
- creates: `internal/judge/questions_test.go`
- creates: `internal/judge/sources.go`
- creates: `internal/judge/grouping_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGroupingQuestionIsTheMeasuredOne|TestGroupingReaderAcceptsOnlyFindingsAndBacklogEntries|TestGroupingSourcesArePreparedAsMeasured|TestGroupingPairsEveryAdoptedSourceWithEveryOpenSource|TestGroupingIsPlannedAtEveryStage|TestGroupingSuggestsAtTheMeasuredThreshold|TestGroupingNeverComparesAnotherModelsAnswer|TestGroupingRequestsCarryOnlySourceText|TestGroupingJudgeLogLine|TestQuestionFileLoads|TestRequestsCarryOnlySpecArtifactText|TestRunRaisesAtTheMeasuredThresholds|TestCitationClaimsFollowTheMeasuredExtraction|TestGoalPairsFollowTheMeasuredSelection)$" ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGroupingQuestionIsTheMeasuredOne TestGroupingReaderAcceptsOnlyFindingsAndBacklogEntries TestGroupingSourcesArePreparedAsMeasured TestGroupingPairsEveryAdoptedSourceWithEveryOpenSource TestGroupingIsPlannedAtEveryStage TestGroupingSuggestsAtTheMeasuredThreshold TestGroupingNeverComparesAnotherModelsAnswer TestGroupingRequestsCarryOnlySourceText TestGroupingJudgeLogLine TestQuestionFileLoads TestRequestsCarryOnlySpecArtifactText TestRunRaisesAtTheMeasuredThresholds TestCitationClaimsFollowTheMeasuredExtraction TestGoalPairsFollowTheMeasuredSelection; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the nine grouping tests do not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; awk 'BEGIN { fence = sprintf("%c%c%c", 96, 96, 96) } $0 == "### The grouping question" { h = 1 } h && $0 == fence "text" { f = 1; next } f && $0 == fence { exit } f { print }' docs/specs/0209-sources-that-share-a-context-share-a-spec/_techspec.md > "$tmp/want" || exit 1; test -s "$tmp/want" || { printf 'no grouping block in the TechSpec\n' >&2; exit 1; }; want="$(tr -d '\n' < "$tmp/want")"; test -f internal/judge/questions.json || { printf 'internal/judge/questions.json is missing\n' >&2; exit 1; }; tr -d '\n' < internal/judge/questions.json | grep -qF -- "$want" || { printf 'internal/judge/questions.json lacks the grouping block\n' >&2; exit 1; }` — expected: exit 0; before this Task the question file holds no `source-grouping` member, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 5; User Stories 4 and 5; Core Features 6, 7 and 8; Success Metrics 3, 4 and 5; Recorded limits
- [_techspec.md](_techspec.md) — The grouping question; Interfaces; The grouping reader; Preparing a source; Grouping pairs; Asking and outcomes; API Contract 3; API Contract 4; Testing Approach 2; Testing Approach 3; Testing Approach 4; Build Order 2
- ADR-0209; ADR-0200; ADR-0201; ADR-0089
