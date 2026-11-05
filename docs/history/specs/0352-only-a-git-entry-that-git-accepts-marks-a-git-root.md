---
schema: roundfix/archive-record/v1
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
title: Only a .git entry that Git accepts marks a Git root
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0352-only-a-git-entry-that-git-accepts-marks-a-git-root
source_revision: f77d3423936dc8ba12da3b9c4f9f8eec0316aa2e
qa_task: task_04
qa_report: qa-report-2026-11-05-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0362
  - ADR-0357
sources:
  - 2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md
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
type: fix
---

# Only a .git entry that Git accepts marks a Git root

The Git-root fix passes the complete authored package and documentation checks, all three public command transcripts, nested repository recovery, independent Git comparisons, eight sabotage probes and documentation and history audits at the audited head. The gate closes as a qualifying partial solely because no Pull Request is open, leaving live approval and Merge-Ready acceptance unobserved. The offline archive preview skips the judge and names the adopted maintainer report as its sole non-core candidate; ADR-0362 and the updated guides already hold the reusable decision, so no additional knowledge needs promotion. No archive, commit, push, Pull Request or release was performed.
