package gatetime

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"roundfix/internal/store"
	"sort"
	"strings"
	"time"
)

type Sample struct {
	At      time.Time `json:"at"`
	Seconds int       `json:"seconds"`
}
type Source map[string][]Sample

func Windows(doc Document, source Source) ([]string, bool) {
	var lines []string
	within := true
	for _, g := range doc.Gates {
		if g.Window.Size == 0 {
			continue
		}
		samples := append([]Sample{}, source[g.ID]...)
		sort.SliceStable(samples, func(i, j int) bool { return samples[i].At.Before(samples[j].At) })
		var values []int
		for _, s := range samples {
			if g.Window.Since != nil && s.At.Before(*g.Window.Since) {
				continue
			}
			values = append(values, s.Seconds)
			if len(values) == g.Window.Size {
				break
			}
		}
		prefix := fmt.Sprintf("gate times: %s, %d of %d passing: ", g.Name, len(values), g.Window.Size)
		if len(values) < g.Window.Size {
			lines = append(lines, prefix+"window open")
			continue
		}
		sort.Ints(values)
		median := values[len(values)/2]
		if len(values)%2 == 0 {
			median = (values[len(values)/2-1] + median) / 2
		}
		p95 := values[int(math.Ceil(float64(len(values))*.95))-1]
		limit := g.Approved
		if doc.Rule != nil && g.Number != nil {
			limit = *g.Number
		}
		state := "within"
		if median > limit.MedianSeconds || p95 > limit.P95Seconds {
			state = "over"
			if doc.Rule != nil {
				within = false
			}
		}
		lines = append(lines, fmt.Sprintf("%smedian %ds of %ds, p95 %ds of %ds, %s", prefix, median, limit.MedianSeconds, p95, limit.P95Seconds, state))
	}
	return lines, within
}
func readOnlyDSN(path string) string {
	return (&url.URL{Scheme: "file", Path: filepath.ToSlash(path), RawQuery: "mode=ro"}).String()
}
func openRunDatabase(path string) (*sql.DB, error) {
	return sql.Open("sqlite", readOnlyDSN(path))
}
func settlementSamples(db *sql.DB, root string, since *time.Time) ([]Sample, error) {
	sinceText := ""
	if since != nil {
		sinceText = since.UTC().Format(time.RFC3339Nano)
	}
	rows, err := db.Query(`SELECT e.run_id,e.created_at,e.payload FROM run_events e JOIN runs r ON r.id=e.run_id WHERE e.kind='daemon.verification' AND r.git_root=? AND r.created_at>=? ORDER BY e.created_at,e.cursor`, root, sinceText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	starts := map[string]time.Time{}
	var result []Sample
	for rows.Next() {
		var run, at, payload string
		if err := rows.Scan(&run, &at, &payload); err != nil {
			return nil, err
		}
		var p struct {
			Command, Phase, Task string
			Attempt              int
			Retry                int
		}
		if err := json.Unmarshal([]byte(payload), &p); err != nil {
			return nil, err
		}
		if p.Command != "make verify-changed" {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339Nano, at)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%s/%s/%d/%d", run, p.Task, p.Attempt, p.Retry)
		switch p.Phase {
		case "started":
			starts[key] = timestamp
		case "command_passed", "command-passed":
			if start, ok := starts[key]; ok {
				result = append(result, Sample{At: timestamp, Seconds: int(math.Ceil(timestamp.Sub(start).Seconds()))})
				delete(starts, key)
			}
		}
	}
	return result, rows.Err()
}
func readSources(doc Document, latest bool) (Source, error) {
	source := Source{}
	active := false
	for _, g := range doc.Gates {
		active = active || g.Window.Since != nil
	}
	if !active && !latest {
		return source, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	root, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, err
	}
	db, err := openRunDatabase(store.DatabasePath(home))
	if err != nil {
		return nil, err
	}
	var since *time.Time
	for _, g := range doc.Gates {
		if g.ID == "task-settlement" {
			since = g.Window.Since
		}
	}
	samples, queryErr := settlementSamples(db, strings.TrimSpace(string(root)), since)
	closeErr := db.Close()
	if queryErr != nil {
		return nil, queryErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	source["task-settlement"] = samples
	data, err := exec.Command("gh", "run", "list", "--workflow", "CI | Verify", "--event", "pull_request", "--status", "success", "--limit", "1000", "--json", "databaseId,createdAt").Output()
	if err != nil {
		return nil, err
	}
	var runs []ciRun
	if err = json.Unmarshal(data, &runs); err != nil {
		return nil, err
	}
	var pullRequest *Gate
	for i := range doc.Gates {
		if doc.Gates[i].ID == "pull-request" {
			pullRequest = &doc.Gates[i]
		}
	}
	runs = runsInWindow(runs, pullRequest)
	for _, run := range runs {
		data, err := exec.Command("gh", "run", "view", fmt.Sprint(run.DatabaseID), "--json", "jobs").Output()
		if err != nil {
			return nil, err
		}
		var view struct {
			Jobs []struct{ StartedAt, CompletedAt time.Time }
		}
		if err = json.Unmarshal(data, &view); err != nil {
			return nil, err
		}
		seconds := 0
		for _, job := range view.Jobs {
			seconds = max(seconds, int(math.Ceil(job.CompletedAt.Sub(job.StartedAt).Seconds())))
		}
		source["pull-request"] = append(source["pull-request"], Sample{At: run.CreatedAt, Seconds: seconds})
	}
	return source, nil
}

// ciRun is one listed run of the pull request verification workflow.
type ciRun struct {
	DatabaseID int
	CreatedAt  time.Time
}

// runsInWindow keeps the runs a window can count: those created at or after its
// start, oldest first, at most its size. Each kept run costs one GitHub API
// call, so reading every listed run would exhaust the rate limit.
func runsInWindow(runs []ciRun, gate *Gate) []ciRun {
	if gate == nil {
		return nil
	}
	kept := make([]ciRun, 0, len(runs))
	for _, run := range runs {
		if gate.Window.Since != nil && run.CreatedAt.Before(*gate.Window.Since) {
			continue
		}
		kept = append(kept, run)
	}
	sort.SliceStable(kept, func(i, j int) bool { return kept[i].CreatedAt.Before(kept[j].CreatedAt) })
	if gate.Window.Size > 0 && len(kept) > gate.Window.Size {
		kept = kept[:gate.Window.Size]
	}
	return kept
}
