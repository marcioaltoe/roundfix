package releaseplan

import "fmt"

// ProposalInput is the classified release impact to calculate from. Breaking
// raises the effective impact to major while preserving the version-zero rule.
type ProposalInput struct {
	Impact                       Impact
	Breaking                     bool
	ManualClassificationRequired bool
}

// Proposal is the deterministic version calculation result.
type Proposal struct {
	State           State
	ProposedVersion string
	Increment       IncrementKind
	Breaking        bool
	Approval        Approval
}

// CompareImpact orders impacts by semantic-version precedence.
func CompareImpact(left Impact, right Impact) int {
	leftRank := impactRank(left)
	rightRank := impactRank(right)
	switch {
	case leftRank < rightRank:
		return -1
	case leftRank > rightRank:
		return 1
	default:
		return 0
	}
}

// MaxImpact returns the highest semantic-version impact in inputs.
func MaxImpact(impacts ...Impact) Impact {
	maximum := ImpactNone
	for _, impact := range impacts {
		if CompareImpact(impact, maximum) > 0 {
			maximum = impact
		}
	}
	return maximum
}

// CalculateProposal applies the Release Plan version calculation contract to a
// stable base version and classified impact.
func CalculateProposal(base Version, input ProposalInput) (Proposal, error) {
	if input.ManualClassificationRequired {
		return Proposal{
			State:     StateManualClassificationRequired,
			Increment: IncrementNone,
			Approval:  noApproval(),
		}, nil
	}
	if !AllowedImpact(input.Impact) {
		return Proposal{}, UnknownImpactError{Impact: input.Impact}
	}

	impact := input.Impact
	breaking := input.Breaking || impact == ImpactMajor
	if input.Breaking {
		impact = MaxImpact(impact, ImpactMajor)
	}

	switch impact {
	case ImpactNone:
		return Proposal{
			State:     StateNoRelease,
			Increment: IncrementNone,
			Breaking:  breaking,
			Approval:  noApproval(),
		}, nil
	case ImpactPatch:
		version := base.IncrementPatch()
		return Proposal{
			State:           StateReady,
			ProposedVersion: version.String(),
			Increment:       IncrementPatch,
			Breaking:        breaking,
			Approval:        noApproval(),
		}, nil
	case ImpactMinor:
		version := base.IncrementMinor()
		return approvalProposal(version, IncrementMinor, breaking), nil
	case ImpactMajor:
		if base.Major() == 0 {
			version := base.IncrementMinor()
			return approvalProposal(version, IncrementMinor, true), nil
		}
		version := base.IncrementMajor()
		return approvalProposal(version, IncrementMajor, true), nil
	default:
		return Proposal{}, UnknownImpactError{Impact: input.Impact}
	}
}

func noApproval() Approval {
	return Approval{Increment: IncrementNone}
}

func approvalProposal(version Version, increment IncrementKind, breaking bool) Proposal {
	proposedVersion := version.String()
	approval := Approval{
		Required:        true,
		Increment:       increment,
		ProposedVersion: proposedVersion,
		Question:        fmt.Sprintf("Approve the %s increment to %s?", increment, proposedVersion),
	}
	return Proposal{
		State:           StateApprovalRequired,
		ProposedVersion: proposedVersion,
		Increment:       increment,
		Breaking:        breaking,
		Approval:        approval,
	}
}

// AllowedImpact reports whether impact is part of the stable schema.
func AllowedImpact(impact Impact) bool {
	switch impact {
	case ImpactNone, ImpactPatch, ImpactMinor, ImpactMajor:
		return true
	default:
		return false
	}
}

func impactRank(impact Impact) int {
	switch impact {
	case ImpactNone:
		return 0
	case ImpactPatch:
		return 1
	case ImpactMinor:
		return 2
	case ImpactMajor:
		return 3
	default:
		return -1
	}
}
