package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
)

// Archived readers may run after the last active Spec left the checkout.
// Existing directories keep the normal configuration validation.
func resolveHistoricalSpecsRoot(loaded roundconfig.Loaded, repo string) (roundconfig.SpecsRoot, error) {
	configured := strings.TrimSpace(loaded.Config.Specs.Root)
	if configured == "" {
		return roundconfig.ResolveSpecsRoot(loaded, repo)
	}
	root := filepath.Clean(configured)
	if !filepath.IsAbs(root) {
		root = filepath.Join(repo, root)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		return roundconfig.ResolveSpecsRoot(loaded, repo)
	}
	rel, err := filepath.Rel(repo, root)
	if err != nil {
		return roundconfig.SpecsRoot{}, fmt.Errorf("locate historical Specs Root: %w", err)
	}
	return roundconfig.SpecsRoot{
		Path:        root,
		BuiltInRoot: root == filepath.Join(repo, "docs", "specs"),
		External:    rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)),
	}, nil
}

func readArchivedSpecAt(ctx context.Context, runner preflight.GitRunner, repo, ref, root, slug string) (spec.ArchivedSpec, error) {
	return spec.ReadArchivedSpecAt(root, slug, func(file string) ([]byte, error) {
		files, err := runner.RunGit(ctx, repo, "ls-tree", "--name-only", ref, "--", file)
		if err != nil {
			return nil, fmt.Errorf("inspect archived path %s: %w", file, err)
		}
		if strings.TrimSpace(files) == "" {
			return nil, os.ErrNotExist
		}
		body, err := runner.RunGit(ctx, repo, "show", ref+":"+file)
		return []byte(body), err
	})
}
