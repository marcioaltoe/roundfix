---
schema: roundfix/archive-record/v1
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
title: A reused owner PID means the owner is gone
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0363-a-reused-owner-pid-means-the-owner-is-gone
source_revision: 017d2b990bc1f41714cbafb5b22a31c807eabe1b
qa_task: task_04
qa_report: qa-report-2026-11-06.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0373
sources:
  - 2026-11-06-doctor-residue-stays-partial-when-an-owner-pid-is-reused.md
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
pull_request: "776"
type: fix
---

# A reused owner PID means the owner is gone

The owner-liveness rule now treats a proven reused PID as the end of the recorded owner: Doctor avoids false partial residue and stale locks, slots and queue owners can be released while matching or unprovable identities remain protected. All runnable checks and all six sabotage probes passed; the only unmet row is the pre-PR Pull Request control, with no Merge-Ready acceptance yet. The offline archive preview (`./bin/roundfix archive 0363-a-reused-owner-pid-means-the-owner-is-gone --plan --offline`, exit 0) is preserved in archive plan (`docs/specs/0363-a-reused-owner-pid-means-the-owner-is-gone/qa/evidence/2026-11-06-owner-reuse/archive-plan.txt` (not reachable from any branch or tag)); it made no archive, judge call, commit or push. The Spec-specific diagnostics and probe transcripts are audit evidence, and no additional durable reference needs promotion.
