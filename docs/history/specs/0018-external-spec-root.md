---
schema: roundfix/archive-record/v1
spec: 0018-external-spec-root
title: External Spec Root
status: archived
created: "2026-07-07"
archived: "2026-07-07"
disposition: pass
source: docs/history/specs/0018-external-spec-root
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-07.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0035
sources: []
regeneration: []
promoted: []
---

# External Spec Root

Roundfix assumes Spec artifacts live at `docs/specs/` inside the repository working tree. Repositories that keep planning artifacts in a knowledge workspace — a separate git repository nested in or beside the code repository, reached through a symlink — break that assumption twice: the per-Task commit stages the task file through a path that crosses a symbolic link, which git refuses (`pathspec beyond a symbolic link`), and Run Worktrees materialize the symlink without its target, so a relative link dangles. External Spec Root makes the Spec artifact location a first-class configuration, so Roundfix reads and writes Specs wherever they live and never stages external artifacts into code-repository commits.
