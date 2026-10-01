// Suite: Spec judge command
// Invariant: the advisory command preserves its transcript, environment, and JSON contracts.
// Boundary IN: public command dispatch, real temporary Spec files and Judge Log
// Boundary OUT: HTTP service (fake RoundTripper), judgment policy (internal/judge)
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/judge"
)

type specJudgeTransport struct {
	t          *testing.T
	calls      int
	failSecond bool
	host, key  string
}

func (f *specJudgeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.calls++
	if f.host == "" {
		f.t.Fatal("unexpected request")
	}
	if r.URL.Host != f.host || r.Header.Get("Authorization") != "Bearer "+f.key {
		f.t.Fatalf("wrong transport or command key: host=%s", r.URL.Host)
	}
	if f.failSecond && f.calls == 2 {
		return &http.Response{StatusCode: 503, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{}"))}, nil
	}
	var body struct {
		Questions map[string]json.RawMessage `json:"questions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		f.t.Fatal(err)
	}
	answer := `{"type":"choice","choice":"supports","probabilities":{"supports":0.93,"contradicts":0.05,"says_nothing":0.02},"confidence":0.93}`
	questions, err := judge.Load()
	if err != nil {
		f.t.Fatal(err)
	}
	id := questions.Citation.QuestionID
	if _, ok := body.Questions[questions.Goal.QuestionID]; ok {
		id = questions.Goal.QuestionID
		noul := 0.9
		if f.calls == 5 {
			noul = 0.12
		}
		answer = fmt.Sprintf(`{"type":"noul","noul":%g}`, noul)
	} else if f.calls == 1 {
		answer = `{"type":"choice","choice":"says_nothing","probabilities":{"supports":0.02,"contradicts":0.05,"says_nothing":0.93},"confidence":0.93}`
	}
	model, cost := "jev-1.13.0", ""
	if f.host == "openrouter.ai" {
		model, cost = "typesafe/jev-1.13-20260917", `,"cost":0.000035364`
	}
	response := fmt.Sprintf(`{"model":%q,"answers":{%q:%s},"usage":{"input_tokens":842,"output_tokens":20%s}}`, model, id, answer, cost)
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}, nil
}

const specJudgeEnglish = "The author reads the evidence for the work and keeps the rule in the documentation. "
const specJudgeAdvisory = "advisory citation-support docs/specs/0300-example/_techspec.md:31 ADR-0035 says_nothing at confidence 0.93: The cited decision keeps every Run Event for ninety days after the Run ends.\n"
const specJudgeNoKey = "Judge: skipped: ROUNDFIX_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY); 5 judgment(s) not asked\n"

func specJudgeFixture(t *testing.T, keyVariable string) (commandEnvironment, *specJudgeTransport) {
	t.Helper()
	repo, home := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "--initial-branch=main")
	repo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	write := func(relative, text string) {
		path := filepath.Join(repo, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		mustWrite(t, path, text)
	}
	write("docs/adr/0035-spec-root.md", "---\nstatus: accepted\n---\n# Spec Root\n"+strings.Repeat(specJudgeEnglish, 4))
	write("docs/specs/0300-example/_prd.md", "# Example\n"+specJudgeEnglish+"\n## Goals\n- The author reads evidence for the work.\n- The author keeps the cache for the work.\n")
	lines := make([]string, 80)
	lines[0], lines[1] = "# Example", specJudgeEnglish
	lines[30] = "ADR-0035 keeps every Run Event for ninety days after the Run ends."
	lines[32] = "ADR-0035 requires the author to read every rule before the work begins."
	lines[34] = "ADR-0035 preserves the evidence for every artifact that the author reads."
	lines[70], lines[71], lines[73] = "## Coverage Map", "- Goal 1 → Reader mechanism", "- Goal 2 → The widget cache"
	lines[75], lines[76], lines[78], lines[79] = "## Reader mechanism", strings.Repeat(specJudgeEnglish, 5), "## The widget cache", strings.Repeat(specJudgeEnglish, 5)
	write("docs/specs/0300-example/_techspec.md", strings.Join(lines, "\n"))
	write("docs/specs/0301-exemplo/_prd.md", "A decisão é uma regra para os autores e não está na sua documentação.\n")
	fake := &specJudgeTransport{t: t}
	env := commandEnvironment{homeDir: home, workDir: repo, dependencies: defaultCommandDependencies(), environ: []string{}}
	if keyVariable != "" {
		env.environ = []string{keyVariable + "=fixture-key"}
		if keyVariable == "ROUNDFIX_OPENROUTER_API_KEY" {
			fake.host = "openrouter.ai"
		}
		if keyVariable == "ROUNDFIX_TYPESAFE_API_KEY" {
			fake.host = "api.typesafe.ai"
		}
		fake.key = "fixture-key"
	}
	env.dependencies.judgeTransport = fake
	env.dependencies.judgeNow = func() time.Time { return time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC) }
	return env, fake
}

func specJudgeRun(t *testing.T, env commandEnvironment, args []string, wantOut, wantErr string, wantCode int) {
	t.Helper()
	var out, err bytes.Buffer
	code := runWithContext(context.Background(), append([]string{"spec", "judge"}, args...), &out, &err, env)
	if code != wantCode || out.String() != wantOut || err.String() != wantErr {
		t.Fatalf("exit=%d want=%d\nstdout=%q\nwant=%q\nstderr=%q\nwant=%q", code, wantCode, out.String(), wantOut, err.String(), wantErr)
	}
}

func TestSpecJudgeReportsAdvisoryJudgments(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	specJudgeRun(t, env, []string{"0300-example", "--stage", "techspec"}, specJudgeAdvisory+"advisory goal-mechanism docs/specs/0300-example/_techspec.md:74 Goal 2 → The widget cache: P(delivers) 0.12\nJudge: 2 advisory, 3 clear, 0 skipped; 5 call(s), 4210 input tokens, US$0.0002; month US$0.0002 of US$5.00; model jev-1.13 via openrouter\n", "", 0)
	if fake.calls != 5 {
		t.Fatalf("calls=%d", fake.calls)
	}
	log, err := os.ReadFile(filepath.Join(env.homeDir, ".roundfix/judge/2026-10.jsonl"))
	if err != nil || bytes.Count(log, []byte("\n")) != 5 || !bytes.Contains(log, []byte(`"model":"typesafe/jev-1.13-20260917"`)) {
		t.Fatalf("Judge Log: %v %s", err, log)
	}
}

func TestSpecJudgeSkipsWithoutAKey(t *testing.T) {
	env, _ := specJudgeFixture(t, "OPENROUTER_API_KEY")
	specJudgeRun(t, env, []string{"0300-example"}, specJudgeNoKey, "", 0)
}

func TestSpecJudgeReadsTheKeyFromItsOwnEnvironment(t *testing.T) {
	// Sequential: this regression intentionally poisons the process environment.
	t.Setenv("ROUNDFIX_OPENROUTER_API_KEY", "process-openrouter")
	t.Setenv("ROUNDFIX_TYPESAFE_API_KEY", "process-typesafe")
	env, _ := specJudgeFixture(t, "OPENROUTER_API_KEY")
	specJudgeRun(t, env, []string{"0300-example"}, specJudgeNoKey, "", 0)
}

func TestSpecJudgeSkipsAtTheMonthlyCeiling(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_TYPESAFE_API_KEY")
	fake.host = ""
	path := filepath.Join(env.homeDir, ".roundfix/judge/2026-10.jsonl")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, path, "{\"cost_usd\":5.0003}\n")
	specJudgeRun(t, env, []string{"0300-example"}, "Judge: skipped: monthly ceiling reached (US$5.0003 of US$5.00); 5 judgment(s) not asked\n", "", 0)
}

func TestSpecJudgeSkipsANonEnglishSpec(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	fake.host = ""
	specJudgeRun(t, env, []string{"0301-exemplo", "--stage", "prd"}, "skipped docs/specs/0301-exemplo/_prd.md: not English\nJudge: 0 advisory, 0 clear, 0 skipped; 0 call(s), 0 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.13 via openrouter\n", "", 0)
}

func TestSpecJudgeStopsWhenTheServiceFails(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_TYPESAFE_API_KEY")
	fake.failSecond = true
	specJudgeRun(t, env, []string{"0300-example", "--stage", "techspec"}, specJudgeAdvisory+"Judge: 1 advisory, 0 clear, 4 skipped; 2 call(s), 842 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.13 via typesafe; stopped: service unavailable (HTTP 503)\n", "", 0)
	if fake.calls != 2 {
		t.Fatalf("calls after stop=%d", fake.calls)
	}
}

func TestSpecJudgeRefusesAnUnknownSpec(t *testing.T) {
	env, _ := specJudgeFixture(t, "")
	specJudgeRun(t, env, []string{"0999-missing"}, "", "roundfix: spec judge failed: unknown active Spec slug \"0999-missing\"\nRun 'roundfix spec judge --help' for usage.\n", 2)
}

func TestSpecJudgeRefusesAnUnknownStageOrFormat(t *testing.T) {
	for _, tc := range []struct {
		name    string
		args    []string
		message string
	}{
		{"stage", []string{"0300-example", "--stage", "qa"}, `unsupported --stage "qa"; use prd or techspec`},
		{"format", []string{"0300-example", "--format=yaml"}, `unsupported --format "yaml"; use text or json`},
		{"flag", []string{"0300-example", "--bogus"}, `unknown flag "--bogus"`},
		{"missing slug", nil, "spec judge requires one Spec slug"},
		{"missing value", []string{"0300-example", "--stage"}, "--stage requires a value"},
		{"extra slug", []string{"0300-example", "another"}, `unexpected argument "another"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env, _ := specJudgeFixture(t, "")
			specJudgeRun(t, env, tc.args, "", "roundfix: spec judge failed: "+tc.message+"\nRun 'roundfix spec judge --help' for usage.\n", 2)
		})
	}
}

func TestSpecJudgePrintsJSON(t *testing.T) {
	env, _ := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	var out, err bytes.Buffer
	code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example", "--format", "json"}, &out, &err, env)
	if code != 0 || err.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &err)
	}
	var doc map[string]json.RawMessage
	if e := json.Unmarshal(out.Bytes(), &doc); e != nil {
		t.Fatal(e)
	}
	fields := strings.Fields("schema spec model transport skipped stopped artifacts_skipped judgments calls input_tokens cost_usd month_cost_usd month_ceiling_usd")
	if len(doc) != len(fields) {
		t.Fatalf("fields=%v", doc)
	}
	for _, field := range fields {
		if _, ok := doc[field]; !ok {
			t.Errorf("missing %s", field)
		}
	}
	if string(doc["schema"]) != `"roundfix/spec-judge/v1"` || string(doc["transport"]) != `"openrouter"` || string(doc["calls"]) != "5" {
		t.Fatalf("document=%s", &out)
	}
	var judgments []map[string]json.RawMessage
	if e := json.Unmarshal(doc["judgments"], &judgments); e != nil {
		t.Fatal(e)
	}
	if len(judgments) != 5 {
		t.Fatalf("judgments=%s", doc["judgments"])
	}
	for i, j := range judgments {
		for _, field := range strings.Fields("kind artifact line target text section_title outcome reason answer probabilities confidence noul model") {
			if _, ok := j[field]; !ok {
				t.Errorf("judgment %d missing %s", i, field)
			}
		}
	}
	if string(judgments[1]["outcome"]) != `"clear"` || string(judgments[1]["noul"]) != "null" || string(judgments[4]["section_title"]) != `"The widget cache"` {
		t.Fatalf("judgments=%s", doc["judgments"])
	}
}

func TestSpecJudgeHelp(t *testing.T) {
	env, _ := specJudgeFixture(t, "")
	specJudgeRun(t, env, []string{"--help"}, specJudgeUsage, "", 0)
	for _, text := range []string{specJudgeUsage, specUsage, usage} {
		if !strings.Contains(text, "roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]") {
			t.Errorf("missing judge synopsis")
		}
	}
	for _, phrase := range []string{"ROUNDFIX_OPENROUTER_API_KEY", "ROUNDFIX_TYPESAFE_API_KEY", "monthly ceiling", "Judge Log", "Never fails for a judgment"} {
		if !strings.Contains(specJudgeUsage, phrase) {
			t.Errorf("missing help phrase %s", phrase)
		}
	}
}

func TestSpecJudgeUsesConfiguredSpecRoot(t *testing.T) {
	env, _ := specJudgeFixture(t, "OPENROUTER_API_KEY")
	root := filepath.Join(t.TempDir(), "specs")
	if err := os.Rename(filepath.Join(env.workDir, "docs/specs"), root); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(env.workDir, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", root))
	specJudgeRun(t, env, []string{"0300-example"}, specJudgeNoKey, "", 0)
}

func TestSpecJudgeRequiresStageArtifacts(t *testing.T) {
	for _, name := range []string{"_prd.md", "_techspec.md"} {
		t.Run(name, func(t *testing.T) {
			env, _ := specJudgeFixture(t, "")
			if err := os.Remove(filepath.Join(env.workDir, "docs/specs/0300-example", name)); err != nil {
				t.Fatal(err)
			}
			var out, stderr bytes.Buffer
			code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example", "--stage=techspec"}, &out, &stderr, env)
			if code != 2 || out.Len() != 0 || !strings.Contains(stderr.String(), "spec judge failed:") {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &stderr)
			}
		})
	}
}

func TestSpecJudgePrintsIndividualSkip(t *testing.T) {
	env, fake := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	fake.host = ""
	mustWrite(t, filepath.Join(env.workDir, "docs/specs/0300-example/_prd.md"), specJudgeEnglish+"\n\nADR-0999 keeps every Run Event for ninety days after the Run ends.\n")
	specJudgeRun(t, env, []string{"0300-example", "--stage=prd"}, "skipped citation-support docs/specs/0300-example/_prd.md:3 ADR-0999: cited decision is not an accepted regular ADR\nJudge: 0 advisory, 0 clear, 1 skipped; 0 call(s), 0 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.13 via openrouter\n", "", 0)
}

func TestSpecJudgeJSONSkipHasNullTransport(t *testing.T) {
	env, _ := specJudgeFixture(t, "OPENROUTER_API_KEY")
	var out, stderr bytes.Buffer
	code := runWithContext(context.Background(), []string{"spec", "judge", "0300-example", "--format=json"}, &out, &stderr, env)
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if code != 0 || stderr.Len() != 0 || string(doc["transport"]) != "null" || string(doc["stopped"]) != "null" || string(doc["artifacts_skipped"]) != "[]" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &stderr)
	}
}
