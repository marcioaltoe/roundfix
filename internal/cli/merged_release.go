package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func loadDeliveryMergedHeads(
	ctx context.Context,
	reader *store.Store,
	repository string,
) ([]runworktree.MergedHead, error) {
	if reader == nil {
		return nil, errors.New("read Delivery Queue merge records: Run Database is unavailable")
	}
	loadedRepository := strings.TrimSpace(repository)
	if loadedRepository == "" {
		return nil, errors.New("read Delivery Queue merge records: repository is required")
	}
	resolvedRepository, err := filepath.EvalSymlinks(loadedRepository)
	if err != nil {
		return nil, fmt.Errorf("resolve repository for Delivery Queue merge records %q: %w", loadedRepository, err)
	}
	repositories := []string{resolvedRepository}
	if filepath.Clean(loadedRepository) != filepath.Clean(resolvedRepository) {
		repositories = append(repositories, loadedRepository)
	}

	merged := make([]runworktree.MergedHead, 0)
	for _, gitRoot := range repositories {
		queue, found, err := reader.DeliveryQueue(ctx, gitRoot)
		if err != nil {
			return nil, fmt.Errorf("read Delivery Queue merge records for repository %q: %w", gitRoot, err)
		}
		if !found {
			continue
		}
		for _, item := range queue.Items {
			if item.Stage != store.DeliveryStageMerged ||
				strings.TrimSpace(item.MergeCommit) == "" ||
				len(item.CandidateCommits) == 0 {
				continue
			}
			head := strings.TrimSpace(item.CandidateCommits[len(item.CandidateCommits)-1])
			if head == "" {
				continue
			}
			merged = append(merged, runworktree.MergedHead{
				SpecSlug:     strings.TrimSpace(item.SpecSlug),
				TargetBranch: strings.TrimSpace(item.Branch),
				Head:         head,
				MergeCommit:  strings.TrimSpace(item.MergeCommit),
				PullRequest:  strings.TrimSpace(item.PullRequestNumber),
			})
		}
	}
	return merged, nil
}
