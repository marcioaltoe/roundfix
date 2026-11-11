---
schema: roundfix/archive-record/v1
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
title: A Machine Verification Capacity bounds Verifications across repositories
status: archived
created: "2026-11-10"
archived: "2026-11-11"
disposition: partial
source: docs/specs/0374-a-machine-verification-capacity-bounds-verifications-across-repositories
source_revision: 62d6b642c69376121822f811a83e3fd7c51afcbc
qa_task: task_05
qa_report: qa-report-2026-11-11.md
qa_verdict: partial
unproven:
  - the operator's runs list and doctor evidence from the first cycle after the release, recorded session evidence.
adrs:
  - ADR-0384
  - ADR-0373
  - ADR-0056
sources:
  - 2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md
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
promoted:
  - docs/references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md
relinked:
  - docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light/_prd.md
pull_request: "821"
type: feat
---

# A Machine Verification Capacity bounds Verifications across repositories

The Archive Record can carry a qualifying declared partial: all runnable checks passed without product findings, while the operator's first post-release cycle with two real repositories remains unproven under the PRD satisfied-by declaration.
