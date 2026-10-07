---
schema: roundfix/archive-record/v1
spec: 0012-npm-distribution
title: npm Distribution and Skill Bundle
status: archived
created: "2026-07-06"
archived: "2026-07-06"
disposition: qa-override
source: docs/history/specs/0012-npm-distribution
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: ""
qa_verdict: ""
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0031
sources: []
regeneration: []
promoted: []
pull_request: "18"
delivery_commit: 5929aea32a0733f4329165cb8ea55afa4e6a9276
---

# npm Distribution and Skill Bundle

Roundfix has no public install path. A user who wants the CLI must clone the repository and build it, and the Upgrade Command can only pull release assets that no workflow produces yet. Node is already a hard prerequisite of the agent layer, so the shortest install every target platform shares is npm. This Spec ships Roundfix through npm as a launcher package with per-platform binary packages, driven by a tag-triggered release workflow that cross-compiles the Go binary for every supported platform and publishes both the npm packages and the GitHub Release assets the Upgrade Command already consumes. The distributed binary also carries the Roundfix skill bundle — the operational Roundfix Skill plus the authorial workflow skills the repository owns — so an installed Roundfix can teach an Agent runtime not just how to operate the CLI but the whole spec-driven method it drives.
