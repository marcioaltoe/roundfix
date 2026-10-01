package judge

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Source can only be populated by the two bounded readers in this file.
type Source struct{ path, text string }

const maxSourceBytes = 1 << 20

var adrNumber = regexp.MustCompile(`^[0-9]{4}$`)
var inactiveStatus = regexp.MustCompile(`(?im)^\s*(?:\*\*)?status(?:\*\*)?:\s*(?:proposed|rejected|deprecated|superseded)\b`)

func readSpecArtifact(specDir, name string) (Source, error) {
	if name != "_prd.md" && name != "_techspec.md" {
		return Source{}, errors.New("not a Spec artifact")
	}
	path := filepath.Join(specDir, name)
	text, err := readRegular(path)
	if err != nil {
		return Source{}, err
	}
	return Source{path: path, text: text}, nil
}

func readADR(repoRoot, number string) (src Source, ok bool, err error) {
	if !adrNumber.MatchString(number) {
		return Source{}, false, errors.New("invalid ADR number")
	}
	// Reject linked directory components as well as linked records.
	for _, path := range []string{filepath.Join(repoRoot, "docs"), filepath.Join(repoRoot, "docs", "adr")} {
		info, err := os.Lstat(path)
		if err != nil {
			return Source{}, false, err
		}
		if !info.IsDir() {
			return Source{}, false, errors.New("ADR directory is not a regular directory")
		}
	}
	dir := filepath.Join(repoRoot, "docs", "adr")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Source{}, false, err
	}
	var path string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), number+"-") && strings.HasSuffix(entry.Name(), ".md") {
			if path != "" {
				return Source{}, false, errors.New("ambiguous ADR number")
			}
			path = filepath.Join(dir, entry.Name())
		}
	}
	if path == "" {
		return Source{}, false, os.ErrNotExist
	}
	text, err := readRegular(path)
	if err != nil {
		return Source{}, false, err
	}
	front, body, _, err := splitFrontMatter(text)
	if err != nil {
		return Source{}, false, err
	}
	var lifecycle struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal([]byte(front), &lifecycle); err != nil {
		return Source{}, false, fmt.Errorf("parse ADR lifecycle: %w", err)
	}
	if lifecycle.Status != "" && lifecycle.Status != "accepted" {
		return Source{}, false, nil
	}
	if lifecycle.Status == "" && inactiveStatus.MatchString(body) {
		return Source{}, false, nil
	}
	return Source{path: path, text: text}, true, nil
}

// Lstat rejects links, and the opened descriptor is checked before reading.
func readRegular(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Size() > maxSourceBytes {
		return "", errors.New("not a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil {
		return "", err
	}
	if !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return "", errors.New("source changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(f, maxSourceBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxSourceBytes {
		return "", errors.New("source exceeds size limit")
	}
	return string(data), nil
}

// The offset preserves original line numbers after removing front matter.
func splitFrontMatter(text string) (front, body string, offset int, err error) {
	lines := strings.Split(text, "\n")
	if strings.TrimSuffix(lines[0], "\r") != "---" {
		return "", text, 0, nil
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSuffix(lines[i], "\r") == "---" {
			return strings.Join(lines[1:i], "\n"), strings.Join(lines[i+1:], "\n"), i + 1, nil
		}
	}
	return "", "", 0, errors.New("missing front matter closing marker")
}
