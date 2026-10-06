package config

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

func TestLightTierUserConfig(t *testing.T) {
	for _, tc := range []struct {
		name, user string
		models     []string
		ceiling    float64
	}{
		{"defaults", "", []string{DefaultLightModel}, DefaultImplementMonthlyCeilingUSD},
		{"custom", "openrouter:\n  light_models: [deepseek/deepseek-v4.1-flash, qwen/qwen3-coder]\n  implement_monthly_ceiling_usd: 25.5\n", []string{DefaultLightModel, "qwen/qwen3-coder"}, 25.5},
		{"off", "openrouter:\n  light_models: []\n", []string{}, DefaultImplementMonthlyCeilingUSD},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts, warnings := jevCeilingFixture(t, tc.user, "")
			loaded, err := Load(opts)
			if err != nil {
				t.Fatal(err)
			}
			got := loaded.Config.OpenRouter
			if !slices.Equal(got.LightModels, tc.models) || got.ImplementMonthlyCeilingUSD != tc.ceiling || warnings.Len() != 0 {
				t.Fatalf("OpenRouter=%+v warnings=%q", got, warnings.String())
			}
			if tc.name == "off" && got.LightModels == nil {
				t.Fatal("explicit empty list must remain non-nil")
			}
		})
	}
	for _, id := range []string{"openai/gpt-6.1-sol", "anthropic/claude-opus-5.5", "openrouter/auto", "deepseek", "deepseek/", "/flash", "deepseek/a/b", "deepseek/a b"} {
		t.Run("refuse model "+id, func(t *testing.T) {
			opts, _ := jevCeilingFixture(t, fmt.Sprintf("openrouter:\n  light_models: [%q]\n", id), "")
			_, err := Load(opts)
			field := "openrouter.light_models[0]"
			want := fmt.Sprintf("%s %q must be an OpenRouter model id <author>/<slug>", field, id)
			if ruleErr := CheckSubscriptionRule(field, "opencode", "openrouter/"+id); ruleErr != nil {
				want = ruleErr.Error()
			}
			want = fmt.Sprintf("parse config %q: %s", filepath.Join(opts.HomeDir, userConfigRelPath), want)
			if err == nil || err.Error() != want {
				t.Fatalf("err=%v want=%s", err, want)
			}
		})
	}
	for _, value := range []string{"0", "-1", ".inf", "-.inf", ".nan", "a string", "'10'", "null", "true", "[]"} {
		t.Run("refuse ceiling "+value, func(t *testing.T) {
			opts, _ := jevCeilingFixture(t, "openrouter:\n  implement_monthly_ceiling_usd: "+value+"\n", "")
			_, err := Load(opts)
			want := fmt.Sprintf("parse config %q: openrouter.implement_monthly_ceiling_usd must be a finite number greater than 0", filepath.Join(opts.HomeDir, userConfigRelPath))
			if err == nil || err.Error() != want {
				t.Fatalf("err=%v want=%s", err, want)
			}
		})
	}
	for _, value := range []string{"null", "deepseek/flash", "{}", "[42]", "[null]", "[{}]"} {
		t.Run("refuse list "+value, func(t *testing.T) {
			opts, _ := jevCeilingFixture(t, "openrouter:\n  light_models: "+value+"\n", "")
			if _, err := Load(opts); err == nil {
				t.Fatalf("Load accepted light_models: %s", value)
			}
		})
	}
	for _, user := range []string{"", "openrouter:\n  light_models: []\n  implement_monthly_ceiling_usd: 30\n"} {
		t.Run("project ignores "+user, func(t *testing.T) {
			opts, warnings := jevCeilingFixture(t, user, "openrouter:\n  light_models: invalid\n  implement_monthly_ceiling_usd: invalid\n")
			loaded, err := Load(opts)
			if err != nil {
				t.Fatal(err)
			}
			wantModels, wantCeiling := []string{DefaultLightModel}, DefaultImplementMonthlyCeilingUSD
			if user != "" {
				wantModels, wantCeiling = []string{}, 30
			}
			if !slices.Equal(loaded.Config.OpenRouter.LightModels, wantModels) || loaded.Config.OpenRouter.ImplementMonthlyCeilingUSD != wantCeiling {
				t.Fatalf("OpenRouter=%+v", loaded.Config.OpenRouter)
			}
			want := ""
			for _, key := range []string{"openrouter.light_models", "openrouter.implement_monthly_ceiling_usd"} {
				want += fmt.Sprintf("config: %s in Project Config is ignored; set %s in User Config\n", key, key)
			}
			if warnings.String() != want {
				t.Fatalf("warnings=%q want=%q", warnings.String(), want)
			}
		})
	}
}

func TestOpenRouterImplementKeyNamesOneVariable(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		present bool
	}{
		{"set", map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": "fake"}, true},
		{"empty", map[string]string{"ROUNDFIX_OPENROUTER_API_KEY": ""}, false},
		{"absent", nil, false},
		{"generic only", map[string]string{"OPENROUTER_API_KEY": "fake"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var lookedUp []string
			variable, present := OpenRouterImplementKey(func(key string) string { lookedUp = append(lookedUp, key); return tc.env[key] })
			if variable != "ROUNDFIX_OPENROUTER_API_KEY" || present != tc.present || !slices.Equal(lookedUp, []string{variable}) {
				t.Fatalf("variable=%q present=%t lookedUp=%v", variable, present, lookedUp)
			}
		})
	}
}
