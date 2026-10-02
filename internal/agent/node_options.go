package agent

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// agentNodeOptions returns the options for an agent process and the absolute
// preload paths removed. An unsplittable value is returned unchanged.
func agentNodeOptions(value string, exists func(string) bool) (kept string, dropped []string, ok bool) {
	words, ok := splitNodeOptions(value)
	if !ok {
		return value, nil, false
	}
	keptWords := make([]string, 0, len(words))
	for index := 0; index < len(words); index++ {
		word := words[index]
		preload := ""
		paired := false
		switch {
		case word == "--require" || word == "-r" || word == "--import":
			if index+1 < len(words) {
				preload = words[index+1]
				paired = true
			}
		case strings.HasPrefix(word, "--require="):
			preload = strings.TrimPrefix(word, "--require=")
		case strings.HasPrefix(word, "--import="):
			preload = strings.TrimPrefix(word, "--import=")
		}
		path := absoluteNodePreloadPath(preload)
		if path != "" && !exists(path) {
			dropped = append(dropped, path)
			if paired {
				index++
			}
			continue
		}
		keptWords = append(keptWords, quoteNodeOption(word))
		if paired {
			index++
			keptWords = append(keptWords, quoteNodeOption(words[index]))
		}
	}
	return strings.Join(keptWords, " "), dropped, true
}

// Node options use double quotes, with backslash escaping inside quotes.
// Single quotes and backslashes outside quotes are literal characters.
func splitNodeOptions(value string) ([]string, bool) {
	var words []string
	var word strings.Builder
	quoted, started := false, false
	for index := 0; index < len(value); index++ {
		char := value[index]
		switch {
		case char == '"':
			quoted = !quoted
			started = true
		case quoted && char == '\\':
			index++
			if index == len(value) {
				return nil, false
			}
			word.WriteByte(value[index])
		case !quoted && char == ' ':
			if started {
				words = append(words, word.String())
				word.Reset()
				started = false
			}
		default:
			word.WriteByte(char)
			started = true
		}
	}
	if quoted {
		return nil, false
	}
	if started {
		words = append(words, word.String())
	}
	return words, true
}

func quoteNodeOption(word string) string {
	if word != "" && !strings.ContainsAny(word, " \"") {
		return word
	}
	word = strings.ReplaceAll(word, "\\", "\\\\")
	word = strings.ReplaceAll(word, "\"", "\\\"")
	return "\"" + word + "\""
}

func absoluteNodePreloadPath(value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	if strings.HasPrefix(value, "file://") {
		parsed, err := url.Parse(value)
		if err == nil && filepath.IsAbs(parsed.Path) {
			return parsed.Path
		}
	}
	return ""
}

// Notices are deduplicated across runners in this Roundfix process. The lock
// also serializes writes when concurrent probes share a notice writer.
var nodePreloadNotices = struct {
	sync.Mutex
	paths map[string]struct{}
}{paths: make(map[string]struct{})}

func (runner *ACPXRunner) noticeNodePreloads(paths []string) {
	if len(paths) == 0 {
		return
	}
	writer := runner.Notices
	if writer == nil {
		writer = os.Stderr
	}
	nodePreloadNotices.Lock()
	defer nodePreloadNotices.Unlock()
	for _, path := range paths {
		if _, seen := nodePreloadNotices.paths[path]; seen {
			continue
		}
		writeNodePreloadNotice(writer, path)
		nodePreloadNotices.paths[path] = struct{}{}
	}
}

func writeNodePreloadNotice(writer io.Writer, path string) {
	if _, err := fmt.Fprintf(writer, "roundfix: notice: NODE_OPTIONS preload %q does not exist; Roundfix left it out of the agent environment\n", path); err != nil {
		// A custom diagnostic sink must not prevent the agent from starting.
		log.Printf("write NODE_OPTIONS preload notice: %v", err)
	}
}

// ACPXEnvironment returns environment with every NODE_OPTIONS preload whose
// file is missing left out, writing one notice per dropped path to notices.
// Processes Roundfix starts outside an ACPXRunner, such as the acpx and node
// version probes, use it so a dead preload cannot stop them either.
func ACPXEnvironment(environment []string, notices io.Writer) []string {
	if environment == nil {
		environment = []string{}
	}
	runner := &ACPXRunner{Environment: environment, Notices: notices}
	return runner.baseEnv()
}
