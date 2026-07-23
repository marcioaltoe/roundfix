// Suite: Baseline root-instruction classification contracts
// Invariant: only a complete digest-bound classification of every sealed source entry can become deterministic planning input.
// Boundary IN: canonical analysis snapshots, strict proposal parsing, manual fallback, and Decision Document normalization.
// Boundary OUT: ACPX process lifecycle, human interaction, and repository mutation.

package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestSealedClassificationSnapshotIsCanonicalAndBounded(t *testing.T) {
	source := classificationTestSource()
	first, err := NewAnalysisSnapshot(source)
	if err != nil {
		t.Fatalf("build first Analysis Snapshot: %v", err)
	}
	second, err := NewAnalysisSnapshot(source)
	if err != nil {
		t.Fatalf("build second Analysis Snapshot: %v", err)
	}
	firstBytes, err := first.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal first Analysis Snapshot: %v", err)
	}
	secondBytes, err := second.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal second Analysis Snapshot: %v", err)
	}
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("equivalent Analysis Snapshots produced different canonical bytes")
	}
	if first.SnapshotDigest == "" || first.SnapshotDigest != second.SnapshotDigest {
		t.Fatalf("snapshot digests differ: %q != %q", first.SnapshotDigest, second.SnapshotDigest)
	}
	if len(first.Entries) != len(source.Entries) ||
		len(first.Destinations) != 2 ||
		first.Destinations[0].Disposition != "rejected" ||
		first.Destinations[1].Path != "docs/agents/repository-rules.md" {
		t.Fatalf("unexpected sealed snapshot contract: %+v", first)
	}

	oversized := source
	oversized.Entries = make([]ReadoptionSourceEntry, AnalysisSnapshotMaxEntries+1)
	for index := range oversized.Entries {
		oversized.Entries[index] = source.Entries[0]
		oversized.Entries[index].ID += strings.Repeat("x", index+1)
	}
	oversized.EntryCount = len(oversized.Entries)
	if _, err := NewAnalysisSnapshot(oversized); err == nil ||
		!strings.Contains(err.Error(), "256 entries") {
		t.Fatalf("oversized snapshot error = %v", err)
	}
}

func TestProposalValidationRejectsIncompleteOrUntrustedOutput(t *testing.T) {
	snapshot := classificationTestSnapshot(t)
	valid := classificationTestProposal(t, snapshot)
	validJSON, err := json.Marshal(valid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseClassificationProposal(validJSON, snapshot); err != nil {
		t.Fatalf("valid proposal rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(ClassificationProposal) []byte
	}{
		{
			name: "extra prose",
			mutate: func(proposal ClassificationProposal) []byte {
				payload, _ := json.Marshal(proposal)
				return append([]byte("proposal:\n"), payload...)
			},
		},
		{
			name: "unknown field",
			mutate: func(proposal ClassificationProposal) []byte {
				payload, _ := json.Marshal(proposal)
				return bytes.Replace(payload, []byte(`"schemaVersion":`),
					[]byte(`"unknown":true,"schemaVersion":`), 1)
			},
		},
		{
			name: "duplicate key",
			mutate: func(proposal ClassificationProposal) []byte {
				payload, _ := json.Marshal(proposal)
				return bytes.Replace(payload, []byte(`"snapshotDigest":`),
					[]byte(`"snapshotDigest":"sha256:`+strings.Repeat("0", 64)+`","snapshotDigest":`), 1)
			},
		},
		{
			name: "digest mismatch",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.SnapshotDigest = "sha256:" + strings.Repeat("0", 64)
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
		{
			name: "missing disposition",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.Dispositions = proposal.Dispositions[:0]
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
		{
			name: "unknown entry",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.Dispositions[0].EntryID = "source-entry." + strings.Repeat("f", 64)
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
		{
			name: "duplicate disposition",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.Dispositions = append(proposal.Dispositions, proposal.Dispositions[0])
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
		{
			name: "unsupported destination",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.Dispositions[0].Disposition = "repository-document"
				proposal.Dispositions[0].Destination = &ReadoptionDestination{
					DocumentType: "agent-guide",
					Path:         "docs/agents/guide.md",
					Digest:       proposal.Dispositions[0].EntryDigest,
				}
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
		{
			name: "changed proposed bytes",
			mutate: func(proposal ClassificationProposal) []byte {
				proposal.Dispositions[0].Destination.ProposedBytes = "Y2hhbmdlZA=="
				payload, _ := json.Marshal(proposal)
				return payload
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ParseClassificationProposal(test.mutate(valid), snapshot); err == nil {
				t.Fatal("invalid proposal was accepted")
			}
		})
	}
}

func TestManualClassificationFallbackReturnsCompleteDestinations(t *testing.T) {
	snapshot := classificationTestSnapshot(t)
	proposal, err := ManualClassificationProposal(snapshot)
	if err != nil {
		t.Fatalf("build manual Classification Proposal: %v", err)
	}
	if len(proposal.Dispositions) != len(snapshot.Entries) {
		t.Fatalf("manual dispositions = %d, want %d", len(proposal.Dispositions), len(snapshot.Entries))
	}
	for index, disposition := range proposal.Dispositions {
		if disposition.EntryID != snapshot.Entries[index].ID {
			t.Fatalf("manual disposition %d targets %q, want %q",
				index, disposition.EntryID, snapshot.Entries[index].ID)
		}
		switch disposition.Disposition {
		case "rejected":
			if disposition.Destination != nil || strings.TrimSpace(disposition.Reason) == "" {
				t.Fatalf("incomplete rejected destination: %+v", disposition)
			}
		case "repository-rules":
			if disposition.Destination == nil ||
				disposition.Destination.Path != "docs/agents/repository-rules.md" ||
				disposition.Destination.ProposedBytes == "" {
				t.Fatalf("incomplete Repository-Specific Normative Rules destination: %+v", disposition)
			}
		default:
			t.Fatalf("unsupported manual disposition: %+v", disposition)
		}
	}
}

func TestSealedClassificationEquivalentProposalsProduceSamePlanDigest(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "keep this repository rule\n")
	commitInspectionRepository(t, repo, "add root instructions")

	unresolved, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModePreservation},
	)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := NewAnalysisSnapshot(unresolved.SourceBaseline)
	if err != nil {
		t.Fatal(err)
	}
	firstProposal, err := ManualClassificationProposal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(firstProposal)
	if err != nil {
		t.Fatal(err)
	}
	secondProposal, err := ParseClassificationProposal(encoded, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	firstDecisions, err := DecisionDocumentFromClassificationProposal(snapshot, firstProposal)
	if err != nil {
		t.Fatal(err)
	}
	secondDecisions, err := DecisionDocumentFromClassificationProposal(snapshot, secondProposal)
	if err != nil {
		t.Fatal(err)
	}

	first := classificationTestPlan(t, repo, firstDecisions)
	second := classificationTestPlan(t, repo, secondDecisions)
	if first.PlanDigest != second.PlanDigest {
		t.Fatalf("equivalent proposal sources changed Plan Digest: %s != %s",
			first.PlanDigest, second.PlanDigest)
	}
}

func classificationTestSource() ReadoptionSourceBaseline {
	content := []byte("keep this rule\n")
	entry := newReadoptionSourceEntry(
		"AGENTS.md",
		"unmarked-span",
		0,
		len(content),
		content,
		"e06863f68ee0f70c25a3e1427568856be1350ed4b067945d9b38570d4c3c097f",
		map[string]any{"markerState": "unmarked"},
	)
	return ReadoptionSourceBaseline{
		ID:               "baseline.readoption." + strings.Repeat("b", 64),
		DeclaredIdentity: "unconfigured",
		Compatibility:    "incompatible",
		Digest:           strings.Repeat("b", 64),
		CarrierCount:     1,
		EntryCount:       1,
		ByteCount:        len(content),
		Entries:          []ReadoptionSourceEntry{entry},
	}
}

func classificationTestSnapshot(t *testing.T) AnalysisSnapshot {
	t.Helper()
	snapshot, err := NewAnalysisSnapshot(classificationTestSource())
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func classificationTestProposal(t *testing.T, snapshot AnalysisSnapshot) ClassificationProposal {
	t.Helper()
	proposal, err := ManualClassificationProposal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	return proposal
}

func classificationTestPlan(t *testing.T, repo string, decisions DecisionDocument) PlanDocument {
	t.Helper()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode:      PreservationModePreservation,
			Decisions: &decisions,
		},
	})
	if err != nil {
		t.Fatalf("build classified plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("classification did not produce a complete plan: %+v", outcome.Result)
	}
	return *outcome.Plan
}
