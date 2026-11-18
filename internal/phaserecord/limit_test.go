package phaserecord

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestReadUsageLimitClassifiesTheHarnessMessages(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, message string
		kind          LimitKind
		limit         string
	}{
		{"session", "You've hit your session limit · resets 3:45pm", "shared", "session"},
		{"weekly", "You've hit your weekly limit · resets Mon 12:00am", "shared", "weekly"},
		{"opus", "You've hit your Opus limit · resets 3:45pm", "model", "Opus"},
		{"sonnet", "You've hit your Sonnet limit · resets 3:45pm", "model", "Sonnet"},
		{"monthly_spend", "You've hit your monthly spend limit", "other", "monthly spend"},
		{"individual_spend", "You've hit your individual spend limit", "other", "individual spend"},
		{"individual_usage", "You've hit your individual usage limit", "other", "individual usage"},
		{"consent", "Fable limit reached · continuing on Fable 5.1 uses usage credits, and the prompt to confirm went unanswered · …", "model", "Fable"},
		{"0404_observation", "You've hit your Fable limit · resets Nov 22 at 8pm", "model", "Fable"},
		{"curly", "You’ve hit your Fable limit", "model", "Fable"},
		{"upper", "YOU'VE HIT YOUR FABLE LIMIT", "model", "FABLE"},
		{"lower", "you've hit your fable limit", "model", "fable"},
		{"reached_upper", "FABLE LIMIT REACHED", "model", "FABLE"},
		{"multi_line", "API error\nYou've hit your Fable limit\nYou've hit your weekly limit", "model", "Fable"},
		{"first_shared", "You've hit your weekly limit\nFable limit reached", "shared", "weekly"},
		{"other_reached", "Monthly spend limit reached", "other", "Monthly spend"},
		{"none", "rate limit exceeded", "none", ""},
		{"not_at_start", "Error: Fable limit reached", "other", "Error: Fable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ReadUsageLimit(tc.message)
			wantLine := tc.message
			if tc.name == "multi_line" {
				wantLine = "You've hit your Fable limit"
			}
			if tc.name == "first_shared" {
				wantLine = "You've hit your weekly limit"
			}
			if tc.kind == "none" {
				wantLine = ""
			}
			if got != (UsageLimit{Kind: tc.kind, Name: tc.limit, Line: wantLine}) {
				t.Fatalf("got %+v, want kind %s name %q line %q", got, tc.kind, tc.limit, wantLine)
			}
		})
	}
	t.Run("controls", func(t *testing.T) {
		t.Parallel()
		got := ReadUsageLimit("\tYou've hit your Fable limit\x00\x1b · reset\r")
		if got.Line != "You've hit your Fable limit   · reset" {
			t.Fatalf("line = %q", got.Line)
		}
	})
	t.Run("300_characters", func(t *testing.T) {
		t.Parallel()
		message := "You've hit your Fable limit " + strings.Repeat("é", 273)
		got := ReadUsageLimit(message)
		if !utf8.ValidString(got.Line) || len([]rune(got.Line)) != 200 || got.Line != string([]rune(message)[:200]) {
			t.Fatalf("line = %q", got.Line)
		}
	})
}
