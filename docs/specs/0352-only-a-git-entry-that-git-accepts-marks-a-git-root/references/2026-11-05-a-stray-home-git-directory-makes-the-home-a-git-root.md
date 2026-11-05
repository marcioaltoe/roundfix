---
type: fix
status: promoted
created: 2026-11-05
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
reason: null
---

# MacBook evidence (reported by the maintainer, 2026-11-05 BRT, pasted terminal output)

Host `macbook`, Roundfix v0.131.0+ installed. `roundfix runs list` / `attach` (and strict commands) refused with
"worktree.location must not be inside the repository tree: ... is inside the Git root "/Users/marcio" (found "/Users/marcio/.git")".

```
❯ ls -lad ~/.git /Users/.git /.git
"/Users/.git": No such file or directory (os error 2)
"/.git": No such file or directory (os error 2)
drwxr-xr-x@ - marcio  7 Oct 13:30  /Users/marcio/.git
❯ git -C ~ log --oneline -5
fatal: not a git repository (or any of the parent directories): .git
❯ git -C ~ remote -v
fatal: not a git repository (or any of the parent directories): .git
❯ git -C ~ status --short | head -20
fatal: not a git repository (or any of the parent directories): .git
❯ ls -la ~/.git
drwxr-xr-x@ - marcio 10 Jun 13:19  gk
drwxr-xr-x@ - marcio  7 Oct 13:30  info
❯ mv ~/.git ~/.git.desativado-20261105
```

So `~/.git` was a plain directory holding only `gk/` (some Git tool's data) and `info/` — no `HEAD`, no `objects/`. Git does not treat it as a repository and keeps walking up; Roundfix's `findGitRoot` (internal/config/config.go:2166, `os.Stat(<dir>/.git) == nil`) does, and adopts the home directory as the Git root. The maintainer renamed the directory as a workaround; the tool that owns `gk/` may recreate `~/.git/gk`, which would bring the refusal back.
