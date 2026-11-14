// Suite: Idea-stage command, exact transcripts, and the one-item data boundary.
// Boundary IN: public command dispatch, temporary Git repository and request files.
// Boundary OUT: judge HTTP is replaced by a per-test transport.
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
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/gittest"
	"roundfix/internal/judge"
)

const ideaFeatureRef = "docs/backlog/2026-11-02-teams-can-see-each-others-run-results.md"
const ideaFixRef = "docs/backlog/2026-11-03-status-prints-the-header-twice.md"
const ideaNoKey = "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set (nor ROUNDFIX_OPENROUTER_API_KEY, nor ROUNDFIX_TYPESAFE_API_KEY)"
const ideaTitle = "Teams can see each other's Run results"

type ideaCommandTransport struct {
	t       *testing.T
	calls   int
	values  [3]float64
	inspect func([]byte)
	status  int
	model   string
}

func (fake *ideaCommandTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	fake.t.Helper()
	fake.calls++
	if r.URL.Host != "api.typesafe.ai" || r.Header.Get("Authorization") != "Bearer fake-judge-key" {
		fake.t.Fatal("wrong transport or environment key")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fake.t.Fatal(err)
	}
	if fake.inspect != nil {
		fake.inspect(body)
	}
	model := fake.model
	if model == "" {
		model = "jev-1.13.0"
	}
	response := fmt.Sprintf(`{"model":%q,"answers":{"shape_open":{"type":"noul","noul":%g},"unverified_assumption":{"type":"noul","noul":%g},"market_question":{"type":"noul","noul":%g}},"usage":{"input_tokens":512,"cost":0.000021504}}`, model, fake.values[0], fake.values[1], fake.values[2])
	status := fake.status
	if status == 0 {
		status = 200
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}, nil
}
func ideaCommandFixture(t *testing.T, withKey bool) (commandEnvironment, *ideaCommandTransport) {
	t.Helper()
	repo, home := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "--initial-branch=main")
	var err error
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	followUpWrite(t, repo, "README.md", "# Widget\n")
	followUpWrite(t, repo, ideaFeatureRef, "---\ntype: feat\nstatus: open\ncreated: 2026-11-02\n---\n# "+ideaTitle+"\n\nMaybe a web page, maybe a chat bot.\n")
	followUpWrite(t, repo, ideaFixRef, "---\ntype: fix\nstatus: open\ncreated: 2026-11-03\n---\n# Status prints the header twice\n")
	if err := os.MkdirAll(filepath.Join(repo, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, repo, "add", ".")
	gitImplement(t, repo, "commit", "-m", "test: seed Idea-stage input")
	gitImplement(t, repo, "remote", "add", "origin", "git@github.com:owner/widget.git")
	fake := &ideaCommandTransport{t: t, values: [3]float64{.98, .84, .17}}
	env := commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
	if withKey {
		env.environ = []string{"ROUNDFIX_TYPESAFE_API_KEY=" + strings.Join([]string{"fake", "judge", "key"}, "-")}
	}
	env.dependencies.judgeTransport = fake
	env.dependencies.judgeNow = func() time.Time { return time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC) }
	return env, fake
}
func ideaCommandError(message string) string {
	return "roundfix: spec judge failed: " + message + "\nRun 'roundfix spec judge --help' for usage.\n"
}

func TestSpecJudgeIdeaStageRequiresAnItem(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, false)
	specJudgeRun(t, env, []string{"--stage", "idea"}, "", ideaCommandError("--stage idea requires --item, the reference of the item to judge, or --request, a file that holds the request"), 2)
	if fake.calls != 0 {
		t.Fatal("request on usage error")
	}
}
func TestSpecJudgeIdeaStageRefusesAnUnknownReference(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	ref := "docs/backlog/2026-11-09-missing.md"
	specJudgeRun(t, env, []string{"--stage", "idea", "--item", ref}, "", ideaCommandError(fmt.Sprintf("no pending item %q; name a reference that 'roundfix follow-up' lists", ref)), 2)
	if fake.calls != 0 {
		t.Fatal("request for unknown reference")
	}
}
func TestSpecJudgeIdeaStageDoesNotJudgeAFix(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	specJudgeRun(t, env, []string{"--stage", "idea", "--item", ideaFixRef}, "Idea stage for "+ideaFixRef+": not judged: the item is typed fix; only a feature is judged for the Idea stage\n", "", 0)
	if fake.calls != 0 {
		t.Fatal("fix was sent")
	}
}
func TestSpecJudgeIdeaStageListsTheCriteriaWhenTheJudgeDoesNotAnswer(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, false)
	specJudgeRun(t, env, []string{"--stage", "idea", "--item", ideaFeatureRef}, ideaUnavailableTranscript(ideaFeatureRef, ideaNoKey), "", 0)
	var out, stderr bytes.Buffer
	code := runWithContext(context.Background(), []string{"spec", "judge", "--stage=idea", "--item=" + ideaFeatureRef, "--format=json"}, &out, &stderr, env)
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"reason", "verdict", "judged_by", "model", "transport", "key_variable"} {
		want := "null"
		if key == "reason" {
			data, _ := json.Marshal(ideaNoKey)
			want = string(data)
		}
		if string(doc[key]) != want {
			t.Fatalf("%s=%s want=%s", key, doc[key], want)
		}
	}
	var criteria []map[string]json.RawMessage
	if err := json.Unmarshal(doc["criteria"], &criteria); err != nil {
		t.Fatal(err)
	}
	if len(criteria) != 3 {
		t.Fatalf("criteria=%s", doc["criteria"])
	}
	for _, criterion := range criteria {
		if string(criterion["probability"]) != "null" || string(criterion["answer"]) != "null" {
			t.Fatalf("criterion=%v", criterion)
		}
	}
	if code != 0 || stderr.Len() != 0 || fake.calls != 0 {
		t.Fatalf("exit=%d stderr=%s calls=%d", code, &stderr, fake.calls)
	}
}
func TestSpecJudgeIdeaStageListsTheCriteriaForARequest(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, false)
	path := filepath.Join(t.TempDir(), "request.txt")
	text := "Teams should see each other's Run results\nMaybe a web page, maybe a chat bot.\n"
	mustWrite(t, path, text)
	specJudgeRun(t, env, []string{"--stage", "idea", "--request", path}, ideaUnavailableTranscript("the request in the session", ideaNoKey), "", 0)
	after, err := os.ReadFile(path)
	if err != nil || string(after) != text || fake.calls != 0 {
		t.Fatalf("request file changed: %v calls=%d", err, fake.calls)
	}
}

func TestSpecJudgeIdeaStagePrintsTheVerdictAndTheRecord(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		values  [3]float64
		answers [3]string
		verdict string
	}{
		{[3]float64{.98, .84, .17}, [3]string{"yes", "yes", "no"}, "indicated"},
		{[3]float64{.12, .17, .07}, [3]string{"no", "no", "no"}, "not indicated"},
		{[3]float64{.42, .20, .05}, [3]string{"unsure", "no", "no"}, "undecided"},
	} {
		t.Run(tc.verdict, func(t *testing.T) {
			// Each format gets a fresh home so its whole report has the same month cost.
			for _, format := range []string{"text", "json"} {
				env, fake := ideaCommandFixture(t, true)
				fake.values = tc.values
				want := ideaAnsweredTranscript(ideaFeatureRef, tc.verdict, tc.answers, tc.values)
				if format == "json" {
					want = ideaAnsweredJSON(ideaFeatureRef, tc.verdict, tc.answers, tc.values)
				}
				specJudgeRun(t, env, []string{"--stage=idea", "--item=" + ideaFeatureRef, "--format=" + format}, want, "", 0)
				rows, err := judge.ReadMonth(context.Background(), env.homeDir, env.dependencies.judgeNow())
				if err != nil {
					t.Fatal(err)
				}
				if fake.calls != 1 || len(rows) != 1 || rows[0].Judgment != "idea-stage" {
					t.Fatalf("calls=%d rows=%+v", fake.calls, rows)
				}
			}
		})
	}
}

func TestIdeaStageItemComesOnlyFromTheInventory(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	loaded, err := loadCommandConfig(env, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	items, _, err := followUpInventory(loaded, "widget", "", env.workDir)
	if err != nil {
		t.Fatal(err)
	}
	item, err := ideaStageItem(items, ideaFeatureRef)
	if err != nil || item.Ref != ideaFeatureRef || item.Title != ideaTitle {
		t.Fatalf("item=%+v err=%v", item, err)
	}
	// The file exists, but is neither an open backlog item nor another inventory source.
	// Its absolute path reaches it from any working directory, so a lookup that
	// fell back to reading the reference from disk would find it.
	followUpWrite(t, env.workDir, "arbitrary.md", "# Never send this file\n")
	for _, ref := range []string{"arbitrary.md", filepath.Join(env.workDir, "arbitrary.md"), "./" + ideaFeatureRef} {
		want := fmt.Sprintf("no pending item %q; name a reference that 'roundfix follow-up' lists", ref)
		if _, err := ideaStageItem(items, ref); err == nil || err.Error() != want {
			t.Fatalf("lookup %q: %v", ref, err)
		}
		specJudgeRun(t, env, []string{"--stage=idea", "--item", ref}, "", ideaCommandError(want), 2)
	}
	followUpWrite(t, env.workDir, "docs/specs/0300-active/_prd.md", "---\nstatus: active\n---\n# Active\n")
	// The inventory discovers active PRDs and reads their Task Graphs.
	followUpWrite(t, env.workDir, "docs/specs/0300-active/_tasks.md", "# Tasks\n")
	items, _, err = followUpInventory(loaded, "widget", "", env.workDir)
	if err != nil {
		t.Fatal(err)
	}
	ref := "docs/specs/0300-active"
	want := fmt.Sprintf("%q is an active Spec; the Idea stage is judged before a PRD exists", ref)
	if _, err := ideaStageItem(items, ref); err == nil || err.Error() != want {
		t.Fatalf("active Spec lookup: %v inventory=%+v", err, items)
	}
	specJudgeRun(t, env, []string{"--stage=idea", "--item", ref}, "", ideaCommandError(want), 2)
	if fake.calls != 0 {
		t.Fatal("unlisted file or Spec was sent")
	}
}

func TestSpecJudgeIdeaStageReadsAMessageEntry(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	followUpWrite(t, env.workDir, "messages.json", `{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"one","from":"other-project","created_at":"2026-11-14","subject":"Share results","text":"Perhaps a page, perhaps a bot.","type":"feat"}]}`)
	fake.inspect = func(data []byte) {
		var body struct {
			State struct {
				Project string            `json:"project"`
				Item    map[string]string `json:"item"`
			} `json:"state"`
		}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"kind": "message", "type": "feat", "origin": "other-project", "title": "Share results", "text": "Perhaps a page, perhaps a bot."}
		if body.State.Project != "widget" || !reflect.DeepEqual(body.State.Item, want) {
			t.Fatalf("state=%+v", body.State)
		}
	}
	specJudgeRun(t, env, []string{"--stage=idea", "--item=message:one", "--messages=messages.json"}, ideaAnsweredTranscript("message:one", "indicated", [3]string{"yes", "yes", "no"}, fake.values), "", 0)
	if fake.calls != 1 {
		t.Fatalf("calls=%d", fake.calls)
	}
	for _, tc := range []struct {
		args    []string
		message string
	}{
		{[]string{"--stage=idea", "--item=x", "--request=x"}, "--item and --request cannot be used together"},
		{[]string{"slug", "--stage=idea", "--item=x"}, "--stage idea takes no Spec slug; name the item with --item or the request with --request"},
		{[]string{"slug", "--item=x"}, "--item applies only to --stage idea"},
		{[]string{"slug", "--request=x"}, "--request applies only to --stage idea"},
		{[]string{"slug", "--messages=x"}, "--messages applies only to --stage idea with --item"},
		{[]string{"--stage=idea", "--request=x", "--messages=x"}, "--messages applies only to --stage idea with --item"},
		{[]string{"--stage=idea", "--messages=x"}, "--messages applies only to --stage idea with --item"},
		{[]string{"--stage=qa"}, `unsupported --stage "qa"; use idea, prd, techspec or tasks`},
		{[]string{"--stage=idea", "--item"}, "--item requires a value"},
		{[]string{"--stage=idea", "--request"}, "--request requires a value"},
		{[]string{"--stage=idea", "--messages"}, "--messages requires a value"},
		{[]string{"--stage=idea", "--item=x", "--format=yaml"}, `unsupported --format "yaml"; use text or json`},
	} {
		specJudgeRun(t, env, tc.args, "", ideaCommandError(tc.message), 2)
	}
	for _, name := range []string{"missing.json", "invalid.json"} {
		path := filepath.Join(env.workDir, name)
		if name == "invalid.json" {
			mustWrite(t, path, "{}")
		}
		_, _, err := followUpMessages(path)
		if err == nil {
			t.Fatal("invalid messages fixture")
		}
		specJudgeRun(t, env, []string{"--stage=idea", "--item=message:one", "--messages", name}, "", ideaCommandError(err.Error()), 2)
	}
	if fake.calls != 1 {
		t.Fatalf("usage errors sent requests: %d", fake.calls)
	}
}

func TestSpecJudgeIdeaStageDoesNotJudgeAWaitingItem(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	followUpWrite(t, env.workDir, ideaFeatureRef, "---\ntype: feat\nstatus: open\nwaiting_on: the user's answers\nsince: 2026-11-13\n---\n# "+ideaTitle+"\n")
	specJudgeRun(t, env, []string{"--stage=idea", "--item", ideaFeatureRef}, "Idea stage for "+ideaFeatureRef+": not judged: the item waits on the user's answers since 2026-11-13\n", "", 0)
	if fake.calls != 0 {
		t.Fatal("Waiting Item sent")
	}
}

func TestSpecJudgeIdeaStageSendsOnlyTheStartOfARequest(t *testing.T) {
	t.Parallel()
	env, fake := ideaCommandFixture(t, true)
	title := strings.Repeat("ç", 230)
	text := "  \n  ### " + title + "\n" + strings.Repeat("é", 4100) + "AFTER_THE_BOUNDARY\n"
	path := filepath.Join(t.TempDir(), "private-session-request.txt")
	mustWrite(t, path, text)
	fake.inspect = func(data []byte) {
		if strings.Contains(string(data), "AFTER_THE_BOUNDARY") || strings.Contains(string(data), path) || strings.Contains(string(data), filepath.Base(path)) {
			t.Fatal("request sent extra content or path")
		}
		var body struct {
			State struct {
				Project string            `json:"project"`
				Item    map[string]string `json:"item"`
			} `json:"state"`
		}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		excerpt := string([]rune(strings.TrimSpace(text))[:4000])
		want := map[string]string{"kind": "request", "type": "request", "origin": "widget", "title": strings.Repeat("ç", 200), "text": excerpt}
		if body.State.Project != "widget" || !reflect.DeepEqual(body.State.Item, want) || utf8.RuneCountInString(body.State.Item["text"]) != 4000 {
			t.Fatalf("bad request boundary: project=%s title runes=%d text runes=%d", body.State.Project, utf8.RuneCountInString(body.State.Item["title"]), utf8.RuneCountInString(body.State.Item["text"]))
		}
	}
	specJudgeRun(t, env, []string{"--stage=idea", "--request", path}, ideaAnsweredTranscript("the request in the session", "indicated", [3]string{"yes", "yes", "no"}, fake.values), "", 0)
	after, err := os.ReadFile(path)
	if err != nil || string(after) != text {
		t.Fatal("command mutated request file")
	}
	// Relative paths use the command's working directory.
	followUpWrite(t, env.workDir, "request.txt", "# A short request\nDetails.")
	item, err := ideaStageRequestItem(filepath.Join(env.workDir, "request.txt"), "widget")
	if err != nil || item.Title != "A short request" {
		t.Fatalf("item=%+v err=%v", item, err)
	}
	fake.inspect = nil
	specJudgeRun(t, env, []string{"--stage=idea", "--request=request.txt"}, ideaAnsweredTranscriptWithMonth("the request in the session", "indicated", [3]string{"yes", "yes", "no"}, fake.values, "0.0000"), "", 0)
	for _, name := range []string{"missing", "large", "blank"} {
		bad := filepath.Join(t.TempDir(), name)
		message := ""
		switch name {
		case "missing":
			_, err := os.Open(bad)
			message = err.Error()
		case "large":
			mustWrite(t, bad, strings.Repeat("x", (1<<20)+1))
			message = "exceeds 1 MiB"
		case "blank":
			mustWrite(t, bad, " \n\t\r\n")
			message = "holds no text"
		}
		specJudgeRun(t, env, []string{"--stage=idea", "--request", bad}, "", ideaCommandError(fmt.Sprintf("request file %q: %s", bad, message)), 2)
	}
	if fake.calls != 2 {
		t.Fatalf("calls=%d", fake.calls)
	}
}

func TestSpecJudgeIdeaHelpPreservesExistingLines(t *testing.T) {
	t.Parallel()
	const synopsis = "roundfix spec judge --stage idea (--item <reference> [--messages <path>] | --request <path>) [--format <text|json>]"
	for _, help := range []string{specJudgeUsage, specUsage, usage} {
		if strings.Count(help, synopsis) != 1 {
			t.Fatalf("synopsis missing or repeated: %s", help)
		}
	}
	for _, phrase := range []string{"--item", "--messages", "--request", "Judge idea, prd, techspec or tasks", "or unknown item", "Of a request file only the first line and the first 4,000 characters are sent. It never gates."} {
		if !strings.Contains(strings.Join(strings.Fields(specJudgeUsage), " "), phrase) {
			t.Fatalf("help missing %q", phrase)
		}
	}
}

func ideaUnavailableTranscript(ref, reason string) string {
	return fmt.Sprintf("Idea stage for %s: not judged: %s\nAnswer the three criteria yourself, each as yes, no or unsure, and say that the judge did not answer:\n  shape_open: the right shape of the solution is still open\n    yes: At least two materially different designs could satisfy the item, and the item does not settle which one to build.\n    no: The item names one concrete design, or only one reasonable design exists.\n  unverified_assumption: it rests on an assumption about its users, the business or the market that it gives no evidence for\n    yes: A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong.\n    no: The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made.\n  market_question: deciding to build it depends on a market or competitive question\n    yes: The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it.\n    no: The decision depends only on the product's own behavior and its current users.\nRule: indicated when any criterion is yes, not indicated when every criterion is no, undecided otherwise.\nRecord:\n  - item: %s\n  - judged-by: agent (%s)\n", ref, reason, ref, reason)
}
func ideaAnsweredTranscript(ref, verdict string, answers [3]string, values [3]float64) string {
	return ideaAnsweredTranscriptWithMonth(ref, verdict, answers, values, "0.0000")
}
func ideaAnsweredTranscriptWithMonth(ref, verdict string, answers [3]string, values [3]float64, month string) string {
	return fmt.Sprintf("Idea stage for %s: %s\n  shape_open: %s (%.2f) \u2014 the right shape of the solution is still open\n  unverified_assumption: %s (%.2f) \u2014 it rests on an assumption about its users, the business or the market that it gives no evidence for\n  market_question: %s (%.2f) \u2014 deciding to build it depends on a market or competitive question\nRule: indicated when any criterion is yes, not indicated when every criterion is no, undecided otherwise; yes at 0.70 or more, no at 0.30 or less.\nRecord:\n  - item: %s\n  - judged-by: judge (jev-1.13.0)\n  - shape_open: %s (%.2f)\n  - unverified_assumption: %s (%.2f)\n  - market_question: %s (%.2f)\n  - verdict: %s\nJudge: 1 call(s), 512 input tokens, US$0.0000; month US$%s of US$5.00; model jev-1.13.0 via typesafe on ROUNDFIX_TYPESAFE_API_KEY\n", ref, verdict, answers[0], values[0], answers[1], values[1], answers[2], values[2], ref, answers[0], values[0], answers[1], values[1], answers[2], values[2], verdict, month)
}
func ideaAnsweredJSON(ref, verdict string, answers [3]string, values [3]float64) string {
	return fmt.Sprintf("{\"schema\":\"roundfix/idea-stage-judgment/v1\",\"item\":{\"ref\":%q,\"kind\":\"backlog\",\"type\":\"feat\",\"title\":\"Teams can see each other's Run results\"},\"state\":\"answered\",\"reason\":null,\"verdict\":%q,\"judged_by\":\"judge\",\"criteria\":[{\"id\":\"shape_open\",\"summary\":\"the right shape of the solution is still open\",\"yes\":\"At least two materially different designs could satisfy the item, and the item does not settle which one to build.\",\"no\":\"The item names one concrete design, or only one reasonable design exists.\",\"probability\":%g,\"answer\":%q},{\"id\":\"unverified_assumption\",\"summary\":\"it rests on an assumption about its users, the business or the market that it gives no evidence for\",\"yes\":\"A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong.\",\"no\":\"The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made.\",\"probability\":%g,\"answer\":%q},{\"id\":\"market_question\",\"summary\":\"deciding to build it depends on a market or competitive question\",\"yes\":\"The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it.\",\"no\":\"The decision depends only on the product's own behavior and its current users.\",\"probability\":%g,\"answer\":%q}],\"thresholds\":{\"yes_at_least\":0.7,\"no_at_most\":0.3},\"model\":\"jev-1.13.0\",\"transport\":\"typesafe\",\"key_variable\":\"ROUNDFIX_TYPESAFE_API_KEY\",\"calls\":1,\"input_tokens\":512,\"cost_usd\":0.000021504,\"month_cost_usd\":0.000021504,\"month_ceiling_usd\":5}\n", ref, verdict, values[0], answers[0], values[1], answers[1], values[2], answers[2])
}

func TestSpecJudgeIdeaUnavailableReportsExitZero(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Run", "ceiling", "unreadable log", "secret", "service", "outside pin"} {
		t.Run(name, func(t *testing.T) {
			env, fake := ideaCommandFixture(t, true)
			reason := ""
			wantCalls := 0
			switch name {
			case "Run":
				env.environ = append(env.environ, "ROUNDFIX_RUN_ID=fixture-run")
				reason = "inside Run fixture-run, named by ROUNDFIX_RUN_ID"
			case "ceiling", "unreadable log":
				path := filepath.Join(env.homeDir, ".roundfix", "judge", "2026-11.jsonl")
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				content := "{\"cost_usd\":5}\n"
				if name == "unreadable log" {
					content = "broken json\n"
				}
				mustWrite(t, path, content)
				if name == "ceiling" {
					reason = "monthly ceiling reached (US$5.0000 of US$5.00)"
				} else {
					_, err := judge.ReadMonth(context.Background(), env.homeDir, env.dependencies.judgeNow())
					reason = "judge log unreadable: " + err.Error()
				}
			case "secret":
				secret := "ghp_" + strings.Repeat("a1", 18)
				followUpWrite(t, env.workDir, ideaFeatureRef, "---\ntype: feat\nstatus: open\n---\n# "+ideaTitle+"\n"+secret)
				reason = "not sent: secret pattern github-token"
			case "service":
				fake.status = 500
				reason = "service unavailable (HTTP 500)"
				wantCalls = 1
			case "outside pin":
				fake.model = "jev-1.14.0"
				reason = "answered by jev-1.14.0, thresholds belong to jev-1.13"
				wantCalls = 1
			}
			want := ideaUnavailableTranscript(ideaFeatureRef, reason)
			if wantCalls > 0 {
				want += "Judge: 1 call(s), 0 input tokens, US$0.0000; month US$0.0000 of US$5.00; model  via typesafe on ROUNDFIX_TYPESAFE_API_KEY\n"
				if name == "outside pin" {
					want = strings.TrimSuffix(want, "Judge: 1 call(s), 0 input tokens, US$0.0000; month US$0.0000 of US$5.00; model  via typesafe on ROUNDFIX_TYPESAFE_API_KEY\n") + "Judge: 1 call(s), 512 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.14.0 via typesafe on ROUNDFIX_TYPESAFE_API_KEY\n"
				}
			}
			specJudgeRun(t, env, []string{"--stage=idea", "--item", ideaFeatureRef}, want, "", 0)
			if fake.calls != wantCalls {
				t.Fatalf("calls=%d want=%d", fake.calls, wantCalls)
			}
		})
	}
}
