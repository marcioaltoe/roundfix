// Suite: Follow-up Command list view.
// Invariant: pending work is ordered deterministically, incomplete sources are visible,
// and a report never mutates its input repositories or Run Database.
// Boundary IN: runCLI, real temporary git repositories, YAML and SQLite.
// Boundary OUT: GitHub and judge HTTP use per-test fakes only.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/judge"
	"roundfix/internal/store"
)

type followUpFakeGH struct {
	t              *testing.T
	root, response string
	err            error
}

func (fake followUpFakeGH) RunGH(_ context.Context, root string, args ...string) (string, error) {
	fake.t.Helper()
	if root != fake.root || !reflect.DeepEqual(args, []string{"pr", "list", "--state", "open", "--limit", "100", "--json", "number,title,headRefName,url,isDraft"}) {
		fake.t.Fatalf("gh root=%s args=%q", root, args)
	}
	return fake.response, fake.err
}

type followUpFakeJudge struct {
	t             *testing.T
	calls         int
	titles        []string
	probabilities map[string][4]float64
	failAt        int
}

func (fake *followUpFakeJudge) RoundTrip(request *http.Request) (*http.Response, error) {
	fake.t.Helper()
	fake.calls++
	if fake.probabilities == nil {
		fake.t.Fatal("unexpected judge request")
	}
	if fake.failAt == fake.calls {
		return &http.Response{StatusCode: 503, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{}"))}, nil
	}
	var body struct {
		Questions map[string]json.RawMessage `json:"questions"`
		State     struct {
			Project string `json:"project"`
			Item    struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			} `json:"item"`
		} `json:"state"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		fake.t.Fatal(err)
	}
	if body.State.Project != "widget" || len(body.Questions) != 4 || request.URL.Host != "api.typesafe.ai" || request.Header.Get("Authorization") != "Bearer fixture-key" {
		fake.t.Fatalf("request: host=%s body=%+v", request.URL.Host, body)
	}
	fake.titles = append(fake.titles, body.State.Item.Title)
	values, ok := fake.probabilities[body.State.Item.Title]
	if !ok {
		fake.t.Fatalf("unexpected item %q", body.State.Item.Title)
	}
	answers := make(map[string]any)
	for index, key := range []string{"is_fix", "blocks_other_project", "is_small", "affects_production"} {
		if _, ok := body.Questions[key]; !ok {
			fake.t.Fatalf("missing question %s", key)
		}
		answers[key] = map[string]any{"type": "noul", "noul": values[index]}
	}
	q, err := judge.Load()
	if err != nil {
		fake.t.Fatal(err)
	}
	model := q.Transports[len(q.Transports)-1].RequestModel
	// The direct transport returns a versioned model accepted by the current pin.
	response, err := json.Marshal(map[string]any{"model": model, "answers": answers, "usage": map[string]any{"input_tokens": 100, "output_tokens": 10}})
	if err != nil {
		fake.t.Fatal(err)
	}
	return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(response))}, nil
}

type followUpFixture struct {
	repo, home, brain string
	judge             *followUpFakeJudge
}

func newFollowUpFixture(t *testing.T) followUpFixture {
	t.Helper()
	repo, home, brain := t.TempDir(), t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "--initial-branch=main")
	var err error
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	followUpWrite(t, repo, "README.md", "# Widget\n")
	if err := os.MkdirAll(filepath.Join(repo, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, repo, "add", ".")
	gitImplement(t, repo, "commit", "-m", "feat: seed widget")
	followUpWrite(t, home, ".roundfix/config.yml", fmt.Sprintf("secondbrain:\n  root: %q\n  inbox_namespace: widget\n", brain))
	if err := os.MkdirAll(filepath.Join(brain, "inbox/widget"), 0755); err != nil {
		t.Fatal(err)
	}
	fake := &followUpFakeJudge{t: t}
	setCommandEnvironmentForTest(t, home, repo)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.environ = func() []string { return []string{} }
		deps.getenv = func(string) string { return "" }
		deps.followUpGHRunner = followUpFakeGH{t: t, root: repo, response: "[]"}
		deps.judgeTransport = fake
		deps.judgeNow = func() time.Time { return time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC) }
	})
	return followUpFixture{repo, home, brain, fake}
}
func followUpWrite(t *testing.T, root, path, content string) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, target, content)
}
func followUpRunCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var out, err bytes.Buffer
	code := runCLI(t, append([]string{"follow-up"}, args...), &out, &err)
	return out.String(), err.String(), code
}
func followUpJSON(t *testing.T, args ...string) followUpDocument {
	t.Helper()
	out, stderr, code := followUpRunCLI(t, append(args, "--format", "json")...)
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
	}
	var doc followUpDocument
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Schema != "roundfix/follow-up/v1" || doc.View != "follow-up" {
		t.Fatalf("document=%+v", doc)
	}
	want := []string{"secondbrain-inbox", "ai-memory-messages", "backlog", "findings", "repository-inbox", "specs", "runs", "delivery-queue", "pull-requests", "git"}
	names := []string{}
	for _, source := range doc.Sources {
		names = append(names, source.Name)
	}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("sources=%v", names)
	}
	return doc
}
func followUpSeedPending(t *testing.T, f followUpFixture) string {
	t.Helper()
	followUpWrite(t, f.brain, "inbox/widget/a-new.md", "---\nstatus: pending\ntype-hint: fix\ncreated_at: 2026-10-22\norigin: fluxus\n---\n# New inbox\nINBOX BODY NEVER PRINTED\n")
	followUpWrite(t, f.brain, "inbox/widget/z-old.md", "---\nstatus: pending\ntype-hint: fix\ncreated_at: 2026-10-20\norigin: widget\n---\n# Old inbox\nINBOX BODY NEVER PRINTED\n")
	followUpWrite(t, f.brain, "inbox/widget/_triaged/old.md", "# Triaged\n")
	followUpWrite(t, f.brain, "inbox/widget/done.md", "---\nstatus: done\n---\n# Already done\n")
	target := filepath.Join(f.brain, "outside.md")
	mustWrite(t, target, "# Linked content\n")
	if err := os.Symlink(target, filepath.Join(f.brain, "inbox/widget/link.md")); err != nil {
		t.Fatal(err)
	}
	followUpWrite(t, f.repo, "docs/backlog/open.md", "---\nstatus: open\ntype: feat\ncreated: 2026-10-19\n---\n# Export widget\nBACKLOG BODY NEVER PRINTED\n")
	followUpWrite(t, f.repo, "docs/backlog/promoted.md", "---\nstatus: promoted\ntype: fix\n---\n# Promoted\n")
	followUpWrite(t, f.repo, "docs/findings/pending.md", "---\nstatus: pending\ncreated_at: 2026-10-18\n---\n# Broken widget\nFINDING BODY NEVER PRINTED\n")
	followUpWrite(t, f.repo, "docs/findings/terminal.md", "---\nstatus: done\n---\n# Terminal\n")
	followUpWrite(t, f.repo, "docs/specs/0001-widget/_prd.md", "---\nstatus: active\ncreated: 2026-10-17\n---\n# Widget Spec\n\nFirst paragraph.\n\nPrivate second paragraph.\n")
	followUpWrite(t, f.repo, "docs/specs/0002-old/_prd.md", "---\nstatus: completed\n---\n# Old Spec\n")
	path := filepath.Join(t.TempDir(), "messages.json")
	mustWrite(t, path, `{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"one","from":"other","created_at":"2026-10-16","subject":"Message subject","text":"MESSAGE BODY NEVER PRINTED","type":"refactor"}]}`)
	return path
}

func TestFollowUpListsEveryPendingItemOnce(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	messages := followUpSeedPending(t, f)
	doc := followUpJSON(t, "--messages", messages)
	if len(doc.Items) != 6 {
		t.Fatalf("items=%+v", doc.Items)
	}
	want := map[string]struct{ kind, typ, date, title, requested string }{
		"inbox/widget/z-old.md":    {"inbox", "fix", "2026-10-20", "Old inbox", ""},
		"inbox/widget/a-new.md":    {"inbox", "fix", "2026-10-22", "New inbox", "fluxus"},
		"docs/backlog/open.md":     {"backlog", "feat", "2026-10-19", "Export widget", ""},
		"docs/findings/pending.md": {"finding", "finding", "2026-10-18", "Broken widget", ""},
		"docs/specs/0001-widget":   {"spec", "spec", "2026-10-17", "Widget Spec", ""},
		"message:one":              {"message", "refactor", "2026-10-16", "Message subject", "other"},
	}
	inbox := []string{}
	for _, item := range doc.Items {
		expected, ok := want[item.Ref]
		if !ok {
			t.Fatalf("unexpected/duplicate item=%+v", item)
		}
		delete(want, item.Ref)
		requested := ""
		if item.RequestedBy != nil {
			requested = *item.RequestedBy
		}
		if item.Kind != expected.kind || item.Type != expected.typ || item.Date != expected.date || item.Title != expected.title || requested != expected.requested {
			t.Fatalf("item=%+v want=%+v", item, expected)
		}
		if item.Kind == "inbox" {
			inbox = append(inbox, item.Ref)
		}
	}
	if len(want) != 0 || !reflect.DeepEqual(inbox, []string{"inbox/widget/z-old.md", "inbox/widget/a-new.md"}) {
		t.Fatalf("missing=%v inbox=%v", want, inbox)
	}
	out, stderr, code := followUpRunCLI(t, "--messages", messages)
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%s", code, stderr)
	}
	if strings.Contains(out, "BODY NEVER PRINTED") || strings.Contains(out, "Private second") {
		t.Fatalf("body in report: %s", out)
	}
	if f.judge.calls != 0 {
		t.Fatal("judge called without key")
	}
}

func TestFollowUpOrdersItemsByTheFourAnswers(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	f.judge.probabilities = map[string][4]float64{
		"unblocking":       {.9, .9, .1, .1},
		"small unblocking": {.1, .9, .9, .1},
		"health":           {.9, .1, .1, .1},
		"unsure":           {.9, .5, .1, .1},
		"small health old": {.9, .1, .9, .1},
		"small health new": {.9, .1, .9, .1},
		"feature":          {.1, .1, .1, .9},
		"other":            {.1, .1, .1, .1},
	}
	entries := []struct{ title, date string }{
		{"other", "2026-10-01"}, {"feature", "2026-10-01"}, {"health", "2026-10-01"}, {"unblocking", "2026-10-01"},
		{"unsure", "2026-10-02"}, {"small health new", "2026-10-20"}, {"small health old", "2026-10-10"}, {"small unblocking", "2026-10-20"},
	}
	for index, entry := range entries {
		followUpWrite(t, f.brain, fmt.Sprintf("inbox/widget/%02d.md", index), fmt.Sprintf("---\ntype-hint: feat\ncreated_at: %s\norigin: fluxus\n---\n# %s\n\nWork description.\n", entry.date, entry.title))
	}
	setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
	doc := followUpJSON(t)
	titles := []string{}
	for _, item := range doc.Items {
		titles = append(titles, item.Title)
		wantClass := "health"
		switch item.Title {
		case "unblocking", "small unblocking":
			wantClass = "unblocking"
		case "feature":
			wantClass = "feature"
		case "other":
			wantClass = "other"
		}
		if item.PriorityName != wantClass || item.Basis != "judge" || item.Reason != nil || len(item.Answers) != 4 {
			t.Fatalf("item=%+v", item)
		}
		blocking := item.Title == "unblocking" || item.Title == "small unblocking"
		if blocking && (item.Blocks == nil || *item.Blocks != "fluxus") || !blocking && item.Blocks != nil {
			t.Fatalf("blocks=%v for %s", item.Blocks, item.Title)
		}
	}
	want := []string{"small unblocking", "unblocking", "small health old", "small health new", "health", "unsure", "feature", "other"}
	if !reflect.DeepEqual(titles, want) || f.judge.calls != len(entries) || doc.Judge.State != "answered" {
		t.Fatalf("titles=%v calls=%d judge=%+v", titles, f.judge.calls, doc.Judge)
	}
	// Inventory order, before judgment, follows inbox created_at then reference.
	wantCalls := []string{"other", "feature", "health", "unblocking", "unsure", "small health old", "small health new", "small unblocking"}
	if !reflect.DeepEqual(f.judge.titles, wantCalls) {
		t.Fatalf("judge order=%v", f.judge.titles)
	}
}

func TestFollowUpOrdersByTypeWithoutTheJudge(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"no key", "Run", "ceiling"} {
		t.Run(mode, func(t *testing.T) {
			f := newFollowUpFixture(t)
			messages := followUpSeedPending(t, f)
			reason := "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"
			switch mode {
			case "Run":
				setCommandEnvForTest(t, "ROUNDFIX_RUN_ID", "fixture-run")
				setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
				reason = "inside Run fixture-run"
			case "ceiling":
				setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
				followUpWrite(t, f.home, ".roundfix/config.yml", fmt.Sprintf("secondbrain:\n  root: %q\n  inbox_namespace: widget\njev:\n  monthly_ceiling_usd: 0.01\n", f.brain))
				followUpWrite(t, f.home, ".roundfix/judge/2026-10.jsonl", "{\"cost_usd\":0.02}\n")
				reason = "monthly ceiling reached"
			}
			doc := followUpJSON(t, "--messages", messages)
			priorities := []string{}
			for _, item := range doc.Items {
				priorities = append(priorities, item.PriorityName)
				if item.Basis != "type" || item.Answers != nil || item.Reason == nil || !strings.Contains(*item.Reason, reason) || item.Blocks != nil {
					t.Fatalf("item=%+v", item)
				}
			}
			if !reflect.DeepEqual(priorities, []string{"health", "health", "health", "feature", "other", "other"}) || f.judge.calls != 0 || doc.Judge.State != "unavailable" {
				t.Fatalf("priorities=%v calls=%d judge=%+v", priorities, f.judge.calls, doc.Judge)
			}
			out, stderr, code := followUpRunCLI(t, "--messages", messages)
			if code != 0 || stderr != "" || !strings.Contains(out, "Judge: unavailable: "+*doc.Judge.Reason+"; ordered by type\n") {
				t.Fatalf("exit=%d stderr=%q out=%s", code, stderr, out)
			}
		})
	}
}

func followUpSeedDatabase(t *testing.T, f followUpFixture, run, queue bool) {
	t.Helper()
	ctx := context.Background()
	writer, err := store.Open(ctx, f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	if run {
		if _, err := writer.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.repo, LocalBranch: "main", SpecSlug: "0001-widget"}); err != nil {
			t.Fatal(err)
		}
	}
	if queue {
		if _, err := writer.CreateDeliveryQueue(ctx, f.repo, []string{"0002-queued"}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFollowUpSaysNothingToDo(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	out, stderr, code := followUpRunCLI(t)
	want := fmt.Sprintf("Follow-up for widget in %s\nNothing to do: no pending item and nothing in progress.\nSources:\n  secondbrain-inbox: checked 0\n  ai-memory-messages: missing: no --messages file\n  backlog: checked 0\n  findings: checked 0\n  repository-inbox: checked 0\n  specs: checked 0\n  runs: checked 0\n  delivery-queue: checked 0\n  pull-requests: checked 0\n  git: checked 0\n", f.repo)
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr, out, want)
	}
	if doc := followUpJSON(t); doc.State != "nothing_to_do" || len(doc.Items) != 0 {
		t.Fatalf("doc=%+v", doc)
	}
	followUpSeedDatabase(t, f, true, false)
	doc := followUpJSON(t)
	if doc.State != "work_found" || len(doc.InProgress.Runs) != 1 || len(doc.Items) != 0 {
		t.Fatalf("doc=%+v", doc)
	}
	out, stderr, code = followUpRunCLI(t)
	if code != 0 || stderr != "" || !strings.Contains(out, "In progress: 1\n  run ") || !strings.Contains(out, "Judge: not asked: no item\n") {
		t.Fatalf("stdout=%s stderr=%s exit=%d", out, stderr, code)
	}
}

func TestFollowUpReportsAMissingSecondbrainRoot(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	followUpWrite(t, f.home, ".roundfix/config.yml", "secondbrain:\n  inbox_namespace: widget\n")
	followUpWrite(t, f.repo, "docs/backlog/item.md", "---\nstatus: open\ntype: fix\n---\n# Repair\n")
	doc := followUpJSON(t)
	if len(doc.Items) != 1 || doc.Sources[0].State != "missing" || doc.Sources[0].Detail != "secondbrain.root is not set in User Config" {
		t.Fatalf("doc=%+v", doc)
	}
}

func TestFollowUpNamespaceDefaultsToTheRepositoryName(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"origin", "directory", "configured"} {
		t.Run(mode, func(t *testing.T) {
			f := newFollowUpFixture(t)
			followUpWrite(t, f.home, ".roundfix/config.yml", "")
			want := filepath.Base(f.repo)
			if mode != "directory" {
				gitImplement(t, f.repo, "remote", "add", "origin", "git@github.com:owner/widget.git")
				want = "widget"
			}
			if mode == "configured" {
				followUpWrite(t, f.repo, ".roundfixrc.yml", "secondbrain:\n  inbox_namespace: configured\n")
				want = "configured"
			}
			doc := followUpJSON(t)
			if doc.Repository.Namespace != want {
				t.Fatalf("namespace=%s want=%s", doc.Repository.Namespace, want)
			}
		})
	}
}

func TestFollowUpReadsNoUnavailableSourceAsAFailure(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.followUpGHRunner = followUpFakeGH{t: t, root: f.repo, err: fmt.Errorf("authentication unavailable\nsecond private line")}
	})
	doc := followUpJSON(t)
	source := doc.Sources[8]
	if source.State != "unavailable" || source.Detail != "authentication unavailable" || source.Count != 0 {
		t.Fatalf("source=%+v", source)
	}
}

func followUpSnapshot(t *testing.T, roots ...string) map[string][]byte {
	t.Helper()
	snapshot := make(map[string][]byte)
	for index, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.Type().IsRegular() {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			snapshot[fmt.Sprintf("%d/%s", index, relative)] = data
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return snapshot
}

func TestFollowUpWritesNothing(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	messages := followUpSeedPending(t, f)
	followUpSeedDatabase(t, f, true, true)
	followUpWrite(t, f.repo, "README.md", "# Widget changed\n")
	before := followUpSnapshot(t, f.repo, f.brain, f.home)
	status := gitImplementOutput(t, f.repo, "--no-optional-locks", "status", "--porcelain")
	for _, format := range []string{"text", "json"} {
		out, stderr, code := followUpRunCLI(t, "--messages", messages, "--format", format)
		if code != 0 || stderr != "" || out == "" {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
		}
		if after := followUpSnapshot(t, f.repo, f.brain, f.home); !reflect.DeepEqual(before, after) {
			for path, data := range before {
				if !bytes.Equal(data, after[path]) {
					t.Errorf("changed bytes: %s (before %d after %d)", path, len(data), len(after[path]))
				}
			}
			for path := range after {
				if _, ok := before[path]; !ok {
					t.Errorf("new file: %s", path)
				}
			}
			t.Fatalf("%s run changed input file bytes", format)
		}
		if after := gitImplementOutput(t, f.repo, "--no-optional-locks", "status", "--porcelain"); after != status {
			t.Fatalf("status changed: before=%q after=%q", status, after)
		}
	}
}

func TestFollowUpRefusesAnUnknownFormat(t *testing.T) {
	t.Parallel()
	newFollowUpFixture(t)
	// Bad configuration and missing messages would fail if validation read first.
	setCommandHomeDirForTest(t, filepath.Join(t.TempDir(), "missing-home"))
	out, stderr, code := followUpRunCLI(t, "--format", "yaml")
	want := "roundfix: follow-up failed: unknown --format \"yaml\"; use text or json\nRun 'roundfix follow-up --help' for usage.\n"
	if out != "" || stderr != want || code != 2 {
		t.Fatalf("stdout=%q stderr=%q exit=%d", out, stderr, code)
	}
}

func TestFollowUpRefusesAnInvalidMessagesFile(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, content string }{
		{"missing", ""}, {"oversize", strings.Repeat("x", (1<<20)+1)},
		{"schema", `{"schema":"wrong","messages":[]}`},
		{"id", `{"schema":"roundfix/follow-up-messages/v1","messages":[{"text":"hello","created_at":"2026-10-01"}]}`},
		{"duplicate", `{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"one","text":"hello","created_at":"2026-10-01"},{"id":"one","text":"world","created_at":"2026-10-02"}]}`},
		{"date", `{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"one","text":"hello","created_at":"yesterday"}]}`},
		{"type", `{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"one","text":"hello","created_at":"2026-10-01","type":"bug"}]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			newFollowUpFixture(t)
			path := filepath.Join(t.TempDir(), "messages.json")
			if tc.name != "missing" {
				mustWrite(t, path, tc.content)
			}
			out, stderr, code := followUpRunCLI(t, "--messages", path)
			if out != "" || code != 2 || !strings.Contains(stderr, path) || !strings.HasSuffix(stderr, "Run 'roundfix follow-up --help' for usage.\n") {
				t.Fatalf("stdout=%s stderr=%s exit=%d", out, stderr, code)
			}
		})
	}
}

func TestFollowUpListHelpMatchesTranscript(t *testing.T) {
	t.Parallel()
	newFollowUpFixture(t)
	const want = `Usage:
  roundfix follow-up [--messages <path>] [--format <text|json>]

List this repository's pending work in Follow-up Priority order.
It is read-only: it never writes a file, commits, pushes or prompts.
Options:
  --messages  JSON file of ai-memory messages to list with the pending work
  --format    text (default) or json

Exits 0 for a report, including one with nothing to do; 2 for invalid usage,
configuration or messages file; 1 for an unreadable Run Database.
`
	out, stderr, code := followUpRunCLI(t, "--help")
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr, out)
	}
}

func TestFollowUpListsInProgressAndParkActions(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	ctx := context.Background()
	writer, err := store.Open(ctx, f.home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := writer.CreateDeliveryQueue(ctx, f.repo, []string{"0001-merged", "0002-parked"})
	if err != nil {
		t.Fatal(err)
	}
	queue.Items[0].Stage = store.DeliveryStageMerged
	queue.Items[1].Stage = store.DeliveryStageParked
	queue.Items[1].Blocker = "checks-timeout"
	for _, item := range queue.Items {
		if err := writer.UpdateDeliveryQueueItem(ctx, f.repo, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.followUpGHRunner = followUpFakeGH{t: t, root: f.repo, response: `[{"number":7,"title":"Open change","headRefName":"feat/widget","url":"https://example.invalid/pull/7","isDraft":true}]`}
	})
	doc := followUpJSON(t)
	if doc.State != "work_found" || len(doc.InProgress.DeliveryQueue) != 1 || len(doc.InProgress.PullRequests) != 1 {
		t.Fatalf("doc=%+v", doc)
	}
	parked := doc.InProgress.DeliveryQueue[0]
	if parked.Spec != "0002-parked" || parked.Position != queue.Items[1].Position || parked.Stage != "parked" || parked.Blocker != "checks-timeout" || parked.ParkClass != "environment" || !strings.Contains(parked.Next, "roundfix deliver retry 0002-parked") {
		t.Fatalf("parked=%+v", parked)
	}
	out, stderr, code := followUpRunCLI(t)
	if code != 0 || stderr != "" || !strings.Contains(out, fmt.Sprintf("  queue %d 0002-parked parked checks-timeout\n", queue.Items[1].Position)) || !strings.Contains(out, "  pull request #7 feat/widget: Open change\n") {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr, out)
	}
}

func TestFollowUpReadsCommittedWALWithoutWriting(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	ctx := context.Background()
	writer, err := store.Open(ctx, f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	if _, err := writer.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.repo, LocalBranch: "main", SpecSlug: "0001-live"}); err != nil {
		t.Fatal(err)
	}

	before := followUpSnapshot(t, f.repo, f.brain, f.home)
	doc := followUpJSON(t)
	if len(doc.InProgress.Runs) != 1 || doc.InProgress.Runs[0].Spec != "0001-live" {
		t.Fatalf("progress=%+v", doc.InProgress)
	}
	out, stderr, code := followUpRunCLI(t)
	if code != 0 || stderr != "" || !strings.Contains(out, "0001-live") {
		t.Fatalf("stdout=%s stderr=%s exit=%d", out, stderr, code)
	}
	if after := followUpSnapshot(t, f.repo, f.brain, f.home); !reflect.DeepEqual(before, after) {
		t.Fatal("live text/JSON reader changed an input file")
	}

}

func TestFollowUpRefusesAnUnreadableRunDatabase(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	followUpWrite(t, f.home, ".roundfix/roundfix.db", "invalid SQLite database")
	out, stderr, code := followUpRunCLI(t)
	if code != 1 || out != "" || !strings.HasPrefix(stderr, "roundfix: follow-up failed:") {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr, out)
	}
}

func TestFollowUpReportsSkippedItemsAndAStoppedJudge(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"secret", "service"} {
		t.Run(mode, func(t *testing.T) {
			f := newFollowUpFixture(t)
			f.judge.probabilities = map[string][4]float64{"First": {.9, .1, .1, .1}, "Second": {.1, .1, .1, .9}, "Third": {.1, .1, .1, .9}}
			if mode == "service" {
				f.judge.failAt = 2
			}
			for index, title := range []string{"First", "Second", "Third"} {
				body := "Work description."
				if mode == "secret" && index == 1 {
					body = "EXAMPLE_" + "API_KEY=" + strings.Repeat("a", 24) + "9"
				}
				followUpWrite(t, f.brain, fmt.Sprintf("inbox/widget/%d.md", index), fmt.Sprintf("---\ntype-hint: feat\n---\n# %s\n\n%s\n", title, body))
			}
			setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
			doc := followUpJSON(t)
			if doc.Judge.State != "partial" || doc.Judge.Calls != 2 || f.judge.calls != 2 || doc.Items[0].Basis != "judge" {
				t.Fatalf("doc=%+v calls=%d", doc, f.judge.calls)
			}
			for _, item := range doc.Items {
				if item.Title == "Second" && (item.Basis != "type" || item.Reason == nil) {
					t.Fatalf("second=%+v", item)
				}
			}
			rows, err := judge.ReadMonth(context.Background(), f.home, time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC))
			if err != nil {
				t.Fatal(err)
			}
			if len(rows) != 2 {
				t.Fatalf("log rows=%d", len(rows))
			}
			for _, row := range rows {
				if row.Outcome == "answered" && row.Model == "" {
					t.Fatal("Judge Log missing model")
				}
			}
		})
	}
}

func TestFollowUpListsRepositoryNotesAndMessageTitleFallback(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	followUpWrite(t, f.repo, "docs/_inbox/2026-10-01-note.md", "# Note title\n\nNote body.\n")
	followUpWrite(t, f.repo, "docs/_inbox/nested/undated.md", "Just a raw note.\n")
	messages := filepath.Join(t.TempDir(), "messages.json")
	title := strings.Repeat("é", 90)
	content, err := json.Marshal(map[string]any{"schema": "roundfix/follow-up-messages/v1", "messages": []map[string]string{{"id": "one", "created_at": "2026-10-02T00:00:00Z", "text": title + "\nbody"}}})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, messages, string(content))
	doc := followUpJSON(t, "--messages", messages)
	if len(doc.Items) != 3 || doc.Sources[4].Count != 2 {
		t.Fatalf("doc=%+v", doc)
	}
	if doc.Items[0].Kind != "note" || doc.Items[0].Date != "2026-10-01" || doc.Items[1].Title != strings.Repeat("é", 80) || doc.Items[2].Title != "undated" || doc.Items[2].Date != "" {
		t.Fatalf("items=%+v", doc.Items)
	}
}

func TestFollowUpSnapshotKeepsCommittedWALAndCleansUp(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	ctx := context.Background()
	writer, err := store.Open(ctx, f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	create := func(slug string) {
		t.Helper()
		if _, err := writer.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.repo, LocalBranch: "main", SpecSlug: slug}); err != nil {
			t.Fatal(err)
		}
	}
	create("0001-first")
	before := followUpSnapshot(t, f.home)
	reader, found, err := openFollowUpReader(ctx, f.home)
	if err != nil || !found {
		t.Fatalf("reader found=%v err=%v", found, err)
	}
	defer func() {
		if reader != nil {
			if err := reader.Close(); err != nil {
				t.Error(err)
			}
		}
	}()
	temporaryHome := reader.temporaryHome
	if temporaryHome == "" {
		t.Fatal("live database did not use a private snapshot")
	}
	if after := followUpSnapshot(t, f.home); !reflect.DeepEqual(before, after) {
		t.Fatal("snapshot modified its source")
	}
	create("0002-later")
	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{GitRoot: f.repo, States: store.StatesActive})
	if err != nil || len(runs) != 1 || runs[0].SpecSlug != "0001-first" {
		t.Fatalf("snapshot runs=%+v err=%v", runs, err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	reader = nil
	if _, err := os.Stat(temporaryHome); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("snapshot directory remains: %v", err)
	}
	doc := followUpJSON(t)
	if len(doc.InProgress.Runs) != 2 {
		t.Fatalf("new report did not see latest commits: %+v", doc.InProgress)
	}
}

func TestFollowUpSnapshotHonorsCancellation(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	writer, err := store.Open(context.Background(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	if _, err := writer.CreateRun(context.Background(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.repo, LocalBranch: "main", SpecSlug: "0001-live"}); err != nil {
		t.Fatal(err)
	}
	before := followUpSnapshot(t, f.home)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	reader, _, err := openFollowUpReader(ctx, f.home)
	if reader != nil {
		if closeErr := reader.Close(); closeErr != nil {
			t.Error(closeErr)
		}
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("snapshot cancellation=%v", err)
	}
	if after := followUpSnapshot(t, f.home); !reflect.DeepEqual(before, after) {
		t.Fatal("canceled snapshot modified its source")
	}
}
