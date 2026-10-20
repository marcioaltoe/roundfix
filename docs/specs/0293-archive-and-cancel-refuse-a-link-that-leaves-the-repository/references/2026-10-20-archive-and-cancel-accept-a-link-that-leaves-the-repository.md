---
type: fix
status: promoted
created: 2026-10-20
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
reason: null
---

# Archive and cancel accept a link that leaves the repository

## Symptom

`roundfix archive` and `roundfix cancel` accept a path inside the repository
that crosses a symbolic link whose target lies outside the repository, and
write through it. On 2026-10-20, at `70c78d15` (v0.100.0), in temporary
repositories:

- with `docs/history` committed as a link to a directory outside the
  repository, `roundfix archive 0001-demo --qa-override --approval probe
  --reason probe` and `roundfix cancel 0001-demo --reason probe` each exited
  0, printed `-> docs/history/specs/0001-demo.md`, wrote the Archive Record
  into the outside directory and removed the Spec folder; `git status` showed
  only the deleted Spec files, so no commit could hold the record;
- with `docs/references` linked outside, `archive --promote
  references/notes.md` wrote the promoted copy outside;
- with `docs/backlog` linked outside, `cancel` printed `returned
  docs/backlog/2026-10-20-entry.md as open` and wrote the returned Backlog
  Entry outside;
- with the Spec folder itself a link to a directory outside, `archive`
  exited 0 and built the record from files Git never held.

## Where

`internal/spec/repository_path.go` `RepositoryPath` (lines 23-24): a path
lexically inside the repository returns before either side is resolved. The
promotion destinations in `internal/spec/archive.go` and
`internal/spec/cancel.go` and the cancel's returned-source destinations are
joined to the repository root and written without any containment check.

## Expected

Both commands refuse, before writing anything, every path whose resolved
target lies outside the repository, with the message they already print for
a path outside the repository, naming the path. A path whose target stays
inside keeps today's bytes, and a checkout or external Spec Root reached
through a symbolic link keeps working.

## Evidence

Spec 0288's pre-PR review raised it as finding F1 on 2026-10-20 and the
finding was dismissed as the pre-existing lexical behavior; the
characterization above was run for this entry. No Finding.
