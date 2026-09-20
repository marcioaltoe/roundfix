---
type: feat
status: promoted
created: 2026-09-08
spec: 0128-release-planning-with-bare-stable-tags
reason: null
---

# Plan releases whose stable tags have no v prefix

The maintainer requested implementation coverage for all remaining Findings.
The QA Rollup still carries the request in section 4 of
[the original observation](../../../history/findings/2026-08-06-the-gate-checks-that-adrs-were-cited-not-that-they-were-obeyed.md).
The current Release Plan accepts only `vMAJOR.MINOR.PATCH`; consumers using
`MAJOR.MINOR.PATCH` cannot plan their releases through the same public command.

## Intended outcome

Accept both strict stable spellings while preserving exact ref identity,
current v-prefixed behavior, semantic-version ordering and the selected base's
spelling in the proposal. An ambiguous highest base requires existing `--from`
to select the exact ref. Reset inventory keeps distinct refs distinct.

ADR-0143 keeps ordinary and reset planning read-only. This intent authorizes
neither a release nor tag creation, deletion, renaming, workflow changes or a
change to Roundfix's own publication convention. Spec 0128 owns the extension.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0128-release-planning-with-bare-stable-tags](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
