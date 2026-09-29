---
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# Gates that refuse only what someone can act on

## Executive Summary

The fix has three parts:

- **Recorded paths.** Every writer of a Task commit (the Daemon's
  `commitTask` and the Settle Command) records the ordinary paths the commit
  carries that the Task did not declare. The record goes in a Daemon-owned
  `## Recorded paths` section of the Task file and in the commit event.
- **The pre-PR Pull Request row.** QA eligibility derives the pre-PR Pull
  Request row from a report's validated Results rows and stops counting it
  against a qualifying declared `partial`.
- **The horizon.** `SC-ADR-RELATED` skips an ADR whose adding commit is not an
  ancestor of the Spec PRD's adding commit.

The skills, the user guide and the glossary follow.

The primary trade-off is in the first part. The QA scope audit stops being a
planning-accuracy check for ordinary paths and becomes a disclosure check.
Governed paths keep their refusal. The pre-PR review of the diff against the
Spec stays the judge of whether a disclosed path belongs to the slice. No
exported function changes its signature, and the Run Database schema is
untouched.

## Project Constraints

- Identifier strategy: applicable — new identifiers are the Task file heading
  `## Recorded paths` and the Task commit event payload key `recorded_paths`,
  following `## Carry-forward provenance` and the snake_case payload keys.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0057, ADR-0080, ADR-0088,
  ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0097, ADR-0104, ADR-0116,
  ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold; this Spec adds ADR-0166,
  ADR-0167 and ADR-0168; ADR-0169 and ADR-0170 do not apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-29 for the `qa-gate` and `write-tasks` skills and the standing grant
  of 2026-09-18 for the Roundfix skill, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Recorded paths

A new file `internal/spec/recorded_paths.go` owns the record, next to
`RecordCarryForward` in `internal/spec/spec.go`, which is its precedent:

```go
// RecordedPathsHeading names the Daemon-owned Task file section (ADR-0166).
const RecordedPathsHeading = "## Recorded paths"

// UndeclaredTaskPaths returns, sorted and unique, the committed paths that are
// not taskFile, not declared as interface: or creates: in task.Context, and
// not reported true by governed.
func UndeclaredTaskPaths(task Task, taskFile string, committed []string, governed func(string) bool) []string

// RecordTaskPaths writes paths as the Task file's recorded section, replacing
// an existing one and preserving every other byte. No paths removes an
// existing section and writes nothing else.
func RecordTaskPaths(taskPath string, paths []string) error

// RecordedTaskPaths reads the recorded section from Task file bytes.
func RecordedTaskPaths(content []byte) []string
```

`internal/spec` does not import `internal/speccheck`, so both callers pass
`speccheck.GovernedPath` as `governed`. An `instruction:` path the Task edited
is undeclared as an edit and is recorded. The section is written at the end of
the file, replacing any earlier one:

```text
## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/store/delivery_test.go`
```

Each path is one backticked bullet, repository-relative with forward slashes. A
path containing a backtick or a line break is refused with an error, because
the Daemon cannot record it faithfully. The Task parser, `CarryForwardInputs`,
the Context limit and every consistency detector ignore the section, because
none of them reads a heading other than their own.

**The Daemon.** In `internal/daemon/task_engine.go`, `taskCommitPreparation`
gains `recorded []string`.

- `prepareTaskCommit` computes it for every Task whose `Type` is not `qa`, with
  `spec.UndeclaredTaskPaths(task, <repo-relative Task file>, committedFiles,
  speccheck.GovernedPath)`.
- `committedFiles` is `stageable` with every untracked directory entry expanded
  into the files `git ls-files --others --exclude-standard -z -- <dir>` lists
  under it, so a new package records its files, not its directory.
- `commitTask`, which runs only for a completed Task, calls
  `spec.RecordTaskPaths` before `Committer.Commit` when `recorded` is not
  empty. The Task file is already in `stageable`, so the section rides in the
  same commit.
- The commit event payload gains `"recorded_paths": recorded` only when it is
  not empty.
- A failed Task is never committed, so it records nothing.
- A hook refusal leaves the section staged with the work.

**The Settle Command.** In `internal/cli/settle.go`, `settleTaskAndCommit`
already stages the surface with `git add --all` and reads the staged files with
`stagedSettlePaths`.

- After that read, and for a non-QA Task, it computes the recorded paths the
  same way, reading deleted paths too.
- It writes the section with `spec.RecordTaskPaths`, then stages the Task file
  again before `Committer.Commit`.
- A surface a Hook Refusal left behind already holds the Daemon's section, and
  settle replaces it with the same content.

Task Carry-Forward cherry-picks the settlement commit, so the section travels
with it. `CarryForwardInputs` reads the Task file from the settlement's parent,
which never holds the section.

## The pre-PR Pull Request row

`internal/spec/qa.go` adds:

```go
const (
	QANoOpenPullRequestStatus = "blocked (environment: no open Pull Request)"
	QAPullRequestRowSource    = "Pull Request row"
)
```

`QAReport` gains `RowsBlockedPrePullRequest int`. `readQAReport` derives it from
the body with a new `qaReportPrePullRequestRows(body []byte) int`. The function
walks the `## Results` section with the same fence, heading and table helpers
`qaReportHollow` uses. It reads each table whose header has both a `Status` and
a `Provenance` column, compared trimmed and case-insensitively. It counts a row
whose trimmed status equals `QANoOpenPullRequestStatus` case-insensitively and
whose provenance contains `QAPullRequestRowSource` case-insensitively. A table
without a `Provenance` column contributes nothing. Frontmatter counts are read
and validated exactly as before.

In `QAReportEligibility`, the `partial` branch replaces its environment check.
Let `excused` be the smaller of `RowsBlockedPrePullRequest` and
`RowsBlockedEnvironment`, and `outside` be `RowsBlockedEnvironment - excused`.
When `outside > 0`:

- with `excused == 0`, the error stays `rows_blocked_environment is %d; expected
  0`;
- otherwise the error is `rows_blocked_environment is %d, %d outside the pre-PR
  Pull Request row; expected 0 outside it`.

Every other check, and its order, is unchanged. A `partial` with no declared
row still refuses with `newest QA Report verdict is "partial"; expected "pass"`.
The `pass` branch is untouched. `roundfix qa-report accept`, QA settlement in
the Daemon, `roundfix settle` and `roundfix archive` all call this one function,
so all four change together.

The Daemon's QA prompt line in `internal/agent/spec_prompt.go` for a resolved,
absent Pull Request becomes:

```text
Pull Request: none open; record each Pull Request journey as blocked (environment: no open Pull Request) and name the Pull Request row in its provenance. That row alone never prevents a qualifying declared partial.
```

It keeps the `Pull Request: none open;` prefix that existing daemon tests read.
The mechanical stage's count cross-check in `internal/speccheck/mechanical.go`
is unchanged: the row still counts in `rows_blocked_environment`.

## The related-ADR horizon

A new file `internal/speccheck/adr_horizon.go`:

```go
type adrHorizon struct {
	repoRoot  string
	prdCommit string            // commit that added the Spec's _prd.md
	addedBy   map[string]string // repository-relative ADR path -> adding commit
}

// newADRHorizon returns false when the full check must apply.
func newADRHorizon(repoRoot, prdPath string) (adrHorizon, bool)
func (horizon adrHorizon) predates(adrPath string) bool
```

`newADRHorizon` runs Git with `mechanicalGitEnvironment()` and returns `false`
when any of the following holds:

- `repoRoot` and the PRD's directory resolve to different
  `git rev-parse --path-format=absolute --git-common-dir`;
- `git -C repoRoot log -1 --diff-filter=A --format=%H -- <prd>` prints nothing,
  which means the PRD is uncommitted;
- any Git command fails.

It fills `addedBy` from one `git log --diff-filter=A --format=%x00%H
--name-only -- docs/adr`, keeping the first, newest commit per path.

`predates` returns:

- `false` for an ADR with no adding commit;
- `true` when `git merge-base --is-ancestor <adrCommit> <prdCommit>` exits `0`;
- `false` when it exits `1`;
- `true` on any other failure, so the check stays full.

`detectADRConsistency` builds the horizon lazily, at most once, on the first
related candidate. It skips a candidate whose ADR `predates` reports `false`.
Nothing else in the function changes, including `SC-ADR-UNLISTED` and
`SC-CITATION-UNSUPPORTED`. Temporary fixtures without Git keep the full check,
so existing characterizations and the corpus golden do not move.

## Skills, guides and glossary

- The `qa-gate` skill (`.agents/skills/qa-gate/SKILL.md`, canonical) makes two
  changes:
  - its table row for a qualifying declared `partial` and its Pull Request
    journey rule say that the pre-PR Pull Request row, recorded as
    `blocked (environment: no open Pull Request)` with the Pull Request row in
    its provenance, never decides a qualifying `partial` and needs no
    declaration;
  - a scope rule counts a path a Task file lists under `## Recorded paths` as
    declared and names it in the scope row, while a Governed Path still needs
    its authorization.
- The `write-tasks` skill says, after its declared-path rule, that the Daemon
  records an undeclared path at commit and that recording discloses and
  reserves nothing. It keeps every phrase
  `TestWriteTasksSkillStatesTheDeclaredPathRules` requires.
- Both skills move both version declarations to `0.0.3`.
- The Roundfix skill's settlement and archive eligibility text names the
  pre-PR row exception, and its Task commit and Settle Command text names the
  recorded section.
- `make skills-sync` regenerates the mirrors, and `make baseline-digests` runs
  after the edits.
- `docs/user-guide/commands.md`:
  - The Task commit and Settle Command text name the recorded section. Task
    `task_01` does this.
  - The `archive` refusal list and the `qa-report accept` text name the pre-PR
    row. Task `task_04` does this.
- `docs/user-guide/context-driven-development.md` states the Pull Request row
  exception and the horizon.
- `CONTEXT.md` gains **Recorded Path**, and **QA Report** and **Spec
  Consistency Check** state the exception and the horizon.

## API Contracts

1. API Contract: Task file record — a completed non-QA Task's commit carries
   `## Recorded paths` listing exactly the committed ordinary paths the Task
   did not declare; no section when there are none; the authored `## Context`
   is byte-identical before and after.
2. API Contract: Task commit event — the `daemon.commit` event for a Task
   carries `recorded_paths` with the same paths when any exist, and no such key
   otherwise.
3. API Contract: `roundfix settle --spec <slug> --task <task_id>` — the Task commit it creates
   carries the same section for the same rule; stdout and exit codes are
   unchanged.
4. API Contract: `roundfix qa-report accept <path>` — exit `0` for a declared
   `partial` whose only environment-blocked rows are pre-PR Pull Request rows;
   exit `1` with the new message when another environment-blocked row remains,
   and with the unchanged message when none of them is the pre-PR row.
5. API Contract: `roundfix spec check` — `SC-ADR-RELATED` is not reported for
   an ADR outside a committed Spec's horizon; exit codes, flags and other codes
   are unchanged.

## Coverage Map

- Goal 1 → Recorded paths; API Contracts 1-3.
- Goal 2 → The pre-PR Pull Request row; API Contract 4.
- Goal 3 → The related-ADR horizon; API Contract 5.
- Goal 4 → Recorded paths; The pre-PR Pull Request row; The related-ADR
  horizon.
- Core Feature 1 → Recorded paths; Skills, guides and glossary.
- Core Feature 2 → The pre-PR Pull Request row; Skills, guides and glossary.
- Core Feature 3 → The related-ADR horizon; Skills, guides and glossary.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2, Testing Approach 5.
- Success Metric 3 → Testing Approach 3, Testing Approach 5.
- Success Metric 4 → Testing Approach 4.
- API Contracts 1-3 → Recorded paths.
- API Contract 4 → The pre-PR Pull Request row.
- API Contract 5 → The related-ADR horizon.

## Integration Points

- **Task commit writers.** `commitTask` and `settleTaskAndCommit` are the only
  code that creates a Task's standard commit. Carry-forward replays an existing
  settlement commit and adds only `## Carry-forward provenance`.
- **QA eligibility.** One function serves the Daemon, `settle`, `archive` and
  `qa-report accept`; the mechanical stage keeps validating the counts.
- **Spec Consistency Check.** The horizon lives inside one detector, so `spec
  check`, the Delivery Plan, queue revalidation and the QA precondition inherit
  it without change.

## Testing Approach

1. **Recorded paths.**
   - Spec unit tests in the new `internal/spec/recorded_paths_test.go`:
     undeclared, declared, governed, `instruction:`, Task file and duplicate
     inputs; replacing a section; preserving every other byte; refusing an
     unrecordable path; a Task file with the section parsing to the same
     Context and the same `CarryForwardInputs`.
   - Daemon tests in the new `internal/daemon/recorded_paths_test.go` over the
     existing Task cycle fixture (`taskFakeRunner` on a `gittest` repository).
     The 0180 shape records its test file. A new untracked package records its
     files. Every-path-declared, governed, QA and failed Tasks record nothing.
     The event payload carries `recorded_paths` only when paths are recorded.
   - A Settle Command test in the new
     `internal/cli/settle_recorded_paths_test.go`.
2. **Pre-PR row.** New `internal/spec/qa_prepr_row_test.go` covers:
   - the qualifying `partial` accepted;
   - another environment row refused with the new message;
   - a no-PR status without the Pull Request source refused with the unchanged
     message;
   - a table without provenance gaining nothing;
   - `pass` unchanged.

   `TestQAReportEligibilityKeepsExistingPartialRefusalPrecedence` stays green
   unchanged, and `internal/agent/spec_prompt_test.go` pins the new prompt
   line.
3. **Horizon.** New `internal/speccheck/adr_horizon_test.go` over `gittest`
   repositories:
   - an ADR committed after the PRD opens no gap;
   - an ADR committed before the PRD, or in the same commit, opens the gap;
   - an uncommitted PRD keeps the full check;
   - an uncommitted ADR against a committed PRD opens no gap;
   - a directory without Git keeps the full check;
   - `SC-ADR-UNLISTED` still fires for a cited ADR outside the horizon.
4. **Documentation.** Phrase checks on the canonical skills, their mirrors, the
   user guide and `CONTEXT.md`, `make skills-sync-check`, and
   `TestWriteTasksSkillStatesTheDeclaredPathRules`.
5. **Real history through the built binary.** The QA gate builds the candidate.
   - It runs `roundfix qa-report accept` on Spec 0179's archived
     `qa-report-2026-09-29.md` and on Oraculum Spec 0027's mirrored report.
   - It replays the ADR-0161 cascade in a disposable clone at `ebeb997f`,
     against a binary built from the starting main.

## Build Order

1. Recorded paths in the spec package, the Daemon and the Settle Command, with
   the commands guide text for Task commits and settle, task_01 (depends on:
   none).
2. The pre-PR Pull Request row in QA eligibility and the QA prompt, task_02
   (depends on: none).
3. The related-ADR horizon, task_03 (depends on: none).
4. Skills, remaining guide text and glossary, task_04 (depends on: 1, 2, 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Daemon tests whose fake Agent writes undeclared files.** Several existing
  daemon tests write paths their fixture Task does not declare. Those Task
  files now gain a section. A test that compares Task file bytes exactly may
  change. task_01 updates only the tests the change invalidates and names each
  in its Result.
- **This Spec runs on the v0.20.0 Daemon.** Its own Tasks get no record, so
  they declare every path they expect to change. Its QA scope row discloses,
  rather than refuses, a `_test.go` file a Task's Result names as invalidated
  by the change. That is the policy this Spec ships, applied by hand once.
- **Gaming the pre-PR exception.** A QA Agent could label another block as the
  Pull Request row. The exception requires both the exact status and the Pull
  Request source in provenance. The mechanical stage still validates every row,
  and the pre-PR review reads the report.
- **Git cost in the consistency sweep.** The horizon runs Git only when a
  related candidate exists, at most one `log` per Spec and one per repository
  call, and it stays inside the corpus budget test.
- **Spec 0182 shares no file with this Spec.** Its ADR-0169 and ADR-0170 are
  listed here as not applying.

## Decisions

- Record in a Daemon-owned section, not in `## Context`. See ADR-0166.
- Derive the pre-PR row from validated Results rows, without a new frontmatter
  key. See ADR-0167.
- Use commit ancestry, not dates, for the horizon. See ADR-0168.
- Keep the mechanical stage's counts and the `pass` rule unchanged.
- Research: the Secondbrain and Exa sources are recorded in the PRD's Research
  basis. The published scope checkers (TaskBound, AgentScope, agentdiff,
  agent-guardrails) informed the disclosure-versus-refusal split.
