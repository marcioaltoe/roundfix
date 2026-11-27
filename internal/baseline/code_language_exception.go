package baseline

import (
	"errors"
	"strings"
)

const codeLanguageExceptionDecisionID = "language.code_exception"
const codeLanguageExceptionDecisionType = "code-language-exception"

// CodeLanguageException records the existing codebase's language and its owner's decision.
type CodeLanguageException struct {
	Language string `json:"language"`
	Source   string `json:"source"`
}

func normalizeCodeLanguageException(value any) (CodeLanguageException, error) {
	object, ok := objectValue(value)
	if !ok {
		return CodeLanguageException{}, errors.New("must be an object with language and source")
	}
	for field := range object {
		if field != "language" && field != "source" {
			return CodeLanguageException{}, errors.New("accepts only language and source")
		}
	}
	source, _ := object["source"].(string)
	if strings.TrimSpace(source) == "" {
		return CodeLanguageException{}, errors.New("requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception")
	}
	language, _ := object["language"].(string)
	if strings.TrimSpace(language) == "" {
		return CodeLanguageException{}, errors.New("requires a language")
	}
	if strings.EqualFold(strings.TrimSpace(language), "English") {
		return CodeLanguageException{}, errors.New("names a language other than English, which is the rule and needs no exception")
	}
	if err := validateStructuredRenderText("code language exception language", language, false); err != nil {
		return CodeLanguageException{}, err
	}
	if err := validateStructuredRenderText("code language exception source", source, false); err != nil {
		return CodeLanguageException{}, err
	}
	return CodeLanguageException{Language: language, Source: source}, nil
}
