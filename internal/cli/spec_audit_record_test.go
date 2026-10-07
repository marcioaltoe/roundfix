package cli

import (
	"os"
	"path/filepath"
	"roundfix/internal/spec"
	"strings"
	"testing"
)

func TestSpecAuditAcceptsAnArchiveRecordSlug(t *testing.T) {
	for _, builtIn := range []bool{true, false} {
		name := "custom root"
		if builtIn {
			name = "built-in root"
		}
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "docs", "specs")
			const slug = "0001-archived"
			archive := spec.ArchiveSpecRoot(root, builtIn)
			if err := os.MkdirAll(archive, 0o755); err != nil {
				t.Fatal(err)
			}
			content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{Spec: slug, Title: "Fixture", Created: "2026-10-06", Archived: "2026-10-06", Disposition: spec.ArchivePass, Source: "docs/specs/" + slug, SourceRevision: strings.Repeat("a", 40), Outcome: "Archived."})
			if err != nil {
				t.Fatal(err)
			}
			path := spec.ArchiveRecordPath(archive, slug)
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := validateSpecAuditSlug(root, builtIn, slug); err != nil {
				t.Fatal(err)
			}
			if err := validateSpecAuditSlug(root, builtIn, "0002-missing"); err == nil {
				t.Fatal("missing slug accepted")
			}
			if err := os.WriteFile(path, []byte("bad record"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := validateSpecAuditSlug(root, builtIn, slug); err == nil {
				t.Fatal("malformed record accepted")
			}
		})
	}
}
