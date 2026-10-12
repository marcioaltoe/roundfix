---
spec: 0260-help-on-every-command-path-exits-zero
prd: _prd.md
created: 2026-10-12
---

# Help on every command path exits zero — Technical Spec

## Executive Summary

Two dispatchers in `internal/cli` read a help flag as an operand. The
`qa-report` dispatcher checks for help only when `accept` is absent, so
`accept --help` reaches the report-path check. The `baseline` dispatcher sends
`capabilities`, `skills` and `assets` to a usage error unless the next word is
one of their commands. Both learn to answer a help request, a new test asks
every command path for help, and the Behavior Surface Record's contract stops
falling back to a parent's help. The trade-off accepted: a group prefix prints
the whole `baseline` usage instead of new text of its own, so no usage text is
written and no fingerprint changes, at the cost of a longer answer.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. Help prints constants; nothing
  reaches a network or credential. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0256: a Behavior Surface is "the
  help of one command path", and its record is computed "from the built
  command help"; API Contract 4 makes the contract read exactly that. Its
  authoring rule needs a skills Task, which task_01 is. ADR-0184: "Task tests
  and QA Requirements name the transcript by number and reuse its text", which
  task_01 and task_02 do. ADR-0259: "Every top-level test in a Parallel Test
  Package calls `t.Parallel()` as its first statement", which both new tests
  do. ADR-0189 and ADR-0233 govern the owned skill's version raise. ADR-0250
  does not apply: no Baseline module version or coverage record changes.
  ADR-0244 is met by the `None.` declarations. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths, measured with
  `speccheck.GovernedPath` through a `go test -overlay` probe that wrote
  nothing, are granted by the cycle AE answer of 2026-10-11. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0260-help-on-every-command-path-exits-zero/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## System Architecture

No new package or seam. The change touches three existing places:

- `internal/cli/qa_report.go` — `runQAReportCommand` gains a help check on the
  arguments after `accept`, using the existing `commandWantsHelp`, before the
  path-count and report-path checks.
- `internal/cli/baseline_profile.go` — `runBaselineCommand` answers a help
  token in the position after `capabilities`, `skills` or `assets` with
  `commandUsage("baseline")`, before the existing usage errors.
- `internal/docscontract/skill_coverage_test.go` — `skillCoverageCLIHelp`
  reads the command path's own help and fails the contract when it exits
  non-zero, instead of walking up to the nearest ancestor that succeeds.

The help sweep is a new test file in `internal/cli`. The `roundfix` skill's
settle and baseline references describe the new behavior.

Current behavior, characterized at `7a8b7f70` with `make build` and a sweep of
`--help` and `-h` over 70 command paths and group prefixes: eight calls fail,
listed in `_prd.md`; 132 exit 0 with a usage block and an empty standard
error. A prototype of this design on a scratch copy passed the sweep, kept
`docs/references/behavior-surfaces.json` byte-identical, and made the
Behavior Surface Record contract fail naming `qa-report accept` when only the
production fix was reverted.

## Implementation Design

### Interfaces

No exported symbol changes. The test-only helper in the new sweep file:

```go
// helpCommandPaths returns every command path the root usage names and every
// group prefix of one, in first-seen order without duplicates.
func helpCommandPaths(rootUsage string) []string
```

The contract helper keeps its signature and loses its loop:

```go
func skillCoverageCLIHelp(t *testing.T, words []string) []byte
```

Invariants:

1. A help flag never reads, writes or validates a report, a repository or a
   journal.
2. A help answer writes the usage to standard output, nothing to standard
   error, and exits 0.
3. A non-help argument keeps today's behavior, streams and exit code on every
   touched path.
4. The Behavior Surface Record's bytes do not change.

### Data Models

None. No stored or serialized shape changes.

### API Contracts

1. API Contract 1: `roundfix qa-report accept` with `--help` or `-h` anywhere
   in its arguments before a `--` terminator, as `commandWantsHelp` decides,
   prints `qaReportUsage` and exits 0. `qa-report accept` with no path, with
   two paths, or with a blank path keeps the usage error and exit 2; a path
   keeps today's eligibility decision and exit 0 or 1.
2. API Contract 2: `roundfix baseline capabilities`, `baseline skills` and
   `baseline assets` followed directly by `--help`, `-h` or `help` print
   `commandUsage("baseline")` and exit 0. `baseline capabilities check`,
   `baseline skills restore`, `baseline skills reconcile` and
   `baseline assets sync` keep their own help. A bare group or an unknown
   command after it keeps today's usage error, text and exit 2.
3. API Contract 3: `TestEveryCommandPathPrintsItsHelp` derives the command
   paths from `roundfix --help`, as the Behavior Surface Record does (words up
   to the first token starting with `-`, `[`, `(` or `<`, with a `<a|b>` token
   after a command word expanded into one path per alternative), adds every
   group prefix, requires at least 60 paths, and for each path and each of
   `--help` and `-h` requires exit 0, standard output starting with
   `Usage:` and an empty standard error. A failure names the command line,
   exit code, first standard-output line and standard error.
4. API Contract 4: `skillCoverageCLIHelp` calls the command path's own
   `--help` once and, when that call does not exit 0, fails the test with
   `roundfix <path> --help: help exited <code>: <stderr>`.
   `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` asserts
   that `command qa-report accept` carries the fingerprint of
   `qa-report accept --help`'s own standard output and that the call exits 0,
   failing with `command qa-report accept: fingerprint does not use its own
   help`.

### Surface Transcripts

1. Surface Transcript: `qa-report accept --help` prints its usage.

   ```transcript
   $ roundfix qa-report accept --help
   stdout:
   Usage:
     roundfix qa-report accept <path>

   Applies the declared-acceptance policy to the selected QA Report. The report's
   Spec directory supplies its declarations. Produces no stdout; exit status is
   the machine-readable result.

   Exit codes:
     0  the report is acceptable
     1  the report is missing, unreadable, or unacceptable
     2  usage error
   stderr:
   exit: 0
   ```

2. Surface Transcript: `qa-report accept -h` prints the same usage.

   ```transcript
   $ roundfix qa-report accept -h
   stdout:
   Usage:
     roundfix qa-report accept <path>

   Applies the declared-acceptance policy to the selected QA Report. The report's
   Spec directory supplies its declarations. Produces no stdout; exit status is
   the machine-readable result.

   Exit codes:
     0  the report is acceptable
     1  the report is missing, unreadable, or unacceptable
     2  usage error
   stderr:
   exit: 0
   ```

3. Surface Transcript: a help flag after the report path still prints help
   and reads no report.

   ```transcript
   $ roundfix qa-report accept docs/specs/absent/qa/qa-report-2026-10-12.md --help
   stdout:
   Usage:
     roundfix qa-report accept <path>

   Applies the declared-acceptance policy to the selected QA Report. The report's
   Spec directory supplies its declarations. Produces no stdout; exit status is
   the machine-readable result.

   Exit codes:
     0  the report is acceptable
     1  the report is missing, unreadable, or unacceptable
     2  usage error
   stderr:
   exit: 0
   ```

4. Surface Transcript: `baseline capabilities --help` prints the `baseline`
   usage.

   ```transcript
   $ roundfix baseline capabilities --help
   stdout:
   Usage:
     roundfix baseline [--repo <path>] [--format <text|json>]
   ...
   stderr:
   exit: 0
   ```

5. Surface Transcript: `baseline skills -h` prints the `baseline` usage.

   ```transcript
   $ roundfix baseline skills -h
   stdout:
   Usage:
     roundfix baseline [--repo <path>] [--format <text|json>]
   ...
   stderr:
   exit: 0
   ```

6. Surface Transcript: `baseline assets --help` prints the `baseline` usage.

   ```transcript
   $ roundfix baseline assets --help
   stdout:
   Usage:
     roundfix baseline [--repo <path>] [--format <text|json>]
   ...
   stderr:
   exit: 0
   ```

Transcripts 4 to 6 elide the rest of the usage with `...`; the implementing
test asserts the whole standard output equals `commandUsage("baseline")`.

## Coverage Map

- Goal 1, User Stories 1-2, Core Features 1-2 → `runQAReportCommand`,
  `runBaselineCommand`; Surface Transcripts 1-6; API Contracts 1-2.
- Goal 2, User Story 3, Core Feature 3, Success Metrics 1-2 →
  `TestEveryCommandPathPrintsItsHelp`, the four new rows of
  `TestRunCommandHelp`; API Contract 3.
- Goal 3, Core Feature 4, Success Metric 3 → `skillCoverageCLIHelp`,
  `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent`; API
  Contract 4.
- Goal 4, Core Feature 5 → the `roundfix` skill's settle and baseline
  references.

## Integration Points

None external. The Behavior Surface Record contract (`make verify-docs`) and
the release plan's `skill-coverage` check read the record, whose bytes stay
the same (Invariant 4).

## Testing Approach

- `TestHelpSurfaceTranscripts`, new in `internal/cli/command_help_test.go`,
  asserts Surface Transcripts 1-6 through `runCLI`: every standard-output line
  in order, the empty standard error and exit 0; for 4-6 it compares standard
  output with `commandUsage("baseline")`. Transcript 3's path does not exist,
  so a regression that reads it fails with exit 1.
- `TestEveryCommandPathPrintsItsHelp`, new in the same file, is the class
  guard of API Contract 3. Before the fix it fails on exactly the eight calls
  in `_prd.md` (prototype measurement).
- `TestRunCommandHelp` in `internal/cli/cli_test.go`, the help-exit test the
  adopted entry names, gains rows for `qa-report accept`,
  `baseline capabilities`, `baseline skills` and `baseline assets`, each
  requiring exit 0, an empty standard error and its usage's first command
  line.
- The two contract tests in `internal/docscontract` change as API Contract 4
  says; they run under `-tags docscontract`.
- Existing `qa_report_test.go` and `baseline_profile_test.go` cases stay
  unchanged and keep passing, which proves Invariant 3.

## Build Order

1. Help on the four command paths, the sweep and transcript tests, the
   help-exit rows, the contract helper without fallback, and the skill
   references with the recorded version raise (task_01).
2. Final QA gate (task_02) (depends on: 1).

## Risks & Considerations

- The sweep runs every command with a help flag in-process. A future command
  that ignores help would run its real behavior inside `runCLI`'s test
  environment; the test's failure message names it, and the environment is
  the per-test one every `internal/cli` test uses.
- The skill version raise is written only by the record command; the Task
  declares its record file.
- The bare-group usage error keeps a doubled `baseline` and a pointer at
  `baseline profile --help`. It is outside this Spec and recorded in
  `docs/backlog/2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md`.

## Glossary

None.

## Decisions

- A group prefix prints the existing `baseline` usage rather than new text, so
  the change writes no usage text and keeps every fingerprint.
- The class guard lives in `internal/cli`, which `make test` runs on every
  change, and the docscontract helper fails instead of masking a broken help.
- No ADR: no decision here is hard to reverse or departs from ADR-0256.
- Research: the Secondbrain and Exa consultations are recorded in
  `_prd.md` → Research record; neither changed the design.
