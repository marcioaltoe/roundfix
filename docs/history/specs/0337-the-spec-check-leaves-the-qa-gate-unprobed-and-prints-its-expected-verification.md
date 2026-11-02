---
schema: roundfix/archive-record/v1
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
title: The Spec Check leaves the QA gate unprobed and prints its expected Verification
status: archived
created: "2026-11-01"
archived: "2026-11-01"
disposition: partial
source: docs/specs/0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
source_revision: 8757d1f9a54f339c3acaeeaa8c4d72086889efe6
qa_task: task_03
qa_report: qa-report-2026-11-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0148
sources:
  - 2026-11-01-the-spec-check-probe-judges-the-qa-gates-derived-verification-vacuous.md
  - 2026-11-01-no-authoring-path-shows-the-qa-gates-derived-verification.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# The Spec Check leaves the QA gate unprobed and prints its expected Verification

Spec Check leaves the QA gate unprobed through the shared Daemon predicate and exposes the exact authored-QA Verification replacement in text and JSON. Both reported adopter cases pass through the built CLI, package and documentation checks pass, and all six isolated mutations are detected by their named tests. QA closes as a qualifying partial solely because no Pull Request exists and Merge-Ready acceptance is not yet available; the offline archive preview succeeds with the judge skipped and both adopted reports listed as candidates. The guide, glossary and ADR already preserve the reusable rules, so no additional promotion is proposed.
