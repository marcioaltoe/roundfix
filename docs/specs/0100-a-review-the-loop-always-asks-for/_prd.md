---
spec: 0100-a-review-the-loop-always-asks-for
status: active
created: 2026-08-12
surfaces: [backend, cli]
---

# A review the loop always asks for

With automatic review disabled by deliberate decision, Roundfix publishes the
review request only for heads it pushes itself. Every other head waits for
evidence nobody asked for: a head pushed by hand, the loop's own closure commit,
and a Pull Request with no findings at all — the clean case that should work
best. Two Runs on one Pull Request spent sixty minutes of wall clock and delivered
nothing while twelve Review Issues were fetchable the whole time. The tool already
owns every piece needed to avoid this: the request command, the publisher, and an
idempotence marker keyed on the head. What it lacks is the decision to use them
when the head did not come from a Round.

## Project Constraints

- Identifier strategy: applicable — Review Source, Review Issue, Round, and Final
  Push are glossary terms whose relationships this Spec changes, and the request's
  idempotence is keyed on a head identifier the glossary already names. The
  closing node checks whether the work changed a term the glossary should carry.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the work reads and publishes through the
  hosting provider's CLI, which holds the credential. This Spec introduces no new
  authentication or transport policy and adds no direct HTTP client; it changes
  when existing calls are made and how their failures are classified. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0111 makes an unobserved Verification
  unknown rather than a verdict, which is the same principle this Spec applies to
  a review signal: an absent or unrecognised review is not an approval. No
  accepted ADR governs when a review request is published, which is the gap this
  Spec fills. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the watch, review-source, and
  preflight packages plus their tests. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. Every head the loop waits on has been asked to be reviewed.
2. A refusal Roundfix can already read ends the wait in the first minute, not at
   the timeout.
3. A transient hosting-provider failure does not discard a wait that was working.
4. Two settings that cannot both be satisfied are refused before a Run starts.

## User Stories

1. As a Supervisor starting a review Run on a head I pushed by hand, I want the
   request published for that head, so that the Run works on the issues instead of
   waiting for a review nobody asked for.
2. As a Supervisor whose Pull Request has no findings at all, I want the review
   requested anyway, so that the cleanest case is not the one that times out.
3. As a Supervisor reading a Run that stopped early, I want a known refusal to be
   reported as an instruction, so that I act on it in the first minute rather than
   after thirty.
4. As a Supervisor configuring the loop, I want a ruleset requiring every thread
   resolved to be refused against a configuration that never fetches some of
   them, so that a Pull Request cannot sit blocked with no actionable signal.

## Core Features

1. **The request follows the head, not the push.** A review is requested for the
   head a Run starts from when the configured Review Source holds no accepted
   evidence for it, with idempotence proven by the existing per-head marker rather
   than by comparing against the previously seen head.
2. **A known refusal fails fast.** A review-source signal that is recognisably a
   refusal, in a configuration that can request a review, ends the wait
   immediately and reports the remedy it already knows, instead of polling to the
   timeout.
3. **A transient provider failure retries.** A hosting-provider command failure
   already carrying a transient classification does not discard the wait; the
   classification is consumed on that path rather than ignored.
4. **A coherence refusal for contradictory review settings.** A branch protection
   requiring every review thread resolved is refused against a configuration that
   excludes a severity class from fetching, naming both settings and the
   correction, in the same shape as the existing coherence refusal.

## User Experience

A Run whose head has no review evidence publishes the request and says so, once
per head. A Run that meets a known refusal ends in the first minute with the
remedy named as a command the maintainer can run. A configuration refusal names
the branch protection setting, the Roundfix setting, and which to change.

## Non-Goals / Out of Scope

- Changing what counts as a review, or weakening the rule that a green check is
  not evidence a review happened. That distinction prevented a merge on a false
  signal twice and must survive this Spec.
- Enabling automatic review in any repository, which is a maintainer decision.
- Changing the Round contract, the resolution batch shape, or the artifact layout.
- Any change to merge behavior.

## Success Metrics

- A Run started on a hand-pushed head reaches a Round rather than timing out,
  measured against a Pull Request whose twelve Review Issues were fetchable during
  a sixty-minute wait on 2026-08-10.
- A Pull Request with no findings requests its review and settles, in a repository
  this Spec did not build.
- A recognisable refusal ends a Run in under one poll interval.
- A contradictory settings pair is refused before a Run is created.
- The refusal to read a green check as review evidence still holds, proven by
  exercising both previously observed refusal bodies.

## Decisions

- Idempotence is proven by the published per-head marker rather than by comparing
  the resolved head against the current one, because the marker already answers
  the question the comparison was approximating.

## Open Questions

- Whether the request is published unconditionally at Run start or only when the
  configured Review Source holds no accepted evidence for the head. The second is
  the default until answered, because it buys no duplicate review in the common
  case where a review already ran.
- Whether the transient-failure retry has a bound, and what happens when it is
  reached. The default is the Run's existing budget, so a retry cannot outlive the
  wait it belongs to.
