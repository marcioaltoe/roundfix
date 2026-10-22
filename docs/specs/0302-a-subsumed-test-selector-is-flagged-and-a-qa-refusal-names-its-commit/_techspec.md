---
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
prd: _prd.md
created: 2026-10-22
---

# A subsumed test selector is flagged and a QA refusal names its commit — Technical Spec

## Executive Summary

No structural Verification check in `internal/speccheck` reads test runner
arguments, so a pytest node id passed beside its own file reaches the prober,
which can only say `vacuous`. A new detector in
`internal/speccheck/verification_selector.go`, wired beside the Verification
path detector, reports `SC-VERIFY-SUBSUMED-SELECTOR` for each node id another
collection argument of the same pytest invocation contains, under pytest 9's
own subsumption rule (ADR-0315). Separately, `spec.WritePreconditionRefusalReport`
and the blocking branch of the Daemon's `mechanicalQAReportContent` write no
`build`; a new `AuditedCommit` field of `spec.AuditorEvidence`, filled by the
Daemon from the Run Worktree's `HEAD`, is written as `build:` right after
`verdict:` in both reports, and the pending seed keeps its bytes.

The trade-off accepted: the selector rule is pytest-only and lexical. A path
spelled two ways that pytest resolves to one file (a symbolic link, an
absolute and a relative form) is not compared, and the check stays silent
there; the prober still calls such a command vacuous.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the check code is a constant and `build` is an existing report
  key. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0315 and
  keeps ADR-0148, and Build Order 3 extends ADR-0132's refusal report, each
  named with its title in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("same as cycle AO"), carrying the cycle AJ grant
  "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2`. Line numbers move when
other Specs merge; the function names do not.

1. `Check` runs the structural Verification detectors from
   `internal/speccheck/constraints.go` (`detectVerificationPaths` at `:205`)
   and `internal/speccheck/citations.go` (`detectTaskCoverageAndContextReferences`:
   `WorkIndependentVerification`, `InvertedExitVerification`,
   `WrapFragileVerification`, `nonHermeticVerificationCommand`,
   `VacuousVerificationCommands`). The stage table `stagedDetectors`
   (`internal/speccheck/coherence.go:48`) lists each code with the stage that
   can decide it. None of them recognizes a test runner or the `::` of a
   node id: `rg -n 'pytest' internal/speccheck -g '!*_test.go'` finds only
   the `pytest.ini` pattern of `governed.go`.
2. `roundfix spec check --run-verification` runs every command through
   `daemon.ProbeCommands` (`internal/daemon/verification_probe.go:45`) in a
   disposable checkout and reports `vacuous` with the fixed text `(exited zero
   before work)` (`renderSpecCheckText`, `internal/cli/spec_check.go`); the
   verdict carries no cause.
3. pytest, measured in a throwaway virtual environment: 9.1.1 exits 0 for a
   missing node id passed with its file, a directory above it, `.` or
   `./<file>`, also through `python -m pytest`; it exits 4 for the node id
   alone, beside a sibling file, beside another node id of the same file, or
   with `--keep-duplicates`. 8.3.5 and 7.4.4 exit 4 in every case. The
   adopted Backlog Entry lists every command.
4. `spec.WritePreconditionRefusalReport` (`internal/spec/qa.go:109`) writes
   `verdict`, `auditing_binary`, `auditor_staleness`, the four blocked counts
   and the two precondition keys; no `build`. Its only production caller is
   `mechanicalQAReportContent` (`internal/daemon/task_engine.go:3509`), which
   for a blocking result without a precondition refusal writes its own
   frontmatter (`verdict`, auditor keys, three counts), also without `build`,
   and for a non-blocking result writes the pending seed the Agent finishes.
5. `writeMechanicalQAReport` (`:3351`) resolves the evidence through
   `resolveQAAuditorEvidence` (`:3388`) before rendering; the gate resolves
   the audited head only afterwards (`qaAuditedHead`, `:3802`, called at
   `:2910`), for the `user_flow_binary` prefix check. A blocking result never
   reaches an Agent (`if !mechanicalResult.Blocking`, `:2942`), so the
   Daemon's report is the final one.
6. `readQAReport` (`internal/spec/qa.go:392`) unmarshals a fixed set of keys
   and ignores others, so a report with `build` stays readable; nothing reads
   `build` today. `parseMechanicalReport` (`internal/speccheck/mechanical.go`)
   reads only the blocked counts and evidence snapshots from the frontmatter.
7. Tests that pin those bytes: `TestWriteMechanicalQAReportWritesThePreconditionRefusal`
   (`internal/daemon/task_engine_test.go:2017`) pins the refusal frontmatter
   prefix and, in "a gate that refused nothing writes the report it wrote
   before", the whole blocking report; both run on a fixture whose Run
   Worktree is a Git repository. The `internal/spec` refusal tests and the
   report sequence tests check substrings or run without a Git repository.

## System Architecture

- `internal/speccheck/verification_selector.go` (new): the code constant and
  `detectSubsumedSelectors`.
- `internal/speccheck/constraints.go`: one call beside `detectVerificationPaths`.
- `internal/speccheck/coherence.go`: one `stagedDetectors` entry.
- `internal/spec/auditor_evidence.go`: the `AuditedCommit` field.
- `internal/spec/qa.go`: the `build` line of the refusal report.
- `internal/daemon/task_engine.go`: `resolveQAAuditorEvidence` and
  `mechanicalQAReportContent`.
- `docs/user-guide/commands/spec.md`, the Roundfix `spec` reference, the
  `write-tasks` and `qa-gate` skills,
  `docs/user-guide/context-driven-development.md`, `CONTEXT.md`.

No package, import edge or schema changes.

## Implementation Design

### Interfaces

```go
// CodeVerifySubsumedSelector identifies a Verification command that passes a
// test selector together with a collection argument that contains it.
const CodeVerifySubsumedSelector = "SC-VERIFY-SUBSUMED-SELECTOR"

func detectSubsumedSelectors(result *Result, repoRoot, specDir string, graph *spec.Graph) error

type AuditorEvidence struct {
	// existing fields unchanged
	AuditedCommit string
}
```

### Data Models

`spec.AuditorEvidence` gains `AuditedCommit`. `spec.QAReport` and
`spec.PreconditionRefusal` are unchanged.

### API Contracts

1. API Contract: the invocation. In each Verification command of a Task whose
   status is `pending`, that is not the manifest's QA Task and whose type is
   not `qa`, the detector reads the command and the body of each `$(...)`
   substitution, splits each into simple commands at `;`, `&&`, `||`, `|`,
   `&`, `(`, `)`, `{` and `}` (the existing `shellWords` tokenizer), and
   recognizes a pytest invocation at command position after any of `!`,
   `rtk`, `env`, `exec` and `NAME=value` words: `pytest` or `py.test` by base
   name (so `.venv/bin/pytest` counts), an interpreter whose base name matches
   `python[0-9.]*` followed by `-m pytest`, or `uv`, `poetry`, `pdm`, `pipenv`
   or `hatch` followed by `run` and, at any later word, one of those runners.
2. API Contract: the collection arguments. After the runner, a word is a
   collection argument unless it starts with `-`, is a redirection operator or
   its target, is the value of `-k`, `-m`, `-p`, `-c`, `-o`, `-W`, `-n`, `-r`,
   `--deselect`, `--ignore`, `--ignore-glob`, `--rootdir`, `--confcutdir`,
   `--basetemp`, `--maxfail`, `--tb`, `--durations`, `--import-mode`,
   `--junitxml`, `--junit-xml`, `--cov`, `--cov-report`, `--log-level` or
   `--override-ini`, or holds `*`, `?`, `[`, `$` or a backtick without `::`.
   An invocation with `--keep-duplicates` is not judged. An argument splits at
   `::` into a path, cleaned with `path.Clean` after a trailing `/` is removed,
   and its parts.
3. API Contract: subsumption, pytest 9's `is_collection_argument_subsumed_by`.
   A node id (an argument with parts) is contained by another argument of the
   same invocation whose text differs when that argument has no parts and its
   path is `.`, equals the node id's path or is a directory above it (a prefix
   followed by `/`), or when both paths are equal and the other argument's
   parts prefix the node id's parts, its last part matching either exactly or
   as the name before `[` of a parametrized part.
4. API Contract: the finding. One error per contained node id, for the first
   argument in command order that contains it: code
   `SC-VERIFY-SUBSUMED-SELECTOR`, summary `Task <id> Verification passes the
   test selector <selector> together with <argument>, which contains it;
   pytest runs <argument> whole and ignores the selector, so a test that does
   not exist yet goes unnoticed`, location the Task file and the Verification
   line holding the command (`sectionLineContaining`), and fix `Name only the
   selectors in one command and run <argument> in a separate command,
   following the write-tasks rule "Test selectors and their containing
   path".` A Spec without a Task Graph records a skip with the manifest path,
   and `stagedDetectors` places the code at the Tasks stage. The detector
   executes nothing.
5. API Contract: the audited commit. `resolveQAAuditorEvidence` sets
   `AuditedCommit` to `qaAuditedHead(ctx, plan.WorkDir)`; when that fails the
   field stays empty and the evidence is still returned, as an auditor signal
   that cannot answer is.
6. API Contract: the refusal report. `WritePreconditionRefusalReport` writes
   `build: <strconv.Quote(commit)>` right after the `verdict:` line when the
   trimmed `AuditedCommit` is not empty, and no `build` key otherwise; every
   other byte is unchanged.
7. API Contract: the blocking mechanical report. `mechanicalQAReportContent`
   writes the same line right after `verdict:` when `result.Blocking` is true
   and the commit is not empty; the non-blocking seed, which the Agent turns
   into the full report with its own `build`, keeps its bytes.
8. API Contract: documentation. The `spec` guide's code table and the
   Roundfix `spec` reference describe `SC-VERIFY-SUBSUMED-SELECTOR`; the
   `write-tasks` skill gains the section "Test selectors and their containing
   path"; the `qa-gate` refusal template and the Context-driven development
   guide's refusal paragraph name `build`; the glossary adds **Subsumed Test
   Selector** and revises **Precondition Refusal** and **QA Report**.

### Surface Transcripts

None. No command, flag, exit code or stream placement changes: `spec check`
gains one finding code in its existing text and JSON shapes, and QA Reports
gain one frontmatter key; the tests of API Contracts 4, 6 and 7 assert the
exact text.

### Invariants

```text
1. The selector check executes no command and reads no test file.
2. A command with no pytest invocation, an invocation with --keep-duplicates, and a node id no other argument contains are never reported.
3. Completed, in-progress and QA Tasks are never reported.
4. The prober's verdicts and texts are unchanged.
5. A final report the Daemon writes names the Run Worktree's HEAD as build, quoted, once, right after verdict.
6. The pending seed and every report already written keep their bytes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariant 1 → the invocation, the collection arguments, subsumption and the finding (API Contracts 1-4)
- Goal 2, Success Metric 2, Invariants 2-4 → the same contracts' exclusions and stage gating (API Contracts 2-4)
- Goal 3, Core Feature 3, Success Metrics 3-4, Invariants 5-6 → the audited commit, the refusal report and the blocking report (API Contracts 5-7)
- Goal 4, Core Features 2 and 4 → documentation (API Contract 8)
- Success Metric 5 → the sabotages of the QA Task over the tests of API Contracts 1-7
- Success Metric 6 → the PRD's Unreachable Acceptance

## Absence Evidence

1. The subsumed selector check — command: `grep -rn -e SUBSUMED -e SubsumedSelector -e 'Subsumed Test Selector' internal docs/user-guide .agents skills CONTEXT.md`; result: no match at `42f523c2`.
2. A reader of test runner arguments — command: `rg -n 'pytest' internal/speccheck -g '!*_test.go'`; result: only `internal/speccheck/governed.go`, the `pytest.ini` Governed Path pattern.
3. A `build` key written by the Daemon — command: `grep -n -e '"build: ' -e 'build: %s' internal/spec/qa.go internal/daemon/task_engine.go`; result: no match at `42f523c2`.

## ADR Alternatives

- ADR-0315: no rejected alternative is taken; the prober gains no cause text, no other runner is recognized, and the finding is an error, not a gap.
- ADR-0148: its single prober and classification stay unchanged; the new check reads the written command and never runs it, so it is not the second classifier of executed commands the ADR rejects.
- ADR-0132: the refusal keeps its terminal row and its report rather than being suppressed; it only gains the audited commit beside the auditor fields.

## Integration Points

The Task files the check reads and the Run Worktree's Git `HEAD` the Daemon
reads. No test, Verification command or QA row runs pytest, reaches the
network, or reads another repository.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, and passes no real-time duration (the Real-Time
Budget rule gives a new file a ceiling of zero).

- `internal/speccheck/verification_selector_test.go` (new), through
  `CheckStage` with the fixture helpers `writeConsumerFixture` and
  `checkProofFixture`:
  - `TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath`: ten commands
    (`pytest <id> <file>`; `uv run pytest -q <file> <id>`; `uv run --group dev
    pytest tests <id>`; `python -m pytest tests/ <class id>`; `python3.12 -m
    pytest . <id>`; `pytest ./<file> <id>`; a class node id prefixing a method
    node id; `test_it` prefixing `test_it[new]`; the same invocation inside
    `out="$(…)"`; `cd tools && .venv/bin/pytest <id> tests`), each with exactly
    one error at both `StageAll` and `StageTasks`, the exact summary and fix,
    and the Verification line.
  - `TestSubsumedSelectorAcceptsSelectorsNoArgumentContains`: ten commands
    with no error (a lone node id; a sibling file; two sibling node ids;
    `tests/test_ab.py::test_new tests/test_a`, which shares only a string
    prefix; `--keep-duplicates`; `--deselect <id> <id>[new]`; `-k <file>
    <id>`; two invocations joined by `;`; `grep -qF -- <id> <file>`; the same
    node id twice).
  - `TestSubsumedSelectorJudgesOnlyPendingNonQATasks`: completed and
    in-progress Tasks yield no error; the PRD stage yields none and records
    the code as skipped.
- `internal/spec/qa_build_test.go` (new):
  `TestPreconditionRefusalReportRecordsItsAuditedCommit`, with the commit
  `1234567e89012345678901234567890123456789`, which reads as a YAML number
  unquoted: the report opens with `---`, `verdict: fail`, the quoted `build`
  and `auditing_binary:`, holds one `build` key, and `ReadQAReport` returns
  the refusal; a blank commit writes no `build` key.
- `internal/daemon/qa_report_build_test.go` (new):
  `TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal`: the refusal
  and the blocking result record `build` after `verdict`, the pending seed
  records none; through `writeMechanicalQAReport` on the task-cycle fixture,
  the written refusal names the fixture Run Worktree's `HEAD`.
- `internal/daemon/task_engine_test.go`: in
  `TestWriteMechanicalQAReportWritesThePreconditionRefusal`, the refusal
  frontmatter prefix and the blocking report of "a gate that refused nothing
  writes the report it wrote before" gain the `build` line of the fixture's
  `HEAD`, a declared break of pinned bytes; nothing else in the file changes.
- Unchanged and run by name: `TestVerificationPathUndeclaredRaisesTheMissingFile`,
  `TestWritePreconditionRefusalReportWritesOneTerminalRow`,
  `TestWritePreconditionRefusalReportIsReadableAsARefusal`,
  `TestPreconditionRefusalReportNamesItsAuditor`,
  `TestMechanicalQAReportSeedStartsPending`,
  `TestWriteMechanicalQAReportRecordsTheRefusal` and
  `TestMechanicalQAReportStartsUnsuffixedOnAFreshDate`.

## Build Order

1. The selector check, its wiring and its tests (no dependency).
2. The `spec` guide row, the Roundfix `spec` reference, the `write-tasks`
   section, both skills' versions, the Behavior Surface Record and **Subsumed
   Test Selector** (depends on: 1, which ships the code this text describes).
3. The audited commit in the refusal and blocking reports, and their tests
   (no dependency; it shares no file with 1 or 2).
4. The `qa-gate` refusal template, the Context-driven development guide, the
   `qa-gate` skill's version, the Behavior Surface Record, **Precondition
   Refusal** and **QA Report** (depends on: 2, which shares `CONTEXT.md`, the
   skill version record and the Behavior Surface Record, and 3, which ships
   the key this text describes).
5. QA gate (depends on: 1-4).

No Spec is listed in `requires`. Specs 0299, 0300 and 0301 deliver first and
raise the Roundfix skill's version; this Spec rebases on them and touches none
of their files apart from that skill's version fields.

## Risks & Considerations

- An adopting repository that runs pytest through a wrapper this rule does not
  know (`tox -e`, `nox`, a Makefile target) is not judged; the prober still
  refuses the vacuous command, as today.
- Lexical comparison misses two spellings of one path; it never reports a
  pair pytest would keep apart, because a containing path is matched only at a
  `/` boundary.
- An adopter on pytest 8 is told about an overlap that its pytest would have
  refused anyway; the selector adds nothing beside its own file, so the fix is
  still right.
- `build` written by the Daemon is the full commit, while Agents often write
  an abbreviated one; both name the audited commit, and nothing compares them.

## Glossary

- adds: **Subsumed Test Selector**
- changes: **Precondition Refusal**
- changes: **QA Report**

## Decisions

- The selector rule is structural and pytest-only; see ADR-0315.
- The audited commit travels in `spec.AuditorEvidence`, the evidence every
  report writer already receives, so no writer signature changes.
- Only reports the Daemon closes without an Agent turn record it; the seed
  keeps its bytes so the Agent's full report never carries two `build` keys.
