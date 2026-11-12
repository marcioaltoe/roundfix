package cli

import (
	"fmt"
	"path/filepath"

	"roundfix/internal/app"
	"roundfix/internal/config"
)

func refuseNewerBaselineRepository(repo string) error {
	absolute, err := filepath.Abs(repo)
	if err != nil {
		return fmt.Errorf("resolve Baseline repository: %w", err)
	}
	root := config.FindGitRoot(absolute)
	if root == "" {
		return nil
	}
	return config.CheckRecordedRoundfixVersion(root, app.Version)
}
