package cli

import (
	"context"
	"encoding/json"
	"reflect"

	"fmt"
	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
	"strings"
	"testing"
)

// Suite: Release Ceiling at the public command boundary.
// Invariant: only sufficient, operative grants for every delivery answer approval.
// Boundary IN: committed temporary repositories and text/JSON command requests.
// Boundary OUT: release mutations, configuration, and network services.
func TestReleasePlanAnswersApprovalFromTheReleaseCeiling(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, ceiling, reason    string
		direct, breaking, manual bool
	}{
		{name: "ceiling minor", ceiling: "minor"},
		{name: "ceiling patch", ceiling: "patch", reason: "impact minor is above the patch ceiling of 0001-x"},
		{name: "no ceiling"},
		{name: "a direct feat commit", ceiling: "minor", direct: true},
		{name: "a breaking footer", ceiling: "major", breaking: true, reason: "the plan is breaking"},
		{name: "classified by hand", ceiling: "minor", manual: true, reason: "the plan was classified by hand"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := newReleasePlanCommandRepo(t, "v0.1.0")
			writeReleasePlanFile(t, repo, "docs/backlog/2026-11-01-example.md", "intent\n")
			gitReleasePlan(t, repo, "add", "-A")
			commitReleasePlan(t, repo, "docs: file intent")
			record := "docs/specs/0001-x/_authorization.md"
			authorization := "---\nstatus: approved\ngranted: 2026-11-01\naction: deliver x\nconsuming: 0001-x\npaths: []\noperations: [implement, release]\n"
			if tt.ceiling != "" {
				authorization += "release_ceiling: " + tt.ceiling + "\nrelease_paths: [VERSION, CHANGELOG.md]\n"
			}
			authorization += "---\n"
			writeReleasePlanFile(t, repo, record, authorization)
			writeReleasePlanFile(t, repo, "docs/specs/0001-x/_prd.md", "# X\n")
			gitReleasePlan(t, repo, "mv", "docs/backlog/2026-11-01-example.md", "docs/specs/0001-x/reference.md")
			gitReleasePlan(t, repo, "add", "-A")
			parent := commitReleasePlan(t, repo, "docs: plan x")
			gitReleasePlan(t, repo, "rm", "-r", "docs/specs/0001-x")
			writeReleasePlanFile(t, repo, "docs/history/specs/0001-x.md", "# X archived\n")
			writeReleasePlanFile(t, repo, "internal/x.go", "package x\n")
			gitReleasePlan(t, repo, "add", "-A")
			body := ""
			if tt.breaking {
				body = "BREAKING CHANGE: remove the old contract"
			}
			delivery := commitReleasePlan(t, repo, "feat: deliver 0001-x (#2)", body)
			if tt.direct {
				writeReleasePlanFile(t, repo, "internal/direct.go", "package x\n")
				gitReleasePlan(t, repo, "add", "-A")
				sha := commitReleasePlan(t, repo, "feat: direct change")
				tt.reason = "commit " + sha[:12] + " sets the impact outside a delivered Spec"
			}
			var args []string
			if tt.manual {
				args = []string{"--impact", "minor", "--reason", "x"}
			}
			wantExit, state := exitUnverified, "approval_required"
			if tt.name == "ceiling minor" {
				wantExit, state = exitOK, "ready"
			}
			outputs := make(map[string]string)
			for _, format := range []string{"text", "json"} {
				code, output, stderr := runReleasePlanCommandInRepo(t, repo, append(args, "--format", format)...)
				if code != wantExit || stderr != "" {
					t.Fatalf("%s exit=%d stderr=%q output=%s", format, code, stderr, output)
				}
				outputs[format] = output
				if format == "text" {
					if !strings.HasPrefix(output, "Decision: "+state+"\n") {
						t.Fatalf("decision: %s", output)
					}
					if tt.ceiling == "" {
						if strings.Contains(output, "Release authorization:") {
							t.Fatalf("unexpected authorization: %s", output)
						}
					} else {
						line := "Release authorization: answered"
						if tt.reason != "" {
							line = "Release authorization: not answered — " + tt.reason
						}
						grant := fmt.Sprintf("- 0001-x ceiling=%s record=%s@%s paths=VERSION,CHANGELOG.md", tt.ceiling, record, parent[:12])
						if !strings.Contains(output, line+"\n") || !strings.Contains(output, grant+"\n") {
							t.Fatalf("authorization missing: %s", output)
						}
						if strings.Index(output, "Release authorization:") < strings.Index(output, "Next action:") {
							t.Fatalf("authorization precedes next action: %s", output)
						}
					}
					if !strings.Contains(output, "Proposed version: v0.2.0\n") {
						t.Fatalf("proposal: %s", output)
					}
				} else {
					var plan struct {
						State         string `json:"state"`
						Authorization *struct {
							Answered bool   `json:"answered"`
							Reason   string `json:"reason"`
							Grants   []struct {
								Spec     string   `json:"spec"`
								SHA      string   `json:"commitSHA"`
								Record   string   `json:"record"`
								Revision string   `json:"revision"`
								Granted  bool     `json:"granted"`
								Ceiling  string   `json:"ceiling"`
								Paths    []string `json:"paths"`
							} `json:"grants"`
						} `json:"releaseAuthorization"`
					}
					if err := json.Unmarshal([]byte(output), &plan); err != nil {
						t.Fatal(err)
					}
					if plan.State != state {
						t.Fatalf("state = %s", plan.State)
					}
					if tt.ceiling == "" {
						if strings.Contains(output, `"releaseAuthorization"`) {
							t.Fatalf("unexpected JSON authorization: %s", output)
						}
					} else {
						a := plan.Authorization
						if a == nil || a.Answered != (wantExit == exitOK) || a.Reason != tt.reason || len(a.Grants) != 1 {
							t.Fatalf("authorization: %s", output)
						}
						g := a.Grants[0]
						if g.Spec != "0001-x" || g.SHA != delivery || g.Record != record || g.Revision != parent || !g.Granted || g.Ceiling != tt.ceiling || strings.Join(g.Paths, ",") != "VERSION,CHANGELOG.md" {
							t.Fatalf("grant: %+v", g)
						}
					}
				}
			}
			if tt.name == "no ceiling" {
				// Replace only the historical plan object so the range retains identical
				// commit IDs. This permits a literal byte comparison with no release grant.
				gitReleasePlan(t, repo, "checkout", "--detach", parent)
				writeReleasePlanFile(t, repo, record, strings.Replace(authorization, "[implement, release]", "[implement]", 1))
				gitReleasePlan(t, repo, "add", "-A")
				gitReleasePlan(t, repo, "commit", "--amend", "--no-edit")
				replacement := strings.TrimSpace(gitReleasePlanOutput(t, repo, "rev-parse", "HEAD"))
				gitReleasePlan(t, repo, "checkout", "--detach", delivery)
				gitReleasePlan(t, repo, "replace", parent, replacement)
				for _, format := range []string{"text", "json"} {
					code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format", format)
					if code != exitUnverified || stderr != "" || output != outputs[format] {
						t.Fatalf("no ceiling changed %s bytes: exit=%d stderr=%q\nwith release:%s\nwithout release:%s", format, code, stderr, outputs[format], output)
					}
				}
			}
		})
	}
}

func TestReleaseGrantsReadEveryDeliveryAtItsParent(t *testing.T) {
	t.Parallel()
	repo := newReleasePlanCommandRepo(t, "v0.1.0")
	cases := []struct {
		slug, operations, extra, detail string
		granted                         bool
	}{
		{"0001-good", "[release]", "release_ceiling: minor\nrelease_paths: [VERSION]\n", "", true},
		{"0002-no-release", "[implement]", "", "record grants no release", false},
		{"0003-no-ceiling", "[release]", "", "record grants release without a ceiling", false},
		{"0004-malformed", "[release]", "release_ceiling: unknown\nrelease_paths: [VERSION]\n", "record refused: malformed_record", false},
		{"0005-withdrawn", "[release]", "", "record refused: status", false},
	}
	for _, tt := range cases {
		status := "approved"
		if tt.slug == "0005-withdrawn" {
			status = "withdrawn"
		}
		record := fmt.Sprintf("---\nstatus: %s\ngranted: 2026-11-01\naction: deliver\nconsuming: %s\npaths: []\noperations: %s\n%s---\n", status, tt.slug, tt.operations, tt.extra)
		writeReleasePlanFile(t, repo, "custom/specs/"+tt.slug+"/_authorization.md", record)
	}
	writeReleasePlanFile(t, repo, "custom/specs/unpaired/_authorization.md", "not an authorization")
	gitReleasePlan(t, repo, "add", "-A")
	parent := commitReleasePlan(t, repo, "docs: plan deliveries")
	gitReleasePlan(t, repo, "tag", "v0.1.1")
	gitReleasePlan(t, repo, "rm", "-r", "custom/specs")
	for _, tt := range cases {
		writeReleasePlanFile(t, repo, "custom/archive/"+tt.slug+".md", "not parsed; source_revision need not be reachable\n")
	}
	writeReleasePlanFile(t, repo, "custom/archive/archive-only.md", "not a delivery\n")
	writeReleasePlanFile(t, repo, "internal/x.go", "package x\n")
	gitReleasePlan(t, repo, "add", "-A")
	sha := commitReleasePlan(t, repo, "feat: deliver several Specs")
	source := newReleasePlanGitSource(repo, preflight.ExecGitRunner{})
	before := snapshotReleasePlanRepo(t, repo)
	grants, err := source.ReleaseGrants(context.Background(), releaseplan.Range{}, []releaseplan.Commit{{SHA: sha}})
	if err != nil {
		t.Fatal(err)
	}
	assertReleasePlanRepoUnchanged(t, repo, before)
	if len(grants) != len(cases) {
		t.Fatalf("grants=%+v", grants)
	}
	for i, g := range grants {
		tt := cases[i]
		ceiling := releaseplan.ImpactNone
		var paths []string
		if tt.granted {
			ceiling = releaseplan.ImpactMinor
			paths = []string{"VERSION"}
		}
		// Refused records retain parsed paths as inspectable evidence.
		if tt.slug == "0004-malformed" {
			paths = []string{"VERSION"}
		}
		if g.Spec != tt.slug || g.CommitSHA != sha || g.Record != "custom/specs/"+tt.slug+"/_authorization.md" || g.Revision != parent || g.Granted != tt.granted || g.Detail != tt.detail || g.Ceiling != ceiling || !reflect.DeepEqual(g.Paths, paths) {
			t.Fatalf("grant=%+v want=%+v", g, tt)
		}
	}
	for _, format := range []string{"text", "json"} {
		code, output, stderr := runReleasePlanCommandInRepo(t, repo, "--format", format)
		if code != exitUnverified || stderr != "" {
			t.Fatalf("%s exit=%d stderr=%q output=%s", format, code, stderr, output)
		}
		reason := "0002-no-release grants no release ceiling"
		if format == "text" {
			if !strings.Contains(output, "Release authorization: not answered — "+reason) || !strings.Contains(output, "- 0002-no-release ceiling=none record=custom/specs/0002-no-release/_authorization.md@"+parent[:12]+" paths=none") {
				t.Fatal(output)
			}
		} else {
			plan := decodeReleasePlanJSON(t, output)
			a := plan.ReleaseAuthorization
			if a == nil || a.Answered || a.Reason != reason || len(a.Grants) != len(cases) || a.Grants[1].Ceiling != "" || a.Grants[1].Detail != "record grants no release" || a.Grants[1].Paths == nil {
				t.Fatal(output)
			}
		}
	}
}
