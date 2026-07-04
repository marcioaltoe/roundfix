---
name: setup-workflow
description: Configure a repo for the CONTEXT-driven spec workflow — scaffold docs/specs/, docs/adr/, and the CONTEXT.md glossary, wire the knowledge workspace when the repo uses one, and configure triage labels when the repo receives external issues. Run once when preparing a repo for the write-prd/write-tasks/implement pipeline.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.4.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Workflow

Scaffold the per-repo configuration the CONTEXT-driven spec workflow assumes. Local markdown under `docs/specs/` is the **only home of planning artifacts** — there is no external tracker.

- **Spec artifacts** — `docs/specs/<feature-slug>/`, read and written by `write-idea`, `write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, and `archive-spec`
- **Domain docs** — `CONTEXT.md` (glossary) and `docs/adr/`, and the consumer rules for reading them
- **Triage labels** (conditional) — only when the repo receives external/incoming issues on its forge (e.g. a public GitHub repo) that the `triage` skill will process

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `AGENTS.md` and `CLAUDE.md` at the repo root — does either exist? Is there already an `## Agent skills` section? Is `CLAUDE.md` a symlink to `AGENTS.md` (the usual convention)?
- **Knowledge workspace**: are `CONTEXT.md` and `docs` symlinks into `.knowledge/`? If so, the docs tree lives in the central knowledge repository — the scaffolding below happens _through_ the symlinks, and every docs commit follows the `knowledge-workspace` skill (`git -C .knowledge …`), never the code repo. If the symlinks are absent but the project should use the workspace, point the user at `scripts/knowledge-bootstrap.sh` before continuing.
- `CONTEXT.md` / `CONTEXT-MAP.md` at the repo root
- `docs/specs/` and `docs/specs/_archived/` — layout already in place? Any active specs?
- `docs/adr/` and any `src/*/docs/adr/` directories
- `docs/agents/` — does this skill's prior output already exist?
- Legacy planning locations (`.scratch/`, `.compozy/`, `docs/tasks/`, `docs/plans/`) — note them as read-only history; new work goes to `docs/specs/`.
- `git remote -v` — a public forge remote means external issues may arrive (Section C).

### 2. Present findings and ask

Summarise what's present and what's missing. Then walk the user through the decisions **one at a time** — present a section, get the answer, move on (use the AskUserQuestion tool or the CLI's equivalent). Assume the user may not know the terms: open each section with a short explainer.

**Section A — Spec artifacts (the core; a confirmation, not a choice).**

> Explainer: The spec workflow skills coordinate through per-feature folders. Each feature gets `docs/specs/<feature-slug>/` holding `_idea.md` (optional), `_prd.md`, `_techspec.md` (optional), the `_tasks.md` dependency graph, one `task_NN.md` per task, and `qa/` evidence. Dependencies live only in `_tasks.md`; task status lives only in each task file's frontmatter. Shipped specs move to `docs/specs/_archived/` so the active folder shows only live work.

The layout is a fixed convention the skills share — confirm the user wants it scaffolded, and whether any legacy planning artifacts found in step 1 should be flagged as read-only history in the summary.

**Section B — Domain docs.**

> Explainer: The skills read `CONTEXT.md` (the glossary — pure ubiquitous language) before writing specs, code names, or test names, and `docs/adr/` for past decisions. They need to know whether the repo has one global context or multiple (a monorepo with separate bounded contexts).

Confirm the layout:

- **Single-context** — one `CONTEXT.md` + `docs/adr/` at the repo root. Most repos are this.
- **Multi-context** — `CONTEXT-MAP.md` at the root pointing to per-context `CONTEXT.md` files (typically a monorepo).

**Section C — Triage labels (only when the repo receives external issues on its forge).**

> Explainer: When outsiders (users, teammates, bots) file issues on the forge (e.g. GitHub issues of a public repo), the `triage` skill moves them through a state machine using five canonical roles. It needs the label strings this repo actually uses. If nothing external lands on the forge, skip this section entirely — spec tasks carry their own status lifecycle (`pending`/`in_progress`/`completed`/`failed`) and never need triage labels.

The five canonical roles: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. Default: each role's string equals its name; ask only if the repo's labels differ.

### 3. Confirm and edit

Show the user a draft of:

- The scaffold actions for spec artifacts (directories to create, `CONTEXT.md` skeleton if missing)
- The `## Agent skills` block to add to whichever of `CLAUDE.md` / `AGENTS.md` is being edited (see step 4 for selection rules)
- The contents of `docs/agents/issue-tracker.md` and `docs/agents/domain.md` (plus `docs/agents/triage-labels.md` when Section C applies)

Let them edit before writing.

### 4. Write

**Pick the file to edit:**

- If `CLAUDE.md` exists as a real file, edit it. If it is a symlink to `AGENTS.md`, edit `AGENTS.md`.
- Else if `AGENTS.md` exists, edit it.
- If neither exists, ask the user which one to create — don't pick for them.

If an `## Agent skills` block already exists in the chosen file, update it in place rather than appending a duplicate. Don't overwrite user edits to surrounding sections.

The block:

```markdown
## Agent skills

### Issue tracker

Tasks live as local markdown under `docs/specs/<feature-slug>/` (the canonical source — no external tracker). See `docs/agents/issue-tracker.md`.

### Domain docs

[one-line summary of layout — "single-context" or "multi-context"]. See `docs/agents/domain.md`.

### Spec artifacts

Feature specs live under `docs/specs/<feature-slug>/` (`_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`). Dependencies live only in `_tasks.md`; task status lives only in each task file's frontmatter. Shipped specs are archived to `docs/specs/_archived/`.
```

Add a `### Triage labels` line only when Section C applies.

Then scaffold the spec artifacts:

- Create `docs/specs/` and `docs/adr/` if missing (add a `.gitkeep` to empty directories so the layout survives a clone). `docs/specs/_archived/` can wait for the first archive.
- If `CONTEXT.md` is missing, create the glossary skeleton below. Do **not** pre-fill terms — a glossary written by hand during real grilling sessions is worth more than a generated one, and generated context files measurably hurt agent output.

```markdown
# <Project name>

<!-- One or two sentences: what this project is and why it exists. -->

## Language

<!-- One entry per project-specific term, added as each term is resolved during grilling/domain-modeling:
**Term**:
One-sentence definition of what it IS.
_Avoid_: rejected synonyms
-->
```

Then write the docs files using the seed templates in this skill folder as a starting point:

- [issue-tracker-local.md](./issue-tracker-local.md) — the canonical local `docs/specs/` conventions
- [triage-labels.md](./triage-labels.md) — label mapping (Section C only)
- [domain.md](./domain.md) — domain doc consumer rules + layout

**Knowledge workspace note:** when `docs` is a symlink into `.knowledge/`, everything written above lands in the knowledge repository — remind the user (or the committing agent) that these files are committed via `git -C .knowledge …` per the `knowledge-workspace` skill, and the code repo commits only the symlinks and `AGENTS.md`.

### 5. Done

Tell the user the setup is complete and which skills now read from these files: the spec pipeline (`write-idea`, `write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, `archive-spec`) plus `triage` when external issues arrive on the forge. Mention they can edit `docs/agents/*.md` and `CONTEXT.md` directly later — re-running this skill is only necessary to restart from scratch.
