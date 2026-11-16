package phaserecord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// RecordedPhase is one phase entry read back.
type RecordedPhase struct {
	Time       time.Time
	Phase      Phase
	Spec       string
	Repository string
}

// Phases returns phase entries at or after since, oldest first, counting invalid
// lines. Prompt entries do not count. A missing directory is an empty record.
func Phases(homeDir string, since time.Time) (entries []RecordedPhase, skipped int, err error) {
	dir := filepath.Join(homeDir, ".roundfix", "phase-record")
	files, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("read phase record directory: %w", err)
	}
	for _, file := range files {
		date, dateErr := time.Parse("2006-01-02.jsonl", file.Name())
		if dateErr != nil || date.Format("2006-01-02.jsonl") != file.Name() || file.Name() < since.UTC().Format("2006-01-02.jsonl") || !file.Type().IsRegular() {
			continue
		}
		data, readErr := os.ReadFile(filepath.Join(dir, file.Name()))
		if readErr != nil {
			return nil, skipped, fmt.Errorf("read phase record day %q: %w", file.Name(), readErr)
		}
		lines := bytes.Split(data, []byte{'\n'})
		if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		}
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			var header struct {
				Schema string `json:"schema"`
				Kind   string `json:"kind"`
			}
			if len(line) == 0 || line[0] != '{' || json.Unmarshal(line, &header) != nil || header.Schema != Schema {
				skipped++
				continue
			}
			if header.Kind != "phase" {
				continue
			}
			var row phaseLine
			if json.Unmarshal(line, &row) != nil {
				skipped++
				continue
			}
			at, timeErr := time.Parse(time.RFC3339, row.Time)
			phase, phaseErr := ParsePhase(string(row.Phase))
			if timeErr != nil || phaseErr != nil {
				skipped++
				continue
			}
			if !at.Before(since) {
				entries = append(entries, RecordedPhase{Time: at, Phase: phase, Spec: row.Spec, Repository: row.Repository})
			}
		}
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Time.Before(entries[j].Time) })
	return entries, skipped, nil
}
