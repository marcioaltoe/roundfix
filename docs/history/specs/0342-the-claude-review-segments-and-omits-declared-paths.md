---
schema: roundfix/archive-record/v1
spec: 0342-the-claude-review-segments-and-omits-declared-paths
title: The Claude review segments a large candidate and omits declared paths
status: archived
created: "2026-11-03"
archived: "2026-11-03"
disposition: partial
source: docs/specs/0342-the-claude-review-segments-and-omits-declared-paths
source_revision: 074fda9ab58097d3c21969829dcc17f0cbe5df80
qa_task: task_05
qa_report: qa-report-2026-11-03.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0227
  - ADR-0282
sources:
  - 2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md
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

# The Claude review segments a large candidate and omits declared paths

The Archive Record can carry this qualifying partial with only the pre-PR Pull Request row unproven. The isolated `./bin/roundfix archive 0342-the-claude-review-segments-and-omits-declared-paths --plan --offline` exited 0 and skipped the judge, naming the adopted Vortex report as its only promotion candidate; see archive-plan.txt (`docs/specs/0342-the-claude-review-segments-and-omits-declared-paths/qa/evidence/2026-11-03-gate/archive-plan.txt` (not reachable from any branch or tag)). No additional promotion is needed because ADR-0352, the amended ADRs, glossary and public guides already hold the reusable decisions. The plan measured the tree at execution; later report and evidence additions change its file and byte counts. No archive, Task status, commit, push or Pull Request mutation was performed.
