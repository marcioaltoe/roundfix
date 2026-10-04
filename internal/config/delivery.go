package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// Declared reports whether either field of the item build is present.
func (declaration ItemBinaryDeclaration) Declared() bool {
	return declaration.Build != "" || declaration.Path != ""
}

// UnmarshalYAML preserves strict nested-key validation for the declaration.
func (declaration *ItemBinaryDeclaration) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.MappingNode {
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index].Value
			switch key {
			case "build", "path":
			default:
				return fmt.Errorf("delivery.item_binary.%s is not a supported config key", key)
			}
		}
	}
	type rawItemBinaryDeclaration ItemBinaryDeclaration
	var raw rawItemBinaryDeclaration
	if err := node.Decode(&raw); err != nil {
		return err
	}
	*declaration = ItemBinaryDeclaration(raw)
	if !declaration.Declared() {
		return errors.New("delivery.item_binary requires build and path")
	}
	return nil
}

func validateItemBinary(declaration ItemBinaryDeclaration) error {
	if !declaration.Declared() {
		return nil
	}
	if strings.TrimSpace(declaration.Build) == "" || strings.TrimSpace(declaration.Path) == "" {
		return errors.New("delivery.item_binary requires build and path")
	}
	if path.IsAbs(declaration.Path) || strings.Contains(declaration.Path, "\\") || declaration.Path == "." || path.Clean(declaration.Path) != declaration.Path {
		return fmt.Errorf("delivery.item_binary has unsafe path %q", declaration.Path)
	}
	for _, segment := range strings.Split(declaration.Path, "/") {
		if segment == ".." {
			return fmt.Errorf("delivery.item_binary has unsafe path %q", declaration.Path)
		}
	}
	return nil
}

func validateDerivedPaths(declarations []DerivedPathDeclaration) error {
	for index, declaration := range declarations {
		if len(declaration.Paths) == 0 || strings.TrimSpace(declaration.Regenerate) == "" {
			return fmt.Errorf("delivery.derived_paths[%d] requires paths and regenerate", index)
		}
		for _, entry := range declaration.Paths {
			clean := strings.TrimSuffix(entry, "/")
			if clean == "" || clean == "." || path.IsAbs(entry) || strings.Contains(entry, "\\") || path.Clean(clean) != clean {
				return fmt.Errorf("delivery.derived_paths[%d] has unsafe path %q", index, entry)
			}
			for _, segment := range strings.Split(entry, "/") {
				if segment == ".." {
					return fmt.Errorf("delivery.derived_paths[%d] has unsafe path %q", index, entry)
				}
			}
			if _, err := path.Match(entry, ""); err != nil {
				return fmt.Errorf("delivery.derived_paths[%d] has invalid pattern %q: %w", index, entry, err)
			}
		}
	}
	return nil
}

// Matches covers an exact file, a directory prefix, or a whole-path Go pattern.
func (declaration DerivedPathDeclaration) Matches(name string) bool {
	for _, entry := range declaration.Paths {
		if strings.HasSuffix(entry, "/") && strings.HasPrefix(name, entry) {
			return true
		}
		if matched, _ := path.Match(entry, name); matched {
			return true
		}
	}
	return false
}

// DeliveryConfigAtCommit reads the project scope from an immutable default
// commit, retaining the user scope beneath it. It never reads the item tree.
func DeliveryConfigAtCommit(ctx context.Context, runner interface {
	RunGit(context.Context, string, ...string) (string, error)
}, workDir, userConfigPath, commit string) (Config, error) {
	listing, err := runner.RunGit(ctx, workDir, "ls-tree", "--name-only", commit, "--", ".roundfixrc.yml")
	if err != nil {
		return Config{}, fmt.Errorf("inspect default Project Config: %w", err)
	}
	var project []byte
	if strings.TrimSpace(listing) != "" {
		content, err := runner.RunGit(ctx, workDir, "show", commit+":.roundfixrc.yml")
		if err != nil {
			return Config{}, fmt.Errorf("read default Project Config: %w", err)
		}
		project = []byte(content)
	}
	user, err := os.ReadFile(userConfigPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("read User Config for conflict: %w", err)
	}
	config, err := ResolveConfigProposal(user, project)
	if err != nil {
		return Config{}, fmt.Errorf("read conflict declarations: %w", err)
	}
	return config, nil
}
