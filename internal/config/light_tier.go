package config

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"

	"roundfix/internal/openrouterkey"
)

const DefaultLightModel = "deepseek/deepseek-v4.1-flash"
const DefaultImplementMonthlyCeilingUSD = 10.0

// OpenRouter holds the machine's light tier configuration.
type OpenRouter struct {
	LightModels                []string // nil means the default list; empty means off
	ImplementMonthlyCeilingUSD float64  // zero means the default
}

type openRouterOverlay struct {
	LightModels                *[]string `yaml:"light_models"`
	ImplementMonthlyCeilingUSD *float64  `yaml:"implement_monthly_ceiling_usd"`
}

// OpenRouterImplementKey names the implementation key without exposing its value.
func OpenRouterImplementKey(environ []string) (variable string, present bool) {
	return openrouterkey.Select(environ, openrouterkey.StageImplement)
}

func prepareLightTierConfig(document *yaml.Node, warnings *configWarnings, source ProfileSource) error {
	for _, key := range []string{"light_models", "implement_monthly_ceiling_usd"} {
		if source == ProfileSourceProject && removeYAMLPath(document, []string{"openrouter", key}) {
			warnings.warnIgnoredProjectSetting("openrouter." + key)
		}
	}
	if value, found := yamlValueAtPath(document, []string{"openrouter", "light_models"}); found {
		if value.Kind != yaml.SequenceNode {
			return errors.New("openrouter.light_models must be a sequence of OpenRouter model ids")
		}
		for index, entry := range value.Content {
			field := fmt.Sprintf("openrouter.light_models[%d]", index)
			author, slug, hasSlash := strings.Cut(entry.Value, "/")
			if entry.Tag != "!!str" || !hasSlash || author == "" || slug == "" || strings.Contains(slug, "/") || strings.ContainsFunc(entry.Value, unicode.IsSpace) {
				return fmt.Errorf("%s %q must be an OpenRouter model id <author>/<slug>", field, entry.Value)
			}
			if err := CheckSubscriptionRule(field, "opencode", "openrouter/"+entry.Value); err != nil {
				return err
			}
		}
	}
	if value, found := yamlValueAtPath(document, []string{"openrouter", "implement_monthly_ceiling_usd"}); found {
		var ceiling float64
		if (value.Tag != "!!int" && value.Tag != "!!float") || value.Decode(&ceiling) != nil || ceiling <= 0 || math.IsNaN(ceiling) || math.IsInf(ceiling, 0) {
			return errors.New("openrouter.implement_monthly_ceiling_usd must be a finite number greater than 0")
		}
	}
	return nil
}
