package baseline

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var secondbrainSearchPhrases = []struct{ name, text string }{
	{"search_with_the_tool", "search with the `pesquisar_acervos` tool of the `indexus` MCP server"},
	{"one_group_per_reached_acervo", "one group per acervo the session's key reaches"},
	{"abrir_conteudo_opens_a_page", "`abrir_conteudo` opens a returned page by its `caminho` as returned"},
	{"wiki_path_relative_to_wiki", "for a `wiki` result is relative to `wiki/`"},
	{"no_wiki_result_is_a_limitation", "An answer without a `wiki` result is a limitation to record"},
	{"no_other_search_tool", "the Secondbrain search is unavailable: substitute no other search tool"},
	{"index_first", "Read `wiki/index.md` first"},
	{"mirrors_are_references", "treat mirrors as references, not workspaces"},
}

var secondbrainDuplicateCheckPhrases = []struct{ name, text string }{
	{"duplicate_check_with_the_tool", "Run the advisory duplicate check first with the Indexus tool `pesquisar_acervos`"},
	{"each_returned_page_opens", "verify that each returned page opens"},
	{"extend_existing_knowledge", "so the owner extends existing knowledge instead of duplicating it"},
	{"no_label_or_rank_decides", "a relevance label or a rank alone never decides"},
}

func assertSecondbrainIndexusPhrases(t *testing.T, text string, phrases []struct{ name, text string }) {
	t.Helper()
	text = strings.Join(strings.Fields(text), " ")
	for _, phrase := range phrases {
		t.Run(phrase.name, func(t *testing.T) {
			if count := strings.Count(text, phrase.text); count != 1 {
				t.Errorf("phrase %q occurs %d times, want once", phrase.text, count)
			}
		})
	}
}

func TestTheSecondbrainClausesNameTheIndexus(t *testing.T) {
	t.Parallel()
	clauses := make(map[string]characterizedBaselineClause)
	for _, clause := range embeddedBaselineClauses(t) {
		clauses[clause.ID] = clause
	}
	search := clauses["clause.secondbrain.02-query-order"]
	duplicate := clauses["clause.secondbrain.research-to-the-knowledge-owner"]
	t.Run("keeps_identity_and_force", func(t *testing.T) {
		for _, id := range []string{"clause.secondbrain.02-query-order", "clause.secondbrain.research-to-the-knowledge-owner"} {
			clause, ok := clauses[id]
			if !ok || clause.ID != id || clause.Enforcement != "mandatory" {
				t.Errorf("clause %s = %+v, want its original identity and mandatory force", id, clause)
			}
		}
		if len(search.Replaces) != 0 {
			t.Errorf("search replacements = %v, want none", search.Replaces)
		}
		if !reflect.DeepEqual(duplicate.Replaces, []string{"clause.secondbrain.research-capture"}) {
			t.Errorf("duplicate-check replacements = %v", duplicate.Replaces)
		}
	})
	t.Run("search_clause", func(t *testing.T) {
		assertSecondbrainIndexusPhrases(t, search.Guidance, secondbrainSearchPhrases)
	})
	t.Run("duplicate_check_clause", func(t *testing.T) {
		assertSecondbrainIndexusPhrases(t, duplicate.Guidance, secondbrainDuplicateCheckPhrases)
	})
	t.Run("no_retired_word", func(t *testing.T) {
		retired := regexp.MustCompile(`(?i)\bqmd\b`)
		for _, clause := range []characterizedBaselineClause{search, duplicate} {
			if retired.MatchString(clause.Guidance) {
				t.Errorf("%s still names the retired search command", clause.ID)
			}
		}
	})
}

func TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	decisions := standardTypeScriptDecisions("make verify")
	for i := range decisions {
		if decisions[i].ID == "secondbrain.enabled" {
			decisions[i].Value = true
		}
	}
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	golden, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md")
	if !ok {
		t.Fatal("catalog has no Secondbrain formatter golden")
	}
	guide, err := os.ReadFile("../../docs/agents/secondbrain.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []struct{ name, content string }{
		{"planned_guide", string(planPostimage(t, plan, "docs/agents/secondbrain.md").Content)},
		{"formatter_golden", string(golden.Data)},
		{"repository_guide", string(guide)},
	} {
		t.Run(text.name, func(t *testing.T) {
			assertSecondbrainIndexusPhrases(t, text.content, secondbrainSearchPhrases)
			assertSecondbrainIndexusPhrases(t, text.content, secondbrainDuplicateCheckPhrases)
		})
	}
}

func TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	earlier := cloneCatalogForRetentionDrift(t, catalog)
	// Delivery Base wording: keep this historical input independent of live Specs.
	for _, clause := range []struct{ id, guidance string }{
		{"clause.secondbrain.02-query-order", "Read `wiki/index.md` first. Then run `qmd query \"<question>\" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces."},
		{"clause.secondbrain.research-to-the-knowledge-owner", "Send substantive external research and relevant findings or articles from Exa MCP as a sourced digest in one message to the project that owns the Secondbrain, for its own ingestion and later reuse as reference material, including relevant articles found outside a dedicated research session. Include each source's title and URL, a concise summary of the content read, the projects or topics it informs, and why it matters; distinguish published evidence from inference. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and name a strong verified match in the message so the owner extends existing knowledge instead of duplicating it; a score alone never decides."},
	} {
		mutateCatalogClause(t, earlier, "secondbrain", clause.id, func(doc document) bool {
			doc["guidance"] = clause.guidance
			return true
		})
	}
	repository := newProjectDecisionPlanRepository(t)
	decisions := standardTypeScriptDecisions("make verify")
	for i := range decisions {
		if decisions[i].ID == "secondbrain.enabled" {
			decisions[i].Value = true
		}
	}
	request := PlanRequest{
		Repository:   repository,
		ProfileID:    "standard-typescript-monorepo",
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	}
	adoption, err := buildPlanWithCatalog(t.Context(), request, earlier)
	if err != nil {
		t.Fatal(err)
	}
	if adoption.Plan == nil {
		t.Fatalf("earlier adoption has no plan: %+v", adoption.Result)
	}
	if _, err := applyPlanWithCatalog(t.Context(), repository, *adoption.Plan, adoption.Plan.PlanDigest, earlier); err != nil {
		t.Fatal(err)
	}
	commitInspectionRepository(t, repository, "adopt earlier Secondbrain wording")
	request.Preservation = RootPreservationRequest{Mode: PreservationModeManagedRefresh}
	refreshed, err := buildPlanWithCatalog(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Plan == nil {
		t.Fatalf("refresh has no plan: %+v", refreshed.Result)
	}
	plan := *refreshed.Plan
	var paths []string
	for _, change := range plan.FileChanges {
		paths = append(paths, change.Path)
	}
	sort.Strings(paths)
	if !reflect.DeepEqual(paths, []string{"docs/agents/secondbrain.md", "docs/agents/setup-context.json"}) {
		t.Errorf("refresh changes %v, want only the guide and the manifest", paths)
	}
	if plan.ClauseDelta == nil {
		t.Fatal("refresh has no clause accounting")
	}
	for id, disposition := range plan.ClauseDelta.Dispositions {
		if disposition == ClauseUnaccounted {
			t.Errorf("unaccounted clause %s", id)
		}
	}
	for _, want := range []struct {
		id          string
		disposition ClauseDisposition
		targets     []string
	}{
		{"clause.secondbrain.02-query-order", ClauseRetained, nil},
		{"clause.secondbrain.research-capture", ClauseReplaced, []string{"clause.secondbrain.research-to-the-knowledge-owner"}},
	} {
		if got := plan.ClauseDelta.Dispositions[want.id]; got != want.disposition {
			t.Errorf("%s disposition = %s, want %s", want.id, got, want.disposition)
		}
		count := 0
		for _, entry := range plan.Retention {
			if entry.FromClause != want.id {
				continue
			}
			count++
			if entry.Enforcement != "mandatory" || entry.Disposition != string(want.disposition) {
				t.Errorf("%s retention = %+v", want.id, entry)
			}
			if want.targets != nil && !reflect.DeepEqual(entry.Targets, want.targets) {
				t.Errorf("%s targets = %v, want %v", want.id, entry.Targets, want.targets)
			}
		}
		if count != 1 {
			t.Errorf("%s has %d retention entries, want one", want.id, count)
		}
	}
	if _, err := applyPlanWithCatalog(t.Context(), repository, plan, plan.PlanDigest, catalog); err != nil {
		t.Fatal(err)
	}
	again, err := buildPlanWithCatalog(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if again.Plan == nil {
		t.Fatalf("second refresh has no plan: %+v", again.Result)
	}
	if len(again.Plan.FileChanges) != 0 {
		t.Errorf("second refresh changes files: %+v", again.Plan.FileChanges)
	}
}

func TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand(t *testing.T) {
	t.Parallel()
	retired := regexp.MustCompile(`(?i)\bqmd\b`)
	// The Source Baseline is an immutable record of what adopters originally held.
	const exempt = "assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md"
	var modules, goldens int
	exemptSeen := false
	for _, root := range []string{"assets", "../../docs/agents"} {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			path = filepath.ToSlash(path)
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.HasPrefix(path, "assets/modules/") {
				modules++
			}
			if strings.HasPrefix(path, "assets/formatter-fixtures/") && strings.Contains(path, "/golden/") {
				goldens++
			}
			if path == exempt {
				exemptSeen = true
				if !retired.Match(content) {
					t.Errorf("%s no longer needs its historical exemption", path)
				}
			} else if retired.Match(content) {
				t.Errorf("%s names the retired search command", path)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if !exemptSeen || modules == 0 || goldens == 0 {
		t.Fatalf("incomplete scan: historical exemption %t, %d modules, %d formatter goldens", exemptSeen, modules, goldens)
	}
}
