package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func (workflow *commandDeliveryWorkflow) numberedFileReport(ctx context.Context, workDir, specsRoot, specSlug string, graph *spec.Graph) (delivery.NumberedFileReport, error) {
	empty := delivery.NumberedFileReport{}
	branch := preflight.DetectDefaultBranch(ctx, workDir, "", workflow.git)
	if branch.Source == preflight.DefaultBranchUndetermined {
		return empty, nil
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	ref := "refs/remotes/" + remote + "/" + branch.Name
	defaultHead, err := workflow.git.RunGit(ctx, workDir, "rev-parse", ref+"^{commit}")
	if err != nil {
		return empty, nil
	}
	config, err := roundconfig.DeliveryConfigAtCommit(ctx, workflow.git, workDir, workflow.loaded.UserConfigPath, strings.TrimSpace(defaultHead))
	if err != nil {
		return empty, err
	}
	if len(config.Delivery.NumberedFiles) == 0 {
		return empty, nil
	}

	// Refresh only the default tracking ref; leave tags and FETCH_HEAD alone.
	if _, err := workflow.git.RunGit(ctx, workDir, "fetch", "--no-tags", "--no-write-fetch-head", remote, "+refs/heads/"+branch.Name+":"+ref); err != nil {
		return empty, fmt.Errorf("fetch numbered-file default branch: %w", err)
	}
	defaultHead, err = workflow.git.RunGit(ctx, workDir, "rev-parse", ref+"^{commit}")
	if err != nil {
		return empty, fmt.Errorf("resolve numbered-file default commit: %w", err)
	}
	defaultHead = strings.TrimSpace(defaultHead)
	config, err = roundconfig.DeliveryConfigAtCommit(ctx, workflow.git, workDir, workflow.loaded.UserConfigPath, defaultHead)
	if err != nil {
		return empty, err
	}
	if len(config.Delivery.NumberedFiles) == 0 {
		return empty, nil
	}
	defaultFiles, err := workflow.numberedFilesAt(ctx, workDir, defaultHead, config.Delivery)
	if err != nil {
		return empty, err
	}
	subject, err := workflow.numberedFileClaims(ctx, workDir, defaultHead, "HEAD", graph, config.Delivery)
	if err != nil {
		return empty, fmt.Errorf("read item claims: %w", err)
	}
	queue, _, err := workflow.store.DeliveryQueue(ctx, workflow.loaded.GitRoot)
	if err != nil {
		return empty, err
	}
	subjectIndex := -1
	for index, item := range queue.Items {
		if item.SpecSlug == specSlug {
			subjectIndex = index
			break
		}
	}
	if subjectIndex < 0 {
		return delivery.FindNumberedFileCollisions(defaultFiles, subject, nil), nil
	}
	relativeRoot, err := filepath.Rel(workDir, specsRoot)
	if err != nil {
		return empty, fmt.Errorf("resolve numbered-file Specs Root: %w", err)
	}
	var ahead []delivery.NumberedFileClaims
	for _, item := range queue.Items[:subjectIndex] {
		if item.Stage == store.DeliveryStageMerged {
			continue
		}
		head := defaultHead
		if item.Branch != "" {
			resolved, err := workflow.git.RunGit(ctx, workDir, "rev-parse", "refs/heads/"+item.Branch+"^{commit}")
			if err == nil {
				head = strings.TrimSpace(resolved)
			}
		}
		planRoot := specsRoot
		if item.Worktree != "" {
			recordedRoot := filepath.Join(item.Worktree, relativeRoot)
			if info, err := os.Stat(filepath.Join(recordedRoot, item.SpecSlug, "_tasks.md")); err == nil && !info.IsDir() {
				planRoot = recordedRoot
			}
		}
		plan, err := spec.Load(planRoot, item.SpecSlug)
		if err != nil {
			// An unavailable plan leaves the branch claims intact.
			plan = nil
		}
		paths, err := workflow.numberedFileClaims(ctx, workDir, defaultHead, head, plan, config.Delivery)
		if err != nil {
			return empty, fmt.Errorf("read earlier item %q claims: %w", item.SpecSlug, err)
		}
		ahead = append(ahead, delivery.NumberedFileClaims{Owner: item.SpecSlug, Paths: paths})
	}
	return delivery.FindNumberedFileCollisions(defaultFiles, subject, ahead), nil
}

func (workflow *commandDeliveryWorkflow) numberedFilesAt(ctx context.Context, workDir, commit string, declaration roundconfig.Delivery) ([]string, error) {
	listing, err := workflow.git.RunGit(ctx, workDir, "ls-tree", "-r", "--name-only", "-z", commit)
	if err != nil {
		return nil, fmt.Errorf("list numbered files at %s: %w", commit, err)
	}
	var paths []string
	for _, name := range strings.Split(listing, "\x00") {
		if declaration.IsNumberedFile(name) {
			paths = append(paths, name)
		}
	}
	return paths, nil
}

func (workflow *commandDeliveryWorkflow) numberedFileClaims(ctx context.Context, workDir, defaultHead, head string, graph *spec.Graph, declaration roundconfig.Delivery) ([]string, error) {
	base, err := workflow.git.RunGit(ctx, workDir, "merge-base", defaultHead, head)
	if err != nil {
		return nil, fmt.Errorf("find numbered-file merge base: %w", err)
	}
	base = strings.TrimSpace(base)
	added, err := workflow.git.RunGit(ctx, workDir, "diff", "--name-only", "--no-renames", "--diff-filter=A", "-z", base, head)
	if err != nil {
		return nil, fmt.Errorf("list added numbered files: %w", err)
	}
	var paths []string
	seen := make(map[string]bool)
	add := func(name string) {
		if declaration.IsNumberedFile(name) && !seen[name] {
			paths = append(paths, name)
			seen[name] = true
		}
	}
	for _, name := range strings.Split(added, "\x00") {
		add(name)
	}
	if graph == nil {
		return paths, nil
	}
	baseFiles, err := workflow.numberedFilesAt(ctx, workDir, base, declaration)
	if err != nil {
		return nil, err
	}
	present := make(map[string]bool, len(baseFiles))
	for _, name := range baseFiles {
		present[name] = true
	}
	for _, task := range graph.Tasks {
		if task.Status == spec.StatusCompleted || task.Type == spec.TaskTypeQA {
			continue
		}
		for _, ref := range task.Context {
			if ref.Kind == spec.ContextKindCreates && !present[ref.Path] {
				add(ref.Path)
			}
		}
	}
	return paths, nil
}
