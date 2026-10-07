---
schema: roundfix/archive-record/v1
spec: 0239-a-glossary-every-spec-keeps-current
title: A glossary every Spec keeps current
status: archived
created: "2026-10-06"
archived: "2026-10-06"
disposition: partial
source: docs/history/specs/0239-a-glossary-every-spec-keeps-current
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-06.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0244
sources:
  - 2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
promoted: []
pull_request: "431"
delivery_commit: 7b766babf077963bd53a949e81f62b74054c6ff2
---

# A glossary every Spec keeps current

Roundfix was built on the CONTEXT-driven method: the glossary (`CONTEXT.md`) and the ADRs are the base, and every Spec, Task and line of copy draws its vocabulary from them. On 2026-10-06 the maintainer observed that the loop no longer keeps the glossary current: it last changed on 2026-10-02, twenty-one Specs archived after it without touching it, and authoring and QA reports named new terms that nobody wrote down ([the adopted Backlog Entry](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)). The domain guide already asks for a glossary check at the close of a Spec, but no gate reads the answer, the autonomous authoring route never reaches `domain-modeling`, and the operator's briefing had told authors not to edit the glossary.
