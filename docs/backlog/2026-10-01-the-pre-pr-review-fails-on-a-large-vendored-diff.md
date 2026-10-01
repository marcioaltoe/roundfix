---
type: fix
status: open
created: 2026-10-01
spec: null
reason: null
---

# The pre-PR review fails on a large vendored diff

## Opportunity

`roundfix review --base origin/main` for Spec 0200 failed three times with `review runtime failure: Agent Selection failed for runtime "codex": agent/protocol error`: twice through the Delivery Queue, once by hand. The candidate changed 72 files with 6,252 insertions, most of them copies of vendored upstream skills (`context7-cli` and three Go skills). The same binary reviewed a one-commit diff, and `codex exec` answered with the review model at the same effort. That points at the size of the review prompt the ACP path sends.

## Value

A Spec that refreshes vendored skills cannot pass its pre-PR review, so its delivery stops and someone has to review it by hand.

## Shape

Measure the prompt size at which the ACP review fails. Then either exclude vendored, regenerable paths from the reviewed diff (listing them in the record as not reviewed) or split the review into bounded parts. Record the adapter's stderr tail in the review record, so a `blocked` outcome names its cause.
