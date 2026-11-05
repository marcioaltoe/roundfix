---
schema: roundfix/archive-record/v1
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
title: runs list reads as a table, and the Run Browser draws, at a terminal
status: archived
created: "2026-11-05"
archived: "2026-11-05"
disposition: partial
source: docs/specs/0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
source_revision: 8b448cc96648ae89fa62fc9c82a1cc6aeb3030b5
qa_task: task_04
qa_report: qa-report-2026-11-05.md
qa_verdict: partial
unproven:
  - the model and command tests this Spec adds, and the maintainer running bare `roundfix attach` on the MacBook after v0.133.0
adrs:
  - ADR-0357
  - ADR-0363
sources:
  - 2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md
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
pull_request: "755"
type: feat
---

# runs list reads as a table, and the Run Browser draws, at a terminal

QA settles Spec 0353 as a qualifying declared partial at head `23c1440e3856402622003253e3d07b787d4cbf1e`. The Run Browser, the attach cockpit and the owning cockpit receive the terminal file and refuse any other writer by type. A Run Browser closed without a selection prints one stderr line and exits `0`. `runs list` keeps its stable line byte for byte when piped or with `--format plain`, and prints an aligned table with an eight-digit short id, local start time and a width-bounded `TARGET` at a terminal, with `--format table` or with `--wide`. `attach` and `runs show` resolve a unique short id and refuse an ambiguous one listing every match.
