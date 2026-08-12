# QA-01 — Governance and protected-tooling provenance

Status: pass

The Task Graph names `task_06` as its unique `type: qa` node. It depends on
`task_07`, the only non-QA leaf, whose transitive closure contains every other
Task. Tasks 01, 02, 03, 04, 05, and 07 are `completed`; Task 06 remains
Daemon-owned and `pending` during this gate.

Both the PRD and TechSpec account for identifier strategy, authentication and
HTTP, active ADR obligations, and tooling authority. Their operative source
paths exist under `docs/agents/`. ADR-0080, ADR-0081, ADR-0091, and the cited
ADR-0093 non-applicability are recorded explicitly.

`rtk git diff-tree --no-commit-id --name-status -r <commit>` was run for every
Spec commit from authoring through build `a10638b`:

- `d8c0403` records the QA-gate authorization in the PRD, TechSpec, and
  authorization source before Task 04's protected mutation.
- `f87db11` changes only the authorized `.agents/skills/qa-gate/SKILL.md`, its
  `skills/qa-gate/SKILL.md` mirror, Task 04, and sanctioned ADR-0081 generated
  artifacts. `45ffc72` precedes it and changes only Task 04's instructions.
- `14ddabb` records the corrective Roundfix Skill authorization and Task 07
  before `a10638b` performs the mutation.
- `a10638b` changes only `docs/user-guide/commands.md`, the authorized
  Roundfix Skill pair, Task 07, and sanctioned generated artifacts. It changes
  no Archive Command Go source.

No authorization, prerequisite tooling fix, or consequent tooling fix is
folded into either protected Task commit. No out-of-scope tooling path appears.

A disposable clone at `/private/tmp/qa70-audit-a10638b` reproduced the derived
state:

- `GOCACHE=/private/tmp/qa70-audit-gocache rtk make baseline-digests` exited 0
  with `changed:false` and reported that derived artifacts already matched
  canonical sources.
- The plan-characterization re-record command passed 7 tests.
- The catalog-diagnostic re-record command passed 2 tests.
- `rtk git -c core.fsmonitor=false status --short` was empty afterward.
- Both the QA Gate Skill mirror and Roundfix Skill mirror compared byte-equal.
