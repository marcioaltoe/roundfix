---
schema: roundfix/archive-record/v1
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
title: Review findings with evidence, and no unselected providers
status: archived
created: "2026-09-28"
archived: "2026-09-29"
disposition: partial
source: docs/history/specs/0179-review-findings-with-evidence-and-no-unselected-providers
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-29-04.md
qa_verdict: partial
unproven:
  - the task_02 tests (every ledger field and each refusal) executed against the built tree, and the first `roundfix review dispose` run on a real findings record after this Spec merges
  - the task_03 tests executed against the built tree, and the first reuse of a real findings record after this Spec merges
  - the captured bodies checked against the dispositions the ledger format can express, and a live replay after this Spec merges
  - the repository checks and the pre-PR review recorded on the pull request at publication, before merge
adrs:
  - ADR-0165
sources:
  - 2026-09-28-unselected-review-providers-receive-no-request.md
  - 2026-09-28-review-findings-get-evidence-backed-dispositions.md
  - 2026-09-28-blocking-review-after-archive-parks-publication.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "274"
delivery_commit: f9c5fcce00daa0dfe36139f6c97e8938242fa6e0
---

# Review findings with evidence, and no unselected providers

A pre-PR review finding today ends in the operator's hands, and its disposition lives nowhere Roundfix can read. On 2026-09-25 and 2026-09-28 the Codex pre-PR review raised findings on Specs 0170, 0171 and 0172. The operator fixed some with corrective Tasks and TechSpec edits, repaired a QA Report's front matter by hand, and dismissed two findings only in pull request bodies: the Windows portability finding on #258 and the stale-binary finding on #259. No per-finding record exists, and four gaps follow from that:
