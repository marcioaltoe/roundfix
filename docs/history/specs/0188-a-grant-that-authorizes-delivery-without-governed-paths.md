---
schema: roundfix/archive-record/v1
spec: 0188-a-grant-that-authorizes-delivery-without-governed-paths
title: A grant that authorizes delivery without Governed Paths
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: pass
source: docs/history/specs/0188-a-grant-that-authorizes-delivery-without-governed-paths
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_02
qa_report: qa-report-2026-09-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0179
sources:
  - 2026-09-30-a-spec-with-no-governed-path-cannot-record-a-grant.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "291"
delivery_commit: 4c2eb637e0bf870d260be18a08183c513c81b944
---

# A grant that authorizes delivery without Governed Paths

On 2026-09-30, `roundfix deliver plan` blocked Spec 0186 with `authorization refused: paths`. The Spec changes only ordinary files, and the authorization reader refuses a record whose `paths` has no entry. Listing the ordinary files instead failed the repository contract `TestEveryBoundedPathIsGoverned` in the Verification gate of PR #287 (run 36701160969), because every bounded path must be governed (ADR-0130). No record was valid, so a Spec that changes no Governed Path could only be delivered by hand.
