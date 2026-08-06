# Constraint and ancestry audit

Build: `9252430f9e6c63332775a90ee9dcb08f7bbccef7`.

## Project Constraints

- Identifier strategy: accounted for in both active artifacts. The PRD records
  no new project-owned Internal Identifier; the TechSpec records the applicable
  stable `SC-*` diagnostic convention. Both cite `docs/agents/domain.md`.
- Authentication and HTTP: not applicable with a reason in both artifacts;
  the feature reads local Spec artifacts and opens no transport. Both cite
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: both artifacts account for accepted ADR-0080,
  ADR-0081, ADR-0091, and ADR-0093 with applicability or relation-only reasons
  and cite `docs/agents/domain.md`.
- Tooling authority: both artifacts name the exact protected Skill trees and
  mirrors. The TechSpec also lists the Baseline module asset and deterministic
  digest paths. The operative records are
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md` and
  `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.

## Authored gate and Task state

`_tasks.md` names `qa: task_06`. `task_06` is `type: qa`, is terminal, and now
depends on `task_07`, the only non-QA leaf. Tasks 01, 02, 03, 04, 05, and 07
all carry `status: completed`. The gate rerun therefore starts after every
non-QA dependency settled. The PRD contains no `## Unreachable Acceptance`
declaration.

## Chronology

Fresh `git merge-base --is-ancestor` checks all exited 0:

- authorization commits `2e560cea` and `9bdaedbe` precede task_01 commit
  `bd38544d` and task_05 commit `a51a94cb`;
- corrective authoring commit `b4347ac6` precedes corrective task_07 commit
  `d603031e`;
- task_07 precedes the evidence-only correction `9252430f`;
- dependency order is preserved: task_01 precedes task_04; task_02 precedes
  task_03; tasks 03 and 04 precede task_05.

The prior F-001 correction made the active PRD path boundary exact after the
first failed gate, but it did not create or backdate authorization: both
maintainer authorization records already existed before either protected
Skill Task commit.

## Changed-path audit

Fresh `git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r`
inventories were read for every Spec implementation, corrective, and authoring
commit. The protected changes are bounded as follows:

- `bd38544d` changes the authorized Baseline module carrier, the guide, its
  Task file, and deterministic Baseline digest/characterization fallout; no Go
  source changes.
- `a51a94cb` changes only the authorized `write-tasks` and `roundfix` Skill
  trees and mirrors, its Task file, and deterministic Baseline fallout; no Go
  source changes.
- `d603031e` changes the same authorized Skill trees and mirrors, the order
  carriers, its Task file, and deterministic Baseline fallout; no Go source
  changes.
- `9252430f` changes only active `task_01.md` and `task_07.md` evidence and
  Verification scope.

The implementation commits `acb3dbc9`, `db1b555e`, and `7eb63c65` change
product Go code, tests/fixtures, and their assigned Task files; they do not
change repository-tooling configuration. No prerequisite or consequent fix is
folded into a protected tooling Task commit. The deterministic fallout is
attributed to `make baseline-digests` and the two named characterization
re-record commands in the Task Results, as sanctioned by ADR-0081.

Result: R01 and R02 pass.
