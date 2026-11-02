---
schema: roundfix/archive-record/v1
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
title: QA evidence links resolve in the committed tree
status: archived
created: "2026-11-01"
archived: "2026-11-02"
disposition: partial
source: docs/specs/0340-qa-evidence-links-resolve-in-the-committed-tree
source_revision: f0a93e483a2b19daf9e6513a8332af661800f850
qa_task: task_04
qa_report: qa-report-2026-11-02.md
qa_verdict: partial
unproven: []
adrs: []
sources:
  - 2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md
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
pull_request: "724"
---

# QA evidence links resolve in the committed tree

The gate reproduced refusal of ignored and missing evidence through the Run-built CLI, accepted retained .txt evidence on fresh invocations, and observed one repair turn with committed evidence or refusal without a QA Report commit in the Daemon tests. All required positive suites and seven sabotage probes satisfied their expectations; the Task scope, authorization, glossary and offline outside-source comparisons passed. The report closes partial solely because no Pull Request exists. The offline archive plan exited 0 with judge skipped and one adopted historical Argus reference candidate; the reusable rule already lives in ADR-0351, CONTEXT and shipped guides, so no new promotion is needed. This QA Agent performed no archive, Task-status change, repository commit, push or Pull Request creation.
