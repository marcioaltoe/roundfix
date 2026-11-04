---
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
prd: _prd.md
created: 2026-11-04
---

# A retry after a QA precondition refusal runs the whole QA gate — Technical Spec

## Executive Summary

The QA Report Import of ADR-0194 brings the newest unintegrated Daemon QA
Report commit of a Spec into the next pass, whatever its report records. A
Precondition Refusal report is such a commit, so a retry imports a refusal as
its previous report: the import adds no row to carry, and in a repository
whose Verification checks the formatting of `qa/` files without a Format
Command configured, it makes the next precondition refuse on the imported
bytes. The design adds one refusal to `copyPriorQAPass`: a candidate commit
that adds or modifies a report recording a Precondition Refusal is refused
with the reason `precondition refusal: <path>`, so the pass runs in full
with no previous report. The format step is unchanged; a new test pins that
the refusal report goes through it. The implement reference, the
configuration guide, `CONTEXT.md` and ADR-0355 state the rules.

The trade-off accepted: a candidate that mixes a matrix report with an older
imported refusal, which only a Roundfix release without this Spec could
write, is refused whole, so that pass carries no row and re-runs every one.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the existing
  `refused` outcome of the `prior_report` event gains a reason text. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, endpoint or
  credential changes. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0194 owns the import, ADR-0194: "a
  failed pass on an unintegrated Run Branch is the next pass's previous
  report"; ADR-0210 refines it, ADR-0210: "This refines ADR-0194 without
  superseding it"; ADR-0249 owns the Format Command, ADR-0249: "An empty value
  runs nothing, and that is today's behavior."; ADR-0132 keeps the refusal
  report as the refused gate's record, ADR-0132: "it emits one terminal row
  naming the refusal that produced it"; ADR-0355 is new. ADR Alternatives
  below relates each. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AX answers of 2026-11-04 quoted in the authorization record; the Governed
  Paths are the Roundfix skill files task_02 changes; no Makefile, workflow,
  module, lint or Baseline change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-04 at `6a0ed904`; `_prd.md` → Validation records the runs.

1. `runQAGate` (`internal/daemon/task_engine.go`) calls `importPriorQAPass`
   before the mechanical stage. `findRecordedPriorQAPass`
   (`internal/daemon/qa_prior_pass.go`) keeps a commit only when a `created`
   QA commit event of an Implement Run of the Spec names it, the Run's Run
   Branch still exists and holds it, its subject is
   `docs: qa report for <slug> (`, it carries `Roundfix-Spec` and no
   `Roundfix-Task`, and it is not an ancestor of `HEAD`; the newest such
   commit by commit date wins. Its files are the `qa/` paths the commit adds
   or modifies; its report is the newest of them. No line reads the report's
   content.
2. `copyPriorQAPass` skips compiled-source evidence files, refuses with
   `not newer than <path>`, `dated after today`, `path differs: <path>` or
   `report shape: <code>`, writes the files it lacks, and removes what it
   created when it refuses. `importPriorQAPass` publishes
   `Prior QA Report <outcome> for Spec <slug>.` with phase `prior_report`,
   outcome `none`, `imported` or `refused`, and the reason.
3. An imported report becomes the pass's previous report
   (`buildQAPromptContext`), the Format Command runs over its files at stage
   `imported`, and the files ride in the pass's QA Report commit
   (`prepareQAReportCommit`).
4. A refused repository Verification sets `PreconditionRefusal`, the gate
   skips the Agent, `mechanicalQAReportContent` writes the report through
   `spec.WritePreconditionRefusalReport` (`rows_blocked_precondition: 1`), and
   the commit path runs `formatQADirectory` at stage `commit` before
   `commitPreparedQAReport`, as for every report. The strict Spec check
   refusal writes the same shape.
5. Task Carry-Forward (`carryForwardTaskCommits`, `internal/cli/carryforward.go`)
   lists only commits whose `Roundfix-Task` trailer names a Task settled
   `completed`; a QA Report commit has no such trailer, so it never leaves its
   Run Branch. A retry whose item worktree already holds every settled Task as
   `completed` carries nothing and prints no `Carried forward from Run` line,
   and the next Run's scheduler skips completed Tasks (`TaskAwaitsWork`).
6. `.agents/skills/roundfix/references/implement.md` names "an imported pass"
   only in the Format Command paragraph; nothing says what the import brings
   in or never brings in.

## System Architecture

No package, command, flag or file outside the listed ones is added. Changes
by component:

- `internal/daemon/qa_prior_pass.go`: the refusal check in `copyPriorQAPass`.
- `internal/daemon/qa_prior_pass_test.go` and
  `internal/daemon/qa_format_test.go`: the new tests.
- The implement reference, the configuration guide, `CONTEXT.md` and
  ADR-0355.

## Implementation Design

### Interfaces

```go
// internal/daemon/qa_prior_pass.go — signature unchanged
func (engine *Engine) copyPriorQAPass(ctx context.Context, plan TaskPlan, pass *priorQAPass) (reason string, err error)

// The reason text the new refusal returns, path as listed in pass.Files:
//   "precondition refusal: " + path
```

### Data Models

None. The `prior_report` event keeps its keys and its three outcomes; only
the reason text is new.

### API Contracts

1. API Contract: the refusal check. In `copyPriorQAPass`, after the copy loop
   and before `speccheck.ReportShapeFindings`, each path of `pass.Files`
   whose base name starts with `qa-report-` and ends with `.md` is read with
   `spec.ReadQAReportFile` from the Run Worktree. The first, in `pass.Files`
   order, that reads without error and has `RowsBlockedPrecondition` greater
   than zero makes the function return `precondition refusal: <path>`; the
   deferred cleanup then removes every file and directory the import created,
   as for the other refusals. A report that does not read keeps today's
   handling: the shape detector decides. `importPriorQAPass` is unchanged and
   publishes outcome `refused` with that reason and `files` null.
2. API Contract: no other candidate. The refusal does not make
   `findRecordedPriorQAPass` look at an older commit; the pass runs with no
   previous report, as after any other refusal.
3. API Contract: the format path. `runQAGate`, `prepareQAReportCommit`,
   `formatQADirectory` and `commitPreparedQAReport` keep their bytes; the
   Precondition Refusal report is formatted at stage `commit` under the
   existing restore rules.
4. API Contract: the implement reference. After the Format Command paragraph
   of step 5, both copies of `references/implement.md` state:
   1. that the Format Command covers every QA Report the Daemon commits, the
      Precondition Refusal report included, that without `verification.format`
      the Daemon commits its reports as written, and that a repository whose
      Verification checks the formatting of `qa/` files configures it;
   2. the QA Report Import: before a pass builds its matrix, the Daemon copies
      byte for byte the newest Daemon QA Report commit of the Spec that a
      recorded Run Branch holds and the Run's history lacks, failed passes
      included, as the pass's previous report; a pass always runs the whole
      gate, and an imported report only lets rows whose Evidence Snapshot
      proves their inputs unmoved carry;
   3. that the import never imports a Precondition Refusal and refuses with
      `precondition refusal: <path>`, and what else it never brings in: a
      commit already in the Run's history, one no recorded Run Branch holds,
      a Task commit or another Spec's commit, a report not newer than the
      newest one in the worktree, one dated after today, one the shape
      detector refuses, a path that differs from an existing file, and
      compiled-source evidence files, which are skipped one by one;
   4. how to start a pass with no previous report: after the earlier Run's
      Tasks are carried forward, remove that Run's Run Worktree and Run
      Branch; this discards its unintegrated commits, so it is a destructive
      Git operation that needs the maintainer's authority, and a refusal
      never needs it.
5. API Contract: the configuration guide. The Format Command section of
   `docs/user-guide/configuration.md` gains the sentence of API Contract 4.1.
6. API Contract: the glossary and the ADR. `CONTEXT.md` adds
   **QA Report Import** and revises Precondition Refusal: a later Run never
   imports it. ADR-0355 is created from the ADR-0355 section below.

### Surface Transcripts

None. No command's standard output, standard error or exit code changes; the
change is a Run Event reason and the files a pass commits.

## Invariants

```text
1. A candidate whose added or modified reports record no Precondition Refusal is imported with the same files, bytes, skips and event payload as before.
2. No Precondition Refusal report reaches a later Run's worktree or QA Report commit through the import.
3. A refused import leaves no file or directory it created.
4. Every QA Report commit, a refusal's included, runs the Format Command under the same restore rules.
5. Task Carry-Forward never carries a QA Report commit.
```

## Declared breaks

- A recorded refusal commit was imported (`imported`); it is now refused
  (`refused`, `precondition refusal: <path>`). No existing test pins the old
  outcome: `rg -n 'rows_blocked_precondition' internal/daemon/qa_prior_pass*_test.go internal/daemon/qa_two_pass_carry_test.go`
  finds nothing.
- No existing test, golden or fixture changes.

## ADR-0355

task_02 creates
`docs/adr/0355-a-precondition-refusal-is-never-imported-as-a-prior-qa-pass.md`
with the lifecycle overlay (`status: accepted`, `created_at` and `updated_at`
`2026-11-04T12:00:00Z`, `deprecated_at: null`, `superseded_by: null`), the
title "A precondition refusal is never imported as a prior QA pass", and:

- Context: ADR-0194 imports the newest unintegrated QA Report commit so that a
  failed pass becomes the next pass's previous report; ADR-0249 rejected not
  importing a failed pass and formats the QA directory instead. On 2026-11-03
  Pantheon, which configures no Format Command, retried after a precondition
  refusal; the import brought the refusal in, the precondition refused on its
  bytes, and the Run recommitted it with no Agent session.
- Decision: a candidate whose added or modified reports record a Precondition
  Refusal is refused, with no older candidate searched; a refusal measured no
  requirement, so it is not a pass and has nothing to hand on; the refusal
  report keeps going through the Format Command at its commit; the refusal
  commit is never carried forward.
- Considered Options, each rejected with its reason: importing it as today;
  searching an older candidate; dropping only the refusal file from a mixed
  commit; carrying the refusal commit forward and superseding it; a flag to
  discard a prior report.
- Consequences: this narrows ADR-0194 without superseding it; a pre-fix mixed
  commit is refused whole; Pantheon's manual workaround is no longer needed.

## Coverage Map

- Goal 1, Goal 2, Core Feature 1, Core Feature 2, Success Metric 1,
  Success Metric 2 → the refusal check (API Contracts 1 and 2), Invariants 2 and 3,
  `TestPriorQAPassRefusesAPreconditionRefusal` and
  `TestARetryAfterAPreconditionRefusalRunsTheWholeGate`.
- Goal 3, Core Feature 3, Success Metric 3 → Invariant 1 and the existing
  import, carry and format tests.
- Goal 4, Core Feature 4, Success Metric 4 → the format path (API Contract 3),
  Invariant 4, `TestQAFormatFormatsAPreconditionRefusalReport`.
- Goal 5, Core Feature 5, Success Metric 5 → the import rules written into both
  copies of the implement reference (API Contract 4), the Format Command
  sentence of the configuration guide (API Contract 5), and the QA Report
  Import and Precondition Refusal entries of `CONTEXT.md` with ADR-0355 (API
  Contract 6).

Deviations from the PRD: None.

## Absence Evidence

1. A Precondition Refusal check in the import — command: `grep -n 'RowsBlockedPrecondition\|ReadQAReportFile\|precondition' internal/daemon/qa_prior_pass.go`; result: no match; the import never reads a report's frontmatter.
2. The import rules in the implement reference — command: `grep -n -i 'never imports\|previous report\|Precondition Refusal' .agents/skills/roundfix/references/implement.md`; result: no match.
3. A test of the refusal report's format path — command: `grep -n 'PreconditionRefus\|rows_blocked_precondition' internal/daemon/qa_format_test.go`; result: no match.
4. The glossary term — command: `grep -n 'QA Report Import' CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0194: narrowed; the rule it replaced, an establishing head that must
  be an ancestor of the current head, stays replaced, and a failed pass that
  built its matrix is still imported; only a refusal, which carries no row,
  is excluded.
- ADR-0210: followed; the per-file record it replaced stays replaced; a
  refusal records no Evidence Snapshot, so the exclusion changes no digest.
- ADR-0249: followed; its rejected alternative "Not importing a failed pass"
  stays rejected, because a Precondition Refusal is not a pass; its other
  rejected alternative, excluding imported files from the precondition, stays
  rejected; formatting stays the repair for every committed report.
- ADR-0132: followed; its rejected alternative, writing no report for a
  refused gate, stays rejected: the refusing Run still writes and commits the
  report, and a later Run only declines to import it.
- ADR-0355: new; records the exclusion and the options above.

## Integration Points

None. The change reads files the import already wrote and Git objects it
already read.

## Testing Approach

All tests are hermetic: temporary Git repositories and Run Databases from
`priorFixture`, the fake Agent runner, shell scripts as Format Command and
repository Verification, no network. Each new top-level test calls
`t.Parallel()` first.

- `internal/daemon/qa_prior_pass_test.go`:
  - `TestPriorQAPassRefusesAPreconditionRefusal`, a table:
    - `repository Verification refusal`: one recorded commit adding a report
      written by `spec.WritePreconditionRefusalReport` with check `make verify`;
    - `strict Spec check refusal`: the same with check
      `speccheck.GatePreconditionCheck`;
    - `older refusal beside a matrix report`: one recorded commit adding the
      refusal as `priorReportName(1)` and `priorReportContent("pass")` as
      `priorReportName(2)`.
    Each asserts `imported` false, the last `prior_report` event `refused`
    with `precondition refusal: <refusal path>`, and no `qa/` directory left.
  - `TestARetryAfterAPreconditionRefusalRunsTheWholeGate`: the Pantheon shape.
    A recorded refusal commit; a repository Verification
    `! grep -rqF -- '| - | --- |' docs/specs/<slug>/qa` run by `ExecVerifier`;
    `GitCommitter`. `TaskCycle` publishes `refused`, prompts the Agent once,
    settles `pass` and accepted, and `HEAD`'s QA Report commit lists none of
    the refusal commit's files.
- `internal/daemon/qa_format_test.go`:
  `TestQAFormatFormatsAPreconditionRefusalReport`, subtests:
  - `formatter runs`: repository Verification `exit 3`, an appending
    formatter; the committed newest report has `rows_blocked_precondition` 1
    and ends with the formatter's bytes, and the last format event is
    `commit` `formatted`;
  - `formatter fails`: a formatter that rewrites the file and exits 1; the
    committed report lacks its bytes, the event is `commit` `failed`, and the
    `roundfix: QA format commit failed:` diagnostic is printed.
- Existing suites that consume the changed contract, unchanged:
  `TestPriorQAPass`, `TestAPriorPassIsImportedOnlyFromARecordedRunCommit`,
  `TestQAGateCarriesARowFromAnUnintegratedFailedPass`,
  `TestQAPromptNamesTheImportedPassHead`, `TestTwoGatePasses`,
  `TestQAFormat` and `TestQAGateRefusesOnFailedRepositoryVerification`.
- Sabotages, each run on a temporary copy, each with the case that makes it
  observable:
  - the refusal check removed: every `TestPriorQAPassRefusesAPreconditionRefusal`
    case fails, and `TestARetryAfterAPreconditionRefusalRunsTheWholeGate`
    fails with no Agent prompt;
  - only the newest report checked: the `older refusal beside a matrix
    report` case fails, because its newest report is the matrix;
  - every `fail` verdict refused: `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte`
    and `TestQAGateCarriesARowFromAnUnintegratedFailedPass` fail on their
    `verdict: fail` matrix reports;
  - the Format Command skipped when the newest report records a refusal
    (`formatQADirectory` returning early): the `formatter runs` subtest
    fails on the committed bytes.

## Build Order

1. The refusal check and its tests, with the pinned format path (task_01).
2. The implement reference and its mirror, the configuration guide,
   `CONTEXT.md`, ADR-0355, the Roundfix skill version and the Behavior Surface
   Record last (task_02) (depends on: 1).
3. QA gate (task_03) (depends on: 2).

task_02 needs task_01 because it documents the reason text and the rule
task_01 makes true. The Tasks share no file.

Files Specs 0346, 0347 and 0348, authored in parallel and delivered after
this Spec, may also change, each kept to a narrow edit here so they rebase
cleanly: the Roundfix `SKILL.md` pair (version fields only),
`skills/testdata/owned-skill-versions.json`,
`docs/references/behavior-surfaces.json` and `CONTEXT.md` (one new entry and
one revised sentence). This Spec owns the implement reference; it touches no
other reference, no Baseline asset, no `docs/agents/` file, no
`internal/speccheck` or `internal/config` file and no other skill.

## Risks & Considerations

- An adopter with a pre-fix mixed commit on a Run Branch loses carried rows
  once; the gate re-runs them, which costs time, never a wrong verdict.
- Pantheon's report stayed unformatted because its Project Config sets no
  `verification` key; after this Spec its retry no longer imports the refusal,
  but its later QA Reports still need a Format Command to pass its own
  formatter check. The reply to Pantheon says so.
- Secondbrain: `wiki/index.md` and the `qmd` query of `_prd.md` → Validation
  returned only Roundfix and skills mirrors, which agree with ADR-0132 and
  ADR-0249 and say nothing about importing a refusal; nothing in it changes
  the design.
- Exa: the Bazel passage of `_prd.md` → Acceptance evidence supports never
  reusing a failure as the next run's result.

## Glossary

None.

## Decisions

- A refusal is recognized by `rows_blocked_precondition` greater than zero.
  See ADR-0355.
- Any report the candidate adds or modifies decides, not only the newest.
  See ADR-0355.
- No older candidate is searched after the refusal. See ADR-0355.
- The format step is pinned, not changed (API Contract 3).
- No new command; the reference names the destructive manual path.
