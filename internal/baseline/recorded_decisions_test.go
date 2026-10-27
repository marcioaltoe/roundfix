// Suite: recorded Baseline decisions
// Invariant: only accepted reachable recorded answers carry; flags and fixed values win.
// Boundary IN: Setup Manifest bytes, selected profile, explicit answers, and catalog graph.
// Boundary OUT: rendering and apply.
package baseline

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects(t *testing.T) {
	t.Parallel()
	catalog := loadManifestInputCatalog(t)
	for _, test := range []struct {
		name         string
		override     bool
		fixed        bool
		invalid      bool
		incompatible bool
	}{
		{name: "carry closure"}, {name: "flag changes closure", override: true}, {name: "fixed profile answer", fixed: true}, {name: "catalog refuses recorded value", invalid: true}, {name: "incompatible recorded profile", incompatible: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			stored := manifestInputDecisionsForProfile(t, root, catalog, "go-cli-tui")
			stored["unselected.decision"] = "old"
			stored["autonomous.enabled"] = true
			if test.invalid {
				stored["domain.layout"] = "invalid"
			}
			manifest := newManifestInputFixture(t, root, catalog, "go-cli-tui", stored)
			if test.incompatible {
				manifest.Profile = "unknown-old-profile"
				manifest.Generator.Baseline = "baseline." + manifest.Profile + "-" + ManifestVersion
			}
			writeManifestInputFixture(t, root, manifest)
			profile, err := ResolveProfile(root, "python", catalog)
			if err != nil {
				t.Fatal(err)
			}
			explicit := []DecisionValue{}
			if test.override {
				explicit = append(explicit, DecisionValue{ID: "autonomous.enabled", Value: false})
			}
			if test.fixed {
				profile.Values = map[string]any{"domain.layout": "single-context"}
			}
			decisions, carried, warnings, err := carryRecordedDecisions(root, profile, explicit, catalog)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Contains(carried, "verification.gate") || !slices.IsSorted(carried) {
				t.Fatalf("carried=%v", carried)
			}
			if slices.Contains(carried, "unselected.decision") {
				t.Fatal("carried unselected decision")
			}
			if test.override {
				if slices.Contains(carried, "autonomous.enabled") || slices.Contains(carried, "runtime.backend") || slices.Contains(carried, "runtime.design") {
					t.Fatalf("wrong closure: %v", carried)
				}
				if decisions[0].Value != false {
					t.Fatal("flag lost")
				}
			} else if !slices.Contains(carried, "runtime.backend") || !slices.Contains(carried, "runtime.design") {
				t.Fatalf("lost dependent answers: %v", carried)
			}
			wants := []string{"recorded decision unselected.decision is not carried: profile python does not select it"}
			if test.fixed {
				wants = append(wants, "recorded decision domain.layout is not carried: profile python fixes it")
			}
			if test.invalid {
				wants = append(wants, "recorded decision domain.layout is not carried: the catalog refuses its recorded value")
			}
			for _, want := range wants {
				count := 0
				for _, warning := range warnings {
					if warning.Message == want && warning.Code == "baseline.decision.not-carried" && warning.Path == manifestPath {
						count++
					}
				}
				if count != 1 {
					t.Errorf("warning %q count=%d: %+v", want, count, warnings)
				}
			}
			if test.fixed || test.invalid {
				if slices.Contains(carried, "domain.layout") {
					t.Fatal("carried refused/fixed value")
				}
			}
		})
	}
}

func TestCarryRecordedDecisionsWithoutAManifestChangesNothing(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	catalog := loadManifestInputCatalog(t)
	profile, err := ResolveProfile(root, "python", catalog)
	if err != nil {
		t.Fatal(err)
	}
	explicit := []DecisionValue{{ID: "autonomous.enabled", Value: false}}
	decisions, carried, warnings, err := carryRecordedDecisions(root, profile, explicit, catalog)
	if err != nil || !reflect.DeepEqual(decisions, explicit) || len(carried) != 0 || len(warnings) != 0 {
		t.Fatalf("decisions=%v carried=%v warnings=%v err=%v", decisions, carried, warnings, err)
	}
	writeManifestInputBytes(t, root, []byte("{invalid"))
	decisions, carried, warnings, err = carryRecordedDecisions(root, profile, explicit, catalog)
	if err != nil || !reflect.DeepEqual(decisions, explicit) || len(carried) != 0 || len(warnings) != 1 || warnings[0].Code != "baseline.decision.manifest-unreadable" || !strings.Contains(warnings[0].Message, "cannot be read") {
		t.Fatalf("unreadable: %v %v %+v %v", decisions, carried, warnings, err)
	}
}
