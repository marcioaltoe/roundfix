package baseline

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func renderProjectDecision(
	decisionID string,
	value any,
	declaration document,
) (string, error) {
	switch decisionID {
	case "identifier.strategy":
		return renderIdentifierStrategy(value)
	case httpContractDecisionID:
		return renderHTTPContract(value, declaration)
	case authProviderDecisionID:
		return renderAuthProvider(value)
	default:
		return "`" + renderDecisionValue(value) + "`", nil
	}
}

func renderIdentifierStrategy(value any) (string, error) {
	object, ok := objectValue(value)
	if !ok {
		return "", errors.New("structured render identifier strategy must be an object")
	}
	kind, _ := object["kind"].(string)
	const sourceContracts = "Preserve external provider identifiers, protocol identifiers, natural keys, and business codes according to their source contracts."
	switch kind {
	case "uuid-v7":
		return strings.Join([]string{
			"## Identifier strategy",
			"",
			"Use UUID version 7 for new project-owned Internal Identifiers only.",
			"",
			sourceContracts,
		}, "\n"), nil
	case "repository-defined":
		guidance, _ := object["guidance"].(string)
		if err := validateStructuredRenderText("identifier guidance", guidance, true); err != nil {
			return "", err
		}
		return strings.Join([]string{
			"## Identifier strategy",
			"",
			guidance,
			"",
			"This repository-defined rule applies to new project-owned Internal Identifiers only.",
			"",
			sourceContracts,
		}, "\n"), nil
	default:
		return "", fmt.Errorf("structured render identifier kind %q is not supported", kind)
	}
}

func renderHTTPContract(value any, declaration document) (string, error) {
	contract, err := normalizeHTTPContract(value, declaration)
	if err != nil {
		return "", fmt.Errorf("structured render HTTP Contract Decision: %w", err)
	}
	if err := validateStructuredRenderText("HTTP mode", contract.Mode, false); err != nil {
		return "", err
	}
	lines := []string{
		"## HTTP contract",
		"",
		"Application HTTP mode: **" + escapeMarkdownText(contract.Mode) + "**.",
	}
	if len(contract.Exceptions) == 0 {
		return strings.Join(append(lines, "", "There are no confirmed HTTP exceptions."), "\n"), nil
	}
	lines = append(lines, "", "Confirmed ordered exceptions:", "")
	for index, exception := range contract.Exceptions {
		rendered, renderErr := renderHTTPException(exception)
		if renderErr != nil {
			return "", fmt.Errorf("structured render HTTP exception %d: %w", index, renderErr)
		}
		lines = append(lines, fmt.Sprintf("%d. %s", index+1, rendered))
	}
	return strings.Join(lines, "\n"), nil
}

func renderAuthProvider(value any) (string, error) {
	provider, err := normalizeAuthProviderDecision(value)
	if err != nil {
		return "", fmt.Errorf("structured render authentication provider: %w", err)
	}
	exception := provider.RouteException
	if err := validateHTTPExceptionRenderText(exception); err != nil {
		return "", err
	}
	return strings.Join([]string{
		"### Better Auth",
		"",
		"Better Auth owns the authentication protocol for " + markdownCode(exception.Scope) + ". " +
			"Its confirmed " + renderHTTPMethods(exception.Methods) +
			" exception preserves this provider contract: " +
			escapeMarkdownText(exception.Reason),
	}, "\n"), nil
}

func renderHTTPException(exception HTTPException) (string, error) {
	if err := validateHTTPExceptionRenderText(exception); err != nil {
		return "", err
	}
	return "**" + escapeMarkdownText(exception.Owner) + "** owns " +
		renderHTTPMethods(exception.Methods) + " for " +
		markdownCode(exception.Scope) + ": " +
		escapeMarkdownText(exception.Reason), nil
}

func validateHTTPExceptionRenderText(exception HTTPException) error {
	for _, field := range []struct {
		name  string
		value string
	}{
		{name: "HTTP exception scope", value: exception.Scope},
		{name: "HTTP exception owner", value: exception.Owner},
		{name: "HTTP exception reason", value: exception.Reason},
	} {
		if err := validateStructuredRenderText(field.name, field.value, false); err != nil {
			return err
		}
	}
	return nil
}

func validateStructuredRenderText(field, value string, multiline bool) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("structured render %s is not valid UTF-8", field)
	}
	if value == "" || value != strings.TrimSpace(value) {
		return fmt.Errorf("structured render %s is non-canonical", field)
	}
	if strings.Contains(value, "\r") || !multiline && strings.Contains(value, "\n") {
		return fmt.Errorf("structured render %s is non-canonical", field)
	}
	lower := strings.ToLower(value)
	if strings.Contains(lower, "<!--") ||
		strings.Contains(lower, "-->") ||
		strings.Contains(lower, "setup-context-driven:") ||
		strings.Contains(lower, "roundfix:repository-rule:") ||
		strings.Contains(value, "{{") ||
		strings.Contains(value, "}}") {
		return fmt.Errorf("structured render %s contains marker-shaped content", field)
	}
	for _, line := range strings.Split(value, "\n") {
		if line != strings.TrimRight(line, " \t") {
			return fmt.Errorf("structured render %s is non-canonical", field)
		}
	}
	for _, character := range value {
		if unicode.IsControl(character) && character != '\n' {
			return fmt.Errorf("structured render %s contains a control character", field)
		}
	}
	return nil
}

func renderHTTPMethods(methods []string) string {
	rendered := make([]string, len(methods))
	for index, method := range methods {
		rendered[index] = markdownCode(method)
	}
	switch len(rendered) {
	case 0:
		return ""
	case 1:
		return rendered[0]
	case 2:
		return rendered[0] + " and " + rendered[1]
	default:
		return strings.Join(rendered[:len(rendered)-1], ", ") + ", and " + rendered[len(rendered)-1]
	}
}

func escapeMarkdownText(value string) string {
	const markdownPunctuation = `\` + "`*_{}[]<>"
	var rendered strings.Builder
	rendered.Grow(len(value))
	for _, character := range value {
		if strings.ContainsRune(markdownPunctuation, character) {
			rendered.WriteByte('\\')
		}
		rendered.WriteRune(character)
	}
	return rendered.String()
}

func markdownCode(value string) string {
	maxRun := 0
	currentRun := 0
	for _, character := range value {
		if character == '`' {
			currentRun++
			maxRun = max(maxRun, currentRun)
			continue
		}
		currentRun = 0
	}
	delimiter := strings.Repeat("`", maxRun+1)
	padding := ""
	if strings.HasPrefix(value, "`") || strings.HasSuffix(value, "`") ||
		strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") {
		padding = " "
	}
	return delimiter + padding + value + padding + delimiter
}
