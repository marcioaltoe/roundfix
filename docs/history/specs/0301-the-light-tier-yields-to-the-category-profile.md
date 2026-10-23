---
schema: roundfix/archive-record/v1
spec: 0301-the-light-tier-yields-to-the-category-profile
title: The Light Tier yields to the category profile when it cannot serve the requested access
status: archived
created: "2026-10-22"
archived: "2026-10-23"
disposition: partial
source: docs/specs/0301-the-light-tier-yields-to-the-category-profile
source_revision: dcd92652ce4165aa0c76840ec4f866b20df58e1d
qa_task: task_04
qa_report: qa-report-2026-10-23.md
qa_verdict: partial
unproven:
  - the operator's notice to Fluxus after release and the next Fluxus Run with `openrouter.light_models` restored in User Config.
adrs:
  - ADR-0314
sources:
  - 2026-10-21-the-light-tier-fails-a-task-under-agent-full-access.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Light Tier yields to the category profile when it cannot serve the requested access

The Archive Record can carry a qualifying declared partial at the audited head: the light-tier access skip, pre-work yield, unchanged fallback and spend behavior, Doctor output and documentation passed the hermetic gate, and all eleven sabotage mutations were detected. The single declared post-release Fluxus criterion remains unproven, alongside the permitted offline outside-evidence and pre-PR environment blocks. The read-only `./bin/roundfix archive 0301-the-light-tier-yields-to-the-category-profile --plan --offline` exited 0 with every Judge Key Variable and TYPESAFE_API_KEY unset; its cut and candidates are recorded in `evidence/2026-10-23-gate/archive-plan.txt`. Durable behavior already lives in ADR-0314, CONTEXT.
