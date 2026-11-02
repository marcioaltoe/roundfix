package speccheck

import (
	"errors"
	"fmt"
	"os"

	"roundfix/internal/spec"
)

// CodePRDType identifies an unsupported PRD Delivery Type.
const CodePRDType = "SC-PRD-TYPE"

func detectPRDDeliveryType(result *Result, repoRoot, prdPath string) error {
	content, err := os.ReadFile(prdPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read PRD type %q: %w", prdPath, err)
	}
	_, _, err = spec.ParsePRDDeliveryType(content)
	if err == nil {
		return nil
	}
	var unsupported *spec.UnsupportedDeliveryTypeError
	if !errors.As(err, &unsupported) {
		return err
	}
	result.Findings = append(result.Findings, Finding{
		Code: CodePRDType, Severity: SeverityError,
		Summary: err.Error(),
		Where:   []Location{{Path: artifactDisplayPath(repoRoot, prdPath)}},
		Fix:     "Set the PRD's type to feat, fix, perf or refactor, or remove it to deliver as feat.",
	})
	return nil
}
