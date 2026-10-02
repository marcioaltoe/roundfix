// Suite: Go, CLI and TUI clause force
// Invariant: every stack rule carries force and the Go preference bans no named library.
// Boundary IN: embedded modules and rendered plans for temporary adopters.
// Boundary OUT: applying a repository refresh and Daemon Verification.

package baseline

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func ruleLevelGuidanceFindings(moduleID string, module map[string]any) []string {
	var findings []string
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if guidance, ok := rule["guidance"].(string); ok && strings.TrimSpace(guidance) != "" {
			findings = append(findings, fmt.Sprintf("module %s rule %s carries rule-level guidance", moduleID, rule["id"]))
		}
	}
	return findings
}

func namedLibraryBanFindings(guide string, libraries []string) []string {
	var findings []string
	ban := regexp.MustCompile(`(?i)\b(do not|must not|never|ban(?:ned|s)?|forbid(?:den|s)?|prohibit(?:ed|s)?)\b`)
	sentences := strings.FieldsFunc(strings.Join(strings.Fields(guide), " "), func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	for _, library := range libraries {
		name := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(library) + `\b`)
		for _, sentence := range sentences {
			if name.MatchString(sentence) && ban.MatchString(sentence) {
				findings = append(findings, "named library ban: "+library)
				break
			}
		}
	}
	return findings
}

func TestEveryGoCLIAndTUIRuleCarriesForce(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, moduleID := range []string{"go", "cli-surface", "tui-surface"} {
		t.Run(moduleID, func(t *testing.T) {
			asset, ok := catalog.Module(moduleID)
			if !ok {
				t.Fatalf("missing module %s", moduleID)
			}
			var module map[string]any
			if err := json.Unmarshal(asset.Data, &module); err != nil {
				t.Fatal(err)
			}
			if findings := ruleLevelGuidanceFindings(moduleID, module); len(findings) != 0 {
				t.Fatalf("rule-level guidance findings: %v", findings)
			}
		})
	}
}

func TestARuleWithoutForceIsReported(t *testing.T) {
	t.Parallel()
	module := map[string]any{
		"schemaVersion": "setup-context-driven/module-v2",
		"rules":         []any{map[string]any{"id": "rule.go.literal", "guidance": "Use stdlib testing."}},
	}
	want := []string{"module go rule rule.go.literal carries rule-level guidance"}
	if got := ruleLevelGuidanceFindings("go", module); !reflect.DeepEqual(got, want) {
		t.Fatalf("rule-level guidance findings = %v, want %v", got, want)
	}
}

func TestTheGoCLIAndTUIGuidesStateTheirClauses(t *testing.T) {
	t.Parallel()
	plan := buildTestPlan(t, newPlanRepository(t))
	rules := []stackWordingRule{
		{guide: "docs/agents/go.md", must: []string{
			"- **prohibited**: Do not hand-edit `go.mod` or `go.sum`; change them through the `go` command.",
			"- **mandatory**: Keep each Go `main` package thin: it parses input, wires dependencies, and calls behavior that lives in cohesive packages.",
			"- **mandatory**: Prefer the Go standard library. Add a third-party module only for a named job that the standard library and the modules already required cannot do, and record that reason in an ADR or another recorded repository decision in the change that adds it.",
			"- **mandatory**: Give every goroutine an owner that waits for it and a cancellation path that stops it.",
			"- **mandatory**: Pass a `context.Context` as the first parameter of blocking and IO work, and stop that work when the context is cancelled.",
			"- **mandatory**: Wrap a returned error with the operation that failed using `%w`, and keep `errors.Is` and `errors.As` matching wherever a caller branches on the error.",
			"- **mandatory**: When a change touches a file with a build constraint, build the affected non-test packages for every operating system its constraints name, for example with `GOOS=windows go build ./...`; a build or test run on the host compiles only the host's files.",
			"- **mandatory**: Test observable package and command behavior through public entry points: stdout, stderr, files, exit codes, cancellation, and failure paths.",
			"- **mandatory**: Run Go tests through `go test` with the standard `testing` package as the harness. An assertion or mocking library is a third-party module and needs its recorded reason.",
		}, mustNot: []string{"Use stdlib `testing`.", "Prefer the standard library; add a dependency only for a named job it cannot perform"}},
		{guide: "docs/agents/cli.md", must: []string{
			"- **mandatory**: Treat command names, flags, stdout and stderr placement, machine-readable fields, and exit codes as public API.",
			"- **mandatory**: Keep stdout for requested output and stderr for diagnostics, progress, and warnings.",
			"- **mandatory**: When the repository ships a skill or a command reference that describes its commands, a change to command behavior updates that skill or reference in the same pull request.",
			"- **mandatory**: Make automation deterministic and non-interactive.",
			"- **mandatory**: Make write operations explicit, replayable, safe by default, and observable; use dry-run, confirmation, or idempotency contracts where the repository requires them.",
		}},
		{guide: "docs/agents/tui.md", must: []string{
			"- **mandatory**: Drive TUI model updates synchronously and assert rendered state, messages, and transitions.",
			"- **mandatory**: Use terminal emulation only when model-level tests cannot prove the behavior.",
			"- **mandatory**: Keep layout and interaction policy in repository-owned design guidance.",
		}},
	}
	rendered := make(map[string]string)
	for _, rule := range rules {
		text := string(planPostimage(t, plan, rule.guide).Content)
		rendered[rule.guide] = text
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "- ") && !strings.HasPrefix(line, "- **mandatory**: ") && !strings.HasPrefix(line, "- **prohibited**: ") {
				t.Errorf("%s: rule line has no force: %s", rule.guide, line)
			}
		}
	}
	if findings := stackWordingFindings(rendered, rules); len(findings) != 0 {
		t.Fatalf("stack clause wording findings: %v", findings)
	}
}

func TestTheGoGuideBansNoLibraryByName(t *testing.T) {
	t.Parallel()
	plan := buildTestPlan(t, newPlanRepository(t))
	guide := string(planPostimage(t, plan, "docs/agents/go.md").Content)
	if findings := namedLibraryBanFindings(guide, []string{"Cobra", "testify", "Viper"}); len(findings) != 0 {
		t.Fatalf("Go named-library ban findings: %v", findings)
	}
}

func TestALibraryBanInTheGoGuideIsReported(t *testing.T) {
	t.Parallel()
	want := []string{"named library ban: Cobra"}
	if got := namedLibraryBanFindings("Do not use Cobra.", []string{"Cobra", "testify", "Viper"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("named-library ban findings = %v, want %v", got, want)
	}
}

func TestTheGoCLIAndTUIClausesCarryTheirForce(t *testing.T) {
	t.Parallel()
	want := []expectedClause{
		{"go", "clause.go.change-module-files-through-go", "prohibited"},
		{"go", "clause.go.keep-entry-points-thin", "mandatory"},
		{"go", "clause.go.record-the-reason-for-a-module", "mandatory"},
		{"go", "clause.go.own-every-goroutine", "mandatory"},
		{"go", "clause.go.pass-context-first", "mandatory"},
		{"go", "clause.go.wrap-errors-with-the-operation", "mandatory"},
		{"go", "clause.go.build-every-constrained-platform", "mandatory"},
		{"go", "clause.go.test-observable-behavior", "mandatory"},
		{"go", "clause.go.test-through-go-test", "mandatory"},
		{"cli-surface", "clause.cli.public-command-contract", "mandatory"},
		{"cli-surface", "clause.cli.separate-output-streams", "mandatory"},
		{"cli-surface", "clause.cli.ship-the-skill-with-the-behavior", "mandatory"},
		{"cli-surface", "clause.cli.deterministic-non-interactive", "mandatory"},
		{"cli-surface", "clause.cli.explicit-safe-writes", "mandatory"},
		{"tui-surface", "clause.tui.drive-models-synchronously", "mandatory"},
		{"tui-surface", "clause.tui.emulate-the-terminal-last", "mandatory"},
		{"tui-surface", "clause.tui.keep-design-policy-in-repository-guidance", "mandatory"},
	}
	if findings := clauseForceFindings(mustEmbeddedCatalog(t), want); len(findings) != 0 {
		t.Fatalf("Go, CLI and TUI clause force findings: %v", findings)
	}
}
