package baseline

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	// ErrNoManifest identifies a repository that has not adopted a Baseline.
	ErrNoManifest = errors.New("setup manifest is absent")
	// ErrManifestIncompatible identifies manifest input that requires adoption.
	ErrManifestIncompatible = errors.New("setup manifest is incompatible")
)

// ManifestInputState is the caller-facing outcome of Setup Manifest resolution.
type ManifestInputState string

const (
	// ManifestInputAbsent means the repository has no Setup Manifest.
	ManifestInputAbsent ManifestInputState = "absent"
	// ManifestInputIncompatible means the manifest cannot supply update inputs.
	ManifestInputIncompatible ManifestInputState = "incompatible"
	// ManifestInputResolved means the manifest supplies every required input.
	ManifestInputResolved ManifestInputState = "resolved"
	// ManifestInputIncomplete means the current catalog requires new decisions.
	ManifestInputIncomplete ManifestInputState = "incomplete"
)

// ManifestInputIncompatibility identifies why adoption is required.
type ManifestInputIncompatibility string

const (
	// ManifestInputManifestUnreadable means the manifest path cannot be read safely.
	ManifestInputManifestUnreadable ManifestInputIncompatibility = "manifest_unreadable"
	// ManifestInputManifestInvalid means the manifest does not match the supported schema.
	ManifestInputManifestInvalid ManifestInputIncompatibility = "manifest_invalid"
	// ManifestInputProfileUnresolved means the recorded Baseline Profile no longer resolves.
	ManifestInputProfileUnresolved ManifestInputIncompatibility = "profile_unresolved"
	// ManifestInputDecisionsInvalid means at least one recorded decision no longer validates.
	ManifestInputDecisionsInvalid ManifestInputIncompatibility = "decisions_invalid"
)

// DecisionSuggestion is one required catalog decision absent from the manifest.
type DecisionSuggestion struct {
	ID             string `json:"id"`
	SuggestedValue any    `json:"suggestedValue"`
	Summary        string `json:"summary"`
}

// UnresolvedProfileKind identifies the repair path for one unresolved
// Baseline Profile recorded by the Setup Manifest.
type UnresolvedProfileKind string

const (
	// UnresolvedProfileRepositoryMissing means the repository-owned location
	// did not supply the recorded Profile.
	UnresolvedProfileRepositoryMissing UnresolvedProfileKind = "repository_profile_missing"
	// UnresolvedProfileCatalogUnknown means the identity is neither in the
	// embedded catalog nor a valid repository-owned Profile identity.
	UnresolvedProfileCatalogUnknown UnresolvedProfileKind = "catalog_identity_unknown"
)

// UnresolvedProfileDiagnosis is the actionable projection of a failed
// recorded Profile resolution.
type UnresolvedProfileDiagnosis struct {
	Kind              UnresolvedProfileKind `json:"kind"`
	Identity          string                `json:"identity"`
	SearchedLocations []string              `json:"searchedLocations"`
	Action            string                `json:"action"`
}

func (diagnosis UnresolvedProfileDiagnosis) message() string {
	locations := "none"
	if len(diagnosis.SearchedLocations) != 0 {
		locations = strings.Join(diagnosis.SearchedLocations, ", ")
	}
	return fmt.Sprintf(
		"recorded Baseline Profile %q is unresolved; searched locations: %s; repair action: %s",
		diagnosis.Identity,
		locations,
		diagnosis.Action,
	)
}

type unresolvedProfileDiagnosisError struct {
	diagnosis UnresolvedProfileDiagnosis
	cause     error
}

func (err *unresolvedProfileDiagnosisError) Error() string {
	return err.diagnosis.message()
}

func (err *unresolvedProfileDiagnosisError) Unwrap() error {
	return err.cause
}

// ManifestInput is a stored Setup Manifest projected into plan inputs.
type ManifestInput struct {
	State                ManifestInputState           `json:"state"`
	Incompatibility      ManifestInputIncompatibility `json:"incompatibility,omitempty"`
	ProfileID            string                       `json:"profileId,omitempty"`
	ProfileDigestChanged bool                         `json:"profileDigestChanged,omitempty"`
	Decisions            []DecisionValue              `json:"decisions"`
	NewDecisions         []DecisionSuggestion         `json:"newDecisions"`
	UnresolvedProfile    *UnresolvedProfileDiagnosis  `json:"unresolvedProfile,omitempty"`
	Manifest             SetupManifest                `json:"manifest"`
}

// ResolveManifestInput reads docs/agents/setup-context.json and resolves it
// against the current catalog. It never prompts and never writes.
func ResolveManifestInput(root string, catalog *Catalog) (ManifestInput, error) {
	if catalog == nil {
		return ManifestInput{}, errors.New("resolve Setup Manifest input: catalog is required")
	}
	rootPath, err := cleanRepositoryRoot(root)
	if err != nil {
		return ManifestInput{}, fmt.Errorf("resolve Setup Manifest input repository: %w", err)
	}
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return ManifestInput{}, fmt.Errorf("open repository root for Setup Manifest input: %w", err)
	}
	defer anchored.Close()

	_, err = anchored.Lstat(filepath.FromSlash(manifestPath))
	if errors.Is(err, fs.ErrNotExist) {
		return ManifestInput{
			State:        ManifestInputAbsent,
			Decisions:    []DecisionValue{},
			NewDecisions: []DecisionSuggestion{},
		}, ErrNoManifest
	}
	if err != nil {
		return incompatibleManifestInput(
			ManifestInput{},
			ManifestInputManifestUnreadable,
			fmt.Errorf("inspect Setup Manifest: %w", err),
		)
	}
	data, err := readRootRegularFile(anchored, manifestPath)
	if err != nil {
		return incompatibleManifestInput(
			ManifestInput{},
			ManifestInputManifestUnreadable,
			err,
		)
	}
	manifest, valid := parseManagedSetupManifest(data)
	if !valid {
		return incompatibleManifestInput(
			ManifestInput{},
			ManifestInputManifestInvalid,
			errors.New("manifest does not match the supported Setup Manifest schema"),
		)
	}

	input := ManifestInput{
		ProfileID:    manifest.Profile,
		Decisions:    []DecisionValue{},
		NewDecisions: []DecisionSuggestion{},
		Manifest:     manifest,
	}
	profile, err := ResolveProfile(rootPath, manifest.Profile, catalog)
	if err != nil {
		diagnosis := diagnoseUnresolvedProfile(manifest.Profile, err)
		input.UnresolvedProfile = &diagnosis
		return incompatibleManifestInput(
			input,
			ManifestInputProfileUnresolved,
			&unresolvedProfileDiagnosisError{diagnosis: diagnosis, cause: err},
		)
	}
	input.ProfileDigestChanged = manifest.ProfileDigest != profile.Digest

	decisions, missing, err := resolveManifestDecisionValues(manifest, profile, catalog)
	if err != nil {
		return incompatibleManifestInput(
			input,
			ManifestInputDecisionsInvalid,
			err,
		)
	}
	input.Decisions = decisions
	if len(missing) == 0 {
		input.State = ManifestInputResolved
		return input, nil
	}
	input.NewDecisions, err = resolveDecisionSuggestions(missing, catalog)
	if err != nil {
		return incompatibleManifestInput(
			input,
			ManifestInputDecisionsInvalid,
			err,
		)
	}
	input.State = ManifestInputIncomplete
	return input, nil
}

func diagnoseUnresolvedProfile(profileID string, err error) UnresolvedProfileDiagnosis {
	diagnosis := UnresolvedProfileDiagnosis{
		Kind:              UnresolvedProfileCatalogUnknown,
		Identity:          profileID,
		SearchedLocations: []string{},
		Action: "run roundfix baseline to adopt a Baseline Profile known to the current catalog, " +
			"then rerun roundfix baseline update",
	}
	var resolutionErr *ProfileResolutionError
	if !errors.As(err, &resolutionErr) || len(resolutionErr.SearchedLocations) == 0 {
		return diagnosis
	}
	diagnosis.Kind = UnresolvedProfileRepositoryMissing
	diagnosis.SearchedLocations = append(diagnosis.SearchedLocations, resolutionErr.SearchedLocations...)
	diagnosis.Action = fmt.Sprintf(
		"restore the repository-owned Baseline Profile at %s, then rerun roundfix baseline update",
		strings.Join(diagnosis.SearchedLocations, ", "),
	)
	return diagnosis
}

func incompatibleManifestInput(
	input ManifestInput,
	reason ManifestInputIncompatibility,
	cause error,
) (ManifestInput, error) {
	input.State = ManifestInputIncompatible
	input.Incompatibility = reason
	if input.Decisions == nil {
		input.Decisions = []DecisionValue{}
	}
	if input.NewDecisions == nil {
		input.NewDecisions = []DecisionSuggestion{}
	}
	return input, fmt.Errorf("%w: %w", ErrManifestIncompatible, cause)
}

func resolveManifestDecisionValues(
	manifest SetupManifest,
	profile ResolvedProfile,
	catalog *Catalog,
) ([]DecisionValue, []string, error) {
	ids := sortedKeys(manifest.Decisions)
	explicit := make([]DecisionValue, 0, len(ids))
	for _, id := range ids {
		stored := manifest.Decisions[id]
		if fixed, ok := profile.Values[id]; ok {
			if !valuesEqual(stored.Value, fixed) {
				return nil, nil, fmt.Errorf(
					"resolve Setup Manifest decision %q: recorded value conflicts with the Baseline Profile",
					id,
				)
			}
			continue
		}
		explicit = append(explicit, DecisionValue{ID: id, Value: cloneJSONValue(stored.Value)})
	}
	decisions, missing, err := normalizePlanDecisions(profile, explicit, catalog)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve Setup Manifest decisions: %w", err)
	}
	return decisions, missing, nil
}

func resolveDecisionSuggestions(ids []string, catalog *Catalog) ([]DecisionSuggestion, error) {
	suggestions := make([]DecisionSuggestion, 0, len(ids))
	for _, id := range ids {
		declaration, ok := catalog.decisions[id]
		if !ok {
			return nil, fmt.Errorf("resolve Setup Manifest suggestion: unknown decision %q", id)
		}
		summary, ok := stringValue(declaration, "summary")
		if !ok || summary == "" {
			return nil, fmt.Errorf("resolve Setup Manifest suggestion: decision %q has no summary", id)
		}
		value, ok := declaration["suggestion"]
		if !ok {
			value, ok = declaration["default"]
		}
		if !ok {
			return nil, fmt.Errorf("resolve Setup Manifest suggestion: decision %q has no suggested value", id)
		}
		if err := validateDecisionValue(declaration, value); err != nil {
			return nil, fmt.Errorf("resolve Setup Manifest suggestion for %q: %w", id, err)
		}
		suggestions = append(suggestions, DecisionSuggestion{
			ID:             id,
			SuggestedValue: cloneJSONValue(value),
			Summary:        summary,
		})
	}
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].ID < suggestions[j].ID
	})
	return suggestions, nil
}
