package judge

// Suite: source grouping at the bounded reader and injected HTTP boundary.
// All filesystem writes are in t.TempDir; clocks and transports are fixed.
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func groupingBlock(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("../../docs/specs/0209-sources-that-share-a-context-share-a-spec/_techspec.md")
	if err != nil {
		b, err = os.ReadFile("../../docs/history/specs/0209-sources-that-share-a-context-share-a-spec/_techspec.md")
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.SplitN(strings.SplitN(strings.SplitN(string(b), "### The grouping question", 2)[1], "```text\n", 2)[1], "```", 2)[0]
}

func groupingFixture(t *testing.T) (Questions, Request) {
	t.Helper()
	q, req := runFixture(t)
	writeFixture(t, req.SpecDir, "_prd.md", englishContext)
	writeFixture(t, req.SpecDir, "_techspec.md", englishContext)
	writeFixture(t, req.SpecDir, "references/_index.md", "| source | type | owner | adopted date | path |\n| --- | --- | --- | --- | --- |\n| original | "+q.Grouping.AdoptedSourceTypes[0]+" | example | today | anchor.md |\n")
	writeFixture(t, req.SpecDir, "references/anchor.md", "---\nstatus: done\nprivate: FRONT_SENTINEL\n---\n"+englishContext+"The run events need a bound.\n")
	for i, name := range []string{"a-prune.md", "b-header.md"} {
		writeFixture(t, req.RepoRoot, "docs/backlog/"+name, "---\nstatus: "+q.Grouping.OpenBacklogStatuses[0]+"\n---\n"+englishContext+fmt.Sprintf("The work for item %d is <visible> & bounded.\n", i))
	}
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return groupingResponse(q, q.Transports[0].RequestModel, q.Grouping.SuggestWhenNoulAtLeast), nil
	})
	return q, req
}

func groupingResponse(q Questions, model string, noul float64) *http.Response {
	if model == q.Transports[0].RequestModel {
		model = q.Transports[1].Name + "/" + q.PinnedModel + "-20260917"
	}
	return fixtureResponse(200, fmt.Sprintf(`{"model":%q,"answers":{%q:{"type":%q,"noul":%g}},"usage":{"input_tokens":920,"cost":0.00003864}}`, model, q.Grouping.QuestionID, q.Grouping.Question.Type, noul))
}

func TestGroupingQuestionIsTheMeasuredOne(t *testing.T) {
	q, req := groupingFixture(t)
	block := groupingBlock(t)
	if !strings.Contains(string(questionFile), strings.TrimSuffix(block, "\n")) {
		t.Fatal("measured member bytes differ")
	}
	var measured struct {
		Grouping GroupingJudgment `json:"source-grouping"`
	}
	if err := json.Unmarshal([]byte("{"+block+"}"), &measured); err != nil {
		t.Fatal(err)
	}
	got := q.Grouping
	got.SourceScrub = nil
	if !reflect.DeepEqual(got, measured.Grouping) || q.Grouping.SourceScrub == nil {
		t.Fatalf("loaded grouping differs: %+v", got)
	}
	if runChecked(t, q, req).Calls != 2 {
		t.Fatal("measured question not asked")
	}
}

func TestGroupingReaderAcceptsOnlyFindingsAndBacklogEntries(t *testing.T) {
	q, req := groupingFixture(t)
	backlog := filepath.Join(req.RepoRoot, "docs/backlog")
	finding := filepath.Join(req.RepoRoot, "docs/findings")
	for _, tc := range []struct {
		dir, name, text string
		ok              bool
		reason          string
	}{
		{backlog, "open.md", "---\nstatus: " + q.Grouping.OpenBacklogStatuses[0] + "\n---\n" + englishContext, true, ""},
		{backlog, "declined.md", "---\nstatus: declined\n---\n" + englishContext, false, ""},
		{backlog, "missing.md", englishContext, false, ""},
		{finding, "pending.md", "---\nstatus: " + q.Grouping.UnresolvedFindingStatuses[0] + "\n---\n" + englishContext, true, ""},
		{finding, "partial.md", "---\nstatus: " + q.Grouping.UnresolvedFindingStatuses[1] + "\n---\n" + englishContext, true, ""},
		{finding, "done.md", "---\nstatus: done\n---\n" + englishContext, false, ""},
		{backlog, "portuguese.md", "---\nstatus: " + q.Grouping.OpenBacklogStatuses[0] + "\n---\nA decisão é uma regra para os autores e não está na sua documentação.\n", false, "not English"},
		{backlog, "big.md", strings.Repeat("x", maxSourceBytes+1), false, "not a regular file in its directory"},
		{backlog, "_index.md", englishContext, false, "not a regular file in its directory"},
		{backlog, "code.go", englishContext, false, "not a regular file in its directory"},
		{filepath.Join(req.RepoRoot, "docs/_inbox"), "note.md", englishContext, false, "only Findings and Backlog Entries are sent"},
	} {
		writeFixture(t, tc.dir, tc.name, tc.text)
		_, ok, err := readGroupingSource(tc.dir, tc.name)
		reason := ""
		if err != nil {
			reason = err.Error()
		}
		if ok != tc.ok || reason != tc.reason {
			t.Fatalf("%s: ok=%v reason=%q", tc.name, ok, reason)
		}
	}
	if err := os.Symlink(filepath.Join(backlog, "open.md"), filepath.Join(backlog, "linked.md")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"linked.md", "../backlog/open.md"} {
		if _, ok, err := readGroupingSource(backlog, name); ok || err == nil || err.Error() != "not a regular file in its directory" {
			t.Fatalf("accepted %s", name)
		}
	}
	index := "| source | type | owner | adopted date | path |\n| --- | --- | --- | --- | --- |\n| note | inbox | example | today | note.md |\n| finding | " + q.Grouping.AdoptedSourceTypes[0] + " | example | today | anchor.md |\n"
	writeFixture(t, req.SpecDir, "references/_index.md", index)
	writeFixture(t, req.SpecDir, "references/note.md", "INBOX_SENTINEL")
	anchors, skips, err := adoptedSources(q, req.SpecDir)
	if err != nil || len(anchors) != 1 || len(skips) != 1 || skips[0].Reason != "only Findings and Backlog Entries are sent" {
		t.Fatalf("anchors=%v skips=%v err=%v", anchors, skips, err)
	}
	runChecked(t, q, req)
	// A similar table cannot grant adoption, and a directory link cannot escape.
	writeFixture(t, req.SpecDir, "references/_index.md", strings.Replace(index, "adopted date", "date", 1))
	anchors, _, err = adoptedSources(q, req.SpecDir)
	if err != nil || len(anchors) != 0 {
		t.Fatal("non-adoption table accepted")
	}
	req.Transport = neverRequest(t)
	if got := runChecked(t, q, req); got.Calls != 0 {
		t.Fatal("wrong header sent sources")
	}
	outside := filepath.Join(t.TempDir(), "references")
	writeFixture(t, outside, "_index.md", index)
	writeFixture(t, outside, "anchor.md", englishContext)
	if err := os.RemoveAll(filepath.Join(req.SpecDir, "references")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(req.SpecDir, "references")); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := readGroupingSource(filepath.Join(req.SpecDir, "references"), "anchor.md"); ok || err == nil {
		t.Fatal("linked directory accepted")
	}
	if got := runChecked(t, q, req); got.Calls != 0 || len(got.ArtifactsSkipped) != 1 {
		t.Fatalf("linked directory report=%+v", got)
	}

}

func TestGroupingSourcesArePreparedAsMeasured(t *testing.T) {
	q, req := groupingFixture(t)
	body := " \nThe Spec 0123 and ADR-0456 describe <the> & contract.\n" + strings.Repeat("é", q.Grouping.SourceMaxChars) + "TAIL"
	writeFixture(t, req.SpecDir, "references/anchor.md", "---\nstatus: done\nsecret: FRONT_SENTINEL\n---\n"+body)
	src, ok, err := readGroupingSource(filepath.Join(req.SpecDir, "references"), "anchor.md")
	if err != nil || !ok {
		t.Fatalf("reader: %v %v", ok, err)
	}
	want := cut(strings.ReplaceAll(strings.ReplaceAll(body, "0123", ""), "ADR-0456", ""), q.Grouping.SourceMaxChars)
	if got := prepareSource(q, src); got != want || len([]rune(got)) != q.Grouping.SourceMaxChars {
		t.Fatalf("preparation differs: %q", got)
	}
	// An empty leading block is also removed without trimming its body.
	writeFixture(t, req.SpecDir, "references/anchor.md", "---\n---\n"+englishContext)
	src, ok, err = readGroupingSource(filepath.Join(req.SpecDir, "references"), "anchor.md")
	if err != nil || !ok || prepareSource(q, src) != englishContext {
		t.Fatal("empty front matter not removed")
	}
	runChecked(t, q, req)
}

func TestGroupingPairsEveryAdoptedSourceWithEveryOpenSource(t *testing.T) {
	q, req := groupingFixture(t)
	writeFixture(t, req.RepoRoot, "docs/findings/c.md", "---\nstatus: "+q.Grouping.UnresolvedFindingStatuses[0]+"\n---\n"+englishContext+"The finding needs the repair.\n")
	index, err := os.ReadFile(filepath.Join(req.SpecDir, "references/_index.md"))
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, req.SpecDir, "references/_index.md", string(index)+"| second | "+q.Grouping.AdoptedSourceTypes[1]+" | example | today | second.md |\n")
	writeFixture(t, req.SpecDir, "references/second.md", englishContext+"The second source describes the repair.\n")
	plan, err := PlanSpec(q, req.RepoRoot, req.SpecDir, "prd")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Pending) != 6 {
		t.Fatalf("pairs=%+v", plan.Pending)
	}
	for i, p := range plan.Pending {
		anchor := []string{"anchor.md", "second.md"}[i/3]
		target := []string{"docs/backlog/a-prune.md", "docs/backlog/b-header.md", "docs/findings/c.md"}[i%3]
		if p.Kind != "source-grouping" || p.Artifact != filepath.Join(req.SpecDir, "references", anchor) || p.Target != filepath.Join(req.RepoRoot, target) || p.Line != 0 {
			t.Fatalf("pair %d: %+v", i, p)
		}
	}
	if runChecked(t, q, req).Calls != 6 {
		t.Fatal("pairs not asked")
	}
	// Same prepared state is asked once, including duplicate adoption rows.
	writeFixture(t, req.SpecDir, "references/second.md", string(mustReadGrouping(t, filepath.Join(req.SpecDir, "references/anchor.md"))))
	if runChecked(t, q, req).Calls != 3 {
		t.Fatal("identical states not deduplicated")
	}
}

func mustReadGrouping(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestGroupingIsPlannedAtEveryStage(t *testing.T) {
	for _, stage := range []Stage{"prd", "techspec", ""} {
		t.Run(string(stage), func(t *testing.T) {
			q, req := groupingFixture(t)
			req.Stage = stage
			if got := runChecked(t, q, req); got.Calls != 2 || len(got.Judgments) != 2 {
				t.Fatalf("stage %q: %+v", stage, got)
			}
			if err := os.Remove(filepath.Join(req.SpecDir, "references/_index.md")); err != nil {
				t.Fatal(err)
			}
			req.Transport = neverRequest(t)
			if got := runChecked(t, q, req); got.Calls != 0 || len(got.Judgments) != 0 {
				t.Fatalf("no index: %+v", got)
			}
		})
	}
}

func TestGroupingSuggestsAtTheMeasuredThreshold(t *testing.T) {
	for _, tc := range []struct {
		noul    float64
		outcome string
	}{{0.30, "suggested"}, {0.29, "clear"}} {
		t.Run(tc.outcome, func(t *testing.T) {
			q, req := groupingFixture(t)
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
				return groupingResponse(q, q.Transports[0].RequestModel, tc.noul), nil
			})
			got := runChecked(t, q, req)
			if got.Calls != 2 {
				t.Fatalf("calls=%d", got.Calls)
			}
			for _, j := range got.Judgments {
				if j.Outcome != tc.outcome || j.Noul == nil || *j.Noul != tc.noul {
					t.Fatalf("judgment=%+v", j)
				}
			}
		})
	}
}

func TestGroupingNeverComparesAnotherModelsAnswer(t *testing.T) {
	q, req := groupingFixture(t)
	req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { return groupingResponse(q, "jev-1.14.0", 1), nil })
	got := runChecked(t, q, req)
	for _, j := range got.Judgments {
		if j.Outcome != "skipped" || j.Reason == nil || *j.Reason != "answered by jev-1.14.0, thresholds belong to "+q.PinnedModel {
			t.Fatalf("unpinned=%+v", j)
		}
	}
	rows := strings.Split(strings.TrimSpace(string(mustReadGrouping(t, newJudgeLog(req.HomeDir, req.Now()).path))), "\n")
	if len(rows) != 2 {
		t.Fatal("unpinned calls not logged")
	}
	for _, row := range rows {
		var j logLine
		if err := json.Unmarshal([]byte(row), &j); err != nil {
			t.Fatal(err)
		}
		if j.Outcome != "skipped" || j.Model != "jev-1.14.0" {
			t.Fatalf("log=%+v", j)
		}
	}
}

func TestGroupingRequestsCarryOnlySourceText(t *testing.T) {
	for _, transport := range loadQuestions(t).Transports {
		t.Run(transport.Name, func(t *testing.T) {
			q, req := groupingFixture(t)
			req.Keys = map[string]string{transport.KeyVariable: "BOUNDARY_KEY"}
			for path, text := range map[string]string{
				"docs/_inbox/note.md":        "INBOX_SENTINEL",
				"docs/backlog/declined.md":   "---\nstatus: declined\n---\nDECLINED_SENTINEL\n" + englishContext,
				"docs/findings/done.md":      "---\nstatus: done\n---\nDONE_SENTINEL\n" + englishContext,
				"source.go":                  "GO_SENTINEL",
				"docs/backlog/portuguese.md": "---\nstatus: " + q.Grouping.OpenBacklogStatuses[0] + "\n---\nA decisão é uma regra para os autores e não está na sua documentação. PORTUGUESE_SENTINEL\n",
				"outside.md":                 "LINK_SENTINEL\n" + englishContext,
			} {
				writeFixture(t, req.RepoRoot, path, text)
			}
			if err := os.Symlink(filepath.Join(req.RepoRoot, "outside.md"), filepath.Join(req.RepoRoot, "docs/backlog/link.md")); err != nil {
				t.Fatal(err)
			}
			anchors, _, err := adoptedSources(q, req.SpecDir)
			if err != nil {
				t.Fatal(err)
			}
			candidates, _, err := openSources(q, req.RepoRoot)
			if err != nil {
				t.Fatal(err)
			}
			if len(anchors) != 1 || len(candidates) != 2 {
				t.Fatalf("accepted %d/%d", len(anchors), len(candidates))
			}
			calls := 0
			req.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				for _, forbidden := range []string{"INBOX_SENTINEL", "DECLINED_SENTINEL", "DONE_SENTINEL", "GO_SENTINEL", "LINK_SENTINEL", "PORTUGUESE_SENTINEL", "FRONT_SENTINEL", "BOUNDARY_KEY"} {
					if strings.Contains(string(b), forbidden) {
						t.Fatalf("request contains %s", forbidden)
					}
				}
				if r.URL.String() != transport.Endpoint || r.Header.Get("Authorization") != "Bearer BOUNDARY_KEY" {
					t.Fatal("wrong transport")
				}
				var payload struct {
					State     map[string]string   `json:"state"`
					Questions map[string]Question `json:"questions"`
				}
				if err := json.Unmarshal(b, &payload); err != nil {
					t.Fatal(err)
				}
				if len(payload.State) != 2 || payload.State["first"] != prepareSource(q, anchors[0]) || payload.State["second"] != prepareSource(q, candidates[calls]) {
					t.Fatalf("state outside prepared sources: %+v", payload.State)
				}
				if len(payload.Questions) != 1 || !reflect.DeepEqual(payload.Questions[q.Grouping.QuestionID], q.Grouping.Question) {
					t.Fatal("wrong grouping question")
				}
				if strings.Contains(string(b), `\u003c`) || strings.Contains(string(b), `\u0026`) {
					t.Fatal("HTML escaped")
				}
				calls++
				return groupingResponse(q, transport.RequestModel, q.Grouping.SuggestWhenNoulAtLeast), nil
			})
			got := runChecked(t, q, req)
			if calls != 2 || got.Calls != 2 || len(got.ArtifactsSkipped) != 2 {
				t.Fatalf("boundary report=%+v calls=%d", got, calls)
			}
			for _, skip := range got.ArtifactsSkipped {
				if !strings.HasPrefix(skip.Artifact, "docs/backlog/") {
					t.Fatalf("wrong skip path: %+v", skip)
				}
			}
		})
	}
}

func TestGroupingJudgeLogLine(t *testing.T) {
	q, req := groupingFixture(t)
	got := runChecked(t, q, req)
	rows := strings.Split(strings.TrimSpace(string(mustReadGrouping(t, newJudgeLog(req.HomeDir, req.Now()).path))), "\n")
	if len(rows) != 2 {
		t.Fatalf("rows=%d", len(rows))
	}
	for i, row := range rows {
		var fields map[string]any
		if err := json.Unmarshal([]byte(row), &fields); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"line", "answer", "probabilities", "confidence"} {
			v, ok := fields[name]
			if !ok || v != nil {
				t.Fatalf("%s=%v", name, v)
			}
		}
		if fields["schema"] != "roundfix/judge-log/v1" || fields["judgment"] != "source-grouping" || fields["artifact"] != got.Judgments[i].Artifact || fields["target"] != got.Judgments[i].Target || fields["question_id"] != q.Grouping.QuestionID || fields["outcome"] != "suggested" || fields["noul"] != q.Grouping.SuggestWhenNoulAtLeast {
			t.Fatalf("log=%v", fields)
		}
		if fields["time"] != req.Now().Format("2006-01-02T15:04:05Z07:00") {
			t.Fatal("clock not fixed")
		}
	}
	b, err := json.Marshal(got.Judgments[0])
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(b, &fields); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"line", "text", "section_title", "answer", "probabilities", "confidence", "reason"} {
		if v, ok := fields[name]; !ok || v != nil {
			t.Fatalf("report %s=%v", name, v)
		}
	}
	cost, err := newJudgeLog(req.HomeDir, req.Now()).monthCost()
	if err != nil || cost != got.CostUSD {
		t.Fatalf("cost=%g err=%v", cost, err)
	}
}

func TestGroupingUsesExistingFailureAndSpendRules(t *testing.T) {
	for _, tc := range []struct {
		name, body    string
		status, calls int
		stop          bool
	}{
		{"malformed", "invalid", 200, 2, false},
		{"missing answer", `{ "model":"jev-1.13.0","answers":{} }`, 200, 2, false},
		{"refused", "", 422, 2, false},
		{"key refused", "", 401, 1, true},
		{"retry exhausted", "", 529, 3, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q, req := groupingFixture(t)
			calls := 0
			req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return fixtureResponse(tc.status, tc.body), nil })
			got := runChecked(t, q, req)
			if calls != tc.calls || (got.Stopped != nil) != tc.stop {
				t.Fatalf("report=%+v calls=%d", got, calls)
			}
			for _, j := range got.Judgments {
				if j.Outcome != "skipped" {
					t.Fatalf("judgment=%+v", j)
				}
			}
			rows := strings.Split(strings.TrimSpace(string(mustReadGrouping(t, newJudgeLog(req.HomeDir, req.Now()).path))), "\n")
			if len(rows) != tc.calls {
				t.Fatal("every attempt must be logged")
			}
		})
	}
	t.Run("ceiling before retry", func(t *testing.T) {
		q, req := groupingFixture(t)
		calls := 0
		req.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			writeFixture(t, req.HomeDir, ".roundfix/judge/2026-10.jsonl", fmt.Sprintf(`{"cost_usd":%g}`, q.MonthlyCeilingUSD)+"\n")
			return fixtureResponse(429, ""), nil
		})
		got := runChecked(t, q, req)
		if calls != 1 || got.Stopped == nil || !strings.HasPrefix(*got.Stopped, "monthly ceiling reached") {
			t.Fatalf("report=%+v calls=%d", got, calls)
		}
	})
	t.Run("grouping follows existing judgments", func(t *testing.T) {
		q, req := groupingFixture(t)
		_, original := runFixture(t)
		for _, name := range []string{"_prd.md", "_techspec.md"} {
			writeFixture(t, req.SpecDir, name, string(mustReadGrouping(t, filepath.Join(original.SpecDir, name))))
		}
		plan, err := PlanSpec(q, req.RepoRoot, req.SpecDir, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(plan.Pending) != 4 || plan.Pending[0].Kind != "citation-support" || plan.Pending[1].Kind != "goal-mechanism" || plan.Pending[2].Kind != "source-grouping" || plan.Pending[3].Kind != "source-grouping" {
			t.Fatalf("order=%+v", plan.Pending)
		}
	})
}
