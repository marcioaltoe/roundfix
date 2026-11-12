package cli

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFollowUpDeclaredPrecedenceHoldsOneEntryPerSender(t *testing.T) {
	t.Parallel()
	items := []followUpItem{
		{Kind: "message", Ref: "message:a3", RequestedBy: followUpString("argus"), Date: "2026-10-25", BlockedWork: "release 1.5 waits on the label"},
		{Kind: "message", Ref: "message:j1", RequestedBy: followUpString("fluxus"), Date: "2026-10-01"},
		{Kind: "message", Ref: "message:q2", RequestedBy: followUpString("quaestor"), Date: "2026-10-20", BlockedWork: "Spec 0120 task_02 waits on the totals rule"},
		{Kind: "message", Ref: "message:a2", RequestedBy: followUpString("argus"), Date: "2026-10-10", BlockedWork: "release 1.4 waits on another answer"},
		{Kind: "message", Ref: "message:a1", RequestedBy: followUpString("argus"), Date: "2026-10-10", BlockedWork: "release 1.4 waits on the answer"},
		{Kind: "message", Ref: "message:a5", RequestedBy: followUpString("argus "), Date: "2026-10-30", BlockedWork: "release 1.6 waits on the label"},
	}
	before := append([]followUpItem{}, items...)
	declared, ranked := followUpDeclaredPrecedence(items)
	if !reflect.DeepEqual(declared, []followUpItem{items[2], items[4]}) {
		t.Fatalf("declared=%+v, want q2 and a1 unchanged in input order", declared)
	}
	wantRanked := []followUpItem{items[0], items[1], items[3], items[5]}
	for index := range wantRanked {
		if index != 1 {
			wantRanked[index].PrecedenceHeldBy = "message:a1"
		}
	}
	if !reflect.DeepEqual(ranked, wantRanked) {
		t.Fatalf("ranked=%+v, want %+v", ranked, wantRanked)
	}
	if !reflect.DeepEqual(items, before) {
		t.Fatal("selection changed its input")
	}
	emptyDeclared, emptyRanked := followUpDeclaredPrecedence(nil)
	if emptyDeclared == nil || emptyRanked == nil || len(emptyDeclared) != 0 || len(emptyRanked) != 0 {
		t.Fatalf("empty selection: declared=%#v ranked=%#v", emptyDeclared, emptyRanked)
	}
}

func TestFollowUpSendsALaterDeclaredMessageToTheJudge(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	f.judge.probabilities = map[string][4]float64{
		"Second question": {.9, .1, .1, .1},
		"Judged feature":  {.1, .1, .1, .9},
	}
	setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
	path := followUpWriteMessages(t, []map[string]string{
		{"id": "d3", "from": "argus", "created_at": "2026-10-25", "subject": "Second question", "text": "Which label?", "blocked_work": "release 1.5 waits on the label"},
		{"id": "j1", "from": "fluxus", "created_at": "2026-10-01", "subject": "Judged feature", "text": "Production work.", "type": "feat"},
		{"id": "d1", "from": "argus", "created_at": "2026-10-10", "subject": "Export field", "text": "Which field holds the export date?", "type": "information", "blocked_work": "release 1.4 waits on the answer"},
	})
	doc := followUpJSON(t, "--messages", path)
	if len(doc.Items) != 3 {
		t.Fatalf("items=%+v", doc.Items)
	}
	for index, ref := range []string{"message:d1", "message:d3", "message:j1"} {
		if item := doc.Items[index]; item.Ref != ref || item.Rank != index+1 {
			t.Fatalf("item %d=%+v, want ref=%s rank=%d", index, item, ref, index+1)
		}
	}
	holder, later, feature := doc.Items[0], doc.Items[1], doc.Items[2]
	if holder.Basis != "declared" || holder.Priority != 1 || holder.PrecedenceHeldBy != "" || holder.Blocks == nil || *holder.Blocks != "argus" || holder.Answers != nil || holder.Reason != nil {
		t.Fatalf("holder=%+v", holder)
	}
	if later.Basis != "judge" || later.Priority != 2 || later.PriorityName != "health" || later.Blocks != nil || later.BlockedWork != "release 1.5 waits on the label" || later.PrecedenceHeldBy != "message:d1" || later.Answers["is_fix"] != "yes" || later.Reason != nil {
		t.Fatalf("later=%+v", later)
	}
	if feature.Basis != "judge" || feature.Priority != 3 || feature.PrecedenceHeldBy != "" {
		t.Fatalf("feature=%+v", feature)
	}
	if f.judge.calls != 2 || doc.Judge.Calls != 2 || doc.Judge.State != "answered" || !reflect.DeepEqual(f.judge.titles, []string{"Second question", "Judged feature"}) {
		t.Fatalf("judge=%+v calls=%d titles=%v", doc.Judge, f.judge.calls, f.judge.titles)
	}
	out := renderFollowUpText(doc)
	if !strings.Contains(out, "  2. health message 2026-10-25 message:d3: Second question — from argus — declares blocked work: release 1.5 waits on the label — precedence held by message:d1\n") || doc.Judge.Model == nil || !strings.Contains(out, "Judge: answered 2 of 2 item(s) with "+*doc.Judge.Model+"\n") {
		t.Fatalf("stdout=%s", out)
	}
}

func TestFollowUpLaterDeclaredMessageTranscript(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	path := followUpWriteMessages(t, []map[string]string{
		{"id": "d3", "from": "argus", "created_at": "2026-10-25", "subject": "Second question", "text": "Which label?", "blocked_work": "release 1.5 waits on the label"},
		{"id": "j1", "from": "fluxus", "created_at": "2026-10-01", "subject": "Broken export", "text": "The export fails.", "type": "fix"},
		{"id": "d2", "from": "quaestor", "created_at": "2026-10-20", "subject": "Totals question", "text": "Which totals rule applies?", "blocked_work": "Spec 0120 task_02 waits on the totals rule"},
		{"id": "d1", "from": "argus", "created_at": "2026-10-10", "subject": "Export field", "text": "Which field holds the export date?", "type": "information", "blocked_work": "release 1.4 waits on the answer"},
	})
	wantOut := fmt.Sprintf(`Follow-up for widget in %s
Items: 4
  1. unblocking information 2026-10-10 message:d1: Export field — from argus — blocks argus: release 1.4 waits on the answer
  2. unblocking message 2026-10-20 message:d2: Totals question — from quaestor — blocks quaestor: Spec 0120 task_02 waits on the totals rule
  3. health fix 2026-10-01 message:j1: Broken export — from fluxus
  4. other message 2026-10-25 message:d3: Second question — from argus — declares blocked work: release 1.5 waits on the label — precedence held by message:d1
In progress: 0
Judge: unavailable: ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set; ordered by type
Sources:
  ai-memory-messages: checked 4
  backlog: checked 0
  findings: checked 0
  repository-inbox: checked 0
  specs: checked 0
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: checked 0
  git: checked 0
`, f.repo)
	out, stderr, code := followUpRunCLI(t, "--messages", path)
	if code != 0 || out != wantOut || stderr != "" {
		t.Fatalf("exit=%d want=0\nstdout=%q\nwant stdout=%q\nstderr=%q want empty", code, out, wantOut, stderr)
	}
	if f.judge.calls != 0 {
		t.Fatalf("judge calls=%d", f.judge.calls)
	}
}
