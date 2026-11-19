// Suite: Profile proof setup deadline command contract
// Invariant: every configured tuple receives the effective setup deadline; invalid config is refused.
// Boundary IN: profiles validate, profiles show, profiles configure and profiles check --apply with isolated config and injected proof runner.
// Boundary OUT: ProbeRequest at the Agent boundary and complete command streams and exit code.
package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func TestProfileProofPassesTheConfiguredSetupDeadline(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, project string
		want          time.Duration
	}{
		{name: "built-in", want: roundconfig.DefaultProfileProofSetupTimeout},
		{name: "project", project: "defaults:\n  profile_proof_setup_timeout: 90s\n", want: 90 * time.Second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home, repo := withCLIWorkspace(t)
			if tt.project != "" {
				mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), tt.project)
			}
			runner := &profileReadinessExactRunner{}
			withAgentRunner(t, runner)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"profiles", "validate", "--json"}, &stdout, &stderr)
			if code != exitOK || stderr.Len() != 0 {
				t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
			}
			response := decodeProfilesValidateResponse(t, stdout.String())
			if !response.OK || len(response.Proofs) == 0 || len(runner.exactRequests) != len(response.Proofs) {
				t.Fatalf("proofs = %+v, requests = %+v", response, runner.exactRequests)
			}
			wantSelections := make(map[roundconfig.AgentSelection]bool)
			for _, resolved := range roundconfig.Builtin().Profiles {
				wantSelections[resolved.Profile.Preferred] = true
				for _, fallback := range resolved.Profile.Fallbacks {
					wantSelections[fallback] = true
				}
			}
			for _, req := range runner.exactRequests {
				if req.SetupTimeout != tt.want {
					t.Fatalf("request %+v timeout = %s, want %s", req.Runtime, req.SetupTimeout, tt.want)
				}
				selection := roundconfig.AgentSelection{Runtime: req.Runtime.ID, Model: req.Runtime.Model, ReasoningEffort: req.Runtime.ReasoningEffort}
				if !wantSelections[selection] {
					t.Fatalf("unexpected or duplicate tuple: %+v", selection)
				}
				delete(wantSelections, selection)
			}
			if len(wantSelections) != 0 {
				t.Fatalf("unproven configured tuples: %+v", wantSelections)
			}
			assertNoRunDatabase(t, home)
		})
	}
}

func TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "defaults:\n  profile_proof_setup_timeout: 0s\n")
	withAgentRunner(t, &fakeAgentRunner{probe: func(agent.ProbeRequest) error { t.Fatal("invalid config must fail before proof"); return nil }})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"profiles", "show"}, &stdout, &stderr)
	const wantStderr = "roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m\nRun 'roundfix profiles --help' for usage.\n"
	if code != exitPreflight || stdout.Len() != 0 || stderr.String() != wantStderr {
		t.Fatalf("transcript: stdout=%q stderr=%q exit=%d; want stdout empty stderr=%q exit=%d", stdout.String(), stderr.String(), code, wantStderr, exitPreflight)
	}
	assertNoRunDatabase(t, home)
}

func TestProfileAdoptionPassesTheConfiguredSetupDeadline(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, scope, user, project string
		checkApply                 bool
		want                       time.Duration
	}{
		{name: "configure/built_in", scope: "project", want: roundconfig.DefaultProfileProofSetupTimeout},
		{name: "configure/project_scope_reads_user", scope: "project", user: "75s", want: 75 * time.Second},
		{name: "configure/user_scope_reads_project", scope: "user", project: "90s", want: 90 * time.Second},
		{name: "configure/project_wins", scope: "user", user: "75s", project: "90s", want: 90 * time.Second},
		{name: "check_apply/project", scope: "project", checkApply: true, project: "90s", want: 90 * time.Second},
		{name: "check_apply/user", scope: "user", checkApply: true, user: "75s", want: 75 * time.Second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var home, repo string
			if tt.checkApply {
				home, repo = writeProfilesCheckFixture(t)
			} else {
				home, repo = withCLIWorkspace(t)
			}
			for _, config := range []struct{ path, timeout string }{
				{filepath.Join(home, ".roundfix", "config.yml"), tt.user},
				{filepath.Join(repo, ".roundfixrc.yml"), tt.project},
			} {
				if config.timeout == "" {
					continue
				}
				content, err := os.ReadFile(config.path)
				if err != nil && !errors.Is(err, os.ErrNotExist) {
					t.Fatal(err)
				}
				mustMkdir(t, filepath.Dir(config.path))
				mustWrite(t, config.path, string(content)+"\ndefaults:\n  profile_proof_setup_timeout: "+config.timeout+"\n")
			}
			args := []string{"profiles", "check", "--apply", "--scope", tt.scope, "--yes", "--json"}
			if !tt.checkApply {
				fragment := filepath.Join(repo, "profile.yml")
				mustWrite(t, fragment, profilesConfigureTestFragment("deadline-preferred", "deadline-fallback"))
				args = []string{"profiles", "configure", "--scope", tt.scope, "--file", fragment, "--yes", "--json"}
			}
			runner := withSuccessfulProfilesConfigureProof(t)
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if len(runner.exactRequests) == 0 {
				t.Fatal("profile adoption made no proof requests")
			}
			for _, req := range runner.exactRequests {
				if req.SetupTimeout != tt.want {
					t.Errorf("request %+v timeout = %s, want %s", req.Runtime, req.SetupTimeout, tt.want)
				}
			}
		})
	}
}

func TestProfilesConfigureRefusesAnInvalidSetupDeadlineTranscript(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "defaults:\n  profile_proof_setup_timeout: 0s\n")
	fragment := filepath.Join(repo, "profile.yml")
	mustWrite(t, fragment, profilesConfigureTestFragment("deadline-preferred", "deadline-fallback"))
	withAgentRunner(t, &fakeAgentRunner{probe: func(agent.ProbeRequest) error {
		t.Fatal("invalid deadline must fail before proof")
		return nil
	}})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"profiles", "configure", "--scope", "user", "--file", fragment, "--yes"}, &stdout, &stderr)
	const wantStderr = "roundfix: profiles failed: defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m\nRun 'roundfix profiles --help' for usage.\n"
	if code != exitPreflight || stdout.Len() != 0 || stderr.String() != wantStderr {
		t.Fatalf("transcript: stdout=%q stderr=%q exit=%d; want stdout empty stderr=%q exit=%d", stdout.String(), stderr.String(), code, wantStderr, exitPreflight)
	}
	if _, err := os.Stat(filepath.Join(home, ".roundfix", "config.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("User Config must not be written: %v", err)
	}
}

func TestProfilesConfigureNamesTheFileAndLineOfABrokenConfig(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	projectConfig := filepath.Join(repo, ".roundfixrc.yml")
	mustWrite(t, projectConfig, "defaults:\n  profile_proof_setup_timeout: 90s\nbroken: a: b\n")
	fragment := filepath.Join(repo, "profile.yml")
	mustWrite(t, fragment, profilesConfigureTestFragment("deadline-preferred", "deadline-fallback"))
	withAgentRunner(t, &fakeAgentRunner{probe: func(agent.ProbeRequest) error {
		t.Fatal("broken YAML must fail before proof")
		return nil
	}})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"profiles", "configure", "--scope", "user", "--file", fragment, "--yes"}, &stdout, &stderr)
	wantStderr := "roundfix: profiles failed: parse config " + strconv.Quote(projectConfig) + ": yaml: line 3: mapping values are not allowed in this context\nRun 'roundfix profiles --help' for usage.\n"
	if code != exitPreflight || stdout.Len() != 0 || stderr.String() != wantStderr {
		t.Fatalf("transcript: stdout=%q stderr=%q exit=%d; want stdout empty stderr=%q exit=%d", stdout.String(), stderr.String(), code, wantStderr, exitPreflight)
	}
	if _, err := os.Stat(filepath.Join(home, ".roundfix", "config.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("User Config must not be written: %v", err)
	}
}
