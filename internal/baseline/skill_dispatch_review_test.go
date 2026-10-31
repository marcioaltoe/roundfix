// Suite: skill dispatch review
// Invariant: person-started skills require a request and frontend bundles separate product UI from public pages.
// Boundary IN: embedded module dispatch triggers and activation bundles.
// Boundary OUT: Agent skill loading and managed-guide rendering.
package baseline

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks(t *testing.T) {
	t.Parallel()
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	for _, tc := range []struct {
		module string
		skill  string
	}{
		{"typescript", "architectural-analysis"},
		{"typescript", "humanizer"},
		{"typescript", "writing-clearly-and-concisely"},
		{"context-workflow", "grill-with-docs"},
		{"frontend", "web-design-guidelines"},
	} {
		t.Run(tc.skill, func(t *testing.T) {
			t.Parallel()
			module := readUpstreamSkillDocument(t, "modules/"+tc.module+".json")
			id := "trigger." + tc.module + "." + tc.skill
			matches := 0
			for _, dispatch := range objectsOrEmpty(module["skillDispatch"]) {
				if dispatch["skill"] != tc.skill {
					continue
				}
				for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
					when, _ := trigger["when"].(string)
					if !strings.HasPrefix(when, "When the person asks") {
						t.Errorf("trigger %s = %q, want prefix %q", trigger["id"], when, "When the person asks")
					}
					if trigger["id"] == id {
						matches++
					}
				}
			}
			if matches != 1 {
				t.Errorf("trigger %s occurs %d times, want 1", id, matches)
			}
			for _, bundle := range objectsOrEmpty(activations["bundles"]) {
				if slices.Contains(stringsOrEmpty(bundle["skills"]), tc.skill) {
					t.Errorf("bundle %s includes person-started skill %s", bundle["id"], tc.skill)
				}
			}
		})
	}
}

func TestTheFrontendBundlesSeparateProductUIFromPublicPages(t *testing.T) {
	t.Parallel()
	activations := readUpstreamSkillDocument(t, "skill-activations.json")
	bundles := objectsOrEmpty(activations["bundles"])
	for _, bundle := range bundles {
		skills := stringsOrEmpty(bundle["skills"])
		for _, skill := range []string{"wcag-audit-patterns", "web-design-guidelines"} {
			if slices.Contains(skills, skill) {
				t.Errorf("bundle %s includes on-request audit or review skill %s", bundle["id"], skill)
			}
		}
		if slices.Contains(skills, "frontend-design") && slices.Contains(skills, "interface-design") {
			t.Errorf("bundle %s combines public-page and product-UI design", bundle["id"])
		}
	}
	for _, tc := range []struct {
		bundle  string
		skills  []string
		trigger string
		when    string
	}{
		{
			"bundle.frontend-public-pages", []string{"frontend-design", "fixing-accessibility"},
			"trigger.frontend.public-pages", "Designing or reviewing marketing or other public pages.",
		},
		{
			"bundle.frontend-ui-quality", []string{"interaction-design", "interface-design", "fixing-accessibility"},
			"trigger.frontend.ui-quality", "Designing or reviewing product UI: application screens, forms, and dashboards.",
		},
	} {
		t.Run(tc.bundle, func(t *testing.T) {
			t.Parallel()
			matches := 0
			for _, bundle := range bundles {
				if bundle["id"] == tc.bundle {
					matches++
					if got := stringsOrEmpty(bundle["skills"]); !reflect.DeepEqual(got, tc.skills) {
						t.Errorf("bundle %s skills = %v, want %v", tc.bundle, got, tc.skills)
					}
				}
			}
			if matches != 1 {
				t.Errorf("bundle %s occurs %d times, want 1", tc.bundle, matches)
			}
			matches = 0
			for _, activation := range objectsOrEmpty(activations["activations"]) {
				if activation["id"] == tc.trigger {
					matches++
					want := document{"id": tc.trigger, "owner": "frontend", "when": tc.when, "bundle": tc.bundle}
					if !reflect.DeepEqual(activation, want) {
						t.Errorf("activation %s = %v, want %v", tc.trigger, activation, want)
					}
				}
			}
			if matches != 1 {
				t.Errorf("activation %s occurs %d times, want 1", tc.trigger, matches)
			}
		})
	}
}
