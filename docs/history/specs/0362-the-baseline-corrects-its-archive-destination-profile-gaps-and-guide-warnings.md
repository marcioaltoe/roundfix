---
schema: roundfix/archive-record/v1
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
title: The Baseline corrects its archive destination, profile gaps and guide warnings
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
source_revision: d1bef5f8a794cf6224dbd3f1dbf80983f4d4edc9
qa_task: task_05
qa_report: qa-report-2026-11-06.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0372
  - ADR-0330
sources:
  - 2026-11-06-the-archive-clause-names-the-retired-secondbrain-inbox.md
  - 2026-11-06-a-repository-profile-without-a-template-decision-fails-opaquely.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/docs-layout.md
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/setup-context-driven/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/setup-context-driven/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# The Baseline corrects its archive destination, profile gaps and guide warnings

The Baseline now routes archive knowledge through ai-memory to the project that owns the Secondbrain, names a repository profile’s …
