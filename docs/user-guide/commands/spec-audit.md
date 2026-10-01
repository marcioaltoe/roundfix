### spec audit

```bash
roundfix spec audit <slug> [--format <text|json>]
```

Audits one active or archived Spec, reading its delivery state and surviving
Git branches and worktrees. The report lists every surviving branch and
worktree with its classification evidence and gives residue an exact reclaim
command. It never changes Git state, the Run Database, or Spec artifacts.
`--format` selects text (the default) or JSON output.

Exit codes:

- `0` — no residue or undelivered work.
- `1` — residue or undelivered work found, or the audit could not run.
- `2` — usage error or unknown Spec slug.

