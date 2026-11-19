package cli

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"

	"roundfix/internal/releaseplan"
)

// ManagedPaths reads only the two committed records at the range endpoints.
func (source releasePlanGitSource) ManagedPaths(ctx context.Context, releaseRange releaseplan.Range) (releaseplan.ManagedPaths, error) {
	var managed releaseplan.ManagedPaths
	root, err := source.git(ctx, source.workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return managed, releaseplan.GitSourceError{
			Operation: "resolve managed paths Git root", Err: err,
			NextAction: "run the Release Plan Command inside a Git repository",
		}
	}
	root = filepath.Clean(strings.TrimSpace(root))
	for index, revision := range []string{releaseRange.Target.CommitSHA, releaseRange.Base.CommitSHA} {
		lock, err := source.managedRecord(ctx, root, revision, "skills-lock.json")
		if err != nil {
			return managed, err
		}
		var record struct {
			Skills map[string]json.RawMessage `json:"skills"`
		}
		if json.Unmarshal([]byte(lock), &record) == nil {
			if index == 0 {
				_, managed.Adopter = record.Skills["roundfix"]
			}
			for name := range record.Skills {
				managed.Skills = append(managed.Skills, name)
			}
		}
		manifest, err := source.managedRecord(ctx, root, revision, "docs/agents/setup-context.json")
		if err != nil {
			return managed, err
		}
		var setup struct {
			ManagedArtifacts []struct {
				Path string `json:"path"`
			} `json:"managedArtifacts"`
		}
		if json.Unmarshal([]byte(manifest), &setup) == nil {
			for _, artifact := range setup.ManagedArtifacts {
				managed.Files = append(managed.Files, artifact.Path)
			}
		}
	}
	slices.Sort(managed.Files)
	managed.Files = slices.Compact(managed.Files)
	slices.Sort(managed.Skills)
	managed.Skills = slices.Compact(managed.Skills)
	return managed, nil
}

func (source releasePlanGitSource) managedRecord(ctx context.Context, root, revision, name string) (string, error) {
	listing, err := source.git(ctx, root, "ls-tree", "-z", revision, "--", name)
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation: "inspect managed record", Ref: revision, Paths: []string{name}, Err: err,
			NextAction: "verify the range references readable committed records",
		}
	}
	if listing == "" {
		return "", nil
	}
	data, err := source.git(ctx, root, "show", revision+":"+name)
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation: "read managed record", Ref: revision, Paths: []string{name}, Err: err,
			NextAction: "verify the range references readable committed records",
		}
	}
	return data, nil
}

var _ releaseplan.ManagedPathSource = releasePlanGitSource{}
