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
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const (
	CapabilityRecheckSchemaVersion = "roundfix/baseline-capability-recheck/v1"

	maxCapabilityPaths     = 32
	maxCapabilityFileBytes = 1024 * 1024
	maxExecutableLinkHops  = 64
	maxHTTPSourceFiles     = 256
	maxHTTPSourceBytes     = 2 * 1024 * 1024

	executableProbeReasonNotFound      = "not-found"
	executableProbeReasonBrokenLink    = "broken-link"
	executableProbeReasonLinkCycle     = "link-cycle"
	executableProbeReasonNotExecutable = "not-executable"
)

// ErrNoResolvableProfile names a capability re-check that cannot select a
// current Baseline Profile from either the request or the Setup Manifest.
var ErrNoResolvableProfile = errors.New("capability re-check: no resolvable Baseline Profile")

type ProfileAlignmentState string

const (
	ProfileAlignmentReady          ProfileAlignmentState = "ready"
	ProfileAlignmentActionRequired ProfileAlignmentState = "action_required"
)

type CapabilityRequirement string

const (
	CapabilityRequired    CapabilityRequirement = "required"
	CapabilityRecommended CapabilityRequirement = "recommended"
	CapabilityOptional    CapabilityRequirement = "optional"
)

type CapabilityEvidenceKind string

const (
	CapabilityEvidenceDeclaredFile   CapabilityEvidenceKind = "declared-file"
	CapabilityEvidenceExecutable     CapabilityEvidenceKind = "executable"
	CapabilityEvidenceInstalledSkill CapabilityEvidenceKind = "installed-skill"
)

type CapabilityEvidenceStrength string

const (
	CapabilityEvidenceNone       CapabilityEvidenceStrength = "none"
	CapabilityEvidenceDeclared   CapabilityEvidenceStrength = "declared"
	CapabilityEvidenceDiscovered CapabilityEvidenceStrength = "discovered"
	CapabilityEvidenceVerified   CapabilityEvidenceStrength = "verified"
)

type CapabilityEvidenceStatus string

const (
	CapabilityEvidenceAbsent  CapabilityEvidenceStatus = "absent"
	CapabilityEvidenceInvalid CapabilityEvidenceStatus = "invalid"
	CapabilityEvidencePresent CapabilityEvidenceStatus = "present"
)

type CapabilityStatus string

const (
	CapabilityInsufficient CapabilityStatus = "insufficient"
	CapabilityMissing      CapabilityStatus = "missing"
	CapabilitySatisfied    CapabilityStatus = "satisfied"
)

type CapabilityEvidenceClassification string

const (
	EvidenceProfileRequirement CapabilityEvidenceClassification = "profile-requirement"
	EvidenceRepositoryContract CapabilityEvidenceClassification = "repository-contract"
	EvidenceImplementation     CapabilityEvidenceClassification = "implementation"
)

type RepositoryCapability struct {
	ID              string
	Title           string
	Requirement     CapabilityRequirement
	EvidenceKind    CapabilityEvidenceKind
	MinimumEvidence CapabilityEvidenceStrength
	Probe           map[string]any
	Explanation     string
	NextAction      string
}

type CapabilityEvidence struct {
	Status         CapabilityEvidenceStatus         `json:"status"`
	Kind           CapabilityEvidenceKind           `json:"kind"`
	Strength       CapabilityEvidenceStrength       `json:"strength"`
	Classification CapabilityEvidenceClassification `json:"classification"`
	SourcePath     string                           `json:"sourcePath,omitempty"`
	SourceDigest   string                           `json:"sourceDigest,omitempty"`
	Detail         string                           `json:"detail,omitempty"`
}

type CapabilityDiagnostic struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	NextAction string `json:"nextAction,omitempty"`
}

type CapabilityOutcome struct {
	ID              string                     `json:"id"`
	Title           string                     `json:"title"`
	Requirement     CapabilityRequirement      `json:"requirement"`
	EvidenceKind    CapabilityEvidenceKind     `json:"evidenceKind"`
	MinimumEvidence CapabilityEvidenceStrength `json:"minimumEvidence"`
	Status          CapabilityStatus           `json:"status"`
	Blocking        bool                       `json:"blocking"`
	Evidence        []CapabilityEvidence       `json:"evidence"`
	Diagnostic      CapabilityDiagnostic       `json:"diagnostic"`
	Explanation     string                     `json:"explanation"`
}

type executableProbeResult struct {
	Candidate string
	Resolved  string
	Reason    string
	HopCount  int
}

type ProfileDivergenceGroup string

const (
	ProfileDivergenceBlocking      ProfileDivergenceGroup = "blocking"
	ProfileDivergenceAdvisory      ProfileDivergenceGroup = "advisory"
	ProfileDivergenceInformational ProfileDivergenceGroup = "informational"

	advisoryNonBlockingStatement = "This advisory does not block readiness or apply."
)

type ProfileCapabilityResolution struct {
	SelectedTechnology    string   `json:"selectedTechnology"`
	RepositoryRemediation string   `json:"repositoryRemediation"`
	ProfileAdaptation     string   `json:"profileAdaptation"`
	RemovedDecisions      []string `json:"removedDecisions"`
}

type ProfileDivergence struct {
	Code                 string                       `json:"code"`
	ID                   string                       `json:"id"`
	Requirement          CapabilityRequirement        `json:"requirement"`
	Group                ProfileDivergenceGroup       `json:"group,omitempty"`
	Blocking             bool                         `json:"blocking"`
	Message              string                       `json:"message"`
	NonBlockingStatement string                       `json:"nonBlockingStatement,omitempty"`
	NextAction           string                       `json:"nextAction,omitempty"`
	Probe                map[string]any               `json:"probe,omitempty"`
	Evidence             []CapabilityEvidence         `json:"evidence,omitempty"`
	CapabilityResolution *ProfileCapabilityResolution `json:"capabilityResolution,omitempty"`
}

// HTTPRouteCandidate records only locally observed route facts. Repository
// policy fields deliberately do not exist in this projection.
type HTTPRouteCandidate struct {
	SourcePath   string   `json:"sourcePath"`
	SourceDigest string   `json:"sourceDigest"`
	Scope        string   `json:"scope"`
	Methods      []string `json:"methods"`
}

type PostgreSQLEvidence struct {
	AcceptedContractPaths []string             `json:"acceptedContractPaths"`
	Contract              *CapabilityEvidence  `json:"contract,omitempty"`
	Implementation        []CapabilityEvidence `json:"implementation"`
}

type VerificationClassification string

const (
	VerificationProfileExpectation VerificationClassification = "profile-expectation"
	VerificationRepositoryCommand  VerificationClassification = "repository-command"
)

type VerificationProjection struct {
	ID                   string                     `json:"id"`
	Role                 string                     `json:"role"`
	Tool                 string                     `json:"tool,omitempty"`
	Command              string                     `json:"command"`
	SatisfiedByCommand   string                     `json:"satisfiedByCommand,omitempty"`
	Classification       VerificationClassification `json:"classification"`
	RepositoryExecutable bool                       `json:"repositoryExecutable"`
	DeclarationPath      string                     `json:"declarationPath,omitempty"`
	DeclarationDigest    string                     `json:"declarationDigest,omitempty"`
}

type ProfileAlignmentRequest struct {
	ProfileID                string
	Decisions                []DecisionValue
	Profile                  *ResolvedProfile
	RemediationProfileID     string
	VerificationRoleMappings map[string]string
}

// ProfileAlignment is the deterministic, read-only result consumed by later
// planning and interaction layers.
type ProfileAlignment struct {
	State          ProfileAlignmentState    `json:"state"`
	Ready          bool                     `json:"ready"`
	Profile        ResolvedProfile          `json:"profile"`
	Decisions      []DecisionValue          `json:"decisions"`
	Capabilities   []CapabilityOutcome      `json:"capabilities"`
	Divergences    []ProfileDivergence      `json:"divergences"`
	HTTPCandidates []HTTPRouteCandidate     `json:"httpRouteCandidates"`
	PostgreSQL     PostgreSQLEvidence       `json:"postgresql"`
	Verification   []VerificationProjection `json:"verification"`
}

// CapabilityRecheckRequest selects the repository and, optionally, an exact
// Baseline Profile. When ProfileID is empty, the current Setup Manifest owns
// selection. Decisions are intentionally not an input to this operation.
type CapabilityRecheckRequest struct {
	Repository string
	ProfileID  string
}

// CapabilityRecheckProfile identifies the exact Profile evaluated by a
// capability re-check without projecting its decision declarations or values.
type CapabilityRecheckProfile struct {
	ID     string                `json:"id"`
	Source BaselineProfileSource `json:"source"`
	Digest string                `json:"digest"`
}

// CapabilityRecheckResult is the read-only capability projection returned to
// automation. It deliberately omits decisions and every write-bearing plan
// field.
type CapabilityRecheckResult struct {
	SchemaVersion string                    `json:"schemaVersion"`
	Operation     string                    `json:"operation"`
	State         string                    `json:"state"`
	Category      string                    `json:"category,omitempty"`
	Message       string                    `json:"message,omitempty"`
	NextAction    string                    `json:"nextAction,omitempty"`
	Profile       *CapabilityRecheckProfile `json:"profile,omitempty"`
	Capabilities  []CapabilityOutcome       `json:"capabilities"`
	Divergences   []ProfileDivergence       `json:"divergences"`
	PostgreSQL    PostgreSQLEvidence        `json:"postgresql"`
}

type profileCapabilityEvaluation struct {
	capabilities []CapabilityOutcome
	divergences  []ProfileDivergence
	postgres     PostgreSQLEvidence
}

var (
	httpMethodCall = regexp.MustCompile(`(?i)\.\s*(get|post|put|patch|delete|options|head|all)\s*\(\s*["'\x60]([^"'` + "\x60" + `\r\n]+)["'\x60]`)
	httpScopeCall  = regexp.MustCompile(`(?i)\.\s*(?:route|basePath)\s*\(\s*["'\x60]([^"'` + "\x60" + `\r\n]+)["'\x60]`)
	makeTargetName = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
)

var universalCapabilities = []RepositoryCapability{
	{
		ID:              "capability.context7",
		Title:           "Context7",
		Requirement:     CapabilityRequired,
		EvidenceKind:    CapabilityEvidenceInstalledSkill,
		MinimumEvidence: CapabilityEvidenceVerified,
		Probe:           map[string]any{"skill": "context7"},
		Explanation:     "Context7 provides current authoritative library and API documentation.",
		NextAction:      "Add the context7 skill to the Repository Skill Set, then rerun capability evaluation.",
	},
	{
		ID:              "capability.exa",
		Title:           "Exa",
		Requirement:     CapabilityRequired,
		EvidenceKind:    CapabilityEvidenceInstalledSkill,
		MinimumEvidence: CapabilityEvidenceVerified,
		Probe:           map[string]any{"skill": "exa-web-search"},
		Explanation:     "Exa provides broad-source searches when local and authoritative sources are insufficient.",
		NextAction:      "Add the exa-web-search skill to the Repository Skill Set, then rerun capability evaluation.",
	},
	{
		ID:              "capability.firecrawl",
		Title:           "Firecrawl",
		Requirement:     CapabilityRecommended,
		EvidenceKind:    CapabilityEvidenceInstalledSkill,
		MinimumEvidence: CapabilityEvidenceVerified,
		Probe:           map[string]any{"skill": "firecrawl"},
		Explanation:     "Firecrawl provides structured web-content extraction for external research.",
		NextAction:      "Add the firecrawl skill if structured web extraction is useful.",
	},
	{
		ID:              "capability.rg",
		Title:           "rg",
		Requirement:     CapabilityRecommended,
		EvidenceKind:    CapabilityEvidenceExecutable,
		MinimumEvidence: CapabilityEvidenceDiscovered,
		Probe:           map[string]any{"executable": "rg"},
		Explanation:     "rg provides fast bounded local repository search.",
		NextAction:      "Install rg and expose it on PATH if faster local search is useful.",
	},
	{
		ID:              "capability.rtk",
		Title:           "rtk",
		Requirement:     CapabilityRecommended,
		EvidenceKind:    CapabilityEvidenceExecutable,
		MinimumEvidence: CapabilityEvidenceDiscovered,
		Probe:           map[string]any{"executable": "rtk"},
		Explanation:     "rtk keeps command evidence compact without changing command behavior.",
		NextAction:      "Install rtk and expose it on PATH if compact command output is useful.",
	},
}

var evidenceRank = map[CapabilityEvidenceStrength]int{
	CapabilityEvidenceNone:       0,
	CapabilityEvidenceDeclared:   1,
	CapabilityEvidenceDiscovered: 2,
	CapabilityEvidenceVerified:   3,
}

// ResolveProfileAlignment evaluates one selected Baseline Profile from
// bounded local evidence. It does not execute repository commands, use the
// network, or mutate repository bytes.
func ResolveProfileAlignment(
	ctx context.Context,
	repoRoot string,
	request ProfileAlignmentRequest,
	catalog *Catalog,
) (ProfileAlignment, error) {
	if err := ctx.Err(); err != nil {
		return ProfileAlignment{}, err
	}
	if catalog == nil {
		return ProfileAlignment{}, errors.New("resolve profile alignment: catalog is required")
	}
	profileID := strings.TrimSpace(request.ProfileID)
	if profileID == "" {
		return ProfileAlignment{}, errors.New("resolve profile alignment: select exactly one Baseline Profile")
	}
	rootPath, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return ProfileAlignment{}, fmt.Errorf("resolve profile alignment repository: %w", err)
	}
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return ProfileAlignment{}, fmt.Errorf("open profile alignment repository: %w", err)
	}
	defer root.Close()

	var profile ResolvedProfile
	if request.Profile != nil {
		profile = *request.Profile
		if profile.ID != profileID {
			return ProfileAlignment{}, fmt.Errorf(
				"resolve selected Baseline Profile %q: in-memory profile identity is %q",
				profileID,
				profile.ID,
			)
		}
	} else {
		profile, err = ResolveProfile(rootPath, profileID, catalog)
		if err != nil {
			return ProfileAlignment{}, fmt.Errorf("resolve selected Baseline Profile %q: %w", profileID, err)
		}
	}
	decisions, decisionDivergences, err := normalizeAlignmentDecisions(profile, request.Decisions, catalog)
	if err != nil {
		return ProfileAlignment{}, err
	}
	remediationProfileID := strings.TrimSpace(request.RemediationProfileID)
	if remediationProfileID == "" {
		remediationProfileID = profile.ID
	}
	capabilityEvaluation, err := evaluateProfileCapabilities(
		ctx,
		root,
		profile,
		remediationProfileID,
		catalog,
	)
	if err != nil {
		return ProfileAlignment{}, err
	}
	divergences := append([]ProfileDivergence(nil), decisionDivergences...)
	divergences = append(divergences, capabilityEvaluation.divergences...)

	httpCandidates, err := collectHTTPRouteCandidates(ctx, root, profile, catalog)
	if err != nil {
		return ProfileAlignment{}, err
	}
	verification, verificationDivergences, err := resolveVerificationProjection(
		root,
		profile,
		decisions,
		request.VerificationRoleMappings,
		catalog,
	)
	if err != nil {
		return ProfileAlignment{}, err
	}
	divergences = append(divergences, verificationDivergences...)
	classifyProfileDivergences(divergences)
	sortProfileDivergences(divergences)

	ready := true
	for _, divergence := range divergences {
		if divergence.Blocking {
			ready = false
			break
		}
	}
	state := ProfileAlignmentReady
	if !ready {
		state = ProfileAlignmentActionRequired
	}
	return ProfileAlignment{
		State:          state,
		Ready:          ready,
		Profile:        profile,
		Decisions:      decisions,
		Capabilities:   capabilityEvaluation.capabilities,
		Divergences:    divergences,
		HTTPCandidates: httpCandidates,
		PostgreSQL:     capabilityEvaluation.postgres,
		Verification:   verification,
	}, nil
}

// RecheckCapabilities evaluates Repository Capabilities without accepting or
// resolving decisions and without constructing a Change Plan. It shares the
// capability evaluator and divergence projection used by full planning.
func RecheckCapabilities(
	ctx context.Context,
	request CapabilityRecheckRequest,
) (CapabilityRecheckResult, error) {
	if err := ctx.Err(); err != nil {
		return CapabilityRecheckResult{}, err
	}
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return CapabilityRecheckResult{}, err
	}
	rootPath, err := cleanRepositoryRoot(request.Repository)
	if err != nil {
		return CapabilityRecheckResult{}, fmt.Errorf("re-check capabilities repository: %w", err)
	}
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return CapabilityRecheckResult{}, fmt.Errorf("open capability re-check repository: %w", err)
	}
	defer root.Close()

	profile, err := resolveCapabilityRecheckProfile(root, rootPath, request.ProfileID, catalog)
	if err != nil {
		return CapabilityRecheckResult{}, err
	}
	evaluation, err := evaluateProfileCapabilities(ctx, root, profile, profile.ID, catalog)
	if err != nil {
		return CapabilityRecheckResult{}, err
	}
	state := ProfileAlignmentReady
	for _, divergence := range evaluation.divergences {
		if divergence.Blocking {
			state = ProfileAlignmentActionRequired
			break
		}
	}
	return CapabilityRecheckResult{
		SchemaVersion: CapabilityRecheckSchemaVersion,
		Operation:     "capabilities.check",
		State:         string(state),
		Profile: &CapabilityRecheckProfile{
			ID: profile.ID, Source: profile.Source, Digest: profile.Digest,
		},
		Capabilities: evaluation.capabilities,
		Divergences:  evaluation.divergences,
		PostgreSQL:   evaluation.postgres,
	}, nil
}

func resolveCapabilityRecheckProfile(
	root *os.Root,
	rootPath string,
	requestedID string,
	catalog *Catalog,
) (ResolvedProfile, error) {
	profileID := strings.TrimSpace(requestedID)
	if profileID == "" {
		data, err := readRootRegularFile(root, manifestPath)
		if err != nil {
			return ResolvedProfile{}, fmt.Errorf("%w: current Setup Manifest is unavailable", ErrNoResolvableProfile)
		}
		manifest, valid := parseManagedSetupManifest(data)
		if !valid {
			return ResolvedProfile{}, fmt.Errorf("%w: current Setup Manifest is invalid", ErrNoResolvableProfile)
		}
		profileID = manifest.Profile
	}
	profile, err := ResolveProfile(rootPath, profileID, catalog)
	if err != nil {
		return ResolvedProfile{}, fmt.Errorf("%w %q: %w", ErrNoResolvableProfile, profileID, err)
	}
	return profile, nil
}

func evaluateProfileCapabilities(
	ctx context.Context,
	root *os.Root,
	profile ResolvedProfile,
	remediationProfileID string,
	catalog *Catalog,
) (profileCapabilityEvaluation, error) {
	capabilities, err := resolvedProfileCapabilities(profile, catalog)
	if err != nil {
		return profileCapabilityEvaluation{}, err
	}
	outcomes, postgres, err := evaluateRepositoryCapabilities(ctx, root, capabilities)
	if err != nil {
		return profileCapabilityEvaluation{}, err
	}
	applyUniversalCapabilityRemediation(outcomes, remediationProfileID)
	divergences := make([]ProfileDivergence, 0)
	for index, outcome := range outcomes {
		if outcome.Status == CapabilitySatisfied {
			continue
		}
		divergences = append(divergences, ProfileDivergence{
			Code:                 outcome.Diagnostic.Code,
			ID:                   outcome.ID,
			Requirement:          outcome.Requirement,
			Blocking:             outcome.Blocking,
			Message:              outcome.Diagnostic.Message,
			NextAction:           outcome.Diagnostic.NextAction,
			Probe:                capabilities[index].Probe,
			Evidence:             outcome.Evidence,
			CapabilityResolution: stackCapabilityResolution(profile, outcome, catalog),
		})
	}
	classifyProfileDivergences(divergences)
	sortProfileDivergences(divergences)
	return profileCapabilityEvaluation{
		capabilities: outcomes,
		divergences:  divergences,
		postgres:     postgres,
	}, nil
}

func applyUniversalCapabilityRemediation(
	outcomes []CapabilityOutcome,
	profileID string,
) {
	skills := map[string]string{
		"capability.context7": "context7",
		"capability.exa":      "exa-web-search",
	}
	for index := range outcomes {
		skill, universal := skills[outcomes[index].ID]
		if !universal || outcomes[index].Status == CapabilitySatisfied {
			continue
		}
		outcomes[index].Diagnostic.NextAction = fmt.Sprintf(
			"Run roundfix baseline skills restore --profile %s --skill %s, review its exact restoration Plan Digest, rerun with --confirm-plan <digest>, then rerun profile alignment.",
			profileID,
			skill,
		)
	}
}

func stackCapabilityResolution(
	profile ResolvedProfile,
	outcome CapabilityOutcome,
	catalog *Catalog,
) *ProfileCapabilityResolution {
	if outcome.Requirement != CapabilityRequired ||
		!strings.HasPrefix(outcome.ID, "capability.stack.") ||
		!slices.Contains(profile.Capabilities, outcome.ID) {
		return nil
	}
	retainedCapabilities := make([]string, 0, len(profile.Capabilities)-1)
	for _, capabilityID := range profile.Capabilities {
		if capabilityID != outcome.ID {
			retainedCapabilities = append(retainedCapabilities, capabilityID)
		}
	}
	beforeRemoval := profileAdaptationDecisions(
		profile.Decisions,
		profile.Modules,
		profile.Capabilities,
		catalog,
	)
	retainedDecisions := stringSet(profileAdaptationDecisions(
		profile.Decisions,
		profile.Modules,
		retainedCapabilities,
		catalog,
	))
	removedDecisions := make([]string, 0)
	for _, decisionID := range beforeRemoval {
		if _, retained := retainedDecisions[decisionID]; !retained {
			removedDecisions = append(removedDecisions, decisionID)
		}
	}
	return &ProfileCapabilityResolution{
		SelectedTechnology:    outcome.Title,
		RepositoryRemediation: outcome.Diagnostic.NextAction,
		ProfileAdaptation: fmt.Sprintf(
			"Remove %s through a reviewed repository-owned Profile adaptation.",
			outcome.Title,
		),
		RemovedDecisions: removedDecisions,
	}
}

func classifyProfileDivergences(divergences []ProfileDivergence) {
	for index := range divergences {
		divergences[index].Group = profileDivergenceGroup(divergences[index].Requirement)
		if divergences[index].Group == ProfileDivergenceAdvisory {
			divergences[index].NonBlockingStatement = advisoryNonBlockingStatement
		}
	}
}

func profileDivergenceGroup(requirement CapabilityRequirement) ProfileDivergenceGroup {
	switch requirement {
	case CapabilityRequired:
		return ProfileDivergenceBlocking
	case CapabilityRecommended:
		return ProfileDivergenceAdvisory
	default:
		return ProfileDivergenceInformational
	}
}

func normalizeAlignmentDecisions(
	profile ResolvedProfile,
	input []DecisionValue,
	catalog *Catalog,
) ([]DecisionValue, []ProfileDivergence, error) {
	selected := stringSet(profile.Decisions)
	values := make(map[string]any, len(profile.Values)+len(input))
	for id, value := range profile.Values {
		values[id] = cloneJSONValue(value)
	}
	for _, decision := range input {
		id := strings.TrimSpace(decision.ID)
		if _, ok := selected[id]; !ok {
			return nil, nil, fmt.Errorf("resolve profile alignment decision %q: not selected by profile %q", id, profile.ID)
		}
		if _, duplicate := values[id]; duplicate {
			return nil, nil, fmt.Errorf("resolve profile alignment decision %q: duplicate answer", id)
		}
		declaration, ok := catalog.decisions[id]
		if !ok {
			return nil, nil, fmt.Errorf("resolve profile alignment decision %q: catalog declaration is missing", id)
		}
		if err := validateDecisionValue(declaration, decision.Value); err != nil {
			if id == authProviderDecisionID || id == httpContractDecisionID {
				return nil, nil, fmt.Errorf(
					"resolve profile alignment: %w",
					projectDecisionError(fmt.Errorf("normalize %q: %w", id, err)),
				)
			}
			return nil, nil, fmt.Errorf("resolve profile alignment decision %q: %w", id, err)
		}
		values[id] = cloneJSONValue(decision.Value)
	}

	decisions := make([]DecisionValue, 0, len(values))
	divergences := make([]ProfileDivergence, 0)
	for _, id := range profile.Decisions {
		value, ok := values[id]
		if !ok {
			divergences = append(divergences, ProfileDivergence{
				Code:        "profile.decision.required",
				ID:          id,
				Requirement: CapabilityRequired,
				Blocking:    true,
				Message:     fmt.Sprintf("selected profile decision %q requires an explicit answer", id),
				NextAction:  fmt.Sprintf("answer %s or select a different Baseline Profile", id),
			})
			continue
		}
		decisions = append(decisions, DecisionValue{ID: id, Value: value})
	}
	sort.Slice(decisions, func(i, j int) bool { return decisions[i].ID < decisions[j].ID })
	if len(divergences) == 0 {
		var err error
		decisions, err = normalizeProjectDecisions(decisions, catalog)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve profile alignment: %w", err)
		}
	}
	return decisions, divergences, nil
}

func cloneJSONValue(value any) any {
	data, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var cloned any
	if err := json.Unmarshal(data, &cloned); err != nil {
		return value
	}
	return cloned
}

func resolvedProfileCapabilities(profile ResolvedProfile, catalog *Catalog) ([]RepositoryCapability, error) {
	capabilities := append([]RepositoryCapability(nil), universalCapabilities...)
	if len(profile.Capabilities) == 0 {
		sort.Slice(capabilities, func(i, j int) bool { return capabilities[i].ID < capabilities[j].ID })
		return capabilities, nil
	}

	declarations := make(map[string]RepositoryCapability)
	for _, profileID := range catalog.ProfileIDs() {
		raw := catalog.profiles[profileID]
		for _, value := range objectsOrEmpty(raw["capabilities"]) {
			capability, err := parseProfileCapability(profileID, value)
			if err != nil {
				return nil, err
			}
			declarations[capability.ID] = capability
		}
	}
	for _, id := range profile.Capabilities {
		capability, ok := declarations[id]
		if !ok {
			return nil, fmt.Errorf("resolve profile alignment capability %q: catalog declaration is missing", id)
		}
		capabilities = append(capabilities, capability)
	}
	sort.Slice(capabilities, func(i, j int) bool { return capabilities[i].ID < capabilities[j].ID })
	return capabilities, nil
}

func parseProfileCapability(profileID string, raw document) (RepositoryCapability, error) {
	id, idOK := stringValue(raw, "id")
	title, titleOK := stringValue(raw, "title")
	requirementText, requirementOK := stringValue(raw, "strength")
	probe, probeOK := objectValue(raw["probe"])
	kindText, kindOK := stringValue(probe, "kind")
	if !idOK || !titleOK || !requirementOK || !probeOK || !kindOK {
		return RepositoryCapability{}, fmt.Errorf("resolve profile %q capability declaration: required field is missing", profileID)
	}
	requirement := CapabilityRequirement(requirementText)
	if requirement != CapabilityRequired && requirement != CapabilityRecommended && requirement != CapabilityOptional {
		return RepositoryCapability{}, fmt.Errorf("resolve profile %q capability %q: unsupported requirement %q", profileID, id, requirement)
	}
	kind := CapabilityEvidenceKind(kindText)
	if kind != CapabilityEvidenceDeclaredFile && kind != CapabilityEvidenceExecutable {
		return RepositoryCapability{}, fmt.Errorf("resolve profile %q capability %q: unsupported evidence kind %q", profileID, id, kind)
	}
	minimum := CapabilityEvidenceDeclared
	if kind == CapabilityEvidenceExecutable {
		minimum = CapabilityEvidenceDiscovered
	}
	return RepositoryCapability{
		ID:              id,
		Title:           title,
		Requirement:     requirement,
		EvidenceKind:    kind,
		MinimumEvidence: minimum,
		Probe:           cloneJSONMap(probe),
		Explanation:     fmt.Sprintf("%s is %s for the %s profile.", title, requirement, profileID),
		NextAction:      fmt.Sprintf("Add compatible local %s evidence and rerun profile alignment.", title),
	}, nil
}

func evaluateRepositoryCapabilities(
	ctx context.Context,
	root *os.Root,
	capabilities []RepositoryCapability,
) ([]CapabilityOutcome, PostgreSQLEvidence, error) {
	outcomes := make([]CapabilityOutcome, 0, len(capabilities))
	postgres := PostgreSQLEvidence{}
	for _, capability := range capabilities {
		if err := ctx.Err(); err != nil {
			return nil, PostgreSQLEvidence{}, err
		}
		var outcome CapabilityOutcome
		var err error
		if capability.ID == "capability.stack.postgresql" {
			outcome, postgres, err = evaluatePostgreSQLCapability(root, capability)
		} else {
			evidence := collectCapabilityEvidence(root, capability)
			outcome = evaluateCapability(capability, evidence)
		}
		if err != nil {
			return nil, PostgreSQLEvidence{}, err
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, postgres, nil
}

func collectCapabilityEvidence(root *os.Root, capability RepositoryCapability) []CapabilityEvidence {
	switch capability.EvidenceKind {
	case CapabilityEvidenceDeclaredFile:
		return collectDeclaredFileEvidence(root, capability)
	case CapabilityEvidenceInstalledSkill:
		return []CapabilityEvidence{collectInstalledSkillEvidence(root, capability)}
	case CapabilityEvidenceExecutable:
		return []CapabilityEvidence{collectExecutableEvidence(capability)}
	default:
		return []CapabilityEvidence{invalidCapabilityEvidence(capability.EvidenceKind, "unsupported capability evidence kind")}
	}
}

func collectDeclaredFileEvidence(root *os.Root, capability RepositoryCapability) []CapabilityEvidence {
	paths := stringsFromAny(capability.Probe["paths"])
	contains, containsOK := capability.Probe["contains"].(string)
	if len(paths) == 0 || len(paths) > maxCapabilityPaths || !containsOK {
		return []CapabilityEvidence{invalidCapabilityEvidence(capability.EvidenceKind, "declared file probe is invalid")}
	}
	slices.Sort(paths)
	evidence := make([]CapabilityEvidence, 0, len(paths))
	for _, relative := range paths {
		data, state, detail := readBoundedRegularFile(root, relative, maxCapabilityFileBytes)
		if state == CapabilityEvidenceInvalid {
			item := invalidCapabilityEvidence(capability.EvidenceKind, detail)
			item.SourcePath = relative
			return append(evidence, item)
		}
		if state != CapabilityEvidencePresent {
			evidence = append(evidence, CapabilityEvidence{
				Status:         CapabilityEvidenceAbsent,
				Kind:           capability.EvidenceKind,
				Strength:       CapabilityEvidenceNone,
				Classification: EvidenceProfileRequirement,
				SourcePath:     relative,
				Detail:         "file not found",
			})
			continue
		}
		if !bytes.Contains(data, []byte(contains)) {
			evidence = append(evidence, CapabilityEvidence{
				Status:         CapabilityEvidenceAbsent,
				Kind:           capability.EvidenceKind,
				Strength:       CapabilityEvidenceNone,
				Classification: EvidenceProfileRequirement,
				SourcePath:     relative,
				SourceDigest:   contentIdentity(data),
				Detail:         "expected content not found",
			})
			continue
		}
		return append(evidence, CapabilityEvidence{
			Status:         CapabilityEvidencePresent,
			Kind:           capability.EvidenceKind,
			Strength:       CapabilityEvidenceDeclared,
			Classification: EvidenceProfileRequirement,
			SourcePath:     relative,
			SourceDigest:   contentIdentity(data),
		})
	}
	return evidence
}

func collectInstalledSkillEvidence(root *os.Root, capability RepositoryCapability) CapabilityEvidence {
	skill, ok := capability.Probe["skill"].(string)
	if !ok || !identifierIsSafe(skill) {
		return invalidCapabilityEvidence(capability.EvidenceKind, "installed skill probe is invalid")
	}
	relative := ".agents/skills/" + skill + "/SKILL.md"
	data, state, detail := readBoundedRegularFile(root, relative, maxCapabilityFileBytes)
	if state == CapabilityEvidenceInvalid {
		return invalidCapabilityEvidence(capability.EvidenceKind, detail)
	}
	if state != CapabilityEvidencePresent {
		return absentCapabilityEvidence(capability.EvidenceKind)
	}
	return CapabilityEvidence{
		Status:         CapabilityEvidencePresent,
		Kind:           capability.EvidenceKind,
		Strength:       CapabilityEvidenceVerified,
		Classification: EvidenceProfileRequirement,
		SourcePath:     relative,
		SourceDigest:   contentIdentity(data),
	}
}

func collectExecutableEvidence(capability RepositoryCapability) CapabilityEvidence {
	executable, ok := capability.Probe["executable"].(string)
	if !ok || !identifierIsSafe(executable) {
		return invalidCapabilityEvidence(capability.EvidenceKind, "executable probe is invalid")
	}
	result := resolveExecutableCandidate(executable)
	if result.Reason == executableProbeReasonNotFound {
		evidence := absentCapabilityEvidence(capability.EvidenceKind)
		evidence.Detail = result.Reason
		return evidence
	}
	if result.Reason != "" {
		return CapabilityEvidence{
			Status:         CapabilityEvidenceInvalid,
			Kind:           capability.EvidenceKind,
			Strength:       CapabilityEvidenceNone,
			Classification: EvidenceProfileRequirement,
			SourcePath:     filepath.ToSlash(result.Candidate),
			Detail:         result.Reason,
		}
	}
	return CapabilityEvidence{
		Status:         CapabilityEvidencePresent,
		Kind:           capability.EvidenceKind,
		Strength:       CapabilityEvidenceDiscovered,
		Classification: EvidenceProfileRequirement,
		SourcePath:     filepath.ToSlash(result.Candidate),
	}
}

// resolveExecutableCandidate inspects PATH candidates without executing them.
func resolveExecutableCandidate(name string) executableProbeResult {
	var firstFailure executableProbeResult
	for _, directory := range filepath.SplitList(os.Getenv("PATH")) {
		if directory == "" {
			directory = "."
		}
		result, exists := inspectExecutableCandidate(filepath.Join(directory, name))
		if !exists {
			continue
		}
		if result.Reason == "" {
			return result
		}
		if firstFailure.Candidate == "" {
			firstFailure = result
		}
	}
	if firstFailure.Candidate != "" {
		return firstFailure
	}
	return executableProbeResult{Reason: executableProbeReasonNotFound}
}

func inspectExecutableCandidate(candidate string) (executableProbeResult, bool) {
	absolute, err := filepath.Abs(candidate)
	if err == nil {
		candidate = absolute
	}
	candidate = filepath.Clean(candidate)
	result := executableProbeResult{Candidate: candidate}
	info, err := os.Lstat(candidate)
	if errors.Is(err, fs.ErrNotExist) {
		return executableProbeResult{}, false
	}
	if err != nil {
		result.Reason = executableProbeReasonNotExecutable
		return result, true
	}

	current := candidate
	visited := make(map[string]struct{}, maxExecutableLinkHops)
	for {
		if _, exists := visited[current]; exists {
			result.Reason = executableProbeReasonLinkCycle
			return result, true
		}
		visited[current] = struct{}{}

		if info.Mode()&fs.ModeSymlink == 0 {
			if !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
				result.Reason = executableProbeReasonNotExecutable
				return result, true
			}
			result.Resolved = current
			return result, true
		}
		if result.HopCount >= maxExecutableLinkHops {
			result.Reason = executableProbeReasonLinkCycle
			return result, true
		}

		target, err := os.Readlink(current)
		if err != nil {
			result.Reason = executableProbeReasonBrokenLink
			return result, true
		}
		result.HopCount++
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(current), target)
		}
		current = filepath.Clean(target)
		info, err = os.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			result.Reason = executableProbeReasonBrokenLink
			return result, true
		}
		if err != nil {
			result.Reason = executableProbeReasonNotExecutable
			return result, true
		}
	}
}

func evaluateCapability(capability RepositoryCapability, evidence []CapabilityEvidence) CapabilityOutcome {
	status := CapabilityMissing
	for _, item := range evidence {
		if item.Status == CapabilityEvidencePresent &&
			evidenceRank[item.Strength] >= evidenceRank[capability.MinimumEvidence] {
			status = CapabilitySatisfied
			break
		}
		if item.Status == CapabilityEvidencePresent || item.Status == CapabilityEvidenceInvalid {
			status = CapabilityInsufficient
		}
	}
	blocking := capability.Requirement == CapabilityRequired && status != CapabilitySatisfied
	diagnostic := CapabilityDiagnostic{
		Code:       capabilityDiagnosticCode(capability.Requirement, status),
		Message:    fmt.Sprintf("%s has no compatible local evidence.", capability.Title),
		NextAction: capability.NextAction,
	}
	if status == CapabilitySatisfied {
		diagnostic = CapabilityDiagnostic{
			Code:    "capability.satisfied",
			Message: fmt.Sprintf("%s has sufficient local evidence.", capability.Title),
		}
	} else if status == CapabilityInsufficient {
		diagnostic.Code = "capability.evidence.insufficient"
		diagnostic.Message = fmt.Sprintf("%s local evidence is present but insufficient.", capability.Title)
	}
	return CapabilityOutcome{
		ID:              capability.ID,
		Title:           capability.Title,
		Requirement:     capability.Requirement,
		EvidenceKind:    capability.EvidenceKind,
		MinimumEvidence: capability.MinimumEvidence,
		Status:          status,
		Blocking:        blocking,
		Evidence:        evidence,
		Diagnostic:      diagnostic,
		Explanation:     capability.Explanation,
	}
}

func capabilityDiagnosticCode(requirement CapabilityRequirement, status CapabilityStatus) string {
	if status == CapabilityInsufficient {
		return "capability.evidence.insufficient"
	}
	switch requirement {
	case CapabilityRequired:
		return "capability.required.missing"
	case CapabilityRecommended:
		return "capability.recommended.missing"
	default:
		return "capability.optional.missing"
	}
}

func evaluatePostgreSQLCapability(
	root *os.Root,
	capability RepositoryCapability,
) (CapabilityOutcome, PostgreSQLEvidence, error) {
	accepted := stringsFromAny(capability.Probe["paths"])
	slices.Sort(accepted)
	postgres := PostgreSQLEvidence{AcceptedContractPaths: accepted}
	contains, _ := capability.Probe["contains"].(string)
	contractEvidence := make([]CapabilityEvidence, 0, len(accepted))
	for _, relative := range accepted {
		data, state, detail := readBoundedRegularFile(root, relative, maxCapabilityFileBytes)
		if state == CapabilityEvidenceInvalid {
			return CapabilityOutcome{}, PostgreSQLEvidence{}, fmt.Errorf("inspect PostgreSQL contract %q: %s", relative, detail)
		}
		if state != CapabilityEvidencePresent {
			contractEvidence = append(contractEvidence, CapabilityEvidence{
				Status:         CapabilityEvidenceAbsent,
				Kind:           CapabilityEvidenceDeclaredFile,
				Strength:       CapabilityEvidenceNone,
				Classification: EvidenceRepositoryContract,
				SourcePath:     relative,
				Detail:         "file not found",
			})
			continue
		}
		if !bytes.Contains(data, []byte(contains)) {
			contractEvidence = append(contractEvidence, CapabilityEvidence{
				Status:         CapabilityEvidenceAbsent,
				Kind:           CapabilityEvidenceDeclaredFile,
				Strength:       CapabilityEvidenceNone,
				Classification: EvidenceRepositoryContract,
				SourcePath:     relative,
				SourceDigest:   contentIdentity(data),
				Detail:         "expected content not found",
			})
			continue
		}
		evidence := CapabilityEvidence{
			Status:         CapabilityEvidencePresent,
			Kind:           CapabilityEvidenceDeclaredFile,
			Strength:       CapabilityEvidenceDeclared,
			Classification: EvidenceRepositoryContract,
			SourcePath:     relative,
			SourceDigest:   contentIdentity(data),
			Detail:         "accepted PostgreSQL repository contract",
		}
		postgres.Contract = &evidence
		contractEvidence = append(contractEvidence, evidence)
		break
	}
	implementation, err := collectPostgreSQLImplementationEvidence(root)
	if err != nil {
		return CapabilityOutcome{}, PostgreSQLEvidence{}, err
	}
	postgres.Implementation = implementation
	if postgres.Contract != nil {
		return evaluateCapability(capability, contractEvidence), postgres, nil
	}

	evidence := append([]CapabilityEvidence(nil), implementation...)
	evidence = append(evidence, contractEvidence...)
	outcome := evaluateCapability(capability, evidence)
	outcome.Status = CapabilityMissing
	outcome.Blocking = capability.Requirement == CapabilityRequired
	outcome.Diagnostic.Code = "capability.contract.missing"
	if len(implementation) == 0 {
		outcome.Diagnostic.Message = "PostgreSQL has no accepted repository contract and no implementation evidence was found."
	} else {
		outcome.Diagnostic.Message = "PostgreSQL implementation evidence was found, but the required repository contract is absent."
	}
	outcome.Diagnostic.NextAction = "Record the PostgreSQL repository contract at one accepted path: " + strings.Join(accepted, ", ")
	return outcome, postgres, nil
}

func collectPostgreSQLImplementationEvidence(root *os.Root) ([]CapabilityEvidence, error) {
	candidates := []string{
		"package.json",
		"packages/backend/package.json",
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
		".env.example",
	}
	evidence := make([]CapabilityEvidence, 0)
	for _, relative := range candidates {
		data, state, detail := readBoundedRegularFile(root, relative, maxCapabilityFileBytes)
		if state == CapabilityEvidenceInvalid {
			return nil, fmt.Errorf("inspect PostgreSQL implementation evidence %q: %s", relative, detail)
		}
		if state != CapabilityEvidencePresent {
			continue
		}
		lower := strings.ToLower(string(data))
		facts := make([]string, 0, 3)
		for _, fact := range []string{"database_url", "drizzle-orm", "postgres"} {
			if strings.Contains(lower, fact) {
				facts = append(facts, fact)
			}
		}
		if len(facts) == 0 {
			continue
		}
		evidence = append(evidence, CapabilityEvidence{
			Status:         CapabilityEvidencePresent,
			Kind:           CapabilityEvidenceDeclaredFile,
			Strength:       CapabilityEvidenceDiscovered,
			Classification: EvidenceImplementation,
			SourcePath:     relative,
			SourceDigest:   contentIdentity(data),
			Detail:         "observed " + strings.Join(facts, ", "),
		})
	}
	return evidence, nil
}

func collectHTTPRouteCandidates(
	ctx context.Context,
	root *os.Root,
	profile ResolvedProfile,
	catalog *Catalog,
) ([]HTTPRouteCandidate, error) {
	if !profileHasHTTPContract(profile, catalog) {
		return nil, nil
	}
	info, err := root.Lstat(filepath.FromSlash("packages/backend"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("inspect HTTP candidate root: %w", err)
	}
	if !info.IsDir() {
		return nil, nil
	}

	candidates := make([]HTTPRouteCandidate, 0)
	sourceFiles := 0
	totalBytes := 0
	err = fs.WalkDir(root.FS(), "packages/backend", func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "node_modules", ".git", "dist", "build", "coverage":
				return fs.SkipDir
			default:
				return nil
			}
		}
		if entry.Type()&fs.ModeSymlink != 0 || !httpSourceExtension(filepath.Ext(filePath)) {
			return nil
		}
		data, state, detail := readBoundedRegularFile(root, filePath, maxCapabilityFileBytes)
		if state == CapabilityEvidenceInvalid {
			return fmt.Errorf("inspect HTTP candidate %q: %s", filePath, detail)
		}
		if state != CapabilityEvidencePresent {
			return nil
		}
		fileCandidates := routeCandidatesFromSource(filePath, data)
		if len(fileCandidates) == 0 {
			return nil
		}
		sourceFiles++
		if sourceFiles > maxHTTPSourceFiles {
			return fmt.Errorf("HTTP candidate source count exceeds %d", maxHTTPSourceFiles)
		}
		totalBytes += len(data)
		if totalBytes > maxHTTPSourceBytes {
			return fmt.Errorf("HTTP candidate source bytes exceed %d", maxHTTPSourceBytes)
		}
		candidates = append(candidates, fileCandidates...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect bounded HTTP route candidates: %w", err)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].SourcePath != candidates[j].SourcePath {
			return candidates[i].SourcePath < candidates[j].SourcePath
		}
		if candidates[i].Scope != candidates[j].Scope {
			return candidates[i].Scope < candidates[j].Scope
		}
		return strings.Join(candidates[i].Methods, ",") < strings.Join(candidates[j].Methods, ",")
	})
	return compactHTTPRouteCandidates(candidates), nil
}

func profileHasHTTPContract(profile ResolvedProfile, catalog *Catalog) bool {
	if _, selected := stringSet(profile.Decisions)["http.contract"]; !selected {
		return false
	}
	if profile.Source == ProfileSourceBuiltIn {
		_, ok := objectValue(catalog.profiles[profile.ID]["httpContract"])
		return ok
	}
	return true
}

func httpSourceExtension(extension string) bool {
	switch strings.ToLower(extension) {
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs":
		return true
	default:
		return false
	}
}

func routeCandidatesFromSource(sourcePath string, data []byte) []HTTPRouteCandidate {
	digest := contentIdentity(data)
	candidates := make([]HTTPRouteCandidate, 0)
	for _, match := range httpMethodCall.FindAllSubmatch(data, -1) {
		candidates = append(candidates, HTTPRouteCandidate{
			SourcePath:   sourcePath,
			SourceDigest: digest,
			Scope:        string(match[2]),
			Methods:      []string{strings.ToUpper(string(match[1]))},
		})
	}
	for _, match := range httpScopeCall.FindAllSubmatch(data, -1) {
		candidates = append(candidates, HTTPRouteCandidate{
			SourcePath:   sourcePath,
			SourceDigest: digest,
			Scope:        string(match[1]),
			Methods:      nil,
		})
	}
	return candidates
}

func compactHTTPRouteCandidates(input []HTTPRouteCandidate) []HTTPRouteCandidate {
	if len(input) < 2 {
		return input
	}
	output := input[:1]
	for _, candidate := range input[1:] {
		previous := output[len(output)-1]
		if previous.SourcePath == candidate.SourcePath &&
			previous.SourceDigest == candidate.SourceDigest &&
			previous.Scope == candidate.Scope &&
			slices.Equal(previous.Methods, candidate.Methods) {
			continue
		}
		output = append(output, candidate)
	}
	return output
}

func resolveVerificationProjection(
	root *os.Root,
	profile ResolvedProfile,
	decisions []DecisionValue,
	roleMappings map[string]string,
	catalog *Catalog,
) ([]VerificationProjection, []ProfileDivergence, error) {
	projections := make([]VerificationProjection, 0)
	divergences := make([]ProfileDivergence, 0)
	if profile.Source == ProfileSourceBuiltIn {
		for _, raw := range objectsOrEmpty(catalog.profiles[profile.ID]["verification"]) {
			id, _ := stringValue(raw, "id")
			role, _ := stringValue(raw, "kind")
			tool, _ := stringValue(raw, "tool")
			command, _ := stringValue(raw, "command")
			resolvedCommand := command
			mappedCommand, mapped := roleMappings[role]
			var declaration commandDeclaration
			var err error
			if mapped {
				resolvedCommand = strings.TrimSpace(mappedCommand)
				declaration, err = validateLocalCommandDeclaration(root, resolvedCommand)
			} else {
				resolvedCommand, declaration, err = resolveProfileVerificationCommand(
					root,
					role,
					command,
				)
			}
			if err != nil {
				return nil, nil, fmt.Errorf("validate profile Verification expectation %q: %w", id, err)
			}
			classification := VerificationProfileExpectation
			if mapped || resolvedCommand != command {
				classification = VerificationRepositoryCommand
			}
			satisfiedByCommand := ""
			if mapped && declaration.Path != "" {
				satisfiedByCommand = resolvedCommand
			}
			projections = append(projections, VerificationProjection{
				ID:                   id,
				Role:                 role,
				Tool:                 tool,
				Command:              resolvedCommand,
				SatisfiedByCommand:   satisfiedByCommand,
				Classification:       classification,
				RepositoryExecutable: declaration.Path != "",
				DeclarationPath:      declaration.Path,
				DeclarationDigest:    declaration.Digest,
			})
			if declaration.Path == "" {
				if mapped {
					divergences = append(divergences, ProfileDivergence{
						Code:        "verification.role-mapping.undeclared",
						ID:          id,
						Requirement: CapabilityRecommended,
						Blocking:    false,
						Message:     fmt.Sprintf("portable %s role maps to repository command %q, but no matching local declaration exists", role, resolvedCommand),
						NextAction:  "declare the mapped repository command or map the portable role to another declared command",
					})
					continue
				}
				divergences = append(divergences, ProfileDivergence{
					Code:        "verification.profile-expectation.unresolved",
					ID:          id,
					Requirement: CapabilityRecommended,
					Blocking:    false,
					Message:     fmt.Sprintf("portable %s role expects %q, but no matching local declaration exists", role, command),
					NextAction:  "map the portable role to a declared repository command or keep it as a profile expectation",
				})
			}
		}
	}
	for _, decision := range decisions {
		if decision.ID != "verification.gate" {
			continue
		}
		command, _ := decision.Value.(string)
		declaration, err := validateLocalCommandDeclaration(root, command)
		if err != nil {
			return nil, nil, fmt.Errorf("validate selected repository Verification command: %w", err)
		}
		projections = append(projections, VerificationProjection{
			ID:                   "verification.gate",
			Role:                 "repository-gate",
			Command:              command,
			Classification:       VerificationRepositoryCommand,
			RepositoryExecutable: declaration.Path != "",
			DeclarationPath:      declaration.Path,
			DeclarationDigest:    declaration.Digest,
		})
		if declaration.Path == "" {
			divergences = append(divergences, ProfileDivergence{
				Code:        "verification.command.undeclared",
				ID:          "verification.gate",
				Requirement: CapabilityRequired,
				Blocking:    true,
				Message:     fmt.Sprintf("selected repository Verification command %q has no matching local declaration", command),
				NextAction:  "select a command declared by the repository or add its local declaration",
			})
		}
	}
	sort.Slice(projections, func(i, j int) bool { return projections[i].ID < projections[j].ID })
	return projections, divergences, nil
}

func resolveProfileVerificationCommand(
	root *os.Root,
	role string,
	command string,
) (string, commandDeclaration, error) {
	declaration, err := validateLocalCommandDeclaration(root, command)
	if err != nil || declaration.Path != "" || role != "format" {
		return command, declaration, err
	}
	fields := strings.Fields(command)
	hasRTK := len(fields) != 0 && fields[0] == "rtk"
	if hasRTK {
		fields = fields[1:]
	}
	if len(fields) != 3 || fields[1] != "run" {
		return command, declaration, nil
	}
	switch fields[0] {
	case "bun", "npm", "pnpm", "yarn":
	default:
		return command, declaration, nil
	}
	for _, script := range []string{"fmt", "format"} {
		if script == fields[2] {
			continue
		}
		candidateFields := []string{fields[0], "run", script}
		if hasRTK {
			candidateFields = append([]string{"rtk"}, candidateFields...)
		}
		candidate := strings.Join(candidateFields, " ")
		candidateDeclaration, candidateErr := validateLocalCommandDeclaration(root, candidate)
		if candidateErr != nil {
			return "", commandDeclaration{}, candidateErr
		}
		if candidateDeclaration.Path != "" {
			return candidate, candidateDeclaration, nil
		}
	}
	return command, declaration, nil
}

type commandDeclaration struct {
	Path   string
	Digest string
}

func validateLocalCommandDeclaration(root *os.Root, command string) (commandDeclaration, error) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return commandDeclaration{}, nil
	}
	if fields[0] == "rtk" {
		fields = fields[1:]
	}
	if len(fields) == 0 {
		return commandDeclaration{}, nil
	}
	if fields[0] == "make" && len(fields) == 2 && makeTargetName.MatchString(fields[1]) {
		return validateMakeTarget(root, fields[1])
	}
	if len(fields) == 3 && fields[1] == "run" {
		switch fields[0] {
		case "bun", "npm", "pnpm", "yarn":
			return validatePackageScript(root, fields[2])
		}
	}
	return commandDeclaration{}, nil
}

func validatePackageScript(root *os.Root, script string) (commandDeclaration, error) {
	data, state, detail := readBoundedRegularFile(root, "package.json", maxCapabilityFileBytes)
	if state == CapabilityEvidenceInvalid {
		return commandDeclaration{}, errors.New(detail)
	}
	if state != CapabilityEvidencePresent {
		return commandDeclaration{}, nil
	}
	var declaration struct {
		Scripts map[string]json.RawMessage `json:"scripts"`
	}
	if err := json.Unmarshal(data, &declaration); err != nil {
		return commandDeclaration{}, fmt.Errorf("parse package.json scripts: %w", err)
	}
	value, ok := declaration.Scripts[script]
	if !ok {
		return commandDeclaration{}, nil
	}
	var command string
	if err := json.Unmarshal(value, &command); err != nil || strings.TrimSpace(command) == "" {
		return commandDeclaration{}, nil
	}
	return commandDeclaration{Path: "package.json", Digest: contentIdentity(data)}, nil
}

func validateMakeTarget(root *os.Root, target string) (commandDeclaration, error) {
	for _, relative := range []string{"Makefile", "makefile", "GNUmakefile"} {
		data, state, detail := readBoundedRegularFile(root, relative, maxCapabilityFileBytes)
		if state == CapabilityEvidenceInvalid {
			return commandDeclaration{}, errors.New(detail)
		}
		if state != CapabilityEvidencePresent {
			continue
		}
		pattern := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(target) + `\s*:(?:[^=]|$)`)
		if pattern.Match(data) {
			return commandDeclaration{Path: relative, Digest: contentIdentity(data)}, nil
		}
	}
	return commandDeclaration{}, nil
}

func readBoundedRegularFile(
	root *os.Root,
	relative string,
	limit int,
) ([]byte, CapabilityEvidenceStatus, string) {
	if !repositoryPathIsSafe(relative) {
		return nil, CapabilityEvidenceInvalid, "path is not normalized and repository-relative"
	}
	info, err := root.Lstat(filepath.FromSlash(relative))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, CapabilityEvidenceAbsent, ""
	}
	if err != nil {
		return nil, CapabilityEvidenceInvalid, "declared file cannot be inspected"
	}
	if !info.Mode().IsRegular() {
		return nil, CapabilityEvidenceInvalid, "declared file must be a regular file"
	}
	if info.Size() > int64(limit) {
		return nil, CapabilityEvidenceInvalid, fmt.Sprintf("declared file exceeds %d bytes", limit)
	}
	file, err := root.Open(filepath.FromSlash(relative))
	if err != nil {
		return nil, CapabilityEvidenceInvalid, "declared file cannot be opened"
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return nil, CapabilityEvidenceInvalid, "declared file changed during inspection"
	}
	data, err := io.ReadAll(io.LimitReader(file, int64(limit)+1))
	if err != nil {
		return nil, CapabilityEvidenceInvalid, "declared file cannot be read"
	}
	if len(data) > limit {
		return nil, CapabilityEvidenceInvalid, fmt.Sprintf("declared file exceeds %d bytes", limit)
	}
	return data, CapabilityEvidencePresent, ""
}

func contentIdentity(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func absentCapabilityEvidence(kind CapabilityEvidenceKind) CapabilityEvidence {
	return CapabilityEvidence{
		Status:         CapabilityEvidenceAbsent,
		Kind:           kind,
		Strength:       CapabilityEvidenceNone,
		Classification: EvidenceProfileRequirement,
	}
}

func invalidCapabilityEvidence(kind CapabilityEvidenceKind, detail string) CapabilityEvidence {
	return CapabilityEvidence{
		Status:         CapabilityEvidenceInvalid,
		Kind:           kind,
		Strength:       CapabilityEvidenceNone,
		Classification: EvidenceProfileRequirement,
		Detail:         detail,
	}
}

func stringsFromAny(value any) []string {
	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil
			}
			values = append(values, text)
		}
		return values
	default:
		return nil
	}
}

func identifierIsSafe(value string) bool {
	if value == "" {
		return false
	}
	for index, character := range value {
		alphaNumeric := character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9'
		if !alphaNumeric && character != '.' && character != '_' && character != '-' {
			return false
		}
		if (character == '.' || character == '_' || character == '-') &&
			(index == 0 || index == len(value)-1) {
			return false
		}
	}
	return true
}

// RenderProfileDivergences renders the projection carried by each divergence.
// It does not inspect the repository or repeat capability evaluation.
func RenderProfileDivergences(divergences []ProfileDivergence) string {
	if len(divergences) == 0 {
		return "Divergences: none\n"
	}
	var rendered strings.Builder
	groups := []struct {
		group   ProfileDivergenceGroup
		heading string
	}{
		{group: ProfileDivergenceBlocking, heading: "Blocking divergences:"},
		{group: ProfileDivergenceAdvisory, heading: "Advisory divergences:"},
		{group: ProfileDivergenceInformational, heading: "Informational divergences:"},
	}
	for _, section := range groups {
		wroteHeading := false
		for _, divergence := range divergences {
			group := divergence.Group
			if group == "" {
				group = profileDivergenceGroup(divergence.Requirement)
			}
			if group != section.group {
				continue
			}
			if !wroteHeading {
				if rendered.Len() != 0 {
					rendered.WriteByte('\n')
				}
				fmt.Fprintln(&rendered, section.heading)
				wroteHeading = true
			}
			fmt.Fprintf(&rendered, "- %s (%s): %s\n", divergence.ID, divergence.Code, divergence.Message)
			renderProfileDivergenceProbe(&rendered, divergence)
			renderProfileCapabilityResolution(&rendered, divergence.CapabilityResolution)
			if group == ProfileDivergenceAdvisory {
				statement := divergence.NonBlockingStatement
				if statement == "" {
					statement = advisoryNonBlockingStatement
				}
				fmt.Fprintf(&rendered, "  %s\n", statement)
			}
			if divergence.NextAction != "" {
				fmt.Fprintf(&rendered, "  Next action: %s\n", divergence.NextAction)
			}
		}
	}
	return rendered.String()
}

func renderProfileDivergenceProbe(rendered *strings.Builder, divergence ProfileDivergence) {
	kind, _ := divergence.Probe["kind"].(string)
	switch CapabilityEvidenceKind(kind) {
	case CapabilityEvidenceDeclaredFile:
		expected, _ := divergence.Probe["contains"].(string)
		inspected := 0
		for _, evidence := range divergence.Evidence {
			if evidence.Kind != CapabilityEvidenceDeclaredFile ||
				evidence.Classification == EvidenceImplementation ||
				evidence.SourcePath == "" {
				continue
			}
			state := string(evidence.Status)
			if evidence.Status == CapabilityEvidenceAbsent && evidence.SourceDigest != "" {
				state = "present"
			}
			if evidence.Detail != "" {
				state += " (" + evidence.Detail + ")"
			}
			fmt.Fprintf(
				rendered,
				"  Inspected path: %s: %s; expected content: %s\n",
				evidence.SourcePath,
				state,
				expected,
			)
			inspected++
		}
		if inspected == 0 {
			fmt.Fprintln(rendered, "  Inspected paths: none recorded")
		}
	case CapabilityEvidenceExecutable:
		executable, _ := divergence.Probe["executable"].(string)
		if executable != "" {
			fmt.Fprintf(rendered, "  Executable probe: %s\n", executable)
		}
		for _, evidence := range divergence.Evidence {
			if evidence.Kind != CapabilityEvidenceExecutable || evidence.SourcePath == "" {
				continue
			}
			detail := ""
			if evidence.Detail != "" {
				detail = " (" + evidence.Detail + ")"
			}
			fmt.Fprintf(rendered, "  Inspected candidate: %s%s\n", evidence.SourcePath, detail)
			return
		}
		fmt.Fprintln(rendered, "  Inspected candidate: none existed")
	case CapabilityEvidenceInstalledSkill:
		skill, _ := divergence.Probe["skill"].(string)
		fmt.Fprintf(rendered, "  Installed-skill probe: %s\n", skill)
		if len(divergence.Evidence) != 0 {
			fmt.Fprintf(rendered, "  Evidence state: %s\n", divergence.Evidence[0].Status)
		}
	case "":
		return
	default:
		fmt.Fprintf(rendered, "  Probe kind: %s\n", kind)
	}
}

func renderProfileCapabilityResolution(
	rendered *strings.Builder,
	resolution *ProfileCapabilityResolution,
) {
	if resolution == nil {
		return
	}
	fmt.Fprintf(rendered, "  Selected technology: %s\n", resolution.SelectedTechnology)
	fmt.Fprintf(rendered, "  Repository remediation: %s\n", resolution.RepositoryRemediation)
	fmt.Fprintf(rendered, "  Profile adaptation: %s\n", resolution.ProfileAdaptation)
	if len(resolution.RemovedDecisions) == 0 {
		fmt.Fprintln(rendered, "  Decision cascade: none")
		return
	}
	fmt.Fprintf(rendered, "  Decision cascade: %s\n", strings.Join(resolution.RemovedDecisions, ", "))
}

func sortProfileDivergences(divergences []ProfileDivergence) {
	sort.Slice(divergences, func(i, j int) bool {
		leftGroup := profileDivergenceGroup(divergences[i].Requirement)
		rightGroup := profileDivergenceGroup(divergences[j].Requirement)
		if leftGroup != rightGroup {
			return profileDivergenceGroupRank(leftGroup) < profileDivergenceGroupRank(rightGroup)
		}
		if divergences[i].ID != divergences[j].ID {
			return divergences[i].ID < divergences[j].ID
		}
		return divergences[i].Code < divergences[j].Code
	})
}

func profileDivergenceGroupRank(group ProfileDivergenceGroup) int {
	switch group {
	case ProfileDivergenceBlocking:
		return 0
	case ProfileDivergenceAdvisory:
		return 1
	default:
		return 2
	}
}
