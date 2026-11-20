---
status: accepted
knowledge: false
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T19:30:00Z
deprecated_at: null
superseded_by: null
---

# A phase entry and the corrections record name their repository by its Git lineage

The identity of a repository that a phase entry records and that keys the
corrections record is the Repository Identity: the digest of the root commits
of its history, the one a Baseline Plan already carries (ADR-0071). It is
computed from Git objects alone, so it is the same in every full clone and
working tree of the repository, it survives the removal of the clone it was
computed in, and it cannot hold a credential.

## Context

ADR-0418 asks that an entry of an Authoring Phase record "an identity of the
repository that outlives the working tree it was written in" and never a
credential. An author works in a scratch clone that is removed once the plan
merges, and the operator's sessions run commands from such clones. The status
view matches an entry by asking the working tree for its `origin` remote
(ADR-0410), which stops when the tree is gone.

`internal/baseline` already derives "a clone-stable identity derived only from
Git lineage, never from a checkout path, branch, upstream, or dirty state": the
sorted root commits of `HEAD`, hashed under the schema
`roundfix/repository-identity/v1`. Measured in a scratch clone of this
repository on 2026-11-17, `git rev-list --max-parents=0 HEAD` returned one
commit in 11 ms.

Published sources disagree on the root commit as an identity, and both sides
were read through Exa on 2026-11-17. The Git documentation says that
`--max-parents=0` "gives all root commits". One tool's source keeps it because
it "survives remote renames, host migrations, and SSH/HTTPS URL changes";
another prefers the remote address because "root commits are shared across
unrelated forks that clone the same template". The Go documentation of
`net/url` warns that authentication information in a URL "has proven to be a
security risk in almost every case where it has been used".

## Decision

**One identity for both records.** A phase entry of an Authoring Phase records
the Repository Identity of the working tree it is written in, as
`repository_id`. The corrections record of a repository is kept under the
first 16 hexadecimal characters of the same digest (ADR-0420). The phase
command, the declaring commands, the Queue Owner and the report compute it the
same way, through one exported function of `internal/baseline`.

**A match needs the Spec too.** The Corrections Report takes an entry for a
Spec when the entry's identity is the checkout's and its Spec slug is the
Spec's. An entry without an identity, as every entry written before this
decision is, matches when its `repository` path is the checkout or one of its
working trees.

**Never a refusal.** When the identity cannot be computed, the phase entry is
written with an empty `repository_id`. A command that needs the corrections
record and cannot compute the identity fails with that reason.

**No remote is read.** No command of this decision asks Git for a remote
address.

## Considered Options

- **The `origin` remote, normalised to host and path.** Rejected: a repository
  without that remote, a clone made from a local path, an SSH host alias and a
  renamed repository each give another identity or none, and the value must be
  stripped of its user name and password before it is recorded.
- **The resolved path of the main checkout**, as the Artifact Directory is
  keyed. Rejected: it dies with the clone, and a declaration made from another
  clone would be written to a record the report never reads.
- **An identifier stored in the repository's Git configuration.** Rejected:
  a clone does not copy it.
- **Both the remote and the lineage.** Rejected: two rules to match by, for
  one question.

## Consequences

- A shallow clone has another identity than a full clone, because its history
  is truncated. An entry written there matches no full clone, and the report
  shows that phase as `unknown`.
- Two repositories that share their root commits, as a fork does, share the
  identity. The Spec slug still tells their entries apart unless both hold the
  same Spec.
- A merge of an unrelated history changes the identity of every branch that
  holds it. An item in delivery at that moment is declared in its own
  worktree.
- The phase command starts Git twice.
- The status view keeps its own rule (ADR-0410) within its seven days.
