---
schema: roundfix/archive-record/v1
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
title: Runs load the dispatched skills, and the Baseline follows the curated catalog
status: archived
created: "2026-10-30"
archived: "2026-10-31"
disposition: partial
source: docs/specs/0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
source_revision: 6f9af76df69c2fc0c879d87e8c79c553e65931ac
qa_task: task_05
qa_report: qa-report-2026-10-31.md
qa_verdict: partial
unproven:
  - the operator reruns `measure-frontend-skill-loads.py` over the Runs started after the release and reports the counts to `default/marcioaltoe-roundfix` through ai-memory
  - the operator's first managed refresh of a TypeScript adopter after the release, reported to `default/marcioaltoe-roundfix` through ai-memory
adrs:
  - ADR-0346
sources:
  - 2026-10-30-the-baseline-follows-the-curated-skills-catalog.md
  - 2026-10-30-the-completed-ui-ux-pro-max-updates-automatically.md
regeneration:
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: roundfix baseline skills restore --profile go-cli-tui --skill context7-cli --skill grilling --skill tui-design --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - .agents/skills/context7-cli/SKILL.md
      - .agents/skills/context7-cli/references/docs.md
      - .agents/skills/grilling/SKILL.md
      - .agents/skills/tui-design/SKILL.md
      - .agents/skills/tui-design/references/evidence.md
      - skills-lock.json
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/backend.json
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/frontend.json
      - internal/baseline/assets/modules/typescript.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "714"
---

# Runs load the dispatched skills, and the Baseline follows the curated catalog

The audited candidate refreshes the Baseline to catalog 86df860, separates product and public-page dispatch, carries the matching-skill load rule through Task/QA prompts and carry-forward, and retires implement-spec while preserving its delivery rules. The positive suites, all 19 named tests, eight negative probes and public CLI journeys passed; frontend usage and a released TypeScript adopter refresh remain unproven under their two PRD declarations. The offline archive plan exited 0 with judge skipped and listed the two adopted catalog messages as candidates; their lasting knowledge is already held in the Baseline and live guides, so no promotion is proposed. No archive, commit, push or Pull Request mutation was performed.
