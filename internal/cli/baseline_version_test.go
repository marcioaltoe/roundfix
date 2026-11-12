// Suite: Baseline Recorded Roundfix Version commands
// Invariant: newer Baseline records refuse before repository work and writer releases converge.
// Boundary IN: command streams, exit codes, manifest bytes, and repository tree identities.
// Boundary OUT: live release lookup and skill acquisition.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/app"
)

func TestBaselineUpdateRecordsTheRunningRoundfixVersion(t *testing.T) {
	t.Parallel()
	repository := staleBaselineUpdateRepository(t)
	manifest := ReadBaselineSetupManifest(t, repository)
	manifest.RoundfixVersion = ""
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	path := filepath.Join(repository, baselineSetupManifestPath)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(before, []byte(`"roundfixVersion"`)) {
		t.Fatal("fixture already records a writer")
	}
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--yes", "--format", "json")
	if code != exitOK || result.State != "verified" || stderr != "" {
		t.Fatalf("exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	expected := `"roundfixVersion": "` + app.NormalizeVersion(app.Version) + `"`
	if !bytes.Contains(after, []byte(expected)) {
		t.Fatalf("manifest omits %s: %s", expected, after)
	}
	result, stdout, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--yes", "--format", "json")
	if code != exitOK || len(result.FileChanges) != 0 || result.State != "verified" || stderr != "" {
		t.Fatalf("second update exit=%d result=%+v stdout=%s stderr=%s", code, result, stdout, stderr)
	}
}

func TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                      string
		args                      []string
		category, action, finding string
	}{
		{name: "adopt", args: []string{"baseline"}, category: "preflight", action: "run roundfix upgrade, then rerun roundfix baseline"},
		{name: "update", args: []string{"baseline", "update"}, category: "version", action: "run roundfix upgrade, then rerun roundfix baseline update"},
		{name: "plan", args: []string{"baseline", "plan"}},
		{name: "apply", args: []string{"baseline", "apply", "--plan", "missing-plan.json", "--confirm-plan", "sha256:" + strings.Repeat("b", 64)}},
		{name: "skills restore", args: []string{"baseline", "skills", "restore", "--profile", "go-cli-tui"}, finding: "restore.roundfix-version-newer", action: "Run roundfix upgrade, then rerun roundfix baseline skills restore."},
		{name: "skills reconcile", args: []string{"baseline", "skills", "reconcile", "--profile", "go-cli-tui", "--source", "example/skills", "--revision", strings.Repeat("a", 40)}, finding: "reconcile.roundfix-version-newer", action: "Run roundfix upgrade, then rerun roundfix baseline skills reconcile."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repository, refusal := newRecordedVersionRepository(t)
			before := baselinePlanTestTree(t, repository)
			args := append(append([]string(nil), tc.args...), "--repo", repository, "--format", "json")
			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
			escaped, err := json.Marshal(refusal)
			if err != nil {
				t.Fatal(err)
			}
			if code != exitPreflight || !bytes.Contains(stdout.Bytes(), escaped) || strings.Contains(stdout.String()+stderr.String(), "future_section") {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			var result struct {
				State      string                                  `json:"state"`
				Category   string                                  `json:"category"`
				NextAction string                                  `json:"nextAction"`
				Finding    *struct{ Code, Message, Action string } `json:"finding"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if tc.category != "" && result.Category != tc.category {
				t.Fatalf("category=%q want=%q", result.Category, tc.category)
			}
			if tc.name == "update" && result.State != "failed" {
				t.Fatalf("state=%q", result.State)
			}
			if tc.finding != "" {
				if result.Finding == nil || result.Finding.Code != tc.finding || result.Finding.Message != refusal || result.Finding.Action != tc.action {
					t.Fatalf("finding=%+v", result.Finding)
				}
			} else if tc.action != "" && result.NextAction != tc.action {
				t.Fatalf("action=%q want=%q", result.NextAction, tc.action)
			}
			if before != baselinePlanTestTree(t, repository) {
				t.Fatal("refused command changed repository tree")
			}
		})
	}
}

func TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix(t *testing.T) {
	t.Parallel()
	repository, refusal := newRecordedVersionRepository(t)
	configureUpdateCheckTest(t, "0.0.0-dev", func(context.Context) ([]app.Release, error) {
		t.Fatal("development build reached release transport")
		return nil, nil
	})
	before := baselinePlanTestTree(t, repository)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"update", "check"}, &stdout, &stderr)
	wantOut := "Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify).\nBaseline: failed: " + refusal + "\n"
	wantErr := "roundfix: warning: " + refusal + "; Project Config not read\n"
	if code != exitOK || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("Surface Transcript 2: exit=%d stdout=%q stderr=%q; want stdout=%q stderr=%q", code, &stdout, &stderr, wantOut, wantErr)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, t.Context(), []string{"update", "check", "--format", "json"}, &stdout, &stderr)
	var report updateCheckReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode report: %v stdout=%s", err, &stdout)
	}
	if code != exitOK || stderr.String() != wantErr || report.Baseline.State != "failed" || report.Baseline.Message != refusal || report.Next != "none" {
		t.Fatalf("exit=%d report=%+v stderr=%q", code, report, &stderr)
	}
	if before != baselinePlanTestTree(t, repository) {
		t.Fatal("Update Check changed repository tree")
	}
}
