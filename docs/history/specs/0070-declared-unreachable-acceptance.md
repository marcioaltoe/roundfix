---
schema: roundfix/archive-record/v1
spec: 0070-declared-unreachable-acceptance
title: Declared unreachable acceptance
status: archived
created: "2026-08-02"
archived: "2026-08-04"
disposition: pass
source: docs/history/specs/0070-declared-unreachable-acceptance
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_06
qa_report: qa-report-2026-08-04-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "111"
delivery_commit: b6ecb00a24c5dfd500facebf567921d3fc906b42
---

# Declared unreachable acceptance

Some acceptance genuinely cannot be reached by any hermetic Verification. Spec 0058's remaining QA row needs a real tagged release published against six live npm trusted-publisher bindings — an irreversible act against a live registry that no gate may perform. Its QA closed `partial` with **zero findings and zero failed rows**, and the Spec still could not archive, because `roundfix archive` requires `verdict: pass`.
