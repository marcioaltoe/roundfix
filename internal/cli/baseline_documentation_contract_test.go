package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

// Suite: Baseline public documentation contract
// Invariant: every public Baseline recipe matches the shipped parser and strict document schemas.
// Boundary IN: CLI help, public guides, README, and canonical public skill recipes.
// Boundary OUT: real repository adoption journeys and Python-runtime removal.

func TestBaselineDocumentationContract(t *testing.T) {
	t.Parallel()
	helpCases := []struct {
		command  string
		args     []string
		snippets []string
	}{
		{
			command: "baseline",
			args:    []string{"baseline", "--help"},
			snippets: []string{
				"roundfix baseline [--repo <path>] [--format <text|json>]",
				"roundfix baseline plan (--profile <id> | --profile-file <draft.json>)",
				"roundfix baseline apply --plan <file> --confirm-plan <digest>",
				"roundfix baseline capabilities check [--profile <id>]",
				"roundfix baseline profile init --id <id>",
				"roundfix baseline skills restore --profile <id>",
				"roundfix baseline assets sync --source-dir <path>",
				"explicit confirmation of the displayed Plan Digest",
			},
		},
		{
			command: "baseline capabilities check",
			args:    []string{"baseline", "capabilities", "check", "--help"},
			snippets: []string{
				"roundfix/baseline-capability-recheck/v1",
				"same evaluator",
				"divergence renderer",
				"accepts and resolves",
				"no decisions",
				"never writes repository or journal bytes",
				"no resolvable Baseline Profile",
				"3  capability evidence evaluated with a blocking divergence",
			},
		},
		{
			command: "baseline plan",
			args:    []string{"baseline", "plan", "--help"},
			snippets: []string{
				"roundfix/baseline-plan/v1",
				"roundfix/baseline-result/v1",
				"--profile-file",
				"mutually exclusive",
				"--decision-file",
				"never prompts",
				"never",
				"uses the network",
				"0  complete Baseline Plan emitted",
				"2  invalid arguments",
				"3  a decision",
			},
		},
		{
			command: "baseline apply",
			args:    []string{"baseline", "apply", "--help"},
			snippets: []string{
				"--confirm-plan",
				"recoverable",
				"transaction",
				"Repository formatter and Verification commands are reported as",
				"never run",
				"0  approved plan applied or already applied",
				"1  apply, verification, output, rollback, or recovery failure",
				"3  confirmation mismatch, stale preimage, or unrelated Git lineage",
			},
		},
		{
			command: "baseline profile",
			args:    []string{"baseline", "profile", "--help"},
			snippets: []string{
				"profile init",
				"profile show",
				"profile validate",
			},
		},
		{
			command: "baseline skills restore",
			args:    []string{"baseline", "skills", "restore", "--help"},
			snippets: []string{
				"complete preimage",
				"--source-dir",
				"--confirm-plan",
				"A non-empty preview exits 3",
				"An empty restoration is an idempotent exit 0",
			},
		},
		{
			command: "baseline assets sync",
			args:    []string{"baseline", "assets", "sync", "--help"},
			snippets: []string{
				"--source-dir",
				"--check",
				"never installs skills",
			},
		},
	}
	for _, test := range helpCases {
		t.Run(test.command, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(test.args, &stdout, &stderr); code != exitOK {
				t.Fatalf("help exit = %d, want 0 stderr=%q", code, stderr.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("help stderr = %q, want empty", stderr.String())
			}
			assertBaselineDocumentationContains(t, test.command+" help", stdout.String(), test.snippets)
		})
	}

	root := baselineDocumentationRepoRoot()
	docs := []struct {
		name     string
		path     string
		snippets []string
	}{
		{
			name: "Context-Driven user guide",
			path: filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
			snippets: []string{
				"### First adoption",
				"Greenfield",
				"Preservation",
				"### Update, profile change, and rejected plans",
				"### Automation",
				"roundfix/baseline-plan/v1",
				"roundfix/baseline-result/v1",
				"### Profiles",
				"### Repository Skill Set restoration",
				"### Canonical asset synchronization",
				"### Recovery and troubleshooting",
				"Interrupted transaction",
				"Incomplete rollback",
				"### Security and execution limits",
				"### Migrate from the script-backed setup skill",
				"profile expectation",
				"repository command",
				"recommendation",
			},
		},
		{
			name: "command reference",
			path: filepath.Join(root, "docs", "user-guide", "commands.md"),
			snippets: []string{
				"### baseline",
				"roundfix baseline plan --profile <id>",
				"roundfix baseline apply --plan <file> --confirm-plan <digest>",
				"stdout",
				"stderr",
				"incomplete-rollback",
			},
		},
		{
			name: "README",
			path: filepath.Join(root, "README.md"),
			snippets: []string{
				"roundfix baseline",
				"complete Baseline adoption",
				"migration",
				"recovery",
			},
		},
		{
			name: "canonical setup skill",
			path: filepath.Join(root, ".agents", "skills", "setup-context-driven", "SKILL.md"),
			snippets: []string{
				"public `roundfix baseline` command family",
				"only runtime authority",
				"roundfix baseline plan",
				"roundfix baseline apply",
				"## Recovery",
			},
		},
		{
			name: "canonical Roundfix skill",
			path: filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Context-Driven Baseline",
				"roundfix baseline plan",
				"roundfix baseline apply",
				"no partial plan",
			},
		},
	}
	for _, doc := range docs {
		t.Run(doc.name, func(t *testing.T) {
			content := readBaselineDocumentation(t, doc.path)
			assertBaselineDocumentationContains(t, doc.name, content, doc.snippets)
		})
	}

	guide := readBaselineDocumentation(
		t,
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
	)
	baselineSection := betweenBaselineDocumentation(
		t,
		guide,
		"## Adopt or update the Context-Driven Baseline",
		"## How Roundfix executes it",
	)
	for _, forbidden := range []string{"/Users/", "/home/", `C:\`, "context_setup.py"} {
		if strings.Contains(baselineSection, forbidden) {
			t.Fatalf("Baseline user guide contains environment-specific or retired runtime text %q", forbidden)
		}
	}
}

func TestGuidanceCompositionDocumentation(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	guide := readBaselineDocumentation(
		t,
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
	)
	assertBaselineDocumentationContains(t, "Context-Driven user guide", guide, []string{
		"### Instruction hierarchy",
		"Universal instructions",
		"Context and documentation",
		"Spec workflow",
		"Autonomous work",
		"Stack guidance",
		"Surface guidance",
		"Optional knowledge sources",
		"A narrower guide may add constraints",
		"### Greenfield composition and update redistribution",
		"exact source bytes",
		"semantic owner",
		"docs/agents/specific-repository.md",
		"### ADR and Findings lifecycle",
		"Only `accepted` is active",
		"`pending`, `partial`, `deferred`, and `done`",
		"### Profile alignment and adaptation",
		"Change Baseline Profile",
		"repository-owned Profile adaptation",
		"Decline without writing",
		"--profile-file",
		"mutually exclusive",
		"roundfix baseline skills restore --repo . --profile",
		"--confirm-plan <digest>",
		"Generate a fresh plan",
	})

	generated := readBaselineDocumentation(
		t,
		filepath.Join(
			root,
			"internal",
			"baseline",
			"assets",
			"formatter-fixtures",
			"standard-typescript-monorepo",
			"golden",
			"docs",
			"agents",
			"docs-layout.md",
		),
	)
	tests := []struct {
		name           string
		guideStart     string
		guideEnd       string
		generatedAfter string
	}{
		{
			name:           "ADR lifecycle template",
			guideStart:     "<!-- baseline-adr-lifecycle-template:start -->",
			guideEnd:       "<!-- baseline-adr-lifecycle-template:end -->",
			generatedAfter: "When creating a new ADR",
		},
		{
			name:           "Findings template",
			guideStart:     "<!-- baseline-findings-template:start -->",
			guideEnd:       "<!-- baseline-findings-template:end -->",
			generatedAfter: "complete copyable Findings Operational Contract",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			published := markdownFenceContent(
				t,
				betweenBaselineDocumentation(t, guide, test.guideStart, test.guideEnd),
			)
			generatedStart := strings.Index(generated, test.generatedAfter)
			if generatedStart < 0 {
				t.Fatalf("generated guidance marker %q is missing", test.generatedAfter)
			}
			want := markdownFenceContent(
				t,
				generated[generatedStart:],
			)
			if published != want {
				t.Fatalf("%s differs from generated guidance\npublished:\n%s\nwant:\n%s", test.name, published, want)
			}
		})
	}
}

func TestProjectConstraintDocumentation(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	guidePath := filepath.Join(
		root,
		"docs",
		"user-guide",
		"context-driven-development.md",
	)
	skillPath := filepath.Join(
		root,
		".agents",
		"skills",
		"setup-context-driven",
		"SKILL.md",
	)
	guide := readBaselineDocumentation(t, guidePath)
	skill := readBaselineDocumentation(t, skillPath)
	guideContract := strings.Join(strings.Fields(guide), " ")
	skillContract := strings.Join(strings.Fields(skill), " ")

	assertBaselineDocumentationContains(t, "Context-Driven user guide", guideContract, []string{
		"UUID version 7 is a visible suggestion",
		"`identifier.strategy`",
		"`auth.provider`",
		"`GET` and `POST` under `/api/auth/*`",
		"Session, OAuth redirect, callback, and related provider protocol",
		"exit `3`",
		"no partial Plan",
		"### Project Constraints",
		"Identifier strategy",
		"Authentication and HTTP",
		"Active ADR obligations",
		"Tooling authority",
		"express maintainer authorization",
		"exact bounded repository-relative files",
		"`docs/agents/agent-instructions.md`",
		"`docs/agents/domain.md`",
		"`docs/agents/backend.md`",
		"`docs/agents/spec-routing.md`",
	})
	assertBaselineDocumentationContains(t, "thin setup skill", skillContract, []string{
		"UUID version 7 is a visible suggestion",
		"`identifier.strategy`",
		"`auth.provider`",
		"`GET` and `POST` under `/api/auth/*`",
		"Project Constraints",
		"express maintainer authorization",
		"exact bounded repository-relative files",
		"does not collect, derive, validate, or render decisions",
	})

	example := betweenBaselineDocumentation(
		t,
		guide,
		"<!-- baseline-decision-document:start -->",
		"<!-- baseline-decision-document:end -->",
	)
	example = strings.TrimSpace(example)
	example = strings.TrimPrefix(example, "```json")
	example = strings.TrimSuffix(example, "```")
	example = strings.TrimSpace(example)
	document, err := baseline.ParseDecisionDocument([]byte(example), guidePath)
	if err != nil {
		t.Fatalf("published project Decision Document does not parse: %v", err)
	}
	if len(document.Decisions) != 14 {
		t.Fatalf("published project Decision Document decisions = %d, want 14", len(document.Decisions))
	}

	byID := make(map[string]any, len(document.Decisions))
	for _, decision := range document.Decisions {
		byID[decision.ID] = decision.Value
	}
	if !reflect.DeepEqual(
		byID["identifier.strategy"],
		map[string]any{"kind": "uuid-v7"},
	) {
		t.Fatalf("published identifier strategy = %#v", byID["identifier.strategy"])
	}
	wantAuth := map[string]any{
		"kind": "better-auth",
		"routeException": map[string]any{
			"scope":   "/api/auth/*",
			"methods": []any{"GET", "POST"},
			"owner":   "Better Auth",
			"reason": "Session, OAuth redirect, callback, and related provider " +
				"protocol routes require provider-owned GET and POST semantics.",
		},
	}
	if !reflect.DeepEqual(byID["auth.provider"], wantAuth) {
		t.Fatalf("published authentication provider = %#v", byID["auth.provider"])
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := baseline.ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions := make([]baseline.DecisionValue, 0, len(document.Decisions))
	for _, decision := range document.Decisions {
		if decision.ID != "preservation.mode" {
			decisions = append(decisions, decision)
		}
	}
	if _, missing, err := baseline.ResolveDecisionInput(profile, decisions, catalog); err != nil {
		t.Fatalf("published project decisions are invalid: %v", err)
	} else if len(missing) != 0 {
		t.Fatalf("published project decisions are incomplete: %v", missing)
	}

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

func TestBaselineExamplesParse(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	paths := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
		filepath.Join(root, ".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md"),
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

func TestBaselineDecisionExamples(t *testing.T) {
	t.Parallel()
	path := filepath.Join(
		baselineDocumentationRepoRoot(),
		"docs",
		"user-guide",
		"context-driven-development.md",
	)
	content := readBaselineDocumentation(t, path)
	example := betweenBaselineDocumentation(
		t,
		content,
		"<!-- baseline-decision-document:start -->",
		"<!-- baseline-decision-document:end -->",
	)
	example = strings.TrimSpace(example)
	example = strings.TrimPrefix(example, "```json")
	example = strings.TrimSuffix(example, "```")
	example = strings.TrimSpace(example)

	document, err := baseline.ParseDecisionDocument([]byte(example), path)
	if err != nil {
		t.Fatalf("published Decision Document does not pass the strict parser: %v", err)
	}
	if document.SchemaVersion != baseline.DecisionDocumentSchemaVersion ||
		document.Version != baseline.DecisionDocumentVersion ||
		len(document.Decisions) != 14 {
		t.Fatalf("published Decision Document parsed unexpectedly: %#v", document)
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := baseline.ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions := make([]baseline.DecisionValue, 0, len(document.Decisions))
	for _, decision := range document.Decisions {
		if decision.ID != "preservation.mode" {
			decisions = append(decisions, decision)
		}
	}
	if _, missing, err := baseline.ResolveDecisionInput(profile, decisions, catalog); err != nil {
		t.Fatalf("published complete Decision Document is invalid for standard-typescript-monorepo: %v", err)
	} else if len(missing) != 0 {
		t.Fatalf("published complete Decision Document is missing required decisions: %v", missing)
	}

	invalid := strings.Replace(example, `"version": "0.0.1",`, `"version": "0.0.1", "unknown": true,`, 1)
	if _, err := baseline.ParseDecisionDocument([]byte(invalid), path+"#invalid"); err == nil {
		t.Fatal("strict parser accepted an unknown field in the published Decision Document shape")
	}
}

func parsePublishedBaselineExample(args []string, workDir string) error {
	switch {
	case len(args) == 0:
		_, err := parseBaselineHumanCommand(nil)
		return err
	case args[0] == "plan":
		_, err := parseBaselinePlanCommand(args[1:])
		return err
	case args[0] == "apply":
		_, err := parseBaselineApplyCommand(args[1:])
		return err
	case len(args) >= 2 && args[0] == "profile" && args[1] == "init":
		var stdout, stderr bytes.Buffer
		if code := runBaselineProfileInitCommand(args[2:], &stdout, &stderr, commandEnvironment{workDir: workDir}); code != exitOK {
			return fmt.Errorf("exit %d: %s", code, stderr.String())
		}
		return nil
	case len(args) >= 2 && args[0] == "profile" && args[1] == "show":
		_, _, err := parseBaselineProfileTargetFormat(args[2:], true)
		return err
	case len(args) >= 2 && args[0] == "profile" && args[1] == "validate":
		_, _, err := parseBaselineProfileTargetFormat(args[2:], false)
		return err
	case len(args) >= 2 && args[0] == "skills" && args[1] == "restore":
		_, err := parseBaselineSkillsRestoreCommand(args[2:])
		return err
	case len(args) >= 2 && args[0] == "assets" && args[1] == "sync":
		_, err := parseBaselineAssetsSyncCommand(args[2:])
		return err
	default:
		_, err := parseBaselineHumanCommand(args)
		return err
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
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertBaselineDocumentationContains(t *testing.T, label string, content string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("%s missing %q", label, snippet)
		}
	}
}

func betweenBaselineDocumentation(t *testing.T, content, start, end string) string {
	t.Helper()
	startIndex := strings.Index(content, start)
	if startIndex < 0 {
		t.Fatalf("documentation marker %q is missing", start)
	}
	startIndex += len(start)
	endIndex := strings.Index(content[startIndex:], end)
	if endIndex < 0 {
		t.Fatalf("documentation marker %q is missing", end)
	}
	return content[startIndex : startIndex+endIndex]
}

func markdownFenceContent(t *testing.T, content string) string {
	t.Helper()
	const opening = "```markdown"
	start := strings.Index(content, opening)
	if start < 0 {
		t.Fatalf("markdown fence %q is missing", opening)
	}
	start += len(opening)
	content = content[start:]
	end := strings.Index(content, "```")
	if end < 0 {
		t.Fatal("closing markdown fence is missing")
	}
	return strings.TrimSpace(content[:end])
}
