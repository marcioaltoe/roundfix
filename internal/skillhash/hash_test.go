// Suite: external skill hash compatibility
// Invariant: path-plus-content digests match skills CLI 1.5.19 under stable en-US collation.
// Boundary IN: already collected relative paths and file bytes.
// Boundary OUT: filesystem traversal and lock rendering, owned by skills and internal/baseline.
package skillhash

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestSumMatchesSkillsCLI1519(t *testing.T) {
	files := []File{
		{Path: "nested/É.md", Content: []byte("nested uppercase accent\n")},
		{Path: "-a", Content: []byte("hyphen\n")},
		{Path: "2-guide.md", Content: []byte("two\n")},
		{Path: "nested/deeper/ß.md", Content: []byte("eszett\n")},
		{Path: "A.md", Content: []byte("uppercase\n")},
		{Path: "nested/ä.md", Content: []byte("umlaut\n")},
		{Path: "_a", Content: []byte("underscore\n")},
		{Path: "10-guide.md", Content: []byte("ten\n")},
		{Path: "é.md", Content: []byte("precomposed accent\n")},
		{Path: "é.md", Content: []byte("decomposed accent\n")},
		{Path: "a.md", Content: []byte("lowercase\n")},
		{Path: "nested/deeper/ss.md", Content: []byte("double-s\n")},
		{Path: "nested/z.md", Content: []byte("zed\n")},
	}
	before := append([]File(nil), files...)

	const want = "2a46b6d704729eafc0148969028b9cc4030813059e1f7524def2f38b433011d4"
	if got := Sum(files); got != want {
		t.Fatalf("Sum() = %q, want skills CLI 1.5.19 digest %q", got, want)
	}
	if !reflect.DeepEqual(files, before) {
		t.Fatalf("Sum() changed caller-owned order:\ngot:  %#v\nwant: %#v", files, before)
	}
}

func TestSumSortsUnderscoreBeforeHyphen(t *testing.T) {
	files := []File{
		{Path: "-a", Content: []byte("hyphen\n")},
		{Path: "_a", Content: []byte("underscore\n")},
	}
	orderedBytes := []byte("_aunderscore\n-ahyphen\n")
	wantBytes := sha256.Sum256(orderedBytes)
	want := hex.EncodeToString(wantBytes[:])

	if got := Sum(files); got != want {
		t.Fatalf("Sum() = %q, want underscore-before-hyphen digest %q", got, want)
	}
}

func TestSumChangesWhenPathOrContentChanges(t *testing.T) {
	base := []File{
		{Path: "SKILL.md", Content: []byte("skill\n")},
		{Path: "references/guide.md", Content: []byte("guide\n")},
	}
	wantDifferentFrom := Sum(base)
	tests := []struct {
		name  string
		files []File
	}{
		{
			name: "path changes",
			files: []File{
				{Path: "SKILL.md", Content: []byte("skill\n")},
				{Path: "references/changed.md", Content: []byte("guide\n")},
			},
		},
		{
			name: "content changes",
			files: []File{
				{Path: "SKILL.md", Content: []byte("skill\n")},
				{Path: "references/guide.md", Content: []byte("guide changed\n")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Sum(test.files); got == wantDifferentFrom {
				t.Fatalf("Sum() = unchanged digest %q", got)
			}
		})
	}
}
