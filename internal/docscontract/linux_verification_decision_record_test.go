//go:build docscontract

//verify:always

package docscontract

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestTheLinuxVerificationDecisionSupersedesTheLaneRecords(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if err != nil {
		t.Fatal(err)
	}
	const title = "# Roundfix verifies on macOS, and the Delivery Queue has no Linux Lane"
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
	retired := []string{"0288", "0294", "0297", "0303", "0316"}
	t.Run("decision_record", func(t *testing.T) {
		h := spec.ReadADRHeader([]byte(body))
		if !h.Frontmatter || h.Status != "accepted" || h.Knowledge == nil || !*h.Knowledge {
			t.Errorf("decision header: %+v", h)
		}
		if !regexp.MustCompile(`Supersedes ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 \(\d{4}-\d{2}-\d{2}\); all are retired, and Git keeps their text\.`).MatchString(body) {
			t.Error("missing dated supersession sentence")
		}
		for _, d := range []string{"2026-10-18", "2026-10-19", "2026-10-20", "2026-10-22"} {
			if !strings.Contains(body, d) {
				t.Errorf("missing original creation date %s", d)
			}
		}
	})
	t.Run("superseded_records_are_deleted", func(t *testing.T) {
		for _, e := range entries {
			for _, n := range retired {
				if strings.HasPrefix(e.Name(), n+"-") {
					t.Errorf("retired record exists: %s", e.Name())
				}
			}
		}
	})
	t.Run("amended_records", func(t *testing.T) {
		for _, n := range []string{"0252", "0285", "0295", "0384", "0385", "0423"} {
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
		for _, term := range []string{"Linux Lane", "Lane Environment Skip", "Lane Container Failure", "Lane Change Failure"} {
			if strings.Contains(s, "**"+term+"**:") {
				t.Errorf("still defined: %s", term)
			}
		}
		for _, term := range []string{"Base Test Tools", "External Test Tool"} {
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
	t.Run("export_list", func(t *testing.T) {
		b, err := os.ReadFile(filepath.Join(root, ".secondbrain-export"))
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(string(b), "\n")
		path := "docs/adr/" + name
		count := 0
		var adrs []string
		for _, line := range lines {
			if line == path {
				count++
			}
			if strings.HasPrefix(line, "docs/adr/") {
				adrs = append(adrs, line)
				for _, n := range retired {
					if strings.HasPrefix(line, "docs/adr/"+n+"-") {
						t.Errorf("retired export: %s", line)
					}
				}
			}
		}
		if count != 1 || !slices.IsSorted(adrs) {
			t.Errorf("decision count=%d, sorted=%v", count, slices.IsSorted(adrs))
		}
	})
}

func TestNoLiveDocumentDescribesTheLinuxLane(t *testing.T) {
	t.Parallel()
	root := characterizationRepositoryRoot(t)
	check := func(path string) {
		f, err := os.Open(path)
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		line := 0
		for scanner.Scan() {
			line++
			s := scanner.Text()
			if strings.Contains(s, "is deprecated and ignored; the Linux Lane was retired, so remove it") {
				continue
			}
			for _, term := range []string{"Linux Lane", "deliver lane", "Linux gate", "linux-lane", "DR-LANE-"} {
				if strings.Contains(s, term) {
					rel, _ := filepath.Rel(root, path)
					t.Errorf("%s:%d: contains %q", filepath.ToSlash(rel), line, term)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			t.Error(err)
		}
	}
	for _, dir := range []string{"docs/user-guide", ".agents/skills/roundfix", "skills/roundfix", "docs/agents", ".github/workflows"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (dir == ".github/workflows" || strings.HasSuffix(path, ".md")) {
				check(path)
			}
			return nil
		})
		if err != nil {
			t.Error(fmt.Errorf("scan %s: %w", dir, err))
		}
	}
	for _, name := range []string{"AGENTS.md", "CLAUDE.md", "README.md", "Makefile", ".roundfixrc.yml"} {
		check(filepath.Join(root, name))
	}
}
