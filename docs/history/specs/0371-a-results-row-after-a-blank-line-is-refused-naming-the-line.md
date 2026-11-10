---
schema: roundfix/archive-record/v1
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
title: A Results row after a blank line is refused, naming the line
status: archived
created: "2026-11-10"
archived: "2026-11-10"
disposition: partial
source: docs/specs/0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
source_revision: bf0c3a13c3a28922c90be539a1878f1e461d4af4
qa_task: task_04
qa_report: qa-report-2026-11-10-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0381
sources:
  - 2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - skills/qa-gate/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "811"
type: fix
---

# A Results row after a blank line is refused, naming the line

A Results row separated from its table by a blank line is refused with the report path, exact row line and repair, and the Daemon instructs the gate to continue the seeded table and offers one shared repair turn. The fresh suites, eight sabotage checks, wording and glossary checks, and Task-scope audit support that result at the audited head; the report closes as a qualifying partial solely because no Pull Request exists for the remaining Merge-Ready and review-ancestry controls. The offline archive plan identifies one adopted provenance reference, with its reusable knowledge already preserved by ADR-0381 and the updated glossary and guides.
