---
spec: 0235-one-qa-partial-policy
prd: _prd.md
created: 2026-10-06
---

# One QA partial policy — Technical Spec

## Executive Summary

The four QA settlement callers already call `spec.QAReportEligibility`, so the
fix lives in that one function and in what `ReadQAReport` derives from the
Results rows. The reader gains two structural counts beside the existing
pre-PR count: network-denied outside-evidence rows and skipped rows. It also
recognizes a Pull Request row whose provenance carries a note. Eligibility
then accepts a `partial` whose unmet rows are only those exempt environment
rows and covered declared rows. A new `QAReport` method names the environment
rows that still need the override. The Delivery Queue reads that method instead of
its own predicate, and `roundfix settle` names the report it judged. The skills,
the two Baseline clauses and the user guides then state the one policy. The
primary trade-off is trusting the gate's written marker for a network denial
that the Daemon cannot observe. The named host makes the claim checkable
afterwards, and the alternative, Daemon-side network observation, would
change the sandbox contract (ADR-0240).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  new names are the status marker `blocked (environment: network denied:
  <host>)`, two `QAReport` fields and one method. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — eligibility reads QA Reports
  from disk and from Git; no credential or network call is added, and every
  test uses temporary repositories, temporary homes and fake runners. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0240 (this Spec): "The policy has
  one rule, in `internal/spec`, which every caller reads", and "Any other
  environment row, including a network-denied row with no outside-evidence
  provenance, still needs declaring or the override". It extends ADR-0167,
  whose pre-PR row "no longer counts against a qualifying `partial`";
  supersedes ADR-0104 in part, whose gate "holds pull request preparation
  until that row is satisfied"; and supersedes ADR-0229 in part, whose queue
  parks `qa-environment-partial` "including when every environment-blocked
  row is one that waits for an open Pull Request". ADR-0080 still decides
  `pass`: an environment-blocked row "does not cap the verdict when the
  report records equivalent observed or supervised evidence". ADR-0154 stands for what stays ineligible: a QA archive override records user authority, not a pass, ADR-0154: "The exception can waive the terminal QA Task's archive". ADR-0088, ADR-0091, ADR-0096,
  ADR-0097, ADR-0155, ADR-0182, ADR-0194, ADR-0195 and ADR-0210 bind the
  gate, its machine stage, its row carry, its matrix and its evidence.
  ADR-0179 bounds the Governed Paths, ADR-0187, ADR-0189 and ADR-0233 bind
  the owned skills' versions, and ADR-0184 binds the Surface Transcripts.
  ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check
  consistency by citation and receipt. ADR-0237 cites ADR-0229 but decides how a Delivery Retry records a merge made outside the queue, which this Spec leaves unchanged, so it does not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_03 changes the skill Governed Paths
  and task_04 the Baseline ones; task_01 and task_02 change none. Express
  maintainer authorization: "Concedo" of 2026-10-06, beside the standing
  grants of 2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário"; "Autorizar os dois"); bounded files:
  `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`,
  `docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/archive-spec/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0235-one-qa-partial-policy/_authorization.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| QA Report reader | `readQAReport`, `qaReportPrePullRequestRows` in `internal/spec/qa.go` | Derives network-denied outside-evidence rows and skipped rows; the Pull Request provenance item may carry a note |
| Eligibility | `QAReportEligibility`, new `QAReport.EnvironmentRowsNeedingOverride` in `internal/spec/qa.go` | One rule for a qualifying partial |
| Settle | `requireSettleQAReportEligibility` and its caller in `internal/cli/settle.go` | Refusal names the report path |
| Queue classification | `qaEnvironmentPartial` in `internal/cli/deliver_workflow.go` | Reads `EnvironmentRowsNeedingOverride` |
| Daemon settlement, archive, `qa-report accept` | `internal/daemon/task_engine.go`, `internal/spec/archive.go`, `internal/cli/qa_report.go` | Unchanged code; they read the same function |
| Skills | qa-gate, archive-spec, Roundfix (`SKILL.md`, `references/archive.md`, `references/settle.md`) | The marker, the table row, the references |
| Baseline guides | `clause.spec.project-constraints-06-outside-evidence` in `spec-workflow.json`; the delivery-order clause of `autonomous-work.json` | One added sentence each, then derived files |
| User guides | `docs/user-guide/commands/qa-report.md`, `archive.md`, `settle.md`, `docs/user-guide/context-driven-development.md` | State the policy and the settle refusal |

## Implementation Design

### Interfaces

```go
// internal/spec/qa.go
const (
	QANetworkDeniedStatusPrefix = "blocked (environment: network denied: "
	QAOutsideEvidenceRowSource  = "outside-evidence row"
)
type QAReport struct {
	// existing fields unchanged
	RowsBlockedNetworkDenied int // derived from Results, like RowsBlockedPrePullRequest
	RowsSkipped              int // derived from Results
}
func (report QAReport) EnvironmentRowsNeedingOverride() int
func QAReportEligibility(specDir string, report QAReport) error // signature unchanged
```

```text
1. A pre-PR Pull Request row is a Results row whose Status cell, trimmed, equals QANoOpenPullRequestStatus exactly and whose Provenance cell has an item (split on ";" and ",", trimmed) that equals QAPullRequestRowSource or starts with it followed by " ", ":" or "(".
2. A network-denied outside-evidence row is a Results row whose Status cell, trimmed, starts with QANetworkDeniedStatusPrefix, ends with ")", and has non-empty text between them after trimming, and whose Provenance cell has an item that equals QAOutsideEvidenceRowSource or starts with it followed by " ", ":" or "(". Matching is case-sensitive, as for the Pull Request row.
3. A skipped row is a Results row whose Status cell, trimmed, equals "skipped" ignoring case. Only Results tables with a Status column count; fenced blocks and tables outside "## Results" never count, as today.
4. EnvironmentRowsNeedingOverride is RowsBlockedEnvironment minus min(RowsBlockedEnvironment, RowsBlockedPrePullRequest + RowsBlockedNetworkDenied), never negative.
5. A pass is judged exactly as today.
6. A partial is refused, in this order: rows_blocked_finding > 0 (message unchanged); EnvironmentRowsNeedingOverride > 0 (message unchanged when RowsBlockedNetworkDenied is 0, else API Contract 2); RowsSkipped > 0 (API Contract 3); rows_blocked_declared greater than the Spec's declarations (message unchanged; declarations are read only when rows_blocked_declared > 0); no exempt environment row and no declared row (message `newest QA Report verdict is "partial"; expected "pass"`, unchanged); a Hollow report (unchanged). Otherwise it qualifies.
7. Every other verdict is refused exactly as today.
8. The Delivery Queue classifies an unresolved Run qa-environment-partial only when the newest report on the Run Branch is partial, has rows_blocked_finding 0 and EnvironmentRowsNeedingOverride() > 0.
9. roundfix settle's eligibility refusal keeps its prefix and appends " (report <path>)", the path relative to the settle surface's working tree.
10. Daemon settlement, roundfix archive and roundfix qa-report accept keep calling QAReportEligibility; no caller adds a condition of its own.
```

### Data Models

No schema, frontmatter or Run Database change. The new `QAReport` fields are
derived from the Results rows, like `RowsBlockedPrePullRequest`. The
frontmatter count `rows_blocked_environment` keeps counting every
environment-blocked row, exempt or not, so the mechanical stage's count
checks and every archived report read as before.

### API Contracts

1. API Contract: `roundfix qa-report accept <path>`, `roundfix archive <slug>`,
   `roundfix settle --spec <slug> --task <qa task>` and Daemon QA settlement
   accept a `partial` that meets Invariant 6 and refuse every other
   `partial` with the same reason text.
2. API Contract: a partial with a non-exempt environment row beside a
   network-denied outside-evidence row is refused with
   `rows_blocked_environment is <n>, <m> outside the pre-PR Pull Request row and network-denied outside-evidence rows; expected 0 outside them`.
3. API Contract: a partial with a skipped row is refused with
   `partial records <n> skipped row(s); a qualifying partial records none`.
4. API Contract: a refused `roundfix settle` of a QA Task prints
   `roundfix: settle QA Report is ineligible: <reason> (report <path>)` on
   stderr and exits 1, leaving the Task file unchanged.
5. API Contract: the Delivery Queue parks an unresolved Run
   `qa-environment-partial` only under Invariant 8, and otherwise
   `run-unresolved` as before.

### Surface Transcripts

1. Surface Transcript: a partial whose only unmet row is the pre-PR Pull
   Request row is accepted.

   ```transcript
   $ roundfix qa-report accept docs/specs/<slug>/qa/qa-report-<date>.md
   stdout:
   stderr:
   exit: 0
   ```

2. Surface Transcript: a network-denied row without an outside-evidence
   provenance still refuses.

   ```transcript
   $ roundfix qa-report accept docs/specs/<slug>/qa/qa-report-<date>.md
   stdout:
   stderr:
   roundfix: QA Report eligibility refused: rows_blocked_environment is 1; expected 0
   exit: 1
   ```

3. Surface Transcript: settle names the report it refused.

   ```transcript
   $ roundfix settle --spec <slug> --task <task>
   stdout:
   ...
   stderr:
   Settle surface: <path>
   roundfix: settle QA Report is ineligible: <reason> (report <path>)
   exit: 1
   ```

## Coverage Map

- Goal 1 → `QAReportEligibility` as the only decision (Invariant 10; API Contract 1)
- Goal 2 → reader counts and eligibility (Invariants 1, 2, 4, 6); queue classification (Invariant 8)
- Goal 3 → Invariant 6's remaining refusals; `EnvironmentRowsNeedingOverride`
- Goal 4 → settle refusal (Invariant 9; API Contract 4; Surface Transcript 3)
- Goal 5 → derived counts only; frontmatter unchanged (Data Models)
- Core Feature 1 → `QAReportEligibility` (Invariants 5 to 7; API Contracts 2, 3)
- Core Feature 2 → Invariant 1
- Core Feature 3 → Invariant 2; the qa-gate skill
- Core Feature 4 → Invariant 8; API Contract 5
- Core Feature 5 → Invariant 9
- Core Feature 6 → skills, Baseline clauses, user guides
- Success Metric 1 → task_01 tests on the 0213 and 0220 shapes; QA replay on the archived reports
- Success Metric 2 → task_02 agreement tests over the four callers
- Success Metric 3 → task_02 queue classification tests
- Success Metric 4 → task_02 settle test; Surface Transcript 3
- Success Metric 5 → the repository Verification at each Task's settlement

## Integration Points

None outside the repository. The Daemon, the Delivery Queue, settle, archive
and `qa-report accept` are in-process readers of the same package.

## Testing Approach

- `internal/spec`: a new `qa_partial_policy_test.go` builds reports through
  `ReadQAReportFile` on temporary Spec directories, with the row shapes of the
  measured corpus: Pull Request rows only (0213, 0220, Oraculum), an
  annotated Pull Request row (0227), declared rows with the Pull Request row
  (Fluxus), a network-denied outside-evidence row, a network-denied row
  without that provenance, another environment row, a skipped row and a
  finding row. In `qa_prepr_row_test.go`, the subtest "partial without a
  declared row keeps its refusal" now expects acceptance. That is a declared
  break, renamed to say so, and every other case there keeps its message.
- `internal/cli`: a new `qa_partial_policy_test.go` runs the same shapes
  through `qa-report accept`, `archive` and `settle` with the existing
  archive and settle harnesses, and asserts that the three agree. The queue
  classification cases "pre-PR only" and "pull request rows only" in
  `deliver_operator_archive_test.go` flip to `false`, a declared break, and
  a new case with a network-denied outside-evidence row also expects
  `false`.
- `internal/daemon`: a new `qa_partial_policy_test.go` drives
  `Engine.TaskCycle` with the existing `taskFakeRunner` and the same shapes,
  and asserts that the Daemon settles each one as the CLI callers do.
- `internal/baseline`: a new test asserts the two added sentences in the
  embedded clauses, in the formatter goldens and in `docs/agents/`. The
  existing `TestAuthorizationClausesStateTodaysRefusals` keeps passing,
  because the old sentence stays and the new one follows it.
- `skills`: `TestSettlementGuidanceIsOneTable` proves that the table row is
  identical in all three skills, and `TestEveryOwnedSkillVersionIsRecorded`
  proves that the raised versions are recorded.

## Build Order

1. One policy in `internal/spec`: reader counts, annotated Pull Request row,
   eligibility, `EnvironmentRowsNeedingOverride`, tests.
2. The callers: settle names the report, the queue reads
   `EnvironmentRowsNeedingOverride`, and agreement tests across the CLI
   callers and the Daemon (depends on: 1).
3. The skills and user guides state the policy and the marker (depends on: 1).
4. The two Baseline clauses, `make baseline-digests` and the Managed Refresh
   (depends on: 3).
5. Final QA gate (depends on: 2, 4).

Steps 2 and 3 share no file. Step 4 follows step 3 so that only one Task at a
time edits guidance and regenerates derived files.

## Risks & Considerations

- A gate can label a reachable row network-denied. The marker names the
  host, the row stays in the archived report, and the QA gate of this Spec
  checks the skill text that asks for it.
- Spec 0236 changes Baseline code, and possibly the same derived files. This
  Spec merges first. The derived files are regenerated, never merged by hand.
- The three owned skills' versions are raised by the record command, which
  picks free versions, so a concurrent raise in 0236 or 0237 is resolved at
  merge (ADR-0233).
- `### QA settlement` is the contract table itself. This Spec changes its
  qualifying-partial row identically in all three skills, which the
  contract test requires. Other Specs leave that section alone so the copies
  never drift; this Spec owns the policy the row states, so it changes all
  three copies in one Task.

## Exact texts

The Baseline sentences are appended after the existing sentence of each
clause:

- `clause.spec.project-constraints-06-outside-evidence`: "An outside-evidence
  row blocked only because the Run sandbox denied network access, recorded as
  `blocked (environment: network denied: <host>)`, is the exception: the
  report records that the source was not reached, the row never decides a
  qualifying `partial`, and whoever needs that proof declares it under
  Unreachable Acceptance."
- the `autonomous-work.json` delivery-order clause, after "never decides a
  qualifying `partial`.": "Neither does an outside-evidence row blocked only
  because the Run sandbox denied network access, recorded as
  `blocked (environment: network denied: <host>)`."

The `### QA settlement` row "qualifying declared `partial`" keeps its label
and Archives cell. Its Settles cell becomes: "Settles the QA Task as
`completed` when no row failed, was skipped or is finding-blocked, every
declared-blocked row is covered by a matching `## Unreachable Acceptance`
declaration, and every environment-blocked row is the pre-PR Pull Request
row, recorded as `blocked (environment: no open Pull Request)` with the Pull
Request row named in its provenance, or an outside-evidence row the Run
sandbox could not reach, recorded as `blocked (environment: network denied:
<host>)` with the outside-evidence row named in its provenance. Neither row
needs an Unreachable Acceptance declaration, and a partial whose only unmet
rows are such rows qualifies."

## Vocabulary Contract

- emits: `internal/spec/qa.go`
  pattern: `network denied: `
  documented-in: `docs/user-guide/commands/qa-report.md`

No glossary term is adopted. `CONTEXT.md`'s **QA Report** entry names only the
Pull Request row as exempt; the QA gate records whether the glossary needs
the network-denied row.

## Decisions

- One rule in `spec.QAReportEligibility`; callers add nothing. See ADR-0240.
- The gate keeps one verdict rule: evidence decides `pass`, and the row kinds
  decide whether a `partial` qualifies. See ADR-0240.
- A skipped row keeps a partial from qualifying, so "only unmet rows" holds.
- The exempt counts are derived from the Results rows; the frontmatter is
  unchanged.
- The Baseline gains sentences rather than rewording existing ones, so the
  pinned clause test and adopters' current text stay valid.
