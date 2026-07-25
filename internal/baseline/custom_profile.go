package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	CustomProfileSchemaVersion   = "roundfix/custom-baseline-profile/v1"
	ResolvedProfileSchemaVersion = "roundfix/baseline-profile/v1"

	customProfileDigestDomain = ResolvedProfileSchemaVersion + "\x00"
	customProfileDirectory    = ".roundfix/baseline/profiles"
)

var ErrUnsafeCustomProfilePath = errors.New("unsafe custom Baseline Profile path")

type BaselineProfileSource string

const (
	ProfileSourceBuiltIn    BaselineProfileSource = "built-in"
	ProfileSourceRepository BaselineProfileSource = "repository"
)

// ResolvedProfile is one validated Baseline Profile normalized against the
// embedded catalog.
type ResolvedProfile struct {
	SchemaVersion string                `json:"schemaVersion"`
	CatalogSchema string                `json:"catalogSchema"`
	ID            string                `json:"id"`
	Source        BaselineProfileSource `json:"source"`
	Modules       []string              `json:"modules"`
	Decisions     []string              `json:"decisions"`
	Capabilities  []string              `json:"capabilities"`
	Templates     []string              `json:"templates"`
	Values        map[string]any        `json:"values"`
	Digest        string                `json:"digest"`
	Path          string                `json:"path,omitempty"`
}

// IdentifierStrategy is the strict project decision for new project-owned
// Internal Identifiers.
type IdentifierStrategy struct {
	Kind     string `json:"kind"`
	Guidance string `json:"guidance,omitempty"`
}

// HTTPException is one provider-owned HTTP route exception.
type HTTPException struct {
	Scope   string   `json:"scope"`
	Methods []string `json:"methods"`
	Owner   string   `json:"owner"`
	Reason  string   `json:"reason"`
}

// AuthProviderDecision is the strict authentication-provider decision.
type AuthProviderDecision struct {
	Kind           string        `json:"kind"`
	RouteException HTTPException `json:"routeException"`
}

// CustomProfileInitResult identifies the repository file created by init.
type CustomProfileInitResult struct {
	Path        string
	FromProfile string
	Profile     ResolvedProfile
}

type customProfileDocument struct {
	SchemaVersion string         `json:"schemaVersion"`
	CatalogSchema string         `json:"catalogSchema"`
	ID            string         `json:"id"`
	Modules       []string       `json:"modules"`
	Decisions     []string       `json:"decisions"`
	Capabilities  []string       `json:"capabilities"`
	Templates     []string       `json:"templates"`
	Values        map[string]any `json:"values"`
}

type profileDigestPayload struct {
	SchemaVersion string         `json:"schemaVersion"`
	CatalogSchema string         `json:"catalogSchema"`
	ID            string         `json:"id"`
	Modules       []string       `json:"modules"`
	Decisions     []string       `json:"decisions"`
	Capabilities  []string       `json:"capabilities"`
	Templates     []string       `json:"templates"`
	Values        map[string]any `json:"values"`
}

// ProfileDraftInputFromDocument binds one strict custom Profile document to
// its only compatible built-in source Profile.
func ProfileDraftInputFromDocument(document []byte, catalog *Catalog) (ProfileDraftInput, error) {
	if catalog == nil {
		return ProfileDraftInput{}, errors.New("resolve custom Profile draft source: catalog is required")
	}
	draft, err := ParseCustomProfile(document, "Profile draft", catalog)
	if err != nil {
		return ProfileDraftInput{}, err
	}
	var sources []string
	for _, sourceID := range catalog.ProfileIDs() {
		source, resolveErr := resolveBuiltInProfile(sourceID, catalog)
		if resolveErr != nil {
			return ProfileDraftInput{}, resolveErr
		}
		if validateProfileAdaptation(source, draft, catalog) == nil {
			sources = append(sources, sourceID)
		}
	}
	switch len(sources) {
	case 0:
		return ProfileDraftInput{}, errors.New(
			"custom.profile.draft.source.unresolved: document is not a valid adaptation of any built-in Profile",
		)
	case 1:
		return ProfileDraftInput{
			SourceProfileID: sources[0],
			Document:        append([]byte(nil), document...),
		}, nil
	default:
		return ProfileDraftInput{}, fmt.Errorf(
			"custom.profile.draft.source.ambiguous: document matches built-in Profiles %s",
			strings.Join(sources, ", "),
		)
	}
}

// NewProfileAdaptationDraft creates a strict draft by removing only reviewed
// modules and profile-specific capabilities from one built-in Profile.
func NewProfileAdaptationDraft(
	sourceProfileID string,
	id string,
	removedModules []string,
	removedCapabilities []string,
	catalog *Catalog,
) (ProfileDraftInput, error) {
	if catalog == nil {
		return ProfileDraftInput{}, errors.New("create Profile adaptation draft: catalog is required")
	}
	sourceProfileID = strings.TrimSpace(sourceProfileID)
	source, err := resolveBuiltInProfile(sourceProfileID, catalog)
	if err != nil {
		return ProfileDraftInput{}, fmt.Errorf("create Profile adaptation draft: %w", err)
	}
	if err := validateCustomProfileID(id); err != nil {
		return ProfileDraftInput{}, err
	}
	moduleRemovals, err := reviewedProfileRemovals("module", removedModules, source.Modules)
	if err != nil {
		return ProfileDraftInput{}, err
	}
	capabilityRemovals, err := reviewedProfileRemovals(
		"capability",
		removedCapabilities,
		source.Capabilities,
	)
	if err != nil {
		return ProfileDraftInput{}, err
	}
	if len(moduleRemovals) == 0 && len(capabilityRemovals) == 0 {
		return ProfileDraftInput{}, errors.New("custom.profile.adaptation.removal.required")
	}

	modules := profileValuesWithout(source.Modules, moduleRemovals)
	capabilities := profileValuesWithout(source.Capabilities, capabilityRemovals)
	decisions := profileAdaptationDecisions(source.Decisions, modules, capabilities, catalog)
	values := make(map[string]any)
	selectedDecisions := stringSet(decisions)
	for decisionID, value := range source.Values {
		if _, selected := selectedDecisions[decisionID]; selected {
			values[decisionID] = cloneJSONValue(value)
		}
	}
	document := customProfileDocument{
		SchemaVersion: CustomProfileSchemaVersion,
		CatalogSchema: catalogSchema,
		ID:            strings.TrimSpace(id),
		Modules:       modules,
		Decisions:     decisions,
		Capabilities:  capabilities,
		Templates:     cloneStrings(source.Templates),
		Values:        values,
	}
	data, err := marshalCustomProfileDocument(document)
	if err != nil {
		return ProfileDraftInput{}, fmt.Errorf("serialize Profile adaptation draft: %w", err)
	}
	draft, err := ParseCustomProfile(data, "Profile adaptation draft", catalog)
	if err != nil {
		return ProfileDraftInput{}, err
	}
	if err := validateProfileAdaptation(source, draft, catalog); err != nil {
		return ProfileDraftInput{}, err
	}
	return ProfileDraftInput{SourceProfileID: sourceProfileID, Document: data}, nil
}

// ResolveProfileDraft validates and normalizes one in-memory Profile draft
// without writing its canonical repository target.
func ResolveProfileDraft(
	repoRoot string,
	input ProfileDraftInput,
	catalog *Catalog,
) (ResolvedProfile, []byte, error) {
	return resolveProfileDraft(repoRoot, input, catalog)
}

func reviewedProfileRemovals(kind string, removed, selected []string) (map[string]struct{}, error) {
	available := stringSet(selected)
	result := make(map[string]struct{}, len(removed))
	for _, raw := range removed {
		id := strings.TrimSpace(raw)
		if _, ok := available[id]; !ok {
			return nil, fmt.Errorf("custom.profile.adaptation.%s.removal.invalid: %s", kind, id)
		}
		if _, duplicate := result[id]; duplicate {
			return nil, fmt.Errorf("custom.profile.adaptation.%s.removal.duplicate: %s", kind, id)
		}
		result[id] = struct{}{}
	}
	return result, nil
}

func profileValuesWithout(values []string, removed map[string]struct{}) []string {
	result := make([]string, 0, len(values)-len(removed))
	for _, value := range values {
		if _, remove := removed[value]; !remove {
			result = append(result, value)
		}
	}
	return result
}

func profileAdaptationDecisions(
	sourceDecisions []string,
	selectedModules []string,
	selectedCapabilities []string,
	catalog *Catalog,
) []string {
	selected := stringSet(sourceDecisions)
	modules := stringSet(selectedModules)
	capabilities := stringSet(selectedCapabilities)
	artifacts := profileModuleArtifacts(selectedModules, catalog)
	for changed := true; changed; {
		changed = false
		for decisionID := range selected {
			remove := false
			for _, capabilityID := range stringsOrEmpty(
				catalog.decisions[decisionID]["requiresCapabilities"],
			) {
				if _, retained := capabilities[capabilityID]; !retained {
					remove = true
					break
				}
			}
			for _, effect := range objectsOrEmpty(catalog.decisions[decisionID]["effects"]) {
				for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
					if _, active := modules[moduleID]; !active {
						remove = true
						break
					}
				}
				if remove {
					break
				}
				for _, artifactID := range decisionEffectArtifacts(effect) {
					if _, active := artifacts[artifactID]; !active {
						remove = true
						break
					}
				}
				if remove {
					break
				}
				for _, requiredID := range stringsOrEmpty(effect["requireDecisions"]) {
					if _, retained := selected[requiredID]; !retained {
						remove = true
						break
					}
				}
				if remove {
					break
				}
			}
			if remove {
				delete(selected, decisionID)
				changed = true
			}
		}
	}
	result := make([]string, 0, len(selected))
	for _, decisionID := range sourceDecisions {
		if _, retained := selected[decisionID]; retained {
			result = append(result, decisionID)
		}
	}
	return result
}

// CatalogSchemaVersion returns the schema identity custom profiles bind to.
func CatalogSchemaVersion() string {
	return catalogSchema
}

// InitCustomProfile creates one repository-owned profile from an embedded
// built-in profile without retaining a profile-to-profile reference.
func InitCustomProfile(repoRoot, id, from string, catalog *Catalog) (CustomProfileInitResult, error) {
	if catalog == nil {
		return CustomProfileInitResult{}, errors.New("initialize custom Baseline Profile: catalog is required")
	}
	id = strings.TrimSpace(id)
	from = strings.TrimSpace(from)
	if err := validateCustomProfileID(id); err != nil {
		return CustomProfileInitResult{}, err
	}
	if _, exists := catalog.Profile(id); exists {
		return CustomProfileInitResult{}, fmt.Errorf("custom.profile.id.built-in: %s", id)
	}
	source, err := resolveBuiltInProfile(from, catalog)
	if err != nil {
		return CustomProfileInitResult{}, err
	}

	document := customProfileDocument{
		SchemaVersion: CustomProfileSchemaVersion,
		CatalogSchema: catalogSchema,
		ID:            id,
		Modules:       cloneStrings(source.Modules),
		Decisions:     cloneStrings(source.Decisions),
		Capabilities:  cloneStrings(source.Capabilities),
		Templates:     cloneStrings(source.Templates),
		Values:        map[string]any{},
	}
	data, err := marshalCustomProfileDocument(document)
	if err != nil {
		return CustomProfileInitResult{}, fmt.Errorf("serialize custom Baseline Profile %q: %w", id, err)
	}

	profilePath, err := repositoryProfilePath(repoRoot, id)
	if err != nil {
		return CustomProfileInitResult{}, err
	}
	if err := prepareCustomProfileDirectory(repoRoot); err != nil {
		return CustomProfileInitResult{}, err
	}
	file, err := os.OpenFile(profilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return CustomProfileInitResult{}, fmt.Errorf("create custom Baseline Profile %q: %w", profilePath, err)
	}
	writeErr := error(nil)
	if _, err := file.Write(data); err != nil {
		writeErr = fmt.Errorf("write custom Baseline Profile %q: %w", profilePath, err)
	}
	if err := file.Close(); err != nil && writeErr == nil {
		writeErr = fmt.Errorf("close custom Baseline Profile %q: %w", profilePath, err)
	}
	if writeErr != nil {
		_ = os.Remove(profilePath)
		return CustomProfileInitResult{}, writeErr
	}

	resolved, err := ParseCustomProfile(data, profilePath, catalog)
	if err != nil {
		_ = os.Remove(profilePath)
		return CustomProfileInitResult{}, fmt.Errorf("validate initialized custom Baseline Profile: %w", err)
	}
	resolved.Path = profilePath
	return CustomProfileInitResult{Path: profilePath, FromProfile: from, Profile: resolved}, nil
}

func resolveProfileDraft(
	repoRoot string,
	input ProfileDraftInput,
	catalog *Catalog,
) (ResolvedProfile, []byte, error) {
	sourceID := strings.TrimSpace(input.SourceProfileID)
	if sourceID == "" {
		return ResolvedProfile{}, nil, errors.New("custom.profile.draft.source.required")
	}
	source, err := resolveBuiltInProfile(sourceID, catalog)
	if err != nil {
		return ResolvedProfile{}, nil, fmt.Errorf("resolve custom Profile draft source: %w", err)
	}
	profile, err := ParseCustomProfile(input.Document, "Profile draft", catalog)
	if err != nil {
		return ResolvedProfile{}, nil, err
	}
	if err := validateProfileAdaptation(source, profile, catalog); err != nil {
		return ResolvedProfile{}, nil, err
	}
	profile.Templates = profileTemplateIDs(catalog, profile.Modules, profile.Decisions)
	relative := path.Join(customProfileDirectory, profile.ID+".json")
	if !safeRelative(relative) {
		return ResolvedProfile{}, nil, fmt.Errorf("%w: %s", ErrUnsafeCustomProfilePath, relative)
	}
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return ResolvedProfile{}, nil, err
	}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return ResolvedProfile{}, nil, fmt.Errorf("open repository for custom Profile draft: %w", err)
	}
	defer anchored.Close()
	if err := validatePathParents(anchored, relative); err != nil {
		return ResolvedProfile{}, nil, fmt.Errorf("%w: %s: %v", ErrUnsafeCustomProfilePath, relative, err)
	}
	if err := validateMutationDestination(anchored, relative); err != nil {
		return ResolvedProfile{}, nil, fmt.Errorf("%w: %s: %v", ErrUnsafeCustomProfilePath, relative, err)
	}
	profile.Path = relative
	profile.Digest, err = profileDigest(profile)
	if err != nil {
		return ResolvedProfile{}, nil, err
	}
	data, err := marshalResolvedCustomProfile(profile)
	if err != nil {
		return ResolvedProfile{}, nil, err
	}
	return profile, data, nil
}

func validateProfileAdaptation(source, draft ResolvedProfile, catalog *Catalog) error {
	if err := validateProfileFieldSubset("modules", draft.Modules, source.Modules); err != nil {
		return err
	}
	if err := validateProfileFieldSubset("decisions", draft.Decisions, source.Decisions); err != nil {
		return err
	}
	if err := validateProfileFieldSubset("capabilities", draft.Capabilities, source.Capabilities); err != nil {
		return err
	}
	selectedModules := stringSet(draft.Modules)
	selectedDecisions := stringSet(draft.Decisions)
	selectedArtifacts := profileModuleArtifacts(draft.Modules, catalog)
	for _, moduleID := range draft.Modules {
		for _, decisionID := range stringsOrEmpty(catalog.modules[moduleID]["requiredDecisions"]) {
			if _, included := selectedDecisions[decisionID]; !included {
				return fmt.Errorf(
					"custom.profile.adaptation.decision.required: module %s requires %s",
					moduleID,
					decisionID,
				)
			}
		}
	}
	for _, decisionID := range draft.Decisions {
		for _, effect := range objectsOrEmpty(catalog.decisions[decisionID]["effects"]) {
			for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
				if _, included := selectedModules[moduleID]; !included {
					return fmt.Errorf(
						"custom.profile.adaptation.module.decision-missing: decision %s can activate %s",
						decisionID,
						moduleID,
					)
				}
			}
			for _, artifactID := range decisionEffectArtifacts(effect) {
				if _, included := selectedArtifacts[artifactID]; !included {
					return fmt.Errorf(
						"custom.profile.adaptation.module.decision-missing: decision %s targets %s from a removed module",
						decisionID,
						artifactID,
					)
				}
			}
			for _, requiredID := range stringsOrEmpty(effect["requireDecisions"]) {
				if _, included := selectedDecisions[requiredID]; !included {
					return fmt.Errorf(
						"custom.profile.adaptation.decision.dependency-missing: %s requires %s",
						decisionID,
						requiredID,
					)
				}
			}
		}
	}
	if err := validateConditionalProfileDecisions(
		draft.Decisions,
		draft.Capabilities,
		catalog,
	); err != nil {
		return err
	}
	return nil
}

func profileModuleArtifacts(modules []string, catalog *Catalog) map[string]struct{} {
	artifacts := make(map[string]struct{})
	for _, moduleID := range modules {
		module := catalog.modules[moduleID]
		for _, field := range []string{"rootBlocks", "supportingGuides", "repositoryExtensions"} {
			for _, declaration := range objectsOrEmpty(module[field]) {
				if artifactID, ok := stringValue(declaration, "id"); ok {
					artifacts[artifactID] = struct{}{}
				}
			}
		}
	}
	return artifacts
}

func decisionEffectArtifacts(effect document) []string {
	var artifacts []string
	artifacts = append(artifacts, stringsOrEmpty(effect["includeArtifacts"])...)
	for _, field := range []string{"selectTemplates", "renderBindings"} {
		for _, declaration := range objectsOrEmpty(effect[field]) {
			if artifactID, ok := stringValue(declaration, "artifact"); ok {
				artifacts = append(artifacts, artifactID)
			}
		}
	}
	return artifacts
}

func validateProfileFieldSubset(field string, selected, source []string) error {
	selectedSet := stringSet(selected)
	normalized := make([]string, 0, len(selected))
	for _, id := range source {
		if _, included := selectedSet[id]; included {
			normalized = append(normalized, id)
			delete(selectedSet, id)
		}
	}
	if len(selectedSet) != 0 {
		return fmt.Errorf(
			"custom.profile.adaptation.%s.addition: %s",
			field,
			strings.Join(sortedKeys(selectedSet), ", "),
		)
	}
	if !equalStringLists(selected, normalized) {
		return fmt.Errorf("custom.profile.adaptation.%s.order.invalid", field)
	}
	return nil
}

func equalStringLists(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func marshalResolvedCustomProfile(profile ResolvedProfile) ([]byte, error) {
	return marshalCustomProfileDocument(customProfileDocument{
		SchemaVersion: CustomProfileSchemaVersion,
		CatalogSchema: profile.CatalogSchema,
		ID:            profile.ID,
		Modules:       cloneStrings(profile.Modules),
		Decisions:     cloneStrings(profile.Decisions),
		Capabilities:  cloneStrings(profile.Capabilities),
		Templates:     cloneStrings(profile.Templates),
		Values:        cloneJSONMap(profile.Values),
	})
}

func marshalCustomProfileDocument(document customProfileDocument) ([]byte, error) {
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// ParseCustomProfile strictly parses and resolves one repository declaration.
func ParseCustomProfile(data []byte, sourcePath string, catalog *Catalog) (ResolvedProfile, error) {
	if catalog == nil {
		return ResolvedProfile{}, errors.New("parse custom Baseline Profile: catalog is required")
	}
	raw, diagnostics := decodeDocument(data, sourcePath)
	if len(diagnostics) != 0 {
		parts := make([]string, 0, len(diagnostics))
		for _, diagnostic := range diagnostics {
			code := strings.Replace(diagnostic.Code, "catalog.json.", "custom.profile.json.", 1)
			parts = append(parts, Diagnostic{Code: code, Path: diagnostic.Path, Info: diagnostic.Info}.String())
		}
		return ResolvedProfile{}, errors.New(strings.Join(parts, "; "))
	}
	if raw == nil {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.json.object.required: %s", sourcePath)
	}
	allowed := stringSet([]string{
		"schemaVersion",
		"catalogSchema",
		"id",
		"modules",
		"decisions",
		"capabilities",
		"templates",
		"values",
	})
	for key := range raw {
		if _, ok := allowed[key]; !ok {
			return ResolvedProfile{}, fmt.Errorf("custom.profile.field.unknown: %s: %s", sourcePath, key)
		}
	}
	for _, field := range []string{
		"schemaVersion",
		"catalogSchema",
		"id",
		"modules",
		"decisions",
		"capabilities",
		"templates",
		"values",
	} {
		if _, ok := raw[field]; !ok {
			return ResolvedProfile{}, fmt.Errorf("custom.profile.field.missing: %s: %s", sourcePath, field)
		}
	}
	schema, _ := raw["schemaVersion"].(string)
	if schema != CustomProfileSchemaVersion {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.schema.invalid: %s: got %q, want %q", sourcePath, schema, CustomProfileSchemaVersion)
	}
	boundCatalog, _ := raw["catalogSchema"].(string)
	if boundCatalog != catalogSchema {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.catalog-schema.invalid: %s: got %q, want %q", sourcePath, boundCatalog, catalogSchema)
	}
	id, _ := raw["id"].(string)
	if err := validateCustomProfileID(id); err != nil {
		return ResolvedProfile{}, err
	}
	if _, exists := catalog.Profile(id); exists {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.id.built-in: %s", id)
	}
	modules, err := requiredUniqueStringList(raw, "modules", sourcePath, true)
	if err != nil {
		return ResolvedProfile{}, err
	}
	decisions, err := requiredUniqueStringList(raw, "decisions", sourcePath, true)
	if err != nil {
		return ResolvedProfile{}, err
	}
	capabilities, err := requiredUniqueStringList(raw, "capabilities", sourcePath, false)
	if err != nil {
		return ResolvedProfile{}, err
	}
	templates, err := requiredUniqueStringList(raw, "templates", sourcePath, false)
	if err != nil {
		return ResolvedProfile{}, err
	}
	values, ok := objectValue(raw["values"])
	if !ok {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.values.invalid: %s", sourcePath)
	}
	if err := validateCustomProfileReferences(catalog, modules, decisions, capabilities, templates, values); err != nil {
		return ResolvedProfile{}, err
	}

	profile := ResolvedProfile{
		SchemaVersion: ResolvedProfileSchemaVersion,
		CatalogSchema: catalogSchema,
		ID:            id,
		Source:        ProfileSourceRepository,
		Modules:       cloneStrings(modules),
		Decisions:     cloneStrings(decisions),
		Capabilities:  cloneStrings(capabilities),
		Templates:     cloneStrings(templates),
		Values:        cloneJSONMap(values),
		Path:          sourcePath,
	}
	digest, err := profileDigest(profile)
	if err != nil {
		return ResolvedProfile{}, err
	}
	profile.Digest = digest
	return profile, nil
}

// LoadRepositoryProfile resolves one ID only from the repository profile
// directory. No user-scoped location participates.
func LoadRepositoryProfile(repoRoot, id string, catalog *Catalog) (ResolvedProfile, error) {
	profilePath, err := repositoryProfilePath(repoRoot, id)
	if err != nil {
		return ResolvedProfile{}, err
	}
	return LoadCustomProfilePath(repoRoot, profilePath, catalog)
}

// LoadCustomProfilePath resolves an explicit path only when it names one
// direct regular file in the repository-owned profile directory.
func LoadCustomProfilePath(repoRoot, profilePath string, catalog *Catalog) (ResolvedProfile, error) {
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return ResolvedProfile{}, err
	}
	if !filepath.IsAbs(profilePath) {
		profilePath = filepath.Join(root, profilePath)
	}
	profilePath = filepath.Clean(profilePath)
	profilesRoot := filepath.Join(root, filepath.FromSlash(customProfileDirectory))
	relative, err := filepath.Rel(profilesRoot, profilePath)
	if err != nil || relative == "." || filepath.Dir(relative) != "." || filepath.Ext(relative) != ".json" {
		return ResolvedProfile{}, fmt.Errorf("%w: %s", ErrUnsafeCustomProfilePath, profilePath)
	}
	if err := rejectSymlinkComponents(root, profilePath); err != nil {
		return ResolvedProfile{}, err
	}
	info, err := os.Lstat(profilePath)
	if err != nil {
		return ResolvedProfile{}, fmt.Errorf("read custom Baseline Profile %q: %w", profilePath, err)
	}
	if !info.Mode().IsRegular() {
		return ResolvedProfile{}, fmt.Errorf("%w: %s is not a regular file", ErrUnsafeCustomProfilePath, profilePath)
	}
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return ResolvedProfile{}, fmt.Errorf("read custom Baseline Profile %q: %w", profilePath, err)
	}
	profile, err := ParseCustomProfile(data, profilePath, catalog)
	if err != nil {
		return ResolvedProfile{}, err
	}
	fileID := strings.TrimSuffix(filepath.Base(profilePath), ".json")
	if profile.ID != fileID {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.id.path-mismatch: %s: got %q, want %q", profilePath, profile.ID, fileID)
	}
	profile.Path = profilePath
	return profile, nil
}

// DiscoverRepositoryProfiles loads direct JSON children of the repository
// profile directory in lexical order.
func DiscoverRepositoryProfiles(repoRoot string, catalog *Catalog) ([]ResolvedProfile, error) {
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return nil, err
	}
	profilesRoot := filepath.Join(root, filepath.FromSlash(customProfileDirectory))
	if err := rejectSymlinkComponents(root, profilesRoot); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	entries, err := os.ReadDir(profilesRoot)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("discover custom Baseline Profiles: %w", err)
	}
	profiles := make([]ResolvedProfile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		profile, err := LoadCustomProfilePath(root, filepath.Join(profilesRoot, entry.Name()), catalog)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// ResolveProfile resolves a built-in ID or one exact repository-owned ID.
func ResolveProfile(repoRoot, id string, catalog *Catalog) (ResolvedProfile, error) {
	if catalog == nil {
		return ResolvedProfile{}, errors.New("resolve Baseline Profile: catalog is required")
	}
	id = strings.TrimSpace(id)
	if _, ok := catalog.Profile(id); ok {
		return resolveBuiltInProfile(id, catalog)
	}
	return LoadRepositoryProfile(repoRoot, id, catalog)
}

func resolveBuiltInProfile(id string, catalog *Catalog) (ResolvedProfile, error) {
	entry, ok := catalog.Profile(id)
	if !ok {
		return ResolvedProfile{}, fmt.Errorf("custom.profile.source.unknown: %s", id)
	}
	raw, diagnostics := decodeDocument(entry.Data, "embedded profile "+id)
	if len(diagnostics) != 0 || raw == nil {
		return ResolvedProfile{}, fmt.Errorf("resolve embedded Baseline Profile %q", id)
	}
	modules, _ := stringList(raw["modules"])
	decisions, _ := stringList(raw["entryDecisions"])
	capabilities := profileCapabilityIDs(raw)
	templates := profileTemplateIDs(catalog, modules, decisions)
	profile := ResolvedProfile{
		SchemaVersion: ResolvedProfileSchemaVersion,
		CatalogSchema: catalogSchema,
		ID:            id,
		Source:        ProfileSourceBuiltIn,
		Modules:       append([]string(nil), modules...),
		Decisions:     append([]string(nil), decisions...),
		Capabilities:  capabilities,
		Templates:     templates,
		Values:        map[string]any{},
	}
	digest, err := profileDigest(profile)
	if err != nil {
		return ResolvedProfile{}, err
	}
	profile.Digest = digest
	return profile, nil
}

func validateCustomProfileReferences(
	catalog *Catalog,
	modules, decisions, capabilities, templates []string,
	values document,
) error {
	selectedModules := stringSet(modules)
	positions := make(map[string]int, len(modules))
	for index, id := range modules {
		positions[id] = index
		module, exists := catalog.modules[id]
		if !exists {
			return fmt.Errorf("custom.profile.module.unknown: %s", id)
		}
		for _, dependency := range stringsOrEmpty(module["dependsOn"]) {
			dependencyPosition, included := positions[dependency]
			if !included || dependencyPosition >= index {
				return fmt.Errorf("custom.profile.module.dependency.invalid: %s -> %s", id, dependency)
			}
		}
		for _, conflict := range stringsOrEmpty(module["conflictsWith"]) {
			if _, included := selectedModules[conflict]; included {
				return fmt.Errorf("custom.profile.module.conflict: %s -> %s", id, conflict)
			}
		}
	}
	selectedDecisions := stringSet(decisions)
	for _, id := range decisions {
		if _, exists := catalog.decisions[id]; !exists {
			return fmt.Errorf("custom.profile.decision.unknown: %s", id)
		}
	}
	knownCapabilities := catalogCapabilityIDs(catalog)
	for _, id := range capabilities {
		if universalCapability(id) {
			return fmt.Errorf("custom.profile.capability.universal: %s", id)
		}
		if _, exists := knownCapabilities[id]; !exists {
			return fmt.Errorf("custom.profile.capability.unknown: %s", id)
		}
	}
	for _, id := range templates {
		if _, exists := catalog.templates[id]; !exists {
			return fmt.Errorf("custom.profile.template.unknown: %s", id)
		}
	}
	for id, value := range values {
		if _, selected := selectedDecisions[id]; !selected {
			return fmt.Errorf("custom.profile.value.decision.unselected: %s", id)
		}
		if err := validateDecisionValue(catalog.decisions[id], value); err != nil {
			return fmt.Errorf("custom.profile.value.invalid: %s: %w", id, err)
		}
	}
	if err := validateConditionalProfileDecisions(decisions, capabilities, catalog); err != nil {
		return err
	}
	return nil
}

func universalCapability(id string) bool {
	for _, capability := range universalCapabilities {
		if capability.ID == id {
			return true
		}
	}
	return false
}

func validateDecisionValue(decision document, value any) error {
	kind, _ := stringValue(decision, "type")
	switch kind {
	case "boolean":
		if _, ok := value.(bool); !ok {
			return errors.New("must be a boolean")
		}
	case "enum":
		text, ok := value.(string)
		if !ok {
			return errors.New("must be a string")
		}
		for _, allowed := range stringsOrEmpty(decision["values"]) {
			if text == allowed {
				return nil
			}
		}
		return fmt.Errorf("%q is not an allowed value", text)
	case "string":
		text, ok := value.(string)
		if !ok || strings.TrimSpace(text) == "" {
			return errors.New("must be a non-empty string")
		}
	case "http-contract":
		_, err := normalizeHTTPContract(value, decision)
		return err
	case "identifier-strategy":
		return validateIdentifierStrategy(value)
	case "auth-provider":
		return validateAuthProviderDecision(value)
	default:
		return fmt.Errorf("decision type %q is unsupported", kind)
	}
	return nil
}

// ValidateDecisionValue validates one value against its embedded catalog
// declaration.
func ValidateDecisionValue(catalog *Catalog, id string, value any) error {
	if catalog == nil {
		return errors.New("validate Baseline decision: catalog is required")
	}
	declaration, ok := catalog.decisions[id]
	if !ok {
		return fmt.Errorf("validate Baseline decision: unknown decision %q", id)
	}
	if err := validateDecisionValue(declaration, value); err != nil {
		return fmt.Errorf("validate Baseline decision %q: %w", id, err)
	}
	return nil
}

func validateIdentifierStrategy(value any) error {
	object, ok := objectValue(value)
	if !ok {
		return errors.New("must be an object")
	}
	kind, _ := object["kind"].(string)
	switch kind {
	case "uuid-v7":
		if !hasExactFields(object, "kind") {
			return errors.New("uuid-v7 requires exactly kind")
		}
	case "repository-defined":
		if !hasExactFields(object, "kind", "guidance") {
			return errors.New("repository-defined requires exactly kind and guidance")
		}
		guidance, ok := object["guidance"].(string)
		if !ok || strings.TrimSpace(guidance) == "" {
			return errors.New("repository-defined guidance must be a non-empty string")
		}
	default:
		return fmt.Errorf("kind %q is not allowed", kind)
	}
	return nil
}

func validateAuthProviderDecision(value any) error {
	_, err := normalizeAuthProviderDecision(value)
	return err
}

func decisionStringList(value any) ([]string, bool) {
	switch values := value.(type) {
	case []string:
		return append([]string(nil), values...), true
	case []any:
		result := make([]string, 0, len(values))
		for _, value := range values {
			text, ok := value.(string)
			if !ok || text == "" {
				return nil, false
			}
			result = append(result, text)
		}
		return result, true
	default:
		return nil, false
	}
}

func validateConditionalProfileDecisions(
	decisions []string,
	capabilities []string,
	catalog *Catalog,
) error {
	selected := stringSet(decisions)
	retained := stringSet(capabilities)
	for _, decisionID := range catalog.DecisionIDs() {
		declaration := catalog.decisions[decisionID]
		requiredCapabilities := stringsOrEmpty(declaration["requiresCapabilities"])
		if len(requiredCapabilities) == 0 {
			continue
		}
		applicable := true
		for _, capabilityID := range requiredCapabilities {
			if _, ok := retained[capabilityID]; !ok {
				applicable = false
				break
			}
		}
		_, included := selected[decisionID]
		switch {
		case applicable && !included:
			return fmt.Errorf(
				"custom.profile.decision.capability.required: %s requires %s",
				strings.Join(requiredCapabilities, ", "),
				decisionID,
			)
		case !applicable && included:
			return fmt.Errorf(
				"custom.profile.decision.capability.unselected: %s requires %s",
				decisionID,
				strings.Join(requiredCapabilities, ", "),
			)
		}
	}
	return nil
}

func requiredUniqueStringList(raw document, field, sourcePath string, requireNonEmpty bool) ([]string, error) {
	values, ok := stringList(raw[field])
	if !ok || (requireNonEmpty && len(values) == 0) {
		return nil, fmt.Errorf("custom.profile.%s.invalid: %s", field, sourcePath)
	}
	if !uniqueStrings(values) {
		return nil, fmt.Errorf("custom.profile.%s.duplicate: %s", field, sourcePath)
	}
	return values, nil
}

func profileCapabilityIDs(profile document) []string {
	ids := make([]string, 0)
	for _, capability := range objectsOrEmpty(profile["capabilities"]) {
		if id, ok := stringValue(capability, "id"); ok {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func catalogCapabilityIDs(catalog *Catalog) map[string]struct{} {
	ids := make(map[string]struct{})
	for _, profile := range catalog.profiles {
		for _, id := range profileCapabilityIDs(profile) {
			ids[id] = struct{}{}
		}
	}
	return ids
}

func profileTemplateIDs(catalog *Catalog, modules, decisions []string) []string {
	seen := make(map[string]struct{})
	add := func(value any) {
		var object document
		switch typed := value.(type) {
		case document:
			object = typed
		default:
			var ok bool
			object, ok = objectValue(value)
			if !ok {
				return
			}
		}
		if id, ok := stringValue(object, "template"); ok {
			seen[id] = struct{}{}
		}
	}
	for _, moduleID := range modules {
		module := catalog.modules[moduleID]
		for _, artifact := range objectsOrEmpty(module["rootBlocks"]) {
			add(artifact)
		}
		for _, artifact := range objectsOrEmpty(module["supportingGuides"]) {
			add(artifact)
		}
		for _, extension := range objectsOrEmpty(module["repositoryExtensions"]) {
			add(extension)
		}
	}
	for _, decisionID := range decisions {
		for _, effect := range objectsOrEmpty(catalog.decisions[decisionID]["effects"]) {
			for _, selection := range objectsOrEmpty(effect["selectTemplates"]) {
				add(selection)
			}
			for _, binding := range objectsOrEmpty(effect["renderBindings"]) {
				add(binding)
			}
		}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func profileDigest(profile ResolvedProfile) (string, error) {
	payload := profileDigestPayload{
		SchemaVersion: profile.SchemaVersion,
		CatalogSchema: profile.CatalogSchema,
		ID:            profile.ID,
		Modules:       profile.Modules,
		Decisions:     profile.Decisions,
		Capabilities:  profile.Capabilities,
		Templates:     profile.Templates,
		Values:        profile.Values,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("serialize Baseline Profile digest: %w", err)
	}
	sum := sha256.Sum256(append([]byte(customProfileDigestDomain), data...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func validateCustomProfileID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" || len(id) > 128 {
		return fmt.Errorf("custom.profile.id.invalid: %q", id)
	}
	for index, character := range id {
		isAlphaNumeric := character >= 'a' && character <= 'z' || character >= '0' && character <= '9'
		if !isAlphaNumeric && character != '-' {
			return fmt.Errorf("custom.profile.id.invalid: %q", id)
		}
		if character == '-' && (index == 0 || index == len(id)-1) {
			return fmt.Errorf("custom.profile.id.invalid: %q", id)
		}
	}
	return nil
}

func repositoryProfilePath(repoRoot, id string) (string, error) {
	if err := validateCustomProfileID(id); err != nil {
		return "", err
	}
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, filepath.FromSlash(customProfileDirectory), id+".json"), nil
}

func cleanRepositoryRoot(repoRoot string) (string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return "", errors.New("repository root is required")
	}
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	return filepath.Clean(root), nil
}

func prepareCustomProfileDirectory(repoRoot string) error {
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return err
	}
	profilesRoot := filepath.Join(root, filepath.FromSlash(customProfileDirectory))
	if err := rejectSymlinkComponents(root, profilesRoot); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(profilesRoot, 0o755); err != nil {
		return fmt.Errorf("create custom Baseline Profile directory: %w", err)
	}
	return rejectSymlinkComponents(root, profilesRoot)
}

func rejectSymlinkComponents(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%w: %s", ErrUnsafeCustomProfilePath, target)
	}
	current := root
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		if component == "" || component == "." {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%w: symlink component %s", ErrUnsafeCustomProfilePath, current)
		}
	}
	return nil
}

func cloneJSONMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	data, err := json.Marshal(input)
	if err != nil {
		return map[string]any{}
	}
	var output map[string]any
	if err := json.Unmarshal(data, &output); err != nil {
		return map[string]any{}
	}
	return output
}

func cloneStrings(input []string) []string {
	return append(make([]string, 0, len(input)), input...)
}
