---
name: implement-spec
description: Implement every pending task of a spec in dependency order — resolve the _tasks.md graph into waves and run the implement-task cycle per task, starting immediately on invocation with no confirmation prompts, and finish by running the qa-gate against the whole feature.
disable-model-invocation: true
argument-hint: "<spec slug or path under docs/specs/> [--from task_NN]"
metadata:
  category: implementation
  tags: [workflow, agents, coding, issues]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Implement Spec

Drive a whole spec to completion by looping over its task graph, one task end-to-end at a time. The artifacts do the coordination: `_tasks.md` supplies the order, each `task_NN.md` supplies the contract and carries the status, and this loop just walks the graph — which is why a killed loop resumes exactly where it stopped.

**Invoking this skill is the authorization to run the whole loop.** Do not ask for confirmation, do not wait for further instructions, and do not present the plan as a question ("Please confirm and I'll start") — the human gates already happened: the breakdown was approved in `write-tasks`, and every task carries its own verification gate. A confirmation prompt here only stalls an AFK loop.

## 1. Resolve the spec

`$ARGUMENTS` names the spec (slug or path under `docs/specs/`). Read `_prd.md`, `_techspec.md` if present, and `_tasks.md`. Parse the graph from the manifest frontmatter and the current `status` of every task file it lists — statuses live in the task files, never in the manifest.

`--from task_NN` skips everything before that task (its `needs` must already be `completed`).

## 2. Plan the waves

Topologically sort the graph: wave 1 = tasks with no `needs`; wave N = tasks whose `needs` are all satisfied by earlier waves. Tasks already `completed` are skipped — re-running the loop is idempotent. A task with a `failed` dependency is blocked, not skippable.

Print the wave plan as a statement, not a question, and **start executing it immediately** — no confirmation before the loop and none between tasks. A ready, unblocked, fully specified task is never a reason to pause; pause only for the stop conditions in section 5.

## 3. Loop

For each ready task, in wave order:

1. Run the full **implement-task** cycle (load → plan → implement → verify with fresh evidence → record → one conventional commit). That skill is the single source of the per-task discipline; do not improvise a lighter version of it inside the loop.
2. On `completed`: continue to the next ready task automatically.
3. On `failed`: stop the wave and surface it — a failed dependency poisons everything downstream, so pressing on builds on sand.

Run tasks **sequentially by default**. Run tasks of the same wave in parallel only when the user asked for it and the tasks touch disjoint files (check their Overviews and diffs-so-far); parallel sessions sharing a worktree corrupt each other's commits.

Context hygiene: each task deserves fresh attention. When subagents are available, hand each task to a fresh one (pointing it at the spec folder and the implement-task skill); in a single session, keep between-task narration minimal so the window stays sharp for the last wave.

## 4. After the last task — QA gate

When every task in the graph is `completed`, run the **qa-gate** skill against the spec. It validates the assembled feature against the PRD's user stories on the running app — browser-driven when the PRD's `surfaces` include `frontend` — and writes its report under `docs/specs/<slug>/qa/`.

- QA verdict **pass** → the spec is complete: run **archive-spec** on it automatically (its gate is completion — all tasks done plus this QA pass — not a merged PR), then report the loop complete. Archiving is the loop's last step; `archive-spec` itself closes by suggesting the publish step, so don't suggest a PR from here.
- QA verdict **fail** → stop and surface the failures. Fixes re-enter as new tasks in the graph or as direct fixes the user directs — the loop does not silently patch and re-run.

## 5. Stop conditions — fail loudly

Stop and surface (never silently skip, fake, or weaken anything) when:

- a task's verification cannot reach green after honest root-cause attempts;
- a ready task's requirements need a product/design decision the spec doesn't contain;
- the work would exceed the task's slice (the PRD marks it out of scope);
- the worktree contains unrelated changes that staging can't cleanly avoid.

Report progress as a checklist (`task_NN ✓ / in_progress / failed / blocked`) only at natural stop points. Resume later with `--from`.

## Guardrails

- One commit per task; one `ma/` branch per spec (create it if missing) — no per-task branches unless the user asks.
- Pushing and PR creation are explicit user actions; the loop never publishes on its own.
- Never run destructive git commands (`reset`, `restore`, `clean`, `checkout --`) without explicit permission.
- Never edit `_prd.md`, `_techspec.md`, or `_tasks.md` from inside the loop — spec amendments are a human decision made outside it.
