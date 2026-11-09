---
task: task_01
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
status: pending
type: backend
complexity: medium
---

# Task 01: The shared check refuses Borrowed Evidence in qa-report accept and settle

## Overview

`internal/spec` gains the one check every QA settlement path reads:
`FindBorrowedEvidence` names each item of Borrowed Evidence a QA Report relies
on, with one of four reasons, and `RequireOwnRunEvidence` turns them into one
refusal text. `roundfix qa-report accept` and the Settle Command for a `qa`
Task apply it after the evidence-link rule. The qa-report and settle guides,
the Roundfix settle reference and `CONTEXT.md` state the rule. This answers
the message reported by quaestor on 2026-11-06, "QA terminal reprovou com
evidência de outro Run"
([references](references/2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md)).

## Requirements

1. MUST create `internal/spec/qa_evidence_binding.go` with the reason
   constants, `BorrowedEvidence`, `BorrowedEvidenceOptions`,
   `FindBorrowedEvidence`, `BorrowedEvidenceError` and
   `RequireOwnRunEvidence` exactly as `_techspec.md` → Interfaces and API
   Contracts 1-3 define them: the worktree root and audited head
   derivation, the evidence links under `<spec dir>/qa/evidence/` read
   through `RelativeLinks`, the Carried Row exemption, the four reasons in
   their order with a symbolic-link-free comparison for
   `user_flow_executables`, the other-Run worktree scan over
   `git worktree list --porcelain`, report-line order, nothing checked
   outside Git, an error when `RunBranchPrefix` is empty, and the text
   `QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>)`
   joined by `; `. Git goes through the package's existing reader
   `archiveGitRead`; a read or Git failure is returned as an error, never as
   an item. No new package import.
2. MUST create `internal/spec/qa_evidence_binding_test.go` with
   `TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory`,
   `TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit`,
   `TestBorrowedEvidenceExemptsACarriedRow`,
   `TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree`,
   `TestBorrowedEvidenceNamesAnotherRunWorktree`,
   `TestBorrowedEvidenceNamesImportedEvidence`,
   `TestBorrowedEvidenceErrorNamesEachItem` and
   `TestBorrowedEvidenceOutsideGitChecksNothing`, each calling `t.Parallel()`
   first, with every case `_techspec.md` → Testing Approach lists for them:
   the Quaestor layout with `head-exclusive.execution.txt` and
   `replay.py.txt`, the untracked and the changed tracked file, the
   `.venv/bin/python` symbolic link to a file outside the worktree, the own
   Task worktree on `roundfix/run-own-task_01`, and the non-Run worktree.
   Worktrees and repositories live under `t.TempDir()` and are created with
   `gittest`.
3. MUST make `runQAReportCommand` (`internal/cli/qa_report.go`) call
   `spec.RequireOwnRunEvidence(ctx, reportPath, spec.BorrowedEvidenceOptions{RunBranchPrefix: store.RunBranchPrefix})`
   after `RequireResolvedEvidenceLinks` accepts and print a refusal as
   `roundfix: QA Report eligibility refused: <error>` with exit 1 and empty
   stdout (API Contract 4). Leave `qaReportUsage` unchanged.
4. MUST make `requireSettleQAReportEligibility` (`internal/cli/settle.go`)
   apply the same call on the same report after the link check accepts and
   wrap a refusal as `<error> (report <path>)` (API Contract 5). Change no
   other function of that file.
5. MUST add `TestQAReportAcceptRefusesBorrowedEvidence` and
   `TestSettleRefusesAQAReportWithBorrowedEvidence` to
   `internal/cli/qa_report_test.go`, each calling `t.Parallel()` first, on a
   workspace from `newImplementWorkspace` whose earlier commit holds
   `evidence/2026-11-06-gate/head-exclusive.execution.txt` and
   `evidence/2026-11-06-gate/replay.py.txt` under the Spec's `qa/` directory,
   and an uncommitted report
   `qa/qa-report-2026-11-09.md` with the content
   `---\nverdict: pass\n---\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| 1 | pass | [run](evidence/2026-11-06-gate/head-exclusive.execution.txt) [replay](evidence/2026-11-06-gate/replay.py.txt) |\n`.
   The exact stderr is
   `roundfix: QA Report eligibility refused: QA Report evidence does not belong to this Run: evidence/2026-11-06-gate/head-exclusive.execution.txt (line 9, reused); evidence/2026-11-06-gate/replay.py.txt (line 9, reused)\n`
   for `qa-report accept`, and
   `Settle surface: <root>\nroundfix: settle QA Report is ineligible: QA Report evidence does not belong to this Run: evidence/2026-11-06-gate/head-exclusive.execution.txt (line 9, reused); evidence/2026-11-06-gate/replay.py.txt (line 9, reused) (report <report path relative to root>)\n`
   for settle, which leaves the Task file unchanged. Each test then writes
   new evidence under `evidence/qa-report-2026-11-09/`, updates both links,
   and asserts exit 0 (and, for settle, the Task `completed`).
   `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink` and
   `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink` keep passing
   unchanged.
6. MUST document the rule in `docs/user-guide/commands/qa-report.md` and in
   the `## QA Report acceptance` section of the canonical
   `.agents/skills/roundfix/references/settle.md`, after the evidence-link
   paragraphs, as one paragraph that begins with the sentence "Every
   evidence file the report relies on must belong to this Run." and states
   the four reasons, the audited head, the Carried Row exemption, that
   outside Git nothing is checked, and the refusal line in a `text` block:
   `roundfix: QA Report eligibility refused: QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>)`.
   In `docs/user-guide/commands/settle.md`, after the evidence-link
   paragraph and its block, state that for a `qa` Task the Settle Command
   applies the same rule after the link rule accepts, with the line
   `roundfix: settle QA Report is ineligible: QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>) (report <path>)`.
   Change no `### QA settlement` section.
7. MUST add **Borrowed Evidence** to `CONTEXT.md` through `domain-modeling`,
   after **Unresolved Evidence Link**, beginning "Evidence a QA Report relies
   on that its own Run did not produce", naming the four reasons and the
   Carried Row exemption, citing ADR-0376, with an `_Avoid_:` line; and
   append to **QA Report**, after its sentence on local links, the sentence
   "It lists the executables its user-flow rows ran as `user_flow_executables`, and it settles only when it relies on no Borrowed Evidence."
   citing ADR-0376.
8. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`,
   recording `skills/testdata/owned-skill-versions.json`. Never hand-edit a
   mirror or a version field.
9. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of any guide, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and keep `docs/references/behavior-surfaces.json`.
10. MUST NOT touch `internal/spec/qa.go`, `internal/spec/qa_evidence_links.go`,
    `internal/spec/archive.go`, `internal/speccheck/`, `internal/daemon/`,
    `internal/agent/`, the `qa-gate` skill, any Baseline module,
    `docs/agents/`, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: an uncommitted `pass` report linking evidence an earlier commit holds is accepted by `qa-report accept` today.
- [ ] Add the shared check and its tests.
- [ ] Wire `qa-report accept` and the Settle Command, with their tests.
- [ ] Update the guides, the reference and `CONTEXT.md`, sync the mirror, record the skill version, re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] `qa-report accept` and the Settle Command refuse the Quaestor layout with the shared text and accept the same report once its evidence is regenerated under the pass's own directory.
- [ ] Carried Rows, Python virtual-environment links inside the worktree, the own Run's worktrees and reports outside Git get no item.
- [ ] The guides, both reference copies and `CONTEXT.md` carry the rule, the mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0376-a-qa-report-relies-only-on-evidence-its-own-run-produced.md`
- instruction: `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/spec/qa_evidence_links.go`
- instruction: `internal/spec/record_links.go`
- creates: `internal/spec/qa_evidence_binding.go`
- creates: `internal/spec/qa_evidence_binding_test.go`
- interface: `internal/cli/qa_report.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/qa_report_test.go`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/settle.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `spec.FindBorrowedEvidence`, `spec.RequireOwnRunEvidence` and `spec.BorrowedEvidenceError` — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory|TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit|TestBorrowedEvidenceExemptsACarriedRow|TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree|TestBorrowedEvidenceNamesAnotherRunWorktree|TestBorrowedEvidenceNamesImportedEvidence|TestBorrowedEvidenceErrorNamesEachItem|TestBorrowedEvidenceOutsideGitChecksNothing)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit TestBorrowedEvidenceExemptsACarriedRow TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree TestBorrowedEvidenceNamesAnotherRunWorktree TestBorrowedEvidenceNamesImportedEvidence TestBorrowedEvidenceErrorNamesEachItem TestBorrowedEvidenceOutsideGitChecksNothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestQAReportAcceptRefusesBorrowedEvidence|TestSettleRefusesAQAReportWithBorrowedEvidence|TestQAReportAcceptRefusesAnUnresolvedEvidenceLink|TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportAcceptRefusesBorrowedEvidence TestSettleRefusesAQAReportWithBorrowedEvidence TestQAReportAcceptRefusesAnUnresolvedEvidenceLink TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec && go test -count=1 -timeout 20m -run 'QAReport|Settle|QAPartial|PartialShape' ./internal/cli` — expected: exit 0. Before this Task none of the ten new tests exists, so the first pass line is missing; afterwards they pass with the existing evidence-link tests, the Spec suite and the QA Report and Settle tests.
- `grep -qF 'RequireOwnRunEvidence(' internal/cli/qa_report.go || { echo 'qa-report accept does not read the shared check' >&2; exit 1; }; grep -qF 'RequireOwnRunEvidence(' internal/cli/settle.go || { echo 'settle does not read the shared check' >&2; exit 1; }; for f in docs/user-guide/commands/qa-report.md .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Every evidence file the report relies on must belong to this Run.' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/settle.md | grep -qF -- 'QA Report evidence does not belong to this Run' || { echo 'missing settle refusal in the settle guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Borrowed Evidence**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'it settles only when it relies on no Borrowed Evidence.' || { echo 'missing QA Report sentence in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**QA Report**:' || { echo 'missing QA Report term in CONTEXT.md' >&2; exit 1; }; cmp .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither CLI path reads the check, so the first check fails. Afterwards both read it, the guides, both reference copies and the glossary carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint and the glossary phrase test among them, pass.

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1, 2 and 4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 8; Glossary
- `_techspec.md` → Current behavior 1-2 and 5; API Contracts 1-5; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1
- ADR-0376, ADR-0351
