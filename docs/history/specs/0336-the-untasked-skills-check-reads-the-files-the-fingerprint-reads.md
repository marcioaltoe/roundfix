---
schema: roundfix/archive-record/v1
spec: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
title: The untasked-skills check reads the files the fingerprint reads
status: archived
created: "2026-10-31"
archived: "2026-10-31"
disposition: partial
source: docs/specs/0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
source_revision: 6fe06ccaaf50f1d3eddcef3272f62847c168bb8c
qa_task: task_02
qa_report: qa-report-2026-10-31.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0256
sources:
  - 2026-10-30-the-untasked-skills-check-counts-test-files-the-fingerprint-ignores.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "717"
---

# The untasked-skills check reads the files the fingerprint reads

The audited candidate excludes Go test files from skills declarations through the same predicate used by the config fingerprint, preserves enforcement for source files, and documents that rule in the guide and both skill references. The CLI reproduction, all required regression suites, four sabotage variants, fingerprint history and bounded-path audit passed. The gate is a qualifying partial solely because no Pull Request exists; Merge-Ready acceptance remains unproven. No new glossary term or additional durable knowledge needs promotion because the rule already lives in its guide and skill reference.
