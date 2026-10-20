---
schema: roundfix/archive-record/v1
spec: 0287-reports-that-count-what-they-say
title: Reports that count what they say
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0287-reports-that-count-what-they-say
source_revision: 119f1cd389e54e4db80f16fb4f9bf9181a971ee0
qa_task: task_04
qa_report: qa-report-2026-10-20-01.md
qa_verdict: partial
unproven:
  - the operator runs `roundfix storage report` and `roundfix deliver lane --all` on `main` after this Spec merges.
adrs: []
sources:
  - 2026-10-20-the-linux-lane-counts-package-patterns-as-packages.md
  - 2026-10-20-the-storage-report-misses-the-default-artifact-roots.md
  - 2026-10-20-an-archive-record-summary-ends-at-the-colon-that-introduces-its-list.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# Reports that count what they say

The gate verified the container-selected Linux Lane count and explicit no-count line, read-only storage reporting under each resolved Run root, and Archive Record summaries that carry their introduced list or paragraph while existing records stay unchanged. The public journeys, named regression tests, documentation contracts and five guard-removal probes passed. The result is a qualifying declared partial with three environment-blocked rows and one combined declared maintainer criterion still unproven.
