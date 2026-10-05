package worktree

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// ItemBranchReconciliation reports a local item branch and its cleanup proof.
type ItemBranchReconciliation struct {
	Branch, SpecSlug, Head, Worktree string
	Releasable                       bool
	Proof, RefusalReason             string
	evidence                         *itemBranchEvidence
}

type itemBranchEvidence struct {
	root, location string
	live           []string
	delivery       deliveryEvidence
	snapshot       ItemBranchReconciliation
}

var itemBranchPattern = regexp.MustCompile(`^roundfix/deliver-(.+)-[0-9a-f]{16}$`)

// InspectItemBranches inspects local item branches without mutating Git state.
// live contains branches recorded by non-merged Delivery Queue items.
func InspectItemBranches(ctx context.Context, gitRoot, location string, live []string) ([]ItemBranchReconciliation, error) {
	runner := execGitRunner{}
	output, err := runner.Run(ctx, gitRoot, "for-each-ref", "--format=%(refname)", "refs/heads/roundfix/deliver-*")
	if err != nil {
		return nil, fmt.Errorf("list item branches: %w", err)
	}
	worktrees, err := listRegisteredWorktrees(ctx, runner, gitRoot)
	if err != nil {
		return nil, fmt.Errorf("inspect item branch worktrees: %w", err)
	}
	results := make([]ItemBranchReconciliation, 0)
	for _, ref := range strings.Fields(output) {
		branch := strings.TrimPrefix(ref, "refs/heads/")
		match := itemBranchPattern.FindStringSubmatch(branch)
		if match == nil {
			continue
		}
		slug := match[1]
		if clean, err := cleanPathSegment(slug); err != nil || clean != slug {
			continue
		}
		result, err := inspectItemBranch(ctx, runner, gitRoot, location, branch, slug, live, worktrees)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func inspectItemBranch(ctx context.Context, runner gitRunner, root, location, branch, slug string, live []string, worktrees []registeredWorktree) (ItemBranchReconciliation, error) {
	result := ItemBranchReconciliation{Branch: branch, SpecSlug: slug, Worktree: registeredBranchPath(worktrees, branch)}
	head, err := runner.Run(ctx, root, "rev-parse", "--verify", "refs/heads/"+branch+"^{commit}")
	if err != nil {
		return result, fmt.Errorf("read item branch %q head: %w", branch, err)
	}
	result.Head = strings.TrimSpace(head)
	if slices.Contains(live, branch) {
		result.RefusalReason = fmt.Sprintf("item branch belongs to live Delivery Queue item %q", slug)
		return result, nil
	}
	delivery, found := provenDeliveryEvidence(ctx, runner, root, slug, result.Head)
	if !found {
		result.RefusalReason = fmt.Sprintf("Spec %q has no merge evidence on the default branch", slug)
		return result, nil
	}
	proof, refusal, proven := supersededByDelivery(ctx, runner, root, slug, result.Head, delivery)
	if !proven {
		if refusal == "" {
			refusal = fmt.Sprintf("item branch %q supersession could not be proven", branch)
		}
		result.RefusalReason = refusal
		return result, nil
	}
	if result.Worktree != "" {
		ref, err := ItemRefFor(root, location, branch)
		if err != nil {
			return result, fmt.Errorf("derive item branch %q worktree: %w", branch, err)
		}
		for _, registered := range worktrees {
			if registered.Branch != branch {
				continue
			}
			if !samePath(registered.Path, ref.Path) {
				result.Worktree = registered.Path
				result.RefusalReason = fmt.Sprintf("item worktree %q is not at derived path %q", result.Worktree, ref.Path)
				return result, nil
			}
		}
		status, err := runner.Run(ctx, result.Worktree, "status", "--porcelain=v1", "-z", "--untracked-files=all")
		if err != nil {
			result.RefusalReason = fmt.Sprintf("item worktree %q cleanliness could not be inspected: %v", result.Worktree, err)
			return result, nil
		}
		if status != "" {
			result.RefusalReason = fmt.Sprintf("item worktree %q is dirty", result.Worktree)
			return result, nil
		}
	}
	result.Proof = strings.Replace(proof, "Run work is", "item branch is", 1)
	result.Releasable = true
	result.evidence = &itemBranchEvidence{root: root, location: location, live: slices.Clone(live), delivery: delivery, snapshot: result}
	return result, nil
}

// ApplyItemBranch re-proves the captured evidence before removing a clean
// registered worktree without force and deleting its item branch.
func ApplyItemBranch(ctx context.Context, result ItemBranchReconciliation) error {
	if !result.Releasable || result.evidence == nil {
		reason := result.RefusalReason
		if reason == "" {
			reason = "item branch has no release evidence"
		}
		return errors.New(reason)
	}
	evidence := result.evidence
	snapshot := result
	snapshot.evidence = nil
	if snapshot != evidence.snapshot {
		return errors.New("item branch evidence is stale: inspection changed")
	}
	runner := execGitRunner{}
	worktrees, err := listRegisteredWorktrees(ctx, runner, evidence.root)
	if err != nil {
		return fmt.Errorf("reinspect item worktrees: %w", err)
	}
	fresh, err := inspectItemBranch(ctx, runner, evidence.root, evidence.location, result.Branch, result.SpecSlug, evidence.live, worktrees)
	if err != nil {
		return fmt.Errorf("item branch evidence is stale: %w", err)
	}
	if !fresh.Releasable || fresh.Head != result.Head || fresh.Worktree != result.Worktree || fresh.evidence.delivery != evidence.delivery {
		return fmt.Errorf("item branch evidence is stale: head, delivery, default head or worktree changed; %s", fresh.RefusalReason)
	}
	if fresh.Worktree != "" {
		if _, err := runner.Run(ctx, evidence.root, "worktree", "remove", fresh.Worktree); err != nil {
			return fmt.Errorf("remove clean item worktree %q: %w", fresh.Worktree, err)
		}
	}
	if _, err := runner.Run(ctx, evidence.root, "branch", "-D", fresh.Branch); err != nil {
		return fmt.Errorf("delete item branch %q: %w", fresh.Branch, err)
	}
	return nil
}
