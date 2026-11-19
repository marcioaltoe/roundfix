package releaseplan

import (
	"context"
	"path"
	"strings"
)

// ManagedPathSource supplies the adopter's committed maintenance records.
type ManagedPathSource interface {
	ManagedPaths(context.Context, Range) (ManagedPaths, error)
}

// ManagedPaths joins the Setup Manifest and skill lock at the base and target.
type ManagedPaths struct {
	Adopter bool     // the target's skills-lock.json records the roundfix skill
	Files   []string // paths recorded by either Setup Manifest
	Skills  []string // skill names recorded by either lock
}

// Contains reports whether a path is an Adopter Maintenance Path.
func (managed ManagedPaths) Contains(changedPath string) bool {
	clean := cleanChangedPath(changedPath)
	if !managed.Adopter || clean == "" {
		return false
	}
	switch clean {
	case "docs/agents/setup-context.json", "skills-lock.json", "CLAUDE.md",
		"docs/agents/specific-repository.md", ".roundfixrc.yml":
		return true
	}
	if path.Base(clean) == "CONTEXT.md" {
		return true
	}
	if path.Dir(clean) == ".claude/agents" && strings.HasPrefix(path.Base(clean), "roundfix-") && strings.HasSuffix(clean, ".md") {
		return true
	}
	for _, file := range managed.Files {
		if recorded := cleanChangedPath(file); recorded != "" && clean == recorded {
			return true
		}
	}
	for _, name := range managed.Skills {
		if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
			continue
		}
		for _, prefix := range []string{".agents/skills/", ".claude/skills/"} {
			directory := prefix + name
			if clean == directory || strings.HasPrefix(clean, directory+"/") {
				return true
			}
		}
	}
	return false
}
