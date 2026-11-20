//go:build docscontract

//verify:always

package docscontract

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

func TestTheIndexusDecisionAmendsTheConsultationRule(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if err != nil {
		t.Fatal(err)
	}
	const title = "# The Secondbrain search and the duplicate check go through the Indexus, with no second path"
	var recordPath string
	var record []byte
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join("docs/adr", entry.Name())
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(content), "\n") {
			if line != title {
				continue
			}
			if recordPath != "" {
				t.Fatalf("decision title occurs in both %s and %s", recordPath, path)
			}
			recordPath, record = filepath.ToSlash(path), content
		}
	}
	if recordPath == "" {
		t.Fatal("no decision record has the Indexus consultation title")
	}
	header := spec.ReadADRHeader(record)
	if !header.Frontmatter || header.Status != "accepted" || header.Knowledge == nil || !*header.Knowledge {
		t.Fatalf("Indexus decision header = %+v, want accepted with true Knowledge Mark", header)
	}
	number, _, ok := strings.Cut(filepath.Base(recordPath), "-")
	if !ok {
		t.Fatalf("decision record has no numeric identifier: %s", recordPath)
	}
	identifier := "ADR-" + number
	var metadata struct {
		CreatedAt string `yaml:"created_at"`
	}
	parts := strings.SplitN(string(record), "---", 3)
	if len(parts) != 3 {
		t.Fatal("decision record lacks lifecycle frontmatter")
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		t.Fatal(err)
	}
	created, err := time.Parse(time.RFC3339, metadata.CreatedAt)
	if err != nil {
		t.Fatalf("decision created_at: %v", err)
	}
	prior, err := os.ReadFile(filepath.Join(root, "docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md"))
	if err != nil {
		t.Fatal(err)
	}
	amendments := spec.ReadADRHeader(prior).AmendedBy
	if !slices.Contains(amendments, identifier) {
		t.Errorf("ADR-0330 amended_by = %v, missing %s", amendments, identifier)
	}
	// This historical paragraph must stay byte-identical; the amendment follows it.
	const firstParagraph = "**The Secondbrain stays the business knowledge base, and it is read-only.**\nThe Baseline keeps its consultation rules: read `wiki/index.md` first, then\n`qmd`, cite every file used, keep secrets out, and ask its owner to ingest a\ndurable change. A session never writes a Secondbrain file, `inbox/` included.\nSubstantive external research reaches the Secondbrain as a digest message to\nthe project that owns it, which ingests it by its own contract."
	noticeLead := "**Amended by [" + identifier + "](" + filepath.Base(recordPath) + ") (" + created.Format("2006-01-02") + ").**"
	decisionStart := "## Decision\n\n" + firstParagraph + "\n\n" + noticeLead
	if !strings.Contains(string(prior), decisionStart) {
		t.Error("ADR-0330 must preserve its first Decision paragraph and follow it immediately with the linked, dated amendment")
	}
	for _, paragraph := range strings.Split(string(prior), "\n\n") {
		if !strings.HasPrefix(paragraph, noticeLead) {
			continue
		}
		text := strings.Join(strings.Fields(paragraph), " ")
		for _, phrase := range []string{
			"after the wiki index uses the Indexus tool `pesquisar_acervos`",
			"the advisory duplicate check uses the same tool",
			"no second search tool is named",
		} {
			if !strings.Contains(text, phrase) {
				t.Errorf("amendment notice lacks %q", phrase)
			}
		}
	}
	if !strings.Contains(string(prior), "updated_at: "+metadata.CreatedAt+"\n") {
		t.Error("ADR-0330 updated_at must equal the Indexus decision's created_at")
	}
	export, err := os.ReadFile(filepath.Join(root, ".secondbrain-export"))
	if err != nil {
		t.Fatal(err)
	}
	const begin = "# roundfix:standard-export:begin"
	const end = "# roundfix:standard-export:end"
	lines := strings.Split(string(export), "\n")
	if strings.Count(string(export), begin) != 1 || strings.Count(string(export), end) != 1 {
		t.Fatal("Standard Export List must have one pair of markers")
	}
	start, stop := slices.Index(lines, begin), slices.Index(lines, end)
	if start < 0 || stop <= start {
		t.Fatal("Standard Export List markers are out of order")
	}
	var adrPaths []string
	count := 0
	for i, line := range lines {
		if line == recordPath {
			count++
			if i <= start || i >= stop {
				t.Error("Indexus decision path is outside the Standard Export List")
			}
		}
		if i > start && i < stop && strings.HasPrefix(line, "docs/adr/") {
			adrPaths = append(adrPaths, line)
		}
	}
	if count != 1 {
		t.Errorf("Indexus decision path occurs %d times in the export, want once", count)
	}
	if !slices.IsSorted(adrPaths) {
		t.Errorf("Standard Export List ADR paths are out of order: %v", adrPaths)
	}
}
