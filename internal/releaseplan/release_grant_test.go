package releaseplan

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// Suite: Release Ceiling approval answer.
// Invariant: every condition is conservative and retains the version and classification.
// Boundary IN: domain plans and normalized grants.
// Boundary OUT: Git, CLI rendering, and release mutation.
func TestAnswerApprovalFollowsEveryCeiling(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, reason string
		alter        func(*Plan, *[]ReleaseGrant)
	}{
		{name: "answered"},
		{name: "breaking", reason: "the plan is breaking", alter: func(p *Plan, g *[]ReleaseGrant) {
			p.Classification.Breaking = true
			p.Classification.ManualReason = "also manual"
		}},
		{name: "classified by hand", reason: "the plan was classified by hand", alter: func(p *Plan, g *[]ReleaseGrant) { p.Classification.ManualReason = "x"; (*g)[0].Ceiling = ImpactPatch }},
		{name: "every Spec grants", reason: "0002-y grants no release ceiling", alter: func(p *Plan, g *[]ReleaseGrant) {
			(*g)[0].Ceiling = ImpactPatch
			*g = append(*g, ReleaseGrant{Spec: "0002-y", CommitSHA: "second"})
		}},
		{name: "every ceiling", reason: "impact minor is above the patch ceiling of 0002-y", alter: func(p *Plan, g *[]ReleaseGrant) {
			*g = append(*g, ReleaseGrant{Spec: "0002-y", CommitSHA: "second", Granted: true, Ceiling: ImpactPatch})
		}},
		{name: "outside impact", reason: "commit 123456789abc sets the impact outside a delivered Spec", alter: func(p *Plan, g *[]ReleaseGrant) {
			p.Changes = append(p.Changes, ChangeEvidence{CommitSHA: "123456789abcdef", AutomaticImpact: ImpactMinor, CrossesMaintenanceOnlyBoundary: true})
		}},
		{name: "no ceiling", alter: func(p *Plan, g *[]ReleaseGrant) { (*g)[0].Granted = false; (*g)[0].Ceiling = ImpactNone }},
		{name: "no deliveries", alter: func(p *Plan, g *[]ReleaseGrant) { *g = nil }},
		{name: "lower impact outside", alter: func(p *Plan, g *[]ReleaseGrant) {
			p.Changes = append(p.Changes, ChangeEvidence{CommitSHA: "direct-fix", AutomaticImpact: ImpactPatch, CrossesMaintenanceOnlyBoundary: true})
		}},
		{name: "planning feature", alter: func(p *Plan, g *[]ReleaseGrant) {
			p.Changes = append(p.Changes, ChangeEvidence{CommitSHA: "plan", AutomaticImpact: ImpactMinor})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			plan := Plan{State: StateApprovalRequired, ProposedVersion: "v0.2.0", Classification: Classification{Impact: ImpactMinor, Source: SourceConventionalCommit}, Approval: Approval{Required: true, Increment: IncrementMinor, ProposedVersion: "v0.2.0", Question: "approve?"}, Changes: []ChangeEvidence{{CommitSHA: "delivery", AutomaticImpact: ImpactMinor, CrossesMaintenanceOnlyBoundary: true}}}
			grants := []ReleaseGrant{{Spec: "0001-x", CommitSHA: "delivery", Granted: true, Ceiling: ImpactMinor}}
			if tt.alter != nil {
				tt.alter(&plan, &grants)
			}
			got := AnswerApproval(plan, grants)
			if tt.name == "no ceiling" || tt.name == "no deliveries" {
				if !reflect.DeepEqual(got, plan) {
					t.Fatalf("no ceiling changed plan: %+v", got)
				}
				return
			}
			a := got.ReleaseAuthorization
			if a == nil || a.Reason != tt.reason || a.Answered != (tt.reason == "") || !reflect.DeepEqual(a.Grants, grants) {
				t.Fatalf("authorization=%+v", a)
			}
			if got.ProposedVersion != plan.ProposedVersion || !reflect.DeepEqual(got.Classification, plan.Classification) || !reflect.DeepEqual(got.Changes, plan.Changes) {
				t.Fatalf("answer changed evidence: %+v", got)
			}
			if tt.reason == "" {
				want := Approval{Increment: IncrementMinor, ProposedVersion: plan.ProposedVersion}
				if got.State != StateReady || got.Approval != want {
					t.Fatalf("answered plan=%+v", got)
				}
			} else if got.State != StateApprovalRequired || got.Approval != plan.Approval {
				t.Fatalf("refused plan=%+v", got)
			}
		})
	}
}

type ceilingBuildSource struct {
	commits []Commit
	calls   int
	err     error
}

func (s *ceilingBuildSource) ResolveRange(context.Context, string, string) (Range, error) {
	v, _ := ParseStableVersion("v0.1.0")
	return Range{Base: VersionRef{Version: v}}, nil
}
func (s *ceilingBuildSource) Commits(context.Context, Range) ([]Commit, error) { return s.commits, nil }
func (s *ceilingBuildSource) ReleaseGrants(context.Context, Range, []Commit) ([]ReleaseGrant, error) {
	s.calls++
	return []ReleaseGrant{{Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery"}}, s.err
}

func TestBuildLooksForGrantsOnlyWhenApprovalIsRequired(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, subject, path string
		calls               int
	}{
		{"patch", "fix: x", "internal/x.go", 0},
		{"maintenance", "docs: plan", "docs/backlog/2026-11-01-x.md", 0},
		{"manual classification", "refactor: x", "internal/x.go", 0},
		{"minor", "feat: deliver x", "internal/x.go", 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			source := &ceilingBuildSource{commits: []Commit{{SHA: "delivery", Subject: tt.subject, ChangedPaths: []string{tt.path}}}}
			plan, err := Build(context.Background(), Request{}, source)
			if err != nil || source.calls != tt.calls {
				t.Fatalf("Build=%+v err=%v calls=%d", plan, err, source.calls)
			}
			if tt.calls == 1 && (plan.State != StateReady || plan.ReleaseAuthorization == nil || !plan.ReleaseAuthorization.Answered) {
				t.Fatalf("Build did not answer: %+v", plan)
			}
		})
	}
	t.Run("source without the optional capability", func(t *testing.T) {
		t.Parallel()
		source := &ceilingBuildSource{commits: []Commit{{SHA: "delivery", Subject: "feat: x", ChangedPaths: []string{"internal/x.go"}}}}
		plan, err := Build(context.Background(), Request{}, struct{ GitSource }{source})
		if err != nil || source.calls != 0 || plan.State != StateApprovalRequired || plan.ReleaseAuthorization != nil {
			t.Fatalf("plan=%+v err=%v calls=%d", plan, err, source.calls)
		}
	})

	t.Run("read error", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("unreadable deliveries")
		source := &ceilingBuildSource{commits: []Commit{{Subject: "feat: x", ChangedPaths: []string{"internal/x.go"}}}, err: sentinel}
		if _, err := Build(context.Background(), Request{}, source); !errors.Is(err, sentinel) {
			t.Fatalf("Build error=%v", err)
		}
	})
}
