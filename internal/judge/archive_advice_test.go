package judge

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func archiveAdviceResponse(t *testing.T, q Questions) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body struct {
			Questions map[string]json.RawMessage `json:"questions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		id := q.Archive.QuestionID
		answer := `{"type":"choice","choice":"reusable_knowledge","probabilities":{"reusable_knowledge":0.94,"repository_record":0.04,"transient_evidence":0.02},"confidence":0.94}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"model":"jev-1.13.0","answers":{"` + id + `":` + answer + `},"usage":{"input_tokens":10,"output_tokens":2}}`))}, nil
	})
}

func TestArchiveAdviceJudgesOnlyCandidateFiles(t *testing.T) {
	q, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for _, name := range []string{"references.md", "_prd.md", "qa/evidence/out.md", "binary.dat"} {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		data := []byte("reusable text")
		if name == "binary.dat" {
			data = []byte{0, 1}
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := AdviseArchive(context.Background(), q, ArchiveAdviceRequest{RepoRoot: dir, SpecDir: dir, Spec: "demo", Files: []string{"references.md", "_prd.md", "qa/evidence/out.md", "binary.dat"}, Keys: map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "key"}, HomeDir: t.TempDir(), Transport: archiveAdviceResponse(t, q)})
	if err != nil || len(got.Advice) != 2 || got.Calls != 1 || got.Advice[0].Choice == nil || *got.Advice[0].Choice != "reusable_knowledge" || got.Advice[1].Reason == nil || *got.Advice[1].Reason != "binary" {
		t.Fatalf("report=%+v err=%v", got, err)
	}
}

func TestArchiveAdviceFailsOpenWithoutAKey(t *testing.T) {
	q, _ := Load()
	got, err := AdviseArchive(context.Background(), q, ArchiveAdviceRequest{SpecDir: t.TempDir(), Files: []string{"candidate.md"}, Keys: nil, HomeDir: t.TempDir(), Transport: neverRequest(t)})
	if err != nil || got.Skipped == nil || got.Calls != 0 {
		t.Fatalf("report=%+v err=%v", got, err)
	}
}

func TestArchiveAdviceStopsAtTheCeiling(t *testing.T) {
	q, _ := Load()
	q = q.WithMonthlyCeiling(1)
	home := t.TempDir()
	logDir := filepath.Join(home, ".roundfix", "judge")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(logDir, "2026-10.jsonl"), []byte(`{"cost_usd":1}`+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := AdviseArchive(context.Background(), q, ArchiveAdviceRequest{SpecDir: t.TempDir(), Files: []string{"candidate.md"}, Keys: map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "key"}, HomeDir: home, Transport: neverRequest(t)})
	if err != nil || got.Skipped == nil || got.Calls != 0 {
		t.Fatalf("report=%+v err=%v", got, err)
	}
}
