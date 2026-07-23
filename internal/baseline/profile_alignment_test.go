// Suite: Baseline Profile alignment
// Invariant: one selected profile turns bounded local facts and explicit answers into deterministic blocking or advisory decisions without executing repository commands.
// Boundary IN: profile resolution, capability evidence, HTTP/PostgreSQL projections, and Verification declaration validation.
// Boundary OUT: portable Plan serialization, CLI prompting/rendering, repository mutation, and network access.

package baseline

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestProfileAlignmentResolvesExactlyOneProfile(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve built-in profile alignment: %v", err)
	}
	if alignment.Profile.ID != "standard-typescript-monorepo" ||
		alignment.Profile.Source != ProfileSourceBuiltIn ||
		!alignment.Ready ||
		alignment.State != ProfileAlignmentReady {
		t.Fatalf("resolved alignment = %+v", alignment)
	}

	custom, err := InitCustomProfile(repository, "repository-typescript", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatalf("initialize repository-owned profile: %v", err)
	}
	repositoryAlignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: custom.Profile.ID,
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve repository-owned profile alignment: %v", err)
	}
	if repositoryAlignment.Profile.Source != ProfileSourceRepository || !repositoryAlignment.Ready {
		t.Fatalf("repository-owned alignment = %+v", repositoryAlignment)
	}

	tests := []struct {
		name      string
		profileID string
		want      string
	}{
		{name: "missing selection", profileID: "", want: "exactly one Baseline Profile"},
		{name: "unknown selection", profileID: "unknown-profile", want: "unknown-profile"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
				ProfileID: tt.profileID,
				Decisions: standardTypeScriptDecisions("make verify"),
			}, catalog)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("ResolveProfileAlignment() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestRequiredDivergencePreventsReadyPlan(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(false, true))

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve required divergence: %v", err)
	}
	if alignment.Ready || alignment.State != ProfileAlignmentActionRequired {
		t.Fatalf("required divergence alignment = %+v", alignment)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.stack.better-auth")
	if !ok || !divergence.Blocking || divergence.Requirement != CapabilityRequired ||
		divergence.Code != "capability.required.missing" {
		t.Fatalf("required divergence = %+v, found=%v", divergence, ok)
	}

	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(true, true))
	resolved, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve remediated alignment: %v", err)
	}
	if !resolved.Ready {
		t.Fatalf("remediated alignment remains blocked: %+v", resolved.Divergences)
	}
}

func TestProfileAlignmentAdvisoryDivergenceNeverBlocksOrInfersPolicy(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve advisory divergence: %v", err)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.firecrawl")
	if !ok || divergence.Blocking || divergence.Requirement != CapabilityRecommended {
		t.Fatalf("recommended divergence = %+v, found=%v", divergence, ok)
	}
	if !alignment.Ready {
		t.Fatalf("advisory divergence blocked alignment: %+v", alignment.Divergences)
	}
	for _, decision := range alignment.Decisions {
		if decision.ID == "capability.firecrawl" {
			t.Fatalf("advisory capability became inferred policy: %+v", alignment.Decisions)
		}
	}
}

func TestProfileAlignmentCapabilityEvidenceRanking(t *testing.T) {
	tests := []struct {
		name        string
		requirement CapabilityRequirement
		minimum     CapabilityEvidenceStrength
		evidence    CapabilityEvidence
		wantStatus  CapabilityStatus
		wantBlock   bool
	}{
		{
			name:        "declared evidence satisfies declared requirement",
			requirement: CapabilityRequired,
			minimum:     CapabilityEvidenceDeclared,
			evidence: CapabilityEvidence{
				Status:   CapabilityEvidencePresent,
				Kind:     CapabilityEvidenceDeclaredFile,
				Strength: CapabilityEvidenceDeclared,
			},
			wantStatus: CapabilitySatisfied,
		},
		{
			name:        "declared evidence is insufficient for verified requirement",
			requirement: CapabilityRequired,
			minimum:     CapabilityEvidenceVerified,
			evidence: CapabilityEvidence{
				Status:   CapabilityEvidencePresent,
				Kind:     CapabilityEvidenceInstalledSkill,
				Strength: CapabilityEvidenceDeclared,
			},
			wantStatus: CapabilityInsufficient,
			wantBlock:  true,
		},
		{
			name:        "missing recommended evidence remains advisory",
			requirement: CapabilityRecommended,
			minimum:     CapabilityEvidenceVerified,
			evidence: CapabilityEvidence{
				Status:   CapabilityEvidenceAbsent,
				Kind:     CapabilityEvidenceInstalledSkill,
				Strength: CapabilityEvidenceNone,
			},
			wantStatus: CapabilityMissing,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outcome := evaluateCapability(RepositoryCapability{
				ID:              "capability.test",
				Title:           "Test capability",
				Requirement:     tt.requirement,
				EvidenceKind:    tt.evidence.Kind,
				MinimumEvidence: tt.minimum,
			}, []CapabilityEvidence{tt.evidence})
			if outcome.Status != tt.wantStatus || outcome.Blocking != tt.wantBlock {
				t.Fatalf("capability outcome = %+v, want status=%s blocking=%v", outcome, tt.wantStatus, tt.wantBlock)
			}
		})
	}
}

func TestHTTPRouteCandidatesContainFactsWithoutNormativeClause(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "packages/backend/src/infra/controllers/http/app.ts", `
const app = new Hono()
app.get("/health", health)
app.post("/api/orders", createOrder)
app.route("/auth", auth)
`)

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve HTTP route candidates: %v", err)
	}
	want := []HTTPRouteCandidate{
		{
			SourcePath:   "packages/backend/src/infra/controllers/http/app.ts",
			SourceDigest: alignment.HTTPCandidates[0].SourceDigest,
			Scope:        "/api/orders",
			Methods:      []string{"POST"},
		},
		{
			SourcePath:   "packages/backend/src/infra/controllers/http/app.ts",
			SourceDigest: alignment.HTTPCandidates[0].SourceDigest,
			Scope:        "/auth",
			Methods:      nil,
		},
		{
			SourcePath:   "packages/backend/src/infra/controllers/http/app.ts",
			SourceDigest: alignment.HTTPCandidates[0].SourceDigest,
			Scope:        "/health",
			Methods:      []string{"GET"},
		},
	}
	if !reflect.DeepEqual(alignment.HTTPCandidates, want) {
		t.Fatalf("HTTP candidates = %+v, want %+v", alignment.HTTPCandidates, want)
	}
	if !strings.HasPrefix(alignment.HTTPCandidates[0].SourceDigest, "sha256:") {
		t.Fatalf("HTTP source digest = %q", alignment.HTTPCandidates[0].SourceDigest)
	}
	data, err := json.Marshal(alignment.HTTPCandidates)
	if err != nil {
		t.Fatalf("marshal HTTP candidates: %v", err)
	}
	for _, forbidden := range []string{"mode", "owner", "reason", "rationale", "normative"} {
		if strings.Contains(strings.ToLower(string(data)), forbidden) {
			t.Fatalf("HTTP candidates inferred %q policy: %s", forbidden, data)
		}
	}
}

func TestPostgreSQLEvidenceSeparatesImplementationAndContract(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	if err := os.Remove(filepath.Join(repository, "DATABASE.md")); err != nil {
		t.Fatalf("remove PostgreSQL contract: %v", err)
	}
	writeProfileAlignmentFile(t, repository, "docker-compose.yml", "services:\n  db:\n    image: postgres:18-alpine\n")

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve PostgreSQL evidence: %v", err)
	}
	if alignment.PostgreSQL.Contract != nil || len(alignment.PostgreSQL.Implementation) == 0 {
		t.Fatalf("PostgreSQL evidence = %+v", alignment.PostgreSQL)
	}
	wantPaths := []string{
		"DATABASE.md",
		"docs/architecture/database.json",
		"docs/architecture/postgresql.json",
	}
	if !slices.Equal(alignment.PostgreSQL.AcceptedContractPaths, wantPaths) {
		t.Fatalf("accepted PostgreSQL contract paths = %v, want %v", alignment.PostgreSQL.AcceptedContractPaths, wantPaths)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.stack.postgresql")
	if !ok || !divergence.Blocking || divergence.Code != "capability.contract.missing" ||
		!strings.Contains(divergence.Message, "implementation evidence") ||
		!strings.Contains(divergence.NextAction, "DATABASE.md") {
		t.Fatalf("PostgreSQL divergence = %+v, found=%v", divergence, ok)
	}

	writeProfileAlignmentFile(t, repository, "DATABASE.md", "# Database\n\nPostgreSQL is the repository database contract.\n")
	resolved, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve PostgreSQL contract: %v", err)
	}
	if resolved.PostgreSQL.Contract == nil || !resolved.Ready {
		t.Fatalf("PostgreSQL contract did not satisfy alignment: %+v", resolved)
	}
}

func TestExecutableVerificationCommandRequiresLocalDeclaration(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(true, false))
	if err := os.Remove(filepath.Join(repository, "Makefile")); err != nil {
		t.Fatalf("remove Makefile: %v", err)
	}

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve undeclared Verification commands: %v", err)
	}
	format, ok := findVerificationProjection(alignment.Verification, "verification.format")
	if !ok || format.RepositoryExecutable || format.Classification != VerificationProfileExpectation {
		t.Fatalf("format projection = %+v, found=%v", format, ok)
	}
	workspace, ok := findVerificationProjection(alignment.Verification, "verification.workspace")
	if !ok || workspace.RepositoryExecutable {
		t.Fatalf("workspace projection = %+v, found=%v", workspace, ok)
	}
	gate, ok := findVerificationProjection(alignment.Verification, "verification.gate")
	if !ok || gate.RepositoryExecutable || gate.Classification != VerificationRepositoryCommand {
		t.Fatalf("selected Verification projection = %+v, found=%v", gate, ok)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "verification.gate")
	if !ok || !divergence.Blocking || divergence.Code != "verification.command.undeclared" {
		t.Fatalf("selected Verification divergence = %+v, found=%v", divergence, ok)
	}

	writeProfileAlignmentFile(t, repository, "Makefile", "verify:\n\t@true\n")
	resolved, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve declared Verification command: %v", err)
	}
	gate, ok = findVerificationProjection(resolved.Verification, "verification.gate")
	if !ok || !gate.RepositoryExecutable || gate.DeclarationPath != "Makefile" ||
		!strings.HasPrefix(gate.DeclarationDigest, "sha256:") {
		t.Fatalf("declared Verification projection = %+v, found=%v", gate, ok)
	}
	if !resolved.Ready {
		t.Fatalf("declared Verification command remains blocked: %+v", resolved.Divergences)
	}
}

func TestProfileAlignmentEquivalentNormalizedDecisions(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	interactive := standardTypeScriptDecisions("make verify")
	fileBased := append([]DecisionValue(nil), interactive...)
	slices.Reverse(fileBased)

	first, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: interactive,
	}, catalog)
	if err != nil {
		t.Fatalf("resolve interactive answers: %v", err)
	}
	second, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: fileBased,
	}, catalog)
	if err != nil {
		t.Fatalf("resolve file answers: %v", err)
	}
	firstJSON, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("marshal first alignment: %v", err)
	}
	secondJSON, err := json.Marshal(second)
	if err != nil {
		t.Fatalf("marshal second alignment: %v", err)
	}
	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("equivalent answers normalized differently:\nfirst=%s\nsecond=%s", firstJSON, secondJSON)
	}
}

func TestCapabilityAuditNoExecution(t *testing.T) {
	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	bin := t.TempDir()
	marker := filepath.Join(t.TempDir(), "executed")
	for _, name := range []string{"docker", "rg", "rtk"} {
		writeProfileAlignmentExecutable(t, bin, name, "#!/bin/sh\nprintf executed > \""+marker+"\"\n")
	}
	t.Setenv("PATH", bin)

	before := snapshotInspectionTree(t, repository)
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve command-free capability audit: %v", err)
	}
	if !alignment.Ready {
		t.Fatalf("command-free capability audit unexpectedly blocked: %+v", alignment.Divergences)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("capability audit executed a discovered command: %v", err)
	}
	after := snapshotInspectionTree(t, repository)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("capability audit changed repository bytes:\nbefore=%+v\nafter=%+v", before, after)
	}
}

func loadProfileAlignmentCatalog(t *testing.T) *Catalog {
	t.Helper()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	return catalog
}

func newAlignedTypeScriptRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeProfileAlignmentFile(t, root, "package.json", typeScriptPackageJSON(true, true))
	writeProfileAlignmentFile(t, root, "packages/frontend/package.json", `{"name":"frontend"}`)
	writeProfileAlignmentFile(t, root, "packages/backend/package.json", `{"name":"backend","dependencies":{"postgres":"latest","drizzle-orm":"latest"}}`)
	writeProfileAlignmentFile(t, root, "DATABASE.md", "# Database\n\nPostgreSQL is the repository database contract.\n")
	writeProfileAlignmentFile(t, root, "Makefile", "verify:\n\t@true\n")
	writeProfileAlignmentFile(t, root, ".agents/skills/context7/SKILL.md", "# Context7\n")
	writeProfileAlignmentFile(t, root, ".agents/skills/exa-web-search/SKILL.md", "# Exa\n")
	return root
}

func typeScriptPackageJSON(includeBetterAuth, includeProfileScripts bool) string {
	dependencies := []string{
		`"@logtape/logtape":"latest"`,
		`"@tanstack/react-query":"latest"`,
		`"@tanstack/react-router":"latest"`,
		`"drizzle-orm":"latest"`,
		`"hono":"latest"`,
		`"oxfmt":"latest"`,
		`"oxlint":"latest"`,
		`"postgres":"latest"`,
		`"react":"latest"`,
		`"shadcn":"latest"`,
		`"tailwindcss":"latest"`,
		`"turbo":"latest"`,
		`"typescript":"latest"`,
		`"vite":"latest"`,
		`"vitest":"latest"`,
		`"zod":"latest"`,
	}
	if includeBetterAuth {
		dependencies = append(dependencies, `"better-auth":"latest"`)
	}
	scripts := `"fmt":"oxfmt ."`
	if includeProfileScripts {
		scripts = `"format":"oxfmt .","lint":"oxlint","test":"vitest","build":"turbo build","verify":"bun run lint"`
	}
	return `{"name":"root","packageManager":"bun@1.3.0","scripts":{` + scripts + `},"dependencies":{` + strings.Join(dependencies, ",") + `}}`
}

func standardTypeScriptDecisions(verification string) []DecisionValue {
	return []DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: verification},
		{ID: "http.contract", Value: map[string]any{"mode": "REST"}},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: false},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: false},
	}
}

func findProfileDivergence(divergences []ProfileDivergence, id string) (ProfileDivergence, bool) {
	for _, divergence := range divergences {
		if divergence.ID == id {
			return divergence, true
		}
	}
	return ProfileDivergence{}, false
}

func findVerificationProjection(projections []VerificationProjection, id string) (VerificationProjection, bool) {
	for _, projection := range projections {
		if projection.ID == id {
			return projection, true
		}
	}
	return VerificationProjection{}, false
}

func writeProfileAlignmentFile(t *testing.T, root, relative, content string) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("create %s parent: %v", relative, err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}

func writeProfileAlignmentExecutable(t *testing.T, root, name, content string) {
	t.Helper()
	target := filepath.Join(root, name)
	if err := os.WriteFile(target, []byte(content), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", name, err)
	}
}
