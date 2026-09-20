package spec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const SupersessionFilename = "_supersession.md"

var (
	ErrNoSupersession     = errors.New("no supersession record")
	ErrSupersessionExists = errors.New("supersession record already exists")
)

// Supersession records that another Spec delivered this Spec's content.
type Supersession struct {
	SupersededBy string `yaml:"superseded_by"`
	Date         string `yaml:"date"`
	Reason       string `yaml:"reason"`
	Explanation  string `yaml:"-"`
}

// ReadSupersession reads and validates a Spec's supersession amendment.
func ReadSupersession(specDir string) (Supersession, error) {
	path := filepath.Join(specDir, SupersessionFilename)
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Supersession{}, fmt.Errorf("%w at %q", ErrNoSupersession, path)
		}
		return Supersession{}, fmt.Errorf("read supersession record %q: %w", path, err)
	}
	frontmatter, body, err := splitFrontmatter(content)
	if err != nil {
		return Supersession{}, fmt.Errorf("parse supersession record %q: %w", path, err)
	}
	var record Supersession
	if err := yaml.Unmarshal(frontmatter, &record); err != nil {
		return Supersession{}, fmt.Errorf("parse supersession record %q frontmatter: %w", path, err)
	}
	record.SupersededBy = strings.TrimSpace(record.SupersededBy)
	record.Date = strings.TrimSpace(record.Date)
	record.Reason = strings.TrimSpace(record.Reason)
	record.Explanation = strings.TrimSpace(string(body))
	if record.SupersededBy == "" {
		return Supersession{}, fmt.Errorf("parse supersession record %q: superseded_by is required", path)
	}
	parsedDate, err := time.Parse(time.DateOnly, record.Date)
	if err != nil || parsedDate.Format(time.DateOnly) != record.Date {
		return Supersession{}, fmt.Errorf("parse supersession record %q: date must be YYYY-MM-DD", path)
	}
	if record.Reason == "" {
		return Supersession{}, fmt.Errorf("parse supersession record %q: reason is required", path)
	}
	if record.Explanation == "" {
		return Supersession{}, fmt.Errorf("parse supersession record %q: explanation body is required", path)
	}
	return record, nil
}

// WriteSupersession creates a Spec's supersession amendment without replacing
// an existing record.
func WriteSupersession(specDir string, supersedingSlug string, reason string, recordedAt time.Time) (Supersession, error) {
	record := Supersession{
		SupersededBy: strings.TrimSpace(supersedingSlug),
		Reason:       strings.TrimSpace(reason),
		Explanation:  strings.TrimSpace(reason),
	}
	if record.SupersededBy == "" {
		return Supersession{}, errors.New("superseding Spec slug is required")
	}
	if record.Reason == "" {
		return Supersession{}, errors.New("supersession reason is required")
	}
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}
	record.Date = recordedAt.UTC().Format(time.DateOnly)

	frontmatter, err := yaml.Marshal(record)
	if err != nil {
		return Supersession{}, fmt.Errorf("encode supersession record: %w", err)
	}
	content := append([]byte("---\n"), frontmatter...)
	content = append(content, []byte("---\n\n")...)
	content = append(content, []byte(record.Explanation+"\n")...)

	path := filepath.Join(specDir, SupersessionFilename)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return Supersession{}, fmt.Errorf("%w at %q", ErrSupersessionExists, path)
		}
		return Supersession{}, fmt.Errorf("create supersession record %q: %w", path, err)
	}
	writtenBytes, writeErr := file.Write(content)
	if writeErr == nil && writtenBytes != len(content) {
		writeErr = io.ErrShortWrite
	}
	if writeErr != nil {
		closeErr := file.Close()
		removeErr := os.Remove(path)
		return Supersession{}, fmt.Errorf("write supersession record %q: %w", path, errors.Join(writeErr, closeErr, removeErr))
	}
	if err := file.Close(); err != nil {
		removeErr := os.Remove(path)
		return Supersession{}, fmt.Errorf("close supersession record %q: %w", path, errors.Join(err, removeErr))
	}
	return record, nil
}
