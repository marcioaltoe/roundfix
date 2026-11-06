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
	Base                  string `json:"base"`
	Head                  string `json:"head"`
	MeasuredFrom          string `json:"measuredFrom,omitempty"`
	DiffBytes             int    `json:"diffBytes"`
	EstimatedPromptTokens int    `json:"estimatedPromptTokens,omitempty"`
}

type reviewSegmentFit func(index int, prompt string, diff string) bool

type reviewSegmentAnswer struct {
	result agent.ExecuteResult
	err    error
	diff   string
}

func planReviewSegments(ctx context.Context, root, base, head string, specRoots []string, git preflight.GitRunner, scope reviewSpecContextResult, fit reviewSegmentFit) ([]reviewSegment, []string, []string, error) {
	commits, err := git.RunGit(ctx, root, "rev-list", "--first-parent", "--reverse", base+".."+head)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("list review candidate commits: %w", err)
	}
	var segments []reviewSegment
	var diffs, prompts []string
	segmentBase, segmentHead, segmentDiff := base, "", ""
	segmentMeasuredFrom := base
	var segmentOmitted []reviewOmittedPath
	appendSegment := func() {
		segment := reviewSegment{Base: segmentBase, Head: segmentHead, DiffBytes: len(segmentDiff)}
		if segmentMeasuredFrom != segmentBase {
			segment.MeasuredFrom = segmentMeasuredFrom
		}
		segments = append(segments, segment)
		diffs = append(diffs, segmentDiff)
		prompts = append(prompts, buildReviewPrompt(segmentBase, segmentHead, segmentDiff, segmentOmitted))
	}
	segmentPrompt := func(index int, segmentBase, segmentHead, measuredFrom, diff string, omitted []reviewOmittedPath) string {
		// Any admitted total is one digit, so "of 4" has the sent header's length.
		return reviewSegmentHeader(index, reviewSegmentLimit, segmentBase, segmentHead, base, head, measuredFrom) + buildReviewPrompt(segmentBase, segmentHead, diff, omitted)
	}
	previous := base
	for _, commit := range strings.Fields(commits) {
		// Keep generated declarations pinned to the candidate base, even when
		// the range starts at an intermediate candidate commit. Every commit
		// must fit on its own, including one a combined diff would shrink.
		ownMeasuredFrom, _, err := reviewMeasuredFrom(ctx, git, root, base, previous, commit)
		if err != nil {
			return nil, nil, nil, err
		}
		ownDiff, ownOmitted, err := reviewScopedDiff(ctx, root, ownMeasuredFrom, commit, specRoots, git, base)
		if err != nil {
			return nil, nil, nil, err
		}
		ownIndex := len(segments) + 1
		if segmentHead != "" {
			// A later commit can start a new segment without the first one's context.
			ownIndex++
		}
		if !fit(ownIndex, segmentPrompt(ownIndex, previous, commit, ownMeasuredFrom, ownDiff, ownOmitted), ownDiff) {
			if scope.estimatedPromptTokens > 0 {
				return nil, nil, nil, reviewSizeRefusalForDiff(ctx, git, root, base, head, scope.omittedPaths, scope.candidateDiff, fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) exceeds the claude review budget of %d tokens, half of its %d-token context window", scope.estimatedPromptTokens, len(scope.omittedPaths), claudeReviewPromptBudget, claudeReviewContextTokens))
			}
			return nil, nil, nil, checkReviewDiffBound(ctx, git, root, base, head, scope.candidateDiff, scope.omittedPaths)
		}
		diff, omitted, measuredFrom := ownDiff, ownOmitted, ownMeasuredFrom
		if segmentBase != previous {
			measuredFrom, _, err = reviewMeasuredFrom(ctx, git, root, base, segmentBase, commit)
			if err != nil {
				return nil, nil, nil, err
			}
			diff, omitted, err = reviewScopedDiff(ctx, root, measuredFrom, commit, specRoots, git, base)
			if err != nil {
				return nil, nil, nil, err
			}
		}
		index := len(segments) + 1
		if !fit(index, segmentPrompt(index, segmentBase, commit, measuredFrom, diff, omitted), diff) {
			appendSegment()
			segmentBase = previous
			diff, omitted, measuredFrom = ownDiff, ownOmitted, ownMeasuredFrom
		}
		segmentHead, segmentDiff, segmentOmitted = commit, diff, omitted
		segmentMeasuredFrom = measuredFrom
		previous = commit
	}
	if segmentHead != "" {
		appendSegment()
	}
	if len(segments) > reviewSegmentLimit {
		if scope.estimatedPromptTokens > 0 {
			return nil, nil, nil, reviewSizeRefusalForDiff(ctx, git, root, base, head, scope.omittedPaths, scope.candidateDiff, fmt.Sprintf("review prompt too large: %d estimated tokens after omitting %d path(s) needs %d review segments; a review allows at most %d", scope.estimatedPromptTokens, len(scope.omittedPaths), len(segments), reviewSegmentLimit))
		}
		return nil, nil, nil, reviewSizeRefusalForDiff(ctx, git, root, base, head, scope.omittedPaths, scope.candidateDiff, fmt.Sprintf("review diff too large: %d bytes after omitting %d path(s) needs %d review segments; a review allows at most %d", scope.diffBytes, len(scope.omittedPaths), len(segments), reviewSegmentLimit))
	}
	if len(segments) == 0 {
		return nil, nil, nil, fmt.Errorf("review candidate has no first-parent commits")
	}
	for i := range prompts {
		prompts[i] = reviewSegmentHeader(i+1, len(segments), segments[i].Base, segments[i].Head, base, head, segments[i].MeasuredFrom) + prompts[i]
	}
	return segments, diffs, prompts, nil
}

func reviewSegmentHeader(index, total int, base, head, candidateBase, candidateHead, measuredFrom string) string {
	header := fmt.Sprintf("Review Segment %d of %d: commits %s..%s of the candidate %s..%s; answer for this segment only.\n", index, total, base, head, candidateBase, candidateHead)
	if measuredFrom != "" && measuredFrom != base {
		header += fmt.Sprintf("This segment's diff starts from %s instead of %s, so it leaves out what the default branch changed.\n", measuredFrom, base)
	}
	return header
}

func classifyReviewSegmentAnswers(record reviewRecord, answers []reviewSegmentAnswer, codePromptCount int) (reviewRecord, int) {
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
			if i >= codePromptCount {
				finding.DataSegment = i - codePromptCount + 1
			} else if len(record.Segments) > 0 {
				finding.Segment = i + 1
			}
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
