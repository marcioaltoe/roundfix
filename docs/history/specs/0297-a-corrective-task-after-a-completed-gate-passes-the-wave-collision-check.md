---
schema: roundfix/archive-record/v1
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
title: A corrective Task after a completed gate passes the Wave collision check
status: archived
created: "2026-10-22"
archived: "2026-10-22"
disposition: partial
source: docs/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
source_revision: 1aa6ae3a87012502654cfa3a18bd3e7b150cab22
qa_task: task_03
qa_report: qa-report-2026-10-22.md
qa_verdict: partial
unproven:
  - the adopter operator's newest QA Report of that Spec after this Spec's release, whose frontmatter carries no `precondition_reason` naming `SC-WAVE-COLLISION`.
adrs: []
sources:
  - 2026-10-21-sc-wave-collision-no-portao-com-needs-e-check-estrito-limpo.md
regeneration:
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
promoted: []
relinked:
  - docs/backlog/2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md
pull_request: "632"
---

# A corrective Task after a completed gate passes the Wave collision check

The gate records a qualifying declared partial: all nine runnable rows pass, the prohibited publication lookup and absent Pull Request remain environment-blocked, and the single post-release adopter declaration remains unproven. The Run-built offline archive plan exited 0 and reported judge skipped offline; its evidence records the proposed Archive Record and disposition. The Daemon retains ownership of Task status and settlement. No archive, audited-checkout commit, push or Pull Request was performed.
