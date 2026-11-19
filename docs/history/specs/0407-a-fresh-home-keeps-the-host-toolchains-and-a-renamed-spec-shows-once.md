---
schema: roundfix/archive-record/v1
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
title: A fresh home keeps the host toolchains, and a renamed Spec shows once
status: archived
created: "2026-11-18"
archived: "2026-11-19"
disposition: partial
source: docs/specs/0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
source_revision: 30f9fd9748a6bfd4df31db13bc5e98217b9587b3
qa_task: task_03
qa_report: qa-report-2026-11-19.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0434
sources:
  - 2026-11-18-the-qa-gate-isolated-home-lacks-the-host-rust-toolchain.md
  - 2026-11-18-maintainer-answers.md
  - 2026-11-18-characterization.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/qa-gate/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# A fresh home keeps the host toolchains, and a renamed Spec shows once

The candidate keeps a renamed Spec to one queue row while retaining a second Spec Number, and both QA skills preserve the host toolchain locations under a fresh home. The live Rust journey reproduced the original rustup failure and passed with preserved locations; the full CLI suite passed in that environment. All seven independent sabotages were detected. QA closes as a qualifying partial because the Spec has no open Pull Request, with no product finding or declared-unreachable acceptance.
