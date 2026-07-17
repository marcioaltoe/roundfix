package releaseplan

import "context"

const SchemaVersion = "roundfix.release-plan/v1"

// State is the final Release Plan decision state.
type State string

const (
	StateReady                        State = "ready"
	StateApprovalRequired             State = "approval_required"
	StateManualClassificationRequired State = "manual_classification_required"
	StateNoRelease                    State = "no_release"
)

// Impact is the release impact assigned to a change or classified range.
type Impact string

const (
	ImpactNone  Impact = "none"
	ImpactPatch Impact = "patch"
	ImpactMinor Impact = "minor"
	ImpactMajor Impact = "major"
)

// IncrementKind is the semantic-version component the proposal increments.
type IncrementKind string

const (
	IncrementNone  IncrementKind = "none"
	IncrementPatch IncrementKind = "patch"
	IncrementMinor IncrementKind = "minor"
	IncrementMajor IncrementKind = "major"
)

// ClassificationSource identifies how the release impact was established.
type ClassificationSource string

const (
	SourceConventionalCommit ClassificationSource = "conventional_commit"
	SourceMaintenanceOnly    ClassificationSource = "maintenance_only"
	SourceManual             ClassificationSource = "manual"
	SourceMixed              ClassificationSource = "mixed"
)

// VersionRef identifies the release base used by the plan.
type VersionRef struct {
	Tag       string
	Version   Version
	CommitSHA string
}

// RevisionRef identifies the target revision analyzed by the plan.
type RevisionRef struct {
	Name      string
	CommitSHA string
}

// Classification summarizes the final release-impact classification.
type Classification struct {
	Source                       ClassificationSource
	Impact                       Impact
	Breaking                     bool
	ManualReason                 string
	ManualClassificationRequired bool
	BlockingCommits              []string
}

// ChangeEvidence carries the per-commit evidence used by classifiers and
// renderers to justify or block a Release Plan.
type ChangeEvidence struct {
	CommitSHA                      string
	Subject                        string
	ConventionalType               string
	Breaking                       bool
	AutomaticImpact                Impact
	CrossesMaintenanceOnlyBoundary bool
}

// Commit is one normalized committed change supplied by a Git adapter. The
// Release Plan domain never reads Git directly.
type Commit struct {
	SHA          string
	Subject      string
	Body         string
	ChangedPaths []string
}

// Range identifies one committed release range after Git ref resolution.
type Range struct {
	Base   VersionRef
	Target RevisionRef
}

// Request is the Release Plan build request before Git range resolution.
type Request struct {
	From         string
	To           string
	ManualImpact Impact
	ManualReason string
}

// GitSource supplies committed release ranges and normalized commits to the
// Release Plan domain. Implementations live at repository boundaries.
type GitSource interface {
	ResolveRange(context.Context, string, string) (Range, error)
	Commits(context.Context, Range) ([]Commit, error)
}

// ClassifyRequest is the normalized input for conservative release-impact
// classification.
type ClassifyRequest struct {
	Commits      []Commit
	ManualImpact Impact
	ManualReason string
}

// ClassificationResult carries the aggregate classification and all commit
// evidence used to derive it.
type ClassificationResult struct {
	State          State
	Classification Classification
	Changes        []ChangeEvidence
}

// Approval captures the human approval boundary for non-patch increments.
type Approval struct {
	Required        bool
	Increment       IncrementKind
	ProposedVersion string
	Question        string
}

// Plan is the versioned Release Plan schema shared by text and JSON renderers.
type Plan struct {
	SchemaVersion   string
	State           State
	Base            VersionRef
	Target          RevisionRef
	Classification  Classification
	ProposedVersion string
	Approval        Approval
	Changes         []ChangeEvidence
}

// NewPlan returns a Release Plan with the supported schema identifier and a
// copied evidence slice.
func NewPlan(base VersionRef, target RevisionRef, classification Classification, proposal Proposal, changes []ChangeEvidence) Plan {
	return Plan{
		SchemaVersion:   SchemaVersion,
		State:           proposal.State,
		Base:            base,
		Target:          target,
		Classification:  classification,
		ProposedVersion: proposal.ProposedVersion,
		Approval:        proposal.Approval,
		Changes:         append([]ChangeEvidence(nil), changes...),
	}
}
