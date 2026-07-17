# Spec workflow routing

The CONTEXT-driven pipeline coordinates through local markdown under
`docs/specs/`.

## Pipeline

```text
write-idea → write-prd → write-techspec → write-tasks →
implement-spec / implement-task → qa-gate → archive-spec
```

Every stage reads and writes `docs/specs/<slug>/`. A fresh session must be able
to continue from the files without relying on prior conversation.

## Entry points

| Change | Route |
| --- | --- |
| Large or fuzzy product initiative | `write-idea` → `write-prd` → `write-techspec` → `write-tasks` |
| Standard feature that changes product behavior | `write-prd` → `write-techspec` → `write-tasks` |
| Refactor or bug fix without product-behavior change | `write-techspec` → `write-tasks` |
| Trivial one-line fix, typo, or configuration tweak | Direct implementation; no Spec folder |

- `brainstorming` precedes creative or feature work and selects the route.
- The refactor/bugfix route creates a minimal `_prd.md` so every downstream
  skill retains one artifact contract.
- A tech spec can be skipped only when the feature has no architectural
  surface.
- When two tiers appear sufficient, begin with the smaller route and add an
  upstream artifact if product questions emerge.
- Implementation always executes from the Task Graph.

## Done

- Every Task reaches `completed` with fresh verification evidence.
- `qa-gate` validates the assembled behavior and writes evidence under `qa/`.
- On QA pass, `archive-spec` stamps and moves the Spec to `_archived/`.
- Merge and release remain separate user-directed actions.
