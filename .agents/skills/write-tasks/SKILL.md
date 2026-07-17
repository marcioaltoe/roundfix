---
name: write-tasks
description: Decompose a spec's PRD/TechSpec into a dependency-ordered task graph — vertical-slice task files plus a machine-parseable _tasks.md DAG manifest under docs/specs/<slug>/, gated by user approval of the breakdown before any file is written.
argument-hint: "<spec slug or path under docs/specs/>"
metadata:
  category: issue-decomposition
  tags: [issues, workflow, prd, agents]
  version: 0.2.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Write Tasks

Turn `docs/specs/<slug>/_prd.md` (and `_techspec.md` when present) into the executable task graph: one `task_NN.md` per task plus `_tasks.md`, the manifest that `implement-task`, `implement-spec`, and any future daemon schedule from. The artifacts are the contract — a fresh agent session must be able to pick up one task file and build it with no other context than the spec folder.

## Preconditions

`$ARGUMENTS` names the spec. `_prd.md` must exist. If `_techspec.md` is missing on a feature with real architectural surface, call that out and say `write-techspec` should run first; use `AskUserQuestion` to get explicit user acceptance before proceeding without the tech spec. Continue only after acceptance, compensating with deeper codebase exploration before decomposing.

## Ownership rules (what lives where, and why)

- **Dependencies live only in the `_tasks.md` graph.** Never in task files. Two copies of an edge is how graphs drift.
- **Status lives only in each task file's frontmatter.** `_tasks.md` owns topology, not progress — schedulers re-read task files for state.
- **Task Type is required in every task file.** Use exactly one of `backend`,
  `frontend`, `data`, `infra`, `docs`, `test`, or `chore`. The Implement
  Command consumes this value for routing and rejects missing or unknown
  values; never leave a placeholder or infer the type later.
- **Content references, never duplicates.** Task files point at PRD/TechSpec sections by name; a task that restates the spec goes stale the first time the spec is amended.
- **Context entries are labeled paths.** Add `## Context` only when the Task
  needs specific instruction or interface paths beyond the standard Spec
  bundle. Use bullets shaped as `- instruction: <path>` or
  `- interface: <path>`. Paths must be clean, repository-relative, and unique;
  a Task may declare at most 50 unique entries. The Daemon reserves those paths
  before filling the 200-path Spec Context Bundle with prior changed files.

## Decomposition rules

- **Vertical slices.** Each task delivers a narrow but complete path through every layer it touches, demoable or verifiable on its own — a tracer bullet, not a layer ("all the schemas" is a wrong task; "expired imports retry and surface their status" is a right one).
- **Prefactoring first.** When a slice needs the ground prepared, make that its own leading task: make the change easy, then make the easy change.
- **Sized for one fresh session.** A task an agent can complete in a single sitting with a fresh context. More than ~7 subtasks or files means split it.
- **Tests embedded, never separated.** Every task's acceptance criteria include its own tests; a trailing "write the tests" task means the earlier tasks were never done.
- **Independently implementable.** Once its `needs` are completed, a task must require no other unfinished work — that's what allows parallel execution across worktrees later.
- **Verification must be hermetic, portable, effect-proving, and Daemon-owned.** Every task's `## Verification` commands must be satisfiable in a fresh worktree using only repository state, declared config, and task-owned setup. Do not depend on untracked local files, prior Runs, interactive prompts, pushed branches, or ambient machine state unless the task explicitly creates that state. Use portable shell forms: prefer `grep` over `rg` in task gates, avoid `wc`-pipeline assertions that vary across platforms, and use repository build flags such as `go build -buildvcs=false ./...` when a build is required. Verification must prove the Task's effect with executable checks, not only prove that unrelated suites still pass. The Daemon runs these commands after the Agent turn and may send one failure-only Verification Feedback prompt; do not tell the Agent to run the authoritative gate itself.
- **Commit and push stay out of task criteria.** The Daemon owns Task commits, Run integration, and any configured push. Never put commit, push, PR creation, or branch-publishing requirements in task Requirements, Subtasks, Acceptance Criteria, or Verification commands.

### Task Type selection

Choose the type from the Task's dominant delivered outcome, not from who will
implement it, the tool used, or configuration choices:

- `backend` — CLI, daemon, service, API, domain, or core application behavior;
- `frontend` — browser UI, TUI, visual, interaction, accessibility, or UX behavior;
- `data` — persistence schema, query, migration, import, or data-pipeline behavior;
- `infra` — build, CI/CD, packaging, deployment, environment, or operational infrastructure;
- `docs` — documentation-only behavior or durable knowledge artifacts;
- `test` — test harness, fixture, evaluation, or coverage work without a product behavior change;
- `chore` — bounded maintenance that changes none of the preceding product surfaces.

When a vertical slice crosses types, use the type of its primary user-visible or
operational outcome. If two outcomes are independently valuable or the dominant
outcome remains ambiguous, split the Task or resolve the type during the approval
gate; do not encode multiple values.

## Process

### 1. Derive the breakdown

Start from the TechSpec's Build Order (its dependency statements become graph edges). Map every PRD user story and core feature onto at least one task; an uncovered story is a hole to fix now.

### 2. Approval gate — the one mandatory human check

Present the breakdown as a table (`id | title | type | complexity | needs`) plus a one-line rationale for the slicing, and iterate until the user explicitly approves. Every `type` cell must already contain one allowed Task Type; a blank, placeholder, or combined value blocks approval. Granularity, dependency correctness, and dominant outcome classification are where human judgment pays most — **write no files before approval.**

### 3. Write the files

From the templates in [references/task-template.md](references/task-template.md), write `_tasks.md` and every `task_NN.md` (numbered from `01` in topological order). Copy the approved Task Type into each task's `type` frontmatter and the `_tasks.md` projection table. Acceptance criteria must be independently verifiable — a criterion nobody can check is a wish, not a criterion. Include a `## Verification` section with exact, hermetic, portable commands that prove the task's effect in a fresh worktree; the Daemon runs them verbatim after Agent work. Add `## Context` only for task-specific instruction or interface paths that the standard Spec Context Bundle would not make obvious.

Durability applies here too: describe behavior and interfaces, not repo file paths (relative references within the spec folder are fine — the folder moves as a unit).

### 4. Validate the graph

Before reporting, verify mechanically — parse, don't eyeball:

- Every `graph.nodes[].file` exists and every task file has parseable frontmatter with `status: pending`.
- Every task file has exactly one `type` value from `backend`, `frontend`,
  `data`, `infra`, `docs`, `test`, or `chore`; the `_tasks.md` projection row
  for that task contains the identical value.
- Every `needs` entry names an existing node id; no cycles (a topological order can be printed).
- Every PRD user story appears in some task's References.

Print the wave plan (wave 1 = no needs; wave N = needs met by earlier waves) as the execution preview.

### 5. Report

Reply with the file list and the wave plan.

## Anti-patterns

- Horizontal slicing (a schema task, an API task, a UI task) — every task should cut through the stack.
- Mega-tasks that "do the feature" — if it can't be verified alone, it isn't a task yet.
- Vague criteria ("works correctly", "handles errors") — name the observable behavior.
- Non-hermetic Verification that relies on local-only files, prior Run state, manual prompts, or remote branch state.
- Non-portable Verification such as `wc` pipeline shape checks, `rg` dependency checks where `grep` works, or Go build commands missing the repository's `-buildvcs=false` flag.
- Verification that only exercises broad suites without proving the Task's specific effect.
- Commit, push, PR, or branch-publishing acceptance criteria — those are Daemon and delivery responsibilities, not task success criteria.
- Editing the graph and the task files out of sync — regenerate both from the approved breakdown.
- Writing files before the user approves the breakdown.
- Missing, placeholder, combined, or invented Task Types — classification is a
  required routing contract, not descriptive free text.

## References

- [references/task-template.md](references/task-template.md) — templates for `_tasks.md` and `task_NN.md`. Read it before writing any file.
