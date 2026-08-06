# Static gate

Build: `7ed6af3cac80c4cdc051c91815b08b5c57b642f6`

Command, run unpiped from the Run Worktree root:

```text
rtk make verify
exit 0
```

Observed summary:

```text
Go test: 3521 passed in 26 packages
Go test: 1 passed in 1 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
go build -buildvcs=false ... ./cmd/roundfix
Spec 0073-skill-versions-decoupled-from-the-binary
No findings.
```

The Spec Consistency Check skipped only the absent optional Vocabulary Contract
section and absent `references/_index.md`; neither is required by this Spec or
used by an acceptance path.

After the report closed, `rtk ./bin/roundfix spec check` exited 0 and again
reported `No findings.` for Spec 0073 with the same two non-acceptance skips.
