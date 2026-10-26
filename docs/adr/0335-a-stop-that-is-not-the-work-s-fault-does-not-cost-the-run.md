---
status: accepted
created_at: 2026-10-26T21:00:00Z
updated_at: 2026-10-26T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A stop that is not the work's fault does not cost the Run

Two rules stopped Runs whose work was sound. ADR-0276 exempts an
outside-evidence row the Spec's own authorization keeps offline only in one
exact status, `blocked (environment: network prohibited by Spec
authorization)`. Specs 0317 and 0309 recorded rows whose authorization forbade
something else, "upstream checkout and network" and "another repository", in
that same form with their own subject. Both partials were refused, both items
parked `qa-environment-partial`, and the maintainer archived both with the QA
Archive Override after the operator checked the blocked source by hand.
ADR-0274 lets the next Fallback Selection take a Task after a Provider System
Error only before the First Handoff, or in a QA Task while its seeded report
still says `verdict: pending`. On 2026-10-23 the QA turns of Specs 0300 and
0302 wrote `verdict: partial` while still checking rows, and the provider
then answered "Selected model is at capacity". The recovery was `exhausted`
although the `qa` profile named `claude/opus/high`, and both items parked
`runtime-infrastructure` until the operator retried them.

**The Network-Prohibited Row names its subject.** A Results row qualifies as
a Network-Prohibited Row when its status is
`blocked (environment: <subject> prohibited by Spec authorization)` and its
provenance names the outside-evidence row as ADR-0240 reads it. The subject
names what the authorization forbids. It is 1 to 80 characters after
trimming, counted as Unicode code points, and holds no parenthesis. The
ADR-0276 status is this form with the subject `network`, so every report that
qualified still qualifies. Each Results row counts toward one exempt kind at
most; a status in this form is not also read as a Network-Denied Row. The
fixed suffix, the outside-evidence provenance, the frontmatter counts, skipped
and finding rows, Unreachable Acceptance coverage and every refusal message
keep their bytes. The Baseline clauses that name the exemption and the QA
skills say to write the subject in the authorization's words, `network` for a
network lookup.

**A Provider System Error falls through the Fallback Chain during any Task
or QA turn.** The window of ADR-0274 is removed for Task and QA scopes: a
Provider System Error in a Task's first turn, in a Verification Feedback
repair turn after the First Handoff, or anywhere in a QA turn, moves the work
to the next selection of the Fallback Chain. The new Agent Session receives
the Task prompt, the existing Provider System Error notice and, when the
failed turn was not the first, that turn's prompt, the way a Lost Rollout
recovery resends a lost turn (ADR-0245). Before a QA fallback session starts,
the Daemon rewrites the seeded report's frontmatter `verdict:` line to
`pending` when the failed turn had written another verdict, and adds
`- discarded_verdict: <verdict>` to the `## Agent runtime fallback` section,
because a verdict written by a turn that did not finish is not a conclusion.
When the chain has no next selection, the Task still settles `failed` as
runtime infrastructure, and the Delivery Queue still parks the item
`runtime-infrastructure` with an uncounted retry. Review Batches, Lost
Rollout recovery, the Light Tier's skip, yield and escalation rules
(ADR-0238, ADR-0314) and the selection-start fallback of ADR-0114 are
unchanged.

## Considered Options

- Tighten only the QA skill's wording so agents write the network phrase:
  rejected. A row blocked because another repository is forbidden would then
  claim the network, which is not what the authorization says, and the gate
  would trade an honest record for a qualifying one.
- Read `_authorization.md` and accept only subjects it names: rejected for
  the reason ADR-0276 gave. The records state prohibitions in prose and have
  no field to match.
- Accept any status that mentions the authorization: still rejected. The
  bounded subject and the fixed prefix and suffix keep the exemption to one
  form, as the Test Anything Protocol keeps its `SKIP` directive fixed while
  its reason is free text.
- Keep the First Handoff boundary for Tasks and widen only QA: rejected. The
  same selection cannot continue after a Provider System Error, so the only
  alternative to the next selection is a park whose retry restarts the Task
  from its base. The Light Tier's escalation (ADR-0238) already lets another
  selection inherit a handed-off attempt for its repair, and a capacity reply
  is not the work's fault.
- Leave the QA report as the failed turn left it and only tell the fallback
  session: rejected. A fallback that reads a verdict may end its turn at once,
  and the Daemon would settle a verdict no finished turn wrote.

## Consequences

- ADR-0276 is amended: its "accept a host suffix or any status that mentions
  the authorization" rejection now holds only for statuses outside the
  subject form. ADR-0274's fallback window is widened to every Task and QA
  turn, and ADR-0245's rejection of a fallback after a First Handoff no longer
  applies to a Provider System Error. ADR-0114's guard against swapping models
  over partially finished work gains this exception, next to ADR-0238's.
- An agent could write a subject for a lookup it was allowed to make. As in
  ADR-0276, the report and the authorization make the claim checkable, and the
  QA skill still says to check offline every source that can be checked.
- A fallback session may repeat QA rows the failed turn had already run, and a
  Task fallback after the First Handoff starts from the tree and feedback the
  preferred selection left. Both cost turns; neither costs the Run.
