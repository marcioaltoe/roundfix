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

<!-- setup-context-driven:begin id=guide.spec-routing version=0.0.1 -->

# Spec routing

- **mandatory**: Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.

- **mandatory**: Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks`; skip the TechSpec only when the feature has no architectural surface.

- **mandatory**: Refactor or bug fix without product-behavior change: run `write-techspec` → `write-tasks` and create the minimal `_prd.md` required by the downstream artifact contract.

- **mandatory**: Trivial one-line fix, typo, or configuration tweak: implement directly without a Spec only when intent, acceptance criteria, and Verification are obvious.

- **mandatory**: Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, and execute implementation from the Task Graph.

- **mandatory**: Before producing a Task Graph, require every active, non-archived, and not already completed Spec PRD and present TechSpec to contain complete Project Constraints: applicability with reasons for identifier strategy, authentication and HTTP, active ADR obligations, and tooling authority, each citing its operative `docs/agents/` source.

- **mandatory**: Refuse a tooling Task unless the active PRD and present TechSpec record express maintainer authorization and the exact bounded repository-relative files; Task assignment, setup approval, or generic implementation approval is not authorization.

- **mandatory**: An authorized tooling Task may mutate only its bounded repository-relative files and its own Task file; stop before any other mutation and fail the Task if changed-file postflight finds another path.

- **mandatory**: Final QA verifies Project Constraint applicability, operative source paths, tooling authorization, and actual changed-file scope from Git evidence; missing authorization, untraceable scope, or out-of-scope tooling changes fails the gate.

- **mandatory**: Keep completed or archived legacy Specs byte-identical. Dependencies remain owned only by the Task Graph, and status remains owned only by each Task file.

<!-- setup-context-driven:end id=guide.spec-routing -->

<!-- roundfix:repository-rule:begin id=rule.4a17d217e1a8c5c732ab54bb64f5b6ce3eafeaaf428fa6da80fe46ac48276555 -->
### Spec routing

Pick the pipeline entry point by the change — large initiative, feature,
refactor/bugfix, or trivial. See `docs/agents/spec-routing.md`.


<!-- roundfix:repository-rule:end id=rule.4a17d217e1a8c5c732ab54bb64f5b6ce3eafeaaf428fa6da80fe46ac48276555 -->

<!-- roundfix:repository-rule:begin id=rule.068b6fb73d68a1c56aa33b0eb68f6732f2a4b84851b29b279dc0912117c18bf9 -->
4. Asking for confirmation before running spec tasks — invocation is the
   authorization

<!-- roundfix:repository-rule:end id=rule.068b6fb73d68a1c56aa33b0eb68f6732f2a4b84851b29b279dc0912117c18bf9 -->
