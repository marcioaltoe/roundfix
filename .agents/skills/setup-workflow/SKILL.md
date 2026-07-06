---
name: setup-workflow
description: Configure a repo for the CONTEXT-driven spec workflow — scaffold docs/specs/, docs/adr/, and the CONTEXT.md glossary, and seed the docs/agents/ usage guides (issue tracker, spec routing, domain docs, triage labels, autonomous work model). Run when preparing a repo for the write-prd/write-tasks/implement pipeline or for Roundfix-driven autonomous work; re-run to refresh — it overwrites the skill-owned docs/agents/ files and prunes deprecated content.
disable-model-invocation: true
metadata:
  category: setup
  tags: [workflow, prd, issues, planning, triage, repository-context, agents]
  version: 0.6.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Setup Workflow

Scaffold the per-repo configuration the CONTEXT-driven spec workflow assumes. Local markdown under `docs/specs/` is the **only home of planning artifacts** — there is no external tracker.

- **Spec artifacts** — `docs/specs/<feature-slug>/`, read and written by `write-idea`, `write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, and `archive-spec`
- **Spec routing** — how an agent picks the pipeline entry point for a given change (large initiative / feature / refactor-bugfix / trivial) and what marks a spec done
- **Domain docs** — `CONTEXT.md` (glossary) and `docs/adr/`, and the consumer rules for reading them
- **Triage labels** (conditional) — only when the repo receives external/incoming issues on its forge (e.g. a public GitHub repo) that the `triage` skill will process
- **Autonomous work model** (conditional) — only when the repo delegates implementation to agent runtimes (e.g. through Roundfix): the orchestrator/implementer split (Fable orchestrates and authors Specs; implementation goes to an ACP Runtime), runtime routing, and the hard rule that makes the split binding

These usage rules are seeded into the repo as `docs/agents/*.md` — the canonical, always-current explanation of how agents work inside the CONTEXT-driven workflow. This skill owns those files: a re-run regenerates them.

This is a prompt-driven skill, not a deterministic script. Explore, present what you found, confirm with the user, then write.

## Process

### 1. Explore

Look at the current repo to understand its starting state. Read whatever exists; don't assume:

- `AGENTS.md` and `CLAUDE.md` at the repo root — does either exist? Is there already an `## Agent skills` section? Is `CLAUDE.md` a symlink to `AGENTS.md` (the usual convention)?
- `CONTEXT.md` / `CONTEXT-MAP.md` at the repo root
- `docs/specs/` and `docs/specs/_archived/` — layout already in place? Any active specs?
- `docs/adr/` and any `src/*/docs/adr/` directories
- `docs/agents/` — does this skill's prior output already exist?
- Legacy planning locations (`.scratch/`, `.compozy/`, `docs/tasks/`, `docs/plans/`) — note them as read-only history; new work goes to `docs/specs/`.
- `git remote -v` — a public forge remote means external issues may arrive (Section C).

### 2. Present findings and ask

Summarise what's present and what's missing. Then walk the user through the decisions **one at a time** — present a section, get the answer, move on (use the AskUserQuestion tool or the CLI's equivalent). Assume the user may not know the terms: open each section with a short explainer.

**Section A — Spec artifacts (the core; a confirmation, not a choice).**

> Explainer: The spec workflow skills coordinate through per-feature folders. Each feature gets `docs/specs/<feature-slug>/` holding `_idea.md` (optional), `_prd.md`, `_techspec.md` (optional), the `_tasks.md` dependency graph, one `task_NN.md` per task, and `qa/` evidence. Dependencies live only in `_tasks.md`; task status lives only in each task file's frontmatter. Completed specs (all tasks done, QA passed) move to `docs/specs/_archived/` so the active folder shows only live work.

The layout is a fixed convention the skills share — confirm the user wants it scaffolded, and whether any legacy planning artifacts found in step 1 should be flagged as read-only history in the summary.

**Section B — Domain docs.**

> Explainer: The skills read `CONTEXT.md` (the glossary — pure ubiquitous language) before writing specs, code names, or test names, and `docs/adr/` for past decisions. They need to know whether the repo has one global context or multiple (a monorepo with separate bounded contexts).

Confirm the layout:

- **Single-context** — one `CONTEXT.md` + `docs/adr/` at the repo root. Most repos are this.
- **Multi-context** — `CONTEXT-MAP.md` at the root pointing to per-context `CONTEXT.md` files (typically a monorepo).

**Section C — Triage labels (only when the repo receives external issues on its forge).**

> Explainer: When outsiders (users, teammates, bots) file issues on the forge (e.g. GitHub issues of a public repo), the `triage` skill moves them through a state machine using five canonical roles. It needs the label strings this repo actually uses. If nothing external lands on the forge, skip this section entirely — spec tasks carry their own status lifecycle (`pending`/`in_progress`/`completed`/`failed`) and never need triage labels.

The five canonical roles: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. Default: each role's string equals its name; ask only if the repo's labels differ.

**Section D — Autonomous work model (only when the repo delegates implementation to agent runtimes).**

> Explainer: When a supervising Claude Code session (Fable) drives this repo autonomously — typically through Roundfix — the roles split hard: Fable orchestrates and authors Specs; implementation is delegated to an ACP Runtime. Fable's usage budget is reserved for judgment; operational work goes to the other models. The split binds every Fable-powered session, interactive or autonomous.

Confirm, one at a time:

- Whether the repo works this way at all — if not, skip the section and don't seed the file.
- The default implementer (default: Codex `gpt-5.5` at `xhigh`, model selection delegated down the acpx → codex chain).
- The design implementer and its scope (default: Claude Code with Opus 4.8 at `high`/`xhigh` for design, UI, UX, and frontend Tasks) — and what this repo's design surface actually is (TUI, web frontend, both), which fills the placeholder in the seed.
- The repo's verification gate name, so the seed's runtime-independence section can name it.

### 3. Confirm and edit

Show the user a draft of:

- The scaffold actions for spec artifacts (directories to create, `CONTEXT.md` skeleton if missing)
- The `## Agent skills` block to add to whichever of `CLAUDE.md` / `AGENTS.md` is being edited (see step 4 for selection rules)
- The contents of `docs/agents/issue-tracker.md`, `docs/agents/spec-routing.md`, and `docs/agents/domain.md` (plus `docs/agents/triage-labels.md` when Section C applies, and `docs/agents/autonomous-work.md` when Section D applies)

Let them edit before writing. On a re-run, existing `docs/agents/` files found in step 1 are inputs to the draft, not something to preserve verbatim: carry forward repo-specific answers (like custom label strings), regenerate the rest from the current seeds.

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

Feature specs live under `docs/specs/<feature-slug>/` (`_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, `task_NN.md`, `qa/`). Dependencies live only in `_tasks.md`; task status lives only in each task file's frontmatter. Completed specs (all tasks done, QA passed) are archived to `docs/specs/_archived/`.

### Spec routing

Pick the pipeline entry point by the change — large initiative, feature, refactor/bugfix, or trivial. See `docs/agents/spec-routing.md`.
```

Add a `### Triage labels` line only when Section C applies.

When Section D applies, add this subsection to the block:

```markdown
### Autonomous work

Fable orchestrates and authors Specs; implementation is delegated to an ACP Runtime — Codex (`gpt-5.5` at `xhigh`) by default, Claude Code (Opus 4.8 at `high`/`xhigh`) for design, UI, UX, and frontend Tasks. Binding for every Fable-powered session. See `docs/agents/autonomous-work.md`.
```

and add a one-line hard-rule pointer to the repo's high-priority rules (the section the repo uses for MUST-level rules), for example:

```markdown
- **HARD RULE — autonomous work model**: binding for every Fable-powered
  session — Fable orchestrates only; implementation is delegated to an ACP
  Runtime per `docs/agents/autonomous-work.md`.
```

Rule bodies live in the seeded doc and in the workflow skills — the agent-instructions file holds only short mandatory pointers, never the full rule text.

**Re-run semantics — this skill owns its seeded `docs/agents/` files and the `## Agent skills` block.** When either already exists, rewrite rather than append: overwrite each seeded `docs/agents/*.md` with the freshly confirmed draft (carrying forward repo-specific answers), delete previously seeded `docs/agents/` files this skill no longer seeds (deprecated guides must not linger and contradict current ones), and regenerate the `## Agent skills` block in place, dropping subsections that no longer apply while carrying forward repo-authored subsections that point at `docs/agents/` files this skill does not seed. Ownership covers only the block and the files seeded from this skill's templates — repo-authored `docs/agents/` files it never seeded and user content elsewhere in `AGENTS.md`/`CLAUDE.md` stay untouched.

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
- [spec-routing.md](./spec-routing.md) — pipeline entry-point routing and the definition of done
- [triage-labels.md](./triage-labels.md) — label mapping (Section C only)
- [domain.md](./domain.md) — domain doc consumer rules + layout
- [autonomous-work.md](./autonomous-work.md) — orchestrator/implementer split, runtime routing, and Spec-authoring behaviors (Section D only; fill the design-surface and verification-gate placeholders with the confirmed answers)

### 5. Done

Tell the user the setup is complete and which skills now read from these files: the spec pipeline (`write-idea`, `write-prd`, `write-techspec`, `write-tasks`, `implement-task`, `implement-spec`, `qa-gate`, `archive-spec`), `triage` when external issues arrive on the forge, and — when Section D applies — every Fable-powered orchestrator session, which is bound by `docs/agents/autonomous-work.md`. Mention they can edit `docs/agents/*.md` and `CONTEXT.md` directly later, but that a re-run of this skill regenerates the seeded `docs/agents/` files and the `## Agent skills` block from the current templates — durable customizations belong in the confirmation answers, `CONTEXT.md`, or sections outside the owned block.
