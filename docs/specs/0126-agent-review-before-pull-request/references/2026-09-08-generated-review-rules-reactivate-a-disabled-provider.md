---
status: done
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
spec: 0126-agent-review-before-pull-request
---

# Baseline — generated review rules reactivate a disabled provider (2026-09-08)

The Pantheon capture reports that generated instructions requested CodeRabbit
reviews after its maintainer had explicitly disabled that provider. The same
provider-specific assumption appears in Roundfix's current generated guidance;
it must retire with the operational integration rather than requiring local
exceptions in every adopting repository.

Source: Secondbrain `inbox/roundfix/2026-09-08-baseline-exige-revisao-manual-de-provedor-desativado.md`.
The original captured evidence and its Pantheon commit remain authoritative for
that observation; no Pantheon configuration is changed by this triage.

## 1. Disabled automatic review is treated as a manual-review obligation

- Symptom / evidence: the capture names two generated clauses and Pantheon commit `08789fd4ee1172342b77aee1a6a977ff10537888`; local `docs/agents/agent-instructions.md` still carries the hand-opened review-marker obligation.
- Root cause: provider-specific operational assumptions are embedded in portable guidance without honoring an explicit provider opt-out.
- Action / suggestion: Spec 0126 retires CodeRabbit requests and matching canonical obligations, preserves required independent review where the repository approves it, and must never restore a disabled provider merely because automatic review is off. An adopting repository's provider opt-out is not permission to mutate its configuration or select a replacement on its behalf.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this evidence for [Spec 0126](../_prd.md).
The source is adopted once. Its done status means routed to implementation,
not that the provider integration has been removed or its gate has passed.

## Capture provenance after triage

The original Inbox path above records the observation at intake. The preserved
capture now lives at `inbox/roundfix/_triaged/2026-09-08-baseline-exige-revisao-manual-de-provedor-desativado.md`
in Secondbrain, with `resolved_to` naming this owned Finding.
