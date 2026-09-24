package spec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
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

func TestReopenGateWritesToTheValidatedTarget(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	validatedTarget := filepath.Join(root, "targets", "validated", "task_qa.md")
	writeFile(t, validatedTarget, md(`---
task: task_qa
spec: demo
status: completed
type: qa
complexity: low
---

# Task QA: Verify the feature

## Verification

- 'go test ./...' — expected: pass.
`))
	unvalidatedTarget := filepath.Join(root, "targets", "unvalidated", "task_qa.md")
	const unvalidatedContent = "unvalidated target must stay unchanged\n"
	writeFile(t, unvalidatedTarget, unvalidatedContent)
	taskPath := filepath.Join(root, "demo", "task_qa.md")
	if err := os.MkdirAll(filepath.Dir(taskPath), 0o755); err != nil {
		t.Fatalf("create Task directory: %v", err)
	}
	if err := os.Symlink(validatedTarget, taskPath); err != nil {
		t.Fatalf("link Task path to target: %v", err)
	}
	resolvedTarget, err := filepath.EvalSymlinks(taskPath)
	if err != nil {
		t.Fatalf("resolve validated Task target: %v", err)
	}
	if err := os.Remove(taskPath); err != nil {
		t.Fatalf("remove original Task symlink: %v", err)
	}
	if err := os.Symlink(unvalidatedTarget, taskPath); err != nil {
		t.Fatalf("redirect Task symlink after validation: %v", err)
	}

	if err := ReopenGate(resolvedTarget, "qa/qa-report-2026-09-19.md", []string{"task_01"}, time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("ReopenGate: %v", err)
	}

	targetBytes, err := os.ReadFile(validatedTarget)
	if err != nil {
		t.Fatalf("read validated target Task: %v", err)
	}
	target := string(targetBytes)
	if !strings.Contains(target, "status: pending") {
		t.Fatalf("validated target status was not rewritten to pending:\n%s", target)
	}
	if !strings.Contains(target, "- QA Report: `qa/qa-report-2026-09-19.md`") {
		t.Fatalf("target does not record the invalidated QA Report:\n%s", target)
	}
	if !strings.Contains(target, "- Dependencies not completed: `task_01`") {
		t.Fatalf("target does not record the stale dependency:\n%s", target)
	}
	if got, err := os.ReadFile(unvalidatedTarget); err != nil {
		t.Fatalf("read unvalidated target Task: %v", err)
	} else if string(got) != unvalidatedContent {
		t.Fatalf("unvalidated target changed:\n%s", got)
	}
	info, err := os.Lstat(taskPath)
	if err != nil {
		t.Fatalf("lstat Task path: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("Task path mode = %v, want symlink", info.Mode())
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

func TestDerivedQAVerificationProvesNewestVerdictReadable(t *testing.T) {
	const slug = "derived-qa-verification"

	tests := []struct {
		name    string
		reports map[string]string
		wantErr bool
	}{
		{
			name: "only passing report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": derivedQAReportFixture(VerdictPass),
			},
		},
		{
			name: "failing verdict remains readable",
			reports: map[string]string{
				"qa-report-2026-08-31.md": derivedQAReportFixture(VerdictFail),
			},
		},
		{
			name: "partial verdict remains readable",
			reports: map[string]string{
				"qa-report-2026-08-31.md": derivedQAReportFixture(VerdictPartial, "rows_blocked_declared: 1"),
			},
		},
		{
			name: "newer unreadable rerun defeats older readable report",
			reports: map[string]string{
				"qa-report-2026-08-31.md":    derivedQAReportFixture(VerdictPass),
				"qa-report-2026-08-31-02.md": "not frontmatter\n",
			},
			wantErr: true,
		},
		{
			name: "numeric rerun order reads ten after two",
			reports: map[string]string{
				"qa-report-2026-08-31-02.md": "not frontmatter\n",
				"qa-report-2026-08-31-10.md": derivedQAReportFixture(VerdictPass),
			},
		},
		{
			name: "malformed rerun loses to readable unsuffixed report",
			reports: map[string]string{
				"qa-report-2026-08-31.md":         derivedQAReportFixture(VerdictPass),
				"qa-report-2026-08-31-ffd6852.md": "not frontmatter\n",
			},
		},
		{
			name: "duplicate verdict is unreadable",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "---\nverdict: pass\nverdict: fail\n---\n",
			},
			wantErr: true,
		},
		{
			name: "body verdict does not replace missing frontmatter verdict",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "---\n---\n\nverdict: pass\n",
			},
			wantErr: true,
		},
		{
			name: "empty verdict is unreadable",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "---\nverdict:   \n---\n",
			},
			wantErr: true,
		},
		{
			name: "unparseable report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": "not frontmatter\n",
			},
			wantErr: true,
		},
		{
			name: "invalid later date loses to readable valid report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": derivedQAReportFixture(VerdictPass),
				"qa-report-2026-13-40.md": "not frontmatter\n",
			},
		},
		{
			name: "non-digit date loses to readable valid report",
			reports: map[string]string{
				"qa-report-2026-08-31.md": derivedQAReportFixture(VerdictPass),
				"qa-report-202A-08-31.md": "not frontmatter\n",
			},
		},
		{
			// The sole report is malformed. Sorting a malformed name last only
			// helps while a well-formed one exists beside it; alone, it won
			// tail -1 and its verdict was read as the gate's.
			name: "a sole malformed report is not a candidate",
			reports: map[string]string{
				"qa-report-notadate.md": derivedQAReportFixture(VerdictPass),
			},
			wantErr: true,
		},
		{
			name: "a sole report with a malformed sequence is not a candidate",
			reports: map[string]string{
				"qa-report-2026-08-31-x2.md": derivedQAReportFixture(VerdictPass),
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
			root, command := derivedQAVerificationFixture(t, slug)
			reportDir := filepath.Join(root, "docs", "specs", slug, "qa")
			for name, report := range tt.reports {
				writeFile(t, filepath.Join(reportDir, name), report)
			}

			output, err := command.CombinedOutput()
			if (err != nil) != tt.wantErr {
				t.Fatalf("derived QA Verification error = %v, wantErr %v; output: %s", err, tt.wantErr, output)
			}
		})
	}
}

func TestDerivedQAVerificationInvokesNoRoundfixBinary(t *testing.T) {
	const slug = "derived-qa-verification"
	commands := DerivedQAVerification(slug)
	if strings.Contains(commands[0], "roundfix") {
		t.Fatalf("DerivedQAVerification() invokes roundfix: %q", commands[0])
	}
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "docs", "specs", slug, "qa", "qa-report-2026-08-31.md"), derivedQAReportFixture(VerdictPass))

	selectorPath := t.TempDir()
	for _, name := range []string{"awk", "cut", "find", "sort", "tail"} {
		target, err := exec.LookPath(name)
		if err != nil {
			t.Fatalf("locate %s for derived QA Verification: %v", name, err)
		}
		if err := os.Symlink(target, filepath.Join(selectorPath, name)); err != nil {
			t.Fatalf("link %s for derived QA Verification: %v", name, err)
		}
	}

	command := exec.Command("sh", "-c", commands[0])
	command.Dir = root
	command.Env = append(os.Environ(), "PATH="+selectorPath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("derived QA Verification required roundfix: %v; output: %s", err, output)
	}
}

func derivedQAVerificationFixture(t *testing.T, slug string) (string, *exec.Cmd) {
	t.Helper()
	commands := DerivedQAVerification(slug)
	if len(commands) != 1 {
		t.Fatalf("DerivedQAVerification() = %q, want one command", commands)
	}
	reportDir := filepath.ToSlash(filepath.Join("docs", "specs", slug, "qa"))
	if !strings.Contains(commands[0], reportDir) {
		t.Fatalf("DerivedQAVerification() = %q, want the Spec's QA directory", commands[0])
	}

	root := t.TempDir()
	command := exec.Command("sh", "-c", commands[0])
	command.Dir = root
	command.Env = os.Environ()
	return root, command
}

func derivedQAReportFixture(verdict string, extraFrontmatter ...string) string {
	frontmatter := ""
	if len(extraFrontmatter) > 0 {
		frontmatter = strings.Join(extraFrontmatter, "\n") + "\n"
	}
	return "---\nverdict: " + verdict + "\n" + frontmatter + "---\n"
}

func TestDerivedQAVerificationPassesTheChecker(t *testing.T) {
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(testFile), "..", ".."))

	// internal/speccheck imports internal/spec, so run the checker-side
	// integration test in a subprocess instead of introducing an import cycle.
	command := exec.Command(
		"go", "test", "-count=1", "./internal/speccheck",
		"-run", "^TestDerivedQAVerificationIsAcceptedByTaskStageChecker$",
	)
	command.Dir = repoRoot
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("derived QA Verification checker contract: %v\n%s", err, output)
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
			// Built here rather than through taskFixture, which routes its
			// whole body through md() and rewrites every single quote into a
			// backtick. The derived command carries single quotes now that the
			// Spec path is shell-quoted, so that helper would corrupt the very
			// command under test.
			content := "---\ntask: task_01\nspec: demo\nstatus: pending\ntype: " + tt.taskType +
				"\ncomplexity: low\n---\n\n# Task 01: Fixture\n\n## Overview\n\nFixture task.\n\n" +
				"## Verification\n\n- `" + tt.verification + "`\n"
			writeFile(t, filepath.Join(specsRoot, relFile), content)

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

func TestTaskParsesIndependentVerification(t *testing.T) {
	t.Parallel()

	t.Run("records the declaration", func(t *testing.T) {
		t.Parallel()

		gitRoot := t.TempDir()
		specsRoot := defaultSpecsRoot(gitRoot)
		relFile := filepath.Join("demo", "task_01.md")
		content := strings.Replace(
			taskFixture("task_01", "Collect every failure", "pending", "backend", defaultVerificationSection),
			"complexity: low\n",
			"complexity: low\nverification: independent\n",
			1,
		)
		writeFile(t, filepath.Join(specsRoot, relFile), content)

		task := &Task{ID: "task_01", File: relFile}
		if err := ReloadTask(specsRoot, task); err != nil {
			t.Fatalf("ReloadTask: %v", err)
		}
		if task.VerificationMode != VerificationModeIndependent {
			t.Fatalf("VerificationMode = %q, want %q", task.VerificationMode, VerificationModeIndependent)
		}
	})

	t.Run("refuses an unknown value by name", func(t *testing.T) {
		t.Parallel()

		gitRoot := t.TempDir()
		specsRoot := defaultSpecsRoot(gitRoot)
		relFile := filepath.Join("demo", "task_01.md")
		content := strings.Replace(
			taskFixture("task_01", "Reject an unknown mode", "pending", "backend", defaultVerificationSection),
			"complexity: low\n",
			"complexity: low\nverification: sequential-groups\n",
			1,
		)
		writeFile(t, filepath.Join(specsRoot, relFile), content)

		err := ReloadTask(specsRoot, &Task{ID: "task_01", File: relFile})
		if err == nil {
			t.Fatal("ReloadTask succeeded, want an unsupported verification error")
		}
		if !strings.Contains(err.Error(), `unsupported verification "sequential-groups"`) {
			t.Fatalf("ReloadTask error = %q, want the unknown value named", err)
		}
	})
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

func TestDerivedQAVerificationQuotesTheSpecPath(t *testing.T) {
	t.Parallel()
	// The slug reaches this command from a directory name, not from a validated
	// identifier, and the command runs under sh -c. An unquoted slug carrying a
	// separator or a quote would change which command runs.
	for _, testCase := range []struct {
		name string
		slug string
	}{
		{name: "command separator", slug: "demo; touch owned"},
		{name: "single quote", slug: "demo'x"},
		{name: "command substitution", slug: "demo$(touch owned)"},
		{name: "backtick substitution", slug: "demo`touch owned`"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			command := DerivedQAVerification(testCase.slug)[0]
			quoted := shellSingleQuoted(filepath.ToSlash(filepath.Join("docs", "specs", testCase.slug, "qa")))
			if !strings.Contains(command, quoted) {
				t.Fatalf("derived command does not carry the quoted path %q", quoted)
			}
			if strings.Contains(command, "find "+filepath.ToSlash(filepath.Join("docs", "specs", testCase.slug, "qa"))+" ") {
				t.Fatalf("derived command interpolates the slug unquoted: %q", command)
			}
		})
	}
}

func TestShellSingleQuotedSurvivesEmbeddedQuotes(t *testing.T) {
	t.Parallel()
	if got := shellSingleQuoted("plain"); got != "'plain'" {
		t.Fatalf("shellSingleQuoted(plain) = %q", got)
	}
	// A single quote closes the word, escapes, and reopens it, so the shell
	// still reads one argument.
	if got := shellSingleQuoted("a'b"); got != `'a'\''b'` {
		t.Fatalf("shellSingleQuoted(a'b) = %q", got)
	}
}
