---
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
prd: _prd.md
created: 2026-10-27
---

# Settle follows the branch's plan and the Task commit's rules, and a declared row matches the criterion it names — Technical Spec

## Executive Summary

Four defects around the QA acceptance and the Settle Command (`_prd.md` →
Sources). First, `spec.QAReportEligibility` checks that
`rows_blocked_declared` does not exceed the number of Unreachable Acceptance
declarations, so one declaration covering two criteria cannot cover two
rows, and a second `criterion:` line silently replaces the first. The
declaration now names several criteria, and a report whose declared-blocked
rows each match one is accepted; today's count stays the fallback, so no
report accepted today is refused (new ADR-0343, widening ADR-0240's
"covered"). Second,
settle loads the Spec from the surface it settles, so a plan amended on the
branch is not read; settle now fast-forwards a Run Worktree the branch is
ahead of and refuses any other Spec-folder drift. Third, settle run inside a
Run Worktree takes it for the user checkout; it now refuses there. Fourth,
settle stages every change, so it commits an untracked executable a Task
commit refuses; it now applies the Task commit's rule (ADR-0340).

The trade-off accepted: two acceptance paths instead of one, so a paraphrased
row still passes when the count allows it, and settle may move a kept Run
Worktree before its Verification. Matching alone would refuse reports
accepted today; counting criteria would accept rows the Spec never declared.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local repository commands only.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the new ADR-0343; ADR-0240, ADR-0104,
  ADR-0167, ADR-0229, ADR-0237, ADR-0274, ADR-0276, ADR-0295, ADR-0312,
  ADR-0335, ADR-0340 and ADR-0157, as `_prd.md` → Project Constraints
  states. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AR answers quoted in the authorization record grant the ten Governed
  Paths below. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/settle.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`. No Makefile, workflow, Go module, lint or formatter change.

## Current behavior

Read and measured on 2026-10-27 at `5e9475bf`, in a clone.

1. `internal/spec/spec.go` `parseUnreachableDeclarations` keeps one
   `Criterion` per declaration: `readField` assigns it, so a second
   `criterion:` line replaces the first, and `appendContinuation` appends to
   it.
2. `internal/spec/qa.go` `QAReportEligibility`, for a `partial` with
   `RowsBlockedDeclared > 0`, reads `Unreachable(specDir)` and refuses
   `rows_blocked_declared is %d, but Spec declares %d unreachable
   acceptance%s; shortfall is %d` when the count exceeds the declarations. It
   never reads which criterion a row names. `readQAReport` derives the exempt
   environment rows and the skipped rows from the Results table
   (`qaReportResultsCounts`) but no declared row. Its callers:
   `internal/cli/qa_report.go` (`qa-report accept`), `internal/cli/settle.go`
   `requireSettleQAReportEligibility`, `internal/daemon/task_engine.go`
   `settleQAVerdict`, and `internal/spec/archive.go`
   `archiveUnprovenActions`, which records each declaration's `satisfied-by`.
3. `internal/speccheck/mechanical.go` `detectMechanicalReportShape` counts a
   Results row as declared when `typedBlockedStatus(lower, "declared")` holds
   (`blocked (declared:` … `)` with a non-empty inside) and reports a
   frontmatter count that disagrees; `internal/daemon/task_engine.go` writes
   `rows_blocked_declared` from the same body count. A row
   `blocked (declared unreachable)`, which Specs 0309, 0317 and 0327 authored
   in their QA Tasks, is reported there as a cause outside the three kinds.
4. The QA prompt (`internal/agent/spec_prompt.go` `qaGateContract`) and the
   `qa-gate` skill name "each PRD Unreachable Acceptance declaration" a
   non-waivable source, and the `write-prd` template says "Declare one
   unreachable criterion with these fields". The `### QA settlement` table
   already says "every declared-blocked row is covered by a matching
   `## Unreachable Acceptance` declaration"; it does not change.
5. `internal/cli/settle.go` `preflightSettle` sets `userRoot`, `workDir` and
   the `runworktree.Ref.UserRoot` to the repository root of the working
   directory. Inside a kept Run Worktree, `latestKeptSettleRun` still finds
   the Run (`store.LatestKeptSpecRun` resolves the repository through
   `RepositoryRoots`), the surface is that worktree, and `worktree.Integrate`
   finds the target branch checked out elsewhere: measured, settle commits in
   the Run Worktree and prints `integration pending — git merge --ff-only
   roundfix/run-<id>`, exit `1`.
6. `resolveSettleSurface` loads `spec.Load` from each candidate's own Spec
   Root (`specsRootForWorkDir`), so the Task file, `_tasks.md` and the PRD
   the QA eligibility reads all come from the surface. Nothing compares the
   surface with the Run's target branch (`run.LocalBranch`). Measured: a QA
   Task settled in a kept Run Worktree after the branch amended the PRD
   exits `1` with the report ineligible against the old PRD.
   `readSettleAuthorization` reads the authorization at `run.HeadSHA` in the
   user root, the record the Run resolved at start.
7. `worktree.Integrate` fast-forwards the target branch onto the Run Branch
   only when the target is an ancestor of it; otherwise settle prints
   `integration pending` (`TestSettleNamesNoRemovalWhenIntegrationIsPending`
   pins that for a user branch that gained a commit after the Run started).
8. `settleTaskAndCommit` stages with `git add --all` (`addAllChanges`) and
   filters nothing. Measured: an untracked `newdir/helper.sh` with mode
   `0755` is committed. The Daemon's Task commit drops it through
   `filterStageablePaths` with `authorizedExecutables(plan.Authorization)`
   (ADR-0340) and keeps a tracked path whatever its mode (ADR-0157).
9. Task Carry-Forward (`internal/cli/carryforward.go`) compares each carried
   Task's declared inputs at its settlement commit; a fast-forward of the Run
   Branch rewrites none of those commits.

## System Architecture

### Declared rows matched by criterion

`UnreachableDeclaration` gains `Criteria []string`; `Criterion` stays the
first one for existing readers. `QAReport` gains `DeclaredRows`, each the
row's first cell and the criterion its status names, read in the same pass
as the exempt rows. `QAReportEligibility` accepts when every listed row
matches a criterion of some declaration, and otherwise applies today's count
unchanged; only when both fail does it refuse, naming the unmatched rows
(API Contracts 1-3). The mechanical stage
calls the same classifier (API Contract 4).

### Settle runs from the user checkout

`preflightSettle` refuses, right after it inspects Git and before it opens
the Run Database for writing or selects a surface, when the working
directory's branch is a `roundfix/run-` branch (API Contract 5). The branch
names the Run, so the refusal reads that Run's recorded `GitRoot` read-only
and tells the caller where to run the same command; nothing is committed,
integrated or removed, and the caller's shell never stands in a directory
settle removes.

### Settle follows the branch

After every other preflight check, when the selected surface belongs to a
kept Run, settle compares the surface's HEAD with the Run's target branch
and acts only on files the branch changed under the Spec folder since their
merge base (API Contract 6).

### Settle commits what a Task commit keeps

Preflight lists the surface's untracked files and asks the Daemon's stage
filter which of them a Task commit would refuse as executables, under the
authorization settle already reads (API Contract 7).

### The rules in writing

Each Task carries its rule into the documents that describe the code it
changes, so no document states a rule before its code exists (API Contracts
8-11). task_01 writes the acceptance rule into the `qa-gate` skill under its
own heading, the `write-prd` PRD template, the Roundfix `settle` reference's
QA Report acceptance section, the `qa-report` guide, the QA prompt and the
**Unreachable Acceptance** entry. task_02, task_03 and task_04 each add one
sentence to `settleUsage`, the `settle` guide and the Roundfix `settle`
reference (task_02 also its `spec-delivery` step) and revise one glossary
entry. Every Task syncs the skill mirrors, records the owned-skill versions
and re-records the Behavior Surface fingerprint as its last step.

## Implementation Design

### Interfaces

```go
// internal/spec
type UnreachableDeclaration struct {
	Criterion   string   // the first criterion, unchanged for existing readers
	Criteria    []string // every criterion, in order
	Reason, SatisfiedBy string
	Line        int
}
type QADeclaredRow struct{ ID, Criterion string }
// QAReport gains DeclaredRows []QADeclaredRow
func QADeclaredRowCriterion(status string) (string, bool)
func CriterionNamed(declared, row string) bool
// internal/daemon
func UnauthorizedExecutables(ctx context.Context, workDir string, paths []string, resolution spec.AuthorizationResolution) []DroppedStagePath
```

### Data Models

The PRD section accepts a repeated field:

```markdown
- criterion: Success Metric 10, ...
  criterion: Success Metric 11, ...
  reason: ...
  satisfied-by: ...
```

### API Contracts

1. API Contract: declarations. Each `criterion:` line starts a new criterion
   of the current declaration; an indented continuation line extends the last
   one. A declaration with an empty criterion, or none, keeps today's
   `missing criterion` error. `Criterion` is `Criteria[0]`.
2. API Contract: declared rows. In the Results table, a row whose Status cell
   is `blocked (declared: <criterion>)` (prefix case-insensitive, closing
   `)`, non-empty after trimming) is a declared row with that criterion and
   its first cell as its ID. `CriterionNamed(declared, row)` normalizes both
   (drop backticks, collapse whitespace, trim spaces and trailing periods,
   lower-case) and holds when they are equal, or when the declared text
   starts with the row text followed by `,`, `:`, `;`, `.`, ` (`, ` —` or
   ` - `.
3. API Contract: eligibility. For a `partial` with `rows_blocked_declared`
   above zero, the declared rows are covered when either holds:
   - matched: the report lists declared rows, every one is named by a
     criterion of some declaration, and `rows_blocked_declared` does not
     exceed their number;
   - today's count: `rows_blocked_declared` does not exceed the number of
     declarations, exactly as at `5e9475bf`.
   When neither holds, the refusal is today's `rows_blocked_declared is <n>,
   but Spec declares <d> unreachable acceptance(s); shortfall is <k>`
   followed, when listed rows are unmatched, by `; no declaration names
   declared-blocked row <id> ("<criterion>")[, …]` for each. Every other
   check and its order is unchanged, so the four callers keep agreeing and
   every report accepted today is still accepted.
4. API Contract: one classifier. `detectMechanicalReportShape` counts a
   declared row exactly when `QADeclaredRowCriterion` accepts its status; its
   findings and counts do not change for any status.
5. API Contract: settle runs from the user checkout. When the working
   directory's branch starts with `roundfix/run-`, preflight refuses, exit
   `2`, nothing written, with `settle runs from the user checkout, not from
   <root> on branch <branch>; run it from <checkout>: roundfix settle --spec
   <slug> --task <task_id>`, where `<checkout>` is the `GitRoot` of the Run
   the branch names (the Run ID before any `-task_…` suffix) or, without that
   record, the repository's main worktree.
6. API Contract: settle follows the branch. When the selected surface is a
   kept Run's Task Worktree or Run Worktree and the target branch resolves in
   the user root, settle reads `H` (surface HEAD), `T` (branch tip) and their
   merge base `B`, and lists `git diff --name-only B T -- <Spec folder>`:
   - empty list or `H == T`: no change from today;
   - a Run Worktree with `B == H`: `git merge --ff-only T` in the surface,
     then stderr `Run Worktree advanced: <H short> -> <T short> (<branch>)`
     and the Spec and Task reloaded from the surface; when Git refuses,
     exit `2` naming the surface, the branch and Git's message, surface
     unchanged;
   - otherwise exit `2`, nothing written: `Spec <slug> changed on <branch>
     after settle surface <path> forked from it: <files>; merge <branch> into
     that worktree, then settle again`.
   It runs after the authorization and executable checks, so a refusal never
   follows an advance; a later Verification failure leaves the advance in
   place and the Task status unchanged.
7. API Contract: settle's executables. Preflight lists
   `git ls-files --others --exclude-standard -z` in the surface and refuses,
   exit `2`, nothing written, every path `UnauthorizedExecutables` returns:
   `refuse Settle commit: executable file <path> (mode <mode>)[, …] is not
   valid Work Item output unless the Spec's _authorization.md lists the path
   under executables`. `UnauthorizedExecutables` returns the drops of
   `filterStageablePaths` with `authorizedExecutables(resolution)` whose
   reason is `executable file`, so a listed path passes and a tracked path
   passes whatever its mode, as in the Task commit.
8. API Contract: the acceptance rule in writing (task_01). This sentence,
   verbatim and without inline code, goes into a new `### Declared rows name
   their criterion` section of the `qa-gate` skill (outside `### QA
   settlement`), the Roundfix `settle` reference's QA Report acceptance
   section, the `qa-report` guide and the `write-prd` PRD template: "One
   Unreachable Acceptance declaration may name several criteria, one
   criterion line each, and each criterion becomes its own QA row; acceptance
   matches every declared-blocked row to a declaration that names its
   criterion, by its text or its leading label." Each place adds that a
   report whose declared-blocked rows do not exceed its declarations is
   accepted as before. The `qa-gate` section also says to write the row as
   `blocked (declared: <criterion>)` with the criterion or its leading label,
   never a paraphrase; its non-waivable list
   and `qaGateContract` say "each criterion a PRD Unreachable Acceptance
   declaration names". `CONTEXT.md` revises **Unreachable Acceptance** with
   "One declaration may name several criteria; each is its own QA row."
9. API Contract: the checkout rule in writing (task_02). This sentence,
   verbatim, goes into `settleUsage`, the `settle` guide, the Roundfix
   `settle` reference and its `spec-delivery` recovery step: "Settle runs
   from the user checkout: inside a Run Worktree or Task Worktree it refuses
   and names the checkout." `CONTEXT.md` revises **Settle Command** with
   "It runs from the user checkout."
10. API Contract: the branch rule in writing (task_03). This sentence,
    verbatim, goes into `settleUsage`, the `settle` guide and the Roundfix
    `settle` reference: "When the Run's branch changed the Spec folder after
    the surface forked from it, settle fast-forwards a Run Worktree the
    branch is ahead of and refuses any other drift." Each place shows the
    `Run Worktree advanced:` line and says a later Verification failure
    keeps the advance. `CONTEXT.md` revises **Run Worktree** with "settle
    fast-forwards it to the Run's branch when that branch changed the Spec
    folder."
11. API Contract: the executables rule in writing (task_04). This sentence,
    verbatim, goes into the `settle` guide and the Roundfix `settle`
    reference: "Settle refuses an untracked executable file that the Spec's
    authorization does not list under executables, as a Task commit does; a
    tracked file whose mode changes is committed by both." `CONTEXT.md`
    revises **Authorized Executable** with "A settle commit keeps the same
    paths."

### Surface Transcripts

None. Every settle outcome needs a kept Run, which the QA gate cannot create
without an Agent; API Contracts 5-7 are asserted by named tests through CLI
dispatch, and `qa-report accept` is exercised by the QA gate with the built
binary.

### Declared breaks

No QA Report accepted today is refused (API Contract 3); a refusal keeps
today's words and may append the unmatched rows. The prototype moved no
existing eligibility, archive, settle or Daemon test.

1. `internal/spec/spec_test.go` compares declarations field by field, since
   `UnreachableDeclaration` holds a slice (test-only).
2. `internal/agent/spec_prompt_test.go` pins the new non-waivable sentence.
3. Settle inside a Run Worktree or Task Worktree exits `2` instead of
   committing there.
4. Settle may fast-forward a kept Run Worktree before Verification.
5. Settle refuses an untracked unlisted executable it used to commit.
6. `settle --help` changes; the Behavior Surface record moves (each Task).

### Invariants

1. Every report today's count accepts is accepted; a report the count
   refuses is accepted only when every declared-blocked row it lists is
   named by a declaration, and its refusal names the rows that are not.
2. The four eligibility callers and the mechanical stage read one rule.
3. Settle never writes before a refusal, never rewrites a Run commit, and
   never advances a Task Worktree or a diverged Run Worktree.
4. No settle commit keeps an untracked executable a Task commit of the same
   Run would refuse.

## Coverage Map

- Goal 1 → Declared rows matched by criterion (API Contracts 1-3)
- Goal 2 → Declared rows matched by criterion (API Contracts 3-4, Invariant 2)
- Goal 3 → Settle follows the branch (API Contract 6)
- Goal 4 → Settle runs from the user checkout (API Contract 5)
- Goal 5 → Settle commits what a Task commit keeps (API Contract 7)
- Goal 6 → The rules in writing (API Contracts 8-11)
- Core Features 1-2 → API Contracts 1-4
- Core Feature 3 → API Contract 6
- Core Feature 4 → API Contract 5
- Core Feature 5 → API Contract 7
- Core Feature 6 → API Contracts 8-11
- Success Metrics 1-3 → API Contracts 1-3, Invariant 1
- Success Metric 4 → API Contract 5
- Success Metric 5 → API Contract 6, Invariant 3
- Success Metric 6 → API Contract 7, Invariant 4
- Success Metric 7 → API Contracts 8-11

Deviations from the PRD: None.

## Absence Evidence

1. Several criteria per declaration — command: `grep -n 'Criteria' internal/spec/spec.go internal/spec/qa.go`;
   result: no match at `5e9475bf`.
2. Declared rows read from the Results table — command: `grep -rn 'DeclaredRows\|QADeclaredRow' internal`;
   result: no match at `5e9475bf`.
3. Settle's checkout refusal and branch advance — command: `grep -n 'RunBranchPrefix\|merge-base\|ff-only' internal/cli/settle.go`;
   result: no match at `5e9475bf`.
4. Settle's executable check — command: `grep -n 'executable\|ls-files' internal/cli/settle.go`;
   result: no match at `5e9475bf`.

## ADR Alternatives

- ADR-0343: written by this Spec; its rejected alternatives (matching
  alone, counting criteria, exact text only, reading the plan
  from the user checkout, refusing instead of fast-forwarding, settling from
  inside the worktree, a grant for tracked mode changes) are not taken.
- ADR-0240: kept as one rule in `internal/spec`; coverage widens to a match
  beside its count, and its rejected alternatives (a forced `pass`, trusting
  only Daemon-observed denials) stay rejected.
- ADR-0104: unchanged; outside-evidence rows other than the exempt kinds
  still hold Pull Request preparation.
- ADR-0167: unchanged; the pre-PR Pull Request row needs no declaration.
- ADR-0229: unchanged; the park decision reads environment rows only.
- ADR-0276: unchanged; the Network-Prohibited Row needs no declaration.
- ADR-0335: unchanged; the prohibited subject form stays one exempt kind.
- ADR-0237: unchanged; a retry still records a merge made outside the
  queue, and settle's fast-forward touches no queue item.
- ADR-0274: unchanged; a turn that did not finish is still refused before
  any QA Report is read.
- ADR-0295: unchanged; reachability still reads the fetched default branch.
- ADR-0312: unchanged; a resumed archived item still starts at review.
- ADR-0340: applied to settle; its rejected alternatives (staging every
  executable, naming executables in the Task file, globs) are not taken.
- ADR-0157: unchanged; a tracked path stays work output whatever its mode,
  in settle too.

## Integration Points

The Delivery Queue archives through `roundfix archive` and settles the QA
Task through Daemon settlement, both of which call `QAReportEligibility`, so
the matching reaches the queue unchanged. The Linux Lane runs the changed
`internal/spec`, `internal/speccheck`, `internal/agent`, `internal/daemon`
and `internal/cli` tests in `golang:1.26`; the settle tests need only Git.

## Testing Approach

- task_01: `internal/spec/qa_declared_rows_test.go` (parser, matcher table,
  eligibility of matched, unmatched and over-counted reports, and a
  paraphrased row accepted by today's count before and after) and
  `internal/cli/qa_report_declared_rows_test.go` through CLI dispatch
  (Quaestor's shape before and after the split).
- task_02: `internal/cli/settle_checkout_test.go` (Run Worktree, Task
  Worktree, no write, checkout named).
- task_03: `internal/cli/settle_plan_test.go` (advance with an amended PRD,
  Git overlap refusal, diverged refusal, Task Worktree refusal, no drift
  unchanged, advance kept after a Verification failure).
- task_04: `internal/cli/settle_executables_test.go` (unlisted refused, also
  in a new directory; listed committed `100755`; tracked made executable
  committed) and a Daemon test in
  `internal/daemon/task_engine_executables_test.go` (a tracked file made
  executable rides the Task commit).
- Each Task: phrase checks, owned-skill record, documentation contracts with
  the Behavior Surface record re-recorded last.
- QA sabotages, each on its own `git archive` copy (lesson 39; the input that
  makes the mutation observable is named in task_05).

## Build Order

1. task_01: declared rows matched by criterion and its written rule (API
   Contracts 1-4 and 8).
2. task_02: settle runs from the user checkout (API Contracts 5 and 9)
   (depends on: 1, for the shared `CONTEXT.md`, Roundfix skill files and
   records).
3. task_03: settle follows the branch (API Contracts 6 and 10) (depends on:
   2, for `internal/cli/settle.go` and the same shared files; its test uses
   task_01's matching).
4. task_04: settle commits what a Task commit keeps (API Contracts 7 and 11)
   (depends on: 3, for `internal/cli/settle.go` and the shared files).
5. task_05: the QA gate (depends on: 4).

## Risks & Considerations

- A gate on an older binary still counts; this Spec's own QA declares one
  criterion and records one row, which both rules accept.
- A paraphrased declared row matches nothing and passes only by today's
  count; the `qa-gate` text says to write the criterion or its label so the
  two-criteria case can match.
- A fast-forward runs Git in a kept worktree the user may have open; Git
  refuses rather than overwrite a local change, and settle relays its
  message.
- Briefing lesson 37 (one declaration, exactly one row) is superseded for
  Specs authored after this release.

## Glossary

- changes: **Unreachable Acceptance**
- changes: **Settle Command**
- changes: **Run Worktree**
- changes: **Authorized Executable**

## Decisions

- Repeated `criterion:` lines; match by text or leading label; today's count
  stays the fallback, so matching only widens acceptance (ADR-0343; operator
  review of 2026-10-27).
- The mechanical stage uses the same classifier.
- Settle refuses inside a Run or Task Worktree rather than settling there.
- Settle fast-forwards only on Spec-folder drift and only a Run Worktree; a
  branch ahead outside the Spec folder keeps today's `integration pending`.
- Settle reuses the Daemon's stage filter for executables; ADR-0157 stands
  for tracked mode changes.
