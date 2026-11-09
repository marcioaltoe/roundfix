---
task: task_02
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
status: pending
type: backend
complexity: medium
---

# Task 02: The Daemon names the Run Worktree, gives one repair turn and refuses Borrowed Evidence at the commit

## Overview

The Daemon's QA settlement joins the rule task_01 added. The QA prompt names
the Run Worktree, the audited head and the pass's evidence directory; after
the QA turn the Daemon reads the newest report for Borrowed Evidence, with the
QA Report Import's copies and the Run Branch it knows, returns the turn once
alongside the evidence-link repair, and at the QA Report commit boundary
settles a report that still borrows `failed` without a commit. This answers
the message reported by quaestor on 2026-11-06, whose Run parked
`run-unresolved` on another Run's evidence
([references](references/2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md)).

## Requirements

1. MUST append, in `runQAGate` (`internal/daemon/task_engine.go`), right
   after the `Seeded QA Report:` lines, the four prompt lines of
   `_techspec.md` → API Contract 8, verbatim:
   `Run Worktree: <plan.WorkDir>`, `Audited head: <auditedHead>`,
   `Evidence directory: <dir>/` and
   `Build and run every executable a row exercises inside this Run Worktree, list their paths in user_flow_executables, write this pass's evidence under the evidence directory, and never copy, replay or cite evidence or an executable an earlier Run left.`,
   where `<dir>` is the seeded report's directory joined with `evidence`
   and the report's file name without `.md`, in slash form.
2. MUST add `borrowedQAEvidence` (the newest report through
   `spec.NewestQAReport`, none on `spec.ErrNoQAReport`, then
   `spec.FindBorrowedEvidence` with `RunBranchPrefix: store.RunBranchPrefix`,
   `RunBranch: plan.RunWorktree.Branch`, `AuditedHead: auditedHead` and,
   when a prior pass was imported, `ImportCommit: prior.Commit` and
   `ImportedPaths: prior.Files`) and `qaBorrowedEvidenceRepairPrompt`, whose
   lines are exactly those of API Contract 6.
3. MUST make the repair point of `runQAGate` (where the evidence-link repair
   runs today) compute the unresolved links and the Borrowed Evidence
   together and, when either exists, run exactly one repair turn through
   `runAgentSession` in the same Agent Session. Its prompt is
   `qaEvidenceLinkRepairPrompt` unchanged when links exist, followed by
   `qaBorrowedEvidenceRepairPrompt` when Borrowed Evidence exists. With
   Borrowed Evidence it publishes the `evidence_binding` QA Run Event
   (`phase` `evidence_binding`, `outcome` `repair`, `items` of `subject`,
   `line`, `reason`) and prints
   `QA evidence does not belong to this Run: <n> item(s); returning the QA turn once.`
   The link event and line keep their bytes. Stop, deadline and Agent errors
   follow the QA turn's handling; no second repair turn runs.
4. MUST add the commit-boundary refusal of API Contract 7: after the Secret
   Scan, the non-image binary check and the link check find nothing,
   remaining Borrowed Evidence sets the QA status `failed`, `accepted` false
   and the reason `spec.BorrowedEvidenceError{Evidence: items}.Error()`,
   publishes `evidence_binding` with outcome `refused`, prints
   `QA Report refused: <n> borrowed evidence item(s)`, and returns before
   `commitPreparedQAReport`. A read or Git error is returned as an
   infrastructure failure.
5. MUST create `internal/daemon/qa_evidence_binding_test.go` with a suite
   header, a wrapper runner over `taskFakeRunner` that recognises the repair
   prompt by its first line `QA evidence does not belong to this Run.`, and
   these tests, each calling `t.Parallel()` first:
   - `TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence`: the fixture
     commits `qa/evidence/2026-11-06-gate/head.txt` under the Spec before the
     cycle; the QA report (`pass`) links it from a Results row; the repair
     writes new evidence under the evidence directory and relinks it; the QA
     Task settles `completed`, the QA Report commit holds the new file, one
     repair prompt names `Run Worktree:`, `Audited head: <HEAD>` and
     `Evidence directory:`, and the `evidence_binding` outcomes are
     `["repair"]`;
   - `TestQABorrowedEvidenceAfterRepairSettlesFailed`: no repair; the QA Task
     settles `failed` with the reason
     `QA Report evidence does not belong to this Run: evidence/2026-11-06-gate/head.txt (line <n>, reused)`,
     the outcomes are `["repair", "refused"]`, no QA Report commit or commit
     event exists, and the report and evidence stay in the worktree;
   - `TestQABorrowedEvidenceNamesARerunRowCitingImportedEvidence`: a prior
     pass recorded with `recordPriorCommit` (as
     `TestQAGateCarriesARowFromAnUnintegratedFailedPass` builds one) is
     imported; the new report's `pass` row links the imported `.bin`
     evidence unchanged; the reason names it with `imported`;
   - `TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory`: the
     first QA prompt holds the four lines of requirement 1 in order.
   `TestQAEvidenceLinkRepairTurnKeepsTheEvidence`,
   `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed`,
   `TestQAReportWithResolvedLinksRunsNoRepairTurn`,
   `TestQAGateCarriesARowFromAnUnintegratedFailedPass`,
   `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte` and every
   other test of the package that pins QA prompt text on the default branch
   when the Task runs keep passing unchanged.
6. MUST document, in `docs/user-guide/commands/implement.md` and the
   canonical `.agents/skills/roundfix/references/implement.md`, after the
   evidence-link paragraph, one paragraph stating that the QA prompt names
   the Run Worktree, the audited head and the evidence directory
   `qa/evidence/<report name>/`; that the Daemon then checks the newest QA
   Report for Borrowed Evidence, counting the QA Report Import's unchanged
   copies as `imported`; that "On Borrowed Evidence the Daemon returns the QA
   turn once", in the same turn as an unresolved link; and that a remaining
   item settles the QA Task `failed` with
   `QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>)`,
   publishes the `evidence_binding` event and creates no QA Report commit.
7. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`,
   recording `skills/testdata/owned-skill-versions.json`.
8. MUST re-record the Behavior Surface Record as the LAST step, after the
   last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and keep `docs/references/behavior-surfaces.json`.
9. MUST NOT change `settleQAVerdict`, `qaSettlementReportEligibility`,
   `prepareQAReportCommit`, `importPriorQAPass`, `copyPriorQAPass`,
   `mechanicalQAReportContent`, `recordQAEvidenceSnapshots`,
   `qaEvidenceLinkRepairPrompt`, `internal/spec/`, `internal/cli/`,
   `internal/agent/`, the `qa-gate` skill, any `### QA settlement` section,
   any Baseline module, `docs/agents/`, the Makefile, CI workflows, `go.mod`
   or `go.sum`.

## Subtasks

- [ ] Append the prompt facts.
- [ ] Add the shared repair turn and the commit-boundary refusal.
- [ ] Add the four Daemon tests.
- [ ] Update the implement guide and reference, sync, record the skill version, re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The QA prompt names the Run Worktree, the audited head and the evidence directory.
- [ ] Borrowed Evidence gets one repair turn; a repaired report settles `completed` with its commit, and an unrepaired one settles `failed` with the shared text and no QA Report commit.
- [ ] A re-run row citing the QA Report Import's unchanged copy is named `imported`, and the evidence-link and carry-forward behavior keeps its bytes.

## Context

- instruction: `docs/adr/0376-a-qa-report-relies-only-on-evidence-its-own-run-produced.md`
- instruction: `internal/daemon/qa_evidence_links_test.go`
- instruction: `internal/daemon/qa_prior_pass_test.go`
- instruction: `internal/daemon/qa_prior_pass.go`
- creates: `internal/daemon/qa_evidence_binding_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Daemon's QA prompt facts, its repair turn and its `evidence_binding` QA Run Event — consumers: `./internal/daemon`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence|TestQABorrowedEvidenceAfterRepairSettlesFailed|TestQABorrowedEvidenceNamesARerunRowCitingImportedEvidence|TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence TestQABorrowedEvidenceAfterRepairSettlesFailed TestQABorrowedEvidenceNamesARerunRowCitingImportedEvidence TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/daemon` — expected: exit 0. Before this Task none of the four tests exists, so the first pass line is missing; afterwards they pass with the whole Daemon suite, the evidence-link and prior-pass tests among them.
- `grep -qF 'qaBorrowedEvidenceRepairPrompt(' internal/daemon/task_engine.go || { echo 'the Daemon has no Borrowed Evidence repair turn' >&2; exit 1; }; grep -qF 'spec.BorrowedEvidenceError{' internal/daemon/task_engine.go || { echo 'the Daemon does not refuse with the shared text' >&2; exit 1; }; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'On Borrowed Evidence the Daemon returns the QA turn once' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the Daemon has no Borrowed Evidence repair prompt, so the first check fails. Afterwards the Daemon reads the shared check, the guide and both reference copies carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Core Feature 4; Success Metric 6; Success Metric 7; Success Metric 8
- `_techspec.md` → Current behavior 3; API Contracts 6-8; Invariants 2, 3 and 5; Testing Approach; Risks & Considerations; Build Order 2
- ADR-0376, ADR-0351
