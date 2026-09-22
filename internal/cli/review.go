package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	roundconfig "roundfix/internal/config"
)

type reviewOutcome string

const (
	reviewOutcomeReviewed reviewOutcome = "reviewed"
	reviewOutcomeFindings reviewOutcome = "findings"
	reviewOutcomeBlocked  reviewOutcome = "blocked"
	reviewOutcomeOmitted  reviewOutcome = "omitted"
)

type reviewRecord struct {
	Repository string        `json:"repository"`
	BaseCommit string        `json:"baseCommit"`
	HeadCommit string        `json:"headCommit"`
	Provider   string        `json:"provider"`
	Source     string        `json:"source"`
	Outcome    reviewOutcome `json:"outcome"`
	Findings   string        `json:"findings,omitempty"`
	Reason     string        `json:"reason,omitempty"`
}

func newReviewRecord(
	repository string,
	baseCommit string,
	headCommit string,
	policy roundconfig.PrePRReview,
	outcome reviewOutcome,
) reviewRecord {
	return reviewRecord{
		Repository: repository,
		BaseCommit: baseCommit,
		HeadCommit: headCommit,
		Provider:   policy.Provider,
		Source:     policy.Source,
		Outcome:    outcome,
	}
}

func writeReviewRecord(writer io.Writer, record reviewRecord) error {
	if err := validateReviewRecord(record); err != nil {
		return err
	}
	if err := json.NewEncoder(writer).Encode(record); err != nil {
		return fmt.Errorf("encode review record: %w", err)
	}
	return nil
}

func validateReviewRecord(record reviewRecord) error {
	switch {
	case strings.TrimSpace(record.Repository) == "":
		return errors.New("review record repository is required")
	case strings.TrimSpace(record.BaseCommit) == "":
		return errors.New("review record base commit is required")
	case strings.TrimSpace(record.HeadCommit) == "":
		return errors.New("review record head commit is required")
	case strings.TrimSpace(record.Provider) == "":
		return errors.New("review record provider is required")
	case strings.TrimSpace(record.Source) == "":
		return errors.New("review record source is required")
	}

	switch record.Outcome {
	case reviewOutcomeReviewed, reviewOutcomeOmitted:
		if record.Findings != "" || record.Reason != "" {
			return fmt.Errorf("review record outcome %q cannot carry findings or a reason", record.Outcome)
		}
	case reviewOutcomeFindings:
		if strings.TrimSpace(record.Findings) == "" {
			return errors.New("review record findings are required for findings outcome")
		}
		if record.Reason != "" {
			return errors.New("review record findings outcome cannot carry a reason")
		}
	case reviewOutcomeBlocked:
		if strings.TrimSpace(record.Reason) == "" {
			return errors.New("review record reason is required for blocked outcome")
		}
		if record.Findings != "" {
			return errors.New("review record blocked outcome cannot carry findings")
		}
	default:
		return fmt.Errorf("review record outcome %q is invalid", record.Outcome)
	}

	return nil
}
