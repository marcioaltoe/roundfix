---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Spec authoring gets an advisory judge that never gates

The Spec Consistency Check proves citations and coverage by word overlap and
structure, so it cannot tell whether a cited ADR says what a Spec attributes
to it, or whether the TechSpec section a Coverage Map line names describes a
mechanism for the goal. On 2026-09-30 two TypeSafe Jev judgments were measured
against the Spec archive: citation support reached AUROC 0.89 and flagged 3 of
200 real citations, all three correctly, and goal to mechanism reached AUROC
0.86 and flagged 3 of 143 real mapped sections. Their recall on real defects is
modest or unmeasured, the thresholds were chosen on the data they are reported
on, and the model is a paid external service.

Roundfix therefore adds the two judgments as an advisory command,
`roundfix spec judge`, that the `write-prd` and `write-techspec` skills run
after the Spec Consistency Check:

- **Advisory, never a gate.** A raised judgment is a finding the authoring
  model answers, by correcting the artifact or by stating why the text stands.
  The command exits `0` whenever it ran, with findings or without, and no other
  command reads its result or changes its exit code because of it.
- **Fail-open.** A missing key, a network or service failure, a reached
  spending ceiling and a non-English artifact each end as a named skip, never
  as a failure and never as a clean result.
- **One pinned model version.** Requests name Jev 1.13, the version the
  thresholds were measured on, by the ID each transport accepts: `jev-1.13`
  on OpenRouter and `jev-1.13.0` on TypeSafe (ADR-0201). An answer counts
  only when the model it reports normalizes to that version:
  `jev-1.13.<patch>` from TypeSafe, or `typesafe/jev-1.13` with an optional
  `-YYYYMMDD` snapshot suffix from OpenRouter. An answer reported by any other
  model is logged and shown as skipped, but never compared with a threshold.
  Moving the pin means measuring again and recording the new thresholds.
- **Only the two measured judgments.** Task Type, the overlap shortlist and
  Task lint were measured too and are not adopted: the first is redundant, the
  second is no better than a lexical baseline, and the third is too weak.

## Consequences

- A Spec can reach decomposition with a raised judgment still standing, when
  the authoring model states why the text is right.
- A judgment that a later model version would answer differently stays
  unchanged until someone measures that version.
- The thresholds were chosen through TypeSafe directly. On 2026-10-01 the
  same 885 states asked through OpenRouter gave equivalent AUROC (0.8899
  against 0.8900 for citation support, 0.8644 against 0.8583 for goal to
  mechanism), so the pin admits OpenRouter's 1.13 snapshot.
- About 2% of advisory flags flip between two runs of the same states, on
  either transport, which is one more reason a raised judgment is answered,
  never enforced.
