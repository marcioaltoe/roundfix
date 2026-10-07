---
schema: roundfix/archive-record/v1
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
title: A repository says who reviews before the Pull Request
status: archived
created: "2026-09-17"
archived: "2026-09-18"
disposition: pass
source: docs/history/specs/0143-a-repository-says-who-reviews-before-the-pull-request
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-18.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "206"
delivery_commit: 46cae32285cc3551d8af63e5f5e5e6e8581e3bce
---

# A repository says who reviews before the Pull Request

Every delivery in this repository is supposed to pass an independent review before its Pull Request opens. The obligation is written as a mandatory clause in the agent guides, which name the supported choices — `codex`, `claude`, `coderabbit`, or an explicit `none` — and say the choice must be resolved before publication.
