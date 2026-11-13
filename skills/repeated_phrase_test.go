// Suite: owned skill repeated phrase contract
// Invariant: no line of an owned skill ends with a phrase of two or more words that the next line begins with
// Boundary: canonical and shipped copies of every owned skill file; a line that is only the phrase and Markdown table rows stay exempt

package skills

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type repeatedLinePhrase struct {
	Line   int
	Phrase string
}

func repeatedLinePhrases(text string) []repeatedLinePhrase {
	lines := strings.Split(text, "\n")
	var phrases []repeatedLinePhrase

	for line := 0; line+1 < len(lines); line++ {
		trailing := strings.Fields(lines[line])
		leading := strings.Fields(lines[line+1])
		if len(trailing) == 0 || len(leading) == 0 ||
			strings.HasPrefix(trailing[0], "|") || strings.HasPrefix(leading[0], "|") {
			continue
		}

		limit := len(trailing)
		if len(leading) < limit {
			limit = len(leading)
		}
		for size := limit; size >= 2; size-- {
			equal := true
			for word := 0; word < size; word++ {
				if strings.ToLower(trailing[len(trailing)-size+word]) != strings.ToLower(leading[word]) {
					equal = false
					break
				}
			}
			if !equal {
				continue
			}
			if size != len(trailing) && size != len(leading) {
				phrases = append(phrases, repeatedLinePhrase{
					Line:   line + 1,
					Phrase: strings.Join(leading[:size], " "),
				})
			}
			break
		}
	}

	return phrases
}

func TestRepeatedLinePhrasesFindsAPhraseRepeatedAcrossALineBreak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []repeatedLinePhrase
	}{
		{
			name: "doubled-instruction",
			text: "rejected. Begin the relation with\nBegin the relation with `not applicable`, or with",
			want: []repeatedLinePhrase{{Line: 1, Phrase: "Begin the relation with"}},
		},
		{
			name: "two-words-in-another-case",
			text: "Record the result in the Task\nthe task file keeps it.",
			want: []repeatedLinePhrase{{Line: 1, Phrase: "the task"}},
		},
		{
			name: "one-repeated-word",
			text: "The guide says that\nthat rule still holds.",
		},
		{
			name: "command-listing",
			text: "roundfix gc compact\nroundfix gc compact --apply",
		},
		{
			name: "table-rows",
			text: "| id | title | |\n| | | done |",
		},
		{
			name: "single-instruction",
			text: "rejected. Begin the relation with\n`not applicable`, or with `não se aplica`",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := repeatedLinePhrases(test.text); !reflect.DeepEqual(got, test.want) {
				t.Errorf("repeatedLinePhrases() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestNoOwnedSkillRepeatsAPhraseAcrossALineBreak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		root string
	}{
		{name: "canonical", root: filepath.Join("..", ".agents", "skills")},
		{name: "shipped", root: "."},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			for _, name := range Names() {
				root := filepath.Join(test.root, name)
				err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}
					if entry.IsDir() {
						return nil
					}
					data, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					for _, phrase := range repeatedLinePhrases(string(data)) {
						t.Errorf("%s:%d ends with %q and the next line begins with it; write the phrase once", filepath.ToSlash(path), phrase.Line, phrase.Phrase)
					}
					return nil
				})
				if err != nil {
					t.Fatalf("walk %s: %v", root, err)
				}
			}
		})
	}
}
