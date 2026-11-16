// Package phaserecord appends the identifiers and times of Supervisor Phases
// and Agent Session prompts to private day files under Roundfix Home.
package phaserecord

import (
	"bufio"
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

type PhaseEntry struct {
	Phase      Phase
	Spec       string
	Repository string
	Harness    string
	Session    string
}

type Boundary struct {
	Crossed       bool
	PreviousPhase Phase
	PreviousSpec  string
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
	Schema        string  `json:"schema"`
	Kind          string  `json:"kind"`
	Time          string  `json:"time"`
	Phase         Phase   `json:"phase"`
	Spec          string  `json:"spec"`
	Repository    string  `json:"repository"`
	Harness       string  `json:"harness"`
	Session       string  `json:"session"`
	Boundary      bool    `json:"boundary"`
	PreviousPhase *Phase  `json:"previous_phase,omitempty"`
	PreviousSpec  *string `json:"previous_spec,omitempty"`
}

// StartPhase appends a phase entry and reports whether the session crossed a boundary.
func (r *Recorder) StartPhase(ctx context.Context, entry PhaseEntry) (Boundary, error) {
	if err := ctx.Err(); err != nil {
		return Boundary{}, err
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
	boundary, err := r.boundary(ctx, now, row)
	if err != nil {
		return Boundary{}, fmt.Errorf("read phase record: %w", err)
	}
	row.Boundary = boundary.Crossed
	if boundary.Crossed {
		row.PreviousPhase = &boundary.PreviousPhase
		row.PreviousSpec = &boundary.PreviousSpec
	}
	data, err := json.Marshal(row)
	if err != nil {
		return Boundary{}, fmt.Errorf("encode phase entry: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return Boundary{}, err
	}
	if err := os.MkdirAll(r.directory, 0700); err != nil {
		return Boundary{}, fmt.Errorf("create phase record directory: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(r.directory, now.Format("2006-01-02")+".jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return Boundary{}, fmt.Errorf("open phase record: %w", err)
	}
	data = append(data, '\n')
	n, writeErr := file.Write(data)
	if writeErr == nil && n != len(data) {
		writeErr = io.ErrShortWrite
	}
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return Boundary{}, fmt.Errorf("write phase entry: %w", err)
	}
	return boundary, nil
}

func (r *Recorder) boundary(ctx context.Context, now time.Time, entry phaseLine) (Boundary, error) {
	if entry.Session == "" {
		return Boundary{}, nil
	}
	var latest phaseLine
	var latestTime time.Time
	for age := 6; age >= 0; age-- {
		path := filepath.Join(r.directory, now.AddDate(0, 0, -age).Format("2006-01-02")+".jsonl")
		err := readDay(ctx, path, func(line []byte) {
			var row phaseLine
			if json.Unmarshal(line, &row) != nil || row.Kind != "phase" || row.Session != entry.Session {
				return
			}
			recorded, err := time.Parse(time.RFC3339, row.Time)
			if err != nil || recorded.Before(latestTime) {
				return
			}
			latest, latestTime = row, recorded
		})
		if err != nil {
			return Boundary{}, err
		}
	}
	if latest.Session == "" || latest.Phase == "follow-up" ||
		(latest.Phase == entry.Phase && latest.Spec == entry.Spec) ||
		(latest.Phase == "delivery" && entry.Phase == "release" && (entry.Spec == latest.Spec || entry.Spec == "")) {
		return Boundary{}, nil
	}
	return Boundary{Crossed: true, PreviousPhase: Phase(bounded(string(latest.Phase))), PreviousSpec: bounded(latest.Spec)}, nil
}

// maxLineBytes bounds one line read back from a day file. A written entry holds
// at most eight strings of 1,024 bytes, which JSON escaping can grow sixfold, so a
// longer line is not an entry: it is discarded in buffer-sized pieces and the
// lines after it are still read.
const maxLineBytes = 64 << 10

// readDay passes each line to visit, which must not keep the slice.
func readDay(ctx context.Context, path string, visit func([]byte)) (err error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	reader := bufio.NewReaderSize(file, maxLineBytes)
	oversized := false
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		line, err := reader.ReadSlice('\n')
		switch {
		case errors.Is(err, bufio.ErrBufferFull):
			oversized = true
			continue
		case oversized:
			oversized = false
		case len(line) > 0:
			visit(line)
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
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
