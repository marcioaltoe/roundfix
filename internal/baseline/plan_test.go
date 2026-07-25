// Suite: Portable Baseline Plans
// Invariant: equivalent bounded repository state and normalized decisions produce one strict, portable, digest-bound plan without writes.
// Boundary IN: plan assembly, codecs, rendering, ledger projection, digest, and clone preimage validation.
// Boundary OUT: transaction apply, rollback, and interactive decision collection.

package baseline

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestPlanDocumentStrictCodecs(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	encoded, err := MarshalPlanDocument(plan)
	if err != nil {
		t.Fatalf("marshal Plan Document: %v", err)
	}
	decoded, err := ParsePlanDocument(encoded)
	if err != nil {
		t.Fatalf("parse Plan Document: %v", err)
	}
	if !reflect.DeepEqual(decoded, plan) {
		t.Fatalf("Plan Document round trip differs")
	}

	var raw map[string]any
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatal(err)
	}
	raw["unknown"] = true
	unknown, _ := json.Marshal(raw)
	if _, err := ParsePlanDocument(unknown); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("unknown field error = %v", err)
	}
	duplicate := bytes.Replace(encoded, []byte(`"schemaVersion":`),
		[]byte(`"schemaVersion":"roundfix/baseline-plan/v1","schemaVersion":`), 1)
	if _, err := ParsePlanDocument(duplicate); err == nil || !strings.Contains(err.Error(), "duplicate JSON key") {
		t.Fatalf("duplicate key error = %v", err)
	}

	result := Result{
		SchemaVersion: ResultSchemaVersion, Operation: "plan", State: "action_required",
		Category: "decision", Message: "missing", NextAction: "answer",
		VerifiedPostimages: []Postimage{}, Warnings: []Finding{}, Recommendations: []string{},
	}
	resultJSON, err := MarshalResult(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	if _, err := ParseResult(resultJSON); err != nil {
		t.Fatalf("parse result: %v", err)
	}
}

func TestFileChangesProjectionRejectsMismatch(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	if len(plan.FileChanges) == 0 {
		t.Fatal("plan has no file changes")
	}
	plan.FileChanges[0].ManagedEntries = append(plan.FileChanges[0].ManagedEntries, "invented")
	plan.PlanDigest, _ = computePlanDigest(plan)
	if err := ValidatePlanDocument(plan); err == nil ||
		!strings.Contains(err.Error(), "fileChanges does not match") {
		t.Fatalf("projection mismatch error = %v", err)
	}
}

func TestPlanDigestBindsExactPostimagesAndIgnoresProjection(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	original := plan.PlanDigest

	projected := plan
	projected.FileChanges = append([]FileChange(nil), plan.FileChanges...)
	projected.FileChanges[0].Action = "invented"
	digest, err := computePlanDigest(projected)
	if err != nil {
		t.Fatal(err)
	}
	if digest != original {
		t.Fatalf("projection changed digest: got %s want %s", digest, original)
	}

	changed := plan
	changed.Postimages = append([]Postimage(nil), plan.Postimages...)
	changed.Postimages[0].Content = append([]byte(nil), changed.Postimages[0].Content...)
	changed.Postimages[0].Content = append(changed.Postimages[0].Content, 'x')
	changed.Postimages[0].ContentIdentity = planContentIdentity(changed.Postimages[0].Content)
	digest, err = computePlanDigest(changed)
	if err != nil {
		t.Fatal(err)
	}
	if digest == original {
		t.Fatal("exact postimage bytes did not change Plan Digest")
	}
}

func TestCrossClonePlanAcceptsMatchingIdentityAndPreimages(t *testing.T) {
	source := newPlanRepository(t)
	plan := buildTestPlan(t, source)
	clone := filepath.Join(t.TempDir(), "clone")
	runPlanGit(t, "", "clone", "--quiet", source, clone)
	if err := ValidatePlanRepository(context.Background(), clone, plan); err != nil {
		t.Fatalf("matching clone rejected portable plan: %v", err)
	}
	writeInspectionFile(t, clone, ".agents/skills/context7/SKILL.md", "changed\n")
	if err := ValidatePlanRepository(context.Background(), clone, plan); err == nil ||
		!strings.Contains(err.Error(), "stale") {
		t.Fatalf("changed bounded preimage error = %v", err)
	}
}

func TestPlanDeterminismAndNoMutation(t *testing.T) {
	repo := newPlanRepository(t)
	before := testRepositoryDigest(t, repo)
	first := buildTestPlan(t, repo)
	second := buildTestPlan(t, repo)
	after := testRepositoryDigest(t, repo)
	if first.PlanDigest != second.PlanDigest {
		t.Fatalf("Plan Digests differ: %s != %s", first.PlanDigest, second.PlanDigest)
	}
	firstJSON, _ := MarshalPlanDocument(first)
	secondJSON, _ := MarshalPlanDocument(second)
	if !bytes.Equal(firstJSON, secondJSON) {
		t.Fatal("equivalent planning produced different JSON bytes")
	}
	if before != after {
		t.Fatalf("planning mutated repository: before=%s after=%s", before, after)
	}
	if strings.Contains(string(firstJSON), filepath.ToSlash(repo)) {
		t.Fatalf("portable plan contains checkout path %q", repo)
	}
	for _, postimage := range first.Postimages {
		if postimage.Path == "" || filepath.IsAbs(postimage.Path) {
			t.Fatalf("non-portable postimage path %q", postimage.Path)
		}
	}
}

func TestInstructionHierarchyRendersActivePointersOnce(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	agents := string(planPostimage(t, plan, "AGENTS.md").Content)

	if !strings.Contains(agents, "### Instruction hierarchy") ||
		!strings.Contains(
			agents,
			"cannot weaken a universal Normative Clause or confirmed project decision",
		) {
		t.Fatalf("AGENTS.md has no precedence contract:\n%s", agents)
	}

	wantPointers := []string{
		"docs/agents/agent-instructions.md",
		"docs/agents/skill-dispatch.md",
		"docs/agents/domain.md",
		"docs/agents/docs-layout.md",
		"docs/agents/spec-routing.md",
		"docs/agents/issue-tracker.md",
		"docs/agents/autonomous-work.md",
		"docs/agents/go.md",
		"docs/agents/cli.md",
		"docs/agents/tui.md",
	}
	position := -1
	for _, pointer := range wantPointers {
		if count := strings.Count(agents, pointer); count != 1 {
			t.Fatalf("AGENTS.md pointer %q count = %d, want 1:\n%s", pointer, count, agents)
		}
		next := strings.Index(agents, pointer)
		if next <= position {
			t.Fatalf("AGENTS.md pointer %q is outside hierarchy order:\n%s", pointer, agents)
		}
		position = next
	}
	for _, inactive := range []string{
		"docs/agents/external-triage.md",
		"docs/agents/secondbrain.md",
		"docs/agents/specific-repository.md",
	} {
		if strings.Contains(agents, inactive) {
			t.Fatalf("AGENTS.md contains inactive pointer %q:\n%s", inactive, agents)
		}
	}
}

func TestInstructionHierarchyPreservesPlanAndResultSchemas(t *testing.T) {
	plan := buildTestPlan(t, newPlanRepository(t))
	planJSON, err := MarshalPlanDocument(plan)
	if err != nil {
		t.Fatalf("MarshalPlanDocument() error = %v", err)
	}
	assertJSONFields(t, planJSON, []string{
		"schemaVersion",
		"repository",
		"catalog",
		"profile",
		"decisions",
		"retention",
		"preimages",
		"postimages",
		"warnings",
		"setupManifest",
		"managedEntries",
		"fileChanges",
		"planDigest",
	})

	resultJSON, err := MarshalResult(planReadyResult(t, plan))
	if err != nil {
		t.Fatalf("MarshalResult() error = %v", err)
	}
	assertJSONFields(t, resultJSON, []string{
		"schemaVersion",
		"operation",
		"state",
		"planDigest",
		"verifiedPostimages",
		"warnings",
		"recommendations",
	})
}

func TestADRLifecycleContract(t *testing.T) {
	docsLayout := string(planPostimage(
		t,
		buildTestPlan(t, newPlanRepository(t)),
		"docs/agents/docs-layout.md",
	).Content)
	required := []string{
		"status: proposed # proposed | accepted | rejected | deprecated | superseded",
		"created_at: YYYY-MM-DDTHH:MM:SSZ",
		"updated_at: YYYY-MM-DDTHH:MM:SSZ",
		"deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ",
		"superseded_by: null # null or ADR-NNNN",
		"Only `accepted` is active.",
		"without lifecycle frontmatter as active unless its body explicitly marks it inactive",
		"Do not rewrite existing ADRs solely to adopt lifecycle metadata.",
		"domain-modeling/ADR-FORMAT.md",
	}
	assertContractFragments(t, docsLayout, required)

	fixtures := []struct {
		name   string
		path   string
		active bool
	}{
		{name: "accepted lifecycle", path: "accepted.md", active: true},
		{name: "deprecated lifecycle", path: "deprecated.md", active: false},
		{name: "legacy active", path: "legacy-active.md", active: true},
		{name: "legacy explicitly inactive", path: "legacy-inactive.md", active: false},
	}
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			content, err := fs.ReadFile(
				compatibilityFixtures,
				"testdata/adr-lifecycle/"+fixture.path,
			)
			if err != nil {
				t.Fatalf("read ADR lifecycle fixture: %v", err)
			}
			if got := adrFixtureIsActive(content); got != fixture.active {
				t.Fatalf("adrFixtureIsActive() = %t, want %t", got, fixture.active)
			}
		})
	}

	existing, err := fs.ReadFile(
		compatibilityFixtures,
		"testdata/adr-lifecycle/legacy-active.md",
	)
	if err != nil {
		t.Fatalf("read existing ADR fixture: %v", err)
	}
	repo := newPlanRepository(t)
	const existingPath = "docs/adr/0001-legacy-active.md"
	writeInspectionFile(t, repo, existingPath, string(existing))
	commitInspectionRepository(t, repo, "add existing ADR fixture")
	plan := buildTestPlan(t, repo)
	for _, postimage := range plan.Postimages {
		if postimage.Path == existingPath {
			t.Fatalf("plan rewrites existing ADR solely for lifecycle metadata")
		}
	}
	after, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(existingPath)))
	if err != nil {
		t.Fatalf("read existing ADR after planning: %v", err)
	}
	if !bytes.Equal(after, existing) {
		t.Fatal("existing ADR fixture bytes changed during planning")
	}
}

func TestFindingsOperationalContract(t *testing.T) {
	docsLayout := string(planPostimage(
		t,
		buildTestPlan(t, newPlanRepository(t)),
		"docs/agents/docs-layout.md",
	).Content)
	required := []string{
		"status: pending # pending | partial | deferred | done",
		"created_at: YYYY-MM-DD",
		"updated_at: YYYY-MM-DD",
		"# <Area> — <short title> (YYYY-MM-DD)",
		"session or investigation",
		"## 1. <Finding title — symptom, not hypothesis>",
		"Symptom / evidence:",
		"Root cause: <proven cause, or `unknown`",
		"Action / suggestion:",
		"route to a Spec",
		"## What worked — keep",
		"## Addendum — YYYY-MM-DD — <short title>",
		"`pending` when the finding is new and has no implementation Spec",
		"`partial` when a linked Spec covers only the selected implementation scope",
		"`deferred` only when the finding will not be implemented",
		"`status: done` as soon as the Spec is created and linked",
		"append evidence and routing links as dated addenda",
		"Update `updated_at` whenever status changes or an evidence addendum is appended",
	}
	assertContractFragments(t, docsLayout, required)
}

func assertContractFragments(t *testing.T, content string, required []string) {
	t.Helper()
	for _, fragment := range required {
		if !strings.Contains(content, fragment) {
			t.Errorf("generated contract missing %q", fragment)
		}
		mutated := strings.ReplaceAll(content, fragment, "")
		if strings.Contains(mutated, fragment) {
			t.Fatalf("contract mutation did not remove %q", fragment)
		}
	}
}

func adrFixtureIsActive(content []byte) bool {
	text := string(content)
	if strings.HasPrefix(text, "---\n") {
		frontmatter, _, found := strings.Cut(strings.TrimPrefix(text, "---\n"), "\n---\n")
		if found {
			for _, line := range strings.Split(frontmatter, "\n") {
				if value, ok := strings.CutPrefix(strings.TrimSpace(line), "status:"); ok {
					return strings.TrimSpace(value) == "accepted"
				}
			}
		}
	}

	lower := strings.ToLower(text)
	for _, inactive := range []string{
		"status: proposed",
		"status: rejected",
		"status: deprecated",
		"status: superseded",
	} {
		if strings.Contains(lower, inactive) {
			return false
		}
	}
	return true
}

func TestGreenfieldRepositoryExtensionDoesNotCreateEmptyCarrier(t *testing.T) {
	repo := newPlanRepository(t)
	plan := buildPlanWithRepositoryExtension(t, repo)

	for _, postimage := range plan.Postimages {
		switch postimage.Path {
		case "docs/agents/specific-repository.md",
			"docs/agents/repository.md",
			"docs/agents/repository-rules.md":
			t.Fatalf("greenfield planned empty repository carrier %q", postimage.Path)
		}
	}
	for _, artifact := range plan.SetupManifest.ManagedArtifacts {
		if artifact.ID == "root.repository-extension" {
			t.Fatal("greenfield manifest linked an absent repository carrier")
		}
	}
	agents := planPostimage(t, plan, "AGENTS.md")
	if bytes.Contains(agents.Content, []byte("specific-repository.md")) ||
		bytes.Contains(agents.Content, []byte("repository.md")) ||
		bytes.Contains(agents.Content, []byte("repository-rules.md")) {
		t.Fatalf("greenfield AGENTS.md linked an absent repository carrier:\n%s", agents.Content)
	}
}

func assertJSONFields(t *testing.T, data []byte, want []string) {
	t.Helper()

	var document map[string]json.RawMessage
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatalf("decode JSON fields: %v", err)
	}
	got := make([]string, 0, len(document))
	for field := range document {
		got = append(got, field)
	}
	slices.Sort(got)
	want = append([]string(nil), want...)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("JSON fields = %v, want %v", got, want)
	}
}

func planReadyResult(t *testing.T, plan PlanDocument) Result {
	t.Helper()

	return Result{
		SchemaVersion:      ResultSchemaVersion,
		Operation:          "plan",
		State:              "ready",
		PlanDigest:         plan.PlanDigest,
		VerifiedPostimages: []Postimage{},
		Warnings:           []Finding{},
		Recommendations:    []string{},
	}
}

func TestRepositoryExtensionMigratesOneLegacyCarrier(t *testing.T) {
	repo := newPlanRepository(t)
	const rules = "# Repository rules\n\nKeep the Fluxus boundary explicit.\n"
	writeInspectionFile(t, repo, "docs/agents/repository.md", rules)
	commitInspectionRepository(t, repo, "seed legacy repository rules")

	plan := buildPlanWithRepositoryExtension(t, repo)
	canonical := planPostimage(t, plan, "docs/agents/specific-repository.md")
	if string(canonical.Content) != rules {
		t.Fatalf("canonical repository rules = %q, want exact legacy bytes %q", canonical.Content, rules)
	}
	legacy := planPostimage(t, plan, "docs/agents/repository.md")
	if legacy.Kind != PreimageMissing {
		t.Fatalf("legacy repository carrier postimage kind = %q, want missing", legacy.Kind)
	}
	agents := planPostimage(t, plan, "AGENTS.md")
	if !bytes.Contains(agents.Content, []byte("docs/agents/specific-repository.md")) ||
		bytes.Contains(agents.Content, []byte("docs/agents/repository.md")) {
		t.Fatalf("AGENTS.md did not converge on the canonical carrier:\n%s", agents.Content)
	}
}

func TestRepositoryExtensionDropsLegacyEmptyScaffold(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(
		t,
		repo,
		"docs/agents/repository.md",
		"# Repository instructions\n\nAdd project-specific hard rules here. Setup preserves this file byte-for-byte.\n",
	)
	commitInspectionRepository(t, repo, "seed empty repository scaffold")

	plan := buildPlanWithRepositoryExtension(t, repo)
	legacy := planPostimage(t, plan, "docs/agents/repository.md")
	if legacy.Kind != PreimageMissing {
		t.Fatalf("legacy scaffold postimage kind = %q, want missing", legacy.Kind)
	}
	for _, postimage := range plan.Postimages {
		if postimage.Path == "docs/agents/specific-repository.md" {
			t.Fatalf("legacy empty scaffold created canonical carrier with bytes %q", postimage.Content)
		}
	}
	agents := planPostimage(t, plan, "AGENTS.md")
	if bytes.Contains(agents.Content, []byte("specific-repository.md")) ||
		bytes.Contains(agents.Content, []byte("repository.md")) {
		t.Fatalf("legacy empty scaffold retained a repository carrier pointer:\n%s", agents.Content)
	}
}

func TestRepositoryExtensionRejectsDivergentLegacyCarriers(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, "docs/agents/repository.md", "first repository rule\n")
	writeInspectionFile(t, repo, "docs/agents/repository-rules.md", "different repository rule\n")
	commitInspectionRepository(t, repo, "seed conflicting repository rules")

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisionsWithRepositoryExtension(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "classification" ||
		!strings.Contains(outcome.Result.Message, "repository-specific rule carriers conflict") {
		t.Fatalf("divergent legacy carrier outcome = %+v", outcome)
	}
}

func TestDisabledRepositoryExtensionRejectsNonemptyLegacyCarrier(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, legacyRepositoryPath, "keep this repository rule\n")
	commitInspectionRepository(t, repo, "seed legacy repository rules")

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "classification" ||
		!strings.Contains(outcome.Result.Message, "require repository.extension.enabled=true") {
		t.Fatalf("disabled legacy carrier migration outcome = %+v", outcome)
	}
}

func TestSemanticRuleDistributionMovesExactBytesAndAccountsLedgers(t *testing.T) {
	repo := newPlanRepository(t)
	semanticRule := []byte("Keep requested CLI output on stdout.\n")
	residualRule := []byte("Keep the repository-specific release name.\n")
	sourceBytes := append(append([]byte(nil), semanticRule...), residualRule...)
	writeInspectionFile(t, repo, "AGENTS.md", string(sourceBytes))
	commitInspectionRepository(t, repo, "seed segmented repository rules")

	unresolved, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModePreservation},
	)
	if err != nil {
		t.Fatal(err)
	}
	segmentation, err := NewRuleSegmentationSnapshot(unresolved.SourceBaseline)
	if err != nil {
		t.Fatal(err)
	}
	entry := unresolved.SourceBaseline.Entries[0]
	proposal := RuleSegmentationProposal{
		SchemaVersion:  RuleSegmentationProposalSchemaVersion,
		SnapshotDigest: segmentation.SnapshotDigest,
		SourceBaseline: ClassificationSource{
			ID:     segmentation.SourceBaseline.ID,
			Digest: segmentation.SourceBaseline.Digest,
		},
		Segments: []RuleSegmentProposal{
			{
				EntryID: entry.ID,
				Start:   0,
				End:     len(semanticRule),
				Digest:  ruleSegmentDigest(semanticRule),
			},
			{
				EntryID: entry.ID,
				Start:   len(semanticRule),
				End:     len(sourceBytes),
				Digest:  ruleSegmentDigest(residualRule),
			},
		},
	}
	classifiedSource, err := MaterializeRuleSegments(segmentation, proposal)
	if err != nil {
		t.Fatal(err)
	}
	semanticEntry := classifiedSource.Entries[0]
	residualEntry := classifiedSource.Entries[1]
	document := decisionDocumentForSource(classifiedSource, []ReadoptionDisposition{
		{
			EntryID:        semanticEntry.ID,
			EntryDigest:    semanticEntry.Digest,
			Classification: "normative-clause",
			Disposition:    "repository-document",
			Destination: &ReadoptionDestination{
				DocumentType: "agent-guide",
				Path:         "docs/agents/cli.md",
				Digest:       semanticEntry.Digest,
			},
			Reason: "The active CLI guide owns this repository policy.",
		},
		{
			EntryID:        residualEntry.ID,
			EntryDigest:    residualEntry.Digest,
			Classification: "normative-clause",
			Disposition:    "repository-rules",
			Destination: &ReadoptionDestination{
				DocumentType:  "repository-rules",
				Path:          specificRepositoryPath,
				Digest:        residualEntry.Digest,
				ProposedBytes: base64.StdEncoding.EncodeToString(residualEntry.SourceBytes),
			},
			Reason: "No active semantic guide owns this repository-specific policy.",
		},
	})

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisionsWithRepositoryExtension(),
		Preservation: RootPreservationRequest{
			Mode:           PreservationModePreservation,
			Decisions:      &document,
			SourceBaseline: &classifiedSource,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil {
		t.Fatalf("semantic distribution returned action: %+v", outcome.Result)
	}
	plan := *outcome.Plan

	guide := planPostimage(t, plan, "docs/agents/cli.md")
	blocks, err := parseRepositoryRuleBlocks(guide.Path, guide.Content)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 || !bytes.Equal(blocks[0].Body, semanticRule) {
		t.Fatalf("semantic block bodies = %+v, want exact bytes %q", blocks, semanticRule)
	}
	residual := planPostimage(t, plan, specificRepositoryPath)
	if !bytes.Equal(residual.Content, residualRule) {
		t.Fatalf("residual bytes = %q, want %q", residual.Content, residualRule)
	}

	retentionBySource := make(map[string]RetentionEvidence)
	for _, evidence := range plan.Retention {
		retentionBySource[evidence.FromClause] = evidence
	}
	if retentionBySource[semanticEntry.ID].Disposition != "repository-document" ||
		retentionBySource[semanticEntry.ID].Targets[0] != "docs/agents/cli.md" ||
		retentionBySource[residualEntry.ID].Disposition != "repository-rules" ||
		retentionBySource[residualEntry.ID].Targets[0] != specificRepositoryPath {
		t.Fatalf("semantic distribution retention = %+v", plan.Retention)
	}
	var repositoryBlockLedger bool
	for _, managed := range plan.ManagedEntries {
		if managed.Path == guide.Path &&
			managed.Kind == "repository-owned" &&
			strings.HasPrefix(managed.ID, "repository-rule:") &&
			managed.ContentIdentity == planContentIdentity(semanticRule) {
			repositoryBlockLedger = true
		}
	}
	if !repositoryBlockLedger {
		t.Fatal("semantic repository-owned block is absent from the managed-entry ledger")
	}
	for _, artifact := range plan.SetupManifest.ManagedArtifacts {
		if strings.HasPrefix(artifact.ID, "repository-rule:") {
			t.Fatal("repository-owned semantic block leaked into SetupManifest.ManagedArtifacts")
		}
	}
}

func decisionDocumentForSource(
	source ReadoptionSourceBaseline,
	dispositions []ReadoptionDisposition,
) DecisionDocument {
	document := DecisionDocument{
		SchemaVersion: DecisionDocumentSchemaVersion,
		Version:       DecisionDocumentVersion,
		Decisions:     []DecisionValue{},
		Readoption: &ReadoptionDecisions{
			Dispositions: dispositions,
		},
	}
	document.Readoption.SourceBaseline.ID = source.ID
	document.Readoption.SourceBaseline.Digest = source.Digest
	return document
}

func buildSemanticCarrierPlan(t *testing.T, repo, carrierPath string) PlanDocument {
	t.Helper()
	unresolved, err := PlanRootPreservation(
		inspectPreservationRepository(t, repo),
		RootPreservationRequest{Mode: PreservationModePreservation},
	)
	if err != nil {
		t.Fatal(err)
	}
	var entry ReadoptionSourceEntry
	for _, candidate := range unresolved.SourceBaseline.Entries {
		if candidate.Path == carrierPath {
			entry = candidate
			break
		}
	}
	if entry.ID == "" {
		t.Fatalf("recognized carrier %q has no Source Baseline Entry", carrierPath)
	}
	document := decisionDocumentForSource(
		unresolved.SourceBaseline,
		[]ReadoptionDisposition{{
			EntryID:        entry.ID,
			EntryDigest:    entry.Digest,
			Classification: "normative-clause",
			Disposition:    "repository-document",
			Destination: &ReadoptionDestination{
				DocumentType: "agent-guide",
				Path:         "docs/agents/cli.md",
				Digest:       entry.Digest,
			},
			Reason: "The active CLI guide owns this repository policy.",
		}},
	)
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode:      PreservationModePreservation,
			Decisions: &document,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan == nil {
		t.Fatalf("semantic carrier plan returned action: %+v", outcome.Result)
	}
	return *outcome.Plan
}

func TestResidualCarrierRemovesAllRecognizedEmptyResults(t *testing.T) {
	for _, carrierPath := range []string{
		specificRepositoryPath,
		legacyRepositoryPath,
		legacyRepositoryRulesPath,
	} {
		t.Run(carrierPath, func(t *testing.T) {
			repo := newPlanRepository(t)
			rule := []byte("Keep CLI diagnostics on stderr.\n")
			writeInspectionFile(t, repo, carrierPath, string(rule))
			commitInspectionRepository(t, repo, "seed recognized repository-rule carrier")

			unresolved, err := PlanRootPreservation(
				inspectPreservationRepository(t, repo),
				RootPreservationRequest{Mode: PreservationModePreservation},
			)
			if err != nil {
				t.Fatal(err)
			}
			if len(unresolved.SourceBaseline.Entries) != 1 ||
				unresolved.SourceBaseline.Entries[0].Path != carrierPath {
				t.Fatalf("recognized carrier source inventory = %+v", unresolved.SourceBaseline.Entries)
			}
			entry := unresolved.SourceBaseline.Entries[0]
			document := decisionDocumentForSource(
				unresolved.SourceBaseline,
				[]ReadoptionDisposition{{
					EntryID:        entry.ID,
					EntryDigest:    entry.Digest,
					Classification: "normative-clause",
					Disposition:    "repository-document",
					Destination: &ReadoptionDestination{
						DocumentType: "agent-guide",
						Path:         "docs/agents/cli.md",
						Digest:       entry.Digest,
					},
					Reason: "The active CLI guide owns this repository policy.",
				}},
			)
			outcome, err := BuildPlan(context.Background(), PlanRequest{
				Repository: repo,
				ProfileID:  "go-cli-tui",
				Decisions:  planTestDecisions(),
				Preservation: RootPreservationRequest{
					Mode:      PreservationModePreservation,
					Decisions: &document,
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			if outcome.Plan == nil {
				t.Fatalf("recognized carrier distribution returned action: %+v", outcome.Result)
			}
			plan := *outcome.Plan
			removed := planPostimage(t, plan, carrierPath)
			if removed.Kind != PreimageMissing {
				t.Fatalf("recognized carrier %q postimage = %+v, want missing", carrierPath, removed)
			}
			for _, postimage := range plan.Postimages {
				if postimage.Path == specificRepositoryPath && postimage.Kind == PreimageRegular {
					t.Fatalf("empty redistribution retained residual bytes %q", postimage.Content)
				}
			}
			agents := planPostimage(t, plan, "AGENTS.md")
			if bytes.Contains(agents.Content, []byte(specificRepositoryPath)) ||
				bytes.Contains(agents.Content, []byte(legacyRepositoryPath)) ||
				bytes.Contains(agents.Content, []byte(legacyRepositoryRulesPath)) {
				t.Fatalf("empty redistribution retained a root pointer:\n%s", agents.Content)
			}
			if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
				t.Fatalf("apply recognized carrier distribution: %v", err)
			}
			if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(carrierPath))); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("recognized carrier remains after apply: %v", err)
			}
		})
	}
}

func TestNestedCarrierRemainsByteIdenticalAndWarningOnly(t *testing.T) {
	repo := newPlanRepository(t)
	const nestedPath = "services/payments/AGENTS.md"
	nested := []byte("Keep this nested repository policy byte-identical.\r\n")
	writeInspectionFile(t, repo, nestedPath, string(nested))
	commitInspectionRepository(t, repo, "seed arbitrary nested carrier")

	plan := buildTestPlan(t, repo)
	var warned bool
	for _, warning := range plan.Warnings {
		if warning.Code == "baseline.inventory.nested-carrier-conflict" &&
			warning.Path == nestedPath {
			warned = true
		}
	}
	if !warned {
		t.Fatalf("nested carrier warnings = %+v", plan.Warnings)
	}
	for _, postimage := range plan.Postimages {
		if postimage.Path == nestedPath {
			t.Fatal("arbitrary nested carrier became a mutation target")
		}
	}
	if _, err := ApplyPlan(context.Background(), repo, plan, plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(nestedPath)))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, nested) {
		t.Fatalf("nested carrier bytes = %q, want %q", after, nested)
	}
}

func TestPlanDeterminismMatchesMaintainedManagedEntryFixture(t *testing.T) {
	fixturePath := filepath.Join(
		"testdata", "parity-corpus", "v1", "fixtures", "greenfield-go-cli-tui.json",
	)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		ManagedEntryLedger []struct {
			ID       string `json:"id"`
			Path     string `json:"path"`
			Kind     string `json:"kind"`
			Module   string `json:"module"`
			Template string `json:"template"`
			Version  string `json:"version"`
			Digest   string `json:"digest"`
		} `json:"managedEntryLedger"`
		PlannedByteSequence []struct {
			Path          string `json:"path"`
			AfterIdentity string `json:"afterIdentity"`
		} `json:"plannedByteSequence"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile, err := ResolveProfile("", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	_, artifacts, err := resolveManagedArtifacts(catalog, profile, planTestDecisions(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != len(fixture.ManagedEntryLedger) {
		t.Fatalf("managed entry count = %d, want %d", len(artifacts), len(fixture.ManagedEntryLedger))
	}
	type managedEntryIdentity struct {
		Path     string
		Kind     string
		Module   string
		Template string
		Version  string
	}
	expected := make(map[string]managedEntryIdentity, len(fixture.ManagedEntryLedger))
	for _, entry := range fixture.ManagedEntryLedger {
		expected[entry.ID] = managedEntryIdentity{
			Path: entry.Path, Kind: entry.Kind, Module: entry.Module,
			Template: entry.Template, Version: entry.Version,
		}
	}
	for _, artifact := range artifacts {
		want, ok := expected[artifact.ID]
		if !ok {
			t.Fatalf("managed entry %q is not in maintained fixture", artifact.ID)
		}
		got := managedEntryIdentity{
			Path: artifact.Path, Kind: artifact.Kind, Module: artifact.Module,
			Template: artifact.Template, Version: artifact.Version,
		}
		if got != want {
			t.Fatalf("managed entry %q = %+v, want %+v", artifact.ID, got, want)
		}
	}

	plan := buildTestPlan(t, newPlanRepository(t))
	expectedPostimages := make(map[string]string)
	for _, entry := range fixture.PlannedByteSequence {
		if entry.Path != manifestPath &&
			entry.Path != "AGENTS.md" &&
			entry.Path != "docs/agents/docs-layout.md" {
			expectedPostimages[entry.Path] = entry.AfterIdentity
		}
	}
	for _, postimage := range plan.Postimages {
		want, ok := expectedPostimages[postimage.Path]
		if !ok {
			continue
		}
		if postimage.ContentIdentity != want {
			t.Fatalf("planned bytes %q identity = %s, want maintained %s",
				postimage.Path, postimage.ContentIdentity, want)
		}
		delete(expectedPostimages, postimage.Path)
	}
	if len(expectedPostimages) != 0 {
		t.Fatalf("maintained planned paths missing from plan: %v", sortedKeys(expectedPostimages))
	}
}

func TestPlanDocumentMissingDecisionsReturnsResultWithoutPartialPlan(t *testing.T) {
	repo := newPlanRepository(t)
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "decision" || outcome.Result.NextAction == "" {
		t.Fatalf("missing-decision outcome = %+v", outcome)
	}
	if _, err := MarshalResult(outcome.Result); err != nil {
		t.Fatalf("missing-decision result schema: %v", err)
	}
}

func TestPlanDocumentIncludesMaintainedUpgradeRetention(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, manifestPath, `{
  "schemaVersion": 1,
  "generator": {
    "skill": "setup-context-driven",
    "version": 1,
    "baseline": "baseline.managed-v2"
  },
  "managedArtifacts": []
}
`)
	plan := buildTestPlan(t, repo)
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	contract, err := catalog.UpgradeRetentionContract("transition.managed-v2-to-portable-v3")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Retention) != len(contract.PriorClauses) {
		t.Fatalf("retention entries = %d, want %d", len(plan.Retention), len(contract.PriorClauses))
	}
	for index, entry := range plan.Retention {
		if entry.FromClause != contract.Accounting[index].FromClause ||
			entry.Enforcement == "" ||
			entry.Disposition != contract.Accounting[index].Disposition ||
			!reflect.DeepEqual(entry.Targets, contract.Accounting[index].Targets) ||
			entry.Reason != contract.Accounting[index].Reason {
			t.Fatalf("retention entry %d = %+v, want contract accounting %+v",
				index, entry, contract.Accounting[index])
		}
	}
}

func TestPlanDocumentRejectsUnknownManifestRetentionWithoutPartialPlan(t *testing.T) {
	repo := newPlanRepository(t)
	writeInspectionFile(t, repo, manifestPath, `{
  "schemaVersion": 1,
  "generator": {
    "skill": "setup-context-driven",
    "version": 1,
    "baseline": "baseline.unknown"
  },
  "managedArtifacts": []
}
`)
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" ||
		outcome.Result.Category != "classification" || outcome.Result.NextAction == "" {
		t.Fatalf("unknown retention outcome = %+v", outcome)
	}
}

func buildTestPlan(t *testing.T, repo string) PlanDocument {
	t.Helper()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisions(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatalf("build plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("build plan returned result: %+v", outcome.Result)
	}
	return *outcome.Plan
}

func planTestDecisions() []DecisionValue {
	return []DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: true},
		{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		{ID: "runtime.design", Value: "claude opus xhigh"},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: false},
	}
}

func planTestDecisionsWithRepositoryExtension() []DecisionValue {
	decisions := planTestDecisions()
	for index := range decisions {
		if decisions[index].ID == "repository.extension.enabled" {
			decisions[index].Value = true
		}
	}
	return decisions
}

func buildPlanWithRepositoryExtension(t *testing.T, repo string) PlanDocument {
	t.Helper()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  planTestDecisionsWithRepositoryExtension(),
		Preservation: RootPreservationRequest{
			Mode: PreservationModeGreenfield,
		},
	})
	if err != nil {
		t.Fatalf("build repository-extension plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("build repository-extension plan returned result: %+v", outcome.Result)
	}
	return *outcome.Plan
}

func planPostimage(t *testing.T, plan PlanDocument, path string) Postimage {
	t.Helper()
	for _, postimage := range plan.Postimages {
		if postimage.Path == path {
			return postimage
		}
	}
	t.Fatalf("plan has no postimage for %q", path)
	return Postimage{}
}

func newPlanRepository(t *testing.T) string {
	t.Helper()
	repo := newInspectionRepository(t)
	writeInspectionFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeInspectionFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeInspectionFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitInspectionRepository(t, repo, "seed portable plan")
	return repo
}

func runPlanGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func testRepositoryDigest(t *testing.T, repo string) string {
	t.Helper()
	command := exec.Command("git", "status", "--short", "--untracked-files=all")
	command.Dir = repo
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git status: %v\n%s", err, output)
	}
	return string(output)
}
