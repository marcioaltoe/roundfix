---
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: active
created: 2026-10-27
surfaces: [backend, cli, docs]
---

# Settle follows the branch's plan and the Task commit's rules, and a declared row matches the criterion it names

Minimal PRD for a bug fix, minted by `write-techspec`. Quaestor's QA gate
recorded two Success Metrics as separate declared-blocked rows that one
Unreachable Acceptance declaration covered, and `roundfix qa-report accept`
refused the report because it compares the number of declared rows with the
number of declarations. After Quaestor split the declaration on its branch,
`roundfix settle` read the old PRD from the kept Run Worktree until Quaestor
moved that worktree by hand. Two more settle defects share the command:
run inside the Run Worktree it takes that worktree for the user checkout, and
it commits an untracked executable that a Task commit refuses. After this
Spec acceptance matches each declared-blocked row to a declaration naming its
criterion, settle reads the plan the branch carries, runs only from the user
checkout, and commits only what a Task commit keeps.

## Sources

Three Backlog Entries, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md`](references/2026-10-27-a-declaration-covering-two-criteria-is-refused-and-settle-reads-the-run-worktree-plan.md),
   `type: fix`, written at birth from the paragraph on Spec 0024 of ai-memory
   message `01a17d0a-cd60-7191-a724-2988affea2d0` from
   `default/gesttione-solutions-quaestor` (2026-10-27 08:07 BRT), quoted
   verbatim. Items 1 and 2 of the same message belong to Spec 0329 and are
   out of scope here.
2. [`references/2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md`](references/2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md),
   `type: fix`, moved from `docs/backlog/`.
3. [`references/2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md`](references/2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md),
   `type: fix`, moved from `docs/backlog/`.

At the authoring base `docs/backlog/` held only entries 2 and 3, and
`docs/findings/` holds no entry.

## Prerequisites

None in `requires`. Authored against `5e9475bf` (v0.125.0) with no active
Spec under `docs/specs/`. Specs 0329 and 0331 are authored in parallel and
deliver before and after this one; neither changes `internal/cli/settle.go`,
`internal/spec/qa.go` or `internal/spec/spec.go`, and if one of them edits a
file a Task here declares, the later Task edits it after that merge and pins
none of its bytes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the PRD's Unreachable Acceptance section accepts a repeated
  `criterion:` field and the QA Report keeps its format. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  network access; the commands read and write a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the new ADR-0343 records the
  decisions. ADR-0240 (one QA partial policy; declared rows covered by the
  declarations) is widened by it, and its count stays as the fallback. ADR-0104, ADR-0167, ADR-0229, ADR-0276
  and ADR-0335 cite ADR-0240 for the other row kinds and keep their rules;
  ADR-0237, ADR-0274, ADR-0295 and ADR-0312 cite those for queue retry, turn
  completion, reachability and resume, and this Spec leaves each unchanged.
  ADR-0340 (Authorized Executables) and ADR-0157 (tracked source is work
  output) are the rules settle now applies.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AR answers of 2026-10-27 quoted in the authorization record, which keep the
  cycle AJ grant for the Governed Paths each Spec declares; the governed files
  are the owned skills `qa-gate`, `write-prd` and `roundfix`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/settle.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A Spec whose single declaration names two criteria has a QA Report with
  two declared-blocked rows accepted; every report accepted today is still
  accepted, and a refused report names the rows no declaration names.
- Settlement, `roundfix qa-report accept` and `roundfix archive` keep one
  acceptance decision, and the QA gate's mechanical stage classifies a
  declared row with the same rule.
- `roundfix settle` reads the plan the Run's branch carries: it moves a Run
  Worktree that the branch is ahead of, and it refuses any other plan drift.
- `roundfix settle` run inside a Run Worktree or Task Worktree refuses and
  names the user checkout.
- `roundfix settle` refuses an untracked executable that a Task commit would
  refuse, and both commit a tracked file whose mode a Task changes.
- `write-prd`, `qa-gate`, the Roundfix skill, the user guides and
  `CONTEXT.md` say each rule.

## Core Features

1. Declared rows matched by criterion. A declaration names one or more
   criteria, one `criterion:` line each. A row
   `blocked (declared: <criterion>)` is covered when a declaration names its
   criterion, by its text or its leading label. Matching only widens
   acceptance: a report today's count accepts (declared-blocked rows not above
   the number of declarations) is still accepted, and a refusal names the rows
   no declaration names.
2. One rule for every reader. The rule lives in `internal/spec`; the
   mechanical stage's declared-row classification calls it.
3. Settle follows the branch. When the Run's target branch changed files
   under the Spec folder after the settle surface forked from it, settle
   fast-forwards a Run Worktree that the branch is ahead of and names the
   move on stderr, and refuses any other drift, naming the files.
4. Settle runs from the user checkout. Settle refuses, with exit `2` and
   nothing written, when its working directory is in a Run Worktree or Task
   Worktree, and names the checkout the Run recorded.
5. Settle commits what a Task commit keeps. Settle refuses an untracked
   executable that the authorization it reads does not list under
   `executables`; a tracked file whose mode changes passes in settle and in
   the Task commit alike.
6. The rules in writing. The `write-prd` PRD template and the `qa-gate`
   skill say that one declaration may name several criteria and that each
   criterion becomes its own QA row; the Roundfix `settle` and
   `spec-delivery` references, the `qa-report` and `settle` guides, both
   command usages, the QA contract in the Agent prompt and `CONTEXT.md` say
   the rest.

## Non-Goals / Out of Scope

- Items 1 and 2 of Quaestor's message (Spec 0329).
- Re-judging archived Specs: no Archive Record changes.
- Matching a paraphrase of a criterion; a paraphrased row is accepted only
  by today's count, as now.
- Advancing a Task Worktree, or merging a diverged branch into a Run
  Worktree: settle refuses and names the files.
- A branch that is ahead of the Run Worktree only outside the Spec folder:
  settle behaves as today and reports `integration pending`.
- Requiring the `executables` grant for a tracked file that gains an execute
  bit (ADR-0157 stands).
- Any text inside a `### QA settlement` section; the Makefile, the CI
  workflows, `go.mod`, `go.sum`; any Baseline module or guide under
  `docs/agents/`.

## Success Metrics

1. Success Metric: `roundfix qa-report accept` exits `0` for a report with
   rows `Q10` and `Q11` recorded `blocked (declared: Success Metric 10)` and
   `blocked (declared: Success Metric 11)` against one declaration with two
   `criterion:` lines; at `5e9475bf` it exits `1` with `rows_blocked_declared
   is 2, but Spec declares 1 unreachable acceptance`. A third row naming
   `Success Metric 12` in place of `Q11` makes it exit `1` with today's
   shortfall words and `declared-blocked row Q12`.
2. Success Metric: a row naming `Success Metric 5` matches the criterion
   `Success Metric 5, after release, …`; a row naming `Success Metric 1` does
   not match `Success Metric 10, …`; backticks, whitespace and a trailing
   period do not change a match.
3. Success Metric: every report today's count accepts is still accepted: a
   row `Q09` recorded `blocked (declared: no later Run fails the test)`
   against the one declaration `no Run delivered after this Spec fails the
   test` matches no criterion and is accepted before and after this Spec, as
   is a report that lists no `blocked (declared: …)` row.
4. Success Metric: `roundfix settle` run inside a kept Run Worktree exits `2`
   with `settle runs from the user checkout` and the user checkout's path,
   and leaves the worktree's HEAD and files as they were; at `5e9475bf` it
   exits `1` with `integration pending`.
5. Success Metric: after the branch amends the PRD, settling the kept Run
   Worktree's QA Task prints `Run Worktree advanced: ` on stderr, accepts
   the report against the amended PRD and exits `0`; at `5e9475bf` it exits
   `1` with the QA Report ineligible. A diverged branch, or a branch change
   that overlaps the surface's uncommitted work, makes settle exit `2` naming
   the files, with nothing written.
6. Success Metric: settle exits `2` naming `newdir/helper.sh` for an
   untracked executable no authorization lists, commits a listed one with
   mode `100755`, and commits a tracked file made executable, as the Daemon's
   Task commit does; at `5e9475bf` the unlisted file is committed.
7. Success Metric: every file `_prd.md` → Core Feature 6 names carries its
   rule, and `CONTEXT.md` revises **Unreachable Acceptance**, **Settle
   Command**, **Run Worktree** and **Authorized Executable**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Git's documentation, read through Exa on 2026-10-27,
  `git-merge` (https://git-scm.com/docs/git-merge): "With `--ff-only`,
  resolve the merge as a fast-forward when possible. When not possible,
  refuse to merge and exit with a non-zero status." and, under PRE-MERGE
  CHECKS, "`git pull` and `git merge` will stop without doing anything when
  local uncommitted changes overlap with files that `git pull`/`git merge`
  may need to update."
- Git's behaviour on the QA host, a tool this Spec did not build, measured at
  authoring with Git 2.56.0: a worktree with a staged new file and an
  unstaged edit fast-forwards to a branch that changed another file, keeping
  both changes; when the branch changed the edited file, `git merge
  --ff-only` prints `error: Your local changes to the following files would
  be overwritten by merge:` and exits `1`.
- The declared-blocked rows of four QA Reports this Spec did not write,
  read from Git at their Archive Records' `source_revision`: Spec 0326's
  `blocked (declared: Success Metric 5)`, Spec 0327's and Spec 0328's rows
  that repeat their criterion, and Spec 0324's
  `blocked (declared: no later Run fails the unreaped-group test)`, which
  paraphrases its criterion. The first three match their declarations under
  the new rule; the last matches none and stays accepted by today's count.
- Quaestor's report, quoted verbatim in the adopted Backlog Entry, written by
  another project before this Spec existed.

## Unreachable Acceptance

- criterion: Quaestor, on the release that carries this Spec, accepts a QA
  Report whose single Unreachable Acceptance declaration names two blocked
  criteria, and settles a Task after amending its PRD without moving the Run
  Worktree by hand.
  reason: it needs another repository's Runs after the release, which no
  Verification or QA row may read.
  satisfied-by: Quaestor's reply through ai-memory to the release message the
  operator sends, recorded session evidence.

## Validation

The decisions were checked on 2026-10-27 against these sources.

- Characterization at `5e9475bf` in a clone, through CLI dispatch on the
  implement test fixtures (`_techspec.md` → Current behavior): settle from
  inside a kept Run Worktree exits `1` with `integration pending`; settle of a
  QA Task after the branch amends the PRD exits `1` with the report
  ineligible against the old PRD; settle commits an untracked `0755`
  `newdir/helper.sh`.
- A prototype of the four Tasks' code in a second clone passed those three
  tests, and across `internal/spec`, `internal/cli`, `internal/daemon`,
  `internal/speccheck` and `internal/agent` moved only the tests
  `_techspec.md` → Declared breaks names, all of them test-only. A variant
  that refused every unmatched row instead of falling back to the count broke
  thirteen `internal/spec` tests, nine `internal/cli` tests and two
  `internal/daemon` tests, and would refuse Spec 0324's archived report; the
  fallback keeps all of them passing unchanged.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix settle stale PRD run worktree" --all --files --min-score 0.3`
  (the Roundfix usage guide mirror),
  `qmd query "unreachable acceptance declaration criteria QA report accept" --all --files --min-score 0.3`
  (QA Reports of two other projects, neither with a declared row) and
  `qmd query "roundfix settle executable file mode" --all --files --min-score 0.3`
  (the mirror of source 3 and the changelog). None records a decision on
  matching declared rows or on the plan settle reads.
- Exa: the `git-merge` documentation quoted above.

## Glossary

- changes: **Unreachable Acceptance**
- changes: **Settle Command**
- changes: **Run Worktree**
- changes: **Authorized Executable**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AR brief for Spec 0330 of 2026-10-27 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-10-27, about 11:40 BRT),
verbatim:

- Scope: "A, B e C (Recommended)" — option text: "Lock + Clean vermelho;
  família do settle + aceite da QA; profile proof sob carga. Todas são health
  e duas se repetem em adotantes."
- QA acceptance: "Casar por linha (Recommended)" — option text: "Uma
  declaração cobre cada criterion que nomeia; o accept confere linha a linha,
  não a contagem. write-prd/qa-gate passam a dizer isso."
- Authorization: "Mesmas regras do AQ (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."

Taken from those answers and the characterization, with no further
interview (ADR-0343):

- A declaration names criteria with repeated `criterion:` lines; a row
  matches by text or leading label; matching only widens acceptance, and
  today's count (declared-blocked rows not above declarations) stays the
  fallback, so no report accepted today is refused (operator review,
  2026-10-27: specs evolve, never regress).
- Settle fast-forwards a Run Worktree when the branch changed the Spec folder
  and is ahead of it, and refuses any other Spec-folder drift.
- Settle refuses inside a Run or Task Worktree instead of settling there.
- Settle applies the Task commit's executable rule with the authorization it
  already reads; a tracked file's mode change stays work output.

## Open Questions

None.
