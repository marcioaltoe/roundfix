# QA-06 — Pending branch

Status: pass.

The public CLI against the current real Spec classified
`roundfix/run-run_20260805T013130Z_31398370e8ba8670` as `pending`. Its evidence
named 10 commits not represented on `origin/HEAD`, 22 branch-only files, and 24
differing shared files. Text and fresh JSON agreed, and neither emitted a
reclaim command.

`TestAuditClassifiesPendingBranch` also passed in the fresh named fixture
selection. Refs and worktree registration remained identical across the public
invocations.
