---
schema: roundfix/archive-record/v1
spec: 0073-skill-versions-decoupled-from-the-binary
title: Skill versions decoupled from the binary
status: archived
created: "2026-08-03"
archived: "2026-08-06"
disposition: partial
source: docs/history/specs/0073-skill-versions-decoupled-from-the-binary
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_06
qa_report: qa-report-2026-08-06-01.md
qa_verdict: partial
unproven:
  - a maintainer grant naming `skills-lock.json` for the three external Go Skills, followed by a bounded Task that records their provenance without creating any setup-snapshot entry
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "132"
delivery_commit: 5646174b8f6ba24e4f15e5b97ad9bad34cb06a58
---

# Skill versions decoupled from the binary

Roundfix pins skill *content* rather than skill *compatibility*. Each setup snapshot carries a `treeDigest` per skill, the catalog digest is computed over those snapshots, and both characterization corpora embed observed digests inside recorded diagnostics. Any skill edit therefore changes what the binary claims, and a Roundfix release asserts a fact about content it does not version.
