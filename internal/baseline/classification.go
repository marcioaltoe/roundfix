package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	AnalysisSnapshotSchemaVersion       = "roundfix/baseline-classification-snapshot/v1"
	ClassificationProposalSchemaVersion = "roundfix/baseline-classification-proposal/v1"
	AnalysisSnapshotMaxEntries          = 256
	AnalysisSnapshotMaxBytes            = 2 << 20
	ClassificationProposalMaxBytes      = 512 << 10

	analysisSnapshotDigestDomain = AnalysisSnapshotSchemaVersion + "\x00"
)

// ClassificationSource identifies the exact Source Baseline presented for
// semantic classification.
type ClassificationSource struct {
	ID     string `json:"id"`
	Digest string `json:"digest"`
}

// ClassificationDestination is one active disposition the sealed analyzer
// may propose for an exact source entry.
type ClassificationDestination struct {
	Disposition     string   `json:"disposition"`
	ManagedID       string   `json:"managedId,omitempty"`
	DocumentType    string   `json:"documentType,omitempty"`
	Path            string   `json:"path,omitempty"`
	Classifications []string `json:"classifications,omitempty"`
}

// ClassificationProposalContract makes the strict response shape explicit
// inside the sealed payload, so no separate prompt or checkout context is
// required.
type ClassificationProposalContract struct {
	RequiredFields            []string `json:"requiredFields"`
	RequiredDispositionFields []string `json:"requiredDispositionFields"`
	Classifications           []string `json:"classifications"`
	Rules                     []string `json:"rules"`
}

// ClassificationSemanticEntry exposes the exact readable source text used for
// semantic analysis without changing the byte evidence that planning admits.
type ClassificationSemanticEntry struct {
	EntryID  string `json:"entryId"`
	Text     string `json:"text"`
	Readable bool   `json:"readable"`
}

// AnalysisSnapshot is the canonical, checkout-free input to one sealed
// semantic classification attempt.
type AnalysisSnapshot struct {
	SchemaVersion    string                         `json:"schemaVersion"`
	Operation        string                         `json:"operation"`
	ProposalSchema   string                         `json:"proposalSchema"`
	ProposalContract ClassificationProposalContract `json:"proposalContract"`
	SourceBaseline   ClassificationSource           `json:"sourceBaseline"`
	Entries          []ReadoptionSourceEntry        `json:"entries"`
	SemanticEntries  []ClassificationSemanticEntry  `json:"semanticEntries"`
	Destinations     []ClassificationDestination    `json:"destinations"`
	SnapshotDigest   string                         `json:"snapshotDigest"`
}

// ClassificationProposal is a complete proposal for every entry in one
// Analysis Snapshot. It becomes authoritative only after deterministic
// validation and human review.
type ClassificationProposal struct {
	SchemaVersion  string                  `json:"schemaVersion"`
	SnapshotDigest string                  `json:"snapshotDigest"`
	Dispositions   []ReadoptionDisposition `json:"dispositions"`
}

// NewAnalysisSnapshot builds the only semantic payload admitted to the ACP
// boundary. It contains source bytes and identities, but no checkout path.
func NewAnalysisSnapshot(
	source ReadoptionSourceBaseline,
	semanticOwners ...SemanticOwnerRegistry,
) (AnalysisSnapshot, error) {
	if source.EntryCount != len(source.Entries) {
		return AnalysisSnapshot{}, fmt.Errorf(
			"build Analysis Snapshot: Source Baseline declares %d entries but contains %d",
			source.EntryCount,
			len(source.Entries),
		)
	}
	if len(source.Entries) > AnalysisSnapshotMaxEntries {
		return AnalysisSnapshot{}, fmt.Errorf(
			"build Analysis Snapshot: source has %d entries; maximum is %d entries",
			len(source.Entries),
			AnalysisSnapshotMaxEntries,
		)
	}
	entryBytes := 0
	for _, entry := range source.Entries {
		entryBytes += len(entry.SourceBytes)
	}
	if source.ByteCount != entryBytes {
		return AnalysisSnapshot{}, fmt.Errorf(
			"build Analysis Snapshot: Source Baseline declares %d bytes but entries contain %d",
			source.ByteCount,
			entryBytes,
		)
	}
	if !strings.HasPrefix(source.ID, "baseline.readoption.") ||
		strings.TrimPrefix(source.ID, "baseline.readoption.") != source.Digest {
		return AnalysisSnapshot{}, errors.New(
			"build Analysis Snapshot: Source Baseline id and digest do not match",
		)
	}
	var owners SemanticOwnerRegistry
	if len(semanticOwners) > 1 {
		return AnalysisSnapshot{}, errors.New(
			"build Analysis Snapshot: at most one semantic owner registry is allowed",
		)
	}
	if len(semanticOwners) == 1 {
		owners = semanticOwners[0]
	}
	snapshot := AnalysisSnapshot{
		SchemaVersion:    AnalysisSnapshotSchemaVersion,
		Operation:        "classify-root-instructions",
		ProposalSchema:   ClassificationProposalSchemaVersion,
		ProposalContract: classificationProposalContract(),
		SourceBaseline:   ClassificationSource{ID: source.ID, Digest: source.Digest},
		Entries:          cloneClassificationEntries(source.Entries),
		SemanticEntries:  classificationSemanticEntries(source.Entries),
		Destinations:     supportedClassificationDestinations(owners),
	}
	digest, err := computeAnalysisSnapshotDigest(snapshot)
	if err != nil {
		return AnalysisSnapshot{}, err
	}
	snapshot.SnapshotDigest = digest
	if err := validateAnalysisSnapshot(snapshot); err != nil {
		return AnalysisSnapshot{}, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return AnalysisSnapshot{}, fmt.Errorf("serialize Analysis Snapshot: %w", err)
	}
	if len(data) > AnalysisSnapshotMaxBytes {
		return AnalysisSnapshot{}, fmt.Errorf(
			"build Analysis Snapshot: canonical input is %d bytes; maximum is %d bytes",
			len(data),
			AnalysisSnapshotMaxBytes,
		)
	}
	return snapshot, nil
}

// CanonicalBytes returns the exact bytes reused by preferred and fallback
// attempts.
func (snapshot AnalysisSnapshot) CanonicalBytes() ([]byte, error) {
	if err := validateAnalysisSnapshot(snapshot); err != nil {
		return nil, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("serialize Analysis Snapshot: %w", err)
	}
	if len(data) > AnalysisSnapshotMaxBytes {
		return nil, fmt.Errorf(
			"serialize Analysis Snapshot: canonical input is %d bytes; maximum is %d bytes",
			len(data),
			AnalysisSnapshotMaxBytes,
		)
	}
	return data, nil
}

// ParseClassificationProposal rejects all output except one strict, complete
// proposal bound to the supplied Analysis Snapshot.
func ParseClassificationProposal(
	data []byte,
	snapshot AnalysisSnapshot,
) (ClassificationProposal, error) {
	if len(data) > ClassificationProposalMaxBytes {
		return ClassificationProposal{}, fmt.Errorf(
			"parse Classification Proposal: output is %d bytes; maximum is %d bytes",
			len(data),
			ClassificationProposalMaxBytes,
		)
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return ClassificationProposal{}, fmt.Errorf("parse Classification Proposal: %w", err)
	}
	var proposal ClassificationProposal
	if err := decodeStrictJSON(data, &proposal); err != nil {
		return ClassificationProposal{}, fmt.Errorf("parse Classification Proposal: %w", err)
	}
	return normalizeClassificationProposal(snapshot, proposal)
}

// ManualClassificationProposal produces a complete editable destination for
// every entry when sealed analysis cannot be admitted.
func ManualClassificationProposal(snapshot AnalysisSnapshot) (ClassificationProposal, error) {
	if err := validateAnalysisSnapshot(snapshot); err != nil {
		return ClassificationProposal{}, err
	}
	source := ReadoptionSourceBaseline{
		ID:         snapshot.SourceBaseline.ID,
		Digest:     snapshot.SourceBaseline.Digest,
		EntryCount: len(snapshot.Entries),
		Entries:    cloneClassificationEntries(snapshot.Entries),
	}
	skeleton := buildDecisionSkeleton(source)
	proposal := ClassificationProposal{
		SchemaVersion:  ClassificationProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Dispositions:   skeleton.Document.Readoption.Dispositions,
	}
	return normalizeClassificationProposal(snapshot, proposal)
}

// DecisionDocumentFromClassificationProposal converts an admitted proposal
// into the same strict deterministic input used by manual classification.
func DecisionDocumentFromClassificationProposal(
	snapshot AnalysisSnapshot,
	proposal ClassificationProposal,
) (DecisionDocument, error) {
	normalized, err := normalizeClassificationProposal(snapshot, proposal)
	if err != nil {
		return DecisionDocument{}, err
	}
	document := DecisionDocument{
		SchemaVersion: DecisionDocumentSchemaVersion,
		Version:       DecisionDocumentVersion,
		Decisions:     []DecisionValue{},
		Readoption:    &ReadoptionDecisions{Dispositions: normalized.Dispositions},
	}
	document.Readoption.SourceBaseline.ID = snapshot.SourceBaseline.ID
	document.Readoption.SourceBaseline.Digest = snapshot.SourceBaseline.Digest
	data, err := json.Marshal(document)
	if err != nil {
		return DecisionDocument{}, fmt.Errorf("serialize normalized Decision Document: %w", err)
	}
	parsed, err := ParseDecisionDocument(data, "sealed-classification-proposal")
	if err != nil {
		return DecisionDocument{}, fmt.Errorf("validate normalized Decision Document: %w", err)
	}
	return parsed, nil
}

func normalizeClassificationProposal(
	snapshot AnalysisSnapshot,
	proposal ClassificationProposal,
) (ClassificationProposal, error) {
	if err := validateAnalysisSnapshot(snapshot); err != nil {
		return ClassificationProposal{}, err
	}
	if proposal.SchemaVersion != ClassificationProposalSchemaVersion {
		return ClassificationProposal{}, fmt.Errorf(
			"validate Classification Proposal: unsupported schema %q",
			proposal.SchemaVersion,
		)
	}
	if proposal.SnapshotDigest != snapshot.SnapshotDigest {
		return ClassificationProposal{}, errors.New(
			"validate Classification Proposal: Analysis Snapshot digest does not match",
		)
	}
	if proposal.Dispositions == nil {
		return ClassificationProposal{}, errors.New(
			"validate Classification Proposal: dispositions must be a JSON array",
		)
	}

	expected := make(map[string]ReadoptionSourceEntry, len(snapshot.Entries))
	order := make(map[string]int, len(snapshot.Entries))
	for index, entry := range snapshot.Entries {
		expected[entry.ID] = entry
		order[entry.ID] = index
	}
	seen := make(map[string]struct{}, len(proposal.Dispositions))
	normalized := make([]ReadoptionDisposition, 0, len(proposal.Dispositions))
	for _, disposition := range proposal.Dispositions {
		entry, exists := expected[disposition.EntryID]
		if !exists {
			return ClassificationProposal{}, fmt.Errorf(
				"validate Classification Proposal: unknown Source Baseline Entry %q",
				disposition.EntryID,
			)
		}
		if _, duplicate := seen[disposition.EntryID]; duplicate {
			return ClassificationProposal{}, fmt.Errorf(
				"validate Classification Proposal: duplicate disposition for %q",
				disposition.EntryID,
			)
		}
		seen[disposition.EntryID] = struct{}{}
		normalizedDisposition := cloneReadoptionDisposition(disposition)
		deriveClassificationByteEvidence(entry, &normalizedDisposition)
		if err := validateClassificationDisposition(
			snapshot,
			entry,
			normalizedDisposition,
		); err != nil {
			return ClassificationProposal{}, err
		}
		normalized = append(normalized, normalizedDisposition)
	}
	for _, entry := range snapshot.Entries {
		if _, exists := seen[entry.ID]; !exists {
			return ClassificationProposal{}, fmt.Errorf(
				"validate Classification Proposal: missing disposition for %q",
				entry.ID,
			)
		}
	}
	sort.Slice(normalized, func(i, j int) bool {
		return order[normalized[i].EntryID] < order[normalized[j].EntryID]
	})
	return ClassificationProposal{
		SchemaVersion:  ClassificationProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Dispositions:   normalized,
	}, nil
}

func deriveClassificationByteEvidence(
	entry ReadoptionSourceEntry,
	disposition *ReadoptionDisposition,
) {
	if disposition.EntryDigest == "" {
		disposition.EntryDigest = entry.Digest
	}
	if disposition.Destination == nil {
		return
	}
	switch disposition.Disposition {
	case "repository-rules":
		if disposition.Destination.Digest == "" {
			disposition.Destination.Digest = entry.Digest
		}
		if disposition.Destination.ProposedBytes == "" {
			disposition.Destination.ProposedBytes =
				base64.StdEncoding.EncodeToString(entry.SourceBytes)
		}
	case "repository-document":
		if disposition.Destination.Digest == "" {
			disposition.Destination.Digest = entry.Digest
		}
	}
}

func validateClassificationDisposition(
	snapshot AnalysisSnapshot,
	entry ReadoptionSourceEntry,
	disposition ReadoptionDisposition,
) error {
	if disposition.EntryDigest != entry.Digest {
		return fmt.Errorf(
			"validate Classification Proposal: entry digest for %q does not match",
			entry.ID,
		)
	}
	if !containsString(
		[]string{"non-governed", "normative-clause", "operational-contract", "recommendation"},
		disposition.Classification,
	) {
		return fmt.Errorf(
			"validate Classification Proposal: unsupported classification %q for %q",
			disposition.Classification,
			entry.ID,
		)
	}
	switch disposition.Disposition {
	case "rejected":
		if disposition.Destination != nil || strings.TrimSpace(disposition.Reason) == "" {
			return fmt.Errorf(
				"validate Classification Proposal: rejected entry %q requires a reason and null destination",
				entry.ID,
			)
		}
	case "repository-rules":
		if disposition.Classification == "non-governed" ||
			entry.Kind == "managed-block" ||
			len(bytes.TrimSpace(entry.SourceBytes)) == 0 ||
			strings.TrimSpace(disposition.Reason) == "" {
			return fmt.Errorf(
				"validate Classification Proposal: entry %q cannot enter Repository-Specific Normative Rules without a reason",
				entry.ID,
			)
		}
		destination := disposition.Destination
		if destination == nil ||
			destination.ManagedID != "" ||
			destination.DocumentType != "repository-rules" ||
			destination.Path != repositoryRulesPath ||
			destination.Digest != entry.Digest {
			return fmt.Errorf(
				"validate Classification Proposal: unsupported destination for %q",
				entry.ID,
			)
		}
		proposed, err := base64.StdEncoding.DecodeString(destination.ProposedBytes)
		if err != nil ||
			base64.StdEncoding.EncodeToString(proposed) != destination.ProposedBytes ||
			!bytes.Equal(proposed, entry.SourceBytes) {
			return fmt.Errorf(
				"validate Classification Proposal: proposed bytes for %q do not match the sealed source",
				entry.ID,
			)
		}
	case "managed-entry":
		destination := disposition.Destination
		if disposition.Classification == "non-governed" ||
			destination == nil ||
			destination.ManagedID == "" ||
			strings.TrimSpace(disposition.Reason) == "" ||
			!classificationDestinationAllowed(
				snapshot.Destinations,
				disposition.Classification,
				ClassificationDestination{
					Disposition: "managed-entry",
					ManagedID:   destination.ManagedID,
				},
			) {
			return fmt.Errorf(
				"validate Classification Proposal: unsupported managed destination for %q",
				entry.ID,
			)
		}
	case "repository-document":
		destination := disposition.Destination
		if disposition.Classification == "non-governed" ||
			entry.Kind == "managed-block" ||
			len(bytes.TrimSpace(entry.SourceBytes)) == 0 ||
			strings.TrimSpace(disposition.Reason) == "" ||
			destination == nil ||
			destination.ManagedID != "" ||
			destination.DocumentType != "agent-guide" ||
			destination.Digest != entry.Digest ||
			!classificationDestinationAllowed(
				snapshot.Destinations,
				disposition.Classification,
				ClassificationDestination{
					Disposition:  "repository-document",
					DocumentType: destination.DocumentType,
					Path:         destination.Path,
				},
			) {
			return fmt.Errorf(
				"validate Classification Proposal: unsupported semantic destination for %q",
				entry.ID,
			)
		}
	default:
		return fmt.Errorf(
			"validate Classification Proposal: unsupported disposition %q for %q",
			disposition.Disposition,
			entry.ID,
		)
	}
	if disposition.Classification == "non-governed" && disposition.Disposition != "rejected" {
		return fmt.Errorf(
			"validate Classification Proposal: non-governed entry %q must be rejected",
			entry.ID,
		)
	}
	return nil
}

func validateAnalysisSnapshot(snapshot AnalysisSnapshot) error {
	if snapshot.SchemaVersion != AnalysisSnapshotSchemaVersion ||
		snapshot.Operation != "classify-root-instructions" ||
		snapshot.ProposalSchema != ClassificationProposalSchemaVersion {
		return errors.New("validate Analysis Snapshot: schema or operation is invalid")
	}
	if !strings.HasPrefix(snapshot.SourceBaseline.ID, "baseline.readoption.") ||
		strings.TrimPrefix(snapshot.SourceBaseline.ID, "baseline.readoption.") != snapshot.SourceBaseline.Digest ||
		!isRawSHA256(snapshot.SourceBaseline.Digest) {
		return errors.New("validate Analysis Snapshot: Source Baseline identity is invalid")
	}
	if !equalClassificationProposalContract(
		snapshot.ProposalContract,
		classificationProposalContract(),
	) {
		return errors.New("validate Analysis Snapshot: proposal contract changed")
	}
	if snapshot.Entries == nil || len(snapshot.Entries) > AnalysisSnapshotMaxEntries {
		return fmt.Errorf(
			"validate Analysis Snapshot: entries must be an array with at most %d entries",
			AnalysisSnapshotMaxEntries,
		)
	}
	if !equalClassificationSemanticEntries(
		snapshot.SemanticEntries,
		classificationSemanticEntries(snapshot.Entries),
	) {
		return errors.New("validate Analysis Snapshot: semantic entries changed")
	}
	if err := validateClassificationDestinations(snapshot.Destinations); err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(snapshot.Entries))
	for _, entry := range snapshot.Entries {
		if _, duplicate := seen[entry.ID]; duplicate {
			return fmt.Errorf("validate Analysis Snapshot: duplicate entry %q", entry.ID)
		}
		seen[entry.ID] = struct{}{}
		if !strings.HasPrefix(entry.ID, "source-entry.") ||
			!isRawSHA256(strings.TrimPrefix(entry.ID, "source-entry.")) ||
			!repositoryPathIsSafe(entry.Path) ||
			entry.Carrier != entry.Path ||
			entry.Start < 0 ||
			entry.End < entry.Start ||
			entry.End-entry.Start != len(entry.SourceBytes) ||
			!isRawSHA256(entry.Digest) ||
			!isRawSHA256(entry.CarrierDigest) ||
			entry.Encoding != "base64" ||
			entry.StructuralProvenance == nil {
			return fmt.Errorf("validate Analysis Snapshot: entry %q is invalid", entry.ID)
		}
		sum := sha256.Sum256(entry.SourceBytes)
		if hex.EncodeToString(sum[:]) != entry.Digest {
			return fmt.Errorf(
				"validate Analysis Snapshot: entry %q bytes do not match their digest",
				entry.ID,
			)
		}
		rebuilt := newReadoptionSourceEntry(
			entry.Path,
			entry.Kind,
			entry.Start,
			entry.End,
			entry.SourceBytes,
			entry.CarrierDigest,
			entry.StructuralProvenance,
		)
		if rebuilt.ID != entry.ID {
			return fmt.Errorf(
				"validate Analysis Snapshot: entry %q identity does not match its sealed evidence",
				entry.ID,
			)
		}
	}
	digest, err := computeAnalysisSnapshotDigest(snapshot)
	if err != nil {
		return err
	}
	if snapshot.SnapshotDigest != digest {
		return errors.New("validate Analysis Snapshot: snapshot digest does not match")
	}
	return nil
}

func computeAnalysisSnapshotDigest(snapshot AnalysisSnapshot) (string, error) {
	payload := snapshot
	payload.SnapshotDigest = ""
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("serialize Analysis Snapshot digest payload: %w", err)
	}
	sum := sha256.Sum256(append([]byte(analysisSnapshotDigestDomain), data...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func supportedClassificationDestinations(
	semanticOwners SemanticOwnerRegistry,
) []ClassificationDestination {
	destinations := []ClassificationDestination{
		{Disposition: "rejected"},
		{
			Disposition:  "repository-rules",
			DocumentType: "repository-rules",
			Path:         repositoryRulesPath,
		},
	}
	documentClassifications := make(map[string]map[string]struct{})
	for _, managedID := range sortedKeys(semanticOwners) {
		owner := semanticOwners[managedID]
		classifications := append([]string(nil), owner.Classifications...)
		sort.Strings(classifications)
		destinations = append(destinations, ClassificationDestination{
			Disposition:     "managed-entry",
			ManagedID:       owner.ManagedID,
			Classifications: classifications,
		})
		if documentClassifications[owner.Path] == nil {
			documentClassifications[owner.Path] = make(map[string]struct{})
		}
		for _, classification := range owner.Classifications {
			documentClassifications[owner.Path][classification] = struct{}{}
		}
	}
	for _, relative := range sortedKeys(documentClassifications) {
		destinations = append(destinations, ClassificationDestination{
			Disposition:     "repository-document",
			DocumentType:    "agent-guide",
			Path:            relative,
			Classifications: sortedKeys(documentClassifications[relative]),
		})
	}
	return destinations
}

func classificationProposalContract() ClassificationProposalContract {
	return ClassificationProposalContract{
		RequiredFields: []string{
			"schemaVersion",
			"snapshotDigest",
			"dispositions",
		},
		RequiredDispositionFields: []string{
			"entryId",
			"classification",
			"disposition",
			"destination",
			"reason",
		},
		Classifications: []string{
			"non-governed",
			"normative-clause",
			"operational-contract",
			"recommendation",
		},
		Rules: []string{
			"Return exactly one JSON object with no prose.",
			"Return exactly one disposition for every entryId.",
			"Use semanticEntries text for meaning; entries sourceBytes remain the authoritative byte evidence.",
			"Treat readable=false as unavailable semantic text and preserve it conservatively.",
			"Do not return entryDigest, destination digest, or proposedBytes; the planner derives exact byte evidence locally.",
			"Use rejected with a null destination and non-empty reason for non-governed, managed, or empty evidence.",
			"Use repository-rules only with the advertised documentType and path and a non-empty reason.",
			"Use repository-document only with one advertised active agent-guide path and a non-empty reason.",
			"Use managed-entry only with one advertised active managedId and a non-empty reason.",
		},
	}
}

func classificationSemanticEntries(
	entries []ReadoptionSourceEntry,
) []ClassificationSemanticEntry {
	result := make([]ClassificationSemanticEntry, len(entries))
	for index, entry := range entries {
		result[index] = ClassificationSemanticEntry{
			EntryID:  entry.ID,
			Readable: utf8.Valid(entry.SourceBytes),
		}
		if result[index].Readable {
			result[index].Text = string(entry.SourceBytes)
		}
	}
	return result
}

func equalClassificationSemanticEntries(
	first []ClassificationSemanticEntry,
	second []ClassificationSemanticEntry,
) bool {
	if first == nil || len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func equalClassificationProposalContract(
	first ClassificationProposalContract,
	second ClassificationProposalContract,
) bool {
	return equalStrings(first.RequiredFields, second.RequiredFields) &&
		equalStrings(first.RequiredDispositionFields, second.RequiredDispositionFields) &&
		equalStrings(first.Classifications, second.Classifications) &&
		equalStrings(first.Rules, second.Rules)
}

func equalStrings(first, second []string) bool {
	if first == nil || len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func equalClassificationDestinations(first, second []ClassificationDestination) bool {
	if first == nil || len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index].Disposition != second[index].Disposition ||
			first[index].ManagedID != second[index].ManagedID ||
			first[index].DocumentType != second[index].DocumentType ||
			first[index].Path != second[index].Path ||
			!equalOptionalStrings(first[index].Classifications, second[index].Classifications) {
			return false
		}
	}
	return true
}

func equalOptionalStrings(first, second []string) bool {
	if (first == nil) != (second == nil) || len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func validateClassificationDestinations(destinations []ClassificationDestination) error {
	if destinations == nil || len(destinations) < 2 {
		return errors.New("validate Analysis Snapshot: supported destinations changed")
	}
	fixed := supportedClassificationDestinations(nil)
	if !equalClassificationDestinations(destinations[:2], fixed) {
		return errors.New("validate Analysis Snapshot: supported destinations changed")
	}
	seenManaged := make(map[string]struct{})
	seenPaths := make(map[string]struct{})
	for index := 2; index < len(destinations); index++ {
		destination := destinations[index]
		if destination.Classifications == nil {
			return errors.New("validate Analysis Snapshot: semantic destination classifications are missing")
		}
		switch destination.Disposition {
		case "managed-entry":
			if destination.ManagedID == "" ||
				destination.DocumentType != "" ||
				destination.Path != "" {
				return errors.New("validate Analysis Snapshot: managed semantic destination is invalid")
			}
			if _, duplicate := seenManaged[destination.ManagedID]; duplicate {
				return errors.New("validate Analysis Snapshot: semantic managed destination is duplicated")
			}
			seenManaged[destination.ManagedID] = struct{}{}
		case "repository-document":
			if destination.ManagedID != "" ||
				destination.DocumentType != "agent-guide" ||
				!repositoryPathIsSafe(destination.Path) {
				return errors.New("validate Analysis Snapshot: repository semantic destination is invalid")
			}
			if _, duplicate := seenPaths[destination.Path]; duplicate {
				return errors.New("validate Analysis Snapshot: semantic repository destination is duplicated")
			}
			seenPaths[destination.Path] = struct{}{}
		default:
			return errors.New("validate Analysis Snapshot: unsupported semantic destination")
		}
	}
	return nil
}

func classificationDestinationAllowed(
	destinations []ClassificationDestination,
	classification string,
	candidate ClassificationDestination,
) bool {
	for _, destination := range destinations {
		if destination.Disposition != candidate.Disposition ||
			destination.ManagedID != candidate.ManagedID ||
			destination.DocumentType != candidate.DocumentType ||
			destination.Path != candidate.Path {
			continue
		}
		// The semantic classifications describe the concern owned by the guide.
		// Readoption's classification describes enforcement shape, so presence
		// in an active advertised destination is the binding admission check.
		return classification != "non-governed"
	}
	return false
}

func cloneClassificationEntries(entries []ReadoptionSourceEntry) []ReadoptionSourceEntry {
	result := make([]ReadoptionSourceEntry, len(entries))
	for index, entry := range entries {
		result[index] = entry
		result[index].SourceBytes = append([]byte(nil), entry.SourceBytes...)
		result[index].StructuralProvenance = cloneMap(entry.StructuralProvenance)
	}
	return result
}

func cloneReadoptionDisposition(disposition ReadoptionDisposition) ReadoptionDisposition {
	cloned := disposition
	if disposition.Destination != nil {
		destination := *disposition.Destination
		cloned.Destination = &destination
	}
	return cloned
}
