package baseline

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const waitingOnContractLine = "# waiting_on: <who or what it waits on, in one line> # optional, with since"
const sinceContractLine = "# since: YYYY-MM-DD # optional, the date the wait began"

var waitingDeclarationClauses = []struct{ id, sentence string }{
	{
		id:       "clause.context.backlog-01-operational-contract",
		sentence: "An `open` entry that waits on someone or something outside the repository may uncomment `waiting_on`, naming who or what in one line, and `since`, the `YYYY-MM-DD` date the wait began.",
	},
	{
		id:       "clause.context.findings-01-frontmatter",
		sentence: "A `pending` or `partial` Finding that waits on someone or something outside the repository may uncomment `waiting_on`, naming who or what in one line, and `since`, the `YYYY-MM-DD` date the wait began.",
	},
}

func TestTheOperationalContractsShowTheWaitingDeclarationCommented(t *testing.T) {
	t.Parallel()
	clauses := embeddedBaselineClauses(t)
	for _, want := range waitingDeclarationClauses {
		t.Run(want.id, func(t *testing.T) {
			var matches []characterizedBaselineClause
			for _, clause := range clauses {
				if clause.ID == want.id {
					matches = append(matches, clause)
				}
			}
			if len(matches) != 1 {
				t.Fatalf("clause occurs %d times, want 1", len(matches))
			}
			clause := matches[0]
			if clause.Enforcement != "mandatory" || len(clause.Replaces) != 0 {
				t.Errorf("clause force/replacements changed: %+v", clause)
			}
			if count := strings.Count(clause.Guidance, want.sentence); count != 1 {
				t.Errorf("sentence occurs %d times, want 1", count)
			}
			frontmatter := firstFrontmatterBlock(t, clause.Guidance)
			if !strings.HasSuffix(frontmatter, waitingOnContractLine+"\n"+sinceContractLine) {
				t.Errorf("frontmatter does not end with both commented keys: %s", frontmatter)
			}
			t.Run("copied YAML", func(t *testing.T) {
				var fields map[string]yaml.Node
				if err := yaml.Unmarshal([]byte(frontmatter), &fields); err != nil {
					t.Fatal(err)
				}
				for _, key := range []string{"waiting_on", "since"} {
					if _, exists := fields[key]; exists {
						t.Errorf("copied frontmatter declares %s", key)
					}
				}
			})
			t.Run("uncommented YAML", func(t *testing.T) {
				uncommented := strings.NewReplacer(waitingOnContractLine, "waiting_on: ERP vendor support answer", sinceContractLine, "since: 2026-10-25").Replace(frontmatter)
				var fields map[string]yaml.Node
				if err := yaml.Unmarshal([]byte(uncommented), &fields); err != nil {
					t.Fatal(err)
				}
				for key, wantValue := range map[string]string{"waiting_on": "ERP vendor support answer", "since": "2026-10-25"} {
					if node := fields[key]; node.Kind != yaml.ScalarNode || node.Value != wantValue {
						t.Errorf("%s = %+v, want scalar %q", key, node, wantValue)
					}
				}
			})
		})
	}
	for name, path := range map[string]string{
		"Standard TypeScript golden": "assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md",
		"repository guide":           "../../docs/agents/docs-layout.md",
	} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			for _, wantLine := range []string{waitingOnContractLine, sinceContractLine} {
				count := 0
				for _, line := range strings.Split(string(data), "\n") {
					if line == wantLine {
						count++
					}
				}
				if count != 2 {
					t.Errorf("commented line %q occurs %d times on its own line, want 2", wantLine, count)
				}
			}
			text := strings.Join(strings.Fields(string(data)), " ")
			for _, want := range waitingDeclarationClauses {
				if count := strings.Count(text, want.sentence); count != 1 {
					t.Errorf("sentence %q occurs %d times, want 1", want.sentence, count)
				}
			}
		})
	}
}

func TestAnAdopterKeepsTheWaitingContractClauses(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	active, artifacts, err := resolveManagedArtifacts(catalog, profile, standardTypeScriptDecisions("make verify"), false)
	if err != nil {
		t.Fatal(err)
	}
	managed := make([]ManifestArtifact, len(artifacts))
	for i, a := range artifacts {
		managed[i] = ManifestArtifact{ID: a.ID, Path: a.Path, Kind: a.Kind, Module: a.Module, Template: a.Template, Version: a.Version, Digest: a.Digest}
	}
	source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	evidence, delta := classifySourceClauseTransition(source, managed, catalog, active)
	for _, want := range waitingDeclarationClauses {
		t.Run(want.id, func(t *testing.T) {
			if got := delta.Dispositions[want.id]; got != ClauseRetained {
				t.Errorf("disposition = %s, want retained", got)
			}
			count := 0
			for _, entry := range evidence {
				if entry.FromClause != want.id {
					continue
				}
				count++
				if entry.Disposition != string(ClauseRetained) || entry.Enforcement != "mandatory" || !reflect.DeepEqual(entry.Targets, []string{want.id}) {
					t.Errorf("retention evidence = %+v, want mandatory retained clause targeting itself", entry)
				}
			}
			if count != 1 {
				t.Errorf("retention evidence occurs %d times, want 1", count)
			}
		})
	}
}

func firstFrontmatterBlock(t *testing.T, guidance string) string {
	t.Helper()
	_, fenced, ok := strings.Cut(guidance, "```markdown\n---\n")
	if !ok {
		t.Fatal("missing fenced frontmatter opening")
	}
	frontmatter, _, ok := strings.Cut(fenced, "\n---\n")
	if !ok {
		t.Fatal("missing fenced frontmatter closing")
	}
	return frontmatter
}
