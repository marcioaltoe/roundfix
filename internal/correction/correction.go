// Package correction owns append-only declarations and delivery sign snapshots.
// It neither starts processes nor reads Git.
package correction

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const Schema = "roundfix/corrections/v1"

type Class string
type Artifact string

const NotACorrection Class = "not-a-correction"

func DefectClasses() []Class {
	return []Class{"plan", "implementation", "environment", "operator", "unclear"}
}
func Answers() []Class { return append(DefectClasses(), NotACorrection) }

type Sign struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Subject string `json:"subject,omitempty"`
}
type Declaration struct {
	Time                 time.Time
	Correction           string
	Class                Class
	Artifact             Artifact
	Reason               string
	Signs                []Sign
	Withdrawn            bool
	By, Harness, Session string
}
type Snapshot struct {
	Time        time.Time
	Head        string
	Signs       []Sign
	Dismissed   int
	Unavailable string
}
type File struct {
	Declarations []Declaration
	Snapshots    []Snapshot
	Skipped      int
}

// Path accepts the digest of the shared Repository Identity.
func Path(homeDir, repositoryID, slug string) string {
	key := strings.TrimPrefix(repositoryID, "sha256:")
	if len(key) > 16 {
		key = key[:16]
	}
	return filepath.Join(homeDir, ".roundfix", "corrections", key, slug+".jsonl")
}

var correctionID = regexp.MustCompile(`^C[1-9][0-9]*$`)
var taskID = regexp.MustCompile(`^task_[0-9]+$`)

func oneLine(value string, max int) bool {
	if !utf8.ValidString(value) || utf8.RuneCountInString(value) < 1 || utf8.RuneCountInString(value) > max {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			return false
		}
	}
	return true
}

// Validate validates a whole declaration; omitted fields of a redeclaration
// are merged by its writer before this boundary.
func Validate(d Declaration) error {
	accepted := false
	for _, class := range Answers() {
		if d.Class == class {
			accepted = true
		}
	}
	if !accepted {
		return fmt.Errorf("class must be plan, implementation, environment, operator, unclear or not-a-correction")
	}
	if d.Class == "plan" {
		if d.Artifact != "prd" && d.Artifact != "techspec" && d.Artifact != "tasks" {
			return fmt.Errorf("class plan requires --artifact prd, techspec or tasks")
		}
	} else if d.Artifact != "" {
		return fmt.Errorf("--artifact prd, techspec or tasks is accepted only with class plan")
	}
	if !oneLine(strings.TrimSpace(d.Reason), 200) {
		return fmt.Errorf("--reason must be 1 to 200 characters on one line with no control character or line or paragraph separator")
	}
	if d.By != "" && !oneLine(strings.TrimSpace(d.By), 64) {
		return fmt.Errorf("--by must be 1 to 64 characters on one line with no control character or line or paragraph separator")
	}
	for _, sign := range d.Signs {
		if sign.Kind == "task" && !taskID.MatchString(sign.ID) {
			return fmt.Errorf("--task must match task_<digits>")
		}
	}
	if d.Correction != "" && !correctionID.MatchString(d.Correction) {
		return fmt.Errorf("--correction must match C<n>, with n a positive whole number")
	}
	if d.Withdrawn && len(d.Signs) != 0 {
		return fmt.Errorf("a withdrawn Correction names no sign")
	}
	return nil
}

func NextID(f File) string {
	max := "0"
	for _, d := range f.Declarations {
		if !correctionID.MatchString(d.Correction) {
			continue
		}
		n := d.Correction[1:]
		if len(n) > len(max) || len(n) == len(max) && n > max {
			max = n
		}
	}
	// Increment as decimal text so even an identifier at the integer limit is never reused.
	digits := []byte(max)
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] != '9' {
			digits[i]++
			return "C" + string(digits)
		}
		digits[i] = '0'
	}
	return "C1" + string(digits)
}

func Standing(f File) []Declaration {
	result := make([]Declaration, 0)
	positions := make(map[string]int)
	for _, d := range f.Declarations {
		d.Signs = append([]Sign(nil), d.Signs...)
		if d.Withdrawn {
			d.Signs = nil
		}
		if index, ok := positions[d.Correction]; ok {
			result[index] = d
		} else {
			positions[d.Correction] = len(result)
			result = append(result, d)
		}
	}
	return result
}
func Undeclared(signs []Sign, standing []Declaration) []Sign {
	type key struct{ kind, id string }
	covered := make(map[key]bool)
	for _, d := range standing {
		if !d.Withdrawn {
			for _, s := range d.Signs {
				covered[key{s.Kind, s.ID}] = true
			}
		}
	}
	var result []Sign
	for _, s := range signs {
		if !covered[key{s.Kind, s.ID}] {
			result = append(result, s)
		}
	}
	return result
}

// Field order in these two structs is the public JSON Lines format.
type declarationLine struct {
	Schema     string    `json:"schema"`
	Kind       string    `json:"kind"`
	Time       time.Time `json:"time"`
	Spec       string    `json:"spec"`
	Correction string    `json:"correction"`
	Class      Class     `json:"class"`
	Artifact   Artifact  `json:"artifact"`
	Reason     string    `json:"reason"`
	Signs      []Sign    `json:"signs"`
	Withdrawn  bool      `json:"withdrawn"`
	By         string    `json:"by"`
	Harness    string    `json:"harness"`
	Session    string    `json:"session"`
}
type snapshotLine struct {
	Schema      string    `json:"schema"`
	Kind        string    `json:"kind"`
	Time        time.Time `json:"time"`
	Spec        string    `json:"spec"`
	Head        string    `json:"head"`
	Signs       []Sign    `json:"signs"`
	Dismissed   int       `json:"dismissed"`
	Unavailable string    `json:"unavailable"`
}

func EncodeDeclaration(slug string, d Declaration) ([]byte, error) {
	signs := make([]Sign, 0, len(d.Signs))
	if !d.Withdrawn {
		for _, s := range d.Signs {
			signs = append(signs, Sign{Kind: s.Kind, ID: s.ID})
		}
	}
	b, err := json.Marshal(declarationLine{Schema, "declaration", d.Time.UTC(), slug, d.Correction, d.Class, d.Artifact, strings.TrimSpace(d.Reason), signs, d.Withdrawn, d.By, d.Harness, d.Session})
	if err != nil {
		return nil, fmt.Errorf("encode declaration: %w", err)
	}
	return append(b, '\n'), nil
}
func EncodeSnapshot(slug string, s Snapshot) ([]byte, error) {
	signs := make([]Sign, 0, len(s.Signs))
	for _, sign := range s.Signs {
		if sign.Kind != "commit" {
			sign.Subject = ""
		}
		signs = append(signs, sign)
	}
	b, err := json.Marshal(snapshotLine{Schema, "signs", s.Time.UTC(), slug, s.Head, signs, s.Dismissed, s.Unavailable})
	if err != nil {
		return nil, fmt.Errorf("encode signs: %w", err)
	}
	return append(b, '\n'), nil
}

func Read(path string) (File, error) {
	input, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("open corrections record: %w", err)
	}
	defer input.Close()
	var result File
	reader := bufio.NewReader(input)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) != 0 {
			decodeLine(line, &result)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return result, fmt.Errorf("read corrections record: %w", err)
		}
	}
	return result, nil
}
func decodeLine(line []byte, f *File) {
	var header struct{ Schema, Kind string }
	if json.Unmarshal(line, &header) != nil || header.Schema != Schema {
		f.Skipped++
		return
	}
	switch header.Kind {
	case "declaration":
		var d declarationLine
		if json.Unmarshal(line, &d) != nil {
			f.Skipped++
			return
		}
		f.Declarations = append(f.Declarations, Declaration{d.Time, d.Correction, d.Class, d.Artifact, d.Reason, d.Signs, d.Withdrawn, d.By, d.Harness, d.Session})
	case "signs":
		var s snapshotLine
		if json.Unmarshal(line, &s) != nil {
			f.Skipped++
			return
		}
		f.Snapshots = append(f.Snapshots, Snapshot{s.Time, s.Head, s.Signs, s.Dismissed, s.Unavailable})
	default:
		f.Skipped++
	}
}

func RemoveBefore(ctx context.Context, homeDir string, cutoff time.Time, dryRun bool) (int, error) {
	root := filepath.Join(homeDir, ".roundfix", "corrections")
	repos, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read corrections directory: %w", err)
	}
	count := 0
	var allErr error
	for _, repo := range repos {
		if err := ctx.Err(); err != nil {
			return count, errors.Join(allErr, err)
		}
		if !repo.IsDir() {
			continue
		}
		dir := filepath.Join(root, repo.Name())
		files, err := os.ReadDir(dir)
		if err != nil {
			allErr = errors.Join(allErr, err)
			continue
		}
		for _, entry := range files {
			if err := ctx.Err(); err != nil {
				return count, errors.Join(allErr, err)
			}
			if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".jsonl") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			f, err := Read(path)
			if err != nil {
				allErr = errors.Join(allErr, err)
				continue
			}
			newest := time.Time{}
			for _, d := range f.Declarations {
				if d.Time.After(newest) {
					newest = d.Time
				}
			}
			for _, s := range f.Snapshots {
				if s.Time.After(newest) {
					newest = s.Time
				}
			}
			if len(f.Declarations)+len(f.Snapshots) == 0 {
				info, err := entry.Info()
				if err != nil {
					allErr = errors.Join(allErr, err)
					continue
				}
				newest = info.ModTime()
			}
			if !newest.Before(cutoff) {
				continue
			}
			if !dryRun {
				if err := os.Remove(path); err != nil {
					allErr = errors.Join(allErr, err)
					continue
				}
			}
			count++
		}
		if !dryRun {
			remaining, err := os.ReadDir(dir)
			if err != nil {
				allErr = errors.Join(allErr, err)
			} else if len(remaining) == 0 {
				allErr = errors.Join(allErr, os.Remove(dir))
			}
		}
	}
	return count, allErr
}
