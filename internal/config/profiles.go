package config

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type AgentSelection struct {
	Runtime         string `json:"runtime"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoning_effort"`
}

type AgentSelectionProfile struct {
	Preferred AgentSelection
	Fallbacks []AgentSelection
}

type WorkCategory string

const (
	CategoryGeneral  WorkCategory = "general"
	CategoryBackend  WorkCategory = "backend"
	CategoryFrontend WorkCategory = "frontend"
	CategoryData     WorkCategory = "data"
	CategoryInfra    WorkCategory = "infra"
	CategoryDocs     WorkCategory = "docs"
	CategoryTest     WorkCategory = "test"
	CategoryChore    WorkCategory = "chore"
	CategoryQA       WorkCategory = "qa"
	CategoryReview   WorkCategory = "review"
)

type ProfileSource string

const (
	ProfileSourceBuiltIn    ProfileSource = "built-in"
	ProfileSourceUser       ProfileSource = "user"
	ProfileSourceProject    ProfileSource = "project"
	ProfileSourceInvocation ProfileSource = "invocation"
)

type ProfileEntry struct {
	Profile AgentSelectionProfile
	Source  ProfileSource
}

type Profiles map[WorkCategory]ProfileEntry

type ResolvedProfile struct {
	Category      WorkCategory
	Source        ProfileSource
	InheritedFrom WorkCategory
	Profile       AgentSelectionProfile
}

var requiredWorkCategories = []WorkCategory{
	CategoryGeneral,
	CategoryBackend,
	CategoryFrontend,
	CategoryQA,
	CategoryReview,
}

var optionalWorkCategories = []WorkCategory{
	CategoryData,
	CategoryInfra,
	CategoryDocs,
	CategoryTest,
	CategoryChore,
}

var allWorkCategories = []WorkCategory{
	CategoryGeneral,
	CategoryBackend,
	CategoryFrontend,
	CategoryData,
	CategoryInfra,
	CategoryDocs,
	CategoryTest,
	CategoryChore,
	CategoryQA,
	CategoryReview,
}

func ParseWorkCategory(value string) (WorkCategory, bool) {
	category := WorkCategory(strings.TrimSpace(value))
	for _, candidate := range allWorkCategories {
		if category == candidate {
			return category, true
		}
	}
	return "", false
}

func RequiredWorkCategories() []WorkCategory {
	return append([]WorkCategory(nil), requiredWorkCategories...)
}

func OptionalWorkCategories() []WorkCategory {
	return append([]WorkCategory(nil), optionalWorkCategories...)
}

func AllWorkCategories() []WorkCategory {
	return append([]WorkCategory(nil), allWorkCategories...)
}

// ConfiguredWorkCategories returns every required Agent Work Category plus each
// optional category the effective configuration defines. An optional category
// present only by inheritance from general is absent, because it contributes no
// distinct Agent Selection tuple to prove. The required categories keep their
// relative order, so a configuration defining no optional category yields
// exactly RequiredWorkCategories. See ADR-0107.
func ConfiguredWorkCategories(config Config) []WorkCategory {
	entries := config.Profiles
	if entries == nil {
		entries = builtinProfiles()
	}
	categories := make([]WorkCategory, 0, len(allWorkCategories))
	for _, category := range allWorkCategories {
		if isOptionalWorkCategory(category) {
			if _, defined := entries[category]; !defined {
				continue
			}
		}
		categories = append(categories, category)
	}
	return categories
}

func ResolveProfile(config Config, category WorkCategory, preferredOverride *AgentSelection) (ResolvedProfile, error) {
	entries := config.Profiles
	if entries == nil {
		entries = builtinProfiles()
	}
	entry, inheritedFrom, ok := resolveProfileEntry(entries, category)
	if !ok {
		return ResolvedProfile{}, fmt.Errorf("Agent Work Category %q is not configured; supported values: %s", category, supportedWorkCategoryList())
	}

	profile := cloneProfile(entry.Profile)
	source := entry.Source
	if preferredOverride != nil {
		selection, err := normalizeSelection("invocation preferred", *preferredOverride, true)
		if err != nil {
			return ResolvedProfile{}, err
		}
		profile.Preferred = selection
		if err := validateAgentSelectionProfile("invocation profile", profile); err != nil {
			return ResolvedProfile{}, err
		}
		source = ProfileSourceInvocation
	}

	return ResolvedProfile{
		Category:      category,
		Source:        source,
		InheritedFrom: inheritedFrom,
		Profile:       profile,
	}, nil
}

func resolveProfileEntry(entries Profiles, category WorkCategory) (ProfileEntry, WorkCategory, bool) {
	if entry, ok := entries[category]; ok {
		return cloneProfileEntry(entry), "", true
	}
	if isOptionalWorkCategory(category) {
		entry, ok := entries[CategoryGeneral]
		return cloneProfileEntry(entry), CategoryGeneral, ok
	}
	return ProfileEntry{}, "", false
}

func isOptionalWorkCategory(category WorkCategory) bool {
	for _, candidate := range optionalWorkCategories {
		if category == candidate {
			return true
		}
	}
	return false
}

func builtinProfiles() Profiles {
	general := AgentSelectionProfile{
		Preferred: AgentSelection{Runtime: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"},
		Fallbacks: []AgentSelection{{Runtime: "codex", Model: "gpt-5.5", ReasoningEffort: "xhigh"}},
	}
	frontend := AgentSelectionProfile{
		Preferred: AgentSelection{Runtime: "claude", Model: "opus", ReasoningEffort: "xhigh"},
		Fallbacks: []AgentSelection{{Runtime: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"}},
	}
	return Profiles{
		CategoryGeneral:  {Profile: cloneProfile(general), Source: ProfileSourceBuiltIn},
		CategoryBackend:  {Profile: cloneProfile(general), Source: ProfileSourceBuiltIn},
		CategoryFrontend: {Profile: cloneProfile(frontend), Source: ProfileSourceBuiltIn},
		CategoryQA:       {Profile: cloneProfile(general), Source: ProfileSourceBuiltIn},
		CategoryReview:   {Profile: cloneProfile(general), Source: ProfileSourceBuiltIn},
	}
}

type profilesOverlay struct {
	entries Profiles
}

func (overlay *profilesOverlay) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.New("profiles must be a mapping")
	}
	if len(node.Content) == 0 {
		return errors.New("profiles must define at least one Agent Selection Profile")
	}

	entries := Profiles{}
	for index := 0; index < len(node.Content); index += 2 {
		key := strings.TrimSpace(node.Content[index].Value)
		category, ok := ParseWorkCategory(key)
		if !ok {
			return fmt.Errorf("profiles.%s is not a supported Agent Work Category; supported values: %s", key, supportedWorkCategoryList())
		}
		if _, duplicate := entries[category]; duplicate {
			return fmt.Errorf("profiles.%s is defined more than once", category)
		}
		profile, err := decodeProfile("profiles."+string(category), node.Content[index+1])
		if err != nil {
			return err
		}
		entries[category] = ProfileEntry{Profile: profile}
	}
	overlay.entries = entries
	return nil
}

func applyProfilesOverlay(config *Config, overlay *profilesOverlay, source ProfileSource) {
	if overlay == nil {
		return
	}
	if config.Profiles == nil {
		config.Profiles = builtinProfiles()
	}
	for category, entry := range overlay.entries {
		config.Profiles[category] = ProfileEntry{
			Profile: cloneProfile(entry.Profile),
			Source:  source,
		}
	}
}

func applyLegacyRuntimeProfiles(config *Config, source ProfileSource) {
	if config.Profiles == nil {
		config.Profiles = builtinProfiles()
	}
	defaults, _ := config.Runtimes.DefaultsFor(config.Defaults.Agent)
	selection := AgentSelection{
		Runtime:         strings.TrimSpace(config.Defaults.Agent),
		Model:           strings.TrimSpace(defaults.Model),
		ReasoningEffort: strings.TrimSpace(defaults.ReasoningEffort),
	}
	builtins := builtinProfiles()
	for _, category := range requiredWorkCategories {
		base := cloneProfile(builtins[category].Profile)
		base.Preferred = selection
		for index, fallback := range base.Fallbacks {
			if fallback == selection {
				base.Fallbacks[index] = builtins[category].Profile.Preferred
			}
		}
		config.Profiles[category] = ProfileEntry{Profile: base, Source: source}
	}
}

func validateProfiles(entries Profiles) error {
	if entries == nil {
		return errors.New("profiles are required")
	}
	for _, category := range requiredWorkCategories {
		entry, ok := entries[category]
		if !ok {
			return fmt.Errorf("profiles.%s is required", category)
		}
		if err := validateAgentSelectionProfile("profiles."+string(category), entry.Profile); err != nil {
			return err
		}
	}
	for category, entry := range entries {
		if _, ok := ParseWorkCategory(string(category)); !ok {
			return fmt.Errorf("profiles.%s is not a supported Agent Work Category; supported values: %s", category, supportedWorkCategoryList())
		}
		if err := validateAgentSelectionProfile("profiles."+string(category), entry.Profile); err != nil {
			return err
		}
	}
	return nil
}

func decodeProfile(path string, node *yaml.Node) (AgentSelectionProfile, error) {
	if node.Kind != yaml.MappingNode {
		return AgentSelectionProfile{}, fmt.Errorf("%s must be a mapping", path)
	}
	if len(node.Content) == 0 {
		return AgentSelectionProfile{}, fmt.Errorf("%s must define preferred and fallbacks", path)
	}

	var preferredNode *yaml.Node
	var fallbacksNode *yaml.Node
	seen := map[string]bool{}
	for index := 0; index < len(node.Content); index += 2 {
		key := strings.TrimSpace(node.Content[index].Value)
		if seen[key] {
			return AgentSelectionProfile{}, fmt.Errorf("%s.%s is defined more than once", path, key)
		}
		seen[key] = true
		switch key {
		case "preferred":
			preferredNode = node.Content[index+1]
		case "fallbacks":
			fallbacksNode = node.Content[index+1]
		default:
			return AgentSelectionProfile{}, fmt.Errorf("%s.%s is not a supported profile key", path, key)
		}
	}
	if preferredNode == nil {
		return AgentSelectionProfile{}, fmt.Errorf("%s.preferred is required", path)
	}
	if fallbacksNode == nil {
		return AgentSelectionProfile{}, fmt.Errorf("%s.fallbacks is required; one additional distinct authorized and proven Agent Selection is required", path)
	}

	preferred, err := decodeSelection(path+".preferred", preferredNode)
	if err != nil {
		return AgentSelectionProfile{}, err
	}
	fallbacks, err := decodeFallbacks(path+".fallbacks", fallbacksNode)
	if err != nil {
		return AgentSelectionProfile{}, err
	}
	profile := AgentSelectionProfile{Preferred: preferred, Fallbacks: fallbacks}
	if err := validateAgentSelectionProfile(path, profile); err != nil {
		return AgentSelectionProfile{}, err
	}
	return profile, nil
}

func decodeFallbacks(path string, node *yaml.Node) ([]AgentSelection, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s must be a sequence", path)
	}
	if len(node.Content) == 0 {
		return nil, fmt.Errorf("%s must include at least one Agent Selection; one additional distinct authorized and proven Agent Selection is required", path)
	}
	fallbacks := make([]AgentSelection, 0, len(node.Content))
	for index, fallbackNode := range node.Content {
		selection, err := decodeSelection(fmt.Sprintf("%s[%d]", path, index), fallbackNode)
		if err != nil {
			return nil, err
		}
		fallbacks = append(fallbacks, selection)
	}
	return fallbacks, nil
}

func decodeSelection(path string, node *yaml.Node) (AgentSelection, error) {
	if node.Kind != yaml.MappingNode {
		return AgentSelection{}, fmt.Errorf("%s must be a mapping", path)
	}
	values := map[string]string{}
	seen := map[string]bool{}
	for index := 0; index < len(node.Content); index += 2 {
		key := strings.TrimSpace(node.Content[index].Value)
		if seen[key] {
			return AgentSelection{}, fmt.Errorf("%s.%s is defined more than once", path, key)
		}
		seen[key] = true
		switch key {
		case "runtime", "model", "reasoning_effort":
			value, err := decodeStringValue(path+"."+key, node.Content[index+1])
			if err != nil {
				return AgentSelection{}, err
			}
			values[key] = value
		default:
			return AgentSelection{}, fmt.Errorf("%s.%s is not a supported Agent Selection key", path, key)
		}
	}
	for _, key := range []string{"runtime", "model", "reasoning_effort"} {
		if !seen[key] {
			return AgentSelection{}, fmt.Errorf("%s.%s is required", path, key)
		}
	}
	return normalizeSelection(path, AgentSelection{
		Runtime:         values["runtime"],
		Model:           values["model"],
		ReasoningEffort: values["reasoning_effort"],
	}, true)
}

func decodeStringValue(path string, node *yaml.Node) (string, error) {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return "", fmt.Errorf("%s must be a string; use \"\" for a model-managed reasoning effort", path)
	}
	var value string
	if err := node.Decode(&value); err != nil {
		return "", fmt.Errorf("%s must be a string: %w", path, err)
	}
	return value, nil
}

func validateAgentSelectionProfile(path string, profile AgentSelectionProfile) error {
	if _, err := normalizeSelection(path+".preferred", profile.Preferred, true); err != nil {
		return err
	}
	if len(profile.Fallbacks) == 0 {
		return fmt.Errorf("%s.fallbacks must include at least one Agent Selection; one additional distinct authorized and proven Agent Selection is required", path)
	}
	seen := map[AgentSelection]bool{}
	preferred, _ := normalizeSelection(path+".preferred", profile.Preferred, true)
	seen[preferred] = true
	for index, fallback := range profile.Fallbacks {
		selection, err := normalizeSelection(fmt.Sprintf("%s.fallbacks[%d]", path, index), fallback, true)
		if err != nil {
			return err
		}
		if seen[selection] {
			return fmt.Errorf("%s contains duplicate Agent Selection %q; one additional distinct authorized and proven Agent Selection is required", path, formatAgentSelection(selection))
		}
		seen[selection] = true
	}
	return nil
}

func normalizeSelection(path string, selection AgentSelection, requireReasoning bool) (AgentSelection, error) {
	normalized := AgentSelection{
		Runtime:         strings.TrimSpace(selection.Runtime),
		Model:           strings.TrimSpace(selection.Model),
		ReasoningEffort: strings.TrimSpace(selection.ReasoningEffort),
	}
	if normalized.Runtime == "" {
		return AgentSelection{}, fmt.Errorf("%s.runtime must not be empty", path)
	}
	if !isSupportedAgent(normalized.Runtime) {
		return AgentSelection{}, fmt.Errorf("%s.runtime %q is invalid; supported values: codex, claude, opencode", path, normalized.Runtime)
	}
	if normalized.Model == "" {
		return AgentSelection{}, fmt.Errorf("%s.model must not be empty", path)
	}
	if err := validateModelManagedReasoning(path, normalized); err != nil {
		return AgentSelection{}, err
	}
	if !requireReasoning {
		return normalized, nil
	}
	return normalized, nil
}

// runtimesManagingReasoning lists the ACP Runtimes whose reasoning effort only
// the Agent Model can set, so Roundfix must not be configured to assign one.
// See ADR-0106.
var runtimesManagingReasoning = map[string]string{
	"opencode": "OpenCode advertises reasoning effort per model and only after an Agent Session's first prompt, " +
		"so Roundfix cannot apply one during a token-free Exact Agent Selection Proof",
}

func validateModelManagedReasoning(path string, selection AgentSelection) error {
	reason, managed := runtimesManagingReasoning[selection.Runtime]
	if !managed || selection.ReasoningEffort == "" {
		return nil
	}
	return fmt.Errorf(
		"%s.reasoning_effort must be empty for runtime %q; %s; use reasoning_effort: \"\" for model-managed reasoning",
		path,
		selection.Runtime,
		reason,
	)
}

func configHasProfilesSection(document *yaml.Node) bool {
	return documentHasTopLevelKey(document, "profiles")
}

func configHasLegacyRuntimeDefaults(document *yaml.Node) bool {
	root := documentRootMapping(document)
	if root == nil {
		return false
	}
	for index := 0; index < len(root.Content); index += 2 {
		key := root.Content[index].Value
		value := root.Content[index+1]
		switch key {
		case "defaults":
			if mappingHasKey(value, "agent") {
				return true
			}
		case "runtimes":
			return true
		}
	}
	return false
}

func documentHasTopLevelKey(document *yaml.Node, key string) bool {
	root := documentRootMapping(document)
	return mappingHasKey(root, key)
}

func documentRootMapping(document *yaml.Node) *yaml.Node {
	if document == nil {
		return nil
	}
	root := document
	if document.Kind == yaml.DocumentNode {
		if len(document.Content) == 0 {
			return nil
		}
		root = document.Content[0]
	}
	if root.Kind != yaml.MappingNode {
		return nil
	}
	return root
}

func mappingHasKey(node *yaml.Node, key string) bool {
	if node == nil || node.Kind != yaml.MappingNode {
		return false
	}
	for index := 0; index < len(node.Content); index += 2 {
		if node.Content[index].Value == key {
			return true
		}
	}
	return false
}

func profileSchemaConflictError(source ProfileSource) error {
	scope := "user"
	if source == ProfileSourceProject {
		scope = "project"
	}
	return fmt.Errorf("config mixes legacy runtime defaults with profiles; migrate with `roundfix profiles configure --scope %s`", scope)
}

func cloneProfileEntry(entry ProfileEntry) ProfileEntry {
	return ProfileEntry{Profile: cloneProfile(entry.Profile), Source: entry.Source}
}

func cloneProfile(profile AgentSelectionProfile) AgentSelectionProfile {
	return AgentSelectionProfile{
		Preferred: profile.Preferred,
		Fallbacks: append([]AgentSelection(nil), profile.Fallbacks...),
	}
}

func formatAgentSelection(selection AgentSelection) string {
	return selection.Runtime + " / " + selection.Model + " / " + selection.ReasoningEffort
}

func supportedWorkCategoryList() string {
	values := make([]string, 0, len(allWorkCategories))
	for _, category := range allWorkCategories {
		values = append(values, string(category))
	}
	return strings.Join(values, ", ")
}
