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

	"gopkg.in/yaml.v3"
)

//go:embed roundfix write-idea write-prd write-techspec write-tasks setup-workflow implement-task implement-spec brainstorming council business-analyst archive-spec qa-gate evidence-gate recommended.txt
var embedded embed.FS

// skillNames is the Roundfix-owned skill bundle shipped in the binary: the
// operational roundfix skill plus the 13 authorial workflow skills. Kept in
// sync with the Makefile OWNED_SKILLS list by make skills-sync/skills-sync-check.
var skillNames = []string{
	"roundfix",
	"write-idea", "write-prd", "write-techspec", "write-tasks",
	"setup-workflow", "implement-task", "implement-spec",
	"brainstorming", "council", "business-analyst",
	"archive-spec", "qa-gate", "evidence-gate",
}

type File struct {
	Skill string
	Path  string
	Data  []byte
}

type Diagnostic struct {
	Path    string
	Message string
}

type openAIManifest struct {
	Name         string `yaml:"name"`
	EntryPoint   string `yaml:"entrypoint"`
	RuntimeHints struct {
		Command string `yaml:"command"`
	} `yaml:"runtime_hints"`
	Runtime struct {
		Hints struct {
			Command string `yaml:"command"`
		} `yaml:"hints"`
	} `yaml:"runtime"`
}

type InstallRequest struct {
	Target     string
	TargetDirs map[string]string
	ProjectDir string
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
	banned := []string{
		"reference project",
		"Reference Project",
	}
	var diagnostics []Diagnostic

	// The operational roundfix skill carries a strict contract: required
	// wording, Roundfix branding, and a valid OpenAI manifest.
	roundfixRequired := map[string][]string{
		"roundfix/SKILL.md": {
			"Prefer `roundfix` commands over manual GitHub scraping.",
			"Report the Run ID",
			"state whenever you summarize progress.",
			"Review Runs (`fetch`, `resolve`, and `watch`) execute in the user's checkout",
			"Branch Integrity Preflight runs before any fetch, Agent Session",
			"`--skip-branch-integrity` is the only bypass",
			"watch ends CleanUnverified, exits `3`",
			"Roundfix publishes Outcome Comments",
			"adapter: ok",
			"commit <path>",
			"`  reason: <one line>`",
			"owner PID is provably dead",
			"`done` becomes `completed`",
			"This Verification Feedback retry never consumes a Round",
			"Do not manually resolve CodeRabbit threads",
			"Read every assigned Review Issue file completely",
			"Update only assigned Review Issue statuses",
			"Do not create commits inside an assigned Batch run.",
			"Do not push inside an assigned Batch run.",
			"Do not call GitHub, CodeRabbit, or other Review Source mutation APIs",
			"Do not edit unassigned Review Issue files.",
			"Do not mark any issue as `duplicated`",
			"rtk bun run --cwd <package-dir> <script> [args...]",
		},
		"roundfix/agents/openai.yaml": {
			"name: roundfix",
			"entrypoint: SKILL.md",
			"command: roundfix watch --source coderabbit --pr <number> --agent <agent> --until-clean",
			"review_run_contract:",
			"watch_outcome_contract:",
			"CleanUnverified with exit code",
			"review_source_contract:",
			"assigned Review Issue files during Batch runs",
			"Run state",
		},
	}
	for path, phrases := range roundfixRequired {
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
		if path == "roundfix/agents/openai.yaml" {
			diagnostics = append(diagnostics, checkOpenAIManifest(path, data)...)
		}
	}

	// The authorial workflow skills get structural validation only: a SKILL.md
	// that parses with a name and carries no reference-project branding. Their
	// generic authorial language and independent versions never trip the check.
	for _, skill := range skillNames {
		if skill == "roundfix" {
			continue
		}
		path := skill + "/SKILL.md"
		data, err := embedded.ReadFile(path)
		if err != nil {
			diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "missing SKILL.md"})
			continue
		}
		text := string(data)
		switch name := frontmatterName(text); {
		case name == "":
			diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "SKILL.md frontmatter is missing a name"})
		case name != skill:
			diagnostics = append(diagnostics, Diagnostic{Path: path, Message: fmt.Sprintf("frontmatter name %q does not match skill directory %q", name, skill)})
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

// frontmatterName returns the name field from a SKILL.md YAML frontmatter block,
// or "" when there is none.
func frontmatterName(text string) string {
	trimmed := strings.TrimSpace(text)
	if !strings.HasPrefix(trimmed, "---") {
		return ""
	}
	rest := strings.TrimPrefix(trimmed, "---")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	var meta struct {
		Name string `yaml:"name"`
	}
	if err := yaml.Unmarshal([]byte(rest[:end]), &meta); err != nil {
		return ""
	}
	return strings.TrimSpace(meta.Name)
}

// Recommended returns the externally-managed skills (from skills-lock.json) that
// Roundfix recommends but does not ship. They install through the user's own
// skills tooling, never through roundfix.
func Recommended() []string {
	data, err := embedded.ReadFile("recommended.txt")
	if err != nil {
		return nil
	}
	var names []string
	for _, line := range strings.Split(string(data), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			names = append(names, s)
		}
	}
	return names
}

func checkOpenAIManifest(path string, data []byte) []Diagnostic {
	var manifest openAIManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return []Diagnostic{{Path: path, Message: fmt.Sprintf("parse skill manifest: %v", err)}}
	}
	var diagnostics []Diagnostic
	if strings.TrimSpace(manifest.Name) != "roundfix" {
		diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "manifest field name must be roundfix"})
	}
	if strings.TrimSpace(manifest.EntryPoint) == "" {
		diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "manifest field entrypoint is required"})
	}
	if strings.TrimSpace(openAIManifestCommand(manifest)) == "" {
		diagnostics = append(diagnostics, Diagnostic{Path: path, Message: "manifest runtime command is required"})
	}
	return diagnostics
}

func openAIManifestCommand(manifest openAIManifest) string {
	if command := strings.TrimSpace(manifest.RuntimeHints.Command); command != "" {
		return command
	}
	return strings.TrimSpace(manifest.Runtime.Hints.Command)
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
		target = "project"
	}
	if target != "project" && target != "all" && target != "codex" && target != "claude" && target != "opencode" {
		return nil, fmt.Errorf("unsupported skill install target %q; supported values: project, codex, claude, opencode, all", target)
	}
	if target == "project" {
		return projectTargetDirs(req)
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

func projectTargetDirs(req InstallRequest) (map[string]string, error) {
	if dir := strings.TrimSpace(req.TargetDirs["project"]); dir != "" {
		return map[string]string{"project": dir}, nil
	}
	projectDir := strings.TrimSpace(req.ProjectDir)
	if projectDir == "" {
		return nil, errors.New("project skill install requires a project directory")
	}
	return map[string]string{"project": filepath.Join(projectDir, ".agents", "skills")}, nil
}

func orderedTargets(targets map[string]string) []string {
	order := []string{"project", "codex", "claude", "opencode"}
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
