# Spec Workflow Routing

How to work inside this repo's CONTEXT-driven spec workflow: which pipeline stages to run for a given change, and what marks the work done. The canonical working model is local markdown under `docs/specs/` — artifact locations and conventions live in `docs/agents/issue-tracker.md`.

## The pipeline

```text
write-idea → write-prd → write-techspec → write-tasks → implement-spec / implement-task → qa-gate → archive-spec
```

Every stage reads and writes `docs/specs/<slug>/`. Downstream stages parse the artifacts, not the conversation — a fresh session must be able to continue from the files alone.

## Pick the entry point by the change

| Change                                                                                            | Route                                                                             |
| ------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| **Large or fuzzy product initiative** — new product area, multi-feature epic, open solution shape | `write-idea` → `write-prd` → `write-techspec` → `write-tasks`                     |
| **Standard feature** — clear scope, changes product behavior                                      | `write-prd` → `write-techspec` → `write-tasks`                                    |
| **Refactor or bug fix** — no product behavior change                                              | `write-techspec` → `write-tasks`                                                  |
| **Trivial change** — one-line fix, typo, config tweak                                             | Direct implementation (`systematic-debugging` + `no-workarounds`), no spec folder |

Notes that keep the routes honest:

- `brainstorming` precedes creative/feature work and routes the outcome to the right entry point; for trivial changes it says so explicitly.
- On the refactor/bugfix route, `write-techspec` mints the numbered spec folder and a **minimal `_prd.md`** (frontmatter, problem statement, goals, core features, non-goals — engineering-framed, no product interview) so `write-tasks`, `qa-gate`, and `archive-spec` keep a single artifact contract.
- A techspec is skippable only when the feature has no real architectural surface — `write-tasks` calls that out and compensates with deeper exploration.
- When in doubt between two tiers, start with the smaller one — a route upgrades cleanly by adding the missing upstream artifact when product questions appear.
- Every route converges on `write-tasks`: implementation always executes from the task graph, never from an ad-hoc plan.

## What "done" means

- `implement-spec` / `implement-task` drive every task to `completed`, each with its own fresh verification evidence.
- After the last task, `qa-gate` validates the assembled feature against the spec's user stories and acceptance criteria on the running app, writing evidence to `docs/specs/<slug>/qa/`.
- On QA pass the spec is **complete**: `archive-spec` runs automatically — it stamps the frontmatter and moves the folder to `docs/specs/_archived/` — and then suggests the publish step (PR). Merge/release is a separate, user-driven action, never a gate for archiving.
