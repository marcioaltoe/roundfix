package runevent

import "unicode/utf8"

// MaxSummaryLength bounds event summaries in bytes so list views and
// plain-text output stay readable when tool output is huge.
const MaxSummaryLength = 2048

// BoundSummary truncates text to MaxSummaryLength bytes on a rune boundary,
// marking truncation with an ellipsis.
func BoundSummary(text string) string {
	if len(text) <= MaxSummaryLength {
		return text
	}
	cut := MaxSummaryLength
	for cut > 0 && !utf8.RuneStart(text[cut]) {
		cut--
	}
	return text[:cut] + "…"
}
