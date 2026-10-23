package lighttier

import (
	"errors"
	"reflect"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess(t *testing.T) {
	t.Parallel()
	models := []string{"first", "second"}
	for _, name := range []string{"full-access", "runtime-default", "factory-error", "second-model-serves", "empty", "selections"} {
		t.Run(name, func(t *testing.T) {
			input := models
			if name == "empty" {
				input = nil
			}
			factory := func(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
				if name == "factory-error" {
					return agent.RuntimeSpec{}, errors.New("factory refused")
				}
				runtime, err := agent.RuntimeFor(agent.RuntimeOptions{Agent: selection.Runtime, Model: selection.Model, EnableFullAccess: name != "runtime-default"})
				if name == "second-model-serves" && selection.Model == "openrouter/second" {
					runtime.SupportedFullAccessMode = "full-access"
				}
				return runtime, err
			}
			if name == "selections" {
				want := []roundconfig.AgentSelection{{Runtime: "opencode", Model: "openrouter/first"}, {Runtime: "opencode", Model: "openrouter/second"}}
				if got := Selections(input); !reflect.DeepEqual(got, want) {
					t.Fatalf("selections = %+v", got)
				}
				return
			}
			kept, refused := ServingSelections(input, factory)
			switch name {
			case "full-access":
				if len(kept) != 0 || len(refused) != 2 {
					t.Fatalf("kept=%v refused=%v", kept, refused)
				}
				for i, refusal := range refused {
					if refusal.Selection != Selections(models)[i] || refusal.Err.Predicate != agent.AccessModeAvailablePredicate || refusal.Err.Runtime != "opencode" || refusal.Err.Policy != agent.AccessPolicyFullAccess {
						t.Fatalf("refusal = %+v", refusal)
					}
					const want = `requested access policy "full-access" failed predicate "runtime_has_access_mode" for runtime "opencode"; recovery: disable full access or select a runtime that supports it`
					if refusal.Err.Error() != want {
						t.Fatalf("error = %q", refusal.Err.Error())
					}
				}
			case "runtime-default", "factory-error":
				if !reflect.DeepEqual(kept, Selections(models)) || len(refused) != 0 {
					t.Fatalf("kept=%v refused=%v", kept, refused)
				}
			case "second-model-serves":
				if len(kept) != 1 || kept[0].Model != "openrouter/second" || len(refused) != 1 || refused[0].Selection.Model != "openrouter/first" {
					t.Fatalf("kept=%v refused=%v", kept, refused)
				}
			case "empty":
				if len(kept) != 0 || len(refused) != 0 {
					t.Fatalf("kept=%v refused=%v", kept, refused)
				}
			}
		})
	}
}
