---
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
prd: _prd.md
created: 2026-10-07
---

# A formatted QA report and a reopen for late dependencies — Technical Spec

## Executive Summary

The QA step is the only place where the Daemon commits files after the last
repository Verification: `runQAGate` runs the precondition
(`internal/daemon/task_engine.go:2798`), lets the Agent write the report, and
commits it (`task_engine.go:2973` → `commitQAReport`, `task_engine.go:3596`).
It also imports an unintegrated prior pass byte for byte before the
precondition (`task_engine.go:2762` → `importPriorQAPass`,
`internal/daemon/qa_prior_pass.go:230`). A repository whose Verification checks
formatting therefore refuses its own next gate, and needs a formatting commit
after every pass.

The fix adds one optional Project Config command, `verification.format`. The
Daemon runs it over the regular files under the Spec's `qa/` directory at two
points: after the mechanical stage and before the precondition, over the
imported pass, and just before the QA Report commit, over the files that commit
stages. A failure restores the original bytes. Reopen gains a second trigger,
a Late Dependency, proven from the Task Graph manifest at the commit that added
the newest QA Report.

The trade-off accepted: Roundfix does not discover the formatter, so an adopter
must configure it, and a formatter that rejects unknown file types must be
given its own ignore option. In exchange no stack is guessed, and no file
outside `qa/` is ever rewritten (ADR-0249).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  Project Config gains `verification.format` and `daemon.qa` gains the phase
  `format`, both snake-case under existing names. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request or credential is added;
  the Format Command and Git run locally, and tests use temporary repositories,
  a temporary home and a shell-script formatter. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0249 (this Spec) governs the Format
  Command, its run points and revert, and the Late Dependency proof. ADR-0194:
  "It copies that commit's report and evidence byte for byte, so a failed pass
  on an unintegrated Run Branch is the next pass's previous report", refined
  by ADR-0249 so the copy and the carry proof stay byte for byte before the
  format. ADR-0195 and ADR-0097 decide what carries and are unchanged.
  ADR-0059: "managed Markdown must survive the repository's selected
  formatter unchanged", a different output this Spec leaves alone. The gate is bound by ADR-0080,
  ADR-0091, ADR-0104, ADR-0156, ADR-0167 and ADR-0240, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by citation and
  receipt. ADR-0182 runs Settlement Checks before each Task commit. ADR-0229
  cites ADR-0167 but decides how an operator archive resumes a park, and ADR-0237
  decides how a Delivery Retry records a merge made outside the queue; this
  Spec changes neither. ADR-0184: "A TechSpec now declares numbered Surface
  Transcripts", answered in the TechSpec with the reason none applies. ADR-0096 decides the gate's machine stage, which
  runs before the format and is unchanged, and ADR-0210 decides the digest an
  Evidence Snapshot records, which the format never reads.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths. Express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário", the standing grant "Concedo" for the Governed Paths each Spec
  declares, and "Pode seguir nessa ordem" (2026-10-07) for this Spec. Bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies/_authorization.md`.

## System Architecture

No new package or command. Four existing components change:

- Config (`internal/config/config.go`): `Verification` (`config.go:177`)
  gains `Format string`; the `verification` overlay (`config.go:443`, key
  switch at `config.go:490`) accepts `format`; the merge (`config.go:1827`)
  lets Project Config replace User Config; the rendered default
  (`config.go:948`) shows `format: ""`.
- Implement wiring (`internal/cli/implement.go`): `executeImplementCycle`
  (`implement.go:1021`, called at `implement.go:559`) passes the value into
  `daemon.TaskPlan` (`task_engine.go:230`), which gains `FormatCommand string`.
- QA step (`internal/daemon`): a new `qa_format.go` holds the format run.
  `runQAGate` calls it between the mechanical stage (`task_engine.go:2775`) and
  `runQARepositoryVerification` (`task_engine.go:2631`, called at 2798) over
  `prior.Files`, and `commitQAReport` calls it after `FilterStageablePaths`
  (`task_engine.go:3631`) and before `Committer.Commit` (`task_engine.go:3651`).
  The import itself (`copyPriorQAPass`, `qa_prior_pass.go:280`) is unchanged.
- Reopen (`internal/cli/reopen.go`, `internal/spec`): today
  `deriveReopenPlan` (`reopen.go:128`) reopens only on `StaleGateError`
  (`internal/spec/spec.go:775-789`) and otherwise refuses in
  `reopenHealthyGateRefusal` (`reopen.go:235`, message at `reopen.go:246`).
  The completed-gate branch now looks for a Late Dependency before refusing.
  A new `internal/spec/gate_closure.go` computes a QA Task's dependency closure
  from manifest bytes through `parseManifestNodes` (`spec.go:598`), and
  `internal/spec/task.go` gains a Late Dependency variant of `ReopenGate`
  (`task.go:177`) and its invalidation record (`appendGateInvalidation`,
  `task.go:213`).

Task settlement is not changed: `settleTask` (`task_engine.go:2084`) settles one
Task, and `roundfix settle` loads the graph with `spec.Load`
(`internal/cli/settle.go:405`), which already refuses a completed gate above a
pending dependency.

## Implementation Design

### Interfaces

```go
// internal/config
type Verification struct { /* existing */ Format string }

// internal/daemon
type TaskPlan struct { /* existing */ FormatCommand string }
const qaFormatTimeout = 5 * time.Minute
// formatQADirectory runs plan.FormatCommand over paths under specQADir.
func (engine *Engine) formatQADirectory(ctx context.Context, plan TaskPlan,
	ordinal int, stage string, paths []string, verdict string) error

// internal/spec
func QAGateClosure(manifestPath string, content []byte) (qaTaskID string, closure []string, err error)
func ReopenGateForLateDependencies(taskPath, reportPath string, taskIDs []string, date time.Time) error
```

### Invariants

1. `verification.format` is a string scalar, default empty; any other YAML
   kind is a config error `verification.format must be a string`. A
   whitespace-only value is empty. Project Config replaces User Config.
2. With `FormatCommand` empty, no format run, Run Event or stderr line exists,
   and both run points leave every byte as today.
3. A format run takes the given paths, keeps only those under
   `<Spec dir>/qa/` that are regular files (by `Lstat`; symlinks and deleted
   paths are dropped), sorts them, and runs nothing when none remain. It
   executes `sh -c '<command> "$@"' roundfix-format <path>...` with the Run
   Worktree as working directory and repository-relative paths, outside the
   Agent sandbox. Paths are arguments, never interpolated into the command.
4. Before the run the Daemon keeps each path's bytes and mode. The run is
   bounded by `qaFormatTimeout`, which cancels and kills the process.
5. The `imported` run happens only when a prior pass was imported, after the
   mechanical stage and before the repository Verification precondition, over
   `prior.Files`. The carry proof has already read the committed bytes.
6. The `commit` run happens in `commitQAReport` after the stageable set is
   final, after the QA Task settled and the evidence snapshots were recorded,
   and immediately before the commit, over the stageable paths. The stageable
   set is not recomputed.
7. The run reverts when the command exits non-zero, cannot start, or hits the
   timeout. It also reverts when the QA Report it covers no longer reads with
   the same verdict through `spec.ReadQAReportFile`: the settled verdict at
   `commit`, or the imported report's verdict at `imported`. To revert, it
   writes every kept path back with its original bytes and mode, and the step
   continues with those bytes. A Stop Request during the run reverts the same
   way and returns the stop.
8. Each run publishes one `daemon.qa` Run Event with phase `format` and the
   payload `stage` (`imported` or `commit`), `outcome` (`formatted`,
   `unchanged`, `failed` or `reverted`), `command`, `paths` (count),
   `changed` (sorted repository-relative paths whose bytes changed), and, for
   `failed` and `reverted`, `reason` and `diagnostics` (the last 2048 bytes of
   combined output). `failed` and `reverted` also write one stderr line
   starting `roundfix: QA format` that names the stage, the outcome and the
   reason.
9. Reopen first behaves as today: a `StaleGateError` reopens over the stale
   dependencies, and a missing or unsettled gate refuses as today.
10. For a completed gate with every dependency completed, reopen resolves the
    newest QA Report (`spec.NewestQAReport`) and its path relative to the Git
    root. It finds the oldest commit in `HEAD`'s history that added that path
    (`git log --diff-filter=A --format=%H HEAD -- <path>`) and reads
    `<spec dir>/_tasks.md` at that commit. A Late Dependency is a Task in the
    QA Task's current closure, read from the working-tree manifest, that the
    recorded closure lacks. An absent commit, an absent or unparsable recorded
    manifest, a report outside the Git root, or no Late Dependency keeps
    today's refusal and message byte for byte.
11. A Late Dependency reopen sets the QA Task `pending` in one atomic
    replacement and appends `## Invalidation` with `- Date:`, `- QA Report:`
    and `- Dependencies added after the QA Report: ` followed by the ids, each
    in backticks. It prints the existing stdout line. The pre-write recheck
    re-derives the same plan, including the Late Dependency ids, or refuses
    with exit 2 as today.

### Data Models

`daemon.qa` gains a phase under an existing kind, and the QA Task file gains an
`## Invalidation` line variant. No schema or journal change.

### API Contracts

1. API Contract: `verification.format`. Per Invariants 1 and 2; `roundfix init`
   renders `format: ""` under `verification:` with a comment.
2. API Contract: QA step formatting. Per Invariants 3 to 8: committed files
   under `qa/` are the formatter's output, or the original bytes after a
   recorded revert.
3. API Contract: reopen over a Late Dependency. `roundfix reopen --spec <slug>`
   exits 0 and prints `reopened <qa id> pending — invalidated <report>` per
   Invariants 10 and 11; every other refusal keeps its exit code and text.

### Surface Transcripts

None. The changed surfaces are a configuration key read only by the Daemon, a
Run Event phase reached inside a QA step, and one new success case of an
existing command whose stdout line is unchanged. They are proved at the config
seam, at the Daemon seam with a temporary Git repository and a shell-script
formatter, through `roundfix implement` with fakes, and through the built
binary for reopen.

## Coverage Map

- Goal 1 → Invariants 3, 6, 8; API Contract 2.
- Goal 2 → Invariant 5; API Contract 2.
- Goal 3 → Invariant 7; API Contract 2.
- Goal 4 → Invariants 9 to 11; API Contract 3.
- Core Feature 1 → API Contract 1.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → API Contract 3.
- Core Feature 4 → Build Order 1.
- Success Metric 1 → Testing Approach 2 and 3.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2 and 5.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.

## Integration Points

- The repository's formatter, invoked as a shell command with file arguments,
  in the pattern lint-staged documents: the command gets the committed paths,
  and filtering is the command's own configuration.
- Git, read-only, for reopen: `log --diff-filter=A` and `show <commit>:<path>`
  in the user checkout.

Research record. The Secondbrain was read first (`wiki/index.md`, then `qmd`
queries on the unformatted report and on reopen). They returned the triaged
adopter report, the adopter's mirrored Spec 0046 and 0047 records, and
Roundfix's Spec 0149 record of the reopen command. Exa returned the lint-staged
README and a report of argument-list limits when a formatter receives
hundreds of paths; a QA directory holds tens.

## Testing Approach

1. Config seam: new `internal/config/verification_format_test.go` covers the
   default, a Project value replacing a User value, a whitespace value, a
   non-string value refused, and the rendered default.
2. Daemon seam: new `internal/daemon/qa_format_test.go`, using the prior-pass
   fixture of `qa_prior_pass_test.go` with `GitCommitter` and a formatter
   script written to a temporary directory. It covers a committed report in
   the formatter's output; an imported unformatted report formatted before a
   precondition that fails on unformatted `qa/` files; a failing formatter; a
   verdict-changing formatter; a Task file and a path outside `qa/` left
   untouched; and an empty command with no `format` event.
3. CLI seam: new `internal/cli/implement_qa_format_test.go` runs
   `roundfix implement` with `verification:\n  format: <script>` in User
   Config and proves the script received the QA Report path.
4. Reopen: new `internal/cli/reopen_late_dependency_test.go` (built binary and
   temporary repository) and `internal/spec/gate_closure_test.go` cover a Late
   Dependency reopen, an uncommitted report, a report committed with the same
   closure, and the record line.
5. Existing `./internal/config`, `./internal/daemon` prior-pass and carry,
   `./internal/cli` reopen and implement, and `./internal/spec` tests pass
   unchanged.

## Build Order

1. The glossary, the configuration and Spec workflow guides, the reopen
   reference and the Roundfix Skill's `implement` and `settle` references
   describe the Format Command and the Late Dependency; the skill version is
   raised and recorded.
2. The Project Config reads `verification.format` (API Contract 1).
3. The QA step runs the Format Command and `implement` wires it (API Contract
   2) (depends on: 2).
4. Reopen proves a Late Dependency (API Contract 3).
5. The final QA gate (depends on: 1, 3, 4).

## Risks & Considerations

- A formatter that also rewrites files it was not given leaves those changes
  uncommitted in the Run Worktree; only the listed paths are committed.
- A formatter that changes a report's tables in a way Roundfix's parsers
  misread is caught by the verdict re-read and reverted. The adopter's third
  Run read a formatter-rewritten report as its previous report and closed.
- A report imported from a pass recorded before the command was configured is
  formatted after its carry proof, so the next pass may re-run those rows.
- Reopen depends on local Git history; a shallow or rewritten history yields
  today's refusal, never a wrong reopen.

## Vocabulary Contract

None. The `roundfix: QA format` stderr line is documented in
`docs/user-guide/configuration.md` by task_01 and asserted by task_03's tests.

Glossary terms adopted in `CONTEXT.md` by task_01: **Format Command** and
**Late Dependency**.

## Glossary

None.

## Decisions

- One configured command, run over `qa/` files only, at the import and the
  commit; revert on failure; see ADR-0249.
- No formatter discovery from the Baseline, `make fmt` or the repository
  Verification; see ADR-0249.
- A Late Dependency is proven from the manifest at the commit that added the
  newest QA Report; the loader is unchanged; see ADR-0249.
