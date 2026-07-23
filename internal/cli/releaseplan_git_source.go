package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
)

type releasePlanGitSource struct {
	workDir string
	runner  preflight.GitRunner
}

type releasePlanResetSource struct {
	git releasePlanGitSource
	gh  preflight.GHRunner
}

func newReleasePlanGitSource(workDir string, runner preflight.GitRunner) releasePlanGitSource {
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	return releasePlanGitSource{workDir: workDir, runner: runner}
}

func newReleasePlanResetSource(workDir string, gitRunner preflight.GitRunner, ghRunner preflight.GHRunner) releasePlanResetSource {
	if ghRunner == nil {
		ghRunner = preflight.ExecGHRunner{}
	}
	return releasePlanResetSource{
		git: newReleasePlanGitSource(workDir, gitRunner),
		gh:  ghRunner,
	}
}

func (source releasePlanResetSource) ResolveTarget(ctx context.Context) (releaseplan.RevisionRef, error) {
	root, err := source.gitRoot(ctx)
	if err != nil {
		return releaseplan.RevisionRef{}, err
	}
	if err := source.git.requireCleanWorktree(ctx, root); err != nil {
		return releaseplan.RevisionRef{}, err
	}
	commitSHA, err := source.git.resolveCommit(ctx, root, "HEAD")
	if err != nil {
		return releaseplan.RevisionRef{}, err
	}
	return releaseplan.RevisionRef{Name: "HEAD", CommitSHA: commitSHA}, nil
}

func (source releasePlanResetSource) Tags(ctx context.Context) ([]releaseplan.TagRef, error) {
	root, err := source.gitRoot(ctx)
	if err != nil {
		return nil, err
	}
	local, err := source.localStableTags(ctx, root)
	if err != nil {
		return nil, err
	}
	remote, err := source.remoteStableTags(ctx, root)
	if err != nil {
		return nil, err
	}
	return append(local, remote...), nil
}

func (source releasePlanResetSource) Releases(ctx context.Context) ([]releaseplan.ReleaseRef, error) {
	root, err := source.gitRoot(ctx)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	output, err := source.gh.RunGH(
		ctx,
		root,
		"api",
		"--method",
		"GET",
		"--paginate",
		"--slurp",
		"repos/{owner}/{repo}/releases?per_page=100",
	)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil {
		return nil, fmt.Errorf("read complete paginated GitHub Release inventory: %w", err)
	}
	releases, err := parseReleasePlanGitHubPages(output)
	if err != nil {
		return nil, fmt.Errorf("read complete paginated GitHub Release inventory: %w", err)
	}
	return releases, nil
}

func (source releasePlanResetSource) localStableTags(ctx context.Context, root string) ([]releaseplan.TagRef, error) {
	output, err := source.git.git(ctx, root, "for-each-ref", "--format=%(refname)", "refs/tags")
	if err != nil {
		return nil, fmt.Errorf("inventory local stable tags: %w", err)
	}
	var tags []releaseplan.TagRef
	for _, ref := range splitNonEmptyLines(output) {
		name := strings.TrimPrefix(ref, "refs/tags/")
		if _, err := releaseplan.ParseStableVersion(name); err != nil {
			continue
		}
		commitSHA, err := source.git.git(ctx, root, "rev-parse", "--verify", ref+"^{commit}")
		if err != nil {
			return nil, fmt.Errorf("inventory local stable tag %q target commit: %w", ref, err)
		}
		commitSHA = strings.TrimSpace(commitSHA)
		if commitSHA == "" {
			return nil, fmt.Errorf("inventory local stable tag %q target commit: git returned an empty commit", ref)
		}
		tags = append(tags, releaseplan.TagRef{
			Name:         name,
			Source:       releaseplan.TagSourceLocal,
			Ref:          ref,
			ImmutableID:  resetTagImmutableID(releaseplan.TagSourceLocal, "", ref, commitSHA),
			TargetCommit: commitSHA,
		})
	}
	return tags, nil
}

func (source releasePlanResetSource) remoteStableTags(ctx context.Context, root string) ([]releaseplan.TagRef, error) {
	output, err := source.git.git(ctx, root, "remote")
	if err != nil {
		return nil, fmt.Errorf("inventory Git remotes for stable tags: %w", err)
	}
	var tags []releaseplan.TagRef
	for _, remote := range splitNonEmptyLines(output) {
		remoteOutput, err := source.git.git(ctx, root, "ls-remote", "--tags", remote)
		if err != nil {
			return nil, fmt.Errorf("inventory remote stable tags from %q: %w", remote, err)
		}
		remoteTags, err := parseReleasePlanRemoteTags(remote, remoteOutput)
		if err != nil {
			return nil, fmt.Errorf("inventory remote stable tags from %q: %w", remote, err)
		}
		tags = append(tags, remoteTags...)
	}
	return tags, nil
}

func (source releasePlanResetSource) gitRoot(ctx context.Context) (string, error) {
	root, err := source.git.git(ctx, source.git.workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release reset Git root",
			NextAction: "run the Release Plan Command inside a Git repository",
			Err:        err,
		}
	}
	root = filepath.Clean(strings.TrimSpace(root))
	if root == "" || root == "." {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release reset Git root",
			NextAction: "run the Release Plan Command inside a Git repository",
			Err:        errors.New("git returned an empty root"),
		}
	}
	return root, nil
}

type releasePlanRemoteTagTarget struct {
	direct string
	peeled string
}

func parseReleasePlanRemoteTags(remote string, output string) ([]releaseplan.TagRef, error) {
	targets := map[string]releasePlanRemoteTagTarget{}
	for _, line := range splitNonEmptyLines(output) {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("malformed ls-remote line %q", line)
		}
		objectID := fields[0]
		ref := fields[1]
		peeled := strings.HasSuffix(ref, "^{}")
		ref = strings.TrimSuffix(ref, "^{}")
		if !strings.HasPrefix(ref, "refs/tags/") {
			return nil, fmt.Errorf("unexpected remote tag ref %q", ref)
		}
		name := strings.TrimPrefix(ref, "refs/tags/")
		if _, err := releaseplan.ParseStableVersion(name); err != nil {
			continue
		}
		target := targets[ref]
		if peeled {
			if target.peeled != "" && target.peeled != objectID {
				return nil, fmt.Errorf("conflicting peeled targets for %q", ref)
			}
			target.peeled = objectID
		} else {
			if target.direct != "" && target.direct != objectID {
				return nil, fmt.Errorf("conflicting direct targets for %q", ref)
			}
			target.direct = objectID
		}
		targets[ref] = target
	}

	tags := make([]releaseplan.TagRef, 0, len(targets))
	for ref, target := range targets {
		commitSHA := target.peeled
		if commitSHA == "" {
			commitSHA = target.direct
		}
		if commitSHA == "" {
			return nil, fmt.Errorf("remote tag %q has no target", ref)
		}
		name := strings.TrimPrefix(ref, "refs/tags/")
		tags = append(tags, releaseplan.TagRef{
			Name:         name,
			Source:       releaseplan.TagSourceRemote,
			Remote:       remote,
			Ref:          ref,
			ImmutableID:  resetTagImmutableID(releaseplan.TagSourceRemote, remote, ref, commitSHA),
			TargetCommit: commitSHA,
		})
	}
	return tags, nil
}

type releasePlanGitHubRelease struct {
	ID              int64  `json:"id"`
	NodeID          string `json:"node_id"`
	Name            string `json:"name"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
}

func parseReleasePlanGitHubPages(output string) ([]releaseplan.ReleaseRef, error) {
	var pages [][]releasePlanGitHubRelease
	if err := json.Unmarshal([]byte(output), &pages); err != nil {
		return nil, fmt.Errorf("decode paginated GitHub Release response: %w", err)
	}
	if len(pages) == 0 {
		return nil, errors.New("GitHub returned no pagination envelope")
	}
	var releases []releaseplan.ReleaseRef
	for pageIndex, page := range pages {
		if page == nil {
			return nil, fmt.Errorf("GitHub Release page %d is null", pageIndex+1)
		}
		for _, item := range page {
			targetCommit := ""
			if isFullGitObjectID(item.TargetCommitish) {
				targetCommit = item.TargetCommitish
			}
			releases = append(releases, releaseplan.ReleaseRef{
				ID:              item.ID,
				NodeID:          item.NodeID,
				Name:            item.Name,
				TagName:         item.TagName,
				TargetCommitish: item.TargetCommitish,
				TargetCommit:    targetCommit,
				ImmutableID:     "github-release:" + strconv.FormatInt(item.ID, 10),
			})
		}
	}
	return releases, nil
}

func resetTagImmutableID(source releaseplan.TagSource, remote string, ref string, targetCommit string) string {
	parts := []string{string(source)}
	if remote != "" {
		parts = append(parts, remote)
	}
	parts = append(parts, ref+"@"+targetCommit)
	return strings.Join(parts, ":")
}

func isFullGitObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	for _, char := range value {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}

func (source releasePlanGitSource) ResolveRange(ctx context.Context, from string, to string) (releaseplan.Range, error) {
	root, err := source.git(ctx, source.workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return releaseplan.Range{}, releaseplan.GitSourceError{
			Operation:  "resolve release plan Git root",
			NextAction: "run the Release Plan Command inside a Git repository",
			Err:        err,
		}
	}
	root = filepath.Clean(strings.TrimSpace(root))
	if root == "" || root == "." {
		return releaseplan.Range{}, releaseplan.GitSourceError{
			Operation:  "resolve release plan Git root",
			NextAction: "run the Release Plan Command inside a Git repository",
			Err:        errors.New("git returned an empty root"),
		}
	}

	if err := source.requireCleanWorktree(ctx, root); err != nil {
		return releaseplan.Range{}, err
	}

	targetName := strings.TrimSpace(to)
	if targetName == "" {
		targetName = "HEAD"
	}
	targetSHA, err := source.resolveCommit(ctx, root, targetName)
	if err != nil {
		return releaseplan.Range{}, err
	}

	baseTag := strings.TrimSpace(from)
	if baseTag == "" {
		baseTag, err = source.latestStableTag(ctx, root, targetSHA)
		if err != nil {
			return releaseplan.Range{}, err
		}
	}
	baseVersion, err := releaseplan.ParseStableVersion(baseTag)
	if err != nil {
		return releaseplan.Range{}, err
	}
	baseSHA, err := source.resolveTagCommit(ctx, root, baseTag)
	if err != nil {
		return releaseplan.Range{}, err
	}

	if err := source.validateRange(ctx, root, baseSHA, targetSHA); err != nil {
		return releaseplan.Range{}, err
	}

	return releaseplan.Range{
		Base: releaseplan.VersionRef{
			Tag:       baseTag,
			Version:   baseVersion,
			CommitSHA: baseSHA,
		},
		Target: releaseplan.RevisionRef{
			Name:      targetName,
			CommitSHA: targetSHA,
		},
	}, nil
}

func (source releasePlanGitSource) Commits(ctx context.Context, releaseRange releaseplan.Range) ([]releaseplan.Commit, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if releaseRange.Base.CommitSHA == "" || releaseRange.Target.CommitSHA == "" {
		return nil, releaseplan.GitSourceError{
			Operation:  "load release range commits",
			NextAction: "resolve a non-empty release range before loading commits",
			Err:        releaseplan.ErrInvalidReleaseRange,
		}
	}
	root, err := source.git(ctx, source.workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, releaseplan.GitSourceError{
			Operation:  "resolve release plan Git root",
			NextAction: "run the Release Plan Command inside a Git repository",
			Err:        err,
		}
	}
	root = filepath.Clean(strings.TrimSpace(root))

	rawList, err := source.git(ctx, root, "rev-list", "--reverse", releaseRange.Base.CommitSHA+".."+releaseRange.Target.CommitSHA)
	if err != nil {
		return nil, releaseplan.GitSourceError{
			Operation:  "load release range commits",
			NextAction: "choose a non-empty range where the base is an ancestor of the target",
			Err:        err,
		}
	}
	shas := splitNonEmptyLines(rawList)
	commits := make([]releaseplan.Commit, 0, len(shas))
	for _, sha := range shas {
		message, err := source.git(ctx, root, "show", "-s", "--format=%B", sha)
		if err != nil {
			return nil, releaseplan.GitSourceError{
				Operation:  "load release range commit message",
				Ref:        sha,
				NextAction: "verify the range references committed objects",
				Err:        err,
			}
		}
		paths, err := source.commitPaths(ctx, root, sha)
		if err != nil {
			return nil, err
		}
		subject, body := splitCommitMessage(message)
		commits = append(commits, releaseplan.Commit{
			SHA:          sha,
			Subject:      subject,
			Body:         body,
			ChangedPaths: paths,
		})
	}
	return commits, nil
}

func (source releasePlanGitSource) requireCleanWorktree(ctx context.Context, root string) error {
	status, err := source.git(ctx, root, "--no-optional-locks", "status", "--porcelain=v1", "-z")
	if err != nil {
		return releaseplan.GitSourceError{
			Operation:  "check release plan worktree status",
			NextAction: "commit, stash, or remove local changes before planning a release",
			Err:        err,
		}
	}
	dirty, err := preflight.ParsePorcelainStatus(status)
	if err != nil {
		return releaseplan.GitSourceError{
			Operation:  "parse release plan worktree status",
			NextAction: "commit, stash, or remove local changes before planning a release",
			Err:        err,
		}
	}
	if len(dirty) == 0 {
		return nil
	}
	return releaseplan.GitSourceError{
		Operation:  "check release plan worktree status",
		Paths:      dirtyPathNames(dirty),
		NextAction: "commit, stash, or remove these paths before planning a release",
		Err:        releaseplan.ErrDirtyWorktree,
	}
}

func (source releasePlanGitSource) resolveCommit(ctx context.Context, root string, ref string) (string, error) {
	if _, err := source.git(ctx, root, "rev-parse", "--verify", ref); err != nil {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release plan target",
			Ref:        ref,
			NextAction: "choose a revision that exists and resolves to a commit",
			Err:        releaseplan.ErrUnresolvedRevision,
		}
	}
	sha, err := source.git(ctx, root, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release plan target",
			Ref:        ref,
			NextAction: "choose a revision that resolves to a commit",
			Err:        releaseplan.ErrNonCommitRevision,
		}
	}
	sha = strings.TrimSpace(sha)
	if sha == "" {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release plan target",
			Ref:        ref,
			NextAction: "choose a revision that resolves to a commit",
			Err:        errors.New("git returned an empty commit"),
		}
	}
	return sha, nil
}

func (source releasePlanGitSource) resolveTagCommit(ctx context.Context, root string, tag string) (string, error) {
	sha, err := source.git(ctx, root, "rev-parse", "--verify", "refs/tags/"+tag+"^{commit}")
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release plan base tag",
			Ref:        tag,
			NextAction: "choose an existing stable vMAJOR.MINOR.PATCH tag",
			Err:        releaseplan.ErrUnresolvedRevision,
		}
	}
	sha = strings.TrimSpace(sha)
	if sha == "" {
		return "", releaseplan.GitSourceError{
			Operation:  "resolve release plan base tag",
			Ref:        tag,
			NextAction: "choose an existing stable vMAJOR.MINOR.PATCH tag",
			Err:        errors.New("git returned an empty commit"),
		}
	}
	return sha, nil
}

func (source releasePlanGitSource) latestStableTag(ctx context.Context, root string, targetSHA string) (string, error) {
	rawTags, err := source.git(ctx, root, "tag", "--merged", targetSHA, "--list", "v*")
	if err != nil {
		return "", releaseplan.GitSourceError{
			Operation:  "find latest stable release tag",
			NextAction: "pass --from with an existing stable vMAJOR.MINOR.PATCH tag",
			Err:        err,
		}
	}
	var selectedTag string
	var selectedVersion releaseplan.Version
	for _, tag := range splitNonEmptyLines(rawTags) {
		version, err := releaseplan.ParseStableVersion(tag)
		if err != nil {
			continue
		}
		if selectedTag == "" || compareVersion(version, selectedVersion) > 0 {
			selectedTag = tag
			selectedVersion = version
		}
	}
	if selectedTag == "" {
		return "", releaseplan.GitSourceError{
			Operation:  "find latest stable release tag",
			NextAction: "create or pass a stable vMAJOR.MINOR.PATCH base tag",
			Err:        releaseplan.ErrNoStableReleaseTag,
		}
	}
	return selectedTag, nil
}

func (source releasePlanGitSource) validateRange(ctx context.Context, root string, baseSHA string, targetSHA string) error {
	if baseSHA == targetSHA {
		return releaseplan.GitSourceError{
			Operation:  "validate release plan range",
			NextAction: "choose a target commit after the base release tag",
			Err:        releaseplan.ErrInvalidReleaseRange,
		}
	}
	ancestor, err := source.isAncestor(ctx, root, baseSHA, targetSHA)
	if err != nil {
		return err
	}
	if ancestor {
		return nil
	}
	reversed, err := source.isAncestor(ctx, root, targetSHA, baseSHA)
	if err != nil {
		return err
	}
	if reversed {
		return releaseplan.GitSourceError{
			Operation:  "validate release plan range",
			NextAction: "choose a base release tag that is an ancestor of the target",
			Err:        releaseplan.ErrInvalidReleaseRange,
		}
	}
	return releaseplan.GitSourceError{
		Operation:  "validate release plan range",
		NextAction: "choose a base release tag reachable from the target",
		Err:        releaseplan.ErrInvalidReleaseRange,
	}
}

func (source releasePlanGitSource) isAncestor(ctx context.Context, root string, ancestor string, descendant string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	_, err := source.runner.RunGit(ctx, root, "merge-base", "--is-ancestor", ancestor, descendant)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return false, ctxErr
	}
	return err == nil, nil
}

func (source releasePlanGitSource) commitPaths(ctx context.Context, root string, sha string) ([]string, error) {
	rawPaths, err := source.git(ctx, root, "diff-tree", "--no-commit-id", "--name-only", "-r", "--root", sha)
	if err != nil {
		return nil, releaseplan.GitSourceError{
			Operation:  "load release range commit paths",
			Ref:        sha,
			NextAction: "verify the range references committed objects",
			Err:        err,
		}
	}
	return splitNonEmptyLines(rawPaths), nil
}

func (source releasePlanGitSource) git(ctx context.Context, workDir string, args ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	runner := source.runner
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	output, err := runner.RunGit(ctx, workDir, args...)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", ctxErr
	}
	return output, err
}

func dirtyPathNames(changes []preflight.ChangedPath) []string {
	paths := make([]string, 0, len(changes))
	for _, change := range changes {
		paths = append(paths, change.Path)
	}
	return paths
}

func splitCommitMessage(message string) (string, string) {
	lines := strings.Split(message, "\n")
	subject := strings.TrimSpace(lines[0])
	if len(lines) == 1 {
		return subject, ""
	}
	return subject, strings.TrimLeft(strings.Join(lines[1:], "\n"), "\n")
}

func splitNonEmptyLines(text string) []string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func compareVersion(left releaseplan.Version, right releaseplan.Version) int {
	for _, pair := range [][2]int{
		{left.Major(), right.Major()},
		{left.Minor(), right.Minor()},
		{left.Patch(), right.Patch()},
	} {
		if pair[0] > pair[1] {
			return 1
		}
		if pair[0] < pair[1] {
			return -1
		}
	}
	return 0
}

var _ releaseplan.GitSource = releasePlanGitSource{}
var _ releaseplan.ResetInventorySource = releasePlanResetSource{}
