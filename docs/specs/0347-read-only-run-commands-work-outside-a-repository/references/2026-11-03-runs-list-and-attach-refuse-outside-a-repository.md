---
type: fix
status: promoted
created: 2026-11-03
spec: 0347-read-only-run-commands-work-outside-a-repository
reason: null
---

# Maintainer report, 2026-11-03 (BRT), Roundfix 0.129.x/0.130.x on the MacBook

Marcio ran, from `~/Google Drive/My Drive` (not a Git repository of his):
- `roundfix runs list`
- `roundfix attach`
Both refused with: `worktree.location must not be inside the repository tree: "/Users/marcio/.roundfix/worktrees"`.
He asked whether this was an error/problem.

Operator diagnosis (this repository, not measured on the MacBook): `findGitRoot` (internal/config/config.go, `func findGitRoot`) walks up from the working directory and accepts the first ancestor holding any `.git` (os.Stat). The likely cause is a stray `.git` in `/Users/marcio` (or `/Users`, `/`) on the MacBook, which makes the home the "repository tree"; `ResolveWorktreeLocation` (same file) then refuses the default worktree location `~/.roundfix/worktrees` because it lies inside that tree. Gaps: the refusal does not name the Git root it detected nor where the `.git` is, and `runs list` / `attach` read the global Run store, so they should not need the current directory's repository at all.
