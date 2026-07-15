---
name: qa-gate
description: Run the self-contained final QA gate for a completed spec — derive a resumable QA matrix from the PRD and task evidence, execute real user flows through every declared surface, probe high-risk user behavior, capture auditable evidence, classify findings by user impact, and write the spec-local dated QA report. Use after the last spec task, before its PR, or when asked to "QA this feature", "validate the spec", or "run the QA gate". Do not use as a substitute for implementation tests or a broad standalone QA knowledge base.
metadata:
  category: qa
  tags: [qa, testing, browser, workflow]
  version: 0.3.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# QA Gate

Validate the assembled feature against the promises in its spec by exercising the public interfaces a real user reaches. This skill owns the complete gate: plan, execution, evidence, findings, report, and final verdict. Use the agent's tools and connected servers directly; do not load another QA skill or create a separate living QA tree.

## Non-negotiables

1. **Real user seat.** Enter through the same frontend, API, CLI, data workflow, or documented operational path as the intended actor. Internal helpers and code inspection can diagnose a failure, but cannot prove a user story passes.
2. **Proof beyond optimistic state.** A pass requires the expected observable, an independent confirmation through a fresh load or another public read path, persistence across refresh/restart when relevant, and captured evidence.
3. **Resumable evidence.** Create the dated report with every row `pending` before the first check. Update it after each row so an interrupted run resumes from disk instead of repeating completed work.
4. **One honest verdict.** Every planned row ends as `pass`, `fail`, `blocked`, or `skipped`; the report closes with zero `pending` rows.

## 1. Resolve scope and preconditions

Resolve `docs/specs/<slug>/`, then read `_prd.md`, every `task_NN.md`, and prior files under `qa/`.

- Require every task status to be `completed`. If any task remains incomplete, stop and list it. Run a partial gate only when the user explicitly requests one; mark its scope and final verdict `partial`.
- Identify the PRD's declared surfaces, user stories, core features, acceptance criteria, user-experience states, and Non-Goals.
- Read each task's `## Result`. Credit a task-level criterion only when it points to named, reproducible evidence and the current static gate passes. Spend live QA effort on assembled user journeys, cross-task seams, persistence, failure behavior, and scope creep.
- On a rerun, start with previously failed or blocked rows, then run the remaining matrix against the current build.

The scope is complete when every promise and explicit exclusion in the spec maps to a planned row or a documented reason for exclusion.

## 2. Build the QA matrix and open the report

Create a collision-safe report path before execution: `docs/specs/<slug>/qa/qa-report-YYYY-MM-DD-<scope-or-build>.md`, or `qa-report-YYYY-MM-DD-NN.md` when no stable scope or build slug exists. If today's report exists with `status: in-progress` for the same build and scope, resume it; otherwise create a new report and preserve older reports as history.

Add a row for:

- every user story, exercised end to end by a named actor;
- every acceptance criterion not safely credited from task evidence;
- every Non-Goal that needs a scope-creep check;
- each mandatory surface sweep below.

For each row record the actor, entry point, surface, steps, expected observable, independent confirmation, persistence check, evidence path, and status `pending`. Order rows by user impact and blast radius. Select 2-5 relevant behavior probes for each high-risk journey: double submit, refresh or back navigation mid-action, deep-link/reopen, invalid or out-of-order input, session expiry, offline/reconnect, concurrent tabs, or locale/accessibility changes. Choose probes that fit the feature; unrelated probes create noise.

The plan is complete when every story and criterion has coverage, every chosen probe has a reason, and the report contains the full pending matrix.

## 3. Run static gates first

Run the repository's full verification pipeline, `make verify`, and record the exact command and result. Do not substitute build, lint, typecheck, or test equivalents. If `make verify` cannot run, or if any formatting, test, or build check fails, record the verification gate as blocked.

When a command fails, diagnose its source before continuing:

- **Code-caused:** record the failure and stop. Flow QA on a broken build hides the first blocker.
- **Environment-caused:** prove the constraint with the error or an unchanged-base reproduction, mark affected rows `blocked`, continue checks that remain valid, and cap the verdict at `partial`.

The static gate is complete when it passes or every environment-caused block has concrete proof and an explicit unblocking requirement.

## 4. Exercise real flows

Start the application through its production-like entry point. Record the build or commit, commands, URL or binary, data setup, credentials profile, and every parity deviation.

For each matrix row, repeat this loop:

1. Enter as the named actor through a real entry point.
2. Act through the public surface and observe the product's user-facing response.
3. Confirm the result through a fresh load, restart, deep link, second public endpoint/command, or another user-visible read path.
4. Capture evidence at the goal state and every divergence. Store screenshots and focused artifacts under `qa/evidence/<YYYY-MM-DD>-<run-slug>/`; keep exact commands and concise outputs in the report.
5. Update the row immediately. A stall, dead control, infinite spinner, or vanished state is a finding, not a reason to work around the product.

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
- `blocked`: an external prerequisite or human-only leg prevents a valid check; include the proof and exact unblocking action;
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
build: <commit-or-artifact>
status: in-progress # in-progress | closed
verdict: pending # pending | pass | fail | partial
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
<!-- Stories and criteria passed/total; probes and frontend sweeps attempted; Non-Goals checked. -->

## Final verdict
<!-- Written last: one actionable sentence, plus counts by status and impact tier. -->
```

Close with `status: closed` and apply the verdict mechanically:

- `fail` when any row failed;
- `partial` when none failed but any row is blocked/skipped, or the user requested a partial run;
- `pass` only when every row passed and all evidence paths resolve.

The gate permits PR preparation only on `pass`. On `fail` or `partial`, state what must change or be verified before rerunning. In a daemon-assigned Roundfix QA step, write the report but never commit or push; the daemon owns the QA report commit. Daemon-assigned steps may also run sandboxed: when an operation outside the workspace fails with a permission error (writes to `$HOME`, network, nested tool state), classify it immediately as environment-caused, mark the affected row `blocked` with the error text, and move on — never retry-loop a sandbox denial — noting in the environment record which checks need a full-access session.

## Decision examples

- A UI says "Saved", but the record disappears after refresh: `fail`, `Data-Loss`; the optimistic message is not proof.
- The full suite passes, but real OAuth needs a human-controlled account: affected rows `blocked`, overall `partial`, with exact human verification steps.
- A task Result names a passing unit test, while the assembled browser journey also persists after refresh with screenshots: credit the task criterion and pass the user-story row from live evidence.

## Anti-patterns

- Declaring pass from tests, source inspection, route rendering, or optimistic UI alone.
- Writing the report only after the run, losing resumability and uncovered rows.
- Skipping failure states, behavior probes, or responsive checks because the happy path passed.
- Capturing hundreds of screenshots instead of goal states and divergences.
- Treating an environment denial as a code failure, or letting it silently become a pass.
- Softening a failure into a note, leaving rows pending, or claiming coverage without a matrix row.
- Fixing during the sweep and continuing as though every row ran against one build.
