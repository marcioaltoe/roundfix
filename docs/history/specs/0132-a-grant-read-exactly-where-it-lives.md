---
schema: roundfix/archive-record/v1
spec: 0132-a-grant-read-exactly-where-it-lives
title: A grant read exactly where it lives
status: archived
created: "2026-09-10"
archived: "2026-09-13"
disposition: pass
source: docs/history/specs/0132-a-grant-read-exactly-where-it-lives
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-09-11-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A grant read exactly where it lives

Spec 0119 made the authorization record the thing that decides whether governed work may proceed. Its pre-Pull-Request review found that the reader accepts a record it should refuse, and refuses records it should accept. A malformed frontmatter delimiter still yields a grant, and a Spec Root configured outside the code repository is unreadable because the record path is assembled from a constant instead of the resolved root. Two of 0119's own tests also pin the record to its active path, so archiving the Spec breaks them.
