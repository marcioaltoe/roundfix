package baseline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	PlanSchemaVersion   = "roundfix/baseline-plan/v1"
	ResultSchemaVersion = "roundfix/baseline-result/v1"
	ManifestSchema      = "setup-context-driven/manifest/0.0.1"
	ManifestVersion     = "0.0.1"

	planDigestDomain = PlanSchemaVersion + "\x00"
	manifestPath     = "docs/agents/setup-context.json"

	specificRepositoryPath      = "docs/agents/specific-repository.md"
	legacyRepositoryPath        = "docs/agents/repository.md"
	legacyRepositoryRulesPath   = "docs/agents/repository-rules.md"
	legacyRepositoryScaffold    = "# Repository instructions\n\nAdd project-specific hard rules here. Setup preserves this file byte-for-byte.\n"
	repositoryExtensionModuleID = "repository-extension"
	repositoryExtensionRootID   = "root.repository-extension"
)

// PlanRequest is the complete normalized non-interactive planning input.
type PlanRequest struct {
	Repository            string
	ProfileID             string
	ProfileDraft          *ProfileDraftInput
	Decisions             []DecisionValue
	Preservation          RootPreservationRequest
	ExecutableDirectories []string
}

// ProfileDraftInput binds one strict repository-owned Profile draft to the
// built-in Profile it is allowed to narrow.
type ProfileDraftInput struct {
	SourceProfileID string
	Document        []byte
}

// PlanOutcome contains either one complete portable plan or an actionable
// result. A nil Plan always means that no partial plan is available.
type PlanOutcome struct {
	Plan   *PlanDocument
	Result Result
}

type specificRepositoryPlan struct {
	IncludeRoot      bool
	CanonicalContent []byte
	DeletePaths      []string
}

type repositoryRuleInventory struct {
	ByPath    map[string][]RepositoryRuleBlock
	Retention []RetentionEvidence
}

type plannedProfileDraft struct {
	ID      string
	Path    string
	Content []byte
}

// ResolveDecisionInput normalizes human or automation answers through the
// same catalog decision graph used by BuildPlan. Missing IDs are returned in
// stable lexical order without producing a partial plan.
func ResolveDecisionInput(
	profile ResolvedProfile,
	input []DecisionValue,
	catalog *Catalog,
) ([]DecisionValue, []string, error) {
	if catalog == nil {
		return nil, nil, errors.New("resolve Baseline decisions: catalog is required")
	}
	return normalizePlanDecisions(profile, input, catalog)
}

// CatalogIdentity binds a plan to the embedded catalog authority.
type CatalogIdentity struct {
	SchemaVersion string `json:"schemaVersion"`
	Digest        string `json:"digest"`
}

// RetentionEvidence records the complete transition accounting considered by
// planning. Initial adoption has an empty ledger.
type RetentionEvidence struct {
	FromClause  string   `json:"fromClause"`
	Enforcement string   `json:"enforcement"`
	Disposition string   `json:"disposition"`
	Targets     []string `json:"targets"`
	Reason      string   `json:"reason"`
}

// Postimage carries exact repository-relative bytes approved by a plan.
type Postimage struct {
	Path            string       `json:"path"`
	Kind            PreimageKind `json:"kind"`
	Mode            uint32       `json:"mode"`
	Content         []byte       `json:"content,omitempty"`
	ContentIdentity string       `json:"contentIdentity,omitempty"`
}

// ManagedEntry is the canonical ordered change ledger. FileChanges is derived
// exclusively from these entries and the preimage/postimage sets.
type ManagedEntry struct {
	Ordinal         int    `json:"ordinal"`
	ID              string `json:"id"`
	Path            string `json:"path"`
	Action          string `json:"action"`
	Kind            string `json:"kind"`
	Module          string `json:"module,omitempty"`
	Template        string `json:"template,omitempty"`
	Version         string `json:"version,omitempty"`
	BeforeIdentity  string `json:"beforeIdentity,omitempty"`
	AfterIdentity   string `json:"afterIdentity,omitempty"`
	ContentIdentity string `json:"contentIdentity,omitempty"`
}

// FileChange is the concise, reproducible projection of one affected path.
type FileChange struct {
	Path           string   `json:"path"`
	Action         string   `json:"action"`
	BeforeIdentity string   `json:"beforeIdentity,omitempty"`
	AfterIdentity  string   `json:"afterIdentity,omitempty"`
	ManagedEntries []string `json:"managedEntries"`
}

// ManifestDecision records one normalized owner decision without a volatile
// confirmation timestamp.
type ManifestDecision struct {
	Value any `json:"value"`
}

// ManifestGenerator identifies the Go-owned Baseline authority.
type ManifestGenerator struct {
	Skill    string `json:"skill"`
	Version  string `json:"version"`
	Baseline string `json:"baseline"`
}

// ManifestArtifact records one setup-owned managed artifact.
type ManifestArtifact struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	Module   string `json:"module"`
	Template string `json:"template"`
	Version  string `json:"version"`
	Digest   string `json:"digest"`
}

// SetupManifest is the current Go-owned setup identity included in every
// portable plan.
type SetupManifest struct {
	SchemaVersion    string                      `json:"schemaVersion"`
	Version          string                      `json:"version"`
	Generator        ManifestGenerator           `json:"generator"`
	Profile          string                      `json:"profile"`
	ProfileDigest    string                      `json:"profileDigest"`
	CatalogDigest    string                      `json:"catalogDigest"`
	Modules          []string                    `json:"modules"`
	Decisions        map[string]ManifestDecision `json:"decisions"`
	ManagedArtifacts []ManifestArtifact          `json:"managedArtifacts"`
	LocalSkills      []string                    `json:"localSkills"`
	Verification     []VerificationProjection    `json:"verification"`
}

// PlanDocument is the complete portable approval artifact.
type PlanDocument struct {
	SchemaVersion  string              `json:"schemaVersion"`
	Repository     RepositoryIdentity  `json:"repository"`
	Catalog        CatalogIdentity     `json:"catalog"`
	Profile        ResolvedProfile     `json:"profile"`
	Decisions      []DecisionValue     `json:"decisions"`
	Retention      []RetentionEvidence `json:"retention"`
	ClauseDelta    *ClauseDelta        `json:"clauseDelta,omitempty"`
	Preimages      []Preimage          `json:"preimages"`
	Postimages     []Postimage         `json:"postimages"`
	Warnings       []Finding           `json:"warnings"`
	SetupManifest  SetupManifest       `json:"setupManifest"`
	ManagedEntries []ManagedEntry      `json:"managedEntries"`
	FileChanges    []FileChange        `json:"fileChanges"`
	PlanDigest     string              `json:"planDigest"`
}

// EvidenceStatus reports whether one result axis produced affirmative evidence.
type EvidenceStatus string

const (
	EvidenceStatusVerified EvidenceStatus = "verified"
	EvidenceStatusNotRun   EvidenceStatus = "not run"
)

// ResultStatusMatrix keeps independent Baseline evidence from collapsing into
// one success state.
type ResultStatusMatrix struct {
	ApprovedPostimages     EvidenceStatus `json:"approvedPostimages"`
	SemanticRetention      EvidenceStatus `json:"semanticRetention"`
	ProfileAlignment       EvidenceStatus `json:"profileAlignment"`
	RepositoryVerification EvidenceStatus `json:"repositoryVerification"`
	Idempotence            EvidenceStatus `json:"idempotence"`
}

// Result is the strict automation result used when no complete plan can be
// emitted and by later Baseline operations.
type Result struct {
	SchemaVersion      string              `json:"schemaVersion"`
	Operation          string              `json:"operation"`
	State              string              `json:"state"`
	Category           string              `json:"category,omitempty"`
	Message            string              `json:"message,omitempty"`
	NextAction         string              `json:"nextAction,omitempty"`
	PlanDigest         string              `json:"planDigest,omitempty"`
	VerifiedPostimages []Postimage         `json:"verifiedPostimages"`
	Warnings           []Finding           `json:"warnings"`
	Recommendations    []string            `json:"recommendations"`
	ClauseDelta        *ClauseDelta        `json:"clauseDelta,omitempty"`
	StatusMatrix       *ResultStatusMatrix `json:"statusMatrix,omitempty"`
}

type plannedArtifact struct {
	ID       string
	Path     string
	Kind     string
	Module   string
	Template string
	Version  string
	Body     string
	Digest   string
}

// BuildPlan produces a portable plan without mutating repository state.
func BuildPlan(ctx context.Context, request PlanRequest) (PlanOutcome, error) {
	if err := ctx.Err(); err != nil {
		return PlanOutcome{}, err
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return PlanOutcome{}, err
	}
	return buildPlanWithCatalog(ctx, request, catalog)
}

func buildPlanWithCatalog(
	ctx context.Context,
	request PlanRequest,
	catalog *Catalog,
) (PlanOutcome, error) {
	initial, err := InspectRepository(ctx, request.Repository, nil)
	if err != nil {
		return PlanOutcome{}, err
	}
	if len(initial.Snapshot.Blocking) != 0 {
		return actionOutcome("preflight", "repository preflight found unsafe bounded carriers",
			"repair every blocking carrier and rerun roundfix baseline plan",
			append(initial.Snapshot.Warnings, initial.Snapshot.Blocking...)), nil
	}
	profileID := strings.TrimSpace(request.ProfileID)
	if profileID != "" && request.ProfileDraft != nil {
		return PlanOutcome{}, errors.New("select exactly one Baseline Profile ID or Profile draft")
	}
	if profileID == "" && request.ProfileDraft == nil {
		return actionOutcome("decision", "Baseline Profile selection is required",
			"rerun with --profile <id>", initial.Snapshot.Warnings), nil
	}
	var profile ResolvedProfile
	var plannedDraft *plannedProfileDraft
	if request.ProfileDraft != nil {
		var profileBytes []byte
		profile, profileBytes, err = resolveProfileDraft(initial.Root, *request.ProfileDraft, catalog)
		if err != nil {
			return PlanOutcome{}, fmt.Errorf("resolve Baseline Profile draft: %w", err)
		}
		plannedDraft = &plannedProfileDraft{
			ID:      profile.ID,
			Path:    profile.Path,
			Content: profileBytes,
		}
	} else {
		profile, err = ResolveProfile(initial.Root, profileID, catalog)
		if err != nil {
			return PlanOutcome{}, fmt.Errorf("resolve Baseline Profile: %w", err)
		}
		profile.Path = portableProfilePath(initial.Root, profile.Path)
	}

	decisions, missing, err := normalizePlanDecisions(profile, request.Decisions, catalog)
	if err != nil {
		return PlanOutcome{}, err
	}
	if len(missing) != 0 {
		return actionOutcome("decision",
			"required Baseline decisions are missing: "+strings.Join(missing, ", "),
			"answer every named decision with --decision or --decision-file and rerun",
			initial.Snapshot.Warnings), nil
	}
	if request.Preservation.Mode == "" {
		return actionOutcome("decision", "instruction-preservation mode is required",
			"rerun with --decision preservation.mode=greenfield or preservation",
			initial.Snapshot.Warnings), nil
	}
	_, carrierArtifacts, err := resolveManagedArtifacts(
		catalog,
		profile,
		decisions,
		false,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	initialClassifications, err := classifyCarriers(
		initial.Root,
		initial.Snapshot.Carriers,
		catalog,
		carrierArtifacts,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	initial.Snapshot.Warnings = warningsForCarrierClassifications(
		initial.Snapshot.Warnings,
		initialClassifications,
	)
	remediationProfileID := profile.ID
	if request.ProfileDraft != nil {
		remediationProfileID = request.ProfileDraft.SourceProfileID
	}
	alignment, err := ResolveProfileAlignment(ctx, initial.Root, ProfileAlignmentRequest{
		ProfileID:             profile.ID,
		Decisions:             profileAlignmentDecisions(profile, decisions),
		Profile:               &profile,
		RemediationProfileID:  remediationProfileID,
		ExecutableDirectories: request.ExecutableDirectories,
	}, catalog)
	if err != nil {
		return PlanOutcome{}, err
	}
	if !alignment.Ready {
		blocking := make([]string, 0)
		nextActions := make([]string, 0)
		for _, divergence := range alignment.Divergences {
			if !divergence.Blocking {
				continue
			}
			blocking = append(blocking, divergence.ID)
			if divergence.NextAction != "" && !slices.Contains(nextActions, divergence.NextAction) {
				nextActions = append(nextActions, divergence.NextAction)
			}
		}
		nextAction := "provide the missing repository evidence or select a different profile"
		if len(nextActions) != 0 {
			nextAction = strings.Join(nextActions, " ")
		}
		return actionOutcome("decision",
			"required profile alignment is unresolved: "+strings.Join(blocking, ", "),
			nextAction,
			initial.Snapshot.Warnings), nil
	}

	semanticOwners, err := ResolveSemanticOwnerRegistry(catalog, profile, decisions)
	if err != nil {
		return PlanOutcome{}, err
	}
	preservationRequest := request.Preservation
	preservationRequest.semanticOwners = semanticOwners
	preservationRequest.managedArtifacts = carrierArtifacts
	preservationRequest.classifyCarriers = true
	preservationRequest.classifications = initialClassifications

	preservation, err := planRootPreservationWithCatalog(initial, preservationRequest, catalog)
	if err != nil {
		return PlanOutcome{}, err
	}
	if preservation.State != PreservationStateReady {
		return actionOutcome("decision", preservation.NextAction, preservation.NextAction,
			append(initial.Snapshot.Warnings, preservation.Findings...)), nil
	}

	repositoryPlan, repositoryFindings, err := planSpecificRepository(
		initial.Root,
		preservation.RepositoryRulesBytes,
		decisionBool(decisions, "repository.extension.enabled"),
		preservationRedistributesRecognizedRepositoryRules(preservation),
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	if len(repositoryFindings) != 0 {
		return actionOutcome(
			"classification",
			repositoryFindings[0].Message,
			"resolve the reported repository-specific rule carrier conflict and rerun Baseline planning",
			append(initial.Snapshot.Warnings, repositoryFindings...),
		), nil
	}

	activeModules, artifacts, err := resolveManagedArtifacts(
		catalog,
		profile,
		decisions,
		repositoryPlan.IncludeRoot,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	manifest := buildSetupManifest(catalog, profile, decisions, activeModules, artifacts, alignment.Verification)
	manifestBytes, err := marshalSetupManifestBytes(manifest)
	if err != nil {
		return PlanOutcome{}, fmt.Errorf("serialize Setup Manifest: %w", err)
	}
	manifestBytes = append(manifestBytes, '\n')

	mutablePaths := []string{manifestPath}
	for _, artifact := range artifacts {
		mutablePaths = append(mutablePaths, artifact.Path)
	}
	for _, backup := range preservation.Backups {
		mutablePaths = append(mutablePaths, backup.Path)
	}
	mutablePaths = append(
		mutablePaths,
		specificRepositoryPath,
		legacyRepositoryPath,
		legacyRepositoryRulesPath,
	)
	mutablePaths = append(mutablePaths, alignmentEvidencePaths(alignment)...)
	if plannedDraft != nil {
		mutablePaths = append(mutablePaths, plannedDraft.Path)
	}
	snapshot, err := inspectRepositorySnapshot(initial.Root, InventoryRequest{MutablePaths: mutablePaths})
	if err != nil {
		return PlanOutcome{}, err
	}
	currentClassifications, err := classifyCarriers(
		initial.Root,
		snapshot.Carriers,
		catalog,
		carrierArtifacts,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	snapshot.Warnings = warningsForCarrierClassifications(
		snapshot.Warnings,
		currentClassifications,
	)
	preservationRequest.classifications = currentClassifications
	if plannedDraft != nil {
		if err := validateProfileDraftTarget(snapshot, *plannedDraft); err != nil {
			return PlanOutcome{}, err
		}
	}
	inspection := RepositoryInspection{Root: initial.Root, Identity: initial.Identity, Snapshot: snapshot}
	preservation, err = planRootPreservationWithCatalog(inspection, preservationRequest, catalog)
	if err != nil {
		return PlanOutcome{}, err
	}
	if preservation.State != PreservationStateReady {
		return actionOutcome("decision", preservation.NextAction, preservation.NextAction,
			append(snapshot.Warnings, preservation.Findings...)), nil
	}
	currentRepositoryPlan, repositoryFindings, err := planSpecificRepository(
		initial.Root,
		preservation.RepositoryRulesBytes,
		decisionBool(decisions, "repository.extension.enabled"),
		preservationRedistributesRecognizedRepositoryRules(preservation),
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	if len(repositoryFindings) != 0 {
		return actionOutcome(
			"classification",
			repositoryFindings[0].Message,
			"resolve the reported repository-specific rule carrier conflict and rerun Baseline planning",
			append(snapshot.Warnings, repositoryFindings...),
		), nil
	}
	if !reflectJSONEqual(repositoryPlan, currentRepositoryPlan) {
		return PlanOutcome{}, errors.New("repository-specific rule carriers changed during planning")
	}
	repositoryRules, err := inventoryRepositoryRuleBlocks(
		initial.Root,
		artifacts,
		preservation.RepositoryRuleBlocks,
	)
	if err != nil {
		return actionOutcome(
			"classification",
			err.Error(),
			"repair the repository-owned rule markers and rerun Baseline planning",
			snapshot.Warnings,
		), nil
	}
	retention, clauseDelta, retentionAction, err := resolvePlanRetention(
		initial.Root,
		snapshot,
		catalog,
		compatibleRetentionProfileIDs(profile.ID, request.ProfileDraft),
		manifest,
		activeModules,
		preservation,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	if retentionAction != "" {
		outcome := actionOutcome(
			"classification",
			retentionAction,
			"restore an exact supported Setup Manifest or add a reviewed Upgrade Retention Contract",
			snapshot.Warnings,
		)
		outcome.Result.ClauseDelta = clauseDelta
		return outcome, nil
	}
	retention = append(retention, repositoryRules.Retention...)

	postimages, ledger, err := assemblePostimages(
		initial.Root,
		snapshot,
		artifacts,
		manifestBytes,
		preservation,
		repositoryPlan,
		repositoryRules,
		plannedDraft,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	doc := PlanDocument{
		SchemaVersion:  PlanSchemaVersion,
		Repository:     initial.Identity,
		Catalog:        CatalogIdentity{SchemaVersion: CatalogSchemaVersion(), Digest: catalog.Digest()},
		Profile:        profile,
		Decisions:      decisions,
		Retention:      retention,
		ClauseDelta:    clauseDelta,
		Preimages:      snapshot.Preimages,
		Postimages:     postimages,
		Warnings:       cloneFindings(preservation.Warnings),
		SetupManifest:  manifest,
		ManagedEntries: ledger,
	}
	doc.FileChanges, err = deriveFileChanges(doc.ManagedEntries, doc.Preimages, doc.Postimages)
	if err != nil {
		return PlanOutcome{}, err
	}
	doc.PlanDigest, err = computePlanDigest(doc)
	if err != nil {
		return PlanOutcome{}, err
	}
	if err := validatePlanDocumentWithCatalog(doc, catalog); err != nil {
		return PlanOutcome{}, fmt.Errorf("validate assembled Baseline Plan: %w", err)
	}
	result := readyResult(doc.PlanDigest, doc.Warnings, alignment.Verification)
	return PlanOutcome{Plan: &doc, Result: result}, nil
}

func alignmentEvidencePaths(alignment ProfileAlignment) []string {
	seen := make(map[string]struct{})
	add := func(relative string) {
		if safeRelative(relative) {
			seen[relative] = struct{}{}
		}
	}
	for _, capability := range alignment.Capabilities {
		for _, evidence := range capability.Evidence {
			add(evidence.SourcePath)
		}
	}
	for _, route := range alignment.HTTPCandidates {
		add(route.SourcePath)
	}
	if alignment.PostgreSQL.Contract != nil {
		add(alignment.PostgreSQL.Contract.SourcePath)
	}
	for _, evidence := range alignment.PostgreSQL.Implementation {
		add(evidence.SourcePath)
	}
	for _, verification := range alignment.Verification {
		add(verification.DeclarationPath)
	}
	return sortedKeys(seen)
}

func decisionBool(decisions []DecisionValue, id string) bool {
	for _, decision := range decisions {
		if decision.ID == id {
			value, _ := decision.Value.(bool)
			return value
		}
	}
	return false
}

func planSpecificRepository(
	rootPath string,
	proposed []byte,
	enabled bool,
	redistribute bool,
) (specificRepositoryPlan, []Finding, error) {
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return specificRepositoryPlan{}, nil,
			fmt.Errorf("open repository root for repository-specific rules: %w", err)
	}
	defer anchored.Close()

	if redistribute {
		plan := specificRepositoryPlan{
			IncludeRoot: len(proposed) != 0,
		}
		for _, relative := range []string{
			specificRepositoryPath,
			legacyRepositoryPath,
			legacyRepositoryRulesPath,
		} {
			_, exists, finding := readSpecificRepositoryCarrier(anchored, relative)
			if finding != nil {
				return specificRepositoryPlan{}, []Finding{*finding}, nil
			}
			if exists && (relative != specificRepositoryPath || len(proposed) == 0) {
				plan.DeletePaths = append(plan.DeletePaths, relative)
			}
		}
		if len(proposed) != 0 {
			if !enabled {
				return specificRepositoryPlan{}, []Finding{{
					Code:    "baseline.repository-rules.disabled",
					Path:    specificRepositoryPath,
					Message: "Repository-Specific Normative Rules require repository.extension.enabled=true",
				}}, nil
			}
			plan.CanonicalContent = append([]byte(nil), proposed...)
		}
		sort.Strings(plan.DeletePaths)
		return plan, nil, nil
	}

	canonical, canonicalExists, finding := readSpecificRepositoryCarrier(
		anchored,
		specificRepositoryPath,
	)
	if finding != nil {
		return specificRepositoryPlan{}, []Finding{*finding}, nil
	}

	type legacyCarrier struct {
		path    string
		content []byte
	}
	var legacy []legacyCarrier
	var deletePaths []string
	for _, relative := range []string{legacyRepositoryPath, legacyRepositoryRulesPath} {
		content, exists, carrierFinding := readSpecificRepositoryCarrier(anchored, relative)
		if carrierFinding != nil {
			return specificRepositoryPlan{}, []Finding{*carrierFinding}, nil
		}
		if !exists {
			continue
		}
		if repositoryCarrierEmpty(content, relative == legacyRepositoryPath) {
			deletePaths = append(deletePaths, relative)
			continue
		}
		legacy = append(legacy, legacyCarrier{path: relative, content: content})
	}

	var selected []byte
	selectedPath := ""
	if canonicalExists && !repositoryCarrierEmpty(canonical, false) {
		selected = canonical
		selectedPath = specificRepositoryPath
	}
	for _, carrier := range legacy {
		if len(selected) == 0 {
			selected = carrier.content
			selectedPath = carrier.path
		} else if !bytes.Equal(selected, carrier.content) {
			return specificRepositoryPlan{}, []Finding{{
				Code: "baseline.repository-rules.conflict",
				Path: carrier.path,
				Message: fmt.Sprintf(
					"repository-specific rule carriers conflict: %s and %s contain different bytes",
					selectedPath,
					carrier.path,
				),
			}}, nil
		}
		deletePaths = append(deletePaths, carrier.path)
	}
	if len(proposed) != 0 {
		if len(selected) == 0 {
			selected = proposed
			selectedPath = "Baseline Readoption"
		} else if !bytes.Equal(selected, proposed) {
			selected = appendRepositoryRules(selected, proposed)
		}
	}
	if !enabled {
		switch {
		case len(proposed) != 0:
			return specificRepositoryPlan{}, []Finding{{
				Code:    "baseline.repository-rules.disabled",
				Path:    specificRepositoryPath,
				Message: "Repository-Specific Normative Rules require repository.extension.enabled=true",
			}}, nil
		case len(legacy) != 0:
			return specificRepositoryPlan{}, []Finding{{
				Code:    "baseline.repository-rules.migration-disabled",
				Path:    legacy[0].path,
				Message: "legacy Repository-Specific Normative Rules require repository.extension.enabled=true for canonical migration",
			}}, nil
		}
		if canonicalExists && repositoryCarrierEmpty(canonical, false) {
			deletePaths = append(deletePaths, specificRepositoryPath)
		}
		sort.Strings(deletePaths)
		return specificRepositoryPlan{DeletePaths: deletePaths}, nil, nil
	}

	plan := specificRepositoryPlan{
		IncludeRoot: len(selected) != 0,
		DeletePaths: deletePaths,
	}
	if canonicalExists && repositoryCarrierEmpty(canonical, false) && len(selected) == 0 {
		plan.DeletePaths = append(plan.DeletePaths, specificRepositoryPath)
	}
	if len(selected) != 0 {
		plan.CanonicalContent = append([]byte(nil), selected...)
	}
	sort.Strings(plan.DeletePaths)
	return plan, nil, nil
}

func preservationRedistributesRecognizedRepositoryRules(
	preservation RootPreservationPlan,
) bool {
	if preservation.Mode != PreservationModePreservation ||
		len(preservation.Dispositions) == 0 {
		return false
	}
	for _, entry := range preservation.SourceBaseline.Entries {
		switch entry.Path {
		case specificRepositoryPath, legacyRepositoryPath, legacyRepositoryRulesPath:
			return true
		}
	}
	return false
}

func appendRepositoryRules(existing, proposed []byte) []byte {
	if bytes.Contains(existing, proposed) {
		return append([]byte(nil), existing...)
	}
	combined := append([]byte(nil), existing...)
	if len(combined) != 0 &&
		!bytes.HasSuffix(combined, []byte("\n")) &&
		!bytes.HasPrefix(proposed, []byte("\n")) {
		combined = append(combined, '\n')
	}
	return append(combined, proposed...)
}

func readSpecificRepositoryCarrier(
	root *os.Root,
	relative string,
) ([]byte, bool, *Finding) {
	info, err := root.Lstat(filepath.FromSlash(relative))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, &Finding{
			Code:    "baseline.repository-rules.unreadable",
			Path:    relative,
			Message: "repository-specific rule carrier cannot be inspected",
		}
	}
	if !info.Mode().IsRegular() {
		return nil, true, &Finding{
			Code:    "baseline.repository-rules.invalid",
			Path:    relative,
			Message: "repository-specific rule carrier must be a regular non-symlink file",
		}
	}
	content, err := readRootRegularFile(root, relative)
	if err != nil {
		return nil, true, &Finding{
			Code:    "baseline.repository-rules.unreadable",
			Path:    relative,
			Message: "repository-specific rule carrier cannot be read",
		}
	}
	return content, true, nil
}

func repositoryCarrierEmpty(content []byte, allowLegacyScaffold bool) bool {
	if len(bytes.TrimSpace(content)) == 0 {
		return true
	}
	return allowLegacyScaffold &&
		strings.TrimSpace(string(content)) == strings.TrimSpace(legacyRepositoryScaffold)
}

type repositoryRuleSpan struct {
	ID    string
	Start int
	End   int
	Body  []byte
}

var repositoryRuleBeginMarker = regexp.MustCompile(
	`^<!-- roundfix:repository-rule:begin id=([a-z0-9][a-z0-9.-]*) -->\n`,
)

func inventoryRepositoryRuleBlocks(
	root string,
	artifacts []plannedArtifact,
	proposed []RepositoryRuleBlock,
) (repositoryRuleInventory, error) {
	activeGuides := make(map[string]struct{})
	for _, artifact := range artifacts {
		if artifact.Kind == "guide" {
			activeGuides[artifact.Path] = struct{}{}
		}
	}
	inventory := repositoryRuleInventory{
		ByPath: make(map[string][]RepositoryRuleBlock),
	}
	existing := make(map[string]string)
	existingBodies := make(map[string][][]byte)
	for _, relative := range sortedKeys(activeGuides) {
		content, err := readOptionalRegular(root, relative)
		if err != nil {
			return repositoryRuleInventory{}, err
		}
		spans, err := parseRepositoryRuleBlocks(relative, content)
		if err != nil {
			return repositoryRuleInventory{}, err
		}
		for _, span := range spans {
			if previous, duplicate := existing[span.ID]; duplicate {
				return repositoryRuleInventory{}, fmt.Errorf(
					"repository-rule marker %q is duplicated in %q and %q",
					span.ID,
					previous,
					relative,
				)
			}
			existing[span.ID] = relative
			existingBodies[relative] = append(
				existingBodies[relative],
				append([]byte(nil), span.Body...),
			)
			inventory.ByPath[relative] = append(
				inventory.ByPath[relative],
				RepositoryRuleBlock{
					ID:   span.ID,
					Path: relative,
					Body: append([]byte(nil), span.Body...),
				},
			)
		}
	}

	proposedIDs := make(map[string]struct{}, len(proposed))
	for _, block := range proposed {
		if _, active := activeGuides[block.Path]; !active {
			return repositoryRuleInventory{}, fmt.Errorf(
				"repository-rule %q targets inactive semantic guide %q",
				block.ID,
				block.Path,
			)
		}
		if _, duplicate := proposedIDs[block.ID]; duplicate {
			return repositoryRuleInventory{}, fmt.Errorf(
				"repository-rule %q is proposed more than once",
				block.ID,
			)
		}
		proposedIDs[block.ID] = struct{}{}
		if currentPath, exists := existing[block.ID]; exists {
			if currentPath != block.Path {
				return repositoryRuleInventory{}, fmt.Errorf(
					"repository-rule %q already belongs to semantic guide %q",
					block.ID,
					currentPath,
				)
			}
			continue
		}
		var bodyExists bool
		for _, body := range existingBodies[block.Path] {
			if bytes.Equal(body, block.Body) {
				bodyExists = true
				break
			}
		}
		if bodyExists {
			continue
		}
		cloned := block
		cloned.Body = append([]byte(nil), block.Body...)
		inventory.ByPath[block.Path] = append(inventory.ByPath[block.Path], cloned)
	}

	for _, relative := range sortedKeys(inventory.ByPath) {
		for _, block := range inventory.ByPath[relative] {
			if _, planned := proposedIDs[block.ID]; planned {
				continue
			}
			inventory.Retention = append(inventory.Retention, RetentionEvidence{
				FromClause:  "repository-rule." + block.ID,
				Enforcement: "repository-owned",
				Disposition: "repository-document",
				Targets:     []string{relative},
				Reason:      "Retain the current repository-owned semantic rule body outside setup-owned markers.",
			})
		}
	}
	return inventory, nil
}

func parseRepositoryRuleBlocks(relative string, content []byte) ([]repositoryRuleSpan, error) {
	var spans []repositoryRuleSpan
	seen := make(map[string]struct{})
	cursor := 0
	const markerPrefix = "<!-- roundfix:repository-rule:"
	for {
		offset := bytes.Index(content[cursor:], []byte(markerPrefix))
		if offset < 0 {
			break
		}
		start := cursor + offset
		match := repositoryRuleBeginMarker.FindSubmatchIndex(content[start:])
		if match == nil || match[0] != 0 {
			return nil, fmt.Errorf(
				"repository-rule marker in %q is malformed or has no matching begin marker",
				relative,
			)
		}
		id := string(content[start+match[2] : start+match[3]])
		if _, duplicate := seen[id]; duplicate {
			return nil, fmt.Errorf("repository-rule marker %q is duplicated in %q", id, relative)
		}
		seen[id] = struct{}{}
		bodyStart := start + match[1]
		endMarker := []byte(
			"\n<!-- roundfix:repository-rule:end id=" + id + " -->",
		)
		endOffset := bytes.Index(content[bodyStart:], endMarker)
		if endOffset < 0 {
			return nil, fmt.Errorf("repository-rule marker %q in %q is unterminated", id, relative)
		}
		bodyEnd := bodyStart + endOffset
		if nested := bytes.Index(content[bodyStart:bodyEnd], []byte(markerPrefix)); nested >= 0 {
			return nil, fmt.Errorf("repository-rule marker %q in %q contains a nested marker", id, relative)
		}
		end := bodyEnd + len(endMarker)
		if end < len(content) && content[end] == '\n' {
			end++
		}
		for _, entry := range partitionRootSource(relative, content) {
			if entry.Kind == "managed-block" && start < entry.End && end > entry.Start {
				return nil, fmt.Errorf(
					"repository-rule marker %q in %q is inside a setup-owned block",
					id,
					relative,
				)
			}
		}
		spans = append(spans, repositoryRuleSpan{
			ID:    id,
			Start: start,
			End:   end,
			Body:  append([]byte(nil), content[bodyStart:bodyEnd]...),
		})
		cursor = end
	}
	return spans, nil
}

func upsertRepositoryRuleBlocks(
	relative string,
	content []byte,
	blocks []RepositoryRuleBlock,
) ([]byte, error) {
	existing, err := parseRepositoryRuleBlocks(relative, content)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(existing))
	for _, span := range existing {
		seen[span.ID] = struct{}{}
	}
	result := append([]byte(nil), content...)
	for _, block := range blocks {
		if _, exists := seen[block.ID]; exists {
			continue
		}
		if len(result) != 0 {
			if !bytes.HasSuffix(result, []byte("\n")) {
				result = append(result, '\n')
			}
			result = append(result, '\n')
		}
		result = append(result,
			[]byte("<!-- roundfix:repository-rule:begin id="+block.ID+" -->\n")...,
		)
		result = append(result, block.Body...)
		result = append(result,
			[]byte("\n<!-- roundfix:repository-rule:end id="+block.ID+" -->\n")...,
		)
		seen[block.ID] = struct{}{}
	}
	if _, err := parseRepositoryRuleBlocks(relative, result); err != nil {
		return nil, err
	}
	return result, nil
}

func actionOutcome(category, message, next string, warnings []Finding) PlanOutcome {
	if message == "" {
		message = "Baseline planning requires another action"
	}
	return PlanOutcome{Result: Result{
		SchemaVersion:      ResultSchemaVersion,
		Operation:          "plan",
		State:              "action_required",
		Category:           category,
		Message:            message,
		NextAction:         next,
		VerifiedPostimages: []Postimage{},
		Warnings:           cloneFindings(warnings),
		Recommendations:    []string{},
	}}
}

func readyResult(digest string, warnings []Finding, verification []VerificationProjection) Result {
	recommendations := make([]string, 0, len(verification))
	for _, item := range verification {
		if item.RepositoryExecutable && item.Command != "" {
			recommendations = append(recommendations, item.Command)
		}
	}
	sort.Strings(recommendations)
	return Result{
		SchemaVersion:      ResultSchemaVersion,
		Operation:          "plan",
		State:              "ready",
		PlanDigest:         digest,
		VerifiedPostimages: []Postimage{},
		Warnings:           cloneFindings(warnings),
		Recommendations:    recommendations,
	}
}

func resolvePlanRetention(
	root string,
	snapshot RepositorySnapshot,
	catalog *Catalog,
	profileIDs []string,
	targetManifest SetupManifest,
	activeModules []string,
	preservation RootPreservationPlan,
) ([]RetentionEvidence, *ClauseDelta, string, error) {
	retention := readoptionRetentionEvidence(preservation.Dispositions)
	preimage, ok := preimagesByPath(snapshot.Preimages)[manifestPath]
	if !ok || !preimage.Exists {
		return retention, nil, "", nil
	}
	if preimage.Kind != PreimageRegular {
		return nil, nil, "", fmt.Errorf("resolve Upgrade Retention: Setup Manifest is not a regular file")
	}
	data, err := readOptionalRegular(root, manifestPath)
	if err != nil {
		return nil, nil, "", fmt.Errorf("resolve Upgrade Retention: %w", err)
	}
	if planContentIdentity(data) != preimage.ContentIdentity {
		return nil, nil, "", errors.New("resolve Upgrade Retention: Setup Manifest changed during planning")
	}
	manifest, diagnostics := decodeDocument(data, manifestPath)
	if len(diagnostics) != 0 || manifest == nil {
		return retention, nil, "the existing Setup Manifest is not valid strict JSON", nil
	}

	generator, _ := objectValue(manifest["generator"])
	declaredBaseline, _ := stringValue(generator, "baseline")
	compatibleBaselines := make(map[string]struct{}, len(profileIDs))
	for _, profileID := range profileIDs {
		compatibleBaselines["baseline."+profileID+"-"+ManifestVersion] = struct{}{}
	}
	if manifest["schemaVersion"] == ManifestSchema && manifest["version"] == ManifestVersion {
		var existing SetupManifest
		if err := strictJSON(data, &existing); err != nil {
			return retention, nil, "the existing Setup Manifest is not valid strict JSON", nil
		}
		if reflectJSONEqual(existing, targetManifest) {
			return retention, nil, "", nil
		}
		if declaredBaseline == targetManifest.Generator.Baseline {
			if reflectJSONEqual(existing.ManagedArtifacts, targetManifest.ManagedArtifacts) {
				return retention, nil, "", nil
			}
			source, found := catalog.retentionSources[declaredBaseline]
			if !found {
				// Artifact drift proves bytes changed, but without an exact source
				// inventory it does not prove that a managed clause disappeared.
				return retention, nil, "", nil
			}
			evidence, delta := classifySourceClauseTransition(
				source,
				existing.ManagedArtifacts,
				catalog,
				activeModules,
			)
			retention = append(retention, evidence...)
			if unaccounted := clauseDeltaIDs(delta, ClauseUnaccounted); len(unaccounted) != 0 {
				return retention,
					&delta,
					fmt.Sprintf(
						"retention transition has %d unaccounted clause(s): %s",
						len(unaccounted),
						strings.Join(unaccounted, ", "),
					),
					nil
			}
			if len(delta.Dispositions) != 0 {
				return retention, &delta, "", nil
			}
			return retention, nil, "", nil
		}
		if _, compatibleBaseline := compatibleBaselines[declaredBaseline]; compatibleBaseline {
			return retention, nil, "", nil
		}
	}
	if currentSetupManifestProfileIsValid(root, manifest, generator, catalog) {
		return retention, nil, "", nil
	}
	// The portable-v3 manifest remains a supported reader contract while the
	// current writer emits the profile-bound 0.0.1 identity.
	if declaredBaseline == "baseline.portable-v3" {
		return retention, nil, "", nil
	}

	var matches []UpgradeRetentionContract
	fingerprint := legacyManifestFingerprint(manifest)
	for _, transitionID := range catalog.TransitionIDs() {
		contract, err := catalog.UpgradeRetentionContract(transitionID)
		if err != nil {
			return nil, nil, "", err
		}
		matched := declaredBaseline != "" && contract.FromBaseline == declaredBaseline
		if declaredBaseline == "" && fingerprint != "" {
			matched = containsString(contract.LegacyManifestFingerprints, fingerprint)
		}
		if matched {
			matches = append(matches, contract)
		}
	}
	if len(matches) != 1 {
		identity := declaredBaseline
		if identity == "" {
			identity = "manifest fingerprint " + fingerprint
		}
		return retention, nil,
			fmt.Sprintf("the existing Setup Manifest identity %q has no unique maintained transition", identity),
			nil
	}
	return append(retention, transitionRetentionEvidence(matches[0])...), nil, "", nil
}

func compatibleRetentionProfileIDs(
	profileID string,
	draft *ProfileDraftInput,
) []string {
	result := []string{profileID}
	if draft != nil &&
		draft.SourceProfileID != "" &&
		draft.SourceProfileID != profileID {
		result = append(result, draft.SourceProfileID)
	}
	return result
}

func currentSetupManifestProfileIsValid(
	root string,
	manifest document,
	generator document,
	catalog *Catalog,
) bool {
	if manifest["schemaVersion"] != ManifestSchema ||
		manifest["version"] != ManifestVersion ||
		generator["skill"] != "setup-context-driven" ||
		generator["version"] != ManifestVersion ||
		manifest["catalogDigest"] != catalog.Digest() {
		return false
	}
	profileID, profileOK := stringValue(manifest, "profile")
	profileDigest, digestOK := stringValue(manifest, "profileDigest")
	declaredBaseline, baselineOK := stringValue(generator, "baseline")
	if !profileOK || !digestOK || !baselineOK ||
		declaredBaseline != "baseline."+profileID+"-"+ManifestVersion {
		return false
	}
	profile, err := ResolveProfile(root, profileID, catalog)
	return err == nil && profile.Digest == profileDigest
}

func classifySourceClauseTransition(
	source SourceBaseline,
	managedArtifacts []ManifestArtifact,
	catalog *Catalog,
	activeModules []string,
) ([]RetentionEvidence, ClauseDelta) {
	managedCarriers := make(map[string]struct{}, len(managedArtifacts))
	for _, artifact := range managedArtifacts {
		if artifact.Kind == "root-block" || artifact.Kind == "guide" {
			managedCarriers[artifact.Path] = struct{}{}
		}
	}
	current := selectedClauseEnforcement(catalog, activeModules)
	delta := newClauseDelta()
	var evidence []RetentionEvidence
	for _, previous := range source.Entries {
		if previous.Kind != "normative-clause" {
			continue
		}
		if _, managed := managedCarriers[previous.Carrier]; !managed {
			continue
		}
		disposition := ClauseUnaccounted
		targets := []string{}
		reason := "The selected Baseline has no clause with the same identity and enforcement."
		if enforcement, retained := current[previous.ID]; retained && enforcement == previous.Enforcement {
			disposition = ClauseRetained
			targets = append(targets, previous.ID)
			reason = "Stable clause identity and enforcement remain in the selected Baseline."
		}
		delta.Dispositions[previous.ID] = disposition
		delta.Counts[disposition]++
		evidence = append(evidence, RetentionEvidence{
			FromClause:  previous.ID,
			Enforcement: previous.Enforcement,
			Disposition: string(disposition),
			Targets:     targets,
			Reason:      reason,
		})
	}
	return evidence, delta
}

func selectedClauseEnforcement(catalog *Catalog, activeModules []string) map[string]string {
	result := make(map[string]string)
	for _, moduleID := range activeModules {
		module := catalog.modules[moduleID]
		for _, rule := range objectsOrEmpty(module["rules"]) {
			for _, clause := range objectsOrEmpty(rule["clauses"]) {
				id, idOK := stringValue(clause, "id")
				enforcement, enforcementOK := stringValue(clause, "enforcement")
				if idOK && enforcementOK {
					result[id] = enforcement
				}
			}
		}
	}
	return result
}

func newClauseDelta() ClauseDelta {
	counts := make(map[ClauseDisposition]int, 7)
	for _, disposition := range allClauseDispositions() {
		counts[disposition] = 0
	}
	return ClauseDelta{
		Dispositions: make(map[string]ClauseDisposition),
		Counts:       counts,
	}
}

func clauseDeltaIDs(delta ClauseDelta, disposition ClauseDisposition) []string {
	var result []string
	for id, observed := range delta.Dispositions {
		if observed == disposition {
			result = append(result, id)
		}
	}
	sort.Strings(result)
	return result
}

// RenderClauseDeltaBeforeLedger renders the accounted semantic projection
// ahead of an already-rendered file ledger. The ledger bytes are appended
// unchanged, and no clause disposition is inferred from another plan field.
func RenderClauseDeltaBeforeLedger(delta *ClauseDelta, fileLedger string) string {
	if delta == nil || len(delta.Dispositions) == 0 {
		return fileLedger
	}

	var rendered strings.Builder
	fmt.Fprintln(&rendered, "Clause-level semantic delta:")
	fmt.Fprintf(&rendered, "Total clauses: %d\n", len(delta.Dispositions))
	fmt.Fprintln(&rendered, "Disposition counts:")
	for _, disposition := range allClauseDispositions() {
		fmt.Fprintf(&rendered, "- %s: %d\n", disposition, delta.Counts[disposition])
	}
	fmt.Fprintln(&rendered, "Clauses:")
	clauseIDs := make([]string, 0, len(delta.Dispositions))
	for clauseID := range delta.Dispositions {
		clauseIDs = append(clauseIDs, clauseID)
	}
	sort.Strings(clauseIDs)
	for _, clauseID := range clauseIDs {
		fmt.Fprintf(&rendered, "- %s: %s\n", clauseID, delta.Dispositions[clauseID])
	}
	rendered.WriteString(fileLedger)
	return rendered.String()
}

func readoptionRetentionEvidence(dispositions []ReadoptionDisposition) []RetentionEvidence {
	result := make([]RetentionEvidence, 0, len(dispositions))
	for _, disposition := range dispositions {
		targets := []string{}
		if disposition.Destination != nil {
			switch disposition.Disposition {
			case "managed-entry":
				targets = append(targets, disposition.Destination.ManagedID)
			default:
				targets = append(targets, disposition.Destination.Path)
			}
		}
		result = append(result, RetentionEvidence{
			FromClause:  disposition.EntryID,
			Enforcement: disposition.Classification,
			Disposition: disposition.Disposition,
			Targets:     targets,
			Reason:      disposition.Reason,
		})
	}
	return result
}

func transitionRetentionEvidence(contract UpgradeRetentionContract) []RetentionEvidence {
	enforcement := make(map[string]string, len(contract.PriorClauses))
	for _, clause := range contract.PriorClauses {
		enforcement[clause.ID] = clause.Enforcement
	}
	result := make([]RetentionEvidence, len(contract.Accounting))
	for index, accounting := range contract.Accounting {
		result[index] = RetentionEvidence{
			FromClause:  accounting.FromClause,
			Enforcement: enforcement[accounting.FromClause],
			Disposition: accounting.Disposition,
			Targets:     append([]string(nil), accounting.Targets...),
			Reason:      accounting.Reason,
		}
	}
	return result
}

func legacyManifestFingerprint(manifest document) string {
	artifacts := objectsOrEmpty(manifest["managedArtifacts"])
	if len(artifacts) == 0 {
		return ""
	}
	normalized := make([]map[string]any, 0, len(artifacts))
	for _, artifact := range artifacts {
		id, idOK := stringValue(artifact, "id")
		template, templateOK := stringValue(artifact, "template")
		digest, digestOK := stringValue(artifact, "digest")
		version, versionOK := artifact["version"].(json.Number)
		if !idOK || id == "" || !templateOK || template == "" ||
			!digestOK || !isRawSHA256(digest) || !versionOK {
			return ""
		}
		versionValue, err := strconv.Atoi(version.String())
		if err != nil || versionValue < 1 {
			return ""
		}
		normalized = append(normalized, map[string]any{
			"id": id, "version": versionValue, "template": template, "digest": digest,
		})
	}
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i]["id"].(string) < normalized[j]["id"].(string)
	})
	data, err := json.Marshal(normalized)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func normalizePlanDecisions(
	profile ResolvedProfile,
	input []DecisionValue,
	catalog *Catalog,
) ([]DecisionValue, []string, error) {
	values := make(map[string]any, len(profile.Values)+len(input))
	for id, value := range profile.Values {
		values[id] = cloneJSONValue(value)
	}
	for _, decision := range input {
		id := strings.TrimSpace(decision.ID)
		if id == "" {
			return nil, nil, errors.New("normalize Baseline decisions: decision id is required")
		}
		if _, duplicate := values[id]; duplicate {
			return nil, nil, fmt.Errorf("normalize Baseline decisions: duplicate decision %q", id)
		}
		declaration, ok := catalog.decisions[id]
		if !ok {
			return nil, nil, fmt.Errorf("normalize Baseline decisions: unknown decision %q", id)
		}
		if err := validateDecisionValue(declaration, decision.Value); err != nil {
			if id == authProviderDecisionID || id == httpContractDecisionID {
				return nil, nil, fmt.Errorf(
					"normalize Baseline decisions: %w",
					projectDecisionError(fmt.Errorf("normalize %q: %w", id, err)),
				)
			}
			return nil, nil, fmt.Errorf("normalize Baseline decision %q: %w", id, err)
		}
		values[id] = cloneJSONValue(decision.Value)
	}
	required := stringSet(profile.Decisions)
	changed := true
	for changed {
		changed = false
		for id := range required {
			value, answered := values[id]
			if !answered {
				continue
			}
			for _, effect := range objectsOrEmpty(catalog.decisions[id]["effects"]) {
				when, _ := objectValue(effect["when"])
				if !planConditionMatches(when, value) {
					continue
				}
				for _, dependent := range stringsOrEmpty(effect["requireDecisions"]) {
					if _, exists := required[dependent]; !exists {
						required[dependent] = struct{}{}
						changed = true
					}
				}
				for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
					for _, dependent := range stringsOrEmpty(catalog.modules[moduleID]["requiredDecisions"]) {
						if _, exists := required[dependent]; !exists {
							required[dependent] = struct{}{}
							changed = true
						}
					}
				}
			}
		}
	}
	for id := range values {
		if _, ok := required[id]; !ok {
			return nil, nil, fmt.Errorf(
				"normalize Baseline decisions: decision %q is not selected by profile %q",
				id, profile.ID,
			)
		}
	}
	ordered := make([]DecisionValue, 0, len(values))
	missing := make([]string, 0)
	for id := range required {
		value, ok := values[id]
		if !ok {
			missing = append(missing, id)
			continue
		}
		ordered = append(ordered, DecisionValue{ID: id, Value: value})
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	sort.Strings(missing)
	if len(missing) != 0 {
		return ordered, missing, nil
	}
	ordered, err := normalizeProjectDecisions(ordered, catalog)
	if err != nil {
		return nil, nil, fmt.Errorf("normalize Baseline decisions: %w", err)
	}
	return ordered, missing, nil
}

func profileAlignmentDecisions(profile ResolvedProfile, decisions []DecisionValue) []DecisionValue {
	selected := stringSet(profile.Decisions)
	filtered := make([]DecisionValue, 0, len(decisions))
	for _, decision := range decisions {
		if _, ok := selected[decision.ID]; ok {
			filtered = append(filtered, decision)
		}
	}
	return filtered
}

func portableProfilePath(root, profilePath string) string {
	if profilePath == "" {
		return ""
	}
	relative, err := filepath.Rel(root, profilePath)
	if err != nil || !safeRelative(filepath.ToSlash(relative)) {
		return ""
	}
	return filepath.ToSlash(relative)
}

func validateProfileDraftTarget(snapshot RepositorySnapshot, draft plannedProfileDraft) error {
	preimage, ok := preimagesByPath(snapshot.Preimages)[draft.Path]
	if !ok {
		return fmt.Errorf("custom Profile draft target %q has no bounded preimage", draft.Path)
	}
	switch preimage.Kind {
	case PreimageMissing:
		return nil
	case PreimageRegular:
		if preimage.ContentIdentity != planContentIdentity(draft.Content) {
			return fmt.Errorf(
				"custom Profile draft target %q conflicts with existing repository bytes",
				draft.Path,
			)
		}
		return nil
	default:
		return fmt.Errorf("%w: %s is %s", ErrUnsafeCustomProfilePath, draft.Path, preimage.Kind)
	}
}

func resolveManagedArtifacts(
	catalog *Catalog,
	profile ResolvedProfile,
	decisions []DecisionValue,
	includeRepositoryExtension bool,
) ([]string, []plannedArtifact, error) {
	values := make(map[string]any, len(decisions))
	for _, decision := range decisions {
		values[decision.ID] = decision.Value
	}
	controlledModules := make(map[string]struct{})
	controlledArtifacts := make(map[string]struct{})
	included := make(map[string]bool)
	excluded := make(map[string]bool)
	templateOverrides := make(map[string]string)
	renderValues := make(map[string]string)
	for _, declaration := range catalog.decisions {
		for _, effect := range objectsOrEmpty(declaration["effects"]) {
			for _, id := range stringsOrEmpty(effect["activateModules"]) {
				controlledModules[id] = struct{}{}
			}
			for _, field := range []string{"includeArtifacts", "excludeArtifacts"} {
				for _, id := range stringsOrEmpty(effect[field]) {
					controlledArtifacts[id] = struct{}{}
				}
			}
			for _, field := range []string{"selectTemplates", "renderBindings"} {
				for _, item := range objectsOrEmpty(effect[field]) {
					id, _ := stringValue(item, "artifact")
					controlledArtifacts[id] = struct{}{}
				}
			}
		}
	}
	active := make(map[string]bool)
	for _, id := range profile.Modules {
		if _, controlled := controlledModules[id]; !controlled {
			active[id] = true
		}
	}
	decisionOrder := append([]string(nil), profile.Decisions...)
	seenDecision := stringSet(decisionOrder)
	for _, decision := range decisions {
		if _, seen := seenDecision[decision.ID]; !seen {
			decisionOrder = append(decisionOrder, decision.ID)
			seenDecision[decision.ID] = struct{}{}
		}
	}
	var matchedEffects []document
	for _, id := range decisionOrder {
		value, answered := values[id]
		if !answered {
			continue
		}
		declaration := catalog.decisions[id]
		for _, effect := range objectsOrEmpty(declaration["effects"]) {
			when, _ := objectValue(effect["when"])
			if !planConditionMatches(when, value) {
				continue
			}
			matchedEffects = append(matchedEffects, effect)
			for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
				active[moduleID] = true
			}
			for _, artifactID := range stringsOrEmpty(effect["includeArtifacts"]) {
				included[artifactID] = true
			}
			for _, artifactID := range stringsOrEmpty(effect["excludeArtifacts"]) {
				excluded[artifactID] = true
			}
			for _, selection := range objectsOrEmpty(effect["selectTemplates"]) {
				artifactID, _ := stringValue(selection, "artifact")
				templateID, _ := stringValue(selection, "template")
				templateOverrides[artifactID] = templateID
			}
			for _, binding := range objectsOrEmpty(effect["renderBindings"]) {
				token, _ := stringValue(binding, "token")
				rendered, err := renderProjectDecision(id, value, declaration)
				if err != nil {
					return nil, nil, fmt.Errorf("render project decision %q: %w", id, err)
				}
				renderValues[token] = rendered
			}
		}
	}
	if includeRepositoryExtension {
		active[repositoryExtensionModuleID] = true
		included[repositoryExtensionRootID] = true
		delete(excluded, repositoryExtensionRootID)
	} else {
		delete(active, repositoryExtensionModuleID)
		excluded[repositoryExtensionRootID] = true
		delete(included, repositoryExtensionRootID)
	}
	activeModules := make([]string, 0, len(profile.Modules))
	for _, moduleID := range profile.Modules {
		if active[moduleID] {
			activeModules = append(activeModules, moduleID)
		}
	}
	paths := managedArtifactPaths(catalog)
	markerVersion := profileMarkerVersion(catalog, profile.ID)
	declarations := make(map[string]document)
	moduleArtifacts := make(map[string][]string)
	for _, moduleID := range activeModules {
		module := catalog.modules[moduleID]
		for _, declaration := range append(objectsOrEmpty(module["rootBlocks"]), objectsOrEmpty(module["supportingGuides"])...) {
			id, _ := stringValue(declaration, "id")
			declaration = cloneJSONMap(declaration)
			declaration["_module"] = moduleID
			declarations[id] = declaration
			moduleArtifacts[moduleID] = append(moduleArtifacts[moduleID], id)
		}
	}
	var orderedIDs []string
	add := func(id string) {
		if _, exists := declarations[id]; !exists || excluded[id] && !included[id] {
			return
		}
		for _, existing := range orderedIDs {
			if existing == id {
				return
			}
		}
		orderedIDs = append(orderedIDs, id)
	}
	for _, moduleID := range activeModules {
		for _, id := range moduleArtifacts[moduleID] {
			if _, controlled := controlledArtifacts[id]; !controlled {
				add(id)
			}
		}
	}
	for _, effect := range matchedEffects {
		for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
			for _, id := range moduleArtifacts[moduleID] {
				add(id)
			}
		}
		for _, id := range stringsOrEmpty(effect["includeArtifacts"]) {
			add(id)
		}
		for _, selection := range objectsOrEmpty(effect["selectTemplates"]) {
			id, _ := stringValue(selection, "artifact")
			add(id)
		}
		for _, binding := range objectsOrEmpty(effect["renderBindings"]) {
			id, _ := stringValue(binding, "artifact")
			add(id)
		}
	}
	orderedIDs = orderRootArtifacts(catalog, orderedIDs)
	var artifacts []plannedArtifact
	for _, id := range orderedIDs {
		declaration := declarations[id]
		moduleID, _ := stringValue(declaration, "_module")
		templateID, _ := stringValue(declaration, "template")
		if override := templateOverrides[id]; override != "" {
			templateID = override
		}
		template, ok := catalog.TemplateContent(templateID)
		if !ok {
			return nil, nil, fmt.Errorf("render managed entry %q: template %q is missing", id, templateID)
		}
		valuesForArtifact := artifactRenderValues(
			catalog,
			declaration,
			activeModules,
			orderedIDs,
			paths,
		)
		for token, rendered := range renderValues {
			valuesForArtifact[token] = rendered
		}
		body, err := renderTemplate(template.Data, valuesForArtifact)
		if err != nil {
			return nil, nil, fmt.Errorf("render managed entry %q: %w", id, err)
		}
		kind := "root-block"
		if _, ok := declaration["path"]; ok {
			kind = "guide"
		}
		version := markerVersion
		if version == "" {
			version = renderDecisionValue(declaration["version"])
		}
		normalized := strings.TrimSpace(body) + "\n"
		sum := sha256.Sum256([]byte(normalized))
		artifacts = append(artifacts, plannedArtifact{
			ID: id, Path: paths[id], Kind: kind, Module: moduleID,
			Template: templateID, Version: version, Body: body,
			Digest: hex.EncodeToString(sum[:]),
		})
	}
	return activeModules, artifacts, nil
}

// ResolveSemanticOwnerRegistry derives the exact semantic destinations active
// for a resolved Profile and its confirmed project decisions.
func ResolveSemanticOwnerRegistry(
	catalog *Catalog,
	profile ResolvedProfile,
	decisions []DecisionValue,
) (SemanticOwnerRegistry, error) {
	if catalog == nil {
		return nil, errors.New("resolve Semantic Owner Registry: catalog is required")
	}
	modules, artifacts, err := resolveManagedArtifacts(catalog, profile, decisions, false)
	if err != nil {
		return nil, fmt.Errorf("resolve Semantic Owner Registry artifacts: %w", err)
	}
	artifactIDs := make([]string, len(artifacts))
	for index, artifact := range artifacts {
		artifactIDs[index] = artifact.ID
	}
	return catalog.SemanticOwnerRegistry(modules, artifactIDs), nil
}

func orderRootArtifacts(catalog *Catalog, artifactIDs []string) []string {
	rank := make(map[string]int)
	for _, level := range catalog.instructionHierarchy {
		for _, rootID := range level.RootBlocks {
			rank[rootID] = len(rank)
		}
	}
	positions := make([]int, 0)
	roots := make([]string, 0)
	for index, artifactID := range artifactIDs {
		if _, ok := rank[artifactID]; !ok {
			continue
		}
		positions = append(positions, index)
		roots = append(roots, artifactID)
	}
	sort.SliceStable(roots, func(i, j int) bool {
		return rank[roots[i]] < rank[roots[j]]
	})
	ordered := append([]string(nil), artifactIDs...)
	for index, position := range positions {
		ordered[position] = roots[index]
	}
	return ordered
}

func planConditionMatches(condition document, value any) bool {
	if expected, ok := condition["equals"]; ok {
		return valuesEqual(expected, value)
	}
	if expected, ok := condition["present"].(bool); ok {
		return (value != nil) == expected
	}
	return false
}

func valuesEqual(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return bytes.Equal(leftJSON, rightJSON)
}

func renderDecisionValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	default:
		data, _ := json.Marshal(typed)
		return string(data)
	}
}

func profileMarkerVersion(catalog *Catalog, profileID string) string {
	profile := catalog.profiles[profileID]
	if value, ok := profile["markerVersion"].(string); ok {
		return value
	}
	return ""
}

func managedArtifactPaths(catalog *Catalog) map[string]string {
	paths := make(map[string]string)
	for _, module := range catalog.modules {
		for _, block := range objectsOrEmpty(module["rootBlocks"]) {
			id, _ := stringValue(block, "id")
			paths[id] = "AGENTS.md"
		}
		for _, guide := range objectsOrEmpty(module["supportingGuides"]) {
			id, _ := stringValue(guide, "id")
			relative, _ := stringValue(guide, "path")
			paths[id] = relative
		}
	}
	return paths
}

func artifactRenderValues(
	catalog *Catalog,
	artifact document,
	activeModules []string,
	activeArtifacts []string,
	artifactPaths map[string]string,
) map[string]string {
	values := make(map[string]string)
	switch artifact["id"] {
	case "guide.domain":
		values["identifier.strategy"] = ""
	case "guide.backend":
		values["http.contract"] = ""
		values["auth.provider"] = ""
	}
	var rules []string
	moduleID := ""
	for id, module := range catalog.modules {
		for _, candidate := range append(objectsOrEmpty(module["rootBlocks"]), objectsOrEmpty(module["supportingGuides"])...) {
			if candidate["id"] == artifact["id"] {
				moduleID = id
			}
		}
	}
	if moduleID != "" {
		module := catalog.modules[moduleID]
		ruleByID := make(map[string]document)
		for _, rule := range objectsOrEmpty(module["rules"]) {
			id, _ := stringValue(rule, "id")
			ruleByID[id] = rule
		}
		for _, ruleID := range stringsOrEmpty(artifact["rules"]) {
			rule := ruleByID[ruleID]
			clauses := objectsOrEmpty(rule["clauses"])
			if len(clauses) == 0 {
				if guidance, ok := stringValue(rule, "guidance"); ok && strings.TrimSpace(guidance) != "" {
					rules = append(rules, "- "+strings.TrimSpace(guidance))
				}
				continue
			}
			sort.Slice(clauses, func(i, j int) bool {
				left, _ := stringValue(clauses[i], "id")
				right, _ := stringValue(clauses[j], "id")
				return left < right
			})
			for _, clause := range clauses {
				enforcement, _ := stringValue(clause, "enforcement")
				guidance, _ := stringValue(clause, "guidance")
				rules = append(rules, fmt.Sprintf("- **%s**: %s", enforcement, strings.TrimSpace(guidance)))
			}
		}
	}
	if len(rules) != 0 {
		values["artifact.rules"] = strings.Join(rules, "\n\n")
	}
	if artifact["id"] == "guide.skill-dispatch" {
		values["active-modules.skill-dispatch"] = renderSkillDispatch(catalog, activeModules)
	}
	if artifact["id"] == "root.core" {
		values["instruction.hierarchy"] = renderInstructionHierarchy(catalog, activeArtifacts)
	}
	for _, reference := range objectsOrEmpty(artifact["references"]) {
		token, _ := stringValue(reference, "token")
		targetID, _ := stringValue(reference, "managedId")
		if target := artifactPaths[targetID]; target != "" {
			values[token] = "`" + target + "`"
			continue
		}
		ownership, _ := stringValue(reference, "ownership")
		repositoryPath, _ := stringValue(reference, "path")
		if ownership == "repository" && safeRelative(repositoryPath) {
			values[token] = "`" + repositoryPath + "`"
		}
	}
	return values
}

func renderInstructionHierarchy(catalog *Catalog, activeArtifacts []string) string {
	active := stringSet(activeArtifacts)
	lines := []string{
		"### Instruction hierarchy",
		"",
		"Apply active guidance in this order. A narrower guide may add constraints for its concern but cannot weaken a universal Normative Clause or confirmed project decision.",
		"",
	}
	ordinal := 1
	for _, level := range catalog.instructionHierarchy {
		levelActive := false
		for _, rootID := range level.RootBlocks {
			if _, ok := active[rootID]; ok {
				levelActive = true
				break
			}
		}
		if !levelActive {
			continue
		}
		lines = append(lines, fmt.Sprintf("%d. **%s**", ordinal, level.Title))
		ordinal++
	}
	return strings.Join(lines, "\n")
}

func renderSkillDispatch(catalog *Catalog, activeModules []string) string {
	active := stringSet(activeModules)
	var lines []string
	if asset, ok := catalog.Asset("skill-activations.json"); ok {
		raw, diagnostics := decodeDocument(asset.Data, asset.Path)
		if len(diagnostics) == 0 && raw != nil {
			bundles := make(map[string][]string)
			for _, bundle := range objectsOrEmpty(raw["bundles"]) {
				id, _ := stringValue(bundle, "id")
				bundles[id] = stringsOrEmpty(bundle["skills"])
			}
			var activationLines []string
			for _, activation := range objectsOrEmpty(raw["activations"]) {
				owner, _ := stringValue(activation, "owner")
				if _, ok := active[owner]; !ok {
					continue
				}
				id, _ := stringValue(activation, "id")
				when, _ := stringValue(activation, "when")
				bundleID, _ := stringValue(activation, "bundle")
				skills := bundles[bundleID]
				rendered := make([]string, len(skills))
				for index, skill := range skills {
					rendered[index] = "`" + skill + "`"
				}
				activationLines = append(activationLines,
					fmt.Sprintf("- `%s`: %s", id, when),
					fmt.Sprintf("  - `%s`: %s", bundleID, strings.Join(rendered, ", ")),
				)
			}
			if len(activationLines) != 0 {
				lines = append(lines, "Exact activation bundles:", "")
				lines = append(lines, activationLines...)
				lines = append(lines, "", "Individual skill triggers:", "")
			}
		}
	}
	dispatchBySkill := make(map[string][]document)
	for _, moduleID := range activeModules {
		module := catalog.modules[moduleID]
		for _, dispatch := range objectsOrEmpty(module["skillDispatch"]) {
			skill, _ := stringValue(dispatch, "skill")
			for _, trigger := range objectsOrEmpty(dispatch["triggers"]) {
				dispatchBySkill[skill] = append(dispatchBySkill[skill], trigger)
			}
		}
	}
	for _, skill := range sortedKeys(dispatchBySkill) {
		triggers := dispatchBySkill[skill]
		sort.Slice(triggers, func(i, j int) bool {
			left, _ := stringValue(triggers[i], "id")
			right, _ := stringValue(triggers[j], "id")
			return left < right
		})
		lines = append(lines, "- `"+skill+"`:")
		for _, trigger := range triggers {
			id, _ := stringValue(trigger, "id")
			when, _ := stringValue(trigger, "when")
			lines = append(lines, fmt.Sprintf("  - `%s`: %s", id, when))
		}
	}
	return strings.Join(lines, "\n")
}

func renderTemplate(data []byte, values map[string]string) (string, error) {
	content := string(data)
	for {
		start := strings.Index(content, "{{")
		if start < 0 {
			break
		}
		endOffset := strings.Index(content[start+2:], "}}")
		if endOffset < 0 {
			return "", errors.New("template contains an unterminated token")
		}
		end := start + 2 + endOffset
		token := content[start+2 : end]
		replacement, ok := values[token]
		if !ok {
			return "", fmt.Errorf("missing render value for token %q", token)
		}
		content = content[:start] + replacement + content[end+2:]
	}
	return content, nil
}

func buildSetupManifest(
	catalog *Catalog,
	profile ResolvedProfile,
	decisions []DecisionValue,
	modules []string,
	artifacts []plannedArtifact,
	verification []VerificationProjection,
) SetupManifest {
	manifestDecisions := make(map[string]ManifestDecision, len(decisions))
	for _, decision := range decisions {
		manifestDecisions[decision.ID] = ManifestDecision{Value: cloneJSONValue(decision.Value)}
	}
	managed := make([]ManifestArtifact, len(artifacts))
	for index, artifact := range artifacts {
		managed[index] = ManifestArtifact{
			ID: artifact.ID, Path: artifact.Path, Kind: artifact.Kind,
			Module: artifact.Module, Template: artifact.Template,
			Version: artifact.Version, Digest: artifact.Digest,
		}
	}
	return SetupManifest{
		SchemaVersion: ManifestSchema, Version: ManifestVersion,
		Generator: ManifestGenerator{
			Skill: "setup-context-driven", Version: ManifestVersion,
			Baseline: "baseline." + profile.ID + "-" + ManifestVersion,
		},
		Profile: profile.ID, ProfileDigest: profile.Digest, CatalogDigest: catalog.Digest(),
		Modules: append([]string(nil), modules...), Decisions: manifestDecisions,
		ManagedArtifacts: managed, LocalSkills: []string{},
		Verification: append([]VerificationProjection{}, verification...),
	}
}

func marshalSetupManifestBytes(manifest SetupManifest) ([]byte, error) {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	return compactSetupManifestMethodArrays(data), nil
}

func compactSetupManifestMethodArrays(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	result := make([]string, 0, len(lines))
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		if strings.TrimSpace(line) != `"methods": [` {
			result = append(result, line)
			continue
		}
		var methods []string
		end := index + 1
		for ; end < len(lines); end++ {
			trimmed := strings.TrimSpace(lines[end])
			if trimmed == "]" || trimmed == "]," {
				break
			}
			var method string
			if err := json.Unmarshal([]byte(strings.TrimSuffix(trimmed, ",")), &method); err != nil {
				methods = nil
				break
			}
			methods = append(methods, method)
		}
		if methods == nil || end >= len(lines) {
			result = append(result, line)
			continue
		}
		encodedMethods := make([]string, 0, len(methods))
		for _, method := range methods {
			encoded, err := json.Marshal(method)
			if err != nil {
				encodedMethods = nil
				break
			}
			encodedMethods = append(encodedMethods, string(encoded))
		}
		if encodedMethods == nil {
			result = append(result, line)
			continue
		}
		suffix := ""
		if strings.TrimSpace(lines[end]) == "]," {
			suffix = ","
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " "))]
		result = append(
			result,
			indent+`"methods": [`+strings.Join(encodedMethods, ", ")+"]"+suffix,
		)
		index = end
	}
	return []byte(strings.Join(result, "\n"))
}

func assemblePostimages(
	root string,
	snapshot RepositorySnapshot,
	artifacts []plannedArtifact,
	manifestBytes []byte,
	preservation RootPreservationPlan,
	repositoryPlan specificRepositoryPlan,
	repositoryRules repositoryRuleInventory,
	profileDraft *plannedProfileDraft,
) ([]Postimage, []ManagedEntry, error) {
	byPath := make(map[string][]plannedArtifact)
	for _, artifact := range artifacts {
		byPath[artifact.Path] = append(byPath[artifact.Path], artifact)
	}
	outputs := make(map[string][]byte)
	removedRepositoryPointer := false
	for relative, grouped := range byPath {
		current, err := readOptionalRegular(root, relative)
		if err != nil {
			return nil, nil, err
		}
		content := string(current)
		if preservationConsumesRootPath(preservation, relative) {
			content = ""
		}
		if relative == "AGENTS.md" && !repositoryPlan.IncludeRoot {
			withoutPointer := removeManagedBlock(string(current), repositoryExtensionRootID)
			removedRepositoryPointer = withoutPointer != string(current)
			withoutPointer = removeManagedBlock(content, repositoryExtensionRootID)
			content = withoutPointer
		}
		for _, artifact := range grouped {
			content = upsertManagedBlock(content, artifact)
		}
		rendered, err := upsertRepositoryRuleBlocks(
			relative,
			[]byte(content),
			repositoryRules.ByPath[relative],
		)
		if err != nil {
			return nil, nil, err
		}
		outputs[relative] = rendered
	}
	outputs[manifestPath] = append([]byte(nil), manifestBytes...)
	for _, backup := range preservation.Backups {
		source, err := readOptionalRegular(root, backup.SourcePath)
		if err != nil {
			return nil, nil, err
		}
		if rendered, exists := outputs[backup.SourcePath]; exists &&
			bytes.Equal(source, rendered) &&
			unchangedSetupManagedGuidance(backup.ContentIdentity, rendered) {
			continue
		}
		outputs[backup.Path] = source
	}
	if len(repositoryPlan.CanonicalContent) != 0 {
		outputs[specificRepositoryPath] = append([]byte(nil), repositoryPlan.CanonicalContent...)
	}
	if profileDraft != nil {
		outputs[profileDraft.Path] = append([]byte(nil), profileDraft.Content...)
	}

	preimages := preimagesByPath(snapshot.Preimages)
	pathSet := make(map[string]struct{}, len(outputs)+len(repositoryPlan.DeletePaths))
	for relative := range outputs {
		pathSet[relative] = struct{}{}
	}
	for _, relative := range repositoryPlan.DeletePaths {
		pathSet[relative] = struct{}{}
	}
	paths := sortedKeys(pathSet)
	postimages := make([]Postimage, 0, len(paths))
	ledger := make([]ManagedEntry, 0)
	for _, relative := range paths {
		before := preimages[relative].ContentIdentity
		postimage := Postimage{Path: relative, Kind: PreimageMissing}
		if content, exists := outputs[relative]; exists {
			postimage = Postimage{
				Path: relative, Kind: PreimageRegular, Mode: 0o644,
				Content: content, ContentIdentity: planContentIdentity(content),
			}
		}
		postimages = append(postimages, postimage)
		action := fileAction(
			preimages[relative],
			postimage.Kind,
			before,
			postimage.ContentIdentity,
		)
		grouped := byPath[relative]
		for _, artifact := range grouped {
			ledger = append(ledger, ManagedEntry{
				ID: artifact.ID, Path: relative, Action: action, Kind: artifact.Kind,
				Module: artifact.Module, Template: artifact.Template, Version: artifact.Version,
				BeforeIdentity: before, AfterIdentity: postimage.ContentIdentity,
				ContentIdentity: "sha256:" + artifact.Digest,
			})
		}
		if content, exists := outputs[relative]; exists &&
			len(repositoryRules.ByPath[relative]) != 0 {
			blocks, err := parseRepositoryRuleBlocks(relative, content)
			if err != nil {
				return nil, nil, err
			}
			for _, block := range blocks {
				ledger = append(ledger, ManagedEntry{
					ID:              "repository-rule:" + block.ID,
					Path:            relative,
					Action:          action,
					Kind:            "repository-owned",
					BeforeIdentity:  before,
					AfterIdentity:   postimage.ContentIdentity,
					ContentIdentity: planContentIdentity(block.Body),
				})
			}
		}
		if relative == manifestPath {
			ledger = append(ledger, ManagedEntry{
				ID: "manifest", Path: relative, Action: action, Kind: "manifest",
				Version: ManifestVersion, BeforeIdentity: before,
				AfterIdentity:   postimage.ContentIdentity,
				ContentIdentity: postimage.ContentIdentity,
			})
		}
		if profileDraft != nil && relative == profileDraft.Path {
			ledger = append(ledger, ManagedEntry{
				ID:   "profile:" + profileDraft.ID,
				Path: relative, Action: action, Kind: "profile",
				BeforeIdentity: before, AfterIdentity: postimage.ContentIdentity,
				ContentIdentity: postimage.ContentIdentity,
			})
		}
		for _, backup := range preservation.Backups {
			if backup.Path == relative {
				ledger = append(ledger, ManagedEntry{
					ID: "backup:" + backup.SourcePath, Path: relative, Action: action, Kind: "backup",
					BeforeIdentity: before, AfterIdentity: postimage.ContentIdentity,
					ContentIdentity: backup.ContentIdentity,
				})
			}
		}
		if relative == specificRepositoryPath && len(repositoryPlan.CanonicalContent) != 0 {
			ledger = append(ledger, ManagedEntry{
				ID: "repository-rules.canonicalize", Path: relative, Action: action,
				Kind: "repository-owned", BeforeIdentity: before,
				AfterIdentity:   postimage.ContentIdentity,
				ContentIdentity: postimage.ContentIdentity,
			})
		}
		if containsString(repositoryPlan.DeletePaths, relative) {
			ledger = append(ledger, ManagedEntry{
				ID: "repository-rules.remove:" + relative, Path: relative, Action: action,
				Kind: "repository-owned", BeforeIdentity: before,
			})
		}
		if relative == "AGENTS.md" && removedRepositoryPointer {
			ledger = append(ledger, ManagedEntry{
				ID: "repository-rules.remove-root-pointer", Path: relative, Action: action,
				Kind: "repository-owned", BeforeIdentity: before,
				AfterIdentity:   postimage.ContentIdentity,
				ContentIdentity: postimage.ContentIdentity,
			})
		}
	}
	for index := range ledger {
		ledger[index].Ordinal = index
	}
	return postimages, ledger, nil
}

func preservationConsumesRootPath(preservation RootPreservationPlan, relative string) bool {
	if preservation.Mode != PreservationModePreservation ||
		preservation.State != PreservationStateReady {
		return false
	}
	if _, consumed := preservation.consumedRootPaths[relative]; consumed {
		return relative == "AGENTS.md"
	}
	if len(preservation.Dispositions) == 0 {
		return false
	}
	for _, entry := range preservation.SourceBaseline.Entries {
		if entry.Path == relative {
			return relative == "AGENTS.md"
		}
	}
	return false
}

func readOptionalRegular(root, relative string) ([]byte, error) {
	target := filepath.Join(root, filepath.FromSlash(relative))
	info, err := os.Lstat(target)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect planned path %q: %w", relative, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("inspect planned path %q: expected a regular file or missing path", relative)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("read planned path %q: %w", relative, err)
	}
	return data, nil
}

func removeManagedBlock(current, id string) string {
	begin := fmt.Sprintf("<!-- setup-context-driven:begin id=%s ", id)
	end := fmt.Sprintf("<!-- setup-context-driven:end id=%s -->", id)
	start := strings.Index(current, begin)
	if start < 0 {
		return current
	}
	finish := strings.Index(current[start:], end)
	if finish < 0 {
		return current
	}
	finish += start + len(end)
	if finish < len(current) && current[finish] == '\n' {
		finish++
	}
	result := current[:start] + current[finish:]
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}
	return strings.TrimLeft(result, "\n")
}

func upsertManagedBlock(current string, artifact plannedArtifact) string {
	begin := fmt.Sprintf("<!-- setup-context-driven:begin id=%s ", artifact.ID)
	end := fmt.Sprintf("<!-- setup-context-driven:end id=%s -->", artifact.ID)
	block := fmt.Sprintf(
		"<!-- setup-context-driven:begin id=%s version=%s -->\n\n%s\n\n<!-- setup-context-driven:end id=%s -->\n",
		artifact.ID, artifact.Version, strings.TrimSpace(artifact.Body), artifact.ID,
	)
	start := strings.Index(current, begin)
	if start >= 0 {
		finish := strings.Index(current[start:], end)
		if finish >= 0 {
			finish += start + len(end)
			if finish < len(current) && current[finish] == '\n' {
				finish++
			}
			return current[:start] + block + current[finish:]
		}
	}
	if current == "" {
		return block
	}
	if !strings.HasSuffix(current, "\n") {
		current += "\n"
	}
	return current + "\n" + block
}

func planContentIdentity(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func preimagesByPath(preimages []Preimage) map[string]Preimage {
	result := make(map[string]Preimage, len(preimages))
	for _, preimage := range preimages {
		result[preimage.Path] = preimage
	}
	return result
}

func fileAction(before Preimage, afterKind PreimageKind, beforeIdentity, afterIdentity string) string {
	if !before.Exists && afterKind != PreimageMissing {
		return "create"
	}
	if before.Exists && afterKind == PreimageMissing {
		return "delete"
	}
	if beforeIdentity == afterIdentity {
		return "unchanged"
	}
	return "update"
}

func deriveFileChanges(
	ledger []ManagedEntry,
	preimages []Preimage,
	postimages []Postimage,
) ([]FileChange, error) {
	preimageByPath := preimagesByPath(preimages)
	postimageByPath := make(map[string]Postimage, len(postimages))
	idsByPath := make(map[string][]string)
	for index, entry := range ledger {
		if entry.Ordinal != index {
			return nil, fmt.Errorf("managed-entry ordinal %d appears at index %d", entry.Ordinal, index)
		}
		if !safeRelative(entry.Path) {
			return nil, fmt.Errorf("managed entry %q has unsafe path %q", entry.ID, entry.Path)
		}
		idsByPath[entry.Path] = append(idsByPath[entry.Path], entry.ID)
	}
	for _, postimage := range postimages {
		if _, duplicate := postimageByPath[postimage.Path]; duplicate {
			return nil, fmt.Errorf("duplicate postimage path %q", postimage.Path)
		}
		postimageByPath[postimage.Path] = postimage
	}
	paths := sortedKeys(idsByPath)
	changes := make([]FileChange, 0, len(paths))
	for _, relative := range paths {
		postimage, ok := postimageByPath[relative]
		if !ok {
			return nil, fmt.Errorf("managed ledger path %q has no postimage", relative)
		}
		before := preimageByPath[relative]
		if before.Exists && before.Kind == postimage.Kind &&
			before.ContentIdentity == postimage.ContentIdentity {
			continue
		}
		changes = append(changes, FileChange{
			Path:           relative,
			Action:         fileAction(before, postimage.Kind, before.ContentIdentity, postimage.ContentIdentity),
			BeforeIdentity: before.ContentIdentity,
			AfterIdentity:  postimage.ContentIdentity,
			ManagedEntries: append([]string(nil), idsByPath[relative]...),
		})
	}
	return changes, nil
}

func computePlanDigest(document PlanDocument) (string, error) {
	payload := document
	payload.PlanDigest = ""
	payload.FileChanges = nil
	normalized, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("serialize Baseline Plan digest payload: %w", err)
	}
	sum := sha256.Sum256(append([]byte(planDigestDomain), normalized...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// ValidatePlanDocument checks ordering, projection, exact postimages, and the
// Plan Digest without consulting a checkout.
func ValidatePlanDocument(document PlanDocument) error {
	if err := validatePlanDocumentShape(document); err != nil {
		return err
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return fmt.Errorf("load Baseline catalog for plan validation: %w", err)
	}
	return validatePlanDocumentAgainstCatalog(document, catalog)
}

func validatePlanDocumentWithCatalog(document PlanDocument, catalog *Catalog) error {
	if err := validatePlanDocumentShape(document); err != nil {
		return err
	}
	return validatePlanDocumentAgainstCatalog(document, catalog)
}

func validatePlanDocumentShape(document PlanDocument) error {
	if document.SchemaVersion != PlanSchemaVersion {
		return fmt.Errorf("unsupported Baseline Plan schema %q", document.SchemaVersion)
	}
	if document.Decisions == nil || document.Retention == nil ||
		document.Preimages == nil || document.Postimages == nil ||
		document.Warnings == nil || document.ManagedEntries == nil ||
		document.FileChanges == nil {
		return errors.New("Baseline Plan collections must be JSON arrays")
	}
	if err := validatePlanRepositoryIdentity(document.Repository); err != nil {
		return err
	}
	return nil
}

func validatePlanDocumentAgainstCatalog(document PlanDocument, catalog *Catalog) error {
	if document.Catalog.SchemaVersion != CatalogSchemaVersion() ||
		document.Catalog.Digest != catalog.Digest() {
		return errors.New("Baseline Plan catalog identity does not match the embedded catalog")
	}
	profileDigest, err := profileDigest(document.Profile)
	if err != nil {
		return err
	}
	if document.Profile.Digest != profileDigest {
		return errors.New("Baseline Plan profile digest mismatch")
	}
	if document.Profile.Path != "" && !safeRelative(document.Profile.Path) {
		return fmt.Errorf("unsafe Baseline Profile path %q", document.Profile.Path)
	}
	for index, decision := range document.Decisions {
		if decision.ID == "" || index > 0 && document.Decisions[index-1].ID >= decision.ID {
			return errors.New("Baseline Plan decisions must have unique IDs in lexical order")
		}
	}
	normalizedProjectDecisions, err := normalizeProjectDecisions(document.Decisions, catalog)
	if err != nil {
		return fmt.Errorf("validate Baseline Plan project decisions: %w", err)
	}
	if !reflectJSONEqual(normalizedProjectDecisions, document.Decisions) {
		return errors.New("Baseline Plan project decisions are not normalized or derived")
	}
	for index, retention := range document.Retention {
		if retention.FromClause == "" || retention.Enforcement == "" ||
			retention.Disposition == "" || retention.Targets == nil ||
			strings.TrimSpace(retention.Reason) == "" {
			return fmt.Errorf("invalid retention entry %d", index)
		}
	}
	if err := validateReadyClauseDelta(document.ClauseDelta, document.Retention); err != nil {
		return err
	}
	preimagePaths := make(map[string]struct{}, len(document.Preimages))
	for index, preimage := range document.Preimages {
		if !safeRelative(preimage.Path) {
			return fmt.Errorf("unsafe preimage path %q", preimage.Path)
		}
		if _, duplicate := preimagePaths[preimage.Path]; duplicate {
			return fmt.Errorf("duplicate preimage path %q", preimage.Path)
		}
		if index > 0 && document.Preimages[index-1].Path >= preimage.Path {
			return errors.New("Baseline Plan preimages must be in lexical path order")
		}
		preimagePaths[preimage.Path] = struct{}{}
	}
	postimagePaths := make(map[string]struct{}, len(document.Postimages))
	for index, postimage := range document.Postimages {
		if !safeRelative(postimage.Path) {
			return fmt.Errorf("unsafe postimage path %q", postimage.Path)
		}
		if _, duplicate := postimagePaths[postimage.Path]; duplicate {
			return fmt.Errorf("duplicate postimage path %q", postimage.Path)
		}
		if index > 0 && document.Postimages[index-1].Path >= postimage.Path {
			return errors.New("Baseline Plan postimages must be in lexical path order")
		}
		postimagePaths[postimage.Path] = struct{}{}
		if _, bounded := preimagePaths[postimage.Path]; !bounded {
			return fmt.Errorf("postimage %q has no bounded preimage", postimage.Path)
		}
		if postimage.Kind == PreimageRegular && planContentIdentity(postimage.Content) != postimage.ContentIdentity {
			return fmt.Errorf("postimage %q content identity mismatch", postimage.Path)
		}
	}
	ledgerPaths := make(map[string]struct{}, len(document.ManagedEntries))
	for index, entry := range document.ManagedEntries {
		if entry.Ordinal != index {
			return fmt.Errorf("managed entry %q has ordinal %d, want %d", entry.ID, entry.Ordinal, index)
		}
		if entry.ID == "" || !safeRelative(entry.Path) {
			return fmt.Errorf("managed entry %d has invalid identity or path", index)
		}
		ledgerPaths[entry.Path] = struct{}{}
	}
	for relative := range postimagePaths {
		if _, represented := ledgerPaths[relative]; !represented {
			return fmt.Errorf("postimage %q has no canonical managed entry", relative)
		}
	}
	if err := validateSetupManifest(document, catalog); err != nil {
		return err
	}
	projection, err := deriveFileChanges(document.ManagedEntries, document.Preimages, document.Postimages)
	if err != nil {
		return err
	}
	expected, _ := json.Marshal(projection)
	actual, _ := json.Marshal(document.FileChanges)
	if !bytes.Equal(expected, actual) {
		return errors.New("fileChanges does not match the canonical managed-entry ledger")
	}
	digest, err := computePlanDigest(document)
	if err != nil {
		return err
	}
	if document.PlanDigest != digest {
		return fmt.Errorf("Baseline Plan digest mismatch: got %q, want %q", document.PlanDigest, digest)
	}
	if err := validatePlanApplyContract(document); err != nil {
		return err
	}
	return nil
}

func validateReadyClauseDelta(delta *ClauseDelta, retention []RetentionEvidence) error {
	if delta == nil {
		return nil
	}
	if len(delta.Dispositions) == 0 || len(delta.Counts) != len(allClauseDispositions()) {
		return errors.New("Baseline Plan clause delta is empty or incomplete")
	}
	allowed := make(map[ClauseDisposition]struct{}, len(allClauseDispositions()))
	calculated := make(map[ClauseDisposition]int, len(allClauseDispositions()))
	for _, disposition := range allClauseDispositions() {
		allowed[disposition] = struct{}{}
	}
	retentionByClause := make(map[string]RetentionEvidence, len(retention))
	for _, entry := range retention {
		retentionByClause[entry.FromClause] = entry
	}
	for clauseID, disposition := range delta.Dispositions {
		if clauseID == "" {
			return errors.New("Baseline Plan clause delta has an empty clause ID")
		}
		if _, ok := allowed[disposition]; !ok {
			return fmt.Errorf("Baseline Plan clause %q has invalid disposition %q", clauseID, disposition)
		}
		if disposition == ClauseUnaccounted {
			return fmt.Errorf("Baseline Plan clause %q is unaccounted", clauseID)
		}
		entry, ok := retentionByClause[clauseID]
		if !ok || entry.Disposition != string(disposition) {
			return fmt.Errorf("Baseline Plan clause %q has no matching retention record", clauseID)
		}
		calculated[disposition]++
	}
	for _, disposition := range allClauseDispositions() {
		if delta.Counts[disposition] != calculated[disposition] {
			return fmt.Errorf("Baseline Plan clause disposition %q count is invalid", disposition)
		}
	}
	return nil
}

func validatePlanRepositoryIdentity(identity RepositoryIdentity) error {
	if identity.SchemaVersion != RepositoryIdentitySchemaVersion ||
		(identity.ObjectFormat != "sha1" && identity.ObjectFormat != "sha256") ||
		len(identity.RootCommits) == 0 {
		return errors.New("Baseline Plan repository identity is invalid")
	}
	objectLength := 40
	if identity.ObjectFormat == "sha256" {
		objectLength = 64
	}
	for index, commit := range identity.RootCommits {
		if len(commit) != objectLength || !isLowerHex(commit) ||
			index > 0 && identity.RootCommits[index-1] >= commit {
			return errors.New("Baseline Plan repository root commits are invalid or unordered")
		}
	}
	payload := struct {
		SchemaVersion string   `json:"schemaVersion"`
		ObjectFormat  string   `json:"objectFormat"`
		RootCommits   []string `json:"rootCommits"`
	}{
		SchemaVersion: identity.SchemaVersion,
		ObjectFormat:  identity.ObjectFormat,
		RootCommits:   identity.RootCommits,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("serialize repository identity validation payload: %w", err)
	}
	sum := sha256.Sum256(append([]byte(repositoryIdentityDigestDomain), data...))
	if identity.Digest != "sha256:"+hex.EncodeToString(sum[:]) {
		return errors.New("Baseline Plan repository identity digest mismatch")
	}
	return nil
}

func validateSetupManifest(document PlanDocument, catalog *Catalog) error {
	manifest := document.SetupManifest
	if manifest.SchemaVersion != ManifestSchema || manifest.Version != ManifestVersion ||
		manifest.Generator.Skill != "setup-context-driven" ||
		manifest.Generator.Version != ManifestVersion ||
		manifest.Generator.Baseline != "baseline."+document.Profile.ID+"-"+ManifestVersion ||
		manifest.Profile != document.Profile.ID ||
		manifest.ProfileDigest != document.Profile.Digest ||
		manifest.CatalogDigest != document.Catalog.Digest {
		return errors.New("Setup Manifest identity does not match the Baseline Plan")
	}
	if manifest.Modules == nil || manifest.Decisions == nil ||
		manifest.ManagedArtifacts == nil || manifest.LocalSkills == nil ||
		manifest.Verification == nil {
		return errors.New("Setup Manifest collections must be arrays or objects")
	}
	includeRepositoryExtension := containsString(manifest.Modules, repositoryExtensionModuleID)
	if includeRepositoryExtension &&
		!decisionBool(document.Decisions, "repository.extension.enabled") {
		return errors.New("Setup Manifest enables repository-specific rules without owner approval")
	}
	modules, artifacts, err := resolveManagedArtifacts(
		catalog,
		document.Profile,
		document.Decisions,
		includeRepositoryExtension,
	)
	if err != nil {
		return fmt.Errorf("resolve Setup Manifest validation inventory: %w", err)
	}
	if !reflectJSONEqual(manifest.Modules, modules) {
		return errors.New("Setup Manifest modules do not match the resolved Decision Plan")
	}
	if len(manifest.Decisions) != len(document.Decisions) {
		return errors.New("Setup Manifest decisions do not match the Baseline Plan")
	}
	for _, decision := range document.Decisions {
		stored, ok := manifest.Decisions[decision.ID]
		if !ok || !valuesEqual(stored.Value, decision.Value) {
			return fmt.Errorf("Setup Manifest decision %q does not match the Baseline Plan", decision.ID)
		}
	}
	artifactsByID := make(map[string]ManifestArtifact, len(manifest.ManagedArtifacts))
	for _, artifact := range manifest.ManagedArtifacts {
		if _, duplicate := artifactsByID[artifact.ID]; duplicate {
			return fmt.Errorf("duplicate Setup Manifest managed artifact %q", artifact.ID)
		}
		artifactsByID[artifact.ID] = artifact
	}
	managedCount := 0
	for _, entry := range document.ManagedEntries {
		if entry.Kind != "root-block" && entry.Kind != "guide" {
			continue
		}
		managedCount++
		artifact, ok := artifactsByID[entry.ID]
		if !ok ||
			artifact.Path != entry.Path ||
			artifact.Kind != entry.Kind ||
			artifact.Module != entry.Module ||
			artifact.Template != entry.Template ||
			artifact.Version != entry.Version ||
			"sha256:"+artifact.Digest != entry.ContentIdentity {
			return fmt.Errorf("Setup Manifest managed artifact %q does not match the ledger", entry.ID)
		}
	}
	if managedCount != len(manifest.ManagedArtifacts) {
		return errors.New("Setup Manifest managed artifacts do not match the ledger")
	}
	if len(artifacts) != len(manifest.ManagedArtifacts) {
		return errors.New("Setup Manifest managed artifacts do not match the resolved Decision Plan")
	}
	for index, artifact := range artifacts {
		stored := manifest.ManagedArtifacts[index]
		if stored.ID != artifact.ID ||
			stored.Path != artifact.Path ||
			stored.Kind != artifact.Kind ||
			stored.Module != artifact.Module ||
			stored.Template != artifact.Template ||
			stored.Version != artifact.Version ||
			stored.Digest != artifact.Digest {
			return fmt.Errorf("Setup Manifest managed artifact %d does not match the resolved Decision Plan", index)
		}
	}
	return nil
}

func reflectJSONEqual(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}

// ValidatePlanRepository accepts the same plan in any matching clone whose
// stable identity and every bounded preimage still match.
func ValidatePlanRepository(ctx context.Context, repository string, document PlanDocument) error {
	if err := ValidatePlanDocument(document); err != nil {
		return err
	}
	return validatePlanRepositoryState(ctx, repository, document)
}

func validatePlanRepositoryWithCatalog(
	ctx context.Context,
	repository string,
	document PlanDocument,
	catalog *Catalog,
) error {
	if catalog == nil {
		return errors.New("validate Baseline Plan repository with catalog: catalog is required")
	}
	if err := validatePlanDocumentWithCatalog(document, catalog); err != nil {
		return err
	}
	return validatePlanRepositoryState(ctx, repository, document)
}

func validatePlanRepositoryState(
	ctx context.Context,
	repository string,
	document PlanDocument,
) error {
	root, identity, err := inspectRepositoryIdentity(ctx, repository, nil)
	if err != nil {
		return err
	}
	if identity.Digest != document.Repository.Digest {
		return errors.New("Baseline Plan repository identity does not match")
	}
	paths := make([]string, len(document.Preimages))
	for index, preimage := range document.Preimages {
		paths[index] = preimage.Path
	}
	snapshot, err := inspectRepositorySnapshot(root, InventoryRequest{MutablePaths: paths})
	if err != nil {
		return err
	}
	current := preimagesByPath(snapshot.Preimages)
	for _, approved := range document.Preimages {
		observed, ok := current[approved.Path]
		if !ok || !samePreimage(approved, observed) {
			return fmt.Errorf("Baseline Plan preimage is stale at %q", approved.Path)
		}
	}
	return nil
}

func samePreimage(left, right Preimage) bool {
	return left.Path == right.Path &&
		left.Exists == right.Exists &&
		left.Kind == right.Kind &&
		left.Mode == right.Mode &&
		left.LinkTarget == right.LinkTarget &&
		left.Bytes == right.Bytes &&
		left.ContentIdentity == right.ContentIdentity
}

func cloneFindings(findings []Finding) []Finding {
	if findings == nil {
		return []Finding{}
	}
	return append([]Finding{}, findings...)
}
