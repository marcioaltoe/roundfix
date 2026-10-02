// Package runcause classifies recorded Verification failures and corrective
// Tasks with a fixed, ordered signature table. It performs no writes.
package runcause

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
)

//go:embed signatures.json
var signatureBytes []byte

type signature struct {
	ID      string `json:"id"`
	Class   string `json:"class"`
	Source  string `json:"source"`
	Pattern string `json:"pattern"`
	re      *regexp.Regexp
}

type trigger struct {
	ID      string `json:"id"`
	Pattern string `json:"pattern"`
	re      *regexp.Regexp
}

// Table is the compiled, embedded classification contract and its byte digest.
type Table struct {
	SHA256              string      `json:"-"`
	Schema              string      `json:"schema"`
	Classes             []string    `json:"classes"`
	RepositoryKnowledge []string    `json:"repository_knowledge"`
	DiagnosticTailBytes int64       `json:"diagnostic_tail_bytes"`
	TaskTextBytes       int         `json:"task_text_bytes"`
	Triggers            []trigger   `json:"triggers"`
	Signatures          []signature `json:"signatures"`
}

// Load validates and compiles the embedded table without reordering it.
func Load() (Table, error) { return load(signatureBytes) }

func load(data []byte) (Table, error) {
	var t Table
	if err := json.Unmarshal(data, &t); err != nil {
		return t, fmt.Errorf("decode cause signatures: %w", err)
	}
	if t.Schema != "roundfix/cause-signatures/v1" || t.DiagnosticTailBytes <= 0 || t.TaskTextBytes <= 0 {
		return t, fmt.Errorf("invalid cause signature schema or text limits")
	}
	for _, class := range t.RepositoryKnowledge {
		if !slices.Contains(t.Classes, class) {
			return t, fmt.Errorf("unknown repository-knowledge class %q", class)
		}
	}
	ids := map[string]bool{}
	compile := func(id, pattern string) (*regexp.Regexp, error) {
		if id == "" || ids[id] {
			return nil, fmt.Errorf("duplicate or empty signature/trigger id %q", id)
		}
		ids[id] = true
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("compile cause pattern %q: %w", id, err)
		}
		return re, nil
	}
	for i := range t.Signatures {
		s := &t.Signatures[i]
		if !slices.Contains(t.Classes, s.Class) {
			return t, fmt.Errorf("unknown cause class %q", s.Class)
		}
		if !slices.Contains([]string{"diagnostic", "command", "task", "any"}, s.Source) {
			return t, fmt.Errorf("unknown cause source %q", s.Source)
		}
		var err error
		s.re, err = compile(s.ID, s.Pattern)
		if err != nil {
			return t, err
		}
	}
	for i := range t.Triggers {
		var err error
		t.Triggers[i].re, err = compile(t.Triggers[i].ID, t.Triggers[i].Pattern)
		if err != nil {
			return t, err
		}
	}
	t.SHA256 = fmt.Sprintf("%x", sha256.Sum256(data))
	return t, nil
}
