---
schema: roundfix/archive-record/v1
spec: 0007-command-lifecycle
title: Command Lifecycle
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0007-command-lifecycle
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0022
  - ADR-0021
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# Command Lifecycle

Roundfix now implements Specs and cleans reviews end to end, but living with it day to day still involves ceremony the tool should own: installing the pinned acpx by hand (and discovering the local-adapter override the hard way), no way to update an installed roundfix or learn a newer one exists, stopping a live Run only from its own terminal with Ctrl-C, and spec Runs that never push even when a repository wants its Clean branches published. This Spec ships the four lifecycle pieces the first dogfood round asked for: setup, upgrade, graceful stop, and configurable spec-Run push.
