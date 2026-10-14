package secretscan

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Match is one Secret Match. It never holds the matched value.
type Match struct {
	Path   string
	Line   int
	Kind   string
	Commit string
}

func (m Match) String() string {
	location := fmt.Sprintf("%s:%d: %s", m.Path, m.Line, m.Kind)
	if m.Commit != "" {
		location += " (commit " + m.Commit[:min(12, len(m.Commit))] + ")"
	}
	return location
}

// MatchesError lists every Secret Match without retaining any matched value.
type MatchesError struct {
	Matches []Match
}

func (err *MatchesError) Error() string {
	lines := make([]string, len(err.Matches))
	for index, match := range err.Matches {
		lines[index] = match.String()
	}
	return strings.Join(lines, "\n")
}

// ScanContent reports at most one Secret Match per text line. Content with a
// NUL byte in its first 8000 bytes is binary and contributes no Matches.
func ScanContent(path string, content []byte) []Match {
	if bytes.IndexByte(content[:min(8000, len(content))], 0) >= 0 {
		return nil
	}
	var matches []Match
	for index, line := range bytes.Split(content, []byte{'\n'}) {
		if kind := matchesKind(line); kind != "" {
			matches = append(matches, Match{Path: path, Line: index + 1, Kind: kind})
		}
	}
	return matches
}

// PushedHistoryArgs selects every commit not reachable from the remote's
// tracking refs, rather than the net difference of the branch.
func PushedHistoryArgs(head, remote string) []string {
	return []string{"log", "-p", "-U0", "--reverse", "--no-color", "--no-ext-diff", "--no-renames", "--format=commit %H", head, "--not", "--remotes=" + remote}
}

// ScanPushedHistory reads the patch stream selected by PushedHistoryArgs.
// Only added lines contribute Matches, with their line in that file version.
func ScanPushedHistory(log []byte) []Match {
	var matches []Match
	var commit, path string
	var lineNumber, oldRemaining, newRemaining int
	for _, raw := range bytes.Split(log, []byte{'\n'}) {
		line := string(raw)
		if oldRemaining > 0 || newRemaining > 0 {
			switch {
			case strings.HasPrefix(line, "+"):
				if path != "" && commit != "" {
					if kind := matchesKind(raw[1:]); kind != "" {
						matches = append(matches, Match{Path: path, Line: lineNumber, Kind: kind, Commit: commit})
					}
				}
				lineNumber++
				newRemaining--
			case strings.HasPrefix(line, "-"):
				oldRemaining--
			case strings.HasPrefix(line, " "):
				lineNumber++
				oldRemaining--
				newRemaining--
			}
			continue
		}
		switch {
		case strings.HasPrefix(line, "commit "):
			commit = strings.TrimPrefix(line, "commit ")
			path = ""
		case strings.HasPrefix(line, "diff --git "):
			path = ""
		case strings.HasPrefix(line, "+++ "):
			path = strings.TrimSuffix(strings.TrimPrefix(line, "+++ "), "\t")
			if strings.HasPrefix(path, `"`) {
				decoded, err := strconv.Unquote(path)
				if err != nil {
					path = ""
					continue
				}
				path = decoded
			}
			if strings.HasPrefix(path, "b/") {
				path = strings.TrimPrefix(path, "b/")
			} else {
				path = ""
			}
		case strings.HasPrefix(line, "@@ "):
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[3] == "@@" {
				_, oldRemaining = hunkRange(fields[1], '-')
				lineNumber, newRemaining = hunkRange(fields[2], '+')
			}
		}
	}
	return matches
}

func hunkRange(field string, prefix byte) (start, count int) {
	if len(field) == 0 || field[0] != prefix {
		return 0, 0
	}
	number, length, explicit := strings.Cut(field[1:], ",")
	start, err := strconv.Atoi(number)
	if err != nil {
		return 0, 0
	}
	count = 1
	if explicit {
		count, err = strconv.Atoi(length)
		if err != nil {
			return 0, 0
		}
	}
	return start, count
}
