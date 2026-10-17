package cli

import (
	"context"
	"fmt"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/preflight"
)

const reviewSegmentLimit = 4

type reviewSegment struct {
	Base      string `json:"base"`
	Head      string `json:"head"`
	DiffBytes int    `json:"diffBytes"`
}

type reviewSegmentAnswer struct {
	result agent.ExecuteResult
	err    error
	diff   string
}

func planReviewSegments(ctx context.Context, root, base, head string, specRoots []string, git preflight.GitRunner, scope reviewSpecContextResult) ([]reviewSegment, []string, []string, error) {
	commits, err := git.RunGit(ctx, root, "rev-list", "--first-parent", "--reverse", base+".."+head)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("list review candidate commits: %w", err)
	}
	var segments []reviewSegment
	var diffs, prompts []string
	segmentBase, segmentHead, segmentDiff := base, "", ""
	var segmentOmitted []reviewOmittedPath
	appendSegment := func() {
		segments = append(segments, reviewSegment{Base: segmentBase, Head: segmentHead, DiffBytes: len(segmentDiff)})
		diffs = append(diffs, segmentDiff)
		prompts = append(prompts, buildReviewPrompt(segmentBase, segmentHead, segmentDiff, segmentOmitted))
	}
	previous := base
	for _, commit := range strings.Fields(commits) {
		// Keep generated declarations pinned to the candidate base, even when
		// the range starts at an intermediate candidate commit. Every commit
		// must fit on its own, including one a combined diff would shrink.
		ownDiff, ownOmitted, err := reviewScopedDiff(ctx, root, previous, commit, specRoots, git, base)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(ownDiff) > reviewDiffBound {
			return nil, nil, nil, checkReviewDiffBound(scope.diffBytes, scope.omittedPaths)
		}
		diff, omitted := ownDiff, ownOmitted
		if segmentBase != previous {
			diff, omitted, err = reviewScopedDiff(ctx, root, segmentBase, commit, specRoots, git, base)
			if err != nil {
				return nil, nil, nil, err
			}
		}
		if len(diff) > reviewDiffBound {
			appendSegment()
			segmentBase = previous
			diff, omitted = ownDiff, ownOmitted
		}
		segmentHead, segmentDiff, segmentOmitted = commit, diff, omitted
		previous = commit
	}
	if segmentHead != "" {
		appendSegment()
	}
	if len(segments) > reviewSegmentLimit {
		return nil, nil, nil, reviewPrePromptError{fmt.Errorf("review diff too large: %d bytes after omitting %d path(s) needs %d review segments; a review allows at most %d", scope.diffBytes, len(scope.omittedPaths), len(segments), reviewSegmentLimit)}
	}
	if len(segments) == 0 {
		return nil, nil, nil, fmt.Errorf("review candidate has no first-parent commits")
	}
	for i := range prompts {
		prompts[i] = fmt.Sprintf("Review Segment %d of %d: commits %s..%s of the candidate %s..%s; answer for this segment only.\n", i+1, len(segments), segments[i].Base, segments[i].Head, base, head) + prompts[i]
	}
	return segments, diffs, prompts, nil
}

func classifyReviewSegmentAnswers(record reviewRecord, answers []reviewSegmentAnswer) (reviewRecord, int) {
	var findings []reviewFinding
	var texts []string
	permissionRefused := false
	for i, answer := range answers {
		classified, code := classifyReviewCommandResult(record, answer.result, answer.err)
		if classified.Outcome == reviewOutcomeFindings || classified.Outcome == reviewOutcomeReviewed {
			classified, code = validateReviewFindingAnchors(classified, answer.diff)
		}
		permissionRefused = permissionRefused || classified.PermissionRefused
		if classified.Outcome == reviewOutcomeBlocked {
			classified.PermissionRefused = permissionRefused
			return classified, code
		}
		for _, finding := range classified.FindingItems {
			finding.ID = fmt.Sprintf("F%d", len(findings)+1)
			finding.Segment = i + 1
			findings = append(findings, finding)
			texts = append(texts, "- "+finding.Text)
		}
	}
	record.PermissionRefused = permissionRefused
	record.Validation = &reviewValidation{Conventions: deliveryConventionsVersion, Validator: "not-needed"}
	record.FindingItems, record.Findings = findings, strings.Join(texts, "\n")
	if len(findings) == 0 {
		record.Outcome = reviewOutcomeReviewed
		return record, exitOK
	}
	record.Outcome = reviewOutcomeFindingsDismissed
	for _, finding := range findings {
		if !reviewFindingDismissedByValidation(finding) {
			record.Outcome = reviewOutcomeFindings
			return record, exitRunFailed
		}
	}
	return record, exitOK
}
