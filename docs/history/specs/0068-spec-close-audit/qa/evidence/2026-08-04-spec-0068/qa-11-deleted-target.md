# QA-11 — Deleted-target content resolution

Status: pass.

Fresh focused command exited 0 with six tests passed:

`rtk go test ./internal/worktree -count=1 -run
'^TestInspectTerminalRun(SafeWhenTargetDeletedAfterSquashMerge|UnintegratedWhenDeletedTargetHasRunOnlyFile|UnintegratedWhenDeletedTargetHasDifferentSharedFile|UnintegratedWhenDeletedTargetContentComparisonFails|UnknownWhenDeletedTargetDefaultBranchCannotBeResolved|ActiveRunDoesNotInspectDeletedTargetContent)$' -v`

It proves full representation resolves `safe`; Run-only or different shared
files remain `unintegrated`; failure preserves; unresolved default stays
`unknown`; reasons carry evidence; and Active Runs bypass the path. The full
worktree and CLI packages then passed 1,041 tests, retaining existing refusals.
