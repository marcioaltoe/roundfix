---
schema: roundfix/archive-record/v1
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
title: The skills lock matches what is installed, and a Run without a QA gate settles green
status: archived
created: "2026-10-27"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
source_revision: 030365692b1f5ade864a97a5d7c711635a57c79d
qa_task: task_05
qa_report: qa-report-2026-10-27-01.md
qa_verdict: partial
unproven:
  - the release notice to `default/gesttione-solutions-quaestor` and that project's next upgrade, reported to this project through ai-memory.
adrs:
  - ADR-0189
  - ADR-0342
sources:
  - 2026-10-27-baseline-update-leaves-the-skills-lock-hash-of-the-catalog.md
  - 2026-10-27-a-run-without-a-qa-gate-ends-clean-with-the-repository-red.md
regeneration:
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - .agents/skills/write-tasks/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-tasks/SKILL.md
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
pull_request: "709"
---

# The skills lock matches what is installed, and a Run without a QA gate settles green

Quaestor reported two frictions on Roundfix v0.124.0 and saw both again on v0.125.0.
