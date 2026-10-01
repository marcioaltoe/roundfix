package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
)

func writeProfilesCheckFixture(t *testing.T) (string, string) {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "profiles:\n  docs:\n    preferred: {runtime: codex, model: docs-choice, reasoning_effort: high}\n    fallbacks:\n      - {runtime: claude, model: sonnet, reasoning_effort: high}\n")
	content := fmt.Sprintf(`profiles:
  backend:
    preferred: {runtime: codex, model: chosen-model, reasoning_effort: high}
    fallbacks:
      - {runtime: codex, model: backup-model, reasoning_effort: xhigh}
  qa:
    preferred: {runtime: codex, model: older-choice, reasoning_effort: high}
    fallbacks:
      - {runtime: claude, model: opus, reasoning_effort: high}
    deviation: {from: "2000-01-01", reason: Older choice}
  review:
    preferred: {runtime: codex, model: pinned-model, reasoning_effort: high}
    fallbacks:
      - {runtime: codex, model: backup-model, reasoning_effort: max}
    deviation: {from: %q, reason: Keep the validated model}
`, roundconfig.ModelRecommendationSnapshotVersion)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), content)
	return home, repo
}

func TestProfilesCheckPrintsEachDifferenceAndTheSummary(t *testing.T) {
	t.Parallel()
	writeProfilesCheckFixture(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	backend, _ := roundconfig.RecommendedProfile(roundconfig.CategoryBackend)
	docs, _ := roundconfig.RecommendedProfile(roundconfig.CategoryDocs)
	qa, _ := roundconfig.RecommendedProfile(roundconfig.CategoryQA)
	want := fmt.Sprintf("recommendations: backend differs from snapshot %[1]s\n  configured:  codex / chosen-model / high, then codex / backup-model / xhigh (project)\n  recommended: %[2]s\nrecommendations: docs differs from snapshot %[1]s\n  configured:  codex / docs-choice / high, then claude / sonnet / high (user)\n  recommended: %[3]s\nrecommendations: qa differs from snapshot %[1]s (its deviation was declared against 2000-01-01)\n  configured:  codex / older-choice / high, then claude / opus / high (project)\n  recommended: %[4]s\nrecommendations: review pinned against snapshot %[1]s: Keep the validated model\nrecommendations: snapshot %[1]s; 2 current, 3 differ, 1 pinned\nrecommendations: adopt with `roundfix profiles check --apply --scope user|project`, or declare a deviation\n", roundconfig.ModelRecommendationSnapshotVersion, formatRecommendationProfile(backend), formatRecommendationProfile(docs), formatRecommendationProfile(qa))
	if stdout.String() != want {
		t.Fatalf("text:\n%s\nwant:\n%s", &stdout, want)
	}
}

func TestProfilesCheckExitsZeroWithDifferences(t *testing.T) {
	t.Parallel()
	writeProfilesCheckFixture(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	if !strings.Contains(stdout.String(), "3 differ, 1 pinned") {
		t.Fatalf("missing difference counts: %s", &stdout)
	}
}

func TestProfilesCheckBuiltinsPrintOnlyTheSummary(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	want := fmt.Sprintf("recommendations: snapshot %s; %d current, 0 differ, 0 pinned\n", roundconfig.ModelRecommendationSnapshotVersion, len(roundconfig.RequiredWorkCategories()))
	if stdout.String() != want {
		t.Fatalf("output = %q, want %q", stdout.String(), want)
	}
}

func TestProfilesCheckJSONIsSchemaV1(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check", "--json"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	var response struct {
		Schema     string `json:"schema"`
		Snapshot   string `json:"snapshot"`
		Current    int    `json:"current"`
		Differ     int    `json:"differ"`
		Pinned     int    `json:"pinned"`
		Categories []struct {
			Category   roundconfig.WorkCategory  `json:"category"`
			Status     string                    `json:"status"`
			Source     roundconfig.ProfileSource `json:"source"`
			Configured struct {
				Preferred roundconfig.AgentSelection   `json:"preferred"`
				Fallbacks []roundconfig.AgentSelection `json:"fallbacks"`
			} `json:"configured"`
			Recommended struct {
				Preferred roundconfig.AgentSelection   `json:"preferred"`
				Fallbacks []roundconfig.AgentSelection `json:"fallbacks"`
			} `json:"recommended"`
			Deviation *roundconfig.ProfileDeviation `json:"deviation"`
		} `json:"categories"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Schema != "roundfix/profiles-check/v1" || response.Snapshot != roundconfig.ModelRecommendationSnapshotVersion || response.Current != 2 || response.Differ != 3 || response.Pinned != 1 {
		t.Fatalf("header = %+v", response)
	}
	wantCategories := []roundconfig.WorkCategory{roundconfig.CategoryGeneral, roundconfig.CategoryBackend, roundconfig.CategoryFrontend, roundconfig.CategoryDocs, roundconfig.CategoryQA, roundconfig.CategoryReview}
	wantStatuses := []string{"current", "differs", "current", "differs", "differs", "pinned"}
	if len(response.Categories) != len(wantCategories) {
		t.Fatalf("categories = %+v", response.Categories)
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range response.Categories {
		configured, err := roundconfig.ResolveProfile(loaded.Config, wantCategories[i], nil)
		if err != nil {
			t.Fatal(err)
		}
		recommended, _ := roundconfig.RecommendedProfile(wantCategories[i])
		if row.Category != wantCategories[i] || row.Status != wantStatuses[i] || row.Source != configured.Source || row.Configured.Preferred != configured.Profile.Preferred || !reflect.DeepEqual(row.Configured.Fallbacks, configured.Profile.Fallbacks) || row.Recommended.Preferred != recommended.Preferred || !reflect.DeepEqual(row.Recommended.Fallbacks, recommended.Fallbacks) || !reflect.DeepEqual(row.Deviation, configured.Deviation) {
			t.Errorf("category %d incorrect: %+v", i, row)
		}
	}
	// Additions such as deviation are absent when there is no declaration.
	var raw struct {
		Categories []map[string]json.RawMessage `json:"categories"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if _, present := raw.Categories[0]["deviation"]; present {
		t.Fatal("absent deviation serialized")
	}
}

func TestProfilesCheckRefusesUnknownFlagsAndArguments(t *testing.T) {
	t.Parallel()
	for _, arg := range []string{"--unknown", "--apply", "--scope", "--dry-run", "--yes", "extra"} {
		t.Run(arg, func(t *testing.T) {
			withCLIWorkspace(t)
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"profiles", "check", arg}, &stdout, &stderr); code != 2 || stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
		})
	}
	t.Run("configuration does not load", func(t *testing.T) {
		_, repo := withCLIWorkspace(t)
		mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "profiles: [broken]\n")
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"profiles", "check", "--json"}, &stdout, &stderr); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "profiles must be a mapping") {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
		}
	})
}

type profilesCheckForbiddenRunner struct{ t *testing.T }

func (r profilesCheckForbiddenRunner) Probe(context.Context, agent.ProbeRequest) error {
	r.t.Fatal("profiles check called runner Probe")
	return nil
}
func (r profilesCheckForbiddenRunner) Run(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	r.t.Fatal("profiles check opened an Agent Session")
	return agent.ExecuteResult{}, nil
}

func (r profilesCheckForbiddenRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	r.t.Fatal("profiles check called runner EndSession")
	return nil
}
func (r profilesCheckForbiddenRunner) ProveProfileSelection(context.Context, agent.ProbeRequest) (agent.SelectionProof, error) {
	r.t.Fatal("profiles check opened a proof Session")
	return agent.SelectionProof{}, nil
}

func profilesCheckFileBytes(t *testing.T, roots ...string) map[string][]byte {
	t.Helper()
	files := make(map[string][]byte)
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files[path] = content
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return files
}

func TestProfilesCheckOpensNoAgentSessionAndWritesNothing(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	withAgentRunner(t, profilesCheckForbiddenRunner{t: t})
	before := profilesCheckFileBytes(t, home, repo)
	for _, args := range [][]string{{"profiles", "check"}, {"profiles", "check", "--json"}} {
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, args, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
			t.Fatalf("exit=%d stderr=%s", code, &stderr)
		}
		after := profilesCheckFileBytes(t, home, repo)
		if !reflect.DeepEqual(before, after) {
			t.Fatalf("configuration or file set changed: before=%v after=%v", before, after)
		}
	}
}

func TestProfilesCheckHelp(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check", "--help"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	for _, want := range []string{"roundfix profiles check [--json]", "roundfix/profiles-check/v1", "Read-only and offline", "Exit 0"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("help missing %q", want)
		}
	}
}
