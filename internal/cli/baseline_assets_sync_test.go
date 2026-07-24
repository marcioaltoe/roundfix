// Suite: public Baseline asset synchronization CLI
// Invariant: asset synchronization exposes one non-interactive text/JSON contract with results on stdout, diagnostics on stderr, and stable exit categories.
// Boundary IN: public dispatch, flag parsing, help, request injection, rendering, and exit mapping.
// Boundary OUT: source proof and recoverable refresh, owned by internal/baseline/assets_sync_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineAssetsSyncCommand(t *testing.T) {
	t.Run("public help exposes the approved contract", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "assets", "sync", "--help",
		}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		for _, want := range []string{
			"--source-dir",
			"--check",
			"--format",
			"Exit codes:",
			"read-only",
			"recoverable Baseline transaction",
		} {
			if !strings.Contains(stdout.String(), want) {
				t.Fatalf("help missing %q:\n%s", want, stdout.String())
			}
		}
	})

	t.Run("JSON drift stays on stdout and exits execution failure", func(t *testing.T) {
		var captured baseline.AssetsSyncRequest
		execute := func(
			_ context.Context,
			request baseline.AssetsSyncRequest,
		) (baseline.AssetsSyncPayload, error) {
			captured = request
			finding := baseline.AssetsSyncFinding{
				Code:      "skills.setup-snapshot.drift",
				Severity:  "error",
				Path:      "assets/setups/go-cli.json",
				ManagedID: "setup.go-cli",
				Message:   "Bundled setup snapshot differs from the canonical source.",
				Action:    "Run roundfix baseline assets sync without --check.",
			}
			return baseline.AssetsSyncPayload{
					SchemaVersion:  baseline.AssetsSyncSchemaVersion,
					OK:             false,
					Summary:        baseline.AssetsSyncSummary{Errors: 1},
					Findings:       []baseline.AssetsSyncFinding{finding},
					PlannedChanges: []baseline.AssetsSyncChange{},
				}, &baseline.AssetsSyncError{
					Category: baseline.AssetsSyncExecution,
					Finding:  finding,
					Err:      errors.New("snapshot drift"),
				}
		}
		source := t.TempDir()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runBaselineAssetsSyncCommandWith(
			context.Background(),
			[]string{"--source-dir", source, "--check", "--format=json"},
			&stdout,
			&stderr,
			execute,
		)
		if code != exitRunFailed {
			t.Fatalf("drift exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "baseline assets sync failed") {
			t.Fatalf("drift stderr = %q", stderr.String())
		}
		var payload baseline.AssetsSyncPayload
		if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
			t.Fatalf("decode drift JSON: %v\n%s", err, stdout.String())
		}
		if len(payload.Findings) != 1 || payload.Findings[0].Code != "skills.setup-snapshot.drift" {
			t.Fatalf("drift payload = %+v", payload)
		}
		wantSource, err := filepath.Abs(source)
		if err != nil {
			t.Fatal(err)
		}
		if captured.SourceDir != wantSource || !captured.Check {
			t.Fatalf("captured request = %+v", captured)
		}
	})

	t.Run("text success is diagnostic free", func(t *testing.T) {
		execute := func(
			_ context.Context,
			_ baseline.AssetsSyncRequest,
		) (baseline.AssetsSyncPayload, error) {
			return baseline.AssetsSyncPayload{
				SchemaVersion: baseline.AssetsSyncSchemaVersion,
				OK:            true,
				Summary:       baseline.AssetsSyncSummary{Info: 1},
				Findings: []baseline.AssetsSyncFinding{{
					Code:      "skills.setup-snapshot.updated",
					Severity:  "info",
					Path:      "assets/setups/go-cli.json",
					ManagedID: "setup.go-cli",
					Message:   "Setup snapshot was synchronized from the canonical source.",
					Action:    "Review the snapshot diff and run asset validation.",
				}},
				PlannedChanges: []baseline.AssetsSyncChange{},
			}, nil
		}
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runBaselineAssetsSyncCommandWith(
			context.Background(),
			[]string{"--source-dir", "."},
			&stdout,
			&stderr,
			execute,
		)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("success exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		for _, want := range []string{
			"setup-context-driven audit: findings",
			"errors=0 decisions=0 warnings=0 info=1",
			"skills.setup-snapshot.updated",
		} {
			if !strings.Contains(stdout.String(), want) {
				t.Fatalf("text result missing %q:\n%s", want, stdout.String())
			}
		}
	})
}
