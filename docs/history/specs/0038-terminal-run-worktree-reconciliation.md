---
schema: roundfix/archive-record/v1
spec: 0038-terminal-run-worktree-reconciliation
title: Terminal Run Worktree reconciliation
status: archived
created: "2026-07-17"
archived: "2026-07-27"
disposition: pass
source: docs/history/specs/0038-terminal-run-worktree-reconciliation
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-27.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0053
sources: []
regeneration: []
promoted: []
pull_request: "39"
delivery_commit: 146faa1456d5f9827c7d8c48041641525fa52b4e
---

# Terminal Run Worktree reconciliation

Terminal spec Runs can retain clean Run Worktrees and Run Branches after their commits have already reached the user's target branch. Roundfix currently compares the Run Branch only with its creation base, hides the residue behind the default Active Run listing, and offers no supported cleanup command. Prior dogfood evidence, now absorbed into this Spec and retained in Git history, showed users having to prove ancestry and run Git commands manually.
