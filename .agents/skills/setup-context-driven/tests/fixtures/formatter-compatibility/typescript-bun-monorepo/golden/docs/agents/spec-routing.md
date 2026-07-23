<!-- setup-context-driven:begin id=guide.spec-routing version=0.0.1 -->

# Spec routing

- **mandatory**: Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.

- **mandatory**: Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks`; skip the TechSpec only when the feature has no architectural surface.

- **mandatory**: Refactor or bug fix without product-behavior change: run `write-techspec` → `write-tasks` and create the minimal `_prd.md` required by the downstream artifact contract.

- **mandatory**: Trivial one-line fix, typo, or configuration tweak: implement directly without a Spec only when intent, acceptance criteria, and Verification are obvious.

- **mandatory**: Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, and execute implementation from the Task Graph.

<!-- setup-context-driven:end id=guide.spec-routing -->
