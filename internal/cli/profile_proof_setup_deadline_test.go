// Suite: Profile proof setup deadline command contract
// Invariant: every configured tuple receives the effective setup deadline; invalid config is refused.
// Boundary IN: profiles validate and profiles show with isolated config and injected proof runner.
// Boundary OUT: ProbeRequest at the Agent boundary and complete command streams and exit code.
package cli

import (
	"bytes"
	"path/filepath"
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
