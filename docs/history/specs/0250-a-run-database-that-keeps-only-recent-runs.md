---
schema: roundfix/archive-record/v1
spec: 0250-a-run-database-that-keeps-only-recent-runs
title: A Run Database that keeps only recent Runs
status: archived
created: "2026-10-08"
archived: "2026-10-08"
disposition: pass
source: docs/specs/0250-a-run-database-that-keeps-only-recent-runs
source_revision: 27ea7cceff0a45ee7f77e8b2fd83b3e07ce99a4f
qa_task: task_05
qa_report: qa-report-2026-10-08.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0255
sources: []
regeneration:
  - command: make skills-sync
promoted: []
---

# A Run Database that keeps only recent Runs

The machine-wide Run Database grows with every Run and never sheds a Run. On 2026-10-08 it was 1.1 GB, of which `run_events` held 1.08 GB for 832k events written in the two weeks before, about 77 MB a day, across 313 Runs. Journal Retention prunes the journal of an old terminal Run but keeps its row and the rows that hang from it. ADR-0033: "Active Runs and their journals are never pruned". Freed pages go back to the filesystem only when someone runs `roundfix gc compact --apply` by hand, and `gc compact` reclaimed 5 MB because the space was live data. The maintainer decided on 2026-10-08, "Automática + gc (Recommended)", that Roundfix keeps only the data of Runs still running, and of any Run a Delivery Queue still references, plus the last N days. The operator then gets a database that stays bounded without a manual step, and `roundfix gc` shows what it would remove before it removes it.
