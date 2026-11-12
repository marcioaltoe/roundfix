package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/app"
)

const SetupManifestPath = "docs/agents/setup-context.json"

var releaseVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// IsReleaseVersion accepts only a plain MAJOR.MINOR.PATCH release number.
func IsReleaseVersion(version string) bool {
	return releaseVersionPattern.MatchString(version)
}

// NewerRepositoryError identifies a Baseline written by a newer Roundfix.
type NewerRepositoryError struct {
	Path, Recorded, Running string
}

func (err *NewerRepositoryError) Error() string {
	return fmt.Sprintf("Setup Manifest %q records Roundfix %s, newer than this binary's %s; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'", err.Path, err.Recorded, err.Running)
}

// RecordedRoundfixVersion reads only the optional record, leaving other Setup
// Manifest validation to the Baseline reader. Missing or unreadable manifests
// and content without an object field do not record a release.
func RecordedRoundfixVersion(gitRoot string) (string, error) {
	if strings.TrimSpace(gitRoot) == "" {
		return "", nil
	}
	path, err := filepath.Abs(filepath.Join(gitRoot, SetupManifestPath))
	if err != nil {
		return "", nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	var manifest map[string]json.RawMessage
	if json.Unmarshal(content, &manifest) != nil {
		return "", nil
	}
	raw, present := manifest["roundfixVersion"]
	if !present {
		return "", nil
	}
	var version string
	if json.Unmarshal(raw, &version) != nil || !IsReleaseVersion(version) {
		return "", fmt.Errorf("Setup Manifest %q records roundfixVersion %s, which is not a Roundfix release version such as 0.140.0; restore the value the Roundfix that wrote it recorded", path, raw)
	}
	return version, nil
}

// CheckRecordedRoundfixVersion refuses only a strictly newer release or a
// malformed record. Development builds skip the read entirely.
func CheckRecordedRoundfixVersion(gitRoot string, running string) error {
	if app.IsDevelopmentVersion(running) {
		return nil
	}
	recorded, err := RecordedRoundfixVersion(gitRoot)
	if err != nil {
		return err
	}
	if recorded == "" || app.CompareVersions(recorded, running) <= 0 {
		return nil
	}
	path, err := filepath.Abs(filepath.Join(gitRoot, SetupManifestPath))
	if err != nil {
		return err
	}
	return &NewerRepositoryError{Path: path, Recorded: recorded, Running: app.NormalizeVersion(running)}
}
