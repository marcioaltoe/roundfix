---
schema: roundfix/archive-record/v1
spec: 0161-deliver-ready-for-real-repositories
title: Deliver, ready for real repositories
status: archived
created: "2026-09-24"
archived: "2026-09-24"
disposition: pass
source: docs/history/specs/0161-deliver-ready-for-real-repositories
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-24-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "242"
delivery_commit: 57f430549a707daba35f906e9b66492dfe842a57
---

# Deliver, ready for real repositories

Spec 0156 shipped `roundfix deliver` with six recorded limits after its corrective ceiling was spent. It also stated that no release may ship the command before this Spec merges. Four of the limits make the loop park or publish wrongly on an ordinary repository; two make it unrecoverable after a crash or a reboot.
