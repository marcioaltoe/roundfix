---
spec: 0214-measure-before-changing
prd: _prd.md
created: 2026-10-01
---

# Measure before changing — Technical Spec

## Executive Summary

A new package, `internal/runcause`, reads a repository's terminal Spec Runs
from the Run Database through the reader that never migrates, finds each
failed Verification attempt and each corrective Task, and classifies its
cause with an embedded, ordered signature table. `roundfix runs causes`
prints that report and, as JSON, the per-Task attempt counts. A flag-gated
test harness in `internal/judge` reads those counts as labels, asks Jev 1.13
one fixed question about each Task's acceptance criteria through the judge
that Spec 0205 delivers, and writes a record a test can recompute. Two
documentation Tasks run the measurements and write three reference documents
with verdicts computed by rules fixed here. The trade-off this design accepts
is keyword classification: it can misclassify an item for the wrong reason,
and in exchange every class is reproducible, free, and shows the signature
that decided it (ADR-0214).

## Project Constraints

- Identifier strategy: not applicable — no persistent identifier is minted.
  Classes, signatures and triggers are fixed words of the embedded table;
  records are named by their schema. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — only the Task acceptance harness
  sends HTTPS requests, only when `-measure-task-acceptance` is passed, and
  only through the judge's transports: OpenRouter's System One API with
  `ROUNDFIX_OPENROUTER_API_KEY` requesting `jev-1.13`, or else TypeSafe with
  `ROUNDFIX_TYPESAFE_API_KEY` requesting `jev-1.13.0`. Keys come from the
  process environment of that explicit run, travel only in their endpoint's
  authorization header, and appear in no output, record or log; the generic
  `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY` are never read. Every unit test
  injects an `http.RoundTripper`. `runs causes` makes no network request.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0214 and ADR-0215 (this Spec)
  decide the classification and the measurement boundary. ADR-0201 bounds
  what the judge may send and already admits Task files, ADR-0201: "Findings,
  Backlog Entries and Task files are within the authorization but not read";
  the harness keeps its transports, keys, Judge Log and ceiling. ADR-0200
  pins Jev 1.13, ADR-0200: "Moving the pin means measuring again and
  recording the new thresholds", and the harness compares only answers whose
  model normalizes to that version. ADR-0004, ADR-0008, ADR-0014, ADR-0033
  and ADR-0038 define the Run Database, its raw payloads, the Daemon's
  Verification, journal retention and the single repair, which together make
  the attempt count a label; nothing here changes them. ADR-0035 makes the
  Spec Root configuration, ADR-0035: "the Spec Root becomes explicit
  configuration", and the command reads Task Graphs from it and from its
  archive. ADR-0120 and ADR-0121 own the history root and archive moves,
  which the archive measurement only reads. ADR-0089 passes the environment
  to code under test explicitly. Under ADR-0184, ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts", so the command is stated as
  transcripts below. The command's description goes in the `runs` files,
  because ADR-0187 split the skill and the command reference by command
  family, ADR-0187: "an entry file plus one file per command family".
  ADR-0189 ties an owned skill's version to its content, ADR-0189: "An owned
  skill's version names its content", so the Roundfix Skill's version rises;
  ADR-0081 and ADR-0149 sanction `make skills-sync` and `make
  baseline-digests`. ADR-0193 names
  Spec 0205 as a prerequisite. ADR-0182, ADR-0178, ADR-0166 and ADR-0130
  govern each Task commit; ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104,
  ADR-0117, ADR-0155, ADR-0156 and ADR-0167 bind the gate; ADR-0093 and
  ADR-0094 check consistency. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and the 2026-10-01 program decision, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runs.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Owner | File |
| --- | --- | --- |
| Events of chosen kinds for one Run | `Store.RunEventsOfKinds` | `internal/store/run_events_of_kinds.go` (new) |
| Signature table, classes, triggers | `signatures.json`, `Load` | `internal/runcause/signatures.json` (new), `internal/runcause/signatures.go` (new) |
| Classifying one item | `Table.Classify`, `Table.Trigger` | `internal/runcause/classify.go` (new) |
| Runs, attempts, corrective Tasks, summary | `Build`, `Report` | `internal/runcause/report.go` (new) |
| The command | `runRunsCausesCommand` | `internal/cli/runs_causes.go` (new), `internal/cli/runs.go`, `internal/cli/cli.go` |
| Record consistency | `TestCausesRecordIsConsistent` | `internal/runcause/record_test.go` (new) |
| Task acceptance harness and statistics | test-only | `internal/judge/task_acceptance_measure_test.go` (new), `internal/judge/task_acceptance_stats_test.go` (new), `internal/judge/testdata/task-acceptance-question.json` (new) |
| Measurements and recommendations | documents | `docs/references/archived-evidence-measurement.md`, `docs/references/corrective-causes-measurement.md`, `docs/references/corrective-causes.json`, `docs/references/task-acceptance-measurement.md`, `docs/references/task-acceptance-remeasurement.json` (all new) |

`internal/runcause` imports `internal/store`, `internal/runevent` and
`internal/spec`, and nothing imports it except `internal/cli`. The harness
lives only in `_test.go` files of `internal/judge`, so the binary gains no
network path (ADR-0215).

## Implementation Design

### The signature table

`internal/runcause/signatures.json` is embedded and is byte-identical to this
block. Patterns are RE2. `source` names the one evidence text a signature
reads: `diagnostic` (the tail of a failed attempt's Verification log),
`command` (the failing Verification command), `task` (a corrective Task's
title and `## Overview`), or `any` (each text the item has, in that order).

```json
{
  "schema": "roundfix/cause-signatures/v1",
  "classes": ["scope-or-authorization", "shared-section-contract", "repository-convention", "implementation-defect", "environment"],
  "repository_knowledge": ["scope-or-authorization", "shared-section-contract", "repository-convention"],
  "diagnostic_tail_bytes": 65536,
  "task_text_bytes": 4096,
  "triggers": [
    {"id": "pre-pr-review", "pattern": "(?i)pre-PR review|independent review|review round|second review"},
    {"id": "qa-gate", "pattern": "(?i)QA finding|QA gate|final QA|\\bF-\\d{3}\\b"},
    {"id": "verification", "pattern": "(?i)repository Verification|make verify|Verification fail"}
  ],
  "signatures": [
    {"id": "database-locked", "class": "environment", "source": "diagnostic", "pattern": "(?i)database is locked|SQLITE_BUSY"},
    {"id": "disk-full", "class": "environment", "source": "diagnostic", "pattern": "(?i)no space left on device"},
    {"id": "network-unreachable", "class": "environment", "source": "diagnostic", "pattern": "(?i)connection refused|no such host|i/o timeout|TLS handshake timeout|network is unreachable"},
    {"id": "killed-or-timed-out", "class": "environment", "source": "diagnostic", "pattern": "(?i)signal: killed|panic: test timed out|context deadline exceeded"},
    {"id": "governed-path", "class": "scope-or-authorization", "source": "diagnostic", "pattern": "(?i)governed path|tooling.authority|not authori[sz]ed|outside (the|its) grant|undeclared path"},
    {"id": "authorization-or-scope", "class": "scope-or-authorization", "source": "task", "pattern": "(?i)authori[sz]ation|governed path|\\bgrant\\b|undeclared path|outside (the|its) (scope|slice|grant)"},
    {"id": "shared-section", "class": "shared-section-contract", "source": "any", "pattern": "(?i)QA settlement|byte.identical|skills-sync-check|drifts from \\.agents/skills|shared section"},
    {"id": "record-or-golden", "class": "repository-convention", "source": "any", "pattern": "(?i)coverage (record|addition|equivalence)|TestCoverageEquivalence|golden|re-?record|owned-skill-versions|skill version|baseline-digests|digest pin|characteri[sz]ation|corpus"},
    {"id": "lint-or-format", "class": "repository-convention", "source": "diagnostic", "pattern": "(?i)gofmt|go vet|staticcheck|golangci|\\blint\\b"},
    {"id": "help-or-guide", "class": "repository-convention", "source": "task", "pattern": "(?i)help text|user guide|command reference|Roundfix skill|describe the|teach the"},
    {"id": "go-test-failure", "class": "implementation-defect", "source": "diagnostic", "pattern": "(?m)^\\s*--- FAIL: |\\[build failed\\]|^panic: |undefined: |cannot use "},
    {"id": "review-defect", "class": "implementation-defect", "source": "task", "pattern": "(?i)found (a|an|one|two|three|four|\\d+) (defects?|gaps?|ways|bugs?)|confirmed in source|reproduced before|blocker|regression"}
  ]
}
```

A change to this file changes the classification contract, and every report
names the SHA-256 of the bytes that produced it. `Load` refuses a class
outside `classes`, a duplicate `id`, an unknown `source` and a pattern that
does not compile.

### Interfaces

```go
// internal/store/run_events_of_kinds.go
// RunEventsOfKinds returns one Run's events whose kind is in kinds, in
// cursor order, with payloads. It filters in SQL and writes nothing.
func (store *Store) RunEventsOfKinds(ctx context.Context, runID string, kinds ...runevent.Kind) ([]JournalEvent, error)
```

```go
// internal/runcause
func Load() (Table, error) // the embedded table; Table.SHA256 is its digest

type Evidence struct{ Diagnostic, Command, Task string } // empty = absent
// Classify returns the first signature whose source text matches, or
// ("unclassified", "", "") when none does.
func (t Table) Classify(e Evidence) (class, signature, source string)
func (t Table) Trigger(taskText string) string // trigger id, or "unknown"

type Reader interface {
	ListRuns(ctx context.Context, query store.ListRunsQuery) ([]store.Run, error)
	RunEventsOfKinds(ctx context.Context, runID string, kinds ...runevent.Kind) ([]store.JournalEvent, error)
}
type Request struct {
	Reader                  Reader
	RepositoryRoot          string // runs.repository_root to match
	SpecsRoot, ArchiveRoot  string // where Task Graphs are read
	Since, Until            time.Time // UTC midnights; zero = open end
}
func Build(ctx context.Context, t Table, req Request) (Report, error)
```

### Building the report

1. Runs: `ListRuns` with `RepositoryRoot` and `StatesTerminal`, kept when
   `spec_slug` is not empty and `Since <= created_at < Until`; oldest first.
2. Events: `RunEventsOfKinds(run, daemon.verification, daemon.task)`.
3. Attempts: an attempt is (Run, Task, `attempt`). It failed when a
   `daemon.verification` event of that attempt has `phase` `failed` or
   `verdict` `failed`. Its check is the `command` of the first `failed`
   event (or the joined `commands`), and its diagnostic is the last
   `diagnostic_tail_bytes` of the file named by that event's
   `diagnostic_path`, read only when it is a regular file (`Lstat`, never
   through a symbolic link) inside the Run's artifact directory; otherwise
   the diagnostic is absent.
4. Tasks: for every Task with a `daemon.task` `started` event in a kept
   Run, the counts of `verdict` `passed` and `failed` events, of
   `verification_feedback` phases, and of Runs. A Task's Spec is read from
   `<SpecsRoot>/<slug>/_tasks.md` or else `<ArchiveRoot>/<slug>/_tasks.md`:
   `qa` is whether its id equals the graph's `qa:` Task, and `corrective` is
   whether it is a graph node, not the QA Task, whose `task_NN` number is
   greater than the QA Task's. When neither file exists, `corrective` is
   null and the slug is listed in `specs_not_found`.
5. Corrective items: one per corrective Task, attached to the first kept Run
   in which it started. Its text is its title and `## Overview`, at most
   `task_text_bytes`; its check is `Trigger(text)`; its class is
   `Classify(Evidence{Task: text})`.
6. Items are ordered by Run, then failed attempts in event order, then
   corrective Tasks in Task order. The summary counts each class,
   `unclassified`, the classified items and the repository-knowledge items.

### Data Models

The JSON report, schema `roundfix/runs-causes/v1`:

```json
{"schema": "roundfix/runs-causes/v1",
 "window": {"since": "2026-09-17", "until": null},
 "signatures_sha256": "<64 hex>", "runs": 1,
 "items": [{"kind": "verification", "spec": "0300-example", "task": "task_01",
   "run_id": "<id>", "attempt": 1, "check": "go test ./internal/store",
   "class": "environment", "signature": "database-locked", "evidence": "diagnostic"}],
 "tasks": [{"spec": "0300-example", "task": "task_01", "qa": false, "corrective": false,
   "verdicts_passed": 1, "verdicts_failed": 1, "feedback_rounds": 1, "runs": 1}],
 "summary": {"scope-or-authorization": 0, "shared-section-contract": 0,
   "repository-convention": 0, "implementation-defect": 0, "environment": 1,
   "unclassified": 0, "items": 1, "classified": 1, "repository_knowledge": 0},
 "specs_not_found": []}
```

`attempt` is null and `kind` is `corrective` for a corrective item;
`signature` and `evidence` are null for `unclassified`. No absolute path,
key or diagnostic text is in the report.

The Task acceptance record, schema `roundfix/task-acceptance-remeasurement/v1`,
carries `status` (`measured` or `blocked`), `reason`, `labels_sha256` (of
the cause record's bytes), `question_sha256`, `pinned_model`, `transport`,
`calls`, `input_tokens`, `cost_usd`, `tasks` (each `spec`, `task`,
`repaired`, `verdicts_failed`, `criteria_chars`, `state_sha256`, `noul`,
`model`, `outcome` and `reason`), `excluded` (each `spec`, `task`,
`reason`), `statistics` and `verdict`.

### The Task acceptance question

`internal/judge/testdata/task-acceptance-question.json`, byte-identical to:

```json
{
  "schema": "roundfix/task-acceptance-question/v1",
  "judgment": "task-acceptance",
  "question_id": "acceptance_observable",
  "question": {
    "type": "noul",
    "instructions": "`acceptance_criteria` is the acceptance-criteria list of one implementation Task from a software specification, and `task_title` is the Task's title. Could a reviewer who did not write the Task decide every criterion as met or not met by running a command or reading a named file or output, without judging intent or quality?",
    "criteria": {
      "true": "Every criterion names an observable result: a command outcome, an exit code, a file, a field, an output line, or a test that passes or fails",
      "false": "At least one criterion depends on judgment, intent or quality, or names no result an observer could check"
    }
  },
  "task_title_max_chars": 200,
  "acceptance_criteria_max_chars": 4000,
  "flag_when_noul_below": 0.5,
  "bootstrap": {"resamples": 2000, "seed": 214, "cluster": "spec", "max_discarded_share": 0.05},
  "random_baseline_seed": 214
}
```

The state is `{"task_title", "acceptance_criteria"}`: the Task file's first
`# ` heading and its `## Acceptance Criteria` section, each truncated to its
limit. Nothing after that section's end is read, so a `## Result` written by
an Agent can never reach a request.

### The harness

1. Labels: every `tasks` entry of the cause record with `qa` false and at
   least one verdict; `repaired` is `verdicts_failed >= 1`. A Task whose file
   is missing in the active or archived Spec, that has no `## Acceptance
   Criteria`, or whose state fails the judge's language gate is `excluded`
   with the reason.
2. Asking: through Spec 0205's transport selection, request, model pin,
   Judge Log and ceiling, with judgment `task-acceptance`. A missing key, a
   stop or the ceiling end the run; with no key the record is `blocked` and no
   request is sent. An answer whose model does not normalize to Jev 1.13 is
   `skipped` and enters no statistic.
3. Statistics, over answered Tasks, score = 1 − `noul`: AUROC by the
   Mann-Whitney count with ties at one half; the 95% interval from
   `resamples` cluster-bootstrap draws of whole Specs with a PCG generator
   seeded `seed`, nearest-rank 2.5th and 97.5th percentiles, draws without
   both classes discarded; the AUROC of a seeded random score and of
   `criteria_chars`; and at `noul < 0.5` the flagged count, precision and
   recall.

### Decision rules

Fixed here, before any data are read:

1. **Memory question** (cause record over its window): `inconclusive` when
   items are fewer than 30 or `unclassified` exceeds a third of them;
   otherwise `reopen` when repository-knowledge items exceed
   `implementation-defect` plus `environment` items; otherwise
   `keep closed`.
2. **Task acceptance**: `inconclusive` when fewer than 150 Tasks were
   answered, fewer than 20 of them were repaired, the record is `blocked`, or
   more than `max_discarded_share` of draws were discarded; otherwise `adopt`
   when the interval's lower bound is at least 0.70, at least 10 Tasks are
   flagged and precision is at least twice the base rate; otherwise
   `do not adopt`.
3. **Archive candidates**: a kind of archived content (QA evidence, design
   and reference binaries) is a candidate when no repository file outside
   the history root reads it and the ablation without it leaves `make verify`
   and `make verify-docs` at exit `0`. Agent reads and the mirror are
   reported as costs, never as vetoes.

Each fired rule writes one open Backlog Entry: `reopen` →
`docs/backlog/<date>-reopen-the-memory-question-for-agent-sessions.md`;
`adopt` → `docs/backlog/<date>-adopt-the-task-acceptance-judgment.md`; a
non-empty candidate set →
`docs/backlog/<date>-remove-archived-evidence-nobody-reads.md`, whose Target
states that removal waits for the maintainer's explicit approval.

### The archive measurement

In the Task's worktree, at its starting commit `C`, read-only:

1. **Inventory.** From `git ls-tree -r -l C -- docs/history`: files and
   bytes in total, of QA evidence (paths under
   `docs/history/specs/*/qa/evidence/`), of QA Reports, of core Spec
   artifacts, of other families, and of binaries (`.png`, `.jpg`, `.jpeg`,
   `.gif`, `.pdf`, `.db`, `.zip`).
2. **Static readers.** Every tracked file outside `docs/history/` that
   names `docs/history` or builds the archive root (`ArchiveDir`,
   `ArchiveSpecRoot`), with what it reads.
3. **Ablation.** A disposable clone outside the repository
   (`git clone --no-local` of `C` into a temporary directory), with the QA
   evidence directories and the binaries removed there only; `make verify`
   and `make verify-docs` run in it with their exit codes and failing tests
   recorded, and `git grep` over the clone timed five times before and after
   the removal.
4. **Agent reads.** From the Run Database opened as
   `file:<home>/.roundfix/roundfix.db?mode=ro&immutable=1`: `agent.tool_started`
   events whose payload names `docs/history/`, counted by kind of path and
   by tool kind, with distinct Runs.
5. **Mirror.** Files and bytes of the Secondbrain mirror's copy of the
   history root, read-only.

### Surface Transcripts

The fixture is a repository whose Spec `0300-example` has QA Task `task_03`
and corrective `task_04`, whose overview reads "Pre-PR review found that the
golden output was not re-recorded". One terminal Run
`run_20261001T120000Z_0000000000000001`, created 2026-10-01, failed
`task_01` attempt 1 with "database is locked" in its log, failed `task_02`
attempt 1 with `--- FAIL: TestCausesFixture`, and failed `task_02` attempt 2
on `make verify` with a log no signature matches.

1. Surface Transcript: a classified window.

   ```transcript
   $ roundfix runs causes --since 2026-10-01
   stdout:
   verification 0300-example/task_01 run_20261001T120000Z_0000000000000001 attempt-1 environment database-locked: go test ./internal/store
   verification 0300-example/task_02 run_20261001T120000Z_0000000000000001 attempt-1 implementation-defect go-test-failure: go test ./internal/cli
   verification 0300-example/task_02 run_20261001T120000Z_0000000000000001 attempt-2 unclassified -: make verify
   corrective 0300-example/task_04 run_20261001T120000Z_0000000000000001 - repository-convention record-or-golden: pre-pr-review
   Causes: 4 item(s) from 1 terminal Spec Run(s) in window 2026-10-01..*; scope-or-authorization 0, shared-section-contract 0, repository-convention 1, implementation-defect 1, environment 1, unclassified 1; repository knowledge 1 of 3 classified; signatures <digest>
   stderr:
   exit: 0
   ```

2. Surface Transcript: no terminal Spec Run in the window, or no Run
   Database.

   ```transcript
   $ roundfix runs causes --until 2026-09-01
   stdout:
   Causes: 0 item(s) from 0 terminal Spec Run(s) in window *..2026-09-01; scope-or-authorization 0, shared-section-contract 0, repository-convention 0, implementation-defect 0, environment 0, unclassified 0; repository knowledge 0 of 0 classified; signatures <digest>
   stderr:
   exit: 0
   ```

3. Surface Transcript: a malformed date.

   ```transcript
   $ roundfix runs causes --since 2026-13-01
   stdout:
   stderr:
   roundfix: runs causes failed: --since must be a date in YYYY-MM-DD form, got "2026-13-01"
   Run 'roundfix runs causes --help' for usage.
   exit: 2
   ```

4. Surface Transcript: outside a Git repository.

   ```transcript
   $ roundfix runs causes
   stdout:
   stderr:
   roundfix: runs causes failed: runs causes requires a Git repository
   Run 'roundfix runs causes --help' for usage.
   exit: 2
   ```

A check is the first line of the command with whitespace collapsed, cut at
80 characters with `…`; `<digest>` is the first 12 hex digits of
`signatures_sha256`.

### API Contracts

1. API Contract: `roundfix runs causes [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]`
   — Surface Transcripts 1 to 4. `--since` is inclusive and `--until`
   exclusive, both UTC. Exit `0` whenever it ran, including an empty window
   and an absent Run Database; exit `2` for an unknown flag, a malformed date,
   `--until` not after `--since`, an unknown format, or no Git repository;
   exit `1` when the Run Database cannot be read. It writes nothing.
2. API Contract: `--format json` prints the document of Data Models with
   every field, including `tasks`.
3. API Contract: help — `roundfix runs --help` and the top-level help name
   the synopsis of API Contract 1, and the runs help lists `causes` with its
   options.
4. API Contract: the Task acceptance record of Data Models, written only by
   the harness, and the cause record, which is the command's JSON output
   saved unchanged.

## Vocabulary Contract

No emitted vocabulary is added to a documented catalog. The classes, triggers
and signature identifiers are defined by the embedded table above and
documented in `docs/user-guide/commands/runs.md` by task_01.

## Coverage Map

- Goal 1 → Building the report; The signature table; API Contract 1.
- Goal 2 → Building the report; Decision rules; Data Models.
- Goal 3 → The Task acceptance question; The harness; Decision rules.
- Goal 4 → The archive measurement; Decision rules.
- Goal 5 → Data Models; Testing Approach 5; Decision rules.
- User Story 1 → Building the report; Surface Transcript 1.
- User Story 2 → The signature table; API Contract 2.
- User Story 3 → The harness; The Task acceptance question.
- User Story 4 → The archive measurement.
- User Story 5 → The archive measurement; Testing Approach 6.
- Core Feature 1 → API Contracts 1 and 2; Surface Transcripts 1 to 4.
- Core Feature 2 → The signature table.
- Core Feature 3 → The signature table; Interfaces.
- Core Feature 4 → Building the report.
- Core Feature 5 → Building the report; Data Models.
- Core Feature 6 → The harness; The Task acceptance question.
- Core Feature 7 → The harness; Testing Approach 4.
- Core Feature 8 → Decision rules.
- Core Feature 9 → The archive measurement.
- Core Feature 10 → The archive measurement; Decision rules.
- Core Feature 11 → Decision rules.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 5.
- Success Metric 3 → Testing Approaches 3, 4 and 5.
- Success Metric 4 → Testing Approach 6.
- Success Metric 5 → Testing Approach 5.

## Integration Points

- **The Run Database.** Read through `store.OpenReader`, which opens it
  read-only and never migrates; an absent database is an empty report. The
  archive measurement's agent-read count opens it immutable through
  `sqlite3`. Nothing writes to it.
- **Run artifact directories.** Failed-attempt logs are read only inside the
  Run's own artifact directory, and only their tail.
- **OpenRouter and TypeSafe.** Reached only by the harness, only with its
  flag, through the judge's transports (Spec 0205).
- **The Secondbrain mirror.** Read-only file counts for the archive
  measurement; nothing is written there.

### Measured outside evidence

On 2026-10-01 the authoring session read the maintainer's Run Database
read-only: 210 Runs since 2026-08-27; 2,618 `daemon.verification` events
since 2026-09-17; 85 failed verdicts, 42 Verification Feedback rounds and
357 Task executions with a verdict across 63 Specs, 34 of them repaired; all
86 failed-attempt logs present; no `classification` value recorded on any
Verification event. 827 `agent.tool_started` events named `docs/history/`,
372 of them a QA evidence path. In the archived Task Graphs, every Task
numbered after its QA Task reads as a corrective one, and the overviews name
their trigger in words such as "Pre-PR review found", "QA finding F-001"
and "Corrective Task from the QA gate", which the trigger table encodes.

## Testing Approach

1. **The table.** Task 01's Verification requires the embedded file to be
   byte-identical to the block above; `TestLoadRefusesABrokenTable` covers
   each refusal of `Load`.
2. **Classification and building.** `internal/runcause/classify_test.go`
   gives one case per signature and per trigger, the first-match order, and
   an unmatched item. `report_test.go` writes a fixture Run Database through
   `store.Open` in a temporary home and requires the items, tasks and
   summary of Surface Transcript 1, a symbolic-link log treated as absent, a
   Spec missing from both roots, and that the Run Database file and the home
   are byte-identical afterwards.
3. **The command.** `internal/cli/runs_causes_test.go` reproduces Surface
   Transcripts 1 to 4 and API Contracts 2 and 3 with a temporary repository
   and home.
4. **The harness.** `internal/judge/task_acceptance_measure_test.go`, with a
   fake transport: requests carry only the two state fields, a sentinel in a
   Task's `## Result` and in a Go file never appears in a body, each key
   only in its header, only the generic keys set gives a `blocked` record
   with no request, a non-pinned model is skipped, and the Judge Log gains
   one line per request. `task_acceptance_stats_test.go` checks AUROC,
   ties, the cluster bootstrap's determinism and each verdict branch.
5. **Records.** `TestCausesRecordIsConsistent` and
   `TestTaskAcceptanceRecordIsConsistent` take a record and its document
   through flags, recompute every count, figure and verdict, and compare;
   without flags they run on fixtures, one of them sabotaged to fail.
6. **No archived file removed.** Task 03 and the QA gate require
   `git diff --name-only --diff-filter=DR <merge-base> -- docs/history` to be
   empty, and the archive document's counts to equal Git's at its commit.

Each new gate is proved to fail: the Task that adds it records the sabotage
and the failing test in its Result, then restores the code.

## Build Order

1. `RunEventsOfKinds`, the table, classification, the report, the command,
   its help, the command reference and the Roundfix Skill, task_01 (depends
   on: none).
2. The Task acceptance harness, question file and statistics, task_02
   (depends on: none). It reads the cause record's schema from this
   TechSpec, and it compiles against Spec 0205's `internal/judge`.
3. The archive measurement and proposal, task_03 (depends on: none).
4. The cause and Task acceptance measurements, their documents and
   follow-up Backlog Entries, task_04 (depends on: 1, 2).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

This Spec is delivered after Spec 0205, whose judge task_02 re-uses.

## Risks & Considerations

- **Wrong-reason matches.** A word such as "golden" in an unrelated log
  classifies an item. Every item names its signature, and a reviewer can
  dispute it from the record alone.
- **Signatures fitted after reading.** The table is fixed by this TechSpec
  and task_04 does not change it. When more than a third of items stay
  `unclassified`, the memory rule yields `inconclusive`, and the document
  lists the most frequent unclassified checks for a later table revision.
- **Few repaired Tasks.** With about 30 repaired Tasks the interval is wide;
  the rule then yields `inconclusive`, which is a valid result.
- **Live measurement needs the keys and the real home.** Task 04 runs where
  the Run Database and the keys are; without a key the record is `blocked`
  and the document says so.
- **Large Run Database.** `RunEventsOfKinds` filters in SQL, so agent
  payloads are never loaded.

## Decisions

- Deterministic signatures from a closed list, read-only. See ADR-0214.
- Measure, record, recommend, remove nothing. See ADR-0215.
- One report reads the Run Database; the harness reads its JSON.
- The harness is test-only, so the binary gains no sender of Task files.
