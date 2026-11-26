//go:build docscontract

// Suite: approved gate time targets.
// Invariant: target constants and the reference table agree with the validated document.
// Boundary IN: the live target JSON and reference document.
// Boundary OUT: measurement collection and elapsed-time verdicts.
//
//verify:always
package docscontract

import (
	"os"
	"path/filepath"
	"roundfix/internal/gatetime"
	"strings"
	"testing"
)

func TestTheGateTimeTargetsHoldTheApprovedTargets(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..")
	doc, err := gatetime.Load(filepath.Join(root, "docs/references/gate-time-targets.json"))
	if err != nil {
		t.Fatal(err)
	}
	approved := map[string]gatetime.Limits{"task-settlement": {MedianSeconds: 120, P95Seconds: 240}, "complete-suite-run": {Seconds: 360}, "before-pull-request": {Seconds: 90}, "pull-request": {MedianSeconds: 480, P95Seconds: 600}, "test-binary": {Seconds: 300}}
	if len(doc.Gates) > 0 {
		lab := doc.Gates[0].Lab
		if lab == nil || lab.Seconds != 45 || lab.RegenerationSeconds != 150 {
			t.Fatal("changed settlement lab targets", lab)
		}
	}
	for _, err := range gatetime.Validate(doc, approved) {
		t.Error(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "docs/references/gate-time-targets.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{"| Task settlement | the mini | median 120 s, p95 240 s |", "| Task settlement lab | the mini | 45 s, 150 s with regeneration |", "| Complete suite run | the mini | 360 s |", "| Before pull request | the mini | 90 s |", "| Pull request | standard macOS runner | median 480 s, p95 600 s |", "| Test binary | standard macOS runner | 300 s |"} {
		if !strings.Contains(string(data), row) {
			t.Errorf("missing target row %s", row)
		}
	}
}
