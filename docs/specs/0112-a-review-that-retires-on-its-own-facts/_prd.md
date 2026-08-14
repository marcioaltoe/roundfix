---
spec: 0112-a-review-that-retires-on-its-own-facts
status: active
created: 2026-08-14
surfaces: [backend, cli]
---

# A Review that retires on its own facts

Spec 0094 gave orphan Review Artifacts a terminal home and decided their
retirement by local Git reachability, so the migration could run offline. Running
that migration against the repository that shipped it produced three different
answers for the same fifty folders within one session, and nothing about the pull
requests changed between them — only the local object store did. A clone with no
pull request refs calls a head undecidable; fetching those refs calls it live,
because a squash merge leaves no ancestor; deleting the refs while the objects
remain calls it abandoned, and only then does anything migrate. Whether a review
retires is currently a property of the machine, not of the review. Separately,
the fifty folders sat at a path the resolver no longer writes to, because the
migration evaluates them for retirement and never for the rename that Spec 0094's
own resolver change implies.

## Project Constraints

- Identifier strategy: applicable — Review Artifact, Round, History Root and
  History Relocation are glossary terms this Spec changes the decision rules of,
  and a recorded settlement fact would be new vocabulary the glossary must own.
  The closing node checks whether the work introduced or changed a term. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the accurate liveness signal is the
  hosting provider's Pull Request state, which Spec 0094 deliberately refused
  because the migration must run offline. This Spec reopens that decision, so
  whether a credential is read at all, and by which command, is the central
  question rather than a detail. It adds no transport of its own and sets no
  authentication policy. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0123 is the decision under revision:
  it settles that a Review retires on local reachability and that an undecidable
  head stays live. This Spec either supersedes its reachability rule or narrows
  it, and must say which rather than leaving the record ambiguous. ADR-0121 keeps
  a relocation a ledger of identities rather than of content, which this Spec does
  not change. ADR-0120 fixes the history root that receives the retired reviews.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the spec and baseline packages plus
  their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. Whether a Review Artifact has retired is a property of the review, not of the
   machine reading it.
2. A repository that squash-merges can retire a finished review.
3. The underscored orphan review root reaches the live root without any liveness
   judgement.
4. Two clones of one repository at one commit give the same answer.

## User Stories

1. As a maintainer running the migration on a fresh clone, I want the same answer
   the migration gives on my working machine, so that the outcome does not depend
   on what I happened to fetch.
2. As a maintainer of a repository that squash-merges every pull request, I want
   finished reviews to retire, so that the feature does something rather than
   retaining everything forever.
3. As a maintainer whose repository still holds the underscored orphan review
   root, I want it renamed to the live root, so that the tool stops writing to one
   path while fifty folders sit at another.
4. As a maintainer reading a retained review, I want to know why it was retained,
   so that a permanent retention is visible rather than silent.

## Core Features

1. **Retirement rests on a recorded fact.** Whether a review is finished is
   answered from something durable about that review rather than from the local
   object store's current contents. Recording the outcome when a Round settles,
   reading the provider once for the orphan family, and treating a closed pull
   request number as terminal are the candidates the design settles between.
2. **A squash merge does not defeat the answer.** Whatever replaces ancestry
   survives a merge that discards the branch's history, because ancestry cannot
   prove integration for a repository that squash-merges.
3. **The underscored root is a legacy layout.** The migration recognises
   `docs/specs/_reviews/` as a shape to rename, independent of any judgement about
   whether the reviews inside it finished.
4. **A retained review says why, once.** A review the rule declines to retire
   reports its reason, and a permanently undecidable one does not become noise on
   every run.
5. **The relocation leaves no root behind.** The source root the relocation
   emptied is removed with the directories beneath it.

## Non-Goals / Out of Scope

- Changing where retired reviews live, or the relocation machinery itself, which
  moved 501 files with every destination verified and needs no repair.
- Changing how a Spec-owned Review Artifact is written or travels with its Spec.
- Making the migration require a network for every family; if a provider read is
  introduced, it is bounded to the orphan review family and its absence degrades
  to retention rather than to failure.
- Retroactively deciding reviews already retired.

## Success Metrics

- A fresh clone and a working checkout of one repository at one commit produce
  the same retirement decision for the same review.
- A repository that squash-merges retires at least one finished review, measured
  against a repository this Spec did not build.
- The underscored orphan root is renamed on a repository that still holds it,
  with no liveness answer involved.
- A permanently retained review reports its reason once rather than on every run.

## Decisions

- The rule is revised rather than removed. Retirement was the right goal; deriving
  it from the object store was the wrong source.

## Open Questions

- Whether the durable fact is recorded at Round settlement or read from the
  provider on demand. Recording at settlement keeps the migration offline, which
  was ADR-0123's reason, and only helps reviews created after it ships — the fifty
  that exist today would still need one provider read or a manual answer.
- Whether an undecidable head should retain or refuse. Retaining is today's
  behaviour and the safe direction; refusing would surface the ambiguity instead
  of hiding it in a warning nobody reads twice.
