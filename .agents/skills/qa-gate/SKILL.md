---
name: qa-gate
description: Run the final QA gate for a completed spec — validate every user story and acceptance criterion against the running application with real flows, driving a real browser (Chrome DevTools MCP) when the PRD surfaces include frontend, and write an evidence report to docs/specs/<slug>/qa/. Use after the last task of a spec completes, before opening the feature's PR, or when the user asks to "QA this feature", "validate the spec", or "run the QA gate".
metadata:
  category: qa
  tags: [qa, testing, browser, workflow]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# QA Gate

Validate that the assembled feature behaves as the PRD promises — exercised the way a user would exercise it, not the way the test suite does. Tests passed before anyone pressed the button; this gate presses the button. It is the last check before a PR, and its report is the evidence the PR carries.

**This skill is self-contained by design.** Everything needed is in this file — use the agent's own tools and connected MCP servers directly; do not load other skills to perform the QA.

## Preconditions

Resolve the spec (`docs/specs/<slug>/`). Read `_prd.md` and every `task_NN.md` status. If any task is not `completed`, stop and report which — the gate validates the whole feature, and a partial pass would be read as a full one. Run partially only when the user explicitly asks, and mark the report `partial`.

## 1. Build the QA plan

- From `_prd.md`: the `surfaces` list, every user story, every core feature, and the Non-Goals (things you must _not_ find implemented).
- From the task files: the union of all acceptance criteria.
- From `qa/`: previous reports — on a re-run, verify previously failed items first.

Every user story gets at least one real-flow check. Write the plan as a checklist before executing it; an unplanned QA session drifts toward the happy path.

## 2. Static gates first

Run the repository's verify pipeline (`make verify`, or the documented build/lint/typecheck/test equivalents). Any failure → report and stop. Flow-level QA on a broken build wastes the effort and hides the real blocker.

## 3. Exercise real flows

Start the application the way a user runs it (dev server, built binary, API process). Then, per surface in the PRD's `surfaces`:

| Surface          | How to exercise                                                                                                                                                                                                                              | Evidence to capture                                                                                           |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `frontend`       | Drive a real browser via the Chrome DevTools MCP tools (`new_page`, `navigate_page`, `take_snapshot`, `click`, `fill`, `wait_for`, `take_screenshot`, `list_console_messages`, `list_network_requests`). Walk each user story as that actor. | Screenshot per key state into `qa/evidence/`; console free of errors; no failed network requests on the flows |
| `backend`        | Real requests against the running API (`curl`) — happy path plus every failure mode the criteria name.                                                                                                                                       | Status codes, response payloads, observable side effects                                                      |
| `cli`            | Run the real binary/commands a user would type.                                                                                                                                                                                              | Exit codes and output                                                                                         |
| `data`           | Run the migration/sync against a scratch database.                                                                                                                                                                                           | Row counts, shape checks, idempotency on re-run                                                               |
| `infra` / `docs` | Apply to a scratch target / read as the intended audience.                                                                                                                                                                                   | Command output / gaps found                                                                                   |

Frontend specifics that catch what component tests miss: exercise the empty, loading, and error states the PRD's User Experience section names; if it mentions breakpoints, resize and re-check; screenshot failures as well as successes — a failing screenshot is the best bug report.

## 4. Verdict per criterion

Each user story and acceptance criterion gets `pass`, `fail`, or `blocked` — with a pointer to its evidence. No "should work", no "looks fine": if the output wasn't captured, the check did not happen. Checking a Non-Goal and finding it implemented is a finding too (scope creep shipped).

## 5. Report

Write `docs/specs/<slug>/qa/qa-report-YYYY-MM-DD.md`:

```markdown
---
spec: <slug>
date: YYYY-MM-DD
verdict: pass # pass | fail | partial
surfaces: [frontend, backend]
---

# QA Report — <feature>

## Summary
<!-- 2-4 sentences: what was exercised, overall verdict, notable findings. -->

## Results

| # | Criterion / user story | Verdict | Evidence |
| - | ---------------------- | ------- | -------- |

## Failures
<!-- One block per failure: expected vs actual, where observed, repro steps. Omit when none. -->

## Environment
<!-- How the app was run, browser/tooling used, data setup. Enough to reproduce this session. -->
```

`verdict: pass` only when **every** criterion passes. On `fail`, list each failure as an actionable item — the fixes become new tasks or direct repairs, and the gate runs again. **The gate blocks the PR**: do not proceed to PR preparation while the latest report fails.

## Anti-patterns

- Declaring pass from the test suite alone — step 2 is the floor, not the gate.
- Skipping browser QA because "the components have tests" — integration is exactly what component tests don't see.
- Happy-path-only screenshots when criteria name error states.
- Softening a fail into "minor issue, noting for later".
- Fixing bugs mid-QA and continuing on the mutated build — record, finish the sweep, then fix and re-run.
