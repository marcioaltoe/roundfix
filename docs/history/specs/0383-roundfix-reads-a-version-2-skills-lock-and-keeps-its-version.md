---
schema: roundfix/archive-record/v1
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
title: Roundfix reads a version 2 skills lock and keeps its version
status: archived
created: "2026-11-12"
archived: "2026-11-12"
disposition: partial
source: docs/specs/0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
source_revision: ce19cc3ff38443232ab62204a725109790208ba5
qa_task: task_04
qa_report: qa-report-2026-11-12-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0393
sources:
  - 2026-11-12-vortex-skills-lock-version-2-refused.md
  - 2026-11-12-pantheon-skills-lock-version-2-refused.md
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
pull_request: "845"
type: fix
---

# Roundfix reads a version 2 skills lock and keeps its version

The audited binary repairs the reported stale Roundfix entry in a version 2 skills lock with the independently computed published digest, removes its ref and preserves its version and repeat-install bytes; version 3 is refused by name, with upgrade next actions proved by the named command tests. All required positive tests, eight sabotage probes, Windows vet, documentation, glossary, outside-source and history checks passed, including the corrected selector that failed the earlier gate; only the absent Pull Request and its unavailable Merge-Ready acceptance keep the verdict a qualifying partial.
