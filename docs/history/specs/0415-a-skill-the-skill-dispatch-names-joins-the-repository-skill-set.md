---
schema: roundfix/archive-record/v1
spec: 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
title: A skill the skill dispatch names joins the Repository Skill Set
status: archived
created: "2026-11-20"
archived: "2026-11-20"
disposition: partial
source: docs/specs/0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
source_revision: f442e41aae13ff6aa763f388132846f09eff479d
qa_task: task_02
qa_report: qa-report-2026-11-20-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0448
sources: []
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/typescript.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
pull_request: "943"
type: fix
---

# A skill the skill dispatch names joins the Repository Skill Set

The Archive Record can carry a qualifying partial: all 17 runnable rows passed, including the corrected TypeScript-only adopter preview and all four negative probes, and the sole unmet row is the pre-PR Pull Request row. The offline archive plan lists the maintainer answers and characterization as two promotion candidates with no judge advice; the accepted ADR-0448 and the Skill Activation definition already retain the confirmed rule, so this gate promotes nothing and performs no archive. Offline archive plan (`docs/specs/0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set/qa/evidence/qa-report-2026-11-20-01/archive-plan.txt` (not reachable from any branch or tag)).
