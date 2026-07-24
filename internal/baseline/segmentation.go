package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	RuleSegmentationSnapshotSchemaVersion = "roundfix/baseline-segmentation-snapshot/v1"
	RuleSegmentationProposalSchemaVersion = "roundfix/baseline-segmentation-proposal/v1"
	RuleSegmentationSnapshotMaxBytes      = 2 << 20
	RuleSegmentationProposalMaxBytes      = 512 << 10
	RuleSegmentationProposalMaxSegments   = AnalysisSnapshotMaxEntries

	ruleSegmentationSnapshotDigestDomain = RuleSegmentationSnapshotSchemaVersion + "\x00"
)

// RuleSegmentationProposalContract makes the sealed byte-range response shape
// explicit without granting access to a checkout or write capability.
type RuleSegmentationProposalContract struct {
	RequiredFields        []string `json:"requiredFields"`
	RequiredSegmentFields []string `json:"requiredSegmentFields"`
	Rules                 []string `json:"rules"`
}

// RuleSegmentationSnapshot is the canonical, checkout-free source presented
// to one sealed segmentation attempt.
type RuleSegmentationSnapshot struct {
	SchemaVersion    string                           `json:"schemaVersion"`
	Operation        string                           `json:"operation"`
	ProposalSchema   string                           `json:"proposalSchema"`
	ProposalContract RuleSegmentationProposalContract `json:"proposalContract"`
	SourceBaseline   ReadoptionSourceBaseline         `json:"sourceBaseline"`
	SnapshotDigest   string                           `json:"snapshotDigest"`
}

// RuleSegmentProposal identifies one byte range inside one structural Source
// Baseline Entry. It never carries proposed or rewritten content.
type RuleSegmentProposal struct {
	EntryID string `json:"entryId"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Digest  string `json:"digest"`
}

// RuleSegmentationProposal is a complete byte-exhaustive proposal bound to one
// immutable Segmentation Snapshot and Source Baseline identity.
type RuleSegmentationProposal struct {
	SchemaVersion  string                `json:"schemaVersion"`
	SnapshotDigest string                `json:"snapshotDigest"`
	SourceBaseline ClassificationSource  `json:"sourceBaseline"`
	Segments       []RuleSegmentProposal `json:"segments"`
}

// NewRuleSegmentationSnapshot seals one current Source Baseline for semantic
// byte-range analysis.
func NewRuleSegmentationSnapshot(
	source ReadoptionSourceBaseline,
) (RuleSegmentationSnapshot, error) {
	if _, err := NewAnalysisSnapshot(source); err != nil {
		return RuleSegmentationSnapshot{}, fmt.Errorf(
			"build Segmentation Snapshot: invalid Source Baseline: %w",
			err,
		)
	}
	snapshot := RuleSegmentationSnapshot{
		SchemaVersion:    RuleSegmentationSnapshotSchemaVersion,
		Operation:        "segment-source-baseline-entries",
		ProposalSchema:   RuleSegmentationProposalSchemaVersion,
		ProposalContract: ruleSegmentationProposalContract(),
		SourceBaseline:   cloneReadoptionSourceBaseline(source),
	}
	digest, err := computeRuleSegmentationSnapshotDigest(snapshot)
	if err != nil {
		return RuleSegmentationSnapshot{}, err
	}
	snapshot.SnapshotDigest = digest
	if err := validateRuleSegmentationSnapshot(snapshot); err != nil {
		return RuleSegmentationSnapshot{}, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return RuleSegmentationSnapshot{}, fmt.Errorf(
			"serialize Segmentation Snapshot: %w",
			err,
		)
	}
	if len(data) > RuleSegmentationSnapshotMaxBytes {
		return RuleSegmentationSnapshot{}, fmt.Errorf(
			"build Segmentation Snapshot: canonical input is %d bytes; maximum is %d bytes",
			len(data),
			RuleSegmentationSnapshotMaxBytes,
		)
	}
	return snapshot, nil
}

// CanonicalBytes returns the exact sealed bytes reused by preferred and
// fallback attempts.
func (snapshot RuleSegmentationSnapshot) CanonicalBytes() ([]byte, error) {
	if err := validateRuleSegmentationSnapshot(snapshot); err != nil {
		return nil, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("serialize Segmentation Snapshot: %w", err)
	}
	if len(data) > RuleSegmentationSnapshotMaxBytes {
		return nil, fmt.Errorf(
			"serialize Segmentation Snapshot: canonical input is %d bytes; maximum is %d bytes",
			len(data),
			RuleSegmentationSnapshotMaxBytes,
		)
	}
	return data, nil
}

// ParseRuleSegmentationProposal rejects every response except one strict,
// complete proposal for the supplied Segmentation Snapshot.
func ParseRuleSegmentationProposal(
	data []byte,
	snapshot RuleSegmentationSnapshot,
) (RuleSegmentationProposal, error) {
	if len(data) > RuleSegmentationProposalMaxBytes {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"parse Segmentation Proposal: output is %d bytes; maximum is %d bytes",
			len(data),
			RuleSegmentationProposalMaxBytes,
		)
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"parse Segmentation Proposal: %w",
			err,
		)
	}
	var proposal RuleSegmentationProposal
	if err := decodeStrictJSON(data, &proposal); err != nil {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"parse Segmentation Proposal: %w",
			err,
		)
	}
	return normalizeRuleSegmentationProposal(snapshot, proposal)
}

func normalizeRuleSegmentationProposal(
	snapshot RuleSegmentationSnapshot,
	proposal RuleSegmentationProposal,
) (RuleSegmentationProposal, error) {
	if err := validateRuleSegmentationSnapshot(snapshot); err != nil {
		return RuleSegmentationProposal{}, err
	}
	if proposal.SchemaVersion != RuleSegmentationProposalSchemaVersion {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"validate Segmentation Proposal: unsupported schema %q",
			proposal.SchemaVersion,
		)
	}
	if proposal.SnapshotDigest != snapshot.SnapshotDigest {
		return RuleSegmentationProposal{}, errors.New(
			"validate Segmentation Proposal: Segmentation Snapshot digest does not match",
		)
	}
	if proposal.SourceBaseline.ID != snapshot.SourceBaseline.ID ||
		proposal.SourceBaseline.Digest != snapshot.SourceBaseline.Digest {
		return RuleSegmentationProposal{}, errors.New(
			"validate Segmentation Proposal: Source Baseline identity does not match",
		)
	}
	if proposal.Segments == nil {
		return RuleSegmentationProposal{}, errors.New(
			"validate Segmentation Proposal: segments must be a JSON array",
		)
	}
	if len(proposal.Segments) > RuleSegmentationProposalMaxSegments {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"validate Segmentation Proposal: proposal has %d segments; maximum is %d",
			len(proposal.Segments),
			RuleSegmentationProposalMaxSegments,
		)
	}
	materializedEntryCount := len(proposal.Segments)
	for _, entry := range snapshot.SourceBaseline.Entries {
		if len(entry.SourceBytes) == 0 {
			materializedEntryCount++
		}
	}
	if materializedEntryCount > AnalysisSnapshotMaxEntries {
		return RuleSegmentationProposal{}, fmt.Errorf(
			"validate Segmentation Proposal: materialized source has %d entries; classification maximum is %d",
			materializedEntryCount,
			AnalysisSnapshotMaxEntries,
		)
	}

	entries := make(map[string]ReadoptionSourceEntry, len(snapshot.SourceBaseline.Entries))
	for _, entry := range snapshot.SourceBaseline.Entries {
		entries[entry.ID] = entry
	}
	for _, segment := range proposal.Segments {
		if _, exists := entries[segment.EntryID]; !exists {
			return RuleSegmentationProposal{}, fmt.Errorf(
				"validate Segmentation Proposal: unknown Source Baseline Entry %q",
				segment.EntryID,
			)
		}
	}

	segmentIndex := 0
	for _, entry := range snapshot.SourceBaseline.Entries {
		if len(entry.SourceBytes) == 0 {
			if segmentIndex < len(proposal.Segments) &&
				proposal.Segments[segmentIndex].EntryID == entry.ID {
				return RuleSegmentationProposal{}, fmt.Errorf(
					"validate Segmentation Proposal: empty Source Baseline Entry %q cannot have a segment",
					entry.ID,
				)
			}
			continue
		}
		if segmentIndex >= len(proposal.Segments) ||
			proposal.Segments[segmentIndex].EntryID != entry.ID {
			return RuleSegmentationProposal{}, fmt.Errorf(
				"validate Segmentation Proposal: ranges for %q do not cover its first byte",
				entry.ID,
			)
		}
		cursor := 0
		for segmentIndex < len(proposal.Segments) &&
			proposal.Segments[segmentIndex].EntryID == entry.ID {
			segment := proposal.Segments[segmentIndex]
			if segment.Start != cursor {
				return RuleSegmentationProposal{}, fmt.Errorf(
					"validate Segmentation Proposal: range [%d,%d) for %q starts at %d; expected %d",
					segment.Start,
					segment.End,
					entry.ID,
					segment.Start,
					cursor,
				)
			}
			if segment.End <= segment.Start {
				return RuleSegmentationProposal{}, fmt.Errorf(
					"validate Segmentation Proposal: range [%d,%d) for %q must be non-empty",
					segment.Start,
					segment.End,
					entry.ID,
				)
			}
			if segment.Start < 0 || segment.End > len(entry.SourceBytes) {
				return RuleSegmentationProposal{}, fmt.Errorf(
					"validate Segmentation Proposal: range [%d,%d) for %q is outside [0,%d)",
					segment.Start,
					segment.End,
					entry.ID,
					len(entry.SourceBytes),
				)
			}
			if segment.Digest != ruleSegmentDigest(entry.SourceBytes[segment.Start:segment.End]) {
				return RuleSegmentationProposal{}, fmt.Errorf(
					"validate Segmentation Proposal: range digest for %q does not match sealed source bytes",
					entry.ID,
				)
			}
			cursor = segment.End
			segmentIndex++
		}
		if cursor != len(entry.SourceBytes) {
			return RuleSegmentationProposal{}, fmt.Errorf(
				"validate Segmentation Proposal: ranges for %q cover %d of %d bytes",
				entry.ID,
				cursor,
				len(entry.SourceBytes),
			)
		}
	}
	if segmentIndex != len(proposal.Segments) {
		return RuleSegmentationProposal{}, errors.New(
			"validate Segmentation Proposal: ranges are not in Source Baseline order",
		)
	}

	normalized := proposal
	normalized.Segments = append([]RuleSegmentProposal(nil), proposal.Segments...)
	return normalized, nil
}

// ManualRuleSegmentationProposal preserves every non-empty structural entry as
// one lossless range when semantic segmentation is unavailable.
func ManualRuleSegmentationProposal(
	snapshot RuleSegmentationSnapshot,
) (RuleSegmentationProposal, error) {
	if err := validateRuleSegmentationSnapshot(snapshot); err != nil {
		return RuleSegmentationProposal{}, err
	}
	proposal := RuleSegmentationProposal{
		SchemaVersion:  RuleSegmentationProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		SourceBaseline: ClassificationSource{
			ID:     snapshot.SourceBaseline.ID,
			Digest: snapshot.SourceBaseline.Digest,
		},
		Segments: make([]RuleSegmentProposal, 0, len(snapshot.SourceBaseline.Entries)),
	}
	for _, entry := range snapshot.SourceBaseline.Entries {
		if len(entry.SourceBytes) == 0 {
			continue
		}
		proposal.Segments = append(proposal.Segments, RuleSegmentProposal{
			EntryID: entry.ID,
			Start:   0,
			End:     len(entry.SourceBytes),
			Digest:  ruleSegmentDigest(entry.SourceBytes),
		})
	}
	return normalizeRuleSegmentationProposal(snapshot, proposal)
}

// MaterializeRuleSegments derives every admitted segment's bytes, digest, and
// Source Baseline Entry identity locally.
func MaterializeRuleSegments(
	snapshot RuleSegmentationSnapshot,
	proposal RuleSegmentationProposal,
) (ReadoptionSourceBaseline, error) {
	normalized, err := normalizeRuleSegmentationProposal(snapshot, proposal)
	if err != nil {
		return ReadoptionSourceBaseline{}, err
	}
	materialized := cloneReadoptionSourceBaseline(snapshot.SourceBaseline)
	materialized.Entries = make(
		[]ReadoptionSourceEntry,
		0,
		len(normalized.Segments),
	)
	segmentIndex := 0
	for _, entry := range snapshot.SourceBaseline.Entries {
		if len(entry.SourceBytes) == 0 {
			materialized.Entries = append(
				materialized.Entries,
				cloneReadoptionSourceEntry(entry),
			)
			continue
		}
		for segmentIndex < len(normalized.Segments) &&
			normalized.Segments[segmentIndex].EntryID == entry.ID {
			segment := normalized.Segments[segmentIndex]
			sourceBytes := entry.SourceBytes[segment.Start:segment.End]
			materialized.Entries = append(
				materialized.Entries,
				newReadoptionSourceEntry(
					entry.Path,
					entry.Kind,
					entry.Start+segment.Start,
					entry.Start+segment.End,
					sourceBytes,
					entry.CarrierDigest,
					entry.StructuralProvenance,
				),
			)
			segmentIndex++
		}
	}
	materialized.EntryCount = len(materialized.Entries)
	if _, err := NewAnalysisSnapshot(materialized); err != nil {
		return ReadoptionSourceBaseline{}, fmt.Errorf(
			"materialize Segmentation Proposal: classification contract rejected entries: %w",
			err,
		)
	}
	return materialized, nil
}

func validateRuleSegmentationSnapshot(snapshot RuleSegmentationSnapshot) error {
	if snapshot.SchemaVersion != RuleSegmentationSnapshotSchemaVersion ||
		snapshot.Operation != "segment-source-baseline-entries" ||
		snapshot.ProposalSchema != RuleSegmentationProposalSchemaVersion {
		return errors.New("validate Segmentation Snapshot: schema or operation is invalid")
	}
	if !equalRuleSegmentationProposalContract(
		snapshot.ProposalContract,
		ruleSegmentationProposalContract(),
	) {
		return errors.New("validate Segmentation Snapshot: proposal contract changed")
	}
	if _, err := NewAnalysisSnapshot(snapshot.SourceBaseline); err != nil {
		return fmt.Errorf(
			"validate Segmentation Snapshot: invalid Source Baseline: %w",
			err,
		)
	}
	digest, err := computeRuleSegmentationSnapshotDigest(snapshot)
	if err != nil {
		return err
	}
	if snapshot.SnapshotDigest != digest {
		return errors.New(
			"validate Segmentation Snapshot: snapshot digest does not match",
		)
	}
	return nil
}

func computeRuleSegmentationSnapshotDigest(
	snapshot RuleSegmentationSnapshot,
) (string, error) {
	payload := snapshot
	payload.SnapshotDigest = ""
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("serialize Segmentation Snapshot digest payload: %w", err)
	}
	sum := sha256.Sum256(append(
		[]byte(ruleSegmentationSnapshotDigestDomain),
		data...,
	))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func ruleSegmentationProposalContract() RuleSegmentationProposalContract {
	return RuleSegmentationProposalContract{
		RequiredFields: []string{
			"schemaVersion",
			"snapshotDigest",
			"sourceBaseline",
			"segments",
		},
		RequiredSegmentFields: []string{
			"entryId",
			"start",
			"end",
			"digest",
		},
		Rules: []string{
			"Return exactly one JSON object with no prose.",
			"Return only ordered non-empty byte ranges for advertised entryId values.",
			"Cover every byte of every non-empty entry exactly once without gaps, overlap, duplication, or reordering.",
			"Copy the lowercase SHA-256 digest of each proposed range; never return source content.",
		},
	}
}

func equalRuleSegmentationProposalContract(
	first RuleSegmentationProposalContract,
	second RuleSegmentationProposalContract,
) bool {
	return equalStrings(first.RequiredFields, second.RequiredFields) &&
		equalStrings(first.RequiredSegmentFields, second.RequiredSegmentFields) &&
		equalStrings(first.Rules, second.Rules)
}

func cloneReadoptionSourceBaseline(
	source ReadoptionSourceBaseline,
) ReadoptionSourceBaseline {
	cloned := source
	cloned.Entries = make([]ReadoptionSourceEntry, len(source.Entries))
	for index, entry := range source.Entries {
		cloned.Entries[index] = cloneReadoptionSourceEntry(entry)
	}
	return cloned
}

func cloneReadoptionSourceEntry(entry ReadoptionSourceEntry) ReadoptionSourceEntry {
	cloned := entry
	cloned.SourceBytes = append([]byte(nil), entry.SourceBytes...)
	cloned.StructuralProvenance = cloneMap(entry.StructuralProvenance)
	return cloned
}

func ruleSegmentDigest(sourceBytes []byte) string {
	sum := sha256.Sum256(sourceBytes)
	return hex.EncodeToString(sum[:])
}
