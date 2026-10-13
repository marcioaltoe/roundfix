---
type: fix
status: done
created: 2026-10-12
spec: 0259-proven-signals-and-a-shared-build-cache
reason: "Spec 0259's corrective task_06 fixes both disposal lock tests: unexported runner-taking removeRegisteredWorktree and addDetachedWorktree, tests that cancel once the Git common directory is resolved and require context.Canceled, no context.WithTimeout left in disposal_test.go, and 50 passes with -race -parallel 16. The same flake had stopped Spec 0259's QA ten-suite measurement. Spec 0262 dropped the work it had planned for this entry."
---

# The worktree admin-lock test fails under load

`TestRemoveRegisteredWorktreeWaitsForTheAdminLock` failed the repository gate of
the Delivery Queue item for Spec 0256 on 2026-10-12, then passed three times out
of three when run alone and passed on the retry.

Full text in Git at `94fbf85b43eb2d28574fd052351d7ed5aba6d832`: `docs/backlog/2026-10-12-the-worktree-admin-lock-test-is-timing-dependent.md`.
