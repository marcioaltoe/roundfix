---
schema: roundfix/archive-record/v1
spec: 0388-a-completion-notice-skips-a-record-the-merge-adopted
title: A completion notice skips a record the merge adopted
status: archived
created: "2026-11-13"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0388-a-completion-notice-skips-a-record-the-merge-adopted
source_revision: 61418487b1c5e807b0794d12951a66533f28b41e
qa_task: task_02
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven:
  - the operator's or Argus's first managed refresh after the release, recorded session evidence.
adrs:
  - ADR-0396
sources:
  - 2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
promoted: []
pull_request: "861"
type: fix
---

# A completion notice skips a record the merge adopted

The completion-notice guidance now keeps adopted records with their consuming Spec, asks the user when adoption is ambiguous, restores the send channel and ships matching canonical/embedded skill version 0.1.229. Installation, regression detection, historical Git evidence and bounded scope are verified locally; post-release adopter refresh remains unproven and must be recorded by the operator or Argus. The offline archive plan exited 0 with the judge skipped and identified only the adopted defect report as a promotion candidate; its reusable correction already belongs to the release reference, so no promotion is proposed. Offline archive plan (`docs/specs/0388-a-completion-notice-skips-a-record-the-merge-adopted/qa/evidence/qa-report-2026-11-13/archive-plan.txt` (not reachable from any branch or tag)).
