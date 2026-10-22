---
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# A subsumed test selector is flagged and a QA refusal names its commit

Minimal PRD for a bug fix, minted by `write-techspec`. Two defects share one
theme: the evidence does not name what it checked. A Verification command that
passes a pytest node id together with its own file runs the file whole under
pytest 9, so a test the Task has not written yet goes unnoticed, and `roundfix
spec check --run-verification` reports the command vacuous without saying why.
And a QA Report the Daemon writes when the gate refuses at its precondition
records no audited commit, so a refusal an earlier Run left on the branch reads
like the current result. After this Spec the Spec Consistency Check reports the
overlapping selector from the written command, the `write-tasks` skill states
the rule, and every QA Report the Daemon closes without an Agent turn records
`build`, the commit it audited, as the gate's full report already does. The
TechSpec holds the design.

## Sources

Two Backlog Entries, adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-22-a-pytest-node-id-beside-its-own-file-passes-before-the-test-exists.md`](references/2026-10-22-a-pytest-node-id-beside-its-own-file-passes-before-the-test-exists.md),
   type `fix`, written by this planning change and adopted at birth from
   Quaestor's report of 2026-10-21 at 19:02 BRT: pytest ignores a node id that
   does not exist yet when the same command also passes its whole file, exits
   0, and the check calls the command vacuous without the cause; with only
   the node ids pytest exits 4.
2. [`references/2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md`](references/2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md),
   type `fix`, open in `docs/backlog/` and promoted here: oraculum's Spec 0087
   carried a refusal report without `build:` into three later Runs, and the
   operator read it as the current result four times.

They are grouped because both make the evidence a reader relies on name what
it checked: the selector a Verification really runs, and the commit a QA
refusal really audited. The other open Backlog Entry at the authoring base,
the documentation contract suite's cost, belongs to Spec 0305;
`docs/findings/` holds no entry.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0), with
Specs 0299 and 0300 active, and delivers after them and after Spec 0301 of
this cycle. None of them changes `internal/speccheck`,
`internal/spec/qa.go`, `internal/spec/auditor_evidence.go` or the QA report
functions of `internal/daemon/task_engine.go`; 0299, 0300 and 0301 raise the
Roundfix skill's version, so the documentation Tasks raise each skill's
version over whatever the tree holds when they start.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new check code is a constant beside the existing
  `SC-VERIFY-*` codes and the report key `build` already exists in the full
  report. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the check reads Task files and the
  Daemon reads the Run Worktree's `HEAD`, and every test uses temporary
  repositories. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0315: "A test selector that
  another argument of its invocation contains is refused at authoring", which
  keeps ADR-0148's "One Verification prober refuses vacuous Tasks before the
  Agent and at authoring" unchanged; and ADR-0132: "A refused gate records the
  refusal as its terminal row", whose report gains the audited commit. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config change.

## Goals

- A pending Task whose Verification passes a pytest node id together with a
  collection argument of the same invocation that contains it is reported by
  `roundfix spec check` before any Run, with the selector, the argument and
  the fix named.
- Every other command, a pytest invocation with `--keep-duplicates`, and the
  prober's verdicts stay as they are.
- A QA Report the Daemon closes without an Agent turn records `build`, the
  commit the gate audited; the seed an Agent finishes and every report already
  written keep their bytes.
- The `spec` guide, the Roundfix `spec` reference, the `write-tasks` and
  `qa-gate` skills, the Context-driven development guide and the glossary
  state both rules.

## Core Features

1. The subsumed selector check. `SC-VERIFY-SUBSUMED-SELECTOR`, an error from
   the Tasks stage, reports each node id a pytest invocation of a pending
   non-QA Task's Verification passes together with its file, a directory
   above it, `.`, or a node id that prefixes it, unless the invocation passes
   `--keep-duplicates`. It reads invocations at command position, behind
   environment assignments and `run` wrappers, and inside `$(...)`
   substitutions; option values, redirection targets and globs are not
   collection arguments. It executes no command.
2. The rule in the authoring skill. The `write-tasks` skill gains a section
   "Test selectors and their containing path" that tells the author to name a
   new test's selector alone and run its containing path in a separate
   command, and names the check.
3. The audited commit in a final Daemon report. The Daemon resolves the Run
   Worktree's `HEAD` with the auditor evidence; the precondition refusal
   report and the blocking mechanical report record it as `build`, quoted,
   right after `verdict`. A commit that cannot be resolved writes no key. The
   pending seed an Agent finishes records none, because the Agent's full
   report writes its own.
4. The rule in the gate skill and the guides. The `qa-gate` skill's refusal
   template carries `build`, the Context-driven development guide names it in
   the refusal paragraph, and the glossary defines Subsumed Test Selector and
   revises Precondition Refusal and QA Report.

## Non-Goals / Out of Scope

- Changing the prober's classification or its text, or the Daemon's pre-work
  refusal (ADR-0148).
- Recognizing another test runner; only pytest was measured to drop a
  selector its file contains.
- Rewriting reports already written, adding a reader of `build`, or making
  any command choose the current QA Report by it.
- Changing the full QA Report template, `### QA settlement`, the report shape
  detector, or which report the mechanical stage reads.
- The documentation contract suite's cost (Spec 0305).

## Success Metrics

1. Success Metric: through `CheckStage`, each of ten commands that pass a
   pytest node id with its file, a directory above it, `.`, a prefixing node
   id or a parametrized form, including one behind `uv run --group dev` and
   one inside `out="$(…)"`, yields exactly one `SC-VERIFY-SUBSUMED-SELECTOR`
   error with the exact summary, fix and Verification line.
2. Success Metric: ten commands that do not overlap (a lone node id, a
   sibling file, two sibling node ids, a path that only shares a string
   prefix, `--keep-duplicates`, a node id given to `--deselect`, a value given
   to `-k`, two separate invocations, a `grep` of the same words and the same
   node id twice) yield none; completed and in-progress Tasks yield none, and
   the PRD stage skips the code.
3. Success Metric: a precondition refusal report records `build:
   "<commit>"` right after `verdict`, quoted, exactly once, and stays readable
   by `ReadQAReport`; a refusal written without a resolved commit has no
   `build` key.
4. Success Metric: the Daemon's precondition refusal and blocking mechanical
   reports record the Run Worktree's `HEAD` as `build`, and the pending seed
   records none.
5. Success Metric: each sabotage QA names makes the test whose input exposes
   it fail.
6. Success Metric: after release, Quaestor's `roundfix spec check <slug>
   --strict` reports `SC-VERIFY-SUBSUMED-SELECTOR` for a plan whose
   Verification passes a new node id with its file, and passes once the
   selector stands alone.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Quaestor's report of 2026-10-21 at 19:02 BRT, transcribed by the operator
  before this Spec existed and recorded in the adopted Backlog Entry: the
  node id beside its file exits 0, and alone exits 4.
- oraculum's QA history of 2026-10-21, read read-only by Spec 0297's authoring
  and recorded in the adopted Backlog Entry before this Spec existed: the
  refusal report `qa-report-2026-10-21-01.md` carries no `build:`, while
  reports `-02` to `-04` carry `8736693`, `7472ad34` and `75a9dd0`.
- pytest 9.0.0's changelog, #12083
  (https://docs.pytest.org/en/9.0.x/changelog.html): "`pytest a/b a/` or
  `pytest a/ a/b` are equivalent to `pytest a`; if an argument overlaps
  another arguments, only the prefix remains", and "If you rely on these
  behaviors, consider using `--keep-duplicates`, which retains its existing
  behavior (including the bug)."

## Unreachable Acceptance

- criterion: Success Metric 6, after release, Quaestor's `roundfix spec
  check <slug> --strict` reports `SC-VERIFY-SUBSUMED-SELECTOR` for a plan
  whose Verification passes a new node id with its file, and passes once the
  selector stands alone.
  reason: only Quaestor's repository holds its plans, and no Verification or
  QA row may read another repository or reach the network.
  satisfied-by: Quaestor's next plan check after the operator's release
  message to that project.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2`. pytest 9.1.1, 8.3.5 and 7.4.4 ran in a
  throwaway virtual environment under the authoring scratch directory (the
  Backlog Entry lists every command and exit code): only pytest 9 exits 0
  when a node id rides with its file, a directory above it or `.`, and
  `--keep-duplicates` restores exit 4. `rg -n 'pytest|::' internal/speccheck
  -g '!*_test.go'` found no reader of test runner arguments. The Daemon writes
  the refusal report through `spec.WritePreconditionRefusalReport` from
  `mechanicalQAReportContent`, whose frontmatter has no `build`; a prototype
  of both changes in a scratch worktree passed the tests the Tasks name and
  failed them under each sabotage the QA Task lists. The TechSpec records the
  code it read.
- Secondbrain, read only. Read `wiki/index.md`; `qmd query "pytest node id
  vacuous verification whole file collection argument" --all --files
  --min-score 0.3` returned the Roundfix mirror's reference
  `projects/roundfix/mirror/docs/references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md`,
  which records that commands naming an existing suite or file passed a clean
  checker and were caught only by the probe, the precedent for a structural
  check that names the cause; `qmd query "QA report precondition refusal
  audited commit build" --all --files --min-score 0.3` returned the mirror's
  `docs/user-guide/commands/qa-report.md`, `CONTEXT.md` and this Backlog
  Entry, and no prior decision about `build` in a refusal. No Inbox Entry was
  written, because the briefing forbids writing there.
- Exa: pytest 9.0.0's changelog and release notes
  (https://github.com/pytest-dev/pytest/releases/tag/9.0.0), pull request
  #13704 ("main: add explicit handling of overlapping collection
  arguments", https://github.com/pytest-dev/pytest/pull/13704) and
  `_pytest/main.py`'s `is_collection_argument_subsumed_by`
  (https://docs.pytest.org/en/stable/_modules/_pytest/main.html), which
  confirmed the rule the check copies, the node id parts and parametrization
  cases, and `--keep-duplicates` as the documented opt-out. Challenging: the
  pull request calls `--keep-duplicates` "still subtly broken", so the check
  only exempts it and never recommends it.

## Glossary

- adds: **Subsumed Test Selector**
- changes: **Precondition Refusal**
- changes: **QA Report**

## Decisions

- 2026-10-22 (questions of 15:25 BRT answered by about 18:30 BRT): the
  cycle AP scope "Correções de adotantes (Recomendado)", the order
  "Correções, follow-up, Python (Recomendado)" and the authorization "same as
  cycle AO".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AP briefing for Spec 0302 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The selector rule is a structural check, not a richer vacuous diagnosis:
  the overlap is decidable from the written command, and the prober stays the
  one shared classification; see ADR-0315.
- The check is an error and reads pytest only; its code and term name the
  class without the runner.
- The audited commit is the Run Worktree's `HEAD`, the same value the full
  report's `user_flow_binary` check reads, written quoted after `verdict` in
  both reports the Daemon closes without an Agent; the pending seed keeps its
  bytes so an Agent never inherits a second `build` key.
- An unresolvable commit omits the key rather than failing the refusal, as an
  unanswerable auditor signal does.

## Open Questions

None.
