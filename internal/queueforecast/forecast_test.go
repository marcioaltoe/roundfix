package queueforecast

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func seconds(first, last int) []time.Duration {
	values := make([]time.Duration, 0, last-first+1)
	for i := first; i <= last; i++ {
		values = append(values, time.Duration(i)*time.Second)
	}
	return values
}

func repeated(value time.Duration, n int) []time.Duration {
	values := make([]time.Duration, n)
	for i := range values {
		values[i] = value
	}
	return values
}

func instant(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func equal[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}

func TestThePercentileIsTheNearestRank(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name             string
		n, likely, limit int
		ok               bool
	}{
		{"eleven_values", 11, 6, 10, true}, {"twenty_values", 20, 10, 17, true}, {"ten_values_give_no_pair", 10, 0, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			likely, limit, ok := Pair(seconds(1, tc.n))
			equal(t, likely, time.Duration(tc.likely)*time.Second)
			equal(t, limit, time.Duration(tc.limit)*time.Second)
			equal(t, ok, tc.ok)
		})
	}
	t.Run("the_argument_is_not_changed", func(t *testing.T) {
		values := seconds(1, 20)
		slices.Reverse(values)
		before := slices.Clone(values)
		likely, limit, ok := Pair(values)
		equal(t, likely, 10*time.Second)
		equal(t, limit, 17*time.Second)
		equal(t, ok, true)
		equal(t, values, before)
	})
	t.Run("empty", func(t *testing.T) {
		likely, limit, ok := Pair(nil)
		equal(t, likely, time.Duration(0))
		equal(t, limit, time.Duration(0))
		equal(t, ok, false)
	})
}

func TestNewBasisMeasuresTheDeliveredSpecs(t *testing.T) {
	t.Parallel()
	base := instant(t, "2026-11-16T10:00:00Z")
	at := func(s int) time.Time { return base.Add(time.Duration(s) * time.Second) }
	for _, tc := range []struct {
		name    string
		history []SpecTimes
		want    Basis
	}{
		{"start_to_merge", []SpecTimes{{Spec: "b", Start: at(10), End: at(30)}, {Spec: "a", Start: at(0), End: at(40)}, {Spec: "c", Start: at(20)}}, Basis{Delivered: 2, FirstEnd: at(30), LastEnd: at(40), StartToMerge: []time.Duration{20 * time.Second, 40 * time.Second}, StartGap: []time.Duration{10 * time.Second, 10 * time.Second}}},
		{"a_gap_needs_the_plan_on_time", []SpecTimes{{Start: at(0)}, {Plan: at(1), Start: at(20)}}, Basis{}},
		{"a_zero_plan_counts", []SpecTimes{{Start: at(0)}, {Start: at(20)}}, Basis{StartGap: []time.Duration{20 * time.Second}}},
		{"sorted_by_start_then_slug", []SpecTimes{{Spec: "c", Plan: at(0), Start: at(20)}, {Spec: "b", Plan: at(11), Start: at(10)}, {Spec: "a", Start: at(10)}}, Basis{StartGap: []time.Duration{10 * time.Second}}},
		{"first_and_last_end", []SpecTimes{{Start: at(10), End: at(70)}, {Start: at(0), End: at(100)}, {Start: at(20), End: at(50)}}, Basis{Delivered: 3, FirstEnd: at(50), LastEnd: at(100), StartToMerge: []time.Duration{30 * time.Second, 60 * time.Second, 100 * time.Second}, StartGap: []time.Duration{10 * time.Second, 10 * time.Second}}},
		{"a_spec_without_a_start_is_left_out", []SpecTimes{{End: at(30)}, {Start: at(10)}}, Basis{}},
		{"whole_seconds", []SpecTimes{{Spec: "a", Start: at(0).Add(900 * time.Millisecond), End: at(2).Add(100 * time.Millisecond)}, {Spec: "b", Plan: at(0).Add(950 * time.Millisecond), Start: at(3).Add(100 * time.Millisecond)}}, Basis{Delivered: 1, FirstEnd: at(2), LastEnd: at(2), StartToMerge: []time.Duration{2 * time.Second}, StartGap: []time.Duration{3 * time.Second}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := slices.Clone(tc.history)
			got := NewBasis(tc.history)
			equal(t, got, tc.want)
			equal(t, tc.history, before)
			if len(got.StartToMerge) > 0 {
				got.StartToMerge[0] = 0
				equal(t, tc.history, before)
			}
		})
	}
}

func TestASetOfSumsIsCutTo512Values(t *testing.T) {
	t.Parallel()
	t.Run("pairwise_sum", func(t *testing.T) {
		a, b := seconds(1, 2), []time.Duration{20 * time.Second, 10 * time.Second}
		equal(t, sum(a, b), []time.Duration{11 * time.Second, 12 * time.Second, 21 * time.Second, 22 * time.Second})
		equal(t, a, seconds(1, 2))
		equal(t, b, []time.Duration{20 * time.Second, 10 * time.Second})
	})
	for _, tc := range []struct {
		name string
		n    int
	}{{"512_values_are_kept", 512}, {"513_values", 513}, {"1024_values", 1024}} {
		t.Run(tc.name, func(t *testing.T) {
			values := seconds(1, tc.n)
			slices.Reverse(values)
			before := slices.Clone(values)
			got := cut(values)
			var want []time.Duration
			switch tc.n {
			case 512:
				want = seconds(1, 512)
			case 513:
				want = seconds(2, 513)
			case 1024:
				for i := 2; i <= 1024; i += 2 {
					want = append(want, time.Duration(i)*time.Second)
				}
			}
			equal(t, got, want)
			equal(t, values, before)
			got[0] = 0
			equal(t, values, before)
		})
	}
	t.Run("empty_sum", func(t *testing.T) { equal(t, len(sum(nil, seconds(1, 2))), 0) })
}

func TestForecastGivesEachRowItsCells(t *testing.T) {
	t.Parallel()
	base := instant(t, "2026-11-16T10:00:00Z")
	clock := func(s int) time.Time { return base.Add(time.Duration(s) * time.Second) }
	actual := func(s int) Cell { return Cell{State: Actual, Actual: clock(s)} }
	forecast := func(likely, limit int) Cell {
		return Cell{State: Forecasted, Likely: clock(likely), Limit: clock(limit)}
	}
	condition := func(text string) Cell { return Cell{State: Condition, Condition: text} }
	durations := Cell{State: Condition, HasAfter: true, LikelyAfter: 100 * time.Second, LimitAfter: 100 * time.Second}
	missing := Cell{State: None, Reason: "not enough history (10 of 11)"}
	normal := Basis{Delivered: 11, StartToMerge: repeated(100*time.Second, 11), StartGap: repeated(20*time.Second, 11)}
	short := Basis{Delivered: 10, StartToMerge: repeated(100*time.Second, 10), StartGap: repeated(20*time.Second, 11)}
	fewGaps := normal
	fewGaps.StartGap = seconds(1, 11)
	fiveLonger := normal
	fiveLonger.StartToMerge = seconds(1, 11)
	fewerGaps := normal
	fewerGaps.StartGap = repeated(20*time.Second, 10)
	head := Row{Label: "head", Planned: true, Queued: true, Started: base}
	planned := func(label string) Row { return Row{Label: label, Planned: true} }
	tcases := []struct {
		name  string
		basis Basis
		now   time.Time
		rows  []Row
		want  []RowForecast
	}{
		{"started", normal, clock(10), []Row{head}, []RowForecast{{actual(0), forecast(100, 100)}}},
		{"started_and_parked", normal, clock(10), []Row{{Started: base, Parked: true, Wait: "ignored"}}, []RowForecast{{actual(0), condition("after its park is answered")}}},
		{"past_the_usual_time", normal, clock(100), []Row{head}, []RowForecast{{actual(0), Cell{State: None, Reason: "past the usual time"}}}},
		{"not_enough_history", short, clock(10), []Row{head, {Wait: "a wait"}, {}, planned("planned")}, []RowForecast{{actual(0), missing}, {condition("a wait"), missing}, {condition("after its plan is merged"), missing}, {forecast(20, 20), missing}}},
		{"own_wait", normal, clock(10), []Row{{Planned: true, Wait: "after prerequisite 0001 is merged"}}, []RowForecast{{condition("after prerequisite 0001 is merged"), durations}}},
		{"not_planned", normal, clock(10), []Row{{}}, []RowForecast{{condition("after its plan is merged"), durations}}},
		{"one_two_and_three_behind", normal, clock(10), []Row{head, planned("one"), planned("two"), planned("three")}, []RowForecast{{actual(0), forecast(100, 100)}, {forecast(20, 20), forecast(120, 120)}, {forecast(40, 40), forecast(140, 140)}, {forecast(60, 60), forecast(160, 160)}}},
		{"fourth_behind", normal, clock(10), []Row{head, planned("one"), planned("two"), planned("three"), planned("four")}, []RowForecast{{actual(0), forecast(100, 100)}, {forecast(20, 20), forecast(120, 120)}, {forecast(40, 40), forecast(140, 140)}, {forecast(60, 60), forecast(160, 160)}, {condition("after three starts"), durations}}},
		{"too_few_gaps_longer", fewGaps, clock(10), []Row{head, planned("one"), planned("two"), planned("three")}, []RowForecast{{actual(0), forecast(100, 100)}, {condition("after head is merged"), durations}, {condition("after head is merged"), durations}, {condition("after head is merged"), durations}}},
		{"own_wait_is_not_counted", normal, clock(10), []Row{head, {Planned: true, Wait: "slot"}, {}, planned("one"), planned("two"), planned("three"), planned("four")}, []RowForecast{{actual(0), forecast(100, 100)}, {condition("slot"), durations}, {condition("after its plan is merged"), durations}, {forecast(20, 20), forecast(120, 120)}, {forecast(40, 40), forecast(140, 140)}, {forecast(60, 60), forecast(160, 160)}, {condition("after three starts"), durations}}},
		{"no_anchor_queued", normal, clock(10), []Row{{Label: "one", Planned: true, Queued: true}, {Planned: true, Wait: "slot"}, {Label: "two", Planned: true, Queued: true}, {Label: "three", Planned: true, Queued: true}}, []RowForecast{{condition("next in the queue"), durations}, {condition("slot"), durations}, {condition("after one starts"), durations}, {condition("after two starts"), durations}}},
		{"no_anchor_not_queued", normal, clock(10), []Row{planned("one"), planned("two"), {Planned: true, Queued: true}}, []RowForecast{{condition("after it is queued"), durations}, {condition("after it is queued"), durations}, {condition("after two starts"), durations}}},
		{"a_parked_row_is_no_anchor", normal, clock(10), []Row{{Started: base, Parked: true}, planned("one")}, []RowForecast{{actual(0), condition("after its park is answered")}, {condition("after it is queued"), durations}}},
		{"the_latest_start_is_the_anchor", normal, clock(10), []Row{{Label: "latest", Started: clock(5)}, head, {Started: clock(8), Parked: true}, planned("one")}, []RowForecast{{actual(5), forecast(105, 105)}, {actual(0), forecast(100, 100)}, {actual(8), condition("after its park is answered")}, {forecast(25, 25), forecast(125, 125)}}},
		{"an_end_behind_needs_eleven_deliveries", short, clock(10), []Row{head, planned("one"), planned("two"), planned("three"), planned("four")}, []RowForecast{{actual(0), missing}, {forecast(20, 20), missing}, {forecast(40, 40), missing}, {forecast(60, 60), missing}, {condition("after three starts"), missing}}},
		{"fewer_than_eleven_gaps", fewerGaps, clock(10), []Row{head, planned("one"), planned("two"), planned("three")}, []RowForecast{{actual(0), forecast(100, 100)}, {condition("after head is merged"), durations}, {condition("after head is merged"), durations}, {condition("after head is merged"), durations}}},
		{"exactly_five_longer_gaps", fewGaps, clock(6), []Row{head, planned("one")}, []RowForecast{{actual(0), forecast(100, 100)}, {forecast(9, 11), forecast(109, 111)}}},
		{"exactly_five_longer", fiveLonger, clock(6), []Row{head}, []RowForecast{{actual(0), forecast(9, 11)}}},
		{"only_four_longer", fiveLonger, clock(7), []Row{head}, []RowForecast{{actual(0), Cell{State: None, Reason: "past the usual time"}}}},
		{"whole_seconds", normal, clock(19).Add(999 * time.Millisecond), []Row{{Started: base.Add(999 * time.Millisecond)}, planned("one")}, []RowForecast{{actual(0), forecast(100, 100)}, {forecast(20, 20), forecast(120, 120)}}},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			beforeRows := slices.Clone(tc.rows)
			beforeBasis := tc.basis
			beforeBasis.StartToMerge = slices.Clone(tc.basis.StartToMerge)
			beforeBasis.StartGap = slices.Clone(tc.basis.StartGap)
			equal(t, Forecast(tc.basis, tc.rows, tc.now), tc.want)
			equal(t, tc.rows, beforeRows)
			equal(t, tc.basis, beforeBasis)
		})
	}
	t.Run("empty", func(t *testing.T) { equal(t, len(Forecast(Basis{}, nil, base)), 0) })
}
