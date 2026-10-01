package mdtree

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestTextReturnsTheEntryAloneWithoutCompanions(t *testing.T) {
	files := fstest.MapFS{"entry.md": {Data: []byte("entry\n")}}
	got, err := Text(files, "entry.md", "references")
	if err != nil || got != "entry\n" {
		t.Fatalf("Text() = %q, %v; want entry alone", got, err)
	}
}

func TestTextAppendsCompanionFilesInLexicalOrder(t *testing.T) {
	files := fstest.MapFS{
		"entry.md":        {Data: []byte("entry\n")},
		"references/z.md": {Data: []byte("last")},
		"references/a.md": {Data: []byte("first\n")},
	}
	got, err := Text(files, "entry.md", "references")
	if err != nil || got != "entry\n\nfirst\n\nlast" {
		t.Fatalf("Text() = %q, %v; want entry and lexical companions", got, err)
	}
}

func TestTextIgnoresNestedDirectoriesAndOtherSuffixes(t *testing.T) {
	files := fstest.MapFS{
		"entry.md":                      {Data: []byte("entry")},
		"references/direct.md":          {Data: []byte("direct")},
		"references/nested.md/child.md": {Data: []byte("nested")},
		"references/other.txt":          {Data: []byte("other")},
		"references/upper.MD":           {Data: []byte("upper")},
	}
	got, err := Text(files, "entry.md", "references")
	if err != nil || got != "entry\ndirect" {
		t.Fatalf("Text() = %q, %v; want only direct .md files", got, err)
	}
}

func TestTextReportsAMissingEntry(t *testing.T) {
	files := fstest.MapFS{"references/a.md": {Data: []byte("companion")}}
	if _, err := Text(files, "entry.md", "references"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Text() error = %v, want missing entry", err)
	}
}
