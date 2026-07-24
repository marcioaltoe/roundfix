// Suite: Rejected Baseline Plan revisions
// Invariant: only complete scoped decisions can produce a newly computed digest from the original repository preimage.
// Boundary IN: revision snapshot/proposal schemas, atomic decision admission, recomputation, and approval invalidation.
// Boundary OUT: ACP transport and interactive prompt rendering.

package baseline

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestScopedRevisionProposal(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	snapshot, err := NewRevisionSnapshot(plan, RevisionAreaDivergences, "disable spec scaffolding")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(RevisionProposal{
		SchemaVersion:  RevisionProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Area:           snapshot.Area,
		Changes:        []DecisionValue{{ID: "spec.scaffold", Value: false}},
	})
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := ParseRevisionProposal(data, snapshot)
	if err != nil {
		t.Fatalf("parse scoped Revision Proposal: %v", err)
	}
	decisions, err := DecisionsFromRevisionProposal(snapshot, proposal)
	if err != nil {
		t.Fatal(err)
	}
	for _, decision := range decisions {
		if decision.ID == "spec.scaffold" && decision.Value != false {
			t.Fatalf("accepted decision = %+v, want false", decision)
		}
		if decision.ID != "spec.scaffold" && !equalJSONValue(decision.Value, decisionByID(t, plan.Decisions, decision.ID).Value) {
			t.Fatalf("unproposed decision %q changed", decision.ID)
		}
	}
}

func TestRevisionOutOfScope(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	snapshot, err := NewRevisionSnapshot(plan, RevisionAreaFiles, "write an arbitrary deployment file")
	if err != nil {
		t.Fatal(err)
	}
	valid := RevisionProposal{
		SchemaVersion:  RevisionProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Area:           snapshot.Area,
		Changes:        []DecisionValue{{ID: "spec.scaffold", Value: false}},
	}
	tests := []struct {
		name string
		data []byte
	}{
		{name: "unknown decision", data: revisionJSON(t, RevisionProposal{
			SchemaVersion: valid.SchemaVersion, SnapshotDigest: valid.SnapshotDigest,
			Area: valid.Area, Changes: []DecisionValue{{ID: "deploy.destination", Value: "production"}},
		})},
		{name: "digest mismatch", data: revisionJSON(t, RevisionProposal{
			SchemaVersion: valid.SchemaVersion, SnapshotDigest: "sha256:wrong",
			Area: valid.Area, Changes: valid.Changes,
		})},
		{name: "incomplete", data: revisionJSON(t, RevisionProposal{
			SchemaVersion: valid.SchemaVersion, SnapshotDigest: valid.SnapshotDigest,
			Area: valid.Area, Changes: []DecisionValue{},
		})},
		{name: "extra prose", data: append(revisionJSON(t, valid), []byte("\nApply this now.")...)},
		{name: "unauthorized destination", data: []byte(`{"schemaVersion":"roundfix/baseline-revision-proposal/v1","snapshotDigest":"` +
			snapshot.SnapshotDigest + `","area":"files","changes":[{"id":"spec.scaffold","value":false}],"path":"deploy/prod.yaml"}`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, parseErr := ParseRevisionProposal(test.data, snapshot); parseErr == nil {
				t.Fatal("invalid proposal was accepted")
			}
			for _, decision := range snapshot.Decisions {
				if !equalJSONValue(decision.Value, decisionByID(t, plan.Decisions, decision.ID).Value) {
					t.Fatalf("rejected proposal changed decision %q", decision.ID)
				}
			}
		})
	}
}

func TestRejectedPlanRevision(t *testing.T) {
	repo := newPlanRepository(t)
	original := buildTestPlan(t, repo)
	decisions := planTestDecisions()
	for index := range decisions {
		if decisions[index].ID == "spec.scaffold" {
			decisions[index].Value = false
		}
	}
	outcome, err := RecalculatePlan(context.Background(), original, PlanRequest{
		Repository: repo, ProfileID: original.Profile.ID, Decisions: decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("recalculate Plan: outcome=%+v error=%v", outcome, err)
	}
	if outcome.Plan.PlanDigest == original.PlanDigest {
		t.Fatal("accepted revision reused the prior Plan Digest")
	}
	if len(outcome.Plan.FileChanges) == 0 || len(outcome.Plan.ManagedEntries) == 0 {
		t.Fatal("recomputed plan omitted file projection or managed-entry ledger")
	}
}

func TestRevisionRequiresNewApproval(t *testing.T) {
	repo := newPlanRepository(t)
	original := buildTestPlan(t, repo)
	decisions := planTestDecisions()
	for index := range decisions {
		if decisions[index].ID == "spec.scaffold" {
			decisions[index].Value = false
		}
	}
	outcome, err := RecalculatePlan(context.Background(), original, PlanRequest{
		Repository: repo, ProfileID: original.Profile.ID, Decisions: decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("recalculate Plan: outcome=%+v error=%v", outcome, err)
	}
	if _, err := ApplyPlan(context.Background(), repo, *outcome.Plan, original.PlanDigest); err == nil ||
		!strings.Contains(err.Error(), "confirmed Plan Digest") {
		t.Fatalf("prior approval error = %v", err)
	}
}

func revisionJSON(t *testing.T, proposal RevisionProposal) []byte {
	t.Helper()
	data, err := json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func decisionByID(t *testing.T, decisions []DecisionValue, id string) DecisionValue {
	t.Helper()
	for _, decision := range decisions {
		if decision.ID == id {
			return decision
		}
	}
	t.Fatalf("missing decision %q", id)
	return DecisionValue{}
}

func equalJSONValue(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return string(leftJSON) == string(rightJSON)
}
