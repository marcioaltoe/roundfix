package spec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSetStatusRewritesOnlyTheStatusValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		statusLine string
		newStatus  Status
		wantLine   string
	}{
		{
			name:       "plain value",
			statusLine: "status: pending",
			newStatus:  StatusCompleted,
			wantLine:   "status: completed",
		},
		{
			name:       "trailing comment preserved",
			statusLine: "status: pending # pending | in_progress | completed | failed",
			newStatus:  StatusInProgress,
			wantLine:   "status: in_progress # pending | in_progress | completed | failed",
		},
		{
			name:       "extra spacing preserved",
			statusLine: "status:   pending",
			newStatus:  StatusFailed,
			wantLine:   "status:   failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := md(`---
task: task_01
spec: demo
` + tt.statusLine + `
type: backend
complexity: low
---

# Task 01: Fixture

## Verification

- 'go test ./...' — expected: pass.
`)
			path := filepath.Join(t.TempDir(), "task_01.md")
			writeFile(t, path, original)

			if err := SetStatus(path, tt.newStatus); err != nil {
				t.Fatalf("SetStatus: %v", err)
			}

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read rewritten file: %v", err)
			}
			want := strings.Replace(original, tt.statusLine, tt.wantLine, 1)
			if string(got) != want {
				t.Errorf("rewritten file = %q, want byte-identical except the status value: %q", got, want)
			}
		})
	}
}

func TestSetStatusRejectsInvalidInput(t *testing.T) {
	t.Parallel()
	valid := md(`---
task: task_01
spec: demo
status: pending
---

# Task 01: Fixture

## Verification

- 'go test ./...' — expected: pass.
`)

	tests := []struct {
		name    string
		content string
		status  Status
		wantMsg string
	}{
		{
			name:    "unsupported status value",
			content: valid,
			status:  Status("done"),
			wantMsg: "not allowed",
		},
		{
			name: "no status field",
			content: md(`---
task: task_01
spec: demo
---

# Task 01: Fixture
`),
			status:  StatusCompleted,
			wantMsg: "no status field",
		},
		{
			name:    "no frontmatter",
			content: "# Task 01: Fixture\n",
			status:  StatusCompleted,
			wantMsg: "frontmatter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "task_01.md")
			writeFile(t, path, tt.content)

			err := SetStatus(path, tt.status)
			if err == nil {
				t.Fatal("SetStatus succeeded, want error")
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error %q does not contain %q", err, tt.wantMsg)
			}

			got, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Fatalf("read file: %v", readErr)
			}
			if string(got) != tt.content {
				t.Errorf("file changed after failed SetStatus")
			}
		})
	}
}

func TestReloadTaskPicksUpAgentEdits(t *testing.T) {
	t.Parallel()
	gitRoot := t.TempDir()
	specsRoot := defaultSpecsRoot(gitRoot)
	relFile := filepath.Join("demo", "task_01.md")
	path := filepath.Join(specsRoot, relFile)
	writeFile(t, path, taskFixture("task_01", "Build the parser", "pending", "backend", defaultVerificationSection))

	task := &Task{ID: "task_01", File: relFile, Needs: []string{"task_00"}}
	if err := ReloadTask(specsRoot, task); err != nil {
		t.Fatalf("ReloadTask: %v", err)
	}
	if task.Status != StatusPending {
		t.Fatalf("Status = %q, want %q", task.Status, StatusPending)
	}

	// Simulate the Agent: settle the status and append a Result section.
	edited := taskFixture("task_01", "Build the parser", "completed", "backend", defaultVerificationSection) + md(`
## Result

- All acceptance criteria verified; 'go test ./...' passed.
`)
	writeFile(t, path, edited)

	if err := ReloadTask(specsRoot, task); err != nil {
		t.Fatalf("ReloadTask after edit: %v", err)
	}
	if task.Status != StatusCompleted {
		t.Errorf("Status = %q, want %q", task.Status, StatusCompleted)
	}
	if task.Title != "Build the parser" {
		t.Errorf("Title = %q, want %q", task.Title, "Build the parser")
	}
	if len(task.Verification) != 1 || task.Verification[0] != "go test ./..." {
		t.Errorf("Verification = %v, want the one fixture command", task.Verification)
	}
	if len(task.Needs) != 1 || task.Needs[0] != "task_00" {
		t.Errorf("Needs = %v, want the manifest-owned value untouched", task.Needs)
	}
}

func TestDerivedQAVerificationRequiresTheNewestReportToPass(t *testing.T) {
	t.Parallel()

	const slug = "derived-qa-verification"
	commands := DerivedQAVerification(slug)
	if len(commands) != 1 {
		t.Fatalf("DerivedQAVerification() = %q, want one command", commands)
	}
	if !strings.Contains(commands[0], filepath.ToSlash(filepath.Join("docs", "specs", slug, "qa"))) {
		t.Fatalf("DerivedQAVerification() = %q, want the Spec's QA directory", commands[0])
	}

	tests := []struct {
		name    string
		reports map[string]string
		wantErr bool
	}{
		{
			name: "only passing report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "pass",
			},
		},
		{
			name: "newer failed rerun defeats older pass",
			reports: map[string]string{
				"qa-report-2026-08-31.md":    "pass",
				"qa-report-2026-08-31-02.md": "fail",
			},
			wantErr: true,
		},
		{
			name: "zero rerun still defeats the unsuffixed report",
			reports: map[string]string{
				"qa-report-2026-08-31.md":    "pass",
				"qa-report-2026-08-31-00.md": "fail",
			},
			wantErr: true,
		},
		{
			name: "numeric rerun order accepts ten after two",
			reports: map[string]string{
				"qa-report-2026-08-31-02.md": "fail",
				"qa-report-2026-08-31-10.md": "pass",
			},
		},
		{
			name: "partial is outside the passing domain",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "partial",
			},
			wantErr: true,
		},
		{
			name: "body text cannot override a failed frontmatter verdict",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "fail\n---\n\nverdict: pass",
			},
			wantErr: true,
		},
		{
			name: "duplicate verdict is not a passing report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "fail\nverdict: pass",
			},
			wantErr: true,
		},
		{
			name: "invalid later date loses to a valid report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "fail",
				"qa-report-2026-13-40.md": "pass",
			},
			wantErr: true,
		},
		{
			name:    "missing report",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			reportDir := filepath.Join(root, "docs", "specs", slug, "qa")
			for name, verdict := range tt.reports {
				writeFile(t, filepath.Join(reportDir, name), "---\nverdict: "+verdict+"\n---\n")
			}

			command := exec.Command("sh", "-c", commands[0])
			command.Dir = root
			output, err := command.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("derived QA Verification error = %v, wantErr %v; output: %s", err, tt.wantErr, output)
			}
		})
	}
}

func TestReloadTaskDerivesOnlyQAVerification(t *testing.T) {
	t.Parallel()

	const slug = "demo"
	derived := DerivedQAVerification(slug)
	tests := []struct {
		name         string
		taskType     string
		verification string
		want         []string
		wantAuthored bool
	}{
		{
			name:         "qa authored command is not effective",
			taskType:     "qa",
			verification: "echo author-controlled",
			want:         derived,
			wantAuthored: true,
		},
		{
			name:         "rendered derived qa command is accepted",
			taskType:     "qa",
			verification: derived[0],
			want:         derived,
		},
		{
			name:         "non-qa command remains authored",
			taskType:     "backend",
			verification: "echo author-controlled",
			want:         []string{"echo author-controlled"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			relFile := filepath.Join(slug, "task_01.md")
			verification := md("## Verification\n\n- '" + tt.verification + "'\n")
			writeFile(t, filepath.Join(specsRoot, relFile), taskFixture("task_01", "Fixture", "pending", tt.taskType, verification))

			task := &Task{ID: "task_01", File: relFile}
			if err := ReloadTask(specsRoot, task); err != nil {
				t.Fatalf("ReloadTask: %v", err)
			}
			if !reflect.DeepEqual(task.Verification, tt.want) {
				t.Errorf("Verification = %q, want %q", task.Verification, tt.want)
			}
			if got := AuthoredQAVerification(*task); got != tt.wantAuthored {
				t.Errorf("AuthoredQAVerification() = %v, want %v", got, tt.wantAuthored)
			}
		})
	}
}

func TestNegativeControlSectionParsesInOrder(t *testing.T) {
	t.Parallel()

	document, err := parseTaskDocument([]byte(taskFixture(
		"task_01",
		"Carry the negative control",
		"pending",
		"backend",
		md(`## Negative Control

- 'go test ./internal/spec -run TestRejectsKnownDefect' — expected: exits non-zero.
- a declaration without a backticked command is skipped
- 'go test ./internal/spec -run TestRejectsMissingSection' — expected: exits non-zero.

`)+defaultVerificationSection,
	)), "task_01.md")
	if err != nil {
		t.Fatalf("parseTaskDocument: %v", err)
	}

	want := []string{
		"go test ./internal/spec -run TestRejectsKnownDefect",
		"go test ./internal/spec -run TestRejectsMissingSection",
	}
	if !reflect.DeepEqual(document.NegativeControl, want) {
		t.Fatalf("NegativeControl = %q, want declarations in source order %q", document.NegativeControl, want)
	}
}

func TestNegativeControlAbsentSectionParsesEmpty(t *testing.T) {
	t.Parallel()

	document, err := parseTaskDocument([]byte(taskFixture(
		"task_01",
		"Carry no negative control",
		"pending",
		"backend",
		defaultVerificationSection,
	)), "task_01.md")
	if err != nil {
		t.Fatalf("parseTaskDocument: %v", err)
	}
	if len(document.NegativeControl) != 0 {
		t.Fatalf("NegativeControl = %q, want an empty declaration list", document.NegativeControl)
	}
}

func TestNegativeControlParsingPreservesVerificationCommands(t *testing.T) {
	t.Parallel()

	document, err := parseTaskDocument([]byte(taskFixture(
		"task_01",
		"Keep verification stable",
		"pending",
		"backend",
		md(`## Negative Control

- 'go test ./internal/spec -run TestRejectsKnownDefect' — expected: exits non-zero.

## Verification

- 'go test ./internal/spec/' — expected: all tests pass.
- a bullet without a command is skipped
- 'go build ./...' — expected: builds cleanly.

## References

- 'go vet ./...' outside the Verification section is not a command.
`),
	)), "task_01.md")
	if err != nil {
		t.Fatalf("parseTaskDocument: %v", err)
	}

	want := []string{"go test ./internal/spec/", "go build ./..."}
	if !reflect.DeepEqual(document.Verification, want) {
		t.Fatalf("Verification = %q, want unchanged commands %q", document.Verification, want)
	}
}

func TestParseTaskDocumentDeclarations(t *testing.T) {
	t.Parallel()

	document, err := parseTaskDocument([]byte(md(`---
task: task_01
spec: demo
status: pending
type: backend
complexity: low
---

# Task 01: Rehearse the gate

## Requirements

1. MUST keep the named gate enabled across
   the full rehearsal.
2. MUST NOT keep the named gate enabled.

### Notes

This nested section is not part of the requirement declaration.

## Rehearsal Cases

- Case: contradictory requirements; Observation: spec check reports the contradiction.
- Case: declared cases; Observation: the focused test records the result.

## Verification

- 'go test ./internal/speccheck' — expected: exit 0.
`)), "task_01.md")
	if err != nil {
		t.Fatalf("parseTaskDocument: %v", err)
	}

	wantRequirements := []TaskDeclaration{
		{Text: "MUST keep the named gate enabled across the full rehearsal.", Line: 13},
		{Text: "MUST NOT keep the named gate enabled.", Line: 15},
	}
	if !reflect.DeepEqual(document.Requirements, wantRequirements) {
		t.Errorf("Requirements = %#v, want %#v", document.Requirements, wantRequirements)
	}
	wantCases := []TaskDeclaration{
		{Text: "Case: contradictory requirements; Observation: spec check reports the contradiction.", Line: 23},
		{Text: "Case: declared cases; Observation: the focused test records the result.", Line: 24},
	}
	if !reflect.DeepEqual(document.RehearsalCases, wantCases) {
		t.Errorf("RehearsalCases = %#v, want %#v", document.RehearsalCases, wantCases)
	}
	if document.TitleLine != 9 {
		t.Errorf("TitleLine = %d, want 9", document.TitleLine)
	}
}

func TestReloadTaskNormalizesStatusValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		rawStatus      string
		wantStatus     Status
		wantNormalized bool
	}{
		{name: "pending canonical", rawStatus: "pending", wantStatus: StatusPending},
		{name: "in progress canonical", rawStatus: "in_progress", wantStatus: StatusInProgress},
		{name: "completed canonical", rawStatus: "completed", wantStatus: StatusCompleted},
		{name: "failed canonical", rawStatus: "failed", wantStatus: StatusFailed},
		{name: "done synonym", rawStatus: "done", wantStatus: StatusCompleted, wantNormalized: true},
		{name: "hyphen in progress synonym", rawStatus: "in-progress", wantStatus: StatusInProgress, wantNormalized: true},
		{name: "space in progress synonym", rawStatus: "in progress", wantStatus: StatusInProgress, wantNormalized: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			relFile := filepath.Join("demo", "task_01.md")
			writeFile(t, filepath.Join(specsRoot, relFile), taskFixture("task_01", "Fixture", tt.rawStatus, "backend", defaultVerificationSection))

			task := &Task{ID: "task_01", File: relFile}
			if err := ReloadTask(specsRoot, task); err != nil {
				t.Fatalf("ReloadTask: %v", err)
			}
			if task.Status != tt.wantStatus {
				t.Fatalf("Status = %q, want %q", task.Status, tt.wantStatus)
			}
			if task.StatusNormalized != tt.wantNormalized {
				t.Fatalf("StatusNormalized = %v, want %v", task.StatusNormalized, tt.wantNormalized)
			}
		})
	}
}

func TestReloadTaskReportsBrokenAgentEdits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, err error)
	}{
		{
			name:    "unsupported status written by the Agent",
			content: taskFixture("task_01", "Fixture", "finished", "backend", defaultVerificationSection),
			check: func(t *testing.T, err error) {
				var taskErr TaskFileError
				if !errors.As(err, &taskErr) {
					t.Fatalf("error = %v, want TaskFileError", err)
				}
				if taskErr.TaskID != "task_01" {
					t.Errorf("TaskID = %q, want %q", taskErr.TaskID, "task_01")
				}
				if !strings.Contains(err.Error(), `unsupported status "finished"`) {
					t.Errorf("error = %q, want unsupported status diagnostic", err)
				}
				if !strings.Contains(err.Error(), "pending, in_progress, completed, failed") {
					t.Errorf("error = %q, want allowed statuses", err)
				}
			},
		},
		{
			name:    "Verification section removed",
			content: taskFixture("task_01", "Fixture", "completed", "backend", "## Result\n\nDone.\n"),
			check: func(t *testing.T, err error) {
				var missing MissingVerificationError
				if !errors.As(err, &missing) {
					t.Fatalf("error = %v, want MissingVerificationError", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRoot := t.TempDir()
			specsRoot := defaultSpecsRoot(gitRoot)
			relFile := filepath.Join("demo", "task_01.md")
			writeFile(t, filepath.Join(specsRoot, relFile), tt.content)

			task := &Task{ID: "task_01", File: relFile}
			err := ReloadTask(specsRoot, task)
			if err == nil {
				t.Fatal("ReloadTask succeeded, want error")
			}
			tt.check(t, err)
		})
	}
}
