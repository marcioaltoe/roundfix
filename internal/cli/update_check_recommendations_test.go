// Suite: Update Check profile recommendation reports
// Invariant: only differing profiles add advice; streams, next, exit, and configuration stay unchanged.
// Boundary IN: real command dispatch, committed Git fixtures, config loading, and Baseline preview.
// Boundary OUT: release transport, replaced by the existing command-context lister.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

const updateCheckOldDocsProfile = `  docs:
    preferred: {runtime: codex, model: gpt-5.6-luna, reasoning_effort: max}
    fallbacks:
      - {runtime: claude, model: sonnet, reasoning_effort: high}
`
const updateCheckOldReviewProfile = `  review:
    preferred: {runtime: codex, model: gpt-5.6-luna, reasoning_effort: max}
    fallbacks:
      - {runtime: codex, model: gpt-6.1-sol, reasoning_effort: high}
    deviation: {from: 2026-09-30, reason: Keep the validated model}
`

func TestUpdateCheckNamesEachProfileThatDiffers(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"transcript", "json", "pinned_and_current_add_nothing", "user_config_source", "no_profiles_keep_the_report_bytes", "newer_baseline_reports_no_member", "user_pin", "collapsed_reason", "bounded_reason", "release_update_preserves_next"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			root := newBaselinePlanTestRepository(t)
			setCommandEnvironmentForTest(t, home, root)
			latest := "99.0.0"
			if name == "release_update_preserves_next" {
				latest = "100.0.0"
			}
			configureUpdateCheckTest(t, "99.0.0", func(context.Context) ([]app.Release, error) {
				return []app.Release{{Version: latest, Titles: []string{"Fixture release"}}}, nil
			})
			profiles := "profiles:\n" + updateCheckOldDocsProfile + updateCheckOldReviewProfile
			baselineState := "action_required"
			baselineMessage := "the repository has no Setup Manifest and cannot be updated"
			wantErr := ""
			authoringJSON, authoringText := "", ""
			source := "project"
			reason := "Keep the validated model"
			nextWithoutProfiles := ""
			switch name {
			case "pinned_and_current_add_nothing":
				docs, ok := roundconfig.RecommendedProfile(roundconfig.CategoryDocs)
				if !ok {
					t.Fatal("missing docs recommendation")
				}
				profiles = "profiles:\n" + updateCheckProfileYAML("docs", docs) + strings.Replace(updateCheckOldReviewProfile, "2026-09-30", roundconfig.ModelRecommendationSnapshotVersion, 1)
			case "no_profiles_keep_the_report_bytes":
				profiles = "# No profile configuration.\n"
			case "user_config_source":
				profiles = "profiles:\n" + updateCheckOldDocsProfile
				source = "user"
			case "user_pin":
				profiles = "profiles:\n" + updateCheckOldReviewProfile
				source = "user"
			case "collapsed_reason":
				profiles = strings.Replace(profiles, "reason: Keep the validated model", `reason: "Keep\n  the\tvalidated model"`, 1)
			case "bounded_reason":
				profiles = strings.Replace(profiles, "Keep the validated model", strings.Repeat("é", 301), 1)
				reason = strings.Repeat("é", 300)
			case "json":
				// Include the preceding report to prove recommendations follow it.
				writeBaselinePlanTestFile(t, root, ".claude/agents/roundfix-prd-author.md", "---\nname: roundfix-prd-author\n---\n\n<!-- roundfix:generated skills.authoring.claude.prd -->\n")
				authoringJSON = `,"authoringModels":{"state":"unapplied","harness":"claude","declared":[],"problems":[".claude/agents/roundfix-prd-author.md has no declared entry"],"next":"roundfix skills install --target project"}`
				authoringText = "Authoring models: the installed files do not match the declaration: .claude/agents/roundfix-prd-author.md has no declared entry; apply with roundfix skills install --target project.\n"
			case "newer_baseline_reports_no_member":
				writeBaselinePlanTestFile(t, root, roundconfig.SetupManifestPath, `{"roundfixVersion":"`+futureRecordedRoundfixVersion+`"}`)
				refusal := fmt.Sprintf("Setup Manifest %q records Roundfix %s, newer than this binary's %s; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'", filepath.Join(root, roundconfig.SetupManifestPath), futureRecordedRoundfixVersion, app.NormalizeVersion(app.Version))
				baselineState, baselineMessage = "failed", refusal
				wantErr = "roundfix: warning: " + refusal + "; Project Config not read\n"
				// Also suppress differences in User Config when Project Config is unread.
				writeUpdateCheckUserConfig(t, home, "profiles:\n"+updateCheckOldReviewProfile)
			}
			if name == "json" {
				writeBaselinePlanTestFile(t, root, ".roundfixrc.yml", "# No profiles yet.\n")
				commitBaselinePlanTestRepository(t, root)
				raw := runUpdateCheckTest(t, "update", "check", "--refresh", "--format", "json")
				var report struct {
					Next string `json:"next"`
				}
				if err := json.Unmarshal([]byte(raw), &report); err != nil {
					t.Fatal(err)
				}
				nextWithoutProfiles = report.Next
			}
			if source == "user" {
				writeUpdateCheckUserConfig(t, home, profiles)
				writeBaselinePlanTestFile(t, root, ".roundfixrc.yml", "# Profiles are in User Config.\n")
			} else {
				writeBaselinePlanTestFile(t, root, ".roundfixrc.yml", profiles)
			}
			commitBaselinePlanTestRepository(t, root)

			// These are the existing report transcripts without a profile comparison.
			prefix := "Roundfix 99.0.0 is current (latest release 99.0.0; update.mode notify).\n"
			releasesJSON, next, nextText := "[]", "none", ""
			state := "current"
			if latest != "99.0.0" {
				prefix = "Roundfix 99.0.0 → 100.0.0 is available (update.mode notify).\nChanges since 99.0.0:\n- 100.0.0: Fixture release\n"
				releasesJSON = `[{"version":"100.0.0","titles":["Fixture release"],"url":""}]`
				state, next = "update_available", "apply"
				nextText = "Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.\n"
			}
			prefix += "Baseline: " + baselineState + ": " + baselineMessage + "\n" + authoringText
			wantJSON := `{"schemaVersion":"roundfix/update-check/v1","mode":"notify","installed":"99.0.0","latest":` + updateCheckRecommendationJSON(t, latest) + `,"state":` + updateCheckRecommendationJSON(t, state) + `,"lookup":{"source":"network","checkedAt":"2026-11-06T15:00:00Z","error":""},"releases":` + releasesJSON + `,"releasesTruncated":false,"baseline":{"state":` + updateCheckRecommendationJSON(t, baselineState) + `,"fileChanges":0,"historyUnits":0,"skillChanges":0,"message":` + updateCheckRecommendationJSON(t, baselineMessage) + `}` + authoringJSON
			profileText, categoriesJSON := "", []string{}
			if name != "pinned_and_current_add_nothing" && name != "no_profiles_keep_the_report_bytes" && name != "newer_baseline_reports_no_member" {
				snapshot := roundconfig.ModelRecommendationSnapshotVersion
				docs, _ := roundconfig.RecommendedProfile(roundconfig.CategoryDocs)
				review, _ := roundconfig.RecommendedProfile(roundconfig.CategoryReview)
				adopt := "roundfix profiles check --apply --scope " + source
				if name != "user_pin" {
					configured := roundconfig.AgentSelectionProfile{Preferred: roundconfig.AgentSelection{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"}, Fallbacks: []roundconfig.AgentSelection{{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"}}}
					profileText += fmt.Sprintf("Profiles: docs (%s) differs from the Recommended Profile of %s: configured codex / gpt-5.6-luna / max, then claude / sonnet / high; recommended %s; adopt with %s.\n", source, snapshot, formatRecommendationProfile(docs), adopt)
					categoriesJSON = append(categoriesJSON, `{"category":"docs","state":"differs","source":`+updateCheckRecommendationJSON(t, source)+`,"configured":`+updateCheckExpectedProfileJSON(t, configured)+`,"recommended":`+updateCheckExpectedProfileJSON(t, docs)+`,"adopt":`+updateCheckRecommendationJSON(t, adopt)+`}`)
				}
				if name != "user_config_source" {
					label := "Project"
					if source == "user" {
						label = "User"
					}
					repin := "profiles.review.deviation.from: " + snapshot
					profileText += fmt.Sprintf("Profiles: review (%s) is pinned to an older snapshot (2026-09-30: %s); recommended %s as of %s; re-pin with %s in %s Config, or adopt with %s.\n", source, reason, formatRecommendationProfile(review), snapshot, repin, label, adopt)
					rawReason := reason
					if name == "collapsed_reason" {
						rawReason = "Keep\n  the\tvalidated model"
					}
					if name == "bounded_reason" {
						rawReason = strings.Repeat("é", 301)
					}
					configured := roundconfig.AgentSelectionProfile{Preferred: roundconfig.AgentSelection{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"}, Fallbacks: []roundconfig.AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}}}
					categoriesJSON = append(categoriesJSON, `{"category":"review","state":"pinned_to_older_snapshot","source":`+updateCheckRecommendationJSON(t, source)+`,"configured":`+updateCheckExpectedProfileJSON(t, configured)+`,"recommended":`+updateCheckExpectedProfileJSON(t, review)+`,"deviation":{"from":"2026-09-30","reason":`+updateCheckRecommendationJSON(t, rawReason)+`},"repin":`+updateCheckRecommendationJSON(t, repin)+`,"adopt":`+updateCheckRecommendationJSON(t, adopt)+`}`)
				}
				wantJSON += `,"recommendations":{"snapshot":` + updateCheckRecommendationJSON(t, snapshot) + `,"categories":[` + strings.Join(categoriesJSON, ",") + `]}`
			}
			if name == "json" && next != nextWithoutProfiles {
				t.Fatalf("next=%q without profiles=%q", next, nextWithoutProfiles)
			}
			wantJSON += `,"next":` + updateCheckRecommendationJSON(t, next) + "}\n"
			// Prime the permitted cache write so snapshots also detect any home-file
			// changes made by the first profile comparison.
			writeUpdateCheckCache(filepath.Join(home, ".roundfix", updateCheckCacheFile), updateCheckCache{
				SchemaVersion: updateCheckCacheSchema, CheckedAt: updateCheckTestTime(),
				Releases: []app.Release{{Version: latest, Titles: []string{"Fixture release"}}},
			})
			beforeHome := snapshotDoctorPath(t, home)
			before := baselinePlanTestTree(t, root)
			for _, format := range []string{"text", "json"} {
				var stdout, stderr bytes.Buffer
				code := runCLIContext(t, t.Context(), []string{"update", "check", "--refresh", "--format", format}, &stdout, &stderr)
				want := prefix + profileText + nextText
				if format == "json" {
					want = wantJSON
				}
				if code != 0 || stdout.String() != want || stderr.String() != wantErr {
					t.Fatalf("%s: exit=%d\nstdout=%s\nwant stdout=%s\nstderr=%q want stderr=%q", format, code, &stdout, want, &stderr, wantErr)
				}
			}
			if after := baselinePlanTestTree(t, root); after != before {
				t.Fatal("Update Check changed repository files or modes")
			}
			// The existing release cache is the command's only permitted write.
			var cachedOut, cachedErr bytes.Buffer
			code := runCLIContext(t, t.Context(), []string{"update", "check", "--format", "json"}, &cachedOut, &cachedErr)
			wantCached := strings.Replace(wantJSON, `"source":"network"`, `"source":"cache"`, 1)
			if code != 0 || cachedOut.String() != wantCached || cachedErr.String() != wantErr {
				t.Fatalf("cached: exit=%d stdout=%s want=%s stderr=%q want=%q", code, &cachedOut, wantCached, &cachedErr, wantErr)
			}
			if !reflect.DeepEqual(beforeHome, snapshotDoctorPath(t, home)) {
				t.Fatal("profile comparison changed User Config or other home files")
			}
		})
	}
}

func updateCheckProfileYAML(category string, profile roundconfig.AgentSelectionProfile) string {
	selection := func(s roundconfig.AgentSelection) string {
		return fmt.Sprintf("{runtime: %s, model: %s, reasoning_effort: %s}", s.Runtime, s.Model, s.ReasoningEffort)
	}
	text := "  " + category + ":\n    preferred: " + selection(profile.Preferred) + "\n    fallbacks:\n"
	for _, fallback := range profile.Fallbacks {
		text += "      - " + selection(fallback) + "\n"
	}
	return text
}

func updateCheckExpectedProfileJSON(t *testing.T, profile roundconfig.AgentSelectionProfile) string {
	t.Helper()
	return `{"preferred":` + updateCheckRecommendationJSON(t, profile.Preferred) + `,"fallbacks":` + updateCheckRecommendationJSON(t, profile.Fallbacks) + `}`
}

func updateCheckRecommendationJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func TestUpdateCheckOmitsAnUnavailableProfileComparison(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"comparison_error", "unread_config"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			loaded := roundconfig.Loaded{Config: roundconfig.Builtin()}
			loaded.Config.Profiles[roundconfig.CategoryReview] = roundconfig.ProfileEntry{
				Profile: roundconfig.AgentSelectionProfile{Preferred: roundconfig.AgentSelection{Runtime: "codex", Model: "old-model", ReasoningEffort: "high"}, Fallbacks: []roundconfig.AgentSelection{{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"}}},
				Source:  roundconfig.ProfileSourceProject,
			}
			if name == "comparison_error" {
				delete(loaded.Config.Profiles, roundconfig.CategoryGeneral)
			} else {
				loaded.RepositoryVersionErr = errors.New("newer Baseline")
			}
			if report := updateCheckProfileRecommendations(loaded); report != nil {
				t.Fatalf("unavailable comparison produced %+v", report)
			}
			var stdout bytes.Buffer
			printUpdateCheckRecommendations(&stdout, nil)
			if stdout.String() != "" {
				t.Fatalf("nil report stdout=%q", &stdout)
			}
		})
	}
}
