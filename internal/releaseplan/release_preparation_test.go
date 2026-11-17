package releaseplan

import (
	"reflect"
	"testing"
)

// Suite: Release-Preparation Commit classification.
// Invariant: only non-empty changes with no impact inside granted release paths prepare a release.
// Boundary IN: normalized commits, declared release paths, exact evidence and classification.
// Boundary OUT: Git, grant discovery, CLI rendering, and release mutation.
func TestReleasePreparationCommitSetsNoImpact(t *testing.T) {
	t.Parallel()
	fix := Commit{SHA: "fix", Subject: "fix: repair x", ChangedPaths: []string{"internal/x.go"}}
	fixEvidence := ChangeEvidence{CommitSHA: "fix", Subject: "fix: repair x", ConventionalType: "fix", AutomaticImpact: ImpactPatch, CrossesMaintenanceOnlyBoundary: true}
	for _, tt := range []struct {
		name         string
		commit       Commit
		releasePaths []string
		withFix      bool
		wantState    State
		wantClass    Classification
		wantEvidence ChangeEvidence
	}{
		{
			name: "inside_release_paths", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}},
			releasePaths: []string{"CHANGELOG.md", "VERSION"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, ReleasePreparation: true},
		},
		{
			name: "cleaned_paths_match", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"./CHANGELOG.md", "VERSION"}},
			releasePaths: []string{"CHANGELOG.md", "VERSION"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, ReleasePreparation: true},
		},
		{
			name: "one_unlisted_path", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION", "README.md"}},
			releasePaths: []string{"CHANGELOG.md", "VERSION"},
			wantState:    StateManualClassificationRequired,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch, ManualClassificationRequired: true, BlockingCommits: []string{"preparation"}},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name: "no_release_paths", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"CHANGELOG.md", "VERSION"}},
			wantState:    StateManualClassificationRequired,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch, ManualClassificationRequired: true, BlockingCommits: []string{"preparation"}},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name: "no_changed_path", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0"},
			releasePaths: []string{"CHANGELOG.md", "VERSION"},
			wantState:    StateManualClassificationRequired,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch, ManualClassificationRequired: true, BlockingCommits: []string{"preparation"}},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name:         "fix_on_release_paths",
			commit:       Commit{SHA: "release-fix", Subject: "fix: correct release notes", ChangedPaths: []string{"CHANGELOG.md"}},
			releasePaths: []string{"CHANGELOG.md"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceConventionalCommit, Impact: ImpactPatch},
			wantEvidence: ChangeEvidence{CommitSHA: "release-fix", Subject: "fix: correct release notes", ConventionalType: "fix", AutomaticImpact: ImpactPatch, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name:         "breaking_on_release_paths",
			commit:       Commit{SHA: "breaking", Subject: "chore!: replace version scheme", ChangedPaths: []string{"CHANGELOG.md"}},
			releasePaths: []string{"CHANGELOG.md"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceConventionalCommit, Impact: ImpactMajor, Breaking: true},
			wantEvidence: ChangeEvidence{CommitSHA: "breaking", Subject: "chore!: replace version scheme", ConventionalType: "chore", Breaking: true, AutomaticImpact: ImpactMajor, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name:         "docs_references_planning_commit",
			commit:       Commit{SHA: "planning", Subject: "docs: plan x", ChangedPaths: []string{"docs/specs/0001-x/_prd.md", "docs/references/evidence.md"}},
			wantState:    StateNoRelease,
			wantClass:    Classification{Source: SourceMaintenanceOnly, Impact: ImpactNone},
			wantEvidence: ChangeEvidence{CommitSHA: "planning", Subject: "docs: plan x", ConventionalType: "docs", AutomaticImpact: ImpactNone},
		},
		{
			name:         "feat_on_release_paths",
			commit:       Commit{SHA: "feature", Subject: "feat: add version metadata", ChangedPaths: []string{"VERSION"}},
			releasePaths: []string{"VERSION"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceConventionalCommit, Impact: ImpactMinor},
			wantEvidence: ChangeEvidence{CommitSHA: "feature", Subject: "feat: add version metadata", ConventionalType: "feat", AutomaticImpact: ImpactMinor, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name:         "perf_on_release_paths",
			commit:       Commit{SHA: "performance", Subject: "perf: reduce version parsing cost", ChangedPaths: []string{"VERSION"}},
			releasePaths: []string{"VERSION"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceConventionalCommit, Impact: ImpactPatch},
			wantEvidence: ChangeEvidence{CommitSHA: "performance", Subject: "perf: reduce version parsing cost", ConventionalType: "perf", AutomaticImpact: ImpactPatch, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name: "preparation_without_release_subject", withFix: true,
			commit:       Commit{SHA: "metadata", Subject: "docs: refresh metadata", ChangedPaths: []string{"VERSION"}},
			releasePaths: []string{"VERSION"},
			wantState:    StateReady,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch},
			wantEvidence: ChangeEvidence{CommitSHA: "metadata", Subject: "docs: refresh metadata", ConventionalType: "docs", AutomaticImpact: ImpactNone, ReleasePreparation: true},
		},
		{
			name: "unlisted_planning_path", withFix: true,
			commit:       Commit{SHA: "preparation", Subject: "chore: release v0.2.0", ChangedPaths: []string{"VERSION", "docs/references/evidence.md"}},
			releasePaths: []string{"VERSION"},
			wantState:    StateManualClassificationRequired,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactPatch, ManualClassificationRequired: true, BlockingCommits: []string{"preparation"}},
			wantEvidence: ChangeEvidence{CommitSHA: "preparation", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, CrossesMaintenanceOnlyBoundary: true},
		},
		{
			name:         "invalid_changed_path",
			commit:       Commit{SHA: "invalid", Subject: "chore: release v0.2.0", ChangedPaths: []string{"../VERSION"}},
			releasePaths: []string{"", "../VERSION", "VERSION"},
			wantState:    StateManualClassificationRequired,
			wantClass:    Classification{Source: SourceMixed, Impact: ImpactNone, ManualClassificationRequired: true, BlockingCommits: []string{"invalid"}},
			wantEvidence: ChangeEvidence{CommitSHA: "invalid", Subject: "chore: release v0.2.0", ConventionalType: "chore", AutomaticImpact: ImpactNone, CrossesMaintenanceOnlyBoundary: true},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commits := []Commit{tt.commit}
			wantChanges := []ChangeEvidence{tt.wantEvidence}
			if tt.withFix {
				commits = append([]Commit{fix}, commits...)
				wantChanges = append([]ChangeEvidence{fixEvidence}, wantChanges...)
			}
			got, err := ClassifyChanges(ClassifyRequest{Commits: commits, ReleasePaths: tt.releasePaths})
			if err != nil {
				t.Fatalf("ClassifyChanges: %v", err)
			}
			want := ClassificationResult{State: tt.wantState, Classification: tt.wantClass, Changes: wantChanges}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("ClassifyChanges = %+v, want %+v", got, want)
			}
		})
	}
}
