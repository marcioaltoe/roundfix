package cli

// Suite: qualifying partial archives through the command entry point.
// Invariant: the new record's unproven links pass the enabled strict history rule.
// Boundary IN: committed temporary main repository, declarations and --promote.
// Boundary OUT: exit codes, the whole unproven list and strict-check diagnostics.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck(t *testing.T) {
	t.Parallel()
	_, root := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Implementation", status: string(spec.StatusCompleted)}})
	if err := os.RemoveAll(filepath.Join(root, "docs/specs", implementTestSlug)); err != nil {
		t.Fatal(err)
	}
	specsRoot := filepath.Join(root, "docs/specs")
	writeImplementSpecAtRoot(t, specsRoot, "0001-demo", []implementSeed{{id: "task_01", title: "Implementation", status: string(spec.StatusCompleted)}})
	writeImplementSpecAtRoot(t, specsRoot, "0002-other", []implementSeed{{id: "task_01", title: "Other implementation", status: string(spec.StatusPending)}})
	folder := filepath.Join(specsRoot, "0001-demo")
	prd := filepath.Join(folder, "_prd.md")
	mustWrite(t, prd, mustRead(t, prd)+"\n## Unreachable Acceptance\n"+
		"\n- criterion: the alert record\n  reason: it needs a live cluster\n  satisfied-by: [`references/removed.md`](references/removed.md)\n"+
		"\n- criterion: the seven-day measurement\n  reason: it needs seven days of production traffic\n  satisfied-by: [`references/promoted.md`](references/promoted.md)\n")
	otherTask := filepath.Join(specsRoot, "0002-other/task_01.md")
	mustWrite(t, otherTask, mustRead(t, otherTask)+"\n## Consumers\n\nNone. This fixture Task introduces no shared contract.\n")
	mustMkdir(t, filepath.Join(folder, "references"))
	mustWrite(t, filepath.Join(folder, "references/removed.md"), "Alert record\n")
	mustWrite(t, filepath.Join(folder, "references/promoted.md"), "Measurement procedure\n")
	mustMkdir(t, filepath.Join(folder, "qa"))
	mustWrite(t, filepath.Join(folder, "qa/qa-report-2026-10-21.md"), "---\nverdict: partial\nrows_total: 3\nrows_passed: 1\nrows_blocked: 2\nrows_blocked_declared: 2\nrows_blocked_finding: 0\nrows_blocked_environment: 0\n---\n\n# QA\n\n## Outcome\n\nImplementation verified; production evidence remains declared.\n")
	mustMkdir(t, filepath.Join(root, "docs/agents"))
	mustWrite(t, filepath.Join(root, "docs/agents/agent-instructions.md"), "# Agent instructions\n\nVerification: go test ./...\n")
	mustWrite(t, filepath.Join(root, "docs/agents/domain.md"), "# Domain\n\nIdentifier strategy: not applicable.\n")
	mustWrite(t, filepath.Join(root, "docs/agents/docs-layout.md"), "<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->\n"+spec.RetirementRuleSentinel+"\n<!-- setup-context-driven:end id=guide.docs-layout -->\n")
	mustWrite(t, filepath.Join(root, "docs/agents/setup-context.json"), `{"managedArtifacts":[{"id":"guide.docs-layout","path":"docs/agents/docs-layout.md"}]}`)
	commitArchiveFixture(t)
	revision := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"archive", "0001-demo", "--promote", "references/promoted.md"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("archive exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
	record := archiveRecordAt(t, filepath.Join(root, "docs/history/specs/0001-demo.md"))
	want := []string{
		"`references/removed.md` (`docs/specs/0001-demo/references/removed.md@" + revision + "`)",
		"[`references/promoted.md`](../../references/promoted.md)",
	}
	if !reflect.DeepEqual(record.Unproven, want) {
		t.Errorf("unproven=%q want=%q", record.Unproven, want)
	}
	gittest.Run(t, root, "add", "-A")
	gittest.Run(t, root, "commit", "-m", "docs: archive demo")
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "check", "--strict", "0002-other"}, &stdout, &stderr)
	if strings.Contains(stdout.String()+stderr.String(), "SC-HISTORY-RETAINED") {
		t.Fatalf("history check reports retained links: stdout=%s stderr=%s", &stdout, &stderr)
	}
	if code != exitOK {
		t.Fatalf("strict check exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
}
