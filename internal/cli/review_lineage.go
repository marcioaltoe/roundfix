package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/preflight"
)

type reviewLineage struct {
	Round                int                        `json:"round"`
	PreviousHead         string                     `json:"previousHead,omitempty"`
	PreviousFindings     []reviewFinding            `json:"previousFindings,omitempty"`
	PreviousDispositions []reviewFindingDisposition `json:"previousDispositions,omitempty"`
	ReviewedHead         string                     `json:"reviewedHead,omitempty"`
	Session              string                     `json:"session,omitempty"`
	Selection            int                        `json:"selection"`
	SessionOpen          bool                       `json:"sessionOpen"`
	ACPSessionIDs        []string                   `json:"acpSessionIds,omitempty"`
	Continued            bool                       `json:"continued"`
}

type reviewLineagePlan struct {
	Lineage   reviewLineage
	Reuse     bool
	Ceiling   bool
	prior     *reviewRecord
	candidate reviewRecord
	ctx       context.Context
}

func decideReviewLineage(ctx context.Context, prior *reviewRecord, candidate reviewRecord, git preflight.GitRunner) (reviewLineagePlan, error) {
	plan := reviewLineagePlan{Lineage: reviewLineage{Round: 1}, candidate: candidate, ctx: ctx}
	if prior == nil || filepath.Clean(prior.Repository) != filepath.Clean(candidate.Repository) || prior.Provider != candidate.Provider || prior.BaseCommit != candidate.BaseCommit {
		return plan, nil
	}
	if prior.HeadCommit != candidate.HeadCommit {
		if _, err := git.RunGit(ctx, candidate.Repository, "merge-base", "--is-ancestor", prior.HeadCommit, candidate.HeadCommit); err != nil {
			// Exit 1 means non-ancestry. Other Git failures must not silently reset the ceiling.
			var exitErr interface{ ExitCode() int }
			if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
				return plan, nil
			}
			return plan, fmt.Errorf("check review lineage ancestry: %w", err)
		}
	}
	plan.prior = prior
	if prior.Lineage != nil {
		plan.Lineage = *prior.Lineage
	}
	if prior.Outcome == reviewOutcomeBlocked {
		return plan, nil
	}
	if prior.HeadCommit == candidate.HeadCommit {
		plan.Reuse = true
		return plan, nil
	}
	if plan.Lineage.Round == 2 || prior.Outcome == reviewOutcomeCeilingClosed {
		plan.Ceiling = true
		return plan, nil
	}
	plan.Lineage = reviewLineage{Round: 2, PreviousHead: prior.HeadCommit, PreviousFindings: prior.FindingItems}
	if prior.Lineage != nil {
		plan.Lineage.ACPSessionIDs = append([]string(nil), prior.Lineage.ACPSessionIDs...)
	}
	return plan, nil
}

func buildRoundTwoPrompt(plan reviewLineagePlan, delta string, previous []reviewFinding, dispositions []reviewFindingDisposition) string {
	header := buildReviewPrompt(plan.candidate.BaseCommit, plan.candidate.HeadCommit, "")
	header, _, _ = strings.Cut(header, "\n\n--- BEGIN CANDIDATE DIFF ---")
	var prompt strings.Builder
	prompt.WriteString(header)
	fmt.Fprintf(&prompt, "\n\nPrevious head: %s\nThis is round 2 of 2 of this review. The diff below is the change since round 1 reviewed %s.\nRound 1 findings:\n", plan.Lineage.PreviousHead, plan.Lineage.PreviousHead)
	for _, finding := range previous {
		validation := "stands"
		if reviewFindingDismissedByValidation(finding) {
			validation = fmt.Sprintf("dismissed-by-validation (%s): %s", finding.Validation.Rule, finding.Validation.Reason)
		}
		disposition := "no disposition"
		for _, d := range dispositions {
			if d.Finding == finding.ID && d.Text == finding.Text {
				switch d.Disposition {
				case "fixed":
					disposition = "fixed by " + d.FixedBy
				case "dismissed":
					disposition = "dismissed: " + d.Evidence
				}
				break
			}
		}
		fmt.Fprintf(&prompt, "- %s | %s | %s | %s\n", finding.ID, validation, disposition, strings.Join(strings.Fields(finding.Text), " "))
	}
	prompt.WriteString("Do not raise again a round-1 finding that was dismissed or fixed unless the change below reintroduces it; raise a standing round-1 finding without a disposition only if the change below leaves it unresolved.\n\n--- BEGIN ROUND 2 DELTA DIFF ---\n")
	prompt.WriteString(delta)
	if !strings.HasSuffix(delta, "\n") {
		prompt.WriteByte('\n')
	}
	prompt.WriteString("--- END ROUND 2 DELTA DIFF ---\n")
	return prompt.String()
}

func reviewDispositionsAtHead(ledger []reviewFindingDisposition, repository, head string) []reviewFindingDisposition {
	result := make([]reviewFindingDisposition, 0)
	for _, d := range ledger {
		if d.Repository == repository && d.HeadCommit == head {
			result = append(result, d)
		}
	}
	return result
}

func closeAtCeiling(plan reviewLineagePlan, ledger []reviewFindingDisposition, git preflight.GitRunner) (reviewRecord, int) {
	record := plan.candidate
	prior := plan.prior
	reviewedHead := prior.HeadCommit
	if prior.Outcome == reviewOutcomeCeilingClosed {
		reviewedHead = prior.Lineage.ReviewedHead
	}
	record.Lineage = &reviewLineage{Round: 2, ReviewedHead: reviewedHead}
	dispositions := reviewDispositionsAtHead(ledger, record.Repository, reviewedHead)
	var matched []reviewFindingDisposition
	var missing []string
	for _, finding := range prior.FindingItems {
		if reviewFindingDismissedByValidation(finding) {
			continue
		}
		matches := 0
		var match reviewFindingDisposition
		for _, d := range dispositions {
			if d.Finding != finding.ID || d.Text != finding.Text {
				continue
			}
			valid := d.Disposition == "dismissed" && strings.TrimSpace(d.Evidence) != "" && d.FixedBy == ""
			if d.Disposition == "fixed" && strings.TrimSpace(d.FixedBy) != "" {
				_, err := git.RunGit(plan.ctx, record.Repository, "merge-base", "--is-ancestor", d.FixedBy, record.HeadCommit)
				valid = err == nil
			}
			if valid {
				matches++
				match = d
			}
		}
		if matches != 1 {
			missing = append(missing, finding.ID)
		} else {
			matched = append(matched, match)
		}
	}
	if len(missing) > 0 {
		record.Outcome = reviewOutcomeBlocked
		record.Reason = fmt.Sprintf("review round ceiling reached: round 2 reviewed %s; dispose finding %s with roundfix review dispose", reviewedHead, strings.Join(missing, ", "))
		return record, exitPreflight
	}
	record.Outcome = reviewOutcomeCeilingClosed
	record.Findings = prior.Findings
	record.FindingItems = prior.FindingItems
	record.Dispositions = matched
	return record, exitOK
}

func validateCeilingClosedReviewRecord(record reviewRecord) error {
	if record.Lineage == nil || record.Lineage.Round != 2 || strings.TrimSpace(record.Lineage.ReviewedHead) == "" || record.Lineage.ReviewedHead == record.HeadCommit {
		return errors.New("ceiling-closed review record requires round 2 and a different reviewed head")
	}
	if strings.TrimSpace(record.Findings) != "" && len(record.FindingItems) == 0 {
		return errors.New("ceiling-closed review record findings require finding items")
	}
	if record.Reason != "" {
		return errors.New("ceiling-closed review record cannot carry a reason")
	}
	standing := 0
	for _, finding := range record.FindingItems {
		if reviewFindingDismissedByValidation(finding) {
			continue
		}
		standing++
		matches := 0
		for _, d := range record.Dispositions {
			if d.Repository != record.Repository || d.HeadCommit != record.Lineage.ReviewedHead || d.Finding != finding.ID || d.Text != finding.Text {
				continue
			}
			if (d.Disposition == "dismissed" && strings.TrimSpace(d.Evidence) != "" && d.FixedBy == "") || (d.Disposition == "fixed" && strings.TrimSpace(d.FixedBy) != "") {
				matches++
			}
		}
		if matches != 1 {
			return fmt.Errorf("ceiling-closed review record finding %q requires one disposition", finding.ID)
		}
	}
	if len(record.Dispositions) != standing {
		return errors.New("ceiling-closed review record requires exactly one disposition per standing finding")
	}
	return nil
}

// End only the session owned by this checkout's record. Closing needs no
// readiness probe or prompt, including when the configured provider changed.
func endOpenReviewSession(ctx context.Context, prior *reviewRecord, runner agent.Runner) error {
	if prior == nil || prior.Lineage == nil || !prior.Lineage.SessionOpen {
		return nil
	}
	runtime, err := agent.RuntimeFor(agent.RuntimeOptions{Agent: prior.Provider})
	if err != nil {
		return err
	}
	return runner.EndSession(context.WithoutCancel(ctx), runtime, agent.SessionRef{Name: prior.Lineage.Session, WorkDir: prior.Repository})
}
