---
schema: roundfix/archive-record/v1
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
title: The Baseline setups stop sending agents to the Secondbrain inbox
status: archived
created: "2026-10-26"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
source_revision: 1398a64d3626e4d168f17a0234b749b2625d03e2
qa_task: task_02
qa_report: qa-report-2026-10-27-01.md
qa_verdict: partial
unproven:
  - the release notice to `default/gesttione-solutions-pantheon` and that project's refresh and restore after the release, reported to this project through ai-memory.
adrs:
  - ADR-0060
  - ADR-0191
  - ADR-0264
sources:
  - 2026-10-26-the-knowledge-workspace-skill-still-writes-to-the-secondbrain-inbox.md
regeneration:
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: roundfix baseline skills restore --profile go-cli-tui --skill knowledge-workspace --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - .agents/skills/knowledge-workspace/SKILL.md
      - skills-lock.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
pull_request: "702"
---

# The Baseline setups stop sending agents to the Secondbrain inbox

The candidate carries the corrected upstream knowledge-workspace tree in all five Setup Snapshots and this repository's installed skill and lock.
