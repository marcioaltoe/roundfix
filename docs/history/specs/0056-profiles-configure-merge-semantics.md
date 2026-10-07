---
schema: roundfix/archive-record/v1
spec: 0056-profiles-configure-merge-semantics
title: Profiles configure merge semantics
status: archived
created: "2026-07-28"
archived: "2026-08-02"
disposition: pass
source: docs/history/specs/0056-profiles-configure-merge-semantics
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-08-01-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0049
sources: []
regeneration: []
promoted: []
pull_request: "68"
delivery_commit: 88550ea2c3c83fd843a6800e8c41e7d4aea6731e
---

# Profiles configure merge semantics

`roundfix profiles configure --scope project --file <fragment>` replaces the entire `profiles:` map with the fragment, so configuring one category silently deletes every other configured profile — a fragment naming only `frontend` destroyed `general`, `backend`, `qa`, and `review` with their Fallback Chains, reported `changed: true` naming only the requested category, and rewrote the file's indentation so the 103-line diff read as formatting. This contradicts the command's documented contract of preserving unrelated config. A second defect compounds it: a declined confirmation in a non-interactive context writes nothing and exits `0`, so automation reads a refusal as success. Evidence: [profiles configure replaces the whole profiles map](../../findings/2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md).
