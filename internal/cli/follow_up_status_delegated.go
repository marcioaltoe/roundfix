package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

type followUpDelegatedDecision struct {
	Spec     string `json:"spec"`
	Line     int    `json:"line"`
	Date     string `json:"date"`
	Phase    string `json:"phase"`
	Question string `json:"question"`
	Options  string `json:"options"`
	Adopted  string `json:"adopted"`
	Under    string `json:"under"`
	Answer   string `json:"answer"`
	Problem  string `json:"problem"`
}

func readFollowUpDelegatedDecisions(root string, active []spec.Spec, skipped []spec.SkippedSpec) ([]followUpDelegatedDecision, error) {
	rows := []followUpDelegatedDecision{}
	var readErr error
	for _, item := range active {
		content, err := os.ReadFile(filepath.Join(item.Dir, "_prd.md"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			readErr = errors.Join(readErr, fmt.Errorf("read PRD for Spec %s: %w", item.Slug, err))
			continue
		}
		for _, entry := range speccheck.DelegatedDecisions(content) {
			rows = append(rows, followUpDelegatedDecision{Spec: item.Slug, Line: entry.Line, Date: entry.Date, Phase: entry.Phase, Question: entry.Question, Options: entry.Options, Adopted: entry.Adopted, Under: entry.Under, Answer: entry.Answer, Problem: entry.Problem})
		}
	}
	// Discovery cannot call a Spec active when its PRD is unreadable. Preserve
	// those read errors on the specs source, but a missing PRD is normal authoring.
	for _, item := range skipped {
		slug := filepath.Base(item.Dir)
		if _, err := os.ReadFile(filepath.Join(root, slug, "_prd.md")); err != nil && !errors.Is(err, os.ErrNotExist) {
			readErr = errors.Join(readErr, fmt.Errorf("read PRD for Spec %s: %w", slug, err))
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Spec != rows[j].Spec {
			return rows[i].Spec < rows[j].Spec
		}
		return rows[i].Line < rows[j].Line
	})
	return rows, readErr
}

func renderFollowUpDelegatedDecisions(rows []followUpDelegatedDecision) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Delegated decisions: %d\n", len(rows))
	for _, row := range rows {
		if row.Problem != "" {
			fmt.Fprintf(&out, "  %s line %d: unreadable: %s\n", row.Spec, row.Line, row.Problem)
			continue
		}
		option, _, _ := strings.Cut(row.Adopted, " — ")
		answer, _, _ := strings.Cut(row.Answer, " — ")
		fmt.Fprintf(&out, "  %s %s %s: %s -> %s (answer: %s)\n", row.Spec, row.Date, row.Phase, row.Question, option, answer)
	}
	return out.String()
}
