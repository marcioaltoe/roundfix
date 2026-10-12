package spec

// Suite: archive enumeration identity, including malformed record declarations.
// Boundary IN: names and temporary filesystem archives.
// Boundary OUT: archive writes and network.

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestOnlyANamedSchemaFileIsAnArchiveRecord(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, content string
		named, record bool
	}{
		{"0256-example.md", "---\nschema: " + ArchiveRecordSchema + "\n---\n", true, true},
		{"0256-malformed.md", "---\nschema: " + ArchiveRecordSchema + "\nbroken: [\n---\n", true, true},
		{"0256-unclosed.md", "---\nschema: " + ArchiveRecordSchema + "\n", true, true},
		{"0256-quoted.md", "---\nschema: '" + ArchiveRecordSchema + "' # record\n---\n", true, true},
		{"0256-key.md", "---\n\"schema\": '" + ArchiveRecordSchema + "'\n---\n", true, true},
		{"0256-indented.md", "---\n  schema: " + ArchiveRecordSchema + "\n---\n", true, true},
		{"README.md", "---\nschema: " + ArchiveRecordSchema + "\n---\n", false, false},
		{"025-example.md", "", false, false},
		{"0256-.md", "", false, false},
		{"0256-dir/file.md", "", false, false},
		{"0256-dir\\file.md", "", false, false},
		{"0256-prose.md", "# Index\nschema: " + ArchiveRecordSchema + "\n", true, false},
		{"0256-body.md", "---\nother: schema\n---\nschema: " + ArchiveRecordSchema + "\n", true, false},
		{"0256-other.md", "---\nschema: other/v1\n---\n", true, false},
		{"0256-nested.md", "---\nother:\n  schema: " + ArchiveRecordSchema + "\n---\n", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsArchiveRecordName(tc.name); got != tc.named {
				t.Fatalf("name=%v want=%v", got, tc.named)
			}
			if got := IsArchiveRecord(tc.name, []byte(tc.content)); got != tc.record {
				t.Fatalf("record=%v want=%v", got, tc.record)
			}
		})
	}
}

func TestArchivedSpecSlugsIgnoreNonRecordMarkdown(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	r := ArchiveRecord{Spec: "0256-example", Title: "Example", Source: "docs/specs/0256-example", SourceRevision: strings.Repeat("a", 40), Disposition: ArchivePass}
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string][]byte{
		r.Spec + ".md": content, "README.md": []byte("# Index\n"), "0256-prose.md": []byte("# Prose\n"),
		"other.md": []byte("---\nschema: " + ArchiveRecordSchema + "\n---\n"),
	} {
		if err := os.WriteFile(filepath.Join(root, name), content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	legacy := filepath.Join(root, "legacy")
	if err := os.Mkdir(legacy, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacy, "_prd.md"), []byte("# Legacy\n"), 0600); err != nil {
		t.Fatal(err)
	}
	slugs, err := ArchivedSpecSlugs(root)
	if err != nil || !slices.Equal(slugs, []string{r.Spec, "legacy"}) {
		t.Fatalf("slugs=%v err=%v", slugs, err)
	}
	if err := os.WriteFile(filepath.Join(root, "0256-malformed.md"), []byte("---\nschema: "+ArchiveRecordSchema+"\n---\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ArchivedSpecSlugs(root); err == nil {
		t.Fatal("malformed record silently skipped")
	}
}
