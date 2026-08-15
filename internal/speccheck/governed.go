package speccheck

import "regexp"

const toolingAuthorityClause = "docs/agents/agent-instructions.md: linter, formatter, typechecker, test-runner, architecture-checker, build-tool, package-manager, code-generator configuration and scripts, ignore files, plugin declarations, and version pins"

type governedPathSetEntry struct {
	kind    string
	clause  string
	pattern *regexp.Regexp
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
	return entry.kind != "" && entry.clause != "" && entry.pattern.MatchString(path)
}
