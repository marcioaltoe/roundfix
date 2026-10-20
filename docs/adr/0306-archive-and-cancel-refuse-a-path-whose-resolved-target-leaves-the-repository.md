---
status: accepted
created_at: 2026-10-20T21:40:00Z
updated_at: 2026-10-20T21:40:00Z
deprecated_at: null
superseded_by: null
---

# Archive and cancel refuse a path whose resolved target leaves the repository

Since 2026-10-20 `roundfix archive` and `roundfix cancel` turn every path into
a repository-relative one with one resolver: lexically first, then with both
the path and the repository root resolved. The lexical comparison accepted a
path before anything resolved it, so a symbolic link inside the repository
whose target lies outside was accepted. On 2026-10-20, in temporary
repositories, a committed `docs/history` link to another directory made both
commands write the Archive Record outside the repository and exit 0 while the
Spec folder was removed, so no commit could hold the record. A
`docs/references` link took the archive's promoted copy outside, a
`docs/backlog` link took the cancel's returned Backlog Entry outside, and a
Spec folder that was itself such a link was archived from files Git never
held. The maintainer chose "Recusar no archive/cancel (Recomendado)" (option
text: "`archive` e `cancel` recusam esse caminho com a mesma mensagem da 0288,
nomeando o caminho. Hoje ele é aceito, então a Spec declara a quebra.").

**A path is inside the repository only when its resolved target is.** Before
anything is written, both commands accept a path, and each destination they
write, only when the path resolved with `config.CanonicalizeResolvedPath` lies
inside the repository root resolved the same way. Both sides are always
resolved, so a checkout reached through a symbolic link, such as macOS `/var`
for `/private/var` or a linked `worktree.location` (ADR-0268), and an external
Spec Root whose repository is its own Git top level stay inside. A
path lexically inside whose target stays inside keeps the lexical relative
path, so every record, relink and message keeps its bytes. A path that fails
the test is refused with the text each step already printed for a path
outside the repository, naming the path: `path "<path>" is outside the
repository "<root>"`, or, from the cancel's own checks, `cancellation path
"<path>" is outside the repository`. The commands exit `2` and write nothing.
Every other reader keeps its behavior; the resolver has no caller outside the
two commands.

This narrows ADR-0269's "Symbolic-link destination parents stay accepted": a
destination parent link whose target stays in the repository, or in the
external Spec Root's repository, stays accepted; one whose target leaves it is
refused.

## Considered Options

- Refusing every path that crosses a symbolic link is what ADR-0268 rejected:
  it would refuse a checkout under `/var` on macOS, a linked Worktree Root and
  the external Spec Root of a knowledge workspace.
- Always returning the resolved relative path would change the bytes of every
  path reached through a link that stays inside, which the lexical-first rule
  of the resolver keeps.
- Writing through `os.Root`, which refuses a symbolic link that escapes its
  directory ("symbolic links may not reference a location outside the root",
  `go doc os.Root`), would refuse only at the first escaping write, after
  earlier writes, and would leave `archive --plan` and the cancel's citation
  advice, which write nothing, accepting the path. The decision asks for a
  refusal before any write.
- Refusing only at the record path would still write promoted copies and
  returned sources outside the repository.

## Consequences

A repository whose `docs/history`, `docs/references`, `docs/backlog`,
`docs/findings` or Spec folder is a symbolic link to a directory outside the
repository can no longer archive or cancel until the link is removed or
retargeted inside; in the Delivery Queue the archive step parks the item as
any failed archive does. A dangling link is resolved as a missing component,
as ADR-0268 resolves a missing worktree directory, and `os.MkdirAll` refuses to
create a directory through it, so the write fails without leaving the
repository.
