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
	repositoryRulesPath           = "docs/agents/repository-rules.md"
)

type PreservationMode string

const (
	PreservationModeGreenfield   PreservationMode = "greenfield"
	PreservationModePreservation PreservationMode = "preservation"
)

type PreservationState string

const (
	PreservationStateReady          PreservationState = "ready"
	PreservationStateActionRequired PreservationState = "action_required"
	PreservationStateBlocked        PreservationState = "blocked"
)

// RootPreservationRequest contains the explicit root-instruction choice and
// any owner-approved Decision Document.
type RootPreservationRequest struct {
	Mode      PreservationMode
	Decisions *DecisionDocument
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

// RootPreservationPlan is the complete read-only preservation decision
// result consumed by later portable-plan assembly.
type RootPreservationPlan struct {
	Mode                 PreservationMode         `json:"mode"`
	State                PreservationState        `json:"state"`
	Backups              []RootBackup             `json:"backups"`
	SourceBaseline       ReadoptionSourceBaseline `json:"sourceBaseline"`
	Dispositions         []ReadoptionDisposition  `json:"dispositions"`
	RepositoryRulesBytes []byte                   `json:"repositoryRulesBytes,omitempty"`
	Warnings             []Finding                `json:"warnings"`
	Findings             []Finding                `json:"findings"`
	DecisionSkeleton     *DecisionSkeleton        `json:"decisionSkeleton,omitempty"`
	NextAction           string                   `json:"nextAction,omitempty"`
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

// PlanRootPreservation resolves only the root-instruction preservation slice.
// It reads the exact sources already selected by RepositoryInspection and
// never mutates the repository.
func PlanRootPreservation(
	inspection RepositoryInspection,
	request RootPreservationRequest,
) (RootPreservationPlan, error) {
	if request.Mode != PreservationModeGreenfield && request.Mode != PreservationModePreservation {
		return RootPreservationPlan{}, fmt.Errorf(
			"plan root-instruction preservation: unsupported mode %q",
			request.Mode,
		)
	}
	plan := RootPreservationPlan{
		Mode:     request.Mode,
		State:    PreservationStateReady,
		Warnings: append([]Finding(nil), inspection.Snapshot.Warnings...),
		Findings: append([]Finding(nil), inspection.Snapshot.Blocking...),
	}
	sources, findings, err := loadRootPreservationSources(inspection)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	plan.Findings = append(plan.Findings, findings...)
	plan.Backups, findings, err = planRootBackups(inspection.Root, sources)
	if err != nil {
		return RootPreservationPlan{}, err
	}
	plan.Findings = append(plan.Findings, findings...)
	plan.SourceBaseline = buildReadoptionSourceBaseline(sources)
	plan.Findings = sortedFindings(plan.Findings)

	if len(plan.Findings) != 0 {
		plan.State = PreservationStateBlocked
		plan.NextAction = "repair every blocking root carrier or backup collision and rerun Baseline planning"
		return plan, nil
	}
	if request.Mode == PreservationModeGreenfield {
		return plan, nil
	}
	if len(plan.SourceBaseline.Entries) == 0 {
		return plan, nil
	}
	if request.Decisions == nil {
		skeleton := buildDecisionSkeleton(plan.SourceBaseline)
		plan.State = PreservationStateActionRequired
		plan.DecisionSkeleton = &skeleton
		plan.NextAction = skeleton.NextAction
		return plan, nil
	}
	dispositions, repositoryRules, decisionFindings := validatePreservationDecisions(
		inspection.Root,
		plan.SourceBaseline,
		*request.Decisions,
	)
	plan.Dispositions = dispositions
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

type rootPreservationSource struct {
	carrierPath     string
	sourcePath      string
	contentIdentity string
	content         []byte
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
	for _, source := range sources {
		writeLengthPrefixed(identity, source.sourcePath)
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(source.content)))
		_, _ = identity.Write(length[:])
		_, _ = identity.Write(source.content)
		byteCount += len(source.content)
		entries = append(entries, partitionRootSource(source.sourcePath, source.content)...)
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
		CarrierCount:     len(sources),
		EntryCount:       len(entries),
		ByteCount:        byteCount,
		Entries:          entries,
	}
}

var managedBeginMarker = regexp.MustCompile(
	`<!--\s*setup-context-driven:begin\s+id=([A-Za-z0-9_.-]+)\s+version=([0-9]+)\s*-->`,
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
		version, _ := strconv.Atoi(versionText)
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
) ([]ReadoptionDisposition, []byte, []Finding) {
	var findings []Finding
	if document.SchemaVersion != DecisionDocumentSchemaVersion ||
		document.Version != DecisionDocumentVersion ||
		document.Readoption == nil {
		findings = append(findings, Finding{
			Code:    "baseline.preservation.decision-document.invalid",
			Path:    ".",
			Message: "preservation requires a complete strict Decision Document with Readoption dispositions",
		})
		return nil, nil, findings
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
		if disposition.Disposition != "repository-rules" && disposition.Disposition != "rejected" {
			findings = append(findings, Finding{
				Code:    "baseline.preservation.disposition.invalid",
				Path:    entry.Path,
				Message: "root instruction preservation accepts only repository-rules or rejected dispositions",
			})
			continue
		}
		if disposition.Disposition == "rejected" {
			if strings.TrimSpace(disposition.Reason) == "" || disposition.Destination != nil {
				findings = append(findings, Finding{
					Code:    "baseline.preservation.disposition.invalid",
					Path:    entry.Path,
					Message: "rejected root rule requires an individual reason and null destination",
				})
				continue
			}
		} else if finding := validateRepositoryRulesDisposition(disposition); finding != nil {
			finding.Path = entry.Path
			findings = append(findings, *finding)
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
		return accepted, nil, findings
	}
	sort.Slice(accepted, func(i, j int) bool {
		return order[accepted[i].EntryID] < order[accepted[j].EntryID]
	})
	var repositoryRules []byte
	for _, disposition := range accepted {
		if disposition.Disposition != "repository-rules" {
			continue
		}
		proposed, _ := base64.StdEncoding.DecodeString(disposition.Destination.ProposedBytes)
		repositoryRules = append(repositoryRules, proposed...)
	}
	if finding := validateRepositoryRulesTarget(rootPath, repositoryRules); finding != nil {
		return accepted, nil, append(findings, *finding)
	}
	return accepted, repositoryRules, nil
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
