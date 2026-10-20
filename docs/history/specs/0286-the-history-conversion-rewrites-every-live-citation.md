---
schema: roundfix/archive-record/v1
spec: 0286-the-history-conversion-rewrites-every-live-citation
title: The history conversion rewrites every live citation
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0286-the-history-conversion-rewrites-every-live-citation
source_revision: d7a0595ac59ccd2ff4f6afcc7ff4251b8d08c8f6
qa_task: task_05
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0299
sources:
  - 2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "601"
---

# The history conversion rewrites every live citation

The conversion rewrites all nine listed transcript citations with the expected locators and lineage date, preserves the pending target and excluded files, and is empty on a fresh second plan. The strict checker reports all seven retained citation errors with file, line and replacement locator. Scope and byte-invariant audits pass. This is a qualifying partial whose only unmet observations are the pre-PR Pull Request controls and three outside web sources prohibited by authorization. The adopted Inbox Entry remains Spec-owned source material, while durable rules are already in ADR-0299, CONTEXT.md and the guides; no promotion or archive mutation is performed by QA.
