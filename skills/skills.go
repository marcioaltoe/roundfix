package skills

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed roundfix-watch/SKILL.md roundfix-watch/agents/openai.yaml roundfix-resolve-round/SKILL.md roundfix-resolve-round/agents/openai.yaml
var embedded embed.FS

var skillNames = []string{"roundfix-watch", "roundfix-resolve-round"}

type File struct {
	Skill string
	Path  string
	Data  []byte
}

type Diagnostic struct {
	Path    string
	Message string
}

type InstallRequest struct {
	Target     string
	TargetDirs map[string]string
}

type InstalledTarget struct {
	Target string
	Dir    string
	Files  int
}

type InstallResult struct {
	Targets []InstalledTarget
}

func Names() []string {
	return append([]string{}, skillNames...)
}

func Files() ([]File, error) {
	files := make([]File, 0)
	for _, skill := range skillNames {
		err := fs.WalkDir(embedded, skill, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			data, err := embedded.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, File{
				Skill: skill,
				Path:  path,
				Data:  data,
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("read embedded skill %s: %w", skill, err)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	return files, nil
}

func Check() []Diagnostic {
	required := map[string][]string{
		"roundfix-watch/SKILL.md": {
			"Prefer `roundfix` commands over manual GitHub scraping.",
			"Report the Run ID, Open Pull Request, Review Source, Agent, and current Run state",
			"Do not manually resolve CodeRabbit threads",
		},
		"roundfix-resolve-round/SKILL.md": {
			"Read every assigned Review Issue file completely",
			"Update only assigned Review Issue statuses",
			"Do not create commits.",
			"Do not push.",
			"Do not call GitHub, CodeRabbit, or other Review Source mutation APIs.",
			"Do not edit unassigned Review Issue files.",
			"Do not mark any issue as `duplicated`",
		},
		"roundfix-watch/agents/openai.yaml":         {"roundfix-watch"},
		"roundfix-resolve-round/agents/openai.yaml": {"roundfix-resolve-round"},
	}
	banned := []string{
		"reference project",
		"Reference Project",
	}

	var diagnostics []Diagnostic
	for path, phrases := range required {
		data, err := embedded.ReadFile(path)
		if err != nil {
			diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "missing required skill artifact"})
			continue
		}
		text := string(data)
		if !strings.Contains(text, "Roundfix") && strings.HasSuffix(path, "SKILL.md") {
			diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "skill must use Roundfix branding"})
		}
		for _, phrase := range phrases {
			if !strings.Contains(text, phrase) {
				diagnostics = append(diagnostics, Diagnostic{Path: path, Message: fmt.Sprintf("missing required wording %q", phrase)})
			}
		}
		for _, phrase := range banned {
			if strings.Contains(text, phrase) {
				diagnostics = append(diagnostics, Diagnostic{Path: path, Message: fmt.Sprintf("contains banned reference branding %q", phrase)})
			}
		}
	}
	sort.Slice(diagnostics, func(i, j int) bool {
		if diagnostics[i].Path == diagnostics[j].Path {
			return diagnostics[i].Message < diagnostics[j].Message
		}
		return diagnostics[i].Path < diagnostics[j].Path
	})
	return diagnostics
}

func Install(ctx context.Context, req InstallRequest) (InstallResult, error) {
	targets, err := targetDirs(req)
	if err != nil {
		return InstallResult{}, err
	}
	files, err := Files()
	if err != nil {
		return InstallResult{}, err
	}

	result := InstallResult{Targets: make([]InstalledTarget, 0, len(targets))}
	for _, target := range orderedTargets(targets) {
		if err := ctx.Err(); err != nil {
			return InstallResult{}, err
		}
		root := targets[target]
		count := 0
		for _, file := range files {
			dest := filepath.Join(root, file.Path)
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return InstallResult{}, fmt.Errorf("create skill directory %q: %w", filepath.Dir(dest), err)
			}
			if err := os.WriteFile(dest, file.Data, 0o644); err != nil {
				return InstallResult{}, fmt.Errorf("write skill file %q: %w", dest, err)
			}
			count++
		}
		result.Targets = append(result.Targets, InstalledTarget{
			Target: target,
			Dir:    root,
			Files:  count,
		})
	}
	return result, nil
}

func DefaultTargetDirs() (map[string]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory for skill install targets: %w", err)
	}
	return map[string]string{
		"codex":    filepath.Join(firstNonEmpty(os.Getenv("CODEX_HOME"), filepath.Join(home, ".codex")), "skills"),
		"claude":   filepath.Join(firstNonEmpty(os.Getenv("CLAUDE_HOME"), filepath.Join(home, ".claude")), "skills"),
		"opencode": filepath.Join(firstNonEmpty(os.Getenv("OPENCODE_HOME"), filepath.Join(home, ".opencode")), "skills"),
	}, nil
}

func targetDirs(req InstallRequest) (map[string]string, error) {
	target := strings.TrimSpace(req.Target)
	if target == "" {
		target = "all"
	}
	if target != "all" && target != "codex" && target != "claude" && target != "opencode" {
		return nil, fmt.Errorf("unsupported skill install target %q; supported values: codex, claude, opencode, all", target)
	}
	defaults, err := DefaultTargetDirs()
	if err != nil {
		return nil, err
	}
	for key, value := range req.TargetDirs {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if _, ok := defaults[key]; !ok {
			return nil, fmt.Errorf("unsupported skill install target %q", key)
		}
		defaults[key] = value
	}
	if target == "all" {
		return defaults, nil
	}
	dir := defaults[target]
	if dir == "" {
		return nil, errors.New("skill install target directory is empty")
	}
	return map[string]string{target: dir}, nil
}

func orderedTargets(targets map[string]string) []string {
	order := []string{"codex", "claude", "opencode"}
	ordered := make([]string, 0, len(targets))
	for _, target := range order {
		if _, ok := targets[target]; ok {
			ordered = append(ordered, target)
		}
	}
	return ordered
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
