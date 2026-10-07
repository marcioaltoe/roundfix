---
schema: roundfix/archive-record/v1
spec: 0050-doctor-skill-readiness-hardening
title: Doctor Skill Readiness Hardening
status: archived
created: ""
archived: "2026-07-26"
disposition: pass
source: docs/history/specs/0050-doctor-skill-readiness-hardening
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-26.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
---

# Doctor Skill Readiness Hardening

Spec 0036 made Repository Skill Set readiness a blocking Doctor Command result, but review found three gaps in that implementation. Repository-level symbolic links can redirect skill or lock reads beyond the intended Git root, the external hash order only approximates the installed skills CLI's `String.localeCompare` order, and the Doctor coordinator can substitute its process working directory when configuration did not resolve a Git root.
