// Suite: Baseline classification contracts
// Invariant: semantic dispositions stay digest-bound, and carrier warnings narrow only from exact managed or repository-extension evidence.
// Boundary IN: canonical analysis snapshots, strict proposal parsing, carrier bytes, Setup Manifest ownership, and retention input.
// Boundary OUT: ACPX process lifecycle, human interaction, and repository mutation.

package baseline

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"path/filepath"
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
		first.Destinations[1].Path != "docs/agents/specific-repository.md" {
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

func TestClassificationSnapshotMakesExactTextAvailableToSealedAnalysis(t *testing.T) {
	snapshot := classificationTestSnapshot(t)
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		t.Fatalf("marshal Analysis Snapshot: %v", err)
	}
	entry := snapshot.Entries[0]
	for _, required := range []string{
		`"semanticEntries"`,
		`"entryId":"` + entry.ID + `"`,
		`"text":"keep this rule\n"`,
	} {
		if !bytes.Contains(canonical, []byte(required)) {
			t.Fatalf("Analysis Snapshot does not expose %q:\n%s", required, canonical)
		}
	}
}

func TestClassificationProposalDerivesByteEvidenceLocally(t *testing.T) {
	snapshot := classificationTestSnapshot(t)
	entry := snapshot.Entries[0]
	destination := snapshot.Destinations[1]
	proposal := map[string]any{
		"schemaVersion":  ClassificationProposalSchemaVersion,
		"snapshotDigest": snapshot.SnapshotDigest,
		"dispositions": []map[string]any{{
			"entryId":        entry.ID,
			"classification": "normative-clause",
			"disposition":    "repository-rules",
			"destination": map[string]any{
				"documentType": destination.DocumentType,
				"path":         destination.Path,
			},
			"reason": "No active semantic guide owns this repository policy.",
		}},
	}
	payload, err := json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseClassificationProposal(payload, snapshot)
	if err != nil {
		t.Fatalf("classification without agent-computed byte evidence was rejected: %v", err)
	}
	got := parsed.Dispositions[0]
	if got.EntryDigest != entry.Digest ||
		got.Destination == nil ||
		got.Destination.Digest != entry.Digest ||
		got.Destination.ProposedBytes != base64.StdEncoding.EncodeToString(entry.SourceBytes) {
		t.Fatalf("classification byte evidence was not derived locally: %+v", got)
	}
}

func TestSemanticRuleDistributionAdmitsOnlyActiveSemanticOwners(t *testing.T) {
	repo := newPlanRepository(t)
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := ResolveProfile(repo, "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	modules, artifacts, err := resolveManagedArtifacts(
		catalog,
		profile,
		planTestDecisions(),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	artifactIDs := make([]string, len(artifacts))
	for index, artifact := range artifacts {
		artifactIDs[index] = artifact.ID
	}
	registry := catalog.SemanticOwnerRegistry(modules, artifactIDs)
	snapshot, err := NewAnalysisSnapshot(classificationTestSource(), registry)
	if err != nil {
		t.Fatal(err)
	}

	var cliDestination *ClassificationDestination
	for index := range snapshot.Destinations {
		destination := &snapshot.Destinations[index]
		if destination.Disposition == "repository-document" &&
			destination.Path == "docs/agents/cli.md" {
			cliDestination = destination
		}
		if destination.Path == "docs/agents/frontend.md" {
			t.Fatal("inactive frontend guide was advertised as a semantic destination")
		}
	}
	if cliDestination == nil {
		t.Fatal("active CLI guide was not advertised as a semantic destination")
	}

	proposal := classificationTestProposal(t, snapshot)
	entry := snapshot.Entries[0]
	proposal.Dispositions[0] = ReadoptionDisposition{
		EntryID:        entry.ID,
		EntryDigest:    entry.Digest,
		Classification: "normative-clause",
		Disposition:    "repository-document",
		Destination: &ReadoptionDestination{
			DocumentType: "agent-guide",
			Path:         cliDestination.Path,
			Digest:       entry.Digest,
		},
		Reason: "The active CLI guide owns this repository policy.",
	}
	encoded, err := json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseClassificationProposal(encoded, snapshot); err != nil {
		t.Fatalf("active semantic destination was rejected: %v", err)
	}

	proposal.Dispositions[0].Destination.Path = "docs/agents/frontend.md"
	encoded, err = json.Marshal(proposal)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseClassificationProposal(encoded, snapshot); err == nil {
		t.Fatal("inactive semantic destination was accepted")
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
				disposition.Destination.Path != "docs/agents/specific-repository.md" ||
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

func TestCarrierClassification(t *testing.T) {
	t.Parallel()

	t.Run("four kinds warn only when human review is needed", func(t *testing.T) {
		t.Parallel()

		repo := newInspectionRepository(t)
		catalog, err := LoadEmbeddedCatalog()
		if err != nil {
			t.Fatal(err)
		}
		profile, err := ResolveProfile(repo, "go-cli-tui", catalog)
		if err != nil {
			t.Fatal(err)
		}
		decisions := planTestDecisions()
		modules, artifacts, err := resolveManagedArtifacts(catalog, profile, decisions, false)
		if err != nil {
			t.Fatal(err)
		}
		guides := make([]plannedArtifact, 0, 2)
		for _, artifact := range artifacts {
			if artifact.Kind == "guide" {
				guides = append(guides, artifact)
				if len(guides) == 2 {
					break
				}
			}
		}
		if len(guides) != 2 {
			t.Fatalf("active guide artifacts = %d, want at least two", len(guides))
		}
		if guides[0].Path == guides[1].Path {
			t.Fatalf("selected guides share path %q, want distinct carriers", guides[0].Path)
		}

		current := "preserved repository bytes\n\n" + upsertManagedBlock("", guides[0])
		staleArtifact := guides[1]
		staleArtifact.Body += "\nstale managed clause"
		manifestBytes, err := marshalSetupManifestBytes(buildSetupManifest(
			catalog,
			profile,
			decisions,
			modules,
			artifacts,
			nil,
		))
		if err != nil {
			t.Fatal(err)
		}
		writeInspectionFile(t, repo, manifestPath, string(append(manifestBytes, '\n')))
		writeInspectionFile(t, repo, guides[0].Path, current)
		writeInspectionFile(t, repo, guides[1].Path, upsertManagedBlock("", staleArtifact))
		writeInspectionFile(t, repo, specificRepositoryPath, "Keep this recognized repository extension.\n")
		const unmanagedPath = "services/payments/AGENTS.md"
		writeInspectionFile(t, repo, unmanagedPath, "Keep this unmanaged nested policy.\n")
		commitInspectionRepository(t, repo, "seed classified carriers")

		inspection, err := InspectRepository(context.Background(), repo, nil)
		if err != nil {
			t.Fatal(err)
		}
		classifications, err := classifyCarriers(repo, inspection.Snapshot.Carriers, catalog, artifacts)
		if err != nil {
			t.Fatal(err)
		}
		got := make(map[string]carrierClassificationKind, len(classifications))
		for _, classification := range classifications {
			got[classification.Path] = classification.Kind
		}
		want := map[string]carrierClassificationKind{
			guides[0].Path:         carrierCurrentManaged,
			guides[1].Path:         carrierStaleManaged,
			specificRepositoryPath: carrierRepositoryExtension,
			unmanagedPath:          carrierUnmanagedNested,
		}
		for carrierPath, wantKind := range want {
			if got[carrierPath] != wantKind {
				t.Errorf("carrier %q classification = %q, want %q", carrierPath, got[carrierPath], wantKind)
			}
		}

		nonCarrier := Finding{Code: "baseline.inventory.other", Path: "other", Message: "unchanged"}
		warnings := append(append([]Finding(nil), inspection.Snapshot.Warnings...), nonCarrier)
		filtered := warningsForCarrierClassifications(warnings, classifications)
		if !hasRepositoryFinding(filtered, "baseline.inventory.nested-carrier-conflict", unmanagedPath) {
			t.Fatalf("unmanaged nested carrier warning missing: %+v", filtered)
		}
		for _, warning := range filtered {
			if warning.Code == "baseline.inventory.nested-carrier-conflict" {
				if warning.Path != unmanagedPath {
					t.Errorf("managed or recognized carrier still warns: %+v", warning)
				}
				if !strings.Contains(warning.Message, "inside setup-context-driven begin/end markers") ||
					!strings.Contains(warning.Message, "outside those markers are preserved") {
					t.Errorf("carrier warning does not name managed/preserved boundary: %q", warning.Message)
				}
			}
		}
		if filtered[len(filtered)-1] != nonCarrier {
			t.Fatalf("non-carrier diagnostic changed: %+v", filtered)
		}

		semanticOwners, err := ResolveSemanticOwnerRegistry(catalog, profile, planTestDecisions())
		if err != nil {
			t.Fatal(err)
		}
		preservation, err := planRootPreservationWithCatalog(
			inspection,
			RootPreservationRequest{
				Mode:             PreservationModeGreenfield,
				semanticOwners:   semanticOwners,
				managedArtifacts: artifacts,
				classifyCarriers: true,
				classifications:  classifications,
			},
			catalog,
		)
		if err != nil {
			t.Fatal(err)
		}
		if preservation.State != PreservationStateActionRequired {
			t.Fatalf("stale managed carrier preservation state = %q, want action_required", preservation.State)
		}
		staleIsRetentionInput := false
		for _, entry := range preservation.SourceBaseline.Entries {
			if entry.Path == guides[1].Path && entry.Kind == "managed-block" {
				staleIsRetentionInput = true
			}
			if entry.Path == guides[0].Path {
				t.Fatalf("current managed carrier became retention input: %+v", entry)
			}
		}
		if !staleIsRetentionInput {
			t.Fatalf("stale managed carrier is absent from retention input: %+v", preservation.SourceBaseline.Entries)
		}
	})

	t.Run("verified apply replans without managed carrier warnings", func(t *testing.T) {
		t.Parallel()

		repo := newBaselinePlanCharacterizationRepository(t, true, true, true)
		request := baselinePlanCharacterizationRequest(repo)
		initial, err := BuildPlan(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}
		if initial.Plan == nil {
			t.Fatalf("initial adoption returned result: %+v", initial.Result)
		}
		if _, err := ApplyPlan(
			context.Background(),
			repo,
			*initial.Plan,
			initial.Plan.PlanDigest,
		); err != nil {
			t.Fatal(err)
		}
		fresh, err := BuildPlan(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}
		if fresh.Plan == nil {
			t.Fatalf("verified re-plan returned result: %+v", fresh.Result)
		}
		if len(fresh.Plan.Warnings) != 0 || len(fresh.Result.Warnings) != 0 {
			t.Fatalf("verified re-plan warnings = plan:%+v result:%+v", fresh.Plan.Warnings, fresh.Result.Warnings)
		}
		if len(fresh.Plan.FileChanges) != 0 {
			t.Fatalf("verified re-plan file changes = %+v", fresh.Plan.FileChanges)
		}
	})
}

func TestUnclassifiableCarrierStillWarns(t *testing.T) {
	t.Parallel()

	t.Run("unknown managed marker retains boundary warning", func(t *testing.T) {
		t.Parallel()

		repo := newInspectionRepository(t)
		const carrierPath = "packages/api/AGENTS.md"
		writeInspectionFile(t, repo, carrierPath, `<!-- setup-context-driven:begin id=unknown.guide version=1 -->

unknown managed-looking bytes

<!-- setup-context-driven:end id=unknown.guide -->
`)
		commitInspectionRepository(t, repo, "seed unclassifiable carrier")

		catalog, err := LoadEmbeddedCatalog()
		if err != nil {
			t.Fatal(err)
		}
		profile, err := ResolveProfile(repo, "go-cli-tui", catalog)
		if err != nil {
			t.Fatal(err)
		}
		_, artifacts, err := resolveManagedArtifacts(catalog, profile, planTestDecisions(), false)
		if err != nil {
			t.Fatal(err)
		}
		inspection, err := InspectRepository(context.Background(), repo, nil)
		if err != nil {
			t.Fatal(err)
		}
		classifications, err := classifyCarriers(repo, inspection.Snapshot.Carriers, catalog, artifacts)
		if err != nil {
			t.Fatal(err)
		}
		for _, classification := range classifications {
			if classification.Path == carrierPath && classification.Kind != "" {
				t.Fatalf("unclassifiable carrier classification = %q", classification.Kind)
			}
		}
		filtered := warningsForCarrierClassifications(inspection.Snapshot.Warnings, classifications)
		if !hasRepositoryFinding(filtered, "baseline.inventory.nested-carrier-conflict", carrierPath) {
			t.Fatalf("unclassifiable carrier warning disappeared: %+v", filtered)
		}
		for _, warning := range filtered {
			if warning.Path == carrierPath &&
				(!strings.Contains(warning.Message, "inside setup-context-driven begin/end markers") ||
					!strings.Contains(warning.Message, "outside those markers are preserved")) {
				t.Fatalf("unclassifiable warning does not name byte boundary: %q", warning.Message)
			}
		}
	})
}

func TestCarrierClassificationReportsRootOpenFailure(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	missingRoot := filepath.Join(t.TempDir(), "missing")
	_, err = classifyCarriers(missingRoot, nil, catalog, nil)
	if !errors.Is(err, fs.ErrNotExist) ||
		!strings.Contains(err.Error(), "open repository root for carrier classification") {
		t.Fatalf("classifyCarriers() error = %v, want root-open failure", err)
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
		Decisions:  planTestDecisionsWithRepositoryExtension(),
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
