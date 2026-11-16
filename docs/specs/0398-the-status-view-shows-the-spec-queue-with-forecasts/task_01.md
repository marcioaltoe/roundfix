---
task: task_01
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
status: completed
type: backend
complexity: high
---

# Task 01: The Queue Forecast is computed from the delivered Specs and replayed over the recorded history

## Overview

A new pure package computes every **Queue Forecast** from Spec times it is
handed: the **Forecast Basis** with its two measures, the likely time and
the limit by nearest rank, the pairwise sum of durations, and the start and
end cell of each row. It reads no file and no clock, so its tests hand it
every instant. The slice is whole on its own: a replay over the 109
delivered Specs of the backtest reproduces the counts the maintainer
accepted the method on, and the first example of the PRD comes out to the
second. No command calls the package yet; tasks 02 and 03 do.

## Requirements

1. MUST add the package `internal/queueforecast` with every declaration of
   `_techspec.md` → Interfaces for the package and for the rows and their
   forecasts, with those names, fields and signatures: the four constants
   `MinimumHistory` 11, `MinimumLonger` 5, `MaximumBehind` 3 and
   `KeepValues` 512; `SpecTimes`; `Basis`; `NewBasis`; `Pair`; `Row`;
   `CellState` with the four values `actual`, `forecast`, `condition` and
   `none`; `Cell`; `RowForecast`; and `Forecast`. The package imports the
   standard library only. Its non-test files read no file and no clock:
   they import neither `os` nor `io/fs` and never call `time.Now`,
   `time.Since` or `time.Until`. Every function returns new values and
   changes none of its arguments.
2. MUST implement the method of `_techspec.md` → API Contract 5, in whole
   seconds: every time is cut to the second before a difference is taken.
   - `NewBasis` leaves out a Spec whose `Start` is zero, sorts the others by
     `Start` and then by `Spec`, and returns `StartToMerge` and `StartGap`
     sorted from the shortest. `Delivered` is the number of Specs with an
     `End`; `FirstEnd` and `LastEnd` are the earliest and the latest `End`,
     and are zero when there is none.
   - `StartGap` holds, for each Spec after the first in that order, its
     `Start` minus the `Start` of the Spec before it, only when its `Plan`
     is zero or not after that earlier `Start`.
   - The likely value of `n` values is the value at position
     `ceil(50 × n ÷ 100)` and the limit the value at position
     `ceil(85 × n ÷ 100)`, counted from 1 over the values sorted from the
     shortest, computed with integers. `Pair` sorts a copy of its argument
     and reports `ok` only with at least `MinimumHistory` values.
   - The sum of two sets is every pairwise sum, sorted. A set of more than
     `KeepValues` values is cut to the values at positions
     `ceil(k × n ÷ 512)` for `k` from 1 to 512; a set of 512 values or
     fewer is kept whole.
3. MUST implement the cells of `_techspec.md` → API Contract 6. `Forecast`
   returns one `RowForecast` for each row, in order, and every field of a
   cell that its state does not use is zero. No row is a Backlog Entry:
   `Row` holds no kind, and the command gives a Backlog Entry its cells
   itself. Where the contract leaves a case open, these hold:
   - The anchor is the row with the latest `Started` among the rows whose
     `Started` is not zero and that are not parked.
   - A started row that is not parked takes its end from the values of
     `StartToMerge` greater than `now` minus its `Started`. With fewer than
     `MinimumHistory` values in `StartToMerge` the end is `none` with the
     reason `not enough history (<n> of 11)`, where `<n>` is `Delivered`.
     With fewer than `MinimumLonger` greater values it is `none` with the
     reason `past the usual time`. Otherwise it is `forecast`, at its
     `Started` plus the likely value and plus the limit of those greater
     values, taken at the two positions of requirement 2 without the
     minimum of 11.
   - A counted row is a row that has not started, is planned and has an
     empty `Wait`. A row with a `Wait` of its own and a row that is not
     planned are never counted.
   - With an anchor, each of the first `MaximumBehind` counted rows takes
     the gaps greater than `now` minus the anchor's `Started` as its first
     set. When `StartGap` holds fewer than `MinimumHistory` values or that
     first set holds fewer than `MinimumLonger`, its start is the condition
     `after <anchor label> is merged` and its end is the durations. This
     holds for each of those rows, not only the first.
   - Otherwise the set of the first counted row is that first set, and the
     set of each later one is the set of the counted row before it summed
     with `StartGap`, cut to `KeepValues` after each sum. Its start is
     `forecast`, at the anchor's `Started` plus the likely value and the
     limit of its set. Its end is `forecast`, at the anchor's `Started` plus
     the likely value and the limit of its set summed with `StartToMerge`
     and cut; when `StartToMerge` holds fewer than `MinimumHistory` values
     the end is `none` with the reason `not enough history (<n> of 11)`, as
     every other end is.
   - A counted row beyond the first `MaximumBehind` has the start condition
     `after <label> starts`, where `<label>` is the `Label` of the counted
     row before it, and the durations as its end.
   - With no anchor, the first counted row has the start condition
     `next in the queue` when `Queued` and `after it is queued` otherwise;
     a later one has `after <label> starts`, with the `Label` of the counted
     row before it, when `Queued`, and `after it is queued` otherwise. Each
     has the durations as its end.
   - The durations are an end whose state is `condition`, whose `Condition`
     is empty, whose `HasAfter` is true and whose `LikelyAfter` and
     `LimitAfter` are the pair of `StartToMerge`; when that pair is not
     `ok`, the end is `none` with the reason
     `not enough history (<n> of 11)`.
   - The eight condition texts and the two reasons are exactly these, with
     nothing added: `after its park is answered`,
     `after its plan is merged`, `after <label> is merged`,
     `after <label> starts`, `next in the queue`, `after it is queued`, the
     row's own `Wait`, and the empty text of the durations;
     `past the usual time` and `not enough history (<n> of 11)`.
4. MUST add the replay fixture
   `internal/queueforecast/testdata/delivery-history.json`: a JSON array of
   109 objects, each with exactly the four string keys `spec`, `plan`,
   `start` and `end`, one object for each row of the table in section 5 of
   [`references/2026-11-16-forecast-backtest.md`](references/2026-11-16-forecast-backtest.md),
   in the table's order, with the four digits and the three times as the
   table writes them. The SHA-256 of its 109 lines
   `<spec> <plan> <start> <end>`, each ended by one line feed and in that
   order, is
   `aaeb9026ec51a61891a19f226a12350bd84ba1ed032a422322a4a043901e62b3`. The
   fixture holds nothing else: no path, no session and no text. No test
   reads the Spec folder.
5. MUST add these tests in the package, each top-level test calling
   `t.Parallel()` first. Every `now` is a value the test writes; no test
   reads the clock, sleeps or waits on a timer.
   - `TestThePercentileIsTheNearestRank`, with the subtests `eleven_values`
     (of 1 s to 11 s, the likely value is 6 s and the limit 10 s),
     `twenty_values` (of 1 s to 20 s, 10 s and 17 s),
     `ten_values_give_no_pair` and `the_argument_is_not_changed`.
   - `TestNewBasisMeasuresTheDeliveredSpecs`, with the subtests
     `start_to_merge`, `a_gap_needs_the_plan_on_time` (a Spec planned after
     the earlier start adds no gap), `a_zero_plan_counts`,
     `sorted_by_start_then_slug`, `first_and_last_end`,
     `a_spec_without_a_start_is_left_out` and `whole_seconds`.
   - `TestASetOfSumsIsCutTo512Values`, with the subtests `pairwise_sum`
     (the sum of 1 s and 2 s with 10 s and 20 s is 11 s, 12 s, 21 s and
     22 s), `512_values_are_kept`, `513_values` (of 1 s to 513 s the cut
     keeps 2 s to 513 s) and `1024_values` (of 1 s to 1,024 s the cut keeps
     every even second).
   - `TestForecastGivesEachRowItsCells`, with one subtest for each case of
     requirement 3, each asserting both cells whole: `started`,
     `started_and_parked`, `past_the_usual_time`, `not_enough_history`,
     `own_wait`, `not_planned`, `one_two_and_three_behind`, `fourth_behind`,
     `too_few_gaps_longer`, `own_wait_is_not_counted`, `no_anchor_queued`,
     `no_anchor_not_queued`, `a_parked_row_is_no_anchor`,
     `the_latest_start_is_the_anchor` and
     `an_end_behind_needs_eleven_deliveries`.
   - `TestTheForecastOfTheRecordedMoment`: with the fixture as the history
     plus one Spec `0397` planned at 2026-11-16T11:18:17Z, started at
     2026-11-16T11:30:44Z and not delivered, `now` 2026-11-16T13:36:00Z, and
     three rows (`0397` planned, queued and started at that instant; `0398`
     planned and not queued; `0399` not planned), the cells are:
     - `0397`: start `actual` 2026-11-16T11:30:44Z; end `forecast`, likely
       2026-11-16T14:59:09Z and limit 2026-11-16T17:12:44Z;
     - `0398`: start `forecast`, likely 2026-11-16T14:20:15Z and limit
       2026-11-16T18:02:58Z; end `forecast`, likely 2026-11-16T17:07:24Z and
       limit 2026-11-16T22:40:19Z;
     - `0399`: start the condition `after its plan is merged`; end the
       durations 6,990 s and 15,711 s.
   - `TestTheReplayReproducesTheBacktest`, in `replay_test.go`, with these
     subtests. A case's history is every fixture Spec whose start is not
     after the case's instant, with its end kept only when that end is not
     after the instant.
     - `fixture`: the file holds 109 objects with only the four keys, in
       the order of their starts, and the digest of requirement 4.
     - `whole_history`: `NewBasis` over the 109 Specs holds 109 times from
       start to merge and 74 start gaps, and `Pair` gives 6,990 s and
       15,711 s for the first and 7,365 s and 13,963 s for the second.
     - `at_its_start`: for each fixture Spec, the instant is its start and
       the one row is that Spec, planned, queued and started at its start.
       The Spec is a case when its history holds at least 11 times from
       start to merge. There are 98 cases; the Spec's end is not after the
       limit of the end cell in 80 and not after its likely time in 46.
     - `one_behind`, `two_behind` and `three_behind`, for a depth of 1, 2
       and 3: in the order of the starts, for each Spec that has at least
       that many Specs before it, the head is the Spec that many places
       before it and the instant is the head's start. It is a case when
       each of the Specs after the head, up to and including it, has a plan
       not after that instant, and the history holds at least 11 times from
       start to merge and at least 11 start gaps. The rows are the head,
       started at the instant, and then one planned row that has not
       started for each of those Specs; the cells compared are those of the
       last row, with the Spec's own start and end. The cases are 63, 39
       and 19. The start is not after its limit in 55, 36 and 19 and not
       after its likely time in 35, 21 and 8; the end is not after its
       limit in 53, 36 and 17 and not after its likely time in 27, 20 and
       10.
6. MUST keep the package building for every platform:
   `GOOS=windows go build -buildvcs=false ./internal/queueforecast` and the
   same for `linux`. No file of the package is platform-specific.
7. MUST change no file outside the four this Task creates, add no module
   dependency and change no Governed Path.

## Subtasks

- [ ] Add the package with the basis, the pair, the pairwise sum and the cut.
- [ ] Add the cells of every row case, with the condition texts and the reasons.
- [ ] Add the fixture from the table of the backtest and check its digest.
- [ ] Add the table tests of the method and of each row case.
- [ ] Add the replay and the test of the recorded moment.

## Acceptance Criteria

- [ ] Over the fixture of 109 delivered Specs the replay reproduces the recorded counts: 80 and 46 of 98 ends forecast at the start; 55 and 35 of 63 starts and 53 and 27 of 63 ends one row behind; 36 and 21 of 39 starts and 36 and 20 of 39 ends two rows behind; 19 and 8 of 19 starts and 17 and 10 of 19 ends three rows behind.
- [ ] A change to the percentile position, to the minimum of 11 or to the pairwise sum changes at least one of those counts, and a change to the cut at 512 fails the test of the cut.
- [ ] With fewer than 11 delivered Specs every end reads `not enough history (<n> of 11)` and no cell holds a time or a duration; a started Spec with fewer than 5 longer deliveries reads `past the usual time`.
- [ ] A planned Spec one, two or three rows behind the Spec in progress that started last receives a start and an end as pairs of times; the fourth, a row with a wait of its own, a row that is not planned and every row with no Spec in progress ahead receive a named condition and two durations, and no clock time.
- [ ] The package reads no file and no clock, imports the standard library only and builds for Windows and Linux.

## Context

- instruction: `docs/adr/0409-a-queue-forecast-is-a-likely-time-and-a-limit-from-the-repository-own-deliveries-or-a-named-condition.md`
- creates: `internal/queueforecast/forecast.go`
- creates: `internal/queueforecast/forecast_test.go`
- creates: `internal/queueforecast/replay_test.go`
- creates: `internal/queueforecast/testdata/delivery-history.json`

## Consumers

None. The Task adds a package that no code calls yet; tasks 02 and 03 add its callers and run the suites that consume it.

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestThePercentileIsTheNearestRank|TestNewBasisMeasuresTheDeliveredSpecs|TestASetOfSumsIsCutTo512Values|TestForecastGivesEachRowItsCells)$' ./internal/queueforecast 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThePercentileIsTheNearestRank/eleven_values TestThePercentileIsTheNearestRank/twenty_values TestThePercentileIsTheNearestRank/ten_values_give_no_pair TestNewBasisMeasuresTheDeliveredSpecs/a_gap_needs_the_plan_on_time TestNewBasisMeasuresTheDeliveredSpecs/sorted_by_start_then_slug TestNewBasisMeasuresTheDeliveredSpecs/whole_seconds TestASetOfSumsIsCutTo512Values/pairwise_sum TestASetOfSumsIsCutTo512Values/513_values TestASetOfSumsIsCutTo512Values/1024_values TestForecastGivesEachRowItsCells/started TestForecastGivesEachRowItsCells/started_and_parked TestForecastGivesEachRowItsCells/past_the_usual_time TestForecastGivesEachRowItsCells/not_enough_history TestForecastGivesEachRowItsCells/own_wait TestForecastGivesEachRowItsCells/not_planned TestForecastGivesEachRowItsCells/one_two_and_three_behind TestForecastGivesEachRowItsCells/fourth_behind TestForecastGivesEachRowItsCells/too_few_gaps_longer TestForecastGivesEachRowItsCells/own_wait_is_not_counted TestForecastGivesEachRowItsCells/no_anchor_queued TestForecastGivesEachRowItsCells/no_anchor_not_queued TestForecastGivesEachRowItsCells/a_parked_row_is_no_anchor TestForecastGivesEachRowItsCells/the_latest_start_is_the_anchor TestForecastGivesEachRowItsCells/an_end_behind_needs_eleven_deliveries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; for f in internal/queueforecast/*.go; do case "$f" in *_test.go) continue ;; esac; if grep -qE 'time[.](Now|Since|Until)[(]|"os"|"io/fs"' "$f"; then echo "the package reads a clock or a file: $f" >&2; exit 1; fi; done; GOOS=windows go build -buildvcs=false ./internal/queueforecast && GOOS=linux go build -buildvcs=false ./internal/queueforecast` — expected: exit 0. Before this Task the package does not exist, so the test command fails; afterwards the nearest rank, the two measures, the pairwise sum, the cut and each row case hold by name, no non-test file reads a clock or a file, and the package builds for Windows and Linux.

- `test -f internal/queueforecast/testdata/delivery-history.json || { echo 'missing: the replay fixture' >&2; exit 1; }; out="$(go test -count=1 -timeout 15m -v -run '^(TestTheReplayReproducesTheBacktest|TestTheForecastOfTheRecordedMoment)$' ./internal/queueforecast 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheReplayReproducesTheBacktest/fixture TestTheReplayReproducesTheBacktest/whole_history TestTheReplayReproducesTheBacktest/at_its_start TestTheReplayReproducesTheBacktest/one_behind TestTheReplayReproducesTheBacktest/two_behind TestTheReplayReproducesTheBacktest/three_behind TestTheForecastOfTheRecordedMoment; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 15m ./internal/queueforecast` — expected: exit 0. Before this Task the fixture does not exist; afterwards it carries the 109 Specs of the backtest by its digest, the replay reproduces every recorded count, the recorded moment gives its instants to the second, and the whole package passes.

## References

- `_prd.md` → Goals; User Stories 2, 3, 4, 6 and 7; Core Features 4, 5, 6, 7, 8 and 15; Success Metric 2; Validation
- `_techspec.md` → Queue forecast package; Interfaces; API Contract 5; API Contract 6; Data Models; Testing Approach; Build Order 1
- [`references/2026-11-16-forecast-backtest.md`](references/2026-11-16-forecast-backtest.md), sections 2, 3 and 5
- ADR-0409

## Result

Skills loaded: implement-task, coding-guidelines, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Added the pure `queueforecast` package with the declared interfaces, whole-second
history measures, nearest-rank pairs, bounded pairwise sums, and the specified
row-cell rules. Inputs are copied before sorting; cell comparisons assert every
field, including unused zero fields. The only production imports are `fmt`,
`slices`, and `time`. No caller is added.

The replay fixture contains exactly the 109 four-key records copied from the
backtest table. Its canonical-line SHA-256 is
`aaeb9026ec51a61891a19f226a12350bd84ba1ed032a422322a4a043901e62b3`.
Tests read only that package-local fixture, never the Spec folder.

Focused implementation evidence:

- `go test -race -count=1 -timeout=2m -v ./internal/queueforecast` exited 0.
  All six prescribed top-level tests and their required subtests ran, including
  the recorded moment, exact cell comparisons, defensive-copy checks, and the
  replay. After strengthening the tied-start ordering case and adding an
  exactly-five-longer-gaps boundary case,
  `go test -race -count=1 -timeout=2m ./internal/queueforecast` exited 0 again.
- The whole-history test measures 109 deliveries and 74 eligible start gaps,
  with pairs 6,990/15,711 seconds and 7,365/13,963 seconds respectively.
  The replay reports the following counts (limit/likely):

  | Case | Cases | Start | End |
  | --- | --- | --- | --- |
  | At its start | 98 | actual | 80/46 |
  | One behind | 63 | 55/35 | 53/27 |
  | Two behind | 39 | 36/21 | 36/20 |
  | Three behind | 19 | 19/8 | 17/10 |

- Four isolated `go test -overlay=<temporary overlay> -count=1 -timeout=2m -v
  -run <focused pattern> ./internal/queueforecast` probes each exited 1 with
  assertion failures, as intended; repository source was never mutated.
  Changing median rank from 50 to 51 changes the one-behind counts to
  `63,55,37,53,28`. Changing `MinimumHistory` to 10 changes the at-start counts
  to `99,81,47`. Adding 100 hours to each pairwise sum changes one-behind
  counts to `63,55,35,63,63`. Changing `KeepValues` to 511 fails all three
  named cutoff cases. The actual assertion transcripts were inspected, not
  merely the exit codes. Logs and overlays are under
  `/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/queueforecast-mutations-0czpokdv/`.
- `not_enough_history` and `an_end_behind_needs_eleven_deliveries` assert empty
  end times and durations with `not enough history (10 of 11)` for the
  relevant started, waiting, unplanned, and counted rows. As requirement 3
  specifies, actual starts and eligible gap-based start forecasts remain
  available independently of end history. `past_the_usual_time`,
  `exactly_five_longer`, and `only_four_longer` exercise strict greater-than
  filtering and the five-value boundary.
- `one_two_and_three_behind`, `fourth_behind`, `own_wait_is_not_counted`, the
  no-anchor cases, the parked-anchor case, and the latest-start case assert
  the complete cells and exact condition text. `too_few_gaps_longer` and
  `fewer_than_eleven_gaps` check fallback for each of the first three rows;
  `exactly_five_longer_gaps` verifies that five retained gaps suffice.
- Focused cross-compilation used `go test -c -buildvcs=false -o <temporary
  binary> ./internal/queueforecast` with `GOOS=windows` and then `GOOS=linux`,
  each with `GOARCH=amd64 CGO_ENABLED=0`; both exited 0 after the final test
  edits. This compiles the production package and its tests without running
  either target binary.
- `go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}'
  ./internal/queueforecast` returned only `roundfix/internal/queueforecast`.
  A focused source inspection confirmed no production file/clock access, no
  test clock/timer access or Spec-folder read, no platform-specific files,
  and `t.Parallel()` first in all six top-level tests. `gofmt -l` returned
  no paths and `git diff --check` exited 0.

The first focused run exposed an arithmetic typo in the new
`not_enough_history` expectation: it expected a 120-second start for a
20-second gap. The test expectation was corrected to the specified 20-second
start; the production calculation was unchanged, and the focused rerun exited 0.

Scope inspection shows only the four declared new package files and this
Task's Result; the pre-existing Daemon change from `pending` to `in_progress`
is preserved. No Task status, Task Graph, other Task, module dependency,
Governed Path, or caller was changed by this implementation.

Declared Verification and repository Verification were not run; the Daemon
owns those checks and settlement. No commit, push, or pull request was made.
No follow-up work was identified within this slice.
