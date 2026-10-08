package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func probeProvenanceCLI(t *testing.T, format string) (int, string, specCheckDocument) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean", "--run-verification", "--format", format}, &stdout, &stderr)
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
	var document specCheckDocument
	if format == "json" {
		if err := json.Unmarshal(stdout.Bytes(), &document); err != nil {
			t.Fatal(err)
		}
	}
	return code, stdout.String(), document
}

func TestSpecCheckReportsAMalformedCommand(t *testing.T) {
	t.Parallel()
	marker := filepath.Join(t.TempDir(), "ran")
	_, _ = newSpecCheckVerificationWorkspace(t, []string{fmt.Sprintf("touch %q; if", marker)})
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			code, output, doc := probeProvenanceCLI(t, format)
			if code != exitRunFailed {
				t.Fatalf("exit = %d; output=%s", code, output)
			}
			if format == "text" {
				if !strings.Contains(output, "task_01: malformed") || !strings.Contains(strings.ToLower(output), "syntax error") {
					t.Fatalf("output = %s", output)
				}
			} else if len(doc.Verification.Commands) != 1 || doc.Verification.Commands[0].Verdict != "malformed" || !strings.Contains(strings.ToLower(doc.Verification.Commands[0].Cause), "syntax error") {
				t.Fatalf("report = %+v", doc.Verification)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("malformed command ran: %v", err)
			}
		})
	}
}

func TestSpecCheckNamesAnUncommittedVerificationSource(t *testing.T) {
	t.Parallel()
	_, repo := newSpecCheckVerificationWorkspace(t, []string{"test -f task-output.txt"})
	dir := filepath.Join(repo, "docs/specs/clean")
	manifest := filepath.Join(dir, "_tasks.md")
	// Keep the modified manifest and make its Task an untracked source. The
	// commands must still execute, despite both provenance states.
	mustWrite(t, manifest, mustRead(t, manifest)+"\nModified graph prose.\n")
	gitImplement(t, repo, "rm", "--cached", "docs/specs/clean/task_01.md")
	gitImplement(t, repo, "commit", "-m", "test: untrack task source")
	marker := filepath.Join(t.TempDir(), "ran")
	taskPath := filepath.Join(dir, "task_01.md")
	mustWrite(t, taskPath, strings.Replace(mustRead(t, taskPath), "test -f task-output.txt", fmt.Sprintf("touch %q; false", marker), 1))
	mustWrite(t, filepath.Join(dir, "unrelated.md"), "unrelated\n")
	want := []specCheckVerificationSource{{Path: "docs/specs/clean/_tasks.md", State: "modified"}, {Path: "docs/specs/clean/task_01.md", State: "untracked"}}
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			code, output, doc := probeProvenanceCLI(t, format)
			if code != exitOK {
				t.Fatalf("exit=%d output=%s", code, output)
			}
			if format == "json" {
				if !reflect.DeepEqual(doc.Verification.Uncommitted, want) {
					t.Fatalf("sources = %+v", doc.Verification.Uncommitted)
				}
			} else {
				expected := "Verification tree: HEAD\nUncommitted Verification source: docs/specs/clean/_tasks.md (modified)\nUncommitted Verification source: docs/specs/clean/task_01.md (untracked)\n"
				if !strings.Contains(output, expected) || strings.Count(output, "Uncommitted Verification source:") != 2 {
					t.Fatalf("output=%s", output)
				}
			}
			if _, err := os.Stat(marker); err != nil {
				t.Fatalf("uncommitted command was not run: %v", err)
			}
		})
	}
}

func TestSpecCheckNamesNoSourceWhenCommitted(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckVerificationWorkspace(t, []string{"test -f task-output.txt"})
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			code, output, doc := probeProvenanceCLI(t, format)
			if code != exitOK {
				t.Fatalf("exit=%d output=%s", code, output)
			}
			if format == "json" {
				if doc.Verification.Uncommitted == nil || len(doc.Verification.Uncommitted) != 0 {
					t.Fatalf("sources=%+v", doc.Verification.Uncommitted)
				}
			} else if strings.Contains(output, "Uncommitted Verification source:") {
				t.Fatalf("output=%s", output)
			}
		})
	}
}
