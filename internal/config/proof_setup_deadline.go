package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// ProfileProofSetupTimeout reads defaults.profile_proof_setup_timeout from
// User Config and Project Config with Load's precedence and range check. It
// reads no other key.
func ProfileProofSetupTimeout(homeDir string, workDir string) (time.Duration, error) {
	homeDir, err := resolveHomeDir(homeDir)
	if err != nil {
		return 0, err
	}
	workDir, err = resolveWorkDir(workDir)
	if err != nil {
		return 0, err
	}
	paths := []string{filepath.Join(homeDir, userConfigRelPath)}
	if root := FindGitRoot(workDir); root != "" {
		paths = append(paths, filepath.Join(root, projectConfigName))
	}
	timeout := DefaultProfileProofSetupTimeout
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return 0, fmt.Errorf("read config %q: %w", path, err)
		}
		var overlay struct {
			Defaults struct {
				ProfileProofSetupTimeout *durationValue `yaml:"profile_proof_setup_timeout"`
			} `yaml:"defaults"`
		}
		if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&overlay); err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			return 0, fmt.Errorf("parse config %q: %w", path, err)
		}
		if value := overlay.Defaults.ProfileProofSetupTimeout; value != nil {
			timeout = value.value
		}
	}
	if timeout <= 0 || timeout > MaxProfileProofSetupTimeout {
		return 0, errors.New("defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m")
	}
	return timeout, nil
}
