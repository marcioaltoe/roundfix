package cli

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles(t *testing.T) {
	t.Parallel()
	repository := languageRulePriorRepository(t)
	before := baselinePlanTestTree(t, repository)
	args := []string{"baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format=json"}
	stopped, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), args...)
	if code != 3 || stderr != "" || stopped.Category != "decision" || !strings.Contains(stdout, "language.reader") {
		t.Fatalf("stop: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if before != baselinePlanTestTree(t, repository) {
		t.Fatal("required-decision stop wrote repository files")
	}
	args = append(args, "--decision", "language.reader=English", "--yes")
	result, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), args...)
	if code != exitOK || stderr != "" {
		t.Fatalf("apply: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var changed []string
	for _, change := range result.FileChanges {
		changed = append(changed, change.Path)
	}
	slices.Sort(changed)
	want := []string{"docs/agents/agent-instructions.md", "docs/agents/domain.md", "docs/agents/setup-context.json"}
	if !reflect.DeepEqual(changed, want) {
		t.Fatalf("changes = %v, want %v", changed, want)
	}
	if result.ClauseDelta == nil || result.ClauseDelta.Counts[baseline.ClauseUnaccounted] != 0 || result.ClauseDelta.Dispositions["clause.domain.canonical-language"] != baseline.ClauseRetained {
		t.Fatalf("clause accounting = %+v", result.ClauseDelta)
	}
	for id, disposition := range result.ClauseDelta.Dispositions {
		if disposition == baseline.ClauseUnaccounted {
			t.Errorf("prior clause %s = %s", id, disposition)
		}
	}
	domain, err := os.ReadFile(filepath.Join(repository, "docs/agents/domain.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(strings.Fields(string(domain)), " "), "- **mandatory**: Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. In a code name, a term takes its Name in Code: the term itself, or the name its Glossary Definition records on a line `_In code_:` when the term was translated to English or kept without accents. A Spec that introduces or changes such a term writes that line. Call out a missing term instead of inventing a competing synonym.") {
		t.Fatal("canonical clause lost text or mandatory force")
	}
	current, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format=json")
	if code != exitOK || stderr != "" || current.State != "current" || len(current.FileChanges) != 0 {
		t.Fatalf("second update: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
}

func TestTheUpdateNamesTheReaderLanguageDecision(t *testing.T) {
	t.Parallel()
	repository := languageRulePriorRepository(t)
	before := baselinePlanTestTree(t, repository)
	_, stdout, stderr, code := runBaselineUpdateTestCommand(t, context.Background(), "baseline", "update", "--repo", repository, "--no-skills", "--no-history", "--format", "text")
	for _, label := range []string{"Prior catalog:", "Current catalog:"} {
		digest := baselineUpdateTextValue(stdout, label)
		if len(digest) != len("sha256:")+64 || !strings.HasPrefix(digest, "sha256:") {
			t.Fatalf("invalid %s %q", label, digest)
		}
		stdout = strings.ReplaceAll(stdout, digest, "<digest>")
	}
	const want = `Baseline update: action required
Category: decision
Result: the current Baseline catalog requires new decisions: language.reader
Prior catalog: <digest>
Current catalog: <digest>
File changes: 0
Misplaced repository rules: 0
Retention evidence: 0
Declared decisions: 0
New decision: language.reader suggested="English" (The Reader Language: the language of what the product's users read.)
Skills: not run
Skills installed: 0
Skills restored: 0
Skills drifted: 0
Next action: answer the named decisions through roundfix baseline, or rerun roundfix baseline update with --adopt-suggested
`
	if stdout != want || stderr != "" || code != 3 {
		t.Fatalf("Surface Transcript 1: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if before != baselinePlanTestTree(t, repository) {
		t.Fatal("transcript stop wrote files")
	}
}

func TestARepositoryProfileWithoutTheReaderLanguageIsRefused(t *testing.T) {
	t.Parallel()
	repository := newBaselineApplyTestRepository(t)
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	created, err := baseline.InitCustomProfile(repository, "language-profile", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatal(err)
	}
	var profile map[string]any
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatal(err)
	}
	decisions := profile["decisions"].([]any)
	profile["decisions"] = slices.DeleteFunc(decisions, func(v any) bool { return v == "language.reader" })
	data, err = json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, created.Path, string(data))
	before := baselinePlanTestTree(t, repository)
	commitBaselinePlanTestRepository(t, repository)
	stdout, stderr, code := runRepositoryProfileCommand(t, "baseline", "plan", "--repo", repository, "--profile", "language-profile", "--decision", "preservation.mode=greenfield", "--format=text")
	if code != exitUnverified || stderr != "" {
		t.Fatalf("profile refusal: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	for _, want := range []string{"language.reader", "template.guide.agent-instructions", `add "language.reader" to "decisions"`} {
		if !strings.Contains(stdout, want) {
			t.Errorf("profile refusal %q omits %q", stdout, want)
		}
	}
	if before != baselinePlanTestTree(t, repository) {
		t.Fatal("profile refusal wrote files")
	}
	profile["decisions"] = append(profile["decisions"].([]any), "language.reader")
	data, err = json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, created.Path, string(data))
	commitBaselinePlanTestRepository(t, repository)
	answers := baselineReleaseDecisionArgs("go-cli-tui", "greenfield")
	args := []string{"baseline", "plan", "--repo", repository, "--profile", "language-profile", "--format=json"}
	for _, answer := range answers {
		args = append(args, "--decision", answer)
	}
	stdout, stderr, code = runRepositoryProfileCommand(t, args...)
	if code != exitOK || stderr != "" {
		t.Fatalf("repaired profile: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	plan, err := baseline.ParsePlanDocument([]byte(stdout))
	if err != nil {
		t.Fatal(err)
	}
	for _, image := range plan.Postimages {
		if image.Path == "docs/agents/agent-instructions.md" {
			guide := strings.Join(strings.Fields(string(image.Content)), " ")
			for _, text := range []string{"The Reader Language is English.", "- **mandatory**: Write code and infrastructure names in English: variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs.", "- **mandatory**: A business term with no faithful English translation keeps its original form inside an English name: an official acronym as it is (`icms_amount`), another term without accents (`apuracao_summary`). Brazilian fiscal and legal terms are the usual case.", "- **mandatory**: A new name follows this rule. An existing name changes only when its code is touched, a large rename takes a Spec of its own, and nothing is renamed in bulk without a request. The Pre-PR Review reports a new name that breaks the rule, and the QA gate reports text a user reads outside the Reader Language.", "- **mandatory**: Write what the product's users read in the Reader Language named above: screens, user messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Specs, ADRs, agent guides, commits and pull requests keep the language they already have."} {
				if !strings.Contains(guide, text) {
					t.Errorf("repaired profile omits %s", text)
				}
			}
			return
		}
	}
	t.Fatal("repaired profile rendered no agent guide")
}

// Reconstruct the prior catalog's managed text and decisions on the existing
// adoption fixture, retaining its Source Baseline accounting and region markers.
func languageRulePriorRepository(t *testing.T) string {
	t.Helper()
	repository := newBaselineReleaseRepository(t, "standard-typescript-monorepo")
	args := []string{"baseline", "plan", "--repo", repository, "--profile", "standard-typescript-monorepo", "--format=json"}
	for _, answer := range baselineReleaseDecisionArgs("standard-typescript-monorepo", "greenfield") {
		args = append(args, "--decision", answer)
	}
	stdout, stderr, code := runRepositoryProfileCommand(t, args...)
	if code != exitOK || stderr != "" {
		t.Fatalf("adopt prior-text fixture: exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	plan, err := baseline.ParsePlanDocument([]byte(stdout))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := baseline.ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	manifest := ReadBaselineSetupManifest(t, repository)
	delete(manifest.Decisions, "language.reader")
	manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)
	for _, path := range []string{"docs/agents/agent-instructions.md", "docs/agents/domain.md"} {
		data, err := os.ReadFile(filepath.Join(repository, path))
		if err != nil {
			t.Fatal(err)
		}
		text := string(data)
		if strings.HasSuffix(path, "agent-instructions.md") {
			text = strings.ReplaceAll(text, "The Reader Language is English.\n", "")
			text = strings.ReplaceAll(text, "The Reader Language is `English`.\n", "")
			for _, line := range []string{"- **mandatory**: Write code and infrastructure names in English: variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs.\n\n", "- **mandatory**: A business term with no faithful English translation keeps its original form inside an English name: an official acronym as it is (`icms_amount`), another term without accents (`apuracao_summary`). Brazilian fiscal and legal terms are the usual case.\n\n", "- **mandatory**: A new name follows this rule. An existing name changes only when its code is touched, a large rename takes a Spec of its own, and nothing is renamed in bulk without a request. The Pre-PR Review reports a new name that breaks the rule, and the QA gate reports text a user reads outside the Reader Language.\n\n", "- **mandatory**: Write what the product's users read in the Reader Language named above: screens, user messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Specs, ADRs, agent guides, commits and pull requests keep the language they already have.\n\n"} {
				text = strings.ReplaceAll(text, line, "")
			}
		} else {
			text = strings.ReplaceAll(text, "Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. In a code name, a term takes its Name in Code: the term itself, or the name its Glossary Definition records on a line `_In code_:` when the term was translated to English or kept without accents. A Spec that introduces or changes such a term writes that line. Call out a missing term instead of inventing a competing synonym.", "Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. Call out a missing term instead of inventing a competing synonym.")
		}
		mustWrite(t, filepath.Join(repository, path), text)
		updateFleetManifestArtifactDigest(t, repository, &manifest, path)
	}
	WriteBaselineSetupManifest(t, repository, manifest)
	commitBaselinePlanTestRepository(t, repository)
	return repository
}
