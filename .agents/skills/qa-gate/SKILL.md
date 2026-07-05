---
name: qa-gate
description: Run the final QA gate for a completed spec — validate every user story and acceptance criterion against the running application with real flows, driving a real browser (Chrome DevTools MCP) when the PRD surfaces include frontend, and write an evidence report to docs/specs/<slug>/qa/. Use after the last task of a spec completes, before opening the feature's PR, or when the user asks to "QA this feature", "validate the spec", or "run the QA gate".
metadata:
  category: qa
  tags: [qa, testing, browser, workflow]
  version: 0.2.0
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
- From the task files: the union of all acceptance criteria — and their `## Result` evidence. When a task's Result carries daemon-verified evidence (a named passing test behind a verification-gated commit), credit that task-level criterion from the pointer once the static gate is green; do not re-prove it by hand. QA effort belongs to what per-task verification cannot see: cross-task integration, user stories as real flows, and Non-Goals.
- From `qa/`: previous reports — on a re-run, verify previously failed and blocked items first.

Every user story gets at least one real-flow check. Write the plan as a checklist before executing it; an unplanned QA session drifts toward the happy path.

## 2. Static gates first

Run the repository's verify pipeline (`make verify`, or the documented build/lint/typecheck/test equivalents).

On failure, diagnose before deciding: is it **code-caused** or **environment-caused** (sandbox limits, host git/gpg config, missing local tools, network policy)? Environment-caused means the same failure disappears in a properly isolated invocation or reproduces identically on the unmodified base commit — capture that proof. Code-caused → report and stop; flow-level QA on a broken build hides the real blocker. Environment-caused → record the affected checks as `blocked` with the proof, continue the sweep, and cap the verdict at `partial`; never let an environment problem masquerade as a code failure, and never let it pass silently either.

## 3. Exercise real flows

Start the application the way a user runs it (dev server, built binary, API process). Then, per surface in the PRD's `surfaces`:

| Surface          | How to exercise                                                                                                                                                                                                                              | Evidence to capture                                                                                           |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `frontend`       | Drive a real browser via the Chrome DevTools MCP tools (`new_page`, `navigate_page`, `take_snapshot`, `click`, `fill`, `wait_for`, `take_screenshot`, `resize_page`, `list_console_messages`, `list_network_requests`). Walk each user story as that actor, then run the three frontend sweeps below. | Screenshot per key state into `qa/evidence/`; console free of errors; no failed network requests on the flows |
| `backend`        | Real requests against the running API (`curl`) — happy path plus every failure mode the criteria name.                                                                                                                                       | Status codes, response payloads, observable side effects                                                      |
| `cli`            | Run the real binary/commands a user would type.                                                                                                                                                                                              | Exit codes and output                                                                                         |
| `data`           | Run the migration/sync against a scratch database.                                                                                                                                                                                           | Row counts, shape checks, idempotency on re-run                                                               |
| `infra` / `docs` | Apply to a scratch target / read as the intended audience.                                                                                                                                                                                   | Command output / gaps found                                                                                   |

### Frontend sweeps (mandatory when `frontend` is in surfaces)

1. **Responsiveness.** Exercise the PRD's named breakpoints; absent them, three canonical viewports via `resize_page` — mobile 375×812, tablet 768×1024, desktop 1440×900. The body never scrolls horizontally, controls never overlap or clip, and the primary flow completes at every size. Screenshot each.
2. **UI/UX.** The empty, loading, and error states the PRD's User Experience section names; visible feedback on every async action; focus order and a keyboard-only path through the primary flow; copy matching the product's vocabulary. A failing screenshot is the best bug report — capture failures as well as successes.
3. **Best practices.** Console free of errors and warnings on exercised flows; no failed network requests; an accessibility pass on key screens from the snapshot tree (inputs labeled, images with alt text, obviously readable contrast); no visibly blocking loads on the happy path (`lighthouse_audit` when available).

**Terminal frontends (TUI)** map the same sweeps: responsiveness = render at small, large, and degenerate terminal sizes without breakage or panics; UI/UX = key flows work as the footer hints promise, focus and follow semantics hold; best practices = readable without color, aligned monospace output, clean exit on interrupt.

### Inside a Roundfix spec Run (daemon-assigned QA step)

- Never commit or push — the Daemon owns the QA Report commit; your job ends at writing the report.
- You may be sandboxed. When an operation outside the workspace fails with a permission error (writes to `$HOME`, network, nested tool state), classify it immediately as environment-caused, record the affected criterion as `blocked` with the error text, and move on — never retry-loop a sandbox denial. Note in the Environment section which checks need a full-access session.

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

## Failures and blocked items
<!-- One block per failure: expected vs actual, where observed, repro steps.
     One line per blocked item: what blocked it (with proof) and what unblocks it. Omit when none. -->

## Environment
<!-- How the app was run, browser/tooling used, data setup, sandbox constraints hit.
     Enough to reproduce this session. -->
```

Verdict mapping: `pass` only when **every** criterion passes; `partial` when any criterion is `blocked` (environment or an explicitly partial run) and none failed; `fail` when any criterion failed. On `fail`, list each failure as an actionable item — the fixes become new tasks or direct repairs, and the gate runs again. **The gate blocks the PR**: do not proceed to PR preparation while the latest report fails.

## Anti-patterns

- Declaring pass from the test suite alone — step 2 is the floor, not the gate.
- Skipping browser QA because "the components have tests" — integration is exactly what component tests don't see.
- Happy-path-only screenshots when criteria name error states.
- Softening a fail into "minor issue, noting for later".
- Fixing bugs mid-QA and continuing on the mutated build — record, finish the sweep, then fix and re-run.
- Retry-looping a sandbox-denied operation instead of recording it as blocked once.
- Re-proving daemon-verified task evidence item by item while user-story flows go unexercised.
- Labeling an environment failure as a code `fail` — or letting it slide into a `pass`.
