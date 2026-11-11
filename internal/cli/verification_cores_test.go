package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/corebudget"
	"roundfix/internal/daemon"
	"roundfix/internal/spec"
)

func withVerificationCoreBudgetInputs(t *testing.T, gomaxprocs string) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.logicalCPUs = func() int { return 12 }
		environ := deps.environ
		deps.environ = func() []string {
			var env []string
			for _, entry := range environ() {
				name, _, _ := strings.Cut(entry, "=")
				if name != corebudget.GoVariable && name != corebudget.Variable {
					env = append(env, entry)
				}
			}
			if gomaxprocs != "" {
				env = append(env, corebudget.GoVariable+"="+gomaxprocs)
			}
			return env
		}
	})
}

type implementCoreBudgetVerifier struct {
	*fakeVerifier
	requests []daemon.VerifyRequest
}

func (v *implementCoreBudgetVerifier) Verify(ctx context.Context, req daemon.VerifyRequest) (daemon.VerifyResult, error) {
	v.requests = append(v.requests, req)
	return v.fakeVerifier.Verify(ctx, req)
}

func TestImplementGivesTheEngineTheVerificationCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, gomaxprocs string
		want             int
	}{{"half divided by two", "", 3}, {"explicit GOMAXPROCS", "7", 7}} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			writeUserConfig(t, home, "runs:\n  max_verifying: 2\ndefaults:\n  verification: repository check\n")
			withVerificationCoreBudgetInputs(t, row.gomaxprocs)
			var requests []agent.ExecuteRequest
			runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, onTask: func(req agent.ExecuteRequest, _ string) error {
				requests = append(requests, req)
				return nil
			}}
			_, base, _, _ := withImplementCollaborators(t, runner)
			verifier := &implementCoreBudgetVerifier{fakeVerifier: base}
			withVerifier(t, verifier)
			var stdout, stderr bytes.Buffer
			if code := runCLIContext(t, t.Context(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr); code != exitOK {
				t.Fatalf("implement exit = %d; stderr=%s stdout=%s", code, stderr.String(), stdout.String())
			}
			if len(verifier.requests) == 0 || len(requests) == 0 {
				t.Fatalf("missing requests: Verification=%d Agent=%d", len(verifier.requests), len(requests))
			}
			for _, req := range verifier.requests {
				if req.Cores != row.want {
					t.Errorf("Verification %q cores = %d, want %d", req.Command, req.Cores, row.want)
				}
			}
			for _, req := range requests {
				if req.Cores != row.want {
					t.Errorf("Agent Session %s cores = %d, want %d", req.Session.Name, req.Cores, row.want)
				}
			}
		})
	}
}
