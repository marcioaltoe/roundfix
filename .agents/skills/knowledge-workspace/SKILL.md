---
name: knowledge-workspace
description: Work with project documentation and the central knowledge repository (marcioaltoe/secondbrain). Use when editing docs/, CONTEXT.md, specs, ADRs, or task files in any repo with a .secondbrain-project file; when deciding where a documentation change should be committed; when touching secondbrain's mirror, notes, shared or wiki areas; and when research, planning or implementation work could consult or feed the second brain (wiki + qmd search).
metadata:
  category: setup
  tags: [workflow, documentation, git, repository-context, second-brain]
  version: 0.6.2
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
---

# Knowledge Workspace (Secondbrain)

The central knowledge repository is **marcioaltoe/secondbrain**: per-project mirrors, editorial notes, shared standards, and an **LLM-maintained knowledge wiki** (the second brain). Code repositories own their documentation as regular files; Secondbrain aggregates it automatically.

## The contract

Documentation (`CONTEXT.md`, `docs/` — specs, ADRs, task files, review artifacts) consists of **regular files in the code repository**. There is ONE commit flow:

- Edit docs and code together and commit together, in the same branch and PR. `task_NN.md` status flips and `## Result` sections are ordinary repo commits — no second repository, no separate docs push.
- After every push to `main` (any merge type or direct push; never feature branches), GitHub Actions mirrors the paths listed in `.secondbrain-export` into `secondbrain/projects/<project>/mirror/`. This is automatic — **never sync manually**.
- The mirror is **read-only**. Never edit `projects/<project>/mirror/` in secondbrain; the next sync overwrites it. Content worth editing belongs either in the code repo (project docs) or in secondbrain's editorial areas.
- Editorial areas are edited directly in secondbrain and pushed to its `main`: `projects/<project>/notes/` (project notes that don't belong to the repo — meeting notes, era archives, client feedback) and `shared/` (cross-project standards). `registry/projects.json` whitelists which repository may publish to which project.
- Consumption is one-way: nothing in secondbrain is ever written back to code repositories. To reference `shared/` material from a repo, link to it — do not copy it in.

## The export contract

| File in the code repo                    | Purpose                                                                                                                                                                                               |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `.secondbrain-project`                   | Exactly the project name — the key under `projects/` and in the secondbrain registry                                                                                                                  |
| `.secondbrain-export`                    | What the repo publishes: one path per line, directories recursive, `!pattern` lines are excludes (rsync/gitignore-like), `#` comments                                                                 |
| `.github/workflows/secondbrain-sync.yml` | Standalone sync workflow, copied from secondbrain's `templates/repo/secondbrain-sync.yml` — the same file for every repo and owner, because a private personal repo cannot expose a reusable workflow |

Only committed files can be exported (the sync reads a checkout of the merged commit), and symlinks are rejected. When adding a new exported artifact (e.g. an OpenAPI file), update `.secondbrain-export` in the same PR.

## Adding a new project

Follow `templates/repo/` and the README in secondbrain: register in `registry/projects.json` (`active: true`), create `projects/<name>/notes/`, add the three files above to the repo, grant it access to the `SECONDBRAIN_APP_CLIENT_ID` variable and `SECONDBRAIN_APP_PRIVATE_KEY` secret (org-level for org repos, repo-level for personal-account repos), merge to `main`.

## The second brain (wiki)

On top of the record, secondbrain hosts a Karpathy-style LLM wiki. The authoritative contract is **secondbrain's own `AGENTS.md`** (auto-loaded by sessions in `~/dev/secondbrain`); this section covers what any other session needs. When the `secondbrain` skill (from secondbrain's `brain/skills/`) is installed, it is the detailed operating guide for consulting, ingesting, and messaging; this skill owns the documentation and mirror contract.

- **Consulting** (research for implementation or planning, from any repo): start with `~/dev/secondbrain/wiki/index.md`, or search everything — wiki, ingested sources, shared standards, and all project mirrors — with qmd:

  ```bash
  qmd query "pergunta em linguagem natural"          # hybrid + rerank, best
  qmd search "keywords" -c wiki                      # scoped: wiki|raw|shared|projects|areas|fiscal
  qmd query "..." --all --files --min-score 0.3      # agent-friendly output
  qmd query "..." -c leis                            # statutory text; opt-in, never in unscoped searches
  ```

- **Feeding business knowledge**: load `ai-memory-messaging` and send research to `default/marcioaltoe-secondbrain` with `memory_message_send`. Include source provenance, dates, findings and limits. The filesystem inbox was retired by Marcio on 2026-10-25; outside sessions do not write to Secondbrain. Research is archived and compiled there through Ingest. A qmd duplicate check is useful before sending.
- **Project coordination and continuity**: use `memory_message_*` for requests and responses between existing projects, and `ai-memory-handoff` for same-project continuity on request. Static clients name sender and recipient scopes. If a mailbox is absent, tell Marcio; do not create a filesystem substitute. Update stale local rules/automation and check the Roundfix Baseline when it still prescribes the old workflow.
- **Territories**: `wiki/`, `raw/`, editorial notes, shared material and areas are maintained by Secondbrain sessions. Mirrors remain read-only; there is no outside-session write carve-out.

## Anti-patterns

- Creating symlinks for `CONTEXT.md`/`docs`, mounting a `.knowledge/` checkout, or running a bootstrap script — that flow is retired everywhere.
- Editing or committing anything under `projects/<project>/mirror/` in secondbrain.
- Splitting a docs change out of the PR that motivated it "to sync faster" — the mirror updates on merge; in-flight state living only in the branch is correct.
- Copying `shared/` content into a code repository — link to it instead; copies drift.
- Pointing `.secondbrain-project` at a different project's folder.
