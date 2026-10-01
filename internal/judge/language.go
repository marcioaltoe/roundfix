package judge

import (
	"strings"
	"unicode"
)

func (gate LanguageGate) isEnglish(text string) bool {
	_, body, _, err := splitFrontMatter(text)
	if err != nil {
		return false
	}
	words := strings.FieldsFunc(strings.ToLower(body), func(r rune) bool { return !unicode.IsLetter(r) })
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
