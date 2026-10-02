package baseline

import (
	"fmt"
	"reflect"
)

// composeSetupSnapshot preserves component order and refuses unequal entries
// sharing a skill name or activation bundle identifier.
func composeSetupSnapshot(id string, components []document) (document, error) {
	skills := []any{}
	bundles := []any{}
	seenSkills := map[string]document{}
	seenBundles := map[string]document{}
	setupIDs := []any{}
	for _, component := range components {
		setupIDs = append(setupIDs, component["id"])
		for _, collection := range []struct {
			field  string
			key    string
			seen   map[string]document
			target *[]any
		}{
			{"skills", "name", seenSkills, &skills},
			{"activationBundles", "id", seenBundles, &bundles},
		} {
			for _, entry := range objectsOrEmpty(component[collection.field]) {
				name, _ := stringValue(entry, collection.key)
				if prior, exists := collection.seen[name]; exists {
					if !reflect.DeepEqual(prior, entry) {
						return nil, fmt.Errorf("components disagree on %s %q", collection.field, name)
					}
					continue
				}
				collection.seen[name] = entry
				*collection.target = append(*collection.target, map[string]any(entry))
			}
		}
	}
	result := document{
		"schemaVersion": "setup-context-driven/setup-snapshot/0.0.1",
		"id":            id,
		"version":       "0.0.1",
		"source":        map[string]any{"type": "composed", "setups": setupIDs},
		"skills":        skills,
	}
	var payload any = skills
	if len(bundles) > 0 {
		result["activationBundles"] = bundles
		payload = document{"skills": skills, "activationBundles": bundles}
	}
	digest, err := canonicalSHA256(payload)
	if err != nil {
		return nil, fmt.Errorf("calculate composed setup digest: %w", err)
	}
	result["digest"] = digest
	return result, nil
}

func setupCompositionComponents(setup document, setups map[string]document) ([]document, error) {
	source, _ := objectValue(setup["source"])
	ids, ok := stringList(source["setups"])
	if !ok || len(ids) < 2 || !uniqueStrings(ids) {
		return nil, fmt.Errorf("composition requires at least two distinct setup identifiers")
	}
	components := make([]document, 0, len(ids))
	for _, id := range ids {
		component, exists := setups[id]
		if !exists {
			return nil, fmt.Errorf("unknown component %q", id)
		}
		componentSource, _ := objectValue(component["source"])
		if componentSource["type"] == "composed" {
			return nil, fmt.Errorf("component %q is composed", id)
		}
		components = append(components, component)
	}
	return components, nil
}
