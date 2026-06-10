package runevent

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestBoundSummaryKeepsShortTextUntouched(t *testing.T) {
	if got := BoundSummary("short summary"); got != "short summary" {
		t.Fatalf("expected identity for short text, got %q", got)
	}
}

func TestBoundSummaryTruncatesOnRuneBoundary(t *testing.T) {
	huge := strings.Repeat("é", MaxSummaryLength)

	bounded := BoundSummary(huge)

	if len(bounded) > MaxSummaryLength+len("…") {
		t.Fatalf("expected bounded length, got %d bytes", len(bounded))
	}
	if !utf8.ValidString(bounded) {
		t.Fatal("expected truncation on a rune boundary")
	}
	if !strings.HasSuffix(bounded, "…") {
		t.Fatalf("expected truncation marker, got %q", bounded[len(bounded)-8:])
	}
}
