package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestProfilesShowStatesTheRecommendationStatus(t *testing.T) {
	t.Parallel()
	writeProfilesCheckFixture(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "show"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	statuses := map[roundconfig.WorkCategory]string{
		roundconfig.CategoryGeneral: "current", roundconfig.CategoryBackend: "differs", roundconfig.CategoryFrontend: "current",
		roundconfig.CategoryData: "inherited", roundconfig.CategoryInfra: "inherited", roundconfig.CategoryDocs: "differs", roundconfig.CategoryTest: "inherited", roundconfig.CategoryChore: "inherited", roundconfig.CategoryQA: "differs", roundconfig.CategoryReview: "pinned",
	}
	for _, category := range roundconfig.AllWorkCategories() {
		start := strings.Index(stdout.String(), "Category: "+string(category)+"\n")
		if start < 0 {
			t.Fatalf("missing category %s", category)
		}
		block := strings.SplitN(stdout.String()[start:], "\n\n", 2)[0]
		statusLine := "Recommendation status: " + statuses[category] + "\n"
		if !strings.Contains(block, statusLine) || strings.Index(block, statusLine) > strings.Index(block, "Recommended profile") {
			t.Errorf("status missing or misplaced: %s", block)
		}
		var deviation string
		if category == roundconfig.CategoryQA {
			deviation = "Deviation: from 2000-01-01 — Older choice"
		}
		if category == roundconfig.CategoryReview {
			deviation = fmt.Sprintf("Deviation: from %s — Keep the validated model", roundconfig.ModelRecommendationSnapshotVersion)
		}
		if deviation != "" && (!strings.Contains(block, deviation) || strings.Index(block, deviation) < strings.Index(block, "Recommended profile")) {
			t.Errorf("deviation missing or misplaced: %s", block)
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, []string{"profiles", "show", "--json"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	var response struct {
		Schema   string `json:"schema"`
		Profiles []struct {
			Category  roundconfig.WorkCategory      `json:"category"`
			Status    string                        `json:"recommendation_status"`
			Deviation *roundconfig.ProfileDeviation `json:"deviation"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != profilesShowSchema || len(response.Profiles) != len(statuses) {
		t.Fatalf("response = %+v", response)
	}
	for _, row := range response.Profiles {
		if row.Status != statuses[row.Category] {
			t.Errorf("%s status = %q", row.Category, row.Status)
		}
		switch row.Category {
		case roundconfig.CategoryReview:
			if row.Deviation == nil || row.Deviation.From != roundconfig.ModelRecommendationSnapshotVersion || row.Deviation.Reason != "Keep the validated model" {
				t.Errorf("pin deviation = %+v", row.Deviation)
			}
		case roundconfig.CategoryQA:
			if row.Deviation == nil || row.Deviation.From != "2000-01-01" || row.Deviation.Reason != "Older choice" {
				t.Errorf("expired deviation = %+v", row.Deviation)
			}
		default:
			if row.Deviation != nil {
				t.Errorf("unexpected deviation: %+v", row)
			}
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, []string{"profiles", "show", "--category", "review", "--json"}, &stdout, &stderr); code != exitOK {
		t.Fatalf("filtered exit=%d stderr=%s", code, &stderr)
	}
	response.Profiles = nil
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Profiles) != 1 || response.Profiles[0].Status != "pinned" {
		t.Fatalf("filter = %+v", response)
	}
}
