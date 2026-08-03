package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func TestProfilesConfigureExitCodes(t *testing.T) {
	// Sequential: overrides package-level test seams.
	t.Run("declined confirmation", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		original := "watch:\n  max_rounds: 4\n"
		mustWrite(t, configPath, original)
		fragmentPath := filepath.Join(repoDir, "declined-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("declined-preferred", "declined-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
			return false, nil
		})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("declined confirmation exit = %d, want %d stderr=%q", code, exitRunFailed, stderr.String())
		}
		response := decodeProfilesConfigureResponse(t, stdout.String())
		if response.Changed || response.Error != "" || response.Scope != "project" || response.Path != configPath || len(response.Profiles) != 1 || len(response.Changes) != 1 {
			t.Fatalf("declined confirmation changed existing machine fields: %+v", response)
		}
		if !decodeProfilesConfigureRefused(t, stdout.String()) {
			t.Fatalf("declined confirmation did not mark refusal: %s", stdout.String())
		}
		if !strings.Contains(stderr.String(), "confirmation declined") {
			t.Fatalf("declined confirmation diagnostic missing from stderr: %q", stderr.String())
		}
		if strings.Contains(stdout.String(), "confirmation declined") {
			t.Fatalf("declined confirmation diagnostic leaked to stdout: %q", stdout.String())
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("declined confirmation mutated config\nwant: %q\n got: %q", original, got)
		}
	})

	t.Run("non-interactive confirmation EOF", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		original := "watch:\n  max_rounds: 4\n"
		mustWrite(t, configPath, original)
		fragmentPath := filepath.Join(repoDir, "non-interactive-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("non-interactive-preferred", "non-interactive-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		withProfilesConfigureConfirm(t, defaultConfirmProfilesConfigure)
		withProfilesConfigureInput(t, "")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("non-interactive confirmation exit = %d, want %d stderr=%q", code, exitRunFailed, stderr.String())
		}
		if !decodeProfilesConfigureRefused(t, stdout.String()) {
			t.Fatalf("non-interactive confirmation did not mark refusal: %s", stdout.String())
		}
		if !strings.Contains(stderr.String(), "confirmation declined") {
			t.Fatalf("non-interactive refusal diagnostic missing from stderr: %q", stderr.String())
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("non-interactive refusal mutated config\nwant: %q\n got: %q", original, got)
		}
	})

	t.Run("validation failures", func(t *testing.T) {
		t.Run("invalid flags", func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, []string{"profiles", "configure", "--scope", "invalid", "--json"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("invalid flag exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
			}
			if decodeProfilesConfigureRefused(t, stdout.String()) {
				t.Fatalf("invalid flag was marked as a refusal: %s", stdout.String())
			}
		})

		t.Run("configured and removed category", func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			configPath := filepath.Join(repoDir, ".roundfixrc.yml")
			original := "watch:\n  max_rounds: 4\n"
			mustWrite(t, configPath, original)
			fragmentPath := filepath.Join(repoDir, "conflicting-profile.yml")
			mustWrite(t, fragmentPath, profilesConfigureTestFragment("conflicting-preferred", "conflicting-fallback"))
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--remove", "backend", "--yes", "--json"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("conflicting category exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
			}
			response := decodeProfilesConfigureResponse(t, stdout.String())
			if response.Changed || !strings.Contains(response.Error, "cannot be both configured and removed") {
				t.Fatalf("unexpected conflicting category response: %+v", response)
			}
			if decodeProfilesConfigureRefused(t, stdout.String()) {
				t.Fatalf("conflicting category was marked as a refusal: %s", stdout.String())
			}
			if got := mustRead(t, configPath); got != original {
				t.Fatalf("conflicting category mutated config\nwant: %q\n got: %q", original, got)
			}
		})

		t.Run("failed proof", func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			configPath := filepath.Join(repoDir, ".roundfixrc.yml")
			original := "watch:\n  max_rounds: 4\n"
			mustWrite(t, configPath, original)
			fragmentPath := filepath.Join(repoDir, "unproven-profile.yml")
			mustWrite(t, fragmentPath, profilesConfigureTestFragment("unproven-preferred", "unproven-fallback"))
			withAgentRunner(t, &profileReadinessExactRunner{prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
				return agent.SelectionProof{}, &agent.SelectionUnsupportedError{
					Kind:            agent.SelectionModelNotAdvertised,
					Runtime:         req.Runtime.ID,
					Model:           req.Runtime.Model,
					ReasoningEffort: req.Runtime.ReasoningEffort,
				}
			}})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--yes", "--json"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("failed proof exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
			}
			if decodeProfilesConfigureRefused(t, stdout.String()) {
				t.Fatalf("failed proof was marked as a refusal: %s", stdout.String())
			}
			if got := mustRead(t, configPath); got != original {
				t.Fatalf("failed proof mutated config\nwant: %q\n got: %q", original, got)
			}
		})
	})

	for _, tt := range []struct {
		name        string
		original    string
		extraArgs   []string
		wantChanged bool
		wantWritten bool
	}{
		{
			name:        "applied write",
			original:    "watch:\n  max_rounds: 4\n",
			extraArgs:   []string{"--yes"},
			wantChanged: true,
			wantWritten: true,
		},
		{
			name: "already-satisfied no-op",
			original: `profiles:
  backend:
    preferred:
      runtime: codex
      model: stable-preferred
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: stable-fallback
        reasoning_effort: xhigh
`,
			extraArgs: []string{"--yes"},
		},
		{
			name:      "dry run",
			original:  "watch:\n  max_rounds: 4\n",
			extraArgs: []string{"--dry-run"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			configPath := filepath.Join(repoDir, ".roundfixrc.yml")
			mustWrite(t, configPath, tt.original)
			fragmentPath := filepath.Join(repoDir, "stable-profile.yml")
			mustWrite(t, fragmentPath, profilesConfigureTestFragment("stable-preferred", "stable-fallback"))
			withSuccessfulProfilesConfigureProof(t)
			withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
				t.Fatal("explicit consent and dry-run paths must not confirm")
				return false, nil
			})
			args := []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}
			args = append(args, tt.extraArgs...)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("%s exit = %d, want %d stderr=%q", tt.name, code, exitOK, stderr.String())
			}
			response := decodeProfilesConfigureResponse(t, stdout.String())
			if response.Changed != tt.wantChanged {
				t.Fatalf("%s changed = %t, want %t response=%+v", tt.name, response.Changed, tt.wantChanged, response)
			}
			if decodeProfilesConfigureRefused(t, stdout.String()) {
				t.Fatalf("%s was marked as a refusal: %s", tt.name, stdout.String())
			}
			got := mustRead(t, configPath)
			if tt.wantWritten {
				if got == tt.original {
					t.Fatalf("%s did not write config", tt.name)
				}
			} else if got != tt.original {
				t.Fatalf("%s mutated config\nwant: %q\n got: %q", tt.name, tt.original, got)
			}
		})
	}
}

func TestProfilesConfigureChangeSummary(t *testing.T) {
	// Sequential: overrides package-level test seams.
	t.Run("added category agrees in text and machine output", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("added-preferred", "added-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		var preview string
		withProfilesConfigureConfirm(t, func(_ context.Context, _ io.Writer, got string) (bool, error) {
			preview = got
			return true, nil
		})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("configure add exit = %d stderr=%q", code, stderr.String())
		}
		assertProfilesConfigureSummaryLines(t, preview, []string{"added: backend"})
		response := decodeProfilesConfigureResponse(t, stdout.String())
		assertProfilesConfigureChanges(t, response.Changes, []profilesConfigureChange{{
			Category: roundconfig.CategoryBackend,
			Kind:     roundconfig.ChangeAdded,
		}})
	})

	t.Run("replacement excludes untouched categories", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		mustWrite(t, configPath, `
profiles:
  backend:
    preferred: {runtime: codex, model: old-backend, reasoning_effort: high}
    fallbacks:
      - {runtime: claude, model: old-backend-fallback, reasoning_effort: xhigh}
  frontend:
    preferred: {runtime: claude, model: untouched-frontend, reasoning_effort: xhigh}
    fallbacks:
      - {runtime: codex, model: untouched-frontend-fallback, reasoning_effort: high}
`)
		fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("replacement-preferred", "replacement-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		var preview string
		withProfilesConfigureConfirm(t, func(_ context.Context, _ io.Writer, got string) (bool, error) {
			preview = got
			return true, nil
		})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("configure replacement exit = %d stderr=%q", code, stderr.String())
		}
		assertProfilesConfigureSummaryLines(t, preview, []string{"replaced: backend"})
		if strings.Contains(preview, "frontend") {
			t.Fatalf("replacement preview mentions untouched frontend category:\n%s", preview)
		}
		response := decodeProfilesConfigureResponse(t, stdout.String())
		assertProfilesConfigureChanges(t, response.Changes, []profilesConfigureChange{{
			Category: roundconfig.CategoryBackend,
			Kind:     roundconfig.ChangeReplaced,
		}})
	})

	t.Run("declared removals include present and absent categories", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		mustWrite(t, configPath, `
profiles:
  frontend:
    preferred: {runtime: claude, model: configured-frontend, reasoning_effort: xhigh}
    fallbacks:
      - {runtime: codex, model: configured-frontend-fallback, reasoning_effort: high}
`)
		withSuccessfulProfilesConfigureProof(t)
		var preview string
		withProfilesConfigureConfirm(t, func(_ context.Context, _ io.Writer, got string) (bool, error) {
			preview = got
			return true, nil
		})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--remove", "frontend", "--remove", "review", "--json"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("configure removals exit = %d stderr=%q", code, stderr.String())
		}
		assertProfilesConfigureSummaryLines(t, preview, []string{"removed: frontend", "removed: review"})
		response := decodeProfilesConfigureResponse(t, stdout.String())
		assertProfilesConfigureChanges(t, response.Changes, []profilesConfigureChange{
			{Category: roundconfig.CategoryFrontend, Kind: roundconfig.ChangeRemoved},
			{Category: roundconfig.CategoryReview, Kind: roundconfig.ChangeRemoved},
		})
	})

	t.Run("dry run reports the same change and preserves bytes", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		original := "watch:\n  max_rounds: 4\n"
		mustWrite(t, configPath, original)
		fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("dry-summary-preferred", "dry-summary-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
			t.Fatal("dry-run must not ask for confirmation")
			return false, nil
		})
		var textOutput bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--dry-run"}, &textOutput, &stderr)
		if code != exitOK {
			t.Fatalf("text dry-run exit = %d stderr=%q", code, stderr.String())
		}
		assertProfilesConfigureSummaryLines(t, textOutput.String(), []string{"added: backend"})
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("text dry-run mutated config\nwant: %q\n got: %q", original, got)
		}

		var machineOutput bytes.Buffer
		stderr.Reset()
		code = runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--dry-run", "--json"}, &machineOutput, &stderr)
		if code != exitOK {
			t.Fatalf("machine dry-run exit = %d stderr=%q", code, stderr.String())
		}
		response := decodeProfilesConfigureResponse(t, machineOutput.String())
		assertProfilesConfigureChanges(t, response.Changes, []profilesConfigureChange{{
			Category: roundconfig.CategoryBackend,
			Kind:     roundconfig.ChangeAdded,
		}})
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("machine dry-run mutated config\nwant: %q\n got: %q", original, got)
		}
	})
}

func TestProfilesConfigureProofScope(t *testing.T) {
	// Sequential: overrides package-level test seams.
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := `
profiles:
  backend:
    preferred: {runtime: codex, model: old-backend, reasoning_effort: high}
    fallbacks:
      - {runtime: claude, model: old-backend-fallback, reasoning_effort: xhigh}
  frontend:
    preferred: {runtime: stale-adapter, model: unprovable-untouched, reasoning_effort: impossible}
    fallbacks:
      - {runtime: stale-adapter, model: unprovable-untouched-fallback, reasoning_effort: impossible}
`
	mustWrite(t, configPath, original)
	fragmentPath := filepath.Join(repoDir, "written-profiles.yml")
	mustWrite(t, fragmentPath, `
profiles:
  backend:
    preferred: {runtime: codex, model: shared-written, reasoning_effort: high}
    fallbacks:
      - {runtime: claude, model: backend-written-fallback, reasoning_effort: xhigh}
  qa:
    preferred: {runtime: codex, model: shared-written, reasoning_effort: high}
    fallbacks:
      - {runtime: claude, model: qa-written-fallback, reasoning_effort: high}
`)
	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			if strings.Contains(req.Runtime.Model, "unprovable-untouched") {
				return agent.SelectionProof{}, errors.New("untouched tuple must not be proven")
			}
			return agent.SelectionProof{Status: agent.SelectionProofStatusProven}, nil
		},
	}
	withAgentRunner(t, runner)
	withProfilesConfigureConfirm(t, func(_ context.Context, _ io.Writer, preview string) (bool, error) {
		if len(runner.exactRequests) != 3 {
			t.Fatalf("confirmation reached after %d exact proofs, want 3 distinct written tuples", len(runner.exactRequests))
		}
		assertProfilesConfigureSummaryLines(t, preview, []string{"replaced: backend", "added: qa"})
		if strings.Contains(preview, "frontend") {
			t.Fatalf("proof preview mentions untouched frontend category:\n%s", preview)
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("config mutated before confirmation\nwant: %q\n got: %q", original, got)
		}
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("proof-scoped configure exit = %d stderr=%q", code, stderr.String())
	}
	wantModels := []string{"shared-written", "backend-written-fallback", "qa-written-fallback"}
	if len(runner.exactRequests) != len(wantModels) {
		t.Fatalf("exact proof requests = %#v, want %v", runner.exactRequests, wantModels)
	}
	for index, want := range wantModels {
		if got := runner.exactRequests[index].Runtime.Model; got != want {
			t.Fatalf("proof %d model = %q, want %q", index, got, want)
		}
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	assertProfilesConfigureChanges(t, response.Changes, []profilesConfigureChange{
		{Category: roundconfig.CategoryBackend, Kind: roundconfig.ChangeReplaced},
		{Category: roundconfig.CategoryQA, Kind: roundconfig.ChangeAdded},
	})
}

func assertProfilesConfigureSummaryLines(t *testing.T, output string, want []string) {
	t.Helper()
	var got []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "added:") || strings.HasPrefix(line, "replaced:") || strings.HasPrefix(line, "removed:") {
			got = append(got, line)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("profile change summary lines mismatch\nwant: %q\n got: %q\noutput:\n%s", want, got, output)
	}
}

func assertProfilesConfigureChanges(t *testing.T, got []profilesConfigureChange, want []profilesConfigureChange) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("profile configure machine changes mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func decodeProfilesConfigureRefused(t *testing.T, output string) bool {
	t.Helper()
	var response struct {
		Refused *bool `json:"refused"`
	}
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		t.Fatalf("decode profiles configure refusal marker %q: %v", output, err)
	}
	if response.Refused == nil {
		t.Fatalf("profiles configure response is missing boolean refused field: %s", output)
	}
	return *response.Refused
}

func TestPrintProfilesConfigureRefusalWritesJSONBeforeDiagnostic(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer

	code := printProfilesConfigureRefusal(
		profilesConfigureRequest{json: true},
		roundconfig.ProfileConfigResult{Scope: "project", Path: ".roundfixrc.yml"},
		&stdout,
		failingWriter{err: errors.New("stderr unavailable")},
	)

	if code != exitRunFailed {
		t.Fatalf("refusal exit = %d, want %d", code, exitRunFailed)
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if !response.Refused || response.Changed {
		t.Fatalf("refusal response = %+v, want refused and unchanged", response)
	}
}
