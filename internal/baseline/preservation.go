package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	DecisionDocumentSchemaVersion = "setup-context-driven/decisions/0.0.1"
	DecisionDocumentVersion       = "0.0.1"
	repositoryRulesPath           = specificRepositoryPath
)

type PreservationMode string

const (
	PreservationModeGreenfield     PreservationMode = "greenfield"
	PreservationModePreservation   PreservationMode = "preservation"
	PreservationModeManagedRefresh PreservationMode = "managed-refresh"
)

type PreservationState string

const (
	PreservationStateReady          PreservationState = "ready"
	PreservationStateActionRequired PreservationState = "action_required"
	PreservationStateBlocked        PreservationState = "blocked"
)

// RootPreservationRequest contains the explicit instruction-preservation
// choice, optional locally materialized segmented source, and any
// owner-approved Decision Document.
type RootPreservationRequest struct {
	Mode           PreservationMode
	Decisions      *DecisionDocument
	SourceBaseline *ReadoptionSourceBaseline

	semanticOwners   SemanticOwnerRegistry
	managedArtifacts []plannedArtifact
	classifyCarriers bool
	classifications  []carrierClassification
}

// RootBackup is one immutable raw-byte backup selected before mutation.
type RootBackup struct {
	CarrierPath     string   `json:"carrierPath"`
	SourcePath      string   `json:"sourcePath"`
	Path            string   `json:"path"`
	ContentIdentity string   `json:"contentIdentity"`
	AlreadyExists   bool     `json:"alreadyExists"`
	Preimage        Preimage `json:"preimage"`
}

// ReadoptionSourceEntry is one byte-evidenced structural root instruction
// entry. SourceBytes uses JSON's canonical base64 byte encoding.
type ReadoptionSourceEntry struct {
	ID                   string         `json:"id"`
	Path                 string         `json:"path"`
	Carrier              string         `json:"carrier"`
	Kind                 string         `json:"kind"`
	Start                int            `json:"start"`
	End                  int            `json:"end"`
	Digest               string         `json:"digest"`
	CarrierDigest        string         `json:"carrierDigest"`
	SourceBytes          []byte         `json:"sourceBytes"`
	Encoding             string         `json:"encoding"`
	StructuralProvenance map[string]any `json:"structuralProvenance"`
}

// ReadoptionSourceBaseline is the immutable incompatible source identity
// presented for owner classification.
type ReadoptionSourceBaseline struct {
	ID               string                  `json:"id"`
	DeclaredIdentity string                  `json:"declaredIdentity"`
	Compatibility    string                  `json:"compatibility"`
	Digest           string                  `json:"digest"`
	CarrierCount     int                     `json:"carrierCount"`
	EntryCount       int                     `json:"entryCount"`
	ByteCount        int                     `json:"byteCount"`
	Entries          []ReadoptionSourceEntry `json:"entries"`
}

// DecisionValue is one catalog or repository decision.
type DecisionValue struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

// ReadoptionDestination is the strict union used by Source Baseline
// dispositions.
type ReadoptionDestination struct {
	ManagedID     string `json:"managedId,omitempty"`
	DocumentType  string `json:"documentType,omitempty"`
	Path          string `json:"path,omitempty"`
	Digest        string `json:"digest,omitempty"`
	ProposedBytes string `json:"proposedBytes,omitempty"`
}

// ReadoptionDisposition accounts for one Source Baseline Entry.
type ReadoptionDisposition struct {
	EntryID        string                 `json:"entryId"`
	EntryDigest    string                 `json:"entryDigest"`
	Classification string                 `json:"classification"`
	Disposition    string                 `json:"disposition"`
	Destination    *ReadoptionDestination `json:"destination"`
	Reason         string                 `json:"reason"`
}

// ReadoptionDecisions binds dispositions to one exact source identity.
type ReadoptionDecisions struct {
	SourceBaseline struct {
		ID     string `json:"id"`
		Digest string `json:"digest"`
	} `json:"sourceBaseline"`
	Dispositions []ReadoptionDisposition `json:"dispositions"`
}

// RepositoryRuleBlock is one accepted byte-exact repository-owned rule
// destined for an active semantic guide.
type RepositoryRuleBlock struct {
	ID             string `json:"id"`
	SourceEntryID  string `json:"sourceEntryId"`
	Classification string `json:"classification"`
	ManagedID      string `json:"managedId"`
	Path           string `json:"path"`
	Body           []byte `json:"body"`
}

// DecisionDocument is the maintained strict setup decision contract.
type DecisionDocument struct {
	SchemaVersion string               `json:"schemaVersion"`
	Version       string               `json:"version"`
	Decisions     []DecisionValue      `json:"decisions"`
	Readoption    *ReadoptionDecisions `json:"readoption,omitempty"`
}

// DecisionSkeleton is an editable, parser-valid manual classification input.
type DecisionSkeleton struct {
	Document   DecisionDocument `json:"document"`
	NextAction string           `json:"nextAction"`
}

// DecisionDocumentDiagnostic identifies one strict input failure.
type DecisionDocumentDiagnostic struct {
	Code    string
	Path    string
	ItemID  string
	Message string
}

// DecisionDocumentError contains every deterministic Decision Document
// diagnostic found at one validation boundary.
type DecisionDocumentError struct {
	Diagnostics []DecisionDocumentDiagnostic
}

func (e *DecisionDocumentError) Error() string {
	messages := make([]string, len(e.Diagnostics))
	for index, diagnostic := range e.Diagnostics {
		messages[index] = fmt.Sprintf(
			"%s: %s: %s: %s",
			diagnostic.Code,
			diagnostic.Path,
			diagnostic.ItemID,
			diagnostic.Message,
		)
	}
	return strings.Join(messages, "\n")
}

// RootPreservationPlan is the complete read-only root and recognized
// repository-rule preservation result consumed by portable-plan assembly.
type RootPreservationPlan struct {
	Mode                 PreservationMode         `json:"mode"`
	State                PreservationState        `json:"state"`
	Backups              []RootBackup             `json:"backups"`
	SourceBaseline       ReadoptionSourceBaseline `json:"sourceBaseline"`
	Dispositions         []ReadoptionDisposition  `json:"dispositions"`
	RepositoryRuleBlocks []RepositoryRuleBlock    `json:"repositoryRuleBlocks"`
	RepositoryRulesBytes []byte                   `json:"repositoryRulesBytes,omitempty"`
	Warnings             []Finding                `json:"warnings"`
	Findings             []Finding                `json:"findings"`
	DecisionSkeleton     *DecisionSkeleton        `json:"decisionSkeleton,omitempty"`
	NextAction           string                   `json:"nextAction,omitempty"`

	consumedRootPaths map[string]struct{}
}

// ParseDecisionDocument parses the maintained strict Decision Document
// schema. It rejects duplicate JSON keys before typed validation.
func ParseDecisionDocument(data []byte, sourcePath string) (DecisionDocument, error) {
	doc, jsonDiagnostics := decodeDocument(data, sourcePath)
	if len(jsonDiagnostics) != 0 {
		diagnostics := make([]DecisionDocumentDiagnostic, 0, len(jsonDiagnostics))
		for _, diagnostic := range jsonDiagnostics {
			code := "decision-file.json.invalid"
			itemID := "decision-file"
			if diagnostic.Code == "catalog.json.key.duplicate" {
				code = "decision-file.json.duplicate-key"
				itemID = diagnostic.Info
			}
			diagnostics = append(diagnostics, DecisionDocumentDiagnostic{
				Code:    code,
				Path:    sourcePath,
				ItemID:  itemID,
				Message: diagnostic.Info,
			})
		}
		return DecisionDocument{}, &DecisionDocumentError{Diagnostics: diagnostics}
	}
	if doc == nil {
		return DecisionDocument{}, decisionDocumentError(
			"decision-file.schema.invalid",
			sourcePath,
			"decision-file",
			"decision document must be a JSON object",
		)
	}
	return parseDecisionDocumentObject(doc, sourcePath)
}

func parseDecisionDocumentObject(doc document, sourcePath string) (DecisionDocument, error) {
	allowed := map[string]struct{}{
		"schemaVersion": {},
		"version":       {},
		"decisions":     {},
		"readoption":    {},
	}
	if !hasExactRequiredFields(doc, allowed, "schemaVersion", "version", "decisions") ||
		doc["schemaVersion"] != DecisionDocumentSchemaVersion ||
		doc["version"] != DecisionDocumentVersion {
		return DecisionDocument{}, decisionDocumentError(
			"decision-file.schema.invalid",
			sourcePath,
			"decision-file",
			"decision document requires schemaVersion, version, and decisions with only the optional readoption field",
		)
	}

	rawDecisions, ok := objectList(doc["decisions"])
	if !ok {
		return DecisionDocument{}, decisionDocumentError(
			"decision-file.decisions.invalid",
			sourcePath,
			"decisions",
			"decisions must be an ordered array",
		)
	}
	result := DecisionDocument{
		SchemaVersion: DecisionDocumentSchemaVersion,
		Version:       DecisionDocumentVersion,
		Decisions:     make([]DecisionValue, 0, len(rawDecisions)),
	}
	seenDecisions := make(map[string]struct{}, len(rawDecisions))
	for index, raw := range rawDecisions {
		itemID := fmt.Sprintf("decisions[%d]", index)
		if !hasExactFields(raw, "id", "value") {
			return DecisionDocument{}, decisionDocumentError(
				"decision-file.decision.invalid", sourcePath, itemID,
				"decision record requires exactly id and value",
			)
		}
		id, ok := raw["id"].(string)
		if !ok || !validBaselineIdentifier(id) {
			return DecisionDocument{}, decisionDocumentError(
				"decision-file.decision.invalid", sourcePath, itemID,
				"decision id is invalid",
			)
		}
		if _, duplicate := seenDecisions[id]; duplicate {
			return DecisionDocument{}, decisionDocumentError(
				"decision-file.decision.duplicate", sourcePath, id,
				"decision id appears more than once",
			)
		}
		seenDecisions[id] = struct{}{}
		result.Decisions = append(result.Decisions, DecisionValue{ID: id, Value: raw["value"]})
	}
	sort.Slice(result.Decisions, func(i, j int) bool { return result.Decisions[i].ID < result.Decisions[j].ID })

	if raw, exists := doc["readoption"]; exists {
		readoption, err := parseReadoptionDecisions(raw, sourcePath)
		if err != nil {
			return DecisionDocument{}, err
		}
		result.Readoption = &readoption
	}
	return result, nil
}

func parseReadoptionDecisions(value any, sourcePath string) (ReadoptionDecisions, error) {
	raw, ok := objectValue(value)
	if !ok || !hasExactFields(raw, "sourceBaseline", "dispositions") {
		return ReadoptionDecisions{}, decisionDocumentError(
			"decision-file.readoption.invalid", sourcePath, "readoption",
			"readoption requires exactly sourceBaseline and dispositions",
		)
	}
	source, ok := objectValue(raw["sourceBaseline"])
	if !ok || !hasExactFields(source, "id", "digest") {
		return ReadoptionDecisions{}, decisionDocumentError(
			"decision-file.readoption.source.invalid", sourcePath, "sourceBaseline",
			"sourceBaseline requires an id and lowercase SHA-256 digest",
		)
	}
	sourceID, idOK := source["id"].(string)
	sourceDigest, digestOK := source["digest"].(string)
	if !idOK || sourceID == "" || !digestOK || !isRawSHA256(sourceDigest) {
		return ReadoptionDecisions{}, decisionDocumentError(
			"decision-file.readoption.source.invalid", sourcePath, "sourceBaseline",
			"sourceBaseline requires an id and lowercase SHA-256 digest",
		)
	}
	rawDispositions, ok := objectList(raw["dispositions"])
	if !ok {
		return ReadoptionDecisions{}, decisionDocumentError(
			"decision-file.readoption.dispositions.invalid", sourcePath, "dispositions",
			"dispositions must be an ordered array",
		)
	}
	var result ReadoptionDecisions
	result.SourceBaseline.ID = sourceID
	result.SourceBaseline.Digest = sourceDigest
	result.Dispositions = make([]ReadoptionDisposition, 0, len(rawDispositions))
	for index, item := range rawDispositions {
		disposition, err := parseReadoptionDisposition(item, sourcePath, index)
		if err != nil {
			return ReadoptionDecisions{}, err
		}
		result.Dispositions = append(result.Dispositions, disposition)
	}
	return result, nil
}

func parseReadoptionDisposition(raw document, sourcePath string, index int) (ReadoptionDisposition, error) {
	itemID := fmt.Sprintf("dispositions[%d]", index)
	if !hasExactFields(
		raw,
		"entryId",
		"entryDigest",
		"classification",
		"disposition",
		"destination",
		"reason",
	) {
		return ReadoptionDisposition{}, decisionDocumentError(
			"readoption.disposition.invalid", sourcePath, itemID,
			"disposition record fields are invalid",
		)
	}
	entryID, entryOK := raw["entryId"].(string)
	entryDigest, digestOK := raw["entryDigest"].(string)
	classification, classificationOK := raw["classification"].(string)
	disposition, dispositionOK := raw["disposition"].(string)
	reason, reasonOK := raw["reason"].(string)
	if !entryOK || !strings.HasPrefix(entryID, "source-entry.") ||
		!isRawSHA256(strings.TrimPrefix(entryID, "source-entry.")) ||
		!digestOK || !isRawSHA256(entryDigest) ||
		!classificationOK || !containsString(
		[]string{"non-governed", "normative-clause", "operational-contract", "recommendation"},
		classification,
	) ||
		!dispositionOK || !containsString(
		[]string{"managed-entry", "rejected", "repository-document", "repository-rules"},
		disposition,
	) ||
		!reasonOK {
		return ReadoptionDisposition{}, decisionDocumentError(
			"readoption.disposition.invalid", sourcePath, itemID,
			"entry evidence, classification, disposition, or reason is invalid",
		)
	}
	if classification == "non-governed" && disposition != "rejected" {
		return ReadoptionDisposition{}, decisionDocumentError(
			"readoption.disposition.invalid", sourcePath, entryID,
			"non-governed evidence must use the rejected disposition",
		)
	}
	if (classification == "non-governed" || disposition == "rejected") && strings.TrimSpace(reason) == "" {
		return ReadoptionDisposition{}, decisionDocumentError(
			"readoption.disposition.reason.required", sourcePath, entryID,
			"non-governed or rejected evidence requires an individual reason",
		)
	}
	destination, err := parseReadoptionDestination(raw["destination"], disposition, sourcePath, entryID)
	if err != nil {
		return ReadoptionDisposition{}, err
	}
	return ReadoptionDisposition{
		EntryID:        entryID,
		EntryDigest:    entryDigest,
		Classification: classification,
		Disposition:    disposition,
		Destination:    destination,
		Reason:         reason,
	}, nil
}

func parseReadoptionDestination(
	value any,
	disposition string,
	sourcePath string,
	entryID string,
) (*ReadoptionDestination, error) {
	if disposition == "rejected" {
		if value != nil {
			return nil, decisionDocumentError(
				"readoption.disposition.invalid", sourcePath, entryID,
				"rejected disposition must have a null destination",
			)
		}
		return nil, nil
	}
	raw, ok := objectValue(value)
	if !ok {
		return nil, decisionDocumentError(
			"readoption.disposition.invalid", sourcePath, entryID,
			"non-rejected disposition requires a typed destination object",
		)
	}
	switch disposition {
	case "managed-entry":
		managedID, ok := raw["managedId"].(string)
		if !ok || managedID == "" || !hasExactFields(raw, "managedId") {
			return nil, decisionDocumentError(
				"readoption.disposition.invalid", sourcePath, entryID,
				"managed-entry destination requires exactly managedId",
			)
		}
		return &ReadoptionDestination{ManagedID: managedID}, nil
	case "repository-document":
		if !hasExactFields(raw, "documentType", "path", "digest") {
			return nil, decisionDocumentError(
				"readoption.disposition.invalid", sourcePath, entryID,
				"repository-document destination fields are invalid",
			)
		}
		documentType, typeOK := raw["documentType"].(string)
		targetPath, pathOK := raw["path"].(string)
		digest, digestOK := raw["digest"].(string)
		if !typeOK || !containsString(
			[]string{"agent-guide", "architecture-decision", "design-contract", "domain-context", "http-contract"},
			documentType,
		) {
			return nil, decisionDocumentError(
				"readoption.destination.document-type.invalid", sourcePath, entryID,
				"repository documentType is not supported",
			)
		}
		if !pathOK || !repositoryPathIsSafe(targetPath) || !digestOK || !isRawSHA256(digest) {
			return nil, decisionDocumentError(
				"readoption.disposition.invalid", sourcePath, entryID,
				"repository-document destination path or digest is invalid",
			)
		}
		return &ReadoptionDestination{
			DocumentType: documentType,
			Path:         targetPath,
			Digest:       digest,
		}, nil
	default:
		if !hasExactFields(raw, "documentType", "path", "proposedBytes", "digest") {
			return nil, decisionDocumentError(
				"readoption.disposition.invalid", sourcePath, entryID,
				"repository-rules destination fields are invalid",
			)
		}
		documentType, typeOK := raw["documentType"].(string)
		targetPath, pathOK := raw["path"].(string)
		proposedText, bytesOK := raw["proposedBytes"].(string)
		digest, digestOK := raw["digest"].(string)
		if pathOK &&
			(targetPath == legacyRepositoryPath || targetPath == legacyRepositoryRulesPath) {
			targetPath = repositoryRulesPath
		}
		if !typeOK || documentType != "repository-rules" ||
			!pathOK || targetPath != repositoryRulesPath ||
			!bytesOK || !digestOK || !isRawSHA256(digest) {
			return nil, decisionDocumentError(
				"readoption.destination.repository-rules.invalid", sourcePath, entryID,
				"Repository-Specific Normative Rules require the default typed path and exact bytes",
			)
		}
		proposed, err := base64.StdEncoding.DecodeString(proposedText)
		if err != nil || base64.StdEncoding.EncodeToString(proposed) != proposedText {
			return nil, decisionDocumentError(
				"readoption.destination.proposed-bytes.invalid", sourcePath, entryID,
				"proposedBytes must be canonical base64",
			)
		}
		sum := sha256.Sum256(proposed)
		if len(proposed) == 0 || hex.EncodeToString(sum[:]) != digest {
			return nil, decisionDocumentError(
				"readoption.destination.proposed-bytes.stale", sourcePath, entryID,
				"proposed bytes are empty or do not match the declared digest",
			)
		}
		if bytes.Contains(proposed, []byte("<!-- setup-context-driven:")) {
			return nil, decisionDocumentError(
				"readoption.destination.proposed-bytes.managed-marker", sourcePath, entryID,
				"Repository-Specific Normative Rules proposed bytes must remain unmarked",
			)
		}
		return &ReadoptionDestination{
			DocumentType:  documentType,
			Path:          targetPath,
			Digest:        digest,
			ProposedBytes: proposedText,
		}, nil
	}
}

// PlanRootPreservation resolves root instructions and the three recognized
// repository-rule carriers. It never mutates the repository.
func PlanRootPreservation(
	inspection RepositoryInspection,
	request RootPreservationRequest,
) (RootPreservationPlan, error) {
	return planRootPreservationWithCatalog(inspection, request, nil)
}

func planRootPreservationWithCatalog(
	inspection RepositoryInspection,
	request RootPreservationRequest,
	catalog *Catalog,
) (RootPreservationPlan, error) {
	if request.Mode != PreservationModeGreenfield &&
		request.Mode != PreservationModePreservation &&
		request.Mode != PreservationModeManagedRefresh {
		return RootPreservationPlan{}, fmt.Errorf(
			"plan root-instruction preservation: unsupported mode %q",
			request.Mode,
		)
	}
	warnings := append([]Finding(nil), inspection.Snapshot.Warnings...)
	var carrierClassifications []carrierClassification
	if request.classifyCarriers && catalog != nil {
		carrierClassifications = request.classifications
		if carrierClassifications == nil {
			var err error
			carrierClassifications, err = classifyCarriers(
				inspection.Root,
				inspection.Snapshot.Carriers,
				catalog,
				request.managedArtifacts,
			)
			if err != nil {
				return RootPreservationPlan{}, err
			}
		}
		warnings = warningsForCarrierClassifications(warnings, carrierClassifications)
	}
	plan := RootPreservationPlan{
		Mode:     request.Mode,
		State:    PreservationStateReady,
		Warnings: warnings,
		Findings: append([]Finding(nil), inspection.Snapshot.Blocking...),
	}
	rootSources, findings, err := loadRootPreservationSources(inspection)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	plan.Findings = append(plan.Findings, findings...)
	if request.Mode == PreservationModeManagedRefresh {
		markerFindings, err := managedRefreshMarkerFindings(inspection.Root)
		if err != nil {
			return RootPreservationPlan{}, err
		}
		plan.Findings = sortedFindings(append(plan.Findings, markerFindings...))
		if len(plan.Findings) != 0 {
			plan.State = PreservationStateBlocked
			plan.NextAction = "repair every blocking root carrier or modified managed marker and rerun Baseline planning"
		}
		return plan, nil
	}
	retainsRepositoryRules, err := currentSetupRetainsRecognizedRepositoryRulesWithCatalog(
		inspection.Root,
		catalog,
	)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	if request.Mode == PreservationModePreservation && retainsRepositoryRules {
		rootSources, plan.consumedRootPaths, err = excludePreviouslyBackedUpRootGuidance(
			inspection.Root,
			rootSources,
		)
		if err != nil {
			return RootPreservationPlan{}, err
		}
	}
	migrationRootSources := rootSourcesWithUnmarkedGuidance(rootSources)
	plan.Backups, findings, err = planRootBackups(inspection.Root, rootSources)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	plan.Findings = append(plan.Findings, findings...)
	var repositorySources []rootPreservationSource
	if !retainsRepositoryRules {
		repositorySources, findings, err = loadRecognizedRepositoryRuleSources(inspection.Root)
		if err != nil {
			return RootPreservationPlan{}, err
		}
		plan.Findings = append(plan.Findings, findings...)
	}
	staleManagedSources, findings, err := loadStaleManagedCarrierSources(
		inspection,
		staleManagedCarrierPaths(carrierClassifications),
	)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	plan.Findings = append(plan.Findings, findings...)
	sources := append(migrationRootSources, repositorySources...)
	sources = append(sources, staleManagedSources...)
	plan.SourceBaseline = buildReadoptionSourceBaseline(sources)
	plan.Findings = sortedFindings(plan.Findings)

	if len(plan.Findings) != 0 {
		plan.State = PreservationStateBlocked
		plan.NextAction = "repair every blocking root carrier or backup collision and rerun Baseline planning"
		return plan, nil
	}
	if request.Mode == PreservationModeGreenfield && len(staleManagedSources) == 0 {
		return plan, nil
	}
	if len(plan.SourceBaseline.Entries) == 0 {
		return plan, nil
	}
	if request.SourceBaseline != nil {
		if err := validateClassifiedSourceBaseline(plan.SourceBaseline, *request.SourceBaseline); err != nil {
			plan.State = PreservationStateActionRequired
			plan.Findings = []Finding{{
				Code:    "baseline.preservation.classified-source.invalid",
				Path:    ".",
				Message: err.Error(),
			}}
			plan.NextAction = "rerun segmentation and classification against the current Source Baseline"
			return plan, nil
		}
		plan.SourceBaseline = cloneReadoptionSourceBaseline(*request.SourceBaseline)
	}
	if request.Decisions == nil {
		skeleton := buildDecisionSkeleton(plan.SourceBaseline)
		plan.State = PreservationStateActionRequired
		plan.DecisionSkeleton = &skeleton
		plan.NextAction = skeleton.NextAction
		return plan, nil
	}
	dispositions, repositoryBlocks, repositoryRules, decisionFindings := validatePreservationDecisions(
		inspection.Root,
		plan.SourceBaseline,
		*request.Decisions,
		request.semanticOwners,
	)
	plan.Dispositions = dispositions
	plan.RepositoryRuleBlocks = repositoryBlocks
	plan.RepositoryRulesBytes = repositoryRules
	plan.Findings = sortedFindings(decisionFindings)
	if len(plan.Findings) != 0 {
		plan.State = PreservationStateActionRequired
		skeleton := buildDecisionSkeleton(plan.SourceBaseline)
		plan.DecisionSkeleton = &skeleton
		plan.NextAction = "complete or correct every root-rule disposition and rerun Baseline planning"
	}
	return plan, nil
}

func managedRefreshMarkerFindings(
	root string,
) ([]Finding, error) {
	manifestBytes, err := readOptionalRegular(root, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("inspect managed-refresh Setup Manifest: %w", err)
	}
	manifest, valid := parseManagedSetupManifest(manifestBytes)
	if !valid {
		return nil, nil
	}
	sourceByPath := manifestArtifactsByPath(manifest.ManagedArtifacts)
	var findings []Finding
	for _, relative := range sortedKeys(sourceByPath) {
		content, err := readOptionalRegular(root, relative)
		if err != nil {
			return nil, err
		}
		if content == nil {
			continue
		}
		entriesByID := make(map[string][]ReadoptionSourceEntry)
		for _, entry := range partitionRootSource(relative, content) {
			if entry.Kind != "managed-block" {
				continue
			}
			managedID, _ := entry.StructuralProvenance["managedId"].(string)
			entriesByID[managedID] = append(entriesByID[managedID], entry)
		}
		modified := false
		for managedID, artifact := range sourceByPath[relative] {
			entries := entriesByID[managedID]
			if len(entries) != 1 || !managedEntryMatchesManifest(entries[0], artifact) {
				modified = true
				break
			}
		}
		if !modified {
			continue
		}
		findings = append(findings, Finding{
			Code:    "baseline.preservation.managed-marker.modified",
			Path:    relative,
			Message: "managed marker bytes no longer match the adopted Setup Manifest",
		})
	}
	return sortedFindings(findings), nil
}

type rootPreservationSource struct {
	carrierPath     string
	sourcePath      string
	contentIdentity string
	content         []byte

	classificationEntries []ReadoptionSourceEntry
}

func loadRootPreservationSources(
	inspection RepositoryInspection,
) ([]rootPreservationSource, []Finding, error) {
	anchored, err := os.OpenRoot(inspection.Root)
	if err != nil {
		return nil, nil, fmt.Errorf("open repository root for preservation: %w", err)
	}
	defer anchored.Close()

	selectedCarriers := make(map[string]InstructionCarrier)
	for _, carrier := range inspection.Snapshot.Carriers {
		if carrier.Scope != "root" || carrier.Kind == CarrierOpaque {
			continue
		}
		selected, exists := selectedCarriers[carrier.TargetPath]
		if !exists ||
			carrier.Kind == CarrierRegular && selected.Kind != CarrierRegular ||
			carrier.Kind == selected.Kind && carrier.Path < selected.Path {
			selectedCarriers[carrier.TargetPath] = carrier
		}
	}
	targets := sortedMapKeys(selectedCarriers)
	var sources []rootPreservationSource
	var findings []Finding
	for _, target := range targets {
		carrier := selectedCarriers[target]
		content, readErr := readRootRegularFile(anchored, carrier.TargetPath)
		if readErr != nil {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.source.stale",
				Path:    carrier.TargetPath,
				Message: readErr.Error(),
			})
			continue
		}
		sum := sha256.Sum256(content)
		identity := "sha256:" + hex.EncodeToString(sum[:])
		if identity != carrier.ContentIdentity {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.source.stale",
				Path:    carrier.TargetPath,
				Message: "root instruction bytes changed after repository inspection",
			})
			continue
		}
		sources = append(sources, rootPreservationSource{
			carrierPath:     carrier.Path,
			sourcePath:      carrier.TargetPath,
			contentIdentity: identity,
			content:         content,
		})
	}
	sort.Slice(sources, func(i, j int) bool {
		if sources[i].sourcePath != sources[j].sourcePath {
			return sources[i].sourcePath < sources[j].sourcePath
		}
		return sources[i].carrierPath < sources[j].carrierPath
	})
	return sources, findings, nil
}

func loadRecognizedRepositoryRuleSources(
	rootPath string,
) ([]rootPreservationSource, []Finding, error) {
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open repository root for recognized rule carriers: %w", err)
	}
	defer anchored.Close()

	var sources []rootPreservationSource
	var findings []Finding
	for _, relative := range []string{
		specificRepositoryPath,
		legacyRepositoryPath,
		legacyRepositoryRulesPath,
	} {
		content, exists, finding := readSpecificRepositoryCarrier(anchored, relative)
		if finding != nil {
			findings = append(findings, *finding)
			continue
		}
		if !exists || repositoryCarrierEmpty(content, relative == legacyRepositoryPath) {
			continue
		}
		sum := sha256.Sum256(content)
		sources = append(sources, rootPreservationSource{
			carrierPath:     relative,
			sourcePath:      relative,
			contentIdentity: "sha256:" + hex.EncodeToString(sum[:]),
			content:         content,
		})
	}
	return sources, findings, nil
}

func loadStaleManagedCarrierSources(
	inspection RepositoryInspection,
	stalePaths map[string]struct{},
) ([]rootPreservationSource, []Finding, error) {
	if len(stalePaths) == 0 {
		return nil, nil, nil
	}
	anchored, err := os.OpenRoot(inspection.Root)
	if err != nil {
		return nil, nil, fmt.Errorf("open repository root for stale managed carriers: %w", err)
	}
	defer anchored.Close()

	var sources []rootPreservationSource
	var findings []Finding
	for _, carrier := range inspection.Snapshot.Carriers {
		if _, stale := stalePaths[carrier.Path]; !stale {
			continue
		}
		content, ok := readCarrierClassificationBytes(anchored, carrier)
		if !ok {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.source.stale",
				Path:    carrier.Path,
				Message: "stale managed carrier bytes changed after repository inspection",
			})
			continue
		}
		sources = append(sources, rootPreservationSource{
			carrierPath:           carrier.Path,
			sourcePath:            carrier.Path,
			contentIdentity:       carrier.ContentIdentity,
			content:               content,
			classificationEntries: partitionRootSource(carrier.Path, content),
		})
	}
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].sourcePath < sources[j].sourcePath
	})
	return sources, findings, nil
}

func currentSetupRetainsRecognizedRepositoryRules(rootPath string) (bool, error) {
	return currentSetupRetainsRecognizedRepositoryRulesWithCatalog(rootPath, nil)
}

func currentSetupRetainsRecognizedRepositoryRulesWithCatalog(
	rootPath string,
	catalog *Catalog,
) (bool, error) {
	manifestBytes, err := readOptionalRegular(rootPath, manifestPath)
	if err != nil {
		return false, fmt.Errorf("read current Setup Manifest ownership: %w", err)
	}
	if len(manifestBytes) == 0 {
		return false, nil
	}
	var manifest SetupManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return false, nil
	}
	if manifest.SchemaVersion != ManifestSchema ||
		manifest.Version != ManifestVersion ||
		manifest.Generator.Skill != "setup-context-driven" ||
		manifest.Generator.Version != ManifestVersion ||
		manifest.Generator.Baseline != "baseline."+manifest.Profile+"-"+ManifestVersion {
		return false, nil
	}
	if catalog == nil {
		catalog, err = LoadEmbeddedCatalog()
		if err != nil {
			return false, fmt.Errorf("load catalog for Setup Manifest ownership: %w", err)
		}
	}
	if manifest.CatalogDigest != catalog.Digest() {
		return false, nil
	}
	profile, err := ResolveProfile(rootPath, manifest.Profile, catalog)
	if err != nil || profile.Digest != manifest.ProfileDigest {
		return false, nil
	}
	return true, nil
}

func rootSourcesWithUnmarkedGuidance(
	sources []rootPreservationSource,
) []rootPreservationSource {
	selected := make([]rootPreservationSource, 0, len(sources))
	for _, source := range sources {
		if len(classificationEntries(source)) == 0 {
			continue
		}
		selected = append(selected, source)
	}
	return selected
}

var agentsBackupName = regexp.MustCompile(`^AGENTS\.([a-f0-9]{64})\.md$`)

func excludePreviouslyBackedUpRootGuidance(
	rootPath string,
	sources []rootPreservationSource,
) ([]rootPreservationSource, map[string]struct{}, error) {
	payloads, err := verifiedAgentsBackupPayloads(rootPath)
	if err != nil {
		return nil, nil, err
	}
	consumed := make(map[string]struct{})
	for index := range sources {
		source := &sources[index]
		if source.carrierPath != "AGENTS.md" ||
			source.sourcePath != "AGENTS.md" ||
			!containsSetupManagedGuidance(source.content) {
			continue
		}
		entries, excluded := excludeBackedUpPayloads(
			unmarkedSourceEntries(source.sourcePath, source.content),
			payloads,
		)
		source.classificationEntries = entries
		if excluded {
			consumed[source.sourcePath] = struct{}{}
		}
	}
	return sources, consumed, nil
}

func verifiedAgentsBackupPayloads(rootPath string) ([][]byte, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("list prior AGENTS backups: %w", err)
	}
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, fmt.Errorf("open repository root for prior AGENTS backups: %w", err)
	}
	defer anchored.Close()

	var payloads [][]byte
	for _, entry := range entries {
		match := agentsBackupName.FindStringSubmatch(entry.Name())
		if len(match) != 2 || entry.Type()&fs.ModeSymlink != 0 {
			continue
		}
		payload, readErr := readRootRegularFile(anchored, entry.Name())
		if readErr != nil {
			continue
		}
		sum := sha256.Sum256(payload)
		if hex.EncodeToString(sum[:]) != match[1] {
			continue
		}
		payloads = append(payloads, payload)
	}
	sort.Slice(payloads, func(i, j int) bool {
		if len(payloads[i]) != len(payloads[j]) {
			return len(payloads[i]) > len(payloads[j])
		}
		return bytes.Compare(payloads[i], payloads[j]) < 0
	})
	return payloads, nil
}

func excludeBackedUpPayloads(
	entries []ReadoptionSourceEntry,
	payloads [][]byte,
) ([]ReadoptionSourceEntry, bool) {
	result := make([]ReadoptionSourceEntry, 0, len(entries))
	var excluded bool
	for _, entry := range entries {
		ranges := [][2]int{{0, len(entry.SourceBytes)}}
		for _, payload := range payloads {
			if len(payload) == 0 {
				continue
			}
			var remaining [][2]int
			for _, current := range ranges {
				cursor := current[0]
				for cursor < current[1] {
					offset := bytes.Index(entry.SourceBytes[cursor:current[1]], payload)
					if offset < 0 {
						break
					}
					start := cursor + offset
					if start > cursor {
						remaining = append(remaining, [2]int{cursor, start})
					}
					cursor = start + len(payload)
					excluded = true
				}
				if cursor < current[1] {
					remaining = append(remaining, [2]int{cursor, current[1]})
				}
			}
			ranges = remaining
		}
		for _, current := range ranges {
			sourceBytes := entry.SourceBytes[current[0]:current[1]]
			if len(bytes.TrimSpace(sourceBytes)) == 0 {
				continue
			}
			result = append(result, newReadoptionSourceEntry(
				entry.Path,
				entry.Kind,
				entry.Start+current[0],
				entry.Start+current[1],
				sourceBytes,
				entry.CarrierDigest,
				entry.StructuralProvenance,
			))
		}
	}
	return result, excluded
}

func containsSetupManagedGuidance(content []byte) bool {
	for _, entry := range partitionRootSource("", content) {
		if entry.Kind == "managed-block" {
			return true
		}
	}
	return false
}

func classificationEntries(source rootPreservationSource) []ReadoptionSourceEntry {
	if source.classificationEntries != nil {
		return source.classificationEntries
	}
	return unmarkedSourceEntries(source.sourcePath, source.content)
}

func validateClassifiedSourceBaseline(
	current ReadoptionSourceBaseline,
	classified ReadoptionSourceBaseline,
) error {
	if classified.ID != current.ID ||
		classified.DeclaredIdentity != current.DeclaredIdentity ||
		classified.Compatibility != current.Compatibility ||
		classified.Digest != current.Digest ||
		classified.CarrierCount != current.CarrierCount ||
		classified.ByteCount != current.ByteCount ||
		classified.EntryCount != len(classified.Entries) {
		return errors.New("classified Source Baseline identity does not match current repository bytes")
	}
	if _, err := NewAnalysisSnapshot(classified); err != nil {
		return fmt.Errorf("classified Source Baseline entries are invalid: %w", err)
	}

	classifiedIndex := 0
	for _, original := range current.Entries {
		cursor := original.Start
		consumed := 0
		for classifiedIndex < len(classified.Entries) {
			entry := classified.Entries[classifiedIndex]
			if entry.Path != original.Path || entry.Start >= original.End && original.End != original.Start {
				break
			}
			if original.Start == original.End {
				if entry.Start != original.Start || entry.End != original.End {
					break
				}
			} else if entry.Start != cursor || entry.End <= entry.Start || entry.End > original.End {
				return fmt.Errorf("classified entry %q leaves a gap, overlaps, or escapes its structural source", entry.ID)
			}
			if entry.Kind != original.Kind ||
				entry.CarrierDigest != original.CarrierDigest ||
				!reflectJSONEqual(entry.StructuralProvenance, original.StructuralProvenance) {
				return fmt.Errorf("classified entry %q changes structural provenance", entry.ID)
			}
			start := entry.Start - original.Start
			end := entry.End - original.Start
			if start < 0 || end < start || end > len(original.SourceBytes) ||
				!bytes.Equal(entry.SourceBytes, original.SourceBytes[start:end]) {
				return fmt.Errorf("classified entry %q changes source bytes", entry.ID)
			}
			cursor = entry.End
			classifiedIndex++
			consumed++
			if cursor == original.End {
				break
			}
		}
		if consumed == 0 || cursor != original.End {
			return fmt.Errorf("classified entries do not cover source entry %q exactly once", original.ID)
		}
	}
	if classifiedIndex != len(classified.Entries) {
		return fmt.Errorf("classified entry %q has no current structural source", classified.Entries[classifiedIndex].ID)
	}
	return nil
}

func planRootBackups(
	rootPath string,
	sources []rootPreservationSource,
) ([]RootBackup, []Finding, error) {
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open repository root for backup planning: %w", err)
	}
	defer anchored.Close()

	backups := make([]RootBackup, 0, len(sources))
	var findings []Finding
	for _, source := range sources {
		digest := strings.TrimPrefix(source.contentIdentity, "sha256:")
		base := strings.TrimSuffix(source.carrierPath, path.Ext(source.carrierPath))
		backupPath := base + "." + digest + ".md"
		backup := RootBackup{
			CarrierPath:     source.carrierPath,
			SourcePath:      source.sourcePath,
			Path:            backupPath,
			ContentIdentity: source.contentIdentity,
			Preimage:        Preimage{Path: backupPath, Kind: PreimageMissing},
		}
		info, lstatErr := anchored.Lstat(filepath.FromSlash(backupPath))
		switch {
		case errors.Is(lstatErr, fs.ErrNotExist):
		case lstatErr != nil:
			findings = append(findings, Finding{
				Code:    "baseline.preservation.backup.unreadable",
				Path:    backupPath,
				Message: "content-addressed backup path cannot be inspected",
			})
		case !info.Mode().IsRegular():
			backup.Preimage = preimageFromInfo(backupPath, info, "")
			findings = append(findings, Finding{
				Code:    "baseline.preservation.backup.collision",
				Path:    backupPath,
				Message: "content-addressed backup path is not a regular file",
			})
		default:
			backup.Preimage = preimageFromInfo(backupPath, info, "")
			existing, readErr := readRootRegularFile(anchored, backupPath)
			if readErr != nil {
				findings = append(findings, Finding{
					Code:    "baseline.preservation.backup.unreadable",
					Path:    backupPath,
					Message: readErr.Error(),
				})
				break
			}
			existingSum := sha256.Sum256(existing)
			backup.Preimage.ContentIdentity = "sha256:" + hex.EncodeToString(existingSum[:])
			if !bytes.Equal(existing, source.content) {
				findings = append(findings, Finding{
					Code:    "baseline.preservation.backup.collision",
					Path:    backupPath,
					Message: "existing content-addressed backup bytes do not match their full SHA-256 identity",
				})
				break
			}
			backup.AlreadyExists = true
		}
		backups = append(backups, backup)
	}
	return backups, findings, nil
}

func readRootRegularFile(root *os.Root, relative string) ([]byte, error) {
	info, err := root.Lstat(filepath.FromSlash(relative))
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", relative, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%q must be a regular non-symlink file", relative)
	}
	file, err := root.Open(filepath.FromSlash(relative))
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", relative, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect opened %q: %w", relative, err)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return nil, fmt.Errorf("%q changed while it was inspected", relative)
	}
	data, err := io.ReadAll(io.LimitReader(file, maxInventoryFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", relative, err)
	}
	if len(data) > maxInventoryFileBytes {
		return nil, fmt.Errorf("%q exceeds %d bytes", relative, maxInventoryFileBytes)
	}
	return data, nil
}

func buildReadoptionSourceBaseline(sources []rootPreservationSource) ReadoptionSourceBaseline {
	identity := sha256.New()
	var entries []ReadoptionSourceEntry
	byteCount := 0
	carrierCount := 0
	for _, source := range sources {
		sourceEntries := classificationEntries(source)
		if len(sourceEntries) == 0 {
			continue
		}
		writeLengthPrefixed(identity, source.sourcePath)
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(source.content)))
		_, _ = identity.Write(length[:])
		_, _ = identity.Write(source.content)
		carrierCount++
		for _, entry := range sourceEntries {
			byteCount += len(entry.SourceBytes)
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Path != entries[j].Path {
			return entries[i].Path < entries[j].Path
		}
		if entries[i].Start != entries[j].Start {
			return entries[i].Start < entries[j].Start
		}
		if entries[i].End != entries[j].End {
			return entries[i].End < entries[j].End
		}
		return entries[i].Kind < entries[j].Kind
	})
	digest := hex.EncodeToString(identity.Sum(nil))
	return ReadoptionSourceBaseline{
		ID:               "baseline.readoption." + digest,
		DeclaredIdentity: "unconfigured",
		Compatibility:    "incompatible",
		Digest:           digest,
		CarrierCount:     carrierCount,
		EntryCount:       len(entries),
		ByteCount:        byteCount,
		Entries:          entries,
	}
}

var managedBeginMarker = regexp.MustCompile(
	`<!--\s*setup-context-driven:begin\s+id=([A-Za-z0-9_.-]+)\s+version=([0-9]+(?:\.[0-9]+)*)\s*-->`,
)

type sourceSpan struct {
	start      int
	end        int
	kind       string
	provenance map[string]any
}

func partitionRootSource(sourcePath string, content []byte) []ReadoptionSourceEntry {
	carrierSum := sha256.Sum256(content)
	carrierDigest := hex.EncodeToString(carrierSum[:])
	var spans []sourceSpan
	cursor := 0
	for {
		location := managedBeginMarker.FindSubmatchIndex(content[cursor:])
		if location == nil {
			break
		}
		start := cursor + location[0]
		openEnd := cursor + location[1]
		managedID := string(content[cursor+location[2] : cursor+location[3]])
		versionText := string(content[cursor+location[4] : cursor+location[5]])
		closePattern := regexp.MustCompile(
			`<!--\s*setup-context-driven:end\s+id=` + regexp.QuoteMeta(managedID) + `\s*-->`,
		)
		closeLocation := closePattern.FindIndex(content[openEnd:])
		if closeLocation == nil {
			cursor = openEnd
			continue
		}
		end := openEnd + closeLocation[1]
		if bytes.HasPrefix(content[end:], []byte("\r\n")) {
			end += 2
		} else if bytes.HasPrefix(content[end:], []byte("\n")) {
			end++
		}
		var version any = versionText
		if integerVersion, err := strconv.Atoi(versionText); err == nil {
			version = integerVersion
		}
		spans = append(spans, sourceSpan{
			start: start,
			end:   end,
			kind:  "managed-block",
			provenance: map[string]any{
				"managedId":     managedID,
				"markerVersion": version,
			},
		})
		cursor = end
	}
	if len(spans) == 0 {
		return []ReadoptionSourceEntry{
			newReadoptionSourceEntry(
				sourcePath,
				"unmarked-span",
				0,
				len(content),
				content,
				carrierDigest,
				map[string]any{"markerState": "unmarked"},
			),
		}
	}
	var entries []ReadoptionSourceEntry
	cursor = 0
	for _, span := range spans {
		if span.start > cursor {
			entries = append(entries, newReadoptionSourceEntry(
				sourcePath,
				"unmarked-span",
				cursor,
				span.start,
				content[cursor:span.start],
				carrierDigest,
				map[string]any{"markerState": "unmarked"},
			))
		}
		entries = append(entries, newReadoptionSourceEntry(
			sourcePath,
			span.kind,
			span.start,
			span.end,
			content[span.start:span.end],
			carrierDigest,
			span.provenance,
		))
		cursor = span.end
	}
	if cursor < len(content) {
		entries = append(entries, newReadoptionSourceEntry(
			sourcePath,
			"unmarked-span",
			cursor,
			len(content),
			content[cursor:],
			carrierDigest,
			map[string]any{"markerState": "unmarked"},
		))
	}
	return entries
}

func unmarkedSourceEntries(sourcePath string, content []byte) []ReadoptionSourceEntry {
	partitioned := partitionRootSource(sourcePath, content)
	entries := make([]ReadoptionSourceEntry, 0, len(partitioned))
	for _, entry := range partitioned {
		if entry.Kind == "managed-block" || len(bytes.TrimSpace(entry.SourceBytes)) == 0 {
			continue
		}
		entries = append(entries, entry)
	}
	return entries
}

func containsOnlySetupManagedGuidance(content []byte) bool {
	entries := partitionRootSource("", content)
	hasManaged := false
	for _, entry := range entries {
		if entry.Kind == "managed-block" {
			hasManaged = true
			continue
		}
		if len(bytes.TrimSpace(entry.SourceBytes)) != 0 {
			return false
		}
	}
	return hasManaged
}

func unchangedSetupManagedGuidance(sourceIdentity string, rendered []byte) bool {
	return sourceIdentity != "" &&
		planContentIdentity(rendered) == sourceIdentity &&
		containsOnlySetupManagedGuidance(rendered)
}

func newReadoptionSourceEntry(
	sourcePath string,
	kind string,
	start int,
	end int,
	sourceBytes []byte,
	carrierDigest string,
	provenance map[string]any,
) ReadoptionSourceEntry {
	sum := sha256.Sum256(sourceBytes)
	digest := hex.EncodeToString(sum[:])
	identity := sha256.New()
	provenanceJSON, _ := json.Marshal(provenance)
	for _, value := range []string{
		sourcePath,
		kind,
		strconv.Itoa(start),
		strconv.Itoa(end),
		digest,
		string(provenanceJSON),
	} {
		writeLengthPrefixed(identity, value)
	}
	return ReadoptionSourceEntry{
		ID:                   "source-entry." + hex.EncodeToString(identity.Sum(nil)),
		Path:                 sourcePath,
		Carrier:              sourcePath,
		Kind:                 kind,
		Start:                start,
		End:                  end,
		Digest:               digest,
		CarrierDigest:        carrierDigest,
		SourceBytes:          append([]byte(nil), sourceBytes...),
		Encoding:             "base64",
		StructuralProvenance: cloneMap(provenance),
	}
}

func buildDecisionSkeleton(source ReadoptionSourceBaseline) DecisionSkeleton {
	document := DecisionDocument{
		SchemaVersion: DecisionDocumentSchemaVersion,
		Version:       DecisionDocumentVersion,
		Decisions:     []DecisionValue{},
		Readoption:    &ReadoptionDecisions{},
	}
	document.Readoption.SourceBaseline.ID = source.ID
	document.Readoption.SourceBaseline.Digest = source.Digest
	document.Readoption.Dispositions = make([]ReadoptionDisposition, 0, len(source.Entries))
	for _, entry := range source.Entries {
		if entry.Kind == "managed-block" || len(bytes.TrimSpace(entry.SourceBytes)) == 0 {
			document.Readoption.Dispositions = append(
				document.Readoption.Dispositions,
				ReadoptionDisposition{
					EntryID:        entry.ID,
					EntryDigest:    entry.Digest,
					Classification: "non-governed",
					Disposition:    "rejected",
					Destination:    nil,
					Reason:         "Structural or previously managed evidence is not proposed as a repository-specific rule.",
				},
			)
			continue
		}
		sum := sha256.Sum256(entry.SourceBytes)
		document.Readoption.Dispositions = append(
			document.Readoption.Dispositions,
			ReadoptionDisposition{
				EntryID:        entry.ID,
				EntryDigest:    entry.Digest,
				Classification: "normative-clause",
				Disposition:    "repository-rules",
				Destination: &ReadoptionDestination{
					DocumentType:  "repository-rules",
					Path:          repositoryRulesPath,
					Digest:        hex.EncodeToString(sum[:]),
					ProposedBytes: base64.StdEncoding.EncodeToString(entry.SourceBytes),
				},
				Reason: "Preserve this source entry as a Repository-Specific Normative Rule for human review.",
			},
		)
	}
	return DecisionSkeleton{
		Document:   document,
		NextAction: "review and edit every proposed root-rule classification, then supply the complete Decision Document",
	}
}

func validatePreservationDecisions(
	rootPath string,
	source ReadoptionSourceBaseline,
	document DecisionDocument,
	semanticOwners SemanticOwnerRegistry,
) ([]ReadoptionDisposition, []RepositoryRuleBlock, []byte, []Finding) {
	var findings []Finding
	if document.SchemaVersion != DecisionDocumentSchemaVersion ||
		document.Version != DecisionDocumentVersion ||
		document.Readoption == nil {
		findings = append(findings, Finding{
			Code:    "baseline.preservation.decision-document.invalid",
			Path:    ".",
			Message: "preservation requires a complete strict Decision Document with Readoption dispositions",
		})
		return nil, nil, nil, findings
	}
	if document.Readoption.SourceBaseline.ID != source.ID ||
		document.Readoption.SourceBaseline.Digest != source.Digest {
		findings = append(findings, Finding{
			Code:    "baseline.preservation.source.stale",
			Path:    ".",
			Message: "Decision Document Source Baseline identity does not match current root bytes",
		})
	}
	expected := make(map[string]ReadoptionSourceEntry, len(source.Entries))
	order := make(map[string]int, len(source.Entries))
	for index, entry := range source.Entries {
		expected[entry.ID] = entry
		order[entry.ID] = index
	}
	seen := make(map[string]int, len(document.Readoption.Dispositions))
	accepted := make([]ReadoptionDisposition, 0, len(source.Entries))
	for _, disposition := range document.Readoption.Dispositions {
		seen[disposition.EntryID]++
		entry, exists := expected[disposition.EntryID]
		if !exists {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.unknown",
				Path:    ".",
				Message: "Decision Document disposition names an unknown root Source Baseline Entry",
			})
			continue
		}
		if seen[disposition.EntryID] > 1 {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.duplicate",
				Path:    entry.Path,
				Message: "root Source Baseline Entry has more than one disposition",
			})
			continue
		}
		if disposition.EntryDigest != entry.Digest {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.stale",
				Path:    entry.Path,
				Message: "root Source Baseline Entry digest does not match current bytes",
			})
			continue
		}
		switch disposition.Disposition {
		case "rejected":
			if strings.TrimSpace(disposition.Reason) == "" || disposition.Destination != nil {
				findings = append(findings, Finding{
					Code:    "baseline.preservation.disposition.invalid",
					Path:    entry.Path,
					Message: "rejected root rule requires an individual reason and null destination",
				})
				continue
			}
		case "managed-entry":
			if finding := validateManagedEntryDisposition(
				entry,
				disposition,
				semanticOwners,
			); finding != nil {
				findings = append(findings, *finding)
				continue
			}
		case "repository-document":
			if _, finding := semanticOwnerForDisposition(
				entry,
				disposition,
				semanticOwners,
			); finding != nil {
				findings = append(findings, *finding)
				continue
			}
		case "repository-rules":
			if finding := validateRepositoryRulesDisposition(disposition); finding != nil {
				finding.Path = entry.Path
				findings = append(findings, *finding)
				continue
			}
			proposed, _ := base64.StdEncoding.DecodeString(disposition.Destination.ProposedBytes)
			if !bytes.Equal(proposed, entry.SourceBytes) {
				findings = append(findings, Finding{
					Code:    "baseline.preservation.repository-rules.invalid",
					Path:    entry.Path,
					Message: "Repository-Specific Normative Rules proposed bytes do not match the current source entry",
				})
				continue
			}
		default:
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.invalid",
				Path:    entry.Path,
				Message: "root instruction preservation disposition is unsupported",
			})
			continue
		}
		accepted = append(accepted, disposition)
	}
	for _, entry := range source.Entries {
		if seen[entry.ID] == 0 {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.missing",
				Path:    entry.Path,
				Message: "root Source Baseline Entry has no disposition",
			})
		}
	}
	if len(findings) != 0 {
		return accepted, nil, nil, findings
	}
	sort.Slice(accepted, func(i, j int) bool {
		return order[accepted[i].EntryID] < order[accepted[j].EntryID]
	})
	var repositoryBlocks []RepositoryRuleBlock
	var repositoryRules []byte
	for _, disposition := range accepted {
		entry := expected[disposition.EntryID]
		switch disposition.Disposition {
		case "repository-document":
			owner, _ := semanticOwnerForDisposition(entry, disposition, semanticOwners)
			repositoryBlocks = append(repositoryBlocks, RepositoryRuleBlock{
				ID:             repositoryRuleStableID(entry, owner),
				SourceEntryID:  entry.ID,
				Classification: disposition.Classification,
				ManagedID:      owner.ManagedID,
				Path:           owner.Path,
				Body:           append([]byte(nil), entry.SourceBytes...),
			})
		case "repository-rules":
			proposed, _ := base64.StdEncoding.DecodeString(disposition.Destination.ProposedBytes)
			repositoryRules = append(repositoryRules, proposed...)
		}
	}
	if finding := validateRepositoryRulesTarget(rootPath, repositoryRules); finding != nil {
		return accepted, repositoryBlocks, nil, append(findings, *finding)
	}
	return accepted, repositoryBlocks, repositoryRules, nil
}

func validateManagedEntryDisposition(
	entry ReadoptionSourceEntry,
	disposition ReadoptionDisposition,
	semanticOwners SemanticOwnerRegistry,
) *Finding {
	destination := disposition.Destination
	_, active := semanticOwners[destinationManagedID(destination)]
	if destination == nil || destination.ManagedID == "" || !active ||
		disposition.Classification == "non-governed" ||
		strings.TrimSpace(disposition.Reason) == "" {
		return &Finding{
			Code:    "baseline.preservation.managed-entry.invalid",
			Path:    entry.Path,
			Message: "managed-entry requires one active managed semantic owner and a review reason",
		}
	}
	return nil
}

func semanticOwnerForDisposition(
	entry ReadoptionSourceEntry,
	disposition ReadoptionDisposition,
	semanticOwners SemanticOwnerRegistry,
) (SemanticOwner, *Finding) {
	destination := disposition.Destination
	if destination == nil ||
		destination.ManagedID != "" ||
		destination.DocumentType != "agent-guide" ||
		destination.Digest != entry.Digest ||
		disposition.Classification == "non-governed" ||
		entry.Kind == "managed-block" ||
		len(bytes.TrimSpace(entry.SourceBytes)) == 0 ||
		strings.TrimSpace(disposition.Reason) == "" {
		return SemanticOwner{}, &Finding{
			Code:    "baseline.preservation.repository-document.invalid",
			Path:    entry.Path,
			Message: "repository-document requires exact non-managed source bytes and an active semantic guide",
		}
	}
	var matched []SemanticOwner
	for _, owner := range semanticOwners {
		if owner.Path == destination.Path {
			matched = append(matched, owner)
		}
	}
	if len(matched) == 0 {
		return SemanticOwner{}, &Finding{
			Code:    "baseline.preservation.repository-document.inactive",
			Path:    destination.Path,
			Message: "repository-document destination is not one active semantic owner",
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ManagedID < matched[j].ManagedID
	})
	return matched[0], nil
}

func destinationManagedID(destination *ReadoptionDestination) string {
	if destination == nil {
		return ""
	}
	return destination.ManagedID
}

func repositoryRuleStableID(entry ReadoptionSourceEntry, owner SemanticOwner) string {
	identity := sha256.New()
	for _, value := range []string{
		entry.ID,
		strconv.Itoa(entry.Start),
		strconv.Itoa(entry.End),
		owner.Path,
	} {
		writeLengthPrefixed(identity, value)
	}
	return "rule." + hex.EncodeToString(identity.Sum(nil))
}

func validateRepositoryRulesDisposition(disposition ReadoptionDisposition) *Finding {
	destination := disposition.Destination
	if destination == nil ||
		destination.DocumentType != "repository-rules" ||
		destination.Path != repositoryRulesPath ||
		!isRawSHA256(destination.Digest) {
		return &Finding{
			Code:    "baseline.preservation.repository-rules.invalid",
			Message: "Repository-Specific Normative Rules require the canonical typed destination",
		}
	}
	proposed, err := base64.StdEncoding.DecodeString(destination.ProposedBytes)
	if err != nil || base64.StdEncoding.EncodeToString(proposed) != destination.ProposedBytes {
		return &Finding{
			Code:    "baseline.preservation.repository-rules.invalid",
			Message: "Repository-Specific Normative Rules proposedBytes must be canonical base64",
		}
	}
	sum := sha256.Sum256(proposed)
	if len(proposed) == 0 || destination.Digest != hex.EncodeToString(sum[:]) ||
		bytes.Contains(proposed, []byte("<!-- setup-context-driven:")) {
		return &Finding{
			Code:    "baseline.preservation.repository-rules.invalid",
			Message: "Repository-Specific Normative Rules proposed bytes are stale, empty, or managed",
		}
	}
	return nil
}

func validateRepositoryRulesTarget(rootPath string, proposed []byte) *Finding {
	if len(proposed) == 0 {
		return nil
	}
	anchored, err := os.OpenRoot(rootPath)
	if err != nil {
		return &Finding{
			Code:    "baseline.preservation.repository-rules.unreadable",
			Path:    repositoryRulesPath,
			Message: "repository root cannot be opened for Repository-Specific Normative Rules validation",
		}
	}
	defer anchored.Close()
	info, err := anchored.Lstat(filepath.FromSlash(repositoryRulesPath))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil || !info.Mode().IsRegular() {
		return &Finding{
			Code:    "baseline.preservation.repository-rules.invalid",
			Path:    repositoryRulesPath,
			Message: "Repository-Specific Normative Rules must be an absent or regular non-symlink file",
		}
	}
	return nil
}

func hasExactRequiredFields(doc document, allowed map[string]struct{}, required ...string) bool {
	for field := range doc {
		if _, ok := allowed[field]; !ok {
			return false
		}
	}
	for _, field := range required {
		if _, ok := doc[field]; !ok {
			return false
		}
	}
	return true
}

func hasExactFields(doc document, fields ...string) bool {
	if len(doc) != len(fields) {
		return false
	}
	for _, field := range fields {
		if _, ok := doc[field]; !ok {
			return false
		}
	}
	return true
}

func decisionDocumentError(code, sourcePath, itemID, message string) error {
	return &DecisionDocumentError{Diagnostics: []DecisionDocumentDiagnostic{{
		Code:    code,
		Path:    sourcePath,
		ItemID:  itemID,
		Message: message,
	}}}
}

func validBaselineIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for index, character := range value {
		valid := character >= 'a' && character <= 'z' ||
			character >= '0' && character <= '9' ||
			index > 0 && (character == '.' || character == '_' || character == '-')
		if !valid {
			return false
		}
	}
	last := value[len(value)-1]
	return last >= 'a' && last <= 'z' || last >= '0' && last <= '9'
}

func isRawSHA256(value string) bool {
	return len(value) == sha256.Size*2 && isLowerHex(value)
}

func writeLengthPrefixed(writer io.Writer, value string) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = writer.Write(length[:])
	_, _ = io.WriteString(writer, value)
}

func cloneMap(value map[string]any) map[string]any {
	result := make(map[string]any, len(value))
	for key, item := range value {
		result[key] = item
	}
	return result
}
