package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestImplementPassesVerificationFormatToTheQAStep(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core"},
		implementQAGateSeed("", "task_01"),
	})
	directory := t.TempDir()
	arguments := filepath.Join(directory, "arguments")
	script := filepath.Join(directory, "formatter.sh")
	body := "#!/bin/sh\nset -eu\nprintf '%s\\n' \"$@\" >> '" + arguments + "'\n"
	if err := os.WriteFile(script, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	writeUserConfig(t, home, fmt.Sprintf("verification:\n  format: %q\n", "sh '"+script+"'"))
	runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, qaReport: implementQAReport("pass")}
	withImplementCollaborators(t, runner)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("implement exit = %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	got, err := os.ReadFile(arguments)
	if err != nil {
		t.Fatal(err)
	}
	report := filepath.ToSlash(runner.qaReportPath)
	if report == "" || !strings.Contains(string(got), report+"\n") {
		t.Fatalf("formatter arguments = %q, missing report %q", got, report)
	}
	for _, path := range strings.Split(strings.TrimSpace(string(got)), "\n") {
		if !strings.HasPrefix(path, "docs/specs/"+implementTestSlug+"/qa/") {
			t.Fatalf("formatter received non-QA file %q", path)
		}
	}
}
