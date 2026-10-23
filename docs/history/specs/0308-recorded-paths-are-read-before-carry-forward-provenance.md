---
schema: roundfix/archive-record/v1
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
title: Recorded paths are read before carry-forward provenance
status: archived
created: "2026-10-22"
archived: "2026-10-22"
disposition: partial
source: docs/specs/0308-recorded-paths-are-read-before-carry-forward-provenance
source_revision: c7a0900ccb1ed165cf084e21f420cdf657ca5276
qa_task: task_02
qa_report: qa-report-2026-10-22.md
qa_verdict: partial
unproven:
  - oraculum's own `roundfix deliver retry` with the released binary, reported back through the Secondbrain inbox.
adrs: []
sources:
  - 2026-10-22-retry-corrective-spec-required-ignora-recorded-paths.md
regeneration: []
promoted: []
---

# Recorded paths are read before carry-forward provenance

A carried Task's recorded paths are accepted by the real corrective-retry proof for both archive forms and by reconcile, while out-of-scope corrections remain refused. Reader and writer fixtures preserve heading ownership, exact bytes and provenance, and all five sabotages are detected. QA has a qualifying declared partial with no findings; the network-prohibited CommonMark source and absent Pull Request remain environment-blocked, and Oraculum's production retry after release remains unproven under the PRD declaration. The offline archive plan succeeded without changing the Spec; archive and Task settlement remain Daemon-owned.
