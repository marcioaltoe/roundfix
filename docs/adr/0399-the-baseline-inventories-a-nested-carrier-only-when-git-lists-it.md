---
status: accepted
created_at: 2026-11-13T14:00:00Z
updated_at: 2026-11-13T14:00:00Z
deprecated_at: null
superseded_by: null
---

# The Baseline inventories a nested carrier only when Git lists it

A nested `AGENTS.md` or `CLAUDE.md`, one below the repository root and
outside `docs/agents/`, enters the Baseline inventory only when
`git ls-files -z --cached --others --exclude-standard` at the repository root
lists it: a tracked file, or an untracked file Git does not ignore. Git
decides, so a copy in an ignored path and a carrier inside a linked worktree,
a nested repository or a submodule stay out of the inventory, raise no
`baseline.inventory.nested-carrier-conflict` and do not enter the Plan
Digest. Root carriers, guides under `docs/agents/` and the paths the walk
already skips keep their behavior.

## Context

On 2026-11-13 quaestor reported that `roundfix baseline update` (0.143.0)
warned `baseline.inventory.nested-carrier-conflict` for the `AGENTS.md` and
`CLAUDE.md` of `.storage/simulacao-deploy/wt-aee32c4/`, a linked worktree of
the same repository that a local script creates under a directory its
`.gitignore` lists. Removing the worktree removed both warnings. The inventory
walks the working tree and skips only `.git`, `.hg`, `.svn`, `.venv`,
`__pycache__`, `node_modules`, `vendor`, `venv`, `skills/` and
`.agents/skills/`; it never asked Git which files belong to the repository.

ADR-0070 says the Baseline "inventories every bounded instruction and
agent-document carrier". A carrier Git ignores, or one another checkout owns,
is not a file of this repository: the Baseline never writes it, no commit of
this repository holds it, and its warning asks the operator to review a
conflict that does not exist.

Characterized on 2026-11-13 with Git 2.56.1 in a temporary repository:
`git check-ignore` fails with exit 128 and `fatal: Pathspec ... is in
submodule` for a path inside a submodule, so asking Git path by path about
ignore rules is not enough; `git ls-files --cached --others
--exclude-standard` lists the tracked and the untracked, not ignored files and
stops at a linked worktree, a nested repository and a submodule. A directory
whose `.git` file points to a repository that no longer exists is ordinary
content to Git, and Git lists its files.

## Decision

- **Git is the authority.** After the walk, when it met at least one nested
  `AGENTS.md` or `CLAUDE.md`, the inventory runs, once, at the repository
  root: `git ls-files -z --cached --others --exclude-standard --
  ':(glob)**/AGENTS.md' ':(glob)**/CLAUDE.md'`. A nested carrier the output
  lists is inventoried exactly as before; any other nested carrier is left
  out entirely: no preimage, no carrier, no warning, no bytes counted against
  the inventory limits. The output is read byte for byte, split on NUL.
- **The rule covers nested `AGENTS.md` and `CLAUDE.md` only.** Root carriers
  are always inventoried and stay mutation targets; guides under
  `docs/agents/` keep the classification of ADR-0372.
- **A failed listing stops the inspection.** The error names
  `list repository instruction carriers` and the Git failure; the Baseline
  commands report it like any Git failure. No carrier is guessed in or out.
- **Nothing is reported for a carrier left out**, as nothing is reported for
  `node_modules` or `vendor` today. A repository whose nested carriers Git
  all lists keeps the same snapshot and the same Plan Digest.

## Considered Options

- **Skip the root of any directory that holds a `.git` entry, then ask
  `git check-ignore` about the rest.** Rejected: two rules where Git already
  has one, and the filesystem rule would also drop a directory whose stale
  `.git` file Git treats as ordinary content, which `git status` still shows
  as untracked.
- **`git check-ignore` alone.** Rejected: it fails with exit 128 for a path
  inside a submodule, and it answers nothing about a linked worktree or a
  nested repository in a path that is not ignored.
- **Parse `.gitignore` in Go.** Rejected: it would miss
  `$GIT_COMMON_DIR/info/exclude`, `core.excludesFile` and Git's rule that a
  file under an excluded directory cannot be re-included, and it would drift
  from Git.
- **Report a count of the carriers left out.** Rejected: the walk reports
  none of the trees it already skips, and a count line would change the
  output and the Plan Digest of every repository that keeps ignored copies,
  for information `git status --ignored` already gives.
- **Keep the warning and document it.** Rejected: an operator cannot act on
  it, and an adopter learns to ignore a warning class that exists to show a
  real conflict.

## Consequences

- A copy or worktree in an ignored path, a linked worktree anywhere in the
  tree, a nested repository and a submodule no longer warn.
- An untracked nested carrier that Git ignores through `info/exclude` or the
  user's `core.excludesFile` is also left out, as `git status` leaves it out.
- An untracked directory or special file named `AGENTS.md` or `CLAUDE.md` is
  left out, because Git lists files only; a tracked path that still exists
  on disk keeps its inventory whatever kind of entry it has become.
- Every inspection that meets a nested carrier runs one more read-only Git
  command.

This decision amends ADR-0070, whose audit now covers the carriers Git lists
as the repository's own, and follows the precedent of ADR-0351, which already
asks the same Git command whether the committed tree would hold a QA Report's
link targets.
