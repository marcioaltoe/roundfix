package baseline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	RevisionSnapshotSchemaVersion = "roundfix/baseline-revision-snapshot/v1"
	RevisionProposalSchemaVersion = "roundfix/baseline-revision-proposal/v1"
	RevisionSnapshotMaxBytes      = 2 << 20
	RevisionProposalMaxBytes      = 512 << 10

	revisionSnapshotDigestDomain = RevisionSnapshotSchemaVersion + "\x00"
)

type RevisionArea string

const (
	RevisionAreaProfile         RevisionArea = "profile"
	RevisionAreaRepositoryRules RevisionArea = "repository-rules"
	RevisionAreaDivergences     RevisionArea = "divergences"
	RevisionAreaFiles           RevisionArea = "files"
)

// RevisionSnapshot is the bounded, digest-bound input for one optional
// semantic translation of maintainer feedback.
type RevisionSnapshot struct {
	SchemaVersion      string          `json:"schemaVersion"`
	Operation          string          `json:"operation"`
	ProposalSchema     string          `json:"proposalSchema"`
	Area               RevisionArea    `json:"area"`
	Suggestion         string          `json:"suggestion"`
	PlanDigest         string          `json:"planDigest"`
	Decisions          []DecisionValue `json:"decisions"`
	AllowedDecisionIDs []string        `json:"allowedDecisionIds"`
	Rules              []string        `json:"rules"`
	SnapshotDigest     string          `json:"snapshotDigest"`
}

// RevisionProposal contains only complete replacements for permitted
// Baseline decisions. Manual is produced locally and is never accepted from
// ACP output.
type RevisionProposal struct {
	SchemaVersion  string          `json:"schemaVersion"`
	SnapshotDigest string          `json:"snapshotDigest"`
	Area           RevisionArea    `json:"area"`
	Changes        []DecisionValue `json:"changes"`
	Manual         bool            `json:"manual,omitempty"`
}

func NewRevisionSnapshot(plan PlanDocument, area RevisionArea, suggestion string) (RevisionSnapshot, error) {
	if strings.TrimSpace(plan.PlanDigest) == "" || len(plan.Decisions) == 0 {
		return RevisionSnapshot{}, errors.New("build Revision Snapshot: complete plan digest and decisions are required")
	}
	if !validRevisionArea(area) {
		return RevisionSnapshot{}, fmt.Errorf("build Revision Snapshot: unsupported decision area %q", area)
	}
	suggestion = strings.TrimSpace(suggestion)
	if suggestion == "" {
		return RevisionSnapshot{}, errors.New("build Revision Snapshot: suggestion is required")
	}
	decisions := cloneDecisionValues(plan.Decisions)
	allowed := make([]string, len(decisions))
	for index, decision := range decisions {
		allowed[index] = decision.ID
	}
	sort.Strings(allowed)
	snapshot := RevisionSnapshot{
		SchemaVersion:      RevisionSnapshotSchemaVersion,
		Operation:          "revise-baseline-plan",
		ProposalSchema:     RevisionProposalSchemaVersion,
		Area:               area,
		Suggestion:         suggestion,
		PlanDigest:         plan.PlanDigest,
		Decisions:          decisions,
		AllowedDecisionIDs: allowed,
		Rules: []string{
			"Return JSON only with schemaVersion, snapshotDigest, area, and changes.",
			"Change only allowed decision IDs and include every changed decision exactly once.",
			"Do not propose repository paths, file bytes, commands, policy, or destinations.",
		},
	}
	digest, err := computeRevisionSnapshotDigest(snapshot)
	if err != nil {
		return RevisionSnapshot{}, err
	}
	snapshot.SnapshotDigest = digest
	if err := validateRevisionSnapshot(snapshot); err != nil {
		return RevisionSnapshot{}, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return RevisionSnapshot{}, fmt.Errorf("serialize Revision Snapshot: %w", err)
	}
	if len(data) > RevisionSnapshotMaxBytes {
		return RevisionSnapshot{}, fmt.Errorf("build Revision Snapshot: canonical input is %d bytes; maximum is %d bytes", len(data), RevisionSnapshotMaxBytes)
	}
	return snapshot, nil
}

func (snapshot RevisionSnapshot) CanonicalBytes() ([]byte, error) {
	if err := validateRevisionSnapshot(snapshot); err != nil {
		return nil, err
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("serialize Revision Snapshot: %w", err)
	}
	if len(data) > RevisionSnapshotMaxBytes {
		return nil, fmt.Errorf("serialize Revision Snapshot: canonical input is %d bytes; maximum is %d bytes", len(data), RevisionSnapshotMaxBytes)
	}
	return data, nil
}

func ParseRevisionProposal(data []byte, snapshot RevisionSnapshot) (RevisionProposal, error) {
	if len(data) > RevisionProposalMaxBytes {
		return RevisionProposal{}, fmt.Errorf("parse Revision Proposal: output is %d bytes; maximum is %d bytes", len(data), RevisionProposalMaxBytes)
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return RevisionProposal{}, fmt.Errorf("parse Revision Proposal: %w", err)
	}
	var proposal RevisionProposal
	if err := decodeStrictJSON(data, &proposal); err != nil {
		return RevisionProposal{}, fmt.Errorf("parse Revision Proposal: %w", err)
	}
	if proposal.Manual {
		return RevisionProposal{}, errors.New("validate Revision Proposal: ACP output cannot request manual mode")
	}
	return normalizeRevisionProposal(snapshot, proposal)
}

func ManualRevisionProposal(snapshot RevisionSnapshot) (RevisionProposal, error) {
	if err := validateRevisionSnapshot(snapshot); err != nil {
		return RevisionProposal{}, err
	}
	return RevisionProposal{
		SchemaVersion:  RevisionProposalSchemaVersion,
		SnapshotDigest: snapshot.SnapshotDigest,
		Area:           snapshot.Area,
		Changes:        []DecisionValue{},
		Manual:         true,
	}, nil
}

func DecisionsFromRevisionProposal(snapshot RevisionSnapshot, proposal RevisionProposal) ([]DecisionValue, error) {
	normalized, err := normalizeRevisionProposal(snapshot, proposal)
	if err != nil {
		return nil, err
	}
	if normalized.Manual {
		return nil, errors.New("apply Revision Proposal: manual revision is required")
	}
	decisions := cloneDecisionValues(snapshot.Decisions)
	index := make(map[string]int, len(decisions))
	for position, decision := range decisions {
		index[decision.ID] = position
	}
	for _, change := range normalized.Changes {
		decisions[index[change.ID]].Value = cloneJSONValue(change.Value)
	}
	return decisions, nil
}

// RecalculatePlan rebuilds the complete plan from normalized input. The
// original plan is the immutable preimage boundary; no prior plan fields are
// patched or reused as output.
func RecalculatePlan(ctx context.Context, original PlanDocument, request PlanRequest) (PlanOutcome, error) {
	if err := ValidatePlanDocument(original); err != nil {
		return PlanOutcome{}, fmt.Errorf("recalculate Baseline Plan: validate original plan: %w", err)
	}
	if err := ValidatePlanRepository(ctx, request.Repository, original); err != nil {
		return PlanOutcome{}, fmt.Errorf(
			"recalculate Baseline Plan: immutable repository snapshot changed; restart Baseline adoption: %w",
			err,
		)
	}
	outcome, err := BuildPlan(ctx, request)
	if err != nil || outcome.Plan == nil {
		return outcome, err
	}
	if !reflect.DeepEqual(outcome.Plan.Repository, original.Repository) {
		return PlanOutcome{}, errors.New("recalculate Baseline Plan: repository identity changed after review; restart Baseline adoption")
	}
	originalByPath := preimagesByPath(original.Preimages)
	for _, current := range outcome.Plan.Preimages {
		prior, known := originalByPath[current.Path]
		if known {
			if !reflect.DeepEqual(prior, current) {
				return PlanOutcome{}, fmt.Errorf("recalculate Baseline Plan: immutable preimage %q changed after review; restart Baseline adoption", current.Path)
			}
			continue
		}
		if current.Exists {
			return PlanOutcome{}, fmt.Errorf("recalculate Baseline Plan: revised plan consulted existing path %q outside the original snapshot; restart Baseline adoption", current.Path)
		}
	}
	return outcome, nil
}

func normalizeRevisionProposal(snapshot RevisionSnapshot, proposal RevisionProposal) (RevisionProposal, error) {
	if err := validateRevisionSnapshot(snapshot); err != nil {
		return RevisionProposal{}, err
	}
	if proposal.SchemaVersion != RevisionProposalSchemaVersion {
		return RevisionProposal{}, fmt.Errorf("validate Revision Proposal: unsupported schema %q", proposal.SchemaVersion)
	}
	if proposal.SnapshotDigest != snapshot.SnapshotDigest {
		return RevisionProposal{}, errors.New("validate Revision Proposal: snapshot digest mismatch")
	}
	if proposal.Area != snapshot.Area {
		return RevisionProposal{}, fmt.Errorf("validate Revision Proposal: decision area %q does not match %q", proposal.Area, snapshot.Area)
	}
	if proposal.Manual {
		if len(proposal.Changes) != 0 {
			return RevisionProposal{}, errors.New("validate Revision Proposal: manual revision cannot contain changes")
		}
		return proposal, nil
	}
	if len(proposal.Changes) == 0 {
		return RevisionProposal{}, errors.New("validate Revision Proposal: at least one complete decision change is required")
	}
	allowed := make(map[string]struct{}, len(snapshot.AllowedDecisionIDs))
	for _, id := range snapshot.AllowedDecisionIDs {
		allowed[id] = struct{}{}
	}
	seen := make(map[string]struct{}, len(proposal.Changes))
	changes := cloneDecisionValues(proposal.Changes)
	for _, change := range changes {
		if _, ok := allowed[change.ID]; !ok {
			return RevisionProposal{}, fmt.Errorf("validate Revision Proposal: unknown or out-of-scope decision %q", change.ID)
		}
		if _, duplicate := seen[change.ID]; duplicate {
			return RevisionProposal{}, fmt.Errorf("validate Revision Proposal: duplicate decision %q", change.ID)
		}
		if change.Value == nil {
			return RevisionProposal{}, fmt.Errorf("validate Revision Proposal: decision %q has an incomplete null value", change.ID)
		}
		seen[change.ID] = struct{}{}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].ID < changes[j].ID })
	proposal.Changes = changes
	return proposal, nil
}

func validateRevisionSnapshot(snapshot RevisionSnapshot) error {
	if snapshot.SchemaVersion != RevisionSnapshotSchemaVersion ||
		snapshot.Operation != "revise-baseline-plan" ||
		snapshot.ProposalSchema != RevisionProposalSchemaVersion {
		return errors.New("validate Revision Snapshot: unsupported contract")
	}
	if !validRevisionArea(snapshot.Area) || strings.TrimSpace(snapshot.Suggestion) == "" ||
		strings.TrimSpace(snapshot.PlanDigest) == "" || len(snapshot.Decisions) == 0 ||
		len(snapshot.AllowedDecisionIDs) == 0 || len(snapshot.Rules) == 0 {
		return errors.New("validate Revision Snapshot: incomplete contract")
	}
	decisionIDs := make([]string, len(snapshot.Decisions))
	seen := make(map[string]struct{}, len(snapshot.Decisions))
	for index, decision := range snapshot.Decisions {
		if strings.TrimSpace(decision.ID) == "" || decision.Value == nil {
			return errors.New("validate Revision Snapshot: incomplete decision")
		}
		if _, duplicate := seen[decision.ID]; duplicate {
			return fmt.Errorf("validate Revision Snapshot: duplicate decision %q", decision.ID)
		}
		seen[decision.ID] = struct{}{}
		decisionIDs[index] = decision.ID
	}
	sort.Strings(decisionIDs)
	if !reflect.DeepEqual(decisionIDs, snapshot.AllowedDecisionIDs) {
		return errors.New("validate Revision Snapshot: allowed decisions do not match the normalized plan decisions")
	}
	expected, err := computeRevisionSnapshotDigest(snapshot)
	if err != nil {
		return err
	}
	if snapshot.SnapshotDigest != expected {
		return errors.New("validate Revision Snapshot: snapshot digest mismatch")
	}
	return nil
}

func computeRevisionSnapshotDigest(snapshot RevisionSnapshot) (string, error) {
	copy := snapshot
	copy.SnapshotDigest = ""
	data, err := json.Marshal(copy)
	if err != nil {
		return "", fmt.Errorf("compute Revision Snapshot digest: %w", err)
	}
	sum := sha256.Sum256(append([]byte(revisionSnapshotDigestDomain), data...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func cloneDecisionValues(values []DecisionValue) []DecisionValue {
	cloned := make([]DecisionValue, len(values))
	for index, value := range values {
		cloned[index] = DecisionValue{ID: value.ID, Value: cloneJSONValue(value.Value)}
	}
	return cloned
}

func validRevisionArea(area RevisionArea) bool {
	switch area {
	case RevisionAreaProfile, RevisionAreaRepositoryRules, RevisionAreaDivergences, RevisionAreaFiles:
		return true
	default:
		return false
	}
}
