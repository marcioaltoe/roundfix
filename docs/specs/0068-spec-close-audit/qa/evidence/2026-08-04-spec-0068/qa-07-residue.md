# QA-07 — Local and remote residue branches

Status: pass.

The built CLI against `/private/tmp/roundfix-qa0068-scratch.D31PJZ` classified
local `ma/scratch-close` and remote `origin/ma/scratch-close` as `residue`,
with content fully represented on `main`. It printed, but did not run:

```text
git branch -d -- 'ma/scratch-close'
git push --delete 'origin' 'ma/scratch-close'
```

Text and JSON exited 1. Fresh branch and worktree reads after repeat execution
proved all surfaces remained. The local-residue and motivating-session remote
backup fixtures also passed in the 12-test assembled selection.
