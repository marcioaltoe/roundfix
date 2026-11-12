// Suite: Recorded Roundfix Version in Baseline plans
// Invariant: a release record moves only upward and only with a file change.
// Boundary IN: plan postimages, verified apply, Managed Refresh, and manifest validation.
// Boundary OUT: CLI refusal and release transport.
package baseline

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAPlanRecordsItsWriterVersionOnlyWithAFileChange(t *testing.T) {
	t.Parallel()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, true)
	request := baselinePlanCharacterizationRequest(repository, baselinePlanCharacterizationBin(t))
	build := func(t *testing.T, writer string) PlanDocument {
		t.Helper()
		request.WriterVersion = writer
		outcome, err := BuildPlan(t.Context(), request)
		if err != nil || outcome.Plan == nil {
			t.Fatalf("writer %q: result=%+v error=%v", writer, outcome.Result, err)
		}
		return *outcome.Plan
	}
	assertRecord := func(t *testing.T, plan PlanDocument, want string) {
		t.Helper()
		if plan.SetupManifest.RoundfixVersion != want {
			t.Fatalf("record=%q want=%q", plan.SetupManifest.RoundfixVersion, want)
		}
		for _, postimage := range plan.Postimages {
			if postimage.Path != manifestPath {
				continue
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(postimage.Content, &fields); err != nil {
				t.Fatal(err)
			}
			raw, present := fields["roundfixVersion"]
			if want == "" {
				if present {
					t.Fatalf("unexpected roundfixVersion: %s", raw)
				}
			} else if string(raw) != `"`+want+`"` {
				t.Fatalf("postimage record=%s want=%q", raw, want)
			}
			return
		}
		t.Fatal("manifest postimage missing")
	}
	withoutWriter := build(t, "")
	assertRecord(t, withoutWriter, "")
	adoption := build(t, "0.140.0")
	assertRecord(t, adoption, "0.140.0")
	result, err := ApplyPlan(t.Context(), repository, adoption, adoption.PlanDigest)
	if err != nil || result.State != "verified" {
		t.Fatalf("apply=%+v error=%v", result, err)
	}
	request.Preservation.Mode = PreservationModeManagedRefresh
	for _, writer := range []string{"0.141.0", "", "0.0.0-dev", "0.141.0-rc.1"} {
		t.Run("converged writer "+writer, func(t *testing.T) {
			plan := build(t, writer)
			assertRecord(t, plan, "0.140.0")
			if len(plan.FileChanges) != 0 {
				t.Fatalf("converged plan proposes changes: %+v", plan.FileChanges)
			}
		})
	}
	guide := filepath.Join(repository, "docs/agents/agent-instructions.md")
	content, err := os.ReadFile(guide)
	if err != nil {
		t.Fatal(err)
	}
	stale := bytes.ReplaceAll(content, []byte("Fix root causes."), []byte("Fix root causes!"))
	if bytes.Equal(stale, content) {
		t.Fatal("fixture lacks root-cause guidance")
	}
	if err := os.WriteFile(guide, stale, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ writer, want string }{
		{"0.139.0", "0.140.0"}, {"0.140.0", "0.140.0"},
		{"0.141.0-rc.1", "0.140.0"}, {"", "0.140.0"},
		{"0.0.0-dev", "0.140.0"}, {"invalid", "0.140.0"},
		{"v0.141.0", "0.141.0"},
	} {
		t.Run("stale writer "+tc.writer, func(t *testing.T) {
			plan := build(t, tc.writer)
			assertRecord(t, plan, tc.want)
			if len(plan.FileChanges) == 0 {
				t.Fatal("stale guide produced no file change")
			}
		})
	}
}

func TestASetupManifestRoundfixVersionMustBeAReleaseVersion(t *testing.T) {
	t.Parallel()
	repository := newBaselinePlanCharacterizationRepository(t, true, true, true)
	request := baselinePlanCharacterizationRequest(repository, baselinePlanCharacterizationBin(t))
	request.WriterVersion = "0.140.0"
	outcome, err := BuildPlan(t.Context(), request)
	if err != nil || outcome.Plan == nil {
		t.Fatalf("plan=%+v error=%v", outcome.Result, err)
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"", "0.140.0", "0.140", "v0.140.0", "0.140.0-rc.1"} {
		t.Run(version, func(t *testing.T) {
			plan := *outcome.Plan
			plan.SetupManifest.RoundfixVersion = version
			data, err := json.Marshal(plan.SetupManifest)
			if err != nil {
				t.Fatal(err)
			}
			_, valid := parseManagedSetupManifest(data)
			wantValid := version == "" || version == "0.140.0"
			if valid != wantValid {
				t.Fatalf("parse version %q valid=%v want=%v", version, valid, wantValid)
			}
			err = validateSetupManifest(plan, catalog)
			if wantValid {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || err.Error() != "Setup Manifest roundfixVersion must be a release version such as 0.140.0" {
				t.Fatalf("invalid version %q error=%v", version, err)
			}
		})
	}
}
