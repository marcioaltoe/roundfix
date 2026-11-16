package phaserecord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type PromptEntry struct {
	Started, Ended time.Time
	RunID, Spec    string
	ScopeKind      string
	ScopeID        string
	Repository     string
	WorkDir        string
	Harness        string
	Session        string
}

// promptLine's declaration order is the public file format.
type promptLine struct {
	Schema     string `json:"schema"`
	Kind       string `json:"kind"`
	Time       string `json:"time"`
	Started    string `json:"started"`
	Ended      string `json:"ended"`
	RunID      string `json:"run_id"`
	Spec       string `json:"spec"`
	ScopeKind  string `json:"scope_kind"`
	ScopeID    string `json:"scope_id"`
	Repository string `json:"repository"`
	WorkDir    string `json:"work_dir"`
	Harness    string `json:"harness"`
	Session    string `json:"session"`
}

// RecordPrompt appends a prompt entry.
func (r *Recorder) RecordPrompt(ctx context.Context, entry PromptEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	now := r.now().UTC()
	row := promptLine{
		Schema: Schema, Kind: "prompt", Time: now.Format(time.RFC3339),
		Started: entry.Started.UTC().Format(time.RFC3339), Ended: entry.Ended.UTC().Format(time.RFC3339),
		RunID: bounded(entry.RunID), Spec: bounded(entry.Spec), ScopeKind: bounded(entry.ScopeKind), ScopeID: bounded(entry.ScopeID),
		Repository: bounded(entry.Repository), WorkDir: bounded(entry.WorkDir), Harness: bounded(entry.Harness), Session: bounded(entry.Session),
	}
	data, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("encode prompt entry: %w", err)
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
		return fmt.Errorf("write prompt entry: %w", err)
	}
	return nil
}

// RemoveBefore removes, or with dryRun counts, the day files dated before cutoff.
func (r *Recorder) RemoveBefore(ctx context.Context, cutoff time.Time, dryRun bool) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	entries, err := os.ReadDir(r.directory)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read phase record directory: %w", err)
	}
	cutoffDate := cutoff.UTC().Format("2006-01-02")
	removed := 0
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return removed, err
		}
		if entry.IsDir() {
			continue
		}
		date, err := time.Parse("2006-01-02.jsonl", entry.Name())
		if err != nil || date.Format("2006-01-02.jsonl") != entry.Name() || date.Format("2006-01-02") >= cutoffDate {
			continue
		}
		if !dryRun {
			if err := os.Remove(filepath.Join(r.directory, entry.Name())); err != nil {
				return removed, fmt.Errorf("remove phase record day file %q: %w", entry.Name(), err)
			}
		}
		removed++
	}
	return removed, nil
}
