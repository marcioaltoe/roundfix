package judge

import (
	"regexp"
	"strings"
	"unicode"
)

// quotedText matches a passage in straight or curly double quotation marks,
// bounded so an unbalanced mark cannot swallow the rest of a document.
var quotedText = regexp.MustCompile(`"[^"]{0,1000}"|\x{201C}[^\x{201D}]{0,1000}\x{201D}`)

// unquoted drops blockquote lines and quoted passages, because a quotation
// keeps the language its speaker used and says nothing about the author's.
func unquoted(body string) string {
	lines := strings.Split(body, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), ">") {
			continue
		}
		kept = append(kept, line)
	}
	return quotedText.ReplaceAllString(strings.Join(kept, "\n"), " ")
}

func (gate LanguageGate) isEnglish(text string) bool {
	_, body, _, err := splitFrontMatter(text)
	if err != nil {
		return false
	}
	words := strings.FieldsFunc(strings.ToLower(unquoted(body)), func(r rune) bool { return !unicode.IsLetter(r) })
	if len(words) == 0 {
		return false
	}
	english, portuguese := 0, 0
	for _, word := range words {
		for _, candidate := range gate.EnglishWords {
			if word == candidate {
				english++
				break
			}
		}
		for _, candidate := range gate.PortugueseWords {
			if word == candidate {
				portuguese++
				break
			}
		}
	}
	return float64(english)/float64(len(words)) >= gate.MinEnglishShare && float64(portuguese)/float64(len(words)) < gate.MaxPortugueseShare
}
