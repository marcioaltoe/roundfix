package releaseplan

import (
	"errors"
	"testing"
)

// Suite: Release Plan domain
// Invariant: classified release impact deterministically produces the documented state and version.
// Boundary IN: stable version parsing, impact ordering, proposal calculation, domain schema values.
// Boundary OUT: Git history, changed-path classification, CLI rendering, repository mutations.

func TestParseStableVersion(t *testing.T) {
	valid := []struct {
		name  string
		tag   string
		major int
		minor int
		patch int
	}{
		{name: "version zero release", tag: "v0.4.0", major: 0, minor: 4, patch: 0},
		{name: "major release", tag: "v1.4.2", major: 1, minor: 4, patch: 2},
		{name: "multi digit release", tag: "v10.20.30", major: 10, minor: 20, patch: 30},
	}

	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStableVersion(tt.tag)
			if err != nil {
				t.Fatalf("ParseStableVersion(%q): %v", tt.tag, err)
			}
			if got.Major() != tt.major || got.Minor() != tt.minor || got.Patch() != tt.patch {
				t.Fatalf("ParseStableVersion(%q) = %d.%d.%d, want %d.%d.%d", tt.tag, got.Major(), got.Minor(), got.Patch(), tt.major, tt.minor, tt.patch)
			}
			if got.String() != tt.tag {
				t.Fatalf("ParseStableVersion(%q).String() = %q, want round-trip without normalization drift", tt.tag, got.String())
			}
		})
	}

	invalid := []struct {
		name string
		tag  string
		want error
	}{
		{name: "empty", tag: "", want: ErrMalformedStableVersion},
		{name: "missing v prefix", tag: "1.2.3", want: ErrMalformedStableVersion},
		{name: "missing patch", tag: "v1.2", want: ErrMalformedStableVersion},
		{name: "extra component", tag: "v1.2.3.4", want: ErrMalformedStableVersion},
		{name: "major leading zero", tag: "v01.2.3", want: ErrMalformedStableVersion},
		{name: "minor leading zero", tag: "v1.02.3", want: ErrMalformedStableVersion},
		{name: "patch leading zero", tag: "v1.2.03", want: ErrMalformedStableVersion},
		{name: "build metadata", tag: "v1.2.3+build", want: ErrMalformedStableVersion},
		{name: "pre-release", tag: "v1.2.3-rc.1", want: ErrPrereleaseVersion},
		{name: "version zero pre-release", tag: "v0.4.0-alpha", want: ErrPrereleaseVersion},
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseStableVersion(tt.tag)
			if err == nil {
				t.Fatalf("ParseStableVersion(%q) succeeded, want error", tt.tag)
			}
			if !errors.Is(err, tt.want) {
				t.Fatalf("ParseStableVersion(%q) error = %v, want errors.Is(..., %v)", tt.tag, err, tt.want)
			}
			var versionErr StableVersionError
			if !errors.As(err, &versionErr) {
				t.Fatalf("ParseStableVersion(%q) error = %T, want StableVersionError", tt.tag, err)
			}
			if versionErr.Input != tt.tag {
				t.Fatalf("StableVersionError.Input = %q, want %q", versionErr.Input, tt.tag)
			}
			if versionErr.NextAction == "" {
				t.Fatalf("StableVersionError.NextAction is empty")
			}
		})
	}
}

func TestCalculateProposal(t *testing.T) {
	tests := []struct {
		name                 string
		base                 string
		input                ProposalInput
		wantState            State
		wantProposedVersion  string
		wantIncrement        IncrementKind
		wantBreaking         bool
		wantApprovalRequired bool
	}{
		{
			name:                 "patch impact is ready",
			base:                 "v0.4.0",
			input:                ProposalInput{Impact: ImpactPatch},
			wantState:            StateReady,
			wantProposedVersion:  "v0.4.1",
			wantIncrement:        IncrementPatch,
			wantApprovalRequired: false,
		},
		{
			name:                 "compatible minor requires approval",
			base:                 "v0.4.1",
			input:                ProposalInput{Impact: ImpactMinor},
			wantState:            StateApprovalRequired,
			wantProposedVersion:  "v0.5.0",
			wantIncrement:        IncrementMinor,
			wantApprovalRequired: true,
		},
		{
			name:                 "version-zero breaking maps to minor",
			base:                 "v0.4.1",
			input:                ProposalInput{Impact: ImpactMajor, Breaking: true},
			wantState:            StateApprovalRequired,
			wantProposedVersion:  "v0.5.0",
			wantIncrement:        IncrementMinor,
			wantBreaking:         true,
			wantApprovalRequired: true,
		},
		{
			name:                 "breaking marker raises patch to version-zero minor",
			base:                 "v0.4.1",
			input:                ProposalInput{Impact: ImpactPatch, Breaking: true},
			wantState:            StateApprovalRequired,
			wantProposedVersion:  "v0.5.0",
			wantIncrement:        IncrementMinor,
			wantBreaking:         true,
			wantApprovalRequired: true,
		},
		{
			name:                 "major breaking increments major at one or later",
			base:                 "v1.4.2",
			input:                ProposalInput{Impact: ImpactMajor, Breaking: true},
			wantState:            StateApprovalRequired,
			wantProposedVersion:  "v2.0.0",
			wantIncrement:        IncrementMajor,
			wantBreaking:         true,
			wantApprovalRequired: true,
		},
		{
			name:                "none produces no release",
			base:                "v1.4.2",
			input:               ProposalInput{Impact: ImpactNone},
			wantState:           StateNoRelease,
			wantProposedVersion: "",
			wantIncrement:       IncrementNone,
		},
		{
			name:                "manual classification required produces no version",
			base:                "v1.4.2",
			input:               ProposalInput{ManualClassificationRequired: true},
			wantState:           StateManualClassificationRequired,
			wantProposedVersion: "",
			wantIncrement:       IncrementNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := mustParseVersion(t, tt.base)
			got, err := CalculateProposal(base, tt.input)
			if err != nil {
				t.Fatalf("CalculateProposal: %v", err)
			}
			if got.State != tt.wantState {
				t.Fatalf("State = %q, want %q", got.State, tt.wantState)
			}
			if got.ProposedVersion != tt.wantProposedVersion {
				t.Fatalf("ProposedVersion = %q, want %q", got.ProposedVersion, tt.wantProposedVersion)
			}
			if got.Increment != tt.wantIncrement {
				t.Fatalf("Increment = %q, want %q", got.Increment, tt.wantIncrement)
			}
			if got.Breaking != tt.wantBreaking {
				t.Fatalf("Breaking = %v, want %v", got.Breaking, tt.wantBreaking)
			}
			if got.Approval.Required != tt.wantApprovalRequired {
				t.Fatalf("Approval.Required = %v, want %v", got.Approval.Required, tt.wantApprovalRequired)
			}
		})
	}
}

func TestCalculateProposalRejectsUnknownImpact(t *testing.T) {
	base := mustParseVersion(t, "v1.4.2")
	_, err := CalculateProposal(base, ProposalInput{Impact: Impact("feature")})
	if err == nil {
		t.Fatal("CalculateProposal succeeded, want error")
	}
	var impactErr UnknownImpactError
	if !errors.As(err, &impactErr) {
		t.Fatalf("CalculateProposal error = %T, want UnknownImpactError", err)
	}
}

func TestImpactOrdering(t *testing.T) {
	ordered := []Impact{ImpactNone, ImpactPatch, ImpactMinor, ImpactMajor}
	for leftIndex, left := range ordered {
		for rightIndex, right := range ordered {
			t.Run(string(left)+" vs "+string(right), func(t *testing.T) {
				got := CompareImpact(left, right)
				switch {
				case leftIndex < rightIndex && got >= 0:
					t.Fatalf("CompareImpact(%q, %q) = %d, want negative", left, right, got)
				case leftIndex == rightIndex && got != 0:
					t.Fatalf("CompareImpact(%q, %q) = %d, want zero", left, right, got)
				case leftIndex > rightIndex && got <= 0:
					t.Fatalf("CompareImpact(%q, %q) = %d, want positive", left, right, got)
				}
			})
		}
	}

	if got := MaxImpact(ImpactPatch, ImpactNone, ImpactMajor, ImpactMinor); got != ImpactMajor {
		t.Fatalf("MaxImpact mixed order = %q, want %q", got, ImpactMajor)
	}
	if got := MaxImpact(ImpactNone, ImpactPatch); got != ImpactPatch {
		t.Fatalf("MaxImpact none and patch = %q, want %q", got, ImpactPatch)
	}
}

func TestApprovalDecision(t *testing.T) {
	tests := []struct {
		name                string
		base                string
		input               ProposalInput
		wantRequired        bool
		wantIncrement       IncrementKind
		wantProposedVersion string
		wantQuestion        string
		wantBreaking        bool
	}{
		{
			name:                "minor approval question",
			base:                "v0.4.1",
			input:               ProposalInput{Impact: ImpactMinor},
			wantRequired:        true,
			wantIncrement:       IncrementMinor,
			wantProposedVersion: "v0.5.0",
			wantQuestion:        "Approve the minor increment to v0.5.0?",
		},
		{
			name:                "version-zero breaking uses minor approval question",
			base:                "v0.4.1",
			input:               ProposalInput{Impact: ImpactMajor, Breaking: true},
			wantRequired:        true,
			wantIncrement:       IncrementMinor,
			wantProposedVersion: "v0.5.0",
			wantQuestion:        "Approve the minor increment to v0.5.0?",
			wantBreaking:        true,
		},
		{
			name:                "major approval question",
			base:                "v1.4.2",
			input:               ProposalInput{Impact: ImpactMajor, Breaking: true},
			wantRequired:        true,
			wantIncrement:       IncrementMajor,
			wantProposedVersion: "v2.0.0",
			wantQuestion:        "Approve the major increment to v2.0.0?",
			wantBreaking:        true,
		},
		{
			name:                "patch needs no approval",
			base:                "v1.4.2",
			input:               ProposalInput{Impact: ImpactPatch},
			wantRequired:        false,
			wantIncrement:       "",
			wantProposedVersion: "",
			wantQuestion:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := mustParseVersion(t, tt.base)
			proposal, err := CalculateProposal(base, tt.input)
			if err != nil {
				t.Fatalf("CalculateProposal: %v", err)
			}
			if proposal.Approval.Required != tt.wantRequired {
				t.Fatalf("Approval.Required = %v, want %v", proposal.Approval.Required, tt.wantRequired)
			}
			if proposal.Approval.Increment != tt.wantIncrement {
				t.Fatalf("Approval.Increment = %q, want %q", proposal.Approval.Increment, tt.wantIncrement)
			}
			if proposal.Approval.ProposedVersion != tt.wantProposedVersion {
				t.Fatalf("Approval.ProposedVersion = %q, want %q", proposal.Approval.ProposedVersion, tt.wantProposedVersion)
			}
			if proposal.Approval.Question != tt.wantQuestion {
				t.Fatalf("Approval.Question = %q, want %q", proposal.Approval.Question, tt.wantQuestion)
			}
			if proposal.Breaking != tt.wantBreaking {
				t.Fatalf("Breaking = %v, want %v", proposal.Breaking, tt.wantBreaking)
			}
		})
	}
}

func TestPlanSchemaVersionAndEvidence(t *testing.T) {
	baseVersion := mustParseVersion(t, "v0.4.0")
	proposal, err := CalculateProposal(baseVersion, ProposalInput{Impact: ImpactPatch})
	if err != nil {
		t.Fatalf("CalculateProposal: %v", err)
	}
	changes := []ChangeEvidence{
		{
			CommitSHA:                      "abc123",
			Subject:                        "fix: repair release artifact check",
			ConventionalType:               "fix",
			AutomaticImpact:                ImpactPatch,
			CrossesMaintenanceOnlyBoundary: true,
		},
	}

	plan := NewPlan(
		VersionRef{Tag: "v0.4.0", Version: baseVersion},
		RevisionRef{Name: "HEAD", CommitSHA: "def456"},
		Classification{Source: SourceConventionalCommit, Impact: ImpactPatch},
		proposal,
		changes,
	)
	changes[0].CommitSHA = "mutated"

	if plan.SchemaVersion != SchemaVersion {
		t.Fatalf("SchemaVersion = %q, want %q", plan.SchemaVersion, SchemaVersion)
	}
	if plan.SchemaVersion != "roundfix.release-plan/v1" {
		t.Fatalf("schema identifier = %q, want roundfix.release-plan/v1", plan.SchemaVersion)
	}
	if len(plan.Changes) != 1 || plan.Changes[0].CommitSHA != "abc123" {
		t.Fatalf("Changes = %+v, want copied evidence with original commit SHA", plan.Changes)
	}
	if plan.ProposedVersion != "v0.4.1" || plan.State != StateReady {
		t.Fatalf("plan decision = state %q version %q, want ready v0.4.1", plan.State, plan.ProposedVersion)
	}
}

func mustParseVersion(t *testing.T, tag string) Version {
	t.Helper()
	version, err := ParseStableVersion(tag)
	if err != nil {
		t.Fatalf("ParseStableVersion(%q): %v", tag, err)
	}
	return version
}
