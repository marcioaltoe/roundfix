package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/codex"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	roundnotify "roundfix/internal/notify"
	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	"roundfix/internal/watch"
	runworktree "roundfix/internal/worktree"
)

func init() {
	newOutcomeNotifier = func(roundconfig.Config) roundnotify.Notifier {
		return testNoopOutcomeNotifier{}
	}
}

type testNoopOutcomeNotifier struct{}

func (testNoopOutcomeNotifier) Notify(context.Context, roundnotify.Outcome) error {
	return nil
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix watch") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix runs\n") {
		t.Fatalf("expected help output to list bare runs command, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix runs list") {
		t.Fatalf("expected help output to list runs command, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	for _, args := range [][]string{{"--version"}, {"version"}, {"-v"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d", code)
			}
			if !strings.HasPrefix(stdout.String(), "roundfix ") {
				t.Fatalf("expected version output, got %q", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

func TestRunInitCreatesProjectConfigWithExplicitScope(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected Project Config %s: %v", configPath, err)
	}
	if !strings.Contains(stdout.String(), "Roundfix config created") || !strings.Contains(stdout.String(), "roundfix fetch --pr <number>") {
		t.Fatalf("expected init success output, got %q", stdout.String())
	}
}

func TestRunInitCreatesUserConfigWithExplicitScope(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "user"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected User Config %s: %v", configPath, err)
	}
	if !strings.Contains(stdout.String(), configPath) {
		t.Fatalf("expected user config path in stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunInitPromptsForScopeAndDefaultsProject(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	prompted := false
	withInitScopePrompt(t, func(context.Context, io.Writer) (string, error) {
		prompted = true
		return roundconfig.InitScopeProject, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !prompted {
		t.Fatal("expected init to prompt for missing scope")
	}
	if _, err := os.Stat(filepath.Join(repoDir, ".roundfixrc.yml")); err != nil {
		t.Fatalf("expected default Project Config: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr from stubbed prompt, got %q", stderr.String())
	}
}

func TestReadInitScopeDefaultsProjectOnBlankInput(t *testing.T) {
	var stderr bytes.Buffer

	scope, err := readInitScope(context.Background(), strings.NewReader("\n"), &stderr)

	if err != nil {
		t.Fatalf("expected blank scope to default, got %v", err)
	}
	if scope != roundconfig.InitScopeProject {
		t.Fatalf("expected project default, got %q", scope)
	}
	if !strings.Contains(stderr.String(), "Scope [project]") {
		t.Fatalf("expected prompt output, got %q", stderr.String())
	}
}

func TestRunInitRejectsExistingConfigWithoutForce(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	mustWrite(t, configPath, "defaults:\n  agent: claude\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--force") {
		t.Fatalf("expected force guidance, got %q", stderr.String())
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "agent: claude") {
		t.Fatalf("expected existing config to remain, got %s", string(content))
	}
}

func TestRunInitForceOverwritesExistingConfig(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	mustWrite(t, configPath, "defaults:\n  agent: claude\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project", "--force"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "profiles:") ||
		!strings.Contains(string(content), "model: gpt-5.6-sol") ||
		!strings.Contains(string(content), "model: claude-opus-5") ||
		strings.Contains(string(content), "agent: claude") ||
		strings.Contains(string(content), "runtimes:") {
		t.Fatalf("expected generated config to replace old content, got %s", string(content))
	}
	if !strings.Contains(stdout.String(), "Roundfix config updated") {
		t.Fatalf("expected updated output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunCommandHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "fetch",
			args:     []string{"fetch", "--help"},
			contains: []string{"roundfix fetch --source coderabbit --pr <number>", "Branch Integrity Preflight", "creates no Run Worktree"},
		},
		{
			name:     "resolve",
			args:     []string{"resolve", "--help"},
			contains: []string{"roundfix resolve --pr <number> [--spec <slug>]", "--agent <agent> --model <model> --reasoning-effort <effort>", "use the review profile", "Branch Integrity Preflight", "clean tracked checkout", "--reasoning-effort", "--no-agent-console", "--detach"},
		},
		{
			name:     "watch",
			args:     []string{"watch", "--help"},
			contains: []string{"roundfix watch --source coderabbit --pr <number> [--spec <slug>]", "--agent <agent> --model <model> --reasoning-effort <effort>", "use the review profile", "Branch Integrity Preflight", "CleanUnverified", "exits 3", "--reasoning-effort", "--until-clean", "Review Source check", "--no-agent-console", "--detach"},
		},
		{
			name:     "setup",
			args:     []string{"setup", "--help"},
			contains: []string{"roundfix setup [--yes] [--no-input]", "--yes", "--no-input"},
		},
		{
			name:     "doctor",
			args:     []string{"doctor", "--help"},
			contains: []string{"roundfix doctor", "the effective adapter, the required\nAgent Selection Profiles", "Repository Skill Set", "skills:", "codex runtime hygiene", "offline", "read-only", "mutates nothing"},
		},
		{
			name:     "upgrade",
			args:     []string{"upgrade", "--help"},
			contains: []string{"roundfix upgrade [--check]", "--check", "atomically"},
		},
		{
			name:     "runs",
			args:     []string{"runs", "--help"},
			contains: []string{"roundfix runs list [--all] [--state <active|terminal|all>] [--limit N]", "--all", "--state", "--limit"},
		},
		{
			name:     "profiles",
			args:     []string{"profiles", "--help"},
			contains: []string{"roundfix profiles show [--category <category>] [--json]", "show"},
		},
		{
			name:     "profiles show",
			args:     []string{"profiles", "show", "--help"},
			contains: []string{"roundfix profiles show [--category <category>] [--json]", "--category", "--json", "read-only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d", code)
			}
			for _, want := range tt.contains {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("expected help output to contain %q, got %q", want, stdout.String())
				}
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

func TestCommandUsageDocumentsProfileLedAndCompleteSelectionOverrides(t *testing.T) {
	for _, name := range []string{"resolve", "watch", "implement"} {
		t.Run(name, func(t *testing.T) {
			got := commandUsage(name)
			for _, want := range []string{
				"Omit all Agent Selection flags",
				"--agent <agent> --model <model> --reasoning-effort <effort>",
				"requires --agent, --model, and --reasoning-effort together",
			} {
				if !strings.Contains(got, want) {
					t.Fatalf("%s help missing %q:\n%s", name, want, got)
				}
			}
		})
	}

	fetchHelp := commandUsage("fetch")
	for _, forbidden := range []string{"--agent", "--model", "--reasoning-effort", "Agent Selection flags"} {
		if strings.Contains(fetchHelp, forbidden) {
			t.Fatalf("fetch help changed to include Agent selection term %q:\n%s", forbidden, fetchHelp)
		}
	}
}

func TestEventsHelpDocumentsAgentSelectionFilter(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"events", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("events help exit = %d, want 0 stderr=%q", code, stderr.String())
	}
	if stderr.String() != "" {
		t.Fatalf("events help stderr = %q, want empty", stderr.String())
	}
	wantCategories := "task-status,batch,verification,outcome,agent-selection"
	if !strings.Contains(stdout.String(), wantCategories) {
		t.Fatalf("events help missing filter categories %q, got %q", wantCategories, stdout.String())
	}

	req, err := parseEventsCommand([]string{"run_123", "--filter", "agent-selection"})
	if err != nil {
		t.Fatalf("parse agent-selection filter: %v", err)
	}
	if !req.filter.Includes(runevent.StreamCategorySelection) {
		t.Fatalf("parsed filter excludes %q: %#v", runevent.StreamCategorySelection, req.filter)
	}

	guide := mustRead(t, filepath.Join(cliTestRepoRoot(t), "docs", "user-guide", "commands.md"))
	if !strings.Contains(guide, wantCategories) || !strings.Contains(guide, "`agent-selection`") {
		t.Fatalf("user guide does not document agent-selection filter")
	}
}

func TestProfilesDocumentationContractMatchesPublicGuidance(t *testing.T) {
	repoRoot := cliTestRepoRoot(t)
	readme := mustRead(t, filepath.Join(repoRoot, "README.md"))
	commands := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "commands.md"))
	usage := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "usage.md"))
	configuration := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "configuration.md"))
	autonomousWork := mustRead(t, filepath.Join(repoRoot, "docs", "agents", "autonomous-work.md"))
	releaseRunbook := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "release-runbook.md"))
	roundfixSkill := mustRead(t, filepath.Join(repoRoot, ".agents", "skills", "roundfix", "SKILL.md"))
	roundfixManifest := mustRead(t, filepath.Join(repoRoot, ".agents", "skills", "roundfix", "agents", "openai.yaml"))

	for _, doc := range []struct {
		name    string
		content string
	}{
		{name: "usage", content: usage},
		{name: "configuration", content: configuration},
		{name: "roundfix skill", content: roundfixSkill},
	} {
		for _, want := range []string{
			"roundfix profiles show",
			"roundfix profiles configure",
			"roundfix profiles validate",
			"gpt-5.6-sol",
			"gpt-5.6-terra",
			"claude-opus-5",
			"claude-fable-5",
			"2026-07-16",
			"category_specific: false",
			"agent_work_started",
			"defaults.agent",
			"runtimes",
			"gpt-5.5",
			"xhigh",
		} {
			if !strings.Contains(doc.content, want) {
				t.Fatalf("%s documentation is missing %q", doc.name, want)
			}
		}
	}

	for _, doc := range []struct {
		name    string
		content string
	}{
		{name: "configuration", content: configuration},
		{name: "usage", content: usage},
		{name: "roundfix skill", content: roundfixSkill},
	} {
		for _, want := range []string{
			"@agentclientprotocol/codex-acp",
			"1.1.4",
			"@zed-industries/codex-acp",
			"exact proof",
			"advisory",
			"gpt-5.6-terra",
			"gpt-5.6-luna",
		} {
			if !strings.Contains(doc.content, want) {
				t.Fatalf("%s readiness documentation is missing %q", doc.name, want)
			}
		}
		if strings.Contains(doc.content, "fallback `codex / gpt-5.6-terra / max`") {
			t.Fatalf("%s still documents Terra/max as the generated fallback", doc.name)
		}
	}

	for _, want := range []string{
		"profiles: ok",
		"roundfix profiles configure",
		"roundfix profiles validate",
		"does not recommend model-managed reasoning",
	} {
		if !strings.Contains(commands, want) {
			t.Fatalf("command documentation is missing %q", want)
		}
	}

	for _, want := range []string{
		"partial Agent Selection override",
		"omit all three selection flags",
		"provide `--agent`, `--model`, and `--reasoning-effort` together",
		"exit `2` before configuration, proof, or Run mutation",
	} {
		if !strings.Contains(releaseRunbook, want) {
			t.Fatalf("release guidance is missing %q", want)
		}
	}

	for _, want := range []string{
		"profiles.general",
		"gpt-5.6-sol",
		"gpt-5.5",
		"Fallback Chain",
	} {
		if !strings.Contains(autonomousWork, want) {
			t.Fatalf("autonomous work guidance is missing %q", want)
		}
	}

	for _, help := range []struct {
		name     string
		snippets []string
	}{
		{name: "setup", snippets: []string{"official Codex adapter", "profile readiness", "before writing"}},
		{name: "doctor", snippets: []string{"Agent Selection Profiles", "profiles:", "Repository Skill Set", "skills:", "read-only", "mutates nothing"}},
		{name: "profiles configure", snippets: []string{"exact Agent Selection", "before confirmation", "without writing"}},
		{name: "profiles validate", snippets: []string{"exact", "Read-only", "disposable ACP Runtime session"}},
	} {
		content := commandUsage(help.name)
		for _, want := range help.snippets {
			if !strings.Contains(content, want) {
				t.Fatalf("%s help is missing %q:\n%s", help.name, want, content)
			}
		}
	}

	for _, doc := range []struct {
		name    string
		content string
	}{
		{name: "README", content: readme},
		{name: "command reference", content: commands},
		{name: "usage", content: usage},
		{name: "autonomous work", content: autonomousWork},
		{name: "roundfix skill", content: roundfixSkill},
		{name: "roundfix manifest", content: roundfixManifest},
	} {
		assertAgentStartingExamplesUseProfilesOrCompleteOverrides(t, doc.name, doc.content)
	}

	for _, want := range []string{
		"roundfix/profiles/v1",
		"roundfix/profiles-configure/v1",
		"roundfix/profiles-validate/v1",
		"notification-first",
		"only then activates",
		"no fallback",
	} {
		if !strings.Contains(roundfixSkill, want) {
			t.Fatalf("roundfix skill is missing %q", want)
		}
	}

	for _, path := range []string{
		filepath.Join(repoRoot, ".agents", "skills", "write-tasks", "SKILL.md"),
		filepath.Join(repoRoot, "skills", "write-tasks", "SKILL.md"),
	} {
		content := mustRead(t, path)
		for _, forbidden := range []string{
			"gpt-5.6",
			"claude-fable",
			"roundfix profiles",
			"profiles:",
			"recommendation",
			"ranking",
			"runtime:",
			"model:",
		} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s must not contain profile policy term %q", path, forbidden)
			}
		}
	}
}

func assertAgentStartingExamplesUseProfilesOrCompleteOverrides(t *testing.T, label string, content string) {
	t.Helper()
	content = strings.ReplaceAll(content, "\\\n", " ")
	for _, line := range strings.Split(content, "\n") {
		if !strings.Contains(line, "roundfix resolve") && !strings.Contains(line, "roundfix watch") && !strings.Contains(line, "roundfix implement") {
			continue
		}
		present := 0
		for _, flag := range []string{"--agent", "--model", "--reasoning-effort"} {
			if strings.Contains(line, flag) {
				present++
			}
		}
		if present != 0 && present != 3 {
			t.Fatalf("%s has a partial Agent Selection example: %q", label, strings.TrimSpace(line))
		}
	}
}

func cliTestRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repository root from %q", dir)
		}
		dir = parent
	}
}

type profilesShowTestResponse struct {
	Schema   string                    `json:"schema"`
	Profiles []profilesShowTestProfile `json:"profiles"`
}

type profilesShowTestProfile struct {
	Category             string                           `json:"category"`
	Source               string                           `json:"source"`
	InheritedFrom        string                           `json:"inherited_from"`
	RecommendationSource string                           `json:"recommendation_source"`
	Preferred            profilesShowTestSelection        `json:"preferred"`
	Fallbacks            []profilesShowTestSelection      `json:"fallbacks"`
	Recommendations      []profilesShowTestRecommendation `json:"recommendations"`
}

type profilesShowTestSelection struct {
	Runtime         string `json:"runtime"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoning_effort"`
}

type profilesShowTestRecommendation struct {
	Category          string                    `json:"category"`
	Rank              int                       `json:"rank"`
	Selection         profilesShowTestSelection `json:"selection"`
	Benchmark         string                    `json:"benchmark"`
	ResultPercent     float64                   `json:"result_percent"`
	AverageCostUSD    float64                   `json:"average_cost_usd"`
	SourceAsOf        string                    `json:"source_as_of"`
	Rationale         string                    `json:"rationale"`
	CategorySpecific  bool                      `json:"category_specific"`
	UnavailableReason string                    `json:"unavailable_reason,omitempty"`
}

func TestProfilesShowJSONRendersProfileAndRecommendations(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
profiles:
  backend:
    preferred:
      runtime: claude
      model: project-backend
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: project-fallback
        reasoning_effort: high
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "show", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles show exit = %d, stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no diagnostics on stderr, got %q", stderr.String())
	}
	response := decodeProfilesShowResponse(t, stdout.String())
	if response.Schema != "roundfix/profiles/v1" {
		t.Fatalf("schema = %q, want roundfix/profiles/v1", response.Schema)
	}
	if len(response.Profiles) != 1 {
		t.Fatalf("len(profiles) = %d, want 1", len(response.Profiles))
	}
	profile := response.Profiles[0]
	if profile.Category != "backend" || profile.Source != "project" || profile.InheritedFrom != "" {
		t.Fatalf("unexpected backend profile metadata: %+v", profile)
	}
	wantPreferred := profilesShowTestSelection{Runtime: "claude", Model: "project-backend", ReasoningEffort: "medium"}
	if profile.Preferred != wantPreferred {
		t.Fatalf("preferred = %+v, want %+v", profile.Preferred, wantPreferred)
	}
	wantFallback := profilesShowTestSelection{Runtime: "codex", Model: "project-fallback", ReasoningEffort: "high"}
	if len(profile.Fallbacks) != 1 || profile.Fallbacks[0] != wantFallback {
		t.Fatalf("fallbacks = %+v, want [%+v]", profile.Fallbacks, wantFallback)
	}
	if profile.RecommendationSource != "backend" {
		t.Fatalf("recommendation_source = %q, want backend", profile.RecommendationSource)
	}
	if len(profile.Recommendations) != 5 {
		t.Fatalf("len(recommendations) = %d, want 5", len(profile.Recommendations))
	}
	first := profile.Recommendations[0]
	if first.Rank != 1 || first.Category != "backend" {
		t.Fatalf("first recommendation identity = %+v, want backend rank 1", first)
	}
	if first.Selection != (profilesShowTestSelection{Runtime: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"}) {
		t.Fatalf("first recommendation selection = %+v", first.Selection)
	}
	if first.Benchmark != "DeepSWE v1.1" || first.ResultPercent != 69 || first.AverageCostUSD != 3.47 || first.SourceAsOf != "2026-07-16" {
		t.Fatalf("first recommendation evidence = %+v", first)
	}
	if first.CategorySpecific {
		t.Fatalf("category_specific = true, want false")
	}
	if !strings.Contains(first.Rationale, "complex repository changes") {
		t.Fatalf("first recommendation rationale missing backend context: %q", first.Rationale)
	}
	if profile.Preferred.Model == first.Selection.Model {
		t.Fatalf("configured preferred must remain primary; recommendation rank one replaced it")
	}
}

func TestProfilesShowOptionalCategoryReportsGeneralRecommendationSource(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "show", "--category", "data", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles show exit = %d, stderr=%q", code, stderr.String())
	}
	response := decodeProfilesShowResponse(t, stdout.String())
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	profile := response.Profiles[0]
	if profile.Category != "data" {
		t.Fatalf("category = %q, want data", profile.Category)
	}
	if profile.InheritedFrom != "general" {
		t.Fatalf("inherited_from = %q, want general", profile.InheritedFrom)
	}
	if profile.RecommendationSource != "general" {
		t.Fatalf("recommendation_source = %q, want general", profile.RecommendationSource)
	}
	if len(profile.Recommendations) != 5 {
		t.Fatalf("len(recommendations) = %d, want 5", len(profile.Recommendations))
	}
	if profile.Recommendations[0].Category != "data" || profile.Recommendations[0].Selection.Model != "gpt-5.6-sol" {
		t.Fatalf("optional recommendation should label data while reusing general order, got %+v", profile.Recommendations[0])
	}
}

func TestProfilesShowTextAndJSONAreByteStableAndConsistent(t *testing.T) {
	withCLIWorkspace(t)
	textFirst, textSecond := runProfilesShowTwice(t, []string{"profiles", "show"})
	if textFirst != textSecond {
		t.Fatalf("text output is not byte-stable\nfirst:\n%s\nsecond:\n%s", textFirst, textSecond)
	}
	jsonFirst, jsonSecond := runProfilesShowTwice(t, []string{"profiles", "show", "--json"})
	if jsonFirst != jsonSecond {
		t.Fatalf("JSON output is not byte-stable\nfirst:\n%s\nsecond:\n%s", jsonFirst, jsonSecond)
	}
	response := decodeProfilesShowResponse(t, jsonFirst)
	wantCategories := []string{"general", "backend", "frontend", "data", "infra", "docs", "test", "chore", "qa", "review"}
	if got := profileCategoriesForTest(response.Profiles); !reflect.DeepEqual(got, wantCategories) {
		t.Fatalf("categories = %v, want %v", got, wantCategories)
	}
	for _, profile := range response.Profiles {
		assertProfilesShowTextContainsProfile(t, textFirst, profile)
	}
}

func TestProfilesShowRejectsUnknownCategory(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "show", "--category", "design", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("exit = %d, want %d", code, exitPreflight)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout on usage error, got %q", stdout.String())
	}
	for _, want := range []string{"unknown profile category", "design", "general, backend, frontend, data, infra, docs, test, chore, qa, review"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr %q does not contain %q", stderr.String(), want)
		}
	}
}

func TestProfilesShowDoesNotMutateConfigOrRunState(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	userConfig := filepath.Join(homeDir, ".roundfix", "config.yml")
	projectConfig := filepath.Join(repoDir, ".roundfixrc.yml")
	mustMkdir(t, filepath.Dir(userConfig))
	mustWrite(t, userConfig, `
profiles:
  backend:
    preferred:
      runtime: codex
      model: user-backend
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: user-fallback
        reasoning_effort: max
`)
	mustWrite(t, projectConfig, `
profiles:
  frontend:
    preferred:
      runtime: claude
      model: project-frontend
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: project-fallback
        reasoning_effort: high
`)
	userBefore := mustRead(t, userConfig)
	projectBefore := mustRead(t, projectConfig)

	for _, args := range [][]string{
		{"profiles", "show", "--category", "backend"},
		{"profiles", "show", "--json"},
	} {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if code := Run(args, &stdout, &stderr); code != exitOK {
			t.Fatalf("Run(%v) exit = %d stderr=%q", args, code, stderr.String())
		}
		if stdout.Len() == 0 {
			t.Fatalf("Run(%v) produced no stdout", args)
		}
		if stderr.Len() != 0 {
			t.Fatalf("Run(%v) stderr = %q, want empty", args, stderr.String())
		}
	}
	if got := mustRead(t, userConfig); got != userBefore {
		t.Fatalf("User Config mutated\nbefore:\n%s\nafter:\n%s", userBefore, got)
	}
	if got := mustRead(t, projectConfig); got != projectBefore {
		t.Fatalf("Project Config mutated\nbefore:\n%s\nafter:\n%s", projectBefore, got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfilesConfigureFileWritesProjectProfileJSON(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runner := withSuccessfulProfilesConfigureProof(t)
	fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
	mustWrite(t, fragmentPath, `
profiles:
  backend:
    preferred:
      runtime: claude
      model: configured-backend
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: configured-fallback
        reasoning_effort: high
`)
	var preview string
	withProfilesConfigureConfirm(t, func(_ context.Context, _ io.Writer, got string) (bool, error) {
		preview = got
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles configure exit = %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr from stubbed confirmation, got %q", stderr.String())
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if response.Schema != "roundfix/profiles-configure/v1" || !response.Changed || response.Scope != "project" {
		t.Fatalf("unexpected configure response: %+v", response)
	}
	if response.Path != filepath.Join(repoDir, ".roundfixrc.yml") {
		t.Fatalf("path = %q, want project config", response.Path)
	}
	if len(response.Profiles) != 1 || response.Profiles[0].Category != roundconfig.CategoryBackend {
		t.Fatalf("profiles = %+v, want one backend profile", response.Profiles)
	}
	if len(runner.exactRequests) != 2 {
		t.Fatalf("file configure exact proofs = %#v, want preferred and fallback", runner.exactRequests)
	}
	for _, want := range []string{"Profile Configure Preview", "Scope: project", "Category: backend", "configured-backend", "configured-fallback"} {
		if !strings.Contains(preview, want) {
			t.Fatalf("preview missing %q in:\n%s", want, preview)
		}
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: homeDir, WorkDir: repoDir})
	if err != nil {
		t.Fatalf("load config after configure: %v", err)
	}
	resolved, err := roundconfig.ResolveProfile(loaded.Config, roundconfig.CategoryBackend, nil)
	if err != nil {
		t.Fatalf("resolve backend after configure: %v", err)
	}
	if resolved.Source != roundconfig.ProfileSourceProject {
		t.Fatalf("source = %q, want project", resolved.Source)
	}
	if resolved.Profile.Preferred.Model != "configured-backend" || resolved.Profile.Fallbacks[0].Model != "configured-fallback" {
		t.Fatalf("unexpected resolved profile: %+v", resolved.Profile)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfilesConfigureYesSkipsConfirmation(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runner := withSuccessfulProfilesConfigureProof(t)
	fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
	mustWrite(t, fragmentPath, `
backend:
  preferred:
    runtime: codex
    model: noninteractive-backend
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: noninteractive-fallback
      reasoning_effort: max
`)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		t.Fatal("--yes must not ask for write confirmation")
		return false, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--yes", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles configure --yes exit = %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr with --yes, got %q", stderr.String())
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if !response.Changed {
		t.Fatalf("expected --yes write to report changed, got %+v", response)
	}
	if len(runner.exactRequests) != 2 {
		t.Fatalf("--yes exact proofs = %#v, want preferred and fallback", runner.exactRequests)
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: homeDir, WorkDir: repoDir})
	if err != nil {
		t.Fatalf("load config after --yes configure: %v", err)
	}
	resolved, err := roundconfig.ResolveProfile(loaded.Config, roundconfig.CategoryBackend, nil)
	if err != nil {
		t.Fatalf("resolve backend after --yes configure: %v", err)
	}
	if resolved.Profile.Preferred.Model != "noninteractive-backend" {
		t.Fatalf("preferred model = %q, want noninteractive-backend", resolved.Profile.Preferred.Model)
	}
}

func TestProfilesJSONSuccessReturnsEncoderFailures(t *testing.T) {
	writeErr := errors.New("write failed")
	writer := failingWriter{err: writeErr}

	configureErr := printProfilesConfigureSuccess(profilesConfigureRequest{json: true}, roundconfig.ProfileConfigResult{Scope: "project"}, writer)
	if !errors.Is(configureErr, writeErr) {
		t.Fatalf("profiles configure encode error = %v, want %v", configureErr, writeErr)
	}

	validateErr := printProfilesValidateSuccess(profilesValidateRequest{json: true}, profileProofResult{}, writer)
	if !errors.Is(validateErr, writeErr) {
		t.Fatalf("profiles validate encode error = %v, want %v", validateErr, writeErr)
	}
}

func TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	runner := withSuccessfulProfilesConfigureProof(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	validFragment := filepath.Join(repoDir, "valid-profile.yml")
	mustWrite(t, validFragment, `
backend:
  preferred:
    runtime: codex
    model: dry-run-backend
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: dry-run-fallback
      reasoning_effort: max
`)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		t.Fatal("dry-run must not ask for write confirmation")
		return false, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", validFragment, "--dry-run", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("dry-run exit = %d stderr=%q", code, stderr.String())
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if response.Changed {
		t.Fatalf("dry-run changed = true, response=%+v", response)
	}
	if len(runner.exactRequests) != 2 {
		t.Fatalf("dry-run exact proofs = %#v, want preferred and fallback", runner.exactRequests)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("dry-run mutated config\nwant: %q\n got: %q", original, got)
	}

	invalidFragment := filepath.Join(repoDir, "invalid-profile.yml")
	mustWrite(t, invalidFragment, `
backend:
  preferred:
    runtime: codex
    model: broken
    reasoning_effort: high
`)
	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"profiles", "configure", "--scope", "project", "--file", invalidFragment, "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("invalid configure exit = %d, want %d", code, exitPreflight)
	}
	response = decodeProfilesConfigureResponse(t, stdout.String())
	if response.Changed || !strings.Contains(response.Error, "fallbacks is required") {
		t.Fatalf("unexpected invalid configure response: %+v", response)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("invalid configure mutated config\nwant: %q\n got: %q", original, got)
	}
}

func TestProfilesConfigureProofRunsBeforeConfirmationAndWrite(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	fragmentPath := filepath.Join(repoDir, "backend-profile.yml")
	mustWrite(t, fragmentPath, `
backend:
  preferred:
    runtime: codex
    model: exact-preferred
    reasoning_effort: high
  fallbacks:
    - runtime: claude
      model: exact-fallback
      reasoning_effort: ""
`)
	runner := withSuccessfulProfilesConfigureProof(t)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		if len(runner.exactRequests) != 2 {
			t.Fatalf("confirmation reached after %d exact proofs, want 2", len(runner.exactRequests))
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("config mutated before confirmation\nwant: %q\n got: %q", original, got)
		}
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("configure exit = %d stderr=%q", code, stderr.String())
	}
	if len(runner.exactRequests) != 2 {
		t.Fatalf("exact proof requests = %#v, want preferred and fallback", runner.exactRequests)
	}
	if got := runner.exactRequests[0].Runtime; got.Model != "exact-preferred" || got.ReasoningEffort != "high" {
		t.Fatalf("preferred proof = %+v", got)
	}
	if got := runner.exactRequests[1].Runtime; got.ID != "claude" || got.Model != "exact-fallback" || got.ReasoningEffort != "" {
		t.Fatalf("fallback proof = %+v", got)
	}
	if response := decodeProfilesConfigureResponse(t, stdout.String()); !response.Changed {
		t.Fatalf("configure response = %+v, want changed", response)
	}
}

func TestProfilesConfigureInteractiveProofKeepsRecommendationsAdvisory(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	withProfilesConfigureInput(t, "backend\ncodex\ninteractive-choice\nhigh\nclaude\ninteractive-fallback\nmedium\n\n")
	runner := withSuccessfulProfilesConfigureProof(t)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) { return false, nil })
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("interactive configure exit = %d stderr=%q", code, stderr.String())
	}
	if len(runner.exactRequests) != 2 {
		t.Fatalf("interactive exact proofs = %#v, want preferred and fallback", runner.exactRequests)
	}
	if got := runner.exactRequests[0].Runtime; got.Model != "interactive-choice" || got.ReasoningEffort != "high" {
		t.Fatalf("interactive preferred proof = %+v", got)
	}
	if got := runner.exactRequests[1].Runtime; got.ID != "claude" || got.Model != "interactive-fallback" || got.ReasoningEffort != "medium" {
		t.Fatalf("interactive fallback proof = %+v", got)
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if response.Changed || len(response.Profiles) != 1 || response.Profiles[0].Preferred.Model != "interactive-choice" || response.Profiles[0].Fallbacks[0].Model != "interactive-fallback" {
		t.Fatalf("interactive response inserted or changed a selection: %+v", response)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("interactive decline mutated config\nwant: %q\n got: %q", original, got)
	}
}

func TestProfilesConfigureFallbackFailurePrecedesProofAndPreservesBytes(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	runner := withSuccessfulProfilesConfigureProof(t)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		t.Fatal("invalid fallback must fail before confirmation")
		return false, nil
	})

	tests := []struct {
		name     string
		fragment string
	}{
		{
			name: "missing fallback",
			fragment: `
backend:
  preferred:
    runtime: codex
    model: only-selection
    reasoning_effort: high
`,
		},
		{
			name: "fallback duplicates preferred",
			fragment: `
backend:
  preferred:
    runtime: codex
    model: only-selection
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: only-selection
      reasoning_effort: high
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fragmentPath := filepath.Join(repoDir, strings.ReplaceAll(tt.name, " ", "-")+".yml")
			mustWrite(t, fragmentPath, tt.fragment)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json", "--yes"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("configure exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
			}
			response := decodeProfilesConfigureResponse(t, stdout.String())
			if response.Changed || !strings.Contains(response.Error, "one additional distinct authorized and proven Agent Selection is required") {
				t.Fatalf("fallback response = %+v", response)
			}
			if got := mustRead(t, configPath); got != original {
				t.Fatalf("fallback failure mutated config\nwant: %q\n got: %q", original, got)
			}
		})
	}
	if len(runner.exactRequests) != 0 {
		t.Fatalf("structural fallback failures created disposable Sessions: %#v", runner.exactRequests)
	}
}

func TestProfilesConfigureProofFailureYesJSONPreservesBytes(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	fragmentPath := filepath.Join(repoDir, "unsupported-profile.yml")
	mustWrite(t, fragmentPath, `
backend:
  preferred:
    runtime: codex
    model: unavailable-model
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: proven-fallback
      reasoning_effort: xhigh
`)
	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{}, &agent.SelectionUnsupportedError{
				Kind:            agent.SelectionModelNotAdvertised,
				Runtime:         req.Runtime.ID,
				Model:           req.Runtime.Model,
				ReasoningEffort: req.Runtime.ReasoningEffort,
			}
		},
	}
	withAgentRunner(t, runner)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		t.Fatal("--yes proof failure must not confirm")
		return false, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--yes", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("proof failure exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if response.Changed || !strings.Contains(response.Error, agent.SelectionModelNotAdvertised) {
		t.Fatalf("proof failure response = %+v", response)
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("proof failure mutated config\nwant: %q\n got: %q", original, got)
	}
}

func TestProfilesConfigureProofCleanupFailureAndDeclinePreserveBytes(t *testing.T) {
	t.Run("cleanup failure", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		original := "watch:\n  max_rounds: 4\n"
		mustWrite(t, configPath, original)
		fragmentPath := filepath.Join(repoDir, "cleanup-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("cleanup-preferred", "cleanup-fallback"))
		withAgentRunner(t, &profileReadinessExactRunner{prove: func(agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{}, &agent.AgentSessionCleanupError{Session: "disposable", Err: errors.New("close failed")}
		}})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json", "--yes"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("cleanup failure exit = %d stderr=%q", code, stderr.String())
		}
		response := decodeProfilesConfigureResponse(t, stdout.String())
		if response.Changed || !strings.Contains(response.Error, agent.SessionCleanupFailed) {
			t.Fatalf("cleanup response = %+v", response)
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("cleanup failure mutated config\nwant: %q\n got: %q", original, got)
		}
	})

	t.Run("decline", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		configPath := filepath.Join(repoDir, ".roundfixrc.yml")
		original := "watch:\n  max_rounds: 4\n"
		mustWrite(t, configPath, original)
		fragmentPath := filepath.Join(repoDir, "decline-profile.yml")
		mustWrite(t, fragmentPath, profilesConfigureTestFragment("decline-preferred", "decline-fallback"))
		withSuccessfulProfilesConfigureProof(t)
		withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) { return false, nil })
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("decline exit = %d stderr=%q", code, stderr.String())
		}
		if response := decodeProfilesConfigureResponse(t, stdout.String()); response.Changed {
			t.Fatalf("decline response = %+v", response)
		}
		if got := mustRead(t, configPath); got != original {
			t.Fatalf("decline mutated config\nwant: %q\n got: %q", original, got)
		}
	})
}

func TestProfilesConfigureJSONOutputFailureDoesNotMutateConfig(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	original := "watch:\n  max_rounds: 4\n"
	mustWrite(t, configPath, original)
	fragmentPath := filepath.Join(repoDir, "output-profile.yml")
	mustWrite(t, fragmentPath, profilesConfigureTestFragment("output-preferred", "output-fallback"))
	withSuccessfulProfilesConfigureProof(t)
	writeErr := errors.New("stdout failed")
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "project", "--file", fragmentPath, "--json", "--yes"}, failingWriter{err: writeErr}, &stderr)

	if code != exitRunFailed {
		t.Fatalf("output failure exit = %d, want %d stderr=%q", code, exitRunFailed, stderr.String())
	}
	if got := mustRead(t, configPath); got != original {
		t.Fatalf("output failure mutated config\nwant: %q\n got: %q", original, got)
	}
}

func TestProfilesConfigureInteractiveRequiresCompleteFallbackBeforeConfirm(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withProfilesConfigureInput(t, "backend\ncodex\ninteractive-backend\nhigh\n\n")
	confirmCalls := 0
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
		confirmCalls++
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "configure", "--scope", "user", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("interactive configure exit = %d, want %d", code, exitPreflight)
	}
	if confirmCalls != 0 {
		t.Fatalf("confirmation called %d time(s) before fallback completed", confirmCalls)
	}
	response := decodeProfilesConfigureResponse(t, stdout.String())
	if response.Changed || !strings.Contains(response.Error, "fallbacks must include at least one complete Agent Selection before confirmation") {
		t.Fatalf("unexpected interactive failure response: %+v", response)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfilesValidateDeduplicatesProofsAndReportsEveryReference(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	runner := &fakeAgentRunner{}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "validate", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles validate exit = %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if response.Schema != "roundfix/profiles-validate/v1" || !response.OK {
		t.Fatalf("unexpected validate response: %+v", response)
	}
	if len(runner.probeRequests) != 3 {
		t.Fatalf("expected three unique tuple probes, got %#v", runner.probeRequests)
	}
	wantModels := []string{"gpt-5.6-sol", "gpt-5.5", "claude-opus-5"}
	for index, want := range wantModels {
		if runner.probeRequests[index].Runtime.Model != want {
			t.Fatalf("probe %d model = %q, want %q", index, runner.probeRequests[index].Runtime.Model, want)
		}
	}
	if len(response.Proofs) != 3 {
		t.Fatalf("len(proofs) = %d, want 3", len(response.Proofs))
	}
	assertProofReferences(t, response.Proofs[0], []string{"general/preferred", "backend/preferred", "frontend/fallback", "qa/preferred", "review/preferred"})
	if runner.calls != 0 {
		t.Fatalf("profiles validate must not send Agent prompts, calls=%d", runner.calls)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfilesValidateFailedProofNamesTupleAffectedCategoriesAndRecovery(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	mustWrite(t, configPath, `
profiles:
  backend:
    preferred:
      runtime: codex
      model: broken-backend
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: backup-backend
        reasoning_effort: ""
`)
	before := mustRead(t, configPath)
	runner := &fakeAgentRunner{
		probe: func(req agent.ProbeRequest) error {
			if req.Runtime.Model == "broken-backend" {
				return &agent.SelectionPreflightError{
					Runtime:         req.Runtime.ID,
					Model:           req.Runtime.Model,
					ReasoningEffort: req.Runtime.ReasoningEffort,
					Err:             errors.New("adapter rejected tuple"),
				}
			}
			return nil
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("profiles validate exit = %d, want %d", code, exitPreflight)
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if response.OK || !strings.Contains(response.Error, "broken-backend") || !strings.Contains(response.Error, "backend preferred") || !strings.Contains(response.Error, "adapter rejected tuple") || !strings.Contains(response.Error, "roundfix profiles configure") {
		t.Fatalf("unexpected failed validation response: %+v", response)
	}
	if len(response.Proofs) == 0 || response.Proofs[0].Status != "failed" || !strings.Contains(response.Proofs[0].Error, "adapter rejected tuple") {
		t.Fatalf("expected failed proof details, got %+v", response.Proofs)
	}
	for _, want := range []string{"broken-backend", "backend preferred", "adapter rejected tuple", "roundfix profiles validate"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr missing %q in %q", want, stderr.String())
		}
	}
	if got := mustRead(t, configPath); got != before {
		t.Fatalf("profiles validate mutated Project Config\nwant: %q\n got: %q", before, got)
	}
	if runner.calls != 0 {
		t.Fatalf("profiles validate must not send Agent prompts, calls=%d", runner.calls)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProveProfileSelectionsDeduplicatesReferencesAndStartsFreshProofPass(t *testing.T) {
	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{
				Runtime:         req.Runtime.ID,
				Model:           req.Runtime.Model,
				ReasoningEffort: req.Runtime.ReasoningEffort,
				Assignment: agent.SelectionAssignment{
					Encoding: agent.SelectionEncodingIndependent,
				},
				Adapter: agent.AdapterEvidence{
					Command: "official-adapter",
					Version: "1.1.4",
				},
				Status: agent.SelectionProofStatusProven,
			}, nil
		},
	}
	categories := roundconfig.RequiredWorkCategories()

	first := proveProfileSelections(context.Background(), roundconfig.Builtin(), categories, "/workspace", runner)
	second := proveProfileSelections(context.Background(), roundconfig.Builtin(), categories, "/workspace", runner)

	if first.Err != nil || second.Err != nil {
		t.Fatalf("profile readiness errors: first=%v second=%v", first.Err, second.Err)
	}
	if len(runner.exactRequests) != 6 {
		t.Fatalf("exact proof requests = %d, want two fresh passes of three tuples", len(runner.exactRequests))
	}
	if runner.probeCalls != 0 {
		t.Fatalf("legacy error-only probes = %d, want exact proof results", runner.probeCalls)
	}
	if len(first.Proofs) != 3 {
		t.Fatalf("proofs = %d, want three deduplicated tuples", len(first.Proofs))
	}
	proof := first.Proofs[0]
	if proof.Encoding != agent.SelectionEncodingIndependent || proof.AdapterCommand != "official-adapter" || proof.AdapterVersion != "1.1.4" {
		t.Fatalf("proof metadata = %+v", proof)
	}
	assertProofReferences(t, proof, []string{"general/preferred", "backend/preferred", "frontend/fallback", "qa/preferred", "review/preferred"})
}

func TestProveProfileSelectionsRetainsStableFallbackPositions(t *testing.T) {
	config := roundconfig.Builtin()
	shared := config.Profiles[roundconfig.CategoryBackend].Profile.Fallbacks[0]
	frontend := config.Profiles[roundconfig.CategoryFrontend]
	frontend.Profile.Fallbacks = []roundconfig.AgentSelection{
		{Runtime: "claude", Model: "frontend-first-fallback", ReasoningEffort: ""},
		shared,
	}
	config.Profiles[roundconfig.CategoryFrontend] = frontend
	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{Runtime: req.Runtime.ID, Model: req.Runtime.Model, ReasoningEffort: req.Runtime.ReasoningEffort}, nil
		},
	}

	result := proveProfileSelections(context.Background(), config, []roundconfig.WorkCategory{
		roundconfig.CategoryBackend,
		roundconfig.CategoryFrontend,
	}, "/workspace", runner)

	if result.Err != nil {
		t.Fatalf("profile readiness error = %v", result.Err)
	}
	if len(runner.exactRequests) != 4 {
		t.Fatalf("exact proof requests = %d, want four unique tuples", len(runner.exactRequests))
	}
	sharedProof := result.Proofs[1]
	if sharedProof.Selection != shared || len(sharedProof.References) != 2 {
		t.Fatalf("shared fallback proof = %+v", sharedProof)
	}
	if got := sharedProof.References[0]; got.Category != roundconfig.CategoryBackend || got.Role != "fallback" || got.FallbackIndex != 1 || got.Source != roundconfig.ProfileSourceBuiltIn {
		t.Fatalf("backend fallback reference = %+v", got)
	}
	if got := sharedProof.References[1]; got.Category != roundconfig.CategoryFrontend || got.Role != "fallback" || got.FallbackIndex != 2 || got.Source != roundconfig.ProfileSourceBuiltIn {
		t.Fatalf("frontend fallback reference = %+v", got)
	}
}

func TestProfileOperationalPreflightMatchesProfilesValidateClassifiedFailure(t *testing.T) {
	classified := &agent.SelectionUnsupportedError{
		Kind:                agent.SelectionReasoningControlNotAdvertised,
		Runtime:             "codex",
		Model:               "gpt-5.6-sol",
		ReasoningEffort:     "high",
		AdvertisedModels:    []string{"gpt-5.6-sol", "gpt-5.5"},
		AdvertisedReasoning: []string{"low", "medium"},
	}
	newRunner := func() *profileReadinessExactRunner {
		return &profileReadinessExactRunner{
			prove: func(agent.ProbeRequest) (agent.SelectionProof, error) {
				return agent.SelectionProof{}, classified
			},
		}
	}
	config := roundconfig.Builtin()
	categories := []roundconfig.WorkCategory{roundconfig.CategoryBackend}

	validation := proveProfileSelections(context.Background(), config, categories, "/workspace", newRunner())
	operational, operationalErr := runProfileOperationalPreflight(context.Background(), commandRequest{name: "implement"}, config, categories, "/workspace", newRunner(), io.Discard)

	if validation.Err == nil || operationalErr == nil {
		t.Fatalf("expected classified failures: validation=%v operational=%v", validation.Err, operationalErr)
	}
	if len(validation.Proofs) == 0 || len(operational.Proofs) == 0 {
		t.Fatalf("missing failed proof evidence: validation=%+v operational=%+v", validation.Proofs, operational.Proofs)
	}
	validateProof := validation.Proofs[0]
	preflightProof := operational.Proofs[0]
	if !reflect.DeepEqual(validateProof, preflightProof) {
		t.Fatalf("consumer proof evidence differs:\nvalidate:   %+v\noperational: %+v", validateProof, preflightProof)
	}
	if validateProof.Classification != agent.SelectionReasoningControlNotAdvertised {
		t.Fatalf("classification = %q", validateProof.Classification)
	}
	if validateProof.NextAction == "" || len(validateProof.AdvertisedModels) != 2 || len(validateProof.AdvertisedReasoning) != 2 {
		t.Fatalf("incomplete bounded failure evidence: %+v", validateProof)
	}

	var jsonStdout bytes.Buffer
	var jsonStderr bytes.Buffer
	if code := printProfilesValidateError(profilesValidateRequest{json: true}, validation, validation.Err, &jsonStdout, &jsonStderr); code != exitPreflight {
		t.Fatalf("JSON failure exit = %d", code)
	}
	response := decodeProfilesValidateResponse(t, jsonStdout.String())
	if response.Schema != profilesValidateSchema || len(response.Proofs) == 0 || !reflect.DeepEqual(response.Proofs[0], validateProof) {
		t.Fatalf("unexpected JSON failure response: %+v", response)
	}
	for _, want := range []string{validateProof.Classification, "backend preferred", validateProof.NextAction} {
		if !strings.Contains(jsonStderr.String(), want) {
			t.Fatalf("JSON stderr missing %q in %q", want, jsonStderr.String())
		}
	}

	var textStderr bytes.Buffer
	if code := printProfilesValidateError(profilesValidateRequest{}, validation, validation.Err, io.Discard, &textStderr); code != exitPreflight {
		t.Fatalf("text failure exit = %d", code)
	}
	for _, want := range []string{validateProof.Classification, "backend preferred", validateProof.NextAction} {
		if !strings.Contains(textStderr.String(), want) {
			t.Fatalf("text stderr missing %q in %q", want, textStderr.String())
		}
	}
}

func TestInvocationProfileOverrideOmittedUsesTaskQAAndReviewProfiles(t *testing.T) {
	runner := &fakeAgentRunner{}
	var stderr bytes.Buffer
	graph := &spec.Graph{Tasks: []spec.Task{
		{ID: "task_frontend", Status: spec.StatusPending, Type: spec.TaskTypeFrontend},
		{ID: "task_backend", Status: spec.StatusPending, Type: spec.TaskTypeBackend},
	}}
	categories := implementProfileCategories(graph, true)

	result, err := runProfileOperationalPreflight(context.Background(), commandRequest{name: "implement"}, roundconfig.Builtin(), categories, "/workspace", runner, &stderr)

	if err != nil {
		t.Fatalf("operational profile preflight error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no warning without invocation override, got %q", stderr.String())
	}
	wantModels := []string{"gpt-5.6-sol", "gpt-5.5", "claude-opus-5"}
	if got := probeRequestModels(runner.probeRequests); !reflect.DeepEqual(got, wantModels) {
		t.Fatalf("probe models = %v, want %v", got, wantModels)
	}
	if len(result.Proofs) != 3 {
		t.Fatalf("len(proofs) = %d, want 3", len(result.Proofs))
	}
	assertProofReferences(t, result.Proofs[0], []string{"backend/preferred", "frontend/fallback", "qa/preferred"})
	assertProofReferences(t, result.Proofs[1], []string{"backend/fallback", "qa/fallback"})
	assertProofReferences(t, result.Proofs[2], []string{"frontend/preferred"})
	for _, request := range runner.probeRequests {
		if request.WorkDir != "/workspace" {
			t.Fatalf("probe WorkDir = %q, want /workspace", request.WorkDir)
		}
	}
	if runner.calls != 0 || len(runner.fallbackModels) != 0 {
		t.Fatalf("profile preflight must not prompt or discover fallbacks, calls=%d fallback=%#v", runner.calls, runner.fallbackModels)
	}

	reviewRunner := &fakeAgentRunner{}
	reviewResult, err := runProfileOperationalPreflight(context.Background(), commandRequest{name: "resolve"}, roundconfig.Builtin(), reviewProfileCategories(), "/workspace", reviewRunner, io.Discard)
	if err != nil {
		t.Fatalf("review profile preflight error = %v", err)
	}
	if reviewResult.Override != nil {
		t.Fatalf("review profile preflight override = %+v, want nil", reviewResult.Override)
	}
	assertProofReferences(t, reviewResult.Proofs[0], []string{"review/preferred"})
	assertProofReferences(t, reviewResult.Proofs[1], []string{"review/fallback"})
}

func TestInvocationProfileOverrideAppliesAcrossCategoriesPreservesFallbacksAndWarns(t *testing.T) {
	runner := &fakeAgentRunner{}
	var stderr bytes.Buffer
	graph := &spec.Graph{Tasks: []spec.Task{
		{ID: "task_backend", Status: spec.StatusPending, Type: spec.TaskTypeBackend},
	}}
	req := commandRequest{
		name:               "implement",
		agent:              "codex",
		agentSet:           true,
		model:              "gpt-5.6-sol",
		modelSet:           true,
		reasoningEffort:    "high",
		reasoningEffortSet: true,
	}
	categories := implementProfileCategories(graph, true)

	result, err := runProfileOperationalPreflight(context.Background(), req, roundconfig.Builtin(), categories, "/workspace", runner, &stderr)

	if err != nil {
		t.Fatalf("operational profile preflight error = %v", err)
	}
	if result.Override == nil || result.Override.Model != "gpt-5.6-sol" || result.Override.ReasoningEffort != "high" {
		t.Fatalf("unexpected invocation override: %+v", result.Override)
	}
	wantModels := []string{"gpt-5.6-sol", "gpt-5.5"}
	if got := probeRequestModels(runner.probeRequests); !reflect.DeepEqual(got, wantModels) {
		t.Fatalf("probe models = %v, want %v", got, wantModels)
	}
	assertProofReferences(t, result.Proofs[0], []string{"backend/preferred", "qa/preferred"})
	if !strings.Contains(stderr.String(), "one invocation Agent Selection override applies to Agent Work Categories backend, qa") ||
		!strings.Contains(stderr.String(), "Fallback Chains are preserved") {
		t.Fatalf("expected cross-category warning, got %q", stderr.String())
	}
	if len(result.Warnings) != 1 || result.Warnings[0] != strings.TrimSpace(stderr.String()) {
		t.Fatalf("warning metadata = %#v, stderr=%q", result.Warnings, stderr.String())
	}
	profiles, err := operationalAgentSelectionProfiles(roundconfig.Builtin(), categories, result.Override)
	if err != nil {
		t.Fatalf("operationalAgentSelectionProfiles() error = %v", err)
	}
	for _, category := range categories {
		withoutOverride, err := roundconfig.ResolveProfile(roundconfig.Builtin(), category, nil)
		if err != nil {
			t.Fatalf("ResolveProfile(%s) error = %v", category, err)
		}
		if got := profiles[category]; !reflect.DeepEqual(got.Profile.Fallbacks, withoutOverride.Profile.Fallbacks) {
			t.Fatalf("%s fallback chain = %+v, want %+v", category, got.Profile.Fallbacks, withoutOverride.Profile.Fallbacks)
		}
	}
}

func TestProfilesShowReportsUnavailableRecommendationWithoutReordering(t *testing.T) {
	unavailableSelection := roundconfig.AgentSelection{
		Runtime:         "codex",
		Model:           "gpt-5.6-terra",
		ReasoningEffort: "max",
	}
	response, err := buildProfilesShowResponseWithAvailability(roundconfig.Builtin(), []roundconfig.WorkCategory{roundconfig.CategoryBackend}, map[roundconfig.AgentSelection]string{
		unavailableSelection: "adapter proof rejected tuple",
	})
	if err != nil {
		t.Fatalf("build profiles show response: %v", err)
	}
	if len(response.Profiles) != 1 {
		t.Fatalf("len(profiles) = %d, want 1", len(response.Profiles))
	}
	recommendations := response.Profiles[0].Recommendations
	if len(recommendations) != 5 {
		t.Fatalf("len(recommendations) = %d, want 5", len(recommendations))
	}
	for index, recommendation := range recommendations {
		if recommendation.Rank != index+1 {
			t.Fatalf("rank at index %d = %d, want %d", index, recommendation.Rank, index+1)
		}
	}
	if recommendations[0].Selection.Model != "gpt-5.6-sol" || recommendations[1].Selection.Model != "gpt-5.6-terra" || recommendations[2].Selection.Model != "claude-fable-5" {
		t.Fatalf("recommendation order changed: %+v", recommendations)
	}
	if recommendations[1].UnavailableReason != "adapter proof rejected tuple" {
		t.Fatalf("unavailable_reason = %q, want proof rejection", recommendations[1].UnavailableReason)
	}
	if recommendations[0].UnavailableReason != "" || recommendations[2].UnavailableReason != "" {
		t.Fatalf("unexpected unavailable marker outside rejected tuple: %+v", recommendations)
	}
}

func TestModelRecommendationsUseOfficialCatalogModels(t *testing.T) {
	for _, category := range roundconfig.AllWorkCategories() {
		t.Run(string(category), func(t *testing.T) {
			recommendations, source, ok := roundconfig.ModelRecommendations(category)
			if !ok {
				t.Fatalf("ModelRecommendations(%q) missing", category)
			}
			if len(recommendations) != 5 {
				t.Fatalf("len(recommendations) = %d, want 5", len(recommendations))
			}
			if source == "" {
				t.Fatal("recommendation source is empty")
			}
			seenModels := map[string]bool{}
			for index, recommendation := range recommendations {
				if recommendation.Category != category {
					t.Fatalf("recommendation category = %q, want %q", recommendation.Category, category)
				}
				if recommendation.Rank != index+1 {
					t.Fatalf("rank = %d, want %d", recommendation.Rank, index+1)
				}
				if seenModels[recommendation.Selection.Model] {
					t.Fatalf("duplicate model %q in %q recommendations", recommendation.Selection.Model, category)
				}
				seenModels[recommendation.Selection.Model] = true
				if !modelCatalogContainsSelection(recommendation.Selection) {
					t.Fatalf("recommendation uses non-catalog official model: %+v", recommendation.Selection)
				}
				if recommendation.Benchmark == "" || recommendation.SourceAsOf != roundconfig.ModelRecommendationSnapshotDate {
					t.Fatalf("recommendation has incomplete source evidence: %+v", recommendation)
				}
				if recommendation.Rationale == "" {
					t.Fatalf("recommendation rationale is empty: %+v", recommendation)
				}
				if recommendation.CategorySpecific {
					t.Fatalf("category_specific = true, want false for initial snapshot")
				}
			}
		})
	}
}

func decodeProfilesShowResponse(t *testing.T, output string) profilesShowTestResponse {
	t.Helper()
	var response profilesShowTestResponse
	decoder := json.NewDecoder(strings.NewReader(output))
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("decode profiles JSON %q: %v", output, err)
	}
	if decoder.More() {
		t.Fatalf("expected one JSON object, got trailing content in %q", output)
	}
	return response
}

func decodeProfilesConfigureResponse(t *testing.T, output string) profilesConfigureResponse {
	t.Helper()
	var response profilesConfigureResponse
	decoder := json.NewDecoder(strings.NewReader(output))
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("decode profiles configure JSON %q: %v", output, err)
	}
	return response
}

func decodeProfilesValidateResponse(t *testing.T, output string) profilesValidateResponse {
	t.Helper()
	var response profilesValidateResponse
	decoder := json.NewDecoder(strings.NewReader(output))
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("decode profiles validate JSON %q: %v", output, err)
	}
	return response
}

func assertProofReferences(t *testing.T, proof profileProofReport, want []string) {
	t.Helper()
	got := make([]string, 0, len(proof.References))
	for _, reference := range proof.References {
		value := string(reference.Category) + "/" + reference.Role
		if reference.Role == "fallback" {
			value = string(reference.Category) + "/fallback"
		}
		got = append(got, value)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("proof references = %v, want %v", got, want)
	}
}

func probeRequestModels(requests []agent.ProbeRequest) []string {
	models := make([]string, 0, len(requests))
	for _, request := range requests {
		models = append(models, request.Runtime.Model)
	}
	return models
}

func runProfilesShowTwice(t *testing.T, args []string) (string, string) {
	t.Helper()
	var firstStdout bytes.Buffer
	var firstStderr bytes.Buffer
	if code := Run(args, &firstStdout, &firstStderr); code != exitOK {
		t.Fatalf("first Run(%v) exit = %d stderr=%q", args, code, firstStderr.String())
	}
	if firstStderr.Len() != 0 {
		t.Fatalf("first Run(%v) stderr = %q, want empty", args, firstStderr.String())
	}
	var secondStdout bytes.Buffer
	var secondStderr bytes.Buffer
	if code := Run(args, &secondStdout, &secondStderr); code != exitOK {
		t.Fatalf("second Run(%v) exit = %d stderr=%q", args, code, secondStderr.String())
	}
	if secondStderr.Len() != 0 {
		t.Fatalf("second Run(%v) stderr = %q, want empty", args, secondStderr.String())
	}
	return firstStdout.String(), secondStdout.String()
}

func profileCategoriesForTest(profiles []profilesShowTestProfile) []string {
	categories := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		categories = append(categories, profile.Category)
	}
	return categories
}

func assertProfilesShowTextContainsProfile(t *testing.T, output string, profile profilesShowTestProfile) {
	t.Helper()
	for _, want := range []string{
		"Category: " + profile.Category,
		"Profile source: " + profile.Source,
		"Profile inherited from: " + emptyDashForTest(profile.InheritedFrom),
		"Preferred Selection: " + selectionStringForTest(profile.Preferred),
		"Recommendation source: " + profile.RecommendationSource,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("profiles text missing %q in:\n%s", want, output)
		}
	}
	for index, fallback := range profile.Fallbacks {
		want := fmt.Sprintf("%d. %s", index+1, selectionStringForTest(fallback))
		if !strings.Contains(output, want) {
			t.Fatalf("profiles text missing fallback %q in:\n%s", want, output)
		}
	}
	for _, recommendation := range profile.Recommendations {
		want := fmt.Sprintf("%d. %s — %s %s, average cost $%.2f, source %s, category_specific=%t", recommendation.Rank, selectionStringForTest(recommendation.Selection), recommendation.Benchmark, formatPercentForTest(recommendation.ResultPercent), recommendation.AverageCostUSD, recommendation.SourceAsOf, recommendation.CategorySpecific)
		if !strings.Contains(output, want) {
			t.Fatalf("profiles text missing recommendation %q in:\n%s", want, output)
		}
	}
}

func selectionStringForTest(selection profilesShowTestSelection) string {
	return selection.Runtime + " / " + selection.Model + " / " + selection.ReasoningEffort
}

func emptyDashForTest(value string) string {
	if value == "" {
		return "—"
	}
	return value
}

func formatPercentForTest(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%d%%", int(value))
	}
	return fmt.Sprintf("%.1f%%", value)
}

func modelCatalogContainsSelection(selection roundconfig.AgentSelection) bool {
	for _, choice := range agent.ModelCatalog(selection.Runtime) {
		if choice.Value == selection.Model {
			return true
		}
	}
	return false
}

// seedRunsForListColumns seeds one terminal and one Active Run in the
// current repository plus one Active Run in another repository, and pins
// the listing clock so durations are byte-stable.
func seedRunsForListColumns(t *testing.T, homeDir, repoDir, otherRepo string) []store.Run {
	t.Helper()
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	runs := seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:        store.KindResolve,
			state:       store.StateClean,
			gitRoot:     repoDir,
			branch:      "feature/review",
			prNumber:    "123",
			createdAt:   base,
			completedAt: base.Add(42 * time.Minute),
		},
		{
			kind:      store.KindImplement,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    "ma/spec-run",
			specSlug:  "0017-run-discovery",
			createdAt: base.Add(2 * time.Minute),
		},
		{
			kind:      store.KindWatch,
			state:     store.StateActive,
			gitRoot:   otherRepo,
			branch:    "feature/other",
			prNumber:  "999",
			createdAt: base.Add(4 * time.Minute),
		},
	})
	withRunsListNow(t, base.Add(14*time.Minute))
	return runs
}

func TestRunRunsListPrintsStableColumnsNewestFirst(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	runs := seedRunsForListColumns(t, homeDir, repoDir, otherRepo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
	}
	want := fmt.Sprintf(
		"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run\n",
		runs[1].ID,
	)
	if stdout.String() != want {
		t.Fatalf("unexpected runs list output:\n got: %q\nwant: %q", stdout.String(), want)
	}
	if stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
		t.Fatalf("expected terminal-hidden note, got %q", stderr.String())
	}
}

func TestRunRunsListStateFlagFiltersAndNotes(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	runs := seedRunsForListColumns(t, homeDir, repoDir, otherRepo)
	activeLine := fmt.Sprintf(
		"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run\n",
		runs[1].ID,
	)
	terminalLine := fmt.Sprintf(
		"%s  Clean  resolve  pr:123  codex  2026-07-06T12:00:00Z  42m  feature/review\n",
		runs[0].ID,
	)
	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "state all includes terminal Runs without a note",
			args:       []string{"runs", "list", "--state", "all"},
			wantStdout: activeLine + terminalLine,
			wantStderr: "",
		},
		{
			name:       "state terminal excludes Active Runs and notes them",
			args:       []string{"runs", "list", "--state", "terminal"},
			wantStdout: terminalLine,
			wantStderr: "(1 active Run(s) hidden; use --state all)\n",
		},
		{
			name: "all flag appends the repository column",
			args: []string{"runs", "list", "--all", "--state", "all"},
			wantStdout: fmt.Sprintf(
				"%s  Active  watch  pr:999  codex  2026-07-06T12:04:00Z  running 10m  feature/other  %s\n",
				runs[2].ID, otherRepo,
			) + fmt.Sprintf(
				"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run  %s\n",
				runs[1].ID, repoDir,
			) + fmt.Sprintf(
				"%s  Clean  resolve  pr:123  codex  2026-07-06T12:00:00Z  42m  feature/review  %s\n",
				runs[0].ID, repoDir,
			),
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
			}
			if stdout.String() != tt.wantStdout {
				t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", stdout.String(), tt.wantStdout)
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("unexpected stderr:\n got: %q\nwant: %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunRunsListLimitBoundsNewestMatches(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	seeds := make([]runListSeed, 0, 25)
	for index := 0; index < 25; index++ {
		seeds = append(seeds, runListSeed{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    fmt.Sprintf("feature/run-%02d", index),
			prNumber:  fmt.Sprintf("%d", 500+index),
			createdAt: base.Add(time.Duration(index) * time.Minute),
		})
	}
	runs := seedRunsForList(t, homeDir, seeds)
	withRunsListNow(t, base.Add(40*time.Minute))
	tests := []struct {
		name       string
		args       []string
		wantLines  int
		wantFirst  string
		wantLast   string
		wantStderr string
	}{
		{
			name:       "default bounds to the 20 newest",
			args:       []string{"runs", "list"},
			wantLines:  20,
			wantFirst:  runs[24].ID,
			wantLast:   runs[5].ID,
			wantStderr: "(5 older Run(s) hidden; use --limit 0)\n",
		},
		{
			name:       "limit zero is unbounded",
			args:       []string{"runs", "list", "--limit", "0"},
			wantLines:  25,
			wantFirst:  runs[24].ID,
			wantLast:   runs[0].ID,
			wantStderr: "",
		},
		{
			name:       "explicit limit bounds the newest",
			args:       []string{"runs", "list", "--limit", "3"},
			wantLines:  3,
			wantFirst:  runs[24].ID,
			wantLast:   runs[22].ID,
			wantStderr: "(22 older Run(s) hidden; use --limit 0)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
			}
			lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			if len(lines) != tt.wantLines {
				t.Fatalf("expected %d lines, got %d:\n%s", tt.wantLines, len(lines), stdout.String())
			}
			if !strings.HasPrefix(lines[0], tt.wantFirst+"  ") {
				t.Fatalf("expected first line for Run %s, got %q", tt.wantFirst, lines[0])
			}
			if !strings.HasPrefix(lines[len(lines)-1], tt.wantLast+"  ") {
				t.Fatalf("expected last line for Run %s, got %q", tt.wantLast, lines[len(lines)-1])
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("unexpected stderr:\n got: %q\nwant: %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunRunsListAllRowsHiddenKeepsSingleEmptyLine(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:        store.KindResolve,
			state:       store.StateClean,
			gitRoot:     repoDir,
			branch:      "feature/only-terminal",
			prNumber:    "321",
			createdAt:   base,
			completedAt: base.Add(time.Minute),
		},
	})
	withRunsListNow(t, base.Add(10*time.Minute))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.String() != "No Runs found.\n" {
		t.Fatalf("expected single empty-result line, got %q", stdout.String())
	}
	if stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
		t.Fatalf("expected terminal-hidden note, got %q", stderr.String())
	}
}

func TestRunRunsWithoutSubcommandHonorsInteractivity(t *testing.T) {
	t.Run("non-interactive exits 2 naming runs list", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, false)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "runs requires a subcommand in non-interactive mode; use 'roundfix runs list'") {
			t.Fatalf("expected non-interactive runs guidance, got %q", stderr.String())
		}
	})
	t.Run("interactive stdin without a TTY exits 2 naming runs list", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "runs requires a subcommand in non-interactive mode; use 'roundfix runs list'") {
			t.Fatalf("expected non-interactive runs guidance, got %q", stderr.String())
		}
	})
	t.Run("interactive terminal opens the Run Browser over the repository's Runs", func(t *testing.T) {
		homeDir, repoDir := withCLIWorkspace(t)
		base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
		runs := seedRunsForList(t, homeDir, []runListSeed{
			{
				kind:        store.KindResolve,
				state:       store.StateClean,
				gitRoot:     repoDir,
				branch:      "feature/terminal",
				prNumber:    "201",
				createdAt:   base,
				completedAt: base.Add(time.Minute),
			},
			{
				kind:      store.KindResolve,
				state:     store.StateActive,
				gitRoot:   repoDir,
				branch:    "feature/active",
				prNumber:  "202",
				createdAt: base.Add(2 * time.Minute),
			},
		})
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "always")
		sessionCalls := withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout after browser cancel, got %q", stdout.String())
		}
		calls := *sessionCalls
		if len(calls) != 1 {
			t.Fatalf("expected one Run Browser session, got %d", len(calls))
		}
		if len(calls[0].active) != 1 || calls[0].active[0].ID != runs[1].ID {
			t.Fatalf("expected the Active Run listed first, got %+v", calls[0].active)
		}
		if len(calls[0].all) != 2 || calls[0].all[0].ID != runs[1].ID || calls[0].all[1].ID != runs[0].ID {
			t.Fatalf("expected all Runs newest first, got %+v", calls[0].all)
		}
	})
	t.Run("interactive terminal without a Run Database opens the empty browser", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "always")
		sessionCalls := withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
		}
		calls := *sessionCalls
		if len(calls) != 1 {
			t.Fatalf("expected one Run Browser session, got %d", len(calls))
		}
		if len(calls[0].active) != 0 || len(calls[0].all) != 0 {
			t.Fatalf("expected empty listings without a Run Database, got %+v", calls[0])
		}
	})
}

func TestRunRunsListEmptyResultExitsZero(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected empty list exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.String() != "No Runs found.\n" {
		t.Fatalf("expected single empty-result line, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunRunsListUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "unknown subcommand", args: []string{"runs", "bogus"}, want: "Run 'roundfix runs --help' for usage."},
		{name: "unexpected list argument", args: []string{"runs", "list", "extra"}, want: "Run 'roundfix runs list --help' for usage."},
		{name: "unknown state value", args: []string{"runs", "list", "--state", "bogus"}, want: `unknown --state "bogus"; use active, terminal, or all`},
		{name: "negative limit", args: []string{"runs", "list", "--limit", "-1"}, want: "--limit must be 0 or a positive count, got -1"},
		{name: "removed active flag", args: []string{"runs", "list", "--active"}, want: "Run 'roundfix runs list --help' for usage."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("expected usage pointer %q, got %q", tt.want, stderr.String())
			}
		})
	}
}

func TestRunRunsListOutsideRepositoryRequiresAll(t *testing.T) {
	homeDir := t.TempDir()
	outsideRepo := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Chdir(outsideRepo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected exit code 2 outside repository, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--all") {
		t.Fatalf("expected --all guidance, got %q", stderr.String())
	}
}

func TestRunDetachRejectsInteractiveWithExistingConflictShape(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "resolve", args: []string{"resolve", "--detach", "--interactive"}},
		{name: "watch", args: []string{"watch", "--detach", "--interactive"}},
		{name: "implement", args: []string{"implement", "--detach", "--interactive"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "--interactive cannot be used with --no-input") {
				t.Fatalf("expected existing no-input conflict shape, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunSetupFreshMachineAcceptsOffers(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["claude-agent-acp"] = "/bin/claude-agent-acp"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: installed",
		"adapter: ok",
		"profile readiness: passed",
		"acpx agents override: installed",
		"User Config: installed",
		"Project Config: installed",
	})
	if got := strings.Join(fake.installCalls, "\n"); got != "npm install -g acpx@"+agent.MinimumACPXVersion {
		t.Fatalf("expected minimum acpx install command, got %q", got)
	}
	if got := strings.Join(fake.initScopes, ","); got != "user,project" {
		t.Fatalf("expected User and Project Config init flows, got %q", got)
	}
	if fake.acpxInitCalls != 0 {
		t.Fatalf("expected in-memory ACPX proposal without config init side effects, got %d init calls", fake.acpxInitCalls)
	}
	acpxConfig := fake.files[fake.acpxConfigPath]
	if strings.Contains(acpxConfig, `"codex"`) || !strings.Contains(acpxConfig, `"command": "claude-agent-acp"`) {
		t.Fatalf("expected only the non-Codex local adapter override, got %s", acpxConfig)
	}
	for _, want := range []string{"acpx agents override diff:", "--- ", "+++ ", "+  \"agents\""} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected setup output to contain %q, got %q", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected --yes to avoid prompts, got stderr %q", stderr.String())
	}
}

func TestRunSetupHealthyMachineIsIdempotent(t *testing.T) {
	assertSetupCommandHealthyMachineIsIdempotent(t)
}

func TestRunSetupNewerACPXIsReadyWithoutInstall(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxVersion = "0.12.1"
	fake.paths["codex-acp"] = "/bin/codex-acp"
	fake.files[fake.userConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.projectConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.acpxConfigPath] = "{\n  \"agents\": {\n    \"codex\": {\n      \"command\": \"codex-acp\"\n    }\n  }\n}\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "acpx: ok (0.12.1 >= "+agent.MinimumACPXVersion+")") {
		t.Fatalf("expected newer acpx readiness report, got %q", stdout.String())
	}
	if len(fake.installCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("expected newer acpx to avoid install and prompt, installs=%v prompts=%v", fake.installCalls, fake.prompts)
	}
}

func TestSetupCommandCompatibility(t *testing.T) {
	assertSetupCommandHealthyMachineIsIdempotent(t)
}

func assertSetupCommandHealthyMachineIsIdempotent(t *testing.T) {
	t.Helper()
	fake := newSetupFakeDeps()
	fake.paths["codex-acp"] = "/bin/codex-acp"
	fake.files[fake.userConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.projectConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.acpxConfigPath] = "{\n  \"agents\": {\n    \"codex\": {\n      \"command\": \"codex-acp\"\n    }\n  }\n}\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: ok",
		"adapter: ok",
		"profile readiness: passed",
		"acpx agents override: ok",
		"User Config: ok",
		"Project Config: ok",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("expected idempotent setup to avoid side effects, installs=%v init=%v writes=%v prompts=%v", fake.installCalls, fake.initScopes, fake.writeCalls, fake.prompts)
	}
	if len(fake.probeRequests) != 3 {
		t.Fatalf("expected three distinct exact profile proofs, got %#v", fake.probeRequests)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSetupReportsAdapterFailuresWithoutWrites(t *testing.T) {
	tests := []struct {
		name       string
		adapterErr error
		want       string
	}{
		{
			name: "legacy lineage",
			adapterErr: &agent.AdapterLineageError{
				Command: "codex-acp",
				Package: "@zed-industries/codex-acp",
				Version: "0.16.0",
			},
			want: "legacy package @zed-industries/codex-acp version 0.16.0",
		},
		{
			name: "unsupported official version",
			adapterErr: &agent.AdapterVersionError{
				Command:         "codex-acp",
				Package:         agent.CodexAdapterPackage,
				FoundVersion:    "1.1.3",
				RequiredVersion: agent.PinnedCodexAdapterVersion,
			},
			want: "package @agentclientprotocol/codex-acp version 1.1.3; required version 1.1.4 or newer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := newSetupFakeDeps()
			fake.files[fake.userConfigPath] = "user config sentinel\n"
			fake.files[fake.projectConfigPath] = "project config sentinel\n"
			fake.files[fake.acpxConfigPath] = "{}\n"
			fake.adapterErr = tt.adapterErr
			before := make(map[string]string, len(fake.files))
			for path, content := range fake.files {
				before[path] = content
			}
			withSetupFakeDeps(t, fake)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

			if code != exitRunFailed {
				t.Fatalf("expected adapter failure exit %d, got %d", exitRunFailed, code)
			}
			for _, want := range []string{
				"adapter: failed",
				"command \"codex-acp\"",
				tt.want,
				agent.CodexAdapterInstallCommand(),
			} {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("expected setup output to contain %q, got %q", want, stdout.String())
				}
			}
			if strings.Contains(stdout.String(), "acpx agents override:") || strings.Contains(stdout.String(), "User Config:") || strings.Contains(stdout.String(), "Project Config:") {
				t.Fatalf("adapter failure continued into mutation surfaces: %q", stdout.String())
			}
			if !reflect.DeepEqual(fake.files, before) || len(fake.writeCalls) != 0 || len(fake.initScopes) != 0 || fake.acpxInitCalls != 0 {
				t.Fatalf("adapter failure changed config state: before=%v after=%v writes=%v init=%v acpxInit=%d", before, fake.files, fake.writeCalls, fake.initScopes, fake.acpxInitCalls)
			}
			if len(fake.probeRequests) != 0 {
				t.Fatalf("adapter failure started profile proof: %#v", fake.probeRequests)
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

func TestRunSetupProfileProofsEveryDistinctTupleOnceBeforePersistence(t *testing.T) {
	fake := newSetupFakeDeps()
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	want := map[roundconfig.AgentSelection]int{
		{Runtime: "codex", Model: "gpt-5.6-sol", ReasoningEffort: "high"}:     1,
		{Runtime: "codex", Model: "gpt-5.5", ReasoningEffort: "xhigh"}:        1,
		{Runtime: "claude", Model: "claude-opus-5", ReasoningEffort: "xhigh"}: 1,
	}
	got := map[roundconfig.AgentSelection]int{}
	for _, request := range fake.probeRequests {
		got[roundconfig.AgentSelection{
			Runtime:         strings.TrimSuffix(request.Runtime.ID, "-custom"),
			Model:           request.Runtime.Model,
			ReasoningEffort: request.Runtime.ReasoningEffort,
		}]++
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("profile proof requests = %#v, want %#v", got, want)
	}
	if len(fake.initScopes) != 2 {
		t.Fatalf("expected proof to precede both config writes, init scopes=%v", fake.initScopes)
	}
}

func TestRunSetupProfileProofFailurePreservesAllTargets(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.files[fake.acpxConfigPath] = "{\n  \"theme\": \"sentinel\"\n}\n"
	fake.probeErr = &agent.SelectionUnsupportedError{
		Kind:                agent.SelectionReasoningControlNotAdvertised,
		Runtime:             "codex",
		Model:               "gpt-5.6-sol",
		ReasoningEffort:     "high",
		AdvertisedModels:    []string{"gpt-5.6-sol"},
		AdvertisedReasoning: []string{"medium"},
	}
	before := cloneSetupFiles(fake.files)
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitRunFailed, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "profile readiness: failed") {
		t.Fatalf("expected profile readiness failure, got %q", stdout.String())
	}
	if !reflect.DeepEqual(fake.files, before) || len(fake.writeCalls) != 0 || len(fake.initScopes) != 0 {
		t.Fatalf("proof failure mutated targets: before=%v after=%v writes=%v init=%v", before, fake.files, fake.writeCalls, fake.initScopes)
	}
}

func TestRunSetupProfileCleanupAndInvalidEvidencePreserveAllTargets(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "cleanup failure",
			err:  &agent.AgentSessionCleanupError{Session: "roundfix-preflight-test", Err: errors.New("close denied")},
			want: agent.SessionCleanupFailed,
		},
		{
			name: "invalid evidence",
			err:  &agent.CapabilityEvidenceError{Issues: []string{agent.CapabilityIssueInvalidOption}},
			want: agent.CapabilityEvidenceInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := newSetupFakeDeps()
			fake.probeErr = tt.err
			before := cloneSetupFiles(fake.files)
			withSetupFakeDeps(t, fake)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

			if code != exitRunFailed {
				t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitRunFailed, stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), "profile readiness: failed") || !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("expected classified profile failure %q, got %q", tt.want, stdout.String())
			}
			if !reflect.DeepEqual(fake.files, before) || len(fake.writeCalls) != 0 {
				t.Fatalf("profile failure mutated targets: before=%v after=%v writes=%v", before, fake.files, fake.writeCalls)
			}
		})
	}
}

func TestRunSetupProfileWriteFailurePreservesNotYetCommittedTargets(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.writeErrors = map[string]error{fake.userConfigPath: errors.New("disk full")}
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitRunFailed, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "User Config: failed") || !strings.Contains(stdout.String(), "disk full") {
		t.Fatalf("expected User Config write failure, got %q", stdout.String())
	}
	if _, ok := fake.files[fake.userConfigPath]; ok {
		t.Fatalf("failed User Config write created target: %v", fake.files)
	}
	if _, ok := fake.files[fake.projectConfigPath]; ok {
		t.Fatalf("later Project Config target changed after User Config failure: %v", fake.files)
	}
}

func TestRunSetupProfilePersistenceMatchesSubsequentValidation(t *testing.T) {
	fake := newSetupFakeDeps()
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	resolved, err := roundconfig.ResolveConfigProposal(
		[]byte(fake.files[fake.userConfigPath]),
		[]byte(fake.files[fake.projectConfigPath]),
	)
	if err != nil {
		t.Fatalf("resolve persisted config: %v", err)
	}
	validationRunner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{
				Runtime: req.Runtime.ID, Model: req.Runtime.Model, ReasoningEffort: req.Runtime.ReasoningEffort,
				Assignment: agent.SelectionAssignment{Encoding: agent.SelectionEncodingIndependent},
				Adapter:    fake.adapterEvidence,
				Status:     agent.SelectionProofStatusProven,
			}, nil
		},
	}
	result := proveProfileSelections(context.Background(), resolved, roundconfig.RequiredWorkCategories(), fake.gitRoot, validationRunner)
	if result.Err != nil {
		t.Fatalf("subsequent profile validation: %v", result.Err)
	}
	if len(result.Proofs) != 3 || len(validationRunner.exactRequests) != 3 {
		t.Fatalf("subsequent validation proofs=%d requests=%d, want three distinct tuples", len(result.Proofs), len(validationRunner.exactRequests))
	}
}

func TestRunSetupNoInputProfileProofCreatesNoTargets(t *testing.T) {
	fake := newSetupFakeDeps()
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	if len(fake.probeRequests) != 3 || len(fake.files) != 0 || len(fake.writeCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("--no-input proof/mutation mismatch: proofs=%d files=%v writes=%v prompts=%v", len(fake.probeRequests), fake.files, fake.writeCalls, fake.prompts)
	}
}

func TestRunSetupAdapterMigrationDeclinePreservesAllTargets(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.files[fake.acpxConfigPath] = "{\n  \"agents\": {\n    \"codex\": {\n      \"command\": \"codex-acp\"\n    }\n  }\n}\n"
	fake.files[fake.userConfigPath] = "# user config sentinel\n{}\n"
	fake.files[fake.projectConfigPath] = "# project config sentinel\n{}\n"
	fake.adapterErr = &agent.AdapterLineageError{
		Command: "codex-acp",
		Package: "@zed-industries/codex-acp",
		Version: "0.16.0",
	}
	fake.confirm = func(context.Context, io.Writer, string) (bool, error) { return false, nil }
	before := cloneSetupFiles(fake.files)
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("declined migration exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	if len(fake.prompts) != 1 || !strings.Contains(fake.prompts[0], "Migrate") {
		t.Fatalf("migration prompts = %v, want one explicit migration offer", fake.prompts)
	}
	if !reflect.DeepEqual(fake.files, before) || len(fake.writeCalls) != 0 || len(fake.initScopes) != 0 {
		t.Fatalf("declined migration mutated targets: before=%v after=%v writes=%v init=%v", before, fake.files, fake.writeCalls, fake.initScopes)
	}
}

func TestRunSetupAdapterMigrationPersistsSupportedCommand(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.files[fake.acpxConfigPath] = "{\n  \"agents\": {\n    \"codex\": {\n      \"command\": \"codex-acp\"\n    }\n  }\n}\n"
	fake.adapterErr = &agent.AdapterLineageError{
		Command: "codex-acp",
		Package: "@zed-industries/codex-acp",
		Version: "0.16.0",
	}
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("setup exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	content := fake.files[fake.acpxConfigPath]
	for _, want := range []string{`"command": "npx"`, `"args": ["-y","` + agent.CodexAdapterPackage + `@` + agent.PinnedCodexAdapterVersion + `"]`} {
		if !strings.Contains(content, want) {
			t.Fatalf("migrated ACPX config missing %q:\n%s", want, content)
		}
	}
	if strings.Contains(content, `"command": "codex-acp"`) {
		t.Fatalf("migrated ACPX config retained bare PATH override:\n%s", content)
	}
	if len(fake.adapterRequests) != 2 || fake.adapterRequests[1].Command != agent.CodexAdapterCommand() {
		t.Fatalf("adapter proof requests = %#v, want current then deterministic proposal", fake.adapterRequests)
	}
}

func TestRunSetupProfileProofUsesProposedProfilesAndWorkDir(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.config.Runtimes.Codex.Model = "repo-codex"
	fake.config.Runtimes.Codex.ReasoningEffort = "repo-xhigh"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stdout %q stderr %q)", code, stdout.String(), stderr.String())
	}
	if len(fake.probeRequests) != 3 {
		t.Fatalf("expected three generated profile proofs, got %#v", fake.probeRequests)
	}
	for _, gotProbe := range fake.probeRequests {
		if gotProbe.WorkDir != fake.gitRoot {
			t.Fatalf("expected Setup proof in repository workdir, got %#v", gotProbe)
		}
	}
}

func TestRunSetupAcceptsConfiguredEmptyReasoningEffort(t *testing.T) {
	fake := newSetupFakeDeps()
	customConfig := strings.ReplaceAll(roundconfig.DefaultConfigYAML(), "reasoning_effort: high", `reasoning_effort: ""`)
	fake.files[fake.userConfigPath] = customConfig
	fake.files[fake.projectConfigPath] = customConfig
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stdout %q stderr %q)", code, stdout.String(), stderr.String())
	}
	found := false
	for _, gotProbe := range fake.probeRequests {
		if gotProbe.Runtime.ID == "codex" && gotProbe.Runtime.Model == "gpt-5.6-sol" && gotProbe.Runtime.ReasoningEffort == "" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected exact model-managed Agent Selection proof, got %#v", fake.probeRequests)
	}
}

func TestRunSetupRejectsMissingConfiguredAgentSelection(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.config.Runtimes.Codex.Model = ""
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected setup selection failure exit %d, got %d", exitRunFailed, code)
	}
	if !strings.Contains(stdout.String(), `adapter: failed (agent selection for runtime "codex" missing model`) {
		t.Fatalf("expected missing selection diagnostic, got %q", stdout.String())
	}
	if len(fake.probeRequests) != 0 {
		t.Fatalf("expected invalid selection to skip readiness probe, got %#v", fake.probeRequests)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSetupMismatchedACPXUpgradeOffer(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxVersion = "0.11.0"
	fake.files[fake.userConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.projectConfigPath] = roundconfig.DefaultConfigYAML()
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "acpx: installed") || !strings.Contains(stdout.String(), "found 0.11.0") {
		t.Fatalf("expected mismatched acpx upgrade report, got %q", stdout.String())
	}
	if got := strings.Join(fake.installCalls, "\n"); got != "npm install -g acpx@"+agent.MinimumACPXVersion {
		t.Fatalf("expected minimum acpx install command, got %q", got)
	}
}

func TestRunSetupMergesACPXAgentsOverridePreservingUnrelatedBytes(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.paths["claude-agent-acp"] = "/bin/claude-agent-acp"
	fake.files[fake.userConfigPath] = roundconfig.DefaultConfigYAML()
	fake.files[fake.projectConfigPath] = roundconfig.DefaultConfigYAML()
	unrelated := "  \"theme\": {\n    \"color\": \"blue\",\n    \"nested\": [1, 2, 3]\n  }"
	existingAgent := "    \"custom\": {\n      \"command\": \"existing-custom\"\n    }"
	fake.files[fake.acpxConfigPath] = "{\n" + unrelated + ",\n  \"agents\": {\n" + existingAgent + "\n  }\n}\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	acpxConfig := fake.files[fake.acpxConfigPath]
	for _, want := range []string{unrelated, existingAgent, `"claude"`, `"command": "claude-agent-acp"`} {
		if !strings.Contains(acpxConfig, want) {
			t.Fatalf("expected merged config to preserve/include %q, got %s", want, acpxConfig)
		}
	}
	if !strings.Contains(stdout.String(), "acpx agents override diff:") || !strings.Contains(stdout.String(), "+    \"claude\"") {
		t.Fatalf("expected before/after diff for override, got %q", stdout.String())
	}
}

func TestRunSetupDeclinedOffersReportNoWrites(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["claude-agent-acp"] = "/bin/claude-agent-acp"
	fake.confirm = func(context.Context, io.Writer, string) (bool, error) {
		return false, nil
	}
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected declined setup offers to exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: offered: declined",
		"adapter: skipped",
		"profile readiness: skipped",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || fake.acpxInitCalls != 0 {
		t.Fatalf("expected no writes after declined offers, installs=%v init=%v writes=%v acpxInit=%d", fake.installCalls, fake.initScopes, fake.writeCalls, fake.acpxInitCalls)
	}
	if len(fake.prompts) != 1 {
		t.Fatalf("expected one prerequisite confirmation prompt, got %d: %v", len(fake.prompts), fake.prompts)
	}
}

func TestRunSetupNoInputSkipsOffers(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["claude-agent-acp"] = "/bin/claude-agent-acp"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected no-input skipped offers to exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: skipped",
		"adapter: skipped",
		"profile readiness: skipped",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("expected --no-input to avoid side effects, installs=%v init=%v writes=%v prompts=%v", fake.installCalls, fake.initScopes, fake.writeCalls, fake.prompts)
	}
}

func TestRunSetupExitCodes(t *testing.T) {
	t.Run("check failure exits one", func(t *testing.T) {
		fake := newSetupFakeDeps()
		fake.nodeVersion = "v20.0.0"
		fake.files[fake.userConfigPath] = "defaults:\n  agent: codex\n"
		fake.files[fake.projectConfigPath] = "defaults:\n  verification: make verify\n"
		withSetupFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"setup"}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("expected setup check failure exit 1, got %d", code)
		}
		if !strings.Contains(stdout.String(), "node: failed") {
			t.Fatalf("expected node failure report, got %q", stdout.String())
		}
	})

	t.Run("usage error exits two", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"setup", "--yes", "--no-input"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected setup usage error exit 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "--yes cannot be used with --no-input") {
			t.Fatalf("expected usage diagnostic, got %q", stderr.String())
		}
	})
}

func TestHealthCheckerReturnsStructuredReadOnlyResults(t *testing.T) {
	ctx := context.Background()
	req := agent.ProbeRequest{
		Runtime: agent.RuntimeSpec{ID: "codex", Model: "gpt-5.5", ReasoningEffort: "xhigh"},
		WorkDir: "/repo/project",
	}
	probed := false
	checker := newHealthChecker(healthCheckDependencies{
		nodeVersion: func(context.Context) (string, error) {
			return " v25.6.1 \n", nil
		},
		acpxVersion: func(context.Context) (string, error) {
			return agent.MinimumACPXVersion + "\n", nil
		},
		checkAdapter: func(context.Context, agent.RuntimeSpec) (agent.AdapterEvidence, error) {
			return agent.AdapterEvidence{
				Command: "npx -y " + agent.CodexAdapterPackage,
				Package: agent.CodexAdapterPackage,
				Version: agent.PinnedCodexAdapterVersion,
			}, nil
		},
		probeAgent: func(_ context.Context, got agent.ProbeRequest) error {
			probed = true
			if !reflect.DeepEqual(got, req) {
				t.Fatalf("expected probe request %#v, got %#v", req, got)
			}
			return nil
		},
	})

	assertCheckResult(t, checker.Node(ctx), CheckResult{
		Name:   HealthCheckNode,
		Status: CheckStatusOK,
		Detail: "v25.6.1 >= " + setupNodeMinimumVersion,
	})
	assertCheckResult(t, checker.ACPX(ctx), CheckResult{
		Name:   HealthCheckACPX,
		Status: CheckStatusOK,
		Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion,
	})
	assertCheckResult(t, checker.Adapter(ctx, req.Runtime), CheckResult{
		Name:   HealthCheckAdapter,
		Status: CheckStatusOK,
		Detail: "command=\"npx -y " + agent.CodexAdapterPackage + "\"; package=" + agent.CodexAdapterPackage + "; version=" + agent.PinnedCodexAdapterVersion,
	})
	assertCheckResult(t, checker.Agent(ctx, req), CheckResult{
		Name:   HealthCheckAgent,
		Status: CheckStatusOK,
		Detail: "codex",
	})
	if !probed {
		t.Fatal("expected agent probe to run")
	}
}

func TestHealthCheckerAcceptsNewerACPXVersion(t *testing.T) {
	checker := newHealthChecker(healthCheckDependencies{
		acpxVersion: func(context.Context) (string, error) {
			return "0.12.1\n", nil
		},
	})

	assertCheckResult(t, checker.ACPX(context.Background()), CheckResult{
		Name:   HealthCheckACPX,
		Status: CheckStatusOK,
		Detail: "0.12.1 >= " + agent.MinimumACPXVersion,
	})
}

func TestHealthCheckerReportsFailedPrerequisitesWithNextActions(t *testing.T) {
	ctx := context.Background()
	probeErr := errors.New("probe denied")
	checker := newHealthChecker(healthCheckDependencies{
		nodeVersion: func(context.Context) (string, error) {
			return "v20.0.0", nil
		},
		acpxVersion: func(context.Context) (string, error) {
			return "0.11.0", nil
		},
		probeAgent: func(context.Context, agent.ProbeRequest) error {
			return probeErr
		},
	})

	assertCheckResult(t, checker.Node(ctx), CheckResult{
		Name:       HealthCheckNode,
		Status:     CheckStatusFailed,
		Detail:     "found v20.0.0; Node.js " + setupNodeMinimumVersion + " or newer is required",
		NextAction: "install Node.js " + setupNodeMinimumVersion + " or newer",
	})
	assertCheckResult(t, checker.ACPX(ctx), CheckResult{
		Name:       HealthCheckACPX,
		Status:     CheckStatusFailed,
		Detail:     "found 0.11.0; acpx " + agent.MinimumACPXVersion + " or newer is required; run " + setupACPXInstallCommand(),
		NextAction: setupACPXInstallCommand(),
	})
	assertCheckResult(t, checker.Agent(ctx, agent.ProbeRequest{Runtime: agent.RuntimeSpec{ID: "codex"}}), CheckResult{
		Name:   HealthCheckAgent,
		Status: CheckStatusFailed,
		Detail: "probe denied",
		Err:    probeErr,
	})
}

func TestHealthCheckerReportsCodexResult(t *testing.T) {
	checker := newHealthChecker(healthCheckDependencies{
		codexInspector: fakeCodexInspector{
			result: codex.Result{
				Status:     codex.StatusFailed,
				Detail:     "/path/codex has an invalid or missing code signature",
				NextAction: codex.ReinstallNextAction,
			},
		},
	})

	assertCheckResult(t, checker.Codex(context.Background()), CheckResult{
		Name:       HealthCheckCodex,
		Status:     CheckStatusFailed,
		Detail:     "/path/codex has an invalid or missing code signature",
		NextAction: codex.ReinstallNextAction,
	})
}

func assertCheckResult(t *testing.T, got CheckResult, want CheckResult) {
	t.Helper()
	if got != want {
		t.Fatalf("unexpected check result:\n got: %#v\nwant: %#v", got, want)
	}
}

type fakeCodexInspector struct {
	result codex.Result
}

func (fake fakeCodexInspector) Inspect(context.Context) codex.Result {
	return fake.result
}

func TestRunSkillsCheck(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "check"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Roundfix skill check passed") {
		t.Fatalf("expected skills check output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix") || strings.Contains(stdout.String(), "roundfix-watch") || strings.Contains(stdout.String(), "roundfix-resolve-round") {
		t.Fatalf("expected consolidated skill name only, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsListSeparatesOwnedFromExternal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "list"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	out := stdout.String()
	for _, want := range []string{
		"Bundled skills",
		"roundfix",
		"write-prd",
		"evidence-gate",
		"Recommended skills",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected skills list to contain %q, got %q", want, out)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCopiesArtifactsToProjectByDefault(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Installed Roundfix skill for project") {
		t.Fatalf("expected project install output, got %q", stdout.String())
	}
	targetDir := filepath.Join(repoDir, ".agents", "skills")
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
	} {
		if _, err := os.Stat(filepath.Join(targetDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCopiesArtifactsToExplicitTarget(t *testing.T) {
	targetDir := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "codex", "--dir", targetDir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Installed Roundfix skill for codex") {
		t.Fatalf("expected install output, got %q", stdout.String())
	}
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
	} {
		if _, err := os.Stat(filepath.Join(targetDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCreatesClaudeProjectSymlink(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustMkdir(t, filepath.Join(repoDir, ".claude", "skills"))
	prompted := false
	withClaudeSkillSymlinkPrompt(t, func(ctx context.Context, stderr io.Writer, linkPath string, target string) (bool, error) {
		prompted = true
		expectedLink := filepath.Join(repoDir, ".claude", "skills", "roundfix")
		if linkPath != expectedLink {
			t.Fatalf("expected prompt link path %q, got %q", expectedLink, linkPath)
		}
		expectedTarget := filepath.Join("..", "..", ".agents", "skills", "roundfix")
		if target != expectedTarget {
			t.Fatalf("expected prompt target %q, got %q", expectedTarget, target)
		}
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !prompted {
		t.Fatal("expected Claude symlink prompt")
	}
	linkPath := filepath.Join(repoDir, ".claude", "skills", "roundfix")
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("expected Claude skill symlink: %v", err)
	}
	expectedTarget := filepath.Join("..", "..", ".agents", "skills", "roundfix")
	if target != expectedTarget {
		t.Fatalf("expected symlink target %q, got %q", expectedTarget, target)
	}
	if !strings.Contains(stdout.String(), "Created Claude skill symlink") {
		t.Fatalf("expected symlink output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallRejectsUnsupportedTarget(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "other"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported skill install target") {
		t.Fatalf("expected unsupported target error, got %q", stderr.String())
	}
}

func TestRunOperationalCommandAcceptsMVPFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedCode   int
		expectedOutput string
	}{
		{
			name:           "fetch",
			args:           []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"},
			expectedCode:   0,
			expectedOutput: "reached Fetched",
		},
		{
			name:           "resolve",
			args:           []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
			expectedCode:   0,
			expectedOutput: "resolve selected 1 downloaded Unresolved Review Issue",
		},
		{
			name:           "watch",
			args:           []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			expectedCode:   0,
			expectedOutput: "Watch Run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != tt.expectedCode {
				t.Fatalf("expected exit code %d, got %d", tt.expectedCode, code)
			}
			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Fatalf("expected output to contain %q, got stdout=%q stderr=%q", tt.expectedOutput, stdout.String(), stderr.String())
			}
			if !strings.Contains(output, builtinArtifactDirForRepo(t, repoDir)) {
				t.Fatalf("expected artifact dir in output, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if strings.Contains(output, "not implemented yet") {
				t.Fatalf("did not expect scaffold message, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if tt.name == "fetch" {
				if !strings.Contains(stdout.String(), "did not start an Agent, commit, or push") {
					t.Fatalf("expected fetch safety confirmation, got %q", stdout.String())
				}
				if !strings.Contains(stdout.String(), "Review Issues: 1") {
					t.Fatalf("expected fetched issue count, got %q", stdout.String())
				}
				if !strings.Contains(stdout.String(), "Artifacts: created new Round") {
					t.Fatalf("expected new artifact confirmation, got %q", stdout.String())
				}
				issuePath := filepath.Join(defaultReviewRootForRepo(repoDir, "123"), "round-001", "issue_001.md")
				issueContent, err := os.ReadFile(issuePath)
				if err != nil {
					t.Fatalf("expected Review Issue artifact %s: %v", issuePath, err)
				}
				if !strings.Contains(string(issueContent), "source_ref: thread:PRRT_test,comment:PRRC_test") {
					t.Fatalf("expected source ref in issue artifact, got %s", string(issueContent))
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
			} else if tt.name == "resolve" {
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Verification passed (attempt 1).") {
					t.Fatalf("expected verification verdict summary, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertNoAgentLogs(t, repoDir)
			} else if tt.name == "watch" {
				if !strings.Contains(stderr.String(), "Review Source status: settled") {
					t.Fatalf("expected fake Review Source status output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Fetched Round 001 with 1 Review Issue") {
					t.Fatalf("expected watch fetch output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Verification passed (attempt 1).") {
					t.Fatalf("expected verification verdict summary, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "reached Clean") {
					t.Fatalf("expected watch Clean terminal outcome, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertNoAgentLogs(t, repoDir)
			}
		})
	}
}

func TestRunFetchWritesReviewArtifactsUnderSpecSelector(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	specSlug := "0001-widget-flow"
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", specSlug))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--spec", specSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", specSlug, "reviews", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected spec-associated Review Issue artifact %s: %v", issuePath, err)
	}
	oldPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "reviews", "pr-123", "round-001", "issue_001.md")
	if _, err := os.Stat(oldPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no default Artifact Directory Review Issue at %s, got err %v", oldPath, err)
	}
}

func TestRunFetchWritesReviewArtifactsUnderTrailerSpec(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	specSlug := "0001-widget-flow"
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", specSlug))
	withReviewSpecGitRunner(t, &fakeReviewSpecGitRunner{message: "subject\n\nRoundfix-Spec: " + specSlug})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", specSlug, "reviews", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected trailer-associated Review Issue artifact %s: %v", issuePath, err)
	}
}

func TestRunFetchWritesReviewArtifactsUnderSpeclessRoot(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-123", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected spec-less Review Issue artifact %s: %v", issuePath, err)
	}
}

func TestReviewSpecSlugAssociation(t *testing.T) {
	repoRoot := t.TempDir()
	explicitSlug := "0001-explicit"
	oldSlug := "0002-old"
	newSlug := "0003-new"
	for _, slug := range []string{explicitSlug, oldSlug, newSlug} {
		mustMkdir(t, filepath.Join(repoRoot, "docs", "specs", slug))
	}

	tests := []struct {
		name     string
		req      reviewSpecRequest
		message  string
		want     string
		wantCall bool
	}{
		{
			name: "explicit selector wins without reading git",
			req: reviewSpecRequest{
				ExplicitSlug: explicitSlug,
				RepoRoot:     repoRoot,
				HeadSHA:      "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: " + newSlug,
			want:     explicitSlug,
			wantCall: false,
		},
		{
			name: "newest trailer wins",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: " + oldSlug + "\nRoundfix-Spec: " + newSlug,
			want:     newSlug,
			wantCall: true,
		},
		{
			name: "unknown trailer slug is no association",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: 9999-missing",
			want:     "",
			wantCall: true,
		},
		{
			name: "missing trailer is no association",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject only",
			want:     "",
			wantCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeReviewSpecGitRunner{message: tt.message}

			got, err := reviewSpecSlug(context.Background(), tt.req, runner)
			if err != nil {
				t.Fatalf("reviewSpecSlug() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("reviewSpecSlug() = %q, want %q", got, tt.want)
			}
			if gotCall := runner.calls > 0; gotCall != tt.wantCall {
				t.Fatalf("git calls = %v, want %v", gotCall, tt.wantCall)
			}
		})
	}
}

func TestReviewArtifactUsesDefaultSpecsRootMatchesConfigString(t *testing.T) {
	if !reviewArtifactUsesDefaultSpecsRoot("docs/specs") {
		t.Fatal("expected built-in default Spec Root to match")
	}
	if !reviewArtifactUsesDefaultSpecsRoot("  docs/specs  ") {
		t.Fatal("expected whitespace-trimmed built-in default Spec Root to match")
	}
	if reviewArtifactUsesDefaultSpecsRoot(`docs\specs`) {
		t.Fatal("expected OS-specific path spelling to be treated as non-default config")
	}
}

type fakeReviewSpecGitRunner struct {
	message string
	calls   int
}

func (runner *fakeReviewSpecGitRunner) RunGit(_ context.Context, _ string, _ ...string) (string, error) {
	runner.calls++
	return runner.message, nil
}

func TestRunFetchWarnsAndIgnoresDeprecatedUserConfig(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
resolve:
  concurrent: 1
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch exit 0, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), "deprecated") {
		t.Fatalf("expected warning to stay off stdout, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Review Issues: 1") || !strings.Contains(stdout.String(), "Artifacts: created new Round") {
		t.Fatalf("expected normal fetch stdout, got %q", stdout.String())
	}
	warning := "config: resolve.concurrent is deprecated and ignored; use worktree.concurrency\n"
	if strings.Count(stderr.String(), warning) != 1 {
		t.Fatalf("expected one deprecated config warning %q, got %q", warning, stderr.String())
	}
	if strings.Contains(stderr.String(), "Preflight failed") {
		t.Fatalf("expected fetch to proceed past Preflight, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunWatchPrintsDeterministicStdoutReport(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, repoDir string)
		args       []string
		wantCode   int
		wantStdout string
	}{
		{
			name: "clean",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
				"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n",
		},
		{
			name: "unresolved",
			setup: func(t *testing.T, _ string) {
				t.Helper()
				withAgentRunner(t, &fakeAgentRunner{status: rounds.StatusFailed})
			},
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantCode: exitRunFailed,
			wantStdout: "" +
				"issue 001 failed — major: handle test issue\n" +
				"This Run (Unresolved after 1 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.\n" +
				"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.\n",
		},
		{
			name: "max rounds reached",
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
				withAgentRunner(t, &fakeAgentRunner{statuses: []string{rounds.StatusResolved, rounds.StatusFailed}})
				withFetchReviewItems(t, []reviewsource.ReviewItem{
					{
						Title:                   "major: handle first issue",
						File:                    "internal/first.go",
						Line:                    12,
						Severity:                "major",
						Author:                  "coderabbitai[bot]",
						Body:                    "First issue.",
						SourceRef:               "thread:PRRT_first,comment:PRRC_first",
						ReviewHash:              "review-hash-first",
						SourceReviewID:          "9001",
						SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
					},
					{
						Title:                   "major: handle second issue",
						File:                    "internal/second.go",
						Line:                    24,
						Severity:                "major",
						Author:                  "coderabbitai[bot]",
						Body:                    "Second issue.",
						SourceRef:               "thread:PRRT_second,comment:PRRC_second",
						ReviewHash:              "review-hash-second",
						SourceReviewID:          "9002",
						SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
					},
				})
			},
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "1", "--no-input"},
			wantStdout: "" +
				"issue 001 resolved — major: handle first issue\n" +
				"issue 002 failed — major: handle second issue\n" +
				"This Run (MaxRoundsReached after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.\n" +
				"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.\n",
		},
		{
			name: "stopped after fetch",
			setup: func(t *testing.T, _ string) {
				t.Helper()
				withAgentRunner(t, &fakeStoppingAgentRunner{})
				withChangedPaths(t, nil)
			},
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantStdout: "" +
				"issue 001 unresolved — major: handle test issue\n" +
				"This Run (Stopped after 1 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 1 unresolved.\n" +
				"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 0 failed, 1 unresolved.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.setup != nil {
				tt.setup(t, repoDir)
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q", tt.wantCode, code, stderr.String())
			}
			if stdout.String() != tt.wantStdout {
				t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", tt.wantStdout, stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunResolvePrintsDeterministicStdoutReport(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q", code, stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", wantStdout, stdout.String(), stderr.String())
	}
}

func TestReviewProfilePreflightResolveAndWatchUseOnlyReviewProfile(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "1", "--no-input"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
profiles:
  review:
    preferred:
      runtime: claude
      model: review-preferred
      reasoning_effort: medium
    fallbacks:
      - runtime: codex
        model: review-fallback
        reasoning_effort: high
`)
			runner := &fakeAgentRunner{}
			withAgentRunner(t, runner)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("%s exit = %d stderr=%q stdout=%q", tt.name, code, stderr.String(), stdout.String())
			}
			wantModels := []string{"review-preferred", "review-fallback"}
			if got := probeRequestModels(runner.probeRequests); !reflect.DeepEqual(got, wantModels) {
				t.Fatalf("%s probe models = %v, want %v", tt.name, got, wantModels)
			}
			if len(runner.runRuntimes) == 0 || runner.runRuntimes[0].Model != "review-preferred" {
				t.Fatalf("%s run runtime = %#v, want review-preferred", tt.name, runner.runRuntimes)
			}
		})
	}
}

func TestReviewProfilePreflightFetchCreatesNoAgentSession(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("fetch exit = %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(runner.probeRequests) != 0 || runner.calls != 0 {
		t.Fatalf("fetch must not prove or invoke Agent sessions, probes=%#v calls=%d", runner.probeRequests, runner.calls)
	}
}

func TestPrintReviewIssueReportSplitsRunAndCumulativeCountsAndReasons(t *testing.T) {
	report := reviewIssueReport{
		runIssues: []rounds.Issue{
			{Title: "major: generated file", Status: rounds.StatusInvalid, TerminalReason: "invalid: generated file"},
			{Title: "major: failing verification", Status: rounds.StatusFailed, TerminalReason: "Verification failed:\nmake verify exit status 1"},
			{Title: "major: still open", Status: rounds.StatusPending},
		},
		cumulativeIssues: []rounds.Issue{
			{Status: rounds.StatusResolved},
			{Status: rounds.StatusInvalid},
			{Status: rounds.StatusFailed},
			{Status: rounds.StatusPending},
		},
	}
	var stdout bytes.Buffer

	printReviewIssueReport(&stdout, store.StateCleanUnverified, 2, report)

	wantStdout := "" +
		"issue 001 invalid — major: generated file — reason: invalid: generated file\n" +
		"issue 002 failed — major: failing verification — reason: Verification failed: make verify exit status 1\n" +
		"issue 003 unresolved — major: still open\n" +
		"This Run (Clean Unverified after 2 Round(s)): 0 resolved, 1 invalid, 0 duplicated, 1 failed, 1 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 1 invalid, 0 duplicated, 1 failed, 1 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q", wantStdout, stdout.String())
	}
}

func TestReviewIssueReportDataMarksCumulativeUnavailableOnLoadFailure(t *testing.T) {
	req := commandRequest{reviewRoot: "["}
	preflightResult := preflight.Result{
		PullRequest: preflight.PullRequest{
			Number:         "123",
			HeadRepository: "owner/project",
			HeadBranch:     "feature/review",
		},
	}
	var stderr bytes.Buffer

	report := reviewIssueReportData(context.Background(), req, preflightResult, []rounds.Issue{{
		Title:  "major: handled issue",
		Status: rounds.StatusResolved,
	}}, &stderr)

	if !report.cumulativeUnavailable {
		t.Fatalf("expected cumulative data unavailable, got %+v", report)
	}
	var stdout bytes.Buffer
	printReviewIssueReport(&stdout, store.StateClean, 1, report)
	wantStdout := "" +
		"issue 001 resolved — major: handled issue\n" +
		"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: unavailable.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Pull Request cumulative Review Issue report unavailable") {
		t.Fatalf("expected load diagnostic on stderr, got %q", stderr.String())
	}
}

func TestRunResolvePersistsEffectiveSelection(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"resolve",
		"--pr", "123",
		"--agent", "codex",
		"--model", "stored-resolve-model",
		"--reasoning-effort", "stored-resolve-reasoning",
		"--round", "all",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q", code, stderr.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "stored-resolve-model" || run.ReasoningEffort != "stored-resolve-reasoning" {
		t.Fatalf("expected stored resolve selection, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: stored-resolve-model") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: stored-resolve-reasoning") {
		t.Fatalf("expected resolve progress to show concrete selection, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "auto") {
		t.Fatalf("expected no auto selection placeholder, got %q", stderr.String())
	}
}

func TestRunResolveAcceptsExplicitEmptyReasoningEffort(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"resolve",
		"--pr", "123",
		"--agent", "codex",
		"--model", "gpt-5.6-sol",
		"--reasoning-effort=",
		"--round", "all",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "gpt-5.6-sol" || run.ReasoningEffort != "" {
		t.Fatalf("expected model-managed resolve selection persisted, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: gpt-5.6-sol") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: model-managed") {
		t.Fatalf("expected resolve progress to show model-managed selection, got %q", stderr.String())
	}
}

func TestRunWatchPersistsEffectiveSelection(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"watch",
		"--source", "coderabbit",
		"--pr", "123",
		"--agent", "codex",
		"--model", "stored-watch-model",
		"--reasoning-effort", "stored-watch-reasoning",
		"--until-clean",
		"--max-rounds", "6",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "stored-watch-model" || run.ReasoningEffort != "stored-watch-reasoning" {
		t.Fatalf("expected stored watch selection, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: stored-watch-model") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: stored-watch-reasoning") {
		t.Fatalf("expected watch progress to show concrete selection, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "auto") {
		t.Fatalf("expected no auto selection placeholder, got %q", stderr.String())
	}
}

func TestRunWatchRendersModelManagedReasoningHeader(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"watch",
		"--source", "coderabbit",
		"--pr", "123",
		"--agent", "codex",
		"--model", "gpt-5.6-sol",
		"--reasoning-effort=",
		"--until-clean",
		"--max-rounds", "1",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "gpt-5.6-sol" || run.ReasoningEffort != "" {
		t.Fatalf("expected model-managed watch selection persisted, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Default Reasoning Effort: model-managed") {
		t.Fatalf("expected watch progress to show model-managed reasoning, got %q", stderr.String())
	}
}

func TestRunOutcomeNotificationsCaptureTerminalResolveWatchAndImplement(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setup        func(t *testing.T) string
		runID        func(t *testing.T, stderr string) string
		wantState    string
		wantKind     string
		wantTarget   string
		wantExitCode int
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := withCLIWorkspace(t)
				withSuccessfulPreflight(t, repoDir)
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
				return homeDir
			},
			runID:      reviewRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindResolve,
			wantTarget: "pr:123",
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := withCLIWorkspace(t)
				withSuccessfulPreflight(t, repoDir)
				return homeDir
			},
			runID:      reviewRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindWatch,
			wantTarget: "pr:123",
		},
		{
			name: "implement",
			args: []string{"implement", "--spec", implementTestSlug, "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
					{id: "task_01", title: "Wire notification outcome"},
				})
				withImplementCollaborators(t, &implementFakeRunner{
					gitRoot: repoDir,
					statusByTask: map[string]spec.Status{
						"task_01": spec.StatusCompleted,
					},
				})
				return homeDir
			},
			runID:      implementRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindImplement,
			wantTarget: "spec:" + implementTestSlug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.setup(t)
			notifier := &recordingOutcomeNotifier{}
			withOutcomeNotifier(t, notifier)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != tt.wantExitCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q stdout=%q", tt.wantExitCode, code, stderr.String(), stdout.String())
			}
			runID := tt.runID(t, stderr.String())
			assertRecordedOutcomes(t, notifier, []roundnotify.Outcome{{
				RunID:  runID,
				State:  tt.wantState,
				Kind:   tt.wantKind,
				Target: tt.wantTarget,
			}})
		})
	}
}

func TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")

	agentStarted := make(chan struct{})
	var startedOnce sync.Once
	releaseAgent := make(chan struct{})
	var releaseOnce sync.Once
	release := func() {
		releaseOnce.Do(func() {
			close(releaseAgent)
		})
	}
	t.Cleanup(release)
	runner := &fakeAgentRunner{onRun: func(agent.ExecuteRequest) error {
		startedOnce.Do(func() {
			close(agentStarted)
		})
		<-releaseAgent
		return nil
	}}
	withAgentRunner(t, runner)
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
		return nil
	}))
	notifier := &recordingOutcomeNotifier{}
	withOutcomeNotifier(t, notifier)

	var resolveStdout bytes.Buffer
	var resolveStderr bytes.Buffer
	resolveDone := make(chan int, 1)
	go func() {
		resolveDone <- RunContext(context.Background(), []string{
			"resolve", "--pr", "123", "--round", "all", "--no-input",
		}, &resolveStdout, &resolveStderr)
	}()

	<-agentStarted
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		release()
		t.Fatalf("open store during completion race: %v", err)
	}
	active, found, err := runStore.ActiveRun(context.Background(), "owner/project", "feature/review")
	closeErr := runStore.Close()
	if err != nil || !found {
		release()
		t.Fatalf("read active Run during completion race: found=%v err=%v", found, err)
	}
	if closeErr != nil {
		release()
		t.Fatalf("close store during completion race: %v", closeErr)
	}

	var stopStdout bytes.Buffer
	var stopStderr bytes.Buffer
	stopCode := RunContext(context.Background(), []string{"stop", "--force", active.ID}, &stopStdout, &stopStderr)
	if stopCode != exitOK {
		release()
		<-resolveDone
		t.Fatalf("Force Stop winner exit = %d, want %d; stdout=%q stderr=%q", stopCode, exitOK, stopStdout.String(), stopStderr.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	release()
	resolveCode := <-resolveDone

	if resolveCode != exitRunFailed {
		t.Fatalf("losing owner exit = %d, want %d; stop stdout=%q stop stderr=%q resolve stderr=%q", resolveCode, exitRunFailed, stopStdout.String(), stopStderr.String(), resolveStderr.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	events := runEvents(t, homeDir, active.ID)
	outcomes := 0
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonOutcome {
			outcomes++
		}
	}
	if outcomes != 1 {
		t.Fatalf("terminal outcome events = %d, want 1; events=%+v", outcomes, events)
	}
	assertRecordedOutcomes(t, notifier, []roundnotify.Outcome{{
		RunID:  active.ID,
		State:  store.StateStopped,
		Kind:   store.KindResolve,
		Target: "pr:123",
	}})

	stopStdout.Reset()
	stopStderr.Reset()
	if replayCode := RunContext(context.Background(), []string{"stop", "--force", active.ID}, &stopStdout, &stopStderr); replayCode != exitOK {
		t.Fatalf("identical Force Stop replay exit = %d, want %d; stderr=%q", replayCode, exitOK, stopStderr.String())
	}
	events = runEvents(t, homeDir, active.ID)
	outcomes = 0
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonOutcome {
			outcomes++
		}
	}
	if outcomes != 1 {
		t.Fatalf("identical replay published another terminal outcome: events=%+v", events)
	}
	assertRecordedOutcomes(t, notifier, []roundnotify.Outcome{{
		RunID:  active.ID,
		State:  store.StateStopped,
		Kind:   store.KindResolve,
		Target: "pr:123",
	}})
}

func TestRunOutcomeNotificationsSkipFetch(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	notifier := &recordingOutcomeNotifier{}
	withOutcomeNotifier(t, notifier)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected fetch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertRecordedOutcomes(t, notifier, nil)
}

func TestNotifyTerminalOutcomeDetachesCanceledParentWithBoundedDeadline(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	notifier := &deadlineRecordingOutcomeNotifier{}
	run := store.Run{
		ID:       "run-123",
		Kind:     store.KindResolve,
		State:    store.StateClean,
		PRNumber: "123",
	}
	var stderr bytes.Buffer
	start := time.Now()

	notifyTerminalOutcome(parent, nil, notifier, &stderr, run)

	if notifier.wasCanceled {
		t.Fatal("expected terminal outcome notification context to ignore parent cancellation")
	}
	if !notifier.hadDeadline {
		t.Fatal("expected terminal outcome notification context to have a deadline")
	}
	if !notifier.deadline.After(start) || notifier.deadline.After(start.Add(outcomeNotificationTimeout+time.Second)) {
		t.Fatalf("expected deadline within notification timeout, got %s from start %s", notifier.deadline, start)
	}
	want := roundnotify.Outcome{
		RunID:  "run-123",
		State:  store.StateClean,
		Kind:   store.KindResolve,
		Target: "pr:123",
	}
	if !reflect.DeepEqual(notifier.outcome, want) {
		t.Fatalf("notification outcome mismatch\nwant: %#v\ngot:  %#v", want, notifier.outcome)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunOutcomeNotificationFailureWarnsAndJournalsWithoutChangingReportOrExit(t *testing.T) {
	wantStdout, _, wantCode, _ := runCleanResolveForOutcomeNotification(t, nil)
	gotStdout, gotStderr, gotCode, homeDir := runCleanResolveForOutcomeNotification(t, errors.New("forced notifier failure"))

	if gotCode != wantCode {
		t.Fatalf("expected exit code to stay %d, got %d stderr=%q", wantCode, gotCode, gotStderr)
	}
	if gotStdout != wantStdout {
		t.Fatalf("stdout changed after notification failure\nwant:\n%q\ngot:\n%q", wantStdout, gotStdout)
	}
	warning := "roundfix: outcome notification failed: forced notifier failure\n"
	if strings.Count(gotStderr, warning) != 1 {
		t.Fatalf("expected one notification warning %q, got stderr=%q", warning, gotStderr)
	}
	runID, events := journaledRunEvents(t, homeDir, gotStderr)
	for _, entry := range events {
		event := entry.Event
		if event.Source == runevent.SourceDaemon &&
			event.Kind == runevent.KindDaemonStatus &&
			strings.Contains(event.Summary, "outcome notification failed") &&
			strings.Contains(string(event.Payload), "forced notifier failure") {
			return
		}
	}
	t.Fatalf("expected Daemon-source notification failure event for Run %s, got %+v", runID, events)
}

func assertCleanCleanupWarningEvent(t *testing.T, homeDir string, stderr string, keptPath string, reason string) {
	t.Helper()
	if strings.Contains(stderr, "failed after Run start") || strings.Contains(stderr, "Run Worktree kept:") {
		t.Fatalf("cleanup warning must not turn Clean into failed/kept diagnostics, got stderr=%q", stderr)
	}
	runID, events := journaledRunEvents(t, homeDir, stderr)
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected cleanup failure to leave Run Clean, got %s", run.State)
	}
	for _, entry := range events {
		event := entry.Event
		if event.Source == runevent.SourceDaemon &&
			event.Kind == runevent.KindDaemonStatus &&
			strings.Contains(event.Summary, "Run Worktree cleanup failed") &&
			strings.Contains(string(event.Payload), keptPath) &&
			strings.Contains(string(event.Payload), reason) {
			return
		}
	}
	t.Fatalf("expected Daemon-source cleanup failure event for Run %s, got %+v", runID, events)
}

func TestRunOutcomeNotificationsDisabledSkipsNotifier(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "notify:\n  enabled: false\n")
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withOutcomeNotifierFactory(t, func(roundconfig.Config) roundnotify.Notifier {
		t.Fatal("disabled notifications must not construct a notifier")
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", wantStdout, stdout.String(), stderr.String())
	}
	if strings.Contains(stderr.String(), "outcome notification failed") {
		t.Fatalf("disabled notifications must be silent, got stderr=%q", stderr.String())
	}
}

func TestRunReviewCommandsRefuseDirtyTrackedCheckout(t *testing.T) {
	for _, command := range []string{"resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			homeDir, repoDir := withReviewGitWorkspace(t)
			withRealReviewPreflight(t, repoDir, true)
			withBranchIntegrity(t, nil, nil)
			mustWrite(t, filepath.Join(repoDir, "README.md"), "dirty tracked change\n")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), branchIntegrityCommandArgs(command), &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit 2, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if !strings.Contains(stderr.String(), "README.md") {
				t.Fatalf("expected dirty tracked path in stderr, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "stash or commit tracked changes") {
				t.Fatalf("expected stash-or-commit next action, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveAllowsUntrackedCheckoutFiles(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withRealReviewPreflight(t, repoDir, true)
	withBranchIntegrity(t, nil, nil)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{}
	withPusher(t, pusher)
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	mustWrite(t, filepath.Join(repoDir, "scratch.txt"), "local notes\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit with untracked file, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.WorkDir != repoDir {
		t.Fatalf("expected Run work_dir to be user checkout %q, got %q", repoDir, run.WorkDir)
	}
	if gitImplementOutput(t, repoDir, "ls-files", "scratch.txt") != "" {
		t.Fatalf("expected pre-existing untracked scratch.txt to remain uncommitted")
	}
	if got := mustRead(t, filepath.Join(repoDir, "scratch.txt")); got != "local notes\n" {
		t.Fatalf("expected untracked file to remain unchanged, got %q", got)
	}
}

func TestRunResolveCommitsOnUserBranchWithoutRunBranch(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withRealReviewPreflight(t, repoDir, true)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{}
	withPusher(t, pusher)
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean Run, got %s", run.State)
	}
	if run.WorkDir != repoDir {
		t.Fatalf("expected Run work_dir to be user checkout %q, got %q", repoDir, run.WorkDir)
	}
	if len(runner.gitRoots) != 1 || runner.gitRoots[0] != repoDir {
		t.Fatalf("expected Agent to run in user checkout %q, got %v", repoDir, runner.gitRoots)
	}
	branch := runworktree.BranchName(runID)
	assertRunBranchRemoved(t, repoDir, branch)
	if strings.Contains(stderr.String(), "Run Worktree:") || strings.Contains(stderr.String(), "Integration command:") {
		t.Fatalf("resolve must not report a Run Worktree or integration command, got %q", stderr.String())
	}
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent work\n" {
		t.Fatalf("expected agent work in user checkout, got %q", got)
	}
	if subjects := gitImplementOutput(t, repoDir, "log", "--format=%s", "--", "agent.txt"); !strings.Contains(subjects, daemon.BatchCommitMessage(1)) {
		t.Fatalf("expected batch commit on user branch, got subjects %q", subjects)
	}
}

func TestRunResolvePushRunsFromUserCheckoutWithoutRunWorktree(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withRealReviewPreflight(t, repoDir, true)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{
		check: func(req daemon.PushRequest) error {
			if req.WorkDir != repoDir {
				return fmt.Errorf("push WorkDir left the user checkout: %q", req.WorkDir)
			}
			if got := gitImplementOutput(t, repoDir, "show", "HEAD:agent.txt"); got != "agent work\n" {
				return fmt.Errorf("push ran before integrated agent work reached HEAD, got %q", got)
			}
			return nil
		},
	}
	withPusher(t, pusher)
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one Final Push after integration, got %d", pusher.calls)
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean Run, got %s", run.State)
	}
	if run.WorkDir != repoDir {
		t.Fatalf("expected Run work_dir to be user checkout %q, got %q", repoDir, run.WorkDir)
	}
	if len(runner.gitRoots) != 1 || runner.gitRoots[0] != repoDir {
		t.Fatalf("expected Agent to run in user checkout %q, got %v", repoDir, runner.gitRoots)
	}
	if len(pusher.workDir) != 1 || pusher.workDir[0] != repoDir {
		t.Fatalf("expected Final Push from the user checkout %q, got %v", repoDir, pusher.workDir)
	}
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if strings.Contains(stderr.String(), "Run Worktree:") || strings.Contains(stderr.String(), "Integration command:") {
		t.Fatalf("resolve must not report a Run Worktree or integration command, got %q", stderr.String())
	}
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent work\n" {
		t.Fatalf("expected integrated agent work in user checkout, got %q", got)
	}
	subject := gitImplementOutput(t, repoDir, "log", "-1", "--format=%s")
	if !strings.HasPrefix(subject, "docs: review round") {
		t.Fatalf("expected review artifact docs commit at HEAD, got %q", subject)
	}
}

func TestRunWatchReusesUserCheckoutAcrossRoundsWithoutRunBranch(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	withRealReviewPreflight(t, repoDir, true)
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusSettled},
			{State: watch.StatusSettled},
		},
	}).Status)
	fetchCalls := 0
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		fetchCalls++
		return []reviewsource.ReviewItem{
			{
				Title:                   fmt.Sprintf("major: handle watch issue %d", fetchCalls),
				File:                    fmt.Sprintf("internal/watch_%d.go", fetchCalls),
				Line:                    12,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "Watch issue.",
				SourceRef:               fmt.Sprintf("thread:PRRT_watch_%d,comment:PRRC_watch_%d", fetchCalls, fetchCalls),
				ReviewHash:              fmt.Sprintf("review-hash-watch-%d", fetchCalls),
				SourceReviewID:          fmt.Sprintf("90%02d", fetchCalls),
				SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, fetchCalls, 0, 0, time.UTC),
			},
		}, nil
	})
	headCheck := &fakeWatchHeadCheck{states: []watch.HeadCheckState{watch.CheckFailure, watch.CheckSuccess}}
	withWatchHeadCheck(t, headCheck.Check)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{}
	withPusher(t, pusher)
	runCount := 0
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			runCount++
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte(fmt.Sprintf("agent round %d\n", runCount)), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "2", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if runner.calls != 2 {
		t.Fatalf("expected two watched resolve rounds, got %d", runner.calls)
	}
	if len(runner.gitRoots) != 2 || runner.gitRoots[0] != repoDir || runner.gitRoots[1] != repoDir {
		t.Fatalf("expected user checkout reused across rounds %q, got %v", repoDir, runner.gitRoots)
	}
	if pusher.calls != 2 {
		t.Fatalf("expected Final Push after each integrated clean round, got %d", pusher.calls)
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean watch Run, got %s", run.State)
	}
	if run.WorkDir != repoDir {
		t.Fatalf("expected recorded work_dir to be user checkout %q, got %q", repoDir, run.WorkDir)
	}
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if strings.Contains(stderr.String(), "Run Worktree:") || strings.Contains(stderr.String(), "Integration command:") {
		t.Fatalf("watch must not report a Run Worktree or integration command, got %q", stderr.String())
	}
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent round 2\n" {
		t.Fatalf("expected latest integrated watch work in user checkout, got %q", got)
	}
	if !strings.Contains(stdout.String(), "This Run (Clean after 2 Round(s)):") {
		t.Fatalf("expected clean two-round stdout report, got %q", stdout.String())
	}
}

func TestRunWatchPrintsBudgetExceededStdoutReport(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	clock := &fakeWatchClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	withWatchTiming(t, clock, &fakeWatchSleeper{clock: clock})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusPending},
			{State: watch.StatusSettled},
		},
	}).Status)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1s
  review_timeout: 10s
  quiet_period: 2s
budget:
  enabled: true
  max_run_duration: 2s
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected BudgetExceeded exit %d, got %d stderr=%q", exitRunFailed, code, stderr.String())
	}
	wantStdout := "" +
		"This Run (BudgetExceeded after 0 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected BudgetExceeded stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "reached BudgetExceeded") {
		t.Fatalf("expected BudgetExceeded terminal outcome on stderr, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("BudgetExceeded before fetch must not fetch or run Agent, got %q", stderr.String())
	}
}

func TestOperationalStdoutReportStartsAfterTerminalRunLine(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setup          func(t *testing.T, repoDir string)
		terminalNeedle string
		wantStdout     string
	}{
		{
			name:           "watch",
			args:           []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"},
			terminalNeedle: "reached Clean",
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
				"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n",
		},
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			},
			terminalNeedle: "reached Clean",
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
				"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.setup != nil {
				tt.setup(t, repoDir)
			}
			recorder := &orderedWriteRecorder{}

			code := RunContext(context.Background(), tt.args, recorder.writer("stdout"), recorder.writer("stderr"))

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, recorder.content("stderr"))
			}
			if got := recorder.content("stdout"); got != tt.wantStdout {
				t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", tt.wantStdout, got, recorder.content("stderr"))
			}
			terminalIndex := recorder.firstIndexContaining("stderr", tt.terminalNeedle)
			if terminalIndex < 0 {
				t.Fatalf("expected terminal stderr line containing %q, got %q", tt.terminalNeedle, recorder.content("stderr"))
			}
			stdoutIndex := recorder.firstIndex("stdout")
			if stdoutIndex < 0 {
				t.Fatal("expected stdout report")
			}
			if stdoutIndex < terminalIndex {
				t.Fatalf("stdout report began before terminal Run line: events=%#v", recorder.events)
			}
		})
	}
}

func TestRunResolveAgentFullAccessIsExplicitOptIn(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		args       []string
		wantAccess string
	}{
		{
			name:       "default sandboxed",
			args:       []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
			wantAccess: "",
		},
		{
			name:       "command flag opts in",
			args:       []string{"resolve", "--pr", "123", "--agent-full-access", "--round", "all", "--no-input"},
			wantAccess: "full-access",
		},
		{
			name: "config opts in",
			config: `
defaults:
  agent_full_access: true
`,
			args:       []string{"resolve", "--pr", "123", "--round", "all", "--no-input"},
			wantAccess: "full-access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &fakeAgentRunner{}
			withAgentRunner(t, runner)
			if tt.config != "" {
				mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), tt.config)
			}
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d, stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if len(runner.probedRuntimes) != 2 {
				t.Fatalf("expected preferred and fallback profile probes, got %#v", runner.probedRuntimes)
			}
			for _, runtime := range runner.probedRuntimes {
				if runtime.FullAccessMode != tt.wantAccess {
					t.Fatalf("expected probe full-access mode %q, got runtime %#v", tt.wantAccess, runtime)
				}
			}
			if len(runner.runRuntimes) != 1 {
				t.Fatalf("expected one Agent run, got %#v", runner.runRuntimes)
			}
			if runner.runRuntimes[0].FullAccessMode != tt.wantAccess {
				t.Fatalf("expected run full-access mode %q, got %q", tt.wantAccess, runner.runRuntimes[0].FullAccessMode)
			}
		})
	}
}

func TestRunFetchReusesMatchingAutoRound(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	args := []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}

	var firstStdout bytes.Buffer
	var firstStderr bytes.Buffer
	firstCode := Run(args, &firstStdout, &firstStderr)
	if firstCode != 0 {
		t.Fatalf("expected first fetch exit 0, got %d stdout=%q stderr=%q", firstCode, firstStdout.String(), firstStderr.String())
	}

	var secondStdout bytes.Buffer
	var secondStderr bytes.Buffer
	secondCode := Run(args, &secondStdout, &secondStderr)
	if secondCode != 0 {
		t.Fatalf("expected second fetch exit 0, got %d stdout=%q stderr=%q", secondCode, secondStdout.String(), secondStderr.String())
	}

	if !strings.Contains(secondStdout.String(), "Artifacts: reused existing matching Round") {
		t.Fatalf("expected second fetch to report reused Round, got %q", secondStdout.String())
	}
	if _, err := os.Stat(filepath.Join(defaultReviewRootForRepo(repoDir, "123"), "round-002")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no duplicate round-002, got err %v", err)
	}
	assertRunCount(t, filepath.Join(homeDir, ".roundfix", "roundfix.db"), 2)
}

func TestRunWatchTimeoutOffersManualReviewWithoutFetching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("watch timeout must not fetch Review Source issues")
		return nil, nil
	})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusPending},
			{State: watch.StatusPending},
			{State: watch.StatusPending},
		},
	}).Status)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1s
  review_timeout: 2s
  quiet_period: 1s
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected watch timeout exit 1, got %d", code)
	}
	wantStdout := "" +
		"This Run (TimedOut after 0 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected timeout stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Review Source status: pending") {
		t.Fatalf("expected pending Review Source status output, got %q", stderr.String())
	}
	if count := strings.Count(stderr.String(), "Review Source status: pending\n"); count != 1 {
		t.Fatalf("expected one unchanged status-poll line, got %d in %q", count, stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached TimedOut") {
		t.Fatalf("expected TimedOut terminal outcome, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "@coderabbitai review") {
		t.Fatalf("expected manual review trigger guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("timeout must not fetch or run Agent, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchMissingHeadCheckEndsCleanUnverified(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1ns
  check_grace_period: 1ns
`)
	withWatchHeadCheck(t, (&fakeWatchHeadCheck{
		states: []watch.HeadCheckState{watch.CheckMissing},
	}).Check)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != exitUnverified {
		t.Fatalf("expected CleanUnverified watch exit %d, got %d stderr=%q", exitUnverified, code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached CleanUnverified") {
		t.Fatalf("expected CleanUnverified terminal outcome, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Next: confirm the pull request's Review Source check before merging.") {
		t.Fatalf("expected CleanUnverified next action, got %q", stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"This Run (Clean Unverified after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected CleanUnverified stdout report %q, got %q", wantStdout, stdout.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchReusesOneAgentSessionAcrossRoundsAndCloses(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	inner := &fakeAgentRunner{statuses: []string{rounds.StatusResolved, rounds.StatusFailed, rounds.StatusResolved}}
	runner := &sessionRecordingRunner{inner: inner}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	fetchCalls := 0
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		fetchCalls++
		switch fetchCalls {
		case 1:
			return []reviewsource.ReviewItem{
				{
					Title:                   "major: handle first issue",
					File:                    "internal/first.go",
					Line:                    12,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "First issue.",
					SourceRef:               "thread:PRRT_first,comment:PRRC_first",
					ReviewHash:              "review-hash-first",
					SourceReviewID:          "9001",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
				},
				{
					Title:                   "major: handle retry issue",
					File:                    "internal/retry.go",
					Line:                    24,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "Retry issue.",
					SourceRef:               "thread:PRRT_retry_1,comment:PRRC_retry_1",
					ReviewHash:              "review-hash-retry",
					SourceReviewID:          "9002",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
				},
			}, nil
		default:
			return []reviewsource.ReviewItem{
				{
					Title:                   "major: handle retry issue",
					File:                    "internal/retry.go",
					Line:                    24,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "Retry issue.",
					SourceRef:               "thread:PRRT_retry_2,comment:PRRC_retry_2",
					ReviewHash:              "review-hash-retry",
					SourceReviewID:          "9003",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 2, 0, 0, time.UTC),
				},
			}, nil
		}
	})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusSettled},
			{State: watch.StatusSettled},
		},
	}).Status)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "2", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit 0, got %d stderr=%q", code, stderr.String())
	}
	if fetchCalls != 2 {
		t.Fatalf("expected two fetched Rounds, got %d", fetchCalls)
	}
	if !strings.Contains(stderr.String(), "after 2 Round(s)") {
		t.Fatalf("expected Watch Run to span two Rounds, got %q", stderr.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	assertRecordedOneSessionForRun(t, runner, runID, inner.calls)
	assertSessionLifecycleEvents(t, events)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchStopRequestBeforeAgentMarksStopped(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	statusCalls := 0
	withWatchStatus(t, func(ctx context.Context, _ reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
		statusCalls++
		if statusCalls > 1 {
			t.Fatalf("Stop Request must prevent another Review Source status call")
		}
		runStore, err := store.Open(ctx, homeDir)
		if err != nil {
			t.Fatalf("open Run Database to request Stop: %v", err)
		}
		defer func() {
			if err := runStore.Close(); err != nil {
				t.Fatalf("close Run Database after Stop Request: %v", err)
			}
		}()
		active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
		if err != nil {
			t.Fatalf("read active watch Run: %v", err)
		}
		if !found {
			t.Fatal("expected active watch Run before Review Source status")
		}
		if err := runStore.RequestStop(ctx, active.ID); err != nil {
			t.Fatalf("request Stop for watch Run: %v", err)
		}
		return reviewsource.WatchStatus{State: watch.StatusPending}, nil
	})
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("Stop Request before Agent must not fetch Review Source issues")
		return nil, nil
	})
	withChangedPaths(t, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	wantStdout := "" +
		"This Run (Stopped after 0 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stopped stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Watch Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Watch Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Changed paths after Stop Request: none") {
		t.Fatalf("expected changed-path report, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("Stop Request before Agent must not fetch or run Agent, got %q", stderr.String())
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertStopRequested(t, homeDir, runID)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveStopRequestDuringAgentPreservesWorkAndSkipsDaemonMutations(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPathSnapshots(t, nil, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 unresolved — major: handle test issue\n" +
		"This Run (Stopped after 1 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 0 failed, 1 unresolved.\n" +
		"Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 0 failed, 1 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stopped stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Resolve Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Resolve Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "partial agent output") {
		t.Fatalf("expected persisted Agent output in stream, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "M src/app.go") {
		t.Fatalf("expected changed path report, got %q", stderr.String())
	}
	if verifier.calls != 0 {
		t.Fatalf("Stop Request must not run verification, got %d calls", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("Stop Request must not create commits, got %d calls", committer.calls)
	}
	if sourceResolver.calls != 0 {
		t.Fatalf("Stop Request must not resolve Review Source threads, got %d calls", sourceResolver.calls)
	}
	if pusher.calls != 0 {
		t.Fatalf("Stop Request must not push, got %d calls", pusher.calls)
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusPending {
		t.Fatalf("Stop Request must preserve assigned issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")

}

func TestRunResolveSIGINTStopReturns130(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPaths(t, nil)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stoppedCode := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)
	code := exitForInterrupt(stoppedCode, true)

	if code != 130 {
		t.Fatalf("expected SIGINT exit 130, got %d", code)
	}
	if !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected Stopped output before SIGINT exit, got %q", stderr.String())
	}
}

func TestExitForWatchOutcome(t *testing.T) {
	tests := []struct {
		outcome string
		code    int
	}{
		{outcome: store.StateClean, code: exitOK},
		{outcome: store.StateCleanUnverified, code: exitUnverified},
		{outcome: store.StateMaxRoundsReached, code: exitOK},
		{outcome: store.StateStopped, code: exitOK},
		{outcome: store.StateBudgetExceeded, code: exitRunFailed},
		{outcome: store.StateTimedOut, code: exitRunFailed},
		{outcome: store.StateFailed, code: exitRunFailed},
		{outcome: store.StateUnresolved, code: exitRunFailed},
		{outcome: store.StateFetched, code: exitRunFailed},
	}

	for _, tt := range tests {
		t.Run(tt.outcome, func(t *testing.T) {
			if code := exitForWatchOutcome(tt.outcome); code != tt.code {
				t.Fatalf("expected exit code %d for %s, got %d", tt.code, tt.outcome, code)
			}
		})
	}
}

func TestRunResolveRejectsMissingCompatibleArtifactsBeforeRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no downloaded Unresolved Review Issues in Compatible Artifacts") {
		t.Fatalf("expected missing Compatible Artifacts error, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix fetch --source coderabbit --pr 123") {
		t.Fatalf("expected fetch guidance, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveHonorsRoundSelector(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--round", "2", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Unresolved resolve exit code 1, got %d", code)
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"This Run (Unresolved after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 1 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected Unresolved stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 1 downloaded Unresolved Review Issue") {
		t.Fatalf("expected one selected issue, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Round scope: 002") {
		t.Fatalf("expected round scope in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push blocked: 1 Unresolved Review Issue") {
		t.Fatalf("expected Final Push to be blocked by unselected unresolved issue, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome for out-of-scope unresolved issue, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveDeduplicatesBeforeBatching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit code 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 duplicated — major: handle test issue\n" +
		"issue 002 resolved — major: handle test issue\n" +
		"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 1 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 1 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected duplicate stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 2 downloaded Unresolved Review Issue(s)") {
		t.Fatalf("expected selected issue count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "assigned 1 newest occurrence(s) into 1 Batch(es)") {
		t.Fatalf("expected deduped assignment count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "associated 1 older duplicate occurrence(s)") {
		t.Fatalf("expected duplicate association count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Marked 1 older duplicate Review Issue occurrence") {
		t.Fatalf("expected older duplicate marker, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
		t.Fatalf("expected one unique source thread resolution, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push completed") {
		t.Fatalf("expected Final Push after all duplicates terminal, got %q", stderr.String())
	}
	if sourceResolver.calls != 2 {
		t.Fatalf("expected duplicate outcome comment plus one unique thread resolve, got %d", sourceResolver.calls)
	}
	if len(sourceResolver.commentRequests) != 1 || !strings.Contains(sourceResolver.commentRequests[0].Body, "duplicated") {
		t.Fatalf("expected duplicate outcome comment, got %+v", sourceResolver.commentRequests)
	}
	if got := len(sourceResolver.resolveRequests); got != 1 {
		t.Fatalf("expected one unique source thread resolve, got %d", got)
	}
	if sourceResolver.resolveRequests[0].SourceRef != "thread:PRRT_test,comment:PRRC_test" {
		t.Fatalf("expected duplicate source ref resolved once, got %#v", sourceResolver.resolveRequests[0])
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveVerificationFailureDoesNotCommit(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{err: errors.New("tests failed")}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	if !strings.Contains(stdout.String(), "issue 001 failed — major: handle test issue — reason: Verification failed") ||
		!strings.Contains(stdout.String(), "tests failed") ||
		!strings.Contains(stdout.String(), "This Run (Unresolved after 1 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.") ||
		!strings.Contains(stdout.String(), "Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.") {
		t.Fatalf("expected failed stdout report with reason, got %q", stdout.String())
	}
	if verifier.calls != 2 {
		t.Fatalf("expected initial and final verification calls, got %d", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("expected no Batch commit after verification failure, got %d", committer.calls)
	}
	if len(sourceResolver.resolveRequests) != 0 {
		t.Fatalf("expected failed Review Issue to stay unresolved on Review Source, got %+v", sourceResolver.resolveRequests)
	}
	if len(sourceResolver.commentRequests) == 0 {
		t.Fatal("expected failed Review Issue outcome comment after verification failure")
	}
	if pusher.calls != 0 {
		t.Fatalf("expected no Final Push after verification failure, got %d", pusher.calls)
	}
	if !strings.Contains(stderr.String(), "tests failed") {
		t.Fatalf("expected verification failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome after failed verification, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveProcessesAllBatchesBeforeFinalPush(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	result := persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle first issue\n" +
		"issue 002 resolved — major: handle second issue\n" +
		"This Run (Clean after 1 Round(s)): 2 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 2 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected clean stdout report %q, got %q", wantStdout, stdout.String())
	}
	if verifier.calls != 2 {
		t.Fatalf("expected one verification call per Batch, got %d", verifier.calls)
	}
	if committer.calls != 2 {
		t.Fatalf("expected one commit per Batch, got %d", committer.calls)
	}
	if len(sourceResolver.resolveRequests) != 2 {
		t.Fatalf("expected one Review Source resolve per Batch, got %+v", sourceResolver.resolveRequests)
	}
	if pusher.calls != 1 {
		t.Fatalf("expected Final Push after all Batches complete, got %d", pusher.calls)
	}
	if committer.messages[0] != daemon.BatchCommitMessage(1) {
		t.Fatalf("expected Batch commit message %q, got %q", daemon.BatchCommitMessage(1), committer.messages[0])
	}
	if committer.messages[1] != daemon.BatchCommitMessage(2) {
		t.Fatalf("expected Batch commit message %q, got %q", daemon.BatchCommitMessage(2), committer.messages[1])
	}
	if !strings.Contains(stderr.String(), "assigned 2 newest occurrence(s) into 2 Batch(es)") {
		t.Fatalf("expected two planned Batches, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Batch: 001/002") || !strings.Contains(stderr.String(), "Batch: 002/002") {
		t.Fatalf("expected both Batch progress lines, got %q", stderr.String())
	}
	first, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(result.IssuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusResolved {
		t.Fatalf("expected first Batch issue resolved, got %q", first.Status)
	}
	if second.Status != rounds.StatusResolved {
		t.Fatalf("expected second Batch issue resolved, got %q", second.Status)
	}
	if strings.Contains(stderr.String(), "Final Push blocked") {
		t.Fatalf("did not expect Final Push blocked after all Batches complete, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push completed: git push origin HEAD:feature/review") {
		t.Fatalf("expected Final Push output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveFinalPushRunsOnceAfterAllUnresolvedTerminal(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	persistCLIReviewItems(t, repoDir, 2, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call for one Batch, got %d", verifier.calls)
	}
	if committer.calls != 1 {
		t.Fatalf("expected one Batch commit, got %d", committer.calls)
	}
	if len(sourceResolver.resolveRequests) != 2 {
		t.Fatalf("expected both assigned terminal issues to resolve source threads, got %+v", sourceResolver.resolveRequests)
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one Final Push, got %d", pusher.calls)
	}
	if pusher.remotes[0] != "origin" || pusher.branches[0] != "feature/review" {
		t.Fatalf("expected push to origin feature/review, got %s %s", pusher.remotes[0], pusher.branches[0])
	}
	if strings.Contains(strings.Join(pusher.args, " "), "--force") {
		t.Fatalf("Final Push must not force-push, got args %#v", pusher.args)
	}
	if !strings.Contains(stderr.String(), "Final Push completed: git push origin HEAD:feature/review") {
		t.Fatalf("expected Final Push output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveResolvesInvalidAssignedIssueSourceThread(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withAgentRunner(t, &fakeAgentRunner{status: rounds.StatusInvalid})
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if len(sourceResolver.commentRequests) != 1 || !strings.Contains(sourceResolver.commentRequests[0].Body, "invalid") {
		t.Fatalf("expected invalid outcome comment before resolution, got %+v", sourceResolver.commentRequests)
	}
	if len(sourceResolver.resolveRequests) != 1 {
		t.Fatalf("expected invalid issue to resolve source thread, got %+v", sourceResolver.resolveRequests)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveProbeFailureDoesNotCreateRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("adapter missing")})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected probe preflight exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "adapter missing") {
		t.Fatalf("expected probe diagnostic, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunResolveSelectionPreflightRejectionReportsTupleAndCreatesNoRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	probeErr := &agent.SelectionPreflightError{
		Runtime:         "codex",
		Model:           "gpt-5.6-sol",
		ReasoningEffort: "xhigh",
		Err:             errors.New("missing model metadata"),
	}
	runner := &fakeAgentRunner{probeErr: probeErr}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
runtimes:
  codex:
    model: gpt-5.6-sol
    reasoning_effort: xhigh
`)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, want := range []string{
		"codex",
		"gpt-5.6-sol",
		"xhigh",
		"missing model metadata",
		"update the ACP Runtime or adapter",
		"choose supported Agent Model and Default Reasoning Effort values",
		`set runtimes.codex.reasoning_effort "" when the model manages reasoning`,
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	if runner.calls != 0 {
		t.Fatalf("expected no durable Agent Session run, got %d run call(s)", runner.calls)
	}
	if len(runner.probeRequests) != 1 {
		t.Fatalf("expected one selection preflight attempt, got %#v", runner.probeRequests)
	}
	if runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected preflight workdir %q, got %q", repoDir, runner.probeRequests[0].WorkDir)
	}
	if got := runner.probeRequests[0].Runtime; got.Model != "gpt-5.6-sol" || got.ReasoningEffort != "xhigh" {
		t.Fatalf("expected exact selection without fallback, got %#v", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunWatchSelectionPreflightFailureCreatesNoRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{probeErr: &agent.SelectionPreflightError{
		Runtime:         "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "unsupported",
		Err:             errors.New("reasoning rejected"),
	}}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "reasoning rejected") || !strings.Contains(stderr.String(), "unsupported") {
		t.Fatalf("expected reasoning rejection diagnostic, got %q", stderr.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected no durable Agent Session run, got %d run call(s)", runner.calls)
	}
	if len(runner.probeRequests) != 1 || runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected one selection preflight in git root %q, got %#v", repoDir, runner.probeRequests)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunReviewAgentCommandsReportProfileProofFailureWithoutCreatingRun(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		needsReviewIssue bool
	}{
		{
			name:             "resolve with explicit reasoning",
			args:             []string{"resolve", "--pr", "123", "--agent", "codex", "--model", "broken-model", "--reasoning-effort", "unsupported", "--no-input"},
			needsReviewIssue: true,
		},
		{
			name: "watch with model-managed reasoning",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--model", "broken-model", "--reasoning-effort", "unsupported", "--no-input"},
		},
		{
			name:             "resolve with non-interactive stderr",
			args:             []string{"resolve", "--pr", "123", "--agent", "codex", "--model", "broken-model", "--reasoning-effort", "unsupported"},
			needsReviewIssue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &fakeAgentRunner{
				probeErr: &agent.SelectionPreflightError{
					Runtime:         "codex",
					Model:           "broken-model",
					ReasoningEffort: "unsupported",
					Err:             errors.New("selection rejected"),
				},
				fallback:   agent.FallbackSelection{Model: "should-not-be-used", ReasoningEffort: "high"},
				fallbackOK: true,
			}
			withAgentRunner(t, runner)
			withFakeWorktree(t)
			if tt.needsReviewIssue {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected selection preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("profile proof report must stay off stdout, got %q", stdout.String())
			}
			for _, want := range []string{
				`profile proof failed for runtime "codex", model "broken-model", reasoning_effort "unsupported"`,
				"review preferred",
				"selection rejected",
				"roundfix profiles configure --scope user|project",
				"roundfix profiles validate",
			} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			for _, forbidden := range []string{"Fallback Selection:", "Probed Agent Models", "Re-run:"} {
				if strings.Contains(stderr.String(), forbidden) {
					t.Fatalf("profile proof must not report dynamic fallback %q in %q", forbidden, stderr.String())
				}
			}
			if runner.calls != 0 {
				t.Fatalf("expected no Agent work, got %d call(s)", runner.calls)
			}
			if len(runner.fallbackModels) != 0 {
				t.Fatalf("profile proof must not discover dynamic fallback candidates, got %#v", runner.fallbackModels)
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveSelectionFailureDoesNotProbeDynamicCandidates(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{probeErr: &agent.SelectionPreflightError{
		Runtime:         "codex",
		Model:           "gpt-5.6-sol",
		ReasoningEffort: "unsupported",
		Err:             errors.New("selection rejected"),
	}}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--model", "gpt-5.6-sol", "--reasoning-effort", "unsupported", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d", exitPreflight, code)
	}
	for _, want := range []string{
		"selection rejected",
		"recovery: update the ACP Runtime or adapter",
		"profile proof failed",
		"review preferred",
		"roundfix profiles configure",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	for _, forbidden := range []string{"Fallback probe", "Probed Agent Models", "Probed reasoning efforts"} {
		if strings.Contains(stderr.String(), forbidden) {
			t.Fatalf("profile proof must not report dynamic candidates %q in %q", forbidden, stderr.String())
		}
	}
	if stdout.Len() != 0 || runner.calls != 0 {
		t.Fatalf("expected no output or Agent work, stdout=%q calls=%d", stdout.String(), runner.calls)
	}
	if len(runner.fallbackModels) != 0 {
		t.Fatalf("profile proof must not discover dynamic fallback candidates, got %#v", runner.fallbackModels)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunReviewAgentCommandsDoNotPromptForDynamicFallback(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		needsReviewIssue bool
	}{
		{
			name:             "resolve",
			args:             []string{"resolve", "--pr", "123"},
			needsReviewIssue: true,
		},
		{
			name: "watch with model-managed reasoning",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			configPath := filepath.Join(repoDir, ".roundfixrc.yml")
			configContent := "runtimes:\n  codex:\n    model: broken-model\n    reasoning_effort: unsupported\n"
			mustWrite(t, configPath, configContent)
			runner := &fakeAgentRunner{
				probeErr: &agent.SelectionPreflightError{
					Runtime:         "codex",
					Model:           "broken-model",
					ReasoningEffort: "unsupported",
					Err:             errors.New("selection rejected"),
				},
				fallback:   agent.FallbackSelection{Model: "should-not-be-used", ReasoningEffort: "high"},
				fallbackOK: true,
			}
			withAgentRunner(t, runner)
			withFallbackConfirmation(t, "yes\n")
			if tt.needsReviewIssue {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected profile proof preflight exit %d, got exit %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if strings.Contains(stderr.String(), "Fallback Selection:") || strings.Contains(stderr.String(), "Use this Fallback Selection for this Run?") {
				t.Fatalf("profile proof must not prompt for dynamic fallback candidates, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "roundfix profiles configure") || !strings.Contains(stderr.String(), "roundfix profiles validate") {
				t.Fatalf("expected profile remediation guidance, got %q", stderr.String())
			}
			if got := mustRead(t, configPath); got != configContent {
				t.Fatalf("failed profile proof must not change Project Config\nwant: %q\n got: %q", configContent, got)
			}
			if stdout.Len() != 0 || runner.calls != 0 || len(runner.fallbackModels) != 0 {
				t.Fatalf("expected no output, Agent work, or dynamic fallback probes; stdout=%q calls=%d fallback=%#v", stdout.String(), runner.calls, runner.fallbackModels)
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveProfileProofFailureIgnoresFallbackConfirmationInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "no", input: "no\n"},
		{name: "empty", input: "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			configPath := filepath.Join(repoDir, ".roundfixrc.yml")
			configContent := "runtimes:\n  codex:\n    model: broken-model\n    reasoning_effort: unsupported\n"
			mustWrite(t, configPath, configContent)
			runner := &fakeAgentRunner{
				probeErr: &agent.SelectionPreflightError{
					Runtime:         "codex",
					Model:           "broken-model",
					ReasoningEffort: "unsupported",
					Err:             errors.New("selection rejected"),
				},
				fallback:   agent.FallbackSelection{Model: "gpt-5.5", ReasoningEffort: "high"},
				fallbackOK: true,
			}
			withAgentRunner(t, runner)
			withFallbackConfirmation(t, tt.input)
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"resolve", "--pr", "123"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected profile proof preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			for _, want := range []string{"Preflight failed", "profile proof failed", "roundfix profiles configure", "roundfix profiles validate"} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected profile proof output to contain %q, got %q", want, stderr.String())
				}
			}
			if strings.Contains(stderr.String(), "Use this Fallback Selection for this Run?") || strings.Contains(stderr.String(), "Fallback Selection:") {
				t.Fatalf("profile proof failure must not prompt for dynamic fallback candidates, got %q", stderr.String())
			}
			if stdout.Len() != 0 || runner.calls != 0 || len(runner.fallbackModels) != 0 {
				t.Fatalf("profile proof failure must start no work or dynamic fallback, stdout=%q Agent calls=%d fallback=%#v", stdout.String(), runner.calls, runner.fallbackModels)
			}
			assertNoRunDatabase(t, homeDir)
			if got := mustRead(t, configPath); got != configContent {
				t.Fatalf("failed profile proof must not change Project Config\nwant: %q\n got: %q", configContent, got)
			}
		})
	}
}

func TestRunResolveNoInputProfileProofFailureReportsRemediation(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{
		probeErr: &agent.SelectionPreflightError{Runtime: "codex", Model: "broken-model", Err: errors.New("selection rejected")},
		fallback: agent.FallbackSelection{Model: "gpt-5.5", ReasoningEffort: "high"}, fallbackOK: true,
	}
	withAgentRunner(t, runner)
	withFallbackConfirmation(t, "yes\n")
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--model", "broken-model", "--reasoning-effort", "xhigh", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected profile proof preflight exit %d, got %d", exitPreflight, code)
	}
	if strings.Contains(stderr.String(), "Use this Fallback Selection for this Run?") {
		t.Fatalf("profile proof failure must not prompt, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix profiles configure") || !strings.Contains(stderr.String(), "roundfix profiles validate") {
		t.Fatalf("expected profile remediation guidance, got %q", stderr.String())
	}
	if stdout.Len() != 0 || runner.calls != 0 || len(runner.fallbackModels) != 0 {
		t.Fatalf("profile proof failure must create no Run, Agent work, or dynamic fallback; stdout=%q calls=%d fallback=%#v", stdout.String(), runner.calls, runner.fallbackModels)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunResolveDetachedChildReportsProfileProofFailure(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{
		probeErr: &agent.SelectionPreflightError{
			Runtime: "codex",
			Model:   "broken-model",
			Err:     errors.New("selection rejected"),
		},
		fallback:   agent.FallbackSelection{Model: "gpt-5.5", ReasoningEffort: "xhigh"},
		fallbackOK: true,
	}
	withAgentRunner(t, runner)
	withFallbackConfirmation(t, "yes\n")
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("create detached child pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = readPipe.Close()
		_ = writePipe.Close()
	})
	t.Setenv(detachHandshakeFDEnv, strconv.Itoa(int(writePipe.Fd())))
	t.Setenv(detachConsoleTempEnv, filepath.Join(t.TempDir(), "console.tmp"))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--model", "broken-model", "--reasoning-effort", "xhigh"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected detached preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "profile proof failed") || !strings.Contains(stderr.String(), "roundfix profiles configure") {
		t.Fatalf("expected profile proof remediation, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Use this Fallback Selection for this Run?") || strings.Contains(stderr.String(), "Re-run:") {
		t.Fatalf("detached profile proof failure must not prompt or suggest dynamic fallback rerun, got %q", stderr.String())
	}
	if stdout.Len() != 0 || runner.calls != 0 || len(runner.fallbackModels) != 0 {
		t.Fatalf("expected detached failure to create no Run, Agent work, or dynamic fallback, stdout=%q calls=%d fallback=%#v", stdout.String(), runner.calls, runner.fallbackModels)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunReviewAgentCommandsPassOneRunSelectionOverridesToPreflight(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--model", "future-model", "--reasoning-effort", "experimental-reasoning", "--no-input"},
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--model", "future-model", "--reasoning-effort", "experimental-reasoning", "--no-input"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &fakeAgentRunner{probeErr: errors.New("stop after selection preflight")}
			withAgentRunner(t, runner)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
			}
			if len(runner.probeRequests) != 1 {
				t.Fatalf("expected one selection preflight, got %#v", runner.probeRequests)
			}
			got := runner.probeRequests[0].Runtime
			if got.Model != "future-model" || got.ReasoningEffort != "experimental-reasoning" {
				t.Fatalf("expected one-Run selection overrides at preflight, got %#v", got)
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunReviewAgentCommandsRejectExplicitEmptySelectionOverrides(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "resolve model",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--model=", "--reasoning-effort", "high", "--no-input"},
			want: "--model cannot be empty",
		},
		{
			name: "watch model",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--model=", "--reasoning-effort", "high", "--no-input"},
			want: "--model cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.want, stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveSelectionOverrideRejectsPartialBeforeConfigLoad(t *testing.T) {
	partials := []struct {
		name string
		args []string
	}{
		{name: "agent only detached", args: []string{"--detach", "--agent", "codex"}},
		{name: "model only", args: []string{"--model", "gpt-5.6-sol"}},
		{name: "reasoning only", args: []string{"--reasoning-effort", "high"}},
		{name: "agent and model", args: []string{"--agent", "codex", "--model", "gpt-5.6-sol"}},
		{name: "agent and reasoning", args: []string{"--agent", "codex", "--reasoning-effort", "high"}},
		{name: "model and reasoning", args: []string{"--model", "gpt-5.6-sol", "--reasoning-effort", "high"}},
	}
	for _, tt := range partials {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"resolve", "--pr", "123", "--no-input"}
			args = append(args, tt.args...)
			assertRunReviewSelectionOverrideRejectsPartialBeforeConfigLoad(t, args)
		})
	}
}

func TestRunWatchSelectionOverrideRejectsPartialBeforeConfigLoad(t *testing.T) {
	assertRunReviewSelectionOverrideRejectsPartialBeforeConfigLoad(t, []string{
		"watch", "--source", "coderabbit", "--pr", "123", "--model", "gpt-5.6-sol", "--no-input",
	})
}

func assertRunReviewSelectionOverrideRejectsPartialBeforeConfigLoad(t *testing.T, args []string) {
	t.Helper()
	homeDir, _ := withCLIWorkspace(t)
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	const invalidConfig = "defaults:\n  agent: [\n"
	mustMkdir(t, filepath.Dir(configPath))
	mustWrite(t, configPath, invalidConfig)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(args, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
	}
	const wantGrammar = "--agent, --model, and --reasoning-effort must be provided together for a one-Run Agent Selection override; omit all three to use Agent Selection Profiles"
	if !strings.Contains(stderr.String(), wantGrammar) {
		t.Fatalf("expected selection grammar error before config load, got %q", stderr.String())
	}
	if got := mustRead(t, configPath); got != invalidConfig {
		t.Fatalf("partial override changed User Config\nwant: %q\n got: %q", invalidConfig, got)
	}
	assertNoRunDatabase(t, homeDir)
	if _, err := os.Stat(filepath.Join(homeDir, ".roundfix", "worktrees")); err == nil {
		t.Fatalf("partial override created a Run Worktree root")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Worktree root: %v", err)
	}
}

func TestRunResolveACPXProbeFailureReportsActionablePreflight(t *testing.T) {
	installCommand := "npm install -g acpx@" + agent.MinimumACPXVersion
	tests := []struct {
		name     string
		runner   agent.Runner
		contains []string
	}{
		{
			name:   "missing binary",
			runner: &agent.ACPXRunner{Command: filepath.Join(t.TempDir(), "missing-acpx")},
			contains: []string{
				"acpx is required but was not found on PATH",
				installCommand,
			},
		},
		{
			name:   "version mismatch",
			runner: &agent.ACPXRunner{Command: fakeACPXVersionCommand(t, "0.11.0")},
			contains: []string{
				"found 0.11.0",
				"requires " + agent.MinimumACPXVersion,
				"upgrade",
				installCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			withAgentRunner(t, tt.runner)
			withFakeWorktree(t)
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected probe preflight exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			for _, want := range tt.contains {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			if strings.Count(stderr.String(), installCommand) != 1 {
				t.Fatalf("expected one actionable install command, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "Roundfix did not create a Run") {
				t.Fatalf("expected no-side-effects line, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveAgentFailureMarksBatchFailed(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{runErr: errors.New("agent crashed")})
	withFakeWorktree(t)
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	if !strings.Contains(stdout.String(), "issue 001 failed — major: handle test issue — reason: Agent failed") ||
		!strings.Contains(stdout.String(), "agent crashed") ||
		!strings.Contains(stdout.String(), "This Run (Unresolved after 1 Round(s)): 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.") ||
		!strings.Contains(stdout.String(), "Pull Request cumulative: 0 resolved, 0 invalid, 0 duplicated, 1 failed, 0 unresolved.") {
		t.Fatalf("expected failed stdout report with reason, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "agent crashed") {
		t.Fatalf("expected Agent failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome after Agent failure, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveAgentFailureContinuesWithLaterBatches(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{runErr: errors.New("agent crashed")}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	result := persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Unresolved exit 1, got %d", code)
	}
	if runner.calls != 2 {
		t.Fatalf("expected later Batches to run after the first failed, got %d Agent calls", runner.calls)
	}
	first, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(result.IssuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusFailed {
		t.Fatalf("expected first Batch issue failed, got %q", first.Status)
	}
	if second.Status != rounds.StatusFailed {
		t.Fatalf("expected second Batch issue failed after its own Agent failure, got %q", second.Status)
	}
	if !strings.Contains(stderr.String(), "Batch 001 failed") || !strings.Contains(stderr.String(), "Batch 002 failed") {
		t.Fatalf("expected both Batch failures reported, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveUsesOneAgentSessionPerRunAndCloses(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	inner := &fakeAgentRunner{}
	runner := &sessionRecordingRunner{inner: inner}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit 0, got %d stderr=%q", code, stderr.String())
	}
	if inner.calls != 2 {
		t.Fatalf("expected two Agent calls for two Batches, got %d", inner.calls)
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	assertRecordedOneSessionForRun(t, runner, runID, inner.calls)
	assertSessionLifecycleEvents(t, events)
}

func TestRunResolveClosesAgentSessionForTerminalOutcomes(t *testing.T) {
	tests := []struct {
		name      string
		inner     agent.Runner
		pusherErr error
		wantCode  int
		wantState string
		closeErr  error
	}{
		{
			name:      "clean",
			inner:     &fakeAgentRunner{},
			wantCode:  1,
			wantState: store.StateFailed,
			closeErr:  errors.New("close failed"),
		},
		{
			name:      "unresolved",
			inner:     &fakeAgentRunner{runErr: errors.New("agent crashed")},
			wantCode:  1,
			wantState: store.StateUnresolved,
		},
		{
			name:      "failed",
			inner:     &fakeAgentRunner{},
			pusherErr: errors.New("push failed"),
			wantCode:  1,
			wantState: store.StateFailed,
		},
		{
			name:      "stopped",
			inner:     &fakeStoppingAgentRunner{},
			wantCode:  0,
			wantState: store.StateStopped,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &sessionRecordingRunner{inner: tt.inner, closeErr: tt.closeErr}
			withAgentRunner(t, runner)
			withFakeWorktree(t)
			if tt.pusherErr != nil {
				withPusher(t, &fakePusher{err: tt.pusherErr})
			}
			if tt.wantState == store.StateStopped {
				withChangedPaths(t, nil)
			}
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q", tt.wantCode, code, stderr.String())
			}
			runID, _ := journaledRunEvents(t, homeDir, stderr.String())
			run := runFromStore(t, homeDir, runID)
			if run.State != tt.wantState {
				t.Fatalf("expected Run state %q, got %q", tt.wantState, run.State)
			}
			assertRecordedOneSessionForRun(t, runner, runID, len(runner.runSessions))
		})
	}
}

func TestRunResolveRejectsIncompatibleArtifacts(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "other-branch")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Head Repository \"owner/project\", PR Head Branch \"feature/review\"") {
		t.Fatalf("expected incompatible artifact context, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunOperationalCommandRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "missing pull request with no input",
			args:     []string{"fetch", "--source", "coderabbit", "--no-input"},
			contains: "missing required --pr because --no-input disables Interactive Input",
		},
		{
			name:     "invalid pull request number",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "abc", "--no-input"},
			contains: "--pr must be a positive integer",
		},
		{
			name:     "unsupported source",
			args:     []string{"fetch", "--source", "other", "--pr", "123", "--no-input"},
			contains: "unsupported Review Source",
		},
		{
			name:     "unsupported agent",
			args:     []string{"resolve", "--pr", "123", "--agent", "other", "--model", "other-model", "--reasoning-effort", "high", "--no-input"},
			contains: "unsupported Agent",
		},
		{
			name:     "invalid max rounds",
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--max-rounds", "0", "--no-input"},
			contains: "--max-rounds must be greater than 0",
		},
		{
			name:     "resolve agent console suppression conflicts with interactive input",
			args:     []string{"resolve", "--pr", "123", "--interactive", "--no-agent-console"},
			contains: "--interactive cannot be used with --no-agent-console",
		},
		{
			name:     "watch agent console suppression conflicts with interactive input",
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--interactive", "--no-agent-console"},
			contains: "--interactive cannot be used with --no-agent-console",
		},
		{
			name:     "unknown flag",
			args:     []string{"fetch", "--unknown"},
			contains: "flag provided but not defined",
		},
		{
			name:     "unexpected argument",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "123", "extra"},
			contains: "unexpected argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.contains) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.contains, stderr.String())
			}
			if !strings.Contains(stderr.String(), "did not create a Run") {
				t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunNoAgentConsoleRejectsInteractiveCockpit(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--no-agent-console", "--no-input"},
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--no-agent-console", "--no-input"},
		},
		{
			name: "implement",
			args: []string{"implement", "--spec", "0001-widget-flow", "--no-agent-console", "--no-input"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			t.Setenv("ROUNDFIX_TUI", "always")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "--no-agent-console cannot be used with the interactive cockpit") {
				t.Fatalf("expected interactive cockpit conflict, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunFetchCollectsMissingPullRequestWithInteractiveInput(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withCurrentPullRequestSuggestion(t, "321")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch exit 0, got %d", code)
	}
	if inputReq.Command != "fetch" {
		t.Fatalf("expected fetch Interactive Input request, got %#v", inputReq)
	}
	if inputReq.PRSuggestion.Value != "321" || inputReq.PRSuggestion.Source != "current" {
		t.Fatalf("expected current PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if !strings.Contains(stderr.String(), "Interactive Input collected command parameters.") {
		t.Fatalf("expected Interactive Input confirmation, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Roundfix fetch") || !strings.Contains(stderr.String(), "State: FetchingIssues") {
		t.Fatalf("expected Live Run View output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "321" {
		t.Fatalf("expected remembered PR 321, got %#v", defaults)
	}
}

func TestRunFetchReportsBlankInteractivePullRequestWithoutSuggestingInteractiveAgain(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withCurrentPullRequestSuggestion(t, "")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		return req.Values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Interactive Input did not collect required --pr") {
		t.Fatalf("expected blank Interactive Input guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "use --interactive") {
		t.Fatalf("did not expect guidance to reopen Interactive Input, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveInteractiveInputSuggestsConfiguredAgentAndRememberedPullRequest(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if err := runStore.RememberInteractiveDefaults(context.Background(), store.InteractiveDefaults{PRNumber: "123", Agent: "claude"}); err != nil {
		t.Fatalf("remember defaults: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withCurrentPullRequestSuggestion(t, "")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--interactive"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected resolve exit 0, got %d", code)
	}
	if inputReq.PRSuggestion.Value != "123" || inputReq.PRSuggestion.Source != "remembered" {
		t.Fatalf("expected remembered PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if inputReq.AgentSuggestion.Value != "codex" || inputReq.AgentSuggestion.Source != "config" {
		t.Fatalf("expected configured Agent suggestion before remembered value, got %#v", inputReq.AgentSuggestion)
	}
	if !strings.Contains(stderr.String(), "State: ResolvingWithAgent") {
		t.Fatalf("expected resolving Live Run View, got %q", stderr.String())
	}
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "123" || defaults.Agent != "codex" {
		t.Fatalf("expected remembered PR and configured Agent after run, got %#v", defaults)
	}
}

func TestRunInteractiveInputRunsPreflightBeforeSideEffects(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("fetch must not run after post-input Preflight Validation failure")
		return nil, nil
	})
	withCurrentPullRequestSuggestion(t, "123")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("post-input preflight failed")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected post-input preflight exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "post-input preflight failed") {
		t.Fatalf("expected preflight failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
	_ = repoDir
}

func TestRunOperationalCommandAppliesConfigAndCLIArtifactDirPrecedence(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  artifact_dir: user-artifacts
`)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
defaults:
  artifact_dir: project-artifacts
`)

	var projectStdout bytes.Buffer
	var projectStderr bytes.Buffer
	projectCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &projectStdout, &projectStderr)

	if projectCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", projectCode)
	}
	if !strings.Contains(projectStdout.String(), filepath.Join(repoDir, "project-artifacts")) {
		t.Fatalf("expected project artifact dir to win over user config, got %q", projectStdout.String())
	}

	var cliStdout bytes.Buffer
	var cliStderr bytes.Buffer
	cliCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", "cli-artifacts", "--no-input"}, &cliStdout, &cliStderr)

	if cliCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", cliCode)
	}
	if !strings.Contains(cliStdout.String(), filepath.Join(repoDir, "cli-artifacts")) {
		t.Fatalf("expected CLI artifact dir to win over project config, got %q", cliStdout.String())
	}
}

func TestRunFetchRejectsDuplicateActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	subdir := filepath.Join(repoDir, "nested")
	mustMkdir(t, subdir)
	t.Chdir(subdir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2 for duplicate Active Run, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Active Run "+active.ID) ||
		!strings.Contains(stderr.String(), "roundfix stop "+active.ID) ||
		!strings.Contains(stderr.String(), "roundfix stop --force "+active.ID) {
		t.Fatalf("expected existing run id in stderr, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestBranchIntegrityPreflightRejectsPendingRunBranchForReviewCommands(t *testing.T) {
	for _, command := range []string{"fetch", "resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			pending := []runworktree.PendingRunWork{{
				Branch:       "roundfix/run-stranded",
				WorktreePath: filepath.Join(repoDir, "..", "run-stranded"),
				AheadCommits: 2,
				FastForward:  false,
			}}
			recorder := withBranchIntegrity(t, pending, nil)
			withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
				t.Fatal("fetch must not run after Branch Integrity Preflight refusal")
				return nil, nil
			})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(branchIntegrityCommandArgs(command), &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected Branch Integrity Preflight exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			for _, want := range []string{
				"Branch Integrity Preflight refused pending Run Branch work",
				"roundfix/run-stranded",
				"ahead_commits=2",
				`integration_command="git merge --ff-only roundfix/run-stranded"`,
				"did not create a Run",
			} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			if len(recorder.integrations) != 0 {
				t.Fatalf("expected no integration attempts on refusal, got %#v", recorder.integrations)
			}
			assertRunCount(t, store.DatabasePath(homeDir), 0)
		})
	}
}

func TestBranchIntegrityPreflightIntegratesFastForwardRunBranchAndJournals(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pending := []runworktree.PendingRunWork{{
		Branch:       "roundfix/run-ready",
		WorktreePath: filepath.Join(repoDir, "..", "run-ready"),
		AheadCommits: 1,
		FastForward:  true,
	}}
	recorder := withBranchIntegrity(t, pending, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch to proceed after auto-integration, got %d stderr=%q", code, stderr.String())
	}
	if len(recorder.integrations) != 1 || recorder.integrations[0] != "roundfix/run-ready" {
		t.Fatalf("expected one integration of roundfix/run-ready, got %#v", recorder.integrations)
	}
	if !strings.Contains(stderr.String(), "Branch Integrity Preflight: integrated roundfix/run-ready") ||
		!strings.Contains(stderr.String(), "git merge --ff-only roundfix/run-ready") {
		t.Fatalf("expected integration report, got %q", stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	if !journalPayloadContains(events, "branch_integrity_auto_integration") {
		t.Fatalf("expected auto-integration Run Event, got %+v", events)
	}
	if !strings.Contains(stdout.String(), "Fetch complete") {
		t.Fatalf("expected fetch success, got %q", stdout.String())
	}
}

// seedOutdatedV9RunDatabase downgrades a freshly created Run Database to
// schema v9 by reversing the v9→v10 migration, matching what a historical
// binary left behind.
func seedOutdatedV9RunDatabase(t *testing.T, homeDir string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("create Run Database: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open Run Database for downgrade: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close downgraded Run Database: %v", err)
		}
	}()
	for _, statement := range []string{
		`ALTER TABLE runs DROP COLUMN owner_identity`,
		`PRAGMA user_version = 9`,
	} {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("downgrade Run Database to v9: %v", err)
		}
	}
}

func TestBranchIntegrityPreflightMigratesOutdatedRunDatabase(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withBranchIntegrity(t, nil, nil)
	seedOutdatedV9RunDatabase(t, homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected fetch to migrate the outdated Run Database and proceed, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Fetch complete") {
		t.Fatalf("expected fetch success, got %q", stdout.String())
	}
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("expected migrated Run Database to open read-only, got %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("close reader after migration: %v", err)
		}
	}()
	version, err := reader.MigrationVersion(context.Background())
	if err != nil {
		t.Fatalf("read migrated schema version: %v", err)
	}
	if version != 10 {
		t.Fatalf("expected schema version 10 after preflight migration, got %d", version)
	}
}

func TestBranchIntegrityPreflightRejectsActiveRunForReviewCommands(t *testing.T) {
	for _, command := range []string{"fetch", "resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			withBranchIntegrity(t, nil, nil)
			runStore, err := store.Open(context.Background(), homeDir)
			if err != nil {
				t.Fatalf("open store: %v", err)
			}
			active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
				Kind:           store.KindWatch,
				HeadRepository: "owner/project",
				HeadBranch:     "feature/review",
				BaseRepository: "owner/project",
				PRNumber:       "123",
				GitRoot:        repoDir,
				LocalBranch:    "feature/review",
				HeadSHA:        "abc123",
				ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
				Agent:          "codex",
			})
			if err != nil {
				t.Fatalf("create active Run: %v", err)
			}
			if err := runStore.Close(); err != nil {
				t.Fatalf("close store: %v", err)
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(branchIntegrityCommandArgs(command), &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected active Run refusal exit 2, got %d", code)
			}
			for _, want := range []string{
				active.ID,
				"roundfix stop " + active.ID,
				"roundfix stop --force " + active.ID,
				"did not create a Run",
			} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			assertRunCount(t, store.DatabasePath(homeDir), 1)
		})
	}
}

func TestBranchIntegrityBypassPublishesAuditBeforeFetch(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pending := []runworktree.PendingRunWork{{
		Branch:       "roundfix/run-ignored",
		WorktreePath: filepath.Join(repoDir, "..", "run-ignored"),
		AheadCommits: 3,
		FastForward:  false,
	}}
	recorder := withBranchIntegrity(t, pending, nil)
	comments := withPullRequestComments(t, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--skip-branch-integrity", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected bypassed fetch to proceed, got %d stderr=%q", code, stderr.String())
	}
	if len(recorder.integrations) != 0 {
		t.Fatalf("expected bypass to skip integrations, got %#v", recorder.integrations)
	}
	if len(comments.calls) != 1 {
		t.Fatalf("expected one bypass audit comment, got %#v", comments.calls)
	}
	if comments.calls[0].repository != "owner/project" {
		t.Fatalf("expected audit posted to the PR's Base Repository, got %q", comments.calls[0].repository)
	}
	body := comments.calls[0].body
	for _, want := range []string{
		"<!-- roundfix: run:",
		"bypass:branch-integrity",
		"pr:123",
		"Run:",
		"Skipped guardrails:",
		"Pending Run Branch work",
		"Active Run bound to target",
		"roundfix/run-ignored",
		"ahead_commits=3",
		`integration_command="git merge --ff-only roundfix/run-ignored"`,
		"Ignored Active Runs:\n- none",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected audit body to contain %q, got %q", want, body)
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	if !journalPayloadContains(events, "branch_integrity_bypass") {
		t.Fatalf("expected bypass Run Event, got %+v", events)
	}
}

func TestBranchIntegrityBypassAuditsActiveRunAndProceeds(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withBranchIntegrity(t, nil, nil)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		Agent:          "codex",
	})
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	comments := withPullRequestComments(t, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--skip-branch-integrity", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected bypassed fetch to proceed despite active Run, got %d stderr=%q", code, stderr.String())
	}
	if len(comments.calls) != 1 {
		t.Fatalf("expected one bypass audit comment, got %#v", comments.calls)
	}
	for _, want := range []string{
		"run_id=" + active.ID,
		`stop_command="roundfix stop ` + active.ID + `"`,
		`force_stop_command="roundfix stop --force ` + active.ID + `"`,
	} {
		if !strings.Contains(comments.calls[0].body, want) {
			t.Fatalf("expected audit body to contain %q, got %q", want, comments.calls[0].body)
		}
	}
	assertRunCount(t, store.DatabasePath(homeDir), 2)
	assertActiveRun(t, homeDir, "owner/project", "feature/review", active.ID)
}

func TestBranchIntegrityBypassFailsWhenAuditCommentPublishFails(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withBranchIntegrity(t, []runworktree.PendingRunWork{{
		Branch:       "roundfix/run-ignored",
		WorktreePath: filepath.Join(repoDir, "..", "run-ignored"),
		AheadCommits: 1,
		FastForward:  false,
	}}, nil)
	withPullRequestComments(t, errors.New("comment denied"))
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("fetch must not run when bypass audit publish fails")
		return nil, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--skip-branch-integrity", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected audit publish failure exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, want := range []string{
		"publish Branch Integrity Preflight bypass audit",
		"comment denied",
		"did not fetch Review Source issues",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestReviewCommandHelpDocumentsSkipBranchIntegrity(t *testing.T) {
	for _, command := range []string{"fetch", "resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{command, "--help"}, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected help exit 0, got %d", code)
			}
			if !strings.Contains(stdout.String(), "skip-branch-integrity") {
				t.Fatalf("expected help to document skip-branch-integrity, got %q", stdout.String())
			}
		})
	}
}

func TestRunStopByRunIDRecordsStopRequest(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", active.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	for _, expected := range []string{active.ID, "Stop Request recorded; the Run stops after the current Work Item settles.", "--force"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("expected stop output to contain %q, got %q", expected, stdout.String())
		}
	}
	assertStopRequested(t, homeDir, active.ID)
	assertActiveRun(t, homeDir, "owner/project", "feature/review", active.ID)
}

func TestRunStopByPullRequestRecordsStopRequestForMatchingActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withStopPullRequestResolver(t, preflight.PullRequest{
		Number:         "123",
		State:          "OPEN",
		BaseRepository: "owner/project",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
	})

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--pr", "123"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), active.ID) {
		t.Fatalf("expected stopped run id in stdout, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Stop Request recorded; the Run stops after the current Work Item settles.") {
		t.Fatalf("expected Stop Request report in stdout, got %q", stdout.String())
	}
	assertStopRequested(t, homeDir, active.ID)
	assertActiveRun(t, homeDir, "owner/project", "feature/review", active.ID)
}

func TestRunStopBySpecRecordsStopRequestForMatchingActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	request := store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/implement-spec",
		SpecSlug:    "0001-widget-flow",
	}
	active, err := runStore.CreateRun(ctx, request)
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--spec", "0001-widget-flow"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), active.ID) ||
		!strings.Contains(stdout.String(), "Stop Request recorded; the Run stops after the current Work Item settles.") ||
		!strings.Contains(stdout.String(), "--force") {
		t.Fatalf("expected Stop Request report for implement run in stdout, got %q", stdout.String())
	}
	runStore, err = store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close reopened store: %v", err)
		}
	}()
	current, found, err := runStore.Run(ctx, active.ID)
	if err != nil || !found {
		t.Fatalf("lookup Run: found=%v err=%v", found, err)
	}
	if current.State != store.StateActive {
		t.Fatalf("expected Run to stay Active until the engine settles, got %q", current.State)
	}
	requested, err := runStore.StopRequested(ctx, active.ID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected Stop Request flag recorded")
	}
	if _, err := runStore.CreateRun(ctx, request); err == nil {
		t.Fatal("expected Stop Request to keep the Active Run lock until engine completion")
	}
}

func TestRunForceStopOwnerExitPrecedesCompletionAndLockRelease(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	proved := false
	withOwnerProcessController(t, ownerProcessControllerFunc(func(ctx context.Context, pid int, recordedIdentity string) error {
		if active.OwnerPID == nil || pid != *active.OwnerPID {
			t.Fatalf("owner process PID = %d, want recorded PID %v", pid, active.OwnerPID)
		}
		if recordedIdentity != active.OwnerIdentity {
			t.Fatalf("owner identity = %q, want recorded identity %q", recordedIdentity, active.OwnerIdentity)
		}
		runStore, err := store.Open(ctx, homeDir)
		if err != nil {
			t.Fatalf("open store during owner exit proof: %v", err)
		}
		defer func() {
			if err := runStore.Close(); err != nil {
				t.Fatalf("close store during owner exit proof: %v", err)
			}
		}()
		current, found, err := runStore.Run(ctx, active.ID)
		if err != nil || !found {
			t.Fatalf("read Run during owner exit proof: found=%v err=%v", found, err)
		}
		if current.State != store.StateActive {
			t.Fatalf("Run state during owner exit proof = %q, want Active", current.State)
		}
		if _, err := runStore.CreateRun(ctx, request); err == nil {
			t.Fatal("Active Run lock was released before owner exit proof")
		}
		proved = true
		return nil
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("force stop exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !proved {
		t.Fatal("owner exit was not proven")
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop success report, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunForceStopOwnerPermissionAndDeadlineFailuresRetainActiveLock(t *testing.T) {
	tests := []struct {
		name string
		step string
		err  error
	}{
		{
			name: "permission",
			step: "send graceful termination",
			err:  errors.New("operation not permitted"),
		},
		{
			name: "deadline",
			step: "prove exit after force kill",
			err:  context.DeadlineExceeded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
			withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
				return store.OwnerProcessControlError{PID: *active.OwnerPID, Step: tt.step, Err: tt.err}
			}))
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

			if code != exitRunFailed {
				t.Fatalf("force stop exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("failed force stop printed success output %q", stdout.String())
			}
			for _, want := range []string{
				active.ID,
				strconv.Itoa(*active.OwnerPID),
				tt.step,
				"remains Active",
				"Active Run lock retained",
			} {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("force stop diagnostic missing %q: %q", want, stderr.String())
				}
			}
			assertRunState(t, homeDir, active.ID, store.StateActive)
			runStore, err := store.Open(context.Background(), homeDir)
			if err != nil {
				t.Fatalf("open store after failed force stop: %v", err)
			}
			defer func() {
				if err := runStore.Close(); err != nil {
					t.Fatalf("close store after failed force stop: %v", err)
				}
			}()
			if _, err := runStore.CreateRun(context.Background(), request); err == nil {
				t.Fatal("failed force stop released the Active Run lock")
			}
		})
	}
}

func TestRunForceStopOwnerProofFailurePreservesAgentSessions(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	cancelCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return nil
	})
	closeCalls := 0
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		closeCalls++
		return nil
	})
	terminateCalls := 0
	withOwnerProcessController(t, ownerProcessControllerStub{
		prove: func(context.Context, int, string) error {
			return store.OwnerProcessControlError{
				PID:  *active.OwnerPID,
				Step: "prove owner process identity",
				Err:  store.ErrOwnerProcessIdentityUnproven,
			}
		},
		terminate: func(context.Context, int, string) error {
			terminateCalls++
			return nil
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("Force Stop failure exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("Force Stop failure printed success output %q", stdout.String())
	}
	if cancelCalls != 0 || closeCalls != 0 {
		t.Fatalf("failed owner proof mutated Agent Sessions: cancel=%d close=%d", cancelCalls, closeCalls)
	}
	if terminateCalls != 0 {
		t.Fatalf("failed owner proof still terminated the owner: %d calls", terminateCalls)
	}
	if strings.Contains(stderr.String(), "Secondary cleanup warning:") {
		t.Fatalf("failed owner proof must report no cleanup warnings, got %q", stderr.String())
	}
	for _, want := range []string{
		active.ID,
		strconv.Itoa(*active.OwnerPID),
		"prove owner process identity",
		"remains Active",
		"Active Run lock retained",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("Force Stop diagnostic missing %q: %q", want, stderr.String())
		}
	}
	assertRunState(t, homeDir, active.ID, store.StateActive)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store after failed Force Stop: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store after failed Force Stop: %v", err)
		}
	}()
	activeScopes, err := runStore.ActiveAgentSelectionScopes(context.Background(), active.ID)
	if err != nil {
		t.Fatalf("read active Agent Selection scopes: %v", err)
	}
	if len(activeScopes) != 1 || activeScopes[0].Status != store.AgentSelectionStatusActive {
		t.Fatalf("Agent Selection lifecycle changed after owner failure: %+v", activeScopes)
	}
	if _, err := runStore.CreateRun(context.Background(), request); err == nil {
		t.Fatal("failed owner proof released the Active Run lock")
	}
}

func TestRunForceStopPrimaryFailurePrecedesSecondaryCleanupWarnings(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, _ := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return errors.New("cancel denied")
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return errors.New("close denied")
	})
	// The owner proof succeeds, so Force Stop reaches session cleanup; the
	// termination that follows fails and must still surface the cleanup
	// warnings it accumulated, behind the primary failure.
	withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
		return store.OwnerProcessControlError{
			PID:  *active.OwnerPID,
			Step: "prove owner exit",
			Err:  context.DeadlineExceeded,
		}
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("Force Stop failure exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("Force Stop failure printed success output %q", stdout.String())
	}
	primaryIndex := strings.Index(stderr.String(), "Stop failed")
	secondaryIndex := strings.Index(stderr.String(), "Secondary cleanup warning:")
	if primaryIndex < 0 || secondaryIndex < 0 || primaryIndex >= secondaryIndex {
		t.Fatalf("primary failure must precede labeled secondary cleanup warnings, got %q", stderr.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateActive)

	events := runEvents(t, homeDir, active.ID)
	primaryCursor := int64(0)
	secondaryCursor := int64(0)
	for _, entry := range events {
		switch {
		case strings.Contains(entry.Event.Summary, "Force Stop primary failure"):
			primaryCursor = entry.Cursor
		case strings.Contains(entry.Event.Summary, "Secondary cleanup warning"):
			if secondaryCursor == 0 {
				secondaryCursor = entry.Cursor
			}
		}
	}
	if primaryCursor == 0 || secondaryCursor == 0 || primaryCursor >= secondaryCursor {
		t.Fatalf("journaled primary failure must precede secondary cleanup warnings, got %+v", events)
	}
}

// TestRunForceStopOwnerPIDReuseFailsClosed lives in orphan_unix_test.go: it
// proves the real identity comparison against a live scratch process whose
// stored owner identity token differs from its genuine one.

func TestRunForceStopStoppedRunIsIdempotentWithoutOwnerOrSessionActions(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, _ := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store to complete Run: %v", err)
	}
	if _, err := runStore.CompleteRun(context.Background(), active.ID, store.StateStopped); err != nil {
		t.Fatalf("complete Run as Stopped: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store after completion: %v", err)
	}
	cancelCalls := 0
	closeCalls := 0
	ownerCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		closeCalls++
		return nil
	})
	withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
		ownerCalls++
		return nil
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("idempotent force stop exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if cancelCalls != 0 || closeCalls != 0 || ownerCalls != 0 {
		t.Fatalf("idempotent force stop actions: cancel=%d close=%d owner=%d", cancelCalls, closeCalls, ownerCalls)
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected idempotent force stop report, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunStopForceAgentSessionCleanupSkipsRunWithoutActiveLifecycle(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	cancelCalls := 0
	closeCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		closeCalls++
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", "--spec", "0001-widget-flow"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "Agent Session") || strings.Contains(stderr.String(), "session roundfix-") {
		t.Fatalf("expected no Agent Session diagnostics without active lifecycle evidence, got %q", stderr.String())
	}
	if cancelCalls != 0 || closeCalls != 0 {
		t.Fatalf("expected zero Agent Session calls, got cancel=%d close=%d", cancelCalls, closeCalls)
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") || !strings.Contains(stdout.String(), active.ID) {
		t.Fatalf("expected force stop report with Run ID, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	next, err := runStore.CreateRun(context.Background(), request)
	if err != nil {
		t.Fatalf("expected force stop to release Spec lock, got %v", err)
	}
	if next.ID == active.ID {
		t.Fatal("expected a new Run after force stop")
	}
}

func TestRunStopForceRegisteredAgentSessionCleanupTargetsActiveScopesInOrder(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", "")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeReview, ScopeID: "batch-002",
		Category: "review", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "claude", Model: "review-model",
		Status: store.AgentSelectionStatusActive,
	})
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeQA, ScopeID: "qa",
		Category: "qa", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRoleFallback, FallbackIndex: 1,
		Runtime: "opencode", Model: "qa-model", Status: store.AgentSelectionStatusActive,
	})
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	calls := []string{}
	withStopAgentSessionCanceler(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		calls = append(calls, "cancel "+runtime.ID+" "+session.Name+" "+session.WorkDir)
		return nil
	})
	withStopAgentSessionCloser(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		calls = append(calls, "close "+runtime.ID+" "+session.Name+" "+session.WorkDir)
		return nil
	})
	withOwnerProcessController(t, ownerProcessControllerStub{
		prove: func(_ context.Context, pid int, _ string) error {
			calls = append(calls, fmt.Sprintf("prove %d", pid))
			return nil
		},
		terminate: func(_ context.Context, pid int, _ string) error {
			calls = append(calls, fmt.Sprintf("owner %d", pid))
			return nil
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	// The read-only owner proof runs first, then every registered Agent
	// Session is cancelled, and only then is the owner terminated.
	wantCalls := []string{
		fmt.Sprintf("prove %d", *active.OwnerPID),
		"cancel codex roundfix-" + active.ID + "-task_01 " + taskRef.Path,
		"close codex roundfix-" + active.ID + "-task_01 " + taskRef.Path,
		"cancel opencode roundfix-" + active.ID + "-qa-fallback-01 " + active.WorkDir,
		"close opencode roundfix-" + active.ID + "-qa-fallback-01 " + active.WorkDir,
		"cancel claude roundfix-" + active.ID + "-review-002 " + repoDir,
		"close claude roundfix-" + active.ID + "-review-002 " + repoDir,
		fmt.Sprintf("owner %d", *active.OwnerPID),
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("unexpected session invocations\nwant: %#v\ngot:  %#v", wantCalls, calls)
	}
	if strings.Contains(stderr.String(), "Agent Session") || strings.Contains(stderr.String(), "session roundfix-") {
		t.Fatalf("expected successful registered cleanup to stay silent, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeTask, "task_01", store.AgentSelectionStatusClosed)
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeQA, "qa", store.AgentSelectionStatusClosed)
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeReview, "batch-002", store.AgentSelectionStatusClosed)
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunStopForceRegisteredAgentSessionAbsenceIsIdempotent(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, _ := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	absent := fmt.Errorf("adapter response: %w", &agent.InfrastructureError{ExitCode: 4, Reason: "missing session"})
	cancelCalls := 0
	closeCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return absent
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		closeCalls++
		return absent
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if cancelCalls != 1 || closeCalls != 1 {
		t.Fatalf("expected one cancel and close attempt, got cancel=%d close=%d", cancelCalls, closeCalls)
	}
	if strings.Contains(stderr.String(), "Agent Session") || strings.Contains(stderr.String(), "session roundfix-") {
		t.Fatalf("expected registered Agent Session absence to stay silent, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeTask, "task_01", store.AgentSelectionStatusClosed)
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunStopForceAgentSessionCleanupFailureRemainsVisibleWithoutClosedLifecycle(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, _ := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return errors.New("cancel denied")
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return errors.New("close denied")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{
		"Agent Session cancel failed for task task_01",
		"cancel denied",
		"Agent Session close failed for task task_01",
		"close denied",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected cleanup diagnostic %q, got %q", want, stderr.String())
		}
	}
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeTask, "task_01", store.AgentSelectionStatusActive)
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunStopForceReapsEmptyRunAndTaskWorktrees(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", "")
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
		return nil
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	for _, expected := range []string{
		"roundfix: reaped terminal Worktree path=" + active.WorkDir + " branch=" + runworktree.BranchName(active.ID),
		"roundfix: reaped terminal Worktree path=" + taskRef.Path + " branch=" + taskRef.Branch,
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	assertRunWorktreeRemoved(t, active.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(active.ID))
	assertRunWorktreeRemoved(t, taskRef.Path)
	assertRunBranchRemoved(t, repoDir, taskRef.Branch)
}

func TestRunStopForceKeepsRunWorktreeWithCommits(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, runRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", "")
	if err := os.WriteFile(filepath.Join(runRef.Path, "kept.txt"), []byte("settled work\n"), 0o644); err != nil {
		t.Fatalf("write kept work: %v", err)
	}
	gitImplement(t, runRef.Path, "add", "kept.txt")
	gitImplement(t, runRef.Path, "commit", "-m", "settled work")
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	withOwnerProcessController(t, ownerProcessControllerFunc(func(context.Context, int, string) error {
		return nil
	}))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "reaped terminal Worktree") {
		t.Fatalf("expected no debris reap report for a Run Branch with commits, got %q", stderr.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	assertRunWorktreeExists(t, active.WorkDir)
	assertRunBranchExists(t, repoDir, runworktree.BranchName(active.ID))
}

func TestRunStopForceReportsCancelFailuresButCompletes(t *testing.T) {
	tests := []struct {
		name          string
		runtime       string
		cancelErr     error
		wantStderr    []string
		wantCanceller bool
		wantStatus    store.AgentSelectionStatus
	}{
		{
			name:          "missing acpx",
			runtime:       "codex",
			cancelErr:     errors.New("acpx not found"),
			wantStderr:    []string{"Agent Session cancel failed", "acpx not found"},
			wantCanceller: true,
			wantStatus:    store.AgentSelectionStatusClosed,
		},
		{
			name:          "unknown runtime",
			runtime:       "gemini",
			wantStderr:    []string{"Agent Session cleanup skipped", `unsupported Agent "gemini"`},
			wantCanceller: false,
			wantStatus:    store.AgentSelectionStatusActive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
			recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
				RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
				Category: "backend", ProfileSource: "project", Attempt: 1,
				SelectionRole: store.AgentSelectionRolePreferred, Runtime: tt.runtime, Model: "task-model",
				Status: store.AgentSelectionStatusActive,
			})
			cancelCalls := 0
			withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
				cancelCalls++
				return tt.cancelErr
			})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
			}
			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			if tt.wantCanceller && cancelCalls != 1 {
				t.Fatalf("expected one cancel attempt, got %d", cancelCalls)
			}
			if !tt.wantCanceller && cancelCalls != 0 {
				t.Fatalf("expected no cancel attempt, got %d", cancelCalls)
			}
			for _, want := range []string{
				"Roundfix Run force-stopped",
				"proved the recorded owner process exited",
				"released its Active Run locks",
			} {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("expected force stop report to contain %q, got %q", want, stdout.String())
				}
			}
			if strings.Contains(stdout.String(), "immediately") {
				t.Fatalf("force stop report must not promise immediate completion, got %q", stdout.String())
			}
			assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeTask, "task_01", tt.wantStatus)
			assertRunState(t, homeDir, active.ID, store.StateStopped)
			runStore, err := store.Open(context.Background(), homeDir)
			if err != nil {
				t.Fatalf("open store: %v", err)
			}
			_, err = runStore.CreateRun(context.Background(), request)
			closeErr := runStore.Close()
			if err != nil {
				t.Fatalf("expected force stop to release Spec lock after cancel failure, got %v", err)
			}
			if closeErr != nil {
				t.Fatalf("close store: %v", closeErr)
			}
		})
	}
}

func TestRunStopGracefulThenForceCompletesImmediately(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	recordAgentSelectionForStop(t, homeDir, store.AgentSelectionAttemptRequest{
		RunID: active.ID, ScopeKind: store.AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: store.AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
		Status: store.AgentSelectionStatusActive,
	})
	var gracefulStdout bytes.Buffer
	var gracefulStderr bytes.Buffer

	gracefulCode := Run([]string{"stop", "--spec", "0001-widget-flow"}, &gracefulStdout, &gracefulStderr)

	if gracefulCode != 0 {
		t.Fatalf("expected graceful stop exit 0, got %d stderr=%q", gracefulCode, gracefulStderr.String())
	}
	if !strings.Contains(gracefulStdout.String(), "Stop Request recorded") {
		t.Fatalf("expected graceful Stop Request report, got %q", gracefulStdout.String())
	}
	assertStopRequested(t, homeDir, active.ID)
	cancelCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return nil
	})
	var forceStdout bytes.Buffer
	var forceStderr bytes.Buffer

	forceCode := Run([]string{"stop", "--force", "--spec", "0001-widget-flow"}, &forceStdout, &forceStderr)

	if forceCode != 0 {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", forceCode, forceStderr.String())
	}
	if strings.Contains(forceStderr.String(), "Agent Session") || strings.Contains(forceStderr.String(), "session roundfix-") {
		t.Fatalf("expected successful registered cleanup to stay silent, got %q", forceStderr.String())
	}
	if cancelCalls != 1 {
		t.Fatalf("expected one cancel after graceful request, got %d", cancelCalls)
	}
	if !strings.Contains(forceStdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", forceStdout.String())
	}
	assertAgentSelectionScopeStatus(t, homeDir, active.ID, store.AgentSelectionScopeTask, "task_01", store.AgentSelectionStatusClosed)
	assertRunState(t, homeDir, active.ID, store.StateStopped)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	if _, err := runStore.CreateRun(context.Background(), request); err != nil {
		t.Fatalf("expected force after graceful request to release Spec lock, got %v", err)
	}
}

func TestRunStopSpecSelectorRejectsOtherSelectors(t *testing.T) {
	withCLIWorkspace(t)
	tests := []struct {
		name string
		args []string
	}{
		{name: "with pull request", args: []string{"stop", "--spec", "0001-widget-flow", "--pr", "123"}},
		{name: "with run id", args: []string{"stop", "--spec", "0001-widget-flow", "run_existing"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected stop validation exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "cannot be combined") {
				t.Fatalf("expected mutual-exclusion message, got %q", stderr.String())
			}
		})
	}
}

func TestRunStopBySpecReportsMissingActiveRun(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--spec", "0002-missing"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected stop validation exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	want := fmt.Sprintf("no Active Run exists for repository %q and Spec %q", repoDir, "0002-missing")
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected missing spec target %q, got %q", want, stderr.String())
	}
}

func TestRunStopHelpExplainsProofBeforeCompletion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop help exit 0, got %d", code)
	}
	for _, want := range []string{
		"roundfix stop --spec <slug>",
		"--spec",
		"--force",
		"graceful",
		"owner exit is proven",
		"Run remains Active",
		"Active Run lock stays retained",
		"roundfix runs list --state active",
		"roundfix stop --force <run-id>",
		"already Stopped",
		"different terminal outcome",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected stop help to contain %q, got %q", want, stdout.String())
		}
	}
	for _, forbidden := range []string{
		"Immediately stop a dead or runaway Run and release its lock",
		"completes the Run Stopped immediately",
	} {
		if strings.Contains(stdout.String(), forbidden) {
			t.Fatalf("stop help must not promise %q, got %q", forbidden, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunStopRejectsAlreadyTerminalRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, active.ID, store.StateClean); err != nil {
		t.Fatalf("complete run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", active.ID}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected stop validation exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "cannot record Stop Request for terminal Run") ||
		!strings.Contains(stderr.String(), "Clean") {
		t.Fatalf("expected named terminal Stop Request refusal, got %q", stderr.String())
	}
}

func TestRunPreflightFailureLeavesBufferOutputPlainByDefault(t *testing.T) {
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("plain preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if strings.Contains(stderr.String(), "\x1b[") {
		t.Fatalf("did not expect ANSI color in buffer output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Preflight failed") || !strings.Contains(stderr.String(), "plain preflight failure") {
		t.Fatalf("expected formatted preflight failure, got %q", stderr.String())
	}
}

func TestRunPreflightFailureColorsOutputWhenForced(t *testing.T) {
	t.Setenv("ROUNDFIX_COLOR", "always")
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("colored preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "\x1b[31m") {
		t.Fatalf("expected forced ANSI color in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "colored preflight failure") {
		t.Fatalf("expected original error text, got %q", stderr.String())
	}
}

func TestRunOperationalCommandRejectsInvalidConfigAndArtifactDirectory(t *testing.T) {
	t.Run("invalid user config YAML", func(t *testing.T) {
		homeDir, _ := withCLIWorkspace(t)
		mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
		mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), "defaults:\n  agent: [")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "parse config") {
			t.Fatalf("expected config parse error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})

	t.Run("artifact path is a file", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		artifactFile := filepath.Join(repoDir, "artifact-file")
		mustWrite(t, artifactFile, "not a directory")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", artifactFile, "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "is not a directory") {
			t.Fatalf("expected artifact dir error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})
}

type setupFakeDeps struct {
	homeDir           string
	gitRoot           string
	config            roundconfig.Config
	userConfigPath    string
	projectConfigPath string
	acpxConfigPath    string
	nodeVersion       string
	nodeErr           error
	acpxVersion       string
	acpxErr           error
	adapterEvidence   agent.AdapterEvidence
	adapterErr        error
	probeErr          error
	paths             map[string]string
	files             map[string]string
	installCalls      []string
	initScopes        []string
	acpxInitCalls     int
	writeCalls        []string
	writeErrors       map[string]error
	adapterRequests   []agent.RuntimeSpec
	probeRequests     []agent.ProbeRequest
	prompts           []string
	confirm           func(context.Context, io.Writer, string) (bool, error)
}

func newSetupFakeDeps() *setupFakeDeps {
	homeDir := "/home/roundfix-test"
	gitRoot := "/repo/project"
	return &setupFakeDeps{
		homeDir:           homeDir,
		gitRoot:           gitRoot,
		config:            roundconfig.Builtin(),
		userConfigPath:    filepath.Join(homeDir, ".roundfix", "config.yml"),
		projectConfigPath: filepath.Join(gitRoot, ".roundfixrc.yml"),
		acpxConfigPath:    filepath.Join(homeDir, ".acpx", "config.json"),
		nodeVersion:       "v25.6.1",
		acpxVersion:       agent.MinimumACPXVersion,
		adapterEvidence: agent.AdapterEvidence{
			Command: "npx -y " + agent.CodexAdapterPackage,
			Package: agent.CodexAdapterPackage,
			Version: agent.PinnedCodexAdapterVersion,
		},
		paths: map[string]string{},
		files: map[string]string{},
	}
}

func cloneSetupFiles(files map[string]string) map[string]string {
	cloned := make(map[string]string, len(files))
	for path, content := range files {
		cloned[path] = content
	}
	return cloned
}

func withSetupFakeDeps(t *testing.T, fake *setupFakeDeps) {
	t.Helper()
	old := setupDeps
	profileRunner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			fake.probeRequests = append(fake.probeRequests, req)
			if fake.probeErr != nil {
				return agent.SelectionProof{}, fake.probeErr
			}
			return agent.SelectionProof{
				Runtime:         req.Runtime.ID,
				Model:           req.Runtime.Model,
				ReasoningEffort: req.Runtime.ReasoningEffort,
				Assignment:      agent.SelectionAssignment{Encoding: agent.SelectionEncodingIndependent},
				Adapter:         fake.adapterEvidence,
				Status:          agent.SelectionProofStatusProven,
			}, nil
		},
	}
	setupDeps = setupDependencies{
		loadConfig: func(roundconfig.LoadOptions) (roundconfig.Loaded, error) {
			return roundconfig.Loaded{
				Config:            fake.config,
				GitRoot:           fake.gitRoot,
				HomeDir:           fake.homeDir,
				UserConfigPath:    fake.userConfigPath,
				ProjectConfigPath: fake.projectConfigPath,
			}, nil
		},
		nodeVersion: func(context.Context) (string, error) {
			return fake.nodeVersion, fake.nodeErr
		},
		acpxVersion: func(context.Context) (string, error) {
			return fake.acpxVersion, fake.acpxErr
		},
		installACPX: func(context.Context) error {
			fake.installCalls = append(fake.installCalls, "npm install -g acpx@"+agent.MinimumACPXVersion)
			return nil
		},
		checkAdapter: func(_ context.Context, runtime agent.RuntimeSpec) (agent.AdapterEvidence, error) {
			fake.adapterRequests = append(fake.adapterRequests, runtime)
			if fake.adapterErr != nil && strings.TrimSpace(runtime.Command) == "" {
				return agent.AdapterEvidence{}, fake.adapterErr
			}
			evidence := fake.adapterEvidence
			if strings.TrimSpace(runtime.Command) != "" {
				evidence.Command = runtime.Command
			}
			return evidence, nil
		},
		profileRunner: profileRunner,
		probeAgent: func(_ context.Context, req agent.ProbeRequest) error {
			fake.probeRequests = append(fake.probeRequests, req)
			return fake.probeErr
		},
		lookPath: func(command string) (string, error) {
			path := fake.paths[command]
			if path == "" {
				return "", os.ErrNotExist
			}
			return path, nil
		},
		exists: func(path string) (bool, error) {
			_, ok := fake.files[path]
			return ok, nil
		},
		readFile: func(path string) ([]byte, error) {
			content, ok := fake.files[path]
			if !ok {
				return nil, os.ErrNotExist
			}
			return []byte(content), nil
		},
		writeFile: func(path string, content []byte) error {
			if err := fake.writeErrors[path]; err != nil {
				fake.writeCalls = append(fake.writeCalls, path)
				return err
			}
			fake.files[path] = string(content)
			fake.writeCalls = append(fake.writeCalls, path)
			switch path {
			case fake.userConfigPath:
				fake.initScopes = append(fake.initScopes, roundconfig.InitScopeUser)
			case fake.projectConfigPath:
				fake.initScopes = append(fake.initScopes, roundconfig.InitScopeProject)
			}
			return nil
		},
		mkdirAll: func(string) error {
			return nil
		},
		confirm: func(ctx context.Context, stderr io.Writer, prompt string) (bool, error) {
			fake.prompts = append(fake.prompts, prompt)
			if fake.confirm != nil {
				return fake.confirm(ctx, stderr, prompt)
			}
			return true, nil
		},
	}
	t.Cleanup(func() {
		setupDeps = old
	})
}

func assertSetupLineOrder(t *testing.T, output string, lines []string) {
	t.Helper()
	last := -1
	for _, line := range lines {
		index := strings.Index(output, line)
		if index < 0 {
			t.Fatalf("expected setup output to contain %q, got %q", line, output)
		}
		if index < last {
			t.Fatalf("expected setup output line %q after previous lines, got %q", line, output)
		}
		last = index
	}
}

type recordingOutcomeNotifier struct {
	mu       sync.Mutex
	err      error
	outcomes []roundnotify.Outcome
}

func (notifier *recordingOutcomeNotifier) Notify(ctx context.Context, outcome roundnotify.Outcome) error {
	if ctx == nil {
		return errors.New("notification context is required")
	}
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	notifier.outcomes = append(notifier.outcomes, outcome)
	return notifier.err
}

func (notifier *recordingOutcomeNotifier) recorded() []roundnotify.Outcome {
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	return append([]roundnotify.Outcome(nil), notifier.outcomes...)
}

type deadlineRecordingOutcomeNotifier struct {
	hadDeadline bool
	deadline    time.Time
	wasCanceled bool
	outcome     roundnotify.Outcome
}

func (notifier *deadlineRecordingOutcomeNotifier) Notify(ctx context.Context, outcome roundnotify.Outcome) error {
	deadline, ok := ctx.Deadline()
	notifier.hadDeadline = ok
	notifier.deadline = deadline
	select {
	case <-ctx.Done():
		notifier.wasCanceled = true
	default:
	}
	notifier.outcome = outcome
	return nil
}

func withOutcomeNotifier(t *testing.T, notifier roundnotify.Notifier) {
	t.Helper()
	withOutcomeNotifierFactory(t, func(roundconfig.Config) roundnotify.Notifier {
		return notifier
	})
}

func withOutcomeNotifierFactory(t *testing.T, factory func(roundconfig.Config) roundnotify.Notifier) {
	t.Helper()
	old := newOutcomeNotifier
	newOutcomeNotifier = factory
	t.Cleanup(func() {
		newOutcomeNotifier = old
	})
}

func assertRecordedOutcomes(t *testing.T, notifier *recordingOutcomeNotifier, want []roundnotify.Outcome) {
	t.Helper()
	got := notifier.recorded()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recorded notification outcomes mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func runCleanResolveForOutcomeNotification(t *testing.T, notifyErr error) (string, string, int, string) {
	t.Helper()
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	notifier := &recordingOutcomeNotifier{err: notifyErr}
	withOutcomeNotifier(t, notifier)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(notifier.recorded()) != 1 {
		t.Fatalf("expected one notification attempt, got %#v", notifier.recorded())
	}
	return stdout.String(), stderr.String(), code, homeDir
}

func withCLIWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	mustMkdir(t, filepath.Join(repoDir, ".git"))
	t.Setenv("HOME", homeDir)
	t.Chdir(repoDir)
	return homeDir, repoDir
}

type runListSeed struct {
	kind        string
	state       string
	gitRoot     string
	branch      string
	prNumber    string
	specSlug    string
	createdAt   time.Time
	completedAt time.Time
}

func seedRunsForList(t *testing.T, homeDir string, seeds []runListSeed) []store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open run store: %v", err)
	}
	runs := make([]store.Run, 0, len(seeds))
	for _, seed := range seeds {
		run, err := runStore.CreateRun(ctx, createRunRequestForList(seed))
		if err != nil {
			t.Fatalf("create listed Run: %v", err)
		}
		switch {
		case seed.state == "" || seed.state == store.StateActive:
		case store.IsTerminalState(seed.state):
			completed, completeErr := runStore.CompleteRun(ctx, run.ID, seed.state)
			err = completeErr
			if err != nil {
				t.Fatalf("complete listed Run as %s: %v", seed.state, err)
			}
			run = completed.Run
		default:
			if err := runStore.UpdateRunState(ctx, run.ID, seed.state); err != nil {
				t.Fatalf("update listed Run state to %s: %v", seed.state, err)
			}
			var found bool
			run, found, err = runStore.Run(ctx, run.ID)
			if err != nil {
				t.Fatalf("read updated listed Run: %v", err)
			}
			if !found {
				t.Fatalf("expected listed Run %s to exist", run.ID)
			}
		}
		runs = append(runs, run)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close run store: %v", err)
	}
	for index, seed := range seeds {
		setListedRunCreatedAt(t, homeDir, runs[index].ID, seed.createdAt)
		if !seed.completedAt.IsZero() {
			setListedRunCompletedAt(t, homeDir, runs[index].ID, seed.completedAt)
		}
	}
	return runs
}

func createRunRequestForList(seed runListSeed) store.CreateRunRequest {
	req := store.CreateRunRequest{
		Kind:            seed.kind,
		GitRoot:         seed.gitRoot,
		LocalBranch:     seed.branch,
		HeadSHA:         "abc123",
		ArtifactDir:     filepath.Join(seed.gitRoot, ".roundfix"),
		Agent:           "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
	}
	if req.Kind == store.KindImplement {
		req.SpecSlug = seed.specSlug
		return req
	}
	req.HeadRepository = "owner/project"
	req.HeadBranch = seed.branch
	req.BaseRepository = "owner/project"
	req.PRNumber = seed.prNumber
	return req
}

func setListedRunCreatedAt(t *testing.T, homeDir string, runID string, createdAt time.Time) {
	t.Helper()
	ctx := t.Context()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open raw Run Database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close raw Run Database: %v", err)
		}
	}()
	result, err := db.ExecContext(
		ctx,
		`UPDATE runs SET created_at = ?, updated_at = ? WHERE id = ?`,
		createdAt.UTC().Format(time.RFC3339Nano),
		createdAt.UTC().Format(time.RFC3339Nano),
		runID,
	)
	if err != nil {
		t.Fatalf("set listed Run creation time: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("read listed Run timestamp update result: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected to update one listed Run timestamp, updated %d", affected)
	}
}

func setListedRunCompletedAt(t *testing.T, homeDir string, runID string, completedAt time.Time) {
	t.Helper()
	ctx := t.Context()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open raw Run Database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close raw Run Database: %v", err)
		}
	}()
	result, err := db.ExecContext(
		ctx,
		`UPDATE runs SET completed_at = ? WHERE id = ?`,
		completedAt.UTC().Format(time.RFC3339Nano),
		runID,
	)
	if err != nil {
		t.Fatalf("set listed Run completion time: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("read listed Run completion update result: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected to update one listed Run completion time, updated %d", affected)
	}
}

func withRunsListNow(t *testing.T, now time.Time) {
	t.Helper()
	old := runsListNow
	runsListNow = func() time.Time { return now }
	t.Cleanup(func() {
		runsListNow = old
	})
}

func withRunsInteractiveInput(t *testing.T, available bool) {
	t.Helper()
	old := runsInteractiveInputAvailable
	runsInteractiveInputAvailable = func() bool {
		return available
	}
	t.Cleanup(func() {
		runsInteractiveInputAvailable = old
	})
}

func withReviewGitWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := t.TempDir()
	gitImplement(t, repoDir, "init", "--initial-branch=main")
	gitImplement(t, repoDir, "config", "user.name", "Roundfix Test")
	gitImplement(t, repoDir, "config", "user.email", "roundfix-test@example.com")
	gitImplement(t, repoDir, "config", "commit.gpgsign", "false")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
defaults:
  verification: test -f agent.txt
watch:
  poll_interval: 1s
  review_timeout: 5s
  quiet_period: 1ms
`)
	mustWrite(t, filepath.Join(repoDir, "README.md"), "seed\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed review repo")
	gitImplement(t, repoDir, "checkout", "-b", "feature/review")
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve review repo dir: %v", err)
	}
	t.Chdir(resolved)
	return homeDir, resolved
}

func withRealReviewPreflight(t *testing.T, repoDir string, pushEnabled bool) {
	t.Helper()
	withPreflight(t, func(ctx context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		if req.pr == "" {
			return preflight.Result{}, errors.New("missing pr in test preflight")
		}
		gitState, err := preflight.InspectGit(ctx, repoDir, preflight.ExecGitRunner{})
		if err != nil {
			return preflight.Result{}, err
		}
		gitState.UpstreamRemote = "origin"
		gitState.UpstreamBranch = "feature/review"
		return preflight.Result{
			Git: gitState,
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{
				Enabled: pushEnabled,
				Remote:  "origin",
				Branch:  "feature/review",
				Command: []string{"git", "push", "origin", "HEAD:feature/review"},
			},
		}, nil
	})
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func withSuccessfulPreflight(t *testing.T, repoDir string) {
	t.Helper()
	withAgentRunner(t, &fakeAgentRunner{})
	withNoReviewRunWorktrees(t)
	withFakeWorktree(t)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withBranchIntegrity(t, nil, nil)
	clock := &fakeWatchClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	withWatchTiming(t, clock, &fakeWatchSleeper{clock: clock})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{{State: watch.StatusSettled}},
	}).Status)
	withWatchHeadCheck(t, (&fakeWatchHeadCheck{
		states: []watch.HeadCheckState{watch.CheckSuccess},
	}).Check)
	withWatchHeadSHA(t, func(context.Context, string) (string, error) {
		return "abc123", nil
	})
	withChangedPaths(t, nil)
	withReviewSpecGitRunner(t, &fakeReviewSpecGitRunner{})
	withFetchReviewItems(t, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal: `$(rm -rf /)`.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		if req.pr == "" {
			return preflight.Result{}, errors.New("missing pr in test preflight")
		}
		return preflight.Result{
			Git: preflight.GitState{
				Root:            repoDir,
				Branch:          "feature/review",
				HEAD:            "abc123",
				UpstreamRemote:  "origin",
				UpstreamBranch:  "feature/review",
				UnpushedCommits: 1,
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{
				Enabled: req.name != "fetch",
				Remote:  "origin",
				Branch:  "feature/review",
				Command: []string{"git", "push", "origin", "HEAD:feature/review"},
			},
		}, nil
	})
}

func withNoReviewRunWorktrees(t *testing.T) {
	t.Helper()
	oldCreate := createRunWorktree
	oldIntegrate := integrateRunWorktree
	oldCleanup := cleanupCleanRunWorktree
	oldPrune := pruneTerminalRunWorktrees
	createRunWorktree = func(_ context.Context, opts runworktree.CreateOptions) (runworktree.Ref, error) {
		return runworktree.Ref{}, fmt.Errorf("review command unexpectedly created a Run Worktree for run %s", opts.RunID)
	}
	integrateRunWorktree = func(context.Context, runworktree.Ref, string, string) (runworktree.IntegrationResult, error) {
		return runworktree.IntegrationResult{}, errors.New("review command unexpectedly integrated a Run Worktree")
	}
	cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
		return fmt.Errorf("review command unexpectedly cleaned a Run Worktree %s", ref.Path)
	}
	pruneTerminalRunWorktrees = func(context.Context, string, string, runworktree.TerminalRunReconciliationStore, runworktree.TerminalRunLookup) ([]runworktree.PrunedRef, error) {
		return nil, nil
	}
	t.Cleanup(func() {
		createRunWorktree = oldCreate
		integrateRunWorktree = oldIntegrate
		cleanupCleanRunWorktree = oldCleanup
		pruneTerminalRunWorktrees = oldPrune
	})
}

func withFetchReviewItems(t *testing.T, items []reviewsource.ReviewItem) {
	t.Helper()
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		return items, nil
	})
}

func withFetchReviewItemsFunc(t *testing.T, fn func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error)) {
	t.Helper()
	old := fetchReviewItems
	fetchReviewItems = fn
	t.Cleanup(func() {
		fetchReviewItems = old
	})
}

func withPreflight(t *testing.T, fn func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error)) {
	t.Helper()
	old := runCommandPreflight
	runCommandPreflight = fn
	t.Cleanup(func() {
		runCommandPreflight = old
	})
}

type branchIntegrityRecorder struct {
	integrations []string
}

func withBranchIntegrity(t *testing.T, pending []runworktree.PendingRunWork, integrateErr error) *branchIntegrityRecorder {
	t.Helper()
	recorder := &branchIntegrityRecorder{}
	oldList := listPendingRunWork
	oldIntegrate := integratePendingRunWork
	oldRefresh := refreshBranchIntegrityHead
	listPendingRunWork = func(context.Context, string, string) ([]runworktree.PendingRunWork, error) {
		return append([]runworktree.PendingRunWork(nil), pending...), nil
	}
	integratePendingRunWork = func(_ context.Context, _ string, _ string, runBranch string) error {
		recorder.integrations = append(recorder.integrations, runBranch)
		return integrateErr
	}
	refreshBranchIntegrityHead = func(_ context.Context, result preflight.Result) (preflight.Result, error) {
		result.Git.HEAD = "integrated-head"
		return result, nil
	}
	t.Cleanup(func() {
		listPendingRunWork = oldList
		integratePendingRunWork = oldIntegrate
		refreshBranchIntegrityHead = oldRefresh
	})
	return recorder
}

type pullRequestCommentRecorder struct {
	calls []pullRequestCommentRecord
}

type pullRequestCommentRecord struct {
	source     string
	repository string
	pr         int
	body       string
}

func withPullRequestComments(t *testing.T, publishErr error) *pullRequestCommentRecorder {
	t.Helper()
	recorder := &pullRequestCommentRecorder{}
	old := commentOnPullRequest
	commentOnPullRequest = func(_ context.Context, source string, repository string, prNumber int, body string) error {
		recorder.calls = append(recorder.calls, pullRequestCommentRecord{source: source, repository: repository, pr: prNumber, body: body})
		return publishErr
	}
	t.Cleanup(func() {
		commentOnPullRequest = old
	})
	return recorder
}

func branchIntegrityCommandArgs(command string) []string {
	switch command {
	case "fetch":
		return []string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}
	case "resolve":
		return []string{"resolve", "--pr", "123", "--no-input"}
	case "watch":
		return []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--no-input"}
	default:
		return []string{command}
	}
}

func journalPayloadContains(events []store.JournalEvent, text string) bool {
	for _, event := range events {
		if strings.Contains(string(event.Event.Payload), text) {
			return true
		}
	}
	return false
}

func withReviewSpecGitRunner(t *testing.T, runner preflight.GitRunner) {
	t.Helper()
	old := reviewSpecGitRunner
	reviewSpecGitRunner = runner
	t.Cleanup(func() {
		reviewSpecGitRunner = old
	})
}

func overrideCollaborators(t *testing.T, mutate func(*engineCollaborators)) {
	t.Helper()
	old := newEngineCollaborators
	current := old()
	mutate(&current)
	newEngineCollaborators = func() engineCollaborators { return current }
	t.Cleanup(func() {
		newEngineCollaborators = old
	})
}

func withAgentRunner(t *testing.T, runner agent.Runner) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.runner = runner
	})
	withRoundfixSessionLister(t, func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error) {
		return nil, nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
}

func withFallbackConfirmation(t *testing.T, input string) {
	t.Helper()
	oldAvailable := fallbackConfirmationAvailable
	oldInput := fallbackConfirmationInput
	fallbackConfirmationAvailable = func(io.Writer) bool { return true }
	fallbackConfirmationInput = func() io.Reader { return strings.NewReader(input) }
	t.Cleanup(func() {
		fallbackConfirmationAvailable = oldAvailable
		fallbackConfirmationInput = oldInput
	})
}

func withProfilesConfigureInput(t *testing.T, input string) {
	t.Helper()
	old := profilesConfigureInput
	profilesConfigureInput = func() io.Reader { return strings.NewReader(input) }
	t.Cleanup(func() {
		profilesConfigureInput = old
	})
}

func withProfilesConfigureConfirm(t *testing.T, confirm func(context.Context, io.Writer, string) (bool, error)) {
	t.Helper()
	old := confirmProfilesConfigure
	confirmProfilesConfigure = confirm
	t.Cleanup(func() {
		confirmProfilesConfigure = old
	})
}

func withSuccessfulProfilesConfigureProof(t *testing.T) *profileReadinessExactRunner {
	t.Helper()
	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			return agent.SelectionProof{
				Runtime:         req.Runtime.ID,
				Model:           req.Runtime.Model,
				ReasoningEffort: req.Runtime.ReasoningEffort,
				Status:          agent.SelectionProofStatusProven,
			}, nil
		},
	}
	withAgentRunner(t, runner)
	return runner
}

func profilesConfigureTestFragment(preferredModel string, fallbackModel string) string {
	return fmt.Sprintf(`
backend:
  preferred:
    runtime: codex
    model: %s
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: %s
      reasoning_effort: xhigh
`, preferredModel, fallbackModel)
}

func withStopAgentSessionCanceler(t *testing.T, cancel func(context.Context, agent.RuntimeSpec, agent.SessionRef) error) {
	t.Helper()
	old := cancelStopAgentSession
	cancelStopAgentSession = cancel
	t.Cleanup(func() {
		cancelStopAgentSession = old
	})
	withRoundfixSessionLister(t, func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error) {
		return nil, nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
}

func withRoundfixSessionLister(t *testing.T, list func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error)) {
	t.Helper()
	old := listRoundfixAgentSessions
	listRoundfixAgentSessions = list
	t.Cleanup(func() {
		listRoundfixAgentSessions = old
	})
}

func withStopAgentSessionCloser(t *testing.T, closeSession func(context.Context, agent.RuntimeSpec, agent.SessionRef) error) {
	t.Helper()
	old := closeStopAgentSession
	closeStopAgentSession = closeSession
	t.Cleanup(func() {
		closeStopAgentSession = old
	})
}

// ownerProcessControllerFunc drives the termination phase only; its owner
// proof always succeeds. Tests that need a failing or observable proof use
// ownerProcessControllerStub.
type ownerProcessControllerFunc func(context.Context, int, string) error

func (controller ownerProcessControllerFunc) ProveOwner(context.Context, int, string) error {
	return nil
}

func (controller ownerProcessControllerFunc) TerminateAndWait(ctx context.Context, pid int, recordedIdentity string) error {
	return controller(ctx, pid, recordedIdentity)
}

// ownerProcessControllerStub drives the owner proof and the termination phase
// independently so tests can assert their order and their failure isolation.
type ownerProcessControllerStub struct {
	prove     func(context.Context, int, string) error
	terminate func(context.Context, int, string) error
}

func (controller ownerProcessControllerStub) ProveOwner(ctx context.Context, pid int, recordedIdentity string) error {
	if controller.prove == nil {
		return nil
	}
	return controller.prove(ctx, pid, recordedIdentity)
}

func (controller ownerProcessControllerStub) TerminateAndWait(ctx context.Context, pid int, recordedIdentity string) error {
	if controller.terminate == nil {
		return nil
	}
	return controller.terminate(ctx, pid, recordedIdentity)
}

func withOwnerProcessController(t *testing.T, controller OwnerProcessController) {
	t.Helper()
	old := ownerProcesses
	ownerProcesses = controller
	t.Cleanup(func() {
		ownerProcesses = old
	})
}

func fakeACPXVersionCommand(t *testing.T, version string) string {
	t.Helper()
	script := filepath.Join(t.TempDir(), "acpx")
	body := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then\n" +
		"  printf '%s\\n' '" + version + "'\n" +
		"  exit 0\n" +
		"fi\n" +
		"printf 'unexpected acpx invocation: %s\\n' \"$*\" >&2\n" +
		"exit 2\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("write fake acpx command: %v", err)
	}
	return script
}

func withVerifier(t *testing.T, verifier daemon.Verifier) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.verifier = verifier
	})
}

func withCommitter(t *testing.T, committer daemon.Committer) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.committer = committer
	})
}

func withSourceResolver(t *testing.T, resolver *fakeSourceResolver) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.source = resolver
	})
}

func withPusher(t *testing.T, pusher daemon.Pusher) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.pusher = pusher
	})
}

// fakeWorktree reports an Agent-made change on every post-Batch snapshot so
// CLI tests exercise the commit path without a real git repository.
type fakeWorktree struct {
	mu    sync.Mutex
	calls int
}

func (worktree *fakeWorktree) Snapshot(context.Context, string) ([]string, error) {
	worktree.mu.Lock()
	defer worktree.mu.Unlock()
	worktree.calls++
	if worktree.calls%2 == 0 {
		return []string{"src/agent-change.go"}, nil
	}
	return nil, nil
}

func withFakeWorktree(t *testing.T) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.worktree = &fakeWorktree{}
	})
}

func withWatchStatus(t *testing.T, fn func(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error)) {
	t.Helper()
	old := watchReviewStatus
	watchReviewStatus = fn
	t.Cleanup(func() {
		watchReviewStatus = old
	})
}

func withWatchHeadCheck(t *testing.T, fn func(context.Context, reviewsource.HeadCheckRequest) (watch.HeadCheckState, error)) {
	t.Helper()
	old := watchHeadCheck
	watchHeadCheck = fn
	t.Cleanup(func() {
		watchHeadCheck = old
	})
}

func withWatchHeadSHA(t *testing.T, fn func(context.Context, string) (string, error)) {
	t.Helper()
	old := watchHeadSHA
	watchHeadSHA = fn
	t.Cleanup(func() {
		watchHeadSHA = old
	})
}

func withWatchTiming(t *testing.T, clock watch.Clock, sleeper watch.Sleeper) {
	t.Helper()
	oldClock := watchClock
	oldSleeper := watchSleeper
	watchClock = clock
	watchSleeper = sleeper
	t.Cleanup(func() {
		watchClock = oldClock
		watchSleeper = oldSleeper
	})
}

func withInteractiveInput(t *testing.T, fn func(context.Context, roundtui.InputRequest) (roundtui.CommandValues, error)) {
	t.Helper()
	old := collectInteractiveInput
	collectInteractiveInput = fn
	t.Cleanup(func() {
		collectInteractiveInput = old
	})
}

func withInitScopePrompt(t *testing.T, fn func(context.Context, io.Writer) (string, error)) {
	t.Helper()
	old := promptInitScope
	promptInitScope = fn
	t.Cleanup(func() {
		promptInitScope = old
	})
}

func withClaudeSkillSymlinkPrompt(t *testing.T, fn func(context.Context, io.Writer, string, string) (bool, error)) {
	t.Helper()
	old := promptProjectClaudeSkillSymlink
	promptProjectClaudeSkillSymlink = fn
	t.Cleanup(func() {
		promptProjectClaudeSkillSymlink = old
	})
}

func withCurrentPullRequestSuggestion(t *testing.T, value string) {
	t.Helper()
	old := suggestCurrentPullRequest
	suggestCurrentPullRequest = func(context.Context, string) (string, error) {
		return value, nil
	}
	t.Cleanup(func() {
		suggestCurrentPullRequest = old
	})
}

func withStopPullRequestResolver(t *testing.T, result preflight.PullRequest) {
	t.Helper()
	old := resolvePullRequestForStop
	resolvePullRequestForStop = func(context.Context, string, string) (preflight.PullRequest, error) {
		return result, nil
	}
	t.Cleanup(func() {
		resolvePullRequestForStop = old
	})
}

func withChangedPaths(t *testing.T, changes []preflight.ChangedPath) {
	t.Helper()
	old := inspectChangedPaths
	inspectChangedPaths = func(context.Context, string) ([]preflight.ChangedPath, error) {
		return changes, nil
	}
	t.Cleanup(func() {
		inspectChangedPaths = old
	})
}

func withChangedPathSnapshots(t *testing.T, snapshots ...[]preflight.ChangedPath) {
	t.Helper()
	old := inspectChangedPaths
	index := 0
	inspectChangedPaths = func(context.Context, string) ([]preflight.ChangedPath, error) {
		if len(snapshots) == 0 {
			return nil, nil
		}
		if index >= len(snapshots) {
			return snapshots[len(snapshots)-1], nil
		}
		changes := snapshots[index]
		index++
		return changes, nil
	}
	t.Cleanup(func() {
		inspectChangedPaths = old
	})
}

func readInteractiveDefaults(t *testing.T, homeDir string) store.InteractiveDefaults {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for defaults: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for defaults: %v", err)
		}
	}()
	defaults, err := runStore.InteractiveDefaults(context.Background())
	if err != nil {
		t.Fatalf("read defaults: %v", err)
	}
	return defaults
}

func assertRunCount(t *testing.T, dbPath string, expected int) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, filepath.Dir(filepath.Dir(dbPath)))
	if err != nil {
		t.Fatalf("open store for count: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for count: %v", err)
		}
	}()
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != expected {
		t.Fatalf("expected %d Run record(s), got %d", expected, count)
	}
}

func assertNoRunDatabase(t *testing.T, homeDir string) {
	t.Helper()
	if _, err := os.Stat(store.DatabasePath(homeDir)); err == nil {
		t.Fatalf("expected no Run Database at %s", store.DatabasePath(homeDir))
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Database: %v", err)
	}
}

func createActiveImplementRunForStop(t *testing.T, homeDir string, repoDir string, specSlug string, agentID string) (store.Run, store.CreateRunRequest) {
	t.Helper()
	ownerPID := os.Getpid()
	ownerIdentity := "cli-test-owner-identity"
	request := store.CreateRunRequest{
		Kind:          store.KindImplement,
		GitRoot:       repoDir,
		LocalBranch:   "ma/implement-spec",
		HeadSHA:       "abc123",
		SpecSlug:      specSlug,
		Agent:         agentID,
		OwnerPID:      ownerPID,
		OwnerIdentity: ownerIdentity,
	}
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	active, err := runStore.CreateRun(context.Background(), request)
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	withOwnerProcessController(t, ownerProcessControllerFunc(func(_ context.Context, pid int, recordedIdentity string) error {
		if pid != ownerPID {
			t.Fatalf("owner process PID = %d, want %d", pid, ownerPID)
		}
		if recordedIdentity != ownerIdentity {
			t.Fatalf("owner identity = %q, want recorded identity %q", recordedIdentity, ownerIdentity)
		}
		return nil
	}))
	return active, request
}

func recordAgentSelectionForStop(t *testing.T, homeDir string, req store.AgentSelectionAttemptRequest) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store to record Agent Selection: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store after recording Agent Selection: %v", err)
		}
	}()
	if _, err := runStore.AppendAgentSelectionAttempt(context.Background(), req); err != nil {
		t.Fatalf("record Agent Selection for %s %s: %v", req.ScopeKind, req.ScopeID, err)
	}
}

func assertAgentSelectionScopeStatus(
	t *testing.T,
	homeDir string,
	runID string,
	scopeKind store.AgentSelectionScopeKind,
	scopeID string,
	want store.AgentSelectionStatus,
) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store to inspect Agent Selection: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store after inspecting Agent Selection: %v", err)
		}
	}()
	attempts, err := runStore.AgentSelectionAttemptsForScope(context.Background(), runID, scopeKind, scopeID)
	if err != nil {
		t.Fatalf("read Agent Selection for %s %s: %v", scopeKind, scopeID, err)
	}
	if len(attempts) == 0 {
		t.Fatalf("expected Agent Selection lifecycle for %s %s", scopeKind, scopeID)
	}
	if got := attempts[len(attempts)-1].Status; got != want {
		t.Fatalf("Agent Selection lifecycle for %s %s: got %q, want %q", scopeKind, scopeID, got, want)
	}
}

func createActiveImplementRunForStopInGitRoot(t *testing.T, homeDir string, gitRoot string, specSlug string) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     gitRoot,
		LocalBranch: "ma/other-work",
		HeadSHA:     "abc123",
		SpecSlug:    specSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	return active
}

func assertNoActiveRun(t *testing.T, homeDir string, headRepository string, headBranch string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for active run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for active run: %v", err)
		}
	}()
	if active, found, err := runStore.ActiveRun(context.Background(), headRepository, headBranch); err != nil {
		t.Fatalf("lookup active run: %v", err)
	} else if found {
		t.Fatalf("expected no Active Run, got %#v", active)
	}
}

func assertActiveRun(t *testing.T, homeDir string, headRepository string, headBranch string, runID string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for active run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for active run: %v", err)
		}
	}()
	active, found, err := runStore.ActiveRun(context.Background(), headRepository, headBranch)
	if err != nil {
		t.Fatalf("lookup active run: %v", err)
	}
	if !found || active.ID != runID {
		t.Fatalf("expected Active Run %s, found=%v active=%#v", runID, found, active)
	}
}

func assertRunState(t *testing.T, homeDir string, runID string, want string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for Run state: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for Run state: %v", err)
		}
	}()
	run, found, err := runStore.Run(context.Background(), runID)
	if err != nil {
		t.Fatalf("lookup Run state: %v", err)
	}
	if !found {
		t.Fatalf("expected Run %s to exist", runID)
	}
	if run.State != want {
		t.Fatalf("expected Run %s state %s, got %s", runID, want, run.State)
	}
}

func reviewRunIDFromStderr(t *testing.T, stderr string) string {
	t.Helper()
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		if value, ok := strings.CutPrefix(trimmed, "Run: "); ok {
			return strings.TrimSpace(value)
		}
		if value, ok := strings.CutPrefix(trimmed, "Watch Run: "); ok {
			return strings.TrimSpace(value)
		}
	}
	t.Fatalf("expected review Run id in stderr, got %q", stderr)
	return ""
}

func assertStopRequested(t *testing.T, homeDir string, runID string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for Stop Request flag: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for Stop Request flag: %v", err)
		}
	}()
	requested, err := runStore.StopRequested(context.Background(), runID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatalf("expected Stop Request recorded for Run %s", runID)
	}
}

func assertAgentLogContains(t *testing.T, repoDir string, expected string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", "*", "agent", "batch-001.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one Agent log, got %#v", matches)
	}
	content, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read Agent log %s: %v", matches[0], err)
	}
	if !strings.Contains(string(content), expected) {
		t.Fatalf("expected Agent log to contain %q, got %s", expected, string(content))
	}
}

func assertNoAgentLogs(t *testing.T, repoDir string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", "*", "agent", "batch-*.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no Agent logs, got %#v", matches)
	}
}

func persistCLIReviewIssue(t *testing.T, repoDir string, roundNumber int, headBranch string) rounds.PersistResult {
	t.Helper()
	return persistCLIReviewItems(t, repoDir, roundNumber, headBranch, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		},
	})
}

func persistCLIReviewItems(t *testing.T, repoDir string, roundNumber int, headBranch string, items []reviewsource.ReviewItem) rounds.PersistResult {
	t.Helper()
	result, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    builtinArtifactDirForRepo(t, repoDir),
		ReviewRoot:     defaultReviewRootForRepo(repoDir, "123"),
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     headBranch,
		HeadSHA:        "abc123",
		Round:          roundNumber,
		CreatedAt:      time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		Items:          items,
	})
	if err != nil {
		t.Fatalf("persist CLI Review Issue artifact: %v", err)
	}
	return result
}

func defaultReviewRootForRepo(repoDir string, prNumber string) string {
	return filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-"+prNumber)
}

func builtinArtifactDirForRepo(t *testing.T, repoDir string) string {
	t.Helper()
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Fatal("HOME is required for builtin Artifact Directory")
	}
	sum := sha256.Sum256([]byte(filepath.Clean(repoDir)))
	repoID := hex.EncodeToString(sum[:])[:16]
	return filepath.Join(homeDir, ".roundfix", "artifacts", repoID)
}

type fakeVerifier struct {
	err      error
	calls    int
	commands []string
}

func (verifier *fakeVerifier) Verify(_ context.Context, req daemon.VerifyRequest) (daemon.VerifyResult, error) {
	verifier.calls++
	verifier.commands = append(verifier.commands, req.Command)
	if req.OutputPath != "" {
		if err := os.MkdirAll(filepath.Dir(req.OutputPath), 0o755); err != nil {
			return daemon.VerifyResult{OutputPath: req.OutputPath}, err
		}
		if err := os.WriteFile(req.OutputPath, []byte("fake verification output\n"), 0o644); err != nil {
			return daemon.VerifyResult{OutputPath: req.OutputPath}, err
		}
	}
	if verifier.err != nil {
		return daemon.VerifyResult{OutputPath: req.OutputPath}, &daemon.VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: verifier.err}
	}
	if req.OutputPath != "" {
		_ = os.Remove(req.OutputPath)
	}
	return daemon.VerifyResult{OutputPath: req.OutputPath}, nil
}

type fakeCommitter struct {
	err         error
	afterCommit func(context.Context, daemon.CommitRequest) error
	calls       int
	workDirs    []string
	messages    []string
}

func (committer *fakeCommitter) Commit(ctx context.Context, req daemon.CommitRequest) error {
	committer.calls++
	committer.workDirs = append(committer.workDirs, req.WorkDir)
	committer.messages = append(committer.messages, req.Message)
	if committer.err != nil {
		return committer.err
	}
	if committer.afterCommit != nil {
		return committer.afterCommit(ctx, req)
	}
	return nil
}

type fakeSourceResolver struct {
	err             error
	calls           int
	requests        []reviewsource.ResolveRequest
	resolveRequests []reviewsource.IssueResolveRequest
	commentRequests []reviewsource.IssueCommentRequest
}

func (resolver *fakeSourceResolver) Resolve(_ context.Context, req reviewsource.ResolveRequest) error {
	return resolver.ResolveIssues(context.Background(), req)
}

func (resolver *fakeSourceResolver) ResolveIssues(_ context.Context, req reviewsource.ResolveRequest) error {
	resolver.calls++
	resolver.requests = append(resolver.requests, req)
	return resolver.err
}

func (resolver *fakeSourceResolver) ResolveIssue(_ context.Context, req reviewsource.IssueResolveRequest) error {
	resolver.calls++
	resolver.resolveRequests = append(resolver.resolveRequests, req)
	return resolver.err
}

func (resolver *fakeSourceResolver) ReplyToIssue(_ context.Context, req reviewsource.IssueCommentRequest) (reviewsource.IssueCommentResult, error) {
	resolver.calls++
	resolver.commentRequests = append(resolver.commentRequests, req)
	if resolver.err != nil {
		return reviewsource.IssueCommentResult{}, resolver.err
	}
	return reviewsource.IssueCommentResult{Posted: true}, nil
}

type fakePusher struct {
	err      error
	calls    int
	workDirs []string
	remotes  []string
	branches []string
	args     []string
}

func (pusher *fakePusher) Push(_ context.Context, req daemon.PushRequest) error {
	pusher.calls++
	pusher.workDirs = append(pusher.workDirs, req.WorkDir)
	pusher.remotes = append(pusher.remotes, req.Remote)
	pusher.branches = append(pusher.branches, req.Branch)
	pusher.args = append(pusher.args, "push", req.Remote, "HEAD:"+req.Branch)
	return pusher.err
}

type checkingPusher struct {
	calls   int
	workDir []string
	check   func(daemon.PushRequest) error
}

func (pusher *checkingPusher) Push(_ context.Context, req daemon.PushRequest) error {
	pusher.calls++
	pusher.workDir = append(pusher.workDir, req.WorkDir)
	if pusher.check != nil {
		return pusher.check(req)
	}
	return nil
}

type fakeWatchStatus struct {
	err      error
	calls    int
	statuses []reviewsource.WatchStatus
}

func (source *fakeWatchStatus) Status(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
	source.calls++
	if source.err != nil {
		return reviewsource.WatchStatus{}, source.err
	}
	if len(source.statuses) == 0 {
		return reviewsource.WatchStatus{State: watch.StatusSettled}, nil
	}
	status := source.statuses[0]
	if len(source.statuses) > 1 {
		source.statuses = source.statuses[1:]
	}
	return status, nil
}

type fakeWatchHeadCheck struct {
	err      error
	calls    int
	states   []watch.HeadCheckState
	headSHAs []string
}

func (source *fakeWatchHeadCheck) Check(_ context.Context, req reviewsource.HeadCheckRequest) (watch.HeadCheckState, error) {
	source.calls++
	source.headSHAs = append(source.headSHAs, req.HeadSHA)
	if source.err != nil {
		return "", source.err
	}
	if len(source.states) == 0 {
		return watch.CheckSuccess, nil
	}
	state := source.states[0]
	if len(source.states) > 1 {
		source.states = source.states[1:]
	}
	return state, nil
}

type fakeWatchClock struct {
	now time.Time
}

func (clock *fakeWatchClock) Now() time.Time {
	return clock.now
}

func (clock *fakeWatchClock) Advance(duration time.Duration) {
	clock.now = clock.now.Add(duration)
}

type fakeWatchSleeper struct {
	clock *fakeWatchClock
}

func (sleeper *fakeWatchSleeper) Sleep(_ context.Context, duration time.Duration) error {
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	return nil
}

type orderedWriteEvent struct {
	stream string
	text   string
}

type orderedWriteRecorder struct {
	events []orderedWriteEvent
}

func (recorder *orderedWriteRecorder) writer(stream string) io.Writer {
	return orderedStreamWriter{recorder: recorder, stream: stream}
}

func (recorder *orderedWriteRecorder) content(stream string) string {
	var builder strings.Builder
	for _, event := range recorder.events {
		if event.stream == stream {
			builder.WriteString(event.text)
		}
	}
	return builder.String()
}

func (recorder *orderedWriteRecorder) firstIndex(stream string) int {
	for index, event := range recorder.events {
		if event.stream == stream {
			return index
		}
	}
	return -1
}

func (recorder *orderedWriteRecorder) firstIndexContaining(stream string, needle string) int {
	for index, event := range recorder.events {
		if event.stream == stream && strings.Contains(event.text, needle) {
			return index
		}
	}
	return -1
}

type orderedStreamWriter struct {
	recorder *orderedWriteRecorder
	stream   string
}

func (writer orderedStreamWriter) Write(data []byte) (int, error) {
	writer.recorder.events = append(writer.recorder.events, orderedWriteEvent{
		stream: writer.stream,
		text:   string(data),
	})
	return len(data), nil
}

type failingWriter struct {
	err error
}

func (writer failingWriter) Write([]byte) (int, error) {
	return 0, writer.err
}

func publishFakeAgentOutput(ctx context.Context, sink runevent.Sink, req agent.ExecuteRequest, text string) error {
	if sink == nil {
		return nil
	}
	payload, err := json.Marshal(struct {
		Text string `json:"text"`
	}{Text: text})
	if err != nil {
		return err
	}
	return sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Summary: text,
		Payload: payload,
	})
}

type fakeAgentRunner struct {
	probeErr       error
	probe          func(agent.ProbeRequest) error
	fallback       agent.FallbackSelection
	fallbackOK     bool
	fallbackErr    error
	runErr         error
	status         string
	terminalReason string
	statuses       []string
	onRun          func(agent.ExecuteRequest) error
	calls          int
	gitRoots       []string
	probeRequests  []agent.ProbeRequest
	fallbackModels []agent.FallbackCandidateSet
	fallbackRuns   []agent.RuntimeSpec
	probedRuntimes []agent.RuntimeSpec
	runRuntimes    []agent.RuntimeSpec
}

type profileReadinessExactRunner struct {
	fakeAgentRunner
	prove         func(agent.ProbeRequest) (agent.SelectionProof, error)
	exactRequests []agent.ProbeRequest
	probeCalls    int
}

func (runner *profileReadinessExactRunner) ProveProfileSelection(_ context.Context, req agent.ProbeRequest) (agent.SelectionProof, error) {
	runner.exactRequests = append(runner.exactRequests, req)
	if runner.prove == nil {
		return agent.SelectionProof{}, nil
	}
	return runner.prove(req)
}

func (runner *profileReadinessExactRunner) Probe(ctx context.Context, req agent.ProbeRequest) error {
	runner.probeCalls++
	return runner.fakeAgentRunner.Probe(ctx, req)
}

func (runner *fakeAgentRunner) Probe(_ context.Context, req agent.ProbeRequest) error {
	runner.probeRequests = append(runner.probeRequests, req)
	runner.probedRuntimes = append(runner.probedRuntimes, req.Runtime)
	if runner.probe != nil {
		return runner.probe(req)
	}
	return runner.probeErr
}

func (runner *fakeAgentRunner) ProbeFallback(_ context.Context, runtime agent.RuntimeSpec, candidates agent.FallbackCandidateSet) (agent.FallbackSelection, bool, error) {
	runner.fallbackRuns = append(runner.fallbackRuns, runtime)
	runner.fallbackModels = append(runner.fallbackModels, candidates)
	return runner.fallback, runner.fallbackOK, runner.fallbackErr
}

func (runner *fakeAgentRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.calls++
	runner.gitRoots = append(runner.gitRoots, req.GitRoot)
	runner.runRuntimes = append(runner.runRuntimes, req.Runtime)
	status := runner.status
	if len(runner.statuses) >= runner.calls {
		status = runner.statuses[runner.calls-1]
	}
	if status == "" {
		status = rounds.StatusResolved
	}
	if runner.onRun != nil {
		if err := runner.onRun(req); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	for _, issue := range req.Batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, status, "", runner.terminalReason); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	output := "fake agent output\n"
	if err := publishFakeAgentOutput(ctx, sink, req, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if runner.runErr != nil {
		return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, runner.runErr
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, nil
}

func (runner *fakeAgentRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type fakeStoppingAgentRunner struct{}

func (runner *fakeStoppingAgentRunner) Probe(context.Context, agent.ProbeRequest) error {
	return nil
}

func (runner *fakeStoppingAgentRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	output := "partial agent output\n"
	if err := publishFakeAgentOutput(ctx, sink, req, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if sink != nil {
		// Mirror the real runner: a stopped Agent publishes its status
		// event before returning.
		if err := sink.Publish(ctx, runevent.RunEvent{
			RunID:   req.RunID,
			Batch:   req.Batch.Number,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentStatus,
			Summary: "SESSION STOPPED\n",
			Payload: []byte(`{"status":"stopped"}`),
		}); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, agent.StopError{
		LogPath: req.LogPath,
		Output:  output,
		Err:     context.Canceled,
	}
}

func (runner *fakeStoppingAgentRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type sessionRecordingRunner struct {
	inner          agent.Runner
	closeErr       error
	runSessions    []agent.SessionRef
	ensureSessions []agent.SessionRef
	closeSessions  []agent.SessionRef
	seenSessions   map[string]bool
}

func (runner *sessionRecordingRunner) Probe(ctx context.Context, req agent.ProbeRequest) error {
	return runner.inner.Probe(ctx, req)
}

func (runner *sessionRecordingRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.runSessions = append(runner.runSessions, req.Session)
	sessionName := strings.TrimSpace(req.Session.Name)
	if sessionName != "" {
		if runner.seenSessions == nil {
			runner.seenSessions = map[string]bool{}
		}
		if !runner.seenSessions[sessionName] {
			runner.seenSessions[sessionName] = true
			runner.ensureSessions = append(runner.ensureSessions, req.Session)
			if err := publishFakeSessionStatus(ctx, sink, req, "session_started"); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
	}
	return runner.inner.Run(ctx, req, sink)
}

func (runner *sessionRecordingRunner) EndSession(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
	runner.closeSessions = append(runner.closeSessions, session)
	if runner.seenSessions != nil {
		delete(runner.seenSessions, strings.TrimSpace(session.Name))
	}
	if err := runner.inner.EndSession(ctx, runtime, session); err != nil {
		return err
	}
	return runner.closeErr
}

func publishFakeSessionStatus(ctx context.Context, sink runevent.Sink, req agent.ExecuteRequest, status string) error {
	if sink == nil {
		return nil
	}
	payload, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: status})
	if err != nil {
		return err
	}
	return sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: "SESSION " + strings.ToUpper(status) + "\n",
		Payload: payload,
	})
}

func assertRecordedOneSessionForRun(t *testing.T, runner *sessionRecordingRunner, runID string, wantRunCalls int) {
	t.Helper()
	if len(runner.ensureSessions) != wantRunCalls {
		t.Fatalf("expected %d owned session ensure record(s), got %#v", wantRunCalls, runner.ensureSessions)
	}
	for _, session := range runner.ensureSessions {
		if strings.TrimSpace(session.WorkDir) == "" {
			t.Fatalf("expected owned session %q to include a working directory, got %#v", session.Name, runner.ensureSessions)
		}
		if session.Name == "roundfix-"+runID || !strings.HasPrefix(session.Name, "roundfix-"+runID+"-") {
			t.Fatalf("expected owned per-work session for run %s, got %#v", runID, runner.ensureSessions)
		}
	}
	if len(runner.closeSessions) != len(runner.ensureSessions) {
		t.Fatalf("expected one close per owned session %#v, got %#v", runner.ensureSessions, runner.closeSessions)
	}
	for index, ensured := range runner.ensureSessions {
		if runner.closeSessions[index] != ensured {
			t.Fatalf("expected close %d to match ensure %#v, got %#v", index, ensured, runner.closeSessions)
		}
	}
	if len(runner.runSessions) != wantRunCalls {
		t.Fatalf("expected %d Agent run session records, got %#v", wantRunCalls, runner.runSessions)
	}
	for index, session := range runner.runSessions {
		if session != runner.ensureSessions[index] {
			t.Fatalf("expected Agent run %d to use owned session %#v, got %#v", index, runner.ensureSessions[index], runner.runSessions)
		}
	}
}

func assertSessionLifecycleEvents(t *testing.T, events []store.JournalEvent) {
	t.Helper()
	started := false
	closed := false
	for _, entry := range events {
		if entry.Event.Kind != runevent.KindAgentStatus {
			continue
		}
		payload := string(entry.Event.Payload)
		started = started || strings.Contains(payload, `"session_started"`)
		closed = closed || strings.Contains(payload, `"session_closed"`)
	}
	if !started || !closed {
		t.Fatalf("expected session_started and session_closed agent.status events, got %+v", events)
	}
}

func runFromStore(t *testing.T, homeDir string, runID string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	run, found, err := runStore.Run(ctx, runID)
	if err != nil {
		t.Fatalf("read run %s: %v", runID, err)
	}
	if !found {
		t.Fatalf("run %s not found", runID)
	}
	return run
}

func journaledRunEvents(t *testing.T, homeDir string, stderr string) (string, []store.JournalEvent) {
	t.Helper()
	runID := ""
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		if value, ok := strings.CutPrefix(trimmed, "Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
		if value, ok := strings.CutPrefix(trimmed, "Watch Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
		if value, ok := strings.CutPrefix(trimmed, "ID: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
		if value, ok := strings.CutPrefix(trimmed, "Implement Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Run id in stderr, got %q", stderr)
	}
	return runID, runEvents(t, homeDir, runID)
}

func runEvents(t *testing.T, homeDir string, runID string) []store.JournalEvent {
	t.Helper()
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	events, err := reader.RunEventsAfter(context.Background(), runID, 0, 100)
	if err != nil {
		t.Fatalf("list journaled Run Events: %v", err)
	}
	return events
}

func TestAgentConsoleDisplaySinkKeepsWriterBytesByDefault(t *testing.T) {
	event := runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Payload: []byte(`{"text":"fake agent output\n"}`),
	}
	var direct bytes.Buffer
	var wrapped bytes.Buffer

	if err := (agent.WriterSink{Writer: &direct}).Publish(context.Background(), event); err != nil {
		t.Fatalf("publish through direct writer sink: %v", err)
	}
	if err := agentConsoleDisplaySink(&wrapped, false).Publish(context.Background(), event); err != nil {
		t.Fatalf("publish through default display sink: %v", err)
	}

	if wrapped.String() != direct.String() {
		t.Fatalf("expected default display bytes %q, got %q", direct.String(), wrapped.String())
	}
}

func TestAgentConsoleDisplaySinkKeepsDistinctToolCallsVisible(t *testing.T) {
	var buffer bytes.Buffer
	sink := agentConsoleDisplaySink(&buffer, false)
	events := []runevent.RunEvent{
		cliToolLifecycleEvent("run-1", runevent.KindAgentToolUpdated, "edit_call_1", "running", compactCLIEditLifecyclePayload("tool_call_update", "edit_call_1", "internal/app/server.go", "old\n", "new\n", "running")),
		cliToolLifecycleEvent("run-1", runevent.KindAgentToolUpdated, "edit_call_2", "running", compactCLIEditLifecyclePayload("tool_call_update", "edit_call_2", "internal/app/server.go", "old\n", "new\n", "running")),
	}

	for _, event := range events {
		if err := sink.Publish(context.Background(), event); err != nil {
			t.Fatalf("publish distinct tool call: %v", err)
		}
	}

	want := "edit internal/app/server.go (+1/-1)\nedit internal/app/server.go (+1/-1)\n"
	if buffer.String() != want {
		t.Fatalf("expected distinct tool calls to remain visible\nwant: %q\n got: %q", want, buffer.String())
	}
}

func TestAgentConsoleDisplaySinkUsesStatefulSinkForNonTTYAndDetachedLogWriter(t *testing.T) {
	ctx := context.Background()
	homeDir, repoDir := withCLIWorkspace(t)
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		_ = runStore.Close()
	})
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		WorkDir:        repoDir,
		Agent:          "codex",
	})
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	consoleLog, err := os.CreateTemp(t.TempDir(), "console-*.log")
	if err != nil {
		t.Fatalf("create Console Log file: %v", err)
	}
	t.Cleanup(func() {
		_ = consoleLog.Close()
	})
	ui, err := startRunUI(ctx, roundtui.LiveRunView{}, run.ID, homeDir, runStore, consoleLog, false)
	if err != nil {
		t.Fatalf("start non-TTY Run UI: %v", err)
	}

	payload := compactCLIEditLifecyclePayload("tool_call_update", "edit_call_1", "internal/app/server.go", "old\n", "new\n", "running")
	events := []runevent.RunEvent{
		cliToolLifecycleEvent(run.ID, runevent.KindAgentToolStarted, "edit_call_1", "running", payload),
		cliToolLifecycleEvent(run.ID, runevent.KindAgentToolUpdated, "edit_call_1", "completed", payload),
	}
	for _, event := range events {
		if err := ui.sink.Publish(ctx, event); err != nil {
			t.Fatalf("publish Run Event through non-TTY UI: %v", err)
		}
	}
	ui.Close()

	if err := consoleLog.Close(); err != nil {
		t.Fatalf("close Console Log file: %v", err)
	}
	consoleBytes, err := os.ReadFile(consoleLog.Name())
	if err != nil {
		t.Fatalf("read Console Log file: %v", err)
	}
	if string(consoleBytes) != "edit internal/app/server.go (+1/-1)\n" {
		t.Fatalf("expected non-TTY Console Log writer to collapse duplicate summary, got %q", string(consoleBytes))
	}
	journaled, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read journaled Run Events: %v", err)
	}
	if len(journaled) != 2 {
		t.Fatalf("expected both lifecycle events to stay journaled, got %d", len(journaled))
	}
	for index, entry := range journaled {
		if entry.Event.Kind != events[index].Kind {
			t.Fatalf("expected journal kind %d unchanged, got %s", index, entry.Event.Kind)
		}
		if entry.Event.ToolID != events[index].ToolID {
			t.Fatalf("expected journal tool id %d unchanged, got %q", index, entry.Event.ToolID)
		}
		if entry.Event.ToolState != events[index].ToolState {
			t.Fatalf("expected journal tool state %d unchanged, got %q", index, entry.Event.ToolState)
		}
		if !bytes.Equal(entry.Event.Payload, events[index].Payload) {
			t.Fatalf("expected journal payload %d unchanged\nwant: %s\ngot:  %s", index, events[index].Payload, entry.Event.Payload)
		}
	}
	if journaled[0].Cursor == 0 || journaled[1].Cursor == 0 || journaled[0].Cursor == journaled[1].Cursor {
		t.Fatalf("expected distinct non-zero journal cursors, got %d and %d", journaled[0].Cursor, journaled[1].Cursor)
	}
	if !bytes.Equal(journaled[0].Event.Payload, journaled[1].Event.Payload) {
		t.Fatalf("expected sanitized lifecycle pair to keep byte-identical payloads")
	}

	timeline := roundtui.NewRunTimeline(10)
	cursor, err := replayRunEvents(ctx, runStore, run.ID, 0, timeline)
	if err != nil {
		t.Fatalf("replay Run Events into Live Run View timeline: %v", err)
	}
	if cursor != journaled[1].Cursor {
		t.Fatalf("expected replay cursor %d, got %d", journaled[1].Cursor, cursor)
	}
	wantLines := []string{
		"edit internal/app/server.go (+1/-1)",
		"edit internal/app/server.go (+1/-1)",
	}
	if got := timeline.Lines(); !reflect.DeepEqual(got, wantLines) {
		t.Fatalf("expected Attach/Live Run View replay to retain duplicate events\nwant: %#v\n got: %#v", wantLines, got)
	}
}

func TestAgentConsoleDisplaySinkKeepsNoAgentConsoleSuppression(t *testing.T) {
	event := runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Payload: []byte(`{"text":"fake agent output\n"}`),
	}
	var buffer bytes.Buffer

	if err := agentConsoleDisplaySink(&buffer, true).Publish(context.Background(), event); err != nil {
		t.Fatalf("publish through suppressed display sink: %v", err)
	}

	if buffer.Len() != 0 {
		t.Fatalf("expected --no-agent-console to suppress Agent output, got %q", buffer.String())
	}
}

func cliToolLifecycleEvent(runID string, kind runevent.Kind, toolID string, toolState string, payload []byte) runevent.RunEvent {
	return runevent.RunEvent{
		RunID:     runID,
		Batch:     1,
		Source:    runevent.SourceAgent,
		Kind:      kind,
		ToolID:    toolID,
		ToolState: toolState,
		Payload:   payload,
	}
}

func compactCLIEditLifecyclePayload(sessionUpdate string, toolID string, path string, oldText string, newText string, status string) []byte {
	return []byte(`{"sessionId":"s","update":{"sessionUpdate":` + strconv.Quote(sessionUpdate) + `,"toolCallId":` + strconv.Quote(toolID) + `,"kind":"edit","title":"edit","status":` + strconv.Quote(status) + `,"locations":[{"path":` + strconv.Quote(path) + `}],"content":[{"type":"diff","path":` + strconv.Quote(path) + `,"oldText":` + strconv.Quote(oldText) + `,"newText":` + strconv.Quote(newText) + `}]}}`)
}

func assertJournalContainsAgentAndDaemonEvents(t *testing.T, events []store.JournalEvent, agentText string) {
	t.Helper()
	agentSeen := false
	daemonSeen := false
	for _, entry := range events {
		switch entry.Event.Source {
		case runevent.SourceAgent:
			if strings.Contains(entry.Event.Summary, agentText) || strings.Contains(string(entry.Event.Payload), agentText) {
				agentSeen = true
			}
		case runevent.SourceDaemon:
			daemonSeen = true
		}
	}
	if !agentSeen {
		t.Fatalf("expected Agent-source event containing %q in the Run Event Journal, got %+v", agentText, events)
	}
	if !daemonSeen {
		t.Fatalf("expected Daemon-source events in the Run Event Journal, got %+v", events)
	}
}

func TestResolveJournalsAgentRunEventsDurably(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeAgentRunner{})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	if len(events) == 0 {
		t.Fatal("expected journaled Run Events after the Resolve Run")
	}
	var event runevent.RunEvent
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindAgentRaw {
			event = entry.Event
			break
		}
	}
	if event.Kind != runevent.KindAgentRaw || event.Source != runevent.SourceAgent {
		t.Fatalf("expected agent.raw journal entry, got %+v", events)
	}
	if event.RunID != runID || event.Batch != 1 {
		t.Fatalf("expected Run identity on journaled event, got %+v", event)
	}
	if !strings.Contains(event.Summary, "fake agent output") {
		t.Fatalf("expected bounded summary, got %q", event.Summary)
	}
	if !strings.Contains(string(event.Payload), "fake agent output") {
		t.Fatalf("expected raw payload preserved, got %s", event.Payload)
	}
	if !strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected console output unchanged, got %q", stderr.String())
	}
}

func TestRunResolveSkipsAgentLogFilesByDefaultAndStillJournals(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	assertNoAgentLogs(t, repoDir)
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunResolveWritesAgentLogFilesWhenEnabledAndStillJournals(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
logs:
  agent: true
`)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	assertAgentLogContains(t, repoDir, "fake agent output")
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunResolveNoAgentConsoleSuppressesAgentDisplayOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeAgentRunner{})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"resolve selected 1 downloaded Unresolved Review Issue",
		"Batch: 001/001 (1 Review Issue(s))",
		"Verification passed (attempt 1).",
		"Batch commit created",
		"Resolved 1 Review Source thread",
		"Resolve Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to keep daemon/progress line %q, got %q", want, stderr.String())
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunWatchNoAgentConsoleSuppressesAgentDisplayOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"Watch Run:",
		"Review Source status: settled",
		"Fetched Round 001 with 1 Review Issue",
		"Batch: 001/001 (1 Review Issue(s))",
		"Verification passed (attempt 1).",
		"Final Push completed",
		"Watch Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to keep daemon/progress line %q, got %q", want, stderr.String())
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestStoppedResolveJournalsStoppedEventBeforeReturning(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPathSnapshots(t, nil, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	stoppedSeen := false
	for _, entry := range events {
		if entry.Event.Kind != runevent.KindAgentStatus {
			continue
		}
		if strings.Contains(string(entry.Event.Payload), `"stopped"`) {
			stoppedSeen = true
		}
	}
	if !stoppedSeen {
		t.Fatalf("expected stopped status event journaled, got %+v", events)
	}
}

func TestWatchSkipsFinalPushWhenAutoPushDisabled(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pusher := &fakePusher{}
	withPusher(t, pusher)
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{
			Git: preflight.GitState{
				Root:            repoDir,
				Branch:          "feature/review",
				HEAD:            "abc123",
				UpstreamRemote:  "origin",
				UpstreamBranch:  "feature/review",
				UnpushedCommits: 1,
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{Enabled: false},
		}, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 0 {
		t.Fatalf("expected Final Push gated off when auto-push is disabled, got %d push(es)", pusher.calls)
	}
	if !strings.Contains(stderr.String(), "Final Push skipped: auto-push disabled or no push target configured.") {
		t.Fatalf("expected gating message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Clean") {
		t.Fatalf("expected Clean watch outcome without push, got %q", stderr.String())
	}
}

func TestWatchFinalPushRunsOncePerCleanRoundThroughEngine(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pusher := &fakePusher{}
	withPusher(t, pusher)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 1 {
		t.Fatalf("expected exactly one Final Push for the clean Round, got %d", pusher.calls)
	}
	if pusher.remotes[0] != "origin" || pusher.branches[0] != "feature/review" {
		t.Fatalf("expected push target from the push plan, got %v %v", pusher.remotes, pusher.branches)
	}
	if !strings.Contains(stderr.String(), "Batch commit created") {
		t.Fatalf("expected engine-driven Batch commit in watch Round, got %q", stderr.String())
	}
}

func runResolveForAttachTest(t *testing.T, repoDir string) (string, *bytes.Buffer) {
	t.Helper()
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("seed resolve run failed: %d stderr=%q", code, stderr.String())
	}
	runID := ""
	for _, line := range strings.Split(stderr.String(), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Run id in resolve output, got %q", stderr.String())
	}
	return runID, &stderr
}

func TestAttachReplaysCompletedRunReadOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	// Attach must never probe or start an Agent; a probing attach would
	// fail loudly here.
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("attach must not probe Agents")})
	dbPath := filepath.Join(homeDir, ".roundfix", "roundfix.db")
	assertRunCount(t, dbPath, 1)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Roundfix attach",
		"ID: " + runID,
		"State: Clean",
		"fake agent output",
		"Run " + runID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
	// Non-mutating: no new Runs, and the Run's terminal state is untouched.
	assertRunCount(t, dbPath, 1)
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	run, found, err := reader.Run(context.Background(), runID)
	if err != nil || !found {
		t.Fatalf("lookup run after attach: found=%v err=%v", found, err)
	}
	if run.State != store.StateClean {
		t.Fatalf("expected attach to leave Run state untouched, got %q", run.State)
	}
}

func TestAttachDisplaysStoredSelectionAfterConfigChanges(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var resolveStdout bytes.Buffer
	var resolveStderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"resolve",
		"--pr", "123",
		"--agent", "codex",
		"--model", "historical-model",
		"--reasoning-effort", "historical-reasoning",
		"--no-input",
	}, &resolveStdout, &resolveStderr)
	if code != exitOK {
		t.Fatalf("seed resolve run failed: %d stderr=%q", code, resolveStderr.String())
	}
	runID := reviewRunIDFromStderr(t, resolveStderr.String())
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
runtimes:
  codex:
    model: changed-config-model
    reasoning_effort: changed-config-reasoning
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code = RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected attach exit 0, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Agent: Codex",
		"Agent Model: historical-model",
		"Default Reasoning Effort: historical-reasoning",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain stored value %q, got:\n%s", expected, output)
		}
	}
	for _, unexpected := range []string{"changed-config-model", "changed-config-reasoning", "auto"} {
		if strings.Contains(output, unexpected) {
			t.Fatalf("expected attach output not to contain %q, got:\n%s", unexpected, output)
		}
	}
	if run := runFromStore(t, homeDir, runID); run.Model != "historical-model" || run.ReasoningEffort != "historical-reasoning" {
		t.Fatalf("expected stored historical selection to remain unchanged, got %#v", run)
	}
}

func TestAgentSelectionAttachReplayRendersPerScopeSelectionState(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open run store: %v", err)
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:            store.KindResolve,
		HeadRepository:  "owner/project",
		HeadBranch:      "feature/review",
		BaseRepository:  "owner/project",
		PRNumber:        "123",
		GitRoot:         repoDir,
		LocalBranch:     "feature/review",
		HeadSHA:         "abc123",
		ArtifactDir:     filepath.Join(repoDir, ".roundfix", "reviews"),
		Agent:           "Codex",
		Model:           "compat-summary-model",
		ReasoningEffort: "compat-summary-reasoning",
	})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	fallback := runevent.RunEvent{
		RunID:  run.ID,
		Batch:  1,
		Source: runevent.SourceDaemon,
		Kind:   runevent.KindDaemonAgentSelectionFallback,
		Payload: []byte(`{
			"event":"agent_selection_fallback",
			"category":"review",
			"scope_kind":"review",
			"scope_id":"batch-001",
			"scope_identity":"review:batch-001",
			"failed_selection":{"runtime":"codex","model":"gpt-5.6-sol","reasoning_effort":"high"},
			"next_selection":{"runtime":"claude","model":"claude-fable-5","reasoning_effort":"xhigh"},
			"fallback_index":1,
			"reason_code":"runtime_unavailable",
			"reason":"runtime unavailable",
			"automatic":true
		}`),
	}
	active := runevent.RunEvent{
		RunID:  run.ID,
		Batch:  1,
		Source: runevent.SourceDaemon,
		Kind:   runevent.KindDaemonAgentSelectionActive,
		Payload: []byte(`{
			"event":"agent_selection_active",
			"scope_kind":"review",
			"scope_id":"batch-001",
			"scope_identity":"review:batch-001",
			"category":"review",
			"profile_source":"project",
			"attempt":2,
			"selection_role":"fallback",
			"fallback_index":1,
			"runtime":"claude",
			"model":"claude-fable-5",
			"reasoning_effort":"xhigh",
			"status":"active"
		}`),
	}
	if _, err := runStore.AppendRunEvents(ctx, []runevent.RunEvent{fallback, active}); err != nil {
		t.Fatalf("append selection events: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateClean); err != nil {
		t.Fatalf("complete run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close run store: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"attach", run.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected attach exit 0, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Agent Model: compat-summary-model",
		"Selections:",
		"review batch-001 (review) fallback failed source unavailable codex/gpt-5.6-sol/high -> fallback 1 claude/claude-fable-5/xhigh reason runtime_unavailable: runtime unavailable",
		"review batch-001 (review) attempt 2 fallback active source project claude/claude-fable-5/xhigh",
		"SELECTION review batch-001 (review) fallback failed...",
		"SELECTION review batch-001 (review) attempt 2 fallb...",
		"Run " + run.ID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach replay to contain %q, got:\n%s", expected, output)
		}
	}
	for _, forbidden := range []string{"prompt", "credential", "token", "cookie", "secret"} {
		if strings.Contains(strings.ToLower(output), forbidden) {
			t.Fatalf("expected attach replay to omit %q, got:\n%s", forbidden, output)
		}
	}
}

func TestFallbackNotificationOrderingSelectionConsoleSink(t *testing.T) {
	var stderr bytes.Buffer
	fanout := runevent.NewFanout([]runevent.Sink{
		selectionConsoleDisplaySink(&stderr),
		agent.NewConsoleDisplaySink(&stderr),
	}, nil)
	fallback := runevent.RunEvent{
		RunID:  "run_9",
		Source: runevent.SourceDaemon,
		Kind:   runevent.KindDaemonAgentSelectionFallback,
		Payload: []byte(`{
			"event":"agent_selection_fallback",
			"category":"backend",
			"scope_kind":"task",
			"scope_id":"task_01",
			"scope_identity":"task:task_01",
			"failed_selection":{"runtime":"codex","model":"gpt-5.6-sol","reasoning_effort":"high"},
			"next_selection":{"runtime":"claude","model":"claude-fable-5","reasoning_effort":"xhigh"},
			"fallback_index":1,
			"reason_code":"runtime_unavailable",
			"reason":"runtime unavailable",
			"automatic":true
		}`),
	}
	active := runevent.RunEvent{
		RunID:  "run_9",
		Source: runevent.SourceDaemon,
		Kind:   runevent.KindDaemonAgentSelectionActive,
		Payload: []byte(`{
			"event":"agent_selection_active",
			"scope_kind":"task",
			"scope_id":"task_01",
			"scope_identity":"task:task_01",
			"category":"backend",
			"profile_source":"project",
			"attempt":2,
			"selection_role":"fallback",
			"fallback_index":1,
			"runtime":"claude",
			"model":"claude-fable-5",
			"reasoning_effort":"xhigh",
			"status":"active"
		}`),
	}
	workStarted := runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Payload: []byte(`{"status":"agent_work_started"}`),
	}
	for _, event := range []runevent.RunEvent{fallback, active, workStarted} {
		if err := fanout.Publish(context.Background(), event); err != nil {
			t.Fatalf("publish event: %v", err)
		}
	}

	output := stderr.String()
	assertCLIContainsInOrder(t, output,
		"SELECTION task task_01 (backend) fallback failed source unavailable codex/gpt-5.6-sol/high -> fallback 1 claude/claude-fable-5/xhigh reason runtime_unavailable: runtime unavailable",
		"SELECTION task task_01 (backend) attempt 2 fallback active source project claude/claude-fable-5/xhigh",
		"SESSION AGENT_WORK_STARTED",
	)
}

func assertCLIContainsInOrder(t *testing.T, haystack string, needles ...string) {
	t.Helper()
	offset := 0
	for _, needle := range needles {
		index := strings.Index(haystack[offset:], needle)
		if index < 0 {
			t.Fatalf("expected %q after byte offset %d in:\n%s", needle, offset, haystack)
		}
		offset += index + len(needle)
	}
}

func TestAttachRunBrowserLoopOpensCockpitAndRefreshes(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	base := time.Date(2026, 7, 6, 14, 0, 0, 0, time.UTC)
	runs := seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:      store.KindResolve,
			state:     store.StateClean,
			gitRoot:   repoDir,
			branch:    "feature/terminal",
			prNumber:  "201",
			createdAt: base.Add(time.Minute),
		},
		{
			kind:      store.KindImplement,
			state:     store.StateVerifying,
			gitRoot:   repoDir,
			branch:    "ma/spec-active",
			specSlug:  "0017-run-discovery",
			createdAt: base.Add(2 * time.Minute),
		},
		{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   otherRepo,
			branch:    "feature/other",
			prNumber:  "999",
			createdAt: base.Add(3 * time.Minute),
		},
	})
	terminalRun, activeRun, otherRepoRun := runs[0], runs[1], runs[2]
	withAttachInteractiveInput(t, true)
	t.Setenv("ROUNDFIX_TUI", "always")
	var createdBehindCockpit store.Run
	cockpitCalls := withBrowserAttachCockpit(t, func(run store.Run, concurrency int) int {
		// Mutating the store while the cockpit is open proves the next
		// browser pass re-queries instead of reusing stale listings.
		createdBehindCockpit = seedRunsForList(t, homeDir, []runListSeed{{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    "feature/started-behind-cockpit",
			prNumber:  "202",
			createdAt: base.Add(4 * time.Minute),
		}})[0]
		return exitOK
	})
	sessionCalls := withRunBrowserSession(t,
		roundtui.BrowserOutcome{RunID: activeRun.ID},
		roundtui.BrowserOutcome{Cancelled: true},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected browser loop exit 0, got %d stderr=%q", code, stderr.String())
	}
	calls := *sessionCalls
	if len(calls) != 2 {
		t.Fatalf("expected browser sessions before and after the cockpit, got %d", len(calls))
	}
	first := calls[0]
	// The browser is machine-wide: Active Runs from every repository,
	// newest first.
	if len(first.active) != 2 || first.active[0].ID != otherRepoRun.ID || first.active[1].ID != activeRun.ID {
		t.Fatalf("expected every repository's Active Runs newest first, got %+v", first.active)
	}
	if len(first.all) != 3 || first.all[0].ID != otherRepoRun.ID || first.all[1].ID != activeRun.ID || first.all[2].ID != terminalRun.ID {
		t.Fatalf("expected every repository's Runs newest first, got %+v", first.all)
	}
	cockpit := *cockpitCalls
	if len(cockpit) != 1 || cockpit[0].runID != activeRun.ID {
		t.Fatalf("expected one cockpit pass for the selected Run, got %+v", cockpit)
	}
	if cockpit[0].concurrency != 2 {
		t.Fatalf("expected the explicit-attach concurrency fallback, got %d", cockpit[0].concurrency)
	}
	second := calls[1]
	if len(second.all) != 4 || second.all[0].ID != createdBehindCockpit.ID {
		t.Fatalf("expected the refreshed browser to list the Run created behind the cockpit, got %+v", second.all)
	}
}

func TestBrowserAttachCockpitIsTheExplicitAttachCockpit(t *testing.T) {
	if reflect.ValueOf(browserAttachCockpit).Pointer() != reflect.ValueOf(runAttachCockpit).Pointer() {
		t.Fatal("expected the browser loop to open the same attach cockpit as explicit attach")
	}
}

func TestAttachRunBrowserCancelExitsZeroWithoutAttaching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	dbPath := filepath.Join(homeDir, ".roundfix", "roundfix.db")
	assertRunCount(t, dbPath, 1)
	withAttachInteractiveInput(t, true)
	t.Setenv("ROUNDFIX_TUI", "always")
	cockpitCalls := withBrowserAttachCockpit(t, nil)
	withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected cancelled browser exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no attach output after browser cancel, got %q", stdout.String())
	}
	if len(*cockpitCalls) != 0 {
		t.Fatalf("expected no cockpit pass after cancel, got %+v", *cockpitCalls)
	}
	assertRunCount(t, dbPath, 1)
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	run, found, err := reader.Run(context.Background(), runID)
	if err != nil || !found {
		t.Fatalf("lookup run after browser cancel: found=%v err=%v", found, err)
	}
	if run.State != store.StateClean {
		t.Fatalf("expected browser cancel to leave Run untouched, got %q", run.State)
	}
}

func TestAttachWithoutRunIDNonInteractiveNamesAllRunsList(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		interactive bool
	}{
		{name: "no tty", args: []string{"attach"}, interactive: false},
		{name: "no input flag", args: []string{"attach", "--no-input"}, interactive: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCLIWorkspace(t)
			withAttachInteractiveInput(t, tt.interactive)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "roundfix runs list --all") {
				t.Fatalf("expected discovery command guidance, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "pass a run id") {
				t.Fatalf("expected run id requirement, got %q", stderr.String())
			}
		})
	}
}

func TestAttachSpecRunReadsTasksFromKeptWorkDir(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	workDir := t.TempDir()
	writeImplementSpec(t, repoDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Read state", status: "pending"}})
	writeImplementSpec(t, workDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Read state", status: "completed"}})
	run := createTerminalAttachSpecRun(t, homeDir, repoDir, workDir, store.StateUnresolved)
	appendAttachConcurrencyEvent(t, homeDir, run.ID, 4)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", run.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Run Worktree: " + workDir,
		"Concurrency: 4",
		"task_01 completed — Read state",
		"Run " + run.ID + " reached Unresolved; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
	if strings.Contains(output, "task_01 pending") {
		t.Fatalf("expected attach to ignore stale GitRoot task status, got:\n%s", output)
	}
}

func TestAttachSpecRunFallsBackToGitRootWhenCleanWorkDirIsPruned(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	workDir := filepath.Join(t.TempDir(), "pruned-worktree")
	writeImplementSpec(t, repoDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Integrated state", status: "completed"}})
	run := createTerminalAttachSpecRun(t, homeDir, repoDir, workDir, store.StateClean)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", run.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit for pruned worktree, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Run Worktree: " + workDir,
		"task_01 completed — Integrated state",
		"Run " + run.ID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestAttachUnknownRunFailsBeforeTUIStart(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", "41"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit for unknown Run, got %d", code)
	}
	want := `Run "41" does not exist; picker numbers are not stable Run ids — pass a run id or run 'roundfix attach' to pick interactively`
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected picker-number error %q, got %q", want, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no replay output before failure, got %q", stdout.String())
	}
}

func TestAttachWithoutRunDatabaseFailsAsCLIError(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", "run_x"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit without a Run Database, got %d", code)
	}
	if !strings.Contains(stderr.String(), "roundfix attach failed:") {
		t.Fatalf("expected CLI error before TUI start, got %q", stderr.String())
	}
}

func TestAttachSkipsUnknownEventKindsOnReplay(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	if _, err := writer.AppendRunEvent(context.Background(), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    "future.unknown",
		Summary: "future event",
		Time:    time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"new":"shape"}`),
	}); err != nil {
		t.Fatalf("append future event: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected unknown kinds skipped, got exit %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "future event") {
		t.Fatalf("expected unknown kind not rendered, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "fake agent output") {
		t.Fatalf("expected known events still replayed, got %q", stdout.String())
	}
}

func createTerminalAttachSpecRun(t *testing.T, homeDir string, repoDir string, workDir string, state string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/widget-flow",
		HeadSHA:     "abc123",
		WorkDir:     workDir,
		SpecSlug:    implementTestSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create implement run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, state)
	if err != nil {
		t.Fatalf("complete implement run: %v", err)
	}
	return completed.Run
}

func appendAttachConcurrencyEvent(t *testing.T, homeDir string, runID string, concurrency int) {
	t.Helper()
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Fatalf("close writer: %v", err)
		}
	}()
	_, err = writer.AppendRunEvent(context.Background(), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: fmt.Sprintf("Task cycle started with concurrency %d.", concurrency),
		Payload: []byte(fmt.Sprintf(`{"concurrency":%d}`, concurrency)),
	})
	if err != nil {
		t.Fatalf("append concurrency event: %v", err)
	}
}

func TestEventsReplayDefaultAndFilterJSONLRecordsOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateClean)
	appendEvents(t, homeDir, run.ID,
		runevent.RunEvent{
			RunID:   run.ID,
			Batch:   2,
			Source:  runevent.SourceDaemon,
			Kind:    runevent.KindDaemonBatch,
			Summary: "Batch 002 started.",
			Payload: []byte(`{"phase":"started","batch":2}`),
		},
		runevent.RunEvent{
			RunID:   run.ID,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentRaw,
			Summary: "raw agent payload",
			Payload: []byte(`{"text":"agent output bytes"}`),
		},
		runevent.RunEvent{
			RunID:       run.ID,
			Batch:       2,
			Source:      runevent.SourceDaemon,
			Kind:        runevent.KindDaemonVerification,
			ReviewIssue: "task_03",
			Summary:     "Verification attempt 1 failed: make verify",
			Payload:     []byte(`{"attempt":1,"phase":"verdict","verdict":"failed","command":"make verify","diagnostic_path":"/tmp/verification.log"}`),
		},
		runevent.RunEvent{
			RunID:       run.ID,
			Batch:       2,
			Source:      runevent.SourceDaemon,
			Kind:        runevent.KindDaemonTask,
			ReviewIssue: "task_03",
			Summary:     "Task task_03 settled completed.",
			Payload:     []byte(`{"task":"task_03","phase":"settled","status":"completed"}`),
		},
		runevent.RunEvent{
			RunID:   run.ID,
			Source:  runevent.SourceDaemon,
			Kind:    runevent.KindDaemonOutcome,
			Summary: "Run reached Clean.",
			Payload: []byte(`{"state":"Clean","remaining":0}`),
		},
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"events", run.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected replay exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for replay, got %q", stderr.String())
	}
	records := decodeEventRecords(t, stdout.String())
	if got := eventRecordCategories(records); strings.Join(got, "|") != "batch|verification|task-status|outcome" {
		t.Fatalf("expected stable categories in cursor order, got %v records=%v", got, records)
	}
	for _, raw := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		for _, forbidden := range []string{"daemon.", "agent output bytes", "make verify", "diagnostic_path", "/tmp/verification.log"} {
			if strings.Contains(raw, forbidden) {
				t.Fatalf("expected records to omit %q, got %s", forbidden, raw)
			}
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{"events", run.ID, "--filter", "verification,outcome"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected filtered replay exit 0, got %d stderr=%q", code, stderr.String())
	}
	records = decodeEventRecords(t, stdout.String())
	if got := eventRecordCategories(records); strings.Join(got, "|") != "verification|outcome" {
		t.Fatalf("expected filtered categories only, got %v records=%v", got, records)
	}
}

func TestEventsReplayLegacyVerificationEvent(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateClean)
	appendEvents(t, homeDir, run.ID, runevent.RunEvent{
		RunID:   run.ID,
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonVerification,
		Summary: "Verification command passed: make verify",
		Payload: []byte(`{"phase":"passed","command":"make verify"}`),
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"events", run.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected legacy replay exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr for legacy replay, got %q", stderr.String())
	}
	records := decodeEventRecords(t, stdout.String())
	if len(records) != 1 || records[0]["category"] != "verification" || records[0]["attempt"] != float64(1) {
		t.Fatalf("expected one normalized legacy verification record, got %v", records)
	}
}

func TestEventsFollowDrainsTerminalWithoutDuplicateBoundary(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateActive)
	appendEvents(t, homeDir, run.ID, runevent.RunEvent{
		RunID:   run.ID,
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonBatch,
		Summary: "Batch 001 started.",
		Payload: []byte(`{"phase":"started","batch":1}`),
	})

	appendedTerminal := false
	withAttachSleep(t, func(context.Context) error {
		if !appendedTerminal {
			appendEvents(t, homeDir, run.ID, runevent.RunEvent{
				RunID:   run.ID,
				Source:  runevent.SourceDaemon,
				Kind:    runevent.KindDaemonOutcome,
				Summary: "Run reached Clean.",
				Payload: []byte(`{"state":"Clean","remaining":0}`),
			})
			completeEventsRun(t, homeDir, run.ID, store.StateClean)
			appendedTerminal = true
		}
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"events", run.ID, "--follow"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected follow exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	records := decodeEventRecords(t, stdout.String())
	if got := eventRecordCategories(records); strings.Join(got, "|") != "batch|outcome" {
		t.Fatalf("expected replay boundary and terminal event once, got %v records=%v", got, records)
	}
	if records[0]["cursor"] == records[1]["cursor"] {
		t.Fatalf("expected distinct cursors without duplicate boundary records, got %v", records)
	}
}

func TestEventsTerminalRunReplaysAndExitsImmediately(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateClean)
	appendEvents(t, homeDir, run.ID, runevent.RunEvent{
		RunID:   run.ID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: "Run reached Clean.",
		Payload: []byte(`{"state":"Clean","remaining":0}`),
	})
	withAttachSleep(t, func(context.Context) error {
		t.Fatal("terminal events follow must not poll")
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"events", run.ID, "--follow"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected terminal follow exit 0, got %d stderr=%q", code, stderr.String())
	}
	if got := eventRecordCategories(decodeEventRecords(t, stdout.String())); strings.Join(got, "|") != "outcome" {
		t.Fatalf("expected terminal replay only, got %v", got)
	}
}

func TestEventsValidationErrorsEmitNoStdout(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateClean)
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing run id", args: []string{"events"}},
		{name: "unknown run", args: []string{"events", "run_missing"}},
		{name: "invalid filter", args: []string{"events", run.ID, "--filter", "raw"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected validation exit 2, got %d stderr=%q", code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout records, got %q", stdout.String())
			}
			if stderr.Len() == 0 {
				t.Fatal("expected validation diagnostic on stderr")
			}
		})
	}
}

func TestEventsMalformedRelevantPayloadFailsNoStdout(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateClean)
	appendEvents(t, homeDir, run.ID, runevent.RunEvent{
		RunID:   run.ID,
		Batch:   1,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonBatch,
		Summary: "Batch 001 started.",
		Payload: []byte(`{"batch":1}`),
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"events", run.ID}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected malformed payload exit 1, got %d stderr=%q", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout records, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "roundfix events failed") {
		t.Fatalf("expected events failure diagnostic, got %q", stderr.String())
	}
}

func TestEventsFollowCancellationExits130WithoutTrailer(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	run := createEventsRun(t, homeDir, repoDir, store.StateActive)
	ctx, cancel := context.WithCancel(context.Background())
	withAttachSleep(t, func(sleepCtx context.Context) error {
		cancel()
		return sleepCtx.Err()
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"events", run.ID, "--follow"}, &stdout, &stderr)

	if code != exitSIGINT {
		t.Fatalf("expected cancellation exit 130, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout trailer on cancellation, got %q", stdout.String())
	}
}

func createEventsRun(t *testing.T, homeDir string, repoDir string, state string) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open events store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close events store: %v", err)
		}
	}()
	run, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:            store.KindImplement,
		GitRoot:         repoDir,
		LocalBranch:     "ma/event-stream",
		HeadSHA:         "abc123",
		WorkDir:         repoDir,
		ArtifactDir:     filepath.Join(repoDir, ".roundfix"),
		SpecSlug:        "0024-context-efficient-runs",
		Agent:           "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
	})
	if err != nil {
		t.Fatalf("create events Run: %v", err)
	}
	if state == "" || state == store.StateActive {
		return run
	}
	return completeEventsRun(t, homeDir, run.ID, state)
}

func completeEventsRun(t *testing.T, homeDir string, runID string, state string) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open events store to complete Run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close events complete store: %v", err)
		}
	}()
	run, err := runStore.CompleteRun(context.Background(), runID, state)
	if err != nil {
		t.Fatalf("complete events Run: %v", err)
	}
	return run.Run
}

func appendEvents(t *testing.T, homeDir string, runID string, events ...runevent.RunEvent) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open events store to append: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close events append store: %v", err)
		}
	}()
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	for index := range events {
		events[index].RunID = runID
		if events[index].Time.IsZero() {
			events[index].Time = now.Add(time.Duration(index) * time.Second)
		}
	}
	if _, err := runStore.AppendRunEvents(context.Background(), events); err != nil {
		t.Fatalf("append events: %v", err)
	}
}

func decodeEventRecords(t *testing.T, output string) []map[string]any {
	t.Helper()
	if strings.TrimSpace(output) == "" {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	records := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode JSONL record %q: %v", line, err)
		}
		if record["schema"] != runevent.StreamSchema {
			t.Fatalf("expected schema %q, got %v in %v", runevent.StreamSchema, record["schema"], record)
		}
		records = append(records, record)
	}
	return records
}

func eventRecordCategories(records []map[string]any) []string {
	categories := make([]string, 0, len(records))
	for _, record := range records {
		if category, ok := record["category"].(string); ok {
			categories = append(categories, category)
		}
	}
	return categories
}

type fakeAttachSource struct {
	mu          sync.Mutex
	run         store.Run
	events      []store.JournalEvent
	version     int64
	eventReads  int
	versionRead int
}

func (source *fakeAttachSource) Run(context.Context, string) (store.Run, bool, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	return source.run, true, nil
}

func (source *fakeAttachSource) RunEventsAfter(_ context.Context, _ string, cursor int64, limit int) ([]store.JournalEvent, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.eventReads++
	page := []store.JournalEvent{}
	for _, entry := range source.events {
		if entry.Cursor > cursor && len(page) < limit {
			page = append(page, entry)
		}
	}
	return page, nil
}

func (source *fakeAttachSource) DataVersion(context.Context) (int64, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.versionRead++
	return source.version, nil
}

func (source *fakeAttachSource) appendEvent(summary string) {
	source.mu.Lock()
	defer source.mu.Unlock()
	cursor := int64(len(source.events) + 1)
	source.events = append(source.events, store.JournalEvent{
		Cursor: cursor,
		Event: runevent.RunEvent{
			RunID:   source.run.ID,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentRaw,
			Summary: summary,
			Payload: []byte(`{"text":` + strconv.Quote(summary) + `}`),
		},
	})
	source.version++
}

func (source *fakeAttachSource) complete(state string) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.run.State = state
	source.version++
}

func TestAttachFollowerAppendsOnlyNewerEventsWithoutDuplicates(t *testing.T) {
	source := &fakeAttachSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	source.appendEvent("backlog one\n")
	source.appendEvent("backlog two\n")
	accepted := []string{}
	steps := 0
	follower := attachFollower{
		source: source,
		sleep: func(context.Context) error {
			steps++
			switch steps {
			case 1:
				source.appendEvent("live three\n")
			case 2:
				source.complete(store.StateClean)
			}
			return nil
		},
		accept: func(entry store.JournalEvent) error {
			accepted = append(accepted, entry.Event.Summary)
			return nil
		},
	}

	final, cursor, err := follower.follow(context.Background(), "run-1", 2)

	if err != nil {
		t.Fatalf("follow: %v", err)
	}
	if final.State != store.StateClean {
		t.Fatalf("expected terminal run returned, got %+v", final)
	}
	if strings.Join(accepted, "") != "live three\n" {
		t.Fatalf("expected only events newer than the replay cursor, got %v", accepted)
	}
	if cursor != 3 {
		t.Fatalf("expected cursor advanced to 3, got %d", cursor)
	}

	// Reconnect from the returned cursor: nothing replays twice.
	reaccepted := []string{}
	again := attachFollower{
		source: source,
		sleep:  func(context.Context) error { return nil },
		accept: func(entry store.JournalEvent) error {
			reaccepted = append(reaccepted, entry.Event.Summary)
			return nil
		},
	}
	if _, _, err := again.follow(context.Background(), "run-1", cursor); err != nil {
		t.Fatalf("reconnect follow: %v", err)
	}
	if len(reaccepted) != 0 {
		t.Fatalf("expected no duplicates across reconnects, got %v", reaccepted)
	}
}

func TestAttachFollowerIdlePollsReadNoEventRows(t *testing.T) {
	source := &fakeAttachSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 7}
	steps := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	follower := attachFollower{
		source: source,
		sleep: func(context.Context) error {
			steps++
			if steps >= 5 {
				cancel()
			}
			return ctx.Err()
		},
		accept: func(store.JournalEvent) error { return nil },
	}

	_, _, err := follower.follow(ctx, "run-1", 0)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected detach via cancellation, got %v", err)
	}
	if source.eventReads != 1 {
		t.Fatalf("expected only the initial drain to read rows; idle polls must not, got %d reads", source.eventReads)
	}
	if source.versionRead < 5 {
		t.Fatalf("expected change-detection polls to continue, got %d", source.versionRead)
	}
}

func TestAttachDetachLeavesRunActive(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	active, err := writer.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/other",
		BaseRepository: "owner/project",
		PRNumber:       "456",
		GitRoot:        repoDir,
		LocalBranch:    "feature/other",
		HeadSHA:        "def456",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Detach on the first follow poll, after backlog replay.
	withAttachSleep(t, func(sleepCtx context.Context) error {
		cancel()
		return sleepCtx.Err()
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"attach", active.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean detach exit, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Following live events.") {
		t.Fatalf("expected live-follow status visible, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Detached; Run "+active.ID+" keeps going.") {
		t.Fatalf("expected detach message, got %q", stdout.String())
	}
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	after, found, err := reader.Run(context.Background(), active.ID)
	if err != nil || !found {
		t.Fatalf("lookup run after detach: found=%v err=%v", found, err)
	}
	if store.IsTerminalState(after.State) {
		t.Fatalf("expected detach to leave the Run active, got %q", after.State)
	}
}

func TestAttachFollowsLiveRunToTerminalState(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	t.Cleanup(func() { _ = writer.Close() })
	active, err := writer.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/other",
		BaseRepository: "owner/project",
		PRNumber:       "456",
		GitRoot:        repoDir,
		LocalBranch:    "feature/other",
		HeadSHA:        "def456",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	withAttachSleep(t, func(ctx context.Context) error { return ctx.Err() })

	writerDone := make(chan error, 1)
	go func() {
		for index := 0; index < 5; index++ {
			if _, err := writer.AppendRunEvent(context.Background(), runevent.RunEvent{
				RunID:   active.ID,
				Source:  runevent.SourceAgent,
				Kind:    runevent.KindAgentRaw,
				Summary: fmt.Sprintf("live line %d\n", index),
				Payload: []byte(fmt.Sprintf(`{"text":"live line %d\n"}`, index)),
			}); err != nil {
				writerDone <- err
				return
			}
		}
		_, err := writer.CompleteRun(context.Background(), active.ID, store.StateStopped)
		writerDone <- err
	}()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"attach", active.ID}, &stdout, &stderr)

	if err := <-writerDone; err != nil {
		t.Fatalf("writer goroutine: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected clean follow exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for index := 0; index < 5; index++ {
		needle := fmt.Sprintf("live line %d", index)
		if count := strings.Count(output, needle); count != 1 {
			t.Fatalf("expected %q exactly once, got %d in:\n%s", needle, count, output)
		}
	}
	if !strings.Contains(output, "Run "+active.ID+" reached Stopped.") {
		t.Fatalf("expected terminal follow exit message, got %q", output)
	}
}

func withAttachSleep(t *testing.T, sleep func(ctx context.Context) error) {
	t.Helper()
	old := attachSleep
	attachSleep = sleep
	t.Cleanup(func() {
		attachSleep = old
	})
}

type browserSessionCall struct {
	active []store.Run
	all    []store.Run
}

// withRunBrowserSession scripts Run Browser outcomes so entry-point tests
// drive the browser loop through the model seam, not terminal emulation.
func withRunBrowserSession(t *testing.T, outcomes ...roundtui.BrowserOutcome) *[]browserSessionCall {
	t.Helper()
	calls := &[]browserSessionCall{}
	old := runBrowserSession
	runBrowserSession = func(_ context.Context, _ io.Writer, active, all []store.Run) (roundtui.BrowserOutcome, error) {
		index := len(*calls)
		*calls = append(*calls, browserSessionCall{active: active, all: all})
		if index >= len(outcomes) {
			t.Fatalf("unexpected Run Browser session call %d", index+1)
		}
		return outcomes[index], nil
	}
	t.Cleanup(func() {
		runBrowserSession = old
	})
	return calls
}

type browserCockpitCall struct {
	runID       string
	concurrency int
}

// withBrowserAttachCockpit replaces the browser loop's cockpit step and
// records each pass; a nil handler just reports success.
func withBrowserAttachCockpit(t *testing.T, handler func(run store.Run, concurrency int) int) *[]browserCockpitCall {
	t.Helper()
	calls := &[]browserCockpitCall{}
	old := browserAttachCockpit
	browserAttachCockpit = func(_ context.Context, _ roundconfig.Loaded, _ *store.Store, run store.Run, concurrency int, _ io.Writer, _ io.Writer) int {
		*calls = append(*calls, browserCockpitCall{runID: run.ID, concurrency: concurrency})
		if handler == nil {
			return exitOK
		}
		return handler(run, concurrency)
	}
	t.Cleanup(func() {
		browserAttachCockpit = old
	})
	return calls
}

func withAttachInteractiveInput(t *testing.T, available bool) {
	t.Helper()
	old := attachInteractiveInputAvailable
	attachInteractiveInputAvailable = func() bool {
		return available
	}
	t.Cleanup(func() {
		attachInteractiveInputAvailable = old
	})
}

func daemonKindSequence(events []store.JournalEvent) []string {
	kinds := []string{}
	for _, entry := range events {
		if entry.Event.Source != runevent.SourceDaemon {
			continue
		}
		kinds = append(kinds, string(entry.Event.Kind))
	}
	return kinds
}

func assertOrderedSubsequence(t *testing.T, haystack []string, expected []string) {
	t.Helper()
	index := 0
	for _, item := range haystack {
		if index < len(expected) && item == expected[index] {
			index++
		}
	}
	if index != len(expected) {
		t.Fatalf("expected ordered subsequence %v in %v (matched %d)", expected, haystack, index)
	}
}

func TestWatchRunJournalsOrderedLoopNarrative(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	kinds := daemonKindSequence(events)
	assertOrderedSubsequence(t, kinds, []string{
		string(runevent.KindDaemonReviewStatus),
		string(runevent.KindDaemonFetch),
		string(runevent.KindDaemonFetch),
		string(runevent.KindDaemonSelection),
		string(runevent.KindDaemonBatch),
		string(runevent.KindDaemonVerification),
		string(runevent.KindDaemonVerification),
		string(runevent.KindDaemonCommit),
		string(runevent.KindDaemonSourceResolution),
		string(runevent.KindDaemonBatch),
		string(runevent.KindDaemonPush),
		string(runevent.KindDaemonOutcome),
	})
	pushed := 0
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonPush && strings.Contains(string(entry.Event.Payload), `"pushed"`) {
			pushed++
		}
	}
	if pushed != 1 {
		t.Fatalf("expected exactly one pushed decision on the clean Round, got %d", pushed)
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Clean"`) {
		t.Fatalf("expected Clean outcome event last, got %+v", last)
	}
}

func TestStoppedRunJournalsStopWithoutLaterUnsafeDaemonEvents(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withChangedPathSnapshots(t, nil, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit, got %d", code)
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	stopIndex := -1
	for index, entry := range events {
		if entry.Event.Kind == runevent.KindAgentStatus && strings.Contains(string(entry.Event.Payload), `"stopped"`) {
			stopIndex = index
		}
	}
	if stopIndex < 0 {
		t.Fatalf("expected stop event journaled, got %+v", events)
	}
	for _, entry := range events[stopIndex+1:] {
		switch entry.Event.Kind {
		case runevent.KindDaemonVerification, runevent.KindDaemonCommit, runevent.KindDaemonPush, runevent.KindDaemonSourceResolution, runevent.KindDaemonFetch:
			t.Fatalf("expected no unsafe daemon events after stop, got %+v", entry.Event)
		}
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Stopped"`) {
		t.Fatalf("expected Stopped outcome event last, got %+v", last)
	}
}

func TestFailedVerificationJournalsFailureWithoutCommitEvents(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{err: errors.New("tests failed")})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code == 0 {
		t.Fatal("expected failed verification to fail the Run")
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	failed := false
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonVerification && strings.Contains(string(entry.Event.Payload), `"failed"`) {
			failed = true
		}
		if entry.Event.Kind == runevent.KindDaemonCommit {
			t.Fatalf("expected no commit events after failed verification, got %+v", entry.Event)
		}
		if entry.Event.Kind == runevent.KindDaemonPush && strings.Contains(string(entry.Event.Payload), `"pushed"`) {
			t.Fatalf("expected no push after failed verification, got %+v", entry.Event)
		}
	}
	if !failed {
		t.Fatalf("expected verification failure event journaled, got %+v", events)
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Unresolved"`) {
		t.Fatalf("expected Unresolved outcome event last, got %+v", last)
	}
}

func TestTriageOnlyBatchJournalsCommitSkipDecision(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	// Identical snapshots: the Agent triaged without touching the worktree.
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.worktree = staticWorktree{}
	})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected triage-only Batch to succeed, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	skipped := false
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonCommit && strings.Contains(string(entry.Event.Payload), `"skipped"`) {
			skipped = true
		}
	}
	if !skipped {
		t.Fatalf("expected commit-skip decision journaled, got %v", daemonKindSequence(events))
	}
}

type staticWorktree struct{}

func (staticWorktree) Snapshot(context.Context, string) ([]string, error) {
	return []string{"user-wip.txt"}, nil
}

func TestAttachRendersWatchDaemonEventsInTimeline(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var watchStdout bytes.Buffer
	var watchStderr bytes.Buffer
	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "6", "--no-input"}, &watchStdout, &watchStderr)
	if code != 0 {
		t.Fatalf("seed watch run failed: %d stderr=%q", code, watchStderr.String())
	}
	runID := ""
	for _, line := range strings.Split(watchStderr.String(), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "Watch Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Watch Run id, got %q", watchStderr.String())
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	attachCode := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if attachCode != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", attachCode, stderr.String())
	}
	for _, expected := range []string{
		"Review Source status: settled",
		"Verification attempt 1 for Batch 001 verdict: passed",
		"Batch commit created:",
		"Final Push completed:",
		"Run reached Clean.",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("expected attach timeline to narrate %q, got:\n%s", expected, stdout.String())
		}
	}
}

func TestResolvePrintsIssueSummaryAfterCompletion(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"This Run (Clean after 1 Round(s)): 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n" +
		"Pull Request cumulative: 1 resolved, 0 invalid, 0 duplicated, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Resolve Run") || !strings.Contains(stderr.String(), "reached Clean") {
		t.Fatalf("expected terminal daemon diagnostics on stderr, got %q", stderr.String())
	}
}

func TestStageableReviewRootClassifiesInsideOutsideAndSymlink(t *testing.T) {
	repoDir := t.TempDir()
	external := t.TempDir()
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-9"))
	mustMkdir(t, filepath.Join(external, "specs", "_reviews", "pr-9"))

	relative, ok := stageableReviewRoot(repoDir, filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-9"))
	if !ok || relative != "docs/specs/_reviews/pr-9" {
		t.Fatalf("expected inside root to stage as docs/specs/_reviews/pr-9, got %q ok=%v", relative, ok)
	}
	if _, ok := stageableReviewRoot(repoDir, filepath.Join(external, "specs", "_reviews", "pr-9")); ok {
		t.Fatal("expected external root to be unstageable")
	}
	linkRepo := t.TempDir()
	mustMkdir(t, filepath.Join(linkRepo, "docs"))
	if err := os.Symlink(filepath.Join(external, "specs"), filepath.Join(linkRepo, "docs", "specs")); err != nil {
		t.Fatalf("create specs symlink: %v", err)
	}
	if _, ok := stageableReviewRoot(linkRepo, filepath.Join(linkRepo, "docs", "specs", "_reviews", "pr-9")); ok {
		t.Fatal("expected symlink-crossing root to be unstageable")
	}
}
