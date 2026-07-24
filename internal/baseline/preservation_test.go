// Suite: Baseline root-instruction preservation
// Invariant: every trusted root source is backed up and every preserved source entry is dispositioned before planning is ready.
// Boundary IN: bounded repository inspection, root backups, Decision Documents, Readoption, Source Baselines, and retention contracts.
// Boundary OUT: profile alignment, portable Plan Documents, file transactions, and ACP classification proposals.

package baseline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGreenfieldPlanBacksUpWithoutImport(t *testing.T) {
	repo := newInspectionRepository(t)
	const instructions = "keep repository policy\n"
	writeInspectionFile(t, repo, "AGENTS.md", instructions)
	commitInspectionRepository(t, repo, "seed")

	inspection := inspectPreservationRepository(t, repo)
	plan, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode: PreservationModeGreenfield,
	})
	if err != nil {
		t.Fatalf("plan greenfield preservation: %v", err)
	}
	if plan.State != PreservationStateReady {
		t.Fatalf("greenfield state = %q, want %q: %+v", plan.State, PreservationStateReady, plan.Findings)
	}
	if len(plan.Backups) != 1 {
		t.Fatalf("greenfield backups = %+v, want one", plan.Backups)
	}
	sum := sha256.Sum256([]byte(instructions))
	wantPath := "AGENTS." + hex.EncodeToString(sum[:]) + ".md"
	if plan.Backups[0].Path != wantPath ||
		plan.Backups[0].ContentIdentity != "sha256:"+hex.EncodeToString(sum[:]) {
		t.Fatalf("greenfield backup = %+v, want path %q with raw-byte identity", plan.Backups[0], wantPath)
	}
	if len(plan.Dispositions) != 0 || len(plan.RepositoryRulesBytes) != 0 {
		t.Fatalf("greenfield imported root rules: dispositions=%+v bytes=%q", plan.Dispositions, plan.RepositoryRulesBytes)
	}
}

func TestPreservationRequiresEveryDisposition(t *testing.T) {
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "first rule\n\nsecond rule\n")
	commitInspectionRepository(t, repo, "seed")
	inspection := inspectPreservationRepository(t, repo)

	unresolved, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode: PreservationModePreservation,
	})
	if err != nil {
		t.Fatalf("plan unresolved preservation: %v", err)
	}
	if unresolved.State != PreservationStateActionRequired ||
		unresolved.DecisionSkeleton == nil ||
		len(unresolved.SourceBaseline.Entries) == 0 {
		t.Fatalf("unresolved preservation plan = %+v", unresolved)
	}

	incomplete := unresolved.DecisionSkeleton.Document
	incomplete.Readoption.Dispositions = incomplete.Readoption.Dispositions[:0]
	stillUnresolved, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode:      PreservationModePreservation,
		Decisions: &incomplete,
	})
	if err != nil {
		t.Fatalf("plan incomplete preservation: %v", err)
	}
	if stillUnresolved.State != PreservationStateActionRequired ||
		!hasPreservationFinding(stillUnresolved.Findings, "baseline.preservation.disposition.missing") {
		t.Fatalf("incomplete classifications became ready: %+v", stillUnresolved)
	}
}

func TestPreservationPlanAcceptsCompleteDecisionDocument(t *testing.T) {
	repo := newInspectionRepository(t)
	const instructions = "preserve this rule\n"
	writeInspectionFile(t, repo, "CLAUDE.md", instructions)
	commitInspectionRepository(t, repo, "seed")
	inspection := inspectPreservationRepository(t, repo)

	unresolved, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode: PreservationModePreservation,
	})
	if err != nil {
		t.Fatalf("plan unresolved preservation: %v", err)
	}
	decisions := unresolved.DecisionSkeleton.Document
	ready, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode:      PreservationModePreservation,
		Decisions: &decisions,
	})
	if err != nil {
		t.Fatalf("plan complete preservation: %v", err)
	}
	if ready.State != PreservationStateReady {
		t.Fatalf("complete preservation state = %q: %+v", ready.State, ready.Findings)
	}
	if len(ready.Dispositions) != len(ready.SourceBaseline.Entries) {
		t.Fatalf("dispositions = %d, entries = %d", len(ready.Dispositions), len(ready.SourceBaseline.Entries))
	}
	if string(ready.RepositoryRulesBytes) != instructions {
		t.Fatalf("Repository-Specific Normative Rules bytes = %q, want %q", ready.RepositoryRulesBytes, instructions)
	}
}

func TestRootBackupIdentityRejectsCollisions(t *testing.T) {
	repo := newInspectionRepository(t)
	const instructions = "root instructions\n"
	writeInspectionFile(t, repo, "AGENTS.md", instructions)
	sum := sha256.Sum256([]byte(instructions))
	backupPath := "AGENTS." + hex.EncodeToString(sum[:]) + ".md"
	writeInspectionFile(t, repo, backupPath, "different bytes\n")
	commitInspectionRepository(t, repo, "seed")

	plan, err := PlanRootPreservation(inspectPreservationRepository(t, repo), RootPreservationRequest{
		Mode: PreservationModeGreenfield,
	})
	if err != nil {
		t.Fatalf("plan colliding backup: %v", err)
	}
	if plan.State != PreservationStateBlocked ||
		!hasPreservationFinding(plan.Findings, "baseline.preservation.backup.collision") {
		t.Fatalf("backup collision did not block: %+v", plan)
	}
}

func TestRootBackupIdentitySafeAliasesBackUpTargetOnce(t *testing.T) {
	repo := newInspectionRepository(t)
	const instructions = "shared root policy\n"
	writeInspectionFile(t, repo, "policy/shared.md", instructions)
	if err := os.Symlink("policy/shared.md", filepath.Join(repo, "AGENTS.md")); err != nil {
		t.Fatalf("create AGENTS alias: %v", err)
	}
	if err := os.Symlink("policy/shared.md", filepath.Join(repo, "CLAUDE.md")); err != nil {
		t.Fatalf("create CLAUDE alias: %v", err)
	}
	commitInspectionRepository(t, repo, "seed")

	plan, err := PlanRootPreservation(inspectPreservationRepository(t, repo), RootPreservationRequest{
		Mode: PreservationModeGreenfield,
	})
	if err != nil {
		t.Fatalf("plan safe aliases: %v", err)
	}
	if plan.State != PreservationStateReady || len(plan.Backups) != 1 {
		t.Fatalf("safe alias plan = %+v, want one ready backup", plan)
	}
	if plan.Backups[0].CarrierPath != "AGENTS.md" ||
		plan.Backups[0].SourcePath != "policy/shared.md" {
		t.Fatalf("safe alias backup = %+v", plan.Backups[0])
	}
}

func TestDecisionDocumentSkeletonPassesStrictParser(t *testing.T) {
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", strings.Join([]string{
		"classify me",
		"<!-- setup-context-driven:begin id=root.old version=2 -->",
		"old managed rule",
		"<!-- setup-context-driven:end id=root.old -->",
		"",
	}, "\n"))
	commitInspectionRepository(t, repo, "seed")
	inspection := inspectPreservationRepository(t, repo)

	plan, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode: PreservationModePreservation,
	})
	if err != nil {
		t.Fatalf("plan preservation skeleton: %v", err)
	}
	data, err := json.Marshal(plan.DecisionSkeleton.Document)
	if err != nil {
		t.Fatalf("marshal decision skeleton: %v", err)
	}
	parsed, err := ParseDecisionDocument(data, "generated-decision-skeleton")
	if err != nil {
		t.Fatalf("strict parser rejected generated skeleton: %v\n%s", err, data)
	}
	if !reflect.DeepEqual(parsed, plan.DecisionSkeleton.Document) {
		t.Fatalf("parsed skeleton differs:\n got=%+v\nwant=%+v", parsed, plan.DecisionSkeleton.Document)
	}
	if !strings.Contains(plan.DecisionSkeleton.NextAction, "review") {
		t.Fatalf("decision skeleton next action = %q, want stable review action", plan.DecisionSkeleton.NextAction)
	}
}

func TestDecisionDocumentSkeletonRejectsMalformedInput(t *testing.T) {
	_, err := ParseDecisionDocument([]byte(`{
	  "schemaVersion":"setup-context-driven/decisions/0.0.1",
	  "version":"0.0.1",
	  "decisions":[],
	  "decisions":[]
	}`), "duplicate.json")
	var decisionErr *DecisionDocumentError
	if err == nil || !strings.Contains(err.Error(), "decision-file.json.duplicate-key") ||
		!errors.As(err, &decisionErr) {
		t.Fatalf("malformed Decision Document error = %v", err)
	}
}

func TestReadoptionCompatibilityMaintainedFixture(t *testing.T) {
	fixturePath := filepath.Join(
		"testdata", "parity-corpus", "v1", "fixtures", "readoption-preservation.json",
	)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read maintained Readoption fixture: %v", err)
	}
	var fixture struct {
		Input struct {
			DecisionDocument json.RawMessage `json:"decisionDocument"`
		} `json:"input"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatalf("decode maintained Readoption fixture: %v", err)
	}
	document, err := ParseDecisionDocument(fixture.Input.DecisionDocument, fixturePath)
	if err != nil {
		t.Fatalf("parse maintained Readoption Decision Document: %v", err)
	}
	if document.Readoption == nil || len(document.Readoption.Dispositions) != 19 {
		t.Fatalf("maintained Readoption dispositions = %+v, want 19", document.Readoption)
	}
	for _, disposition := range document.Readoption.Dispositions {
		if disposition.Disposition == "repository-rules" &&
			disposition.Destination.Path != specificRepositoryPath {
			t.Fatalf(
				"legacy Readoption destination normalized to %q, want %q",
				disposition.Destination.Path,
				specificRepositoryPath,
			)
		}
	}

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	sourceBaseline, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatalf("load maintained Source Baseline: %v", err)
	}
	if len(sourceBaseline.Entries) != 60 ||
		sourceBaseline.Identity.EntryCount != 60 ||
		len(sourceBaseline.Accounting) != 51 {
		t.Fatalf(
			"maintained Source Baseline counts = identity %d entries %d accounting %d",
			sourceBaseline.Identity.EntryCount,
			len(sourceBaseline.Entries),
			len(sourceBaseline.Accounting),
		)
	}
	for _, transitionID := range catalog.TransitionIDs() {
		contract, err := catalog.UpgradeRetentionContract(transitionID)
		if err != nil {
			t.Fatalf("load maintained retention contract %s: %v", transitionID, err)
		}
		if len(contract.PriorClauses) == 0 ||
			len(contract.PriorClauses) != len(contract.Accounting) {
			t.Fatalf("retention contract %s is incomplete: prior=%d accounting=%d", transitionID, len(contract.PriorClauses), len(contract.Accounting))
		}
	}
}

func TestNestedCarrierWarningLeavesNestedSourcesOutOfPreservation(t *testing.T) {
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "packages/api/AGENTS.md", "nested policy\n")
	commitInspectionRepository(t, repo, "seed")

	plan, err := PlanRootPreservation(inspectPreservationRepository(t, repo), RootPreservationRequest{
		Mode: PreservationModePreservation,
	})
	if err != nil {
		t.Fatalf("plan nested carrier: %v", err)
	}
	if len(plan.Backups) != 0 || len(plan.SourceBaseline.Entries) != 0 {
		t.Fatalf("nested carrier entered mutation plan: backups=%+v entries=%+v", plan.Backups, plan.SourceBaseline.Entries)
	}
	if !hasPreservationFinding(plan.Warnings, "baseline.inventory.nested-carrier-conflict") {
		t.Fatalf("nested conflict warning missing: %+v", plan.Warnings)
	}
}

func TestNestedCarrierWarningUnsafeAliasRemainsNonBlocking(t *testing.T) {
	repo := newInspectionRepository(t)
	if err := os.MkdirAll(filepath.Join(repo, "packages", "api"), 0o755); err != nil {
		t.Fatalf("create nested carrier directory: %v", err)
	}
	if err := os.Symlink("../../../outside.md", filepath.Join(repo, "packages", "api", "AGENTS.md")); err != nil {
		t.Fatalf("create escaping nested alias: %v", err)
	}
	commitInspectionRepository(t, repo, "seed")

	inspection := inspectPreservationRepository(t, repo)
	if len(inspection.Snapshot.Blocking) != 0 {
		t.Fatalf("nested unsafe alias became apply-blocking: %+v", inspection.Snapshot.Blocking)
	}
	plan, err := PlanRootPreservation(inspection, RootPreservationRequest{
		Mode: PreservationModeGreenfield,
	})
	if err != nil {
		t.Fatalf("plan nested unsafe alias: %v", err)
	}
	if plan.State != PreservationStateReady || len(plan.Backups) != 0 {
		t.Fatalf("nested unsafe alias entered preservation plan: %+v", plan)
	}
	if !hasPreservationFinding(plan.Warnings, "baseline.inventory.nested-carrier-conflict") {
		t.Fatalf("nested unsafe alias warning missing: %+v", plan.Warnings)
	}
}

func inspectPreservationRepository(t *testing.T, repo string) RepositoryInspection {
	t.Helper()
	inspection, err := InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatalf("inspect preservation repository: %v", err)
	}
	return inspection
}

func hasPreservationFinding(findings []Finding, code string) bool {
	for _, finding := range findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
