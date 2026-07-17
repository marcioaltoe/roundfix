package config

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const profileConfigTempPattern = ".roundfix-profiles-*"

type ProfileConfigOptions struct {
	Scope    string
	HomeDir  string
	WorkDir  string
	Profiles Profiles
	DryRun   bool
}

type ProfileConfigResult struct {
	Scope    string
	Path     string
	Changed  bool
	Profiles Profiles
}

func ParseProfilesFragment(content []byte) (Profiles, error) {
	var document yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&document); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("profile fragment must not be empty")
		}
		return nil, fmt.Errorf("parse profile fragment: %w", err)
	}
	if err := requireSingleYAMLDocument(decoder); err != nil {
		return nil, fmt.Errorf("parse profile fragment: %w", err)
	}
	root := documentRootMapping(&document)
	if root == nil {
		return nil, errors.New("profile fragment must be a mapping")
	}

	profilesNode := root
	if wrapped, ok, err := onlyProfilesWrapper(root); err != nil {
		return nil, err
	} else if ok {
		profilesNode = wrapped
	}

	var overlay profilesOverlay
	if err := overlay.UnmarshalYAML(profilesNode); err != nil {
		return nil, err
	}
	return NormalizeProfilesFragment(overlay.entries)
}

func NormalizeProfilesFragment(profiles Profiles) (Profiles, error) {
	if len(profiles) == 0 {
		return nil, errors.New("profiles must define at least one Agent Selection Profile")
	}
	normalized := Profiles{}
	for _, category := range allWorkCategories {
		entry, ok := profiles[category]
		if !ok {
			continue
		}
		profile := cloneProfile(entry.Profile)
		if err := validateAgentSelectionProfile("profiles."+string(category), profile); err != nil {
			return nil, err
		}
		normalized[category] = ProfileEntry{Profile: profile}
	}
	for category := range profiles {
		if _, ok := ParseWorkCategory(string(category)); !ok {
			return nil, fmt.Errorf("profiles.%s is not a supported Agent Work Category; supported values: %s", category, supportedWorkCategoryList())
		}
		if _, ok := normalized[category]; !ok {
			return nil, fmt.Errorf("profiles.%s is not a supported Agent Work Category; supported values: %s", category, supportedWorkCategoryList())
		}
	}
	return normalized, nil
}

func WriteProfilesConfig(ctx context.Context, opts ProfileConfigOptions) (ProfileConfigResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return ProfileConfigResult{}, err
	}
	scope, path, source, err := profileConfigTarget(opts.Scope, opts.HomeDir, opts.WorkDir)
	if err != nil {
		return ProfileConfigResult{}, err
	}
	profiles, err := NormalizeProfilesFragment(opts.Profiles)
	if err != nil {
		return ProfileConfigResult{}, err
	}
	result := ProfileConfigResult{
		Scope:    scope,
		Path:     path,
		Profiles: cloneProfilesForWrite(profiles),
	}

	current, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ProfileConfigResult{}, fmt.Errorf("read config %q: %w", path, err)
	}
	next, err := mergeProfilesConfigContent(path, current, profiles, source)
	if err != nil {
		return ProfileConfigResult{}, err
	}
	if bytes.Equal(current, next) || opts.DryRun {
		return result, nil
	}
	if err := writeFileAtomic(ctx, path, next); err != nil {
		return ProfileConfigResult{}, err
	}
	result.Changed = true
	return result, nil
}

func ProfileConfigTarget(scope string, homeDir string, workDir string) (string, string, error) {
	normalizedScope, path, _, err := profileConfigTarget(scope, homeDir, workDir)
	return normalizedScope, path, err
}

func profileConfigTarget(scope string, homeDir string, workDir string) (string, string, ProfileSource, error) {
	normalizedScope := strings.TrimSpace(scope)
	switch normalizedScope {
	case InitScopeUser:
	case InitScopeProject:
	default:
		return "", "", "", fmt.Errorf("unsupported profiles scope %q; supported values: user, project", scope)
	}
	resolvedHome, err := resolveHomeDir(homeDir)
	if err != nil {
		return "", "", "", err
	}
	resolvedWork, err := resolveWorkDir(workDir)
	if err != nil {
		return "", "", "", err
	}
	if normalizedScope == InitScopeUser {
		return normalizedScope, filepath.Join(resolvedHome, userConfigRelPath), ProfileSourceUser, nil
	}
	gitRoot := findGitRoot(resolvedWork)
	if gitRoot == "" {
		return "", "", "", errors.New("project profiles configure requires a Git root; use --scope user outside a repository")
	}
	return normalizedScope, filepath.Join(gitRoot, projectConfigName), ProfileSourceProject, nil
}

func mergeProfilesConfigContent(path string, content []byte, profiles Profiles, source ProfileSource) ([]byte, error) {
	document, root, err := decodeConfigDocumentForProfiles(path, content)
	if err != nil {
		return nil, err
	}
	hasProfiles := configHasProfilesSection(document)
	hasLegacy := configHasLegacyRuntimeDefaults(document)
	if hasProfiles && hasLegacy {
		return nil, fmt.Errorf("parse config %q: %w", path, profileSchemaConflictError(source))
	}
	if hasLegacy {
		return nil, fmt.Errorf("parse config %q: legacy runtime defaults cannot be combined with profiles; migrate with `roundfix profiles configure --scope %s` after removing defaults.agent and runtimes", path, scopeForProfileSource(source))
	}
	setTopLevelYAMLNode(root, "profiles", profilesYAMLNode(profiles))
	encoded, err := encodeYAMLNode(document)
	if err != nil {
		return nil, fmt.Errorf("encode config %q: %w", path, err)
	}
	return encoded, nil
}

func decodeConfigDocumentForProfiles(path string, content []byte) (*yaml.Node, *yaml.Node, error) {
	document := &yaml.Node{}
	if len(bytes.TrimSpace(content)) == 0 {
		document.Kind = yaml.DocumentNode
		document.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		return document, document.Content[0], nil
	}
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(document); err != nil {
		return nil, nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if err := requireSingleYAMLDocument(decoder); err != nil {
		return nil, nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	root := documentRootMapping(document)
	if root == nil {
		return nil, nil, fmt.Errorf("parse config %q: config must be a mapping", path)
	}
	return document, root, nil
}

func requireSingleYAMLDocument(decoder *yaml.Decoder) error {
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return errors.New("multiple YAML documents are not supported")
}

func onlyProfilesWrapper(root *yaml.Node) (*yaml.Node, bool, error) {
	var profilesNode *yaml.Node
	for index := 0; index+1 < len(root.Content); index += 2 {
		key := strings.TrimSpace(root.Content[index].Value)
		if key != "profiles" {
			continue
		}
		if profilesNode != nil {
			return nil, false, errors.New("profile fragment profiles is defined more than once")
		}
		profilesNode = root.Content[index+1]
	}
	if profilesNode == nil {
		return nil, false, nil
	}
	if len(root.Content) != 2 {
		return nil, false, errors.New("profile fragment with profiles wrapper must not include other top-level keys")
	}
	return profilesNode, true, nil
}

func setTopLevelYAMLNode(root *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index+1 < len(root.Content); index += 2 {
		if root.Content[index].Value == key {
			root.Content[index+1] = value
			return
		}
	}
	root.Content = append(root.Content, yamlStringNode(key), value)
}

func profilesYAMLNode(profiles Profiles) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, category := range allWorkCategories {
		entry, ok := profiles[category]
		if !ok {
			continue
		}
		node.Content = append(node.Content, yamlStringNode(string(category)), profileYAMLNode(entry.Profile))
	}
	return node
}

func profileYAMLNode(profile AgentSelectionProfile) *yaml.Node {
	fallbacks := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, selection := range profile.Fallbacks {
		fallbacks.Content = append(fallbacks.Content, selectionYAMLNode(selection))
	}
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			yamlStringNode("preferred"), selectionYAMLNode(profile.Preferred),
			yamlStringNode("fallbacks"), fallbacks,
		},
	}
}

func selectionYAMLNode(selection AgentSelection) *yaml.Node {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			yamlStringNode("runtime"), yamlStringNode(selection.Runtime),
			yamlStringNode("model"), yamlStringNode(selection.Model),
			yamlStringNode("reasoning_effort"), yamlStringNode(selection.ReasoningEffort),
		},
	}
}

func yamlStringNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func writeFileAtomic(ctx context.Context, path string, content []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config directory %q: %w", dir, err)
	}
	temp, err := os.CreateTemp(dir, profileConfigTempPattern)
	if err != nil {
		return fmt.Errorf("create temporary config in %q: %w", dir, err)
	}
	tempPath := temp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tempPath)
		}
	}()
	if _, err := temp.Write(content); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary config %q: %w", tempPath, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("write temporary config %q: %w", tempPath, err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace config %q atomically: %w", path, err)
	}
	cleanup = false
	return nil
}

func cloneProfilesForWrite(profiles Profiles) Profiles {
	cloned := Profiles{}
	for category, entry := range profiles {
		cloned[category] = ProfileEntry{Profile: cloneProfile(entry.Profile), Source: entry.Source}
	}
	return cloned
}

func scopeForProfileSource(source ProfileSource) string {
	if source == ProfileSourceProject {
		return InitScopeProject
	}
	return InitScopeUser
}
