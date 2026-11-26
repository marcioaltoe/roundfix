//go:build docscontract

//verify:always

package docscontract

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestTheTimeRulesRetirementDecisionRecord(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if err != nil {
		t.Fatal(err)
	}
	const title = "# Two time rules are retired until the Gate Time Rule is written"
	var name, body string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, "docs/adr", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if slices.Contains(strings.Split(string(b), "\n"), title) {
			if name != "" {
				t.Fatal("duplicate decision title")
			}
			name, body = e.Name(), string(b)
		}
	}
	if name == "" {
		t.Fatal("decision title not found")
	}
	number, _, ok := strings.Cut(name, "-")
	if !ok || !regexp.MustCompile(`^\d{4}$`).MatchString(number) {
		t.Fatalf("bad record name: %s", name)
	}
	id := "ADR-" + number
	t.Run("decision_record", func(t *testing.T) {
		h := spec.ReadADRHeader([]byte(body))
		if !h.Frontmatter || h.Status != "accepted" || h.Knowledge == nil || *h.Knowledge {
			t.Errorf("decision header: %+v", h)
		}
		if !strings.Contains(body, "Supersedes ADR-0387 (2026-11-10); it is retired, and Git keeps its text.") {
			t.Error("missing dated supersession sentence")
		}
	})
	t.Run("superseded_record_is_deleted", func(t *testing.T) {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "0387-") {
				t.Errorf("retired record exists: %s", e.Name())
			}
		}
	})
	t.Run("amended_records", func(t *testing.T) {
		for _, n := range []string{"0452", "0253", "0252", "0259", "0257", "0182"} {
			matches, err := filepath.Glob(filepath.Join(root, "docs/adr", n+"-*.md"))
			if err != nil || len(matches) != 1 {
				t.Errorf("ADR-%s: %v (%v)", n, matches, err)
				continue
			}
			b, err := os.ReadFile(matches[0])
			if err != nil {
				t.Error(err)
				continue
			}
			h := spec.ReadADRHeader(b)
			if h.Status != "accepted" || !slices.Contains(h.AmendedBy, id) {
				t.Errorf("ADR-%s header: %+v", n, h)
			}
			marker := "**Amended by [" + id + "](" + name + ") ("
			if strings.Count(string(b), marker) != 1 {
				t.Errorf("ADR-%s requires one amendment notice", n)
			}
		}
	})
	t.Run("glossary", func(t *testing.T) {
		b, err := os.ReadFile(filepath.Join(root, "CONTEXT.md"))
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		for _, term := range []string{"Gate Time Target", "Gate Time Report", "Gate Time Rule", "Complete Suite Run", "Measured Dossier", "Release Provenance", "Full Contract Run"} {
			start := strings.Index(s, "**"+term+"**:\n")
			if start < 0 {
				t.Errorf("missing %s", term)
				continue
			}
			part := strings.SplitN(s[start:], "\n\n", 2)[0]
			if !strings.Contains(part, name) {
				t.Errorf("%s lacks decision link", term)
			}
		}
	})
}
