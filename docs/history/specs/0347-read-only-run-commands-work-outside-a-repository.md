---
schema: roundfix/archive-record/v1
spec: 0347-read-only-run-commands-work-outside-a-repository
title: Read-only Run commands work outside a repository
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0347-read-only-run-commands-work-outside-a-repository
source_revision: 9f01b36a191bb4b12f7ba82a9c656ed12ef4b691
qa_task: task_05
qa_report: qa-report-2026-11-04.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0357
sources:
  - 2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md
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
  - command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
    outputs:
      - docs/references/coverage-record.json
promoted: []
type: fix
---

# Read-only Run commands work outside a repository

The Run-built CLI now lists every repository outside Git, keeps repository scope inside Git, and prints one warning instead of a worktree-location refusal on Read-Only Run Commands; strict commands still refuse with the detected root and `.git` named. The required suites, exact CLI transcripts, eight negative probes, documentation and history audits passed. The report closes as qualifying partial solely because no Pull Request exists and Merge-Ready acceptance is not yet available.
