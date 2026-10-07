---
schema: roundfix/archive-record/v1
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
title: Owned skills and a release step that follow the bundle
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: pass
source: docs/history/specs/0195-owned-skills-and-a-release-step-that-follow-the-bundle
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0189
sources:
  - 2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md
  - 2026-09-30-a-release-does-not-check-skills-and-guides.md
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "299"
delivery_commit: cad8bb25e36b3e32db042900e797252ae059156f
---

# Owned skills and a release step that follow the bundle

Roundfix ships 14 owned skills inside its binary and compares each installed copy with a minimum version. The minimum is one literal for all of them, and the bundle is already ahead of it. A repository can therefore run an older `qa-gate` than the binary carries while Doctor reports `skills: ok` and the `baseline update` preview reports `current`. A skill's content can also change while its version stays the same, so the comparison says nothing about content. Two setups omit the Roundfix skill, and no profile tells an Agent when to load it. No release re-reads skills or guides before it ships.
