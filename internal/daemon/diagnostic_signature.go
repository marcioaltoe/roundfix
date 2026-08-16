package daemon

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
)

var diagnosticNormalizers = []struct {
	pattern     *regexp.Regexp
	replacement []byte
}{
	{
		pattern:     regexp.MustCompile(`\b[0-9]{4}-[0-9]{2}-[0-9]{2}[T ][0-9]{2}:[0-9]{2}:[0-9]{2}([.,][0-9]+)?(Z|[+-][0-9]{2}:?[0-9]{2})?\b`),
		replacement: []byte("<timestamp>"),
	},
	{
		pattern:     regexp.MustCompile(`(?i)(/private)?/tmp/[^[:space:]"'<>]+|/var/tmp/[^[:space:]"'<>]+|(/private)?/var/folders/[^[:space:]"'<>]+/T/[^[:space:]"'<>]+`),
		replacement: []byte("<temporary-path>"),
	},
	{
		pattern:     regexp.MustCompile(`(?i)[A-Z]:\\[^[:space:]"'<>]*\\(AppData\\Local\\Temp|Windows\\Temp|Temp)\\[^[:space:]"'<>]+`),
		replacement: []byte("<temporary-path>"),
	},
	{
		pattern:     regexp.MustCompile(`\b([0-9]+([.][0-9]+)?(ns|us|µs|ms|s|m|h))+\b`),
		replacement: []byte("<duration>"),
	},
	{
		pattern:     regexp.MustCompile(`\b[0-9]{1,2}:[0-9]{2}:[0-9]{2}([.,][0-9]+)?\b`),
		replacement: []byte("<duration>"),
	},
	{
		pattern:     regexp.MustCompile(`(?i)\b(pid|process([ _-]?id)?)[[:space:]:=#-]*[0-9]+\b`),
		replacement: []byte("<process-id>"),
	},
	{
		pattern:     regexp.MustCompile(`\brun_[0-9]{8}T[0-9]{6}Z_[0-9A-Za-z]+\b`),
		replacement: []byte("<run-id>"),
	},
	{
		pattern:     regexp.MustCompile(`(?i)\brun[ _-]?id[[:space:]:=#]+[A-Za-z0-9][A-Za-z0-9._-]*`),
		replacement: []byte("<run-id>"),
	},
}

// DiagnosticSignature returns a command-scoped hash of a diagnostic after
// replacing the volatile spans named by ADR-0136.
func DiagnosticSignature(command string, diagnostic []byte) string {
	normalized := diagnostic
	for _, normalizer := range diagnosticNormalizers {
		normalized = normalizer.pattern.ReplaceAll(normalized, normalizer.replacement)
	}

	payload := strconv.AppendInt(nil, int64(len(command)), 10)
	payload = append(payload, ':')
	payload = append(payload, command...)
	payload = append(payload, normalized...)
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
