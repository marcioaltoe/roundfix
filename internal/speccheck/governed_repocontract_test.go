//go:build repocontract

// Suite: governed-path authorization history contract.
// Invariant: every path bounded by an authorization record remains governed.
// Boundary IN: authorization records and the public governed-path predicate.
// Boundary OUT: changed-path audit integration and sanctioned regeneration.
package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const authorizationRecordsDir = "docs/workflow/authorizations"

func TestEveryBoundedPathIsGoverned(t *testing.T) {
	t.Run("repository records are governed", func(t *testing.T) {
		repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatalf("resolve repository root: %v", err)
		}

		findings, recordsExist, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if !recordsExist {
			t.Skipf("no authorization records under %s", authorizationRecordsDir)
		}
		if len(findings) != 0 {
			t.Fatalf("bounded-path contract failed:\n%s", strings.Join(findings, "\n"))
		}
	})

	t.Run("unmatched path names path and record", func(t *testing.T) {
		repoRoot := t.TempDir()
		const record = "docs/workflow/authorizations/synthetic.md"
		writeGovernedAuthorizationRecord(t, repoRoot, record, `---
granted: 2026-08-15
action: synthetic contract probe
paths:
  - README.md
consuming: synthetic-spec
---
`)

		findings, recordsExist, err := auditBoundedPathsAreGoverned(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
		if !recordsExist {
			t.Fatal("authorization record was not read")
		}
		if len(findings) != 1 || !strings.Contains(findings[0], "README.md") || !strings.Contains(findings[0], record) {
			t.Fatalf("findings = %q, want one finding naming path %q and record %q", findings, "README.md", record)
		}
	})

	t.Run("no authorization records skips", func(t *testing.T) {
		findings, recordsExist, err := auditBoundedPathsAreGoverned(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if recordsExist || len(findings) != 0 {
			t.Fatalf("recordsExist = %t, findings = %q, want no records and no findings", recordsExist, findings)
		}
		t.Skipf("no authorization records under %s", authorizationRecordsDir)
	})
}

func auditBoundedPathsAreGoverned(repoRoot string) ([]string, bool, error) {
	directory := filepath.Join(repoRoot, filepath.FromSlash(authorizationRecordsDir))
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read authorization records %q: %w", directory, err)
	}

	var findings []string
	recordsExist := false
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		recordsExist = true
		record := filepath.ToSlash(filepath.Join(authorizationRecordsDir, entry.Name()))
		content, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			return nil, true, fmt.Errorf("read authorization record %q: %w", record, err)
		}

		bounded := parseMechanicalAuthorizationPaths(content)
		paths := make([]string, 0, len(bounded))
		for path := range bounded {
			// Legacy prose can cite a regeneration command in the same bullet as
			// a bounded file. Commands are not repository-relative paths.
			if strings.ContainsAny(path, " \t") {
				continue
			}
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			if !GovernedPath(path) {
				findings = append(findings, fmt.Sprintf("%s bounds %s, which is not governed", record, path))
			}
		}
	}
	return findings, recordsExist, nil
}

func writeGovernedAuthorizationRecord(t *testing.T, repoRoot, relativePath, content string) {
	t.Helper()

	path := filepath.Join(repoRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
