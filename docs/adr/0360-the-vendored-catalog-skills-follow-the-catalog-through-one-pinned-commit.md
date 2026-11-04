---
status: accepted
created_at: 2026-11-04T11:00:00Z
updated_at: 2026-11-04T11:00:00Z
deprecated_at: null
superseded_by: null
---

# The vendored catalog skills follow the catalog through one pinned commit

Amends ADR-0224 (2026-11-04): the expected upstream revision is read from
`skills-lock.json`, not from a hand-pinned constant.

On 2026-11-04 the skills catalog `marcioaltoe/skills` reported that this
repository's 33 catalog entries in `skills-lock.json` named seven different
refs: six the branch `main` and the rest commits from 25 July to 1 November.
A branch ref cannot be reproduced, and `systematic-debugging` still carried
five files the catalog's `main` (`c8b2e85`) had removed. The five upstream
Setup Snapshots pin `86df860`, and ADR-0224 holds this repository's required
skills to the `go` snapshot, so a skill cannot move without the snapshots.
Specs 0326 and 0333 each moved both by hand-run commands and re-pinned the
test constants that name the upstream revision and three skill trees. The
maintainer decided the same day: "Bump fixado e repetível (Recommended)" —
"Um alvo de Makefile move as 33 entradas para o commit atual da main do
catálogo (ref = SHA, não 'main'), reinstala e atualiza o hash; mantém o lock
reprodutível e os digests da Baseline determinísticos. Rodo o alvo agora
(c8b2e85) e a cada ciclo."

**One Catalog Revision.** Every catalog entry of `skills-lock.json` and every
upstream Setup Snapshot, with every skill it pins, name one full commit of
the catalog, never a branch. The Roundfix-owned skills stay out of the lock:
they are authored here and the catalog copies them from this repository.

**The Catalog Skill Bump is the only writer.** `make skills-third-party-bump`
resolves the catalog's `main` to its commit, or takes `SKILLS_REF`, and
fetches that one commit from `SKILLS_CATALOG`; only these two steps may use
the network, and a local clone serves them offline. It then runs the
existing writers in their order: the asset sync of the five snapshots and its
check (ADR-0191), a build that embeds them, the confirmed restore of every
catalog entry from the same checkout (ADR-0264), `make baseline-digests`
and the managed refresh. It refuses, before any write, a lock that pins a
Roundfix-owned skill, a lock and a `.agents/skills` set that differ, a
branch name and a commit the catalog does not serve. After the restore it
refuses a vendored folder that differs from the catalog's tree or a lock
entry left at another commit. At the same commit a second run changes
nothing. A bump that the sync or the suites refuse, because the catalog
renamed, added or removed a skill of a setup, stops there; that move is a
Spec. `make skills-update`, which floated the entries through another
installer, is retired.

**The tests read the records the bump writes.** The upstream revision the
setup tests expect is the one commit the lock's catalog entries share, and
the three tree tests that guarded a fix keep guarding it by refusing the
tree each fix replaced, instead of pinning the tree each sync publishes. A
documentation contract compares each vendored folder with the tree digest
the snapshots pin for it at that commit, so a local edit fails even when its
lock entry is re-hashed.

## Considered options

- Bumping only the lock and the vendored folders was rejected: ADR-0224's
  repository test and Doctor would report `systematic-debugging` as a Setup
  Snapshot Difference, and adopters would keep the old tree.
- Keeping `ref: main` on every entry was rejected: the lock would no longer
  reproduce what it installed, which is the reported defect.
- Re-pinning the revision and tree constants on every bump was rejected: the
  bump would never be repeatable without a hand edit of Go tests.

## Consequences

The bump moves Baseline content that adopters receive: after a release that
carries it, an adopter's managed refresh reports each moved skill as a Setup
Snapshot Difference, which the adopter restores deliberately. A Run without
network takes the commit from a clone the operator prepares and names in the
Task, as Specs 0326 and 0333 did.
