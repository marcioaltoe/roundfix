---
status: accepted
knowledge: true
created_at: 2026-11-19T15:00:00Z
updated_at: 2026-11-19T15:00:00Z
deprecated_at: null
superseded_by: null
---

# An ADR's header says whether the Secondbrain wants it and which ADRs amend it

An ADR carries two more fields in its header. The Knowledge Mark,
`knowledge: true` or `knowledge: false`, says whether the decision is
knowledge the Secondbrain wants, and only a marked ADR is published. The
Amendment List, `amended_by: [ADR-NNNN]`, names the ADRs that change part of
it. Both are written into the ADRs a repository already has, which the
lifecycle rule forbade until now.

## Context

A count of thirteen project mirrors on 2026-11-18 found 1,106 ADRs, 1,083 of
them `accepted`. Many decide how a repository is built or tested, which no
reader outside that repository uses. The project that reads the mirrors asked
whether "purely technical ADRs leave, or an ADR marks whether it interests the
knowledge base".

The same count found two ways an expired decision reads as current. Four
superseded ADRs stay in two projects although the rule deletes them. And in
this repository 64 of 364 ADRs carry a line such as "Amended by" or "Narrowed
by" in their body while their status stays `accepted`; in 63 of the 64 that
line comes after the first twelve lines of the file. A public discussion of an
ADR format with the same optional status line describes what an Agent does:
"An agent greps the ADR folder, opens a hit, reads the first few lines and
stops. If the supersede note is buried in the body, it never sees it"
(https://github.com/mattpocock/skills/discussions/847).

The maintainer answered on 2026-11-18. Which ADRs are published: "Só as
marcadas", against the author's recommendation to publish every accepted ADR.
Who marks: "Agente marca, você vê no PR (Recommended)". The ADRs that already
exist: "O Agente da sessão do update (Recommended)". An ADR that no longer
holds: "Apagar substituída, emendada no topo (Recommended)".

The docs-layout guide said until now: "Do not rewrite existing ADRs solely to
adopt lifecycle metadata." Both answers write into existing ADRs.

## Decision

**The Knowledge Mark has two values.** `knowledge: true` when the decision
states a rule of the business or of the product that someone outside the
repository's development can use. `knowledge: false` when it is about how the
repository is built, tested, organized or delivered, or about a tool it uses.
The field sits on the line after `status`.

**Only a marked ADR is published.** The Standard Export List names an ADR
only when its mark is true
([ADR-0438](0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md)).
An ADR without the mark is not exported, and the Doctor Command counts the
ADRs that lack it.

**The author marks a new ADR, and an Agent marks the existing ones once.** The
Agent of the session that applies the Baseline update that brings this rule
marks each existing ADR by the same test, in one commit that a person reviews
in a Pull Request. No command decides a mark: telling a rule of the business
from a decision about tooling needs the text read, and sending every ADR to
the judge would send a repository's decisions to an outside service.

**The mark applies where the Secondbrain is turned on.** Its rule lives in the
Secondbrain guide, and a repository without that guide marks nothing.

**The Amendment List names what changed the decision.** An accepted ADR that a
later accepted ADR changes in part, without replacing it, stays `accepted`
and lists that ADR in `amended_by`, in ascending order, on the line after
`status` or after the Knowledge Mark when it has one. Its body keeps the
dated notice after the text that changed. The author of the later ADR writes
both in the same change.

**The update fills the Amendment List from the notices.** `roundfix baseline
update` reads each ADR's body, outside code, for "Amended by", "Narrowed by",
"Extended by" and "Superseded in part by" followed by one or more ADR
identifiers, and writes the identifiers the header lacks. The Doctor Command
reports a notice the header does not list.

**A header field is added and nothing else changes.** Neither field moves
`updated_at` or any other line. An ADR without frontmatter gains a header
that holds only the field written, and it stays a legacy ADR: its status is
still read from its body.

**The lifecycle rule gains one exception.** The sentence that forbids
rewriting existing ADRs for lifecycle metadata stays, and a header field that
a guide of the repository defines for existing ADRs is the exception.

**A superseded ADR is still deleted.** The retirement rule and the history
cleanup of the update are unchanged; the Doctor Command now reports a
rejected, deprecated or superseded ADR that is still in the tree.

## Considered Options

- **Publish every accepted ADR.** The author's recommendation, because a
  decision missing from the knowledge base fails in silence. Rejected by the
  maintainer: decisions about tooling are noise there.
- **Decide the mark from the path or from keywords.** Rejected: a decision
  about a tax rule and a decision about a test helper share a folder and a
  form.
- **The judge marks inside the command.** Rejected: it sends each ADR's text
  out of the machine, in every repository.
- **A status `amended`.** Rejected: an amended ADR is still in force, and
  every reader of `status` treats anything but `accepted` as inactive.
- **Keep the notice in the body alone.** Rejected by the count above and by
  the maintainer's answer.
- **Keep a superseded ADR with its status.** Rejected by the maintainer, who
  kept the rule that deletes it.

## Consequences

- Each repository carries one change with a mark on every ADR, about 1,100
  across the projects counted, and its review is the only guard against a
  rule of the business marked as tooling.
- Between the update and the merge of that change, a repository publishes no
  ADR.
- A new ADR written without the mark is not published until someone marks it.
- The update fills the Amendment List only from the four phrases; an amendment
  described in other words is listed by hand.
- No reader of an ADR's header in Roundfix decodes it strictly, so the two
  fields break none of them.
