// Package suiteguardcontract holds repository contracts shared by the suite
// guard and the audits that verify its installation and regeneration behavior.
package suiteguardcontract

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const sanctionedRegenerationHeading = "Sanctioned regeneration"

// SanctionedRegeneration binds one declared command to any outputs the record
// still enumerates. An empty Outputs slice leaves output resolution to the
// command's repository-owned declaration.
type SanctionedRegeneration struct {
	Command string
	Outputs []string
}

// ReadSanctionedRegenerations reads every authorization record under root.
func ReadSanctionedRegenerations(root string) ([]SanctionedRegeneration, error) {
	authorizationRoot := filepath.Join(root, "docs", "workflow", "authorizations")
	var declarations []SanctionedRegeneration
	err := filepath.WalkDir(authorizationRoot, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("inspect authorization record %q: %w", filePath, walkErr)
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read authorization record %q: %w", filePath, err)
		}
		declarations = append(declarations, ParseSanctionedRegenerations(content)...)
		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("walk authorization records %q: %w", authorizationRoot, err)
	}
	sort.Slice(declarations, func(i, j int) bool {
		return declarations[i].Command < declarations[j].Command
	})
	return declarations, nil
}

// ParseSanctionedRegenerations reads the "Sanctioned regeneration" YAML
// blocks used by both the suite guard and the changed-path audit.
func ParseSanctionedRegenerations(content []byte) []SanctionedRegeneration {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	inSanctionedRegeneration := false
	var declarations []SanctionedRegeneration
	for index := 0; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "## ") {
			inSanctionedRegeneration = strings.EqualFold(
				strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")),
				sanctionedRegenerationHeading,
			)
			continue
		}
		if !inSanctionedRegeneration || trimmed != "```yaml" {
			continue
		}

		start := index + 1
		for index++; index < len(lines) && strings.TrimSpace(lines[index]) != "```"; index++ {
		}
		if index >= len(lines) {
			return declarations
		}
		var declaration struct {
			Command string    `yaml:"command"`
			Outputs *[]string `yaml:"outputs"`
		}
		if err := yaml.Unmarshal([]byte(strings.Join(lines[start:index], "\n")), &declaration); err != nil {
			continue
		}
		declaration.Command = strings.TrimSpace(declaration.Command)
		if declaration.Command == "" {
			continue
		}

		var outputs []string
		if declaration.Outputs != nil {
			outputs = make([]string, 0, len(*declaration.Outputs))
			seen := make(map[string]bool, len(*declaration.Outputs))
			for _, output := range *declaration.Outputs {
				output = strings.TrimSpace(output)
				clean := cleanRepositoryPath(output)
				if clean == "" || clean != output || strings.ContainsAny(clean, "*?") || seen[clean] {
					continue
				}
				seen[clean] = true
				outputs = append(outputs, clean)
			}
			if len(outputs) == 0 {
				continue
			}
			sort.Strings(outputs)
		}
		declarations = append(declarations, SanctionedRegeneration{
			Command: declaration.Command,
			Outputs: outputs,
		})
	}
	return declarations
}

func cleanRepositoryPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, `\`) {
		return ""
	}
	clean := filepath.ToSlash(filepath.Clean(value))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return ""
	}
	return clean
}
