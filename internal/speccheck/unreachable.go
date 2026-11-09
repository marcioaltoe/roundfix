package speccheck

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"roundfix/internal/spec"
)

// CodeUnreachableUnread identifies section text the QA declaration reader drops.
const CodeUnreachableUnread = "SC-UNREACHABLE-UNREAD"

func detectUnreachableDeclarations(result *Result, repoRoot, prdPath string) error {
	content, err := os.ReadFile(prdPath)
	if err != nil {
		return fmt.Errorf("read Spec PRD %q: %w", prdPath, err)
	}
	unread, err := spec.UnreachableUnread(content, prdPath)
	var declarationErr spec.UnreachableDeclarationError
	if err != nil && !errors.As(err, &declarationErr) {
		return fmt.Errorf("read Unreachable Acceptance in %q: %w", prdPath, err)
	}
	const fix = `Write each entry in the fielded form: a top-level "- criterion: <text the QA row repeats>" line, then indented "reason: <why no Verification or QA row can reach it>" and "satisfied-by: <evidence or follow-up that satisfies it>" lines; start the section with None. when no criterion is unreachable.`
	display := artifactDisplayPath(repoRoot, prdPath)
	for index, line := range unread {
		if index > 0 && line.Line == unread[index-1].Line+1 {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeUnreachableUnread, Severity: SeverityError,
			Summary: fmt.Sprintf("%s line %d under ## Unreachable Acceptance is not part of a fielded declaration, so the QA gate does not read it: %s", display, line.Line, strconv.Quote(line.Text)),
			Where:   []Location{{Path: display, Line: line.Line}}, Fix: fix,
		})
	}
	if err != nil {
		result.Findings = append(result.Findings, Finding{
			Code: CodeUnreachableUnread, Severity: SeverityError,
			Summary: fmt.Sprintf("%s line %d starts an Unreachable Acceptance declaration without its %s: field", display, declarationErr.Line, declarationErr.Field),
			Where:   []Location{{Path: display, Line: declarationErr.Line}}, Fix: fix,
		})
	}
	return nil
}
