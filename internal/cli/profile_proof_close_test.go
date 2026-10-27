// Suite: profile proof cleanup warnings.
// Invariant: a passed proof reports its unclosed session in text and JSON.
package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"roundfix/internal/agent"
)

type unclosedProfileProofRunner struct{ fakeAgentRunner }

func (runner *unclosedProfileProofRunner) ProveProfileSelection(_ context.Context, req agent.ProbeRequest) (agent.SelectionProof, error) {
	return agent.SelectionProof{Status: agent.SelectionProofStatusProven, EffectiveAccessPolicy: req.Runtime.RequestedPolicy(), Unclosed: &agent.UnclosedProofSession{Session: "recorded-proof", Warning: "disposable session recorded for the next proof"}}, nil
}
func TestProfilesValidateReportsAnUnclosedProofSession(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"text", "json"} {
		t.Run(mode, func(t *testing.T) {
			withCLIWorkspace(t)
			withAgentRunner(t, &unclosedProfileProofRunner{})
			args := []string{"profiles", "validate", "--category", "backend"}
			if mode == "json" {
				args = append(args, "--json")
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, args, &stdout, &stderr)
			if code != exitOK || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			const warning = "disposable session recorded for the next proof"
			if mode == "json" {
				response := decodeProfilesValidateResponse(t, stdout.String())
				if !response.OK || len(response.Proofs) == 0 {
					t.Fatalf("response=%+v", response)
				}
				for _, proof := range response.Proofs {
					if proof.Status != "passed" || proof.Warning != warning {
						t.Fatalf("proof=%+v", proof)
					}
				}
				if !strings.Contains(stdout.String(), `"warning":"`+warning+`"`) {
					t.Fatalf("missing JSON field: %s", stdout.String())
				}
			} else if !strings.Contains(stdout.String(), "   - backend preferred source=built-in\n   warning: "+warning+"\n") {
				t.Fatalf("warning must follow references: %s", stdout.String())
			}
		})
	}
}
