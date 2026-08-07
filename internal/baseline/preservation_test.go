// Suite: Baseline root-instruction preservation
// Invariant: preservation modes either prove managed-only refresh safety or account for every root source before planning is ready.
// Boundary IN: bounded repository inspection, root backups, Decision Documents, Readoption, Source Baselines, and retention contracts.
// Boundary OUT: profile alignment, portable Plan Documents, file transactions, and ACP classification proposals.

package baseline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGreenfieldPlanBacksUpWithoutImport(t *testing.T) {
	t.Parallel()

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

func TestManagedRefreshPlanNeedsNoClassificationInputOrBackup(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", "repository-authored policy\n")
	commitInspectionRepository(t, repo, "seed managed refresh root")

	plan, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModeManagedRefresh},
	)
	if err != nil {
		t.Fatalf("plan managed refresh preservation: %v", err)
	}
	if plan.State != PreservationStateReady ||
		len(plan.SourceBaseline.Entries) != 0 ||
		plan.DecisionSkeleton != nil ||
		len(plan.Backups) != 0 {
		t.Fatalf("managed refresh preservation plan = %+v", plan)
	}
}

func TestManagedRefreshUnsafeRootCarrierStillBlocks(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	if err := os.Symlink("../outside.md", filepath.Join(repo, "AGENTS.md")); err != nil {
		t.Fatalf("create escaping root alias: %v", err)
	}
	commitInspectionRepository(t, repo, "seed unsafe root carrier")

	plan, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModeManagedRefresh},
	)
	if err != nil {
		t.Fatalf("plan managed refresh with unsafe root carrier: %v", err)
	}
	if plan.State != PreservationStateBlocked ||
		!hasRepositoryFinding(plan.Findings, "baseline.inventory.unsafe-alias", "AGENTS.md") {
		t.Fatalf("unsafe root carrier did not block managed refresh: %+v", plan)
	}
}

func TestPreservationRequiresEveryDisposition(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

func TestDecisionDocumentSkeletonDoesNotProposeManagedSemanticVersionBytes(t *testing.T) {
	t.Parallel()

	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, "AGENTS.md", strings.Join([]string{
		"<!-- setup-context-driven:begin id=root.core version=0.0.1 -->",
		"managed root guidance",
		"<!-- setup-context-driven:end id=root.core -->",
		"",
	}, "\n"))
	commitInspectionRepository(t, repo, "seed managed root")

	plan, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModePreservation},
	)
	if err != nil {
		t.Fatalf("plan managed-root preservation: %v", err)
	}
	if plan.State != PreservationStateReady ||
		plan.DecisionSkeleton != nil ||
		len(plan.SourceBaseline.Entries) != 0 {
		t.Fatalf("setup-managed bytes entered classification: %+v", plan)
	}
}

func TestDecisionDocumentSkeletonRejectsMalformedInput(t *testing.T) {
	t.Parallel()

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

// The maintained Source Baseline's expected shape. Named here so a legitimate
// corpus change moves one declared value instead of hunting literals, and so
// the diff says what moved.
const (
	maintainedSourceBaselineEntries    = 118
	maintainedSourceBaselineAccounting = 51
)

func TestReadoptionCompatibilityMaintainedFixture(t *testing.T) {
	// Sequential: can rewrite shared digest artifacts when the update flag is enabled.
	if !*updateBaselineDigests {
		t.Parallel()
	}
	if *updateBaselineDigests {
		regenerateMaintainedSourceBaseline(t)
		return
	}
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
		t.Fatalf(
			"load embedded catalog: %v; %s to regenerate stale derived artifacts",
			err,
			baselineDigestRegenerationHint,
		)
	}
	sourceBaseline, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatalf(
			"load maintained Source Baseline: %v; %s",
			err,
			baselineDigestRegenerationHint,
		)
	}
	// This compatibility fixture protects both internal count agreement and
	// the maintained corpus's exact shape. A legitimate corpus change must move
	// the named expectations above through the sanctioned re-recording workflow.
	if sourceBaseline.Identity.EntryCount != len(sourceBaseline.Entries) ||
		len(sourceBaseline.Entries) != maintainedSourceBaselineEntries ||
		len(sourceBaseline.Accounting) != maintainedSourceBaselineAccounting {
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

func regenerateMaintainedSourceBaseline(t *testing.T) {
	t.Helper()
	const baselineID = "baseline.standard-typescript-monorepo-0.0.1"
	baselineRoot := filepath.Join("assets", "source-baselines", baselineID)
	manifestPath := filepath.Join(baselineRoot, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read maintained Source Baseline manifest: %v", err)
	}
	var manifest sourceBaselineManifest
	if err := strictJSON(manifestData, &manifest); err != nil {
		t.Fatalf("decode maintained Source Baseline manifest: %v", err)
	}
	if manifest.ID != baselineID {
		t.Fatalf("maintained Source Baseline manifest id = %q", manifest.ID)
	}
	for index := range manifest.Entries {
		entry := &manifest.Entries[index]
		if !safeRelative(entry.Path) || !strings.HasPrefix(entry.Path, "corpus/") {
			t.Fatalf("Source Baseline entry %q has unsafe path %q", entry.ID, entry.Path)
		}
		sourcePath := filepath.Join(baselineRoot, filepath.FromSlash(entry.Path))
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			t.Fatalf("read Source Baseline entry %q: %v", entry.ID, err)
		}
		start, end, err := sourceBaselineMarkerSpan(source, entry.ID)
		if err != nil {
			t.Fatalf("locate Source Baseline entry %q: %v", entry.ID, err)
		}
		sum := sha256.Sum256(source[start:end])
		entry.Start = start
		entry.End = end
		entry.Digest = hex.EncodeToString(sum[:])
	}
	if err := validateMaintainedSourceBaselineEntries(baselineRoot, manifest.Entries); err != nil {
		t.Fatalf("self-validate regenerated Source Baseline: %v", err)
	}

	updatedManifest, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode maintained Source Baseline manifest: %v", err)
	}
	updatedManifest = append(updatedManifest, '\n')
	manifestSum := sha256.Sum256(updatedManifest)
	manifestDigest := hex.EncodeToString(manifestSum[:])
	corpusDigest := maintainedSourceBaselineCorpusDigest(t, filepath.Join(baselineRoot, "corpus"))

	identityPath := filepath.Join(baselineRoot, "baseline.json")
	identityData, err := os.ReadFile(identityPath)
	if err != nil {
		t.Fatalf("read maintained Source Baseline identity: %v", err)
	}
	var identity map[string]any
	identityDecoder := json.NewDecoder(bytes.NewReader(identityData))
	identityDecoder.UseNumber()
	if err := identityDecoder.Decode(&identity); err != nil {
		t.Fatalf("decode maintained Source Baseline identity: %v", err)
	}
	if err := identityDecoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("maintained Source Baseline identity has trailing JSON")
	}
	identity["entryCount"] = json.Number(fmt.Sprint(len(manifest.Entries)))
	identity["manifestDigest"] = manifestDigest
	identity["corpusDigest"] = corpusDigest
	updatedIdentity, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		t.Fatalf("encode maintained Source Baseline identity: %v", err)
	}
	updatedIdentity = append(updatedIdentity, '\n')

	indexPath := filepath.Join("assets", "source-baselines", "index.json")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read maintained Source Baseline index: %v", err)
	}
	var baselineIndex sourceBaselineIndex
	if err := strictJSON(indexData, &baselineIndex); err != nil {
		t.Fatalf("decode maintained Source Baseline index: %v", err)
	}
	found := false
	for position := range baselineIndex.Baselines {
		record := &baselineIndex.Baselines[position]
		if record.ID != baselineID {
			continue
		}
		record.EntryIDs = make([]string, len(manifest.Entries))
		for index, entry := range manifest.Entries {
			record.EntryIDs[index] = entry.ID
		}
		record.EntryCount = len(manifest.Entries)
		record.ManifestDigest = manifestDigest
		record.CorpusDigest = corpusDigest
		found = true
		break
	}
	if !found {
		t.Fatalf("maintained Source Baseline index has no record %q", baselineID)
	}
	updatedIndex, err := json.MarshalIndent(baselineIndex, "", "  ")
	if err != nil {
		t.Fatalf("encode maintained Source Baseline index: %v", err)
	}
	updatedIndex = append(updatedIndex, '\n')

	writeBaselineDerivedArtifact(t, manifestPath, updatedManifest)
	writeBaselineDerivedArtifact(t, identityPath, updatedIdentity)
	writeBaselineDerivedArtifact(t, indexPath, updatedIndex)
	regenerateCatalogCompatibilityFromAssets(t)
}

func sourceBaselineMarkerSpan(source []byte, entryID string) (int, int, error) {
	opening := []byte(fmt.Sprintf("<!-- source-baseline-entry: %s -->", entryID))
	closing := []byte(fmt.Sprintf("<!-- /source-baseline-entry: %s -->", entryID))
	if count := bytes.Count(source, opening); count != 1 {
		return 0, 0, fmt.Errorf("opening marker count = %d, want 1", count)
	}
	if count := bytes.Count(source, closing); count != 1 {
		return 0, 0, fmt.Errorf("closing marker count = %d, want 1", count)
	}
	start := bytes.Index(source, opening) + len(opening)
	if start >= len(source) || source[start] != '\n' {
		return 0, 0, errors.New("opening marker is not followed by a newline")
	}
	start++
	end := bytes.Index(source[start:], closing)
	if end < 0 {
		return 0, 0, errors.New("closing marker precedes opening marker")
	}
	end += start
	if end <= start {
		return 0, 0, errors.New("marker span is empty")
	}
	return start, end, nil
}

func validateMaintainedSourceBaselineEntries(
	baselineRoot string,
	entries []SourceBaselineEntry,
) error {
	for _, entry := range entries {
		if !safeRelative(entry.Path) || !strings.HasPrefix(entry.Path, "corpus/") {
			return fmt.Errorf("entry %q has unsafe path %q", entry.ID, entry.Path)
		}
		sourcePath := filepath.Join(baselineRoot, filepath.FromSlash(entry.Path))
		source, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("read entry %q: %w", entry.ID, err)
		}
		start, end, err := sourceBaselineMarkerSpan(source, entry.ID)
		if err != nil {
			return fmt.Errorf("locate entry %q: %w", entry.ID, err)
		}
		if entry.Start != start || entry.End != end {
			return fmt.Errorf(
				"entry %q span = %d:%d, marker span = %d:%d",
				entry.ID,
				entry.Start,
				entry.End,
				start,
				end,
			)
		}
		sum := sha256.Sum256(source[start:end])
		if got := hex.EncodeToString(sum[:]); entry.Digest != got {
			return fmt.Errorf(
				"entry %q digest = %q, marker span digest = %q",
				entry.ID,
				entry.Digest,
				got,
			)
		}
	}
	return nil
}

func maintainedSourceBaselineCorpusDigest(t *testing.T, corpusRoot string) string {
	t.Helper()
	files := make(map[string][]byte)
	if err := filepath.WalkDir(
		corpusRoot,
		func(filePath string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			if !entry.Type().IsRegular() {
				return fmt.Errorf("corpus entry %q is not a regular file", filePath)
			}
			relative, err := filepath.Rel(corpusRoot, filePath)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			files[filepath.ToSlash(relative)] = data
			return nil
		},
	); err != nil {
		t.Fatalf("walk maintained Source Baseline corpus: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("maintained Source Baseline corpus is empty")
	}
	return portableFileDigest(files)
}

func TestSourceBaselineRegenerationRejectsCorruptedSpan(t *testing.T) {
	t.Parallel()

	const baselineID = "baseline.standard-typescript-monorepo-0.0.1"
	baselineRoot := filepath.Join("assets", "source-baselines", baselineID)
	data, err := os.ReadFile(filepath.Join(baselineRoot, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest sourceBaselineManifest
	if err := strictJSON(data, &manifest); err != nil {
		t.Fatal(err)
	}
	corrupted := append([]SourceBaselineEntry(nil), manifest.Entries...)
	corrupted[0].Start++
	if err := validateMaintainedSourceBaselineEntries(baselineRoot, corrupted); err == nil {
		t.Fatal("corrupted Source Baseline span passed regeneration self-validation")
	}
}

func TestSourceBaselineGuidanceComposition(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	sourceBaseline, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatalf("load maintained Source Baseline: %v", err)
	}

	entries := make(map[string]SourceBaselineEntry, len(sourceBaseline.Entries))
	carriers := make(map[string]bool)
	for _, entry := range sourceBaseline.Entries {
		entries[entry.ID] = entry
		carriers[entry.Carrier] = true
	}
	for _, entryID := range []string{
		"clause.autonomous.delegate-through-roundfix",
		"clause.context.adr-01-template",
		"clause.core.research-authoritative-external-sources",
		"clause.core.require-tooling-authorization",
		"clause.secondbrain.01-consult-triggers",
		"clause.spec.routing-01-large-initiative",
		"rule.backend.boundary-contracts",
		"rule.external-triage",
		"rule.monorepo.context-boundaries",
	} {
		if _, ok := entries[entryID]; !ok {
			t.Errorf("Source Baseline entry %q is missing", entryID)
		}
	}
	for _, carrier := range []string{
		"docs/agents/external-triage.md",
		"docs/agents/monorepo.md",
		"docs/agents/skill-dispatch.md",
	} {
		if !carriers[carrier] {
			t.Errorf("semantic destination %q has no Source Baseline evidence", carrier)
		}
	}
	for _, accounting := range sourceBaseline.Accounting {
		for _, target := range accounting.Targets {
			if _, ok := entries[target]; !ok {
				t.Errorf("accounting target %q is not a current Source Baseline entry", target)
			}
		}
	}
}

func TestNestedCarrierWarningLeavesNestedSourcesOutOfPreservation(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
