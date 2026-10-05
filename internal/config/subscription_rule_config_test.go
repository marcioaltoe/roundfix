package config

import (
	"fmt"
	"strings"
	"testing"
)

func subscriptionProfileYAML(runtime, model string, fallback bool) string {
	selection := fmt.Sprintf("{runtime: %s, model: %q, reasoning_effort: high}", runtime, model)
	other := `{runtime: claude, model: opus, reasoning_effort: high}`
	if fallback {
		selection, other = other, selection
	}
	return fmt.Sprintf("profiles:\n  docs:\n    preferred: %s\n    fallbacks: [%s]\n", selection, other)
}

func TestSubscriptionRuleRefusesEveryConfigScope(t *testing.T) {
	t.Parallel()
	for _, model := range []string{
		"roundfix-openrouter/typesafe/jev-router",
		"openrouter/openai/gpt-6.1-sol",
		"openrouter/anthropic/claude-opus-5.5",
		"openrouter/openrouter/auto",
	} {
		t.Run(model, func(t *testing.T) {
			message := fmt.Sprintf("%q can reach OpenAI or Anthropic models through OpenRouter; %s", model, SubscriptionRule)
			if strings.HasPrefix(model, "roundfix-openrouter/") {
				message = fmt.Sprintf("%q names the retired Jev Router; %s", model, SubscriptionRule)
			}
			for _, scope := range []string{"user", "project", "user-shadowed"} {
				for _, placement := range []string{"preferred", "fallback", "legacy"} {
					t.Run(scope+"/"+placement, func(t *testing.T) {
						content := subscriptionProfileYAML("opencode", model, placement == "fallback")
						field := "profiles.docs.preferred.model"
						if placement == "fallback" {
							field = "profiles.docs.fallbacks[0].model"
						}
						if placement == "legacy" {
							content = fmt.Sprintf("defaults:\n  agent: opencode\nruntimes:\n  opencode:\n    model: %q\n    reasoning_effort: high\n", model)
							field = "profiles.general.preferred.model"
						}
						user, project := "", ""
						if scope != "project" {
							user = content
						} else {
							project = content
						}
						if scope == "user-shadowed" {
							project = "defaults:\n  agent: codex\nruntimes:\n  codex:\n    model: gpt-6.1-sol\n    reasoning_effort: high\n"
						}
						opts, _ := jevCeilingFixture(t, user, project)
						_, err := Load(opts)
						if err == nil || !strings.HasSuffix(err.Error(), field+" "+message) {
							t.Fatalf("error=%v, want suffix %q", err, field+" "+message)
						}
					})
				}
			}
			t.Run("override", func(t *testing.T) {
				opts, _ := jevCeilingFixture(t, "", "")
				loaded, err := Load(opts)
				if err != nil {
					t.Fatal(err)
				}
				override := AgentSelection{Runtime: " opencode ", Model: " " + model + " ", ReasoningEffort: "high"}
				_, err = ResolveProfile(loaded.Config, CategoryDocs, &override)
				if err == nil || err.Error() != "invocation preferred.model "+message {
					t.Fatalf("error=%v, want %q", err, "invocation preferred.model "+message)
				}
			})
		})
	}
	for _, selection := range []AgentSelection{
		{Runtime: "opencode", Model: "openrouter/deepseek/deepseek-v4-pro"},
		{Runtime: "codex", Model: "gpt-6.1-sol"},
	} {
		for _, scope := range []string{"user", "project"} {
			t.Run(scope+"/allowed/"+selection.Model, func(t *testing.T) {
				content := subscriptionProfileYAML(selection.Runtime, selection.Model, false)
				user, project := "", ""
				if scope == "user" {
					user = content
				} else {
					project = content
				}
				opts, _ := jevCeilingFixture(t, user, project)
				loaded, err := Load(opts)
				if err != nil {
					t.Fatal(err)
				}
				resolved, err := ResolveProfile(loaded.Config, CategoryDocs, &selection)
				if err != nil || resolved.Profile.Preferred.Model != selection.Model {
					t.Fatalf("profile=%+v error=%v", resolved, err)
				}
			})
		}
	}
}

func TestRetiredRouterCreditFloorIsDeprecated(t *testing.T) {
	t.Parallel()
	const warning = "config: jev.router_min_credit_usd is deprecated and ignored; the Jev Router was retired, so remove it\n"
	for _, scope := range []string{"user", "project", "both"} {
		for _, value := range []string{"20", "0", "-1", ".inf", ".nan", "null", "true", "nope", "[]", "{}", "'50'"} {
			t.Run(scope+"/"+value, func(t *testing.T) {
				content := "jev:\n  router_min_credit_usd: " + value + "\n"
				user, project := "", ""
				if scope != "project" {
					user = content
				}
				if scope != "user" {
					project = content
				}
				opts, stderr := jevCeilingFixture(t, user, project)
				if _, err := Load(opts); err != nil {
					t.Fatal(err)
				}
				if stderr.String() != warning {
					t.Fatalf("stderr=%q, want %q", stderr.String(), warning)
				}
			})
		}
	}
}
