---
schema: roundfix/archive-record/v1
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
title: A large authored data file is reviewed apart, not omitted
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
source_revision: 68e7619304191c887b2885b06265689856ae32ad
qa_task: task_05
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0368
  - ADR-0282
  - ADR-0352
sources:
  - 2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md
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
type: feat
---

# A large authored data file is reviewed apart, not omitted

A declared authored catalog is now reviewed apart from code in bounded Data Review Segments, with valid file-line anchors and distinguishable findings in both providers' records and answers. All runnable checks and all five sabotage probes passed; the only unmet row is the pre-PR Pull Request row because no Pull Request is open. The offline archive plan ran without judge access (`docs/specs/0358-a-large-authored-data-file-is-reviewed-apart-not-omitted/qa/evidence/2026-11-05-data-review/archive-plan.txt` (not reachable from any branch or tag)); reusable terminology and guidance already live in CONTEXT.md, the review/configuration guides and accepted ADR-0368, so this QA adds no knowledge promotion.
