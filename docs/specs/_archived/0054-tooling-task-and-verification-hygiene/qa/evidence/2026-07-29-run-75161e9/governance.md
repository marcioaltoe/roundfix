# Governance evidence

Build: `75161e9c3a5f7554cd1e0b9290bce6c61820b5c7`.

## Project Constraints and Task completion

- `_prd.md` and `_techspec.md` each contain all four required Project
  Constraint rows. Identifier strategy and authentication/HTTP are
  inapplicable with reasons. Active ADR obligations and tooling authority are
  applicable with operative sources under `docs/agents/`.
- The tooling rows match byte-for-byte on their exact bounded protected paths:
  `Makefile`, `.gitignore`, both Roundfix Skill copies, all three setup
  snapshots, the normalized catalog and digest, and both parity-corpus files.
- `docs/agents/domain.md`, `docs/agents/cli.md`, and
  `docs/agents/agent-instructions.md` exist and contain the cited contracts.
- `_tasks.md` declares five Tasks. Every `task_01.md` through `task_05.md` is
  present with `status: completed` and a Result section containing named,
  reproducible evidence for its acceptance criteria.

## Protected-tooling history audit

Fresh commands:

```text
rtk proxy git -c core.fsmonitor=false log --reverse --format=... --name-only origin/main..HEAD
rtk proxy git -c core.fsmonitor=false diff-tree --no-commit-id --name-only -r <commit>
rtk proxy git -c core.fsmonitor=false status --porcelain=v1 --untracked-files=all
```

One-pass result:

- `397227f` authored the original PRD/TechSpec authorization before every
  Task commit.
- `f022ae0` (`task_02`) changes exactly `.gitignore`, `Makefile`, and
  `task_02.md`.
- `8e6c707` adds `go-cli.json` and `rust-cli.json` to both authorization
  records before `b6b59a5` regenerates them.
- `b6b59a5` contains exactly the seven authorized generated setup,
  catalog, and parity paths in its own reconciliation commit.
- `c02feaa` (`task_05`) contains the two authorized Skill copies, the five
  target-produced Skill fallout paths, repository-authored guidance, and its
  own Task file.
- Consequent repairs are separate and follow their causes:
  `980e2d1` follows `task_04`; `fa4c5ef`, `a9c1d3c`, and `9b8b7b9` follow
  `task_05`; `a9c1d3c` and `75161e9` keep Makefile repairs separate after
  `task_02`.
- No prerequisite fix or authorization is folded into a Daemon-owned Task
  commit. No protected current-worktree delta exists; only this Spec's QA
  directory is untracked.

No missing, late, untraceable, or out-of-scope protected-tooling mutation was
found on the current history.
