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
	grants  []ReleaseGrant
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
	if s.grants != nil {
		return s.grants, s.err
	}
	return []ReleaseGrant{{Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"CHANGELOG.md", "VERSION"}}}, s.err
}

func TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease(t *testing.T) {
	t.Parallel()
	delivery := Commit{SHA: "delivery", Subject: "feat: deliver x", ChangedPaths: []string{"internal/x.go"}}
	for _, tt := range []struct {
		name     string
		commits  []Commit
		grants   []ReleaseGrant
		request  Request
		reason   string
		calls    int
		state    State
		answered bool
	}{
		{
			name: "patch", calls: 0, state: StateReady,
			commits: []Commit{{SHA: "delivery", Subject: "fix: x", ChangedPaths: []string{"internal/x.go"}}},
		},
		{
			name: "maintenance", calls: 0, state: StateNoRelease,
			commits: []Commit{{SHA: "plan", Subject: "docs: plan", ChangedPaths: []string{"docs/backlog/2026-11-01-x.md"}}},
		},
		{
			name: "manual_classification", calls: 1, state: StateManualClassificationRequired,
			commits: []Commit{{SHA: "refactor", Subject: "refactor: x", ChangedPaths: []string{"internal/x.go"}}},
		},
		{
			name: "minor", calls: 1, state: StateReady, answered: true,
			commits: []Commit{delivery},
		},
		{
			name: "candidate_in_a_minor_range", calls: 1, state: StateManualClassificationRequired,
			commits: []Commit{delivery, {SHA: "candidate", Subject: "chore: adjust x", ChangedPaths: []string{"internal/x.go"}}},
		},
		{
			name: "preparation_is_answered", calls: 1, state: StateReady, answered: true,
			commits: []Commit{delivery, {SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}}},
		},
		{
			name: "paths_from_every_granted_record", calls: 1, state: StateReady, answered: true,
			commits: []Commit{delivery, {SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}}},
			grants: []ReleaseGrant{
				{Spec: "0001-x", Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"CHANGELOG.md"}},
				{Spec: "0002-y", Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"VERSION"}},
			},
		},
		{
			name: "ungranted_paths_do_not_prepare", calls: 1, state: StateManualClassificationRequired,
			commits: []Commit{delivery, {SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}}},
			grants: []ReleaseGrant{
				{Spec: "0001-x", Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"CHANGELOG.md"}},
				{Spec: "0002-y", CommitSHA: "delivery", Paths: []string{"VERSION"}},
			},
		},
		{
			name: "empty_early_grants_are_not_reloaded", calls: 1, state: StateApprovalRequired,
			commits: []Commit{delivery, {SHA: "candidate", Subject: "chore: adjust x", ChangedPaths: []string{"internal/x.go"}}},
			grants:  []ReleaseGrant{},
			request: Request{ManualImpact: ImpactMinor, ManualReason: "compatible behavior"},
		},
		{
			name: "manual_preparation_is_not_answered", calls: 1, state: StateApprovalRequired,
			commits: []Commit{delivery, {SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}}},
			request: Request{ManualImpact: ImpactMinor, ManualReason: "compatible behavior"},
			reason:  "the plan was classified by hand",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			source := &ceilingBuildSource{commits: tt.commits, grants: tt.grants}
			plan, err := Build(context.Background(), tt.request, source)
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if source.calls != tt.calls {
				t.Fatalf("ReleaseGrants calls = %d, want %d", source.calls, tt.calls)
			}
			if plan.State != tt.state {
				t.Fatalf("State = %q, want %q", plan.State, tt.state)
			}
			if tt.answered || tt.reason != "" {
				grants := tt.grants
				if grants == nil {
					grants = []ReleaseGrant{{Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"CHANGELOG.md", "VERSION"}}}
				}
				want := &ReleaseAuthorization{Answered: tt.answered, Reason: tt.reason, Grants: grants}
				if !reflect.DeepEqual(plan.ReleaseAuthorization, want) || plan.Approval.Required != !tt.answered {
					t.Fatalf("authorization = %+v, approval = %+v, want %+v", plan.ReleaseAuthorization, plan.Approval, want)
				}
			} else if plan.ReleaseAuthorization != nil {
				t.Fatalf("unexpected authorization = %+v", plan.ReleaseAuthorization)
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
		plan, err := Build(context.Background(), Request{}, source)
		if !errors.Is(err, sentinel) || source.calls != 1 || !reflect.DeepEqual(plan, Plan{}) {
			t.Fatalf("Build=%+v error=%v calls=%d", plan, err, source.calls)
		}
	})
	t.Run("candidate read error", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("unreadable candidate grants")
		source := &ceilingBuildSource{commits: []Commit{{Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md"}}}, err: sentinel}
		plan, err := Build(context.Background(), Request{}, source)
		if !errors.Is(err, sentinel) || err.Error() != "load release grants: unreadable candidate grants" || source.calls != 1 || !reflect.DeepEqual(plan, Plan{}) {
			t.Fatalf("Build=%+v error=%v calls=%d", plan, err, source.calls)
		}
	})
}

// Without a preparation candidate, grant discovery still follows the proposal's
// approval boundary, even when the request supplies a manual impact.
func TestBuildLooksForGrantsOnlyWhenApprovalIsRequired(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		request Request
		state   State
		version string
		calls   int
	}{
		{name: "automatic_patch", state: StateReady, version: "v0.1.1"},
		{
			name:    "manual_patch",
			request: Request{ManualImpact: ImpactPatch, ManualReason: "compatible fix"},
			state:   StateReady, version: "v0.1.1",
		},
		{
			name:    "manual_minor",
			request: Request{ManualImpact: ImpactMinor, ManualReason: "compatible addition"},
			state:   StateApprovalRequired, version: "v0.2.0", calls: 1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			source := &ceilingBuildSource{commits: []Commit{{SHA: "delivery", Subject: "fix: repair x", ChangedPaths: []string{"internal/x.go"}}}}
			plan, err := Build(context.Background(), tt.request, source)
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if source.calls != tt.calls || plan.State != tt.state || plan.ProposedVersion != tt.version {
				t.Fatalf("calls=%d state=%s version=%s, want calls=%d state=%s version=%s", source.calls, plan.State, plan.ProposedVersion, tt.calls, tt.state, tt.version)
			}
			var wantAuthorization *ReleaseAuthorization
			if tt.calls == 1 {
				wantAuthorization = &ReleaseAuthorization{
					Reason: "the plan was classified by hand",
					Grants: []ReleaseGrant{{Granted: true, Ceiling: ImpactMinor, CommitSHA: "delivery", Paths: []string{"CHANGELOG.md", "VERSION"}}},
				}
			}
			if !reflect.DeepEqual(plan.ReleaseAuthorization, wantAuthorization) {
				t.Fatalf("authorization=%+v, want %+v", plan.ReleaseAuthorization, wantAuthorization)
			}
			wantChanges := []ChangeEvidence{{CommitSHA: "delivery", Subject: "fix: repair x", ConventionalType: "fix", AutomaticImpact: ImpactPatch, CrossesMaintenanceOnlyBoundary: true}}
			if !reflect.DeepEqual(plan.Changes, wantChanges) {
				t.Fatalf("changes=%+v, want %+v", plan.Changes, wantChanges)
			}
		})
	}
}
