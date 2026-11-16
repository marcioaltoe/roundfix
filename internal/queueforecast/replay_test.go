package queueforecast

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const historyDigest = "aaeb9026ec51a61891a19f226a12350bd84ba1ed032a422322a4a043901e62b3"

func deliveryHistory(t *testing.T) ([]SpecTimes, []map[string]string) {
	t.Helper()
	data, err := os.ReadFile("testdata/delivery-history.json")
	if err != nil {
		t.Fatal(err)
	}
	var records []map[string]string
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatal(err)
	}
	history := make([]SpecTimes, len(records))
	for i, r := range records {
		history[i] = SpecTimes{Spec: r["spec"], Plan: instant(t, r["plan"]), Start: instant(t, r["start"]), End: instant(t, r["end"])}
	}
	return history, records
}

// historyAt reveals only starts and deliveries observed by the supplied instant.
func historyAt(history []SpecTimes, now time.Time) Basis {
	var visible []SpecTimes
	for _, spec := range history {
		if spec.Start.After(now) {
			continue
		}
		if spec.End.After(now) {
			spec.End = time.Time{}
		}
		visible = append(visible, spec)
	}
	return NewBasis(visible)
}

func TestTheReplayReproducesTheBacktest(t *testing.T) {
	t.Parallel()
	history, records := deliveryHistory(t)
	t.Run("fixture", func(t *testing.T) {
		equal(t, len(records), 109)
		var lines strings.Builder
		for i, r := range records {
			equal(t, len(r), 4)
			for _, key := range []string{"spec", "plan", "start", "end"} {
				if _, ok := r[key]; !ok {
					t.Fatalf("record %d missing %q", i, key)
				}
			}
			if len(r["spec"]) != 4 || strings.Trim(r["spec"], "0123456789") != "" {
				t.Fatalf("invalid Spec: %q", r["spec"])
			}
			for _, key := range []string{"plan", "start", "end"} {
				equal(t, instant(t, r[key]).Format(time.RFC3339), r[key])
			}
			if i > 0 && history[i].Start.Before(history[i-1].Start) {
				t.Fatalf("record %d is out of start order", i)
			}
			fmt.Fprintf(&lines, "%s %s %s %s\n", r["spec"], r["plan"], r["start"], r["end"])
		}
		equal(t, fmt.Sprintf("%x", sha256.Sum256([]byte(lines.String()))), historyDigest)
	})
	t.Run("whole_history", func(t *testing.T) {
		basis := NewBasis(history)
		equal(t, basis.Delivered, 109)
		equal(t, len(basis.StartToMerge), 109)
		equal(t, len(basis.StartGap), 74)
		likely, limit, ok := Pair(basis.StartToMerge)
		equal(t, ok, true)
		equal(t, likely, 6990*time.Second)
		equal(t, limit, 15711*time.Second)
		likely, limit, ok = Pair(basis.StartGap)
		equal(t, ok, true)
		equal(t, likely, 7365*time.Second)
		equal(t, limit, 13963*time.Second)
	})
	t.Run("at_its_start", func(t *testing.T) {
		cases, limitHeld, likelyHeld := 0, 0, 0
		for _, spec := range history {
			basis := historyAt(history, spec.Start)
			if len(basis.StartToMerge) < MinimumHistory {
				continue
			}
			end := Forecast(basis, []Row{{Label: spec.Spec, Planned: true, Queued: true, Started: spec.Start}}, spec.Start)[0].End
			equal(t, end.State, Forecasted)
			cases++
			if !spec.End.After(end.Limit) {
				limitHeld++
			}
			if !spec.End.After(end.Likely) {
				likelyHeld++
			}
		}
		equal(t, []int{cases, limitHeld, likelyHeld}, []int{98, 80, 46})
		t.Logf("cases=%d end limit=%d likely=%d", cases, limitHeld, likelyHeld)
	})
	for _, tc := range []struct {
		name  string
		depth int
		want  []int
	}{
		{"one_behind", 1, []int{63, 55, 35, 53, 27}},
		{"two_behind", 2, []int{39, 36, 21, 36, 20}},
		{"three_behind", 3, []int{19, 19, 8, 17, 10}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counts := make([]int, 5)
			for i := tc.depth; i < len(history); i++ {
				spec, head := history[i], history[i-tc.depth]
				basis := historyAt(history, head.Start)
				if len(basis.StartToMerge) < MinimumHistory || len(basis.StartGap) < MinimumHistory {
					continue
				}
				rows := []Row{{Label: head.Spec, Planned: true, Queued: true, Started: head.Start}}
				eligible := true
				for _, behind := range history[i-tc.depth+1 : i+1] {
					if behind.Plan.After(head.Start) {
						eligible = false
						break
					}
					rows = append(rows, Row{Label: behind.Spec, Planned: true, Queued: true})
				}
				if !eligible {
					continue
				}
				last := Forecast(basis, rows, head.Start)[tc.depth]
				equal(t, last.Start.State, Forecasted)
				equal(t, last.End.State, Forecasted)
				counts[0]++
				if !spec.Start.After(last.Start.Limit) {
					counts[1]++
				}
				if !spec.Start.After(last.Start.Likely) {
					counts[2]++
				}
				if !spec.End.After(last.End.Limit) {
					counts[3]++
				}
				if !spec.End.After(last.End.Likely) {
					counts[4]++
				}
			}
			equal(t, counts, tc.want)
			t.Logf("cases=%d start limit=%d likely=%d end limit=%d likely=%d", counts[0], counts[1], counts[2], counts[3], counts[4])
		})
	}
}

func TestTheForecastOfTheRecordedMoment(t *testing.T) {
	t.Parallel()
	history, _ := deliveryHistory(t)
	started := instant(t, "2026-11-16T11:30:44Z")
	history = append(history, SpecTimes{Spec: "0397", Plan: instant(t, "2026-11-16T11:18:17Z"), Start: started})
	rows := []Row{{Label: "0397", Planned: true, Queued: true, Started: started}, {Label: "0398", Planned: true}, {Label: "0399"}}
	want := []RowForecast{
		{Start: Cell{State: Actual, Actual: started}, End: Cell{State: Forecasted, Likely: instant(t, "2026-11-16T14:59:09Z"), Limit: instant(t, "2026-11-16T17:12:44Z")}},
		{Start: Cell{State: Forecasted, Likely: instant(t, "2026-11-16T14:20:15Z"), Limit: instant(t, "2026-11-16T18:02:58Z")}, End: Cell{State: Forecasted, Likely: instant(t, "2026-11-16T17:07:24Z"), Limit: instant(t, "2026-11-16T22:40:19Z")}},
		{Start: Cell{State: Condition, Condition: "after its plan is merged"}, End: Cell{State: Condition, HasAfter: true, LikelyAfter: 6990 * time.Second, LimitAfter: 15711 * time.Second}},
	}
	equal(t, Forecast(NewBasis(history), rows, instant(t, "2026-11-16T13:36:00Z")), want)
}
