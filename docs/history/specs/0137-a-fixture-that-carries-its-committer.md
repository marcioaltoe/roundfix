---
schema: roundfix/archive-record/v1
spec: 0137-a-fixture-that-carries-its-committer
title: A fixture that carries its committer
status: archived
created: "2026-09-14"
archived: "2026-09-14"
disposition: pass
source: docs/history/specs/0137-a-fixture-that-carries-its-committer
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_02
qa_report: qa-report-2026-09-14.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A fixture that carries its committer

Spec 0133 replaced the Daemon Task-cycle fixture's per-test repository setup with one shared seed copied per test. The setup it replaced wrote a committer identity into the repository's own config; the seed does not. Repository hardening writes only maintenance keys, and the identity the test helper passes is a per-invocation override that no copy inherits.
