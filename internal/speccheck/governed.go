package speccheck

import "regexp"

const (
	toolingAuthorityClause           = "docs/agents/agent-instructions.md: linter, formatter, typechecker, test-runner, architecture-checker, build-tool, package-manager, code-generator configuration and scripts, ignore files, plugin declarations, and version pins"
	historicalToolingAuthorityClause = "docs/adr/0130-the-audit-judges-governed-paths-and-history-keeps-the-set-honest.md: every path an authorization has bounded remains governed"
)

type governedPathSetEntry struct {
	kind    string
	clause  string
	pattern *regexp.Regexp
	paths   map[string]struct{}
}

// governedPathSet is compiled into the checker so changing the governed class
// requires the same code review as changing GovernedPath. Each entry carries
// the universal tooling-authority clause that puts its kind in the set.
var governedPathSet = []governedPathSetEntry{
	{
		kind:    "linter configuration and scripts",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)(\.golangci\.(yml|yaml|toml|json)|\.eslintrc(\.[^/]*)?|eslint\.config\.(js|mjs|cjs|ts)|ruff\.toml)$|^\.coderabbit\.ya?ml$`),
	},
	{
		kind:    "formatter configuration and scripts",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)(\.editorconfig|\.prettierrc(\.[^/]*)?|prettier\.config\.(js|mjs|cjs|ts)|rustfmt\.toml)$`),
	},
	{
		kind:    "typechecker configuration and scripts",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)(tsconfig[^/]*\.json|pyrightconfig\.json|mypy\.ini)$`),
	},
	{
		kind:    "test-runner configuration and scripts",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)(pytest\.ini|(jest|vitest)\.config\.(js|mjs|cjs|ts)|playwright\.config\.(js|ts))$`),
	},
	{
		kind:    "architecture-checker configuration and scripts",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)(\.dependency-cruiser\.(js|cjs|json)|dependency-cruiser\.config\.(js|cjs|ts)|archunit\.properties)$`),
	},
	{
		kind:   "build-tool configuration and scripts",
		clause: toolingAuthorityClause,
		pattern: regexp.MustCompile(
			`(^|/)(GNUmakefile|[Mm]akefile|Taskfile\.(yml|yaml)|[Jj]ustfile)$|^\.github/workflows/|^\.roundfixrc\.ya?ml$`,
		),
	},
	{
		kind:   "package-manager configuration and scripts",
		clause: toolingAuthorityClause,
		pattern: regexp.MustCompile(
			`(^|/)(go\.(mod|sum)|package(-lock)?\.json|pnpm-lock\.yaml|yarn\.lock|bun\.lockb?|Cargo\.(toml|lock)|pyproject\.toml|poetry\.lock|uv\.lock)$`,
		),
	},
	{
		kind:   "code-generator configuration and scripts",
		clause: toolingAuthorityClause,
		pattern: regexp.MustCompile(
			`(^|/)(buf\.gen\.yaml|sqlc\.ya?ml|gqlgen\.ya?ml|oapi-codegen\.ya?ml|scripts/(generate|codegen)(\.[^/]*)?)$|^internal/baseline/assets/`,
		),
	},
	{
		kind:    "ignore file",
		clause:  toolingAuthorityClause,
		pattern: regexp.MustCompile(`(^|/)\.[^/]*ignore$`),
	},
	{
		kind:   "plugin declaration",
		clause: toolingAuthorityClause,
		pattern: regexp.MustCompile(
			`(^|/)(\.codex-plugin/plugin\.json|plugin\.json|plugins\.ya?ml)$|^\.agents/skills/|^skills/[^/]+/(SKILL\.md|agents/.*)$`,
		),
	},
	{
		kind:   "version pin",
		clause: toolingAuthorityClause,
		pattern: regexp.MustCompile(
			`(^|/)(\.tool-versions|\.(go|node|python|ruby|java)-version|\.nvmrc|rust-toolchain(\.toml)?)$`,
		),
	},
	{
		kind:   "historically bounded path",
		clause: historicalToolingAuthorityClause,
		paths: exactGovernedPaths(
			"docs/agents/agent-instructions.md",
			"docs/agents/autonomous-work.md",
			"docs/agents/docs-layout.md",
			"docs/agents/secondbrain.md",
			"docs/agents/setup-context.json",
			"docs/agents/skill-dispatch.md",
			"docs/agents/spec-routing.md",
			"docs/agents/specific-repository.md",
			"docs/backlog/2026-08-03-verification-performance-contract.md",
			"docs/backlog/2026-08-06-event-journal-payload-economics.md",
			"docs/backlog/2026-08-06-two-stage-qa-gate-economics.md",
			"docs/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md",
			"docs/findings/2026-08-06-a-promoted-backlog-entry-has-nowhere-valid-to-go.md",
			"docs/references/coverage-record.json",
			"docs/specs/0080-cheap-detectors-run-before-the-gate/references/2026-08-03-verification-performance-contract.md",
			"docs/specs/0080-cheap-detectors-run-before-the-gate/references/2026-08-06-two-stage-qa-gate-economics.md",
			"docs/specs/0080-cheap-detectors-run-before-the-gate/references/_index.md",
			"docs/specs/0081-a-journal-cheap-to-write-and-keep/references/2026-08-06-event-journal-payload-economics.md",
			"docs/specs/0081-a-journal-cheap-to-write-and-keep/references/_index.md",
			"docs/workflow/authorizations/2026-08-06-promoted-backlog-destination.md",
			"internal/baseline/derived_ownership_test.go",
			"internal/baseline/derived_regeneration_repocontract_test.go",
			"internal/baseline/plan_test.go",
			"internal/baseline/repository_test.go",
			"internal/cli/baseline_documentation_contract_test.go",
			"internal/cli/baseline_human_test.go",
			"internal/cli/baseline_plan_test.go",
			"internal/cli/baseline_release_gate_test.go",
			"internal/cli/cli_test.go",
			"internal/cli/releaseplan_documentation_contract_test.go",
			"internal/docscontract",
			"internal/docscontract/doc.go",
			"internal/docscontract/publicdocs_test.go",
			"internal/docscontract/testdata/corpus-golden.json",
			"internal/gittest/gittest.go",
			"internal/spec/archive.go",
			"internal/spec/archive_layout_characterization_test.go",
			"internal/spec/archive_test.go",
			"internal/spec/coverage_test.go",
			"internal/speccheck/backlog.go",
			"internal/speccheck/backlog_test.go",
			"internal/speccheck/coherence.go",
			"internal/speccheck/constraints.go",
			"internal/speccheck/constraints_characterization_test.go",
			"internal/speccheck/testdata/corpus-golden.json",
			"skills/baseline_skill_contract_integration_test.go",
			"skills/baseline_skill_contract_test.go",
			"skills/owned_skill_edit_repocontract_test.go",
		),
	},
}

// GovernedPath reports whether the tooling rules bind path. An ordinary source
// file is not governed and is never audited against a grant. ADR-0130 explains
// why the declared set is held to the historical record.
func GovernedPath(path string) bool {
	clean := cleanMechanicalPath(path)
	if clean == "" {
		return false
	}

	for _, entry := range governedPathSet {
		if entry.matches(clean) {
			return true
		}
	}
	return false
}

func (entry governedPathSetEntry) matches(path string) bool {
	if entry.kind == "" || entry.clause == "" {
		return false
	}
	if entry.pattern != nil && entry.pattern.MatchString(path) {
		return true
	}
	_, ok := entry.paths[path]
	return ok
}

func exactGovernedPaths(paths ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		set[path] = struct{}{}
	}
	return set
}
