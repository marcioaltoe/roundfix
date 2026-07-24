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
)

// PlanRequest is the complete normalized non-interactive planning input.
type PlanRequest struct {
	Repository   string
	ProfileID    string
	Decisions    []DecisionValue
	Preservation RootPreservationRequest
}

// PlanOutcome contains either one complete portable plan or an actionable
// result. A nil Plan always means that no partial plan is available.
type PlanOutcome struct {
	Plan   *PlanDocument
	Result Result
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
	Preimages      []Preimage          `json:"preimages"`
	Postimages     []Postimage         `json:"postimages"`
	Warnings       []Finding           `json:"warnings"`
	SetupManifest  SetupManifest       `json:"setupManifest"`
	ManagedEntries []ManagedEntry      `json:"managedEntries"`
	FileChanges    []FileChange        `json:"fileChanges"`
	PlanDigest     string              `json:"planDigest"`
}

// Result is the strict automation result used when no complete plan can be
// emitted and by later Baseline operations.
type Result struct {
	SchemaVersion      string      `json:"schemaVersion"`
	Operation          string      `json:"operation"`
	State              string      `json:"state"`
	Category           string      `json:"category,omitempty"`
	Message            string      `json:"message,omitempty"`
	NextAction         string      `json:"nextAction,omitempty"`
	PlanDigest         string      `json:"planDigest,omitempty"`
	VerifiedPostimages []Postimage `json:"verifiedPostimages"`
	Warnings           []Finding   `json:"warnings"`
	Recommendations    []string    `json:"recommendations"`
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
	if profileID == "" {
		return actionOutcome("decision", "Baseline Profile selection is required",
			"rerun with --profile <id>", initial.Snapshot.Warnings), nil
	}
	profile, err := ResolveProfile(initial.Root, profileID, catalog)
	if err != nil {
		return PlanOutcome{}, fmt.Errorf("resolve Baseline Profile: %w", err)
	}
	profile.Path = portableProfilePath(initial.Root, profile.Path)

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

	preservation, err := PlanRootPreservation(initial, request.Preservation)
	if err != nil {
		return PlanOutcome{}, err
	}
	if preservation.State != PreservationStateReady {
		return actionOutcome("decision", preservation.NextAction, preservation.NextAction,
			append(initial.Snapshot.Warnings, preservation.Findings...)), nil
	}
	alignment, err := ResolveProfileAlignment(ctx, initial.Root, ProfileAlignmentRequest{
		ProfileID: profile.ID,
		Decisions: profileAlignmentDecisions(profile, decisions),
	}, catalog)
	if err != nil {
		return PlanOutcome{}, err
	}
	if !alignment.Ready {
		next := make([]string, 0)
		for _, divergence := range alignment.Divergences {
			if divergence.Blocking {
				next = append(next, divergence.ID)
			}
		}
		return actionOutcome("decision",
			"required profile alignment is unresolved: "+strings.Join(next, ", "),
			"provide the missing repository evidence or select a different profile",
			initial.Snapshot.Warnings), nil
	}

	activeModules, artifacts, err := resolveManagedArtifacts(catalog, profile, decisions)
	if err != nil {
		return PlanOutcome{}, err
	}
	manifest := buildSetupManifest(catalog, profile, decisions, activeModules, artifacts, alignment.Verification)
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
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
	mutablePaths = append(mutablePaths, alignmentEvidencePaths(alignment)...)
	if len(preservation.RepositoryRulesBytes) != 0 {
		mutablePaths = append(mutablePaths, repositoryRulesPath)
	}
	snapshot, err := inspectRepositorySnapshot(initial.Root, InventoryRequest{MutablePaths: mutablePaths})
	if err != nil {
		return PlanOutcome{}, err
	}
	inspection := RepositoryInspection{Root: initial.Root, Identity: initial.Identity, Snapshot: snapshot}
	preservation, err = PlanRootPreservation(inspection, request.Preservation)
	if err != nil {
		return PlanOutcome{}, err
	}
	if preservation.State != PreservationStateReady {
		return actionOutcome("decision", preservation.NextAction, preservation.NextAction,
			append(snapshot.Warnings, preservation.Findings...)), nil
	}
	retention, retentionAction, err := resolvePlanRetention(
		initial.Root,
		snapshot,
		catalog,
		profile.ID,
		preservation,
	)
	if err != nil {
		return PlanOutcome{}, err
	}
	if retentionAction != "" {
		return actionOutcome(
			"classification",
			retentionAction,
			"restore an exact supported Setup Manifest or add a reviewed Upgrade Retention Contract",
			snapshot.Warnings,
		), nil
	}

	postimages, ledger, err := assemblePostimages(initial.Root, snapshot, artifacts, manifestBytes, preservation)
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
	if err := ValidatePlanDocument(doc); err != nil {
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
		if item.Command != "" {
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
	profileID string,
	preservation RootPreservationPlan,
) ([]RetentionEvidence, string, error) {
	retention := readoptionRetentionEvidence(preservation.Dispositions)
	preimage, ok := preimagesByPath(snapshot.Preimages)[manifestPath]
	if !ok || !preimage.Exists {
		return retention, "", nil
	}
	if preimage.Kind != PreimageRegular {
		return nil, "", fmt.Errorf("resolve Upgrade Retention: Setup Manifest is not a regular file")
	}
	data, err := readOptionalRegular(root, manifestPath)
	if err != nil {
		return nil, "", fmt.Errorf("resolve Upgrade Retention: %w", err)
	}
	if planContentIdentity(data) != preimage.ContentIdentity {
		return nil, "", errors.New("resolve Upgrade Retention: Setup Manifest changed during planning")
	}
	manifest, diagnostics := decodeDocument(data, manifestPath)
	if len(diagnostics) != 0 || manifest == nil {
		return retention, "the existing Setup Manifest is not valid strict JSON", nil
	}

	generator, _ := objectValue(manifest["generator"])
	declaredBaseline, _ := stringValue(generator, "baseline")
	currentBaseline := "baseline." + profileID + "-" + ManifestVersion
	if manifest["schemaVersion"] == ManifestSchema &&
		manifest["version"] == ManifestVersion &&
		declaredBaseline == currentBaseline {
		return retention, "", nil
	}
	if currentSetupManifestProfileIsValid(root, manifest, generator, catalog) {
		return retention, "", nil
	}
	// The portable-v3 manifest remains a supported reader contract while the
	// current writer emits the profile-bound 0.0.1 identity.
	if declaredBaseline == "baseline.portable-v3" {
		return retention, "", nil
	}

	var matches []UpgradeRetentionContract
	fingerprint := legacyManifestFingerprint(manifest)
	for _, transitionID := range catalog.TransitionIDs() {
		contract, err := catalog.UpgradeRetentionContract(transitionID)
		if err != nil {
			return nil, "", err
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
		return retention,
			fmt.Sprintf("the existing Setup Manifest identity %q has no unique maintained transition", identity),
			nil
	}
	return append(retention, transitionRetentionEvidence(matches[0])...), "", nil
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

func resolveManagedArtifacts(
	catalog *Catalog,
	profile ResolvedProfile,
	decisions []DecisionValue,
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
				renderValues[token] = "`" + renderDecisionValue(value) + "`"
			}
		}
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
		valuesForArtifact := artifactRenderValues(catalog, declaration, activeModules, paths)
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
	artifactPaths map[string]string,
) map[string]string {
	values := make(map[string]string)
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

func assemblePostimages(
	root string,
	snapshot RepositorySnapshot,
	artifacts []plannedArtifact,
	manifestBytes []byte,
	preservation RootPreservationPlan,
) ([]Postimage, []ManagedEntry, error) {
	byPath := make(map[string][]plannedArtifact)
	for _, artifact := range artifacts {
		byPath[artifact.Path] = append(byPath[artifact.Path], artifact)
	}
	outputs := make(map[string][]byte)
	for relative, grouped := range byPath {
		current, err := readOptionalRegular(root, relative)
		if err != nil {
			return nil, nil, err
		}
		content := string(current)
		for _, artifact := range grouped {
			content = upsertManagedBlock(content, artifact)
		}
		outputs[relative] = []byte(content)
	}
	outputs[manifestPath] = append([]byte(nil), manifestBytes...)
	for _, backup := range preservation.Backups {
		source, err := readOptionalRegular(root, backup.SourcePath)
		if err != nil {
			return nil, nil, err
		}
		outputs[backup.Path] = source
	}
	if len(preservation.RepositoryRulesBytes) != 0 {
		outputs[repositoryRulesPath] = append([]byte(nil), preservation.RepositoryRulesBytes...)
	}

	preimages := preimagesByPath(snapshot.Preimages)
	paths := sortedKeys(outputs)
	postimages := make([]Postimage, 0, len(paths))
	ledger := make([]ManagedEntry, 0)
	for _, relative := range paths {
		content := outputs[relative]
		identity := planContentIdentity(content)
		postimages = append(postimages, Postimage{
			Path: relative, Kind: PreimageRegular, Mode: 0o644,
			Content: content, ContentIdentity: identity,
		})
		before := preimages[relative].ContentIdentity
		action := fileAction(preimages[relative], PreimageRegular, before, identity)
		grouped := byPath[relative]
		for _, artifact := range grouped {
			ledger = append(ledger, ManagedEntry{
				ID: artifact.ID, Path: relative, Action: action, Kind: artifact.Kind,
				Module: artifact.Module, Template: artifact.Template, Version: artifact.Version,
				BeforeIdentity: before, AfterIdentity: identity,
				ContentIdentity: "sha256:" + artifact.Digest,
			})
		}
		if relative == manifestPath {
			ledger = append(ledger, ManagedEntry{
				ID: "manifest", Path: relative, Action: action, Kind: "manifest",
				Version: ManifestVersion, BeforeIdentity: before, AfterIdentity: identity,
				ContentIdentity: identity,
			})
		}
		for _, backup := range preservation.Backups {
			if backup.Path == relative {
				ledger = append(ledger, ManagedEntry{
					ID: "backup:" + backup.SourcePath, Path: relative, Action: action, Kind: "backup",
					BeforeIdentity: before, AfterIdentity: identity, ContentIdentity: backup.ContentIdentity,
				})
			}
		}
		if relative == repositoryRulesPath {
			ledger = append(ledger, ManagedEntry{
				ID: "repository-rules.readoption", Path: relative, Action: action,
				Kind: "repository-owned", BeforeIdentity: before, AfterIdentity: identity,
				ContentIdentity: identity,
			})
		}
	}
	for index := range ledger {
		ledger[index].Ordinal = index
	}
	return postimages, ledger, nil
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
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return fmt.Errorf("load Baseline catalog for plan validation: %w", err)
	}
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
	for index, retention := range document.Retention {
		if retention.FromClause == "" || retention.Enforcement == "" ||
			retention.Disposition == "" || retention.Targets == nil ||
			strings.TrimSpace(retention.Reason) == "" {
			return fmt.Errorf("invalid retention entry %d", index)
		}
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
	modules, artifacts, err := resolveManagedArtifacts(catalog, document.Profile, document.Decisions)
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
