# QA-01 — Governance and tooling provenance

Status: pass.

- `_tasks.md` names `task_07` as the sole `qa` node. It is terminal and needs
  both non-QA leaves, `task_05` and `task_06`.
- A fresh status sweep found `task_01` through `task_06` `completed` and only
  the Daemon-owned `task_07` `pending`.
- Both `_prd.md` and `_techspec.md` contain Project Constraints rows for
  identifier strategy, authentication and HTTP, active ADR obligations, and
  tooling authority. Each row states applicability and cites an operative
  source under `docs/agents/`.
- The applicable implementation obligations were checked against accepted
  ADR-0052, ADR-0053, and ADR-0081. The PRD additionally accounts for the QA
  workflow obligations from ADR-0080 and ADR-0091 and dismisses ADR-0093's
  relation candidate with a subsystem-specific reason.
- Commit `9bdaedbe20102634f3b34500668ca1bf51e06ced` records the standing
  authorization and explicitly names Spec 0068 before Spec authoring commit
  `2b87f051a438e276d6d0e1c6a626c42344b70627` and tooling Task commit
  `1346d83d4213e10b73a89bae6796d6d95dda6c18`.
- `git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r
  1346d83d4213e10b73a89bae6796d6d95dda6c18` found only the authorized
  canonical and mirrored Roundfix Skill, `task_06.md`, and deterministic
  artifacts under `DERIVED_DIGEST_PATHS`. The four plan-characterization
  goldens are the two explicitly required post-regeneration corpora named by
  the standing authorization and Task.
- No prerequisite or consequent repair commit exists between authorization,
  authoring, and the Task commit; the Task declared its deterministic fallout
  before execution.

Fresh reproduction:

- `rtk make skills-sync-check` — exit 0; four Skill tests passed.
- `rtk go test ./internal/baseline -count=1 -run
  'TestBaselinePlanCharacterization|TestCatalogDiagnosticCharacterization'` —
  exit 0; nine tests passed against the committed regenerated outputs.
- `rtk make verify` — exit 0 on the same build, including the Repository Skill
  Set and Spec Consistency Check.

The current delta before flow QA contains only this gate's report and evidence.
