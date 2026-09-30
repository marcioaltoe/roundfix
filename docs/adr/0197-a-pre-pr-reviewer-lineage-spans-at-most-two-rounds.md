---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A pre-PR reviewer lineage spans at most two rounds

Each `roundfix review` opened a fresh reviewer session and reread the whole
candidate diff, even when the head had only gained the fix for the last
round's findings. The reviewer did not know what it had raised before, nor how
each finding was disposed. The maintainer's rule of 2026-09-09 allows two
review rounds per candidate, but it lived only in operating notes.

Now reviews of one candidate form a Reviewer Lineage: the same checkout,
provider and merge base, each head descending from the one reviewed before.

- Round 1 reviews the full candidate diff.
- Round 2 reviews only the diff from the round-1 head to the current head. Its
  prompt always carries the round-1 findings with their validation and
  disposition, so it never depends on the reviewer's memory.
- When round 1 leaves standing findings, the round-1 Agent Session is left open
  and round 2 continues it through acpx, on the same selection. Whether the
  runtime resumed it is measured, not assumed: the round records the ACP
  session id of each round, and `continued` is true only when both ids match.
  A runtime that cannot resume gets a new session, and round 2 still has the
  recorded conversation.
- A third review in the lineage never calls the reviewer. It closes the
  lineage as `ceiling-closed` when every standing round-2 finding has an
  operator disposition: fixed by a commit the head contains, or dismissed with
  evidence. Otherwise it prints a blocked result that names the findings to
  dispose. It keeps the round-2 record in place, so `roundfix review dispose`
  still reads it.

## Consequences

A rebase moves the merge base, so it starts a new lineage at round 1. A head
that does not descend from the recorded one also starts a new lineage, and
Roundfix closes the old lineage's open session first. The session is closed
after round 2, after any round that leaves nothing standing, and when a reuse
at the same head resolves to `findings-dismissed`.

A pre-PR review is not a Run. ADR-0018 and ADR-0051 still forbid carrying a
session across Runs, and this lineage does not touch Run sessions.

`ceiling-closed` is not a reviewer pass. It records that the change after
round 2 is the final correction, with the operator dispositions that close it.
It names the head the reviewer last saw. The Delivery Queue advances on it as
it advances on `findings-dismissed`.

This decision refines
[ADR-0153](0153-pre-pr-review-is-an-explicit-provider-policy.md),
[ADR-0018](0018-one-agent-session-per-run.md) and
[ADR-0051](0051-tasks-and-qa-own-agent-sessions.md). It supersedes none of
them.
