package releaseplan

import (
	"context"
	"fmt"
)

// ReleaseGrant is the release authorization of one Spec a range delivers.
type ReleaseGrant struct {
	Spec      string
	CommitSHA string
	Record    string
	Revision  string
	Granted   bool
	Ceiling   Impact
	Paths     []string
	Detail    string
}

// ReleaseGrantSource is an optional GitSource capability for delivered records.
type ReleaseGrantSource interface {
	ReleaseGrants(context.Context, Range, []Commit) ([]ReleaseGrant, error)
}

// ReleaseAuthorization records the grants inspected and their approval answer.
type ReleaseAuthorization struct {
	Answered bool
	Reason   string
	Grants   []ReleaseGrant
}

// AnswerApproval applies the recorded Release Ceilings to an approval proposal.
func AnswerApproval(plan Plan, grants []ReleaseGrant) Plan {
	hasGrant := false
	for _, grant := range grants {
		hasGrant = hasGrant || grant.Granted
	}
	if !hasGrant {
		return plan
	}
	authorization := &ReleaseAuthorization{Grants: grants}
	plan.ReleaseAuthorization = authorization
	switch {
	case plan.Classification.Breaking:
		authorization.Reason = "the plan is breaking"
	case plan.Classification.ManualReason != "":
		authorization.Reason = "the plan was classified by hand"
	default:
		// Each condition examines all grants before the next condition begins.
		for _, grant := range grants {
			if !grant.Granted {
				authorization.Reason = fmt.Sprintf("%s grants no release ceiling", grant.Spec)
				return plan
			}
		}
		deliveries := make(map[string]bool, len(grants))
		for _, grant := range grants {
			if CompareImpact(grant.Ceiling, plan.Classification.Impact) < 0 {
				authorization.Reason = fmt.Sprintf("impact %s is above the %s ceiling of %s", plan.Classification.Impact, grant.Ceiling, grant.Spec)
				return plan
			}
			deliveries[grant.CommitSHA] = true
		}
		for _, change := range plan.Changes {
			if change.AutomaticImpact == plan.Classification.Impact && change.CrossesMaintenanceOnlyBoundary && !deliveries[change.CommitSHA] {
				sha := change.CommitSHA
				if len(sha) > 12 {
					sha = sha[:12]
				}
				authorization.Reason = fmt.Sprintf("commit %s sets the impact outside a delivered Spec", sha)
				return plan
			}
		}
		authorization.Answered = true
		plan.State = StateReady
		plan.Approval = Approval{Increment: plan.Approval.Increment, ProposedVersion: plan.ProposedVersion}
	}
	return plan
}
