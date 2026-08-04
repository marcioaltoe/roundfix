package spec

// Suite: Spec archive eligibility corpus
// Invariant: pass Specs remain eligible and QA overrides retain failed evidence.
// Boundary IN: archived Task status, QA metadata, and the archive QA precondition.
// Boundary OUT: command I/O and filesystem movement in internal/cli/archive_test.go.

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestArchivedPassCorpusRemainsArchiveEligible(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	pattern := filepath.Join(filepath.Dir(testFile), "..", "..", "docs", "specs", "_archived", "*", "qa", "qa-report-*.md")
	reportPaths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived QA Reports: %v", err)
	}
	if len(reportPaths) == 0 {
		t.Fatal("archived QA Report corpus is empty")
	}

	seen := make(map[string]struct{})
	passSpecs := 0
	for _, reportPath := range reportPaths {
		specDir := filepath.Dir(filepath.Dir(reportPath))
		if _, ok := seen[specDir]; ok {
			continue
		}
		seen[specDir] = struct{}{}

		report, err := ReadQAReport(specDir)
		if err != nil {
			t.Fatalf("ReadQAReport(%q): %v", specDir, err)
		}
		if report.Verdict != VerdictPass {
			continue
		}
		passSpecs++

		taskPaths, err := filepath.Glob(filepath.Join(specDir, "task_*.md"))
		if err != nil {
			t.Fatalf("find archived Tasks for %q: %v", specDir, err)
		}
		if len(taskPaths) == 0 {
			t.Fatalf("archived pass Spec %q has no Tasks", specDir)
		}
		for _, taskPath := range taskPaths {
			content, err := os.ReadFile(taskPath)
			if err != nil {
				t.Fatalf("read archived Task %q: %v", taskPath, err)
			}
			frontmatterBytes, _, err := splitFrontmatter(content)
			if err != nil {
				t.Fatalf("parse archived Task %q: %v", taskPath, err)
			}
			var frontmatter struct {
				Status string `yaml:"status"`
			}
			if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
				t.Fatalf("parse archived Task %q frontmatter: %v", taskPath, err)
			}
			if status := NormalizeStatus(frontmatter.Status); status != string(StatusCompleted) {
				t.Fatalf("archived pass Task %q status is %q; expected %q", taskPath, status, StatusCompleted)
			}
		}

		unproven, err := archiveUnprovenActions(specDir, report)
		if err != nil {
			t.Fatalf("archive precondition rejected archived pass Spec %q: %v", specDir, err)
		}
		if len(unproven) != 0 {
			t.Fatalf("archive precondition added unproven actions to pass Spec %q: %q", specDir, unproven)
		}
	}
	if passSpecs == 0 {
		t.Fatal("archived corpus has no pass Specs")
	}
	t.Logf("checked archive eligibility for %d archived pass Specs", passSpecs)
}

func TestArchivedQAOverrideCorpusIncludesFailedSpec(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	pattern := filepath.Join(filepath.Dir(testFile), "..", "..", "docs", "specs", "_archived", "*", "_prd.md")
	prdPaths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived Spec PRDs: %v", err)
	}
	for _, prdPath := range prdPaths {
		content, err := os.ReadFile(prdPath)
		if err != nil {
			t.Fatalf("read archived Spec PRD %q: %v", prdPath, err)
		}
		frontmatterBytes, _, err := splitFrontmatter(content)
		if err != nil {
			t.Fatalf("parse archived Spec PRD %q: %v", prdPath, err)
		}
		var frontmatter struct {
			Status     string `yaml:"status"`
			QAOverride bool   `yaml:"qa_override"`
		}
		if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
			t.Fatalf("parse archived Spec PRD %q frontmatter: %v", prdPath, err)
		}
		if !frontmatter.QAOverride || frontmatter.Status != "archived" {
			continue
		}
		report, err := ReadQAReport(filepath.Dir(prdPath))
		if errors.Is(err, ErrNoQAReport) {
			continue
		}
		if err != nil {
			t.Fatalf("read QA override Spec report for %q: %v", prdPath, err)
		}
		if report.Verdict == VerdictFail {
			t.Logf("found archived failed QA override Spec %s", filepath.Base(filepath.Dir(prdPath)))
			return
		}
	}
	t.Fatal("archived corpus has no qa_override Spec whose newest QA verdict is fail")
}
