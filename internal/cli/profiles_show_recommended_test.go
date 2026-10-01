package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	config "roundfix/internal/config"
)

func TestProfilesShowPrintsTheRecommendedProfile(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "show"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit %d: %s", code, &stderr)
	}
	for _, category := range config.AllWorkCategories() {
		rows, _ := config.ModelRecommendations(category)
		start := strings.Index(stdout.String(), "Category: "+string(category)+"\n")
		if start < 0 {
			t.Fatalf("missing category %s", category)
		}
		block := strings.SplitN(stdout.String()[start:], "\n\n", 2)[0]
		if !strings.Contains(block, "Recommended profile (snapshot "+config.ModelRecommendationSnapshotVersion+"):") {
			t.Fatalf("missing snapshot: %s", block)
		}
		var configure bytes.Buffer
		printProfilesConfigureRecommendations(&configure, category)
		if !strings.Contains(configure.String(), "advisory only") {
			t.Fatal("configure lost advisory label")
		}
		for _, row := range rows {
			want := fmt.Sprintf("  %d. %s %s\n     rationale: %s\n", row.Rank, row.Role, formatProfileSelection(row.Selection), row.Rationale)
			for _, output := range []string{block + "\n", configure.String()} {
				if !strings.Contains(output, want) {
					t.Fatalf("%s missing row %q in %s", category, want, output)
				}
			}
		}
	}
	for _, removed := range []string{"Recommendation source:", "Recommendations snapshot:", "average cost", "category_specific"} {
		if strings.Contains(stdout.String(), removed) {
			t.Fatalf("removed text %q", removed)
		}
	}
}

func TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "show", "--category", "docs", "--json"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit %d: %s", code, &stderr)
	}
	var response struct {
		Schema   string                       `json:"schema"`
		Profiles []map[string]json.RawMessage `json:"profiles"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "roundfix/profiles/v2" || len(response.Profiles) != 1 {
		t.Fatalf("response %+v", response)
	}
	profile := response.Profiles[0]
	if len(profile) != 7 {
		t.Fatalf("profile fields: %v", profile)
	}
	for _, field := range []string{"category", "source", "inherited_from", "preferred", "fallbacks", "recommendations", "recommendation_status"} {
		if _, ok := profile[field]; !ok {
			t.Fatalf("missing %s", field)
		}
	}
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(profile["recommendations"], &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows %v", rows)
	}
	for i, row := range rows {
		if len(row) != 6 {
			t.Fatalf("row fields: %v", row)
		}
		for _, field := range []string{"category", "rank", "role", "selection", "source_as_of", "rationale"} {
			if _, ok := row[field]; !ok {
				t.Fatalf("missing %s", field)
			}
		}
		var role string
		if err := json.Unmarshal(row["role"], &role); err != nil {
			t.Fatal(err)
		}
		want := []string{"preferred", "fallback"}[i]
		if role != want {
			t.Fatalf("role %q want %q", role, want)
		}
	}
	decoded := decodeProfilesShowResponse(t, stdout.String()).Profiles[0]
	wantSelections := []profilesShowTestSelection{
		{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"},
		{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"},
	}
	for i, row := range decoded.Recommendations {
		if row.Category != "docs" || row.Rank != i+1 || row.Selection != wantSelections[i] || row.SourceAsOf != config.ModelRecommendationSnapshotVersion || row.Rationale == "" {
			t.Fatalf("docs recommendation %d = %+v", i, row)
		}
	}
	for _, removed := range []string{"benchmark", "result_percent", "average_cost_usd", "category_specific", "recommendation_source"} {
		if strings.Contains(stdout.String(), fmt.Sprintf("%q", removed)) {
			t.Fatalf("removed key %s", removed)
		}
	}
}
