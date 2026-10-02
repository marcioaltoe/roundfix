package baseline

import (
	"fmt"
	"path"
	"sort"
)

// TrailingSetupSkills compares installed required external skill trees with
// the immutable Setup Snapshot used by RestoreSkills. Missing trees are left
// to Repository Skill Set readiness; this comparison does not read the lock.
func TrailingSetupSkills(repoRoot, profileID string, required []string) ([]string, error) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return nil, fmt.Errorf("load catalog for snapshot comparison: %w", err)
	}
	profile, err := loadRestoreProfile(catalog, profileID)
	if err != nil {
		return nil, fmt.Errorf("load profile for snapshot comparison: %w", err)
	}
	trailing := make(map[string]struct{})
	for _, name := range required {
		contract, external := profile.Skills[name]
		_, isRequired := profile.RequiredSkills[name]
		if !external || !isRequired || contract.TreeDigest == "" {
			continue
		}
		files, exists, err := inspectRestoreTarget(repoRoot, path.Join(".agents/skills", name))
		if err != nil {
			return nil, fmt.Errorf("read installed skill %q for snapshot comparison: %w", name, err)
		}
		if exists && portableRestoreDigest(files) != contract.TreeDigest {
			trailing[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(trailing))
	for name := range trailing {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}
