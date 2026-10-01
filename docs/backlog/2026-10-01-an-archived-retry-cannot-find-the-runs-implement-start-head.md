---
type: fix
status: open
created: 2026-10-01
spec: null
reason: null
---

# An archived retry cannot find the Run's Implement start head

## Opportunity

Spec 0201 lets `deliver retry` resume an item the operator archived after an environment-only QA partial. On its first real use, Spec 0204 parked `qa-environment-partial`. The operator followed the printed answer: carry-forward, operator evidence, then `archive --qa-override`. The retry then refused with `archived retry proof: … Run "run_20261001T094428Z_f4d3047c47b08ac0" has no Implement start head for repository "/Users/marcio/dev/roundfix"`, and with `archived item head … differs from candidate head ""`. The Run was started by the queue owner built from the 0201 merge.

## Value

The class 0201 added for this exact case still ends in a manual Pull Request. The supervisor also has to run the pre-PR review by hand, because the park happens before review.

## Shape

Reproduce with a queue-started Run that parks `qa-environment-partial`, then find why the proof finds no Implement start head: the journal event it reads, the repository key, or the Run it chooses. Add an end-to-end test that drives the whole printed answer through `deliver retry` to `reviewing`.
