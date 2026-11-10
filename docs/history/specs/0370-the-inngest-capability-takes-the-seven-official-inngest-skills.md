---
schema: roundfix/archive-record/v1
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
title: The Inngest capability takes the seven official Inngest skills
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: partial
source: docs/specs/0370-the-inngest-capability-takes-the-seven-official-inngest-skills
source_revision: b6a73415fccc5e115175c05ec98e7726883ece75
qa_task: task_03
qa_report: qa-report-2026-11-09.md
qa_verdict: partial
unproven:
  - the operator runs the update in one TypeScript adopter after the patch release and records its `Skills installed` and `Skills deleted` lines through ai-memory, recorded session evidence.
adrs:
  - ADR-0350
  - ADR-0380
  - ADR-0206
  - ADR-0360
sources:
  - 2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/backend.json
      - internal/baseline/module-versions.json
  - command: make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0370/skills SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce
    outputs:
      - .agents/skills/context7-cli/references/setup.md
      - docs/agents/setup-context.json
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
      - skills-lock.json
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: roundfix baseline skills restore --profile go-cli-tui --skill <each catalog entry> --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - .agents/skills/context7-cli/references/setup.md
      - skills-lock.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# The Inngest capability takes the seven official Inngest skills

The backend Baseline requires and dispatches the seven official Inngest skills, retires inngest with deletion of an unedited copy and preservation of an edited one, and pins the upstream snapshots and catalog skills to e45466d2. Fresh tests, Delivery Base controls, three sabotages, an unchanged repeat bump and the updated documentation all support the local result; the qualifying partial retains the post-release adopter migration as unproven and records the absent Pull Request separately.
