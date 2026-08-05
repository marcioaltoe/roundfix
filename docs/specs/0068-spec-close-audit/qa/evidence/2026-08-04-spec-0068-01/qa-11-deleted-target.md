# QA-11 — Deleted-target content resolution

Status: pass.

The built `roundfix reconcile` text and JSON entry points both exited 0 in
read-only mode. This checkout had no matching retained Run in the
Agent-visible database, so the exact deleted-target cases used disposable
real-Git fixtures.

Fresh command:

```text
rtk proxy go test ./internal/worktree -count=1 -run '^TestInspectTerminalRun(SafeWhenTargetDeletedAfterSquashMerge|UnintegratedWhenDeletedTargetHasRunOnlyFile|UnintegratedWhenDeletedTargetHasDifferentSharedFile|UnintegratedWhenDeletedTargetContentComparisonFails|UnknownWhenDeletedTargetDefaultBranchCannotBeResolved|ActiveRunDoesNotInspectDeletedTargetContent)$' -v
```

Exit 0; all six named cases passed. Full representation resolves `safe`;
Run-only or differing shared content resolves `unintegrated`; comparison
failure preserves; an unresolved default stays `unknown`; the reasons name
the evidence; and Active Runs bypass the path.
