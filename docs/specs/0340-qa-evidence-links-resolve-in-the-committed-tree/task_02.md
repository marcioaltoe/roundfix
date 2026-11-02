---
task: task_02
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: pending
type: backend
complexity: medium
---

# Task 02: The Daemon gives the QA gate one repair turn and refuses an unresolved evidence link at the commit

## Overview

The Daemon's QA settlement reads the shared check from task_01. When the QA
turn ends with an Unresolved Evidence Link in the newest QA Report, the
Daemon returns the same Agent Session once with each link and reason. At the
QA Report commit boundary it checks again, and a remaining link settles the
QA Task `failed`, publishes one `evidence_links` QA Run Event and creates no
QA Report commit, as a Secret Match and a non-image binary do. The implement
guide and the Roundfix implement reference state the rule. This answers the
adopted Argus report of 2026-11-01
([references](references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md)),
whose report settled `completed` in a Daemon Run.

## Requirements

1. MUST add, in `runQAGate` (`internal/daemon/task_engine.go`), inside the
   branch that ran the QA Agent and after its background-task check, the
   repair turn of `_techspec.md` → API Contract 6: `spec.UnresolvedEvidenceLinks`
   on the newest report; when it returns links, the `evidence_links` event
   with outcome `repair` and its `links` payload, the progress line
   `QA evidence links unresolved: <n>; returning the QA turn once.`, and one
   turn through `runAgentSession` in the same Agent Session with the prompt a
   new `qaEvidenceLinkRepairPrompt` builds. Its first line is
   `QA evidence links do not resolve in the committed tree.`; it names the
   report path and, per link, destination, line and reason, and carries the
   three instructions API Contract 6 lists. Stop, deadline and Agent errors
   follow the QA turn's handling. No second repair turn runs.
2. MUST add the commit-boundary refusal of API Contract 7: after the Secret
   Scan and the non-image binary check find nothing, a remaining link sets the
   QA status `failed`, `accepted` false and the reason to
   `spec.EvidenceLinksError{Links: links}.Error()`, publishes the
   `evidence_links` event with outcome `refused`, prints
   `QA Report refused: <n> unresolved evidence link(s)`, and returns before
   `commitPreparedQAReport`. A mechanically closed report gets the check
   without a repair turn; a read or Git error is returned as an
   infrastructure failure.
3. MUST create `internal/daemon/qa_evidence_links_test.go` with a suite
   header, a wrapper runner over `taskFakeRunner` that recognises the repair
   prompt by its first line, and
   `TestQAEvidenceLinkRepairTurnKeepsTheEvidence`,
   `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed` and
   `TestQAReportWithResolvedLinksRunsNoRepairTurn`, each calling
   `t.Parallel()` first, with the fixtures and assertions `_techspec.md` →
   Testing Approach lists: a committed `.gitignore` holding `*.log`, a `pass`
   report with a Results row linking the evidence, the repair prompt count,
   the event payloads, the exact settled reason, the QA Report commit's
   presence or absence and the files left in the worktree.
4. MUST document, in `docs/user-guide/commands/implement.md` and the canonical
   `.agents/skills/roundfix/references/implement.md`, after the non-image
   binary paragraph, one paragraph stating that the Daemon checks every local
   link of the newest QA Report, that on an unresolved link the Daemon
   returns the QA turn once, naming each link and its reason, and that a
   remaining link settles the QA Task `failed` with
   `QA Report evidence links do not resolve in the committed tree: <link> (line <n>, <reason>)`,
   publishes the `evidence_links` event and creates no QA Report commit. The
   phrase "the Daemon returns the QA turn once" appears verbatim.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`.
6. MUST re-record the Behavior Surface Record as the LAST step, after the
   last guide edit, with `go test -count=1 -tags docscontract
   ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
   -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
7. MUST NOT change `settleQAVerdict`, `qaSettlementReportEligibility`,
   `prepareQAReportCommit`, `scanQAReportPaths`, `nonImageQAReportEvidence`,
   `repairTaskVerification`, `internal/spec/`, `internal/cli/`, the `qa-gate`
   skill, any `### QA settlement` section, any Baseline module, `docs/agents/`,
   the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: a Daemon QA turn that writes a `pass` report linking an ignored `.log` is accepted today and its commit lacks the file.
- [ ] Add the repair turn and its prompt, then the commit-boundary refusal.
- [ ] Add the three Daemon tests.
- [ ] Update the implement guide and reference, sync the mirror, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] An ignored evidence link gets exactly one repair turn; a repaired report is accepted and its commit holds the kept evidence.
- [ ] A link still unresolved after the repair settles the QA Task `failed` with the shared reason, one `refused` event, and no QA Report commit.
- [ ] The implement guide and both reference copies carry the rule, the mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- creates: `internal/daemon/qa_evidence_links_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Daemon's QA settlement and its `evidence_links` QA Run Event — consumers: `./internal/daemon`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestQAEvidenceLinkRepairTurnKeepsTheEvidence|TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed|TestQAReportWithResolvedLinksRunsNoRepairTurn)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAEvidenceLinkRepairTurnKeepsTheEvidence TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed TestQAReportWithResolvedLinksRunsNoRepairTurn; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/daemon` — expected: exit 0. Before this Task none of the three tests exists, so the first pass line is missing; afterwards they pass with the whole Daemon suite.
- `grep -qF 'qaEvidenceLinkRepairPrompt(' internal/daemon/task_engine.go || { echo 'the Daemon has no QA evidence link repair turn' >&2; exit 1; }; grep -qF 'spec.EvidenceLinksError{' internal/daemon/task_engine.go || { echo 'the Daemon does not refuse with the shared text' >&2; exit 1; }; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'the Daemon returns the QA turn once' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the Daemon has no repair prompt, so the first check fails. Afterwards the Daemon reads the shared check, the guide and both reference copies carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 3; Core Feature 4; Success Metric 4; Success Metric 5; Success Metric 7
- `_techspec.md` → Current behavior 1-2; API Contracts 6-8; Invariants 1, 4 and 5; Testing Approach; Risks & Considerations; Build Order 2
