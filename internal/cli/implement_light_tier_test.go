// Boundary: implement command with temporary User Config, Home, repository and fake Agent runner.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/lighttier"
	"roundfix/internal/openrouterkey"
	"roundfix/internal/spec"
)

func TestImplementBuildsTheLightTierPlan(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing-key", "configured-light", "override", "disabled-models", "configured-ceiling"} {
		t.Run(name, func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", taskType: "docs"}, {id: "task_02", taskType: "backend", needs: []string{"task_01"}}})
			for _, id := range []string{"task_01", "task_02"} {
				path := filepath.Join(repo, "docs", "specs", implementTestSlug, id+".md")
				mustWrite(t, path, strings.Replace(mustRead(t, path), "status: pending\n", "status: pending\ncomplexity: low\n", 1))
			}
			gitImplement(t, repo, "add", ".")
			gitImplement(t, repo, "commit", "-m", "seed low tasks")
			model := "x-ai/grok-4.5"
			config := "openrouter:\n  light_models: [" + model + "]\n  implement_monthly_ceiling_usd: 0.75\n"
			if name == "disabled-models" {
				config = "openrouter:\n  light_models: []\n"
			}
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), config)
			var requests []agent.ExecuteRequest
			runner := &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}, onTask: func(req agent.ExecuteRequest, _ string) error { requests = append(requests, req); return nil }}
			withImplementCollaborators(t, runner)
			environment := commandEnvironmentForTest(t)
			variable := openrouterkey.Shared
			environment.environ = withEnvValue(environment.environ, variable, "")
			if name != "missing-key" {
				environment.environ = withEnvValue(environment.environ, variable, "fixture-secret")
			}
			args := []string{"implement", "--spec", implementTestSlug, "--no-input"}
			if name == "override" {
				args = append(args, "--agent", "codex", "--model", "gpt-5.6-sol", "--reasoning-effort", "high")
			}
			if name == "configured-ceiling" {
				if err := lighttier.AppendSpend(home, environment.dependencies.implementBudgetNow(), lighttier.SpendLine{CostUSD: 0.8, CostSource: "opencode"}); err != nil {
					t.Fatal(err)
				}
			}
			var stdout, stderr bytes.Buffer
			code := runWithContext(context.Background(), args, &stdout, &stderr, environment)
			if code != exitOK {
				t.Fatalf("exit %d: stderr=%s stdout=%s", code, stderr.String(), stdout.String())
			}
			if len(requests) != 2 {
				t.Fatalf("work requests = %d", len(requests))
			}
			if strings.Contains(stdout.String()+stderr.String(), "fixture-secret") {
				t.Fatal("key value reached command output")
			}
			warning := "roundfix: warning: light tier skipped for Task task_01: " + openrouterkey.Implement + " is not set; it runs on its docs profile"
			if name == "missing-key" {
				if !strings.Contains(stderr.String(), warning) || strings.Contains(stdout.String(), "light tier skipped") {
					t.Fatalf("Surface Transcript 3 mismatch: %s", stderr.String())
				}
			} else if name == "configured-light" {
				path := filepath.Join(home, ".roundfix", "openrouter", "implement", environment.dependencies.implementBudgetNow().UTC().Format("2006-01")+".jsonl")
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")
				if len(lines) != 2 || strings.Contains(string(data), "fixture-secret") {
					t.Fatal("unexpected spend lines or key value in log")
				}
				for _, line := range lines {
					var row lighttier.SpendLine
					if err := json.Unmarshal([]byte(line), &row); err != nil {
						t.Fatal(err)
					}
					if row.Repository != repo || row.Model != "openrouter/"+model || row.CostSource != "unreported" {
						t.Fatalf("spend identity = %+v", row)
					}
				}
				for _, req := range requests {
					if req.Runtime.ID != "opencode" || req.Runtime.Model != "openrouter/"+model || req.Runtime.ReasoningEffort != "" || req.Runtime.OpenRouterKeyVariable != variable {
						t.Fatalf("loaded light selection = %+v", req.Runtime)
					}
				}
			} else {
				if name == "configured-ceiling" && !strings.Contains(stderr.String(), "this month's light spend US$0.8000 reached the ceiling of US$0.75") {
					t.Fatalf("loaded ceiling missing: %s", stderr.String())
				}
				for _, req := range requests {
					if req.Runtime.ID != "codex" || req.Runtime.OpenRouterKeyVariable != "" {
						t.Fatalf("standard selection = %+v", req.Runtime)
					}
				}
				if name != "configured-ceiling" && strings.Contains(stderr.String(), "light tier skipped") {
					t.Fatal("disabled tier printed a skip warning")
				}
				if name == "override" {
					for _, req := range requests {
						if req.Runtime.Model != "gpt-5.6-sol" || req.Runtime.ReasoningEffort != "high" {
							t.Fatalf("override = %+v", req.Runtime)
						}
					}
				}
			}
		})
	}
}
