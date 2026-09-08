// Suite: public Baseline skills restoration CLI
// Invariant: restoration exposes one deterministic preview/confirm contract with results on stdout, diagnostics on stderr, and stable exit categories.
// Boundary IN: public dispatch, flag parsing, help, request injection, text/JSON rendering, and exit mapping.
// Boundary OUT: immutable Git acquisition and recoverable mutation, which internal/baseline/skills_restore_test.go owns.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineSkillsRestoreCommand(t *testing.T) {
	t.Parallel()
	t.Run("public help exposes the approved contract", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "skills", "restore", "--help",
		}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		for _, want := range []string{
			"--repo",
			"--profile",
			"--skill",
			"--source-dir",
			"--confirm-plan",
			"--format",
			"Exit codes:",
			"recoverable Baseline",
			"transaction to update",
		} {
			if !strings.Contains(stdout.String(), want) {
				t.Fatalf("help missing %q:\n%s", want, stdout.String())
			}
		}
	})

	t.Run("JSON preview stays on stdout and exits action required", func(t *testing.T) {
		digest := strings.Repeat("a", 64)
		var captured baseline.SkillsRestoreRequest
		execute := func(
			_ context.Context,
			request baseline.SkillsRestoreRequest,
		) (baseline.SkillsRestorePayload, error) {
			captured = request
			finding := baseline.RestoreFinding{
				Code:    "plan.confirmation.required",
				Message: "Restoration requires confirmation of this exact Change Plan.",
				Action:  "Review plannedChanges and rerun with --confirm-plan planDigest.",
			}
			return baseline.SkillsRestorePayload{
					SchemaVersion: baseline.SkillsRestoreSchemaVersion,
					Profile:       "rust-cli",
					Setup:         stringPointerForCLI("rust-cli"),
					Acquisitions: []baseline.RestoreAcquisition{{
						Provider: "github", Repository: "example/skills", Ref: strings.Repeat("b", 40),
					}},
					Skills: []baseline.RestoreSkill{},
					PlannedChanges: []baseline.RestorePlannedChange{{
						Action: "create", Path: ".agents/skills/example/SKILL.md", Skill: "example",
					}},
					PlanDigest: &digest,
					Finding:    &finding,
				}, &baseline.SkillsRestoreError{
					Category: baseline.SkillsRestoreAction,
					Finding:  finding,
					Err:      errors.New("confirmation required"),
				}
		}
		source := t.TempDir()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runBaselineSkillsRestoreCommandWith(
			context.Background(),
			[]string{
				"--repo", ".",
				"--profile", "rust-cli",
				"--skill", "example",
				"--source-dir", source,
				"--format=json",
			},
			&stdout,
			&stderr,
			execute,
		)
		if code != exitUnverified {
			t.Fatalf("preview exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "baseline skills restore failed") {
			t.Fatalf("preview stderr = %q", stderr.String())
		}
		var payload baseline.SkillsRestorePayload
		if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
			t.Fatalf("decode preview JSON: %v\n%s", err, stdout.String())
		}
		if payload.Finding == nil || payload.Finding.Code != "plan.confirmation.required" ||
			payload.PlanDigest == nil || *payload.PlanDigest != digest {
			t.Fatalf("preview payload = %+v", payload)
		}
		wantSource, err := filepath.Abs(source)
		if err != nil {
			t.Fatal(err)
		}
		if captured.ProfileID != "rust-cli" ||
			!reflect.DeepEqual(captured.Skills, []string{"example"}) ||
			captured.SourceDir != wantSource {
			t.Fatalf("captured request = %+v", captured)
		}
	})

	t.Run("text success is diagnostic free", func(t *testing.T) {
		digest := strings.Repeat("c", 64)
		execute := func(
			_ context.Context,
			_ baseline.SkillsRestoreRequest,
		) (baseline.SkillsRestorePayload, error) {
			return baseline.SkillsRestorePayload{
				SchemaVersion: baseline.SkillsRestoreSchemaVersion,
				OK:            true,
				Applied:       true,
				Profile:       "rust-cli",
				Setup:         stringPointerForCLI("rust-cli"),
				Acquisitions:  []baseline.RestoreAcquisition{},
				Skills:        []baseline.RestoreSkill{},
				PlannedChanges: []baseline.RestorePlannedChange{{
					Action: "refresh", Path: ".agents/skills/example/SKILL.md", Skill: "example",
				}},
				PlanDigest: &digest,
				Finding: &baseline.RestoreFinding{
					Code:    "restore.completed",
					Message: "Selected skills match.",
					Action:  "Run baseline plan.",
				},
			}, nil
		}
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runBaselineSkillsRestoreCommandWith(
			context.Background(),
			[]string{"--profile", "rust-cli"},
			&stdout,
			&stderr,
			execute,
		)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("success exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
		for _, want := range []string{
			"Baseline skills restore: applied",
			"restore.completed",
			"Plan Digest: " + digest,
			"refresh .agents/skills/example/SKILL.md [example]",
		} {
			if !strings.Contains(stdout.String(), want) {
				t.Fatalf("text result missing %q:\n%s", want, stdout.String())
			}
		}
	})
}

func stringPointerForCLI(value string) *string {
	return &value
}
