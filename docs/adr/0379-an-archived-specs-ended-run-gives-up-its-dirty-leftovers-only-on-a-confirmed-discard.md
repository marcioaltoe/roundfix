---
status: accepted
created_at: 2026-11-09T15:00:00Z
updated_at: 2026-11-09T15:00:00Z
deprecated_at: null
superseded_by: null
---

# An archived Spec's ended Run gives up its dirty leftovers only on a confirmed discard

ADR-0053 keeps dirty work intact, and ADR-0289 releases a dirty Run of a
delivered Spec only when every uncommitted path lies inside the archived
Spec's declared and recorded paths. Its rejected option says why a path
outside that scope stays: "A stray edit outside the Spec's scope is work only
its author can judge". Nothing let that judge act. On 2026-11-09 Fiscus,
asked by the maintainer to end with only `main` and no worktree, found eleven
Runs that ended Failed or Unresolved for Specs that other Runs had delivered,
merged and archived. `roundfix reconcile` kept all eleven as `dirty`
("<path> is outside the archived Spec's declared and recorded paths"): the
leftovers were renumbered migrations, tests and components of the failed
attempt. There was no supported switch, and the Roundfix skill forbids manual
cleanup, so Fiscus ran `git reset --hard` and `git clean -fd` in each Run
Worktree by hand; the next reconcile proved 5 `safe` and 6 `superseded`.

The same report kept one Run `unknown` for ever: its recorded Run Worktree
directory was not registered in Git, held no `.git`, and held only a 46 MB
test database volume under `.storage/`, which the repository ignores. An empty
repository directory under the Worktree Root also stayed. Both were removed
with `rm -rf`. The maintainer chose, on 2026-11-09: "Spec fix: modo explícito
com confirmação (ex. --discard-dirty) só para Run terminal cuja Spec está
arquivada no default branch, gravando Branch Disposition com os paths
descartados; diretório não registrado só com artefatos ignorados vira
liberável; grupos vazios removidos."

## Decision

**A Dirty Discard is explicit, confirmed and recorded.** `roundfix reconcile`
offers a terminal Run's dirty leftovers for discard only when all of these
hold:

- the Run ended `Failed`, `Unresolved`, `Stopped`, `BudgetExceeded` or
  `TimedOut`, and no live process tree is attributed to it;
- its Spec's Archive Record is on the default branch head;
- the commit proof of ADR-0212 and ADR-0289 at the chosen head is `safe` or
  `superseded`, so every commit the Run Branch holds is already represented,
  and the only thing that keeps the Run is an uncommitted path outside the
  archived Spec's declared and recorded paths.

The dry run lists every uncommitted path of each such Run and prints one
SHA-256 digest of the whole set. `--discard-dirty <digest>` acts only when the
digest recomputed in that invocation is the confirmed one; otherwise nothing
changes. For each Run it first writes a Branch Disposition record with the
discarded paths beside the Run's artifacts, then re-proves the Run head, the
proof head and the exact path set, and only then removes the Run Worktree with
`--force` and deletes the Run Branch. Paths inside the archived Spec's scope
are discarded with the rest because ADR-0289 already proves them superseded;
committed work is never discarded, because a failed commit proof keeps the Run
out of the set. The Delivery Queue never discards.

**A recorded Run Worktree that Git does not register and that holds only
ignored files is `safe`.** The directory's name is the Run ID, the Run Branch
is absent, no registered worktree holds that branch, there is no `.git` entry
at its top, and every file below it is reported by
`git check-ignore --stdin` run in the repository, which reads the checkout's
ignore rules and never reports a tracked path. An empty directory qualifies.
`--apply` re-lists the files, refuses a changed list, and removes the
directory without a Git command; nothing is recorded in the Run Database,
because no Run Branch or head exists to record.

**`--apply` on a full scan removes empty repository directories.** A
directory directly under the configured Worktree Root whose name has the shape
Roundfix gives repository directories, `<slug>-<8 hex>`, and that holds
nothing is removed after the Run releases of the same invocation; the dry run
lists it. A Run ID limits the invocation to that Run and removes none.

## Considered Options

- Releasing every dirty Run of an archived Spec without a confirmation, as
  `--apply` does for a proof: rejected for the reason ADR-0289 gives. The
  confirmation is the judgement that rejected option asked for, and it binds
  the exact paths the operator saw.
- A confirmation by Run ID alone, as `--carry-forward` takes: rejected,
  because it does not bind the path set the dry run showed, and a full scan of
  eleven Runs would need eleven commands. The digest follows the Plan Digest
  confirmation of `roundfix baseline`.
- Cleaning the leftovers with `git reset --hard` and `git clean` and leaving
  the release to `--apply`: rejected, because the discard would then delete
  files while the Run Worktree stays, and a second command must re-prove what
  the first just proved. Removing the worktree with `--force` after the record
  is one step with one proof.
- Discarding when the commit proof fails: rejected. That deletes committed
  work the delivery does not represent; ADR-0053 keeps it.
- Proving "ignored" from the default branch's tree, or with
  `--no-index`: rejected. The checkout's rules are the ones that made the
  files ignored, and without the index a tracked path that matches a pattern
  would count as ignored.
- Removing every empty directory under the Worktree Root: rejected. A
  configured location can be a shared directory; only the repository-directory
  name shape is Roundfix's.

## Consequences

This amends ADR-0053, whose dirty work now has one explicit exit; ADR-0289,
whose rejected option stays rejected for every automatic path; and ADR-0290,
whose Run Worktree row gains the Dirty Discard as a reclaimer. Every other
reclaimer of a Run Worktree still needs a proof, and none reclaims by age. The
Delivery Queue's release after a merge reuses the terminal-Run inspection, so
it also releases an unregistered directory that holds only ignored files. The
Branch Disposition record keeps the Run head, so a discarded Run Branch stays
recoverable from Git's object store until Git prunes it.
