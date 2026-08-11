# Project Constraint audit

Build: `ef6eb44ad8951112b1c3641bb7fd21793b440f95`

## Task status

All seven canonical Task files carry `status: completed`.

## Active artifacts

The PRD and TechSpec each account for all mandatory Project Constraint axes:

- Identifier strategy: not applicable because this Spec reuses existing Run,
  process, Agent Session, Work Item, and scope identities.
- Authentication and HTTP: not applicable because the changed local process,
  Run Database, CLI, and polling behavior introduces no authentication or HTTP
  contract.
- Active ADR obligations: applicable. ADR-0022, ADR-0044, ADR-0051, and
  ADR-0052 resolve and bind Stop Request transport, owner-death proof,
  Work Item-scoped Agent Sessions, and terminal compare-and-set behavior.
- Tooling authority: applicable. Both artifacts expressly authorize the
  canonical/generated Roundfix Skill pair and the five deterministic
  Skill-digest artifacts listed below.

All four operative source citations resolve under `docs/agents/`.

## Daemon-owned changed paths

Task 06 commit `87d0fd0c2adfc8c30348f54e68437181d3ee3003`:

```text
docs/specs/0037-terminal-outcome-integrity/task_06.md
docs/user-guide/commands.md
docs/user-guide/usage.md
internal/cli/cli.go
internal/cli/cli_test.go
```

This commit changes no protected repository-tooling path.

Task 07 commit `9fc2ba491839ce46aae9b70631c5f9cd3c9f05ba`:

```text
.agents/skills/roundfix/SKILL.md
docs/specs/0037-terminal-outcome-integrity/task_07.md
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
skills/roundfix/SKILL.md
```

The two Skill paths and all five derived digest paths are expressly bounded in
both active Spec artifacts. The only additional path is Task 07's own file.

Follow-up commit `ef6eb44ad8951112b1c3641bb7fd21793b440f95`:

```text
docs/specs/0037-terminal-outcome-integrity/_prd.md
docs/specs/0037-terminal-outcome-integrity/_techspec.md
docs/specs/0037-terminal-outcome-integrity/task_07.md
```

It records the exact maintainer authorization in the two active artifacts and
aligns Task 07's evidence contract; it changes no tooling artifact.

`rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited zero
with no output. The pair is byte-identical.

Verdict: pass.
