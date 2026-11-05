---
status: accepted
created_at: 2026-11-04T09:30:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A Read-Only Run Command needs no repository and only warns on worktree.location

The commands that only read the Run Database run from any directory. Outside
a Git repository `runs list` lists every repository, and a `worktree.location`
that other commands refuse is reported as one warning instead of failing them.

## Context

Every spec Run executes in its own git worktree under `worktree.location`,
`~/.roundfix/worktrees` by default, and the configuration loader refuses a
`worktree.location` inside the Git root of the working directory. That root is the first ancestor holding any `.git`, so a stray `.git` in the home
directory makes the whole home a repository and puts the default
`~/.roundfix/worktrees` inside it.

On 2026-11-03 the maintainer ran `roundfix runs list` and `roundfix attach`
from `~/Google Drive/My Drive`, which is not a repository he uses, and both
refused with `worktree.location must not be inside the repository tree:
"/Users/marcio/.roundfix/worktrees"`. Neither command creates a Run Worktree:
both read the Run Database under Roundfix Home. The refusal named neither the
Git root nor the `.git` that made it one. On 2026-11-04 he answered:

> runs list/attach funcionam fora de repositório (leem o store global), e a
> recusa do worktree.location nomeia a raiz Git detectada e o caminho do .git,
> para o usuário achar o .git perdido.

## Decision

A Read-Only Run Command is bare `runs`, `runs list`, `runs show`, `attach` or
`events`. It loads configuration with the worktree location optional: when
`worktree.location` cannot be resolved, the command prints the refusal once on
stderr as `roundfix: warning: <refusal>; commands that create a Run Worktree
refuse it`, leaves the location unset and goes on. Every other command keeps
refusing.

Outside a Git repository `runs list` behaves as `runs list --all`, with the
repository column and exit `0`. Inside a Git repository, including one that a
stray `.git` makes, its scope stays the repository.

The refusal keeps its text as a prefix and adds the Git root and the `.git`
path it found: `worktree.location must not be inside the repository tree:
"<location>" is inside the Git root "<root>" (found "<root>/.git")`.

How the Git root is found does not change: any `.git` entry, directory or
file, still marks a repository.

**Amended by [ADR-0362](0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md) (2026-11-05).** A `.git` entry marks a Git root only when Git accepts it: a directory that passes Git's repository test or a file starting with `gitdir:`. A stray `.git` that Git rejects is skipped. The refusal text and the Read-Only Run Command contract still hold.

## Considered Options

- Keep refusing `runs list` outside a repository: rejected because the
  maintainer asked for it to work there, and outside a repository there is no
  narrower scope than every repository.
- Skip the worktree-location check silently: rejected because a stray `.git`
  would then scope `runs list` to the home directory and print `No Runs found.`
  with nothing pointing at the cause.
- Make every command lenient: rejected because `implement`, `deliver`,
  `reconcile`, `stop` and `gc` create, inspect or remove Run Worktrees under
  that location.
- Recognize a `.git` only when Git would, or refuse one at the home directory:
  rejected in this change because it alters repository detection for every
  command, and a `.git` file of a linked worktree or submodule must keep
  marking a repository. ADR-0362 takes the first half of this option,
  recognizing a `.git` only when Git would.
- Include `runs causes`: rejected because it reads the repository's own
  Archive Records and keeps requiring one.

## Consequences

The documented exit `2` of `runs list` outside a repository becomes exit `0`
with every repository listed. Under a stray `.git` the read-only commands
succeed and print the warning first on stderr; the commands that need the
location still fail, now naming the root and the `.git` to inspect.
