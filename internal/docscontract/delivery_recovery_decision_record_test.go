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

func TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if err != nil {
		t.Fatal(err)
	}
	const title = "# A Delivery Retry ends with a live Queue Owner, and status reports a queue nobody advances"
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
		t.Fatal("no decision record has the delivery recovery title")
	}
	header := spec.ReadADRHeader(record)
	if !header.Frontmatter || header.Status != "accepted" || header.Knowledge == nil || !*header.Knowledge {
		t.Fatalf("delivery recovery decision header = %+v, want accepted with true Knowledge Mark", header)
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
	for _, priorName := range []string{
		"0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md",
		"0395-a-queue-item-stops-only-for-its-own-work.md",
	} {
		prior, err := os.ReadFile(filepath.Join(root, "docs/adr", priorName))
		if err != nil {
			t.Fatal(err)
		}
		priorHeader := spec.ReadADRHeader(prior)
		if priorHeader.Status != "accepted" || !slices.Contains(priorHeader.AmendedBy, identifier) {
			t.Errorf("%s header = %+v", priorName, priorHeader)
		}
		noticeLead := "**Amended by [" + identifier + "](" + filepath.Base(recordPath) + ") (" + created.Format("2006-01-02") + ").**"
		anchor := "A queue with no recorded owner adds no line."
		if strings.HasPrefix(priorName, "0395-") {
			anchor = "**An owner that stops with an item left keeps its record.**"
			start := strings.Index(string(prior), anchor)
			if start < 0 {
				t.Fatal("owner paragraph missing")
			}
			paragraph, _, _ := strings.Cut(string(prior)[start:], "\n\n")
			anchor = paragraph
		}
		if !strings.Contains(string(prior), anchor+"\n\n"+noticeLead) {
			t.Errorf("%s must place its dated amendment directly after the changed paragraph", priorName)
		}
		if !strings.Contains(string(prior), "updated_at: "+metadata.CreatedAt+"\n") {
			t.Errorf("%s updated_at differs from decision created_at", priorName)
		}
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
				t.Error("delivery recovery decision path is outside the Standard Export List")
			}
		}
		if i > start && i < stop && strings.HasPrefix(line, "docs/adr/") {
			adrPaths = append(adrPaths, line)
		}
	}
	if count != 1 {
		t.Errorf("delivery recovery decision path occurs %d times in the export, want once", count)
	}
	if !slices.IsSorted(adrPaths) {
		t.Errorf("Standard Export List ADR paths are out of order: %v", adrPaths)
	}
}
