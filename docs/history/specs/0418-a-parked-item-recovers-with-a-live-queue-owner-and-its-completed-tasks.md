---
schema: roundfix/archive-record/v1
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
title: A parked item recovers with a live Queue Owner and its completed Tasks
status: archived
created: "2026-11-23"
archived: "2026-11-24"
disposition: partial
source: docs/specs/0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
source_revision: 7f8565feaeb1ba6bfdb7f474e85db5d1e17647f3
qa_task: task_05
qa_report: qa-report-2026-11-24-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0279
  - ADR-0395
  - ADR-0450
sources: []
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
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "951"
type: fix
---

# A parked item recovers with a live Queue Owner and its completed Tasks

The retry and status recovery contracts passed the runnable QA sweep and all eight regression probes, with the repository Verification supplied passing by the Daemon; the report closes as a qualifying partial solely because the pre-PR Pull Request row lacks an open PR and Merge-Ready acceptance. The offline archive plan exited successfully, skipped the judge and listed the characterization and source capture as candidates; their reusable rules already live in ADR-0450, the glossary and deliver documentation, so this gate promotes neither historical capture and performs no archive.
