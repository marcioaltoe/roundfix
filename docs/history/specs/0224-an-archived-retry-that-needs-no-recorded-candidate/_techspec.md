---
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
prd: _prd.md
created: 2026-10-04
---

# An archived retry that needs no recorded candidate — Technical Spec

## Executive Summary

Two predicates decide the 0220 defect the PRD describes. `Engine.Retry` in
`internal/delivery/engine.go` admits an archived item without a candidate head
only when its blocker is `qa-environment-partial`, and
`commandDeliveryWorkflow.qaEnvironmentPartial` in
`internal/cli/deliver_workflow.go` sets the environment partial only when the
environment-blocked rows outnumber the rows waiting for an open Pull Request.
This Spec drops the blocker condition from the first and lowers the second to
"at least one environment-blocked row". No interface, schema, blocker, Park
Class or status line changes. The trade-off accepted is a weaker anchor: the
Implement start head of the item's Run is an ancestor of every commit on the
item branch, so the QA Archive Override, not the ancestry proof, is what
authorizes the retry; review, repository gate, delivery authorization and
required checks still run on the recorded head (ADR-0229).

## Project Constraints

- Identifier strategy: not applicable — no identifier changes; queue items,
  Runs and blockers keep their names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — both changes read local Git, the
  Run Branch's QA Report and the Run Database; no credential, forge read or
  network call is added, and no test reaches GitHub. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0229 (this Spec) decides both
  rules, ADR-0229: "including when every environment-blocked row is one that
  waits for an open Pull Request". ADR-0223 keeps a delivery's work across
  archive, requeue and review, and its descendant rule for an item with a
  candidate is unchanged, ADR-0223: "records that head as a new candidate and
  returns the item to review". ADR-0167 still holds that the pre-PR Pull
  Request row never decides a qualifying partial for archive, ADR-0167: "Every
  other environment-blocked row keeps today's rule". ADR-0154 keeps the
  override a record of user authority, ADR-0154: "Preserve its approval
  source/date and actual QA outcome or absence", and ADR-0165 keeps a review finding
  after archive parked for a corrective Spec, ADR-0165: "an archived Spec is
  never edited to absorb a finding". ADR-0187 and ADR-0189 govern the
  Roundfix Skill edit and its version. ADR-0184 has this TechSpec state its
  changed command surfaces, ADR-0184: "A TechSpec now declares numbered
  Surface Transcripts". The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0104, ADR-0155 and ADR-0156, and ADR-0093, ADR-0117, ADR-0168,
  ADR-0176 and ADR-0183 check this Spec's consistency. ADR-0178, ADR-0179
  and ADR-0166 decide each Task commit's grant and record undeclared paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage
  and row carry, ADR-0182 cites ADR-0117 but settles a Task on the facts its
  gate will check, ADR-0169 cites ADR-0165 but decides the review's merge-base
  diff, ADR-0196 cites ADR-0169 but decides when a review finding parks,
  ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row
  records, when it is observed again and its evidence snapshot, and ADR-0192
  cites ADR-0178 but decides derived-path conflicts; this Spec changes none of
  them, so none applies.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and
  its `SKILL.md` mirror are Governed Paths under the maintainer's skill
  authorization and this Spec's grant of 2026-10-04. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0224-an-archived-retry-that-needs-no-recorded-candidate/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`.

## System Architecture

The Delivery Queue's `Engine` runs each item through its stages and parks it
with a blocker. `runCandidate` parks an unresolved Run as
`qa-environment-partial` when the workflow's `RunResult.QAEnvironmentPartial`
is true and as `run-unresolved` otherwise. `commandDeliveryWorkflow.runResult`
sets that flag from the newest QA Report on the Run Branch, read through
`spec.ReadQAReportFile`.

`Engine.Retry` returns a parked item to a stage. Its branches, in order, are:
a prerequisite park; `corrective-spec-required`; `pull-request-conflict`; an
archived Spec; an active Spec with Task Carry-Forward. In the archived branch
today, `candidateHead` fails for an item without candidate commits, and the
retry refuses with `candidate head is missing` unless the item is an
operator-archived `qa-environment-partial`. For that one park the anchor is
`History.RunStart(item.RunID)`, the Implement start head the Run Database
recorded, and `History.Descends` proves the item head descends from it.

Both components keep their place and seams. Nothing new is added.

## Implementation Design

### Interfaces

No signature changes. The two predicates change as follows; the rest of each
function is unchanged.

```go
// Engine.Retry, archived branch (internal/delivery/engine.go)
operatorArchive := state.QAOverride

// commandDeliveryWorkflow.qaEnvironmentPartial (internal/cli/deliver_workflow.go)
return err == nil && report.Verdict == spec.VerdictPartial &&
	report.RowsBlockedFinding == 0 && report.RowsBlockedEnvironment > 0
```

### The archived retry without a candidate

1. The archived branch is reached only after the prerequisite,
   `corrective-spec-required` and `pull-request-conflict` branches, so those
   parks keep their rules.
2. When the item has no candidate head and its archive records
   `qa_override: true` and the engine has a `History`, the anchor is
   `History.RunStart(ctx, gitRoot, item.RunID)`, whatever the item's blocker.
3. When `History.Descends(ctx, workDir, anchor, head)` is true, the retry
   appends the head to `CandidateCommits`, sets the stage `reviewing` and
   hands the item to the owner, as it does today for the environment park.
4. When the anchor cannot be read, the head does not descend, or the engine
   has no `History`, the retry logs the proof error as today and refuses with
   the existing text, leaving the item unchanged.
5. Without `qa_override: true`, an item with no candidate still refuses with
   `candidate head is missing`.

### The environment-only partial

1. After an unresolved Run, the newest QA Report on the Run Branch sets the
   environment partial when its verdict is `partial`,
   `rows_blocked_finding` is `0` and `rows_blocked_environment` is at least
   `1`.
2. The rows counted as waiting for an open Pull Request
   (`RowsBlockedPrePullRequest`) no longer reduce that count. The field stays
   in `spec.QAReport` for the Archive Command's eligibility (ADR-0167).
3. A missing or unreadable report, a finding-blocked partial, a partial with
   no environment-blocked row, and any other verdict keep `run-unresolved`.

### Data Models

None. The queue item, the Run record and the QA Report keep their fields.

### API Contracts

1. API Contract: Delivery Retry of an operator-archived item without a
   candidate. Input: a parked item with no candidate commits, a recorded Run
   ID, and an archived Spec whose `_prd.md` records `qa_override: true`.
   Output: stage `reviewing`, candidate commits holding exactly the item head,
   blocker cleared, the retry recorded and the owner handed the item, for any
   blocker that reaches the archived branch. Failure: a head that does not
   descend from the Run start, an unreadable Run start or no item history
   refuses with
   `retry Delivery Queue item "<slug>": archived item head "<head>" differs from candidate head ""`
   and leaves the item unchanged; an archive without the override refuses
   with `retry Delivery Queue item "<slug>": candidate head is missing`.
2. API Contract: environment-only partial. An unresolved Run whose newest
   Run-Branch QA Report is `partial` with `rows_blocked_finding: 0` and at
   least one environment-blocked row parks `qa-environment-partial`; every
   other unresolved Run parks `run-unresolved`.

### Surface Transcripts

None. `roundfix deliver retry` and `roundfix deliver status` print the lines
they print today; this Spec changes which retries succeed and which blocker a
park records, not any output text. Reproducing a successful retry through the
built binary would start a detached queue owner that runs the review, so
API Contracts 1 and 2 are asserted through the engine and the workflow with
real Git, a real archive and a disposable Run Database.

## Vocabulary Contract

No new glossary term. The Spec uses **Delivery Retry**, **QA Archive
Override**, **Park Class** and the Implement start head as the `deliver`
guide already does. task_03 documents both rules in the `deliver` guide and
the Roundfix Skill's `deliver` reference.

## Coverage Map

- Goal 1 → The archived retry without a candidate; API Contract 1.
- Goal 2 → The environment-only partial; API Contract 2.
- Goal 3 → The archived retry without a candidate, steps 1, 4 and 5; The
  environment-only partial, step 3.
- Core Feature 1 → The archived retry without a candidate; API Contract 1.
- Core Feature 2 → The environment-only partial; API Contract 2.
- Core Feature 3 → Build Order 3.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 1, 2 and 3.

## Integration Points

- Local Git through `History.Descends` (`git merge-base --is-ancestor`) and
  `git show` of the Run Branch's QA Report; both exist today.
- The Run Database through `History.RunStart`; unchanged.
- No forge, provider or network boundary.

## Testing Approach

1. **Engine retry**, in the new file
   `internal/delivery/operator_archive_any_park_test.go`, with the fakes of
   `operator_archive_retry_test.go` (`operatorArchiveHistory`,
   `fakeItemRecovery`) and a real queue store: an item parked
   `run-unresolved`, `qa-environment-partial` or `run-budget-exceeded`, with
   no candidate and an archive with the override, resumes at `reviewing`
   with candidate commits holding exactly the item head, after `RunStart` of
   its Run ID; the same `run-unresolved` item without the override refuses
   with `candidate head is missing`, and with a head that does not descend
   refuses with API Contract 1's text, both leaving the item unchanged. The
   tests of `operator_archive_retry_test.go` pass unedited.
2. **Real archive**, in `internal/cli/deliver_archived_retry_test.go`: the
   body of `TestArchivedRetryOfAQueueStartedRunReturnsToReview` becomes a
   helper taking the blocker, and a new test runs it for `run-unresolved`
   with no candidate, through `roundfix archive --qa-override` in a linked
   item worktree and a Run recorded in a disposable Run Database. Before this
   Spec the new test fails with `candidate head is missing`, the 0220
   message.
3. **Classification**, in `internal/cli/deliver_operator_archive_test.go`
   through `commandDeliveryWorkflow.runResult` on a disposable Run Branch: a
   new test commits a report shaped like 0220's (two rows with status
   `blocked (environment: no open Pull Request)` and provenance
   `Pull Request row`, `rows_blocked_environment: 2`) and expects the
   environment partial, a declared-only partial and a finding beside those
   rows expect none. The fixture is written inline; no test reads an archived
   Spec. The `pre-PR only` case of
   `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch` changes its
   expectation to `true`, the one declared break.

## Build Order

1. The archived retry without a candidate in `Engine.Retry`, with Testing
   Approach 1 and 2, task_01 (depends on: none).
2. The environment-only partial in `qaEnvironmentPartial`, with Testing
   Approach 3, task_02 (depends on: 3, which writes the guide of the CLI
   surface it changes).
3. The `deliver` guide, the Roundfix Skill's `deliver` reference, the skill
   version and its mirrors, written from this TechSpec, task_03 (depends on:
   none).
4. Terminal QA, task_04 (depends on: 1, 2, 3, 5).
5. Corrective: the override alone admits the archived branch, so an archive
   without a candidate and without History refuses with API Contract 1's
   archived-head text, task_05 (depends on: 1).

## Risks & Considerations

- **A weak anchor.** Any item-branch commit descends from the Run start, so a
  retry after an override could carry operator commits the gate never saw.
  The retry records that head as the candidate and returns to `reviewing`, so
  the pre-PR review, the repository gate and the required checks judge it.
- **A partial the gate should have passed.** A report blocked only by Pull
  Request rows that the gate could have covered with equivalents now parks
  for the operator instead of re-running. The operator still decides whether
  to override, and the status line names the recovery.
- **Another park without a Run.** An item parked before any Run has no Run ID;
  `RunStart` fails and the retry refuses as today.
- **Concurrent skill edits.** Specs 0222 and 0223 may raise the same skill
  version; task_03 raises it by one patch level from the tree it starts on.

## Decisions

- The override, not the park, authorizes the Run start anchor. See ADR-0229.
- A zero-finding partial with any environment-blocked row is an environment
  park, and the Archive Command's eligibility is unchanged. See ADR-0229.
- The refusal texts stay as they are, so status, logs and existing tests keep
  their bytes.
- The override alone admits an archive without a candidate to the
  archived-head comparison; a missing History is a failed proof that refuses
  with API Contract 1's archived-head text, not `candidate head is missing`
  (finding F1 of the 2026-10-04 QA Report, task_05).
