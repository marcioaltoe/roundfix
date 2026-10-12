package spec

import (
	"fmt"
	"os"
	"regexp"
)

var archivePullRequestNumber = regexp.MustCompile(`^[1-9][0-9]{0,9}$`)

// RecordArchivePullRequest stores the delivery pull request in an Archive
// Record. Repeating the same number preserves the file byte for byte; a
// different stored number is refused without writing.
func RecordArchivePullRequest(recordPath, number string) (bool, error) {
	if !archivePullRequestNumber.MatchString(number) {
		return false, fmt.Errorf("pull request number must match ^[1-9][0-9]{0,9}$")
	}
	content, err := os.ReadFile(recordPath)
	if err != nil {
		return false, fmt.Errorf("read Archive Record: %w", err)
	}
	record, err := ParseArchiveRecord(content)
	if err != nil {
		return false, fmt.Errorf("parse Archive Record: %w", err)
	}
	if record.PullRequest == number {
		return false, nil
	}
	if record.PullRequest != "" {
		return false, fmt.Errorf("Archive Record already records pull request %s; cannot replace it with %s", record.PullRequest, number)
	}
	record.PullRequest = number
	content, err = RenderArchiveRecord(record)
	if err != nil {
		return false, fmt.Errorf("render Archive Record: %w", err)
	}
	if err := os.WriteFile(recordPath, content, 0644); err != nil {
		return false, fmt.Errorf("write Archive Record: %w", err)
	}
	return true, nil
}
