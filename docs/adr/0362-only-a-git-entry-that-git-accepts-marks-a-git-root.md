---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# Only a .git entry that Git accepts marks a Git root

Roundfix finds the Git root of the working directory by walking up to the
first ancestor whose `.git` entry Git itself would accept. A `.git` that Git
rejects is skipped and the walk goes on to the parent, as Git's own discovery
does.

## Context

Configuration resolution (`config.Load`, `roundfix init --scope project`,
`roundfix profiles configure --scope project`) and the project scope of `skills
install` found the Git root by taking the first ancestor holding any `.git`
entry. On 2026-11-05 the maintainer's MacBook had a `~/.git` directory that
held only `gk/` and `info/`, left by another tool, with no `HEAD` and no
`objects/`. `git -C ~ status` answered `fatal: not a git repository`, but
Roundfix took the home directory for a repository: every command that creates
a Run Worktree refused the default `~/.roundfix/worktrees` as inside the
repository tree, and `roundfix init --scope project` wrote `.roundfixrc.yml`
into the home directory. ADR-0357 made the refusal name the root and its
`.git`, and kept the detection unchanged. The maintainer renamed the
directory, and the tool that owns `gk/` may create it again. On 2026-11-05 he
answered:

> findGitRoot passa a aceitar só um .git que o Git reconhece: pasta com HEAD
> e objects/, ou arquivo 'gitdir:'. Um teste cobre um ~/.git falso.

## Decision

A `.git` entry marks a Git root when it is one of these:

- A directory that passes Git's own repository test (`is_git_directory` in
  Git's `setup.c`): its `HEAD` is a symbolic link whose target starts with
  `refs/`, or a file that starts with `ref:` followed by optional whitespace
  and `refs/`, or with a hexadecimal object name of 40 or 64 digits; and its
  common directory holds an `objects` directory and a `refs` directory. The
  common directory is the path in the entry's `commondir` file, relative to
  the entry when not absolute, or the entry itself when it has none.
- A regular file whose content starts with `gitdir:`, as a linked worktree or
  a submodule writes it. The file marks the root whether or not its target
  exists, because Git stops its own walk at that file too and reports the
  broken link there.

Any other `.git` entry is skipped and the walk continues with the parent: a
directory that fails the test, a file without `gitdir:`, and an entry that is
neither a file nor a directory. The rule is read from the file system, with
no `git` process, so loading configuration stays fast and needs no `git`
binary. Configuration resolution and the project scope of `skills install`
use this one walk. `GIT_DIR`, `GIT_COMMON_DIR`, `GIT_OBJECT_DIRECTORY`,
`GIT_CEILING_DIRECTORIES`, bare repositories and `safe.directory` stay
outside it, as before.

## Considered Options

- Ask Git with `git rev-parse --show-toplevel`: rejected because every
  configuration load would start a process and need a `git` binary, and the
  commands that read only the Run Database must keep working without one.
- Accept a directory with `HEAD` and `objects/` only, as the maintainer's
  option text says: narrowed to Git's exact test, which also requires `refs/`
  and a `HEAD` that names a ref or an object, so Roundfix never takes for a
  repository a directory that Git rejects.
- Stop the walk at any regular `.git` file, as Git does when the file has no
  `gitdir:` line and fails with `invalid gitfile format`: rejected because the
  approved decision names the `gitdir:` file, and a file without it is not a
  repository Git can use; Roundfix skips it and goes on to the parent.
- Skip a `gitdir:` file whose target is missing: rejected because Git stops
  there, and walking past a broken linked worktree could adopt an outer
  repository's configuration.
- Refuse a `.git` in the home directory only: rejected because a home
  directory can be a real repository, such as a dotfiles checkout, and a stray
  `.git` can sit in any ancestor.

## Consequences

Under a `.git` that Git rejects, Roundfix finds no Git root there: commands
that create a Run Worktree accept the default location, Read-Only Run
Commands print no warning, and `roundfix init --scope project` refuses with
`project init requires a Git root; use --scope user outside a repository`.
A repository Git accepts, a linked worktree and a submodule resolve as
before. Tests that faked a repository with an empty `.git` directory create
the smallest `.git` directory Git accepts instead. This decision amends
ADR-0357, whose Decision kept any `.git` entry as a marker.
