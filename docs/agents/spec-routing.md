# Spec workflow routing

How to work inside this repo's CONTEXT-driven spec workflow: which pipeline stages to run for a
change, and what marks the work done. The canonical working model is local markdown under
`docs/specs/`; artifact locations and conventions live in `docs/agents/issue-tracker.md`.

## Pipeline

```text
write-idea -> write-prd -> write-techspec -> write-tasks -> implement-spec / implement-task -> qa-gate -> roundfix archive
```

Every stage reads and writes `docs/specs/<slug>/`. Downstream stages parse the artifacts, not the
conversation, so a fresh session must be able to continue from the files alone.

## Entry points

| Change | Route |
| --- | --- |
| Large or fuzzy product initiative: new product area, multi-feature epic, open solution shape | `write-idea` -> `write-prd` -> `write-techspec` -> `write-tasks` |
| Standard feature: clear scope, changes product behavior | `write-prd` -> `write-techspec` -> `write-tasks` |
| Refactor or bug fix: no product behavior change | `write-techspec` -> `write-tasks` |
| Trivial change: one-line fix, typo, config tweak | Direct implementation; no spec folder |

Notes that keep the routes honest:

- `brainstorming` precedes creative or feature work and routes the outcome to the right entry
  point. For trivial changes it says so explicitly.
- On the refactor or bug-fix route, `write-techspec` creates the numbered spec folder and a
  minimal `_prd.md` so `write-tasks`, `qa-gate`, and `archive-spec` keep one artifact contract.
- A tech spec can be skipped only when the feature has no real architectural surface.
  `write-tasks` calls that out and compensates with deeper exploration.
- When in doubt between two tiers, start with the smaller one. A route upgrades by adding the
  missing upstream artifact when product questions appear.
- Every route converges on `write-tasks`: implementation executes from the Task Graph, not from
  an ad-hoc plan.

## Done

- `implement-spec` or `implement-task` drives every Task to `completed`, each with fresh
  verification evidence.
- After the last Task, `qa-gate` validates the assembled feature against the Spec's user stories
  and acceptance criteria, writing evidence to `docs/specs/<slug>/qa/`.
- On QA pass, the Archive Command (`roundfix archive <slug>`) stamps archive
  metadata and moves the folder to `docs/specs/_archived/`. Merge and release
  are separate, user-driven actions.
