---
schema: roundfix/archive-record/v1
spec: 0036-doctor-skill-readiness
title: Doctor Skill Readiness
status: archived
created: "2026-07-17"
archived: "2026-07-26"
disposition: pass
source: docs/history/specs/0036-doctor-skill-readiness
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-26.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "37"
delivery_commit: 95a2c0d25b22db019507bee36cceb2218844a2f4
---

# Doctor Skill Readiness

Roundfix depends on a repository-local set of Agent Skills, but the Doctor Command currently validates only the runtime toolchain. A repository can pass Doctor while a Roundfix-owned skill is stale relative to the running binary, an externally managed skill is missing, or installed external content no longer matches `skills-lock.json`. That creates a dangerous split: the CLI executes one contract while the Supervisor or Agent follows another.
