---
schema: roundfix/archive-record/v1
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
title: The Baseline sends handoffs and messages through ai-memory
status: archived
created: "2026-10-24"
archived: "2026-10-25"
disposition: qa-override
source: docs/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
source_revision: ef1a4bab3f4b23370377de142c40ff97aa61e631
qa_task: task_05
qa_report: qa-report-2026-10-25.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'Marcio, AskUserQuestion 2026-10-25 ~12:40 BRT: "Autorizo o override (Recomendado)"'
qa_override_reason: 'QA partial with no failed, skipped or finding-blocked row; QA-08 (upstream setup lists) was recorded as ''upstream checkout and network prohibited by Spec authorization'', outside the accepted wording; the operator verified it: the embedded go, rust and typescript setups equal setups/*.txt of marcioaltoe/skills 65d5f365 and record that commit. QA-09/QA-13 are the accepted outside-evidence and Pull Request rows; QA-10 is the declared Unreachable Acceptance.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: ef1a4bab3f4b23370377de142c40ff97aa61e631
adrs:
  - ADR-0330
  - ADR-0095
sources:
  - 2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md
regeneration:
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: roundfix baseline skills restore --profile go-cli-tui --skill ai-memory-durable-pages --skill ai-memory-handoff --skill ai-memory-learning-maintenance --skill ai-memory-messaging --skill ai-memory-retrieval --skill ai-memory-routing-install --skill context7-cli --skill golang-cli --skill grilling --skill knowledge-workspace --skill systematic-debugging --skill typesafe-ai --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - .agents/skills/ai-memory-durable-pages/LICENSE
      - .agents/skills/ai-memory-durable-pages/SKILL.md
      - .agents/skills/ai-memory-handoff/LICENSE
      - .agents/skills/ai-memory-handoff/SKILL.md
      - .agents/skills/ai-memory-learning-maintenance/LICENSE
      - .agents/skills/ai-memory-learning-maintenance/SKILL.md
      - .agents/skills/ai-memory-messaging/LICENSE
      - .agents/skills/ai-memory-messaging/SKILL.md
      - .agents/skills/ai-memory-retrieval/LICENSE
      - .agents/skills/ai-memory-retrieval/SKILL.md
      - .agents/skills/ai-memory-routing-install/LICENSE
      - .agents/skills/ai-memory-routing-install/SKILL.md
      - .agents/skills/context7-cli/SKILL.md
      - .agents/skills/context7-cli/references/docs.md
      - .agents/skills/context7-cli/references/setup.md
      - .agents/skills/context7-cli/references/skills.md
      - .agents/skills/golang-cli/assets/examples/go.mod
      - .agents/skills/grilling/SKILL.md
      - .agents/skills/knowledge-workspace/SKILL.md
      - .agents/skills/systematic-debugging/SKILL.md
      - .agents/skills/typesafe-ai/SKILL.md
      - skills-lock.json
  - command: roundfix baseline skills reconcile --profile go-cli-tui --source marcioaltoe/skills --revision 65d5f365cace27c4b468db1723213723b10c870b --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - skills-lock.json
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Baseline sends handoffs and messages through ai-memory

The only archive command was the Run-built binary's `archive 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory --plan --offline`, with fresh HOME and TMPDIR and every Judge Key Variable plus TYPESAFE_API_KEY unset. It exited 0, reported an archive plan and `judge: skipped (offline)`, and applied nothing. The archive outcome remains a partial QA report with Success Metric 6 unproven; the official eligibility reader subsequently rejected the authored upstream-source status, so normal settlement is not currently established. Any later settlement or archive belongs to the Daemon.
