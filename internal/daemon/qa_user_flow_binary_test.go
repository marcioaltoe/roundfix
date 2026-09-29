// Suite: QA user-flow binary settlement.
// Invariant: the Daemon owns seeded auditor fields and, only in a self-audit, accepts public-CLI evidence from the audited head's build.
// Boundary IN: the seeded QA Report, injected Auditing Binary, audited Git head, Agent-authored report, prompt, and QA Task settlement.
// Boundary OUT: archive and standalone acceptance commands, which continue to use spec.QAReportEligibility.
package daemon

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/app"
	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const qaForeignBuildCommit = "0123456789abcdef0123456789abcdef01234567"

type qaUserFlowBinaryFixture struct {
	qa            *qaAuditorStalenessFixture
	binary        app.AuditingBinary
	auditingLine  string
	stalenessLine string
	auditedHead   string
}

func newQAUserFlowBinaryFixture(t *testing.T, selfAudit bool) *qaUserFlowBinaryFixture {
	t.Helper()
	qa := newQAAuditorStalenessFixture(t, 1)
	commit := qaForeignBuildCommit
	ancestry := app.AncestryUnknown
	if selfAudit {
		commit = qa.deliveryBase
		ancestry = app.AncestryNotOlder
	}
	binary := app.AuditingBinary{Version: "0.17.0", Commit: commit, Built: "2026-09-28T12:00:00Z"}
	return &qaUserFlowBinaryFixture{
		qa:            qa,
		binary:        binary,
		auditingLine:  binary.String(),
		stalenessLine: binary.StalenessLine("", ancestry),
		auditedHead:   strings.TrimSpace(gittest.Run(t, qa.taskFixture.gitRoot, "rev-parse", "HEAD")),
	}
}

func (fixture *qaUserFlowBinaryFixture) run(t *testing.T, report string) (*taskFakeRunner, TaskCycleResult) {
	t.Helper()
	runner := &taskFakeRunner{
		calls:    fixture.qa.taskFixture.calls,
		gitRoot:  fixture.qa.taskFixture.gitRoot,
		qaReport: report,
	}
	engine := fixture.qa.engine(t, runner, &qaGateRecordingVerifier{}, fixture.binary, speccheck.MechanicalResult{})
	result, err := engine.TaskCycle(context.Background(), fixture.qa.plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	return runner, result
}

func qaUserFlowPassReport(auditingBinary, auditorStaleness, userFlowBinary string) string {
	var userFlowLine string
	if userFlowBinary != "" {
		userFlowLine = "user_flow_binary: " + strconv.Quote(userFlowBinary) + "\n"
	}
	return "---\n" +
		"verdict: pass\n" +
		"auditing_binary: " + strconv.Quote(auditingBinary) + "\n" +
		"auditor_staleness: " + strconv.Quote(auditorStaleness) + "\n" +
		userFlowLine +
		"rows_blocked_environment: 0\n" +
		"rows_blocked_finding: 0\n" +
		"rows_blocked_declared: 0\n" +
		"---\n\n# QA Report\n\n## Results\n\n" +
		"| # | Status | Provenance |\n| - | --- | --- |\n| 1 | pass | observed |\n"
}

func (fixture *qaUserFlowBinaryFixture) assertSettlement(t *testing.T, want spec.Status, cause string) {
	t.Helper()
	qaTaskID := fixture.qa.taskFixture.graph.QATaskID
	if got := taskStatusOnDisk(t, fixture.qa.taskFixture.gitRoot, qaTaskID); got != string(want) {
		t.Fatalf("QA Task status = %s, want %s", got, want)
	}
	reason := ""
	for _, event := range taskEventsOfKind(fixture.qa.taskFixture.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == qaTaskID && payload["phase"] == "settled" {
			reason, _ = payload["reason"].(string)
		}
	}
	if cause == "" {
		if reason != "" {
			t.Fatalf("QA Task settlement reason = %q, want empty", reason)
		}
		return
	}
	if !strings.Contains(reason, cause) {
		t.Fatalf("QA Task settlement reason = %q, want cause containing %q", reason, cause)
	}
}

func TestQASettlementAcceptsTheSeededAuditorFields(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, false)
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, ""))
	fixture.assertSettlement(t, spec.StatusCompleted, "")
}

func TestQASettlementRefusesARewrittenAuditingBinary(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, false)
	fixture.run(t, qaUserFlowPassReport("roundfix 9.9.9 (7654321, built later)", fixture.stalenessLine, ""))
	fixture.assertSettlement(t, spec.StatusFailed, "auditor fields are Daemon-owned")
}

func TestQASettlementRefusesARewrittenAuditorStaleness(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, false)
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, "current: rewritten by the Agent", ""))
	fixture.assertSettlement(t, spec.StatusFailed, "auditor fields are Daemon-owned")
}

func TestSelfAuditSettlementAcceptsAUserFlowBinaryBuiltFromTheAuditedHead(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, true)
	userFlowBinary := "roundfix 0.17.0 (" + fixture.auditedHead[:12] + "-dirty, built 2026-09-28T12:30:00Z)"
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, userFlowBinary))
	fixture.assertSettlement(t, spec.StatusCompleted, "")
}

func TestSelfAuditSettlementRefusesAMissingUserFlowBinary(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, true)
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, ""))
	fixture.assertSettlement(t, spec.StatusFailed, "user_flow_binary")
}

func TestSelfAuditSettlementRefusesAUserFlowBinaryFromAnotherCommit(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, true)
	userFlowBinary := "roundfix 0.17.0 (" + fixture.qa.deliveryBase[:12] + ", built 2026-09-28T12:30:00Z)"
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, userFlowBinary))
	fixture.assertSettlement(t, spec.StatusFailed, "user_flow_binary")
}

func TestSelfAuditSettlementRefusesAUserFlowBinaryWithoutABuildCommit(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, true)
	fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, "roundfix 0.17.0"))
	fixture.assertSettlement(t, spec.StatusFailed, "user_flow_binary")
}

func TestSettlementOutsideASelfAuditIgnoresTheUserFlowBinary(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name           string
		userFlowBinary string
	}{
		{name: "missing"},
		{name: "foreign", userFlowBinary: "roundfix 0.17.0 (fedcba987654321, built elsewhere)"},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newQAUserFlowBinaryFixture(t, false)
			fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, test.userFlowBinary))
			fixture.assertSettlement(t, spec.StatusCompleted, "")
		})
	}
}

func TestSelfAuditQAPromptNamesTheUserFlowBinary(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, true)
	userFlowBinary := "roundfix 0.17.0 (" + fixture.auditedHead[:12] + ", built 2026-09-28T12:30:00Z)"
	runner, _ := fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, userFlowBinary))
	const line = "Self-audit: build roundfix from this Run Worktree with make build, run every public-CLI row with ./bin/roundfix and never a roundfix found on PATH, and record its --version line as user_flow_binary."
	if len(runner.qaPrompts) != 1 || !strings.Contains(runner.qaPrompts[0], line) {
		t.Fatalf("self-audit QA prompt omitted %q:\n%s", line, strings.Join(runner.qaPrompts, "\n"))
	}
}

func TestQAPromptOutsideASelfAuditOmitsTheUserFlowBinary(t *testing.T) {
	t.Parallel()
	fixture := newQAUserFlowBinaryFixture(t, false)
	runner, _ := fixture.run(t, qaUserFlowPassReport(fixture.auditingLine, fixture.stalenessLine, ""))
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("QA prompts = %d, want one", len(runner.qaPrompts))
	}
	if strings.Contains(runner.qaPrompts[0], "Self-audit:") || strings.Contains(runner.qaPrompts[0], "user_flow_binary") {
		t.Fatalf("non-self-audit QA prompt names user_flow_binary:\n%s", runner.qaPrompts[0])
	}
}
