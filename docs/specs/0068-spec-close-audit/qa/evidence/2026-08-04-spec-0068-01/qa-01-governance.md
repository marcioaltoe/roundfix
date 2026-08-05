# QA-01 — Graph, constraints, and tooling provenance

Status: pass.

The manifest names `task_07` as the only `type: qa` Task. It is terminal and
depends on `task_08`, the only non-QA leaf. Tasks 01–06 and 08 are `completed`;
the Daemon-owned current gate remains `pending` while this report is written.

Both `_prd.md` and `_techspec.md` account for the four Project Constraints:

- Identifier strategy is not applicable because the feature creates no
  project-owned Internal Identifier; `docs/agents/domain.md` is the operative
  source.
- Authentication and HTTP are not applicable because the audit opens no
  transport and cannot invent policy; `docs/agents/agent-instructions.md` is
  the operative source.
- Active ADR obligations are applicable. ADR-0052 protects Active Runs,
  ADR-0053 requires proof-based reconciliation, ADR-0080 owns blocked-row
  verdict semantics, ADR-0081 owns regenerated digest fallout, ADR-0091 owns
  the authored QA node, and ADR-0093 is explicitly considered and dismissed
  for the Spec consistency subsystem.
- Tooling authority is applicable only to the Roundfix Skill pair. The
  standing CLI-synchronisation grant at
  `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` names Spec
  0068 and bounds `.agents/skills/roundfix/**` plus `skills/roundfix/**`.

Git provenance:

- Authorization commit `8dda1969` predates the Spec authoring commit
  `2b87f051`, which predates authorized Task commit `1346d83d`; both ancestry
  checks exited 0.
- `git diff-tree --no-commit-id --name-only -r 1346d83d` shows only the Skill
  pair, task_06, and ADR-0081 derived files under `DERIVED_DIGEST_PATHS`,
  including the two characterization corpora the authorization requires.
- Task commits `c9c2c562`, `8b651ccf`, `e647538c`, `c765c8e2`, `25335522`, and
  `30ec663c` contain their assigned Task file and bounded implementation or
  fixture paths. Corrective authoring/verification commits `b1d03c86` and
  `c48d2e22` precede task_08's implementation commit.
- `diff -qr .agents/skills/roundfix skills/roundfix` exited 0,
  `rtk make skills-sync-check` passed four tests, and the full gate validated
  the Repository Skill Set and regenerated digests.
- The current worktree delta contains only this rerun's QA report and evidence.

No authorization, ordering, derived-pin, or out-of-scope tooling defect was
found.
