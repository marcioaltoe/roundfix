<!-- setup-context-driven:begin id=guide.spec-routing version=0.0.1 -->

# Spec routing

- **mandatory**: Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.

- **mandatory**: Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks`; skip the TechSpec only when the feature has no architectural surface.

- **mandatory**: Refactor or bug fix without product-behavior change: run `write-techspec` → `write-tasks` and create the minimal `_prd.md` required by the downstream artifact contract.

- **mandatory**: Trivial one-line fix, typo, or configuration tweak: implement directly without a Spec only when intent, acceptance criteria, and Verification are obvious.

- **mandatory**: Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, and execute implementation from the Task Graph.

- **mandatory**: For each Task, run the selected incremental Verification named in `docs/agents/agent-instructions.md` to answer whether the current slice remains valid before handoff. CI must run the selected repository Verification from a fresh run to answer whether the assembled tree satisfies the repository contract. A missing incremental selection is a Baseline decision to answer, never a license to skip the local tier or a waiver to repeat in each Spec.

- **mandatory**: Before producing a Task Graph, require every active, non-archived, and not already completed Spec PRD and present TechSpec to contain complete Project Constraints: applicability with reasons for identifier strategy, authentication and HTTP, active ADR obligations, and tooling authority, each citing its operative `docs/agents/` source.

- **mandatory**: Refuse a tooling Task unless the active PRD and present TechSpec record express maintainer authorization and the exact bounded repository-relative files; Task assignment, setup approval, or generic implementation approval is not authorization. The operative record lives at `<spec-root>/<slug>/_authorization.md` and carries the approval state (`status`), maintainer decision date (`granted`), permitted actions (`action`), the closed `operations` vocabulary (`implement`, `commit`, `push`, `pull_request`, `merge`, and `release`), exact bounded repository paths (`paths`), the consuming Spec (`consuming`), and any sanctioned regeneration. An absent `operations` list grants no operation. A Spec that changes no Governed Path records `paths: []`, which grants its listed operations and bounds no path; a record without a `paths` key is refused.

- **mandatory**: An authorized tooling Task may change a Governed Path, that is a protected tooling path, only when its grant bounds that path; stop before any other governed mutation, because the Task fails when its commit changes a Governed Path outside the grant. An ordinary path a Task changes without declaring it is not refused: the Daemon records it in the Task file under `## Recorded paths`, and the QA gate discloses it. A new or widened grant must land in the delivery target's ancestry before the consuming squash delivery; a separate grant commit inside that consuming pull request is insufficient.

- **mandatory**: Final QA verifies Project Constraint applicability, operative source paths, tooling authorization, and actual changed-file scope from Git evidence; missing authorization, untraceable scope, or out-of-scope tooling changes fails the gate.

- **mandatory**: Keep completed or archived legacy Specs byte-identical. Dependencies remain owned only by the Task Graph, and status remains owned only by each Task file.

- **mandatory**: Rest a Spec's acceptance, in at least one named row, on evidence originating outside the Spec's own artifacts: a repository the Spec did not build, a measurement it did not design, or published literature. Record in that row where the evidence came from, so a later reader can tell it apart from a rehearsal of the Spec's own premise. When the outside source cannot be obtained during authoring, record the row as blocked with that reason and continue: decomposition never stalls and never asks a person. The QA gate then holds Pull Request preparation until the row is satisfied or carried forward on declared unmoved evidence.

<!-- setup-context-driven:end id=guide.spec-routing -->
