package cli

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestFollowUpRanksADeclaredBlockingMessageFirstWithoutTheJudge(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	f.judge.probabilities = map[string][4]float64{
		"Judged blocker": {.9, .9, .9, .1},
		"Judged feature": {.1, .1, .1, .9},
	}
	setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
	path := followUpWriteMessages(t, []map[string]string{
		{"id": "d2", "from": "quaestor", "created_at": "2026-10-20", "subject": "Totals question", "text": "Which totals rule applies?", "blocked_work": "Spec 0120 task_02 waits on the totals rule"},
		{"id": "j1", "from": "fluxus", "created_at": "2026-10-01", "subject": "Judged blocker", "text": "Our implementation cannot continue until you answer.", "type": "fix"},
		{"id": "d1", "from": "argus", "created_at": "2026-10-10", "subject": "Export field", "text": "Which field holds the export date?", "type": "information", "blocked_work": "  release 1.4 waits on the answer  "},
		{"id": "j2", "from": "fluxus", "created_at": "2026-10-02", "subject": "Judged feature", "text": "Production work.", "type": "feat"},
	})
	doc := followUpJSON(t, "--messages", path)
	if len(doc.Items) != 4 {
		t.Fatalf("items=%+v", doc.Items)
	}
	wantRefs := []string{"message:d1", "message:d2", "message:j1", "message:j2"}
	for index, item := range doc.Items {
		if item.Ref != wantRefs[index] || item.Rank != index+1 {
			t.Fatalf("item %d=%+v, want ref=%s rank=%d", index, item, wantRefs[index], index+1)
		}
		if index < 2 {
			wantWork := []string{"release 1.4 waits on the answer", "Spec 0120 task_02 waits on the totals rule"}[index]
			if item.Priority != 1 || item.PriorityName != "unblocking" || item.Basis != "declared" || item.Blocks == nil || item.RequestedBy == nil || *item.Blocks != *item.RequestedBy || item.BlockedWork != wantWork || item.Answers != nil || item.Reason != nil {
				t.Fatalf("declared item=%+v", item)
			}
		} else if item.Basis != "judge" || item.BlockedWork != "" {
			t.Fatalf("judged item=%+v", item)
		}
	}
	if f.judge.calls != 2 || !reflect.DeepEqual(f.judge.titles, []string{"Judged blocker", "Judged feature"}) || doc.Judge.State != "answered" {
		t.Fatalf("judge=%+v calls=%d titles=%v", doc.Judge, f.judge.calls, f.judge.titles)
	}
	out := renderFollowUpText(doc)
	if !strings.Contains(out, "  1. unblocking information 2026-10-10 message:d1: Export field — from argus — blocks argus: release 1.4 waits on the answer\n") || doc.Judge.Model == nil || !strings.Contains(out, "Judge: answered 2 of 2 item(s) with "+*doc.Judge.Model+"\n") {
		t.Fatalf("stdout=%s", out)
	}
}

func TestFollowUpRefusesAMalformedBlockedWork(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, from, work, reason string }{
		{"line break", "argus", "first\nsecond", "blocked_work must be one line"},
		{"without from", "  ", "work", "blocked_work needs from, the project whose work is blocked"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			newFollowUpFixture(t)
			path := followUpWriteMessages(t, []map[string]string{{"id": "b", "from": tc.from, "created_at": "2026-10-01", "text": "t", "blocked_work": tc.work}})
			out, stderr, code := followUpRunCLI(t, "--messages", path)
			wantErr := fmt.Sprintf("roundfix: follow-up failed: messages file %q: message %q: %s\nRun 'roundfix follow-up --help' for usage.\n", path, "b", tc.reason)
			if code != 2 || out != "" || stderr != wantErr {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want exit=2 stdout=empty stderr=%q", code, out, stderr, wantErr)
			}
		})
	}
}

func TestFollowUpOnlyDeclaredItemsAskNoJudge(t *testing.T) {
	t.Parallel()
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("all=%t", all), func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			setCommandEnvForTest(t, "ROUNDFIX_TYPESAFE_API_KEY", "fixture-key")
			path := followUpWriteMessages(t, []map[string]string{
				{"id": "b", "from": "argus", "created_at": "2026-10-10", "text": "Question", "blocked_work": "work b"},
				{"id": "a", "from": "fluxus", "created_at": "2026-10-10", "text": "Question", "blocked_work": "work a"},
			})
			args := []string{"--messages", path}
			if all {
				followUpWrite(t, f.repo, "docs/backlog/waiting.md", "---\nstatus: open\ntype: feat\nwaiting_on: supplier\nsince: 2026-10-25\n---\n# Waiting\n")
				followUpCommitWaitingFixture(t, f)
				args = append(args, "--all")
			}
			doc := followUpJSON(t, args...)
			if f.judge.calls != 0 || doc.Judge.Calls != 0 || doc.Judge.State != "not_asked" || doc.Judge.Reason == nil || *doc.Judge.Reason != "only declared or waiting items" {
				t.Fatalf("judge=%+v calls=%d", doc.Judge, f.judge.calls)
			}
			wantCount := 2
			if all {
				wantCount++
			}
			if len(doc.Items) != wantCount || doc.Items[0].Ref != "message:a" || doc.Items[1].Ref != "message:b" {
				t.Fatalf("items=%+v", doc.Items)
			}
			if all && (doc.Items[2].Basis != "waiting" || doc.Items[2].Rank != 3) {
				t.Fatalf("waiting item=%+v", doc.Items[2])
			}
		})
	}
}

func TestFollowUpDeclaredBlockedWorkTranscripts(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"declared first", "without from"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := newFollowUpFixture(t)
			messages := []map[string]string{
				{"id": "j1", "from": "fluxus", "created_at": "2026-10-01", "subject": "Broken export", "text": "The export fails.", "type": "fix"},
				{"id": "d2", "from": "quaestor", "created_at": "2026-10-20", "subject": "Totals question", "text": "Which totals rule applies?", "blocked_work": "Spec 0120 task_02 waits on the totals rule"},
				{"id": "d1", "from": "argus", "created_at": "2026-10-10", "subject": "Export field", "text": "Which field holds the export date?", "type": "information", "blocked_work": "release 1.4 waits on the answer"},
			}
			if name == "without from" {
				messages = []map[string]string{{"id": "b", "created_at": "2026-10-01", "text": "t", "blocked_work": "work"}}
			}
			path := followUpWriteMessages(t, messages)
			wantOut, wantErr, wantCode := "", fmt.Sprintf("roundfix: follow-up failed: messages file %q: message \"b\": blocked_work needs from, the project whose work is blocked\nRun 'roundfix follow-up --help' for usage.\n", path), 2
			if name == "declared first" {
				wantCode, wantErr = 0, ""
				wantOut = fmt.Sprintf(`Follow-up for widget in %s
Items: 3
  1. unblocking information 2026-10-10 message:d1: Export field — from argus — blocks argus: release 1.4 waits on the answer
  2. unblocking message 2026-10-20 message:d2: Totals question — from quaestor — blocks quaestor: Spec 0120 task_02 waits on the totals rule
  3. health fix 2026-10-01 message:j1: Broken export — from fluxus
In progress: 0
Judge: unavailable: ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set; ordered by type
Sources:
  ai-memory-messages: checked 3
  backlog: checked 0
  findings: checked 0
  repository-inbox: checked 0
  specs: checked 0
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: checked 0
  git: checked 0
`, f.repo)
			}
			out, stderr, code := followUpRunCLI(t, "--messages", path)
			if code != wantCode || out != wantOut || stderr != wantErr {
				t.Fatalf("exit=%d want=%d\nstdout=%q\nwant stdout=%q\nstderr=%q\nwant stderr=%q", code, wantCode, out, wantOut, stderr, wantErr)
			}
			if f.judge.calls != 0 {
				t.Fatalf("judge calls=%d", f.judge.calls)
			}
		})
	}
}
