// Suite: Context7 capability compatibility
// Invariant: current and prior skill names satisfy Context7; restoration uses the current setup member.
// Boundary IN: temporary repositories, installed-skill evidence, alignment, and embedded setup membership.
// Boundary OUT: installation, executables, and network.
package baseline

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func context7TestCapability(t *testing.T) RepositoryCapability {
	t.Helper()
	for _, capability := range universalCapabilities {
		if capability.ID == "capability.context7" {
			return capability
		}
	}
	t.Fatal("Context7 capability missing")
	return RepositoryCapability{}
}

func TestTheContext7CapabilityIsSatisfiedByTheCurrentSkill(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		both bool
	}{
		{name: "current only"},
		{name: "current takes precedence", both: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			repository := t.TempDir()
			writeProfileAlignmentFile(t, repository, ".agents/skills/context7-cli/SKILL.md", "# Current Context7\n")
			if test.both {
				writeProfileAlignmentFile(t, repository, ".agents/skills/context7/SKILL.md", "# Prior Context7\n")
			}
			root, err := os.OpenRoot(repository)
			if err != nil {
				t.Fatal(err)
			}
			defer root.Close()
			capability := context7TestCapability(t)
			evidence := collectInstalledSkillEvidence(root, capability)
			if evidence.SourcePath != ".agents/skills/context7-cli/SKILL.md" {
				t.Fatalf("current evidence = %+v", evidence)
			}
			if outcome := evaluateCapability(capability, []CapabilityEvidence{evidence}); outcome.Status != CapabilitySatisfied {
				t.Fatalf("current outcome = %+v", outcome)
			}
		})
	}
}

func TestTheContext7CapabilityIsSatisfiedByItsPriorSkill(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	writeProfileAlignmentFile(t, repository, ".agents/skills/context7/SKILL.md", "# Prior Context7\n")
	root, err := os.OpenRoot(repository)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	capability := context7TestCapability(t)
	evidence := collectInstalledSkillEvidence(root, capability)
	if evidence.Status != CapabilityEvidencePresent || evidence.Strength != CapabilityEvidenceVerified || evidence.SourcePath != ".agents/skills/context7/SKILL.md" {
		t.Fatalf("prior evidence = %+v", evidence)
	}
	if outcome := evaluateCapability(capability, []CapabilityEvidence{evidence}); outcome.Status != CapabilitySatisfied {
		t.Fatalf("prior outcome = %+v", outcome)
	}
	capability.Probe = map[string]any{"skill": "context7-cli", "priorSkills": []string{"../context7"}}
	if got := collectInstalledSkillEvidence(root, capability); got.Status != CapabilityEvidenceInvalid {
		t.Fatalf("unsafe prior evidence = %+v", got)
	}
}

func TestAMissingContext7SkillIsRemediatedWithTheCurrentName(t *testing.T) {
	t.Parallel()
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	if err := os.Remove(filepath.Join(repository, ".agents/skills/context7/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{ProfileID: "standard-typescript-monorepo", Decisions: standardTypeScriptDecisions("make verify")}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	outcome, ok := findCapabilityOutcome(alignment.Capabilities, "capability.context7")
	if !ok || outcome.Status == CapabilitySatisfied {
		t.Fatalf("missing capability = %+v", outcome)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.context7")
	if !ok || !divergence.Blocking || !strings.Contains(divergence.NextAction, "--skill context7-cli") {
		t.Fatalf("missing remediation = %+v", divergence)
	}
}

func TestEveryRestorableCapabilitySkillIsRequiredAndInEverySetup(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for id, name := range universalCapabilityRestoreSkills {
		found := false
		for _, capability := range universalCapabilities {
			if capability.ID == id {
				found = true
				if capability.Probe["skill"] != name {
					t.Errorf("%s restore %s differs from probe %v", id, name, capability.Probe)
				}
			}
		}
		if !found {
			t.Errorf("restore capability %s missing", id)
		}
		if !slices.Contains(stringsOrEmpty(catalog.modules["core"]["requiredSkills"]), name) {
			t.Errorf("core does not require %s", name)
		}
		for profileID, profile := range catalog.profiles {
			setupID, _ := stringValue(profile, "setup")
			external := false
			for _, skill := range objectsOrEmpty(catalog.setups[setupID]["skills"]) {
				source, _ := objectValue(skill["source"])
				if skill["name"] == name && source["type"] == "github" {
					external = true
				}
			}
			if !external {
				t.Errorf("profile %s setup %s lacks external %s", profileID, setupID, name)
			}
		}
	}
}
