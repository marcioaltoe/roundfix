# Project Constraint and protected-tooling audit

Build: `bdf6ff8d4d680188a97986ee1340ab9ff052a499`.

## Spec shape

Fresh `rtk rg` inspection found all four Project Constraint rows in both
`_prd.md` and `_techspec.md`. Identifier strategy and authentication/HTTP are
reasoned non-applicable and cite operative sources under `docs/agents/`.
Active ADR obligations and Tooling Authority are applicable and cite their
operative sources.

The active cited decisions were read from disk: ADR-0029 retains review
artifact placement, ADR-0036 owns the separate artifact-only docs commit,
ADR-0054 keeps Review Source Evidence authoritative, ADR-0080 owns typed QA
blocking and equivalent evidence, ADR-0081 sanctions deterministic digest
fallout, ADR-0091 owns the authored terminal QA node, and ADR-0093 is expressly
recorded as a non-applicable relation candidate. No conflict was found.

## Authorization and chronology

The authorization record
`docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` names Spec 0078,
records the maintainer's 2026-08-05 `skill + mirror` authorization, and bounds
the cause to `.agents/skills/roundfix/**` plus `skills/roundfix/**` for shipped
CLI-contract synchronisation. It names `make skills-sync`, `make
baseline-digests`, and the two characterization-corpus re-records.

`rtk git -c core.fsmonitor=false merge-base --is-ancestor f253cc92 bdf6ff8d`
exited 0. Chronological history begins with the separate authorization commit
`f253cc92`, continues through the Spec and implementation commits, and ends
with Task 05 commit `bdf6ff8d`. The authorization is neither late nor folded
into the Task commit. No prerequisite or consequent tooling-fix commit exists.

## Exact Task 05 paths

`rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-status -r
bdf6ff8d` exited 0 and listed only:

- `.agents/skills/roundfix/SKILL.md`;
- `skills/roundfix/SKILL.md`;
- `docs/specs/0078-roundfix-asks-for-the-review/task_05.md`;
- deterministic files under `internal/baseline/assets/setups/` and
  `internal/baseline/testdata/`, including the explicitly required
  characterization corpus.

The canonical Skill and mirror are the expressly authorized cause. The Task
file is the assigned Task artifact. Every remaining path is inside
`DERIVED_DIGEST_PATHS`, which the Makefile resolves to Baseline asset/testdata
trees and ADR-0081 treats as sanctioned deterministic fallout. No Go source,
Makefile, ignore file, third-party Skill, or other protected path changed in
the Task commit.

Fresh reproducibility is closed by the report's authoritative `make verify`
row and Skill-sync row; this audit does not rely on the Task's prior claim for
their current result.
