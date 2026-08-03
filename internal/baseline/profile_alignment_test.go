// Suite: Baseline Profile alignment
// Invariant: one selected profile turns bounded local facts and explicit answers into deterministic blocking or advisory decisions without executing repository commands.
// Boundary IN: profile resolution, capability evidence, HTTP/PostgreSQL projections, and Verification declaration validation.
// Boundary OUT: portable Plan serialization, CLI prompting/rendering, repository mutation, and network access.

package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestProfileAlignmentResolvesExactlyOneProfile(t *testing.T) {
	t.Parallel()

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

	t.Run("relative search directory is rejected without reading the process working directory", func(t *testing.T) {
		result := resolveExecutableCandidate("probe", []string{"relative-bin"})
		if result.Candidate != filepath.Join("relative-bin", "probe") ||
			result.Reason != executableProbeReasonRelativePath {
			t.Fatalf("relative executable directory result = %+v", result)
		}
	})
}

func TestCapabilityRecheckMatchesFullPlan(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	executableDirectories := []string{t.TempDir()}

	fullPlanAlignment, err := ResolveProfileAlignment(
		context.Background(),
		repository,
		ProfileAlignmentRequest{
			ProfileID:             "standard-typescript-monorepo",
			Decisions:             standardTypeScriptDecisions("make verify"),
			ExecutableDirectories: executableDirectories,
		},
		catalog,
	)
	if err != nil {
		t.Fatalf("resolve full-plan alignment: %v", err)
	}

	recheck, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
		Repository:            repository,
		ProfileID:             "standard-typescript-monorepo",
		ExecutableDirectories: executableDirectories,
	})
	if err != nil {
		t.Fatalf("re-check capabilities: %v", err)
	}
	if !reflect.DeepEqual(recheck.Capabilities, fullPlanAlignment.Capabilities) {
		t.Fatalf("re-check capabilities differ from full-plan alignment:\nre-check=%#v\nfull-plan=%#v",
			recheck.Capabilities, fullPlanAlignment.Capabilities)
	}

	capabilityIDs := make(map[string]struct{}, len(fullPlanAlignment.Capabilities))
	for _, outcome := range fullPlanAlignment.Capabilities {
		capabilityIDs[outcome.ID] = struct{}{}
	}
	fullPlanDivergences := make([]ProfileDivergence, 0, len(fullPlanAlignment.Divergences))
	for _, divergence := range fullPlanAlignment.Divergences {
		if _, capability := capabilityIDs[divergence.ID]; capability {
			fullPlanDivergences = append(fullPlanDivergences, divergence)
		}
	}
	if !reflect.DeepEqual(recheck.Divergences, fullPlanDivergences) {
		t.Fatalf("re-check divergences differ from full-plan alignment:\nre-check=%#v\nfull-plan=%#v",
			recheck.Divergences, fullPlanDivergences)
	}
	if got, want := RenderProfileDivergences(recheck.Divergences), RenderProfileDivergences(fullPlanDivergences); got != want {
		t.Fatalf("re-check rendering differs from full-plan rendering:\nre-check=%s\nfull-plan=%s", got, want)
	}
}

func TestCapabilityRecheck(t *testing.T) {
	t.Parallel()

	t.Run("names a missing Profile", func(t *testing.T) {
		repository := t.TempDir()
		_, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
			Repository: repository,
		})
		if !errors.Is(err, ErrNoResolvableProfile) {
			t.Fatalf("RecheckCapabilities() error = %v, want ErrNoResolvableProfile", err)
		}
	})

	t.Run("resolves the current Setup Manifest Profile without decisions", func(t *testing.T) {
		catalog := loadProfileAlignmentCatalog(t)
		repository := newAlignedTypeScriptRepository(t)
		profile, err := ResolveProfile(repository, "standard-typescript-monorepo", catalog)
		if err != nil {
			t.Fatalf("resolve manifest Profile: %v", err)
		}
		manifest := buildSetupManifest(
			catalog,
			profile,
			nil,
			profile.Modules,
			nil,
			[]VerificationProjection{},
		)
		data, err := marshalSetupManifestBytes(manifest)
		if err != nil {
			t.Fatalf("marshal Setup Manifest: %v", err)
		}
		writeProfileAlignmentFile(t, repository, manifestPath, string(data))

		result, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
			Repository: repository,
		})
		if err != nil {
			t.Fatalf("re-check current Setup Manifest Profile: %v", err)
		}
		if result.Profile == nil || result.Profile.ID != profile.ID || result.Profile.Digest != profile.Digest {
			t.Fatalf("re-check Profile = %#v, want %q %q", result.Profile, profile.ID, profile.Digest)
		}
	})

	t.Run("rejects an invalid current Setup Manifest", func(t *testing.T) {
		repository := t.TempDir()
		writeProfileAlignmentFile(t, repository, manifestPath, "{}")

		_, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
			Repository: repository,
		})
		if !errors.Is(err, ErrNoResolvableProfile) ||
			!strings.Contains(err.Error(), "current Setup Manifest is invalid") {
			t.Fatalf("RecheckCapabilities() error = %v, want invalid manifest ErrNoResolvableProfile", err)
		}
	})

	t.Run("wraps an explicit Profile resolution failure", func(t *testing.T) {
		repository := t.TempDir()
		_, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
			Repository: repository,
			ProfileID:  "missing-profile",
		})
		if !errors.Is(err, ErrNoResolvableProfile) || !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("RecheckCapabilities() error = %v, want ErrNoResolvableProfile and fs.ErrNotExist", err)
		}
	})

	t.Run("normalizes omitted PostgreSQL evidence to arrays", func(t *testing.T) {
		result, err := RecheckCapabilities(context.Background(), CapabilityRecheckRequest{
			Repository: t.TempDir(),
			ProfileID:  "go-cli-tui",
		})
		if err != nil {
			t.Fatalf("RecheckCapabilities() error = %v", err)
		}
		if result.PostgreSQL.AcceptedContractPaths == nil || result.PostgreSQL.Implementation == nil {
			t.Fatalf("PostgreSQL evidence contains nil slices: %#v", result.PostgreSQL)
		}
		encoded, err := json.Marshal(result.PostgreSQL)
		if err != nil {
			t.Fatalf("marshal PostgreSQL evidence: %v", err)
		}
		if !strings.Contains(string(encoded), `"acceptedContractPaths":[]`) ||
			!strings.Contains(string(encoded), `"implementation":[]`) {
			t.Fatalf("PostgreSQL evidence JSON = %s, want empty arrays", encoded)
		}
	})
}

func TestRequiredDivergencePreventsReadyPlan(t *testing.T) {
	t.Parallel()

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

func TestDivergenceCarriesProbeEvidence(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(false, true))
	executableDirectories := []string{t.TempDir()}

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: executableDirectories,
	}, catalog)
	if err != nil {
		t.Fatalf("resolve alignment with capability divergences: %v", err)
	}
	capabilities, err := resolvedProfileCapabilities(alignment.Profile, catalog)
	if err != nil {
		t.Fatalf("resolve evaluated capability definitions: %v", err)
	}
	definitions := make(map[string]RepositoryCapability, len(capabilities))
	for _, capability := range capabilities {
		definitions[capability.ID] = capability
	}

	for _, outcome := range alignment.Capabilities {
		if outcome.Status == CapabilitySatisfied {
			continue
		}
		divergence, ok := findProfileDivergence(alignment.Divergences, outcome.ID)
		if !ok {
			t.Errorf("unsatisfied capability %q has no divergence", outcome.ID)
			continue
		}
		definition, ok := definitions[outcome.ID]
		if !ok {
			t.Errorf("unsatisfied capability %q has no evaluated definition", outcome.ID)
			continue
		}
		wantProbe, err := json.Marshal(definition.Probe)
		if err != nil {
			t.Fatalf("marshal evaluated probe for %q: %v", outcome.ID, err)
		}
		gotProbe, err := json.Marshal(divergence.Probe)
		if err != nil {
			t.Fatalf("marshal divergence probe for %q: %v", outcome.ID, err)
		}
		if !slices.Equal(gotProbe, wantProbe) {
			t.Errorf("divergence probe for %q = %s, want evaluated bytes %s", outcome.ID, gotProbe, wantProbe)
		}
		if !reflect.DeepEqual(divergence.Evidence, outcome.Evidence) {
			t.Errorf("divergence evidence for %q = %+v, want verdict evidence %+v", outcome.ID, divergence.Evidence, outcome.Evidence)
		}
	}

	legacy, err := json.Marshal(ProfileDivergence{
		Code:        "profile.decision.required",
		ID:          "language.generated",
		Requirement: CapabilityRequired,
		Blocking:    true,
		Message:     "answer required",
		NextAction:  "answer the decision",
	})
	if err != nil {
		t.Fatalf("marshal legacy divergence: %v", err)
	}
	const wantLegacy = `{"code":"profile.decision.required","id":"language.generated","requirement":"required","blocking":true,"message":"answer required","nextAction":"answer the decision"}`
	if string(legacy) != wantLegacy {
		t.Fatalf("legacy divergence JSON = %s, want %s", legacy, wantLegacy)
	}
}

func TestProfileAlignmentAdvisoryDivergenceNeverBlocksOrInfersPolicy(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

func TestHTTPRouteCandidateLimitIgnoresUnrelatedSourceFiles(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	for index := range maxHTTPSourceFiles + 1 {
		writeProfileAlignmentFile(
			t,
			repository,
			fmt.Sprintf("packages/backend/src/domain/unrelated-%03d.ts", index),
			"export const value = true\n",
		)
	}
	writeProfileAlignmentFile(
		t,
		repository,
		"packages/backend/src/infra/controllers/http/app.ts",
		`app.post("/api/orders", createOrder)`,
	)

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve HTTP route candidates in a large codebase: %v", err)
	}
	if got := len(alignment.HTTPCandidates); got != 1 {
		t.Fatalf("HTTP candidate count = %d, want 1", got)
	}
}

func TestHTTPRouteCandidateLimitRejectsTooManyRelevantSources(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	for index := range maxHTTPSourceFiles + 1 {
		writeProfileAlignmentFile(
			t,
			repository,
			fmt.Sprintf("packages/backend/src/infra/controllers/routes/route-%03d.ts", index),
			fmt.Sprintf(`app.get("/route-%03d", handler)`, index),
		)
	}

	_, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err == nil || !strings.Contains(err.Error(), "HTTP candidate source count exceeds 256") {
		t.Fatalf("HTTP relevant-source limit error = %v", err)
	}
}

func TestPostgreSQLEvidenceSeparatesImplementationAndContract(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	if !ok ||
		!format.RepositoryExecutable ||
		format.Command != "bun run fmt" ||
		format.Classification != VerificationRepositoryCommand {
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

func TestPortableVerificationRoleMapping(t *testing.T) {
	t.Parallel()

	const (
		role                = "workspace"
		portableCommand     = "bun run verify"
		declaredCommand     = "make verify"
		undeclaredCommand   = "make missing"
		legacyMessage       = "portable workspace role expects \"bun run verify\", but no matching local declaration exists"
		legacyNextAction    = "map the portable role to a declared repository command or keep it as a profile expectation"
		mappedAbsentMessage = "portable workspace role maps to repository command \"make missing\", but no matching local declaration exists"
		mappedAbsentAction  = "declare the mapped repository command or map the portable role to another declared command"
		executionMarker     = "verification-role-executed"
	)

	newRepository := func(t *testing.T) string {
		t.Helper()
		repository := newAlignedTypeScriptRepository(t)
		writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(true, false))
		writeProfileAlignmentFile(t, repository, "Makefile", "verify:\n\t@touch "+executionMarker+"\n")
		return repository
	}
	resolve := func(t *testing.T, repository string, mappings map[string]string) ProfileAlignment {
		t.Helper()
		alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
			ProfileID:                "standard-typescript-monorepo",
			Decisions:                standardTypeScriptDecisions(declaredCommand),
			VerificationRoleMappings: mappings,
		}, loadProfileAlignmentCatalog(t))
		if err != nil {
			t.Fatalf("resolve portable Verification role mapping: %v", err)
		}
		return alignment
	}

	t.Run("mapped role is satisfied by a present declaration without execution", func(t *testing.T) {
		repository := newRepository(t)
		alignment := resolve(t, repository, map[string]string{role: declaredCommand})

		projection, ok := findVerificationProjection(alignment.Verification, "verification.workspace")
		if !ok || projection.Command != declaredCommand ||
			projection.SatisfiedByCommand != declaredCommand ||
			projection.Classification != VerificationRepositoryCommand ||
			!projection.RepositoryExecutable || projection.DeclarationPath != "Makefile" ||
			!strings.HasPrefix(projection.DeclarationDigest, "sha256:") {
			t.Fatalf("mapped workspace projection = %+v, found=%v", projection, ok)
		}
		if divergence, exists := findProfileDivergence(alignment.Divergences, "verification.workspace"); exists {
			t.Fatalf("mapped workspace role remained divergent: %+v", divergence)
		}
		if _, err := os.Stat(filepath.Join(repository, executionMarker)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("mapped repository command was executed: %v", err)
		}
	})

	t.Run("unmapped role keeps its legacy divergence", func(t *testing.T) {
		alignment := resolve(t, newRepository(t), nil)

		projection, ok := findVerificationProjection(alignment.Verification, "verification.workspace")
		if !ok || projection.Command != portableCommand ||
			projection.SatisfiedByCommand != "" ||
			projection.Classification != VerificationProfileExpectation ||
			projection.RepositoryExecutable {
			t.Fatalf("unmapped workspace projection = %+v, found=%v", projection, ok)
		}
		divergence, ok := findProfileDivergence(alignment.Divergences, "verification.workspace")
		if !ok || divergence.Code != "verification.profile-expectation.unresolved" ||
			divergence.Requirement != CapabilityRecommended || divergence.Blocking ||
			divergence.Message != legacyMessage || divergence.NextAction != legacyNextAction {
			t.Fatalf("unmapped workspace divergence = %+v, found=%v", divergence, ok)
		}
	})

	t.Run("mapped role rejects an absent declaration with a reason", func(t *testing.T) {
		alignment := resolve(t, newRepository(t), map[string]string{role: undeclaredCommand})

		projection, ok := findVerificationProjection(alignment.Verification, "verification.workspace")
		if !ok || projection.Command != undeclaredCommand ||
			projection.SatisfiedByCommand != "" ||
			projection.Classification != VerificationRepositoryCommand ||
			projection.RepositoryExecutable {
			t.Fatalf("undeclared mapped workspace projection = %+v, found=%v", projection, ok)
		}
		divergence, ok := findProfileDivergence(alignment.Divergences, "verification.workspace")
		if !ok || divergence.Code != "verification.role-mapping.undeclared" ||
			divergence.Requirement != CapabilityRecommended || divergence.Blocking ||
			divergence.Message != mappedAbsentMessage || divergence.NextAction != mappedAbsentAction {
			t.Fatalf("undeclared mapped workspace divergence = %+v, found=%v", divergence, ok)
		}
	})

	t.Run("mapped role rejects an empty command", func(t *testing.T) {
		repository := newRepository(t)
		_, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
			ProfileID:                "standard-typescript-monorepo",
			Decisions:                standardTypeScriptDecisions(declaredCommand),
			VerificationRoleMappings: map[string]string{role: "   "},
		}, loadProfileAlignmentCatalog(t))
		if err == nil || !strings.Contains(err.Error(), `mapped repository command for role "workspace" is empty`) {
			t.Fatalf("empty role mapping error = %v", err)
		}
	})
}

func TestProfileAlignmentDiscoversDeclaredRepositoryFormatter(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(true, false))

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID: "standard-typescript-monorepo",
		Decisions: standardTypeScriptDecisions("make verify"),
	}, catalog)
	if err != nil {
		t.Fatalf("resolve repository formatter: %v", err)
	}
	format, ok := findVerificationProjection(alignment.Verification, "verification.format")
	if !ok ||
		format.Command != "bun run fmt" ||
		!format.RepositoryExecutable ||
		format.DeclarationPath != "package.json" ||
		format.Classification != VerificationRepositoryCommand {
		t.Fatalf("repository formatter projection = %+v, found=%v", format, ok)
	}
	if divergence, exists := findProfileDivergence(alignment.Divergences, "verification.format"); exists {
		t.Fatalf("declared repository formatter remained divergent: %+v", divergence)
	}
}

func TestProfileAlignmentEquivalentNormalizedDecisions(t *testing.T) {
	t.Parallel()

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

func TestExecutableCandidateResolution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		arrange      func(*testing.T, string, string) string
		wantResolved func(string, string) string
		wantReason   string
		wantHops     int
	}{
		{
			name: "direct regular executable",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				writeProfileAlignmentExecutable(t, bin, name, "executable")
				return filepath.Join(bin, name)
			},
			wantResolved: func(bin, name string) string { return filepath.Join(bin, name) },
		},
		{
			name: "one-hop symlink",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				target := filepath.Join(bin, "target")
				writeProfileAlignmentExecutable(t, bin, "target", "executable")
				if err := os.Symlink(filepath.Base(target), filepath.Join(bin, name)); err != nil {
					t.Fatalf("create one-hop executable symlink: %v", err)
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(bin, _ string) string { return filepath.Join(bin, "target") },
			wantHops:     1,
		},
		{
			name: "multi-hop symlink",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				writeProfileAlignmentExecutable(t, bin, "target", "executable")
				if err := os.Symlink("target", filepath.Join(bin, "middle")); err != nil {
					t.Fatalf("create middle executable symlink: %v", err)
				}
				if err := os.Symlink("middle", filepath.Join(bin, name)); err != nil {
					t.Fatalf("create first executable symlink: %v", err)
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(bin, _ string) string { return filepath.Join(bin, "target") },
			wantHops:     2,
		},
		{
			name: "link cycle",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				if err := os.Symlink("second", filepath.Join(bin, name)); err != nil {
					t.Fatalf("create first cycle link: %v", err)
				}
				if err := os.Symlink(name, filepath.Join(bin, "second")); err != nil {
					t.Fatalf("create second cycle link: %v", err)
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(_, _ string) string { return "" },
			wantReason:   executableProbeReasonLinkCycle,
			wantHops:     2,
		},
		{
			name: "link depth exceeded",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				writeProfileAlignmentExecutable(t, bin, "target", "executable")
				for index := range maxExecutableLinkHops + 1 {
					link := name
					if index != 0 {
						link = fmt.Sprintf("link-%d", index)
					}
					target := "target"
					if index != maxExecutableLinkHops {
						target = fmt.Sprintf("link-%d", index+1)
					}
					if err := os.Symlink(target, filepath.Join(bin, link)); err != nil {
						t.Fatalf("create executable symlink %d: %v", index, err)
					}
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(_, _ string) string { return "" },
			wantReason:   "link-depth-exceeded",
			wantHops:     maxExecutableLinkHops,
		},
		{
			name: "broken link",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				if err := os.Symlink("missing-target", filepath.Join(bin, name)); err != nil {
					t.Fatalf("create broken executable symlink: %v", err)
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(_, _ string) string { return "" },
			wantReason:   executableProbeReasonBrokenLink,
			wantHops:     1,
		},
		{
			name: "non-executable target",
			arrange: func(t *testing.T, bin, name string) string {
				t.Helper()
				if err := os.WriteFile(filepath.Join(bin, "target"), []byte("not executable"), 0o644); err != nil {
					t.Fatalf("write non-executable target: %v", err)
				}
				if err := os.Symlink("target", filepath.Join(bin, name)); err != nil {
					t.Fatalf("create non-executable target symlink: %v", err)
				}
				return filepath.Join(bin, name)
			},
			wantResolved: func(_, _ string) string { return "" },
			wantReason:   executableProbeReasonNotExecutable,
			wantHops:     1,
		},
		{
			name: "absent candidate",
			arrange: func(_ *testing.T, _, _ string) string {
				return ""
			},
			wantResolved: func(_, _ string) string { return "" },
			wantReason:   executableProbeReasonNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bin := t.TempDir()
			const name = "probe"
			wantCandidate := tt.arrange(t, bin, name)

			result := resolveExecutableCandidate(name, []string{bin})

			if result.Candidate != wantCandidate ||
				result.Resolved != tt.wantResolved(bin, name) ||
				result.Reason != tt.wantReason ||
				result.HopCount != tt.wantHops {
				t.Fatalf("resolveExecutableCandidate(%q) = %+v", name, result)
			}
		})
	}
}

func TestExecutableCandidateNeverExecutes(t *testing.T) {
	t.Parallel()

	bin := t.TempDir()
	marker := filepath.Join(t.TempDir(), "executed")
	target := filepath.Join(bin, "target")
	content := "#!/bin/sh\nprintf executed > \"" + marker + "\"\n"
	writeProfileAlignmentExecutable(t, bin, filepath.Base(target), content)
	if err := os.Symlink(filepath.Base(target), filepath.Join(bin, "probe")); err != nil {
		t.Fatalf("create executable probe symlink: %v", err)
	}
	evidence := collectExecutableEvidence(RepositoryCapability{
		EvidenceKind: CapabilityEvidenceExecutable,
		Probe:        map[string]any{"executable": "probe"},
	}, []string{bin})

	if evidence.Status != CapabilityEvidencePresent || evidence.SourcePath != filepath.ToSlash(filepath.Join(bin, "probe")) {
		t.Fatalf("collectExecutableEvidence(probe) = %+v", evidence)
	}
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("executable candidate was invoked: %v", err)
	}
}

func TestExecutableEvidenceDistinguishesFailureFromAbsence(t *testing.T) {
	t.Parallel()

	const name = "probe"
	capability := RepositoryCapability{
		EvidenceKind: CapabilityEvidenceExecutable,
		Probe:        map[string]any{"executable": name},
	}

	t.Run("failed candidate", func(t *testing.T) {
		bin := t.TempDir()
		candidate := filepath.Join(bin, name)
		if err := os.Symlink("missing-target", candidate); err != nil {
			t.Fatalf("create broken executable symlink: %v", err)
		}

		evidence := collectExecutableEvidence(capability, []string{bin})

		if evidence.Status != CapabilityEvidenceInvalid ||
			evidence.SourcePath != filepath.ToSlash(candidate) ||
			evidence.Detail != executableProbeReasonBrokenLink {
			t.Fatalf("failed executable evidence = %+v", evidence)
		}
	})

	t.Run("absent candidate", func(t *testing.T) {
		evidence := collectExecutableEvidence(capability, []string{t.TempDir()})

		if evidence.Status != CapabilityEvidenceAbsent || evidence.SourcePath != "" || evidence.Detail != executableProbeReasonNotFound {
			t.Fatalf("absent executable evidence = %+v", evidence)
		}
	})
}

func TestCapabilityAuditNoExecution(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	bin := t.TempDir()
	marker := filepath.Join(t.TempDir(), "executed")
	for _, name := range []string{"docker", "rg", "rtk"} {
		writeProfileAlignmentExecutable(t, bin, name, "#!/bin/sh\nprintf executed > \""+marker+"\"\n")
	}
	before := snapshotInspectionTree(t, repository)
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{bin},
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

func TestProfileDivergenceResolution(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(
		t,
		repository,
		"package.json",
		`{"name":"root","packageManager":"bun@1.3.0","scripts":{"verify":"true"},"dependencies":{"hono":"latest","typescript":"latest","zod":"latest"}}`,
	)
	if err := os.RemoveAll(filepath.Join(repository, "packages", "frontend")); err != nil {
		t.Fatalf("remove frontend fixture: %v", err)
	}

	source, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatalf("resolve source Profile: %v", err)
	}
	decisions := standardTypeScriptDecisions("make verify")
	before, err := ResolveProfileAlignment(
		context.Background(),
		repository,
		ProfileAlignmentRequest{
			ProfileID:            source.ID,
			Decisions:            decisions,
			RemediationProfileID: source.ID,
		},
		catalog,
	)
	if err != nil {
		t.Fatalf("resolve source divergence: %v", err)
	}
	if before.Ready {
		t.Fatal("backend-only fixture unexpectedly aligned with the full TypeScript Profile")
	}

	profileCapabilities := stringSet(source.Capabilities)
	removedCapabilities := make([]string, 0)
	for _, divergence := range before.Divergences {
		if !divergence.Blocking {
			continue
		}
		if _, profileSpecific := profileCapabilities[divergence.ID]; profileSpecific {
			removedCapabilities = append(removedCapabilities, divergence.ID)
		}
	}
	input, err := NewProfileAdaptationDraft(
		source.ID,
		"guided-backend",
		[]string{"frontend", "autonomous-work"},
		removedCapabilities,
		catalog,
	)
	if err != nil {
		t.Fatalf("construct reviewed Profile adaptation: %v", err)
	}
	adapted, canonical, err := ResolveProfileDraft(repository, input, catalog)
	if err != nil {
		t.Fatalf("resolve reviewed Profile adaptation: %v", err)
	}
	input.Document = canonical
	after, err := ResolveProfileAlignment(
		context.Background(),
		repository,
		ProfileAlignmentRequest{
			ProfileID:            adapted.ID,
			Decisions:            decisionsForResolvedProfile(decisions, adapted),
			Profile:              &adapted,
			RemediationProfileID: input.SourceProfileID,
		},
		catalog,
	)
	if err != nil {
		t.Fatalf("re-audit reviewed Profile adaptation: %v", err)
	}
	if !after.Ready || after.Profile.ID != "guided-backend" {
		t.Fatalf("reviewed Profile adaptation remains unresolved: %+v", after.Divergences)
	}
	for _, requiredID := range []string{"capability.context7", "capability.exa"} {
		outcome, ok := findCapabilityOutcome(after.Capabilities, requiredID)
		if !ok || outcome.Requirement != CapabilityRequired || outcome.Status != CapabilitySatisfied {
			t.Fatalf("universal requirement %s became a waiver: %+v found=%v", requiredID, outcome, ok)
		}
	}
	if _, err := os.Stat(filepath.Join(repository, filepath.FromSlash(adapted.Path))); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("alignment wrote Profile draft target: %v", err)
	}
}

func TestUniversalCapabilityRemediation(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	if err := os.Remove(filepath.Join(repository, ".agents", "skills", "context7", "SKILL.md")); err != nil {
		t.Fatalf("remove Context7 fixture: %v", err)
	}

	alignment, err := ResolveProfileAlignment(
		context.Background(),
		repository,
		ProfileAlignmentRequest{
			ProfileID:            "standard-typescript-monorepo",
			Decisions:            standardTypeScriptDecisions("make verify"),
			RemediationProfileID: "standard-typescript-monorepo",
		},
		catalog,
	)
	if err != nil {
		t.Fatalf("resolve universal capability remediation: %v", err)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.context7")
	if !ok || !divergence.Blocking || divergence.Requirement != CapabilityRequired {
		t.Fatalf("Context7 divergence = %+v found=%v", divergence, ok)
	}
	for _, want := range []string{
		"roundfix baseline skills restore",
		"--profile standard-typescript-monorepo",
		"--skill context7",
		"--confirm-plan <digest>",
	} {
		if !strings.Contains(divergence.NextAction, want) {
			t.Fatalf("universal remediation missing %q: %s", want, divergence.NextAction)
		}
	}
}

func TestDivergenceRendersProbe(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(false, true))
	if err := os.Remove(filepath.Join(repository, "packages", "frontend", "package.json")); err != nil {
		t.Fatalf("remove declared-file fixture: %v", err)
	}
	emptyExecutableDirectory := t.TempDir()

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{emptyExecutableDirectory},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve divergence rendering fixture: %v", err)
	}
	rendered := RenderProfileDivergences(alignment.Divergences)
	for _, want := range []string{
		`package.json: present (expected content not found); expected content: "better-auth"`,
		`packages/frontend/package.json: absent (file not found); expected content: "better-auth"`,
		`packages/backend/package.json: present (expected content not found); expected content: "better-auth"`,
		"Inspected candidate: none existed",
		"Selected technology: Better Auth",
		"Repository remediation:",
		"Profile adaptation:",
		"Decision cascade: auth.provider",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered divergences missing %q:\n%s", want, rendered)
		}
	}

	betterAuth, ok := findProfileDivergence(alignment.Divergences, "capability.stack.better-auth")
	if !ok || betterAuth.CapabilityResolution == nil {
		t.Fatalf("Better Auth divergence resolution = %+v, found=%v", betterAuth, ok)
	}
	if betterAuth.CapabilityResolution.SelectedTechnology != "Better Auth" ||
		betterAuth.CapabilityResolution.RepositoryRemediation == "" ||
		betterAuth.CapabilityResolution.ProfileAdaptation == "" ||
		!slices.Equal(betterAuth.CapabilityResolution.RemovedDecisions, []string{"auth.provider"}) {
		t.Fatalf("Better Auth resolution = %+v", betterAuth.CapabilityResolution)
	}

	withoutShadcn := strings.Replace(typeScriptPackageJSON(false, true), `"shadcn":"latest",`, "", 1)
	writeProfileAlignmentFile(t, repository, "package.json", withoutShadcn)
	alignment, err = ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{emptyExecutableDirectory},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve decision-free cascade fixture: %v", err)
	}
	shadcn, ok := findProfileDivergence(alignment.Divergences, "capability.stack.shadcn")
	if !ok || shadcn.CapabilityResolution == nil {
		t.Fatalf("shadcn divergence resolution = %+v, found=%v", shadcn, ok)
	}
	if len(shadcn.CapabilityResolution.RemovedDecisions) != 0 ||
		!strings.Contains(RenderProfileDivergences([]ProfileDivergence{shadcn}), "Decision cascade: none") {
		t.Fatalf("shadcn decision cascade = %+v", shadcn.CapabilityResolution)
	}

	bin := t.TempDir()
	writeProfileAlignmentFile(t, bin, "docker", "not executable")
	alignment, err = ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{bin},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve rejected executable fixture: %v", err)
	}
	docker, ok := findProfileDivergence(alignment.Divergences, "capability.optional.docker")
	if !ok {
		t.Fatal("Docker divergence is missing")
	}
	rendered = RenderProfileDivergences([]ProfileDivergence{docker})
	wantCandidate := "Inspected candidate: " + filepath.ToSlash(filepath.Join(bin, "docker")) + " (not-executable)"
	if !strings.Contains(rendered, wantCandidate) {
		t.Fatalf("executable divergence missing inspected candidate %q:\n%s", wantCandidate, rendered)
	}
}

func TestCapabilityTextRendersProbe(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(false, true))
	if err := os.Remove(filepath.Join(repository, "packages", "frontend", "package.json")); err != nil {
		t.Fatalf("remove declared-file fixture: %v", err)
	}

	bin := t.TempDir()
	candidate := filepath.Join(bin, "rtk")
	if err := os.Symlink("missing-rtk-target", candidate); err != nil {
		t.Fatalf("create broken rtk candidate: %v", err)
	}
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{bin},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve rejected executable fixture: %v", err)
	}
	rendered := RenderProfileDivergences(alignment.Divergences)
	for _, want := range []string{
		"Inspected candidate: " + filepath.ToSlash(candidate) + " (broken-link)",
		"Repair the inspected candidate " + filepath.ToSlash(candidate) + " (broken-link)",
		`packages/frontend/package.json: absent (file not found); expected content: "better-auth"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("capability text missing %q:\n%s", want, rendered)
		}
	}
	if strings.Contains(rendered, "Install rtk") {
		t.Errorf("rejected rtk candidate incorrectly recommends installation:\n%s", rendered)
	}

	emptyExecutableDirectory := t.TempDir()
	alignment, err = ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{emptyExecutableDirectory},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve absent executable fixture: %v", err)
	}
	rtkDivergence, ok := findProfileDivergence(alignment.Divergences, "capability.rtk")
	if !ok {
		t.Fatal("absent rtk divergence is missing")
	}
	rendered = RenderProfileDivergences([]ProfileDivergence{rtkDivergence})
	if !strings.Contains(rendered, "Inspected candidate: none existed") ||
		!strings.Contains(rendered, "Install rtk") {
		t.Fatalf("absent candidate text does not distinguish absence:\n%s", rendered)
	}
}

func TestCapabilityTextAndJSONAgree(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	bin := t.TempDir()
	candidate := filepath.Join(bin, "rtk")
	if err := os.Symlink("missing-rtk-target", candidate); err != nil {
		t.Fatalf("create broken rtk candidate: %v", err)
	}
	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{bin},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve text and JSON fixture: %v", err)
	}
	divergence, ok := findProfileDivergence(alignment.Divergences, "capability.rtk")
	if !ok || len(divergence.Evidence) != 1 {
		t.Fatalf("rtk divergence = %+v, found=%v", divergence, ok)
	}

	machineBefore, err := json.Marshal(divergence)
	if err != nil {
		t.Fatalf("marshal machine divergence: %v", err)
	}
	rendered := RenderProfileDivergences([]ProfileDivergence{divergence})
	evidence := divergence.Evidence[0]
	if !strings.Contains(rendered, evidence.SourcePath) || !strings.Contains(rendered, evidence.Detail) {
		t.Fatalf("text and machine evidence disagree:\njson=%s\ntext=%s", machineBefore, rendered)
	}
	machineAfter, err := json.Marshal(divergence)
	if err != nil {
		t.Fatalf("marshal machine divergence after text rendering: %v", err)
	}
	if string(machineAfter) != string(machineBefore) {
		t.Fatalf("text rendering changed machine bytes:\nbefore=%s\nafter=%s", machineBefore, machineAfter)
	}
}

func TestDivergenceGroupsByRequirement(t *testing.T) {
	t.Parallel()

	catalog := loadProfileAlignmentCatalog(t)
	repository := newAlignedTypeScriptRepository(t)
	writeProfileAlignmentFile(t, repository, "package.json", typeScriptPackageJSON(false, true))
	emptyExecutableDirectory := t.TempDir()

	alignment, err := ResolveProfileAlignment(context.Background(), repository, ProfileAlignmentRequest{
		ProfileID:             "standard-typescript-monorepo",
		Decisions:             standardTypeScriptDecisions("make verify"),
		ExecutableDirectories: []string{emptyExecutableDirectory},
	}, catalog)
	if err != nil {
		t.Fatalf("resolve grouped divergences: %v", err)
	}
	if alignment.Ready || alignment.State != ProfileAlignmentActionRequired {
		t.Fatalf("grouping changed blocking readiness: %+v", alignment)
	}

	groupRank := map[ProfileDivergenceGroup]int{
		ProfileDivergenceBlocking:      0,
		ProfileDivergenceAdvisory:      1,
		ProfileDivergenceInformational: 2,
	}
	lastRank := -1
	seen := make(map[ProfileDivergenceGroup]bool)
	for _, divergence := range alignment.Divergences {
		wantGroup := profileDivergenceGroup(divergence.Requirement)
		if divergence.Group != wantGroup {
			t.Errorf("divergence %q group = %q, want %q", divergence.ID, divergence.Group, wantGroup)
		}
		rank := groupRank[divergence.Group]
		if rank < lastRank {
			t.Fatalf("divergence groups are not ordered: %+v", alignment.Divergences)
		}
		lastRank = rank
		seen[divergence.Group] = true
		if divergence.Group == ProfileDivergenceAdvisory {
			if divergence.NonBlockingStatement != advisoryNonBlockingStatement {
				t.Errorf("advisory %q statement = %q", divergence.ID, divergence.NonBlockingStatement)
			}
			encoded, marshalErr := json.Marshal(divergence)
			if marshalErr != nil {
				t.Fatalf("marshal advisory %q: %v", divergence.ID, marshalErr)
			}
			statementAt := strings.Index(string(encoded), `"nonBlockingStatement"`)
			nextActionAt := strings.Index(string(encoded), `"nextAction"`)
			if statementAt < 0 || nextActionAt >= 0 && statementAt > nextActionAt {
				t.Errorf("advisory machine statement does not precede next action: %s", encoded)
			}
		}
	}
	for _, group := range []ProfileDivergenceGroup{
		ProfileDivergenceBlocking,
		ProfileDivergenceAdvisory,
		ProfileDivergenceInformational,
	} {
		if !seen[group] {
			t.Errorf("group %q is absent from machine divergences", group)
		}
	}

	rendered := RenderProfileDivergences(alignment.Divergences)
	blockingAt := strings.Index(rendered, "Blocking divergences:")
	advisoryAt := strings.Index(rendered, "Advisory divergences:")
	informationalAt := strings.Index(rendered, "Informational divergences:")
	if blockingAt < 0 || advisoryAt <= blockingAt || informationalAt <= advisoryAt {
		t.Fatalf("rendered divergence groups are not ordered:\n%s", rendered)
	}
	advisorySection := rendered[advisoryAt:informationalAt]
	statementAt := strings.Index(advisorySection, advisoryNonBlockingStatement)
	nextActionAt := strings.Index(advisorySection, "Next action:")
	if statementAt < 0 || nextActionAt >= 0 && statementAt > nextActionAt {
		t.Fatalf("rendered advisory statement does not precede next action:\n%s", advisorySection)
	}
}

func decisionsForResolvedProfile(input []DecisionValue, profile ResolvedProfile) []DecisionValue {
	selected := stringSet(profile.Decisions)
	result := make([]DecisionValue, 0, len(input))
	for _, decision := range input {
		if _, ok := selected[decision.ID]; ok {
			result = append(result, decision)
		}
	}
	return result
}

func findCapabilityOutcome(outcomes []CapabilityOutcome, id string) (CapabilityOutcome, bool) {
	for _, outcome := range outcomes {
		if outcome.ID == id {
			return outcome, true
		}
	}
	return CapabilityOutcome{}, false
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
		{ID: "identifier.strategy", Value: map[string]any{"kind": "uuid-v7"}},
		{ID: "http.contract", Value: map[string]any{"mode": "REST"}},
		{ID: "auth.provider", Value: completeAuthProviderDecision()},
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
