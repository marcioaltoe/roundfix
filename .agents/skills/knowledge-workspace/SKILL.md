---
name: knowledge-workspace
description: Operate the external knowledge workspace — the central knowledge repository mounted at .knowledge/ and exposed through versioned symlinks (CONTEXT.md, docs/, legacy .compozy/.scratch). Use when bootstrapping the workspace, when CONTEXT.md or docs/ are symlinks into .knowledge, and before committing ANY documentation change (specs, ADRs, glossary, task-file statuses) — those commits happen inside .knowledge, never in the code repository.
metadata:
  category: setup
  tags: [workflow, documentation, git, repository-context]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Knowledge Workspace

Long-lived documentation lives in one central repository (`gesttione-solutions/knowledge`, `projects/<project>/`), not in each code repository. Code repos mount their slice as a sparse checkout at `.knowledge/` (gitignored) and expose it through **versioned symlinks**. The payoff: docs history in one place, code history clean, and every agent session sees the same spec state.

## The layout contract

| Path in the code repo            | What it is                                                                              |
| -------------------------------- | --------------------------------------------------------------------------------------- |
| `.knowledge/`                    | Sparse git checkout of the knowledge repo — **gitignored, never committed**             |
| `CONTEXT.md`                     | Symlink → `.knowledge/projects/<project>/CONTEXT.md`                                    |
| `docs`                           | Symlink → `.knowledge/projects/<project>/docs` (includes `docs/specs/` and `docs/adr/`) |
| `.compozy`, `.scratch`           | Legacy symlinks (`.scratch` usually local-only via `.git/info/exclude`)                 |
| `scripts/knowledge-bootstrap.sh` | Symlink → the bootstrap script inside the workspace                                     |

The code repository versions the symlinks themselves (plus `AGENTS.md`/`DESIGN.md` as real files). Never replace a symlink with a real file, and never point it anywhere else.

## Bootstrap

When `.knowledge/` or the symlinks are missing:

```bash
scripts/knowledge-bootstrap.sh                # new repo, no local docs yet
scripts/knowledge-bootstrap.sh --adopt-local  # existing CONTEXT.md/docs get MOVED into the workspace
KNOWLEDGE_PROJECT=<name> scripts/knowledge-bootstrap.sh   # when the folder name != project name
```

`make docs-bootstrap` runs the same script where the target exists. `--adopt-local` moves existing local docs into `.knowledge` — if the destination already has content, the script stops and asks for a manual merge; do that merge deliberately, never by overwriting. After bootstrap, confirm `.knowledge/` is in `.gitignore`.

## The two-repository commit flow

Documentation and code have different homes, so one piece of work often produces **two commits**:

1. **Code changes** → committed in the code repository, as always.
2. **Documentation changes** — anything under the symlinks: `CONTEXT.md` glossary entries, `docs/adr/`, `docs/specs/` artifacts (`_prd.md`, `_techspec.md`, `_tasks.md`, `task_NN.md` **including status flips**), QA reports — → committed inside the workspace:

```bash
git -C .knowledge status
git -C .knowledge add projects/<project>
git -C .knowledge commit -m "docs: <what changed>"
git -C .knowledge push origin main
```

Consequences for the spec workflow:

- `implement-task`'s per-task commit is a **code** commit; the task file's `status: completed` flip and `## Result` section are **docs** — commit and push them in `.knowledge` in the same breath. An unpushed status flip is invisible to every other session and to any scheduler reading the graph.
- `archive-spec`'s move happens inside the workspace: `git -C .knowledge mv projects/<project>/docs/specs/<slug> projects/<project>/docs/specs/_archived/<slug>`.
- The knowledge repo pushes to `main` directly — a failed push (non-fast-forward) means another session moved first: `git -C .knowledge pull --ff-only` and re-apply, never force-push.

## Anti-patterns

- Committing `.knowledge/` — or any docs content — in the code repository.
- Replacing a symlink with a real file "to make the edit stick"; the edit already sticks — through the symlink into the workspace.
- Flipping a task status without pushing the knowledge commit — the coordination state exists only when pushed.
- Mixing the docs commit into the code commit message ("feat: X + update docs") — two repos, two commits, each with its own message.
- Running `--adopt-local` casually when `.knowledge` already has content for the project — merge manually as the script demands.
