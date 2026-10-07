---
schema: roundfix/archive-record/v1
spec: 0034-release-plan
title: Release Plan
status: archived
created: "2026-07-16"
archived: "2026-07-17"
disposition: pass
source: docs/history/specs/0034-release-plan
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-17.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0048
sources: []
regeneration: []
promoted: []
pull_request: "32"
delivery_commit: d7fdd8cf3f986bda0cbd3683e239c799dcf6b41a
---

# Release Plan

Roundfix validates release tags and publishes matching artifacts, but it does not determine which semantic-version increment the committed changes require or identify when that decision needs explicit human approval. The Release Plan gives maintainers and Agents one deterministic, read-only assessment before any changelog, tag, push, package, or GitHub Release mutation.
