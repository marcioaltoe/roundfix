# QA-06 — Pending branch

Status: pass.

The built CLI against the real Spec classified the Run Branch `pending`, named
six unrepresented commits, two branch-only files, and 22 differing shared
files, and emitted no reclaim command. Fresh JSON returned the same evidence;
refs and worktree registration remained unchanged.

`TestAuditClassifiesPendingBranch` also passed inside the fresh 12-test fixture
selection and requires evidence naming unrepresented commit and content state.
