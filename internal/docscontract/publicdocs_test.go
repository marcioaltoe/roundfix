// Suite: public documentation contracts
// Invariant: README, the user guides, the agent guidance, and the skill
// mirrors state what the shipped binary actually does.
// Boundary IN: real repository markdown and the public packages it documents.
// Boundary OUT: CLI behavior journeys and parser-level example validation,
// which stay in internal/cli.

// The docscontract tag keeps this invalidation domain out of go test ./...;
// make verify-docs runs it at the pull request boundary.
//go:build docscontract

package docscontract

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
	"roundfix/internal/agent"
	"roundfix/internal/baseline"
	"roundfix/internal/cli"
	"roundfix/internal/spec"
)

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
			command: "baseline update",
			args:    []string{"baseline", "update", "--help"},
			snippets: []string{
				"roundfix baseline update [--repo <path>] [--format <text|json>]",
				"--yes | --confirm-plan <digest>",
				"--adopt-suggested",
				"--no-skills",
				"--skills-source-dir <path>",
				"roundfix/baseline-update-result/v1",
				"0  repository already current, or approved managed refresh applied and verified",
				"1  apply, verification, output, rollback, or recovery failure",
				"2  invalid input, incompatible manifest, or unsafe repository",
				"3  adoption, a new decision, confirmation, or retention action is required",
				"130 operation canceled",
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
			if code := cli.Run(test.args, &stdout, &stderr); code != 0 {
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
				"roundfix baseline update --repo . --yes --format json",
				"roundfix/baseline-update-result/v1",
				"Every byte outside managed boundaries remains byte-identical",
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
				"roundfix baseline update --repo . --yes --format json",
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
				"roundfix baseline update --repo . --yes --format json",
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
	if len(document.Decisions) != 15 {
		t.Fatalf("published project Decision Document decisions = %d, want 15", len(document.Decisions))
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
		len(document.Decisions) != 15 {
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

func TestProfilesDocumentationContractMatchesPublicGuidance(t *testing.T) {
	t.Parallel()
	repoRoot := baselineDocumentationRepoRoot()
	readme := mustRead(t, filepath.Join(repoRoot, "README.md"))
	commands := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "commands.md"))
	usage := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "usage.md"))
	configuration := mustRead(t, filepath.Join(repoRoot, "docs", "user-guide", "configuration.md"))
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
			"sonnet",
			"claude-fable-5",
			"2026-08-07",
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
			agent.PinnedCodexAdapterVersion,
			"official lineage proof",
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
		"Required profiles are `general`, `backend`, `frontend`, `qa`, and `review`.",
		"gpt-5.6-sol",
		"gpt-5.5",
		"Fallback Chain",
	} {
		if !strings.Contains(configuration, want) {
			t.Fatalf("profile configuration guidance is missing %q", want)
		}
	}

	for _, doc := range []struct {
		name    string
		content string
	}{
		{name: "README", content: readme},
		{name: "command reference", content: commands},
		{name: "usage", content: usage},
		{name: "configuration", content: configuration},
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

func TestReviewHistoryConfiguration(t *testing.T) {
	t.Parallel()
	configurationPath := filepath.Join(baselineDocumentationRepoRoot(), ".coderabbit.yaml")
	historyRoot := path.Dir(spec.ArchiveDir(spec.ArchiveKindSpec))
	wantComment := "# history root: " + historyRoot + "/ (must match internal/spec.ArchiveDir)"
	if content := mustRead(t, configurationPath); !strings.Contains(content, wantComment) {
		t.Fatalf("%s does not name the resolved history root in comment %q", configurationPath, wantComment)
	}
	if err := validateReviewHistoryConfiguration(readReviewConfiguration(t)); err != nil {
		t.Fatal(err)
	}
}

func TestReviewHistoryConfigurationRejectsReachableRuleSources(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "history root",
			pattern: "**/AGENTS.md",
		},
		{
			name:    "Spec-owned review directory",
			pattern: "docs/specs/**/reviews/**/*.md",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configuration := readReviewConfiguration(t)
			configuration.KnowledgeBase.CodeGuidelines.FilePatterns = append(
				append([]string(nil), configuration.KnowledgeBase.CodeGuidelines.FilePatterns...),
				test.pattern,
			)
			err := validateReviewHistoryConfiguration(configuration)
			if err == nil {
				t.Fatalf("rule-source pattern %q unexpectedly passed", test.pattern)
			}
			if !strings.Contains(err.Error(), test.pattern) {
				t.Fatalf("error %q does not name offending pattern %q", err, test.pattern)
			}
		})
	}
}

func TestReviewHistoryConfigurationRequiresExclusions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		pattern string
		path    string
	}{
		{
			name:    "history root",
			pattern: "!docs/history/**",
			path:    "docs/history/specs/example/_prd.md",
		},
		{
			name:    "Spec-owned review directory",
			pattern: "!docs/specs/**/reviews/**",
			path:    "docs/specs/example/reviews/round-01/issue.md",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configuration := readReviewConfiguration(t)
			configuration.Reviews.PathFilters = reviewPatternsWithout(
				configuration.Reviews.PathFilters,
				test.pattern,
			)
			excluded, matchErr := reviewPathExcluded(configuration.Reviews.PathFilters, test.path)
			if matchErr != nil {
				t.Fatal(matchErr)
			}
			if excluded {
				t.Fatalf("path %q remains excluded without %q", test.path, test.pattern)
			}
			err := validateReviewHistoryConfiguration(configuration)
			if err == nil {
				t.Fatalf("configuration without %q unexpectedly passed", test.pattern)
			}
			if !strings.Contains(err.Error(), test.pattern) || !strings.Contains(err.Error(), test.path) {
				t.Fatalf("error %q does not name required pattern %q and protected path %q", err, test.pattern, test.path)
			}
		})
	}
}

type reviewConfiguration struct {
	Reviews struct {
		PathFilters []string `yaml:"path_filters"`
	} `yaml:"reviews"`
	KnowledgeBase struct {
		CodeGuidelines struct {
			FilePatterns []string `yaml:"filePatterns"`
		} `yaml:"code_guidelines"`
	} `yaml:"knowledge_base"`
}

func readReviewConfiguration(t *testing.T) reviewConfiguration {
	t.Helper()
	configurationPath := filepath.Join(baselineDocumentationRepoRoot(), ".coderabbit.yaml")
	var configuration reviewConfiguration
	if err := yaml.Unmarshal([]byte(mustRead(t, configurationPath)), &configuration); err != nil {
		t.Fatalf("parse %s: %v", configurationPath, err)
	}
	return configuration
}

func validateReviewHistoryConfiguration(configuration reviewConfiguration) error {
	historyRoot := path.Dir(spec.ArchiveDir(spec.ArchiveKindSpec))
	requiredExclusions := []struct {
		pattern string
		path    string
	}{
		{
			pattern: "!" + historyRoot + "/**",
			path:    spec.ArchiveDir(spec.ArchiveKindSpec) + "/example/_prd.md",
		},
		{
			pattern: "!docs/specs/**/reviews/**",
			path:    "docs/specs/example/reviews/round-01/issue.md",
		},
	}
	for _, required := range requiredExclusions {
		if !reviewPatternPresent(configuration.Reviews.PathFilters, required.pattern) {
			return fmt.Errorf("reviews.path_filters must contain %q to exclude %q", required.pattern, required.path)
		}
		excluded, err := reviewPathExcluded(configuration.Reviews.PathFilters, required.path)
		if err != nil {
			return err
		}
		if !excluded {
			return fmt.Errorf("reviews.path_filters do not exclude %q", required.path)
		}
	}

	protectedRoots := []string{
		historyRoot,
		"docs/specs/example/reviews",
	}
	for _, pattern := range configuration.KnowledgeBase.CodeGuidelines.FilePatterns {
		for _, root := range protectedRoots {
			matched, err := reviewGlobCanReach(pattern, root)
			if err != nil {
				return fmt.Errorf("knowledge_base.code_guidelines.filePatterns pattern %q: %w", pattern, err)
			}
			if matched {
				return fmt.Errorf(
					"knowledge_base.code_guidelines.filePatterns pattern %q reaches protected review tree under %q",
					pattern,
					root,
				)
			}
		}
	}
	return nil
}

func reviewPatternPresent(patterns []string, required string) bool {
	for _, pattern := range patterns {
		if pattern == required {
			return true
		}
	}
	return false
}

func reviewPatternsWithout(patterns []string, removed string) []string {
	result := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		if pattern != removed {
			result = append(result, pattern)
		}
	}
	return result
}

func reviewPathExcluded(patterns []string, candidate string) (bool, error) {
	for _, pattern := range patterns {
		if !strings.HasPrefix(pattern, "!") {
			continue
		}
		matched, err := matchReviewGlob(strings.TrimPrefix(pattern, "!"), candidate)
		if err != nil {
			return false, fmt.Errorf("reviews.path_filters pattern %q: %w", pattern, err)
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func reviewGlobCanReach(pattern, protectedRoot string) (bool, error) {
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	rootSegments := strings.Split(strings.Trim(protectedRoot, "/"), "/")
	for _, segment := range patternSegments {
		if segment == "**" {
			continue
		}
		if _, err := path.Match(segment, ""); err != nil {
			return false, err
		}
	}

	var search func(int, int) bool
	search = func(patternIndex, rootIndex int) bool {
		if rootIndex == len(rootSegments) {
			return patternIndex < len(patternSegments)
		}
		if patternIndex == len(patternSegments) {
			return false
		}

		patternSegment := patternSegments[patternIndex]
		if patternSegment == "**" {
			return search(patternIndex+1, rootIndex) || search(patternIndex, rootIndex+1)
		}

		matched, _ := path.Match(patternSegment, rootSegments[rootIndex])
		if !matched {
			return false
		}
		return search(patternIndex+1, rootIndex+1)
	}
	return search(0, 0), nil
}

func matchReviewGlob(pattern, candidate string) (bool, error) {
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	candidateSegments := strings.Split(strings.Trim(candidate, "/"), "/")
	for _, segment := range patternSegments {
		if segment == "**" {
			continue
		}
		if _, err := path.Match(segment, ""); err != nil {
			return false, err
		}
	}

	var match func(int, int) bool
	match = func(patternIndex, candidateIndex int) bool {
		if patternIndex == len(patternSegments) {
			return candidateIndex == len(candidateSegments)
		}
		if patternSegments[patternIndex] == "**" {
			for index := candidateIndex; index <= len(candidateSegments); index++ {
				if match(patternIndex+1, index) {
					return true
				}
			}
			return false
		}
		if candidateIndex == len(candidateSegments) {
			return false
		}
		matched, _ := path.Match(patternSegments[patternIndex], candidateSegments[candidateIndex])
		return matched && match(patternIndex+1, candidateIndex+1)
	}
	return match(0, 0), nil
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

func TestReleasePlanDocumentationContract(t *testing.T) {
	t.Parallel()
	var rootStdout, rootStderr bytes.Buffer
	if code := cli.Run([]string{"--help"}, &rootStdout, &rootStderr); code != 0 {
		t.Fatalf("root help exit = %d, want 0 stderr=%q", code, rootStderr.String())
	}
	if rootStderr.String() != "" {
		t.Fatalf("root help stderr = %q, want empty", rootStderr.String())
	}
	assertReleasePlanDocumentationContains(t, "root help", rootStdout.String(), []string{
		"roundfix release plan [--from <tag>] [--to <revision>] [--format <text|json>]",
		"release    Plan the next release version without mutating repository or release state",
	})

	var commandStdout, commandStderr bytes.Buffer
	if code := cli.Run([]string{"release", "plan", "--help"}, &commandStdout, &commandStderr); code != 0 {
		t.Fatalf("release plan help exit = %d, want 0 stderr=%q", code, commandStderr.String())
	}
	if commandStderr.String() != "" {
		t.Fatalf("release plan help stderr = %q, want empty", commandStderr.String())
	}
	assertReleasePlanDocumentationContains(t, "release plan help", commandStdout.String(), []string{
		"roundfix release plan [--from <tag>] [--to <revision>] [--impact <none|patch|minor|major> --reason <text>] [--format <text|json>]",
		"ready",
		"approval_required",
		"manual_classification_required",
		"no_release",
		"--from",
		"--to",
		"--impact",
		"--reason",
		"--format",
		"creates no Run",
		"never mutates",
	})

	repoRoot := releasePlanDocumentationRepoRoot()
	docs := []struct {
		name     string
		path     string
		snippets []string
	}{
		{
			name: "release runbook",
			path: filepath.Join(repoRoot, "docs", "user-guide", "release-runbook.md"),
			snippets: []string{
				"roundfix release plan",
				"before changelog edits",
				"version-file edits",
				"GitHub Release creation",
				"generic release request authorizes",
				"patch proposal",
				"minor, major, or version-zero breaking",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
				"tag-triggered",
				"publication workflow",
			},
		},
		{
			name: "usage command index",
			path: filepath.Join(repoRoot, "docs", "user-guide", "usage.md"),
			snippets: []string{
				"Release planning before publication",
				"roundfix release plan",
				"generic release",
				"authorizes only a conclusive patch plan",
				"version-zero",
				"breaking plans require explicit human approval",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
				"| `release plan` | Classify committed release changes without mutating release state |",
			},
		},
		{
			name: "root Agent pointer",
			path: filepath.Join(repoRoot, "AGENTS.md"),
			snippets: []string{
				"Repository-owned rules",
				"docs/agents/specific-repository.md",
			},
		},
		{
			name: "canonical release-planning rule",
			path: filepath.Join(repoRoot, "docs", "agents", "agent-instructions.md"),
			snippets: []string{
				"Release work starts with the read-only release plan",
				"before any changelog, version, tag, push, package, asset",
				"conclusive patch plan",
				"minor, major, version-zero breaking",
				"manual",
				"classification outcomes require the maintainer decisions",
			},
		},
		{
			name: "repository-specific rules",
			path: filepath.Join(repoRoot, "docs", "agents", "specific-repository.md"),
			snippets: []string{
				"roundfix release plan",
				"docs/user-guide/release-runbook.md",
			},
		},
		{
			name: "canonical Roundfix skill",
			path: filepath.Join(repoRoot, ".agents", "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Release planning",
				"roundfix release plan",
				"before changelog edits",
				"version-file edits",
				"GitHub Release creation",
				"generic release request authorizes only a conclusive patch plan",
				"minor,",
				"major, or version-zero breaking",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
			},
		},
		{
			name: "embedded Roundfix skill",
			path: filepath.Join(repoRoot, "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Release planning",
				"roundfix release plan",
				"generic release request authorizes only a conclusive patch plan",
				"--impact <none|patch|minor|major> --reason <text>",
			},
		},
	}
	for _, doc := range docs {
		t.Run(doc.name, func(t *testing.T) {
			content := readReleasePlanDocumentation(t, doc.path)
			assertReleasePlanDocumentationContains(t, doc.name, content, doc.snippets)
		})
	}

	canonical := readReleasePlanDocumentation(t, filepath.Join(repoRoot, ".agents", "skills", "roundfix", "SKILL.md"))
	embedded := readReleasePlanDocumentation(t, filepath.Join(repoRoot, "skills", "roundfix", "SKILL.md"))
	if canonical != embedded {
		t.Fatal("embedded Roundfix skill differs from canonical .agents/skills/roundfix/SKILL.md; run make skills-sync")
	}
}

func releasePlanDocumentationRepoRoot() string {
	return filepath.Clean(filepath.Join("..", ".."))
}

func readReleasePlanDocumentation(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertReleasePlanDocumentationContains(t *testing.T, label string, content string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("%s missing %q", label, snippet)
		}
	}
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

func mustRead(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}
