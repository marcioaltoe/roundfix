---
type: feat
status: promoted
created: 2026-09-28
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
reason: null
---

# An unselected review provider receives no request from Roundfix or its guidance

## Opportunity

Spec 0126 Core Feature 4: CodeRabbit is still assumed in places (watch/resolve coherence, guidance) even when the Pre-PR Review Policy selects another provider or `none`.

## Value

A repository that does not use a provider is never asked to have it, and migration never silently opts a repository into a service or into `none`.

## Shape

Remove universal CodeRabbit requirements; keep CodeRabbit fully available when selected; keep historical PR-feedback commands distinguishable from the local pre-PR interface.

Evidence: carried from Spec 0126 (`docs/history/specs/0126-agent-review-before-pull-request/_prd.md`) when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
