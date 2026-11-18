package phaserecord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FallbackEntry records the declared model, its limit and the chosen fallback.
type FallbackEntry struct {
	Phase                                        Phase
	Spec, Repository, RepositoryID               string
	Harness, Session                             string
	DeclaredModel, DeclaredEffort, Model, Effort string
	Limit, Message                               string
}

type RecordedFallback struct {
	Time time.Time
	FallbackEntry
}

// Declaration order is the public fallback line format, including empty fields.
type fallbackLine struct {
	Schema         string `json:"schema"`
	Kind           string `json:"kind"`
	Time           string `json:"time"`
	Phase          Phase  `json:"phase"`
	Spec           string `json:"spec"`
	Repository     string `json:"repository"`
	RepositoryID   string `json:"repository_id"`
	Harness        string `json:"harness"`
	Session        string `json:"session"`
	DeclaredModel  string `json:"declared_model"`
	DeclaredEffort string `json:"declared_effort"`
	Model          string `json:"model"`
	Effort         string `json:"effort"`
	Limit          string `json:"limit"`
	Message        string `json:"message"`
}

// RecordFallback appends once, without reading or changing earlier entries.
func (r *Recorder) RecordFallback(ctx context.Context, entry FallbackEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	now := r.now().UTC()
	row := fallbackLine{
		Schema: Schema, Kind: "fallback", Time: now.Format(time.RFC3339), Phase: entry.Phase,
		Spec: bounded(entry.Spec), Repository: bounded(entry.Repository), RepositoryID: bounded(entry.RepositoryID),
		Harness: bounded(entry.Harness), Session: bounded(entry.Session),
		DeclaredModel: bounded(entry.DeclaredModel), DeclaredEffort: bounded(entry.DeclaredEffort),
		Model: bounded(entry.Model), Effort: bounded(entry.Effort), Limit: bounded(entry.Limit), Message: bounded(entry.Message),
	}
	data, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("encode fallback entry: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(r.directory, 0700); err != nil {
		return fmt.Errorf("create phase record directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(r.directory, now.Format("2006-01-02")+".jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open phase record: %w", err)
	}
	data = append(data, '\n')
	n, writeErr := file.Write(data)
	if writeErr == nil && n != len(data) {
		writeErr = io.ErrShortWrite
	}
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("write fallback entry: %w", err)
	}
	return nil
}

// Fallbacks reads all fallback entries oldest first. Other entry kinds do not
// count as skipped. A missing directory is an empty record.
func Fallbacks(homeDir string) (entries []RecordedFallback, skipped int, err error) {
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
		if dateErr != nil || date.Format("2006-01-02.jsonl") != file.Name() || !file.Type().IsRegular() {
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
			if header.Kind != "fallback" {
				continue
			}
			var row fallbackLine
			if json.Unmarshal(line, &row) != nil {
				skipped++
				continue
			}
			at, timeErr := time.Parse(time.RFC3339, row.Time)
			if timeErr != nil || !row.Phase.IsAuthoring() {
				skipped++
				continue
			}
			entries = append(entries, RecordedFallback{Time: at, FallbackEntry: FallbackEntry{
				Phase: row.Phase, Spec: row.Spec, Repository: row.Repository, RepositoryID: row.RepositoryID,
				Harness: row.Harness, Session: row.Session, DeclaredModel: row.DeclaredModel, DeclaredEffort: row.DeclaredEffort,
				Model: row.Model, Effort: row.Effort, Limit: row.Limit, Message: row.Message,
			}})
		}
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].Time.Before(entries[j].Time) })
	return entries, skipped, nil
}
