package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func runProfilesApply(t *testing.T, flags ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	args := append([]string{"profiles", "check", "--apply"}, flags...)
	code := runCLI(t, args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func requireProfilesApplyUnchanged(t *testing.T, before map[string][]byte, home, repo string) {
	t.Helper()
	if after := profilesCheckFileBytes(t, home, repo); !reflect.DeepEqual(before, after) {
		t.Fatalf("files changed: before=%v after=%v", before, after)
	}
}

func TestProfilesCheckApplyWritesTheRecommendedProfileOfEachDifferingCategory(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	userBefore := mustRead(t, filepath.Join(home, ".roundfix", "config.yml"))
	runner := withSuccessfulProfilesConfigureProof(t)
	code, out, errOut := runProfilesApply(t, "--scope", "project", "--yes", "--json")
	if code != exitOK || errOut != "" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	response := decodeProfilesConfigureResponse(t, out)
	if response.Schema != profilesConfigureSchema || !response.Changed || len(response.Profiles) != 3 || len(response.Changes) != 3 {
		t.Fatalf("response=%+v", response)
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	wantedProofs := map[roundconfig.AgentSelection]bool{}
	for _, category := range []roundconfig.WorkCategory{roundconfig.CategoryBackend, roundconfig.CategoryDocs, roundconfig.CategoryQA} {
		got, err := roundconfig.ResolveProfile(loaded.Config, category, nil)
		if err != nil {
			t.Fatal(err)
		}
		want, _ := roundconfig.RecommendedProfile(category)
		if !reflect.DeepEqual(got.Profile, want) || got.Deviation != nil {
			t.Fatalf("%s profile=%+v want=%+v", category, got, want)
		}
		wantedProofs[want.Preferred] = true
		for _, fallback := range want.Fallbacks {
			wantedProofs[fallback] = true
		}
	}
	proved := map[roundconfig.AgentSelection]bool{}
	for _, req := range runner.exactRequests {
		proved[roundconfig.AgentSelection{Runtime: req.Runtime.ID, Model: req.Runtime.Model, ReasoningEffort: req.Runtime.ReasoningEffort}] = true
	}
	if !reflect.DeepEqual(proved, wantedProofs) {
		t.Fatalf("proved=%v want=%v", proved, wantedProofs)
	}
	if mustRead(t, filepath.Join(home, ".roundfix", "config.yml")) != userBefore {
		t.Fatal("project adoption changed User Config")
	}
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"profiles", "check"}, &stdout, &stderr); code != exitOK || !strings.Contains(stdout.String(), "0 differ, 1 pinned") {
		t.Fatalf("following check exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
}

func TestProfilesCheckApplyLeavesAPinnedCategoryUntouched(t *testing.T) {
	t.Parallel()
	_, repo := writeProfilesCheckFixture(t)
	path := filepath.Join(repo, ".roundfixrc.yml")
	original := mustRead(t, path)
	pinned := original[strings.Index(original, "  review:"):]
	withSuccessfulProfilesConfigureProof(t)
	code, out, errOut := runProfilesApply(t, "--scope", "project", "--yes", "--json")
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	response := decodeProfilesConfigureResponse(t, out)
	for _, profile := range response.Profiles {
		if profile.Category == roundconfig.CategoryReview || profile.Category == roundconfig.CategoryGeneral || profile.Category == roundconfig.CategoryFrontend {
			t.Fatalf("selected pinned/current category: %+v", profile)
		}
	}
	after := mustRead(t, path)
	start := strings.Index(after, "  review:")
	if start < 0 || len(after)-start < len(pinned) || after[start:start+len(pinned)] != pinned {
		t.Fatalf("pinned bytes changed: %s", after)
	}
}

func TestProfilesCheckApplyProvesBeforeItWrites(t *testing.T) {
	t.Parallel()
	for _, fail := range []bool{false, true} {
		name := "success"
		if fail {
			name = "failed proof"
		}
		t.Run(name, func(t *testing.T) {
			home, repo := writeProfilesCheckFixture(t)
			before := profilesCheckFileBytes(t, home, repo)
			calls := 0
			runner := &profileReadinessExactRunner{prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
				calls++
				requireProfilesApplyUnchanged(t, before, home, repo)
				if fail {
					return agent.SelectionProof{}, errors.New("proof refused")
				}
				return agent.SelectionProof{Runtime: req.Runtime.ID, Model: req.Runtime.Model, ReasoningEffort: req.Runtime.ReasoningEffort, Status: agent.SelectionProofStatusProven}, nil
			}}
			withAgentRunner(t, runner)
			withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
				if calls == 0 {
					t.Fatal("confirmation before proof")
				}
				requireProfilesApplyUnchanged(t, before, home, repo)
				return true, nil
			})
			code, out, errOut := runProfilesApply(t, "--scope", "project")
			if calls == 0 {
				t.Fatal("no exact proof attempted")
			}
			if fail {
				if code != exitPreflight || !strings.Contains(errOut, "proof refused") {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
				}
				requireProfilesApplyUnchanged(t, before, home, repo)
			} else {
				if code != exitOK {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
				}
				if reflect.DeepEqual(before, profilesCheckFileBytes(t, home, repo)) {
					t.Fatal("no write after successful proof and confirmation")
				}
			}
		})
	}
}

func TestProfilesCheckApplyDryRunWritesNothing(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	before := profilesCheckFileBytes(t, home, repo)
	runner := withSuccessfulProfilesConfigureProof(t)
	withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) { t.Fatal("dry-run prompted"); return false, nil })
	code, out, errOut := runProfilesApply(t, "--scope", "project", "--dry-run")
	if code != exitOK || !strings.Contains(out, "Profile Configure Preview") || !strings.Contains(out, "Profile configuration dry run:") || len(runner.exactRequests) == 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	requireProfilesApplyUnchanged(t, before, home, repo)
}

func TestProfilesCheckApplyDeclinedConfirmationWritesNothing(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	before := profilesCheckFileBytes(t, home, repo)
	runner := withSuccessfulProfilesConfigureProof(t)
	withProfilesConfigureConfirm(t, defaultConfirmProfilesConfigure)
	withProfilesConfigureInput(t, "n\n")
	code, out, errOut := runProfilesApply(t, "--scope", "project", "--json")
	if code != exitRunFailed || !decodeProfilesConfigureRefused(t, out) || !strings.Contains(errOut, "confirmation declined") || len(runner.exactRequests) == 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	requireProfilesApplyUnchanged(t, before, home, repo)
}

func TestProfilesCheckApplyUserScopeNamesProjectDefinedCategories(t *testing.T) {
	t.Parallel()
	home, repo := writeProfilesCheckFixture(t)
	projectBefore := mustRead(t, filepath.Join(repo, ".roundfixrc.yml"))
	withSuccessfulProfilesConfigureProof(t)
	code, out, errOut := runProfilesApply(t, "--scope", "user", "--yes", "--json")
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	for _, category := range []string{"backend", "qa"} {
		if !strings.Contains(errOut, "recommendations: "+category+" is defined by Project Config; use --scope project") {
			t.Fatalf("missing skip advice: %s", errOut)
		}
	}
	response := decodeProfilesConfigureResponse(t, out)
	if len(response.Profiles) != 1 || response.Profiles[0].Category != roundconfig.CategoryDocs {
		t.Fatalf("response=%+v", response)
	}
	if mustRead(t, filepath.Join(repo, ".roundfixrc.yml")) != projectBefore {
		t.Fatal("user adoption changed Project Config bytes")
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	docs, err := roundconfig.ResolveProfile(loaded.Config, roundconfig.CategoryDocs, nil)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := roundconfig.RecommendedProfile(roundconfig.CategoryDocs)
	if !reflect.DeepEqual(docs.Profile, want) || docs.Deviation != nil {
		t.Fatalf("docs=%+v", docs)
	}
}

func TestProfilesCheckApplyWithNothingToAdoptChangesNothing(t *testing.T) {
	t.Parallel()
	for _, fixture := range []string{"builtins", "pinned", "project excluded"} {
		t.Run(fixture, func(t *testing.T) {
			home, repo := withCLIWorkspace(t)
			scope := "project"
			if fixture != "builtins" {
				// Reuse the mixed fixture, then retain only its pinned or project-defined category.
				home, repo = writeProfilesCheckFixture(t)
				path := filepath.Join(repo, ".roundfixrc.yml")
				original := mustRead(t, path)
				if fixture == "pinned" {
					mustWrite(t, path, "profiles:\n"+original[strings.Index(original, "  review:"):])
				} else {
					mustWrite(t, path, "profiles:\n"+original[strings.Index(original, "  backend:"):strings.Index(original, "  qa:")])
					scope = "user"
				}
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "watch:\n  max_rounds: 4\n")
			}
			before := profilesCheckFileBytes(t, home, repo)
			withAgentRunner(t, profilesCheckForbiddenRunner{t: t})
			withProfilesConfigureConfirm(t, func(context.Context, io.Writer, string) (bool, error) {
				t.Fatal("empty adoption prompted")
				return false, nil
			})
			for _, jsonOutput := range []bool{false, true} {
				flags := []string{"--scope", scope}
				if jsonOutput {
					flags = append(flags, "--json")
				}
				code, out, errOut := runProfilesApply(t, flags...)
				if code != exitOK {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, errOut)
				}
				if jsonOutput {
					response := decodeProfilesConfigureResponse(t, out)
					if response.Schema != profilesConfigureSchema || response.Changed || response.Refused || response.Error != "" || len(response.Profiles) != 0 || len(response.Changes) != 0 || !strings.Contains(out, "\"profiles\":[]") || !strings.Contains(out, "\"changes\":[]") {
						t.Fatalf("response=%s", out)
					}
				} else if out != "Profile configuration unchanged: nothing to adopt\n" {
					t.Fatalf("output=%q", out)
				}
				requireProfilesApplyUnchanged(t, before, home, repo)
			}
		})
	}
}

func TestProfilesCheckApplyFlagRules(t *testing.T) {
	t.Parallel()
	for _, flags := range [][]string{
		{"--apply"}, {"--scope", "project"}, {"--dry-run"}, {"--yes"}, {"--apply", "--scope", "invalid"},
		{"--apply", "--scope", "project", "extra"}, {"--apply", "--scope", "project", "--unknown"},
		{"--apply=false", "--scope", "project", "--yes"}, {"--apply", "--scope", "project", "--yes=invalid"},
	} {
		t.Run(strings.Join(flags, " "), func(t *testing.T) {
			home, repo := writeProfilesCheckFixture(t)
			before := profilesCheckFileBytes(t, home, repo)
			withAgentRunner(t, profilesCheckForbiddenRunner{t: t})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, append([]string{"profiles", "check"}, flags...), &stdout, &stderr)
			if code != exitPreflight || stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			requireProfilesApplyUnchanged(t, before, home, repo)
		})
	}
}
