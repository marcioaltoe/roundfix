---
name: setup-workflow
description: Configure a repo for a CONTEXT-driven spec workflow — set up its issue tracker, triage label vocabulary, domain doc layout, and the docs/specs/ artifact scaffold (CONTEXT.md glossary, docs/adr/, spec folders) before using planning, decomposition, or implementation-loop skills. Run once when preparing a repo for PRD, spec, task, triage, or agent workflow skills.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.2.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Workflow

Scaffold the per-repo configuration that a CONTEXT-driven spec workflow assumes:

- **Issue tracker** — where issues live (GitHub by default; local markdown is also supported out of the box)
- **Triage labels** — the strings used for the five canonical triage roles
- **Domain docs** — where `CONTEXT.md` and ADRs live, and the consumer rules for reading them
- **Spec artifacts** — the `docs/specs/<feature-slug>/` layout that `write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, and `archive-spec` read and write

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `git remote -v` and `.git/config` — is this a GitHub repo? Which one?
- `AGENTS.md` and `CLAUDE.md` at the repo root — does either exist? Is there already an `## Agent skills` section in either?
- `CONTEXT.md` and `CONTEXT-MAP.md` at the repo root
- `docs/adr/` and any `src/*/docs/adr/` directories
- `docs/agents/` — does this skill's prior output already exist?
- `docs/specs/` and `docs/specs/_archived/` — is the spec artifact layout already in place? Any active specs?
- `.scratch/` — sign that a local-markdown issue tracker convention is already in use

### 2. Present findings and ask

Summarise what's present and what's missing. Then walk the user through the four decisions **one at a time** — present a section, get the user's answer, then move to the next. Don't dump them all at once.

Assume the user does not know what these terms mean. Each section starts with a short explainer (what it is, why these skills need it, what changes if they pick differently). Then show the choices and the default.

**Section A — Issue tracker.**

> Explainer: The "issue tracker" is where issues live for this repo. Skills like `to-issues`, `triage`, `to-prd`, and `qa` read from and write to it — they need to know whether to call `gh issue create`, write a markdown file under `.scratch/`, or follow some other workflow you describe. Pick the place you actually track work for this repo.

Default posture: PRD-to-issues workflows usually work best when the issue tracker is explicit before planning starts. If a `git remote` points at GitHub, propose that. If a `git remote` points at GitLab (`gitlab.com` or a self-hosted host), propose GitLab. Otherwise (or if the user prefers), offer:

- **GitHub** — issues live in the repo's GitHub Issues (uses the `gh` CLI)
- **GitLab** — issues live in the repo's GitLab Issues (uses the [`glab`](https://gitlab.com/gitlab-org/cli) CLI)
- **Local markdown** — issues live as files under `.scratch/<feature>/` in this repo (good for solo projects or repos without a remote)
- **Linear** — issues live in a Linear workspace, driven through the **Linear MCP** (no CLI). Initiatives = products, projects = deliverables, issues = stories; PRDs map to projects. See [issue-tracker-linear.md](./issue-tracker-linear.md) for the full structure. Linear isn't detectable from `git remote`, so propose it when the user says they track work in Linear.
- **Other** (Jira, etc.) — ask the user to describe the workflow in one paragraph; the skill will record it as freeform prose

**Section B — Triage label vocabulary.**

> Explainer: When the `triage` skill processes an incoming issue, it moves it through a state machine — needs evaluation, waiting on reporter, ready for an AFK agent to pick up, ready for a human, or won't fix. To do that, it needs to apply labels (or the equivalent in your issue tracker) that match strings _you've actually configured_. If your repo already uses different label names (e.g. `bug:triage` instead of `needs-triage`), map them here so the skill applies the right ones instead of creating duplicates.

The five canonical roles:

- `needs-triage` — maintainer needs to evaluate
- `needs-info` — waiting on reporter
- `ready-for-agent` — fully specified, AFK-ready (an agent can pick it up with no human context)
- `ready-for-human` — needs human implementation
- `wontfix` — will not be actioned

Default: each role's string equals its name. Ask the user if they want to override any. If their issue tracker has no existing labels, the defaults are fine.

**Section C — Domain docs.**

> Explainer: Some skills (`improve-codebase-architecture`, `diagnosing-bugs`, `tdd`) read a `CONTEXT.md` file to learn the project's domain language, and `docs/adr/` for past architectural decisions. They need to know whether the repo has one global context or multiple (e.g. a monorepo with separate frontend/backend contexts) so they look in the right place.

Confirm the layout:

- **Single-context** — one `CONTEXT.md` + `docs/adr/` at the repo root. Most repos are this.
- **Multi-context** — `CONTEXT-MAP.md` at the root pointing to per-context `CONTEXT.md` files (typically a monorepo).

**Section D — Spec artifacts.**

> Explainer: The spec workflow skills (`write-prd` → `write-techspec` → `write-tasks` → `implement-task`/`implement-spec` → `qa-gate` → `archive-spec`) coordinate through per-feature artifact folders. Each feature gets `docs/specs/<feature-slug>/` holding `_prd.md`, optionally `_techspec.md`, the `_tasks.md` dependency graph, one `task_NN.md` per task, and `qa/` for QA evidence. Shipped specs move to `docs/specs/_archived/<feature-slug>/` so the active folder always shows only live work.

This section is a confirmation, not an open choice — the layout is a fixed convention the skills share. Confirm the user wants it scaffolded, and whether any existing planning artifacts (e.g. `docs/tasks/`, `docs/plans/`) should be noted as legacy locations in the summary.

### 3. Confirm and edit

Show the user a draft of:

- The `## Agent skills` block to add to whichever of `CLAUDE.md` / `AGENTS.md` is being edited (see step 4 for selection rules)
- The contents of `docs/agents/issue-tracker.md`, `docs/agents/triage-labels.md`, `docs/agents/domain.md`
- The scaffold actions for spec artifacts (directories to create, `CONTEXT.md` skeleton if missing)

Let them edit before writing.

### 4. Write

**Pick the file to edit:**

- If `CLAUDE.md` exists, edit it.
- Else if `AGENTS.md` exists, edit it.
- If neither exists, ask the user which one to create — don't pick for them.

Never create `AGENTS.md` when `CLAUDE.md` already exists (or vice versa) — always edit the one that's already there.

If an `## Agent skills` block already exists in the chosen file, update its contents in-place rather than appending a duplicate. Don't overwrite user edits to the surrounding sections.

The block:

```markdown
## Agent skills

### Issue tracker

[one-line summary of where issues are tracked]. See `docs/agents/issue-tracker.md`.

### Triage labels

[one-line summary of the label vocabulary]. See `docs/agents/triage-labels.md`.

### Domain docs

[one-line summary of layout — "single-context" or "multi-context"]. See `docs/agents/domain.md`.

### Spec artifacts

Feature specs live under `docs/specs/<feature-slug>/` (`_prd.md`, `_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`). Dependencies live only in `_tasks.md`; task status lives only in each task file's frontmatter. Shipped specs are archived to `docs/specs/_archived/`.
```

Then write the three docs files using the seed templates in this skill folder as a starting point:

- [issue-tracker-github.md](./issue-tracker-github.md) — GitHub issue tracker
- [issue-tracker-gitlab.md](./issue-tracker-gitlab.md) — GitLab issue tracker
- [issue-tracker-local.md](./issue-tracker-local.md) — local-markdown issue tracker
- [issue-tracker-linear.md](./issue-tracker-linear.md) — Linear issue tracker (via Linear MCP)
- [triage-labels.md](./triage-labels.md) — label mapping
- [domain.md](./domain.md) — domain doc consumer rules + layout

For "other" issue trackers, write `docs/agents/issue-tracker.md` from scratch using the user's description.

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

### 5. Done

Tell the user the setup is complete and which skills now read from these files: the spec pipeline (`write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, `archive-spec`) plus the PRD, issue-decomposition, and triage skills. Mention they can edit `docs/agents/*.md` and `CONTEXT.md` directly later — re-running this skill is only necessary to switch issue trackers or restart from scratch.
