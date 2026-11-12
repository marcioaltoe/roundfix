// Suite: skills lock version compatibility
// Invariant: supported lock versions retain their published folder digests and unknown versions are refused by name.
// Boundary IN: numeric version text and already collected relative paths and bytes.
// Boundary OUT: filesystem traversal, lock readers and rendering.
package skillhash

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestSumForLockVersionMatchesThePublishedSkillsCLI(t *testing.T) {
	t.Parallel()
	edge := []File{
		{Path: "SKILL.md", Content: []byte("---\nname: edge\n---\nbody\r\nline two\r\n")},
		{Path: "metadata.json", Content: []byte("meta\n")},
		{Path: "__pycache__/x.pyc", Content: []byte("cache\n")},
		{Path: "references/B.md", Content: []byte("upper\n")},
		{Path: "references/a.md", Content: []byte("lower\n")},
		{Path: "_x.md", Content: []byte("under\n")},
		{Path: "node_modules/pkg/index.js", Content: []byte("dep\n")},
	}
	tests := []struct {
		name    string
		version int
		files   []File
		want    string
	}{
		{"version 2 empty", 2, nil, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"version 2 single file", 2, []File{{Path: "SKILL.md", Content: []byte("# hi\n")}}, "93c36952d8843b33e68ae970ebabcc03e1687c9bb0dd1ecaf357060318952f55"},
		{"version 2 nested files", 2, []File{{Path: "docs/b.md", Content: []byte("b\n")}, {Path: "SKILL.md", Content: []byte("a\n")}}, "b22c4d22aa319d2924312fa18fcbfedc88fa89c74dbd8f298c683e1caf7b7f24"},
		{"version 2 edge folder", 2, edge, "e5adbe14f1b9b04eb73f28a1fc7e58f083b72328dbb7271c4a2813cbcadb42ac"},
		{"version 1 edge folder", 1, edge[:len(edge)-1], "416d74f34b15134cda0d15e3b9e7dd32f3b0fb6aaf9369c51d9931b37138dc73"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := make([]File, len(test.files))
			for i, file := range test.files {
				before[i] = File{Path: file.Path, Content: append([]byte(nil), file.Content...)}
			}
			got, err := SumForLockVersion(test.version, test.files)
			if err != nil || got != test.want {
				t.Fatalf("digest = %q, %v; want %q", got, err, test.want)
			}
			for i := range test.files {
				if !reflect.DeepEqual(test.files[i], before[i]) {
					t.Fatal("changed caller-owned files")
				}
			}
			if test.version == 1 && got != Sum(test.files) {
				t.Fatal("version 1 changed its legacy digest")
			}
		})
	}
	t.Run("version 2 exclusions and binary content", func(t *testing.T) {
		files := []File{{Path: "data.bin", Content: []byte("a\x00b\r\n")}, {Path: "text.md", Content: []byte("a\nb\n")}}
		want, err := SumForLockVersion(2, files)
		if err != nil {
			t.Fatal(err)
		}
		input := append([]File(nil), files...)
		input[1].Content = []byte("a\r\nb\r\n")
		for _, path := range []string{".git/config", "nested/.git/config", "nested/__pycache__/x.pyc", "__pypackages__/x.py", "nested/metadata.json", "metadata.json/x"} {
			input = append(input, File{Path: path, Content: []byte("ignored")})
		}
		got, err := SumForLockVersion(2, input)
		if err != nil || got != want {
			t.Fatalf("filtered digest = %q, %v; want %q", got, err, want)
		}
		input[0].Content = []byte("a\x00b\n")
		got, err = SumForLockVersion(2, input)
		if err != nil || got == want {
			t.Fatalf("binary CRLF must remain significant: %q, %v", got, err)
		}
	})
}

func TestParseLockVersionReadsOnlyVersionsOneAndTwo(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		text string
		want int
	}{{"1", 1}, {"2", 2}, {"1.0", 1}, {"2e0", 2}, {"0.2e1", 2}} {
		t.Run(test.text, func(t *testing.T) {
			got, err := ParseLockVersion(test.text)
			if err != nil || got != test.want {
				t.Fatalf("ParseLockVersion(%q) = %d, %v", test.text, got, err)
			}
		})
	}
	for _, text := range []string{"3", "0", "1.5", "-1", "", "null", "word", "1e9999", "1.0000000000000001", "1.9999999999999999", "+1", "0x1", "1/1"} {
		t.Run("refuse "+text, func(t *testing.T) {
			_, err := ParseLockVersion(text)
			want := fmt.Sprintf("skills-lock.json version %s is not supported; this Roundfix reads versions 1 and 2", text)
			var versionErr UnsupportedLockVersionError
			if !errors.Is(err, ErrUnsupportedLockVersion) || !errors.As(err, &versionErr) || versionErr.Version != text || err.Error() != want {
				t.Fatalf("refusal = %v; want %q and typed version", err, want)
			}
		})
	}
	if _, err := SumForLockVersion(3, nil); !errors.Is(err, ErrUnsupportedLockVersion) {
		t.Fatalf("unsupported digest version = %v", err)
	}
	if (UnsupportedLockVersionError{Version: "3"}).Is(errors.New("unrelated")) {
		t.Fatal("unsupported version matches unrelated error")
	}
}
