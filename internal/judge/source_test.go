package judge

// Suite: bounded text sources and language.
// Invariant: only regular, bounded artifacts and active ADRs become Sources.
// Boundary IN: temporary files and both readers.
// Boundary OUT: planner state, covered in pairs_test.go.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFixture(t *testing.T, root, name, text string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadersAcceptOnlySpecArtifactsAndADRs(t *testing.T) {
	for _, name := range []string{"_prd.md", "_techspec.md"} {
		t.Run("accept "+name, func(t *testing.T) {
			root := t.TempDir()
			path := writeFixture(t, root, name, "text")
			src, err := readSpecArtifact(root, name)
			if err != nil || src.path != path || src.text != "text" {
				t.Fatalf("source = %+v, %v", src, err)
			}
		})
	}
	for _, name := range []string{"task_01.md", "../_prd.md", "nested/_prd.md", "/_prd.md"} {
		t.Run("refuse name "+name, func(t *testing.T) {
			if _, err := readSpecArtifact(t.TempDir(), name); err == nil {
				t.Fatal("accepted another artifact")
			}
		})
	}
	for _, reader := range []string{"spec", "adr"} {
		for _, kind := range []string{"missing", "directory", "symlink", "oversized", "exact limit"} {
			t.Run(reader+" "+kind, func(t *testing.T) {
				root := t.TempDir()
				name := "_prd.md"
				if reader == "adr" {
					name = "docs/adr/0123-fixture.md"
				}
				path := filepath.Join(root, name)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				switch kind {
				case "directory":
					if err := os.Mkdir(path, 0o755); err != nil {
						t.Fatal(err)
					}
				case "symlink":
					target := writeFixture(t, root, "target.md", "text")
					if err := os.Symlink(target, path); err != nil {
						t.Fatal(err)
					}
				case "oversized":
					writeFixture(t, root, name, strings.Repeat("x", maxSourceBytes+1))
				case "exact limit":
					writeFixture(t, root, name, strings.Repeat("x", maxSourceBytes))
				}
				var err error
				ok := false
				if reader == "spec" {
					_, err = readSpecArtifact(root, "_prd.md")
					ok = err == nil
				} else {
					_, ok, err = readADR(root, "0123")
				}
				if ok != (kind == "exact limit") {
					t.Fatalf("accepted=%v error=%v", ok, err)
				}
			})
		}
	}
	for _, status := range []string{"accepted", "proposed", "rejected", "deprecated", "superseded", "unknown", "", "legacy inactive", "no status", "malformed"} {
		t.Run("ADR status "+status, func(t *testing.T) {
			root := t.TempDir()
			text := "# The decision\nThe record keeps the gate."
			switch status {
			case "":
			case "legacy inactive":
				text += "\n**Status**: superseded"
			case "no status":
				text = "---\ncreated: yesterday\n---\n" + text
			case "malformed":
				text = "---\nstatus: [\n---\n" + text
			default:
				text = "---\nstatus: " + status + "\n---\n" + text
			}
			writeFixture(t, root, "docs/adr/0123-decision.md", text)
			_, ok, err := readADR(root, "0123")
			want := status == "accepted" || status == "" || status == "no status"
			if ok != want {
				t.Fatalf("accepted=%v want=%v err=%v", ok, want, err)
			}
		})
	}
	t.Run("ADR paths and numbers", func(t *testing.T) {
		root := t.TempDir()
		writeFixture(t, root, "docs/adr/nested/0123-decision.md", "text")
		writeFixture(t, root, "docs/adr/0123.txt", "text")
		writeFixture(t, root, "docs/adr/01234-decision.md", "text")
		for _, number := range []string{"0123", "123", "../0123", "01234"} {
			if _, ok, _ := readADR(root, number); ok {
				t.Fatalf("accepted %q", number)
			}
		}
	})
	t.Run("linked ADR directory", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, "records")
		writeFixture(t, root, "records/0123-decision.md", "text")
		if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(root, "docs/adr")); err != nil {
			t.Fatal(err)
		}
		if _, ok, _ := readADR(root, "0123"); ok {
			t.Fatal("followed ADR directory link")
		}
	})
}

func TestLanguageGateSeparatesEnglishFromPortuguese(t *testing.T) {
	q := loadQuestions(t)
	for _, tc := range []struct {
		name, text string
		want       bool
	}{
		{"English", "The decision keeps the gate and it is the rule for the author.", true},
		{"Portuguese", "A decisão é uma regra para os autores e não está na sua documentação.", false},
		{"Unicode letters", "ÉéTHEÉé ééOFéé ééANDéé", false},
		{"case and punctuation", "THE,OF!AND?TO-IN;IS", true},
		{"empty", "1234 -", false},
		{"front matter excluded", "---\ntext: de que não uma para com os das dos\n---\nThe gate is the rule for the author.", true},
		{"quoted Portuguese excluded", "The maintainer asked for this rule: \"podemos ter em uma mesma spec um ou mais findings, não é obrigatório e nem desejável uma spec para cada um, o cuidado é não ter specs grandes\". The rule is that it is grouped.", true},
		{"curly-quoted Portuguese excluded", "The maintainer asked: \u201cpodemos alterar e atualizar um finding que ainda não foi implementado para que não seja necessário criar outro\u201d. It is the rule.", true},
		{"blockquoted Portuguese excluded", "The maintainer asked for it.\n> podemos ter em uma mesma spec um ou mais findings, não é obrigatório para os autores\nThe rule is that it is grouped.", true},
		{"quoted English in Portuguese stays Portuguese", "A decisão é uma regra para os autores e não está na documentação, como \"the rule is that it is the gate for the author\" diz.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := q.Language.isEnglish(tc.text); got != tc.want {
				t.Fatalf("English=%v want=%v", got, tc.want)
			}
		})
	}
	t.Run("measured boundaries", func(t *testing.T) {
		// Read the measured words and thresholds rather than copying them.
		english := q.Language.EnglishWords[0]
		portuguese := q.Language.PortugueseWords[0]
		words := make([]string, 100)
		for i := range words {
			words[i] = "xyz"
		}
		for i := 0; i < int(q.Language.MinEnglishShare*100); i++ {
			words[i] = english
		}
		if !q.Language.isEnglish(strings.Join(words, " ")) {
			t.Fatal("minimum English share rejected")
		}
		words[0] = "xyz"
		if q.Language.isEnglish(strings.Join(words, " ")) {
			t.Fatal("below minimum English share accepted")
		}
		for i := range words {
			words[i] = english
		}
		for i := 0; i < int(q.Language.MaxPortugueseShare*100); i++ {
			words[i] = portuguese
		}
		if q.Language.isEnglish(strings.Join(words, " ")) {
			t.Fatal("Portuguese ceiling accepted")
		}
	})
}
