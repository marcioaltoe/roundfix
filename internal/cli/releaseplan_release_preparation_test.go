package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/releaseplan"
)

// Suite: planning references and Release-Preparation Commits at the command boundary.
// Invariant: maintenance preserves the proposal without widening a Release Ceiling.
// Boundary IN: committed fixture repositories and this repository's release history.
// Boundary OUT: network services and release mutations.
func TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name          string
		omitReference bool
		preparation   bool
		unlistedPath  bool
		noCeiling     bool
		fix           bool
		manual        bool
	}{
		{name: "planning_reference_is_answered"},
		{name: "without_the_reference_the_same_decision", omitReference: true},
		{name: "preparation_keeps_the_decision", preparation: true},
		{name: "preparation_with_an_unlisted_path", preparation: true, unlistedPath: true},
		{name: "record_without_a_ceiling", preparation: true, noCeiling: true},
		{name: "fix_on_release_paths", fix: true},
		{name: "classified_by_hand", manual: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newReleasePlanCommandRepo(t, "v0.1.0")
			record := "---\nstatus: approved\ngranted: 2026-11-01\naction: deliver x\nconsuming: 0001-x\npaths: []\noperations: [implement, release]\n"
			if !tt.noCeiling {
				record += "release_ceiling: minor\nrelease_paths: [VERSION, CHANGELOG.md]\n"
			}
			writeReleasePlanFile(t, repo, "docs/specs/0001-x/_authorization.md", record+"---\n")
			writeReleasePlanFile(t, repo, "docs/specs/0001-x/_prd.md", "# X\n")
			if !tt.omitReference {
				writeReleasePlanFile(t, repo, "docs/references/evidence.md", "outside evidence\n")
			}
			gitReleasePlan(t, repo, "add", "-A")
			planning := commitReleasePlan(t, repo, "docs: plan x")
			gitReleasePlan(t, repo, "rm", "-r", "docs/specs/0001-x")
			writeReleasePlanFile(t, repo, "docs/history/specs/0001-x.md", "# X archived\n")
			writeReleasePlanFile(t, repo, "internal/x.go", "package x\n")
			gitReleasePlan(t, repo, "add", "-A")
			delivery := commitReleasePlan(t, repo, "feat: deliver 0001-x (#2)")
			var extra string
			if tt.preparation {
				writeReleasePlanFile(t, repo, "VERSION", "0.2.0\n")
				writeReleasePlanFile(t, repo, "CHANGELOG.md", "# v0.2.0\n")
				if tt.unlistedPath {
					writeReleasePlanFile(t, repo, "README.md", "unlisted change\n")
				}
				gitReleasePlan(t, repo, "add", "-A")
				extra = commitReleasePlan(t, repo, "chore: release v0.2.0")
			}
			if tt.fix {
				// A second base lets text expose the patch's own determining line too.
				gitReleasePlan(t, repo, "tag", "v0.1.1")
				writeReleasePlanFile(t, repo, "CHANGELOG.md", "corrected release notes\n")
				gitReleasePlan(t, repo, "add", "-A")
				extra = commitReleasePlan(t, repo, "fix: correct release notes")
			}

			wantState, wantVersion, wantExit := releaseplan.StateReady, "v0.2.0", exitOK
			wantAuthorization := "Release authorization: answered\n"
			if tt.unlistedPath || tt.noCeiling {
				wantState, wantVersion, wantExit = releaseplan.StateManualClassificationRequired, "", exitUnverified
				wantAuthorization = ""
			}
			args := []string{"--from", "v0.1.0"}
			if tt.manual {
				args = append(args, "--impact", "minor", "--reason", "x")
				wantState, wantExit = releaseplan.StateApprovalRequired, exitUnverified
				wantAuthorization = "Release authorization: not answered — the plan was classified by hand\n"
			}
			before := snapshotReleasePlanRepo(t, repo)
			defer assertReleasePlanRepoUnchanged(t, repo, before)
			for _, format := range []string{"text", "json"} {
				code, output, stderr := runReleasePlanCommandInRepo(t, repo, append(args, "--format", format)...)
				if code != wantExit || stderr != "" {
					t.Fatalf("%s exit=%d stderr=%q output=%s", format, code, stderr, output)
				}
				if format == "text" {
					if !strings.HasPrefix(output, "Decision: "+string(wantState)+"\n") {
						t.Fatalf("decision: %s", output)
					}
					version := wantVersion
					if version == "" {
						version = "none"
					}
					if !strings.Contains(output, "Proposed version: "+version+"\n") ||
						(wantAuthorization != "" && !strings.Contains(output, wantAuthorization)) {
						t.Fatalf("proposal or authorization: %s", output)
					}
					if wantState == releaseplan.StateManualClassificationRequired {
						want := "Blocking commits:\n- " + extra[:12] + " chore: release v0.2.0 [impact=none, type=chore, manual]\n"
						if !strings.HasSuffix(output, want) {
							t.Fatalf("blocking commit: %s", output)
						}
					} else {
						want := "Determining commits:\n- " + delivery[:12] + " feat: deliver 0001-x (#2) [impact=minor, type=feat]\n"
						if strings.Contains(output, "Blocking commits:") || !strings.HasSuffix(output, want) {
							t.Fatalf("determining commits: %s", output)
						}
					}
					continue
				}

				plan := decodeReleasePlanJSON(t, output)
				if plan.State != wantState || plan.ProposedVersion != wantVersion {
					t.Fatalf("proposal: %s", output)
				}
				if wantAuthorization != "" {
					a := plan.ReleaseAuthorization
					if a == nil || a.Answered == tt.manual || len(a.Grants) != 1 ||
						(tt.manual && a.Reason != "the plan was classified by hand") || (!tt.manual && a.Reason != "") {
						t.Fatalf("authorization: %s", output)
					}
				} else if plan.ReleaseAuthorization != nil {
					t.Fatalf("unexpected authorization: %s", output)
				}
				if wantState == releaseplan.StateManualClassificationRequired {
					if !plan.Classification.ManualClassificationRequired || strings.Join(plan.Classification.BlockingCommits, ",") != extra {
						t.Fatalf("blocking commits: %s", output)
					}
				} else if plan.Classification.Impact != releaseplan.ImpactMinor || len(plan.Classification.BlockingCommits) != 0 {
					t.Fatalf("classification: %s", output)
				}
				// Raw keys distinguish omission from an explicitly encoded false value.
				var raw struct {
					Changes []map[string]json.RawMessage `json:"changes"`
				}
				if err := json.Unmarshal([]byte(output), &raw); err != nil {
					t.Fatal(err)
				}
				wantCount := 2
				if extra != "" {
					wantCount++
				}
				if len(plan.Changes) != wantCount || len(raw.Changes) != wantCount || plan.Changes[0].CommitSHA != planning || plan.Changes[1].CommitSHA != delivery ||
					(extra != "" && plan.Changes[2].CommitSHA != extra) {
					t.Fatalf("changes: %s", output)
				}
				for i, change := range plan.Changes {
					marker, present := raw.Changes[i]["releasePreparation"]
					wantPreparation := tt.preparation && !tt.unlistedPath && !tt.noCeiling && change.CommitSHA == extra
					if present != wantPreparation || (present && string(marker) != "true") || (wantPreparation && change.CrossesMaintenanceOnlyBoundary) {
						t.Fatalf("preparation evidence for %s: %s", change.CommitSHA, output)
					}
					if tt.fix && change.CommitSHA == extra && (change.AutomaticImpact != releaseplan.ImpactPatch || change.ConventionalType != "fix" || !change.CrossesMaintenanceOnlyBoundary) {
						t.Fatalf("fix lost its impact: %s", output)
					}
				}
			}
			if tt.fix {
				for _, format := range []string{"text", "json"} {
					code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--from", "v0.1.1", "--format", format)
					if code != exitOK || stderr != "" {
						t.Fatalf("patch-only %s exit=%d stderr=%q output=%s", format, code, stderr, output)
					}
					if format == "text" {
						if !strings.HasPrefix(output, "Decision: ready\n") || !strings.Contains(output, "Impact: patch\n") || !strings.Contains(output, "Proposed version: v0.1.2\n") ||
							!strings.HasSuffix(output, "- "+extra[:12]+" fix: correct release notes [impact=patch, type=fix]\n") {
							t.Fatalf("patch-only text: %s", output)
						}
					} else {
						plan := decodeReleasePlanJSON(t, output)
						if plan.State != releaseplan.StateReady || plan.Classification.Impact != releaseplan.ImpactPatch || plan.ProposedVersion != "v0.1.2" {
							t.Fatalf("patch-only JSON: %s", output)
						}
					}
				}
			}
		})
	}
}

func TestTheV0143ReleaseRangeIsReadyWithItsPreparationCommit(t *testing.T) {
	t.Parallel()
	root := strings.TrimSpace(gitReleasePlanOutput(t, ".", "rev-parse", "--show-toplevel"))
	repo := filepath.Join(t.TempDir(), "history")
	gitReleasePlan(t, root, "clone", "--shared", "--no-checkout", root, repo)
	gittest.Harden(t, repo)
	// A checkout that fetched no tags still
	// holds the commits: name them with the two release tags when they are absent.
	for _, release := range [][2]string{{"v0.143.3", "251ed81e81937353af9cdf4727599257c3900e8f"}, {"v0.143.4", "4f8a5b7db8385f7b3745f63d6cb5668dd576f662"}} {
		if gitReleasePlanOutput(t, repo, "tag", "--list", release[0]) == "" {
			gitReleasePlan(t, repo, "tag", release[0], release[1])
		}
	}
	gitReleasePlan(t, repo, "checkout", "--detach", "v0.143.4")
	// This checkout includes directory symlinks, which the file-only fixture
	// snapshot cannot read. Git checks tracked, untracked and ignored paths here.
	snapshot := func() [5]string {
		return [5]string{
			gitReleasePlanOutput(t, repo, "rev-parse", "HEAD"),
			gitReleasePlanOutput(t, repo, "for-each-ref", "--format=%(refname) %(objectname)"),
			gitReleasePlanOutput(t, repo, "remote", "-v"),
			gitReleasePlanOutput(t, repo, "config", "--local", "--list"),
			gitReleasePlanOutput(t, repo, "--no-optional-locks", "status", "--porcelain=v1", "--ignored", "-z"),
		}
	}
	before := snapshot()
	if before[4] != "" {
		t.Fatalf("historical clone is not clean: %q", before[4])
	}
	defer func() {
		if after := snapshot(); after != before {
			t.Fatalf("historical clone changed: before=%q after=%q", before, after)
		}
	}()

	setCommandWorkDirForTest(t, repo)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"release", "plan", "--from", "v0.143.3", "--to", "v0.143.4"}, &stdout, &stderr)
	output := stdout.String()
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%s", code, stderr.String(), output)
	}
	wantPrefix := "Decision: ready\n" +
		"Base: v0.143.3 (251ed81e8193)\n" +
		"Target: v0.143.4 (4f8a5b7db838)\n" +
		"Impact: patch\n" +
		"Breaking: false\n" +
		"Proposed version: v0.143.4\n" +
		"Approval required: no\n" +
		"Next action: release may proceed for v0.143.4 after independent release verification.\n"
	wantSuffix := "Determining commits:\n" +
		"- 3bc58f41a3ad fix: deliver 0391-the-grouping-bound-becomes-a-criterion (#872) [impact=patch, type=fix]\n"
	if !strings.HasPrefix(output, wantPrefix) || !strings.HasSuffix(output, wantSuffix) {
		t.Fatalf("Surface Transcript 1 mismatch:\n%s", output)
	}
}
