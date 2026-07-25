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

func TestFormatterComposition(t *testing.T) {
	repository := newAlignedTypeScriptRepository(t)
	runPlanGit(t, repository, "init", "-q")
	runPlanGit(t, repository, "config", "user.email", "fixture@example.invalid")
	runPlanGit(t, repository, "config", "user.name", "Fixture Test")
	runPlanGit(t, repository, "config", "commit.gpgsign", "false")
	runPlanGit(t, repository, "add", ".")
	runPlanGit(t, repository, "commit", "-qm", "seed formatter fixture")

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository: repository,
		ProfileID:  "standard-typescript-monorepo",
		Decisions: []DecisionValue{
			{ID: "language.generated", Value: "English"},
			{ID: "verification.gate", Value: "make verify"},
			{ID: "identifier.strategy", Value: map[string]any{"kind": "uuid-v7"}},
			{ID: "http.contract", Value: map[string]any{"mode": "Post-only"}},
			{ID: "auth.provider", Value: completeAuthProviderDecision()},
			{ID: "spec.scaffold", Value: true},
			{ID: "domain.layout", Value: "single-context"},
			{ID: "triage.external", Value: true},
			{ID: "autonomous.enabled", Value: true},
			{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
			{ID: "runtime.design", Value: "claude opus xhigh"},
			{ID: "secondbrain.enabled", Value: true},
			{ID: "repository.extension.enabled", Value: false},
		},
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatalf("BuildPlan() formatter composition error = %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("BuildPlan() formatter composition result = %+v", outcome.Result)
	}
	plan := *outcome.Plan

	catalog := mustEmbeddedCatalog(t)
	profile, ok := catalog.Profile("standard-typescript-monorepo")
	if !ok {
		t.Fatal("standard TypeScript Profile is missing")
	}
	var contract struct {
		Formatter struct {
			FixturePaths []string `json:"fixturePaths"`
		} `json:"formatter"`
	}
	if err := json.Unmarshal(profile.Data, &contract); err != nil {
		t.Fatalf("decode formatter contract: %v", err)
	}
	postimages := make(map[string][]byte, len(plan.Postimages))
	for _, postimage := range plan.Postimages {
		if postimage.Kind != PreimageMissing {
			postimages[postimage.Path] = postimage.Content
		}
	}
	const goldenSeparator = "/golden/"
	for _, fixturePath := range contract.Formatter.FixturePaths {
		position := strings.Index(fixturePath, goldenSeparator)
		if position < 0 {
			t.Fatalf("formatter fixture %q is outside the golden root", fixturePath)
		}
		generatedPath := fixturePath[position+len(goldenSeparator):]
		generated, ok := postimages[generatedPath]
		if !ok {
			t.Errorf("Plan has no generated postimage for formatter fixture %q", generatedPath)
			continue
		}
		fixture, ok := catalog.Asset(fixturePath)
		if !ok {
			t.Errorf("formatter fixture %q is missing", fixturePath)
			continue
		}
		if !bytes.Equal(generated, fixture.Data) {
			t.Errorf("generated output %q differs from its formatter fixture", generatedPath)
		}
	}
	if _, ok := postimages[specificRepositoryPath]; ok {
		t.Fatalf("greenfield formatter composition created %q", specificRepositoryPath)
	}

	if _, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatalf("ApplyPlan() formatter composition error = %v", err)
	}
	result, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("ApplyPlan() empty formatter reapply error = %v", err)
	}
	if result.State != "verified" || !strings.Contains(result.Message, "already applied") {
		t.Fatalf("empty formatter reapply result = %+v", result)
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
			entry.Path != "docs/agents/agent-instructions.md" &&
			entry.Path != "docs/agents/docs-layout.md" &&
			entry.Path != "docs/agents/skill-dispatch.md" {
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

func TestToolingAuthorityClause(t *testing.T) {
	t.Parallel()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}

	const clause = "Do not create, edit, rename, move, or delete any linter, formatter, typechecker, test-runner, architecture-checker, build-tool, package-manager, code-generator, or other repository-tooling configuration, script, ignore file, plugin declaration, or version pin without express maintainer authorization. Setup completion, a Profile, a narrower guide, or a generic implementation request does not grant that authorization."
	for _, profileID := range catalog.ProfileIDs() {
		profileID := profileID
		t.Run(profileID, func(t *testing.T) {
			profile, err := ResolveProfile(t.TempDir(), profileID, catalog)
			if err != nil {
				t.Fatalf("resolve Profile %q: %v", profileID, err)
			}
			_, artifacts, err := resolveManagedArtifacts(catalog, profile, []DecisionValue{
				{ID: "language.generated", Value: "English"},
				{ID: "verification.gate", Value: "rtk make verify"},
			}, false)
			if err != nil {
				t.Fatalf("render Profile %q artifacts: %v", profileID, err)
			}
			for _, artifact := range artifacts {
				if artifact.ID != "guide.agent-instructions" {
					continue
				}
				if !strings.Contains(artifact.Body, clause) {
					t.Fatalf("Profile %q tooling-authority clause is incomplete:\n%s", profileID, artifact.Body)
				}
				return
			}
			t.Fatalf("Profile %q did not render guide.agent-instructions", profileID)
		})
	}
}

func TestIdentifierStrategyDecision(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatalf("resolve Standard TypeScript Monorepo Profile: %v", err)
	}
	if !slices.Contains(profile.Decisions, "identifier.strategy") {
		t.Fatalf("Profile decisions = %v, want identifier.strategy", profile.Decisions)
	}

	decisions := standardTypeScriptDecisions("make verify")
	decisions = slices.DeleteFunc(decisions, func(decision DecisionValue) bool {
		return decision.ID == "identifier.strategy"
	})
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   newPlanRepository(t),
		ProfileID:    profile.ID,
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatalf("build Plan without identifier strategy: %v", err)
	}
	if outcome.Plan != nil ||
		outcome.Result.State != "action_required" ||
		outcome.Result.Category != "decision" ||
		!strings.Contains(outcome.Result.Message, "identifier.strategy") {
		t.Fatalf("missing identifier strategy outcome = %+v", outcome)
	}
}

func TestAuthProviderDecision(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	source, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatalf("resolve Standard TypeScript Monorepo Profile: %v", err)
	}
	if !slices.Contains(source.Decisions, "auth.provider") {
		t.Fatalf("Profile decisions = %v, want auth.provider", source.Decisions)
	}

	decisions := standardTypeScriptDecisions("make verify")
	decisions = slices.DeleteFunc(decisions, func(decision DecisionValue) bool {
		return decision.ID == "auth.provider"
	})
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   newPlanRepository(t),
		ProfileID:    source.ID,
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatalf("build Plan without Better Auth provider: %v", err)
	}
	if outcome.Plan != nil ||
		outcome.Result.State != "action_required" ||
		!strings.Contains(outcome.Result.Message, "auth.provider") {
		t.Fatalf("missing Better Auth provider outcome = %+v", outcome)
	}

	input, err := NewProfileAdaptationDraft(
		source.ID,
		"without-better-auth",
		[]string{"autonomous-work"},
		[]string{"capability.stack.better-auth"},
		catalog,
	)
	if err != nil {
		t.Fatalf("create Profile adaptation without Better Auth: %v", err)
	}
	resolved, _, err := ResolveProfileDraft(t.TempDir(), input, catalog)
	if err != nil {
		t.Fatalf("resolve Profile adaptation without Better Auth: %v", err)
	}
	if slices.Contains(resolved.Capabilities, "capability.stack.better-auth") {
		t.Fatalf("adapted Profile retained Better Auth: %v", resolved.Capabilities)
	}
	if slices.Contains(resolved.Decisions, "auth.provider") {
		t.Fatalf("adapted Profile decisions = %v, want auth.provider omitted", resolved.Decisions)
	}
	if !slices.Contains(resolved.Decisions, "identifier.strategy") {
		t.Fatalf("adapted Profile decisions = %v, want identifier.strategy retained", resolved.Decisions)
	}

	var document map[string]any
	if err := json.Unmarshal(input.Document, &document); err != nil {
		t.Fatalf("decode Profile adaptation: %v", err)
	}
	selected, ok := document["decisions"].([]any)
	if !ok {
		t.Fatalf("Profile adaptation decisions = %#v", document["decisions"])
	}
	withAuth := make([]any, 0, len(selected)+1)
	for _, decisionID := range selected {
		withAuth = append(withAuth, decisionID)
		if decisionID == "http.contract" {
			withAuth = append(withAuth, "auth.provider")
		}
	}
	document["decisions"] = withAuth
	invalid, err := json.Marshal(document)
	if err != nil {
		t.Fatalf("encode invalid Profile adaptation: %v", err)
	}
	input.Document = invalid
	if _, _, err := ResolveProfileDraft(t.TempDir(), input, catalog); err == nil ||
		!strings.Contains(err.Error(), "custom.profile.decision.capability.unselected") {
		t.Fatalf("Profile retaining auth.provider without Better Auth error = %v", err)
	}
}

func TestDeriveBetterAuthHTTPContract(t *testing.T) {
	repository := newProjectDecisionPlanRepository(t)
	decisions := standardTypeScriptDecisions("make verify")
	provider := completeAuthProviderDecision()
	providerException := provider["routeException"].(map[string]any)
	providerException["methods"] = []any{"POST", "GET"}
	for index := range decisions {
		switch decisions[index].ID {
		case "auth.provider":
			decisions[index].Value = provider
		case "http.contract":
			decisions[index].Value = map[string]any{
				"mode": "Post-only",
				"exceptions": []any{
					map[string]any{
						"scope":   "/health",
						"methods": []any{"GET"},
						"owner":   "Operations",
						"reason":  "Expose repository health checks.",
					},
					map[string]any{
						"scope":   "/api/auth/*",
						"methods": []any{"POST", "GET"},
						"owner":   "Better Auth",
						"reason":  providerException["reason"],
					},
				},
			}
		}
	}

	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   repository,
		ProfileID:    "standard-typescript-monorepo",
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatalf("build Better Auth Plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("build Better Auth Plan returned result: %+v", outcome.Result)
	}

	httpDecision := planDecisionValue(t, outcome.Plan.Decisions, "http.contract")
	httpJSON, err := json.Marshal(httpDecision)
	if err != nil {
		t.Fatalf("marshal normalized HTTP Contract Decision: %v", err)
	}
	wantHTTP := `{"exceptions":[{"methods":["GET","POST"],"owner":"Better Auth","reason":"Session, OAuth redirect, callback, and related provider protocol routes require provider-owned GET and POST semantics.","scope":"/api/auth/*"},{"methods":["GET"],"owner":"Operations","reason":"Expose repository health checks.","scope":"/health"}],"mode":"Post-only"}`
	if string(httpJSON) != wantHTTP {
		t.Fatalf("normalized HTTP Contract Decision = %s, want %s", httpJSON, wantHTTP)
	}
	if stored := outcome.Plan.SetupManifest.Decisions["http.contract"].Value; !valuesEqual(stored, httpDecision) {
		t.Fatalf("Setup Manifest HTTP Contract Decision = %#v, want %#v", stored, httpDecision)
	}
	authDecision := planDecisionValue(t, outcome.Plan.Decisions, "auth.provider")
	if stored := outcome.Plan.SetupManifest.Decisions["auth.provider"].Value; !valuesEqual(stored, authDecision) {
		t.Fatalf("Setup Manifest auth provider = %#v, want %#v", stored, authDecision)
	}
}

func TestProjectDecisionRendering(t *testing.T) {
	repository := newProjectDecisionPlanRepository(t)
	if _, err := os.Stat(filepath.Join(repository, "docs", "adr")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("repository ADR fixture state = %v, want absent", err)
	}

	decisions := standardTypeScriptDecisions("make verify")
	for index := range decisions {
		if decisions[index].ID == "http.contract" {
			decisions[index].Value = map[string]any{"mode": "Post-only"}
		}
	}
	plan := buildProjectDecisionPlan(t, repository, decisions)
	catalog := mustEmbeddedCatalog(t)
	for _, fixture := range []struct {
		path      string
		assetPath string
	}{
		{
			path:      "docs/agents/domain.md",
			assetPath: "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md",
		},
		{
			path:      "docs/agents/backend.md",
			assetPath: "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md",
		},
	} {
		golden, ok := catalog.Asset(fixture.assetPath)
		if !ok {
			t.Fatalf("formatter fixture %q is missing", fixture.assetPath)
		}
		rendered := planPostimage(t, plan, fixture.path).Content
		if !bytes.Equal(rendered, golden.Data) {
			t.Errorf("rendered %q differs from formatter fixture", fixture.path)
		}
	}

	if _, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest); err != nil {
		t.Fatalf("apply project-decision Plan: %v", err)
	}
	result, err := ApplyPlan(context.Background(), repository, plan, plan.PlanDigest)
	if err != nil {
		t.Fatalf("empty project-decision reapply: %v", err)
	}
	if result.State != "verified" ||
		!strings.Contains(result.Message, "already applied") ||
		len(result.VerifiedPostimages) != len(plan.Postimages) {
		t.Fatalf("empty project-decision reapply result = %+v", result)
	}
}

func TestIdentifierGuidance(t *testing.T) {
	for _, layout := range []string{"single-context", "multi-context"} {
		t.Run("UUID version 7 scopes the rule for "+layout, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			for index := range decisions {
				if decisions[index].ID == "domain.layout" {
					decisions[index].Value = layout
				}
			}
			plan := buildProjectDecisionPlan(
				t,
				newProjectDecisionPlanRepository(t),
				decisions,
			)
			domain := string(planPostimage(t, plan, "docs/agents/domain.md").Content)
			for _, required := range []string{
				"Use UUID version 7 for new project-owned Internal Identifiers only.",
				"external provider identifiers",
				"protocol identifiers",
				"natural keys",
				"business codes",
				"source contracts",
			} {
				if !strings.Contains(domain, required) {
					t.Errorf("domain guidance is missing %q:\n%s", required, domain)
				}
			}
		})
	}

	t.Run("repository-defined guidance renders exactly", func(t *testing.T) {
		const guidance = "Use the repository's stable AccountId format."
		decisions := standardTypeScriptDecisions("make verify")
		for index := range decisions {
			if decisions[index].ID == "identifier.strategy" {
				decisions[index].Value = map[string]any{
					"kind":     "repository-defined",
					"guidance": guidance,
				}
			}
		}
		plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
		domain := string(planPostimage(t, plan, "docs/agents/domain.md").Content)
		if strings.Count(domain, guidance) != 1 {
			t.Fatalf("repository-defined guidance count = %d, want 1:\n%s",
				strings.Count(domain, guidance), domain)
		}
		if !strings.Contains(domain, "new project-owned Internal Identifiers only") ||
			!strings.Contains(domain, "source contracts") {
			t.Fatalf("repository-defined guidance lost identifier scope exceptions:\n%s", domain)
		}
	})
}

func TestBetterAuthGuidance(t *testing.T) {
	decisions := standardTypeScriptDecisions("make verify")
	for index := range decisions {
		if decisions[index].ID == "http.contract" {
			decisions[index].Value = map[string]any{
				"mode": "Post-only",
				"exceptions": []any{
					map[string]any{
						"scope":   "/health",
						"methods": []any{"GET"},
						"owner":   "Operations *team*",
						"reason":  "Expose [repository] health checks.",
					},
				},
			}
		}
	}
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	backend := string(planPostimage(t, plan, "docs/agents/backend.md").Content)
	for _, required := range []string{
		"Application HTTP mode: **Post-only**.",
		"**Better Auth**",
		"`GET` and `POST`",
		"`/api/auth/*`",
		"Session, OAuth redirect, callback, and related provider protocol routes require provider-owned GET and POST semantics.",
		"**Operations \\*team\\***",
		"`/health`",
		"Expose \\[repository\\] health checks.",
		"Better Auth owns the authentication protocol",
	} {
		if !strings.Contains(backend, required) {
			t.Errorf("backend guidance is missing %q:\n%s", required, backend)
		}
	}
	if auth, health := strings.Index(backend, "`/api/auth/*`"), strings.Index(backend, "`/health`"); auth < 0 || health < 0 || auth >= health {
		t.Fatalf("HTTP exceptions are not rendered in normalized order:\n%s", backend)
	}
}

func TestStructuredRenderSafety(t *testing.T) {
	tests := []struct {
		name   string
		mutate func([]DecisionValue)
	}{
		{
			name: "managed marker in repository-defined guidance",
			mutate: func(decisions []DecisionValue) {
				for index := range decisions {
					if decisions[index].ID == "identifier.strategy" {
						decisions[index].Value = map[string]any{
							"kind": "repository-defined",
							"guidance": "Use AccountId.\n\n" +
								"<!-- setup-context-driven:begin id=guide.domain version=0.0.1 -->",
						}
					}
				}
			},
		},
		{
			name: "non-canonical repository-defined guidance",
			mutate: func(decisions []DecisionValue) {
				for index := range decisions {
					if decisions[index].ID == "identifier.strategy" {
						decisions[index].Value = map[string]any{
							"kind":     "repository-defined",
							"guidance": " Use AccountId.",
						}
					}
				}
			},
		},
		{
			name: "template token in provider rationale",
			mutate: func(decisions []DecisionValue) {
				for index := range decisions {
					if decisions[index].ID != "auth.provider" {
						continue
					}
					provider := completeAuthProviderDecision()
					provider["routeException"].(map[string]any)["reason"] =
						"Preserve {{artifact.rules}} provider semantics."
					decisions[index].Value = provider
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			test.mutate(decisions)
			outcome, err := BuildPlan(context.Background(), PlanRequest{
				Repository:   newProjectDecisionPlanRepository(t),
				ProfileID:    "standard-typescript-monorepo",
				Decisions:    decisions,
				Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
			})
			if err == nil || outcome.Plan != nil || !strings.Contains(err.Error(), "structured render") {
				t.Fatalf("unsafe structured content produced Plan=%v error=%v", outcome.Plan != nil, err)
			}
		})
	}
}

func TestHTTPContractConflict(t *testing.T) {
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatalf("resolve Standard TypeScript Monorepo Profile: %v", err)
	}
	provider := completeAuthProviderDecision()
	providerReason := provider["routeException"].(map[string]any)["reason"]
	tests := []struct {
		name       string
		exceptions []any
	}{
		{
			name: "duplicate normalized methods",
			exceptions: []any{map[string]any{
				"scope": "/api/auth/*", "methods": []any{"GET", "get"},
				"owner": "Better Auth", "reason": providerReason,
			}},
		},
		{
			name: "unsupported method",
			exceptions: []any{map[string]any{
				"scope": "/api/auth/*", "methods": []any{"PATCH"},
				"owner": "Better Auth", "reason": providerReason,
			}},
		},
		{
			name: "missing rationale",
			exceptions: []any{map[string]any{
				"scope": "/api/auth/*", "methods": []any{"GET", "POST"},
				"owner": "Better Auth", "reason": " ",
			}},
		},
		{
			name: "provider projection conflicts with HTTP policy",
			exceptions: []any{map[string]any{
				"scope": "/api/auth/*", "methods": []any{"GET"},
				"owner": "Better Auth", "reason": providerReason,
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decisions := standardTypeScriptDecisions("make verify")
			for index := range decisions {
				if decisions[index].ID == "http.contract" {
					decisions[index].Value = map[string]any{
						"mode":       "Post-only",
						"exceptions": test.exceptions,
					}
				}
			}
			_, _, err := ResolveDecisionInput(profile, decisions, catalog)
			if err == nil ||
				!strings.Contains(err.Error(), "auth.provider") ||
				!strings.Contains(err.Error(), "http.contract") {
				t.Fatalf("ResolveDecisionInput() error = %v, want both decision IDs", err)
			}
		})
	}

	conflicting := standardTypeScriptDecisions("make verify")
	for index := range conflicting {
		if conflicting[index].ID == "http.contract" {
			conflicting[index].Value = map[string]any{
				"mode": "Post-only",
				"exceptions": []any{map[string]any{
					"scope": "/api/auth/*", "methods": []any{"GET"},
					"owner": "Better Auth", "reason": providerReason,
				}},
			}
		}
	}
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   newProjectDecisionPlanRepository(t),
		ProfileID:    profile.ID,
		Decisions:    conflicting,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err == nil || outcome.Plan != nil {
		t.Fatalf("conflicting provider/HTTP pair produced Plan=%v error=%v", outcome.Plan != nil, err)
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

func TestProfileDraftPlanIncludesCanonicalRepositoryProfile(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	request, draftTemplates := backendProfileDraftPlanRequest(t, repo)

	outcome, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("BuildPlan() profile draft error = %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("BuildPlan() profile draft result = %+v", outcome.Result)
	}
	plan := *outcome.Plan
	const profilePath = ".roundfix/baseline/profiles/backend-only.json"
	if _, err := os.Lstat(filepath.Join(repo, filepath.FromSlash(profilePath))); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("profile draft planning mutated repository: %v", err)
	}
	if plan.Profile.ID != "backend-only" ||
		plan.Profile.Source != ProfileSourceRepository ||
		plan.Profile.Path != profilePath ||
		slices.Equal(plan.Profile.Templates, draftTemplates) {
		t.Fatalf("resolved draft profile = %+v, input templates = %v", plan.Profile, draftTemplates)
	}

	profilePostimage := planPostimage(t, plan, profilePath)
	resolved, err := ParseCustomProfile(profilePostimage.Content, profilePath, mustEmbeddedCatalog(t))
	if err != nil {
		t.Fatalf("parse planned profile postimage: %v", err)
	}
	resolved.Path = profilePath
	if resolved.Digest != plan.Profile.Digest ||
		!slices.Equal(resolved.Templates, plan.Profile.Templates) {
		t.Fatalf("planned profile postimage = %+v, plan profile = %+v", resolved, plan.Profile)
	}
	assertProfilePlanLedgers(t, plan, profilePath, profilePostimage.ContentIdentity)
	if plan.SetupManifest.Profile != plan.Profile.ID ||
		plan.SetupManifest.ProfileDigest != plan.Profile.Digest {
		t.Fatalf("Setup Manifest profile identity = %q %q, want %q %q",
			plan.SetupManifest.Profile,
			plan.SetupManifest.ProfileDigest,
			plan.Profile.ID,
			plan.Profile.Digest,
		)
	}

	cloneParent := t.TempDir()
	clone := filepath.Join(cloneParent, "portable-profile-plan")
	runPlanGit(t, cloneParent, "clone", "--quiet", repo, clone)
	if err := ValidatePlanRepository(context.Background(), clone, plan); err != nil {
		t.Fatalf("ValidatePlanRepository() portable profile plan error = %v", err)
	}
}

func TestProfileDraftPlanAcceptsMatchingSourceBaselineWithoutTransition(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	catalog := mustEmbeddedCatalog(t)
	source, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	staleManifest := SetupManifest{
		SchemaVersion: ManifestSchema,
		Version:       ManifestVersion,
		Generator: ManifestGenerator{
			Skill:    "setup-context-driven",
			Version:  ManifestVersion,
			Baseline: "baseline." + source.ID + "-" + ManifestVersion,
		},
		Profile:          source.ID,
		ProfileDigest:    source.Digest,
		CatalogDigest:    "sha256:" + strings.Repeat("0", 64),
		Modules:          []string{},
		Decisions:        map[string]ManifestDecision{},
		ManagedArtifacts: []ManifestArtifact{},
		LocalSkills:      []string{},
		Verification:     []VerificationProjection{},
	}
	data, err := json.MarshalIndent(staleManifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeInspectionFile(t, repo, manifestPath, string(append(data, '\n')))
	commitInspectionRepository(t, repo, "seed source Baseline manifest")

	request, _ := backendProfileDraftPlanRequest(t, repo)
	outcome, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("BuildPlan() matching source Baseline error = %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("matching source Baseline blocked reviewed Profile adaptation: %+v", outcome.Result)
	}
	assertProfilePlanLedgers(
		t,
		*outcome.Plan,
		".roundfix/baseline/profiles/backend-only.json",
		planPostimage(t, *outcome.Plan, ".roundfix/baseline/profiles/backend-only.json").ContentIdentity,
	)
}

func TestProfileDraftPlanRejectsSimultaneousAndConflictingInputs(t *testing.T) {
	t.Run("simultaneous profile inputs", func(t *testing.T) {
		repo := newBackendProfileRepository(t, true)
		request, _ := backendProfileDraftPlanRequest(t, repo)
		request.ProfileID = "go-cli-tui"

		_, err := BuildPlan(context.Background(), request)
		if err == nil || !strings.Contains(err.Error(), "exactly one Baseline Profile") {
			t.Fatalf("BuildPlan() simultaneous profile error = %v", err)
		}
	})

	t.Run("conflicting canonical target", func(t *testing.T) {
		repo := newBackendProfileRepository(t, true)
		const profilePath = ".roundfix/baseline/profiles/backend-only.json"
		writeInspectionFile(t, repo, profilePath, "{}\n")
		commitInspectionRepository(t, repo, "seed conflicting profile target")
		before := snapshotVisibleTree(t, repo)
		request, _ := backendProfileDraftPlanRequest(t, repo)

		_, err := BuildPlan(context.Background(), request)
		if err == nil || !strings.Contains(err.Error(), "conflicts with existing repository bytes") {
			t.Fatalf("BuildPlan() conflicting profile target error = %v", err)
		}
		assertVisibleTree(t, repo, before)
	})
}

func TestProfileAdaptationCannotRemoveUniversalRequiredCapabilities(t *testing.T) {
	repo := newBackendProfileRepository(t, false)
	request, _ := backendProfileDraftPlanRequest(t, repo)

	outcome, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("BuildPlan() universal capability error = %v", err)
	}
	if outcome.Plan != nil || outcome.Result.State != "action_required" {
		t.Fatalf("BuildPlan() universal capability outcome = %+v", outcome)
	}
	for _, capabilityID := range []string{"capability.context7", "capability.exa"} {
		if !strings.Contains(outcome.Result.Message, capabilityID) {
			t.Errorf("universal capability outcome does not name %q: %+v", capabilityID, outcome.Result)
		}
	}
}

func TestProfileAdaptationRejectsInvalidDraftBoundaries(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ProfileDraftInput, *customProfileDocument)
		want   string
	}{
		{
			name: "unknown source Profile",
			mutate: func(input *ProfileDraftInput, _ *customProfileDocument) {
				input.SourceProfileID = "missing-profile"
			},
			want: "custom.profile.source.unknown",
		},
		{
			name: "module addition",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.Modules = []string{"core", "go", "typescript", "bun", "backend"}
			},
			want: "custom.profile.adaptation.modules.addition",
		},
		{
			name: "missing module dependency",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.Modules = []string{"core", "bun", "backend"}
			},
			want: "custom.profile.module.dependency.invalid",
		},
		{
			name: "decision addition",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.Decisions = append(document.Decisions, "runtime.backend")
			},
			want: "custom.profile.adaptation.decisions.addition",
		},
		{
			name: "universal capability override",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.Capabilities = append(document.Capabilities, "capability.context7")
			},
			want: "custom.profile.capability.universal",
		},
		{
			name: "unknown template",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.Templates = []string{"template.missing"}
			},
			want: "custom.profile.template.unknown",
		},
		{
			name: "stale catalog binding",
			mutate: func(_ *ProfileDraftInput, document *customProfileDocument) {
				document.CatalogSchema = "roundfix/baseline-catalog/stale"
			},
			want: "custom.profile.catalog-schema.invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newBackendProfileRepository(t, true)
			request, _ := backendProfileDraftPlanRequest(t, repo)
			var document customProfileDocument
			if err := json.Unmarshal(request.ProfileDraft.Document, &document); err != nil {
				t.Fatalf("decode base Profile draft: %v", err)
			}
			test.mutate(request.ProfileDraft, &document)
			data, err := json.MarshalIndent(document, "", "  ")
			if err != nil {
				t.Fatalf("marshal invalid Profile draft: %v", err)
			}
			request.ProfileDraft.Document = append(data, '\n')

			_, err = BuildPlan(context.Background(), request)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("BuildPlan() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestProfileDraftPlanRejectsUnsafeTargetParent(t *testing.T) {
	repo := newBackendProfileRepository(t, true)
	if err := os.Symlink(t.TempDir(), filepath.Join(repo, ".roundfix")); err != nil {
		t.Fatalf("create unsafe Profile parent: %v", err)
	}
	before := snapshotVisibleTree(t, repo)
	request, _ := backendProfileDraftPlanRequest(t, repo)

	_, err := BuildPlan(context.Background(), request)
	if !errors.Is(err, ErrUnsafeCustomProfilePath) {
		t.Fatalf("BuildPlan() unsafe Profile target error = %v, want ErrUnsafeCustomProfilePath", err)
	}
	assertVisibleTree(t, repo, before)
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

func backendProfileDraftPlanRequest(t *testing.T, repo string) (PlanRequest, []string) {
	t.Helper()
	document := customProfileDocument{
		SchemaVersion: CustomProfileSchemaVersion,
		CatalogSchema: CatalogSchemaVersion(),
		ID:            "backend-only",
		Modules:       []string{"core", "typescript", "bun", "backend"},
		Decisions:     []string{"language.generated", "verification.gate"},
		Capabilities: []string{
			"capability.stack.bun",
			"capability.stack.hono",
			"capability.stack.typescript",
			"capability.stack.zod",
		},
		Templates: []string{"template.root.core"},
		Values:    map[string]any{},
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatalf("marshal profile draft: %v", err)
	}
	data = append(data, '\n')
	return PlanRequest{
		Repository: repo,
		ProfileDraft: &ProfileDraftInput{
			SourceProfileID: "standard-typescript-monorepo",
			Document:        data,
		},
		Decisions: []DecisionValue{
			{ID: "language.generated", Value: "English"},
			{ID: "verification.gate", Value: "make verify"},
		},
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	}, append([]string(nil), document.Templates...)
}

func newBackendProfileRepository(t *testing.T, includeUniversalSkills bool) string {
	t.Helper()
	repo := newInspectionRepository(t)
	if includeUniversalSkills {
		writeInspectionFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
		writeInspectionFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	}
	writeInspectionFile(t, repo, "Makefile", "verify:\n\t@true\n")
	writeInspectionFile(t, repo, "package.json", `{
  "packageManager": "bun@1.0.0",
  "dependencies": {
    "hono": "1.0.0",
    "typescript": "1.0.0",
    "zod": "1.0.0"
  }
}
`)
	commitInspectionRepository(t, repo, "seed backend profile repository")
	return repo
}

func mustEmbeddedCatalog(t *testing.T) *Catalog {
	t.Helper()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedCatalog() error = %v", err)
	}
	return catalog
}

func assertProfilePlanLedgers(t *testing.T, plan PlanDocument, profilePath, identity string) {
	t.Helper()
	entryID := "profile:" + plan.Profile.ID
	var managed bool
	for _, entry := range plan.ManagedEntries {
		if entry.ID == entryID &&
			entry.Path == profilePath &&
			entry.Kind == "profile" &&
			entry.AfterIdentity == identity &&
			entry.ContentIdentity == identity {
			managed = true
			break
		}
	}
	if !managed {
		t.Fatalf("profile managed-entry ledger is missing %q at %q", entryID, profilePath)
	}
	for _, change := range plan.FileChanges {
		if change.Path == profilePath && slices.Contains(change.ManagedEntries, entryID) {
			return
		}
	}
	t.Fatalf("profile file-change ledger is missing %q at %q", entryID, profilePath)
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

func newProjectDecisionPlanRepository(t *testing.T) string {
	t.Helper()
	repository := newAlignedTypeScriptRepository(t)
	runPlanGit(t, repository, "init", "-q")
	runPlanGit(t, repository, "config", "user.email", "fixture@example.invalid")
	runPlanGit(t, repository, "config", "user.name", "Fixture Test")
	runPlanGit(t, repository, "config", "commit.gpgsign", "false")
	runPlanGit(t, repository, "add", ".")
	runPlanGit(t, repository, "commit", "-qm", "seed project decisions")
	return repository
}

func buildProjectDecisionPlan(
	t *testing.T,
	repository string,
	decisions []DecisionValue,
) PlanDocument {
	t.Helper()
	outcome, err := BuildPlan(context.Background(), PlanRequest{
		Repository:   repository,
		ProfileID:    "standard-typescript-monorepo",
		Decisions:    decisions,
		Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
	})
	if err != nil {
		t.Fatalf("build project-decision Plan: %v", err)
	}
	if outcome.Plan == nil {
		t.Fatalf("build project-decision Plan returned result: %+v", outcome.Result)
	}
	return *outcome.Plan
}

func planDecisionValue(t *testing.T, decisions []DecisionValue, id string) any {
	t.Helper()
	for _, decision := range decisions {
		if decision.ID == id {
			return decision.Value
		}
	}
	t.Fatalf("Plan has no decision %q", id)
	return nil
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
