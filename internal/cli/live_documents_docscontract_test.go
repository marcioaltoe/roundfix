// Suite: live repository skill, guide, help, and release-plan documentation contracts.
// Invariant: shipped documentation and skills agree with CLI parsers, help, and read-only checks.
// Boundary IN: live repository skills, public guides, command help, and adopted repository fixtures.
// Boundary OUT: external services and release execution.
// These tests read the live documentation trees and run under the docscontract tag per ADR-0326.
//go:build docscontract

//verify:always

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/internal/mdtree"
	"roundfix/internal/runevent"
	"roundfix/skills"
)

func TestThisRepositoryHoldsEveryRequiredExternalSkill(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	readiness := checkThisRepositorySkillSet(t, root, true)
	if !readiness.Ready() {
		t.Fatalf("repository skill set is not ready: missing=%v outdated=%v owned=%v", readiness.MissingExternal, readiness.OutdatedExternal, readiness.MissingOwned)
	}
}

func TestARepositoryMissingARequiredExternalSkillIsReported(t *testing.T) {
	t.Parallel()
	root := copyThisRepositorySkillSetFixture(t)
	external, ok, err := resolveExternalSkillRequirement(root)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || len(external) == 0 {
		t.Fatalf("fixture external requirement = %v, manifest readable = %t", external, ok)
	}
	missing := external[0]
	if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", missing)); err != nil {
		t.Fatalf("remove fixture skill %q: %v", missing, err)
	}

	readiness := checkThisRepositorySkillSet(t, root, false)
	if readiness.Ready() {
		t.Fatal("repository with a missing required external skill was reported ready")
	}
	found := false
	for _, name := range readiness.MissingExternal {
		if name == missing {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing external skills = %v, want %q", readiness.MissingExternal, missing)
	}
}

func checkThisRepositorySkillSet(t *testing.T, root string, wantReady bool) skills.RepositoryReadiness {
	t.Helper()
	external, ok, err := resolveExternalSkillRequirement(root)
	if err != nil {
		t.Fatalf("resolve external skill requirement: %v", err)
	}
	if !ok {
		t.Fatal("Setup Manifest did not resolve an external skill requirement")
	}
	readiness, err := skills.CheckRepositoryWithExternal(t.Context(), root, external)
	if err != nil {
		t.Fatalf("check repository skill set: %v", err)
	}
	for _, name := range readiness.MissingExternal {
		t.Logf("missing external skill: %s", name)
	}
	for _, name := range readiness.OutdatedExternal {
		t.Logf("outdated external skill: %s", name)
	}
	if wantReady && !readiness.Ready() {
		t.Errorf("repository skill set is not ready: missing=%v outdated=%v", readiness.MissingExternal, readiness.OutdatedExternal)
	}

	lockBytes, err := os.ReadFile(filepath.Join(root, "skills-lock.json"))
	if err != nil {
		t.Fatalf("read skills lock: %v", err)
	}
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	for _, name := range []string{"context7", "feature-systems-pattern", "rust", "review", "triage", "the-fool", "autoresearch", "council"} {
		if _, exists := lock.Skills[name]; exists {
			t.Errorf("obsolete skill remains in skills-lock.json: %s", name)
		}
		if _, err := os.Stat(filepath.Join(root, ".agents", "skills", name)); err == nil {
			t.Errorf("obsolete skill directory remains: %s", name)
		} else if !os.IsNotExist(err) {
			t.Errorf("inspect obsolete skill directory %s: %v", name, err)
		}
	}
	return readiness
}

func copyThisRepositorySkillSetFixture(t *testing.T) string {
	t.Helper()
	source, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir()
	for _, relative := range []string{
		"docs/agents/setup-context.json",
		"skills-lock.json",
		".agents/skills",
	} {
		copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(destination, relative))
	}
	return destination
}

func TestBaselineExamplesParse(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	paths := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
		filepath.Join(root, ".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(root, ".agents/skills/roundfix/SKILL.md"),
	}

	temporaryRepo := t.TempDir()
	if err := os.Mkdir(filepath.Join(temporaryRepo, ".git"), 0o755); err != nil {
		t.Fatalf("create documentation-test Git marker: %v", err)
	}
	exampleCount := 0
	profileCounter := 0
	for _, path := range paths {
		content := readBaselineDocumentation(t, path)
		for _, command := range baselineBashExamples(content) {
			exampleCount++
			args := strings.Fields(command)
			if len(args) < 2 || args[0] != "roundfix" || args[1] != "baseline" {
				t.Fatalf("%s: invalid Baseline command example %q", path, command)
			}
			if len(args) >= 4 && args[2] == "profile" && args[3] == "init" {
				profileCounter++
				args = normalizeBaselineProfileInitExample(args, profileCounter)
			}
			if err := parsePublishedBaselineExample(args[2:], temporaryRepo); err != nil {
				t.Fatalf("%s: command example %q does not parse: %v", path, command, err)
			}
		}
	}
	if exampleCount < 20 {
		t.Fatalf("parsed %d Baseline examples, want at least 20 public recipes", exampleCount)
	}
}

// TestProjectConstraintPlanningExamplesParse keeps the parser-level half of
// the project-constraint documentation contract beside the parsers it needs:
// the string contracts on the same two documents live in
// internal/docscontract.
func TestProjectConstraintPlanningExamplesParse(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	guide := readBaselineDocumentation(t, filepath.Join(
		root, "docs", "user-guide", "context-driven-development.md",
	))
	skill := readBaselineDocumentation(t, filepath.Join(
		root, ".agents", "skills", "setup-context-driven", "SKILL.md",
	))
	parsedPlanningExamples := 0
	for _, content := range []string{guide, skill} {
		for _, command := range baselineBashExamples(content) {
			if !strings.Contains(command, "baseline plan") ||
				!strings.Contains(command, "--decision-file") {
				continue
			}
			args := strings.Fields(command)
			if err := parsePublishedBaselineExample(args[2:], ""); err != nil {
				t.Fatalf("project-decision command %q does not parse: %v", command, err)
			}
			parsedPlanningExamples++
		}
	}
	if parsedPlanningExamples < 2 {
		t.Fatalf("parsed project-decision planning examples = %d, want at least 2", parsedPlanningExamples)
	}
}

func normalizeBaselineProfileInitExample(args []string, sequence int) []string {
	normalized := append([]string(nil), args...)
	for index := 0; index+1 < len(normalized); index++ {
		switch normalized[index] {
		case "--id":
			normalized[index+1] = fmt.Sprintf("documentation-profile-%d", sequence)
		case "--from":
			if strings.HasPrefix(normalized[index+1], "<") {
				normalized[index+1] = "go-cli-tui"
			}
		}
	}
	return normalized
}

func baselineBashExamples(content string) []string {
	var examples []string
	inBash := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "```bash" {
			inBash = true
			continue
		}
		if trimmed == "```" {
			inBash = false
			continue
		}
		if inBash && (trimmed == "roundfix baseline" || strings.HasPrefix(trimmed, "roundfix baseline ")) {
			examples = append(examples, trimmed)
		}
	}
	return examples
}

func baselineDocumentationRepoRoot() string {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return filepath.Clean(filepath.Join("..", ".."))
	}
	return root
}

func readBaselineDocumentation(t *testing.T, path string) string {
	t.Helper()
	if filepath.Base(path) == "SKILL.md" && filepath.Base(filepath.Dir(path)) == "roundfix" {
		text, err := mdtree.Text(os.DirFS(filepath.Dir(path)), "SKILL.md", "references")
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		return text
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func TestEventsHelpDocumentsAgentSelectionFilter(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"events", "--help"}, &stdout, &stderr)

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

	guide, err := mdtree.Text(os.DirFS(cliTestRepoRoot(t)), "docs/user-guide/commands.md", "docs/user-guide/commands")
	if err != nil {
		t.Fatalf("read command reference: %v", err)
	}
	if !strings.Contains(guide, wantCategories) || !strings.Contains(guide, "`agent-selection`") {
		t.Fatalf("user guide does not document agent-selection filter")
	}
}

func TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault(t *testing.T) {
	t.Parallel()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := catalog.Decision("http.contract")
	if !ok {
		t.Fatal("HTTP decision is absent")
	}
	var decision struct {
		Default struct {
			Mode string `json:"mode"`
		} `json:"default"`
		Modes []string `json:"modes"`
	}
	if err := json.Unmarshal(entry.Data, &decision); err != nil {
		t.Fatal(err)
	}
	if decision.Default.Mode == "" || len(decision.Modes) == 0 {
		t.Fatal("HTTP decision has no default or modes")
	}
	contract, ok := catalog.Asset("contract-v1.json")
	if !ok {
		t.Fatal("Baseline contract is absent")
	}
	guide := readBaselineDocumentation(t, filepath.Join(baselineDocumentationRepoRoot(), "docs/user-guide/context-driven-development.md"))
	statements := []httpDefaultStatement{
		{"contract-v1.json", string(contract.Data), "the interactive workflow suggests %s until"},
		{"context-driven-development.md", guide, "| HTTP contract | `%s` |"},
	}
	if findings := httpDefaultStatementFindings(decision.Default.Mode, decision.Modes, statements); len(findings) != 0 {
		t.Fatalf("HTTP statement findings: %v", findings)
	}
}

func TestReleasePlanChecksOnAnAdoptedRepository(t *testing.T) {
	t.Parallel()
	for _, stale := range []bool{false, true} {
		t.Run(fmt.Sprint(stale), func(t *testing.T) {
			repo := newBaselineUpdateRepository(t)
			if stale {
				repo = staleBaselineUpdateRepository(t)
			}
			copyReleaseCheckSkills(t, repo)
			commitBaselinePlanTestRepository(t, repo)
			addReleaseCheckRange(t, repo)
			code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format=json")
			if code != exitOK || stderr != "" {
				t.Fatalf("exit=%d output=%s stderr=%s", code, output, stderr)
			}
			checks := decodeReleaseChecks(t, output)
			if checks.Skills.Status != "ok" || checks.Skills.NextAction != "" {
				t.Fatalf("skills=%+v", checks.Skills)
			}
			if !stale {
				if checks.Baseline.Status != "current" || checks.Baseline.NextAction != "" {
					t.Fatalf("Surface Transcript 2: %+v", checks)
				}
				if strings.Contains(output, `"nextAction"`) {
					t.Fatal("empty nextAction must be omitted")
				}
			} else {
				result, _, stderr, _ := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repo, "--no-skills", "--format=json")
				want := fmt.Sprintf("%d file change(s)", len(result.FileChanges))
				if len(result.HistoryMoves) > 0 {
					want += fmt.Sprintf(", %d history move(s)", len(result.HistoryMoves))
				}
				if stderr != "" || checks.Baseline.Status != "plan_ready" || checks.Baseline.Detail != want {
					t.Fatalf("baseline=%+v update=%+v stderr=%s", checks.Baseline, result, stderr)
				}
			}
		})
	}
}

func TestReleasePlanChecksWriteNothing(t *testing.T) {
	t.Parallel()
	repo := staleBaselineUpdateRepository(t)
	copyReleaseCheckSkills(t, repo)
	commitBaselinePlanTestRepository(t, repo)
	addReleaseCheckRange(t, repo)
	before := snapshotReleasePlanRepo(t, repo)
	for _, format := range []string{"text", "json"} {
		code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format="+format)
		if code != exitOK || stderr != "" || !strings.Contains(output, "plan_ready") {
			t.Fatalf("exit=%d output=%s stderr=%s", code, output, stderr)
		}
		assertReleasePlanRepoUnchanged(t, repo, before)
	}
}

func addReleaseCheckRange(t *testing.T, repo string) {
	t.Helper()
	gitReleasePlan(t, repo, "tag", "v0.4.0")
	writeReleasePlanFile(t, repo, "internal/output.go", "package output\n")
	gitReleasePlan(t, repo, "add", "-A")
	commitReleasePlan(t, repo, "fix: correct output")
}

func copyReleaseCheckSkills(t *testing.T, repo string) {
	t.Helper()
	lock, err := os.ReadFile(filepath.Join(baselineDocumentationRepoRoot(), "skills-lock.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "skills-lock.json"), lock, 0644); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(baselineDocumentationRepoRoot(), ".agents", "skills")
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(repo, ".agents", "skills", relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
