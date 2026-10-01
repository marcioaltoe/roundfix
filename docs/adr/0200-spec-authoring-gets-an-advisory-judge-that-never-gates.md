---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
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
- **One pinned model version.** Requests name `jev-1.13.0`, the version the
  thresholds were measured on. An answer reported by any other version is
  logged and shown, but never compared with a threshold. Moving the pin means
  measuring again and recording the new thresholds.
- **Only the two measured judgments.** Task Type, the overlap shortlist and
  Task lint were measured too and are not adopted: the first is redundant, the
  second is no better than a lexical baseline, and the third is too weak.

## Consequences

- A Spec can reach decomposition with a raised judgment still standing, when
  the authoring model states why the text is right.
- A judgment that a later model version would answer differently stays
  unchanged until someone measures that version.
