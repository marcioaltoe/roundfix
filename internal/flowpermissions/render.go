package flowpermissions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const codexHeader = "# Roundfix Flow Command List for Codex. Written by roundfix permissions; do not edit.\n" +
	"# To remove it: roundfix permissions --harness codex --remove\n"

// ClaudeRule renders one allow rule; Guards supplies separate ask rules.
func ClaudeRule(entry Entry) string { return "Bash(" + strings.Join(entry.Words, " ") + " *)" }

// PlainEntries returns the list in plain words for the OpenCode report.
func PlainEntries() []string { return shipped.plainEntries() }

func (l list) plainEntries() []string {
	result := make([]string, 0, len(l.entries))
	for _, entry := range l.entries {
		result = append(result, strings.Join(entry.Words, " "))
	}
	return result
}

// RenderCodexRules renders the complete Codex file, with no Guard Entries.
func RenderCodexRules() []byte { return shipped.renderCodexRules() }

func (l list) renderCodexRules() []byte {
	var out bytes.Buffer
	out.WriteString(codexHeader)
	for _, group := range []string{"roundfix", "git"} {
		var alternatives []string
		for _, entry := range l.entries {
			if len(entry.Words) >= 2 && entry.Words[0] == group {
				alternatives = append(alternatives, strconv.Quote(entry.Words[1]))
			}
		}
		if len(alternatives) > 0 {
			writeCodexBlock(&out, "["+strconv.Quote(group)+", ["+strings.Join(alternatives, ", ")+"]]", codexJustification(group))
		}
	}
	// Normally these are GitHub entries. Retaining an unexpected entry in the
	// rendering lets the invariant checker expose a forbidden addition.
	for _, entry := range l.entries {
		if len(entry.Words) > 0 && (entry.Words[0] == "roundfix" || entry.Words[0] == "git") {
			continue
		}
		words := make([]string, 0, len(entry.Words))
		for _, word := range entry.Words {
			words = append(words, strconv.Quote(word))
		}
		writeCodexBlock(&out, "["+strings.Join(words, ", ")+"]", codexJustification("gh"))
	}
	return out.Bytes()
}

func codexJustification(group string) string {
	switch group {
	case "roundfix":
		return "Roundfix Flow Command List: a roundfix command of the flow"
	case "git":
		return "Roundfix Flow Command List: a Git read the skills name"
	default:
		return "Roundfix Flow Command List: a GitHub read the skills name"
	}
}

func writeCodexBlock(out *bytes.Buffer, pattern, justification string) {
	fmt.Fprintf(out, "prefix_rule(\n    pattern = %s,\n    decision = \"allow\",\n    justification = %s,\n)\n", pattern, strconv.Quote(justification))
}

var codexWord = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ReadCodexRules accepts only the rendered file grammar, including its headers,
// allow decision and group justification. It can read an older shipped list.
func ReadCodexRules(data []byte) (entries []string, ok bool) {
	lines := strings.Split(string(data), "\n")
	lines = slices.DeleteFunc(lines, func(line string) bool { return line == "" })
	headers := strings.Split(strings.TrimSuffix(codexHeader, "\n"), "\n")
	if len(lines) < len(headers) || !slices.Equal(lines[:len(headers)], headers) {
		return nil, false
	}
	lines = lines[len(headers):]
	for i := 0; i < len(lines); {
		if len(lines)-i < 5 || lines[i] != "prefix_rule(" || lines[i+2] != "    decision = \"allow\"," || lines[i+4] != ")" {
			return nil, false
		}
		pattern, found := strings.CutPrefix(lines[i+1], "    pattern = ")
		if !found || !strings.HasSuffix(pattern, ",") {
			return nil, false
		}
		pattern = strings.TrimSuffix(pattern, ",")
		var parts []json.RawMessage
		if json.Unmarshal([]byte(pattern), &parts) != nil || len(parts) < 2 {
			return nil, false
		}
		var group string
		if json.Unmarshal(parts[0], &group) != nil {
			return nil, false
		}
		var blockEntries []string
		var canonical []string
		switch group {
		case "roundfix", "git":
			var alternatives []string
			if len(parts) != 2 || json.Unmarshal(parts[1], &alternatives) != nil || len(alternatives) == 0 {
				return nil, false
			}
			for _, word := range alternatives {
				if !codexWord.MatchString(word) {
					return nil, false
				}
				blockEntries = append(blockEntries, group+" "+word)
				canonical = append(canonical, strconv.Quote(word))
			}
			if pattern != "["+strconv.Quote(group)+", ["+strings.Join(canonical, ", ")+"]]" {
				return nil, false
			}
		case "gh":
			if len(parts) > 3 {
				return nil, false
			}
			words := []string{group}
			canonical = []string{strconv.Quote(group)}
			for _, part := range parts[1:] {
				var word string
				if json.Unmarshal(part, &word) != nil || !codexWord.MatchString(word) {
					return nil, false
				}
				words = append(words, word)
				canonical = append(canonical, strconv.Quote(word))
			}
			if pattern != "["+strings.Join(canonical, ", ")+"]" {
				return nil, false
			}
			blockEntries = append(blockEntries, strings.Join(words, " "))
		default:
			return nil, false
		}
		if lines[i+3] != "    justification = "+strconv.Quote(codexJustification(group))+"," {
			return nil, false
		}
		entries = append(entries, blockEntries...)
		i += 5
	}
	return entries, len(entries) > 0
}
