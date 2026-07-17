package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

const roundfixSessionPrefix = "roundfix-"

// RoundfixSession is a discovered acpx Agent Session owned by Roundfix.
type RoundfixSession struct {
	Name   string
	RunID  string
	TaskID string
}

func SessionRefForQA(runID string, workDir string) SessionRef {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return SessionRef{}
	}
	return SessionRef{Name: "roundfix-" + runID + "-qa", WorkDir: strings.TrimSpace(workDir)}
}

func SessionRefForReview(runID string, batchNumber int, workDir string) SessionRef {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return SessionRef{}
	}
	if batchNumber > 0 {
		return SessionRef{Name: fmt.Sprintf("roundfix-%s-review-%03d", runID, batchNumber), WorkDir: strings.TrimSpace(workDir)}
	}
	return SessionRef{Name: "roundfix-" + runID + "-review", WorkDir: strings.TrimSpace(workDir)}
}

// ParseRoundfixSessions extracts roundfix-named Agent Sessions from acpx
// sessions list output. It accepts JSON-shaped output and human table output,
// but only returns names with the roundfix- prefix.
func ParseRoundfixSessions(output string) []RoundfixSession {
	var sessions []RoundfixSession
	seen := map[string]struct{}{}
	add := func(name string) {
		session, ok := ParseRoundfixSessionName(name)
		if !ok {
			return
		}
		if _, exists := seen[session.Name]; exists {
			return
		}
		seen[session.Name] = struct{}{}
		sessions = append(sessions, session)
	}

	raw := strings.TrimSpace(output)
	if raw == "" {
		return nil
	}
	var decoded any
	if json.Unmarshal([]byte(raw), &decoded) == nil {
		collectRoundfixSessionNames(decoded, add, true)
		return sessions
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		for _, token := range sessionNameTokens(scanner.Text()) {
			add(token)
		}
	}
	return sessions
}

func ParseRoundfixSessionName(name string) (RoundfixSession, bool) {
	name = cleanSessionNameToken(name)
	rest, ok := strings.CutPrefix(name, roundfixSessionPrefix)
	if !ok || rest == "" {
		return RoundfixSession{}, false
	}
	session := RoundfixSession{Name: name}
	if index := strings.LastIndex(rest, "-task_"); index >= 0 {
		runID := strings.TrimSpace(rest[:index])
		taskID := strings.TrimSpace(rest[index+1:])
		if runID == "" || taskID == "" {
			return RoundfixSession{}, false
		}
		session.RunID = runID
		session.TaskID = taskID
		return session, true
	}
	switch {
	case strings.HasSuffix(rest, "-qa"):
		rest = strings.TrimSuffix(rest, "-qa")
	case strings.HasSuffix(rest, "-review"):
		rest = strings.TrimSuffix(rest, "-review")
	default:
		if index := strings.LastIndex(rest, "-review-"); index >= 0 {
			suffix := rest[index+len("-review-"):]
			if len(suffix) == 3 && allDigits(suffix) {
				rest = rest[:index]
			}
		}
	}
	if rest == "" {
		return RoundfixSession{}, false
	}
	session.RunID = rest
	return session, true
}

func allDigits(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return value != ""
}

func collectRoundfixSessionNames(value any, add func(string), allowString bool) {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			collectRoundfixSessionNames(item, add, true)
		}
	case map[string]any:
		for key, item := range typed {
			normalized := strings.ToLower(strings.ReplaceAll(key, "_", ""))
			if text, ok := item.(string); ok && (normalized == "name" || normalized == "session" || normalized == "sessionname" || normalized == "id") {
				add(text)
				continue
			}
			collectRoundfixSessionNames(item, add, normalized == "sessions")
		}
	case string:
		if allowString {
			add(typed)
		}
	}
}

func sessionNameTokens(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	fields := strings.Fields(line)
	if strings.Contains(line, "|") {
		fields = strings.Split(line, "|")
	}
	tokens := []string{}
	for _, field := range fields {
		token := cleanSessionNameToken(field)
		if strings.HasPrefix(token, roundfixSessionPrefix) {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func cleanSessionNameToken(value string) string {
	return strings.Trim(strings.TrimSpace(value), " \t\r\n|,;:\"'`")
}
