---
type: fix
status: promoted
created: 2026-10-16
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
reason: null
---

# A QA gate's archive plan contacts the judge an offline QA forbids

## Symptom

The `qa-gate` skill's Outcome follow-up asks the Agent to run
`roundfix archive <spec> --plan`. The archive plan calls the archive-value
judge, which sends each non-core candidate file, up to 6000 bytes, to the
provider whenever a judge key is in the environment. A QA Task that forbids
provider requests cannot satisfy both, and the Agent fails its own gate.

## Where

The `qa-gate` skill's Outcome follow-up, the archive plan's advisory judge, and
the Delivery Queue, which hands the operator's environment, judge keys
included, to every Run.

## Expected

The archive plan has an offline form, or the `qa-gate` follow-up skips the
advisory judge, so a QA Task that forbids provider requests can still list its
archive candidates.

## Evidence

Spec 0273's first QA attempt on 2026-10-16 recorded QA-001: two archive-value
requests through the inherited judge key (response model
`typesafe/jev-1.13-20260917`) for the Spec's two `references/` files, about
2 KB in all. All 13 runnable product rows passed. A retry with the judge keys
unset passed and the Spec merged in pull request #547.
