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
