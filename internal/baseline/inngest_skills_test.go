// Suite: official Inngest skills in the Baseline catalog.
// Invariant: backend requires and dispatches all seven replacements and retires inngest.
// Boundary IN: embedded module, restore profiles and Retired Skill set.
// Boundary OUT: installed skill trees, network acquisition and CLI migration.

package baseline

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestTheBackendModuleTakesTheSevenOfficialInngestSkills(t *testing.T) {
	t.Parallel()
	want := map[string]string{
		"inngest-cli":               "Installing or running the Inngest CLI and its local Dev Server.",
		"inngest-durable-functions": "Building or changing Inngest functions, their triggers, retries, and schedules.",
		"inngest-events":            "Designing Inngest events, fan-out, and idempotent event handling.",
		"inngest-flow-control":      "Configuring Inngest concurrency, throttling, rate limits, debounce, or batching.",
		"inngest-middleware":        "Adding Inngest middleware for logging, tracing, error tracking, or dependency injection.",
		"inngest-setup":             "Setting up Inngest in a project: the SDK, client configuration, environment variables, and the serve endpoint.",
		"inngest-steps":             "Building or changing the steps of Inngest functions: step runs, sleeps, waits for events, and invocations.",
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := catalog.Module("backend")
	if !ok {
		t.Fatal("backend module is missing")
	}
	var module struct {
		RequiredSkills []string `json:"requiredSkills"`
		SkillDispatch  []struct {
			Skill    string `json:"skill"`
			Triggers []struct {
				ID   string `json:"id"`
				When string `json:"when"`
			} `json:"triggers"`
		} `json:"skillDispatch"`
	}
	if err := json.Unmarshal(entry.Data, &module); err != nil {
		t.Fatal(err)
	}
	if slices.Contains(module.RequiredSkills, "inngest") {
		t.Error("backend still requires inngest")
	}
	counts := make(map[string]int)
	for _, dispatch := range module.SkillDispatch {
		if dispatch.Skill == "inngest" {
			t.Error("backend still dispatches inngest")
		}
		when, expected := want[dispatch.Skill]
		if !expected {
			continue
		}
		counts[dispatch.Skill]++
		if len(dispatch.Triggers) != 1 {
			t.Errorf("%s triggers = %v, want exactly one", dispatch.Skill, dispatch.Triggers)
			continue
		}
		trigger := dispatch.Triggers[0]
		if trigger.ID != "trigger.backend."+dispatch.Skill || trigger.When != when {
			t.Errorf("%s trigger = %+v, want trigger.backend.%s: %s", dispatch.Skill, trigger, dispatch.Skill, when)
		}
	}
	for name := range want {
		if !slices.Contains(module.RequiredSkills, name) || counts[name] != 1 {
			t.Errorf("%s required = %t, dispatch count = %d; want true, 1", name, slices.Contains(module.RequiredSkills, name), counts[name])
		}
	}
	for _, id := range []string{"standard-typescript-monorepo", "go-cli-typescript-monorepo"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			profile, err := loadRestoreProfile(catalog, id)
			if err != nil {
				t.Fatal(err)
			}
			if _, exists := profile.Skills["inngest"]; exists {
				t.Error("profile still pins inngest")
			}
			if _, exists := profile.RequiredSkills["inngest"]; exists {
				t.Error("profile still requires inngest")
			}
			for name := range want {
				if _, exists := profile.RequiredSkills[name]; !exists {
					t.Errorf("profile does not require %s", name)
				}
				skill, exists := profile.Skills[name]
				if !exists || skill.Source.Path != "skills/05-implementation-loop/"+name || skill.TreeDigest == "" {
					t.Errorf("profile skill %s = %+v, exists = %t; want catalog path and tree digest", name, skill, exists)
				}
			}
		})
	}
	if !slices.Contains(RetiredSkills(), "inngest") {
		t.Error("inngest is not a Retired Skill")
	}
}
