package releaseplan

import (
	"errors"
	"reflect"
	"testing"
)

// Suite: Release Plan classifier
// Invariant: normalized committed changes classify to the highest required impact without Git or I/O.
// Boundary IN: Conventional Commit evidence, maintenance-only paths, manual impact validation.
// Boundary OUT: Git range loading, version proposal rendering, CLI exit mapping, repository mutation checks.

func TestClassifyCommit(t *testing.T) {
	tests := []struct {
		name         string
		commit       Commit
		wantType     string
		wantImpact   Impact
		wantBreaking bool
		wantCrosses  bool
	}{
		{
			name: "breaking marker in subject",
			commit: Commit{
				SHA:          "break-subject",
				Subject:      "feat!: remove legacy flag",
				ChangedPaths: []string{"internal/cli/release.go"},
			},
			wantType:     "feat",
			wantImpact:   ImpactMajor,
			wantBreaking: true,
			wantCrosses:  true,
		},
		{
			name: "breaking footer in body",
			commit: Commit{
				SHA:          "break-footer",
				Subject:      "fix: repair config parsing",
				Body:         "Repair config handling.\n\nBREAKING CHANGE: old config keys are rejected",
				ChangedPaths: []string{"internal/config/config.go"},
			},
			wantType:     "fix",
			wantImpact:   ImpactMajor,
			wantBreaking: true,
			wantCrosses:  true,
		},
		{
			name: "scoped feature is minor",
			commit: Commit{
				SHA:          "feature",
				Subject:      "feat(cli): add release plan",
				ChangedPaths: []string{"cmd/roundfix/main.go"},
			},
			wantType:    "feat",
			wantImpact:  ImpactMinor,
			wantCrosses: true,
		},
		{
			name: "fix is patch",
			commit: Commit{
				SHA:          "fix",
				Subject:      "fix: keep artifact names stable",
				ChangedPaths: []string{"internal/release/release.go"},
			},
			wantType:    "fix",
			wantImpact:  ImpactPatch,
			wantCrosses: true,
		},
		{
			name: "perf is patch",
			commit: Commit{
				SHA:          "perf",
				Subject:      "perf: reduce release scan allocations",
				ChangedPaths: []string{"internal/release/release.go"},
			},
			wantType:    "perf",
			wantImpact:  ImpactPatch,
			wantCrosses: true,
		},
		{
			name: "other conventional type has no automatic increment",
			commit: Commit{
				SHA:          "docs",
				Subject:      "docs: record release plan decision",
				ChangedPaths: []string{"docs/specs/0034-release-plan/_techspec.md"},
			},
			wantType:   "docs",
			wantImpact: ImpactNone,
		},
		{
			name: "malformed subject has no automatic increment",
			commit: Commit{
				SHA:          "malformed",
				Subject:      "repair release docs",
				ChangedPaths: []string{"docs/user-guide/commands.md"},
			},
			wantImpact:  ImpactNone,
			wantCrosses: true,
		},
		{
			name: "malformed type with bang is not guessed",
			commit: Commit{
				SHA:          "malformed-bang",
				Subject:      "123!: not a conventional commit type",
				ChangedPaths: []string{"internal/cli/release.go"},
			},
			wantImpact:  ImpactNone,
			wantCrosses: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyCommit(tt.commit)
			if got.CommitSHA != tt.commit.SHA {
				t.Fatalf("CommitSHA = %q, want %q", got.CommitSHA, tt.commit.SHA)
			}
			if got.Subject != tt.commit.Subject {
				t.Fatalf("Subject = %q, want %q", got.Subject, tt.commit.Subject)
			}
			if got.ConventionalType != tt.wantType {
				t.Fatalf("ConventionalType = %q, want %q", got.ConventionalType, tt.wantType)
			}
			if got.AutomaticImpact != tt.wantImpact {
				t.Fatalf("AutomaticImpact = %q, want %q", got.AutomaticImpact, tt.wantImpact)
			}
			if got.Breaking != tt.wantBreaking {
				t.Fatalf("Breaking = %v, want %v", got.Breaking, tt.wantBreaking)
			}
			if got.CrossesMaintenanceOnlyBoundary != tt.wantCrosses {
				t.Fatalf("CrossesMaintenanceOnlyBoundary = %v, want %v", got.CrossesMaintenanceOnlyBoundary, tt.wantCrosses)
			}
		})
	}
}

func TestMaintenanceOnly(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "docs/specs/0034-release-plan/task_02.md", want: true},
		{path: "docs/adr/0048-release-planning-is-read-only-and-confirmation-gated.md", want: true},
		{path: "docs/findings/2026-07-16-release-version-strategy-and-approval-gates.md", want: true},
		{path: "docs/handoffs/2026-07-16-release-plan.md", want: true},
		{path: "internal/releaseplan/classifier_test.go", want: true},
		{path: ".agents/skills/setup-context-driven/tests/test_assets.py", want: true},
		{path: "internal/releaseplan/testdata/commit.txt", want: true},
		{path: "internal/releaseplan/fixtures/plan.json", want: true},
		{path: ".github/workflows/ci-conventions.yml", want: true},
		{path: ".github/workflows/secondbrain-sync.yml", want: true},
		{path: ".github/workflows/release.yml", want: false},
		{path: "internal/cli/release.go", want: false},
		{path: "cmd/roundfix/main.go", want: false},
		{path: ".agents/skills/roundfix/SKILL.md", want: false},
		{path: "package.json", want: false},
		{path: "go.mod", want: false},
		{path: "README.md", want: false},
		{path: "docs/user-guide/commands.md", want: false},
		{path: "../outside.md", want: false},
		{path: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := MaintenanceOnly(tt.path); got != tt.want {
				t.Fatalf("MaintenanceOnly(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestClassifyChanges(t *testing.T) {
	breakingCommit := Commit{SHA: "breaking", Subject: "feat!: remove deprecated flag", ChangedPaths: []string{"internal/cli/release.go"}}
	featureCommit := Commit{SHA: "feature", Subject: "feat: add release plan command", ChangedPaths: []string{"internal/cli/release.go"}}
	fixCommit := Commit{SHA: "fix", Subject: "fix: keep artifact validation deterministic", ChangedPaths: []string{"internal/release/release.go"}}
	maintenanceCommit := Commit{SHA: "docs", Subject: "docs: record release plan task", ChangedPaths: []string{"docs/specs/0034-release-plan/task_02.md"}}
	maintenanceFixCommit := Commit{SHA: "fix-test", Subject: "fix: repair release plan test", ChangedPaths: []string{"internal/releaseplan/classifier_test.go"}}
	maintenanceBreakingCommit := Commit{SHA: "breaking-test", Subject: "feat!: remove obsolete test helper", ChangedPaths: []string{"internal/releaseplan/classifier_test.go"}}

	tests := []struct {
		name                string
		request             ClassifyRequest
		wantState           State
		wantImpact          Impact
		wantBreaking        bool
		wantSource          ClassificationSource
		wantManualRequired  bool
		wantBlockingCommits []string
	}{
		{
			name:         "breaking outranks feature fix and none",
			request:      ClassifyRequest{Commits: []Commit{maintenanceCommit, fixCommit, featureCommit, breakingCommit}},
			wantState:    StateReady,
			wantImpact:   ImpactMajor,
			wantBreaking: true,
			wantSource:   SourceMixed,
		},
		{
			name:         "breaking outranks every order",
			request:      ClassifyRequest{Commits: []Commit{breakingCommit, featureCommit, fixCommit, maintenanceCommit}},
			wantState:    StateReady,
			wantImpact:   ImpactMajor,
			wantBreaking: true,
			wantSource:   SourceMixed,
		},
		{
			name:       "feature outranks fix and none",
			request:    ClassifyRequest{Commits: []Commit{fixCommit, maintenanceCommit, featureCommit}},
			wantState:  StateReady,
			wantImpact: ImpactMinor,
			wantSource: SourceMixed,
		},
		{
			name:       "fix outranks none",
			request:    ClassifyRequest{Commits: []Commit{maintenanceCommit, fixCommit}},
			wantState:  StateReady,
			wantImpact: ImpactPatch,
			wantSource: SourceMixed,
		},
		{
			name:       "maintenance only produces no release",
			request:    ClassifyRequest{Commits: []Commit{maintenanceCommit}},
			wantState:  StateNoRelease,
			wantImpact: ImpactNone,
			wantSource: SourceMaintenanceOnly,
		},
		{
			name:       "maintenance only fix produces no release",
			request:    ClassifyRequest{Commits: []Commit{maintenanceFixCommit}},
			wantState:  StateNoRelease,
			wantImpact: ImpactNone,
			wantSource: SourceMaintenanceOnly,
		},
		{
			name:       "maintenance only breaking feature produces no release",
			request:    ClassifyRequest{Commits: []Commit{maintenanceBreakingCommit}},
			wantState:  StateNoRelease,
			wantImpact: ImpactNone,
			wantSource: SourceMaintenanceOnly,
		},
		{
			name: "ambiguous shipped surface requires manual classification",
			request: ClassifyRequest{Commits: []Commit{
				{SHA: "ambiguous", Subject: "chore: adjust release behavior", ChangedPaths: []string{"internal/cli/release.go"}},
				maintenanceCommit,
			}},
			wantState:           StateManualClassificationRequired,
			wantImpact:          ImpactNone,
			wantSource:          SourceMixed,
			wantManualRequired:  true,
			wantBlockingCommits: []string{"ambiguous"},
		},
		{
			name: "valid manual classification resolves ambiguity",
			request: ClassifyRequest{
				Commits: []Commit{
					{SHA: "ambiguous", Subject: "chore: adjust release behavior", ChangedPaths: []string{"internal/cli/release.go"}},
					fixCommit,
				},
				ManualImpact: ImpactMinor,
				ManualReason: "compatible release plan behavior was added behind the new command",
			},
			wantState:           StateReady,
			wantImpact:          ImpactMinor,
			wantSource:          SourceMixed,
			wantManualRequired:  false,
			wantBlockingCommits: []string{"ambiguous"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClassifyChanges(tt.request)
			if err != nil {
				t.Fatalf("ClassifyChanges: %v", err)
			}
			if got.State != tt.wantState {
				t.Fatalf("State = %q, want %q", got.State, tt.wantState)
			}
			if got.Classification.Impact != tt.wantImpact {
				t.Fatalf("Impact = %q, want %q", got.Classification.Impact, tt.wantImpact)
			}
			if got.Classification.Breaking != tt.wantBreaking {
				t.Fatalf("Breaking = %v, want %v", got.Classification.Breaking, tt.wantBreaking)
			}
			if got.Classification.Source != tt.wantSource {
				t.Fatalf("Source = %q, want %q", got.Classification.Source, tt.wantSource)
			}
			if got.Classification.ManualClassificationRequired != tt.wantManualRequired {
				t.Fatalf("ManualClassificationRequired = %v, want %v", got.Classification.ManualClassificationRequired, tt.wantManualRequired)
			}
			if !reflect.DeepEqual(got.Classification.BlockingCommits, tt.wantBlockingCommits) {
				t.Fatalf("BlockingCommits = %v, want %v", got.Classification.BlockingCommits, tt.wantBlockingCommits)
			}
			if len(got.Changes) != len(tt.request.Commits) {
				t.Fatalf("Changes length = %d, want %d", len(got.Changes), len(tt.request.Commits))
			}
		})
	}
}

func TestValidateManualImpact(t *testing.T) {
	tests := []struct {
		name    string
		impact  Impact
		reason  string
		minimum Impact
		wantErr error
	}{
		{
			name:    "valid equal impact",
			impact:  ImpactPatch,
			reason:  "ambiguous shipped change is a compatible fix",
			minimum: ImpactPatch,
		},
		{
			name:    "valid higher impact",
			impact:  ImpactMinor,
			reason:  "ambiguous shipped change adds compatible behavior",
			minimum: ImpactPatch,
		},
		{
			name:    "empty reason",
			impact:  ImpactPatch,
			reason:  " \t",
			minimum: ImpactNone,
			wantErr: ErrManualReasonRequired,
		},
		{
			name:    "impact below automatic minimum",
			impact:  ImpactPatch,
			reason:  "try to classify a feature as a fix",
			minimum: ImpactMinor,
			wantErr: ErrManualImpactTooLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateManualImpact(tt.impact, tt.reason, tt.minimum)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("ValidateManualImpact: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateManualImpact succeeded, want %v", tt.wantErr)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateManualImpact error = %v, want errors.Is(..., %v)", err, tt.wantErr)
			}
		})
	}
}

func TestClassifyChangesRejectsInvalidManualImpactWithoutPartialResult(t *testing.T) {
	tests := []struct {
		name    string
		request ClassifyRequest
		wantErr error
	}{
		{
			name: "manual reason without impact",
			request: ClassifyRequest{
				Commits:      []Commit{{SHA: "ambiguous", Subject: "chore: adjust release behavior", ChangedPaths: []string{"internal/cli/release.go"}}},
				ManualReason: "compatible fix",
			},
			wantErr: ErrManualImpactRequired,
		},
		{
			name: "manual impact below automatic minimum",
			request: ClassifyRequest{
				Commits:      []Commit{{SHA: "feature", Subject: "feat: add release planning", ChangedPaths: []string{"internal/cli/release.go"}}},
				ManualImpact: ImpactPatch,
				ManualReason: "try to downgrade a feature to patch",
			},
			wantErr: ErrManualImpactTooLow,
		},
		{
			name: "unsupported manual impact",
			request: ClassifyRequest{
				Commits:      []Commit{{SHA: "ambiguous", Subject: "chore: adjust release behavior", ChangedPaths: []string{"internal/cli/release.go"}}},
				ManualImpact: Impact("revision"),
				ManualReason: "unsupported impact name",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClassifyChanges(tt.request)
			if err == nil {
				t.Fatal("ClassifyChanges succeeded, want error")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("ClassifyChanges error = %v, want errors.Is(..., %v)", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				var impactErr UnknownImpactError
				if !errors.As(err, &impactErr) {
					t.Fatalf("ClassifyChanges error = %T, want UnknownImpactError", err)
				}
			}
			if !reflect.DeepEqual(got, ClassificationResult{}) {
				t.Fatalf("ClassifyChanges result = %+v, want zero result on invalid manual impact", got)
			}
		})
	}
}
