// Suite: Waiting Items in the Follow-up Command.
// Invariant: declarations hold items back without judging or losing source accounting.
// Boundary IN: public CLI, committed disposable Git repositories and YAML.
// Boundary OUT: per-test fake GitHub and judge transports, no network.
package cli

import (
	"encoding/json"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func followUpCommitWaitingFixture(t *testing.T, f followUpFixture) {
	t.Helper()
	gitImplement(t, f.repo, "add", ".")
	// Replace the disposable seed commit with a reproducible root commit so
	// the characterization substitutes only the repository root, never HEAD.
	cmd := exec.Command("git", append(gitConfigArgsForTest(), "commit", "--amend", "--reset-author", "-m", "feat: seed follow-up fixture")...)
	cmd.Dir = f.repo
	cmd.Env = append(isolatedGitEnvForTest(), "GIT_AUTHOR_DATE=2026-10-23T12:00:00Z", "GIT_COMMITTER_DATE=2026-10-23T12:00:00Z")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit fixture: %v\n%s", err, out)
	}
}

func TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	followUpWrite(t, f.repo, "docs/backlog/open.md", "---\nstatus: open\ntype: feat\ncreated: 2026-10-19\n---\n# Export widget\n")
	followUpWrite(t, f.repo, "docs/findings/pending.md", "---\nstatus: pending\ncreated_at: 2026-10-18\n---\n# Broken widget\n")
	followUpWrite(t, f.repo, "docs/_inbox/2026-10-21-note.md", "# Visible note\n")
	followUpCommitWaitingFixture(t, f)
	messages := followUpWriteMessages(t, []map[string]string{{"id": "one", "from": "fluxus", "created_at": "2026-10-20", "subject": "Handed-in message", "text": "Pending work.", "type": "fix"}})
	// Recorded from the starting tree, before any production code change.
	const text = `Follow-up for widget in <root>
Items: 4
  1. health finding 2026-10-18 docs/findings/pending.md: Broken widget
  2. health fix 2026-10-20 message:one: Handed-in message — from fluxus
  3. feature feat 2026-10-19 docs/backlog/open.md: Export widget
  4. other note 2026-10-21 docs/_inbox/2026-10-21-note.md: Visible note
In progress: 0
Judge: unavailable: ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set; ordered by type
Sources:
  ai-memory-messages: checked 1
  backlog: checked 1
  findings: checked 1
  repository-inbox: checked 1
  specs: checked 0
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: checked 0
  git: checked 0
`
	const json = `{"schema":"roundfix/follow-up/v1","view":"follow-up","state":"work_found","repository":{"root":"<root>","namespace":"widget"},"sources":[{"name":"ai-memory-messages","state":"checked","count":1,"detail":""},{"name":"backlog","state":"checked","count":1,"detail":""},{"name":"findings","state":"checked","count":1,"detail":""},{"name":"repository-inbox","state":"checked","count":1,"detail":""},{"name":"specs","state":"checked","count":0,"detail":""},{"name":"runs","state":"checked","count":0,"detail":""},{"name":"delivery-queue","state":"checked","count":0,"detail":""},{"name":"pull-requests","state":"checked","count":0,"detail":""},{"name":"git","state":"checked","count":0,"detail":""}],"judge":{"state":"unavailable","reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set","model":null,"calls":0,"cost_usd":0},"items":[{"rank":1,"priority":2,"priority_name":"health","kind":"finding","ref":"docs/findings/pending.md","type":"finding","date":"2026-10-18","title":"Broken widget","requested_by":null,"blocks":null,"basis":"type","answers":null,"reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"},{"rank":2,"priority":2,"priority_name":"health","kind":"message","ref":"message:one","type":"fix","date":"2026-10-20","title":"Handed-in message","requested_by":"fluxus","blocks":null,"basis":"type","answers":null,"reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"},{"rank":3,"priority":3,"priority_name":"feature","kind":"backlog","ref":"docs/backlog/open.md","type":"feat","date":"2026-10-19","title":"Export widget","requested_by":null,"blocks":null,"basis":"type","answers":null,"reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"},{"rank":4,"priority":4,"priority_name":"other","kind":"note","ref":"docs/_inbox/2026-10-21-note.md","type":"note","date":"2026-10-21","title":"Visible note","requested_by":null,"blocks":null,"basis":"type","answers":null,"reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"}],"in_progress":{"runs":[],"delivery_queue":[],"pull_requests":[],"git":{"branch":"main","head":"76cb09c","upstream":"","unpushed_commits":0,"changed_paths":[]}}}
`
	for _, tc := range []struct{ format, want string }{{"text", text}, {"json", json}} {
		format := tc.format
		out, stderr, code := followUpRunCLI(t, "--messages", messages, "--format", format)
		if code != 0 || stderr != "" {
			t.Fatalf("format=%s exit=%d stderr=%q stdout=%s", format, code, stderr, out)
		}
		if want := strings.ReplaceAll(tc.want, "<root>", f.repo); out != want {
			t.Fatalf("format=%s stdout=%q want=%q", format, out, want)
		}
	}
}

func TestFollowUpHoldsBackWaitingItems(t *testing.T) {
	t.Parallel()
	t.Run("mixed inventory", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/findings/waiting.md", "---\nstatus: pending\nwaiting_on: '  ERP vendor support answer  '\nsince: 2026-10-25\n---\n# ERP totals differ\n")
		followUpWrite(t, f.repo, "docs/findings/ordinary.md", "---\nstatus: partial\n---\n# Broken import\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		f.judge.probabilities = map[string][4]float64{"Broken import": {.9, .1, .1, .1}}
		out, stderr, code := followUpRunCLI(t, "--format", "json")
		if code != 0 || stderr != "" {
			t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
		}
		var doc followUpDocument
		if err := json.Unmarshal([]byte(out), &doc); err != nil {
			t.Fatal(err)
		}
		if len(doc.Items) != 1 || doc.Items[0].Title != "Broken import" || doc.Items[0].Basis != "judge" || doc.Judge.State != "answered" {
			t.Fatalf("report=%+v", doc)
		}
		if source := doc.Sources[2]; source.State != "checked" || source.Count != 2 || source.Waiting != 1 {
			t.Fatalf("findings=%+v", source)
		}
		if f.judge.calls != 1 || !reflect.DeepEqual(f.judge.titles, []string{"Broken import"}) {
			t.Fatalf("judge calls=%d titles=%v", f.judge.calls, f.judge.titles)
		}
		var raw struct {
			Items []map[string]json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal([]byte(out), &raw); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"waiting", "waiting_on", "since", "waiting_error"} {
			if _, ok := raw.Items[0][key]; ok {
				t.Errorf("ordinary item carries %s: %s", key, out)
			}
		}
	})
	t.Run("only waiting items", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/backlog/waiting.md", "---\nstatus: open\ntype: feat\nwaiting_on: supplier sample data\nsince: 2026-10-25\n---\n# Waiting backlog\n")
		followUpWrite(t, f.repo, "docs/findings/waiting.md", "---\nstatus: partial\nwaiting_on: ERP vendor\nsince: 2026-10-26\n---\n# Waiting finding\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		doc := followUpJSON(t)
		if doc.State != "nothing_to_do" || len(doc.Items) != 0 || doc.Judge.State != "not_asked" || doc.Judge.Reason == nil || *doc.Judge.Reason != "no item" || f.judge.calls != 0 || len(f.judge.titles) != 0 {
			t.Fatalf("report=%+v judge titles=%v", doc, f.judge.titles)
		}
		for _, index := range []int{1, 2} {
			if source := doc.Sources[index]; source.Count != 1 || source.Waiting != 1 || source.State != "checked" {
				t.Fatalf("source=%+v", source)
			}
		}
		out, stderr, code := followUpRunCLI(t)
		if code != 0 || stderr != "" || !strings.Contains(out, "Nothing to do: no pending item and nothing in progress.\n") || !strings.Contains(out, "backlog: checked 1: 1 waiting, held back without --all\n") || !strings.Contains(out, "findings: checked 1: 1 waiting, held back without --all\n") {
			t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
		}
	})
	t.Run("notes and specs remain ranked", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/_inbox/note.md", "---\nwaiting_on: ignored on a note\nsince: 2026-10-25\n---\n# Note\n")
		followUpWrite(t, f.repo, "docs/specs/0001-widget/_prd.md", "---\nstatus: active\nwaiting_on: [ignored on a spec]\nsince: [invalid]\n---\n# Spec\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		f.judge.probabilities = map[string][4]float64{"Note": {.1, .1, .1, .1}, "Spec": {.1, .1, .1, .1}}
		doc := followUpJSON(t)
		if len(doc.Items) != 2 || !reflect.DeepEqual(f.judge.titles, []string{"Note", "Spec"}) {
			t.Fatalf("items=%+v judge titles=%v", doc.Items, f.judge.titles)
		}
		for _, item := range doc.Items {
			if item.WaitingOn != "" || item.Since != "" || item.WaitingError != "" || item.Basis != "judge" {
				t.Fatalf("item=%+v", item)
			}
		}
	})
}

func TestFollowUpAllListsWaitingItemsLast(t *testing.T) {
	t.Parallel()
	t.Run("ranked then since then reference", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/findings/a-later.md", "---\nstatus: pending\ncreated_at: 2026-10-01\nwaiting_on: ERP vendor support answer\nsince: 2026-10-26\n---\n# Later wait\n")
		followUpWrite(t, f.repo, "docs/findings/z-earlier.md", "---\nstatus: partial\ncreated_at: 2026-10-23\nwaiting_on: '  supplier sample data  '\nsince: '2026-10-25'\n---\n# Earlier wait\n")
		followUpWrite(t, f.repo, "docs/backlog/earlier.md", "---\nstatus: open\ntype: feat\nwaiting_on: maintainer answer\nsince: 2026-10-25\n---\n# Same date wait\n")
		followUpWrite(t, f.repo, "docs/backlog/ordinary.md", "---\nstatus: open\ntype: refactor\n---\n# Ordinary\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		f.judge.probabilities = map[string][4]float64{"Ordinary": {.1, .1, .1, .1}}
		doc := followUpJSON(t, "--all")
		if len(doc.Items) != 4 || doc.Items[0].Title != "Ordinary" || doc.Items[0].PriorityName != "other" || doc.Judge.State != "answered" || f.judge.calls != 1 || !reflect.DeepEqual(f.judge.titles, []string{"Ordinary"}) {
			t.Fatalf("report=%+v judge titles=%v", doc, f.judge.titles)
		}
		want := []struct{ ref, waitingOn, since string }{
			{"docs/backlog/earlier.md", "maintainer answer", "2026-10-25"},
			{"docs/findings/z-earlier.md", "supplier sample data", "2026-10-25"},
			{"docs/findings/a-later.md", "ERP vendor support answer", "2026-10-26"},
		}
		for index, expected := range want {
			item := doc.Items[index+1]
			if item.Rank != index+2 || item.Ref != expected.ref || item.Priority != 5 || item.PriorityName != "waiting" || item.Basis != "waiting" || item.WaitingOn != expected.waitingOn || item.Since != expected.since || item.WaitingError != "" || item.Answers != nil || item.Blocks != nil || item.Reason != nil {
				t.Fatalf("item=%+v want=%+v", item, expected)
			}
		}
		if doc.Sources[1].Count != 2 || doc.Sources[1].Waiting != 1 || doc.Sources[2].Count != 2 || doc.Sources[2].Waiting != 2 {
			t.Fatalf("sources=%+v", doc.Sources)
		}
		out, stderr, code := followUpRunCLI(t, "--all")
		if code != 0 || stderr != "" || !strings.Contains(out, "Judge: answered 1 of 1 item(s) with ") || strings.Contains(out, "ordered by type") || strings.Contains(out, "held back without --all") {
			t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
		}
		if !reflect.DeepEqual(f.judge.titles, []string{"Ordinary", "Ordinary"}) {
			t.Fatalf("text judge titles=%v", f.judge.titles)
		}
	})
	t.Run("partial judge totals exclude waiting", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/findings/waiting.md", "---\nstatus: pending\nwaiting_on: vendor\nsince: 2026-10-25\n---\n# Waiting\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		f.judge.probabilities = map[string][4]float64{"First": {.9, .1, .1, .1}, "Second": {.1, .1, .1, .1}, "Third": {.1, .1, .1, .1}}
		f.judge.failAt = 2
		messages := followUpWriteMessages(t, []map[string]string{
			{"id": "one", "created_at": "2026-10-21", "subject": "First", "text": "Work.", "type": "fix"},
			{"id": "two", "created_at": "2026-10-22", "subject": "Second", "text": "Work.", "type": "fix"},
			{"id": "three", "created_at": "2026-10-23", "subject": "Third", "text": "Work.", "type": "fix"},
		})
		doc := followUpJSON(t, "--all", "--messages", messages)
		if doc.Judge.State != "partial" || len(doc.Items) != 4 || doc.Items[3].PriorityName != "waiting" || f.judge.calls != 2 || !reflect.DeepEqual(f.judge.titles, []string{"First"}) {
			t.Fatalf("report=%+v calls=%d titles=%v", doc, f.judge.calls, f.judge.titles)
		}
		// Start another stopped judgment for the text rendering of the same inventory.
		f.judge.calls, f.judge.titles = 0, nil
		out, stderr, code := followUpRunCLI(t, "--all", "--messages", messages)
		if code != 0 || stderr != "" || !strings.Contains(out, "Judge: answered 1 of 3 item(s) with ") || !strings.Contains(out, "; 2 ordered by type\n") || !reflect.DeepEqual(f.judge.titles, []string{"First"}) {
			t.Fatalf("exit=%d stderr=%q stdout=%s judge titles=%v", code, stderr, out, f.judge.titles)
		}
	})
	t.Run("only waiting items", func(t *testing.T) {
		t.Parallel()
		f := newFollowUpFixture(t)
		followUpWrite(t, f.repo, "docs/findings/waiting.md", "---\nstatus: pending\nwaiting_on: ERP vendor support answer\nsince: 2026-10-25\n---\n# Waiting\n")
		followUpCommitWaitingFixture(t, f)
		setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
		doc := followUpJSON(t, "--all")
		if doc.State != "work_found" || len(doc.Items) != 1 || doc.Items[0].Priority != 5 || doc.Items[0].Rank != 1 || doc.Judge.State != "not_asked" || doc.Judge.Reason == nil || *doc.Judge.Reason != "only waiting items" || doc.Judge.Calls != 0 || f.judge.calls != 0 || len(f.judge.titles) != 0 {
			t.Fatalf("report=%+v judge titles=%v", doc, f.judge.titles)
		}
		out, stderr, code := followUpRunCLI(t, "--all")
		if code != 0 || stderr != "" || !strings.Contains(out, "Judge: not asked: only waiting items\n") {
			t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
		}
	})
}

func TestFollowUpReportsAnInvalidWaitingDeclaration(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, fields, waitingError string }{
		{"blank", "waiting_on: '  '\nsince: 2026-10-25\n", "waiting_on must be one line of non-empty text"},
		{"list", "waiting_on: [vendor]\nsince: 2026-10-25\n", "waiting_on must be one line of non-empty text"},
		{"map", "waiting_on: {vendor: answer}\n", "waiting_on must be one line of non-empty text"},
		{"two lines", "waiting_on: |-\n  vendor\n  answer\nsince: 2026-10-25\n", "waiting_on must be one line of non-empty text"},
		{"trailing newline", "waiting_on: |\n  vendor\nsince: 2026-10-25\n", "waiting_on must be one line of non-empty text"},
		{"missing since", "waiting_on: vendor\n", "since is missing"},
		{"null since", "waiting_on: vendor\nsince: null\n", "since is missing"},
		{"impossible date", "waiting_on: vendor\nsince: 2026-02-30\n", "since must be a YYYY-MM-DD date"},
		{"list date", "waiting_on: vendor\nsince: [x]\n", "since must be a YYYY-MM-DD date"},
		{"map date", "waiting_on: vendor\nsince: {date: x}\n", "since must be a YYYY-MM-DD date"},
		{"blank date", "waiting_on: vendor\nsince: ''\n", "since must be a YYYY-MM-DD date"},
		{"datetime", "waiting_on: vendor\nsince: 2026-10-25T12:00:00Z\n", "since must be a YYYY-MM-DD date"},
		{"since alone", "since: 2026-10-25\n", "since needs waiting_on"},
		{"since alone invalid", "since: [x]\n", "since needs waiting_on"},
		{"null waiting on", "waiting_on: null\nsince: 2026-10-25\n", "since needs waiting_on"},
		{"validation order", "waiting_on: []\nsince: []\n", "waiting_on must be one line of non-empty text"},
		{"null pair", "waiting_on: null\nsince: ~\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			followUpWrite(t, f.repo, "docs/backlog/item.md", "---\nstatus: open\ntype: feat\n"+tc.fields+"---\n# Backlog\n")
			followUpWrite(t, f.repo, "docs/findings/item.md", "---\nstatus: pending\n"+tc.fields+"---\n# Finding\n")
			followUpCommitWaitingFixture(t, f)
			setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
			f.judge.probabilities = map[string][4]float64{"Backlog": {.1, .1, .1, .9}, "Finding": {.9, .1, .1, .1}}
			messages := followUpWriteMessages(t, []map[string]string{})
			doc := followUpJSON(t, "--messages", messages)
			if len(doc.Items) != 2 || !reflect.DeepEqual(f.judge.titles, []string{"Backlog", "Finding"}) {
				t.Fatalf("items=%+v judge titles=%v", doc.Items, f.judge.titles)
			}
			for _, source := range doc.Sources {
				if source.State != "checked" || source.Waiting != 0 {
					t.Errorf("source=%+v", source)
				}
			}
			if doc.Sources[1].Count != 1 || doc.Sources[2].Count != 1 {
				t.Fatalf("sources=%+v", doc.Sources)
			}
			for _, item := range doc.Items {
				wantPriority := "health"
				if item.Kind == "backlog" {
					wantPriority = "feature"
				}
				if item.WaitingError != tc.waitingError || item.WaitingOn != "" || item.Since != "" || item.PriorityName != wantPriority || item.Basis != "judge" {
					t.Errorf("item=%+v want error=%q", item, tc.waitingError)
				}
			}
			out, stderr, code := followUpRunCLI(t, "--messages", messages)
			wantSuffixes := 2
			if tc.waitingError == "" {
				wantSuffixes = 0
			}
			if code != 0 || stderr != "" || strings.Count(out, " — invalid waiting declaration: "+tc.waitingError) != wantSuffixes {
				t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr, out)
			}
		})
	}
}

func followUpWaitingTranscriptFixture(t *testing.T) followUpFixture {
	t.Helper()
	f := newFollowUpFixture(t)
	followUpWrite(t, f.repo, "docs/_inbox/.gitkeep", "")
	followUpWrite(t, f.repo, "docs/backlog/2026-10-21-export-widget.md", "---\nstatus: open\ntype: feat\ncreated: 2026-10-21\n---\n# Export widget\n")
	followUpWrite(t, f.repo, "docs/findings/2026-10-20-erp-totals.md", "---\nstatus: pending\ncreated_at: 2026-10-20\nwaiting_on: ERP vendor support answer\nsince: 2026-10-25\n---\n# ERP totals differ\n")
	followUpWrite(t, f.repo, "docs/findings/2026-10-22-broken-import.md", "---\nstatus: pending\ncreated_at: 2026-10-22\n---\n# Broken import\n")
	followUpWrite(t, f.repo, "docs/findings/2026-10-23-late-feed.md", "---\nstatus: pending\ncreated_at: 2026-10-23\nwaiting_on: supplier sample data\n---\n# Late feed\n")
	followUpCommitWaitingFixture(t, f)
	return f
}

func TestFollowUpHoldsBackWaitingItemsTranscript(t *testing.T) {
	t.Parallel()
	f := followUpWaitingTranscriptFixture(t)
	const transcript = `Follow-up for <namespace> in <root>
Items: 3
  1. health finding 2026-10-22 docs/findings/2026-10-22-broken-import.md: Broken import
  2. health finding 2026-10-23 docs/findings/2026-10-23-late-feed.md: Late feed — invalid waiting declaration: since is missing
  3. feature feat 2026-10-21 docs/backlog/2026-10-21-export-widget.md: Export widget
In progress: 0
Judge: unavailable: <reason>; ordered by type
Sources:
  ai-memory-messages: missing: no --messages file
  backlog: checked 1
  findings: checked 3: 1 waiting, held back without --all
  repository-inbox: checked 0
  specs: checked 0
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: <pull request source>
  git: checked 0
`
	want := strings.NewReplacer("<namespace>", "widget", "<root>", f.repo, "<reason>", "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set", "<pull request source>", "checked 0").Replace(transcript)
	out, stderr, code := followUpRunCLI(t)
	if out != want || stderr != "" || code != 0 {
		t.Fatalf("stdout=%q stderr=%q exit=%d, want stdout=%q stderr=empty exit=0", out, stderr, code, want)
	}
}

func TestFollowUpAllListsWaitingItemsTranscript(t *testing.T) {
	t.Parallel()
	f := followUpWaitingTranscriptFixture(t)
	const transcript = `Follow-up for <namespace> in <root>
Items: 4
  1. health finding 2026-10-22 docs/findings/2026-10-22-broken-import.md: Broken import
  2. health finding 2026-10-23 docs/findings/2026-10-23-late-feed.md: Late feed — invalid waiting declaration: since is missing
  3. feature feat 2026-10-21 docs/backlog/2026-10-21-export-widget.md: Export widget
  4. waiting finding 2026-10-20 docs/findings/2026-10-20-erp-totals.md: ERP totals differ — waiting on ERP vendor support answer since 2026-10-25
In progress: 0
Judge: unavailable: <reason>; ordered by type
Sources:
  ai-memory-messages: missing: no --messages file
  backlog: checked 1
  findings: checked 3: 1 waiting
  repository-inbox: checked 0
  specs: checked 0
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: <pull request source>
  git: checked 0
`
	want := strings.NewReplacer("<namespace>", "widget", "<root>", f.repo, "<reason>", "ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set", "<pull request source>", "checked 0").Replace(transcript)
	out, stderr, code := followUpRunCLI(t, "--all")
	if out != want || stderr != "" || code != 0 {
		t.Fatalf("stdout=%q stderr=%q exit=%d, want stdout=%q stderr=empty exit=0", out, stderr, code, want)
	}
}

func TestFollowUpStatusRefusesAllTranscript(t *testing.T) {
	t.Parallel()
	newFollowUpFixture(t)
	out, stderr, code := followUpRunCLI(t, "--status", "--all")
	const want = "roundfix: follow-up failed: --all applies only to the follow-up list; drop it with --status\nRun 'roundfix follow-up --help' for usage.\n"
	if out != "" || stderr != want || code != 2 {
		t.Fatalf("stdout=%q stderr=%q exit=%d, want stdout=empty stderr=%q exit=2", out, stderr, code, want)
	}
}
