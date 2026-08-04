# QA-01 — Governance and provenance

Status: pass

- `_tasks.md` names `qa: task_06`; `task_06` is `type: qa`, is the only QA
  node, is terminal, and depends on `task_05`, the only non-QA leaf.
  `task_01` through `task_05` are `completed`; `task_06` is `pending`, as
  required while the gate runs.
- Both `_prd.md` and `_techspec.md` account for identifier strategy,
  authentication and HTTP, active ADR obligations, and tooling authority with
  operative `docs/agents/` citations. ADR-0080, ADR-0081, and ADR-0091 are
  accepted and applicable; ADR-0093 is expressly accounted for as
  non-applicable.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` expressly
  authorizes `.agents/skills/qa-gate/**` and `skills/qa-gate/**` for Spec
  0070. Commit `d8c0403` recorded that authorization and the Spec artifacts
  before the tooling Task. Commit `45ffc72` changed only `task_04.md` to name
  the required regeneration chain before the tooling commit `f87db11`.
- Exact `rtk git diff-tree --no-commit-id --name-only -r <commit>` inspection:
  `7099ab9`, `765cf97`, `0adef28`, and `b6ea034` contain only their Task file
  plus product/test slice paths. `f87db11` contains the two authorized QA-skill
  files, `task_04.md`, and twelve generated files under
  `internal/baseline/assets/setups/` or `internal/baseline/testdata/`; both
  roots are included by `DERIVED_DIGEST_PATHS` in `Makefile:103`.
- `rtk cmp -s .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md`
  exited 0.
- A disposable local clone at build `b6ea034` reproduced the sanctioned
  generation without changing the QA worktree. `rtk make baseline-digests`
  exited 0 and reported `changed:false`. The two characterization commands
  first hit the sandbox-inaccessible host Go cache; rerunning them with a
  disposable cache under `/private/tmp` passed 7 and 2 tests. A final
  `rtk git status --short` in the clone was empty, proving the committed
  derived artifacts reproduce from their canonical sources.
- No earlier file existed under this Spec's `qa/` directory before this run.
  The current PRD contains no `## Unreachable Acceptance` section, so the
  current gate has zero pre-run declarations it may match.
