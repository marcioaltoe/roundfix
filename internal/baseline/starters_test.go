// Suite: confirmed Starter Files
// Invariant: only absent, confirmed paths are written; failure removes created files.
// Boundary IN: embedded catalog, temporary repositories, and explicit profile selection.
// Boundary OUT: CLI, repository commands, and external tools.
package baseline

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestStartersPreviewsAndCreatesOnlyAbsentFiles(t *testing.T) {
	t.Parallel()
	repository := newPlanRepository(t)
	catalog := mustEmbeddedCatalog(t)
	original, err := os.ReadFile(filepath.Join(repository, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	request := StartersRequest{Repository: repository, ProfileID: "python"}
	preview, err := Starters(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if preview.SchemaVersion != StartersSchemaVersion || preview.Profile != "python" || preview.State != "preview" || !strings.HasPrefix(preview.PlanDigest, "sha256:") || len(preview.PlanDigest) != 71 {
		t.Fatalf("preview = %+v", preview)
	}
	files, err := catalog.StarterFiles([]string{"python"})
	if err != nil {
		t.Fatal(err)
	}
	want := []StarterEntry{
		{Path: "Makefile", Action: "keep", ContentIdentity: files[0].ContentIdentity},
		{Path: ".roundfixrc.yml", Action: "create", ContentIdentity: files[1].ContentIdentity},
		{Path: ".github/workflows/verify.yml", Action: "create", ContentIdentity: files[2].ContentIdentity},
	}
	if !reflect.DeepEqual(preview.Starters, want) {
		t.Fatalf("entries = %+v, want %+v", preview.Starters, want)
	}
	for _, relative := range []string{".roundfixrc.yml", ".github"} {
		if _, err := os.Lstat(filepath.Join(repository, relative)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("preview created %s: %v", relative, err)
		}
	}
	t.Run("repository-owned profile retains starters", func(t *testing.T) {
		t.Parallel()
		customRepository := newPlanRepository(t)
		custom, err := InitCustomProfile(customRepository, "repository-python", "python", catalog)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Starters(t.Context(), StartersRequest{Repository: customRepository, ProfileID: custom.Profile.ID}, catalog)
		if err != nil || got.State != "preview" || got.Profile != custom.Profile.ID || !reflect.DeepEqual(got.Starters, preview.Starters) {
			t.Fatalf("repository profile preview = %+v, %v", got, err)
		}
	})
	request.ConfirmPlan = "sha256:unconfirmed"
	if _, err := Starters(t.Context(), request, catalog); !errors.Is(err, ErrStartersStale) {
		t.Fatalf("wrong digest error = %v", err)
	}
	for _, relative := range []string{".roundfixrc.yml", ".github"} {
		if _, err := os.Lstat(filepath.Join(repository, relative)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("wrong digest created %s: %v", relative, err)
		}
	}
	request.ConfirmPlan = preview.PlanDigest
	applied, err := Starters(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if applied.State != "applied" || applied.PlanDigest != preview.PlanDigest || !reflect.DeepEqual(applied.Starters, preview.Starters) {
		t.Fatalf("apply = %+v", applied)
	}
	for _, file := range files[1:] {
		data, err := os.ReadFile(filepath.Join(repository, filepath.FromSlash(file.Path)))
		if err != nil || !bytes.Equal(data, file.Content) {
			t.Fatalf("created %s: %q, %v", file.Path, data, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(repository, "Makefile"))
	if err != nil || !bytes.Equal(data, original) {
		t.Fatalf("Makefile changed: %q, %v", data, err)
	}
	request.ConfirmPlan = ""
	current, err := Starters(t.Context(), request, catalog)
	if err != nil || current.State != "current" {
		t.Fatalf("current = %+v, %v", current, err)
	}
	for _, entry := range current.Starters {
		if entry.Action != "keep" {
			t.Fatalf("current entry = %+v", entry)
		}
	}
	request.ConfirmPlan = preview.PlanDigest
	current, err = Starters(t.Context(), request, catalog)
	if err != nil || current.State != "current" {
		t.Fatalf("confirmed current = %+v, %v", current, err)
	}
}

func TestStartersRefuseADigestThatNoLongerMatches(t *testing.T) {
	t.Parallel()
	repository := newPlanRepository(t)
	catalog := mustEmbeddedCatalog(t)
	request := StartersRequest{Repository: repository, ProfileID: "python"}
	preview, err := Starters(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	config := []byte("repository owns this config\n")
	if err := os.WriteFile(filepath.Join(repository, ".roundfixrc.yml"), config, 0o644); err != nil {
		t.Fatal(err)
	}
	request.ConfirmPlan = preview.PlanDigest
	if _, err := Starters(t.Context(), request, catalog); !errors.Is(err, ErrStartersStale) {
		t.Fatalf("confirmation error = %v", err)
	}
	if _, err := os.Lstat(filepath.Join(repository, ".github")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stale confirmation wrote workflow directory: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(repository, ".roundfixrc.yml"))
	if err != nil || !bytes.Equal(data, config) {
		t.Fatalf("existing config changed: %q, %v", data, err)
	}
}

func TestStartersRemoveWhatTheyCreatedWhenAWriteFails(t *testing.T) {
	t.Parallel()
	repository := newPlanRepository(t)
	catalog := mustEmbeddedCatalog(t)
	obstruction := []byte("existing regular file\n")
	if err := os.WriteFile(filepath.Join(repository, ".github"), obstruction, 0o644); err != nil {
		t.Fatal(err)
	}
	request := StartersRequest{Repository: repository, ProfileID: "python"}
	preview, err := Starters(t.Context(), request, catalog)
	if err != nil {
		t.Fatal(err)
	}
	request.ConfirmPlan = preview.PlanDigest
	if _, err := Starters(t.Context(), request, catalog); err == nil || !strings.Contains(err.Error(), "starters.write.failed: .github/workflows/verify.yml:") {
		t.Fatalf("write error = %v", err)
	}
	if _, err := os.Lstat(filepath.Join(repository, ".roundfixrc.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("created config remains: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(repository, ".github"))
	if err != nil || !bytes.Equal(data, obstruction) {
		t.Fatalf("obstruction changed: %q, %v", data, err)
	}
}

func TestTheUndeclaredGateNamesTheStartersCommand(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"python", "rust-cli"} {
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			repository := newPlanRepository(t)
			if err := os.Remove(filepath.Join(repository, "Makefile")); err != nil {
				t.Fatal(err)
			}
			catalog := mustEmbeddedCatalog(t)
			resolved, err := ResolveProfile(repository, profile, catalog)
			if err != nil {
				t.Fatal(err)
			}
			decisions := slices.DeleteFunc(standardTypeScriptDecisions("make verify"), func(decision DecisionValue) bool {
				return !slices.Contains(resolved.Decisions, decision.ID)
			})
			alignment, err := ResolveProfileAlignment(t.Context(), repository, ProfileAlignmentRequest{
				ProfileID: profile, Decisions: decisions, ExecutableDirectories: []string{t.TempDir()},
			}, catalog)
			if err != nil {
				t.Fatal(err)
			}
			divergence, ok := findProfileDivergence(alignment.Divergences, "verification.gate")
			want := "select a command declared by the repository or add its local declaration"
			if profile == "python" {
				want += "; `roundfix baseline starters --profile python` previews this profile's Starter Files"
			}
			if !ok || divergence.Code != "verification.command.undeclared" || divergence.NextAction != want {
				t.Fatalf("divergence = %+v, found=%v; want next action %q", divergence, ok, want)
			}
		})
	}
}
