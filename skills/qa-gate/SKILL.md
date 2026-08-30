---
name: qa-gate
description: Execute the self-contained final QA gate as a Spec's authored terminal `qa` Task — derive a resumable QA matrix from the PRD and task evidence, execute real user flows through every declared surface, probe high-risk user behavior, capture auditable evidence, classify findings by user impact, and write the spec-local dated QA report. Do not use as a substitute for implementation tests or a broad standalone QA knowledge base.
metadata:
  category: qa
  tags: [qa, testing, browser, workflow]
  version: 0.0.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.0.2
---

# QA Gate

Validate the assembled feature against the promises in its spec by exercising the public interfaces a real user reaches. The Daemon runs this skill from the unique authored terminal `qa` Task named by `_tasks.md` frontmatter; it is part of the Task Graph, not a per-run request. That node depends on every non-QA leaf and therefore becomes runnable only after the graph it closes settles. This skill owns the complete gate: plan, execution, evidence, findings, report, and final verdict. Use the agent's tools and connected servers directly; do not load another QA skill or create a separate living QA tree.

## Non-negotiables

1. **Real user seat.** Enter through the same frontend, API, CLI, data workflow, or documented operational path as the intended actor. Internal helpers and code inspection can diagnose a failure, but cannot prove a user story passes.
2. **Proof beyond optimistic state.** A pass requires the expected observable, an independent confirmation through a fresh load or another public read path, persistence across refresh/restart when relevant, and captured evidence.
3. **Resumable evidence.** Create the dated report with every row `pending` before the first check. Update it after each row so an interrupted run resumes from disk instead of repeating completed work.
4. **One honest verdict.** Every planned row ends as `pass`, `fail`, `blocked`, or `skipped`; the report closes with zero `pending` rows.
5. **Typed blocked causes.** Record a row that is unreachable for a proved
   environmental cause as `blocked (environment: <cause>)` and count it in
   `rows_blocked_environment`; record a row stopped by a finding as `blocked
   (finding: <id> — waits on <named failure>)` and count it in
   `rows_blocked_finding`; record a row that is unreachable for the reason in
   a matching, pre-run Spec declaration as `blocked (declared: <criterion>)`
   and count it in `rows_blocked_declared`. Keep the three causes separate:
   never fold one into another to make the report or verdict look cleaner.

## 1. Resolve scope and preconditions

Resolve `docs/specs/<slug>/`, then read `_tasks.md`, `_prd.md`, every
`task_NN.md`, and prior files under `qa/`. Confirm that the manifest's `qa:`
field names the current `type: qa` Task and that this node is terminal and
depends on every non-QA leaf. A missing or mismatched authored node is a graph
defect; do not run the gate outside that node.

- A clean authoring check is a precondition of the gate, not a substitute for
  it. For every active, non-legacy Spec, run:

  ```bash
  roundfix spec check <slug> --strict
  ```

  The gate keeps this non-probing form because the probe asks whether a
  command already passes before its work exists, while the gate runs after
  every Task is complete. Probing completed Tasks would therefore report their
  already-passing commands as vacuous and refuse the gate. The three authoring
  skills use `--run-verification` while the work is still being authored, when
  that question is answerable.

  Stop before building the matrix when the command fails, and write the refusal
  before you stop. A gate that stops here measured no requirement, so its whole
  report is the refusal:

  ```markdown
  ---
  verdict: fail
  rows_blocked_precondition: 1
  auditing_binary: "<version-and-build-identity>"
  auditor_staleness: "<state>: <reason>"
  rows_blocked_environment: 0
  rows_blocked_finding: 0
  rows_blocked_declared: 0
  precondition_check: "roundfix spec check <slug> --strict"
  precondition_reason: "<every refusing code and the sentence beside it>"
  ---

  # QA Report

  ## Results

  | # | Status | Provenance |
  | - | --- | --- |
  | 0 | blocked | precondition |

  ## Precondition refusal

  - check: roundfix spec check <slug> --strict
  - reason: <every refusing code and the sentence beside it>
  ```

  Row `0` is the entire matrix: `blocked` because nothing ran, provenance
  `precondition` because the stop is a refusal and not a measurement. Name every
  refusing `SC-*` code in `precondition_reason`, code first and its sentence
  after, joined with `; ` and collapsed onto one line, so a maintainer reading
  the report knows whether to fix the Spec or the tree. Justify the row in prose
  as a list, never as a second table under `## Results`. An empty Results table
  is not an option: it refuses every later run on this report instead of on the
  Spec, and the prescribed repair — materialize every planned row — is one a run
  that never built a matrix cannot perform.

  A gate that stopped here writes nothing else: no matrix, and no rows for the
  checks that never executed.

  When the command passes, record the clean result in the report's scope and
  environment, but do not create QA matrix rows that repeat the authoring rules
  in the table below. A completed or archived legacy Spec is exempt from forced
  constraint backfill; keep every legacy artifact byte-identical and retain any
  governance row that has no clean authoring result to replace it. Absence of a
  Project Constraints section by itself is not proof of legacy status.

| Authoring rule removed from the QA matrix | `spec check` equivalent |
| --- | --- |
| The authored Task Graph parses, maps PRD promises to Tasks, and resolves declared repository paths. | The Task Graph load required by `SC-COVERAGE-UNTASKED`, plus `SC-COVERAGE-UNTASKED` and `SC-REF-UNRESOLVED` |
| PRD promises map into the TechSpec. | `SC-COVERAGE-UNMAPPED` |
| Every required Project Constraint row exists. | `SC-CONSTRAINT-MISSING` |
| A declared Project Constraint applicability has a reason. | `SC-CONSTRAINT-UNREASONED` |
| A Project Constraint row cites an operative source. | `SC-CONSTRAINT-SOURCE` |
| A cited tooling authorization names the Spec. | `SC-TOOLING-UNAUTHORIZED` |
| Applicable tooling authority declares bounded files. | `SC-TOOLING-UNBOUNDED` |
| A tooling authorization record states its grant in readable fields. | `SC-TOOLING-UNTYPED` |
| Active ADR obligations are listed, related decisions are accounted for, and attributed claims match the cited record. | `SC-ADR-UNLISTED`, `SC-ADR-RELATED`, and `SC-CITATION-UNSUPPORTED` |
| Task requirements do not contradict each other, rehearsals are declared, and Verification can distinguish Task work from no work. | `SC-REQUIREMENT-CONTRADICTORY`, `SC-REHEARSAL-UNDECLARED`, and `SC-VERIFY-WORK-INDEPENDENT` |
| Emitted vocabulary is documented through the TechSpec's Vocabulary Contract. | `SC-VOCABULARY-UNDOCUMENTED` |
| Repository loop declarations and Finding, Rollup, archive, and promoted Backlog lifecycles agree. | `SC-LOOP-ORDER-DIVERGENT`, `SC-FINDING-LIFECYCLE`, `SC-ROLLUP-MEMBER`, `SC-ARCHIVE-LICENSE`, and `SC-BACKLOG-UNMOVED` |

A named detector in the checker's skipped list did not run. It is not an
equivalent, so retain the corresponding QA row. Do not extend this mapping by
analogy: applicability that the checker did not parse, missing tooling
authorization, outside evidence, Non-Goals, the report contract, current Task
status, and live control or chronology behavior stay in the gate unless a
named checker rule decides them.

- Keep the commit-dependent tooling audit as matrix rows. Identify every Task
  that actually creates, edits, renames, moves, or deletes repository-tooling
  configuration, scripts, ignore files, plugin declarations, or version pins.
  Execute each row as commands, not as a judgement over Spec or Result prose:
  resolve the actual paths from the Daemon-owned Task commit and any current
  worktree delta, then resolve the authorization, prerequisite-fix, and
  consequent-fix commits in chronological ancestry. Use
  `git diff-tree --no-commit-id --name-only -r <commit>` for every committed
  change rather than trusting a reported file list.
- Report every post-commit authorization-shape problem together in the same
  failed command row: missing, late, or untraceable authorization;
  authorization or a prerequisite fix folded into the Task commit; a
  consequent fix folded into or ordered before the change that caused it; a
  claimed derived pin without reproducible evidence from the sanctioned
  regeneration command; and every path outside the exact bounded list plus the
  assigned Task file. Do not stop at the first problem and defer the rest to
  another QA rerun. Any problem blocks flow QA after the complete command audit
  has been reported.
- QA writes only its report and evidence. It never changes Task status or Task
  Graph dependencies.
- Require every dependency of the authored `qa` Task to be `completed`. If any
  dependency remains incomplete, stop and list it. The Daemon owns the current
  gate Task's status and settles it from the final verdict. Run a partial gate
  only when the authored QA Task explicitly limits its scope; mark its scope
  and final verdict `partial`.
- Identify the PRD's declared surfaces, user stories, core features, acceptance criteria, user-experience states, and Non-Goals.
- Read every entry under the PRD's `## Unreachable Acceptance` section before
  the run. Treat each declaration as an author's claim to test, not as proof
  and not as permission to omit the row. The gate may match a blocked row only
  to one of these pre-run declarations; it must never declare unreachability
  itself or infer a declaration from a failed attempt.
- Read each task's `## Result`. Credit a task-level criterion only when it
  points to named, reproducible evidence and the static checks it depends on
  pass. An unrelated static failure does not invalidate that evidence. Spend
  live QA effort on assembled user journeys, cross-task seams, persistence,
  failure behavior, and scope creep.
- On a rerun, start with previously failed or blocked rows, then run the remaining matrix against the current build.

The scope is complete when every promise and explicit exclusion in the spec maps to a planned row or a documented reason for exclusion.

## 2. Build the QA matrix and open the report

Create a collision-safe report path before execution:
`docs/specs/<slug>/qa/qa-report-YYYY-MM-DD.md` for the day's first report, then
`qa-report-YYYY-MM-DD-NN.md` with the next unused numeric `-NN` suffix for
same-day reruns. Numeric same-day suffixes are the only allowed suffixes; never
use a scope or build slug. Resume an existing `status: in-progress` report only
when it is for the same build; otherwise create the next numeric sibling and
preserve older reports as history.

Read the Pull Request fact in the Roundfix QA prompt before planning Pull
Request journeys:

- When the fact names an Open Pull Request, those journeys are runnable.
  Observe that Pull Request read-only through the existing `gh` and Review
  Source boundaries. Check approval, checks and status evidence, unresolved
  review threads, Merge-Ready acceptance, and the Daemon's separate
  review-artifact commit ancestry where the Spec requires them. Never mutate
  the Pull Request, resolve threads, commit, or push.
- When the fact says no Pull Request is open, record every Pull Request journey
  as `blocked (environment: no open Pull Request)` and count it in
  `rows_blocked_environment`. When the fact says the Pull Request could not be
  resolved, the absence is unproven: record the cause as
  `blocked (environment: Pull Request unresolved)` and never write it up as a
  confirmed absence. Do not try to resolve a Pull Request from the Run
  Worktree branch: that per-Run branch is never pushed and has no Pull Request
  of its own.

Add a row for:

- every user story, exercised end to end by a named actor;
- every acceptance criterion not safely credited from task evidence;
- the Spec's outside-evidence acceptance row, when no row above already carries it;
- every Non-Goal that needs a scope-creep check;
- each mandatory surface sweep below.

For each row record the actor, entry point, surface, steps, expected observable,
independent confirmation, persistence check, evidence path, and status `pending`.
Order rows by user impact and blast radius. Select 2-5 relevant behavior probes
for each high-risk journey: double submit, refresh or back navigation
mid-action, deep-link/reopen, invalid or out-of-order input, session expiry,
offline/reconnect, concurrent tabs, or locale/accessibility changes. Choose
probes that fit the feature; unrelated probes create noise.

### Row input declaration

A row opts into future evidence-scoped carry-forward by adding a non-empty,
typed `inputs:` declaration to its detailed evidence block. Place the block
under a `### <row-id>` heading whose row identifier matches the row's `#` cell
in the Results table, and use one fenced `yaml` block per row; a block under
any other heading is ignored and the row is never carried. Each entry has a
`kind` and a `ref`:

```yaml
inputs:
  - kind: repository_path
    ref: <repository-relative path or glob>
```

Use one entry for every input the row's truth depends on:

| Kind | Use when |
| --- | --- |
| `repository_path` | The evidence depends on content at a literal path or glob in the current repository. |
| `external_repository` | The evidence depends on content or state in another repository. The row is never carriable. |
| `live_service` | The evidence depends on state observed from a live service. The row is never carriable. |
| `elapsed_time` | The evidence depends on elapsed time, age, duration, or a time window. The row is never carriable. |

Declare every applicable kind. A mixed list containing any non-repository
input is never carriable. A row with no `inputs:` declaration or an empty list
is also never carriable and must be re-observed, so carry-forward remains
opt-in and fail-closed.

A report without `inputs:` behaves exactly as it does today; existing rows,
counts, statuses, verdict rules, and report naming do not change. When a future
round closes, a `pass` row whose declaration carries only repository inputs,
whose ancestry is proven, and whose evidence is byte-identical at the later
head is materialized as `carried (established by: <report>; head: <sha>)` and
is not re-observed.

One row is the Spec's outside-evidence row: the acceptance row that rests on
evidence originating outside the Spec's own artifacts — a repository the Spec
did not build, a measurement it did not design, or published literature. Record
in that row where its evidence came from, named precisely enough for a later
reader to reach the same source, so the result cannot be read as a rehearsal of
the Spec's own premise. When that source cannot be obtained, record the row as
`blocked (environment: <cause>)` with the reason it was unreachable and count it
in `rows_blocked_environment`. Never drop the row, and never satisfy it with
evidence the Spec authored. A blocked or partial outside-evidence row blocks
pull request preparation until the row is satisfied or carried forward on
declared unmoved evidence under ADR-0097. Task authoring never stalls on it —
decomposition records the blocked row and proceeds — so the obligation lands
here, at the gate, where the Spec is asked to account for it. See ADR-0104.

The plan is complete when every story and criterion has coverage, every chosen probe has a reason, and the report contains the full pending matrix.

## 3. Run static gates first

Run the repository's full verification pipeline, `make verify`, and record the exact command and result. Do not substitute build, lint, typecheck, or test equivalents. If `make verify` cannot run at all, record the verification gate as blocked. A formatting, test, or build check that runs and fails is a `fail`, not a block — classify it with the code-caused and environment-caused distinction below before recording anything.

When a command fails, diagnose its source before continuing:

- **Code-caused:** record a finding and trace the failed check to the matrix
  rows whose entry point, observable, or evidence depends on what failed.
  Block only those implicated rows as `blocked (finding: <id> — waits on
  <named failing check>)`; continue every unimplicated row whose result
  remains valid. A failure may block the whole matrix only when every row
  genuinely depends on it, such as a build failure that leaves no runnable
  artifact; record that dependency on each row instead of asserting it once
  for the matrix.
- **Environment-caused:** prove the constraint with the error or an unchanged-base reproduction, record affected rows as `blocked (environment: <cause>)`, continue checks that remain valid, and apply the typed blocked-cause verdict rule in section 6.

The static gate is complete when it passes or every failure is classified,
every implicated row names the check it waits on, every environment-caused
block has concrete proof and an explicit unblocking requirement, and the
remaining rows are identified as safe to run.

## 4. Exercise real flows

Start the application through its production-like entry point. Record the build or commit, commands, URL or binary, data setup, credentials profile, and every parity deviation.

For each matrix row, repeat this loop:

1. Enter as the named actor through a real entry point.
2. Act through the public surface and observe the product's user-facing response.
3. Confirm the result through a fresh load, restart, deep link, second public endpoint/command, or another user-visible read path.
4. Capture evidence at the goal state and every divergence. Store screenshots and focused artifacts under `qa/evidence/<YYYY-MM-DD>-<run-slug>/`; keep exact commands and concise outputs in the report.
5. Update the row immediately. A stall, dead control, infinite spinner, or vanished state is a finding, not a reason to work around the product.

When a row cannot reach its expected observable, compare the row and the
proved cause with the pre-run declarations. Use `blocked (declared:
<criterion>)` only when the row matches a declared criterion and the observed
constraint matches that declaration's reason. A different circumstance stays
environment-blocked or finding-blocked under its actual cause. A blocked row
with no matching declaration keeps its existing cause and remains blocking
under the existing verdict rules.

If the gate can reach a declared row, run it normally and report the
declaration as a wrongly-declared-row finding. Do not accept, ignore, or remove
the declaration merely because the journey passes: the reachable journey is
evidence that the Spec's claim was wrong, and that finding prevents `pass`.

### Surface protocol

| Surface    | Exercise                                                                                                      | Minimum evidence                                                                                                                       |
| ---------- | ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| `frontend` | Drive a real browser through each actor's journey; re-read interactive state after navigation or DOM changes. | Key-state and failure screenshots, independent confirmation after refresh/deep-link, console errors/warnings, failed network requests. |
| `backend`  | Send real requests to the running documented API, including named failure modes.                              | Exact request, status, payload, and a public read that confirms side effects.                                                          |
| `cli`      | Run the installed or built command exactly as an operator would.                                              | Command, exit code, output, and a second command/read confirming state.                                                                |
| `data`     | Apply the real migration or sync to a scratch database and rerun it.                                          | Shape/row checks, preserved distinctions, side effects, and idempotency evidence.                                                      |
| `infra`    | Apply or inspect through a safe scratch/dry-run path that matches deployment.                                 | Command output, resulting configuration/state, and parity gaps.                                                                        |
| `docs`     | Follow the documentation as its intended reader without undocumented knowledge.                               | Commands attempted, result, and any missing or misleading step.                                                                        |

### Frontend sweep

When `frontend` is declared, cover all of these in the matrix:

- **Viewport and input:** use named breakpoints or mobile 375×812, tablet 768×1024, and desktop 1440×900. Complete the primary flow at each size; check clipping, overlap, horizontal scrolling, keyboard-only operation, visible focus, and 200% zoom on key screens.
- **States and recovery:** exercise the PRD's empty, loading, success, validation, and error states. Async actions show prompt feedback, failures preserve input, and recovery offers a concrete next step.
- **Accessibility and language:** inputs and controls have accessible names, heading/focus order is coherent, color is not the only signal, images have appropriate alternatives, and product copy uses the spec's vocabulary.
- **Runtime health:** exercised flows leave no unexplained console errors/warnings or failed requests. Run a focused accessibility or Lighthouse audit when available; disclose when unavailable.

For a TUI, map the same sweep to small/large/degenerate terminal sizes, keyboard flows, focus/follow behavior, color-independent readability, aligned output, and clean interrupt handling.

## 5. Classify results and findings

Assign each row exactly one terminal status:

- `pass`: the expected observable, independent confirmation, persistence check, and evidence all exist;
- `fail`: observed product behavior contradicts the spec or a Non-Goal shipped;
- `blocked (environment: <cause>)`: the gate environment makes the journey
  unreachable; include the proof, exact unblocking action, and equivalent
  observed or supervised evidence when available;
- `blocked (finding: <id> — waits on <named failure>)`: a finding prevents the
  journey from reaching the expected observable; name the failed check or
  behavior the row waits on, then link the finding and affected evidence;
- `blocked (declared: <criterion>)`: the row is unreachable for the reason in
  the matching pre-run Spec declaration; link the declaration, the proof that
  the gate could not reach it, and the named human action that would satisfy
  it;
- `skipped`: allowed only for an explicitly partial gate or a risk-based cut; state why and what remains unverified.

Classify each finding by user impact:

- `Blocks-Completion`: the actor cannot finish a value-delivering journey;
- `Data-Loss`: user data is destroyed, corrupted, or silently discarded;
- `Trust-Damage`: contradictory, inaccessible, or unrecoverable behavior undermines confidence;
- `Friction`: the journey completes with avoidable confusion, delay, or repetition;
- `Cosmetic`: presentation differs without affecting completion or trust.

For every failure, record the actor, journey step, expected and actual behavior, reproduction from the public entry point, evidence, impact tier, and affected rows. Search previous reports before naming a finding as new; mark a repeated or regressed symptom instead of splitting its history.

Finish the planned sweep before fixing code. A changed build invalidates remaining comparisons. After repairs, start a new run or clearly version the build, rerun the static gate, the failed flow, and one adjacent canary.

## 6. Close the report

Use this structure:

```markdown
---
spec: <slug>
date: YYYY-MM-DD
build: <audited-commit-or-artifact>
auditing_binary: "<version-and-build-identity>"
auditor_staleness: "<state>: <reason>" # state is current|stale|unknown; the reason names the signal that answered, such as commit ancestry or the declared tree version
status: in-progress # in-progress | closed
verdict: pending # pending | pass | fail | partial
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_blocked_declared: 0
surfaces: [frontend, backend]
---

# QA report — <feature>

## Scope and environment
<!-- Full or partial scope, app start commands, build, tools, data, actors, parity gaps. -->

## Static gate
<!-- Exact commands and results. -->

## Results
| # | Story / criterion / sweep | Actor and surface | Status | Evidence |
| - | --- | --- | --- | --- |

## Findings
<!-- One block per finding: impact, expected/actual, reproduction, evidence, affected rows. -->

## Blocked and skipped
<!-- Proof, unblocking action, and uncovered scope. Omit only when empty. -->

## Coverage
<!-- Stories and criteria passed/total; rows_blocked_environment, rows_blocked_finding, and rows_blocked_declared; probes and frontend sweeps attempted; Non-Goals checked. -->

## Final verdict
<!-- Written last: one actionable sentence, plus counts by status and impact tier. -->
```

Close with `status: closed` and apply the verdict mechanically:

- `fail` when any row failed or the gate found a wrongly declared row;
- `partial` when none failed and no declaration was wrong, but any row is
  declared-blocked, finding-blocked, or skipped, the authored QA Task defines
  a partial run, or an environment-blocked row lacks equivalent observed or
  supervised evidence;
- `pass` when every runnable row passed, every environment-blocked row records
  its cause and equivalent observed or supervised evidence, no row is
  declared-blocked, finding-blocked, or skipped, no declaration was wrong, and
  all evidence paths resolve. A nonzero
  `rows_blocked_environment` count does not by itself prevent `pass`; a nonzero
  `rows_blocked_finding` or `rows_blocked_declared` count does.

Set `rows_blocked_environment`, `rows_blocked_finding`, and
`rows_blocked_declared` to the exact number of rows with each blocked cause.
Keep all three keys in every closed report, including when any count is zero.

A report that carries the precondition refusal row from section 1 carries
`rows_blocked_precondition` as well, set to the exact number of `blocked` rows
whose provenance is `precondition`, plus `precondition_check` and
`precondition_reason`. A gate that reached its matrix writes none of these three
keys: the refusal is an added shape, not a fourth count every report owes.

The gate permits PR preparation only on `pass`. On `fail` or `partial`, state what must change or be verified before rerunning. In a daemon-assigned Roundfix QA step, write the report but never commit or push; the daemon owns the QA report commit. Daemon-assigned steps may also run sandboxed: when an operation outside the workspace fails with a permission error (writes to `$HOME`, network, nested tool state), classify it immediately as environment-caused, mark the affected row `blocked (environment: <error>)`, and move on — never retry-loop a sandbox denial — noting in the environment record which checks need a full-access session.

## Decision examples

- A UI says "Saved", but the record disappears after refresh: `fail`, `Data-Loss`; the optimistic message is not proof.
- The full suite passes, but real OAuth needs a human-controlled account and no equivalent supervised evidence exists: affected rows `blocked (environment: human-controlled account)`, overall `partial`, with exact human verification steps.
- The Spec declares that publishing a real release is unreachable, and the
  gate proves that its hermetic environment cannot perform that irreversible
  action: the matching row is `blocked (declared: publish the real release)`,
  `rows_blocked_declared` is incremented, and the verdict is `partial`.
- The Spec declares a journey unreachable, but the gate can exercise it: run
  the journey, record a wrongly-declared-row finding, and return `fail` even if
  the journey's product behavior passes.
- A governance check fails but does not affect the runnable application's
  behavior: block only the governance rows that wait on that named check and
  continue the functional journeys.
- The outside-evidence row names repositories this environment does not hold:
  record `blocked (environment: <cause>)` with the attempted lookup as proof and
  count it in `rows_blocked_environment` — never substitute a rehearsal the Spec
  authored.
- The prompt names an Open Pull Request and read-only observation proves approval, Merge-Ready acceptance, and review-artifact ancestry: pass those Pull Request journeys without commit, push, or Pull Request mutation authority.
- A task Result names a passing unit test, while the assembled browser journey also persists after refresh with screenshots: credit the task criterion and pass the user-story row from live evidence.

## Anti-patterns

- Declaring pass from tests, source inspection, route rendering, or optimistic UI alone.
- Writing the report only after the run, losing resumability and uncovered rows.
- Skipping failure states, behavior probes, or responsive checks because the happy path passed.
- Capturing hundreds of screenshots instead of goal states and divergences.
- Treating an environment denial as a code failure, or crediting it without a recorded cause and equivalent evidence.
- Softening a failure into a note, leaving rows pending, or claiming coverage without a matrix row.
- Fixing during the sweep and continuing as though every row ran against one build.
