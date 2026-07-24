package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
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
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return CustomProfileInitResult{}, fmt.Errorf("serialize custom Baseline Profile %q: %w", id, err)
	}
	data = append(data, '\n')

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
		Modules:       append([]string(nil), modules...),
		Decisions:     append([]string(nil), decisions...),
		Capabilities:  append([]string(nil), capabilities...),
		Templates:     append([]string(nil), templates...),
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
	return nil
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
		object, ok := objectValue(value)
		if !ok {
			return errors.New("must be an object")
		}
		mode, _ := object["mode"].(string)
		if mode == "" {
			return errors.New("mode is required")
		}
		for _, allowed := range stringsOrEmpty(decision["modes"]) {
			if mode == allowed {
				return nil
			}
		}
		return fmt.Errorf("mode %q is not allowed", mode)
	default:
		return fmt.Errorf("decision type %q is unsupported", kind)
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
		object, ok := objectValue(value)
		if !ok {
			return
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
