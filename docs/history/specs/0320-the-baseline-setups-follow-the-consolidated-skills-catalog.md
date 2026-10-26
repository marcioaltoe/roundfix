---
schema: roundfix/archive-record/v1
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
title: The Baseline setups follow the consolidated skills catalog
status: archived
created: "2026-10-26"
archived: "2026-10-26"
disposition: partial
source: docs/specs/0320-the-baseline-setups-follow-the-consolidated-skills-catalog
source_revision: acc02e5cd09cfbc1cfba8ac882cfc23efa6c3e49
qa_task: task_02
qa_report: qa-report-2026-10-26.md
qa_verdict: partial
unproven:
  - the operator's next managed refresh of a TypeScript adopter after the release, reported to this project through ai-memory.
adrs:
  - ADR-0191
sources:
  - 2026-10-26-the-baseline-setups-trail-the-consolidated-skills-catalog.md
regeneration:
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
---

# The Baseline setups follow the consolidated skills catalog

The Archive Record can carry a qualifying declared partial for audited head `9eda3f904aa6034f170c15f66dad56657b1edcee`, with 11 passing rows, Q09's matching Unreachable Acceptance and its `satisfied-by` evidence marked unproven, and Q12's pre-PR environment block. The required `./bin/roundfix archive 0320-the-baseline-setups-follow-the-consolidated-skills-catalog --plan --offline` exited 0, listed the Spec cut and its adopted source candidate, and reported `judge: skipped (offline)` (archive plan (`docs/specs/0320-the-baseline-setups-follow-the-consolidated-skills-catalog/qa/evidence/archive-plan.txt` (not reachable from any branch or tag))).
