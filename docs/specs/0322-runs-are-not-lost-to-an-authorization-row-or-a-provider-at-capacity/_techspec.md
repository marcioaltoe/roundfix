---
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
prd: _prd.md
created: 2026-10-26
---

# Runs are not lost to an authorization row or a provider at capacity — Technical Spec

## Executive Summary

Two shared policies stop sound work. The QA partial policy in `internal/spec`
exempts an outside-evidence row the Spec's authorization keeps offline only
when its status equals one network sentence; this Spec reads the status as a
form, `blocked (environment: <subject> prohibited by Spec authorization)`,
with a bounded subject, and counts each Results row toward one exempt kind at
most. The Agent Session owner in `internal/daemon` lets the next Fallback
Selection take a Task after a Provider System Error only before the First
Handoff, and a QA turn only while its report says `verdict: pending`; this
Spec removes that window for Task and QA scopes, resends a failed repair
turn's prompt, and reopens a QA verdict the failed turn wrote. The two
Baseline clauses that describe the exemption are reworded and regenerated,
and the skills, guides, ADR-0335 and the glossary follow.

The trade-off accepted: an agent can name a subject for a lookup it was
allowed to make, and a fallback can repeat work the failed turn had done.
The first is the trust ADR-0276 already gives the status, checkable against
the authorization; the second costs turns, never the Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — fake runtimes and temporary
  repositories only. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0335, ADR-0274, ADR-0276 and
  ADR-0295, as `_prd.md` → Project Constraints states. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization recorded
  in `docs/specs/0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity/_authorization.md`
  (Spec-contained authorization record). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/references/settle.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-10-26 at `d08f8557`.

1. `internal/spec/qa.go` `qaReportResultsCounts` walks the `## Results`
   tables outside fences. For each row with a `Status` and a `Provenance`
   column it counts a Pull Request row when the status equals
   `QANoOpenPullRequestStatus`, a Network-Prohibited Row when the status
   equals `QANetworkProhibitedStatus`, and, in a separate `if`, a
   Network-Denied Row when the status starts with
   `QANetworkDeniedStatusPrefix`, ends with `)` and names a host; the last two
   also need `qaReportProvenanceNames(provenance, QAOutsideEvidenceRowSource)`.
2. `QAReport.EnvironmentRowsNeedingOverride` subtracts the three exempt counts
   from `rows_blocked_environment`, and `QAReportEligibility` refuses a
   partial with any row left, naming the kinds it found. `roundfix settle`,
   Daemon settlement (`settleQAVerdict`), `roundfix qa-report accept`,
   `roundfix archive` and the Delivery Queue's `qa-environment-partial`
   decision all call it.
3. Spec 0258's Results table (four exact network rows, a Pull Request row, a
   covered declaration, `rows_blocked_environment: 5`) qualifies. Spec 0309's
   (R06 and R08 exact, R07 `another repository`, R12 Pull Request,
   `rows_blocked_environment: 4`) and Spec 0317's (QA-08 `upstream checkout
   and network`, QA-09 exact, QA-13 Pull Request, `rows_blocked_environment:
   3`) leave one row each and are refused with
   `rows_blocked_environment is <n>, 1 outside the pre-PR Pull Request row and
   the network-denied and network-prohibited outside-evidence rows; expected 0
   outside them`.
4. `internal/daemon/agent_session_owner.go` `Run` sends a turn; on a
   Provider System Error (`agent.DescribeProviderSystemError`) outside review
   it calls `recoverProviderSystemError`, then resends
   `owner.taskPrompt + "\n\n" + providerSystemErrorNotice`.
   `recoverProviderSystemError` closes the session, records the failed
   attempt, and chooses `fallback` only when a next candidate exists and
   either the scope is `task` and `owner.handedOff` is false, or the scope is
   `qa` and `owner.reportPending()` reads `verdict: pending`; otherwise
   `exhausted`, which returns `providerSystemErrorExhaustedError`. A QA
   fallback appends `## Agent runtime fallback` with
   `- provider_system_error:`, `- fallback:` and `- retry_spent: false`.
5. `owner.handedOff` turns true after the first turn returns without error;
   a Verification failure sends Verification Feedback through the same owner,
   so a Provider System Error in that repair turn is `exhausted`.
   `task_engine.go` sets `reportPending` for the QA owner to
   `spec.ReadQAReportFile(...).Verdict == spec.VerdictPending`.
6. On `exhausted`, the Task or QA Task settles `failed` with
   `runtime infrastructure: <scope> <id> hit a provider system error:
   <detail>`; `internal/cli/deliver_workflow.go` `runResult` reads the
   `daemon.task` event (phase `provider_system_error`, recovery `exhausted`)
   and the queue parks `runtime-infrastructure` with an uncounted retry.
7. The Lost Rollout recovery (`recoverLostRollout`) uses the same First
   Handoff and pending-report window to choose a fallback, and otherwise
   continues the same selection in a new session; it resends the lost turn's
   prompt after its notice.
8. The Baseline clauses `clause.spec.project-constraints-06-outside-evidence`
   (rule `rule.spec.project-constraints`, version 6, module `spec-workflow`
   version 21) and `clause.autonomous.loop-01-qa-once` (rule
   `rule.autonomous.loop`, version 7, module `autonomous-work` version 16)
   each carry one sentence naming the exact network status;
   `internal/baseline/network_prohibited_row_clause_test.go` pins both in the
   embedded catalog, the Standard TypeScript formatter goldens and this
   repository's `docs/agents/spec-routing.md` and
   `docs/agents/autonomous-work.md`.

## System Architecture

- `internal/spec/qa.go`: the subject-form constants and reader, and the
  exclusive classification in `qaReportResultsCounts` (task_01).
- `internal/daemon/agent_session_owner.go`: the fall-through, the resent
  prompt and the QA verdict reset (task_02).
- `internal/baseline/assets/modules/spec-workflow.json` and
  `autonomous-work.json`, their derived files, this repository's managed
  guides and Setup Manifest, and the clause test (task_03).
- Skills, guides, glossary, the skill contract test and the Behavior Surface
  record (task_04).
- Tests: `internal/spec/qa_prohibited_subject_test.go` (new),
  `internal/cli/qa_partial_policy_test.go`,
  `internal/daemon/qa_partial_policy_test.go`,
  `internal/daemon/provider_system_error_recovery_test.go`,
  `internal/daemon/provider_system_error_unfinished_qa_test.go` (new),
  `internal/baseline/network_prohibited_row_clause_test.go`.

## Implementation Design

### Interfaces

In `internal/spec/qa.go`, beside the existing outside-evidence constants:

```go
// The Network-Prohibited Row's status form: the prefix, a subject that names
// what the Spec's own authorization forbids, and the fixed suffix. The
// canonical status QANetworkProhibitedStatus is this form with the subject
// "network".
const (
	QAProhibitedStatusPrefix    = "blocked (environment: "
	QAProhibitedStatusSuffix    = " prohibited by Spec authorization)"
	QAProhibitedSubjectMaxRunes = 80
)

// qaProhibitedRowSubject returns the subject of a status in the
// Network-Prohibited Row's form: non-empty after trimming, at most
// QAProhibitedSubjectMaxRunes Unicode code points, and free of parentheses.
func qaProhibitedRowSubject(status string) (string, bool)
```

`QANetworkProhibitedStatus`, `QANetworkDeniedStatusPrefix`,
`QAOutsideEvidenceRowSource`, the `QAReport` fields and every refusal message
keep their names and bytes.

In `internal/daemon/agent_session_owner.go`:

```go
// resetUnfinishedQAVerdict returns a QA report's frontmatter verdict to
// pending when the turn that wrote it did not finish, and names the verdict
// it discarded. A report without a readable frontmatter verdict keeps its
// bytes and yields "".
func resetUnfinishedQAVerdict(path string) (string, error)
```

### Data Models

The QA report's `## Agent runtime fallback` section after a Provider System
Error in a QA turn, appended by the Daemon:

```text
## Agent runtime fallback

- provider_system_error: <selection>: <detail>
- discarded_verdict: <verdict>        (only when the failed turn wrote a verdict other than pending)
- fallback: <next selection>
- retry_spent: false
```

The `daemon.task` event of phase `provider_system_error` keeps its payload
keys (`scope_kind`, `scope_id`, `selection`, `detail`, `recovery`,
`next_selection` for a fallback, `retry_spent: false`).

### API Contracts

1. API Contract: a Results row counts as a Network-Prohibited Row when
   `qaProhibitedRowSubject(status)` succeeds and the provenance names the
   outside-evidence row. The subject is trimmed, 1 to 80 code points
   (`utf8.RuneCountInString`), and holds neither `(` nor `)`; prefix and
   suffix are case-sensitive.
2. API Contract: a row counted as a Network-Prohibited Row is not also
   counted as a Network-Denied Row (`else if`); the Pull Request row, the
   Network-Denied Row rule, skipped counting, frontmatter validation,
   `EnvironmentRowsNeedingOverride` and `QAReportEligibility` keep their code
   and messages.
3. API Contract: `recoverProviderSystemError` chooses `fallback` whenever the
   scope is `task` or `qa` and a next candidate exists, and `exhausted`
   otherwise; the review scope keeps returning the error before recovery.
4. API Contract: after a fallback, `Run` resends
   `owner.taskPrompt + "\n\n" + providerSystemErrorNotice`, followed by
   `"\n\n" + <failed turn's prompt>` when that prompt differs from
   `owner.taskPrompt`.
5. API Contract: for a QA fallback with a report path, the Daemon calls
   `resetUnfinishedQAVerdict` before appending the section, rewrites only the
   frontmatter line `verdict: <v>` to `verdict: pending` when `<v>` is not
   `pending`, and adds `- discarded_verdict: <v>` between the
   `provider_system_error` and `fallback` lines; a pending report gets the
   section without that line, as today.
6. API Contract: the Baseline sentences become, in
   `clause.spec.project-constraints-06-outside-evidence`, "So is an
   outside-evidence row the QA gate did not look up because the Spec's own
   authorization forbids that lookup, recorded as `blocked (environment:
   <subject> prohibited by Spec authorization)`, where `<subject>` names what
   the authorization forbids, such as `network` or `another repository`." and,
   in `clause.autonomous.loop-01-qa-once`, "Nor does an outside-evidence row
   the QA gate did not look up because the Spec's own authorization forbids
   that lookup, recorded as `blocked (environment: <subject> prohibited by
   Spec authorization)`, where `<subject>` names what the authorization
   forbids, such as `network` or `another repository`." The clauses keep
   their ids, enforcement and `replaces`; each rule's `version` rises by one.
7. API Contract: the `### Rows the Spec's authorization keeps offline`
   section of the `qa-gate`, `archive-spec` and `roundfix` skills becomes,
   in each, the text below, and `skills/acceptance_rules_test.go` pins it as
   `networkProhibitedSettlementSection`:

   ```markdown
   ### Rows the Spec's authorization keeps offline

   The qualifying declared `partial` above also accepts an outside-evidence row
   the gate did not look up because the Spec's own authorization forbids that
   lookup, recorded as `blocked (environment: <subject> prohibited by Spec authorization)`
   with the outside-evidence row named in its provenance. The subject names
   what the authorization forbids, in its words: `network` for a network
   lookup, or `another repository`; it is 1 to 80 characters and holds no
   parenthesis. Like the network-denied row, it never decides a qualifying
   partial and needs no Unreachable Acceptance declaration, a row counts
   toward one exempt kind at most, and settlement, `roundfix qa-report accept`
   and `roundfix archive` read it through the same policy (ADR-0276,
   ADR-0335).
   ```

   The Roundfix `archive` and `settle` references carry, in place of today's
   sentence and pinned as `networkProhibitedReferenceSentence`: "So does an
   outside-evidence row the gate did not look up because the Spec's own
   authorization forbids that lookup, recorded as `blocked (environment:
   <subject> prohibited by Spec authorization)` with the outside-evidence row
   named in its provenance, where the subject names what the authorization
   forbids in 1 to 80 characters without parentheses (ADR-0276, ADR-0335)."
   The settle, qa-report and context-driven guides carry the same sentence
   without the parenthesized ADRs.
8. API Contract: the usage and configuration guides and the Roundfix
   `profiles` reference replace "A Provider System Error lets the next
   Fallback Selection take the Task under the same condition, gets no retry
   on the same selection, does not count toward the Lost Rollout budget, and
   otherwise settles the Task as runtime infrastructure." with "A Provider
   System Error lets the next Fallback Selection take the Task in any Task or
   QA turn, a repair turn after the First Handoff and a QA turn that already
   wrote a verdict included, gets no retry on the same selection, does not
   count toward the Lost Rollout budget, and settles the Task as runtime
   infrastructure only when the chain has no next selection." The usage
   guide's QA fallback paragraph adds that a fallback after a written verdict
   first rewrites the report's `verdict:` line to `pending` and records
   `- discarded_verdict: <verdict>` before the `fallback` line.

### Surface Transcripts

None. No command output changes; the report section line and the clause
sentences are fixed above and asserted by named tests.

### Invariants

```text
1. Every QA Report that qualifies at d08f8557 still qualifies, and every one that is refused for a reason other than an outside-evidence row in the subject form is still refused with the same message.
2. A Results row counts toward one exempt kind at most.
3. A Provider System Error in a Task or QA scope never ends the work while the Fallback Chain has a next selection.
4. A QA fallback session never starts from a verdict the failed turn wrote.
5. An exhausted chain still settles the Task failed as runtime infrastructure and parks the item runtime-infrastructure with an uncounted retry.
6. Review Batches, Lost Rollout recovery, the Light Tier's skip, yield and escalation, and the selection-start fallback behave as before.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1 and 3 → API Contract 1, Invariant 1, task_01
- Goal 2, Core Feature 1, Success Metric 2 → API Contracts 1 and 2, Invariants 1 and 2, task_01
- Goal 3, Core Feature 2, Success Metric 4 → API Contracts 3 and 4, Invariants 3 and 5, task_02
- Goal 4, Core Feature 2, Success Metric 5 → API Contract 5, Invariant 4, task_02
- Success Metric 6 → Invariant 6, task_02
- Goal 5, Core Feature 3, Success Metric 7 → API Contract 6, task_03
- Goal 5, Core Feature 4, Success Metric 7 → API Contracts 7 and 8, Interfaces (skills, guides, glossary), task_04

## Absence Evidence

1. A subject-form reader for the Network-Prohibited Row — command: `grep -rn "QAProhibitedStatusPrefix\|qaProhibitedRowSubject\|prohibited by Spec authorization)\"" internal`; result: only `QANetworkProhibitedStatus` in `internal/spec/qa.go` at `d08f8557`.
2. A QA verdict reset on fallback — command: `grep -rn "discarded_verdict\|resetUnfinishedQAVerdict" internal docs .agents`; result: no match at `d08f8557`.

## ADR Alternatives

- ADR-0335: introduces both rules; its rejected alternatives (skill wording
  only, reading `_authorization.md`, any status naming the authorization,
  widening only QA, telling the fallback without reopening the verdict) are
  the ones this Spec does not build.
- ADR-0276: amended; its rejection of "any status that mentions the
  authorization" holds outside the subject form, its rejection of reading
  `_authorization.md` stands, and its Delivery Base half is untouched.
- ADR-0274: amended; its rejection of retrying the same selection stands, and
  its fallback window widens to every Task and QA turn. The Lost Rollout
  recovery it was modeled on keeps its own window, because the same selection
  can continue a lost rollout but not a provider failure; the Light Tier
  escalation after a failed Verification is the precedent for another
  selection inheriting a handed-off attempt.
- ADR-0295: unchanged; it cites ADR-0276 only for the Delivery Base, which
  this Spec does not touch.

## Integration Points

None. No network, provider, forge or other repository is reached; every test
uses temporary repositories, fake runtimes, fake runners and fake
verifiers.

## Testing Approach

Every test builds temporary repositories or reports; none reads this
repository's Specs, and none reaches the network.

- `internal/spec/qa_prohibited_subject_test.go` (new), with a
  `Suite:`/`Invariant:`/`Boundary:` header, reusing `readPartialPolicyReport`,
  `partialPolicyTable`, `assertPartialPolicyEligibility` and
  `writePrePullRequestDeclaration`:
  - `TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids`, subtests
    `accepts Spec 0309 R07 another repository`, `accepts Spec 0317 QA-08
    upstream checkout and network`, `accepts the canonical network form`,
    `accepts a subject of 80 multibyte characters` (80 × `é`), `refuses an
    empty subject`, `refuses a blank subject`, `refuses a subject of 81
    characters`, `refuses a subject with a parenthesis`, `refuses a changed
    suffix case`, `refuses a row without outside evidence`, `a row in both
    forms counts once` (subject `network denied: docs.example.com` beside a
    `credentials unavailable` row, `rows_blocked_environment: 2`, refused with
    today's two-kind message), `a skipped row still refuses`.
  - `TestTheArchivedAuthorizationRowsQualify`, subtests `Spec 0258`,
    `Spec 0309`, `Spec 0317`: the archived Results rows of
    `_prd.md` → Validation with their environment counts and one covered
    declaration each.
- `internal/cli/qa_partial_policy_test.go` and
  `internal/daemon/qa_partial_policy_test.go`: `qaPartialShapes` gains, in
  both copies, `Spec 0309 another repository prohibited with PR and
  declaration` (accepted), `Spec 0317 upstream checkout and network
  prohibited with PR and declaration` (accepted) and `authorization subject
  with a parenthesis` (refused, `rows_blocked_environment is 1; expected 0`);
  every existing shape keeps its expectation (declared change: three added
  shapes per table).
- `internal/daemon/provider_system_error_recovery_test.go` (declared
  changes):
  - new `TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback`: the
    repair turn fails, three turns run, the third on `fallback-model` in a
    `-fallback-01` session whose prompt starts with the Task prompt, holds
    the notice and ends with the failed repair prompt; recovery `fallback`;
    one `agent_selection_fallback` event with `reason_code`
    `provider_system_error`.
  - `TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure`
    keeps its name and now runs without a fallback selection: two turns, one
    failed Task with today's reason, recovery `exhausted`.
  - `TestProviderSystemErrorRecoveryEligibility`: its case `QA with a
    finished report exhausts` becomes `QA with a finished report falls back`,
    asserting two turns, the second on `fallback-model`, recovery `fallback`.
- `internal/daemon/provider_system_error_unfinished_qa_test.go` (new): a
  `verdictBeforeFailureRunner` that rewrites the seeded report's verdict
  during the failing turn, and
  - `TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt`, subtests
    `partial`, `pass`, `fail` (the fallback receives `verdict: pending`, the
    section has `- discarded_verdict: <verdict>` once with the other three
    lines, the QA is accepted, one failed attempt, recovery `fallback`) and
    `a pending report records no discarded verdict`;
  - `TestQAProviderSystemErrorAfterAVerdictWithoutFallbackSettlesRuntimeInfrastructure`:
    one selection, one turn, recovery `exhausted`, QA not accepted.
- `internal/baseline/network_prohibited_row_clause_test.go`: both sentence
  constants become API Contract 6's sentences (declared change).
- `skills/acceptance_rules_test.go`: `networkProhibitedSettlementSection`
  and `networkProhibitedReferenceSentence` become API Contract 7's texts
  (declared change); `TestTheSettlementSkillsNameTheNetworkProhibitedRow`
  keeps its name and checks.
- Unchanged and run: the existing partial-policy tests of `internal/spec`,
  `TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback`,
  `TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure`,
  `TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt`,
  `TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting`,
  `TestProviderSystemErrorDoesNotUseLostRolloutBudget`, the Light Tier tests
  of `light_tier_test.go` and `light_tier_access_test.go`, the Lost Rollout
  tests, `TestCoverageEquivalence` (no top-level test is removed) and the
  Delivery Queue's provider park test in
  `internal/cli/deliver_provider_system_error_test.go`.
- A prototype in scratch clones passed all of them, and each mutation failed
  the named case, whose input makes it observable:
  - exact `==` comparison kept: `accepts Spec 0309 R07 another repository`,
    `TestTheArchivedAuthorizationRowsQualify/Spec 0309` and the two CLI and
    Daemon shapes;
  - subject length counted in bytes: `accepts a subject of 80 multibyte
    characters` (160 bytes);
  - parenthesis rule dropped: `refuses a subject with a parenthesis` and the
    CLI and Daemon parenthesis shape;
  - `else if` replaced by a second `if`: `a row in both forms counts once`;
  - `!owner.handedOff` kept for Tasks:
    `TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback`;
  - `reportPending()` kept for QA: `QA with a finished report falls back` and
    every verdict subtest of
    `TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt`;
  - the verdict reset skipped: the `partial`, `pass` and `fail` subtests;
  - the failed repair prompt not resent:
    `TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback`.

## Build Order

1. The subject form and its tests (task_01; no dependency).
2. The fall-through, the resent prompt, the verdict reset and their tests
   (task_02; no dependency; shares no file with task_01).
3. The two Baseline clauses, their regeneration and the clause test
   (task_03; no dependency: skill edits do not move the Baseline's derived
   files, measured by running `make baseline-digests` after a skill edit on a
   prototype, which reported `"changed":false`).
4. Skills, guides, glossary and the Behavior Surface record re-recorded last
   (task_04; no dependency: its Verification reads only documents, and it
   shares no file with tasks 01 to 03).
5. QA gate (task_05) (depends on: 1, 2, 3, 4).

No Spec is listed in `requires`. Specs 0319 and 0320 deliver first; this Spec
reads neither's code.

## Risks & Considerations

- A subject is free text: an agent could name a lookup its authorization
  allowed. The report must name the authorization text (the `qa-gate` rule
  stands) and the archive keeps both for a reader to check.
- A fallback after the First Handoff inherits a tree and feedback another
  model produced. The repair prompt is resent whole, as Lost Rollout
  recovery does, and the Task still settles only on its declared
  Verification.
- Rewriting the report's verdict line edits a file the agent wrote. Only the
  frontmatter `verdict:` line changes, the discarded value is recorded, and a
  report whose frontmatter cannot be read keeps its bytes.
- The Baseline rewording reaches every adopter on their next update as a
  changed clause in the same rules; the ids, enforcement and `replaces` stay,
  so no Source Baseline retention changes.

## Glossary

- changes: **Network-Prohibited Row**
- changes: **Provider System Error**
- changes: **First Handoff**
- changes: **Fallback Chain**

## Decisions

- Subject form in the parser, bounded at 80 code points without
  parentheses, one exempt kind per row; skills teach it.
- Baseline clauses reworded in place, rule and module versions raised, no
  `replaces`.
- Provider System Error falls through for Task and QA turns at any point;
  review and Lost Rollout unchanged.
- QA verdict reopened with a recorded `discarded_verdict`.
- ADR-0335, with dated amendments to ADR-0274 and ADR-0276.
