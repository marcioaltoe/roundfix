---
description: Execute a single .scratch PRD/issue task end-to-end, with mandatory PRD + CONTEXT.md reading and the make verify gate
argument-hint: <ONE issue path> — a SINGLE task only, e.g. .scratch/<feature>/issues/04-*.md (never a glob or list)
---

**Input:** `$ARGUMENTS` = path to exactly **one** task/issue file in the `.scratch/` tracker
(e.g. `.scratch/production-readiness-ai/issues/04-record-real-staging-cd-evidence.md`).

This command drives the local markdown issue tracker documented in `docs/agents/issue-tracker.md`.
Tracker shape: one feature per `.scratch/<feature-slug>/`, the PRD at
`.scratch/<feature-slug>/PRD.md`, and issues at `.scratch/<feature-slug>/issues/<NN>-<slug>.md` with
a `Status:` line near the top (vocabulary in `docs/agents/triage-labels.md`).

## 0. Single-task guard (do this first)

1. If `$ARGUMENTS` is empty, names more than one path, or contains a glob/`*`/comma-separated list:
   **STOP** and ask the user to point at exactly one issue file. This command operates on a single
   task by design — batching tasks defeats the per-task acceptance and verification gates.
2. Confirm the path resolves to a real file under `.scratch/<feature>/issues/`. If it points at a
   `PRD.md` or a directory, stop and ask for the specific issue file.

## 1. Mandatory reading (no skipping)

Read these **in full before any planning, design, or code**. This is a blocking requirement.

1. **The task file** at `$ARGUMENTS` — note its `Status:`, `## Parent`, `## What to build`,
   `## Acceptance criteria`, and `## Blocked by`.
2. **The parent PRD** — resolve from the task's `## Parent` line; if absent, fall back to
   `.scratch/<feature>/PRD.md` (same feature directory as the task). Read it whole.
3. **The project `CONTEXT.md`** at the repo root (single source of truth for the Brazilian fiscal
   domain vocabulary). Read it whole.

If any of these three cannot be found, stop and report what is missing rather than guessing.

## 2. Blocked-by gate (prevents forcing externally-gated work)

Inspect the task's `## Blocked by` section and its `Status:`:

- If it lists **unmet prerequisites** (another open issue, or human/external evidence such as a
  provisioned environment, a live URL, a real CD run, owner acceptance, or secret material that does
  not exist in the repo), and that evidence is not present:
  **do NOT fabricate evidence, do NOT check acceptance boxes, do NOT mark it done.**
  Record the precise missing prerequisite as a blocker (in the task's `## Comments` and, if a
  `.codex/loop/<feature>/` tracker exists, via its `update-tracking.py --block-task/--blocker`),
  then stop. Never substitute local tests for required live/staging evidence.
- A `Status:` of `ready-for-human` or `wontfix` is a stop signal: surface it and confirm with the
  user before doing any implementation.
- Only proceed to implementation when the task is genuinely agent-completable from the current repo
  state and provided evidence.

## 3. Activate required skills (before writing anything)

Follow the Agent Skill Dispatch Protocol in the root `CLAUDE.md` (§Step 1–2). Scan the task and
target files for the domain keywords and activate **every** matching skill before planning or code.
Common pairings for this repo:

- HTTP endpoints → `hono-api-best-practices` + `hono` + `zod` (+ `drizzle-*` if DB).
- DB/schema/migrations → `drizzle-postgres` + `drizzle-orm` (+ `drizzle-safe-migrations`).
- Frontend → `ui-craft` + `feature-systems-pattern` + `react` (read `DESIGN.md`).
- Tests → `testing-boss` (+ `vitest`).
- Bug/debug → `systematic-debugging` + `no-workarounds`.
- Security work → `security-best-practices`.

Skipping skill activation invalidates the task per project rules.

## 4. Plan

Use TodoWrite to break the task into steps derived from its **acceptance criteria** (one todo per
criterion is a good default). Capture: affected files, tests to add/update, the exact verification
command, and any artifact updates the acceptance criteria demand.

## 5. Implement

Implement against the PRD intent and the CONTEXT.md domain vocabulary, following the repo
architecture rules (backend Clean Architecture, frontend `systems/<domain>/`, REST/Hono contracts,
etc.). Write/extend tests in the correct location (`__tests__/` beside source; `tests/` for
integration/E2E). Use pt-BR for any domain dialogue with the user; keep code/identifiers in English.

## 6. Quality gate (evidence before claims)

Activate `verification-before-completion`, then run the full gate and read the actual output:

```bash
rtk make verify
```

It must finish with **zero errors and zero warnings** (oxfmt → oxlint → typecheck → tests).
If anything fails, fix the root cause and re-run until clean. Do not claim completion without fresh
passing output in this session.

## 7. Update artifacts

1. In the task file: check off the acceptance criteria that are now satisfied (`- [x]`), and update
   the `Status:` line to the correct triage label (`docs/agents/triage-labels.md`). Only mark a task
   fully done when **every** criterion is checked and `rtk make verify` passed.
2. Append a dated entry under a `## Comments` heading summarizing what was done, the verify evidence,
   and any remaining/blocked items. Never paste secrets, tokens, private URLs, or raw logs.
3. If the feature has a `.codex/loop/<feature>/` tracker, reflect the transition with the loop
   scripts (`update-tracking.py --start-task/--complete-task/--block-task`, with `--memory-written`
   and `--verify-pass`), and validate with `validate-tracking.py <feature>`.

## 8. Summary

Report concisely (in pt-BR for domain content):

- What was built or recorded, and which acceptance criteria are now satisfied.
- The `rtk make verify` result (with the observed numbers).
- The single most relevant next step — or, if blocked, the exact external/human prerequisite that
  must arrive before this task can advance.
