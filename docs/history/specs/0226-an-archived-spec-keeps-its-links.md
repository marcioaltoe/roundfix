---
schema: roundfix/archive-record/v1
spec: 0226-an-archived-spec-keeps-its-links
title: An archived Spec keeps its links
status: archived
created: "2026-10-04"
archived: "2026-10-04"
disposition: pass
source: docs/history/specs/0226-an-archived-spec-keeps-its-links
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-04.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0230
sources:
  - 2026-10-03-archive-does-not-rewrite-relative-links.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "379"
delivery_commit: 69ace45de96d7323ae9c10729670ce94d568292c
---

# An archived Spec keeps its links

`roundfix archive` moves a Spec one directory deeper, from the Spec Root to its archive root. A relative Markdown link in the Spec that leaves it, such as a PRD's link to a Finding or a Task's link to an ADR, then resolves one level short and breaks. fluxus measured it on 2026-10-02: Spec 0096 reached `main` with links that can no longer be fixed in place, because an archived Spec is immutable, and Specs 0097 and 0098 needed three links fixed by hand in their archive commits. Roundfix's own history carries the same damage. This Spec makes the archive keep every outward link reaching the file it reached, and refuse, naming the link, when one reaches nothing.
