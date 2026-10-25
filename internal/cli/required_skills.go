package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	roundskills "roundfix/skills"
)

// missingRequiredSkills names each external skill the adopted modules require
// that the inspection proved absent from .agents/skills or skills-lock.json,
// or no longer matching its lock entry. An error naming no skill adds nothing.
func missingRequiredSkills(profile string, readiness roundskills.RepositoryReadiness, checkErr error) []baselineUpdateSkillDrift {
	reasons := make(map[string]string)
	add := func(names []string, reason string) {
		for _, name := range names {
			if _, exists := reasons[name]; !exists {
				reasons[name] = reason + "; restore with roundfix baseline skills restore --profile " + profile + " --skill " + name
			}
		}
	}
	const absent = "required by the Baseline and not installed"
	add(readiness.MissingExternal, absent)
	var readinessErr *roundskills.RepositoryReadinessError
	if errors.As(checkErr, &readinessErr) {
		add(readinessErr.MissingExternal, absent)
	}
	add(readiness.OutdatedExternal, "required by the Baseline and no longer matches its skills-lock.json entry")
	var result []baselineUpdateSkillDrift
	for name, reason := range reasons {
		result = append(result, baselineUpdateSkillDrift{Skill: name, Reason: reason})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Skill < result[j].Skill })
	return result
}

func requiredSkillsManifestProfile(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath)))
	if err != nil {
		return "", fmt.Errorf("read Setup Manifest profile: %w", err)
	}
	var manifest struct {
		Profile string `json:"profile"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", fmt.Errorf("read Setup Manifest profile: %w", err)
	}
	return manifest.Profile, nil
}
