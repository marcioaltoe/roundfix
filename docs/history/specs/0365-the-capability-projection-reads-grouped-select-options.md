---
schema: roundfix/archive-record/v1
spec: 0365-the-capability-projection-reads-grouped-select-options
title: The capability projection reads grouped select options
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: partial
source: docs/specs/0365-the-capability-projection-reads-grouped-select-options
source_revision: 91120d1a57d7bceef2a06a6bc5b76a3821a26e43
qa_task: task_03
qa_report: qa-report-2026-11-09.md
qa_verdict: partial
unproven:
  - the operator's `roundfix profiles validate` run with codex-acp 2.6.0 after the release, recorded session evidence.
adrs:
  - ADR-0375
  - ADR-0293
sources:
  - 2026-11-09-codex-acp-2-6-0-grouped-select-options-break-the-capability-proof.md
  - codex-acp-2.6.0-config-options.json
  - codex-acp-2.5.0-config-options.json
regeneration: []
promoted: []
pull_request: "789"
type: fix
---

# The capability projection reads grouped select options

The recorded codex-acp 2.6.0 grouped payload now projects the five mode values in advertised order, and the fake ACPX proof succeeds with and without full access while preserving the 2.5.0 behavior. Full access selects agent-full-access, invalid shapes retain their named refusals, and six independent sabotages are detected. The operator must still validate every Codex tuple with the released patch and real codex-acp 2.6.0; this report leaves that outcome unproven and the pre-PR control blocked.
