package judge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// The default HTTP transport reaches only this local server, including the
// empty-marker control that proves the fixture can actually send a request.
func runMarkerServer(t *testing.T, q *Questions) *atomic.Int32 {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var body struct {
			Questions map[string]json.RawMessage `json:"questions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		answers := make(map[string]any)
		for id := range body.Questions {
			switch id {
			case q.Goal.QuestionID, q.Grouping.QuestionID:
				answers[id] = map[string]any{"type": "noul", "noul": 0.9}
			case q.Archive.QuestionID:
				answers[id] = map[string]any{"type": "choice", "choice": "reusable_knowledge", "probabilities": map[string]float64{"reusable_knowledge": 0.94, "repository_record": 0.04, "transient_evidence": 0.02}, "confidence": 0.94}
			default:
				answers[id] = map[string]any{"type": "choice", "choice": "supports", "probabilities": map[string]float64{"supports": 0.94, "contradicts": 0.04, "says_nothing": 0.02}, "confidence": 0.94}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"model": "jev-1.13.0", "answers": answers, "usage": map[string]int{"input_tokens": 10, "output_tokens": 2}}); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)
	for i := range q.Transports {
		q.Transports[i].Endpoint = server.URL
	}
	return &calls
}

func assertRunMarkerNoLog(t *testing.T, home string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(home, ".roundfix", "judge")); !os.IsNotExist(err) {
		t.Fatalf("Judge Log directory exists or could not be checked: %v", err)
	}
}

func TestJudgeSendsNothingInsideARun(t *testing.T) {
	t.Parallel()
	for _, runID := range []string{"run-marker-fixture", ""} {
		t.Run(fmt.Sprintf("RunID=%q", runID), func(t *testing.T) {
			t.Parallel()
			q, req := runFixture(t)
			calls := runMarkerServer(t, &q)
			req.RunID, req.Transport = runID, nil
			if runID != "" {
				req.Now = func() time.Time { t.Fatal("Run Marker reached Judge Log clock"); return time.Time{} }
			}
			got := runChecked(t, q, req)
			if runID == "" {
				if got.Skipped != nil || got.Calls != 2 || calls.Load() != 2 {
					t.Fatalf("outside Run: report=%+v calls=%d", got, calls.Load())
				}
				if len(readLog(t, req)) != 2 {
					t.Fatal("outside Run lost Judge Log")
				}
				return
			}
			want := "inside Run run-marker-fixture, named by ROUNDFIX_RUN_ID"
			if got.Skipped == nil || *got.Skipped != want || got.Calls != 0 || calls.Load() != 0 || got.Transport != nil || got.KeyVariable != nil || got.Stopped != nil {
				t.Fatalf("report=%+v calls=%d", got, calls.Load())
			}
			if len(got.Judgments) != 2 {
				t.Fatalf("judgments=%+v", got.Judgments)
			}
			for _, judgment := range got.Judgments {
				if judgment.Outcome != "skipped" || judgment.Reason == nil || *judgment.Reason != want {
					t.Fatalf("judgment=%+v", judgment)
				}
			}
			assertRunMarkerNoLog(t, req.HomeDir)
		})
	}
}

func TestArchiveAdviceSendsNothingInsideARun(t *testing.T) {
	t.Parallel()
	for _, runID := range []string{"run-marker-fixture", ""} {
		t.Run(fmt.Sprintf("RunID=%q", runID), func(t *testing.T) {
			t.Parallel()
			q, fixture := runFixture(t)
			writeFixture(t, fixture.SpecDir, "references/lesson.md", "A reusable lesson.")
			calls := runMarkerServer(t, &q)
			req := ArchiveAdviceRequest{RunID: runID, RepoRoot: fixture.RepoRoot, SpecDir: fixture.SpecDir, Spec: fixture.Spec, HomeDir: fixture.HomeDir, Keys: fixture.Keys, Now: fixture.Now, Files: []string{"_prd.md", "qa/evidence/out.txt", "references/lesson.md"}}
			if runID != "" {
				req.Now = func() time.Time { t.Fatal("Run Marker reached Judge Log clock"); return time.Time{} }
			}
			got, err := AdviseArchive(t.Context(), q, req)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Advice) != 1 {
				t.Fatalf("advice=%+v", got.Advice)
			}
			if runID == "" {
				if got.Skipped != nil || got.Calls != 1 || calls.Load() != 1 || got.Advice[0].Choice == nil || *got.Advice[0].Choice != "reusable_knowledge" {
					t.Fatalf("outside Run: report=%+v calls=%d", got, calls.Load())
				}
				if len(readLog(t, fixture)) != 1 {
					t.Fatal("outside Run lost Judge Log")
				}
				return
			}
			want := "inside Run run-marker-fixture, named by ROUNDFIX_RUN_ID"
			if got.Skipped == nil || *got.Skipped != want || got.Calls != 0 || calls.Load() != 0 || got.Stopped != nil || got.Advice[0].Outcome != "skipped" || got.Advice[0].Reason == nil || *got.Advice[0].Reason != want {
				t.Fatalf("report=%+v calls=%d", got, calls.Load())
			}
			assertRunMarkerNoLog(t, req.HomeDir)
		})
	}
}
