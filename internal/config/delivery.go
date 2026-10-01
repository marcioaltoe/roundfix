package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

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
