---
schema: roundfix/archive-record/v1
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
title: A Supervisor session ends only when the maintainer asks
status: archived
created: "2026-11-17"
archived: "2026-11-17"
disposition: qa-override
source: docs/specs/0401-a-supervisor-session-ends-only-when-the-maintainer-asks
source_revision: 63136c0f0ecd167ad8a868152dd373d02b1f70ad
qa_task: task_04
qa_report: ""
qa_verdict: ""
unproven: []
qa_override: true
qa_override_approval: Marcio, standing authorization 2026-08-09 (qa_override por ambiente, rows_blocked_finding=0) and 2026-11-17 unattended grant
qa_override_reason: 'QA fail only by environment (no finding-blocked rows): isolated HOME had no rustup toolchain for TestGuidanceCompositionJourney/rust-cli (same failure on the unchanged Delivery Base); ./internal/cli, ./internal/phaserecord and make verify-changed pass with the host toolchain on 63136c0f; the PR row is satisfied by the delivery checks.'
qa_override_qa_outcome: missing
qa_override_qa_task_status: pending
qa_override_revision: 63136c0f0ecd167ad8a868152dd373d02b1f70ad
adrs:
  - ADR-0414
  - ADR-0413
sources:
  - 2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md
  - 2026-11-16-maintainer-answers.md
  - 2026-11-17-phase-record-boundary-measurement.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/autonomous-work.md
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-idea/SKILL.md
      - .agents/skills/write-prd/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - .agents/skills/write-techspec/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
      - skills/write-prd/SKILL.md
      - skills/write-tasks/SKILL.md
      - skills/write-techspec/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "899"
type: fix
---

# A Supervisor session ends only when the maintainer asks

Spec 0397 put a Phase Boundary after every Authoring Phase and after a Spec's delivery and release: the four authoring skills, the Roundfix skill's `--follow-up` and `--unattended` modes and the `autonomous-work` Baseline clause tell the Supervisor to write the Session Handoff and end the session there, and `roundfix phase start` prints the Boundary Notice whenever the same session records a different phase or Spec. The maintainer found first that the stop after each phase slowed the work, and then that the stop at the end of each Spec also harms autonomous work, because an unattended session halts and asks to be restarted. He decided that no phase of a Spec ends the session. After this Spec the Supervisor follows its queue from phase to phase and from Spec to Spec in one session, and writes the Session Handoff and ends only when the maintainer asks or the work stops; `roundfix phase start` keeps recording every phase and prints no notice.
