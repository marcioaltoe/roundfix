// Package phaserecord appends the identifiers and times of Supervisor Phases
// and Agent Session prompts to private day files under Roundfix Home.
package phaserecord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const Schema = "roundfix/phase-record/v1"

// Phase is one of the seven Supervisor Phases.
type Phase string

func ParsePhase(name string) (Phase, error) {
	switch name {
	case "follow-up", "idea", "prd", "techspec", "tasks", "delivery", "release":
		return Phase(name), nil
	default:
		return "", fmt.Errorf("unknown phase %q; use follow-up, idea, prd, techspec, tasks, delivery or release", name)
	}
}

// IsAuthoring reports whether the phase produces a Spec authoring artifact.
func (p Phase) IsAuthoring() bool {
	switch p {
	case "idea", "prd", "techspec", "tasks":
		return true
	default:
		return false
	}
}

type PhaseEntry struct {
	Phase        Phase
	Spec         string
	Repository   string
	Harness      string
	Session      string
	RepositoryID string
	Model        string
	ModelSource  string
	Effort       string
	EffortSource string
	Runner       string
}

type Recorder struct {
	directory string
	now       func() time.Time
}

func New(homeDir string, now func() time.Time) *Recorder {
	return &Recorder{directory: filepath.Join(homeDir, ".roundfix", "phase-record"), now: now}
}

// phaseLine's declaration order is the public file format.
type phaseLine struct {
	Schema     string `json:"schema"`
	Kind       string `json:"kind"`
	Time       string `json:"time"`
	Phase      Phase  `json:"phase"`
	Spec       string `json:"spec"`
	Repository string `json:"repository"`
	Harness    string `json:"harness"`
	Session    string `json:"session"`
	Boundary   bool   `json:"boundary"`
	*phaseModel
}

// phaseModel is present only on authoring entries; its order follows phaseLine.
type phaseModel struct {
	RepositoryID string `json:"repository_id"`
	Model        string `json:"model"`
	ModelSource  string `json:"model_source"`
	Effort       string `json:"effort"`
	EffortSource string `json:"effort_source"`
	Runner       string `json:"runner"`
}

func sourceOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return bounded(value)
}

// StartPhase appends a phase entry without reading earlier entries.
func (r *Recorder) StartPhase(ctx context.Context, entry PhaseEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	now := r.now().UTC()
	row := phaseLine{
		Schema:     Schema,
		Kind:       "phase",
		Time:       now.Format(time.RFC3339),
		Phase:      Phase(bounded(string(entry.Phase))),
		Spec:       bounded(entry.Spec),
		Repository: bounded(entry.Repository),
		Harness:    bounded(entry.Harness),
		Session:    bounded(entry.Session),
	}
	if entry.Phase.IsAuthoring() {
		row.phaseModel = &phaseModel{
			RepositoryID: bounded(entry.RepositoryID), Model: bounded(entry.Model), ModelSource: sourceOrUnknown(entry.ModelSource),
			Effort: bounded(entry.Effort), EffortSource: sourceOrUnknown(entry.EffortSource), Runner: sourceOrUnknown(entry.Runner),
		}
	}
	data, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("encode phase entry: %w", err)
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
		return fmt.Errorf("write phase entry: %w", err)
	}
	return nil
}

// Keep cuts on UTF-8 boundaries so JSON decoding cannot expand a broken rune.
func bounded(value string) string {
	value = strings.ToValidUTF8(value, "�")
	if len(value) <= 1024 {
		return value
	}
	end := 1024
	for !utf8.RuneStart(value[end]) {
		end--
	}
	return value[:end]
}
