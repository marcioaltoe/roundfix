---
schema: roundfix/archive-record/v1
spec: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
title: The Catalog Skill Bump refreshes with a binary built after the digests
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
source_revision: 67e4e9bef1a0e93d972bc9e498cdd6e41c444375
qa_task: task_02
qa_report: qa-report-2026-11-11.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0360
sources:
  - 2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
type: fix
---

# The Catalog Skill Bump refreshes with a binary built after the digests

The Catalog Skill Bump rebuilds its refresh binary after digest generation so the Setup Manifest reflects the regenerated assets. Both sabotage cases and the Delivery Base fail the regression checks, while the real rebuilt binary leaves a current manifest. The offline archive plan (`docs/specs/0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests/qa/evidence/qa-report-2026-11-11/archive-plan.txt` (not reachable from any branch or tag)) lists the cut and candidates with the judge skipped; the existing glossary and ADR already own the durable knowledge, so no promotion is needed. Current-candidate review and Pull Request controls remain for the Supervisor after this qualifying partial QA gate.
