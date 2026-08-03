package tui

import (
	"strings"
	"testing"

	"roundfix/internal/runevent"

	"charm.land/lipgloss/v2"
)

// Lip Gloss v2 renders styles with deterministic ANSI (downsampling happens
// at the program writer), so these tests pin the escape sequences directly:
// the forced color profile the techspec requires.
func ansi256(index string) string {
	return "\x1b[38;5;" + index + "m"
}

func TestResolveTokensStyledRendersDocumentedColors(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(true)
	tests := []struct {
		name  string
		style lipgloss.Style
		color string
	}{
		{"section label cyan", tokens.SectionLabel, "39"},
		{"done green", tokens.Done, "78"},
		{"running amber", tokens.Running, "214"},
		{"waiting amber", tokens.Waiting, "214"},
		{"pending amber", tokens.Pending, "214"},
		{"blocked red", tokens.Blocked, "203"},
		{"failed red", tokens.Failed, "203"},
		{"locked red", tokens.Locked, "203"},
		{"muted gray", tokens.Muted, "244"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rendered := test.style.Render("state")
			if !strings.Contains(rendered, ansi256(test.color)) {
				t.Fatalf("rendered %q, want color sequence %q", rendered, ansi256(test.color))
			}
			if stripANSI(rendered) != "state" {
				t.Fatalf("styled render changed text: %q", stripANSI(rendered))
			}
		})
	}
}

func TestResolveTokensIdentityWithoutColor(t *testing.T) {
	t.Parallel()
	tokens := ResolveTokens(false)
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"section label", tokens.SectionLabel},
		{"done", tokens.Done},
		{"running", tokens.Running},
		{"waiting", tokens.Waiting},
		{"pending", tokens.Pending},
		{"blocked", tokens.Blocked},
		{"failed", tokens.Failed},
		{"locked", tokens.Locked},
		{"muted", tokens.Muted},
	}
	const text = "[locked] docs/specs/task_01.md 12:04:05"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.style.Render(text); got != text {
				t.Fatalf("identity render = %q, want unstyled %q", got, text)
			}
		})
	}
}

func TestResolveTokensBorderTokensKeepStructure(t *testing.T) {
	t.Parallel()
	styled := ResolveTokens(true)
	identity := ResolveTokens(false)
	tests := []struct {
		name     string
		styled   lipgloss.Style
		identity lipgloss.Style
	}{
		{"active border", styled.ActiveBorder, identity.ActiveBorder},
		{"selection", styled.Selection, identity.Selection},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			colored := test.styled.Render("card")
			plain := test.identity.Render("card")
			if !strings.Contains(colored, ansi256("39")) {
				t.Fatalf("styled border missing cyan sequence: %q", colored)
			}
			if strings.Contains(plain, "\x1b") {
				t.Fatalf("identity border carries ANSI: %q", plain)
			}
			if stripANSI(colored) != plain {
				t.Fatalf("border structure differs between modes:\nstyled: %q\nidentity: %q", stripANSI(colored), plain)
			}
			if !strings.Contains(plain, "│") {
				t.Fatalf("identity border lost its frame: %q", plain)
			}
		})
	}
}

func TestEventSummaryRendersOneBoundedLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		summary string
		payload string
		width   int
		want    string
	}{
		{
			name:    "multi-line oversized summary",
			summary: strings.Repeat("a", 50) + "\nsecond line\nthird line",
			width:   20,
			want:    strings.Repeat("a", 19) + "…",
		},
		{
			name:    "short summary unchanged",
			summary: "verify passed",
			width:   40,
			want:    "verify passed",
		},
		{
			name:    "crlf line endings",
			summary: "batch 002 settled\r\ndetail",
			width:   40,
			want:    "batch 002 settled",
		},
		{
			name:    "payload never rendered",
			summary: "",
			payload: `{"text":"raw payload text"}`,
			width:   40,
			want:    "",
		},
		{
			name:    "zero width",
			summary: "anything",
			width:   0,
			want:    "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event := runevent.RunEvent{
				Source:  runevent.SourceDaemon,
				Kind:    runevent.KindDaemonBatch,
				Summary: test.summary,
				Payload: []byte(test.payload),
			}
			got := EventSummary(event, test.width)
			if got != test.want {
				t.Fatalf("EventSummary = %q, want %q", got, test.want)
			}
			if strings.ContainsAny(got, "\r\n") {
				t.Fatalf("EventSummary returned more than one line: %q", got)
			}
			if width := displayWidth(got); width > test.width {
				t.Fatalf("EventSummary width %d exceeds bound %d", width, test.width)
			}
		})
	}
}
