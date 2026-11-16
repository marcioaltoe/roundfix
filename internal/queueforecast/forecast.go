// Package queueforecast computes queue forecasts from supplied history and instants.
package queueforecast

import (
	"fmt"
	"slices"
	"time"
)

const (
	MinimumHistory = 11
	MinimumLonger  = 5
	MaximumBehind  = 3
	KeepValues     = 512
)

// SpecTimes records a plan, first Implement Run and delivery merge.
// A zero End means not delivered; a zero Plan predates the history window.
type SpecTimes struct {
	Spec             string
	Plan, Start, End time.Time
}

// Basis holds the sorted start-to-merge durations and eligible start gaps.
type Basis struct {
	Delivered              int
	FirstEnd, LastEnd      time.Time
	StartToMerge, StartGap []time.Duration
}

// NewBasis measures history without changing it. Specs without starts are omitted.
func NewBasis(history []SpecTimes) Basis {
	ordered := make([]SpecTimes, 0, len(history))
	for _, spec := range history {
		if spec.Start.IsZero() {
			continue
		}
		spec.Plan = spec.Plan.Truncate(time.Second)
		spec.Start = spec.Start.Truncate(time.Second)
		spec.End = spec.End.Truncate(time.Second)
		ordered = append(ordered, spec)
	}
	slices.SortFunc(ordered, func(a, b SpecTimes) int {
		if c := a.Start.Compare(b.Start); c != 0 {
			return c
		}
		if a.Spec < b.Spec {
			return -1
		}
		if a.Spec > b.Spec {
			return 1
		}
		return 0
	})
	var basis Basis
	for i, spec := range ordered {
		if !spec.End.IsZero() {
			basis.Delivered++
			basis.StartToMerge = append(basis.StartToMerge, spec.End.Sub(spec.Start))
			if basis.FirstEnd.IsZero() || spec.End.Before(basis.FirstEnd) {
				basis.FirstEnd = spec.End
			}
			if basis.LastEnd.IsZero() || spec.End.After(basis.LastEnd) {
				basis.LastEnd = spec.End
			}
		}
		if i > 0 && (spec.Plan.IsZero() || !spec.Plan.After(ordered[i-1].Start)) {
			basis.StartGap = append(basis.StartGap, spec.Start.Sub(ordered[i-1].Start))
		}
	}
	slices.Sort(basis.StartToMerge)
	slices.Sort(basis.StartGap)
	return basis
}

// Pair returns the median and 85th percentile by nearest rank, with enough history.
func Pair(values []time.Duration) (likely, limit time.Duration, ok bool) {
	if len(values) < MinimumHistory {
		return 0, 0, false
	}
	ordered := slices.Clone(values)
	slices.Sort(ordered)
	likely, limit = rankedPair(ordered)
	return likely, limit, true
}

func rankedPair(ordered []time.Duration) (time.Duration, time.Duration) {
	n := len(ordered)
	return ordered[(50*n+99)/100-1], ordered[(85*n+99)/100-1]
}

// sum retains the empirical distribution of every pairwise sum, bounded in size.
func sum(a, b []time.Duration) []time.Duration {
	values := make([]time.Duration, 0, len(a)*len(b))
	for _, x := range a {
		for _, y := range b {
			values = append(values, x+y)
		}
	}
	return cut(values)
}

func cut(values []time.Duration) []time.Duration {
	ordered := slices.Clone(values)
	slices.Sort(ordered)
	if len(ordered) <= KeepValues {
		return ordered
	}
	kept := make([]time.Duration, KeepValues)
	for k := 1; k <= KeepValues; k++ {
		kept[k-1] = ordered[(k*len(ordered)+KeepValues-1)/KeepValues-1]
	}
	return kept
}

// Row is a Spec in queue order; Backlog Entries are handled by the caller.
type Row struct {
	Label   string
	Planned bool
	Queued  bool
	Started time.Time
	Parked  bool
	Wait    string
}

// CellState identifies which fields of a Cell carry a value.
type CellState string

const (
	Actual     CellState = "actual"
	Forecasted CellState = "forecast"
	Condition  CellState = "condition"
	None       CellState = "none"
)

// Cell is an actual instant, two forecast instants, a condition or an absence.
// Conditional ends can carry durations after their start instead of instants.
type Cell struct {
	State                   CellState
	Actual                  time.Time
	Likely, Limit           time.Time
	Condition               string
	LikelyAfter, LimitAfter time.Duration
	HasAfter                bool
	Reason                  string
}

// RowForecast carries the start and end cells of a Spec.
type RowForecast struct{ Start, End Cell }

// Forecast computes one forecast per row without changing the basis or rows.
func Forecast(basis Basis, rows []Row, now time.Time) []RowForecast {
	now = now.Truncate(time.Second)
	anchor := -1
	for i, row := range rows {
		if !row.Started.IsZero() && !row.Parked && (anchor < 0 || row.Started.Truncate(time.Second).After(rows[anchor].Started.Truncate(time.Second))) {
			anchor = i
		}
	}
	missing := Cell{State: None, Reason: fmt.Sprintf("not enough history (%d of %d)", basis.Delivered, MinimumHistory)}
	durationEnd := missing
	if likely, limit, ok := Pair(basis.StartToMerge); ok {
		durationEnd = Cell{State: Condition, HasAfter: true, LikelyAfter: likely, LimitAfter: limit}
	}
	var gaps []time.Duration
	var anchorStart time.Time
	if anchor >= 0 {
		anchorStart = rows[anchor].Started.Truncate(time.Second)
		gaps = longer(basis.StartGap, now.Sub(anchorStart))
	}
	gapsOK := len(basis.StartGap) >= MinimumHistory && len(gaps) >= MinimumLonger
	forecasts := make([]RowForecast, len(rows))
	counted := 0
	previous := ""
	for i, row := range rows {
		start := row.Started.Truncate(time.Second)
		f := RowForecast{End: durationEnd}
		switch {
		case !row.Started.IsZero():
			f.Start = Cell{State: Actual, Actual: start}
			switch {
			case row.Parked:
				f.End = Cell{State: Condition, Condition: "after its park is answered"}
			case len(basis.StartToMerge) < MinimumHistory:
				f.End = missing
			default:
				remaining := longer(basis.StartToMerge, now.Sub(start))
				if len(remaining) < MinimumLonger {
					f.End = Cell{State: None, Reason: "past the usual time"}
				} else {
					f.End = at(start, remaining)
				}
			}
		case row.Wait != "":
			f.Start = Cell{State: Condition, Condition: row.Wait}
		case !row.Planned:
			f.Start = Cell{State: Condition, Condition: "after its plan is merged"}
		default:
			counted++
			switch {
			case anchor < 0:
				condition := "after it is queued"
				if row.Queued {
					condition = "next in the queue"
					if counted > 1 {
						condition = "after " + previous + " starts"
					}
				}
				f.Start = Cell{State: Condition, Condition: condition}
			case counted > MaximumBehind:
				f.Start = Cell{State: Condition, Condition: "after " + previous + " starts"}
			case !gapsOK:
				f.Start = Cell{State: Condition, Condition: "after " + rows[anchor].Label + " is merged"}
			default:
				if counted > 1 {
					gaps = sum(gaps, basis.StartGap)
				}
				f.Start = at(anchorStart, gaps)
				if len(basis.StartToMerge) < MinimumHistory {
					f.End = missing
				} else {
					f.End = at(anchorStart, sum(gaps, basis.StartToMerge))
				}
			}
			previous = row.Label
		}
		forecasts[i] = f
	}
	return forecasts
}

func longer(values []time.Duration, elapsed time.Duration) []time.Duration {
	remaining := make([]time.Duration, 0, len(values))
	for _, value := range values {
		if value > elapsed {
			remaining = append(remaining, value)
		}
	}
	slices.Sort(remaining)
	return remaining
}

func at(start time.Time, ordered []time.Duration) Cell {
	likely, limit := rankedPair(ordered)
	return Cell{State: Forecasted, Likely: start.Add(likely), Limit: start.Add(limit)}
}
